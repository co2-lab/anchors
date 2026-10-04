// @anchors
//   code: DPTSD
//   ref: DUPLC

package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/shell"
)

// fakeJscpd stands in for the tool: each run appends a line to runs.txt and writes
// `report` (when not empty) as jscpd-report.json, then exits with `code`.
func fakeJscpd(t *testing.T, root, report string, code int) {
	t.Helper()
	resetDuplicationCache()
	prev := duplicationCommand
	t.Cleanup(func() { duplicationCommand = prev; resetDuplicationCache() })
	duplicationCommand = func(r, outDir string) *exec.Cmd {
		script := `echo run >> "$1/runs.txt"; echo "fake jscpd done"`
		// The report goes by FILE, not as an argument: on Windows an argument is
		// re-parsed on its way to the shell, and the `\\` of a JSON path did not survive.
		src := ""
		if report != "" {
			src = filepath.Join(t.TempDir(), "report.json")
			if err := os.WriteFile(src, []byte(report), 0o644); err != nil {
				t.Fatal(err)
			}
			script += `; cp "$3" "$2/jscpd-report.json"`
		}
		script += fmt.Sprintf("; exit %d", code)
		cmd, err := shell.Command(script, "sh", root, outDir, src)
		if err != nil {
			t.Skip(err)
		}
		cmd.Dir = r
		return cmd
	}
}

func runsOf(t *testing.T, root string) int {
	b, _ := os.ReadFile(filepath.Join(root, "runs.txt"))
	return strings.Count(string(b), "run")
}

// clonesReport: a.go:1-30 ≡ b.go:5-34, and c.go:10-30 ≡ c.go:50-70, at `pct` percent.
func clonesReport(first string, pct string) string {
	name, _ := json.Marshal(first) // escaped as jscpd writes it: a Windows path carries `\`
	return `{"duplicates":[
	  {"lines":30,"firstFile":{"name":` + string(name) + `,"start":1,"end":30},"secondFile":{"name":"pkg/b.go","start":5,"end":34}},
	  {"lines":21,"firstFile":{"name":"pkg/c.go","start":10,"end":30},"secondFile":{"name":"pkg/c.go","start":50,"end":70}}
	],"statistics":{"total":{"percentage":` + pct + `}}}`
}

func dupNode(id string) mapx.Node { return mapx.Node{ID: id, Kind: mapx.KindCode} }

func TestDuplication_FileInNoClonePasses(t *testing.T) {
	t.Run("DUPLC-B01: A file in no clone passes", func(t *testing.T) {})
	root := t.TempDir()
	fakeJscpd(t, root, clonesReport("pkg/a.go", "3.5"), 0)
	if v, msg := checkDuplication("", dupNode("pkg/d.go"), root, &mapx.Graph{}, nil); v != Pass {
		t.Fatalf("a file in no clone must pass, got %v (%s)", v, msg)
	}
}

func TestDuplication_BothSidesFailNamingTheOther(t *testing.T) {
	t.Run("DUPLC-B02: A file in a clone fails naming the other side, from either side", func(t *testing.T) {})
	root := t.TempDir()
	fakeJscpd(t, root, clonesReport("pkg/a.go", "3.5"), 0)
	g := &mapx.Graph{}
	v, msg := checkDuplication("", dupNode("pkg/a.go"), root, g, nil)
	if v != Fail || !strings.Contains(msg, ":1-30 ≡ pkg/b.go:5-34 (30 lines)") {
		t.Fatalf("the first side must fail naming the second, got %v (%s)", v, msg)
	}
	v, msg = checkDuplication("", dupNode("pkg/b.go"), root, g, nil)
	if v != Fail || !strings.Contains(msg, ":5-34 ≡ pkg/a.go:1-30 (30 lines)") {
		t.Fatalf("the second side must fail naming the first, got %v (%s)", v, msg)
	}
}

func TestDuplication_SameFileNamesOnlyRanges(t *testing.T) {
	t.Run("DUPLC-B03: A clone inside one file names only the line ranges", func(t *testing.T) {})
	root := t.TempDir()
	fakeJscpd(t, root, clonesReport("pkg/a.go", "3.5"), 0)
	v, msg := checkDuplication("", dupNode("pkg/c.go"), root, &mapx.Graph{}, nil)
	if v != Fail || !strings.Contains(msg, ":10-30 ≡ :50-70 (21 lines)") || strings.Contains(msg, "≡ pkg/c.go") {
		t.Fatalf("a same-file clone must name only its ranges, got %v (%s)", v, msg)
	}
}

func TestDuplication_ThresholdTolerates(t *testing.T) {
	t.Run("DUPLC-B04: Clones within the declared threshold are reported, over it they fail", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".jscpd.json"), []byte(`{"minLines": 20, "threshold": 3.5}`), 0o644); err != nil {
		t.Fatal(err)
	}
	fakeJscpd(t, root, clonesReport("pkg/a.go", "3.5"), 0)
	v, msg := checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil)
	if v != Diverge || !strings.Contains(msg, "pkg/b.go:5-34") {
		t.Fatalf("at the threshold the clones are reported, not failed, got %v (%s)", v, msg)
	}
	fakeJscpd(t, root, clonesReport("pkg/a.go", "3.6"), 1)
	if v, msg := checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil); v != Fail {
		t.Fatalf("over the threshold the clones fail, got %v (%s)", v, msg)
	}
}

func TestDuplication_NoThresholdAnyCloneFails(t *testing.T) {
	t.Run("DUPLC-B05: Without a declared threshold any clone fails whatever the exit code", func(t *testing.T) {})
	for _, cfg := range []string{"", `{"minLines": 20}`, `not json`} {
		root := t.TempDir()
		if cfg != "" {
			if err := os.WriteFile(filepath.Join(root, ".jscpd.json"), []byte(cfg), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		fakeJscpd(t, root, clonesReport("pkg/a.go", "0.1"), 0)
		if v, msg := checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil); v != Fail {
			t.Fatalf("config %q: exit 0 with a clone and no threshold must fail, got %v (%s)", cfg, v, msg)
		}
	}
}

func TestDuplication_OneRunPerScan(t *testing.T) {
	t.Run("DUPLC-B06: jscpd runs once per scan", func(t *testing.T) {})
	root := t.TempDir()
	fakeJscpd(t, root, clonesReport("pkg/a.go", "3.5"), 0)
	g := &mapx.Graph{}
	checkDuplication("", dupNode("pkg/a.go"), root, g, nil)
	checkDuplication("", dupNode("pkg/b.go"), root, g, nil)
	if n := runsOf(t, root); n != 1 {
		t.Fatalf("the same root and map must reuse the report, jscpd ran %d times", n)
	}
	checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil)
	if n := runsOf(t, root); n != 2 {
		t.Fatalf("a new map must run jscpd again, it ran %d times", n)
	}
}

func TestDuplication_AbsolutePathsAreRelativeToRoot(t *testing.T) {
	t.Run("DUPLC-B07: An absolute path in the report is read relative to the root", func(t *testing.T) {})
	root := t.TempDir()
	fakeJscpd(t, root, clonesReport(filepath.Join(root, "pkg", "a.go"), "3.5"), 0)
	v, msg := checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil)
	if v != Fail || !strings.Contains(msg, "pkg/b.go:5-34") {
		t.Fatalf("the absolute path must match the node, got %v (%s)", v, msg)
	}
	v, msg = checkDuplication("", dupNode("pkg/b.go"), root, &mapx.Graph{}, nil)
	if v != Fail || !strings.Contains(msg, "≡ pkg/a.go:1-30") {
		t.Fatalf("the other side must be named relative to the root, got %v (%s)", v, msg)
	}
}

func TestDuplication_NoReportIsPending(t *testing.T) {
	t.Run("DUPLC-E01: No report leaves the check Pending naming why", func(t *testing.T) {})
	root := t.TempDir()
	fakeJscpd(t, root, "", 1)
	v, msg := checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil)
	if v != Pending || !strings.Contains(msg, "fake jscpd done") {
		t.Fatalf("no report must be Pending naming jscpd's last line, got %v (%s)", v, msg)
	}
	fakeJscpd(t, root, "{not json", 0)
	v, msg = checkDuplication("", dupNode("pkg/a.go"), root, &mapx.Graph{}, nil)
	if v != Pending || !strings.Contains(msg, "unreadable") {
		t.Fatalf("an unreadable report must be Pending, got %v (%s)", v, msg)
	}
}

func TestDuplication_pinnedRelease(t *testing.T) {
	t.Run("DUPLC-B08: The gate runs a pinned jscpd release", func(t *testing.T) {})
	cmd := duplicationCommand(t.TempDir(), t.TempDir())
	if !regexp.MustCompile(`^jscpd@\d+\.\d+\.\d+$`).MatchString(cmd.Args[2]) {
		t.Errorf("the package carries an exact version, got %v", cmd.Args)
	}
}

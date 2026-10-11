// @anchors
//   code: ATTFF
//   ref: AFTSF

package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testlist"
)

// affectedProject writes the files now, and gives what they were at the last commit.
func affectedProject(t *testing.T, now map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range now {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func headOf(m map[string]string) func(string) (string, bool) {
	return func(rel string) (string, bool) { s, ok := m[rel]; return s, ok }
}

var affectedTests = []testlist.Test{
	{File: "t/e2e_test.go", Line: 10, Title: "JSRUN-B01: adds"},
	{File: "t/e2e_test.go", Line: 20, Title: "JSRUN-B02: subtracts"},
	{File: "t/unit_test.go", Line: 5, Title: "helper works"},
}

func TestAffected_rules(t *testing.T) {
	t.Run("AFTSF-B01: A changed rule reaches its tests", func(t *testing.T) {})
	before := "# Js\n\n## Effects\n\n| `JSRUN-B01` | adds |\n| `JSRUN-B02` | subtracts |\n"
	now := "# Js runtime\n\n## Effects\n\n| `JSRUN-B01` | adds two numbers |\n| `JSRUN-B02` | subtracts |\n"
	root := affectedProject(t, map[string]string{"src/js.spec.md": now})
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "src/js.spec.md", Kind: mapx.KindSpec}}}
	a := affectedBy(root, g, &config.Config{}, []string{"src/js.spec.md"}, headOf(map[string]string{"src/js.spec.md": before}))
	if got := selectAffected(a, affectedTests); len(got) != 1 || got[0].Title != "JSRUN-B01: adds" {
		t.Errorf("got %+v", got)
	}
}

func TestAffected_regions(t *testing.T) {
	t.Run("AFTSF-B02: A changed region reaches its tests, and a change outside reaches the closure", func(t *testing.T) {})
	before := "package js\n// #region [JSRUN-B01]\nfunc add() {}\n// #endregion [JSRUN-B01]\n// #region [JSRUN-B02]\nfunc sub() {}\n// #endregion [JSRUN-B02]\n"
	inside := strings.Replace(before, "func add() {}", "func add() { return }", 1)
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "src/js.go", Kind: mapx.KindCode}, {ID: "t/e2e_test.go", Kind: mapx.KindTest}, {ID: "t/unit_test.go", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{{From: "t/unit_test.go", To: "src/js.go", Type: mapx.EdgeDependsOn}},
	}
	head := headOf(map[string]string{"src/js.go": before})
	root := affectedProject(t, map[string]string{"src/js.go": inside})
	got := selectAffected(affectedBy(root, g, &config.Config{}, []string{"src/js.go"}, head), affectedTests)
	if len(got) != 1 || got[0].Title != "JSRUN-B01: adds" {
		t.Errorf("inside a region: %+v", got)
	}
	root = affectedProject(t, map[string]string{"src/js.go": strings.Replace(before, "package js", "package js // v2", 1)})
	got = selectAffected(affectedBy(root, g, &config.Config{}, []string{"src/js.go"}, head), affectedTests)
	if len(got) != 1 || got[0].Title != "helper works" {
		t.Errorf("outside the regions, the closure: %+v", got)
	}
	if got := selectAffected(affectedBy(root, g, &config.Config{}, []string{"src/js.go"}, headOf(nil)), affectedTests); len(got) != 1 || got[0].File != "t/unit_test.go" {
		t.Errorf("a new file reaches the closure: %+v", got)
	}
}

func TestAffected_testsAndFrozen(t *testing.T) {
	t.Run("AFTSF-B03: A changed test runs whole, and a frozen file reaches nothing", func(t *testing.T) {})
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "src/js.go", Kind: mapx.KindCode, Frozen: "QLDTQ-F01"}, {ID: "t/e2e_test.go", Kind: mapx.KindTest}, {ID: "t/unit_test.go", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{{From: "t/unit_test.go", To: "src/js.go", Type: mapx.EdgeDependsOn}},
	}
	root := affectedProject(t, map[string]string{"src/js.go": "x\n", "t/e2e_test.go": "y\n"})
	got := selectAffected(affectedBy(root, g, &config.Config{}, []string{"src/js.go", "t/e2e_test.go"}, headOf(map[string]string{"src/js.go": "z\n"})), affectedTests)
	if len(got) != 2 || got[0].File != "t/e2e_test.go" || got[1].File != "t/e2e_test.go" {
		t.Errorf("got %+v", got)
	}
}

func TestAffected_formats(t *testing.T) {
	t.Run("AFTSF-B04: The tests are printed for the runner", func(t *testing.T) {})
	tests := []testlist.Test{{File: "b_test.go", Line: 3, Title: "TestB"}, {File: "a_test.go", Line: 9, Title: "TestA (x+1)"}, {File: "a_test.go", Line: 2, Title: "TestA0"}}
	a := affected{whole: map[string]bool{"a_test.go": true, "b_test.go": true}, codes: map[string]bool{}}
	for format, want := range map[string]string{
		"lines": "TestA0\nTestA (x+1)\nTestB\n",
		"regex": "^(TestA0|TestA \\(x\\+1\\)|TestB)$\n",
		"files": "a_test.go\nb_test.go\n",
	} {
		var out, errw strings.Builder
		if err := printAffected(&out, &errw, a, tests, true, format); err != nil || out.String() != want {
			t.Errorf("%s: %q (%v), want %q", format, out.String(), err, want)
		}
	}
	var out, errw strings.Builder
	if err := printAffected(&out, &errw, a, tests, true, "json"); err != nil || !strings.Contains(out.String(), `"title": "TestA0"`) || strings.Index(out.String(), "TestA0") > strings.Index(out.String(), "TestB") {
		t.Errorf("json: %s", out.String())
	}
	dir := qProject(t, "version: 7\nlayers: {}\n", map[string]string{"t/x_test.go": "x\n"}, &mapx.Graph{Nodes: []mapx.Node{{ID: "t/x_test.go", Kind: mapx.KindTest}}})
	got, err := runQ(t, newAffectedTestsCmd(), "--root", dir, "t/x_test.go", "--format", "files")
	if err != nil || !strings.Contains(got, "t/x_test.go") {
		t.Errorf("the command: %v %q", err, got)
	}
}

func TestAffected_noTestsSource(t *testing.T) {
	t.Run("AFTSF-B05: With no tests source, the test files are printed", func(t *testing.T) {})
	var out, errw strings.Builder
	a := affected{whole: map[string]bool{"t/e2e.ml": true}, codes: map[string]bool{"JSRUN-B01": true}}
	if err := printAffected(&out, &errw, a, nil, false, "lines"); err != nil || out.String() != "t/e2e.ml\n" || !strings.Contains(errw.String(), "dialect.tests") {
		t.Errorf("%q %q %v", out.String(), errw.String(), err)
	}
}

func TestAffected_unknownFormat(t *testing.T) {
	t.Run("AFTSF-E01: An unknown format is refused", func(t *testing.T) {})
	var out, errw strings.Builder
	err := printAffected(&out, &errw, affected{whole: map[string]bool{}}, nil, true, "tap")
	if err == nil || !strings.Contains(err.Error(), "lines, regex, files or json") {
		t.Errorf("%v", err)
	}
}

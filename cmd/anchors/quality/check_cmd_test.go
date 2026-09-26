package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// Gates are informative on purpose: a blocking fail ends `check` in os.Exit(1), which an
// in-process test cannot survive. The verdicts are still reported and recorded.
const checkYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
  spec:
    kind: spec
    pattern: "*.spec.md"
gates:
  - name: code-ok
    on: [code]
    run: "true"
  - name: code-flagged
    on: [code]
    run: "echo flagged by the tool; exit 1"
  - name: spec-only
    on: [spec]
    when: [pre-push]
    run: "true"
  - name: never-applies
    on: [feature]
    run: "true"
`

func checkGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "a.go", Kind: mapx.KindCode, Rev: "r1"},
			{ID: "b.go", Kind: mapx.KindCode, Rev: "r2"},
			{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "s1"},
		},
		Edges: []mapx.Edge{{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies}},
	}
}

func checkFiles() map[string]string {
	return map[string]string{"a.go": "package a\n", "b.go": "package a\n", "a.spec.md": "# A\n"}
}

func TestCheckAllReportsWithoutRecording(t *testing.T) {
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	before := readQ(t, filepath.Join(dir, mapx.DefaultPath))

	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"check --all — 3 nodes, 4 gates",
		"code-ok",
		"code-flagged",
		// a declared gate that no node reached is named, not silently dropped
		"    never-applies\n",
		"output mirrored to .anchors/",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if after := readQ(t, filepath.Join(dir, mapx.DefaultPath)); after != before {
		t.Errorf("--no-record must not stamp the map")
	}
}

// The incremental check confronts the impact path of the change and stamps the edges it
// confronted.
func TestCheckChangedStampsTheImpactPath(t *testing.T) {
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())

	out, err := runQ(t, newCheckCmd(), "--root", dir, "--changed", filepath.Join(dir, "a.spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "2 nodes, 4 gates") {
		t.Errorf("the impact path of the spec is the spec and its code:\n%s", out)
	}
	if strings.Contains(out, "b.go") {
		t.Errorf("a node off the impact path is not confronted:\n%s", out)
	}
	g, err := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if st := g.Edges[0].Stamp; st == nil || st.ValidatedFromRev != "s1" || st.ValidatedToRev != "r1" {
		t.Errorf("the confronted edge must be stamped at the current revs, got %+v", st)
	}
}

func TestCheckFiltersAndRefusals(t *testing.T) {
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())

	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--category", "nothing-declares-this")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `no gate to run for this slice (phase="" category="nothing-declares-this")`) {
		t.Errorf("an empty slice says so and runs nothing:\n%s", out)
	}

	// a gate declared for pre-push is not charged at pre-commit
	out, err = runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--phase", "pre-commit")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "spec-only") || !strings.Contains(out, "3 gates") {
		t.Errorf("the pre-push gate must be out of a pre-commit run:\n%s", out)
	}

	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--skip-rule", "code-ok"); err == nil || !strings.Contains(err.Error(), "invalid --skip-rule") {
		t.Errorf("a waiver without a reason must be refused; got %v", err)
	}

	judgeOnly := qProject(t, "version: 2\nlayers: {}\ngates:\n  - name: j\n    on: [code]\n    measures: judgment\n", checkFiles(), checkGraph())
	out, err = runQ(t, newCheckCmd(), "--root", judgeOnly, "--all", "--deterministic")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no deterministic gate to run") {
		t.Errorf("--deterministic with only judgment gates runs nothing:\n%s", out)
	}

	noGates := qProject(t, "version: 2\nlayers: {}\n", checkFiles(), checkGraph())
	if _, err := runQ(t, newCheckCmd(), "--root", noGates, "--all"); err == nil || !strings.Contains(err.Error(), "no gate declared") {
		t.Errorf("a project without gates has no pipeline; got %v", err)
	}
	noMap := qProject(t, checkYAML, checkFiles(), nil)
	if _, err := runQ(t, newCheckCmd(), "--root", noMap, "--all"); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map; got %v", err)
	}
	if _, err := runQ(t, newCheckCmd(), "--root", t.TempDir(), "--all"); err == nil || !strings.Contains(err.Error(), "load config") {
		t.Errorf("no config; got %v", err)
	}
}

// The commit message waives a rule when it carries the reason, and is refused without it.
func TestCheckReadsTheWaiverFromTheCommitMessage(t *testing.T) {
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	msg := filepath.Join(dir, "MSG")
	write := func(body string) {
		if err := os.WriteFile(msg, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("feat: x\n\n[skip-code-flagged: the linter is broken upstream]\n")
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--commit-msg", msg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "flagged by the tool") || !strings.Contains(out, "3 nodes, 3 gates") {
		t.Errorf("the waived gate must not run:\n%s", out)
	}

	write("feat: x\n\n[skip-code-flagged: ]\n")
	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--commit-msg", msg); err == nil || !strings.Contains(err.Error(), "invalid --skip-rule") {
		t.Errorf("a marker without a reason must be refused; got %v", err)
	}
}

func TestNormalizeChanged(t *testing.T) {
	root := filepath.FromSlash("/repo")
	got := normalizeChanged([]string{filepath.Join(root, "src", "a.go"), "./b.go", "c/../d.go"}, root)
	if strings.Join(got, ",") != "src/a.go,b.go,d.go" {
		t.Errorf("normalizeChanged = %v", got)
	}
}

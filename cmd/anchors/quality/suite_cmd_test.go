package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// suiteGraph: a spec that specifies a.go, tested by a_test.go — the impact path of a
// change to a.go reaches the code, its test and the spec.
func suiteGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "a.go", Kind: mapx.KindCode, Rev: "r1"},
			{ID: "a_test.go", Kind: mapx.KindTest, Rev: "t1"},
			{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "s1"},
		},
		Edges: []mapx.Edge{
			{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies},
			{From: "a.go", To: "a_test.go", Type: mapx.EdgeTestedBy},
		},
	}
}

func suiteFiles() map[string]string {
	return map[string]string{"a.go": "package a\n", "a_test.go": "package a\n", "a.spec.md": "# A\n"}
}

const suiteLayers = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
`

// The declared command runs at the root and the lcov it writes reaches the map.
func TestSuiteRunsAndIngestsTheReport(t *testing.T) {
	yaml := suiteLayers + `tests:
  - workspace: backend
    layer: unit
    run: "printf 'SF:a.go\nDA:1,1\nDA:2,0\nend_of_record\n' > out/lcov.info"
    lcov: "out/lcov.info"
  - layer: e2e
    run: "echo e2e ran"
`
	files := suiteFiles()
	files["out/.keep"] = ""
	dir := qProject(t, yaml, files, suiteGraph())

	out, err := runQ(t, newTestCmd(), "--root", dir, "unit")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.Contains(out, "━━━ test [backend/unit] ━━━") {
		t.Errorf("the header names workspace and layer:\n%s", out)
	}
	if strings.Contains(out, "e2e ran") {
		t.Errorf("only the requested layer runs:\n%s", out)
	}
	g, err := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	var sig *mapx.TestSignal
	for _, n := range g.Nodes {
		if n.ID == "a.go" {
			sig = n.Signal
		}
	}
	if sig == nil || sig.TotalLines != 2 || sig.CoveredLines != 1 {
		t.Fatalf("the lcov of this run must reach a.go in the map, got %+v", sig)
	}

	// a suite without a declared report says nothing was ingested
	out, err = runQ(t, newTestCmd(), "--root", dir, "e2e")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "[e2e] passed, but the suite declares no report — nothing was ingested.") {
		t.Errorf("a run with no report must say the map got nothing:\n%s", out)
	}
}

// A failing run still ingests its report (it is exactly the run with something to say),
// and stops the chain; a report older than the run is never ingested.
func TestSuiteFailureAndStaleReport(t *testing.T) {
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "printf 'SF:a.go\nDA:1,0\nend_of_record\n' > lcov.info; exit 1"
    lcov: "lcov.info"
  - layer: e2e
    run: "echo should not run"
  - layer: old
    run: "true"
    lcov: "old.info"
`
	files := suiteFiles()
	files["old.info"] = "SF:a.go\nDA:1,1\nend_of_record\n"
	dir := qProject(t, yaml, files, suiteGraph())
	past := mustTime(t, "2020-01-01T00:00:00Z")
	if err := os.Chtimes(filepath.Join(dir, "old.info"), past, past); err != nil {
		t.Fatal(err)
	}

	out, err := runQ(t, newTestCmd(), "--root", dir, "unit", "e2e")
	if err == nil || !strings.Contains(err.Error(), `layer "unit" failed`) {
		t.Fatalf("a failing suite must fail the command; got %v", err)
	}
	if strings.Contains(out, "should not run") {
		t.Errorf("the chain stops at the first failure:\n%s", out)
	}
	g, _ := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	for _, n := range g.Nodes {
		if n.ID == "a.go" && (n.Signal == nil || n.Signal.TotalLines != 1 || n.Signal.CoveredLines != 0) {
			t.Errorf("the failing run's report must still be ingested, got %+v", n.Signal)
		}
	}

	out, err = runQ(t, newTestCmd(), "--root", dir, "old")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "report OLDER than this run") {
		t.Errorf("a report from before the run is not this run's number:\n%s", out)
	}
	g, _ = mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	for _, n := range g.Nodes {
		if n.ID == "a.go" && n.Signal.CoveredLines != 0 {
			t.Errorf("the old report (1 covered) must not overwrite the map, got %+v", n.Signal)
		}
	}
}

func TestSuiteSelectionErrors(t *testing.T) {
	yaml := suiteLayers + `tests:
  - workspace: backend
    layer: unit
    run: "true"
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())

	_, err := runQ(t, newTestCmd(), "--root", dir, "e2e")
	if err == nil || !strings.Contains(err.Error(), "not declared in `tests:`") || !strings.Contains(err.Error(), "e2e") || !strings.Contains(err.Error(), "declared layers:     unit") {
		t.Errorf("an undeclared layer lists what is declared; got %v", err)
	}
	_, err = runQ(t, newTestCmd(), "--root", dir, "unit", "-w", "backend,mobile")
	if err == nil {
		t.Error("an undeclared workspace must be refused")
	}

	// no suite declared at all: show how to declare, and fail
	bare := qProject(t, suiteLayers, suiteFiles(), suiteGraph())
	out, err := runQ(t, newMutationCmd(), "--root", bare)
	if err == nil || !strings.Contains(err.Error(), "no suite declared in `mutation:`") {
		t.Errorf("an empty section must fail; got %v", err)
	}
	if !strings.Contains(out, "No mutation suite declared") || !strings.Contains(out, "report: the mutation JSON") {
		t.Errorf("the answer shows how to declare the mutation section:\n%s", out)
	}
	out, _ = runQ(t, newTestCmd(), "--root", bare)
	if !strings.Contains(out, "junit:/lcov: the reports the run leaves behind") {
		t.Errorf("the test section asks for junit/lcov:\n%s", out)
	}

	if _, err := runQ(t, newTestCmd(), "--root", t.TempDir()); err == nil || !strings.Contains(err.Error(), "load anchors.yaml") {
		t.Errorf("no config: must fail loading it; got %v", err)
	}
}

// --changed runs the run_changed over the impact path, with absolute paths; mutation
// receives only the code, never the test.
func TestSuiteIncrementalUsesTheImpactPath(t *testing.T) {
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "echo FULL > got.txt"
    run_changed: "echo {{files}} > got.txt"
mutation:
  - layer: unit
    run: "echo FULL > mut.txt"
    run_changed: "echo {{files}} > mut.txt"
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())

	out, err := runQ(t, newTestCmd(), "--root", dir, "--changed", "a.go")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.Contains(out, "incremental: 2 file(s) on the impact path") {
		t.Errorf("code and test are on the impact path:\n%s", out)
	}
	got := readQ(t, filepath.Join(dir, "got.txt"))
	root := filepath.ToSlash(dir)
	if !strings.Contains(got, root+"/a.go") || !strings.Contains(got, root+"/a_test.go") || strings.Contains(got, "FULL") {
		t.Errorf("run_changed must receive the absolute impact files, got %q", got)
	}

	if _, err := runQ(t, newMutationCmd(), "--root", dir, "--changed", "a.go"); err != nil {
		t.Fatal(err)
	}
	mut := readQ(t, filepath.Join(dir, "mut.txt"))
	if !strings.Contains(mut, root+"/a.go") || strings.Contains(mut, "a_test.go") {
		t.Errorf("mutation mutates the code, never the test; got %q", mut)
	}
}

func TestSuiteArgvCeiling(t *testing.T) {
	t.Setenv("ANCHORS_ARGV_MAX", "10")
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "echo a command longer than ten characters"
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())
	_, err := runQ(t, newTestCmd(), "--root", dir)
	if err == nil || !strings.Contains(err.Error(), "(ceiling 10 on this platform)") {
		t.Errorf("a line over the ceiling is refused before running; got %v", err)
	}
	if n := suiteArgvLimit(); n != 10 {
		t.Errorf("suiteArgvLimit honors ANCHORS_ARGV_MAX, got %d", n)
	}
}

// --then runs the chained Anchors command only after the suite passed.
func TestSuiteThenChainsCoverage(t *testing.T) {
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "true"
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())
	out, err := runQ(t, newTestCmd(), "--root", dir, "--then", "coverage")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "━━━ then: anchors coverage ━━━") || !strings.Contains(out, "== line coverage") {
		t.Errorf("coverage should run after the suite:\n%s", out)
	}
}

func TestSuiteLabelsAndHelpers(t *testing.T) {
	if firstWord("  yarn test  ") != "yarn" || firstWord("jest") != "jest" {
		t.Error("firstWord takes the program of the command line")
	}
	if joinOrDash(nil) != "—" || joinOrDash([]string{"a", "b"}) != "a, b" {
		t.Error("joinOrDash prints a dash for an unfiltered axis")
	}
	if reportLine("mutation") == reportLine("tests") {
		t.Error("each section asks for its own report")
	}
}

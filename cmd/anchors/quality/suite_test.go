package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// A `run:` that uses {{target}} cannot run without a target. Replacing it with nothing
// would make `npx stryker run --mutate ` mutate the WHOLE PROJECT: hours of run, and none
// of what was asked. Failing here costs a second.
func TestTargetRequiredWhenDeclared(t *testing.T) {
	t.Run("STPRS-E05: A target placeholder without a target is refused", func(t *testing.T) {})
	_, err := buildCommand("npx stryker run --mutate {{target}}", "")
	if err == nil {
		t.Fatal("without --target, the command with {{target}} must fail")
	}
	if !strings.Contains(err.Error(), "--target") {
		t.Errorf("the error must say what to do; got %q", err)
	}
}

func TestTargetIsReplaced(t *testing.T) {
	t.Run("STPRS-B08: The target fills the placeholder and is ignored without one", func(t *testing.T) {})
	got, err := buildCommand("stryker run --mutate {{target}}", "business-logic/dedup.ts")
	if err != nil {
		t.Fatal(err)
	}
	if got != "stryker run --mutate business-logic/dedup.ts" {
		t.Errorf("wrong replacement: %q", got)
	}
}

// Most suites take no target (`yarn test:unit` runs everything). Requiring --target on
// them, or appending it at the end, would break the command.
func TestCommandWithoutPlaceholderIgnoresTarget(t *testing.T) {
	t.Run("STPRS-X01: The declared command runs as declared", func(t *testing.T) {})
	got, err := buildCommand("yarn test:unit", "any/thing.ts")
	if err != nil {
		t.Fatal(err)
	}
	if got != "yarn test:unit" {
		t.Errorf("a command without a placeholder must be left alone; got %q", got)
	}
}

// anchors.yaml declares a project path ("packages/backend/reports/x.json"), and the
// command may run from any directory. Resolving against the root is what makes
// `anchors mutation` work from inside a subdirectory, like the rest of the CLI.
func TestReportPathIsRelativeToTheRoot(t *testing.T) {
	t.Run("STPRS-B02: The report this run wrote is ingested into the map", func(t *testing.T) {})
	root := t.TempDir()
	got := absPath(root, "packages/backend/reports/mutation.json")
	if !filepath.IsAbs(got) {
		t.Fatalf("it should become absolute; got %q", got)
	}
	if !strings.HasSuffix(filepath.ToSlash(got), "packages/backend/reports/mutation.json") {
		t.Errorf("the path suffix was lost: %q", got)
	}
}

// Empty must stay empty. If it became the root itself, ingestion would try to read a
// directory as if it were a report.
func TestNoReportInventsNoPath(t *testing.T) {
	if got := absPath(t.TempDir(), "  "); got != "" {
		t.Errorf("an undeclared report must stay empty; got %q", got)
	}
}

// `--then` accepts only what makes sense after ingesting. Accepting any name would turn
// the flag into a disguised shell runner, which is exactly what these commands are NOT.
func TestChainRefusesAnUnknownCommand(t *testing.T) {
	t.Run("STPRS-E08: A chain naming another command is refused", func(t *testing.T) {})
	err := runChained("deploy", t.TempDir(), nil)
	if err == nil {
		t.Fatal("--then deploy should have been refused")
	}
	if !strings.Contains(err.Error(), "check") {
		t.Errorf("the error must say what IS accepted; got %q", err)
	}
}

// The opt-in is the default. Without `--then`, it runs and stops.
func TestEmptyChainDoesNothing(t *testing.T) {
	t.Run("STPRS-B09: A passing run chains coverage, and an empty chain does nothing", func(t *testing.T) {})
	if err := runChained("", t.TempDir(), nil); err != nil {
		t.Errorf("without --then there is nothing to chain; got %v", err)
	}
	if err := runChained("  ,  ", t.TempDir(), nil); err != nil {
		t.Errorf("empty separators are not a command; got %v", err)
	}
}

// Declaring `run:` without `junit:` is legitimate (a shortcut) and cannot be treated as
// an error. The command says nothing was ingested, because silence here would make the
// user think the gate was about to change colour.
func TestSuiteThatPassesWithoutReportDoesNotBreak(t *testing.T) {
	t.Run("STPRS-B01: The selected suite runs at the root under a header naming it", func(t *testing.T) {})
	root := t.TempDir()
	mark := filepath.Join(root, "ran.txt")
	cs := suiteCommand{nome: "test", secao: "tests"}
	s := []config.Suite{{Layer: "unit", Run: "printf ok > " + filepath.ToSlash(mark)}}

	if err := runSuites(cs, s, root, "", nil); err != nil {
		t.Fatalf("a suite without a report should not fail: %v", err)
	}
	if _, err := os.Stat(mark); err != nil {
		t.Errorf("the declared command did not run: %v", err)
	}
}

// Stopping at the first failure is the rule: the layers depend on each other, and running
// e2e over a red unit only produces noise over a base already broken.
func TestFailingSuiteStopsTheRest(t *testing.T) {
	t.Run("STPRS-B04: A failing suite is still ingested and stops the rest", func(t *testing.T) {})
	root := t.TempDir()
	after := filepath.Join(root, "should-not-exist.txt")
	cs := suiteCommand{nome: "test", secao: "tests"}
	s := []config.Suite{
		{Layer: "unit", Run: "exit 3"},
		{Layer: "e2e", Run: "printf x > " + filepath.ToSlash(after)},
	}

	err := runSuites(cs, s, root, "", nil)
	if err == nil {
		t.Fatal("the unit failure should have stopped the run")
	}
	if !strings.Contains(err.Error(), "unit") {
		t.Errorf("the error must name the layer that failed; got %q", err)
	}
	if _, serr := os.Stat(after); serr == nil {
		t.Error("the next layer must not have run")
	}
}

// ── the two modes: full and incremental ──────────────────────────────────────

// The default is the full run, as in `check`.
func TestWithoutChangedTheFullCommandRuns(t *testing.T) {
	t.Run("STPRS-B05: Without changed files the full command runs", func(t *testing.T) {})
	s := config.Suite{Run: "jest", RunChanged: "jest --findRelatedTests {{files}}"}
	got, err := pickCommand(s, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "jest" {
		t.Errorf("without --changed the full command runs; got %q", got)
	}
}

// And the files go where the project declared them, not appended at the end: the position
// of the cut belongs to the command, not to us.
func TestChangedUsesTheIncrementalCommand(t *testing.T) {
	t.Run("STPRS-B06: The incremental command receives the impact path where it declares it", func(t *testing.T) {})
	s := config.Suite{Run: "jest", RunChanged: "jest --findRelatedTests {{files}} --ci"}
	got, err := pickCommand(s, []string{"a/x.ts", "b/y.ts"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "jest --findRelatedTests a/x.ts b/y.ts --ci" {
		t.Errorf("wrong {{files}} replacement: %q", got)
	}
	abs := impactFiles([]mapx.Node{{ID: "packages/backend/x.ts", Kind: mapx.KindCode}}, "tests", t.TempDir())
	if len(abs) != 1 || !strings.HasPrefix(abs[0], "/") && !strings.Contains(abs[0], ":/") || strings.Contains(abs[0], `\`) {
		t.Errorf("an impact file is passed absolute, with forward slashes: %v", abs)
	}
}

// The most important decision of the two modes: falling back to the FULL run would be
// expensive and, worse, a lie — the user would read "passed" thinking their cut ran.
func TestChangedWithoutRunChangedRefuses(t *testing.T) {
	t.Run("STPRS-E06: The incremental mode without an incremental command is refused", func(t *testing.T) {})
	s := config.Suite{Run: "yarn test:unit"}
	_, err := pickCommand(s, []string{"a/x.ts"}, "")
	if err == nil {
		t.Fatal("without run_changed, the incremental mode must refuse")
	}
	if !strings.Contains(err.Error(), "run_changed") || !strings.Contains(err.Error(), "{{files}}") {
		t.Errorf("the error must show what to declare; got %q", err)
	}
}

// Mutation changes the RULE. Mutating the test would invert the experiment: the test is
// the measuring instrument, not the object measured.
func TestMutationDoesNotMutateTheTest(t *testing.T) {
	t.Run("STPRS-B07: Tests get code and tests, mutation gets only code", func(t *testing.T) {})
	nodes := []mapx.Node{
		{ID: "a/rule.ts", Kind: mapx.KindCode},
		{ID: "a/rule.test.ts", Kind: mapx.KindTest},
		{ID: "a/rule.spec.md", Kind: mapx.KindSpec},
		{ID: "a/rule.feature", Kind: mapx.KindFeature},
	}
	if got := impactFiles(nodes, "mutation", ""); len(got) != 1 || got[0] != "a/rule.ts" {
		t.Errorf("mutation receives only code; got %v", got)
	}
}

// `--findRelatedTests` and its equivalents expect source files; a spec or a feature means
// nothing to a runner.
func TestTestReceivesCodeAndTest(t *testing.T) {
	t.Run("STPRS-B07: Tests get code and tests, mutation gets only code", func(t *testing.T) {})
	nodes := []mapx.Node{
		{ID: "a/rule.ts", Kind: mapx.KindCode},
		{ID: "a/rule.test.ts", Kind: mapx.KindTest},
		{ID: "a/rule.spec.md", Kind: mapx.KindSpec},
	}
	got := impactFiles(nodes, "tests", "")
	if len(got) != 2 || got[0] != "a/rule.ts" || got[1] != "a/rule.test.ts" {
		t.Errorf("test receives code and test, no spec; got %v", got)
	}
}

// The command runs inside `sh -c`, where the backslash is an ESCAPE. Measured: with the
// native separator the path reaches the runner mangled and the run finds 0 tests — no
// error, just an empty cut that reads as "there was nothing to run". It is the most
// expensive failure mode: silent and optimistic.
func TestImpactPathCarriesNoBackslash(t *testing.T) {
	t.Run("STPRS-B06: The incremental command receives the impact path where it declares it", func(t *testing.T) {})
	root := t.TempDir()
	got := impactFiles([]mapx.Node{{ID: "packages/backend/x.ts", Kind: mapx.KindCode}}, "tests", root)
	if len(got) != 1 {
		t.Fatalf("expected 1 file; got %v", got)
	}
	if strings.Contains(got[0], `\`) {
		t.Errorf("a path passed to the shell cannot carry a backslash: %q", got[0])
	}
	if !strings.HasSuffix(got[0], "packages/backend/x.ts") {
		t.Errorf("the path lost its suffix: %q", got[0])
	}
}

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
	t.Run("STPRS-B01: The selected suite runs at the root under a header naming it", func(t *testing.T) {})
	t.Run("STPRS-B02: The report this run wrote is ingested into the map", func(t *testing.T) {})
	t.Run("STPRS-B03: A passing suite with no report says nothing was ingested", func(t *testing.T) {})
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
	t.Run("STPRS-B04: A failing suite is still ingested and stops the rest", func(t *testing.T) {})
	t.Run("STPRS-I01: A report older than the run is never ingested", func(t *testing.T) {})
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
	t.Run("STPRS-E01: The suite commands without configuration fail", func(t *testing.T) {})
	t.Run("STPRS-E02: A section with no suite shows how to declare it and fails", func(t *testing.T) {})
	t.Run("STPRS-E03: A filter that names nothing declared is refused with what is declared", func(t *testing.T) {})
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
	t.Run("STPRS-B07: Tests get code and tests, mutation gets only code", func(t *testing.T) {})
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
	t.Run("STPRS-E07: A command line over the ceiling is refused before running", func(t *testing.T) {})
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
	t.Run("STPRS-B09: A passing run chains coverage, and an empty chain does nothing", func(t *testing.T) {})
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
	for _, empty := range []string{"", "  ,  "} {
		if err := runChained(empty, dir, nil); err != nil {
			t.Errorf("an empty chain %q runs nothing; got %v", empty, err)
		}
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

func TestSuiteIncrementalWithNoCodeOnThePathRunsNothing(t *testing.T) {
	t.Run("STPRS-B10: An impact path with no code file runs nothing", func(t *testing.T) {})
	yaml := suiteLayers + `mutation:
  - layer: unit
    run: "echo FULL > mut.txt"
    run_changed: "echo {{files}} > mut.txt"
`
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "lone.spec.md", Kind: mapx.KindSpec, Rev: "s1"}}}
	dir := qProject(t, yaml, map[string]string{"lone.spec.md": "# L\n"}, g)
	out, err := runQ(t, newMutationCmd(), "--root", dir, "--changed", "lone.spec.md")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.Contains(out, "the impact path reaches no code file — nothing to run.") {
		t.Errorf("a path with no code has nothing to mutate:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, "mut.txt")); err == nil {
		t.Error("no suite may run when the impact path has no code")
	}
}

func TestSuiteValidFiltersMatchingNothingFail(t *testing.T) {
	t.Run("STPRS-E04: Declared filters that match no suite together are refused", func(t *testing.T) {})
	yaml := suiteLayers + `tests:
  - workspace: backend
    layer: unit
    run: "true"
  - workspace: mobile
    layer: e2e
    run: "true"
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())
	_, err := runQ(t, newTestCmd(), "--root", dir, "unit", "-w", "mobile")
	if err == nil || !strings.Contains(err.Error(), "no suite matches layer(s) unit, workspace(s) mobile") {
		t.Errorf("a declared layer and workspace that never meet must be refused; got %v", err)
	}
}

func TestSuiteChangedWithoutMapFails(t *testing.T) {
	t.Run("STPRS-E09: The incremental mode without a map points at the map build", func(t *testing.T) {})
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "true"
    run_changed: "echo {{files}}"
`
	dir := qProject(t, yaml, suiteFiles(), nil)
	_, err := runQ(t, newTestCmd(), "--root", dir, "--changed", "a.go")
	if err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("--changed needs the map; got %v", err)
	}
}

// `--then check` hands the check the suite's own scope: the full sweep after a full run,
// the same changed files after an incremental one. A chained check with neither refused
// ("provide --changed or --all") every time.
func TestSuiteThenChainsCheckWithTheSuitesScope(t *testing.T) {
	t.Run("STPRS-B11: A passing run chains the check over the suite's own scope", func(t *testing.T) {})
	englishOutput(t)
	yaml := suiteLayers + `gates:
  - name: always-fine
    on: [code]
    run: "true"
tests:
  - layer: unit
    run: "true"
    run_changed: "true"
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())
	out, err := runQ(t, newTestCmd(), "--root", dir, "--then", "check")
	if err != nil {
		t.Fatalf("a full run chaining check must pass: %v\n%s", err, out)
	}
	if !strings.Contains(out, "━━━ then: anchors check ━━━") || !strings.Contains(out, "check --all — ") {
		t.Errorf("the chained check of a full run is the full sweep:\n%s", out)
	}
	out, err = runQ(t, newTestCmd(), "--root", dir, "--changed", "a.go", "--then", "check")
	if err != nil {
		t.Fatalf("an incremental run chaining check must pass: %v\n%s", err, out)
	}
	if !strings.Contains(out, "━━━ then: anchors check ━━━") || strings.Contains(out, "check --all — ") {
		t.Errorf("the chained check of an incremental run is incremental:\n%s", out)
	}
}

func TestSuiteIncrementalRespectsSuitePaths(t *testing.T) {
	t.Run("STPRS-B12: An incremental run hands each suite only its own impact files", func(t *testing.T) {})
	yaml := suiteLayers + `tests:
  - workspace: here
    layer: unit
    run: "echo FULL"
    run_changed: "echo {{files}} > here.txt"
    paths: ["*.go"]
  - workspace: there
    layer: unit
    run: "echo FULL"
    run_changed: "echo {{files}} > there.txt"
    paths: ["web/**"]
`
	dir := qProject(t, yaml, suiteFiles(), suiteGraph())
	out, err := runQ(t, newTestCmd(), "--root", dir, "--changed", "a.go")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := readQ(t, filepath.Join(dir, "here.txt")); !strings.Contains(got, "a.go") {
		t.Errorf("the suite whose paths cover the change receives it, got %q", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "there.txt")); err == nil || !strings.Contains(out, "[there/unit] the impact path reaches none of this suite's") {
		t.Errorf("the other suite runs nothing and says why:\n%s", out)
	}
}

// On Linux a file's mtime comes from a coarse clock, a tick behind `time.Now()`: a report
// written in the run's first milliseconds carried a time BEFORE the start, and was dropped
// as older — in CI, one of the two suite tests above failed on most pushes.
func TestIngestIfRecentToleratesACoarseClock(t *testing.T) {
	t.Run("STPRS-B13: A report stamped by a coarse clock just before the start is this run's", func(t *testing.T) {})
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "true"
    lcov: "lcov.info"
`
	files := suiteFiles()
	files["lcov.info"] = "SF:a.go\nDA:1,1\nDA:2,0\nend_of_record\n"
	dir := qProject(t, yaml, files, suiteGraph())
	start := time.Now()
	just := start.Add(-500 * time.Millisecond)
	if err := os.Chtimes(filepath.Join(dir, "lcov.info"), just, just); err != nil {
		t.Fatal(err)
	}
	s := config.Suite{Layer: "unit"}
	if err := ingestIfRecent(dir, "", filepath.Join(dir, "lcov.info"), "", s, start, false); err != nil {
		t.Fatal(err)
	}
	g, _ := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	for _, n := range g.Nodes {
		if n.ID == "a.go" && (n.Signal == nil || n.Signal.CoverageBySuite["lcov.info"].TotalLines != 2) {
			t.Fatalf("a report half a second before the start is this run's, got %+v", n.Signal)
		}
	}
	old := start.Add(-10 * time.Second)
	if err := os.Chtimes(filepath.Join(dir, "lcov.info"), old, old); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() { _ = ingestIfRecent(dir, "", filepath.Join(dir, "lcov.info"), "", s, start, false) })
	if !strings.Contains(out, "report OLDER than this run") {
		t.Errorf("a report ten seconds before the start is still older, got %q", out)
	}
}

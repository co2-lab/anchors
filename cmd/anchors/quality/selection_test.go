package quality

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// testNode is a test file with a result: passed/failed cases, measured at rev atRev.
func testNode(id string, passed, failed int, fresh bool) mapx.Node {
	at := "r-" + id
	if !fresh {
		at = "old"
	}
	return mapx.Node{ID: id, Kind: mapx.KindTest, Rev: "r-" + id,
		Signal: &mapx.TestSignal{Passed: passed, Failed: failed, AtRev: at}}
}

func TestSelection_FourBoxes(t *testing.T) {
	t.Run("SLCTN-B01: A file is placed in one of four boxes", func(t *testing.T) {})
	g := &mapx.Graph{}
	for _, c := range []struct {
		n    mapx.Node
		want fileState
	}{
		{testNode("a", 3, 0, true), freshPassing},
		{testNode("b", 3, 0, false), stalePassing},
		{testNode("c", 2, 1, true), freshBelow},
		{testNode("d", 2, 1, false), staleBelow},
		{mapx.Node{ID: "e", Kind: mapx.KindTest}, unmeasured},
		{mapx.Node{ID: "f", Kind: mapx.KindTest, Signal: &mapx.TestSignal{SecondsBySuite: map[string]float64{"k": 1}}}, unmeasured},
	} {
		if got := testState(g, c.n, "unit"); got != c.want {
			t.Errorf("%s: want %v, got %v", c.n.ID, c.want, got)
		}
	}
}

func TestSelection_DefaultAndFlags(t *testing.T) {
	t.Run("SLCTN-B02: The default takes stale below the minimum and never measured, and each flag opens a side", func(t *testing.T) {})
	all := []fileState{freshPassing, stalePassing, freshBelow, staleBelow, unmeasured}
	takes := func(sel runSelection) []fileState {
		var out []fileState
		for _, st := range all {
			if sel.wants(st) {
				out = append(out, st)
			}
		}
		return out
	}
	for _, c := range []struct {
		sel  runSelection
		want []fileState
	}{
		{runSelection{}, []fileState{staleBelow, unmeasured}},
		{runSelection{IncludeFresh: true}, []fileState{freshBelow, staleBelow, unmeasured}},
		{runSelection{IncludePassing: true}, []fileState{stalePassing, staleBelow, unmeasured}},
		{runSelection{IncludeFresh: true, IncludePassing: true}, all},
		{runSelection{SkipUnmeasured: true}, []fileState{staleBelow}},
	} {
		if got := takes(c.sel); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v: want %v, got %v", c.sel, c.want, got)
		}
	}
}

func TestSelection_TestState(t *testing.T) {
	t.Run("SLCTN-B03: A test file is stale through what it exercises and passes by its own layer", func(t *testing.T) {})
	exercised := testNode("a_test.go", 2, 0, true)
	exercised.Signal.ClosureRev = map[string]string{"a.go": "before"}
	g := &mapx.Graph{Nodes: []mapx.Node{exercised, {ID: "a.go", Kind: mapx.KindCode, Rev: "after"}}}
	if got := testState(g, exercised, "unit"); got != stalePassing {
		t.Errorf("a test whose exercised code changed is stale, got %v", got)
	}
	byLayer := testNode("b_test.go", 5, 1, true)
	byLayer.Signal.ByLayer = map[string]mapx.LayerExec{"unit": {Passed: 4}, "integration": {Passed: 1, Failed: 1}}
	if got := testState(&mapx.Graph{}, byLayer, "unit"); got != freshPassing {
		t.Errorf("in its unit layer the file passes, got %v", got)
	}
	if got := testState(&mapx.Graph{}, byLayer, "integration"); got != freshBelow {
		t.Errorf("in its integration layer the file fails, got %v", got)
	}
	skipped := mapx.Node{ID: "c", Kind: mapx.KindTest, Rev: "r", Signal: &mapx.TestSignal{Skipped: 3, AtRev: "r"}}
	if got := testState(&mapx.Graph{}, skipped, "unit"); got != freshBelow {
		t.Errorf("a file whose cases were all skipped is not passing, got %v", got)
	}
}

func codeNode(id string, killed, survived, timedOut int, low float64, fresh bool) mapx.Node {
	at := "r-" + id
	if !fresh {
		at = "old"
	}
	score := 0.0
	if killed+survived > 0 {
		score = 100 * float64(killed) / float64(killed+survived)
	}
	return mapx.Node{ID: id, Kind: mapx.KindCode, Rev: "r-" + id, Signal: &mapx.TestSignal{AtRev: at,
		MutantsKilled: killed, MutantsSurvived: survived, MutantsTimedOut: timedOut, MutationScore: score, MutationLow: low}}
}

func TestSelection_MutationState(t *testing.T) {
	t.Run("SLCTN-B04: A mutation result under load counts as stale and passes at the floor", func(t *testing.T) {})
	for _, c := range []struct {
		n    mapx.Node
		want fileState
	}{
		{codeNode("load", 74, 4, 65, 0, true), stalePassing},
		{codeNode("atceiling", 100, 0, 20, 0, true), freshPassing}, // exactly 20% timed out
		{codeNode("overceiling", 100, 0, 21, 0, true), stalePassing},
		{codeNode("edge", 21, 9, 0, 0, true), freshPassing}, // exactly 70%
		{codeNode("low", 20, 10, 0, 0, true), freshBelow},   // 66.7%
		{codeNode("floor", 6, 4, 0, 60, true), freshPassing},
		{codeNode("stale", 1, 9, 0, 0, false), staleBelow},
		{mapx.Node{ID: "none", Kind: mapx.KindCode, Rev: "r", Signal: &mapx.TestSignal{AtRev: "r", MutantsNoCoverage: 3}}, freshPassing},
		{mapx.Node{ID: "never", Kind: mapx.KindCode}, unmeasured},
	} {
		if got := mutationState(c.n, config.DefaultTimeoutCeiling); got != c.want {
			t.Errorf("%s: want %v, got %v", c.n.ID, c.want, got)
		}
	}
}

func TestSelection_Ownership(t *testing.T) {
	t.Run("SLCTN-B05: A test file belongs to the suite that ran it", func(t *testing.T) {})
	timedHere := mapx.Node{ID: "a", Signal: &mapx.TestSignal{SecondsBySuite: map[string]float64{"out/unit.xml": 1}}}
	timedElsewhere := mapx.Node{ID: "b", Signal: &mapx.TestSignal{SecondsBySuite: map[string]float64{"web/unit.xml": 1}}}
	layerHere := mapx.Node{ID: "c", Signal: &mapx.TestSignal{ByLayer: map[string]mapx.LayerExec{"unit": {Passed: 1}}}}
	layerElsewhere := mapx.Node{ID: "d", Signal: &mapx.TestSignal{ByLayer: map[string]mapx.LayerExec{"e2e": {Passed: 1}}}}
	for n, want := range map[*mapx.Node]bool{&timedHere: true, &timedElsewhere: false, &layerHere: true, &layerElsewhere: false, {ID: "e"}: true} {
		if got := ownedBy(*n, "out/unit.xml", "unit"); got != want {
			t.Errorf("%s: want %v, got %v", n.ID, want, got)
		}
	}
}

func TestSelection_NeverSupportOrNoSignal(t *testing.T) {
	t.Run("SLCTN-B06: Support files and no_signal targets never run", func(t *testing.T) {})
	support := testNode("helpers_test.go", 0, 1, false)
	support.Support = true
	g := &mapx.Graph{Nodes: []mapx.Node{support, testNode("a_test.go", 0, 1, false),
		codeNode("gen.go", 0, 5, 0, 0, false), codeNode("a.go", 0, 5, 0, 0, false)}}
	cfg := &config.Config{Gates: []config.Gate{{Name: "mutation-score", Check: "mutation-score", NoSignal: map[string]string{"gen.go": "generated"}}}}
	if run, _ := selectFiles(g, cfg, false, "k", "unit", runSelection{}); !reflect.DeepEqual(run, []string{"a_test.go"}) {
		t.Errorf("the support file is never run, got %v", run)
	}
	if run, _ := selectFiles(g, cfg, true, "k", "unit", runSelection{}); !reflect.DeepEqual(run, []string{"a.go"}) {
		t.Errorf("the no_signal target is never run, got %v", run)
	}
	// A test file another suite owns is not this suite's to run.
	other := testNode("web_test.go", 0, 1, false)
	other.Signal.SecondsBySuite = map[string]float64{"web/junit.xml": 1}
	g.Nodes = append(g.Nodes, other)
	if run, _ := selectFiles(g, cfg, false, "k", "unit", runSelection{}); !reflect.DeepEqual(run, []string{"a_test.go"}) {
		t.Errorf("another suite's test file is left to it, got %v", run)
	}
}

func TestSelection_ArgvBatches(t *testing.T) {
	t.Run("SLCTN-B10: The selected files run in as few batches as the ceiling allows", func(t *testing.T) {})
	s := config.Suite{RunChanged: "run {{files}}"} // 4 characters without the placeholder
	got := argvBatches(s, []string{"aaaa", "bbbb", "cccc", strings.Repeat("x", 40), "dd"}, 20)
	want := [][]string{{"aaaa", "bbbb", "cccc"}, {strings.Repeat("x", 40)}, {"dd"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	// A first file already over the ceiling goes alone, with no empty batch before it.
	got = argvBatches(s, []string{strings.Repeat("y", 40), "dd"}, 20)
	if want := [][]string{{strings.Repeat("y", 40)}, {"dd"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

// selectionProject: a suite with run_changed over four test files, one in each box.
func selectionProject(t *testing.T, runChanged string) string {
	t.Helper()
	g := &mapx.Graph{Nodes: []mapx.Node{
		testNode("fp_test.go", 1, 0, true), testNode("sp_test.go", 1, 0, false),
		testNode("fb_test.go", 0, 1, true), testNode("sb_test.go", 0, 1, false),
		{ID: "new_test.go", Kind: mapx.KindTest, Rev: "r"},
	}}
	return budgetProject(t, "tests", runChanged, g)
}

func TestSelection_RunSaysWhatItTookAndLeft(t *testing.T) {
	t.Run("SLCTN-B07: The run says what it selected and what it left out", func(t *testing.T) {})
	dir := selectionProject(t, "sh run.sh {{files}}")
	out, err := runQ(t, newTestCmd(), "--root", dir)
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := readLog(t, dir); !reflect.DeepEqual(got, []string{"new_test.go sb_test.go"}) {
		t.Fatalf("by default only the stale failing and the never measured run, in one batch, got %v", got)
	}
	for _, want := range []string{
		"[unit] selected 2 file(s) to run",
		"1 fresh and passing (--include-fresh --include-passing)",
		"1 stale but passing (--include-passing)",
		"1 fresh and below the minimum (--include-fresh)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestSelection_SuiteWithoutRunChanged(t *testing.T) {
	t.Run("SLCTN-B08: A suite without run_changed runs whole", func(t *testing.T) {})
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "echo whole-run"
    junit: "out/junit.xml"
`
	dir := qProject(t, yaml, map[string]string{"a_test.go": "x", "out/.keep": ""}, &mapx.Graph{Nodes: []mapx.Node{testNode("a_test.go", 0, 1, false)}})
	out, err := runQ(t, newTestCmd(), "--root", dir)
	if err != nil || !strings.Contains(out, "cannot run a subset, and runs whole") || !strings.Contains(out, "whole-run") {
		t.Fatalf("the suite runs whole saying why, got %v\n%s", err, out)
	}
	if _, err := runQ(t, newTestCmd(), "--root", dir, "--budget", "60s"); err == nil || !strings.Contains(err.Error(), "run_changed") {
		t.Fatalf("with a budget it is refused, got %v", err)
	}
}

func TestSelection_NothingToRun(t *testing.T) {
	t.Run("SLCTN-B09: A selection that takes nothing runs nothing", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{testNode("a_test.go", 1, 0, true)}}
	dir := budgetProject(t, "tests", "sh run.sh {{files}}", g)
	out, err := runQ(t, newTestCmd(), "--root", dir)
	if err != nil || !strings.Contains(out, "nothing to run") {
		t.Fatalf("a selection that takes nothing says so, got %v\n%s", err, out)
	}
	if got := readLog(t, dir); len(got) != 0 {
		t.Fatalf("no command runs, got %v", got)
	}
}

func TestSelection_AllAndConflicts(t *testing.T) {
	t.Run("SLCTN-B11: --all runs whole and does not combine with the other choices", func(t *testing.T) {})
	dir := selectionProject(t, "sh run.sh {{files}}")
	out, err := runQ(t, newTestCmd(), "--root", dir, "--all")
	if err != nil || !strings.Contains(out, "echo full") {
		t.Fatalf("--all runs the whole command, got %v\n%s", err, out)
	}
	if got := readLog(t, dir); len(got) != 0 {
		t.Fatalf("--all does not go through run_changed, got %v", got)
	}
	for _, args := range [][]string{
		{"--all", "--changed", "a_test.go"},
		{"--all", "--include-fresh"},
		{"--include-passing", "--changed", "a_test.go"},
	} {
		if _, err := runQ(t, newTestCmd(), append([]string{"--root", dir}, args...)...); err == nil {
			t.Errorf("%v must be refused", args)
		}
	}
}

func TestSelection_NoMap(t *testing.T) {
	t.Run("SLCTN-E01: A selective run without a map is refused", func(t *testing.T) {})
	dir := selectionProject(t, "sh run.sh {{files}}")
	if err := os.Remove(filepath.Join(dir, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	if _, err := runQ(t, newTestCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "--all") {
		t.Fatalf("refused saying to build the map or use --all, got %v", err)
	}
}

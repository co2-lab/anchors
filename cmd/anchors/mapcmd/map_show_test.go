package mapcmd

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// The neighbourhood of the code node: governed by the guide and specified by the spec
// (going up), and nothing below it — a leaf.
func TestMapShow_neighbourhoodOfANode(t *testing.T) {
	root := fixtureProject(t)
	out := runCmd(t, newMapCmd(), "show", "src/login.ts", "--root", root)

	for _, want := range []string{
		"↑ governed by / validates against (2):",
		"src/login.ts  ←governs←  guides/CODE.md",
		"src/login.ts  ←specifies←  src/login.spec.md",
		"↓ propagates to (0):",
		"(none — leaf)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	spec := runCmd(t, newMapCmd(), "show", "src/login.spec.md", "--root", root)
	for _, want := range []string{"(none — top/root)", "src/login.spec.md  —specifies→  src/login.ts", "src/login.spec.md  —tested-by→  src/login.test.ts"} {
		if !strings.Contains(spec, want) {
			t.Errorf("missing %q in the spec's neighbourhood:\n%s", want, spec)
		}
	}
}

func TestMapShow_orphansAndStats(t *testing.T) {
	root := fixtureProject(t)

	orph := runCmd(t, newMapCmd(), "show", "--orphans", "--root", root)
	if !strings.Contains(orph, "(1):") || !strings.Contains(orph, "[guide] guides/LONELY.md") {
		t.Errorf("the lonely guide is the only orphan:\n%s", orph)
	}
	if strings.Contains(orph, "guides/CODE.md") {
		t.Errorf("a guide that governs something is not an orphan:\n%s", orph)
	}

	stats := runCmd(t, newMapCmd(), "show", "--stats", "--root", root)
	for _, want := range []string{"map: 5 nodes, 3 edges", "guide     2", "code      1", "governs     1", "specifies   1", "tested-by   1"} {
		if !strings.Contains(stats, want) {
			t.Errorf("missing %q in the stats:\n%s", want, stats)
		}
	}
}

// Topological order: the ruler before what it governs, the spec before its code.
func TestMapShow_worklistIsTopological(t *testing.T) {
	root := fixtureProject(t)
	out := runCmd(t, newMapCmd(), "show", "--worklist", "--root", root)

	pos := func(id string) int { return strings.Index(out, "  "+id+"\n") }
	for _, id := range []string{"guides/CODE.md", "src/login.spec.md", "src/login.ts", "src/login.test.ts"} {
		if pos(id) < 0 {
			t.Fatalf("%s is missing from the worklist:\n%s", id, out)
		}
	}
	if pos("guides/CODE.md") > pos("src/login.ts") || pos("src/login.spec.md") > pos("src/login.ts") {
		t.Errorf("a parent came after its child:\n%s", out)
	}
	if !strings.Contains(out, "5 node(s), in topological order.") {
		t.Errorf("the count is wrong:\n%s", out)
	}
}

// With --pending the gates run, and only the nodes with a FAILING gate are listed — the
// judgment gate is only pending, and does not count.
func TestMapShow_worklistPendingListsOnlyFailures(t *testing.T) {
	root := fixtureProject(t)
	out := runCmd(t, newMapCmd(), "show", "--worklist", "--pending", "--root", root)

	if !strings.Contains(out, "  src/login.ts\n") {
		t.Errorf("the code node fails `always-red` and was not listed:\n%s", out)
	}
	for _, id := range []string{"src/login.spec.md", "guides/CODE.md"} {
		if strings.Contains(out, "  "+id+"\n") {
			t.Errorf("%s has no failing gate and was listed:\n%s", id, out)
		}
	}
	if !strings.Contains(out, "1 file(s) with pending items") {
		t.Errorf("the count is wrong:\n%s", out)
	}
}

func TestMapShow_errors(t *testing.T) {
	root := fixtureProject(t)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"show", "--root", root}, "provide <file>"},
		{[]string{"show", "src/ghost.ts", "--root", root}, `ghost.ts" is not in the map`},
		{[]string{"show", "--stats", "--root", t.TempDir()}, "run `anchors map build`"},
	} {
		_, err := runCmdErr(newMapCmd(), t, c.args...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: expected an error with %q, got %v", c.args, c.want, err)
		}
	}
	// --pending needs the config: without it the worklist refuses instead of guessing.
	cfgless := t.TempDir()
	if err := mapx.Save(&mapx.Graph{}, cfgless+"/anchors.graph.yaml"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmdErr(newMapCmd(), t, "show", "--worklist", "--pending", "--root", cfgless); err == nil ||
		!strings.Contains(err.Error(), "load config") {
		t.Errorf("--pending without anchors.yaml: expected a config error, got %v", err)
	}
}

// The summary shows the triad's types first and EVERY other type after, alphabetically —
// a fixed list used to hide `realizes`, `depends-on`, and the rest.
func TestPrintEdgeSummary_showsTypesOutsideTheKnownList(t *testing.T) {
	useEnglish(t)
	g := &mapx.Graph{Edges: []mapx.Edge{
		{Type: mapx.EdgeType("realizes")},
		{Type: mapx.EdgeType("depends-on")},
		{Type: mapx.EdgeType("depends-on")},
		{Type: mapx.EdgeSpecifies},
	}}
	out := capturaStdout(t, func() { printEdgeSummary(g) })
	iSpec, iDep, iReal := strings.Index(out, "specifies   1"), strings.Index(out, "depends-on  2"), strings.Index(out, "realizes    1")
	if iSpec < 0 || iDep < 0 || iReal < 0 {
		t.Fatalf("a type is missing from the summary:\n%s", out)
	}
	if !(iSpec < iDep && iDep < iReal) {
		t.Errorf("order must be the triad first, then alphabetical:\n%s", out)
	}
	if got := capturaStdout(t, func() { printEdgeSummary(&mapx.Graph{}) }); got != "" {
		t.Errorf("no edges, no summary: %q", got)
	}
}

// Rebuilding the map keeps the flow graph (which `map build` does not scan), and warns
// when a judgment stamp was lost.
func TestMapBuild_keepsTheFlowAndWarnsOfLostStamps(t *testing.T) {
	root := fixtureProject(t)
	g := loadMap(t, root)
	g.Flow = &mapx.FlowGraph{States: []mapx.FlowState{{Code: "WORKR-P01", Title: "pull", Flow: "flows/w.flow.md"}}}
	// A stamp on an edge that the rebuild will not produce: it is lost, and must be told.
	g.Edges = append(g.Edges, mapx.Edge{From: "guides/LONELY.md", To: "src/gone.ts", Type: mapx.EdgeGoverns,
		Julgamentos: []mapx.Judgment{{Gate: "review", Verdict: "ok"}}})
	if err := mapx.Save(g, root+"/anchors.graph.yaml"); err != nil {
		t.Fatal(err)
	}

	out := runCmd(t, newMapCmd(), "build", "--root", root)
	if !strings.Contains(out, "LOST judgment stamp") || !strings.Contains(out, "review") {
		t.Errorf("the lost stamp was not reported:\n%s", out)
	}
	after := loadMap(t, root)
	if after.Flow == nil || len(after.Flow.States) != 1 {
		t.Errorf("the rebuild erased the flow graph: %+v", after.Flow)
	}
}

func TestMapBuild_withoutConfigPointsAtInit(t *testing.T) {
	_, err := runCmdErr(newMapCmd(), t, "build", "--root", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "anchors init") {
		t.Errorf("expected an error pointing at `anchors init`, got %v", err)
	}
}

// @anchors
//   code: MPTSM
//   ref: MPCMM

package mapcmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// ── build ───────────────────────────────────────────────────────────────────────

// A judgment recorded before a rebuild is still there after it: the workflow runs
// `map build` before every `check`, and a rebuild that started from nothing would undo
// every judgment of the previous step.
func TestMapBuild_keepsTheJudgments(t *testing.T) {
	t.Run("MPCMM-B01: A rebuild keeps the judgment recorded on the guide edge", func(t *testing.T) {})
	t.Run("MPCMM-I01: A judgment survives a rebuild of the map", func(t *testing.T) {})
	root := fixtureProject(t)
	runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "too big")

	runCmd(t, newMapCmd(), "build", "--root", root)

	if got := judgmentsOn(t, root, "src/login.ts", "atomic"); len(got) != 1 || got[0] != "guides/CODE.md→issue" {
		t.Errorf("the rebuild lost the judgment on the guide edge: %v", got)
	}
	if v, ok := loadMap(t, root).JudgedBy("src/login.ts", "atomic"); !ok || v != "issue" {
		t.Errorf("the judgment is no longer answered after the rebuild: %q %v", v, ok)
	}
}

// Rebuilding the map keeps the flow graph (which `map build` does not scan), and warns
// when a judgment stamp was lost.
func TestMapBuild_keepsTheFlowAndWarnsOfLostStamps(t *testing.T) {
	t.Run("MPCMM-B02: A rebuild keeps the flow graph", func(t *testing.T) {})
	root := fixtureProject(t)
	g := loadMap(t, root)
	g.Flow = &mapx.FlowGraph{States: []mapx.FlowState{{Code: "WORKR-T01", Title: "pull", Flow: "flows/w.flow.md"}}}
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

// STAMP LOSS is silent, and the map stays VALID. Measured in the reference app
// (co2-lab/anchors#12): a conflict resolved with `checkout --theirs` + `map build` lost the
// `review` stamp of a spec — found by chance, counting 18 where 19 were expected. The cost is
// the REPORT, which lives in the command's `--reason`, not in the file.
func TestStampLossWarning(t *testing.T) {
	t.Run("MPCMM-B03: A lost judgment stamp is reported per gate", func(t *testing.T) {})
	withJudgments := func(gates ...string) *mapx.Graph {
		g := &mapx.Graph{}
		for _, gate := range gates {
			g.Edges = append(g.Edges, mapx.Edge{
				Julgamentos: []mapx.Judgment{{Gate: gate}},
			})
		}
		return g
	}

	t.Run("a loss WARNS, counted per gate", func(t *testing.T) {
		old := withJudgments("review", "review", "rule-fulfilled")
		updated := withJudgments("review", "rule-fulfilled")

		warning := stampLossWarning(old, updated)
		if warning == "" {
			t.Fatal("a review stamp was lost and nothing was said")
		}
		if !strings.Contains(warning, "review") || !strings.Contains(warning, "2 → 1") {
			t.Errorf("the warning does not say WHICH gate and how much:\n%s", warning)
		}
		// per gate, not the total: "lost 1 stamp" does not say what to redo
		if strings.Contains(warning, "rule-fulfilled") {
			t.Errorf("a gate that lost nothing was reported:\n%s", warning)
		}
		// the warning has to say what is really lost
		if !strings.Contains(warning, "report") {
			t.Errorf("the warning does not say the REPORT is lost:\n%s", warning)
		}
	})

	// A warning that always shows stops being read.
	t.Run("no loss, silence", func(t *testing.T) {
		same := withJudgments("review", "rule-fulfilled")
		if got := stampLossWarning(same, withJudgments("review", "rule-fulfilled")); got != "" {
			t.Errorf("warned with no loss: %q", got)
		}
		// GAINING a stamp is silent too: it is the normal case of judging something new
		more := withJudgments("review", "review", "rule-fulfilled")
		if got := stampLossWarning(same, more); got != "" {
			t.Errorf("warned on a GAIN: %q", got)
		}
	})

	t.Run("an empty previous map does not warn", func(t *testing.T) {
		if got := stampLossWarning(&mapx.Graph{}, withJudgments("review")); got != "" {
			t.Errorf("warned with an empty previous map: %q", got)
		}
	})
}

// The summary shows the unit's types first and EVERY other type after, alphabetically —
// a fixed list used to hide `realizes`, `depends-on`, and the rest.
func TestPrintEdgeSummary_showsTypesOutsideTheKnownList(t *testing.T) {
	t.Run("MPCMM-B04: The edge summary shows every type, the unit's first", func(t *testing.T) {})
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
		t.Errorf("order must be the unit first, then alphabetical:\n%s", out)
	}
	if got := capturaStdout(t, func() { printEdgeSummary(&mapx.Graph{}) }); got != "" {
		t.Errorf("no edges, no summary: %q", got)
	}
}

// The warning groups by PAIR of layers, not by file: in a monorepo the same pair produces
// hundreds of identical lines, and dumping them hides the pattern that has to be seen.
func TestLayerAmbiguity_warningGroupsByPairOfLayers(t *testing.T) {
	t.Run("MPCMM-B05: The layer ambiguity warning is grouped by pair of layers", func(t *testing.T) {})
	amb := []scan.LayerAmbiguity{
		{Arquivo: "a.test.ts", Vencedora: "shared", Perdedoras: []string{"test"}},
		{Arquivo: "b.test.ts", Vencedora: "shared", Perdedoras: []string{"test"}},
		{Arquivo: "c.test.ts", Vencedora: "lambdas", Perdedoras: []string{"test"}},
	}
	out := capturaStdout(t, func() { printLayerAmbiguities(amb) })

	if !strings.Contains(out, "3 file(s)") {
		t.Errorf("the total is missing:\n%s", out)
	}
	if !strings.Contains(out, "shared beat test") || !strings.Contains(out, "(2 file(s), e.g.: a.test.ts)") {
		t.Errorf("the pair `shared beat test` with count 2 and an example is missing:\n%s", out)
	}
	if !strings.Contains(out, "lambdas beat test") {
		t.Errorf("the second pair is missing:\n%s", out)
	}
	// the consequence has to be said: without it the warning looks cosmetic
	if !strings.Contains(out, "EVERY gate") {
		t.Errorf("the warning does not say what is lost:\n%s", out)
	}
	// silence when there is no ambiguity: a warning that always shows stops being read.
	if s := capturaStdout(t, func() { printLayerAmbiguities(nil) }); s != "" {
		t.Errorf("printed with no ambiguity: %q", s)
	}
}

// THE WIRING, not only the function. `scan.Ambiguities` was complete and tested, and
// `map build` did not call it (reference app): a test of the function alone passes with the
// mechanism switched off. Running the real command over an ambiguous project confronts the
// wiring — and the warning must not fail the build.
func TestMapBuild_warnsOfTheLayerAmbiguity(t *testing.T) {
	t.Run("MPCMM-X01: The warnings do not fail the build", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "packages", "shared"), 0o755); err != nil {
		t.Fatal(err)
	}
	// both layers match the same file, and neither declares `priority`
	yaml := `version: 1
layers:
  test:
    pattern: "**/*.test.ts"
    kind: test
  shared:
    pattern: "packages/shared/**/*.ts"
    kind: code
`
	if err := os.WriteFile(filepath.Join(root, config.DefaultFile), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "packages", "shared", "AreaStatus.test.ts")
	if err := os.WriteFile(target, []byte("export const x = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := capturaStdout(t, func() {
		cmd := newMapCmd()
		cmd.SetArgs([]string{"build", "--root", root})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("map build: %v", err)
		}
	})

	if !strings.Contains(out, "HEURISTIC") {
		t.Errorf("`map build` did NOT warn of the ambiguity — `scan.Ambiguities` is not wired:\n%s", out)
	}
	if !strings.Contains(out, "shared beat test") {
		t.Errorf("the pair of layers is missing from the output:\n%s", out)
	}
}

func TestMapBuild_withoutConfigPointsAtInit(t *testing.T) {
	t.Run("MPCMM-E01: Building with no configuration points at init", func(t *testing.T) {})
	_, err := runCmdErr(newMapCmd(), t, "build", "--root", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "anchors init") {
		t.Errorf("expected an error pointing at `anchors init`, got %v", err)
	}
}

// ── show ────────────────────────────────────────────────────────────────────────

// The neighbourhood of the code node: governed by the guide and specified by the spec
// (going up), and nothing below it — a leaf.
func TestMapShow_neighbourhoodOfANode(t *testing.T) {
	t.Run("MPCMM-B06: Showing the login code lists what governs it and marks it a leaf", func(t *testing.T) {})
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
	t.Run("MPCMM-B07: The orphans are the nodes with no edge", func(t *testing.T) {})
	t.Run("MPCMM-B08: The statistics count nodes by kind and edges by type", func(t *testing.T) {})
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
	t.Run("MPCMM-B09: The worklist puts rulers and specs before the code they govern", func(t *testing.T) {})
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
	t.Run("MPCMM-B10: The pending worklist lists only nodes with a failing gate", func(t *testing.T) {})
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
	t.Run("MPCMM-E02: Showing with no map points at the build", func(t *testing.T) {})
	t.Run("MPCMM-E03: Showing a file that is not in the map is refused", func(t *testing.T) {})
	t.Run("MPCMM-E04: Showing with no selector is refused", func(t *testing.T) {})
	t.Run("MPCMM-E05: The pending worklist with no configuration is refused", func(t *testing.T) {})
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

// THE AMBIGUITY IS DETECTED, and a declared priority silences it. `scan.Ambiguities` existed
// — complete and tested — and nobody called it: the comment of the `priority` field promised
// "declare it when the heuristic is wrong; `check` warns where it decided alone", and it did
// not warn.
//
// Cost measured in the reference app: `**/*.test.*` and `packages/shared/**/*.ts` matched the
// same `AreaStatus.test.ts`, the length tie-break chose `shared`, and the project had ZERO
// `kind: test` nodes with 70 green tests. In cascade FOUR gates went blind, including the
// blocking `test-traceable`. It was only found by counting the map's nodes by hand.
func TestLayerAmbiguity_thePairOfLayersIsReportedUntilAPriorityIsDeclared(t *testing.T) {
	t.Run("MPCMM-B05: The layer ambiguity warning is grouped by pair of layers", func(t *testing.T) {})
	t.Run("MPCMM-B11: A declared priority silences the layer ambiguity warning", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"test":   {Pattern: "**/*.test.ts", Kind: "test"},
		"shared": {Pattern: "packages/shared/**/*.ts", Kind: "code"},
	}}
	files := []scan.File{
		{Path: "packages/shared/AreaStatus.test.ts"},
		{Path: "packages/shared/QueryScope.test.ts"},
	}

	amb := scan.Ambiguities(files, cfg)
	if len(amb) != 2 {
		t.Fatalf("Ambiguities returned %d, want 2 — both files match both layers", len(amb))
	}
	// `packages/shared/**/*.ts` (23) is longer than `**/*.test.ts` (12), so the length
	// tie-break hands the TEST to the code layer — the defect of #100.
	if amb[0].Vencedora != "shared" {
		t.Errorf("winner = %q; the pattern length favours `shared`", amb[0].Vencedora)
	}

	// declaring `priority` on the right layer SILENCES the warning: the project has already
	// decided, and there is no guess to report.
	cfg.Layers["test"] = config.Layer{Pattern: "**/*.test.ts", Kind: "test", Priority: 100}
	if amb := scan.Ambiguities(files, cfg); len(amb) != 0 {
		t.Errorf("with `priority` declared the warning stayed: %+v", amb)
	}
}

// capturaStdout collects what the function writes to `os.Stdout`.
//
// The commands print with `fmt.Println` directly (and not through `cmd.OutOrStdout()`), so
// testing the OUTPUT requires intercepting the descriptor. It is what lets a test assert on
// the text the user reads, instead of only on the return value.
func capturaStdout(t *testing.T, f func()) string {
	t.Helper()
	return testkit.CaptureStdout(t, f)
}

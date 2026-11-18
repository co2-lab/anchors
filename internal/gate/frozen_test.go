// @anchors
//   code: FGTFR
//   ref: FRGTF

package gate

import (
	"github.com/co2-lab/anchors/internal/testkit"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// frozenRun writes the files under a temporary root and resets the run's freezes after the
// test, so no other test sees them.
func frozenRun(t *testing.T, files map[string]string) string {
	t.Helper()
	t.Cleanup(func() { setFreezes(nil) })
	root := t.TempDir()
	for rel, body := range files {
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

func frozenGraph() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.ts", Kind: mapx.KindCode, Frozen: "FRZQA-F01"},
		{ID: "b.ts", Kind: mapx.KindCode},
		{ID: "product/q.doctrine.md", Kind: mapx.KindProduct, Freezes: map[string]string{"FRZQA-F01": "the plan mode ships next quarter"}},
	}}
}

func TestFrozen_gateAnswersFrozen(t *testing.T) {
	t.Run("FRGTF-B01: A gate on a frozen part answers frozen, with the reason", func(t *testing.T) {})
	root := frozenRun(t, map[string]string{"a.ts": "", "b.ts": "", "product/q.doctrine.md": "x"})
	g := frozenGraph()
	nonEmpty := config.Gate{Name: "non-empty", Check: "non-empty", On: []string{"code"}}
	res := RunWithWaiver([]config.Gate{nonEmpty}, g.Nodes, root, g, testkit.WithComments(&config.Config{}), false, Waiver{})
	got := map[string]Result{}
	for _, r := range res {
		got[r.Target] = r
	}
	if a := got["a.ts"]; a.Verdict != Frozen || !strings.Contains(a.Detail, "FRZQA-F01 — the plan mode ships next quarter") {
		t.Errorf("the frozen file: %+v", a)
	}
	if b := got["b.ts"]; b.Verdict != Fail {
		t.Errorf("the other file fails as before: %+v", b)
	}
	batch := config.Gate{Name: "lint", On: []string{"code"}, Scope: config.ScopeBatch, Run: "echo {{files}}; exit 1"}
	res = RunWithWaiver([]config.Gate{batch}, g.Nodes, root, g, testkit.WithComments(&config.Config{}), false, Waiver{})
	if len(res) != 1 || strings.Contains(res[0].Detail, "a.ts") || !strings.Contains(res[0].Detail, "b.ts") {
		t.Errorf("the batch gate is given the unfrozen file alone: %+v", res)
	}
	g.Nodes[1].Frozen = "a plain cause"
	if res = RunWithWaiver([]config.Gate{batch}, g.Nodes, root, g, testkit.WithComments(&config.Config{}), false, Waiver{}); len(res) != 0 {
		t.Errorf("every target frozen, the batch gate does not run: %+v", res)
	}
}

func TestFrozen_ruleLeavesTheQuestions(t *testing.T) {
	t.Run("FRGTF-B02: A frozen rule leaves the questions asked of a rule", func(t *testing.T) {})
	frozenRun(t, nil)
	setFreezes(&mapx.Graph{Nodes: []mapx.Node{
		{ID: "s.spec.md", Kind: mapx.KindSpec, FrozenRules: map[string]string{"FRZQB-B02": "FRZQA-F01"}},
		{ID: "product/q.doctrine.md", Kind: mapx.KindProduct, FrozenRules: map[string]string{"FRZQD-R02": "FRZQA-F01"}},
	}})
	spec := "# S\n\n## Effects\n\n| `FRZQB-B01` | live |\n| `FRZQB-B02` | frozen |\n"
	if got := definedRequirements(spec); len(got) != 1 || got[0] != "FRZQB-B01" {
		t.Errorf("only the live rule is required: %v", got)
	}
	doctrine := "### FRZQD-R01 — one\n### FRZQD-R02 — two\n### FRZQD-F01 — the plan mode ships next quarter\n"
	n := mapx.Node{ID: "product/q.doctrine.md", Kind: mapx.KindProduct}
	g := &mapx.Graph{Nodes: []mapx.Node{n}, Edges: []mapx.Edge{{From: "s.spec.md", To: n.ID, Type: mapx.EdgeRealizes, Method: "FRZQD-R01"}}}
	if v, msg := checkDoctrineRealized(doctrine, n, "", g, testkit.WithComments(nil)); v != Pass {
		t.Errorf("no realizer is asked of a freeze entry or a frozen rule: %v %s", v, msg)
	}
	if texts := ruleTexts(doctrine); texts["FRZQD-F01"] != "" {
		t.Errorf("a freeze entry is no rule to repeat: %v", texts)
	}
}

func TestFrozen_consistent(t *testing.T) {
	t.Run("FRGTF-B03: The freeze marks and records are confronted", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.ts", Kind: mapx.KindCode, Frozen: "FRZQA-F02"},
		{ID: "b.ts", Kind: mapx.KindCode, Frozen: "a plain cause"},
		{ID: "product/q.doctrine.md", Kind: mapx.KindProduct, Freezes: map[string]string{"FRZQA-F01": "the plan mode"}},
	}}
	v, msg := checkFreezeConsistent("", mapx.Node{}, "", g, testkit.WithComments(nil))
	if v != Fail || !strings.Contains(msg, "a.ts → FRZQA-F02") || !strings.Contains(msg, "FRZQA-F01 (product/q.doctrine.md)") || strings.Contains(msg, "plain") {
		t.Errorf("unknown mark and orphan record: %v %s", v, msg)
	}
	g.Nodes[0].Frozen = "FRZQA-F01"
	if v, msg := checkFreezeConsistent("", mapx.Node{}, "", g, testkit.WithComments(nil)); v != Pass {
		t.Errorf("consistent: %v %s", v, msg)
	}
	if v, _ := checkFreezeConsistent("", mapx.Node{}, "", &mapx.Graph{}, testkit.WithComments(nil)); v != Skip {
		t.Errorf("no freeze: %v", v)
	}
}

func TestFrozen_touched(t *testing.T) {
	t.Run("FRGTF-B04: A change to a frozen part is told and goes on", func(t *testing.T) {})
	frozenRun(t, nil)
	restore := SetChangedSource(func(_, rel string) (bool, bool) { return rel != "c.ts", true })
	defer restore()
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.ts", Kind: mapx.KindCode, Frozen: "FRZQA-F01"},
		{ID: "s.spec.md", Kind: mapx.KindSpec, FrozenRules: map[string]string{"FRZQB-B02": "FRZQA-F01"}},
		{ID: "c.ts", Kind: mapx.KindCode, Frozen: "FRZQA-F01"},
		{ID: "d.ts", Kind: mapx.KindCode},
		{ID: "product/q.doctrine.md", Kind: mapx.KindProduct, Freezes: map[string]string{"FRZQA-F01": "the plan mode"}},
	}}
	setFreezes(g)
	for i, want := range []Verdict{Diverge, Diverge, Pass, Skip} {
		v, msg := checkFrozenTouched("", g.Nodes[i], "", g, testkit.WithComments(nil))
		if v != want {
			t.Errorf("%s: %v (%s), want %v", g.Nodes[i].ID, v, msg, want)
		}
		if v == Diverge && (!strings.Contains(msg, "FRZQA-F01 — the plan mode") || !strings.Contains(msg, DeclaredMarker)) {
			t.Errorf("%s names the freeze and informs: %s", g.Nodes[i].ID, msg)
		}
	}
	tg := config.Gate{Name: "frozen-touched", Check: "frozen-touched", On: []string{"code"}, Blocking: config.Bool(true)}
	r := Result{Gate: tg.Name, Verdict: Diverge, Blocking: true, Detail: "x " + DeclaredMarker}
	markSeverity(tg, &r)
	if r.Blocks() {
		t.Error("a declared divergence never bars")
	}
}

func TestFrozen_countedApart(t *testing.T) {
	t.Run("FRGTF-B05: Frozen verdicts are counted apart", func(t *testing.T) {})
	p := Aggregate([]Result{
		{Gate: "non-empty", Target: "a.ts", Verdict: Frozen, Blocking: true},
		{Gate: "non-empty", Target: "b.ts", Verdict: Pass, Blocking: true},
	})
	if s := p.ByGate["non-empty"]; s.Frozen != 1 || s.Pass != 1 || !p.Passed || len(p.Failures) != 0 {
		t.Errorf("counted apart: %+v passed=%v", s, p.Passed)
	}
}

func TestFrozen_thawGivesBack(t *testing.T) {
	t.Run("FRGTF-I01: Thawing gives back the verdicts", func(t *testing.T) {})
	root := frozenRun(t, map[string]string{"a.ts": ""})
	g := frozenGraph()
	nonEmpty := config.Gate{Name: "non-empty", Check: "non-empty", On: []string{"code"}}
	nodes := g.Nodes[:1]
	if res := RunWithWaiver([]config.Gate{nonEmpty}, nodes, root, g, testkit.WithComments(&config.Config{}), false, Waiver{}); res[0].Verdict != Frozen {
		t.Fatalf("frozen: %+v", res)
	}
	g.Nodes[0].Frozen = ""
	if res := RunWithWaiver([]config.Gate{nonEmpty}, g.Nodes[:1], root, g, testkit.WithComments(&config.Config{}), false, Waiver{}); res[0].Verdict != Fail {
		t.Errorf("thawed, it fails as before: %+v", res)
	}
}

func TestFrozen_noMap(t *testing.T) {
	t.Run("FRGTF-E01: With no map the freezes are not confronted", func(t *testing.T) {})
	if v, msg := checkFreezeConsistent("", mapx.Node{}, "", nil, testkit.WithComments(nil)); v != Pending || msg == "" {
		t.Errorf("no map: %v %q", v, msg)
	}
}

func TestFrozen_emptyFolder(t *testing.T) {
	t.Run("FRGTF-B06: A frozen folder that holds no file is named", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "ir/tax.ts", Kind: mapx.KindCode, Frozen: "out of the plan"}}}
	cfg := testkit.WithComments(&config.Config{FrozenPaths: []config.FrozenPath{
		{Freeze: "out of the plan", Paths: []string{"ir/**", "irr/**"}},
	}})
	v, msg := checkFreezeConsistent("", mapx.Node{}, "", g, cfg)
	if v != Fail || !strings.Contains(msg, "irr/** (out of the plan)") || strings.Contains(msg, "ir/** ") {
		t.Errorf("the empty glob is named: %v %s", v, msg)
	}
	cfg.FrozenPaths[0].Paths = []string{"ir/**"}
	if v, msg := checkFreezeConsistent("", mapx.Node{}, "", g, cfg); v != Pass {
		t.Errorf("a folder with files is consistent: %v %s", v, msg)
	}
}

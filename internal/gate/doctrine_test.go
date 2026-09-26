package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// projectWithDoctrine builds a project where one spec realizes one product rule.
func projectWithDoctrine(t *testing.T, doctrine string) (string, *mapx.Graph, mapx.Node) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "product"), 0o755)
	if err := os.WriteFile(filepath.Join(root, "product/l.doctrine.md"), []byte(doctrine), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "product/l.doctrine.md", Kind: mapx.KindProduct, Code: "LIMIT"}},
		Edges: []mapx.Edge{{
			From: "s.spec.md", To: "product/l.doctrine.md",
			Type: mapx.EdgeRealizes, Method: "LIMIT-R03", Dep: "CRED-V01",
		}},
	}
	return root, g, mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}
}

// The doctrine carries SEVERAL rules on purpose: the similarity ruler is IDF, and with a
// tiny corpus every shared word weighs zero — the gate would measure nothing and say so
// (Pending). Four rules is the measured floor where the ruler starts to discriminate.
const doctrineRule = "### LIMIT-R03 — o limite de credito nunca e excedido em nenhuma operacao\n" +
	"### LIMIT-R04 — a sessao expira apos trinta minutos sem interacao\n" +
	"### LIMIT-R05 — o relatorio mensal inclui apenas transacoes aprovadas\n" +
	"### LIMIT-R06 — o usuario anonimo e redirecionado para a tela inicial\n"

// The defect this whole axis exists to eliminate: the spec COPIES the doctrine's text
// instead of referencing it. Two texts saying the same thing diverge at the first change,
// and no other gate sees it — both are well-formed, both catalogue their rules, both have
// a complete triad.
func TestDoctrineNotDuplicated_copyFails(t *testing.T) {
	root, g, n := projectWithDoctrine(t, doctrineRule)
	spec := "### CRED-V01 — o limite de credito nunca e excedido em nenhuma operacao    @realizes LIMIT-R03\n"
	v, msg := checkDoctrineNotDuplicated(spec, n, root, g, nil)
	if v != Fail {
		t.Fatalf("expected Fail for a copy, got %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "CRED-V01") || !strings.Contains(msg, "LIMIT-R03") {
		t.Errorf("the verdict must name both ends: %s", msg)
	}
}

// The legitimate case, and the one that must NOT be accused: the spec says what is
// specific to its unit and lets `@realizes` carry the shared rule.
func TestDoctrineNotDuplicated_specificTextPasses(t *testing.T) {
	root, g, n := projectWithDoctrine(t, doctrineRule)
	spec := "### CRED-V01 — desabilita o botao enviar quando o campo valor esta vazio    @realizes LIMIT-R03\n"
	if v, msg := checkDoctrineNotDuplicated(spec, n, root, g, nil); v != Pass {
		t.Errorf("unit-specific text must not be accused: %v (%s)", v, msg)
	}
}

// `@TBD` on the line is DEBT, not a waiver: the wording is still being worked out.
// Charging it now would push whoever is writing to paraphrase for the gate instead of for
// the reader — which is the opposite of what the axis wants.
func TestDoctrineNotDuplicated_tbdDefersTheCharge(t *testing.T) {
	root, g, n := projectWithDoctrine(t, doctrineRule)
	spec := "### CRED-V01 — o limite de credito nunca e excedido em nenhuma operacao    @realizes LIMIT-R03 @TBD: redacao em revisao\n"
	if v, msg := checkDoctrineNotDuplicated(spec, n, root, g, nil); v == Fail {
		t.Errorf("@TBD on the line must defer the charge: %s", msg)
	}
}

func TestDoctrineNotDuplicated_onlySpec(t *testing.T) {
	root, g, _ := projectWithDoctrine(t, doctrineRule)
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindFeature, mapx.KindTest, mapx.KindProduct} {
		n := mapx.Node{ID: "s.spec.md", Kind: k}
		if v, _ := checkDoctrineNotDuplicated("x", n, root, g, nil); v != Skip {
			t.Errorf("kind %v: expected Skip, got %v", k, v)
		}
	}
}

// The ruler is SIMILARITY, not equality, and this is the case that proves why: whoever
// copies almost always adjusts a word. Exact comparison would catch only the laziest copy
// and report green on every other — which is the worst failure for a measuring
// instrument, because the gate looks like it is working.
func TestDoctrineNotDuplicated_nearCopyWithAWordChangedFails(t *testing.T) {
	root, g, n := projectWithDoctrine(t, doctrineRule)
	// Same sentence as the doctrine, with "credito" -> "emprestimo" and one word dropped.
	spec := "### CRED-V01 — o limite de emprestimo nunca e excedido em operacao    @realizes LIMIT-R03\n"
	v, msg := checkDoctrineNotDuplicated(spec, n, root, g, nil)
	if v != Fail {
		t.Fatalf("a near copy must be accused — exact comparison would let it through: %v (%s)", v, msg)
	}
}

// The demand is DECLARED BY LAYER, never universal.
//
// Measured in this repository: of 841 catalogued rules the overwhelming majority is local
// to its unit. Demanding `@realizes` everywhere would force inventing umbrella doctrine
// just to silence the gate — the vice `placeholder-filled` exists to catch. In a product
// application the proportion inverts, and only the Structure knows which case it is.
func TestSpecRealizesDoctrine_onlyWhereTheLayerDemands(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"screen": {RequiresDoctrine: true},
		"gate":   {},
	}}
	spec := "### CRED-V01 — a rule with no declaration\n"

	demanding := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Tags: []string{"screen"}}
	v, msg := checkSpecRealizesDoctrine(spec, demanding, "", nil, cfg)
	if v != Fail {
		t.Errorf("a demanding layer must accuse a naked rule: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "CRED-V01") {
		t.Errorf("the verdict must name the rule: %s", msg)
	}

	quiet := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Tags: []string{"gate"}}
	if v, msg := checkSpecRealizesDoctrine(spec, quiet, "", nil, cfg); v != Skip {
		t.Errorf("a layer that does not demand must stay quiet: %v (%s)", v, msg)
	}
}

// A rule that DECLARES is not accused, and one deferred with `@TBD` becomes debt rather
// than a failure — the doctrine it will realize is still being written.
func TestSpecRealizesDoctrine_declaredPassesAndTbdDefers(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{"screen": {RequiresDoctrine: true}}}
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Tags: []string{"screen"}}

	if v, msg := checkSpecRealizesDoctrine("### CRED-V01 — a rule    @realizes LIMIT-R03\n", n, "", nil, cfg); v != Pass {
		t.Errorf("a declared rule must pass: %v (%s)", v, msg)
	}
	if v, msg := checkSpecRealizesDoctrine("### CRED-V01 — a rule    @TBD: doctrine being written\n", n, "", nil, cfg); v != Pending {
		t.Errorf("`@TBD` must be debt, not failure: %v (%s)", v, msg)
	}
}

// An OPEN QUESTION (`-Q`) is not a realizable rule, and demanding a realizer for one
// would ask a spec to concretise the very thing still undecided.
//
// Found against the real doctrine: `VGRUP-Q01` records a decision about how the grouping
// filter should work, and the gate accused it of having no realizer. The only way to
// comply would be to answer the question inside a spec — which is where it must NOT be
// answered. `open-questions-resolved` is the gate that charges a `-Q`, and it charges the
// right thing: that someone DECIDES.
func TestDoctrineRealized_openQuestionIsNotARule(t *testing.T) {
	doctrine := "### LIMIT-R03 — a rule    @TBD: not realized yet\n\n" +
		"| `LIMIT-Q01` | what should happen when the limit changes mid-flow? | product | a rule |\n"
	root := t.TempDir()
	n := mapx.Node{ID: "product/l.doctrine.md", Kind: mapx.KindProduct, Code: "LIMIT"}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}

	v, msg := checkDoctrineRealized(doctrine, n, root, g, nil)
	if strings.Contains(msg, "LIMIT-Q01") {
		t.Errorf("an open question must not be charged for a realizer: %s", msg)
	}
	if v == Fail {
		t.Errorf("only the deferred rule remains, so this is debt and not failure: %v (%s)", v, msg)
	}
}

// planWithDoctrines builds a root holding the given doctrine files.
func planWithDoctrines(t *testing.T, existing ...string) string {
	t.Helper()
	root := t.TempDir()
	for _, d := range existing {
		p := filepath.Join(root, d)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("### LIMIT-R01 — x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestPlanDoctrineExists(t *testing.T) {
	plan := mapx.Node{ID: "plans/p.plan.md", Kind: mapx.KindPlan}

	t.Run("only a plan is confronted", func(t *testing.T) {
		root := planWithDoctrines(t)
		if v, _ := checkPlanDoctrineExists("seeds `product/a.doctrine.md`\n", mapx.Node{ID: "x", Kind: mapx.KindSpec}, root, nil, nil); v != Skip {
			t.Errorf("a spec is not a plan: %v", v)
		}
	})

	t.Run("a cited doctrine that exists passes", func(t *testing.T) {
		root := planWithDoctrines(t, "product/a.doctrine.md")
		if v, msg := checkPlanDoctrineExists("seeds `product/a.doctrine.md`\n", plan, root, nil, nil); v != Pass {
			t.Errorf("got %v (%s)", v, msg)
		}
	})

	t.Run("every missing doctrine on a line is named, once", func(t *testing.T) {
		root := planWithDoctrines(t, "product/a.doctrine.md")
		content := "seeds `product/z.doctrine.md` and `product/b.doctrine.md`\n" +
			"again `product/z.doctrine.md`, and `product/a.doctrine.md`\n"
		v, msg := checkPlanDoctrineExists(content, plan, root, nil, nil)
		if v != Fail {
			t.Fatalf("got %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "product/b.doctrine.md, product/z.doctrine.md") || strings.Contains(msg, "a.doctrine") {
			t.Errorf("both missing doctrines, sorted, and only those: %s", msg)
		}
		if !strings.Contains(msg, "2") || strings.Contains(msg, "3") {
			t.Errorf("the count is of distinct missing doctrines: %s", msg)
		}
	})

	t.Run("a missing doctrine on a @TBD line is debt", func(t *testing.T) {
		root := planWithDoctrines(t)
		v, msg := checkPlanDoctrineExists("seeds `product/a.doctrine.md` @TBD: written next cycle\n", plan, root, nil, nil)
		if v != Pending || !strings.Contains(msg, "product/a.doctrine.md") {
			t.Errorf("got %v (%s)", v, msg)
		}
	})

	t.Run("prose citations and templates seed nothing", func(t *testing.T) {
		root := planWithDoctrines(t)
		content := "the `x.doctrine.md` convention, see `product/_TEMPLATE.doctrine.md`\n"
		if v, msg := checkPlanDoctrineExists(content, plan, root, nil, nil); v != Skip {
			t.Errorf("got %v (%s)", v, msg)
		}
	})
}

func TestDoctrineRealized(t *testing.T) {
	n := mapx.Node{ID: "product/l.doctrine.md", Kind: mapx.KindProduct}
	doctrine := "### LIMIT-R01 — one\n### LIMIT-R02 — two\n"
	graph := func(methods ...string) *mapx.Graph {
		g := &mapx.Graph{Nodes: []mapx.Node{n}}
		for _, m := range methods {
			g.Edges = append(g.Edges, mapx.Edge{From: "s.spec.md", To: n.ID, Type: mapx.EdgeRealizes, Method: m})
		}
		// an edge of another type carrying a rule code must not count as a realizer
		g.Edges = append(g.Edges, mapx.Edge{From: "t.spec.md", To: n.ID, Type: mapx.EdgeSpecifies, Method: "LIMIT-R02"})
		return g
	}

	if v, _ := checkDoctrineRealized(doctrine, mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, "", graph(), nil); v != Skip {
		t.Errorf("only a doctrine is confronted: %v", v)
	}
	if v, msg := checkDoctrineRealized(doctrine, n, "", nil, nil); v != Pending || msg != i18n.T("gate.no_map_loaded") {
		t.Errorf("without a map there is nothing to confront: %v (%s)", v, msg)
	}
	if v, _ := checkDoctrineRealized("prose, no rules\n", n, "", graph(), nil); v != Skip {
		t.Errorf("a doctrine with no rules: %v", v)
	}
	if v, msg := checkDoctrineRealized(doctrine, n, "", graph("LIMIT-R01", "LIMIT-R02"), nil); v != Pass {
		t.Errorf("every rule realized: %v (%s)", v, msg)
	}
	v, msg := checkDoctrineRealized(doctrine, n, "", graph("LIMIT-R01"), nil)
	if v != Fail || !strings.Contains(msg, "LIMIT-R02") || strings.Contains(msg, "LIMIT-R01") {
		t.Errorf("the unrealized rule, and only it, fails: %v (%s)", v, msg)
	}
}

func TestSpecDoctrineExists(t *testing.T) {
	root, g, n := projectWithDoctrine(t, doctrineRule) // one realizes edge, Method LIMIT-R03
	// an edge of another type is never read as a doctrine
	g.Edges = append(g.Edges, mapx.Edge{From: n.ID, To: "s.go", Type: mapx.EdgeSpecifies, Method: "LIMIT-R09"})

	if v, _ := checkSpecDoctrineExists("@realizes LIMIT-R03\n", mapx.Node{ID: "x.go", Kind: mapx.KindCode}, root, g, nil); v != Skip {
		t.Errorf("only a spec is confronted: %v", v)
	}
	if v, _ := checkSpecDoctrineExists("@realizes LIMIT-R03\n", n, root, nil, nil); v != Pending {
		t.Errorf("without a map there is nothing to confront: %v", v)
	}
	if v, _ := checkSpecDoctrineExists("### CRED-V01 — no tag\n", n, root, g, nil); v != Skip {
		t.Errorf("a spec that declares nothing: %v", v)
	}
	if v, msg := checkSpecDoctrineExists("### CRED-V01 — x @realizes LIMIT-R03\n", n, root, g, nil); v != Pass {
		t.Errorf("the rule the edge resolved: %v (%s)", v, msg)
	}
	// Two rules of the same doctrine on ONE line: the map holds one edge per FILE, and the
	// rule the edge does not name is confirmed by reading the doctrine (LIMIT-R06 is its
	// last rule, so the whole file must be searched).
	if v, msg := checkSpecDoctrineExists("x @realizes LIMIT-R03 @realizes LIMIT-R06\n", n, root, g, nil); v != Pass {
		t.Errorf("a rule present in the resolved doctrine: %v (%s)", v, msg)
	}
	v, msg := checkSpecDoctrineExists("x @realizes LIMIT-R03 @realizes LIMIT-R09\n", n, root, g, nil)
	if v != Fail || !strings.Contains(msg, "LIMIT-R09") || strings.Contains(msg, "LIMIT-R03") {
		t.Errorf("a renamed rule fails, named: %v (%s)", v, msg)
	}
	if v, msg := checkSpecDoctrineExists("x @realizes LIMIT-R09 @TBD: doctrine being written\n", n, root, g, nil); v != Skip {
		t.Errorf("a deferred declaration is not charged: %v (%s)", v, msg)
	}
}

// With exactly the floor of four rules the ruler measures: three doctrine rules and the
// copy make a corpus of four.
func TestDoctrineNotDuplicated_theFloorItselfMeasures(t *testing.T) {
	three := "### LIMIT-R03 — o limite de credito nunca e excedido em nenhuma operacao\n" +
		"### LIMIT-R04 — a sessao expira apos trinta minutos sem interacao\n" +
		"### LIMIT-R05 — o relatorio mensal inclui apenas transacoes aprovadas\n"
	root, g, n := projectWithDoctrine(t, three)
	spec := "### CRED-V01 — o limite de credito nunca e excedido em nenhuma operacao    @realizes LIMIT-R03\n"
	if v, msg := checkDoctrineNotDuplicated(spec, n, root, g, nil); v != Fail {
		t.Errorf("a corpus of %d rules is measured: %v (%s)", minCorpusForIDF, v, msg)
	}
}

func TestDoctrineNotDuplicated_reportsTheScoreAsAPercentage(t *testing.T) {
	root, g, n := projectWithDoctrine(t, doctrineRule)
	spec := "### CRED-V01 — o limite de credito nunca e excedido em nenhuma operacao    @realizes LIMIT-R03\n"
	_, msg := checkDoctrineNotDuplicated(spec, n, root, g, nil)
	if !strings.Contains(msg, "(100%)") {
		t.Errorf("an identical copy scores 100%%: %s", msg)
	}
}

func TestRuleTexts_aRuleWithOnlyATagHasNoText(t *testing.T) {
	got := ruleTexts("### CRED-V01 @realizes LIMIT-R03\n### CRED-V02 — real text @realizes LIMIT-R04\n")
	if len(got) != 1 || got["CRED-V02"] != "real text" {
		t.Errorf("got %q", got)
	}
}

func TestParseRealizesWithLines(t *testing.T) {
	content := "### CRED-V01 — text\n" +
		"continued, @realizes LIMIT-R03 and @realizes LIMIT-R04 @TBD: wording\n" +
		"\n" +
		"@realizes LIMIT-R05\n"
	want := []realizesOnLine{
		{from: "CRED-V01", to: "LIMIT-R03", deferred: true},
		{from: "CRED-V01", to: "LIMIT-R04", deferred: true},
		{from: "", to: "LIMIT-R05"},
	}
	if got := parseRealizesWithLines(content); !reflect.DeepEqual(got, want) {
		t.Errorf("a non-blank line keeps the rule, a blank one ends it:\n got %+v\nwant %+v", got, want)
	}
}

func TestSpecRealizesDoctrine_aTagOnTheNextLineDeclares(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{"screen": {RequiresDoctrine: true}}}
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Tags: []string{"screen"}}
	if v, msg := checkSpecRealizesDoctrine("### CRED-V01 — a rule\n  @realizes LIMIT-R03\n", n, "", nil, cfg); v != Pass {
		t.Errorf("the tag below the rule declares for it: %v (%s)", v, msg)
	}
}

func TestLayerRequiresDoctrine_resolvesTheLayerOfTheTarget(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"screen": {Pattern: "screens/**/*.tsx", Kind: "code", RequiresDoctrine: true},
		"lib":    {Pattern: "lib/**/*.ts", Kind: "code"},
	}}

	t.Run("by the specifies edge", func(t *testing.T) {
		spec := mapx.Node{ID: "docs/home.spec.md", Kind: mapx.KindSpec}
		g := &mapx.Graph{
			Nodes: []mapx.Node{spec, {ID: "app/home.tsx", Kind: mapx.KindCode, Tags: []string{"screen"}},
				{ID: "app/other.tsx", Kind: mapx.KindCode, Tags: []string{"screen"}}},
			Edges: []mapx.Edge{{From: spec.ID, To: "app/home.tsx", Type: mapx.EdgeSpecifies}},
		}
		if !layerRequiresDoctrine(spec, "", g, cfg) {
			t.Error("the specified target is a screen")
		}
		g.Edges[0].Type = mapx.EdgeRealizes
		if layerRequiresDoctrine(spec, "", g, cfg) {
			t.Error("only a specifies edge names the target")
		}
		g.Edges[0] = mapx.Edge{From: spec.ID, To: "lib/x.ts", Type: mapx.EdgeSpecifies}
		if layerRequiresDoctrine(spec, "", g, cfg) {
			t.Error("the edge points at a node that is not the screen")
		}
	})

	t.Run("by path, before the code exists", func(t *testing.T) {
		if !layerRequiresDoctrine(mapx.Node{ID: "screens/home.spec.md", Kind: mapx.KindSpec}, "", nil, cfg) {
			t.Error("screens/home.tsx would be a screen")
		}
		if layerRequiresDoctrine(mapx.Node{ID: "lib/x.spec.md", Kind: mapx.KindSpec}, "", nil, cfg) {
			t.Error("lib/x.ts is not a demanding layer")
		}
		if layerRequiresDoctrine(mapx.Node{ID: "screens/home.md", Kind: mapx.KindSpec}, "", nil, cfg) {
			t.Error("a file that is not a spec has no target path")
		}
	})
}

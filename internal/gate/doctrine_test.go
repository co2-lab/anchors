package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
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

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O CASO REAL: um revert apagou a `DTSTD-B10` — spec, codigo e teste — e a `.feature`
// manteve o cenario dela, porque no ponto do revert aquele cenario ainda nao existia e um
// PR paralelo o reintroduziu sem conflito de git.
//
// Sondado antes deste gate: o `spec-feature-match` respondia Pass com mensagem VAZIA.

const featureComOrfao = "@UNITX\nFeature: X\n\n" +
	"  @UNITX-B01\n  Scenario: a que ficou\n\n" +
	"  @UNITX-B10\n  Scenario: a da regra REVERTIDA\n"

func featNodeRev() mapx.Node {
	return mapx.Node{ID: "x.feature", Kind: mapx.KindFeature, Code: "UNITX"}
}

// monta um projeto com a spec ligada a' feature por `covered-by`.
func comSpec(t *testing.T, spec string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.spec.md"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{
		Nodes: []mapx.Node{featNodeRev()},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy}},
	}
}

// --- feature-spec-match ---

func TestFeatureSpecMatch_cenarioSemRegraEAcusado(t *testing.T) {
	t.Run("FSPMT-B01: a scenario whose rule the spec no longer declares is reported", func(t *testing.T) {})
	root, g := comSpec(t, "### UNITX-B01 — a regra que ficou\n")

	v, msg := checkFeatureSpecMatch(featureComOrfao, featNodeRev(), root, g, nil)
	if v != Fail {
		t.Fatalf("o cenario B10 nao tem regra e veio %v", v)
	}
	if !strings.Contains(msg, "UNITX-B10") {
		t.Errorf("o veredito nao nomeia o orfao: %q", msg)
	}
	if strings.Contains(msg, "UNITX-B01") {
		t.Errorf("acusou o cenario que TEM regra: %q", msg)
	}
}

func TestFeatureSpecMatch_todosComRegraPassa(t *testing.T) {
	t.Run("FSPMT-B02: every scenario backed by a declared rule passes", func(t *testing.T) {})
	root, g := comSpec(t, "### UNITX-B01 — ficou\n\n### UNITX-B10 — tambem ficou\n")

	if v, msg := checkFeatureSpecMatch(featureComOrfao, featNodeRev(), root, g, nil); v != Pass {
		t.Errorf("as duas regras existem e veio %v: %s", v, msg)
	}
}

// O cenario cita a regra de OUTRA unidade para dizer contra o que roda. Cobrar isso da
// spec local faria o gate pedir o impossivel — o erro que o `scenario-coverage` ja' mediu,
// com 18 cenarios cobrados de uma spec que definia 6.
func TestFeatureSpecMatch_naoCobraCodigoDeOutraUnidade(t *testing.T) {
	t.Run("FSPMT-X01: a code from another unit is not charged of this spec", func(t *testing.T) {})
	root, g := comSpec(t, "### UNITX-B01 — a regra\n")
	comVizinha := "@UNITX\nFeature: X\n\n  @UNITX-B01 @OUTRA-B07\n  Scenario: roda contra a vizinha\n"

	if v, msg := checkFeatureSpecMatch(comVizinha, featNodeRev(), root, g, nil); v != Pass {
		t.Errorf("o codigo da vizinha foi cobrado: %v / %s", v, msg)
	}
}

func TestFeatureSpecMatch_skips(t *testing.T) {
	naoEFeature := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFeatureSpecMatch(featureComOrfao, naoEFeature, "", nil, nil); v != Skip {
		t.Errorf("nao e' feature e veio %v", v)
	}
	// Sem spec ligada: Pending, porque quem cobra a existencia da spec e' a co-locacao.
	semAresta := &mapx.Graph{Nodes: []mapx.Node{featNodeRev()}}
	if v, _ := checkFeatureSpecMatch(featureComOrfao, featNodeRev(), t.TempDir(), semAresta, nil); v != Pending {
		t.Errorf("sem spec ligada esperava Pending, veio %v", v)
	}
}

// --- test-feature-match ---

func testNodeRev() mapx.Node {
	return mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest, Code: "UNITX"}
}

func comFeature(t *testing.T, feat string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.feature"), []byte(feat), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{
		Nodes: []mapx.Node{testNodeRev()},
		Edges: []mapx.Edge{{From: "x.feature", To: "x.test.ts", Type: mapx.EdgeTestedBy}},
	}
}

// Um teste verde sobre regra revertida e' pior que um teste ausente, porque ele ATESTA.
func TestTestFeatureMatch_testeProvaCodigoSemCenario(t *testing.T) {
	t.Run("TFTMT-B01: a code the test claims to prove and no scenario declares is reported", func(t *testing.T) {})
	root, g := comFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: a que ficou\n")
	teste := "it('UNITX-B01 ok', () => {})\nit('UNITX-B10 orfao', () => {})\n"

	v, msg := checkTestFeatureMatch(teste, testNodeRev(), root, g, nil)
	if v != Fail {
		t.Fatalf("o teste prova B10 sem cenario e veio %v", v)
	}
	if !strings.Contains(msg, "UNITX-B10") {
		t.Errorf("o veredito nao nomeia o orfao: %q", msg)
	}
}

// COMENTARIO FORA: um codigo citado em comentario e' referencia, nao prova.
func TestTestFeatureMatch_comentarioNaoConta(t *testing.T) {
	t.Run("TFTMT-X01: a code cited only in a comment is not a claim of proof", func(t *testing.T) {})
	root, g := comFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: a que ficou\n")
	teste := "// UNITX-B10 foi revertida, ver #811\nit('UNITX-B01 ok', () => {})\n"

	if v, msg := checkTestFeatureMatch(teste, testNodeRev(), root, g, nil); v != Pass {
		t.Errorf("o codigo em COMENTARIO foi cobrado: %v / %s", v, msg)
	}
}

func TestTestFeatureMatch_naoCobraCodigoDeOutraUnidade(t *testing.T) {
	root, g := comFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: a que ficou\n")
	teste := "it('UNITX-B01 ok', () => { montaFixture(OUTRA-B07) })\n"

	if v, msg := checkTestFeatureMatch(teste, testNodeRev(), root, g, nil); v != Pass {
		t.Errorf("o codigo da vizinha foi cobrado: %v / %s", v, msg)
	}
}

func TestTestFeatureMatch_skips(t *testing.T) {
	naoETeste := mapx.Node{ID: "x.feature", Kind: mapx.KindFeature}
	if v, _ := checkTestFeatureMatch("", naoETeste, "", nil, nil); v != Skip {
		t.Errorf("nao e' teste e veio %v", v)
	}
	semAresta := &mapx.Graph{Nodes: []mapx.Node{testNodeRev()}}
	if v, _ := checkTestFeatureMatch("it('x')", testNodeRev(), t.TempDir(), semAresta, nil); v != Pending {
		t.Errorf("sem feature ligada esperava Pending, veio %v", v)
	}
}

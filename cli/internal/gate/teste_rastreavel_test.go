package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaRastreavel(t *testing.T, teste, feature string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.feature"), []byte(feature), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "u.feature", Kind: mapx.KindFeature}, {ID: "u.test.ts", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{{From: "u.feature", To: "u.test.ts", Type: mapx.EdgeTestedBy}},
	}
	return checkTesteRastreavel(teste, mapx.Node{ID: "u.test.ts", Kind: mapx.KindTest}, root, g, nil)
}

const featFix = "@PHA1X-S01 @nivel-unit\nCenário: a\n@PHA1X-A01 @nivel-unit\nCenário: b\n"

// O caso real: o teste PROVA o comportamento, mas com o código TROCADO (`PH1A` em vez de
// `PHA1X` — dígito e letra invertidos). Três telas tinham teste e nenhuma estava
// rastreada; a busca por código não as alcançava.
func TestTesteRastreavel_codigoTrocadoReprova(t *testing.T) {
	teste := "it('PH1A: renderiza e avança', () => {})"
	v, msg := rodaRastreavel(t, teste, featFix)
	if v != Fail {
		t.Fatalf("código trocado deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "PHA1X-S01") {
		t.Errorf("a mensagem deve mostrar os códigos esperados: %s", msg)
	}
}

// O outro caso real: o teste cobre 5 comportamentos e não cita código nenhum.
func TestTesteRastreavel_semCodigoReprova(t *testing.T) {
	teste := "it('addAccount adiciona com id gerado', () => {})"
	if v, _ := rodaRastreavel(t, teste, featFix); v != Fail {
		t.Errorf("teste sem código é invisível aos gates relacionais: %v", v)
	}
}

// A régua é a mais FRACA possível: um código basta. Cobrar um por caso é do
// `feature-test-match`, e duplicá-lo aqui faria o mesmo débito aparecer duas vezes.
func TestTesteRastreavel_umCodigoBasta(t *testing.T) {
	teste := "// @nivel-unit — PHA1X-S01\nit('PHA1X-S01: renderiza', () => {})"
	if v, msg := rodaRastreavel(t, teste, featFix); v != Pass {
		t.Errorf("um código já torna o teste rastreável: %v (%s)", v, msg)
	}
}

// Teste SEM feature ligada não tem cenário a citar — exigir código dele seria pedir
// referência a nada.
func TestTesteRastreavel_semFeaturePula(t *testing.T) {
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.test.ts", Kind: mapx.KindTest}}}
	v, _ := checkTesteRastreavel("it('x', () => {})",
		mapx.Node{ID: "u.test.ts", Kind: mapx.KindTest}, t.TempDir(), g, nil)
	if v != Skip {
		t.Errorf("sem feature não há o que citar: %v", v)
	}
}

// Feature sem código de cenário também não dá o que citar.
func TestTesteRastreavel_featureSemCodigoPula(t *testing.T) {
	if v, _ := rodaRastreavel(t, "it('x', () => {})", "Funcionalidade: X\n"); v != Skip {
		t.Errorf("feature sem código não exige citação: %v", v)
	}
}

// A cobrança é do TESTE — rodar sobre a feature acusaria o arquivo errado.
func TestTesteRastreavel_soRodaSobreTeste(t *testing.T) {
	n := mapx.Node{ID: "u.feature", Kind: mapx.KindFeature}
	if v, _ := checkTesteRastreavel("", n, "", &mapx.Graph{}, nil); v != Skip {
		t.Errorf("a rastreabilidade é do teste: %v", v)
	}
}

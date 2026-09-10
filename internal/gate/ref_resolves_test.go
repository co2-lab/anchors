package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaRefResolves(t *testing.T, arquivo, conteudo, specNome, specConteudo string) (Verdict, string) {
	t.Helper()
	dir := t.TempDir()
	if specNome != "" {
		if err := os.WriteFile(filepath.Join(dir, specNome), []byte(specConteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	kind := mapx.KindCode
	if strings.HasSuffix(arquivo, ".feature") {
		kind = mapx.KindFeature
	} else if strings.Contains(arquivo, ".test.") {
		kind = mapx.KindTest
	}
	return checkRefResolves(conteudo, mapx.Node{ID: arquivo, Kind: kind}, dir, nil, nil)
}

// O caso real: 49 arquivos de modelo com `ref: DTAXX` — a identidade de quando os modelos
// viviam num arquivo só. Depois da desfusão cada um ganhou spec própria, e nenhum `ref:`
// foi propagado. `header-conforme` ficou verde nos 49: ele confere que o campo EXISTE.
func TestRefResolvesRefactoracaoNaoPropagada(t *testing.T) {
	v, d := rodaRefResolves(t,
		"TaxReceipt.ts", "// @anchors\n//   ref: DTAXX\n//   layer: schema-model\n",
		"TaxReceipt.spec.md", "<!-- @anchors\n  code: TRT1X\n-->\n")
	if v != Fail {
		t.Fatalf("ref divergente deveria reprovar, foi %s (%s)", v, d)
	}
	for _, esperado := range []string{"DTAXX", "TRT1X"} {
		if !strings.Contains(d, esperado) {
			t.Errorf("a mensagem não mostra o par (falta %q): %s", esperado, d)
		}
	}
}

func TestRefResolvesQuandoCasa(t *testing.T) {
	casos := map[string][2]string{
		"código":  {"TaxReceipt.ts", "// @anchors\n//   ref: TRT1X\n"},
		"feature": {"TaxReceipt.feature", "# @anchors\n#   ref: TRT1X\n"},
		"teste":   {"TaxReceipt.test.ts", "// @anchors\n//   ref: TRT1X\n"},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			v, d := rodaRefResolves(t, c[0], c[1], "TaxReceipt.spec.md", "<!-- @anchors\n  code: TRT1X\n-->\n")
			if v != Pass {
				t.Fatalf("ref correto deveria passar, foi %s (%s)", v, d)
			}
		})
	}
}

// Cada gate acusa UMA coisa: sem spec irmã, quem cobra é o trinca-completa. Dois gates
// sobre o mesmo defeito viram ruído e o usuário desliga os dois.
func TestRefResolvesSemSpecIrmaEhDoOutroGate(t *testing.T) {
	if v, d := rodaRefResolves(t, "Solto.ts", "// @anchors\n//   ref: XXXXX\n", "", ""); v != Skip {
		t.Fatalf("sem spec irmã deveria ser Skip, foi %s (%s)", v, d)
	}
}

// Sem `ref:` declarado, a ausência é do header-conforme — não deste gate.
func TestRefResolvesSemRefEhDoHeaderConforme(t *testing.T) {
	if v, _ := rodaRefResolves(t, "X.ts", "// nada aqui\n", "X.spec.md", "<!-- @anchors\n  code: AAAAX\n-->\n"); v != Skip {
		t.Fatalf("sem ref deveria ser Skip, foi %s", v)
	}
}

// A spec é DONA (`code:`), não referencia — o gate não se aplica a ela.
func TestRefResolvesNaoSeAplicaASpec(t *testing.T) {
	dir := t.TempDir()
	v, _ := checkRefResolves("<!-- @anchors\n  code: AAAAX\n-->\n",
		mapx.Node{ID: "X.spec.md", Kind: mapx.KindSpec}, dir, nil, nil)
	if v != Skip {
		t.Fatalf("a spec é dona do código — deveria ser Skip, foi %s", v)
	}
}

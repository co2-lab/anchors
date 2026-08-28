package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaCatalogado(t *testing.T, spec, codigo string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}, {ID: "u.ts", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "u.spec.md", To: "u.ts", Type: mapx.EdgeSpecifies}},
	}
	return checkCodigoCatalogado(spec, mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}, root, g, nil)
}

// O caso real: a spec catalogava 2 regras para 7 funções exportadas, e nenhum gate
// perguntou pelas 5 restantes. É o inverso do `regra-implementada`.
func TestCodigoCatalogado_simboloForaDoCatalogoReprova(t *testing.T) {
	spec := "| `INVAX-B01` | `calcStockoutRisk` classifica o risco |\n"
	codigo := `export function calcStockoutRisk() {}
export function calcBestMonthToBuy() {}
export function calcPriceVariation() {}`

	v, msg := rodaCatalogado(t, spec, codigo)
	if v != Fail {
		t.Fatalf("símbolo fora do catálogo deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "calcBestMonthToBuy") {
		t.Errorf("a mensagem deve nomear o órfão: %s", msg)
	}
	if strings.Contains(msg, "calcStockoutRisk") {
		t.Errorf("o que a spec cataloga não pode ser acusado: %s", msg)
	}
}

// A dispensa por SÍMBOLO fecha o caso legítimo: nem toda exportação merece regra.
func TestCodigoCatalogado_noRuleDispensa(t *testing.T) {
	spec := "| `INVAX-B01` | `calcRisco` classifica |\n"
	codigo := `export function calcRisco() {}
// @no-rule: formatação pura, sem decisão de negócio
export function formatarMoeda() {}`

	if v, msg := rodaCatalogado(t, spec, codigo); v != Pass {
		t.Errorf("`@no-rule` com razão dispensa o símbolo: %v (%s)", v, msg)
	}
}

// Marcador NU não dispensa — senão vira um jeito silencioso de calar o gate.
func TestCodigoCatalogado_noRuleExigeRazao(t *testing.T) {
	spec := "| `INVAX-B01` | `calcRisco` classifica |\n"
	codigo := `export function calcRisco() {}
// @no-rule
export function formatarMoeda() {}`

	if v, _ := rodaCatalogado(t, spec, codigo); v != Fail {
		t.Error("marcador sem razão não pode dispensar")
	}
}

// Spec que cataloga tudo passa.
func TestCodigoCatalogado_tudoCatalogadoPassa(t *testing.T) {
	spec := "| `X-B01` | `a` faz |\n| `X-B02` | `b` faz |\n"
	if v, msg := rodaCatalogado(t, spec, "export const a = 1\nexport const b = 2\n"); v != Pass {
		t.Errorf("tudo catalogado deveria passar: %v (%s)", v, msg)
	}
}

// Sem código ligado a ausência é de outro gate — acusar nos dois duplicaria o débito.
func TestCodigoCatalogado_semCodigoLigadoPula(t *testing.T) {
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}}}
	v, _ := checkCodigoCatalogado("| `X-B01` | x |", mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec},
		t.TempDir(), g, nil)
	if v != Skip {
		t.Errorf("ausência de código é do trinca-completa: %v", v)
	}
}

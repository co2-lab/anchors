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

// O `@no-rule` VALE NO BLOCO DE COMENTÁRIO acima do símbolo, e não só na linha dele.
//
// A declaração é documentação: quem a escreve põe junto com a explicação, e a explicação
// raramente cabe numa linha. O gate olhava só a linha anterior, então uma declaração na
// segunda linha de um comentário de duas era ignorada — e ele seguia acusando.
//
// Medido implementando um adaptador: três símbolos declarados e três acusados. Pior, o
// mesmo padrão PASSAVA noutro arquivo por acaso — o nome do símbolo aparecia no texto da
// spec, então a declaração nunca era lida e ninguém notava que ela não funcionava.
func TestNoRuleValeNoComentarioAcima(t *testing.T) {
	casos := map[string]string{
		"na mesma linha":       "export function x() {} // @no-rule: porta de saída\n",
		"uma linha acima":      "// @no-rule: porta de saída\nexport function x() {}\n",
		"bloco de duas linhas": "// A porta existe para o teste não alcançar a rede;\n// @no-rule: o que trafega é regra de quem a usa\nexport function x() {}\n",
		"bloco com parágrafos": "// Explicação longa.\n//\n// @no-rule: sem comportamento próprio\nexport function x() {}\n",
		"doc comment":          "/**\n * @no-rule: forma de entrada\n */\nexport function x() {}\n",
	}
	for nome, codigo := range casos {
		simbolos := simbolosComLinha(codigo)
		if len(simbolos) == 0 {
			t.Fatalf("%s: nenhum símbolo reconhecido", nome)
		}
		if !noRuleRE.MatchString(simbolos[0].linha) {
			t.Errorf("%s: a declaração deveria valer, e o gate a ignorou.\ncontexto lido:\n%s",
				nome, simbolos[0].linha)
		}
	}
}

// A subida NÃO pode engolir o arquivo: um símbolo sem comentário acima não herda a
// declaração de outro símbolo mais acima. Herdar faria um `@no-rule` isentar o arquivo
// inteiro, que é o oposto do que ele é.
func TestNoRuleNaoVazaEntreSimbolos(t *testing.T) {
	codigo := "// @no-rule: este sim\nexport function comDeclaracao() {}\n\nexport function semDeclaracao() {}\n"
	simbolos := simbolosComLinha(codigo)
	if len(simbolos) != 2 {
		t.Fatalf("esperava 2 símbolos, veio %d", len(simbolos))
	}
	if !noRuleRE.MatchString(simbolos[0].linha) {
		t.Error("o primeiro tem declaração e deveria valer")
	}
	if noRuleRE.MatchString(simbolos[1].linha) {
		t.Error("o segundo NÃO tem declaração — herdar a do primeiro isentaria o arquivo " +
			"inteiro com um marcador só")
	}
}

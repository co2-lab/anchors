package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
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
	cfg := &config.Config{Derived: &config.Derived{ExportDetect: exportedREDefaultTS}}
	return checkCodeCataloged(spec, mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}, root, g, cfg)
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
	v, _ := checkCodeCataloged("| `X-B01` | x |", mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec},
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
		simbolos := symbolsWithLine(codigo, regexp.MustCompile(exportedREDefaultTS))
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
	simbolos := symbolsWithLine(codigo, regexp.MustCompile(exportedREDefaultTS))
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

// rodaCatalogadoCfg é a variante que passa config — o gate agnóstico depende dela.
func rodaCatalogadoCfg(t *testing.T, spec, codigo, arquivo string, cfg *config.Config) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, arquivo), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}, {ID: arquivo, Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "u.spec.md", To: arquivo, Type: mapx.EdgeSpecifies}},
	}
	return checkCodeCataloged(spec, mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}, root, g, cfg)
}

// O GATE NAO PODE APROVAR O QUE NAO SABE LER.
//
// O padrao de export estava EMBUTIDO como sintaxe TS/JS. Num projeto Go, Python ou Ruby
// o regex casava zero simbolos e o gate reportava VERDE -- sendo `blocking: true`,
// carimbava aprovacao sobre o que nao tinha conferido. E a pior falha possivel num
// medidor, e a mesma que o comentario do `mock_detect` ja condenava por escrito.
//
// Sem `export_detect` declarado, o gate PULA e DIZ que pulou.
func TestCodeCatalogedPulaQuandoNaoSabeLerALinguagem(t *testing.T) {
	// Go: `func Publica()` nao casa nenhuma sintaxe de export TS/JS.
	codigo := "package u\n\nfunc Publica() int { return 1 }\n"
	v, msg := rodaCatalogadoCfg(t, "# U\n\n## UUUUU-B01 — algo\n", codigo, "u.go", nil)
	if v == Pass {
		t.Errorf("o gate APROVOU um arquivo que nao sabe ler (%s) — verde sobre o que "+
			"nao conferiu e' pior que vermelho honesto", msg)
	}
	if v != Skip && v != Pending {
		t.Errorf("esperava Skip/Pending, veio %v: %s", v, msg)
	}
	if !strings.Contains(msg, "export_detect") {
		t.Errorf("a mensagem nao diz COMO habilitar (`export_detect`): %q", msg)
	}
}

// Com o padrao declarado, o gate confronta de verdade -- em qualquer linguagem.
func TestCodeCatalogedUsaOPadraoDoProjeto(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{ExportDetect: `(?m)^func\s+([A-Z]\w*)`}}

	codigo := "package u\n\nfunc Publica() int { return 1 }\n"
	v, msg := rodaCatalogadoCfg(t, "# U\n\n## UUUUU-B01 — algo\n", codigo, "u.go", cfg)
	if v != Fail {
		t.Errorf("`Publica` nao esta na spec e o gate nao reprovou (%v): %s", v, msg)
	}
	if !strings.Contains(msg, "Publica") {
		t.Errorf("a mensagem nao nomeia o simbolo: %q", msg)
	}
}

func TestCodeCatalogedUsaDialetoFamily(t *testing.T) {
	cfg := &config.Config{Dialect: &config.Dialect{Family: "go"}}

	codigo := "package u\n\nfunc Publica() int { return 1 }\n"
	v, msg := rodaCatalogadoCfg(t, "# U\n\n## UUUUU-B01 — algo\n", codigo, "u.go", cfg)
	if v != Fail {
		t.Errorf("`Publica` em Go via dialect family deveria reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "Publica") {
		t.Errorf("a mensagem deveria nomear Publica: %q", msg)
	}
}

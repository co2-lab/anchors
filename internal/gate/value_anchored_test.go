package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaAncora(t *testing.T, spec, codigo string, cfg *config.Config) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}, {ID: "u.ts", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "u.spec.md", To: "u.ts", Type: mapx.EdgeSpecifies}},
	}
	return checkValueAnchored(spec, mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}, root, g, cfg)
}

func cfgAncora() *config.Config {
	return &config.Config{Derived: &config.Derived{
		ExportDetect: exportedREDefaultTS,
		ValueAnchor:  `@code-reference-\[([^\]]+)\]-\[([^\]]+)\]`,
	}}
}

// CADA VALOR DE UM CONJUNTO FECHADO E UMA DECISAO -- e uma dispensa sobre o SIMBOLO
// libera todos de uma vez.
//
// MEDIDO: o `QueryScope` do blue-eyes declara `JANELAS = ['15m','1h','6h','24h']` com um
// `@no-rule` na linha do `export`. Os quatro valores nunca foram confrontados
// individualmente, e duas telas passaram a declarar `5m`, `30m` e `1d` -- que o contrato
// NAO aceita. O backend cai no padrao `1h` em silencio: a tela mostra `5m` exibindo dados
// de uma hora.
func TestValorSemAncoraReprova(t *testing.T) {
	t.Run("VLANV-B02: A value of a closed set with no anchor is failed", func(t *testing.T) {})
	t.Run("VLANV-B03: The verdict names the unanchored value", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n  '15m',\n  '1h',\n] as const;\n"
	v, msg := rodaAncora(t, "# U\n\n## UUUUU-B01 — janelas\n", codigo, cfgAncora())
	if v != Fail {
		t.Fatalf("valor sem ancora deveria reprovar, veio %v: %s", v, msg)
	}
	if !strings.Contains(msg, "15m") {
		t.Errorf("a mensagem nao nomeia o valor sem ancora: %q", msg)
	}
}

// A ancora completa vale: chave da regra + valor esperado.
func TestValorComAncoraPassa(t *testing.T) {
	t.Run("VLANV-B04: A value whose anchor carries rule key and value passes", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n" +
		"  // @code-reference-[WINDOW-001]-[15m]\n  '15m',\n" +
		"  // @code-reference-[WINDOW-002]-[1h]\n  '1h',\n] as const;\n"
	if v, msg := rodaAncora(t, "# U\n\n## UUUUU-B01 — janelas\n", codigo, cfgAncora()); v != Pass {
		t.Errorf("ancora completa deveria passar, veio %v: %s", v, msg)
	}
}

// O SEGUNDO COLCHETE E O QUE SEPARA CITAR DE PROVAR.
//
// Uma ancora que so aponta a regra apodrece em silencio: alguem troca o valor e o
// comentario continua parecendo correto. Com o valor escrito nela, o gate confronta o que
// a ancora AFIRMA contra o que a linha DIZ -- e isso nao depende de linguagem, porque
// compara duas partes do proprio comentario com a linha que ele anota.
func TestAncoraQueMenteSobreOValorReprova(t *testing.T) {
	t.Run("VLANV-B05: An anchor that asserts one value while the line says another is failed", func(t *testing.T) {})
	t.Run("VLANV-B06: The verdict of a lying anchor shows both sides of the divergence", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n" +
		"  // @code-reference-[WINDOW-001]-[15m]\n  '5m',\n] as const;\n"
	v, msg := rodaAncora(t, "# U\n\n## UUUUU-B01 — janelas\n", codigo, cfgAncora())
	if v != Fail {
		t.Fatalf("a ancora afirma `15m` e a linha diz `5m` — deveria reprovar, veio %v: %s", v, msg)
	}
	if !strings.Contains(msg, "15m") || !strings.Contains(msg, "5m") {
		t.Errorf("a mensagem precisa mostrar OS DOIS lados da divergencia: %q", msg)
	}
}

// Sem `value_anchor` declarado o gate PULA -- nao sabe ler, nao aprova.
func TestValorAncoradoPulaSemPadrao(t *testing.T) {
	t.Run("VLANV-B08: Without a declared value anchor pattern the gate skips", func(t *testing.T) {})
	t.Run("VLANV-B09: The skip names the setting that enables the gate", func(t *testing.T) {})
	codigo := "export const JANELAS = ['15m'] as const;\n"
	v, msg := rodaAncora(t, "# U\n", codigo, nil)
	if v == Pass || v == Fail {
		t.Errorf("sem padrao declarado o gate nao pode julgar (%v): %s", v, msg)
	}
	if !strings.Contains(msg, "value_anchor") {
		t.Errorf("a mensagem nao diz COMO habilitar: %q", msg)
	}
}

// O gate alcanca o codigo ATRAVES da spec -- sobre os outros tipos ele nao tem jurisdicao
// propria.
func TestValorAncoradoSoOlhaSpec(t *testing.T) {
	t.Run("VLANV-B01: An artifact that is not a spec leaves without a verdict", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.ts", Kind: mapx.KindCode}}}
	n := mapx.Node{ID: "u.ts", Kind: mapx.KindCode}
	if v, msg := checkValueAnchored("", n, t.TempDir(), g, cfgAncora()); v != Skip {
		t.Errorf("no de codigo deveria pular, veio %v: %s", v, msg)
	}
}

// A ANCORA QUE MENTE vem primeiro: e pior que a ausente. A ausente se ve; esta parece
// rastreabilidade e aponta para o lugar errado.
func TestAncoraMentirosaVemAntesDaAusente(t *testing.T) {
	t.Run("VLANV-B07: Lying anchors are reported before the unanchored ones", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n" +
		"  // @code-reference-[WINDOW-001]-[15m]\n  '5m',\n" +
		"  '24h',\n] as const;\n"
	v, msg := rodaAncora(t, "# U\n", codigo, cfgAncora())
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v: %s", v, msg)
	}
	iMentira := strings.Index(msg, "5m")
	i24 := strings.Index(msg, "24h")
	if iMentira < 0 || i24 < 0 {
		t.Fatalf("o laudo precisa citar os dois achados: %q", msg)
	}
	if iMentira > i24 {
		t.Errorf("a ancora mentirosa tem de vir ANTES da ausente:\n%s", msg)
	}
}

// Um `export const X = 1` nao abre lista: nao e conjunto fechado e nao ha o que ancorar.
func TestEscalarNaoEhConjuntoFechado(t *testing.T) {
	t.Run("VLANV-B10: A declaration that opens no list is not a closed set", func(t *testing.T) {})
	codigo := "export const JANELA_PADRAO = '1h';\n"
	if v, msg := rodaAncora(t, "# U\n", codigo, cfgAncora()); v != Pass {
		t.Errorf("escalar nao e conjunto, deveria passar, veio %v: %s", v, msg)
	}
}

// Sem o SEGUNDO grupo a ancora nao afirma valor nenhum -- e o confronto que da razao ao
// gate nao pode acontecer. Vale como padrao nao declarado.
func TestPadraoComUmGrupoNaoHabilita(t *testing.T) {
	t.Run("VLANV-I01: An anchor pattern with a single capture group does not enable the gate", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{
		ExportDetect: exportedREDefaultTS,
		ValueAnchor:  `@code-reference-\[([^\]]+)\]`,
	}}
	codigo := "export const JANELAS = [\n  '15m',\n] as const;\n"
	v, msg := rodaAncora(t, "# U\n", codigo, cfg)
	if v == Pass || v == Fail {
		t.Errorf("padrao de um grupo nao e verificavel, o gate nao pode julgar (%v): %s", v, msg)
	}
}

// A ancora `@code-reference-[WINDOW-001]-[15m]` carrega dois `]`, e trata-los como fim de
// lista encerrava o conjunto na primeira ancora: o gate saia do bloco e nao confrontava
// valor nenhum. Um caso que deveria reprovar passava, porque o confronto nunca acontecia.
func TestAncoraNaoFechaOConjunto(t *testing.T) {
	t.Run("VLANV-I02: A line carrying an anchor is never read as the end of the list", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n" +
		"  // @code-reference-[WINDOW-001]-[15m]\n  '15m',\n" +
		"  // @code-reference-[WINDOW-002]-[1h]\n  '6h',\n] as const;\n"
	v, msg := rodaAncora(t, "# U\n", codigo, cfgAncora())
	if v != Fail {
		t.Fatalf("o segundo valor mente e tem de ser alcancado, veio %v: %s", v, msg)
	}
	if !strings.Contains(msg, "6h") {
		t.Errorf("o confronto parou na primeira ancora — o segundo valor nao foi visto:\n%s", msg)
	}
}

// Aprovar sem poder olhar carimbaria o que nao foi medido.
func TestValorAncoradoSemMapaNaoAprova(t *testing.T) {
	t.Run("VLANV-I03: With no built map the verdict is pending", func(t *testing.T) {})
	n := mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}
	if v, msg := checkValueAnchored("# U\n", n, t.TempDir(), nil, cfgAncora()); v == Pass {
		t.Errorf("sem mapa o gate nao pode aprovar, veio %v: %s", v, msg)
	}
}

// A regua e que a decisao TENHA endereco e que o endereco NAO MINTA. Se o valor pertence
// ao dominio e julgamento, e julgamento e de quem conhece o produto.
func TestNaoJulgaSeOValorEhBom(t *testing.T) {
	t.Run("VLANV-X01: The gate does not judge whether the value is a good one", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n" +
		"  // @code-reference-[WINDOW-999]-[999y]\n  '999y',\n] as const;\n"
	if v, msg := rodaAncora(t, "# U\n", codigo, cfgAncora()); v != Pass {
		t.Errorf("valor absurdo mas ancorado e coerente deveria passar, veio %v: %s", v, msg)
	}
}

// A extracao do literal e deliberadamente simples. Falso-negativo aqui e melhor que
// falso-positivo em massa, que treina o time a ignorar o gate.
func TestNaoAcusaLinhaSemLiteral(t *testing.T) {
	t.Run("VLANV-X02: The gate does not accuse a line whose literal it cannot read", func(t *testing.T) {})
	codigo := "export const JANELAS = [\n  DEFAULT_WINDOW,\n  computeWindow(base),\n] as const;\n"
	if v, msg := rodaAncora(t, "# U\n", codigo, cfgAncora()); v != Pass {
		t.Errorf("linha sem literal legivel nao e acusada, veio %v: %s", v, msg)
	}
}

// As duas formas -- o que e simbolo publico e o que e ancora -- sao declaradas pelo
// projeto. Um gate que as inventasse cobraria uma convencao que ninguem adotou.
func TestNaoInventaAFormaDoSimboloPublico(t *testing.T) {
	t.Run("VLANV-X03: The gate does not decide what an anchor or a public symbol looks like", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{
		ExportDetect: `^public\s+static\s+final\s+\w+\s+(\w+)`, // outra linguagem: nao casa com o `export const`
		ValueAnchor:  `@code-reference-\[([^\]]+)\]-\[([^\]]+)\]`,
	}}
	codigo := "export const JANELAS = [\n  '15m',\n] as const;\n"
	if v, msg := rodaAncora(t, "# U\n", codigo, cfg); v != Pass {
		t.Errorf("o gate nao acha o conjunto pelo padrao declarado — nao acusa, veio %v: %s", v, msg)
	}
}

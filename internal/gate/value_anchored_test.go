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
	codigo := "export const JANELAS = ['15m'] as const;\n"
	v, msg := rodaAncora(t, "# U\n", codigo, nil)
	if v == Pass || v == Fail {
		t.Errorf("sem padrao declarado o gate nao pode julgar (%v): %s", v, msg)
	}
	if !strings.Contains(msg, "value_anchor") {
		t.Errorf("a mensagem nao diz COMO habilitar: %q", msg)
	}
}

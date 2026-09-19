package gate

import (
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestRouteDeclared(t *testing.T) {
	screen := func(body string) string {
		return "<!-- @anchors\n  code: HOMEX\n  layer: screen\n-->\n" + body
	}
	cases := []struct {
		name    string
		content string
		want    Verdict
	}{
		{
			"tela com rota e navegação concreta",
			screen("> **Rota**: `Home`\n\n### Entrada\n| Origem | Tela |\n| --- | --- |\n| MainTabs | HomeScreen |\n"),
			Pass,
		},
		{
			"tela sem rota",
			screen("## Visão Geral\nsem linha de rota\n"),
			Fail,
		},
		{
			"tela com rota mas navegação genérica",
			screen("> **Rota**: `Home`\n\n### Saída\n| Destino | Tela |\n| --- | --- |\n| botão | Próxima tela |\n"),
			Fail,
		},
		{
			"hook não é cobrado (Skip)",
			"<!-- @anchors\n  layer: hook\n-->\n## useAuth\nsem rota, tudo bem\n",
			Skip,
		},
		{
			"business-logic não é cobrado (Skip)",
			"<!-- @anchors\n  layer: business-logic\n-->\n### FOO-B01\n",
			Skip,
		},
		{
			"English screen with route and concrete navigation",
			screen("> **Route**: `Home`\n\n### In\n| Origin | Screen |\n| --- | --- |\n| MainTabs | HomeScreen |\n"),
			Pass,
		},
		{
			"English screen with route but generic navigation",
			screen("> **Route**: `Home`\n\n### Out\n| Target | Screen |\n| --- | --- |\n| button | Next screen |\n"),
			Fail,
		},
	}
	for _, c := range cases {
		got, _ := checkRouteDeclared(c.content, mapx.Node{})
		if got != c.want {
			t.Errorf("%s: checkRouteDeclared = %v, quer %v", c.name, got, c.want)
		}
	}
}

func TestLayerOf_headerWinsOverTags(t *testing.T) {
	content := "<!-- @anchors\n  layer: screen\n-->\n"
	if got := layerOf(mapx.Node{Tags: []string{"hook"}}, content); got != "screen" {
		t.Errorf("layerOf = %q, quer screen (header manda)", got)
	}
	// header omisso → cai nas tags do nó
	if got := layerOf(mapx.Node{Tags: []string{"screen"}}, "sem header\n"); got != "screen" {
		t.Errorf("layerOf (fallback tag) = %q, quer screen", got)
	}
}

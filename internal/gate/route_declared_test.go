package gate

import (
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestRouteDeclared(t *testing.T) {
	t.Run("RTDCL-B01: an artifact that is not a screen leaves without a verdict", func(t *testing.T) {})
	t.Run("RTDCL-B02: a screen with no named route fails", func(t *testing.T) {})
	t.Run("RTDCL-B03: a navigation row carrying a generic term fails", func(t *testing.T) {})
	t.Run("RTDCL-B04: a screen with a named route and concrete neighbours passes", func(t *testing.T) {})
	t.Run("RTDCL-B05: route and navigation are recognised in either declared language", func(t *testing.T) {})
	screen := func(body string) string {
		return "<!-- @anchors\n  code: HOMEX\n  layer: screen\n-->\n" + body
	}
	cases := []struct {
		name    string
		content string
		want    Verdict
	}{
		{
			"RTDCL-B04: tela com rota e navegação concreta",
			screen("> **Rota**: `Home`\n\n### Entrada\n| Origem | Tela |\n| --- | --- |\n| MainTabs | HomeScreen |\n"),
			Pass,
		},
		{
			"RTDCL-B02: tela sem rota",
			screen("## Visão Geral\nsem linha de rota\n"),
			Fail,
		},
		{
			"RTDCL-B03: tela com rota mas navegação genérica",
			screen("> **Rota**: `Home`\n\n### Saída\n| Destino | Tela |\n| --- | --- |\n| botão | Próxima tela |\n"),
			Fail,
		},
		{
			"RTDCL-B01: hook não é cobrado (Skip)",
			"<!-- @anchors\n  layer: hook\n-->\n## useAuth\nsem rota, tudo bem\n",
			Skip,
		},
		{
			"RTDCL-B01: business-logic não é cobrado (Skip)",
			"<!-- @anchors\n  layer: business-logic\n-->\n### FOO-B01\n",
			Skip,
		},
		{
			"RTDCL-B05: English screen with route and concrete navigation",
			screen("> **Route**: `Home`\n\n### In\n| Origin | Screen |\n| --- | --- |\n| MainTabs | HomeScreen |\n"),
			Pass,
		},
		{
			"RTDCL-B05: English screen with route but generic navigation",
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

// RTDCL-I01: the header's declared layer is the source of truth of identity; the node's
// tags are only the fallback. The header is what the author wrote on purpose.
func TestLayerOf_headerWinsOverTags(t *testing.T) {
	t.Run("RTDCL-I01: the header's declared layer wins over the node's tags", func(t *testing.T) {})
	content := "<!-- @anchors\n  layer: screen\n-->\n"
	if got := layerOf(mapx.Node{Tags: []string{"hook"}}, content); got != "screen" {
		t.Errorf("layerOf = %q, quer screen (header manda)", got)
	}
	// header omisso → cai nas tags do nó
	if got := layerOf(mapx.Node{Tags: []string{"screen"}}, "sem header\n"); got != "screen" {
		t.Errorf("layerOf (fallback tag) = %q, quer screen", got)
	}
}

// RTDCL-I02: only TABLE ROWS of the navigation sections declare edges. A sentence
// mentioning a generic destination is the author writing, not an edge being declared —
// accusing it would charge prose for the shape of a table.
func TestRouteDeclared_I02_genericTermInProseIsNotAccused(t *testing.T) {
	t.Run("RTDCL-I02: a generic term in prose is not accused", func(t *testing.T) {})
	content := "<!-- @anchors\n  code: HOMEX\n  layer: screen\n-->\n" +
		"> **Route**: `Home`\n\n### Out\n" +
		"The user reaches the next screen from here, once the form is valid.\n\n" +
		"| Target | Screen |\n| --- | --- |\n| button | DetailScreen |\n"
	if v, msg := checkRouteDeclared(content, mapx.Node{}); v != Pass {
		t.Errorf("a generic term in prose must not be accused; got %v (%s)", v, msg)
	}
}

// RTDCL-X01: no layer other than `screen` is charged for a route. This was the legacy
// validator's vice — it knew only screen and component, so every non-component was
// treated as a screen and every unit that legitimately has no route was accused.
func TestRouteDeclared_X01_onlyScreenIsCharged(t *testing.T) {
	t.Run("RTDCL-X01: no layer other than screen is charged for a route", func(t *testing.T) {})
	for _, layer := range []string{"hook", "business-logic", "store", "dao", "component", "gate"} {
		content := "<!-- @anchors\n  layer: " + layer + "\n-->\n## no route here, on purpose\n"
		v, msg := checkRouteDeclared(content, mapx.Node{})
		if v != Skip {
			t.Errorf("layer %q: expected Skip, got %v (%s)", layer, v, msg)
		}
		if msg == "" {
			t.Errorf("layer %q: the Skip must state its reason", layer)
		}
	}
}

// RTDCL-X02: the ruler is the spec's INTERNAL coherence — a route is named and the
// neighbours are concrete. Confronting that name against the real routing table is a
// different confrontation, over a different artifact.
func TestRouteDeclared_X02_routeIsNotConfrontedAgainstRouter(t *testing.T) {
	t.Run("RTDCL-X02: the declared route is not confronted against the real router", func(t *testing.T) {})
	content := "<!-- @anchors\n  code: HOMEX\n  layer: screen\n-->\n" +
		"> **Route**: `ARouteNoRouterDefines`\n\n### In\n| Origin | Screen |\n| --- | --- |\n| MainTabs | HomeScreen |\n"
	if v, msg := checkRouteDeclared(content, mapx.Node{}); v != Pass {
		t.Errorf("the gate must not confront the route against the router; got %v (%s)", v, msg)
	}
}

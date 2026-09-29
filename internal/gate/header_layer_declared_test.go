package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func layersCfg() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{"landing-feature": {Kind: "code"}, "api": {Kind: "code"}}}
}

func header(layer string) string {
	return "// @anchors\n//   layer: " + layer + "\n\nexport const x = 1\n"
}

func TestHeaderLayerDeclared_unknownLayer(t *testing.T) {
	t.Run("HDLYD-B01: A layer the Estrutura does not have fails", func(t *testing.T) {})
	n := mapx.Node{ID: "NR01Section.tsx", Kind: mapx.KindCode}
	if v, msg := checkHeaderLayerDeclared(header("landing-feature"), n, "", nil, layersCfg()); v != Pass {
		t.Errorf("a declared layer passes, got %v: %s", v, msg)
	}
	v, msg := checkHeaderLayerDeclared(header("landing-component"), n, "", nil, layersCfg())
	if v != Fail || !strings.Contains(msg, "`landing-component`") || !strings.Contains(msg, "(it has: api, landing-feature)") {
		t.Errorf("an undeclared layer fails naming it and the declared ones, got %v: %s", v, msg)
	}
	body := "// @anchors\n//   updated_at: 2026-09-28\n\n// layer: bogus\n"
	if v, _ := checkHeaderLayerDeclared(body, n, "", nil, layersCfg()); v != Skip {
		t.Errorf("a body line is not the header's layer, got %v", v)
	}
}

func TestHeaderLayerDeclared_case(t *testing.T) {
	t.Run("HDLYD-B02: A layer that differs only in case fails naming both", func(t *testing.T) {})
	v, msg := checkHeaderLayerDeclared(header("Landing-Feature"), mapx.Node{Kind: mapx.KindCode}, "", nil, layersCfg())
	if v != Fail || !strings.Contains(msg, "`Landing-Feature`") || !strings.Contains(msg, "`landing-feature`: the case differs") {
		t.Errorf("a case difference fails naming both, got %v: %s", v, msg)
	}
}

func TestHeaderLayerDeclared_skips(t *testing.T) {
	t.Run("HDLYD-B03: Nothing to confront is skipped", func(t *testing.T) {})
	n := mapx.Node{Kind: mapx.KindCode}
	for _, c := range []struct {
		content string
		cfg     *config.Config
	}{
		{"// @anchors\n//   updated_at: 2026-09-28\n", layersCfg()},
		{header("TODO"), layersCfg()},
		{header("todo-later"), layersCfg()},
		{header("api"), &config.Config{}},
		{header("api"), nil},
	} {
		if v, _ := checkHeaderLayerDeclared(c.content, n, "", nil, c.cfg); v != Skip {
			t.Errorf("%q is skipped, got %v", c.content, v)
		}
	}
}

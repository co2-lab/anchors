package initx

import (
	"strings"
	"testing"
)

func TestRenderHeaderGuideDialect(t *testing.T) {
	ts, _ := PresetByName("node-ts")
	out := RenderHeaderGuide(ts, []string{"auth", "billing"})
	if !strings.Contains(out, "// @anchors") {
		t.Error("stack C-like deveria usar dialeto //")
	}
	if !strings.Contains(out, "auth, billing") {
		t.Error("deveria listar as features reais do projeto")
	}
	py, _ := PresetByName("django")
	if !strings.Contains(RenderHeaderGuide(py, nil), "# @anchors") {
		t.Error("stack Python deveria usar dialeto #")
	}
}

func TestRenderHeaderGuideAlwaysHasEssentials(t *testing.T) {
	// mesmo sem preset (Preset{}), o guia menciona o mínimo: code + o gate
	out := RenderHeaderGuide(Preset{}, nil)
	for _, want := range []string{"code:", "updated_at:", "header-conforme", "anchors guide header"} {
		if !strings.Contains(out, want) {
			t.Errorf("guia de header deveria mencionar %q", want)
		}
	}
}

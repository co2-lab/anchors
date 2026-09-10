package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func planCfg() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{
		"dao":            {Pattern: "packages/backend/models/**/*.ts", Regime: "declarativo"},
		"business-logic": {Pattern: "packages/backend/business-logic/**/*.ts"},
	}}
}

func planNode() mapx.Node { return mapx.Node{ID: "plans/x.md", Kind: mapx.KindPlan} }

func TestPlanSeeds_semeiaEmCamadaDeclarativaReprova(t *testing.T) {
	c := "- `packages/backend/models/metadata.spec.md` — **nasce**.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), "/tmp", nil, planCfg())
	if v != Fail {
		t.Fatalf("esperava Fail, got %v", v)
	}
	if !strings.Contains(msg, "RECONHECIDA") {
		t.Errorf("mensagem deveria explicar o motivo: %q", msg)
	}
}

func TestPlanSeeds_camadaRegidaPassa(t *testing.T) {
	c := "- `packages/backend/business-logic/calc.spec.md` — **nasce**.\n"
	if v, msg := checkPlanSeedsValid(c, planNode(), "/tmp", nil, planCfg()); v != Pass {
		t.Errorf("camada regida deveria passar: %v (%s)", v, msg)
	}
}

func TestPlanSeeds_nomeSoltoEmProsaNaoEhSemeadura(t *testing.T) {
	// "ver SubscriptionScreen.spec.md" é referência, não promessa de criar
	c := "Siga o padrão de `SubscriptionScreen.spec.md` e do `_TEMPLATE_SCREEN.spec.md`.\n"
	if v, msg := checkPlanSeedsValid(c, planNode(), "/tmp", nil, planCfg()); v == Fail {
		t.Errorf("citação em prosa não deveria reprovar: %s", msg)
	}
}

func TestPlanSeeds_naoPlanoPula(t *testing.T) {
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkPlanSeedsValid("`packages/backend/models/x.spec.md`", n, "/tmp", nil, planCfg()); v != Skip {
		t.Errorf("só planos são confrontados, esperava Skip, got %v", v)
	}
}

package gate

import (
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestEvaluateObligations_oneStatusPerDuty(t *testing.T) {
	root := t.TempDir()
	// purge.ts mentions only the fulfilled models
	writePackFile(t, root, "purge.ts", "const t = [PAID_ONE, PAID_TWO]\n")
	files := map[string]string{
		"models/PaidOne.spec.md": "<!-- @anchors\n  carries: pii\n-->\n",
		"models/PaidTwo.spec.md": "<!-- @anchors\n  carries: pii\n-->\n",
		"models/Waived.spec.md":  "<!-- @anchors\n  carries: pii\n  obligation_waived: pii-purgavel — shared data\n-->\n",
		"models/Owed.spec.md":    "<!-- @anchors\n  carries: pii\n  obligation_pending: pii-purgavel — phase 3\n-->\n",
		"models/Zeta.spec.md":    "<!-- @anchors\n  carries: pii\n-->\n",
		"models/Alpha.spec.md":   "<!-- @anchors\n  carries: pii\n-->\n",
		"models/Neutral.spec.md": "<!-- @anchors\n  code: ABCDE\n-->\n",
	}
	g := &mapx.Graph{}
	for id, body := range files {
		writePackFile(t, root, id, body)
		g.Nodes = append(g.Nodes, mapx.Node{ID: id, Kind: mapx.KindSpec})
	}
	g.Nodes = append(g.Nodes, mapx.Node{ID: "models/Gone.spec.md", Kind: mapx.KindSpec}) // not on disk

	obs := append(obligCfg().Obligations, config.Obligation{Name: "untriggered", MustAppearIn: []string{"purge.ts"}})
	got := EvaluateObligations(root, &config.Config{}, g, obs)

	want := []ObligationStatus{
		{
			Name:      "pii-purgavel",
			Targets:   []string{"purge.ts"},
			Subject:   6,
			Fulfilled: 3, // two that appear, one waived with a reason
			Debt:      1,
			Waived:    1,
			Missing:   []string{"models/Alpha.spec.md", "models/Zeta.spec.md"},
		},
		{Name: "untriggered", Targets: []string{"purge.ts"}}, // no `when`: nobody is subject
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("EvaluateObligations:\n got %+v\nwant %+v", got, want)
	}
}

func TestObligationsInForce_isTheResolvedList(t *testing.T) {
	cfg := obligCfg()
	got := ObligationsInForce(t.TempDir(), cfg)
	if !reflect.DeepEqual(got, cfg.Obligations) {
		t.Errorf("got %+v", got)
	}
}

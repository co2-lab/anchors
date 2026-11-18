// @anchors
//   code: DCTSG
//   ref: INDCN

package initx

import (
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func sampleConfig() *config.Config {
	return &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":         {Kind: "spec", Tags: []string{"spec"}},
			"mobile-code":  {Kind: "code", Tags: []string{"frontend", "mobile"}},
			"backend-code": {Kind: "code", Tags: []string{"backend"}},
			"generated":    {Kind: "code", Tags: []string{"generated"}},
			"guide":        {Kind: "guide", Tags: []string{"guide"}},
		},
	}
}

func TestCodeLayerNames(t *testing.T) {
	t.Run("INDCN-B01: The code layer names are listed sorted, and only code layers", func(t *testing.T) {})
	got := CodeLayerNames(sampleConfig())
	want := []string{"backend-code", "generated", "mobile-code"} // sorted, kind=code only
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CodeLayerNames = %v, want %v", got, want)
	}
}

func TestPruneCodeLayers(t *testing.T) {
	t.Run("INDCN-B02: Pruning removes the code layers not kept and never an artifact layer", func(t *testing.T) {})
	cfg := sampleConfig()
	// the user keeps only mobile-code and backend-code (unchecks generated)
	PruneCodeLayers(cfg, map[string]bool{"mobile-code": true, "backend-code": true})

	if _, ok := cfg.Layers["generated"]; ok {
		t.Error("the unchosen code layer 'generated' should have been removed")
	}
	// the chosen code layers remain
	for _, name := range []string{"mobile-code", "backend-code"} {
		if _, ok := cfg.Layers[name]; !ok {
			t.Errorf("the chosen code layer %q was wrongly removed", name)
		}
	}
	// artifact layers are NEVER removed by PruneCodeLayers
	for _, name := range []string{"spec", "guide"} {
		if _, ok := cfg.Layers[name]; !ok {
			t.Errorf("artifact layer %q should not be affected", name)
		}
	}
}

func TestTags(t *testing.T) {
	t.Run("INDCN-B03: The candidate tags are the union of every layer's tags, deduplicated and sorted", func(t *testing.T) {})
	got := Tags(sampleConfig())
	want := []string{"backend", "frontend", "generated", "guide", "mobile", "spec"} // deduplicated, sorted
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tags = %v, want %v", got, want)
	}
}

func TestBuildGovernRules(t *testing.T) {
	t.Run("INDCN-B04: One governs rule per guide answered with a tag, ordered by guide, skipping the unanswered and none", func(t *testing.T) {})
	answers := map[string]string{
		"guides/FRONTEND_GUIDE.md": "frontend",
		"guides/BACKEND_GUIDE.md":  "backend",
		"guides/SENTRY_GUIDE.md":   NoneTag, // skipped
		"guides/ASSETS_GUIDE.md":   "",      // skipped (no answer)
	}
	got := BuildGovernRules(answers)

	// only 2 rules (the skipped ones are left out), ordered by guide
	want := []config.GovernRule{
		{From: "guides/BACKEND_GUIDE.md", Governs: "backend"},
		{From: "guides/FRONTEND_GUIDE.md", Governs: "frontend"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildGovernRules = %+v, want %+v", got, want)
	}
}

func TestBuildGovernRules_empty(t *testing.T) {
	t.Run("INDCN-B05: No answer gives no governs rule", func(t *testing.T) {})
	if got := BuildGovernRules(map[string]string{}); len(got) != 0 {
		t.Fatalf("expected no rule, got %+v", got)
	}
}

func TestPruneThenListKeepsExactlyTheKept(t *testing.T) {
	t.Run("INDCN-I01: After pruning, the code layer names are exactly the kept code layers", func(t *testing.T) {})
	cfg := PruneCodeLayers(sampleConfig(), map[string]bool{"mobile-code": true, "spec": true, "ghost": true})
	if got := CodeLayerNames(cfg); !reflect.DeepEqual(got, []string{"mobile-code"}) {
		t.Errorf("CodeLayerNames after pruning = %v, want [mobile-code]", got)
	}
}

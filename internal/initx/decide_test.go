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
	got := CodeLayerNames(sampleConfig())
	want := []string{"backend-code", "generated", "mobile-code"} // ordenado, só kind=code
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CodeLayerNames = %v, want %v", got, want)
	}
}

func TestPruneCodeLayers(t *testing.T) {
	cfg := sampleConfig()
	// usuário mantém só mobile-code e backend-code (desmarca generated)
	PruneCodeLayers(cfg, map[string]bool{"mobile-code": true, "backend-code": true})

	if _, ok := cfg.Layers["generated"]; ok {
		t.Error("layer de código não escolhida 'generated' deveria ter sido removida")
	}
	// camadas de código escolhidas permanecem
	for _, name := range []string{"mobile-code", "backend-code"} {
		if _, ok := cfg.Layers[name]; !ok {
			t.Errorf("layer de código escolhida %q foi removida indevidamente", name)
		}
	}
	// camadas de artefato NUNCA são removidas por PruneCodeLayers
	for _, name := range []string{"spec", "guide"} {
		if _, ok := cfg.Layers[name]; !ok {
			t.Errorf("layer de artefato %q não deveria ser afetada", name)
		}
	}
}

func TestTags(t *testing.T) {
	got := Tags(sampleConfig())
	want := []string{"backend", "frontend", "generated", "guide", "mobile", "spec"} // dedup, ordenado
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tags = %v, want %v", got, want)
	}
}

func TestBuildGovernRules(t *testing.T) {
	answers := map[string]string{
		"guides/FRONTEND_GUIDE.md": "frontend",
		"guides/BACKEND_GUIDE.md":  "backend",
		"guides/SENTRY_GUIDE.md":   NoneTag, // pulado
		"guides/ASSETS_GUIDE.md":   "",      // pulado (sem resposta)
	}
	got := BuildGovernRules(answers)

	// só 2 regras (as pulas ficam de fora), ordenadas por guide
	want := []config.GovernRule{
		{From: "guides/BACKEND_GUIDE.md", Governs: "backend"},
		{From: "guides/FRONTEND_GUIDE.md", Governs: "frontend"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildGovernRules = %+v, want %+v", got, want)
	}
}

func TestBuildGovernRules_empty(t *testing.T) {
	if got := BuildGovernRules(map[string]string{}); len(got) != 0 {
		t.Fatalf("esperava nenhuma regra, veio %+v", got)
	}
}

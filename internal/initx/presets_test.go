package initx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestCatalogIsWellFormed(t *testing.T) {
	if len(Presets) == 0 {
		t.Fatal("catálogo de presets vazio")
	}
	seen := map[string]bool{}
	for _, p := range Presets {
		if p.Name == "" || p.Title == "" {
			t.Errorf("preset sem name/title: %+v", p)
		}
		if seen[p.Name] {
			t.Errorf("nome de preset duplicado: %s", p.Name)
		}
		seen[p.Name] = true
		if len(p.Layers) == 0 {
			t.Errorf("preset %s sem layers", p.Name)
		}
		hasTest := false
		for _, l := range p.Layers {
			if l.Pattern == "" {
				t.Errorf("preset %s: layer %s sem pattern", p.Name, l.Name)
			}
			if l.Kind == "test" {
				hasTest = true
			}
		}
		if !hasTest {
			t.Errorf("preset %s: nenhuma layer de teste", p.Name)
		}
		if p.Modular && p.ModuleGlob == "" {
			t.Errorf("preset %s modular sem ModuleGlob", p.Name)
		}
	}
}

func TestToLayersDefaultsKindToCode(t *testing.T) {
	p := Preset{Layers: []PresetLayer{
		{Name: "a", Pattern: "a/**", Kind: ""}, // sem kind → code
		{Name: "t", Pattern: "t/**", Kind: "test"},
	}}
	layers := p.ToLayers()
	if layers["a"].Kind != "code" {
		t.Errorf("kind vazio deveria virar code, veio %q", layers["a"].Kind)
	}
	if layers["t"].Kind != "test" {
		t.Errorf("kind test deveria ser preservado, veio %q", layers["t"].Kind)
	}
}

func TestApplyPresetWritesLayers(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"guide": {Kind: "guide"}, // layer preexistente não deve sumir
	}}
	preset, _ := PresetByName("go")
	ApplyPreset(cfg, preset, nil)
	for _, want := range []string{"cmd", "internal", "pkg", "test"} {
		if _, ok := cfg.Layers[want]; !ok {
			t.Errorf("layer %q do preset go deveria existir", want)
		}
	}
	if _, ok := cfg.Layers["guide"]; !ok {
		t.Error("layer preexistente guide não deveria ser removida")
	}
}

func TestDeduceModulePrefixesUniqueness(t *testing.T) {
	// auth e audit gerariam ambos "AU" — o segundo deve ser ajustado
	mods := []string{"src/features/auth", "src/features/audit", "src/features/family"}
	pfx := DeduceModulePrefixes(mods)
	if len(pfx) != 3 {
		t.Fatalf("esperava 3 prefixos, veio %d", len(pfx))
	}
	seen := map[string]bool{}
	for m, p := range pfx {
		if len(p) != 2 {
			t.Errorf("prefixo de %s tem %d chars (esperava 2)", m, len(p))
		}
		if seen[p] {
			t.Errorf("prefixo duplicado %q (colisão não resolvida)", p)
		}
		seen[p] = true
	}
	// family → FM (determinístico)
	if pfx["family"] != "FM" {
		t.Errorf("family → %q, esperava FM", pfx["family"])
	}
}

func TestDeduceModulePrefixesDeterministic(t *testing.T) {
	mods := []string{"a/x", "a/y", "a/z"}
	if got1, got2 := DeduceModulePrefixes(mods), DeduceModulePrefixes(mods); len(got1) != len(got2) {
		t.Fatal("não determinístico no tamanho")
	} else {
		for k, v := range got1 {
			if got2[k] != v {
				t.Errorf("não determinístico: %s = %q vs %q", k, v, got2[k])
			}
		}
	}
}

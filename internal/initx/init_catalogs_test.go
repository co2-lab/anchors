package initx

import (
	"strings"
	"testing"
)

func TestCatalogIsWellFormed(t *testing.T) {
	t.Run("INCTN-I01: Every preset has a unique name, a title, patterned layers and a test layer", func(t *testing.T) {})
	t.Run("INCTN-B05: A modular preset declares the directory of its modules", func(t *testing.T) {})
	if len(Presets) == 0 {
		t.Fatal("empty preset catalog")
	}
	seen := map[string]bool{}
	for _, p := range Presets {
		if p.Name == "" || p.Title == "" {
			t.Errorf("preset with no name/title: %+v", p)
		}
		if seen[p.Name] {
			t.Errorf("duplicate preset name: %s", p.Name)
		}
		seen[p.Name] = true
		if len(p.Layers) == 0 {
			t.Errorf("preset %s has no layers", p.Name)
		}
		hasTest := false
		for _, l := range p.Layers {
			if l.Pattern == "" {
				t.Errorf("preset %s: layer %s has no pattern", p.Name, l.Name)
			}
			if l.Kind == "test" {
				hasTest = true
			}
		}
		if !hasTest {
			t.Errorf("preset %s: no test layer", p.Name)
		}
		if p.Modular && p.ModuleGlob == "" {
			t.Errorf("modular preset %s has no ModuleGlob", p.Name)
		}
	}
}

func TestCatalogLeavesCodePrefixEmpty(t *testing.T) {
	t.Run("INCTN-X01: No preset layer carries an identity prefix", func(t *testing.T) {})
	for _, p := range Presets {
		for _, l := range p.Layers {
			if l.CodePrefix != "" {
				t.Errorf("preset %s layer %s carries code prefix %q; it is deduced per module at init",
					p.Name, l.Name, l.CodePrefix)
			}
		}
	}
}

func TestToLayersDefaultsKindToCode(t *testing.T) {
	t.Run("INCTN-B01: A preset layer with no kind becomes a code layer", func(t *testing.T) {})
	p := Preset{Layers: []PresetLayer{
		{Name: "a", Pattern: "a/**", Kind: ""}, // no kind → code
		{Name: "t", Pattern: "t/**", Kind: "test", Tags: []string{"x"}},
	}}
	layers := p.ToLayers()
	if layers["a"].Kind != "code" {
		t.Errorf("an empty kind should become code, got %q", layers["a"].Kind)
	}
	if layers["t"].Kind != "test" || layers["t"].Pattern != "t/**" || len(layers["t"].Tags) != 1 {
		t.Errorf("a declared kind, pattern and tags should be preserved, got %+v", layers["t"])
	}
}

func TestPresetLookup(t *testing.T) {
	t.Run("INCTN-B02: Looking up a preset by an unknown name finds nothing", func(t *testing.T) {})
	t.Run("INCTN-B03: The preset names are listed in catalog order", func(t *testing.T) {})
	if p, ok := PresetByName("go"); !ok || p.Name != "go" {
		t.Errorf("the go preset should be found, got %+v %v", p, ok)
	}
	if p, ok := PresetByName("no-such-stack"); ok || p.Name != "" || len(p.Layers) != 0 {
		t.Errorf("an unknown name should give (zero, false), got %+v %v", p, ok)
	}
	names := PresetNames()
	if len(names) != len(Presets) {
		t.Fatalf("PresetNames has %d names for %d presets", len(names), len(Presets))
	}
	for i, p := range Presets {
		if names[i] != p.Name {
			t.Errorf("name %d is %q, catalog order says %q", i, names[i], p.Name)
		}
	}
}

// The instruction SAYS WHAT TO ANSWER, and says NOT to pass.
//
// A text that only mentioned `@TBD` would satisfy the judgment-gate test without solving
// the problem: the whole point is that the answer is not `pass`.
func TestTBDInstructionForbidsPassAndNamesTheAbsence(t *testing.T) {
	t.Run("INCTN-B04: The @TBD instruction forbids pass, orders a waiver naming the absence, and names the piece asked about", func(t *testing.T) {})
	t.Run("INCTN-I02: The @TBD instruction demands checking that the @TBD is still true", func(t *testing.T) {})
	got := tbdInstruction("o código")
	for _, required := range []string{"@TBD", "DISPENSADO", "pass", "nomeando a ausência"} {
		if !strings.Contains(got, required) {
			t.Errorf("the instruction does not mention %q:\n%s", required, got)
		}
	}
	// The piece goes into the text: without it the instruction would speak of "the code"
	// in a gate that asks about a test.
	if !strings.Contains(tbdInstruction("o teste"), "o teste") {
		t.Error("the instruction does not use the piece it received")
	}
	// The stale `@TBD` is the other half: a piece that came to exist with the marker still
	// in the file makes every gate that reads it waive what it should charge.
	if !strings.Contains(got, "desatualizado") {
		t.Errorf("the instruction does not cover the `@TBD` that stopped being true:\n%s", got)
	}
}

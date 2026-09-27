package initx

import (
	"fmt"
	"testing"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
)

func TestApplyPresetWritesLayers(t *testing.T) {
	t.Run("APPRP-B01: Applying a preset adds its layers and keeps the layers of other names", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"guide": {Kind: "guide"}, // a pre-existing layer must not disappear
	}}
	preset, _ := PresetByName("go")
	ApplyPreset(cfg, preset, nil)
	for _, want := range []string{"cmd", "internal", "pkg", "test"} {
		if _, ok := cfg.Layers[want]; !ok {
			t.Errorf("layer %q of the go preset should exist", want)
		}
	}
	if _, ok := cfg.Layers["guide"]; !ok {
		t.Error("the pre-existing guide layer should not be removed")
	}

	// A configuration with no layer set at all gets one.
	empty := &config.Config{}
	ApplyPreset(empty, preset, nil)
	if len(empty.Layers) != len(preset.Layers) {
		t.Errorf("a configuration with no layers should receive the %d preset layers, got %d",
			len(preset.Layers), len(empty.Layers))
	}
}

func TestApplyPresetReplacesSameNamedLayer(t *testing.T) {
	t.Run("APPRP-B02: A layer with the same name as a preset layer is replaced by the preset's", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"internal": {Kind: "doc", Pattern: "old/**"},
	}}
	preset, _ := PresetByName("go")
	ApplyPreset(cfg, preset, nil)
	if got := cfg.Layers["internal"]; got.Kind != "code" || got.Pattern == "old/**" {
		t.Errorf("the same-named layer should carry the preset's declaration, got %+v", got)
	}
}

func TestApplyPresetReturnsModulePrefixes(t *testing.T) {
	t.Run("APPRP-B05: Applying a preset returns the prefix deduced for each detected module", func(t *testing.T) {})
	cfg := &config.Config{}
	preset, _ := PresetByName("node-ts")
	got := ApplyPreset(cfg, preset, []string{"src/modules/family"})
	if len(got) != 1 || got["family"] != "FM" {
		t.Errorf("ApplyPreset should return {family: FM}, got %v", got)
	}
	if got := ApplyPreset(cfg, preset, nil); len(got) != 0 {
		t.Errorf("no detected module gives no prefix, got %v", got)
	}
}

func TestDeduceModulePrefixesUniqueness(t *testing.T) {
	t.Run("APPRP-B03: Each module receives a two-letter prefix keyed by its directory name", func(t *testing.T) {})
	mods := []string{"src/features/auth", "src/features/audit", "src/features/family/"}
	pfx := DeduceModulePrefixes(mods)
	if len(pfx) != 3 {
		t.Fatalf("expected 3 prefixes, got %d", len(pfx))
	}
	for m, p := range pfx {
		if len(p) != 2 {
			t.Errorf("prefix of %s has %d chars (expected 2)", m, len(p))
		}
		if p != code.ModulePrefix(m) {
			t.Errorf("with no collision, %s should take the identity prefix %q, got %q", m, code.ModulePrefix(m), p)
		}
	}
	// family → FM, keyed by the directory name even with a trailing slash
	if pfx["family"] != "FM" {
		t.Errorf("family → %q, expected FM", pfx["family"])
	}
}

func TestDeduceModulePrefixesResolvesCollision(t *testing.T) {
	t.Run("APPRP-B04: A colliding prefix keeps its first letter and takes the first free second letter", func(t *testing.T) {})
	if code.ModulePrefix("atlas") != code.ModulePrefix("auth") {
		t.Fatalf("precondition: atlas and auth must share a prefix (%s vs %s)",
			code.ModulePrefix("atlas"), code.ModulePrefix("auth"))
	}
	pfx := DeduceModulePrefixes([]string{"m/auth", "m/atlas"})
	// atlas comes first alphabetically and keeps AT; auth moves to the first free AA..AZ.
	if pfx["atlas"] != "AT" || pfx["auth"] != "AA" {
		t.Errorf("expected atlas=AT auth=AA, got %v", pfx)
	}
}

func TestDeduceModulePrefixesDeterministic(t *testing.T) {
	t.Run("APPRP-I01: The same modules in any order give the same prefixes", func(t *testing.T) {})
	a := DeduceModulePrefixes([]string{"m/auth", "m/atlas", "m/family"})
	b := DeduceModulePrefixes([]string{"m/family", "m/atlas", "m/auth"})
	if len(a) != len(b) {
		t.Fatal("not deterministic in size")
	}
	for k, v := range a {
		if b[k] != v {
			t.Errorf("not deterministic: %s = %q vs %q", k, v, b[k])
		}
	}
}

func TestDeduceModulePrefixesReadsNoDisk(t *testing.T) {
	t.Run("APPRP-X01: Prefixes are deduced from the given paths without reading the disk", func(t *testing.T) {})
	got := DeduceModulePrefixes([]string{"/does/not/exist/family"})
	if got["family"] != "FM" {
		t.Errorf("a module that does not exist on disk still gets its prefix, got %v", got)
	}
}

// Two modules with the same folder name are two modules. Keyed by the basename alone,
// the second overwrote the first and one module vanished from the mapping.
func TestDeduceModulePrefixesKeepsModulesWithTheSameName(t *testing.T) {
	t.Run("APPRP-B06: Modules sharing a folder name are each keyed by their path, with distinct prefixes", func(t *testing.T) {})
	pfx := DeduceModulePrefixes([]string{"packages/auth", "apps/auth/", "m/family"})
	if len(pfx) != 3 {
		t.Fatalf("every module must be in the mapping, got %v", pfx)
	}
	a, b := pfx["apps/auth"], pfx["packages/auth"]
	if a == "" || b == "" || a == b {
		t.Errorf("the two auth modules need distinct prefixes under their paths, got %v", pfx)
	}
	if pfx["family"] != "FM" {
		t.Errorf("a module with a unique name stays keyed by its name, got %v", pfx)
	}
}

// After the 26 second letters of an initial are taken, the next module used to keep the
// prefix ALREADY taken — two modules with one identity, and nothing said so.
func TestDeduceModulePrefixesNeverReusesAPrefix(t *testing.T) {
	t.Run("APPRP-I02: No two modules ever share a prefix while a free one exists", func(t *testing.T) {})
	var mods []string
	for i := range 30 {
		mods = append(mods, fmt.Sprintf("m/a%02d", i)) // every one starts with A
	}
	pfx := DeduceModulePrefixes(mods)
	seen := map[string]string{}
	for m, p := range pfx {
		if len(p) != 2 {
			t.Errorf("%s has prefix %q", m, p)
		}
		if other, dup := seen[p]; dup {
			t.Errorf("%s and %s share the prefix %s", m, other, p)
		}
		seen[p] = m
	}
	if len(pfx) != 30 {
		t.Errorf("expected 30 modules, got %d", len(pfx))
	}
}

// The search order of a free prefix: the same first letter with the second from A to Z,
// then every other first letter from A, each with its second letter from A to Z. Each case
// pins one end of that order: the second letter after A, the last second letter Z, the
// first letter after the colliding one, and the last first letter Z.
func TestFirstFreePrefixSearchOrder(t *testing.T) {
	t.Run("APPRP-B04: A colliding prefix keeps its first letter and takes the first free second letter", func(t *testing.T) {})
	takeAll := func(taken map[string]bool, first byte) {
		for c := byte('A'); c <= 'Z'; c++ {
			taken[string([]byte{first, c})] = true
		}
	}
	cases := []struct {
		name  string
		first byte
		taken func() map[string]bool
		want  string
	}{
		{"second letter moves forward from A", 'A', func() map[string]bool {
			return map[string]bool{"AA": true}
		}, "AB"},
		{"second letter reaches Z", 'A', func() map[string]bool {
			m := map[string]bool{}
			takeAll(m, 'A')
			delete(m, "AZ")
			return m
		}, "AZ"},
		{"other first letters are tried forward from A", 'B', func() map[string]bool {
			m := map[string]bool{}
			takeAll(m, 'A')
			takeAll(m, 'B')
			return m
		}, "CA"},
		{"first letter reaches Z", 'A', func() map[string]bool {
			m := map[string]bool{}
			for f := byte('A'); f <= 'Z'; f++ {
				takeAll(m, f)
			}
			delete(m, "ZZ")
			return m
		}, "ZZ"},
	}
	for _, tc := range cases {
		if got := firstFreePrefix(tc.first, tc.taken()); got != tc.want {
			t.Errorf("%s: firstFreePrefix(%q) = %q, want %q", tc.name, tc.first, got, tc.want)
		}
	}
}

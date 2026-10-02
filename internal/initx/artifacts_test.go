// @anchors
//   ref: ARCHR

package initx

import (
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestApplyArtifactChoice_addsAndRemoves(t *testing.T) {
	t.Run("ARCHR-B04: An artifact layer that was not chosen is removed, and code layers are untouched", func(t *testing.T) {})
	t.Run("ARCHR-B05: A chosen artifact layer that already exists is kept as declared", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"spec":        {Kind: "spec", Tags: []string{"spec"}, Pattern: "docs/**/*.spec.md"}, // already exists
		"test":        {Kind: "test", Pattern: "**/*_test.go"},
		"mobile-code": {Kind: "code", Tags: []string{"mobile"}},
	}}
	// the user wants spec + feature (but NOT test/guide)
	ApplyArtifactChoice(cfg, map[string]bool{"spec": true, "feature": true}, nil, "")

	if _, ok := cfg.Layers["feature"]; !ok {
		t.Error("the chosen feature should have been added")
	}
	if got := cfg.Layers["spec"].Pattern; got != "docs/**/*.spec.md" {
		t.Errorf("the chosen spec should remain as declared, pattern = %q", got)
	}
	if _, ok := cfg.Layers["test"]; ok {
		t.Error("the unchosen test should not exist")
	}
	// a code layer is not touched by ApplyArtifactChoice
	if _, ok := cfg.Layers["mobile-code"]; !ok {
		t.Error("the code layer should not be affected")
	}
}

func TestApplyArtifactChoice_emptyProject(t *testing.T) {
	t.Run("ARCHR-B03: A chosen artifact layer is created even when nothing was detected", func(t *testing.T) {})
	// an empty project: nothing detected, but the user chooses spec+feature+test.
	cfg := &config.Config{}
	ApplyArtifactChoice(cfg, map[string]bool{"spec": true, "feature": true, "test": true}, nil, "")

	for _, name := range []string{"spec", "feature", "test"} {
		l, ok := cfg.Layers[name]
		if !ok {
			t.Errorf("in an empty project, the chosen layer %q should be created", name)
			continue
		}
		if l.Kind != name || len(l.Tags) != 1 || l.Tags[0] != name || l.Pattern == "" {
			t.Errorf("layer %q should carry its kind, its default pattern and its name as tag, got %+v", name, l)
		}
	}
}

func TestApplyArtifactChoice_guideDir(t *testing.T) {
	t.Run("ARCHR-B06: The guide and plan layers take the detected directory, otherwise the default pattern", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"guide": true}, map[string]string{"guide": "docs/guides"}, "")
	if got := cfg.Layers["guide"].Pattern; got != "docs/guides/*.md" {
		t.Errorf("guide pattern = %q, want docs/guides/*.md", got)
	}
}

func TestApplyArtifactChoice_planLayer(t *testing.T) {
	// plan chosen, no detected dir → default pattern plans/*.md (specific, to beat the
	// doc wildcard so the watcher suggests 'specify' instead of 'triage').
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"plan": true}, nil, "")
	if got := cfg.Layers["plan"].Pattern; got != "plans/*.md" {
		t.Errorf("plan pattern = %q, want plans/*.md", got)
	}
	if got := cfg.Layers["plan"].Kind; got != "plan" {
		t.Errorf("plan kind = %q, want plan", got)
	}
	// with a detected dir → the dir is respected
	cfg2 := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg2, map[string]bool{"plan": true}, map[string]string{"plan": "docs/plans"}, "")
	if got := cfg2.Layers["plan"].Pattern; got != "docs/plans/*.md" {
		t.Errorf("plan pattern with dir = %q, want docs/plans/*.md", got)
	}
}

func TestApplyColocation(t *testing.T) {
	t.Run("ARCHR-B07: Colocation is declared from the spec as anchor, and the spec is never a derivative", func(t *testing.T) {})
	t.Run("ARCHR-B08: No colocation is declared when it is not wanted, when spec is not chosen, or when nothing derives from the spec", func(t *testing.T) {})
	cfg := &config.Config{}
	// The ANCHOR is the spec, and it is not among the derivatives: from the spec the code,
	// the feature and the test are born. With the four chosen, THREE derivatives remain.
	ApplyColocation(cfg, true, map[string]bool{
		"spec": true, "code": true, "feature": true, "test": true,
	}, "")
	if cfg.Derived == nil || len(cfg.Derived.Files) != 3 {
		t.Fatalf("expected derived with 3 files, got %+v", cfg.Derived)
	}
	if cfg.Derived.Anchor != "spec" {
		t.Errorf("the anchor is the spec, got %q", cfg.Derived.Anchor)
	}
	if _, ehDerivada := cfg.Derived.Files["spec"]; ehDerivada {
		t.Error("the spec is the ORIGIN — listing it as a derivative inverts the doctrine")
	}
	if got := cfg.Derived.Files["feature"]; len(got) != 1 || got[0] != "{{dir}}/{{name}}.feature" {
		t.Errorf("the feature derivative sits beside the spec, got %v", got)
	}

	// With no spec chosen there is no anchor, and colocation is not declared: a project
	// with no spec has no unit to locate.
	semSpec := &config.Config{}
	ApplyColocation(semSpec, true, map[string]bool{"code": true, "test": true}, "")
	if semSpec.Derived != nil {
		t.Error("with no spec there is no anchor — there should be no derived")
	}
	// turning colocation off → derived disappears
	ApplyColocation(cfg, false, map[string]bool{"spec": true}, "")
	if cfg.Derived != nil {
		t.Error("colocation turned off should clear derived")
	}
	// colocation on but nothing derives from the spec (spec + guide) → no derived
	ApplyColocation(cfg, true, map[string]bool{"spec": true, "guide": true}, "")
	if cfg.Derived != nil {
		t.Error("with nothing derived from the spec, derived should stay nil")
	}
}

func TestArtifactNamesInStableOrder(t *testing.T) {
	t.Run("ARCHR-B01: The artifact options are spec, feature, test, guide, plan and code, in that order", func(t *testing.T) {})
	want := []string{"spec", "feature", "test", "guide", "plan", "code"}
	for i := 0; i < 3; i++ {
		if got := ArtifactNames(); !reflect.DeepEqual(got, want) {
			t.Fatalf("ArtifactNames = %v, want %v", got, want)
		}
	}
}

func TestDetectedArtifactsPreChecksWhatInferenceFound(t *testing.T) {
	t.Run("ARCHR-B02: Each artifact inference found is pre-checked, and nothing else", func(t *testing.T) {})
	p := &Proposal{HasSpecMD: true, HasTest: true, PlanDir: "plans", CodeDirs: []string{"src/app"}}
	got := p.DetectedArtifacts()
	want := map[string]bool{"spec": true, "test": true, "plan": true, "code": true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DetectedArtifacts = %v, want %v", got, want)
	}
	p = &Proposal{HasFeature: true, GuideDir: "guides"}
	if got := p.DetectedArtifacts(); !reflect.DeepEqual(got, map[string]bool{"feature": true, "guide": true}) {
		t.Errorf("DetectedArtifacts = %v, want feature and guide", got)
	}
	if got := (&Proposal{}).DetectedArtifacts(); len(got) != 0 {
		t.Errorf("nothing detected pre-checks nothing, got %v", got)
	}
}

// `code` is an artifact CHOICE, not a layer: the code layers come from inference, one per
// folder that holds code. Choosing it must create nothing, and leaving it out must
// remove nothing — PruneCodeLayers decides which code layers stay.
func TestApplyArtifactChoice_codeIsAChoiceWithoutALayer(t *testing.T) {
	t.Run("ARCHR-B09: Choosing code or leaving it out creates and removes no layer", func(t *testing.T) {})
	for _, chosen := range []map[string]bool{{"code": true}, {}} {
		cfg := &config.Config{Layers: map[string]config.Layer{
			"code":     {Kind: "code", Pattern: "src/**/*.go"},
			"app-code": {Kind: "code", Pattern: "app/**/*.go"},
		}}
		ApplyArtifactChoice(cfg, chosen, nil, "")
		if len(cfg.Layers) != 2 || cfg.Layers["code"].Pattern != "src/**/*.go" {
			t.Errorf("choice %v: the code layers must stay as they were, got %v", chosen, cfg.Layers)
		}
	}
}

// The chain the option exists for: the TUI and `--artifacts` offer ArtifactNames, and
// colocation declares the code beside the spec only when `code` was chosen.
func TestApplyColocation_codeFromTheOfferedChoice(t *testing.T) {
	t.Run("ARCHR-B07: Colocation is declared from the spec as anchor, and the spec is never a derivative", func(t *testing.T) {})
	chosen := map[string]bool{}
	for _, n := range ArtifactNames() {
		chosen[n] = true
	}
	cfg := &config.Config{}
	ApplyColocation(cfg, true, chosen, "")
	if cfg.Derived == nil {
		t.Fatal("colocation with every offered artifact must be declared")
	}
	if got := cfg.Derived.Files["code"]; len(got) != 1 || got[0] != "{{dir}}/{{name}}.{{ext}}" {
		t.Errorf("choosing every offered artifact must put the code beside the spec, got %v", cfg.Derived.Files)
	}
}

// The colocation owns only its part of `derived`: the test handle the inference found
// survives colocation on and off.
func TestApplyColocation_keepsTheInferredTestHandle(t *testing.T) {
	t.Run("ARCHR-B10: Colocation keeps the rest of derived", func(t *testing.T) {})
	all := map[string]bool{"spec": true, "feature": true, "test": true}
	on := &config.Config{Derived: &config.Derived{TestHandle: "testID"}}
	ApplyColocation(on, true, all, "")
	if on.Derived == nil || on.Derived.TestHandle != "testID" || on.Derived.Anchor != "spec" {
		t.Errorf("colocation on must keep the test handle: %+v", on.Derived)
	}
	off := &config.Config{Derived: &config.Derived{TestHandle: "testID"}}
	ApplyColocation(off, false, all, "")
	if off.Derived == nil || off.Derived.TestHandle != "testID" || len(off.Derived.Files) != 0 {
		t.Errorf("colocation off must keep only the test handle: %+v", off.Derived)
	}
	none := &config.Config{}
	ApplyColocation(none, false, all, "")
	if none.Derived != nil {
		t.Errorf("with nothing to keep, derived stays empty: %+v", none.Derived)
	}
}

func TestApplyArtifactChoice_testLayerTakesTheProjectPattern(t *testing.T) {
	t.Run("ARCHR-B11: A created test layer takes the project's test pattern", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"test": true}, nil, "**/*_test.go")
	if got := cfg.Layers["test"].Pattern; got != "**/*_test.go" {
		t.Errorf("pattern = %q, want **/*_test.go", got)
	}
	cfg = &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"test": true}, nil, "")
	if got := cfg.Layers["test"].Pattern; got != "**/*.test.*" {
		t.Errorf("pattern = %q, want the default **/*.test.*", got)
	}
}

func TestApplyColocation_testTemplateIsTheProjects(t *testing.T) {
	t.Run("ARCHR-B12: The colocated test template is the project's", func(t *testing.T) {})
	for tmpl, want := range map[string]string{
		"{{dir}}/{{name}}_test.go": "{{dir}}/{{name}}_test.go",
		"":                         "{{dir}}/{{name}}.test.{{ext}}",
	} {
		cfg := &config.Config{}
		ApplyColocation(cfg, true, map[string]bool{"spec": true, "test": true}, tmpl)
		if got := cfg.Derived.Files["test"]; len(got) != 1 || got[0] != want {
			t.Errorf("template %q: test derivative = %v, want %s", tmpl, got, want)
		}
	}
}

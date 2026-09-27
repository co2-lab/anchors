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
	ApplyArtifactChoice(cfg, map[string]bool{"spec": true, "feature": true}, nil)

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
	ApplyArtifactChoice(cfg, map[string]bool{"spec": true, "feature": true, "test": true}, nil)

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
	ApplyArtifactChoice(cfg, map[string]bool{"guide": true}, map[string]string{"guide": "docs/guides"})
	if got := cfg.Layers["guide"].Pattern; got != "docs/guides/*.md" {
		t.Errorf("guide pattern = %q, want docs/guides/*.md", got)
	}
}

func TestApplyArtifactChoice_planLayer(t *testing.T) {
	// plan chosen, no detected dir → default pattern plans/*.md (specific, to beat the
	// doc wildcard so the watcher suggests 'specify' instead of 'triage').
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"plan": true}, nil)
	if got := cfg.Layers["plan"].Pattern; got != "plans/*.md" {
		t.Errorf("plan pattern = %q, want plans/*.md", got)
	}
	if got := cfg.Layers["plan"].Kind; got != "plan" {
		t.Errorf("plan kind = %q, want plan", got)
	}
	// with a detected dir → the dir is respected
	cfg2 := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg2, map[string]bool{"plan": true}, map[string]string{"plan": "docs/plans"})
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
	})
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
	// with no spec has no triad to locate.
	semSpec := &config.Config{}
	ApplyColocation(semSpec, true, map[string]bool{"code": true, "test": true})
	if semSpec.Derived != nil {
		t.Error("with no spec there is no anchor — there should be no derived")
	}
	// turning colocation off → derived disappears
	ApplyColocation(cfg, false, map[string]bool{"spec": true})
	if cfg.Derived != nil {
		t.Error("colocation turned off should clear derived")
	}
	// colocation on but nothing derives from the spec (spec + guide) → no derived
	ApplyColocation(cfg, true, map[string]bool{"spec": true, "guide": true})
	if cfg.Derived != nil {
		t.Error("with nothing derived from the spec, derived should stay nil")
	}
}

func TestArtifactNamesInStableOrder(t *testing.T) {
	t.Run("ARCHR-B01: The artifact options are spec, feature, test, guide and plan, in that order", func(t *testing.T) {})
	want := []string{"spec", "feature", "test", "guide", "plan"}
	for i := 0; i < 3; i++ {
		if got := ArtifactNames(); !reflect.DeepEqual(got, want) {
			t.Fatalf("ArtifactNames = %v, want %v", got, want)
		}
	}
}

func TestDetectedArtifactsPreChecksWhatInferenceFound(t *testing.T) {
	t.Run("ARCHR-B02: Each artifact inference found is pre-checked, and nothing else", func(t *testing.T) {})
	p := &Proposal{HasSpecMD: true, HasTest: true, PlanDir: "plans"}
	got := p.DetectedArtifacts()
	want := map[string]bool{"spec": true, "test": true, "plan": true}
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

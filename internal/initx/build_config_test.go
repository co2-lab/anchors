package initx

import (
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestBuildConfigCodeLayerPerDir(t *testing.T) {
	t.Run("BLCNB-B01: Each detected code directory becomes a code layer named after its last segment", func(t *testing.T) {})
	t.Run("BLCNB-B03: A proposed code layer excludes specs, features and test files", func(t *testing.T) {})
	p := &Proposal{CodeDirs: []string{"apps/mobile", "packages/backend"}, CodeExts: []string{".ts"}}
	c := p.buildConfig()
	if len(c.Layers) != 2 {
		t.Fatalf("expected one layer per code dir, got %v", c.Layers)
	}
	for name, dir := range map[string]string{"mobile-code": "apps/mobile", "backend-code": "packages/backend"} {
		l, ok := c.Layers[name]
		if !ok {
			t.Errorf("layer %q missing", name)
			continue
		}
		if l.Kind != "code" || !reflect.DeepEqual(l.Tags, []string{name}) {
			t.Errorf("layer %q should be a code layer tagged with its name, got %+v", name, l)
		}
		if l.Pattern != dir+"/**/*.ts" {
			t.Errorf("layer %q pattern = %q", name, l.Pattern)
		}
		if !reflect.DeepEqual(l.Exclude, []string{"**/*.spec.md", "**/*.feature", "**/*.test.*"}) {
			t.Errorf("layer %q exclude = %v", name, l.Exclude)
		}
	}
}

func TestBuildConfigPatternOfSeveralExts(t *testing.T) {
	t.Run("BLCNB-B02: With several detected extensions the pattern lists them as a set", func(t *testing.T) {})
	p := &Proposal{CodeDirs: []string{"a/b"}, CodeExts: []string{".ts", ".tsx"}}
	if got := p.buildConfig().Layers["b-code"].Pattern; got != "a/b/**/*.{ts,tsx}" {
		t.Errorf("pattern = %q, want a/b/**/*.{ts,tsx}", got)
	}
}

func TestBuildConfigColocation(t *testing.T) {
	t.Run("BLCNB-B04: Colocation is proposed only when detected, with templates only for the detected kinds", func(t *testing.T) {})
	p := &Proposal{Colocated: true, HasFeature: true}
	c := p.buildConfig()
	if c.Derived == nil || c.Derived.Anchor != "spec" {
		t.Fatalf("detected colocation should propose derived anchored on the spec, got %+v", c.Derived)
	}
	if !reflect.DeepEqual(c.Derived.Files, map[string]config.Padroes{"feature": {"{{dir}}/{{name}}.feature"}}) {
		t.Errorf("only the detected feature should get a template, got %v", c.Derived.Files)
	}
	p = &Proposal{Colocated: false, HasFeature: true, HasTest: true}
	if c := p.buildConfig(); c.Derived != nil {
		t.Errorf("no colocation detected proposes no derived, got %+v", c.Derived)
	}
}

func TestBuildConfigTestHandle(t *testing.T) {
	t.Run("BLCNB-B05: The test handle is proposed only when inference found one", func(t *testing.T) {})
	p := &Proposal{TestHandle: "testID"}
	c := p.buildConfig()
	if c.Derived == nil || c.Derived.TestHandle != "testID" || len(c.Derived.Files) != 0 {
		t.Errorf("a found handle is proposed even without colocation, got %+v", c.Derived)
	}
	if c := (&Proposal{}).buildConfig(); c.Derived != nil {
		t.Errorf("no handle found proposes no default handle, got %+v", c.Derived)
	}
	p = &Proposal{Colocated: true, HasTest: true, TestHandle: "data-testid"}
	c = p.buildConfig()
	if c.Derived.TestHandle != "data-testid" || len(c.Derived.Files) != 1 {
		t.Errorf("with colocation the handle joins the colocation templates, got %+v", c.Derived)
	}
}

func TestBuildConfigProposesNoArtifactLayerNorGoverns(t *testing.T) {
	t.Run("BLCNB-X01: The proposal creates no artifact layer and no governs rule", func(t *testing.T) {})
	p := &Proposal{
		CodeDirs: []string{"src/app"}, CodeExts: []string{".go"},
		HasSpecMD: true, HasFeature: true, HasTest: true, GuideDir: "guides", PlanDir: "plans",
		GuideFiles: []string{"guides/A.md"},
	}
	c := p.buildConfig()
	if !reflect.DeepEqual(CodeLayerNames(c), []string{"app-code"}) || len(c.Layers) != 1 {
		t.Errorf("only the code layer should be proposed, got %v", c.Layers)
	}
	if len(c.Governs) != 0 {
		t.Errorf("governs should start empty, got %+v", c.Governs)
	}
}

// The spec is the ANCHOR of colocation, never one of its own derivatives. Detected
// specs used to be listed among the templates, so the proposal derived the spec from itself.
func TestBuildConfigColocationNeverDerivesTheSpec(t *testing.T) {
	t.Run("BLCNB-B04: Colocation is proposed only when detected, with templates only for the detected kinds", func(t *testing.T) {})
	p := &Proposal{Colocated: true, HasSpecMD: true, HasFeature: true, HasTest: true}
	c := p.buildConfig()
	if c.Derived == nil || c.Derived.Anchor != "spec" {
		t.Fatalf("detected colocation should be anchored on the spec, got %+v", c.Derived)
	}
	if _, derived := c.Derived.Files["spec"]; derived {
		t.Errorf("the spec is the anchor and must not be among the templates, got %v", c.Derived.Files)
	}
	if len(c.Derived.Files) != 2 {
		t.Errorf("feature and test are the derivatives, got %v", c.Derived.Files)
	}
}

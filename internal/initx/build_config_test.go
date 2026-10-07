// @anchors
//   code: BCTBL
//   ref: BLCNB

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
	// A folder with another code folder beneath it covers only its own files: a
	// recursive pattern put every handler in two layers.
	p = &Proposal{CodeDirs: []string{"src/handlers", "src", "."}, CodeExts: []string{".go"}}
	ls := p.buildConfig().Layers
	for name, want := range map[string]string{"handlers-code": "src/handlers/**/*.go", "src-code": "src/*.go", "root-code": "*.go"} {
		if got := ls[name].Pattern; got != want {
			t.Errorf("%s pattern = %q, want %q", name, got, want)
		}
	}
}

func TestBuildConfigNamesFoldersSharingASegmentByPath(t *testing.T) {
	t.Run("BLCNB-B01: Folders sharing the last segment are named by their whole path, and the root is root-code", func(t *testing.T) {})
	p := &Proposal{CodeDirs: []string{"src/handlers", "lib/handlers", "."}, CodeExts: []string{".go"}}
	ls := p.buildConfig().Layers
	for _, name := range []string{"src-handlers-code", "lib-handlers-code", "root-code"} {
		if _, ok := ls[name]; !ok {
			t.Errorf("layer %q missing, got %v", name, keysOf(ls))
		}
	}
	if len(ls) != 3 {
		t.Errorf("one layer per folder, got %v", keysOf(ls))
	}
}

func TestBuildConfigExcludesTestsByConvention(t *testing.T) {
	t.Run("BLCNB-B03: A code layer excludes the test files by the project's convention", func(t *testing.T) {})
	goTests := &Proposal{CodeDirs: []string{"pkg"}, CodeExts: []string{".go"}, TestConventions: []TestConvention{{Suffix: "_test.go", Ext: "go", Family: "go"}}}
	if got := goTests.buildConfig().Layers["pkg-code"].Exclude; !reflect.DeepEqual(got, []string{"**/*.spec.md", "**/*.feature", "**/*_test.go"}) {
		t.Errorf("Go tests: exclude = %v", got)
	}
	py := &Proposal{CodeDirs: []string{"app"}, CodeExts: []string{".py"}, Family: "python"}
	if got := py.buildConfig().Layers["app-code"].Exclude; !reflect.DeepEqual(got, []string{"**/*.spec.md", "**/*.feature", "**/test_*.py"}) {
		t.Errorf("Python with no test: exclude = %v", got)
	}
}

func TestBuildConfigTestTemplateByConvention(t *testing.T) {
	t.Run("BLCNB-B06: The colocated test template follows the project's test convention", func(t *testing.T) {})
	p := &Proposal{Colocated: true, HasTest: true, TestConventions: []TestConvention{{Suffix: "_test.go", Ext: "go", Family: "go"}}}
	if got := p.buildConfig().Derived.Files["test"]; !reflect.DeepEqual(got, config.Padroes{"{{dir}}/{{name}}_test.go"}) {
		t.Errorf("template = %v, want {{dir}}/{{name}}_test.go", got)
	}
	p = &Proposal{Colocated: true, HasTest: true}
	if got := p.buildConfig().Derived.Files["test"]; !reflect.DeepEqual(got, config.Padroes{"{{dir}}/{{name}}.test.{{ext}}"}) {
		t.Errorf("template = %v, want the generic form", got)
	}
}

func TestBuildConfigDialectFamily(t *testing.T) {
	t.Run("BLCNB-B07: The proposal's dialect is the family inference found", func(t *testing.T) {})
	if c := (&Proposal{Family: "go"}).buildConfig(); c.Dialect == nil || c.Dialect.Family != "go" {
		t.Errorf("dialect = %+v, want family go", c.Dialect)
	}
	if c := (&Proposal{}).buildConfig(); c.Dialect != nil {
		t.Errorf("no family proposes no dialect, got %+v", c.Dialect)
	}
}

func TestTestPatternByConvention(t *testing.T) {
	t.Run("BLCNB-B08: The test layer's pattern is the project's test convention", func(t *testing.T) {})
	two := &Proposal{TestConventions: []TestConvention{{Suffix: ".spec.ts", Ext: "ts"}, {Suffix: ".test.ts", Ext: "ts"}}}
	for _, c := range []struct {
		p    *Proposal
		want string
	}{
		{two, "{**/*.spec.ts,**/*.test.ts}"},
		{&Proposal{Family: "go"}, "**/*_test.go"},
		{&Proposal{}, "**/*.test.*"},
	} {
		if got := c.p.TestPattern(); got != c.want {
			t.Errorf("TestPattern = %q, want %q", got, c.want)
		}
	}
}

func keysOf[V any](m map[string]V) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
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

func TestBuildConfigSeedsTheFamilysFalliblePatterns(t *testing.T) {
	t.Run("BLCNB-B09: A new project's dialect carries what its family knows can fail", func(t *testing.T) {})
	c := (&Proposal{Family: "ts"}).buildConfig()
	if c.Dialect == nil || len(c.Dialect.FalliblePatterns) == 0 || !reflect.DeepEqual(c.Dialect.FalliblePatterns, config.FamilyFalliblePatterns("ts")) {
		t.Errorf("the ts project starts with the family's fallible patterns: %+v", c.Dialect)
	}
	if c := (&Proposal{}).buildConfig(); c.Dialect != nil {
		t.Errorf("no family, no dialect: %+v", c.Dialect)
	}
}

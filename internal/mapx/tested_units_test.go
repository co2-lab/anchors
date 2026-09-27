package mapx

import (
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func derivedCfg(anchor string, files map[string]config.Padroes, ov ...config.DerivedOverride) *config.Config {
	return &config.Config{Derived: &config.Derived{Anchor: anchor, Files: files, Overrides: ov}}
}

func graphOf(code, tests []string, layer string) *Graph {
	g := &Graph{}
	for _, c := range code {
		g.Nodes = append(g.Nodes, Node{ID: c, Kind: KindCode, Layer: layer})
	}
	for _, t := range tests {
		g.Nodes = append(g.Nodes, Node{ID: t, Kind: KindTest})
	}
	return g
}

func TestTestedUnits_CodeAnchor(t *testing.T) {
	t.Run("TSUNT-B01: With code as the anchor a code file tests through its own name", func(t *testing.T) {})
	cfg := derivedCfg("code", map[string]config.Padroes{"test": {"{{dir}}/{{name}}.test.{{ext}}"}})
	g := graphOf([]string{"utils/fmt.ts", "ui/Home.tsx"}, []string{"utils/fmt.test.ts", "ui/Home.test.tsx", "ui/Other.test.tsx"}, "x")
	want := map[string][]string{"utils/fmt.test.ts": {"utils/fmt.ts"}, "ui/Home.test.tsx": {"ui/Home.tsx"}}
	if got := TestedUnits(g, cfg); !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestTestedUnits_SpecAnchor(t *testing.T) {
	t.Run("TSUNT-B02: With the spec as the anchor the code path gives the variables", func(t *testing.T) {})
	cfg := derivedCfg("spec", map[string]config.Padroes{"code": {"{{dir}}/{{name}}.go"}, "test": {"{{dir}}/{{name}}_test.go"}})
	g := graphOf([]string{"internal/queue/queue.go"}, []string{"internal/queue/queue_test.go"}, "x")
	if got := TestedUnits(g, cfg)["internal/queue/queue_test.go"]; !reflect.DeepEqual(got, []string{"internal/queue/queue.go"}) {
		t.Fatalf("the code template's name finds the test, got %v", got)
	}
}

func TestTestedUnits_LayerOverride(t *testing.T) {
	t.Run("TSUNT-B03: An override of the code's layer places its test elsewhere", func(t *testing.T) {})
	cfg := derivedCfg("code", map[string]config.Padroes{"test": {"{{dir}}/{{name}}.test.{{ext}}"}},
		config.DerivedOverride{When: "model", Files: map[string]config.Padroes{"test": {"backend/__tests__/models/{{name}}.test.ts"}}})
	g := graphOf([]string{"backend/models/user.ts"}, []string{"backend/__tests__/models/user.test.ts", "backend/models/user.test.ts"}, "model")
	g.Nodes = append(g.Nodes, Node{ID: "ui/Home.tsx", Kind: KindCode, Layer: "screen"}, Node{ID: "ui/Home.test.tsx", Kind: KindTest})
	got := TestedUnits(g, cfg)
	if !reflect.DeepEqual(got["backend/__tests__/models/user.test.ts"], []string{"backend/models/user.ts"}) || got["backend/models/user.test.ts"] != nil {
		t.Fatalf("the model's test is in the override's tree only, got %v", got)
	}
	if !reflect.DeepEqual(got["ui/Home.test.tsx"], []string{"ui/Home.tsx"}) {
		t.Fatalf("another layer keeps the default, got %v", got)
	}
}

func TestTestedUnits_Globs(t *testing.T) {
	t.Run("TSUNT-B04: A glob test template matches the tests it covers, a glob code template is not reversed", func(t *testing.T) {})
	cfg := derivedCfg("code", map[string]config.Padroes{"test": {"{{dir}}/__tests__/**/{{name}}.test.ts"}})
	g := graphOf([]string{"lib/a.ts"}, []string{"lib/__tests__/deep/a.test.ts", "lib/__tests__/deep/b.test.ts"}, "x")
	if got := TestedUnits(g, cfg); !reflect.DeepEqual(got, map[string][]string{"lib/__tests__/deep/a.test.ts": {"lib/a.ts"}}) {
		t.Fatalf("the glob test template finds a's test only, got %v", got)
	}
	cfg = derivedCfg("spec", map[string]config.Padroes{"code": {"src/**/{{name}}.go"}, "test": {"src/{{name}}_test.go"}})
	if got := TestedUnits(graphOf([]string{"src/x/a.go"}, []string{"src/a_test.go"}, "x"), cfg); len(got) != 0 {
		t.Fatalf("a glob code template matches no ordinary path, got %v", got)
	}
	cfg = derivedCfg("spec", map[string]config.Padroes{"code": {"app/[slug]/{{name}}.tsx"}, "test": {"app/[slug]/{{name}}.test.tsx"}})
	if got := TestedUnits(graphOf([]string{"app/[slug]/page.tsx"}, []string{"app/[slug]/page.test.tsx"}, "x"), cfg); !reflect.DeepEqual(got["app/[slug]/page.test.tsx"], []string{"app/[slug]/page.tsx"}) {
		t.Fatalf("a bracketed directory in the code template matches itself, got %v", got)
	}
}

func TestTestedUnits_OnlyTheMapsTests(t *testing.T) {
	t.Run("TSUNT-B05: Only the map's tests are answered, each unit once, in order", func(t *testing.T) {})
	cfg := derivedCfg("code", map[string]config.Padroes{"test": {"{{dir}}/{{name}}.test.ts", "{{dir}}/{{name}}.test.{{ext}}"}})
	g := graphOf([]string{"b/a.ts", "a/a.ts"}, []string{"b/a.test.ts"}, "x")
	got := TestedUnits(g, cfg)
	if !reflect.DeepEqual(got, map[string][]string{"b/a.test.ts": {"b/a.ts"}}) {
		t.Fatalf("a/a.test.ts is not in the map and b/a.ts is listed once, got %v", got)
	}
	two := derivedCfg("spec", map[string]config.Padroes{"code": {"{{dir}}/{{name}}.go", "{{dir}}/{{name}}_extra.go"}, "test": {"{{dir}}/{{name}}_test.go"}})
	if got := TestedUnits(graphOf([]string{"x/a_extra.go", "x/a.go"}, []string{"x/a_test.go"}, "x"), two)["x/a_test.go"]; !reflect.DeepEqual(got, []string{"x/a.go", "x/a_extra.go"}) {
		t.Fatalf("a test's units come in path order, got %v", got)
	}
	if TestedUnits(g, &config.Config{}) != nil || TestedUnits(nil, cfg) != nil {
		t.Fatal("without a derivation or a map there is no answer")
	}
}

// @anchors
//   code: BLTSA
//   ref: GRBLG

package mapx

import (
	"reflect"
	"slices"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

func testCfg() *config.Config {
	return &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":   {Kind: "spec", Tags: []string{"spec"}},
			"screen": {Kind: "code", Tags: []string{"frontend"}},
			"guide":  {Kind: "guide", Tags: []string{"guide"}},
		},
		Derived: &config.Derived{
			Anchor: "code",
			Files: map[string]config.Padroes{
				"spec":    {"{{dir}}/{{name}}.spec.md"},
				"feature": {"{{dir}}/{{name}}.feature"},
				"test":    {"{{dir}}/{{name}}.test.{{ext}}"},
			},
		},
		Governs: []config.GovernRule{
			{From: "guides/SPEC_GUIDE.md", Governs: "spec"},
			{From: "guides/FRONTEND_GUIDE.md", Governs: "frontend"},
		},
	}
}

func testFiles() []scan.File {
	return []scan.File{
		{Path: "src/Login.tsx", Layer: "screen", Kind: "code", Rev: "a"},
		{Path: "src/Login.spec.md", Layer: "spec", Kind: "spec", Rev: "b", Codes: []string{"LOGIX-A01"}},
		{Path: "src/Login.feature", Layer: "feature", Kind: "feature", Rev: "c", Codes: []string{"LOGIX-A01"}},
		{Path: "src/Login.test.tsx", Layer: "test", Kind: "test", Rev: "d", Codes: []string{"LOGIX-A01"}},
		{Path: "guides/SPEC_GUIDE.md", Layer: "guide", Kind: "guide", Rev: "e"},
		{Path: "guides/FRONTEND_GUIDE.md", Layer: "guide", Kind: "guide", Rev: "f"},
	}
}

func hasEdge(g *Graph, from, to string, typ EdgeType) bool {
	for _, e := range g.Edges {
		if e.From == from && e.To == to && e.Type == typ {
			return true
		}
	}
	return false
}

func nodeByID(g *Graph, id string) *Node {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return &g.Nodes[i]
		}
	}
	return nil
}

// ── Nodes and identity ──────────────────────────────────────────────────────────────────────

func TestBuild_nodesCarryTheirLayer(t *testing.T) {
	t.Run("GRBLG-B01: Each scanned file becomes a node with its layer's tags and regime", func(t *testing.T) {})
	t.Run("GRBLG-B06: Only a header marks the identity as declared", func(t *testing.T) {})
	t.Run("GRBLG-B19: Nodes and relations are sorted", func(t *testing.T) {})
	cfg := testCfg()
	l := cfg.Layers["screen"]
	l.Regime = "comportamental"
	cfg.Layers["screen"] = l
	files := testFiles()
	files[1].HeaderCode = "LOGIX"
	g := Build(files, cfg, map[string]string{"src/Login.tsx": "2026-09-01"})

	if len(g.Nodes) != len(files) {
		t.Fatalf("one node per scanned file: got %d for %d files", len(g.Nodes), len(files))
	}
	code := nodeByID(g, "src/Login.tsx")
	if code.Kind != KindCode || code.Rev != "a" || code.Layer != "screen" || code.UpdatedAt != "2026-09-01" ||
		!slices.Equal(code.Tags, []string{"frontend"}) || code.Regime != "comportamental" {
		t.Errorf("the node should carry its file and layer facts, got %+v", *code)
	}
	if !nodeByID(g, "src/Login.spec.md").CodeDeclarado {
		t.Error("the spec declares its code in the header and should be marked declared")
	}
	if n := nodeByID(g, "src/Login.test.tsx"); n.CodeDeclarado {
		t.Errorf("a test with no header only cites a code, it does not declare one: %+v", *n)
	}
	// sorted: nodes by path, edges by (from, to, type)
	for i := 1; i < len(g.Nodes); i++ {
		if g.Nodes[i-1].ID > g.Nodes[i].ID {
			t.Errorf("nodes out of order: %s before %s", g.Nodes[i-1].ID, g.Nodes[i].ID)
		}
	}
	for i := 1; i < len(g.Edges); i++ {
		a, b := g.Edges[i-1], g.Edges[i]
		if a.From > b.From || (a.From == b.From && (a.To > b.To || (a.To == b.To && a.Type > b.Type))) {
			t.Errorf("edges out of order: %+v before %+v", a, b)
		}
	}
}

func TestBuild_noDatesMeansEmptyDates(t *testing.T) {
	t.Run("GRBLG-X01: With no dates given, nodes carry no date", func(t *testing.T) {})
	g := Build(testFiles(), testCfg(), nil)
	for _, n := range g.Nodes {
		if n.UpdatedAt != "" {
			t.Errorf("with no dates from the caller, %s should have none, got %q", n.ID, n.UpdatedAt)
		}
	}
}

// The node's identity is the one DECLARED in the header. It used to be inferred from the first
// scenario code in the text — and a spec that CITES another unit before defining its own entered
// the map with the wrong identity. Measured: a model spec opening with a reference to `DTAXX-B11`
// was registered as the owner of `DTAXX`, and every relational gate confronted the wrong unit.
func TestNodeCodePreferHeader(t *testing.T) {
	t.Run("GRBLG-B02: The declared identity wins over cited codes and over the anchor", func(t *testing.T) {})
	t.Run("GRBLG-B04: With no header and no anchor, the first code's root is the identity", func(t *testing.T) {})
	cases := []struct {
		name string
		f    scan.File
		want string
	}{
		{"the declared header wins over the citation that comes first",
			scan.File{HeaderCode: "MTENX", Codes: []string{"DTAXX-B11", "MTENX-B01"}}, "MTENX"},
		{"with no header, it is inferred from the first code (fallback)",
			scan.File{Codes: []string{"ABCDX-B01"}}, "ABCDX"},
		{"no header and no code: empty",
			scan.File{}, ""},
		{"the declared header holds even with no scenario code in the body",
			scan.File{HeaderCode: "WXYZX"}, "WXYZX"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nodeCode(c.f, nil); got != c.want {
				t.Fatalf("code = %q, want %q", got, c.want)
			}
		})
	}
}

// THE SIBLING ANCHOR beats the inference from the text, for a derived artefact.
//
// Measured in the reference app: `GoLiveChecklist.test.ts` cited `ELKAD-B01` in a data string, and the
// file entered the map as `code: ELKAD`. `scenario-coverage` then charged 21 scenarios of OTHER
// specs, and the three invariants that test really proved appeared unproven.
func TestNodeCode_siblingAnchorBeatsTheDerivedText(t *testing.T) {
	t.Run("GRBLG-B03: A derived file takes its sibling anchor's identity", func(t *testing.T) {})
	f := scan.File{
		Path:  "packages/infra/GoLiveChecklist.test.ts",
		Kind:  "test",
		Codes: []string{"ELKAD-B01", "GLCGL-I01"},
	}
	anchors := map[string]string{"packages/infra/GoLiveChecklist.test.ts": "GLCGL"}

	if got := nodeCode(f, anchors); got != "GLCGL" {
		t.Errorf("code = %q, want GLCGL — the test data was read as a declaration", got)
	}
}

// The HEADER still beats everything: whoever declares is not overridden by the anchor.
func TestNodeCode_headerBeatsTheAnchor(t *testing.T) {
	t.Run("GRBLG-B02: The declared identity wins over cited codes and over the anchor", func(t *testing.T) {})
	f := scan.File{Path: "x/Y.test.ts", HeaderCode: "DECLR", Codes: []string{"OTHER-B01"}}
	anchors := map[string]string{"x/Y.test.ts": "ANCOR"}

	if got := nodeCode(f, anchors); got != "DECLR" {
		t.Errorf("code = %q, want DECLR — the header is where the author says whose file it is", got)
	}
}

// A derived file WITHOUT a resolved anchor falls back: its sibling spec declares no identity,
// and refusing there would leave the node with no code at all.
func TestNodeCode_noAnchorFallsBack(t *testing.T) {
	t.Run("GRBLG-B04: With no header and no anchor, the first code's root is the identity", func(t *testing.T) {})
	f := scan.File{Path: "x/Y.test.ts", Kind: "test", Codes: []string{"ABCDX-B01"}}
	if got := nodeCode(f, map[string]string{}); got != "ABCDX" {
		t.Errorf("code = %q, want ABCDX (fallback)", got)
	}
}

// anchorCodeByDerived links each derived file to its anchor's code by the STEM.
func TestAnchorCodeByDerived_linksByTheStem(t *testing.T) {
	t.Run("GRBLG-B03: A derived file takes its sibling anchor's identity", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{Anchor: "spec"}}
	files := []scan.File{
		{Path: "packages/infra/GoLive.spec.md", Kind: "spec", HeaderCode: "GLCGL"},
		{Path: "packages/infra/GoLive.ts", Kind: "code"},
		{Path: "packages/infra/GoLive.test.ts", Kind: "test"},
		{Path: "packages/infra/GoLive.feature", Kind: "feature"},
		// Another unit in the SAME directory: it cannot receive the neighbour's code.
		{Path: "packages/infra/Outra.test.ts", Kind: "test"},
	}

	got := anchorCodeByDerived(files, cfg)

	for _, p := range []string{"packages/infra/GoLive.ts", "packages/infra/GoLive.test.ts",
		"packages/infra/GoLive.feature"} {
		if got[p] != "GLCGL" {
			t.Errorf("%s → %q, want GLCGL", p, got[p])
		}
	}
	if c, ok := got["packages/infra/Outra.test.ts"]; ok {
		t.Errorf("Outra.test.ts received %q — the stem does not match, the neighbour's code is not its own", c)
	}
}

// A VENDORED pipeline has no local identity, whatever its text quotes: the reference app's
// `anchors-board.yml` entered the map as `FNDTN`, the code of a plan it merely cites.
func TestNodeCode_upstreamHasNone(t *testing.T) {
	t.Run("GRBLG-B05: A vendored file has no local identity", func(t *testing.T) {})
	f := scan.File{Path: ".github/workflows/anchors-board.yml", Upstream: true,
		HeaderCode: "FNDTN", Codes: []string{"FNDTN-W04"}}
	if got := nodeCode(f, map[string]string{f.Path: "ABCDX"}); got != "" {
		t.Fatalf("an upstream-owned file gets no local code, got %q", got)
	}
	g := Build([]scan.File{f}, &config.Config{}, nil)
	if n := g.Nodes[0]; !n.Upstream || n.Code != "" {
		t.Fatalf("the node should be upstream with no code, got upstream=%v code=%q", n.Upstream, n.Code)
	}
}

// ── Co-location ─────────────────────────────────────────────────────────────────────────────

func TestBuild_colocation(t *testing.T) {
	t.Run("GRBLG-B07: The pieces of one unit are linked", func(t *testing.T) {})
	g := Build(testFiles(), testCfg(), nil)

	// the unit links: spec→code (specifies), spec→feature (covered-by), feature→test
	cases := []struct {
		from, to string
		typ      EdgeType
	}{
		{"src/Login.spec.md", "src/Login.tsx", EdgeSpecifies},
		{"src/Login.spec.md", "src/Login.feature", EdgeCoveredBy},
		{"src/Login.feature", "src/Login.test.tsx", EdgeTestedBy},
	}
	for _, c := range cases {
		if !hasEdge(g, c.from, c.to, c.typ) {
			t.Errorf("missing co-location edge: %s →%s→ %s", c.from, c.typ, c.to)
		}
	}
}

func TestBuild_derivedOverride(t *testing.T) {
	t.Run("GRBLG-B08: With the code as anchor, the relations still go down from the spec", func(t *testing.T) {})
	t.Run("GRBLG-B10: A layer override reaches the spec through the layer it declares", func(t *testing.T) {})
	// Centralised layout: the handler anchor lives in functions/<module>/handler.ts; spec and
	// feature are co-located, but the TEST lives in a central region
	// (__tests__/unit/lambdas/<module>.test.ts). The override + {{module}} must find the test
	// there — pure co-location (a sibling) would not.
	cfg := &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":    {Kind: "spec", Tags: []string{"spec"}},
			"handler": {Kind: "code", Tags: []string{"backend", "handler"}},
		},
		Derived: &config.Derived{
			Anchor: "code",
			Files: map[string]config.Padroes{
				"spec":    {"{{dir}}/{{name}}.spec.md"},
				"feature": {"{{dir}}/{{name}}.feature"},
				"test":    {"{{dir}}/{{name}}.test.{{ext}}"},
			},
			Overrides: []config.DerivedOverride{
				{When: "handler", Files: map[string]config.Padroes{
					"test": {"__tests__/unit/lambdas/{{module}}.test.ts"},
				}},
			},
		},
	}
	files := []scan.File{
		{Path: "functions/run-audits/handler.ts", Layer: "handler", Kind: "code", Rev: "a"},
		{Path: "functions/run-audits/handler.spec.md", Layer: "spec", Kind: "spec", Rev: "b"},
		{Path: "functions/run-audits/handler.feature", Layer: "feature", Kind: "feature", Rev: "c"},
		{Path: "__tests__/unit/lambdas/run-audits.test.ts", Layer: "test", Kind: "test", Rev: "d"},
	}
	g := Build(files, cfg, nil)

	// spec→feature still links by co-location (siblings)
	if !hasEdge(g, "functions/run-audits/handler.spec.md", "functions/run-audits/handler.feature", EdgeCoveredBy) {
		t.Error("the spec should cover the co-located feature")
	}
	// feature→test links through the {{module}} override (test NOT co-located)
	if !hasEdge(g, "functions/run-audits/handler.feature", "__tests__/unit/lambdas/run-audits.test.ts", EdgeTestedBy) {
		t.Error("the {{module}} override should link the feature to the central test")
	}
}

func TestBuild_specAnchor(t *testing.T) {
	t.Run("GRBLG-B09: Without a feature, the spec is tested by the test", func(t *testing.T) {})
	// The spec as anchor, the canonical form: files derive from it. The code's extension is
	// NOT in the spec's name (`Login.spec.md`), so the derived template has to state it
	// literally — {{ext}} would be empty here.
	cfg := &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":   {Kind: "spec", Tags: []string{"spec"}},
			"screen": {Kind: "code", Tags: []string{"frontend"}},
			"test":   {Kind: "test"},
		},
		Derived: &config.Derived{
			Anchor: "spec",
			Files: map[string]config.Padroes{
				"code": {"{{dir}}/{{name}}.ts"},
				"test": {"{{dir}}/{{name}}.test.ts"},
			},
		},
	}
	files := []scan.File{
		{Path: "src/Login.spec.md", Layer: "spec", Kind: "spec", Rev: "a"},
		{Path: "src/Login.ts", Layer: "screen", Kind: "code", Rev: "b"},
		{Path: "src/Login.test.ts", Layer: "test", Kind: "test", Rev: "c"},
	}
	g := Build(files, cfg, nil)

	if !hasEdge(g, "src/Login.spec.md", "src/Login.ts", EdgeSpecifies) {
		t.Error("the anchor spec should specify the derived code")
	}
	// With no feature declared, the test derives straight from the spec — otherwise a project
	// that uses no feature would lose the whole spec→test link.
	if !hasEdge(g, "src/Login.spec.md", "src/Login.test.ts", EdgeTestedBy) {
		t.Error("with no feature, the spec should link straight to the test")
	}
	// And the direction does NOT flip: the code does not specify the spec.
	if hasEdge(g, "src/Login.ts", "src/Login.spec.md", EdgeSpecifies) {
		t.Error("the code cannot specify the spec — the spec is the anchor")
	}
}

// THE PER-LAYER OVERRIDE must match the layer of the UNIT, not of the file.
//
// A spec matches `**/*.spec.md`, and its `Layer` is `spec`. A `when: screen` compared against it
// never matches — and the per-layer override would be useless precisely for the anchor, which is
// the one that looks it up to resolve its derived files.
func TestLayerOfUnit_theHeaderBeatsTheFile(t *testing.T) {
	t.Run("GRBLG-B10: A layer override reaches the spec through the layer it declares", func(t *testing.T) {})
	cases := []struct {
		name string
		f    scan.File
		want string
	}{
		{"the spec declares the unit's layer", scan.File{Layer: "spec", HeaderLayer: "screen"}, "screen"},
		{"with no header, the file is already the unit", scan.File{Layer: "screen"}, "screen"},
		{"an empty header does not override", scan.File{Layer: "lambdas", HeaderLayer: ""}, "lambdas"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := layerOfUnit(c.f); got != c.want {
				t.Errorf("layerOfUnit = %q, want %q", got, c.want)
			}
		})
	}
}

// And the effect on the MAP: a screen's derived code is `.tsx`, because the override matches.
// Measured in the reference project: `unit-complete` answered "the code is missing" with the
// file on disk.
func TestBuild_layerOverrideReachesTheSpec(t *testing.T) {
	t.Run("GRBLG-B10: A layer override reaches the spec through the layer it declares", func(t *testing.T) {})
	cfg := &config.Config{
		Layers: map[string]config.Layer{
			"spec":   {Pattern: "**/*.spec.md", Kind: "spec"},
			"screen": {Pattern: "app/screens/**/*.tsx", Kind: "code"},
		},
		Derived: &config.Derived{
			Anchor: "spec",
			Files: map[string]config.Padroes{
				"code": {"{{dir}}/{{name}}.ts"},
			},
			Overrides: []config.DerivedOverride{
				{When: "screen", Files: map[string]config.Padroes{
					"code": {"{{dir}}/{{name}}.tsx"},
				}},
			},
		},
	}
	files := []scan.File{
		{Path: "app/screens/Tela.spec.md", Layer: "spec", Kind: "spec",
			HeaderCode: "TELAX", HeaderLayer: "screen"},
		{Path: "app/screens/Tela.tsx", Layer: "screen", Kind: "code"},
	}

	g := Build(files, cfg, nil)

	if !hasEdge(g, "app/screens/Tela.spec.md", "app/screens/Tela.tsx", EdgeSpecifies) {
		t.Errorf("the spec does not point at the `.tsx` — the per-layer override did not reach the anchor.\nEdges: %+v", g.Edges)
	}
}

// A configuration spec governs several files scattered around (`TypeScriptConfig` describes six
// `tsconfig.json`). The per-code override REPLACES the templates of the kinds it declares:
// inheriting `{{name}}.ts` would have the map look for a file nobody will write. A kind it
// does not declare falls back to the default template — it used to be dropped, so an
// override naming only `code:` silently lost the unit's feature and test.
func TestBuild_codeOverrideReplacesTheTemplates(t *testing.T) {
	t.Run("GRBLG-B11: A code override replaces the templates of the kinds it declares, and the others fall back to the default", func(t *testing.T) {})
	t.Run("GRBLG-B12: A bracketed directory is literal and a template wildcard expands", func(t *testing.T) {})
	build := func(ov config.DerivedOverride) *Graph {
		cfg := &config.Config{
			Derived: &config.Derived{
				Anchor: "spec",
				Files: map[string]config.Padroes{
					"code":    {"{{dir}}/{{name}}.ts"},
					"feature": {"{{dir}}/{{name}}.feature"},
					"test":    {"{{dir}}/{{name}}.test.ts"},
				},
				Overrides: []config.DerivedOverride{ov},
			},
		}
		files := []scan.File{
			{Path: "config/TypeScriptConfig.spec.md", Kind: "spec", HeaderCode: "TSCTY"},
			{Path: "config/TypeScriptConfig.ts", Kind: "code"},
			{Path: "config/TypeScriptConfig.feature", Kind: "feature"},
			{Path: "config/TypeScriptConfig.test.ts", Kind: "test"},
			{Path: "config/tsconfig.feature", Kind: "feature"},
			{Path: "packages/a/tsconfig.json", Kind: "code"},
			{Path: "packages/b/tsconfig.json", Kind: "code"},
		}
		return Build(files, cfg, nil)
	}
	spec := "config/TypeScriptConfig.spec.md"

	// only `code:` declared: code replaced, feature and test from the default
	g := build(config.DerivedOverride{Code: "TSCTY", Files: map[string]config.Padroes{"code": {"packages/*/tsconfig.json"}}})
	for _, want := range []string{"packages/a/tsconfig.json", "packages/b/tsconfig.json"} {
		if !hasEdge(g, spec, want, EdgeSpecifies) {
			t.Errorf("the wildcard template should link %s", want)
		}
	}
	if hasEdge(g, spec, "config/TypeScriptConfig.ts", EdgeSpecifies) {
		t.Errorf("the code override replaces the default code template; got %+v", g.Edges)
	}
	if !hasEdge(g, spec, "config/TypeScriptConfig.feature", EdgeCoveredBy) ||
		!hasEdge(g, "config/TypeScriptConfig.feature", "config/TypeScriptConfig.test.ts", EdgeTestedBy) {
		t.Errorf("the kinds the override does not declare must fall back to the default; got %+v", g.Edges)
	}

	// a declared kind replaces the default one; an empty list declares "none"
	g = build(config.DerivedOverride{Code: "TSCTY", Files: map[string]config.Padroes{
		"code":    {"packages/*/tsconfig.json"},
		"feature": {"config/tsconfig.feature"},
		"test":    {},
	}})
	if !hasEdge(g, spec, "config/tsconfig.feature", EdgeCoveredBy) || hasEdge(g, spec, "config/TypeScriptConfig.feature", EdgeCoveredBy) {
		t.Errorf("a declared feature must replace the default one; got %+v", g.Edges)
	}
	for _, e := range g.Edges {
		if e.To == "config/TypeScriptConfig.test.ts" {
			t.Errorf("an empty test list must link no test; got %+v", e)
		}
	}
}

// A directory with brackets (a Next.js route `app/selo/[slug]/`) is a literal path, not a
// character class: the spec there finds its code, feature and test. Unescaped, `[slug]` matched
// one letter and `unit-complete` reported code and feature missing (reported from the reference app,
// 2026-09-25).
func TestBuild_colocationInABracketDirectory(t *testing.T) {
	t.Run("GRBLG-B12: A bracketed directory is literal and a template wildcard expands", func(t *testing.T) {})
	d := "apps/landing-page/src/app/selo/[slug]"
	files := []scan.File{
		{Path: d + "/SeloClient.tsx", Layer: "screen", Kind: "code", Rev: "a"},
		{Path: d + "/SeloClient.spec.md", Layer: "spec", Kind: "spec", Rev: "b", Codes: []string{"SELOC-A01"}},
		{Path: d + "/SeloClient.feature", Layer: "feature", Kind: "feature", Rev: "c", Codes: []string{"SELOC-A01"}},
		{Path: d + "/SeloClient.test.tsx", Layer: "test", Kind: "test", Rev: "d", Codes: []string{"SELOC-A01"}},
	}
	g := Build(files, testCfg(), nil)
	for _, c := range []struct {
		from, to string
		typ      EdgeType
	}{
		{d + "/SeloClient.spec.md", d + "/SeloClient.tsx", EdgeSpecifies},
		{d + "/SeloClient.spec.md", d + "/SeloClient.feature", EdgeCoveredBy},
		{d + "/SeloClient.feature", d + "/SeloClient.test.tsx", EdgeTestedBy},
	} {
		if !hasEdge(g, c.from, c.to, c.typ) {
			t.Errorf("missing co-location edge in a bracket directory: %s →%s→ %s", c.from, c.typ, c.to)
		}
	}
}

func TestGlobEscape(t *testing.T) {
	t.Run("GRBLG-B12: A bracketed directory is literal and a template wildcard expands", func(t *testing.T) {})
	for in, want := range map[string]string{
		"app/selo/[slug]": `app/selo/\[slug\]`,
		"a*b?c":           `a\*b\?c`,
		`x\y`:             `x\\y`,
		"plain/dir":       "plain/dir",
	} {
		if got := globEscape(in); got != want {
			t.Errorf("globEscape(%q) = %q, want %q", in, got, want)
		}
	}
}

// ── Scenario and declared relations ─────────────────────────────────────────────────────────

func TestBuild_noScenarioDupOnColocated(t *testing.T) {
	t.Run("GRBLG-B13: A shared own code links spec and test across directories only", func(t *testing.T) {})
	// LOGIX-A01 is in a co-located spec/feature/test; it must not add an identity edge
	// duplicating the co-location.
	g := Build(testFiles(), testCfg(), nil)
	for _, e := range g.Edges {
		if e.Origin == OriginInferred {
			t.Errorf("there should be no inferred edge between co-located files; got %+v", e)
		}
	}
}

// Only the unit's OWN code links. Citing is not declaring: measured in the reference project,
// 40 tested-by edges into one test, from specs that only mentioned it in prose.
func TestBuild_scenarioEdgesCrossDirectoriesForOwnCodes(t *testing.T) {
	t.Run("GRBLG-B13: A shared own code links spec and test across directories only", func(t *testing.T) {})
	files := []scan.File{
		{Path: "specs/Order.spec.md", Kind: "spec", HeaderCode: "ORDER", Codes: []string{"ORDER-B01"}},
		{Path: "tests/order_flow.test.ts", Kind: "test", Codes: []string{"ORDER-B01"}},
		{Path: "specs/Cites.spec.md", Kind: "spec", HeaderCode: "CITES", Codes: []string{"CITES-B01", "ORDER-B01"}},
	}
	g := Build(files, &config.Config{}, nil)
	var found bool
	for _, e := range g.Edges {
		if e.From == "specs/Order.spec.md" && e.To == "tests/order_flow.test.ts" && e.Type == EdgeTestedBy && e.Origin == OriginInferred {
			found = true
		}
	}
	if !found {
		t.Errorf("the spec and the test sharing ORDER-B01 across directories should be linked; got %+v", g.Edges)
	}
	if hasEdge(g, "specs/Cites.spec.md", "tests/order_flow.test.ts", EdgeTestedBy) {
		t.Error("a spec that only cites ORDER-B01 is not tested by ORDER's test")
	}
}

func TestBuild_scenarioEdgesReadTheUnitFromTheRef(t *testing.T) {
	t.Run("GRBLG-B25: A file with a code of its own and a ref links across directories through the unit it refs", func(t *testing.T) {})
	files := []scan.File{
		{Path: "lib/stack.feature", Kind: "feature", HeaderCode: "ICSFN", HeaderRefs: []string{"LPSTI"}, Codes: []string{"LPSTI-B01"}},
		{Path: "test/stack.test.ts", Kind: "test", HeaderCode: "ICTNF", HeaderRefs: []string{"LPSTI"}, Codes: []string{"LPSTI-B01", "OTHER-B01"}},
		{Path: "other/x.feature", Kind: "feature", HeaderCode: "OTHFT", HeaderRefs: []string{"OTHER"}, Codes: []string{"OTHER-B01"}},
	}
	g := Build(files, &config.Config{}, nil)
	if !hasEdge(g, "lib/stack.feature", "test/stack.test.ts", EdgeTestedBy) {
		t.Errorf("the feature and the test of LPSTI are linked across directories; got %+v", g.Edges)
	}
	if hasEdge(g, "other/x.feature", "test/stack.test.ts", EdgeTestedBy) {
		t.Error("a code the test only cites does not link it to that unit's feature")
	}
}

func TestBuild_governsByTag(t *testing.T) {
	t.Run("GRBLG-B14: A guide governs the layers of its tag only, and never itself", func(t *testing.T) {})
	g := Build(testFiles(), testCfg(), nil)

	// SPEC_GUIDE (tag spec) governs the spec
	if !hasEdge(g, "guides/SPEC_GUIDE.md", "src/Login.spec.md", EdgeGoverns) {
		t.Error("SPEC_GUIDE should govern the spec (tag spec)")
	}
	// FRONTEND_GUIDE (tag frontend) governs the screen code
	if !hasEdge(g, "guides/FRONTEND_GUIDE.md", "src/Login.tsx", EdgeGoverns) {
		t.Error("FRONTEND_GUIDE should govern the code (tag frontend)")
	}
	// FRONTEND_GUIDE does NOT govern the spec (wrong tag) — no cartesian product
	if hasEdge(g, "guides/FRONTEND_GUIDE.md", "src/Login.spec.md", EdgeGoverns) {
		t.Error("FRONTEND_GUIDE should NOT govern the spec (it has no spec tag)")
	}

	// A guide whose own layer carries the tag it rules governs its siblings, never itself.
	cfg := testCfg()
	cfg.Layers["guide"] = config.Layer{Kind: "guide", Tags: []string{"guide", "spec"}}
	g = Build(testFiles(), cfg, nil)
	if hasEdge(g, "guides/SPEC_GUIDE.md", "guides/SPEC_GUIDE.md", EdgeGoverns) {
		t.Error("a guide never governs itself")
	}
	if !hasEdge(g, "guides/SPEC_GUIDE.md", "guides/FRONTEND_GUIDE.md", EdgeGoverns) {
		t.Error("the other guide carries the tag and should be governed")
	}
}

func TestDependsOnEdges(t *testing.T) {
	t.Run("GRBLG-B15: A dependency row becomes a relation to an existing file", func(t *testing.T) {})
	t.Run("GRBLG-X02: No relation points to a file that was not scanned", func(t *testing.T) {})
	files := []scan.File{
		{Path: "src/Login.spec.md", Layer: "spec", Kind: "spec", Rev: "b",
			Deps: []scan.Dep{
				{Code: "DEP1", File: "src/auth.store.ts", Method: "useAuthStore", Layer: "store"},
				{Code: "DEP2", File: "src/useAuth.ts", Method: "signIn", Layer: "hook"},
				{Code: "DEP3", File: "src/missing.ts", Method: "x", Layer: "hook"}, // absent target → no edge
			}},
		{Path: "src/auth.store.ts", Layer: "store", Kind: "code", Rev: "c"},
		{Path: "src/useAuth.ts", Layer: "hook", Kind: "code", Rev: "d"},
	}
	edges := dependsOnEdges(files)
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges (the 3rd points at a missing file), got %d: %+v", len(edges), edges)
	}
	// the edge goes from the SPEC to the FILE, with method+dep as metadata, origin declared
	var e1 *Edge
	for i := range edges {
		if edges[i].To == "src/auth.store.ts" {
			e1 = &edges[i]
		}
	}
	if e1 == nil {
		t.Fatal("edge to auth.store.ts not built")
	}
	if e1.From != "src/Login.spec.md" || e1.Type != EdgeDependsOn || e1.Origin != OriginDeclared {
		t.Errorf("malformed edge: %+v", *e1)
	}
	if e1.Method != "useAuthStore" || e1.Dep != "DEP1" {
		t.Errorf("method/dep metadata lost: %+v", *e1)
	}
}

func TestDependsOnEdgesIntegratedInBuild(t *testing.T) {
	t.Run("GRBLG-B15: A dependency row becomes a relation to an existing file", func(t *testing.T) {})
	files := append(testFiles(),
		scan.File{Path: "src/useLogin.ts", Layer: "hook", Kind: "code", Rev: "g"},
	)
	// injects a dep on the Login spec pointing at the hook
	for i := range files {
		if files[i].Path == "src/Login.spec.md" {
			files[i].Deps = []scan.Dep{{Code: "DEP1", File: "src/useLogin.ts", Method: "run", Layer: "hook"}}
		}
	}
	g := Build(files, testCfg(), nil)
	if !hasEdge(g, "src/Login.spec.md", "src/useLogin.ts", EdgeDependsOn) {
		t.Error("Build should hold the depends-on edge spec→hook")
	}
}

// A plan cites a spec in two legitimate ways — by path, or by NAME only, which is how prose
// writes it. Measured in a real repository: 10 of 26 citations are by name.
func TestSeedResolvedByName(t *testing.T) {
	t.Run("GRBLG-B16: A seed path is exact and a bare name must be unique", func(t *testing.T) {})
	files := []scan.File{
		{Path: "plans/p.md", Kind: "plan", Seeds: []string{"Tela.spec.md", "apps/y/Outra.spec.md", "Ambigua.spec.md"}},
		{Path: "apps/x/Tela.spec.md", Kind: "spec"},
		{Path: "apps/y/Outra.spec.md", Kind: "spec"},
		// two with the same name: the citation is ambiguous and must NOT produce an edge —
		// choosing one would invent a link the author did not declare.
		{Path: "apps/a/Ambigua.spec.md", Kind: "spec"},
		{Path: "apps/b/Ambigua.spec.md", Kind: "spec"},
	}
	targets := map[string]bool{}
	for _, e := range seedEdges(files) {
		targets[e.To] = true
	}
	if !targets["apps/x/Tela.spec.md"] {
		t.Error("a citation by NAME should resolve to the unique path")
	}
	if !targets["apps/y/Outra.spec.md"] {
		t.Error("a citation by PATH should keep working")
	}
	if targets["apps/a/Ambigua.spec.md"] || targets["apps/b/Ambigua.spec.md"] {
		t.Error("a duplicated name is ambiguous — it cannot produce an edge")
	}
}

// THE REAL CASE, measured in the reference app: plan 0010 (Redis) seeds
// `packages/lambdas/redis/InstanceList.spec.md`, and plan 0009 (Database) had already delivered
// `packages/lambdas/database/InstanceList.spec.md`. The by-name fallback — made for citations in
// PROSE — was applied to a seed carrying the WHOLE PATH, and the map linked plan 0010 to plan
// 0009's spec. A DECLARED path is a path.
func TestBuild_seedWithPathDoesNotResolveByNameInAnotherDirectory(t *testing.T) {
	t.Run("GRBLG-B16: A seed path is exact and a bare name must be unique", func(t *testing.T) {})
	files := []scan.File{
		{Path: "plans/0010-redis.md", Kind: "plan", Seeds: []string{"packages/lambdas/redis/InstanceList.spec.md"}},
		{Path: "packages/lambdas/database/InstanceList.spec.md", Kind: "spec"},
	}

	g := Build(files, &config.Config{}, nil)

	for _, e := range g.Edges {
		if e.Type == EdgeSeeds && e.From == "plans/0010-redis.md" && e.To == "packages/lambdas/database/InstanceList.spec.md" {
			t.Fatalf("the REDIS plan ended up seeding the DATABASE spec:\n  %s → %s\n"+
				"  the declared path was packages/lambdas/redis/InstanceList.spec.md", e.From, e.To)
		}
	}
}

// A seed whose declared path EXISTS still becomes an edge — the normal case.
func TestBuild_seedWithExistingPathBecomesAnEdge(t *testing.T) {
	t.Run("GRBLG-B16: A seed path is exact and a bare name must be unique", func(t *testing.T) {})
	files := []scan.File{
		{Path: "plans/0010-redis.md", Kind: "plan", Seeds: []string{"packages/lambdas/redis/InstanceList.spec.md"}},
		{Path: "packages/lambdas/redis/InstanceList.spec.md", Kind: "spec"},
	}

	g := Build(files, &config.Config{}, nil)

	if !hasEdge(g, "plans/0010-redis.md", "packages/lambdas/redis/InstanceList.spec.md", EdgeSeeds) {
		t.Error("the seed whose path exists did not become an edge")
	}
}

// And a path that exists NOWHERE does not become an edge by guesswork either: the plan promises
// to create a file, and the map must not invent that it is already there.
func TestBuild_missingSeedDoesNotLinkAnotherFile(t *testing.T) {
	t.Run("GRBLG-B16: A seed path is exact and a bare name must be unique", func(t *testing.T) {})
	files := []scan.File{
		{Path: "plans/0010-redis.md", Kind: "plan", Seeds: []string{"packages/lambdas/redis/Thing.spec.md"}},
		{Path: "packages/other/place/Thing.spec.md", Kind: "spec"},
	}

	g := Build(files, &config.Config{}, nil)

	for _, e := range g.Edges {
		if e.Type == EdgeSeeds && e.To == "packages/other/place/Thing.spec.md" {
			t.Errorf("a declared path that does not exist was resolved by name: %s → %s", e.From, e.To)
		}
	}
}

func TestBuild_needsLinksOnlyExistingPlans(t *testing.T) {
	t.Run("GRBLG-B17: A need links only an existing plan", func(t *testing.T) {})
	files := []scan.File{
		{Path: "plans/0002.md", Kind: "plan", Needs: []string{"plans/0001.md", "plans/0009.md"}},
		{Path: "plans/0001.md", Kind: "plan"},
	}
	g := Build(files, &config.Config{}, nil)
	if !hasEdge(g, "plans/0002.md", "plans/0001.md", EdgeNeeds) {
		t.Error("a need on an existing plan should become a needs edge")
	}
	for _, e := range g.Edges {
		if e.To == "plans/0009.md" {
			t.Errorf("a need on a missing plan must not become an edge, got %+v", e)
		}
	}
}

func TestBuild_realizesAndGatedByResolveByUnitCode(t *testing.T) {
	t.Run("GRBLG-B18: Realized doctrine and flag scenarios resolve by unit code", func(t *testing.T) {})
	files := []scan.File{
		{Path: "src/Credit.spec.md", Kind: "spec", HeaderCode: "CREDT",
			Realizes: []scan.Realizes{{From: "CREDT-V01", To: "LIMIT-R03"}, {From: "CREDT-V02", To: "NOPE-R01"}},
			GatedBy:  []scan.Realizes{{From: "CREDT-V03", To: "CHKUT-G02"}}},
		{Path: "product/limit.doctrine.md", Kind: string(KindProduct), HeaderCode: "LIMIT"},
		{Path: "flags/new-checkout.flag.md", Kind: string(KindFlag), HeaderCode: "CHKUT"},
	}
	g := Build(files, &config.Config{}, nil)
	var realizes, gated *Edge
	for i := range g.Edges {
		e := &g.Edges[i]
		switch e.Type {
		case EdgeRealizes:
			realizes = e
		case EdgeGatedBy:
			gated = e
		}
	}
	if realizes == nil || realizes.From != "src/Credit.spec.md" || realizes.To != "product/limit.doctrine.md" ||
		realizes.Dep != "CREDT-V01" || realizes.Method != "LIMIT-R03" {
		t.Errorf("the spec should realize the doctrine, carrying both rules; got %+v", realizes)
	}
	if gated == nil || gated.To != "flags/new-checkout.flag.md" || gated.Dep != "CREDT-V03" || gated.Method != "CHKUT-G02" {
		t.Errorf("the spec should be gated by the flag, carrying both rules; got %+v", gated)
	}
	n := 0
	for _, e := range g.Edges {
		if e.Type == EdgeRealizes {
			n++
		}
	}
	if n != 1 {
		t.Errorf("an unknown unit code (NOPE) links nothing; got %d realizes edges", n)
	}
}

func TestBuild_noEdgeToAnUnscannedFile(t *testing.T) {
	t.Run("GRBLG-X02: No relation points to a file that was not scanned", func(t *testing.T) {})
	files := []scan.File{
		{Path: "src/A.spec.md", Kind: "spec", Deps: []scan.Dep{{Code: "DEP1", File: "src/gone.ts"}}},
		{Path: "plans/p.md", Kind: "plan", Seeds: []string{"src/Gone.spec.md"}, Needs: []string{"plans/gone.md"}},
	}
	g := Build(files, &config.Config{}, nil)
	for _, e := range g.Edges {
		if nodeByID(g, e.To) == nil || nodeByID(g, e.From) == nil {
			t.Errorf("edge to a file that was not scanned: %+v", e)
		}
	}
}

// ── Order and rebuild ───────────────────────────────────────────────────────────────────────

func TestBuild_isIndependentOfTheFileOrder(t *testing.T) {
	t.Run("GRBLG-I01: The build does not depend on the order of the files", func(t *testing.T) {})
	files := testFiles()
	files = append(files, scan.File{Path: "tests/login_e2e.test.ts", Kind: "test", Codes: []string{"LOGIX-A01"}})
	files[1].HeaderCode = "LOGIX"
	reversed := slices.Clone(files)
	slices.Reverse(reversed)
	a := Build(files, testCfg(), nil)
	b := Build(reversed, testCfg(), nil)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("the same files in another order built another graph:\n%+v\n%+v", a, b)
	}
}

// `map build` rebuilds the graph, and `check --all` runs build BEFORE confronting: if the rebuild
// erases the judgment, the verdict dies between one command and the next.
func TestJudgmentSurvivesTheRebuild(t *testing.T) {
	t.Run("GRBLG-B20: A rebuild keeps the stamps and judgments of surviving relations", func(t *testing.T) {})
	old := &Graph{
		Nodes: []Node{{ID: "a.spec.md", Rev: "r1"}, {ID: "a.test.ts", Rev: "r1"}},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}
	old.StampNodeByGate("a.spec.md", "ok", "t0", "my-gate")

	// the build assembles a new graph, with the same edges and no stamp
	rebuilt := &Graph{
		Nodes: []Node{{ID: "a.spec.md", Rev: "r1"}, {ID: "a.test.ts", Rev: "r1"}},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}
	PreserveStamps(rebuilt, old)

	if v, ok := rebuilt.JudgedBy("a.spec.md", "my-gate"); !ok || v != "ok" {
		t.Errorf("the judgment should survive the rebuild: verdict=%q ok=%v", v, ok)
	}
}

// Measured: `anchors stale` reported 9,975 of 9,975 edges "never validated" in a repository whose
// 590 specs had been confronted dozens of times — each rebuild threw away what was known.
func TestPreserveStampsMatchesByTypeAndEnds(t *testing.T) {
	t.Run("GRBLG-B20: A rebuild keeps the stamps and judgments of surviving relations", func(t *testing.T) {})
	st := &Stamp{ValidatedFromRev: "r1", ValidatedToRev: "r1", ChangedAt: "2026-09-01", Verdict: "ok"}
	old := &Graph{Edges: []Edge{
		{From: "a.spec.md", To: "a.test.ts", Type: EdgeTestedBy, Stamp: st},
		{From: "a.spec.md", To: "a.go", Type: EdgeSpecifies, Stamp: st},
	}}
	rebuilt := &Graph{Edges: []Edge{
		{From: "a.spec.md", To: "a.test.ts", Type: EdgeTestedBy},
		{From: "a.spec.md", To: "a.go", Type: EdgeDependsOn}, // same ends, another reason
	}}
	PreserveStamps(rebuilt, old)
	if rebuilt.Edges[0].Stamp == nil || *rebuilt.Edges[0].Stamp != *st {
		t.Errorf("the surviving edge should keep its stamp, got %+v", rebuilt.Edges[0].Stamp)
	}
	if rebuilt.Edges[1].Stamp != nil {
		t.Errorf("an edge that changed type is another relation and keeps nothing, got %+v", rebuilt.Edges[1].Stamp)
	}
}

// Measured: after an `ingest --junit` matched 352 files, one `map build` left the tests-pass gate
// at ⚠494 ✓0, as if nobody had ever run the suite. The cut is the REV: a signal of a file that
// changed is no longer valid.
func TestPreserveSignalsOnlyForUnchangedFiles(t *testing.T) {
	t.Run("GRBLG-B21: A rebuild keeps a signal only for an unchanged file", func(t *testing.T) {})
	old := &Graph{Nodes: []Node{
		{ID: "same_test.go", Rev: "r1", Signal: &TestSignal{Passed: 3, AtRev: "r1"}},
		{ID: "edited_test.go", Rev: "r1", Signal: &TestSignal{Passed: 2, AtRev: "r1"}},
		{ID: "kept.spec.md", Rev: "r1", EvidenceKept: []EvidenceKeep{{Reason: "only a date"}}},
		{ID: "moved.spec.md", Rev: "r1", EvidenceKept: []EvidenceKeep{{Reason: "old"}}},
	}}
	rebuilt := &Graph{Nodes: []Node{{ID: "same_test.go", Rev: "r1"}, {ID: "edited_test.go", Rev: "r2"},
		{ID: "kept.spec.md", Rev: "r1"}, {ID: "moved.spec.md", Rev: "r2"}}}
	PreserveStamps(rebuilt, old)
	if k := rebuilt.Nodes[2]; len(k.EvidenceKept) != 1 || k.Signal != nil {
		t.Errorf("an unchanged file keeps its declarations and gains no signal, got %+v", k)
	}
	if k := rebuilt.Nodes[3]; len(k.EvidenceKept) != 0 {
		t.Errorf("an edited file's declarations speak of another revision, got %+v", k.EvidenceKept)
	}
	if s := rebuilt.Nodes[0].Signal; s == nil || s.Passed != 3 {
		t.Errorf("the unchanged file should keep its signal, got %+v", s)
	}
	if s := rebuilt.Nodes[1].Signal; s != nil {
		t.Errorf("the edited file's signal is no longer valid, got %+v", s)
	}
}

func TestBuild_carriesTheSupportMark(t *testing.T) {
	t.Run("GRBLG-B22: A support file becomes a node marked as support", func(t *testing.T) {})
	files := []scan.File{
		{Path: "flows/utils/login.yaml", Layer: "e2e", Kind: "test", Support: true},
		{Path: "flows/screens/home.yaml", Layer: "e2e", Kind: "test"},
	}
	g := Build(files, &config.Config{Layers: map[string]config.Layer{"e2e": {Pattern: "flows/**", Kind: "test"}}}, nil)
	by := map[string]Node{}
	for _, n := range g.Nodes {
		by[n.ID] = n
	}
	if !by["flows/utils/login.yaml"].Support || by["flows/screens/home.yaml"].Support ||
		by["flows/utils/login.yaml"].Kind != KindTest {
		t.Errorf("the support mark must reach the node and the kind stay test, got %+v", by)
	}
}

func TestBuild_edgesAreInATotalOrder(t *testing.T) {
	t.Run("GRBLG-I02: The edges are in a total order", func(t *testing.T) {})
	for _, deps := range [][]scan.Dep{
		{{File: "src/Login.tsx", Method: "`m`", Code: "DEP2"}, {File: "src/Login.tsx", Method: "`m`", Code: "DEP1"}}, // only the code differs
		{{File: "src/Login.tsx", Method: "`b`", Code: "DEP1"}, {File: "src/Login.tsx", Method: "`a`", Code: "DEP1"}}, // only the method differs
	} {
		one := testFiles()
		one[1].Deps = deps
		two := testFiles()
		two[1].Deps = []scan.Dep{deps[1], deps[0]}
		a, b := Build(one, testCfg(), nil), Build(two, testCfg(), nil)
		if !reflect.DeepEqual(a.Edges, b.Edges) {
			t.Fatalf("the same edges in another build order are written in another order:\n%+v\n%+v", a.Edges, b.Edges)
		}
	}
}

func TestFillSignals(t *testing.T) {
	t.Run("GRBLG-B23: Signals are filled from another map at the same revision", func(t *testing.T) {})
	src := &Graph{Nodes: []Node{
		{ID: "a", Rev: "r1", Signal: &TestSignal{Passed: 1}, EvidenceKept: []EvidenceKeep{{Reason: "src"}}},
		{ID: "b", Rev: "r1", Signal: &TestSignal{Passed: 2}},
		{ID: "c", Rev: "r0", Signal: &TestSignal{Passed: 3}},
	}}
	novo := &Graph{Nodes: []Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1", Signal: &TestSignal{Passed: 9}}, {ID: "c", Rev: "r1"}, {ID: "d", Rev: "r1"}}}
	FillSignals(novo, src)
	if s := novo.Nodes[0]; s.Signal == nil || s.Signal.Passed != 1 || len(s.EvidenceKept) != 1 {
		t.Errorf("a node without signal gets the source's, got %+v", s)
	}
	if novo.Nodes[1].Signal.Passed != 9 || novo.Nodes[2].Signal != nil || novo.Nodes[3].Signal != nil {
		t.Errorf("an own signal stays, another revision and an unknown file get nothing, got %+v", novo.Nodes)
	}
	FillSignals(nil, src)
	FillSignals(novo, nil)
}

func TestBuild_aFilesOwnCodeAndItsUnit(t *testing.T) {
	t.Run("GRBLG-B24: A file's own code is its file code, and a ref keeps its unit", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{Anchor: "spec"}}
	files := []scan.File{
		{Path: "ui/Arena.spec.md", Kind: "spec", HeaderCode: "ARENA"},
		{Path: "ui/Arena.tsx", Kind: "code", HeaderCode: "ARSCR", HeaderRefs: []string{"ARENA"}},
		{Path: "theme/tokens.ts", Kind: "code", HeaderCode: "TOKNS"},
	}
	g := Build(files, cfg, nil)
	got := map[string]Node{}
	for _, n := range g.Nodes {
		got[n.ID] = n
	}
	if n := got["ui/Arena.tsx"]; n.Code != "ARENA" || n.FileCode != "ARSCR" || n.CodeDeclarado {
		t.Errorf("the screen's unit is ARENA, its file code ARSCR: %+v", n)
	}
	if n := got["theme/tokens.ts"]; n.Code != "TOKNS" || n.FileCode != "TOKNS" || !n.CodeDeclarado {
		t.Errorf("the helper owns its unit: %+v", n)
	}
}

func TestBuild_flagEdgesByCode(t *testing.T) {
	t.Run("GRBLG-B26: The @dep and @navigates flags become edges to the file whose own code they name", func(t *testing.T) {})
	files := []scan.File{
		{Path: "ui/Arena.tsx", Kind: "code", HeaderCode: "ARNSC",
			CodeDeps:  []scan.CodeDep{{Code: "TOKNS", Symbols: []string{"PALETTE", "withAlpha"}}, {Code: "GHOST"}, {Waiver: "types"}},
			Navigates: []scan.Navigation{{Codes: []string{"WLLTW"}, Rule: "ARNAA-A02"}, {Codes: []string{"NOWHR"}}}},
		{Path: "theme/tokens.ts", Kind: "code", HeaderCode: "TOKNS"},
		{Path: "ui/Wallet.spec.md", Kind: "spec", HeaderCode: "WLLTW"},
	}
	g := Build(files, &config.Config{}, nil)
	var dep, nav *Edge
	for i, e := range g.Edges {
		if e.From == "ui/Arena.tsx" && e.To == "theme/tokens.ts" && e.Type == EdgeDependsOn {
			dep = &g.Edges[i]
		}
		if e.From == "ui/Arena.tsx" && e.To == "ui/Wallet.spec.md" && e.Type == EdgeNavigatesTo {
			nav = &g.Edges[i]
		}
	}
	if dep == nil || dep.Method != "PALETTE, withAlpha" || dep.Origin != OriginDeclared {
		t.Errorf("the @dep is a declared depends-on carrying the symbols: %+v", dep)
	}
	if nav == nil || nav.Method != "ARNAA-A02" {
		t.Errorf("the @navigates is a navigates-to carrying the rule: %+v", nav)
	}
	for _, e := range g.Edges {
		if e.From == "ui/Arena.tsx" && (e.Type == EdgeDependsOn || e.Type == EdgeNavigatesTo) && e.To != "theme/tokens.ts" && e.To != "ui/Wallet.spec.md" {
			t.Errorf("a code no file owns makes no edge: %+v", e)
		}
	}
}

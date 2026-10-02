package doct

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// The anchor is what makes the index work. Getting it wrong produces a clickable link that
// goes nowhere — the browser stays where it is, and the index LOOKS like it works.
func TestGitHubAnchor(t *testing.T) {
	t.Run("DCLND-B01: A heading's anchor follows the GitHub convention and keeps accents", func(t *testing.T) {})
	cases := []struct{ heading, want string }{
		{"GLCGL-B01 — Cada item tem um artefato que o PROVA",
			"glcgl-b01--cada-item-tem-um-artefato-que-o-prova"},
		{"GLCGL — GoLiveChecklist — a régua que registra o que não foi feito",
			"glcgl--golivechecklist--a-régua-que-registra-o-que-não-foi-feito"},
		// ACCENTS STAY: GitHub keeps `é` in the anchor, and transliterating would produce a
		// link that does not resolve. It is the mistake a slug library makes by default.
		{"A régua e a dívida", "a-régua-e-a-dívida"},
		{"Visão Geral", "visão-geral"},
		// Punctuation goes; the spaces around it each become a hyphen, which is GitHub's
		// behaviour.
		{"O que é, e o que não é", "o-que-é-e-o-que-não-é"},
		{"`código` entre crases", "código-entre-crases"},
		{"  espaços nas pontas  ", "espaços-nas-pontas"},
		{"Camada: infra", "camada-infra"},
		{"snake_case name", "snake-case-name"},
	}
	for _, c := range cases {
		if got := GitHubAnchor(c.heading); got != c.want {
			t.Errorf("GitHubAnchor(%q)\n  = %q\n  want %q", c.heading, got, c.want)
		}
	}
}

// EVERY GENERATED LINK MUST RESOLVE. This is the test that was missing when 483 of 812
// links came out broken: the index always built the rule's anchor, and on big layers the
// target page summarizes — the rule is not a heading there, and the anchor does not exist.
func TestBuild_everyGeneratedLinkResolves(t *testing.T) {
	t.Run("DCLND-I01: Every generated link resolves on a small layer and on a big one", func(t *testing.T) {})
	// A SMALL layer (full page) and a BIG one (summary page): the defect only shows on the
	// second, and a test with a single layer would pass falsely.
	specs := map[string]string{}
	for i := 0; i < DefaultLayout().MaxUnits+2; i++ {
		code := fmt.Sprintf("BIG%02d", i)
		specs[fmt.Sprintf("big/U%02d.spec.md", i)] = fmt.Sprintf(`---
code: %s
layer: grande
---

# %s — a unidade %d

## Visão Geral

Texto.

## Regras

### %s-B01 — a regra da unidade %d

Corpo.
`, code, code, i, code, i)
	}
	specs["peq/Only.spec.md"] = specDeExemplo // layer: infra, with rules

	root, g := projetoDeTeste(t, specs)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := c.InitScaffolds(false); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}

	// Both layers must exist, on opposite sides of the cut-off — otherwise the test does
	// not exercise the case that broke.
	if !c.layerIsBig("grande") {
		t.Fatal("layer `grande` should be above the cut-off — the test does not cover the defect")
	}
	if c.layerIsBig("infra") {
		t.Fatal("layer `infra` should be below the cut-off")
	}

	headings, links := collectDocs(t, filepath.Join(root, OutDir))
	if len(links) == 0 {
		t.Fatal("no link generated — the test proves nothing")
	}
	for _, l := range links {
		anc, ok := headings[l.target]
		if !ok {
			t.Errorf("%s points at `%s`, which does not exist", l.origin, l.target)
			continue
		}
		if !anc[l.anchor] {
			t.Errorf("%s points at `%s#%s` — the anchor does not exist on the target",
				l.origin, l.target, l.anchor)
		}
	}
}

type generatedLink struct{ origin, target, anchor string }

var (
	headingRE = regexp.MustCompile(`(?m)^#{1,6}\s+(.+?)\s*$`)
	linkRE    = regexp.MustCompile(`\]\(([^)#]+)#([^)]+)\)`)
)

// collectDocs returns the anchors of each page and every internal link generated.
func collectDocs(t *testing.T, dir string) (map[string]map[string]bool, []generatedLink) {
	t.Helper()
	headings := map[string]map[string]bool{}
	var links []generatedLink
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		rel = filepath.ToSlash(rel)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		set := map[string]bool{}
		for _, m := range headingRE.FindAllStringSubmatch(string(b), -1) {
			set[GitHubAnchor(m[1])] = true
		}
		headings[rel] = set
		base := path.Dir(rel)
		for _, m := range linkRE.FindAllStringSubmatch(string(b), -1) {
			links = append(links, generatedLink{
				origin: rel,
				target: path.Clean(path.Join(base, m[1])),
				anchor: m[2],
			})
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return headings, links
}

// compilerWithLayout is a compiler over the example spec whose layout is set so the layer
// is small (big=false) or big (big=true).
func compilerWithLayout(t *testing.T, big bool) *Compiler {
	t.Helper()
	root, g := projetoComFeature(t)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if big {
		c.Layout = Layout{MaxUnits: 0, MaxLines: 0}
	}
	return c
}

func TestRuleLink(t *testing.T) {
	t.Run("DCLND-B02: A rule's link points at the rule on a small layer and at its unit on a big one", func(t *testing.T) {
		r := Rule{Code: "GLCGL-B01", Titulo: "r"}
		small := compilerWithLayout(t, false)
		if got := small.fnRuleLink("infra", small.specs[0], r); got != "camadas/infra.md#glcgl-b01--r" {
			t.Errorf("small layer: %q, want the rule's own heading", got)
		}
		big := compilerWithLayout(t, true)
		s := big.specs[0]
		want := "camadas/infra.md#" + GitHubAnchor(s.Code+" — "+s.Titulo)
		if got := big.fnRuleLink("infra", s, r); got != want {
			t.Errorf("big layer: %q, want the unit's heading %q", got, want)
		}
	})
}

func TestScenarioLink(t *testing.T) {
	t.Run("DCLND-B03: A scenario's link falls back to the unit, then to the bare page, on a big layer", func(t *testing.T) {
		known := Scenario{Code: "GLCGL-B01", Titulo: "s", Spec: "GLCGL"}
		unknown := Scenario{Code: "NOPEX-B01", Titulo: "u", Spec: "NOPEX"}

		small := compilerWithLayout(t, false)
		if got := small.fnScenarioLink("infra", known); got != "camadas/infra.md#glcgl-b01--s" {
			t.Errorf("small layer, known: %q", got)
		}
		if got := small.fnScenarioLink("infra", unknown); got != "camadas/infra.md#nopex-b01--u" {
			t.Errorf("small layer, unknown spec: %q", got)
		}

		big := compilerWithLayout(t, true)
		s := big.specs[0]
		if got := big.fnScenarioLink("infra", known); got != "camadas/infra.md#"+GitHubAnchor(s.Code+" — "+s.Titulo) {
			t.Errorf("big layer, known: %q, want the unit's heading", got)
		}
		if got := big.fnScenarioLink("infra", unknown); got != "camadas/infra.md" {
			t.Errorf("big layer, unknown spec: %q, want the bare page", got)
		}
	})
}

// Memoizing `fnSize` is pure optimization: the output must be IDENTICAL with and without
// the cache. The risk is not the cache being wrong, it is it returning ANOTHER selection's
// answer — so different selections must stay different.
func TestFnSize_cachePerSelection(t *testing.T) {
	t.Run("DCLND-B04: A selection's size counts units, rules, lines and scenarios, and is remembered per selection", func(t *testing.T) {})

	// What is counted.
	fr, fg := projetoComFeature(t)
	fc, err := New(fr, fg)
	if err != nil {
		t.Fatal(err)
	}
	got, err := fc.fnSize("layer=infra")
	if err != nil {
		t.Fatal(err)
	}
	want := Size{Units: 1, Rules: 3, Lines: strings.Count(specDeExemplo, "\n"), Scenes: 2}
	if got != want {
		t.Errorf("size = %+v, want %+v", got, want)
	}

	root, g := projectTwoSelections(t)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	a1, err := c.fnSize("layer=alfa")
	if err != nil {
		t.Fatal(err)
	}
	b1, err := c.fnSize("layer=beta")
	if err != nil {
		t.Fatal(err)
	}
	if a1.Units == b1.Units {
		t.Fatalf("the fixture does not tell the selections apart: alfa=%d beta=%d", a1.Units, b1.Units)
	}
	// Second round: now from the cache, and it must match the first.
	a2, _ := c.fnSize("layer=alfa")
	b2, _ := c.fnSize("layer=beta")
	if a2 != a1 {
		t.Errorf("selection alfa changed when read from the cache: %+v vs %+v", a2, a1)
	}
	if b2 != b1 {
		t.Errorf("selection beta changed when read from the cache: %+v vs %+v", b2, b1)
	}
}

// projectTwoSelections builds two layers with DIFFERENT counts, so a cache that mixed the
// selections up would show.
func projectTwoSelections(t *testing.T) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	spec := func(code, layer string) string {
		return "---\ncode: " + code + "\nlayer: " + layer + "\n---\n\n# T — t\n\n## Visão Geral\n\nx.\n"
	}
	os.WriteFile(filepath.Join(root, "pkg/A.spec.md"), []byte(spec("AAAAA", "alfa")), 0o644)
	os.WriteFile(filepath.Join(root, "pkg/B.spec.md"), []byte(spec("BBBBB", "beta")), 0o644)
	os.WriteFile(filepath.Join(root, "pkg/C.spec.md"), []byte(spec("CCCCC", "beta")), 0o644)
	return root, &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/A.spec.md", Kind: mapx.KindSpec, Code: "AAAAA", Layer: "alfa"},
		{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB", Layer: "beta"},
		{ID: "pkg/C.spec.md", Kind: mapx.KindSpec, Code: "CCCCC", Layer: "beta"},
	}}
}

func TestFnSize_emptySelection(t *testing.T) {
	t.Run("DCLND-B05: A selection that matches no spec has size zero and is not big", func(t *testing.T) {
		c := compilerWithLayout(t, false)
		c.Layout = Layout{MaxUnits: 0, MaxLines: 0} // any measured unit would be big
		got, err := c.fnSize("layer=nothere")
		if err != nil || got != (Size{}) {
			t.Errorf("size = %+v, %v; want zero and no error", got, err)
		}
		big, err := c.fnBig("layer=nothere")
		if err != nil || big {
			t.Errorf("big = %v, %v; want false and no error", big, err)
		}
	})
	t.Run("DCLND-X01: The bigness asked by a template has no threshold of its own", func(t *testing.T) {
		c := compilerWithLayout(t, false)
		if big, _ := c.fnBig("layer=infra"); big {
			t.Fatal("one unit is not big under the default layout")
		}
		c.Layout = Layout{MaxUnits: 0, MaxLines: 0}
		if big, _ := c.fnBig("layer=infra"); !big {
			t.Error("the template's answer did not follow the compiler's layout")
		}
	})
}

// The diagram's ARROWS come from the map, aggregated by layer.
func TestLayerDeps_aggregatesByLayer(t *testing.T) {
	t.Run("DCLND-B06: Arrows between layers aggregate dependency edges between specs and carry their count", func(t *testing.T) {})
	specs := map[string]string{
		"a/S1.spec.md": "---\ncode: SCRNA\nlayer: screen\n---\n\n# S1 — tela um\n",
		"a/S2.spec.md": "---\ncode: SCRNB\nlayer: screen\n---\n\n# S2 — tela dois\n",
		"b/H1.spec.md": "---\ncode: SHRDA\nlayer: shared\n---\n\n# H1 — util\n",
	}
	root, g := projetoDeTeste(t, specs)
	// Two screens depend on the same util: one arrow, weight 2.
	g.Edges = []mapx.Edge{
		{From: "a/S1.spec.md", To: "b/H1.spec.md", Type: "needs"},
		{From: "a/S2.spec.md", To: "b/H1.spec.md", Type: "needs"},
		// An edge leaving the specs is not a layer arrow: a plan is not a layer, and drawing
		// it would produce an arrow from a node the page does not show.
		{From: "plans/0001.md", To: "a/S1.spec.md", Type: "needs"},
		// Nor is an edge of another type: `tested-by` is the unit, not a dependency.
		{From: "a/S1.spec.md", To: "b/H1.spec.md", Type: "tested-by"},
		// Nor one inside a layer: a loop on one node says nothing about the architecture.
		{From: "a/S1.spec.md", To: "a/S2.spec.md", Type: "depends-on"},
	}
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}

	deps := c.fnLayerDeps()
	if len(deps) != 1 {
		t.Fatalf("deps = %+v, want 1 (screen→shared)", deps)
	}
	if deps[0].From != "screen" || deps[0].To != "shared" {
		t.Errorf("arrow = %s→%s", deps[0].From, deps[0].To)
	}
	// The WEIGHT tells a single dependency from a systemic one — without it both would look
	// the same in the diagram.
	if deps[0].Count != 2 {
		t.Errorf("weight = %d, want 2 (two screens hold the arrow)", deps[0].Count)
	}
}

// A dependency INSIDE one layer is not an arrow.
func TestLayerDeps_ignoresInternalDependency(t *testing.T) {
	specs := map[string]string{
		"a/S1.spec.md": "---\ncode: SCRNA\nlayer: screen\n---\n\n# S1 — um\n",
		"a/S2.spec.md": "---\ncode: SCRNB\nlayer: screen\n---\n\n# S2 — dois\n",
	}
	root, g := projetoDeTeste(t, specs)
	g.Edges = []mapx.Edge{{From: "a/S1.spec.md", To: "a/S2.spec.md", Type: "needs"}}
	c, _ := New(root, g)
	if deps := c.fnLayerDeps(); len(deps) != 0 {
		t.Errorf("deps = %+v — the dependency is internal to the layer", deps)
	}
}

func TestLayerDeps_order(t *testing.T) {
	t.Run("DCLND-B07: The arrows are ordered by source layer then target layer", func(t *testing.T) {
		specs := map[string]string{
			"x/A.spec.md": "---\ncode: AAAAA\nlayer: a\n---\n\n# A — a\n",
			"x/B.spec.md": "---\ncode: BBBBB\nlayer: b\n---\n\n# B — b\n",
			"x/C.spec.md": "---\ncode: CCCCC\nlayer: c\n---\n\n# C — c\n",
		}
		root, g := projetoDeTeste(t, specs)
		g.Edges = []mapx.Edge{
			{From: "x/B.spec.md", To: "x/A.spec.md", Type: "needs"},
			{From: "x/A.spec.md", To: "x/C.spec.md", Type: "depends-on"},
			{From: "x/A.spec.md", To: "x/B.spec.md", Type: "needs"},
		}
		c, err := New(root, g)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, d := range c.fnLayerDeps() {
			got = append(got, d.From+"→"+d.To)
		}
		if strings.Join(got, ",") != "a→b,a→c,b→a" {
			t.Errorf("order = %v, want a→b, a→c, b→a", got)
		}
	})
}

// A Mermaid ID cannot carry a hyphen or an accent: the parser breaks, and the WHOLE BLOCK
// stops rendering — it shows as raw text. One `feature-hook` layer would be enough.
func TestMermaidID(t *testing.T) {
	t.Run("DCLND-B08: A diagram node identifier has only ASCII letters, digits and underscores", func(t *testing.T) {})
	cases := map[string]string{
		"feature-hook": "n_feature_hook",
		"app-infra":    "n_app_infra",
		"lambdas":      "n_lambdas",
		"padrão":       "n_padr_o",
	}
	for in, want := range cases {
		if got := MermaidID(in); got != want {
			t.Errorf("MermaidID(%q) = %q, want %q", in, got, want)
		}
	}
}

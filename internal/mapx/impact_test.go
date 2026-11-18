// @anchors
//   code: IMTSA
//   ref: IMANM

package mapx

import (
	"reflect"
	"slices"
	"testing"
)

// graph: a guide governs the Login spec AND the Home spec (a shared ruler). The Login unit is
// complete. It proves that climbing from Login reaches the guide (validate), but does NOT go
// down from the guide to Home (climbing does not re-propagate).
func impactGraph() *Graph {
	return &Graph{
		Version: 1,
		Nodes: []Node{
			{ID: "Login.spec.md", Kind: KindSpec},
			{ID: "Login.tsx", Kind: KindCode},
			{ID: "Login.feature", Kind: KindFeature},
			{ID: "Login.test.tsx", Kind: KindTest},
			{ID: "SPEC_GUIDE.md", Kind: KindGuide},
			{ID: "Home.spec.md", Kind: KindSpec},
			{ID: "Home.tsx", Kind: KindCode},
		},
		Edges: []Edge{
			{From: "Login.spec.md", To: "Login.tsx", Type: EdgeSpecifies},
			{From: "Login.spec.md", To: "Login.feature", Type: EdgeCoveredBy},
			{From: "Login.feature", To: "Login.test.tsx", Type: EdgeTestedBy},
			{From: "SPEC_GUIDE.md", To: "Login.spec.md", Type: EdgeGoverns},
			{From: "SPEC_GUIDE.md", To: "Home.spec.md", Type: EdgeGoverns},
			{From: "Home.spec.md", To: "Home.tsx", Type: EdgeSpecifies},
		},
	}
}

// Changing the Login SPEC: propagates (down) to code+feature+test; validates (up) against the
// guide. It must NOT touch Home, neither down nor up.
func TestAnalyzeImpact_specChange(t *testing.T) {
	t.Run("IMANM-B01: Changing a spec propagates down to its code, feature and test", func(t *testing.T) {})
	t.Run("IMANM-B04: Climbing to a shared guide does not reach the sibling unit", func(t *testing.T) {})
	imp := impactGraph().AnalyzeImpact("Login.spec.md")

	wantProp := []string{"Login.feature", "Login.test.tsx", "Login.tsx"}
	if !slices.Equal(imp.Propagate, wantProp) {
		t.Errorf("Propagate = %v, want %v", imp.Propagate, wantProp)
	}
	wantVal := []string{"SPEC_GUIDE.md"}
	if !slices.Equal(imp.Validate, wantVal) {
		t.Errorf("Validate = %v, want %v", imp.Validate, wantVal)
	}
	// the sibling Home NEVER appears — climbing to the guide does not come back down to siblings
	if slices.Contains(imp.Propagate, "Home.spec.md") || slices.Contains(imp.Validate, "Home.spec.md") ||
		slices.Contains(imp.Propagate, "Home.tsx") || slices.Contains(imp.Validate, "Home.tsx") {
		t.Error("the sibling Home should not be reached (climbing does not re-propagate)")
	}
}

// Changing the CODE (a leaf of the unit): propagates to nobody (nothing below depends on the
// code); validates upward against the spec and the guide.
func TestAnalyzeImpact_codeChange(t *testing.T) {
	t.Run("IMANM-B03: Changing code is validated upward against its spec and the spec's guide", func(t *testing.T) {})
	imp := impactGraph().AnalyzeImpact("Login.tsx")
	if len(imp.Propagate) != 0 {
		t.Errorf("changing the code should not propagate down; got %v", imp.Propagate)
	}
	wantVal := []string{"Login.spec.md", "SPEC_GUIDE.md"}
	if !slices.Equal(imp.Validate, wantVal) {
		t.Errorf("Validate = %v, want %v", imp.Validate, wantVal)
	}
}

// Changing the GUIDE (a high-degree parent): propagates down to EVERY spec it governs — the
// global wave. That is the right behaviour when a ruler changes.
func TestAnalyzeImpact_guideChange(t *testing.T) {
	t.Run("IMANM-B01: Changing a spec propagates down to its code, feature and test", func(t *testing.T) {})
	imp := impactGraph().AnalyzeImpact("SPEC_GUIDE.md")
	for _, want := range []string{"Login.spec.md", "Home.spec.md", "Login.tsx", "Home.tsx"} {
		if !slices.Contains(imp.Propagate, want) {
			t.Errorf("changing the guide should propagate to %q (global wave); got %v", want, imp.Propagate)
		}
	}
}

// @noPropagation: a marked child does not let the wave descend THROUGH it.
func TestAnalyzeImpact_noPropagation(t *testing.T) {
	t.Run("IMANM-B02: A no-propagation child is reached but the wave stops there", func(t *testing.T) {})
	g := impactGraph()
	for i := range g.Nodes {
		if g.Nodes[i].ID == "Login.feature" {
			g.Nodes[i].NoPropagation = true
		}
	}
	imp := g.AnalyzeImpact("Login.spec.md")
	if !slices.Contains(imp.Propagate, "Login.feature") {
		t.Error("the feature (a direct child) must still be reached")
	}
	if slices.Contains(imp.Propagate, "Login.test.tsx") {
		t.Error("the test should NOT be reached — the @noPropagation feature does not let the wave through")
	}
}

func TestAnalyzeImpact_isolated(t *testing.T) {
	t.Run("IMANM-B05: A node with no edges has no impact", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "loose.md"}}}
	imp := g.AnalyzeImpact("loose.md")
	if len(imp.Propagate) != 0 || len(imp.Validate) != 0 {
		t.Errorf("an isolated node has no impact; got prop=%v val=%v", imp.Propagate, imp.Validate)
	}
}

func TestAnalyzeImpact_cycleNeverListsTheOrigin(t *testing.T) {
	t.Run("IMANM-I01: The changed node is never in its own lists, even in a cycle", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "A"}, {ID: "B"}},
		Edges: []Edge{{From: "A", To: "B", Type: EdgeDependsOn}, {From: "B", To: "A", Type: EdgeDependsOn}},
	}
	imp := g.AnalyzeImpact("A")
	if !slices.Equal(imp.Propagate, []string{"B"}) || !slices.Equal(imp.Validate, []string{"B"}) {
		t.Errorf("the cycle should reach only B in both directions; got prop=%v val=%v", imp.Propagate, imp.Validate)
	}
}

func TestAnalyzeImpact_leavesTheGraphUntouched(t *testing.T) {
	t.Run("IMANM-X01: The analysis leaves the graph as it was", func(t *testing.T) {})
	g := impactGraph()
	g.Nodes[2].NoPropagation = true
	before := impactGraph()
	before.Nodes[2].NoPropagation = true
	g.AnalyzeImpact("Login.spec.md")
	g.AnalyzeImpact("Login.tsx")
	if !reflect.DeepEqual(g, before) {
		t.Errorf("the analysis changed the graph:\nbefore %+v\nafter  %+v", before, g)
	}
}

func TestAnalyzeImpact_validationStopsAtTypesOnly(t *testing.T) {
	t.Run("IMANM-B10: A parent that imports only types is validated and not climbed past", func(t *testing.T) {})
	g := dependentsGraph()
	g.Edges = append(g.Edges, Edge{From: "deep.tsx", To: "other.tsx", Type: EdgeDependsOn},
		Edge{From: "other.spec.md", To: "other.tsx", Type: EdgeSpecifies})
	g.Nodes = append(g.Nodes, Node{ID: "deep.tsx", Kind: KindCode}, Node{ID: "other.spec.md", Kind: KindSpec})
	val := g.AnalyzeImpact("util.ts").Validate
	if !slices.Contains(val, "other.tsx") {
		t.Errorf("the importer of types is still validated: %v", val)
	}
	for _, past := range []string{"deep.tsx", "other.spec.md"} {
		if slices.Contains(val, past) {
			t.Errorf("%s is past an import of types and must not be reached: %v", past, val)
		}
	}
	if !slices.Contains(val, "app.tsx") {
		t.Errorf("a value import is still climbed: %v", val)
	}
}

// dependentsGraph is a screen imported by a navigator that every screen imports for its types,
// a test and a capture of the screen, its spec, and a guide.
func dependentsGraph() *Graph {
	return &Graph{
		Nodes: []Node{
			{ID: "screen.tsx", Kind: KindCode}, {ID: "nav.tsx", Kind: KindCode}, {ID: "other.tsx", Kind: KindCode},
			{ID: "app.tsx", Kind: KindCode}, {ID: "util.ts", Kind: KindCode},
			{ID: "screen.test.tsx", Kind: KindTest}, {ID: "flow.yaml", Kind: KindTest}, {ID: "nav.test.tsx", Kind: KindTest},
			{ID: "screen.spec.md", Kind: KindSpec}, {ID: "GUIDE.md", Kind: KindGuide},
		},
		Edges: []Edge{
			{From: "screen.tsx", To: "util.ts", Type: EdgeDependsOn},
			{From: "nav.tsx", To: "screen.tsx", Type: EdgeDependsOn},
			{From: "other.tsx", To: "nav.tsx", Type: EdgeDependsOn, NoSideEffect: true},
			{From: "app.tsx", To: "nav.tsx", Type: EdgeDependsOn},
			{From: "screen.test.tsx", To: "screen.tsx", Type: EdgeDependsOn},
			{From: "nav.test.tsx", To: "nav.tsx", Type: EdgeDependsOn},
			{From: "flow.yaml", To: "screen.tsx", Type: EdgeCaptures},
			{From: "screen.spec.md", To: "screen.tsx", Type: EdgeSpecifies},
			{From: "GUIDE.md", To: "screen.spec.md", Type: EdgeGoverns},
		},
	}
}

func ids(ds []Dependent) []string {
	var out []string
	for _, d := range ds {
		out = append(out, d.ID)
	}
	return out
}

func TestDependents_whoDependsOnTheFile(t *testing.T) {
	t.Run("IMANM-B06: The dependents are the side that depends, by level", func(t *testing.T) {})
	g := dependentsGraph()
	direct := g.Dependents([]string{"screen.tsx"}, 1, nil)
	if got := ids(direct); !reflect.DeepEqual(got, []string{"flow.yaml", "nav.tsx", "screen.test.tsx"}) {
		t.Errorf("direct: the importer, the test and the capture — not what the screen imports, nor its spec: %v", got)
	}
	indirect := g.Dependents([]string{"screen.tsx"}, 2, nil)
	if got := ids(indirect); !reflect.DeepEqual(got, []string{"flow.yaml", "nav.tsx", "screen.test.tsx", "app.tsx", "nav.test.tsx", "other.tsx"}) {
		t.Errorf("indirect adds what depends on the navigator: %v", got)
	}
	for _, d := range indirect {
		if d.ID == "app.tsx" && (d.Level != 2 || d.From != "nav.tsx") {
			t.Errorf("each item says its level and where it was reached from: %+v", d)
		}
	}
	if got := ids(g.Dependents([]string{"screen.spec.md"}, 1, nil)); !reflect.DeepEqual(got, []string{"screen.tsx"}) {
		t.Errorf("a spec's dependents are what it specifies: %v", got)
	}
}

func TestDependents_whereTheWalkStops(t *testing.T) {
	t.Run("IMANM-B07: A tie of types alone, a stop and @noPropagation are listed and not walked past", func(t *testing.T) {})
	g := dependentsGraph()
	g.Edges = append(g.Edges, Edge{From: "deep.tsx", To: "other.tsx", Type: EdgeDependsOn})
	g.Nodes = append(g.Nodes, Node{ID: "deep.tsx", Kind: KindCode})
	all := g.Dependents([]string{"screen.tsx"}, 0, nil)
	for _, d := range all {
		if d.ID == "other.tsx" && !d.TypesOnly {
			t.Errorf("an import of types is marked: %+v", d)
		}
		if d.ID == "deep.tsx" {
			t.Errorf("no change at run time crosses an import of types: %+v", d)
		}
	}
	stopped := g.Dependents([]string{"screen.tsx"}, 0, func(id string) bool { return id == "nav.tsx" })
	if got := ids(stopped); !reflect.DeepEqual(got, []string{"flow.yaml", "nav.tsx", "screen.test.tsx"}) {
		t.Errorf("the stop is listed and not climbed from: %v", got)
	}
	g.Nodes[1].NoPropagation = true
	if got := ids(g.Dependents([]string{"screen.tsx"}, 0, nil)); slices.Contains(got, "app.tsx") {
		t.Errorf("@noPropagation is not climbed from: %v", got)
	}
}

func TestDependencies_theOtherWay(t *testing.T) {
	t.Run("IMANM-B08: The dependencies are the other end of the same ties", func(t *testing.T) {})
	g := dependentsGraph()
	if got := ids(g.Dependencies([]string{"screen.tsx"}, 1, nil)); !reflect.DeepEqual(got, []string{"screen.spec.md", "util.ts"}) {
		t.Errorf("what the screen imports and the spec that rules it: %v", got)
	}
	if got := ids(g.Dependencies([]string{"screen.tsx"}, 0, nil)); !reflect.DeepEqual(got, []string{"screen.spec.md", "util.ts", "GUIDE.md"}) {
		t.Errorf("the cascade reaches the guide of the spec: %v", got)
	}
	if g.DependencyCount("screen.tsx") != 2 || g.DependentCount("nav.tsx") != 3 {
		t.Errorf("counts: %d dependencies, %d dependents", g.DependencyCount("screen.tsx"), g.DependentCount("nav.tsx"))
	}
}

func TestRulersAndTestsHolding(t *testing.T) {
	t.Run("IMANM-B09: The rulers are what governs the file, and the tests are those reached or holding it", func(t *testing.T) {})
	g := dependentsGraph()
	if got := g.Rulers("screen.tsx"); !reflect.DeepEqual(got, []string{"GUIDE.md", "screen.spec.md"}) {
		t.Errorf("the spec and its guide, nothing that depends on it: %v", got)
	}
	// The navigator's test reaches the screen through the navigator it imports, and the capture
	// holds what the captured screen uses that no capture of its own proves.
	if got := g.TestsHolding(map[string]bool{"screen.tsx": true}); !reflect.DeepEqual(got, []string{"flow.yaml", "nav.test.tsx", "screen.test.tsx"}) {
		t.Errorf("the tests whose evidence holds the screen: %v", got)
	}
	if got := g.TestsHolding(map[string]bool{"util.ts": true}); !reflect.DeepEqual(got, []string{"flow.yaml", "nav.test.tsx", "screen.test.tsx"}) {
		t.Errorf("the tests whose evidence holds what the screen imports: %v", got)
	}
	if got := g.TestsHolding(map[string]bool{"app.tsx": true, "nav.test.tsx": true}); !reflect.DeepEqual(got, []string{"nav.test.tsx"}) {
		t.Errorf("a test reached is a test to run: %v", got)
	}
}

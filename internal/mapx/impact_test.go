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

// @anchors
//   ref: GRQRG

package mapx

import (
	"reflect"
	"strings"
	"testing"
)

func queryGraph() *Graph {
	return &Graph{
		Nodes: []Node{
			{ID: "spec", Kind: KindSpec},
			{ID: "code", Kind: KindCode},
			{ID: "feature", Kind: KindFeature},
			{ID: "guide", Kind: KindGuide},
			{ID: "island", Kind: KindDoc}, // no edges
		},
		Edges: []Edge{
			{From: "guide", To: "spec", Type: EdgeGoverns},
			{From: "spec", To: "code", Type: EdgeSpecifies},
			{From: "spec", To: "feature", Type: EdgeCoveredBy},
		},
	}
}

func TestNeighbors(t *testing.T) {
	t.Run("GRQRG-B03: The neighbourhood is one step each way, sorted", func(t *testing.T) {})
	g := queryGraph()
	// stored out of order: the answer must still be sorted by (from, to, type)
	g.Edges[1], g.Edges[2] = g.Edges[2], g.Edges[1]
	nb := g.Neighbors("spec")
	// incoming: guide→spec (governs); outgoing: spec→code, spec→feature
	if len(nb.In) != 1 || nb.In[0].From != "guide" {
		t.Errorf("In of spec = %+v, expected only guide", nb.In)
	}
	if len(nb.Out) != 2 || nb.Out[0].To != "code" || nb.Out[1].To != "feature" {
		t.Errorf("Out of spec = %+v, expected code then feature", nb.Out)
	}
}

func TestNeighbors_leaf(t *testing.T) {
	t.Run("GRQRG-B03: The neighbourhood is one step each way, sorted", func(t *testing.T) {})
	nb := queryGraph().Neighbors("code")
	if len(nb.Out) != 0 {
		t.Error("code is a leaf, it should have no outgoing edge")
	}
	if len(nb.In) != 1 || nb.In[0].From != "spec" {
		t.Errorf("In of code = %+v, expected spec", nb.In)
	}
}

func TestOrphans(t *testing.T) {
	t.Run("GRQRG-B04: Orphans are the nodes with no edge, sorted", func(t *testing.T) {})
	orphs := queryGraph().Orphans()
	if len(orphs) != 1 || orphs[0].ID != "island" {
		t.Errorf("Orphans = %+v, expected only 'island'", orphs)
	}
	g := queryGraph()
	g.Nodes = append(g.Nodes, Node{ID: "zeta"}, Node{ID: "alpha"})
	orphs = g.Orphans()
	if len(orphs) != 3 || orphs[0].ID != "alpha" || orphs[1].ID != "island" || orphs[2].ID != "zeta" {
		t.Errorf("Orphans should be sorted by id, got %+v", orphs)
	}
}

func TestStatistics(t *testing.T) {
	t.Run("GRQRG-B05: Statistics count nodes and edges by kind and type", func(t *testing.T) {})
	s := queryGraph().Statistics()
	if s.Nodes != 5 || s.Edges != 3 {
		t.Errorf("stats nodes/edges = %d/%d, want 5/3", s.Nodes, s.Edges)
	}
	if s.NodesByKind[KindSpec] != 1 || s.NodesByKind[KindCode] != 1 {
		t.Errorf("wrong count by kind: %+v", s.NodesByKind)
	}
	if s.EdgesByType[EdgeGoverns] != 1 || s.EdgesByType[EdgeSpecifies] != 1 {
		t.Errorf("wrong count by type: %+v", s.EdgesByType)
	}
}

func TestGovernsAndSummary(t *testing.T) {
	t.Run("GRQRG-B01: A guide governs only the targets of its governs edges, sorted", func(t *testing.T) {})
	t.Run("GRQRG-B02: The governance summary counts governs edges only", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "guides/A.md", Kind: KindGuide},
			{ID: "x.tsx", Kind: KindCode}, {ID: "y.tsx", Kind: KindCode},
			{ID: "z.spec.md", Kind: KindSpec},
		},
		Edges: []Edge{
			{From: "guides/A.md", To: "y.tsx", Type: EdgeGoverns},
			{From: "guides/A.md", To: "x.tsx", Type: EdgeGoverns},
			{From: "z.spec.md", To: "x.tsx", Type: EdgeSpecifies}, // not governs
		},
	}
	gov := g.Governs("guides/A.md")
	if len(gov) != 2 {
		t.Fatalf("A.md should govern 2, got %d (%v)", len(gov), gov)
	}
	if gov[0] != "x.tsx" || gov[1] != "y.tsx" {
		t.Errorf("Governs should come sorted, got %v", gov)
	}
	if other := g.Governs("z.spec.md"); len(other) != 0 {
		t.Errorf("a specifies edge is not governance, got %v", other)
	}
	sum := g.GovernanceSummary()
	if sum["guides/A.md"] != 2 {
		t.Errorf("summary of A.md = %d, want 2", sum["guides/A.md"])
	}
	if _, ok := sum["z.spec.md"]; ok {
		t.Error("specifies should not count as governance")
	}
}

func TestTopoOrderParentsBeforeChildren(t *testing.T) {
	t.Run("GRQRG-B06: Parents come before their children", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "test", Kind: KindTest}, {ID: "code", Kind: KindCode},
			{ID: "feature", Kind: KindFeature}, {ID: "spec", Kind: KindSpec},
			{ID: "guide", Kind: KindGuide},
		},
		Edges: []Edge{
			{From: "guide", To: "spec", Type: EdgeGoverns},
			{From: "spec", To: "code", Type: EdgeSpecifies},
			{From: "spec", To: "feature", Type: EdgeCoveredBy},
			{From: "feature", To: "test", Type: EdgeTestedBy},
		},
	}
	order := g.TopoOrder()
	pos := map[string]int{}
	for i, n := range order {
		pos[n.ID] = i
	}
	// every parent comes BEFORE its child
	checks := [][2]string{{"guide", "spec"}, {"spec", "code"}, {"spec", "feature"}, {"feature", "test"}}
	for _, c := range checks {
		if pos[c[0]] >= pos[c[1]] {
			t.Errorf("%s (pos %d) should come before %s (pos %d)", c[0], pos[c[0]], c[1], pos[c[1]])
		}
	}
}

func TestTopoOrderHandlesCycle(t *testing.T) {
	t.Run("GRQRG-B08: Cycle members are appended at the end", func(t *testing.T) {})
	// cycle A→B→A: it does not stall, both are appended at the end.
	g := &Graph{
		Nodes: []Node{{ID: "B", Kind: KindCode}, {ID: "A", Kind: KindCode}, {ID: "C", Kind: KindSpec}},
		Edges: []Edge{
			{From: "A", To: "B", Type: EdgeSpecifies},
			{From: "B", To: "A", Type: EdgeSpecifies}, // cycle
		},
	}
	order := g.TopoOrder()
	if len(order) != 3 || order[0].ID != "C" || order[1].ID != "A" || order[2].ID != "B" {
		t.Fatalf("every node must appear, the cycle last and sorted: got %+v", order)
	}
}

func TestTopoOrderInformativeEdgesImposeNoOrder(t *testing.T) {
	t.Run("GRQRG-B07: A depends-on edge imposes no order", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "a"}, {ID: "b"}},
		Edges: []Edge{{From: "b", To: "a", Type: EdgeDependsOn}, {From: "b", To: "a", Type: EdgeReferences}},
	}
	order := g.TopoOrder()
	if order[0].ID != "a" || order[1].ID != "b" {
		t.Errorf("an informative edge must not reorder the nodes, got %s, %s", order[0].ID, order[1].ID)
	}
}

func TestTopoOrderIsStable(t *testing.T) {
	t.Run("GRQRG-I01: The order does not depend on how nodes are stored", func(t *testing.T) {})
	edges := []Edge{{From: "g", To: "s2", Type: EdgeGoverns}, {From: "g", To: "s1", Type: EdgeGoverns}}
	one := &Graph{Nodes: []Node{{ID: "s2"}, {ID: "g"}, {ID: "s1"}, {ID: "x"}}, Edges: edges}
	two := &Graph{Nodes: []Node{{ID: "x"}, {ID: "s1"}, {ID: "g"}, {ID: "s2"}}, Edges: edges}
	ids := func(ns []Node) []string {
		var out []string
		for _, n := range ns {
			out = append(out, n.ID)
		}
		return out
	}
	a, b := ids(one.TopoOrder()), ids(two.TopoOrder())
	if !reflect.DeepEqual(a, b) || !reflect.DeepEqual(a, []string{"g", "s1", "s2", "x"}) {
		t.Errorf("the order must be g, s1, s2, x whatever the storage: got %v and %v", a, b)
	}
}

func TestQueriesLeaveTheGraphUntouched(t *testing.T) {
	t.Run("GRQRG-X01: Queries leave the graph as it was", func(t *testing.T) {})
	g := queryGraph()
	g.Edges[1], g.Edges[2] = g.Edges[2], g.Edges[1]
	before := queryGraph()
	before.Edges[1], before.Edges[2] = before.Edges[2], before.Edges[1]
	g.Governs("guide")
	g.GovernanceSummary()
	g.Neighbors("spec")
	g.Orphans()
	g.Statistics()
	g.TopoOrder()
	if !reflect.DeepEqual(g, before) {
		t.Errorf("a query changed the graph:\nbefore %+v\nafter  %+v", before, g)
	}
}

func TestUnitCodesOf(t *testing.T) {
	t.Run("GRQRG-B09: A node's unit codes come from its identity edges", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "pay.spec.md", Kind: KindSpec, Code: "PAYMT"},
			{ID: "share.spec.md", Kind: KindSpec, Code: "SHARE"},
			{ID: "pay.feature", Kind: KindFeature}, {ID: "pay.go", Kind: KindCode}, {ID: "pay_test.go", Kind: KindTest}, {ID: "lone.go", Kind: KindCode},
		},
		Edges: []Edge{
			{From: "pay.spec.md", To: "pay.go", Type: EdgeSpecifies},
			{From: "share.spec.md", To: "pay.go", Type: EdgeSpecifies},
			{From: "pay.spec.md", To: "pay.feature", Type: EdgeCoveredBy},
			{From: "pay.feature", To: "pay_test.go", Type: EdgeTestedBy},
		},
	}
	if got := strings.Join(g.UnitCodesOf("pay.go"), ","); got != "PAYMT,SHARE" {
		t.Errorf("pay.go: %s", got)
	}
	if got := strings.Join(g.UnitCodesOf("pay_test.go"), ","); got != "PAYMT" {
		t.Errorf("pay_test.go: %s", got)
	}
	if got := g.UnitCodesOf("lone.go"); len(got) != 0 {
		t.Errorf("lone.go: %v", got)
	}
}

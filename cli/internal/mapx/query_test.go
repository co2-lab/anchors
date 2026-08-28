package mapx

import "testing"

func queryGraph() *Graph {
	return &Graph{
		Nodes: []Node{
			{ID: "spec", Kind: KindSpec},
			{ID: "code", Kind: KindCode},
			{ID: "feature", Kind: KindFeature},
			{ID: "guide", Kind: KindGuide},
			{ID: "ilha", Kind: KindDoc}, // sem arestas
		},
		Edges: []Edge{
			{From: "guide", To: "spec", Type: EdgeGoverns},
			{From: "spec", To: "code", Type: EdgeSpecifies},
			{From: "spec", To: "feature", Type: EdgeCoveredBy},
		},
	}
}

func TestNeighbors(t *testing.T) {
	nb := queryGraph().Neighbors("spec")
	// entrada: guide→spec (governs); saída: spec→code, spec→feature
	if len(nb.In) != 1 || nb.In[0].From != "guide" {
		t.Errorf("In de spec = %+v, esperava só guide", nb.In)
	}
	if len(nb.Out) != 2 {
		t.Errorf("Out de spec = %+v, esperava 2 (code, feature)", nb.Out)
	}
}

func TestNeighbors_folha(t *testing.T) {
	nb := queryGraph().Neighbors("code")
	if len(nb.Out) != 0 {
		t.Error("code é folha, não deveria ter saída")
	}
	if len(nb.In) != 1 || nb.In[0].From != "spec" {
		t.Errorf("In de code = %+v, esperava spec", nb.In)
	}
}

func TestOrphans(t *testing.T) {
	orphs := queryGraph().Orphans()
	if len(orphs) != 1 || orphs[0].ID != "ilha" {
		t.Errorf("Orphans = %+v, esperava só 'ilha'", orphs)
	}
}

func TestStatistics(t *testing.T) {
	s := queryGraph().Statistics()
	if s.Nodes != 5 || s.Edges != 3 {
		t.Errorf("stats nodes/edges = %d/%d, want 5/3", s.Nodes, s.Edges)
	}
	if s.NodesByKind[KindSpec] != 1 || s.NodesByKind[KindCode] != 1 {
		t.Errorf("contagem por kind errada: %+v", s.NodesByKind)
	}
	if s.EdgesByType[EdgeGoverns] != 1 || s.EdgesByType[EdgeSpecifies] != 1 {
		t.Errorf("contagem por tipo errada: %+v", s.EdgesByType)
	}
}

func TestGovernsAndSummary(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{ID: "guides/A.md", Kind: KindGuide},
			{ID: "x.tsx", Kind: KindCode}, {ID: "y.tsx", Kind: KindCode},
			{ID: "z.spec.md", Kind: KindSpec},
		},
		Edges: []Edge{
			{From: "guides/A.md", To: "x.tsx", Type: EdgeGoverns},
			{From: "guides/A.md", To: "y.tsx", Type: EdgeGoverns},
			{From: "z.spec.md", To: "x.tsx", Type: EdgeSpecifies}, // não é governs
		},
	}
	reg := g.Governs("guides/A.md")
	if len(reg) != 2 {
		t.Fatalf("A.md deveria reger 2, veio %d (%v)", len(reg), reg)
	}
	// ordenado
	if reg[0] != "x.tsx" || reg[1] != "y.tsx" {
		t.Errorf("Governs deveria vir ordenado, veio %v", reg)
	}
	sum := g.GovernanceSummary()
	if sum["guides/A.md"] != 2 {
		t.Errorf("summary de A.md = %d, quer 2", sum["guides/A.md"])
	}
	if _, ok := sum["z.spec.md"]; ok {
		t.Error("specifies não deveria contar como governança")
	}
}

func TestTopoOrderParentsBeforeChildren(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{ID: "test", Kind: KindTest}, {ID: "code", Kind: KindCode},
			{ID: "feature", Kind: KindFeature}, {ID: "spec", Kind: KindSpec},
			{ID: "guide", Kind: KindGuide},
		},
		Edges: []Edge{
			{From: "guide", To: "spec", Type: EdgeGoverns},  // guide rege spec
			{From: "spec", To: "code", Type: EdgeSpecifies}, // spec rege code
			{From: "spec", To: "feature", Type: EdgeCoveredBy},
			{From: "feature", To: "test", Type: EdgeTestedBy}, // feature rege test
		},
	}
	order := g.TopoOrder()
	pos := map[string]int{}
	for i, n := range order {
		pos[n.ID] = i
	}
	// cada pai vem ANTES do filho
	checks := [][2]string{{"guide", "spec"}, {"spec", "code"}, {"spec", "feature"}, {"feature", "test"}}
	for _, c := range checks {
		if pos[c[0]] >= pos[c[1]] {
			t.Errorf("%s (pos %d) deveria vir antes de %s (pos %d)", c[0], pos[c[0]], c[1], pos[c[1]])
		}
	}
}

func TestTopoOrderHandlesCycle(t *testing.T) {
	// ciclo A→B→A: não trava, anexa ambos ao fim.
	g := &Graph{
		Nodes: []Node{{ID: "A", Kind: KindCode}, {ID: "B", Kind: KindCode}, {ID: "C", Kind: KindSpec}},
		Edges: []Edge{
			{From: "A", To: "B", Type: EdgeSpecifies},
			{From: "B", To: "A", Type: EdgeSpecifies}, // ciclo
		},
	}
	order := g.TopoOrder()
	if len(order) != 3 {
		t.Fatalf("todos os nós devem aparecer mesmo com ciclo, veio %d", len(order))
	}
}

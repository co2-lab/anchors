package mapx

import "testing"

// grafo: spec(rev a) ──specifies──▶ code(rev b), e code ──tested-by──▶ test(rev c)
func stampGraph() *Graph {
	return &Graph{
		Version: 1,
		Nodes: []Node{
			{ID: "A.spec.md", Kind: KindSpec, Rev: "a"},
			{ID: "A.tsx", Kind: KindCode, Rev: "b"},
			{ID: "A.test.tsx", Kind: KindTest, Rev: "c"},
		},
		Edges: []Edge{
			{From: "A.spec.md", To: "A.tsx", Type: EdgeSpecifies},
			{From: "A.tsx", To: "A.test.tsx", Type: EdgeTestedBy},
		},
	}
}

func TestStampOnlyEdgesWithBothEndsConfronted(t *testing.T) {
	g := stampGraph()
	// só spec e code foram confrontados; test não entrou nesta rodada
	n := g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "2026-08-07T00:00:00Z")
	if n != 1 {
		t.Fatalf("esperava carimbar 1 aresta (spec→code), veio %d", n)
	}
	if g.Edges[0].Stamp == nil {
		t.Fatal("aresta spec→code deveria ter carimbo")
	}
	if g.Edges[1].Stamp != nil {
		t.Fatal("aresta code→test NÃO deveria ter carimbo (test não confrontado)")
	}
}

func TestStampVerdictIssueWhenEndFails(t *testing.T) {
	g := stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md", Failed: true}, {ID: "A.tsx"}}, "t")
	if got := g.Edges[0].Stamp.Verdict; got != "issue" {
		t.Fatalf("ponta falha → verdict issue, veio %q", got)
	}
}

func TestStampVerdictOkWhenBothPass(t *testing.T) {
	g := stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	if got := g.Edges[0].Stamp.Verdict; got != "ok" {
		t.Fatalf("ambos passam → verdict ok, veio %q", got)
	}
}

func TestStaleBecomesValidatedAfterStamp(t *testing.T) {
	g := stampGraph()
	// antes: nunca validada → stale
	if !g.Stale(g.Edges[0]) {
		t.Fatal("aresta nunca validada deveria ser stale")
	}
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	// depois do carimbo com as revs atuais → não stale
	if g.Stale(g.Edges[0]) {
		t.Fatal("aresta recém-carimbada não deveria ser stale")
	}
}

func TestStaleReappearsWhenRevAdvances(t *testing.T) {
	g := stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	// a spec é editada → rev avança → a aresta volta a ficar stale
	g.Nodes[0].Rev = "a2"
	if !g.Stale(g.Edges[0]) {
		t.Fatal("após a ponta avançar de rev, a aresta deveria voltar a ser stale")
	}
}

func TestStampEdgeSingle(t *testing.T) {
	g := stampGraph()
	// carimba só a aresta spec→code com verdict issue (julgamento por IA)
	ok := g.StampEdge("A.spec.md", "A.tsx", "issue", "2026-08-07T00:00:00Z")
	if !ok {
		t.Fatal("StampEdge deveria achar a aresta A.spec.md→A.tsx")
	}
	if g.Edges[0].Stamp == nil || g.Edges[0].Stamp.Verdict != "issue" {
		t.Fatalf("aresta deveria ter carimbo issue, veio %+v", g.Edges[0].Stamp)
	}
	// a outra aresta permanece sem carimbo
	if g.Edges[1].Stamp != nil {
		t.Fatal("só a aresta pedida deveria ser carimbada")
	}
	// aresta inexistente → false
	if g.StampEdge("X", "Y", "ok", "t") {
		t.Fatal("StampEdge de aresta inexistente deveria devolver false")
	}
}

func TestStampEdgeGoesStaleOnRevChange(t *testing.T) {
	g := stampGraph()
	g.StampEdge("A.spec.md", "A.tsx", "ok", "t")
	if g.Stale(g.Edges[0]) {
		t.Fatal("recém-carimbada não deveria ser stale")
	}
	g.Nodes[1].Rev = "b2" // o alvo (código) mudou
	if !g.Stale(g.Edges[0]) {
		t.Fatal("o veredito de IA deveria envelhecer quando o alvo muda")
	}
}

func TestStaleEdgesListing(t *testing.T) {
	g := stampGraph()
	if len(g.StaleEdges()) != 2 {
		t.Fatalf("no início as 2 arestas são stale, veio %d", len(g.StaleEdges()))
	}
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	if got := len(g.StaleEdges()); got != 1 {
		t.Fatalf("após carimbar 1, resta 1 stale, veio %d", got)
	}
}

// O CARIMBO NÃO PODE MUDAR SOZINHO.
//
// `last_validated` responde "quando esta relação foi confrontada" — pergunta de auditoria
// — e NÃO entra na regra de staleness, que compara `rev` (PROPAGATION.md §3). Com
// precisão de segundo, cada `anchors check` reescrevia as 26 linhas de carimbo do mapa:
// conflito em todo PR onde duas pessoas rodaram o check, e um diff que muda sozinho sem
// dizer nada.
//
// Duas validações no MESMO DIA, com o mesmo veredito e as mesmas revs, têm de produzir o
// mesmo mapa. O que muda o carimbo é o veredito ou a rev — que é quando há o que registrar.
func TestCarimboIgualNaoReescreveOMapa(t *testing.T) {
	monta := func() *Graph {
		return &Graph{
			Nodes: []Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1"}},
			Edges: []Edge{{From: "a", To: "b", Type: EdgeSpecifies}},
		}
	}
	g1, g2 := monta(), monta()
	v := []NodeVerdict{{ID: "a"}, {ID: "b"}}

	// Mesmo DIA, instantes diferentes — é o caso real: dois `check` seguidos.
	g1.StampEdges(v, "2026-08-30")
	g2.StampEdges(v, "2026-08-30")

	if g1.Edges[0].Stamp.LastValidated != g2.Edges[0].Stamp.LastValidated {
		t.Fatalf("mesmo dia deveria produzir o mesmo carimbo: %q vs %q",
			g1.Edges[0].Stamp.LastValidated, g2.Edges[0].Stamp.LastValidated)
	}
	// E o campo não pode ficar vazio: perder a resposta seria trocar um problema por outro.
	if g1.Edges[0].Stamp.LastValidated == "" {
		t.Error("o carimbo ainda tem de dizer QUANDO — reduzir a precisão não é apagar")
	}
}

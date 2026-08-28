package mapx

import "sort"

// Consultas puras sobre o grafo — insumo do `map show`. Não tocam disco; operam só
// sobre o Graph carregado.

// Neighborhood é a vizinhança direta de um nó: as arestas de entrada (quem o rege /
// contra quem ele valida — sobe) e de saída (o que ele propaga — desce).
type Neighborhood struct {
	Node string
	In   []Edge // arestas onde o nó é `to` (os pais)
	Out  []Edge // arestas onde o nó é `from` (os filhos)
}

// Governs devolve os IDs dos nós que um guide rege DIRETAMENTE (as arestas governs
// que saem dele) — não a onda transitiva. É o que dimensiona uma auditoria de
// julgamento "por guide": cada um destes é um alvo a confrontar contra a régua.
func (g *Graph) Governs(guideID string) []string {
	var out []string
	for _, e := range g.Edges {
		if e.Type == EdgeGoverns && e.From == guideID {
			out = append(out, e.To)
		}
	}
	sort.Strings(out)
	return out
}

// GovernanceSummary devolve, para cada guide que rege algo, quantos nós ele rege
// diretamente. Ordenado por contagem desc (via GovernanceCounts para saída estável).
func (g *Graph) GovernanceSummary() map[string]int {
	counts := map[string]int{}
	for _, e := range g.Edges {
		if e.Type == EdgeGoverns {
			counts[e.From]++
		}
	}
	return counts
}

// Neighbors devolve a vizinhança direta de um nó (um passo em cada direção).
func (g *Graph) Neighbors(id string) Neighborhood {
	n := Neighborhood{Node: id}
	for _, e := range g.Edges {
		if e.To == id {
			n.In = append(n.In, e)
		}
		if e.From == id {
			n.Out = append(n.Out, e)
		}
	}
	sortEdges(n.In)
	sortEdges(n.Out)
	return n
}

// Orphans devolve os nós sem NENHUMA aresta (ilhas) — desconectados do organismo.
// É a forma estrutural do órfão (TRACEABILITY §6): a peça que não se liga a nada.
func (g *Graph) Orphans() []Node {
	connected := map[string]bool{}
	for _, e := range g.Edges {
		connected[e.From] = true
		connected[e.To] = true
	}
	var out []Node
	for _, n := range g.Nodes {
		if !connected[n.ID] {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Stats é o resumo do grafo.
type Stats struct {
	Nodes       int
	Edges       int
	NodesByKind map[Kind]int
	EdgesByType map[EdgeType]int
}

// Statistics agrega as contagens do grafo.
func (g *Graph) Statistics() Stats {
	s := Stats{Nodes: len(g.Nodes), Edges: len(g.Edges), NodesByKind: map[Kind]int{}, EdgesByType: map[EdgeType]int{}}
	for _, n := range g.Nodes {
		s.NodesByKind[n.Kind]++
	}
	for _, e := range g.Edges {
		s.EdgesByType[e.Type]++
	}
	return s
}

func sortEdges(es []Edge) {
	sort.Slice(es, func(i, j int) bool {
		if es[i].From != es[j].From {
			return es[i].From < es[j].From
		}
		if es[i].To != es[j].To {
			return es[i].To < es[j].To
		}
		return es[i].Type < es[j].Type
	})
}

// TopoOrder devolve os nós ordenados TOPOLOGICAMENTE de cima para baixo na árvore de
// dependência: as RÉGUAS/PAIS primeiro (guides, specs — quem rege), os REGIDOS depois
// (código, features, testes). Assim quem processa em ordem nunca conserta um filho
// antes do pai que o afeta (evita retrabalho — mudar o filho, depois o pai o invalida).
//
// A direção das arestas que descem (governs/specifies/covered-by/tested-by) define
// pai→filho: uma aresta from→to significa "from rege/precede to", então `to` depende
// de `from` e vem DEPOIS. Ordenação de Kahn; ciclos (raros) são anexados ao fim numa
// ordem estável, para nunca travar.
func (g *Graph) TopoOrder() []Node {
	// só as arestas de DEPENDÊNCIA vertical/co-location contam para a ordem
	// (references é informativa e não impõe precedência).
	depEdge := func(t EdgeType) bool {
		switch t {
		case EdgeGoverns, EdgeSpecifies, EdgeCoveredBy, EdgeTestedBy:
			return true
		}
		return false
	}
	// indegree: quantos pais cada nó tem (arestas de dependência chegando).
	indeg := map[string]int{}
	children := map[string][]string{}
	nodeByID := map[string]Node{}
	for _, n := range g.Nodes {
		indeg[n.ID] = 0
		nodeByID[n.ID] = n
	}
	seen := map[string]bool{}
	for _, e := range g.Edges {
		if !depEdge(e.Type) {
			continue
		}
		key := e.From + "\x00" + e.To
		if seen[key] {
			continue // não conta arestas duplicadas
		}
		seen[key] = true
		if _, ok := indeg[e.To]; ok {
			indeg[e.To]++
			children[e.From] = append(children[e.From], e.To)
		}
	}
	// fila dos sem-pai (o topo), em ordem estável
	var ready []string
	for id, d := range indeg {
		if d == 0 {
			ready = append(ready, id)
		}
	}
	sort.Strings(ready)

	var out []Node
	placed := map[string]bool{}
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		if placed[id] {
			continue
		}
		placed[id] = true
		out = append(out, nodeByID[id])
		kids := append([]string(nil), children[id]...)
		sort.Strings(kids)
		var freed []string
		for _, c := range kids {
			indeg[c]--
			if indeg[c] == 0 {
				freed = append(freed, c)
			}
		}
		ready = append(ready, freed...)
		sort.Strings(ready)
	}
	// nós em ciclo (indegree nunca zerou) — anexa ao fim, ordem estável.
	if len(out) < len(g.Nodes) {
		var rest []Node
		for _, n := range g.Nodes {
			if !placed[n.ID] {
				rest = append(rest, n)
			}
		}
		sort.Slice(rest, func(i, j int) bool { return rest[i].ID < rest[j].ID })
		out = append(out, rest...)
	}
	return out
}

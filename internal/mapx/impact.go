// @anchors
//   code: IMMPM
//   ref: IMANM

package mapx

import "sort"

// Impact é o resultado da análise de impacto (PROPAGATION §3) — CONSULTA pura.
// São DUAS dimensões com propósitos distintos (decisão de doutrina):
//
//   - Propaga  (descer, pai→filho): os filhos que DEPENDEM do nó alterado e
//     precisam ser refeitos. A onda de fato. Para em nós `@noPropagation`.
//   - Valida   (subir, filho→pai): os pais contra os quais o nó deve ser
//     CONFRONTADO — não é propagação, é validação (uma divergência aqui vira
//     issue, conciliável no filho OU no pai; CONCEPT §2/§5). Sobe, mas NÃO
//     re-propaga (para no pai — não pega irmãos).
//
// O impact NÃO abre issue nem muta nada — só mostra o caminho. Abrir a issue é do
// gate; executar a onda é do `propagate`.
type Impact struct {
	Origin    string   // o nó alterado (ponto de partida)
	Propagate []string // filhos alcançados descendo (a onda), ordenados
	Validate  []string // pais a confrontar subindo, ordenados
}

type adjacency struct {
	out map[string][]Edge // from → arestas que saem (para os filhos)
	in  map[string][]Edge // to   → arestas que entram (dos pais)
}

func (g *Graph) adjacency() adjacency {
	a := adjacency{out: map[string][]Edge{}, in: map[string][]Edge{}}
	for _, e := range g.Edges {
		a.out[e.From] = append(a.out[e.From], e)
		a.in[e.To] = append(a.in[e.To], e)
	}
	return a
}

func (g *Graph) noPropSet() map[string]bool {
	m := map[string]bool{}
	for _, n := range g.Nodes {
		if n.NoPropagation {
			m[n.ID] = true
		}
	}
	return m
}

// AnalyzeImpact computa as duas dimensões do impacto de alterar `start`.
func (g *Graph) AnalyzeImpact(start string) Impact {
	adj := g.adjacency()
	noProp := g.noPropSet()

	return Impact{
		Origin:    start,
		Propagate: g.propagateDown(start, adj, noProp),
		Validate:  g.validateUp(start, adj),
	}
}

// propagateDown: BFS descendo (from→to). O filho depende do pai → propaga. Poda em
// nós `@noPropagation` (não descem — o filho declarou que não depende do pai). É a
// "propagação podada" da doutrina.
func (g *Graph) propagateDown(start string, adj adjacency, noProp map[string]bool) []string {
	visited := map[string]bool{start: true}
	var reached []string
	queue := []string{start}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, e := range adj.out[node] {
			child := e.To
			if visited[child] {
				continue
			}
			visited[child] = true
			reached = append(reached, child)
			// o filho propaga adiante SÓ se ele mesmo não for @noPropagation
			if !noProp[child] {
				queue = append(queue, child)
			}
		}
	}
	sort.Strings(reached)
	return reached
}

// validateUp: sobe UM passo de cada vez (filho→pai) para confrontar, mas NÃO
// re-propaga a partir do pai (subir é validação, não propagação — para no pai, não
// pega irmãos). Segue transitivamente para cima (o pai também é validado contra o
// avô), mas nunca vira descida.
func (g *Graph) validateUp(start string, adj adjacency) []string {
	visited := map[string]bool{start: true}
	var reached []string
	queue := []string{start}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, e := range adj.in[node] {
			parent := e.From
			if visited[parent] {
				continue
			}
			visited[parent] = true
			reached = append(reached, parent)
			// A parent that imports only types is confronted, and the climb stops there: it
			// cannot change with the file below it. Without this a sentence edited in the legal
			// text reached, through the navigator every screen imports for its param types,
			// every screen of the app and its spec (measured on a project: 195 nodes, 63 specs).
			if e.NoSideEffect {
				continue
			}
			queue = append(queue, parent) // sobe mais (valida contra o avô), mas só SOBE
		}
	}
	sort.Strings(reached)
	return reached
}

// WHO DEPENDS ON A FILE — the impact the `impact` command answers (DESIGN in MPCTI).
//
// AnalyzeImpact descends every outgoing edge, and a `depends-on` goes from the file that
// imports to the file imported: from a screen it reached what the screen USES, and through the
// navigator the screen imports, every screen of the app (measured on a project: 833 "to redo" for a
// change to one screen, whose evidence closure holds five tests). Each edge type says which
// end depends on the other; the walk below always goes to the end that depends.

// dependsAtTo are the edge types whose `to` depends on its `from`: what governs, specifies,
// is covered, tested or seeded rules what it points at.
var dependsAtTo = map[EdgeType]bool{EdgeGoverns: true, EdgeSpecifies: true, EdgeCoveredBy: true, EdgeTestedBy: true, EdgeSeeds: true}

// dependsAtFrom are the edge types whose `from` depends on its `to`: the importer, the reader,
// the capture, the realizer, the gated rule, the plan that needs another.
var dependsAtFrom = map[EdgeType]bool{EdgeDependsOn: true, EdgeReadsFrom: true, EdgeCaptures: true, EdgeRealizes: true, EdgeGatedBy: true, EdgeNeeds: true}

// Dependent is a node reached from the files asked about: its level (1 is direct), the node it
// was reached from, and the edge type that ties the two.
type Dependent struct {
	ID    string
	Level int
	From  string
	Type  EdgeType
	// TypesOnly is a tie with no side effect on the tests, as the project's flags say — an
	// import of types: listed, and not walked past, since no change at run time crosses it.
	TypesOnly bool
}

// directDependents are the nodes that depend on id, each with the edge that ties them, as an
// edge from id to the dependent.
func (g *Graph) directDependents(id string, adj adjacency) []Edge {
	var out []Edge
	for _, e := range adj.out[id] {
		if dependsAtTo[e.Type] {
			out = append(out, Edge{From: id, To: e.To, Type: e.Type, NoSideEffect: e.NoSideEffect})
		}
	}
	for _, e := range adj.in[id] {
		if dependsAtFrom[e.Type] {
			out = append(out, Edge{From: id, To: e.From, Type: e.Type, NoSideEffect: e.NoSideEffect})
		}
	}
	return out
}

// directDependencies are the nodes id depends on — the other end of the same ties —, as an
// edge from id to the dependency.
func (g *Graph) directDependencies(id string, adj adjacency) []Edge {
	var out []Edge
	for _, e := range adj.out[id] {
		if dependsAtFrom[e.Type] {
			out = append(out, Edge{From: id, To: e.To, Type: e.Type, NoSideEffect: e.NoSideEffect})
		}
	}
	for _, e := range adj.in[id] {
		if dependsAtTo[e.Type] {
			out = append(out, Edge{From: id, To: e.From, Type: e.Type, NoSideEffect: e.NoSideEffect})
		}
	}
	return out
}

// DependentCount is how many nodes depend directly on id.
func (g *Graph) DependentCount(id string) int {
	return len(g.directDependents(id, g.evidenceIndex().adj))
}

// DependencyCount is how many nodes id depends on directly.
func (g *Graph) DependencyCount(id string) int {
	return len(g.directDependencies(id, g.evidenceIndex().adj))
}

// Dependents are the nodes that depend on the starts, up to maxLevel levels (0 has no limit),
// each at the level it is first reached, sorted by level and id. A node marked
// `@noPropagation`, one `stop` names, or one reached by a tie of types alone, is listed and not
// climbed from. A start is never among them.
func (g *Graph) Dependents(starts []string, maxLevel int, stop func(string) bool) []Dependent {
	return g.walkTies(starts, maxLevel, stop, g.directDependents, true)
}

// Dependencies are the nodes the starts depend on, the same walk the other way: what they
// import, read, capture, realize, and the specs, features and guides that rule them. A node
// `stop` names is listed and not descended from.
func (g *Graph) Dependencies(starts []string, maxLevel int, stop func(string) bool) []Dependent {
	return g.walkTies(starts, maxLevel, stop, g.directDependencies, false)
}

func (g *Graph) walkTies(starts []string, maxLevel int, stop func(string) bool, step func(string, adjacency) []Edge, prune bool) []Dependent {
	adj := g.evidenceIndex().adj
	noProp := map[string]bool{}
	if prune {
		noProp = g.noPropSet()
	}
	seen := map[string]bool{}
	frontier := append([]string(nil), starts...)
	sort.Strings(frontier)
	for _, s := range frontier {
		seen[s] = true
	}
	var out []Dependent
	for level := 1; len(frontier) > 0 && (maxLevel == 0 || level <= maxLevel); level++ {
		var next []string
		for _, cur := range frontier {
			for _, e := range step(cur, adj) {
				if seen[e.To] {
					continue
				}
				seen[e.To] = true
				out = append(out, Dependent{ID: e.To, Level: level, From: cur, Type: e.Type, TypesOnly: e.NoSideEffect})
				if !e.NoSideEffect && !noProp[e.To] && (stop == nil || !stop(e.To)) {
					next = append(next, e.To)
				}
			}
		}
		sort.Strings(next)
		frontier = next
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Level != out[j].Level {
			return out[i].Level < out[j].Level
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// Rulers are what start is confronted with: the nodes that govern, specify, cover, test or
// seed it, and the doctrine rules and flag scenarios it realizes or is gated by — climbing,
// transitively, through those ties alone. Nothing that depends on start is a ruler.
func (g *Graph) Rulers(start string) []string {
	adj := g.evidenceIndex().adj
	seen := map[string]bool{start: true}
	var out []string
	queue := []string{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		var up []string
		for _, e := range adj.in[cur] {
			if dependsAtTo[e.Type] {
				up = append(up, e.From)
			}
		}
		for _, e := range adj.out[cur] {
			if e.Type == EdgeRealizes || e.Type == EdgeGatedBy {
				up = append(up, e.To)
			}
		}
		for _, p := range up {
			if !seen[p] {
				seen[p] = true
				out = append(out, p)
				queue = append(queue, p)
			}
		}
	}
	sort.Strings(out)
	return out
}

// TestsHolding are the test files to run again when the nodes change: those among them, and
// those whose evidence closure — the one `evidence-fresh` judges — holds one of them. Sorted.
func (g *Graph) TestsHolding(ids map[string]bool) []string {
	var out []string
	for _, n := range g.Nodes {
		if n.Kind != KindTest {
			continue
		}
		if ids[n.ID] {
			out = append(out, n.ID)
			continue
		}
		for f := range g.EvidenceClosure(n.ID) {
			if ids[f] {
				out = append(out, n.ID)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

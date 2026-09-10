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
			queue = append(queue, parent) // sobe mais (valida contra o avô), mas só SOBE
		}
	}
	sort.Strings(reached)
	return reached
}

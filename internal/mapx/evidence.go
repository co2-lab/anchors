// @anchors
//   code: EVMPV
//   ref: EVFRA

package mapx

import (
	"sort"
	"strings"
)

// Frescor de EVIDÊNCIA — distinto do frescor de CONFRONTO.
//
// O `Stamp` de uma aresta carimba o confronto do `check`: texto contra texto (a spec diz
// X, o código faz X). O `Signal` de um nó carimba a EXECUÇÃO: este teste rodou, e rodou
// contra a rev que ficou em `AtRev`. São perguntas diferentes com a mesma cara, e
// confundi-las faz acreditar que um placar de teste continua válido depois de o código
// mudar — que é o oposto de ter prova.
//
// `SignalStale` (ingest.go) responde pelo nó ISOLADO: "o próprio arquivo de teste mudou?".
// É pouco, e o quanto é pouco foi medido: `utils/login.yaml` é composto por 290 roteiros.
// Ao tocá-lo, a evidência dos 290 devia vencer, e nenhum sinal deles se altera — o teste
// não mudou, mudou o que ele executa.
//
// Daí o fecho: a evidência de um teste vence quando QUALQUER nó que ele alcança avança de
// rev. "Alcança" = descendo pelas arestas de saída (o que o teste compõe e do que depende),
// que é a direção em que uma mudança o afeta.
type EvidenceStale struct {
	Test    string   // o nó de teste cuja evidência venceu
	AtRev   string   // a rev carimbada na ingestão (do próprio teste)
	Culprit []string // os nós do fecho que avançaram de rev desde então
	Own     bool     // o próprio arquivo de teste mudou (o caso que SignalStale já pegava)
}

// EvidenceStaleFor devolve o veredito de frescor da evidência de um nó de teste.
//
// Devolve nil quando não há o que julgar: sem sinal (nunca foi ingerido — ausência de
// prova, não prova vencida) ou fecho inteiro na mesma rev.
//
// O fecho é comparado contra as revs GRAVADAS no momento da ingestão, guardadas em
// `Signal.ClosureRev`. Comparar contra a rev ATUAL de cada nó não diria nada: seria
// comparar o presente consigo mesmo. Quando `ClosureRev` está vazio (sinal ingerido por
// uma versão anterior, que não gravava o fecho), o julgamento cai para o nó isolado — a
// alternativa seria declarar vencida toda evidência antiga, transformando uma melhoria de
// precisão em 1400 falsos vencidos no dia em que entrasse.
func (g *Graph) EvidenceStaleFor(id string) *EvidenceStale {
	n := g.node(id)
	if n == nil || n.Signal == nil || n.Signal.AtRev == "" {
		return nil
	}
	out := &EvidenceStale{Test: id, AtRev: n.Signal.AtRev, Own: n.Signal.AtRev != n.Rev}

	revAtual := map[string]string{}
	for _, x := range g.Nodes {
		revAtual[x.ID] = x.Rev
		for rule, r := range x.OutRows {
			revAtual[OutRowKey(x.ID, rule)] = r
		}
	}
	// A stored closure is read through today's rule: what the rule leaves out — a screen on the
	// path, what a util reaches — stales nothing, even in a proof stamped before the rule.
	leftOut := map[string]bool{}
	now := g.evidenceClosure(id, true)
	for f := range g.evidenceClosure(id, false) {
		if _, ok := now[f]; !ok {
			leftOut[f] = true
		}
	}
	for alvo, revNaIngestao := range n.Signal.ClosureRev {
		if leftOut[alvo] {
			continue
		}
		atual, ok := revAtual[alvo]
		if ok && atual != revNaIngestao {
			out.Culprit = append(out.Culprit, alvo)
		}
		// A navigation row the test asserted that its spec no longer has: the navigation
		// is gone, and the assertion with it.
		if spec, _, isRow := strings.Cut(alvo, "#out:"); isRow && !ok && g.node(spec) != nil {
			out.Culprit = append(out.Culprit, alvo)
		}
	}
	out.Culprit = append(out.Culprit, g.divergedParts(n)...)
	sort.Strings(out.Culprit)
	if !out.Own && len(out.Culprit) == 0 {
		return nil
	}
	return out
}

// EvidenceClosure — os nós que este teste ALCANÇA descendo pelas arestas de saída, com a
// rev atual de cada um. É o que a ingestão carimba junto do sinal.
//
// Desce (não sobe) porque a pergunta é "do que este teste depende para significar o que
// significa": o roteiro compõe utils, o teste importa a unidade sob teste. Subir levaria à
// spec e à feature — o requisito que o teste prova, não os insumos que o fazem provar.
//
// Poda em `@noPropagation`, pelo mesmo motivo do impacto: um nó que declarou não propagar
// afirma que mudanças nele não descem, e respeitar isso aqui evita que um arquivo
// deliberadamente volátil vença a evidência de metade da suíte.
func (g *Graph) EvidenceClosure(id string) map[string]string {
	return g.evidenceClosure(id, true)
}

// evidenceClosure is the closure by today's rule — the asserts — or, with asserts false, by the
// rule before it, which followed every edge: the difference is what the rule leaves out.
func (g *Graph) evidenceClosure(id string, asserts bool) map[string]string {
	adj := g.adjacency()
	noProp := g.noPropSet()
	revs := map[string]string{}
	for _, x := range g.Nodes {
		revs[x.ID] = x.Rev
	}

	// What a test asserts, not the wiring that gets it there (DESIGN-evidence-follows-the-
	// asserts.md): an edge with no side effect on the tests — a navigation, an import of
	// types, as the project's flags say — is never followed (an asserted navigation is in the
	// closure by its Out row); and a test file the test depends on (a util) counts as a file,
	// with the test files it composes, and the walk does not descend from it into code.
	kind := map[string]Kind{}
	for _, x := range g.Nodes {
		kind[x.ID] = x.Kind
	}
	visto := map[string]bool{id: true}
	fila := []string{id}
	out := map[string]string{}
	for len(fila) > 0 {
		atual := fila[0]
		fila = fila[1:]
		wiring := asserts && atual != id && kind[atual] == KindTest
		for _, e := range adj.out[atual] {
			filho := e.To
			if visto[filho] || (asserts && e.NoSideEffect) || (wiring && kind[filho] != KindTest) {
				continue
			}
			visto[filho] = true
			if r, ok := revs[filho]; ok {
				out[filho] = r
			}
			// A capture reaches the unit's own file and images, and what the unit depends on
			// that has no capture of its own — a hook, a store, a service, transitively — up
			// to the next unit that is captured. A component the screen uses has its own
			// capture, and its change stales that one, not this.
			if e.Type == EdgeCaptures {
				for d, r := range g.uncapturedDeps(filho, adj, revs, asserts) {
					if !visto[d] && d != id {
						visto[d] = true
						out[d] = r
					}
				}
				continue
			}
			if !noProp[filho] && e.Type != EdgeComposes {
				fila = append(fila, filho)
			}
		}
	}
	// The navigation rows the test asserts, each with the revision of the row alone.
	if t := g.node(id); t != nil {
		for _, key := range t.Asserts {
			spec, rule, _ := strings.Cut(key, "#out:")
			if sp := g.node(spec); sp != nil {
				if r, ok := sp.OutRows[rule]; ok {
					out[key] = r
				}
			}
		}
	}
	return out
}

// divergedParts are the components of the unit a capture test captures whose own capture
// FAILED — diverged from its baseline beyond the tool's threshold — after this test last
// ran. A component's change stales only its own capture; a component whose capture
// diverged probably changed the look of every unit made of it, and their captures are asked
// again (`composes`, from the unit's Parts Used).
func (g *Graph) divergedParts(n *Node) []string {
	if n.Signal == nil {
		return nil
	}
	adj := g.adjacency()
	var out []string
	for _, c := range adj.out[n.ID] {
		if c.Type != EdgeCaptures {
			continue
		}
		for _, sp := range adj.in[c.To] {
			if sp.Type != EdgeSpecifies {
				continue
			}
			for _, part := range adj.out[sp.From] {
				if part.Type != EdgeComposes {
					continue
				}
				for _, t := range adj.in[part.To] {
					if t.Type != EdgeCaptures {
						continue
					}
					if tn := g.node(t.From); tn != nil && tn.Signal != nil && tn.Signal.Failed > 0 &&
						tn.Signal.IngestedAt > n.Signal.IngestedAt {
						out = append(out, part.To)
						break
					}
				}
			}
		}
	}
	sort.Strings(out)
	return out
}

// uncapturedDeps are what a captured unit's code depends on that no test captures: the
// declared dependencies of the spec that specifies it (`depends-on`), transitively through
// the specs of those, stopping at any file a test captures — that one is proven by its own
// capture —, and, by today's rule, past no dependency with no side effect on the tests.
func (g *Graph) uncapturedDeps(code string, adj adjacency, revs map[string]string, asserts bool) map[string]string {
	captured := func(id string) bool {
		for _, e := range adj.in[id] {
			if e.Type == EdgeCaptures {
				return true
			}
		}
		return false
	}
	out := map[string]string{}
	seen := map[string]bool{code: true}
	queue := []string{code}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		// What the file declares it depends on: its spec's Dependencies table, and its own
		// header's `dep:` lines.
		owners := []string{cur}
		for _, sp := range adj.in[cur] {
			if sp.Type == EdgeSpecifies {
				owners = append(owners, sp.From)
			}
		}
		for _, owner := range owners {
			for _, d := range adj.out[owner] {
				if d.Type != EdgeDependsOn || seen[d.To] || (asserts && d.NoSideEffect) {
					continue
				}
				seen[d.To] = true
				if captured(d.To) {
					continue
				}
				if r, ok := revs[d.To]; ok {
					out[d.To] = r
				}
				queue = append(queue, d.To)
			}
		}
	}
	return out
}

func (g *Graph) node(id string) *Node {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return &g.Nodes[i]
		}
	}
	return nil
}

// Node is the map's node with this id, or nil.
func (g *Graph) Node(id string) *Node { return g.node(id) }

// CapturesReaching are the capture tests whose evidence closure holds one of the files: a
// change to them is what those captures must be run again for — the unit's own file and
// images, or a hook, a store, a service it depends on that no capture of its own proves.
func (g *Graph) CapturesReaching(files []string) []string {
	want := map[string]bool{}
	for _, f := range files {
		want[f] = true
	}
	var out []string
	for _, n := range g.Nodes {
		if n.Kind != KindTest {
			continue
		}
		captures := false
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type == EdgeCaptures {
				captures = true
				break
			}
		}
		if !captures {
			continue
		}
		for f := range g.EvidenceClosure(n.ID) {
			if want[f] {
				out = append(out, n.ID)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

// OutRowKey names one navigation row of a spec in an evidence closure.
func OutRowKey(spec, rule string) string { return spec + "#out:" + rule }

package mapx

import "sort"

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
	}
	for alvo, revNaIngestao := range n.Signal.ClosureRev {
		if atual, ok := revAtual[alvo]; ok && atual != revNaIngestao {
			out.Culprit = append(out.Culprit, alvo)
		}
	}
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
	adj := g.adjacency()
	noProp := g.noPropSet()
	revs := map[string]string{}
	for _, x := range g.Nodes {
		revs[x.ID] = x.Rev
	}

	visto := map[string]bool{id: true}
	fila := []string{id}
	out := map[string]string{}
	for len(fila) > 0 {
		atual := fila[0]
		fila = fila[1:]
		for _, e := range adj.out[atual] {
			filho := e.To
			if visto[filho] {
				continue
			}
			visto[filho] = true
			if r, ok := revs[filho]; ok {
				out[filho] = r
			}
			if !noProp[filho] {
				fila = append(fila, filho)
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

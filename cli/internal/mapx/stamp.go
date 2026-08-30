package mapx

// A gravação do carimbo (o loop check→carimbo, PROPAGATION §3 + QUALITY §5). O gate
// roda POR NÓ; o carimbo é POR ARESTA. Esta é a cola: dado o veredito de cada nó
// confrontado, carimba as arestas cujas DUAS pontas foram confrontadas — porque só
// então a relação foi de fato revalidada. O carimbo grava as revs das pontas e o
// veredito, e é isso que destrava Stale(): na próxima vez que uma ponta avançar de
// rev, a aresta volta a ficar stale sozinha.

// NodeVerdict é o resultado agregado dos gates sobre UM nó, como o check o vê.
type NodeVerdict struct {
	ID     string
	Failed bool // reprovou ao menos um gate BLOQUEANTE
}

// StampEdges carimba, no grafo, todas as arestas cujas duas pontas estão em
// `verdicts`. Uma aresta recebe verdict "issue" se qualquer ponta falhou, senão
// "ok". `now` é a data (carimbada por quem chama — o pacote não inventa tempo).
// Devolve quantas arestas foram carimbadas.

// carimbo monta o Stamp preservando a data quando NADA mudou.
//
// A regra vale nos três pontos que carimbam (`StampEdges`, `StampEdge`, `StampNode`), e
// por isso mora aqui: repetida em cada um, ela se perderia no próximo que nascesse — foi
// assim que o `StampNodeByGate` passou despercebido na primeira tentativa.
func carimbo(anterior *Stamp, fromRev, toRev, verdict, now string) *Stamp {
	quando := now
	// `anterior.ChangedAt != ""` não é detalhe: um carimbo SEM data preservaria o vazio
	// para sempre — o buraco se perpetuaria justamente porque nada muda, e o campo sumiria
	// do mapa (`omitempty`). Sem data, a de hoje é a melhor resposta disponível: é quando
	// se soube que a relação estava assim.
	//
	// É o que também faz o mapa de uma versão anterior (quando o campo se chamava
	// `last_validated`) se converter sozinho no primeiro `check`. Mas a conferência NÃO é
	// código de migração e não tem prazo: ela vale para qualquer carimbo sem data, venha
	// de onde vier.
	if anterior != nil && anterior.ChangedAt != "" &&
		anterior.ValidatedFromRev == fromRev &&
		anterior.ValidatedToRev == toRev && anterior.Verdict == verdict {
		quando = anterior.ChangedAt
	}
	return &Stamp{
		ValidatedFromRev: fromRev,
		ValidatedToRev:   toRev,
		ChangedAt:        quando,
		Verdict:          verdict,
	}
}

func (g *Graph) StampEdges(verdicts []NodeVerdict, now string) int {
	failed := map[string]bool{}
	seen := map[string]bool{}
	for _, v := range verdicts {
		seen[v.ID] = true
		if v.Failed {
			failed[v.ID] = true
		}
	}
	stamped := 0
	for i := range g.Edges {
		e := &g.Edges[i]
		// só carimba se AMBAS as pontas foram confrontadas nesta rodada
		if !seen[e.From] || !seen[e.To] {
			continue
		}
		verdict := "ok"
		if failed[e.From] || failed[e.To] {
			verdict = "issue"
		}
		e.Stamp = carimbo(e.Stamp, g.nodeRev(e.From), g.nodeRev(e.To), verdict, now)
		stamped++
	}
	return stamped
}

// StampEdge carimba UMA aresta específica (from→to) com um veredito — usado pelo
// julgamento por IA, onde o confronto é da régua (guide) contra o alvo. Devolve
// false se a aresta não existe. O carimbo leva as revs atuais das pontas, então o
// veredito de IA envelhece (fica stale) se o alvo mudar depois — mesmo anti-drift.
func (g *Graph) StampEdge(from, to, verdict, now string) bool {
	for i := range g.Edges {
		e := &g.Edges[i]
		if e.From == from && e.To == to {
			e.Stamp = carimbo(e.Stamp, g.nodeRev(e.From), g.nodeRev(e.To), verdict, now)
			return true
		}
	}
	return false
}

// StampNode carimba TODAS as arestas que tocam um nó — o veredito de quem julgou aquela
// unidade, não um par dela.
//
// O `StampEdges` exige que AMBAS as pontas tenham sido confrontadas na mesma rodada, o
// que é certo para o `check` (que percorre muitos nós) e impossível para o `anchors
// judge`, que julga UM alvo: nenhuma aresta tem as duas pontas na lista, e o carimbo saía
// sempre zero. Medido: um nó com 42 arestas, `carimbado: 0`.
//
// A consequência não era cosmética. O veredito da IA abria a issue e não tocava o grafo —
// então um `anchors check` posterior não enxergava o achado, e as arestas do alvo
// continuavam "nunca validadas". O review sobrevivia à sessão como arquivo, e não virava
// pressão no pipeline.
func (g *Graph) StampNode(id, verdict, now string) int {
	return g.StampNodeByGate(id, verdict, now, "")
}

// StampNodeByGate registra o veredito de UM gate de julgamento sobre as arestas do
// alvo — em campo PRÓPRIO (`Edge.Julgamentos`), não no Stamp.
//
// A separação é necessária, e foi medida: o `check` reescreve o Stamp inteiro a cada
// rodada (ele resume o pior veredito de todos os gates daquele nó), então um veredito
// de IA guardado ali era apagado no primeiro check seguinte, e o julgamento voltava a
// ser perguntado como se ninguém tivesse lido. Carimbei 16 alvos e o contador não
// desceu de 16.
//
// Ainda carimba o Stamp também: o veredito de IA é confronto de verdade, e as arestas
// do alvo deixam de estar "nunca validadas".
func (g *Graph) StampNodeByGate(id, verdict, now, gateName string) int {
	stamped := 0
	for i := range g.Edges {
		e := &g.Edges[i]
		if e.From != id && e.To != id {
			continue
		}
		// O `Gate` entra depois: o helper decide a data, e o gate é de quem julgou.
		// Um gate diferente sobre o mesmo estado NÃO é mudança da relação — é outra
		// pergunta sobre ela —, então ele não faz a data avançar.
		st := carimbo(e.Stamp, g.nodeRev(e.From), g.nodeRev(e.To), verdict, now)
		st.Gate = gateName
		e.Stamp = st
		if gateName != "" {
			j := Julgamento{
				Gate:             gateName,
				Verdict:          verdict,
				ValidatedFromRev: g.nodeRev(e.From),
				ValidatedToRev:   g.nodeRev(e.To),
				ChangedAt:        now,
			}
			// Mesma regra do carimbo: rejulgar e achar o mesmo não é fato novo. Sem
			// isto, cada `anchors judge` reescreveria a data de todos os julgamentos.
			for k := range e.Julgamentos {
				a := e.Julgamentos[k]
				if a.Gate == gateName && a.Verdict == verdict &&
					a.ValidatedFromRev == j.ValidatedFromRev && a.ValidatedToRev == j.ValidatedToRev {
					j.ChangedAt = a.ChangedAt
				}
			}
			trocou := false
			for k := range e.Julgamentos {
				if e.Julgamentos[k].Gate == gateName {
					e.Julgamentos[k] = j
					trocou = true
					break
				}
			}
			if !trocou {
				e.Julgamentos = append(e.Julgamentos, j)
			}
		}
		stamped++
	}
	return stamped
}

// JulgadoPor diz se o nó já recebeu veredito DESTE gate e se ele ainda vale — isto é,
// se nenhuma das pontas mudou desde o julgamento.
//
// Basta UMA aresta viva: o `judge` registra em todas as que tocam o alvo, então
// qualquer uma responde. Se o alvo mudou depois, o veredito envelhece e volta a ser
// pergunta — julgamento não é selo permanente, é leitura datada.
func (g *Graph) JulgadoPor(id, gateName string) (verdict string, valido bool) {
	for _, e := range g.Edges {
		if e.From != id && e.To != id {
			continue
		}
		for _, j := range e.Julgamentos {
			if j.Gate != gateName {
				continue
			}
			if j.ValidatedFromRev != g.nodeRev(e.From) || j.ValidatedToRev != g.nodeRev(e.To) {
				continue // envelheceu: o alvo mudou depois do julgamento
			}
			return j.Verdict, true
		}
	}
	return "", false
}

// StaleEdges devolve as arestas atualmente stale (nunca validadas ou com uma ponta
// avançada desde o último confronto). É o que o comando `stale` lista.
func (g *Graph) StaleEdges() []Edge {
	var out []Edge
	for _, e := range g.Edges {
		if g.Stale(e) {
			out = append(out, e)
		}
	}
	return out
}

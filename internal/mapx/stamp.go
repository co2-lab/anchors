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

// stamp monta o Stamp preservando a data quando NADA mudou.
//
// A regra vale nos três pontos que carimbam (`StampEdges`, `StampEdge`, `StampNode`), e
// por isso mora aqui: repetida em cada um, ela se perderia no próximo que nascesse — foi
// assim que o `StampNodeByGate` passou despercebido na primeira tentativa.
func stamp(anterior *Stamp, fromRev, toRev, verdict, now string) *Stamp {
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
		// A WAIVER IS A PERSON'S DECISION, and the mechanical check does not confirm it:
		// both ends passing their gates says nothing about why a rule was waived. This loop
		// rewrote `waived` as `ok` on every edge whose ends it confronted — measured in
		// the reference app: one commit touching a feature turned 21 waived plan-chain judgments
		// into `ok`, dated today, with nobody having looked at them. The waiver is kept
		// as it was; when an end changes rev it goes stale by itself, and a person decides
		// again.
		if e.Stamp != nil && e.Stamp.Verdict == "waived" {
			continue
		}
		verdict := "ok"
		if failed[e.From] || failed[e.To] {
			verdict = "issue"
		}
		e.Stamp = stamp(e.Stamp, g.nodeRev(e.From), g.nodeRev(e.To), verdict, now)
		stamped++
	}
	return stamped
}

// StampEdge carimba UMA aresta específica (from→to) com um veredito — usado pelo
// julgamento por IA, onde o confronto é da régua (guide) contra o alvo. Devolve
// false se a aresta não existe. O carimbo leva as revs atuais das pontas, então o
// veredito de IA envelhece (fica stale) se o alvo mudar depois — mesmo anti-drift.
func (g *Graph) StampEdge(from, to, verdict, now string) bool {
	return g.StampEdgeByGate(from, to, verdict, now, "")
}

// StampEdgeByGate is StampEdge that also records WHICH gate judged the edge, as
// StampNodeByGate does. Without the record, a judgment gate that declares `guide:` stamped
// its guide→target edge and `JudgedBy` never saw it: the next check asked the same
// judgment again, right after it was answered.
func (g *Graph) StampEdgeByGate(from, to, verdict, now, gateName string) bool {
	for i := range g.Edges {
		e := &g.Edges[i]
		if e.From == from && e.To == to {
			st := stamp(e.Stamp, g.nodeRev(e.From), g.nodeRev(e.To), verdict, now)
			if gateName != "" {
				st.Gate = gateName
				g.recordJudgment(e, gateName, verdict, now)
			}
			if !keepsWaiver(e.Stamp, verdict, gateName) {
				e.Stamp = st
			}
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
		st := stamp(e.Stamp, g.nodeRev(e.From), g.nodeRev(e.To), verdict, now)
		st.Gate = gateName
		if !keepsWaiver(e.Stamp, verdict, gateName) {
			e.Stamp = st
		}
		if gateName != "" {
			g.recordJudgment(e, gateName, verdict, now)
		}
		stamped++
	}
	return stamped
}

// keepsWaiver diz se um julgamento deve deixar de pé o `waived` que a aresta carrega.
//
// O `StampEdges` já não reescreve um waiver (o check mecânico não confirma a decisão de
// uma pessoa); o `judge` reescrevia, e as duas portas divergiam sem nada dizer. A regra
// escolhida para o julgamento: um waiver responde à pergunta de UM gate. O veredito de
// OUTRO gate responde outra pergunta e não o desfaz — o julgamento dele vai para
// `Julgamentos`, e o carimbo fica. O MESMO gate julgando de novo é a pessoa respondendo
// de novo à mesma pergunta, e substitui o próprio waiver; um waiver novo sempre grava.
func keepsWaiver(atual *Stamp, verdict, gateName string) bool {
	return atual != nil && atual.Verdict == "waived" && verdict != "waived" &&
		(gateName == "" || atual.Gate != gateName)
}

// JudgedBy diz se o nó já recebeu veredito DESTE gate e se ele ainda vale — isto é,
// se nenhuma das pontas mudou desde o julgamento.
//
// Basta UMA aresta viva: o `judge` registra em todas as que tocam o alvo, então
// qualquer uma responde. Se o alvo mudou depois, o veredito envelhece e volta a ser
// pergunta — julgamento não é selo permanente, é leitura datada.
func (g *Graph) JudgedBy(id, gateName string) (verdict string, valido bool) {
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

// O CARIMBO casa pelo nome EXATO do gate.
//
// Houve uma normalização aqui (`mesmoGate` + `ResolveGateName`), que aceitava o nome
// legado em português. Ela tinha um defeito de assimetria: a LEITURA normalizava e a
// ESCRITA não — um projeto que renomeasse o gate ganhava um SEGUNDO carimbo em vez de
// atualizar o primeiro.
//
// O contrato de formato tornou a normalização desnecessária: o formato 2 só tem nome
// canônico, e quem está no 1 é recusado com a mensagem que manda migrar. A conversão é o
// passo `1→2` (`internal/migra/formato_2.go`) e roda uma vez.

// recordJudgment records on the edge the verdict a named gate gave it, at the ends' current
// revs — what `JudgedBy` reads to know the judgment was answered. One entry per gate.
func (g *Graph) recordJudgment(e *Edge, gateName, verdict, now string) {
	j := Judgment{
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

func edgeKey(e *Edge) string { return string(e.Type) + "\x00" + e.From + "\x00" + e.To }

// EdgeStamps is a snapshot of the stamps of the map's edges, by edge: what a command saw
// before stamping, so ApplyStampChanges can tell its own changes from another process's.
func (g *Graph) EdgeStamps() map[string]*Stamp {
	out := make(map[string]*Stamp, len(g.Edges))
	for i := range g.Edges {
		if g.Edges[i].Stamp != nil {
			s := *g.Edges[i].Stamp
			out[edgeKey(&g.Edges[i])] = &s
		}
	}
	return out
}

// ApplyStampChanges carries to g — the map as it is on disk now — the stamps `src` changed
// since `before` (src's snapshot taken before it stamped). A stamp that another process
// changed in g since `before` is kept: it is newer than what src saw. An edge g does not
// have is left out. It returns how many stamps it applied and how many it kept.
//
// A command that confronts for minutes — `check --all` — cannot hold the map's lock while
// it runs, so it stamps its own copy and brings over only what it changed.
func (g *Graph) ApplyStampChanges(src *Graph, before map[string]*Stamp) (applied, kept int) {
	here := make(map[string]*Edge, len(g.Edges))
	for i := range g.Edges {
		here[edgeKey(&g.Edges[i])] = &g.Edges[i]
	}
	for i := range src.Edges {
		e := &src.Edges[i]
		k := edgeKey(e)
		if sameStamp(e.Stamp, before[k]) {
			continue // src did not change it
		}
		d, ok := here[k]
		if !ok {
			continue
		}
		if !sameStamp(d.Stamp, before[k]) {
			kept++ // changed by another process meanwhile: theirs is newer
			continue
		}
		if e.Stamp == nil {
			d.Stamp = nil
		} else {
			s := *e.Stamp
			d.Stamp = &s
		}
		applied++
	}
	return applied, kept
}

func sameStamp(a, b *Stamp) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// RebaseRev moves a file from one revision to another and carries along everything measured
// at the first: its signals, proofs, coverage and mutation, the closures of the tests that
// reach it, and the stamps and judgments of its edges. It is for a change that proves
// nothing new nor undoes anything proven — the commit hook writing the day into the
// header's `updated_at`: without it, the map committed right after held the old revision,
// the next rebuild dropped the file's proofs over a date, and the CI found the committed
// map different from the one the build makes.
func (g *Graph) RebaseRev(id, from, to string) {
	if from == "" || from == to {
		return
	}
	move := func(r *string) {
		if *r == from {
			*r = to
		}
	}
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if s := n.Signal; s != nil {
			if n.ID == id {
				move(&s.AtRev)
				move(&s.MutationAtRev)
				for k, v := range s.ProvenRevBySuite {
					if v == from {
						s.ProvenRevBySuite[k] = to
					}
				}
				for k, c := range s.CoverageBySuite {
					move(&c.AtRev)
					s.CoverageBySuite[k] = c
				}
				for k, m := range s.MutationByScope {
					move(&m.AtRev)
					s.MutationByScope[k] = m
				}
			}
			if v, ok := s.ClosureRev[id]; ok && v == from {
				s.ClosureRev[id] = to
			}
		}
		if n.ID == id {
			move(&n.Rev)
		}
	}
	for i := range g.Edges {
		e := &g.Edges[i]
		if e.Stamp != nil {
			if e.From == id {
				move(&e.Stamp.ValidatedFromRev)
			}
			if e.To == id {
				move(&e.Stamp.ValidatedToRev)
			}
		}
		for k := range e.Julgamentos {
			if e.From == id {
				move(&e.Julgamentos[k].ValidatedFromRev)
			}
			if e.To == id {
				move(&e.Julgamentos[k].ValidatedToRev)
			}
		}
	}
}

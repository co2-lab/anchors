// @anchors
//   ref: PRFLO

package gate

import (
	"sort"
	"time"
)

// Profile é a agregação dos vereditos (QUALITY §6): a qualidade não é um número, é
// o conjunto de vereditos por gate. Deriva a decisão de promoção ("todos os
// bloqueantes passaram") e as issues (todo fail é uma issue — QUALITY §2).
type Profile struct {
	Results  []Result
	ByGate   map[string]GateSummary // resumo por gate
	Passed   bool                   // todos os gates BLOQUEANTES passaram?
	Failures []Result               // os fails (viram issues)
	Blocked  []Result               // os fails de gates bloqueantes (barram a promoção)
	Judged   []Result               // os que aguardam julgamento de IA (verdict Judge)
}

// GateSummary — contagem de vereditos de um gate.
type GateSummary struct {
	Gate     string
	Blocking bool
	Pass     int
	Fail     int
	Skip     int
	Pending  int // could not measure
	Diverge  int // measured and found something short of wrong
	Judge    int // aguardando julgamento de IA
	// Duracao e' o tempo somado de todas as confrontacoes deste gate, e Pior e' a mais
	// cara delas.
	//
	// As duas medidas juntas separam o que uma so' esconde: um gate pode custar caro por
	// VOLUME (barato por alvo, multiplicado por centenas) ou por ALVO (uma execucao que
	// sobe um compilador). O total sozinho nao distingue os dois casos, e a otimizacao
	// de cada um e' oposta — reduzir alvos no primeiro, trocar a ferramenta no segundo.
	Duracao time.Duration
	Pior    time.Duration
}

// Aggregate monta o perfil a partir dos resultados brutos.
func Aggregate(results []Result) Profile {
	p := Profile{Results: results, ByGate: map[string]GateSummary{}, Passed: true}
	for _, r := range results {
		s := p.ByGate[r.Gate]
		s.Gate = r.Gate
		s.Blocking = r.Blocking
		switch r.Verdict {
		case Pass:
			s.Pass++
		case Fail:
			s.Fail++
			if !r.Ignored() {
				p.Failures = append(p.Failures, r)
			}
			if r.Blocks() {
				p.Blocked = append(p.Blocked, r)
				p.Passed = false
			}
		case Skip:
			s.Skip++
		case Pending, Diverge:
			if r.Verdict == Diverge {
				s.Diverge++
			} else {
				s.Pending++
			}
			// A divergence or a pending item bars when its level's state is `block` (the gate's
			// `severity`, see `markSeverity`). "Nothing to confront" is Skip, not Pending — when
			// it was Pending, making Pending block failed 411 nodes at once.
			if r.Blocks() {
				p.Blocked = append(p.Blocked, r)
				p.Passed = false
			}
		case Judge:
			s.Judge++
			p.Judged = append(p.Judged, r)
		}
		s.Duracao += r.Duracao
		if r.Duracao > s.Pior {
			s.Pior = r.Duracao
		}
		p.ByGate[r.Gate] = s
	}
	return p
}

// NodeVerdict é o resultado agregado por NÓ (um nó pode ser tocado por vários
// gates). Failed = barrou a promoção em ao menos um gate BLOQUEANTE (Fail, ou Pending
// que `Impede`). É o insumo que o mapa usa
// para carimbar as arestas (loop check→carimbo).
type NodeVerdict struct {
	ID     string
	Failed bool
}

// NodeVerdicts colapsa os Results (por gate×alvo) em um veredito por nó. Só entram
// nós que foram efetivamente confrontados (têm ao menos um Result não-skip).
func (p Profile) NodeVerdicts() []NodeVerdict {
	confronted := map[string]bool{}
	failed := map[string]bool{}
	for _, r := range p.Results {
		if r.Verdict == Skip || r.Verdict == Judge {
			continue // skip: não se aplica. Judge: confronto ainda não concluído (a
			// IA não julgou) — só carimba quando `anchors judge` gravar o veredito.
		}
		confronted[r.Target] = true
		// The same test `Aggregate` uses to block promotion: a blocking gate's Fail, or
		// its Pending marked `Impede`. Only the Fail was counted here, so a node whose
		// impeding pending refused the promotion was stamped as not failed on the map.
		if r.Blocks() {
			failed[r.Target] = true
		}
	}
	var out []NodeVerdict
	for id := range confronted {
		out = append(out, NodeVerdict{ID: id, Failed: failed[id]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// GateNames devolve os nomes dos gates do perfil, ordenados (para saída estável).
func (p Profile) GateNames() []string {
	var names []string
	for name := range p.ByGate {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

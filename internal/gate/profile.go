package gate

import "sort"

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
	Pending  int
	Judge    int // aguardando julgamento de IA
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
			p.Failures = append(p.Failures, r)
			if r.Blocking {
				p.Blocked = append(p.Blocked, r)
				p.Passed = false
			}
		case Skip:
			s.Skip++
		case Pending:
			s.Pending++
			// Pendência NÃO barra por si — e a medição mostra por quê: num repositório real,
			// fazer `Pending` barrar em gate bloqueante reprovou 411 nós de uma vez. O
			// `feature-test-match` (bloqueante) tem 410 pendências que significam "não tive o
			// que confrontar", não "há decisão por tomar".
			//
			// São dois sentidos no mesmo veredito, e só o segundo deveria impedir a promoção.
			// Quem distingue é o GATE, que sabe o que mediu — ver `Impede`.
			if r.Blocking && r.Impede {
				p.Blocked = append(p.Blocked, r)
				p.Passed = false
			}
		case Judge:
			s.Judge++
			p.Judged = append(p.Judged, r)
		}
		p.ByGate[r.Gate] = s
	}
	return p
}

// NodeVerdict é o resultado agregado por NÓ (um nó pode ser tocado por vários
// gates). Failed = reprovou ao menos um gate BLOQUEANTE. É o insumo que o mapa usa
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
		if r.Verdict == Fail && r.Blocking {
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

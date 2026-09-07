package board

import "testing"

// A REGRA 1: retomar o próprio vence a prioridade do board.
//
// Um card já reivindicado por este agente carrega o contexto da sessão dele. Entregar
// outro joga fora o que já foi lido — e deixa dois agentes com metade do entendimento.
func TestClaim_retomarOProprioVenceAPrioridade(t *testing.T) {
	meu := Card{Number: 7, Owner: "maq/sessao-1", Labels: []string{StateInProgress}}
	if liveState(meu) != StateInProgress {
		t.Errorf("o estado vivo do card próprio não foi reconhecido: %q", liveState(meu))
	}
}

// A REGRA 3: `needs-user` nunca é entregue.
//
// É o card waiting esperando decisão de gente, e dá-lo a um agente o faz decidir sozinho —
// que é o que o `escalate` existe para impedir.
func TestParado_recusaOsDoisNomesDeNeedsUser(t *testing.T) {
	for _, label := range []string{StateNeedsUser, StateNeedsUserPt} {
		c := Card{Labels: []string{StateToDo, label}}
		if !waiting(c) {
			t.Errorf("card com %q não foi reconhecido como waiting", label)
		}
	}
	if waiting(Card{Labels: []string{StateToDo}}) {
		t.Error("card sem needs-user foi tratado como waiting")
	}
}

// O DONO é o ÚLTIMO comentário `anchors-owner:`.
//
// A posse muda de mão — um card rejeitado volta ao autor original — e cada reivindicação
// fica registrada. Ler o primeiro devolveria um dono que já passou o trabalho adiante.
func TestLastOwner_oUltimoComentarioVence(t *testing.T) {
	r := rawCard{Comments: []struct{ Body string }{
		{Body: "anchors-owner: maq-a/sessao-1"},
		{Body: "um comentário qualquer no meio"},
		{Body: "anchors-owner: maq-b/sessao-2"},
	}}
	if got := lastOwner(r); got != "maq-b/sessao-2" {
		t.Errorf("dono lido: %q — esperava o último", got)
	}
}

// Card sem comentário de posse não tem dono, e é livre.
func TestLastOwner_semComentarioNaoTemDono(t *testing.T) {
	if got := lastOwner(rawCard{}); got != "" {
		t.Errorf("card sem comentário devolveu dono %q", got)
	}
}

// `workflow.labels` vazio é ERRO, não default.
//
// Sem ele o claim puxaria qualquer issue do repositório — inclusive as de produto, que não
// têm a forma que o ciclo espera. O `config.go` já declara isso: "vazio no modo github é
// erro de configuração, não default".
func TestList_labelsVazioEhErro(t *testing.T) {
	c := Client{Repo: "org/repo"}
	if _, err := c.list(StateToDo); err == nil {
		t.Error("labels vazio passou — o claim puxaria issue de produto")
	}
}

// liveState devolve o estado mais AVANÇADO do card.
//
// Um card pode carregar mais de uma label de estado por engano (um workflow que
// acrescentou sem remover). Reportar o mais avançado é o que descreve onde o trabalho
// está — reportar `to-do` num card em revisão faria o agente refazer.
func TestEstadoVivo_reportaOMaisAvancado(t *testing.T) {
	c := Card{Labels: []string{StateToDo, StateInReview}}
	if got := liveState(c); got != StateInReview {
		t.Errorf("estado reportado: %q — esperava o mais avançado", got)
	}
}

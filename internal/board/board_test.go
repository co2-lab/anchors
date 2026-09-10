package board

import "testing"

// A REGRA 1: retomar o próprio vence a prioridade do board.
//
// Um card já reivindicado por este agente carrega o contexto da sessão dele. Entregar
// outro joga fora o que já foi lido — e deixa dois agentes com metade do entendimento.
func TestMine_retomarOProprioVenceAPrioridade(t *testing.T) {
	meu := Card{Number: 7, Owner: "maq/sessao-1", Labels: []string{StateInProgress}}
	if liveState(meu) != StateInProgress {
		t.Errorf("o estado vivo do card próprio não foi reconhecido: %q", liveState(meu))
	}
}

// A REGRA 3: `needs-user` só vai para quem DECLAROU que decide o produto.
//
// Antes ele era recusado por todos, e os escalonados só saíam do board por intervenção
// manual. A declines agora depende de quem pergunta: num projeto com vários devs, cada um
// roda o seu agente, e um agente que pega um card escalonado e pergunta a quem o está
// rodando obtém uma resposta que pode não ser a do dono do projeto.
//
// O padrão é fechado — o zero-value do `Client` não atua —, e o custo de errar para o lado
// aberto é alguém decidir o produto sem autoridade, que é invisível depois do fato.
func TestRecusa_escalonadoDependeDaDeclaracao(t *testing.T) {
	for _, label := range []string{StateNeedsUser, StateNeedsUserPt} {
		card := Card{Labels: []string{StateToDo, label}}

		semDeclarar := Client{}
		if !semDeclarar.declines(card) {
			t.Errorf("%q: o cliente que não declarou nada NÃO pode pegar escalonado", label)
		}

		optOut := Client{UserIssues: false}
		if !optOut.declines(card) {
			t.Errorf("%q: quem declarou que não decide continua recusando", label)
		}

		optIn := Client{UserIssues: true}
		if optIn.declines(card) {
			t.Errorf("%q: quem declarou que decide o produto pode pegá-lo", label)
		}
	}
}

// O card COMUM não é recusado por ninguém — a declaração só governa o escalonado.
func TestRecusa_oCardComumPassaNosDoisModos(t *testing.T) {
	card := Card{Labels: []string{StateToDo}}
	for _, c := range []Client{{}, {UserIssues: true}} {
		if c.declines(card) {
			t.Errorf("card comum recusado por Client{UserIssues:%v}", c.UserIssues)
		}
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

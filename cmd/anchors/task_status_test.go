package main

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/board"
)

// O comando existe por um relato que estava CERTO e insuficiente: o agente de outro dev
// diagnosticou uma falha de CI, consertou, empurrou, e encerrou o turno com "aguardando a
// nova rodada" — com o card `in-progress` no nome dele e o veredito do check sem ninguém
// para ler.
//
// Estes testes cobram o que o formato tem de dizer nos estados em que o silêncio custa.

func testCard(n int, state, titulo string) *board.Card {
	return &board.Card{Number: n, Title: titulo, State: state}
}

func TestTaskStatus_prSemCheckNaoPassaPorConferido(t *testing.T) {
	// O caso medido nesta sessão: o PR #368 abriu e NENHUM check rodou. "PR #368 open"
	// sozinho lê-se como trabalho conferido, e é a família de defeito que este projeto já
	// mediu três vezes — o verde que não fez o trabalho.
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "[INDTN] implementar spec"),
		Branch: "impl-x", Clean: true,
		PR: &branchPR{Number: 368, Estado: "OPEN", Total: 0, Checks: map[string]int{}},
	})
	// Duas asserções para DOIS lugares distintos, e a distinção importa: a linha do PR é
	// o que quem lê rápido vê, e o próximo passo é o que o agente segue. Mutar só a linha
	// do PR (`· nenhum check rodou` → `· ok`) passava batido enquanto as duas asserções
	// compartilhavam a mesma busca no texto todo — o próximo passo repete a frase, e
	// cobria a ausência na linha de cima.
	linhaDoPR := ""
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, "PR ") {
			linhaDoPR = l
			break
		}
	}
	if !strings.Contains(linhaDoPR, "nenhum check rodou") {
		t.Errorf("a LINHA DO PR tem de dizer que nenhum check rodou; era %q", linhaDoPR)
	}
	if !strings.Contains(out, "não confie no verde que não existe") {
		t.Errorf("o próximo passo deveria recusar o PR sem check; saída:\n%s", out)
	}
}

func TestTaskStatus_checkEmCursoNaoEhCheckQuePassou(t *testing.T) {
	// A confusão exata que encerrou o turno errado. Um relato que soma "em curso" ao
	// "passou" produz "4/4 passou" para um PR que ainda não terminou de rodar.
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "x"),
		Branch: "impl-x", Clean: true,
		PR: &branchPR{Number: 368, Estado: "OPEN", Total: 4,
			Checks: map[string]int{"passou": 3, "em curso": 1}},
	})
	if strings.Contains(out, "4/4") {
		t.Errorf("3 passaram e 1 está em curso: não é 4/4; saída:\n%s", out)
	}
	if !strings.Contains(out, "em curso") {
		t.Errorf("o check em curso tem de aparecer; saída:\n%s", out)
	}
	if !strings.Contains(out, "--watch") {
		t.Errorf("o próximo passo de um CI rodando é ESPERAR o veredito; saída:\n%s", out)
	}
	if !strings.Contains(out, "não encerre o turno") {
		t.Errorf("o formato tem de negar que esperar é uma parada; saída:\n%s", out)
	}
}

func TestTaskStatus_checkReprovadoEhTrabalhoDesteCard(t *testing.T) {
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "x"),
		Branch: "impl-x", Clean: true,
		PR: &branchPR{Number: 368, Estado: "OPEN", Total: 4,
			Checks: map[string]int{"passou": 3, "reprovou": 1}},
	})
	if !strings.Contains(out, "trabalho DESTE card") {
		t.Errorf("o vermelho não é um card novo; saída:\n%s", out)
	}
	// A ordem importa: quem lê rápido tem de bater no problema, não no que deu certo.
	if strings.Index(out, "reprovou") > strings.Index(out, "passou") {
		t.Errorf("a reprovação deveria vir antes do que passou; saída:\n%s", out)
	}
}

func TestTaskStatus_naoAconselhaAbrirPRDeCardEmRevisao(t *testing.T) {
	// A primeira versão dizia "abrir o PR" para qualquer card aberto sem PR no branch, e
	// num `in-review` isso é conselho errado — o trabalho já foi entregue. Um próximo passo
	// errado é pior que nenhum: ele parece derivado do estado.
	out := renderTaskStatus(taskState{
		Card:   testCard(301, "anchors:in-review", "[X] tela"),
		Branch: "develop", Clean: true,
	})
	if strings.Contains(out, "abrir o PR") {
		t.Errorf("um card em revisão não tem PR a abrir; saída:\n%s", out)
	}
	if !strings.Contains(out, "em revisão") {
		t.Errorf("o estado do card tem de aparecer traduzido; saída:\n%s", out)
	}
}

func TestTaskStatus_escalonadosAparecemSeparados(t *testing.T) {
	// `needs-user` é a única pendência que NÃO se resolve continuando a trabalhar. Um
	// relato que a omite convida o leitor a esperar por algo que só ele pode destravar.
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "x"),
		Branch: "impl-x", Clean: true,
		Blocked: []board.Card{
			{Number: 99, Title: "[DEC1] qual vocabulário?"},
			{Number: 12, Title: "[DEC2] onde fica o limite?"},
		},
	})
	if !strings.Contains(out, "Esperando decisão de pessoa") {
		t.Errorf("os escalatedCards precisam de seção própria; saída:\n%s", out)
	}
	if !strings.Contains(out, "#99") || !strings.Contains(out, "#12") {
		t.Errorf("os dois números têm de aparecer; saída:\n%s", out)
	}
	if !strings.Contains(out, "nenhum agente resolve") {
		t.Errorf("o formato tem de dizer POR QUE estas não avançam; saída:\n%s", out)
	}
}

func TestTaskStatus_semEscalonadoNaoInventaSecao(t *testing.T) {
	out := renderTaskStatus(taskState{
		Card: testCard(303, "anchors:in-progress", "x"), Branch: "impl-x", Clean: true,
	})
	if strings.Contains(out, "Esperando decisão de pessoa") {
		t.Errorf("sem escalonado a seção não deveria existir; saída:\n%s", out)
	}
}

func TestTaskStatus_asDuasLacunasSaoExplicitas(t *testing.T) {
	// O que a máquina não sabe, ela PEDE. Uma lacuna nomeada é preenchida; uma seção
	// ausente não é notada — e era a ausência dela que deixava cada agente improvisar.
	out := renderTaskStatus(taskState{Branch: "develop", Clean: true})
	for _, want := range []string{"O que provei", "O que ficou de fora"} {
		if !strings.Contains(out, want) {
			t.Errorf("o formato deveria pedir %q; saída:\n%s", want, out)
		}
	}
}

func TestTaskStatus_mudancaNaoCommitadaVemAntesDeQualquerOutroPasso(t *testing.T) {
	// Trabalho não commitado é o único estado em que TODO conselho seguinte é prematuro.
	out := renderTaskStatus(taskState{
		Card: testCard(303, "anchors:in-progress", "x"), Branch: "impl-x", Clean: false,
	})
	if !strings.Contains(out, "não commitada") {
		t.Errorf("a mudança pendente tem de aparecer; saída:\n%s", out)
	}
	if strings.Contains(out, "abrir o PR") {
		t.Errorf("com trabalho não commitado, abrir PR é prematuro; saída:\n%s", out)
	}
}

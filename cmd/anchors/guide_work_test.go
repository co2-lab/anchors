package main

import (
	"strings"
	"testing"
)

// E o guia PRECISA cobrir o que o card deixou de dizer, senão a centralização perdeu
// justamente a informação que motivou tudo.
func TestGuiaDeTrabalhoCobreOAchadoDoAgente(t *testing.T) {
	// O que se confronta são NOMES DE COMANDO E FLAG — o vocabulário do próprio sistema,
	// que só muda quando a interface muda (e aí o teste falhar é o comportamento certo).
	//
	// PROSA fica de fora de propósito. A primeira versão exigia "MESMO PR", e quando a
	// frase quebrou de linha eu "consertei" embutindo o `\n` na assertiva — deixando o
	// teste preso à LARGURA DO PARÁGRAFO. Qualquer reflow, reescrita ou tradução do guia
	// o quebraria, sem que nada tivesse piorado de verdade.
	for _, exigido := range []string{
		"anchors escalate",        // como registrar o achado
		"--card",                  // o vínculo, sem o qual ele nasce solto
		"--for-user",              // a saída para quando muda a direção
		"anchors judge --pending", // o que barra o commit
		"@TBD",                    // a peça que ainda não nasceu
		"anchors:sob-",            // a label que amarra o achado ao trabalho
	} {
		if !strings.Contains(workGuide, exigido) {
			t.Errorf("o guia de trabalho deve cobrir %q — o card não diz mais isso", exigido)
		}
	}
}

// O guia tinha de dizer que EMPURRAR NÃO É ENTREGAR, e não dizia. Um agente de outro dev
// diagnosticou certo por que o CI reprovou (mapa gerado antes do commit que consertava a
// config), reconstruiu, empurrou — e encerrou o turno com "aguardando a nova rodada". O
// card ficou `in-progress` com o nome dele nele.
//
// O guia terminava no `pr-body`: nada dizia como esperar, nem que a espera não é uma
// parada. Este teste cobra as duas coisas que faltavam — o comando que BLOQUEIA até o
// veredito, e a recusa explícita da terceira saída (relatar e parar).
func TestGuiaDeTrabalhoCobraEsperarOVereditoDoCI(t *testing.T) {
	for _, want := range []string{
		"--watch",     // o comando que espera pelo agente
		"work review", // o que falta no PR que é seu
		"escalate",    // a única saída legítima que não é o verde
		"in-progress", // o custo de parar no meio: o card fica com seu nome
		"metade",      // o diagnóstico correto é metade do trabalho
	} {
		if !strings.Contains(workGuide, want) {
			t.Errorf("guia de trabalho deveria cobrir a espera do CI (esperava %q)", want)
		}
	}
}

// A regra é sobre CONTINUAR, e o teste acima casaria mesmo que o guia dissesse o oposto.
// Este confronta a afirmação: o guia tem de negar que empurrar encerra o turno.
func TestGuiaDeTrabalhoNegaQueEmpurrarEncerraOTurno(t *testing.T) {
	if !strings.Contains(workGuide, "não é um ponto de parada") {
		t.Error("o guia deveria NEGAR explicitamente que empurrar um commit é um ponto de parada")
	}
	if !strings.Contains(workGuide, "não fecha o card") {
		t.Error("o guia deveria dizer que abrir o PR não fecha o card")
	}
}

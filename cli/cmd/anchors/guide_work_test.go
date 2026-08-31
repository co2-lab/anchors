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
		"--para-usuario",          // a saída para quando muda a direção
		"anchors judge --pending", // o que barra o commit
		"@TBD",                    // a peça que ainda não nasceu
		"anchors:sob-",            // a label que amarra o achado ao trabalho
	} {
		if !strings.Contains(workGuide, exigido) {
			t.Errorf("o guia de trabalho deve cobrir %q — o card não diz mais isso", exigido)
		}
	}
}

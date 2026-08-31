package main

import (
	"strings"
	"testing"
)

// E o guia PRECISA cobrir o que o card deixou de dizer, senão a centralização perdeu
// justamente a informação que motivou tudo.
func TestGuiaDeTrabalhoCobreOAchadoDoAgente(t *testing.T) {
	for _, exigido := range []string{
		"anchors escalate", // como registrar
		"--card",           // o vínculo, sem o qual o achado nasce solto
		"MESMO\nPR",        // que se entrega junto (a linha quebra no meio)
		"--para-usuario",   // a saída para quando muda a direção
		"judge --pending",  // o que barra o commit
	} {
		if !strings.Contains(workGuide, exigido) {
			t.Errorf("o guia de trabalho deve cobrir %q — o card não diz mais isso", exigido)
		}
	}
}

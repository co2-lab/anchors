package main

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// A ISSUE tem de responder três coisas, e a terceira é a que costuma faltar: como
// destravar. Sem ela o card fica parado esperando alguém adivinhar o protocolo.
func TestEscalada_dizPorQueParouEComoDestravar(t *testing.T) {
	corpo := corpoDaEscalada("A spec pede cache; o plano diz que não haveria cache.",
		"plans/0001-fundacao.md", "12")

	for _, exigido := range []string{
		"A spec pede cache",         // o motivo, com as palavras de quem viu
		"plans/0001-fundacao.md",    // onde
		"R0001",                     // como registrar a decisão
		initx.LabelPrecisaDoUsuario, // o que remover para destravar
		"#12",                       // onde o trabalho parou
	} {
		if !strings.Contains(corpo, exigido) {
			t.Errorf("o corpo da escalada deve conter %q; veio:\n%s", exigido, corpo)
		}
	}
}

// Sem `--card` o comando ainda serve: nem toda incoerência é achada com um card na mão
// (o revisor lendo um PR, por exemplo). O corpo não pode citar um card que não existe.
func TestEscalada_semCardNaoInventaReferencia(t *testing.T) {
	corpo := corpoDaEscalada("O plano contradiz a si mesmo entre F02 e F04.", "", "")
	if strings.Contains(corpo, "#") && strings.Contains(corpo, "Trabalho parado") {
		t.Errorf("sem --card não pode citar card; veio:\n%s", corpo)
	}
	if !strings.Contains(corpo, "contradiz a si mesmo") {
		t.Error("o motivo tem de sobreviver mesmo sem card")
	}
}

// O TÍTULO da issue é uma linha. Um motivo longo não pode virar um título ilegível na
// lista de issues, que é onde o usuário vai encontrá-lo.
func TestEscalada_tituloCabeEmUmaLinha(t *testing.T) {
	longo := strings.Repeat("uma explicação bem detalhada da incoerência ", 5)
	got := primeiraLinhaDoMotivo(longo)
	if len(got) > 70 {
		t.Errorf("título com %d chars; deveria caber em 70: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("o corte deve sinalizar que há mais: %q", got)
	}
	// Motivo de várias linhas: o título é a primeira.
	if got := primeiraLinhaDoMotivo("primeira linha\nsegunda linha"); got != "primeira linha" {
		t.Errorf("o título é a primeira linha, veio %q", got)
	}
}

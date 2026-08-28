package main

import (
	"strings"
	"testing"
)

// Quem cai na guarda de TTY — um agente, um pipe, o CI — precisa saber que existe saída.
// Sem citar o modo não-interativo, o comando parece um beco: "é interativo e não há
// terminal" descreve o problema e não oferece nada.
func TestErroSemTTYOfereceOModoNaoInterativo(t *testing.T) {
	msg := erroSemTTY("Nada foi escrito.").Error()

	if !strings.Contains(msg, "--non-interactive") {
		t.Errorf("a mensagem tem de nomear a saída:\n%s", msg)
	}
	if !strings.Contains(msg, "Nada foi escrito.") {
		t.Errorf("o contexto de quem chamou tem de aparecer:\n%s", msg)
	}
	// As DUAS chamadas da mesma flag: a sem respostas (que pergunta) e a com respostas
	// (que aplica). Mostrar só uma deixaria metade do contrato invisível — e é a metade
	// que o leitor precisa primeiro.
	if strings.Count(msg, "--non-interactive") < 2 {
		t.Errorf("as duas chamadas têm de aparecer:\n%s", msg)
	}
	if !strings.Contains(msg, "--artifacts") {
		t.Errorf("a segunda chamada precisa de um exemplo de resposta:\n%s", msg)
	}
}

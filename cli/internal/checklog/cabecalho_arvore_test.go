package checklog

import (
	"strings"
	"testing"
	"time"
)

// O cabeçalho carimba o estado da árvore para quem RELÊ o relatório depois. Sem
// repositório não há como contar os arquivos sujos — e escrever "limpa" ali afirma um
// estado que ninguém verificou, para um leitor que não tem como desconfiar.
func TestCabecalhoNaoAfirmaArvoreLimpaSemSaber(t *testing.T) {
	quando := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)

	h := Cabecalho("anchors check", "abc123", "assunto", -1, quando)

	if strings.Contains(h, "árvore: limpa") {
		t.Errorf("afirmou limpeza sem ter conseguido contar: %s", h)
	}
	if !strings.Contains(h, "desconhecida") {
		t.Errorf("o cabeçalho tem de dizer que não sabe: %s", h)
	}
}

// A contrapartida: contagem de verdade continua sendo relatada como antes.
func TestCabecalhoRelataContagemReal(t *testing.T) {
	quando := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)

	if h := Cabecalho("c", "h", "a", 0, quando); !strings.Contains(h, "árvore: limpa") {
		t.Errorf("0 sujos É árvore limpa: %s", h)
	}
	if h := Cabecalho("c", "h", "a", 1, quando); !strings.Contains(h, "1 arquivo modificado") {
		t.Errorf("1 sujo: %s", h)
	}
	if h := Cabecalho("c", "h", "a", 5, quando); !strings.Contains(h, "5 arquivos modificados") {
		t.Errorf("5 sujos: %s", h)
	}
}

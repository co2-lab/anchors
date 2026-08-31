package main

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// A SINTAXE É DA PLATAFORMA, e por isso vive num mapa — acrescentar uma é acrescentar
// uma linha, e o gerador não precisa saber quantas existem.
//
// O que NÃO pode é a sintaxe vazar para a doutrina: o Anchors é multi-idioma, e um gate
// que obrigue o corpo do PR a estar em inglês não é régua do Anchors — é uma exigência do
// GitHub disfarçada de regra.
func TestSintaxeDeFechamentoEhPorPlataforma(t *testing.T) {
	if _, ok := sintaxeDeFechamento["github"]; !ok {
		t.Fatal("o github precisa ter sintaxe declarada — é a plataforma do modo `github`")
	}
	// O formato tem de conter `%s`: sem ele o número do card não entra, e o comando
	// imprimiria a mesma linha para todos.
	for plataforma, forma := range sintaxeDeFechamento {
		if !strings.Contains(forma, "%s") {
			t.Errorf("a sintaxe de %q não tem onde pôr o número do card: %q", plataforma, forma)
		}
	}
}

// `--cards` aceita as formas que uma pessoa escreve: com `#`, sem, com espaço.
// Recusar "#44" por causa do sustenido seria atrito sem razão — é como o card aparece
// em todo lugar do GitHub.
func TestCardsPedidosAceitaAsFormasQueSeEscreve(t *testing.T) {
	cfg := &config.Config{}
	for _, entrada := range []string{"44", "#44", " 44 ", "#44 "} {
		got := cardsPedidos(entrada, cfg)
		if len(got) != 1 || got[0] != "44" {
			t.Errorf("%q deveria virar [44], veio %v", entrada, got)
		}
	}
	// Vários de uma vez: o trabalho fecha o card E os achados que nasceram sob ele.
	if got := cardsPedidos("44, #49,50", cfg); len(got) != 3 {
		t.Errorf("três cards deveriam virar três entradas, veio %v", got)
	}
	// Vazio não inventa card: sem `--cards` e sem agente, quem chama recebe erro em vez
	// de um PR que não fecha nada.
	if got := cardsPedidos("  ", cfg); len(got) != 0 {
		t.Errorf("entrada vazia não pode inventar card, veio %v", got)
	}
}

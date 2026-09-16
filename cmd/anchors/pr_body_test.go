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
func TestSintaxeDeVinculoEhPorPlataforma(t *testing.T) {
	if _, ok := linkSyntax["github"]; !ok {
		t.Fatal("o github precisa ter sintaxe declarada — é a plataforma do modo `github`")
	}
	// O formato tem de conter `%s`: sem ele o número do card não entra, e o comando
	// imprimiria a mesma linha para todos.
	for plataforma, forma := range linkSyntax {
		if !strings.Contains(forma, "%s") {
			t.Errorf("a sintaxe de %q não tem onde pôr o número do card: %q", plataforma, forma)
		}
	}

	// VINCULAR E NÃO FECHAR — e esta é a régua que não deixa voltar.
	//
	// Isto gerava `Closes #N`, que faz duas afirmações de uma vez: o pipeline lia "o PR
	// concluiu a implementação" e o GitHub lia "feche a issue". A esteira não acaba em
	// `ready-to-test` — vêm `in-test`, `ready-to-release` e `production`, que são do CD e
	// o Anchors não rastreia. O card fechava com três estados à frente.
	//
	// MEDIDO no projeto de referência: 71 cards fechados em `ready-to-test` contra 7
	// abertos, 35 num só dia. A coluna que devia acumular o que espera teste mostrava só
	// o resíduo, e ninguém testava porque ninguém via.
	//
	// As palavras que o GitHub trata como fechamento estão em
	// docs.github.com/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue
	// e nenhuma delas pode aparecer no que o `pr-body` GERA. Elas seguem ACEITAS na
	// leitura do pipeline (PR antigo, PR escrito à mão) — o que se proíbe é gerar.
	for plataforma, forma := range linkSyntax {
		palavra := strings.ToLower(strings.Fields(forma)[0])
		for _, fecha := range []string{
			"close", "closes", "closed", "fix", "fixes", "fixed",
			"resolve", "resolves", "resolved",
		} {
			if palavra == fecha {
				t.Errorf("a sintaxe de %q gera %q, que FECHA a issue no merge — o card "+
					"precisa ficar aberto em `ready-to-test` porque a esteira segue no CD",
					plataforma, forma)
			}
		}
	}
}

// `--cards` aceita as formas que uma pessoa escreve: com `#`, sem, com espaço.
// Recusar "#44" por causa do sustenido seria atrito sem razão — é como o card aparece
// em todo lugar do GitHub.
func TestCardsPedidosAceitaAsFormasQueSeEscreve(t *testing.T) {
	cfg := &config.Config{}
	for _, entrada := range []string{"44", "#44", " 44 ", "#44 "} {
		got := requestedCards(entrada, cfg)
		if len(got) != 1 || got[0] != "44" {
			t.Errorf("%q deveria virar [44], veio %v", entrada, got)
		}
	}
	// Vários de uma vez: o trabalho fecha o card E os achados que nasceram sob ele.
	if got := requestedCards("44, #49,50", cfg); len(got) != 3 {
		t.Errorf("três cards deveriam virar três entradas, veio %v", got)
	}
	// Vazio não inventa card: sem `--cards` e sem agente, quem chama recebe erro em vez
	// de um PR que não fecha nada.
	if got := requestedCards("  ", cfg); len(got) != 0 {
		t.Errorf("entrada vazia não pode inventar card, veio %v", got)
	}
}

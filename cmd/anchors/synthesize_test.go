package main

import (
	"strings"
	"testing"
)

// O CARD DE SÍNTESE NÃO ESCOLHE LADO, e é isso que o distingue de "resolva o conflito".
//
// Conflito de conteúdo é duas pessoas escrevendo coisas diferentes sobre o mesmo lugar.
// Escolher um lado por automação é escolher sem ler o outro — e num caso real medido (os
// PRs #693 e #556 do projeto de referência), os dois lados estavam CERTOS e eram sobre
// coisas diferentes: um documentava que o corpo da regra contradizia a decisão, o outro que
// o marcador de dispensa fazia o gate responder INDETERMINADO. As duas coisas entravam.
//
// O QUE O CARD PEDE é o que nenhum dos dois PRs entrega sozinho.
func TestCardDeSinteseNaoEscolheLado(t *testing.T) {
	corpo := corpoDaSintese("693", "556", "666", "198", "fix(DSHBR): a B01", "feat: o índice", "x.spec.md")

	// OS DOIS PRs e os dois cards aparecem: sem eles o card não é rastreável para trás.
	for _, ref := range []string{"#693", "#556", "#666", "#198"} {
		if !strings.Contains(corpo, ref) {
			t.Errorf("o card não cita %s — a rastreabilidade se perde quando os PRs fecham", ref)
		}
	}
	// E O QUE ELE PEDE não pode ser "escolha um".
	if !strings.Contains(corpo, "melhor de cada um") {
		t.Error("o card não pede a síntese — se pedisse escolha, a automação já teria escolhido")
	}
	// O AVISO sobre descartar sem dizer por quê: é o modo de falha desta tarefa.
	if !strings.Contains(corpo, "sem dizer por quê") {
		t.Error("o card não avisa contra descartar um lado em silêncio — o trabalho dos " +
			"dois está fechado, e o que se perder ninguém vai saber o que era")
	}
}

// COM UM LADO SÓ o card nasce assim mesmo, e diz que falta.
//
// Quando o conflito é contra o branch de integração, o "outro lado" é trabalho já mesclado.
// Achá-lo exige ler o histórico do arquivo, e adivinhar erraria — um card com um lado só
// ainda é melhor que um PR parado sem dono.
func TestCardDeSinteseComUmLadoSo(t *testing.T) {
	corpo := corpoDaSintese("693", "", "666", "", "fix(DSHBR): a B01", "", "x.spec.md")

	if strings.Contains(corpo, "| # |") || strings.Contains(corpo, "#  ") {
		t.Error("o card cita um PR vazio — o segundo lado não existe, e a tabela mente")
	}
	// DIZER QUE FALTA é o que torna o card acionável: quem o pega sabe que tem de
	// descobrir o outro lado, e onde procurar.
	if !strings.Contains(corpo, "não foi identificado") {
		t.Error("o card não diz que o outro lado falta — quem o pegar vai procurar um PR " +
			"que não existe")
	}
	if !strings.Contains(corpo, "git log") {
		t.Error("o card não diz ONDE procurar o outro lado — a instrução sem o caminho " +
			"transfere o trabalho de descobrir para quem já foi interrompido")
	}
	// E A INSTRUÇÃO no singular: "leia os dois PRs fechados" seria falso com um só.
	if strings.Contains(corpo, "Leia os dois PRs fechados") {
		t.Error("o card manda ler DOIS PRs quando só um foi fechado")
	}
}

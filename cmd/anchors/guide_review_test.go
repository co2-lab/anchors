package main

import (
	"strings"
	"testing"
)

// O guia de review dizia como julgar e nunca dizia QUEM MOVE O CARD. O pipeline sabe a
// regra — `anchors-pr-checks.yml` só escreve `ready-to-test` depois do merge —, e o
// revisor não tinha onde lê-la.
//
// Medido: duas revisões independentes rodaram em paralelo sobre o mesmo PR. A primeira
// aprovou e moveu o card para `ready-to-test` à mão, com o PR aberto. A segunda achou um
// defeito real que a primeira não cobriu e NÃO mexeu no estado — "já está `ready-to-test`
// pela outra revisão". A segunda agiu certo; a primeira criou o fato que a travou.
func TestGuiaDeReview_dizQueORevisorNaoMoveOCard(t *testing.T) {
	if !strings.Contains(reviewGuide, "NÃO MOVE O CARD") {
		t.Error("o guia de review deveria dizer explicitamente que o revisor não move o card")
	}
}

// A regra é sobre o FATO que move cada estado. Sem isso, "não mova" vira proibição sem
// razão — o tipo de instrução que se descumpre na primeira vez que parece atrapalhar.
func TestGuiaDeReview_dizQualFatoMoveCadaEstado(t *testing.T) {
	if !strings.Contains(reviewGuide, "MERGEADO") {
		t.Error("o guia deveria dizer que é o MERGE que leva a `ready-to-test`")
	}
	// A quebra de linha do texto cai no meio de "checks passam" — a asserção tem de
	// normalizar o espaço em branco, senão ela testa a largura da coluna, não a régua.
	achatado := strings.Join(strings.Fields(reviewGuide), " ")
	if !strings.Contains(achatado, "checks passam") {
		t.Error("o guia deveria dizer que é o check que leva a `ready-to-review`")
	}
}

// A consequência, e não só a regra: `ready-to-test` é o fim da alçada do Anchors. Um card
// ali com PR aberto diz que o trabalho entrou quando ele não entrou, e nenhum gate o
// confronta mais.
func TestGuiaDeReview_dizPorQueOEstadoErradoCustaCaro(t *testing.T) {
	if !strings.Contains(reviewGuide, "fim da alçada") {
		t.Error("o guia deveria dizer que `ready-to-test` sai do alcance dos gates")
	}
	// E a assimetria que decide o comportamento diante da dúvida: atrasado se corrige
	// sozinho no próximo evento do pipeline; errado, não.
	if !strings.Contains(reviewGuide, "mais caro que o estado atrasado") {
		t.Error("o guia deveria dizer qual dos dois erros é o pior")
	}
}

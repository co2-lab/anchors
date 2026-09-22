package governance

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
	if !strings.Contains(reviewGuide, "YOU DO NOT MOVE THE CARD") {
		t.Error("o guia de review deveria dizer explicitamente que o revisor não move o card")
	}
}

// A regra é sobre o FATO que move cada estado. Sem isso, "não mova" vira proibição sem
// razão — o tipo de instrução que se descumpre na primeira vez que parece atrapalhar.
func TestGuiaDeReview_dizQualFatoMoveCadaEstado(t *testing.T) {
	if !strings.Contains(reviewGuide, "MERGED") {
		t.Error("o guia deveria dizer que é o MERGE que leva a `ready-to-test`")
	}
	// A quebra de linha do texto cai no meio de "checks passam" — a asserção tem de
	// normalizar o espaço em branco, senão ela testa a largura da coluna, não a régua.
	achatado := strings.Join(strings.Fields(reviewGuide), " ")
	if !strings.Contains(achatado, "when the checks") {
		t.Error("o guia deveria dizer que é o check que leva a `ready-to-review`")
	}
}

// A consequência, e não só a regra: `ready-to-test` é o fim da alçada do Anchors. Um card
// ali com PR aberto diz que o trabalho entrou quando ele não entrou, e nenhum gate o
// confronta mais.
func TestGuiaDeReview_dizPorQueOEstadoErradoCustaCaro(t *testing.T) {
	if !strings.Contains(reviewGuide, "end of Anchors' jurisdiction") {
		t.Error("o guia deveria dizer que `ready-to-test` sai do alcance dos gates")
	}
	// E a assimetria que decide o comportamento diante da dúvida: atrasado se corrige
	// sozinho no próximo evento do pipeline; errado, não.
	if !strings.Contains(reviewGuide, "more expensive than the late state") {
		t.Error("o guia deveria dizer qual dos dois erros é o pior")
	}
}

// A LISTA existe para que nada passe por ESQUECIMENTO.
//
// O review é o lugar do que exige julgamento, e julgamento não vira gate — mas a memória
// de QUAIS julgamentos fazer não deveria depender de lembrar. É a mesma inversão que os
// pontos de conformidade já fazem nos outros guias: em vez de "leia a prosa e lembre do
// que importa", a lista diz "estes são os itens, um a um".
//
// O guia de review era o único que governa e não tinha a seção — o `guide-checklist`
// cobra `## Pontos de conformidade` de toda régua, e esta escapava por viver no binário.
func TestGuiaDeReview_temPontosDeConformidade(t *testing.T) {
	if !strings.Contains(reviewGuide, "## Pontos de conformidade") {
		t.Fatal("o guia de review não tem a seção obrigatória de pontos de conformidade")
	}
}

// Cada ponto tem CÓDIGO, e a numeração é contínua: o relatório precisa referenciar o item
// específico, e um buraco na sequência é item apagado sem ninguém notar.
func TestGuiaDeReview_pontosNumeradosSemBuraco(t *testing.T) {
	for i := 1; i <= 13; i++ {
		codigo := "REV-CK" + itoa(i) + ":"
		if !strings.Contains(reviewGuide, codigo) {
			t.Errorf("falta o ponto %s — a numeração tem buraco", codigo)
		}
	}
}

// ANCORADO NA PROSA: cada ponto destila uma régua que o corpo do guia já explica. Um
// ponto sem explicação acima é invenção — e foi assim que três deles nasceram nesta
// mesma escrita, antes de ganharem a prosa que lhes faltava.
func TestGuiaDeReview_cadaPontoTemProsaAcima(t *testing.T) {
	corpo := reviewGuide[:strings.Index(reviewGuide, "## Pontos de conformidade")]
	// Os termos que cada ponto afirma têm de aparecer ANTES da lista.
	ancoras := map[string]string{
		"REV-CK1":  "do the checks EXIST?",
		"REV-CK3":  "decide what it needed to decide",
		"REV-CK4":  "realize the rule, or only cite it",
		"REV-CK5":  "@TBD",
		"REV-CK6":  "contradict each other",
		"REV-CK9":  "Checked:",
		"REV-CK10": "PROVE, or only execute",
		"REV-CK11": "WHICH requirement it proves",
		"REV-CK12": "change without saying",
		"REV-CK13": "DO NOT MOVE THE CARD",
	}
	for ck, ancora := range ancoras {
		if !strings.Contains(corpo, ancora) {
			t.Errorf("%s não tem prosa acima da lista (procurei %q) — ou falta a explicação, ou o ponto é invenção", ck, ancora)
		}
	}
}

// A lista NÃO substitui os checks. Confundir as duas coisas faria o revisor gastar-se
// refazendo à mão o que o script já confrontou melhor.
func TestGuiaDeReview_listaNaoSubstituiOsChecks(t *testing.T) {
	achatado := strings.Join(strings.Fields(reviewGuide), " ")
	if !strings.Contains(achatado, "NOT a substitute for the checks") {
		t.Error("a lista deveria dizer que não substitui os checks")
	}
}

// itoa sem importar strconv só para isto.
func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

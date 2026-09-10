package gate

import "testing"

// O FORMATO COM TÍTULO DE SEÇÃO, medido no blue-eyes.
//
// O guia mostra a revisão como `> **CODIGO-R0001:** o que mudou`, e o regex exige os dois
// pontos. Mas uma spec pode registrá-la como TÍTULO de seção — que é o formato natural
// quando há mais de uma e elas ganham corpo:
//
//	## Revisões
//
//	### SRMTS-R0001 — a `B06` afirmava um vocabulário que não existe
//
// Medido: a `ServiceMetrics` tinha a `R0001` assim e a `R0002` no cabeçalho (com `:`). O
// gate contou UMA revisão e viu a maior como `-R0002` — e reprovou por "não sequencial".
//
// O diagnóstico que ele deu estava certo sobre o sintoma e errado sobre a causa: as
// revisões ERAM sequenciais; uma delas não foi vista.
func TestRevisionsOf_reconheceOTituloDeSecao(t *testing.T) {
	conteudo := "## Revisões\n\n" +
		"### SRMTS-R0001 — a `B06` afirmava um vocabulário que não existe\n\n" +
		"A regra dizia que as severidades já eram conhecidas.\n"

	revs := RevisionsOf(conteudo)

	if len(revs) != 1 {
		t.Fatalf("a revisão em título de seção não foi vista: %d encontrada(s)", len(revs))
	}
	if revs[0].Numero != 1 || revs[0].Codigo != "SRMTS" {
		t.Errorf("leu %s-R%04d", revs[0].Codigo, revs[0].Numero)
	}
}

// E os dois formatos convivem no mesmo arquivo: é o estado real de uma spec que ganhou a
// segunda revisão no cabeçalho depois de ter a primeira numa seção.
func TestRevisionsOf_osDoisFormatosConvivem(t *testing.T) {
	conteudo := "> **SRMTS-R0002:** a `B06` passou a citar o vocabulário publicado\n\n" +
		"## Revisões\n\n### SRMTS-R0001 — a `B06` afirmava um vocabulário que não existe\n"

	if revs := RevisionsOf(conteudo); len(revs) != 2 {
		t.Errorf("esperava 2 revisões, veio %d — a numeração pareceria não-sequencial", len(revs))
	}
}

// @anchors
//   ref: INCTN

package initx

// --- a instrução que todo gate de julgamento sobre PEÇA AUSENTE precisa carregar ---
//
// A spec nasce ANTES do código: é o fluxo normal do Anchors, porque a spec é a âncora.
// Para declarar isso honestamente existe o `@TBD: code,feature,test` — "estas peças estão
// decididas e ainda não foram escritas".
//
// Os gates `measures: judgment` perguntam sobre o código (ou o teste) que realiza a regra.
// Diante de `@TBD`, a pergunta não tem sujeito: não há trecho a confrontar. E a saída
// fácil é dar PASS para destravar o commit — que é o pior desfecho possível, porque o
// carimbo fica no mapa parecendo verificação real. É pior que pendência aberta: uma
// pendência diz "ninguém olhou"; um PASS afirma que alguém olhou e aprovou.
//
// Medido no app de referência: a `MutationHarness.spec.md` declarava `@TBD: code,feature,test`,
// o identificador MTHRN não existia em arquivo de código nenhum, e o `anchors check`
// BARROU pedindo o julgamento de "o trecho REALIZA o que a regra descreve?".
//
// POR QUE INSTRUIR, E NÃO FILTRAR
//
// A primeira proposta foi o Anchors não PERGUNTAR quando houvesse `@TBD`. É pior, e o
// caso que a derruba é simples: uma spec com `@TBD: test` (falta só o teste) TEM código, e
// a pergunta do `regra-cumprida` é perfeitamente válida. Um filtro por presença do
// marcador a suprimiria — trocaria um julgamento impossível por uma dispensa cega, e o
// gate deixaria de cobrar o que devia.
//
// O `Requires:` também não serve: ele restringe o gate aos alvos que CONTÊM o texto, ou
// seja, faria o inverso — rodaria só onde há `@TBD`.
//
// Quem julga já tem a informação toda: o marcador está no arquivo que ele lê, e ele
// distingue "declarou que o código não existe" de "declarou que falta o teste". Essa
// leitura é o que nenhum filtro faz.

// tbdInstruction returns the instruction appended to the `ask:` of a judgment gate.
//
// `piece` is what the gate asks about, as the `ask:` already names it ("the code", "the
// test") — the text has to read as the continuation of the question, not as a notice
// pinned to the end.
//
// In English, like the rest of the `ask:` it closes: the text is written into every new
// project's anchors.yaml and read by the judge in any language. It used to be Portuguese,
// and told the judge to answer "DISPENSADO" — a word the verdict vocabulary only accepts
// as a legacy alias of `waived`.
func tbdInstruction(piece string) string {
	return " BEFORE ANSWERING, check whether this unit declares `@TBD` for the piece " +
		"this question asks about (`@TBD: code`, `@TBD: code,test`, …). If it does, " +
		piece + " does NOT exist yet by a recorded decision, and there is nothing to " +
		"confront: answer `waived`, naming the absence (\"the spec declares @TBD for " +
		"this piece and it does not exist in the repository\"). Do NOT answer `pass` — " +
		"`pass` is a statement ABOUT " + piece + ", and without " + piece + " it states " +
		"what nobody verified; the stamp stays in the map looking like real verification. " +
		"Also check that the `@TBD` is TRUE: if the piece already exists in the " +
		"repository, the marker is stale, the judgment is due as usual, and the stale " +
		"declaration is a finding — while it stays there, every gate that reads it " +
		"waives what it should demand."
}

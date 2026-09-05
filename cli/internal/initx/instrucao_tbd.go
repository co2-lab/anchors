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
// Medido no blue-eyes (#76): a `MutationHarness.spec.md` declarava `@TBD: code,feature,test`,
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

// instrucaoTBD devolve a instrução a acrescentar ao `ask:` de um gate de julgamento.
//
// `peca` é o que o gate interroga, na forma como o `ask:` já fala dele ("o código", "o
// teste") — o texto tem de soar como continuação da pergunta, não como aviso pregado ao
// fim.
func instrucaoTBD(peca string) string {
	return " ANTES DE RESPONDER, confira se esta unidade declara `@TBD` para a peça que " +
		"esta pergunta interroga (`@TBD: code`, `@TBD: code,test`, …). Se declarar, " +
		peca + " ainda NÃO existe por decisão registrada, e não há o que confrontar: " +
		"responda DISPENSADO, nomeando a ausência (\"a spec declara @TBD para esta peça " +
		"e ela não existe no repositório\"). NÃO responda `pass` — o `pass` é uma " +
		"afirmação SOBRE " + peca + ", e sem " + peca + " ele afirma o que ninguém " +
		"verificou; o carimbo fica no mapa parecendo verificação real. " +
		"Confira também se o `@TBD` é VERDADE: se a peça já existe no repositório, o " +
		"marcador está desatualizado, o julgamento é devido normalmente, e a declaração " +
		"obsoleta é um achado — enquanto ela estiver lá, todo gate que a lê dispensa o " +
		"que devia cobrar."
}

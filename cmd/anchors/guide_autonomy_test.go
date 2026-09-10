package main

import (
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/settings"
)

// O GUIA MUDA conforme a declaração local, e a diferença é o ponto.
//
// O `settings user-issues` fecha a porta do claim: quem não decide o produto não recebe
// card escalonado. Mas a porta que mais se usa é outra, e não tem label — é o agente
// PERGUNTAR ao dev que o está rodando.
//
// A pergunta parece inofensiva e não é: o dev conhece o código e vai responder, a resposta
// é razoável, e vira decisão de produto tomada por quem não tinha autoridade — sem passar
// pelo plano, sem revisão, e sem rastro de que foi decidido ali.
func TestAutonomyGuide_mudaComADeclaracao(t *testing.T) {
	semAutoridade := t.TempDir()
	if err := settings.Save(semAutoridade, settings.Settings{
		UserIssues: settings.Bool(false),
	}); err != nil {
		t.Fatal(err)
	}
	comAutoridade := t.TempDir()
	if err := settings.Save(comAutoridade, settings.Settings{
		UserIssues: settings.Bool(true),
	}); err != nil {
		t.Fatal(err)
	}

	texto := autonomyGuide(semAutoridade)
	if !strings.Contains(texto, "Não pergunte a quem está rodando você") {
		t.Errorf("quem não decide o produto precisa ler a proibição explícita:\n%s", texto)
	}
	if !strings.Contains(texto, "--for-user") {
		t.Error("a instrução não diz o que fazer no lugar de perguntar")
	}
	// A razão importa mais que a proibição: uma regra sem porquê é a primeira a ser
	// contornada quando atrapalha.
	if !strings.Contains(texto, "autoridade") {
		t.Error("a instrução proíbe sem dizer por quê")
	}

	outro := autonomyGuide(comAutoridade)
	if strings.Contains(outro, "Não pergunte a quem está rodando você") {
		t.Error("quem declarou que decide o produto não recebe a proibição")
	}
	if !strings.Contains(outro, "se escreve, não se") {
		t.Error("mesmo quem decide precisa ler que a decisão fica REGISTRADA")
	}
}

// SEM DECLARAÇÃO, a régua é a fechada. O padrão é o mesmo do claim, e pela mesma razão: o
// custo de errar para o lado aberto é alguém decidir o produto sem autoridade, e isso é
// invisível depois do fato.
func TestAutonomyGuide_semDeclaracaoEhFechado(t *testing.T) {
	texto := autonomyGuide(t.TempDir())
	if !strings.Contains(texto, "Não pergunte a quem está rodando você") {
		t.Errorf("sem declaração, o padrão tem de ser o fechado:\n%s", texto)
	}
	// E o texto não pode AFIRMAR uma declaração que não houve: quem nunca declarou iria
	// procurar em `.anchors/settings.yaml` o registro que não está lá.
	if strings.Contains(texto, "ficou declarado") {
		t.Error("o texto afirma uma declaração que não aconteceu")
	}
	if !strings.Contains(texto, "não declarou") {
		t.Error("o texto não diz que a decisão está em aberto")
	}
}

// `decidesProduct` erra para o lado FECHADO quando não consegue ler.
//
// Um erro de leitura que resultasse em "pode perguntar" seria a falha silenciosa mais cara
// deste mecanismo: o agente perguntaria, e ninguém saberia que a régua não foi aplicada.
func TestDecidesProduct_erraParaOLadoFechado(t *testing.T) {
	if decidesProduct(t.TempDir()) {
		t.Error("sem arquivo, o agente não decide o produto")
	}
	if decidesProduct("/caminho/que/nao/existe") {
		t.Error("com erro de leitura, o agente não decide o produto")
	}
}

// `/dev/null` NÃO é terminal, e o `ModeCharDevice` sozinho não sabe disso.
//
// Foi o defeito que impediu um agente de trabalhar: `/dev/null` também é char device, então
// `anchors next < /dev/null` — o caso normal de execução não-interativa — passava pela
// guarda, caía na pergunta do perfil e morria com `erro: ler a resposta: EOF`.
//
// Um comando que morre por EOF não diz o que fazer: parece defeito da ferramenta, e quem lê
// vai procurar o problema no ambiente.
func TestInteractiveTerminal_devNullNaoEhTerminal(t *testing.T) {
	nul, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer nul.Close()

	original := os.Stdin
	os.Stdin = nul
	defer func() { os.Stdin = original }()

	if interactiveTerminal() {
		t.Error("`/dev/null` foi tratado como terminal — o agente cairia na pergunta e " +
			"morreria por EOF")
	}
}

// E um arquivo comum também não: `anchors next < resposta.txt` não é alguém do outro lado.
func TestInteractiveTerminal_arquivoNaoEhTerminal(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	original := os.Stdin
	os.Stdin = f
	defer func() { os.Stdin = original }()

	if interactiveTerminal() {
		t.Error("um arquivo comum foi tratado como terminal")
	}
}

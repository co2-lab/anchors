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

// UM DEV NOVO pediu ao agente dele para contribuir. O agente leu o CONTRIBUTING, montou o
// plano de onboarding CORRETO — instalar, `pnpm install`, declarar perfil, `doctor --fix`,
// `guide work` — e parou no passo 4 pedindo autorização: "o passo 4 altera o repositório
// remoto (branch protection, labels) (...) quero teu OK explícito".
//
// A cautela estava certa. O que faltava era saber que `doctor --fix` é IDEMPOTENTE: num
// repositório já montado ele não muda nada. Sem isso, o agente escolhe entre dois erros —
// pedir autorização para tudo (e não começar) ou não pedir para nada.
//
// Ele também temeu que `settings role` fosse interativo. É, sem argumento — e aceita o
// perfil e a data como argumento, o que nenhum texto dizia.
func TestAutonomyGuide_preparacaoNaoPedeAutorizacao(t *testing.T) {
	dir := t.TempDir()
	g := autonomyGuide(dir)

	if !strings.Contains(g, "não pede autorização") {
		t.Error("o guia deveria dizer que preparar o ambiente não pede autorização")
	}
	if !strings.Contains(g, "idempotente") && !strings.Contains(g, "IDEMPOTENTE") {
		t.Error("a razão tem de estar dita: os comandos conferem antes de agir")
	}
	if !strings.Contains(g, "doctor --fix") {
		t.Error("o comando que causou a parada tem de ser nomeado")
	}
	// A régua tem de ser NOMEADA, senão o agente generaliza errado — "tocar o remoto"
	// incluiria `git push`, e ele voltaria a pedir OK para entregar trabalho.
	if !strings.Contains(g, "REVERSIBILIDADE") {
		t.Error("o guia deveria dizer QUAL é a régua, não só listar exceções")
	}
	// E o contra-exemplo: sem ele, "preparação não pede autorização" lê-se como
	// "nada pede autorização".
	if !strings.Contains(g, "PEDE autorização") {
		t.Error("o guia deveria dizer o que AINDA pede autorização")
	}
}

func TestAutonomyGuide_dizQueSettingsRoleAceitaArgumento(t *testing.T) {
	// O agente parou também por isto: "o passo 3 pode ser interativo (...) eu paro e
	// devolvo o comando pra você rodar". O comando aceita perfil e data como argumento, e
	// nenhum texto dizia — ele teria de rodar `--help` para descobrir.
	g := autonomyGuide(t.TempDir())
	if !strings.Contains(g, "--date") {
		t.Error("o guia deveria mostrar a forma NÃO interativa do `settings role`")
	}
	if !strings.Contains(g, "espera") && !strings.Contains(g, "esperando") {
		t.Error("a consequência tem de estar dita: um agente sem terminal fica esperando")
	}
}

// A seção vale para TODO perfil, e isso não é detalhe: a primeira versão a colocou depois
// da bifurcação que separa quem decide produto de quem não decide, e o `architect` — que
// retorna cedo — nunca a via. Quem monta o ambiente não é só o dev.
func TestAutonomyGuide_preparacaoValeParaTodoPerfil(t *testing.T) {
	for _, perfil := range []string{"dev", "architect", "product-owner", "qa", "reviewer"} {
		dir := t.TempDir()
		s := settings.Settings{Role: settings.Role(perfil), DecidedAt: "2026-09-10"}
		if err := settings.Save(dir, s); err != nil {
			t.Fatalf("salvar settings de %s: %v", perfil, err)
		}
		if g := autonomyGuide(dir); !strings.Contains(g, "não pede autorização") {
			t.Errorf("perfil %s não vê a seção de preparação", perfil)
		}
	}
}

package main

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// A ISSUE tem de responder três coisas, e a terceira é a que costuma faltar: como
// destravar. Sem ela o card fica parado esperando alguém adivinhar o protocolo.
func TestEscalada_dizPorQueParouEComoDestravar(t *testing.T) {
	corpo := corpoDaEscalada("A spec pede cache; o plano diz que não haveria cache.",
		"plans/0001-fundacao.md", "12", true)

	for _, exigido := range []string{
		"A spec pede cache",         // o motivo, com as palavras de quem viu
		"plans/0001-fundacao.md",    // onde
		"R0001",                     // como registrar a decisão
		initx.LabelPrecisaDoUsuario, // o que remover para destravar
		"#12",                       // onde o trabalho parou
	} {
		if !strings.Contains(corpo, exigido) {
			t.Errorf("o corpo da escalada deve conter %q; veio:\n%s", exigido, corpo)
		}
	}
}

// Sem `--card` o comando ainda serve: nem toda incoerência é achada com um card na mão
// (o revisor lendo um PR, por exemplo). O corpo não pode citar um card que não existe.
func TestEscalada_semCardNaoInventaReferencia(t *testing.T) {
	corpo := corpoDaEscalada("O plano contradiz a si mesmo entre F02 e F04.", "", "", true)
	if strings.Contains(corpo, "#") && strings.Contains(corpo, "Trabalho parado") {
		t.Errorf("sem --card não pode citar card; veio:\n%s", corpo)
	}
	if !strings.Contains(corpo, "contradiz a si mesmo") {
		t.Error("o motivo tem de sobreviver mesmo sem card")
	}
}

// O TÍTULO da issue é uma linha. Um motivo longo não pode virar um título ilegível na
// lista de issues, que é onde o usuário vai encontrá-lo.
func TestEscalada_tituloCabeEmUmaLinha(t *testing.T) {
	longo := strings.Repeat("uma explicação bem detalhada da incoerência ", 5)
	got := primeiraLinhaDoMotivo(longo)
	if len(got) > 70 {
		t.Errorf("título com %d chars; deveria caber em 70: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("o corte deve sinalizar que há mais: %q", got)
	}
	// Motivo de várias linhas: o título é a primeira.
	if got := primeiraLinhaDoMotivo("primeira linha\nsegunda linha"); got != "primeira linha" {
		t.Errorf("o título é a primeira linha, veio %q", got)
	}
}

// A SAÍDA PADRÃO é card comum, e ela tem de ser visivelmente diferente da decisão: se as
// duas issues lessem igual, quem abre a lista não saberia qual espera por ele.
func TestEscalada_cardComumNaoPedeDecisao(t *testing.T) {
	corpo := corpoDaEscalada("O plano não cobre configuração e execução de migrations.",
		"plans/0001-fundacao.md", "12", false)

	if strings.Contains(corpo, initx.LabelPrecisaDoUsuario) {
		t.Errorf("card comum não pode mandar remover a label de decisão; veio:\n%s", corpo)
	}
	if strings.Contains(corpo, "Trabalho parado") {
		t.Errorf("card comum NÃO para o trabalho; veio:\n%s", corpo)
	}
	// E tem de ensinar a saída de emergência: quem for mexer pode descobrir, ao mexer,
	// que a mudança era maior do que quem abriu julgou.
	if !strings.Contains(corpo, "--para-usuario") {
		t.Errorf("o card comum deve dizer o que fazer se a mudança se revelar de direção; "+
			"veio:\n%s", corpo)
	}
	if !strings.Contains(corpo, "migrations") {
		t.Error("o motivo tem de sobreviver na saída padrão")
	}
}

// O caso que motivou a revisão do desenho: o gatilho não é só incoerência. Uma LACUNA
// (plano coerente, mas incompleto) segue o mesmo fluxo, e o texto não pode sugerir que
// só serve para contradição — um agente que lesse isso concluiria que o mecanismo não é
// para o caso dele.
func TestEscalada_naoPressupoeIncoerencia(t *testing.T) {
	for _, paraUsuario := range []bool{true, false} {
		corpo := corpoDaEscalada("O plano não previu migrations.", "plans/0001.md", "", paraUsuario)
		if strings.Contains(corpo, "correção mudaria") {
			t.Errorf("o corpo não pode pressupor que houve erro a corrigir (para-usuario=%v):\n%s",
				paraUsuario, corpo)
		}
	}
}

// O VÍNCULO com o trabalho de origem é LABEL, não texto no corpo.
//
// A primeira versão escrevia "Descoberto durante o card #44" na descrição, e isso não se
// consulta: não dá para listar o que pende sob um card, nem para o board desenhar a
// relação, nem para saber o que precisa entrar no mesmo PR.
//
// Também não é a sub-issue nativa do GitHub: ela aceita UM nível, e a hierarquia deste
// fluxo tem mais (plano → fase → spec → achado).
func TestVinculoComOTrabalhoDeOrigemEhLabel(t *testing.T) {
	if got := initx.LabelSob("44"); got != "anchors:sob-44" {
		t.Fatalf("a label liga o achado ao card; veio %q", got)
	}
	// E o corpo NOMEIA a relação, para quem lê a issue saber que ela não é solta.
	corpo := corpoDaEscalada("o jest mede só src/", "jest.config.js", "44", false)
	if !strings.Contains(corpo, "anchors:sob-44") {
		t.Errorf("o corpo deve citar a label, senão quem lê não sabe como achar as irmãs;\n%s", corpo)
	}
	if !strings.Contains(corpo, "mesmo PR") {
		t.Error("o corpo deve dizer que os dois se entregam juntos — é a razão de o achado " +
			"nascer preso em vez de solto na fila")
	}
}

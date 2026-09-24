package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// A ISSUE tem de responder três coisas, e a terceira é a que costuma faltar: como
// destravar. Sem ela o card fica parado esperando alguém adivinhar o protocolo.
func TestEscalada_dizPorQueParouEComoDestravar(t *testing.T) {
	corpo := escalationBody("A spec pede cache; o plano diz que não haveria cache.",
		"plans/0001-fundacao.md", "12", true, false)

	for _, exigido := range []string{
		"A spec pede cache",      // o motivo, com as palavras de quem viu
		"plans/0001-fundacao.md", // onde
		"R0001",                  // como registrar a decisão
		initx.LabelNeedsUser,     // o que remover para destravar
		"#12",                    // onde o trabalho parou
	} {
		if !strings.Contains(corpo, exigido) {
			t.Errorf("o corpo da escalada deve conter %q; veio:\n%s", exigido, corpo)
		}
	}
}

// Sem `--card` o comando ainda serve: nem toda incoerência é achada com um card na mão
// (o revisor lendo um PR, por exemplo). O corpo não pode citar um card que não existe.
func TestEscalada_semCardNaoInventaReferencia(t *testing.T) {
	corpo := escalationBody("O plano contradiz a si mesmo entre F02 e F04.", "", "", true, false)
	if strings.Contains(corpo, "#") && strings.Contains(corpo, "Work stopped") {
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
	got := firstLineOfReason(longo)
	if len(got) > 70 {
		t.Errorf("título com %d chars; deveria caber em 70: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("o corte deve sinalizar que há mais: %q", got)
	}
	// Motivo de várias linhas: o título é a primeira.
	if got := firstLineOfReason("primeira linha\nsegunda linha"); got != "primeira linha" {
		t.Errorf("o título é a primeira linha, veio %q", got)
	}
}

// A SAÍDA PADRÃO é card comum, e ela tem de ser visivelmente diferente da decisão: se as
// duas issues lessem igual, quem abre a lista não saberia qual espera por ele.
func TestEscalada_cardComumNaoPedeDecisao(t *testing.T) {
	corpo := escalationBody("O plano não cobre configuração e execução de migrations.",
		"plans/0001-fundacao.md", "12", false, false)

	if strings.Contains(corpo, initx.LabelNeedsUser) {
		t.Errorf("card comum não pode mandar remover a label de decisão; veio:\n%s", corpo)
	}
	if strings.Contains(corpo, "Work stopped") {
		t.Errorf("card comum NÃO para o trabalho; veio:\n%s", corpo)
	}
	// E tem de ensinar a saída de emergência: quem for mexer pode descobrir, ao mexer,
	// que a mudança era maior do que quem abriu julgou.
	if !strings.Contains(corpo, "--for-user") {
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
		corpo := escalationBody("O plano não previu migrations.", "plans/0001.md", "", paraUsuario, false)
		if strings.Contains(corpo, "correction would change") {
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
	if got := initx.LabelSob("44"); got != "anchors:under-44" {
		t.Fatalf("a label liga o achado ao card; veio %q", got)
	}
	// E o corpo NOMEIA a relação, para quem lê a issue saber que ela não é solta.
	corpo := escalationBody("o jest mede só src/", "jest.config.js", "44", false, false)
	if !strings.Contains(corpo, "anchors:under-44") {
		t.Errorf("o corpo deve citar a label, senão quem lê não sabe como achar as irmãs;\n%s", corpo)
	}
	if !strings.Contains(corpo, "same PR") {
		t.Error("o corpo deve dizer que os dois se entregam juntos — é a razão de o achado " +
			"nascer preso em vez de solto na fila")
	}
}

// O VÍNCULO É COM O CARD, E O PR É O CAMINHO — não um segundo eixo.
//
// Um achado que nasce revisando trabalho alheio não tem `anchors-owner`: o agente não pegou
// card nenhum, está lendo o de outro. Sem procedência, ninguém consegue perguntar "o que
// gerou esta decisão?".
//
// MEDIDO no projeto de referência: os cards #786 e #788 nasceram assim, citando o PR
// revisado em PROSA — exatamente o que a doutrina do `under-<n>` condena, porque frase no
// corpo não se consulta.
//
// POR QUE NÃO UMA LABEL DE PR: o PR já declara a issue dele (`Refs`/`Closes`), e guardar as
// duas pontas duplicaria o que a plataforma relaciona — quem abre o PR vê a issue. O que
// faltava era o comando LER essa declaração em vez de o agente procurar.
//
// Conferido nos dois casos reais: o PR #556 declara `Closes #198`, e o #783 declara
// `Refs #735` — o card cuja decisão gerou a contradição que o #788 escalou.
func TestVinculoDoPRSaiDoCorpoDele(t *testing.T) {
	for _, caso := range []struct {
		nome, corpo, quer string
	}{
		{"o que o pr-body gera", "texto\n\nRefs #735\n", "735"},
		{"a palavra de fechamento", "Closes #198", "198"},
		{"maiúscula ou minúscula", "closes #42", "42"},
		{"fixes também", "Fixes #7", "7"},
		// SÓ NO INÍCIO DA LINHA: um número citado na prosa não é vínculo, e foi o defeito
		// que o `pr-checks` já corrigiu uma vez — "o #12 não é fechado por este PR".
		{"citado na prosa não conta", "o card refs #99 é de outro", ""},
		{"sem vínculo nenhum", "só descrição", ""},
	} {
		m := vinculoNoCorpoRE.FindStringSubmatch(caso.corpo)
		got := ""
		if m != nil {
			got = m[1]
		}
		if got != caso.quer {
			t.Errorf("%s: de %q esperava %q, veio %q", caso.nome, caso.corpo, caso.quer, got)
		}
	}
}

// OS DOIS VÍNCULOS COEXISTEM, e não se substituem.
//
//	anchors:under-198     o achado pertence ao trabalho do card #198
//	anchors:from-pr-556   foi lendo o PR #556 que alguém o viu
//
// O card #647 do projeto de referência tem os dois na história e só um virou label: a prosa
// diz "ao revisar o PR #556" e a label diz `under-198` — o card que aquele PR fecha.
// Guardar só o card perde por onde o achado apareceu; guardar só o PR perde onde ele se
// entrega.
//
// O card É DERIVADO do PR (pelo `Refs`/`Closes`), e isso não torna a segunda label
// redundante: derivar é ler uma vez. Se o PR for reescrito depois, a declaração muda e o
// vínculo histórico se perde — a label registra o que foi lido, quando foi lido.
func TestEscalateGuardaOCardEOPRJuntos(t *testing.T) {
	fonte := leFonte(t, "escalate.go")

	// AS DUAS LABELS no mesmo caminho de criação, não uma OU outra.
	for _, peca := range []string{"initx.LabelDePR(", "initx.LabelSob(card)"} {
		if !strings.Contains(fonte, peca) {
			t.Errorf("o `escalate` não aplica %q — um dos dois vínculos se perde", peca)
		}
	}

	// E NENHUMA DENTRO DO `else` DA OUTRA: se o PR excluísse o card, revisar um PR
	// deixaria o achado sem saber por onde se entrega.
	iPR := strings.Index(fonte, `if revisandoPR != "" {`)
	iCard := strings.Index(fonte, `if card != "" {
				labels = append(labels, initx.LabelSob(card))`)
	if iPR < 0 || iCard < 0 {
		t.Fatal("não achei os dois blocos de vínculo — o comando mudou de forma")
	}
	if iPR > iCard {
		t.Error("o bloco do PR vem depois do bloco do card: como o card é DERIVADO do PR, " +
			"a derivação precisa acontecer antes de o card ser lido")
	}
}

// A bug is not a decision (blue-eyes, 2026-09-24: six "decisions" in two hours, all bugs).
// The body says so first, never asks for one, and names the card only as the flag says.
func TestBugBody_isNotADecision(t *testing.T) {
	b := bugBody("the review job picks the wrong card", ".github/workflows/anchors-pr-checks.yml", "966", true)
	for _, want := range []string{"nothing to decide", "anchors-pr-checks.yml", "Card #966 waits for it", "close this issue"} {
		if !strings.Contains(b, want) {
			t.Errorf("the bug body should say %q:\n%s", want, b)
		}
	}
	for _, bad := range []string{"not the agent's", "needs-user", "decide, and record"} {
		if strings.Contains(b, bad) {
			t.Errorf("a bug body must not read as a decision (%q):\n%s", bad, b)
		}
	}
	if nb := bugBody("x", "", "966", false); strings.Contains(nb, "waits for it") || !strings.Contains(nb, "which goes on") {
		t.Errorf("without --blocking the card goes on, and the body must say so:\n%s", nb)
	}
}

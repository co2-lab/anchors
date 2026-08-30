package main

import (
	"strings"
	"testing"
)

func TestAssuntoAceitaOFormatoQueOChangelogLe(t *testing.T) {
	for _, ok := range []string{
		"feat: a coisa nova",
		"fix(board): a contagem do refresh",
		"feat(gate)!: `--changed` passa a exigir caminho relativo",
		"chore: sobe a versão",
	} {
		if !assuntoRE().MatchString(ok) {
			t.Errorf("deveria aceitar %q", ok)
		}
	}
}

// O CASO REAL que motivou o gate: o squash do PR entrou como "[MTUAO] Plano 0017 —
// mutação", porque o GitHub usa o TÍTULO DO PR como mensagem do squash e o título estava
// no formato do card. Esse commit — o que introduziu o plano — não apareceria no
// changelog, e ninguém notaria até a primeira versão ser gerada.
func TestAssuntoRecusaOQueSumiriaDoChangelog(t *testing.T) {
	for _, ruim := range []string{
		"[MTUAO] Plano 0017 — mutação, revisando o plano 0001",
		"ajustes",
		"feat sem os dois pontos",
		"bugfix: tipo que não existe na lista",
		"feat:sem espaço depois",
		"feat: ", // tipo certo, assunto vazio
	} {
		if assuntoRE().MatchString(ruim) {
			t.Errorf("deveria recusar %q — sumiria do changelog", ruim)
		}
	}
}

// O QUE O GIT GERA passa: ninguém escreveu essas mensagens, e barrá-las quebraria
// operações normais (merge, revert, rebase com fixup) em vez de melhorar o histórico.
func TestMensagemDoGitNaoEhBarrada(t *testing.T) {
	for _, m := range []string{
		"Merge pull request #13 from acme/feat-x",
		"Merge branch 'develop' into feat-y",
		"Revert \"feat: a coisa\"",
		"fixup! feat: a coisa",
		"squash! fix: outra",
	} {
		if !geradaPeloGit(m) {
			t.Errorf("%q é gerada pelo git e não pode barrar", m)
		}
	}
	// Mas uma mensagem HUMANA que começa parecido não escapa: "Mergeando o trabalho" não
	// é do git, e passar por causa do prefixo abriria a porta para qualquer coisa.
	if geradaPeloGit("Mergeando o trabalho da semana") {
		t.Error("só o formato exato do git escapa — senão o prefixo vira brecha")
	}
}

// O assunto é a primeira linha ÚTIL: o arquivo que o git passa vem cheio dos comentários
// dele, e ler a primeira linha crua acusaria todo commit feito pelo editor.
func TestAssuntoPulaOsComentariosDoGit(t *testing.T) {
	arquivo := "\n# Please enter the commit message for your changes.\n#\n" +
		"feat: a coisa\n\n# On branch main\n"
	if got := primeiraLinhaUtil(arquivo); got != "feat: a coisa" {
		t.Errorf("deveria achar o assunto entre os comentários, veio %q", got)
	}
	// Mensagem só de comentário (commit abortado): sem assunto, e quem barra é o git.
	if got := primeiraLinhaUtil("# tudo comentado\n#\n"); got != "" {
		t.Errorf("sem assunto deveria devolver vazio, veio %q", got)
	}
}

// O laudo tem de ENSINAR o formato, com exemplo. Um "formato inválido" seco obrigaria
// quem foi barrado a procurar a convenção em outro lugar — e o caminho barato aí é
// contornar o hook.
func TestLaudoEnsinaOFormato(t *testing.T) {
	cmd := newCommitMsgCmd()
	if !strings.Contains(cmd.Long, "changelog") {
		t.Error("o texto deve dizer POR QUE o formato importa, não só que é obrigatório")
	}
}

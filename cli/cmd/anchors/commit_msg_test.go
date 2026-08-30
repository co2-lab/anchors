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

// CONFRONTADO COM O COMMITLINT, que é a ferramenta madura desta régua. Cada caso aqui foi
// verificado nos dois: o veredito tem de bater, senão a régua embutida ensina um formato
// e o projeto com Node cobra outro.
func TestReguaBateComOCommitlint(t *testing.T) {
	casos := []struct {
		assunto string
		passa   bool
		porque  string
	}{
		{"feat: ok", true, ""},
		{"fix(board): a contagem do refresh", true, ""},
		{"feat(gate)!: muda o contrato", true, ""},
		{"Feat: maiúscula no tipo", false, "`Feat` e `feat` viram grupos separados no changelog"},
		{"feat(): escopo vazio", false, "escopo vazio parece que alguém ia dizer algo"},
		{"feat: termina em ponto.", false, "o changelog emenda o assunto a outros textos"},
		{"bugfix: tipo que não existe", false, "tipo livre vira sinônimo e desfaz o agrupamento"},
		{"[MTUAO] Plano 0017", false, "sumiria do changelog — o caso real"},
		// A SIGLA passa, e é onde esta régua DIVERGE do commitlint de propósito: o
		// `subject-case` dele reprova isto, sem distinguir sigla de frase capitalizada.
		{"feat: SBOM sai da pasta ignorada", true, "sigla legítima não é frase capitalizada"},
	}
	for _, c := range casos {
		got := problemaNoAssunto(c.assunto) == ""
		if got != c.passa {
			verbo := "deveria passar"
			if !c.passa {
				verbo = "deveria barrar"
			}
			t.Errorf("%q %s — %s (laudo: %q)", c.assunto, verbo, c.porque,
				problemaNoAssunto(c.assunto))
		}
	}
}

// O ASSUNTO LONGO é cortado na lista de commits e no changelog. O detalhe vai no corpo,
// que não tem limite — e o laudo tem de dizer isso, senão quem foi barrado só encurta e
// perde a informação.
func TestAssuntoLongoBarraEExplicaOndeCabeODetalhe(t *testing.T) {
	longo := "feat(x): " + strings.Repeat("a", LimiteDoAssunto)
	p := problemaNoAssunto(longo)
	if p == "" {
		t.Fatalf("assunto de %d caracteres deveria barrar (limite %d)", len(longo), LimiteDoAssunto)
	}
	if !strings.Contains(p, "CORPO") {
		t.Errorf("o laudo deve dizer que o detalhe cabe no corpo; veio: %s", p)
	}
	// E o que está no limite passa: barrar em cima da linha seria arbitrário.
	noLimite := "feat: " + strings.Repeat("a", LimiteDoAssunto-6)
	if p := problemaNoAssunto(noLimite); p != "" {
		t.Errorf("assunto de exatamente %d deveria passar; veio: %s", LimiteDoAssunto, p)
	}
}

// CADA DEFEITO tem laudo PRÓPRIO. Um "formato inválido" genérico obriga quem foi barrado a
// adivinhar qual das seis regras quebrou — e adivinhar três vezes é o que faz alguém
// desligar o hook.
func TestCadaDefeitoTemLaudoProprio(t *testing.T) {
	vistos := map[string]string{}
	for _, a := range []string{
		"Feat: x", "feat(): x", "feat: x.", "bugfix: x", "sem formato nenhum",
	} {
		p := problemaNoAssunto(a)
		if p == "" {
			t.Fatalf("%q deveria ter defeito", a)
		}
		if antes, repetido := vistos[p]; repetido {
			t.Errorf("%q e %q recebem o MESMO laudo (%q) — o laudo tem de dizer o que "+
				"consertar", a, antes, p)
		}
		vistos[p] = a
	}
}

// A DIREÇÃO DA DIVERGÊNCIA importa mais que o número dela.
//
// Confrontado com o commitlint, o único caso em que os dois discordam é `feat(): x`:
// ele aceita escopo vazio, esta régua barra. Ser MAIS estrita é seguro — o que passa aqui
// passa lá —, então um projeto que troque esta régua pelo commitlint não descobre um
// histórico que a ferramenta nova reprova. O contrário produziria exatamente isso.
func TestOndeDivergeEhParaOLadoMaisEstrito(t *testing.T) {
	if problemaNoAssunto("feat(): x") == "" {
		t.Error("escopo vazio deve barrar: parece que alguém ia dizer algo e parou")
	}
	// E as formas que o commitlint aceita continuam aceitas aqui — divergir para o lado
	// FROUXO seria o problema.
	for _, ok := range []string{"feat: x", "feat(a): x", "feat(a)!: x", "feat!: x"} {
		if p := problemaNoAssunto(ok); p != "" {
			t.Errorf("%q passa no commitlint e tem de passar aqui; veio: %s", ok, p)
		}
	}
}

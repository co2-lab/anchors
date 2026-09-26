package ops

import (
	"bytes"
	"github.com/co2-lab/anchors/internal/i18n"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O CARD DE SÍNTESE NÃO ESCOLHE LADO, e é isso que o distingue de "resolva o conflito".
//
// Conflito de conteúdo é duas pessoas escrevendo coisas diferentes sobre o mesmo lugar.
// Escolher um lado por automação é escolher sem ler o outro — e num caso real medido (os
// PRs #693 e #556 do projeto de referência), os dois lados estavam CERTOS e eram sobre
// coisas diferentes: um documentava que o corpo da regra contradizia a decisão, o outro que
// o marcador de dispensa fazia o gate responder INDETERMINADO. As duas coisas entravam.
//
// O QUE O CARD PEDE é o que nenhum dos dois PRs entrega sozinho.
func TestCardDeSinteseNaoEscolheLado(t *testing.T) {
	corpo := corpoDaSintese("693", "556", "666", "198", "fix(DSHBR): a B01", "feat: o índice", "x.spec.md")

	// OS DOIS PRs e os dois cards aparecem: sem eles o card não é rastreável para trás.
	for _, ref := range []string{"#693", "#556", "#666", "#198"} {
		if !strings.Contains(corpo, ref) {
			t.Errorf("o card não cita %s — a rastreabilidade se perde quando os PRs fecham", ref)
		}
	}
	// E O QUE ELE PEDE não pode ser "escolha um".
	if !strings.Contains(corpo, "best of each") {
		t.Error("o card não pede a síntese — se pedisse escolha, a automação já teria escolhido")
	}
	// O AVISO sobre descartar sem dizer por quê: é o modo de falha desta tarefa.
	if !strings.Contains(corpo, "without saying why") {
		t.Error("o card não avisa contra descartar um lado em silêncio — o trabalho dos " +
			"dois está fechado, e o que se perder ninguém vai saber o que era")
	}
}

// COM UM LADO SÓ o card nasce assim mesmo, e diz que falta.
//
// Quando o conflito é contra o branch de integração, o "outro lado" é trabalho já mesclado.
// Achá-lo exige ler o histórico do arquivo, e adivinhar erraria — um card com um lado só
// ainda é melhor que um PR parado sem dono.
func TestCardDeSinteseComUmLadoSo(t *testing.T) {
	corpo := corpoDaSintese("693", "", "666", "", "fix(DSHBR): a B01", "", "x.spec.md")

	if strings.Contains(corpo, "| # |") || strings.Contains(corpo, "#  ") {
		t.Error("o card cita um PR vazio — o segundo lado não existe, e a tabela mente")
	}
	// DIZER QUE FALTA é o que torna o card acionável: quem o pega sabe que tem de
	// descobrir o outro lado, e onde procurar.
	if !strings.Contains(corpo, "was not identified") {
		t.Error("o card não diz que o outro lado falta — quem o pegar vai procurar um PR " +
			"que não existe")
	}
	if !strings.Contains(corpo, "git log") {
		t.Error("o card não diz ONDE procurar o outro lado — a instrução sem o caminho " +
			"transfere o trabalho de descobrir para quem já foi interrompido")
	}
	// E A INSTRUÇÃO no singular: "leia os dois PRs fechados" seria falso com um só.
	if strings.Contains(corpo, "Read both closed PRs") {
		t.Error("o card manda ler DOIS PRs quando só um foi fechado")
	}
}

// --- the command, against a fake `gh` (nothing reaches GitHub) ---

// synthGH answers like GitHub for PRs 693 (card #666) and 556 (card #198). `extra` is
// bash inserted before the defaults, to make one call fail.
func synthGH(t *testing.T, extra string) string {
	t.Helper()
	return fakeGH(t, extra+`
case "$*" in
"pr view 693"*"--json body"*) echo "Refs #666" ;;
"pr view 556"*"--json body"*) echo "Closes #198" ;;
"pr view 693"*"--json title"*) echo "fix(DSHBR): the B01" ;;
"pr view 556"*"--json title"*) echo "feat: the index" ;;
"issue create"*) echo "https://github.com/acme/app/issues/700" ;;
esac`)
}

func synthProject(t *testing.T, cfg string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, cfg)
	return root
}

func runSynth(t *testing.T, args ...string) (err error, out, errOut string) {
	t.Helper()
	c := newSynthesizeCmd()
	var o, e bytes.Buffer
	c.SetOut(&o)
	c.SetErr(&e)
	c.SetArgs(args)
	err = c.Execute()
	return err, o.String(), e.String()
}

func TestSynthesizeExistsOnlyInGithubModeAndNeedsAPR(t *testing.T) {
	log := synthGH(t, "")
	if err, _, _ := runSynth(t, "--root", synthProject(t, "version: 1\n"), "--pr-a", "693"); err == nil ||
		!strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode: %v", err)
	}
	if err, _, _ := runSynth(t, "--root", synthProject(t, githubModeConfig)); err == nil ||
		!strings.Contains(err.Error(), "--pr-a") {
		t.Errorf("without --pr-a: %v", err)
	}
	if calls := readLog(t, log); calls != "" {
		t.Errorf("a refused call still talked to GitHub:\n%s", calls)
	}
}

// --dry-run shows the card — both PRs, both cards, both titles — and touches nothing.
func TestSynthesizeDryRunShowsTheCardOnly(t *testing.T) {
	log := synthGH(t, "")
	err, out, _ := runSynth(t, "--root", synthProject(t, githubModeConfig),
		"--pr-a", "#693", "--pr-b", "556", "--files", "Dashboards.spec.md", "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"[synthesis] what PRs #693 and #556 deliver, in one",
		"#666", "#198", "fix(DSHBR): the B01", "feat: the index", "Dashboards.spec.md"} {
		if !strings.Contains(out, want) {
			t.Errorf("the dry run misses %q:\n%s", want, out)
		}
	}
	calls := readLog(t, log)
	for _, forbidden := range []string{"issue create", "pr close", "pr comment", "label create"} {
		if strings.Contains(calls, forbidden) {
			t.Errorf("--dry-run ran %q:\n%s", forbidden, calls)
		}
	}
}

// The real run links the five ends: the card (labelled under both origin cards), both
// PRs commented and closed, both origin cards pointed at the new one.
func TestSynthesizeLinksTheFiveEnds(t *testing.T) {
	log := synthGH(t, "")
	err, out, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig),
		"--pr-a", "693", "--pr-b", "556")
	if err != nil {
		t.Fatalf("synthesize: %v\n%s", err, errOut)
	}
	calls := readLog(t, log)
	for _, want := range []string{
		"label create anchors:under-666 --repo acme/app",
		"label create anchors:under-198 --repo acme/app",
		"--label anchors --label anchors:to-do --label anchors:under-666 --label anchors:under-198",
		"pr comment 693 --repo acme/app", "pr close 693 --repo acme/app",
		"pr comment 556 --repo acme/app", "pr close 556 --repo acme/app",
		"issue comment 666 --repo acme/app", "issue comment 198 --repo acme/app",
	} {
		if !strings.Contains(calls, want) {
			t.Errorf("missing call %q:\n%s", want, calls)
		}
	}
	if !strings.Contains(calls, "https://github.com/acme/app/issues/700") {
		t.Errorf("the comments do not point at the new card:\n%s", calls)
	}
	if !strings.Contains(out, "synthesis card: https://github.com/acme/app/issues/700") ||
		!strings.Contains(out, "point at each other") {
		t.Errorf("output:\n%s", out)
	}
	if errOut != "" {
		t.Errorf("no warning was expected:\n%s", errOut)
	}
}

// With one side only, the other end is the integration branch, and only that PR closes.
// Every failed link is WARNED — the card already exists.
func TestSynthesizeOneSideWarnsOnEachFailedLink(t *testing.T) {
	log := synthGH(t, `case "$*" in "pr comment"*|"pr close"*|"issue comment"*) echo nope >&2; exit 1 ;; esac`)
	err, _, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693")
	if err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	calls := readLog(t, log)
	if !strings.Contains(calls, "what PR #693 delivers, reconciled with what already landed") {
		t.Errorf("the one-sided title is missing:\n%s", calls)
	}
	if !strings.Contains(calls, "what is already on the integration branch") {
		t.Errorf("the PR comment does not name the other side:\n%s", calls)
	}
	if strings.Contains(calls, "pr close 556") || strings.Contains(calls, "under-198") {
		t.Errorf("a side that does not exist was touched:\n%s", calls)
	}
	for _, want := range []string{"#693 did not receive the comment", "#693 was not closed", "card #666 did not receive the link"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("no warning %q:\n%s", want, errOut)
		}
	}
}

func TestSynthesizeFailsWhenTheCardCannotBeOpened(t *testing.T) {
	log := synthGH(t, `case "$*" in "issue create"*) echo "HTTP 403" >&2; exit 1 ;; esac`)
	err, _, _ := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693", "--pr-b", "556")
	if err == nil || !strings.Contains(err.Error(), "open the synthesis card") || !strings.Contains(err.Error(), "HTTP 403") {
		t.Fatalf("want the card error with gh's message, got %v", err)
	}
	if strings.Contains(readLog(t, log), "pr close") {
		t.Error("the PRs were closed although the synthesis card does not exist")
	}
}

// With one side only, the closing line tells what happened: ONE PR was closed. It used to
// print "PRs #693 and # were closed, and the four ends point at each other" — a second PR
// that does not exist, and four ends where there are three.
func TestSynthesizeOneSideReportsOnePRClosed(t *testing.T) {
	synthGH(t, "")
	err, out, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693")
	if err != nil {
		t.Fatalf("synthesize: %v\n%s", err, errOut)
	}
	if strings.Contains(out, "and #") || strings.Contains(out, "four ends") {
		t.Errorf("the one-sided run reports a second PR:\n%s", out)
	}
	if !strings.Contains(out, "PR #693 was closed") {
		t.Errorf("the one-sided run does not say PR #693 was closed:\n%s", out)
	}
}

// The closing line names only what gh actually closed: a PR whose close failed is
// reported as still open, never as closed.
func TestSynthesizeReportsAFailedCloseAsStillOpen(t *testing.T) {
	synthGH(t, `case "$*" in "pr close 556"*) echo "GraphQL: cannot close" >&2; exit 1 ;; esac`)
	err, out, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693", "--pr-b", "556")
	if err != nil {
		t.Fatalf("synthesize: %v\n%s", err, errOut)
	}
	if strings.Contains(out, "#556 were closed") || strings.Contains(out, "PRs #693, #556") {
		t.Errorf("a PR whose close failed was reported closed:\n%s", out)
	}
	if !strings.Contains(out, "PR #693 was closed") || !strings.Contains(out, "still OPEN") || !strings.Contains(out, "#556") {
		t.Errorf("the line must say #693 closed and #556 still open:\n%s", out)
	}
}

// The synthesis card speaks the project's language: title, body and comments come from
// the message catalog, not from English written into the code.
func TestSynthesisCardFollowsTheLanguage(t *testing.T) {
	prev := i18n.Current()
	if err := i18n.Set("pt-BR"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = i18n.Set(prev) }()
	body := corpoDaSintese("693", "556", "666", "198", "fix", "feat", "src/a.ts")
	for _, want := range []string{"Dois PRs escreveram", "## Como entregar", "`Refs #693` e `Refs #556`", "- `src/a.ts`"} {
		if !strings.Contains(body, want) {
			t.Errorf("the pt-BR card must carry %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "How to deliver") {
		t.Errorf("English leaked into the pt-BR card:\n%s", body)
	}
}

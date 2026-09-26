package ops

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// identidadeGitNoTeste dá ao git um autor SÓ para este teste.
//
// Sem isto o `git commit` do `iniciaGit` falha com "Author identity unknown" em qualquer
// máquina sem `user.name`/`user.email` globais — e os três testes deste arquivo passavam
// na máquina de quem os escreveu e falhavam no CI. Descoberto na PRIMEIRA execução do
// workflow de CI: três FAILs, todos aqui.
//
// É a violação do invariante que o próprio projeto declara para teste (o `TSHRT-B06` do
// blue-eyes, e a doutrina geral): um teste não depende do ambiente. Um teste que só passa
// numa máquina configurada não prova o código — prova a máquina.
//
// Por ENV e não por `git config`: a variável vale só para os processos deste teste, então
// nada é escrito no repositório temporário nem na configuração de quem roda. `t.Setenv`
// restaura no fim, e falha se o teste for paralelo — o que é a garantia certa, porque
// mexer em env compartilhado sob paralelismo é corrida.
//
// A correção é no TESTE, não no `iniciaGit`: exigir identidade é o comportamento certo do
// produto, e a mensagem de erro dele orienta o usuário. Injetar autor no produto
// commitaria como alguém que ninguém escolheu.
func identidadeGitNoTeste(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_AUTHOR_NAME", "anchors-test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@anchors.invalid")
	t.Setenv("GIT_COMMITTER_NAME", "anchors-test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@anchors.invalid")
}

// O caminho em que o usuário ACEITA: o repo tem de sair utilizável de verdade — com
// HEAD. Um `git init` sem commit deixa `git log`/`git diff` sem contra o que comparar,
// e metade do Anchors seguiria desligada achando que foi resolvida.
func TestIniciaGitDeixaRepoComHEAD(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	identidadeGitNoTeste(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".DS_Store"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := initGit(dir, initx.GitNaoIniciado); err != nil {
		t.Fatalf("iniciaGit: %v", err)
	}

	if e := initx.DetectGit(dir, true); e != initx.GitPronto {
		t.Fatalf("depois de iniciar, o estado tem de ser GitPronto, foi %v", e)
	}
	out, err := runGit(dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("git log falhou — não há HEAD: %v (%s)", err, out)
	}
	if strings.TrimSpace(out) != initx.FirstCommitMessage {
		t.Errorf("assunto do commit = %q, queria %q", out, initx.FirstCommitMessage)
	}
	b, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore não foi semeado: %v", err)
	}
	if !strings.Contains(string(b), ".DS_Store") {
		t.Error(".gitignore semeado não cobre .DS_Store")
	}
}

// Um .gitignore que já existe é do USUÁRIO e vale mais que o nosso padrão.
func TestIniciaGitNaoSobrescreveGitignoreExistente(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	identidadeGitNoTeste(t)
	dir := t.TempDir()
	meu := "# meu\n*.log\n"
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(meu), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := initGit(dir, initx.GitNaoIniciado); err != nil {
		t.Fatalf("iniciaGit: %v", err)
	}

	b, _ := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if string(b) != meu {
		t.Errorf(".gitignore do usuário foi sobrescrito: %q", string(b))
	}
}

// Num repo que já existe sem commit, só falta o commit — `git init` de novo seria
// desnecessário, e o estado final é o mesmo: HEAD existindo.
func TestIniciaGitSoCommitaQuandoRepoJaExiste(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	identidadeGitNoTeste(t)
	dir := t.TempDir()
	if out, err := runGit(dir, "init"); err != nil {
		t.Fatalf("preparo: %s", out)
	}

	if err := initGit(dir, initx.GitSemCommit); err != nil {
		t.Fatalf("iniciaGit: %v", err)
	}
	if e := initx.DetectGit(dir, true); e != initx.GitPronto {
		t.Fatalf("estado final = %v, queria GitPronto", e)
	}
}

// --- gitStep: the first step of `anchors init` ---

// A machine without git gets the warning and no question: offering `git init` there would
// fail the moment it is accepted. And init goes on — git missing is a warning, not an error.
func TestGitStepWithoutGitWarnsAndProceeds(t *testing.T) {
	resetPromptError(t)
	t.Setenv("PATH", t.TempDir()) // no git on this PATH
	root := t.TempDir()

	var ok bool
	out := captureStdout(t, func() { ok = gitStep(root) })
	if !ok {
		t.Fatal("gitStep stopped init on a machine without git; it must only warn")
	}
	if !strings.Contains(out, initx.AvisoGit(initx.GitNaoInstalado)) {
		t.Errorf("the not-installed warning is missing:\n%s", out)
	}
	if erroDePrompt {
		t.Error("a prompt ran although there is nothing to offer without git")
	}
}

// A repository with a commit is ready: nothing is printed and nothing is asked.
func TestGitStepOnAReadyRepoIsSilent(t *testing.T) {
	resetPromptError(t)
	root := newGitRepo(t)
	var ok bool
	out := captureStdout(t, func() { ok = gitStep(root) })
	if !ok || out != "" {
		t.Errorf("gitStep on a ready repo = %v, printed %q; want true and silence", ok, out)
	}
}

// A repository WITHOUT a commit is offered the first commit. When the offer cannot be
// asked, gitStep says so (false) and does not commit on its own.
func TestGitStepOnARepoWithoutCommitDoesNotActWithoutAnAnswer(t *testing.T) {
	skipIfTerminal(t)
	resetPromptError(t)
	isolateGit(t)
	root := t.TempDir()
	if out, err := runGit(root, "init", "-q"); err != nil {
		t.Fatalf("git init: %s", out)
	}

	var ok bool
	out := captureStdout(t, func() { ok = gitStep(root) })
	if ok {
		t.Error("with the prompt failing, gitStep must tell init to abort")
	}
	if !strings.Contains(out, initx.AvisoGit(initx.GitSemCommit)) {
		t.Errorf("the no-commit warning is missing:\n%s", out)
	}
	if _, err := runGit(root, "rev-parse", "HEAD"); err == nil {
		t.Error("a commit was made although nobody accepted")
	}
}

// Without an identity the commit fails, and the error names the fix instead of pasting
// git's long message.
func TestInitGitWithoutIdentityNamesTheFix(t *testing.T) {
	withoutGitIdentity(t)
	root := t.TempDir()

	var err error
	captureStdout(t, func() { err = initGit(root, initx.GitNaoIniciado) })
	if err == nil {
		t.Fatal("the commit cannot succeed without an identity")
	}
	msg := err.Error()
	if !strings.Contains(msg, "git does not know who you are") ||
		!strings.Contains(msg, `git config --global user.email`) {
		t.Errorf("the error does not name the fix: %v", err)
	}
	if strings.Count(msg, "\n") > 4 {
		t.Errorf("only the first line of git's output belongs in the detail: %v", err)
	}
}

func TestFirstLineCutsAtTheFirstNewline(t *testing.T) {
	for in, want := range map[string]string{
		"one\ntwo\nthree": "one",
		"single":          "single",
		"":                "",
	} {
		if got := firstLine(in); got != want {
			t.Errorf("firstLine(%q) = %q, want %q", in, got, want)
		}
	}
}

// answerGitOffer runs gitStep on a directory outside git, answering the offer with
// `reply`. TERM=dumb puts the prompt in its line-based mode, which reads os.Stdin instead
// of opening a terminal — so the answer is deterministic.
func answerGitOffer(t *testing.T, reply string) (dir string, ok bool, out string) {
	t.Helper()
	resetPromptError(t)
	isolateGit(t)
	t.Setenv("TERM", "dumb")
	dir = t.TempDir()
	withStdin(t, reply, func() {
		out = captureStdout(t, func() { ok = gitStep(dir) })
	})
	return dir, ok, out
}

// Declining is legitimate and init goes on — but it names NOW what stays incomplete.
func TestGitStepDeclinedNamesWhatStaysOff(t *testing.T) {
	dir, ok, out := answerGitOffer(t, "n\n")
	if !ok {
		t.Error("declining git must not stop init")
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		t.Error("a repository was created although the offer was declined")
	}
	for _, want := range []string{"proceeding without git", "coverage --diff", "install-hooks", "`git init` later"} {
		if !strings.Contains(out, want) {
			t.Errorf("the consequence %q is not named:\n%s", want, out)
		}
	}
}

// Accepting leaves a repository with HEAD, the .gitignore seeded, and init going on.
func TestGitStepAcceptedLeavesARepoWithHead(t *testing.T) {
	dir, ok, out := answerGitOffer(t, "y\n")
	if !ok {
		t.Error("accepting git must not stop init")
	}
	if e := initx.DetectGit(dir, true); e != initx.GitPronto {
		t.Fatalf("after accepting, the state is %v, want GitPronto\n%s", e, out)
	}
	if !strings.Contains(out, "✓ first commit") || !strings.Contains(out, "✓ .gitignore seeded") {
		t.Errorf("the output does not report what was done:\n%s", out)
	}
}

// A failure while initializing does not stop init: it is reported with its cause (here,
// git without an identity).
func TestGitStepReportsAFailedInitialization(t *testing.T) {
	resetPromptError(t)
	withoutGitIdentity(t)
	t.Setenv("TERM", "dumb")
	dir := t.TempDir()
	var ok bool
	var out string
	withStdin(t, "y\n", func() {
		out = captureStdout(t, func() { ok = gitStep(dir) })
	})
	if !ok {
		t.Error("a failed initialization must not stop init")
	}
	if !strings.Contains(out, "could not initialize git") || !strings.Contains(out, "git does not know who you are") {
		t.Errorf("the failure and its cause are not reported:\n%s", out)
	}
}

// withoutGitIdentity leaves git with no author: no global config but `useConfigOnly`, and
// none of the identity variables — so a commit fails the way it does on a fresh machine.
func withoutGitIdentity(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	global := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(global, []byte("[user]\n\tuseConfigOnly = true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", global)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL", "EMAIL"} {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}
}

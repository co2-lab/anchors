package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// O caminho em que o usuário ACEITA: o repo tem de sair utilizável de verdade — com
// HEAD. Um `git init` sem commit deixa `git log`/`git diff` sem contra o que comparar,
// e metade do Anchors seguiria desligada achando que foi resolvida.
func TestIniciaGitDeixaRepoComHEAD(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".DS_Store"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := iniciaGit(dir, initx.GitNaoIniciado); err != nil {
		t.Fatalf("iniciaGit: %v", err)
	}

	if e := initx.DetectaGit(dir, true); e != initx.GitPronto {
		t.Fatalf("depois de iniciar, o estado tem de ser GitPronto, foi %v", e)
	}
	out, err := rodaGit(dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("git log falhou — não há HEAD: %v (%s)", err, out)
	}
	if strings.TrimSpace(out) != initx.MensagemPrimeiroCommit {
		t.Errorf("assunto do commit = %q, queria %q", out, initx.MensagemPrimeiroCommit)
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
	dir := t.TempDir()
	meu := "# meu\n*.log\n"
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(meu), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := iniciaGit(dir, initx.GitNaoIniciado); err != nil {
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
	dir := t.TempDir()
	if out, err := rodaGit(dir, "init"); err != nil {
		t.Fatalf("preparo: %s", out)
	}

	if err := iniciaGit(dir, initx.GitSemCommit); err != nil {
		t.Fatalf("iniciaGit: %v", err)
	}
	if e := initx.DetectaGit(dir, true); e != initx.GitPronto {
		t.Fatalf("estado final = %v, queria GitPronto", e)
	}
}

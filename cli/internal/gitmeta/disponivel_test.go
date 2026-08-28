package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func semRepoAcima(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if Verifica(dir) == Disponível {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	return dir
}

// A razão de existir do diagnóstico: `exit status 128` não diz se falta INSTALAR ou
// INICIAR, e os consertos são opostos. Cada causa tem de aparecer com o seu.
func TestExplicaDizOConsertoDeCadaCausa(t *testing.T) {
	semBin := Explica(SemBinário, "listar os arquivos staged")
	if !strings.Contains(semBin, "não está instalado") {
		t.Errorf("sem binário tem de mandar instalar: %s", semBin)
	}
	if strings.Contains(semBin, "git init") {
		t.Errorf("`git init` não roda sem o binário — mandá-lo manda ao lugar errado: %s", semBin)
	}

	semRepo := Explica(SemRepo, "listar os arquivos staged")
	if !strings.Contains(semRepo, "git init") {
		t.Errorf("sem repo tem de mandar iniciar: %s", semRepo)
	}
	if strings.Contains(semRepo, "instalado") {
		t.Errorf("o git ESTÁ instalado neste caso: %s", semRepo)
	}

	// A ação do comando entra na frase: é o que liga o sintoma à causa numa linha.
	if !strings.Contains(semRepo, "listar os arquivos staged") {
		t.Errorf("a mensagem tem de nomear a ação que ficou incompleta: %s", semRepo)
	}
}

// Com git disponível, a falha é OUTRA — inventar uma causa de git esconderia a real.
func TestExplicaCalaQuandoGitEstaDisponivel(t *testing.T) {
	if msg := Explica(Disponível, "qualquer coisa"); msg != "" {
		t.Errorf("git disponível não deve produzir explicação: %q", msg)
	}
}

func TestVerificaDistingueRepoDeSemRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := semRepoAcima(t)
	if d := Verifica(dir); d != SemRepo {
		t.Fatalf("diretório sem repo = %v, queria SemRepo", d)
	}

	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if d := Verifica(dir); d != Disponível {
		t.Fatalf("com .git = %v, queria Disponível", d)
	}
}

// Subpasta de um repo está versionada: o `.git` fica acima, e reclamar ali mandaria
// criar repo aninhado.
func TestVerificaEnxergaRepoAcima(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := semRepoAcima(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "pacotes", "app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if d := Verifica(sub); d != Disponível {
		t.Fatalf("subpasta de repo = %v, queria Disponível", d)
	}
}

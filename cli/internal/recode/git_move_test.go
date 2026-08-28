package recode

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gitmeta"
)

func prepararRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"},
	} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("preparo %v: %s", args, out)
		}
	}
	return dir
}

// FORA de repositório o rename é o comportamento CERTO: não há histórico a preservar,
// e exigir git ali seria transformar uma degradação legítima em bloqueio.
func TestMoveSemRepoUsaRenameSimples(t *testing.T) {
	dir := t.TempDir()
	if gitmeta.Verifica(dir) == gitmeta.Disponível {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gitMove(dir, "a.go", "sub/b.go"); err != nil {
		t.Fatalf("sem repositório o rename deveria funcionar: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "sub", "b.go")); err != nil {
		t.Errorf("o arquivo não chegou ao destino: %v", err)
	}
}

// Num repo, um arquivo rastreado move via `git mv` — e o índice registra a renomeação.
func TestMoveComRepoUsaGitMv(t *testing.T) {
	dir := prepararRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-m", "base"}} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("preparo %v: %s", args, out)
		}
	}

	if err := gitMove(dir, "a.go", "sub/b.go"); err != nil {
		t.Fatalf("git mv deveria funcionar: %v", err)
	}

	c := exec.Command("git", "status", "--porcelain")
	c.Dir = dir
	out, _ := c.Output()
	// `R` (staged rename) é o que prova que o histórico foi preservado; um rename por
	// fora apareceria como `D` + `??`.
	if !strings.Contains(string(out), "R ") {
		t.Errorf("a renomeação não foi registrada no índice: %s", out)
	}
}

// Arquivo ainda NÃO RASTREADO não é recusa a respeitar: o git nunca soube dele, então
// não há histórico a preservar nem índice a manter coerente. O rename direto faz o que
// o usuário pediu.
//
// Este é o caso de um projeto recém-inicializado — nada commitado ainda —, que é
// justamente onde o recode tem mais chance de rodar.
func TestMoveDeArquivoNaoRastreadoAcontece(t *testing.T) {
	dir := prepararRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "solto.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gitMove(dir, "solto.go", "sub/solto.go"); err != nil {
		t.Fatalf("arquivo não rastreado deveria ser movido direto: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "sub", "solto.go")); err != nil {
		t.Errorf("o arquivo não chegou ao destino: %v", err)
	}
}

// O bug que o `-k` escondia: `git mv -k` PULA em silêncio o que não consegue mover e
// sai com status 0. O rename não acontecia, `err` era nil, e o Apply contava como
// feito — um projeto sem nada commitado veria "✓ N arquivos reescritos" com os N
// parados no lugar. Este teste fixa que "sucesso" significa "o arquivo mudou de sítio".
func TestMoveQueRelataSucessoMoveuDeVerdade(t *testing.T) {
	dir := prepararRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "solto.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gitMove(dir, "solto.go", "movido.go"); err != nil {
		t.Fatalf("gitMove: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "solto.go")); err == nil {
		t.Error("gitMove disse que moveu, mas a origem continua lá — é o silêncio do `-k`")
	}
	if _, err := os.Stat(filepath.Join(dir, "movido.go")); err != nil {
		t.Errorf("gitMove disse que moveu, mas o destino não existe: %v", err)
	}
}

// Uma recusa REAL do git (aqui: destino já existente e rastreado) não é contornada por
// os.Rename. O git recusa por um motivo, e mover à revelia deixa o índice divergindo
// do disco — em silêncio, num comando que mexe em massa.
func TestMoveNaoContornaRecusaDoGit(t *testing.T) {
	dir := prepararRepo(t)
	for _, f := range []string{"a.go", "b.go"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("package a\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-m", "base"}} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("preparo %v: %s", args, out)
		}
	}

	// Destino já existe e é rastreado: `git mv` recusa sem `-f`.
	err := gitMove(dir, "a.go", "b.go")

	if err == nil {
		t.Fatal("git recusou o move — contorná-lo com os.Rename sobrescreveria um arquivo rastreado")
	}
	if !strings.Contains(err.Error(), "recusou") {
		t.Errorf("a mensagem tem de dizer que foi o git que recusou: %v", err)
	}
	// O destino tem de continuar sendo o ORIGINAL: nada foi sobrescrito por fora.
	if _, statErr := os.Stat(filepath.Join(dir, "a.go")); statErr != nil {
		t.Error("a origem sumiu numa operação que falhou")
	}
}

package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func foraDeRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if Verifica(dir) == Disponível {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	return dir
}

// O silêncio que tranquiliza: sem repositório, `DirtyCount` devolvia 0 — e 0 vira
// "# árvore: limpa" no cabeçalho do relatório. Um relatório que AFIRMA limpeza sem ter
// conseguido olhar carimba uma foto que nunca existiu.
func TestDirtyCountNaoAfirmaLimpezaQueNaoVerificou(t *testing.T) {
	dir := foraDeRepo(t)

	if n := DirtyCount(dir); n >= 0 {
		t.Fatalf("sem repositório, DirtyCount tem de dizer que não sabe (<0), devolveu %d", n)
	}
}

// A contrapartida: num repo de verdade a contagem continua sendo a contagem.
func TestDirtyCountContaNoRepoDeVerdade(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("preparo: %s", out)
	}

	if n := DirtyCount(dir); n != 0 {
		t.Fatalf("repo recém-criado e vazio tem 0 arquivos sujos, veio %d", n)
	}
	if err := os.WriteFile(filepath.Join(dir, "novo.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if n := DirtyCount(dir); n != 1 {
		t.Fatalf("com 1 arquivo novo, esperava 1, veio %d", n)
	}
}

// `false` de HasUncommittedChanges tinha DOIS significados: "está limpo" e "não deu
// para perguntar". Quem decide a partir disso precisa dos dois separados.
func TestUncommittedChangesSeparaLimpoDeNaoSabido(t *testing.T) {
	dir := foraDeRepo(t)

	mudou, sabido := UncommittedChanges(dir, "qualquer.go")

	if sabido {
		t.Error("sem repositório a pergunta não pôde ser feita — dizer que foi é afirmar sem ter lido")
	}
	if mudou {
		t.Error("sem resposta, não há mudança a afirmar")
	}
}

func TestUncommittedChangesSabeQuandoHaRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("preparo: %s", out)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	mudou, sabido := UncommittedChanges(dir, "a.go")

	if !sabido {
		t.Fatal("com repositório a pergunta pôde ser feita")
	}
	if !mudou {
		t.Error("arquivo novo e não commitado É uma alteração pendente")
	}
}

// A compatibilidade: o wrapper booleano continua valendo para quem só quer o sim/não
// e para quem o `false` não engana.
func TestHasUncommittedChangesSegueValendo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("preparo: %s", out)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !HasUncommittedChanges(dir, "a.go") {
		t.Error("arquivo novo no repo é alteração pendente")
	}
}

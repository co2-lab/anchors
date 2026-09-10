package queue

import (
	"os"
	"path/filepath"
	"testing"
)

// TestClaimNaoReivindicaOQueJaTemDono guarda o pior caso da troca do rename por O_EXCL:
// se o processo morrer entre criar o claimed__X e apagar o pending__X, os dois arquivos
// coexistem. O que NÃO pode acontecer é o pending resíduo ser servido como se a task
// estivesse livre — seria o double-claim de volta, por outra porta.
func TestClaimNaoReivindicaOQueJaTemDono(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "fa.spec.md")
	if _, err := Enqueue(root, task("residuo-a", "fa.spec.md", "spec", "implement")); err != nil {
		t.Fatal(err)
	}

	// Simula a morte no meio: o claimed já existe, o pending ficou para trás.
	d := dirFor(root)
	claimed := filepath.Join(d, fileName(Claimed, "residuo-a"))
	if err := os.WriteFile(claimed, []byte("id: residuo-a\nstate: claimed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	pendente := filepath.Join(d, fileName(Pending, "residuo-a"))
	if _, err := os.Stat(pendente); err != nil {
		t.Fatalf("o cenário exige o pending ainda no disco: %v", err)
	}

	got, err := Claim(root, "w2", "2026-08-22T00:00:00-03:00")
	if err != nil {
		t.Fatalf("Claim devolveu erro: %v", err)
	}
	if got != nil {
		t.Fatalf("reivindicou %q, que já tem dono", got.ID)
	}
	// E o resíduo se limpa sozinho: quem topa com ele apaga, em vez de deixar a fila
	// mostrando para sempre um pendente que ninguém pode pegar.
	if _, err := os.Stat(pendente); !os.IsNotExist(err) {
		t.Error("o pending resíduo devia ter sido removido ao encontrar o claimed")
	}
}

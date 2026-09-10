package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func gitRepo(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"}, {"config", "user.email", "t@t.co"}, {"config", "user.name", "t"},
	} {
		if err := exec.Command("git", append([]string{"-C", d}, args...)...).Run(); err != nil {
			t.Skip("git indisponível")
		}
	}
	return d
}

func TestLastCommitDate(t *testing.T) {
	d := gitRepo(t)
	os.WriteFile(filepath.Join(d, "a.txt"), []byte("x"), 0o644)
	exec.Command("git", "-C", d, "add", "-A").Run()
	exec.Command("git", "-C", d, "commit", "-q", "-m", "c").Run()
	date, ok := LastCommitDate(d, "a.txt")
	if !ok || len(date) != 10 {
		t.Fatalf("esperava data AAAAX-MM-DD, veio %q ok=%v", date, ok)
	}
	// arquivo inexistente → não ok
	if _, ok := LastCommitDate(d, "nao-existe.txt"); ok {
		t.Error("arquivo sem commit deveria devolver ok=false")
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	d := gitRepo(t)
	p := filepath.Join(d, "a.txt")
	os.WriteFile(p, []byte("x"), 0o644)
	exec.Command("git", "-C", d, "add", "-A").Run()
	exec.Command("git", "-C", d, "commit", "-q", "-m", "c").Run()
	if HasUncommittedChanges(d, "a.txt") {
		t.Error("recém-commitado não deveria ter mudanças")
	}
	os.WriteFile(p, []byte("y"), 0o644)
	if !HasUncommittedChanges(d, "a.txt") {
		t.Error("após editar, deveria ter mudanças não-commitadas")
	}
}

func TestToday(t *testing.T) {
	if len(Today()) != 10 {
		t.Errorf("Today deveria ser AAAAX-MM-DD, veio %q", Today())
	}
}

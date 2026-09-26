package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
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

// commitOn commits everything staged-or-not with the given author date.
func commitOn(t *testing.T, d, date, msg string) {
	t.Helper()
	exec.Command("git", "-C", d, "add", "-A").Run()
	c := exec.Command("git", "-C", d, "commit", "-q", "-m", msg)
	c.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+date+"T12:00:00", "GIT_COMMITTER_DATE="+date+"T12:00:00")
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("commit: %v %s", err, out)
	}
}

func TestAllCommitDates(t *testing.T) {
	d := gitRepo(t)
	os.MkdirAll(filepath.Join(d, "pkg"), 0o755)
	os.WriteFile(filepath.Join(d, "a.txt"), []byte("1"), 0o644)
	os.WriteFile(filepath.Join(d, "pkg", "b.txt"), []byte("1"), 0o644)
	commitOn(t, d, "2026-01-10", "first")
	os.WriteFile(filepath.Join(d, "a.txt"), []byte("2"), 0o644)
	commitOn(t, d, "2026-03-05", "second")

	got := AllCommitDates(d)
	// Each file carries the date of the MOST RECENT commit that touched it.
	want := map[string]string{"a.txt": "2026-03-05", "pkg/b.txt": "2026-01-10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AllCommitDates = %v, want %v", got, want)
	}

	if got := AllCommitDates(t.TempDir()); len(got) != 0 {
		t.Fatalf("outside a repository the map must be empty, got %v", got)
	}
}

func TestHead(t *testing.T) {
	d := gitRepo(t)
	if _, _, ok := Head(d); ok {
		t.Fatal("a repository without commits has no HEAD to report")
	}
	if _, _, ok := Head(t.TempDir()); ok {
		t.Fatal("outside a repository Head must report ok=false")
	}

	os.WriteFile(filepath.Join(d, "a.txt"), []byte("x"), 0o644)
	commitOn(t, d, "2026-02-01", "feat: the subject line")
	short, subject, ok := Head(d)
	want, _ := exec.Command("git", "-C", d, "rev-parse", "--short", "HEAD").Output()
	if !ok || short != strings.TrimSpace(string(want)) || subject != "feat: the subject line" {
		t.Fatalf("Head = %q, %q, %v; want %q, the subject, true", short, subject, ok, want)
	}
}

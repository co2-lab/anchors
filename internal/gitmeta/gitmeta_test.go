package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func gitRepo(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"}, {"config", "user.email", "t@t.co"}, {"config", "user.name", "t"},
	} {
		if err := exec.Command("git", append([]string{"-C", d}, args...)...).Run(); err != nil {
			t.Skip("git unavailable")
		}
	}
	return d
}

func TestLastCommitDate(t *testing.T) {
	t.Run("GTMTG-B02: The last commit date is the day of the most recent commit of the file", func(t *testing.T) {})
	d := gitRepo(t)
	os.WriteFile(filepath.Join(d, "a.txt"), []byte("1"), 0o644)
	commitOn(t, d, "2026-01-10", "first")
	os.WriteFile(filepath.Join(d, "a.txt"), []byte("2"), 0o644)
	commitOn(t, d, "2026-03-05", "second")
	date, ok := LastCommitDate(d, "a.txt")
	if !ok || date != "2026-03-05" {
		t.Fatalf("want the day of the most recent commit 2026-03-05, got %q ok=%v", date, ok)
	}
	// a file never committed → not ok
	if _, ok := LastCommitDate(d, "nao-existe.txt"); ok {
		t.Error("a file without a commit should give ok=false")
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	t.Run("GTMTG-B05: The shortcut says yes after an edit and no after a commit", func(t *testing.T) {})
	d := gitRepo(t)
	p := filepath.Join(d, "a.txt")
	os.WriteFile(p, []byte("x"), 0o644)
	exec.Command("git", "-C", d, "add", "-A").Run()
	exec.Command("git", "-C", d, "commit", "-q", "-m", "c").Run()
	if HasUncommittedChanges(d, "a.txt") {
		t.Error("a just-committed file should have no changes")
	}
	os.WriteFile(p, []byte("y"), 0o644)
	if !HasUncommittedChanges(d, "a.txt") {
		t.Error("after an edit it should have uncommitted changes")
	}
}

func TestToday(t *testing.T) {
	t.Run("GTMTG-B01: Today is the system date", func(t *testing.T) {})
	before := time.Now().Format("2006-01-02")
	got := Today()
	after := time.Now().Format("2006-01-02")
	if got != before && got != after {
		t.Errorf("Today = %q, want the system date %q", got, before)
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
	t.Run("GTMTG-B03: The bulk reader gives each file its most recent commit day", func(t *testing.T) {})
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
	t.Run("GTMTG-B06: HEAD is the short hash and subject, and unknown without a commit", func(t *testing.T) {})
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

// The reassuring silence: without a repository, `DirtyCount` used to return 0 — and 0 becomes
// "tree: clean" in the report header. A report that CLAIMS a clean tree without having been
// able to look stamps a snapshot that never existed.
func TestDirtyCountDoesNotClaimACleanTreeItDidNotCheck(t *testing.T) {
	t.Run("GTMTG-X01: A tree that could not be counted is not reported clean", func(t *testing.T) {})
	dir := outsideAnyRepo(t)

	if n := DirtyCount(dir); n >= 0 {
		t.Fatalf("without a repository DirtyCount must say it does not know (<0), got %d", n)
	}
}

// The counterpart: in a real repo the count is still the count.
func TestDirtyCountCountsInARealRepo(t *testing.T) {
	t.Run("GTMTG-B07: The dirty count counts the modified files of a real repository", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("setup: %s", out)
	}

	if n := DirtyCount(dir); n != 0 {
		t.Fatalf("a new empty repo has 0 dirty files, got %d", n)
	}
	if err := os.WriteFile(filepath.Join(dir, "novo.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if n := DirtyCount(dir); n != 1 {
		t.Fatalf("with 1 new file, want 1, got %d", n)
	}
}

// A `false` from HasUncommittedChanges had TWO meanings: "it is clean" and "could not ask".
// Whoever decides from it needs the two apart.
func TestUncommittedChangesTellsCleanFromUnknown(t *testing.T) {
	t.Run("GTMTG-X02: Outside a repository the pending-change question is not known", func(t *testing.T) {})
	dir := outsideAnyRepo(t)

	changed, known := UncommittedChanges(dir, "qualquer.go")

	if known {
		t.Error("without a repository the question could not be asked — saying it was is claiming without reading")
	}
	if changed {
		t.Error("without an answer there is no change to claim")
	}
}

func TestUncommittedChangesKnowsWhenThereIsARepo(t *testing.T) {
	t.Run("GTMTG-B04: A new file in a repository has pending changes, and the question was asked", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("setup: %s", out)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	changed, known := UncommittedChanges(dir, "a.go")

	if !known {
		t.Fatal("with a repository the question could be asked")
	}
	if !changed {
		t.Error("a new, uncommitted file IS a pending change")
	}
}

// Compatibility: the boolean wrapper still holds for whoever only wants the yes/no.
func TestHasUncommittedChangesStillHolds(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("setup: %s", out)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !HasUncommittedChanges(dir, "a.go") {
		t.Error("a new file in the repo is a pending change")
	}
}

func TestAtHead(t *testing.T) {
	t.Run("GTMTG-B08: A file's content at the last commit", func(t *testing.T) {})
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", dir}, args...)...)
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := c.CombinedOutput(); err != nil {
			t.Skipf("git: %v %s", err, out)
		}
	}
	run("init", "-q")
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("v1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-q", "-m", "one")
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("v2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got, ok := AtHead(dir, "a.md"); !ok || got != "v1\n" {
		t.Errorf("the committed version, got %q %v", got, ok)
	}
	if _, ok := AtHead(dir, "new.md"); ok {
		t.Error("a file HEAD does not have has no committed version")
	}
	if _, ok := AtHead(t.TempDir(), "a.md"); ok {
		t.Error("no repository, no committed version")
	}
}

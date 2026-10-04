// @anchors
//   code: CHTSE
//   ref: CHLGC

package checklog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMirrorCopiesOutputToTheFile(t *testing.T) {
	t.Run("CHLGC-B02: The header comes first and the output is copied after it", func(t *testing.T) {})
	dir := t.TempDir()
	e := Open(dir, true, "# header\n\n")
	if e == nil {
		t.Fatal("Open returned nil")
	}
	fmt.Println("report line")
	e.Close()

	b, err := os.ReadFile(filepath.Join(dir, Dir, "check-all.txt"))
	if err != nil {
		t.Fatalf("reading the mirror: %v", err)
	}
	got := string(b)
	if !strings.HasPrefix(got, "# header") {
		t.Errorf("the header is not first:\n%s", got)
	}
	if !strings.Contains(got, "report line") {
		t.Errorf("output not mirrored:\n%s", got)
	}
}

// `--all` takes minutes and `--changed` runs on every commit. If both wrote to the same
// file, the pre-commit would erase the full snapshot — the expensive one to reproduce.
func TestScopesDoNotOverwriteEachOther(t *testing.T) {
	t.Run("CHLGC-B01: Each scope is mirrored to its own file", func(t *testing.T) {})
	t.Run("CHLGC-I01: The changed-files mirror leaves the full snapshot intact", func(t *testing.T) {})
	dir := t.TempDir()

	e := Open(dir, true, "# all\n")
	fmt.Println("full snapshot")
	e.Close()

	e = Open(dir, false, "# changed\n")
	fmt.Println("incremental")
	e.Close()

	all, err := os.ReadFile(filepath.Join(dir, Dir, "check-all.txt"))
	if err != nil {
		t.Fatalf("reading check-all: %v", err)
	}
	if !strings.Contains(string(all), "full snapshot") {
		t.Errorf("--changed overwrote the --all snapshot:\n%s", all)
	}
	chg, err := os.ReadFile(filepath.Join(dir, Dir, "check-changed.txt"))
	if err != nil {
		t.Fatalf("reading check-changed: %v", err)
	}
	if strings.Contains(string(chg), "full snapshot") || !strings.Contains(string(chg), "incremental") {
		t.Errorf("scopes mixed in the same file:\n%s", chg)
	}
}

// Output larger than the pipe buffer (64KB on Linux/macOS): with nobody draining the other
// end, the command would hang writing its own report.
func TestLongOutputDoesNotHang(t *testing.T) {
	t.Run("CHLGC-B03: A long output is copied whole without hanging", func(t *testing.T) {})
	dir := t.TempDir()
	e := Open(dir, true, "")
	line := strings.Repeat("x", 200)
	for i := 0; i < 2000; i++ { // ~400KB
		fmt.Println(line)
	}
	done := make(chan struct{})
	go func() { e.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Close hung — the pipe was not being drained")
	}

	b, err := os.ReadFile(filepath.Join(dir, Dir, "check-all.txt"))
	if err != nil {
		t.Fatalf("reading the mirror: %v", err)
	}
	if n := strings.Count(string(b), line); n != 2000 {
		t.Errorf("mirror truncated: %d lines of 2000", n)
	}
}

// The mirror is a convenience. If the folder cannot be created, the check must keep writing
// to the screen — trading the scan for the report would trade the essential for the accessory.
func TestFailureToOpenDoesNotStopTheCheck(t *testing.T) {
	t.Run("CHLGC-B04: A mirror that cannot be opened does not stop the check", func(t *testing.T) {})
	dir := t.TempDir()
	// A FILE where `.anchors/` should be: MkdirAll fails.
	if err := os.WriteFile(filepath.Join(dir, Dir), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	e := Open(dir, true, "# nothing\n")
	if e != nil {
		t.Error("Open returned a mirror where it could not create the folder")
	}
	e.Close() // safe with nil
	if c := e.Path(); c != "" {
		t.Errorf("Path() = %q, want empty", c)
	}
	if os.Stdout == nil {
		t.Error("stdout became invalid after the failure")
	}
}

func TestHeaderRecordsTheContext(t *testing.T) {
	t.Run("CHLGC-B05: The header records the command, the moment and the HEAD", func(t *testing.T) {})
	when := time.Date(2026, 8, 18, 14, 32, 0, 0, time.UTC)
	h := Header("anchors check --all", "abc1234", "fix: something", 3, when)

	for _, want := range []string{
		"anchors check --all",
		"2026-08-18 14:32:00",
		"abc1234 fix: something",
		"3 modified file(s)",
	} {
		if !strings.Contains(h, want) {
			t.Errorf("header without %q:\n%s", want, h)
		}
	}

	// Without git (a new repo, or git missing) the header must not invent a HEAD.
	if noGit := Header("c", "", "", 0, when); strings.Contains(noGit, "HEAD:") {
		t.Errorf("HEAD invented without git:\n%s", noGit)
	}
}

// The header stamps the tree state for whoever REREADS the report later. Without a
// repository the dirty files cannot be counted — and writing "clean" there asserts a state
// nobody checked, to a reader who has no way to suspect it.
func TestHeaderDoesNotClaimACleanTreeWithoutKnowing(t *testing.T) {
	t.Run("CHLGC-X01: A tree that could not be counted is never called clean", func(t *testing.T) {})
	when := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)

	h := Header("anchors check", "abc123", "subject", -1, when)

	if strings.Contains(h, "tree: clean") {
		t.Errorf("claimed a clean tree without being able to count: %s", h)
	}
	if !strings.Contains(h, "unknown") {
		t.Errorf("the header must say it does not know: %s", h)
	}
}

// The counterpart: a real count is still reported.
func TestHeaderReportsTheRealCount(t *testing.T) {
	t.Run("CHLGC-B06: The tree line reports the real count", func(t *testing.T) {})
	when := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)

	if h := Header("c", "h", "a", 0, when); !strings.Contains(h, "tree: clean") {
		t.Errorf("0 dirty files IS a clean tree: %s", h)
	}
	if h := Header("c", "h", "a", 1, when); !strings.Contains(h, "tree: 1 modified file (") {
		t.Errorf("1 dirty file: %s", h)
	}
	if h := Header("c", "h", "a", 5, when); !strings.Contains(h, "5 modified file(s)") {
		t.Errorf("5 dirty files: %s", h)
	}
}

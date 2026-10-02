// @anchors
//   ref: DCLDF

package testsig

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseUnifiedDiff(t *testing.T) {
	t.Run("DCLDF-B01: Added lines are recorded under the file of the new-file header", func(t *testing.T) {})
	t.Run("DCLDF-I01: Removals do not shift the new-side numbering", func(t *testing.T) {})
	diff := `diff --git a/src/A.tsx b/src/A.tsx
--- a/src/A.tsx
+++ b/src/A.tsx
@@ -10,0 +11,2 @@
+const x = 1
+const y = 2
@@ -20,1 +22,1 @@
-old line
+new line`
	changed := parseUnifiedDiff(diff)
	lines := changed["src/A.tsx"]
	if lines == nil {
		t.Fatal("src/A.tsx should have changed lines")
	}
	for _, want := range []int{11, 12, 22} {
		if !lines[want] {
			t.Errorf("line %d should be marked as changed; got %v", want, lines)
		}
	}
	if lines[21] {
		t.Error("line 21 was not touched")
	}
}

func TestParseDiffFileDeletedFile(t *testing.T) {
	t.Run("DCLDF-B04: A deleted file records nothing", func(t *testing.T) {})
	// a deleted file (+++ /dev/null) must produce no entry
	diff := "--- a/gone.ts\n+++ /dev/null\n@@ -1,2 +0,0 @@\n-a\n-b\n"
	p := filepath.Join(t.TempDir(), "d.diff")
	os.WriteFile(p, []byte(diff), 0o644)
	changed, _ := ParseDiffFile(p)
	if len(changed) != 0 {
		t.Errorf("a deleted file should have no changed lines, got %v", changed)
	}
}

func TestGitDiff(t *testing.T) {
	t.Run("DCLDF-B06: Git compares the working copy with the current commit or with a reference", func(t *testing.T) {})
	t.Run("DCLDF-E01: A directory that is not a repository fails the git diff", func(t *testing.T) {})
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Skipf("git unavailable: %v %s", err, out)
		}
	}
	git("init", "-q")
	git("config", "user.email", "t@t.co")
	git("config", "user.name", "t")
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	os.WriteFile(filepath.Join(root, "pkg", "a.go"), []byte("l1\nl2\nl3\n"), 0o644)
	os.WriteFile(filepath.Join(root, "gone.go"), []byte("x\n"), 0o644)
	git("add", "-A")
	git("commit", "-q", "-m", "base")

	// Working tree vs HEAD: line 2 changed, line 4 added; gone.go only loses lines.
	os.WriteFile(filepath.Join(root, "pkg", "a.go"), []byte("l1\nL2\nl3\nl4\n"), 0o644)
	os.Remove(filepath.Join(root, "gone.go"))
	got, err := GitDiff(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if want := (ChangedLines{"pkg/a.go": {2: true, 4: true}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("GitDiff(worktree) = %v, want %v", got, want)
	}

	// Against a ref: committed changes since HEAD~1.
	git("add", "-A")
	git("commit", "-q", "-m", "change")
	got, err = GitDiff(root, "HEAD~1")
	if err != nil {
		t.Fatal(err)
	}
	if want := (ChangedLines{"pkg/a.go": {2: true, 4: true}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("GitDiff(HEAD~1) = %v, want %v", got, want)
	}
	if got, err := GitDiff(root, ""); err != nil || len(got) != 0 {
		t.Fatalf("a clean tree has no changed lines, got %v, %v", got, err)
	}

	if _, err := GitDiff(t.TempDir(), ""); err == nil {
		t.Fatal("GitDiff outside a repository must fail")
	}
}

func TestParseDiffFileMissing(t *testing.T) {
	t.Run("DCLDF-E02: A missing diff file surfaces the read error", func(t *testing.T) {})
	if _, err := ParseDiffFile(filepath.Join(t.TempDir(), "none.diff")); !os.IsNotExist(err) {
		t.Fatalf("a missing diff file must surface the read error, got %v", err)
	}
}

func TestHunkNewStart(t *testing.T) {
	t.Run("DCLDF-B02: A hunk header sets the starting line of the new side", func(t *testing.T) {})
	for hunk, want := range map[string]int{
		"@@ -1,2 +10,3 @@ func x()": 10,
		"@@ -5 +7 @@":               7,
		"@@ -5 +12":                 12,
		"@@ malformed @@":           0,
	} {
		if got := hunkNewStart(hunk); got != want {
			t.Errorf("hunkNewStart(%q) = %d, want %d", hunk, got, want)
		}
	}
}

// The path is cleaned of prefixes and timestamps, context lines only advance the cursor,
// and a diff file needs no repository.
func TestDiffPathsContextAndFile(t *testing.T) {
	diff := "--- a/x.go\t2020-01-01\n+++ b/x.go\t2021-01-01 10:00\n@@ -1,2 +1,4 @@\n ctx\n+add\n ctx2\n+add2\n"
	got := parseUnifiedDiff(diff)
	t.Run("DCLDF-B03: Path prefixes and a tab-separated timestamp are stripped", func(t *testing.T) {
		if _, ok := got["x.go"]; !ok || len(got) != 1 {
			t.Errorf("want the lines under x.go, got %v", got)
		}
	})
	t.Run("DCLDF-B05: Context lines advance the numbering without being recorded", func(t *testing.T) {
		if want := map[int]bool{2: true, 4: true}; !reflect.DeepEqual(got["x.go"], want) {
			t.Errorf("x.go = %v, want %v", got["x.go"], want)
		}
	})
	t.Run("DCLDF-X01: A diff file is read without git", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "c.diff")
		os.WriteFile(p, []byte("--- /dev/null\n+++ b/x.go\n@@ -0,0 +1,2 @@\n+a\n+b\n"), 0o644)
		changed, err := ParseDiffFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if want := (ChangedLines{"x.go": {1: true, 2: true}}); !reflect.DeepEqual(changed, want) {
			t.Errorf("got %v, want %v", changed, want)
		}
	})
}

// An added line whose text starts with `++ ` reads `+++ …` in the diff. It was taken as a
// new-file header: the line was lost and every later line was booked under a bogus file.
func TestDiffAddedLineLooksLikeHeader(t *testing.T) {
	t.Run("DCLDF-B07: An added line that starts with two plus signs is a line, not a header", func(t *testing.T) {
		// y.go's hunk mixes a context line, a removal and an addition, and z.go's header has
		// no closing @@: each hunk ends exactly where its counts say, so the headers after
		// it are read as headers again, and z.go's "+++ z" is still a line.
		diff := "--- a/x.go\n+++ b/x.go\n@@ -1,0 +1,3 @@\n+a\n+++ counter\n+b\n" +
			"--- a/y.go\n+++ b/y.go\n@@ -1,2 +1,2 @@\n ctx\n--- old\n+new\n" +
			"--- a/z.go\n+++ b/z.go\n@@ -0,0 +1\n+++ z\n"
		want := ChangedLines{"x.go": {1: true, 2: true, 3: true}, "y.go": {2: true}, "z.go": {1: true}}
		if got := parseUnifiedDiff(diff); !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

// bufio.Scanner refuses a line longer than its buffer; a minified bundle or a lockfile can
// put hundreds of kilobytes on one line, and every line after it would be lost.
func TestDiffVeryLongLine(t *testing.T) {
	t.Run("DCLDF-B08: A very long line in the diff is read like any other", func(t *testing.T) {
		long := "+" + strings.Repeat("x", 200*1024)
		diff := "--- a/min.js\n+++ b/min.js\n@@ -0,0 +1,2 @@\n" + long + "\n+after\n"
		want := ChangedLines{"min.js": {1: true, 2: true}}
		if got := parseUnifiedDiff(diff); !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

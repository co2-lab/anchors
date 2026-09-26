package testsig

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseUnifiedDiff(t *testing.T) {
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
		t.Fatal("deveria ter linhas mudadas em src/A.tsx")
	}
	for _, want := range []int{11, 12, 22} {
		if !lines[want] {
			t.Errorf("linha %d deveria estar marcada como mudada; veio %v", want, lines)
		}
	}
	if lines[21] {
		t.Error("linha 21 não foi tocada")
	}
}

func TestUncoveredIn(t *testing.T) {
	fc := FileCoverage{Lines: map[int]bool{10: true, 11: false, 12: true}}
	changed := map[int]bool{10: true, 11: true, 13: true} // 13 não é instrumentada
	un := fc.UncoveredIn(changed)
	if len(un) != 1 || un[0] != 11 {
		t.Fatalf("esperava [11] descoberta, veio %v", un)
	}
	if fc.InstrumentedIn(changed) != 2 { // 10 e 11 (13 não é instrumentada)
		t.Errorf("esperava 2 instrumentadas no diff, veio %d", fc.InstrumentedIn(changed))
	}
}

func TestParseDiffFileDeletedFile(t *testing.T) {
	// arquivo deletado (+++ /dev/null) não deve gerar entrada
	diff := "--- a/gone.ts\n+++ /dev/null\n@@ -1,2 +0,0 @@\n-a\n-b\n"
	p := filepath.Join(t.TempDir(), "d.diff")
	os.WriteFile(p, []byte(diff), 0o644)
	changed, _ := ParseDiffFile(p)
	if len(changed) != 0 {
		t.Errorf("arquivo deletado não deveria ter linhas mudadas, veio %v", changed)
	}
}

func TestGitDiff(t *testing.T) {
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
	if _, err := ParseDiffFile(filepath.Join(t.TempDir(), "none.diff")); !os.IsNotExist(err) {
		t.Fatalf("a missing diff file must surface the read error, got %v", err)
	}
}

func TestHunkNewStart(t *testing.T) {
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

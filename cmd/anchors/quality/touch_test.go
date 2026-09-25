package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

const touchHeader = "// @anchors\n//   code: CODEX\n//   updated_at: 2026-09-01\n"

// What `anchors touch` does with one file: it does not lie about what changed.
func TestDecideTouch(t *testing.T) {
	base := touchHeader + "export const a = 1\n"
	for _, c := range []struct {
		name          string
		current, base string
		hasBase       bool
		bump          bool
		skip          touchSkip
	}{
		{"real change", touchHeader + "export const a = 2\n", base, true, true, ""},
		{"no header", "export const a = 2\n", "export const a = 1\n", true, false, skipNoHeader},
		{"no change", base, base, true, false, skipUnchanged},
		{"only the date changed", strings.Replace(base, "2026-09-01", "2026-09-20", 1), base, true, false, skipDateOnly},
		{"already dated", strings.Replace(touchHeader, "2026-09-01", "2026-09-25", 1) + "export const a = 2\n", base, true, false, skipAlready},
		{"new file", touchHeader + "export const b = 1\n", "", false, true, ""},
		{"updated_at outside a header", "// updated_at: 2026-09-01 in a comment\n", "", false, false, skipNoHeader},
	} {
		d := decideTouch("a.ts", c.current, c.base, c.hasBase, "2026-09-25")
		if d.Bump != c.bump || d.Skip != c.skip {
			t.Errorf("%s: bump=%v skip=%q, want bump=%v skip=%q", c.name, d.Bump, d.Skip, c.bump, c.skip)
		}
		if d.Bump && !strings.Contains(d.Content, "updated_at: 2026-09-25") {
			t.Errorf("%s: the new content does not carry the date:\n%s", c.name, d.Content)
		}
	}
}

func touchRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", root, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	for _, f := range []string{"a.ts", "b.ts", "gen/c.ts", "d.ts"} {
		touchWrite(t, root, f, touchHeader+"export const x = 1\n")
	}
	git("add", ".")
	git("commit", "-qm", "base")
	return root
}

func touchWrite(t *testing.T, root, f, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, f)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, f), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runTouch(t *testing.T, args ...string) string {
	t.Helper()
	cmd := newTouchCmd()
	cmd.SetArgs(args)
	return captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
}

// --changed (the default): real changes and new files are bumped; a date-only change and
// an excluded file are not, and each skip says why.
func TestTouch_changed(t *testing.T) {
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")                                                 // real change
	touchWrite(t, root, "b.ts", strings.Replace(touchHeader, "2026-09-01", "2026-09-10", 1)+"export const x = 1\n") // date only
	touchWrite(t, root, "gen/c.ts", touchHeader+"export const x = 2\n")                                             // excluded
	touchWrite(t, root, "new.ts", touchHeader+"export const n = 1\n")                                               // untracked

	out := runTouch(t, "--root", root, "--date", "2026-09-25", "--exclude", "gen/**")
	for _, want := range []string{
		"bumped a.ts  (2026-09-01 → 2026-09-25)",
		"bumped new.ts",
		"skipped b.ts — " + string(skipDateOnly),
		"skipped gen/c.ts — " + string(skipExcluded),
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the output should say %q:\n%s", want, out)
		}
	}
	for f, want := range map[string]string{"a.ts": "2026-09-25", "new.ts": "2026-09-25", "b.ts": "2026-09-10", "gen/c.ts": "2026-09-01", "d.ts": "2026-09-01"} {
		b, _ := os.ReadFile(filepath.Join(root, f))
		if !strings.Contains(string(b), "updated_at: "+want) {
			t.Errorf("%s should carry updated_at %s:\n%s", f, want, b)
		}
	}
}

// --staged: a clean staged file is bumped and re-staged; one with changes outside the
// index is skipped, because re-staging it would commit changes nobody chose.
func TestTouch_staged(t *testing.T) {
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, root, "d.ts", touchHeader+"export const x = 2\n")
	if out, err := exec.Command("git", "-C", root, "add", "a.ts", "d.ts").CombinedOutput(); err != nil {
		t.Fatalf("git add: %v %s", err, out)
	}
	touchWrite(t, root, "d.ts", touchHeader+"export const x = 3\n") // unstaged on top

	out := runTouch(t, "--root", root, "--staged", "--date", "2026-09-25")
	if !strings.Contains(out, "bumped a.ts") || !strings.Contains(out, "skipped d.ts — "+string(skipUnstaged)) {
		t.Fatalf("staged: a.ts bumped, d.ts skipped for its unstaged change:\n%s", out)
	}
	idx, _ := exec.Command("git", "-C", root, "show", ":a.ts").Output()
	if !strings.Contains(string(idx), "updated_at: 2026-09-25") {
		t.Errorf("the bumped a.ts must be re-staged:\n%s", idx)
	}
	idxD, _ := exec.Command("git", "-C", root, "show", ":d.ts").Output()
	if strings.Contains(string(idxD), "export const x = 3") {
		t.Errorf("the unstaged change of d.ts must not reach the index:\n%s", idxD)
	}
}

// The pre-commit phase dates the staged files first — on by default, off with
// `touch.pre_commit: false`, and never outside the pre-commit over the index.
func TestPreCommitTouches(t *testing.T) {
	off := false
	on := true
	for _, c := range []struct {
		name   string
		staged bool
		phase  string
		cfg    *config.Config
		want   bool
	}{
		{"default: no config", true, "pre-commit", nil, true},
		{"default: no touch block", true, "pre-commit", &config.Config{}, true},
		{"declared on", true, "pre-commit", &config.Config{Touch: &config.Touch{PreCommit: &on}}, true},
		{"declared off", true, "pre-commit", &config.Config{Touch: &config.Touch{PreCommit: &off}}, false},
		{"not staged", false, "pre-commit", nil, false},
		{"another phase", true, "ci", nil, false},
	} {
		if got := preCommitTouches(c.staged, c.phase, c.cfg); got != c.want {
			t.Errorf("%s: preCommitTouches = %v, want %v", c.name, got, c.want)
		}
	}
}

// What the pre-commit prints: never a silent rewrite.
func TestPrintPreCommitTouch(t *testing.T) {
	out := captureStdout(t, func() {
		printPreCommitTouch([]touchDecision{{File: "a.ts"}, {File: "b.ts"}},
			map[touchSkip][]string{skipUnstaged: {"d.ts"}}, nil)
	})
	for _, want := range []string{"dated 2 staged file(s) today", "a.ts, b.ts", "d.ts not dated"} {
		if !strings.Contains(out, want) {
			t.Errorf("the pre-commit should say %q:\n%s", want, out)
		}
	}
	if out := captureStdout(t, func() { printPreCommitTouch(nil, map[touchSkip][]string{}, nil) }); out != "" {
		t.Errorf("nothing dated, nothing said; printed:\n%s", out)
	}
}

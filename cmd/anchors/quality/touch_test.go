// @anchors
//   code: TCTST
//   ref: HDTHD

package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
)

const touchHeader = "// @anchors\n//   code: CODEX\n//   updated_at: 2026-09-01\n"

// What `anchors touch` does with one file: it does not lie about what changed.
func TestDecideTouch(t *testing.T) {
	t.Run("HDTHD-B01: Only a date inside the header at the top of the file is bumped", func(t *testing.T) {})
	t.Run("HDTHD-B02: A real change and a new file are bumped to the date", func(t *testing.T) {})
	t.Run("HDTHD-B03: An unchanged file, a date-only change and a file already at the date are not bumped", func(t *testing.T) {})
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
		// The incident: header text inside a Go constant of a test file — not a header,
		// and the pre-commit rewrote it and broke the test.
		{"header text inside code", "package x\n\nconst h = \"// @anchors\\n//   updated_at: 2026-09-01\\n\"\n", "", false, false, skipNoHeader},
		{"HTML header of a spec", "<!-- @anchors\n  code: CODEX\n  updated_at: 2026-09-01\n-->\n# Spec\n", "", false, true, ""},
		{"header below the top", strings.Repeat("x\n", 12) + touchHeader, "", false, false, skipNoHeader},
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
	t.Run("HDTHD-B04: Without the staged flag the candidates are the worktree changes and the untracked files", func(t *testing.T) {})
	t.Run("HDTHD-B10: Each bump and each skip is listed with its reason, then the total", func(t *testing.T) {})
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
	t.Run("HDTHD-B05: With the staged flag the index is dated and re-staged", func(t *testing.T) {})
	t.Run("HDTHD-B06: A staged file with changes outside the index is skipped", func(t *testing.T) {})
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

// A PARTIAL COMMIT (`git commit -- <paths>`) runs the hook on a temporary index
// (`next-index-<pid>.lock`) and prepares the real one in `index.lock`. Staging only in the
// temporary index dated the commit and left the real index with the old date (`MM`,
// reported from the reference app). Both must end dated. The state git builds is simulated here: both
// indexes hold the pre-hook worktree version.
func TestTouch_stagedInAPartialCommit(t *testing.T) {
	t.Run("HDTHD-B07: In a partial commit the real index is dated too", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	gitDir := filepath.Join(root, ".git")
	temp := filepath.Join(gitDir, "next-index-4242.lock")
	real := filepath.Join(gitDir, "index.lock")
	for _, idx := range []string{temp, real} {
		b, err := os.ReadFile(filepath.Join(gitDir, "index"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(idx, b, 0o644); err != nil {
			t.Fatal(err)
		}
		c := exec.Command("git", "-C", root, "add", "a.ts")
		c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("stage in %s: %v %s", idx, err, out)
		}
	}
	t.Setenv("GIT_INDEX_FILE", temp)

	if _, _, err := touchRun(root, true, false, "2026-09-25", nil, nil); err != nil {
		t.Fatal(err)
	}
	for _, idx := range []string{temp, real} {
		c := exec.Command("git", "-C", root, "show", ":a.ts")
		c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
		out, _ := c.Output()
		if !strings.Contains(string(out), "updated_at: 2026-09-25") {
			t.Errorf("%s must hold the dated a.ts:\n%s", filepath.Base(idx), out)
		}
	}
}

func TestTouch_configExcludeAddsToTheFlag(t *testing.T) {
	t.Run("HDTHD-B08: The exclude globs of the flag and of the configuration add up", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "anchors.yaml", "version: 2\nlayers: {}\ntouch:\n  exclude: [\"gen/**\"]\n")
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, root, "gen/c.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, root, "d.ts", touchHeader+"export const x = 2\n")

	out := runTouch(t, "--root", root, "--date", "2026-09-25", "--exclude", "d.ts")
	for _, want := range []string{"bumped a.ts", "skipped gen/c.ts — " + string(skipExcluded), "skipped d.ts — " + string(skipExcluded)} {
		if !strings.Contains(out, want) {
			t.Errorf("the output should say %q:\n%s", want, out)
		}
	}
}

func TestTouch_dryRunWritesNothingAndDefaultsToToday(t *testing.T) {
	t.Run("HDTHD-B09: The dry run says what it would bump and writes nothing", func(t *testing.T) {})
	t.Run("HDTHD-B11: Without a date the day of the run is written", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	out := runTouch(t, "--root", root, "--dry-run")
	today := gitmeta.Today()
	if !strings.Contains(out, "would bump a.ts  (2026-09-01 → "+today+")") || !strings.Contains(out, "would bump 1 file(s) to "+today+".") {
		t.Errorf("the dry run names the bump with today's date:\n%s", out)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "a.ts")); !strings.Contains(string(b), "updated_at: 2026-09-01") {
		t.Errorf("the dry run must not write:\n%s", b)
	}
}

func TestTouch_secondRunBumpsNothing(t *testing.T) {
	t.Run("HDTHD-I01: Touching twice bumps nothing the second time", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	runTouch(t, "--root", root, "--date", "2026-09-25")
	out := runTouch(t, "--root", root, "--date", "2026-09-25")
	if !strings.Contains(out, "skipped a.ts — "+string(skipAlready)) || !strings.Contains(out, "bumped 0 file(s) to 2026-09-25.") {
		t.Errorf("a file already at the date is left alone:\n%s", out)
	}
}

func TestTouch_fileWithoutHeaderIsNotTouchedNorListed(t *testing.T) {
	t.Run("HDTHD-X01: A changed file with no dated header is neither touched nor listed", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "plain.ts", "export const p = 1\n")
	out := runTouch(t, "--root", root, "--date", "2026-09-25")
	if strings.Contains(out, "plain.ts") {
		t.Errorf("a file without a dated header is not the command's business:\n%s", out)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "plain.ts")); string(b) != "export const p = 1\n" {
		t.Errorf("plain.ts was changed:\n%s", b)
	}
}

func TestTouch_outsideGitFails(t *testing.T) {
	t.Run("HDTHD-E01: Outside a git repository touch fails", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	dir := t.TempDir()
	if temRepoAcima(dir) {
		t.Skipf("the temporary directory %s is inside a git repository", dir)
	}
	cmd := newTouchCmd()
	cmd.SetArgs([]string{"--root", dir})
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	var err error
	captureStdout(t, func() { err = cmd.Execute() })
	if err == nil || !strings.Contains(err.Error(), "git diff HEAD") {
		t.Errorf("without a repository there is no change to date; got %v", err)
	}
}

func TestTouch_unreadableCandidateIsSkipped(t *testing.T) {
	t.Run("HDTHD-E02: A changed file that cannot be read is skipped and named", func(t *testing.T) {})
	testkit.SkipWithoutPOSIXPermissions(t)
	if os.Getuid() == 0 {
		t.Skip("root reads everything")
	}
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, root, "b.ts", touchHeader+"export const x = 2\n")
	if err := os.Chmod(filepath.Join(root, "b.ts"), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(filepath.Join(root, "b.ts"), 0o644) })
	out := runTouch(t, "--root", root, "--date", "2026-09-25")
	if !strings.Contains(out, "skipped b.ts — "+string(skipUnreadable)) || !strings.Contains(out, "bumped a.ts") {
		t.Errorf("an unreadable file is named and the rest goes on:\n%s", out)
	}
}

func TestTouchOnPreCommitIsOnByDefault(t *testing.T) {
	t.Run("HDTHD-B12: The pre-commit bump is on unless the project turns it off", func(t *testing.T) {})
	off, on := false, true
	for _, c := range []struct {
		name string
		cfg  *config.Config
		want bool
	}{
		{"no config", nil, true},
		{"no touch block", &config.Config{}, true},
		{"touch block without pre_commit", &config.Config{Touch: &config.Touch{}}, true},
		{"declared on", &config.Config{Touch: &config.Touch{PreCommit: &on}}, true},
		{"declared off", &config.Config{Touch: &config.Touch{PreCommit: &off}}, false},
	} {
		if got := touchOnPreCommit(c.cfg); got != c.want {
			t.Errorf("%s: touchOnPreCommit = %v, want %v", c.name, got, c.want)
		}
	}
}

// The project root may be a subdirectory of the git repository (a monorepo package). git
// names the changed files from the repository's top, and the touch works relative to the
// project root: every path must be the project's, and files outside it are not its own.
func TestTouch_projectRootBelowTheRepositoryTop(t *testing.T) {
	t.Run("HDTHD-B13: A project root below the repository top dates its own files and nothing outside it", func(t *testing.T) {})
	repo := touchRepo(t)
	for _, f := range []string{"sub/x.ts", "sub/y.ts"} {
		touchWrite(t, repo, f, touchHeader+"export const x = 1\n")
	}
	if out, err := exec.Command("git", "-C", repo, "-c", "user.email=t@t", "-c", "user.name=t", "add", ".").CombinedOutput(); err != nil {
		t.Fatalf("git add: %v %s", err, out)
	}
	if out, err := exec.Command("git", "-C", repo, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "sub").CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v %s", err, out)
	}
	sub := filepath.Join(repo, "sub")

	// --staged: x staged clean, y staged with a change on top, a.ts (outside) staged too
	touchWrite(t, repo, "sub/x.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, repo, "sub/y.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, repo, "a.ts", touchHeader+"export const x = 2\n")
	if out, err := exec.Command("git", "-C", repo, "add", "sub/x.ts", "sub/y.ts", "a.ts").CombinedOutput(); err != nil {
		t.Fatalf("git add: %v %s", err, out)
	}
	touchWrite(t, repo, "sub/y.ts", touchHeader+"export const x = 3\n")

	out := runTouch(t, "--root", sub, "--staged", "--date", "2026-09-25")
	if !strings.Contains(out, "bumped x.ts") || !strings.Contains(out, "skipped y.ts — "+string(skipUnstaged)) {
		t.Errorf("staged under a sub-root: x.ts bumped, y.ts skipped:\n%s", out)
	}
	if strings.Contains(out, "a.ts") {
		t.Errorf("a file outside the project root is not the project's:\n%s", out)
	}
	idx, _ := exec.Command("git", "-C", repo, "show", ":sub/x.ts").Output()
	if !strings.Contains(string(idx), "updated_at: 2026-09-25") {
		t.Errorf("the bumped sub/x.ts must be re-staged:\n%s", idx)
	}
	if _, err := os.Stat(filepath.Join(sub, "sub")); err == nil {
		t.Error("a path from the repository top was joined to the project root")
	}

	// the worktree mode: the changed and the untracked files of the project, by its paths
	touchWrite(t, repo, "sub/new.ts", touchHeader+"export const n = 1\n")
	out = runTouch(t, "--root", sub, "--date", "2026-09-26", "--dry-run")
	for _, want := range []string{"x.ts", "new.ts"} {
		if !strings.Contains(out, "would bump "+want) {
			t.Errorf("the worktree mode under a sub-root should bump %s:\n%s", want, out)
		}
	}
	if strings.Contains(out, "a.ts") {
		t.Errorf("a file outside the project root is not the project's:\n%s", out)
	}
}

func TestTouch_namedFilesNarrowIt(t *testing.T) {
	t.Run("HDTHD-B14: Named files narrow the touch to them", func(t *testing.T) {})
	root := touchRepo(t)
	for _, f := range []string{"a.ts", "b.ts", "gen/c.ts"} {
		touchWrite(t, root, f, touchHeader+"export const x = 2\n")
	}
	t.Chdir(root)
	cmd := newTouchCmd()
	cmd.SetArgs([]string{"a.ts", filepath.Join(root, "gen"), "--date", "2026-09-25"})
	var err error
	testkit.CaptureStdout(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	for f, dated := range map[string]bool{"a.ts": true, "gen/c.ts": true, "b.ts": false} {
		b, _ := os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
		if strings.Contains(string(b), "2026-09-25") != dated {
			t.Errorf("%s dated %v, want %v:\n%s", f, !dated, dated, b)
		}
	}
	cmd = newTouchCmd()
	cmd.SetArgs([]string{filepath.Join(t.TempDir(), "x.ts")})
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	testkit.CaptureStdout(t, func() { err = cmd.Execute() })
	if err == nil || !strings.Contains(err.Error(), "outside the project") {
		t.Errorf("a path outside the project is an error, got %v", err)
	}
	bumped, _, err := touchRun(root, false, true, "2026-09-26", nil, nil)
	if err != nil || len(bumped) != 3 {
		t.Errorf("with nothing named every changed file is a candidate, got %d %v", len(bumped), err)
	}
}

func TestTouch_everyNamedFileGetsAVerdict(t *testing.T) {
	t.Run("HDTHD-B15: Every named file gets a verdict", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	touchWrite(t, root, "plain.ts", "export const y = 1\n")
	bumped, skipped, err := touchRun(root, false, false, "2026-09-25", nil, []string{"a.ts", "b.ts", "plain.ts", "nope.ts", "gen"})
	if err != nil {
		t.Fatal(err)
	}
	if len(bumped) != 1 || bumped[0].File != "a.ts" {
		t.Errorf("the changed file is dated, got %+v", bumped)
	}
	if got := skipped[skipUnchanged]; len(got) != 1 || got[0] != "b.ts" {
		t.Errorf("the unchanged file says so, got %v", skipped)
	}
	if got := skipped[skipUnreadable]; len(got) != 1 || got[0] != "nope.ts" {
		t.Errorf("the missing file says so, got %v", skipped)
	}
	if got := skipped[skipNoHeader]; len(got) != 1 || got[0] != "plain.ts" {
		t.Errorf("the named file with no header says so, got %v", skipped)
	}
	if _, skipped, _ := touchRun(root, false, true, "2026-09-25", nil, nil); len(skipped[skipNoHeader]) != 0 {
		t.Errorf("with nothing named a file with no header stays unlisted, got %v", skipped[skipNoHeader])
	}
	for _, list := range skipped {
		for _, f := range list {
			if strings.HasPrefix(f, "gen") {
				t.Errorf("an unchanged folder brings nothing, got %s", f)
			}
		}
	}
}

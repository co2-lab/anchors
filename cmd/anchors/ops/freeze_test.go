package ops

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/config"
)

const githubModeConfig = "version: 1\nworkflow:\n  mode: github\n  repo: acme/app\n  labels: [anchors]\n"

// frozenRepo is a repository with anchors.yaml committed and pushed to a local bare
// remote, so the freeze's own push has somewhere to go.
func frozenRepo(t *testing.T, cfg string) (root, remote string) {
	t.Helper()
	root = newGitRepo(t)
	remote = t.TempDir()
	writeFile(t, root, config.DefaultFile, cfg)
	for _, args := range [][]string{
		{"init", "-q", "--bare", remote},
		{"remote", "add", "origin", remote},
		{"add", config.DefaultFile},
		{"commit", "-q", "-m", "config"},
		{"push", "-q", "-u", "origin", "main"},
	} {
		dir := root
		if args[0] == "init" {
			dir = remote
		}
		if out, err := runGit(dir, args...); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
	return root, remote
}

func runCmd(t *testing.T, c *cobra.Command, args ...string) (error, string) {
	t.Helper()
	c.SetArgs(args)
	c.SetOut(io.Discard)
	c.SetErr(io.Discard)
	var err error
	out := captureStdout(t, func() { err = c.Execute() })
	return err, out
}

func remoteSubject(t *testing.T, remote string) string {
	t.Helper()
	out, err := runGit(remote, "log", "-1", "--format=%s", "main")
	if err != nil {
		t.Fatalf("git log on the remote: %s", out)
	}
	return out
}

// The reason is mandatory: without it, the freeze is indistinguishable from broken config.
func TestFreezeRequiresAReason(t *testing.T) {
	root, _ := frozenRepo(t, "version: 1\n")
	err, _ := runCmd(t, newFreezeCmd(), "--root", root, "--reason", "  ")
	if err == nil || !strings.Contains(err.Error(), "--reason is required") {
		t.Fatalf("want the reason refusal, got %v", err)
	}
	b, _ := os.ReadFile(filepath.Join(root, config.DefaultFile))
	if strings.Contains(string(b), "enabled:") {
		t.Error("the file was frozen without a reason")
	}
}

// The three layers in one operation: the file (pushed), the ruleset, and the issue.
// Then `thaw` undoes all three.
func TestFreezeAndThawDriveTheThreeLayers(t *testing.T) {
	root, remote := frozenRepo(t, githubModeConfig)
	log := fakeGH(t, `case "$*" in
"issue create"*) echo "https://github.com/acme/app/issues/99" ;;
*"--method POST"*) cat > /dev/null ;;
*"rulesets --jq"*) echo 42 ;;
"issue list"*) echo 99 ;;
esac`)

	reason := "leaked credential: rotate first"
	err, out := runCmd(t, newFreezeCmd(), "--root", root, "--reason", reason)
	if err != nil {
		t.Fatalf("freeze: %v\n%s", err, out)
	}
	cfg, lerr := config.Load(filepath.Join(root, config.DefaultFile))
	if lerr != nil {
		t.Fatalf("the frozen file no longer loads (the reason has a colon): %v", lerr)
	}
	if !cfg.Frozen() || cfg.FreezeReason != reason {
		t.Errorf("frozen = %v, reason = %q", cfg.Frozen(), cfg.FreezeReason)
	}
	if got := remoteSubject(t, remote); !strings.HasPrefix(got, "chore(anchors): freeze — leaked credential") {
		t.Errorf("the freeze did not reach the remote; last subject = %q", got)
	}
	calls := readLog(t, log)
	if !strings.Contains(calls, "api --method POST repos/acme/app/rulesets --input -") {
		t.Errorf("no ruleset was created:\n%s", calls)
	}
	if !strings.Contains(calls, "issue create --repo acme/app --title "+freezeIssueTitle) {
		t.Errorf("no freeze issue was opened:\n%s", calls)
	}
	if !strings.Contains(out, "freeze issue: https://github.com/acme/app/issues/99") {
		t.Errorf("the issue URL is not reported:\n%s", out)
	}

	// Freezing again changes nothing and says so.
	err, out = runCmd(t, newFreezeCmd(), "--root", root, "--reason", "other")
	if err != nil || !strings.Contains(out, "ALREADY frozen: "+reason) {
		t.Errorf("a second freeze: %v\n%s", err, out)
	}

	err, out = runCmd(t, newThawCmd(), "--root", root)
	if err != nil {
		t.Fatalf("thaw: %v\n%s", err, out)
	}
	b, _ := os.ReadFile(filepath.Join(root, config.DefaultFile))
	if string(b) != githubModeConfig {
		t.Errorf("thaw must remove exactly the two lines; file:\n%s", b)
	}
	if got := remoteSubject(t, remote); !strings.HasPrefix(got, "chore(anchors): thaw") {
		t.Errorf("the thaw did not reach the remote; last subject = %q", got)
	}
	calls = readLog(t, log)
	if !strings.Contains(calls, "api --method DELETE repos/acme/app/rulesets/42") {
		t.Errorf("the ruleset found by name was not deleted:\n%s", calls)
	}
	if !strings.Contains(calls, "issue close 99 --repo acme/app") {
		t.Errorf("the freeze issue was not closed:\n%s", calls)
	}

	err, out = runCmd(t, newThawCmd(), "--root", root)
	if err != nil || !strings.Contains(out, "not frozen") {
		t.Errorf("thawing twice: %v\n%s", err, out)
	}
}

// In local mode there is no remote to lock: no ruleset, no issue. --no-push keeps the
// change local, uncommitted.
func TestFreezeLocalModeNoPushTouchesOnlyTheFile(t *testing.T) {
	root, remote := frozenRepo(t, "version: 1\n")
	log := fakeGH(t, `exit 1`)
	err, out := runCmd(t, newFreezeCmd(), "--root", root, "--reason", "stop", "--no-push")
	if err != nil {
		t.Fatalf("freeze: %v\n%s", err, out)
	}
	if calls := readLog(t, log); calls != "" {
		t.Errorf("local mode called gh:\n%s", calls)
	}
	if got := remoteSubject(t, remote); got != "config" {
		t.Errorf("--no-push still pushed: %q", got)
	}
	if st, _ := runGit(root, "status", "--porcelain"); !strings.Contains(st, config.DefaultFile) {
		t.Errorf("the frozen file should be left as a local change; status = %q", st)
	}
}

// Each failing layer is WARNED, and the freeze of the file still holds: a brake that
// half-applies must say which half.
func TestFreezeAndThawWarnWhenALayerFails(t *testing.T) {
	root := newGitRepo(t) // no remote: the push fails
	writeFile(t, root, config.DefaultFile, githubModeConfig)
	fakeGH(t, `echo "forbidden" >&2; exit 1`)

	err, out := runCmd(t, newFreezeCmd(), "--root", root, "--reason", "x")
	if err != nil {
		t.Fatalf("a failing layer must not fail the freeze: %v", err)
	}
	for _, want := range []string{"could not push", "ruleset not created", "issue not opened"} {
		if !strings.Contains(out, want) {
			t.Errorf("no warning %q:\n%s", want, out)
		}
	}
	if cfg, _ := config.Load(filepath.Join(root, config.DefaultFile)); !cfg.Frozen() {
		t.Error("the local brake must hold even when the remote layers fail")
	}

	err, out = runCmd(t, newThawCmd(), "--root", root)
	if err != nil {
		t.Fatalf("thaw: %v", err)
	}
	for _, want := range []string{"could not push", "ruleset not removed", "issue not closed"} {
		if !strings.Contains(out, want) {
			t.Errorf("no warning %q:\n%s", want, out)
		}
	}
}

// No ruleset and no open issue is not an error: there is simply nothing to remove.
func TestThawWithNothingOnTheRemoteRemovesNothing(t *testing.T) {
	log := fakeGH(t, `exit 0`)
	if err := deleteRuleset("acme/app"); err != nil {
		t.Errorf("deleteRuleset: %v", err)
	}
	if err := closeFreezeIssue("acme/app"); err != nil {
		t.Errorf("closeFreezeIssue: %v", err)
	}
	calls := readLog(t, log)
	if strings.Contains(calls, "DELETE") || strings.Contains(calls, "issue close") {
		t.Errorf("something was removed although nothing was found:\n%s", calls)
	}
}

func TestFreezeWithoutConfigFails(t *testing.T) {
	for _, c := range []*cobra.Command{newFreezeCmd(), newThawCmd()} {
		args := []string{"--root", t.TempDir()}
		if c.Name() == "freeze" {
			args = append(args, "--reason", "x")
		}
		if err, _ := runCmd(t, c, args...); err == nil || !strings.Contains(err.Error(), "load anchors.yaml") {
			t.Errorf("%s without anchors.yaml: %v", c.Name(), err)
		}
	}
}

// removeFreeze removes the two lines and ONLY them, wherever they are.
func TestRemoveFreezeKeepsEveryOtherLine(t *testing.T) {
	p := writeFile(t, t.TempDir(), "anchors.yaml", "version: 1\n")
	if err := writeFreeze(p, `a: "quoted" reason`); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if !strings.HasPrefix(string(b), "enabled: false\nfreeze_reason: \"a: \\\"quoted\\\" reason\"\n") {
		t.Errorf("the freeze block is not on top, quoted:\n%s", b)
	}
	if err := removeFreeze(p); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(p); string(b) != "version: 1\n" {
		t.Errorf("after removeFreeze: %q", b)
	}
	if writeFreeze(filepath.Join(t.TempDir(), "missing"), "x") == nil || removeFreeze(filepath.Join(t.TempDir(), "missing")) == nil {
		t.Error("a missing file must be an error")
	}
}

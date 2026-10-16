// @anchors
//   code: FRTSF
//   ref: FRZEX

package ops

import (
	"github.com/co2-lab/anchors/internal/testkit"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func strconvQuote(s string) string { return strconv.Quote(s) }

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
	t.Run("FRZEX-B01: A blank reason is refused", func(t *testing.T) {})
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
	t.Run("FRZEX-B02: The freeze writes two quoted lines on top and the file still loads", func(t *testing.T) {})
	t.Run("FRZEX-B05: In github mode the freeze creates the rule and opens the issue", func(t *testing.T) {})
	t.Run("FRZEX-B07: Freezing twice changes nothing and shows the reason", func(t *testing.T) {})
	t.Run("FRZEX-B09: The thaw undoes the three layers", func(t *testing.T) {})
	t.Run("FRZEX-B10: Thawing a project that is not frozen is a no-op", func(t *testing.T) {})
	t.Run("FRZEX-I01: Freeze and thaw give back the file byte for byte", func(t *testing.T) {})
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
	t.Run("FRZEX-B04: No-push leaves the frozen file as a local change", func(t *testing.T) {})
	t.Run("FRZEX-B05: In github mode the freeze creates the rule and opens the issue", func(t *testing.T) {})
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
	t.Run("FRZEX-B08: Each failing layer is a warning and the local brake holds", func(t *testing.T) {})
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
	t.Run("FRZEX-B11: With nothing on the remote the thaw removes nothing", func(t *testing.T) {})
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
	t.Run("FRZEX-E01: Freeze and thaw refuse a project without config", func(t *testing.T) {})
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
	t.Run("FRZEX-B02: The freeze writes two quoted lines on top and the file still loads", func(t *testing.T) {})
	t.Run("FRZEX-X01: The configuration is never reserialized", func(t *testing.T) {})
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

// The freeze commits and pushes PAST the hooks: the hooks it has just armed would refuse
// the freeze itself. The deprecated Portuguese flags still drive the command, and
// --no-ruleset leaves the remote rules alone while the issue is still opened.
func TestFreezeGoesPastTheHooksAndHonorsTheOldFlagNames(t *testing.T) {
	t.Run("FRZEX-B03: The freeze is committed and pushed past refusing hooks", func(t *testing.T) {})
	t.Run("FRZEX-B06: No-ruleset skips the ruleset and still opens the issue", func(t *testing.T) {})
	t.Run("FRZEX-B12: The deprecated Portuguese flag names still work", func(t *testing.T) {})
	root, remote := frozenRepo(t, githubModeConfig)
	for _, hook := range []string{"pre-commit", "commit-msg", "pre-push"} {
		writeFile(t, root, filepath.Join(".git", "hooks", hook), "#!/bin/sh\necho refused >&2\nexit 1\n")
		if err := os.Chmod(filepath.Join(root, ".git", "hooks", hook), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	log := fakeGH(t, `case "$*" in
"issue create"*) echo "https://github.com/acme/app/issues/7" ;;
esac`)
	err, out := runCmd(t, newFreezeCmd(), "--root", root, "--motivo", "release blocked", "--sem-ruleset")
	if err != nil {
		t.Fatalf("freeze with the old flag names: %v\n%s", err, out)
	}
	if got := remoteSubject(t, remote); got != "chore(anchors): freeze — release blocked" {
		t.Errorf("the freeze did not go past the hooks to the remote; last subject = %q\n%s", got, out)
	}
	calls := readLog(t, log)
	if strings.Contains(calls, "rulesets") {
		t.Errorf("--sem-ruleset still touched the ruleset:\n%s", calls)
	}
	if !strings.Contains(calls, "issue create") {
		t.Errorf("the issue must be opened even without the ruleset:\n%s", calls)
	}

	// --sem-push on thaw: the change stays local.
	err, out = runCmd(t, newThawCmd(), "--root", root, "--sem-push")
	if err != nil {
		t.Fatalf("thaw --sem-push: %v\n%s", err, out)
	}
	if got := remoteSubject(t, remote); !strings.HasPrefix(got, "chore(anchors): freeze") {
		t.Errorf("thaw --sem-push still pushed: %q", got)
	}
}

// A config that already declares `enabled: true` (and a stale freeze_reason) used to get a
// SECOND `enabled:` key on top: yaml refuses the duplicate, the config stopped loading, and
// with it every command — including the thaw that would undo it. The freeze now replaces
// the top-level lines instead of stacking on them.
func TestFreezeReplacesAnExistingEnabledLine(t *testing.T) {
	t.Run("FRZEX-B13: A config that already declares enabled is frozen with a single key and still loads", func(t *testing.T) {})
	const orig = "version: 1\nenabled: true\nfreeze_reason: \"old\"\nlang: en\n"
	root, _ := frozenRepo(t, orig)
	err, out := runCmd(t, newFreezeCmd(), "--root", root, "--reason", "stop now", "--no-push")
	if err != nil {
		t.Fatalf("freeze: %v\n%s", err, out)
	}
	p := filepath.Join(root, config.DefaultFile)
	b, _ := os.ReadFile(p)
	if n := strings.Count(string(b), "enabled:"); n != 1 {
		t.Errorf("want exactly one enabled: key, got %d:\n%s", n, b)
	}
	if n := strings.Count(string(b), "freeze_reason:"); n != 1 {
		t.Errorf("want exactly one freeze_reason: key, got %d:\n%s", n, b)
	}
	cfg, lerr := config.Load(p)
	if lerr != nil {
		t.Fatalf("the frozen config no longer loads: %v\n%s", lerr, b)
	}
	if !cfg.Frozen() || cfg.FreezeReason != "stop now" {
		t.Errorf("frozen = %v, reason = %q", cfg.Frozen(), cfg.FreezeReason)
	}

	// The thaw takes both lines away, and the project loads enabled (the default).
	if err, out := runCmd(t, newThawCmd(), "--root", root, "--no-push"); err != nil {
		t.Fatalf("thaw: %v\n%s", err, out)
	}
	b, _ = os.ReadFile(p)
	if string(b) != "version: 1\nlang: en\n" {
		t.Errorf("after thaw: %q", b)
	}
	if cfg, lerr := config.Load(p); lerr != nil || cfg.Frozen() {
		t.Errorf("after thaw: load err = %v, frozen = %v", lerr, cfg != nil && cfg.Frozen())
	}
	// A hand-written block reason goes with its continuation lines, not leaving them orphaned.
	writeFile(t, root, config.DefaultFile, "version: 1\nenabled: true\nfreeze_reason: |\n  multi\n  line\nlang: en\n")
	if err := writeFreeze(p, "x"); err != nil {
		t.Fatal(err)
	}
	if cfg, lerr := config.Load(p); lerr != nil || cfg.FreezeReason != "x" {
		b, _ := os.ReadFile(p)
		t.Errorf("block reason: load err = %v\n%s", lerr, b)
	}
}

// --- the freeze with a scope ---

var scopedYAML = "version: 8\ndialect:\n" + testkit.CommentsYAML + "layers:\n" +
	"  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n" +
	"  code:\n    pattern: \"src/**/*.ts\"\n    kind: code\n" +
	"  produto:\n    pattern: \"product/*.doctrine.md\"\n    kind: product\n" +
	"derived:\n  anchor: spec\n  files:\n    code: [\"{{dir}}/{{name}}.ts\"]\n"

// scopedProject is a project with two units, a file of its own and, unless bare, one
// doctrine document.
func scopedProject(t *testing.T, doctrine bool) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		config.DefaultFile:      scopedYAML,
		"src/pay/Pay.spec.md":   "<!-- @anchors\n  code: PAYUN\n  layer: spec\n-->\n# Pay\n\n## Effects\n\n| `PAYUN-B01` | pays |\n",
		"src/pay/Pay.ts":        "// @anchors\n//   code: PAYCD\n//   ref: PAYUN\n\nexport const pay = 1\n",
		"src/cart/Cart.spec.md": "<!-- @anchors\n  code: CARTU\n-->\n# Cart\n\n### CARTU-B01 — coupons\n\n| `CARTU-B02` | totals |\n",
		"src/flags.ts":          "// @anchors\n//   code: FLAGS\n\nexport const on = false\n",
	}
	if doctrine {
		files["product/q.doctrine.md"] = "<!-- @anchors\n  code: QLDTQ\n-->\n# Q\n\n## Rules\n\n### QLDTQ-R01 — a rule\n"
	}
	for rel, body := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// runScoped runs a command and gives what it printed.
func runScoped(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var b strings.Builder
	c.SetArgs(args)
	c.SetOut(&b)
	c.SetErr(io.Discard)
	err := c.Execute()
	return b.String(), err
}

func readScoped(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestFreezeScoped_recordsAndMarks(t *testing.T) {
	t.Run("FRZEX-B14: A freeze with a scope is recorded in the doctrine and marks every place", func(t *testing.T) {})
	root := scopedProject(t, true)
	out, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "the plan mode ships next quarter",
		"--unit", "PAYUN", "--rule", "CARTU-B01", "--rule", "CARTU-B02", "--file", "src/flags.ts")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	for rel, want := range map[string]string{
		"product/q.doctrine.md": "## Freezes\n\n- **QLDTQ-F01** — the plan mode ships next quarter\n",
		"src/pay/Pay.spec.md":   "  code: PAYUN\n  frozen: QLDTQ-F01\n",
		"src/cart/Cart.spec.md": "### CARTU-B01 — coupons @frozen: QLDTQ-F01\n\n| `CARTU-B02` | totals @frozen: QLDTQ-F01 |\n",
		"src/flags.ts":          "//   code: FLAGS\n//   frozen: QLDTQ-F01\n",
	} {
		if got := readScoped(t, root, rel); !strings.Contains(got, want) {
			t.Errorf("%s:\n%s\nwant %q", rel, got, want)
		}
	}
	if !strings.Contains(out, "QLDTQ-F01 recorded in product/q.doctrine.md") {
		t.Errorf("says what it recorded:\n%s", out)
	}
	// The map reads what the command wrote: the unit's code is covered through its spec.
	g, err := freezeGraph(root)
	if err != nil {
		t.Fatal(err)
	}
	parts := map[string]string{}
	for _, p := range g.FrozenParts("QLDTQ-F01") {
		parts[p.Path+" "+p.Rule] = p.Value
	}
	for _, want := range []string{"src/pay/Pay.spec.md ", "src/pay/Pay.ts ", "src/flags.ts ", "src/cart/Cart.spec.md CARTU-B01", "src/cart/Cart.spec.md CARTU-B02"} {
		if parts[want] != "QLDTQ-F01" {
			t.Errorf("the map misses %q: %v", want, parts)
		}
	}
	// A second freeze takes the next code.
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "waits for legal", "--rule", "PAYUN-B01"); err != nil {
		t.Fatal(err)
	}
	if got := readScoped(t, root, "product/q.doctrine.md"); !strings.Contains(got, "- **QLDTQ-F01** — the plan mode ships next quarter\n- **QLDTQ-F02** — waits for legal\n") {
		t.Errorf("the next freeze:\n%s", got)
	}
}

func TestFreezeScoped_refusalsAndNoDoctrine(t *testing.T) {
	t.Run("FRZEX-B15: Without a doctrine the marks carry the reason, and what cannot be frozen is refused", func(t *testing.T) {})
	bare := scopedProject(t, false)
	out, err := runScoped(t, newFreezeCmd(), "--root", bare, "--reason", "flag off until launch", "--file", "src/flags.ts")
	if err != nil || !strings.Contains(out, "not traceable") || !strings.Contains(readScoped(t, bare, "src/flags.ts"), "//   frozen: flag off until launch\n") {
		t.Errorf("no doctrine: %v\n%s", err, out)
	}
	root := scopedProject(t, true)
	if err := os.WriteFile(filepath.Join(root, "product", "r.doctrine.md"), []byte("<!-- @anchors\n  code: RLDTR\n-->\n# R\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--unit", "PAYUN"}, "needs its reason"},
		{[]string{"--reason", "x", "--unit", "PAYUN"}, "--in"},
		{[]string{"--reason", "x", "--unit", "NOPEX", "--in", "RLDTR"}, "NOPEX has no owner"},
		{[]string{"--reason", "x", "--rule", "PAYUN-B09", "--in", "RLDTR"}, "PAYUN-B09 is not defined"},
		{[]string{"--reason", "x", "--file", "src/none.ts", "--in", "RLDTR"}, "does not exist"},
	} {
		if _, err := runScoped(t, newFreezeCmd(), append([]string{"--root", root}, c.args...)...); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: %v, want %q", c.args, err, c.want)
		}
	}
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "x", "--file", "src/flags.ts", "--in", "product/r.doctrine.md"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readScoped(t, root, "product/r.doctrine.md"), "- **RLDTR-F01** — x") {
		t.Error("--in records the freeze in the document it names")
	}
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "y", "--file", "src/flags.ts", "--in", "RLDTR"); err == nil || !strings.Contains(err.Error(), "already frozen") {
		t.Errorf("a part frozen by another freeze is refused: %v", err)
	}
}

func TestFreezeScoped_addScope(t *testing.T) {
	t.Run("FRZEX-B16: A recorded freeze gains a scope by its code", func(t *testing.T) {})
	root := scopedProject(t, true)
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "the plan mode", "--unit", "PAYUN"); err != nil {
		t.Fatal(err)
	}
	if _, err := runScoped(t, newFreezeCmd(), "QLDTQ-F01", "--root", root, "--unit", "CARTU"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readScoped(t, root, "src/cart/Cart.spec.md"), "  code: CARTU\n  frozen: QLDTQ-F01\n") || strings.Contains(readScoped(t, root, "product/q.doctrine.md"), "F02") {
		t.Errorf("the scope is added, no freeze recorded:\n%s", readScoped(t, root, "src/cart/Cart.spec.md"))
	}
	if _, err := runScoped(t, newFreezeCmd(), "QLDTQ-F09", "--root", root, "--unit", "CARTU"); err == nil || !strings.Contains(err.Error(), "QLDTQ-F09 is not a freeze") {
		t.Errorf("an unrecorded code is refused: %v", err)
	}
}

func TestFreezeScoped_thaw(t *testing.T) {
	t.Run("FRZEX-B17: Thawing a freeze removes its marks and its record, and nothing else", func(t *testing.T) {})
	root := scopedProject(t, true)
	before := readScoped(t, root, "src/pay/Pay.spec.md")
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "the plan mode", "--unit", "PAYUN", "--rule", "CARTU-B02", "--file", "src/flags.ts"); err != nil {
		t.Fatal(err)
	}
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "legal", "--rule", "CARTU-B01"); err != nil {
		t.Fatal(err)
	}
	out, err := runScoped(t, newThawCmd(), "QLDTQ-F01", "--root", root)
	if err != nil || !strings.Contains(out, "QLDTQ-F01 thawed: 3 mark(s) removed") {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := readScoped(t, root, "src/pay/Pay.spec.md"); got != before {
		t.Errorf("the unit's spec comes back as it was:\n%s", got)
	}
	cart := readScoped(t, root, "src/cart/Cart.spec.md")
	if !strings.Contains(cart, "| `CARTU-B02` | totals |\n") || !strings.Contains(cart, "### CARTU-B01 — coupons @frozen: QLDTQ-F02") {
		t.Errorf("the other freeze's mark stays:\n%s", cart)
	}
	if doc := readScoped(t, root, "product/q.doctrine.md"); strings.Contains(doc, "QLDTQ-F01") || !strings.Contains(doc, "QLDTQ-F02") {
		t.Errorf("only its record leaves:\n%s", doc)
	}
	if _, err := runScoped(t, newThawCmd(), "--root", root, "--rule", "CARTU-B01"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(readScoped(t, root, "src/cart/Cart.spec.md"), "@frozen") {
		t.Error("a thaw by scope removes the mark at the place")
	}
	if _, err := runScoped(t, newThawCmd(), "QLDTQ-F02", "--root", root); err != nil {
		t.Fatal(err)
	}
	if got := readScoped(t, root, "product/q.doctrine.md"); got != "<!-- @anchors\n  code: QLDTQ\n-->\n# Q\n\n## Rules\n\n### QLDTQ-R01 — a rule\n" {
		t.Errorf("the last thaw gives the document back as it was:\n%q", got)
	}
	if _, err := runScoped(t, newThawCmd(), "QLDTQ-F07", "--root", root); err == nil {
		t.Error("a code neither recorded nor marked is refused")
	}
}

func TestFreezeScoped_show(t *testing.T) {
	t.Run("FRZEX-B18: The freezes are listed, each with its reason and places", func(t *testing.T) {})
	root := scopedProject(t, true)
	if out, _ := runScoped(t, newFreezeCmd(), "show", "--root", root); !strings.Contains(out, "nothing is frozen") {
		t.Errorf("nothing frozen: %s", out)
	}
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "the plan mode", "--unit", "PAYUN"); err != nil {
		t.Fatal(err)
	}
	if err := addFreezeEntry(root, "product/q.doctrine.md", "QLDTQ-F02", "orphan"); err != nil {
		t.Fatal(err)
	}
	out, err := runScoped(t, newFreezeCmd(), "show", "--root", root)
	if err != nil || !strings.Contains(out, "❄ QLDTQ-F01 — the plan mode\n  src/pay/Pay.spec.md\n  src/pay/Pay.ts\n") || !strings.Contains(out, "❄ QLDTQ-F02 — orphan\n  (marks nothing)") {
		t.Errorf("%v\n%s", err, out)
	}
	if out, _ := runScoped(t, newFreezeCmd(), "show", "QLDTQ-F01", "--root", root); out != "QLDTQ-F01 — the plan mode\n  src/pay/Pay.spec.md\n  src/pay/Pay.ts\n" {
		t.Errorf("one freeze:\n%q", out)
	}
}

func TestFreezeScoped_hook(t *testing.T) {
	t.Run("FRZEX-B19: Before an edit, the agent is told a file is frozen, and the edit goes on", func(t *testing.T) {})
	root := scopedProject(t, true)
	if _, err := runScoped(t, newFreezeCmd(), "--root", root, "--reason", "the plan mode", "--file", "src/flags.ts", "--rule", "CARTU-B02"); err != nil {
		t.Fatal(err)
	}
	g, err := freezeGraph(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	call := func(path string) string {
		return freezeHook([]byte(`{"hook_event_name":"PreToolUse","cwd":` + strconvQuote(root) + `,"tool_input":{"file_path":` + strconvQuote(path) + `}}`))
	}
	flags := call(filepath.Join(root, "src", "flags.ts"))
	if !strings.Contains(flags, `"additionalContext":"FROZEN: src/flags.ts is frozen (QLDTQ-F01 — the plan mode)`) || !strings.Contains(flags, "This edit is NOT blocked and goes through: if no person ordered it, undo it and ask.") ||
		strings.Contains(flags, "permissionDecision") {
		t.Errorf("a frozen file is told, never denied: %s", flags)
	}
	if cart := call("src/cart/Cart.spec.md"); !strings.Contains(cart, "holds frozen rules") || !strings.Contains(cart, "CARTU-B02 (QLDTQ-F01 — the plan mode)") {
		t.Errorf("a file with frozen rules names them: %s", cart)
	}
	for _, none := range []string{call(filepath.Join(root, "src", "pay", "Pay.ts")), call(filepath.Join(t.TempDir(), "x.ts")), freezeHook([]byte("{broken"))} {
		if none != "" {
			t.Errorf("nothing to tell: %s", none)
		}
	}
}

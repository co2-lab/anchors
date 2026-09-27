package quality

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// `verify` re-executes its own binary as `check`. Under `go test` that binary is the test
// binary, so TestMain plays the child: when the variable is set it records the argv it
// was given and exits with the requested code, instead of running the tests.
func TestMain(m *testing.M) {
	if out := os.Getenv("ANCHORS_TEST_CHILD_ARGS"); out != "" {
		_ = os.WriteFile(out, []byte(strings.Join(os.Args[1:], "\n")), 0o644)
		code, _ := strconv.Atoi(os.Getenv("ANCHORS_TEST_CHILD_EXIT"))
		os.Exit(code)
	}
	os.Exit(m.Run())
}

// fakeChild makes the next runSubcommand record its argv (one per line) in the returned
// file and exit with the given code.
func fakeChild(t *testing.T, exit int) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "child-args")
	t.Setenv("ANCHORS_TEST_CHILD_ARGS", out)
	t.Setenv("ANCHORS_TEST_CHILD_EXIT", strconv.Itoa(exit))
	return out
}

func gitIn(t *testing.T, root string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", root, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
}

// The pre-commit mode, in a real repository: the staged file is dated first, then handed
// to `check` with the flags of an automatic phase.
func TestVerifyStagedDelegatesToCheck(t *testing.T) {
	t.Run("VPFVR-B03: Verify hands the files to check in a child process", func(t *testing.T) {})
	t.Run("VPFVR-B06: The pre-commit over the index dates the staged files first", func(t *testing.T) {})
	t.Run("VPFVR-X01: Verify holds no verdict of its own", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	gitIn(t, root, "add", "a.ts")
	argsFile := fakeChild(t, 0)

	out, err := runQ(t, newVerifyCmd(), "--root", root, "--phase", "pre-commit", "--staged")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "· touch: dated 1 staged file(s) today (updated_at): a.ts") {
		t.Errorf("the pre-commit dates the staged file before checking it:\n%s", out)
	}
	child := strings.Split(readQ(t, argsFile), "\n")
	want := []string{"check", "--root", root, "--changed", "a.ts", "--phase", "pre-commit", "--deterministic", "--only-issues"}
	if strings.Join(child, " ") != strings.Join(want, " ") {
		t.Errorf("the child check got %q, want %q", child, want)
	}
}

// Declared off, the pre-commit dates nothing; and a staged file that also has changes
// outside the index is named instead of dated.
func TestVerifyPreCommitDatingCanBeTurnedOffAndSkipsUnstaged(t *testing.T) {
	t.Run("VPFVR-B07: A project that turned pre-commit dating off gets no dating", func(t *testing.T) {})
	t.Run("VPFVR-B08: A staged file with changes outside the index is named, not dated", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 2\n")
	gitIn(t, root, "add", "a.ts")
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 3\n")
	fakeChild(t, 0)

	out, err := runQ(t, newVerifyCmd(), "--root", root, "--phase", "pre-commit", "--staged")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "· touch: a.ts not dated — it has changes outside the index") || strings.Contains(out, "dated 1") {
		t.Errorf("a file with unstaged changes must be named and not dated:\n%s", out)
	}

	off := touchRepo(t)
	touchWrite(t, off, "anchors.yaml", "version: 2\nlayers: {}\ntouch:\n  pre_commit: false\n")
	touchWrite(t, off, "b.ts", touchHeader+"export const x = 2\n")
	gitIn(t, off, "add", "b.ts")
	out, err = runQ(t, newVerifyCmd(), "--root", off, "--phase", "pre-commit", "--staged")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "· touch:") {
		t.Errorf("with touch.pre_commit false nothing is dated:\n%s", out)
	}
	if !strings.Contains(readQ(t, filepath.Join(off, "b.ts")), "updated_at: 2026-09-01") {
		t.Error("with touch.pre_commit false the staged file keeps its date")
	}
}

// The child's exit 3 ("not governed") must cross the process boundary.
func TestVerifyTranslatesTheChildsNotGoverned(t *testing.T) {
	t.Run("VPFVR-E01: The child's not-governed exit stays not-governed", func(t *testing.T) {})
	fakeChild(t, 3)
	_, err := runQ(t, newVerifyCmd(), "--root", t.TempDir(), "--changed", "package.json")
	var ng errNotGoverned
	if !errors.As(err, &ng) {
		t.Errorf("exit 3 of the child must become errNotGoverned, got %v (%T)", err, err)
	}
}

func TestVerifyGuards(t *testing.T) {
	t.Run("VPFVR-B02: Nothing staged has nothing to verify", func(t *testing.T) {})
	t.Run("VPFVR-E03: Verify with no scope is refused", func(t *testing.T) {})
	englishOutput(t)
	root := touchRepo(t)

	out, err := runQ(t, newVerifyCmd(), "--root", root, "--staged")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "nothing staged — nothing to verify.") {
		t.Errorf("an empty index has nothing to verify:\n%s", out)
	}

	if _, err := runQ(t, newVerifyCmd(), "--root", root); err == nil || !strings.Contains(err.Error(), "specify --staged, --changed <file> or --all") {
		t.Errorf("no scope at all must be refused; got %v", err)
	}
}

func TestStagedFiles(t *testing.T) {
	t.Run("VPFVR-B01: The staged scope is the added, copied, modified and renamed files of the index", func(t *testing.T) {})
	t.Run("VPFVR-E04: The staged scope outside a git repository is refused", func(t *testing.T) {})
	root := touchRepo(t)
	touchWrite(t, root, "a.ts", touchHeader+"export const x = 3\n")
	touchWrite(t, root, "new.ts", "export const n = 1\n")
	touchWrite(t, root, "untracked.ts", "export const u = 1\n")
	gitIn(t, root, "add", "a.ts", "new.ts")
	gitIn(t, root, "rm", "-q", "b.ts")

	got, err := stagedFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	// modified and added — never the deletion (nothing to verify) nor the untracked file
	if strings.Join(got, ",") != "a.ts,new.ts" {
		t.Errorf("stagedFiles = %v, want [a.ts new.ts]", got)
	}

	dir := t.TempDir()
	if temRepoAcima(dir) {
		t.Skip("the temporary directory is inside a git repository")
	}
	if _, err := stagedFiles(dir); err == nil {
		t.Error("outside a repository there is no index: must fail")
	}
}

// A touch failure warns and does not block: the updated-at gate still checks the files.
func TestPrintPreCommitTouchWarnsOnFailure(t *testing.T) {
	t.Run("VPFVR-E05: A dating failure warns and does not stop the verify", func(t *testing.T) {})
	out := captureStdout(t, func() { printPreCommitTouch(nil, nil, errors.New("boom")) })
	if !strings.Contains(out, "· touch: could not date the staged files (boom)") {
		t.Errorf("a touch failure must be said:\n%s", out)
	}
}

// The child's exit code crosses the process boundary.
//
// `verify` is a facade that re-executes its own binary. `c.Run()` returns a generic
// `*exec.ExitError`, and `main` — which turns `errNotGoverned` into `ExitNotGoverned` —
// did not recognise it. `check` exited 3 ("no jurisdiction"), `verify` translated it to 1,
// and the pre-commit barred a configuration-only commit: exactly the case exit 3 exists
// to let through.
func TestExitNotGovernedCrossesTheSubprocess(t *testing.T) {
	// `sh -c 'exit 3'` reproduces the child's code.
	err := translateChildOutput(exec.Command("sh", "-c", "exit 3").Run())

	var nr errNotGoverned
	if !errors.As(err, &nr) {
		t.Fatalf("exit 3 of the child did not become errNotGoverned: %v (%T)", err, err)
	}
}

// Any other code stays a failure: only 3 has its own handling.
func TestOtherCodesStayFailures(t *testing.T) {
	t.Run("VPFVR-E02: Any other failing exit of the child stays a failure", func(t *testing.T) {})
	err := translateChildOutput(exec.Command("sh", "-c", "exit 1").Run())

	var nr errNotGoverned
	if errors.As(err, &nr) {
		t.Errorf("exit 1 was treated as ungoverned — a failure would be swallowed")
	}
	if err == nil {
		t.Error("exit 1 returned no error")
	}
}

func TestSuccessReturnsNoError(t *testing.T) {
	if err := translateChildOutput(exec.Command("sh", "-c", "exit 0").Run()); err != nil {
		t.Errorf("exit 0 returned an error: %v", err)
	}
}

// THE HOOK MUST NOT DUMP THE WHOLE REPORT — let alone twice.
//
// Measured in the reference app: `pre-commit` and `commit-msg` call the SAME
// `anchors verify --phase pre-commit --staged` (on purpose: one reports early, the other
// bars, and only the second has the message in hand). Without `--only-issues` the table of
// 42 gates came out in both, ~88 lines per commit, and the push's `remote:` — the one line
// the person was waiting for — scrolled off the screen.
func TestAutomaticPhaseDoesNotDumpTheWholeReport(t *testing.T) {
	t.Run("VPFVR-B04: An automatic phase asks check for computable gates and issues only", func(t *testing.T) {})
	for _, phase := range []string{"pre-commit", "pre-push", "ci"} {
		args := checkArgs(".", phase, "", "", []string{"a.spec.md"}, false, false, true)
		if !slices.Contains(args, "--only-issues") || !slices.Contains(args, "--deterministic") {
			t.Errorf("phase %q calls check WITHOUT --deterministic --only-issues: the hook dumps the whole gate table\n  args: %v", phase, args)
		}
	}
}

// In the MANUAL phase the report is the product: the person typed the command to see the
// picture, and omitting the clean gates would hide what they asked for. With no phase at
// all, the same.
func TestManualPhaseShowsTheFullReport(t *testing.T) {
	t.Run("VPFVR-B05: The manual phase, or no phase, asks check for the full report", func(t *testing.T) {})
	for _, phase := range []string{"manual", ""} {
		args := checkArgs(".", phase, "", "", []string{"a.spec.md"}, false, false, true)
		if slices.Contains(args, "--only-issues") {
			t.Errorf("phase %q omitted the clean gates — the person asked for the picture\n  args: %v", phase, args)
		}
		if slices.Contains(args, "--deterministic") {
			t.Errorf("phase %q became deterministic — there AI judgment is the point", phase)
		}
	}
}

// What already worked keeps working: extracting `checkArgs` from the RunE must not have
// lost a flag on the way.
func TestCheckArgsKeepsTheFacadeFlags(t *testing.T) {
	t.Run("VPFVR-B09: The facade's flags reach check unchanged", func(t *testing.T) {})
	t.Run("VPFVR-I01: Verify never asks check for both the full sweep and a file list", func(t *testing.T) {})
	args := checkArgs("/r", "pre-commit", "/tmp/MSG", "types", []string{"a.md", "b.md"}, false, true, true)

	for _, want := range []string{
		"check", "--root", "/r",
		"--changed", "a.md", "--changed", "b.md",
		"--phase", "pre-commit",
		"--commit-msg", "/tmp/MSG",
		"--category", "types",
		"--skip-slow", "--no-record",
	} {
		if !slices.Contains(args, want) {
			t.Errorf("flag %q was lost in the extraction\n  args: %v", want, args)
		}
	}

	// `--all` and `--changed` are exclusive: asking for both would make check sweep the
	// whole project in a hook meant to be incremental.
	if slices.Contains(args, "--all") {
		t.Errorf("--changed came together with --all\n  args: %v", args)
	}
	all := checkArgs("/r", "", "", "", []string{"a.md"}, true, false, false)
	if !slices.Contains(all, "--all") || slices.Contains(all, "--changed") {
		t.Errorf("--all must replace the file list\n  args: %v", all)
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

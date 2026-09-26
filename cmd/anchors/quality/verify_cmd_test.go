package quality

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
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

// The child's exit 3 ("not governed") must cross the process boundary.
func TestVerifyTranslatesTheChildsNotGoverned(t *testing.T) {
	fakeChild(t, 3)
	_, err := runQ(t, newVerifyCmd(), "--root", t.TempDir(), "--changed", "package.json")
	var ng errNotGoverned
	if !errors.As(err, &ng) {
		t.Errorf("exit 3 of the child must become errNotGoverned, got %v (%T)", err, err)
	}
}

func TestVerifyGuards(t *testing.T) {
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
	out := captureStdout(t, func() { printPreCommitTouch(nil, nil, errors.New("boom")) })
	if !strings.Contains(out, "· touch: could not date the staged files (boom)") {
		t.Errorf("a touch failure must be said:\n%s", out)
	}
}

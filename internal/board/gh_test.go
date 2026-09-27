package board

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// fakeGhOnPath puts on the PATH a `gh` running the given shell body.
func fakeGhOnPath(t *testing.T, body string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// A `gh` that decides to prompt must find its input closed and give up, not wait forever.
func TestRunGH_givesGhNoInput(t *testing.T) {
	t.Run("GHRNG-B01: gh gets no input", func(t *testing.T) {})
	fakeGhOnPath(t, `if read line; then echo "got: $line"; else exit 7; fi`)
	_, err := runGH("auth", "login")
	if err == nil || !strings.Contains(err.Error(), "exit status 7") {
		t.Fatalf("runGH = %v, want the end of input (exit 7)", err)
	}
}

// `gh` exit status 4 is NOT AUTHENTICATED, and alone it says nothing. Measured with a new dev:
// `anchors next` answered `exit status 4` and they had no way to know login was missing.
func TestAuthHint_namesTheCauseOfCode4(t *testing.T) {
	t.Run("GHRNG-B02: Exit code 4 is explained as not authenticated", func(t *testing.T) {})
	fakeGhOnPath(t, "exit 4")
	_, err := runGH("issue", "list")
	if err == nil {
		t.Fatal("the command should fail with 4")
	}
	for _, want := range []string{"exit status 4", "NOT AUTHENTICATED", "gh auth login", "interactive"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not mention %q:\n%s", want, err)
		}
	}
}

// ANOTHER exit code does NOT get the hint: telling "login missing" to someone with a wrong
// repository in `workflow.repo` sends them to a login that is already done.
func TestAuthHint_onlyCode4(t *testing.T) {
	t.Run("GHRNG-B03: Other failures get no authentication hint", func(t *testing.T) {})
	for _, code := range []string{"1", "2", "3", "5", "127"} {
		err := exec.Command("sh", "-c", "exit "+code).Run()
		if hint := authHint(err, nil); hint != "" {
			t.Errorf("exit %s got the auth hint:\n%s", code, hint)
		}
	}
	// And an error that is not an exit (a command that does not exist) does not either.
	if hint := authHint(exec.Command("command-that-does-not-exist").Run(), nil); hint != "" {
		t.Errorf("an error that is not an ExitError got the hint:\n%s", hint)
	}
}

// The CEILING exists so `gh` cannot hang Anchors; generous on purpose, since a tight ceiling
// would turn a slow network into an error. (This checks the constant only: the timeout path
// itself takes the full ceiling to run.)
func TestGhTimeout_existsAndIsGenerous(t *testing.T) {
	if ghTimeout == 0 {
		t.Fatal("without a ceiling, a `gh` waiting for input hangs Anchors forever")
	}
	if ghTimeout.Seconds() < 10 {
		t.Errorf("a ceiling of %s is tight — a slow network would become an error", ghTimeout)
	}
	if ghTimeout.Minutes() > 2 {
		t.Errorf("a ceiling of %s is long — the point is not to cost the session", ghTimeout)
	}
}

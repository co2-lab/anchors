// @anchors
//   ref: PSXSH

package shell

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// fake replaces what the lookup consults, for one test.
func fake(t *testing.T, onPath string, execPath string, files ...string) {
	t.Helper()
	pl, ge, sf := lookPath, gitExecPath, statFile
	t.Cleanup(func() { lookPath, gitExecPath, statFile = pl, ge, sf })
	lookPath = func(string) (string, error) {
		if onPath == "" {
			return "", exec.ErrNotFound
		}
		return onPath, nil
	}
	gitExecPath = func() (string, error) {
		if execPath == "" {
			return "", errors.New("no git")
		}
		return execPath, nil
	}
	have := map[string]bool{}
	for _, f := range files {
		have[f] = true
	}
	statFile = func(p string) bool { return have[p] }
}

func TestPath_onPath(t *testing.T) {
	t.Run("PSXSH-B01: The sh on PATH is the shell", func(t *testing.T) {})
	for _, goos := range []string{"linux", "darwin", "windows"} {
		fake(t, "/bin/sh", "")
		if got, err := resolve(goos); err != nil || got != "/bin/sh" {
			t.Errorf("%s: want /bin/sh, got %q %v", goos, got, err)
		}
	}
}

func TestPath_besideGit(t *testing.T) {
	t.Run("PSXSH-B02: On Windows, the shell beside git", func(t *testing.T) {})
	git := filepath.Join("C:", "Git")
	execPath := filepath.Join(git, "mingw64", "libexec", "git-core")
	usr := filepath.Join(git, "usr", "bin", "sh.exe")
	bin := filepath.Join(git, "bin", "sh.exe")
	fake(t, "", execPath, usr)
	if got, err := resolve("windows"); err != nil || got != usr {
		t.Errorf("want %s, got %q %v", usr, got, err)
	}
	fake(t, "", execPath, usr, bin)
	if got, _ := resolve("windows"); got != bin {
		t.Errorf("bin\\sh.exe comes first, got %q", got)
	}
	fake(t, "", execPath, bin)
	if _, err := resolve("linux"); !errors.Is(err, ErrNoShell) {
		t.Errorf("only Windows looks beside git, got %v", err)
	}
}

func TestCommand_runsTheScript(t *testing.T) {
	t.Run("PSXSH-B03: A command is sh -c with its arguments", func(t *testing.T) {})
	if runtime.GOOS == "windows" {
		if _, err := Path(); err != nil {
			t.Skip("no POSIX shell on this Windows machine")
		}
	}
	cmd, err := Command(`printf '%s|%s' "$0" "$1"`, "zero", "one")
	if err != nil {
		t.Fatal(err)
	}
	if out, err := cmd.Output(); err != nil || string(out) != "zero|one" {
		t.Errorf("the arguments reach the script, got %q %v", out, err)
	}
	if cmd.Args[1] != "-c" {
		t.Errorf("the shell runs the script with -c, got %v", cmd.Args)
	}
}

func TestPath_noShell(t *testing.T) {
	t.Run("PSXSH-E01: No shell is an environment error", func(t *testing.T) {})
	fake(t, "", "")
	for _, goos := range []string{"linux", "windows"} {
		_, err := resolve(goos)
		if !errors.Is(err, ErrNoShell) || !strings.Contains(err.Error(), "POSIX shell") {
			t.Errorf("%s: want the environment error, got %v", goos, err)
		}
	}
	if _, err := resolve("windows"); !strings.Contains(err.Error(), "Git for Windows") {
		t.Errorf("on Windows the error says how to get one, got %v", err)
	}
	fake(t, "", filepath.Join("C:", "Git", "mingw64", "libexec", "git-core"))
	if _, err := resolve("windows"); !errors.Is(err, ErrNoShell) {
		t.Errorf("git with no shell beside it is still no shell, got %v", err)
	}
	fake(t, "", "")
	if _, err := Command("true"); !errors.Is(err, ErrNoShell) {
		t.Errorf("a command with no shell is the environment error, got %v", err)
	}
}

// On a real Windows, with only git on PATH, the shell is found beside it: the default of
// Git for Windows, where `sh.exe` is not on PATH.
func TestPath_realGitOnWindows(t *testing.T) {
	t.Run("PSXSH-B02: On Windows, the shell beside git", func(t *testing.T) {})
	if runtime.GOOS != "windows" {
		t.Skip("proves Git for Windows' layout")
	}
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("no git")
	}
	t.Setenv("PATH", filepath.Dir(git))
	sh, err := Path()
	if err != nil || !strings.HasSuffix(strings.ToLower(sh), "sh.exe") {
		t.Fatalf("the shell beside git, got %q %v", sh, err)
	}
	cmd, err := Command("printf ok")
	if err != nil {
		t.Fatal(err)
	}
	if out, err := cmd.Output(); err != nil || string(out) != "ok" {
		t.Errorf("git's shell runs a script, got %q %v", out, err)
	}
}

// With nothing on PATH, on any system, the lookup says so as an environment error.
func TestPath_realEmptyPath(t *testing.T) {
	t.Run("PSXSH-E01: No shell is an environment error", func(t *testing.T) {})
	t.Setenv("PATH", t.TempDir())
	if _, err := Path(); !errors.Is(err, ErrNoShell) {
		t.Errorf("an empty PATH has no shell, got %v", err)
	}
}

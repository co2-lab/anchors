// @anchors
//   ref: PSXSH

// Package shell finds the POSIX shell that runs the commands a project declares — a gate's
// `run:`, a suite's `run:`, a tests `script:`. They are written in POSIX shell (`"$@"`,
// `&&`, `VAR=x cmd`), so they need `sh` on every system.
//
// On Linux and macOS `sh` is always there. On Windows it comes with Git for Windows, but its
// default install puts only `git.exe` on PATH, not the `sh.exe` beside it: the commands then
// failed with a bare "executable file not found", read as the gate failing. Here the shell
// is looked up on PATH first and then next to the git that is installed, and when neither
// has one the error says what is missing and how to get it.
package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrNoShell is the environment error of a system with no POSIX shell to run the project's
// commands.
var ErrNoShell = errors.New("no POSIX shell (`sh`) found to run the project's commands")

// What the lookup consults, replaceable in tests.
var (
	lookPath    = exec.LookPath
	gitExecPath = func() (string, error) {
		out, err := exec.Command("git", "--exec-path").Output()
		return strings.TrimSpace(string(out)), err
	}
	statFile = func(p string) bool {
		fi, err := os.Stat(p)
		return err == nil && !fi.IsDir()
	}
)

// Path is the shell to run a project command with.
func Path() (string, error) { return resolve(runtime.GOOS) }

func resolve(goos string) (string, error) {
	if p, err := lookPath("sh"); err == nil {
		return p, nil
	}
	if goos != "windows" {
		return "", fmt.Errorf("%w: install a POSIX shell and put `sh` on PATH", ErrNoShell)
	}
	// Git for Windows: `git --exec-path` is `<git>\mingw64\libexec\git-core`, and the shell
	// is `<git>\bin\sh.exe` or `<git>\usr\bin\sh.exe`.
	if exe, err := gitExecPath(); err == nil && exe != "" {
		gitRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Clean(exe))))
		for _, c := range []string{filepath.Join(gitRoot, "bin", "sh.exe"), filepath.Join(gitRoot, "usr", "bin", "sh.exe")} {
			if statFile(c) {
				return c, nil
			}
		}
	}
	return "", fmt.Errorf("%w: install Git for Windows (it brings `sh.exe`), or put its `usr\\bin` on PATH", ErrNoShell)
}

// Command is `sh -c script`, with `args` after it as `$0 $1…`, through the shell Path finds.
func Command(script string, args ...string) (*exec.Cmd, error) {
	sh, err := Path()
	if err != nil {
		return nil, err
	}
	return exec.Command(sh, append([]string{"-c", script}, args...)...), nil //nolint:gosec // the command is declared by the project
}

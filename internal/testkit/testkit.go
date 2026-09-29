// Package testkit holds what the tests of every package share to run the same on Linux,
// macOS and Windows. It is imported only by tests.
package testkit

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/co2-lab/anchors/internal/shell"
)

// CaptureStdout returns what fn printed on standard output.
//
// The pipe is drained WHILE fn runs. Reading it only afterwards deadlocked on Windows,
// whose pipe buffer is a few kilobytes: a command printing a long prompt blocked on the
// write, and the test sat there until the ten-minute timeout.
func CaptureStdout(t *testing.T, fn func()) string {
	t.Helper()
	return capture(t, &os.Stdout, fn)
}

// CaptureStderr returns what fn printed on standard error.
func CaptureStderr(t *testing.T, fn func()) string {
	t.Helper()
	return capture(t, &os.Stderr, fn)
}

func capture(t *testing.T, std **os.File, fn func()) string {
	t.Helper()
	orig := *std
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	*std = w
	defer func() { *std = orig }()
	fn()
	w.Close()
	*std = orig
	return <-done
}

// FakeBin writes an executable named `name` into dir that runs the POSIX shell `script`,
// and returns its path.
//
// On Windows a script with no extension does not run, and `exec.LookPath` would find the
// real tool first. A `.cmd` wrapper was tried and lost arguments: `cmd.exe` reads `&`, `<`,
// `>` and `|` inside them — a GitHub URL with `?a=1&b=2` became two commands. So on Windows
// the fake is `name.exe`, a small launcher built once per test run, which hands its exact
// argv to the shell Git brings, running the script beside it.
func FakeBin(t *testing.T, dir, name, script string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	body := script
	if !strings.HasPrefix(body, "#!") {
		body = "#!/bin/sh\n" + body
	}
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		sh, err := shell.Path()
		if err != nil {
			t.Skipf("no POSIX shell to run the fake %s: %v", name, err)
		}
		exe, err := launcher()
		if err != nil {
			t.Fatalf("build the fake launcher: %v", err)
		}
		b, err := os.ReadFile(exe)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p+".exe", b, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("ANCHORS_TESTKIT_SH", sh)
	}
	return p
}

// launcherSource runs the script that sits beside it, with the same name and no `.exe`,
// through the shell named in ANCHORS_TESTKIT_SH, passing every argument as it came.
const launcherSource = `package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	self, err := os.Executable()
	if err != nil {
		os.Exit(127)
	}
	script := strings.TrimSuffix(self, ".exe")
	cmd := exec.Command(os.Getenv("ANCHORS_TESTKIT_SH"), append([]string{script}, os.Args[1:]...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(127)
	}
}
`

var (
	launcherOnce sync.Once
	launcherPath string
	launcherErr  error
)

// launcher builds the fake launcher once per test binary.
func launcher() (string, error) {
	launcherOnce.Do(func() {
		dir, err := os.MkdirTemp("", "anchors-testkit-")
		if err != nil {
			launcherErr = err
			return
		}
		src := filepath.Join(dir, "main.go")
		if err := os.WriteFile(src, []byte(launcherSource), 0o644); err != nil {
			launcherErr = err
			return
		}
		out := filepath.Join(dir, "launcher.exe")
		cmd := exec.Command("go", "build", "-o", out, src)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GO111MODULE=off")
		if b, err := cmd.CombinedOutput(); err != nil {
			launcherErr = fmt.Errorf("%v: %s", err, b)
			return
		}
		launcherPath = out
	})
	return launcherPath, launcherErr
}

// OnPath puts dir first on PATH for the rest of the test.
func OnPath(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// SkipWithoutPOSIXPermissions skips a test that needs a file the process cannot read or
// write, or a file mode kept across a rewrite: Windows does not enforce POSIX permissions
// (`chmod` only toggles read-only), so there is nothing there to observe.
func SkipWithoutPOSIXPermissions(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not enforce POSIX file permissions")
	}
}

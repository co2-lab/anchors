// Package testkit holds what the tests of every package share to run the same on Linux,
// macOS and Windows. It is imported only by tests.
package testkit

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
// and returns its path. On Windows a script with no extension does not run, and
// `exec.LookPath` would find the real tool first; a `name.cmd` beside it runs the script
// through the shell Git for Windows brings, so the same fake serves every system.
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
		wrapper := "@\"" + sh + "\" \"%~dp0" + name + "\" %*\r\n"
		if err := os.WriteFile(p+".cmd", []byte(wrapper), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

// OnPath puts dir first on PATH for the rest of the test.
func OnPath(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

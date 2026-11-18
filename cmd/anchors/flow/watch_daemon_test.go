// @anchors
//   code: WDTWT
//   ref: WTDMW

package flow

import (
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// This file builds on every platform, because it is the one test file linked to the
// WatchDaemon rules: a `_windows_test.go` or `//go:build` split would leave the half it
// excludes with no link. Each platform's half runs on its platform and is skipped on the
// other, and it reads nothing that exists on one platform only: the process group comes
// from `ps`, and the Windows creation flags are read by field name.

// The watcher started by `watch start` must survive the terminal: on unix it leads its own
// session, so it is the leader of a process group that is not the terminal's.
func TestDetach_childLeadsItsOwnSession(t *testing.T) {
	t.Run("WTDMW-B01: On unix a detached child leads its own process group", func(t *testing.T) {})
	if runtime.GOOS == "windows" {
		t.Skip("unix-only: Windows has no sessions (WTDMW-B02 is its half)")
	}
	c := exec.Command("sleep", "30")
	detach(c)
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Process.Kill(); _, _ = c.Process.Wait() }()
	pgid := processGroupOf(t, c.Process.Pid)
	if pgid != c.Process.Pid {
		t.Errorf("a detached child leads its own group: pgid %d, pid %d", pgid, c.Process.Pid)
	}
	if mine := processGroupOf(t, os.Getpid()); pgid == mine {
		t.Error("the child must not stay in the parent's process group")
	}
}

// Detaching only prepares the child: the caller starts it and records its pid.
func TestDetach_doesNotStartTheChild(t *testing.T) {
	t.Run("WTDMW-X01: Detaching does not start the child", func(t *testing.T) {})
	c := exec.Command("sleep", "30")
	detach(c)
	if c.Process != nil {
		_ = c.Process.Kill()
		t.Fatal("detach started the child")
	}
	if c.SysProcAttr == nil {
		t.Fatal("detach left no attribute for the start")
	}
}

// Exactly one of the two platform files is built for a target: none would leave `detach`
// undefined, both would define it twice.
func TestDetach_oneImplementationPerPlatform(t *testing.T) {
	t.Run("WTDMW-I01: Exactly one detachment implementation builds per platform", func(t *testing.T) {})
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not on PATH")
	}
	files := func(goos string) string {
		c := exec.Command("go", "list", "-f", "{{.GoFiles}}", ".")
		c.Env = append(os.Environ(), "GOOS="+goos, "GOARCH=amd64")
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("go list for %s: %v\n%s", goos, err, out)
		}
		return string(out)
	}
	for goos, want := range map[string]string{"linux": "watch_unix.go", "darwin": "watch_unix.go", "windows": "watch_windows.go"} {
		got := files(goos)
		other := "watch_windows.go"
		if want == other {
			other = "watch_unix.go"
		}
		if !strings.Contains(got, want) || strings.Contains(got, other) {
			t.Errorf("%s must build %s and not %s: %s", goos, want, other, got)
		}
	}
}

// processGroupOf reads a process's group id through `ps`, which answers the same on every
// unix-like system; `syscall.Getpgid` would keep this file from building on Windows.
func processGroupOf(t *testing.T, pid int) int {
	t.Helper()
	if _, err := exec.LookPath("ps"); err != nil {
		t.Skip("ps not on PATH")
	}
	out, err := exec.Command("ps", "-o", "pgid=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		t.Fatalf("ps for pid %d: %v", pid, err)
	}
	pgid, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatalf("ps gave no process group for pid %d: %q", pid, out)
	}
	return pgid
}

// On Windows there is no session: the child is created in a new process group, so the
// console's CTRL_C_EVENT sent to the parent's group does not reach it.
//
// The flag is read by field name because `SysProcAttr.CreationFlags` and
// `syscall.CREATE_NEW_PROCESS_GROUP` exist only in the Windows build of `syscall`.
func TestDetach_childGetsANewProcessGroup(t *testing.T) {
	t.Run("WTDMW-B02: On Windows a detached child is created in a new process group", func(t *testing.T) {})
	if runtime.GOOS != "windows" {
		t.Skip("Windows-only: unix-like systems detach with a new session (WTDMW-B01)")
	}
	const createNewProcessGroup = 0x00000200 // syscall.CREATE_NEW_PROCESS_GROUP
	c := exec.Command("cmd", "/c", "exit")
	detach(c)
	if c.SysProcAttr == nil {
		t.Fatal("detach left no attribute for the start")
	}
	flags := reflect.ValueOf(c.SysProcAttr).Elem().FieldByName("CreationFlags")
	if !flags.IsValid() || flags.Uint()&createNewProcessGroup == 0 {
		t.Errorf("the child must be created with CREATE_NEW_PROCESS_GROUP: %+v", c.SysProcAttr)
	}
}

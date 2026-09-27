//go:build !windows

package flow

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

// The watcher started by `watch start` must survive the terminal: on unix it leads its own
// session, so it is the leader of a process group that is not the terminal's.
func TestDetach_childLeadsItsOwnSession(t *testing.T) {
	t.Run("WTDMW-B01: On unix a detached child leads its own process group", func(t *testing.T) {})
	c := exec.Command("sleep", "30")
	detach(c)
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Process.Kill(); _, _ = c.Process.Wait() }()
	pgid, err := syscall.Getpgid(c.Process.Pid)
	if err != nil {
		t.Fatal(err)
	}
	if pgid != c.Process.Pid {
		t.Errorf("a detached child leads its own group: pgid %d, pid %d", pgid, c.Process.Pid)
	}
	if mine, _ := syscall.Getpgid(os.Getpid()); pgid == mine {
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

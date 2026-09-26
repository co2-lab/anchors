package daemon

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPathsFor(t *testing.T) {
	p := PathsFor("/proj")
	want := Paths{
		Dir:    filepath.Join("/proj", ".anchors"),
		PID:    filepath.Join("/proj", ".anchors", "watch.pid"),
		Log:    filepath.Join("/proj", ".anchors", "watch.log"),
		Paused: filepath.Join("/proj", ".anchors", "watch.paused"),
		Meta:   filepath.Join("/proj", ".anchors", "watch.meta"),
	}
	if p != want {
		t.Fatalf("PathsFor = %+v, want %+v", p, want)
	}
}

func TestRunning(t *testing.T) {
	p := PathsFor(t.TempDir())

	if pid := Running(p); pid != 0 {
		t.Fatalf("no PID file: Running = %d, want 0", pid)
	}

	if err := WritePID(p, os.Getpid()); err != nil {
		t.Fatal(err)
	}
	if pid := Running(p); pid != os.Getpid() {
		t.Fatalf("live PID: Running = %d, want %d", pid, os.Getpid())
	}

	os.WriteFile(p.PID, []byte("garbage"), 0o644)
	if pid := Running(p); pid != 0 {
		t.Fatalf("unparseable PID: Running = %d, want 0", pid)
	}

	// A PID that no longer exists: Running answers 0 and removes the stale file.
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Skip("cannot spawn a process:", err)
	}
	WritePID(p, cmd.Process.Pid)
	if pid := Running(p); pid != 0 {
		t.Fatalf("dead PID: Running = %d, want 0", pid)
	}
	if _, err := os.Stat(p.PID); !os.IsNotExist(err) {
		t.Fatalf("a stale PID file must be removed: %v", err)
	}
}

func TestMeta(t *testing.T) {
	p := PathsFor(t.TempDir())
	if got := ReadMeta(p); got != "" {
		t.Fatalf("missing meta: ReadMeta = %q, want empty", got)
	}
	os.MkdirAll(p.Dir, 0o755)
	started := time.Date(2026, 9, 26, 10, 30, 0, 0, time.UTC)
	if err := WriteMeta(p, started, "/proj"); err != nil {
		t.Fatal(err)
	}
	if got, want := ReadMeta(p), "started=2026-09-26T10:30:00Z\nroot=/proj\n"; got != want {
		t.Fatalf("ReadMeta = %q, want %q", got, want)
	}
}

func TestPauseResumeCleanup(t *testing.T) {
	p := PathsFor(t.TempDir())
	os.MkdirAll(p.Dir, 0o755)
	if IsPaused(p) {
		t.Fatal("a fresh project must not be paused")
	}
	if err := Pause(p); err != nil {
		t.Fatal(err)
	}
	if !IsPaused(p) {
		t.Fatal("Pause must leave the watcher paused")
	}
	if err := Resume(p); err != nil {
		t.Fatal(err)
	}
	if IsPaused(p) {
		t.Fatal("Resume must clear the pause")
	}

	Pause(p)
	WritePID(p, 12345)
	Cleanup(p)
	for _, f := range []string{p.PID, p.Paused} {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			t.Errorf("Cleanup must remove %s: %v", f, err)
		}
	}
}

func TestStop(t *testing.T) {
	p := PathsFor(t.TempDir())
	if err := Stop(p); err == nil || !strings.Contains(err.Error(), "not running") {
		t.Fatalf("Stop without a daemon = %v, want 'not running'", err)
	}

	// A real child stands in for the watcher; Stop must terminate it and drop the PID.
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Skip("cannot spawn sleep:", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	if err := WritePID(p, cmd.Process.Pid); err != nil {
		t.Fatal(err)
	}
	if err := Stop(p); err != nil {
		t.Fatalf("Stop = %v", err)
	}
	select {
	case err := <-done:
		var ee *exec.ExitError
		if !errors.As(err, &ee) || ee.ExitCode() != -1 {
			t.Fatalf("the child must end by a signal, got %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Stop did not terminate the process")
	}
	if _, err := os.Stat(p.PID); !os.IsNotExist(err) {
		t.Fatalf("Stop must remove the PID file: %v", err)
	}
}

// @anchors
//   code: DRTDM
//   ref: DMRND

package daemon

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func TestAlive_liveAndExitedProcesses(t *testing.T) {
	t.Run("DMRND-B01: A live process is alive and an exited one is not", func(t *testing.T) {})
	if !alive(os.Getpid()) {
		t.Fatal("this test process must be reported alive")
	}
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Skip("cannot spawn a process:", err)
	}
	if alive(cmd.Process.Pid) {
		t.Fatalf("an exited process (pid %d) must not be reported alive", cmd.Process.Pid)
	}
}

func TestTerminate_sendsACatchableSIGTERM(t *testing.T) {
	t.Run("DMRND-B02: Termination on Unix-like systems is a catchable SIGTERM", func(t *testing.T) {})
	if runtime.GOOS == "windows" {
		t.Skip("Windows has no SIGTERM: the process is killed")
	}
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Skip("cannot spawn sleep:", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	if err := terminate(cmd.Process); err != nil {
		t.Fatalf("terminate = %v", err)
	}
	select {
	case err := <-done:
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("the child must end by a signal, got %v", err)
		}
		ws, ok := ee.Sys().(syscall.WaitStatus)
		if !ok || !ws.Signaled() || ws.Signal() != syscall.SIGTERM {
			t.Fatalf("the child must end by SIGTERM, got %v", ee)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("terminate did not end the process")
	}
}

func TestStop_removesThePIDFileWithoutTheLoop(t *testing.T) {
	t.Run("DMRND-I01: Stopping leaves no PID file even when the process cleans nothing", func(t *testing.T) {})
	p := PathsFor(t.TempDir())
	// `sleep` is not the watcher: it never runs the loop's cleanup.
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Skip("cannot spawn sleep:", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	go func() { _ = cmd.Wait() }()

	if err := WritePID(p, cmd.Process.Pid); err != nil {
		t.Fatal(err)
	}
	if err := Stop(p); err != nil {
		t.Fatalf("Stop = %v", err)
	}
	if _, err := os.Stat(p.PID); !os.IsNotExist(err) {
		t.Fatalf("the PID file must be gone after Stop: %v", err)
	}
}

func TestAlive_aProcessWeMayNotSignalIsAlive(t *testing.T) {
	t.Run("DMRND-B03: A live process of another user is alive", func(t *testing.T) {})
	if runtime.GOOS == "windows" {
		t.Skip("the permission probe is the Unix one")
	}
	if os.Geteuid() == 0 {
		t.Skip("root may signal every process: no EPERM to observe")
	}
	// PID 1 (init/launchd) is alive and owned by root, so signal 0 answers EPERM.
	pid1, err := os.FindProcess(1)
	if err != nil {
		t.Skip("cannot find pid 1:", err)
	}
	if err := pid1.Signal(syscall.Signal(0)); !errors.Is(err, syscall.EPERM) {
		t.Skipf("signal 0 to pid 1 = %v, not EPERM", err)
	}
	if !alive(1) {
		t.Fatal("pid 1 answers EPERM, so it exists: it must be reported alive")
	}
}

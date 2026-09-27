//go:build !windows

package daemon

import (
	"errors"
	"os"
	"syscall"
)

// alive: o processo existe? (signal 0 não envia nada, só testa.)
//
// EPERM means the process EXISTS and belongs to someone we may not signal. It used to count
// as dead: a live watcher of another user was reported stale and its PID file deleted,
// which lets a second watcher start beside it.
func alive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}

// terminate envia SIGTERM — o loop do watcher trata o sinal para sair limpo
// (ver watch.go, signal.Notify).
func terminate(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}

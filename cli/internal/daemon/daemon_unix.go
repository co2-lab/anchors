//go:build !windows

package daemon

import (
	"os"
	"syscall"
)

// alive: o processo existe? (signal 0 não envia nada, só testa.)
func alive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// terminate envia SIGTERM — o loop do watcher trata o sinal para sair limpo
// (ver watch.go, signal.Notify).
func terminate(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}

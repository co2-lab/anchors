// @anchors
//   code: DWIDM
//   ref: DMRND

//go:build windows

package daemon

import (
	"os"
	"syscall"
)

// stillActive is the exit code Windows reports for a process that has not exited.
const stillActive = 259

// alive says whether the process is running. A handle alone does not say it: Windows keeps
// an exited process's object while any handle to it is open, so opening it succeeded for a
// process that had already exited — the parent's own handle was enough. The exit code
// does say it: STILL_ACTIVE until the process ends.
func alive(pid int) bool {
	h, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(h)
	var code uint32
	if err := syscall.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code == stillActive
}

// terminate mata o processo diretamente — Windows não tem SIGTERM; Process.Kill
// chama TerminateProcess. O watcher não recebe chance de encerrar o loop de
// forma graciosa nesta plataforma (ver watch.go).
func terminate(proc *os.Process) error {
	return proc.Kill()
}

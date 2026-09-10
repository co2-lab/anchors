//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// detach desanexa o processo filho da sessão do terminal atual (nova sessão
// via Setsid), para sobreviver ao encerramento do terminal.
func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

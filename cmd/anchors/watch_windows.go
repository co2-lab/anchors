//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// detach desanexa o processo filho do console atual — Windows não tem
// sessões/Setsid; CREATE_NEW_PROCESS_GROUP evita que o filho receba o
// CTRL_C_EVENT do console pai, o mais próximo do comportamento de
// "sobrevive ao terminal" disponível de forma portável.
func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

//go:build !windows

package quality

import (
	"os/exec"
	"syscall"
)

// ownProcessGroup starts the command in a group of its own, so the deadline can stop the
// runner's workers along with the shell.
func ownProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup stops the command's whole group.
func killProcessGroup(cmd *exec.Cmd) error {
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

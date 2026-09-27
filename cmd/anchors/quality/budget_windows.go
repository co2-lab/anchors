//go:build windows

package quality

import "os/exec"

// ownProcessGroup: Windows has no process group to kill as one; the command runs as is.
func ownProcessGroup(*exec.Cmd) {}

// killProcessGroup stops the command's process.
func killProcessGroup(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

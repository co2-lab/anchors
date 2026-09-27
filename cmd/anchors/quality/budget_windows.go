//go:build windows

package quality

import (
	"os/exec"
	"time"
)

// ownProcessGroup: Windows has no process group to signal as one; the command runs as is.
func ownProcessGroup(*exec.Cmd) {}

// stopProcessGroup stops the command's process. Windows has no TERM a tool could trap to
// undo what it changed, so a mutation tool working in place can be cut mid-mutant here.
func stopProcessGroup(cmd *exec.Cmd, _ time.Duration) error {
	return cmd.Process.Kill()
}

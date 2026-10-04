// @anchors
//   code: BWCBD
//   ref: BDGRN

//go:build windows

package quality

import (
	"os/exec"
	"strconv"
	"time"
)

// ownProcessGroup: Windows has no process group to signal as one; the command runs as is.
func ownProcessGroup(*exec.Cmd) {}

// stopProcessGroup stops the command's process and every process it started. Windows has
// no group to signal, but `taskkill /T` — part of Windows — ends the whole tree; killing
// only the shell left the runner's workers going. There is no TERM a tool could trap to
// undo what it changed, so a mutation tool working in place can be cut mid-mutant here.
func stopProcessGroup(cmd *exec.Cmd, _ time.Duration) error {
	if cmd.Process == nil {
		return nil
	}
	if err := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run(); err != nil {
		return cmd.Process.Kill()
	}
	return nil
}

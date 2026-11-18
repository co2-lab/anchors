// @anchors
//   code: BUCBD
//   ref: BDGRN

//go:build !windows

package quality

import (
	"os/exec"
	"syscall"
	"time"
)

// ownProcessGroup starts the command in a group of its own, so the deadline can stop the
// runner's workers along with the shell.
func ownProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// stopProcessGroup asks the command's whole group to stop (TERM), and kills what is left
// of it once the grace is over. The group is watched rather than killed blindly: once it
// is gone there is nothing to kill, and the id is not reused within the grace.
func stopProcessGroup(cmd *exec.Cmd, grace time.Duration) error {
	pgid := -cmd.Process.Pid
	if err := syscall.Kill(pgid, syscall.SIGTERM); err != nil {
		return err
	}
	go func() {
		end := time.Now().Add(grace)
		for time.Now().Before(end) {
			if syscall.Kill(pgid, 0) != nil {
				return // the group is gone: it stopped on the TERM
			}
			time.Sleep(100 * time.Millisecond)
		}
		_ = syscall.Kill(pgid, syscall.SIGKILL) // @resilient: a group that exited between the check and the kill has nothing left to stop
	}()
	return nil
}

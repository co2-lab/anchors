// @anchors
//   ref: MPLCK

//go:build !windows

package mapx

import (
	"errors"
	"syscall"
)

// processAlive says whether a process of this machine still exists: signal 0 checks it
// without touching it. EPERM means it exists under another user.
func processAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

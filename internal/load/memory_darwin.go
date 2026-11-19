// @anchors
//   code: LDMDV
//   ref: LDQUS

//go:build darwin

package load

import "golang.org/x/sys/unix"

// totalMemory is the machine's physical memory in bytes; 0 when it cannot be read.
func totalMemory() int64 {
	n, _ := unix.SysctlUint64("hw.memsize") // 0 when it cannot be read
	return int64(n)
}

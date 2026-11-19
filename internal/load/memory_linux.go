// @anchors
//   code: LDMLV
//   ref: LDQUS

//go:build linux

package load

import "golang.org/x/sys/unix"

// totalMemory is the machine's physical memory in bytes; 0 when it cannot be read.
func totalMemory() int64 {
	var info unix.Sysinfo_t
	if unix.Sysinfo(&info) != nil {
		return 0
	}
	return int64(info.Totalram) * int64(info.Unit)
}

// @anchors
//   code: RTDRW
//   ref: PRRPR

package config

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

// diskPath is the path the kernel holds for the folder or file: its case as the disk names
// it, and its links resolved (fcntl F_GETPATH).
func diskPath(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	const fGetPath = 50
	buf := make([]byte, 1024)
	if _, _, errno := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), fGetPath, uintptr(unsafe.Pointer(&buf[0]))); errno != 0 {
		return "", false
	}
	if i := bytes.IndexByte(buf, 0); i >= 0 {
		buf = buf[:i]
	}
	return string(buf), len(buf) > 0
}

// @anchors
//   code: LDMWV
//   ref: LDQUS

//go:build windows

package load

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// memoryStatusEx is the MEMORYSTATUSEX GlobalMemoryStatusEx fills.
type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

var globalMemoryStatusEx = windows.NewLazySystemDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")

// totalMemory is the machine's physical memory in bytes; 0 when it cannot be read.
func totalMemory() int64 {
	m := memoryStatusEx{Length: uint32(unsafe.Sizeof(memoryStatusEx{}))}
	if r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&m))); r == 0 {
		return 0
	}
	return int64(m.TotalPhys)
}

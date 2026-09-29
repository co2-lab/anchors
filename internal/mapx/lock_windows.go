//go:build windows

package mapx

import (
	"errors"
	"syscall"
)

// processAlive says whether a process of this machine still exists. Windows has no signal
// 0: the process is opened for a query and its exit code read — a live one answers
// STILL_ACTIVE. Access denied means it exists under another user.
//
// It answered "alive" for every pid before, so a crashed writer held the map until the
// lock aged out, a minute later.
func processAlive(pid int) bool {
	const (
		queryLimitedInformation = 0x1000 // PROCESS_QUERY_LIMITED_INFORMATION
		stillActive             = 259    // STILL_ACTIVE
	)
	h, err := syscall.OpenProcess(queryLimitedInformation, false, uint32(pid))
	if err != nil {
		return errors.Is(err, syscall.ERROR_ACCESS_DENIED)
	}
	defer syscall.CloseHandle(h)
	var code uint32
	if err := syscall.GetExitCodeProcess(h, &code); err != nil {
		return true // cannot tell: the lock's age decides
	}
	return code == stillActive
}

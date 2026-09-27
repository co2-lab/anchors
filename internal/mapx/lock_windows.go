//go:build windows

package mapx

// processAlive cannot ask Windows cheaply whether a pid still exists; the lock's age
// decides instead (lockStaleAfter).
func processAlive(pid int) bool { return true }

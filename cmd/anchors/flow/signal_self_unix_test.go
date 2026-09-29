//go:build !windows

package flow

import (
	"os"
	"syscall"
)

// signalSelf sends SIGTERM to this process — how a test ends the watch loop it runs.
func signalSelf() error { return syscall.Kill(os.Getpid(), syscall.SIGTERM) }

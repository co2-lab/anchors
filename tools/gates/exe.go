package main

import "runtime"

// exeSuffix is the extension of an executable on this system: `.exe` on Windows.
func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

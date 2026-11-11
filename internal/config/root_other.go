//go:build !darwin && !windows

// @anchors
//   code: RTOTH
//   ref: PRRPR

package config

// diskPath says nothing on Linux and the other Unixes: their file systems tell case apart.
func diskPath(string) (string, bool) { return "", false }

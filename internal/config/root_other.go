//go:build !darwin

// @anchors
//   code: RTOTH
//   ref: PRRPR

package config

// diskPath says nothing outside macOS: Linux tells case apart, and a Windows working
// directory comes from the system in its own case.
func diskPath(string) (string, bool) { return "", false }

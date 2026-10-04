// @anchors
//   code: SSWTS
//   ref: WTCHA

//go:build windows

package flow

import "errors"

// signalSelf cannot be done on Windows: a process does not send itself SIGTERM there. The
// tests that need it skip on Windows before reaching it.
func signalSelf() error { return errors.New("no SIGTERM on Windows") }

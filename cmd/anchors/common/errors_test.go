package common

import (
	"errors"
	"strings"
	"testing"
)

// The message names the path that is not governed, so whoever reads it knows WHICH file the
// Structure does not reach — and the exit code is the one scripts branch on.
func TestErrNotGoverned_namesThePath(t *testing.T) {
	var err error = ErrNotGoverned{Path: "docs/notes.txt"}
	msg := err.Error()
	if !strings.Contains(msg, `"docs/notes.txt"`) {
		t.Errorf("the message does not name the path: %q", msg)
	}
	if !strings.Contains(msg, "not governed") {
		t.Errorf("the message does not say the path is not governed: %q", msg)
	}
	var ng ErrNotGoverned
	if !errors.As(err, &ng) || ng.Path != "docs/notes.txt" {
		t.Errorf("errors.As lost the path: %+v", ng)
	}
	if ExitNotGoverned != 3 {
		t.Errorf("ExitNotGoverned = %d; scripts branch on 3", ExitNotGoverned)
	}
}

// @anchors
//   code: HLTSH
//   layer: teste

package common

import (
	"io"
	"os"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"
)

// captureStdout collects what f writes to os.Stdout. The helpers of this package print
// with fmt.Println directly, so asserting on the text the user reads needs the descriptor.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	return testkit.CaptureStdout(t, f)
}

// withStdin replaces os.Stdin with a pipe that carries input, for the duration of the test.
func withStdin(t *testing.T, input string) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(w, input); err != nil {
		t.Fatal(err)
	}
	w.Close()
	old := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = old; r.Close() })
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

package quality

import (
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	return testkit.CaptureStdout(t, fn)
}

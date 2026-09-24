package mapcmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/mapx"
)

// The report's mtime against the file's: a file changed after the report was written is
// not the file the report measured. The old integration lcov ingested today was stamped
// fresh, and its old line numbers diluted the fresh unit coverage (97% read as 69%).
func TestMarkPredating_fileEditedAfterTheReport(t *testing.T) {
	root := t.TempDir()
	report := filepath.Join(root, "lcov.info")
	edited := filepath.Join(root, "edited.ts")
	untouched := filepath.Join(root, "untouched.ts")
	for _, p := range []string{report, edited, untouched} {
		if err := os.WriteFile(p, []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t0 := time.Now().Add(-2 * time.Hour)
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.Chtimes(untouched, t0, t0))
	must(os.Chtimes(report, t0.Add(time.Minute), t0.Add(time.Minute)))
	must(os.Chtimes(edited, t0.Add(time.Hour), t0.Add(time.Hour)))

	byFile := map[string]mapx.FileCov{"edited.ts": {Total: 1}, "untouched.ts": {Total: 1}, "gone.ts": {Total: 1}}
	markPredating(byFile, root, report)
	if !byFile["edited.ts"].Predates {
		t.Error("a file edited after the report was not flagged")
	}
	if byFile["untouched.ts"].Predates {
		t.Error("a file older than the report was flagged")
	}
	if byFile["gone.ts"].Predates {
		t.Error("a file that is not on disk was flagged")
	}
}

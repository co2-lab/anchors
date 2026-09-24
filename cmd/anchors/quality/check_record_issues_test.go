package quality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
)

// A LOCAL check stamps the map but opens no issue: it runs on work in progress, and in
// blue-eyes an agent's local check filed two `[docs-fresh]` cards for a state that existed
// only on its machine. With issues on (CI, or --record-issues) the failure is filed.
func TestRecordCheck_issuesOnlyWhenOn(t *testing.T) {
	issue.UseFiles()
	fail := gate.Profile{Results: []gate.Result{{
		Gate: "docs-fresh", Target: "a.spec.md", Verdict: gate.Fail, Blocking: true, Detail: "stale",
	}}}
	run := func(on bool) (string, int) {
		root := t.TempDir()
		mapPath := filepath.Join(root, mapx.DefaultPath)
		g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "r1"}}}
		if err := recordCheck(root, mapPath, g, fail, on); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(mapPath); err != nil {
			t.Fatalf("the map was not saved (on=%v): %v", on, err)
		}
		n := 0
		_ = filepath.Walk(filepath.Join(root, issue.Dir), func(p string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				n++
			}
			return nil
		})
		return root, n
	}
	if _, n := run(false); n != 0 {
		t.Errorf("a local check (issues off) opened %d issue(s)", n)
	}
	if _, n := run(true); n == 0 {
		t.Error("with issues on, the blocking failure opened no issue")
	}
}

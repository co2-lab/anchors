package mapcmd

import (
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// The suite key must be the same on every machine. A report OUTSIDE the repository
// (CI writing `/tmp/junit.xml`) used to become `../../../../tmp/junit.xml`, whose depth
// depends on where the runner checked the repository out — so each machine wrote its own
// entry and none ever replaced another.
func TestSuiteKey_stableAcrossMachines(t *testing.T) {
	a := suiteKey("/home/runner/work/anchors/anchors", "/tmp/junit.xml")
	b := suiteKey("/Users/dev/code/anchors", "/tmp/junit.xml")
	if a != b {
		t.Errorf("the same external report got two keys: %q and %q", a, b)
	}
	if a != "external/junit.xml" {
		t.Errorf("external report key = %q, want external/junit.xml", a)
	}

	root := t.TempDir()
	in := suiteKey(root, filepath.Join(root, "apps", "mobile", "junit.xml"))
	if in != "apps/mobile/junit.xml" {
		t.Errorf("a report inside the repo should be keyed by its path from the root, got %q", in)
	}
}

// Two report paths resolving to the SAME node: the result must not depend on map order.
func TestResolveByFile_collisionIsDeterministic(t *testing.T) {
	// The workspace-relative path SORTS FIRST here ("src/…" < "web/…"), so the rule that
	// the exact node ID wins is what decides — not alphabetical luck.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "web/src/x.tsx", Kind: mapx.KindCode}}}
	byFile := map[string]int{
		"src/x.tsx":     1, // workspace-relative
		"web/src/x.tsx": 2, // already the node ID
	}
	for i := 0; i < 50; i++ {
		out := resolveByFile(g, mapx.KindCode, byFile, "/r", "/r/web/lcov.info")
		if out["web/src/x.tsx"] != 2 {
			t.Fatalf("run %d: the path that IS the node ID must win, got %d", i, out["web/src/x.tsx"])
		}
	}
}

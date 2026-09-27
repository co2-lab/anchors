package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// projectWithMap sets up a minimal project and returns its root, graph and config.
func projectWithMap(t *testing.T, content string) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(
		"version: 2\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "U.spec.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(filepath.Join(root, "anchors.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	g := mapx.Build(files, cfg, nil)
	return root, g, cfg
}

// An OLD map has to be reported — and "old" is different from "absent".
//
// The pre-commit already refused a governed file outside the map (the new file). What got
// through was the file that IS in the map with the revision of an earlier version: the map
// was built, and the work went on after it. Measured in the reference project: seven of the
// twelve rejections of the `gates` pipeline were this.
func TestStaleMap_reportsTheFileEditedAfterTheBuild(t *testing.T) {
	t.Run("MPSTM-B01: A spec edited after the map build is named as stale", func(t *testing.T) {})
	t.Run("MPSTM-I01: A freshly rebuilt map is never stale", func(t *testing.T) {})
	t.Run("MPSTM-X01: Asking about staleness leaves the map as it was", func(t *testing.T) {})
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — the rule.\n")
	oldRev := g.Nodes[0].Rev

	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Fatalf("a freshly built map was reported as stale: %v", v)
	}

	// The edit after the build — the case of the seven.
	if err := os.WriteFile(filepath.Join(root, "U.spec.md"),
		[]byte("# U\n\nUMA-B01 — the rule.\nUMA-B02 — another.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	v := StaleMapNodes(root, g, cfg)
	if len(v) != 1 || !strings.HasSuffix(v[0], "U.spec.md") {
		t.Errorf("the edit after `map build` was not reported: %v", v)
	}
	if g.Nodes[0].Rev != oldRev {
		t.Errorf("asking about staleness changed the map's revision: %q → %q", oldRev, g.Nodes[0].Rev)
	}
}

// ABSENT is not OLD, and mixing them would print a message that asks for the wrong fix.
func TestStaleMap_ignoresAFileThatIsGone(t *testing.T) {
	t.Run("MPSTM-B02: A file removed after the map build is not named", func(t *testing.T) {})
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — the rule.\n")
	if err := os.Remove(filepath.Join(root, "U.spec.md")); err != nil {
		t.Fatal(err)
	}
	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Errorf("a REMOVED file was reported as a stale map: %v — they are different cases", v)
	}
}

// The comparison is by HASH, not by date: `updated_at` moves on a checkout without the
// content moving, and does not move on an edit that preserves the mtime.
func TestStaleMap_doesNotLookAtTheDate(t *testing.T) {
	t.Run("MPSTM-B03: Staleness follows the content, not the modification time", func(t *testing.T) {})
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — the rule.\n")
	file := filepath.Join(root, "U.spec.md")

	// mtime changed, CONTENT intact: not a stale map.
	old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(file, old, old); err != nil {
		t.Fatal(err)
	}
	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Errorf("a changed mtime with intact content was reported as stale: %v", v)
	}

	// CONTENT changed, mtime restored to the old one: it IS a stale map.
	if err := os.WriteFile(file, []byte("# U\n\nUMA-B01 — something else.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(file, old, old); err != nil {
		t.Fatal(err)
	}
	if v := StaleMapNodes(root, g, cfg); len(v) != 1 {
		t.Errorf("changed content with an old mtime was NOT reported: %v", v)
	}
}

func TestStaleMap_emptyRevisionIsNotCompared(t *testing.T) {
	t.Run("MPSTM-B04: A node with an empty recorded revision is not named", func(t *testing.T) {})
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — the rule.\n")
	g.Nodes[0].Rev = ""
	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Errorf("a node with no recorded revision was reported: %v", v)
	}
}

func TestStaleMap_noMapAndUnwalkableRoot(t *testing.T) {
	t.Run("MPSTM-B05: No map names nothing", func(t *testing.T) {})
	t.Run("MPSTM-E01: A root that cannot be walked names nothing and raises nothing", func(t *testing.T) {})
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — the rule.\n")
	if v := StaleMapNodes(root, nil, cfg); v != nil {
		t.Errorf("no map must name nothing: %v", v)
	}
	if v := StaleMapNodes(filepath.Join(root, "does-not-exist"), g, cfg); v != nil {
		t.Errorf("an unwalkable root must name nothing: %v", v)
	}
}

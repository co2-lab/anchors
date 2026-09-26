package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// The two legitimate origins of an argument: relative to the ROOT (what Anchors' prompts
// print) and relative to the CWD or absolute (what a shell hands over). Both must land on
// the node ID, with forward slashes.
func TestRelTo_resolvesToTheNodeID(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "pkg", "a"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "pkg", "a", "x.spec.md"), "# x\n")

	if got := RelTo(root, "pkg/a/../a/x.spec.md"); got != "pkg/a/x.spec.md" {
		t.Errorf("root-relative path: got %q", got)
	}
	if got := RelTo(root, filepath.Join(root, "pkg", "a", "x.spec.md")); got != "pkg/a/x.spec.md" {
		t.Errorf("absolute path: got %q", got)
	}
	// Not under the root and not existing: resolved against the CWD, relative to the root.
	cwd, _ := os.Getwd()
	want, _ := filepath.Rel(root, filepath.Join(cwd, "ghost.ts"))
	if got := RelTo(root, "ghost.ts"); got != filepath.ToSlash(want) {
		t.Errorf("cwd-relative path: got %q, want %q", got, filepath.ToSlash(want))
	}
}

func TestNodeExists(t *testing.T) {
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a/x.spec.md"}, {ID: "a/x.ts"}}}
	if !NodeExists(g, "a/x.ts") {
		t.Error("a node in the map was not found")
	}
	if NodeExists(g, "a/y.ts") {
		t.Error("a node absent from the map was found")
	}
	if NodeExists(nil, "a/x.ts") {
		t.Error("a nil map has no node")
	}
}

func TestRelSlug_dropsOnlyTheLastExtension(t *testing.T) {
	for in, want := range map[string]string{
		"pkg/a/x.ts":      "pkg/a/x",
		"pkg/a/x.spec.md": "pkg/a/x.spec",
		"Makefile":        "Makefile",
	} {
		if got := RelSlug(in); got != want {
			t.Errorf("RelSlug(%q) = %q, want %q", in, got, want)
		}
	}
}

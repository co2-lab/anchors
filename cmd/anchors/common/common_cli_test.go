// @anchors
//   ref: CMCLC

package common

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// The message names the path that is not governed, so whoever reads it knows WHICH file the
// Structure does not reach — and the exit code is the one scripts branch on.
func TestErrNotGoverned_namesThePath(t *testing.T) {
	t.Run("CMCLC-B01: The not-governed error names the path", func(t *testing.T) {})
	t.Run("CMCLC-B02: The not-governed exit code is 3", func(t *testing.T) {})
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

// The two legitimate origins of an argument: relative to the ROOT (what Anchors' prompts
// print) and relative to the CWD or absolute (what a shell hands over). Both must land on
// the node ID, with forward slashes.
func TestRelTo_resolvesToTheNodeID(t *testing.T) {
	t.Run("CMCLC-B03: A root-relative path is kept, cleaned", func(t *testing.T) {})
	t.Run("CMCLC-B04: A path absent under the root is resolved from the working directory", func(t *testing.T) {})
	t.Run("CMCLC-I01: Root-relative and absolute names of one file resolve to one node", func(t *testing.T) {})
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
	want, err := filepath.Rel(root, filepath.Join(cwd, "ghost.ts"))
	if err != nil {
		want = "ghost.ts" // on another drive (Windows) there is no relative path: kept as given
	}
	if got := RelTo(root, "ghost.ts"); got != filepath.ToSlash(want) {
		t.Errorf("cwd-relative path: got %q, want %q", got, filepath.ToSlash(want))
	}
}

// A path that cannot be expressed relative to the root comes back as given: a relative
// root and an absolute argument cannot be related, and an empty answer would hide why.
func TestRelTo_unrelatablePathIsKept(t *testing.T) {
	t.Run("CMCLC-E01: A path that cannot be related to the root is kept as given", func(t *testing.T) {})
	if got := RelTo("relroot", "/abs/x"); got != "/abs/x" {
		t.Errorf("got %q, want the argument as given", got)
	}
}

func TestNodeExists(t *testing.T) {
	t.Run("CMCLC-B05: A node exists only by its exact identifier", func(t *testing.T) {})
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
	t.Run("CMCLC-B06: A task slug drops only the last extension", func(t *testing.T) {})
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

// The test binary is built without -ldflags, so it carries the defaults an unstamped
// build reports.
func TestVersion_unstampedBuildSaysDev(t *testing.T) {
	t.Run("CMCLC-B07: An unstamped build says it is a development build", func(t *testing.T) {})
	if Version != "dev" || Commit != "none" || Date != "unknown" {
		t.Errorf("unstamped build reports %q / %q / %q, want dev / none / unknown", Version, Commit, Date)
	}
}

// A Windows-native path given on the command line becomes the map's form.
func TestRelTo_windowsSeparators(t *testing.T) {
	t.Run("CMCLC-B03: A root-relative path is kept, cleaned", func(t *testing.T) {})
	if runtime.GOOS != "windows" {
		t.Skip("`\\` is a separator only on Windows")
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "src", "a.ts"), "x\n")
	if got := RelTo(root, `src\a.ts`); got != "src/a.ts" {
		t.Errorf("a native relative path, got %q", got)
	}
	if got := RelTo(root, filepath.Join(root, "src", "a.ts")); got != "src/a.ts" {
		t.Errorf("a native absolute path, got %q", got)
	}
}

func TestFileArgs_theListOfFiles(t *testing.T) {
	t.Run("CMCLC-B08: A list of files is read the same way by every command", func(t *testing.T) {})
	want := []string{"a.ts", "b.ts", "c.ts"}
	if got := FileArgs([]string{"a.ts", "b.ts, ,c.ts", "a.ts"}); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("each file once, in order, got %v", got)
	}
	var seen []string
	cmd := TakesFiles(&cobra.Command{Use: "x <file>...", RunE: func(_ *cobra.Command, args []string) error { seen = args; return nil }})
	cmd.SetArgs([]string{"a.ts,b.ts", "c.ts"})
	if err := cmd.Execute(); err != nil || strings.Join(seen, "|") != strings.Join(want, "|") {
		t.Errorf("the command gets the files, got %v %v", seen, err)
	}
	if cmd.Annotations[FilesAnnotation] != "true" {
		t.Error("the command is marked as taking files")
	}
}

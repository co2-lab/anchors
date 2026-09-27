package config

import (
	"reflect"
	"testing"
)

func TestMarkersFor(t *testing.T) {
	t.Run("CMMRC-B01: A known extension answers with its line-comment prefixes", func(t *testing.T) {})
	t.Run("CMMRC-B02: An unknown extension answers with no prefix", func(t *testing.T) {})
	if got := MarkersFor(".php"); !reflect.DeepEqual(got, []string{"//", "#"}) {
		t.Errorf("MarkersFor(.php) = %v, want [// #]", got)
	}
	if got := MarkersFor(".sql"); !reflect.DeepEqual(got, []string{"--"}) {
		t.Errorf("MarkersFor(.sql) = %v, want [--]", got)
	}
	if got := MarkersFor(".unknown"); got != nil {
		t.Errorf("MarkersFor(.unknown) = %v, want nil", got)
	}
}

func TestLineCommentFor(t *testing.T) {
	t.Run("CMMRC-B03: The line comment of a path follows its last extension, ignoring case", func(t *testing.T) {})
	t.Run("CMMRC-B04: A markup file gets the opening of a block comment", func(t *testing.T) {})
	t.Run("CMMRC-B05: A path with no extension or an unknown one gets the hash comment", func(t *testing.T) {})
	t.Run("CMMRC-B06: An extension with several prefixes gets the first one declared", func(t *testing.T) {})
	for path, want := range map[string]string{
		"internal/x/y.go":      "//",
		"web/a.test.ts":        "//", // the real extension is the last one
		"scripts/run_test.py":  "#",
		"db/001.SQL":           "--", // case-insensitive
		".SQL":                 "--", // the extension is the whole name: still lowered
		"docs/README.md":       "<!--",
		"Makefile":             "#", // no extension: the visible-mistake default
		"assets/logo.svg":      "#", // unknown extension: same default
		"src/lib.rs":           "//",
		"config/anchors.yaml":  "#",
		"src/Main.hs":          "--",
		"templates/index.html": "<!--",
		"web/index.php":        "//", // `//` and `#` are both valid: the first declared wins
	} {
		if got := LineCommentFor(path); got != want {
			t.Errorf("LineCommentFor(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestLineCommentFor_agreesWithTheTable(t *testing.T) {
	t.Run("CMMRC-I01: For every extension of the table, the line comment is the table's first prefix", func(t *testing.T) {})
	for ext, markers := range CommentMarkers {
		if got := LineCommentFor("file" + ext); got != markers[0] {
			t.Errorf("LineCommentFor(file%s) = %q, the table says %q", ext, got, markers[0])
		}
	}
}

func TestMarkersFor_exactLookup(t *testing.T) {
	t.Run("CMMRC-X01: The prefix lookup takes the extension as given and does not normalise it", func(t *testing.T) {})
	for _, in := range []string{".GO", "go", "main.go"} {
		if got := MarkersFor(in); got != nil {
			t.Errorf("MarkersFor(%q) = %v, want nil (only the exact extension is looked up)", in, got)
		}
	}
}

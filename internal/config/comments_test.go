package config

import (
	"reflect"
	"testing"
)

func TestMarkersFor(t *testing.T) {
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
	for path, want := range map[string]string{
		"internal/x/y.go":      "//",
		"web/a.test.ts":        "//", // the real extension is the last one
		"scripts/run_test.py":  "#",
		"db/001.SQL":           "--", // case-insensitive
		"docs/README.md":       "<!--",
		"Makefile":             "#", // no extension: the visible-mistake default
		"assets/logo.svg":      "#", // unknown extension: same default
		"src/lib.rs":           "//",
		"config/anchors.yaml":  "#",
		"src/Main.hs":          "--",
		"templates/index.html": "<!--",
	} {
		if got := LineCommentFor(path); got != want {
			t.Errorf("LineCommentFor(%q) = %q, want %q", path, got, want)
		}
	}
}

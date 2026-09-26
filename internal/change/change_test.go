package change

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestKeyAndPath(t *testing.T) {
	c := Change{Stage: "code", Unit: "internal/gate/mock stamped.go"}
	if got, want := c.Key(), "code--internal-gate-mock-stamped"; got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}
	root := t.TempDir()
	if got, want := c.Path(root), filepath.Join(root, "changes", "code--internal-gate-mock-stamped.md"); got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
	// Windows separators and dotted names collapse the same way.
	if got := slug(`a\b.c.go`); got != "a-b-c" {
		t.Fatalf("slug = %q, want a-b-c", got)
	}
}

func TestRender_filledSections(t *testing.T) {
	c := Change{
		Stage:     "test",
		Unit:      "internal/x/y.go",
		Files:     []string{"internal/x/y_test.go", "internal/x/y.feature"},
		Intent:    "  Proves the Y rules.  \n",
		Decisions: []string{"kept the old name"},
		Uncovered: []string{"Y-B03 has no test"},
		Date:      "2026-09-26",
		Agent:     "worker-2",
	}
	out := c.Render()
	for _, want := range []string{
		"<!-- @anchors\n  layer: change\n  stage: test\n  unit: internal/x/y.go\n  date: 2026-09-26\n  agent: worker-2\n-->\n",
		"# Delivery: test of `internal/x/y.go`\n",
		"## What was done\n\nProves the Y rules.\n\n",
		"## Files\n\n- `internal/x/y_test.go`\n- `internal/x/y.feature`\n",
		"## Decisions the ruler did not decide\n\n- kept the old name\n",
		"## What is NOT proven\n\n- Y-B03 has no test\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Render() lacks %q\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, "none —") || strings.Contains(out, "nothing —") {
		t.Errorf("filled sections must not carry the empty-section sentence:\n%s", out)
	}
}

func TestRender_emptySectionsAreStillEmitted(t *testing.T) {
	out := Change{Stage: "spec", Unit: "a.go", Date: "2026-01-02"}.Render()
	if strings.Contains(out, "agent:") {
		t.Errorf("an empty Agent must not be written:\n%s", out)
	}
	for _, want := range []string{
		"## Decisions the ruler did not decide\n\nnone — ",
		"## What is NOT proven\n\nnothing — ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Render() lacks %q\n---\n%s", want, out)
		}
	}
}

func TestSavePendingMarkReviewed(t *testing.T) {
	root := t.TempDir()

	// No changes/ directory yet: nothing pending, and no error.
	if got, err := Pending(root); err != nil || got != nil {
		t.Fatalf("Pending on a fresh root = %v, %v; want nil, nil", got, err)
	}

	b := Change{Stage: "spec", Unit: "b.go", Intent: "first", Date: "2026-01-01"}
	a := Change{Stage: "code", Unit: "a.go", Intent: "second", Date: "2026-01-01"}
	pb, err := Save(root, b)
	if err != nil {
		t.Fatal(err)
	}
	pa, err := Save(root, a)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(pb)
	if err != nil || string(data) != b.Render() {
		t.Fatalf("saved file = %q (%v), want the rendered change", data, err)
	}

	// Noise that Pending must skip: a non-markdown file and a subdirectory.
	os.WriteFile(filepath.Join(root, Dir, "notes.txt"), []byte("x"), 0o644)
	os.MkdirAll(filepath.Join(root, Dir, "sub.md"), 0o755)

	got, err := Pending(root)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{pa, pb}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Pending = %v, want %v (sorted, .md files only)", got, want)
	}

	// A second delivery of the same stage and unit replaces the first.
	b.Intent = "rewritten"
	if p, err := Save(root, b); err != nil || p != pb {
		t.Fatalf("re-Save = %q, %v; want the same path %q", p, err, pb)
	}
	if got, _ := Pending(root); len(got) != 2 {
		t.Fatalf("re-Save must replace, not add: %v", got)
	}

	dest, err := MarkReviewed(root, pa)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(root, ReviewedDir, filepath.Base(pa)); dest != want {
		t.Fatalf("MarkReviewed dest = %q, want %q", dest, want)
	}
	if _, err := os.Stat(pa); !os.IsNotExist(err) {
		t.Fatalf("the reviewed record must leave changes/: %v", err)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("the reviewed record must exist in the history: %v", err)
	}
	if got, _ := Pending(root); !reflect.DeepEqual(got, []string{pb}) {
		t.Fatalf("Pending after review = %v, want [%s]", got, pb)
	}

	if _, err := MarkReviewed(root, filepath.Join(root, Dir, "missing.md")); err == nil {
		t.Fatal("MarkReviewed of a missing record must fail")
	}
}

func TestPending_unreadableDirIsAnError(t *testing.T) {
	root := t.TempDir()
	// changes/ exists as a FILE: ReadDir fails with something other than NotExist.
	os.WriteFile(filepath.Join(root, Dir), []byte("x"), 0o644)
	if _, err := Pending(root); err == nil {
		t.Fatal("Pending must report a changes/ it cannot read")
	}
	if _, err := Save(root, Change{Stage: "code", Unit: "a.go"}); err == nil {
		t.Fatal("Save must report a changes/ it cannot create")
	}
}

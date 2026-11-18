// @anchors
//   code: CHTSB
//   ref: CHRCC

package change

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"
)

func TestKeyAndPath(t *testing.T) {
	t.Run("CHRCC-B01: The key joins the stage and the normalised unit", func(t *testing.T) {})
	t.Run("CHRCC-B02: A pending record lives in the changes folder under its key", func(t *testing.T) {})
	c := Change{Stage: "code", Unit: "internal/gate/mock stamped.go"}
	// The slug, then a short hash of the unit that keeps the key injective (CHRCC-I02).
	if got, want := c.Key(), "code--internal-gate-mock-stamped-4dc4c23f"; got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}
	root := t.TempDir()
	if got, want := c.Path(root), filepath.Join(root, "changes", "code--internal-gate-mock-stamped-4dc4c23f.md"); got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
	// Windows separators and dotted names collapse the same way.
	if got := slug(`a\b.c.go`); got != "a-b-c" {
		t.Fatalf("slug = %q, want a-b-c", got)
	}
}

func TestRender_filledSections(t *testing.T) {
	t.Run("CHRCC-B03: The header carries the stage, the unit, the date and the agent only when named", func(t *testing.T) {})
	t.Run("CHRCC-B04: The record states the intent and the touched files", func(t *testing.T) {})
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
	t.Run("CHRCC-B05: Empty decision and proof sections are still written", func(t *testing.T) {})
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
	t.Run("CHRCC-B06: Saving the same stage and unit again replaces the record", func(t *testing.T) {})
	t.Run("CHRCC-B07: The pending list is the sorted markdown files of the changes folder", func(t *testing.T) {})
	t.Run("CHRCC-B08: Marking a record reviewed moves it to the history under the same name", func(t *testing.T) {})
	t.Run("CHRCC-I01: A reviewed record leaves the pending list and stays in the history", func(t *testing.T) {})
	t.Run("CHRCC-X01: A subfolder of changes is never listed as pending", func(t *testing.T) {})
	t.Run("CHRCC-E03: Marking a missing record reviewed fails", func(t *testing.T) {})
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
	t.Run("CHRCC-E01: A changes folder that cannot be read is an error", func(t *testing.T) {})
	t.Run("CHRCC-E02: A changes folder that cannot be created fails the save", func(t *testing.T) {})
	testkit.SkipWithoutPOSIXPermissions(t)
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

func TestKey_differentUnitsNeverCollide(t *testing.T) {
	t.Run("CHRCC-I02: Two different units never share a key", func(t *testing.T) {})
	// Before: all three slugged to "a-b", so one delivery replaced another unit's record.
	units := []string{"a/b.go", "a-b.go", "a/b.ts", "a.b.go", "a b.go"}
	seen := map[string]string{}
	for _, u := range units {
		k := Change{Stage: "code", Unit: u}.Key()
		if other, dup := seen[k]; dup {
			t.Fatalf("units %q and %q share the key %q", other, u, k)
		}
		seen[k] = u
	}
	// The two separators still name the same unit.
	if a, b := (Change{Stage: "code", Unit: `a\b.go`}).Key(), (Change{Stage: "code", Unit: "a/b.go"}).Key(); a != b {
		t.Fatalf("`a\\b.go` and `a/b.go` must share a key: %q vs %q", a, b)
	}
	root := t.TempDir()
	p1, _ := Save(root, Change{Stage: "code", Unit: "a/b.go", Intent: "one"})
	p2, _ := Save(root, Change{Stage: "code", Unit: "a-b.go", Intent: "two"})
	if got, _ := Pending(root); len(got) != 2 || p1 == p2 {
		t.Fatalf("two units must keep two records, got %v", got)
	}
}

func TestMarkReviewed_keepsAnEarlierReviewOfTheSameKey(t *testing.T) {
	t.Run("CHRCC-B09: A second review of the same key keeps the first in the history", func(t *testing.T) {})
	root := t.TempDir()
	c := Change{Stage: "code", Unit: "a.go", Intent: "first"}
	p, _ := Save(root, c)
	d1, err := MarkReviewed(root, p)
	if err != nil {
		t.Fatal(err)
	}
	c.Intent = "second"
	p, _ = Save(root, c)
	d2, err := MarkReviewed(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if d1 == d2 {
		t.Fatalf("the second review landed on the first one's file %q", d1)
	}
	b1, _ := os.ReadFile(d1)
	b2, _ := os.ReadFile(d2)
	if !strings.Contains(string(b1), "first") || !strings.Contains(string(b2), "second") {
		t.Fatalf("history lost a review:\n%s\n---\n%s", b1, b2)
	}
	if !strings.HasSuffix(d2, ".md") || filepath.Dir(d2) != filepath.Join(root, ReviewedDir) {
		t.Fatalf("the second review must stay a record in the history, got %q", d2)
	}
}

func TestRender_singleLineFieldsCannotBreakTheHeader(t *testing.T) {
	t.Run("CHRCC-B10: A line break or a comment close in a field cannot break the header", func(t *testing.T) {})
	out := Change{Stage: "code", Unit: "a.go\nlayer: forged\n-->", Date: "2026-01-02", Agent: "w\r\nx"}.Render()
	header, _, ok := strings.Cut(out, "-->\n")
	if !ok {
		t.Fatalf("no header close:\n%s", out)
	}
	if strings.Contains(header, "\nlayer: forged") {
		t.Fatalf("the unit forged a header line:\n%s", header)
	}
	if !strings.Contains(header, `unit: a.go\nlayer: forged`) || !strings.Contains(header, `agent: w\r\nx`) {
		t.Fatalf("the field must survive escaped on its own line:\n%s", header)
	}
	if !strings.Contains(header, "  date: 2026-01-02\n") {
		t.Fatalf("the date line must follow intact:\n%s", header)
	}
	title := strings.SplitN(strings.TrimPrefix(out, header+"-->\n"), "\n", 2)[0]
	if !strings.HasPrefix(title, "# Delivery: code of `a.go\\nlayer: forged") {
		t.Fatalf("the title must stay one line, got %q", title)
	}
}

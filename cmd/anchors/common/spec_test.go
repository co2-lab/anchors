package common

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestCodeDoHeaderSpec(t *testing.T) {
	for _, c := range []struct{ name, content, want string }{
		{"html header", "<!-- @anchors\n  code: RLSGR\n  updated_at: 2026-09-21\n-->\n# R\n", "RLSGR"},
		{"go comment", "// code: ALFDL\npackage x\n", "ALFDL"},
		{"no header", "# Just a title\n\nsome code: lower\n", ""},
	} {
		if got := CodeDoHeaderSpec(c.content); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

// The unit's code comes from the MAP: the exact node first, then any node of the same
// stem — the code file of a unit answers with the spec's code.
func TestCodeOfUnit_readsTheMap(t *testing.T) {
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/alpha.spec.md", Kind: mapx.KindSpec, Code: "ALPHA"},
		{ID: "pkg/alpha.ts", Kind: mapx.KindCode},
		{ID: "pkg/beta.ts", Kind: mapx.KindCode, Code: "BETAX"},
	}}
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	for unit, want := range map[string]string{
		"pkg/beta.ts":       "BETAX", // exact node
		"pkg/alpha.ts":      "ALPHA", // same stem, the spec answers
		"pkg/alpha.spec.md": "ALPHA",
		"pkg/gamma.ts":      "",
	} {
		if got := CodeOfUnit(root, unit); got != want {
			t.Errorf("CodeOfUnit(%q) = %q, want %q", unit, got, want)
		}
	}
	if got := CodeOfUnit(t.TempDir(), "pkg/beta.ts"); got != "" {
		t.Errorf("without a map there is no code, got %q", got)
	}
}

// Only the unit's own codes count, once each; a code of another unit is only CITED.
func TestCodesInFileOfUnit_keepsOnlyTheUnitsCodes(t *testing.T) {
	p := filepath.Join(t.TempDir(), "alpha.spec.md")
	writeFile(t, p, "# Alpha\n\nALPHA-B01 — first rule.\nALPHA-B02 — second, see BETAX-B07.\nALPHA-B01 again.\n")

	got, err := CodesInFileOfUnit(p, "ALPHA")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"ALPHA-B01", "ALPHA-B02"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	all, _ := CodesInFileOfUnit(p, "")
	if len(all) != 3 {
		t.Errorf("with no unit every code counts, once each: %v", all)
	}
	if _, err := CodesInFileOfUnit(filepath.Join(t.TempDir(), "missing.md"), "ALPHA"); err == nil {
		t.Error("a missing file must be an error")
	}
}

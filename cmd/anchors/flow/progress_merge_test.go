package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The merge driver writes the union into OUR side, and `[x]` beats `[ ]` in both
// directions: what either branch delivered stays delivered.
func TestProgressMergeCmd_writesTheUnionIntoOurSide(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.md")
	ours := filepath.Join(dir, "ours.md")
	theirs := filepath.Join(dir, "theirs.md")
	writeFile(t, dir, "base.md", "# P\n- [ ] a.spec.md\n- [ ] b.spec.md\n- [ ] c.spec.md\n")
	writeFile(t, dir, "ours.md", "# P\n- [x] a.spec.md\n- [ ] b.spec.md\n- [ ] c.spec.md\n")
	writeFile(t, dir, "theirs.md", "# P\n- [ ] a.spec.md\n- [x] b.spec.md\n- [ ] c.spec.md\n")

	cmd := newProgressMergeCmd()
	cmd.SetArgs([]string{base, ours, theirs})
	var err error
	msg := stderrOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(ours)
	for _, want := range []string{"- [x] a.spec.md", "- [x] b.spec.md", "- [ ] c.spec.md"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the merged file lacks %q:\n%s", want, got)
		}
	}
	// The count says what came from the other side, so a merge that unmarked nothing and
	// gained one delivery reads as such.
	if !strings.Contains(msg, "2 done") || !strings.Contains(msg, "(1 from the other side)") {
		t.Errorf("the report must count the result and what the other side added: %q", msg)
	}
}

func TestProgressMergeCmd_reportsWhichSideCouldNotBeRead(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ours.md", "- [x] a.spec.md\n")
	missing := filepath.Join(dir, "nope.md")

	cmd := newProgressMergeCmd()
	cmd.SetArgs([]string{missing, missing, filepath.Join(dir, "ours.md")})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "read our side") {
		t.Errorf("a missing ours must be named, got %v", err)
	}

	cmd = newProgressMergeCmd()
	cmd.SetArgs([]string{missing, filepath.Join(dir, "ours.md"), missing})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "read the other side") {
		t.Errorf("a missing theirs must be named, got %v", err)
	}
}

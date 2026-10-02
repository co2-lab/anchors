// @anchors
//   ref: PRMRP

package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/scan"
)

func runProgressMerge(t *testing.T, args ...string) (stderr string, err error) {
	t.Helper()
	cmd := newProgressMergeCmd()
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	cmd.SetArgs(args)
	stderr = stderrOf(t, func() { err = cmd.Execute() })
	return stderr, err
}

// The merge driver writes the union into OUR side, and `[x]` beats `[ ]` in both
// directions: what either branch delivered stays delivered.
func TestProgressMergeCmd_writesTheUnionIntoOurSide(t *testing.T) {
	t.Run("PRMRP-B02: The union is written into our side", func(t *testing.T) {})
	t.Run("PRMRP-B04: The report counts the done items and what the other side added", func(t *testing.T) {})
	t.Run("PRMRP-I01: An item done on our side stays done", func(t *testing.T) {})
	t.Run("PRMRP-X01: The driver applies the progress reader's union", func(t *testing.T) {})
	dir := t.TempDir()
	base := filepath.Join(dir, "base.md")
	ours := filepath.Join(dir, "ours.md")
	theirs := filepath.Join(dir, "theirs.md")
	writeFile(t, dir, "base.md", "# P\n- [ ] a.spec.md\n- [ ] b.spec.md\n- [ ] c.spec.md\n")
	oursText := "# P\n- [x] a.spec.md\n- [ ] b.spec.md\n- [ ] c.spec.md\n"
	theirsText := "# P\n- [ ] a.spec.md\n- [x] b.spec.md\n- [ ] c.spec.md\n"
	writeFile(t, dir, "ours.md", oursText)
	writeFile(t, dir, "theirs.md", theirsText)

	msg, err := runProgressMerge(t, base, ours, theirs)
	if err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(ours)
	for _, want := range []string{"- [x] a.spec.md", "- [x] b.spec.md", "- [ ] c.spec.md"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the merged file lacks %q:\n%s", want, got)
		}
	}
	if string(got) != scan.MergeProgress(oursText, theirsText) {
		t.Errorf("the result must be the progress reader's union:\n%s", got)
	}
	// The count says what came from the other side, so a merge that unmarked nothing and
	// gained one delivery reads as such.
	if !strings.Contains(msg, "2 done") || !strings.Contains(msg, "(1 from the other side)") {
		t.Errorf("the report must count the result and what the other side added: %q", msg)
	}

	// Nothing new from the other side: the count alone, no "from the other side".
	writeFile(t, dir, "theirs.md", "# P\n- [ ] a.spec.md\n")
	writeFile(t, dir, "ours.md", "# P\n- [x] a.spec.md\n")
	msg, err = runProgressMerge(t, base, ours, theirs)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(msg, "1 done") || strings.Contains(msg, "from the other side") {
		t.Errorf("a merge that gained nothing must say only the count: %q", msg)
	}
}

// Git calls the driver with %O %A %B; any other shape is a misconfigured driver.
func TestProgressMergeCmd_takesExactlyThreeArguments(t *testing.T) {
	t.Run("PRMRP-B01: The merge driver refuses any count other than three arguments", func(t *testing.T) {})
	dir := t.TempDir()
	writeFile(t, dir, "a.md", "- [x] a\n")
	a := filepath.Join(dir, "a.md")
	for _, args := range [][]string{{a, a}, {a, a, a, a}} {
		if _, err := runProgressMerge(t, args...); err == nil {
			t.Errorf("%d arguments must be refused", len(args))
		}
	}
}

// The base is never read: a merge cannot legitimately unmark, so there is nothing to
// learn from it.
func TestProgressMergeCmd_ignoresTheBase(t *testing.T) {
	t.Run("PRMRP-B03: A missing base does not stop the merge", func(t *testing.T) {})
	dir := t.TempDir()
	writeFile(t, dir, "ours.md", "- [x] a.spec.md\n- [ ] b.spec.md\n")
	writeFile(t, dir, "theirs.md", "- [ ] a.spec.md\n- [x] b.spec.md\n")
	ours := filepath.Join(dir, "ours.md")
	if _, err := runProgressMerge(t, filepath.Join(dir, "no-base.md"), ours, filepath.Join(dir, "theirs.md")); err != nil {
		t.Fatalf("a missing base must not stop the merge: %v", err)
	}
	if got, _ := os.ReadFile(ours); !strings.Contains(string(got), "- [x] b.spec.md") {
		t.Errorf("the union must still be written:\n%s", got)
	}
}

func TestProgressMergeCmd_reportsWhichSideCouldNotBeRead(t *testing.T) {
	t.Run("PRMRP-E01: An unreadable our side is named", func(t *testing.T) {})
	t.Run("PRMRP-E02: An unreadable other side is named and our side is untouched", func(t *testing.T) {})
	dir := t.TempDir()
	writeFile(t, dir, "ours.md", "- [x] a.spec.md\n")
	missing := filepath.Join(dir, "nope.md")

	if _, err := runProgressMerge(t, missing, missing, filepath.Join(dir, "ours.md")); err == nil || !strings.Contains(err.Error(), "read our side") {
		t.Errorf("a missing ours must be named, got %v", err)
	}

	if _, err := runProgressMerge(t, missing, filepath.Join(dir, "ours.md"), missing); err == nil || !strings.Contains(err.Error(), "read the other side") {
		t.Errorf("a missing theirs must be named, got %v", err)
	}
	if got, _ := os.ReadFile(filepath.Join(dir, "ours.md")); string(got) != "- [x] a.spec.md\n" {
		t.Errorf("our side must be untouched after a failed read:\n%s", got)
	}
}

func TestProgressMergeCmd_unwritableResultFails(t *testing.T) {
	t.Run("PRMRP-E03: A result that cannot be written fails the merge", func(t *testing.T) {})
	if os.Geteuid() == 0 {
		t.Skip("root writes read-only files")
	}
	dir := t.TempDir()
	writeFile(t, dir, "ours.md", "- [x] a.spec.md\n")
	writeFile(t, dir, "theirs.md", "- [x] b.spec.md\n")
	ours := filepath.Join(dir, "ours.md")
	if err := os.Chmod(ours, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := runProgressMerge(t, ours, ours, filepath.Join(dir, "theirs.md")); err == nil || !strings.Contains(err.Error(), "write the result") {
		t.Errorf("an unwritable result must fail the merge, got %v", err)
	}
}

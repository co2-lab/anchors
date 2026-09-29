package mapcmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"
)

type repo struct {
	t   *testing.T
	dir string
}

func newRepo(t *testing.T) repo {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	r := repo{t, t.TempDir()}
	r.git("init", "-b", "main")
	r.git("config", "user.email", "t@t")
	r.git("config", "user.name", "t")
	r.write("anchors.yaml", "version: 2\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	return r
}

func (r repo) git(args ...string) {
	r.t.Helper()
	c := exec.Command("git", args...)
	c.Dir = r.dir
	if out, err := c.CombinedOutput(); err != nil {
		r.t.Fatalf("git %v: %s", args, out)
	}
}

func (r repo) write(name, body string) {
	r.t.Helper()
	if err := os.WriteFile(filepath.Join(r.dir, name), []byte(body), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

func (r repo) read(name string) string {
	r.t.Helper()
	b, err := os.ReadFile(filepath.Join(r.dir, name))
	if err != nil {
		r.t.Fatal(err)
	}
	return string(b)
}

// renumber runs the command in the repo and returns its stdout.
func (r repo) renumber(args ...string) (string, error) {
	r.t.Helper()
	return runCmdErr(newRenumberCmd(), r.t, append([]string{"--root", r.dir}, args...)...)
}

const pricingHead = "<!-- @anchors\n  code: PRICX\n-->\n# Pricing\n\n> **PRICX-R0001:** first\n"

// forkedRepo is a repo where `main` and the branch `feat` each added a `PRICX-R0002` to
// the same spec after forking; HEAD is `feat`, and nothing of feat's is committed yet
// when uncommitted is true.
func forkedRepo(t *testing.T, uncommitted bool) repo {
	t.Helper()
	r := newRepo(t)
	r.write("Pricing.spec.md", pricingHead+"\n## Overview\n")
	r.git("add", "anchors.yaml", "Pricing.spec.md")
	r.git("commit", "-m", "base")
	r.git("checkout", "-b", "feat")
	r.git("checkout", "main")
	r.write("Pricing.spec.md", pricingHead+"\n## Overview\n\n> **PRICX-R0002:** the other PR\n")
	r.git("commit", "-am", "other PR")
	r.git("checkout", "feat")
	r.write("Pricing.spec.md", pricingHead+"> **PRICX-R0002:** this branch\n\n## Overview\n")
	if !uncommitted {
		r.git("commit", "-am", "branch revision")
	}
	return r
}

// After the rebase both `R0002` sit in the spec, and the base's is cited by a line that
// existed where the branch forked. Only the branch's revision and the branch's citation move.
func TestRenumber_afterRebaseLeavesTheBaseCitation(t *testing.T) {
	t.Run("RNMBR-B04: Citations move only on the lines the branch added", func(t *testing.T) {})
	t.Run("RNMBR-I01: The base's revisions keep their meaning after a renumber", func(t *testing.T) {})
	r := newRepo(t)
	r.write("Pricing.spec.md", pricingHead+"> **PRICX-R0002:** the other PR\n\n## Overview\n")
	r.write("pricing_test.go", "package p\n// PRICX-R0002 is the other PR's\n")
	r.git("add", "anchors.yaml", "Pricing.spec.md", "pricing_test.go")
	r.git("commit", "-m", "base with the other PR merged")

	r.git("checkout", "-b", "feat")
	r.write("Pricing.spec.md", pricingHead+"> **PRICX-R0002:** the other PR\n> **PRICX-R0002:** this branch\n\n## Overview\n")
	r.write("pricing_test.go", "package p\n// PRICX-R0002 is the other PR's\n// PRICX-R0002 is this branch's\n")
	r.git("commit", "-am", "rebased branch revision")

	if _, err := r.renumber("--base", "main"); err != nil {
		t.Fatal(err)
	}
	spec := r.read("Pricing.spec.md")
	if !strings.Contains(spec, "> **PRICX-R0002:** the other PR\n> **PRICX-R0003:** this branch") {
		t.Errorf("only the branch's R0002 should become R0003:\n%s", spec)
	}
	test := r.read("pricing_test.go")
	if !strings.Contains(test, "// PRICX-R0002 is the other PR's\n// PRICX-R0003 is this branch's") {
		t.Errorf("the base's citation must stay and the branch's must move:\n%s", test)
	}
}

// The whole path through git: two branches each take `R0002`, and `anchors renumber` on
// the second moves only its own revision and its own citation.
func TestRenumber_movesOnlyTheBranchRevision(t *testing.T) {
	t.Run("RNMBR-B03: Only the branch's colliding revision moves to the next free number", func(t *testing.T) {})
	t.Run("RNMBR-B06: The dry run writes nothing", func(t *testing.T) {})
	t.Run("RNMBR-X01: The rewrite is left unstaged for review", func(t *testing.T) {})
	r := newRepo(t)
	r.write("Pricing.spec.md", pricingHead+"\n## Overview\n")
	r.write("pricing_test.go", "package p\n// PRICX-R0001 is proven below\n")
	r.write("untouched.go", "package p\n// PRICX-R0002 per the base\n")
	r.git("add", "anchors.yaml", "Pricing.spec.md", "pricing_test.go", "untouched.go")
	r.git("commit", "-m", "base")

	r.git("checkout", "-b", "feat")
	r.write("Pricing.spec.md", pricingHead+"> **PRICX-R0002:** this branch\n\n## Overview\n")
	r.write("pricing_test.go", "package p\n// PRICX-R0001 is proven below\n// PRICX-R0002 is proven here\n")
	r.git("commit", "-am", "branch revision")

	r.git("checkout", "main")
	r.write("Pricing.spec.md", pricingHead+"\n## Overview\n\n> **PRICX-R0002:** the other PR\n")
	r.git("commit", "-am", "other PR")
	r.git("checkout", "feat")

	dry, err := r.renumber("--base", "main", "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dry, "PRICX-R0002 → PRICX-R0003") || !strings.Contains(dry, "(dry-run — nothing was written.)") {
		t.Errorf("the dry run must show the move:\n%s", dry)
	}
	if strings.Contains(r.read("Pricing.spec.md"), "R0003") {
		t.Fatal("--dry-run wrote to disk")
	}

	if _, err := r.renumber("--base", "main"); err != nil {
		t.Fatal(err)
	}
	spec := r.read("Pricing.spec.md")
	if !strings.Contains(spec, "> **PRICX-R0003:** this branch") || strings.Contains(spec, "PRICX-R0002") {
		t.Errorf("the branch's R0002 should now be R0003:\n%s", spec)
	}
	test := r.read("pricing_test.go")
	if !strings.Contains(test, "// PRICX-R0003 is proven here") || !strings.Contains(test, "// PRICX-R0001 is proven below") {
		t.Errorf("only the branch's citation should move:\n%s", test)
	}
	if got := r.read("untouched.go"); !strings.Contains(got, "PRICX-R0002 per the base") {
		t.Errorf("a file the branch did not change was rewritten:\n%s", got)
	}
	if staged, err := gitOut(r.dir, "diff", "--cached", "--name-only"); err != nil || staged != "" {
		t.Errorf("the rewrite must be left unstaged, got staged %q (%v)", staged, err)
	}
	if dirty, _ := gitOut(r.dir, "diff", "--name-only"); !strings.Contains(dirty, "Pricing.spec.md") {
		t.Errorf("the rewrite must be in the working tree: %q", dirty)
	}
}

// With no base given, the remote copy of the integration branch wins over the local one.
func TestRenumber_defaultBase(t *testing.T) {
	t.Run("RNMBR-B01: The default base is the remote integration branch, else the local one", func(t *testing.T) {})
	r := forkedRepo(t, false)
	out, err := r.renumber("--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "renumber against main:") {
		t.Errorf("with no remote, the base is the local branch:\n%s", out)
	}
	// A remote copy of main at the fork point: nothing of the other PR is there yet.
	r.git("update-ref", "refs/remotes/origin/main", "main~1")
	out, err = r.renumber("--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "collides with origin/main") {
		t.Errorf("with a remote copy, the base is origin/main:\n%s", out)
	}
}

// Uncommitted and untracked specs are part of what the branch changed.
func TestRenumber_examinesUncommittedAndUntrackedSpecs(t *testing.T) {
	t.Run("RNMBR-B02: Specs changed but not committed, and untracked specs, are examined", func(t *testing.T) {})
	r := forkedRepo(t, true)
	out, err := r.renumber("--base", "main", "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "PRICX-R0002 → PRICX-R0003   (Pricing.spec.md)") {
		t.Errorf("an uncommitted revision was not examined:\n%s", out)
	}

	// An untracked spec whose file the base created with the same revision number.
	u := newRepo(t)
	u.write("Tax.spec.md", "# Tax\n")
	u.git("add", "anchors.yaml", "Tax.spec.md")
	u.git("commit", "-m", "base")
	u.git("checkout", "-b", "feat")
	u.git("checkout", "main")
	u.write("Other.spec.md", "<!-- @anchors\n  code: TAXXX\n-->\n# Other\n\n> **TAXXX-R0001:** the base's\n")
	u.git("add", "Other.spec.md")
	u.git("commit", "-m", "other")
	u.git("checkout", "feat")
	u.write("Other.spec.md", "<!-- @anchors\n  code: TAXXX\n-->\n# Other\n\n> **TAXXX-R0001:** this branch\n")
	out, err = u.renumber("--base", "main", "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "TAXXX-R0001 → TAXXX-R0002   (Other.spec.md)") {
		t.Errorf("an untracked spec was not examined:\n%s", out)
	}
}

func TestRenumber_binaryModeAndNothingToDo(t *testing.T) {
	t.Run("RNMBR-B05: A changed binary file is not rewritten", func(t *testing.T) {})
	t.Run("RNMBR-B08: A rewritten file keeps its permission bits", func(t *testing.T) {})
	t.Run("RNMBR-B07: No collision says there is nothing to renumber", func(t *testing.T) {})
	testkit.SkipWithoutPOSIXPermissions(t)
	r := forkedRepo(t, true)
	bin := "PRICX-R0002\x00binary"
	r.write("blob.bin", bin)
	r.write("run.sh", "#!/bin/sh\n# PRICX-R0002 in a script\n")
	if err := os.Chmod(filepath.Join(r.dir, "run.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := r.renumber("--base", "main"); err != nil {
		t.Fatal(err)
	}
	if got := r.read("blob.bin"); got != bin {
		t.Errorf("a binary file was rewritten: %q", got)
	}
	if got := r.read("run.sh"); !strings.Contains(got, "PRICX-R0003 in a script") {
		t.Fatalf("the script's citation was not rewritten:\n%s", got)
	}
	if info, err := os.Stat(filepath.Join(r.dir, "run.sh")); err != nil || info.Mode().Perm() != 0o755 {
		t.Errorf("the rewrite lost the file mode: %v %v", info.Mode(), err)
	}

	// Run again: the branch's revision is R0003 now, and nothing collides.
	out, err := r.renumber("--base", "main")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "✓ no revision this branch added collides with main — nothing to renumber.") {
		t.Errorf("a second run must have nothing to do:\n%s", out)
	}
}

func TestRenumber_onlyTheGivenFiles(t *testing.T) {
	t.Run("RNMBR-B09: With files given, only those specs are examined", func(t *testing.T) {})
	r := forkedRepo(t, false)
	other := "<!-- @anchors\n  code: TAXXX\n-->\n# Tax\n\n"
	r.git("checkout", "main")
	r.write("Tax.spec.md", other+"> **TAXXX-R0001:** main's\n")
	r.git("add", "Tax.spec.md")
	r.git("commit", "-m", "tax on main")
	r.git("checkout", "feat")
	r.write("Tax.spec.md", other+"> **TAXXX-R0001:** this branch\n")

	out, err := r.renumber("--base", "main", "--dry-run", filepath.Join(r.dir, "Tax.spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "TAXXX-R0001 → TAXXX-R0002") || strings.Contains(out, "PRICX-R0002 →") {
		t.Errorf("only the given spec must be examined:\n%s", out)
	}
}

func TestRenumber_refusals(t *testing.T) {
	t.Run("RNMBR-E01: A project with no configuration is refused", func(t *testing.T) {})
	t.Run("RNMBR-E02: A base with no merge base is refused naming it", func(t *testing.T) {})
	t.Run("RNMBR-E03: An unreadable spec given by hand is refused naming it", func(t *testing.T) {})
	r := forkedRepo(t, false)
	if _, err := r.renumber("--base", "no-such-branch"); err == nil || !strings.Contains(err.Error(), "no merge base between no-such-branch and HEAD") {
		t.Errorf("an unknown base: got %v", err)
	}
	before := r.read("Pricing.spec.md")
	if _, err := r.renumber("--base", "main", filepath.Join(r.dir, "Missing.spec.md"), filepath.Join(r.dir, "Pricing.spec.md")); err == nil || !strings.Contains(err.Error(), "read Missing.spec.md") {
		t.Errorf("a missing spec given by hand: got %v", err)
	}
	if r.read("Pricing.spec.md") != before {
		t.Error("a refused renumber wrote to disk")
	}
	if _, err := runCmdErr(newRenumberCmd(), t, "--root", t.TempDir()); err == nil || !strings.Contains(err.Error(), "anchors.yaml") {
		t.Errorf("no config: got %v", err)
	}
}

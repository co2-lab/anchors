package mapcmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

// After the rebase both `R0002` sit in the spec, and the base's is cited by a line that
// existed where the branch forked. Only the branch's revision and the branch's citation move.
func TestRenumber_afterRebaseLeavesTheBaseCitation(t *testing.T) {
	r := newRepo(t)
	head := "<!-- @anchors\n  code: PRICX\n-->\n# Pricing\n\n> **PRICX-R0001:** first\n"
	r.write("Pricing.spec.md", head+"> **PRICX-R0002:** the other PR\n\n## Overview\n")
	r.write("pricing_test.go", "package p\n// PRICX-R0002 is the other PR's\n")
	r.git("add", "anchors.yaml", "Pricing.spec.md", "pricing_test.go")
	r.git("commit", "-m", "base with the other PR merged")

	r.git("checkout", "-b", "feat")
	r.write("Pricing.spec.md", head+"> **PRICX-R0002:** the other PR\n> **PRICX-R0002:** this branch\n\n## Overview\n")
	r.write("pricing_test.go", "package p\n// PRICX-R0002 is the other PR's\n// PRICX-R0002 is this branch's\n")
	r.git("commit", "-am", "rebased branch revision")

	cmd := newRenumberCmd()
	cmd.SetArgs([]string{"--root", r.dir, "--base", "main"})
	if err := cmd.Execute(); err != nil {
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
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	read := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	head := "<!-- @anchors\n  code: PRICX\n-->\n# Pricing\n\n> **PRICX-R0001:** first\n"
	git("init", "-b", "main")
	git("config", "user.email", "t@t")
	git("config", "user.name", "t")
	write("anchors.yaml", "version: 2\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	write("Pricing.spec.md", head+"\n## Overview\n")
	write("pricing_test.go", "package p\n// PRICX-R0001 is proven below\n")
	git("add", "anchors.yaml", "Pricing.spec.md", "pricing_test.go")
	git("commit", "-m", "base")

	git("checkout", "-b", "feat")
	write("Pricing.spec.md", head+"> **PRICX-R0002:** this branch\n\n## Overview\n")
	write("pricing_test.go", "package p\n// PRICX-R0001 is proven below\n// PRICX-R0002 is proven here\n")
	git("commit", "-am", "branch revision")

	git("checkout", "main")
	write("Pricing.spec.md", head+"\n## Overview\n\n> **PRICX-R0002:** the other PR\n")
	git("commit", "-am", "other PR")
	git("checkout", "feat")

	dry := newRenumberCmd()
	dry.SetArgs([]string{"--root", dir, "--base", "main", "--dry-run"})
	if err := dry.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(read("Pricing.spec.md"), "R0003") {
		t.Fatal("--dry-run wrote to disk")
	}

	cmd := newRenumberCmd()
	cmd.SetArgs([]string{"--root", dir, "--base", "main"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	spec := read("Pricing.spec.md")
	if !strings.Contains(spec, "> **PRICX-R0003:** this branch") || strings.Contains(spec, "PRICX-R0002") {
		t.Errorf("the branch's R0002 should now be R0003:\n%s", spec)
	}
	test := read("pricing_test.go")
	if !strings.Contains(test, "// PRICX-R0003 is proven here") || !strings.Contains(test, "// PRICX-R0001 is proven below") {
		t.Errorf("only the branch's citation should move:\n%s", test)
	}
}

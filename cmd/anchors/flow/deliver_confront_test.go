// @anchors
//   ref: DLCND

package flow

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The worst silence the confrontation had: without git it returned `nil`, and `nil` is the
// same value as "I checked and everything is fine". Whoever read an output with no warning
// concluded the declared files checked out — a claim nobody had verified.
func TestConfront_saysWhenItCouldNotLook(t *testing.T) {
	t.Run("DLCND-B02: Without git the output says the confrontation did not happen", func(t *testing.T) {})
	t.Run("DLCND-I01: Not looking and finding nothing are different answers", func(t *testing.T) {})
	dir := t.TempDir()
	if _, _, ok := repoAbove(dir); ok {
		t.Skipf("the temporary directory %s is inside a git repository", dir)
	}
	declared := []string{filepath.Join(dir, "a.go")}
	if err := os.WriteFile(declared[0], []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	warnings, confronted := untouchedFiles(dir, declared)
	if confronted {
		t.Fatal("without a repository there is no way to confront — saying it confronted is a lie")
	}
	if len(warnings) != 0 {
		t.Errorf("without a confrontation there is no finding to report: %v", warnings)
	}
	out := stdoutOf(t, func() { confrontDelivery(dir, declared, "") })
	if !strings.Contains(out, "could not confront the declared files") || !strings.Contains(out, "nobody verified") ||
		strings.Contains(out, "do NOT appear in the diff") {
		t.Errorf("outside git the output must say the confrontation did not happen:\n%s", out)
	}

	// The other half: in a repository where every declared file changed, it DID look.
	root := t.TempDir()
	writeFile(t, root, "a.go", "package a\n")
	gitRepo(t, root)
	writeFile(t, root, "a.go", "package a // changed\n")
	warnings, confronted = untouchedFiles(root, []string{"a.go"})
	if !confronted || len(warnings) != 0 {
		t.Errorf("inside git with every file touched: confronted=%v warnings=%v", confronted, warnings)
	}
}

// The counterpart: in a real repository the confrontation HAPPENS, and a declared file
// that nobody touched is still accused.
func TestConfront_accusesADeclaredAndUntouchedFile(t *testing.T) {
	t.Run("DLCND-B03: A committed and untouched declared file is accused", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "quiet.go", "package a\n")
	gitRepo(t, root)
	quiet := filepath.Join(root, "quiet.go")

	warnings, confronted := untouchedFiles(root, []string{quiet})
	if !confronted {
		t.Fatal("with a repository the confrontation must happen")
	}
	if len(warnings) != 1 || warnings[0] != "quiet.go" {
		t.Fatalf("a declared and untouched file should be accused by its relative path, got %v", warnings)
	}
	out := stdoutOf(t, func() { confrontDelivery(root, []string{quiet}, "") })
	section := out[strings.Index(out, "declared files that do NOT appear in the diff"):]
	if !strings.Contains(section, "quiet.go") || !strings.Contains(section, "the disk does not confirm") {
		t.Errorf("the untouched file must be listed with the warning:\n%s", out)
	}
}

// `git status --porcelain` collapses a new directory into one line (`?? dir/`); comparing
// exact paths blinded the confrontation on every delivery that created a directory.
func TestConfront_modifiedAndNewDirectoryFilesAreTouched(t *testing.T) {
	t.Run("DLCND-B04: Modified files and files in a new directory are not accused", func(t *testing.T) {})
	t.Run("DLCND-X01: The confrontation writes nothing", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	writeFile(t, root, "src/quiet.ts", "export const q = 1\n")
	gitRepo(t, root)
	writeFile(t, root, "src/pricing.ts", "export const p = 2\n")
	writeFile(t, root, "src/newdir/handler.ts", "export const h = 1\n")

	status := func() string {
		c := exec.Command("git", "-C", root, "status", "--porcelain", "--untracked-files=all")
		out, err := c.Output()
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	before := status()
	warnings, confronted := untouchedFiles(root, []string{"src/pricing.ts", "src/newdir/handler.ts"})
	if !confronted || len(warnings) != 0 {
		t.Errorf("a modified file and a file in a new directory were touched: confronted=%v %v", confronted, warnings)
	}
	stdoutOf(t, func() { confrontDelivery(root, []string{"src/pricing.ts", "src/quiet.ts"}, "src/pricing.ts") })
	if after := status(); after != before {
		t.Errorf("the confrontation changed the working tree:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// A renamed file is touched, and so is a file under a `--root` that is a subdirectory of
// the repository: git names paths from the top of the repository, and the declared files
// are relative to the root. Before, `R  old -> new` was read as one path and every file
// under a subdirectory root was accused.
func TestConfront_renamedFilesAndSubdirectoryRoot(t *testing.T) {
	t.Run("DLCND-B08: A renamed file and a file under a subdirectory root are not accused", func(t *testing.T) {})
	repo := t.TempDir()
	writeFile(t, repo, "app/src/old.ts", "export const o = 1\n")
	writeFile(t, repo, "app/src/pricing.ts", "export const p = 1\n")
	writeFile(t, repo, "app/src/quiet.ts", "export const q = 1\n")
	writeFile(t, repo, "other/pricing.ts", "export const p = 1\n")
	gitRepo(t, repo)
	if out, err := exec.Command("git", "-C", repo, "mv", "app/src/old.ts", "app/src/renamed.ts").CombinedOutput(); err != nil {
		t.Fatalf("git mv: %s", out)
	}
	writeFile(t, repo, "app/src/pricing.ts", "export const p = 2\n")
	writeFile(t, repo, "other/pricing.ts", "export const p = 2\n")

	// from the top of the repository, the renamed file is touched
	warnings, confronted := untouchedFiles(repo, []string{"app/src/renamed.ts", "app/src/pricing.ts"})
	if !confronted || len(warnings) != 0 {
		t.Errorf("a renamed and a modified file are touched: confronted=%v %v", confronted, warnings)
	}
	// with the root a subdirectory, paths are relative to it
	root := filepath.Join(repo, "app")
	warnings, confronted = untouchedFiles(root, []string{"src/renamed.ts", "src/pricing.ts", "src/quiet.ts"})
	if !confronted || strings.Join(warnings, ",") != "src/quiet.ts" {
		t.Errorf("under a subdirectory root only the untouched file is accused, got confronted=%v %v", confronted, warnings)
	}
	// a change outside the root does not clear a file of the same name inside it
	if w, _ := untouchedFiles(filepath.Join(repo, "other"), []string{"pricing.ts"}); len(w) != 0 {
		t.Errorf("other/pricing.ts is modified: %v", w)
	}
}

// The confrontation prints and never refuses: blocking would push authors to declare less.
func TestConfront_neverBlocksTheDelivery(t *testing.T) {
	t.Run("DLCND-B01: The confrontation never blocks the delivery", func(t *testing.T) {})
	t.Run("DLCND-B06: A failing informative gate is listed by its first sentence", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", deliverLocalYAML)
	writeFile(t, root, "anchors.graph.yaml", deliverMap)
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	writeFile(t, root, "src/quiet.ts", "export const q = 1\n")
	gitRepo(t, root)

	out, err := runDeliver(t, "--root", root, "--stage", "code", "--unit", "src/pricing.ts",
		"--file", "src/pricing.ts,src/quiet.ts", "--intent", "implements PRICX-B01", "--date", "2026-09-26")
	if err != nil {
		t.Fatalf("warnings must not fail the delivery: %v", err)
	}
	if recs, _ := filepath.Glob(filepath.Join(root, "changes", "*.md")); len(recs) != 1 {
		t.Errorf("the record must exist despite the warnings: %v", recs)
	}
	if !strings.Contains(out, "do NOT appear in the diff") {
		t.Errorf("the untouched files are still warned:\n%s", out)
	}
	if !strings.Contains(out, "always-red @ src/pricing.ts — the loop drops the last page.") ||
		strings.Contains(out, "Details follow here") {
		t.Errorf("the failing informative gate is shown by its first sentence:\n%s", out)
	}
}

func TestRedGates_nothingToRunListsNothing(t *testing.T) {
	t.Run("DLCND-B07: Without gates, map or delivered node no gate is listed", func(t *testing.T) {})
	noGates := t.TempDir()
	writeFile(t, noGates, "anchors.yaml", "version: 1\nlayers:\n  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n")
	writeFile(t, noGates, "anchors.graph.yaml", deliverMap)
	if got := redGates(noGates, []string{"src/pricing.ts"}, "src/pricing.ts"); got != nil {
		t.Errorf("without gates nothing runs: %v", got)
	}
	noMap := t.TempDir()
	writeFile(t, noMap, "anchors.yaml", deliverLocalYAML)
	if got := redGates(noMap, []string{"src/pricing.ts"}, "src/pricing.ts"); got != nil {
		t.Errorf("without a map nothing runs: %v", got)
	}
	noNode := t.TempDir()
	writeFile(t, noNode, "anchors.yaml", deliverLocalYAML)
	writeFile(t, noNode, "anchors.graph.yaml", deliverMap)
	writeFile(t, noNode, "src/pricing.ts", "export const p = 1\n")
	if got := redGates(noNode, []string{"src/other.ts"}, "src/other.ts"); got != nil {
		t.Errorf("without a delivered node nothing runs: %v", got)
	}
	// And the same project with the delivered node does list the red gate — the three
	// empty answers above are not the gate being silent.
	if got := redGates(noNode, []string{"src/pricing.ts"}, "src/pricing.ts"); len(got) != 1 {
		t.Errorf("the delivered node must be run through the gate: %v", got)
	}
}

// repoAbove is the test's local helper — it avoids depending on `git status` to know
// whether the TempDir landed inside a repository.
func repoAbove(root string) (string, string, bool) {
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, "", true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", false
		}
		dir = parent
	}
}

// A mutation signal already ingested silences the warning; no test, no warning either.
func TestMutationNotMeasured(t *testing.T) {
	t.Run("DLCND-B05: A tested unit without a mutation signal is warned", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "src/pricing.ts", "x\n")
	if got := mutationNotMeasured(root, "src/pricing.ts"); got != "" {
		t.Errorf("without a test the subject is another: %q", got)
	}
	writeFile(t, root, "src/pricing.test.ts", "x\n")
	writeFile(t, root, "anchors.graph.yaml", "version: 6\nnodes:\n"+
		"  - id: src/pricing.ts\n    kind: code\n    rev: a\n    signal:\n      mutants_killed: 3\nedges: []\n")
	if got := mutationNotMeasured(root, "src/pricing.ts"); got != "" {
		t.Errorf("a measured unit must not be warned: %q", got)
	}
	if got := mutationNotMeasured(root, ""); got != "" {
		t.Errorf("no unit, no warning: %q", got)
	}
}

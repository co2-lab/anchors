// @anchors
//   ref: RCPLR

package recode

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
)

// writeProject writes a minimal synthetic project to exercise BuildPlan.
func writeProject(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// testCfg declares a layer for .tsx/.spec.md/.yaml and the recode block.
func testCfg() *config.Config {
	return &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec": {Pattern: "**/*.spec.md", Kind: "spec"},
			"code": {Pattern: "**/*.tsx", Kind: "code"},
			"e2e":  {Pattern: "**/*.yaml", Kind: "test"},
		},
		Recode: &config.Recode{
			TestID:       "lower",
			FilePatterns: []string{"**/{{code}}-*.yaml"},
		},
	}
}

// plainCfg is testCfg without the recode dialect.
func plainCfg() *config.Config {
	c := testCfg()
	c.Recode = nil
	return c
}

func TestBuildPlan_refusesMalformedOrEqualCodes(t *testing.T) {
	t.Run("RCPLR-B01: A malformed source or target code is refused", func(t *testing.T) {})
	t.Run("RCPLR-B02: Renaming a code to itself is refused", func(t *testing.T) {})
	root := writeProject(t, map[string]string{"Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n"})
	for _, c := range [][2]string{{"abcdx", "WXYZX"}, {"ABCDX", "WX"}, {"ABC-X", "WXYZX"}} {
		if _, err := BuildPlan(root, plainCfg(), c[0], c[1]); err == nil || !strings.Contains(err.Error(), "invalid") {
			t.Errorf("BuildPlan(%s → %s) = %v, want a refusal of the malformed code", c[0], c[1], err)
		}
	}
	if _, err := BuildPlan(root, plainCfg(), "ABCDX", "ABCDX"); err == nil || !strings.Contains(err.Error(), "same code") {
		t.Errorf("BuildPlan(ABCDX → ABCDX) = %v, want the same-code refusal", err)
	}
}

func TestBuildPlan_refusesACodeThatAppearsNowhere(t *testing.T) {
	t.Run("RCPLR-B07: A code that appears nowhere is refused", func(t *testing.T) {})
	root := writeProject(t, map[string]string{"Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n"})
	_, err := BuildPlan(root, testCfg(), "QQQQX", "WXYZX")
	if err == nil || !strings.Contains(err.Error(), "does not appear in any file") {
		t.Fatalf("BuildPlan of an absent code = %v, want the not-found refusal", err)
	}
}

func TestBuildPlan_plansEveryFileSortedAndWritesNothing(t *testing.T) {
	t.Run("RCPLR-B03: The plan holds every file with the code, sorted, with its new content", func(t *testing.T) {})
	t.Run("RCPLR-X02: Planning writes nothing", func(t *testing.T) {})
	files := map[string]string{
		"z/Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-B01 x\n",
		"a/Foo.tsx":     "// @anchors\n//   ref: ABCDX\n",
		"m/Other.tsx":   "// @anchors\n//   ref: MNOPX\n",
	}
	root := writeProject(t, files)
	plan, err := BuildPlan(root, plainCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, f := range plan.Files {
		paths = append(paths, f.Path)
		if strings.Contains(f.NewContent, "ABCDX") || !strings.Contains(f.NewContent, "WXYZX") {
			t.Errorf("%s: new content not rewritten: %q", f.Path, f.NewContent)
		}
	}
	if want := []string{"a/Foo.tsx", "z/Foo.spec.md"}; !reflect.DeepEqual(paths, want) {
		t.Errorf("plan files = %v, want %v (only the files with the code, sorted)", paths, want)
	}
	if plan.Total != 3 {
		t.Errorf("plan total = %d, want 3 content replacements", plan.Total)
	}
	for rel, content := range files {
		if b, _ := os.ReadFile(filepath.Join(root, rel)); string(b) != content {
			t.Errorf("BuildPlan wrote %s: %q", rel, b)
		}
	}
}

func TestBuildPlan_dialect_conformingProject(t *testing.T) {
	t.Run("RCPLR-B04: The declared testID prefix is rewritten and counted apart", func(t *testing.T) {})
	t.Run("RCPLR-B05: Files named after the code are planned for renaming", func(t *testing.T) {})
	// a project that FOLLOWS the convention: testID = the code in lower case (abcdx-).
	root := writeProject(t, map[string]string{
		"Foo.spec.md":          "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-S01 x\n",
		"Foo.tsx":              "// @anchors\n//   ref: ABCDX\nconst x = <View testID=\"abcdx-root\" />\n",
		"flows/ABCDX-S01.yaml": "tags:\n  - ABCDX-S01\n",
	})
	plan, err := BuildPlan(root, testCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	if plan.TestIDs != 1 {
		t.Errorf("want 1 testID rewritten (abcdx-→wxyzx-), got %d", plan.TestIDs)
	}
	if len(plan.Renames) != 1 || plan.Renames[0].To != "flows/WXYZX-S01.yaml" {
		t.Errorf("want the rename ABCDX-S01.yaml→WXYZX-S01.yaml, got %+v", plan.Renames)
	}
	if plan.TestIDLegacy != "" {
		t.Errorf("a conforming project should NOT carry the legacy warning: %q", plan.TestIDLegacy)
	}
	// Without the dialect, neither the testID nor the file name is touched.
	plain, err := BuildPlan(root, plainCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	if plain.TestIDs != 0 || len(plain.Renames) != 0 {
		t.Errorf("without the recode block: %d testIDs, %v renames; want none", plain.TestIDs, plain.Renames)
	}
}

func TestBuildPlan_dialect_legacyWarns(t *testing.T) {
	t.Run("RCPLR-B06: A divergent testID prefix is warned about and never rewritten", func(t *testing.T) {})
	// a LEGACY project: code ABCDX but the testID uses a divergent prefix (zzzz-).
	root := writeProject(t, map[string]string{
		"Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-S01 x\n",
		"Foo.tsx":     "// @anchors\n//   ref: ABCDX\nconst x = <View testID=\"zzzz-root\" />\n",
	})
	plan, err := BuildPlan(root, testCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	if plan.TestIDs != 0 {
		t.Errorf("expected prefix absent → 0 rewrites, got %d", plan.TestIDs)
	}
	if plan.TestIDLegacy == "" {
		t.Error("a legacy project should carry the divergent-testID warning")
	}
	for _, f := range plan.Files {
		if f.Path == "Foo.tsx" && !strings.Contains(f.NewContent, `testID="zzzz-root"`) {
			t.Errorf("the divergent testID must not be rewritten: %q", f.NewContent)
		}
	}
}

func TestApply_keepsModesAndRenames(t *testing.T) {
	t.Run("RCPLR-B08: Applying writes each file with its mode and performs the renames", func(t *testing.T) {})
	testkit.SkipWithoutPOSIXPermissions(t)
	root := writeProject(t, map[string]string{
		"Foo.spec.md":          "<!-- @anchors\n  code: ABCDX\n-->\n",
		"run.tsx":              "// @anchors\n//   ref: ABCDX\n",
		"flows/ABCDX-S01.yaml": "tags: []\n",
	})
	if gitmeta.Check(root) == gitmeta.Disponível {
		t.Skipf("the temporary folder %s is inside a git repository", root)
	}
	if err := os.Chmod(filepath.Join(root, "run.tsx"), 0o755); err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPlan(root, testCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	n, err := plan.Apply(root)
	if err != nil || n != 3 {
		t.Fatalf("Apply = %d, %v; want 3 (two files written, one renamed)", n, err)
	}
	info, err := os.Stat(filepath.Join(root, "run.tsx"))
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Errorf("run.tsx mode = %v (%v), want 0755 kept", info.Mode().Perm(), err)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "run.tsx")); !strings.Contains(string(b), "WXYZX") {
		t.Errorf("run.tsx not rewritten: %q", b)
	}
	if _, err := os.Stat(filepath.Join(root, "flows", "WXYZX-S01.yaml")); err != nil {
		t.Errorf("the flow file was not renamed: %v", err)
	}
}

func TestApply_reportsAFileItCannotWrite(t *testing.T) {
	t.Run("RCPLR-E05: A file that cannot be written stops the apply, naming it", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "adir"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &Plan{Files: []FileChange{{Path: "adir", NewContent: "x"}}}
	n, err := p.Apply(root)
	if err == nil || !strings.Contains(err.Error(), "writing adir") || n != 0 {
		t.Fatalf("Apply onto a folder = %d, %v; want 0 and an error naming adir", n, err)
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"},
	} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("setup %v: %s", args, out)
		}
	}
	return dir
}

// OUTSIDE a repository the rename is the RIGHT behaviour: there is no history to keep, and
// demanding git there would turn a legitimate degradation into a block.
func TestMoveWithoutRepoUsesAPlainRename(t *testing.T) {
	t.Run("RCPLR-B09: Moves go through git inside a repository and are plain renames outside", func(t *testing.T) {})
	dir := t.TempDir()
	if gitmeta.Check(dir) == gitmeta.Disponível {
		t.Skipf("the temporary folder %s is inside a git repository", dir)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gitMoves(dir, "a.go", "sub/b.go"); err != nil {
		t.Fatalf("without a repository the rename should work: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "sub", "b.go")); err != nil {
		t.Errorf("the file did not reach the destination: %v", err)
	}
}

// In a repo, a tracked file moves via `git mv` — and the index records the rename.
func TestMoveInRepoUsesGitMv(t *testing.T) {
	t.Run("RCPLR-B09: Moves go through git inside a repository and are plain renames outside", func(t *testing.T) {})
	dir := initRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-m", "base"}} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("setup %v: %s", args, out)
		}
	}

	if err := gitMoves(dir, "a.go", "sub/b.go"); err != nil {
		t.Fatalf("git mv should work: %v", err)
	}

	c := exec.Command("git", "status", "--porcelain")
	c.Dir = dir
	out, _ := c.Output()
	// `R` (a staged rename) is what proves the history was kept; a rename behind git's back
	// would show as `D` + `??`.
	if !strings.Contains(string(out), "R ") {
		t.Errorf("the rename was not recorded in the index: %s", out)
	}
}

// A file NOT YET TRACKED is not a refusal to respect: git never knew it, so there is no
// history to keep nor index to keep consistent. The direct rename does what the user asked.
//
// This is the case of a freshly initialised project — nothing committed yet — which is exactly
// where the recode is most likely to run.
func TestMoveOfAnUntrackedFileHappens(t *testing.T) {
	t.Run("RCPLR-B10: An untracked file inside a repository is moved by a plain rename", func(t *testing.T) {})
	dir := initRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "loose.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gitMoves(dir, "loose.go", "sub/loose.go"); err != nil {
		t.Fatalf("an untracked file should be moved directly: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "sub", "loose.go")); err != nil {
		t.Errorf("the file did not reach the destination: %v", err)
	}
}

// The bug `-k` hid: `git mv -k` silently SKIPS what it cannot move and exits 0. The rename did
// not happen, `err` was nil, and Apply counted it as done — a project with nothing committed
// would see "✓ N files rewritten" with the N still in place. This test fixes that "success"
// means "the file moved".
func TestMoveThatReportsSuccessReallyMoved(t *testing.T) {
	t.Run("RCPLR-B10: An untracked file inside a repository is moved by a plain rename", func(t *testing.T) {})
	dir := initRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "loose.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gitMoves(dir, "loose.go", "moved.go"); err != nil {
		t.Fatalf("gitMove: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "loose.go")); err == nil {
		t.Error("gitMove said it moved, but the source is still there — the silence of `-k`")
	}
	if _, err := os.Stat(filepath.Join(dir, "moved.go")); err != nil {
		t.Errorf("gitMove said it moved, but the destination does not exist: %v", err)
	}
}

// A REAL git refusal (here: the destination exists and is tracked) is not bypassed by
// os.Rename. Git refuses for a reason, and moving anyway leaves the index diverging from the
// disk — silently, in a command that works in mass.
func TestMoveDoesNotBypassAGitRefusal(t *testing.T) {
	t.Run("RCPLR-X01: A git refusal is surfaced and never bypassed", func(t *testing.T) {})
	dir := initRepo(t)
	for _, f := range []string{"a.go", "b.go"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("package "+strings.TrimSuffix(f, ".go")+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-m", "base"}} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("setup %v: %s", args, out)
		}
	}

	// The destination exists and is tracked: `git mv` refuses without `-f`.
	err := gitMoves(dir, "a.go", "b.go")

	if err == nil {
		t.Fatal("git refused the move — bypassing it with os.Rename would overwrite a tracked file")
	}
	if !strings.Contains(err.Error(), "refused") {
		t.Errorf("the message must say it was git that refused: %v", err)
	}
	// Nothing was moved behind git's back: the source is there, the destination untouched.
	if _, statErr := os.Stat(filepath.Join(dir, "a.go")); statErr != nil {
		t.Error("the source vanished in an operation that failed")
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "b.go")); string(b) != "package b\n" {
		t.Errorf("the destination was overwritten: %q", b)
	}
}

// The collision check compared the bare target code against the file's scenario codes
// (`WXYZX-B01`), which never equal a bare code: renaming onto a code another unit owns
// was planned, merging two identities.
func TestBuildPlan_refusesATargetOwnedByAnotherUnit(t *testing.T) {
	t.Run("RCPLR-B11: A target code another unit already owns is refused", func(t *testing.T) {})
	for name, other := range map[string]string{
		"declared in the header":  "<!-- @anchors\n  code: WXYZX\n-->\n# Bar\n",
		"carried by its scenario": "# Bar\n\n### WXYZX-B01 x\n",
	} {
		t.Run(name, func(t *testing.T) {
			root := writeProject(t, map[string]string{
				"a/Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-B01 x\n",
				"b/Bar.spec.md": other,
			})
			_, err := BuildPlan(root, plainCfg(), "ABCDX", "WXYZX")
			if err == nil || !strings.Contains(err.Error(), "b/Bar.spec.md") {
				t.Fatalf("BuildPlan onto a code owned by b/Bar.spec.md = %v, want the collision refusal naming it", err)
			}
		})
	}
	// A code that merely CONTAINS the target is not a collision.
	root := writeProject(t, map[string]string{
		"a/Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-B01 x\n",
		"b/Bar.spec.md": "<!-- @anchors\n  code: WXYZY\n-->\n### WXYZY-B01 x\n",
	})
	if _, err := BuildPlan(root, plainCfg(), "ABCDX", "WXYZX"); err != nil {
		t.Errorf("an unrelated code was taken for a collision: %v", err)
	}
}

// The malformed-code refusal said "expected 4 chars" whatever the project declared; the
// lengths come from `code_lengths`.
func TestBuildPlan_malformedCodeNamesTheDeclaredLengths(t *testing.T) {
	t.Run("RCPLR-B12: The malformed-code refusal names the lengths the project declares", func(t *testing.T) {})
	prev := config.CodeLengths
	t.Cleanup(func() { config.SetCodeLengths(prev) })
	root := writeProject(t, map[string]string{"Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n"})
	for _, c := range []struct {
		lengths []int
		valid   string
		want    string
	}{{[]int{5}, "ABCDX", "expected 5 characters"}, {[]int{4, 6}, "ABCD", "expected 4/6 characters"}} {
		config.SetCodeLengths(c.lengths)
		_, err := BuildPlan(root, plainCfg(), "ABC", "WXYZX")
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("lengths %v: BuildPlan(ABC) = %v, want it to say %q", c.lengths, err, c.want)
		}
		_, err = BuildPlan(root, plainCfg(), c.valid, "WX")
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("lengths %v: BuildPlan(→ WX) = %v, want it to say %q", c.lengths, err, c.want)
		}
	}
}

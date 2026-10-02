package initx

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// writeFixture creates a toy project in dir: a colocated unit in mobile, code in the
// backend, guides, and noise to ignore (node_modules).
func writeFixture(t *testing.T, dir string) {
	t.Helper()
	files := map[string]string{
		// colocated unit (a screen in mobile)
		"apps/mobile/src/screens/Login.tsx":      "export const Login = () => null // LOGIX-A01",
		"apps/mobile/src/screens/Login.spec.md":  "> **Código**: `LOGIX`\n### LOGIX-A01: entrar",
		"apps/mobile/src/screens/Login.feature":  "@LOGIX-A01\nCenário: entrar",
		"apps/mobile/src/screens/Login.test.tsx": "it('LOGIX-A01: entra', () => {})",
		// more targets, to reach the colocation threshold (>=3)
		"apps/mobile/src/screens/Home.tsx":      "export const Home = () => null",
		"apps/mobile/src/screens/Home.spec.md":  "> **Código**: `HOMEX`",
		"apps/mobile/src/screens/Home.test.tsx": "it('HOMEX-S01', () => {})",
		"apps/mobile/src/screens/Prof.tsx":      "export const Prof = () => null",
		"apps/mobile/src/screens/Prof.spec.md":  "> **Código**: `PROF`",
		// backend (code only, no unit)
		"packages/backend/handlers/auth.ts": "export function auth() {}",
		"packages/backend/repos/user.ts":    "export function getUser() {}",
		// guides
		"guides/FRONTEND_GUIDE.md": "# Frontend",
		"guides/BACKEND_GUIDE.md":  "# Backend",
		// noise: must be ignored
		"node_modules/react/index.js": "module.exports = {}",
	}
	// more code in the backend and in mobile
	for i := range 12 {
		files["packages/backend/gen/f"+itoa(i)+".ts"] = "export const x = 1"
	}
	for i := range 12 {
		files["apps/mobile/src/extra/g"+itoa(i)+".ts"] = "export const y = 1"
	}
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func itoa(i int) string { return string(rune('0'+i/10)) + string(rune('0'+i%10)) }

func TestInfer(t *testing.T) {
	t.Run("INPRN-B02: The presence of specs, features and tests is detected", func(t *testing.T) {})
	t.Run("INPRN-B04: The first guides directory is detected with every guide file in it", func(t *testing.T) {})
	t.Run("INPRN-B07: Colocation is detected when at least three stems pair code with a spec, test or feature", func(t *testing.T) {})
	dir := t.TempDir()
	writeFixture(t, dir)

	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !p.HasSpecMD || !p.HasFeature || !p.HasTest {
		t.Errorf("should detect spec/feature/test; got spec=%v feature=%v test=%v",
			p.HasSpecMD, p.HasFeature, p.HasTest)
	}
	if !p.Colocated {
		t.Error("should detect colocation (there is a unit beside the code)")
	}
	if p.GuideDir != "guides" {
		t.Errorf("GuideDir = %q, want \"guides\"", p.GuideDir)
	}
	if len(p.GuideFiles) != 2 {
		t.Errorf("expected 2 guides, got %d (%v)", len(p.GuideFiles), p.GuideFiles)
	}
	if !contains(p.CodeDirs, "apps/mobile") || !contains(p.CodeDirs, "packages/backend") {
		t.Errorf("CodeDirs should contain apps/mobile and packages/backend; got %v", p.CodeDirs)
	}
	if !contains(p.CodeExts, ".ts") || !contains(p.CodeExts, ".tsx") {
		t.Errorf("CodeExts should contain .ts and .tsx; got %v", p.CodeExts)
	}
}

func TestInfer_ignoresNodeModules(t *testing.T) {
	t.Run("INPRN-B01: Dependency, build and tool directories are not walked", func(t *testing.T) {})
	dir := t.TempDir()
	writeFixture(t, dir)
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	// node_modules must not become a code dir
	for _, d := range p.CodeDirs {
		if d == "node_modules" || filepath.Base(d) == "react" {
			t.Errorf("node_modules should be ignored, but appeared in CodeDirs: %v", p.CodeDirs)
		}
	}
}

func TestInfer_buildsConfig(t *testing.T) {
	t.Run("INPRN-B09: Inference carries a proposed configuration with code layers and no artifact layer", func(t *testing.T) {})
	dir := t.TempDir()
	writeFixture(t, dir)
	p, _ := Infer(dir)
	cfg := p.Config

	// ARTIFACT layers are NOT pre-created by inference — they come from the user's choice
	// at init (ApplyArtifactChoice), which is always asked.
	for _, name := range []string{"spec", "feature", "test", "guide"} {
		if _, ok := cfg.Layers[name]; ok {
			t.Errorf("inference should not pre-create the artifact layer %q (it is the user's choice)", name)
		}
	}
	// but the DETECTED artifacts are reported (to pre-check the options)
	det := p.DetectedArtifacts()
	for _, name := range []string{"spec", "feature", "test", "guide"} {
		if !det[name] {
			t.Errorf("DetectedArtifacts should contain %q (it is in the fixture)", name)
		}
	}
	// there is at least one code layer (inference proposes those as default)
	if len(CodeLayerNames(cfg)) == 0 {
		t.Error("the proposed config should have code layers")
	}
	// governs starts empty (filled in the Q&A)
	if len(cfg.Governs) != 0 {
		t.Errorf("governs should start empty, got %+v", cfg.Governs)
	}
}

func contains(s []string, v string) bool {
	return slices.Contains(s, v)
}

// writeTree writes each file (relative path → content) under dir.
func writeTree(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// codeFiles returns n code files named <dir>/f<i><ext>, each with content.
func codeFiles(dir, ext, content string, n int) map[string]string {
	out := map[string]string{}
	for i := range n {
		out[fmt.Sprintf("%s/f%03d%s", dir, i, ext)] = content
	}
	return out
}

func TestInfer_ignoresEveryToolDir(t *testing.T) {
	t.Run("INPRN-B01: Dependency, build and tool directories are not walked", func(t *testing.T) {})
	dir := t.TempDir()
	for _, ignored := range []string{"node_modules", ".git", "dist", "build", "vendor", ".next", "coverage", ".expo", ".anchors"} {
		files := codeFiles(ignored+"/pkg", ".ts", "", 12)
		files[ignored+"/x/a.spec.md"] = "spec"
		writeTree(t, dir, files)
	}
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.CodeDirs) != 0 || len(p.CodeExts) != 0 || p.HasSpecMD {
		t.Errorf("nothing under an ignored dir should be seen, got dirs=%v exts=%v spec=%v",
			p.CodeDirs, p.CodeExts, p.HasSpecMD)
	}
}

func TestInfer_planWinsOverGuide(t *testing.T) {
	t.Run("INPRN-B03: A markdown file in a plans directory is a plan, even inside a guides directory", func(t *testing.T) {})
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		"docs/guides/plans/0001.md": "# plan",
		"docs/guides/STYLE.md":      "# guide",
	})
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.PlanDir != "docs/guides/plans" {
		t.Errorf("PlanDir = %q, want docs/guides/plans", p.PlanDir)
	}
	if !reflect.DeepEqual(p.GuideFiles, []string{"docs/guides/STYLE.md"}) {
		t.Errorf("the plan must not be counted as a guide, GuideFiles = %v", p.GuideFiles)
	}
}

// A layered project keeps a layer in a folder of two files as readily as in one of fifty.
// The old minimum of ten dropped those layers from the proposal.
func TestInfer_everyFolderWithCodeIsACandidate(t *testing.T) {
	t.Run("INPRN-B05: Every folder holding code is a candidate code directory, with no minimum, ordered by volume", func(t *testing.T) {})
	dir := t.TempDir()
	files := codeFiles("small/x/deep", ".go", "", 2)
	for k, v := range codeFiles("mid/y/deep", ".go", "", 10) {
		files[k] = v
	}
	for k, v := range codeFiles("big/z", ".go", "", 15) {
		files[k] = v
	}
	files["main.go"] = ""
	writeTree(t, dir, files)
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.CodeDirs, []string{"big/z", "mid/y", "small/x", "."}) {
		t.Errorf("CodeDirs = %v, want [big/z mid/y small/x .]", p.CodeDirs)
	}
}

func TestInfer_topFiveExtensions(t *testing.T) {
	t.Run("INPRN-B06: The code extensions are the five most frequent, most frequent first", func(t *testing.T) {})
	dir := t.TempDir()
	files := map[string]string{}
	for i, ext := range []string{".ts", ".go", ".py", ".rb", ".java", ".rs"} {
		for k, v := range codeFiles("src/"+ext[1:], ext, "", 12-i) {
			files[k] = v
		}
	}
	writeTree(t, dir, files)
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.CodeExts, []string{".ts", ".go", ".py", ".rb", ".java"}) {
		t.Errorf("CodeExts = %v, want the five most frequent in order", p.CodeExts)
	}
}

func TestInfer_colocationNeedsThreeStems(t *testing.T) {
	t.Run("INPRN-B07: Colocation is detected when at least three stems pair code with a spec, test or feature", func(t *testing.T) {})
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		"src/a.ts": "", "src/a.test.ts": "",
		"src/b.ts": "", "src/b.spec.md": "",
		"src/c.ts": "", "other/c.feature": "", // a different directory: not the same stem
	})
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Colocated {
		t.Error("two paired stems are a coincidence, not colocation")
	}
}

func TestInfer_testHandle(t *testing.T) {
	t.Run("INPRN-B08: The test handle is the known attribute used most, and only from five uses", func(t *testing.T) {})
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"dominant", `testID="a" testID="b" data-testid="c" data-testid="d" data-testid="e" data-testid="f" data-testid="g"`, "data-testid"},
		{"below five uses", `testID="a" testID="b" testID="c" testID="d"`, ""},
		{"tie goes to the earlier known attribute", `data-testid="a" testID="b"`, "testID"},
	}
	for _, c := range cases {
		dir := t.TempDir()
		// the tie case reaches the threshold by repetition over files
		n := 1
		if c.name != "below five uses" {
			n = 5
		}
		writeTree(t, dir, codeFiles("src", ".tsx", c.content, n))
		p, err := Infer(dir)
		if err != nil {
			t.Fatal(err)
		}
		if p.TestHandle != c.want {
			t.Errorf("%s: TestHandle = %q, want %q", c.name, p.TestHandle, c.want)
		}
	}
}

func TestInfer_handleSampleIsBounded(t *testing.T) {
	t.Run("INPRN-X01: Only the first three hundred code files are read in search of the test handle", func(t *testing.T) {})
	dir := t.TempDir()
	files := codeFiles("a/plain", ".ts", "", 300)
	for k, v := range codeFiles("z/marked", ".ts", `testID="x"`, 20) {
		files[k] = v
	}
	writeTree(t, dir, files)
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.TestHandle != "" {
		t.Errorf("the handle only appears after the sample, so none is proposed; got %q", p.TestHandle)
	}
}

func TestInfer_walkFailure(t *testing.T) {
	t.Run("INPRN-E01: A root that cannot be walked fails the inference with no proposal", func(t *testing.T) {})
	p, err := Infer(filepath.Join(t.TempDir(), "missing"))
	if err == nil || p != nil {
		t.Errorf("a missing root should fail with no proposal, got %+v %v", p, err)
	}
}

// A Go project names its tests `foo_test.go`. The stem of that file must be `foo`, or it
// never meets `foo.go` and a colocated Go project reads as a separate tree.
func TestInfer_testOfEveryDialectPairsWithItsCode(t *testing.T) {
	t.Run("INPRN-B10: A test named in any known dialect pairs with the code of the same stem", func(t *testing.T) {})
	for _, c := range []struct{ code, prefix, test string }{
		{".go", "", "_test.go"},
		{".py", "", "_test.py"},
		{".py", "test_", ".py"},
		{".ts", "", ".spec.ts"},
		{".java", "", "Test.java"},
	} {
		dir := t.TempDir()
		files := map[string]string{}
		for _, stem := range []string{"a", "b", "c"} {
			files["pkg/"+stem+c.code] = ""
			files["pkg/"+c.prefix+stem+c.test] = ""
		}
		writeTree(t, dir, files)
		p, err := Infer(dir)
		if err != nil {
			t.Fatal(err)
		}
		if !p.HasTest || !p.Colocated {
			t.Errorf("%s*%s beside %s: HasTest=%v Colocated=%v, want both true", c.prefix, c.test, c.code, p.HasTest, p.Colocated)
		}
	}
}

// Map order is random. A tie decided by it gave a different layer pattern, and a
// different order of code layers, from one run to the next over the same tree.
func TestInfer_tiesGiveTheSameProposalOnEveryRun(t *testing.T) {
	t.Run("INPRN-I02: The same tree always gives the same code extensions and code directories, ties included", func(t *testing.T) {})
	dir := t.TempDir()
	files := map[string]string{}
	for _, ext := range []string{".ts", ".js", ".go", ".py"} {
		for k, v := range codeFiles("src/"+ext[1:], ext, "", 11) {
			files[k] = v
		}
	}
	writeTree(t, dir, files)
	wantExts := []string{".go", ".js", ".py", ".ts"}
	wantDirs := []string{"src/go", "src/js", "src/py", "src/ts"}
	for range 20 {
		p, err := Infer(dir)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(p.CodeExts, wantExts) || !reflect.DeepEqual(p.CodeDirs, wantDirs) {
			t.Fatalf("tied counts must break by name: exts %v dirs %v, want %v %v", p.CodeExts, p.CodeDirs, wantExts, wantDirs)
		}
		if got := p.Config.Layers["go-code"].Pattern; got != "src/go/**/*.{go,js,py,ts}" {
			t.Fatalf("layer pattern = %q, want the extensions in name order", got)
		}
	}
}

func TestInfer_familyFromManifestOrExtension(t *testing.T) {
	t.Run("INPRN-B11: The language family comes from the root manifest, or from the most frequent extension", func(t *testing.T) {})
	for _, c := range []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"manifest wins over extensions", merge(map[string]string{"go.mod": "module x"}, codeFiles("web", ".ts", "", 5)), "go"},
		{"no manifest: the most frequent extension", merge(codeFiles("app", ".py", "", 3), codeFiles("tools", ".go", "", 1)), "python"},
		{"neither", map[string]string{"README.md": "# x"}, ""},
	} {
		dir := t.TempDir()
		writeTree(t, dir, c.files)
		p, err := Infer(dir)
		if err != nil {
			t.Fatal(err)
		}
		if p.Family != c.want {
			t.Errorf("%s: Family = %q, want %q", c.name, p.Family, c.want)
		}
	}
}

func TestInfer_testConventionsByUse(t *testing.T) {
	t.Run("INPRN-B12: The test conventions are the forms the project's tests follow, most followed first", func(t *testing.T) {})
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		"src/a.spec.tsx": "", "src/b.spec.tsx": "", "src/c.spec.tsx": "",
		"src/d.test.ts": "",
	})
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range p.TestConventions {
		got = append(got, c.Prefix+"*"+c.Suffix)
	}
	if !reflect.DeepEqual(got, []string{"*.spec.tsx", "*.test.ts"}) {
		t.Errorf("conventions = %v, want [*.spec.tsx *.test.ts]", got)
	}
}

func merge(ms ...map[string]string) map[string]string {
	out := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

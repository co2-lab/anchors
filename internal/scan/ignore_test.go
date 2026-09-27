package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// TestGitignoreAnchoring guards the bug that erased 129 units from the map in silence.
//
// The project's `/data/` line means, for git, "the `data` directory AT THE ROOT". By
// dropping the leading slash, Anchors read it as "any segment named `data`" — and the whole
// of `amplify/data/models/` (the schema's specs and models) left the scan. The map was not
// wrong with noise: it was smaller, and nothing reported the absence.
func TestGitignoreAnchoring(t *testing.T) {
	t.Run("SCIGS-B06: A slash anchors a gitignore pattern at the root, no slash matches at any depth", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(strings.Join([]string{
		"/data/",       // anchored: only at the root
		"node_modules", // no slash: any depth
		"*.log",
		"probe*.ts",
		"docs/build", // inner slash: anchored
		"docs/saida", // same, with a name outside the built-in list
	}, "\n")), 0o644))

	ig := LoadIgnore(dir)
	cases := []struct {
		rel   string
		isDir bool
		want  bool
		why   string
	}{
		{"data", true, true, "`/data/` matches the directory at the root"},
		{"data/dump.json", false, true, "below an ignored directory"},
		{"amplify/data", true, false, "`/data/` is ANCHORED — it does not match a nested `data`"},
		{"amplify/data/models/AiUsage.ts", false, false, "the real case that vanished from the map"},
		{"a/node_modules", true, true, "no slash matches at any depth"},
		{"a/b/c.log", false, true, "`*.log` matches the basename at any depth"},
		{"probe1.ts", false, true, "the reviewer's probe that became a task"},
		{"docs/build", true, true, "inner slash: anchored, and it matches"},
		// `apps/docs/build` is NOT an anchoring case: `build` is in the built-in list, which
		// holds in every project and is applied before the `.gitignore`. A nested pair that
		// proves anchoring needs a name outside that list.
		{"docs/saida", true, true, "anchored, at the root: matches"},
		{"apps/docs/saida", true, false, "anchored: `docs/saida` does not match nested"},
	}
	for _, c := range cases {
		var got bool
		if c.isDir {
			got = ig.SkipDir(filepath.Base(c.rel), c.rel)
		} else {
			got = ig.SkipFile(c.rel)
		}
		if got != c.want {
			t.Errorf("%q (dir=%v): ignored=%v, want %v — %s", c.rel, c.isDir, got, c.want, c.why)
		}
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// TestBuiltInListIsDefeatable: `build`/`dist` are compiled output in most projects and a
// SOURCE FOLDER in some. A built-in list that cannot be contested is the framework deciding,
// for a project it does not know, what in it is disposable — and the error is silent: the
// layer vanishes from the map and no gate reports the absence.
func TestBuiltInListIsDefeatable(t *testing.T) {
	t.Run("SCIGS-B01: The built-in directories are skipped when nothing is declared", func(t *testing.T) {})
	t.Run("SCIGS-B02: A layer pointing inside a built-in directory re-enables it, a catch-all does not", func(t *testing.T) {})
	t.Run("SCIGS-B03: A gitignore negation re-enables a built-in directory", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.log\n"), 0o644))

	// nothing declared: the default holds
	if !LoadIgnore(dir).SkipDir("build", "build") {
		t.Error("with nothing declared, `build` follows the default (ignored)")
	}

	// the STRUCTURE declares there is code there
	cfg := &config.Config{Layers: map[string]config.Layer{
		"core": {Pattern: "build/**/*.ts", Kind: "code"},
	}}
	if LoadIgnoreFor(dir, cfg).SkipDir("build", "build") {
		t.Error("a layer pointing inside `build` must DEFEAT the default")
	}

	// a catch-all does NOT re-enable — otherwise `**/*.ts` would bring `node_modules` back
	broad := &config.Config{Layers: map[string]config.Layer{
		"all": {Pattern: "**/*.ts", Kind: "code"},
	}}
	if !LoadIgnoreFor(dir, broad).SkipDir("node_modules", "node_modules") {
		t.Error("a catch-all cannot re-enable `node_modules`")
	}

	// the negation in the `.gitignore` defeats it too
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("!build/\n"), 0o644))
	if LoadIgnore(dir).SkipDir("build", "build") {
		t.Error("`!build/` in the .gitignore must defeat the default")
	}
}

// The machinery is NOT defeatable, under any declaration.
func TestMachineryIsNeverScanned(t *testing.T) {
	t.Run("SCIGS-I01: No declaration re-enables the machinery directories", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("!.git/\n!.anchors/\n"), 0o644))
	ig := LoadIgnore(dir)
	for _, d := range []string{".git", ".anchors"} {
		if !ig.SkipDir(d, d) {
			t.Errorf("%q is machinery — no declaration re-enables it", d)
		}
	}
}

// What Anchors writes is not work for Anchors: the queue fed on its own output.
func TestAnchorsRecordsAreNeverScanned(t *testing.T) {
	t.Run("SCIGS-B04: The records Anchors writes are never scanned", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("!issues/\n"), 0o644))
	ig := LoadIgnore(dir)
	if !ig.SkipDir("issues", "issues") {
		t.Error("`issues` is Anchors' own output, even with a negation")
	}
	if !ig.SkipDir("changes", "a/changes") {
		t.Error("`changes` is skipped at any depth")
	}
}

// TestEphemeraNeverBecomeWork guards the noise measured in a real E2E: the watcher queued a
// task for `amplify/data/.!21662!resource.spec.md` — a file the editor creates for
// milliseconds during an atomic save and that matches the spec glob. Three tasks had to be
// discarded by hand.
//
// The project's `.gitignore` does not cover this, and it should not: it is not the
// project's decision, it is file-system noise.
func TestEphemeraNeverBecomeWork(t *testing.T) {
	t.Run("SCIGS-B05: Editor and system ephemera never become files to scan", func(t *testing.T) {})
	ig := LoadIgnore(t.TempDir())
	noise := []string{
		"amplify/data/.!21662!resource.spec.md", // the real case
		"src/a.ts.swp", "src/a.ts~", "src/.#a.ts", "src/#a.ts#",
		"src/x.tmp", ".DS_Store",
	}
	for _, r := range noise {
		if !ig.SkipFile(r) {
			t.Errorf("%q is editor/system noise and cannot become work", r)
		}
	}
	material := []string{
		"src/a.ts", "src/a.spec.md", "amplify/data/resource.spec.md",
		"src/tmp/util.ts", // `tmp` in the PATH is not `.tmp` in the name
	}
	for _, m := range material {
		if ig.SkipFile(m) {
			t.Errorf("%q is project material and was discarded", m)
		}
	}
}

// A trailing slash restricts what the pattern MATCHES to directories, not what it covers.
func TestDirOnlyPatternCoversBelowButNotAFile(t *testing.T) {
	t.Run("SCIGS-B07: A trailing slash ignores the directory and what is below it, but not a file of that name", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("data/\n"), 0o644))
	ig := LoadIgnore(dir)
	if !ig.SkipDir("data", "data") {
		t.Error("`data/` ignores the directory")
	}
	if !ig.SkipFile("data/dump.json") {
		t.Error("`data/` covers the files below it")
	}
	if ig.SkipFile("data") {
		t.Error("`data/` matches directories only, not a plain file named `data`")
	}
}

// Git's order: the last rule that matches wins, and a negation re-includes.
func TestLastMatchingRuleWins(t *testing.T) {
	t.Run("SCIGS-B08: The last matching gitignore rule decides", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("# comment\n\n*.log\n!keep.log\n"), 0o644))
	ig := LoadIgnore(dir)
	if !ig.SkipFile("a.log") {
		t.Error("`*.log` ignores a.log")
	}
	if ig.SkipFile("keep.log") {
		t.Error("the later `!keep.log` re-includes keep.log")
	}
}

func TestNilIgnoreKeepsTheFixedExclusions(t *testing.T) {
	t.Run("SCIGS-B09: Without a loaded ignore set the fixed exclusions still hold", func(t *testing.T) {})
	var ig *Ignore
	if !ig.SkipDir("node_modules", "node_modules") || !ig.SkipDir("issues", "issues") {
		t.Error("the built-in directories and Anchors' records are skipped without an ignore set")
	}
	if !ig.SkipFile("a.swp") {
		t.Error("ephemera are skipped without an ignore set")
	}
	if ig.SkipFile("a.ts") {
		t.Error("nothing else is ignored without an ignore set")
	}
}

func TestNestedGitignoreIsNotRead(t *testing.T) {
	t.Run("SCIGS-X01: A nested gitignore does not change what the scan sees", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(dir, "sub"), 0o755))
	must(t, os.WriteFile(filepath.Join(dir, "sub", ".gitignore"), []byte("x.ts\n"), 0o644))
	if LoadIgnore(dir).SkipFile("sub/x.ts") {
		t.Error("only the root .gitignore is read")
	}
}

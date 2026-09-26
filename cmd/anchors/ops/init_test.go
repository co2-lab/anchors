package ops

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
)

// --- helpers shared by the tests of this package ---

// captureStdout runs fn with os.Stdout redirected and returns what it printed. The pipe is
// drained while fn runs, so a long output (the work order, a JSON dump) cannot fill the
// pipe buffer and block the test.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	os.Stdout = w
	defer func() { os.Stdout = orig }()
	fn()
	w.Close()
	os.Stdout = orig
	return <-done
}

// skipIfTerminal skips a test that drives the interactive prompts. `huh` opens /dev/tty,
// so under a developer's terminal the prompt would WAIT for a keypress instead of failing
// the way it does in CI, a pipe or an agent — which is the path these tests prove.
func skipIfTerminal(t *testing.T) {
	t.Helper()
	if f, err := os.Open("/dev/tty"); err == nil {
		f.Close()
		t.Skip("a controlling terminal is available: the prompts would wait for input")
	}
}

// resetPromptError isolates the package-level `erroDePrompt` flag: a prompt that failed in
// one test must not make the next one abort.
func resetPromptError(t *testing.T) {
	t.Helper()
	erroDePrompt = false
	t.Cleanup(func() { erroDePrompt = false })
}

// isolateGit keeps git away from the machine's configuration: no global or system config
// (a global core.hooksPath would send the hooks outside the temporary repository), and an
// author identity valid only for this test.
func isolateGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	global := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(global, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", global)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	identidadeGitNoTeste(t)
}

// newGitRepo creates a temporary repository with one commit (so HEAD exists).
func newGitRepo(t *testing.T) string {
	t.Helper()
	isolateGit(t)
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"commit", "-q", "--allow-empty", "-m", "root"},
	} {
		if out, err := runGit(dir, args...); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

// writeFile writes rel under root, creating the directories.
func writeFile(t *testing.T, root, rel, content string) string {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return full
}

// fakeGH puts a fake `gh` first on the PATH. `body` is the bash that answers each call
// (the arguments are in "$@"); every call is also appended, one per line, to the returned
// log. Nothing reaches GitHub.
func fakeGH(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	logPath := filepath.Join(dir, "gh.log")
	script := "#!/bin/bash\necho \"$*\" >> '" + logPath + "'\n" + body + "\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return logPath
}

func readLog(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return string(b)
}

// --- init.go ---

// An existing anchors.yaml is only overwritten on an explicit yes. Without a terminal the
// confirmation cannot run, and the answer must count as NO: the file stays byte for byte.
func TestRunInitKeepsAnExistingConfigWhenTheConfirmationCannotRun(t *testing.T) {
	skipIfTerminal(t)
	resetPromptError(t)
	root := t.TempDir()
	const original = "version: 1\n# mine\n"
	writeFile(t, root, config.DefaultFile, original)

	var err error
	out := captureStdout(t, func() { err = runInit(root) })
	if err != nil {
		t.Fatalf("runInit: %v", err)
	}
	b, _ := os.ReadFile(filepath.Join(root, config.DefaultFile))
	if string(b) != original {
		t.Errorf("anchors.yaml was touched without a confirmation:\n%s", b)
	}
	if strings.Contains(out, "written") {
		t.Errorf("the output announces a write that must not happen:\n%s", out)
	}
}

// Outside git, the first question is the `git init` offer. When it cannot be asked, init
// aborts BEFORE touching anything: no repository created, no anchors.yaml.
func TestRunInitWithoutTTYDoesNotTouchGitNorWrite(t *testing.T) {
	skipIfTerminal(t)
	resetPromptError(t)
	isolateGit(t)
	root := t.TempDir()

	var err error
	out := captureStdout(t, func() { err = runInit(root) })
	if err == nil {
		t.Fatal("init without a terminal must fail, not proceed on answers nobody gave")
	}
	if !strings.Contains(err.Error(), "git was not touched") {
		t.Errorf("the error must say git was not touched, got: %v", err)
	}
	if !strings.Contains(err.Error(), "--non-interactive") {
		t.Errorf("the error must point at the non-interactive mode, got: %v", err)
	}
	if _, serr := os.Stat(filepath.Join(root, ".git")); serr == nil {
		t.Error("a repository was created although nobody accepted the offer")
	}
	if _, serr := os.Stat(filepath.Join(root, config.DefaultFile)); serr == nil {
		t.Error("anchors.yaml was written")
	}
	if !strings.Contains(out, initx.AvisoGit(initx.GitNaoIniciado)) {
		t.Errorf("the git warning was not printed before the question:\n%s", out)
	}
}

// Inside a ready repository init scans, reports what it found and asks. With every prompt
// failing, each answer comes back empty, and writing would produce a config that loads and
// governs nothing — so it must refuse to write anchors.yaml.
func TestRunInitRefusesToWriteAConfigBuiltFromFailedPrompts(t *testing.T) {
	skipIfTerminal(t)
	resetPromptError(t)
	root := newGitRepo(t)
	// Ten code files: below that a directory is anecdotal and init does not call it code.
	for i := 0; i < 10; i++ {
		writeFile(t, root, "src/app/m"+string(rune('a'+i))+".ts", "export const x = 1\n")
	}
	writeFile(t, root, "src/app/calc.ts", "export const x = 1\n")
	writeFile(t, root, "src/app/calc.spec.md", "# Calc\n")
	writeFile(t, root, "src/app/calc.feature", "Feature: calc\n")
	writeFile(t, root, "src/app/calc.test.ts", "test('x', () => {})\n")
	writeFile(t, root, "guides/STYLE_GUIDE.md", "# Style\n")

	var err error
	out := captureStdout(t, func() { err = runInit(root) })
	if err == nil {
		t.Fatal("with every prompt failing, init must not report success")
	}
	if !strings.Contains(err.Error(), "0 layers and 0 gates") {
		t.Errorf("the error must say why nothing was written, got: %v", err)
	}
	if _, serr := os.Stat(filepath.Join(root, config.DefaultFile)); serr == nil {
		t.Error("anchors.yaml was written from answers nobody gave")
	}
	// The findings come from the disk, not from the prompts: they are printed anyway.
	for _, want := range []string{"specs (*.spec.md)", "features (*.feature)", "tests (*.test.*)", "guides in guides/", "code in:"} {
		if !strings.Contains(out, want) {
			t.Errorf("the findings do not report %q:\n%s", want, out)
		}
	}
	// A project with code does not need the DISCOVER phase: no work order.
	if strings.Contains(out, "DISCOVER phase") {
		t.Errorf("a project with code was sent to the DISCOVER phase:\n%s", out)
	}
}

// An empty directory under git: the DISCOVER phase is due, and since nobody can answer a
// prompt here, the reader is treated as an agent and gets the work order.
func TestRunInitInAnEmptyRepoPrintsTheWorkOrder(t *testing.T) {
	skipIfTerminal(t)
	resetPromptError(t)
	root := newGitRepo(t)

	var err error
	out := captureStdout(t, func() { err = runInit(root) })
	if err == nil {
		t.Fatal("init must not succeed on failed prompts")
	}
	if !strings.Contains(out, "anchors guide project") {
		t.Errorf("the empty project did not get the DISCOVER work order:\n%s", out)
	}
	if !strings.Contains(out, "Found:") {
		t.Errorf("the findings header is missing:\n%s", out)
	}
}

// Every prompt helper reports the failure through `erroDePrompt` and hands back the
// DEFAULT (or nothing) — never a choice nobody made.
func TestPromptHelpersFlagTheFailureAndReturnNoChoice(t *testing.T) {
	skipIfTerminal(t)

	resetPromptError(t)
	if got := askText("repo?"); got != "" || !erroDePrompt {
		t.Errorf("askText = %q, flag = %v; want empty and flagged", got, erroDePrompt)
	}
	resetPromptError(t)
	if got := askConfirmDefault("ok?", false); got || !erroDePrompt {
		t.Errorf("askConfirmDefault(false) = %v, flag = %v; want the default and flagged", got, erroDePrompt)
	}
	resetPromptError(t)
	if got := askMultiSelectPre("which?", []string{"a", "b"}, map[string]bool{"a": true}); len(got) != 0 || !erroDePrompt {
		t.Errorf("askMultiSelectPre = %v, flag = %v; want nothing chosen and flagged", got, erroDePrompt)
	}
	resetPromptError(t)
	if got := askMultiSelect("which?", []string{"a"}); len(got) != 0 || !erroDePrompt {
		t.Errorf("askMultiSelect = %v, flag = %v", got, erroDePrompt)
	}
	resetPromptError(t)
	// The select pre-fills its first option, so the tag itself means nothing here: what
	// counts is that the failure is flagged, which makes runInit refuse to save.
	got := askGovernAnswers([]string{"guides/A_GUIDE.md"}, []string{"backend"})
	if _, ok := got["guides/A_GUIDE.md"]; !ok || !erroDePrompt {
		t.Errorf("askGovernAnswers = %v, flag = %v; want an entry per guide, flagged", got, erroDePrompt)
	}
	resetPromptError(t)
	if p, ok := askPreset(); ok || p.Name != "" || !erroDePrompt {
		t.Errorf("askPreset = (%q, %v), flag = %v; want no preset, flagged", p.Name, ok, erroDePrompt)
	}
}

// The module prefixes of a modular preset come from the module directories that EXIST
// under its glob. Files there are not modules, and a non-modular preset has none.
func TestDetectModulesListsOnlyTheModuleDirectories(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "src/modules/users/users.service.ts", "")
	writeFile(t, root, "src/modules/billing/billing.service.ts", "")
	writeFile(t, root, "src/modules/README.md", "")

	var nodeTS initx.Preset
	for _, p := range initx.Presets {
		if p.Name == "node-ts" {
			nodeTS = p
		}
	}
	if !nodeTS.Modular {
		t.Fatal("the node-ts preset is expected to be modular")
	}
	got := strings.Join(detectModules(root, nodeTS), ",")
	if got != "src/modules/billing,src/modules/users" {
		t.Errorf("modules = %q, want the two directories (and not README.md)", got)
	}
	if m := detectModules(root, initx.Preset{Name: "flat"}); m != nil {
		t.Errorf("a non-modular preset must have no modules, got %v", m)
	}
}

// printFindings says only what was found: an empty proposal lists nothing.
func TestPrintFindingsReportsOnlyWhatExists(t *testing.T) {
	out := captureStdout(t, func() {
		printFindings(&initx.Proposal{Colocated: true, CodeDirs: []string{"src"}, CodeExts: []string{".go"}})
	})
	if !strings.Contains(out, "code in: [src]") || !strings.Contains(out, "co-location") {
		t.Errorf("the findings miss what the proposal has:\n%s", out)
	}
	for _, absent := range []string{"specs", "features", "tests", "guides"} {
		if strings.Contains(out, absent) {
			t.Errorf("the findings report %q, which the proposal does not have:\n%s", absent, out)
		}
	}
}

// The `init` command routes --non-interactive to the flag-driven mode instead of the TUI.
func TestInitCommandNonInteractiveRoutesToTheJSONMode(t *testing.T) {
	root := t.TempDir()
	cmd := newInitCmd()
	cmd.SetArgs([]string{"--root", root, "--non-interactive"})
	cmd.SetOut(io.Discard)
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatalf("init --non-interactive: %v", err)
	}
	if !strings.Contains(out, `"perguntas"`) {
		t.Errorf("the command did not emit the questions:\n%s", out)
	}
}

// With TERM=dumb huh switches to line mode, reads os.Stdin, and at end-of-input returns
// each prompt's DEFAULT with no error — so the `erroDePrompt` guard never tripped, and
// `TERM=dumb anchors init </dev/null` wrote an anchors.yaml with 0 artifacts and 0 gates
// under a "✓ written". End-of-input is not an answer: init must refuse to write.
func TestRunInitInLineModeRefusesToWriteOnEndOfInput(t *testing.T) {
	resetPromptError(t)
	root := newGitRepo(t)
	for i := 0; i < 10; i++ {
		writeFile(t, root, "src/app/m"+string(rune('a'+i))+".ts", "export const x = 1\n")
	}
	writeFile(t, root, "src/app/calc.spec.md", "# Calc\n")
	if out, err := runGit(root, "add", "-A"); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
	if out, err := runGit(root, "commit", "-q", "-m", "code"); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}

	t.Setenv("TERM", "dumb")
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()
	origStdin := os.Stdin
	os.Stdin = devNull
	defer func() { os.Stdin = origStdin }()

	var runErr error
	out := captureStdout(t, func() { runErr = runInit(root) })
	if runErr == nil {
		t.Fatalf("init on an input that ended must not succeed:\n%s", out)
	}
	if !strings.Contains(runErr.Error(), "--non-interactive") {
		t.Errorf("the error must point at the non-interactive mode, got: %v", runErr)
	}
	if _, serr := os.Stat(filepath.Join(root, config.DefaultFile)); serr == nil {
		t.Error("anchors.yaml was written from answers nobody gave")
	}
	// Nothing at all: the header guide used to be seeded before the refusal.
	if _, serr := os.Stat(filepath.Join(root, "guides", "HEADER_GUIDE.md")); serr == nil {
		t.Error("an init that refused left guides/HEADER_GUIDE.md behind")
	}
	if strings.Contains(out, "written") {
		t.Errorf("the output announces a write that did not happen:\n%s", out)
	}
}

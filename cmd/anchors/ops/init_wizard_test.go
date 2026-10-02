package ops

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/initx"
)

// --- helpers shared by the tests of this package ---

// captureStdout runs fn with os.Stdout redirected and returns what it printed. The pipe is
// drained while fn runs, so a long output (the work order, a JSON dump) cannot fill the
// pipe buffer and block the test.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	return testkit.CaptureStdout(t, fn)
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
	gitIdentityForTest(t)
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
	testkit.FakeBin(t, dir, "gh", script)
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
	t.Run("INWZN-B01: An existing config is kept when the overwrite is not confirmed", func(t *testing.T) {})
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
	t.Run("INWZN-X01: Git is never initialized nor committed without a yes", func(t *testing.T) {})
	t.Run("INWZN-I01: A prompt that cannot run makes the init write nothing", func(t *testing.T) {})
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
	t.Run("INWZN-I01: A prompt that cannot run makes the init write nothing", func(t *testing.T) {})
	t.Run("INWZN-B08: The findings report only what exists on disk", func(t *testing.T) {})
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
	for _, want := range []string{"specs (*.spec.md)", "features (*.feature)", "tests (**/*.test.ts)", "language: ts", "guides in guides/", "code in:"} {
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
	t.Run("INWZN-B10: An AI operator gets the DISCOVER work order", func(t *testing.T) {})
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
	t.Run("INWZN-I01: A prompt that cannot run makes the init write nothing", func(t *testing.T) {})
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
}

// printFindings says only what was found: an empty proposal lists nothing.
func TestPrintFindingsReportsOnlyWhatExists(t *testing.T) {
	t.Run("INWZN-B08: The findings report only what exists on disk", func(t *testing.T) {})
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
	// No code directory, no code line: not even an empty list.
	out = captureStdout(t, func() { printFindings(&initx.Proposal{}) })
	if strings.Contains(out, "code in") {
		t.Errorf("the findings report code for a proposal without it:\n%s", out)
	}
}

// The `init` command routes --non-interactive to the flag-driven mode instead of the TUI.
func TestInitCommandNonInteractiveRoutesToTheJSONMode(t *testing.T) {
	t.Run("INWZN-B14: Non-interactive routes to the JSON mode", func(t *testing.T) {})
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
	t.Run("INWZN-E02: End of input in line mode is refused, not taken as the defaults", func(t *testing.T) {})
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

// promptStep is one answer of a scripted terminal: `answer` is typed once `waitFor`
// shows on the screen.
type promptStep struct{ waitFor, answer string }

// driveLineMode runs fn in huh's line mode (TERM=dumb) against a scripted terminal:
// stdin and stdout are pipes, and each answer is typed only after its prompt is on the
// screen. The timing matters: every prompt builds a fresh line reader, so an answer typed
// ahead would be swallowed by the reader of the prompt before it. After the last step, or
// when an expected prompt does not show within the deadline, the input is closed, so any
// further prompt hits end of input and the init refuses to write. It returns the whole
// output and the prompts that never showed.
func driveLineMode(t *testing.T, steps []promptStep, fn func()) (string, []string) {
	t.Helper()
	t.Setenv("TERM", "dumb")
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	var (
		mu     sync.Mutex
		screen strings.Builder
	)
	shown := make(chan struct{}, 1)
	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		buf := make([]byte, 4096)
		for {
			n, rerr := outR.Read(buf)
			mu.Lock()
			screen.Write(buf[:n])
			mu.Unlock()
			select {
			case shown <- struct{}{}:
			default:
			}
			if rerr != nil {
				return
			}
		}
	}()
	stop := make(chan struct{})
	var unmet []string
	driveDone := make(chan struct{})
	go func() {
		defer close(driveDone)
		defer inW.Close()
		cursor := 0
		for i, s := range steps {
			for {
				mu.Lock()
				rest := screen.String()[cursor:]
				mu.Unlock()
				if at := strings.Index(rest, s.waitFor); at >= 0 {
					cursor += at + len(s.waitFor)
					if _, werr := inW.WriteString(s.answer + "\n"); werr != nil {
						unmet = append(unmet, s.waitFor)
						return
					}
					break
				}
				select {
				case <-shown:
				case <-stop:
					return
				case <-time.After(10 * time.Second):
					for _, left := range steps[i:] {
						unmet = append(unmet, left.waitFor)
					}
					return
				}
			}
		}
	}()
	origIn, origOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = inR, outW
	fn()
	os.Stdin, os.Stdout = origIn, origOut
	close(stop)
	<-driveDone
	outW.Close()
	<-readDone
	inR.Close()
	return screen.String(), unmet
}

// repoWithCode is a repository whose only commit holds ten code files and a spec.
func repoWithCode(t *testing.T) string {
	t.Helper()
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
	return root
}

// answerEverything is the scripted terminal of an init where every question gets an
// answer: the only path that writes.
func answerEverything() []promptStep {
	return []promptStep{
		{"already exists. Overwrite?", "y"},
		{"Seed guides/HEADER_GUIDE.md", "y"},
		{"Seed CONTRIBUTING.md", "y"},
		{"Which anchor kinds", "0"},
		{"default gate(s)", "y"},
		{"derivatives", "n"},
		{"Which code directories", "0"},
		{"work queue live in GitHub", "n"},
		{"write its findings as files", "y"},
	}
}

// The path where every question gets an answer: the only one that writes. It proves the
// yes of the overwrite, the header guide landing in guides/ when the project has no guide
// directory, the accepted gates in the file, the note on layers and the seeded
// CONTRIBUTING.md.
func TestRunInitWritesWhatTheAnswersChose(t *testing.T) {
	t.Run("INWZN-B01: An existing config is kept when the overwrite is not confirmed", func(t *testing.T) {})
	t.Run("INWZN-B16: The header guide is seeded in guides/ when the project has no guide directory", func(t *testing.T) {})
	t.Run("INWZN-B17: The default gates are offered only when the chosen artifacts have any, and accepted ones are written", func(t *testing.T) {})
	t.Run("INWZN-B22: The code-layer question is preceded by the note on layers", func(t *testing.T) {})
	t.Run("INWZN-B23: CONTRIBUTING.md is seeded when absent, and an existing one is shown, not touched", func(t *testing.T) {})
	t.Run("INWZN-B24: The family's coverage hint is printed", func(t *testing.T) {})
	t.Run("INWZN-B25: Accepting the gates tells to write gate steps in the project's language", func(t *testing.T) {})
	resetPromptError(t)
	root := repoWithCode(t)
	writeFile(t, root, "go.mod", "module x\n")
	const original = "version: 1\n# mine\n"
	writeFile(t, root, config.DefaultFile, original)

	var runErr error
	out, unmet := driveLineMode(t, answerEverything(), func() { runErr = runInit(root) })
	if runErr != nil {
		t.Fatalf("init with every question answered failed: %v\n%s", runErr, out)
	}
	if len(unmet) > 0 {
		t.Errorf("prompts that never showed: %q\n%s", unmet, out)
	}
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		t.Fatalf("the written anchors.yaml does not load: %v", err)
	}
	if len(cfg.Gates) == 0 {
		t.Errorf("the accepted default gates are not in anchors.yaml")
	}
	if _, ok := cfg.Layers["app-code"]; !ok {
		t.Errorf("the project's own folder src/app is not a code layer: %v", cfg.Layers)
	}
	if _, err := os.Stat(filepath.Join(root, "guides", "HEADER_GUIDE.md")); err != nil {
		t.Errorf("the header guide is not in guides/: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "HEADER_GUIDE.md")); err == nil {
		t.Error("the header guide was written at the root")
	}
	note := i18n.T("init.layers_note")
	if at, q := strings.Index(out, note), strings.Index(out, "Which code directories"); at < 0 || at > q {
		t.Errorf("the note on layers should come before the code-layer question:\n%s", out)
	}
	if b, err := os.ReadFile(filepath.Join(root, "CONTRIBUTING.md")); err != nil || !strings.Contains(string(b), "`app-code`") {
		t.Errorf("CONTRIBUTING.md should be seeded naming the code layers: %v\n%s", err, b)
	}
	if !strings.Contains(out, initx.CoverageHint("go")) {
		t.Errorf("the go coverage hint was not printed:\n%s", out)
	}
	if !strings.Contains(out, "not in shell") || !strings.Contains(out, "go run ./tools/gates <gate>") {
		t.Errorf("accepting the gates should tell to write gate steps in the project's language:\n%s", out)
	}
}

// A CONTRIBUTING.md the project already has is the project's: init shows what it would
// add and writes nothing into it.
func TestRunInitLeavesAnExistingContributingAlone(t *testing.T) {
	t.Run("INWZN-B23: CONTRIBUTING.md is seeded when absent, and an existing one is shown, not touched", func(t *testing.T) {})
	resetPromptError(t)
	root := repoWithCode(t)
	writeFile(t, root, config.DefaultFile, "version: 1\n")
	writeFile(t, root, "CONTRIBUTING.md", "ours\n")
	var runErr error
	out, _ := driveLineMode(t, answerEverything(), func() { runErr = runInit(root) })
	if runErr != nil {
		t.Fatalf("init: %v\n%s", runErr, out)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "CONTRIBUTING.md")); string(b) != "ours\n" {
		t.Errorf("the project's CONTRIBUTING.md was touched:\n%s", b)
	}
	if !strings.Contains(out, "already exists — left untouched") || !strings.Contains(out, "## Working with Anchors") {
		t.Errorf("the init should print the section it would add:\n%s", out)
	}
}

// --preset proposed a stack's structure; the flag is gone, and whoever still passes it is
// told why instead of reading "unknown flag".
func TestInitRefusesThePresetFlag(t *testing.T) {
	t.Run("INWZN-B21: --preset is refused with the reason, and nothing is written", func(t *testing.T) {})
	for _, mode := range [][]string{nil, {"--non-interactive"}} {
		root := t.TempDir()
		cmd := newInitCmd()
		cmd.SetArgs(append([]string{"--root", root, "--preset=go"}, mode...))
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		err := cmd.Execute()
		if err == nil || !strings.Contains(err.Error(), "proposes no project structure") {
			t.Errorf("%v: --preset should be refused with the reason, got %v", mode, err)
		}
		if _, serr := os.Stat(filepath.Join(root, config.DefaultFile)); serr == nil {
			t.Errorf("%v: --preset wrote anchors.yaml", mode)
		}
	}
}

// runInitAtEndOfInput runs the init in line mode on an input that already ended. Every
// prompt still prints its question before it fails, so the output lists which questions
// the init asked; it then refuses to write, which the caller need not check again.
func runInitAtEndOfInput(t *testing.T, root string) string {
	t.Helper()
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
	return out
}

// An empty project is told it is new, gets no question about code layers (there are none,
// and it is told so) and, with no artifact chosen, no offer of zero gates.
func TestRunInitOnAnEmptyProjectSkipsTheQuestionsWithNothingToAsk(t *testing.T) {
	t.Run("INWZN-B18: A project with no code, spec, feature or test is announced as new", func(t *testing.T) {})
	t.Run("INWZN-B19: The code-layer question is asked only when there are code layers, and a new project is told to declare them later", func(t *testing.T) {})
	t.Run("INWZN-B17: The default gates are offered only when the chosen artifacts have any, and accepted ones are written", func(t *testing.T) {})
	resetPromptError(t)
	out := runInitAtEndOfInput(t, newGitRepo(t))
	for _, want := range []string{i18n.T("init.empty_project"), i18n.T("init.no_code_yet")} {
		if !strings.Contains(out, want) {
			t.Errorf("an empty project was not told %q:\n%s", want, out)
		}
	}
	for _, absent := range []string{"Which code directories", "default gate(s)"} {
		if strings.Contains(out, absent) {
			t.Errorf("an empty project was asked %q:\n%s", absent, out)
		}
	}
}

// A project with code and a guide is asked which code directories are layers and which
// tag the guide governs, and is not called new.
func TestRunInitOnAProjectWithCodeAndAGuideAsksAboutBoth(t *testing.T) {
	t.Run("INWZN-B18: A project with no code, spec, feature or test is announced as new", func(t *testing.T) {})
	t.Run("INWZN-B19: The code-layer question is asked only when there are code layers, and a new project is told to declare them later", func(t *testing.T) {})
	t.Run("INWZN-B20: Each guide found is asked which tag it governs", func(t *testing.T) {})
	resetPromptError(t)
	root := repoWithCode(t)
	writeFile(t, root, "guides/STYLE_GUIDE.md", "# Style\n")
	out := runInitAtEndOfInput(t, root)
	for _, want := range []string{"Which code directories", "The guide STYLE_GUIDE.md governs which tag?"} {
		if !strings.Contains(out, want) {
			t.Errorf("the init did not ask %q:\n%s", want, out)
		}
	}
	for _, absent := range []string{i18n.T("init.empty_project"), i18n.T("init.no_code_yet")} {
		if strings.Contains(out, absent) {
			t.Errorf("a project with code was told %q:\n%s", absent, out)
		}
	}
}

// --- init_discover.go ---

// clearAgentEnv removes every variable that identifies an AI agent, so the test decides
// who is operating instead of inheriting it from whoever runs the suite.
func clearAgentEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"CLAUDE_CODE_ENTRYPOINT", "CLAUDECODE", "CURSOR_TRACE_ID",
		"AIDER_MODEL", "GEMINI_CLI", "CODEX_SANDBOX", "AI_AGENT"} {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}
}

// A project that already has code (or PROJECT.md) does not need the DISCOVER phase: the
// step says nothing and lets init go on.
func TestDiscoverStepIsSilentWhenThePhaseAlreadyHappened(t *testing.T) {
	t.Run("INWZN-B09: The DISCOVER step is silent when the phase already happened", func(t *testing.T) {})
	root := t.TempDir()
	for name, p := range map[string]*initx.Proposal{
		"with code":  {CodeDirs: []string{"src"}},
		"with specs": {HasSpecMD: true},
	} {
		var ok bool
		out := captureStdout(t, func() { ok = discoverStep(root, p) })
		if !ok || out != "" {
			t.Errorf("%s: discoverStep = %v, printed %q; want true and silence", name, ok, out)
		}
	}
	writeFile(t, root, "PROJECT.md", "# decisions\n")
	var ok bool
	out := captureStdout(t, func() { ok = discoverStep(root, &initx.Proposal{}) })
	if !ok || out != "" {
		t.Errorf("with PROJECT.md: discoverStep = %v, printed %q; want true and silence", ok, out)
	}
}

// An AI operating the CLI gets a WORK ORDER — the task, in order — and init goes on.
func TestDiscoverStepGivesAnAgentTheWorkOrder(t *testing.T) {
	t.Run("INWZN-B10: An AI operator gets the DISCOVER work order", func(t *testing.T) {})
	clearAgentEnv(t)
	t.Setenv("CLAUDECODE", "1")
	var ok bool
	out := captureStdout(t, func() { ok = discoverStep(t.TempDir(), &initx.Proposal{}) })
	if !ok {
		t.Error("the work order must not stop init")
	}
	for _, want := range []string{"YOU (the agent operating this CLI)", "anchors guide project",
		"PROJECT.md", "INSIGHTS.md", "never in a background worker"} {
		if !strings.Contains(out, want) {
			t.Errorf("the work order misses %q:\n%s", want, out)
		}
	}
}

// A person with no known AI installed gets the step-by-step, with the whole prompt ready
// to paste — and init goes on.
func TestInstructPersonWithoutAKnownAIPrintsThePromptToPaste(t *testing.T) {
	t.Run("INWZN-B11: A person without a known AI gets the prompt to paste", func(t *testing.T) {})
	clearAgentEnv(t)
	resetPromptError(t)
	var ok bool
	out := captureStdout(t, func() { ok = instructPerson(t.TempDir()) })
	if !ok {
		t.Error("the instructions must not stop init")
	}
	if !strings.Contains(out, "To do it yourself") {
		t.Errorf("the step-by-step is missing:\n%s", out)
	}
	// The prompt is reflowed, so compare word by word.
	if strings.Join(strings.Fields(out), " ") == "" ||
		!strings.Contains(strings.Join(strings.Fields(out), " "), strings.Join(strings.Fields(initx.PromptDescobrir), " ")) {
		t.Errorf("the full prompt is not in the output:\n%s", out)
	}
	if strings.Contains(out, "I detected") {
		t.Errorf("no AI is installed, yet the output offers to open one:\n%s", out)
	}
}

// With a known AI the person is offered to open it. When that question cannot be asked,
// instructPerson answers false — init must stop instead of guessing.
func TestInstructPersonStopsWhenTheOfferCannotBeAsked(t *testing.T) {
	t.Run("INWZN-I01: A prompt that cannot run makes the init write nothing", func(t *testing.T) {})
	skipIfTerminal(t)
	clearAgentEnv(t)
	resetPromptError(t)
	t.Setenv("GEMINI_CLI", "1")
	var ok bool
	out := captureStdout(t, func() { ok = instructPerson(t.TempDir()) })
	if ok {
		t.Error("with the prompt failing, instructPerson must tell init to stop")
	}
	if !strings.Contains(out, "I detected Gemini CLI") {
		t.Errorf("the detected AI is not named:\n%s", out)
	}
}

// openAI runs the tool in the project root, with the prompt as ONE argument (no shell).
func TestOpenAIRunsInTheRootWithoutAShell(t *testing.T) {
	t.Run("INWZN-B12: A detected AI is opened in the root with the prompt, or declined for the step-by-step", func(t *testing.T) {})
	root := t.TempDir()
	prompt := `a "quoted" (prompt) → with $HOME`
	script := testkit.FakeBin(t, t.TempDir(), "fake-ai", "#!/bin/sh\npwd > seen-dir\nprintf '%s' \"$1\" > seen-arg\n")
	if err := openAI(root, []string{script, prompt}); err != nil {
		t.Fatalf("openAI: %v", err)
	}
	arg, err := os.ReadFile(filepath.Join(root, "seen-arg"))
	if err != nil {
		t.Fatalf("the tool did not run in the root: %v", err)
	}
	if string(arg) != prompt {
		t.Errorf("the prompt arrived as %q, want it verbatim %q", arg, prompt)
	}
	if err := openAI(root, []string{filepath.Join(root, "does-not-exist")}); err == nil {
		t.Error("a missing tool must be an error")
	}
}

// breakAt reflows at the width and prefixes every continuation line, without losing words.
func TestBreakAtWrapsAtTheWidthAndKeepsEveryWord(t *testing.T) {
	t.Run("INWZN-B11: A person without a known AI gets the prompt to paste", func(t *testing.T) {})
	in := "alpha beta gamma delta epsilon zeta eta theta"
	got := breakAt(in, 12, "> ")
	lines := strings.Split(got, "\n")
	if len(lines) < 3 {
		t.Fatalf("expected several lines, got %q", got)
	}
	for i, l := range lines {
		body := l
		if i > 0 {
			if !strings.HasPrefix(l, "> ") {
				t.Errorf("line %d lacks the prefix: %q", i, l)
			}
			body = strings.TrimPrefix(l, "> ")
		}
		if len(body) > 12 {
			t.Errorf("line %d is wider than 12: %q", i, body)
		}
	}
	if strings.Join(strings.Fields(strings.ReplaceAll(got, "> ", "")), " ") != in {
		t.Errorf("words were lost or reordered: %q", got)
	}
	if breakAt("", 10, "> ") != "" {
		t.Error("empty text must stay empty")
	}
}

// Accepting the offer opens the detected AI in the project root with the interview
// prompt. A fake `gemini` on the PATH records how it was called.
func TestInstructPersonOpensTheDetectedAIWhenAccepted(t *testing.T) {
	t.Run("INWZN-B12: A detected AI is opened in the root with the prompt, or declined for the step-by-step", func(t *testing.T) {})
	clearAgentEnv(t)
	resetPromptError(t)
	t.Setenv("GEMINI_CLI", "1")
	t.Setenv("TERM", "dumb") // line-based prompt: reads os.Stdin, never a terminal
	bin := t.TempDir()
	testkit.FakeBin(t, bin, "gemini", "#!/bin/sh\nprintf '%s' \"$1\" > \"$PWD/opened-with\"\n")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	root := t.TempDir()

	var ok bool
	withStdin(t, "y\n", func() {
		captureStdout(t, func() { ok = instructPerson(root) })
	})
	if !ok {
		t.Error("opening the AI must let init go on")
	}
	b, err := os.ReadFile(filepath.Join(root, "opened-with"))
	if err != nil {
		t.Fatalf("the AI was not opened in the project root: %v", err)
	}
	if string(b) != initx.PromptDescobrir {
		t.Errorf("the AI got %q, want the interview prompt", b)
	}
}

// Declining the offer falls back to the step-by-step.
func TestInstructPersonDeclinedPrintsTheStepByStep(t *testing.T) {
	t.Run("INWZN-B12: A detected AI is opened in the root with the prompt, or declined for the step-by-step", func(t *testing.T) {})
	clearAgentEnv(t)
	resetPromptError(t)
	t.Setenv("GEMINI_CLI", "1")
	t.Setenv("TERM", "dumb")
	var ok bool
	var out string
	withStdin(t, "n\n", func() {
		out = captureStdout(t, func() { ok = instructPerson(t.TempDir()) })
	})
	if !ok || !strings.Contains(out, "To do it yourself") {
		t.Errorf("declined: ok = %v\n%s", ok, out)
	}
}

// --- init_git.go ---

// gitIdentityForTest gives git an author for THIS test only.
//
// Without it the `git commit` of `initGit` fails with "Author identity unknown" on any
// machine without a global `user.name`/`user.email` — and the tests passed on the machine of
// whoever wrote them and failed in CI. Found on the FIRST run of the CI workflow: three
// FAILs, all here.
//
// A test does not depend on its environment: one that only passes on a configured machine
// proves the machine, not the code.
//
// Through ENV and not `git config`: the variables hold only for this test's processes, so
// nothing is written to the temporary repository nor to the configuration of whoever runs
// it. `t.Setenv` restores them at the end, and fails if the test is parallel — the right
// guarantee, because touching shared env under parallelism is a race.
//
// The fix belongs in the TEST, not in `initGit`: requiring an identity is the product's
// right behaviour, and its error message guides the user. Injecting an author in the
// product would commit as someone nobody chose.
func gitIdentityForTest(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_AUTHOR_NAME", "anchors-test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@anchors.invalid")
	t.Setenv("GIT_COMMITTER_NAME", "anchors-test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@anchors.invalid")
}

// The path where the user ACCEPTS: the repository must come out really usable — with
// HEAD. A `git init` without a commit leaves `git log`/`git diff` with nothing to compare
// against, and half of Anchors would stay off believing it was solved.
func TestInitGitLeavesARepoWithHead(t *testing.T) {
	t.Run("INWZN-B05: Accepting git leaves a repository with HEAD", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	gitIdentityForTest(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".DS_Store"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := initGit(dir, initx.GitNaoIniciado); err != nil {
		t.Fatalf("initGit: %v", err)
	}

	if e := initx.DetectGit(dir, true); e != initx.GitPronto {
		t.Fatalf("after initializing, the state must be GitPronto, got %v", e)
	}
	out, err := runGit(dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("git log failed — there is no HEAD: %v (%s)", err, out)
	}
	if strings.TrimSpace(out) != initx.FirstCommitMessage {
		t.Errorf("commit subject = %q, want %q", out, initx.FirstCommitMessage)
	}
	b, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore was not seeded: %v", err)
	}
	if !strings.Contains(string(b), ".DS_Store") {
		t.Error("the seeded .gitignore does not cover .DS_Store")
	}
}

// A .gitignore that already exists is the USER's and outranks our default.
func TestInitGitKeepsAnExistingGitignore(t *testing.T) {
	t.Run("INWZN-B06: An existing gitignore is kept", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	gitIdentityForTest(t)
	dir := t.TempDir()
	mine := "# mine\n*.log\n"
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := initGit(dir, initx.GitNaoIniciado); err != nil {
		t.Fatalf("initGit: %v", err)
	}

	b, _ := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if string(b) != mine {
		t.Errorf("the user's .gitignore was overwritten: %q", string(b))
	}
}

// In a repository that exists without a commit only the commit is missing — another
// `git init` is unnecessary, and the end state is the same: HEAD exists.
func TestInitGitOnlyCommitsWhenTheRepoExists(t *testing.T) {
	t.Run("INWZN-B05: Accepting git leaves a repository with HEAD", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	gitIdentityForTest(t)
	dir := t.TempDir()
	if out, err := runGit(dir, "init"); err != nil {
		t.Fatalf("setup: %s", out)
	}

	if err := initGit(dir, initx.GitSemCommit); err != nil {
		t.Fatalf("initGit: %v", err)
	}
	if e := initx.DetectGit(dir, true); e != initx.GitPronto {
		t.Fatalf("end state = %v, want GitPronto", e)
	}
}

// --- gitStep: the first step of `anchors init` ---

// A machine without git gets the warning and no question: offering `git init` there would
// fail the moment it is accepted. And init goes on — git missing is a warning, not an error.
func TestGitStepWithoutGitWarnsAndProceeds(t *testing.T) {
	t.Run("INWZN-B03: Without git the step warns and the init goes on", func(t *testing.T) {})
	resetPromptError(t)
	t.Setenv("PATH", t.TempDir()) // no git on this PATH
	root := t.TempDir()

	var ok bool
	out := captureStdout(t, func() { ok = gitStep(root) })
	if !ok {
		t.Fatal("gitStep stopped init on a machine without git; it must only warn")
	}
	if !strings.Contains(out, initx.AvisoGit(initx.GitNaoInstalado)) {
		t.Errorf("the not-installed warning is missing:\n%s", out)
	}
	if erroDePrompt {
		t.Error("a prompt ran although there is nothing to offer without git")
	}
}

// A repository with a commit is ready: nothing is printed and nothing is asked.
func TestGitStepOnAReadyRepoIsSilent(t *testing.T) {
	t.Run("INWZN-B02: A ready repository makes the git step silent", func(t *testing.T) {})
	resetPromptError(t)
	root := newGitRepo(t)
	var ok bool
	out := captureStdout(t, func() { ok = gitStep(root) })
	if !ok || out != "" {
		t.Errorf("gitStep on a ready repo = %v, printed %q; want true and silence", ok, out)
	}
}

// A repository WITHOUT a commit is offered the first commit. When the offer cannot be
// asked, gitStep says so (false) and does not commit on its own.
func TestGitStepOnARepoWithoutCommitDoesNotActWithoutAnAnswer(t *testing.T) {
	t.Run("INWZN-X01: Git is never initialized nor committed without a yes", func(t *testing.T) {})
	skipIfTerminal(t)
	resetPromptError(t)
	isolateGit(t)
	root := t.TempDir()
	if out, err := runGit(root, "init", "-q"); err != nil {
		t.Fatalf("git init: %s", out)
	}

	var ok bool
	out := captureStdout(t, func() { ok = gitStep(root) })
	if ok {
		t.Error("with the prompt failing, gitStep must tell init to abort")
	}
	if !strings.Contains(out, initx.AvisoGit(initx.GitSemCommit)) {
		t.Errorf("the no-commit warning is missing:\n%s", out)
	}
	if _, err := runGit(root, "rev-parse", "HEAD"); err == nil {
		t.Error("a commit was made although nobody accepted")
	}
}

// Without an identity the commit fails, and the error names the fix instead of pasting
// git's long message.
func TestInitGitWithoutIdentityNamesTheFix(t *testing.T) {
	t.Run("INWZN-B07: A failed git initialization is reported and names the fix", func(t *testing.T) {})
	withoutGitIdentity(t)
	root := t.TempDir()

	var err error
	captureStdout(t, func() { err = initGit(root, initx.GitNaoIniciado) })
	if err == nil {
		t.Fatal("the commit cannot succeed without an identity")
	}
	msg := err.Error()
	if !strings.Contains(msg, "git does not know who you are") ||
		!strings.Contains(msg, `git config --global user.email`) {
		t.Errorf("the error does not name the fix: %v", err)
	}
	if strings.Count(msg, "\n") > 4 {
		t.Errorf("only the first line of git's output belongs in the detail: %v", err)
	}
}

func TestFirstLineCutsAtTheFirstNewline(t *testing.T) {
	for in, want := range map[string]string{
		"one\ntwo\nthree": "one",
		"single":          "single",
		"":                "",
	} {
		if got := firstLine(in); got != want {
			t.Errorf("firstLine(%q) = %q, want %q", in, got, want)
		}
	}
}

// answerGitOffer runs gitStep on a directory outside git, answering the offer with
// `reply`. TERM=dumb puts the prompt in its line-based mode, which reads os.Stdin instead
// of opening a terminal — so the answer is deterministic.
func answerGitOffer(t *testing.T, reply string) (dir string, ok bool, out string) {
	t.Helper()
	resetPromptError(t)
	isolateGit(t)
	t.Setenv("TERM", "dumb")
	dir = t.TempDir()
	withStdin(t, reply, func() {
		out = captureStdout(t, func() { ok = gitStep(dir) })
	})
	return dir, ok, out
}

// Declining is legitimate and init goes on — but it names NOW what stays incomplete.
func TestGitStepDeclinedNamesWhatStaysOff(t *testing.T) {
	t.Run("INWZN-B04: Declining git names what stays off", func(t *testing.T) {})
	dir, ok, out := answerGitOffer(t, "n\n")
	if !ok {
		t.Error("declining git must not stop init")
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		t.Error("a repository was created although the offer was declined")
	}
	for _, want := range []string{"proceeding without git", "coverage --diff", "install-hooks", "`git init` later"} {
		if !strings.Contains(out, want) {
			t.Errorf("the consequence %q is not named:\n%s", want, out)
		}
	}
}

// Accepting leaves a repository with HEAD, the .gitignore seeded, and init going on.
func TestGitStepAcceptedLeavesARepoWithHead(t *testing.T) {
	t.Run("INWZN-B05: Accepting git leaves a repository with HEAD", func(t *testing.T) {})
	dir, ok, out := answerGitOffer(t, "y\n")
	if !ok {
		t.Error("accepting git must not stop init")
	}
	if e := initx.DetectGit(dir, true); e != initx.GitPronto {
		t.Fatalf("after accepting, the state is %v, want GitPronto\n%s", e, out)
	}
	if !strings.Contains(out, "✓ first commit") || !strings.Contains(out, "✓ .gitignore seeded") {
		t.Errorf("the output does not report what was done:\n%s", out)
	}
}

// A failure while initializing does not stop init: it is reported with its cause (here,
// git without an identity).
func TestGitStepReportsAFailedInitialization(t *testing.T) {
	t.Run("INWZN-B07: A failed git initialization is reported and names the fix", func(t *testing.T) {})
	resetPromptError(t)
	withoutGitIdentity(t)
	t.Setenv("TERM", "dumb")
	dir := t.TempDir()
	var ok bool
	var out string
	withStdin(t, "y\n", func() {
		out = captureStdout(t, func() { ok = gitStep(dir) })
	})
	if !ok {
		t.Error("a failed initialization must not stop init")
	}
	if !strings.Contains(out, "could not initialize git") || !strings.Contains(out, "git does not know who you are") {
		t.Errorf("the failure and its cause are not reported:\n%s", out)
	}
}

// withoutGitIdentity leaves git with no author: no global config but `useConfigOnly`, and
// none of the identity variables — so a commit fails the way it does on a fresh machine.
func withoutGitIdentity(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	global := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(global, []byte("[user]\n\tuseConfigOnly = true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", global)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL", "EMAIL"} {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}
}

// Whoever hits the TTY guard — an agent, a pipe, CI — needs to know there is a way out.
// Without naming the non-interactive mode the command looks like a dead end: "it is
// interactive and there is no terminal" describes the problem and offers nothing.
func TestNoTTYErrorOffersTheNonInteractiveMode(t *testing.T) {
	t.Run("INWZN-E01: The no-terminal error offers the non-interactive mode", func(t *testing.T) {})
	msg := errNoTTY("Nothing was written.").Error()

	if !strings.Contains(msg, "--non-interactive") {
		t.Errorf("the message must name the way out:\n%s", msg)
	}
	if !strings.Contains(msg, "Nothing was written.") {
		t.Errorf("the caller's context must appear:\n%s", msg)
	}
	// BOTH calls of the same flag: the one with no answers (which asks) and the one with
	// answers (which applies). Showing only one would hide half of the contract — the half
	// the reader needs first.
	if strings.Count(msg, "--non-interactive") < 2 {
		t.Errorf("both calls must appear:\n%s", msg)
	}
	if !strings.Contains(msg, "--artifacts") {
		t.Errorf("the second call needs an example answer:\n%s", msg)
	}
}

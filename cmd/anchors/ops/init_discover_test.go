package ops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

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
	root := t.TempDir()
	prompt := `a "quoted" (prompt) → with $HOME`
	script := writeFile(t, t.TempDir(), "fake-ai", "#!/bin/sh\npwd > seen-dir\nprintf '%s' \"$1\" > seen-arg\n")
	if err := os.Chmod(script, 0o755); err != nil {
		t.Fatal(err)
	}
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
	clearAgentEnv(t)
	resetPromptError(t)
	t.Setenv("GEMINI_CLI", "1")
	t.Setenv("TERM", "dumb") // line-based prompt: reads os.Stdin, never a terminal
	bin := t.TempDir()
	fake := writeFile(t, bin, "gemini", "#!/bin/sh\nprintf '%s' \"$1\" > \"$PWD/opened-with\"\n")
	if err := os.Chmod(fake, 0o755); err != nil {
		t.Fatal(err)
	}
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

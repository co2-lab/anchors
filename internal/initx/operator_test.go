// @anchors
//   ref: OPDTP

package initx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeEnv simulates environment variables without touching the real ones — detection must be
// testable on both sides, and the machine running the suite has its own.
func fakeEnv(pairs map[string]string) func(string) string {
	return func(k string) string { return pairs[k] }
}

// The distinction that shapes the bootstrap: what Anchors SAYS depends on who reads it. An
// agent variable is an explicit declaration and outweighs the absence of a terminal — which
// any pipe also produces.
func TestOperatorAgentVariableBeatsTheMissingTerminal(t *testing.T) {
	t.Run("OPDTP-B01: A known agent variable makes the operator an AI even with a terminal", func(t *testing.T) {})
	t.Run("OPDTP-B03: No terminal and no agent variable is an AI", func(t *testing.T) {})
	t.Run("OPDTP-B04: A terminal and no agent variable is a person", func(t *testing.T) {})
	withAgent := fakeEnv(map[string]string{"CLAUDE_CODE_ENTRYPOINT": "cli"})

	// Even WITH a terminal, the variable decides: an AI can run in a PTY.
	if o := DetectOperator(true, withAgent); o != OperadorIA {
		t.Error("with an agent variable declared it is an AI — even with a terminal")
	}
	// And with nothing, a terminal is a person.
	if o := DetectOperator(true, fakeEnv(nil)); o != OperadorHumano {
		t.Error("a terminal and no agent variable: it is a person")
	}
	// Without a terminal what is left is a pipe, CI or an agent — none is someone typing.
	if o := DetectOperator(false, fakeEnv(nil)); o != OperadorIA {
		t.Error("without a terminal there is no person reading the output and answering")
	}
}

// `AI_AGENT` is not in the list of names, and still declares that nobody is typing.
func TestOperatorGenericAgentVariableIsAnAI(t *testing.T) {
	t.Run("OPDTP-B02: The generic AI_AGENT variable makes the operator an AI even with a terminal", func(t *testing.T) {})
	env := fakeEnv(map[string]string{"AI_AGENT": "1"})
	if o := DetectOperator(true, env); o != OperadorIA {
		t.Error("AI_AGENT declares an agent even with a terminal")
	}
	if n := AgentName(env); n != "" {
		t.Errorf("AI_AGENT names no tool, got %q", n)
	}
}

// Naming the tool is what allows saying "shall I open Claude Code?" instead of "open your
// AI" — the difference between an actionable offer and advice.
func TestOperatorNamesTheDetectedTool(t *testing.T) {
	t.Run("OPDTP-B05: The tool is named after the known variable that is set", func(t *testing.T) {})
	if n := AgentName(fakeEnv(map[string]string{"CURSOR_TRACE_ID": "x"})); n != "Cursor" {
		t.Errorf("expected Cursor, got %q", n)
	}
	if n := AgentName(fakeEnv(nil)); n != "" {
		t.Errorf("with no known variable there is no name to give, got %q", n)
	}
}

// Two known variables set at once: the variable name that sorts first wins, every time.
func TestOperatorNameIsStableWithSeveralAgents(t *testing.T) {
	t.Run("OPDTP-B06: Several known variables resolve by the first variable name in order", func(t *testing.T) {})
	env := fakeEnv(map[string]string{"CLAUDECODE": "1", "AIDER_MODEL": "x"})
	for i := 0; i < 20; i++ {
		if n := AgentName(env); n != "Aider" {
			t.Fatalf("run %d: expected Aider (AIDER_MODEL sorts first), got %q", i, n)
		}
	}
}

// The command is an ARGV, not a shell line: the prompt has quotes, parentheses and arrows,
// and passing it through `sh -c` would leave escaping as the only thing between the text and
// the interpreter.
func TestOperatorOpenCommandIsArgvNotShell(t *testing.T) {
	t.Run("OPDTP-B09: Only the tools with a stable command line get a command", func(t *testing.T) {})
	t.Run("OPDTP-I01: The whole prompt is one argument of the command", func(t *testing.T) {})
	t.Run("OPDTP-X01: No command is offered for a tool without a stable command line", func(t *testing.T) {})
	want := map[string][]string{
		"CLAUDE_CODE_ENTRYPOINT": {"claude"},
		"GEMINI_CLI":             {"gemini"},
		"AIDER_MODEL":            {"aider", "--message"},
	}
	for variable, prefix := range want {
		argv := CommandToOpenAI(fakeEnv(map[string]string{variable: "x"}))
		if len(argv) != len(prefix)+1 {
			t.Fatalf("%s: expected %v plus the prompt, got %v", variable, prefix, argv)
		}
		for i, p := range prefix {
			if argv[i] != p {
				t.Errorf("%s: argument %d = %q, want %q", variable, i, argv[i], p)
			}
		}
		// The whole prompt in one argument — no shell quotes glued to it.
		if argv[len(argv)-1] != PromptDescobrir {
			t.Errorf("%s: the last argument is not the prompt verbatim: %q", variable, argv[len(argv)-1])
		}
	}
	// Unknown or unstable tool: nothing to offer. A command that does not exist is worse
	// than no offer.
	for _, env := range []map[string]string{nil, {"CURSOR_TRACE_ID": "x"}, {"CODEX_SANDBOX": "x"}} {
		if argv := CommandToOpenAI(fakeEnv(env)); argv != nil {
			t.Errorf("%v: no command should be offered, got %v", env, argv)
		}
	}
}

// The prompt tells to READ the guide, it does not summarise it: `anchors guide project` is
// the source of truth of the interview, and a summary would go stale silently the first
// time it changed.
func TestOperatorPromptPointsAtTheGuideInsteadOfSummarisingIt(t *testing.T) {
	t.Run("OPDTP-B10: The discovery prompt points at the guide and back at init", func(t *testing.T) {})
	if !strings.Contains(PromptDescobrir, "anchors guide project") {
		t.Error("the prompt must tell to run the guide — it is the ruler")
	}
	if !strings.Contains(PromptDescobrir, "anchors init") {
		t.Error("the prompt must close the loop pointing back at init")
	}
}

// The phase is needed only when there is NOTHING to infer from. A PROJECT.md is the proof
// that it ran; code on disk makes it unnecessary (init infers from there).
func TestOperatorDiscoveryOnlyWhenEmpty(t *testing.T) {
	t.Run("OPDTP-B07: The discovery phase is due only when nothing is found and nothing is described", func(t *testing.T) {})
	dir := t.TempDir()

	if !PrecisaDescobrir(dir, &Proposal{}) {
		t.Error("with no PROJECT.md and no code, the phase has not happened")
	}
	if !PrecisaDescobrir(dir, nil) {
		t.Error("a nil proposal must not mask the emptiness")
	}
	// Something on disk: there is something to infer from.
	for name, p := range map[string]*Proposal{
		"code":    {CodeDirs: []string{"src"}},
		"spec":    {HasSpecMD: true},
		"feature": {HasFeature: true},
		"test":    {HasTest: true},
	} {
		if PrecisaDescobrir(dir, p) {
			t.Errorf("with a %s found, init infers from the disk — the phase is not needed", name)
		}
	}
	// PROJECT.md present: the phase already ran, even with no code yet (which is exactly
	// the state it MUST leave the project in).
	if err := os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte("# P\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if PrecisaDescobrir(dir, &Proposal{}) {
		t.Error("PROJECT.md is the proof that the interview happened")
	}
}

// A project that already used `project.md` must not be sent to redo the interview.
func TestOperatorProjectDescriptionSpellings(t *testing.T) {
	t.Run("OPDTP-B08: The project description is recognised under its three spellings", func(t *testing.T) {})
	for _, name := range []string{"project.md", "Project.md"} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, name), []byte("# P\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		// Only a case-sensitive disk tells the spellings apart; on a case-insensitive one
		// every spelling is found anyway, and the check still holds.
		if !HasProjectMD(dir) {
			t.Errorf("%s was not recognised", name)
		}
	}
	if HasProjectMD(t.TempDir()) {
		t.Error("an empty folder has no project description")
	}
}

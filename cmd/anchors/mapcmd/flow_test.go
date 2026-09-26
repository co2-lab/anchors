package mapcmd

import (
	"strings"
	"testing"
)

// A worker flow in miniature: a step that fits an action and routes two of its three
// results, a step with a conditional exit, a terminal, a step stuck with no exit, and a
// step that fits an action nobody wrote.
const workerFlow = "# Worker\n\n" +
	"### WORKR-P01 — PULL: take the next card\n\nFits: `ACPUL`\n\nResults:\n" +
	"- `ACPUL-R01` GOT → `WORKR-P02`\n- `ACPUL-R02` EMPTY → `WORKR-P03`\n\n" +
	"### WORKR-P02 — WORK: do the \"card\"\n\nFits: `ACGHO`\n\nExits:\n" +
	"- `WORKR-P03` when the gates are green\n- `WORKR-P04` when a gate is red\n\n" +
	"### WORKR-P03 — DONE: the work ends\n\n> @terminal\n\n" +
	"### WORKR-P04 — STUCK: nobody knows where to go\n"

const pullAction = "# pull\n\n" +
	"### ACPUL-R01 — GOT: a card was served\n\nSugere: `anchors work`, and read the card.\n\n" +
	"### ACPUL-R02 — EMPTY: the queue is empty\n\n" +
	"### ACPUL-R03 — BROKEN: nobody routes this one\n"

func flowProject(t *testing.T) string {
	t.Helper()
	root := fixtureProject(t)
	writeProjectFile(t, root, "flows/worker.flow.md", workerFlow)
	writeProjectFile(t, root, "flows/actions/pull.action.md", pullAction)
	return root
}

// The build writes the flow into the map, counts steps and results, and names the result
// no flow routes.
func TestFlowBuild_writesTheFlowAndNamesTheHole(t *testing.T) {
	root := flowProject(t)
	out := runCmd(t, newFlowCmd(), "build", "--root", root)

	if !strings.Contains(out, "flow built: 4 step(s), 3 result(s)") {
		t.Errorf("unexpected counts:\n%s", out)
	}
	if !strings.Contains(out, "1 result(s) no flow handles") || !strings.Contains(out, "ACPUL-R03") {
		t.Errorf("the unrouted result was not named:\n%s", out)
	}
	g := loadMap(t, root)
	if g.Flow == nil || len(g.Flow.States) != 7 {
		t.Errorf("the flow graph did not reach the map: %+v", g.Flow)
	}
	// The map's own nodes survive: the flow is added, nothing replaced.
	if len(g.Nodes) != 5 {
		t.Errorf("writing the flow changed the nodes: %d", len(g.Nodes))
	}
}

func TestFlowBuild_withoutFlowsOrMap(t *testing.T) {
	root := fixtureProject(t)
	out := runCmd(t, newFlowCmd(), "build", "--root", root)
	if !strings.Contains(out, "no flow declared") {
		t.Errorf("a project with no flows/ is not an error, and says so:\n%s", out)
	}

	bare := t.TempDir()
	writeProjectFile(t, bare, "flows/worker.flow.md", workerFlow)
	if _, err := runCmdErr(newFlowCmd(), t, "build", "--root", bare); err == nil ||
		!strings.Contains(err.Error(), "build the map first") {
		t.Errorf("a flow with no map must ask for the map first, got %v", err)
	}
}

// `flow next` answers the valid exits of a step, with the suggestion the result carries.
func TestFlowNext_theValidExits(t *testing.T) {
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)

	out := runCmd(t, newFlowCmd(), "next", "workr-p01", "--root", root)
	for _, want := range []string{
		"WORKR-P01 — PULL: take the next card",
		"fits: ACPUL (anchors pull)",
		"valid exits (2):",
		"ACPUL-R01    → WORKR-P02 (WORK: do the \"card\")",
		"↳ `anchors work`",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	work := runCmd(t, newFlowCmd(), "next", "WORKR-P02", "--root", root)
	if !strings.Contains(work, "⚠ no action file declares it") {
		t.Errorf("a step fitting an unwritten action must say so:\n%s", work)
	}
	if !strings.Contains(work, "—            → WORKR-P03") || !strings.Contains(work, "when the gates are green") {
		t.Errorf("a conditional exit shows the dash label and its condition:\n%s", work)
	}

	if done := runCmd(t, newFlowCmd(), "next", "WORKR-P03", "--root", root); !strings.Contains(done, "▣ terminal") {
		t.Errorf("the terminal step must say the work ends:\n%s", done)
	}
	if stuck := runCmd(t, newFlowCmd(), "next", "WORKR-P04", "--root", root); !strings.Contains(stuck, "whoever arrives here is stuck") {
		t.Errorf("a step with no exit and no @terminal is the defect to show:\n%s", stuck)
	}
	if _, err := runCmdErr(newFlowCmd(), t, "next", "WORKR-P99", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "not in the flow graph") {
		t.Errorf("an unknown step: got %v", err)
	}
}

// The text drawing: steps in file order, results hidden, exits sorted, terminal marked.
func TestFlowShow_drawsTheFlow(t *testing.T) {
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)

	out := runCmd(t, newFlowCmd(), "show", "worker", "--root", root)
	for _, want := range []string{
		"● WORKR-P01  PULL: take the next card   [ACPUL]",
		"├─ ACPUL-R01    → WORKR-P02",
		"└─ ACPUL-R02    → WORKR-P03",
		"└─ sempre       → WORKR-P04  when a gate is red",
		"▣ WORKR-P03  DONE: the work ends",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Index(out, "WORKR-P01") > strings.Index(out, "WORKR-P04") {
		t.Errorf("the steps must keep the file's order:\n%s", out)
	}
	if strings.Contains(out, "● ACPUL") {
		t.Errorf("results are not steps, and must not be drawn as one:\n%s", out)
	}

	if _, err := runCmdErr(newFlowCmd(), t, "show", "nonexistent", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "available:") {
		t.Errorf("an unknown flow must list the available ones, got %v", err)
	}
	if _, err := runCmdErr(newFlowCmd(), t, "show", "--root", fixtureProject(t)); err == nil ||
		!strings.Contains(err.Error(), "anchors flow build") {
		t.Errorf("a map without a flow must ask for `flow build`, got %v", err)
	}
}

// The mermaid drawing: native `\n` breaks, the stadium shape for terminals, labels from
// the result's short name, and the entry and terminals highlighted.
func TestFlowShow_mermaid(t *testing.T) {
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)

	out := runCmd(t, newFlowCmd(), "show", "--mermaid", "--direction", "lr", "--root", root)
	for _, want := range []string{
		"flowchart LR",
		`WORKR_P01["PULL: take the next card\nACPUL"]`,
		`WORKR_P02["WORK: do the 'card'\nACGHO"]`,
		`WORKR_P03(["DONE: the work ends"])`,
		"WORKR_P01 -->|GOT| WORKR_P02",
		"WORKR_P01 -->|EMPTY| WORKR_P03",
		"WORKR_P02 -->|when a gate is red| WORKR_P04",
		"class WORKR_P01 inicio",
		"class WORKR_P03 fim",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// An unknown direction falls back to top-down instead of emitting invalid mermaid.
	td := runCmd(t, newFlowCmd(), "show", "--mermaid", "--direction", "sideways", "--root", root)
	if !strings.HasPrefix(td, "flowchart TD\n") {
		t.Errorf("an unknown direction must fall back to TD:\n%s", td)
	}
}

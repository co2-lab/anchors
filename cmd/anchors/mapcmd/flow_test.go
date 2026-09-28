package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A worker flow in miniature: a step that fits an action and routes two of its three
// results, a step with a conditional exit, a terminal, a step stuck with no exit, and a
// step that fits an action nobody wrote.
const workerFlow = "# Worker\n\n" +
	"### WORKR-T01 — PULL: take the next card\n\nFits: `ACPUL`\n\nResults:\n" +
	"- `ACPUL-O01` GOT → `WORKR-T02`\n- `ACPUL-O02` EMPTY → `WORKR-T03`\n\n" +
	"### WORKR-T02 — WORK: do the \"card\"\n\nFits: `ACGHO`\n\nExits:\n" +
	"- `WORKR-T03` when the gates are green\n- `WORKR-T04` when a gate is red\n\n" +
	"### WORKR-T03 — DONE: the work ends\n\n> @terminal\n\n" +
	"### WORKR-T04 — STUCK: nobody knows where to go\n"

const pullAction = "# pull\n\n" +
	"### ACPUL-O01 — GOT: a card was served\n\nSugere: `anchors work`, and read the card.\n\n" +
	"### ACPUL-O02 — EMPTY: the queue is empty\n\n" +
	"### ACPUL-O03 — BROKEN: nobody routes this one\n"

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
	t.Run("FLWOX-B02: Building the flow writes it into the map and keeps the map nodes", func(t *testing.T) {})
	t.Run("FLWOX-B03: A result no flow routes is named as a warning and the build still succeeds", func(t *testing.T) {})
	root := flowProject(t)
	out := runCmd(t, newFlowCmd(), "build", "--root", root)

	if !strings.Contains(out, "flow built: 4 step(s), 3 result(s)") {
		t.Errorf("unexpected counts:\n%s", out)
	}
	if !strings.Contains(out, "1 result(s) no flow handles") || !strings.Contains(out, "ACPUL-O03") {
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
	t.Run("FLWOX-B01: A project with no flow declared is told so and the build is not an error", func(t *testing.T) {})
	t.Run("FLWOX-E01: Building a flow without a map is refused and no map is created", func(t *testing.T) {})
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
	if _, err := os.Stat(filepath.Join(bare, mapx.DefaultPath)); !os.IsNotExist(err) {
		t.Errorf("a refused flow build must not create a map (stat: %v)", err)
	}
}

// `flow next` answers the valid exits of a step, with the suggestion the result carries.
func TestFlowNext_theValidExits(t *testing.T) {
	t.Run("FLWOX-B04: The next command lists the valid exits of a step with the suggestion of each result", func(t *testing.T) {})
	t.Run("FLWOX-B05: A step that fits an action no file declares is flagged", func(t *testing.T) {})
	t.Run("FLWOX-B06: A terminal step says the work ends and a step with no exit and no terminal mark says it is stuck", func(t *testing.T) {})
	t.Run("FLWOX-E02: Asking the exits of a step that is not in the flow graph fails", func(t *testing.T) {})
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)

	out := runCmd(t, newFlowCmd(), "next", "workr-t01", "--root", root)
	for _, want := range []string{
		"WORKR-T01 — PULL: take the next card",
		"fits: ACPUL (anchors pull)",
		"valid exits (2):",
		"ACPUL-O01    → WORKR-T02 (WORK: do the \"card\")",
		"↳ `anchors work`",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	work := runCmd(t, newFlowCmd(), "next", "WORKR-T02", "--root", root)
	if !strings.Contains(work, "⚠ no action file declares it") {
		t.Errorf("a step fitting an unwritten action must say so:\n%s", work)
	}
	if !strings.Contains(work, "—            → WORKR-T03") || !strings.Contains(work, "when the gates are green") {
		t.Errorf("a conditional exit shows the dash label and its condition:\n%s", work)
	}

	if done := runCmd(t, newFlowCmd(), "next", "WORKR-T03", "--root", root); !strings.Contains(done, "▣ terminal") {
		t.Errorf("the terminal step must say the work ends:\n%s", done)
	}
	if stuck := runCmd(t, newFlowCmd(), "next", "WORKR-T04", "--root", root); !strings.Contains(stuck, "whoever arrives here is stuck") {
		t.Errorf("a step with no exit and no @terminal is the defect to show:\n%s", stuck)
	}
	if _, err := runCmdErr(newFlowCmd(), t, "next", "WORKR-T99", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "not in the flow graph") {
		t.Errorf("an unknown step: got %v", err)
	}
}

// The text drawing: steps in file order, results hidden, exits sorted, terminal marked.
func TestFlowShow_drawsTheFlow(t *testing.T) {
	t.Run("FLWOX-B07: The text drawing keeps the file order, hides the results and sorts each step exits", func(t *testing.T) {})
	t.Run("FLWOX-E03: Showing a flow no name matches fails listing the available flows", func(t *testing.T) {})
	t.Run("FLWOX-E04: Showing or navigating a map with no flow asks for the flow build", func(t *testing.T) {})
	t.Run("FLWOX-I01: What the flow build writes is what show and next read", func(t *testing.T) {})
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)

	out := runCmd(t, newFlowCmd(), "show", "worker", "--root", root)
	for _, want := range []string{
		"● WORKR-T01  PULL: take the next card   [ACPUL]",
		"├─ ACPUL-O01    → WORKR-T02",
		"└─ ACPUL-O02    → WORKR-T03",
		"└─ sempre       → WORKR-T04  when a gate is red",
		"▣ WORKR-T03  DONE: the work ends",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Index(out, "WORKR-T01") > strings.Index(out, "WORKR-T04") {
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
	if _, err := runCmdErr(newFlowCmd(), t, "next", "WORKR-T01", "--root", fixtureProject(t)); err == nil ||
		!strings.Contains(err.Error(), "anchors flow build") {
		t.Errorf("navigating a map without a flow must ask for `flow build`, got %v", err)
	}
	if part := runCmd(t, newFlowCmd(), "show", "ork", "--root", root); !strings.Contains(part, "WORKR-T01") {
		t.Errorf("a part of the flow's name must select it:\n%s", part)
	}
}

// The mermaid drawing: native `\n` breaks, the stadium shape for terminals, labels from
// the result's short name, and the entry and terminals highlighted.
func TestFlowShow_mermaid(t *testing.T) {
	t.Run("FLWOX-B08: The mermaid drawing emits valid identifiers, native line breaks, the stadium shape for terminals and highlights the entry and the ends", func(t *testing.T) {})
	t.Run("FLWOX-B09: An unknown mermaid direction falls back to top-down", func(t *testing.T) {})
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)

	out := runCmd(t, newFlowCmd(), "show", "--mermaid", "--direction", "lr", "--root", root)
	for _, want := range []string{
		"flowchart LR",
		`WORKR_T01["PULL: take the next card\nACPUL"]`,
		`WORKR_T02["WORK: do the 'card'\nACGHO"]`,
		`WORKR_T03(["DONE: the work ends"])`,
		"WORKR_T01 -->|GOT| WORKR_T02",
		"WORKR_T01 -->|EMPTY| WORKR_T03",
		"WORKR_T02 -->|when a gate is red| WORKR_T04",
		"class WORKR_T01 inicio",
		"class WORKR_T03 fim",
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
	for _, dir := range []string{"rl", "BT"} {
		got := runCmd(t, newFlowCmd(), "show", "--mermaid", "--direction", dir, "--root", root)
		if !strings.HasPrefix(got, "flowchart "+strings.ToUpper(dir)+"\n") {
			t.Errorf("direction %s must be kept:\n%s", dir, got)
		}
	}
}

// Only `flow build` writes: drawing and navigating read the map and leave it as it was.
func TestFlow_showAndNextWriteNothing(t *testing.T) {
	t.Run("FLWOX-X01: Showing and navigating the flow leave the map unchanged", func(t *testing.T) {})
	root := flowProject(t)
	runCmd(t, newFlowCmd(), "build", "--root", root)
	path := filepath.Join(root, mapx.DefaultPath)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	runCmd(t, newFlowCmd(), "show", "--root", root)
	runCmd(t, newFlowCmd(), "show", "--mermaid", "--root", root)
	runCmd(t, newFlowCmd(), "next", "WORKR-T01", "--root", root)
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Error("show or next changed the map")
	}
}

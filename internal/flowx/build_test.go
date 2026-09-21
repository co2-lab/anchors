package flowx

import (
	"os"
	"path/filepath"
	"testing"
)

// The flow that originated the feature, in miniature: an unblocking rotation.
const destravar = `<!-- @anchors
  code: DSTRV
  updated_at: 2026-09-21
-->
# Unblocking

### DSTRV-N01 — OPEN: the problem is described and nobody attacked it

Exits:
- ` + "`DSTRV-N02`" + ` always — measuring is the first attack, and the cheapest

### DSTRV-N02 — ATTACKING: one approach is under way

Exits:
- ` + "`DSTRV-N03`" + ` when the measurement isolated the term and there is a number
- ` + "`DSTRV-N04`" + ` when TWO attacks already failed on this approach

### DSTRV-N03 — MEASURED: the term was isolated and there is a number

Exits:
- ` + "`DSTRV-N05`" + ` when the measurement points at a concrete target

### DSTRV-N04 — ROTATION: the approach is exhausted, another has to step in

Exits:
- ` + "`DSTRV-N02`" + ` when there is still an approach nobody tried

### DSTRV-N05 — FIXED: the fix is in place and measured

> @terminal
`

func projectWithFlow(t *testing.T, content string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, Dir, "unblock"+FlowSuffix), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// The exit→state association is POSITIONAL: an exit belongs to the last state declared
// above it. That is what lets the flow be written the way it reads.
func TestBuild_exitsBelongToTheStateAbove(t *testing.T) {
	g, err := Build(projectWithFlow(t, destravar))
	if err != nil || g == nil {
		t.Fatalf("build failed: %v", err)
	}
	if len(g.States) != 5 {
		t.Fatalf("expected 5 states, got %d", len(g.States))
	}
	got := Next(g, "DSTRV-N02")
	if len(got) != 2 {
		t.Fatalf("N02 has two exits, got %d: %+v", len(got), got)
	}
	// The condition travels with the exit: it is what whoever works reads to choose.
	if got[1].To != "DSTRV-N04" || got[1].When == "" {
		t.Errorf("the exit must carry its condition: %+v", got[1])
	}
}

// This is the inversion the whole feature exists for: from a state, the answer is the set
// of valid exits — not the catalogue of every approach with its "when this one is right".
func TestNext_answersTheSetAndNotTheCatalogue(t *testing.T) {
	g, _ := Build(projectWithFlow(t, destravar))
	if n := len(Next(g, "DSTRV-N01")); n != 1 {
		t.Errorf("from OPEN there is one way out, got %d", n)
	}
	if n := len(Next(g, "DSTRV-N05")); n != 0 {
		t.Errorf("a terminal state has no exit, got %d", n)
	}
}

// Terminal is DECLARED, never deduced from "has no exit": a state with no exit may be the
// end of the work or an oversight, and only the author knows which.
func TestBuild_terminalIsDeclared(t *testing.T) {
	g, _ := Build(projectWithFlow(t, destravar))
	s, ok := StateByCode(g, "DSTRV-N05")
	if !ok || !s.Terminal {
		t.Errorf("N05 declares @terminal: %+v", s)
	}
	if s, _ := StateByCode(g, "DSTRV-N03"); s.Terminal {
		t.Error("N03 has an exit and does not declare @terminal — it is not terminal")
	}
}

// The finding a hand-written flow produces easily: the state is created with ITS exits,
// and the exit that ARRIVES at it is forgotten. It stays written, documented, unreachable.
func TestUnreachable_findsTheStateNobodyArrivesAt(t *testing.T) {
	orphan := destravar + "\n### DSTRV-N09 — nobody points here\n\n> @terminal\n"
	g, _ := Build(projectWithFlow(t, orphan))
	got := Unreachable(g)
	if len(got) != 1 || got[0].Code != "DSTRV-N09" {
		t.Errorf("expected exactly N09 unreachable, got %+v", got)
	}
}

// A broken exit is worse than an absent one: it LOOKS like a path, and whoever follows it
// arrives nowhere.
func TestDangling_findsTheExitThatLeadsNowhere(t *testing.T) {
	broken := destravar + "\n### DSTRV-N08 — a state with a broken exit\n\nExits:\n- `DSTRV-N77` when whatever\n"
	g, _ := Build(projectWithFlow(t, broken))
	got := Dangling(g)
	if len(got) != 1 || got[0].To != "DSTRV-N77" {
		t.Errorf("expected exactly the exit to N77, got %+v", got)
	}
}

// A project with no `flows/` has not declared any flow — that is not an error.
func TestBuild_projectWithoutFlowsIsNotAnError(t *testing.T) {
	g, err := Build(t.TempDir())
	if err != nil || g != nil {
		t.Errorf("expected nil graph and no error, got %+v / %v", g, err)
	}
}

// The PIECE shape is what makes this finding possible, and the linear drawing hid it: an
// action declares everything it can answer, and a flow that ignores one of those answers
// leaves a hole — whoever gets that result has nowhere to go, and improvises.
func TestUnhandled_findsTheResultNoFlowRoutes(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ActionsDir), 0o755)
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"### FLOWX-P01 — one step\n\nFits: `ACTST`\n\nResults:\n- `ACTST-R01` FINE → `FLOWX-P02`\n\n"+
			"### FLOWX-P02 — the end\n\n> @terminal\n"), 0o644)
	os.WriteFile(filepath.Join(root, ActionsDir, "a"+ActionSuffix), []byte(
		"### ACTST-R01 — FINE: it worked\n\n### ACTST-R02 — BROKEN: nobody routes this one\n"), 0o644)

	g, err := Build(root)
	if err != nil || g == nil {
		t.Fatalf("build: %v", err)
	}
	got := Unhandled(g)
	if len(got) != 1 || got[0].Code != "ACTST-R02" {
		t.Errorf("expected exactly ACTST-R02 unhandled, got %+v", got)
	}
	// A result is DECLARED by its action, never ARRIVED at: asking it "who reaches you?"
	// would accuse every result a flow happens not to route.
	for _, s := range Unreachable(g) {
		if IsResult(s.Code) {
			t.Errorf("a result must not be charged as unreachable: %s", s.Code)
		}
	}
}

// `Fits:` is what turns the flow into assembly rather than redrawing: the step names the
// piece, and the piece declares its own results. Without it every flow would repeat the
// description of `map build`, and the copies would diverge at the first change.
func TestBuild_readsThePieceEachStepFits(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, Dir), 0o755)
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"### FLOWX-P01 — a step\n\nFits: `ACMAP`\n\n### FLOWX-P02 — no piece\n\n> @terminal\n"), 0o644)
	g, _ := Build(root)
	s, _ := StateByCode(g, "FLOWX-P01")
	if s.Fits != "ACMAP" {
		t.Errorf("expected the step to fit ACMAP, got %q", s.Fits)
	}
	if s2, _ := StateByCode(g, "FLOWX-P02"); s2.Fits != "" {
		t.Errorf("a step that fits nothing must carry no piece, got %q", s2.Fits)
	}
}

// The THIRD category: a result that neither generates work by itself nor ends the matter —
// it SUGGESTS what to do, and the choice belongs to whoever works.
//
// Measured in the message catalog: 18 occurrences of "run `anchors <command>`" inside gate
// verdicts. Each is a flow edge hidden in prose, where it depends on somebody reading and
// remembering.
func TestBuild_readsTheSuggestedReaction(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ActionsDir), 0o755)
	os.MkdirAll(filepath.Join(root, Dir), 0o755)
	os.WriteFile(filepath.Join(root, ActionsDir, "a"+ActionSuffix), []byte(
		"### ACTST-R01 — STALE: the map aged\n\nSugere: `anchors map build`, and confront again.\n\n"+
			"### ACTST-R02 — FINE: nothing to do\n"), 0o644)
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"### FLOWX-P01 — a step\n\nFits: `ACTST`\n\nResults:\n"+
			"- `ACTST-R01` STALE → `FLOWX-P01`\n- `ACTST-R02` FINE → `FLOWX-P02`\n\n"+
			"### FLOWX-P02 — done\n\n> @terminal\n"), 0o644)

	g, _ := Build(root)
	r1, _ := StateByCode(g, "ACTST-R01")
	if r1.Suggests == "" {
		t.Error("the suggested reaction was not read")
	}
	// A result with no suggestion carries none — the field is not filled by inheritance
	// from the result above it.
	if r2, _ := StateByCode(g, "ACTST-R02"); r2.Suggests != "" {
		t.Errorf("a result with no suggestion must carry none, got %q", r2.Suggests)
	}
}

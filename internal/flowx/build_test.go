package flowx

import (
	"os"
	"path/filepath"
	"strings"
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
	t.Run("FLBLF-B02: Every coded heading is a state of its file", func(t *testing.T) {})
	t.Run("FLBLF-B03: Exits belong to the state above them", func(t *testing.T) {})
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

// Terminal is DECLARED, never deduced from "has no exit": a state with no exit may be the
// end of the work or an oversight, and only the author knows which.
func TestBuild_terminalIsDeclared(t *testing.T) {
	t.Run("FLBLF-B05: Terminal is declared, not inferred", func(t *testing.T) {})
	g, _ := Build(projectWithFlow(t, destravar))
	s, ok := StateByCode(g, "DSTRV-N05")
	if !ok || !s.Terminal {
		t.Errorf("N05 declares @terminal: %+v", s)
	}
	if s, _ := StateByCode(g, "DSTRV-N03"); s.Terminal {
		t.Error("N03 has an exit and does not declare @terminal — it is not terminal")
	}
}

// A project with no `flows/` has not declared any flow — that is not an error.
func TestBuild_projectWithoutFlowsIsNotAnError(t *testing.T) {
	t.Run("FLBLF-B01: Actions then flows are read in a stable order, and no flows give no graph", func(t *testing.T) {})
	g, err := Build(t.TempDir())
	if err != nil || g != nil {
		t.Errorf("expected nil graph and no error, got %+v / %v", g, err)
	}
}

// `Fits:` is what turns the flow into assembly rather than redrawing: the step names the
// piece, and the piece declares its own results. Without it every flow would repeat the
// description of `map build`, and the copies would diverge at the first change.
func TestBuild_readsThePieceEachStepFits(t *testing.T) {
	t.Run("FLBLF-B06: The piece a step fits is read in any language", func(t *testing.T) {})
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
	t.Run("FLBLF-B07: A suggested reaction is recorded as a suggestion", func(t *testing.T) {})
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

// The flow keyword comes from the TRANSLATION CATALOG, never hardcoded.
//
// The first version read `(?:Encaixa|Fits)` — two languages nailed into the pattern, which
// is exactly what the catalog exists to avoid: a project writing in Spanish had no way to
// write in its own language, and adding one would mean editing the engine.
func TestBuild_theFlowKeywordIsTranslated(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, Dir), 0o755)
	// Written in Spanish — a language this engine has no clause for.
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"### FLOWX-P01 — un paso\n\nEncaja: `ACTST`\n\n### FLOWX-P02 — el final\n\n> @terminal\n"), 0o644)
	g, _ := Build(root)
	s, ok := StateByCode(g, "FLOWX-P01")
	if !ok || s.Fits != "ACTST" {
		t.Errorf("the keyword must be read in any catalogued language, got %+v", s)
	}
}

func TestBuild_routedResultsAndStrayExits(t *testing.T) {
	t.Run("FLBLF-B03: Exits belong to the state above them", func(t *testing.T) {})
	t.Run("FLBLF-B04: A line with two codes routes a result", func(t *testing.T) {})
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, Dir), 0o755)
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"- `FLOWX-P09` an exit above every state\n\n"+
			"### FLOWX-P01 — a step\n\nResults:\n- `ACTST-R01` STALE → `FLOWX-P01`\n"), 0o644)
	g, err := Build(root)
	if err != nil || g == nil {
		t.Fatalf("build: %v", err)
	}
	if len(g.Transitions) != 1 {
		t.Fatalf("want only the routed result as a transition (the stray exit has no owner), got %+v", g.Transitions)
	}
	tr := g.Transitions[0]
	if tr.From != "FLOWX-P01" || tr.On != "ACTST-R01" || tr.To != "FLOWX-P01" {
		t.Errorf("routed result = %+v, want from FLOWX-P01 on ACTST-R01 to FLOWX-P01", tr)
	}
}

func TestBuild_readsActionsThenFlowsSorted(t *testing.T) {
	t.Run("FLBLF-B01: Actions then flows are read in a stable order, and no flows give no graph", func(t *testing.T) {})
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ActionsDir), 0o755)
	for name, code := range map[string]string{"b": "ACTBB", "a": "ACTAA"} {
		os.WriteFile(filepath.Join(root, ActionsDir, name+ActionSuffix), []byte("### "+code+"-R01 — a result\n"), 0o644)
	}
	for name, code := range map[string]string{"z": "FLOWZ", "m": "FLOWM"} {
		os.WriteFile(filepath.Join(root, Dir, name+FlowSuffix), []byte("### "+code+"-P01 — a step\n"), 0o644)
	}
	g, err := Build(root)
	if err != nil || g == nil {
		t.Fatalf("build: %v", err)
	}
	var files []string
	for _, s := range g.States {
		files = append(files, s.Flow)
	}
	want := []string{"flows/actions/a.action.md", "flows/actions/b.action.md", "flows/m.flow.md", "flows/z.flow.md"}
	if strings.Join(files, ",") != strings.Join(want, ",") {
		t.Errorf("read order = %v, want %v", files, want)
	}
}

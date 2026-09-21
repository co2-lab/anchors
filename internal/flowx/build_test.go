package flowx

import (
	"os"
	"path/filepath"
	"testing"
)

// The flow that originated the feature, in miniature: the galaxy unblocking rotation.
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
	if err := os.WriteFile(filepath.Join(root, Dir, "unblock"+SufixoFluxo), []byte(content), 0o644); err != nil {
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

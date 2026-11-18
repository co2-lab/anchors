// @anchors
//   code: MDTSM
//   ref: FLMDF

package flowx

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// This is the inversion the whole feature exists for: from a state, the answer is the set
// of valid exits — not the catalogue of every approach with its "when this one is right".
func TestNext_answersTheSetAndNotTheCatalogue(t *testing.T) {
	t.Run("FLMDF-B01: From a state, only its declared exits", func(t *testing.T) {})
	g, _ := Build(projectWithFlow(t, destravar))
	if n := len(Next(g, "DSTRV-N01")); n != 1 {
		t.Errorf("from OPEN there is one way out, got %d", n)
	}
	if n := len(Next(g, "DSTRV-N05")); n != 0 {
		t.Errorf("a terminal state has no exit, got %d", n)
	}
}

// The finding a hand-written flow produces easily: the state is created with ITS exits,
// and the exit that ARRIVES at it is forgotten. It stays written, documented, unreachable.
func TestUnreachable_findsTheStateNobodyArrivesAt(t *testing.T) {
	t.Run("FLMDF-B07: A step nobody arrives at is unreachable", func(t *testing.T) {})
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
	t.Run("FLMDF-B09: An exit to a state that does not exist is dangling", func(t *testing.T) {})
	broken := destravar + "\n### DSTRV-N08 — a state with a broken exit\n\nExits:\n- `DSTRV-N77` when whatever\n"
	g, _ := Build(projectWithFlow(t, broken))
	got := Dangling(g)
	if len(got) != 1 || got[0].To != "DSTRV-N77" {
		t.Errorf("expected exactly the exit to N77, got %+v", got)
	}
}

// The PIECE shape is what makes this finding possible, and the linear drawing hid it: an
// action declares everything it can answer, and a flow that ignores one of those answers
// leaves a hole — whoever gets that result has nowhere to go, and improvises.
func TestUnhandled_findsTheResultNoFlowRoutes(t *testing.T) {
	t.Run("FLMDF-B08: A result no flow routes is unhandled", func(t *testing.T) {})
	t.Run("FLMDF-B07: A step nobody arrives at is unreachable", func(t *testing.T) {})
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ActionsDir), 0o755)
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"### FLOWX-T01 — one step\n\nFits: `ACTST`\n\nResults:\n- `ACTST-O01` FINE → `FLOWX-T02`\n\n"+
			"### FLOWX-T02 — the end\n\n> @terminal\n"), 0o644)
	os.WriteFile(filepath.Join(root, ActionsDir, "a"+ActionSuffix), []byte(
		"### ACTST-O01 — FINE: it worked\n\n### ACTST-O02 — BROKEN: nobody routes this one\n"), 0o644)

	g, err := Build(root)
	if err != nil || g == nil {
		t.Fatalf("build: %v", err)
	}
	got := Unhandled(g)
	if len(got) != 1 || got[0].Code != "ACTST-O02" {
		t.Errorf("expected exactly ACTST-O02 unhandled, got %+v", got)
	}
	// A result is DECLARED by its action, never ARRIVED at: asking it "who reaches you?"
	// would accuse every result a flow happens not to route.
	for _, s := range Unreachable(g) {
		if IsResult(s.Code) {
			t.Errorf("a result must not be charged as unreachable: %s", s.Code)
		}
	}
}

// pieceProject writes an action file (results R01, R02) and a flow that fits it.
func pieceProject(t *testing.T) *mapx.FlowGraph {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ActionsDir), 0o755)
	os.WriteFile(filepath.Join(root, ActionsDir, "check-it"+ActionSuffix), []byte(
		"### ACTST-O01 — FINE: it worked\n\n### ACTST-O02 — BROKEN: it failed\n"), 0o644)
	os.WriteFile(filepath.Join(root, Dir, "f"+FlowSuffix), []byte(
		"### FLOWX-T03 — the first step\n\nFits: `ACTST`\n\nResults:\n- `ACTST-O01` FINE → `FLOWX-T01`\n\n"+
			"### FLOWX-T01 — a later step\n\nExits:\n- `FLOWX-T03` again\n"), 0o644)
	g, err := Build(root)
	if err != nil || g == nil {
		t.Fatalf("build: %v", err)
	}
	return g
}

func TestStateByCodeAndIsResult(t *testing.T) {
	t.Run("FLMDF-B02: A state is found by its code", func(t *testing.T) {})
	t.Run("FLMDF-B05: A result is a code whose letter is O", func(t *testing.T) {})
	g := pieceProject(t)
	if s, ok := StateByCode(g, "FLOWX-T01"); !ok || s.Title != "a later step" {
		t.Errorf("StateByCode(FLOWX-T01) = %+v, %v", s, ok)
	}
	if _, ok := StateByCode(g, "FLOWX-T99"); ok {
		t.Error("a code with no state must not be found")
	}
	for code, want := range map[string]bool{"ACTST-O01": true, "FLOWX-T01": false, "DSTRV-N02": false, "FLOWX-XO01": false, "ACTST": false} {
		if got := IsResult(code); got != want {
			t.Errorf("IsResult(%q) = %v, want %v", code, got, want)
		}
	}
}

func TestStatesOfAndFlows_keepTheFileOrder(t *testing.T) {
	t.Run("FLMDF-B03: States keep the file's order and flows are listed once", func(t *testing.T) {})
	g := pieceProject(t)
	var codes []string
	for _, s := range StatesOf(g, "flows/f.flow.md") {
		codes = append(codes, s.Code)
	}
	if want := []string{"FLOWX-T03", "FLOWX-T01"}; !reflect.DeepEqual(codes, want) {
		t.Errorf("StatesOf = %v, want the file order %v, never alphabetical", codes, want)
	}
	if got, want := Flows(g), []string{"flows/actions/check-it.action.md", "flows/f.flow.md"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Flows = %v, want %v", got, want)
	}
}

func TestEntry_isTheFirstDeclaredStep(t *testing.T) {
	t.Run("FLMDF-B04: The entry is the first declared step, never a result", func(t *testing.T) {})
	g := pieceProject(t)
	if s, ok := Entry(g, "flows/f.flow.md"); !ok || s.Code != "FLOWX-T03" {
		t.Errorf("Entry(flow) = %+v, %v; want FLOWX-T03", s, ok)
	}
	if s, ok := Entry(g, "flows/actions/check-it.action.md"); ok {
		t.Errorf("an action file declares only results and has no entry step, got %+v", s)
	}
}

func TestActionTitle_isTheActionFileName(t *testing.T) {
	t.Run("FLMDF-B06: An action's title is its file name", func(t *testing.T) {})
	g := pieceProject(t)
	if got, ok := ActionTitle(g, "ACTST"); !ok || got != "check-it" {
		t.Errorf("ActionTitle(ACTST) = %q, %v; want check-it", got, ok)
	}
	if _, ok := ActionTitle(g, "NOPEX"); ok {
		t.Error("an action nobody wrote has no title")
	}
}

func TestNilGraphAnswersNothing(t *testing.T) {
	t.Run("FLMDF-X01: A project without flows answers nothing to every question", func(t *testing.T) {})
	if Next(nil, "X") != nil || StatesOf(nil, "f") != nil || Flows(nil) != nil ||
		Unreachable(nil) != nil || Unhandled(nil) != nil || Dangling(nil) != nil {
		t.Error("a nil graph must answer nothing")
	}
	if _, ok := StateByCode(nil, "X"); ok {
		t.Error("a nil graph has no state")
	}
	if _, ok := Entry(nil, "f"); ok {
		t.Error("a nil graph has no entry")
	}
	if _, ok := ActionTitle(nil, "X"); ok {
		t.Error("a nil graph has no action")
	}
}

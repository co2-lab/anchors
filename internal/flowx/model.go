// Package flowx builds and traverses the graph of WORK FLOWS — the states a task passes
// through and the valid transitions between them.
//
// The PROBLEM this solves, measured in a real project: a 232-line guide cataloguing eight
// approaches for getting unstuck, with the two-attempts rule written on line 11 — and a
// dossier that opens by warning "five passes, each one re-aimed the previous one's
// target... two already sent a round to the wrong site".
//
// The cause is not a missing document: it is the NATURE of the document. A catalogue of
// "what to do in case X" requires whoever works to REMEMBER to consult it, CHOOSE right
// among eight, and NOT SKIP a step — three chances to err per round, forever. A flow
// inverts the burden: instead of "choose among eight", it says "from here the exits are
// these three". The rule stops depending on memory and becomes the only way out.
//
// It is the same inversion the gates already perform on the artifact ("remember to write
// the test" becomes `triad-complete`), applied to the PROCESS.
//
// The GRAPH lives in the same file as the map (`anchors.graph.yaml`), under a separate key
// — see `mapx.FlowGraph`, where the reasons for both decisions are written. Here lives the
// LOGIC: scan the `flows/*.flow.md`, assemble, and answer "what comes after here".
package flowx

import (
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
)

// SufixoAcao marks an ACTION file — the puzzle piece.
//
// An action declares WHAT IT DOES and WHICH RESULTS it offers, and deliberately does not
// know who comes next: that is what makes it reusable. `map build` is the clearest case —
// it is a step of the worker cycle AND of adoption, and as a piece it is written once.
//
// A flow does not redraw the work: it FITS actions together and says what to do with each
// result. The regra that today is prose ("NEVER close with a blocking gate red") becomes
// the absence of a fitting — the BARRED result has no link to `done`.
const ActionSuffix = ".action.md"

// DirAcoes is where the reusable pieces live.
const ActionsDir = "flows/actions"

// Dir is where the project's flows live.
const Dir = "flows"

// SufixoFluxo is the extension that marks a flow file.
const FlowSuffix = ".flow.md"

// Next returns the VALID exits of a state — the heart of the feature.
//
// It is what separates "having a flow" from "having one more document": instead of listing
// the eight approaches and asking someone to choose, it answers "from here, these three".
func Next(g *mapx.FlowGraph, code string) []mapx.FlowTransition {
	if g == nil {
		return nil
	}
	var out []mapx.FlowTransition
	for _, t := range g.Transitions {
		if t.From == code {
			out = append(out, t)
		}
	}
	return out
}

// StateByCode finds a state by its code.
func StateByCode(g *mapx.FlowGraph, code string) (mapx.FlowState, bool) {
	if g == nil {
		return mapx.FlowState{}, false
	}
	for _, s := range g.States {
		if s.Code == code {
			return s, true
		}
	}
	return mapx.FlowState{}, false
}

// StatesOf returns the states of ONE flow, in the order they were declared.
//
// The order is the file's, never alphabetical: whoever writes the flow writes it in the
// sequence in which it happens, and reordering would scramble the thread of reading.
func StatesOf(g *mapx.FlowGraph, flow string) []mapx.FlowState {
	if g == nil {
		return nil
	}
	var out []mapx.FlowState
	for _, s := range g.States {
		if s.Flow == flow {
			out = append(out, s)
		}
	}
	return out
}

// Flows returns the names of the flows present in the graph, without repetition and in
// order of appearance.
func Flows(g *mapx.FlowGraph) []string {
	if g == nil {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, s := range g.States {
		if !seen[s.Flow] {
			seen[s.Flow] = true
			out = append(out, s.Flow)
		}
	}
	return out
}

// ActionTitle gives the readable name of an ACTION from its code (`ACHCK`).
//
// The action's own code is not a node: what lives in the graph are its RESULTS
// (`ACHCK-R01`…). So the title is taken from the file that declares them — which is also
// what keeps a step from pointing at a piece nobody wrote.
func ActionTitle(g *mapx.FlowGraph, action string) (string, bool) {
	if g == nil {
		return "", false
	}
	for _, s := range g.States {
		if pre, _, ok := strings.Cut(s.Code, "-"); ok && pre == action {
			// The action's file is `flows/actions/<name>.action.md`; the name is what the
			// reader recognises, not the code.
			base := s.Flow
			if i := strings.LastIndex(base, "/"); i >= 0 {
				base = base[i+1:]
			}
			return strings.TrimSuffix(base, ActionSuffix), true
		}
	}
	return "", false
}

// Entry gives the flow's ENTRY step — the first one declared in the file.
//
// The entry is positional and not deduced from "nobody points at it": a cycle can make
// several steps have no incoming edge (or none at all), and the reader knows where the
// work starts by where the author wrote it first.
func Entry(g *mapx.FlowGraph, flow string) (mapx.FlowState, bool) {
	for _, s := range StatesOf(g, flow) {
		if !IsResult(s.Code) {
			return s, true
		}
	}
	return mapx.FlowState{}, false
}

// IsResult says whether a code names a RESULT of an action (`-R`) rather than a step.
//
// The letter is the role, and it changes which question makes sense: a step is REACHED by
// a transition, while a result is DECLARED by the action that offers it. Asking a result
// "who arrives here?" would accuse every result no flow happens to route.
func IsResult(code string) bool {
	_, rest, ok := strings.Cut(code, "-")
	return ok && strings.HasPrefix(rest, "R")
}

// Unreachable returns the STEPS no transition reaches, minus the first of each flow —
// which is the entry, and by definition has nobody pointing at it.
//
// Results are out of scope here: they are declared by their action, not arrived at. The
// question that does make sense for them is the opposite one — see `Unhandled`.
//
// It is the finding a hand-written flow produces easily: someone creates the step, creates
// ITS exits, and forgets the exit that ARRIVES at it. The step stays written, documented,
// and unreachable — and nobody notices, because it is right there.
func Unreachable(g *mapx.FlowGraph) []mapx.FlowState {
	if g == nil {
		return nil
	}
	reached := map[string]bool{}
	for _, t := range g.Transitions {
		reached[t.To] = true
	}
	entry := map[string]bool{}
	for _, f := range Flows(g) {
		if st := StatesOf(g, f); len(st) > 0 {
			entry[st[0].Code] = true
		}
	}
	var out []mapx.FlowState
	for _, s := range g.States {
		if IsResult(s.Code) {
			continue
		}
		if !reached[s.Code] && !entry[s.Code] {
			out = append(out, s)
		}
	}
	return out
}

// Unhandled returns the RESULTS that no flow routes anywhere.
//
// It is the finding the puzzle shape makes possible, and the linear drawing hid: an action
// declares everything it can answer, and a flow that ignores one of those answers leaves a
// hole — whoever gets that result has nowhere to go, and improvises.
//
// Measured while writing the first flow: `anchors done` with no id, `map build` losing a
// validation stamp, and `work` over a target outside the Structure were all declared and
// routed by nobody.
//
// It is a FINDING and not a failure: a project may legitimately not cover every branch of
// every action. What it must not do is not know.
func Unhandled(g *mapx.FlowGraph) []mapx.FlowState {
	if g == nil {
		return nil
	}
	handled := map[string]bool{}
	for _, t := range g.Transitions {
		if t.On != "" {
			handled[t.On] = true
		}
	}
	var out []mapx.FlowState
	for _, s := range g.States {
		if IsResult(s.Code) && !handled[s.Code] {
			out = append(out, s)
		}
	}
	return out
}

// Dangling returns the transitions pointing at a state that DOES NOT EXIST.
//
// Typically a state renamed without updating the exits that pointed at it — the same
// defect `ref-resolves` catches in artifacts, here in the process. A broken exit is worse
// than an absent one: it LOOKS like a path, and whoever follows it arrives nowhere.
func Dangling(g *mapx.FlowGraph) []mapx.FlowTransition {
	if g == nil {
		return nil
	}
	existe := map[string]bool{}
	for _, s := range g.States {
		existe[s.Code] = true
	}
	var out []mapx.FlowTransition
	for _, t := range g.Transitions {
		if !existe[t.To] {
			out = append(out, t)
		}
	}
	return out
}

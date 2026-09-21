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

import "github.com/co2-lab/anchors/internal/mapx"

// Dir is where the project's flows live.
const Dir = "flows"

// SufixoFluxo is the extension that marks a flow file.
const SufixoFluxo = ".flow.md"

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

// Unreachable returns the states NO transition reaches, minus the first of each flow —
// which is the entry, and by definition has nobody pointing at it.
//
// It is the finding a hand-written flow produces easily: someone creates the state,
// creates ITS exits, and forgets to create the exit that ARRIVES at it. The state stays
// written, documented, and unreachable — and nobody notices, because it is right there.
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
		if !reached[s.Code] && !entry[s.Code] {
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

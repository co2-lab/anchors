package flowx

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
)

// stateRE matches a STATE declaration: `### DSTRV-N03 — MEASURED: the term was isolated`.
//
// The HEADING form only, not the three forms a spec rule accepts. A state opens a section:
// its exits come below it, and a table row has no "below".
var stateRE = regexp.MustCompile(`(?m)^#{1,6}\s+([A-Z0-9]{3,6}-[A-Z]{1,2}[0-9]{2})\s*(?:—|-)?\s*(.*)$`)

// transitionRE matches an EXIT: a bullet carrying the destination code in backticks.
//
// The condition is the rest of the line, in PROSE and deliberately so: whoever judges
// whether it holds is whoever works, by reading. Demanding a computable condition would
// require Anchors to understand each project's domain — and a flow that only accepts the
// computable does not cover the real case ("the measurement points at a concrete target").
// What Anchors guarantees is the SET: from here, these exits, and no others.
var transitionRE = regexp.MustCompile("(?m)^\\s*[-*]\\s+`([A-Z0-9]{3,6}-[A-Z]{1,2}[0-9]{2})`\\s*(.*)$")

// terminalRE marks the state there is no leaving.
//
// Declared, never deduced from "has no exit": a state with no exit may be the end of the
// work OR an oversight, and the gate has to tell the two apart. Whoever writes asserts.
var terminalRE = regexp.MustCompile(`(?im)^\s*>?\s*@terminal\b`)

// Build scans the project's `flows/*.flow.md` and assembles the flow graph.
//
// Separate from `mapx.Build` because it scans other files and produces other content — but
// it writes into the SAME graph, under its own key. See `mapx.FlowGraph`.
func Build(root string) (*mapx.FlowGraph, error) {
	dir := filepath.Join(root, Dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		// A project with no `flows/` is not an error: it is a project that has not declared
		// any flow yet.
		return nil, nil
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), SufixoFluxo) {
			names = append(names, e.Name())
		}
	}
	// STABLE order across runs: without it, two builds of the same repository would produce
	// different files and the diff would turn into noise.
	sort.Strings(names)

	g := &mapx.FlowGraph{}
	for _, nome := range names {
		rel := filepath.ToSlash(filepath.Join(Dir, nome))
		b, err := os.ReadFile(filepath.Join(dir, nome))
		if err != nil {
			continue
		}
		states, transitions := parse(string(b), rel)
		g.States = append(g.States, states...)
		g.Transitions = append(g.Transitions, transitions...)
	}
	if len(g.States) == 0 {
		return nil, nil
	}
	return g, nil
}

// parse reads ONE flow file.
//
// The exit→state association is POSITIONAL: an exit belongs to the last state declared
// above it. That is what allows writing the flow the way it reads — the state, and right
// below it where it leads — instead of repeating the origin code on every line.
func parse(content, flowPath string) ([]mapx.FlowState, []mapx.FlowTransition) {
	var states []mapx.FlowState
	var transitions []mapx.FlowTransition
	current := ""

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if m := stateRE.FindStringSubmatch(line); m != nil {
			current = m[1]
			states = append(states, mapx.FlowState{
				Code:     m[1],
				Title:    strings.TrimSpace(m[2]),
				Flow:     flowPath,
				Terminal: terminalDeclaredAfter(lines, i),
			})
			continue
		}
		if current == "" {
			continue // an exit before any state has no owner
		}
		if m := transitionRE.FindStringSubmatch(line); m != nil {
			transitions = append(transitions, mapx.FlowTransition{
				From: current,
				To:   m[1],
				When: strings.TrimSpace(m[2]),
				Flow: flowPath,
			})
		}
	}
	return states, transitions
}

// terminalDeclaredAfter looks for `@terminal` in the state's body — between its heading
// and the next one.
func terminalDeclaredAfter(lines []string, start int) bool {
	for i := start + 1; i < len(lines); i++ {
		if stateRE.MatchString(lines[i]) {
			return false
		}
		if terminalRE.MatchString(lines[i]) {
			return true
		}
	}
	return false
}

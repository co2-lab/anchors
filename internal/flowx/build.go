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

// fitsRE matches what a flow step FITS: `Encaixa: ` + "`ACHCK`" + `.
//
// It is what turns the flow into assembly rather than redrawing: the step says which piece
// it uses, and the piece declares its own results. Without it, every flow would repeat the
// description of `map build` — and they would diverge at the first change.
var fitsRE = regexp.MustCompile("(?im)^\\s*(?:Encaixa|Fits)\\s*:\\s*`?([A-Z0-9]{3,6})`?")

// resultLinkRE matches a RESULT being routed to a next step:
// `- ` + "`ACHCK-R02`" + ` BARRADO → ` + "`WORKR-P03`" + `.
//
// Two codes on one line: the result that arrived, and where it goes. The arrow may be
// `→`, `->` or nothing — what identifies the destination is being the SECOND code.
var resultLinkRE = regexp.MustCompile("(?m)`?([A-Z0-9]{3,6}-[A-Z]{1,2}[0-9]{2})`?[^`\\n]*?`([A-Z0-9]{3,6}-[A-Z]{1,2}[0-9]{2})`")

// suggestsRE matches the REACTION a result suggests: `Sugere: <o que fazer>`.
//
// It is the third category, and the most common one. A result does not always generate
// work by itself (the automatic reaction) nor end the matter — it often SUGGESTS what to
// do next, and the choice belongs to whoever works.
//
// Measured in the message catalog: 18 occurrences of "run `anchors <command>`" inside gate
// verdicts — `doctor --fix` (5), `map build` (3), `check --fix` (2), `ingest` (3). Each of
// those is a flow edge hidden in prose, where it depends on somebody reading and
// remembering.
//
// It is recorded as a SUGGESTION and not as a transition on purpose: presenting it as an
// ordinary exit would make the flow lie about who decides. The verdict of
// `spec-feature-match` says it plainly — "write the scenario, OR waive with
// `@no-scenario: <reason>`". Two legitimate reactions, and choosing between them requires
// knowing whether the requirement is real.
var suggestsRE = regexp.MustCompile(`(?im)^\s*(?:Sugere|Suggests)\s*:\s*(\S.*)$`)

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
		if !e.IsDir() && strings.HasSuffix(e.Name(), FlowSuffix) {
			names = append(names, e.Name())
		}
	}
	// STABLE order across runs: without it, two builds of the same repository would produce
	// different files and the diff would turn into noise.
	sort.Strings(names)

	// The ACTIONS live in a subfolder and enter the SAME graph: they are nodes like the
	// steps, and the difference is the role, not the structure. A result (`ACHCK-R02`) is
	// the target of a transition just like a step — what changes is who declares it.
	actionEntries, _ := os.ReadDir(filepath.Join(root, ActionsDir))
	var actionNames []string
	for _, e := range actionEntries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ActionSuffix) {
			actionNames = append(actionNames, e.Name())
		}
	}
	sort.Strings(actionNames)

	g := &mapx.FlowGraph{}
	for _, nome := range actionNames {
		rel := filepath.ToSlash(filepath.Join(ActionsDir, nome))
		b, err := os.ReadFile(filepath.Join(root, ActionsDir, nome))
		if err != nil {
			continue
		}
		states, transitions := parse(string(b), rel)
		g.States = append(g.States, states...)
		g.Transitions = append(g.Transitions, transitions...)
	}
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

// parse reads ONE flow or action file.
//
// The exit→state association is POSITIONAL: what comes below a state belongs to it. That
// is what allows writing the flow the way it reads — the step, and right below it the
// piece it fits and where each result goes.
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
				Fits:     fitsDeclaredAfter(lines, i),
				Suggests: suggestsDeclaredAfter(lines, i),
			})
			continue
		}
		if current == "" {
			continue // an exit before any state has no owner
		}
		// A ROUTED RESULT (two codes on the line) wins over a plain exit: the first code
		// is the result that arrived, the second is where it goes. Reading it as a plain
		// exit would make the flow point at the RESULT instead of at the next step — an
		// edge to a node that lives in another file and is not a step.
		if m := resultLinkRE.FindStringSubmatch(line); m != nil {
			transitions = append(transitions, mapx.FlowTransition{
				From: current, To: m[2], On: m[1],
				When: conditionOf(line), Flow: flowPath,
			})
			continue
		}
		if m := transitionRE.FindStringSubmatch(line); m != nil {
			transitions = append(transitions, mapx.FlowTransition{
				From: current, To: m[1],
				When: strings.TrimSpace(m[2]), Flow: flowPath,
			})
		}
	}
	return states, transitions
}

// conditionOf keeps the prose of the line minus the two codes — what whoever works reads
// to recognise the result.
func conditionOf(line string) string {
	out := resultLinkRE.ReplaceAllString(line, "")
	out = strings.TrimSpace(strings.Trim(strings.TrimSpace(out), "-*→>` "))
	return out
}

// fitsDeclaredAfter reads which PIECE a step fits — between its heading and the next.
func fitsDeclaredAfter(lines []string, start int) string {
	for i := start + 1; i < len(lines); i++ {
		if stateRE.MatchString(lines[i]) {
			return ""
		}
		if m := fitsRE.FindStringSubmatch(lines[i]); m != nil {
			return m[1]
		}
	}
	return ""
}

// suggestsDeclaredAfter reads the reaction a result suggests — between its heading and the
// next one.
func suggestsDeclaredAfter(lines []string, start int) string {
	for i := start + 1; i < len(lines); i++ {
		if stateRE.MatchString(lines[i]) {
			return ""
		}
		if m := suggestsRE.FindStringSubmatch(lines[i]); m != nil {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
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

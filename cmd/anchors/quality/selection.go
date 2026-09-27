package quality

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// fileState is where a file stands for a suite, read from the signal the map keeps.
//
// Two questions, each with two answers, make the four boxes a run chooses from:
//
//	                 fresh            stale
//	passing          freshPassing     stalePassing
//	below minimum    freshBelow       staleBelow
//
// "Fresh" means measured at the current version — of the file and, for a test, of what
// it exercises. "Passing" means a test file with no failing case, or a code file whose
// mutation score is at the floor or above. A file never measured is `unmeasured`.
type fileState int

const (
	unmeasured fileState = iota
	freshPassing
	stalePassing
	freshBelow
	staleBelow
)

// runSelection is which states a run takes. By default only what is stale AND below the
// minimum runs, with what was never measured: everything else is known, and a light run
// skips it. Each flag opens one side of the square.
type runSelection struct {
	IncludeFresh   bool // also fresh files below the minimum (and, with IncludePassing, fresh passing ones)
	IncludePassing bool // also stale files at the minimum or above (and, with IncludeFresh, fresh passing ones)
	SkipUnmeasured bool // leave out the files never measured
}

// wants says whether the run takes a file in this state.
func (s runSelection) wants(st fileState) bool {
	switch st {
	case unmeasured:
		return !s.SkipUnmeasured
	case staleBelow:
		return true
	case freshBelow:
		return s.IncludeFresh
	case stalePassing:
		return s.IncludePassing
	default: // freshPassing
		return s.IncludeFresh && s.IncludePassing
	}
}

// testState reads a test file's state. Stale is the test file changed since its result,
// or something it exercises did (the evidence closure). A file whose cases were all
// skipped is not passing: nothing proved anything.
func testState(g *mapx.Graph, n mapx.Node, layer string) fileState {
	s := n.Signal
	if s == nil {
		return unmeasured
	}
	passed, failed, skipped := s.Passed, s.Failed, s.Skipped
	if le, ok := s.ByLayer[layer]; ok {
		// The suite's own layer, when the file was run by several: a unit failure is not
		// an integration one.
		passed, failed, skipped = le.Passed, le.Failed, le.Skipped
	}
	if passed+failed+skipped == 0 {
		return unmeasured
	}
	stale := n.SignalStale() || g.EvidenceStaleFor(n.ID) != nil
	passing := failed == 0 && passed > 0
	return square(stale, passing)
}

// ownedBy says whether a test file is one this suite runs. A suite that timed the file
// owns it; before any time was recorded, the suite whose layer ran it does. A file no
// suite ever ran is offered to every suite — in a monorepo, once, until a run times it.
func ownedBy(n mapx.Node, key, layer string) bool {
	s := n.Signal
	if s == nil {
		return true
	}
	if len(s.SecondsBySuite) > 0 {
		_, ok := s.SecondsBySuite[key]
		return ok
	}
	if len(s.ByLayer) > 0 {
		_, ok := s.ByLayer[layer]
		return ok
	}
	return true
}

// mutationState reads a code file's mutation state. A measurement taken under load (more
// timed-out mutants than the ceiling allows) is not a measurement to keep, so it counts
// as stale, whatever its revision.
func mutationState(n mapx.Node, ceiling float64) fileState {
	s := n.Signal
	if s == nil || s.MutantsKilled+s.MutantsSurvived+s.MutantsNoCoverage+s.MutantsIgnored == 0 {
		return unmeasured
	}
	ran := s.MutantsKilled + s.MutantsSurvived
	underLoad := ran > 0 && float64(s.MutantsTimedOut)/float64(ran) > ceiling
	stale := n.MutationStale() || underLoad
	floor := s.MutationLow
	if floor <= 0 {
		floor = 70
	}
	passing := ran == 0 || s.MutationScore >= floor
	return square(stale, passing)
}

func square(stale, passing bool) fileState {
	switch {
	case !stale && passing:
		return freshPassing
	case stale && passing:
		return stalePassing
	case !stale:
		return freshBelow
	default:
		return staleBelow
	}
}

// selectFiles picks, for a suite, the files a run takes and counts the ones it leaves out
// by state. Support files and targets the matching gate declares with nothing to measure
// (`no_signal`) are never candidates.
func selectFiles(g *mapx.Graph, cfg *config.Config, mutation bool, key string, suite config.Suite, sel runSelection) (run []string, left map[fileState]int) {
	layer := suite.Layer
	kind, check := mapx.KindTest, "tests-pass"
	if mutation {
		kind, check = mapx.KindCode, "mutation-score"
	}
	gateOf := config.Gate{}
	ceiling := config.DefaultTimeoutCeiling
	if cfg != nil {
		for _, gt := range cfg.Gates {
			if gt.Check == check {
				gateOf = gt
				ceiling = gt.TimeoutCeilingOrDefault()
				break
			}
		}
	}
	left = map[fileState]int{}
	for _, n := range g.Nodes {
		if n.Kind != kind || n.Support || !suite.Covers(n.ID) {
			continue
		}
		if _, none := gateOf.NoSignalFor(n.ID); none {
			continue
		}
		var st fileState
		if mutation {
			// A mutation signal does not say which suite measured it: every code file is a
			// candidate of each mutation suite.
			st = mutationState(n, ceiling)
		} else {
			if !ownedBy(n, key, layer) {
				continue
			}
			st = testState(g, n, layer)
		}
		if sel.wants(st) {
			run = append(run, n.ID)
		} else {
			left[st]++
		}
	}
	sort.Strings(run)
	return run, left
}

// describeLeftOut says what a selective run did not take, and the flag that takes it.
func describeLeftOut(left map[fileState]int) string {
	var parts []string
	for _, c := range []struct {
		st   fileState
		what string
		flag string
	}{
		{freshPassing, "fresh and passing", "--include-fresh --include-passing"},
		{stalePassing, "stale but passing", "--include-passing"},
		{freshBelow, "fresh and below the minimum", "--include-fresh"},
		{unmeasured, "never measured", "drop --skip-unmeasured"},
	} {
		if n := left[c.st]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s (%s)", n, c.what, c.flag))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "left out: " + strings.Join(parts, "; ")
}

// runSelective is the default run: for each suite, the files the selection takes — stale
// and below the minimum, and never measured, unless the flags open more of the square —
// through the suite's `run_changed:`, each batch ingested as a partial run. With a
// budget, fastest first until it is spent. A suite with no `run_changed:` cannot run a
// subset: it runs whole, and says so. With a budget such a suite is refused, as any
// budget run refuses it, before anything runs.
func runSelective(cs suiteCommand, suites []config.Suite, cfg *config.Config, absRoot, target string, budget time.Duration, picked runSelection) error {
	mutation := cs.secao == "mutation"
	var subset, whole []config.Suite
	for _, s := range suites {
		if strings.TrimSpace(s.RunChanged) == "" {
			if budget > 0 {
				return fmt.Errorf("layer %q: %w", s.Layer, budgetRunnable(cs, s))
			}
			whole = append(whole, s)
		} else {
			subset = append(subset, s)
		}
	}
	var only []map[string]bool
	var files [][]string
	if len(subset) > 0 {
		g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
		if err != nil {
			return fmt.Errorf("load map: %w (run `anchors map build`, or `--all` to run the suites whole)", err)
		}
		for _, s := range subset {
			if budget > 0 {
				if err := budgetRunnable(cs, s); err != nil {
					return fmt.Errorf("layer %q: %w", s.Layer, err)
				}
			}
			report := s.JUnit
			if mutation {
				report = s.Report
			}
			run, left := selectFiles(g, cfg, mutation, mapcmd.SuiteKey(absRoot, absPath(absRoot, report)), s, picked)
			fmt.Printf("[%s%s] selected %d file(s) to run", workspaceLabel(s), s.Layer, len(run))
			if d := describeLeftOut(left); d != "" {
				fmt.Printf("; %s", d)
			}
			fmt.Println()
			set := make(map[string]bool, len(run))
			for _, id := range run {
				set[id] = true
			}
			only = append(only, set)
			files = append(files, run)
		}
		fmt.Println()
	}
	for _, s := range whole {
		fmt.Printf("[%s%s] declares no `run_changed:`: it cannot run a subset, and runs whole\n\n", workspaceLabel(s), s.Layer)
	}
	if budget > 0 {
		return budgetSuites(cs, subset, only, absRoot, target, time.Now().Add(budget))
	}
	for i, s := range subset {
		if len(files[i]) == 0 {
			fmt.Printf("[%s%s] nothing to run: every file is left out by its state (`--all` runs them)\n\n", workspaceLabel(s), s.Layer)
			continue
		}
		if err := runFiles(cs, s, absRoot, target, files[i]); err != nil {
			return err
		}
	}
	if len(whole) > 0 {
		return runSuites(cs, whole, absRoot, target, nil)
	}
	return nil
}

// runFiles runs a suite's `run_changed:` over the files, in as few batches as the
// command-line ceiling allows, and ingests each as a partial run. It stops at the first
// batch that fails, as a whole run stops at the first failing suite.
func runFiles(cs suiteCommand, s config.Suite, absRoot, target string, ids []string) error {
	for _, batch := range argvBatches(s, absFiles(absRoot, ids), suiteArgvLimit()) {
		linha, err := pickCommand(s, batch, target)
		if err != nil {
			return fmt.Errorf("layer %q: %w", s.Layer, err)
		}
		fmt.Printf("━━━ %s [%s%s] %d file(s) ━━━\n%s\n\n", cs.nome, workspaceLabel(s), s.Layer, len(batch), linha)
		start := time.Now()
		errRun := execAtRoot(linha, absRoot)
		if err := ingestIfRecent(absRoot, absPath(absRoot, s.JUnit), absPath(absRoot, s.Lcov), absPath(absRoot, s.Report), s, start, true); err != nil {
			return fmt.Errorf("layer %q: ingest report: %w", s.Layer, err)
		}
		if errRun != nil {
			return fmt.Errorf("layer %q failed: %w", s.Layer, errRun)
		}
	}
	return nil
}

// argvBatches splits the files so each batch's command line stays under the ceiling. A
// single file longer than the ceiling still goes alone: refusing it would leave it never
// run.
func argvBatches(s config.Suite, files []string, ceiling int) [][]string {
	base := len(strings.ReplaceAll(s.RunChanged, "{{files}}", ""))
	var out [][]string
	var cur []string
	size := base
	for _, f := range files {
		if len(cur) > 0 && size+len(f)+1 > ceiling {
			out = append(out, cur)
			cur, size = nil, base
		}
		cur = append(cur, f)
		size += len(f) + 1
	}
	if len(cur) > 0 {
		out = append(out, cur)
	}
	return out
}

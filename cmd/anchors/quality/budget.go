package quality

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// timedFile is a file with the time its last run took in a suite.
type timedFile struct {
	ID      string
	Seconds float64
}

// budgetPlan orders a suite's files for a run with a time budget: the ones this suite
// has timed, fastest first (ties by path), then the ones no suite has timed yet, by path.
// A file only other suites have timed belongs to them and is left out. Support files are
// not run.
//
// There is no priority besides time: the project chooses the budget, and the run fits in
// as much as it can.
func budgetPlan(g *mapx.Graph, kind mapx.Kind, key string, suite config.Suite) (timed []timedFile, untimed []string) {
	for _, n := range g.Nodes {
		if n.Kind != kind || n.Support || !suite.Covers(n.ID) {
			continue
		}
		var bySuite map[string]float64
		if n.Signal != nil {
			bySuite = n.Signal.SecondsBySuite
		}
		switch sec, ok := bySuite[key]; {
		case ok:
			timed = append(timed, timedFile{ID: n.ID, Seconds: sec})
		case len(bySuite) == 0:
			untimed = append(untimed, n.ID)
		}
	}
	sort.Slice(timed, func(i, j int) bool {
		if timed[i].Seconds != timed[j].Seconds {
			return timed[i].Seconds < timed[j].Seconds
		}
		return timed[i].ID < timed[j].ID
	})
	sort.Strings(untimed)
	return timed, untimed
}

// nextBatch takes, from the front of the plan, the timed files whose times add up within
// what remains of the budget. When even the fastest remaining one does not fit, the batch
// is empty: running it would only be cut by the deadline. After the timed files, the
// untimed go one per batch — nobody knows how long they take, and each one's run is what
// times it. `oneAtATime` makes every batch a single file, for runs Anchors times itself.
func nextBatch(timed []timedFile, untimed []string, remaining time.Duration, oneAtATime bool) (batch []string, restTimed []timedFile, restUntimed []string) {
	left := remaining.Seconds()
	i := 0
	for i < len(timed) && timed[i].Seconds <= left {
		left -= timed[i].Seconds
		batch = append(batch, timed[i].ID)
		i++
		if oneAtATime {
			break
		}
	}
	if len(batch) > 0 || i < len(timed) {
		return batch, timed[i:], untimed
	}
	if len(untimed) > 0 && remaining > 0 {
		return []string{untimed[0]}, nil, untimed[1:]
	}
	return nil, nil, untimed
}

// budgetOutcome is what a run with a budget did, for the report at the end.
type budgetOutcome struct {
	Ran, Left int
	Failed    bool
	Cut       bool // the deadline stopped a batch in the middle
}

// runWithBudget runs each selected suite's files in batches, fastest first, until the
// budget is spent, and ingests each batch as a partial run. A batch still running when
// the deadline comes is stopped with its whole process group, and nothing of it is
// ingested: a run cut in the middle measured nothing whole.
func runWithBudget(cs suiteCommand, suites []config.Suite, absRoot, target string, budget time.Duration) error {
	// Every suite is checked before any runs: refusing the third after the first two ran
	// would spend the budget on a run the user then has to redo.
	for _, s := range suites {
		if err := budgetRunnable(cs, s); err != nil {
			return fmt.Errorf("layer %q: %w", s.Layer, err)
		}
	}
	deadline := time.Now().Add(budget)
	only := make([]map[string]bool, len(suites))
	return budgetSuites(cs, suites, only, absRoot, target, deadline)
}

// budgetSuites runs each suite under the shared deadline, each over the files `only`
// allows (nil: every file), and reports what ran and what was left.
func budgetSuites(cs suiteCommand, suites []config.Suite, only []map[string]bool, absRoot, target string, deadline time.Time) error {
	var failed []string
	for i, s := range suites {
		out, err := runSuiteWithBudget(cs, s, absRoot, target, deadline, only[i])
		if err != nil {
			return fmt.Errorf("layer %q: %w", s.Layer, err)
		}
		fmt.Printf("[%s%s] budget: ran %d file(s); %d left for a later run%s\n\n", workspaceLabel(s), s.Layer,
			out.Ran, out.Left, map[bool]string{true: " (the last batch was cut by the deadline)", false: ""}[out.Cut])
		if out.Failed {
			failed = append(failed, s.Layer)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("a batch failed in: %s", strings.Join(failed, ", "))
	}
	return nil
}

func runSuiteWithBudget(cs suiteCommand, s config.Suite, absRoot, target string, deadline time.Time, only map[string]bool) (budgetOutcome, error) {
	var out budgetOutcome
	report, kind, oneAtATime := s.JUnit, mapx.KindTest, false
	if cs.secao == "mutation" {
		// A mutation report carries no time per file, so Anchors times each file's run
		// itself — one file per run.
		report, kind, oneAtATime = s.Report, mapx.KindCode, true
	}
	key := mapcmd.SuiteKey(absRoot, absPath(absRoot, report))
	mapPath := filepath.Join(absRoot, mapx.DefaultPath)
	g, err := mapx.Load(mapPath)
	if err != nil {
		return out, fmt.Errorf("load map: %w (run `anchors map build`)", err)
	}
	timed, untimed := budgetPlan(g, kind, key, s)
	if only != nil {
		timed, untimed = keepOnly(timed, untimed, only)
	}
	for {
		remaining := time.Until(deadline)
		var batch []string
		batch, timed, untimed = nextBatch(timed, untimed, remaining, oneAtATime)
		if len(batch) == 0 {
			break
		}
		linha, err := pickCommand(s, absFiles(absRoot, batch), target)
		if err != nil {
			return out, err
		}
		fmt.Printf("━━━ %s [%s%s] budget batch of %d, %s left ━━━\n%s\n\n", cs.nome, workspaceLabel(s), s.Layer,
			len(batch), remaining.Round(time.Second), linha)
		start := time.Now()
		errRun, cut := execUntil(linha, absRoot, deadline)
		if cut {
			out.Cut = true
			out.Left += len(batch)
			break
		}
		if err := ingestIfRecent(absRoot, absPath(absRoot, s.JUnit), absPath(absRoot, s.Lcov), absPath(absRoot, s.Report), s, start, true); err != nil {
			return out, fmt.Errorf("ingest report: %w", err)
		}
		if oneAtATime && errRun == nil {
			if err := recordRunSeconds(mapPath, batch[0], key, time.Since(start).Seconds()); err != nil {
				return out, err
			}
		}
		out.Ran += len(batch)
		if errRun != nil {
			out.Failed = true
		}
	}
	out.Left += len(timed) + len(untimed)
	return out, nil
}

// recordRunSeconds keeps the time Anchors measured for a file's run in the map.
func recordRunSeconds(mapPath, id, key string, seconds float64) error {
	g, err := mapx.Load(mapPath)
	if err != nil {
		return fmt.Errorf("load map: %w", err)
	}
	g.RecordRunSeconds(id, key, seconds)
	return mapx.Save(g, mapPath)
}

// absFiles turns map IDs into the absolute forward-slash paths the suite's command
// receives, as the incremental mode does.
func absFiles(absRoot string, ids []string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, filepath.ToSlash(filepath.Join(absRoot, filepath.FromSlash(id))))
	}
	return out
}

// stopGrace is how long a cut batch has, after the TERM, to undo what it changed before
// its group is killed. A variable so the tests do not wait ten seconds.
var stopGrace = 10 * time.Second

// execUntil runs the suite's command at the root until the deadline. A command still
// running then is stopped with its whole process group — the runner's workers too, not
// only the shell — and `cut` says so.
//
// The group gets a TERM first, and a KILL only if it is still there after stopGrace. A
// mutation tool that works IN PLACE rewrites the source and restores it when it exits
// normally or on a signal it can trap; a KILL runs no trap, and the cut batch left a
// mutated file in the tree (reported from the reference app: Stryker with `inPlace`).
func execUntil(linha, absRoot string, deadline time.Time) (err error, cut bool) {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", linha) //nolint:gosec // the command is declared by the project, as a gate's `run:`
	cmd.Dir = absRoot
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	ownProcessGroup(cmd)
	cmd.Cancel = func() error { return stopProcessGroup(cmd, stopGrace) }
	cmd.WaitDelay = stopGrace + 5*time.Second
	err = cmd.Run()
	return err, ctx.Err() == context.DeadlineExceeded
}

// budgetRunnable says why a suite cannot run with a budget: it needs `run_changed:` to
// receive the batches, and a report to learn the times from.
func budgetRunnable(cs suiteCommand, s config.Suite) error {
	if strings.TrimSpace(s.RunChanged) == "" {
		return fmt.Errorf("`--budget` runs the files in batches through `run_changed:`, and this suite declares none")
	}
	report := s.JUnit
	if cs.secao == "mutation" {
		report = s.Report
	}
	if strings.TrimSpace(report) == "" {
		return fmt.Errorf("`--budget` needs the suite's report to ingest each batch, and this suite declares none")
	}
	return nil
}

// keepOnly narrows a plan to the files a selection took, keeping its order.
func keepOnly(timed []timedFile, untimed []string, only map[string]bool) ([]timedFile, []string) {
	var t []timedFile
	for _, f := range timed {
		if only[f.ID] {
			t = append(t, f)
		}
	}
	var u []string
	for _, id := range untimed {
		if only[id] {
			u = append(u, id)
		}
	}
	return t, u
}

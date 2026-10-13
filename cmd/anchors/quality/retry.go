// @anchors
//   code: RTTSB
//   ref: STPRS

package quality

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/runs"
	"github.com/co2-lab/anchors/internal/testsig"
)

// --- one retry of a test that ran out of time under load (DESIGN-retry-on-load.md) ---

// loadSampler keeps the highest load seen while a run lasts: the load at the moment of a
// failure is what tells a test too slow for a loaded machine from a test that is wrong.
type loadSampler struct {
	mu       sync.Mutex
	peak     float64
	measured bool
	stop     chan struct{}
	done     chan struct{}
}

// loadSampleEvery is how often the sampler reads the load; a variable so a test is quick.
var loadSampleEvery = 2 * time.Second

// loadAverage reads the load; a variable so a test can play a loaded machine.
var loadAverage = runs.LoadAverage

func sampleLoad() *loadSampler {
	ls := &loadSampler{stop: make(chan struct{}), done: make(chan struct{})}
	ls.read()
	go func() {
		defer close(ls.done)
		t := time.NewTicker(loadSampleEvery)
		defer t.Stop()
		for {
			select {
			case <-ls.stop:
				return
			case <-t.C:
				ls.read()
			}
		}
	}()
	return ls
}

func (ls *loadSampler) read() {
	if v, ok := loadAverage(); ok {
		ls.mu.Lock()
		ls.measured = true
		if v > ls.peak {
			ls.peak = v
		}
		ls.mu.Unlock()
	}
}

// Stop ends the sampling and says the highest load, read once more at the end.
func (ls *loadSampler) Stop() (float64, bool) {
	ls.read()
	close(ls.stop)
	<-ls.done
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return ls.peak, ls.measured
}

// retryOnLoad runs again, once per retry the suite allows, the tests of a failed run that ran
// out of time, when the run's load passed the suite's threshold. It says whether the run now
// passes: every failure was a timeout, and the retry passed them. A failed assertion is never
// repeated, and a test that passes on a retry is recorded as such.
func retryOnLoad(s config.Suite, absRoot, target, junit, lcov string, peak float64, measured bool) (bool, error) {
	if s.RetryOnLoad == nil || junit == "" {
		return false, nil
	}
	factor, err := s.RetryOnLoad.Factor()
	if err != nil {
		return false, fmt.Errorf("layer %q: %w", s.Layer, err)
	}
	pat := testsig.DefaultTimeoutPattern
	if s.TimeoutPattern != "" {
		pat = s.TimeoutPattern
	}
	pattern, err := regexp.Compile(pat)
	if err != nil {
		return false, fmt.Errorf("layer %q: timeout_pattern %q: %v", s.Layer, pat, err)
	}
	cpus := runtime.NumCPU()
	load := fmt.Sprintf("%.1f/%d", peak, cpus)
	for attempt := 1; attempt <= s.RetryOnLoad.Retries(); attempt++ {
		rep, err := testsig.ParseJUnit(junit)
		if err != nil {
			return false, nil
		}
		var timed []testsig.CaseResult
		failed := 0
		for _, c := range rep.Cases {
			if !c.Failed {
				continue
			}
			failed++
			if testsig.TimedOut(c, pattern, absRoot) && c.File != "" {
				timed = append(timed, c)
			}
		}
		if failed == 0 {
			return true, nil
		}
		if len(timed) == 0 {
			return false, nil
		}
		if !measured || peak < factor*float64(cpus) {
			fmt.Printf("  [%s] %d test(s) ran out of time, but the load stayed at %s, under %gx the CPUs: no retry\n\n",
				s.Layer, len(timed), load, factor)
			return false, nil
		}
		run := s.RetryOnLoad.Run
		if strings.TrimSpace(run) == "" {
			run = s.RunChanged
		}
		if strings.TrimSpace(run) == "" {
			fmt.Printf("  [%s] %d test(s) ran out of time under load, but the suite has no `retry_on_load.run` nor `run_changed` to run them alone\n\n", s.Layer, len(timed))
			return false, nil
		}
		seen := map[string]bool{}
		var files []string
		for _, c := range timed {
			abs := filepath.ToSlash(filepath.Join(absRoot, filepath.FromSlash(c.File)))
			if !seen[abs] {
				seen[abs] = true
				files = append(files, abs)
			}
		}
		sort.Strings(files)
		line, err := buildCommand(strings.ReplaceAll(run, "{{files}}", strings.Join(files, " ")), target)
		if err != nil {
			return false, fmt.Errorf("layer %q: %w", s.Layer, err)
		}
		fmt.Printf("━━━ retry %d/%d [%s%s] %d test(s) ran out of time under load %s ━━━\n%s\n\n",
			attempt, s.RetryOnLoad.Retries(), workspaceLabel(s), s.Layer, len(timed), load, line)
		start := time.Now()
		errRetry := execAtRoot(line, absRoot)
		if err := ingestIfRecent(absRoot, junit, lcov, "", s, start, true); err != nil {
			return false, fmt.Errorf("layer %q: ingest the retry's report: %w", s.Layer, err)
		}
		again, err := testsig.ParseJUnit(junit)
		if err != nil {
			return false, nil
		}
		passed := map[string]bool{}
		for _, c := range again.Cases {
			if !c.Failed && !c.Skipped {
				passed[c.File+"\x00"+c.Name] = true
			}
		}
		var record []testsig.Retried
		var names []string
		for _, c := range timed {
			if passed[c.File+"\x00"+c.Name] {
				record = append(record, testsig.Retried{Name: c.Name, File: c.File, At: time.Now(), Load: peak, CPUs: cpus})
				names = append(names, c.Name)
			}
		}
		_ = testsig.AddRetried(absRoot, record)
		if len(names) > 0 {
			fmt.Printf("  [%s] %d test(s) passed on retry (load %s, timeout): %s\n\n", s.Layer, len(names), load, strings.Join(names, "; "))
		}
		if errRetry == nil && len(timed) == failed && len(names) == len(timed) {
			return true, nil
		}
		if len(timed) < failed {
			return false, nil // a failure that was no timeout stays a failure
		}
	}
	return false, nil
}

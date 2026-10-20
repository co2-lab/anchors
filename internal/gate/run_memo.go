// @anchors
//   code: RNMEM
//   ref: GTENG

package gate

import (
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
)

// --- what a run answers once ---
//
// A run confronts every target with every gate, and two questions repeat across them with
// the same answer: whether the configuration declares a gate's premises — read through the
// whole configuration as YAML, 28% of a `check --all` at MIF, asked once per gate and target —
// and the content of a project file, read again by each gate of the target (another 10%).
// Both are answered once for the run, and forgotten when it ends.

var (
	runMu       sync.Mutex
	runCfg      *config.Config
	runActive   bool
	runPremise  map[string][2][]string
	runContents map[string][]byte
)

// startRun opens the run's memory for cfg, and returns what closes it.
func startRun(cfg *config.Config) func() {
	runMu.Lock()
	runCfg, runActive = cfg, true
	runPremise, runContents = map[string][2][]string{}, map[string][]byte{}
	runMu.Unlock()
	return func() {
		runMu.Lock()
		runCfg, runActive, runPremise, runContents = nil, false, nil, nil
		runMu.Unlock()
	}
}

// forgetReads drops the contents the run read: a tool the project runs may have written
// the files since.
func forgetReads() {
	runMu.Lock()
	if runActive {
		runContents = map[string][]byte{}
	}
	runMu.Unlock()
}

// runPremises is cfg.Premises(g), answered once per run for each gate's premises.
func runPremises(g config.Gate, cfg *config.Config) (missing, waived []string) {
	key := strings.Join(g.Presupposes, "\x00")
	runMu.Lock()
	if runActive && runCfg == cfg {
		if p, ok := runPremise[key]; ok {
			runMu.Unlock()
			return p[0], p[1]
		}
	}
	runMu.Unlock()
	missing, waived = cfg.Premises(g)
	runMu.Lock()
	if runActive && runCfg == cfg {
		runPremise[key] = [2][]string{missing, waived}
	}
	runMu.Unlock()
	return missing, waived
}

// readOnce is the run's copy of a file's content, read through read the first time.
func readOnce(key string, read func() ([]byte, error)) ([]byte, error) {
	runMu.Lock()
	if !runActive {
		runMu.Unlock()
		return read()
	}
	if b, ok := runContents[key]; ok {
		runMu.Unlock()
		// A copy: whoever reads may change what it got, and the next reader must not see it.
		return append([]byte(nil), b...), nil
	}
	runMu.Unlock()
	b, err := read()
	if err != nil {
		return b, err
	}
	runMu.Lock()
	if runActive {
		runContents[key] = append([]byte(nil), b...)
	}
	runMu.Unlock()
	return b, nil
}

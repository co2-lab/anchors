// @anchors
//   code: LDQRV
//   ref: LDQRS

package quality

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/load"
)

// admitRun waits for the project's load queue to admit a suite's run over its files — the
// run and the ingestion after it, which weighs too —, and returns what releases it. With no
// `load:`, in the commit hook, or when the queue cannot be read, the run starts at once.
func admitRun(absRoot, task string, s config.Suite, files []string) func() {
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil || !load.Enabled(cfg.Load) {
		return func() {}
	}
	rel := make([]string, 0, len(files))
	for _, f := range files {
		if r, err := filepath.Rel(absRoot, f); err == nil && filepath.IsAbs(f) {
			f = r
		}
		rel = append(rel, filepath.ToSlash(f))
	}
	// A whole run names no file: it is placed by its workspace, which the scopes name.
	if len(rel) == 0 && s.Workspace != "" {
		rel = []string{filepath.ToSlash(filepath.Join(s.Workspace, "_"))}
	}
	level := cfg.Load.LevelFor(task, rel)
	cpu, mem, _ := cfg.Load.LevelCost(level)
	release, err := load.Acquire(load.Request{
		Project: absRoot, Task: task + " " + s.Layer, Level: level,
		Cost: load.Cost{CPU: cpu, Memory: mem}, WaitLimit: cfg.Load.WaitLimitFor(task),
		Notify: func(msg string) { fmt.Println(msg) },
	})
	if err != nil {
		fmt.Printf("load: the queue could not be read (%v): the run starts at once\n", err)
	}
	return release
}

// skipLoadQueue keeps this process and its children out of the load queue: a commit never
// waits behind another worktree's heavy run.
func skipLoadQueue() { _ = os.Setenv("ANCHORS_LOAD", "off") }

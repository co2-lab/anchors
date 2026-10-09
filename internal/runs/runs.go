// @anchors
//   code: RNRCR
//   ref: PRCRN

// Package runs keeps the record of the project's long processes — test suites, mutation runs,
// builds, the commit hook — and reads them from the operating system, so a monitor can say
// when one starts, stalls, dies or ends (DESIGN-process-monitor.md).
package runs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Dir is where the records live, under the project root: machine state, like the rest of
// `.anchors/`.
const Dir = ".anchors/runs"

// Who wrote a record.
const (
	ByAnchors = "anchors" // one of Anchors' own long commands
	ByAgent   = "agent"   // the agent hook, for a command the agent launched
	ByRun     = "run"     // `anchors monitor run -- <command>`
	ByOS      = "os"      // the monitor, for a process it found in the process table
)

// Run is one long process: what it is, where, who launched it, where its output goes, and —
// once it ended — how. A record written by whoever launched the process knows its exit; one
// the monitor wrote from the process table knows only that it ended.
type Run struct {
	ID      string     `json:"id"`
	Kind    string     `json:"kind"`
	Name    string     `json:"name"`
	Command string     `json:"command"`
	Dir     string     `json:"dir,omitempty"`
	PID     int        `json:"pid,omitempty"`
	By      string     `json:"by"`
	Started time.Time  `json:"started"`
	Output  string     `json:"output,omitempty"`
	Ended   *time.Time `json:"ended,omitempty"`
	Exit    *int       `json:"exit,omitempty"`
	Summary string     `json:"summary,omitempty"`
	// State is the monitor's last reading: running, stalled, died, finished or ended.
	State string `json:"state,omitempty"`
}

// Done says whether the run ended, by any account.
func (r Run) Done() bool { return r.Ended != nil }

// NewID is a run's id: when it started and its process, so two runs never share one.
func NewID(at time.Time, pid int) string {
	return fmt.Sprintf("%s-%d", at.UTC().Format("20060102T150405"), pid)
}

func dirOf(root string) string { return filepath.Join(root, filepath.FromSlash(Dir)) }

func pathOf(root, id string) string { return filepath.Join(dirOf(root), id+".json") }

// Save writes the record, whole, through a temporary file: a reader never sees half of it.
func Save(root string, r Run) error {
	if err := os.MkdirAll(dirOf(root), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	tmp := pathOf(root, r.ID) + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, pathOf(root, r.ID))
}

// Load reads one record.
func Load(root, id string) (Run, error) {
	var r Run
	b, err := os.ReadFile(pathOf(root, id))
	if err != nil {
		return r, err
	}
	return r, json.Unmarshal(b, &r)
}

// List reads every record, oldest first. A file that does not read as a record is skipped:
// a record half-written by a process killed mid-write says nothing.
func List(root string) []Run {
	entries, err := os.ReadDir(dirOf(root))
	if err != nil {
		return nil
	}
	var out []Run
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".json") || name == stateFile {
			continue
		}
		r, err := Load(root, strings.TrimSuffix(name, ".json"))
		if err != nil || r.ID == "" {
			continue
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].Started.Equal(out[j].Started) {
			return out[i].Started.Before(out[j].Started)
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// Finish records how a run ended: when, its exit code when known, and its summary.
func Finish(root, id string, at time.Time, exit *int, summary string) error {
	r, err := Load(root, id)
	if err != nil {
		return err
	}
	r.Ended, r.Exit = &at, exit
	if summary != "" {
		r.Summary = summary
	}
	if r.State == "" || r.State == StateRunning || r.State == StateStalled {
		r.State = StateFinished
		if exit == nil {
			r.State = StateEnded
		}
	}
	return Save(root, r)
}

// Prune keeps the records of runs still going, and of the ended ones the latest `keep` that
// ended within `maxAge`; it removes the rest, and returns how many.
func Prune(root string, keep int, maxAge time.Duration, now time.Time) int {
	all := List(root)
	var ended []Run
	for _, r := range all {
		if r.Done() {
			ended = append(ended, r)
		}
	}
	sort.Slice(ended, func(i, j int) bool { return ended[i].Ended.After(*ended[j].Ended) })
	removed := 0
	for i, r := range ended {
		if (keep > 0 && i >= keep) || (maxAge > 0 && now.Sub(*r.Ended) > maxAge) {
			if os.Remove(pathOf(root, r.ID)) == nil {
				removed++
			}
		}
	}
	return removed
}

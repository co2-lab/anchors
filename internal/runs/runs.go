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
	"regexp"
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
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Dir     string `json:"dir,omitempty"`
	PID     int    `json:"pid,omitempty"`
	By      string `json:"by"`
	// Session is who launched the run: the identity `anchors progress` writes as a task's
	// holder, so the monitor can tell a held task's holder from its activity
	// (DESIGN-agent-activity.md).
	Session string `json:"session,omitempty"`
	// Agent is the subagent a run belongs to: the subagent itself, for a run of KindAgent,
	// or the one that launched a command. A subagent waiting on a command of its own is not
	// stalled (DESIGN-agent-activity.md).
	Agent   string     `json:"agent,omitempty"`
	Started time.Time  `json:"started"`
	Output  string     `json:"output,omitempty"`
	Ended   *time.Time `json:"ended,omitempty"`
	Exit    *int       `json:"exit,omitempty"`
	Summary string     `json:"summary,omitempty"`
	// State is the monitor's last reading: running, stalled, died, finished or ended.
	State string `json:"state,omitempty"`
}

// KindAgent is the kind of a subagent's run: no process of its own; it lives from the
// agent's SubagentStart hook to its SubagentStop, and moves while its transcript grows.
const KindAgent = "agent"

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

// seenDir holds, per session, when its agent last ran a command: an agent that works with
// quick commands only — an edit, a grep, a short test — starts no long run, and is busy all
// the same (reported from MIF).
const seenDir = "seen"

// MarkSeen records that session was active at the moment given.
func MarkSeen(root, session string, at time.Time) error {
	if session == "" {
		return nil
	}
	dir := filepath.Join(dirOf(root), seenDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	name := regexp.MustCompile(`[^A-Za-z0-9._-]`).ReplaceAllString(session, "_")
	b, _ := json.Marshal(map[string]any{"session": session, "at": at})
	tmp := filepath.Join(dir, name+".tmp")
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, name+".json"))
}

// LastSeen is when each session was last marked active.
func LastSeen(root string) map[string]time.Time {
	out := map[string]time.Time{}
	entries, _ := os.ReadDir(filepath.Join(dirOf(root), seenDir))
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var m struct {
			Session string    `json:"session"`
			At      time.Time `json:"at"`
		}
		b, err := os.ReadFile(filepath.Join(dirOf(root), seenDir, e.Name()))
		if err != nil || json.Unmarshal(b, &m) != nil || m.Session == "" {
			continue
		}
		if m.At.After(out[m.Session]) {
			out[m.Session] = m.At
		}
	}
	return out
}

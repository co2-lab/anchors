// @anchors
//   code: RNMNR
//   ref: MNTRS

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

// States of a run, as the monitor reads them.
const (
	StateRunning  = "running"
	StateStalled  = "stalled"
	StateFinished = "finished" // ended, with its exit or its summary saying how
	StateDied     = "died"     // gone without an exit record or a summary
	StateEnded    = "ended"    // gone, seen only in the process table: how is unknown
)

// stateFile is the monitor's own state, beside the records: one per watcher, so two
// monitors of the same tree — each agent its own — neither overwrite nor silence each
// other's memory. An unnamed watcher keeps `monitor.json`.
func stateFile(watcher string) string {
	if w := watcherChars.ReplaceAllString(watcher, "-"); w != "" {
		return "monitor." + w + ".json"
	}
	return "monitor.json"
}

var watcherChars = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// isStateFile says whether a file of the records' folder is a watcher's state, not a run.
func isStateFile(name string) bool {
	return name == "monitor.json" || strings.HasPrefix(name, "monitor.") && strings.HasSuffix(name, ".json")
}

// Timing is when the monitor looks, and when it speaks without an event.
type Timing struct {
	Every     time.Duration // how often it reads the processes
	Heartbeat time.Duration // how often it says what runs, even when nothing changed
	Progress  time.Duration // the least time between two progress lines of one run
	Stall     time.Duration // quiet time before a run of no declared runner is stalled
	// AgentStall is the quiet time of a subagent's transcript, with no command of its own
	// alive, before it is stalled: an agent thinks for minutes (the user, 2026-10-13).
	AgentStall time.Duration
}

// DefaultTiming comes from MIF's experience with its own monitors: a tick of 30 s caught
// every report within one tick, and a heartbeat of 5 min was frequent enough to catch an
// agent that stopped and sparse enough not to flood the conversation.
func DefaultTiming() Timing {
	return Timing{Every: 30 * time.Second, Heartbeat: 5 * time.Minute, Progress: 5 * time.Minute, Stall: 5 * time.Minute,
		AgentStall: 10 * time.Minute}
}

// Seen is what the monitor remembers of one run between ticks.
type Seen struct {
	PID          int       `json:"pid,omitempty"`
	CPU          float64   `json:"cpu"`
	OutSize      int64     `json:"out_size"`
	LastMove     time.Time `json:"last_move"`
	LastProgress time.Time `json:"last_progress"`
	Kids         string    `json:"kids,omitempty"`
	SawAlive     bool      `json:"saw_alive,omitempty"`
	Stalled      bool      `json:"stalled,omitempty"`
	Done         bool      `json:"done,omitempty"`
}

// State is a watcher's memory, kept in `.anchors/runs/monitor[.<watcher>].json`: a re-armed
// monitor resumes from it, saying what happened while its watcher was away, and nothing
// twice — what another watcher was told is not this one's.
type State struct {
	Cursor        time.Time            `json:"cursor"`
	LastHeartbeat time.Time            `json:"last_heartbeat"`
	Runs          map[string]*Seen     `json:"runs"`
	Reports       map[string]time.Time `json:"reports"`
	// ProgressSum is the progress summary as the monitor last said it, ProgressWas its total,
	// and LastSummary when: the summary goes out when it changes, with what it was
	// (DESIGN-monitor-signal.md), and again only after Settings.SummaryEvery.
	ProgressSum string    `json:"progress_sum,omitempty"`
	ProgressWas string    `json:"progress_was,omitempty"`
	LastSummary time.Time `json:"last_summary,omitempty"`
	// SaidIdle and SaidHigh are the heartbeat of the action level: the project falling idle,
	// and a high load, each said once per stretch.
	SaidIdle bool `json:"said_idle,omitempty"`
	// CPUTime and CPUAt are the CPU seconds of every process, and when they were read: the
	// next tick measures the CPUs' use from them.
	CPUTime  float64   `json:"cpu_time,omitempty"`
	CPUAt    time.Time `json:"cpu_at,omitempty"`
	SaidHigh bool      `json:"said_high,omitempty"`
	// IdleTasks are the tasks in progress said to have no activity from their holder, until
	// it moves again (DESIGN-agent-activity.md).
	IdleTasks map[string]bool `json:"idle_tasks,omitempty"`
	// Watching is when the monitor began to watch: a holder's silence is never counted from
	// before it — at MIF a first tick said "no activity for 11h08m" of agents at work.
	Watching time.Time `json:"watching,omitempty"`
}

// LoadState reads the watcher's state; a missing or unreadable one is a fresh start.
func LoadState(root, watcher string) *State {
	st := &State{}
	if b, err := os.ReadFile(filepath.Join(dirOf(root), stateFile(watcher))); err == nil {
		_ = json.Unmarshal(b, st)
	}
	if st.Runs == nil {
		st.Runs = map[string]*Seen{}
	}
	if st.Reports == nil {
		st.Reports = map[string]time.Time{}
	}
	return st
}

// SaveState writes the watcher's state, through a temporary file of its own: another
// monitor's rename cannot take it.
func SaveState(root, watcher string, st *State) error {
	if err := os.MkdirAll(dirOf(root), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dirOf(root), stateFile(watcher)+".*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	_, werr := f.Write(b)
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		os.Remove(tmp)
		return werr
	}
	return os.Rename(tmp, filepath.Join(dirOf(root), stateFile(watcher)))
}

// Report is a test report file as the monitor last read it.
type Report struct {
	Path     string
	Modified time.Time
	Tests    int
	Failures int
	First    string // the first failure: its test and message
}

// Snapshot is what one tick sees: the time, the processes, the records, the load, the
// reports and a way to read an output file. It is built by the caller, so the monitor's
// reading can be proven without a real process table.
type Snapshot struct {
	Now     time.Time
	Root    string
	Procs   []Proc
	Self    int // the monitor's process: it and what is above it are never a run
	Records []Run
	Load    float64
	HasLoad bool
	CPUs    int
	// Busy is the share of the CPUs in use since the last tick, from 0 to 100, and HasBusy
	// whether it was measured: the load average counts what waits — for a CPU, but on macOS
	// for a disk too — and stays high after a peak, so a load of 10 on 10 CPUs came with
	// the CPUs 65% busy (reported by the user).
	Busy    float64
	HasBusy bool
	Reports []Report
	Output  func(path string) OutputInfo
}

// Event is one line the monitor prints, and the run it is about.
type Event struct {
	Kind string
	Run  string
	Line string
}

// Tick reads one snapshot against what the monitor remembers, and returns the events it
// calls for — and the records to write: a run the monitor found or whose state changed.
func Tick(st *State, snap Snapshot, runners []Runner, t Timing) ([]Event, []Run) {
	tree := NewTree(snap.Procs)
	excluded := map[int]bool{snap.Self: true}
	for _, a := range tree.Ancestors(snap.Self) {
		excluded[a] = true
	}
	var events []Event
	var writes []Run

	// The runs: the records, and the runners the process table shows under the project
	// that no record accounts for.
	runs := map[string]*Run{}
	var order []string
	claimed := map[int]bool{}
	for i := range snap.Records {
		r := snap.Records[i]
		runs[r.ID] = &r
		order = append(order, r.ID)
		if s := st.Runs[r.ID]; s != nil && s.PID > 0 && r.PID == 0 {
			r.PID = s.PID
			runs[r.ID].PID = s.PID
		}
		if r.PID > 0 && !r.Done() {
			claimed[r.PID] = true
			for _, d := range tree.Descendants(r.PID) {
				claimed[d] = true
			}
		}
	}
	// A record of the agent has no process: the hook runs before the command starts. It takes
	// the topmost process under the project that runs the program its command names.
	for _, id := range order {
		r := runs[id]
		if r.Done() || r.PID > 0 || r.By != ByAgent || r.Kind == KindAgent {
			continue
		}
		if pid := findProcess(tree, snap.Root, r.Command, runners, excluded, claimed); pid > 0 {
			r.PID = pid
			writes = append(writes, *r)
			claimed[pid] = true
			for _, d := range tree.Descendants(pid) {
				claimed[d] = true
			}
		}
	}
	tops := projectRunners(tree, snap.Root, runners, excluded, claimed)
	for _, p := range tops {
		if p.PID == 0 || claimed[p.PID] {
			continue
		}
		rn, _ := RunnerFor(runners, p.Cmd)
		r := Run{ID: NewID(snap.Now, p.PID), Kind: rn.Kind, Name: rn.Name, Command: p.Cmd, Dir: p.Cwd,
			PID: p.PID, By: ByOS, Started: snap.Now, State: StateRunning}
		// A process the monitor already follows keeps its id.
		for id, s := range st.Runs {
			if s.PID == p.PID && !s.Done {
				if old, ok := runs[id]; ok {
					r = *old
				} else {
					r.ID = id
				}
				break
			}
		}
		if _, ok := runs[r.ID]; !ok {
			runs[r.ID] = &r
			order = append(order, r.ID)
		}
	}

	// The subagents waiting on a command of their own: while it is alive, the subagent's
	// silence is waiting, and the command is the one watched for a stall.
	waiting := map[string]bool{}
	for _, id := range order {
		r := runs[id]
		if r.Agent != "" && r.Kind != KindAgent && !r.Done() && (st.Runs[id] == nil || !st.Runs[id].Done) {
			waiting[r.Agent] = true
		}
	}

	for _, id := range order {
		r := runs[id]
		s := st.Runs[id]
		fresh := s == nil
		if fresh {
			if r.Done() && r.Ended.Before(st.Cursor) {
				continue // ended before the monitor's memory began: old news
			}
			s = &Seen{PID: r.PID, LastMove: snap.Now, LastProgress: snap.Now}
			st.Runs[id] = s
		}
		if s.Done {
			continue
		}
		rn, known := RunnerFor(runners, FirstLine(r.Command))
		if !known {
			rn = Runner{Name: firstWord(FirstLine(r.Command)), Kind: "other"}
		}
		if r.Name == "" {
			r.Name = rn.Name
		}
		if r.Kind == "" {
			r.Kind = rn.Kind
		}
		out := OutputInfo{}
		if r.Output != "" && snap.Output != nil {
			out = snap.Output(r.Output)
		}
		_, alive := tree.ByPID[r.PID]
		if r.PID > 0 {
			s.PID = r.PID
		}
		if alive || (out.OK && out.Size != s.OutSize) {
			s.SawAlive = true
		}
		if fresh {
			events = append(events, Event{"started", id, startedLine(*r)})
			if r.By == ByOS {
				writes = append(writes, *r)
			}
		}

		// A command of the agent no process runs any more, past its first tick, ended: once
		// its output — when the monitor knows it — stopped growing.
		gone := r.By == ByAgent && r.Kind != KindAgent && r.PID == 0 && !fresh && (!out.OK || out.Size == s.OutSize)
		if r.Done() || (r.PID > 0 && !alive) || gone {
			ev, w := endOf(*r, rn, out, snap.Now, s.SawAlive)
			events = append(events, ev)
			s.Done = true
			writes = append(writes, w)
			continue
		}

		// Moving: CPU advanced, the output grew, or the processes under it changed — a script
		// that waits on short commands uses almost no CPU of its own, and starts new ones.
		cpu, kids := 0.0, ""
		if r.PID > 0 {
			cpu = tree.TreeCPU(r.PID)
			kids = kidsKey(tree, r.PID)
		}
		moved := cpu > s.CPU || newKid(s.Kids, kids) || (out.OK && out.Size != s.OutSize)
		s.Kids = kids
		if fresh || (r.Kind == KindAgent && waiting[r.Agent]) {
			moved = true
		}
		s.CPU = cpu
		if out.OK {
			s.OutSize = out.Size
		}
		if moved {
			if s.Stalled {
				events = append(events, Event{"resumed", id, resumedLine(*r)})
				s.Stalled = false
				r.State = StateRunning
				writes = append(writes, *r)
			}
			s.LastMove = snap.Now
		} else if quiet := snap.Now.Sub(s.LastMove); r.Kind == KindAgent && !s.Stalled && quiet >= t.AgentStall {
			s.Stalled = true
			r.State = StateStalled
			writes = append(writes, *r)
			events = append(events, Event{"stalled", id, fmt.Sprintf("⏸ %s: its transcript quiet for %s, and no command of its own running", r.Name, round(quiet))})
		} else if quiet := snap.Now.Sub(s.LastMove); r.Kind != KindAgent && !s.Stalled && quiet >= stallOf(rn, t) {
			s.Stalled = true
			r.State = StateStalled
			writes = append(writes, *r)
			events = append(events, Event{"stalled", id, fmt.Sprintf("⏸ %s: no output or CPU for %s (pid %d) — %s", r.Name, round(quiet), r.PID, Brief(r.Command))})
		}
		if rn.Progress != nil && out.OK && snap.Now.Sub(s.LastProgress) >= t.Progress && !fresh {
			s.LastProgress = snap.Now
			done := len(rn.Progress.FindAllString(out.Tail, -1))
			if done > 0 {
				events = append(events, Event{"progress", id, fmt.Sprintf("… %s: %d done in the last part of its output, %s running", r.Name, done, round(snap.Now.Sub(r.Started)))})
			}
		}
	}

	// Reports that landed since the last tick.
	for _, rep := range snap.Reports {
		last, seen := st.Reports[rep.Path]
		st.Reports[rep.Path] = rep.Modified
		if seen && !rep.Modified.After(last) {
			continue
		}
		if !seen && !rep.Modified.After(st.Cursor) {
			continue // there before the monitor's memory began
		}
		// A report with no test is no result: a runner that empties its report when it starts
		// and writes it at the end (vitest) left a "✓ 0 test(s)" mid-run, which read as a pass
		// (reported from MIF). It is said once it holds tests.
		if rep.Tests == 0 {
			continue
		}
		mark := "✓"
		if rep.Failures > 0 {
			mark = "✗"
		}
		line := fmt.Sprintf("%s report %s: %d test(s), %d failure(s)", mark, rep.Path, rep.Tests, rep.Failures)
		if rep.First != "" {
			line += " — " + rep.First
		}
		events = append(events, Event{"report", rep.Path, line})
	}

	if t.Heartbeat > 0 && snap.Now.Sub(st.LastHeartbeat) >= t.Heartbeat {
		events = append(events, Event{"heartbeat", "", heartbeat(st, runs, order, snap)})
		st.LastHeartbeat = snap.Now
	}
	st.Cursor = snap.Now
	return events, writes
}

// projectRunners are the processes under the project that a runner recognizes, each the
// topmost of its kind — `npx jest` and the `node` workers it starts are one run.
func projectRunners(tree Tree, root string, runners []Runner, excluded, claimed map[int]bool) []Proc {
	var cand []Proc
	for _, p := range tree.ByPID {
		if excluded[p.PID] || claimed[p.PID] || anchorsItself.MatchString(p.Cmd) || launchesClaimed(tree, p.PID, claimed) {
			continue
		}
		if _, ok := RunnerFor(runners, p.Cmd); ok {
			cand = append(cand, p)
		}
	}
	isCand := map[int]bool{}
	for _, p := range cand {
		isCand[p.PID] = true
	}
	var pids []int
	for _, p := range cand {
		pids = append(pids, p.PID)
	}
	cwd := Cwds(pids)
	var out []Proc
	for _, p := range cand {
		under := false
		for _, a := range tree.Ancestors(p.PID) {
			if isCand[a] {
				under = true
				break
			}
		}
		if under {
			continue
		}
		if p.Cwd == "" {
			p.Cwd = cwd[p.PID]
		}
		if InProject(p, root) {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PID < out[j].PID })
	return out
}

// kidsKey names the processes under pid, so a tick tells whether they changed.
func kidsKey(tree Tree, pid int) string {
	d := tree.Descendants(pid)
	sort.Ints(d)
	parts := make([]string, len(d))
	for i, k := range d {
		parts[i] = fmt.Sprint(k)
	}
	return strings.Join(parts, ",")
}

// newKid says whether a process appeared under the run since the last tick: one that ended
// is no sign of life.
func newKid(before, now string) bool {
	had := map[string]bool{}
	for _, k := range strings.Split(before, ",") {
		had[k] = true
	}
	for _, k := range strings.Split(now, ",") {
		if k != "" && !had[k] {
			return true
		}
	}
	return false
}

// launchesClaimed says whether a process has, under it, a process a run already accounts
// for: it is that run's launcher — the shell an agent ran `anchors check` in names the same
// command, and was found as a second run that ended with its exit unknown (reported from MIF).
func launchesClaimed(tree Tree, pid int, claimed map[int]bool) bool {
	for _, d := range tree.Descendants(pid) {
		if claimed[d] {
			return true
		}
	}
	return false
}

// anchorsItself are Anchors' own wrapping and watching processes: their command lines name
// what they wrap or watch, and they are never a run of their own.
var anchorsItself = regexp.MustCompile(`(^|[\s/])anchors (run|monitor|watch)(\s|$)`)

// findProcess is the topmost process under the project, not yet claimed, whose command line
// runs the program a recorded command names — or the runner its first line is —; 0 when there
// is none.
func findProcess(tree Tree, root, command string, runners []Runner, excluded, claimed map[int]bool) int {
	first := FirstLine(command)
	want, isRunner := RunnerFor(runners, first)
	runs := func(p Proc) bool {
		if sharesProgram(first, p.Cmd) {
			return true
		}
		got, ok := RunnerFor(runners, p.Cmd)
		return isRunner && ok && got.Name == want.Name
	}
	var cand []int
	for _, p := range tree.ByPID {
		if !excluded[p.PID] && !claimed[p.PID] && !anchorsItself.MatchString(p.Cmd) && !launchesClaimed(tree, p.PID, claimed) && runs(p) {
			cand = append(cand, p.PID)
		}
	}
	sort.Ints(cand)
	isCand := map[int]bool{}
	for _, c := range cand {
		isCand[c] = true
	}
	cwd := Cwds(cand)
	for _, c := range cand {
		top := true
		for _, a := range tree.Ancestors(c) {
			if isCand[a] {
				top = false
				break
			}
		}
		p := tree.ByPID[c]
		if p.Cwd == "" {
			p.Cwd = cwd[c]
		}
		if top && InProject(p, root) {
			return c
		}
	}
	return 0
}

// The levels of what the monitor says (DESIGN-monitor-signal.md): everything, or only what
// calls for a reaction — each line wakes a turn of the agent watching.
const (
	LevelAll    = "all"
	LevelAction = "action"
)

// ForAction says whether an event goes out at the action level: a failure, a death, a stall
// and its end, a run the agent launched ending either way — it waits for it —, a report
// with failures; and the heartbeat only when it says something: the project fell idle, or
// the load is high, each once per stretch. idle and high are the project's state now.
func (st *State) ForAction(e Event, r *Run, idle, high bool) bool {
	switch e.Kind {
	case "died", "stalled", "resumed":
		return true
	case "report":
		return strings.HasPrefix(e.Line, "✗")
	case "finished", "ended":
		if strings.HasPrefix(e.Line, "✗") {
			return true
		}
		if r == nil || (r.By != ByAgent && r.By != ByRun) {
			return false
		}
		return !strings.Contains(e.Line, "before the monitor saw it running")
	case "heartbeat":
		keep := (idle && !st.SaidIdle) || (high && !st.SaidHigh)
		if keep {
			st.SaidIdle = st.SaidIdle || idle
			st.SaidHigh = st.SaidHigh || high
		}
		return keep
	}
	return false
}

// NoteState clears what the action level said once, when it stops being true: the project
// runs something again, the load came down.
func (st *State) NoteState(idle, high bool) {
	if !idle {
		st.SaidIdle = false
	}
	if !high {
		st.SaidHigh = false
	}
}

// ExitNothingToConfront is the exit of `anchors check` and `anchors verify` when no file of
// the change is theirs to confront — a commit of configuration, or of Anchors' own files.
const ExitNothingToConfront = 3

// NothingToConfront says whether a run is one of Anchors' checks that found nothing to
// confront: no failure, though its exit is not 0. Read as a failure, every commit of such
// files showed ✗ (reported from MIF).
func NothingToConfront(r Run) bool {
	return r.By == ByAnchors && r.Exit != nil && *r.Exit == ExitNothingToConfront &&
		(r.Name == "anchors check" || r.Name == "anchors verify")
}

// endOf is the line, and the record, of a run that ended: by its exit when it has one, by its
// output's summary when it does not, and as died when neither says how — unless nobody but
// the process table knew of it, and then how it ended is simply unknown.
func endOf(r Run, rn Runner, out OutputInfo, now time.Time, sawAlive bool) (Event, Run) {
	if r.Ended == nil {
		r.Ended = &now
	}
	elapsed := round(r.Ended.Sub(r.Started))
	summary := r.Summary
	if summary == "" && out.OK {
		summary = SummaryLine(rn, out.Tail)
	}
	r.Summary = summary
	switch {
	case r.Kind == KindAgent:
		// The SubagentStop hook says that it ended, and with what last message; not whether
		// it succeeded — the hook carries no verdict (measured on a real session).
		r.State = StateEnded
		return Event{"ended", r.ID, fmt.Sprintf("⏹ %s ended (%s)%s", r.Name, elapsed, suffix(r.Summary))}, r
	case r.Exit != nil:
		r.State = StateFinished
		mark := "✓"
		if *r.Exit != 0 {
			mark = "✗"
		}
		if NothingToConfront(r) {
			mark = "·"
			if summary == "" {
				summary = "nothing to confront"
			}
		}
		return Event{"finished", r.ID, fmt.Sprintf("%s %s ended (exit %d, %s)%s", mark, r.Name, *r.Exit, elapsed, suffix(summary))}, r
	case out.OK:
		if failed, ok := Verdict(rn, out.Tail); ok {
			r.State = StateFinished
			mark := "✓"
			if failed {
				mark = "✗"
			}
			return Event{"finished", r.ID, fmt.Sprintf("%s %s ended (%s)%s", mark, r.Name, elapsed, suffix(summary))}, r
		}
		r.State = StateDied
		return Event{"died", r.ID, fmt.Sprintf("☠ %s gone without a summary (%s) — last lines: %s", r.Name, elapsed, strings.Join(LastLines(out.Tail, 3), " ⏎ "))}, r
	case r.By == ByOS:
		r.State = StateEnded
		return Event{"ended", r.ID, fmt.Sprintf("⏹ %s ended (%s) — exit unknown: seen only in the process table", r.Name, elapsed)}, r
	case r.By == ByAgent && r.PID == 0 && !sawAlive:
		// Launched while no monitor ran, and gone before one did: old news, said as such —
		// not as a run that just ended after hours (reported from MIF).
		r.State = StateEnded
		return Event{"ended", r.ID, fmt.Sprintf("⏹ %s ended before the monitor saw it running (launched %s): %s", r.Name, r.Started.Local().Format("15:04"), Brief(r.Command))}, r
	case r.By == ByAgent && r.PID == 0:
		r.State = StateEnded
		return Event{"ended", r.ID, fmt.Sprintf("⏹ %s is not in the process table (%s) — ended, exit unknown: %s", r.Name, elapsed, Brief(r.Command))}, r
	default:
		r.State = StateDied
		return Event{"died", r.ID, fmt.Sprintf("☠ %s gone without recording how it ended (%s) — %s", r.Name, elapsed, Brief(r.Command))}, r
	}
}

func heartbeat(st *State, runs map[string]*Run, order []string, snap Snapshot) string {
	counts := map[string]int{}
	ended := 0
	for _, id := range order {
		s := st.Runs[id]
		if s == nil {
			continue
		}
		if s.Done {
			if r := runs[id]; r.Ended != nil && r.Ended.After(st.LastHeartbeat) {
				ended++
			}
			continue
		}
		name := runs[id].Name
		if s.Stalled {
			name += " (stalled)"
		}
		counts[name]++
	}
	var names []string
	for n, c := range counts {
		names = append(names, fmt.Sprintf("%s %d", n, c))
	}
	sort.Strings(names)
	running := "NOTHING"
	if len(names) > 0 {
		running = strings.Join(names, ", ")
	}
	line := fmt.Sprintf("%s running=[%s]", snap.Now.Local().Format("15:04"), running)
	if ended > 0 && !st.LastHeartbeat.IsZero() {
		line += fmt.Sprintf(" — %d run(s) ended since %s", ended, st.LastHeartbeat.Local().Format("15:04"))
	}
	// The load is said against the CPUs: a bare 8.4 says nothing to whoever does not know
	// the machine has 10 of them — busy, with nothing waiting —, nor that 12 would mean a queue.
	if snap.HasLoad {
		if snap.CPUs > 0 {
			line += fmt.Sprintf(" load=%.1f/%d", snap.Load, snap.CPUs)
		} else {
			line += fmt.Sprintf(" load=%.1f", snap.Load)
		}
	}
	if snap.HasBusy {
		line += fmt.Sprintf(" cpu=%.0f%%", snap.Busy)
	}
	if snap.Overloaded() {
		line += " (HIGH: processes wait for a CPU — timeouts may be load)"
	}
	return line
}

// BusyThreshold is how busy the CPUs must be, with the load past their number, for the machine
// to be overloaded.
const BusyThreshold = 90

// Overloaded says whether processes wait for a CPU: the load past the CPUs and, when measured,
// the CPUs at least BusyThreshold busy; or the load past twice the CPUs, whatever the CPUs'
// use. The load alone overstated it: it counts what waits for a disk too, and lags behind a
// peak. The use, measured from the processes' CPU seconds, errs low — a process that ends
// takes its seconds with it (measured against top: 16% for 34%) —, so a load that high is
// overload on its own.
func (snap Snapshot) Overloaded() bool {
	if !snap.HasLoad || snap.CPUs <= 0 || snap.Load <= float64(snap.CPUs) {
		return false
	}
	return !snap.HasBusy || snap.Busy >= BusyThreshold || snap.Load > 2*float64(snap.CPUs)
}

// MeasureBusy gives the share of the CPUs in use since the last reading, from the CPU seconds
// of every process now, and keeps the reading for the next; false on the first reading or
// when no time passed. A process that ended in between takes its seconds with it, so the
// share errs low, never high.
func (st *State) MeasureBusy(procs []Proc, now time.Time, cpus int) (float64, bool) {
	total := 0.0
	for _, p := range procs {
		total += p.CPU
	}
	prev, at := st.CPUTime, st.CPUAt
	st.CPUTime, st.CPUAt = total, now
	elapsed := now.Sub(at).Seconds()
	if at.IsZero() || elapsed <= 0 || cpus <= 0 || len(procs) == 0 {
		return 0, false
	}
	busy := (total - prev) / (elapsed * float64(cpus)) * 100
	return max(0, min(100, busy)), true
}

func resumedLine(r Run) string {
	if r.Kind == KindAgent {
		return fmt.Sprintf("▶ %s moving again", r.Name)
	}
	return fmt.Sprintf("▶ %s moving again (pid %d)", r.Name, r.PID)
}

func startedLine(r Run) string {
	who := map[string]string{ByAnchors: "by anchors", ByAgent: "by the agent", ByRun: "by anchors monitor run", ByOS: "found running"}[r.By]
	pid := ""
	if r.PID > 0 {
		pid = fmt.Sprintf(" (pid %d)", r.PID)
	}
	return fmt.Sprintf("▶ %s%s — %s: %s", r.Name, pid, who, Brief(r.Command))
}

func stallOf(rn Runner, t Timing) time.Duration {
	if rn.Stall > 0 {
		return rn.Stall
	}
	return t.Stall
}

func round(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	return d.Round(time.Minute).String()
}

func suffix(s string) string {
	if s == "" {
		return ""
	}
	return " — " + s
}

// launchers are the words a command line starts with that are not the program it runs.
var launchers = map[string]bool{"cd": true, "env": true, "npx": true, "npm": true, "yarn": true, "pnpm": true,
	"exec": true, "run": true, "sh": true, "bash": true, "zsh": true, "node": true, "time": true, "nohup": true,
	"sudo": true, "&&": true, "||": true, ";": true, "-c": true, "go": true, "python": true, "python3": true}

// FirstLine is a command's first line: what a heredoc or a script carries after it is no
// command to recognize — a file edit that mentions jest is no jest run (reported from MIF).
func FirstLine(cmd string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(cmd), "\n")
	return strings.TrimSpace(line)
}

// briefMax is how much of a command's first line a line of the monitor shows.
const briefMax = 100

// Brief is a command as a line of the monitor shows it: its first line, cut at briefMax, and
// how many lines it left out — a 60-line heredoc flooded the notification and cut the event
// before its useful part.
func Brief(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	first := FirstLine(cmd)
	if r := []rune(first); len(r) > briefMax {
		first = string(r[:briefMax]) + "…"
	}
	if more := strings.Count(cmd, "\n"); more > 0 {
		first += fmt.Sprintf(" [+%d lines]", more)
	}
	return first
}

// setup are the commands that prepare a command line rather than run its work.
var setup = map[string]bool{"cd": true, "export": true, "source": true, ".": true, "set": true, "pushd": true, "ulimit": true}

// segmentRE splits a command line where one command ends and the next begins.
var segmentRE = regexp.MustCompile(`&&|\|\||;|\n`)

// firstWord is the program a command line runs: the first word, past the launchers in front
// of it, of its first segment that is no setup — `cd app && npx jest` runs jest.
func firstWord(cmd string) string {
	for _, seg := range segmentRE.Split(cmd, -1) {
		words := strings.Fields(seg)
		if len(words) == 0 || setup[words[0]] {
			continue
		}
		for _, w := range words {
			if strings.Contains(w, "=") || launchers[w] || strings.HasPrefix(w, "-") {
				continue
			}
			return filepath.Base(strings.Trim(w, `"'`))
		}
	}
	return "command"
}

// sharesProgram says whether a process's command line runs the program a recorded command
// names.
func sharesProgram(recorded, proc string) bool {
	w := firstWord(recorded)
	if w == "command" || len(w) < 3 {
		return false
	}
	for _, f := range strings.Fields(proc) {
		if filepath.Base(strings.Trim(f, `"'`)) == w {
			return true
		}
	}
	return false
}

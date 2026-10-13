// @anchors
//   code: MCTMN
//   ref: MNCMD

package flow

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/progress"
	"github.com/co2-lab/anchors/internal/runs"
	"github.com/spf13/cobra"
)

// THE MONITOR. An agent that launches a long process waits on a watch that speaks only at
// the end; when the process dies or hangs, the end never comes (DESIGN-process-monitor.md).
// Anchors cannot wake the agent, but the agent's watch tool can, one notification per line:
// `anchors monitor` is the command it watches.
func newMonitorCmd() *cobra.Command {
	var root, every, heartbeat, progress, stall, summaryEvery, level, idle string
	var untilDone, once, summary bool
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Watch the project's long processes: one line when one starts, stalls, dies or ends",
		Long: `Reads, on a loop, the project's long processes — test suites, mutation runs, builds, the
commit hook, whatever the agent launched — from the operating system and from the records
their launchers write in .anchors/runs/, and prints one line for each event that calls for a
reaction:

  ▶ started    ⏸ stalled (no output nor CPU for the runner's threshold)    ▶ moving again
  ☠ died (gone without an exit or a summary)    ✓/✗ ended, with its exit or its summary
  ✓/✗ a test report landed, with its failures    … progress, sparse
  and a heartbeat with what runs — "running=[NOTHING]" included — and the load.

With --summary (or "summary: true" in the monitor block), it says the summary of
anchors.progress.md — each demand and its stages — when it changes, with what the total was;
--summary-every 30m repeats it unchanged at most that often. Off by default; --summary=false
silences a summary the block turns on.

A task in progress whose holder has no process alive and no activity for --idle (default 10m,
"idle:" in the block, off to silence it) is said once — "⏸ D01.T7 in progress — no activity
from fix-datas for 12m" — and again when it moves.

--level action (or "level: action") says only what calls for a reaction: failures, deaths,
stalls, the runs the agent launched ending, failing reports, a change of progress, and the
heartbeat only when the project falls idle or the load is high. For an agent that watches —
each line wakes a turn.

Run it inside your watch tool, after launching the long command in the background, and react
to each line. Its state stays in .anchors/runs/: re-armed, it says what happened while nobody
watched, and nothing twice.

Timing comes from the defaults, then the "monitor:" block of anchors.yaml, then the flags.`,
		Example: `  anchors monitor                     stream until stopped
  anchors monitor --until-done        exit when the runs it watches have all ended
  anchors monitor --every 10s --heartbeat 2m`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile))
			set, err := runs.Configure(cfg)
			if err != nil {
				return err
			}
			for _, f := range []struct {
				val string
				to  *time.Duration
			}{{every, &set.Timing.Every}, {heartbeat, &set.Timing.Heartbeat}, {progress, &set.Timing.Progress}, {stall, &set.Timing.Stall}} {
				if f.val == "" {
					continue
				}
				d, err := time.ParseDuration(f.val)
				if err != nil || d <= 0 {
					return fmt.Errorf("%q is not a positive duration", f.val)
				}
				*f.to = d
			}
			if cmd.Flags().Changed("summary") {
				set.Summary = summary
			}
			if summaryEvery != "" {
				d, err := time.ParseDuration(summaryEvery)
				if err != nil || d <= 0 {
					return fmt.Errorf("%q is not a positive duration", summaryEvery)
				}
				set.SummaryEvery = d
			}
			switch idle {
			case "":
			case "off":
				set.Idle = 0
			default:
				d, err := time.ParseDuration(idle)
				if err != nil || d <= 0 {
					return fmt.Errorf("--idle %q is not a positive duration, nor off", idle)
				}
				set.Idle = d
			}
			if level != "" {
				if err := runs.CheckLevel(level); err != nil {
					return fmt.Errorf("--level: %w", err)
				}
				set.Level = level
			}
			return runMonitor(cmd.OutOrStdout(), absRoot, set, untilDone, once)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&every, "every", "", "how often to read the processes (default 30s, or monitor.every)")
	cmd.Flags().StringVar(&heartbeat, "heartbeat", "", "how often to say what runs, even when nothing changed (default 5m)")
	cmd.Flags().StringVar(&progress, "progress", "", "least time between two progress lines of one run (default 5m)")
	cmd.Flags().StringVar(&stall, "stall", "", "quiet time before a run of no declared runner is stalled (default 5m)")
	cmd.Flags().BoolVar(&summary, "summary", false, "say the progress summary when it changes (default off, or monitor.summary)")
	cmd.Flags().StringVar(&summaryEvery, "summary-every", "", "repeat the progress summary unchanged at most this often (default never)")
	cmd.Flags().StringVar(&level, "level", "", "all (default), or action: only what calls for a reaction")
	cmd.Flags().StringVar(&idle, "idle", "", "how long a task in progress may go without activity from its holder before it is said (default 10m; off)")
	cmd.Flags().BoolVar(&untilDone, "until-done", false, "exit when the runs it watches have all ended")
	cmd.Flags().BoolVar(&once, "once", false, "read once, print what changed, and exit")
	cmd.AddCommand(newMonitorHookCmd(), newRunCmd())
	return cmd
}

// runMonitor is the loop: each tick reads a snapshot, prints the events and writes the
// records and the state. A tick that cannot read the processes says so — silence would read
// as "nothing changed".
func runMonitor(w io.Writer, root string, set runs.Settings, untilDone, once bool) error {
	st := runs.LoadState(root)
	sawRun := false
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	defer signal.Stop(stop)
	for {
		now := time.Now()
		snap := runs.Snapshot{Now: now, Root: root, Self: os.Getpid(), Records: runs.List(root),
			CPUs: runtime.NumCPU(), Output: runs.ReadOutput}
		procs, err := runs.ListProcs()
		if err != nil {
			fmt.Fprintf(w, "⚠ monitor: cannot read the process table (%v) — runs are judged by their records alone\n", err)
		}
		snap.Procs = procs
		snap.Load, snap.HasLoad = runs.LoadAverage()
		snap.Reports = runs.ReadReports(root, set.Reports, st.Reports)
		events, writes := runs.Tick(st, snap, set.Runners, set.Timing)
		for _, r := range writes {
			_ = runs.Save(root, r)
		}
		live := 0
		for _, s := range st.Runs {
			if !s.Done {
				live++
			}
		}
		idle := live == 0
		high := snap.HasLoad && snap.CPUs > 0 && snap.Load > float64(snap.CPUs)
		st.NoteState(idle, high)
		byID := map[string]*runs.Run{}
		for i := range snap.Records {
			byID[snap.Records[i].ID] = &snap.Records[i]
		}
		for i := range writes {
			byID[writes[i].ID] = &writes[i]
		}
		for _, e := range events {
			if set.Level == runs.LevelAction && !st.ForAction(e, byID[e.Run], idle, high) {
				continue
			}
			fmt.Fprintln(w, e.Line)
		}
		recs := append([]runs.Run(nil), snap.Records...)
		recs = append(recs, writes...)
		for _, l := range idleTaskLines(st, root, recs, procs, now, set.Idle) {
			fmt.Fprintln(w, l)
		}
		if set.Summary {
			for _, l := range progressLines(st, root, now, set.SummaryEvery) {
				fmt.Fprintln(w, l)
			}
		}
		runs.Prune(root, set.Keep, set.MaxAge, now)
		if err := runs.SaveState(root, st); err != nil {
			fmt.Fprintf(w, "⚠ monitor: cannot write its state (%v) — a re-armed monitor may repeat events\n", err)
		}
		if live > 0 {
			sawRun = true
		}
		if once {
			return nil
		}
		if untilDone && live == 0 {
			if !sawRun {
				fmt.Fprintln(w, "· nothing running in the project")
			} else {
				fmt.Fprintln(w, "· every run it watched has ended")
			}
			return nil
		}
		select {
		case <-stop:
			return nil
		case <-time.After(set.Timing.Every):
		}
	}
}

// SummaryInstruction is the line the monitor puts before the progress summary, apart from it:
// a summary the agent only reads never reaches the user, and an instruction inside the summary
// would be shown to the user with it (the user, 2026-10-12). It counts the lines and asks for a
// code block: asked to show them "as they come", agents rewrote them, and with one demand they
// folded the total into it and dropped its line (measured with MIF's summary, 3 of 3); asked
// for N lines verbatim in a code block, the total included, 4 of 4 kept it.
func SummaryInstruction(lines int) string {
	return fmt.Sprintf("· to the agent: copy the next %d lines to the user verbatim, inside a code block, the first line (the total) included — this line is for you, not for them", lines)
}

// progressLines are the summary of anchors.progress.md the monitor says
// (DESIGN-monitor-signal.md): when it changed since it last said it, with what its total was,
// and the same again only after every (0: never). A summary said each time the monitor spoke
// woke the watching agent every minute or two to repeat 9 unchanged lines (measured at MIF, a
// 3-hour session at 50%). A file edited by hand is one warning. No file, nothing.
func progressLines(st *runs.State, root string, now time.Time, every time.Duration) []string {
	if _, err := os.Stat(progress.Path(root)); err != nil {
		st.ProgressSum, st.ProgressWas = "", ""
		return nil
	}
	f, err := progress.Load(root)
	if err != nil {
		b, _ := os.ReadFile(progress.Path(root))
		h := sha256.Sum256(b)
		key := "edited " + hex.EncodeToString(h[:6])
		if key == st.ProgressSum {
			return nil
		}
		st.ProgressSum = key
		return []string{"⚠ progress: " + err.Error()}
	}
	var buf strings.Builder
	_ = printSummary(&buf, f, false)
	text := strings.TrimRight(buf.String(), "\n")
	h := sha256.Sum256([]byte(text))
	key := hex.EncodeToString(h[:6])
	changed := key != st.ProgressSum
	due := every > 0 && now.Sub(st.LastSummary) >= every
	if !changed && !due {
		return nil
	}
	done, total := progress.Totals(f.AllTasks())
	was := fmt.Sprintf("%s (%s/%s)", progress.Percent(done, total), num(done), num(total))
	lines := strings.Split(text, "\n")
	lines[0] = "▸ " + lines[0]
	if st.ProgressWas != "" && st.ProgressWas != was {
		lines[0] += " — was " + st.ProgressWas
	}
	st.ProgressSum, st.ProgressWas, st.LastSummary = key, was, now
	return append([]string{SummaryInstruction(len(lines))}, lines...)
}

// idleTaskLines confront the tasks in progress with their holders' activity
// (DESIGN-agent-activity.md): `[~]` says only that someone ran `start`. A holder is active while
// a run of theirs is alive; else their last activity is the latest of a run of theirs
// starting, moving or ending, and of a task of theirs started. Quiet for idle, the task is
// said once; moving again, it is said again. At MIF a teammate sat idle for hours, its task
// "in progress", until the user asked.
func idleTaskLines(st *runs.State, root string, recs []runs.Run, procs []runs.Proc, now time.Time, idle time.Duration) []string {
	if idle <= 0 {
		return nil
	}
	f, err := progress.Load(root)
	if err != nil {
		return nil
	}
	if st.Watching.IsZero() {
		st.Watching = now
	}
	last, live := map[string]time.Time{}, map[string]bool{}
	alive := map[int]bool{}
	for _, p := range procs {
		alive[p.PID] = true
	}
	seen := func(who string, at time.Time) {
		if at.After(last[who]) {
			last[who] = at
		}
	}
	// The other worktrees of the repository: agents work in their own, and their runs and
	// commands are recorded there.
	for _, dir := range gitmeta.Worktrees(root) {
		if same, _ := sameDir(dir, root); !same {
			recs = append(recs, runs.List(dir)...)
		}
		for who, at := range runs.LastSeen(dir) {
			seen(who, at)
		}
	}
	for _, r := range recs {
		if r.Session == "" {
			continue
		}
		seen(r.Session, r.Started)
		if r.Ended != nil {
			seen(r.Session, *r.Ended)
		}
		if s := st.Runs[r.ID]; s != nil {
			seen(r.Session, s.LastMove)
			if !s.Done && r.Ended == nil {
				live[r.Session] = true
			}
		} else if r.Ended == nil && r.PID > 0 && alive[r.PID] {
			// A run of another worktree, which this monitor does not follow: alive while its
			// process is (reported from MIF: a test of 10 minutes in a worktree read as silence).
			live[r.Session] = true
		}
		// Its output growing is activity, wherever the run is.
		if r.Output != "" {
			if fi, err := os.Stat(r.Output); err == nil {
				seen(r.Session, fi.ModTime())
			}
		}
	}
	var doing []*progress.Task
	for _, t := range f.AllTasks() {
		if t.State == progress.Doing {
			doing = append(doing, t)
			if at, ok := t.Since(); ok {
				seen(t.By, at)
			}
		}
	}
	if st.IdleTasks == nil {
		st.IdleTasks = map[string]bool{}
	}
	held := map[string]bool{}
	var out []string
	for _, t := range doing {
		held[t.ID] = true
		// Silence is counted from the last activity, never from before the monitor watched.
		from := last[t.By]
		if from.Before(st.Watching) {
			from = st.Watching
		}
		quiet := !live[t.By] && now.Sub(from) >= idle
		switch {
		case quiet && !st.IdleTasks[t.ID]:
			st.IdleTasks[t.ID] = true
			out = append(out, fmt.Sprintf("⏸ %s in progress — no activity from %s for %s", t.ID, t.By, minutes(now.Sub(from))))
		case !quiet && st.IdleTasks[t.ID]:
			delete(st.IdleTasks, t.ID)
			out = append(out, fmt.Sprintf("▶ %s moving again (%s)", t.ID, t.By))
		}
	}
	for id := range st.IdleTasks {
		if !held[id] {
			delete(st.IdleTasks, id)
		}
	}
	return out
}

// sameDir says whether two paths are the same folder, through symlinks.
func sameDir(a, b string) (bool, error) {
	ra, err := filepath.EvalSymlinks(a)
	if err != nil {
		return false, err
	}
	rb, err := filepath.EvalSymlinks(b)
	if err != nil {
		return false, err
	}
	return ra == rb, nil
}

// minutes says a duration to the minute: 12m, 2h05m.
func minutes(d time.Duration) string {
	m := int(d.Round(time.Minute) / time.Minute)
	if m < 60 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh%02dm", m/60, m%60)
}

// ownRunCommands are Anchors' long commands, which record their own runs: the monitor then
// knows their exit, and the commit hook's run (`verify`) is seen.
var ownRunCommands = map[string]bool{
	"test": true, "mutation": true, "check": true, "verify": true, "ingest": true,
	"map build": true, "docs build": true,
}

var ownRun *runs.Run
var ownRoot string

// BeginOwnRun records the run of one of Anchors' long commands, in a project that has an
// anchors.yaml; anything else records nothing. EndOwnRun closes it with the exit code.
func BeginOwnRun(cmd *cobra.Command) {
	path := strings.TrimPrefix(cmd.CommandPath(), cmd.Root().Name()+" ")
	if !ownRunCommands[path] {
		return
	}
	root := "."
	if f := cmd.Flags().Lookup("root"); f != nil && f.Value.String() != "" {
		root = f.Value.String()
	}
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return
	}
	if _, err := os.Stat(filepath.Join(absRoot, config.DefaultFile)); err != nil {
		return
	}
	now := time.Now()
	r := runs.Run{ID: runs.NewID(now, os.Getpid()), Kind: "anchors", Name: "anchors " + path,
		Command: strings.Join(os.Args, " "), Dir: absRoot, PID: os.Getpid(), By: runs.ByAnchors, Session: agentID(),
		Started: now, State: runs.StateRunning}
	if runs.Save(absRoot, r) == nil {
		ownRun, ownRoot = &r, absRoot
	}
	// The records do not wait for a monitor to be bounded.
	cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if set, err := runs.Configure(cfg); err == nil {
		runs.Prune(absRoot, set.Keep, set.MaxAge, now)
	}
}

// EndOwnRun records how the command's run ended.
func EndOwnRun(exit int) {
	if ownRun == nil {
		return
	}
	_ = runs.Finish(ownRoot, ownRun.ID, time.Now(), &exit, "")
	ownRun = nil
}

// newRunCmd wraps any command so its run is recorded: its process, its output and its exit.
func newRunCmd() *cobra.Command {
	var root, name string
	cmd := &cobra.Command{
		Use:   "run -- <command> [args...]",
		Short: "Run a command and record its run for the monitor: process, output and exit",
		Long: `Runs the command, passing its output through, and records the run in .anchors/runs/ —
its process, a copy of its output and its exit code —, so "anchors monitor" says how it
ended even when nobody else watched it. The command's exit code is this command's.`,
		Example: `  anchors monitor run -- npx jest --ci
  anchors monitor run --name seed -- ./scripts/seed.sh`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			code, err := recordedRun(absRoot, name, args, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			if err != nil {
				return err
			}
			if code != 0 {
				return common.ExitCode{Code: code}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&name, "name", "", "the run's name (default: the program it runs)")
	return cmd
}

// recordedRun runs the command with its run recorded, and returns its exit code.
func recordedRun(root, name string, args []string, in io.Reader, out, errOut io.Writer) (int, error) {
	now := time.Now()
	id := runs.NewID(now, os.Getpid())
	logPath := filepath.Join(root, filepath.FromSlash(runs.Dir), id+".log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return 0, err
	}
	logf, err := os.Create(logPath)
	if err != nil {
		return 0, err
	}
	defer logf.Close()
	c := exec.Command(args[0], args[1:]...)
	c.Dir, c.Stdin = root, in
	c.Stdout, c.Stderr = io.MultiWriter(out, logf), io.MultiWriter(errOut, logf)
	if err := c.Start(); err != nil {
		return 0, err
	}
	if name == "" {
		name = filepath.Base(args[0])
	}
	r := runs.Run{ID: id, Name: name, Command: strings.Join(args, " "), Dir: root, PID: c.Process.Pid,
		By: runs.ByRun, Session: agentID(), Started: now, Output: logPath, State: runs.StateRunning}
	_ = runs.Save(root, r)
	code := 0
	if err := c.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			code = 1
		}
	}
	_ = runs.Finish(root, id, time.Now(), &code, "")
	return code, nil
}

// THE AGENT HOOK. Claude Code runs a hook before and after each command of its agent; the
// one `install-hooks --agent` writes calls this. Before a command sent to the background, it
// records the run; after, the file the command's output goes to. It never blocks the agent:
// whatever happens, it exits 0 and prints nothing.
func newMonitorHookCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "hook",
		Short:  "Record the agent's long commands (called by the agent's hook)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, _ := io.ReadAll(bufio.NewReader(cmd.InOrStdin()))
			agentHook(b, time.Now())
			return nil
		},
	}
}

// hookInput is what the agent's hook receives on stdin, the fields this hook reads.
type hookInput struct {
	Event     string `json:"hook_event_name"`
	Tool      string `json:"tool_name"`
	ToolUseID string `json:"tool_use_id"`
	Cwd       string `json:"cwd"`
	Input     struct {
		Command    string `json:"command"`
		Background bool   `json:"run_in_background"`
		Timeout    int    `json:"timeout"`
	} `json:"tool_input"`
	Response json.RawMessage `json:"tool_response"`
	Output   json.RawMessage `json:"tool_output"`
	// The subagent the call comes from, and what its own hooks carry (measured on a real
	// session: SubagentStart has agent_id, agent_type and the session's transcript;
	// SubagentStop adds the subagent's transcript, its last message and the background tasks,
	// whose description names it; neither carries a verdict).
	AgentID         string `json:"agent_id"`
	AgentType       string `json:"agent_type"`
	Transcript      string `json:"transcript_path"`
	AgentTranscript string `json:"agent_transcript_path"`
	LastMessage     string `json:"last_assistant_message"`
	BackgroundTasks []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"background_tasks"`
}

// idChars are the characters a record's file name keeps of an id the harness gave.
var idChars = regexp.MustCompile(`[^A-Za-z0-9_-]`)

// subagentHook records a subagent as a run (DESIGN-agent-activity.md): SubagentStart opens it
// — no process of its own, its transcript as its output —, SubagentStop closes it with the
// first line of its last message, named by its task's description when the hook carries it.
func subagentHook(in hookInput, now time.Time) {
	root := projectRootFrom(in.Cwd)
	if root == "" || in.AgentID == "" {
		return
	}
	id := "subagent-" + idChars.ReplaceAllString(in.AgentID, "")
	switch in.Event {
	case "SubagentStart":
		r := runs.Run{ID: id, Kind: runs.KindAgent, Name: "agent " + in.AgentType, Command: "agent " + in.AgentType,
			Dir: in.Cwd, By: runs.ByAgent, Session: agentID(), Agent: in.AgentID, Started: now, State: runs.StateRunning}
		// The subagent's transcript, beside the session's: <session>/subagents/agent-<id>.jsonl.
		if in.Transcript != "" {
			r.Output = strings.TrimSuffix(in.Transcript, ".jsonl") + "/subagents/agent-" + in.AgentID + ".jsonl"
		}
		_ = runs.Save(root, r)
	case "SubagentStop":
		r, err := runs.Load(root, id)
		if err != nil {
			return
		}
		for _, t := range in.BackgroundTasks {
			if t.ID == in.AgentID && t.Description != "" {
				r.Name = "agent " + strconv.Quote(t.Description)
			}
		}
		if in.AgentTranscript != "" {
			r.Output = in.AgentTranscript
		}
		_ = runs.Save(root, r)
		_ = runs.Finish(root, id, now, nil, briefMessage(in.LastMessage))
	}
}

// briefMessage is a message's first line, cut to fit a line of the monitor.
func briefMessage(s string) string {
	s = strings.TrimSpace(runs.FirstLine(strings.TrimSpace(s)))
	if r := []rune(s); len(r) > 100 {
		s = string(r[:100]) + "…"
	}
	return s
}

// outputPathRE finds, in what the harness answered for a background command, the file its
// output goes to: the answer says it in words ("Output is being written to: <path>").
var outputPathRE = regexp.MustCompile(`written to:?\s*([^\s"\\]+)`)

// agentHook records, from one hook call, what it can of the agent's command.
func agentHook(raw []byte, now time.Time) {
	var in hookInput
	if json.Unmarshal(raw, &in) != nil {
		return
	}
	if in.Event == "SubagentStart" || in.Event == "SubagentStop" {
		subagentHook(in, now)
		return
	}
	if in.Tool != "Bash" || in.ToolUseID == "" {
		return
	}
	root := projectRootFrom(in.Cwd)
	if root == "" {
		return
	}
	id := "agent-" + idChars.ReplaceAllString(in.ToolUseID, "")
	switch in.Event {
	case "PreToolUse":
		// Every command, quick ones too, says its agent is at work.
		_ = runs.MarkSeen(root, agentID(), now)
		cmdLine := in.Input.Command
		// Only a command sent to the background is a run of the agent's: a foreground one
		// blocks the agent until it ends, and a long one is seen in the process table. Every
		// foreground command recorded — a file edit, a grep — left a run behind when the
		// command failed and no hook came after it (reported from MIF).
		if !in.Input.Background || anchorsOwn.MatchString(runs.FirstLine(cmdLine)) {
			return
		}
		name := ""
		if rn, ok := runs.RunnerFor(runs.DefaultRunners(), runs.FirstLine(cmdLine)); ok {
			name = rn.Name
		}
		_ = runs.Save(root, runs.Run{ID: id, Name: name, Command: cmdLine, Dir: in.Cwd, By: runs.ByAgent, Session: agentID(), Agent: in.AgentID, Started: now, State: runs.StateRunning})
	case "PostToolUse":
		r, err := runs.Load(root, id)
		if err != nil {
			return
		}
		if m := outputPathRE.FindStringSubmatch(string(in.Response) + string(in.Output)); m != nil {
			r.Output = strings.TrimRight(m[1], ".,;")
			_ = runs.Save(root, r)
		}
	}
}

// anchorsOwn are the agent's commands that need no record of the hook: Anchors' long
// commands record their own runs, and the monitor watches rather than runs.
var anchorsOwn = regexp.MustCompile(`(^|[\s/;&(])anchors (test|mutation|check|verify|ingest|map build|docs build|monitor|run)(\s|$)`)

// projectRootFrom is the nearest folder, from dir up, that holds an anchors.yaml.
func projectRootFrom(dir string) string {
	if dir == "" {
		return ""
	}
	for d := filepath.Clean(dir); ; d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, config.DefaultFile)); err == nil {
			return d
		}
		if filepath.Dir(d) == d {
			return ""
		}
	}
}

// @anchors
//   code: MCTMN
//   ref: MNCMD

package flow

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/runs"
	"github.com/spf13/cobra"
)

// THE MONITOR. An agent that launches a long process waits on a watch that speaks only at
// the end; when the process dies or hangs, the end never comes (DESIGN-process-monitor.md).
// Anchors cannot wake the agent, but the agent's watch tool can, one notification per line:
// `anchors monitor` is the command it watches.
func newMonitorCmd() *cobra.Command {
	var root, every, heartbeat, progress, stall string
	var untilDone, once bool
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
			return runMonitor(cmd.OutOrStdout(), absRoot, set, untilDone, once)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&every, "every", "", "how often to read the processes (default 30s, or monitor.every)")
	cmd.Flags().StringVar(&heartbeat, "heartbeat", "", "how often to say what runs, even when nothing changed (default 5m)")
	cmd.Flags().StringVar(&progress, "progress", "", "least time between two progress lines of one run (default 5m)")
	cmd.Flags().StringVar(&stall, "stall", "", "quiet time before a run of no declared runner is stalled (default 5m)")
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
		for _, e := range events {
			fmt.Fprintln(w, e.Line)
		}
		runs.Prune(root, set.Keep, set.MaxAge, now)
		if err := runs.SaveState(root, st); err != nil {
			fmt.Fprintf(w, "⚠ monitor: cannot write its state (%v) — a re-armed monitor may repeat events\n", err)
		}
		live := 0
		for _, s := range st.Runs {
			if !s.Done {
				live++
			}
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
		Command: strings.Join(os.Args, " "), Dir: absRoot, PID: os.Getpid(), By: runs.ByAnchors,
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
		By: runs.ByRun, Started: now, Output: logPath, State: runs.StateRunning}
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
}

// outputPathRE finds, in what the harness answered for a background command, the file its
// output goes to: the answer says it in words ("Output is being written to: <path>").
var outputPathRE = regexp.MustCompile(`written to:?\s*([^\s"\\]+)`)

// agentHook records, from one hook call, what it can of the agent's command.
func agentHook(raw []byte, now time.Time) {
	var in hookInput
	if json.Unmarshal(raw, &in) != nil || in.Tool != "Bash" || in.ToolUseID == "" {
		return
	}
	root := projectRootFrom(in.Cwd)
	if root == "" {
		return
	}
	id := "agent-" + regexp.MustCompile(`[^A-Za-z0-9_-]`).ReplaceAllString(in.ToolUseID, "")
	switch in.Event {
	case "PreToolUse":
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
		_ = runs.Save(root, runs.Run{ID: id, Name: name, Command: cmdLine, Dir: in.Cwd, By: runs.ByAgent, Started: now, State: runs.StateRunning})
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

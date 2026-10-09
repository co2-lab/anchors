// @anchors
//   code: MNTRT
//   ref: MNTRS

package runs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
)

// The process ids of the synthetic tables are high, so no real process shares one.
const (
	pidAgent = 900001
	pidSelf  = 900002
	pidNpx   = 900010
	pidNode  = 900011
	pidOther = 900020
	pidWrap  = 900030
	pidKid   = 900040
)

var root = filepath.Join(string(filepath.Separator), "work", "proj")

func tick(st *State, now time.Time, procs []Proc, recs []Run, out map[string]OutputInfo) ([]Event, []Run) {
	snap := Snapshot{Now: now, Root: root, Procs: procs, Self: pidSelf, Records: recs, CPUs: 4,
		Output: func(p string) OutputInfo { return out[p] }}
	return Tick(st, snap, DefaultRunners(), DefaultTiming())
}

func lines(evs []Event) string {
	var b []string
	for _, e := range evs {
		b = append(b, e.Line)
	}
	return strings.Join(b, "\n")
}

func kinds(evs []Event, kind string) int {
	n := 0
	for _, e := range evs {
		if e.Kind == kind {
			n++
		}
	}
	return n
}

func agentAndSelf() []Proc {
	return []Proc{{PID: pidAgent, PPID: 1, Cmd: "claude", Cwd: root}, {PID: pidSelf, PPID: pidAgent, Cmd: "anchors monitor", Cwd: root}}
}

func TestTick_findsRunnersUnderTheProject(t *testing.T) {
	t.Run("MNTRS-B01: A runner under the project that no record accounts for is a run found running", func(t *testing.T) {})
	procs := append(agentAndSelf(),
		Proc{PID: pidNpx, PPID: pidAgent, Cmd: "npm exec jest --ci", Cwd: root},
		Proc{PID: pidNode, PPID: pidNpx, Cmd: "node node_modules/.bin/jest --ci", Cwd: root},
		Proc{PID: pidOther, PPID: 1, Cmd: "node jest", Cwd: root + "-other"},
		Proc{PID: pidWrap, PPID: pidAgent, Cmd: "anchors monitor run -- ./bin/jest", Cwd: root},
		// A shell the agent typed `anchors check` in, and the check, which records its own run.
		Proc{PID: pidKid, PPID: pidAgent, Cmd: "/bin/zsh -c anchors check --changed a.go", Cwd: root},
		Proc{PID: pidKid + 1, PPID: pidKid, Cmd: "anchors check --changed a.go", Cwd: root},
	)
	st := &State{Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	check := Run{ID: "check", Name: "anchors check", Command: "anchors check --changed a.go", PID: pidKid + 1, By: ByAnchors, Started: t0}
	st.Runs["check"] = &Seen{PID: pidKid + 1, LastMove: t0}
	evs, writes := tick(st, t0, procs, []Run{check}, nil)
	if kinds(evs, "started") != 1 || len(writes) != 1 || writes[0].PID != pidNpx || writes[0].By != ByOS || writes[0].Name != "jest" {
		t.Errorf("one run, the topmost npx jest, by the process table:\n%s\n%+v", lines(evs), writes)
	}
}

func TestTick_stalledThenMoving(t *testing.T) {
	t.Run("MNTRS-B02: A quiet run is stalled once, and moving again is said", func(t *testing.T) {})
	st := &State{Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	rec := Run{ID: "r", Name: "jest", Command: "npx jest", PID: pidNpx, By: ByRun, Started: t0}
	base := append(agentAndSelf(), Proc{PID: pidNpx, PPID: pidAgent, CPU: 5, Cmd: "npx jest", Cwd: root})
	withKid := append(append([]Proc{}, base...), Proc{PID: pidKid, PPID: pidNpx, Cmd: "sleep 1", Cwd: root})
	tick(st, t0, withKid, []Run{rec}, nil)
	// The child ends meanwhile: no sign of life.
	evs, _ := tick(st, t0.Add(11*time.Minute), base, []Run{rec}, nil)
	if kinds(evs, "stalled") != 1 {
		t.Fatalf("stalled past jest's 10 minutes, a child ending is no life:\n%s", lines(evs))
	}
	if evs, _ := tick(st, t0.Add(12*time.Minute), base, []Run{rec}, nil); kinds(evs, "stalled") != 0 {
		t.Errorf("said once:\n%s", lines(evs))
	}
	moved := append(append([]Proc{}, base...), Proc{PID: pidKid + 1, PPID: pidNpx, Cmd: "node worker", Cwd: root})
	if evs, _ := tick(st, t0.Add(13*time.Minute), moved, []Run{rec}, nil); kinds(evs, "resumed") != 1 {
		t.Errorf("a new process under it is life:\n%s", lines(evs))
	}
}

func TestTick_endings(t *testing.T) {
	t.Run("MNTRS-B03: A run that ended is said once, by its exit, its summary, or nothing", func(t *testing.T) {})
	exit1 := 1
	end := t0.Add(time.Minute)
	recs := []Run{
		{ID: "exit", Name: "go test", Command: "go test ./...", By: ByRun, Started: t0, Ended: &end, Exit: &exit1},
		{ID: "sum", Name: "jest", Command: "npx jest", PID: pidNpx, By: ByRun, Started: t0, Output: "sum.out"},
		{ID: "cut", Name: "jest", Command: "npx jest", PID: pidNpx + 1, By: ByRun, Started: t0, Output: "cut.out"},
		{ID: "silent", Name: "anchors test", Command: "anchors test", PID: pidNpx + 2, By: ByAnchors, Started: t0},
		{ID: "os", Name: "jest", Command: "node jest", PID: pidNpx + 3, By: ByOS, Started: t0},
	}
	alive := append(agentAndSelf(),
		Proc{PID: pidNpx, PPID: 1, Cmd: "npx jest"}, Proc{PID: pidNpx + 1, PPID: 1, Cmd: "npx jest"},
		Proc{PID: pidNpx + 2, PPID: 1, Cmd: "anchors test"}, Proc{PID: pidNpx + 3, PPID: 1, Cmd: "node jest"})
	out := map[string]OutputInfo{
		"sum.out": {OK: true, Size: 10, Tail: "PASS a\nTests:       1 failed, 3 passed, 4 total\n"},
		"cut.out": {OK: true, Size: 10, Tail: "PASS a\nPASS b\n"},
	}
	st := &State{Cursor: t0.Add(-time.Hour), Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	recs[0].Ended = nil
	tick(st, t0, alive, recs[:5], out)
	recs[0].Ended = &end
	evs, _ := tick(st, t0.Add(2*time.Minute), agentAndSelf(), recs, out)
	got := lines(evs)
	for _, want := range []string{"✗ go test ended (exit 1", "✗ jest ended (2m0s) — Tests:       1 failed", "☠ jest gone without a summary", "PASS a ⏎ PASS b",
		"☠ anchors test gone without recording how it ended", "⏹ jest ended (2m0s) — exit unknown"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestTick_agentRecordTakesItsProcess(t *testing.T) {
	t.Run("MNTRS-B04: A command of the agent takes the process that runs its program", func(t *testing.T) {})
	st := &State{Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	rec := Run{ID: "agent-1", Command: "cd app && npx jest --ci", By: ByAgent, Started: t0}
	procs := append(agentAndSelf(), Proc{PID: pidNpx, PPID: pidAgent, Cmd: "npm exec jest --ci", Cwd: filepath.Join(root, "app")},
		Proc{PID: pidNode, PPID: pidNpx, Cmd: "node jest", Cwd: filepath.Join(root, "app")})
	evs, writes := tick(st, t0, procs, []Run{rec}, nil)
	if kinds(evs, "started") != 1 || len(writes) == 0 || writes[0].PID != pidNpx {
		t.Fatalf("the record takes the npx jest, and no run is found apart:\n%s\n%+v", lines(evs), writes)
	}
	// A command no process runs: it ends once its output stops growing.
	quiet := Run{ID: "agent-2", Command: "./scripts/seed.sh", By: ByAgent, Started: t0, Output: "seed.out"}
	out := map[string]OutputInfo{"seed.out": {OK: true, Size: 5, Tail: "seeding\n"}}
	tick(st, t0, agentAndSelf(), []Run{quiet}, out)
	out["seed.out"] = OutputInfo{OK: true, Size: 9, Tail: "seeding\ndone\n"}
	if evs, _ := tick(st, t0.Add(30*time.Second), agentAndSelf(), []Run{quiet}, out); kinds(evs, "died")+kinds(evs, "finished")+kinds(evs, "ended") != 0 {
		t.Errorf("still writing, still running:\n%s", lines(evs))
	}
	if evs, _ := tick(st, t0.Add(time.Minute), agentAndSelf(), []Run{quiet}, out); kinds(evs, "died") != 1 {
		t.Errorf("no process and its output stopped: it ended, without a summary:\n%s", lines(evs))
	}
}

func TestTick_reports(t *testing.T) {
	t.Run("MNTRS-B05: A report that landed is said with its failures", func(t *testing.T) {})
	st := &State{Cursor: t0, Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	snap := Snapshot{Now: t0.Add(time.Minute), Root: root, Self: pidSelf, Reports: []Report{
		{Path: "old.xml", Modified: t0.Add(-time.Hour), Tests: 3},
		{Path: "new.xml", Modified: t0.Add(30 * time.Second), Tests: 4, Failures: 1, First: "AuthScreen-A01"},
	}}
	evs, _ := Tick(st, snap, DefaultRunners(), Timing{})
	if kinds(evs, "report") != 1 || !strings.Contains(lines(evs), "✗ report new.xml: 4 test(s), 1 failure(s) — AuthScreen-A01") {
		t.Errorf("only the new report:\n%s", lines(evs))
	}
}

func TestTick_heartbeat(t *testing.T) {
	t.Run("MNTRS-B06: The heartbeat says what runs, what ended, and the load", func(t *testing.T) {})
	end := t0.Add(2 * time.Minute)
	st := &State{LastHeartbeat: t0, Runs: map[string]*Seen{
		"a": {PID: pidNpx, LastMove: t0.Add(5 * time.Minute)}, "b": {PID: pidNpx + 1, Stalled: true, LastMove: t0}, "c": {Done: true},
	}, Reports: map[string]time.Time{}}
	recs := []Run{
		{ID: "a", Name: "jest", Command: "npx jest", PID: pidNpx, By: ByRun, Started: t0},
		{ID: "b", Name: "maestro", Command: "maestro test", PID: pidNpx + 1, By: ByRun, Started: t0},
		{ID: "c", Name: "go test", Command: "go test", By: ByRun, Started: t0, Ended: &end},
	}
	procs := append(agentAndSelf(), Proc{PID: pidNpx, PPID: 1, Cmd: "npx jest"}, Proc{PID: pidNpx + 1, PPID: 1, Cmd: "maestro test"})
	snap := Snapshot{Now: t0.Add(6 * time.Minute), Root: root, Self: pidSelf, Procs: procs, Records: recs, Load: 9, HasLoad: true, CPUs: 4}
	evs, _ := Tick(st, snap, DefaultRunners(), DefaultTiming())
	got := lines(evs)
	if !strings.Contains(got, "running=[jest 1, maestro (stalled) 1]") || !strings.Contains(got, "1 run(s) ended since") || !strings.Contains(got, "load=9.0 (HIGH") {
		t.Errorf("names, ended and load:\n%s", got)
	}
	st2 := &State{Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	evs, _ = Tick(st2, Snapshot{Now: t0, Root: root, Self: pidSelf}, DefaultRunners(), DefaultTiming())
	if !strings.Contains(lines(evs), "running=[NOTHING]") {
		t.Errorf("nothing running is said:\n%s", lines(evs))
	}
}

func TestTick_rearmed(t *testing.T) {
	t.Run("MNTRS-B07: A re-armed monitor says what happened meanwhile, once", func(t *testing.T) {})
	dir := t.TempDir()
	rec := Run{ID: "r", Name: "jest", Command: "npx jest", PID: pidNpx, By: ByRun, Started: t0}
	st := LoadState(dir)
	tick(st, t0, append(agentAndSelf(), Proc{PID: pidNpx, PPID: 1, Cmd: "npx jest"}), []Run{rec}, nil)
	if err := SaveState(dir, st); err != nil {
		t.Fatal(err)
	}
	again := LoadState(dir)
	evs, _ := tick(again, t0.Add(10*time.Minute), agentAndSelf(), []Run{rec}, nil)
	if kinds(evs, "died") != 1 || kinds(evs, "started") != 0 {
		t.Fatalf("the death meanwhile is said, the start not again:\n%s", lines(evs))
	}
	if evs, _ := tick(again, t0.Add(11*time.Minute), agentAndSelf(), []Run{rec}, nil); kinds(evs, "died") != 0 {
		t.Errorf("said once:\n%s", lines(evs))
	}
	old := t0.Add(-2 * time.Hour)
	fresh := &State{Cursor: t0, Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	if evs, _ := tick(fresh, t0.Add(time.Minute), agentAndSelf(), []Run{{ID: "o", Command: "x", Started: old, Ended: &old}}, nil); len(evs) != 1 || evs[0].Kind != "heartbeat" {
		t.Errorf("a run ended before the memory began is old news:\n%s", lines(evs))
	}
}

func TestTick_progressIsSparse(t *testing.T) {
	t.Run("MNTRS-B08: Progress is sparse", func(t *testing.T) {})
	st := &State{Runs: map[string]*Seen{}, Reports: map[string]time.Time{}}
	rec := Run{ID: "r", Name: "jest", Command: "npx jest", PID: pidNpx, By: ByRun, Started: t0, Output: "o"}
	procs := append(agentAndSelf(), Proc{PID: pidNpx, PPID: 1, Cmd: "npx jest"})
	out := map[string]OutputInfo{"o": {OK: true, Size: 1, Tail: "PASS a\nFAIL b\n"}}
	tick(st, t0, procs, []Run{rec}, out)
	out["o"] = OutputInfo{OK: true, Size: 2, Tail: "PASS a\nFAIL b\nPASS c\n"}
	if evs, _ := tick(st, t0.Add(time.Minute), procs, []Run{rec}, out); kinds(evs, "progress") != 0 {
		t.Errorf("within the interval, no progress line:\n%s", lines(evs))
	}
	evs, _ := tick(st, t0.Add(6*time.Minute), procs, []Run{rec}, out)
	if kinds(evs, "progress") != 1 || !strings.Contains(lines(evs), "jest: 3 done") {
		t.Errorf("past it, one line counting the marks:\n%s", lines(evs))
	}
}

func TestConfigure(t *testing.T) {
	t.Run("MNTRS-B09: The settings are the defaults, then the project's block", func(t *testing.T) {})
	s, err := Configure(nil)
	if err != nil || s.Timing != DefaultTiming() || s.Keep != 50 || s.MaxAge != 7*24*time.Hour || len(s.Runners) != len(DefaultRunners()) {
		t.Fatalf("the defaults: %+v %v", s, err)
	}
	cfg := &config.Config{Tests: []config.Suite{{JUnit: "out/junit.xml"}}, Monitor: &config.Monitor{
		Every: "10s", Heartbeat: "2m", Keep: 5, KeepDays: 1, Reports: []string{"e2e/*.xml"},
		Runners: []config.MonitorRunner{{Name: "jest", Kind: "unit", Match: "craco test", Stall: "1m"}, {Name: "seed", Match: "seed\\.sh"}},
	}}
	s, err = Configure(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if s.Timing.Every != 10*time.Second || s.Timing.Heartbeat != 2*time.Minute || s.Timing.Stall != 5*time.Minute || s.Keep != 5 || s.MaxAge != 24*time.Hour {
		t.Errorf("the block's timing and bounds: %+v", s)
	}
	if s.Runners[0].Name != "jest" || s.Runners[0].Match.String() != "craco test" || s.Runners[1].Name != "seed" || s.Runners[1].Kind != "other" || len(s.Runners) != len(DefaultRunners())+1 {
		t.Errorf("the project's runners first, jest replaced: %+v", s.Runners[:2])
	}
	if strings.Join(s.Reports, " ") != "out/junit.xml e2e/*.xml" {
		t.Errorf("the suites' reports and the block's: %v", s.Reports)
	}
	for _, bad := range []config.Monitor{{Every: "soon"}, {Runners: []config.MonitorRunner{{Name: "x", Match: "("}}}, {Runners: []config.MonitorRunner{{Name: "y"}}}} {
		b := bad
		if _, err := Configure(&config.Config{Monitor: &b}); err == nil {
			t.Errorf("%+v does not read", bad)
		}
	}
}

func TestConfigure_badKeyNamed(t *testing.T) {
	t.Run("MNTRS-E01: A monitor value that does not read is an error naming its key", func(t *testing.T) {})
	if _, err := Configure(&config.Config{Monitor: &config.Monitor{Heartbeat: "-1m"}}); err == nil || !strings.Contains(err.Error(), "monitor.heartbeat") {
		t.Errorf("the key is named: %v", err)
	}
}

func TestOutputs(t *testing.T) {
	t.Run("MNTRS-B10: An output's end gives its verdict, its summary and its last lines", func(t *testing.T) {})
	jest, _ := RunnerFor(DefaultRunners(), "npx jest")
	if failed, ok := Verdict(jest, "Tests:       2 failed, 8 passed, 10 total\n"); !failed || !ok {
		t.Error("a failure beside a pass is failed")
	}
	if failed, ok := Verdict(jest, "Tests:       8 passed, 8 total\n"); failed || !ok {
		t.Error("a pass alone is passed")
	}
	if _, ok := Verdict(jest, "PASS a\n"); ok {
		t.Error("neither says nothing")
	}
	if got := SummaryLine(jest, "x\nTests:       8 passed, 8 total\ny\n"); got != "Tests:       8 passed, 8 total" {
		t.Errorf("summary line: %q", got)
	}
	if got := LastLines("a\n\nb\nc\n\n", 2); strings.Join(got, "|") != "b|c" {
		t.Errorf("last lines: %v", got)
	}
	p := filepath.Join(t.TempDir(), "o.log")
	_ = os.WriteFile(p, []byte("hello\n"), 0o644)
	if o := ReadOutput(p); !o.OK || o.Size != 6 || o.Tail != "hello\n" {
		t.Errorf("size and end: %+v", o)
	}
	if o := ReadOutput(""); o.OK {
		t.Error("no path, nothing read")
	}
}

func TestReadReports(t *testing.T) {
	t.Run("MNTRS-B11: Reports are read by their globs", func(t *testing.T) {})
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "out"), 0o755)
	xml := `<testsuites><testsuite name="s"><testcase name="ok"/><testcase name="AuthScreen-A01"><failure message="x"/></testcase></testsuite></testsuites>`
	_ = os.WriteFile(filepath.Join(dir, "out", "junit.xml"), []byte(xml), 0o644)
	reps := ReadReports(dir, []string{"out/*.xml"}, map[string]time.Time{})
	if len(reps) != 1 || reps[0].Path != "out/junit.xml" || reps[0].Tests != 2 || reps[0].Failures != 1 || reps[0].First != "AuthScreen-A01" {
		t.Fatalf("parsed: %+v", reps)
	}
	again := ReadReports(dir, []string{"out/*.xml"}, map[string]time.Time{"out/junit.xml": reps[0].Modified})
	if len(again) != 1 || again[0].Tests != 0 || !again[0].Modified.Equal(reps[0].Modified) {
		t.Errorf("known: its time alone: %+v", again)
	}
}

func TestDefaultRunners(t *testing.T) {
	t.Run("MNTRS-B12: The built-in runners recognize the common runners", func(t *testing.T) {})
	for cmd, want := range map[string]string{
		"node node_modules/.bin/jest --ci": "jest", "npx vitest run": "vitest", "go test ./...": "go test",
		"python -m pytest": "pytest", "maestro test flows/": "maestro", "npx playwright test": "playwright",
		"npx stryker run": "stryker", "gremlins unleash": "gremlins", "anchors test --all": "anchors",
	} {
		if r, ok := RunnerFor(DefaultRunners(), cmd); !ok || r.Name != want {
			t.Errorf("%q: %q %v, want %q", cmd, r.Name, ok, want)
		}
	}
	if _, ok := RunnerFor(DefaultRunners(), "vim notes.md"); ok {
		t.Error("an editor is no runner")
	}
}

func TestTick_everyEndSaidOnce(t *testing.T) {
	t.Run("MNTRS-I01: Every run that ended is said once", func(t *testing.T) {})
	dir := t.TempDir()
	exit0 := 0
	recs := []Run{
		{ID: "a", Command: "go test", PID: pidNpx, By: ByRun, Started: t0},
		{ID: "b", Command: "npx jest", PID: pidNpx + 1, By: ByRun, Started: t0, Output: "b"},
		{ID: "c", Command: "anchors test", PID: pidNpx + 2, By: ByAnchors, Started: t0},
		{ID: "d", Command: "node jest", PID: pidNpx + 3, By: ByOS, Started: t0},
	}
	procs := append(agentAndSelf(), Proc{PID: pidNpx, PPID: 1, Cmd: "go test"}, Proc{PID: pidNpx + 1, PPID: 1, Cmd: "npx jest"},
		Proc{PID: pidNpx + 2, PPID: 1, Cmd: "anchors test"}, Proc{PID: pidNpx + 3, PPID: 1, Cmd: "node jest"})
	out := map[string]OutputInfo{"b": {OK: true, Size: 1, Tail: "Tests: 1 passed\n"}}
	ends := map[string]int{}
	count := func(evs []Event) {
		for _, e := range evs {
			if e.Kind == "finished" || e.Kind == "died" || e.Kind == "ended" {
				ends[e.Run]++
			}
		}
	}
	st := LoadState(dir)
	evs, _ := tick(st, t0, procs, recs, out)
	count(evs)
	end := t0.Add(time.Minute)
	recs[0].Ended, recs[0].Exit = &end, &exit0
	evs, _ = tick(st, t0.Add(time.Minute), procs[:4], recs, out)
	count(evs)
	_ = SaveState(dir, st)
	st = LoadState(dir)
	for i := 0; i < 3; i++ {
		evs, _ = tick(st, t0.Add(time.Duration(2+i)*time.Minute), agentAndSelf(), recs, out)
		count(evs)
	}
	for _, id := range []string{"a", "b", "c", "d"} {
		if ends[id] != 1 {
			t.Errorf("run %s ended %d time(s): %v", id, ends[id], ends)
		}
	}
}

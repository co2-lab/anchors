// @anchors
//   code: MNCMT
//   ref: MNCMD

package flow

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/runs"
	"github.com/spf13/cobra"
)

// TestMain lets the test binary stand in for a command `anchors monitor run` wraps: with
// MONITOR_HELPER set, it prints a line and exits with that code.
func TestMain(m *testing.M) {
	if code := os.Getenv("MONITOR_HELPER"); code != "" {
		fmt.Println("helper says hi")
		var c int
		fmt.Sscanf(code, "%d", &c)
		os.Exit(c)
	}
	os.Exit(m.Run())
}

func monitorProject(t *testing.T, yaml string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte("version: 7\nlayers: {}\n"+yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func runCobra(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

func TestMonitor_onceAndUntilDone(t *testing.T) {
	t.Run("MNCMD-B01: The monitor reads once, or until every run ended", func(t *testing.T) {})
	dir := monitorProject(t, "")
	end, exit := time.Now(), 0
	if err := runs.Save(dir, runs.Run{ID: "r", Name: "jest", Command: "npx jest", By: runs.ByRun, Started: end.Add(-time.Minute), Ended: &end, Exit: &exit}); err != nil {
		t.Fatal(err)
	}
	out, err := runCobra(t, newMonitorCmd(), "--root", dir, "--once", "--heartbeat", "1h")
	if err != nil || !strings.Contains(out, "✓ jest ended (exit 0") {
		t.Fatalf("the run's end, once: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".anchors", "runs", "monitor.json")); err != nil {
		t.Error("the memory is written")
	}
	out, err = runCobra(t, newMonitorCmd(), "--root", dir, "--until-done", "--every", "10ms", "--heartbeat", "1h")
	if err != nil || !strings.Contains(out, "nothing running") || strings.Contains(out, "jest ended") {
		t.Errorf("nothing runs, it says so and exits, without repeating the end: %v\n%s", err, out)
	}
}

func TestMonitor_timing(t *testing.T) {
	t.Run("MNCMD-B02: The timing comes from the defaults, the block and the flags", func(t *testing.T) {})
	dir := monitorProject(t, "monitor:\n  heartbeat: 1h\n")
	if _, err := runCobra(t, newMonitorCmd(), "--root", dir, "--once"); err != nil {
		t.Errorf("a block that reads: %v", err)
	}
	if _, err := runCobra(t, newMonitorCmd(), "--root", dir, "--once", "--stall", "0s"); err == nil || !strings.Contains(err.Error(), `"0s"`) {
		t.Errorf("a flag that is no positive duration: %v", err)
	}
	bad := monitorProject(t, "monitor:\n  every: soon\n")
	if _, err := runCobra(t, newMonitorCmd(), "--root", bad, "--once"); err == nil || !strings.Contains(err.Error(), "monitor.every") {
		t.Errorf("a block that does not read: %v", err)
	}
}

func TestMonitor_badFlagNamed(t *testing.T) {
	t.Run("MNCMD-E01: A timing flag that does not read fails naming it", func(t *testing.T) {})
	dir := monitorProject(t, "")
	if _, err := runCobra(t, newMonitorCmd(), "--root", dir, "--once", "--every", "soon"); err == nil || !strings.Contains(err.Error(), `"soon"`) {
		t.Errorf("named: %v", err)
	}
}

func TestRun_recordsAndPassesTheExit(t *testing.T) {
	t.Run("MNCMD-B03: A wrapped command's run is recorded with its output and its exit", func(t *testing.T) {})
	dir := monitorProject(t, "")
	t.Setenv("MONITOR_HELPER", "3")
	out, err := runCobra(t, newRunCmd(), "--root", dir, "--", os.Args[0])
	var ec ExitCode
	if !errors.As(err, &ec) || ec.Code != 3 || !strings.Contains(out, "helper says hi") {
		t.Fatalf("its output passes, its code is the command's: %v\n%s", err, out)
	}
	recs := runs.List(dir)
	if len(recs) != 1 || recs[0].PID == 0 || recs[0].Exit == nil || *recs[0].Exit != 3 || recs[0].By != runs.ByRun {
		t.Fatalf("its record: %+v", recs)
	}
	if b, err := os.ReadFile(recs[0].Output); err != nil || !strings.Contains(string(b), "helper says hi") {
		t.Errorf("a copy of its output: %q %v", b, err)
	}
}

func hookCall(event, id, cwd, command string, background bool, answer string) []byte {
	resp := ""
	if answer != "" {
		resp = fmt.Sprintf(`,"tool_response":%s`, answer)
	}
	return []byte(fmt.Sprintf(`{"hook_event_name":%q,"tool_name":"Bash","tool_use_id":%q,"cwd":%q,"tool_input":{"command":%q,"run_in_background":%v}%s}`,
		event, id, cwd, command, background, resp))
}

func TestAgentHook_before(t *testing.T) {
	t.Run("MNCMD-B04: Before a command, the hook records the agent's long ones", func(t *testing.T) {})
	dir := monitorProject(t, "")
	sub := filepath.Join(dir, "apps")
	_ = os.MkdirAll(sub, 0o755)
	now := time.Now()
	agentHook(hookCall("PreToolUse", "toolu_1", sub, "./scripts/seed.sh", true, ""), now)
	agentHook(hookCall("PreToolUse", "toolu_2", dir, "ls -la", false, ""), now)
	agentHook(hookCall("PreToolUse", "toolu_3", dir, "anchors test --all", true, ""), now)
	agentHook(hookCall("PreToolUse", "toolu_4", t.TempDir(), "npx jest", true, ""), now)
	recs := runs.List(dir)
	if len(recs) != 1 || recs[0].ID != "agent-toolu_1" || recs[0].By != runs.ByAgent || recs[0].Command != "./scripts/seed.sh" {
		t.Errorf("only the background command, by its id: %+v", recs)
	}
}

func TestAgentHook_after(t *testing.T) {
	t.Run("MNCMD-B05: After a command, the hook records its output file or its exit", func(t *testing.T) {})
	dir := monitorProject(t, "")
	now := time.Now()
	agentHook(hookCall("PreToolUse", "toolu_bg", dir, "npx jest", true, ""), now)
	agentHook(hookCall("PreToolUse", "toolu_fg", dir, "npx jest --ci", false, ""), now)
	agentHook(hookCall("PostToolUse", "toolu_bg", dir, "npx jest", true, `{"stdout":"Command running in background with ID: b1. Output is being written to: /tmp/tasks/b1.output"}`), now)
	agentHook(hookCall("PostToolUse", "toolu_fg", dir, "npx jest --ci", false, `{"exit_code":1,"stdout":"x"}`), now)
	bg, _ := runs.Load(dir, "agent-toolu_bg")
	fg, _ := runs.Load(dir, "agent-toolu_fg")
	if bg.Output != "/tmp/tasks/b1.output" || bg.Done() {
		t.Errorf("the background command's output file: %+v", bg)
	}
	if !fg.Done() || fg.Exit == nil || *fg.Exit != 1 {
		t.Errorf("the foreground command ended with its exit: %+v", fg)
	}
}

func TestOwnRun(t *testing.T) {
	t.Run("MNCMD-B06: Anchors' long commands record their own run", func(t *testing.T) {})
	dir := monitorProject(t, "")
	bare := t.TempDir()
	tree := func() (*cobra.Command, map[string]*cobra.Command) {
		root := &cobra.Command{Use: "anchors"}
		subs := map[string]*cobra.Command{}
		for _, n := range []string{"test", "status"} {
			c := &cobra.Command{Use: n}
			c.Flags().String("root", ".", "")
			root.AddCommand(c)
			subs[n] = c
		}
		return root, subs
	}
	_, subs := tree()
	for _, c := range []struct {
		cmd  *cobra.Command
		root string
	}{{subs["test"], dir}, {subs["status"], dir}, {subs["test"], bare}} {
		_ = c.cmd.Flags().Set("root", c.root)
		BeginOwnRun(c.cmd)
		EndOwnRun(2)
	}
	recs := runs.List(dir)
	if len(recs) != 1 || recs[0].Name != "anchors test" || recs[0].Exit == nil || *recs[0].Exit != 2 || recs[0].By != runs.ByAnchors {
		t.Errorf("only anchors test in the project, ended with exit 2: %+v", recs)
	}
	if len(runs.List(bare)) != 0 {
		t.Error("a folder with no anchors.yaml records nothing")
	}
}

func TestAgentHook_neverBlocks(t *testing.T) {
	t.Run("MNCMD-I01: The agent's hook never blocks the agent", func(t *testing.T) {})
	dir := monitorProject(t, "")
	for _, in := range []string{"{not json", `{"tool_name":"Edit","tool_use_id":"x"}`,
		string(hookCall("PreToolUse", "toolu_x", t.TempDir(), "npx jest", true, "")),
		string(hookCall("PreToolUse", "toolu_y", dir, "npx jest", true, ""))} {
		c := newMonitorCmd()
		var out bytes.Buffer
		c.SetOut(&out)
		c.SetErr(&out)
		c.SetIn(strings.NewReader(in))
		c.SetArgs([]string{"hook"})
		if err := c.Execute(); err != nil || out.Len() != 0 {
			t.Errorf("%q: %v, printed %q", in, err, out.String())
		}
	}
}

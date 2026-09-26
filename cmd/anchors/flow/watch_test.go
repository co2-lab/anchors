package flow

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/daemon"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/fsnotify/fsnotify"
)

// watchCfg: `logic` keeps the whole triad; `model` waives the feature and the test.
func watchCfg() *config.Config {
	return &config.Config{
		Layers: map[string]config.Layer{
			"logic":   {Pattern: "src/**/*.ts", Kind: "code"},
			"model":   {Pattern: "models/**/*.ts", Kind: "code", OptionalTriadEdges: []string{"covered-by", "tested-by"}},
			"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
			"feature": {Pattern: "**/*.feature", Kind: "feature"},
			"test":    {Pattern: "src/**/*.test.ts", Kind: "test"},
		},
		Derived: &config.Derived{Anchor: "code", Files: map[string]config.Padroes{
			"spec":    {"{{dir}}/{{name}}.spec.md"},
			"feature": {"{{dir}}/{{name}}.feature"},
			"test":    {"{{dir}}/{{name}}.test.ts"},
		}},
	}
}

// useIgnore sets the watcher's ignore list for the test and restores it after.
func useIgnore(t *testing.T, root string, cfg *config.Config) {
	t.Helper()
	prev := watchIgnore
	watchIgnore = scan.LoadIgnoreFor(root, cfg)
	t.Cleanup(func() { watchIgnore = prev })
}

// queued returns the pending tasks as "changed→next".
func queued(t *testing.T, root string) []string {
	t.Helper()
	tasks, err := queue.List(root)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, tk := range tasks {
		out = append(out, tk.Changed+"→"+tk.SuggestedNext)
	}
	return out
}

// A change in a governed file queues the next piece of the chain — and skips the piece
// that already exists, so nobody is told to rewrite a finished feature.
func TestHandleChange_queuesTheNextMissingPiece(t *testing.T) {
	root := t.TempDir()
	cfg := watchCfg()
	useIgnore(t, root, cfg)
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	writeFile(t, root, "src/pricing.feature", "Feature: x\n")
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "src/pricing.ts", Rev: "old"}}}

	out := stdoutOf(t, func() { handleChange(root, cfg, g, "src/pricing.ts") })

	if got := queued(t, root); len(got) != 1 || got[0] != "src/pricing.ts→test" {
		t.Fatalf("the feature exists, so the next piece is the test; got %v\n%s", got, out)
	}
	if g.Nodes[0].Rev == "old" {
		t.Error("the node's revision must follow the file's content")
	}
	if !strings.Contains(out, "task queued") || !strings.Contains(out, "queue: 1 task(s)") {
		t.Errorf("the queueing is reported:\n%s", out)
	}

	// the same change again: deduplicated, and said
	out = stdoutOf(t, func() { handleChange(root, cfg, g, "src/pricing.ts") })
	if !strings.Contains(out, "already in the queue (test)") || len(queued(t, root)) != 1 {
		t.Errorf("a live duplicate must not be queued twice:\n%s", out)
	}
}

// A layer that waives the feature and the test goes straight to the review: queueing the
// waived piece contradicts `anchors work`, which answers STOP for it.
func TestHandleChange_skipsThePiecesTheLayerWaives(t *testing.T) {
	root := t.TempDir()
	cfg := watchCfg()
	useIgnore(t, root, cfg)
	writeFile(t, root, "models/user.ts", "export type U = {}\n")
	stdoutOf(t, func() { handleChange(root, cfg, &mapx.Graph{}, "models/user.ts") })
	if got := queued(t, root); len(got) != 1 || got[0] != "models/user.ts→review" {
		t.Errorf("feature and test are waived, so the review comes next; got %v", got)
	}

	// the spec of a waived unit resolves the unit's layer through the code file
	if l := pieceTargetLayer(root, "models/user.spec.md", cfg); l != "model" {
		t.Errorf("the spec's unit is in `model`, got %q", l)
	}
	if l := pieceTargetLayer(root, "models/ghost.spec.md", cfg); l != "" {
		t.Errorf("a spec with no code has no unit layer, got %q", l)
	}
	if l := pieceTargetLayer(root, "README.md", cfg); l != "" {
		t.Errorf("a file outside the structure has no layer, got %q", l)
	}
	if l := pieceTargetLayer(root, "src/x.ts", nil); l != "" {
		t.Errorf("no config, no layer, got %q", l)
	}
}

// What is not work never becomes a task: a file outside the structure, a file already
// gone, an editor's temporary, and what the project's .gitignore declares disposable.
func TestHandleChange_whatIsNotWorkIsNotQueued(t *testing.T) {
	root := t.TempDir()
	cfg := watchCfg()
	writeFile(t, root, ".gitignore", "probes/\n")
	writeFile(t, root, "README.md", "# x\n")
	writeFile(t, root, "src/probes/probe1.ts", "x\n")
	writeFile(t, root, "src/.!21662!pricing.spec.md", "x\n")
	useIgnore(t, root, cfg)
	g := &mapx.Graph{}
	stdoutOf(t, func() {
		for _, rel := range []string{"README.md", "src/gone.ts", "src/probes/probe1.ts", "src/.!21662!pricing.spec.md"} {
			handleChange(root, cfg, g, rel)
		}
	})
	if got := queued(t, root); len(got) != 0 {
		t.Errorf("nothing here is work; queued %v", got)
	}
}

// A delivery record triggers the review — but only once the unit has code AND test, and a
// plan's record triggers the review of the whole.
func TestHandleChange_deliveryRecordsTriggerTheReview(t *testing.T) {
	root := t.TempDir()
	cfg := watchCfg()
	useIgnore(t, root, cfg)
	g := &mapx.Graph{}

	writeFile(t, root, "changes/a.md", "stage: spec\nunit: src/pricing.ts\n")
	out := stdoutOf(t, func() { handleChange(root, cfg, g, "changes/a.md") })
	if len(queued(t, root)) != 0 || !strings.Contains(out, "the review waits for the triad to close") {
		t.Fatalf("with no code yet the review must wait:\n%s", out)
	}

	writeFile(t, root, "src/pricing.ts", "x\n")
	stdoutOf(t, func() { handleChange(root, cfg, g, "changes/a.md") })
	if len(queued(t, root)) != 0 {
		t.Fatal("code without test: the triad is still open")
	}

	writeFile(t, root, "src/pricing.test.ts", "x\n")
	stdoutOf(t, func() { handleChange(root, cfg, g, "changes/a.md") })
	if got := queued(t, root); len(got) != 1 || got[0] != "changes/a.md→review" {
		t.Fatalf("code and test exist: the review is queued; got %v", got)
	}

	writeFile(t, root, "changes/plan.md", "stage: plan\nunit: plans/0001.md\n")
	stdoutOf(t, func() { handleChange(root, cfg, g, "changes/plan.md") })
	if got := queued(t, root); len(got) != 2 || !containsStr(got, "changes/plan.md→review-plan") {
		t.Errorf("a plan's delivery asks for the review of the whole; got %v", got)
	}

	// a record already reviewed does not reopen the queue
	writeFile(t, root, "changes/reviewed/old.md", "stage: code\nunit: src/pricing.ts\n")
	stdoutOf(t, func() { handleChange(root, cfg, g, "changes/reviewed/old.md") })
	if got := queued(t, root); len(got) != 2 {
		t.Errorf("a reviewed record is history, not work; got %v", got)
	}
}

// A model whose test is waived is reviewed as soon as its code exists — otherwise it
// would never be.
func TestMissingPieceToReview(t *testing.T) {
	root := t.TempDir()
	cfg := watchCfg()
	writeFile(t, root, "models/user.ts", "x\n")
	if missingPieceToReview(root, "models/user.ts", cfg) {
		t.Error("the test is waived in `model`: the code alone closes the triad")
	}
	if missingPieceToReview(root, "", cfg) {
		t.Error("a record with no readable unit does not hold the review")
	}
	writeFile(t, root, "svc/handler.go", "x\n")
	writeFile(t, root, "svc/handler_test.go", "x\n")
	if missingPieceToReview(root, "svc/handler.go", nil) {
		t.Error("a Go test next to the code closes the triad")
	}
}

func TestChangeDelivered(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "changes/a.md", "stage: code\nunit: src/x.ts\n")
	writeFile(t, root, "changes/b.md", "no unit line\n")
	if u, ok := changeDelivered(root, "changes/a.md"); !ok || u != "src/x.ts" {
		t.Errorf("got %q %v", u, ok)
	}
	if u, ok := changeDelivered(root, "changes/b.md"); !ok || u != "" {
		t.Errorf("a record without a unit is still a delivery: got %q %v", u, ok)
	}
	for _, rel := range []string{"changes/missing.md", "changes/a.txt", "docs/a.md"} {
		if _, ok := changeDelivered(root, rel); ok {
			t.Errorf("%s is not a pending delivery record", rel)
		}
	}
}

func TestTaskID_isStableAndPathSafe(t *testing.T) {
	a, b := taskID("src/x/pricing.ts", "test"), taskID("src/x/pricing.ts", "test")
	if a != b {
		t.Errorf("the id must be deterministic: %q vs %q", a, b)
	}
	if !strings.HasPrefix(a, "src-x-pricing-test-") || strings.Contains(a, "/") {
		t.Errorf("the id is a slug of the path and the step: %q", a)
	}
	if taskID("src/x/pricing.ts", "review") == a {
		t.Error("another step is another task")
	}
}

// The watcher adds the whole tree except what it ignores, and a new folder's files are
// swept once so the ones born before the watch took effect are not lost.
func TestAddTreeToWatcherAndFilesBornIn(t *testing.T) {
	root := t.TempDir()
	useIgnore(t, root, watchCfg())
	writeFile(t, root, "src/a/x.ts", "x\n")
	writeFile(t, root, "node_modules/lib/y.js", "y\n")
	writeFile(t, root, "src/b.ts", "b\n")

	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	// root, src, src/a — node_modules and everything under it stay out
	if n := addTreeToWatcher(w, root); n != 3 {
		t.Errorf("expected 3 directories watched, got %d: %v", n, w.WatchList())
	}
	for _, d := range w.WatchList() {
		if strings.Contains(d, "node_modules") {
			t.Errorf("an ignored folder is not watched: %s", d)
		}
	}

	if got := filesBornIn(filepath.Join(root, "src")); len(got) != 1 || filepath.Base(got[0]) != "b.ts" {
		t.Errorf("only the files of the first level: %v", got)
	}
	if got := filesBornIn(filepath.Join(root, "nope")); got != nil {
		t.Errorf("a missing folder has no files: %v", got)
	}
}

// The whole loop, in process: a file created after the start becomes a task, including
// one in a folder born after the start; the signal ends the loop and cleans the state.
func TestRunWatchLoop_queuesChangesUntilSignalled(t *testing.T) {
	root := t.TempDir()
	cfg := watchCfg()
	writeFile(t, root, "src/keep.ts", "x\n")
	prev := watchIgnore
	t.Cleanup(func() { watchIgnore = prev })
	p := daemon.PathsFor(root)
	if err := daemon.WritePID(p, os.Getpid()); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	out := stdoutOf(t, func() {
		go func() { done <- runWatchLoop(root, cfg, &mapx.Graph{}, 10*time.Millisecond) }()
		waitFor(t, func() bool { return watchIgnore != nil }) // the loop is set up
		time.Sleep(100 * time.Millisecond)                    // fsnotify registration
		writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
		writeFile(t, root, "src/fresh/handler.ts", "export const h = 1\n")
		waitFor(t, func() bool { return len(queued(t, root)) >= 2 })

		// the loop owns SIGTERM now; the signal ends it
		if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
			t.Fatal(err)
		}
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("a signal is a clean exit, got %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("the loop did not stop on SIGTERM")
		}
	})
	got := queued(t, root)
	if !containsStr(got, "src/pricing.ts→feature") || !containsStr(got, "src/fresh/handler.ts→feature") {
		t.Errorf("both new files must be queued, got %v\n%s", got, out)
	}
	if !strings.Contains(out, "watching") || !strings.Contains(out, "shutting down") {
		t.Errorf("the loop reports its start and its end:\n%s", out)
	}
	if _, err := os.Stat(p.PID); err == nil {
		t.Error("the loop cleans its pid file on exit")
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not reached in 10s")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func containsStr(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// sleeper is a live process standing in for the daemon.
func sleeper(t *testing.T) *exec.Cmd {
	t.Helper()
	c := exec.Command("sleep", "30")
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Process.Kill(); _, _ = c.Process.Wait() })
	return c
}

func runWatchSub(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newWatchCmd()
	cmd.SetArgs(args)
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	return out, err
}

// The controls read and write the daemon's state in the project: stopped, running,
// paused, resumed, stopped for real, and its log.
func TestWatchControls(t *testing.T) {
	root := t.TempDir()
	p := daemon.PathsFor(root)

	if out, _ := runWatchSub(t, "status", "--root", root); !strings.Contains(out, "watcher: stopped") {
		t.Errorf("no pid, stopped:\n%s", out)
	}
	if _, err := runWatchSub(t, "pause", "--root", root); err == nil || !strings.Contains(err.Error(), "not running") {
		t.Errorf("pausing nothing is an error, got %v", err)
	}
	if _, err := runWatchSub(t, "logs", "--root", root); err == nil || !strings.Contains(err.Error(), "no log") {
		t.Errorf("no log yet, got %v", err)
	}

	s := sleeper(t)
	if err := daemon.WritePID(p, s.Process.Pid); err != nil {
		t.Fatal(err)
	}
	_ = daemon.WriteMeta(p, time.Now(), root)
	if out, _ := runWatchSub(t, "status", "--root", root); !strings.Contains(out, "watcher: running (pid "+strconv.Itoa(s.Process.Pid)+")") ||
		!strings.Contains(out, "root="+root) {
		t.Errorf("a live pid is running, with its meta:\n%s", out)
	}
	if _, err := runWatchSub(t, "start", "--root", root); err == nil || !strings.Contains(err.Error(), "already running") {
		t.Errorf("a second start must be refused, got %v", err)
	}
	if out, err := runWatchSub(t, "pause", "--root", root); err != nil || !daemon.IsPaused(p) || !strings.Contains(out, "paused") {
		t.Errorf("pause must leave the flag: %v\n%s", err, out)
	}
	if out, _ := runWatchSub(t, "status", "--root", root); !strings.Contains(out, "watcher: paused") {
		t.Errorf("status tells paused:\n%s", out)
	}
	if _, err := runWatchSub(t, "resume", "--root", root); err != nil || daemon.IsPaused(p) {
		t.Errorf("resume must remove the flag: %v", err)
	}
	writeFile(t, root, ".anchors/watch.log", "● src/a.ts [code] → task queued\n")
	if out, _ := runWatchSub(t, "logs", "--root", root); out != "● src/a.ts [code] → task queued\n" {
		t.Errorf("logs prints the log as is: %q", out)
	}

	if out, err := runWatchSub(t, "stop", "--root", root); err != nil || !strings.Contains(out, "watcher shut down") {
		t.Fatalf("stop: %v\n%s", err, out)
	}
	if err := s.Wait(); err == nil {
		t.Error("the daemon must have been terminated by the signal")
	}
	if _, err := os.Stat(p.PID); err == nil {
		t.Error("stop removes the pid file")
	}
	if _, err := runWatchSub(t, "stop", "--root", root); err == nil {
		t.Error("stopping a stopped watcher is an error")
	}
}

// `watch run` loads the config and the map before looping, and says which one is missing.
func TestWatchRun_needsConfigAndMap(t *testing.T) {
	root := t.TempDir()
	if _, err := runWatchSub(t, "run", "--root", root); err == nil || !strings.Contains(err.Error(), "load config") {
		t.Errorf("without anchors.yaml, got %v", err)
	}
	writeFile(t, root, "anchors.yaml", "version: 1\nlayers:\n  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n")
	writeFile(t, root, "anchors.graph.yaml", "version: 4\nnodes: [\n")
	if _, err := runWatchSub(t, "run", "--root", root); err == nil || !strings.Contains(err.Error(), "load map") {
		t.Errorf("with a broken map, got %v", err)
	}
}

// detach puts the child in a session of its own: it leads its process group, so closing
// the terminal does not take it down.
func TestDetach_childLeadsItsOwnSession(t *testing.T) {
	c := exec.Command("sleep", "30")
	detach(c)
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Process.Kill(); _, _ = c.Process.Wait() }()
	pgid, err := syscall.Getpgid(c.Process.Pid)
	if err != nil {
		t.Fatal(err)
	}
	if pgid != c.Process.Pid {
		t.Errorf("a detached child leads its own group: pgid %d, pid %d", pgid, c.Process.Pid)
	}
	if mine, _ := syscall.Getpgid(os.Getpid()); pgid == mine {
		t.Error("the child must not stay in the parent's process group")
	}
}

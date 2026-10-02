// @anchors
//   ref: BDGRN

package quality

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

func timedNode(id string, kind mapx.Kind, bySuite map[string]float64) mapx.Node {
	n := mapx.Node{ID: id, Kind: kind, Rev: "r-" + id}
	if bySuite != nil {
		n.Signal = &mapx.TestSignal{SecondsBySuite: bySuite}
	}
	return n
}

func TestBudgetPlan(t *testing.T) {
	t.Run("BDGRN-B01: The plan is the timed files fastest first, then the untimed", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		timedNode("slow_test.go", mapx.KindTest, map[string]float64{"out/junit.xml": 3}),
		timedNode("b_test.go", mapx.KindTest, map[string]float64{"out/junit.xml": 1}),
		timedNode("a_test.go", mapx.KindTest, map[string]float64{"out/junit.xml": 1}),
		timedNode("other_test.go", mapx.KindTest, map[string]float64{"web/junit.xml": 0.1}),
		timedNode("z_new_test.go", mapx.KindTest, nil),
		timedNode("m_new_test.go", mapx.KindTest, nil),
		{ID: "helpers_test.go", Kind: mapx.KindTest, Support: true},
		timedNode("a.go", mapx.KindCode, nil),
	}}
	timed, untimed := budgetPlan(g, mapx.KindTest, "out/junit.xml", config.Suite{})
	want := []timedFile{{"a_test.go", 1}, {"b_test.go", 1}, {"slow_test.go", 3}}
	if !reflect.DeepEqual(timed, want) || !reflect.DeepEqual(untimed, []string{"m_new_test.go", "z_new_test.go"}) {
		t.Fatalf("want %v then [m_new z_new], got %v then %v", want, timed, untimed)
	}
}

func TestBudgetBatches(t *testing.T) {
	t.Run("BDGRN-B02: Batches take what fits and then the untimed one by one", func(t *testing.T) {})
	timed := []timedFile{{"a", 1}, {"b", 2}, {"c", 4}}
	untimed := []string{"x", "y"}
	batch, restT, restU := nextBatch(timed, untimed, 3500*time.Millisecond, false)
	if !reflect.DeepEqual(batch, []string{"a", "b"}) || len(restT) != 1 || len(restU) != 2 {
		t.Fatalf("3.5s takes a and b, got %v (%v, %v)", batch, restT, restU)
	}
	if batch, _, _ := nextBatch(timed, nil, 3*time.Second, false); !reflect.DeepEqual(batch, []string{"a", "b"}) {
		t.Fatalf("files whose times add up to exactly what remains fit, got %v", batch)
	}
	if batch, _, restU := nextBatch(restT, restU, 3*time.Second, false); len(batch) != 0 || len(restU) != 2 {
		t.Fatalf("when the fastest left does not fit, nothing runs and the untimed wait, got %v", batch)
	}
	batch, restT, restU = nextBatch(nil, untimed, time.Second, false)
	if !reflect.DeepEqual(batch, []string{"x"}) || restT != nil || !reflect.DeepEqual(restU, []string{"y"}) {
		t.Fatalf("after the timed, one untimed per batch, got %v (%v)", batch, restU)
	}
	if batch, _, _ := nextBatch(nil, untimed, 0, false); batch != nil {
		t.Fatalf("with no time left nothing runs, got %v", batch)
	}
	batch, restT, _ = nextBatch(timed, nil, time.Hour, true)
	if !reflect.DeepEqual(batch, []string{"a"}) || len(restT) != 2 {
		t.Fatalf("one at a time takes a single file, got %v", batch)
	}
}

func TestBudgetDeadlineStopsTheGroup(t *testing.T) {
	t.Run("BDGRN-B03: A batch still running at the deadline is stopped with its group", func(t *testing.T) {})
	dir := t.TempDir()
	// The child outlives the shell: only a group kill stops it before it writes `late`.
	start := time.Now()
	_, cut := execUntil("(sleep 2; touch late) & sleep 30", dir, time.Now().Add(300*time.Millisecond))
	if !cut {
		t.Fatal("the run must be reported as cut by the deadline")
	}
	if took := time.Since(start); took > 8*time.Second {
		t.Fatalf("the deadline must stop the run promptly, it took %v", took)
	}
	time.Sleep(2500 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(dir, "late")); err == nil {
		t.Fatal("the child must be stopped with the group, but it wrote its file")
	}
	if err, cut := execUntil("true", dir, time.Now().Add(10*time.Second)); err != nil || cut {
		t.Fatalf("a run that ends in time is not cut, got %v %v", err, cut)
	}

	// The group gets a TERM first: a tool that restores the source on it gets to. Windows
	// has no TERM a command can trap: there the stop is immediate, and this is not asked.
	if runtime.GOOS == "windows" {
		return
	}
	prev := stopGrace
	stopGrace = time.Second
	t.Cleanup(func() { stopGrace = prev })
	if _, cut := execUntil("trap 'touch restored; exit 0' TERM; sleep 30", dir, time.Now().Add(300*time.Millisecond)); !cut {
		t.Fatal("the run must be reported as cut")
	}
	if _, err := os.Stat(filepath.Join(dir, "restored")); err != nil {
		t.Fatal("the TERM must reach the command, so its restore runs")
	}
	// A child that ignores the TERM is killed once the grace is over, before it finishes.
	start = time.Now()
	execUntil("(trap '' TERM; sleep 4; touch stubborn) & wait", dir, time.Now().Add(300*time.Millisecond))
	time.Sleep(4500*time.Millisecond - time.Since(start))
	if _, err := os.Stat(filepath.Join(dir, "stubborn")); err == nil {
		t.Fatal("a child ignoring the TERM must be killed after the grace")
	}
}

// budgetRunner is the fake runner: it logs each batch's files and writes a JUnit report
// for them, each case taking `time`.
const budgetRunner = `#!/bin/sh
echo "batch: $*" >> order.log
{
  echo '<testsuites><testsuite name="s">'
  for f in "$@"; do echo "<testcase name=\"case\" file=\"$f\" time=\"0.4\"/>"; done
  echo '</testsuite></testsuites>'
} > out/junit.xml
`

func budgetProject(t *testing.T, section, runChanged string, g *mapx.Graph) string {
	t.Helper()
	report := `    junit: "out/junit.xml"`
	if section == "mutation" {
		report = `    report: "out/mutation.json"`
	}
	yaml := suiteLayers + `  tests:
    kind: test
    pattern: "*_test.go"
` + section + `:
  - layer: unit
    run: "echo full"
    run_changed: "` + runChanged + `"
` + report + "\n"
	files := map[string]string{"run.sh": budgetRunner, "mut.sh": budgetMutator, "out/.keep": ""}
	// The run reads each file's revision from its content: a node measured at its own rev
	// keeps being so under the file's real one.
	rev := scan.ShortHash([]byte("package a\n"))
	for i := range g.Nodes {
		n := &g.Nodes[i]
		files[n.ID] = "package a\n"
		if s := n.Signal; s != nil {
			if s.AtRev == n.Rev {
				s.AtRev = rev
			}
			if s.MutationAtRev == n.Rev {
				s.MutationAtRev = rev
			}
		}
		n.Rev = rev
	}
	return qProject(t, yaml, files, g)
}

func readLog(t *testing.T, dir string) []string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(dir, "order.log"))
	var out []string
	for _, l := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		if l == "" {
			continue
		}
		var names []string
		for _, f := range strings.Fields(strings.TrimPrefix(l, "batch: ")) {
			names = append(names, filepath.Base(f))
		}
		out = append(out, strings.Join(names, " "))
	}
	return out
}

func TestBudgetRunsFastestFirst(t *testing.T) {
	t.Run("BDGRN-B04: A budget runs batches fastest first and reports what ran and what was left", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		timedNode("slow_test.go", mapx.KindTest, map[string]float64{"out/junit.xml": 2}),
		timedNode("fast_test.go", mapx.KindTest, map[string]float64{"out/junit.xml": 0.5}),
		timedNode("new_test.go", mapx.KindTest, nil),
	}}
	dir := budgetProject(t, "tests", "sh run.sh {{files}}", g)
	out, err := runQ(t, newTestCmd(), "--root", dir, "--budget", "60s")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := readLog(t, dir); !reflect.DeepEqual(got, []string{"fast_test.go slow_test.go", "new_test.go"}) {
		t.Fatalf("the timed run fastest first in one batch, then the untimed, got %v", got)
	}
	if !strings.Contains(out, "budget: ran 3 file(s); 0 left") {
		t.Errorf("the report says what ran and what was left:\n%s", out)
	}
	m, err := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range m.Nodes {
		if n.ID == "new_test.go" && (n.Signal == nil || n.Signal.SecondsBySuite["out/junit.xml"] != 0.4) {
			t.Fatalf("each batch is ingested, and the untimed file now has its time, got %+v", n.Signal)
		}
	}
	// A budget the slow file does not fit in leaves it, and the untimed, for later.
	dir = budgetProject(t, "tests", "sh run.sh {{files}}", g)
	out, err = runQ(t, newTestCmd(), "--root", dir, "--budget", "1s")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := readLog(t, dir); !reflect.DeepEqual(got, []string{"fast_test.go"}) || !strings.Contains(out, "ran 1 file(s); 2 left") {
		t.Fatalf("only the fast file fits in 1s, got %v\n%s", got, out)
	}
}

// budgetMutator is the fake mutation runner: it logs each batch and writes a mutation
// report for the file it was given, one killed mutant.
const budgetMutator = `#!/bin/sh
echo "batch: $*" >> order.log
sleep 0.2
printf '{"schemaVersion":"1","thresholds":{"high":80,"low":60},"files":{"%s":{"language":"go","source":"","mutants":[{"id":"1","mutatorName":"m","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}},"status":"Killed"}]}}}' "$1" > out/mutation.json
`

func TestBudgetMutationOneFileAtATime(t *testing.T) {
	t.Run("BDGRN-B05: A mutation budget runs one file per batch and records its time", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{timedNode("a.go", mapx.KindCode, nil), timedNode("b.go", mapx.KindCode, nil)}}
	dir := budgetProject(t, "mutation", "sh mut.sh {{files}}", g)
	out, err := runQ(t, newMutationCmd(), "--root", dir, "--budget", "60s")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := readLog(t, dir); !reflect.DeepEqual(got, []string{"a.go", "b.go"}) {
		t.Fatalf("each code file runs alone, got %v", got)
	}
	m, err := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range m.Nodes {
		if n.Signal == nil || n.Signal.SecondsBySuite["out/mutation.json"] < 0.2 || n.Signal.MutantsKilled != 1 {
			t.Fatalf("%s must be ingested and timed by Anchors (at least the 0.2s it slept), got %+v", n.ID, n.Signal)
		}
	}
}

func TestBudgetFailedBatchDoesNotStop(t *testing.T) {
	t.Run("BDGRN-B06: A failed batch does not stop the budget", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{timedNode("a_test.go", mapx.KindTest, nil), timedNode("b_test.go", mapx.KindTest, nil)}}
	dir := budgetProject(t, "tests", "sh run.sh {{files}}; case {{files}} in *a_test.go) exit 1;; esac", g)
	out, err := runQ(t, newTestCmd(), "--root", dir, "--budget", "60s")
	if err == nil || !strings.Contains(err.Error(), "a batch failed in: unit") {
		t.Fatalf("the command fails at the end naming the suite, got %v\n%s", err, out)
	}
	if got := readLog(t, dir); !reflect.DeepEqual(got, []string{"a_test.go", "b_test.go"}) {
		t.Fatalf("the next batch still runs after a failed one, got %v", got)
	}
}

func TestBudgetRefusesWhatItCannotRun(t *testing.T) {
	t.Run("BDGRN-B07: What the budget cannot run is refused before anything runs", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{timedNode("a_test.go", mapx.KindTest, nil)}}
	yaml := suiteLayers + `tests:
  - layer: unit
    run: "echo full"
    junit: "out/junit.xml"
  - layer: e2e
    run: "echo full"
    run_changed: "sh run.sh {{files}}"
`
	dir := qProject(t, yaml, map[string]string{"run.sh": budgetRunner, "a_test.go": "x"}, g)
	for _, c := range []struct {
		args []string
		why  string
	}{
		{[]string{"unit"}, "declares none"},
		{[]string{"e2e"}, "needs the suite's report"},
		{[]string{"unit", "e2e"}, "declares none"},
		{[]string{"unit", "--changed", "a_test.go"}, "use one"},
	} {
		args := append([]string{"--root", dir, "--budget", "60s"}, c.args...)
		if _, err := runQ(t, newTestCmd(), args...); err == nil || !strings.Contains(err.Error(), c.why) {
			t.Errorf("%v: want a refusal naming %q, got %v", c.args, c.why, err)
		}
	}
	if got := readLog(t, dir); len(got) != 0 {
		t.Fatalf("no command may run before a refusal, got %v", got)
	}
}

func TestBudgetWithoutAMap(t *testing.T) {
	t.Run("BDGRN-E02: A budget without a map is refused saying to build it", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{timedNode("a_test.go", mapx.KindTest, nil)}}
	dir := budgetProject(t, "tests", "sh run.sh {{files}}", g)
	if err := os.Remove(filepath.Join(dir, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	if _, err := runQ(t, newTestCmd(), "--root", dir, "--budget", "60s"); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Fatalf("without a map the budget is refused saying to build it, got %v", err)
	}
	if got := readLog(t, dir); len(got) != 0 {
		t.Fatalf("no command may run, got %v", got)
	}
}

func TestBudgetPlan_SuitePaths(t *testing.T) {
	t.Run("BDGRN-B08: The plan holds only the suite's own files", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{timedNode("apps/mobile/a.test.ts", mapx.KindTest, nil), timedNode("apps/landing/b.test.ts", mapx.KindTest, nil)}}
	timed, untimed := budgetPlan(g, mapx.KindTest, "k", config.Suite{Paths: []string{"apps/mobile/**"}})
	if len(timed) != 0 || !reflect.DeepEqual(untimed, []string{"apps/mobile/a.test.ts"}) {
		t.Fatalf("only the mobile file, got %v %v", timed, untimed)
	}
}

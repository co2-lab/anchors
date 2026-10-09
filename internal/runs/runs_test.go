// @anchors
//   code: PRCRT
//   ref: PRCRN

package runs

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var t0 = time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

func TestRecords_writtenWholeAndListed(t *testing.T) {
	t.Run("PRCRN-B01: A record is written whole and read back as written", func(t *testing.T) {})
	root := t.TempDir()
	a := Run{ID: "b", Name: "jest", Command: "npx jest", By: ByAgent, Started: t0.Add(time.Minute)}
	b := Run{ID: "a", Name: "go test", Command: "go test ./...", PID: 42, By: ByRun, Started: t0}
	for _, r := range []Run{a, b} {
		if err := Save(root, r); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, ".anchors", "runs", "broken.json"), []byte("{half"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SaveState(root, &State{}); err != nil {
		t.Fatal(err)
	}
	if id := NewID(t0, 42); id != "20261009T120000-42" {
		t.Errorf("the id: when it started and its process, got %q", id)
	}
	got := List(root)
	if len(got) != 2 || !reflect.DeepEqual(got[0], b) || !reflect.DeepEqual(got[1], a) {
		t.Errorf("both records, oldest first, as written: %+v", got)
	}
}

func TestFinish(t *testing.T) {
	t.Run("PRCRN-B02: Finishing a run records how it ended", func(t *testing.T) {})
	root := t.TempDir()
	for _, id := range []string{"x", "y"} {
		if err := Save(root, Run{ID: id, Started: t0, State: StateRunning}); err != nil {
			t.Fatal(err)
		}
	}
	code := 1
	if err := Finish(root, "x", t0.Add(time.Minute), &code, "2 failed"); err != nil {
		t.Fatal(err)
	}
	if err := Finish(root, "y", t0.Add(time.Minute), nil, ""); err != nil {
		t.Fatal(err)
	}
	x, _ := Load(root, "x")
	y, _ := Load(root, "y")
	if x.State != StateFinished || x.Exit == nil || *x.Exit != 1 || x.Summary != "2 failed" || !x.Done() {
		t.Errorf("finished with its exit: %+v", x)
	}
	if y.State != StateEnded || y.Exit != nil || !y.Done() {
		t.Errorf("ended with no exit: %+v", y)
	}
}

func TestFinish_noRecord(t *testing.T) {
	t.Run("PRCRN-E01: Finishing a run that has no record fails", func(t *testing.T) {})
	root := t.TempDir()
	if err := Finish(root, "ghost", t0, nil, ""); err == nil {
		t.Error("finishing nothing is an error")
	}
	if len(List(root)) != 0 {
		t.Error("nothing is written")
	}
}

func TestPrune(t *testing.T) {
	t.Run("PRCRN-B03: Pruning keeps the runs going and the latest ended ones", func(t *testing.T) {})
	root := t.TempDir()
	at := func(d time.Duration) *time.Time { x := t0.Add(d); return &x }
	_ = Save(root, Run{ID: "live", Started: t0})
	_ = Save(root, Run{ID: "e1", Started: t0, Ended: at(time.Minute)})
	_ = Save(root, Run{ID: "e2", Started: t0, Ended: at(2 * time.Minute)})
	_ = Save(root, Run{ID: "e3", Started: t0, Ended: at(3 * time.Minute)})
	_ = Save(root, Run{ID: "old", Started: t0, Ended: at(-30 * 24 * time.Hour)})
	if n := Prune(root, 2, 7*24*time.Hour, t0.Add(time.Hour)); n != 2 {
		t.Errorf("two removed, got %d", n)
	}
	var ids []string
	for _, r := range List(root) {
		ids = append(ids, r.ID)
	}
	if !reflect.DeepEqual(ids, []string{"e2", "e3", "live"}) && !reflect.DeepEqual(ids, []string{"live", "e2", "e3"}) {
		t.Errorf("the running one and the two latest ended: %v", ids)
	}
}

func TestProcessTable_read(t *testing.T) {
	t.Run("PRCRN-B04: The process table is read as the system prints it", func(t *testing.T) {})
	for in, want := range map[string]float64{"0:01.50": 1.5, "1:02:03": 3723, "2-00:00:01": 172801, "00:00:07": 7, "x": 0} {
		if got := parseCPUTime(in); got != want {
			t.Errorf("%q: %v, want %v", in, got, want)
		}
	}
	procs := parsePS("  12     1   0:01.00 /usr/bin/node  jest  --ci\n  13    12   0:00.50 sh -c x\nbad line\n")
	if len(procs) != 2 || procs[0].PID != 12 || procs[0].PPID != 1 || procs[0].CPU != 1 || procs[0].Cmd != "/usr/bin/node  jest  --ci" || procs[1].Cmd != "sh -c x" {
		t.Errorf("each line a process, the command with its own spacing: %+v", procs)
	}
	if got := parseLsofCwd("p12\nfcwd\nn/proj/app\np13\nn/tmp\n"); got[12] != "/proj/app" || got[13] != "/tmp" {
		t.Errorf("each process its folder: %v", got)
	}
}

func TestInProject(t *testing.T) {
	t.Run("PRCRN-B05: A process belongs to the project by its folder or its command line", func(t *testing.T) {})
	root := filepath.Join(string(filepath.Separator), "work", "proj")
	cases := []struct {
		p    Proc
		want bool
	}{
		{Proc{Cwd: root}, true},
		{Proc{Cwd: filepath.Join(root, "apps", "mobile")}, true},
		{Proc{Cwd: root + "-other"}, false},
		{Proc{Cmd: "node " + filepath.Join(root, "node_modules", "jest")}, true},
		{Proc{Cmd: "node /elsewhere/jest"}, false},
	}
	for _, c := range cases {
		if got := InProject(c.p, root); got != c.want {
			t.Errorf("%+v: %v, want %v", c.p, got, c.want)
		}
	}
}

func TestTree(t *testing.T) {
	t.Run("PRCRN-B06: The process tree gives descendants, ancestors and the tree's CPU time", func(t *testing.T) {})
	tr := NewTree([]Proc{{PID: 1, PPID: 0, CPU: 1}, {PID: 2, PPID: 1, CPU: 2}, {PID: 3, PPID: 2, CPU: 4}, {PID: 9, PPID: 0}})
	if d := tr.Descendants(1); !reflect.DeepEqual(d, []int{2, 3}) {
		t.Errorf("descendants at any depth: %v", d)
	}
	if a := tr.Ancestors(3); !reflect.DeepEqual(a, []int{2, 1}) {
		t.Errorf("ancestors nearest first: %v", a)
	}
	if c := tr.TreeCPU(1); c != 7 {
		t.Errorf("the tree's CPU: %v", c)
	}
}

func TestListProcs_includesItself(t *testing.T) {
	t.Run("PRCRN-B07: Listing the system's processes includes the process asking", func(t *testing.T) {})
	procs, err := ListProcs()
	if err != nil {
		t.Skipf("no process table here: %v", err)
	}
	for _, p := range procs {
		if p.PID == os.Getpid() {
			return
		}
	}
	t.Errorf("its own process %d is not among %d", os.Getpid(), len(procs))
}

func TestSave_whole(t *testing.T) {
	t.Run("PRCRN-I01: A reader never sees half a record", func(t *testing.T) {})
	root := t.TempDir()
	_ = Save(root, Run{ID: "r", Command: "first"})
	if err := Save(root, Run{ID: "r", Command: "second"}); err != nil {
		t.Fatal(err)
	}
	if r, err := Load(root, "r"); err != nil || r.Command != "second" {
		t.Errorf("read back whole: %+v %v", r, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".anchors", "runs", "r.json.tmp")); err == nil {
		t.Error("no temporary file is left")
	}
}

func TestSystemReadings(t *testing.T) {
	t.Run("PRCRN-B08: The system gives a process's folder, whether it lives, and the load", func(t *testing.T) {})
	cwd, _ := os.Getwd()
	got := Cwds([]int{os.Getpid()})[os.Getpid()]
	if runtime.GOOS == "windows" {
		if got != "" {
			t.Errorf("Windows gives no folder: %q", got)
		}
	} else if real, _ := filepath.EvalSymlinks(cwd); got != cwd && got != real {
		t.Errorf("its folder: %q, want %q", got, cwd)
	}
	if len(Cwds(nil)) != 0 {
		t.Error("no process asked, no folder")
	}
	if !Alive(os.Getpid()) || Alive(0) {
		t.Error("it lives, process 0 does not")
	}
	if _, ok := firstFloat(""); ok {
		t.Error("an empty reading is no load")
	}
	if v, ok := firstFloat("2.50 1.0 0.5"); !ok || v != 2.5 {
		t.Errorf("the first figure: %v %v", v, ok)
	}
	load, ok := LoadAverage()
	if (runtime.GOOS == "windows") == ok || (ok && load < 0) {
		t.Errorf("the load where the system has one: %v %v", load, ok)
	}
}

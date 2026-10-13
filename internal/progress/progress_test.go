// @anchors
//   code: PFTPR
//   ref: DMLDD

package progress

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func sample(t *testing.T) *File {
	t.Helper()
	f := &File{Unit: "hours"}
	d, _ := f.AddDemand("Checkout", "plans/0012-checkout.md")
	must := func(_ string, err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(f.AddTask(d, "kick-off", 1, false))
	s, _ := f.AddStage(d, "Payment")
	must(f.AddTask(s, "cart spec", 3, false))
	must(f.AddTask(s, "cart code", 2, false))
	must(f.AddTask(s, "payment provider", 8, false))
	must(f.AddTask(s, "refund rules", 5, false))
	must(f.AddTask(s, "legacy docs", 1, true))
	if err := f.Done("D01.T2", "agent-a", "2026-10-12"); err != nil {
		t.Fatal(err)
	}
	if err := f.Done("D01.T3", "agent-b", "2026-10-12"); err != nil {
		t.Fatal(err)
	}
	if err := f.Start("D01.T4", "agent-c", "2026-10-12", false); err != nil {
		t.Fatal(err)
	}
	if err := f.Cancel("D01.T5", "the provider covers it"); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestRender_readsBack(t *testing.T) {
	t.Run("DMLDD-B01: The file is rendered and read back", func(t *testing.T) {})
	out := Render(sample(t))
	for _, want := range []string{
		headerLine + "\n# Progress — 33% (5/15 hours)\n",
		"## D01 — Checkout (plan: plans/0012-checkout.md) — 33% (5/15)\n- [ ] D01.T1 — kick-off · 1\n### D01.S1 — Payment — 35% (5/14)\n",
		"- [x] D01.T2 — cart spec · 3 · done 2026-10-12 by agent-a\n",
		"- [~] D01.T4 — payment provider · 8 · doing since 2026-10-12 by agent-c\n",
		"- [-] D01.T5 — refund rules · 5 · cancelled: the provider covers it\n",
		"- [ ] D01.T6 — legacy docs · 1 (estimate)\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if !sumRE.MatchString(strings.TrimRight(out[strings.LastIndex(strings.TrimRight(out, "\n"), "\n")+1:], "\n")) {
		t.Errorf("the last line is the sum:\n%s", out)
	}
	back, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if Render(back) != out {
		t.Errorf("parsing gives back the same file:\n%s\n---\n%s", out, Render(back))
	}
	if crlf, err := Parse(strings.ReplaceAll(out, "\n", "\r\n")); err != nil || Render(crlf) != out {
		t.Errorf("CRLF line endings read the same: %v", err)
	}
}

func TestTotals_weights(t *testing.T) {
	t.Run("DMLDD-B02: Progress counts the done weight over the weight not cancelled", func(t *testing.T) {})
	tasks := []*Task{{Weight: 3, State: Done}, {Weight: 2, State: Done}, {Weight: 8, State: Doing}, {Weight: 5, State: Cancelled}}
	done, total := Totals(tasks)
	if done != 5 || total != 13 || Percent(done, total) != "38%" {
		t.Errorf("%v/%v %s", done, total, Percent(done, total))
	}
	if Percent(0, 0) != "—" {
		t.Error("nothing counted is a dash")
	}
}

func TestIDs_inOrder(t *testing.T) {
	t.Run("DMLDD-B03: Records are numbered in order", func(t *testing.T) {})
	f := &File{}
	d, _ := f.AddDemand("A", "")
	s1, _ := f.AddStage(d, "one")
	s2, _ := f.AddStage(d, "two")
	t1, _ := f.AddTask(s1, "a", 1, false)
	t2, _ := f.AddTask(d, "b", 1, false)
	t3, _ := f.AddTask(s2, "c", 1, false)
	if got := []string{d, s1, s2, t1, t2, t3}; !reflect.DeepEqual(got, []string{"D01", "D01.S1", "D01.S2", "D01.T1", "D01.T2", "D01.T3"}) {
		t.Errorf("ids %v", got)
	}
	if _, err := f.AddTask(d, "zero", 0, false); err == nil {
		t.Error("a task needs a positive weight")
	}
	if d2, _ := f.AddDemand("B", ""); d2 != "D02" {
		t.Errorf("second demand %s", d2)
	}
}

func TestTask_states(t *testing.T) {
	t.Run("DMLDD-B04: A task moves through its states", func(t *testing.T) {})
	f := &File{}
	d, _ := f.AddDemand("A", "")
	id, _ := f.AddTask(d, "x", 2, true)
	if err := f.Start(id, "a1", "2026-10-12", false); err != nil {
		t.Fatal(err)
	}
	if err := f.Start(id, "a2", "2026-10-12", false); err == nil || !strings.Contains(err.Error(), "held by a1") {
		t.Errorf("another agent is refused: %v", err)
	}
	if err := f.Start(id, "a2", "2026-10-13", true); err != nil || f.Task(id).By != "a2" {
		t.Errorf("taken over: %v", err)
	}
	if err := f.Done(id, "a2", "2026-10-14"); err != nil || f.Task(id).State != Done || f.Task(id).At != "2026-10-14" {
		t.Errorf("done: %v %+v", err, f.Task(id))
	}
	if err := f.Reopen(id); err != nil || f.Task(id).State != Todo || f.Task(id).By != "" {
		t.Errorf("reopened: %+v", f.Task(id))
	}
	if err := f.Cancel(id, "out of scope"); err != nil || f.Task(id).Reason != "out of scope" {
		t.Errorf("cancelled: %+v", f.Task(id))
	}
	if err := f.Done(id, "a", "x"); err == nil {
		t.Error("a cancelled task is reopened before it is done")
	}
	if err := f.Weigh(id, 5); err != nil || f.Task(id).Weight != 5 || f.Task(id).Estimate {
		t.Errorf("weighed: %+v", f.Task(id))
	}
	if len(f.AllTasks()) != 1 {
		t.Error("nothing is deleted")
	}
}

func TestUpdate_concurrentChangesAreKept(t *testing.T) {
	t.Run("DMLDD-B05: Changes made at once are all kept", func(t *testing.T) {})
	root := t.TempDir()
	if err := Update(root, "points", func(f *File) error { _, err := f.AddDemand("A", ""); return err }); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- Update(root, "points", func(f *File) error { _, err := f.AddTask("D01", "task", 1, false); return err })
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	f, err := Load(root)
	if err != nil || len(f.AllTasks()) != 20 {
		t.Errorf("kept %d of 20: %v", len(f.AllTasks()), err)
	}
	if empty, err := Load(t.TempDir()); err != nil || len(empty.Demands) != 0 {
		t.Error("an absent file is an empty one")
	}
}

func TestSummary_demandsAndStages(t *testing.T) {
	t.Run("DMLDD-B06: The summary shows demands and stages, not tasks", func(t *testing.T) {})
	f := sample(t)
	b, _ := f.AddDemand("Bugs", "")
	_, _ = f.AddTask(b, "flaky", 1, false)
	rows := Summary(f)
	if len(rows) != 3 || rows[0].ID != "D01" || rows[1].ID != "D01.S1" || !rows[1].Stage || rows[2].ID != "D02" {
		t.Fatalf("rows %+v", rows)
	}
	if rows[0].Percent != "33%" || rows[0].Open != 3 || rows[0].Doing != 1 || rows[0].Estimate != 1 || rows[2].Open != 1 {
		t.Errorf("rows %+v", rows)
	}
}

func TestMerge_byID(t *testing.T) {
	t.Run("DMLDD-B07: Two branches merge by id", func(t *testing.T) {})
	base := &File{}
	d, _ := base.AddDemand("A", "")
	_, _ = base.AddTask(d, "first", 1, false)
	clone := func() *File { f, _ := Parse(Render(base)); return f }
	ours, theirs := clone(), clone()
	_ = ours.Done("D01.T1", "a", "2026-10-12")
	_, _ = ours.AddTask(d, "ours", 2, false)
	_, _ = theirs.AddTask(d, "theirs", 3, false)
	s, _ := theirs.AddStage(d, "later")
	_, _ = theirs.AddTask(s, "staged", 1, false)
	m := Merge(base, ours, theirs)
	if m.Task("D01.T1").State != Done {
		t.Error("done beats open")
	}
	titles := map[string]string{}
	for _, task := range m.AllTasks() {
		titles[task.ID] = task.Title
	}
	if titles["D01.T2"] != "ours" || titles["D01.T3"] != "theirs" || len(m.Demand("D01").Stages) != 1 {
		t.Errorf("both kept, theirs renumbered, their stage joined: %v", titles)
	}
	if _, err := Parse(Render(m)); err != nil {
		t.Errorf("the merge renders a readable file: %v", err)
	}
	if !IsFile(Render(m)) || IsFile("# Progress\n") {
		t.Error("the file is recognized by its header")
	}
}

func TestPlanDemand(t *testing.T) {
	t.Run("DMLDD-B08: A plan becomes a demand", func(t *testing.T) {})
	f := &File{}
	plan := "<!-- @anchors\n  code: CHKPL\n-->\n# Checkout\n\n## CHKPL-W01 — Foundation\n\n## CHKPL-W02 — Payment\n"
	id, err := f.PlanDemand("plans/0012.md", plan, "")
	if err != nil {
		t.Fatal(err)
	}
	d := f.Demand(id)
	if d.Title != "Checkout" || d.Plan != "plans/0012.md" || len(d.Stages) != 2 || d.Stages[1].Title != "CHKPL-W02 — Payment" {
		t.Errorf("demand %+v", d)
	}
	if _, err := f.PlanDemand("plans/0012.md", plan, ""); err == nil {
		t.Error("a tracked plan is refused")
	}
}

func TestMigrateProject(t *testing.T) {
	t.Run("DMLDD-B09: Per-plan progress files are migrated", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(root, "plans", "0001-x.md"), []byte("# X\n"), 0o644)
	legacy := "<!-- anchors:progress -->\n\n# Progress — XXXXX\n\n## XXXXX-W01 — Base\n\n- [x] `a.spec.md` — done\n- [ ] `b.spec.md` — open\n- [ ] TODO: one item per spec\n"
	_ = os.WriteFile(filepath.Join(root, "plans", "0001-x-progress.md"), []byte(legacy), 0o644)
	if m, err := MigrateProject(root, "", "2026-10-12", true); err != nil || len(m) != 1 || m[0].Tasks != 2 {
		t.Fatalf("dry run: %+v %v", m, err)
	}
	if _, err := os.Stat(Path(root)); err == nil {
		t.Error("a dry run writes nothing")
	}
	m, err := MigrateProject(root, "", "2026-10-12", false)
	if err != nil || len(m) != 1 || m[0].Done != 1 {
		t.Fatalf("%+v %v", m, err)
	}
	f, _ := Load(root)
	d := f.Demands[0]
	if d.Plan != "plans/0001-x.md" || d.Title != "X" || len(d.Stages) != 1 || len(d.Stages[0].Tasks) != 2 || !d.Stages[0].Tasks[0].Estimate || d.Stages[0].Tasks[0].State != Done {
		t.Errorf("demand %+v", d)
	}
	if _, err := os.Stat(filepath.Join(root, "plans", "0001-x-progress.md")); err == nil {
		t.Error("the per-plan file is deleted")
	}
}

func TestParse_handEditRefused(t *testing.T) {
	t.Run("DMLDD-E01: A hand edit is refused", func(t *testing.T) {})
	root := t.TempDir()
	if err := Update(root, "", func(f *File) error { _, err := f.AddDemand("A", ""); return err }); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(Path(root))
	edited := strings.Replace(string(b), "## D01 — A", "## D01 — B", 1)
	_ = os.WriteFile(Path(root), []byte(edited), 0o644)
	if _, err := Load(root); !errors.Is(err, ErrEdited) {
		t.Errorf("read: %v", err)
	}
	if err := Update(root, "", func(f *File) error { return nil }); !errors.Is(err, ErrEdited) {
		t.Errorf("change: %v", err)
	}
	if after, _ := os.ReadFile(Path(root)); string(after) != edited {
		t.Error("nothing is written over a hand edit")
	}
}

func TestChanges_refused(t *testing.T) {
	t.Run("DMLDD-E02: What the format cannot read back is refused", func(t *testing.T) {})
	f := &File{}
	if _, err := f.AddDemand("a · b", ""); err == nil {
		t.Error("the separator is refused")
	}
	d, _ := f.AddDemand("A", "")
	id, _ := f.AddTask(d, "x", 1, false)
	if err := f.Cancel(id, " "); err == nil {
		t.Error("an empty reason is refused")
	}
	if err := f.Done("D09.T1", "a", "x"); err == nil || !strings.Contains(err.Error(), "no task D09.T1") {
		t.Errorf("an unknown task: %v", err)
	}
}

func TestTask_since(t *testing.T) {
	t.Run("DMLDD-B10: A task in progress says when it was started", func(t *testing.T) {})
	f := &File{}
	id, _ := f.AddDemand("A", "")
	tid, _ := f.AddTask(id, "x", 1, false)
	if err := f.Start(tid, "a", "2026-10-12 21:31", false); err != nil {
		t.Fatal(err)
	}
	g, err := Parse(Render(f))
	if err != nil {
		t.Fatal(err)
	}
	at, ok := g.Task(tid).Since()
	if !ok || at.Format(StartedAt) != "2026-10-12 21:31" || !strings.Contains(Render(g), "doing since 2026-10-12 21:31 by a") {
		t.Errorf("the minute is kept: %v %v\n%s", at, ok, Render(g))
	}
	old := &Task{At: "2026-10-12"}
	if at, ok := old.Since(); !ok || at.Format(StartedAt) != "2026-10-12 00:00" {
		t.Errorf("a date alone counts from the start of its day: %v %v", at, ok)
	}
	text := strings.Replace(Render(g), "2026-10-12 21:31", "2026-10-12", 1)
	if _, err := Parse(fixSum(text)); err != nil {
		t.Errorf("a line with only the date still reads: %v", err)
	}
}

// fixSum gives a text edited by a test the sum Anchors would have written.
func fixSum(text string) string {
	i := strings.LastIndex(text, "<!-- anchors:progress sum=")
	content := text[:i]
	return content + "<!-- anchors:progress sum=" + sumOf(content) + " -->\n"
}

func TestRename(t *testing.T) {
	t.Run("DMLDD-B11: A demand, a stage or a task is renamed", func(t *testing.T) {})
	f := &File{}
	id, _ := f.AddDemand("A", "")
	sid, _ := f.AddStage(id, "Sob demanda")
	tid, _ := f.AddTask(sid, "x", 1, false)
	if err := f.Done(tid, "a", "2026-10-13"); err != nil {
		t.Fatal(err)
	}
	for _, c := range [][2]string{{id, "Checkout"}, {sid, "Não planejadas"}, {tid, "y"}} {
		if err := f.Rename(c[0], c[1]); err != nil {
			t.Fatalf("%s: %v", c[0], err)
		}
	}
	g, err := Parse(Render(f))
	if err != nil {
		t.Fatal(err)
	}
	d := g.Demand(id)
	task := g.Task(tid)
	if d.Title != "Checkout" || d.Stages[0].Title != "Não planejadas" || d.Stages[0].ID != sid || task.Title != "y" || task.State != Done || task.By != "a" {
		t.Errorf("renamed, the rest kept: %+v %+v", d, task)
	}
	if err := f.Rename("D09", "z"); err == nil {
		t.Error("an id that names nothing is refused")
	}
	if err := f.Rename(tid, "a · b"); err == nil {
		t.Error("a title with the separator is refused, as a new one")
	}
}

func TestMarkUnplanned(t *testing.T) {
	t.Run("DMLDD-B12: A task is marked unplanned", func(t *testing.T) {})
	f := &File{}
	id, _ := f.AddDemand("A", "")
	sid, _ := f.AddStage(id, "S")
	t1, _ := f.AddTask(sid, "fix", 2, false)
	t2, _ := f.AddTask(sid, "helper", 3, true)
	t3, _ := f.AddTask(sid, "gone", 5, false)
	_, _ = f.AddTask(sid, "planned", 10, false)
	for _, x := range []string{t1, t2, t3} {
		if err := f.MarkUnplanned(x, true); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Cancel(t3, "no need"); err != nil {
		t.Fatal(err)
	}
	text := Render(f)
	if !strings.Contains(text, "fix · 2 (unplanned)") || !strings.Contains(text, "helper · 3 (estimate, unplanned)") {
		t.Errorf("the marks in the lines:\n%s", text)
	}
	g, err := Parse(text)
	if err != nil || !g.Task(t1).Unplanned || !g.Task(t2).Unplanned || !g.Task(t2).Estimate {
		t.Fatalf("read back: %v", err)
	}
	if rows := Summary(g); rows[0].Unplanned != 5 || rows[1].Unplanned != 5 {
		t.Errorf("the unplanned weight, cancelled left out: %+v", rows)
	}
	if err := g.MarkUnplanned(t1, false); err != nil || g.Task(t1).Unplanned {
		t.Errorf("the mark comes off: %v", err)
	}
	base, _ := Parse(text)
	ours, _ := Parse(text)
	theirs, _ := Parse(text)
	_ = theirs.MarkUnplanned(t1, false)
	merged := Merge(base, ours, theirs)
	if merged.Task(t1).Unplanned {
		t.Error("the merge keeps the side that changed the mark")
	}
	bad := strings.Replace(text, "fix · 2 (unplanned)", "fix · 2 (urgent)", 1)
	if _, err := Parse(fixSum(bad)); err == nil || !strings.Contains(err.Error(), `"urgent"`) {
		t.Errorf("an unknown mark is refused, named: %v", err)
	}
	if !strings.Contains(UnplannedAdvice, "--unplanned") {
		t.Error("the advice says the command")
	}
}

func TestMoveAndDropStage(t *testing.T) {
	t.Run("DMLDD-B13: A task moves within its demand, and an empty stage is removed", func(t *testing.T) {})
	f := &File{}
	id, _ := f.AddDemand("A", "")
	s1, _ := f.AddStage(id, "Payment")
	_, _ = f.AddTask(s1, "refunds", 5, false)
	found, _ := f.AddStage(id, "Found on the way")
	tid, _ := f.AddTask(found, "fix the parser", 2, false)
	_ = f.MarkUnplanned(tid, true)
	_ = f.Done(tid, "a", "2026-10-13")
	other, _ := f.AddDemand("B", "")
	if err := f.Relocate(tid, s1); err != nil {
		t.Fatal(err)
	}
	g, err := Parse(Render(f))
	if err != nil {
		t.Fatalf("%v\n%s", err, Render(f))
	}
	st := g.Demand(id).Stages[0]
	if len(st.Tasks) != 2 || st.Tasks[1].ID != tid || st.Tasks[1].State != Done || !st.Tasks[1].Unplanned || len(g.Demand(id).Stages[1].Tasks) != 0 {
		t.Errorf("moved to the stage, kept: %+v", st.Tasks)
	}
	if err := f.Relocate(tid, id); err != nil || len(f.Demand(id).Tasks) != 1 {
		t.Errorf("to the demand itself: %v", err)
	}
	if err := f.Relocate(tid, other); err == nil {
		t.Error("another demand is refused")
	}
	if err := f.DropStage(s1); err == nil {
		t.Error("a stage still holding a task is refused")
	}
	if err := f.DropStage(found); err != nil || len(f.Demand(id).Stages) != 1 {
		t.Errorf("the empty stage goes: %v", err)
	}
	if _, err := Parse(Render(f)); err != nil {
		t.Errorf("the file still reads: %v", err)
	}
}

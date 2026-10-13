// @anchors
//   code: PCTPR
//   ref: LDCML

package flow

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/progress"
)

// progressProject is a project naming its unit, with a plan of two phases.
func progressProject(t *testing.T, unit string) string {
	t.Helper()
	root := t.TempDir()
	yaml := "version: 7\nlayers: {}\n"
	if unit != "" {
		yaml += "progress:\n  unit: " + unit + "\n"
	}
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plan := "<!-- @anchors\n  code: CHKPL\n-->\n# Checkout\n\n## CHKPL-W01 — Foundation\n\n## CHKPL-W02 — Payment\n"
	if err := os.WriteFile(filepath.Join(root, "plans", "0012.md"), []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// prog runs `anchors progress` with the args, in the project.
func prog(t *testing.T, root string, args ...string) (string, error) {
	t.Helper()
	var out strings.Builder
	c := newProgressTrackCmd()
	c.SetArgs(append(args, "--root", root))
	c.SetOut(&out)
	c.SetErr(io.Discard)
	err := c.Execute()
	return out.String(), err
}

func mustProg(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := prog(t, root, args...)
	if err != nil {
		t.Fatalf("progress %v: %v", args, err)
	}
	return out
}

func TestProgressCmd_adds(t *testing.T) {
	t.Run("LDCML-B01: Demands, stages and tasks are added", func(t *testing.T) {})
	root := progressProject(t, "hours")
	if out := mustProg(t, root, "add", "demand", "Bug triage"); !strings.Contains(out, "demand D01 — Bug triage") {
		t.Errorf("%q", out)
	}
	if out := mustProg(t, root, "add", "demand", "--plan", "plans/0012.md"); !strings.Contains(out, "demand D02 — the plan plans/0012.md") {
		t.Errorf("%q", out)
	}
	mustProg(t, root, "add", "stage", "D01", "Triage")
	if out := mustProg(t, root, "add", "task", "D01.S1", "flaky login", "--weight", "2"); !strings.Contains(out, "task D01.T1 — flaky login") {
		t.Errorf("%q", out)
	}
	f, _ := progress.Load(root)
	if d := f.Demand("D02"); d == nil || d.Title != "Checkout" || len(d.Stages) != 2 {
		t.Errorf("the plan's demand: %+v", d)
	}
}

func TestProgressCmd_moves(t *testing.T) {
	t.Run("LDCML-B02: Tasks move by command", func(t *testing.T) {})
	root := progressProject(t, "")
	mustProg(t, root, "add", "demand", "A")
	mustProg(t, root, "add", "task", "D01", "x", "--weight", "3")
	mustProg(t, root, "start", "D01.T1", "--by", "a1")
	if _, err := prog(t, root, "start", "D01.T1", "--by", "a2"); err == nil || !strings.Contains(err.Error(), "held by a1") {
		t.Errorf("another agent: %v", err)
	}
	mustProg(t, root, "start", "D01.T1", "--by", "a2", "--take")
	mustProg(t, root, "done", "D01.T1", "--by", "a2")
	f, _ := progress.Load(root)
	if task := f.Task("D01.T1"); task.State != progress.Done || task.By != "a2" || task.At == "" {
		t.Errorf("done: %+v", task)
	}
	mustProg(t, root, "reopen", "D01.T1")
	mustProg(t, root, "cancel", "D01.T1", "--reason", "dropped")
	mustProg(t, root, "weigh", "D01.T1", "5")
	f, _ = progress.Load(root)
	if task := f.Task("D01.T1"); task.State != progress.Cancelled || task.Reason != "dropped" || task.Weight != 5 {
		t.Errorf("after: %+v", task)
	}
}

func TestProgressCmd_summary(t *testing.T) {
	t.Run("LDCML-B03: The summary shows demands and stages", func(t *testing.T) {})
	root := progressProject(t, "hours")
	if out := mustProg(t, root, "summary"); !strings.Contains(out, "nothing in progress") {
		t.Errorf("empty: %q", out)
	}
	mustProg(t, root, "add", "demand", "--plan", "plans/0012.md")
	mustProg(t, root, "add", "task", "D01.S1", "cart", "--weight", "3")
	mustProg(t, root, "add", "task", "D01.S2", "provider", "--weight", "5", "--estimate")
	mustProg(t, root, "done", "D01.T1", "--by", "a")
	out := mustProg(t, root, "summary")
	for _, want := range []string{
		"Progress — 37% (3/8 hours)",
		"  D01 Checkout (plans/0012.md) — 37% (3/8) · 1 open, 1 weight(s) still estimates",
		"      D01.S1 CHKPL-W01 — Foundation — 100% (3/3) · 0 open",
		"      D01.S2 CHKPL-W02 — Payment — 0% (0/5) · 1 open",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "cart") || strings.Contains(out, "provider") {
		t.Errorf("no task in the summary:\n%s", out)
	}
	var rows []progress.Row
	if err := json.Unmarshal([]byte(mustProg(t, root, "summary", "--json")), &rows); err != nil || len(rows) != 3 {
		t.Errorf("json rows: %v %+v", err, rows)
	}
}

func TestProgressCmd_show(t *testing.T) {
	t.Run("LDCML-B04: Show prints the tasks", func(t *testing.T) {})
	root := progressProject(t, "")
	mustProg(t, root, "add", "demand", "A")
	mustProg(t, root, "add", "task", "D01", "alpha", "--weight", "1")
	mustProg(t, root, "add", "demand", "B")
	mustProg(t, root, "add", "task", "D02", "beta", "--weight", "1")
	all := mustProg(t, root, "show")
	if !strings.Contains(all, "- [ ] D01.T1 — alpha · 1") || !strings.Contains(all, "- [ ] D02.T1 — beta · 1") || strings.Contains(all, "anchors:progress") {
		t.Errorf("all:\n%s", all)
	}
	if one := mustProg(t, root, "show", "D02"); strings.Contains(one, "alpha") || !strings.Contains(one, "beta") {
		t.Errorf("one demand:\n%s", one)
	}
}

func TestProgressCmd_unit(t *testing.T) {
	t.Run("LDCML-B05: The unit is the project's", func(t *testing.T) {})
	root := progressProject(t, "")
	mustProg(t, root, "add", "demand", "A")
	if b, _ := os.ReadFile(progress.Path(root)); !strings.Contains(string(b), "(0/0 points)") {
		t.Errorf("%s", b)
	}
}

func TestProgressCmd_refused(t *testing.T) {
	t.Run("LDCML-E01: A refused change writes nothing", func(t *testing.T) {})
	root := progressProject(t, "")
	mustProg(t, root, "add", "demand", "A")
	before, _ := os.ReadFile(progress.Path(root))
	if _, err := prog(t, root, "done", "D01.T9"); err == nil || !strings.Contains(err.Error(), "no task D01.T9") {
		t.Errorf("unknown task: %v", err)
	}
	if _, err := prog(t, root, "add", "demand"); err == nil || !strings.Contains(err.Error(), "needs a title, or --plan") {
		t.Errorf("no title: %v", err)
	}
	if after, _ := os.ReadFile(progress.Path(root)); string(after) != string(before) {
		t.Error("a refused change writes nothing")
	}
}

func TestProgressCmd_rename(t *testing.T) {
	t.Run("LDCML-B06: A stage is renamed by command", func(t *testing.T) {})
	root := progressProject(t, "")
	mustProg(t, root, "add", "demand", "Checkout")
	mustProg(t, root, "add", "stage", "D01", "Sob demanda")
	mustProg(t, root, "add", "task", "D01.S1", "x", "--weight", "2")
	if out := mustProg(t, root, "rename", "D01.S1", "Não planejadas"); !strings.Contains(out, "✎ D01.S1 — Não planejadas") {
		t.Errorf("said: %s", out)
	}
	if out := mustProg(t, root, "summary"); !strings.Contains(out, "D01.S1 Não planejadas — 0% (0/2)") {
		t.Errorf("the new title, the same id and task: %s", out)
	}
}

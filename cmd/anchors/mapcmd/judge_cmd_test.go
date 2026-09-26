package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
)

// judgmentsOn returns the stamps the map carries for gate on the edges that reach target.
func judgmentsOn(t *testing.T, root, target, gate string) []string {
	t.Helper()
	var out []string
	for _, e := range loadMap(t, root).Edges {
		if e.To != target {
			continue
		}
		for _, j := range e.Julgamentos {
			if j.Gate == gate {
				out = append(out, e.From+"→"+j.Verdict)
			}
		}
	}
	return out
}

// edgeStamp returns the verdict of the stamp on the edge from→to, or "" when unstamped.
func edgeStamp(t *testing.T, root, from, to string) string {
	t.Helper()
	for _, e := range loadMap(t, root).Edges {
		if e.From == from && e.To == to && e.Stamp != nil {
			return e.Stamp.Verdict
		}
	}
	return ""
}

func issueFiles(t *testing.T, root, state string) []string {
	t.Helper()
	entries, _ := os.ReadDir(filepath.Join(root, "issues", state))
	var out []string
	for _, e := range entries {
		out = append(out, e.Name())
	}
	return out
}

// The whole cycle on a gate that declares its guide: FAIL stamps the guide→target edge
// and opens the issue with the report; the same finding again changes nothing; a NEW
// finding reopens it; PASS resolves it and stamps `ok`.
func TestJudge_failThenPassOnTheGuideEdge(t *testing.T) {
	root := fixtureProject(t)
	report := "## Report\n- the login assembles its block inline (l.3)"

	out := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "FAIL", "--reason", report)
	if !strings.Contains(out, "issue opened at issues/todo/") || !strings.Contains(out, "stamped: 1 edge(s)") {
		t.Errorf("unexpected output of the fail:\n%s", out)
	}
	if got := edgeStamp(t, root, "guides/CODE.md", "src/login.ts"); got != "issue" {
		t.Errorf("the guide→target edge was not stamped `issue`: %q", got)
	}
	todo := issueFiles(t, root, "todo")
	if len(todo) != 1 {
		t.Fatalf("expected one issue in todo, got %v", todo)
	}
	body, _ := os.ReadFile(filepath.Join(root, "issues", "todo", todo[0]))
	if !strings.Contains(string(body), "assembles its block inline") {
		t.Errorf("the issue body is not the report:\n%s", body)
	}

	same := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", report)
	if !strings.Contains(same, "same finding already recorded") {
		t.Errorf("the same finding twice must change nothing:\n%s", same)
	}
	again := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "a second, different defect")
	if !strings.Contains(again, "NEW finding added") {
		t.Errorf("a new finding must reopen and append:\n%s", again)
	}

	pass := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "pass")
	if !strings.Contains(pass, "previous issue resolved") {
		t.Errorf("a pass must resolve the open issue:\n%s", pass)
	}
	if len(issueFiles(t, root, "todo")) != 0 || len(issueFiles(t, root, "done")) != 1 {
		t.Errorf("the issue did not move to done: todo=%v done=%v", issueFiles(t, root, "todo"), issueFiles(t, root, "done"))
	}
	if got := edgeStamp(t, root, "guides/CODE.md", "src/login.ts"); got != "ok" {
		t.Errorf("the pass did not restamp the edge `ok`: %q", got)
	}
	// A second pass has nothing to resolve.
	if again := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "pass"); !strings.Contains(again, "✓ judged PASS\n") {
		t.Errorf("a pass with no issue open is a plain PASS:\n%s", again)
	}
}

// WAIVED is its own stamp and its own word — never announced as PASS.
func TestJudge_waivedIsNotAPass(t *testing.T) {
	root := fixtureProject(t)
	runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "x")
	out := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "waived",
		"--reason", "the spec declares @TBD: code")
	if !strings.Contains(out, "judged WAIVED") || !strings.Contains(out, "previous issue resolved") || strings.Contains(out, "PASS") {
		t.Errorf("unexpected output of the waiver:\n%s", out)
	}
	if got := edgeStamp(t, root, "guides/CODE.md", "src/login.ts"); got != "waived" {
		t.Errorf("the stamp must be `waived`, not `ok`: %q", got)
	}
	plain := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "waived", "--reason", "still @TBD")
	if !strings.Contains(plain, "○ judged WAIVED — there was nothing to confront\n") {
		t.Errorf("a waiver with nothing to resolve:\n%s", plain)
	}
}

// The legacy `--verdict dispensado` is a waiver too, from ValidateVerdict to the stamp.
func TestJudge_legacyDispensadoIsAWaiver(t *testing.T) {
	root := fixtureProject(t)
	out := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "dispensado",
		"--reason", "the spec declares @TBD: code")
	if strings.Contains(out, "PASS") || !strings.Contains(out, "WAIVED") {
		t.Errorf("the legacy value must judge WAIVED, never PASS:\n%s", out)
	}
	if got := edgeStamp(t, root, "guides/CODE.md", "src/login.ts"); got != "waived" {
		t.Errorf("the legacy value must stamp `waived`, not `ok`: %q", got)
	}
}

// `review` is accepted without being declared, and stamps per NODE (no guide edge); the
// judge task in the queue is closed.
func TestJudge_reviewStampsTheNodeAndClosesTheTask(t *testing.T) {
	root := fixtureProject(t)
	if _, err := queue.Enqueue(root, queue.Task{ID: "judge-review-src-login", Changed: "src/login.ts", SuggestedNext: "judge", Reason: "delivered"}); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Enqueue(root, queue.Task{ID: "judge-atomic-src-login", Changed: "src/login.test.ts", SuggestedNext: "judge", Reason: "changed"}); err != nil {
		t.Fatal(err)
	}
	pending := runCmd(t, newJudgeCmd(), "--pending", "--root", root)
	if !strings.Contains(pending, "○ src/login.ts\n    delivered") || !strings.Contains(pending, "2 target(s)") {
		t.Errorf("--pending does not list the queued judgments:\n%s", pending)
	}

	out := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "review", "--verdict", "pass")
	if !strings.Contains(out, "stamped: 2 edge(s)") {
		t.Errorf("review stamps every edge of the node:\n%s", out)
	}
	if got := judgmentsOn(t, root, "src/login.ts", "review"); len(got) != 2 {
		t.Errorf("expected the two edges into the code stamped by review: %v", got)
	}
	if _, err := os.Stat(filepath.Join(root, queue.DoneDir, "done__judge-review-src-login.yaml")); err != nil {
		t.Errorf("the judge task was not closed: %v", err)
	}
	left := runCmd(t, newJudgeCmd(), "--pending", "--root", root)
	if strings.Contains(left, "○ src/login.ts\n") || !strings.Contains(left, "1 target(s)") {
		t.Errorf("the closed task is still pending:\n%s", left)
	}
}

func TestJudge_nothingPending(t *testing.T) {
	root := fixtureProject(t)
	if out := runCmd(t, newJudgeCmd(), "--pending", "--root", root); !strings.Contains(out, "no target awaiting judgment") {
		t.Errorf("an empty queue must say so:\n%s", out)
	}
}

// The code of a new unit is not written yet: the judgment lands on the piece that exists.
func TestJudge_fallsBackToTheUnitsExistingPiece(t *testing.T) {
	root := fixtureProject(t)
	t.Chdir(root)
	out := runCmd(t, newJudgeCmd(), "src/login.tsx", "--gate", "review", "--verdict", "pass")
	if !strings.Contains(out, "recording on `src/login.spec.md`") {
		t.Errorf("the judgment did not fall back to the spec:\n%s", out)
	}
	if _, err := runCmdErr(newJudgeCmd(), t, "src/signup.tsx", "--gate", "review", "--verdict", "pass"); err == nil ||
		!strings.Contains(err.Error(), `"src/signup.tsx" is not in the map`) {
		t.Errorf("a unit with no piece in the map: got %v", err)
	}
}

// Manual mode stamps the map and prints the report, and writes no issue unless asked.
func TestJudge_manualModeWritesNoIssue(t *testing.T) {
	root := fixtureProjectWith(t, "workflow:\n  mode: manual\n", fixtureSpec)
	out := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail",
		"--reason", "  the report itself  ")
	if !strings.Contains(out, "no issue written (mode: manual") || !strings.Contains(out, "\nthe report itself\n") {
		t.Errorf("manual mode must print the report instead of writing it:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "issues")); err == nil {
		t.Error("manual mode wrote an issue without --record-issues")
	}
	if got := edgeStamp(t, root, "guides/CODE.md", "src/login.ts"); got != "issue" {
		t.Errorf("manual mode must still stamp the map: %q", got)
	}
	for verdict, want := range map[string]string{"pass": "✓ judged PASS", "waived": "judged WAIVED"} {
		o := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", verdict, "--reason", "r")
		if !strings.Contains(o, want) {
			t.Errorf("%s in manual mode:\n%s", verdict, o)
		}
	}
	runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "r", "--record-issues")
	if len(issueFiles(t, root, "todo")) != 1 {
		t.Error("--record-issues must write the issue in manual mode")
	}
}

// --patch turns the fix into an applicable suggestion.
func TestJudge_patchOpensASuggestion(t *testing.T) {
	root := fixtureProject(t)
	patch := filepath.Join(t.TempDir(), "fix.diff")
	if err := os.WriteFile(patch, []byte("--- a/src/login.ts\n+++ b/src/login.ts\n@@ -1 +1 @@\n-x\n+y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "inline block", "--patch", patch)
	if !strings.Contains(out, "anchors suggest show judge-atomic-src-login") {
		t.Errorf("the suggestion was not announced:\n%s", out)
	}
	b, err := os.ReadFile(filepath.Join(root, "suggestions", "pending", "judge-atomic-src-login.md"))
	if err != nil || !strings.Contains(string(b), "+y") {
		t.Errorf("the suggestion does not carry the patch: %v\n%s", err, b)
	}
	if _, err := runCmdErr(newJudgeCmd(), t, "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "r", "--patch", patch+".missing"); err == nil ||
		!strings.Contains(err.Error(), "read the patch") {
		t.Errorf("a missing patch file: got %v", err)
	}
}

func TestJudge_refusals(t *testing.T) {
	root := fixtureProject(t)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--root", root}, "provide the target"},
		{[]string{"src/login.ts", "--root", root}, "--gate is mandatory"},
		{[]string{"src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "maybe"}, "--verdict must be"},
		{[]string{"src/login.ts", "--root", root, "--gate", "always-red", "--verdict", "pass"}, "not a judgment gate"},
		{[]string{"src/login.ts", "--root", t.TempDir(), "--gate", "atomic", "--verdict", "pass"}, "anchors map build"},
	} {
		_, err := runCmdErr(newJudgeCmd(), t, c.args...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: expected %q, got %v", c.args, c.want, err)
		}
	}
	// A map with no config beside it: the gate cannot be validated.
	noCfg := t.TempDir()
	writeProjectFile(t, noCfg, "src/login.ts", "export const login = () => true\n")
	if err := mapx.Save(loadMap(t, root), filepath.Join(noCfg, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmdErr(newJudgeCmd(), t, "src/login.ts", "--root", noCfg, "--gate", "atomic", "--verdict", "pass"); err == nil ||
		!strings.Contains(err.Error(), "load config") {
		t.Errorf("no config: got %v", err)
	}
}

// A judgment gate that declares `guide:` is answered once judged: the guide→target edge
// records the gate, and the next check does not ask again.
func TestJudge_guideGateCountsAsAnswered(t *testing.T) {
	root := fixtureProject(t)
	runCmd(t, newJudgeCmd(), "src/login.ts", "--root", root, "--gate", "atomic", "--verdict", "fail", "--reason", "too big")
	if v, ok := loadMap(t, root).JudgedBy("src/login.ts", "atomic"); !ok || v != "issue" {
		t.Fatalf("the judged gate must be answered on the guide edge: verdict %q, answered %v", v, ok)
	}
}

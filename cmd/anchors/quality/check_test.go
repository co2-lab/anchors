package quality

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
)

// The check ends a blocking failure (and a queued judgment on `--changed`) with
// os.Exit(1), which an in-process test cannot survive. TestCheckChildProcess plays the
// check in a child copy of the test binary; runCheckInChild starts it and reads the
// exit code and everything it printed.
func TestCheckChildProcess(t *testing.T) {
	args := os.Getenv("ANCHORS_CHECK_CHILD")
	if args == "" {
		t.Skip("only runs as the child of runCheckInChild")
	}
	cmd := newCheckCmd()
	cmd.SetArgs(strings.Split(args, "\x1f"))
	cmd.SilenceUsage = true
	if err := cmd.Execute(); err != nil {
		os.Stdout.WriteString("child error: " + err.Error() + "\n")
		os.Exit(2)
	}
	os.Exit(0)
}

func runCheckInChild(t *testing.T, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestCheckChildProcess$")
	cmd.Env = append(os.Environ(), "ANCHORS_CHECK_CHILD="+strings.Join(args, "\x1f"))
	out, err := cmd.CombinedOutput()
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), string(out)
	}
	if err != nil {
		t.Fatal(err)
	}
	return 0, string(out)
}

const blockingYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
gates:
  - name: must-pass
    on: [code]
    blocking: true
    run: "echo broken; exit 1"
`

func TestCheckBlockingFailureExitsOne(t *testing.T) {
	t.Run("CGPCH-B39: A blocking failure exits 1 with the mirror complete", func(t *testing.T) {})
	dir := qProject(t, blockingYAML, map[string]string{"a.go": "package a\n"},
		&mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}})
	code, out := runCheckInChild(t, "--root", dir, "--all", "--no-record")
	if code != 1 {
		t.Fatalf("a blocking failure must exit 1, got %d:\n%s", code, out)
	}
	mirror, err := filepath.Glob(filepath.Join(dir, ".anchors", "check*.txt"))
	if err != nil || len(mirror) == 0 {
		t.Fatalf("no output mirror: %v %v", mirror, err)
	}
	if body := readQ(t, mirror[0]); !strings.Contains(body, "✗ blocked — 1 blocking gate(s) failed") {
		t.Errorf("the mirror must hold the end of the report, the verdict included:\n%s", body)
	}
}

const judgmentYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
gates:
  - name: rule-kept
    on: [code]
    measures: judgment
    guide: SPEC.md
    ask: "Does the code do what the rule says?"
`

func TestCheckQueuedJudgmentBarsOnlyTheChangedCheck(t *testing.T) {
	t.Run("CGPCH-B21: A pending judgment becomes one task in the local queue", func(t *testing.T) {})
	t.Run("CGPCH-B22: A queued judgment bars the incremental check only", func(t *testing.T) {})
	t.Run("CGPCH-B69: A judge task suggests the review stage, a verb the work command composes", func(t *testing.T) {})
	dir := qProject(t, judgmentYAML, map[string]string{"a.go": "package a\n"},
		&mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}})
	code, out := runCheckInChild(t, "--root", dir, "--all")
	if code != 0 {
		t.Fatalf("a queued judgment must not bar --all, got %d:\n%s", code, out)
	}
	tasks, _ := queue.List(dir)
	if len(tasks) != 1 || tasks[0].ID != "judge-rule-kept-a" || tasks[0].Kind != "judgment" {
		t.Fatalf("one judge task per judged target, got %+v", tasks)
	}
	// whoever claims it with `anchors next` composes the prompt with `anchors work <verb>`:
	// a verb `work` refuses leaves the task with no route.
	if !queue.ValidWorkArtifact(tasks[0].SuggestedNext) || tasks[0].SuggestedNext != "review" {
		t.Fatalf("the judge task must suggest the review stage, a verb `anchors work` accepts; got %q", tasks[0].SuggestedNext)
	}
	if !strings.Contains(tasks[0].Reason, "anchors judge a.go --gate rule-kept") {
		t.Errorf("the reason must name the command that closes the task, with its gate: %q", tasks[0].Reason)
	}
	// the gate declares a guide and a question, and whoever judges needs both
	if !strings.Contains(tasks[0].Reason, "SPEC.md") || !strings.Contains(tasks[0].Reason, "Does the code do what the rule says?") {
		t.Errorf("the reason must carry the gate's guide and question: %q", tasks[0].Reason)
	}
	code, out = runCheckInChild(t, "--root", dir, "--changed", "a.go")
	if code != 1 || !strings.Contains(out, "awaiting judgment") {
		t.Errorf("a queued judgment must bar --changed with exit 1, got %d:\n%s", code, out)
	}
}

func TestCheckJudgmentBriefInGitHubModeEvenWithoutRecording(t *testing.T) {
	t.Run("CGPCH-B25: In github mode the brief is printed without recording and nothing is queued", func(t *testing.T) {})
	englishOutput(t)
	t.Setenv("GITHUB_ACTIONS", "")
	defer issue.UseFiles()
	dir := qProject(t, judgmentYAML+githubWorkflowYAML, map[string]string{"a.go": "package a\n"},
		&mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}})
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"⚖  1 target(s) awaiting judgment, across 1 gate(s).",
		"   rule-kept — 1 target(s)\n      guide: SPEC.md\n      asks: Does the code do what the rule says?\n      a.go\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// and it queues nothing in github mode, recording or not: the queue there is the board
	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all"); err != nil {
		t.Fatal(err)
	}
	if tasks, _ := queue.List(dir); len(tasks) != 0 {
		t.Errorf("github mode must not queue judgment locally, got %+v", tasks)
	}
}

const githubWorkflowYAML = "workflow:\n  mode: github\n  repo: o/r\n  labels: [anchors]\n"

func TestCheckPrintsTheLocalBacklogOnTheFullLocalSweepOnly(t *testing.T) {
	t.Run("CGPCH-B36: The local backlog is printed on the full local sweep only", func(t *testing.T) {})
	englishOutput(t)
	t.Setenv("GITHUB_ACTIONS", "")
	defer issue.UseFiles()
	withIssue := func(yaml string) string {
		dir := qProject(t, checkYAML+yaml, checkFiles(), checkGraph())
		if _, _, err := issue.Open(dir, issue.Issue{Kind: issue.Violation, Gate: "g", Target: "a.go", Date: "2026-09-25"}); err != nil {
			t.Fatal(err)
		}
		return dir
	}
	local := withIssue("")
	out, err := runQ(t, newCheckCmd(), "--root", local, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "still open locally") {
		t.Errorf("the full local sweep must print the backlog:\n%s", out)
	}
	out, err = runQ(t, newCheckCmd(), "--root", local, "--changed", "a.go", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "still open locally") {
		t.Errorf("an incremental check must not print the backlog:\n%s", out)
	}
	gh := withIssue(githubWorkflowYAML)
	out, err = runQ(t, newCheckCmd(), "--root", gh, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "still open locally") {
		t.Errorf("in github mode the backlog is the board, not printed here:\n%s", out)
	}
}

func TestCheckGovernanceTipsOnTheFullSweepOnly(t *testing.T) {
	t.Run("CGPCH-B37: Governance tips appear on the full sweep only", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "ℹ Governance Tip") || !strings.Contains(out, "anchors doctor") {
		t.Errorf("the full sweep must print the governance tips:\n%s", out)
	}
	out, err = runQ(t, newCheckCmd(), "--root", dir, "--changed", "a.go", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "Governance Tip") {
		t.Errorf("an incremental check must not print governance tips:\n%s", out)
	}
}

func TestRecordCheck_passResolvesAndPendingOpensDecisionOrDebt(t *testing.T) {
	t.Run("CGPCH-B19: Passes resolve, decisions and debts open in their folders", func(t *testing.T) {})
	t.Run("CGPCH-B77: The record summary counts what the record did", func(t *testing.T) {})
	englishOutput(t)
	issue.UseFiles()
	root := t.TempDir()
	mapPath := filepath.Join(root, mapx.DefaultPath)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.ts", Kind: mapx.KindCode, Rev: "r1"}, {ID: "b.spec.md", Kind: mapx.KindSpec, Rev: "s1"}}}
	old := issue.Issue{Kind: issue.Violation, Gate: "header-conforms", Target: "a.ts", Date: "2026-09-11"}
	if _, _, err := issue.Open(root, old); err != nil {
		t.Fatal(err)
	}
	p := gate.Profile{Results: []gate.Result{
		{Gate: "header-conforms", Target: "a.ts", Verdict: gate.Pass},
		{Gate: "open-questions-resolved", Target: "b.spec.md", Verdict: gate.Pending, Decisão: true, Detail: "who decides?"},
		{Gate: "debt-declared", Target: "a.ts", Verdict: gate.Pending, Divida: true, Prazo: "next sprint", Detail: "debt"},
		{Gate: "undetermined", Target: "a.ts", Verdict: gate.Pending, Detail: "nothing to say"},
	}}
	out := captureStdout(t, func() {
		if err := recordCheck(root, mapPath, g, p, true, false); err != nil {
			t.Fatal(err)
		}
	})
	if st, _ := issue.Exists(root, old.Key()); st != issue.Done {
		t.Errorf("a pass must resolve the open issue of the same gate and target, it is %s", st)
	}
	dec := issue.Issue{Kind: issue.Decision, Gate: "open-questions-resolved", Target: "b.spec.md"}
	if st, _ := issue.Exists(root, dec.Key()); st != issue.Todo {
		t.Errorf("an open decision must open in todo, it is %q", st)
	}
	if l, _ := issue.ListByOwner(root, issue.Todo, issue.DonoUsuário); len(l) != 1 {
		t.Errorf("the decision is the user's, got %v", l)
	}
	debt := issue.Issue{Kind: issue.Violation, Gate: "debt-declared", Target: "a.ts"}
	if st, _ := issue.Exists(root, debt.Key()); st != issue.Future {
		t.Errorf("an assumed debt must open in future/, it is %q", st)
	}
	undetermined := issue.Issue{Kind: issue.Violation, Gate: "undetermined", Target: "a.ts"}
	if st, _ := issue.Exists(root, undetermined.Key()); st != "" {
		t.Errorf("a plain pending opens no issue, found it in %q", st)
	}
	if !strings.Contains(out, "; 1 new issue(s), 1 resolved") || !strings.Contains(out, "1 assumed debt(s) recorded") {
		t.Errorf("the record summary must count what it did:\n%s", out)
	}

	// the decision closes by its own kind when the gate passes again
	pass := gate.Profile{Results: []gate.Result{{Gate: "open-questions-resolved", Target: "b.spec.md", Verdict: gate.Pass}}}
	out = captureStdout(t, func() {
		if err := recordCheck(root, mapPath, g, pass, true, false); err != nil {
			t.Fatal(err)
		}
	})
	if st, _ := issue.Exists(root, dec.Key()); st != issue.Done {
		t.Errorf("a pass of open-questions-resolved must resolve the decision, it is %q", st)
	}
	// the closed decision is counted as resolved, and no debt line appears without a debt
	if !strings.Contains(out, "; 0 new issue(s), 1 resolved") || strings.Contains(out, "assumed debt") {
		t.Errorf("the summary must count the resolved decision and nothing else:\n%s", out)
	}
}

func TestSelectNodesBatchesAndRefusals(t *testing.T) {
	t.Run("CGPCH-B06: An ungoverned file does not taint a batch", func(t *testing.T) {})
	t.Run("CGPCH-E05: The check with no scope is refused", func(t *testing.T) {})
	dir := t.TempDir()
	for _, f := range []string{"src/hooks/a.ts", "package.json", "yarn.lock"} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, f)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, f), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "src/hooks/a.ts", Kind: mapx.KindCode}}}
	cfg := governanceConfig()

	if _, _, err := selectNodes(g, cfg, false, nil, dir); err == nil || !strings.Contains(err.Error(), "provide --changed") {
		t.Errorf("neither --changed nor --all must be refused, got %v", err)
	}
	nodes, scope, err := selectNodes(g, cfg, false, []string{"src/hooks/a.ts", "package.json"}, dir)
	if err != nil || len(nodes) != 1 || scope != "--changed (2 files)" {
		t.Errorf("an ungoverned file does not taint the batch: nodes=%v scope=%q err=%v", nodes, scope, err)
	}
	_, _, err = selectNodes(g, cfg, false, []string{"package.json", "yarn.lock"}, dir)
	var ng errNotGoverned
	if !errors.As(err, &ng) || !strings.Contains(err.Error(), "package.json (and 1 more)") {
		t.Errorf("an all-ungoverned batch names one file and counts the rest, got %v", err)
	}
	// a lone ungoverned file has no others to count
	_, _, err = selectNodes(g, cfg, false, []string{"package.json"}, dir)
	if !errors.As(err, &ng) || ng.target != "package.json" {
		t.Errorf("a lone ungoverned file is named alone, got %v", err)
	}
}

func TestCheckReadsTheWaiverFromTheEnvironment(t *testing.T) {
	t.Run("CGPCH-B13: A waiver in the environment drops the gate and says why", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	t.Setenv("ANCHORS_SKIP_RULES", "code-flagged=the linter is broken upstream")
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "○ waived: code-flagged — the linter is broken upstream") || !strings.Contains(out, "3 nodes, 3 gates") {
		t.Errorf("the environment waiver must drop the gate and say why:\n%s", out)
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	return testkit.CaptureStderr(t, fn)
}

func TestCheckWarnsAboutRevStaleMapNodesOnChangedOnly(t *testing.T) {
	t.Run("CGPCH-B32: Nodes edited after the map build are warned about on the incremental check", func(t *testing.T) {})
	t.Run("CGPCH-X01: The stale-map warning does not bar the check", func(t *testing.T) {})
	englishOutput(t)
	g := checkGraph()
	g.Nodes[0].Rev = "a-revision-from-before-the-edit"
	dir := qProject(t, checkYAML, checkFiles(), g)
	var runErr error
	errOut := captureStderr(t, func() {
		_, runErr = runQ(t, newCheckCmd(), "--root", dir, "--changed", "a.go", "--no-record")
	})
	if runErr != nil {
		t.Fatal(runErr)
	}
	if !strings.Contains(errOut, "the map is OLDER than") || !strings.Contains(errOut, "a.go") {
		t.Errorf("a node whose content moved since the map build must be warned about:\n%s", errOut)
	}
	errOut = captureStderr(t, func() {
		_, runErr = runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	})
	if runErr != nil {
		t.Fatal(runErr)
	}
	if strings.Contains(errOut, "the map is OLDER than") {
		t.Errorf("the full sweep does not warn about rev-stale nodes:\n%s", errOut)
	}

	// a map built after the last edit has nothing to warn about, even on --changed
	fresh := qProject(t, checkYAML, checkFiles(), nil)
	cfg, err := config.Load(filepath.Join(fresh, config.DefaultFile))
	if err != nil {
		t.Fatal(err)
	}
	files, err := scan.Walk(fresh, cfg)
	if err != nil {
		t.Fatal(err)
	}
	rev := map[string]string{}
	for _, f := range files {
		rev[f.Path] = f.Rev
	}
	built := checkGraph()
	for i := range built.Nodes {
		built.Nodes[i].Rev = rev[built.Nodes[i].ID]
	}
	if err := mapx.Save(built, filepath.Join(fresh, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	errOut = captureStderr(t, func() {
		_, runErr = runQ(t, newCheckCmd(), "--root", fresh, "--changed", "a.go", "--no-record")
	})
	if runErr != nil {
		t.Fatal(runErr)
	}
	if strings.Contains(errOut, "the map is OLDER than") {
		t.Errorf("a map built after the last edit must not be warned about:\n%s", errOut)
	}
}

func TestPrintProfileFooterSaysWhatIsStillOpen(t *testing.T) {
	t.Run("CGPCH-B58: The verdict line says what is still open", func(t *testing.T) {})
	englishOutput(t)
	footer := func(p gate.Profile) string { return captureStdout(t, func() { printProfile(p, false, false) }) }
	info := gate.Result{Gate: "g", Target: "a", Verdict: gate.Fail}
	cases := []struct {
		name string
		p    gate.Profile
		want string
	}{
		{"informative failure", gate.Profile{Passed: true, Results: []gate.Result{info}, Failures: []gate.Result{info}}, "INFORMATIVE finding(s) open"},
		{"divergence", gate.Profile{Passed: true, Results: []gate.Result{{Gate: "g", Target: "a", Verdict: gate.Pending, Detail: "x"}}}, "1 PENDING item(s) open"},
		{"not confronted", gate.Profile{Passed: true, Results: []gate.Result{{Gate: "g", Target: "a", Verdict: gate.Skip}}}, "1 confrontation(s) DID NOT HAPPEN"},
		{"clean", gate.Profile{Passed: true, Results: []gate.Result{{Gate: "g", Target: "a", Verdict: gate.Pass}}}, "all gates passed, no open findings"},
		{"blocked", gate.Profile{Blocked: []gate.Result{{Gate: "g"}}, Results: []gate.Result{info}, Failures: []gate.Result{info}}, "✗ blocked — 1 blocking gate(s) failed (+1 informative finding(s))"},
	}
	for _, c := range cases {
		if out := footer(c.p); !strings.Contains(out, c.want) {
			t.Errorf("%s: the footer must say %q:\n%s", c.name, c.want, out)
		}
	}

	// the pending items ride along with the informative findings only when there are any
	drift := gate.Result{Gate: "g", Target: "b", Verdict: gate.Pending, Detail: "x"}
	withDrift := footer(gate.Profile{Passed: true, Results: []gate.Result{info, drift}, Failures: []gate.Result{info}})
	if !strings.Contains(withDrift, "(+1 pending item(s)") {
		t.Errorf("informative findings with a pending item must count it:\n%s", withDrift)
	}
	if out := footer(gate.Profile{Passed: true, Results: []gate.Result{info}, Failures: []gate.Result{info}}); strings.Contains(out, "pending item(s) —") {
		t.Errorf("informative findings with no pending item must not mention pending items:\n%s", out)
	}
	// a failed check with no informative finding says only the blocking count
	block := gate.Result{Gate: "g", Target: "a", Verdict: gate.Fail, Blocking: true}
	if out := footer(gate.Profile{Blocked: []gate.Result{block}, Results: []gate.Result{block}, Failures: []gate.Result{block}}); !strings.Contains(out, "✗ blocked — 1 blocking gate(s) failed\n") {
		t.Errorf("a failed check with no informative finding must not add a count of them:\n%s", out)
	}
}

func TestSkipReasonsShownOnASmallScan(t *testing.T) {
	t.Run("CGPCH-B51: A small scan lists the reason of each skip", func(t *testing.T) {})
	englishOutput(t)
	p := profileOf(gate.GateSummary{Gate: "g", Pass: 1, Skip: 1})
	p.Results = []gate.Result{{Gate: "g", Verdict: gate.Skip, Detail: "not a screen", Target: "n.ts"}}
	out := captureStdout(t, func() { printProfile(p, false, false) })
	if !strings.Contains(out, "~ 1 indeterminate — not a failure") || !strings.Contains(out, "  ~ g @ n.ts\n      not a screen\n") {
		t.Errorf("a small scan lists the reason of each ~:\n%s", out)
	}

	// forty results is still "at most forty"
	for len(p.Results) < maxResultsForSkipDetail {
		p.Results = append(p.Results, gate.Result{Gate: "g", Verdict: gate.Pass, Target: "p.ts"})
	}
	out = captureStdout(t, func() { printProfile(p, false, false) })
	if !strings.Contains(out, "  ~ g @ n.ts\n      not a screen\n") {
		t.Errorf("a scan of exactly forty results lists the reasons:\n%s", out)
	}

	// no skip with a reason, nothing to list
	clean := profileOf(gate.GateSummary{Gate: "g", Pass: 1})
	clean.Results = []gate.Result{{Gate: "g", Verdict: gate.Pass, Target: "n.ts"}}
	if out := captureStdout(t, func() { printProfile(clean, false, false) }); strings.Contains(out, "indeterminate — not a failure") {
		t.Errorf("with no reason to list, no indeterminate block:\n%s", out)
	}
}

// ── the command end to end ──────────────────────────────────────────────────────────

// Gates are informative on purpose: a blocking fail ends `check` in os.Exit(1), which an
// in-process test cannot survive. The verdicts are still reported and recorded.
const checkYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
  spec:
    kind: spec
    pattern: "*.spec.md"
gates:
  - name: code-ok
    on: [code]
    run: "true"
  - name: code-flagged
    on: [code]
    run: "echo flagged by the tool; exit 1"
  - name: spec-only
    on: [spec]
    when: [pre-push]
    run: "true"
  - name: never-applies
    on: [feature]
    run: "true"
`

func checkGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "a.go", Kind: mapx.KindCode, Rev: "r1"},
			{ID: "b.go", Kind: mapx.KindCode, Rev: "r2"},
			{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "s1"},
		},
		Edges: []mapx.Edge{{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies}},
	}
}

func checkFiles() map[string]string {
	return map[string]string{"a.go": "package a\n", "b.go": "package a\n", "a.spec.md": "# A\n"}
}

func TestCheckAllReportsWithoutRecording(t *testing.T) {
	t.Run("CGPCH-B01: The full sweep confronts every node of the map", func(t *testing.T) {})
	t.Run("CGPCH-B17: The no-record mode leaves the map untouched", func(t *testing.T) {})
	t.Run("CGPCH-B35: A declared gate with nothing to measure is named", func(t *testing.T) {})
	t.Run("CGPCH-B38: The report is mirrored to a file", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	before := readQ(t, filepath.Join(dir, mapx.DefaultPath))

	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"check --all — 3 nodes, 4 gates",
		"code-ok",
		"code-flagged",
		// a declared gate that no node reached is named, not silently dropped
		"    never-applies\n",
		"output mirrored to .anchors/",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if after := readQ(t, filepath.Join(dir, mapx.DefaultPath)); after != before {
		t.Errorf("--no-record must not stamp the map")
	}
}

// The incremental check confronts the impact path of the change and stamps the edges it
// confronted.
func TestCheckChangedStampsTheImpactPath(t *testing.T) {
	t.Run("CGPCH-B02: The incremental check confronts only the impact path of the change", func(t *testing.T) {})
	t.Run("CGPCH-B16: The confronted edges are stamped at the current revisions", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())

	out, err := runQ(t, newCheckCmd(), "--root", dir, "--changed", filepath.Join(dir, "a.spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "2 nodes, 4 gates") {
		t.Errorf("the impact path of the spec is the spec and its code:\n%s", out)
	}
	if strings.Contains(out, "b.go") {
		t.Errorf("a node off the impact path is not confronted:\n%s", out)
	}
	g, err := mapx.Load(filepath.Join(dir, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if st := g.Edges[0].Stamp; st == nil || st.ValidatedFromRev != "s1" || st.ValidatedToRev != "r1" {
		t.Errorf("the confronted edge must be stamped at the current revs, got %+v", st)
	}
}

func TestCheckFiltersAndRefusals(t *testing.T) {
	t.Run("CGPCH-B07: Phase and category select the gates charged", func(t *testing.T) {})
	t.Run("CGPCH-B14: The deterministic mode drops the judgment gates", func(t *testing.T) {})
	t.Run("CGPCH-E01: The check without configuration fails", func(t *testing.T) {})
	t.Run("CGPCH-E02: A project with no gate has no pipeline", func(t *testing.T) {})
	t.Run("CGPCH-E03: The check without a map points at the map build", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())

	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--category", "nothing-declares-this")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `no gate to run for this slice (phase="" category="nothing-declares-this")`) {
		t.Errorf("an empty slice says so and runs nothing:\n%s", out)
	}

	// a gate declared for pre-push is not charged at pre-commit
	out, err = runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--phase", "pre-commit")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "spec-only") || !strings.Contains(out, "3 gates") {
		t.Errorf("the pre-push gate must be out of a pre-commit run:\n%s", out)
	}

	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--skip-rule", "code-ok"); err == nil || !strings.Contains(err.Error(), "invalid --skip-rule") {
		t.Errorf("a waiver without a reason must be refused; got %v", err)
	}

	judgeOnly := qProject(t, "version: 2\nlayers: {}\ngates:\n  - name: j\n    on: [code]\n    measures: judgment\n", checkFiles(), checkGraph())
	out, err = runQ(t, newCheckCmd(), "--root", judgeOnly, "--all", "--deterministic")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no deterministic gate to run") {
		t.Errorf("--deterministic with only judgment gates runs nothing:\n%s", out)
	}

	noGates := qProject(t, "version: 2\nlayers: {}\n", checkFiles(), checkGraph())
	if _, err := runQ(t, newCheckCmd(), "--root", noGates, "--all"); err == nil || !strings.Contains(err.Error(), "no gate declared") {
		t.Errorf("a project without gates has no pipeline; got %v", err)
	}
	noMap := qProject(t, checkYAML, checkFiles(), nil)
	if _, err := runQ(t, newCheckCmd(), "--root", noMap, "--all"); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map; got %v", err)
	}
	if _, err := runQ(t, newCheckCmd(), "--root", t.TempDir(), "--all"); err == nil || !strings.Contains(err.Error(), "load config") {
		t.Errorf("no config; got %v", err)
	}
}

// The commit message waives a rule when it carries the reason, and is refused without it.
func TestCheckReadsTheWaiverFromTheCommitMessage(t *testing.T) {
	t.Run("CGPCH-B12: A commit message marker with a reason waives a gate", func(t *testing.T) {})
	t.Run("CGPCH-E04: A waiver without a reason is refused", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	msg := filepath.Join(dir, "MSG")
	write := func(body string) {
		if err := os.WriteFile(msg, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("feat: x\n\n[skip-code-flagged: the linter is broken upstream]\n")
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--commit-msg", msg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "flagged by the tool") || !strings.Contains(out, "3 nodes, 3 gates") {
		t.Errorf("the waived gate must not run:\n%s", out)
	}

	write("feat: x\n\n[skip-code-flagged: ]\n")
	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--commit-msg", msg); err == nil || !strings.Contains(err.Error(), "invalid --skip-rule") {
		t.Errorf("a marker without a reason must be refused; got %v", err)
	}
}

func TestNormalizeChanged(t *testing.T) {
	t.Run("CGPCH-B04: Changed paths are normalised to the map's form", func(t *testing.T) {})
	root := filepath.FromSlash("/repo")
	got := normalizeChanged([]string{filepath.Join(root, "src", "a.go"), "./b.go", "c/../d.go"}, root)
	if strings.Join(got, ",") != "src/a.go,b.go,d.go" {
		t.Errorf("normalizeChanged = %v", got)
	}
}

// ── scope: which nodes the check confronts ─────────────────────────────────────────

func governanceConfig() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{
		"hook":    {Pattern: "src/hooks/**/*.ts", Kind: "code"},
		"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
		"feature": {Pattern: "**/*.feature", Kind: "feature"},
		"test":    {Pattern: "**/*.test.ts", Kind: "test"},
		"doc":     {Pattern: "**/*.md", Kind: "doc"},
	}}
}

// The HOLE this test locks: `selectNodes` answered the SAME thing for two opposite
// situations — "a new governed file, still outside the map" and "a file Anchors does not
// govern at all" — and the pre-commit, unable to tell them apart, treated both as benign.
// The result was the worst case: a NEW hook/screen/service committed clean with no spec,
// no feature and no test, because outside the map no gate confronts it.
func TestSelectNodesTellsGovernedFromUngoverned(t *testing.T) {
	t.Run("CGPCH-B05: Files the project does not govern are recognised as not governed", func(t *testing.T) {})
	t.Run("CGPCH-E06: A governed file outside the map bars the check", func(t *testing.T) {})
	t.Run("CGPCH-E07: A path on neither disk nor map is an error", func(t *testing.T) {})
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src", "hooks"), 0o755)

	governed := "src/hooks/useNew.ts"
	ungoverned := "package.json"
	os.WriteFile(filepath.Join(dir, governed), []byte("export const x = 1\n"), 0o644)
	os.WriteFile(filepath.Join(dir, ungoverned), []byte("{}\n"), 0o644)

	// An EMPTY map: the real condition right after creating a file, before `map build`.
	g := &mapx.Graph{}
	cfg := governanceConfig()

	t.Run("a governed file outside the map is an OFFENCE (it bars)", func(t *testing.T) {
		_, _, err := selectNodes(g, cfg, false, []string{governed}, dir)
		if err == nil {
			t.Fatal("passed with no error — a governed file outside the map must bar the commit")
		}
		var nr errNotGoverned
		if errors.As(err, &nr) {
			t.Fatalf("classified as UNGOVERNED (it would exit %d and the hook would continue): %v", common.ExitNotGoverned, err)
		}
		if !strings.Contains(err.Error(), "GOVERNED") {
			t.Fatalf("the message does not say the file is governed: %v", err)
		}
	})

	t.Run("an ungoverned file exits with its own code (benign)", func(t *testing.T) {
		_, _, err := selectNodes(g, cfg, false, []string{ungoverned}, dir)
		if err == nil {
			t.Fatal("expected the ungoverned signal, got nil")
		}
		var nr errNotGoverned
		if !errors.As(err, &nr) {
			t.Fatalf("it did not signal ungoverned — the hook would bar package.json: %v", err)
		}
	})

	t.Run("a missing path is not ungoverned", func(t *testing.T) {
		// A typo in the path cannot become "benign" and vanish: it would exit 3 and the
		// pre-commit would continue, silencing the mistake.
		_, _, err := selectNodes(g, cfg, false, []string{"src/hooks/doesNotExist.ts"}, dir)
		if err == nil {
			t.Fatal("expected an error for a missing path")
		}
		var nr errNotGoverned
		if errors.As(err, &nr) {
			t.Fatalf("a missing path was classified as ungoverned: %v", err)
		}
	})

	t.Run("the Anchors record (issues/) is not governed", func(t *testing.T) {
		// `issues/` matches the `doc` layer (`**/*.md`) but the scanner NEVER indexes it:
		// it is Anchors' own OUTPUT. Without consulting the ignore, the path became
		// "governed outside the map" — and `map build` never added it, so the commit was
		// barred forever. Measured when committing the issues `check` had resolved.
		os.MkdirAll(filepath.Join(dir, "issues", "done"), 0o755)
		iss := "issues/done/2026-08-15--violation--x.md"
		os.WriteFile(filepath.Join(dir, iss), []byte("# issue\n"), 0o644)
		_, _, err := selectNodes(g, cfg, false, []string{iss}, dir)
		var nr errNotGoverned
		if !errors.As(err, &nr) {
			t.Fatalf("issues/ should be ungoverned (exit %d), got: %v", common.ExitNotGoverned, err)
		}
	})
}

// The `-progress.md` is the same deadlock `issues/`/`changes/` already solved, through
// another door: it matches the `plan` layer (`plans/*.md`) and the scanner NEVER indexes
// it — on purpose, because a file that exists to CHANGE cannot be confronted by gates that
// charge a justification for change.
//
// Without the exclusion the result was the worst case: the file said "governed", absent
// from the map, and `map build` never adding it — the commit barred forever by the very
// mechanism that separated decision from state.
//
// Measured in the reference app when committing the 17 progress files `anchors new progress` had
// just created.
func TestSelectNodes_progressIsNotGoverned(t *testing.T) {
	t.Run("CGPCH-B67: A plan's progress companion is not governed", func(t *testing.T) {})
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "plans"), 0o755)

	plan := "plans/0002-platform.md"
	progress := "plans/0002-platform-progress.md"
	os.WriteFile(filepath.Join(dir, plan), []byte("# Plan\n"), 0o644)
	os.WriteFile(filepath.Join(dir, progress), []byte("# Progress\n\n- [x] done\n"), 0o644)

	cfg := governanceConfig()
	cfg.Layers["plan"] = config.Layer{Pattern: "plans/*.md", Kind: "plan"}
	g := &mapx.Graph{} // an empty map: the real condition right after creating the file

	_, _, err := selectNodes(g, cfg, false, []string{progress}, dir)
	if err == nil {
		t.Fatal("passed with no error — errNotGoverned was expected, not silence")
	}
	var nr errNotGoverned
	if !errors.As(err, &nr) {
		t.Fatalf("the progress file was treated as GOVERNED and the commit would be barred "+
			"forever (`map build` never adds it): %v", err)
	}

	// and the PLAN stays governed: the exclusion is the companion's, not the layer's.
	_, _, err = selectNodes(g, cfg, false, []string{plan}, dir)
	if err == nil {
		t.Fatal("the plan outside the map should bar")
	}
	if errors.As(err, &nr) {
		t.Fatal("the plan was classified as ungoverned — the exclusion leaked into the layer")
	}
}

// `check --changed <module>` brings in the tests whose `@contract` stamps point at the
// module: the pre-commit of whoever changes a function then runs `mock-stamped` on the
// doubles of that function, and a stale double blocks the change that made it stale.
func TestImpactOfBringsTheTestsThatStampTheChangedFile(t *testing.T) {
	t.Run("CGPCH-B03: The tests that stamp a changed module enter its impact path", func(t *testing.T) {})
	root := t.TempDir()
	files := map[string]string{
		"src/hooks/balance.ts":       "export function useBalance(id) {\n  return id\n}\n",
		"src/screens/Home.test.tsx":  "// @contract: src/hooks/balance.ts | export function useBalance(id) { | 3 | deadbeef\njest.mock('@/src/hooks/balance')\n",
		"src/screens/Other.test.tsx": "it('x', () => {})\n",
	}
	g := &mapx.Graph{}
	for name, body := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		kind := mapx.KindCode
		if filepath.Ext(filepath.Base(name)) == ".tsx" {
			kind = mapx.KindTest
		}
		g.Nodes = append(g.Nodes, mapx.Node{ID: name, Kind: kind})
	}

	ids, err := impactOf(g, &config.Config{}, filepath.Join(root, "src/hooks/balance.ts"), root)
	if err != nil {
		t.Fatal(err)
	}
	has := map[string]bool{}
	for _, id := range ids {
		has[id] = true
	}
	if !has["src/screens/Home.test.tsx"] {
		t.Errorf("the test stamping the changed module is not in its impact: %v", ids)
	}
	if has["src/screens/Other.test.tsx"] {
		t.Errorf("a test with no stamp on the module entered its impact: %v", ids)
	}
}

// ── gate selection ──────────────────────────────────────────────────────────────────

func names(gs []config.Gate) []string {
	out := make([]string, 0, len(gs))
	for _, g := range gs {
		out = append(out, g.Name)
	}
	return out
}

// `skip_on` is an EXCLUSION list: a gate that declares nothing runs in both perspectives.
// Were it an inclusion list, every gate written before this field existed would stop
// running — silence can never switch verification off.
func TestFilterGatesSkipOnIsPermissiveByDefault(t *testing.T) {
	t.Run("CGPCH-B08: A gate without skip_on runs in both perspectives", func(t *testing.T) {})
	gates := []config.Gate{{Name: "undeclared"}}

	for _, p := range []string{config.PerspectiveChange, config.PerspectiveAll} {
		got := filterGates(gates, "", "", false, p, gate.Waiver{})
		if len(got) != 1 {
			t.Errorf("perspective %q: a gate without `skip_on` must run, got %v", p, names(got))
		}
	}
}

// The use case that motivated the field: a gate that only answers well over the whole
// picture (an orphan detector asked about a slice would accuse everything the slice does
// not reach) declares itself out of `--changed`.
func TestFilterGatesSkipOnChange(t *testing.T) {
	t.Run("CGPCH-B09: A gate that skips the change perspective runs only on the full sweep", func(t *testing.T) {})
	gates := []config.Gate{
		{Name: "all-only", SkipOn: []string{config.PerspectiveChange}},
		{Name: "always"},
	}

	onChange := filterGates(gates, "", "", false, config.PerspectiveChange, gate.Waiver{})
	if len(onChange) != 1 || onChange[0].Name != "always" {
		t.Errorf("`skip_on: [change]` must leave --changed, got %v", names(onChange))
	}

	onAll := filterGates(gates, "", "", false, config.PerspectiveAll, gate.Waiver{})
	if len(onAll) != 2 {
		t.Errorf("`skip_on: [change]` must stay in --all, got %v", names(onAll))
	}
}

// The reverse: a gate too expensive for the full sweep leaves `--all` without leaving the
// commit. That is what separates this axis from `cost: slow`, which removes the gate from
// both.
func TestFilterGatesSkipOnAll(t *testing.T) {
	t.Run("CGPCH-B10: A gate that skips the full sweep runs only on changed files", func(t *testing.T) {})
	gates := []config.Gate{{Name: "slice-only", SkipOn: []string{config.PerspectiveAll}}}

	if got := filterGates(gates, "", "", false, config.PerspectiveAll, gate.Waiver{}); len(got) != 0 {
		t.Errorf("`skip_on: [all]` must leave --all, got %v", names(got))
	}
	if got := filterGates(gates, "", "", false, config.PerspectiveChange, gate.Waiver{}); len(got) != 1 {
		t.Errorf("`skip_on: [all]` must stay in --changed, got %v", names(got))
	}
}

// Both together switch the gate off — an explicit declaration, not an accident, so the
// filter honours it instead of treating it as a contradiction.
func TestFilterGatesSkipOnBothSwitchesOff(t *testing.T) {
	t.Run("CGPCH-B11: A gate that skips both perspectives is switched off", func(t *testing.T) {})
	gates := []config.Gate{{
		Name:   "off",
		SkipOn: []string{config.PerspectiveChange, config.PerspectiveAll},
	}}

	for _, p := range []string{config.PerspectiveChange, config.PerspectiveAll} {
		if got := filterGates(gates, "", "", false, p, gate.Waiver{}); len(got) != 0 {
			t.Errorf("perspective %q: it should be off, got %v", p, names(got))
		}
	}
}

// The axes are INDEPENDENT: declaring a perspective cannot change the effect of
// phase/category/cost, or the categorisation would stop being adoptable bit by bit.
func TestFilterGatesAxesAreIndependent(t *testing.T) {
	t.Run("CGPCH-I01: Declaring a perspective does not change the cost axis", func(t *testing.T) {})
	gates := []config.Gate{
		{Name: "slow", Cost: "slow", SkipOn: []string{config.PerspectiveChange}},
		{Name: "fast"},
	}

	// on --all the slow one runs; with --skip-slow it leaves by COST, not by perspective.
	if got := filterGates(gates, "", "", false, config.PerspectiveAll, gate.Waiver{}); len(got) != 2 {
		t.Errorf("without skip-slow both run on --all, got %v", names(got))
	}
	if got := filterGates(gates, "", "", true, config.PerspectiveAll, gate.Waiver{}); len(got) != 1 {
		t.Errorf("skip-slow must remove the slow one, got %v", names(got))
	}
}

// ── recording: stamps and issues ────────────────────────────────────────────────────

// Which check writes issues, per mode: local always; github only in CI or when asked;
// manual only when asked.
func TestIssuesOnFor(t *testing.T) {
	t.Run("CGPCH-B15: The issue policy follows the workflow mode", func(t *testing.T) {})
	manual := &config.Config{Workflow: &config.Workflow{Mode: config.ModeManual}}
	gh := &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub, Repo: "o/r", Labels: []string{"anchors"}}}
	for _, c := range []struct {
		name         string
		cfg          *config.Config
		record, inCI bool
		want         bool
	}{
		{"local", &config.Config{}, false, false, true},
		{"manual", manual, false, false, false},
		{"manual in CI", manual, false, true, false},
		{"manual --record-issues", manual, true, false, true},
		{"github local run", gh, false, false, false},
		{"github in CI", gh, false, true, true},
		{"github --record-issues", gh, true, false, true},
	} {
		if got := issuesOnFor(c.cfg, c.record, c.inCI); got != c.want {
			t.Errorf("%s: issuesOnFor = %v, want %v", c.name, got, c.want)
		}
	}
}

// A LOCAL check stamps the map but opens no issue: it runs on work in progress, and in
// the reference app an agent's local check filed two `[docs-fresh]` cards for a state that existed
// only on its machine. With issues on (CI, or --record-issues) the failure is filed.
func TestRecordCheck_issuesOnlyWhenOn(t *testing.T) {
	t.Run("CGPCH-B18: A blocking failure is filed only when issues are on", func(t *testing.T) {})
	t.Run("CGPCH-B77: The record summary counts what the record did", func(t *testing.T) {})
	issue.UseFiles()
	fail := gate.Profile{Results: []gate.Result{{
		Gate: "docs-fresh", Target: "a.spec.md", Verdict: gate.Fail, Blocking: true, Detail: "stale",
	}}}
	englishOutput(t)
	var out string
	run := func(on bool) (string, int) {
		root := t.TempDir()
		mapPath := filepath.Join(root, mapx.DefaultPath)
		g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "r1"}}}
		out = captureStdout(t, func() {
			if err := recordCheck(root, mapPath, g, fail, on, false); err != nil {
				t.Fatal(err)
			}
		})
		if _, err := os.Stat(mapPath); err != nil {
			t.Fatalf("the map was not saved (on=%v): %v", on, err)
		}
		n := 0
		_ = filepath.Walk(filepath.Join(root, issue.Dir), func(p string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				n++
			}
			return nil
		})
		return root, n
	}
	if _, n := run(false); n != 0 {
		t.Errorf("a local check (issues off) opened %d issue(s)", n)
	}
	if _, n := run(true); n == 0 {
		t.Error("with issues on, the blocking failure opened no issue")
	}
	// the opened issue is counted, and a record with no debt says nothing about debts
	if !strings.Contains(out, "; 1 new issue(s), 0 resolved") || strings.Contains(out, "assumed debt") {
		t.Errorf("the summary must count the opened issue and no debt:\n%s", out)
	}
}

// The FULL check closes the open violations it did not reproduce; a partial check does not
// (it did not confront everything, so it cannot say a violation is gone).
func TestRecordCheck_fullCheckClosesWhatItDidNotReproduce(t *testing.T) {
	t.Run("CGPCH-B20: The full check closes the violations it did not reproduce", func(t *testing.T) {})
	t.Run("CGPCH-B77: The record summary counts what the record did", func(t *testing.T) {})
	issue.UseFiles()
	englishOutput(t)
	pass := gate.Profile{Results: []gate.Result{{Gate: "header-conforms", Target: "a.ts", Verdict: gate.Pass}}}
	var out string
	run := func(full bool) issue.State {
		root := t.TempDir()
		mapPath := filepath.Join(root, mapx.DefaultPath)
		g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.ts", Kind: mapx.KindCode, Rev: "r1"}}}
		old := issue.Issue{Kind: issue.Violation, Gate: "header-conforme", Target: "a.ts", Date: "2026-09-11"}
		if _, _, err := issue.Open(root, old); err != nil {
			t.Fatal(err)
		}
		out = captureStdout(t, func() {
			if err := recordCheck(root, mapPath, g, pass, true, full); err != nil {
				t.Fatal(err)
			}
		})
		st, _ := issue.Exists(root, old.Key())
		return st
	}
	if st := run(true); st != issue.Done {
		t.Errorf("a full check must close the violation it did not reproduce, it is %s", st)
	}
	// the closed violation is announced and counted among the resolved
	if !strings.Contains(out, "1 open violation issue(s) in issues/ not reproduced") || !strings.Contains(out, "; 0 new issue(s), 1 resolved") {
		t.Errorf("the full check must announce and count what it closed:\n%s", out)
	}
	if st := run(false); st != issue.Todo {
		t.Errorf("a partial check must leave it alone, it is %s", st)
	}

	// a full check with nothing left to close does not announce a closing
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.ts", Kind: mapx.KindCode, Rev: "r1"}}}
	out = captureStdout(t, func() {
		if err := recordCheck(root, filepath.Join(root, mapx.DefaultPath), g, pass, true, true); err != nil {
			t.Fatal(err)
		}
	})
	if strings.Contains(out, "not reproduced by this full check") {
		t.Errorf("nothing was closed, and the closing line appeared:\n%s", out)
	}
}

// ── the judgment queue ──────────────────────────────────────────────────────────────

// The queue is persistent and the set of applicable targets is not. A gate that starts
// declaring `requires`, or a target that loses the mark that made it applicable, leaves the
// old task orphaned — and the queue starts lying about the size of the work.
func TestDropStaleJudgments(t *testing.T) {
	t.Run("CGPCH-B23: Stale judge tasks leave the queue", func(t *testing.T) {})
	root := t.TempDir()
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "no-test-proof-real", Measures: config.MeasuresJudgment},
	}}
	// the targets must exist: `queue.List` already drops a task whose target is gone, and
	// what is tested here is the other case — a target that exists and a gate that no
	// longer applies.
	for _, f := range []string{"a.spec.md", "b.spec.md", "c.spec.md"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("# spec\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	live := queue.Task{
		ID: judgeTaskID("no-test-proof-real", "a.spec.md"), Changed: "a.spec.md",
		Kind: "judgment", Origin: "check",
	}
	stale := queue.Task{
		ID: judgeTaskID("no-test-proof-real", "b.spec.md"), Changed: "b.spec.md",
		Kind: "judgment", Origin: "check",
	}
	// work from another origin must NOT be dropped by this path
	foreign := queue.Task{
		ID: judgeTaskID("no-test-proof-real", "c.spec.md"), Changed: "c.spec.md",
		Kind: "judgment", Origin: "human",
	}
	for _, tk := range []queue.Task{live, stale, foreign} {
		if _, err := queue.Enqueue(root, tk); err != nil {
			t.Fatal(err)
		}
	}

	// this round's check enqueued only `a`
	p := gate.Profile{Judged: []gate.Result{
		{Gate: "no-test-proof-real", Target: "a.spec.md"},
	}}
	knownJudgmentGates = []string{"no-test-proof-real"}
	dropStaleJudgments(root, cfg, p)

	left := map[string]bool{}
	tasks, err := queue.List(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range tasks {
		left[tk.ID] = true
	}
	if !left[live.ID] {
		t.Error("the task of a target still applicable must NOT be dropped")
	}
	if left[stale.ID] {
		t.Error("the task of a target the gate no longer enqueues should have left the queue")
	}
	if !left[foreign.ID] {
		t.Error("a task of another origin is not this path's — it should stay")
	}
}

// The ID is `judge-<gate>-<slug>`, and the gate name contains `-`: it is read by
// rebuilding the ID from each known gate and the task's target, not by splitting on the
// separator.
func TestGateOfJudgeTask(t *testing.T) {
	knownJudgmentGates = []string{"no-test-proof-real", "atomic-design", "review", "review-deep"}
	cases := []struct{ id, changed, want string }{
		{"judge-no-test-proof-real-apps-x-y.spec", "apps/x/y.spec.md", "no-test-proof-real"},
		{"judge-atomic-design-apps-x", "apps/x.tsx", "atomic-design"},
		{"judge-review-deep-src-x", "src/x.go", "review-deep"},
		{"judge-review-deep-src-x", "deep/src/x.go", "review"},
		{"judge-gate-that-does-not-exist-x", "x.go", ""},
		{"something-else", "x.go", ""},
	}
	for _, c := range cases {
		if got := gateDaTaskJudge(c.id, c.changed); got != c.want {
			t.Errorf("%s (%s) → %q, want %q", c.id, c.changed, got, c.want)
		}
	}
}

// On `--changed` the cleanup does NOT run — which is why it needs the mode.
//
// `dropStaleJudgments` concludes "not enqueued now, so it is obsolete". The inference holds
// when the check looked at EVERY node; on `--changed` it looked at one file, and the other
// targets of the same gate look obsolete only because they were not looked at.
//
// Measured in the reference app: two `check --changed` on different tests left ONE task in the
// queue, and `judge --pending` answered "no target awaiting" with five pending on
// `check --all`.
//
// The neighbouring test calls `dropStaleJudgments` directly and so never exercised WHEN to
// call it — this one covers that decision.
func TestEnqueueJudgments_incrementalKeepsTheOtherTargets(t *testing.T) {
	t.Run("CGPCH-B24: An incremental check keeps the judgments it did not look at", func(t *testing.T) {})
	root := t.TempDir()
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "no-test-proof-real", Measures: config.MeasuresJudgment},
	}}
	for _, f := range []string{"a.spec.md", "b.spec.md"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("# spec\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// already queued: the judgment of `b`, from an earlier round
	if _, err := queue.Enqueue(root, queue.Task{
		ID: judgeTaskID("no-test-proof-real", "b.spec.md"), Changed: "b.spec.md",
		Kind: "judgment", Origin: "check",
	}); err != nil {
		t.Fatal(err)
	}
	knownJudgmentGates = []string{"no-test-proof-real"}
	// this round is incremental and saw only `a`
	p := gate.Profile{Judged: []gate.Result{
		{Gate: "no-test-proof-real", Target: "a.spec.md"},
	}}

	enqueueJudgments(root, cfg, p, false /* full sweep */)

	tasks, err := queue.List(root)
	if err != nil {
		t.Fatal(err)
	}
	left := map[string]bool{}
	for _, tk := range tasks {
		left[tk.ID] = true
	}
	if !left[judgeTaskID("no-test-proof-real", "b.spec.md")] {
		t.Error("`--changed` erased the judgment of a target it did not even look at")
	}

	// and on the FULL sweep the cleanup holds again: there "not enqueued" really means
	// obsolete, the case `dropStaleJudgments` exists for.
	enqueueJudgments(root, cfg, p, true /* full sweep */)
	tasks, err = queue.List(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range tasks {
		if tk.ID == judgeTaskID("no-test-proof-real", "b.spec.md") {
			t.Error("`--all` should have dropped the target the gate no longer enqueues")
		}
	}
}

// ── the judgment brief ─────────────────────────────────────────────────────────────

// A COUNT IS NOT AN ADDRESS.
//
// The earlier version printed only "%d targets awaiting judgment", and the reviewer knew
// there was work without knowing which, where, or what to ask. Measured in the reference
// app: 59 of 85 specs never received a verdict, with CI running on every PR.
func judgedFixture() []gate.Result {
	return []gate.Result{
		{Gate: "rule-fulfilled", Target: "a/One.spec.md", Verdict: gate.Judge},
		{Gate: "rule-fulfilled", Target: "b/Two.spec.md", Verdict: gate.Judge},
		{Gate: "test-proves", Target: "c/Three.test.ts", Verdict: gate.Judge},
	}
}

func brief(t *testing.T) string {
	t.Helper()
	englishOutput(t)
	return captureStdout(t, func() {
		printJudgmentBrief(judgedFixture(),
			map[string]string{"rule-fulfilled": "SPEC.md"},
			map[string]string{"rule-fulfilled": "Does the excerpt DO what the rule describes?"})
	})
}

// The TARGET: without it the reviewer does not know where to look.
func TestJudgmentBrief_namesEveryTarget(t *testing.T) {
	t.Run("CGPCH-B26: The judgment brief names every target", func(t *testing.T) {})
	out := brief(t)
	for _, target := range []string{"a/One.spec.md", "b/Two.spec.md", "c/Three.test.ts"} {
		if !strings.Contains(out, target) {
			t.Errorf("the brief does not name the target %q:\n%s", target, out)
		}
	}
}

// The QUESTION: it is the `ask` the gate declares, and without it the reviewer invents
// the criterion.
func TestJudgmentBrief_carriesTheQuestionAndTheGuide(t *testing.T) {
	t.Run("CGPCH-B27: The judgment brief carries the question and the guide", func(t *testing.T) {})
	out := brief(t)
	if !strings.Contains(out, "DO what the rule describes") {
		t.Errorf("the brief does not carry the gate's question:\n%s", out)
	}
	if !strings.Contains(out, "SPEC.md") {
		t.Errorf("the brief does not say where the ruler is:\n%s", out)
	}
}

// GROUPED BY GATE: a gate asks the same thing of several targets, and repeating the
// question per target would make the reviewer read it three times to answer once.
func TestJudgmentBrief_groupsByGate(t *testing.T) {
	t.Run("CGPCH-B28: The judgment brief groups the targets by gate", func(t *testing.T) {})
	out := brief(t)
	if n := strings.Count(out, "DO what the rule describes"); n != 1 {
		t.Errorf("the question appears %d times, want 1 (grouped by gate):\n%s", n, out)
	}
	if !strings.Contains(out, "rule-fulfilled") || !strings.Contains(out, "test-proves") {
		t.Errorf("the brief does not name both gates:\n%s", out)
	}
}

// WHO JUDGES: the pipeline hands over the list, the reviewer decides. Without saying so,
// the list reads as the machine's backlog — and nobody works through it.
func TestJudgmentBrief_saysThePipelineDoesNotJudge(t *testing.T) {
	t.Run("CGPCH-B29: The judgment brief says who judges", func(t *testing.T) {})
	out := brief(t)
	flat := strings.Join(strings.Fields(out), " ")
	if !strings.Contains(flat, "NOT judge") {
		t.Errorf("the brief does not say the pipeline does not judge:\n%s", out)
	}
	if !strings.Contains(out, "REV-CK5") {
		t.Errorf("the brief does not point at the review checklist item:\n%s", out)
	}
}

// UP TO TEN, and the rest counted: sixty paths drown the question that comes before them.
func TestJudgmentBrief_truncatesALongList(t *testing.T) {
	t.Run("CGPCH-B30: The judgment brief lists ten targets per gate and counts the rest", func(t *testing.T) {})
	englishOutput(t)
	var many []gate.Result
	for i := 0; i < 25; i++ {
		many = append(many, gate.Result{Gate: "rule-fulfilled",
			Target: "u/" + string(rune('a'+i)) + ".spec.md", Verdict: gate.Judge})
	}
	out := captureStdout(t, func() {
		printJudgmentBrief(many, map[string]string{}, map[string]string{})
	})
	if !strings.Contains(out, "… and 15 more") {
		t.Errorf("the brief does not count the ones left out:\n%s", out)
	}
	if strings.Count(out, ".spec.md") > 11 {
		t.Errorf("the brief listed more than ten targets:\n%s", out)
	}
}

// ── warnings that do not bar ────────────────────────────────────────────────────────

// THE WARNING IS ABOUT CONTENT, NOT ABOUT THE CLOCK.
//
// The first version compared the map's mtime with the source files'. Cheap and wrong in
// CI: `checkout` writes the WHOLE repository at clone time, and every file ends up
// microseconds newer than the map. Measured on the first run of the gates pipeline — 26
// files "changed" in a repository where nothing had changed, and the same job's log said,
// two lines above, that the map matched the repository.
//
// Noise like that costs more than it seems: a warning on every green PR teaches people to
// ignore it, and then it is useless when it is true.
func TestMapStaleDoesNotDependOnMtime(t *testing.T) {
	t.Run("CGPCH-B31: Governed files missing from the map make a stale-map warning", func(t *testing.T) {})
	dir := t.TempDir()
	spec := filepath.Join(dir, "a.spec.md")
	if err := os.WriteFile(spec, []byte("# A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mapPath := filepath.Join(dir, "anchors.graph.yaml")
	if err := os.WriteFile(mapPath, []byte("version: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The map OLDER than the file — the situation the old heuristic accused.
	old := mustTime(t, "2020-01-01T00:00:00Z")
	if err := os.Chtimes(mapPath, old, old); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Layers: map[string]config.Layer{
		"spec": {Kind: "spec", Pattern: "*.spec.md"},
	}}
	// The file IS in the map: even with an "older" map there is nothing to warn about.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md"}}}

	out := capturaSaida(t, func() { warnIfMapStale(dir, mapPath, cfg, g) })
	if strings.Contains(out, "STALE") {
		t.Errorf("the map knows the file — the mtime cannot raise a warning.\ngot: %s", out)
	}

	// Now the REAL case: a governed file the map does not know.
	if err := os.WriteFile(filepath.Join(dir, "b.spec.md"), []byte("# B\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out = capturaSaida(t, func() { warnIfMapStale(dir, mapPath, cfg, g) })
	if !strings.Contains(out, "STALE") {
		t.Error("a governed file outside the map is INVISIBLE to the check — it must warn")
	}
	// The message says what was measured: saying "changed" would send the reader looking
	// for an edit that does not exist — the common case is a NEW file, never in the map.
	if strings.Contains(out, "changed") {
		t.Errorf("the message must say 'not in the map', not 'changed'.\ngot: %s", out)
	}
	// It counts the files and names one: alone, with nothing else to count.
	if !strings.HasPrefix(out, "⚠ STALE map: 1 governed file(s) are not in the map — e.g.: b.spec.md\n") {
		t.Errorf("one missing file must be counted and named alone.\ngot: %s", out)
	}
	if err := os.WriteFile(filepath.Join(dir, "c.spec.md"), []byte("# C\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out = capturaSaida(t, func() { warnIfMapStale(dir, mapPath, cfg, g) })
	if !strings.HasPrefix(out, "⚠ STALE map: 2 governed file(s) are not in the map — e.g.: b.spec.md (and 1 more)\n") {
		t.Errorf("two missing files: one named, the other counted.\ngot: %s", out)
	}
}

// With no graph the warning keeps quiet: there is nothing to compare, and guessing would
// be worse.
func TestMapStaleWithoutGraphIsQuiet(t *testing.T) {
	if s := capturaSaida(t, func() {
		warnIfMapStale(t.TempDir(), "missing.yaml", &config.Config{}, nil)
	}); s != "" {
		t.Errorf("with no graph there is nothing to check; got: %s", s)
	}
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// capturaSaida collects what the function writes to stdout. (The timing tests of the
// `timing-metrics` flag use it too.)
func capturaSaida(t *testing.T, fn func()) string {
	t.Helper()
	return testkit.CaptureStdout(t, fn)
}

// THE DEFECT THIS CATCHES: an older binary does not fail — it writes the format it knows
// and UNDOES what the newer version wrote. Measured after renaming a stamp field: the
// earlier local build reverted 26 lines on every `check`, and the map swung between two
// formats with a conflict on every PR.
func TestWarnsWhenAnotherVersionWroteTheMap(t *testing.T) {
	t.Run("CGPCH-B33: A map written by another version is warned about", func(t *testing.T) {})
	w := staleBinaryWarning("0.1.9", "0.1.8")
	if w == "" {
		t.Fatal("different versions must warn")
	}
	// The warning must NAME both, or the reader does not know which one to update.
	if !strings.Contains(w, "0.1.9") || !strings.Contains(w, "0.1.8") {
		t.Errorf("the warning must name both versions; got: %s", w)
	}
	// And say what to do: a warning that only states the problem is noise.
	if !strings.Contains(w, "map build") {
		t.Errorf("the warning must say how to fix it; got: %s", w)
	}
}

// COMPARED BY EQUALITY, not by order — which is what makes the warning work in the real
// case. The two builds that produced the defect both called themselves "dev": ordering
// versions would have caught neither, because "dev" cannot be ordered against "0.1.9".
func TestTwoDevBuildsAreIndistinguishable(t *testing.T) {
	t.Run("CGPCH-I02: The version warning compares names by equality", func(t *testing.T) {})
	// The same name, even for distinct builds: nothing can be done here, and it is the
	// known limit. What must NOT happen is the reverse — keeping quiet when the names differ.
	if staleBinaryWarning("dev", "dev") != "" {
		t.Error("equal names cannot be told apart — warning here would be noise on every local build")
	}
	if staleBinaryWarning("dev", "0.1.9") == "" {
		t.Error("a map written by a local build and a published binary is exactly the case to warn about")
	}
	if staleBinaryWarning("0.1.9", "dev") == "" {
		t.Error("the reverse too: whoever runs a local build over a published map must know")
	}
}

// A MAP WITHOUT THE FIELD must not warn: every map generated before this version lacks
// `gerado_por`, and accusing them would make the first `check` of every existing project
// shout about something nobody can fix.
func TestOldMapDoesNotWarn(t *testing.T) {
	t.Run("CGPCH-B34: A map with no writer version raises no warning", func(t *testing.T) {})
	if staleBinaryWarning("", "0.1.9") != "" {
		t.Error("a map without `gerado_por` is every earlier version's — warning would accuse whoever did nothing wrong")
	}
}

// ── the profile table ───────────────────────────────────────────────────────────────

// The header counts everything the check found, by kind: failures (blocking and
// informative) and divergences. It said "N issue(s) — divergences recorded:" over the
// failures alone, which read as N issues plus the divergences, with the divergences
// counted nowhere.
func TestFindingsSummaryCountsEveryKind(t *testing.T) {
	t.Run("CGPCH-B57: The findings heading counts every kind", func(t *testing.T) {})
	englishOutput(t)
	block := gate.Result{Gate: "g1", Target: "a", Verdict: gate.Fail, Blocking: true}
	info := gate.Result{Gate: "g2", Target: "b", Verdict: gate.Fail}
	drift := gate.Result{Gate: "g3", Target: "c", Verdict: gate.Pending, Detail: "diverged"}
	p := gate.Profile{Results: []gate.Result{block, info, info, drift}, Failures: []gate.Result{block, info, info}}

	out := captureStdout(t, func() { printProfile(p, false, false) })
	for _, want := range []string{"3 failure(s) — 1 blocking, 2 informative — and 1 divergence(s)", "--show-drift"} {
		if !strings.Contains(out, want) {
			t.Errorf("the summary should say %q:\n%s", want, out)
		}
	}

	// Only failures: the heading is printed all the same.
	out = captureStdout(t, func() {
		printProfile(gate.Profile{Results: []gate.Result{info}, Failures: []gate.Result{info}}, false, false)
	})
	if !strings.Contains(out, "1 failure(s) — 0 blocking, 1 informative — and 0 divergence(s)") {
		t.Errorf("failures alone must still be counted:\n%s", out)
	}
	if strings.Contains(out, "issue(s) — divergences recorded") {
		t.Errorf("the old ambiguous header is back:\n%s", out)
	}

	// Only divergences: still counted, with no failure list.
	out = captureStdout(t, func() { printProfile(gate.Profile{Results: []gate.Result{drift}}, false, false) })
	if !strings.Contains(out, "0 failure(s)") || !strings.Contains(out, "1 divergence(s)") {
		t.Errorf("divergences alone must still be counted:\n%s", out)
	}
}

// profileOf builds a Profile with the counters the format cares about.
func profileOf(gates ...gate.GateSummary) gate.Profile {
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}, Passed: true}
	for _, g := range gates {
		p.ByGate[g.Gate] = g
	}
	return p
}

// captureProfile collects what `fn` prints, in the English catalog. `printProfile` writes
// straight to stdout (it is a terminal command, not a library), so testing the FORMAT
// means intercepting it.
func captureProfile(t *testing.T, fn func()) string {
	t.Helper()
	englishOutput(t)
	return captureStdout(t, fn)
}

// The table aligns to be COMPARED by eye: in a list of 49 gates, a crooked column is what
// makes the reader lose the number that matters.
func TestColumnsAlignWithALongName(t *testing.T) {
	t.Run("CGPCH-B40: The name column fits the longest name", func(t *testing.T) {})
	names := []string{"eslint", "handler-ddb-inline-passivo", "circular"}
	w := nameWidth(names)
	if w < len("handler-ddb-inline-passivo") {
		t.Fatalf("width %d cuts the longest name (%d)", w, len("handler-ddb-inline-passivo"))
	}
	// A floor: with short names the column does not shrink until it touches the verdict.
	if got := nameWidth([]string{"a", "bb"}); got < 18 {
		t.Errorf("nameWidth with no floor: %d", got)
	}
}

// ONE width for every column wastes space: if `~` reaches 582 and `✗` never passes 1, the
// fail column would reserve three places for nothing.
func TestWidthIsPerColumn(t *testing.T) {
	t.Run("CGPCH-B41: Each counter column has its own width", func(t *testing.T) {})
	p := profileOf(
		gate.GateSummary{Gate: "a", Pass: 1116, Fail: 0, Skip: 0},
		gate.GateSummary{Gate: "b", Pass: 1, Fail: 1, Skip: 582},
	)
	w := computeWidths(p)

	if w.pass != 4 {
		t.Errorf("pass: %d, want 4 (because of 1116)", w.pass)
	}
	if w.fail != 1 {
		t.Errorf("fail: %d, want 1 — a column must not inherit another's width", w.fail)
	}
	if w.skip != 3 {
		t.Errorf("skip: %d, want 3 (because of 582)", w.skip)
	}

	// The skip column's number is the skipped plus the pending, less the drift: a hundred
	// pending items need three places, and twelve pending that are all drift need one.
	if w := computeWidths(profileOf(gate.GateSummary{Gate: "a", Pending: 100})); w.skip != 3 {
		t.Errorf("skip: %d, want 3 (because of 100 pending)", w.skip)
	}
	allDrift := profileOf(gate.GateSummary{Gate: "a", Pending: 12})
	for i := 0; i < 12; i++ {
		allDrift.Results = append(allDrift.Results, gate.Result{Gate: "a", Verdict: gate.Pending, Detail: "diverged", Target: fmt.Sprint(i)})
	}
	if w := computeWidths(allDrift); w.skip != 1 || w.drift != 2 {
		t.Errorf("skip %d / drift %d, want 1 / 2 (the twelve pending are drift)", w.skip, w.drift)
	}
}

// The ALWAYS-present columns have a floor of 1: `%*d` with width 0 would print glued to
// the symbol. The drift column is the exception — it is born 0 and only opens on real drift.
func TestMinimumWidthIsOne(t *testing.T) {
	t.Run("CGPCH-B42: The always-present columns are at least one wide", func(t *testing.T) {})
	w := computeWidths(profileOf(gate.GateSummary{Gate: "a"}))
	for name, got := range map[string]int{
		"pass": w.pass, "fail": w.fail, "skip": w.skip, "judge": w.judge,
	} {
		if got != 1 {
			t.Errorf("%s: %d, want 1", name, got)
		}
	}
	// With no drift the column does not exist: reserving it would leave a hole in the
	// middle of every line with nothing to justify it — and a table without drift is the
	// common case, not the exception.
	if w.drift != 0 {
		t.Errorf("drift: %d, want 0 — the column must not exist without drift", w.drift)
	}
}

// The ⚠ column leaves the whole table when no gate has drift, and the separator leaves with
// it: otherwise two spaces would be left between `✗` and `~`.
func TestWithNoDriftTheColumnDoesNotExist(t *testing.T) {
	t.Run("CGPCH-B43: Without drift the drift column does not exist", func(t *testing.T) {})
	p := profileOf(
		gate.GateSummary{Gate: "one", Pass: 1},
		gate.GateSummary{Gate: "two", Pass: 583, Skip: 12},
	)
	out := captureProfile(t, func() { printProfile(p, false, false) })

	for _, line := range gateLines(out) {
		// Between the `✗` counter and the `~` there can only be the 2-space separator.
		i := runeIndex([]rune(line), '✗')
		j := runeIndex([]rune(line), '~')
		if i < 0 || j < 0 {
			continue
		}
		middle := string([]rune(line)[i:j])
		if strings.Count(middle, " ") > 2+len("0") {
			t.Errorf("a hole between ✗ and ~ in a table with no drift: %q", line)
		}
	}
}

// `--only-issues` omits whoever passed everything AND left nothing pending. A gate with
// `~` is not clean: it confronted nothing, and that is information.
func TestCleanGateRequiresNothingPending(t *testing.T) {
	t.Run("CGPCH-B45: A clean gate has nothing pending of any kind", func(t *testing.T) {})
	cases := []struct {
		name  string
		s     gate.GateSummary
		drift int
		clean bool
	}{
		{"all zero with passes", gate.GateSummary{Pass: 10}, 0, true},
		{"with a fail", gate.GateSummary{Pass: 10, Fail: 1}, 0, false},
		{"with drift", gate.GateSummary{Pass: 10}, 3, false},
		{"with skip", gate.GateSummary{Pass: 10, Skip: 2}, 0, false},
		{"with pending", gate.GateSummary{Pass: 10, Pending: 2}, 0, false},
		{"with as many skips as pending", gate.GateSummary{Pass: 10, Skip: 2, Pending: 2}, 0, false},
		{"awaiting AI", gate.GateSummary{Pass: 10, Judge: 1}, 0, false},
	}
	for _, c := range cases {
		if got := cleanGate(c.s, c.drift); got != c.clean {
			t.Errorf("%s: cleanGate = %v, want %v", c.name, got, c.clean)
		}
	}
}

// The default hides NOTHING: the full table is the proof that the 49 gates ran. Hiding by
// default would trade that proof for brevity.
func TestDefaultShowsACleanGate(t *testing.T) {
	t.Run("CGPCH-B46: The default table shows the clean gates", func(t *testing.T) {})
	p := profileOf(
		gate.GateSummary{Gate: "clean", Pass: 5},
		gate.GateSummary{Gate: "dirty", Pass: 1, Fail: 1},
	)
	out := captureProfile(t, func() { printProfile(p, false, false) })

	if !strings.Contains(out, "clean") {
		t.Errorf("the clean gate vanished from the default:\n%s", out)
	}
	if strings.Contains(out, "omitted by --only-issues") {
		t.Errorf("the omission footer appeared without the flag:\n%s", out)
	}
}

func TestOnlyIssuesOmitsCleanButCountsInTheFooter(t *testing.T) {
	t.Run("CGPCH-B47: Only-issues omits the clean gates and counts them", func(t *testing.T) {})
	p := profileOf(
		gate.GateSummary{Gate: "clean-one", Pass: 5},
		gate.GateSummary{Gate: "clean-two", Pass: 7},
		gate.GateSummary{Gate: "dirty", Pass: 1, Fail: 1},
	)
	out := captureProfile(t, func() { printProfile(p, true, false) })

	if strings.Contains(out, "clean-one") || strings.Contains(out, "clean-two") {
		t.Errorf("--only-issues did not omit the clean gate:\n%s", out)
	}
	if !strings.Contains(out, "dirty") {
		t.Errorf("--only-issues omitted a gate with a finding:\n%s", out)
	}
	// The omitted gate RAN. Without the number, the output would look like a smaller scan.
	if !strings.Contains(out, "2 gate(s) with nothing to report") {
		t.Errorf("the footer did not count the omitted ones:\n%s", out)
	}
}

// The `~` column has to fall in the SAME position with and without `⚠`. Before, a line
// without drift omitted the whole column and the `~` moved left on that line only — the
// same column existing in two places in the same table.
func TestSkipColumnDoesNotMoveWithOrWithoutDrift(t *testing.T) {
	t.Run("CGPCH-I03: The skip column does not move with or without drift", func(t *testing.T) {})
	p := profileOf(
		gate.GateSummary{Gate: "with-drift", Pass: 110, Pending: 423},
		gate.GateSummary{Gate: "without-drift", Pass: 2537, Skip: 103},
	)
	// `driftCount` reads Results: without them the ⚠ column would not even exist, and the
	// test would measure the wrong table.
	p.Results = []gate.Result{
		{Gate: "with-drift", Verdict: gate.Pending, Detail: "diverged", Target: "x"},
	}
	if without, with := driftColumn(0, 3), driftColumn(407, 3); len([]rune(without)) != len([]rune(with)) {
		t.Errorf("⚠ cell with different widths: without=%d runes, with=%d runes (%q vs %q)",
			len([]rune(without)), len([]rune(with)), without, with)
	}

	out := captureProfile(t, func() { printProfile(p, false, false) })
	var positions []int
	// `TrimSpace` over the whole block would eat the indentation of the FIRST line only —
	// exactly the column being measured. The cut is per line.
	for _, line := range gateLines(out) {
		// The position is the terminal COLUMN, counted in runes. `strings.Index` returns an
		// offset in BYTES, and `✓`/`✗`/`~` are multi-byte.
		if i := runeIndex([]rune(line), '~'); i >= 0 {
			positions = append(positions, i)
		}
	}
	if len(positions) < 2 {
		t.Fatalf("expected 2 lines with ~:\n%s", out)
	}
	for _, pos := range positions[1:] {
		if pos != positions[0] {
			t.Errorf("the ~ changed column between lines (%v):\n%s", positions, out)
		}
	}
}

// `⚠` takes 1 terminal column and 3 bytes in Go. Reserving the blank with `len()`
// misaligns exactly what it exists to align.
func TestDriftCellIsMeasuredInColumnsNotBytes(t *testing.T) {
	t.Run("CGPCH-B44: An empty drift cell is measured in terminal columns", func(t *testing.T) {})
	if got := len([]rune(driftColumn(0, 3))); got != 4 {
		t.Errorf("empty cell with %d columns, want 4 (symbol + 3 digits)", got)
	}
}

// runeIndex returns the position of `target` counted in runes (columns), not bytes.
func runeIndex(runes []rune, target rune) int {
	for i, r := range runes {
		if r == target {
			return i
		}
	}
	return -1
}

// gateLines keeps the lines of the TABLE. The legend and the detail blocks also contain
// `✓`/`✗`/`~`, and measuring them as columns would accuse misalignment where there is no
// table at all.
func gateLines(out string) []string {
	var lines []string
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "blocking") || strings.Contains(l, "informative") || strings.Contains(l, "judgment") {
			lines = append(lines, l)
		}
	}
	return lines
}

// With the flag, drift does NOT inherit the `~` noise cut: the `~` is suppressed on a
// large scan because thousands of "does not apply" are noise, but whoever asked for the
// addresses of the pending items wants them all, and the size of the scan does not change
// that.
func TestShowDriftDoesNotInheritTheSkipCut(t *testing.T) {
	t.Run("CGPCH-B49: Show-drift is not cut by the size of the scan", func(t *testing.T) {})
	p := profileOf(gate.GateSummary{Gate: "layer-boundary", Pass: 1, Pending: 1})
	// A volume above the cut that suppresses the `~` detail.
	for i := 0; i < maxResultsForSkipDetail*3; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "other", Verdict: gate.Skip, Detail: "does not apply", Target: "n",
		})
	}
	p.Results = append(p.Results, gate.Result{
		Gate: "layer-boundary", Verdict: gate.Pending,
		Detail: "boundary being migrated", Target: "AlertSheet.tsx",
	})

	out := captureProfile(t, func() { printProfile(p, false, true) })

	if !strings.Contains(out, "pending item(s)") {
		t.Errorf("the drift block is missing on a large scan:\n%s", out)
	}
	if !strings.Contains(out, "AlertSheet.tsx") {
		t.Errorf("drift with no address — it is the only actionable part:\n%s", out)
	}
}

// The `~` stays cut on a large scan: 2,430 lines of "does not apply" are noise, and the
// counter is enough.
func TestSkipStaysSuppressedOnALargeScan(t *testing.T) {
	t.Run("CGPCH-B52: A large scan does not list the skip reasons", func(t *testing.T) {})
	p := profileOf(gate.GateSummary{Gate: "g", Pass: 1, Skip: 1})
	for i := 0; i < maxResultsForSkipDetail*3; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "g", Verdict: gate.Skip, Detail: "not a screen", Target: "n",
		})
	}
	out := captureProfile(t, func() { printProfile(p, false, false) })

	if strings.Contains(out, "indeterminate — not a failure") {
		t.Errorf("the ~ detail should be suppressed on a large volume:\n%s", out)
	}
}

// Whoever asks for the addresses wants to act on them: truncating the list would bring
// back the problem it exists to solve. Either it is complete, or the table's counter was
// enough.
func TestShowDriftListsAllWithNoCap(t *testing.T) {
	t.Run("CGPCH-B48: Show-drift lists every drift item", func(t *testing.T) {})
	const n = 120
	p := profileOf(gate.GateSummary{Gate: "g", Pass: 1, Pending: n})
	for i := 0; i < n; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "g", Verdict: gate.Pending, Detail: "diverged",
			Target: fmt.Sprintf("target-%03d", i),
		})
	}
	out := captureProfile(t, func() { printProfile(p, false, true) })

	for _, i := range []int{0, n / 2, n - 1} {
		target := fmt.Sprintf("target-%03d", i)
		if !strings.Contains(out, target) {
			t.Errorf("%s was not listed — the list is truncated", target)
		}
	}
	if strings.Contains(out, "… and") {
		t.Errorf("list truncated under --show-drift:\n%s", out)
	}
}

// Without the flag, the table's counter is all: on a large scan there are thousands of
// lines, and dumping them unasked would bury the issues just below.
func TestWithoutShowDriftNoDetailIsListed(t *testing.T) {
	t.Run("CGPCH-B50: Without show-drift only the counter appears", func(t *testing.T) {})
	p := profileOf(gate.GateSummary{Gate: "g", Pass: 1, Pending: 1})
	p.Results = append(p.Results, gate.Result{
		Gate: "g", Verdict: gate.Pending, Detail: "diverged", Target: "single-target",
	})
	out := captureProfile(t, func() { printProfile(p, false, false) })

	if strings.Contains(out, "single-target") {
		t.Errorf("detail listed without the flag:\n%s", out)
	}
	// The signal does not vanish: the table keeps counting.
	if !strings.Contains(out, "⚠1") {
		t.Errorf("the ⚠ counter vanished from the table:\n%s", out)
	}
}

// The legend explains only what the table used: listing `⚠` on a scan with no drift
// teaches people to skip the whole legend.
func TestLegendShowsOnlyTheSymbolsUsed(t *testing.T) {
	t.Run("CGPCH-B53: The legend explains only the symbols used", func(t *testing.T) {})
	noDrift := captureProfile(t, func() {
		printProfile(profileOf(gate.GateSummary{Gate: "g", Pass: 1}), false, false)
	})
	if strings.Contains(noDrift, "⚠  diverged") {
		t.Errorf("the legend explained ⚠ on a table with no drift:\n%s", noDrift)
	}
	if !strings.Contains(noDrift, "✓  passed") || !strings.Contains(noDrift, "~  indeterminate") {
		t.Errorf("incomplete legend:\n%s", noDrift)
	}
	if strings.Contains(noDrift, "⏳") {
		t.Errorf("the legend explained ⏳ on a table with no judgment pending:\n%s", noDrift)
	}

	withJudge := captureProfile(t, func() {
		printProfile(profileOf(gate.GateSummary{Gate: "g", Judge: 2}), false, false)
	})
	if !strings.Contains(withJudge, "⏳ awaiting AI judgment") {
		t.Errorf("a legend without ⏳ on a table that uses it:\n%s", withJudge)
	}
}

// 832 pending items with the IDENTICAL reason are not 832 problems — they are one (the
// ingestion nobody ran). Repeating the same paragraph 832 times hides that reading instead
// of revealing it.
func TestDriftGroupsARepeatedReason(t *testing.T) {
	t.Run("CGPCH-B54: A repeated drift reason is written once with its targets", func(t *testing.T) {})
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}}
	for i := 0; i < 50; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "mutation-score", Verdict: gate.Pending,
			Detail: "no mutation signal ingested", Target: fmt.Sprintf("target-%02d", i),
		})
	}
	out := captureProfile(t, func() { printDrift(driftResults(p)) })

	if n := strings.Count(out, "no mutation signal ingested"); n != 1 {
		t.Errorf("the reason repeated %d times, want 1:\n%s", n, out)
	}
	if !strings.Contains(out, "50 target(s):") {
		t.Errorf("the group heading is missing:\n%s", out)
	}
	// One target per LINE: joining them with commas produced a 36-thousand-character line
	// in the real project — the address was there and nobody could read it.
	for _, l := range strings.Split(out, "\n") {
		if len(l) > 200 {
			t.Errorf("a %d-character line — the targets were concatenated again", len(l))
			break
		}
	}
	// Grouping cannot cost the ADDRESS — it is the actionable part.
	for _, target := range []string{"target-00", "target-25", "target-49"} {
		if !strings.Contains(out, target) {
			t.Errorf("%s vanished in the grouping", target)
		}
	}
	if !strings.Contains(out, "mutation-score — 50") {
		t.Errorf("the gate heading has no total:\n%s", out)
	}
}

// Where each reason is UNIQUE (the divergence is specific to the unit), the list stays
// target by target: compacting there would lose exactly what tells one finding from another.
func TestDriftDoesNotGroupUniqueReasons(t *testing.T) {
	t.Run("CGPCH-B55: Distinct drift reasons stay target by target", func(t *testing.T) {})
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}}
	for i := 0; i < 3; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "feature-test-match", Verdict: gate.Pending,
			Detail: fmt.Sprintf("scenario %d diverges", i), Target: fmt.Sprintf("f%d.feature", i),
		})
	}
	out := captureProfile(t, func() { printDrift(driftResults(p)) })

	for i := 0; i < 3; i++ {
		if !strings.Contains(out, fmt.Sprintf("scenario %d diverges", i)) {
			t.Errorf("reason %d vanished:\n%s", i, out)
		}
	}
	if strings.Contains(out, "target(s):") {
		t.Errorf("it compacted reasons that differ from each other:\n%s", out)
	}
}

// The heading says in how many GATES the pending items are: it is the first useful reading
// — 2,430 spread over 5 gates is a different picture from 2,430 in one.
func TestDriftHeadingCountsGates(t *testing.T) {
	t.Run("CGPCH-B56: The drift heading counts items and gates", func(t *testing.T) {})
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}}
	for _, g := range []string{"a", "b", "a", "c"} {
		p.Results = append(p.Results, gate.Result{
			Gate: g, Verdict: gate.Pending, Detail: "x", Target: "t",
		})
	}
	out := captureProfile(t, func() { printDrift(driftResults(p)) })
	if !strings.Contains(out, "4 pending item(s) in 3 gate(s)") {
		t.Errorf("wrong heading:\n%s", out)
	}
}

// Ten gates build the message with `strings.Join(findings, "; ")`. Five violations in the
// same file came out on an 800-character line: the reader cannot tell where one ends and
// the next begins, and the rule's text — identical in all five — drowns the only datum
// that varies, the line number.
func TestOccurrencesOfTheSameFileBreakIntoLines(t *testing.T) {
	t.Run("CGPCH-B59: Occurrences of a detail are printed one per line", func(t *testing.T) {})
	detail := "line 8: hex colour; line 10: hex colour; line 12: hex colour"
	got := indent(detail, "    ")

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d:\n%s", len(lines), got)
	}
	// The indentation holds for ALL of them: otherwise the second occurrence onwards would
	// touch the margin and leave the block it belongs to.
	for _, l := range lines {
		if !strings.HasPrefix(l, "    ") {
			t.Errorf("a line with no indentation: %q", l)
		}
	}
}

// ONE separator is enough: two occurrences already blur on one line, and it is the most
// common case. (With the cut at >= 2 separators, the pair did not break.)
func TestTwoOccurrencesAlreadyBreak(t *testing.T) {
	t.Run("CGPCH-B60: Two occurrences already break", func(t *testing.T) {})
	got := breakOccurrences("line 8: error; line 10: error")
	if !strings.Contains(got, "\n") {
		t.Errorf("the pair did not break: %q", got)
	}
}

// A one-sentence detail stays whole — there is nothing to separate.
func TestASingleOccurrenceDoesNotBreak(t *testing.T) {
	t.Run("CGPCH-B61: A single occurrence stays whole", func(t *testing.T) {})
	got := breakOccurrences("the spec does not declare `## Open Decisions`")
	if strings.Contains(got, "\n") {
		t.Errorf("it broke where there was no separator: %q", got)
	}
}

// The break is at the `"; "` separator (with a space), not at any `;`. Messages carry regex
// and code excerpts — the `;` of a `#[0-9A-Fa-f]{3};` separates no occurrence, and breaking
// there would cut the message in half.
func TestAGluedSemicolonDoesNotBreak(t *testing.T) {
	t.Run("CGPCH-B62: A glued semicolon does not break", func(t *testing.T) {})
	got := breakOccurrences("the layer cannot contain `a;b;c` — use a token")
	if strings.Contains(got, "\n") {
		t.Errorf("it broke at a glued `;`, which does not separate occurrences: %q", got)
	}
}

// 17 gates build the message with `strings.Join(list, ", ")` without cutting. A large list
// becomes a line nobody reads or greps.
func TestALongCommaListBreaks(t *testing.T) {
	t.Run("CGPCH-B63: A long list of items breaks one per line", func(t *testing.T) {})
	items := make([]string, 30)
	for i := range items {
		items[i] = fmt.Sprintf("src/features/some/path/File%02d.tsx", i)
	}
	got := indent("symbols with no catalogue: "+strings.Join(items, ", "), "    ")

	for _, l := range strings.Split(strings.TrimRight(got, "\n"), "\n") {
		if len([]rune(l)) > listBreakThreshold+20 {
			t.Errorf("a %d-rune line — the list did not break:\n%s", len([]rune(l)), l)
		}
	}
	if !strings.Contains(got, "File29.tsx") {
		t.Errorf("the break lost items:\n%s", got)
	}
}

// In ordinary text the comma separates clauses: cutting it would destroy the sentence.
func TestAShortSentenceWithCommasDoesNotBreak(t *testing.T) {
	t.Run("CGPCH-B64: A short sentence with commas stays whole", func(t *testing.T) {})
	sentence := "the spec exists, the code exists, and the two reference each other"
	if got := breakOccurrences(sentence); strings.Contains(got, "\n") {
		t.Errorf("it broke an ordinary sentence: %q", got)
	}
}

// The list break cannot cut PROSE. The first version used "a long line with 3+ commas" and
// broke the `rule-implemented` message itself into
// "the spec exists,\n the code exists,\n the two reference each other".
func TestLongProseWithCommasDoesNotBreak(t *testing.T) {
	t.Run("CGPCH-B65: Long prose with commas stays whole", func(t *testing.T) {})
	prose := "A catalogued rule with no implementer crosses the whole pipeline — the spec exists, " +
		"the code exists, the two reference each other through the header, and every gate turns green " +
		"over work that was never done."
	if got := breakOccurrences(prose); strings.Contains(got, "\n") {
		t.Errorf("it cut a sentence into clauses:\n%s", got)
	}
}

// And it still breaks a real list — items with no inner space.
func TestAListOfPathsStillBreaks(t *testing.T) {
	t.Run("CGPCH-B66: A list of paths still breaks", func(t *testing.T) {})
	items := make([]string, 20)
	for i := range items {
		items[i] = fmt.Sprintf("apps/mobile/src/features/Module%02d.spec.md", i)
	}
	if got := breakOccurrences(strings.Join(items, ", ")); !strings.Contains(got, "\n") {
		t.Errorf("a list of paths did not break:\n%s", got)
	}
}

// A gate whose name is a prefix of another's (`review`, `review-deep`) must not steal the
// other's tasks: the gate is read from the ID and the task's target together, never by the
// first known name that prefixes the ID.
func TestDropStaleJudgments_prefixGateNameDoesNotStealTheTask(t *testing.T) {
	t.Run("CGPCH-B68: A judge task is read back as the gate that queued it even when another gate's name prefixes it", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "src", "x.go"), []byte("package src\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "review", Measures: config.MeasuresJudgment},
		{Name: "review-deep", Measures: config.MeasuresJudgment},
	}}
	deep := queue.Task{ID: "judge-review-deep-src-x", Changed: "src/x.go", Kind: "judgment", Origin: "check"}
	if _, err := queue.Enqueue(root, deep); err != nil {
		t.Fatal(err)
	}
	// this round, review-deep judged src/x.go again: its task is live
	p := gate.Profile{Judged: []gate.Result{{Gate: "review-deep", Target: "src/x.go"}}}
	knownJudgmentGates = []string{"review", "review-deep"}
	dropStaleJudgments(root, cfg, p)
	tasks, err := queue.List(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].ID != deep.ID {
		t.Fatalf("the live review-deep task was dropped as a stale `review` task: %+v", tasks)
	}
}

// A queue written before the judge task suggested `review` still holds `suggested_next:
// judge`: those tasks are judgments too, and still bar the incremental check.
func TestQueuedJudgmentsCountsTheLegacyVerb(t *testing.T) {
	t.Run("CGPCH-B69: A judge task suggests the review stage, a verb the work command composes", func(t *testing.T) {})
	root := t.TempDir()
	for _, f := range []string{"a.go", "b.go", "c.go"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("package a\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, tk := range []queue.Task{
		{ID: "judge-g-a", Changed: "a.go", Kind: "judgment", Origin: "check", SuggestedNext: "review"},
		{ID: "judge-g-b", Changed: "b.go", Kind: "judgment", Origin: "check", SuggestedNext: "judge"},
		{ID: "001-change-c", Changed: "c.go", Kind: "change", Origin: "watch", SuggestedNext: "review"},
	} {
		if _, err := queue.Enqueue(root, tk); err != nil {
			t.Fatal(err)
		}
	}
	if n := queuedJudgments(root); n != 2 {
		t.Errorf("the two judge tasks (one with the legacy verb) are judgments, the watch review is not; got %d", n)
	}
}

// The unit's pieces enter the impact by identity, whatever the language: a Go test is
// `<stem>_test.go`, and the check of it must reach its spec, feature and code even when
// the map holds no edge between them.
func TestImpactOfBringsTheUnitsPiecesOfAGoUnit(t *testing.T) {
	t.Run("CGPCH-B70: The pieces of the changed file's unit enter its impact path, for a Go unit as for a TypeScript one", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/foo.go", Kind: mapx.KindCode},
		{ID: "pkg/foo_test.go", Kind: mapx.KindTest},
		{ID: "pkg/foo.spec.md", Kind: mapx.KindSpec},
		{ID: "pkg/foo.feature", Kind: mapx.KindFeature},
		{ID: "pkg/bar.go", Kind: mapx.KindCode},
		{ID: "web/x.ts", Kind: mapx.KindCode},
		{ID: "web/x.test.ts", Kind: mapx.KindTest},
		{ID: "web/x.spec.md", Kind: mapx.KindSpec},
		{ID: "py/baz.py", Kind: mapx.KindCode},
		{ID: "py/baz_test.py", Kind: mapx.KindTest},
	}}
	root := t.TempDir()
	cases := map[string][]string{
		"pkg/foo_test.go": {"pkg/foo.go", "pkg/foo.spec.md", "pkg/foo.feature"},
		"pkg/foo.go":      {"pkg/foo_test.go", "pkg/foo.spec.md", "pkg/foo.feature"},
		"pkg/foo.spec.md": {"pkg/foo.go", "pkg/foo_test.go", "pkg/foo.feature"},
		"web/x.test.ts":   {"web/x.ts", "web/x.spec.md"},
		// a language with no fixed name: the code and the test carry the file's own extension
		"py/baz_test.py": {"py/baz.py"},
		"py/baz.py":      {"py/baz_test.py"},
	}
	for changed, want := range cases {
		ids, err := impactOf(g, &config.Config{}, filepath.Join(root, changed), root)
		if err != nil {
			t.Fatal(err)
		}
		has := map[string]bool{}
		for _, id := range ids {
			has[id] = true
		}
		for _, w := range want {
			if !has[w] {
				t.Errorf("--changed %s: the unit's piece %s is not in the impact: %v", changed, w, ids)
			}
		}
		if has["pkg/bar.go"] {
			t.Errorf("--changed %s: another unit entered the impact: %v", changed, ids)
		}
	}
}

// ── the command's announcements ─────────────────────────────────────────────────────

// The version warning reaches the command's error output: the pure function above says
// what the warning is, this proves the check prints it.
func TestCheckPrintsTheVersionWarning(t *testing.T) {
	t.Run("CGPCH-B33: A map written by another version is warned about", func(t *testing.T) {})
	englishOutput(t)
	prev := mapx.GeneratedBy
	mapx.GeneratedBy = "0.0.1-elsewhere"
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	mapx.GeneratedBy = prev
	var runErr error
	errOut := captureStderr(t, func() {
		_, runErr = runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	})
	if runErr != nil {
		t.Fatal(runErr)
	}
	if !strings.Contains(errOut, "0.0.1-elsewhere") {
		t.Errorf("a map written by another version must be warned about on the error output:\n%s", errOut)
	}
}

func TestCheckPrintsNoJudgmentBriefWithNothingToJudge(t *testing.T) {
	t.Run("CGPCH-B71: With nothing awaiting judgment there is no judgment brief", func(t *testing.T) {})
	englishOutput(t)
	t.Setenv("GITHUB_ACTIONS", "")
	defer issue.UseFiles()
	dir := qProject(t, checkYAML+githubWorkflowYAML, checkFiles(), checkGraph())
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "awaiting judgment") {
		t.Errorf("no gate judges here, and a judgment brief was printed:\n%s", out)
	}
}

func TestCheckWarnsOnlyWhenTheRecordFails(t *testing.T) {
	t.Run("CGPCH-B72: A record that fails is warned about and the check still reports", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, checkYAML, checkFiles(), checkGraph())
	var runErr error
	var out string
	errOut := captureStderr(t, func() {
		out, runErr = runQ(t, newCheckCmd(), "--root", dir, "--changed", "a.spec.md")
	})
	if runErr != nil {
		t.Fatal(runErr)
	}
	if strings.Contains(errOut, "failed to record") || !strings.Contains(out, "recorded: 1 edge(s) stamped") {
		t.Errorf("a record that succeeded must raise no warning:\nstdout:\n%s\nstderr:\n%s", out, errOut)
	}

	// the map lives in a folder the check cannot write
	locked := t.TempDir()
	mapPath := filepath.Join(locked, "graph.yaml")
	if err := mapx.Save(checkGraph(), mapPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(mapPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(locked, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755); _ = os.Chmod(mapPath, 0o644) })
	errOut = captureStderr(t, func() {
		out, runErr = runQ(t, newCheckCmd(), "--root", dir, "--map", mapPath, "--changed", "a.spec.md")
	})
	if runErr != nil {
		t.Fatal(runErr)
	}
	if !strings.Contains(errOut, "warning: failed to record") || !strings.Contains(out, "can promote") {
		t.Errorf("a record that failed must be warned about, and the verdict still printed:\nstdout:\n%s\nstderr:\n%s", out, errOut)
	}
}

func TestCheckSaysHowManyJudgmentsItQueued(t *testing.T) {
	t.Run("CGPCH-B73: The check says how many judgments it queued, and nothing when it queued none", func(t *testing.T) {})
	englishOutput(t)
	dir := qProject(t, judgmentYAML, map[string]string{"a.go": "package a\n"},
		&mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}})
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "1 target(s) awaiting AI judgment") {
		t.Errorf("the first run queues one judgment and must say so:\n%s", out)
	}
	// the task already waits in the queue: this run queues nothing new
	out, err = runQ(t, newCheckCmd(), "--root", dir, "--all")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "target(s) awaiting AI judgment") {
		t.Errorf("a run that queued nothing must not announce a queued judgment:\n%s", out)
	}
}

func TestCheckWithoutGovernanceTipsPrintsNoPointer(t *testing.T) {
	t.Run("CGPCH-B74: With no governance tip the full sweep prints neither a tip nor the pointer to the doctor", func(t *testing.T) {})
	englishOutput(t)
	// a project of specs only: no code, no test, nothing the tips ask about
	yaml := "version: 2\nlayers:\n  spec:\n    kind: spec\n    pattern: \"*.spec.md\"\ngates:\n  - name: spec-ok\n    on: [spec]\n    run: \"true\"\n"
	dir := qProject(t, yaml, map[string]string{"a.spec.md": "# A\n"},
		&mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec}}})
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "Governance Tip") || strings.Contains(out, "anchors doctor") {
		t.Errorf("with no tip to give, the full sweep must print neither a tip nor the pointer:\n%s", out)
	}
}

func TestRecordSaysWhyNoIssueWasWritten(t *testing.T) {
	t.Run("CGPCH-B75: When the check writes no issue it says why", func(t *testing.T) {})
	englishOutput(t)
	issue.UseFiles()
	root := qProject(t, "version: 2\nworkflow:\n  mode: manual\n", nil, nil)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode, Rev: "r1"}}}
	out := captureStdout(t, func() {
		if err := recordCheck(root, filepath.Join(root, mapx.DefaultPath), g, gate.Profile{}, false, false); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "no issue written (mode: manual") {
		t.Errorf("in manual mode the record must say the mode is why no issue was written:\n%s", out)
	}
}

// fakeBoard puts on the PATH a `gh` that answers nothing and logs every call, so a test
// can tell whether the check talked to the board without ever reaching GitHub.
func fakeBoard(t *testing.T) string {
	t.Helper()
	bin := t.TempDir()
	log := filepath.Join(bin, "calls.log")
	script := "#!/bin/sh\necho \"gh $*\" >> \"" + log + "\"\nexit 1\n"
	testkit.FakeBin(t, bin, "gh", script)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return log
}

func TestCheckRecordsOnTheBoardInGitHubMode(t *testing.T) {
	t.Run("CGPCH-B75: When the check writes no issue it says why", func(t *testing.T) {})
	t.Run("CGPCH-B76: In github mode the issues go to the board, never to the local folders", func(t *testing.T) {})
	englishOutput(t)
	t.Setenv("GITHUB_ACTIONS", "")
	defer issue.UseFiles()
	log := fakeBoard(t)
	calls := func() string { b, _ := os.ReadFile(log); return string(b) }

	dir := qProject(t, checkYAML+githubWorkflowYAML, checkFiles(), checkGraph())
	// a local run leaves the board to CI, and says so
	out, err := runQ(t, newCheckCmd(), "--root", dir, "--all")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "board issues are left to CI") || calls() != "" {
		t.Errorf("a local run in github mode must leave the board to CI:\n%s\ngh calls:\n%s", out, calls())
	}
	// asked to record, it records on the board, never in the local folders
	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--record-issues"); err != nil {
		t.Fatal(err)
	}
	if calls() == "" {
		t.Error("with --record-issues in github mode the check must record on the board")
	}
	if _, err := os.Stat(filepath.Join(dir, issue.Dir)); err == nil {
		t.Errorf("in github mode no issue folder may be written")
	}

}

func TestTableIndeterminateCounterLeavesTheDriftOut(t *testing.T) {
	t.Run("CGPCH-B78: The indeterminate counter is the skipped and pending less the drift", func(t *testing.T) {})
	p := profileOf(gate.GateSummary{Gate: "g", Pass: 1, Skip: 1, Pending: 2})
	p.Results = []gate.Result{{Gate: "g", Verdict: gate.Pending, Detail: "diverged", Target: "x"}}
	out := captureProfile(t, func() { printProfile(p, false, false) })
	if !strings.Contains(out, "⚠1  ~2\n") {
		t.Errorf("one skip and two pending, one of them drift: ⚠1 and ~2:\n%s", out)
	}
}

func TestFindingsListEachFailureWithItsDetail(t *testing.T) {
	t.Run("CGPCH-B82: Each failure is listed with its detail, and no finding means no findings heading", func(t *testing.T) {})
	block := gate.Result{Gate: "g1", Target: "a.go", Verdict: gate.Fail, Blocking: true, Detail: "broken on line 3"}
	info := gate.Result{Gate: "g2", Target: "b.go", Verdict: gate.Fail}
	p := gate.Profile{Results: []gate.Result{block, info}, Failures: []gate.Result{block, info}}
	out := captureProfile(t, func() { printProfile(p, false, false) })
	if !strings.Contains(out, "  ✗ [BLOCKS] g1 @ a.go\n      broken on line 3\n  ✗ [informative] g2 @ b.go\n") {
		t.Errorf("each failure with its mark, and the detail under the one that has it:\n%s", out)
	}

	pass := gate.Result{Gate: "g", Target: "a.go", Verdict: gate.Pass}
	out = captureProfile(t, func() {
		printProfile(gate.Profile{Passed: true, Results: []gate.Result{pass}}, false, false)
	})
	if strings.Contains(out, "failure(s)") {
		t.Errorf("with no failure and no divergence there is no findings heading:\n%s", out)
	}
}

func TestMaturationNamesTheCleanInformativeGates(t *testing.T) {
	t.Run("CGPCH-B83: A clean informative gate is named as ready to become blocking", func(t *testing.T) {})
	englishOutput(t)
	clean := profileOf(gate.GateSummary{Gate: "lint", Pass: 3})
	out := captureStdout(t, func() { rememberMaturation(clean, false) })
	if !strings.Contains(out, "1 informative gate(s) CLEAN") || !strings.Contains(out, "lint") {
		t.Errorf("a clean informative gate must be named as ready to become blocking:\n%s", out)
	}
	blocking := profileOf(gate.GateSummary{Gate: "lint", Pass: 3, Blocking: true})
	if out := captureStdout(t, func() { rememberMaturation(blocking, false) }); out != "" {
		t.Errorf("with no gate to promote nothing is said:\n%s", out)
	}
}

// The limits of the list break, each at its boundary.
func TestListBreakBoundaries(t *testing.T) {
	t.Run("CGPCH-B63: A long list of items breaks one per line", func(t *testing.T) {})
	item := func(n int, c string) string { return strings.Repeat(c, n) }
	// four tokens of 26 make exactly 110 characters: not longer than the threshold
	at := strings.Join([]string{item(26, "a"), item(26, "b"), item(26, "c"), item(26, "d")}, ", ")
	if len([]rune(at)) != listBreakThreshold {
		t.Fatalf("fixture is %d runes", len([]rune(at)))
	}
	if strings.Contains(breakOccurrences(at), "\n") {
		t.Errorf("a line of exactly %d characters must not break", listBreakThreshold)
	}
	// four tokens of 27: longer, and four items is enough
	four := strings.Join([]string{item(27, "a"), item(27, "b"), item(27, "c"), item(27, "d")}, ", ")
	if strings.Count(breakOccurrences(four), "\n") != 3 {
		t.Errorf("four items over the threshold must break one per line:\n%s", breakOccurrences(four))
	}
	// half of the items with an inner space is not "most of them without"
	half := strings.Join([]string{item(27, "a"), item(27, "b"), "some words " + item(20, "c"), "more words " + item(20, "d")}, ", ")
	if strings.Contains(breakOccurrences(half), "\n") {
		t.Errorf("a line whose items are only half tokens must not break:\n%s", breakOccurrences(half))
	}
}

// The SCENARIOS of the `timing-metrics` flag (`flags/timing-metrics.flag.md`).
//
// The flag governs whether `check` measures and prints the time per gate. The three
// scenarios are here, each with its code in a subtest NAME: that is how `flag-covered`
// recognises them as proven once the execution is ingested.
//
// `--timing` was born to find what makes a scan expensive, and it found it: `docs-fresh`
// was 97% of a 6m49s run. But it was born with no test at all, and it was `flag-covered`
// itself that charged it — the first work the flag axis found.

// profileWithTime builds a Profile with measured time, which is what `printTiming` reads.
func profileWithTime() gate.Profile {
	return gate.Profile{
		ByGate: map[string]gate.GateSummary{
			"docs-fresh": {Gate: "docs-fresh", Pass: 3, Duracao: 900 * time.Millisecond, Pior: 800 * time.Millisecond},
			"build":      {Gate: "build", Pass: 1, Duracao: 50 * time.Millisecond, Pior: 50 * time.Millisecond},
		},
		Results: []gate.Result{
			{Gate: "docs-fresh", Target: "a.spec.md", Duracao: 800 * time.Millisecond},
			{Gate: "build", Target: "b.go", Duracao: 50 * time.Millisecond},
		},
	}
}

// TIMNG-G01 — the value is `off`: check prints only the verdicts, and measures NO time.
//
// The proof is the absence: with the flag off neither the timing header nor the target
// list may print. A test that only checked "nothing broke" would pass with the whole
// table on screen.
func TestTimingG01_offPrintsNoTiming(t *testing.T) {
	t.Run("TIMNG-G01: with the flag off, check prints no timing at all", func(t *testing.T) {})
	// It confronts the VARIABLE the command reads, not a literal: `if off := false` would
	// be a tautology — it would pass with `printTiming` called unconditionally in
	// production, which is exactly the regression this scenario exists to catch.
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{"--timing=false"}); err != nil {
		t.Fatal(err)
	}

	on, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	// The PRODUCTION function, not a copy of the condition: `RunE` calls exactly this.
	out := capturaSaida(t, func() { reportTiming(on, profileWithTime()) })
	if out != "" {
		t.Errorf("with the flag `off` check printed timing:\n%s", out)
	}
}

// TIMNG-G02 — the value is `on`: check also prints time per gate, and the slowest
// targets.
//
// Both blocks are charged, not just one: the per-gate header and the target list answer
// different questions ("which gate costs" and "which file costs"), and it was the second
// that pointed at `fnSize` reading ~43,000 files.
func TestTimingG02_onPrintsTimePerGateAndTargets(t *testing.T) {
	t.Run("TIMNG-G02: with the flag on, check prints time per gate and the slowest targets", func(t *testing.T) {})
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{"--timing"}); err != nil {
		t.Fatal(err)
	}

	on, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	out := capturaSaida(t, func() { reportTiming(on, profileWithTime()) })
	if out == "" {
		t.Fatal("with the flag `on` check printed nothing")
	}
	for _, want := range []string{"docs-fresh", "build"} {
		if !strings.Contains(out, want) {
			t.Errorf("the output does not name the gate %q:\n%s", want, out)
		}
	}
	// The most expensive individual TARGET — the block the per-gate average does not show.
	if !strings.Contains(out, "a.spec.md") {
		t.Errorf("the output does not list the slowest target:\n%s", out)
	}
	// And the most expensive comes first: the order is what makes the table actionable.
	if strings.Index(out, "docs-fresh") > strings.Index(out, "build") {
		t.Errorf("the most expensive gate did not come first:\n%s", out)
	}
}

// TIMNG-G03 — the value is ABSENT: the same as `off`. Measuring is opt-in, never a
// default cost.
//
// It is the scenario `flag-scenarios-complete` charges and the one nobody writes. It
// matters here for a concrete reason: an inverted default would raise no error at all — it
// would just spend everyone's time, forever, in silence.
func TestTimingG03_absentMeansTheDefaultWhichIsOff(t *testing.T) {
	t.Run("TIMNG-G03: with the flag absent, the declared default holds — measuring is opt-in", func(t *testing.T) {})
	// The absence of the VALUE is the absence of the flag on the command line. Parsing an
	// argv without `--timing` leaves the flag in the state the scenario describes — and it
	// is that state, not a literal, that the call site reads.
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	on, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	out := capturaSaida(t, func() { reportTiming(on, profileWithTime()) })
	if out != "" {
		t.Errorf("without `--timing` check measured and printed — measuring must be opt-in:\n%s", out)
	}
	// And the DEFAULT declared in cobra must be the same. It is the other half of the
	// scenario: the block above proves the behaviour with the value absent, and this one
	// proves that absent is really what the command delivers when nobody passes `--timing`.
	f := newCheckCmd().Flags().Lookup("timing")
	if f == nil {
		t.Fatal("the `--timing` flag disappeared from the command")
	}
	if f.DefValue != "false" {
		t.Errorf("the default of `--timing` is %q — measuring stopped being opt-in", f.DefValue)
	}
}

// timingTable prints the time table of three gates: one with a target of every verdict,
// one with five targets and one with a single target.
func timingTable(t *testing.T) string {
	t.Helper()
	englishOutput(t)
	p := gate.Profile{
		ByGate: map[string]gate.GateSummary{
			"alpha-long-name": {Gate: "alpha-long-name", Pass: 1, Fail: 2, Skip: 3, Pending: 4, Judge: 5, Duracao: 900 * time.Millisecond, Pior: 300 * time.Millisecond},
			"b":               {Gate: "b", Pass: 5, Duracao: 100 * time.Millisecond, Pior: 50 * time.Millisecond},
			"one":             {Gate: "one", Pass: 1, Duracao: 500 * time.Millisecond, Pior: 500 * time.Millisecond},
		},
		// the slowest target belongs to the gate whose name sorts last
		Results: []gate.Result{
			{Gate: "b", Target: "x.go", Duracao: 50 * time.Millisecond},
			{Gate: "alpha-long-name", Target: "y.go", Duracao: 300 * time.Millisecond},
			{Gate: "one", Target: "z.go", Duracao: 500 * time.Millisecond},
		},
	}
	return capturaSaida(t, func() { printTiming(p) })
}

func TestTimingTableCountsAndAlignsTheTargets(t *testing.T) {
	t.Run("CGPCH-B79: The time table counts each gate's targets and aligns its columns", func(t *testing.T) {})
	out := timingTable(t)
	for _, want := range []string{
		// every verdict counts: 1+2+3+4+5
		"    alpha-long-name     900ms  15 targets  worst    300ms\n",
		// the name padded to the longest, the number to the widest
		"    b                   100ms   5 targets  worst     50ms\n",
		// one target: singular, and no worst time to compare
		"    one                 500ms   1 target\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestTimingListsTheSlowestTargetsFirst(t *testing.T) {
	t.Run("CGPCH-B80: The slowest targets are listed slowest first, and only when a time was recorded", func(t *testing.T) {})
	out := timingTable(t)
	z, y, x := strings.Index(out, "one @ z.go"), strings.Index(out, "alpha-long-name @ y.go"), strings.Index(out, "b @ x.go")
	if z < 0 || y < 0 || x < 0 || !(z < y && y < x) {
		t.Errorf("the targets must come slowest first (z.go, y.go, x.go):\n%s", out)
	}

	englishOutput(t)
	byGate := map[string]gate.GateSummary{"g": {Gate: "g", Pass: 1}}
	for name, results := range map[string][]gate.Result{
		"no result":        nil,
		"no time recorded": {{Gate: "g", Target: "a.go"}},
	} {
		out := capturaSaida(t, func() { printTiming(gate.Profile{ByGate: byGate, Results: results}) })
		if strings.Contains(out, "slowest targets") {
			t.Errorf("%s: the list of slowest targets must be left out:\n%s", name, out)
		}
	}
}

func TestTimingRoundsToWhatADecisionNeeds(t *testing.T) {
	t.Run("CGPCH-B81: Times are rounded to what a decision needs", func(t *testing.T) {})
	for d, want := range map[time.Duration]string{
		1234567890 * time.Nanosecond: "1.23s",
		1234567 * time.Nanosecond:    "1.2ms",
		1234 * time.Nanosecond:       "1µs",
	} {
		if got := arredonda(d); got != want {
			t.Errorf("%v rounds to %q, want %q", d, got, want)
		}
	}
}

// The race reported from a project running ingests in parallel with the check: the check
// read the map, an ingest wrote its signals, and the check wrote its whole copy back —
// the ingestion disappeared.
func TestRecordCheckKeepsWhatAnotherProcessWrote(t *testing.T) {
	t.Run("CGPCH-B84: The check's stamps do not erase what another process wrote meanwhile", func(t *testing.T) {})
	englishOutput(t)
	issue.UseFiles()
	root := qProject(t, "version: 2\nworkflow:\n  mode: manual\n", nil, nil)
	mapPath := filepath.Join(root, mapx.DefaultPath)
	base := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "s1"}, {ID: "a.go", Kind: mapx.KindCode, Rev: "c1"}},
		Edges: []mapx.Edge{{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies}},
	}
	if err := mapx.Save(base, mapPath); err != nil {
		t.Fatal(err)
	}
	g, err := mapx.Load(mapPath) // what the check confronts
	if err != nil {
		t.Fatal(err)
	}
	if err := mapx.Update(mapPath, func(d *mapx.Graph) error { // an ingest, meanwhile
		d.RecordRunSeconds("a.go", "out/mutation.json", 3)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	p := gate.Profile{Results: []gate.Result{
		{Gate: "g", Target: "a.spec.md", Verdict: gate.Pass, Blocking: true},
		{Gate: "g", Target: "a.go", Verdict: gate.Pass, Blocking: true},
	}}
	captureStdout(t, func() {
		if err := recordCheck(root, mapPath, g, p, false, false); err != nil {
			t.Fatal(err)
		}
	})
	disk, err := mapx.Load(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	if disk.Edges[0].Stamp == nil || disk.Edges[0].Stamp.Verdict != "ok" {
		t.Errorf("the check's stamp reaches the map, got %+v", disk.Edges[0].Stamp)
	}
	for _, n := range disk.Nodes {
		if n.ID == "a.go" && (n.Signal == nil || n.Signal.SecondsBySuite["out/mutation.json"] != 3) {
			t.Errorf("the ingestion written meanwhile is kept, got %+v", n.Signal)
		}
	}
}

func TestCheck_indexOutsideARepository(t *testing.T) {
	t.Run("CGPCH-B85: The index flag reads what the commit records", func(t *testing.T) {})
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte("version: 5\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runQ(t, newCheckCmd(), "--root", dir, "--all", "--no-record", "--index"); err == nil || !strings.Contains(err.Error(), "read the git index") {
		t.Errorf("outside a repository --index fails saying so, got %v", err)
	}
}

const indexYAML = "version: 5\nlayers:\n  code:\n    pattern: \"src/*.ts\"\n    kind: code\ngates:\n  - name: updated-at-atual\n    check: updated-at-atual\n    on: [code]\n    blocking: true\n"

// indexRepo is a project (at `sub` below the repository's top when `below`) with two files
// committed in the past, dated as committed, and a map built from them.
func indexRepo(t *testing.T, below bool) (top, root string, git func(args ...string)) {
	t.Helper()
	top = t.TempDir()
	root = top
	if below {
		root = filepath.Join(top, "sub")
	}
	git = func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", root, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
		c.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-09-01T12:00:00", "GIT_COMMITTER_DATE=2026-09-01T12:00:00")
		if out, err := c.CombinedOutput(); err != nil {
			t.Skipf("git %v: %v %s", args, err, out)
		}
	}
	touchWrite(t, root, "anchors.yaml", indexYAML)
	for _, f := range []string{"src/screen.ts", "src/other.ts"} {
		touchWrite(t, root, f, "// @anchors\n//   updated_at: 2026-09-01\nexport const x = 1\n")
	}
	if err := exec.Command("git", "-C", top, "init", "-q").Run(); err != nil {
		t.Skip("no git")
	}
	git("add", ".")
	git("commit", "-qm", "base")
	cfg, err := config.Load(filepath.Join(root, "anchors.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files, _ := scan.Walk(root, cfg)
	if err := mapx.Save(mapx.Build(files, cfg, gitmeta.AllCommitDates(root)), filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	return top, root, git
}

func updatedAtFails(out string) bool {
	return strings.Contains(out, "✗ [BLOCKS] updated-at-atual")
}

// The case reported from the reference app, simulated in each state a file can be in: the
// commit is judged by what it records, whatever the tree holds.
func TestCheckIndex_updatedAtJudgesTheCommit(t *testing.T) {
	t.Run("CGPCH-B87: Under --index the date is judged by what the commit records, in every staging state", func(t *testing.T) {})
	for _, below := range []bool{false, true} {
		_, root, git := indexRepo(t, below)
		// another session's edit, not staged
		touchWrite(t, root, "src/screen.ts", "// @anchors\n//   updated_at: 2026-09-01\nexport const x = 2 // elsewhere\n")
		touchWrite(t, root, "src/other.ts", "// @anchors\n//   updated_at: "+gitmeta.Today()+"\nexport const x = 3\n")
		git("add", "src/other.ts")

		_, out := runCheckInChild(t, "--root", root, "--changed", "src/screen.ts", "--no-record")
		if !updatedAtFails(out) {
			t.Errorf("below=%v: over the tree, the unstaged edit with an old date fails:\n%s", below, out)
		}
		for _, args := range [][]string{
			{"--changed", "src/screen.ts"},
			{"--all"},
		} {
			_, out := runCheckInChild(t, append([]string{"--root", root, "--index", "--no-record"}, args...)...)
			if updatedAtFails(out) {
				t.Errorf("below=%v %v: under --index the unstaged edit is not the commit's:\n%s", below, args, out)
			}
		}

		// partly staged: the commit itself changes the file, with the old date
		touchWrite(t, root, "src/screen.ts", "// @anchors\n//   updated_at: 2026-09-01\nexport const x = 4 // staged\n")
		git("add", "src/screen.ts")
		touchWrite(t, root, "src/screen.ts", "// @anchors\n//   updated_at: "+gitmeta.Today()+"\nexport const x = 5 // and more, not staged\n")
		_, out = runCheckInChild(t, "--root", root, "--index", "--changed", "src/screen.ts", "--no-record")
		if !updatedAtFails(out) {
			t.Errorf("below=%v: a staged change with the old date fails, whatever the tree says:\n%s", below, out)
		}
	}
}

func TestCheckIndex_noScopeJudgesTheStaged(t *testing.T) {
	t.Run("CGPCH-B86: --index with no scope judges the staged files", func(t *testing.T) {})
	_, root, git := indexRepo(t, false)
	code, out := runCheckInChild(t, "--root", root, "--index", "--no-record")
	if code != 0 || !strings.Contains(out, i18n.T("verify.nothing_staged")) {
		t.Errorf("nothing staged says so, got %d:\n%s", code, out)
	}
	touchWrite(t, root, "src/other.ts", "// @anchors\n//   updated_at: 2026-09-01\nexport const x = 9\n")
	git("add", "src/other.ts")
	_, out = runCheckInChild(t, "--root", root, "--index", "--no-record")
	if !updatedAtFails(out) || !strings.Contains(out, "src/other.ts") {
		t.Errorf("with no scope the staged file is judged:\n%s", out)
	}
}

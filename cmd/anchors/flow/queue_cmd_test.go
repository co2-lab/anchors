package flow

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/co2-lab/anchors/internal/settings"
)

const queueLocalYAML = "version: 1\nlayers:\n" +
	"  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n" +
	"  plan:\n    pattern: \"plans/*.md\"\n    kind: plan\n" +
	"  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n" +
	"docs:\n  required:\n    - kind: openapi\n      path: docs/openapi.yaml\n      trigger: [logic]\n" +
	"gates:\n  - name: informative-one\n    id: INFRM\n    on: [code]\n    blocking: false\n    run: \"true\"\n"

func queueProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", queueLocalYAML)
	return root
}

func runFlowCmd(t *testing.T, cmd interface {
	SetArgs([]string)
	Execute() error
}, args ...string) (string, error) {
	t.Helper()
	cmd.SetArgs(args)
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	return out, err
}

// enqueue puts a task in the queue, creating its file: the queue drops tasks whose file
// no longer exists.
func enqueue(t *testing.T, root string, tk queue.Task) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, tk.Changed)); err != nil {
		writeFile(t, root, tk.Changed, "x\n")
	}
	if tk.CreatedAt == "" {
		tk.CreatedAt = nowStamp()
	}
	if _, err := queue.Enqueue(root, tk); err != nil {
		t.Fatal(err)
	}
}

// `queue` lists what is live, marks the claimed, and gives the hygiene hints.
func TestQueueCmd_listsTheLiveTasks(t *testing.T) {
	root := queueProject(t)
	if out, err := runFlowCmd(t, newQueueCmd(), "--root", root); err != nil || !strings.Contains(out, "empty queue") {
		t.Fatalf("an empty queue says so: %v\n%s", err, out)
	}

	enqueue(t, root, queue.Task{ID: "a-code", Changed: "src/a.ts", Kind: "code", Origin: "watch", SuggestedNext: "feature", Reason: "r1"})
	enqueue(t, root, queue.Task{ID: "b-x", Changed: "odd.bin", Kind: "odd", Origin: "watch", SuggestedNext: "triage", Reason: "r2"})
	if _, err := queue.Claim(root, "w1", nowStamp()); err != nil {
		t.Fatal(err)
	}

	out, err := runFlowCmd(t, newQueueCmd(), "--root", root)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"2 task(s) in the queue", "◐ [claimed]", "by:         w1",
		"○ [pending]", "changed:    odd.bin (odd)", "1 claimed", "1 in 'triage'"} {
		if !strings.Contains(out, want) {
			t.Errorf("the listing lacks %q:\n%s", want, out)
		}
	}
}

// `next` in local mode claims the next task and says how to close it, with the project's
// notifications on top and the reminder of the informative gates at the end.
func TestNextCmd_localClaimsTheNextTask(t *testing.T) {
	root := queueProject(t)
	writeFile(t, root, "notifications.md", "<!-- explain -->\nFreeze on pricing until Monday.\n")
	enqueue(t, root, queue.Task{ID: "a-code", Changed: "src/a.ts", Kind: "code", Origin: "watch", SuggestedNext: "feature", Reason: "the code changed"})

	out, err := runFlowCmd(t, newNextCmd(), "--root", root, "--worker", "w9")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Freeze on pricing until Monday.", "task claimed: a-code",
		"suggestion: feature", "anchors done a-code", "1 informative gate(s) declared"} {
		if !strings.Contains(out, want) {
			t.Errorf("the output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "explain") {
		t.Errorf("HTML comments of the notification are not printed:\n%s", out)
	}
	if strings.Index(out, "Freeze") > strings.Index(out, "task claimed") {
		t.Error("the notification comes before the task")
	}
	tasks, _ := queue.List(root)
	if len(tasks) != 1 || tasks[0].State != queue.Claimed || tasks[0].ClaimedBy != "w9" {
		t.Errorf("the task must be claimed by the worker: %+v", tasks)
	}
}

// A cold start: an empty queue with a plan that still has specs to be born is seeded with
// that plan — one plan, however many are pending.
func TestNextCmd_emptyQueueSeedsFromThePlans(t *testing.T) {
	root := queueProject(t)
	writeFile(t, root, "plans/0001-a.md", "- [ ] `src/A.spec.md` — a\n- [ ] `src/B.spec.md` — b\n")
	writeFile(t, root, "plans/0002-b.md", "- [ ] `src/C.spec.md` — c\n")
	writeFile(t, root, "plans/0003-done.md", "- [x] `src/D.spec.md` — d\n")
	writeFile(t, root, "src/A.spec.md", "# a\n")
	writeFile(t, root, "src/D.spec.md", "# d\n")

	out, err := runFlowCmd(t, newNextCmd(), "--root", root, "--worker", "w1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "seeded with 1 plan(s)") || !strings.Contains(out, "task claimed:") {
		t.Fatalf("the cold start seeds and claims:\n%s", out)
	}
	if !strings.Contains(out, "1 of 2 spec(s) of this plan do not exist yet") {
		t.Errorf("the seed's tally is recomputed from the disk:\n%s", out)
	}
	tasks, _ := queue.List(root)
	if len(tasks) != 1 || tasks[0].Kind != "plan" || tasks[0].Origin != "seed" {
		t.Errorf("exactly one plan is seeded: %+v", tasks)
	}

	// nothing left to seed once every plan's specs exist
	root2 := queueProject(t)
	writeFile(t, root2, "plans/0003-done.md", "- [x] `src/D.spec.md` — d\n")
	writeFile(t, root2, "src/D.spec.md", "# d\n")
	if out, _ := runFlowCmd(t, newNextCmd(), "--root", root2); !strings.Contains(out, "empty queue — nothing to do") {
		t.Errorf("a fulfilled plan seeds nothing:\n%s", out)
	}
}

// A seed cited by NAME counts when the name is unique, and not when two files share it.
func TestSeedExists_byNameOnlyWhenUnique(t *testing.T) {
	root := t.TempDir()
	if seedExists(root, "Pricing.spec.md", nil) {
		t.Error("nothing on disk")
	}
	files := scanFiles("a/Pricing.spec.md")
	if !seedExists(root, "Pricing.spec.md", files) {
		t.Error("a unique name is found")
	}
	if seedExists(root, "Pricing.spec.md", scanFiles("a/Pricing.spec.md", "b/Pricing.spec.md")) {
		t.Error("two files with the name make the citation ambiguous")
	}
}

// `done` closes by id or in batch by file, kind or all.
func TestDoneCmd_closesByIdAndInBatch(t *testing.T) {
	root := queueProject(t)
	for _, tk := range []queue.Task{
		{ID: "t1", Changed: "src/a.ts", Kind: "code"},
		{ID: "t2", Changed: "src/a.ts", Kind: "feature"},
		{ID: "t3", Changed: "src/b.ts", Kind: "code"},
		{ID: "t4", Changed: "src/c.spec.md", Kind: "spec"},
		{ID: "t5", Changed: "src/d.ts", Kind: "test"},
	} {
		tk.SuggestedNext = "next-of-" + tk.ID // the queue dedups same file + same step
		enqueue(t, root, tk)
	}
	ids := func() string {
		ts, _ := queue.List(root)
		var s []string
		for _, tk := range ts {
			s = append(s, tk.ID)
		}
		return strings.Join(s, ",")
	}

	if _, err := runFlowCmd(t, newDoneCmd(), "--root", root); err == nil || !strings.Contains(err.Error(), "provide the <id>") {
		t.Errorf("no id and no filter is an error, got %v", err)
	}
	if out, err := runFlowCmd(t, newDoneCmd(), "--root", root, "t4"); err != nil || !strings.Contains(out, "task done: t4") {
		t.Fatalf("by id: %v\n%s", err, out)
	}
	if out, _ := runFlowCmd(t, newDoneCmd(), "--root", root, "--file", "src/a.ts"); !strings.Contains(out, "2 task(s) done") {
		t.Errorf("by file closes both tasks of src/a.ts:\n%s", out)
	}
	if got := ids(); got != "t3,t5" {
		t.Errorf("left: %s", got)
	}
	if out, _ := runFlowCmd(t, newDoneCmd(), "--root", root, "--kind", "spec"); !strings.Contains(out, "no task matched the filter") {
		t.Errorf("no match is said:\n%s", out)
	}
	if out, _ := runFlowCmd(t, newDoneCmd(), "--root", root, "--kind", "code"); !strings.Contains(out, "1 task(s) done") || ids() != "t5" {
		t.Errorf("by kind closes only code:\n%s", out)
	}
	if _, err := runFlowCmd(t, newDoneCmd(), "--root", root, "--all"); err != nil || ids() != "" {
		t.Errorf("--all empties the queue: %v, left %s", err, ids())
	}
	if _, err := runFlowCmd(t, newDoneCmd(), "--root", root, "nope"); err == nil {
		t.Error("an unknown id is an error")
	}
}

func TestDropCmd_removesWithoutArchiving(t *testing.T) {
	root := queueProject(t)
	enqueue(t, root, queue.Task{ID: "junk", Changed: "x", Kind: "odd", SuggestedNext: "triage"})
	out, err := runFlowCmd(t, newDropCmd(), "--root", root, "junk")
	if err != nil || !strings.Contains(out, "task discarded: junk") {
		t.Fatalf("drop: %v\n%s", err, out)
	}
	if ts, _ := queue.List(root); len(ts) != 0 {
		t.Errorf("the task must be gone: %+v", ts)
	}
	if _, err := runFlowCmd(t, newDropCmd(), "--root", root, "junk"); err == nil {
		t.Error("dropping what is not there is an error")
	}
}

// A task claimed a moment ago stays with its worker, and the zero explains itself; --force
// takes it back anyway.
func TestReclaimCmd_respectsTheLiveWorkerUnlessForced(t *testing.T) {
	root := queueProject(t)
	enqueue(t, root, queue.Task{ID: "t1", Changed: "src/a.ts", Kind: "code", SuggestedNext: "feature"})
	if _, err := queue.Claim(root, "w1", nowStamp()); err != nil {
		t.Fatal(err)
	}
	out, err := runFlowCmd(t, newReclaimCmd(), "--root", root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "0 task(s) returned") || !strings.Contains(out, "1 task(s) claimed RECENTLY stayed") {
		t.Errorf("the zero must explain itself:\n%s", out)
	}
	out, err = runFlowCmd(t, newReclaimCmd(), "--root", root, "--force")
	if err != nil || !strings.Contains(out, "1 task(s) returned") {
		t.Fatalf("--force releases the live claim: %v\n%s", err, out)
	}
	if ts, _ := queue.List(root); len(ts) != 1 || ts[0].State != queue.Pending {
		t.Errorf("the task is pending again: %+v", ts)
	}
}

func TestDefaultWorkerID_isPidAtHost(t *testing.T) {
	id := defaultWorkerID()
	if !regexp.MustCompile(`^\d+@.+$`).MatchString(id) || !strings.HasPrefix(id, strconv.Itoa(os.Getpid())+"@") {
		t.Errorf("got %q", id)
	}
}

// github mode: the queue is the board. Without a session the claim is refused; with one,
// the agent's own card is resumed before anything is asked of the pipeline, with the
// notifications of the integration branch on top.
func TestNextCmd_githubResumesTheAgentsCard(t *testing.T) {
	root := githubProject(t)
	writeFile(t, root, "anchors.graph.yaml", "version: 4\nnodes:\n  - id: src/pricing.ts\n    kind: code\n    rev: a\n    code: PRICX\nedges: []\n")
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	calls := scriptedGH(t,
		ghRule{match: "api -H Accept: application/vnd.github.raw *notifications.md*", out: "Release on Friday."},
		ghRule{match: "api graphql*", out: `{"number":42,"title":"[PRICX] Implementar spec — pricing",` +
			`"body":"Unidade: ` + "`src`" + `\nCódigo: ` + "`PRICX`" + `",` +
			`"labels":[{"name":"anchors"},{"name":"anchors:in-progress"}],"comments":[{"body":"anchors-owner: ` + host + `/dev3"}]}`},
	)

	t.Setenv("ANCHORS_SESSION", "")
	if _, err := runFlowCmd(t, newNextCmd(), "--root", root); err == nil || !strings.Contains(err.Error(), "ANCHORS_SESSION is not set") {
		t.Fatalf("the board is not claimed without a session, got %v", err)
	}

	t.Setenv("ANCHORS_SESSION", "dev3")
	out, err := runFlowCmd(t, newNextCmd(), "--root", root)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Release on Friday.", "(no terminal to ask", "resuming your card",
		"card claimed: #42", "state:    in-progress", "owner:    " + host + "/dev3",
		"1. anchors work code    --for src/pricing.ts", "YOU DO NOT DECIDE THE DIRECTION",
		"anchors pr-body --cards 42"} {
		if !strings.Contains(out, want) {
			t.Errorf("the output lacks %q:\n%s", want, out)
		}
	}
	if len(callsWith(calls(), "workflow run")) != 0 {
		t.Error("with a card in hand nothing is asked of the pipeline")
	}
	// the non-interactive run records nothing: the role stays undeclared
	if s, _ := settings.Load(root); s.Decided() {
		t.Errorf("no terminal, no answer to record: %+v", s)
	}
}

func TestNextFromBoard_needsTheRepo(t *testing.T) {
	cfg := cfgGitHub()
	cfg.Workflow.Repo = ""
	if err := nextFromBoard(t.TempDir(), cfg, "h/a"); err == nil || !strings.Contains(err.Error(), "workflow.repo is empty") {
		t.Errorf("got %v", err)
	}
}

// Each end of a claim without a card says which run to follow, and only a failed run is
// an error.
func TestReportClaimWithoutCard(t *testing.T) {
	var err error
	out := stdoutOf(t, func() { err = reportClaimWithoutCard(board.ClaimOutcome{TimedOut: true}) })
	if err != nil || !strings.Contains(out, "did not finish within") || !strings.Contains(out, "gh run list --workflow "+board.ClaimWorkflow) {
		t.Errorf("timed out without a run: %v\n%s", err, out)
	}
	out = stdoutOf(t, func() {
		err = reportClaimWithoutCard(board.ClaimOutcome{Run: &board.ClaimRun{ID: 77, Conclusion: "success"}})
	})
	if err != nil || !strings.Contains(out, "no free card on the board") || !strings.Contains(out, "gh run view 77 --log") {
		t.Errorf("a successful empty claim is not an error: %v\n%s", err, out)
	}
	err = reportClaimWithoutCard(board.ClaimOutcome{Run: &board.ClaimRun{ID: 78, Conclusion: "failure"}})
	if err == nil || !strings.Contains(err.Error(), `#78 ended as "failure"`) {
		t.Errorf("a failed claim is an error, got %v", err)
	}
	if err := reportClaimWithoutCard(board.ClaimOutcome{}); err != nil {
		t.Errorf("nothing to report, got %v", err)
	}
}

// The card's state decides the deliverable, and the title names it.
func TestPrintBoardWork_theDeliverableFollowsTheCard(t *testing.T) {
	root := queueProject(t)
	writeFile(t, root, "anchors.graph.yaml", "version: 4\nnodes:\n  - id: src/pricing.ts\n    kind: code\n    rev: a\n    code: PRICX\n    layer: logic\nedges: []\n")
	writeFile(t, root, "src/pricing.ts", "x\n")

	plan := stdoutOf(t, func() {
		printBoardWork(root, &board.Card{Number: 5, Title: "[plano] Implementar plan — F02"}, "h/a")
	})
	if !strings.Contains(plan, "DELIVERABLE: every spec the plan seeds.") || strings.Contains(plan, "What to propagate") {
		t.Errorf("a plan card delivers its specs, with no unit to propagate:\n%s", plan)
	}

	spec := stdoutOf(t, func() {
		printBoardWork(root, &board.Card{Number: 6, Title: "[PRICX] Implementar spec — pricing", Body: "Código: `PRICX`"}, "h/a")
	})
	for _, want := range []string{"DELIVERABLE: code + feature + test + documentation.",
		"3. anchors work test    --for src/pricing.ts", "Changing `src/pricing.ts` REQUIRES touching:",
		"docs/openapi.yaml", "anchors docs duties --layer logic", "What to propagate: anchors impact src/pricing.ts"} {
		if !strings.Contains(spec, want) {
			t.Errorf("the spec card lacks %q:\n%s", want, spec)
		}
	}

	// an unknown code falls back to the folder of the body
	other := stdoutOf(t, func() {
		printBoardWork(root, &board.Card{Number: 7, Title: "something else", Body: "Unidade: `packages/infra, packages/x`\nCódigo: `GONE1`"}, "h/a")
	})
	if !strings.Contains(other, "DELIVERABLE: see the body of the card.") || !strings.Contains(other, "anchors impact packages/infra") {
		t.Errorf("the folder is the fallback target:\n%s", other)
	}

	// the product owner is not told to escalate what they decide
	if err := settings.Save(root, settings.Settings{Role: settings.RolePO}); err != nil {
		t.Fatal(err)
	}
	po := stdoutOf(t, func() { printBoardWork(root, &board.Card{Number: 8, Title: "x"}, "h/a") })
	if strings.Contains(po, "YOU DO NOT DECIDE") {
		t.Errorf("whoever decides the product is not told to escalate:\n%s", po)
	}

	// a card in review gets the review deliverable
	rev := stdoutOf(t, func() {
		printBoardWork(root, &board.Card{Number: 9, Title: "x", State: board.StateInReview, Body: "Código: `PRICX`"}, "h/a")
	})
	if !strings.Contains(rev, "DELIVERABLE: the REVIEW") || !strings.Contains(rev, "anchors work review --for src/pricing.ts") {
		t.Errorf("a card in review is reviewed:\n%s", rev)
	}
}

// Documentation duties are silent where the project declares none.
func TestPrintDocDuties_silentWithoutDuties(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "version: 1\nlayers:\n  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n")
	if out := stdoutOf(t, func() { printDocDuties(root, "src/a.ts") }); out != "" {
		t.Errorf("no duties, nothing printed: %q", out)
	}
	if out := stdoutOf(t, func() { printDocDuties(t.TempDir(), "src/a.ts") }); out != "" {
		t.Errorf("no config, nothing printed: %q", out)
	}
	if layerOfPath(root, nil, "") != "" {
		t.Error("an empty path has no layer")
	}
}

// Whoever cannot be asked is not asked: a pipe and /dev/null are not a terminal.
func TestInteractiveTerminal_pipeAndDevNullAreNotATerminal(t *testing.T) {
	prev := os.Stdin
	defer func() { os.Stdin = prev }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	os.Stdin = r
	if interactiveTerminal() {
		t.Error("a pipe is not a terminal")
	}

	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer null.Close()
	os.Stdin = null
	if interactiveTerminal() {
		t.Error("/dev/null is a char device and still not a terminal")
	}
}

// A declared decision is not asked again.
func TestEnsureLocalDecision_keepsWhatWasDeclared(t *testing.T) {
	root := t.TempDir()
	if err := settings.Save(root, settings.Settings{Role: settings.RoleDev}); err != nil {
		t.Fatal(err)
	}
	var s settings.Settings
	var err error
	out := stdoutOf(t, func() { s, err = ensureLocalDecision(root) })
	if err != nil || s.Role != settings.RoleDev || out != "" {
		t.Errorf("a declared role is kept silently: %+v %v %q", s, err, out)
	}
	if decidesProduct(root) {
		t.Error("a dev does not decide the product")
	}
}

func scanFiles(paths ...string) []scan.File {
	var out []scan.File
	for _, p := range paths {
		out = append(out, scan.File{Path: p})
	}
	return out
}

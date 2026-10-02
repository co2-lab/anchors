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
	t.Run("WRQUW-B01: queue lists the live tasks with the hygiene hints", func(t *testing.T) {})
	t.Run("WRQUW-X01: queue claims nothing", func(t *testing.T) {})
	root := queueProject(t)
	if out, err := runFlowCmd(t, newQueueCmd(), "--root", root); err != nil || !strings.Contains(out, "empty queue") {
		t.Fatalf("an empty queue says so: %v\n%s", err, out)
	}

	enqueue(t, root, queue.Task{ID: "a-code", Changed: "src/a.ts", Kind: "code", Origin: "watch", SuggestedNext: "feature", Reason: "r1"})
	enqueue(t, root, queue.Task{ID: "b-x", Changed: "odd.bin", Kind: "odd", Origin: "watch", SuggestedNext: "triage", Reason: "r2"})
	if _, err := queue.Claim(root, "w1", nowStamp()); err != nil {
		t.Fatal(err)
	}

	states := func() string {
		ts, _ := queue.List(root)
		var s []string
		for _, tk := range ts {
			s = append(s, tk.ID+"="+string(tk.State))
		}
		return strings.Join(s, ",")
	}
	before := states()
	out, err := runFlowCmd(t, newQueueCmd(), "--root", root)
	if err != nil {
		t.Fatal(err)
	}
	if after := states(); after != before {
		t.Errorf("queue must claim nothing: before %s, after %s", before, after)
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
	t.Run("WRQUW-B02: next claims the next task in local mode", func(t *testing.T) {})
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
	t.Run("WRQUW-B03: A cold start seeds a plan that still has work", func(t *testing.T) {})
	t.Run("WRQUW-I01: A cold start seeds one plan", func(t *testing.T) {})
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
	t.Run("WRQUW-B04: A cited spec exists by path or by a unique name", func(t *testing.T) {})
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
	t.Run("WRQUW-B06: done closes by id and in batch", func(t *testing.T) {})
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
	t.Run("WRQUW-B07: drop deletes a task without archiving it", func(t *testing.T) {})
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
	t.Run("WRQUW-B08: reclaim respects a live worker unless forced", func(t *testing.T) {})
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
	t.Run("WRQUW-B09: The default worker is pid at host", func(t *testing.T) {})
	id := defaultWorkerID()
	if !regexp.MustCompile(`^\d+@.+$`).MatchString(id) || !strings.HasPrefix(id, strconv.Itoa(os.Getpid())+"@") {
		t.Errorf("got %q", id)
	}
}

// github mode: the queue is the board. Without a session the claim is refused; with one,
// the agent's own card is resumed before anything is asked of the pipeline, with the
// notifications of the integration branch on top.
func TestNextCmd_githubResumesTheAgentsCard(t *testing.T) {
	t.Run("WRQUW-B12: The agent's own card is resumed before asking the pipeline", func(t *testing.T) {})
	t.Run("WRQUW-B14: The claimed card is printed with its state and owner", func(t *testing.T) {})
	t.Run("WRQUW-B19: Other cards end naming the pull request body command", func(t *testing.T) {})
	root := githubProject(t)
	writeFile(t, root, "anchors.graph.yaml", "version: 6\nnodes:\n  - id: src/pricing.ts\n    kind: code\n    rev: a\n    code: PRICX\nedges: []\n")
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
		"anchors pr-body --cards 42", "brings the `Refs` that links the card; it stays open"} {
		if !strings.Contains(out, want) {
			t.Errorf("the output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Closes") {
		t.Errorf("pr-body links with Refs and does not close the card:\n%s", out)
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
	t.Run("WRQUW-X02: The board is not claimed without a repository", func(t *testing.T) {})
	cfg := cfgGitHub()
	cfg.Workflow.Repo = ""
	if err := nextFromBoard(t.TempDir(), cfg, "h/a"); err == nil || !strings.Contains(err.Error(), "workflow.repo is empty") {
		t.Errorf("got %v", err)
	}
}

// Each end of a claim without a card says which run to follow, and only a failed run is
// an error.
func TestReportClaimWithoutCard(t *testing.T) {
	t.Run("WRQUW-B13: A claim without a card names the run to follow", func(t *testing.T) {})
	t.Run("WRQUW-E01: A failed claim run is an error", func(t *testing.T) {})
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
	t.Run("WRQUW-B15: The deliverable follows the card's title", func(t *testing.T) {})
	t.Run("WRQUW-B17: Only an agent that does not decide the product is told to escalate", func(t *testing.T) {})
	root := queueProject(t)
	writeFile(t, root, "anchors.graph.yaml", "version: 6\nnodes:\n  - id: src/pricing.ts\n    kind: code\n    rev: a\n    code: PRICX\n    layer: logic\nedges: []\n")
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
	t.Run("WRQUW-B16: The documentation duties follow the unit", func(t *testing.T) {})
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
	t.Run("WRQUW-B20: The role is asked at most once, and never without a terminal", func(t *testing.T) {})
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

// THE NUMBER MUST NOT LIVE IN THE `reason`.
//
// It is derivable (the seeds are recounted from the disk), and storing it made the text age
// in the queue: the task was born saying "6 of 7", two specs were delivered, and it kept
// saying 6 (co2-lab/anchors#11) — a stale reason sent a maintainer chasing a defect that
// did not exist.
func TestSeedTally(t *testing.T) {
	t.Run("WRQUW-B05: The seed count is recomputed when printed", func(t *testing.T) {})
	build := func(t *testing.T, existing ...string) string {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
			t.Fatal(err)
		}
		yaml := "version: 1\nlayers:\n" +
			"  plan:\n    pattern: \"plans/*.md\"\n    kind: plan\n" +
			"  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"
		if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
			t.Fatal(err)
		}
		plan := "- [ ] `src/A.spec.md` — a\n- [ ] `src/B.spec.md` — b\n- [ ] `src/C.spec.md` — c\n"
		if err := os.WriteFile(filepath.Join(root, "plans", "0001-x.md"), []byte(plan), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, f := range existing {
			p := filepath.Join(root, f)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte("# spec\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	task := queue.Task{
		Changed: "plans/0001-x.md", Kind: "plan", Origin: "seed",
		Reason: "a plan was seeded — generate the specs it lists",
	}

	t.Run("counts what is missing NOW, not at creation", func(t *testing.T) {
		root := build(t)
		if got := seedTally(root, task); !strings.Contains(got, "3 of 3") {
			t.Errorf("with zero specs on disk, want \"3 of 3\"; got %q", got)
		}
		// the SAME task, after two deliveries
		root = build(t, "src/A.spec.md", "src/B.spec.md")
		if got := seedTally(root, task); !strings.Contains(got, "1 of 3") {
			t.Errorf("with two specs delivered, want \"1 of 3\"; got %q", got)
		}
	})

	// An OLD task already carries the count in its reason — adding the recomputed one would
	// print the number twice. Detected by the phrase, not by version: the field does not
	// record who wrote it, and a task survives any number of binary updates.
	t.Run("an old task with a stored count is not doubled", func(t *testing.T) {
		root := build(t)
		old := task
		old.Reason += " — 6 of 7 spec(s) of this plan do not exist yet"
		if got := seedTally(root, old); got != "" {
			t.Errorf("it would double the count: %q", got)
		}
	})

	// A judgment reason does not age — the question is the gate's, not the state's.
	t.Run("only a seed task gets a count", func(t *testing.T) {
		root := build(t)
		for _, other := range []queue.Task{
			{Changed: "x.spec.md", Kind: "judgment", Origin: "check", Reason: "gate 'y' — question: ..."},
			{Changed: "plans/0001-x.md", Kind: "plan", Origin: "human", Reason: "someone asked"},
		} {
			if got := seedTally(root, other); got != "" {
				t.Errorf("kind=%s origin=%s got a count: %q", other.Kind, other.Origin, got)
			}
		}
	})
}

// A review card whose body names no unit (a plan card, a finding under another card)
// printed `anchors work review --for ` with an empty target — seen in the reference app on #931.
// It now points to the PR that references the card.
func TestPrintReviewWork_cardWithNoUnitPointsToItsPR(t *testing.T) {
	t.Run("WRQUW-B18: A card under review asks for a verdict", func(t *testing.T) {})
	out := stdoutOf(t, func() {
		printReviewWork(t.TempDir(), &board.Card{Number: 931, Title: "[plano] x", Body: "no unit here"}, "host/dev3")
	})
	if strings.Contains(out, "--for \n") || strings.Contains(out, "--for  ") {
		t.Errorf("the review command still has an empty target:\n%s", out)
	}
	if !strings.Contains(out, `#931 in:body`) {
		t.Errorf("the output does not point to the card's PR:\n%s", out)
	}
}

// A review ends with the verdict line on the PR, and `anchors next` is the only place the
// reviewer is told so (reference app: reviewers ended with "Veredito: OK",
// `judge`, `decided`, and the card kept coming back). The output names the exact line for
// this agent, and no longer sends a reviewer to open a PR of their own.
func TestPrintReviewWork_endsWithTheVerdictLine(t *testing.T) {
	out := stdoutOf(t, func() {
		printReviewWork(t.TempDir(), &board.Card{Number: 869, Title: "[plano] x", Body: "no unit here"}, "host/dev3")
	})
	for _, want := range []string{"anchors-review: approved by host/dev3", "anchors-review: rejected by host/dev3", "releases you"} {
		if !strings.Contains(out, want) {
			t.Errorf("the review work should say %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "anchors pr-body") {
		t.Errorf("a reviewer is not told to open a PR of their own:\n%s", out)
	}
}

// Without ANCHORS_SESSION the identity falls back to the OS user, and two agents of that
// user collide (reference app). The fallback is announced; a declared session is not.
func TestAgentID_fallbackIsAnnounced(t *testing.T) {
	t.Run("WRQUW-B11: The board identity falls back to the OS user, and says so", func(t *testing.T) {})
	t.Setenv("ANCHORS_SESSION", "")
	t.Setenv("USER", "alice")
	var id string
	out := stderrOf(t, func() { id = agentID() })
	if !strings.HasSuffix(id, "/alice") {
		t.Fatalf("fallback identity = %q, want the OS user", id)
	}
	if !strings.Contains(out, "ANCHORS_SESSION is not set") {
		t.Errorf("the fallback identity was used silently:\n%q", out)
	}

	t.Setenv("USER", "")
	t.Setenv("USERNAME", "bob")
	if id = agentID(); !strings.HasSuffix(id, "/bob") {
		t.Errorf("with no USER, the Windows USERNAME, got %q", id)
	}

	t.Setenv("ANCHORS_SESSION", "devA")
	out = stderrOf(t, func() { id = agentID() })
	if !strings.HasSuffix(id, "/devA") || out != "" {
		t.Errorf("a declared session: id=%q, stderr=%q (want no warning)", id, out)
	}
}

// Claiming from the board requires a declared session: the fallback identity is shared by
// every agent of the same user on one machine (reference app).
func TestRequireSession(t *testing.T) {
	t.Run("WRQUW-B10: The board is not claimed without a session", func(t *testing.T) {})
	t.Setenv("ANCHORS_SESSION", "")
	err := requireSession()
	if err == nil || !strings.Contains(err.Error(), "export ANCHORS_SESSION=") {
		t.Errorf("claiming without ANCHORS_SESSION was allowed, or the error does not say how to fix it: %v", err)
	}
	t.Setenv("ANCHORS_SESSION", "devA")
	if err := requireSession(); err != nil {
		t.Errorf("a declared session was refused: %v", err)
	}
}

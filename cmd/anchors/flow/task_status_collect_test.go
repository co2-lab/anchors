package flow

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/telemetry"
)

// The round's state is discovered, not typed: the card and its label, the people-bound
// decisions, the PR and its checks, what the lock reverted, and the working tree.
func TestCollectTaskState_discoversWhatTheMachineKnows(t *testing.T) {
	root := githubProject(t)
	gitRepo(t, root)
	writeFile(t, root, "src/dirty.ts", "export const x = 1\n")

	scriptedGH(t,
		ghRule{match: "issue view 303 --json number,title,labels,state*",
			out: `{"number":303,"title":"[INDTN] implement spec","state":"OPEN",` +
				`"labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]}`},
		ghRule{match: "issue list --label anchors:needs-user*",
			out: `[{"number":701,"title":"which currency"},{"number":702,"title":"which region"}]`},
		ghRule{match: "pr view main --json number,state,statusCheckRollup*",
			out: `{"number":368,"state":"OPEN","statusCheckRollup":[` +
				`{"conclusion":"SUCCESS","status":"COMPLETED"},{"conclusion":"SKIPPED","status":"COMPLETED"},` +
				`{"conclusion":"FAILURE","status":"COMPLETED"}]}`},
		ghRule{match: "issue view 303 --json comments*",
			out: `{"comments":[` +
				`{"body":"` + initx.MarcadorDeReversao + ` **Reverted: closed by hand** (by ` + "`bob`" + `)\n\nlong rule text","author":{"login":"github-actions"}},` +
				`{"body":"` + initx.MarcadorDeReversao + ` I locked this","author":{"login":"bob"}}]}`},
	)
	cfg, err := config.Load(root + "/anchors.yaml")
	if err != nil {
		t.Fatal(err)
	}

	e := collectTaskState(root, cfg, 303)

	if e.Branch != "main" {
		t.Errorf("branch: got %q", e.Branch)
	}
	if e.Clean {
		t.Error("an untracked file makes the tree NOT clean")
	}
	if e.Card == nil || e.Card.Number != 303 || e.Card.State != "anchors:in-progress" {
		t.Fatalf("the card and its state label: got %+v", e.Card)
	}
	if len(e.Blocked) != 2 || e.Blocked[1].Number != 702 || e.Blocked[1].Title != "which region" {
		t.Errorf("the needs-user cards: got %+v", e.Blocked)
	}
	if e.PR == nil || e.PR.Number != 368 || e.PR.Total != 3 ||
		e.PR.Checks["passou"] != 2 || e.PR.Checks["reprovou"] != 1 {
		t.Errorf("the PR and its verdict: got %+v", e.PR)
	}
	// Only the bot's reversal counts, cut to its first line and without markup.
	if len(e.Reverted) != 1 || e.Reverted[0] != "Reverted: closed by hand (by bob)" {
		t.Errorf("reversals: got %q", e.Reverted)
	}
}

// A CLOSED card has no state label worth telling: the leftover label would lie.
func TestCardByNumber_closedCardSaysClosed(t *testing.T) {
	scriptedGH(t, ghRule{match: "issue view 9 *",
		out: `{"number":9,"title":"x","state":"CLOSED","labels":[{"name":"anchors:ready-to-review"}]}`})
	c := cardByNumber(boardClientFor(), 9)
	if c == nil || c.State != "closed" {
		t.Fatalf("a closed card must read as closed, got %+v", c)
	}
	if len(c.Labels) != 1 || c.Labels[0] != "anchors:ready-to-review" {
		t.Errorf("the labels are still reported: %v", c.Labels)
	}
}

// Every source fails silently: a partial report beats none.
func TestCollectTaskState_everySourceFailingYieldsAnEmptyReport(t *testing.T) {
	scriptedGH(t, ghRule{match: "*", out: "not json"})
	if c := cardByNumber(boardClientFor(), 9); c != nil {
		t.Errorf("unreadable card: got %+v", c)
	}
	if b := escalatedCards(boardClientFor()); b != nil {
		t.Errorf("unreadable list: got %+v", b)
	}
	if p := currentBranchPR(t.TempDir(), "acme/app"); p != nil {
		t.Errorf("unreadable PR: got %+v", p)
	}
	if r := revertedOn("acme/app", 9); r != nil {
		t.Errorf("unreadable comments: got %+v", r)
	}

	scriptedGH(t, ghRule{match: "*", code: 1})
	e := collectTaskState(t.TempDir(), &config.Config{}, 0)
	if e.Card != nil || e.PR != nil || e.Blocked != nil || e.Branch != "" {
		t.Errorf("with no source answering, nothing is invented: %+v", e)
	}
}

// Without --card, the card is the one the board says this agent owns.
func TestTaskStatusCmd_findsTheAgentsCardOnTheBoard(t *testing.T) {
	root := githubProject(t)
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	t.Setenv("ANCHORS_SESSION", "dev7")
	scriptedGH(t,
		ghRule{match: "api graphql*",
			out: `{"number":41,"title":"[ABCDE] other","body":"","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}],"comments":[{"body":"anchors-owner: someone/else"}]}` + "\n" +
				`{"number":42,"title":"[ABCDE] mine","body":"","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}],"comments":[{"body":"anchors-owner: ` + host + `/dev7"}]}`},
	)
	cmd := newTaskStatusCmd()
	cmd.SetArgs([]string{"--root", root})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Task  #42 · mine") {
		t.Errorf("the report must be about the agent's own card:\n%s", out)
	}
}

// The turn-ended event carries numbers and vocabulary only: the state without its
// prefix, and the check counts.
func TestEmitTurnEnded_sendsTheStateTheTurnEndedIn(t *testing.T) {
	got := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got <- string(b)
	}))
	defer srv.Close()
	prev := common.Emitter
	common.Emitter = telemetry.NewEmitter(telemetry.Config{Enabled: true, Endpoint: srv.URL}, "test")
	defer func() { common.Emitter = prev }()

	emitTurnEnded(taskState{
		Card: testCard(303, "anchors:in-review", "secret title"),
		PR:   &branchPR{Number: 368, State: "OPEN", Total: 4, Checks: map[string]int{"reprovou": 3}},
	})
	common.Emitter.Flush()

	body := <-got
	for _, want := range []string{"in-review", "checks_reprovaram", "open"} {
		if !strings.Contains(body, want) {
			t.Errorf("the event lacks %q: %s", want, body)
		}
	}
	for _, leak := range []string{"secret title", "anchors:in-review"} {
		if strings.Contains(body, leak) {
			t.Errorf("the event must not carry %q: %s", leak, body)
		}
	}
}

func boardClientFor() board.Client {
	return board.Client{Repo: "acme/app", Labels: []string{"anchors"}}
}

// A check still running has no conclusion yet. It is "em curso", not a failure: counting
// it as failed makes task-status tell the agent to fix a PR whose CI simply has not
// finished.
func TestCurrentBranchPR_runningChecksAreInProgressNotFailed(t *testing.T) {
	root := t.TempDir()
	gitRepo(t, root)
	scriptedGH(t, ghRule{match: "pr view *",
		out: `{"number":368,"state":"OPEN","statusCheckRollup":[` +
			`{"conclusion":"SUCCESS","status":"COMPLETED"},` +
			`{"conclusion":"","status":"IN_PROGRESS"},` +
			`{"conclusion":"","status":"QUEUED"},` +
			`{"conclusion":"FAILURE","status":"COMPLETED"},` +
			`{"state":"PENDING"},{"state":"SUCCESS"}]}`})

	p := currentBranchPR(root, "acme/app")
	if p == nil {
		t.Fatal("the PR must be read")
	}
	if p.Checks["em curso"] != 3 || p.Checks["reprovou"] != 1 || p.Checks["passou"] != 2 || p.Total != 6 {
		t.Errorf("3 running (one a pending commit status), 1 failed, 2 passed: got %+v", p.Checks)
	}
	steps := strings.Join(nextStep(taskState{Clean: true, PR: p}), "\n")
	if !strings.Contains(steps, "is running") || strings.Contains(steps, "failed") {
		t.Errorf("a running CI must be waited on, not fixed:\n%s", steps)
	}
}

// task-status reads the CONFIGURED repo, never the one the working directory points at,
// and the branch of --root, never the branch of the cwd. From a fork, or with --root
// elsewhere, anything else reports another repo's card and PR.
func TestCollectTaskState_readsTheConfiguredRepoAndTheRootsBranch(t *testing.T) {
	root := githubProject(t)
	gitRepo(t, root)
	c := exec.Command("git", "checkout", "-q", "-b", "feat-303")
	c.Dir = root
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("git checkout: %s", out)
	}
	calls := scriptedGH(t,
		ghRule{match: "issue view 303 --json number,title,labels,state*",
			out: `{"number":303,"title":"x","state":"OPEN","labels":[{"name":"anchors:in-progress"}]}`},
		ghRule{match: "pr view *", out: `{"number":368,"state":"OPEN","statusCheckRollup":[]}`},
	)
	cfg, err := config.Load(root + "/anchors.yaml")
	if err != nil {
		t.Fatal(err)
	}

	collectTaskState(root, cfg, 303)

	got := calls()
	for _, want := range [][]string{
		{"issue view 303 ", "--json number,title,labels,state"},
		{"issue list ", "anchors:needs-user"},
		{"pr view ", "statusCheckRollup"},
		{"issue view 303 ", "--json comments"},
	} {
		hits := callsWith(got, want...)
		if len(hits) != 1 || !strings.Contains(hits[0], "--repo acme/app") {
			t.Errorf("%v must name the configured repo; calls: %v", want, got)
		}
	}
	if pr := callsWith(got, "pr view "); len(pr) != 1 || !strings.Contains(pr[0], "pr view feat-303 ") {
		t.Errorf("the PR looked up must be the one of --root's branch; calls: %v", got)
	}
}

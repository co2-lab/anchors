// @anchors
//   code: TSTTS
//   ref: TSSTT

package flow

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
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
	t.Run("TSSTT-B01: The working tree is read from the root", func(t *testing.T) {})
	t.Run("TSSTT-B02: A card given by number carries its state, and a closed card reads closed", func(t *testing.T) {})
	t.Run("TSSTT-B04: The decisions waiting on a person are listed", func(t *testing.T) {})
	t.Run("TSSTT-B06: Only the lock's own reversal counts, reduced to a clean first line", func(t *testing.T) {})
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
		e.PR.Checks[checkPassed] != 2 || e.PR.Checks[checkFailed] != 1 {
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
	t.Run("TSSTT-E01: Failed and unreadable lookups leave the part absent", func(t *testing.T) {})
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
	if e.Card != nil || e.PR != nil || e.Blocked != nil || e.Reverted != nil || e.Branch != "" {
		t.Errorf("with no source answering, nothing is invented: %+v", e)
	}
}

// Without --card, the card is the one the board says this agent owns.
func TestTaskStatusCmd_findsTheAgentsCardOnTheBoard(t *testing.T) {
	t.Run("TSSTT-B03: Without a number the agent's own card is found on the board", func(t *testing.T) {})
	t.Run("TSSTT-B07: The command prints the report", func(t *testing.T) {})
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
	t.Run("TSSTT-B08: The turn-ended event carries numbers and vocabulary only", func(t *testing.T) {})
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
		PR:   &branchPR{Number: 368, State: "OPEN", Total: 4, Checks: map[string]int{checkFailed: 3}},
	})
	common.Emitter.Flush()

	body := <-got
	// the attribute names are stable English identifiers, whatever the user's language
	for _, want := range []string{"in-review", "open", `"key":"checks_failed"`, `"key":"checks_running"`,
		`"key":"card_state"`, `"key":"pr_state"`, `"key":"has_card"`, `"key":"clean_tree"`} {
		if !strings.Contains(body, want) {
			t.Errorf("the event lacks %q: %s", want, body)
		}
	}
	for _, old := range []string{"reprovaram", "rodando", "estado", "tem_card", "arvore_limpa", "nao_enviado"} {
		if strings.Contains(body, old) {
			t.Errorf("the attribute names are English identifiers, found %q: %s", old, body)
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

// A check still running has no conclusion yet. It is running, not a failure: counting
// it as failed makes task-status tell the agent to fix a PR whose CI simply has not
// finished.
func TestCurrentBranchPR_runningChecksAreInProgressNotFailed(t *testing.T) {
	t.Run("TSSTT-B05: Running checks are running, not failed", func(t *testing.T) {})
	t.Run("TSSTT-I01: The check classes add up to the total", func(t *testing.T) {})
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
	if p.Checks[checkRunning] != 3 || p.Checks[checkFailed] != 1 || p.Checks[checkPassed] != 2 || p.Total != 6 ||
		p.Checks[checkRunning]+p.Checks[checkFailed]+p.Checks[checkPassed] != p.Total {
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
	t.Run("TSSTT-X01: Every lookup names the configured repository and the root's branch", func(t *testing.T) {})
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

// Outside github mode a card is not an issue, and no card is looked up.
func TestCollectTaskState_localModeLooksUpNoCard(t *testing.T) {
	t.Run("TSSTT-X02: Local mode looks up no card", func(t *testing.T) {})
	root := localProject(t)
	calls := scriptedGH(t, ghRule{match: "issue view 303 *",
		out: `{"number":303,"title":"x","state":"OPEN","labels":[{"name":"anchors:in-progress"}]}`})
	cfg, err := config.Load(root + "/anchors.yaml")
	if err != nil {
		t.Fatal(err)
	}
	e := collectTaskState(root, cfg, 303)
	if e.Card != nil || e.Blocked != nil {
		t.Errorf("local mode has no card to report: %+v", e)
	}
	if got := callsWith(calls(), "issue "); len(got) != 0 {
		t.Errorf("no card lookup may be made in local mode: %v", got)
	}
}

// THE REVERSION PARSER only accepts what the LOCK wrote.
//
// Two conditions, and both matter: the marker AND the author being the bot. A person's
// comment that happens to start with the same symbol is not a reversion — treating it as
// one would make the report accuse something that did not happen.
func TestRevertedOn_countsOnlyWhatTheLockWrote(t *testing.T) {
	t.Run("TSSTT-B06: Only the lock's own reversal counts, reduced to a clean first line", func(t *testing.T) {})
	cases := []struct {
		name   string
		body   string
		author string
		counts bool
	}{
		{"the real reversion", initx.MarcadorDeReversao + " **Reverted: closed by hand**", "github-actions", true},
		{"a person using the same symbol", initx.MarcadorDeReversao + " I locked this here", "someone", false},
		{"the bot saying something else", "▶️ Back to the queue.", "github-actions", false},
		{"an ordinary comment", "working on it", "someone", false},
	}
	for _, c := range cases {
		// Calls WHAT THE CODE USES. The first version rewrote the condition here, and
		// survived the mutation that removed the author check — the test measured itself.
		if ehReversao(c.body, c.author) != c.counts {
			t.Errorf("%s: expected counts=%v", c.name, c.counts)
		}
	}
}

// THE FIRST LINE is enough, without the bold markup. The whole comment has six paragraphs
// explaining the rule — dumping them in the terminal would turn the report into a wall.
func TestRevertedOn_showsOnlyTheCleanFirstLine(t *testing.T) {
	body := initx.MarcadorDeReversao + " **Reverted: closed by hand** (by `someone`), and the card was reopened.\n" +
		"\nThe card moves by FACT, not by hand: `in-progress` because the claim delivered...\n" +
		"\n**What this prevents:** a card closed before the merge leaves `ready-to-review`..."
	scriptedGH(t, ghRule{match: "issue view 9 --json comments*",
		out: `{"comments":[{"body":` + strconv.Quote(body) + `,"author":{"login":"github-actions"}}]}`})

	got := revertedOn("acme/app", 9)
	if len(got) != 1 || got[0] != "Reverted: closed by hand (by someone), and the card was reopened." {
		t.Errorf("only the first line, without markup nor marker; got %q", got)
	}
}

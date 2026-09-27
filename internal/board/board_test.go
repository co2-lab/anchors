package board

import (
	"strconv"
	"strings"
	"testing"
)

// boardJSON is one card as the board query's `--jq` emits it: one JSON object per line.
func boardJSON(number int, title, body, owner string, labels ...string) string {
	var ls []string
	for _, l := range labels {
		ls = append(ls, `{"name":"`+l+`"}`)
	}
	comments := ""
	if owner != "" {
		comments = `{"body":"anchors-owner: ` + owner + `"}`
	}
	return `{"number":` + strconv.Itoa(number) + `,"title":"` + title + `","body":"` + body +
		`","labels":[` + strings.Join(ls, ",") + `],"comments":[` + comments + `]}` + "\n"
}

// recorder is a fake `gh` that records each call and answers with the given output.
type recorder struct {
	calls  [][]string
	answer string
}

func (r *recorder) run(args ...string) ([]byte, error) {
	r.calls = append(r.calls, args)
	return []byte(r.answer), nil
}

func TestList_aCardNeedsEveryLabelAndTheState(t *testing.T) {
	t.Run("BRCRB-B01: A card needs every configured label and the asked state", func(t *testing.T) {})
	rec := &recorder{answer: boardJSON(1, "no team", "", "", "anchors", StateInProgress) +
		boardJSON(2, "to-do", "", "", "anchors", "team-x", StateToDo) +
		boardJSON(3, "in progress", "", "", "anchors", "team-x", StateInProgress)}
	c := Client{Repo: "o/r", Labels: []string{"anchors", "team-x"}, run: rec.run}
	cards, err := c.list(StateInProgress)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 || cards[0].Number != 3 {
		t.Fatalf("list(in-progress) = %+v, want only card 3", cards)
	}
}

// The OWNER is the LAST `anchors-owner:` comment: ownership changes hands — a rejected card
// goes back to the original author — and every claim stays recorded.
func TestLastOwner_theLastCommentWins(t *testing.T) {
	t.Run("BRCRB-B02: The owner is the last ownership comment", func(t *testing.T) {})
	r := rawCard{Comments: []struct{ Body string }{
		{Body: "anchors-owner: maq-a/sessao-1"},
		{Body: "some comment in the middle"},
		{Body: "anchors-owner: maq-b/sessao-2"},
	}}
	if got := lastOwner(r); got != "maq-b/sessao-2" {
		t.Errorf("owner read: %q — want the last one", got)
	}
	// A card with no ownership comment has no owner, and is free.
	if got := lastOwner(rawCard{}); got != "" {
		t.Errorf("a card with no comment gave the owner %q", got)
	}
}

// RULE 3: `needs-user` only goes to whoever DECLARED they decide the product. The default is
// closed — the zero-value `Client` does not act — and the cost of erring open is someone
// deciding the product without authority, invisible after the fact.
func TestDeclines_escalatedDependsOnTheDeclaration(t *testing.T) {
	t.Run("BRCRB-B03: An escalated card goes only to whoever decides the product", func(t *testing.T) {})
	for _, label := range []string{StateNeedsUser, StateNeedsUserPt} {
		card := Card{Labels: []string{StateToDo, label}}
		if !(Client{}).declines(card) {
			t.Errorf("%q: a client that declared nothing must NOT take an escalated card", label)
		}
		if !(Client{UserIssues: false}).declines(card) {
			t.Errorf("%q: whoever declared they do not decide keeps declining", label)
		}
		if (Client{UserIssues: true}).declines(card) {
			t.Errorf("%q: whoever declared they decide the product may take it", label)
		}
	}
	// The ORDINARY card is declined by nobody — the declaration only governs the escalated one.
	card := Card{Labels: []string{StateToDo}}
	for _, c := range []Client{{}, {UserIssues: true}} {
		if c.declines(card) {
			t.Errorf("an ordinary card declined by Client{UserIssues:%v}", c.UserIssues)
		}
	}
}

// The agent FINISHES WHAT IS WITH IT before taking something new: in-progress is half-done
// implementation, in-review is a review to do. What separates "mine" from "someone else's"
// is the owner, checked first.
func TestMine_resumesTheAgentsWorkUnderWay(t *testing.T) {
	t.Run("BRCRB-B04: The agent resumes its own work under way", func(t *testing.T) {})
	const me = "maq/sessao-1"
	rec := &recorder{answer: boardJSON(4, "escalated", "", me, "anchors", StateInProgress, StateNeedsUser) +
		boardJSON(5, "someone else's", "", "other/sessao", "anchors", StateInProgress) +
		boardJSON(6, "my to-do", "", me, "anchors", StateToDo) +
		boardJSON(7, "my review", "", me, "anchors", StateToDo, StateInReview)}
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: rec.run}
	card, err := c.Mine(me)
	if err != nil {
		t.Fatal(err)
	}
	if card == nil || card.Number != 7 || card.State != StateInReview {
		t.Fatalf("Mine = %+v, want card 7 in the in-review state", card)
	}
	// With only stopped cards of its own, the agent has no work under way.
	rec.answer = boardJSON(6, "my to-do", "", me, "anchors", StateToDo) +
		boardJSON(8, "my ready", "", me, "anchors", StateReadyToReview)
	if card, err := c.Mine(me); err != nil || card != nil {
		t.Fatalf("Mine with only stopped cards = %+v, %v; want none", card, err)
	}
}

func TestFindByCode_matchesTheTitleAndNotTheBody(t *testing.T) {
	t.Run("BRCRB-B05: A unit's card is found by the code in its title", func(t *testing.T) {})
	rec := &recorder{answer: `[{"number":1,"title":"[ALTPL] depends on GLCGL","body":"depends on [GLCGL]","labels":[]},
		{"number":2,"title":"[GLCGLX] other","body":"","labels":[]},
		{"number":3,"title":"[GLCGL] Implementar spec","body":"","labels":[]}]`}
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: rec.run}
	card, err := c.FindByCode("glcgl")
	if err != nil || card.Number != 3 {
		t.Fatalf("FindByCode(glcgl) = %+v, %v; want card 3", card, err)
	}
}

// An empty code is refused: without it the search would match "[]" and return any card.
func TestFindByCode_refusesAnEmptyCode(t *testing.T) {
	t.Run("BRCRB-E03: An empty code is refused", func(t *testing.T) {})
	c := Client{Repo: "x/y", Labels: []string{"anchors"}}
	for _, empty := range []string{"", "   ", "\t"} {
		if _, err := c.FindByCode(empty); err == nil {
			t.Errorf("the code %q passed — it would match any card", empty)
		}
	}
}

func TestFindByCode_refusesACodeNoTitleHolds(t *testing.T) {
	t.Run("BRCRB-E04: A code no title holds is refused", func(t *testing.T) {})
	rec := &recorder{answer: `[{"number":12,"title":"[ALTPL] other","body":"","labels":[]}]`}
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: rec.run}
	if _, err := c.FindByCode("RIMRD"); err == nil || !strings.Contains(err.Error(), "[RIMRD]") {
		t.Fatalf("FindByCode of an absent code = %v, want a refusal naming [RIMRD]", err)
	}
}

// Two open cards with the same code (reference app): the code alone cannot say which one the
// work belongs to, and picking the first recorded the delivery on the wrong card. The lookup
// refuses and names both, so the agent chooses with `--card <n>`.
func TestFindByCode_refusesAnAmbiguousCode(t *testing.T) {
	t.Run("BRCRB-E05: An ambiguous code is refused naming every card", func(t *testing.T) {})
	answer := ""
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: func(args ...string) ([]byte, error) {
		return []byte(answer), nil
	}}

	answer = `[{"number":966,"title":"[RIMRD] review of PR #884","body":"","labels":[]},
	           {"number":887,"title":"[RIMRD] B10 escopo efetivo","body":"","labels":[]},
	           {"number":12,"title":"[ALTPL] other","body":"","labels":[]}]`
	_, err := c.FindByCode("RIMRD")
	if err == nil {
		t.Fatal("two cards with [RIMRD] and the lookup picked one — the delivery could land on the wrong card")
	}
	for _, want := range []string{"#966", "#887", "--card"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal should name %s: %v", want, err)
		}
	}

	answer = `[{"number":887,"title":"[RIMRD] B10","body":"","labels":[]},{"number":12,"title":"[ALTPL] x","body":"","labels":[]}]`
	if card, err := c.FindByCode("RIMRD"); err != nil || card.Number != 887 {
		t.Errorf("a single match must still be found: %v %+v", err, card)
	}
}

// `deliver --card <n>` records on a card named directly: plan cards carry no `[CODE]` in the
// title. The card must still be OPEN and an Anchors card — otherwise the record goes where
// nobody reads it, or onto someone else's issue.
func TestFindOpenByNumber(t *testing.T) {
	t.Run("BRCRB-B06: An open Anchors card is returned by number", func(t *testing.T) {})
	t.Run("BRCRB-E06: A closed card or a non-Anchors issue is refused by number", func(t *testing.T) {})
	answer := ""
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: func(args ...string) ([]byte, error) {
		if args[0] != "issue" || args[1] != "view" || args[2] != "801" {
			t.Fatalf("unexpected gh call: %v", args)
		}
		return []byte(answer), nil
	}}

	answer = `{"number":801,"title":"[plano] docs","body":"","state":"OPEN","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]}`
	card, err := c.FindOpenByNumber(801)
	if err != nil || card.Number != 801 {
		t.Fatalf("an open Anchors card was refused: %v %+v", err, card)
	}

	answer = `{"number":801,"title":"x","body":"","state":"CLOSED","labels":[{"name":"anchors"}]}`
	if _, err := c.FindOpenByNumber(801); err == nil || !strings.Contains(err.Error(), "closed") {
		t.Errorf("a closed card was accepted: %v", err)
	}

	answer = `{"number":801,"title":"x","body":"","state":"OPEN","labels":[{"name":"bug"}]}`
	if _, err := c.FindOpenByNumber(801); err == nil || !strings.Contains(err.Error(), "not an Anchors card") {
		t.Errorf("an issue without the Anchors label was accepted: %v", err)
	}
}

func TestGh_apiCallsCarryNoRepoFlag(t *testing.T) {
	t.Run("BRCRB-B07: Board queries do not name the repository, other calls do", func(t *testing.T) {})
	rec := &recorder{}
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: rec.run}
	if _, err := c.Mine("me"); err != nil {
		t.Fatal(err)
	}
	rec.answer = `{"number":801,"title":"x","body":"","state":"OPEN","labels":[{"name":"anchors"}]}`
	if _, err := c.FindOpenByNumber(801); err != nil {
		t.Fatal(err)
	}
	if len(rec.calls) != 2 {
		t.Fatalf("want two gh calls, got %v", rec.calls)
	}
	api, view := strings.Join(rec.calls[0], " "), strings.Join(rec.calls[1], " ")
	if !strings.HasPrefix(api, "api graphql") || strings.Contains(api, "--repo") {
		t.Errorf("the board query must carry no --repo: %s", api)
	}
	if !strings.HasSuffix(view, "--repo o/r") {
		t.Errorf("the issue view must name the repository: %s", view)
	}
}

func TestAsk_dispatchesTheClaimPipeline(t *testing.T) {
	t.Run("BRCRB-B08: Asking for work dispatches the claim pipeline", func(t *testing.T) {})
	rec := &recorder{}
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: rec.run}
	if err := c.Ask("maq/1"); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(rec.calls[0], " "); got != "workflow run "+ClaimWorkflow+" -f agent=maq/1 --repo o/r" {
		t.Errorf("Ask called gh %q", got)
	}
}

// `workflow.labels` empty is an ERROR, not a default: without it the claim would pull any issue
// of the repository — product issues included.
func TestList_emptyLabelsIsAnError(t *testing.T) {
	t.Run("BRCRB-E01: Empty labels are refused", func(t *testing.T) {})
	c := Client{Repo: "org/repo"}
	if _, err := c.list(StateToDo); err == nil {
		t.Error("empty labels passed — the claim would pull product issues")
	}
}

func TestList_repoMustBeOwnerSlashName(t *testing.T) {
	t.Run("BRCRB-E02: A repository that is not owner/name is refused", func(t *testing.T) {})
	rec := &recorder{}
	c := Client{Repo: "just-a-name", Labels: []string{"anchors"}, run: rec.run}
	if _, err := c.list(""); err == nil || !strings.Contains(err.Error(), "just-a-name") {
		t.Fatalf("list with repo just-a-name = %v, want the refusal naming it", err)
	}
	if len(rec.calls) != 0 {
		t.Errorf("no gh call may happen: %v", rec.calls)
	}
}

func TestComment_refusesAnEmptyBody(t *testing.T) {
	t.Run("BRCRB-E07: An empty comment is refused", func(t *testing.T) {})
	c := Client{Repo: "x/y"}
	if err := c.Comment(1, "   "); err == nil {
		t.Error("an empty body passed — a blank comment records nothing")
	}
}

// @anchors
//   code: DCTSD
//   ref: DCDDE

package flow

import (
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// `escalate --for-user` wrote a state it did not know how to revert.
//
// Measured (co2-lab/anchors#10): the decision came out, the revisions were applied to the
// 4 files, the decision issue was closed — and the card kept `anchors:needs-user`, the label
// that makes the claim skip it. It was removed by hand.
func TestDecided_requiresCardAndResolution(t *testing.T) {
	t.Run("DCDDE-B01: decided without a card is refused", func(t *testing.T) {})
	cases := []struct {
		name  string
		args  []string
		names string // what the error message must say
	}{
		{"no card", []string{"--resolution", "ABCDE-R0001: x"}, "--card"},
		{"no resolution", []string{"--card", "4"}, "--resolution"},
		{"blank resolution", []string{"--card", "4", "--resolution", "   "}, "--resolution"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmd := newDecidedCmd()
			cmd.SetArgs(c.args)
			cmd.SetOut(os.Stderr)
			cmd.SetErr(os.Stderr)
			err := cmd.Execute()
			if err == nil {
				t.Fatal("it accepted without the required argument")
			}
			if !strings.Contains(err.Error(), c.names) {
				t.Errorf("the error does not say what is missing (%q):\n%s", c.names, err)
			}
		})
	}
}

// The mandatory `--resolution` is not bureaucracy: the exit of an open decision is ONE —
// the answer becomes a RULE, with a code. Releasing the card without saying which revision
// was born from it would leave the decision without a trace.
func TestDecided_theMessageSaysWhyTheResolutionIsRequired(t *testing.T) {
	t.Run("DCDDE-B02: decided without a resolution is refused and says why", func(t *testing.T) {})
	for _, args := range [][]string{{"--card", "4"}, {"--card", "4", "--resolution", "   "}} {
		cmd := newDecidedCmd()
		cmd.SetArgs(args)
		cmd.SetOut(os.Stderr)
		cmd.SetErr(os.Stderr)
		err := cmd.Execute()
		if err == nil {
			t.Fatalf("%v: accepted without --resolution", args)
		}
		// Saying it is missing is not enough: it must say WHY, or the next person passes
		// any string to satisfy the command.
		if !strings.Contains(err.Error(), "--resolution") || !strings.Contains(err.Error(), "becomes a RULE") {
			t.Errorf("%v: the error does not explain the reason for the requirement:\n%s", args, err)
		}
	}
}

// In LOCAL mode the decision lives in `issues/`, and resolving it is moving the file —
// there is no label to remove. Saying so beats failing inside `gh`.
func TestDecided_refusesOutsideGithubMode(t *testing.T) {
	t.Run("DCDDE-B03: decided in local mode is refused", func(t *testing.T) {})
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--root", localProject(t), "--card", "4", "--resolution", "ABCDE-R0001: x"})
	cmd.SetOut(os.Stderr)
	cmd.SetErr(os.Stderr)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("it accepted outside github mode")
	}
	if !strings.Contains(err.Error(), "issues/") {
		t.Errorf("the error does not point to where the decision lives in local mode:\n%s", err)
	}
}

// The guard reverts a close by a person that lacks `anchors:manual`, and `decided` runs
// under a person's account: without the label first, the closed decision came back.
func TestDecided_closeLabelsManualBeforeClosing(t *testing.T) {
	calls := closeDecisionArgs("o/r", 832, "X-R0001: y")
	if len(calls) != 2 {
		t.Fatalf("expected label then close, got %v", calls)
	}
	label, closing := strings.Join(calls[0], " "), strings.Join(calls[1], " ")
	if !strings.Contains(label, "issue edit 832") || !strings.Contains(label, "--add-label "+initx.LabelManual) {
		t.Errorf("the first call must put %s on the issue: %q", initx.LabelManual, label)
	}
	if !strings.Contains(closing, "issue close 832") || !strings.Contains(closing, "X-R0001: y") {
		t.Errorf("the second call must close with the resolution: %q", closing)
	}
}

// indexOf returns the position of the first call holding every fragment, or -1.
func indexOf(calls []string, fragments ...string) int {
	for i, c := range calls {
		if len(callsWith([]string{c}, fragments...)) == 1 {
			return i
		}
	}
	return -1
}

// The whole release: the needs-user label and every blocked-by label leave together, the
// trace of who blocked it stays in the comment, and each decision issue under the card
// is labelled manual and closed with the resolution.
func TestDecidedCmd_releasesTheCardAndClosesItsDecisions(t *testing.T) {
	t.Run("DCDDE-B05: The release removes needs-user and every blocked-by label at once", func(t *testing.T) {})
	t.Run("DCDDE-B06: The comment keeps the resolution and who blocked the card", func(t *testing.T) {})
	t.Run("DCDDE-B07: The decisions under the card are labelled manual and closed", func(t *testing.T) {})
	t.Run("DCDDE-I01: Blockers are read first and the card is released before its decisions close", func(t *testing.T) {})
	t.Run("DCDDE-X01: Only the decisions under this card are closed", func(t *testing.T) {})
	root := githubProject(t)
	blocked := initx.LabelBlockedBy("701")
	calls := scriptedGH(t,
		ghRule{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "[]"},
		ghRule{match: "issue view 4 *--json labels*", out: blocked + "\n" + initx.LabelBlockedBy("702")},
		ghRule{match: "issue list *" + initx.LabelSob("4") + "*", out: `[{"number":701},{"number":702}]`},
		ghRule{match: "issue close 702 *", code: 1},
	)
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--root", root, "--card", "4", "--resolution", "PLTFR-R0004: mTLS is built in F03"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}

	c := calls()
	edit := callsWith(c, "issue edit 4 ")
	if len(edit) != 1 {
		t.Fatalf("expected one release of card 4; calls: %v", c)
	}
	want := "--remove-label " + initx.LabelNeedsUser + "," + blocked + "," + initx.LabelBlockedBy("702")
	if !strings.Contains(edit[0], want) {
		t.Errorf("the release must remove needs-user AND every blocked-by: %s", edit[0])
	}
	comment := callsWith(c, "issue comment 4 ")
	if len(comment) != 1 || !strings.Contains(comment[0], "PLTFR-R0004") ||
		!strings.Contains(comment[0], "**Was blocked by:** #701, #702") {
		t.Errorf("the comment must keep the resolution and the trace of the blockers: %v", comment)
	}
	// Decisions are those with needs-user AND this card's under label — nothing else.
	if len(callsWith(c, "issue list", "--state open", "--label "+initx.LabelNeedsUser, "--label "+initx.LabelSob("4"))) != 1 {
		t.Errorf("the decisions must be filtered by needs-user and the card's under label; calls: %v", c)
	}
	// 701 closes; 702 fails its close, so only one counts.
	if len(callsWith(c, "issue edit 701", initx.LabelManual)) != 1 || len(callsWith(c, "issue close 701", "PLTFR-R0004")) != 1 {
		t.Errorf("decision 701 must be labelled manual and closed with the resolution; calls: %v", c)
	}
	if !strings.Contains(out, "card #4 released") || !strings.Contains(out, "1 decision issue closed") {
		t.Errorf("the output must say the card was released and how many decisions closed:\n%s", out)
	}
	// THE ORDER: the blockers are read while the labels still exist; the card is free
	// before anything else happens, so a later failure never leaves it stopped.
	view, release, note := indexOf(c, "issue view 4"), indexOf(c, "issue edit 4 "), indexOf(c, "issue comment 4 ")
	manual, closing := indexOf(c, "issue edit 701"), indexOf(c, "issue close 701")
	if view < 0 || !(view < release && release < note && note < manual && manual < closing) {
		t.Errorf("order must be view(%d) < release(%d) < comment(%d) < manual(%d) < close(%d); calls: %v",
			view, release, note, manual, closing, c)
	}
}

// Moved from backfill: a reading of decided.go's source. The behaviour it guards (read the
// blockers before removing them, and record them) is proven by the order and the comment
// checked in TestDecidedCmd_releasesTheCardAndClosesItsDecisions.
//
// UNBLOCKING DOES NOT ERASE HISTORY: the label is the only place the link lives, and the
// comment is immutable and dated.
func TestUnblockRecordsTheLinkBeforeRemoving(t *testing.T) {
	source := leFonte(t, "decided.go")

	// The READ comes before the removal: afterwards the information exists nowhere.
	iRead := strings.Index(source, "labelsDeBloqueio(cfg.Workflow.Repo, card)")
	iRemove := strings.Index(source, `"--remove-label"`)
	if iRead < 0 {
		t.Fatal("`decided` does not read the blocking labels — they would hang on a card " +
			"already free, and the board would show a block that no longer exists")
	}
	if iRemove < 0 {
		t.Fatal("the label removal in `decided` was not found — the command changed shape")
	}
	if iRead > iRemove {
		t.Error("`decided` reads the blocking labels AFTER removing them: at that point the " +
			"information no longer exists, and the trace comes out empty")
	}

	// The TRACE in a comment, with the numbers: saying "it was blocked" without saying by
	// whom records that there was a wait and loses the way back to the decision.
	if !strings.Contains(source, "**Was blocked by:**") {
		t.Error("`decided` removes the block without recording who the card waited for")
	}
}

// A decision that generated work holds the card until that work is delivered.
func TestDecidedCmd_refusesWhileAnUnblockCardIsOpen(t *testing.T) {
	t.Run("DCDDE-B04: decided refuses while an unblock card is open", func(t *testing.T) {})
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: `[{"number":512},{"number":513}]`},
	)
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--root", root, "--card", "4", "--resolution", "X-R0001: y"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("the card must not be released while its unblock work is open")
	}
	if !strings.Contains(err.Error(), "#512, #513") || !strings.Contains(err.Error(), initx.LabelNeedsUser) {
		t.Errorf("the refusal must name the pending cards and the label that stays: %v", err)
	}
	if got := callsWith(calls(), "issue edit"); len(got) != 0 {
		t.Errorf("nothing may be edited on a refusal: %v", got)
	}
}

func TestDecidedCmd_failedReleaseIsAnError(t *testing.T) {
	t.Run("DCDDE-E02: A failed release fails the command", func(t *testing.T) {})
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "[]"},
		ghRule{match: "issue list *" + initx.LabelSob("4") + "*", out: `[{"number":701}]`},
		ghRule{match: "issue edit 4 *", out: "no permission", code: 1})
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--root", root, "--card", "4", "--resolution", "X-R0001: y"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err == nil || !strings.Contains(err.Error(), "release card #4") || !strings.Contains(err.Error(), "no permission") {
		t.Fatalf("a failed release must fail the command with gh's answer, got %v", err)
	}
	if strings.Contains(out, "released") {
		t.Errorf("a failed release must not be announced:\n%s", out)
	}
	c := calls()
	if len(callsWith(c, "issue comment")) != 0 || len(callsWith(c, "issue close")) != 0 {
		t.Errorf("a card that was not released is neither commented nor are its decisions closed: %v", c)
	}
}

// Zero and many decisions are worded apart from one.
func TestDecidedCmd_countsTheClosedDecisions(t *testing.T) {
	t.Run("DCDDE-B08: The output counts the decisions closed", func(t *testing.T) {})
	for _, c := range []struct {
		list, want string
	}{
		{"[]", "no open decision issue under card #4"},
		{`[{"number":1},{"number":2}]`, "2 decision issues closed"},
	} {
		root := githubProject(t)
		scriptedGH(t,
			ghRule{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "[]"},
			ghRule{match: "issue list *" + initx.LabelSob("4") + "*", out: c.list})
		cmd := newDecidedCmd()
		cmd.SetArgs([]string{"--root", root, "--card", "4", "--resolution", "X-R0001: y"})
		var err error
		out := stdoutOf(t, func() { err = cmd.Execute() })
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("with %s decisions the output must say %q:\n%s", c.list, c.want, out)
		}
	}
}

func TestLabelsDeBloqueio_readsEveryBlockerAndToleratesFailure(t *testing.T) {
	t.Run("DCDDE-E04: Unreadable labels report no blocker", func(t *testing.T) {})
	scriptedGH(t, ghRule{match: "issue view 9 *", out: "anchors:blocked-by-1\n\n  anchors:blocked-by-2 "})
	got := labelsDeBloqueio("acme/app", "9")
	if strings.Join(got, ",") != "anchors:blocked-by-1,anchors:blocked-by-2" {
		t.Errorf("every blocked-by label must be read, got %v", got)
	}

	scriptedGH(t, ghRule{match: "issue view 9 *", code: 1})
	if got := labelsDeBloqueio("acme/app", "9"); got != nil {
		t.Errorf("a failed read yields nothing, got %v", got)
	}

	// And the release still happens: only needs-user leaves.
	calls := scriptedGH(t,
		ghRule{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "[]"},
		ghRule{match: "issue view 4 *", code: 1})
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "--card", "4", "--resolution", "X-R0001: y"})
	var err error
	stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if got := callsWith(calls(), "issue edit 4 ", "--remove-label "+initx.LabelNeedsUser); len(got) != 1 ||
		strings.Contains(got[0], "blocked-by") {
		t.Errorf("with unreadable labels the card is still released of needs-user alone: %v", got)
	}
}

func TestCloseDecisionsUnder_unreadableListClosesNothing(t *testing.T) {
	t.Run("DCDDE-E03: An unreadable decision list closes nothing", func(t *testing.T) {})
	calls := scriptedGH(t, ghRule{match: "issue list *", out: "not json"})
	if n := closeDecisionsUnder("acme/app", "4", "X-R0001: y"); n != 0 {
		t.Errorf("an unreadable list closes nothing, got %d", n)
	}
	if got := callsWith(calls(), "issue close"); len(got) != 0 {
		t.Errorf("nothing may be closed: %v", got)
	}
}

// Not knowing whether unblock work is still open must not become "can release": a failed
// or unreadable lookup refuses, says why, and edits nothing.
func TestDecidedCmd_refusesWhenTheUnblockLookupFails(t *testing.T) {
	t.Run("DCDDE-E01: An unblock lookup that fails or is unreadable refuses", func(t *testing.T) {})
	for _, r := range []ghRule{
		{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "HTTP 502", code: 1},
		{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "not json"},
	} {
		root := githubProject(t)
		calls := scriptedGH(t, r)
		cmd := newDecidedCmd()
		cmd.SetArgs([]string{"--root", root, "--card", "4", "--resolution", "X-R0001: y"})
		err := cmd.Execute()
		if err == nil {
			t.Fatalf("with %q the card must not be released", r.out)
		}
		if !strings.Contains(err.Error(), initx.LabelDesbloqueia("4")) || !strings.Contains(err.Error(), initx.LabelNeedsUser) {
			t.Errorf("the refusal must say what could not be read and which label stays: %v", err)
		}
		if got := callsWith(calls(), "issue edit"); len(got) != 0 {
			t.Errorf("nothing may be edited when the lookup failed: %v", got)
		}
	}
}

// @anchors
//   ref: DSCRD

package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// Without a card there is nothing to discard: the command refuses before any call.
func TestDiscardCmd_requiresACard(t *testing.T) {
	t.Run("DSCRD-B01: Discard without any card argument is refused", func(t *testing.T) {})
	calls := scriptedGH(t)
	cmd := newDiscardCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "--reason", "obsolete"})
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	if err := cmd.Execute(); err == nil {
		t.Fatal("discard with no card must fail")
	}
	if c := calls(); len(c) != 0 {
		t.Errorf("no call may reach the platform without a card: %v", c)
	}
}

// The ORDER matters: label, reason, close.
//
// Recording the reason BEFORE applying the label is the worst possible inconsistency — if
// the label fails, the card carries a comment saying it was discarded AND stays on the
// board, stating two contrary things at once.
//
// Closing comes last, and only after everything worked: a discarded card that stays open
// would be served by `claim`.
func TestDiscard_appliesTheLabelBeforeRecordingTheReason(t *testing.T) {
	t.Run("DSCRD-B05: A card is labelled, then commented, then closed", func(t *testing.T) {})
	t.Run("DSCRD-X01: Discarding never deletes the issue", func(t *testing.T) {})
	calls := scriptedGH(t)
	if err := descarta("acme/project", "42", "the file no longer exists"); err != nil {
		t.Fatal(err)
	}
	c := calls()
	if len(c) != 3 {
		t.Fatalf("expected label, comment and close; got: %v", c)
	}
	if !strings.Contains(c[0], "issue edit 42") || !strings.Contains(c[0], "--add-label "+initx.LabelDiscarded) {
		t.Errorf("the FIRST call must apply the label: %q", c[0])
	}
	if !strings.Contains(c[1], "issue comment 42") {
		t.Errorf("the reason comes AFTER the label: %q", c[1])
	}
	if !strings.Contains(c[2], "issue close 42") {
		t.Errorf("closing comes last: %q", c[2])
	}
	if got := callsWith(c, "delete"); len(got) != 0 {
		t.Errorf("discarding must never delete the issue: %v", got)
	}
}

// The REASON goes on the card, not only on the screen of whoever ran the command.
func TestDiscard_recordsTheReasonOnTheCard(t *testing.T) {
	t.Run("DSCRD-B06: The comment carries the reason and the way back", func(t *testing.T) {})
	calls := scriptedGH(t)
	if err := descarta("acme/project", "42", "spec deleted in the cleanup of plan 0004"); err != nil {
		t.Fatal(err)
	}
	comments := callsWith(calls(), "issue comment")
	if len(comments) != 1 || !strings.Contains(comments[0], "spec deleted in the cleanup") {
		t.Fatalf("the reason did not reach the card: %v", comments)
	}
	// And the way BACK: a discard with no way to undo it is a removal under another name.
	if !strings.Contains(comments[0], "remove the label `"+initx.LabelDiscarded+"`") {
		t.Error("the comment does not say which label to remove to bring the card back")
	}
}

// A card that is already closed (a test card, a finding about a deleted file) is still
// discarded: the close is the last step and its failure does not undo the other two.
func TestDiscard_alreadyClosedCardIsStillDiscarded(t *testing.T) {
	t.Run("DSCRD-B08: A card that is already closed is still discarded", func(t *testing.T) {})
	scriptedGH(t, ghRule{match: "issue close 42 *", out: "issue already closed", code: 1})
	if err := descarta("acme/project", "42", "obsolete"); err != nil {
		t.Errorf("a failed close must not fail the discard, got %v", err)
	}
}

// The command discards every card it is given, and the ones that fail are named in the
// error instead of stopping the others.
func TestDiscardCmd_discardsEachCardAndNamesTheOnesThatFailed(t *testing.T) {
	t.Run("DSCRD-E01: A card whose label fails is named in the error while the others are discarded", func(t *testing.T) {})
	t.Run("DSCRD-B07: The leading hash is stripped and a blank argument is skipped", func(t *testing.T) {})
	t.Run("DSCRD-B09: Each discarded card is reported on standard output", func(t *testing.T) {})
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue edit 43 *", out: "HTTP 404: issue not found", code: 1},
	)
	cmd := newDiscardCmd()
	cmd.SetArgs([]string{"--root", root, "--reason", "the file no longer exists", "#42", "43", " "})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })

	if err == nil {
		t.Fatal("a card that was not discarded must make the command fail")
	}
	if !strings.Contains(err.Error(), "1 card(s) were not discarded") || !strings.Contains(err.Error(), "#43") ||
		!strings.Contains(err.Error(), "issue not found") {
		t.Errorf("the error must name the failed card and why: %v", err)
	}
	if !strings.Contains(out, "#42 discarded") {
		t.Errorf("the card that worked is reported as discarded:\n%s", out)
	}
	if strings.Contains(out, "#43 discarded") {
		t.Errorf("the card that failed was reported as discarded:\n%s", out)
	}
	c := calls()
	if got := callsWith(c, "issue edit 42 "); len(got) != 1 {
		t.Errorf("card #42 must be edited by its bare number; calls: %v", c)
	}
	// Two cards and the label: a blank argument makes no call of its own.
	if got := callsWith(c, "issue edit"); len(got) != 2 {
		t.Errorf("only the two real cards are edited; calls: %v", c)
	}
	if got := callsWith(c, "issue close 42"); len(got) != 1 {
		t.Errorf("card 42 must be closed once; calls: %v", c)
	}
	if got := callsWith(c, "issue comment 43"); len(got) != 0 {
		t.Errorf("card 43 failed its label and must not be commented; calls: %v", c)
	}
	if got := callsWith(c, "issue close 43"); len(got) != 0 {
		t.Errorf("card 43 failed its label and must not be closed; calls: %v", c)
	}
}

// The label is created on demand, before any card is edited: `gh issue edit` with a
// missing label fails the whole call. When it already exists, creating it fails, and that
// failure is not the command's.
func TestDiscardCmd_createsTheLabelFirstAndToleratesAnExistingOne(t *testing.T) {
	t.Run("DSCRD-B04: The discard label is created before any card is touched, and an existing label does not stop it", func(t *testing.T) {})
	calls := scriptedGH(t,
		ghRule{match: "label create *", out: "label already exists", code: 1},
	)
	cmd := newDiscardCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "--reason", "obsolete", "42"})
	var err error
	stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatalf("an existing label must not stop the discard, got %v", err)
	}
	c := calls()
	if len(c) == 0 || !strings.Contains(c[0], "label create "+initx.LabelDiscarded) {
		t.Errorf("the first call must create the discard label; calls: %v", c)
	}
	if got := callsWith(c, "issue close 42"); len(got) != 1 {
		t.Errorf("card 42 must still be discarded; calls: %v", c)
	}
}

func TestDiscardCmd_requiresAReasonAndGithubMode(t *testing.T) {
	t.Run("DSCRD-B02: Discard with a blank reason is refused", func(t *testing.T) {})
	t.Run("DSCRD-B03: Discard in local mode is refused", func(t *testing.T) {})
	scriptedGH(t)
	cmd := newDiscardCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "--reason", "  ", "42"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "--reason") {
		t.Errorf("a blank reason must be refused, got %v", err)
	}

	cmd = newDiscardCmd()
	cmd.SetArgs([]string{"--root", localProject(t), "--reason", "obsolete", "42"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode has no card to discard, got %v", err)
	}
}

// When the reason cannot be recorded the card is NOT closed: a closed card with no reason
// is indistinguishable from a lost one.
func TestDiscard_failedCommentStopsBeforeClosing(t *testing.T) {
	t.Run("DSCRD-E02: A card whose reason cannot be recorded is not closed", func(t *testing.T) {})
	t.Run("DSCRD-I01: A card is never closed before its label and reason are recorded", func(t *testing.T) {})
	calls := scriptedGH(t,
		ghRule{match: "issue comment 42 *", out: "rate limited", code: 1},
		ghRule{match: "issue edit 43 *", out: "HTTP 404", code: 1},
	)
	err := descarta("acme/app", "42", "obsolete")
	if err == nil || !strings.Contains(err.Error(), "record the reason") || !strings.Contains(err.Error(), "rate limited") {
		t.Fatalf("the failed comment must be reported, got %v", err)
	}
	if err := descarta("acme/app", "43", "obsolete"); err == nil {
		t.Fatal("the failed label must be reported")
	}
	if got := callsWith(calls(), "issue close"); len(got) != 0 {
		t.Errorf("a card was closed without its label or its reason: %v", got)
	}
}

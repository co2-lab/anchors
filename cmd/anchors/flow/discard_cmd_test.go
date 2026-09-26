package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// The command discards every card it is given, and the ones that fail are named in the
// error instead of stopping the others.
func TestDiscardCmd_discardsEachCardAndNamesTheOnesThatFailed(t *testing.T) {
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
	// The label is created on demand, before any card is edited: `gh issue edit` with a
	// missing label fails the whole call.
	if len(c) == 0 || !strings.Contains(c[0], "label create "+initx.LabelDiscarded) {
		t.Errorf("the first call must create the discard label; calls: %v", c)
	}
	if got := callsWith(c, "issue close 42"); len(got) != 1 {
		t.Errorf("card 42 must be closed once; calls: %v", c)
	}
	if got := callsWith(c, "issue close 43"); len(got) != 0 {
		t.Errorf("card 43 failed its label and must not be closed; calls: %v", c)
	}
}

func TestDiscardCmd_requiresAReasonAndGithubMode(t *testing.T) {
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
func TestDescarta_failedCommentStopsBeforeClosing(t *testing.T) {
	calls := scriptedGH(t,
		ghRule{match: "issue comment 42 *", out: "rate limited", code: 1},
	)
	err := descarta("acme/app", "42", "obsolete")
	if err == nil || !strings.Contains(err.Error(), "record the reason") || !strings.Contains(err.Error(), "rate limited") {
		t.Fatalf("the failed comment must be reported, got %v", err)
	}
	if got := callsWith(calls(), "issue close"); len(got) != 0 {
		t.Errorf("the card was closed without its reason: %v", got)
	}
}

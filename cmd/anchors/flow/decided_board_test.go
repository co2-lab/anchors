package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// The whole release: the needs-user label and every blocked-by label leave together, the
// trace of who blocked it stays in the comment, and each decision issue under the card
// is labelled manual and closed with the resolution.
func TestDecidedCmd_releasesTheCardAndClosesItsDecisions(t *testing.T) {
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
	// 701 closes; 702 fails its close, so only one counts.
	if len(callsWith(c, "issue edit 701", initx.LabelManual)) != 1 || len(callsWith(c, "issue close 701")) != 1 {
		t.Errorf("decision 701 must be labelled manual and closed; calls: %v", c)
	}
	if !strings.Contains(out, "card #4 released") || !strings.Contains(out, "1 decision issue closed") {
		t.Errorf("the output must say the card was released and how many decisions closed:\n%s", out)
	}
}

// A decision that generated work holds the card until that work is delivered.
func TestDecidedCmd_refusesWhileAnUnblockCardIsOpen(t *testing.T) {
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
	root := githubProject(t)
	scriptedGH(t,
		ghRule{match: "issue list *" + initx.LabelDesbloqueia("4") + "*", out: "[]"},
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
}

// Zero and many decisions are worded apart from one.
func TestDecidedCmd_countsTheClosedDecisions(t *testing.T) {
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
	scriptedGH(t, ghRule{match: "issue view 9 *", out: "anchors:blocked-by-1\n\n  anchors:blocked-by-2 "})
	got := labelsDeBloqueio("acme/app", "9")
	if strings.Join(got, ",") != "anchors:blocked-by-1,anchors:blocked-by-2" {
		t.Errorf("every blocked-by label must be read, got %v", got)
	}

	scriptedGH(t, ghRule{match: "issue view 9 *", code: 1})
	if got := labelsDeBloqueio("acme/app", "9"); got != nil {
		t.Errorf("a failed read yields nothing, got %v", got)
	}
}

func TestCloseDecisionsUnder_unreadableListClosesNothing(t *testing.T) {
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

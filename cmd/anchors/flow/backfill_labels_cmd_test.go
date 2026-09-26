package flow

import (
	"bytes"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// backfillBoard is a board of open decisions that exercises every path of the command:
//
//	#701 under-44   → origin #44 is an OPEN work card: the unequivocal link
//	#702 under-45   → origin #45 is closed: skipped
//	#703 under-46   → origin #46 is itself an open decision: precedence is not inferred
//	#704 under-704  → a decision does not hold itself
//	#705 (no under) → cites "PR #88", whose body declares `Refs #50`: provenance recovered
//	#706 (no under) → cites "#99" without the word PR: prose, not a link
func backfillBoard(t *testing.T) func() []string {
	t.Helper()
	under := func(n string) string { return `{"name":"` + initx.LabelSob(n) + `"}` }
	return scriptedGH(t,
		ghRule{match: "issue list *--label " + initx.LabelNeedsUser + "*", out: `[` +
			`{"number":701,"labels":[` + under("44") + `]},` +
			`{"number":702,"labels":[` + under("45") + `]},` +
			`{"number":703,"labels":[` + under("46") + `]},` +
			`{"number":704,"labels":[` + under("704") + `]},` +
			`{"number":705,"labels":[{"name":"anchors"}]},` +
			`{"number":706,"labels":[]}]`},
		ghRule{match: "issue view 705 *--json body*", out: "Seen while reviewing PR #88: the loop truncates."},
		ghRule{match: "issue view 706 *--json body*", out: "Related to #99 somehow."},
		ghRule{match: "pr view 88 *", out: "Some summary\n\nRefs #50\n"},
		ghRule{match: "issue view 44 *--json state,labels*", out: `{"state":"OPEN","labels":[{"name":"anchors:in-progress"}]}`},
		ghRule{match: "issue view 45 *--json state,labels*", out: `{"state":"CLOSED","labels":[]}`},
		ghRule{match: "issue view 46 *--json state,labels*", out: `{"state":"OPEN","labels":[{"name":"` + initx.LabelNeedsUser + `"}]}`},
	)
}

func TestBackfillLabelsCmd_writesOnlyTheUnequivocalLinks(t *testing.T) {
	root := githubProject(t)
	calls := backfillBoard(t)
	cmd := newBackfillLabelsCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs([]string{"--root", root})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	c := calls()
	edits := callsWith(c, "issue edit")
	// Exactly two cards are written: #44 gets blocked-by-701, #705 gets its provenance.
	if len(edits) != 2 {
		t.Fatalf("expected two edits, got %v", edits)
	}
	if len(callsWith(edits, "issue edit 44 ", "--add-label "+initx.LabelBlockedBy("701"))) != 1 {
		t.Errorf("#44 must receive blocked-by-701; edits: %v", edits)
	}
	if len(callsWith(edits, "issue edit 705 ", initx.LabelDePR("88")+","+initx.LabelSob("50"))) != 1 {
		t.Errorf("#705 must receive from-pr-88 and under-50; edits: %v", edits)
	}
	for _, n := range []string{"45", "46", "704", "706"} {
		if len(callsWith(edits, "issue edit "+n+" ")) != 0 {
			t.Errorf("#%s must not be touched; edits: %v", n, edits)
		}
	}
	s := out.String()
	for _, want := range []string{
		"#705 ← PR #88, card #50",
		"1 provenance(s) recovered from the cited PR",
		"#46 is also an open decision",
		"#44 ← blocked by #701",
		"1 link(s) written · 2 skipped",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("the output lacks %q:\n%s", want, s)
		}
	}
}

// The dry run names what it would do and touches nothing.
func TestBackfillLabelsCmd_dryRunTouchesNothing(t *testing.T) {
	root := githubProject(t)
	calls := backfillBoard(t)
	cmd := newBackfillLabelsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--root", root, "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	c := calls()
	if got := append(callsWith(c, "issue edit"), callsWith(c, "label create")...); len(got) != 0 {
		t.Errorf("a dry run must write nothing: %v", got)
	}
	s := out.String()
	for _, want := range []string{
		"#705 would receive",
		"#44 would receive `" + initx.LabelBlockedBy("701") + "`",
		"(nothing was touched)",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("the dry run lacks %q:\n%s", want, s)
		}
	}
}

func TestBackfillLabelsCmd_nothingToRecoverAndFailures(t *testing.T) {
	root := githubProject(t)
	scriptedGH(t, ghRule{match: "issue list *", out: `[{"number":1,"labels":[]}]`})
	cmd := newBackfillLabelsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--root", root})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "no link to recover: 1 open decision(s)") {
		t.Errorf("an empty result must say so:\n%s", out.String())
	}

	scriptedGH(t, ghRule{match: "issue list *", code: 1})
	cmd = newBackfillLabelsCmd()
	cmd.SetArgs([]string{"--root", root})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "list the open decisions") {
		t.Errorf("a failed list must fail the command, got %v", err)
	}

	scriptedGH(t, ghRule{match: "issue list *", out: "not json"})
	cmd = newBackfillLabelsCmd()
	cmd.SetArgs([]string{"--root", root})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "read the list") {
		t.Errorf("an unreadable list must fail the command, got %v", err)
	}

	cmd = newBackfillLabelsCmd()
	cmd.SetArgs([]string{"--root", localProject(t)})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode must be refused, got %v", err)
	}
}

// A failed edit is reported on stderr and counted as skipped, not as written.
func TestBackfillLabelsCmd_failedEditIsSkipped(t *testing.T) {
	root := githubProject(t)
	scriptedGH(t,
		ghRule{match: "issue list *", out: `[{"number":701,"labels":[{"name":"` + initx.LabelSob("44") + `"}]}]`},
		ghRule{match: "issue view 44 *", out: `{"state":"OPEN","labels":[]}`},
		ghRule{match: "issue edit 44 *", out: "forbidden", code: 1},
	)
	cmd := newBackfillLabelsCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs([]string{"--root", root})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(errOut.String(), "#44") || !strings.Contains(errOut.String(), "forbidden") {
		t.Errorf("the failure must be reported with gh's answer: %q", errOut.String())
	}
	if !strings.Contains(out.String(), "0 link(s) written · 1 skipped") {
		t.Errorf("a failed edit is skipped, not written:\n%s", out.String())
	}
}

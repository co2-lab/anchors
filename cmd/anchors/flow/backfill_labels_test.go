// @anchors
//   ref: BCLBB

package flow

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// leFonte reads a file of this package so a test can confront what it declares.
//
// It is the last resort: prefer measuring BEHAVIOUR. It is kept for the few guards whose
// subject is where a decision lives in the code, not what the command prints.
func leFonte(t *testing.T, nome string) string {
	t.Helper()
	b, err := os.ReadFile(nome)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// THE BACKFILL DOES NOT INVENT A LINK — a reading of the source, kept as a guard over its
// shape: the link comes from the one spelling of the under label, and the dry run exists.
// The behaviour itself (self-links, closed origins) is proven by the command tests below.
//
// MEASURED in the reference project: 29 open decisions, 5 recoverable links, 4 skipped. The
// 4 skipped are the right behaviour, not a coverage failure.
func TestBackfill_recoversTheLinkWithoutGuessing(t *testing.T) {
	cmd := newBackfillLabelsCmd()

	// The DRY RUN must exist: writing labels in bulk on a production board without being
	// able to look first is the kind of command one runs once and regrets.
	if cmd.Flags().Lookup("dry-run") == nil {
		t.Error("`--dry-run` is missing: this command writes into N issues at once, and " +
			"whoever runs it needs to see what it would do first")
	}

	// THE SPELLING IS ONE. There was a requirement to also read the previous one
	// (`sob-<n>`), and it left: accepting two forms of the same thing is surface for
	// divergence.
	source := leFonte(t, "backfill_labels.go")
	if !strings.Contains(source, "PrefixoLabelSob") {
		t.Error("the backfill does not read `PrefixoLabelSob` — it is where the link comes " +
			"from, and without it the command has nothing to recover")
	}
	if strings.Contains(source, "PrefixoLabelSobAntigo") {
		t.Error("the backfill reads the previous spelling again — one form only")
	}
}

// PRECEDENCE BETWEEN DECISIONS IS NOT INFERRED — but it MAY BE DECLARED.
//
// The guard lives in `backfill`, which infers from `under-<n>` (provenance). `escalate`
// DECLARES the link at the moment, and must not have the guard: suppressing it there would
// erase an order someone established. A source reading, because the subject is WHICH of the
// two commands holds the guard; the behaviour of the backfill side is
// TestBackfillLabelsCmd_writesOnlyTheUnequivocalLinks (#46).
func TestBackfill_doesNotInferPrecedenceBetweenDecisions(t *testing.T) {
	if g := "jaTem[initx.LabelNeedsUser]"; !strings.Contains(leFonte(t, "backfill_labels.go"), g) {
		t.Errorf("the backfill lost the guard %q — it would infer precedence from "+
			"`under-<n>` again, which states provenance: measured, wrong in 4 of 5 cases", g)
	}
	if strings.Contains(leFonte(t, "escalate.go"), "!jaEra") {
		t.Error("`escalate` got the backfill's guard — there the link is DECLARED by whoever " +
			"escalated, and suppressing it would erase the order it established")
	}
}

// backfillBoard is a board of open decisions that exercises every path of the command:
//
//	#701 under-44   → origin #44 is an OPEN work card: the unequivocal link
//	#702 under-45   → origin #45 is closed: skipped
//	#703 under-46   → origin #46 is itself an open decision: precedence is not inferred
//	#704 under-704  → a decision does not hold itself
//	#705 (no under) → cites "PR #88", whose body declares `Refs #50`: provenance recovered
//	#706 (no under) → cites "#99" without the word PR: prose, not a link
//	#707 (no under) → cites "PR #89", whose body declares no card: nothing to recover
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
			`{"number":706,"labels":[]},` +
			`{"number":707,"labels":[]}]`},
		ghRule{match: "issue view 705 *--json body*", out: "Seen while reviewing PR #88: the loop truncates."},
		ghRule{match: "issue view 706 *--json body*", out: "Related to #99 somehow."},
		ghRule{match: "issue view 707 *--json body*", out: "Seen in PR #89."},
		ghRule{match: "pr view 88 *", out: "Some summary\n\nRefs #50\n"},
		ghRule{match: "pr view 89 *", out: "A pull request that names no card.\n"},
		ghRule{match: "issue view 44 *--json state,labels*", out: `{"state":"OPEN","labels":[{"name":"anchors:in-progress"}]}`},
		ghRule{match: "issue view 45 *--json state,labels*", out: `{"state":"CLOSED","labels":[]}`},
		ghRule{match: "issue view 46 *--json state,labels*", out: `{"state":"OPEN","labels":[{"name":"` + initx.LabelNeedsUser + `"}]}`},
	)
}

func TestBackfillLabelsCmd_writesOnlyTheUnequivocalLinks(t *testing.T) {
	t.Run("BCLBB-B02: The open decisions are read in one listing", func(t *testing.T) {})
	t.Run("BCLBB-B04: The provenance is recovered from a cited pull request that declares a card", func(t *testing.T) {})
	t.Run("BCLBB-B05: Only an open origin card receives the block", func(t *testing.T) {})
	t.Run("BCLBB-B06: Precedence between decisions is not inferred", func(t *testing.T) {})
	t.Run("BCLBB-B08: A label is created before the card is edited with it", func(t *testing.T) {})
	t.Run("BCLBB-B10: The output counts what was recovered, written and skipped", func(t *testing.T) {})
	t.Run("BCLBB-I01: A decision under itself is never blocked by itself", func(t *testing.T) {})
	t.Run("BCLBB-X01: Nothing is removed and nothing unsupported is written", func(t *testing.T) {})
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
	if got := callsWith(c, "issue list"); len(got) != 1 {
		t.Errorf("the open decisions are read in ONE listing; got %v", got)
	}
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
	for _, n := range []string{"45", "46", "704", "706", "707"} {
		if len(callsWith(edits, "issue edit "+n+" ")) != 0 {
			t.Errorf("#%s must not be touched; edits: %v", n, edits)
		}
	}
	if got := callsWith(c, "remove-label"); len(got) != 0 {
		t.Errorf("the backfill never removes a label: %v", got)
	}
	create := indexOf(c, "label create "+initx.LabelBlockedBy("701"))
	if edit := indexOf(c, "issue edit 44 "); create < 0 || create > edit {
		t.Errorf("the label must be created before #44 is edited with it; calls: %v", c)
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

// Every decision holding an origin card is written — the claim holds the card while ANY is
// open — and the one already there is not written again.
func TestBackfillLabelsCmd_everyHolderOnceAndNotAgain(t *testing.T) {
	t.Run("BCLBB-B03: Every decision under an origin card holds it", func(t *testing.T) {})
	t.Run("BCLBB-B07: A blocked-by label already present is not written again", func(t *testing.T) {})
	under44 := `[{"name":"` + initx.LabelSob("44") + `"}]`
	board := ghRule{match: "issue list *", out: `[{"number":701,"labels":` + under44 + `},{"number":708,"labels":` + under44 + `}]`}

	calls := scriptedGH(t, board, ghRule{match: "issue view 44 *", out: `{"state":"OPEN","labels":[]}`})
	cmd := newBackfillLabelsCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{"--root", githubProject(t)})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	c := calls()
	for _, d := range []string{"701", "708"} {
		if len(callsWith(c, "issue edit 44 ", "--add-label "+initx.LabelBlockedBy(d))) != 1 {
			t.Errorf("#44 must receive blocked-by-%s; calls: %v", d, c)
		}
	}

	calls = scriptedGH(t, board, ghRule{match: "issue view 44 *",
		out: `{"state":"OPEN","labels":[{"name":"` + initx.LabelBlockedBy("701") + `"}]}`})
	cmd = newBackfillLabelsCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{"--root", githubProject(t)})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	edits := callsWith(calls(), "issue edit")
	if len(edits) != 1 || !strings.Contains(edits[0], initx.LabelBlockedBy("708")) {
		t.Errorf("only the missing blocked-by-708 is written; edits: %v", edits)
	}
}

// The dry run names what it would do and touches nothing.
func TestBackfillLabelsCmd_dryRunTouchesNothing(t *testing.T) {
	t.Run("BCLBB-B09: The dry run names what it would do and touches nothing", func(t *testing.T) {})
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
	t.Run("BCLBB-E01: A failed listing fails the command", func(t *testing.T) {})
	t.Run("BCLBB-E02: An unreadable listing fails the command", func(t *testing.T) {})
	t.Run("BCLBB-B01: backfill-labels in local mode is refused", func(t *testing.T) {})
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
	t.Run("BCLBB-E03: A failed edit is reported and skipped", func(t *testing.T) {})
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

// An origin card whose state cannot be read is not labelled: it might be closed.
func TestBackfillLabelsCmd_unreadableOriginIsSkipped(t *testing.T) {
	t.Run("BCLBB-E04: An origin card whose state cannot be read is skipped", func(t *testing.T) {})
	calls := scriptedGH(t,
		ghRule{match: "issue list *", out: `[{"number":701,"labels":[{"name":"` + initx.LabelSob("44") + `"}]}]`},
		ghRule{match: "issue view 44 *", out: "HTTP 502", code: 1},
	)
	cmd := newBackfillLabelsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--root", githubProject(t)})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := callsWith(calls(), "issue edit"); len(got) != 0 {
		t.Errorf("an origin whose state is unknown must not be edited: %v", got)
	}
	if !strings.Contains(out.String(), "0 link(s) written · 1 skipped") {
		t.Errorf("the unreadable origin is counted as skipped:\n%s", out.String())
	}
}

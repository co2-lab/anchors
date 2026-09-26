package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

const escalatedURL = "https://github.com/acme/app/issues/900"

// runEscalate runs `escalate` over a github project with the given gh rules, and returns
// stdout, stderr, the error and the recorded gh calls.
func runEscalate(t *testing.T, rules []ghRule, args ...string) (string, string, error, []string) {
	t.Helper()
	root := githubProject(t)
	calls := scriptedGH(t, append(rules, ghRule{match: "issue create *", out: escalatedURL})...)
	cmd := newEscalateCmd()
	cmd.SetArgs(append([]string{"--root", root}, args...))
	var err error
	var out string
	errOut := stderrOf(t, func() {
		out = stdoutOf(t, func() { err = cmd.Execute() })
	})
	return out, errOut, err, calls()
}

func onlyCall(t *testing.T, calls []string, fragments ...string) string {
	t.Helper()
	got := callsWith(calls, fragments...)
	if len(got) != 1 {
		t.Fatalf("expected one call with %q, got %v\nall calls: %v", fragments, got, calls)
	}
	return got[0]
}

// The ordinary card: born in to-do under the origin card, and the origin card goes on
// with a trace of it.
func TestEscalateCmd_ordinaryCardIsBornUnderTheOriginAndStopsNothing(t *testing.T) {
	out, _, err, calls := runEscalate(t, nil, "--card", "44", "the plan misses migrations\nmore detail")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	for _, want := range []string{"--title [plan] the plan misses migrations --body-file",
		"--label anchors", "--label anchors:to-do", "--label " + initx.LabelSob("44")} {
		if !strings.Contains(create, want) {
			t.Errorf("the issue must carry %q: %s", want, create)
		}
	}
	if strings.Contains(create, initx.LabelNeedsUser) {
		t.Errorf("an ordinary card waits for nobody: %s", create)
	}
	onlyCall(t, calls, "label create "+initx.LabelSob("44"))
	if len(callsWith(calls, "issue edit 44")) != 0 {
		t.Errorf("an ordinary finding must not stop the origin card: %v", calls)
	}
	if !strings.Contains(onlyCall(t, calls, "issue comment 44"), "Plan change recorded from this work — "+escalatedURL) {
		t.Error("the origin card keeps a trace of the finding")
	}
	if !strings.Contains(out, "finding recorded: "+escalatedURL) {
		t.Errorf("the output names the kind and the url:\n%s", out)
	}
}

// A decision stops the origin card with needs-user AND the link to the decision, and
// tells how to resume.
func TestEscalateCmd_decisionStopsTheCardAndLinksIt(t *testing.T) {
	out, _, err, calls := runEscalate(t, nil, "--card", "44", "--for-user", "which currency?")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	if !strings.Contains(create, "--title [decision] which currency?") || !strings.Contains(create, "--label "+initx.LabelNeedsUser) {
		t.Errorf("a decision is titled and labelled as one: %s", create)
	}
	onlyCall(t, calls, "label create "+initx.LabelBlockedBy("900"))
	edit := onlyCall(t, calls, "issue edit 44")
	if !strings.Contains(edit, "--add-label "+initx.LabelNeedsUser+","+initx.LabelBlockedBy("900")) {
		t.Errorf("the card waits for THIS decision: %s", edit)
	}
	if !strings.Contains(onlyCall(t, calls, "issue comment 44"), "Stopped: there is an open decision — "+escalatedURL) {
		t.Error("the origin card says why it stopped")
	}
	for _, want := range []string{"decision opened: " + escalatedURL, "card #44 stopped until the decision",
		"anchors decided --card 44"} {
		if !strings.Contains(out, want) {
			t.Errorf("the output lacks %q:\n%s", want, out)
		}
	}
}

// `--unsure` stops the card like a decision, but is titled and labelled as a framing
// question.
func TestEscalateCmd_unsureIsAFramingQuestion(t *testing.T) {
	_, _, err, calls := runEscalate(t, nil, "--card", "44", "--unsure", "is this mine?")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	for _, want := range []string{"--title [framing] is this mine?", "--label " + initx.LabelNeedsUser,
		"--label " + initx.LabelNeedsFraming} {
		if !strings.Contains(create, want) {
			t.Errorf("the framing question must carry %q: %s", want, create)
		}
	}
}

// A blocking bug holds the card by blocked-by alone — no needs-user, nobody to decide.
func TestEscalateCmd_blockingBugHoldsTheCardWithoutADecision(t *testing.T) {
	out, _, err, calls := runEscalate(t, nil, "--card", "44", "--bug", "--blocking", "claim picks the wrong card")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	if !strings.Contains(create, "--title [bug] claim picks the wrong card") || !strings.Contains(create, "--label "+initx.LabelBug) {
		t.Errorf("a bug is titled and labelled as one: %s", create)
	}
	edit := onlyCall(t, calls, "issue edit 44")
	if strings.Contains(edit, initx.LabelNeedsUser) || !strings.Contains(edit, "--add-label "+initx.LabelBlockedBy("900")) {
		t.Errorf("the card waits for the bug by blocked-by only: %s", edit)
	}
	if !strings.Contains(onlyCall(t, calls, "issue comment 44"), "a bug in the pipeline or the tool blocks this card") {
		t.Error("the origin card says a bug blocks it")
	}
	if !strings.Contains(out, "bug reported: "+escalatedURL) || !strings.Contains(out, "stopped until the bug is fixed") {
		t.Errorf("the output must say the card waits for the fix:\n%s", out)
	}
}

func TestEscalateCmd_nonBlockingBugLetsTheCardGoOn(t *testing.T) {
	_, _, err, calls := runEscalate(t, nil, "--card", "44", "--bug", "gate misreads a file")
	if err != nil {
		t.Fatal(err)
	}
	if len(callsWith(calls, "issue edit 44")) != 0 {
		t.Errorf("a non-blocking bug stops nothing: %v", calls)
	}
	if !strings.Contains(onlyCall(t, calls, "issue comment 44"), "Bug reported from this work (the card goes on)") {
		t.Error("the origin card keeps a trace of the bug")
	}
}

// A card that could not be labelled is warned about, not silently left free.
func TestEscalateCmd_failedLabelOnTheCardIsWarned(t *testing.T) {
	out, _, err, _ := runEscalate(t, []ghRule{{match: "issue edit 44 *", code: 1}},
		"--card", "44", "--for-user", "which currency?")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "could not label card #44") {
		t.Errorf("the failed label must be warned:\n%s", out)
	}
}

// The card of the reviewed PR is read from its `Refs`, and the PR is kept as provenance.
func TestEscalateCmd_reviewingPRDerivesTheCard(t *testing.T) {
	_, errOut, err, calls := runEscalate(t,
		[]ghRule{{match: "pr view 556 *", out: "Summary\nCloses #198\nRefs #200"}},
		"--reviewing-pr", "#556", "the loop truncates")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	if !strings.Contains(create, "--label "+initx.LabelDePR("556")) || !strings.Contains(create, "--label "+initx.LabelSob("198")) {
		t.Errorf("the finding is born under the PR's FIRST card, with the PR as provenance: %s", create)
	}
	if !strings.Contains(errOut, "PR #556 declares card #198") {
		t.Errorf("the derivation is said: %q", errOut)
	}
}

func TestEscalateCmd_withoutProvenanceItSaysSo(t *testing.T) {
	_, errOut, err, calls := runEscalate(t,
		[]ghRule{{match: "pr view 556 *", out: "no link line here, only #12 in prose"}},
		"--reviewing-pr", "556", "the loop truncates")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(errOut, "PR #556 declares no card") || !strings.Contains(errOut, "WITHOUT provenance") {
		t.Errorf("a finding with no origin must be announced: %q", errOut)
	}
	if strings.Contains(onlyCall(t, calls, "issue create"), initx.PrefixoLabelSob) {
		t.Error("no card was found, so no under- label may be invented")
	}
}

// With two cards in hand, the command refuses to guess under which one it was born.
func TestEscalateCmd_twoCardsInHandIsAmbiguous(t *testing.T) {
	root := githubProject(t)
	calls := scriptedGH(t, ghRule{match: "issue list *--json number,title,labels,comments*",
		out: "44\tfirst\tanchors:in-progress\n45\tsecond\tanchors:in-progress"})
	t.Setenv("ANCHORS_AGENT", "host/dev1")
	cmd := newEscalateCmd()
	cmd.SetArgs([]string{"--root", root, "something"})
	var err error
	stderrOf(t, func() { err = cmd.Execute() })
	if err == nil || !strings.Contains(err.Error(), "2 cards in hand (#44, #45)") {
		t.Fatalf("two cards in hand must be refused, got %v", err)
	}
	if len(callsWith(calls(), "issue create")) != 0 {
		t.Error("nothing may be created on an ambiguous origin")
	}
}

func TestEscalateCmd_theOneCardInHandIsTheOrigin(t *testing.T) {
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue list *--json number,title,labels,comments*", out: "44\tfirst\tanchors:in-progress"},
		ghRule{match: "issue create *", out: escalatedURL},
	)
	t.Setenv("ANCHORS_AGENT", "host/dev1")
	cmd := newEscalateCmd()
	cmd.SetArgs([]string{"--root", root, "something"})
	var err error
	errOut := stderrOf(t, func() { stdoutOf(t, func() { err = cmd.Execute() }) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(errOut, "born under card #44 (from `anchors-owner`)") {
		t.Errorf("the discovered origin must be said: %q", errOut)
	}
	if !strings.Contains(onlyCall(t, calls(), "issue create"), "--label "+initx.LabelSob("44")) {
		t.Error("the finding must be born under the card in hand")
	}
}

// A target that already has open cards is warned about, up to three, before creating.
func TestEscalateCmd_warnsWhenTheTargetAlreadyHasCards(t *testing.T) {
	cards := `[{"number":1,"title":"a src/x.ts","body":""},{"number":2,"title":"b","body":"on src/x.ts"},` +
		`{"number":3,"title":"c src/x.ts","body":""},{"number":4,"title":"d src/x.ts","body":""},` +
		`{"number":5,"title":"unrelated","body":"src/x.tsx"}]`
	_, errOut, err, _ := runEscalate(t, []ghRule{{match: "issue list *--search src/x.ts*", out: cards}},
		"--card", "44", "--about", "src/x.ts", "dup")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"#1 a src/x.ts", "#2 b", "#3 c src/x.ts"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the warning must list %q: %q", want, errOut)
		}
	}
	if strings.Contains(errOut, "#4 d") {
		t.Errorf("only three are listed, the rest are counted: %q", errOut)
	}
}

func TestEscalateCmd_refusals(t *testing.T) {
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--bug", "--for-user", "x"}, "`--bug` is not a decision"},
		{[]string{"--bug", "--unsure", "x"}, "`--bug` is not a decision"},
		{[]string{"--blocking", "x"}, "`--blocking` goes with `--bug`"},
	} {
		_, _, err, calls := runEscalate(t, nil, c.args...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: want %q, got %v", c.args, c.want, err)
		}
		if len(calls) != 0 {
			t.Errorf("%v: a refusal calls nothing, got %v", c.args, calls)
		}
	}

	scriptedGH(t)
	cmd := newEscalateCmd()
	cmd.SetArgs([]string{"--root", localProject(t), "x"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode must be refused, got %v", err)
	}

	_, _, err, _ := runEscalate(t, []ghRule{{match: "issue create *", out: "label not found", code: 1}}, "--card", "44", "x")
	if err == nil || !strings.Contains(err.Error(), "open the issue") || !strings.Contains(err.Error(), "label not found") {
		t.Errorf("a failed create must fail with gh's answer, got %v", err)
	}
}

func TestNumeroDaIssue(t *testing.T) {
	for in, want := range map[string]string{
		"https://github.com/acme/app/issues/900":   "900",
		"https://github.com/acme/app/issues/900\n": "900",
		"https://github.com/acme/app/issues/abc":   "",
		"https://github.com/acme/app/issues/":      "",
		"no url":                                   "",
	} {
		if got := numeroDaIssue(in); got != want {
			t.Errorf("numeroDaIssue(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCardDoPR_noRepoOrPRAsksNothing(t *testing.T) {
	calls := scriptedGH(t)
	if cardDoPR("", "5") != "" || cardDoPR("acme/app", " # ") != "" {
		t.Error("without repo or PR there is no card")
	}
	if len(calls()) != 0 {
		t.Errorf("nothing to ask gh: %v", calls())
	}
	scriptedGH(t, ghRule{match: "pr view *", code: 1})
	if got := cardDoPR("acme/app", "5"); got != "" {
		t.Errorf("a failed read yields no card, got %q", got)
	}
}

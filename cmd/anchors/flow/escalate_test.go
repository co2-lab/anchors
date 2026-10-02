// @anchors
//   ref: SCLTE

package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// THE ISSUE must answer three things, and the third is the one usually missing: how to
// unblock. Without it the card stays stopped waiting for someone to guess the protocol.
func TestEscalationBody_saysWhyItStoppedAndHowToUnblock(t *testing.T) {
	t.Run("SCLTE-B10: The body says why and how to go on for each exit", func(t *testing.T) {})
	body := escalationBody("The spec asks for a cache; the plan said there would be none.",
		"plans/0001-foundation.md", "12", true, false)

	for _, required := range []string{
		"The spec asks for a cache", // the reason, in the words of whoever saw it
		"plans/0001-foundation.md",  // where
		"R0001",                     // how to record the decision
		initx.LabelNeedsUser,        // what to remove to unblock
		"#12",                       // where the work stopped
		"not the agent's",           // why it stopped here
	} {
		if !strings.Contains(body, required) {
			t.Errorf("the escalation body must contain %q; got:\n%s", required, body)
		}
	}

	// A framing question asks about the framing first, and says how to hand it back.
	framing := escalationBody("Is this mine?", "", "12", true, true)
	for _, required := range []string{"THE FIRST QUESTION IS THE FRAMING", "`anchors:to-do`", "Work stopped on card #12"} {
		if !strings.Contains(framing, required) {
			t.Errorf("the framing body must contain %q; got:\n%s", required, framing)
		}
	}
}

// Without `--card` the command still serves: not every incoherence is found with a card in
// hand (a reviewer reading a PR, for example). The body must not cite a card that does not
// exist.
func TestEscalationBody_withoutCardInventsNoReference(t *testing.T) {
	body := escalationBody("The plan contradicts itself between F02 and F04.", "", "", true, false)
	if strings.Contains(body, "#") && strings.Contains(body, "Work stopped") {
		t.Errorf("without --card no card may be cited; got:\n%s", body)
	}
	if !strings.Contains(body, "contradicts itself") {
		t.Error("the reason must survive even without a card")
	}
}

// THE TITLE of the issue is one line. A long reason must not become an unreadable title in
// the issue list, which is where the user will find it.
func TestEscalate_titleFitsInOneLine(t *testing.T) {
	t.Run("SCLTE-B07: The title is the exit's prefix and the reason's first line", func(t *testing.T) {})
	long := strings.Repeat("a very detailed explanation of the incoherence ", 5)
	got := firstLineOfReason(long)
	if len(got) > 70 {
		t.Errorf("title with %d chars; it should fit in 70: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("the cut must signal there is more: %q", got)
	}
	// A reason with several lines: the title is the first.
	if got := firstLineOfReason("first line\nsecond line"); got != "first line" {
		t.Errorf("the title is the first line, got %q", got)
	}
	// And every exit puts its prefix before it.
	for flags, prefix := range map[string]string{"": "[plan] ", "--for-user": "[decision] ",
		"--unsure": "[framing] ", "--bug": "[bug] "} {
		args := []string{"--card", "44"}
		if flags != "" {
			args = append(args, flags)
		}
		_, _, err, calls := runEscalate(t, nil, append(args, "the plan misses migrations\nmore")...)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(onlyCall(t, calls, "issue create"), "--title "+prefix+"the plan misses migrations --body-file") {
			t.Errorf("%q: the title must be %q + the first line", flags, prefix)
		}
	}
}

// THE DEFAULT EXIT is an ordinary card, and it must look visibly different from a decision:
// if the two issues read the same, whoever opens the list would not know which one waits
// for them.
func TestEscalationBody_ordinaryCardAsksForNoDecision(t *testing.T) {
	body := escalationBody("The plan does not cover configuring and running migrations.",
		"plans/0001-foundation.md", "12", false, false)

	if strings.Contains(body, initx.LabelNeedsUser) {
		t.Errorf("an ordinary card must not ask to remove the decision label; got:\n%s", body)
	}
	if strings.Contains(body, "Work stopped") {
		t.Errorf("an ordinary card does NOT stop the work; got:\n%s", body)
	}
	// And it must teach the emergency exit: whoever works on it may find out, while working,
	// that the change was bigger than whoever opened it judged.
	if !strings.Contains(body, "--for-user") {
		t.Errorf("the ordinary card must say what to do if the change turns out to be of direction; "+
			"got:\n%s", body)
	}
	if !strings.Contains(body, "migrations") {
		t.Error("the reason must survive in the default exit")
	}
}

// The trigger is not only incoherence. A GAP (a coherent but incomplete plan) follows the
// same flow, and the text must not suggest it only serves contradictions.
func TestEscalationBody_doesNotPresumeIncoherence(t *testing.T) {
	for _, forUser := range []bool{true, false} {
		body := escalationBody("The plan did not foresee migrations.", "plans/0001.md", "", forUser, false)
		if strings.Contains(body, "correction would change") {
			t.Errorf("the body must not presume there was an error to correct (for-user=%v):\n%s",
				forUser, body)
		}
	}
}

// THE LINK to the origin work is a LABEL, not text in the body — a sentence cannot be
// queried. And the body NAMES the relation, so whoever reads the issue knows it is not loose.
func TestEscalationBody_theLinkToTheOriginWorkIsALabel(t *testing.T) {
	if got := initx.LabelSob("44"); got != "anchors:under-44" {
		t.Fatalf("the label links the finding to the card; got %q", got)
	}
	body := escalationBody("jest only measures src/", "jest.config.js", "44", false, false)
	if !strings.Contains(body, "anchors:under-44") {
		t.Errorf("the body must cite the label, or whoever reads it cannot find the siblings;\n%s", body)
	}
	if !strings.Contains(body, "same PR") {
		t.Error("the body must say the two are delivered together")
	}
}

// THE LINK IS WITH THE CARD, AND THE PR IS THE WAY — the command reads the PR's own
// declaration instead of the agent searching for it. Only at the start of a line: a number
// cited in prose is not a link.
func TestEscalate_theLinkOfThePRComesFromItsBody(t *testing.T) {
	for _, c := range []struct {
		name, body, want string
	}{
		{"what pr-body generates", "text\n\nRefs #735\n", "735"},
		{"the closing word", "Closes #198", "198"},
		{"upper or lower case", "closes #42", "42"},
		{"fixes too", "Fixes #7", "7"},
		{"cited in prose does not count", "the card refs #99 is another's", ""},
		{"no link at all", "only a description", ""},
	} {
		m := vinculoNoCorpoRE.FindStringSubmatch(c.body)
		got := ""
		if m != nil {
			got = m[1]
		}
		if got != c.want {
			t.Errorf("%s: from %q expected %q, got %q", c.name, c.body, c.want, got)
		}
	}
}

// THE TWO LINKS COEXIST — a reading of the source kept as a guard over its shape: the PR
// block comes before the card block (the card is DERIVED from the PR), and neither sits in
// the other's `else`. The behaviour is TestEscalateCmd_reviewingPRDerivesTheCard.
func TestEscalate_keepsTheCardAndThePRTogether(t *testing.T) {
	source := leFonte(t, "escalate.go")
	for _, piece := range []string{"initx.LabelDePR(", "initx.LabelSob(card)"} {
		if !strings.Contains(source, piece) {
			t.Errorf("`escalate` does not apply %q — one of the two links is lost", piece)
		}
	}
	iPR := strings.Index(source, `if revisandoPR != "" {`)
	iCard := strings.Index(source, `if card != "" {
				labels = append(labels, initx.LabelSob(card))`)
	if iPR < 0 || iCard < 0 {
		t.Fatal("the two link blocks were not found — the command changed shape")
	}
	if iPR > iCard {
		t.Error("the PR block comes after the card block: the card is DERIVED from the PR")
	}
}

// A bug is not a decision (reference app, 2026-09-24: six "decisions" in two hours, all bugs).
// The body says so first, never asks for one, and names the card only as the flag says.
func TestBugBody_isNotADecision(t *testing.T) {
	b := bugBody("the review job picks the wrong card", ".github/workflows/anchors-pr-checks.yml", "966", true)
	for _, want := range []string{"nothing to decide", "anchors-pr-checks.yml", "Card #966 waits for it", "close this issue"} {
		if !strings.Contains(b, want) {
			t.Errorf("the bug body should say %q:\n%s", want, b)
		}
	}
	for _, bad := range []string{"not the agent's", "needs-user", "decide, and record"} {
		if strings.Contains(b, bad) {
			t.Errorf("a bug body must not read as a decision (%q):\n%s", bad, b)
		}
	}
	if nb := bugBody("x", "", "966", false); strings.Contains(nb, "waits for it") || !strings.Contains(nb, "which goes on") {
		t.Errorf("without --blocking the card goes on, and the body must say so:\n%s", nb)
	}
}

// THE FINDING IS BORN WITH PROVENANCE, and when it is not the command SAYS so — a reading
// of the source kept as a guard: the discovery reuses `pr-body`'s piece, ambiguity refuses,
// and the loose finding is announced. The behaviour is TestEscalateCmd_twoCardsInHandIsAmbiguous,
// TestEscalateCmd_theOneCardInHandIsTheOrigin and TestEscalateCmd_withoutProvenanceItSaysSo.
func TestEscalate_discoversTheOriginCard(t *testing.T) {
	source := leFonte(t, "escalate.go")
	if !strings.Contains(source, `requestedCards("", cfg)`) {
		t.Error("`escalate` does not discover the origin card through `pr-body`'s piece")
	}
	if !strings.Contains(source, "provide `--card <n>`") {
		t.Error("`escalate` does not handle several cards in hand")
	}
	if !strings.Contains(source, "anchors-owner` — the finding is born WITHOUT ") {
		t.Error("`escalate` accepts a finding without provenance in silence")
	}
}

// Without a reason there is nothing to title the card with.
func TestEscalateCmd_requiresAReason(t *testing.T) {
	t.Run("SCLTE-B01: escalate without a reason is refused", func(t *testing.T) {})
	root := githubProject(t)
	calls := scriptedGH(t)
	cmd := newEscalateCmd()
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	cmd.SetArgs([]string{"--root", root, "--card", "44"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("escalate without a reason must fail")
	}
	if got := callsWith(calls(), "issue create"); len(got) != 0 {
		t.Errorf("no card may be created: %v", got)
	}
}

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
	t.Run("SCLTE-B09: The origin card and the reviewed pull request are both labels", func(t *testing.T) {})
	t.Run("SCLTE-B15: An ordinary card and a non-blocking bug let the origin card go on", func(t *testing.T) {})
	t.Run("SCLTE-B12: The output names the kind and the address", func(t *testing.T) {})
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

// `#44` is how people write a card; it reached the label as `under-#44`, which no
// `--label anchors:under-44` query finds.
func TestEscalateCmd_cardWithHashIsTheSameCard(t *testing.T) {
	t.Run("SCLTE-B17: A card written with its hash is the same card", func(t *testing.T) {})
	_, _, err, calls := runEscalate(t, nil, "--card", "#44", "--for-user", "which currency?")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	if !strings.Contains(create, "--label "+initx.LabelSob("44")) || strings.Contains(create, "#44") {
		t.Errorf("the new card must carry %s: %s", initx.LabelSob("44"), create)
	}
	onlyCall(t, calls, "label create "+initx.LabelSob("44"))
	onlyCall(t, calls, "issue edit 44 ")
}

// A decision stops the origin card with needs-user AND the link to the decision, and
// tells how to resume.
func TestEscalateCmd_decisionStopsTheCardAndLinksIt(t *testing.T) {
	t.Run("SCLTE-B13: A decision stops the origin card and prints the way back", func(t *testing.T) {})
	t.Run("SCLTE-B08: The labels follow the exit", func(t *testing.T) {})
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
	// `--unsure` stops the card like a decision.
	edit := onlyCall(t, calls, "issue edit 44")
	if !strings.Contains(edit, "--add-label "+initx.LabelNeedsUser+","+initx.LabelBlockedBy("900")) {
		t.Errorf("a framing question stops the card like a decision: %s", edit)
	}
}

// A blocking bug holds the card by blocked-by alone — no needs-user, nobody to decide.
func TestEscalateCmd_blockingBugHoldsTheCardWithoutADecision(t *testing.T) {
	t.Run("SCLTE-B14: A blocking bug holds the card by blocked-by alone", func(t *testing.T) {})
	t.Run("SCLTE-X01: A bug never carries needs-user", func(t *testing.T) {})
	out, _, err, calls := runEscalate(t, nil, "--card", "44", "--bug", "--blocking", "claim picks the wrong card")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create")
	if !strings.Contains(create, "--title [bug] claim picks the wrong card") || !strings.Contains(create, "--label "+initx.LabelBug) ||
		strings.Contains(create, initx.LabelNeedsUser) {
		t.Errorf("a bug is titled and labelled as one, and waits for nobody's decision: %s", create)
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
	out, _, err, calls := runEscalate(t, nil, "--card", "44", "--bug", "gate misreads a file")
	if err != nil {
		t.Fatal(err)
	}
	if len(callsWith(calls, "issue edit 44")) != 0 {
		t.Errorf("a non-blocking bug stops nothing: %v", calls)
	}
	if !strings.Contains(onlyCall(t, calls, "issue comment 44"), "Bug reported from this work (the card goes on)") {
		t.Error("the origin card keeps a trace of the bug")
	}
	if !strings.Contains(out, "bug reported: "+escalatedURL) {
		t.Errorf("the output names the kind and the url:\n%s", out)
	}
}

// A card that could not be labelled is warned about, not silently left free.
func TestEscalateCmd_failedLabelOnTheCardIsWarned(t *testing.T) {
	t.Run("SCLTE-E02: A card that cannot be stopped is warned about", func(t *testing.T) {})
	out, _, err, _ := runEscalate(t, []ghRule{{match: "issue edit 44 *", code: 1}},
		"--card", "44", "--for-user", "which currency?")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "could not label card #44") || strings.Contains(out, "stopped until the decision") {
		t.Errorf("the failed label must be warned, and the card not announced as stopped:\n%s", out)
	}
}

// The card of the reviewed PR is read from its `Refs`, and the PR is kept as provenance.
func TestEscalateCmd_reviewingPRDerivesTheCard(t *testing.T) {
	t.Run("SCLTE-B04: The card of the reviewed pull request is the origin", func(t *testing.T) {})
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
	t.Run("SCLTE-B06: A finding with no origin is created and warned", func(t *testing.T) {})
	t.Run("SCLTE-I01: Without an origin card no other card is touched", func(t *testing.T) {})
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
	if got := append(callsWith(calls, "issue edit"), callsWith(calls, "issue comment")...); len(got) != 0 {
		t.Errorf("without an origin card no other card is touched: %v", got)
	}
}

// With two cards in hand, the command refuses to guess under which one it was born.
func TestEscalateCmd_twoCardsInHandIsAmbiguous(t *testing.T) {
	t.Run("SCLTE-B05: The one card in hand is the origin, and two are ambiguous", func(t *testing.T) {})
	root := githubProject(t)
	calls := scriptedGH(t, ghRule{match: "issue list *--json number,title,labels,comments*",
		out: "44\tfirst\tanchors:in-progress\n45\tsecond\tanchors:in-progress"})
	t.Setenv("ANCHORS_AGENT", "host/dev1")
	cmd := newEscalateCmd()
	cmd.SetArgs([]string{"--root", root, "something"})
	var err error
	stderrOf(t, func() { err = cmd.Execute() })
	if err == nil || !strings.Contains(err.Error(), "2 cards in hand (#44, #45)") || !strings.Contains(err.Error(), "--card <n>") {
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
	t.Run("SCLTE-B11: A target with open cards is warned about, three at most", func(t *testing.T) {})
	cards := `[{"number":1,"title":"a src/x.ts","body":""},{"number":2,"title":"b","body":"on src/x.ts"},` +
		`{"number":3,"title":"c src/x.ts","body":""},{"number":4,"title":"d src/x.ts","body":""},` +
		`{"number":5,"title":"unrelated","body":"src/x.tsx"}]`
	_, errOut, err, calls := runEscalate(t, []ghRule{{match: "issue list *--search src/x.ts*", out: cards}},
		"--card", "44", "--about", "src/x.ts", "dup")
	if err != nil {
		t.Fatal(err)
	}
	onlyCall(t, calls, "issue list", "--repo acme/app") // the project's board, not the cwd's
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
	t.Run("SCLTE-B02: Contradictory exits are refused before any call", func(t *testing.T) {})
	t.Run("SCLTE-B03: escalate in local mode is refused", func(t *testing.T) {})
	t.Run("SCLTE-E01: A card the platform refuses to create fails the command", func(t *testing.T) {})
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
	t.Run("SCLTE-B16: The new card's number is read only from an all-digit address tail", func(t *testing.T) {})
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

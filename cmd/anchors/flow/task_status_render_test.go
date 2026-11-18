// @anchors
//   code: TSRTT
//   ref: TSRTS

package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/i18n"
)

// The command exists because of a report that was RIGHT and insufficient: an agent
// diagnosed a CI failure, fixed it, pushed, and ended its turn with "waiting for the new
// run" — with the card `in-progress` under its name and the check's verdict unread.
//
// These tests charge what the format must say in the states where silence costs.

func testCard(n int, state, title string) *board.Card {
	return &board.Card{Number: n, Title: title, State: state}
}

// prLine returns the report's pull request line.
func prLine(out string) string {
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, "PR ") {
			return l
		}
	}
	return ""
}

func TestTaskStatus_sectionsComeInTheOrderThatDecides(t *testing.T) {
	t.Run("TSRTS-B01: The sections come in the order that decides", func(t *testing.T) {})
	t.Run("TSRTS-X01: Rendering looks nothing up", func(t *testing.T) {})
	calls := scriptedGH(t)
	out := renderTaskStatus(taskState{
		Card:     testCard(303, "anchors:in-progress", "[INDTN] implement spec"),
		Reverted: []string{"Reverted: closed by hand"},
		PR:       &branchPR{Number: 368, State: "OPEN", Total: 1, Checks: map[string]int{checkPassed: 1}},
		Branch:   "impl-x", Clean: true,
		Blocked: []board.Card{{Number: 99, Title: "[DEC1] which vocabulary?"}},
	})
	last := -1
	for _, section := range []string{"Task  #303", "UNDONE", "\nPR ", "\nGit ", "Waiting on a person's decision",
		"What I proved", "What was left out", "\nNext"} {
		i := strings.Index(out, section)
		if i < 0 || i < last {
			t.Fatalf("section %q is missing or out of order:\n%s", section, out)
		}
		last = i
	}
	if c := calls(); len(c) != 0 {
		t.Errorf("rendering must look nothing up: %v", c)
	}
}

func TestTaskStatus_cardLineDropsTheCodeAndTranslatesTheState(t *testing.T) {
	t.Run("TSRTS-B02: The card line drops the code and translates the state", func(t *testing.T) {})
	c := testCard(301, "anchors:in-review", "[X] screen")
	c.Owner = "host/dev1"
	out := renderTaskStatus(taskState{Card: c, Branch: "develop", Clean: true})
	for _, want := range []string{"Task  #301 · screen\n", "under review", "owner: host/dev1"} {
		if !strings.Contains(out, want) {
			t.Errorf("the card line lacks %q:\n%s", want, out)
		}
	}
	if out := renderTaskStatus(taskState{Branch: "develop", Clean: true}); !strings.Contains(out, "no card found — provide `--card N`") {
		t.Errorf("without a card the report says how to name one:\n%s", out)
	}
}

func TestTaskStatus_prWithoutChecksDoesNotPassAsChecked(t *testing.T) {
	t.Run("TSRTS-B04: The pull request line never hides missing or mixed checks", func(t *testing.T) {})
	// The case measured: PR #368 opened and NO check ran. "PR #368 open" alone reads as
	// checked work — the green that did not do the work.
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "[INDTN] implement spec"),
		Branch: "impl-x", Clean: true,
		PR: &branchPR{Number: 368, State: "OPEN", Total: 0, Checks: map[string]int{}},
	})
	// Two assertions for TWO different places: the PR line is what a fast reader sees, and
	// the next step is what the agent follows. The next step repeats the phrase, so one
	// search over the whole text would cover a missing PR line.
	if line := prLine(out); !strings.Contains(line, "no check ran") {
		t.Errorf("the PR LINE must say no check ran; it was %q", line)
	}
	if !strings.Contains(out, "do not trust the green that does not exist") {
		t.Errorf("the next step should refuse the PR with no check; output:\n%s", out)
	}
	if line := prLine(renderTaskStatus(taskState{PR: &branchPR{Number: 5, State: "OPEN", Total: 4,
		Checks: map[string]int{checkPassed: 4}}})); !strings.Contains(line, "checks: 4/4 passed") {
		t.Errorf("one class alone reads as count over total: %q", line)
	}
	if line := prLine(renderTaskStatus(taskState{})); line != "PR    none for this branch" {
		t.Errorf("without a PR the line says so: %q", line)
	}
}

// The check classes are printed in the user's language. They were the Portuguese words
// "passou", "em curso" and "reprovou" whatever the language, since they were the map keys.
func TestTaskStatus_checkClassesSpeakTheUsersLanguage(t *testing.T) {
	t.Run("TSRTS-B12: The check classes are printed in the user's language", func(t *testing.T) {})
	prev := i18n.Current()
	t.Cleanup(func() { _ = i18n.Set(prev) })
	e := taskState{PR: &branchPR{Number: 5, State: "OPEN", Total: 3,
		Checks: map[string]int{checkPassed: 1, checkRunning: 1, checkFailed: 1}}}
	for lang, want := range map[string]string{
		"en":    "1 failed, 1 running, 1 passed",
		"pt-BR": "1 reprovou, 1 em curso, 1 passou",
		"es":    "1 falló, 1 en curso, 1 pasó",
	} {
		if err := i18n.Set(lang); err != nil {
			t.Fatal(err)
		}
		if line := prLine(renderTaskStatus(e)); !strings.Contains(line, want) {
			t.Errorf("%s: the classes must read %q: %q", lang, want, line)
		}
	}
}

func TestTaskStatus_runningCheckIsNotAPassedCheck(t *testing.T) {
	t.Run("TSRTS-I01: A running check never reads as passed", func(t *testing.T) {})
	// The exact confusion that ended the turn wrongly. A report that adds "running" to
	// "passed" produces "4/4 passed" for a PR that has not finished running.
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "x"),
		Branch: "impl-x", Clean: true,
		PR: &branchPR{Number: 368, State: "OPEN", Total: 4,
			Checks: map[string]int{checkPassed: 3, checkRunning: 1}},
	})
	if strings.Contains(out, "4/4") {
		t.Errorf("3 passed and 1 is running: it is not 4/4; output:\n%s", out)
	}
	if !strings.Contains(prLine(out), "1 running, 3 passed") {
		t.Errorf("the running check must appear, before the passed ones; output:\n%s", out)
	}
}

func TestTaskStatus_failedCheckIsWorkOfThisCard(t *testing.T) {
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "x"),
		Branch: "impl-x", Clean: true,
		PR: &branchPR{Number: 368, State: "OPEN", Total: 4,
			Checks: map[string]int{checkPassed: 3, checkFailed: 1}},
	})
	if !strings.Contains(out, "work of THIS card") {
		t.Errorf("the red is not a new card; output:\n%s", out)
	}
	// The order matters: a fast reader must hit the problem, not what went right.
	if line := prLine(out); strings.Index(line, "failed") > strings.Index(line, "passed") {
		t.Errorf("the failure should come before what passed; line: %q", line)
	}
}

// The next step with an open pull request follows its checks.
func TestTaskStatus_nextStepFollowsTheChecks(t *testing.T) {
	t.Run("TSRTS-B10: With an open pull request the step follows the checks", func(t *testing.T) {})
	card := testCard(303, "anchors:in-progress", "x")
	for _, c := range []struct {
		checks map[string]int
		total  int
		want   []string
	}{
		{map[string]int{}, 0, []string{"no check ran on PR #368", "do not trust the green that does not exist"}},
		{map[string]int{checkPassed: 3, checkRunning: 1}, 4, []string{"`gh pr checks 368 --watch`", "do not end the turn here"}},
		{map[string]int{checkPassed: 3, checkFailed: 1}, 4, []string{"1 check(s) failed on PR #368", "work of THIS card"}},
		{map[string]int{checkPassed: 4}, 4, []string{"the checks of PR #368 passed — the review is missing"}},
	} {
		steps := strings.Join(nextStep(taskState{Card: card, Clean: true,
			PR: &branchPR{Number: 368, State: "OPEN", Total: c.total, Checks: c.checks}}), "\n")
		for _, want := range c.want {
			if !strings.Contains(steps, want) {
				t.Errorf("checks %v: the next step lacks %q:\n%s", c.checks, want, steps)
			}
		}
	}
}

func TestTaskStatus_doesNotAdviseOpeningAPRForACardInReview(t *testing.T) {
	t.Run("TSRTS-B09: Without a pull request the step depends on where the card is", func(t *testing.T) {})
	// The first version said "open the PR" for any open card without a PR on the branch,
	// and for an `in-review` card that is wrong advice — the work was already delivered.
	for _, state := range []string{"anchors:in-review", "anchors:ready-to-review"} {
		out := renderTaskStatus(taskState{Card: testCard(301, state, "[X] screen"), Branch: "develop", Clean: true})
		if strings.Contains(out, "open the PR") {
			t.Errorf("%s: a card in review has no PR to open; output:\n%s", state, out)
		}
		if !strings.Contains(out, "`gh pr list --search 301` finds it") {
			t.Errorf("%s: the report must say where its PR is; output:\n%s", state, out)
		}
	}
	if steps := strings.Join(nextStep(taskState{Card: testCard(303, "anchors:in-progress", "x"), Clean: true}), "\n"); !strings.Contains(steps, "open the PR") ||
		!strings.Contains(steps, "writes the lines that link the card") || strings.Contains(steps, "closing") {
		t.Errorf("a card in progress on a clean pushed tree is told to open the PR, whose lines link the card without closing it:\n%s", steps)
	}
	if steps := strings.Join(nextStep(taskState{Card: testCard(303, "anchors:to-do", "x"), Clean: true}), "\n"); strings.Contains(steps, "open the PR") {
		t.Errorf("a card nobody took has no PR to open:\n%s", steps)
	}
}

func TestTaskStatus_closedAndIdleStates(t *testing.T) {
	t.Run("TSRTS-B11: A closed card and an idle state have their steps", func(t *testing.T) {})
	if steps := strings.Join(nextStep(taskState{Card: testCard(9, "closed", "x"), Clean: true}), "\n"); !strings.Contains(steps, "`anchors next` asks the claim pipeline for the next one") {
		t.Errorf("a closed card is told to ask for the next one:\n%s", steps)
	}
	if steps := strings.Join(nextStep(taskState{Clean: true}), "\n"); !strings.Contains(steps, "`anchors status` says where the project is") {
		t.Errorf("with nothing to do the fallback applies:\n%s", steps)
	}
}

func TestTaskStatus_escalatedCardsAppearApart(t *testing.T) {
	t.Run("TSRTS-B06: The decisions waiting on a person have their own section", func(t *testing.T) {})
	// `needs-user` is the only pending work that does NOT move by continuing to work. A
	// report that omits it invites the reader to wait for something only they can unblock.
	out := renderTaskStatus(taskState{
		Card:   testCard(303, "anchors:in-progress", "x"),
		Branch: "impl-x", Clean: true,
		Blocked: []board.Card{
			{Number: 99, Title: "[DEC1] which vocabulary?"},
			{Number: 12, Title: "[DEC2] where is the limit?"},
		},
	})
	if !strings.Contains(out, "Waiting on a person's decision") {
		t.Errorf("the escalated cards need a section of their own; output:\n%s", out)
	}
	if !strings.Contains(out, "#99    which vocabulary?") || !strings.Contains(out, "#12    where is the limit?") {
		t.Errorf("both numbers must appear, without their codes; output:\n%s", out)
	}
	if !strings.Contains(out, "no agent resolves these") {
		t.Errorf("the format must say WHY these do not move; output:\n%s", out)
	}
}

func TestTaskStatus_withoutEscalatedCardsNoSectionIsInvented(t *testing.T) {
	out := renderTaskStatus(taskState{
		Card: testCard(303, "anchors:in-progress", "x"), Branch: "impl-x", Clean: true,
	})
	if strings.Contains(out, "Waiting on a person's decision") {
		t.Errorf("without escalated cards the section should not exist; output:\n%s", out)
	}
}

func TestTaskStatus_theTwoGapsAreExplicit(t *testing.T) {
	t.Run("TSRTS-B07: The two gaps are always explicit", func(t *testing.T) {})
	// What the machine does not know, it ASKS for. A named gap gets filled; an absent
	// section goes unnoticed.
	out := renderTaskStatus(taskState{Branch: "develop", Clean: true})
	for _, want := range []string{"What I proved", "What was left out"} {
		if !strings.Contains(out, want) {
			t.Errorf("the format should ask for %q; output:\n%s", want, out)
		}
	}
}

func TestTaskStatus_workingTreeLine(t *testing.T) {
	t.Run("TSRTS-B05: The working-tree line says what is not yet shared", func(t *testing.T) {})
	for _, c := range []struct {
		e    taskState
		want string
	}{
		{taskState{Branch: "impl-x"}, "Git   impl-x · uncommitted change\n"},
		{taskState{Branch: "impl-x", Clean: true, Unpushed: 2}, "Git   impl-x · 2 unpushed commit(s)\n"},
		{taskState{Branch: "impl-x", Clean: true}, "Git   impl-x · up to date with the remote\n"},
	} {
		if out := renderTaskStatus(c.e); !strings.Contains(out, c.want) {
			t.Errorf("want %q in:\n%s", c.want, out)
		}
	}
}

func TestTaskStatus_uncommittedChangeComesBeforeAnyOtherStep(t *testing.T) {
	t.Run("TSRTS-B08: Uncommitted and unpushed work come first", func(t *testing.T) {})
	// Uncommitted work is the one state in which EVERY following advice is premature.
	steps := nextStep(taskState{Card: testCard(303, "anchors:in-progress", "x"), Branch: "impl-x", Unpushed: 1})
	if len(steps) < 2 || !strings.Contains(steps[0], "uncommitted change") || !strings.Contains(steps[1], "`git push`") {
		t.Errorf("uncommitted first, unpushed second; got %q", steps)
	}
	if strings.Contains(strings.Join(steps, "\n"), "open the PR") {
		t.Errorf("with uncommitted work, opening a PR is premature; got %q", steps)
	}
}

// THE REVERSION must appear in the report, and BEFORE everything else.
//
// Measured: an agent closed the card by hand, the state lock undid it the same minute, and
// it ended the turn writing "issue closed and resolved, nothing else to do" — without
// knowing. This report is the last place where the information still changes the outcome.
func TestTaskStatus_theReversionAppearsBeforeTheRest(t *testing.T) {
	t.Run("TSRTS-B03: A reversion appears before the verdict, with the way to authorise it", func(t *testing.T) {})
	out := renderTaskStatus(taskState{
		Card:   testCard(483, "anchors:in-progress", "[DTSTD] the spec excludes the state"),
		Branch: "fix-483", Clean: true,
		Reverted: []string{"Reverted: closed by hand (by someone), and the card was reopened."},
	})

	if !strings.Contains(out, "UNDONE") {
		t.Fatal("the reversion must appear — the agent did not see it on the card")
	}
	if !strings.Contains(out, "closed by hand") {
		t.Error("the report should say WHAT was undone")
	}
	// BEFORE the PR: a reversion changes what the agent thinks it did.
	if i, j := strings.Index(out, "UNDONE"), strings.Index(out, "PR "); i < 0 || j < 0 || i > j {
		t.Error("the reversion should come BEFORE the PR verdict")
	}
	// And HOW TO AUTHORISE: without it the agent concludes the pipeline is broken.
	if !strings.Contains(out, "anchors:manual") {
		t.Error("the report should say how to authorise the deliberate move")
	}
	// pr-body links with `Refs`: the merge does not close the card.
	if strings.Contains(out, "closes") || !strings.Contains(out, "the merge moves it to ready-to-test") {
		t.Errorf("the merge moves the card to ready-to-test and does not close it:\n%s", out)
	}
}

// WITHOUT A REVERSION the section does not exist. A section that always appears — empty
// most of the time — trains the reader to skip it.
func TestTaskStatus_withoutReversionNoSectionIsInvented(t *testing.T) {
	out := renderTaskStatus(taskState{
		Card: testCard(303, "anchors:in-progress", "x"), Branch: "impl-x", Clean: true,
	})
	if strings.Contains(out, "UNDONE") {
		t.Error("without a reversion the section should not exist")
	}
}

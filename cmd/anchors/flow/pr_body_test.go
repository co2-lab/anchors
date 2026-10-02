// @anchors
//   ref: PRBDP

package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
)

// THE SYNTAX IS THE PLATFORM'S, and that is why it lives in a map — adding one is adding a
// line, and the generator does not need to know how many exist.
//
// What must NOT happen is the syntax leaking into the doctrine: Anchors is multi-language,
// and a gate that forces the PR body to be in English is not a rule of Anchors — it is a
// GitHub requirement disguised as a rule.
func TestLinkSyntax_isPerPlatformAndNeverCloses(t *testing.T) {
	t.Run("PRBDP-B02: The link syntax links and never closes", func(t *testing.T) {})
	if _, ok := linkSyntax["github"]; !ok {
		t.Fatal("github must have a declared syntax — it is the platform of the `github` mode")
	}
	// The format must hold `%s`: without it the card number does not go in, and the
	// command would print the same line for every card.
	for platform, form := range linkSyntax {
		if !strings.Contains(form, "%s") {
			t.Errorf("the syntax of %q has no place for the card number: %q", platform, form)
		}
	}

	// LINK AND DO NOT CLOSE — and this is the rule that keeps it from coming back.
	//
	// This used to generate `Closes #N`, which makes two claims at once: the pipeline read
	// "the PR finished the implementation" and GitHub read "close the issue". The pipeline
	// does not end at `ready-to-test` — `in-test`, `ready-to-release` and `production`
	// follow, which belong to the CD and Anchors does not track. The card closed with three
	// states ahead of it.
	//
	// MEASURED in the reference project: 71 cards closed at `ready-to-test` against 7 open,
	// 35 in a single day.
	//
	// The words GitHub treats as closing are at
	// docs.github.com/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue
	// and none of them may appear in what `pr-body` GENERATES.
	for platform, form := range linkSyntax {
		word := strings.ToLower(strings.Fields(form)[0])
		for _, closing := range []string{
			"close", "closes", "closed", "fix", "fixes", "fixed",
			"resolve", "resolves", "resolved",
		} {
			if word == closing {
				t.Errorf("the syntax of %q generates %q, which CLOSES the issue on merge — the card "+
					"must stay open at `ready-to-test` because the pipeline continues in the CD",
					platform, form)
			}
		}
	}
}

// `--cards` accepts the forms a person writes: with `#`, without, with spaces. Refusing
// "#44" because of the hash would be friction for no reason — it is how the card appears
// everywhere on GitHub.
func TestRequestedCards_acceptsTheFormsAPersonWrites(t *testing.T) {
	t.Run("PRBDP-B03: The requested cards accept the forms a person writes", func(t *testing.T) {})
	cfg := &config.Config{}
	for _, in := range []string{"44", "#44", " 44 ", "#44 "} {
		got := requestedCards(in, cfg)
		if len(got) != 1 || got[0] != "44" {
			t.Errorf("%q should become [44], got %v", in, got)
		}
	}
	// Several at once: the work delivers the card AND the findings born under it.
	if got := requestedCards("44, #49,50", cfg); strings.Join(got, ",") != "44,49,50" {
		t.Errorf("three cards should become three entries, got %v", got)
	}
	// Blank does not invent a card: without `--cards` and without an agent, the caller gets
	// an error instead of a PR that links nothing.
	if got := requestedCards("  ", cfg); len(got) != 0 {
		t.Errorf("a blank input must not invent a card, got %v", got)
	}
}

// Each card drags what was born under it, and the lines come out in numeric order — so
// #101 does not sort before #45 as text would. A finding that is also a root is linked once.
func TestPRBodyCmd_linksTheCardsAndWhatWasBornUnderThem(t *testing.T) {
	t.Run("PRBDP-B06: Each root drags the open findings born under it", func(t *testing.T) {})
	t.Run("PRBDP-B07: The lines come out in numeric order", func(t *testing.T) {})
	t.Run("PRBDP-I01: A card that is both a root and a finding is linked once", func(t *testing.T) {})
	t.Run("PRBDP-X01: pr-body writes nothing to the platform", func(t *testing.T) {})
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue list *--label " + initx.LabelSob("44") + " *", out: `[{"number":101},{"number":45},{"number":50}]`},
		ghRule{match: "issue list *--label " + initx.LabelSob("50") + " *", out: `[]`},
	)
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", root, "--cards", "#44, 50"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if out != "Refs #44\nRefs #45\nRefs #50\nRefs #101\n" {
		t.Errorf("unexpected body:\n%s", out)
	}
	c := calls()
	if len(callsWith(c, "--repo acme/app", "--state open")) != 2 {
		t.Errorf("the open findings are looked up in the configured repo: %v", c)
	}
	for _, call := range c {
		if !strings.HasPrefix(call, "issue list ") {
			t.Errorf("pr-body only reads the board, but called: %s", call)
		}
	}
}

// `--so-sob` prints only what was born under the given cards — what CI checks is present.
func TestPRBodyCmd_onlyUnderLeavesTheRootsOut(t *testing.T) {
	t.Run("PRBDP-B08: The only-under switch leaves the roots out", func(t *testing.T) {})
	root := githubProject(t)
	scriptedGH(t, ghRule{match: "issue list *", out: `[{"number":45}]`})
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", root, "--cards", "44", "--so-sob"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if out != "Refs #45\n" {
		t.Errorf("only the finding under #44 is printed, got:\n%s", out)
	}
}

// Without --cards, the cards are the agent's own on the board.
func TestPRBodyCmd_discoversTheAgentsCards(t *testing.T) {
	t.Run("PRBDP-B04: Without requested cards the agent's own card is linked", func(t *testing.T) {})
	root := githubProject(t)
	scriptedGH(t,
		ghRule{match: "issue list *--json number,title,labels,comments*", out: "12\tmine\tanchors:in-progress"},
		ghRule{match: "issue list *", out: `[]`},
	)
	t.Setenv("ANCHORS_AGENT", "host/dev1")
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", root})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if out != "Refs #12\n" {
		t.Errorf("the agent's card must be linked, got:\n%s", out)
	}
}

func TestPRBodyCmd_refusals(t *testing.T) {
	t.Run("PRBDP-B05: pr-body with no card from either source is refused", func(t *testing.T) {})
	t.Run("PRBDP-B01: pr-body in local mode is refused", func(t *testing.T) {})
	scriptedGH(t)
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", githubProject(t)})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "no card") ||
		!strings.Contains(err.Error(), "--cards") || !strings.Contains(err.Error(), "ANCHORS_AGENT") {
		t.Errorf("without cards nothing can be linked, got %v", err)
	}

	cmd = newPRBodyCmd()
	cmd.SetArgs([]string{"--root", localProject(t), "--cards", "4"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode must be refused, got %v", err)
	}
}

func TestCardsUnder_failureYieldsNothing(t *testing.T) {
	t.Run("PRBDP-E01: A failed findings lookup contributes nothing and the root is still linked", func(t *testing.T) {})
	for _, r := range []ghRule{{match: "*", code: 1}, {match: "*", out: "not json"}} {
		scriptedGH(t, r)
		if got := cardsUnder(cfgGitHub(), "44"); got != nil {
			t.Errorf("a failed lookup yields nothing, got %v", got)
		}
		cmd := newPRBodyCmd()
		cmd.SetArgs([]string{"--root", githubProject(t), "--cards", "44"})
		var err error
		out := stdoutOf(t, func() { err = cmd.Execute() })
		if err != nil || out != "Refs #44\n" {
			t.Errorf("the root is still linked when its findings cannot be read, got %q, %v", out, err)
		}
	}
}

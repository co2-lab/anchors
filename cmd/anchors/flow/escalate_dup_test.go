// @anchors
//   ref: ESDPS

package flow

import (
	"strings"
	"testing"
)

// ESCALATING A TARGET THAT ALREADY HAS A CARD must warn.
//
// Two agents delivered the SAME work on the same day in the reference project: one took the
// gate's card for `MetricCard.spec.md` at 11:38; the other, working on another card, found
// the same problem at 12:02 and opened a NEW card. Both PRs added the same section to the
// same document.
func TestOpenCardsAbout_warnsWhenTheTargetAlreadyHasAnOpenCard(t *testing.T) {
	t.Run("ESDPS-B02: The board is searched for open cards with the label and the target", func(t *testing.T) {})
	t.Run("ESDPS-B03: A hit counts when the exact target is in its title or its body", func(t *testing.T) {})
	t.Run("ESDPS-B04: Each card found carries its number and its title", func(t *testing.T) {})
	t.Run("ESDPS-X01: The lookup only reads the board", func(t *testing.T) {})
	target := "apps/mobile/src/components/MetricCard.spec.md"
	calls := scriptedGH(t, ghRule{match: "issue list *", out: `[` +
		`{"number":405,"title":"[doc-required] Violation @ ` + target + `","body":"body"},` +
		`{"number":406,"title":"a human finding","body":"the section is missing in ` + target + `"}]`})
	got := openCardsAbout(target, "anchors", "acme/app")
	// The CONTENT of the warning, not just the count. Whoever reads it needs to know WHICH
	// card to look at — a warning that says "there is a card" without saying which sends
	// them searching the whole queue, and that is when they give up and create the new
	// card anyway.
	want := []string{"#405 [doc-required] Violation @ " + target, "#406 a human finding"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("the open cards of the target:\n got  %v\n want %v", got, want)
	}
	c := calls()
	if len(c) != 1 || len(callsWith(c, "issue list", "--repo acme/app", "--state open", "--label anchors", "--search "+target)) != 1 {
		t.Errorf("one search of open cards with the label and the target; calls: %v", c)
	}
}

// GITHUB'S SEARCH IS APPROXIMATE, and without confirmation it matches the wrong target:
// `MetricCard.spec.md` would match `MetricCardList.spec.md` in a text search, and the warning
// would point at work that is not the same — teaching people to ignore the warning.
func TestOpenCardsAbout_doesNotConfuseTheTargetWithOneThatContainsIt(t *testing.T) {
	t.Run("ESDPS-I01: A card about a file that merely contains the target's name is not reported", func(t *testing.T) {})
	scriptedGH(t, ghRule{match: "issue list *", out: `[{"number":999,"title":"[doc-required] Violation @ apps/mobile/src/components/MetricCardList.spec.md","body":"another target"}]`})
	if v := openCardsAbout("apps/mobile/src/components/MetricCard.spec.md", "anchors", "acme/app"); len(v) != 0 {
		t.Errorf("it matched a DIFFERENT target that only contains the name: %v", v)
	}
}

// WITHOUT A TARGET there is nothing to ask — and calling `gh` for nothing would slow every
// `escalate` without `--about`.
func TestOpenCardsAbout_withoutTargetAsksNothing(t *testing.T) {
	t.Run("ESDPS-B01: Without a target or a label the board is not asked", func(t *testing.T) {})
	calls := scriptedGH(t, ghRule{match: "*", out: `[{"number":1,"title":"X.spec.md","body":"X.spec.md"}]`})
	if v := openCardsAbout("", "anchors", "acme/app"); len(v) != 0 {
		t.Errorf("it asked without a target: %v", v)
	}
	if v := openCardsAbout("X.spec.md", "", "acme/app"); len(v) != 0 {
		t.Errorf("it asked without a label: %v", v)
	}
	if c := calls(); len(c) != 0 {
		t.Errorf("no call may reach the board: %v", c)
	}
}

// A `gh` THAT FAILS must not take `escalate` down: the check is auxiliary, and blocking the
// record because of it would be worse than the duplicate it prevents.
func TestOpenCardsAbout_failedLookupDoesNotStopTheEscalate(t *testing.T) {
	t.Run("ESDPS-E01: A failed or unreadable lookup yields nothing", func(t *testing.T) {})
	for _, r := range []ghRule{{match: "*", code: 1}, {match: "*", out: "not json"}} {
		scriptedGH(t, r)
		if v := openCardsAbout("X.spec.md", "anchors", "acme/app"); v != nil {
			t.Errorf("it returned %v when the lookup failed — it should go on without a warning", v)
		}
	}
}

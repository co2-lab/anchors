package scan

import (
	"strings"
	"testing"
)

func TestProgress_recognisesTheSuffix(t *testing.T) {
	t.Run("PRFLP-B01: Only the progress suffix marks a progress file", func(t *testing.T) {})
	cases := map[string]bool{
		"plans/0017-mutacao-progress.md": true,
		"plans/0017-mutacao.md":          false,
		"0001-progress.md":               true,
		// Not the suffix: `progress` in the middle of the name does not make a state file.
		"docs/progress-notes.md": false,
		"plans/progress.md":      false,
	}
	for path, want := range cases {
		if got := IsProgressFile(path); got != want {
			t.Errorf("%q: %v, want %v", path, got, want)
		}
	}
}

func TestProgressPathFor_replacesTheLastSegmentsExtension(t *testing.T) {
	t.Run("PRFLP-B02: The companion path replaces the extension of the last segment", func(t *testing.T) {})
	for plan, want := range map[string]string{
		"plans/x.md": "plans/x-progress.md",
		"a.b/plan":   "a.b/plan-progress.md", // the dot is in the directory, not an extension
		"x":          "x-progress.md",
	} {
		if got := ProgressPathFor(plan); got != want {
			t.Errorf("ProgressPathFor(%q) = %q, want %q", plan, got, want)
		}
	}
}

// THE REAL CASE, reproduced from the three conflicts measured in blue-eyes (#217, #221 and
// the next one): two branches tick NEIGHBOURING checkboxes of the same plan, git cannot
// resolve it, and the manual resolution is always the same.
func TestMergeProgress_joinsBothSides(t *testing.T) {
	t.Run("PRFLP-B03: Two sides ticking neighbouring items keep both ticks", func(t *testing.T) {})
	ours := "- [x] `a.spec.md` — the first\n- [ ] `b.spec.md` — the second\n"
	theirs := "- [ ] `a.spec.md` — the first\n- [x] `b.spec.md` — the second\n"

	got := MergeProgress(ours, theirs)

	if n := ProgressDone(got); n != 2 {
		t.Errorf("the union lost a tick: %d done, want 2\n%s", n, got)
	}
}

// `[x]` BEATS `[ ]` — ticking done is a fact (the spec exists, the check stamped it, the PR
// merged). Unticking by merge would erase the fact, and `progress-honest` would start
// accusing a file that exists of not being done.
func TestMergeProgress_aTickIsNeverUndone(t *testing.T) {
	t.Run("PRFLP-I02: A merge never unticks an item", func(t *testing.T) {})
	got := MergeProgress(
		"- [x] `a.spec.md` — done\n",
		"- [ ] `a.spec.md` — done\n",
	)
	if !strings.Contains(got, "- [x] `a.spec.md`") {
		t.Errorf("the merge UNTICKED a done item:\n%s", got)
	}
}

// An item that exists on one side only cannot vanish: it is the spec seeded by one branch
// while the other ticked something else.
func TestMergeProgress_keepsAnItemFromOneSideOnly(t *testing.T) {
	t.Run("PRFLP-B04: An item only on their side is added after our last item", func(t *testing.T) {})
	got := MergeProgress(
		"- [x] `a.spec.md` — ours\n",
		"- [x] `a.spec.md` — ours\n- [x] `c.spec.md` — theirs only\n",
	)
	if !strings.Contains(got, "c.spec.md") {
		t.Errorf("the item that existed only on the other side vanished:\n%s", got)
	}
	if n := ProgressDone(got); n != 2 {
		t.Errorf("want 2 done, got %d:\n%s", n, got)
	}
	// The position: after our LAST item, not at the end of the text.
	if got := MergeProgress("# h\n- [ ] a\ntail", "- [x] b\n- [ ] a"); got != "# h\n- [ ] a\n- [x] b\ntail" {
		t.Errorf("their item should come right after our last item, got %q", got)
	}
	// With no item on our side, it goes to the end.
	if got := MergeProgress("# h\nprose", "- [x] b"); got != "# h\nprose\n- [x] b" {
		t.Errorf("with no item on our side their item goes to the end, got %q", got)
	}
}

// Everything that is NOT an item survives intact — heading, prose, sections. The file is
// read by people, and a merge that rewrites the prose is worse than the conflict it avoids.
func TestMergeProgress_keepsWhatIsNotAnItem(t *testing.T) {
	t.Run("PRFLP-B05: The lines that are not items survive the merge", func(t *testing.T) {})
	ours := "# Progress — plan 0007\n\n> a note that explains something\n\n## Phase 1\n\n- [ ] `a.spec.md` — a\n"
	theirs := "# Progress — plan 0007\n\n> a note that explains something\n\n## Phase 1\n\n- [x] `a.spec.md` — a\n"

	got := MergeProgress(ours, theirs)

	for _, want := range []string{"# Progress — plan 0007", "> a note that explains something", "## Phase 1"} {
		if !strings.Contains(got, want) {
			t.Errorf("the merge ate %q:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "- [x] `a.spec.md`") {
		t.Errorf("the other side's tick was not promoted:\n%s", got)
	}
}

// Our side's INDENTATION is kept when the tick is promoted: nested items exist in plans,
// and rewriting them without the indent would change the document's structure.
func TestMergeProgress_keepsIndentation(t *testing.T) {
	t.Run("PRFLP-B06: A promoted tick keeps our indentation", func(t *testing.T) {})
	got := MergeProgress(
		"  - [ ] `a.spec.md` — nested\n",
		"  - [x] `a.spec.md` — nested\n",
	)
	if !strings.Contains(got, "  - [x] `a.spec.md`") {
		t.Errorf("the indentation was lost in the promotion:\n%q", got)
	}
}

// Identical sides are idempotent: the driver runs on every merge, including the ones with
// no conflict at all, and cannot introduce a difference where there was none.
func TestMergeProgress_equalSidesChangeNothing(t *testing.T) {
	t.Run("PRFLP-I01: Merging a side with itself changes nothing", func(t *testing.T) {})
	s := "## Phase 1\n\n- [x] `a.spec.md` — a\n- [ ] `b.spec.md` — b\n"
	if got := MergeProgress(s, s); got != strings.TrimSuffix(s, "\n") && got != s {
		t.Errorf("equal sides produced a different output:\n%q\n%q", s, got)
	}
	if n := ProgressDone(MergeProgress(s, s)); n != 1 {
		t.Errorf("the count changed in an idempotent merge: %d", n)
	}
}

func TestProgressDone_countsBothTickLetters(t *testing.T) {
	t.Run("PRFLP-B07: The done count reads both tick letters at any indentation", func(t *testing.T) {})
	if n := ProgressDone("- [x] a\n- [X] b\n- [ ] c\n  - [x] d"); n != 3 {
		t.Errorf("want 3 done, got %d", n)
	}
}

func TestMergeProgress_differentWordingIsTwoItems(t *testing.T) {
	t.Run("PRFLP-X01: An item whose text differs between the sides is kept twice", func(t *testing.T) {})
	got := MergeProgress("- [ ] a — first wording\n", "- [x] a — second wording\n")
	if !strings.Contains(got, "- [ ] a — first wording") || !strings.Contains(got, "- [x] a — second wording") {
		t.Errorf("both wordings should be kept:\n%s", got)
	}
}

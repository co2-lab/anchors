// @anchors
//   ref: SYCMS

package ops

import (
	"bytes"
	"github.com/co2-lab/anchors/internal/i18n"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// THE SYNTHESIS CARD DOES NOT PICK A SIDE, and that is what sets it apart from "resolve the
// conflict".
//
// A content conflict is two people writing different things about the same place. Picking
// a side by automation is choosing without reading the other — and in a real measured case
// (PRs #693 and #556 of the reference project) both sides were RIGHT and about different
// things: one documented that the rule's body contradicted the decision, the other that the
// waiver marker made the gate answer INDETERMINATE. Both belonged.
//
// WHAT THE CARD ASKS FOR is what neither PR delivers alone.
func TestSynthesisCardDoesNotPickASide(t *testing.T) {
	t.Run("SYCMS-B03: The card cites both PRs and cards and asks for the best of each", func(t *testing.T) {})
	t.Run("SYCMS-X01: The card never picks a side", func(t *testing.T) {})
	body := corpoDaSintese("693", "556", "666", "198", "fix(DSHBR): the B01", "feat: the index", "x.spec.md")

	// BOTH PRs and both cards appear: without them the card cannot be traced back.
	for _, ref := range []string{"#693", "#556", "#666", "#198", "fix(DSHBR): the B01", "feat: the index", "`x.spec.md`"} {
		if !strings.Contains(body, ref) {
			t.Errorf("the card does not cite %s — traceability is lost when the PRs close", ref)
		}
	}
	// And WHAT IT ASKS FOR cannot be "pick one".
	if !strings.Contains(body, "best of each") {
		t.Error("the card does not ask for the synthesis — if it asked for a choice, the automation would already have chosen")
	}
	// The WARNING against discarding without saying why: it is how this task fails.
	if !strings.Contains(body, "without saying why") {
		t.Error("the card does not warn against silently discarding a side — both works are " +
			"closed, and whatever is lost nobody will know what it was")
	}
}

// WITH ONE SIDE ONLY the card is still born, and says what is missing.
//
// When the conflict is against the integration branch, the "other side" is work already
// merged. Finding it takes reading the file's history, and guessing would be wrong — a
// one-sided card is still better than a PR stalled with no owner.
func TestSynthesisCardWithOneSideOnly(t *testing.T) {
	t.Run("SYCMS-B04: A one-sided card says the other side is missing and where to look", func(t *testing.T) {})
	body := corpoDaSintese("693", "", "666", "", "fix(DSHBR): the B01", "", "x.spec.md")

	if strings.Contains(body, "| # |") || strings.Contains(body, "#  ") {
		t.Error("the card cites an empty PR — the second side does not exist, and the table lies")
	}
	// SAYING WHAT IS MISSING makes the card actionable: whoever takes it knows they must
	// find the other side, and where to look.
	if !strings.Contains(body, "was not identified") {
		t.Error("the card does not say the other side is missing — whoever takes it will look for a PR " +
			"that does not exist")
	}
	if !strings.Contains(body, "git log") {
		t.Error("the card does not say WHERE to look for the other side — the instruction without the path " +
			"hands the discovery to whoever was already interrupted")
	}
	// And the instruction in the singular: "read both closed PRs" would be false with one.
	if strings.Contains(body, "Read both closed PRs") {
		t.Error("the card says to read TWO PRs when only one was closed")
	}
}

// --- the command, against a fake `gh` (nothing reaches GitHub) ---

// synthGH answers like GitHub for PRs 693 (card #666) and 556 (card #198). `extra` is
// bash inserted before the defaults, to make one call fail.
func synthGH(t *testing.T, extra string) string {
	t.Helper()
	return fakeGH(t, extra+`
case "$*" in
"pr view 693"*"--json body"*) echo "Refs #666" ;;
"pr view 556"*"--json body"*) echo "Closes #198" ;;
"pr view 693"*"--json title"*) echo "fix(DSHBR): the B01" ;;
"pr view 556"*"--json title"*) echo "feat: the index" ;;
"issue create"*) echo "https://github.com/acme/app/issues/700" ;;
esac`)
}

func synthProject(t *testing.T, cfg string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, cfg)
	return root
}

func runSynth(t *testing.T, args ...string) (err error, out, errOut string) {
	t.Helper()
	c := newSynthesizeCmd()
	var o, e bytes.Buffer
	c.SetOut(&o)
	c.SetErr(&e)
	c.SetArgs(args)
	err = c.Execute()
	return err, o.String(), e.String()
}

func TestSynthesizeExistsOnlyInGithubModeAndNeedsAPR(t *testing.T) {
	t.Run("SYCMS-B01: Only github mode is accepted and the first PR is required", func(t *testing.T) {})
	log := synthGH(t, "")
	if err, _, _ := runSynth(t, "--root", synthProject(t, "version: 1\n"), "--pr-a", "693"); err == nil ||
		!strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode: %v", err)
	}
	if err, _, _ := runSynth(t, "--root", synthProject(t, githubModeConfig)); err == nil ||
		!strings.Contains(err.Error(), "--pr-a") {
		t.Errorf("without --pr-a: %v", err)
	}
	if calls := readLog(t, log); calls != "" {
		t.Errorf("a refused call still talked to GitHub:\n%s", calls)
	}
}

// --dry-run shows the card — both PRs, both cards, both titles — and touches nothing.
func TestSynthesizeDryRunShowsTheCardOnly(t *testing.T) {
	t.Run("SYCMS-B02: A leading # on a PR number is dropped", func(t *testing.T) {})
	t.Run("SYCMS-B09: A dry run shows the card and changes nothing", func(t *testing.T) {})
	log := synthGH(t, "")
	err, out, _ := runSynth(t, "--root", synthProject(t, githubModeConfig),
		"--pr-a", "#693", "--pr-b", "556", "--files", "Dashboards.spec.md", "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"[synthesis] what PRs #693 and #556 deliver, in one",
		"#666", "#198", "fix(DSHBR): the B01", "feat: the index", "Dashboards.spec.md"} {
		if !strings.Contains(out, want) {
			t.Errorf("the dry run misses %q:\n%s", want, out)
		}
	}
	calls := readLog(t, log)
	for _, forbidden := range []string{"issue create", "pr close", "pr comment", "label create"} {
		if strings.Contains(calls, forbidden) {
			t.Errorf("--dry-run ran %q:\n%s", forbidden, calls)
		}
	}
}

// The real run links the five ends: the card (labelled under both origin cards), both
// PRs commented and closed, both origin cards pointed at the new one.
func TestSynthesizeLinksTheFiveEnds(t *testing.T) {
	t.Run("SYCMS-B05: The card is labelled for the board and under each origin card", func(t *testing.T) {})
	t.Run("SYCMS-B06: Each PR is commented and closed and each origin card is pointed at the new card", func(t *testing.T) {})
	log := synthGH(t, "")
	err, out, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig),
		"--pr-a", "693", "--pr-b", "556")
	if err != nil {
		t.Fatalf("synthesize: %v\n%s", err, errOut)
	}
	calls := readLog(t, log)
	for _, want := range []string{
		"label create anchors:under-666 --repo acme/app",
		"label create anchors:under-198 --repo acme/app",
		"--label anchors --label anchors:to-do --label anchors:under-666 --label anchors:under-198",
		"pr comment 693 --repo acme/app", "pr close 693 --repo acme/app",
		"pr comment 556 --repo acme/app", "pr close 556 --repo acme/app",
		"issue comment 666 --repo acme/app", "issue comment 198 --repo acme/app",
	} {
		if !strings.Contains(calls, want) {
			t.Errorf("missing call %q:\n%s", want, calls)
		}
	}
	if !strings.Contains(calls, "https://github.com/acme/app/issues/700") {
		t.Errorf("the comments do not point at the new card:\n%s", calls)
	}
	if !strings.Contains(out, "synthesis card: https://github.com/acme/app/issues/700") ||
		!strings.Contains(out, "point at each other") {
		t.Errorf("output:\n%s", out)
	}
	if errOut != "" {
		t.Errorf("no warning was expected:\n%s", errOut)
	}
}

// With one side only, the other end is the integration branch, and only that PR closes.
// Every failed link is WARNED — the card already exists.
func TestSynthesizeOneSideWarnsOnEachFailedLink(t *testing.T) {
	t.Run("SYCMS-B06: Each PR is commented and closed and each origin card is pointed at the new card", func(t *testing.T) {})
	t.Run("SYCMS-B07: Each failed link is a warning and the command succeeds", func(t *testing.T) {})
	log := synthGH(t, `case "$*" in "pr comment"*|"pr close"*|"issue comment"*) echo nope >&2; exit 1 ;; esac`)
	err, _, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693")
	if err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	calls := readLog(t, log)
	if !strings.Contains(calls, "what PR #693 delivers, reconciled with what already landed") {
		t.Errorf("the one-sided title is missing:\n%s", calls)
	}
	if !strings.Contains(calls, "what is already on the integration branch") {
		t.Errorf("the PR comment does not name the other side:\n%s", calls)
	}
	if strings.Contains(calls, "pr close 556") || strings.Contains(calls, "under-198") {
		t.Errorf("a side that does not exist was touched:\n%s", calls)
	}
	for _, want := range []string{"#693 did not receive the comment", "#693 was not closed", "card #666 did not receive the link"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("no warning %q:\n%s", want, errOut)
		}
	}
}

func TestSynthesizeFailsWhenTheCardCannotBeOpened(t *testing.T) {
	t.Run("SYCMS-E01: A card that cannot be opened fails with the host's message", func(t *testing.T) {})
	t.Run("SYCMS-I01: No PR is closed unless the synthesis card exists", func(t *testing.T) {})
	log := synthGH(t, `case "$*" in "issue create"*) echo "HTTP 403" >&2; exit 1 ;; esac`)
	err, _, _ := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693", "--pr-b", "556")
	if err == nil || !strings.Contains(err.Error(), "open the synthesis card") || !strings.Contains(err.Error(), "HTTP 403") {
		t.Fatalf("want the card error with gh's message, got %v", err)
	}
	if strings.Contains(readLog(t, log), "pr close") {
		t.Error("the PRs were closed although the synthesis card does not exist")
	}
}

// With one side only, the closing line tells what happened: ONE PR was closed. It used to
// print "PRs #693 and # were closed, and the four ends point at each other" — a second PR
// that does not exist, and four ends where there are three.
func TestSynthesizeOneSideReportsOnePRClosed(t *testing.T) {
	t.Run("SYCMS-B08: The closing line names only the PRs actually closed", func(t *testing.T) {})
	synthGH(t, "")
	err, out, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693")
	if err != nil {
		t.Fatalf("synthesize: %v\n%s", err, errOut)
	}
	if strings.Contains(out, "and #") || strings.Contains(out, "four ends") {
		t.Errorf("the one-sided run reports a second PR:\n%s", out)
	}
	if !strings.Contains(out, "PR #693 was closed") {
		t.Errorf("the one-sided run does not say PR #693 was closed:\n%s", out)
	}
}

// The closing line names only what gh actually closed: a PR whose close failed is
// reported as still open, never as closed.
func TestSynthesizeReportsAFailedCloseAsStillOpen(t *testing.T) {
	t.Run("SYCMS-B08: The closing line names only the PRs actually closed", func(t *testing.T) {})
	synthGH(t, `case "$*" in "pr close 556"*) echo "GraphQL: cannot close" >&2; exit 1 ;; esac`)
	err, out, errOut := runSynth(t, "--root", synthProject(t, githubModeConfig), "--pr-a", "693", "--pr-b", "556")
	if err != nil {
		t.Fatalf("synthesize: %v\n%s", err, errOut)
	}
	if strings.Contains(out, "#556 were closed") || strings.Contains(out, "PRs #693, #556") {
		t.Errorf("a PR whose close failed was reported closed:\n%s", out)
	}
	if !strings.Contains(out, "PR #693 was closed") || !strings.Contains(out, "still OPEN") || !strings.Contains(out, "#556") {
		t.Errorf("the line must say #693 closed and #556 still open:\n%s", out)
	}
}

// The synthesis card speaks the project's language: title, body and comments come from
// the message catalog, not from English written into the code.
func TestSynthesisCardFollowsTheLanguage(t *testing.T) {
	t.Run("SYCMS-B10: The card follows the project language", func(t *testing.T) {})
	prev := i18n.Current()
	if err := i18n.Set("pt-BR"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = i18n.Set(prev) }()
	body := corpoDaSintese("693", "556", "666", "198", "fix", "feat", "src/a.ts")
	for _, want := range []string{"Dois PRs escreveram", "## Como entregar", "`Refs #693` e `Refs #556`", "- `src/a.ts`"} {
		if !strings.Contains(body, want) {
			t.Errorf("the pt-BR card must carry %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "How to deliver") {
		t.Errorf("English leaked into the pt-BR card:\n%s", body)
	}
}

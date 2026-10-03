// @anchors
//   ref: GVGDG

package governance_test

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/cmd/anchors/flow"
	"github.com/co2-lab/anchors/cmd/anchors/governance"
	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/cmd/anchors/ops"
	"github.com/co2-lab/anchors/cmd/anchors/quality"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/settings"
	"github.com/spf13/cobra"
)

// These tests live in the external package on purpose: the ops package imports this one,
// so only from outside can the whole command tree be built, and the guides are read the
// way a user reads them — through `anchors guide <name>`, not through the constants.

// guideTitles is the set of guide subcommands and the first line each one prints.
var guideTitles = map[string]string{
	"code":       "# Code guide",
	"feature":    "# Feature guide",
	"flag":       "# Feature flag guide",
	"flow":       "# Flow guide",
	"guide":      "# Guide guide",
	"header":     "# Header guide",
	"changelog":  "# Changelog guide",
	"plan":       "# Plan guide",
	"product":    "# Product doctrine guide",
	"project":    "# Project guide",
	"review":     "# Review guide",
	"spec":       "# Spec guide",
	"test":       "# Test guide",
	"work":       "# Work guide",
	"report-bug": "# Report-bug guide",
}

// captureOut collects what fn writes to os.Stdout: the guides print with fmt.Print.
func captureOut(t *testing.T, fn func()) string {
	t.Helper()
	return testkit.CaptureStdout(t, fn)
}

// governanceRoot is a root holding only what this package registers.
func governanceRoot() *cobra.Command {
	root := &cobra.Command{Use: "anchors", SilenceUsage: true, SilenceErrors: true}
	governance.Register(root)
	return root
}

// guideOut runs `anchors guide <args>` and returns what it printed.
func guideOut(t *testing.T, args ...string) string {
	t.Helper()
	root := governanceRoot()
	root.SetArgs(append([]string{"guide"}, args...))
	var err error
	out := captureOut(t, func() { err = root.Execute() })
	if err != nil {
		t.Fatalf("anchors guide %v: %v", args, err)
	}
	return out
}

// guideIn prints a guide the way a project with no declaration would see it.
func guideIn(t *testing.T, name string) string {
	t.Helper()
	if name == "review" || name == "work" {
		return guideOut(t, name, "--root", t.TempDir())
	}
	return guideOut(t, name)
}

func flat(s string) string { return strings.Join(strings.Fields(s), " ") }

func TestBareGuidePrintsThePlaybook(t *testing.T) {
	t.Run("GVGDG-B01: Bare guide prints the operating playbook", func(t *testing.T) {})
	out := guideOut(t)
	if !strings.HasPrefix(out, "# Operating Anchors (guide for AI agents)\n") {
		t.Errorf("bare `guide` should print the playbook:\n%.200s", out)
	}
	if !strings.Contains(out, "## The development flow") || !strings.Contains(out, "## Command reference") {
		t.Errorf("the playbook should carry the flow and the command reference:\n%.200s", out)
	}
}

func TestEachGuideSubcommandPrintsItsGuide(t *testing.T) {
	t.Run("GVGDG-B02: Each guide subcommand prints its own guide", func(t *testing.T) {})
	guide, _, err := governanceRoot().Find([]string{"guide"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range guide.Commands() {
		got = append(got, c.Name())
		if c.Short == "" {
			t.Errorf("`guide %s` needs a short description: it is what --help shows", c.Name())
		}
	}
	var want []string
	for name := range guideTitles {
		want = append(want, name)
	}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("guide subcommands = %v, want %v", got, want)
	}
	seen := map[string]string{}
	for name, title := range guideTitles {
		out := guideIn(t, name)
		if !strings.HasPrefix(out, title) {
			t.Errorf("`guide %s` should start with %q, got:\n%.120s", name, title, out)
		}
		if other, dup := seen[out]; dup {
			t.Errorf("`guide %s` prints the same text as `guide %s`", name, other)
		}
		seen[out] = name
	}
}

// The Long help's hand-written list named 8 of the 13 subcommands.
func TestGuideHelpListsEverySubcommand(t *testing.T) {
	t.Run("GVGDG-B14: guide --help lists every subcommand", func(t *testing.T) {})
	guide, _, err := governanceRoot().Find([]string{"guide"})
	if err != nil {
		t.Fatal(err)
	}
	listed := map[string]int{}
	for _, l := range strings.Split(guide.Long, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(l), "anchors guide "); ok {
			listed[strings.Fields(rest)[0]]++
		}
	}
	for _, c := range guide.Commands() {
		if listed[c.Name()] != 1 {
			t.Errorf("`anchors guide --help` must list `guide %s` once, listed %d time(s)", c.Name(), listed[c.Name()])
		}
		delete(listed, c.Name())
	}
	for name := range listed {
		t.Errorf("`anchors guide --help` lists `guide %s`, which is not a subcommand", name)
	}
}

func TestReviewAndWorkAppendTheAutonomySection(t *testing.T) {
	t.Run("GVGDG-B03: The review and work guides append the autonomy section of the root they are given", func(t *testing.T) {})
	dev := t.TempDir()
	if err := settings.Save(dev, settings.Settings{Role: "dev", DecidedAt: "2026-09-26"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"review", "work"} {
		out := guideOut(t, name, "--root", dev)
		if !strings.HasPrefix(out, guideTitles[name]) {
			t.Errorf("`guide %s` should print its guide first:\n%.120s", name, out)
		}
		if !strings.Contains(out, "## When you do not know") || !strings.Contains(out, "**Your role (Dev) does NOT decide") {
			t.Errorf("`guide %s --root` should append the section for the role declared there:\n%s", name, out)
		}
		if strings.Index(out, "## When you do not know") < len(guideTitles[name]) {
			t.Errorf("the autonomy section comes after the guide in `guide %s`", name)
		}
		// a root with no project at all is not a failure: the section is read from a
		// directory with no declaration, and says so
		out = guideOut(t, name, "--root", filepath.Join(t.TempDir(), "nowhere"))
		if !strings.Contains(out, "**You did not declare a role**") {
			t.Errorf("`guide %s` outside a project should still print, with no role:\n%s", name, out)
		}
	}
}

func TestOnlyReviewAndWorkDependOnTheProject(t *testing.T) {
	t.Run("GVGDG-X01: Only the review and work guides depend on the project", func(t *testing.T) {})
	guide, _, err := governanceRoot().Find([]string{"guide"})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range guide.Commands() {
		hasRoot := c.Flags().Lookup("root") != nil
		wants := c.Name() == "review" || c.Name() == "work"
		if hasRoot != wants {
			t.Errorf("`guide %s` takes --root = %v, want %v", c.Name(), hasRoot, wants)
		}
		if wants {
			continue
		}
		if out := guideOut(t, c.Name()); strings.Contains(out, "## When you do not know") {
			t.Errorf("`guide %s` should print its text alone, without the autonomy section", c.Name())
		}
	}
}

func TestRegisterAddsTheGovernanceCommands(t *testing.T) {
	t.Run("GVGDG-B04: Register adds exactly the four governance commands", func(t *testing.T) {})
	var got []string
	for _, c := range governanceRoot().Commands() {
		got = append(got, c.Name())
	}
	slices.Sort(got)
	if want := []string{"audit", "compliance", "governs", "guide"}; !slices.Equal(got, want) {
		t.Errorf("Register added %v, want %v", got, want)
	}
}

// fullTree is every command name the real root registers, at any depth.
func fullTree() map[string]bool {
	root := &cobra.Command{Use: "anchors"}
	governance.Register(root)
	mapcmd.Register(root)
	quality.Register(root)
	flow.Register(root)
	ops.Register(root)
	names := map[string]bool{}
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		names[c.Name()] = true
		for _, s := range c.Commands() {
			walk(s)
		}
	}
	walk(root)
	return names
}

// A guide that sends the reader to a command that does not exist leaves them with
// `unknown command` and nowhere to go. "the anchors <word>" is prose about the documents,
// not a command.
func TestEveryCommandAGuideCitesExists(t *testing.T) {
	t.Run("GVGDG-I01: Every command a guide cites exists in the command tree", func(t *testing.T) {})
	registered := fullTree()
	re := regexp.MustCompile(`anchors ([a-z][a-z-]*)`)
	texts := map[string]string{"(playbook)": guideOut(t)}
	for name := range guideTitles {
		texts[name] = guideIn(t, name)
	}
	cited := 0
	for name, text := range texts {
		for _, m := range re.FindAllStringSubmatchIndex(text, -1) {
			if strings.HasSuffix(text[:m[0]], "the ") {
				continue
			}
			cited++
			if cmd := text[m[2]:m[3]]; !registered[cmd] {
				t.Errorf("guide %s tells the reader to run `anchors %s`, and no such command exists", name, cmd)
			}
		}
	}
	if cited == 0 {
		t.Fatal("no `anchors <command>` found in any guide — the pattern broke and the test would pass empty")
	}
}

// The label name in the guide has to be the REAL one. The guide once said
// `anchors:sob-<n>`, the name before the migration to English, and whoever followed it
// looked on the board for a label that does not exist.
func TestWorkGuideUsesTheRealLabelNames(t *testing.T) {
	t.Run("GVGDG-I02: The work guide uses the real label names", func(t *testing.T) {})
	work := guideIn(t, "work")
	if strings.Contains(work, "anchors:sob-") {
		t.Error("the work guide uses the OLD label name")
	}
	for _, label := range []string{initx.PrefixoLabelSob, initx.LabelNeedsUser} {
		if !strings.Contains(work, label) {
			t.Errorf("the work guide should use the real label %q", label)
		}
	}
}

// THE VERDICT IS A LINE THE PIPELINE READS (reference app). The closed cycle: each
// verdict line the guide shows, with a name in place of <you>, is parsed by the very
// expression `anchors-pr-checks.yml` runs — a guide teaching another spelling would leave
// every reviewed PR pending forever.
func TestReviewVerdictLineIsWhatThePipelineParses(t *testing.T) {
	t.Run("GVGDG-I03: The verdict line the review guide teaches is the one the pipeline parses", func(t *testing.T) {})
	wf, err := os.ReadFile(filepath.Join("..", "..", "..", "internal", "initx", "workflows", "anchors-pr-checks.yml"))
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`sed -nE 's/(\^anchors-review: [^/]*)/`).FindSubmatch(wf)
	if m == nil {
		t.Fatal("the pipeline's verdict expression was not found in anchors-pr-checks.yml")
	}
	pipeline := regexp.MustCompile(string(m[1]))
	review := guideIn(t, "review")
	for _, verdict := range []string{"approved", "rejected"} {
		taught := "    anchors-review: " + verdict + " by <you>\n"
		if !strings.Contains(review, taught) {
			t.Errorf("the review guide should show %q as a copyable block", strings.TrimSpace(taught))
			continue
		}
		line := strings.Replace(strings.TrimSpace(taught), "<you>", "agent-b", 1)
		got := pipeline.FindStringSubmatch(line)
		if got == nil || got[1] != verdict || got[2] != "agent-b" {
			t.Errorf("the pipeline does not read the taught line %q as %s by agent-b (got %v)", line, verdict, got)
		}
	}
}

// WHO counts: the assigned reviewer, after the assignment — otherwise the reviewer of an
// earlier round, or any agent, could sign.
func TestReviewGuideTeachesWhoCountsAndWhichVerdictWins(t *testing.T) {
	t.Run("GVGDG-B05: The review guide teaches who counts and which verdict wins", func(t *testing.T) {})
	review := flat(guideIn(t, "review"))
	for _, want := range []string{
		"Only the reviewer the claim assigned counts",
		"only a line posted after the assignment",
		"the last one wins",
		"Inside a code block the line is an example, not a verdict",
	} {
		if !strings.Contains(review, want) {
			t.Errorf("the review guide should say %q", want)
		}
	}
}

// The review guide said how to judge and never said WHO MOVES THE CARD. Measured: two
// independent reviews ran in parallel over the same PR. The first approved and moved the
// card to `ready-to-test` by hand, with the PR open. The second found a real defect and
// did NOT touch the state — the first had created the fact that blocked it.
func TestReviewGuideSaysTheReviewerDoesNotMoveTheCard(t *testing.T) {
	t.Run("GVGDG-B06: The review guide says the reviewer does not move the card", func(t *testing.T) {})
	review := guideIn(t, "review")
	if !strings.Contains(review, "YOU DO NOT MOVE THE CARD") {
		t.Error("the review guide should say explicitly that the reviewer does not move the card")
	}
	// the FACT that moves each state — without it "do not move" is a ban with no reason
	if !strings.Contains(review, "MERGED") || !strings.Contains(flat(review), "when the checks") {
		t.Error("the review guide should say the checks move to `ready-to-review` and the MERGE to `ready-to-test`")
	}
	// the consequence, and which of the two errors is worse
	if !strings.Contains(review, "end of Anchors' jurisdiction") || !strings.Contains(review, "more expensive than the late state") {
		t.Error("the review guide should say why the wrong state costs more than the late one")
	}
}

// The list exists so that nothing is missed by FORGETTING. Each point has a code with no
// gap in the numbering (a gap is a point deleted without anyone noticing), and each one
// distils a ruler the prose above already explains.
func TestReviewGuideCarriesAnchoredConformancePoints(t *testing.T) {
	t.Run("GVGDG-B07: The review guide carries continuous conformance points anchored in its prose", func(t *testing.T) {})
	review := guideIn(t, "review")
	const heading = "## Pontos de conformidade"
	at := strings.Index(review, heading)
	if at < 0 {
		t.Fatal("the review guide has no conformance points section")
	}
	// the heading must be one the `guide-checklist` gate recognises
	if !slices.Contains(i18n.AllTranslations("section.title.compliance_points"), strings.TrimPrefix(heading, "## ")) {
		t.Errorf("%q is not a heading the guide-checklist gate recognises", heading)
	}
	for i := 1; i <= 18; i++ {
		if code := "REV-CK" + itoa(i) + ":"; !strings.Contains(review, code) {
			t.Errorf("point %s is missing — the numbering has a gap", code)
		}
	}
	if strings.Contains(review, "REV-CK19:") {
		t.Error("a 19th point appeared: extend this test and its anchors together")
	}
	body := review[:at]
	anchors := map[string]string{
		"REV-CK1":  "do the checks EXIST?",
		"REV-CK3":  "decide what it needed to decide",
		"REV-CK4":  "realize the rule, or only cite it",
		"REV-CK5":  "AWAITING JUDGMENT?",
		"REV-CK6":  "@TBD",
		"REV-CK7":  "contradict each other",
		"REV-CK10": "Checked:",
		"REV-CK11": "PROVE, or only execute",
		"REV-CK12": "WHICH requirement it proves",
		"REV-CK13": "change without saying",
		"REV-CK14": "DO NOT MOVE THE CARD",
		"REV-CK15": "YOUR VERDICT IS A LINE ON THE PR",
		"REV-CK16": "walk its variations, or stop at the first draft?",
		"REV-CK17": "state the INTENT, or the mechanism?",
		"REV-CK18": "FIX carries its rule",
	}
	for ck, anchor := range anchors {
		if !strings.Contains(body, anchor) {
			t.Errorf("%s has no prose above the list (looked for %q)", ck, anchor)
		}
	}
	// the list does NOT replace the checks
	if !strings.Contains(flat(review), "NOT a substitute for the checks") {
		t.Error("the list should say it is not a substitute for the checks")
	}
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

// THE GUIDE HAS TO TEACH THE CLAIM — it is the first command. Measured in the reference
// project: 38 of 44 PRs were opened with the card still in `to-do`, and the review queue
// showed ONE item while 46 pieces of work waited.
func TestWorkGuideTeachesTheClaimFirst(t *testing.T) {
	t.Run("GVGDG-B08: The work guide teaches the claim as the first step", func(t *testing.T) {})
	work := guideIn(t, "work")
	claim := strings.Index(work, "anchors next")
	order := strings.Index(work, "## The order")
	if claim < 0 || order < 0 {
		t.Fatalf("the work guide should teach `anchors next` and the board order (next at %d, order at %d)", claim, order)
	}
	if claim > order {
		t.Error("`anchors next` is taught after the board order — it is the FIRST command")
	}
	if !strings.Contains(work, "ANCHORS_SESSION") {
		t.Error("the work guide should say to declare ANCHORS_SESSION before claiming")
	}
	// and WHY: without the reason the claim becomes ceremony, and ceremony is skipped
	if !strings.Contains(work, "find it empty") {
		t.Error("the work guide should say why the claim matters (the review queue left empty)")
	}
}

// PUSHING IS NOT DELIVERING. An agent diagnosed correctly why the CI failed, rebuilt,
// pushed — and ended the turn with "awaiting the new round". The card stayed
// `in-progress` with its name on it.
func TestWorkGuideMakesTheAgentWaitForTheCI(t *testing.T) {
	t.Run("GVGDG-B09: The work guide makes the agent wait for the CI verdict", func(t *testing.T) {})
	work := guideIn(t, "work")
	for _, want := range []string{
		"--watch",       // the command that waits for the agent
		"work review",   // what is missing in the PR that is yours
		"escalate",      // the only legitimate exit that is not green
		"in-progress",   // the cost of stopping halfway
		"half the work", // the correct diagnosis is half the work
		"is not a stopping point",
		"does not close the card",
	} {
		if !strings.Contains(work, want) {
			t.Errorf("the work guide should cover waiting for the CI (expected %q)", want)
		}
	}
}

// WHO CLOSES THE CARD IS THE MERGE. Measured: a project accumulated 25 green PRs awaiting
// review with ZERO cards in `ready-to-review` — they had been closed by hand a minute
// before the PR was opened.
func TestWorkGuideForbidsClosingTheCardByHand(t *testing.T) {
	t.Run("GVGDG-B10: The work guide forbids closing the card by hand", func(t *testing.T) {})
	work := guideIn(t, "work")
	for _, want := range []string{"YOU DO NOT CLOSE THE CARD", "Who closes it is the MERGE", "review queue empties", "NEW work", "25 green PRs"} {
		if !strings.Contains(work, want) {
			t.Errorf("the work guide should say %q", want)
		}
	}
}

// A finding that is not the card's own is recorded, tied to the card, and — before a new
// decision is opened — checked against the decision queue that already exists (measured:
// 23 open escalations were seven decisions).
func TestWorkGuideRoutesAFindingThroughEscalate(t *testing.T) {
	t.Run("GVGDG-B11: The work guide sends a finding through escalate and reads the decision queue first", func(t *testing.T) {})
	work := guideIn(t, "work")
	for _, want := range []string{"anchors escalate", "--card", "--for-user", "anchors judge --pending", "@TBD", "gh issue list"} {
		if !strings.Contains(work, want) {
			t.Errorf("the work guide should cover %q", want)
		}
	}
	query := strings.Index(work, "gh issue list --label "+initx.LabelNeedsUser)
	escalate := strings.Index(work, "--for-user")
	if query < 0 || escalate < 0 || query > escalate {
		t.Errorf("the decision queue must be read before `--for-user` is taught (query at %d, --for-user at %d)", query, escalate)
	}
}

// The discover phase runs in the conversation, with no gate to confront it later: if the
// guide loses a piece, nobody finds out by a failure. And the phase only serves if the
// playbook sends the agent there BEFORE planning.
func TestProjectGuideCoversTheDiscoverPhase(t *testing.T) {
	t.Run("GVGDG-B12: The project guide covers the discover phase and the playbook points to it before planning", func(t *testing.T) {})
	project := guideIn(t, "project")
	for _, want := range []string{
		"PROJECT.md", "INSIGHTS.md", "Discarded", "inconsistency review", "CONVERSATION",
		"ONE question at a time", "OPINIONATED", "anchors init",
		"Stage 1 — Purpose and form",
		"Stage 2 — Language and runtime",
		"Stage 3 — Architecture and paradigm",
		"Stage 4 — Macro structure and file conventions",
		"Stage 5 — Tooling and formatting",
		"indentation", "extensions", "editors", "paradigm", "co-location",
	} {
		if !strings.Contains(project, want) {
			t.Errorf("the project guide should cover %q", want)
		}
	}
	playbook := guideOut(t)
	for _, want := range []string{"anchors guide project", "PROJECT.md", "INSIGHTS.md"} {
		if !strings.Contains(playbook, want) {
			t.Errorf("the playbook should cite %q", want)
		}
	}
	discover := strings.Index(playbook, "### 0.5. DISCOVER")
	plan := strings.Index(playbook, "### 1. PLAN")
	if discover < 0 || plan < 0 || discover > plan {
		t.Errorf("DISCOVER should come before PLAN in the playbook (discover at %d, plan at %d)", discover, plan)
	}
}

// A boundary rule applies to every input; the guide says which instrument proves it for
// each shape of input space (reference app), and teaches the refresh of the stamp.
func TestTestGuideNamesInstrumentsAndTheStampRefresh(t *testing.T) {
	t.Run("GVGDG-B13: The test guide names the instrument per input shape and teaches the stamp refresh", func(t *testing.T) {})
	test := guideIn(t, "test")
	for _, want := range []string{
		"SMALL AND CLOSED", "EXHAUSTIVE", "LARGE BUT STRUCTURED", "TABLE OF CLASSES", "OPEN", "never a\n  list of the forbidden",
		"anchors stamp --refresh", "mock-stamped", "same commit",
	} {
		if !strings.Contains(test, want) {
			t.Errorf("the test guide should say %q", want)
		}
	}
}

func TestGuidesTellAFixFromABug(t *testing.T) {
	t.Run("GVGDG-B15: The guides tell a fix from a bug and ask for the marker", func(t *testing.T) {})
	work := guideIn(t, "work")
	for _, want := range []string{"## When what you deliver fixes something", "SHIPPED", "fix(<scope>):",
		"Bug: <where it was seen", "Its own commit", "Reproduce first", "GAP IN THE SPEC"} {
		if !strings.Contains(work, want) {
			t.Errorf("the work guide should say %q", want)
		}
	}
	if !strings.Contains(guideIn(t, "code"), "When what you deliver fixes something") {
		t.Error("the code guide should point to the fix section")
	}
	if !strings.Contains(guideIn(t, "review"), "Is each fix its own commit, and is each bug marked?") {
		t.Error("the review guide should ask whether each fix is its own commit and each bug is marked")
	}
}

func TestChangelogGuideTellsTechnicalFromProduct(t *testing.T) {
	t.Run("GVGDG-B16: The changelog guide tells the technical changelog from the product one", func(t *testing.T) {})
	g := guideIn(t, "changelog")
	for _, want := range []string{"TECHNICAL changelog", "## The product changelog (recommended)",
		"Write it with an agent, from the technical changelog", "**Bugs fixed**",
		"**Fixes** without", "**Chores**", "If it matters to\n  the product"} {
		if !strings.Contains(g, want) {
			t.Errorf("the changelog guide should say %q", want)
		}
	}
}

func TestGuides_theSpecIsWrittenInFourPasses(t *testing.T) {
	t.Run("GVGDG-B17: The spec guide asks for four passes and a review", func(t *testing.T) {})
	g := guideIn(t, "spec")
	for _, want := range []string{"## Writing it: every section, every variation, then a review",
		"1. WALK EVERY SECTION", "Every INPUT the unit accepts", "Every EFFECT", "Every FAILURE",
		"Every STATE the unit reads that another unit also reads",
		"2. DERIVE THE VARIATIONS", "3. GENERALIZE", "an invariant", "4. REVIEW WHAT YOU WROTE",
		"the INTENT, or the mechanism", "WHEN A DEFECT REACHES YOU"} {
		if !strings.Contains(g, want) {
			t.Errorf("the spec guide should say %q", want)
		}
	}
	if strings.Index(g, "## Writing it") > strings.Index(g, "## Spec rules") {
		t.Error("the passes come before the spec rules, right after the sections they walk")
	}
}

func TestReportBugGuide_tellsReportsAndGoesOn(t *testing.T) {
	t.Run("GVGDG-B18: The report-bug guide says how to tell, report and go on", func(t *testing.T) {})
	for _, want := range []string{"Is it Anchors, or the project?", "MINIMAL CASE, made up", "PUBLIC",
		"--dry-run", "While the fix does not come", "Do not edit Anchors' installed files", "[skip-<gate>@<CODE>: Anchors issue"} {
		if !strings.Contains(guideOut(t, "report-bug"), want) {
			t.Errorf("the report-bug guide lacks %q", want)
		}
	}
	for _, want := range []string{"anchors report-bug", "anchors guide report-bug"} {
		if !strings.Contains(guideOut(t), want) {
			t.Errorf("the playbook does not point to %q", want)
		}
	}
}

func TestGuides_recommendVisualRegressionPerState(t *testing.T) {
	t.Run("GVGDG-B19: The guides recommend visual regression for every state of a visual unit", func(t *testing.T) {})
	test := guideOut(t, "test")
	for _, want := range []string{"## Visual regression", "EVERY state", "{CODE}-VR-<state>", "CAPTURE test", "vr-states-covered"} {
		if !strings.Contains(test, want) {
			t.Errorf("the test guide lacks %q", want)
		}
	}
	if !strings.Contains(guideOut(t, "spec"), "VISUAL REGRESSION") {
		t.Error("the spec guide does not point a visual unit's states to visual regression")
	}
	if !strings.Contains(guideOut(t, "feature"), "{CODE}-VR") {
		t.Error("the feature guide does not ask for the VR scenario")
	}
}

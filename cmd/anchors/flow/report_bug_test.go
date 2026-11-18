// @anchors
//   code: RBTRP
//   ref: RPBUG

package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runReportBug runs `report-bug` in root with the given gh rules and standard input, and
// returns stdout, the error and the recorded gh calls.
func runReportBug(t *testing.T, root string, rules []ghRule, stdin string, args ...string) (string, error, []string) {
	t.Helper()
	calls := scriptedGH(t, rules...)
	cmd := newReportBugCmd()
	cmd.SilenceErrors = true
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(append([]string{"--root", root}, args...))
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	return out, err, calls()
}

func TestAnchorsBug_hasTheBugFormSections(t *testing.T) {
	t.Run("RPBUG-B01: The report has the bug form's sections", func(t *testing.T) {})
	b := anchorsBug{What: "gate misreads a file\nmore detail", Expected: "it reads it", Repro: "anchors check --all"}
	if b.title() != "[bug] gate misreads a file" {
		t.Errorf("title = %q", b.title())
	}
	body := b.body()
	for _, want := range []string{"### What happened", "more detail", "### What should happen\n\nit reads it",
		"### Minimal case\n\n```shell\nanchors check --all\n```", "### Version", "### Platform", "anchors report-bug"} {
		if !strings.Contains(body, want) {
			t.Errorf("the body lacks %q:\n%s", want, body)
		}
	}
	if bare := (anchorsBug{What: "x"}).body(); strings.Contains(bare, "What should happen") || strings.Contains(bare, "Minimal case") {
		t.Errorf("a section with nothing to say is left out:\n%s", bare)
	}
}

func TestReportBug_seenAgainInsteadOfADuplicate(t *testing.T) {
	t.Run("RPBUG-B02: An open issue with the same title is told it was seen again", func(t *testing.T) {})
	out, err, calls := runReportBug(t, t.TempDir(),
		[]ghRule{{match: "issue list --repo co2-lab/anchors *", out: "[BUG] Gate misreads a file\t" + upstreamURL}}, "",
		"gate misreads a file", "--expected", "it reads it")
	if err != nil {
		t.Fatal(err)
	}
	if len(callsWith(calls, "issue create")) != 0 || !strings.Contains(onlyCall(t, calls, "issue comment "+upstreamURL), "Seen again") {
		t.Errorf("the open issue gets a comment and no new one: %v", calls)
	}
	if !strings.Contains(out, "already reported to Anchors: "+upstreamURL) {
		t.Errorf("output:\n%s", out)
	}
}

func TestReportBug_createsTheIssueUpstream(t *testing.T) {
	t.Run("RPBUG-B03: A new bug becomes an issue at co2-lab/anchors", func(t *testing.T) {})
	t.Run("RPBUG-B06: It works without an anchors.yaml", func(t *testing.T) {})
	out, err, calls := runReportBug(t, t.TempDir(),
		[]ghRule{{match: "issue create --repo co2-lab/anchors *", out: upstreamURL}}, "",
		"gate misreads a file", "--expected", "it reads it")
	if err != nil {
		t.Fatal(err)
	}
	create := onlyCall(t, calls, "issue create --repo co2-lab/anchors")
	if !strings.Contains(create, "--title [bug] gate misreads a file") || !strings.Contains(create, "--label bug") ||
		!strings.Contains(create, "it reads it") {
		t.Errorf("the issue carries the title, the bug label and the body: %s", create)
	}
	if !strings.Contains(out, "reported to Anchors: "+upstreamURL) {
		t.Errorf("output:\n%s", out)
	}
}

func TestReportBug_dryRunAndTheMinimalCase(t *testing.T) {
	t.Run("RPBUG-B04: A dry run prints the issue and sends nothing", func(t *testing.T) {})
	t.Run("RPBUG-B05: The minimal case comes from a file or from standard input", func(t *testing.T) {})
	root := t.TempDir()
	caso := filepath.Join(root, "case.yaml")
	if err := os.WriteFile(caso, []byte("layers:\n  logic: {pattern: \"src/**\"}"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err, calls := runReportBug(t, root, nil, "", "x", "--expected", "y", "--repro", caso, "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if len(callsWith(calls, "issue")) != 0 || !strings.Contains(out, "[bug] x") || !strings.Contains(out, "logic: {pattern") {
		t.Errorf("a dry run prints the issue with the case and sends nothing:\n%s\n%v", out, calls)
	}
	out, err, _ = runReportBug(t, root, nil, "anchors map build\n", "x", "--expected", "y", "--repro", "-", "--dry-run")
	if err != nil || !strings.Contains(out, "```shell\nanchors map build\n```") {
		t.Errorf("the case comes from standard input with -: %v\n%s", err, out)
	}
}

func TestReportBug_refusesWhatNamesTheProject(t *testing.T) {
	t.Run("RPBUG-X01: A report that names the project is refused", func(t *testing.T) {})
	root := githubProject(t) // workflow.repo: acme/app
	_, err, calls := runReportBug(t, root, nil, "", "acme/app breaks the gate", "--expected", "y")
	if err == nil || !strings.Contains(err.Error(), "acme/app") {
		t.Errorf("the project's repository is refused by name, got %v", err)
	}
	_, err, calls2 := runReportBug(t, root, nil, "", "fails under "+root, "--expected", "y")
	if err == nil {
		t.Error("the project's path is refused")
	}
	if len(callsWith(append(calls, calls2...), "issue")) != 0 {
		t.Errorf("nothing is sent: %v %v", calls, calls2)
	}
}

func TestEscalate_upstreamThatNamesTheProjectIsNotSent(t *testing.T) {
	t.Run("RPBUG-B07: An escalation whose reason names the project is not sent upstream", func(t *testing.T) {})
	root := localProject(t)
	calls := scriptedGH(t)
	cmd := newEscalateCmd()
	cmd.SetArgs([]string{"--root", root, "--bug", "--upstream", "breaks under " + root})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if len(callsWith(calls(), "issue create")) != 0 || !strings.Contains(out, "not reported to Anchors") {
		t.Errorf("a reason that names the project stays home:\n%s\n%v", out, calls())
	}
}

func TestReportBug_errors(t *testing.T) {
	t.Run("RPBUG-E01: Without what should happen the report is refused", func(t *testing.T) {})
	t.Run("RPBUG-E02: A refused report leaves the link to file it", func(t *testing.T) {})
	t.Run("RPBUG-E03: An unreadable minimal case is refused", func(t *testing.T) {})
	_, err, calls := runReportBug(t, t.TempDir(), nil, "", "x")
	if err == nil || !strings.Contains(err.Error(), "--expected") || len(callsWith(calls, "issue")) != 0 {
		t.Errorf("--expected is required and nothing is sent: %v %v", err, calls)
	}
	out, err, _ := runReportBug(t, t.TempDir(),
		[]ghRule{{match: "issue create --repo co2-lab/anchors *", out: "HTTP 403", code: 1}}, "", "x", "--expected", "y")
	if err == nil || !strings.Contains(err.Error(), "could not report to Anchors") ||
		!strings.Contains(out, "https://github.com/co2-lab/anchors/issues/new?") {
		t.Errorf("a refusal fails and leaves the link: %v\n%s", err, out)
	}
	_, err, calls = runReportBug(t, t.TempDir(), nil, "", "x", "--expected", "y", "--repro", "missing.txt")
	if err == nil || len(callsWith(calls, "issue")) != 0 {
		t.Errorf("an unreadable case is refused before sending: %v %v", err, calls)
	}
}

func TestRepoOfRemote(t *testing.T) {
	for in, want := range map[string]string{
		"git@github.com:acme/app.git":              "acme/app",
		"https://github.com/acme/app":              "acme/app",
		"ssh://git@github.com/co2-lab/anchors.git": "co2-lab/anchors",
		"nonsense": "",
	} {
		if got := repoOfRemote(in); got != want {
			t.Errorf("repoOfRemote(%q) = %q, want %q", in, got, want)
		}
	}
}

package flow

import (
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// The unblock card is born linked to the blocked one, and the blocked one is told who it
// waits for — the two halves of the link.
func TestUnblockCmd_createsTheLinkedCardAndTellsTheBlockedOne(t *testing.T) {
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue create *", out: "https://github.com/acme/app/issues/512"},
	)
	cmd := newUnblockCmd()
	cmd.SetArgs([]string{"--root", root, "#311",
		"--reason", "the dispatch loop needs try/catch per token\nsecond paragraph",
		"--about", "src/push/Dispatcher.ts"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}

	c := calls()
	label := initx.LabelDesbloqueia("311")
	if got := callsWith(c, "label create "+label); len(got) != 1 {
		t.Errorf("the link label must be created on demand; calls: %v", c)
	}
	create := callsWith(c, "issue create")
	if len(create) != 1 {
		t.Fatalf("expected one card created; calls: %v", c)
	}
	for _, want := range []string{
		"--title [unblocks #311] the dispatch loop needs try/catch per token --body-file",
		"--label anchors", "--label anchors:to-do", "--label " + label,
	} {
		if !strings.Contains(create[0], want) {
			t.Errorf("the card must carry %q: %s", want, create[0])
		}
	}
	comment := callsWith(c, "issue comment 311")
	if len(comment) != 1 || !strings.Contains(comment[0], "https://github.com/acme/app/issues/512") {
		t.Errorf("the blocked card must name the card it waits for; calls: %v", c)
	}
	if !strings.Contains(out, "unblock card created: https://github.com/acme/app/issues/512") ||
		!strings.Contains(out, "#311 remains stopped") {
		t.Errorf("the output must give the new card and the state of the blocked one:\n%s", out)
	}
}

// The card was created; a failed comment is a warning, not an error that hides it.
func TestUnblockCmd_failedCommentWarnsWithoutFailing(t *testing.T) {
	root := githubProject(t)
	scriptedGH(t,
		ghRule{match: "issue create *", out: "https://github.com/acme/app/issues/512"},
		ghRule{match: "issue comment *", code: 1},
	)
	cmd := newUnblockCmd()
	cmd.SetArgs([]string{"--root", root, "311", "--reason", "x"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatalf("the card exists; the command must not fail: %v", err)
	}
	if !strings.Contains(out, "could not comment on #311") {
		t.Errorf("the failed comment must be said:\n%s", out)
	}
}

func TestUnblockCmd_refusals(t *testing.T) {
	scriptedGH(t, ghRule{match: "issue create *", out: "create failed", code: 1})

	cmd := newUnblockCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "311", "--reason", " "})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "--reason") {
		t.Errorf("a blank reason must be refused, got %v", err)
	}

	cmd = newUnblockCmd()
	cmd.SetArgs([]string{"--root", localProject(t), "311", "--reason", "x"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode must be refused, got %v", err)
	}

	cmd = newUnblockCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "311", "--reason", "x"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "create the card") ||
		!strings.Contains(err.Error(), "create failed") {
		t.Errorf("a failed create must say so with gh's answer, got %v", err)
	}
}

// The body says what to do, where, and how the blocked card gets back to the queue.
func TestCorpoDoDesbloqueio(t *testing.T) {
	p, err := corpoDoDesbloqueio("311", "add the retry rule", "src/a.ts")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(p)
	b, _ := os.ReadFile(p)
	body := string(b)
	for _, want := range []string{
		"unblocks #311", "add the retry rule", "**Where:** `src/a.ts`",
		"remove the label `anchors:needs-user` from #311",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the body lacks %q:\n%s", want, body)
		}
	}

	p2, err := corpoDoDesbloqueio("311", "x", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(p2)
	b2, _ := os.ReadFile(p2)
	if strings.Contains(string(b2), "**Where:**") {
		t.Errorf("without --about there is no Where line:\n%s", b2)
	}
}

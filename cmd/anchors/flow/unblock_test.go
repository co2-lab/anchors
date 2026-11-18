// @anchors
//   code: UNTSN
//   ref: NBLCK

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
	t.Run("NBLCK-B05: The new card carries the unblock title and the three labels", func(t *testing.T) {})
	t.Run("NBLCK-B07: The blocked card is told which card it waits for", func(t *testing.T) {})
	t.Run("NBLCK-B08: The output gives the new card and the state of the blocked one", func(t *testing.T) {})
	t.Run("NBLCK-B09: The leading hash of the card is stripped", func(t *testing.T) {})
	t.Run("NBLCK-I01: The blocked card keeps its needs-user label", func(t *testing.T) {})
	t.Run("NBLCK-X01: The temporary body file does not survive the command", func(t *testing.T) {})
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
	if len(comment) != 1 || !strings.Contains(comment[0], "https://github.com/acme/app/issues/512") ||
		!strings.Contains(comment[0], "`anchors:needs-user` STAYS") {
		t.Errorf("the blocked card must name the card it waits for and keep its label; calls: %v", c)
	}
	if got := callsWith(c, "remove-label"); len(got) != 0 {
		t.Errorf("unblock must never remove a label from the blocked card: %v", got)
	}
	// The body travels in a temporary file that is gone once the command returns.
	bodyFile := strings.Fields(strings.SplitN(create[0], "--body-file ", 2)[1])[0]
	if _, err := os.Stat(bodyFile); !os.IsNotExist(err) {
		t.Errorf("the temporary body file %s survived the command (stat: %v)", bodyFile, err)
	}
	if !strings.Contains(out, "unblock card created: https://github.com/acme/app/issues/512") ||
		!strings.Contains(out, "#311 remains stopped") {
		t.Errorf("the output must give the new card and the state of the blocked one:\n%s", out)
	}
}

// The card was created; a failed comment is a warning, not an error that hides it.
func TestUnblockCmd_failedCommentWarnsWithoutFailing(t *testing.T) {
	t.Run("NBLCK-E02: A failed comment on the blocked card is a warning", func(t *testing.T) {})
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
	t.Run("NBLCK-B02: Unblock with a blank reason is refused", func(t *testing.T) {})
	t.Run("NBLCK-B03: Unblock in local mode is refused", func(t *testing.T) {})
	t.Run("NBLCK-E01: A card the platform refuses to create fails the command", func(t *testing.T) {})
	calls := scriptedGH(t, ghRule{match: "issue create *", out: "create failed", code: 1})

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
	if got := callsWith(calls(), "issue comment"); len(got) != 0 {
		t.Errorf("no card was created, so the blocked one must not be told to wait: %v", got)
	}
}

// The body says what to do, where, and how the blocked card gets back to the queue.
func TestUnblockBody(t *testing.T) {
	t.Run("NBLCK-B06: The body says what to do, where it came from and how the blocked card returns", func(t *testing.T) {})
	p, err := corpoDoDesbloqueio("311", "add the retry rule", "src/a.ts")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(p)
	b, _ := os.ReadFile(p)
	body := string(b)
	for _, want := range []string{
		"unblocks #311", "add the retry rule", "**Where:** `src/a.ts`", "**Where it came from:** #311",
		"**What to do:**",
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

// Exactly one card: the link label and the comment are about ONE blocked card.
func TestUnblockCmd_takesExactlyOneCard(t *testing.T) {
	t.Run("NBLCK-B01: Unblock with no card or with two cards is refused", func(t *testing.T) {})
	calls := scriptedGH(t)
	for _, args := range [][]string{{}, {"311", "312"}} {
		cmd := newUnblockCmd()
		cmd.SilenceErrors, cmd.SilenceUsage = true, true
		cmd.SetArgs(append([]string{"--root", githubProject(t), "--reason", "x"}, args...))
		if err := cmd.Execute(); err == nil {
			t.Errorf("unblock with cards %v must fail", args)
		}
	}
	if got := callsWith(calls(), "issue create"); len(got) != 0 {
		t.Errorf("no card may be created: %v", got)
	}
}

// The link label is one per blocked card; from the second unblock card on it already
// exists, and creating it again fails without stopping anything.
func TestUnblockCmd_existingLinkLabelDoesNotStopIt(t *testing.T) {
	t.Run("NBLCK-B04: The link label is created on demand, and an existing one does not stop the command", func(t *testing.T) {})
	calls := scriptedGH(t,
		ghRule{match: "label create *", out: "already exists", code: 1},
		ghRule{match: "issue create *", out: "https://github.com/acme/app/issues/512"},
	)
	cmd := newUnblockCmd()
	cmd.SetArgs([]string{"--root", githubProject(t), "311", "--reason", "x"})
	var err error
	stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatalf("an existing link label must not stop the command: %v", err)
	}
	c := calls()
	if len(c) < 2 || !strings.Contains(c[0], "label create "+initx.LabelDesbloqueia("311")) ||
		!strings.Contains(c[1], "issue create") {
		t.Errorf("the link label is created before the card; calls: %v", c)
	}
}

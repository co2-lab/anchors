package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/queue"
)

// The full check says what is STILL OPEN locally — not only what this run opened: the
// issues in todo/doing (with those waiting for the user apart) and the tasks not done.
// Before, a backlog from earlier runs stayed silent under a clean-looking check.
func TestLocalBacklog(t *testing.T) {
	root := t.TempDir()
	if b := readLocalBacklog(root); !b.empty() {
		t.Fatalf("an empty project has no backlog: %+v", b)
	}
	if out := captureStdout(t, func() { printLocalBacklog(readLocalBacklog(root)) }); out != "" {
		t.Errorf("no backlog must print nothing, printed:\n%s", out)
	}

	for _, i := range []issue.Issue{
		{Kind: issue.Violation, Gate: "g", Target: "a.ts", Date: "2026-09-25"},
		{Kind: issue.Decision, Target: "b.spec.md", Date: "2026-09-25", Dono: issue.DonoUsuário},
	} {
		if _, _, err := issue.Open(root, i); err != nil {
			t.Fatal(err)
		}
	}
	doing := filepath.Join(root, issue.Dir, string(issue.Doing))
	if err := os.MkdirAll(doing, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(doing, "c.md"), []byte("# c\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The queue ignores a task whose target is gone, so the targets exist.
	for i, c := range []string{"x.ts", "y.ts"} {
		if err := os.WriteFile(filepath.Join(root, c), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		id := []string{"0001-code-x", "0002-code-y"}[i]
		if _, err := queue.Enqueue(root, queue.Task{ID: id, Changed: c, Kind: "code", Origin: "manual", CreatedAt: "2026-09-25T10:00:00Z"}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := queue.Claim(root, "w1", "2026-09-25T10:00:00Z"); err != nil {
		t.Fatal(err)
	}

	b := readLocalBacklog(root)
	if b.Todo != 2 || b.ForUser != 1 || b.Doing != 1 || b.Pending != 1 || b.Claimed != 1 {
		t.Fatalf("backlog = %+v, want todo 2 (1 for the user), doing 1, 1 pending, 1 claimed", b)
	}
	out := captureStdout(t, func() { printLocalBacklog(b) })
	for _, want := range []string{"issues/todo/", "issues/doing/", "anchors next"} {
		if !strings.Contains(out, want) {
			t.Errorf("the backlog should mention %q:\n%s", want, out)
		}
	}
}

// Wired on the full sweep of the LOCAL mode only: in github mode the backlog is the board,
// and on `--changed` (every pre-commit) the same lines would be noise.
func TestLocalBacklogIsPrintedOnTheFullLocalCheckOnly(t *testing.T) {
	b, err := os.ReadFile("check.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	i := strings.Index(s, "printLocalBacklog(readLocalBacklog(absRoot))")
	if i < 0 {
		t.Fatal("check no longer prints the local backlog")
	}
	guard := s[max(0, i-120):i]
	if !strings.Contains(guard, "if all && !cfg.GitHubMode() {") {
		t.Errorf("the backlog must be printed under `all && !GitHubMode()`, found:\n%s", guard)
	}
}

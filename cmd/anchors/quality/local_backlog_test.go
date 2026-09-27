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
	t.Run("LCBCL-B01: The issues in todo and doing are counted, with the user-owned ones apart", func(t *testing.T) {})
	t.Run("LCBCL-B02: The pending and claimed tasks are counted", func(t *testing.T) {})
	t.Run("LCBCL-B03: A project with nothing open prints no backlog", func(t *testing.T) {})
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

func TestLocalBacklogPrintsOnlyTheSideWithWork(t *testing.T) {
	t.Run("LCBCL-B04: Only the side that has something open gets its line", func(t *testing.T) {})
	out := captureStdout(t, func() { printLocalBacklog(localBacklog{Pending: 2}) })
	if !strings.Contains(out, "anchors next") {
		t.Errorf("the tasks line is missing:\n%s", out)
	}
	if strings.Contains(out, "/todo/") || strings.Contains(out, "/doing/") {
		t.Errorf("with no issue open, the issues line must not be printed:\n%s", out)
	}
	out = captureStdout(t, func() { printLocalBacklog(localBacklog{Doing: 1}) })
	if !strings.Contains(out, "/doing/") || strings.Contains(out, "anchors next") {
		t.Errorf("with only an issue in doing, only the issues line is printed:\n%s", out)
	}
}

func TestLocalBacklogSubsetCountsAloneAreEmpty(t *testing.T) {
	t.Run("LCBCL-I01: User-owned and past-window counts alone do not make a backlog", func(t *testing.T) {})
	b := localBacklog{ForUser: 3, Old: 2}
	if !b.empty() {
		t.Fatalf("user-owned and past-window counts alone must leave the backlog empty: %+v", b)
	}
	if out := captureStdout(t, func() { printLocalBacklog(b) }); out != "" {
		t.Errorf("an empty backlog prints nothing, printed:\n%s", out)
	}
}

func TestLocalBacklogReadsWithoutChanging(t *testing.T) {
	t.Run("LCBCL-X01: Reading the backlog changes no issue and no task", func(t *testing.T) {})
	root := t.TempDir()
	if _, _, err := issue.Open(root, issue.Issue{Kind: issue.Violation, Gate: "g", Target: "a.ts", Date: "2026-09-25"}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x.ts"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Enqueue(root, queue.Task{ID: "0001-code-x", Changed: "x.ts", Kind: "code", Origin: "manual", CreatedAt: "2026-09-25T10:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	snapshot := func() map[string]string {
		m := map[string]string{}
		filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				b, _ := os.ReadFile(p)
				m[p] = string(b)
			}
			return nil
		})
		return m
	}
	before := snapshot()
	b := readLocalBacklog(root)
	if b.Todo != 1 || b.Pending != 1 {
		t.Fatalf("backlog = %+v, want 1 todo and 1 pending", b)
	}
	after := snapshot()
	if len(before) != len(after) {
		t.Fatalf("reading the backlog changed the files: %d before, %d after", len(before), len(after))
	}
	for p, c := range before {
		if after[p] != c {
			t.Errorf("reading the backlog changed %s", p)
		}
	}
}

func TestLocalBacklogUnlistableFolderCountsZero(t *testing.T) {
	t.Run("LCBCL-E01: An issue folder that cannot be listed counts as zero", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, issue.Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, issue.Dir, string(issue.Todo)), []byte("not a folder"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := issue.List(root, issue.Todo); err == nil {
		t.Fatal("precondition: listing a todo that is a file must fail")
	}
	if b := readLocalBacklog(root); b.Todo != 0 || !b.empty() {
		t.Errorf("an unlistable todo must count zero, got %+v", b)
	}
}

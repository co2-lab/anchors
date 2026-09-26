package ops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/suggestion"
)

// suggestRepo is a repository with `a.txt` committed and one PENDING suggestion whose
// patch turns "old" into "new".
func suggestRepo(t *testing.T) string {
	t.Helper()
	root := newGitRepo(t)
	writeFile(t, root, "a.txt", "old\n")
	if out, err := runGit(root, "add", "a.txt"); err != nil {
		t.Fatal(out)
	}
	if out, err := runGit(root, "commit", "-q", "-m", "a"); err != nil {
		t.Fatal(out)
	}
	writeFile(t, root, "a.txt", "new\n")
	patch, err := runGit(root, "diff")
	if err != nil {
		t.Fatal(patch)
	}
	writeFile(t, root, "a.txt", "old\n")
	if _, _, err := suggestion.Open(root, suggestion.Suggestion{
		ID: "fix-a", Gate: "demo", Target: "a.txt", Origin: suggestion.FromGate,
		Why: "the word is outdated", Patch: patch,
	}); err != nil {
		t.Fatal(err)
	}
	return root
}

func readA(t *testing.T, root string) string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(root, "a.txt"))
	return string(b)
}

func stateOf(t *testing.T, root, id string) suggestion.State {
	t.Helper()
	for _, st := range []suggestion.State{suggestion.Pending, suggestion.Approved, suggestion.Rejected} {
		if _, err := os.Stat(filepath.Join(root, suggestion.Dir, string(st), id+".md")); err == nil {
			return st
		}
	}
	return ""
}

func TestSuggestListAndShow(t *testing.T) {
	err, out := runCmd(t, newSuggestCmd(), "list", "--root", t.TempDir())
	if err != nil || !strings.Contains(out, "no suggestion in pending") {
		t.Errorf("empty list: %v\n%s", err, out)
	}

	root := suggestRepo(t)
	err, out = runCmd(t, newSuggestCmd(), "list", "--root", root)
	if err != nil || !strings.Contains(out, "1 suggestion(s) in pending") || !strings.Contains(out, "  fix-a") {
		t.Errorf("list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "anchors suggest show <id>") {
		t.Errorf("the pending list should point at `show`:\n%s", out)
	}

	err, out = runCmd(t, newSuggestCmd(), "show", "--root", root, "fix-a")
	if err != nil || !strings.Contains(out, "the word is outdated") || !strings.Contains(out, "+new") {
		t.Errorf("show must print the why and the diff: %v\n%s", err, out)
	}
	if err, _ := runCmd(t, newSuggestCmd(), "show", "--root", root, "nope"); err == nil ||
		!strings.Contains(err.Error(), `"nope" not found in pending`) {
		t.Errorf("show of a missing id: %v", err)
	}
}

// --dry-run only checks: the file and the state stay as they were.
func TestSuggestApplyDryRunChangesNothing(t *testing.T) {
	root := suggestRepo(t)
	err, out := runCmd(t, newSuggestCmd(), "apply", "--root", root, "fix-a", "--dry-run")
	if err != nil || !strings.Contains(out, "applies cleanly") {
		t.Fatalf("dry run: %v\n%s", err, out)
	}
	if readA(t, root) != "old\n" || stateOf(t, root, "fix-a") != suggestion.Pending {
		t.Errorf("--dry-run changed something: a.txt = %q, state = %s", readA(t, root), stateOf(t, root, "fix-a"))
	}
}

// Applying changes the file FIRST and only then approves, with the default reason.
func TestSuggestApplyPatchesAndApproves(t *testing.T) {
	root := suggestRepo(t)
	err, out := runCmd(t, newSuggestCmd(), "apply", "--root", root, "fix-a")
	if err != nil {
		t.Fatalf("apply: %v\n%s", err, out)
	}
	if readA(t, root) != "new\n" {
		t.Errorf("the patch was not applied: %q", readA(t, root))
	}
	if st := stateOf(t, root, "fix-a"); st != suggestion.Approved {
		t.Fatalf("state = %s, want approved", st)
	}
	b, _ := os.ReadFile(filepath.Join(root, suggestion.Dir, "approved", "fix-a.md"))
	if !strings.Contains(string(b), "patch applied without reservations") {
		t.Errorf("the approval carries no reason:\n%s", b)
	}
}

// A patch that no longer matches fails BEFORE touching anything, and says why.
func TestSuggestApplyStalePatchTouchesNothing(t *testing.T) {
	root := suggestRepo(t)
	writeFile(t, root, "a.txt", "changed meanwhile\n")
	err, _ := runCmd(t, newSuggestCmd(), "apply", "--root", root, "fix-a")
	if err == nil || !strings.Contains(err.Error(), "no longer matches") {
		t.Fatalf("want the stale-patch error, got %v", err)
	}
	if readA(t, root) != "changed meanwhile\n" || stateOf(t, root, "fix-a") != suggestion.Pending {
		t.Error("a stale patch changed the file or the state")
	}
}

// A suggestion IS a patch: outside git the error says git is what is missing.
func TestSuggestApplyOutsideGitNamesTheCause(t *testing.T) {
	root := suggestRepo(t)
	if err := os.RemoveAll(filepath.Join(root, ".git")); err != nil {
		t.Fatal(err)
	}
	err := gitApply(root, "not a patch", true)
	if err == nil || !strings.Contains(err.Error(), "apply the suggestion patch") {
		t.Errorf("want the explanation naming the suggestion patch, got %v", err)
	}
}

// Rejecting requires the reason, and the rejected one is KEPT, with it.
func TestSuggestRejectKeepsTheRecordWithTheReason(t *testing.T) {
	root := suggestRepo(t)
	if err, _ := runCmd(t, newSuggestCmd(), "reject", "--root", root, "fix-a"); err == nil {
		t.Error("a rejection without a reason must be refused")
	}
	if stateOf(t, root, "fix-a") != suggestion.Pending {
		t.Fatal("the refused rejection moved the suggestion")
	}
	err, out := runCmd(t, newSuggestCmd(), "reject", "--root", root, "fix-a", "--reason", "the old word is the brand", "--auto")
	if err != nil || !strings.Contains(out, "fix-a rejected") {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	b, _ := os.ReadFile(filepath.Join(root, suggestion.Dir, "rejected", "fix-a.md"))
	if !strings.Contains(string(b), "the old word is the brand") || !strings.Contains(string(b), "auto_judgment") {
		t.Errorf("the record lacks the reason or the automatic marker:\n%s", b)
	}
	err, out = runCmd(t, newSuggestCmd(), "list", "--root", root, "--state", "rejected")
	if err != nil || !strings.Contains(out, "fix-a") || strings.Contains(out, "suggest show") {
		t.Errorf("list --state rejected: %v\n%s", err, out)
	}
}

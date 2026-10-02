// @anchors
//   ref: SGSTS

package suggestion

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fresh(id string) Suggestion {
	return Suggestion{
		ID: id, Gate: "mock-detect-covers-dialect", Target: "anchors.yaml",
		Origin: FromJudgment,
		Why:    "the declared pattern does not match `unittest.mock.patch`, used in 12 tests",
		Patch:  "--- a/anchors.yaml\n+++ b/anchors.yaml\n@@\n-  mock_detect: \"x\"\n+  mock_detect: \"y\"",
	}
}

func TestOpenWritesToPending(t *testing.T) {
	t.Run("SGSTS-B01: A new suggestion is born in pending with its context and reason", func(t *testing.T) {})
	t.Run("SGSTS-B02: A suggestion with a patch carries it in a diff block", func(t *testing.T) {})
	root := t.TempDir()
	created, p, err := Open(root, fresh("s1"))
	if err != nil || !created {
		t.Fatalf("should create: %v %v", created, err)
	}
	if p != filepath.Join(root, "suggestions", "pending", "s1.md") {
		t.Errorf("a suggestion is born in pending: %s", p)
	}
	b, _ := os.ReadFile(p)
	for _, want := range []string{
		"# SUGGESTION: anchors.yaml", "**gate:** mock-detect-covers-dialect", "**origin:** judgment",
		"**target:** anchors.yaml", "## Why\n\nthe declared pattern", "```diff\n--- a/anchors.yaml",
	} {
		if !strings.Contains(string(b), want) {
			t.Errorf("the file lacks %q:\n%s", want, b)
		}
	}
}

// Idempotent: the same proposal reappearing in a later scan does not duplicate.
func TestOpenDoesNotDuplicate(t *testing.T) {
	t.Run("SGSTS-B04: Opening the same pending suggestion twice creates it once", func(t *testing.T) {})
	root := t.TempDir()
	_, first, _ := Open(root, fresh("s1"))
	created, p, err := Open(root, fresh("s1"))
	if err != nil {
		t.Fatal(err)
	}
	if created || p != first {
		t.Errorf("the same suggestion cannot be created twice (created=%v, path %s)", created, p)
	}
}

// An already REJECTED proposal does not go back to pending. Reopening would erase the
// decision of whoever already looked, and it would come back at every scan as new.
func TestOpenDoesNotReopenADecidedOne(t *testing.T) {
	t.Run("SGSTS-B05: A rejected suggestion is never reopened", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("s1"))
	if err := Decide(root, "s1", Rejected, "the dialect is intentional here", false); err != nil {
		t.Fatal(err)
	}
	created, _, _ := Open(root, fresh("s1"))
	if created {
		t.Error("an already decided suggestion cannot reopen")
	}
	pend, _ := List(root, Pending)
	if len(pend) != 0 {
		t.Errorf("pending should be empty: %v", pend)
	}
}

// Moving is deciding — the state is the folder, as with issues. There is no status field in
// the file that could disagree with where it is.
func TestDecideMovesAndRecords(t *testing.T) {
	t.Run("SGSTS-B06: Deciding moves the suggestion and records the reason and the decider", func(t *testing.T) {})
	t.Run("SGSTS-I01: A decided suggestion is no longer pending", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("s1"))
	if err := Decide(root, "s1", Approved, "the pattern really was wrong", false); err != nil {
		t.Fatal(err)
	}
	ap, _ := List(root, Approved)
	if len(ap) != 1 || ap[0] != "s1" {
		t.Fatalf("should be in approved: %v", ap)
	}
	if pend, _ := List(root, Pending); len(pend) != 0 {
		t.Fatalf("a decided suggestion must leave pending: %v", pend)
	}
	b, _ := os.ReadFile(filepath.Join(root, Dir, string(Approved), "s1.md"))
	for _, want := range []string{"## Decision", "**state:** approved", "the pattern really was wrong", "by:** pessoa"} {
		if !strings.Contains(string(b), want) {
			t.Errorf("the decision must record %q: %s", want, b)
		}
	}
}

// A decision without a reason is the same failure as a bare `@no-test`: the trace of why
// someone chose disappears, and the choice becomes indistinguishable from carelessness.
func TestDecideRequiresAReason(t *testing.T) {
	t.Run("SGSTS-E03: A decision without a reason is refused", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("s1"))
	if err := Decide(root, "s1", Approved, "   ", false); err == nil {
		t.Error("a decision without a reason should fail")
	}
	if pend, _ := List(root, Pending); !reflect.DeepEqual(pend, []string{"s1"}) {
		t.Errorf("a refused decision must leave the suggestion pending: %v", pend)
	}
}

// Under `auto_judgment` the decision is the AI's, and that stays MARKED. Erasing the
// distinction would make "nobody looked at this" look like "someone approved".
func TestDecideMarksAutomaticJudgment(t *testing.T) {
	t.Run("SGSTS-B07: An automatic decision is marked as the AI's", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("s1"))
	if err := Decide(root, "s1", Approved, "pattern clearly incomplete", true); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(root, Dir, string(Approved), "s1.md"))
	if !strings.Contains(string(b), "by:** IA (auto_judgment)") {
		t.Errorf("an automatic decision must be distinguishable from a human one: %s", b)
	}
}

func TestPatchOfExtractsTheDiff(t *testing.T) {
	t.Run("SGSTS-B09: The patch comes out clean for git apply", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("s1"))
	p, err := PatchOf(root, "s1", Pending)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(p, "--- a/anchors.yaml") || strings.Contains(p, "```") {
		t.Errorf("the patch must come out clean, applicable with git apply: %q", p)
	}
}

// A finding whose fix is NOT mechanical is still worth a recorded diagnosis — the suggestion
// without a patch says so explicitly instead of faking a fix.
func TestSuggestionWithoutPatchIsADiagnosis(t *testing.T) {
	t.Run("SGSTS-B03: A suggestion without a patch says the fix needs a human decision", func(t *testing.T) {})
	root := t.TempDir()
	s := fresh("s2")
	s.Patch = ""
	_, p, _ := Open(root, s)
	b, _ := os.ReadFile(p)
	if strings.Contains(string(b), "```diff") {
		t.Error("without a patch there must be no diff block")
	}
	if !strings.Contains(string(b), "needs a human decision") {
		t.Errorf("it must say the fix is not mechanical: %s", b)
	}
}

func TestList_sortedIDsOfAState(t *testing.T) {
	t.Run("SGSTS-B08: Listing a state gives its sorted IDs", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("b"))
	_, _, _ = Open(root, fresh("a"))
	os.WriteFile(filepath.Join(root, Dir, string(Pending), "notes.txt"), []byte("x"), 0o644)
	got, err := List(root, Pending)
	if err != nil || !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("List(pending) = %v, %v; want [a b]", got, err)
	}
	if got, err := List(root, Rejected); err != nil || got != nil {
		t.Errorf("List of a missing state folder = %v, %v; want nothing, no error", got, err)
	}
}

func TestOpen_refusesAnEmptyID(t *testing.T) {
	t.Run("SGSTS-E01: A suggestion without an ID is refused", func(t *testing.T) {})
	root := t.TempDir()
	if _, _, err := Open(root, fresh("")); err == nil || !strings.Contains(err.Error(), "without ID") {
		t.Fatalf("Open without ID = %v, want the refusal", err)
	}
	if _, err := os.Stat(filepath.Join(root, Dir)); !os.IsNotExist(err) {
		t.Errorf("nothing may be written for a refused suggestion: %v", err)
	}
}

func TestDecide_refusesAnUnknownState(t *testing.T) {
	t.Run("SGSTS-E02: A decision to an unknown state is refused", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, fresh("s1"))
	if err := Decide(root, "s1", Pending, "a reason", false); err == nil || !strings.Contains(err.Error(), "invalid decision state") {
		t.Fatalf("Decide to pending = %v, want the refusal", err)
	}
	if pend, _ := List(root, Pending); !reflect.DeepEqual(pend, []string{"s1"}) {
		t.Errorf("the suggestion must stay pending: %v", pend)
	}
}

func TestDecide_refusesWhatIsNotPending(t *testing.T) {
	t.Run("SGSTS-E04: Deciding a suggestion that is not pending is refused", func(t *testing.T) {})
	if err := Decide(t.TempDir(), "ghost", Approved, "a reason", false); err == nil || !strings.Contains(err.Error(), "pending suggestion not found") {
		t.Fatalf("Decide of a missing suggestion = %v, want the refusal", err)
	}
}

func TestPatchOf_refusesASuggestionWithoutPatch(t *testing.T) {
	t.Run("SGSTS-E05: Asking the patch of a suggestion without one is refused", func(t *testing.T) {})
	root := t.TempDir()
	s := fresh("s2")
	s.Patch = ""
	_, _, _ = Open(root, s)
	if _, err := PatchOf(root, "s2", Pending); err == nil || !strings.Contains(err.Error(), "has no patch") {
		t.Fatalf("PatchOf without a patch = %v, want the refusal", err)
	}
}

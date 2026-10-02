// @anchors
//   ref: TSQUT

package queue

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func task(id, changed, kind, next string) Task {
	return Task{ID: id, Changed: changed, Kind: kind, Origin: "watch", SuggestedNext: next, CreatedAt: "2026-08-07T00:00:00Z"}
}

// withTarget creates the target file on disk. Since `List` discards a task whose target does
// not exist (the ghost task left over in every real run), the target has to be material.
func withTarget(t *testing.T, root, rel string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestEnqueueAndList(t *testing.T) {
	t.Run("TSQUT-B01: An enqueued task is born pending", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "AddItem.spec.md")
	created, err := Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	if err != nil || !created {
		t.Fatalf("enqueue: created=%v err=%v", created, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".anchors", "tasks", "pending__1-spec-add.yaml")); err != nil {
		t.Fatalf("the task file must be pending__1-spec-add.yaml: %v", err)
	}
	tasks, err := List(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].State != Pending || tasks[0].Changed != "AddItem.spec.md" {
		t.Fatalf("want 1 pending task; got %+v", tasks)
	}
}

func TestList_sortedAndEmptyWithoutAQueue(t *testing.T) {
	t.Run("TSQUT-B03: Listing gives the live tasks sorted, and nothing without a queue", func(t *testing.T) {})
	root := t.TempDir()
	if tasks, err := List(root); err != nil || tasks != nil {
		t.Fatalf("no queue folder: %v, %v; want nothing, no error", tasks, err)
	}
	withTarget(t, root, "a.md")
	withTarget(t, root, "b.md")
	_, _ = Enqueue(root, task("2-b", "b.md", "doc", "triage"))
	_, _ = Enqueue(root, task("1-a", "a.md", "doc", "triage"))
	tasks, _ := List(root)
	if len(tasks) != 2 || tasks[0].ID != "1-a" || tasks[1].ID != "2-b" {
		t.Fatalf("want 1-a then 2-b, got %+v", tasks)
	}
}

func TestEnqueueDedup(t *testing.T) {
	t.Run("TSQUT-B02: The same target and step are not enqueued twice", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "AddItem.spec.md")
	_, _ = Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	// the same (changed, suggested_next) with a different ID → NOT duplicated
	created, err := Enqueue(root, task("2-spec-add", "AddItem.spec.md", "spec", "implement"))
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("must not duplicate a live task for the same target and step")
	}
	if n, _ := PendingCount(root); n != 1 {
		t.Fatalf("want 1 pending, got %d", n)
	}

	// Two judgment gates on the same spec are two tasks; a judgment is not a review of
	// another kind of the same file.
	j1 := Task{ID: "judge-no-test-proof-real-addItem", Changed: "AddItem.spec.md", Kind: KindJudgment, SuggestedNext: "review", Reason: "r"}
	j2 := Task{ID: "judge-rule-fulfilled-addItem", Changed: "AddItem.spec.md", Kind: KindJudgment, SuggestedNext: "review", Reason: "r"}
	w := Task{ID: "3-spec-add", Changed: "AddItem.spec.md", Kind: "spec", SuggestedNext: "review", Reason: "r"}
	for _, tk := range []Task{w, j1, j2} {
		if created, err := Enqueue(root, tk); err != nil || !created {
			t.Errorf("%s must be enqueued, created=%v err=%v", tk.ID, created, err)
		}
	}
	if created, _ := Enqueue(root, j1); created {
		t.Error("the same judgment enqueued twice must not duplicate")
	}
}

func TestClaimEmpty(t *testing.T) {
	t.Run("TSQUT-B05: Claiming takes a pending task and records the worker", func(t *testing.T) {})
	root := t.TempDir()
	got, err := Claim(root, "w1", "2026-08-13T00:00:00-03:00")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("an empty queue should give nil, got %+v", got)
	}
}

func TestClaimMovesToClaimed(t *testing.T) {
	t.Run("TSQUT-B05: Claiming takes a pending task and records the worker", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "AddItem.spec.md")
	_, _ = Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	got, err := Claim(root, "worker-A", "2026-08-13T00:00:00-03:00")
	if err != nil || got == nil {
		t.Fatalf("claim: %+v err=%v", got, err)
	}
	if got.State != Claimed || got.ClaimedBy != "worker-A" || got.ClaimedAt != "2026-08-13T00:00:00-03:00" {
		t.Fatalf("want claimed by worker-A at the given moment, got %+v", got)
	}
	// the pending file is gone, the claimed one exists
	d := dirFor(root)
	if _, err := os.Stat(filepath.Join(d, fileName(Pending, "1-spec-add"))); !os.IsNotExist(err) {
		t.Fatal("the pending file should be gone")
	}
	if _, err := os.Stat(filepath.Join(d, fileName(Claimed, "1-spec-add"))); err != nil {
		t.Fatal("the claimed file should exist")
	}
}

// N concurrent workers over M tasks never take the same task. It is the guarantee that allows
// two terminals running `anchors next`.
func TestClaimAtomicNoDoubleClaim(t *testing.T) {
	t.Run("TSQUT-I01: Concurrent workers never claim the same task", func(t *testing.T) {})
	root := t.TempDir()
	const M = 20
	for i := 0; i < M; i++ {
		id := filepath.Base(t.Name()) + "-" + string(rune('a'+i))
		target := "f" + string(rune('a'+i)) + ".spec.md"
		withTarget(t, root, target)
		_, _ = Enqueue(root, task(id, target, "spec", "implement"))
	}
	var mu sync.Mutex
	claimed := map[string]int{}
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for {
				got, err := Claim(root, "w", "2026-08-13T00:00:00-03:00")
				if err != nil || got == nil {
					return
				}
				mu.Lock()
				claimed[got.ID]++
				mu.Unlock()
			}
		}(w)
	}
	wg.Wait()
	if len(claimed) != M {
		t.Fatalf("want %d tasks claimed, got %d", M, len(claimed))
	}
	for id, n := range claimed {
		if n != 1 {
			t.Fatalf("task %s claimed %d times (double claim!)", id, n)
		}
	}
}

// The worst case of replacing the rename with O_EXCL: if the process dies between creating
// claimed__X and deleting pending__X, both files coexist. What must NOT happen is the pending
// residue being served as if the task were free — the double claim back through another door.
func TestClaimDoesNotTakeWhatAlreadyHasAnOwner(t *testing.T) {
	t.Run("TSQUT-B06: A pending residue of a dead claim is not served and is cleaned", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "fa.spec.md")
	if _, err := Enqueue(root, task("residuo-a", "fa.spec.md", "spec", "implement")); err != nil {
		t.Fatal(err)
	}

	// Simulates the death midway: the claimed file exists, the pending one stayed behind.
	d := dirFor(root)
	claimed := filepath.Join(d, fileName(Claimed, "residuo-a"))
	if err := os.WriteFile(claimed, []byte("id: residuo-a\nstate: claimed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	pending := filepath.Join(d, fileName(Pending, "residuo-a"))
	if _, err := os.Stat(pending); err != nil {
		t.Fatalf("the scenario needs the pending file still on disk: %v", err)
	}

	got, err := Claim(root, "w2", "2026-08-22T00:00:00-03:00")
	if err != nil {
		t.Fatalf("Claim returned an error: %v", err)
	}
	if got != nil {
		t.Fatalf("claimed %q, which already has an owner", got.ID)
	}
	// And the residue cleans itself: whoever meets it deletes it, instead of leaving the queue
	// forever showing a pending task nobody can take.
	if _, err := os.Stat(pending); !os.IsNotExist(err) {
		t.Error("the pending residue should have been removed on meeting the claimed file")
	}
}

func TestMarkDoneMovesToHistory(t *testing.T) {
	t.Run("TSQUT-B07: A done task moves to the history", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "AddItem.spec.md")
	_, _ = Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	_, _ = Claim(root, "w1", "2026-08-13T00:00:00-03:00")
	if err := MarkDone(root, "1-spec-add"); err != nil {
		t.Fatal(err)
	}
	// gone from the live queue
	if n, _ := PendingCount(root); n != 0 {
		t.Fatalf("want an empty queue after done, got %d", n)
	}
	// it appeared in the history
	donePath := filepath.Join(root, DoneDir, fileName(Done, "1-spec-add"))
	if _, err := os.Stat(donePath); err != nil {
		t.Fatalf("the done task should be in %s: %v", DoneDir, err)
	}
	// and a new task for the same target IS allowed now (the previous one left the queue)
	created, _ := Enqueue(root, task("2-spec-add", "AddItem.spec.md", "spec", "implement"))
	if !created {
		t.Fatal("after done, a new enqueue of the same target should be allowed")
	}
}

func TestMarkDone_refusesAnUnknownTask(t *testing.T) {
	t.Run("TSQUT-E01: Marking an unknown task done is refused", func(t *testing.T) {})
	root := t.TempDir()
	if err := MarkDone(root, "ghost"); err == nil || !strings.Contains(err.Error(), "task not found") {
		t.Fatalf("MarkDone of an unknown task = %v, want the refusal", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, DoneDir)); len(entries) != 0 {
		t.Errorf("nothing may reach the history: %v", entries)
	}
}

func TestDropRemovesFromQueue(t *testing.T) {
	t.Run("TSQUT-B08: Dropping deletes without history", func(t *testing.T) {})
	t.Run("TSQUT-E02: Dropping an unknown task is refused", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "plans/x.md")
	_, _ = Enqueue(root, task("1-doc-x", "plans/x.md", "doc", "triage"))
	if err := Drop(root, "1-doc-x"); err != nil {
		t.Fatal(err)
	}
	if n, _ := PendingCount(root); n != 0 {
		t.Fatalf("after drop the queue should be empty, got %d", n)
	}
	// drop creates no history (unlike done)
	if _, err := os.Stat(filepath.Join(root, DoneDir)); !os.IsNotExist(err) {
		t.Fatal("drop should not create .anchors/done/")
	}
	// dropping a task that does not exist → error
	if err := Drop(root, "nao-existe"); err == nil {
		t.Fatal("dropping a task that does not exist should fail")
	}
}

func TestReclaimReturnsClaimedToPending(t *testing.T) {
	t.Run("TSQUT-B09: Old and unstamped claims are returned to pending", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "A.spec.md")
	withTarget(t, root, "B.spec.md")
	_, _ = Enqueue(root, task("1-spec-a", "A.spec.md", "spec", "implement"))
	_, _ = Enqueue(root, task("2-spec-b", "B.spec.md", "spec", "implement"))
	// Two workers take both — LONG AGO (outside the work window), which is the case of a
	// worker that died without closing.
	old := time.Now().Add(-2 * JanelaDeTrabalho).Format(time.RFC3339)
	_, _ = Claim(root, "dead-worker-1", old)
	_, _ = Claim(root, "dead-worker-2", old)
	n, err := Reclaim(root)
	if err != nil || n != 2 {
		t.Fatalf("reclaim: n=%d err=%v (want 2)", n, err)
	}
	// both are back to pending with no claimed_by
	tasks, _ := List(root)
	for _, tk := range tasks {
		if tk.State != Pending {
			t.Errorf("task %s should be pending, got %s", tk.ID, tk.State)
		}
		if tk.ClaimedBy != "" {
			t.Errorf("task %s should have claimed_by cleared, got %q", tk.ID, tk.ClaimedBy)
		}
	}
	// and they can be claimed again
	got, _ := Claim(root, "new-worker", "2026-08-13T00:00:00-03:00")
	if got == nil {
		t.Fatal("a returned task should be claimable")
	}
}

func TestSuggestNext(t *testing.T) {
	t.Run("TSQUT-B12: The next step follows the kind that changed", func(t *testing.T) {})
	cases := map[string]string{
		// The PROMOTED plan (already reviewed, in `plans/`) seeds work; the DRAFT
		// (`plans/review/`) goes to review first. They are two kinds because the state is the
		// FOLDER — otherwise editing the plan during execution would trigger review again
		// every time, until the step became noise.
		// The verbs are the ARTIFACTS of `anchors work` — see ArtefatosDeTrabalho. They used to
		// be `specify`/`implement`/`verify`, which `work` refuses.
		"plan":       "spec",
		"plan-draft": "review-plan-draft", "spec": "code", "feature": "test",
		// `test` closes the unit — and that is where the work LOOKS done. The chain does not
		// end in verifying: it calls the REVIEW. Measured in three rounds of a real E2E, 7
		// serious defects passed with every gate green; none was found by a gate.
		"code": "feature", "test": "review", "guide": "review",
		"mistério": "triage",
	}
	for kind, want := range cases {
		got, reason := SuggestNext(kind)
		if got != want {
			t.Errorf("SuggestNext(%q) = %q, want %q", kind, got, want)
		}
		if reason == "" {
			t.Errorf("SuggestNext(%q) gives no reason", kind)
		}
	}
}

// `reclaim` RESPECTS whoever took the task recently. Without it, it returns everything —
// including what an ACTIVE worker is doing — and two agents start writing the same file
// without knowing. It happened, measured: a subagent ran 90 minutes on one step, looked stuck
// from outside, someone reclaimed, and the work was duplicated.
func TestReclaimRespectsARecentClaim(t *testing.T) {
	t.Run("TSQUT-X01: A recent claim is not reclaimed", func(t *testing.T) {})
	t.Run("TSQUT-B10: Forced reclaiming returns even recent claims", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "A.spec.md")
	_, _ = Enqueue(root, task("1-spec-a", "A.spec.md", "spec", "implement"))

	now := time.Now().Format(time.RFC3339)
	if _, err := Claim(root, "worker-ativo", now); err != nil {
		t.Fatal(err)
	}

	n, err := Reclaim(root)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("reclaim returned %d task(s) of an ACTIVE worker — that is how two agents end up "+
			"in the same file", n)
	}

	// `--force` is the explicit way out for whoever KNOWS the worker stopped.
	if n, err := ReclaimForce(root); err != nil || n != 1 {
		t.Fatalf("reclaim --force: n=%d err=%v (want 1)", n, err)
	}
}

// A task with no stamp of when it was claimed has nothing to respect — it is returned.
func TestReclaimWithoutStampReturns(t *testing.T) {
	t.Run("TSQUT-B09: Old and unstamped claims are returned to pending", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "A.spec.md")
	_, _ = Enqueue(root, task("1-spec-a", "A.spec.md", "spec", "implement"))
	if _, err := Claim(root, "worker-without-stamp", ""); err != nil {
		t.Fatal(err)
	}
	if n, _ := Reclaim(root); n != 1 {
		t.Fatalf("a claim with no stamp should be returned, got %d", n)
	} // ClaimIsOld answers the same test the reclaim uses.
	if !ClaimIsOld(Task{}) || !ClaimIsOld(Task{ClaimedAt: "not a time"}) {
		t.Error("a claim with no stamp, or an unreadable one, is old")
	}
	if ClaimIsOld(Task{ClaimedAt: time.Now().Format(time.RFC3339)}) {
		t.Error("a claim taken now is not old")
	}
	if !ClaimIsOld(Task{ClaimedAt: time.Now().Add(-2 * JanelaDeTrabalho).Format(time.RFC3339)}) {
		t.Error("a claim older than the window is old")
	}
}

// Locks the divergence that broke routing in a real E2E: the queue suggested
// `specify`/`implement`/`verify`, and `anchors work` refuses all three ("unknown artifact").
// Whoever pulled the task could not compose the prompt and had to translate the verbs on
// their own — the only point where the cycle did not route itself.
func TestQueueSuggestionIsComposableByWork(t *testing.T) {
	t.Run("TSQUT-B13: The suggestions of the unit kinds are composable by the work command", func(t *testing.T) {})
	// every mapped kind must suggest a verb `work` accepts; only an unmapped kind gets
	// `triage`, the queue's marker for "decide by hand". `guide` used to suggest
	// `review-governed`, which `work` refuses.
	for _, k := range []string{"plan-draft", "plan", "spec", "feature", "code", "test", "guide"} {
		verb, why := SuggestNext(k)
		if verb == "" {
			continue // a kind with no next step is legitimate
		}
		if !ValidWorkArtifact(verb) {
			t.Errorf("kind %q suggests %q, which `anchors work` refuses — whoever pulls the task "+
				"cannot compose the prompt", k, verb)
		}
		if why == "" {
			t.Errorf("kind %q suggests %q without saying why", k, verb)
		}
	}
}

// Guards the noise measured in four real runs: the watcher enqueues on CHANGE and never
// dequeues on DELETION, so a reviewer's probe (`__probe.test.tsx`), deleted right after, left a
// task alive forever.
func TestTaskWithAMissingTargetIsDiscarded(t *testing.T) {
	t.Run("TSQUT-B04: A task whose target was deleted is removed", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "vive.ts")
	if _, err := Enqueue(root, task("1-code-vive", "vive.ts", "code", "feature")); err != nil {
		t.Fatal(err)
	}
	if _, err := Enqueue(root, task("2-code-sonda", "sonda.test.ts", "code", "review")); err != nil {
		t.Fatal(err)
	}

	tasks, err := List(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].Changed != "vive.ts" {
		t.Fatalf("the task of the deleted target should be gone; got %+v", tasks)
	}
	// And gone from DISK, not only from the listing: filtering without removing would make the
	// ghost reappear in the next `list`.
	left, _ := os.ReadDir(filepath.Join(root, ".anchors", "tasks"))
	for _, e := range left {
		if strings.Contains(e.Name(), "sonda") {
			t.Error("the ghost task is still on disk — it would reappear in the next listing")
		}
	}
}

// The ZERO of reclaim needs to explain itself: `anchors reclaim` answered "0 task(s) returned"
// with a task visibly `claimed` in the queue. The number was RIGHT — it was claimed minutes ago,
// within JanelaDeTrabalho — and the zero alone looks like a defect.
func TestRecentlyHeld(t *testing.T) {
	t.Run("TSQUT-B11: Recently held claims are counted", func(t *testing.T) {})
	root := t.TempDir()
	// the TARGET must exist: List discards a task of a deleted file.
	if err := os.WriteFile(filepath.Join(root, "x.md"), []byte("# x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	now := time.Now().Format(time.RFC3339)
	old := time.Now().Add(-2 * JanelaDeTrabalho).Format(time.RFC3339)

	// The `claimed__` files are written directly: `Claim` takes the FIRST pending one per call,
	// so building three distinct states through it would depend on the scan order.
	for _, c := range []struct{ id, when string }{
		{"recent-1", now},
		{"recent-2", now},
		{"old", old},
	} {
		task := Task{
			ID: c.id, Changed: "x.md", Kind: "plan",
			State: Claimed, ClaimedBy: "worker", ClaimedAt: c.when,
		}
		data, err := yaml.Marshal(task)
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(dirFor(root), fileName(Claimed, c.id))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if got := RecentlyHeld(root); got != 2 {
		t.Errorf("RecentlyHeld = %d, want 2 (the two from now)", got)
	}

	// Reclaim takes only the one past the window
	n, err := Reclaim(root)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("Reclaim returned %d, want 1", n)
	}
	if got := RecentlyHeld(root); got != 2 {
		t.Errorf("after Reclaim, RecentlyHeld = %d, want 2", got)
	}

	// and force takes both
	if n, _ := ReclaimForce(root); n != 2 {
		t.Errorf("ReclaimForce returned %d, want 2", n)
	}
	if got := RecentlyHeld(root); got != 0 {
		t.Errorf("after force, RecentlyHeld = %d, want 0", got)
	}
}

func TestPendingCount_countsPendingAndClaimed(t *testing.T) {
	t.Run("TSQUT-B14: The pending count counts pending and claimed tasks", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "a.md")
	withTarget(t, root, "b.md")
	_, _ = Enqueue(root, task("1-a", "a.md", "doc", "triage"))
	_, _ = Enqueue(root, task("2-b", "b.md", "doc", "triage"))
	_, _ = Claim(root, "w", time.Now().Format(time.RFC3339))
	if n, err := PendingCount(root); err != nil || n != 2 {
		t.Fatalf("PendingCount = %d, %v; want 2", n, err)
	}
}

func TestList_surfacesACorruptedTaskFile(t *testing.T) {
	t.Run("TSQUT-E03: A corrupted task file is listed as triage, and can be claimed and dropped", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "a.md")
	_, _ = Enqueue(root, task("1-a", "a.md", "doc", "triage"))
	if err := os.WriteFile(filepath.Join(dirFor(root), "pending__0-bad.yaml"), []byte("id: [unclosed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Before: it was skipped — never listed, never claimed, never cleaned.
	tasks, err := List(root)
	if err != nil || len(tasks) != 2 || tasks[0].ID != "0-bad" || tasks[1].ID != "1-a" {
		t.Fatalf("List = %+v, %v; want 0-bad and 1-a, no error", tasks, err)
	}
	bad := tasks[0]
	if bad.State != Pending || bad.SuggestedNext != "triage" || !strings.Contains(bad.Reason, "pending__0-bad.yaml") {
		t.Fatalf("the corrupted file must surface as a pending triage task naming its file, got %+v", bad)
	}
	c, err := Claim(root, "w", "2026-08-07T00:00:00Z")
	if err != nil || c == nil || c.ID != "0-bad" {
		t.Fatalf("Claim = %+v, %v; want the corrupted task", c, err)
	}
	if err := Drop(root, "0-bad"); err != nil {
		t.Fatalf("Drop of the corrupted task: %v", err)
	}
	if tasks, _ := List(root); len(tasks) != 1 || tasks[0].ID != "1-a" {
		t.Fatalf("after the drop only 1-a remains, got %+v", tasks)
	}
}

func TestList_keepsATaskWithAnAbsoluteTarget(t *testing.T) {
	t.Run("TSQUT-B15: A task whose target is an absolute path that exists is kept", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "src/a.go")
	abs := filepath.Join(root, "src", "a.go")
	if created, err := Enqueue(root, task("1-code-a", abs, "code", "feature")); err != nil || !created {
		t.Fatalf("enqueue: %v %v", created, err)
	}
	// Before: the absolute target was joined under the root, never found, and the task was
	// deleted as a ghost by the very next List.
	for i := 0; i < 2; i++ {
		tasks, err := List(root)
		if err != nil || len(tasks) != 1 || tasks[0].Changed != abs {
			t.Fatalf("List #%d = %+v, %v; want the task on %s", i+1, tasks, err, abs)
		}
	}
	if _, err := os.Stat(filepath.Join(dirFor(root), "pending__1-code-a.yaml")); err != nil {
		t.Fatalf("the task file must stay on disk: %v", err)
	}
}

func TestEnqueue_refusesAnIDHeldByAnotherTarget(t *testing.T) {
	t.Run("TSQUT-E04: Enqueuing an ID a live task already holds for another target is refused", func(t *testing.T) {})
	root := t.TempDir()
	withTarget(t, root, "a.go")
	withTarget(t, root, "b.go")
	if _, err := Enqueue(root, task("1-x", "a.go", "code", "feature")); err != nil {
		t.Fatal(err)
	}
	created, err := Enqueue(root, task("1-x", "b.go", "code", "feature"))
	if err == nil || created {
		t.Fatalf("Enqueue over a live ID = %v, %v; want a refusal", created, err)
	}
	if !strings.Contains(err.Error(), "1-x") {
		t.Errorf("the refusal must name the ID: %v", err)
	}
	tasks, _ := List(root)
	if len(tasks) != 1 || tasks[0].Changed != "a.go" {
		t.Fatalf("the first task must survive untouched, got %+v", tasks)
	}
	// A claimed holder counts too.
	if _, err := Claim(root, "w", "2026-08-07T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if created, err := Enqueue(root, task("1-x", "b.go", "code", "feature")); err == nil || created {
		t.Fatalf("Enqueue over a claimed ID = %v, %v; want a refusal", created, err)
	}
}

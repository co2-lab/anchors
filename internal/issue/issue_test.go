package issue

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func viol() Issue {
	return Issue{
		Kind: Violation, Target: "features/x/A.spec.md", Gate: "spec-sections",
		Detail: "falta a seção Regras", Date: "2026-08-07",
	}
}

func TestKeyIsStableAcrossDates(t *testing.T) {
	t.Run("ISLFS-B01: The key is stable across dates and distinct per gate", func(t *testing.T) {})
	a := viol()
	b := viol()
	b.Date = "2027-01-01" // another date
	if a.Key() != b.Key() {
		t.Fatalf("the Key must be stable in time: %q vs %q", a.Key(), b.Key())
	}
	if a.ID() == b.ID() {
		t.Fatal("the ID (file name) must vary with the date")
	}
	// a different gate → a different Key (the same target may violate two gates)
	c := viol()
	c.Gate = "spec-has-code"
	if a.Key() == c.Key() {
		t.Fatal("distinct gates should give distinct Keys")
	}
	if got := (Issue{Kind: Stale, Anchor: "a/x.spec.md", Target: "a/x.go"}).Key(); got != "stale--a-x.spec.md--vs--a-x.go" {
		t.Errorf("an anchored edge key = %q, want stale--a-x.spec.md--vs--a-x.go", got)
	}
}

func TestIDIsLegibleAndSanitized(t *testing.T) {
	t.Run("ISLFS-B02: The file name is the date and the key, with no slash", func(t *testing.T) {})
	id := viol().ID()
	if !strings.HasPrefix(id, "2026-08-07--violation--") || !strings.HasSuffix(id, ".md") {
		t.Fatalf("unexpected ID: %s", id)
	}
	if strings.Contains(id, "/") {
		t.Fatalf("the ID must not hold a slash: %s", id)
	}
}

func TestOpenCreatesInTodo(t *testing.T) {
	t.Run("ISLFS-B04: A new issue is opened in todo", func(t *testing.T) {})
	t.Run("ISLFS-B03: The body names the kind, the target, the gate and the detail", func(t *testing.T) {})
	root := t.TempDir()
	created, at, err := Open(root, viol())
	if err != nil || !created || at != Todo {
		t.Fatalf("open: created=%v at=%v err=%v", created, at, err)
	}
	ids, _ := List(root, Todo)
	if len(ids) != 1 {
		t.Fatalf("want 1 issue in todo/, got %v", ids)
	}
	body, _ := os.ReadFile(filepath.Join(root, Dir, string(Todo), ids[0]))
	s := string(body)
	for _, want := range []string{"# VIOLATION: features/x/A.spec.md", "**gate:** spec-sections", "**owner:** agente",
		"**detected on:** 2026-08-07", "## Violated invariant\n\nfalta a seção Regras"} {
		if !strings.Contains(s, want) {
			t.Errorf("the body lacks %q:\n%s", want, s)
		}
	}
}

func TestOpenIsIdempotentAcrossDates(t *testing.T) {
	t.Run("ISLFS-B05: The same issue is not opened twice", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, viol())
	// confronting again on ANOTHER day does not duplicate (same Key)
	later := viol()
	later.Date = "2026-09-15"
	created, at, _ := Open(root, later)
	if created {
		t.Fatal("must not recreate an existing issue (same Key, another date)")
	}
	if at != Todo {
		t.Fatalf("the current state should be todo, got %v", at)
	}
	if ids, _ := List(root, Todo); len(ids) != 1 {
		t.Fatalf("want 1 issue, got %d", len(ids))
	}
}

func TestResolveMovesToDone(t *testing.T) {
	t.Run("ISLFS-B08: Resolving moves a live issue to done, and only once", func(t *testing.T) {})
	t.Run("ISLFS-I01: A resolved issue is moved, not copied, and never resurrected", func(t *testing.T) {})
	root := t.TempDir()
	_, _, _ = Open(root, viol())
	// the confrontation passes again → resolve
	ok, err := Resolve(root, viol().Key())
	if err != nil || !ok {
		t.Fatalf("resolve: ok=%v err=%v", ok, err)
	}
	if ids, _ := List(root, Todo); len(ids) != 0 {
		t.Fatalf("todo/ should be empty after resolve, got %v", ids)
	}
	if ids, _ := List(root, Done); len(ids) != 1 {
		t.Fatalf("done/ should hold 1, got %v", ids)
	}
	// resolving again is a no-op
	ok2, _ := Resolve(root, viol().Key())
	if ok2 {
		t.Fatal("resolving an already resolved issue should be a no-op")
	}
	// and opening again does not resurrect it
	if created, at, _ := Open(root, viol()); created || at != Done {
		t.Fatalf("opening a resolved issue = %v, %q; want nothing created, done", created, at)
	}
}

func TestResolveFromDoing(t *testing.T) {
	t.Run("ISLFS-B08: Resolving moves a live issue to done, and only once", func(t *testing.T) {})
	root := t.TempDir()
	// simulates an issue in doing/ (someone took it)
	i := viol()
	doingDir := filepath.Join(root, Dir, string(Doing))
	_ = os.MkdirAll(doingDir, 0o755)
	_ = os.WriteFile(filepath.Join(doingDir, i.ID()), []byte(i.Body()), 0o644)
	// the check passes → resolve even from doing/
	ok, _ := Resolve(root, i.Key())
	if !ok {
		t.Fatal("should resolve an issue that was in doing/")
	}
	if ids, _ := List(root, Done); len(ids) != 1 {
		t.Fatalf("done/ should hold 1, got %v", ids)
	}
}

func TestOpenDoesNotResurrectResolvedIssue(t *testing.T) {
	t.Run("ISLFS-B05: The same issue is not opened twice", func(t *testing.T) {})
	root := t.TempDir()
	i := viol()
	// already resolved (in done/)
	doneDir := filepath.Join(root, Dir, string(Done))
	_ = os.MkdirAll(doneDir, 0o755)
	_ = os.WriteFile(filepath.Join(doneDir, i.ID()), []byte("resolved"), 0o644)
	created, at, _ := Open(root, i)
	if created {
		t.Fatal("must not reopen an issue already in done/")
	}
	if at != Done {
		t.Fatalf("should report done, got %v", at)
	}
	if todos, _ := List(root, Todo); len(todos) != 0 {
		t.Fatalf("todo/ should stay empty, got %v", todos)
	}
}

// The lifecycle of an assumed debt. `obligation_pending` used to be a line in a file header:
// visible only to whoever opened it, with no state, no way to be paid, no way to fall due.
func TestDebtIsBornInFutureAndClosesWhenPaid(t *testing.T) {
	t.Run("ISLFS-B06: An assumed debt is born in future and shows when it is due", func(t *testing.T) {})
	root := t.TempDir()
	iss := Issue{
		Kind: Violation, Target: "infra/models/X.spec.md", Gate: "obligation-honored",
		Detail: "não alcançada pelo purge", Date: "2026-08-13",
		Prazo: "`lgpd-eliminacao` — na etapa de código da Fase 1",
	}

	created, at, err := OpenAt(root, iss, Future)
	if err != nil || !created {
		t.Fatalf("the debt should be born: created=%v err=%v", created, err)
	}
	if at != Future {
		t.Errorf("a debt is born in `future/`, not in %q — whoever reads `todo/` asks "+
			"'what do I do NOW'", at)
	}

	// Confronting again neither duplicates nor promotes it to `todo/`: it stays deferred.
	if created, at, _ := OpenAt(root, iss, Future); created || at != Future {
		t.Errorf("confronting again must not duplicate; got created=%v at=%q", created, at)
	}

	// Paying the debt (the gate passes again) closes the issue — the same cycle as the others.
	ok, err := Resolve(root, iss.Key())
	if err != nil || !ok {
		t.Fatalf("a paid debt must close: ok=%v err=%v", ok, err)
	}
	if st, exists := Exists(root, iss.Key()); !exists || st != Done {
		t.Errorf("once paid the debt lives in `done/`; got %q (exists=%v)", st, exists)
	}

	// The due moment must be in the body: a debt with no visible due date is a TODO with a
	// better name.
	body := iss.Body()
	for _, want := range []string{"When it will be paid", "na etapa de código da Fase 1", "ASSUMED debt"} {
		if !strings.Contains(body, want) {
			t.Errorf("the debt body must hold %q", want)
		}
	}
}

func TestDecisionBodySaysHowToCloseIt(t *testing.T) {
	t.Run("ISLFS-B07: A decision explains how to close it", func(t *testing.T) {})
	body := Issue{Kind: Decision, Target: "b.spec.md", Gate: "open-questions-resolved", Detail: "which cache?", Date: "2026-08-29"}.Body()
	for _, want := range []string{"## The open question", "**How to close it:**", "PROMOTE the answer to a rule", "**What NOT to do:**"} {
		if !strings.Contains(body, want) {
			t.Errorf("the decision body lacks %q:\n%s", want, body)
		}
	}
}

// Guards the defect measured in an E2E: a second `anchors judge --verdict fail` on the SAME
// unit only printed "issue already recorded" and discarded the `--reason`. Two distinct reports
// were lost like that. Idempotence is the right policy for the SAME problem detected twice;
// it is not for a DIFFERENT problem in the same place.
func TestNewFindingDoesNotVanishInSilence(t *testing.T) {
	t.Run("ISLFS-B09: A new finding reopens the issue and keeps the old report", func(t *testing.T) {})
	root := t.TempDir()
	base := Issue{Kind: Violation, Target: "a.ts", Gate: "review", Date: "2026-08-13"}

	first := base
	first.Detail = "## Laudo A\nperda de dado na leitura"
	if created, _, err := Open(root, first); err != nil || !created {
		t.Fatalf("the first issue should be born: %v", err)
	}
	if _, err := Resolve(root, base.Key()); err != nil {
		t.Fatal(err)
	}

	// a DIFFERENT finding on the same target: appended, keeping the previous one
	second := base
	second.Detail = "## Laudo B\ncontradição entre duas regras"
	reopened, err := Reopen(root, second)
	if err != nil || !reopened {
		t.Fatalf("a new finding should reopen: reopened=%v err=%v", reopened, err)
	}
	st, name, _ := byKey(root, base.Key())
	if st != Todo {
		t.Fatalf("the reopened issue must be in todo/, is in %s", st)
	}
	body, _ := os.ReadFile(pathFor(root, Todo, name))
	for _, want := range []string{"Laudo A", "Laudo B", "## Additional finding", "perda de dado", "contradição"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("the body must keep %q — old AND new report", want)
		}
	}

	// the SAME finding again: idempotent, no duplicate
	if reopened, _ := Reopen(root, second); reopened {
		t.Error("the same finding cannot be appended twice")
	}
	body2, _ := os.ReadFile(pathFor(root, Todo, name))
	if strings.Count(string(body2), "Laudo B") != 1 {
		t.Errorf("`Laudo B` must appear once; it appeared %d times", strings.Count(string(body2), "Laudo B"))
	}
}

// WHOSE issue it is is an axis INDEPENDENT of the kind: an agent that hits a question only a
// person answers must hand the issue over WITHOUT it ceasing to be the violation it was.
func TestOwnerFiltersAndReassigns(t *testing.T) {
	t.Run("ISLFS-B10: Issues are listed by owner, and no owner means the agent", func(t *testing.T) {})
	t.Run("ISLFS-B11: Reassigning hands the issue over with its reason", func(t *testing.T) {})
	root := t.TempDir()

	agents := Issue{Kind: Violation, Target: "src/a.ts", Gate: "layer-boundary",
		Detail: "imported what it should not", Date: "2026-08-29"}
	users := Issue{Kind: Decision, Target: "src/b.spec.md", Gate: "open-questions-resolved",
		Detail: "1 open decision", Date: "2026-08-29", Dono: DonoUsuário}
	for _, i := range []Issue{agents, users} {
		if _, _, err := Open(root, i); err != nil {
			t.Fatal(err)
		}
	}

	// The FILTER is what makes the list usable: the decider's list cannot come mixed with the
	// agent's work, or both stop being read.
	forUser, err := ListByOwner(root, Todo, DonoUsuário)
	if err != nil {
		t.Fatal(err)
	}
	if len(forUser) != 1 || !strings.Contains(forUser[0], "b.spec.md") {
		t.Fatalf("the user filter should bring only the decision, got %v", forUser)
	}
	forAgent, _ := ListByOwner(root, Todo, DonoAgente)
	if len(forAgent) != 1 || !strings.Contains(forAgent[0], "a.ts") {
		t.Fatalf("the agent filter should bring only the violation, got %v", forAgent)
	}

	// REASSIGN: the agent tried, hit a wall, and hands it over. The issue is still the violation
	// it was — the owner changes, not the kind.
	if err := Reassign(root, Todo, forAgent[0], DonoUsuário,
		"the boundary depends on which layer owns the cache, and that is not decided"); err != nil {
		t.Fatal(err)
	}
	after, _ := ListByOwner(root, Todo, DonoUsuário)
	if len(after) != 2 {
		t.Errorf("both should be the user's now, got %d", len(after))
	}
	if left, _ := ListByOwner(root, Todo, DonoAgente); len(left) != 0 {
		t.Errorf("the agent should have nothing left, got %v", left)
	}
	// The WHY goes along: whoever receives an issue without context asks what was tried.
	b, _ := os.ReadFile(filepath.Join(root, Dir, string(Todo), forAgent[0]))
	if !strings.Contains(string(b), "owns the cache") {
		t.Error("the reason for the handover should be recorded in the issue")
	}
	if !strings.Contains(string(b), "violation") {
		t.Error("the owner changed, not the kind — it is still the violation it was")
	}
}

// An issue written BEFORE the field existed is the agent's, which was the only case. Reading it
// as the user's would fill the decider's list with work that is not theirs.
func TestMissingOwnerIsTheAgents(t *testing.T) {
	t.Run("ISLFS-B10: Issues are listed by owner, and no owner means the agent", func(t *testing.T) {})
	root := t.TempDir()
	dir := filepath.Join(root, Dir, string(Todo))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "2026-01-01--violation--x--y.md")
	if err := os.WriteFile(old, []byte("# VIOLATION: y\n\n- **kind:** violation\n- **alvo (regido):** y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if d := FileOwner(old); d != DonoAgente {
		t.Errorf("an issue without the field is the agent's, got %q", d)
	}
}

// `Reassign` must find its anchor in an OLD issue — the one with the Portuguese label. The
// fallback branch ONLY runs on an issue older than the `owner` field, which was written by the
// old binary with `- **alvo`. And the failure is SILENT: the replace does not match, the owner
// line does not go in, and Reassign ends with no error and the issue without an owner.
func TestReassignFindsTheAnchorInAnIssueWithTheOldLabel(t *testing.T) {
	t.Run("ISLFS-B11: Reassigning hands the issue over with its reason", func(t *testing.T) {})
	root := t.TempDir()
	dir := filepath.Join(root, Dir, string(Todo))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// An issue as the OLD binary wrote it: Portuguese label and NO owner field.
	old := "# algo quebrou\n\n- **alvo (regido):** `src/a.ts`\n- **detectada em:** 2026-01-01\n"
	name := "0001-algo.md"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Reassign(root, Todo, name, DonoUsuário, "only a person decides this"); err != nil {
		t.Fatalf("reassign failed: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, "- **owner:**") {
		t.Fatalf("the owner line was not inserted into an issue with the old label — the reassign "+
			"passed silently and the issue kept no owner:\n%s", text)
	}
	// And the owner read back must be the one asked for: writing the line and not being able to
	// read it back would be the same defect one step later.
	if got := Owner(strings.TrimSpace(string(ownerRE.FindSubmatch(b)[1]))); got != DonoUsuário {
		t.Errorf("owner read back = %q, want %q", got, DonoUsuário)
	}
}

// A full check closes every open violation it did not reproduce — and nothing else. The renamed
// gate is the measured case: 40 `header-conforme` issues nothing could close once the gate
// became `header-conforms`.
func TestReconcileViolations(t *testing.T) {
	t.Run("ISLFS-B12: A full check closes the violations it no longer reproduces", func(t *testing.T) {})
	t.Run("ISLFS-X01: Reconciling spares decisions and the user's violations", func(t *testing.T) {})
	UseFiles()
	root := t.TempDir()
	alive := Issue{Kind: Violation, Gate: "header-conforms", Target: "a.ts", Date: "2026-09-25"}
	renamed := Issue{Kind: Violation, Gate: "header-conforme", Target: "a.ts", Date: "2026-09-11"}
	users := Issue{Kind: Violation, Gate: "spec-complete", Target: "b.spec.md", Date: "2026-09-11", Dono: DonoUsuário}
	decision := Issue{Kind: Decision, Gate: "open-questions-resolved", Target: "c.spec.md", Date: "2026-09-11"}
	for _, i := range []Issue{alive, renamed, users, decision} {
		if _, _, err := Open(root, i); err != nil {
			t.Fatal(err)
		}
	}
	// One more, already being worked on.
	inProgress := Issue{Kind: Violation, Gate: "feature-test-match", Target: "d.feature", Date: "2026-09-11"}
	if _, _, err := OpenAt(root, inProgress, Doing); err != nil {
		t.Fatal(err)
	}

	closed, err := ReconcileViolations(root, map[string]bool{alive.Key(): true})
	if err != nil {
		t.Fatal(err)
	}
	if len(closed) != 2 {
		t.Fatalf("the renamed-gate and the doing violation should close, closed %v", closed)
	}
	for _, i := range []Issue{renamed, inProgress} {
		if st, _ := Exists(root, i.Key()); st != Done {
			t.Errorf("%s should be done, is %s", i.Key(), st)
		}
	}
	for _, i := range []Issue{alive, users, decision} {
		if st, _ := Exists(root, i.Key()); st != Todo {
			t.Errorf("%s must stay open (reproduced, the user's, or not a violation), is %s", i.Key(), st)
		}
	}
}

// The DEFAULT backend is files. A local project, or any caller that configured nothing, cannot
// end up talking to the network without asking.
func TestDefaultBackendIsFiles(t *testing.T) {
	t.Run("ISLFS-B13: Issues are files unless GitHub is configured", func(t *testing.T) {})
	UseFiles()
	if target != nil {
		t.Fatal("without configuration, the backend must be files")
	}
	UseGitHub("acme/x", "anchors")
	if target == nil || target.Repo != "acme/x" {
		t.Fatal("UseGitHub should route to the declared repository")
	}
	UseFiles() // does not leak into the other tests of the package
	if target != nil {
		t.Fatal("UseFiles should route back to files")
	}
}

func TestOpen_reportsAStateFolderItCannotCreate(t *testing.T) {
	t.Run("ISLFS-E01: An issue that cannot be written is reported", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, Dir), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if created, _, err := Open(root, viol()); err == nil || created {
		t.Fatalf("Open where issues/ is a file = %v, %v; want the error and nothing created", created, err)
	}
}

func TestReassign_reportsAMissingIssue(t *testing.T) {
	t.Run("ISLFS-E02: Reassigning a missing issue is reported", func(t *testing.T) {})
	if err := Reassign(t.TempDir(), Todo, "nope.md", DonoUsuário, "x"); err == nil {
		t.Fatal("reassigning a missing issue must fail")
	}
}

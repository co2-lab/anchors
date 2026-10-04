// @anchors
//   code: GTTSC
//   ref: GHIGT

package issue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

	"github.com/co2-lab/anchors/internal/initx"
)

// DEDUPLICATION on github comes from the marker in the BODY, and the confirmation must be exact.
//
// GitHub's search is textual and returns near misses: without checking the marker, a finding
// about `Foo.spec.md` would match the card of `FooBar.spec.md` — and Anchors would close the
// wrong card, which is worse than closing none.
func TestMarkerIsExactAndNotAPrefix(t *testing.T) {
	t.Run("GHIGT-B01: A card whose marker only starts with the key is not the issue's", func(t *testing.T) {})
	otherBody := fmt.Sprintf(KeyMarker, "unit-complete:packages/Foo.spec.md:violation")
	wanted := fmt.Sprintf(KeyMarker, "unit-complete:packages/Foo.spec.md")

	// The other card's marker CONTAINS the wanted prefix, and still must not match: they are
	// findings of different targets.
	if strings.Contains(otherBody, wanted) {
		t.Error("the marker must close with `-->`, or a finding matches the card of another " +
			"whose target merely STARTS the same")
	}
}

// The title names the card without depending on the key's format — tying it to the key would
// make it unreadable in the issue list, which is where it shows.
func TestTitleSaysWhatItIsWithoutTheKey(t *testing.T) {
	t.Run("GHIGT-B02: The title names the gate, the kind and the target", func(t *testing.T) {})
	g := GitHub{Repo: "acme/x", Label: "anchors"}
	for kind, want := range map[Kind]string{
		Violation: "Violação",
		Decision:  "Decisão",
		Stale:     "Desatualizado",
		Conflict:  "Conflito",
	} {
		got := g.title(Issue{Kind: kind, Gate: "unit-complete", Target: "a/b.spec.md"})
		if got != "[unit-complete] "+want+" @ a/b.spec.md" {
			t.Errorf("the title of %s = %q, want the gate, %q and the target", kind, got, want)
		}
	}
}

// fakeGH puts on the PATH a `gh` that records each call (one line per call) and answers
// `issue list` with the given JSON. `failOn` is a shell glob: a call that matches it
// exits 1 with an error message. Nothing reaches the network.
func fakeGH(t *testing.T, listJSON, failOn string) (calls func() []string) {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.txt")
	script := "#!/bin/sh\n" +
		"printf '%s' \"$*\" | tr '\\n' ' ' >> '" + log + "'\nprintf '\\n' >> '" + log + "'\n" +
		"case \"$*\" in\n"
	if failOn != "" {
		script += failOn + ") echo 'HTTP 502' >&2; exit 1 ;;\n"
	}
	script += "'issue list'*) cat <<'__EOF__'\n" + listJSON + "\n__EOF__\n;;\nesac\nexit 0\n"
	testkit.FakeBin(t, dir, "gh", script)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return func() []string {
		b, _ := os.ReadFile(log)
		return strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	}
}

func cardJSON(number int, state, key string) string {
	body, _ := json.Marshal("some text\n\n" + fmt.Sprintf(KeyMarker, key) + "\n")
	return fmt.Sprintf(`{"number":%d,"state":%q,"body":%s}`, number, state, body)
}

var ghFinding = Issue{Kind: Violation, Gate: "unit-complete", Target: "a/Foo.spec.md", Detail: "no feature", Date: "2026-09-26"}

func TestGitHubOpen_createsCardWithLabels(t *testing.T) {
	t.Run("GHIGT-B01: A card whose marker only starts with the key is not the issue's", func(t *testing.T) {})
	t.Run("GHIGT-B03: The search covers every state in the configured repository", func(t *testing.T) {})
	t.Run("GHIGT-B04: A new card carries the marker and the labels", func(t *testing.T) {})
	key := ghFinding.Key()
	// The search returns a near miss (a longer key that CONTAINS ours): it must not match.
	calls := fakeGH(t, "["+cardJSON(7, "OPEN", key+"-other")+"]", "")
	g := GitHub{Repo: "acme/x", Label: "anchors"}

	i := ghFinding
	i.Dono = DonoUsuário
	created, at, err := g.Open(i, Todo)
	if err != nil || !created || at != Todo {
		t.Fatalf("Open = %v, %q, %v; want a new card in todo", created, at, err)
	}
	got := calls()
	if len(got) != 2 {
		t.Fatalf("want a search and a create, got %q", got)
	}
	if !strings.HasPrefix(got[0], "issue list --state all --limit 500 --search "+key) || !strings.HasSuffix(got[0], "--repo acme/x") {
		t.Errorf("search call = %q", got[0])
	}
	for _, frag := range []string{
		"issue create --title [unit-complete] Violação @ a/Foo.spec.md",
		fmt.Sprintf(KeyMarker, key),
		"--label anchors --label " + initx.LabelToDo + " --label " + initx.LabelNeedsUser + " --repo acme/x",
	} {
		if !strings.Contains(got[1], frag) {
			t.Errorf("create call lacks %q:\n%s", frag, got[1])
		}
	}
}

func TestGitHubOpen_futureDebtHasNoFlowLabel(t *testing.T) {
	t.Run("GHIGT-B05: An assumed debt's card has no flow label", func(t *testing.T) {})
	calls := fakeGH(t, "[]", "")
	created, at, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Future)
	if err != nil || !created || at != Future {
		t.Fatalf("Open = %v, %q, %v; want a new card in future", created, at, err)
	}
	create := calls()[1]
	if strings.Contains(create, initx.LabelToDo) || strings.Contains(create, initx.LabelNeedsUser) {
		t.Errorf("a future debt must carry only the workflow label:\n%s", create)
	}
}

func TestGitHubOpen_existingOpenCardIsLeftAlone(t *testing.T) {
	t.Run("GHIGT-B06: An open card is left alone", func(t *testing.T) {})
	calls := fakeGH(t, "["+cardJSON(12, "OPEN", ghFinding.Key())+"]", "")
	created, at, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Todo)
	if err != nil || created || at != Todo {
		t.Fatalf("Open = %v, %q, %v; want no new card", created, at, err)
	}
	if got := calls(); len(got) != 1 {
		t.Fatalf("only the search may run, got %q", got)
	}
}

func TestGitHubOpen_closedCardIsReopenedWithTheNewReport(t *testing.T) {
	t.Run("GHIGT-B07: A closed card is reopened with the new report", func(t *testing.T) {})
	calls := fakeGH(t, "["+cardJSON(12, "CLOSED", ghFinding.Key())+"]", "")
	created, at, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Todo)
	if err != nil || !created || at != Todo {
		t.Fatalf("Open = %v, %q, %v; want the card reopened", created, at, err)
	}
	got := calls()
	if len(got) != 3 || got[1] != "issue reopen 12 --repo acme/x" ||
		!strings.HasPrefix(got[2], "issue comment 12 --body ⟲ **O achado voltou.**") || !strings.Contains(got[2], "no feature") {
		t.Fatalf("want reopen + comment with the detail, got %q", got)
	}
}

func TestGitHubOpen_propagatesGhFailures(t *testing.T) {
	t.Run("GHIGT-E01: A gh failure is reported with gh's output", func(t *testing.T) {})
	for name, tc := range map[string]struct{ list, failOn string }{
		"search fails":  {"[]", "'issue list'*"},
		"bad search":    {"not json", ""},
		"create fails":  {"[]", "'issue create'*"},
		"reopen fails":  {"[" + cardJSON(3, "CLOSED", ghFinding.Key()) + "]", "'issue reopen'*"},
		"comment fails": {"[" + cardJSON(3, "CLOSED", ghFinding.Key()) + "]", "'issue comment'*"},
	} {
		t.Run(name, func(t *testing.T) {
			fakeGH(t, tc.list, tc.failOn)
			created, _, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Todo)
			if err == nil || created {
				t.Fatalf("Open = %v, %v; want the failure reported", created, err)
			}
			if tc.failOn != "" && !strings.Contains(err.Error(), "HTTP 502") {
				t.Errorf("the error must carry gh's output, got %v", err)
			}
		})
	}
}

func TestGitHubResolve(t *testing.T) {
	t.Run("GHIGT-B08: Resolving closes only an open card", func(t *testing.T) {})
	t.Run("GHIGT-E01: A gh failure is reported with gh's output", func(t *testing.T) {})
	g := GitHub{Repo: "acme/x", Label: "anchors"}
	key := ghFinding.Key()

	calls := fakeGH(t, "["+cardJSON(5, "OPEN", key)+"]", "")
	if ok, err := g.Resolve(key); !ok || err != nil {
		t.Fatalf("Resolve of an open card = %v, %v; want closed", ok, err)
	}
	if got := calls(); len(got) != 2 || !strings.HasPrefix(got[1], "issue close 5 --comment ✓ Resolvido") {
		t.Fatalf("want the card closed with a comment, got %q", got)
	}

	for name, list := range map[string]string{
		"already closed": "[" + cardJSON(5, "CLOSED", key) + "]",
		"not found":      "[]",
	} {
		calls = fakeGH(t, list, "")
		if ok, err := g.Resolve(key); ok || err != nil {
			t.Errorf("%s: Resolve = %v, %v; want nothing done", name, ok, err)
		}
		if got := calls(); len(got) != 1 {
			t.Errorf("%s: only the search may run, got %q", name, got)
		}
	}

	fakeGH(t, "["+cardJSON(5, "OPEN", key)+"]", "'issue close'*")
	if ok, err := g.Resolve(key); ok || err == nil {
		t.Fatalf("a failed close = %v, %v; want the error", ok, err)
	}
}

// EVERY LABEL THIS PACKAGE APPLIES MUST BE ONE `anchors init` CREATES.
//
// `gh` refuses the WHOLE command over a missing label:
//
//	could not add label: 'anchors:precisa-do-usuario' not found
//
// So a divergent literal does not degrade — it erases the record. Measured: the gate
// `open-questions-resolved` failed a spec with two open questions (the right behaviour) and the
// issue was NOT created. A gate that accuses and does not record is worse than no gate.
func TestAppliedLabelsAreTheOnesInitCreates(t *testing.T) {
	t.Run("GHIGT-B09: The labels applied are the ones init creates", func(t *testing.T) {})
	created := map[string]bool{initx.LabelNeedsUser: true}
	for _, s := range initx.WorkStates {
		created[s] = true
	}

	// The labels `GitHub.Open` applies besides the project's root label (that one comes from the
	// configuration and `init` creates it from there) — read from a real create call.
	calls := fakeGH(t, "[]", "")
	i := ghFinding
	i.Dono = DonoUsuário
	if _, _, err := (GitHub{Repo: "acme/x", Label: "anchors"}).Open(i, Todo); err != nil {
		t.Fatal(err)
	}
	args := strings.Fields(calls()[1])
	var applied []string
	for k := 0; k+1 < len(args); k++ {
		if args[k] == "--label" && args[k+1] != "anchors" {
			applied = append(applied, args[k+1])
		}
	}
	if len(applied) != 2 {
		t.Fatalf("want the to-do and needs-user labels applied, got %v", applied)
	}
	for _, l := range applied {
		if !created[l] {
			t.Errorf("the label %q is applied when creating a card and `init` does not create it — "+
				"`gh` will refuse the whole command and the gate's finding will not be recorded", l)
		}
	}
}

// The LEGACY name must not be applied again: it does not exist in the repositories, and it is
// exactly the literal that erased the record.
func TestLegacyLabelIsNotApplied(t *testing.T) {
	t.Run("GHIGT-B09: The labels applied are the ones init creates", func(t *testing.T) {})
	if initx.LabelNeedsUser == initx.LabelNeedsUserLegacy {
		t.Fatal("the canonical and the legacy became equal — the compatibility bridge lost its point")
	}
	if strings.Contains(initx.LabelNeedsUser, "precisa") {
		t.Errorf("LabelNeedsUser = %q — back to the Portuguese name", initx.LabelNeedsUser)
	}
}

// The PROSE THAT TEACHES must not cite the legacy name either. The rule above protects what the
// code APPLIES; this one protects what the documentation TEACHES — which is how the defect came
// back: a guide told to run `gh issue list --label anchors:sob-44`, which no repository
// initialised by the current version has. Whoever followed it got an EMPTY list, and an empty
// list does not look like an error, it looks like no work.
//
// The `*Legacy`/`*Antigo` constants still exist and are still right — they READ old cards. What
// this test forbids is TEACHING the old name to whoever creates new work.
func TestDocsDoNotTeachTheLegacyLabel(t *testing.T) {
	legacy := map[string]string{
		initx.PrefixoLabelSobAntigo: initx.PrefixoLabelSob,
		initx.LabelNeedsUserLegacy:  initx.LabelNeedsUser,
	}

	root := repoRoot(t)
	docs, err := filepath.Glob(filepath.Join(root, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) == 0 {
		t.Fatal("no .md at the root — the glob broke and the test would pass empty")
	}

	for _, doc := range docs {
		b, err := os.ReadFile(doc)
		if err != nil {
			t.Fatal(err)
		}
		for old, current := range legacy {
			// The block that EXPLAINS the migration may cite the old name — it is what it is
			// about. What it may not do is teach the old one as if it were current, and the
			// difference is the mention of the new one on the same line.
			for _, line := range strings.Split(string(b), "\n") {
				if strings.Contains(line, old) && !strings.Contains(line, current) {
					t.Errorf("%s teaches the legacy label %q (the current one is %q):\n\t%s",
						filepath.Base(doc), old, current, strings.TrimSpace(line))
				}
			}
		}
	}
}

// repoRoot climbs to the `go.mod`: the test runs with the PACKAGE folder as the current one,
// and the documents are at the root.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("did not find go.mod climbing from the package")
	return ""
}

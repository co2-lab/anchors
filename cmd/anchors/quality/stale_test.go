package quality

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// englishOutput pins the catalog for the test: the output goes through i18n, and the
// language must not depend on what a previous test left in the global state.
func englishOutput(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { _ = i18n.Set(i18n.Default) })
	if err := i18n.Set("en"); err != nil {
		t.Fatal(err)
	}
}

func TestStaleSeparatesNeverValidatedFromDrift(t *testing.T) {
	englishOutput(t)
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "s2"},
			{ID: "a.go", Kind: mapx.KindCode, Rev: "c1"},
			{ID: "b.go", Kind: mapx.KindCode, Rev: "b1"},
			// a test whose evidence expired: its dependency a.go advanced since ingestion
			{ID: "a_test.go", Kind: mapx.KindTest, Rev: "t1", Signal: &mapx.TestSignal{
				AtRev: "t1", ClosureRev: map[string]string{"a.go": "c0"},
			}},
		},
		Edges: []mapx.Edge{
			// validated at s1, the spec is at s2 now: drift
			{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies, Stamp: &mapx.Stamp{ValidatedFromRev: "s1", ValidatedToRev: "c1"}},
			// never confronted
			{From: "a.spec.md", To: "b.go", Type: mapx.EdgeSpecifies},
			// validated at the current revs: not stale
			{From: "a_test.go", To: "a.go", Type: mapx.EdgeTestedBy, Stamp: &mapx.Stamp{ValidatedFromRev: "t1", ValidatedToRev: "c1"}},
		},
	}
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, g)

	out, err := runQ(t, newStaleCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"1 EXPIRED test evidence(s)",
		"  a_test.go\n      1 dependencie(s) changed, e.g.: a.go\n",
		"2 of 3 stale edge(s):",
		"  a.spec.md ──specifies──▶ a.go  (advanced rev)\n",
		"  a.spec.md ──specifies──▶ b.go  (never validated)\n",
		"1 never validated, 1 with rev drift",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "a_test.go ──") {
		t.Errorf("an edge stamped at the current revs is not stale:\n%s", out)
	}
	if strings.Contains(out, "(and the test itself)") {
		t.Errorf("the test file did not change; only its dependency did:\n%s", out)
	}
}

func TestStaleWithEverythingValidated(t *testing.T) {
	englishOutput(t)
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "a.spec.md", Rev: "s1"}, {ID: "a.go", Rev: "c1"},
			// own file changed since ingestion, no closure recorded
			{ID: "a_test.go", Kind: mapx.KindTest, Rev: "t2", Signal: &mapx.TestSignal{AtRev: "t1"}}},
		Edges: []mapx.Edge{{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies, Stamp: &mapx.Stamp{ValidatedFromRev: "s1", ValidatedToRev: "c1"}}},
	}
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, g)

	out, err := runQ(t, newStaleCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "  a_test.go\n      the test file itself changed\n") {
		t.Errorf("a test whose own file changed has expired evidence:\n%s", out)
	}
	if !strings.Contains(out, "no stale edges (1 edges, all validated)") {
		t.Errorf("every edge validated at the current revs:\n%s", out)
	}
}

func TestStaleWithoutMapFails(t *testing.T) {
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, nil)
	if _, err := runQ(t, newStaleCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: must point at map build; got %v", err)
	}
}

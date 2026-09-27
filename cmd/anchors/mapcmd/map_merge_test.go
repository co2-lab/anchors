package mapcmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// saveGraph writes g to dir/name and returns the path.
func saveGraph(t *testing.T, dir, name string, g *mapx.Graph) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := mapx.Save(g, p); err != nil {
		t.Fatal(err)
	}
	return p
}

// runMerge runs the merge driver with the three paths git passes, and returns what it
// wrote to stderr.
func runMerge(t *testing.T, base, ours, theirs string) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	prev := os.Stderr
	os.Stderr = w
	cmd := newMapMergeCmd()
	cmd.SetArgs([]string{base, ours, theirs})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SilenceUsage = true
	runErr := cmd.Execute()
	os.Stderr = prev
	w.Close()
	b, _ := io.ReadAll(r)
	return string(b), runErr
}

func judgedEdge(from, to string, gates ...string) mapx.Edge {
	e := mapx.Edge{Type: "specifies", From: from, To: to}
	for _, g := range gates {
		e.Julgamentos = append(e.Julgamentos, mapx.Judgment{Gate: g, Verdict: "ok"})
	}
	return e
}

// git treats `anchors.graph.yaml` as TEXT, and it is derived. Measured in the reference app
// (co2-lab/anchors#12): a `git merge origin/develop` merged the map WITHOUT A CONFLICT and
// erased 62 judgment stamps. A stamp is state of the work, and the REPORT lives in the
// `--reason` of `anchors judge`, not in the file.
func TestMapMerge_unitesTheStampsOfBothSides(t *testing.T) {
	t.Run("MPMRM-B01: The merged map is written onto our side's file", func(t *testing.T) {})
	t.Run("MPMRM-B04: An edge created only on the other branch arrives with its judgment", func(t *testing.T) {})
	dir := t.TempDir()

	// both sides share edge A (judged on both) and each has an exclusive one
	base := saveGraph(t, dir, "base.yaml", &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts")}})
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{Edges: []mapx.Edge{
		judgedEdge("a.spec.md", "a.ts", "review"),
		judgedEdge("b.spec.md", "b.ts", "rule-fulfilled"), // ours only
	}})
	theirs := saveGraph(t, dir, "theirs.yaml", &mapx.Graph{Edges: []mapx.Edge{
		judgedEdge("a.spec.md", "a.ts", "review"),
		judgedEdge("c.spec.md", "c.ts", "review"), // theirs only
	}})

	if _, err := runMerge(t, base, ours, theirs); err != nil {
		t.Fatalf("map merge: %v", err)
	}

	// the result goes to `ours`, which is what git expects
	g, err := mapx.Load(ours)
	if err != nil {
		t.Fatal(err)
	}
	by := map[string][]string{}
	for _, e := range g.Edges {
		for _, j := range e.Julgamentos {
			by[e.From] = append(by[e.From], j.Gate)
		}
	}
	for _, c := range []struct{ from, gate string }{
		{"a.spec.md", "review"},         // both
		{"b.spec.md", "rule-fulfilled"}, // ours only — must not vanish
		{"c.spec.md", "review"},         // theirs only — must arrive
	} {
		found := false
		for _, g := range by[c.from] {
			if g == c.gate {
				found = true
			}
		}
		if !found {
			t.Errorf("the %q judgment of %s vanished in the union: %v", c.gate, c.from, by)
		}
	}
}

// Union, not replacement: the side that has the stamp must NOT lose it to a side that has
// none. That is the merge that erased the 62 — an empty side would win.
func TestMapMerge_anEmptySideDoesNotEraseTheOther(t *testing.T) {
	t.Run("MPMRM-B05: A side with no judgment on a shared edge does not erase the other side's", func(t *testing.T) {})
	dir := t.TempDir()
	stamped := &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts", "review")}}
	unstamped := &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts")}}

	for _, c := range []struct {
		name         string
		ours, theirs *mapx.Graph
	}{
		{"the stamp is on our side", stamped, unstamped},
		{"the stamp is on the other side", unstamped, stamped},
	} {
		t.Run(c.name, func(t *testing.T) {
			base := saveGraph(t, dir, "base-"+c.name+".yaml", unstamped)
			ours := saveGraph(t, dir, "ours-"+c.name+".yaml", c.ours)
			theirs := saveGraph(t, dir, "theirs-"+c.name+".yaml", c.theirs)
			if _, err := runMerge(t, base, ours, theirs); err != nil {
				t.Fatal(err)
			}
			g, err := mapx.Load(ours)
			if err != nil {
				t.Fatal(err)
			}
			if n := countJudgments(g); n != 1 {
				t.Errorf("result with %d judgment(s), want 1 — the empty side erased the other", n)
			}
		})
	}
}

// THE DRIVER UNITED EDGES AND STAMPS, AND NEVER THE NODES. Measured in the reference app:
// base 329 nodes, ours 330, theirs 332, result 330 — and git reports "Automatic merge went
// well". The file stays governed and is not in the map, so no gate confronts it.
func TestMapMerge_keepsTheNodesOnlyTheOtherSideHas(t *testing.T) {
	t.Run("MPMRM-B02: A node created only on the other branch reaches the merged map", func(t *testing.T) {})
	t.Run("MPMRM-I01: Nothing of either side is missing from the merged map", func(t *testing.T) {})
	dir := t.TempDir()
	node := func(id string) mapx.Node { return mapx.Node{ID: id, Kind: "spec", Rev: "r1"} }

	base := saveGraph(t, dir, "base.yaml", &mapx.Graph{Nodes: []mapx.Node{node("common.spec.md")}})
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{
		Nodes: []mapx.Node{node("common.spec.md"), node("ours-only.spec.md")},
		Edges: []mapx.Edge{judgedEdge("ours-only.spec.md", "o.ts")},
	})
	theirs := saveGraph(t, dir, "theirs.yaml", &mapx.Graph{
		Nodes: []mapx.Node{node("common.spec.md"), node("theirs-only-1.spec.md"), node("theirs-only-2.spec.md")},
		Edges: []mapx.Edge{judgedEdge("theirs-only-1.spec.md", "t.ts")},
	})

	if _, err := runMerge(t, base, ours, theirs); err != nil {
		t.Fatalf("map merge: %v", err)
	}
	g, err := mapx.Load(ours)
	if err != nil {
		t.Fatal(err)
	}
	has := map[string]bool{}
	for _, n := range g.Nodes {
		has[n.ID] = true
	}
	for _, id := range []string{"common.spec.md", "ours-only.spec.md", "theirs-only-1.spec.md", "theirs-only-2.spec.md"} {
		if !has[id] {
			t.Errorf("node %q vanished in the union (the map kept %d nodes: %v)", id, len(g.Nodes), has)
		}
	}
	edges := map[string]bool{}
	for _, e := range g.Edges {
		edges[e.From] = true
	}
	if !edges["ours-only.spec.md"] || !edges["theirs-only-1.spec.md"] {
		t.Errorf("an edge of one side is missing from the result: %v", edges)
	}
}

func TestMapMerge_aCommonNodeKeepsOurRevision(t *testing.T) {
	t.Run("MPMRM-B03: A node both sides have keeps our side's revision", func(t *testing.T) {})
	dir := t.TempDir()
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: "spec", Rev: "r-ours"}}})
	theirs := saveGraph(t, dir, "theirs.yaml", &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: "spec", Rev: "r-theirs"}}})
	if _, err := runMerge(t, ours, ours, theirs); err != nil {
		t.Fatal(err)
	}
	g, err := mapx.Load(ours)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Nodes) != 1 || g.Nodes[0].Rev != "r-ours" {
		t.Errorf("the common node must keep our revision, got %+v", g.Nodes)
	}
}

func TestMapMerge_reportsOnStderr(t *testing.T) {
	t.Run("MPMRM-B06: The driver reports on the error stream what it kept and what came from the other side", func(t *testing.T) {})
	dir := t.TempDir()
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts", "review")}})
	theirs := saveGraph(t, dir, "theirs.yaml", &mapx.Graph{Edges: []mapx.Edge{judgedEdge("b.spec.md", "b.ts", "review")}})
	out, err := runMerge(t, ours, ours, theirs)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "2 judgment(s) preserved (1 came from the other side)") {
		t.Errorf("unexpected report:\n%s", out)
	}
	// Nothing new from the other side: no "came from" clause.
	again, _ := runMerge(t, ours, ours, theirs)
	if !strings.Contains(again, "2 judgment(s) preserved\n") {
		t.Errorf("a merge that brought nothing must not claim anything came from the other side:\n%s", again)
	}
}

func TestMapMerge_neverReadsTheBase(t *testing.T) {
	t.Run("MPMRM-X01: The base version is never read", func(t *testing.T) {})
	dir := t.TempDir()
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts")}})
	theirs := saveGraph(t, dir, "theirs.yaml", &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts", "review")}})
	if _, err := runMerge(t, filepath.Join(dir, "no-such-base.yaml"), ours, theirs); err != nil {
		t.Errorf("a missing base must not matter: %v", err)
	}
}

func TestMapMerge_failures(t *testing.T) {
	t.Run("MPMRM-E01: An unreadable side fails the merge naming the side", func(t *testing.T) {})
	t.Run("MPMRM-E02: A call with two paths is refused", func(t *testing.T) {})
	dir := t.TempDir()
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{Edges: []mapx.Edge{judgedEdge("a.spec.md", "a.ts", "review")}})
	before, _ := os.ReadFile(ours)

	if _, err := runMerge(t, ours, ours, filepath.Join(dir, "missing.yaml")); err == nil ||
		!strings.Contains(err.Error(), "read the other side") {
		t.Errorf("a missing other side: got %v", err)
	}
	if after, _ := os.ReadFile(ours); string(after) != string(before) {
		t.Error("a failed merge rewrote our side")
	}
	bad := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(bad, []byte(":::\n\tnot yaml"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runMerge(t, ours, bad, ours); err == nil || !strings.Contains(err.Error(), "read our side") {
		t.Errorf("an unparseable side of ours: got %v", err)
	}

	cmd := newMapMergeCmd()
	cmd.SetArgs([]string{ours, ours})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "accepts 3 arg") {
		t.Errorf("two arguments must be refused: %v", err)
	}
}

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

// Two branches can judge the SAME edge with DIFFERENT gates. The driver used to hand
// the edge the other side's whole judgment list, so our gate's verdict vanished.
func TestMapMerge_joinsTheJudgmentsOfDifferentGatesOnASharedEdge(t *testing.T) {
	t.Run("MPMRM-B07: Judgments of different gates on a shared edge are joined, and one gate judged on both sides keeps the latest", func(t *testing.T) {})
	dir := t.TempDir()
	edge := func(js ...mapx.Judgment) mapx.Edge {
		e := judgedEdge("a.spec.md", "a.ts")
		e.Julgamentos = js
		return e
	}
	oursEdge := edge(
		mapx.Judgment{Gate: "review", Verdict: "ok", ChangedAt: "2026-09-01"},
		mapx.Judgment{Gate: "atomic", Verdict: "issue", ChangedAt: "2026-09-01"},
		mapx.Judgment{Gate: "tie", Verdict: "ok", ChangedAt: "2026-09-05"},
	)
	oursEdge.Stamp = &mapx.Stamp{Verdict: "issue", ChangedAt: "2026-09-01"}
	theirsEdge := edge(
		mapx.Judgment{Gate: "rule-fulfilled", Verdict: "ok", ChangedAt: "2026-09-02"},
		mapx.Judgment{Gate: "atomic", Verdict: "ok", ChangedAt: "2026-09-10"},
		mapx.Judgment{Gate: "tie", Verdict: "issue", ChangedAt: "2026-09-05"},
	)
	theirsEdge.Stamp = &mapx.Stamp{Verdict: "ok", ChangedAt: "2026-09-10"}
	ours := saveGraph(t, dir, "ours.yaml", &mapx.Graph{Edges: []mapx.Edge{oursEdge}})
	theirs := saveGraph(t, dir, "theirs.yaml", &mapx.Graph{Edges: []mapx.Edge{theirsEdge}})
	if _, err := runMerge(t, ours, ours, theirs); err != nil {
		t.Fatal(err)
	}
	g, err := mapx.Load(ours)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Edges) != 1 {
		t.Fatalf("want the one shared edge, got %+v", g.Edges)
	}
	got := map[string]string{}
	for _, j := range g.Edges[0].Julgamentos {
		got[j.Gate] = j.Verdict
	}
	if st := g.Edges[0].Stamp; st == nil || st.Verdict != "ok" {
		t.Errorf("the edge's stamp = %+v, want the other side's newer one", st)
	}
	want := map[string]string{"review": "ok", "rule-fulfilled": "ok", "atomic": "ok", "tie": "ok"}
	if len(got) != len(want) {
		t.Errorf("judgments = %v, want %v", got, want)
	}
	for gate, v := range want {
		if got[gate] != v {
			t.Errorf("gate %s: verdict %q, want %q (judgments %v)", gate, got[gate], v, got)
		}
	}
}

// The flow graph and the observed failures are part of the map too. The driver kept
// only our side's: a flow built on the other branch, or a failure ingested there on a
// node both sides have, vanished in a conflict-free merge.
func TestMapMerge_keepsTheOtherSidesFlowAndFailures(t *testing.T) {
	t.Run("MPMRM-B08: The other side's flow states and transitions and the failures of a shared node reach the merged map", func(t *testing.T) {})
	dir := t.TempDir()
	st := func(code string) mapx.FlowState {
		return mapx.FlowState{Code: code, Title: code, Flow: "flows/w.flow.md"}
	}
	tr := mapx.FlowTransition{From: "WORKR-T01", To: "WORKR-T02", Flow: "flows/w.flow.md"}
	node := func(fs ...mapx.FailureSignal) mapx.Node {
		return mapx.Node{ID: "a.spec.md", Kind: "spec", Rev: "r1", Failures: fs}
	}

	for _, c := range []struct {
		name     string
		oursFlow *mapx.FlowGraph
	}{{"our side has no flow", nil}, {"our side has its own flow", &mapx.FlowGraph{States: []mapx.FlowState{st("WORKR-T01")}}}} {
		t.Run(c.name, func(t *testing.T) {
			ours := saveGraph(t, dir, "ours-"+c.name+".yaml", &mapx.Graph{
				Nodes: []mapx.Node{node(
					mapx.FailureSignal{Rule: "AAAAA-E01", Count: 1, Last: "2026-09-01"},
					mapx.FailureSignal{Rule: "AAAAA-E02", Count: 1, Last: "2026-09-01"},
				)},
				Flow: c.oursFlow,
			})
			theirs := saveGraph(t, dir, "theirs-"+c.name+".yaml", &mapx.Graph{
				Nodes: []mapx.Node{node(
					mapx.FailureSignal{Rule: "AAAAA-E02", Count: 7, Last: "2026-09-09"},
					mapx.FailureSignal{Rule: "AAAAA-E03", Count: 2, Last: "2026-09-02"},
				)},
				Flow: &mapx.FlowGraph{States: []mapx.FlowState{st("WORKR-T01"), st("WORKR-T02")}, Transitions: []mapx.FlowTransition{tr}},
			})
			if _, err := runMerge(t, ours, ours, theirs); err != nil {
				t.Fatal(err)
			}
			g, err := mapx.Load(ours)
			if err != nil {
				t.Fatal(err)
			}
			if g.Flow == nil || len(g.Flow.States) != 2 || len(g.Flow.Transitions) != 1 {
				t.Errorf("the flow of the other side did not arrive whole: %+v", g.Flow)
			}
			counts := map[string]int{}
			for _, f := range g.Nodes[0].Failures {
				counts[f.Rule] = f.Count
			}
			if want := map[string]int{"AAAAA-E01": 1, "AAAAA-E02": 7, "AAAAA-E03": 2}; len(counts) != 3 ||
				counts["AAAAA-E01"] != 1 || counts["AAAAA-E02"] != 7 || counts["AAAAA-E03"] != 2 {
				t.Errorf("failures of the shared node = %v, want %v", counts, want)
			}
		})
	}
}

func TestMergeCommonNodes_keptEvidence(t *testing.T) {
	t.Run("MPMRM-B08: The other side's flow states and transitions and the failures of a shared node reach the merged map", func(t *testing.T) {})
	ours := &mapx.Graph{Nodes: []mapx.Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1", EvidenceKept: []mapx.EvidenceKeep{{Reason: "ours"}}}, {ID: "c", Rev: "r2"}}}
	theirs := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a", Rev: "r1", EvidenceKept: []mapx.EvidenceKeep{{Reason: "theirs"}}},
		{ID: "b", Rev: "r1", EvidenceKept: []mapx.EvidenceKeep{{Reason: "theirs"}}},
		{ID: "c", Rev: "r1", EvidenceKept: []mapx.EvidenceKeep{{Reason: "other rev"}}},
	}}
	mergeCommonNodes(ours, theirs)
	if k := ours.Nodes[0].EvidenceKept; len(k) != 1 || k[0].Reason != "theirs" {
		t.Errorf("the other side's declaration comes when ours has none, got %+v", k)
	}
	if k := ours.Nodes[1].EvidenceKept; k[0].Reason != "ours" {
		t.Errorf("ours stays when it has one, got %+v", k)
	}
	if len(ours.Nodes[2].EvidenceKept) != 0 {
		t.Error("a declaration of another revision does not come")
	}
}

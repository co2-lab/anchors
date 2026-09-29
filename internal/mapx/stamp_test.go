package mapx

import (
	"fmt"
	"strings"
	"testing"
)

// graph: spec(rev a) ──specifies──▶ code(rev b), and code ──tested-by──▶ test(rev c)
func stampGraph() *Graph {
	return &Graph{
		Version: 1,
		Nodes: []Node{
			{ID: "A.spec.md", Kind: KindSpec, Rev: "a"},
			{ID: "A.tsx", Kind: KindCode, Rev: "b"},
			{ID: "A.test.tsx", Kind: KindTest, Rev: "c"},
		},
		Edges: []Edge{
			{From: "A.spec.md", To: "A.tsx", Type: EdgeSpecifies},
			{From: "A.tsx", To: "A.test.tsx", Type: EdgeTestedBy},
		},
	}
}

func TestStampOnlyEdgesWithBothEndsConfronted(t *testing.T) {
	t.Run("EDSTD-B01: Only relations with both ends confronted are stamped", func(t *testing.T) {})
	g := stampGraph()
	// only spec and code were confronted; test was not in this round
	n := g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "2026-08-07T00:00:00Z")
	if n != 1 {
		t.Fatalf("expected to stamp 1 edge (spec→code), got %d", n)
	}
	if g.Edges[0].Stamp == nil {
		t.Fatal("the spec→code edge should be stamped")
	}
	if g.Edges[1].Stamp != nil {
		t.Fatal("the code→test edge should NOT be stamped (test not confronted)")
	}
}

func TestStampVerdictIssueWhenEndFails(t *testing.T) {
	t.Run("EDSTD-B02: The verdict is issue when an end failed and ok when both passed", func(t *testing.T) {})
	g := stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md", Failed: true}, {ID: "A.tsx"}}, "t")
	if got := g.Edges[0].Stamp.Verdict; got != "issue" {
		t.Fatalf("a failed end → verdict issue, got %q", got)
	}
	g = stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx", Failed: true}}, "t")
	if got := g.Edges[0].Stamp.Verdict; got != "issue" {
		t.Fatalf("a failed target end → verdict issue, got %q", got)
	}
}

func TestStampVerdictOkWhenBothPass(t *testing.T) {
	t.Run("EDSTD-B02: The verdict is issue when an end failed and ok when both passed", func(t *testing.T) {})
	g := stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	if got := g.Edges[0].Stamp.Verdict; got != "ok" {
		t.Fatalf("both pass → verdict ok, got %q", got)
	}
}

func TestStaleBecomesValidatedAfterStamp(t *testing.T) {
	t.Run("EDSTD-B04: A stamp is fresh until an end moves", func(t *testing.T) {})
	g := stampGraph()
	// before: never validated → stale
	if !g.Stale(g.Edges[0]) {
		t.Fatal("a never-validated edge should be stale")
	}
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	// after the stamp with the current revs → not stale
	if g.Stale(g.Edges[0]) {
		t.Fatal("a freshly stamped edge should not be stale")
	}
}

func TestStaleReappearsWhenRevAdvances(t *testing.T) {
	t.Run("EDSTD-B04: A stamp is fresh until an end moves", func(t *testing.T) {})
	g := stampGraph()
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	// the spec is edited → its rev moves → the edge goes stale again
	g.Nodes[0].Rev = "a2"
	if !g.Stale(g.Edges[0]) {
		t.Fatal("after an end moved, the edge should be stale again")
	}
}

func TestStampEdgeSingle(t *testing.T) {
	t.Run("EDSTD-B07: One relation is stamped by its ends, a missing one answers false", func(t *testing.T) {})
	g := stampGraph()
	// stamps only the spec→code edge with verdict issue (an AI judgment)
	ok := g.StampEdge("A.spec.md", "A.tsx", "issue", "2026-08-07T00:00:00Z")
	if !ok {
		t.Fatal("StampEdge should find the edge A.spec.md→A.tsx")
	}
	if g.Edges[0].Stamp == nil || g.Edges[0].Stamp.Verdict != "issue" {
		t.Fatalf("the edge should carry an issue stamp, got %+v", g.Edges[0].Stamp)
	}
	// the other edge stays unstamped
	if g.Edges[1].Stamp != nil {
		t.Fatal("only the edge asked for should be stamped")
	}
	// missing edge → false
	if g.StampEdge("X", "Y", "ok", "t") {
		t.Fatal("StampEdge of a missing edge should return false")
	}
}

func TestStampEdgeGoesStaleOnRevChange(t *testing.T) {
	t.Run("EDSTD-B04: A stamp is fresh until an end moves", func(t *testing.T) {})
	g := stampGraph()
	g.StampEdge("A.spec.md", "A.tsx", "ok", "t")
	if g.Stale(g.Edges[0]) {
		t.Fatal("freshly stamped, it should not be stale")
	}
	g.Nodes[1].Rev = "b2" // the target (code) changed
	if !g.Stale(g.Edges[0]) {
		t.Fatal("the AI verdict should expire when the target changes")
	}
}

func TestStaleEdgesListing(t *testing.T) {
	t.Run("EDSTD-B13: The stale relations are listed", func(t *testing.T) {})
	g := stampGraph()
	if len(g.StaleEdges()) != 2 {
		t.Fatalf("at first both edges are stale, got %d", len(g.StaleEdges()))
	}
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "t")
	if got := len(g.StaleEdges()); got != 1 {
		t.Fatalf("after stamping 1, 1 stays stale, got %d", got)
	}
}

// THE STAMP RECORDS THE CHANGE, NOT THE VERIFICATION.
//
// `changed_at` answers "since when is this relation as it is". Confronting again and finding the
// same result is NOT a new fact — and recording each confrontation made the map change on its
// own: 26 lines per check run in the reference project, PR conflicts where two people ran the
// check, and a dirty `git status` all the time.
func TestStampDateDoesNotMoveWhenNothingChanged(t *testing.T) {
	t.Run("EDSTD-B05: The date moves only when revisions or verdict change", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1"}},
		Edges: []Edge{{From: "a", To: "b", Type: EdgeSpecifies}},
	}
	v := []NodeVerdict{{ID: "a"}, {ID: "b"}}

	g.StampEdges(v, "2026-08-30")
	// Days LATER, same rev and same verdict: nothing changed, and the date cannot move.
	g.StampEdges(v, "2026-09-15")
	if got := g.Edges[0].Stamp.ChangedAt; got != "2026-08-30" {
		t.Fatalf("with no change the date must stay where it was, got %q", got)
	}

	// The REV changes → the relation changed, and the date moves.
	g.Nodes[1].Rev = "r2"
	g.StampEdges(v, "2026-09-20")
	if got := g.Edges[0].Stamp.ChangedAt; got != "2026-09-20" {
		t.Fatalf("a new rev is a change and the date must move, got %q", got)
	}

	// The VERDICT changes → also a change, even with the same revs.
	g.StampEdges([]NodeVerdict{{ID: "a"}, {ID: "b", Failed: true}}, "2026-09-25")
	if got := g.Edges[0].Stamp.ChangedAt; got != "2026-09-25" {
		t.Fatalf("a new verdict is a change and the date must move, got %q", got)
	}
}

// THE STAMP CANNOT CHANGE ON ITS OWN.
//
// Two validations on the SAME DAY, with the same verdict and the same revs, must produce the same
// map. What changes the stamp is the verdict or the rev — which is when there is something to
// record.
func TestStampSameInTwoRuns(t *testing.T) {
	t.Run("EDSTD-I01: The same round on the same day gives the same stamps", func(t *testing.T) {})
	build := func() *Graph {
		return &Graph{
			Nodes: []Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1"}},
			Edges: []Edge{{From: "a", To: "b", Type: EdgeSpecifies}},
		}
	}
	g1, g2 := build(), build()
	v := []NodeVerdict{{ID: "a"}, {ID: "b"}}

	// Same DAY — the real case: two `check` runs in a row.
	g1.StampEdges(v, "2026-08-30")
	g2.StampEdges(v, "2026-08-30")

	if *g1.Edges[0].Stamp != *g2.Edges[0].Stamp {
		t.Fatalf("the same day should produce the same stamp: %+v vs %+v", *g1.Edges[0].Stamp, *g2.Edges[0].Stamp)
	}
	// And the field cannot be empty: losing the answer would trade one problem for another.
	if g1.Edges[0].Stamp.ChangedAt == "" {
		t.Error("the stamp still has to say WHEN — less precision is not erasure")
	}
}

// AN UNDATED STAMP cannot perpetuate the gap.
//
// Measured: when the field was renamed, the existing map carried `last_validated` and the new
// parser read an empty `changed_at`. As nothing else changed, the empty value was kept on every
// run — and the field simply vanished from the map (`omitempty`), for good.
func TestStampWithoutDateTakesToday(t *testing.T) {
	t.Run("EDSTD-B06: An undated stamp takes today's date", func(t *testing.T) {})
	t.Run("EDSTD-X01: The stamp carries the caller's date", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1"}},
		Edges: []Edge{{From: "a", To: "b", Type: EdgeSpecifies,
			// The inherited stamp: same revs and verdict, but NO date.
			Stamp: &Stamp{ValidatedFromRev: "r1", ValidatedToRev: "r1", Verdict: "ok"}}},
	}
	g.StampEdges([]NodeVerdict{{ID: "a"}, {ID: "b"}}, "2026-09-01")
	if got := g.Edges[0].Stamp.ChangedAt; got != "2026-09-01" {
		t.Fatalf("an undated stamp must take today's date, or the gap persists; got %q", got)
	}
}

// A WAIVER is a person's decision; the mechanical check does not re-validate it. One commit in
// the reference app turned 21 waived judgments into `ok` through this loop.
func TestStampKeepsAWaiver(t *testing.T) {
	t.Run("EDSTD-B03: A waived stamp is kept as it was", func(t *testing.T) {})
	g := stampGraph()
	waiver := &Stamp{ValidatedFromRev: "old", ValidatedToRev: "old", ChangedAt: "2026-09-01", Verdict: "waived"}
	g.Edges[0].Stamp = waiver
	g.StampEdges([]NodeVerdict{{ID: "A.spec.md"}, {ID: "A.tsx"}}, "2026-09-24")
	if got := g.Edges[0].Stamp; got.Verdict != "waived" || got.ChangedAt != "2026-09-01" || got.ValidatedFromRev != "old" {
		t.Fatalf("the waiver was rewritten by the mechanical stamp: %+v", got)
	}
	if !g.Stale(g.Edges[0]) {
		t.Error("the ends changed rev since the waiver, and the edge does not read stale — a person must decide again")
	}
}

// The judgment gate computes nothing: it is the stamp that says whether someone already read it.
// And the stamp carries the ends' revs, so the verdict expires if the target changes — then it
// is a question again, which is the judgment's anti-drift.
func TestJudgedBy(t *testing.T) {
	t.Run("EDSTD-B10: A judgment holds only at the revisions it was given and only for its gate", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Rev: "r1"},
			{ID: "a.test.ts", Rev: "r1"},
		},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}

	// before judging: nobody answered
	if _, ok := g.JudgedBy("a.spec.md", "my-gate"); ok {
		t.Error("with no stamp, JudgedBy should say there is no verdict")
	}

	// judged: the verdict holds
	if n := g.StampNodeByGate("a.spec.md", "ok", "now", "my-gate"); n != 1 {
		t.Fatalf("stamped %d edges, wanted 1", n)
	}
	v, ok := g.JudgedBy("a.spec.md", "my-gate")
	if !ok || v != "ok" {
		t.Errorf("verdict = %q (ok=%v), wanted ok/true", v, ok)
	}

	// ANOTHER gate's stamp does not answer for this one
	if _, ok := g.JudgedBy("a.spec.md", "other-gate"); ok {
		t.Error("one gate's stamp does NOT answer for another")
	}

	// the target changed afterwards → the verdict expires and is a question again
	g.Nodes[0].Rev = "r2"
	if _, ok := g.JudgedBy("a.spec.md", "my-gate"); ok {
		t.Error("target changed after the stamp: the verdict should have expired")
	}
}

// The `check` rewrites the Stamp on every round (it sums up the worst verdict of all the node's
// gates). The judgment must NOT die there — that is what erased 16 verdicts and kept the counter
// at 16.
func TestJudgmentSurvivesTheCheckStamp(t *testing.T) {
	t.Run("EDSTD-B11: A judgment survives the next check round", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Rev: "r1"},
			{ID: "a.test.ts", Rev: "r1"},
		},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}
	g.StampNodeByGate("a.spec.md", "ok", "t0", "my-gate")

	// the check runs and rewrites the Stamp of the two confronted ends
	g.StampEdges([]NodeVerdict{{ID: "a.spec.md"}, {ID: "a.test.ts"}}, "t1")

	v, ok := g.JudgedBy("a.spec.md", "my-gate")
	if !ok || v != "ok" {
		t.Errorf("the judgment should survive the check stamp: verdict=%q ok=%v", v, ok)
	}
}

// A judgment gate that declares `guide:` stamps its guide→target edge; without the record,
// JudgedBy never saw it and the next check asked the same judgment again.
func TestStampEdgeByGateRecordsTheJudgment(t *testing.T) {
	t.Run("EDSTD-B08: Stamping a relation by gate records the gate and a judgment", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "GUIDE.md", Rev: "g1"}, {ID: "a.spec.md", Rev: "s1"}},
		Edges: []Edge{{From: "GUIDE.md", To: "a.spec.md", Type: EdgeGoverns}},
	}
	if !g.StampEdgeByGate("GUIDE.md", "a.spec.md", "ok", "2026-09-26", "guide-judge") {
		t.Fatal("the edge exists and should be stamped")
	}
	if st := g.Edges[0].Stamp; st == nil || st.Gate != "guide-judge" || st.Verdict != "ok" {
		t.Errorf("the stamp should name the gate, got %+v", st)
	}
	if v, ok := g.JudgedBy("a.spec.md", "guide-judge"); !ok || v != "ok" {
		t.Errorf("the target should count as judged by the gate: verdict=%q ok=%v", v, ok)
	}
	// with no gate named, no judgment is recorded
	g2 := &Graph{Nodes: g.Nodes, Edges: []Edge{{From: "GUIDE.md", To: "a.spec.md", Type: EdgeGoverns}}}
	g2.StampEdge("GUIDE.md", "a.spec.md", "ok", "2026-09-26")
	if len(g2.Edges[0].Julgamentos) != 0 {
		t.Errorf("a stamp with no gate records no judgment, got %+v", g2.Edges[0].Julgamentos)
	}
}

// `anchors judge` judges ONE target, so no edge has both ends in a round: stamping the node is
// the only way its edges stop reading "never validated" (measured: a node with 42 edges, 0 stamped).
func TestStampNodeByGateStampsEveryTouchingEdge(t *testing.T) {
	t.Run("EDSTD-B09: Stamping a node stamps every relation touching it", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "a.spec.md", Rev: "s"}, {ID: "a.go", Rev: "c"}, {ID: "a_test.go", Rev: "t"}, {ID: "x", Rev: "x"}, {ID: "y", Rev: "y"}},
		Edges: []Edge{
			{From: "a.spec.md", To: "a.go", Type: EdgeSpecifies},
			{From: "a.go", To: "a_test.go", Type: EdgeTestedBy},
			{From: "x", To: "y", Type: EdgeDependsOn},
		},
	}
	if n := g.StampNodeByGate("a.go", "ok", "2026-09-26", "code-judge"); n != 2 {
		t.Fatalf("the incoming and the outgoing edge should be stamped, got %d", n)
	}
	for _, e := range g.Edges[:2] {
		if e.Stamp == nil || len(e.Julgamentos) != 1 || e.Julgamentos[0].Gate != "code-judge" {
			t.Errorf("edge %s→%s should carry the stamp and the code-judge judgment, got %+v / %+v", e.From, e.To, e.Stamp, e.Julgamentos)
		}
	}
	if e := g.Edges[2]; e.Stamp != nil || len(e.Julgamentos) != 0 {
		t.Errorf("an edge that does not touch the node must stay untouched, got %+v", e)
	}
}

// Re-judging and finding the same is not a new fact: without keeping the date, every
// `anchors judge` would rewrite the date of every judgment.
func TestJudgmentIsOnePerGateAndKeepsItsDate(t *testing.T) {
	t.Run("EDSTD-B12: One judgment per gate, its date kept when nothing changed", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "a", Rev: "r1"}, {ID: "b", Rev: "r1"}},
		Edges: []Edge{{From: "a", To: "b"}},
	}
	g.StampNodeByGate("a", "ok", "t0", "my-gate")
	g.StampNodeByGate("a", "ok", "t1", "my-gate")
	js := g.Edges[0].Julgamentos
	if len(js) != 1 || js[0].ChangedAt != "t0" {
		t.Fatalf("the same verdict at the same revs keeps one judgment dated t0, got %+v", js)
	}
	g.StampNodeByGate("a", "issue", "t2", "my-gate")
	js = g.Edges[0].Julgamentos
	if len(js) != 1 || js[0].Verdict != "issue" || js[0].ChangedAt != "t2" {
		t.Fatalf("a new verdict replaces the gate's judgment and dates it t2, got %+v", js)
	}
}

// Judging one node rewrote the stamp of every relation touching it, a waiver included:
// the check keeps a waiver (EDSTD-B03) and the judge did not. A waiver answers ONE gate's
// question; another gate's verdict does not answer it, the same gate judging again does.
func TestJudgingANodeKeepsAnotherGatesWaiver(t *testing.T) {
	t.Run("EDSTD-B14: Judging keeps a waiver another gate recorded, and the gate that waived replaces its own", func(t *testing.T) {})
	waived := func(g *Graph) {
		g.StampNodeByGate("A.tsx", "waived", "2026-01-01", "atomic")
	}

	// another gate on the node: the stamp stays waived, the new gate's judgment is recorded
	g := stampGraph()
	waived(g)
	g.StampNodeByGate("A.tsx", "ok", "2026-02-02", "review")
	for _, e := range g.Edges {
		if e.Stamp == nil || e.Stamp.Verdict != "waived" || e.Stamp.Gate != "atomic" || e.Stamp.ChangedAt != "2026-01-01" {
			t.Errorf("%s → %s: another gate rewrote the waiver: %+v", e.From, e.To, e.Stamp)
		}
	}
	if v, ok := g.JudgedBy("A.tsx", "review"); !ok || v != "ok" {
		t.Errorf("the other gate's judgment was not recorded: %q %v", v, ok)
	}
	// the same through the single-edge path (a gate with a guide)
	g.StampEdgeByGate("A.spec.md", "A.tsx", "issue", "2026-02-03", "rule-fulfilled")
	if st := g.Edges[0].Stamp; st.Verdict != "waived" || st.Gate != "atomic" {
		t.Errorf("judging one edge with another gate rewrote the waiver: %+v", st)
	}

	// the gate that waived judges again: its own waiver is replaced
	g = stampGraph()
	waived(g)
	g.StampNodeByGate("A.tsx", "ok", "2026-02-02", "atomic")
	for _, e := range g.Edges {
		if e.Stamp == nil || e.Stamp.Verdict != "ok" {
			t.Errorf("%s → %s: the gate that waived could not replace its waiver: %+v", e.From, e.To, e.Stamp)
		}
	}

	// a new waiver always lands
	g = stampGraph()
	waived(g)
	g.StampNodeByGate("A.tsx", "waived", "2026-02-02", "review")
	if st := g.Edges[0].Stamp; st.Verdict != "waived" || st.Gate != "review" {
		t.Errorf("a new waiver did not land: %+v", st)
	}
}

func TestApplyStampChanges(t *testing.T) {
	t.Run("EDSTD-B15: A round carries only the stamps it changed to the map on disk", func(t *testing.T) {})
	edges := func() []Edge {
		return []Edge{
			{From: "a", To: "b", Type: EdgeTestedBy},
			{From: "a", To: "c", Type: EdgeTestedBy, Stamp: &Stamp{Verdict: "ok", ChangedAt: "d0"}},
			{From: "a", To: "d", Type: EdgeTestedBy},
			{From: "a", To: "e", Type: EdgeTestedBy, Stamp: &Stamp{Verdict: "ok", ChangedAt: "d0"}},
		}
	}
	copyOf := &Graph{Edges: edges()}
	before := copyOf.EdgeStamps()
	copyOf.Edges[0].Stamp = &Stamp{Verdict: "ok", ChangedAt: "d1"}    // the round stamps a→b
	copyOf.Edges[1].Stamp = &Stamp{Verdict: "issue", ChangedAt: "d1"} // and a→c
	copyOf.Edges = append(copyOf.Edges, Edge{From: "a", To: "z", Type: EdgeTestedBy, Stamp: &Stamp{Verdict: "ok"}})

	disk := &Graph{Edges: edges()}
	disk.Edges[1].Stamp = &Stamp{Verdict: "waived", ChangedAt: "d1"} // a judge waived a→c meanwhile
	disk.Edges[2].Stamp = &Stamp{Verdict: "ok", ChangedAt: "d1"}     // and another check stamped a→d
	applied, kept := disk.ApplyStampChanges(copyOf, before)
	if applied != 1 || kept != 1 {
		t.Fatalf("one applied (a→b), one kept (a→c), got %d applied, %d kept", applied, kept)
	}
	if disk.Edges[0].Stamp == nil || disk.Edges[0].Stamp.ChangedAt != "d1" {
		t.Errorf("the round's own stamp reaches the map on disk, got %+v", disk.Edges[0].Stamp)
	}
	if disk.Edges[1].Stamp.Verdict != "waived" || disk.Edges[2].Stamp.Verdict != "ok" || disk.Edges[3].Stamp.ChangedAt != "d0" {
		t.Errorf("what the round did not change stays as the disk has it: %+v %+v %+v", disk.Edges[1].Stamp, disk.Edges[2].Stamp, disk.Edges[3].Stamp)
	}
	if len(disk.Edges) != 4 {
		t.Errorf("an edge the disk does not have is not added, got %d edges", len(disk.Edges))
	}
	disk.Edges[0].Stamp.Verdict = "mutated"
	if copyOf.Edges[0].Stamp.Verdict != "ok" || before[edgeKey(&copyOf.Edges[1])].Verdict != "ok" {
		t.Error("the applied stamp and the snapshot are copies, not shared pointers")
	}
}

func TestRebaseRev(t *testing.T) {
	t.Run("EDSTD-B16: A new revision that proves nothing new keeps what was measured", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Rev: "r1", Signal: &TestSignal{AtRev: "r1", MutationAtRev: "r1",
				ProvenRevBySuite: map[string]string{"u": "r1", "old": "r0"},
				CoverageBySuite:  map[string]SuiteCoverage{"u": {AtRev: "r1"}, "old": {AtRev: "r0"}},
				MutationByScope:  map[string]MutationScope{"isolated": {AtRev: "r1"}}}},
			{ID: "a_test.go", Rev: "t1", Signal: &TestSignal{ClosureRev: map[string]string{"a.spec.md": "r1"}}},
			{ID: "b.spec.md", Rev: "r1"},
		},
		Edges: []Edge{{From: "a.spec.md", To: "a_test.go", Stamp: &Stamp{ValidatedFromRev: "r1", ValidatedToRev: "t1"},
			Julgamentos: []Judgment{{Gate: "g", ValidatedFromRev: "r1", ValidatedToRev: "t1"}}}},
	}
	g.RebaseRev("a.spec.md", "r1", "r2")
	s := g.Nodes[0].Signal
	if g.Nodes[0].Rev != "r2" || s.AtRev != "r2" || s.MutationAtRev != "r2" || s.ProvenRevBySuite["u"] != "r2" ||
		s.CoverageBySuite["u"].AtRev != "r2" || s.MutationByScope["isolated"].AtRev != "r2" {
		t.Errorf("what was measured at r1 moves to r2, got %+v", g.Nodes[0])
	}
	if s.ProvenRevBySuite["old"] != "r0" || s.CoverageBySuite["old"].AtRev != "r0" {
		t.Error("what was measured at another revision stays")
	}
	if g.Nodes[1].Signal.ClosureRev["a.spec.md"] != "r2" || g.Nodes[2].Rev != "r1" {
		t.Error("the closure follows the file; another file with the same revision stays")
	}
	e := g.Edges[0]
	if e.Stamp.ValidatedFromRev != "r2" || e.Stamp.ValidatedToRev != "t1" || e.Julgamentos[0].ValidatedFromRev != "r2" || e.Julgamentos[0].ValidatedToRev != "t1" {
		t.Errorf("the edge's stamp and judgment follow the file's end only, got %+v", e)
	}
	g.RebaseRev("a.spec.md", "r9", "r3")
	if g.Nodes[0].Rev != "r2" {
		t.Error("a revision the file is not at moves nothing")
	}
}

func TestKeepEvidence(t *testing.T) {
	t.Run("EDSTD-B17: A declared change keeps what was proven, and the lines only when asked", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Rev: "r2", Signal: &TestSignal{AtRev: "r1",
				ProvenRevBySuite: map[string]string{"u": "r1", "e2e": "r0"}}},
			{ID: "a.ts", Rev: "c1", Signal: &TestSignal{AtRev: "c1", MutationAtRev: "c1", TotalLines: 9,
				CoverageBySuite: map[string]SuiteCoverage{"u": {AtRev: "c1"}},
				MutationByScope: map[string]MutationScope{"isolated": {AtRev: "c1"}}}},
			{ID: "a_test.go", Rev: "t1", Signal: &TestSignal{AtRev: "t1", ClosureRev: map[string]string{"a.spec.md": "r1", "a.ts": "c1"}}},
		},
		Edges: []Edge{
			{From: "a.spec.md", To: "a_test.go", Stamp: &Stamp{ValidatedFromRev: "r1", ValidatedToRev: "t1"},
				Julgamentos: []Judgment{{Gate: "g", ValidatedFromRev: "r0", ValidatedToRev: "t1"}}},
			{From: "a.ts", To: "a.spec.md", Stamp: &Stamp{ValidatedFromRev: "c1", ValidatedToRev: "r1"}},
		},
	}
	got := g.KeepEvidence("a.spec.md", "r3", "only @realizes added", "2026-09-29", false)
	if strings.Join(got, ",") != "r0,r1" {
		t.Errorf("every earlier revision is carried, sorted, got %v", got)
	}
	s := g.Nodes[0].Signal
	if g.Nodes[0].Rev != "r3" || s.AtRev != "r3" || s.ProvenRevBySuite["u"] != "r3" || s.ProvenRevBySuite["e2e"] != "r3" {
		t.Errorf("the proofs move to the current revision, got %+v", s)
	}
	if g.Nodes[2].Signal.ClosureRev["a.spec.md"] != "r3" || g.Edges[0].Stamp.ValidatedFromRev != "r3" ||
		g.Edges[0].Julgamentos[0].ValidatedFromRev != "r3" || g.Edges[1].Stamp.ValidatedToRev != "r3" {
		t.Errorf("closures, stamps and judgments follow, got %+v %+v", g.Nodes[2].Signal, g.Edges)
	}
	if k := g.Nodes[0].EvidenceKept; len(k) != 1 || k[0].From != "r0 r1" || k[0].To != "r3" || k[0].Reason != "only @realizes added" ||
		k[0].At != "2026-09-29" || k[0].Lines {
		t.Errorf("the declaration is recorded, got %+v", k)
	}

	g.KeepEvidence("a.ts", "c2", "a comment reworded", "d", false)
	c := g.Nodes[1].Signal
	if g.Nodes[1].Rev != "c2" || c.AtRev != "c1" || c.MutationAtRev != "c1" || c.CoverageBySuite["u"].AtRev != "c1" ||
		c.MutationByScope["isolated"].AtRev != "c1" {
		t.Errorf("without lines, coverage and mutation stay, got %+v", c)
	}
	if g.Nodes[2].Signal.ClosureRev["a.ts"] != "c2" || g.Edges[1].Stamp.ValidatedFromRev != "c2" {
		t.Errorf("without lines, the closure and stamps still move")
	}
	g.KeepEvidence("a.ts", "c3", "same lines", "d", true)
	if c.AtRev != "c3" || c.MutationAtRev != "c3" || c.CoverageBySuite["u"].AtRev != "c3" || c.MutationByScope["isolated"].AtRev != "c3" ||
		!g.Nodes[1].EvidenceKept[1].Lines {
		t.Errorf("with lines, coverage and mutation move too, got %+v", c)
	}

	if got := g.KeepEvidence("a.spec.md", "r3", "again", "d", false); got != nil || len(g.Nodes[0].EvidenceKept) != 1 {
		t.Errorf("nothing at an earlier revision carries and records nothing, got %v", got)
	}
	if g.KeepEvidence("nope.md", "x", "r", "d", false) != nil {
		t.Error("an unknown file is left alone")
	}
	bare := &Graph{Nodes: []Node{{ID: "b.ts", Rev: "b1"}, {ID: "b_test.go", Signal: &TestSignal{ClosureRev: map[string]string{"b.ts": "b0"}}}}}
	if got := bare.KeepEvidence("b.ts", "b2", "r", "d", false); len(got) != 1 || bare.Nodes[0].Signal != nil || len(bare.Nodes[0].EvidenceKept) != 1 {
		t.Errorf("a file with no signal records the declaration on the node and gets no signal, got %v %+v", got, bare.Nodes[0])
	}
	for i := 0; i < 7; i++ {
		g.KeepEvidence("a.spec.md", fmt.Sprintf("r%d", 10+i), fmt.Sprint(i), "d", false)
	}
	if k := g.Nodes[0].EvidenceKept; len(k) != 5 || k[4].Reason != "6" || k[0].Reason != "2" {
		t.Errorf("the latest five declarations are kept, got %+v", k)
	}
}

package mapx

import "testing"

func reviewedGraph() *Graph {
	return &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "r1"}}}
}

func TestReviewOf_holdsAtItsRevision(t *testing.T) {
	t.Run("MPRVM-B01: A review holds only at the revision it looked at", func(t *testing.T) {})
	g := reviewedGraph()
	g.RecordReview("a.spec.md", Review{Gate: "rule-fulfilled", By: "human:ana"})
	if _, ok := g.ReviewOf("a.spec.md", "rule-fulfilled"); !ok {
		t.Error("the review at the current revision holds")
	}
	if _, ok := g.ReviewOf("a.spec.md", "other"); ok {
		t.Error("another gate's review does not")
	}
	if _, ok := g.ReviewOf("nope", "rule-fulfilled"); ok {
		t.Error("a node not in the map has none")
	}
	g.Nodes[0].Rev = "r2"
	if _, ok := g.ReviewOf("a.spec.md", "rule-fulfilled"); ok {
		t.Error("a review of an earlier revision no longer holds")
	}
}

func TestRecordReview_replacesTheSameGate(t *testing.T) {
	t.Run("MPRVM-B02: Recording a review replaces the same gate's and keeps the others", func(t *testing.T) {})
	t.Run("MPRVM-I01: A node holds at most one review per gate", func(t *testing.T) {})
	g := reviewedGraph()
	g.RecordReview("a.spec.md", Review{Gate: "g1", By: "x"})
	g.RecordReview("a.spec.md", Review{Gate: "g2", By: "y"})
	g.Nodes[0].Rev = "r2"
	g.RecordReview("a.spec.md", Review{Gate: "g1", By: "z"})
	rs := g.Nodes[0].Reviews
	if len(rs) != 2 {
		t.Fatalf("one review per gate, got %+v", rs)
	}
	for _, r := range rs {
		if r.Gate == "g1" && (r.By != "z" || r.Rev != "r2") {
			t.Errorf("g1's review is the new one, at the current revision: %+v", r)
		}
		if r.Gate == "g2" && r.By != "y" {
			t.Errorf("g2's review is kept: %+v", r)
		}
	}
	if g.RecordReview("nope", Review{Gate: "g1"}) {
		t.Error("recording on a node not in the map records nothing")
	}
}

func TestPreserveStamps_keepsReviews(t *testing.T) {
	t.Run("MPRVM-B03: A rebuild keeps the reviews whatever the revision", func(t *testing.T) {})
	t.Run("MPRVM-X01: Who reviewed is kept as given", func(t *testing.T) {})
	old := reviewedGraph()
	old.RecordReview("a.spec.md", Review{Gate: "g", By: "agent:openai/gpt"})
	novo := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "r9"}}}
	PreserveStamps(novo, old)
	if len(novo.Nodes[0].Reviews) != 1 || novo.Nodes[0].Reviews[0].By != "agent:openai/gpt" {
		t.Errorf("the rebuilt node keeps the review as given: %+v", novo.Nodes[0].Reviews)
	}
}

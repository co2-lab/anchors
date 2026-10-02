package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func reviewCfg() *config.Config {
	return &config.Config{Gates: []config.Gate{
		{Name: "rule-fulfilled", On: []string{"spec"}, Measures: config.MeasuresJudgment, Ask: "q", Review: &config.Review{Ask: "look again"}},
		{Name: "plain", On: []string{"spec"}, Check: "non-empty"},
	}}
}

func reviewGraph() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "r1"}, {ID: "b.spec.md", Kind: mapx.KindSpec, Rev: "r1"}}}
}

func TestReviewsDue_targetsWithNoReviewAtTheirRevision(t *testing.T) {
	t.Run("RVDUR-B01: The targets of a reviewed gate with no review at their revision are to review", func(t *testing.T) {})
	t.Run("RVDUR-B03: A change to a reviewed target makes it to review again", func(t *testing.T) {})
	t.Run("RVDUR-X01: The list never blocks", func(t *testing.T) {})
	g := reviewGraph()
	g.RecordReview("a.spec.md", mapx.Review{Gate: "rule-fulfilled", By: "human:ana"})
	due := ReviewsDue(reviewCfg(), g, t.TempDir(), "")
	if len(due) != 1 || due[0].Target != "b.spec.md" || due[0].Ask != "look again" {
		t.Fatalf("only the unreviewed spec is due, with the question: %+v", due)
	}
	g.Nodes[0].Rev = "r2"
	if due := ReviewsDue(reviewCfg(), g, t.TempDir(), ""); len(due) != 2 {
		t.Errorf("a changed target is to review again: %+v", due)
	}
}

func TestReviewsDue_onlyReviewedGates(t *testing.T) {
	t.Run("RVDUR-B02: A gate with no review has nothing to review, and a named gate lists its own", func(t *testing.T) {})
	if due := ReviewsDue(reviewCfg(), reviewGraph(), t.TempDir(), "plain"); len(due) != 0 {
		t.Errorf("a gate with no review lists nothing: %+v", due)
	}
	if due := ReviewsDue(reviewCfg(), reviewGraph(), t.TempDir(), "rule-fulfilled"); len(due) != 2 {
		t.Errorf("the named reviewed gate lists its targets: %+v", due)
	}
	if ReviewsDue(nil, reviewGraph(), "", "") != nil || ReviewsDue(reviewCfg(), nil, "", "") != nil {
		t.Error("with no configuration or no map there is nothing to review")
	}
}

func TestReviewsDue_changesNoVerdict(t *testing.T) {
	t.Run("RVDUR-I01: Listing what is to review changes no verdict", func(t *testing.T) {})
	cfg := reviewCfg()
	g := reviewGraph()
	n := g.Nodes[0]
	before := runOne(cfg.Gates[0], n, t.TempDir(), g, cfg)
	g.RecordReview(n.ID, mapx.Review{Gate: "rule-fulfilled", By: "human:ana"})
	after := runOne(cfg.Gates[0], g.Nodes[0], t.TempDir(), g, cfg)
	if before.Verdict != Judge || after.Verdict != Judge {
		t.Errorf("the judgment is asked with or without the review: %v, %v", before.Verdict, after.Verdict)
	}
}

func TestReviewsDue_orderedByGateThenTarget(t *testing.T) {
	t.Run("RVDUR-B01: The targets of a reviewed gate with no review at their revision are to review", func(t *testing.T) {})
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "zeta", On: []string{"spec"}, Review: &config.Review{Ask: "z"}},
		{Name: "alpha", On: []string{"spec"}, Review: &config.Review{Ask: "a"}},
	}}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "b.spec.md", Kind: mapx.KindSpec, Rev: "r"}, {ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "r"}}}
	var got []string
	for _, d := range ReviewsDue(cfg, g, t.TempDir(), "") {
		got = append(got, d.Gate+":"+d.Target)
	}
	want := "alpha:a.spec.md,alpha:b.spec.md,zeta:a.spec.md,zeta:b.spec.md"
	if s := strings.Join(got, ","); s != want {
		t.Errorf("order = %s, want %s", s, want)
	}
}

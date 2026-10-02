package gate

import (
	"sort"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// DueReview is a target to review for a gate: no review recorded at its current revision.
type DueReview struct {
	Gate   string
	Target string
	// Ask is what the reviewer is asked to look at.
	Ask string
}

// ReviewsDue are the targets of every gate that declares `review:` with no review recorded
// at their current revision, by gate and target. Only the gate named in `only` when one is
// given.
//
// A review is apart from the gate's verdict: the gate still measures as it declares, and
// this list says what a reviewer has not looked at since it last changed. It is informed,
// never blocking — the project decides when the reviewer comes.
func ReviewsDue(cfg *config.Config, g *mapx.Graph, root, only string) []DueReview {
	if cfg == nil || g == nil {
		return nil
	}
	var out []DueReview
	for _, gt := range cfg.Gates {
		if gt.Review == nil || (only != "" && gt.Name != only) {
			continue
		}
		for _, n := range g.Nodes {
			if !applies(gt, n, root) {
				continue
			}
			if _, ok := g.ReviewOf(n.ID, gt.Name); ok {
				continue
			}
			out = append(out, DueReview{Gate: gt.Name, Target: n.ID, Ask: gt.ReviewAsk()})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Gate != out[j].Gate {
			return out[i].Gate < out[j].Gate
		}
		return out[i].Target < out[j].Target
	})
	return out
}

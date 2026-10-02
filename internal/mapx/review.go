package mapx

// ReviewOf is the review of a node for a gate that holds now: recorded at the node's
// current revision. A review of an earlier revision is history, and the node is to review
// again.
func (g *Graph) ReviewOf(id, gate string) (Review, bool) {
	for _, n := range g.Nodes {
		if n.ID != id {
			continue
		}
		for _, r := range n.Reviews {
			if r.Gate == gate && r.Rev == n.Rev {
				return r, true
			}
		}
	}
	return Review{}, false
}

// RecordReview records a review of a node, at its current revision, replacing the node's
// earlier review for the same gate. It says whether the node is in the map.
func (g *Graph) RecordReview(id string, r Review) bool {
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if n.ID != id {
			continue
		}
		r.Rev = n.Rev
		kept := n.Reviews[:0:0]
		for _, old := range n.Reviews {
			if old.Gate != r.Gate {
				kept = append(kept, old)
			}
		}
		n.Reviews = append(kept, r)
		return true
	}
	return false
}

// preserveReviews carries each node's reviews into the rebuilt map, whatever its revision:
// a review of an earlier revision is the history of who looked and when, and `ReviewOf`
// is what reads it as no longer holding.
func preserveReviews(novo, antigo *Graph) {
	byID := make(map[string][]Review, len(antigo.Nodes))
	for _, n := range antigo.Nodes {
		byID[n.ID] = n.Reviews
	}
	for i := range novo.Nodes {
		if rs, ok := byID[novo.Nodes[i].ID]; ok && len(novo.Nodes[i].Reviews) == 0 {
			novo.Nodes[i].Reviews = rs
		}
	}
}

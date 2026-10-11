// @anchors
//   code: FSMFR
//   ref: FRSCF

package mapx

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/scan"
)

// --- frozen parts: which freeze covers what (DESIGN-frozen-units.md) ---
//
// The scan reads the marks where they are written; the map carries them to everything they
// cover, because it is the one that knows who is whose: a unit's spec frozen covers the
// unit's code, feature and tests; a doctrine document or rule frozen covers the spec rules
// that realize it. The gates, the reports and the guides then ask the node, and never read
// a mark themselves.

// unitOfFile is the unit a file belongs to by declaration — its own header's `code:`, the
// anchor beside it, its `ref:` —, never the code inferred from its text: a file that cites
// a frozen unit in a string is not part of it.
func unitOfFile(f scan.File, anchors map[string]string) string {
	if f.Upstream {
		return ""
	}
	if f.HeaderCode != "" && len(f.HeaderRefs) == 0 {
		return f.HeaderCode
	}
	if c := anchors[f.Path]; c != "" {
		return c
	}
	if len(f.HeaderRefs) > 0 {
		return f.HeaderRefs[0]
	}
	return ""
}

// applyFreezes writes, on each node, the freeze that covers it and the rules of it a freeze
// covers. Its own mark wins over the one it inherits.
func applyFreezes(g *Graph, files []scan.File, anchors map[string]string) {
	byPath := map[string]int{}
	for i, n := range g.Nodes {
		byPath[n.ID] = i
	}
	unitFrozen := map[string]string{}
	for _, f := range files {
		i, ok := byPath[f.Path]
		if !ok {
			continue
		}
		n := &g.Nodes[i]
		n.Frozen = f.Frozen
		n.FrozenRules = copyMap(f.FrozenRules)
		n.Freezes = copyMap(f.Freezes)
		// A unit is frozen by its owner's header: the spec, or the file that owns its code.
		if f.Frozen != "" && f.HeaderCode != "" && len(f.HeaderRefs) == 0 {
			unitFrozen[f.HeaderCode] = f.Frozen
		}
	}
	for _, f := range files {
		i, ok := byPath[f.Path]
		if !ok || g.Nodes[i].Frozen != "" {
			continue
		}
		if v := unitFrozen[unitOfFile(f, anchors)]; v != "" {
			g.Nodes[i].Frozen = v
		}
	}
	// A frozen doctrine document, or rule, freezes the spec rules that realize it.
	for _, e := range g.Edges {
		if e.Type != EdgeRealizes || e.Dep == "" {
			continue
		}
		di, dok := byPath[e.To]
		si, sok := byPath[e.From]
		if !dok || !sok {
			continue
		}
		v := g.Nodes[di].FrozenRules[e.Method]
		if v == "" {
			v = g.Nodes[di].Frozen
		}
		if v == "" || g.Nodes[si].FrozenRules[e.Dep] != "" {
			continue
		}
		if g.Nodes[si].FrozenRules == nil {
			g.Nodes[si].FrozenRules = map[string]string{}
		}
		g.Nodes[si].FrozenRules[e.Dep] = v
	}
}

func copyMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// FreezeRecord is a freeze's entry: its reason and the doctrine document that records it.
type FreezeRecord struct {
	Code   string
	Reason string
	Doc    string
}

// FreezeRecords are the freezes the doctrine records, by code.
func (g *Graph) FreezeRecords() map[string]FreezeRecord {
	out := map[string]FreezeRecord{}
	for _, n := range g.Nodes {
		for c, r := range n.Freezes {
			out[c] = FreezeRecord{Code: c, Reason: r, Doc: n.ID}
		}
	}
	return out
}

// DescribeFreeze is how a mark's value reads to whoever reaches the part: a freeze's code
// with its recorded reason (`QLDTQ-F01 — the plan mode ships next quarter`), or the reason
// itself where the mark carries it.
func (g *Graph) DescribeFreeze(value string) string {
	if !scan.IsFreezeCode(value) {
		return value
	}
	for _, n := range g.Nodes {
		if r, ok := n.Freezes[value]; ok && r != "" {
			return value + " — " + r
		}
	}
	return value
}

// FrozenRules are every rule a freeze covers in the project, each with the freeze's value.
func (g *Graph) FrozenRules() map[string]string {
	out := map[string]string{}
	for _, n := range g.Nodes {
		for r, v := range n.FrozenRules {
			out[r] = v
		}
	}
	return out
}

// FrozenPart is one place a freeze covers: a whole file, or one rule of it.
type FrozenPart struct {
	Path  string
	Rule  string
	Value string
}

// FrozenParts are every place a freeze covers, by path and then rule; with a value, only
// the places of that freeze.
func (g *Graph) FrozenParts(value string) []FrozenPart {
	var out []FrozenPart
	for _, n := range g.Nodes {
		if n.Frozen != "" && (value == "" || n.Frozen == value) {
			out = append(out, FrozenPart{Path: n.ID, Value: n.Frozen})
		}
		if n.Frozen != "" {
			continue // the whole file is covered: its rules are no news
		}
		for r, v := range n.FrozenRules {
			if value == "" || v == value {
				out = append(out, FrozenPart{Path: n.ID, Rule: r, Value: v})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Rule < out[j].Rule
	})
	return out
}

// FrozenAt is the freeze that covers the whole file at path, or empty. A file with only
// some rules frozen is not frozen: those rules are in its FrozenRules.
func (g *Graph) FrozenAt(path string) string {
	path = strings.TrimPrefix(path, "./")
	for _, n := range g.Nodes {
		if n.ID == path {
			return n.Frozen
		}
	}
	return ""
}

// Unfrozen is the map without what a freeze covers whole — its nodes, and the edges that
// touch them —, for the reports that measure the project: a frozen part is kept as it is
// by decision, and its stale proofs, survived mutants and uncovered lines are no news.
func (g *Graph) Unfrozen() *Graph {
	frozen := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Frozen != "" {
			frozen[n.ID] = true
		}
	}
	if len(frozen) == 0 {
		return g
	}
	out := *g
	out.Nodes = nil
	out.Edges = nil
	for _, n := range g.Nodes {
		if !frozen[n.ID] {
			out.Nodes = append(out.Nodes, n)
		}
	}
	for _, e := range g.Edges {
		if !frozen[e.From] && !frozen[e.To] {
			out.Edges = append(out.Edges, e)
		}
	}
	return &out
}

// WithoutFrozenRules are the codes a freeze does not cover, in their order.
func (n Node) WithoutFrozenRules(codes []string) []string {
	if len(n.FrozenRules) == 0 {
		return codes
	}
	var out []string
	for _, c := range codes {
		if _, frozen := n.FrozenRules[c]; !frozen {
			out = append(out, c)
		}
	}
	return out
}

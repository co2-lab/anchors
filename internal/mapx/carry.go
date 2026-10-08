// @anchors
//   code: CRYEV
//   ref: EDSTD

package mapx

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// A proof goes stale when what it proves changes, not when the file does
// (DESIGN-evidence-by-what-it-proves.md). A rebuild used to drop the signal of every file whose
// revision moved, and a date, a navigation row or a flag comment cost a rerun of the suites
// that proved it — 51 MIF specs lost their e2e proofs over navigation tables, because no one
// ran `keep-evidence` on them. The evidence revisions the scan reads leave out what proves
// nothing; when they held, the evidence goes along with the file, with nobody to remember it.

func (n Node) evidenceRev() string {
	if n.EvidenceRev != "" {
		return n.EvidenceRev
	}
	return n.Rev
}

func (n Node) lineRev() string {
	if n.LineRev != "" {
		return n.LineRev
	}
	return n.Rev
}

// setEvidence gives the node its revision and the evidence revisions of its content, each
// left empty when it equals the revision.
func (n *Node) setEvidence(rev string, e scan.Evidence) {
	n.Rev, n.EvidenceRev, n.LineRev, n.RuleRevs, n.RestRev = rev, e.Rev, e.LineRev, e.Rules, e.Rest
	if n.EvidenceRev == rev {
		n.EvidenceRev = ""
	}
	if n.LineRev == rev {
		n.LineRev = ""
	}
}

// AdvanceTo moves a node to its file's content as it is now, carrying the evidence its
// evidence revisions show unchanged (CarryUnchangedEvidence). It says whether the revision
// moved, and whether the evidence went along; an unknown file is left alone.
func (g *Graph) AdvanceTo(id string, content []byte, cfg *config.Config) (moved, carried bool) {
	i := g.nodeIndex(id)
	if i < 0 {
		return false, false
	}
	rev := scan.ShortHash(content)
	if g.Nodes[i].Rev == rev {
		return false, false
	}
	n := g.Nodes[i]
	n.setEvidence(rev, scan.EvidenceOf(string(n.Kind), content, cfg))
	carried = len(CarryUnchangedEvidence(&Graph{Nodes: []Node{n}}, g)) > 0
	g.Nodes[i].setEvidence(rev, scan.EvidenceOf(string(n.Kind), content, cfg))
	return true, carried
}

// FillOldEvidence gives a map written before evidence revisions the evidence of each file the
// new build moved: its content at the map's revision is read by `read` — the commit's copy,
// usually —, and taken only when its hash is that revision. Without it, the first build with
// evidence revisions had nothing to compare a spec's edit with, and the upgrade carried
// nothing. A map that has evidence revisions is left as it is.
func FillOldEvidence(novo, old *Graph, read func(id string) ([]byte, bool), cfg *config.Config) {
	if novo == nil || old == nil || read == nil {
		return
	}
	for _, o := range old.Nodes {
		if o.EvidenceRev != "" || o.RestRev != "" {
			return
		}
	}
	for _, n := range novo.Nodes {
		i := old.nodeIndex(n.ID)
		if i < 0 || old.Nodes[i].Rev == n.Rev {
			continue
		}
		b, ok := read(n.ID)
		if !ok || scan.ShortHash(b) != old.Nodes[i].Rev {
			continue
		}
		o := &old.Nodes[i]
		o.setEvidence(o.Rev, scan.EvidenceOf(string(o.Kind), b, cfg))
	}
}

// CarryUnchangedEvidence moves, in `old`, to the revision each file has in `novo` the evidence
// of the files whose revision moved and whose evidence revision did not: their proofs, the
// closures of the tests that reach them, their edges' stamps — what held at the old revision
// only, as a mechanical repair does. Coverage and mutation go along only when the lines held
// too. A spec whose rules changed and nothing else carries its evidence too, with the
// scenarios of the rules that changed marked stale. It runs before what `old` keeps is given to `novo`, and
// returns the files it carried.
func CarryUnchangedEvidence(novo, old *Graph) []string {
	if novo == nil || old == nil {
		return nil
	}
	var carried []string
	for _, n := range novo.Nodes {
		i := old.nodeIndex(n.ID)
		if i < 0 {
			continue
		}
		o := old.Nodes[i]
		if o.Rev == "" || o.Rev == n.Rev {
			continue
		}
		if o.evidenceRev() == n.evidenceRev() {
			old.keepEvidence(n.ID, n.Rev, "", "", o.lineRev() == n.lineRev(), true, false)
			carried = append(carried, n.ID)
			continue
		}
		changed, ok := changedRules(o, n)
		if !ok {
			continue
		}
		old.keepEvidence(n.ID, n.Rev, "", "", false, true, false)
		if s := old.Nodes[i].Signal; s != nil {
			s.markStale(changed)
		}
		carried = append(carried, n.ID)
	}
	return carried
}

// changedRules are the rules of a spec whose definition changed between two of its versions,
// when nothing else its evidence reads did; it says false when something else did, or when
// either version has no rule revisions to tell.
func changedRules(o, n Node) (map[string]bool, bool) {
	if o.RestRev == "" || o.RestRev != n.RestRev || len(o.RuleRevs) == 0 || len(n.RuleRevs) == 0 {
		return nil, false
	}
	changed := map[string]bool{}
	for c, r := range o.RuleRevs {
		if n.RuleRevs[c] != r {
			changed[c] = true
		}
	}
	for c := range n.RuleRevs {
		if _, ok := o.RuleRevs[c]; !ok {
			changed[c] = true
		}
	}
	return changed, true
}

// markStale marks stale the proven scenarios of the given rules — each variant `#NN` with
// its rule. Their proof stays where it was: a run that proves them again, or the author
// keeping the evidence, makes them fresh.
func (s *TestSignal) markStale(rules map[string]bool) {
	if len(rules) == 0 {
		return
	}
	seen := map[string]bool{}
	for _, c := range s.StaleCodes {
		seen[c] = true
	}
	for _, c := range s.ProvenCodes {
		base, _, _ := strings.Cut(c, "#")
		if rules[base] && !seen[c] {
			seen[c] = true
			s.StaleCodes = append(s.StaleCodes, c)
		}
	}
	sort.Strings(s.StaleCodes)
}

// refreshStale keeps stale only what is still proven and was not proven again now.
func (s *TestSignal) refreshStale(provenNow []string) {
	if len(s.StaleCodes) == 0 {
		return
	}
	still, again := map[string]bool{}, map[string]bool{}
	for _, c := range s.ProvenCodes {
		still[c] = true
	}
	for _, c := range provenNow {
		again[c] = true
	}
	var out []string
	for _, c := range s.StaleCodes {
		if still[c] && !again[c] {
			out = append(out, c)
		}
	}
	s.StaleCodes = out
}

// ClearStale makes fresh the scenarios marked stale, and returns them: the author declared
// the change to their rules proves nothing new (`anchors keep-evidence`).
func (s *TestSignal) ClearStale() []string {
	if s == nil {
		return nil
	}
	out := s.StaleCodes
	s.StaleCodes = nil
	return out
}

// FreshProven are the proven scenarios whose proof still stands: the proven ones, without
// those whose rule changed since (StaleCodes).
func (s *TestSignal) FreshProven() []string {
	if s == nil {
		return nil
	}
	if len(s.StaleCodes) == 0 {
		return s.ProvenCodes
	}
	stale := map[string]bool{}
	for _, c := range s.StaleCodes {
		stale[c] = true
	}
	var out []string
	for _, c := range s.ProvenCodes {
		if !stale[c] {
			out = append(out, c)
		}
	}
	return out
}

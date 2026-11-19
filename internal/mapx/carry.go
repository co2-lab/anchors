// @anchors
//   code: CRYEV
//   ref: EDSTD

package mapx

import (
	"sort"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// A proof goes stale when what it proves changes, not when the file does.
// A rebuild used to drop the signal of every file whose
// revision moved, and a date, a navigation row or a flag comment cost a rerun of the suites
// that proved it — 51 specs of a project lost their e2e proofs over navigation tables, because no one
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
	n.KeptLines, n.CoverageDirectives = "", e.Directives
	if n.EvidenceRev == rev {
		n.EvidenceRev = ""
	}
	if n.LineRev == rev {
		n.LineRev = ""
	} else if n.LineRev != "" {
		n.KeptLines = encodeRanges(e.Kept)
	}
}

// keptLines are the numbers of the lines the node's evidence reads, and whether they are known:
// every line when its line revision is its revision, the recorded ones otherwise.
func (n Node) keptLines(count int) ([]int, bool) {
	if n.LineRev == "" {
		out := make([]int, count)
		for i := range out {
			out[i] = i + 1
		}
		return out, true
	}
	if n.KeptLines == "" {
		return nil, false
	}
	return decodeRanges(n.KeptLines), true
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
	n.setEvidence(rev, scan.EvidenceOf(string(n.Kind), n.ID, content, cfg))
	carried = len(CarryUnchangedEvidence(&Graph{Nodes: []Node{n}}, g)) > 0
	g.Nodes[i].setEvidence(rev, scan.EvidenceOf(string(n.Kind), n.ID, content, cfg))
	return true, carried
}

// TakeContent gives the node the revision and the evidence revisions of its file's content
// as it is now, carrying nothing: what moved its evidence by declaration (`keep-evidence`,
// `check --fix`) moved its revision alone, and the evidence revisions left those of the
// previous content. The next rebuild compared them with the new content's, found the
// evidence changed, and dropped what the declaration had kept (reported from a project: a header
// date after `keep-evidence --lines` took the mutation of 271 files).
func (g *Graph) TakeContent(id string, content []byte, cfg *config.Config) {
	if i := g.nodeIndex(id); i >= 0 {
		n := &g.Nodes[i]
		n.setEvidence(scan.ShortHash(content), scan.EvidenceOf(string(n.Kind), n.ID, content, cfg))
	}
}

// FillOldEvidence gives each file the new build moved the evidence of its content at the
// map's revision, read again by today's rule: the content is read by `read` — the commit's
// copy, usually —, and taken only when its hash is that revision. A map written before
// evidence revisions had nothing to compare an edit with; a map written by an older rule —
// before comment lines left the evidence — compared an edit with a revision today's rule
// never makes, and carried nothing.
func FillOldEvidence(novo, old *Graph, read func(id string) ([]byte, bool), cfg *config.Config) {
	if novo == nil || old == nil || read == nil {
		return
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
		o.setEvidence(o.Rev, scan.EvidenceOf(string(o.Kind), o.ID, b, cfg))
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
			lines := o.lineRev() == n.lineRev()
			lineTo, remap := lineMapping(o, n)
			old.keepEvidence(n.ID, n.Rev, "", "", lines, true, false)
			if !lines && remap {
				old.remapCoverage(n.ID, o.Rev, n.Rev, lineTo)
			}
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

// FollowMoves carries, in `old`, each file git moved to its new path — the node, the edges at it,
// the closures and doubles of the tests that name it —, so the rebuild carries its evidence as
// it carries any file's: a node is known by its path, and a move lost every proof of a file
// whose content had not changed (reported from a project: 96 files moved to a package, every one
// asking for its mutation again; the user: "se um arquivo foi movido no git o anchors pode
// identificar isso e mover no mapa também"). The content then decides, as for any edit: the
// same revision keeps everything, the same evidence carries it, another drops it. A move is
// taken only from a path `novo` no longer has to one `old` does not have. It returns the new
// paths moved, sorted.
func FollowMoves(novo, old *Graph, moves map[string]string) []string {
	if novo == nil || old == nil || len(moves) == 0 {
		return nil
	}
	froms := make([]string, 0, len(moves))
	for from := range moves {
		froms = append(froms, from)
	}
	sort.Strings(froms)
	var moved []string
	for _, from := range froms {
		to := moves[from]
		if old.nodeIndex(from) < 0 || old.nodeIndex(to) >= 0 || novo.nodeIndex(to) < 0 || novo.nodeIndex(from) >= 0 {
			continue
		}
		old.renameNode(from, to)
		moved = append(moved, to)
	}
	sort.Strings(moved)
	return moved
}

// renameNode gives the node `from` the path `to`, and every reference to it.
func (g *Graph) renameNode(from, to string) {
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if n.ID == from {
			n.ID = to
		}
		if s := n.Signal; s != nil {
			if r, ok := s.ClosureRev[from]; ok {
				delete(s.ClosureRev, from)
				s.ClosureRev[to] = r
			}
		}
		for j, d := range n.Doubles {
			if d == from {
				n.Doubles[j] = to
			}
		}
	}
	g.dropNeighborIndex()
	for i := range g.Edges {
		if g.Edges[i].From == from {
			g.Edges[i].From = to
		}
		if g.Edges[i].To == from {
			g.Edges[i].To = to
		}
	}
}

// lineMapping moves a line of the old content to the new one, when the two hold the same code
// in the same order and the same coverage directives: the i-th line the old evidence reads is
// the i-th of the new. It answers false when the layout of either is unknown — a map written
// before it was recorded —, or a comment that steers the coverage tool changed.
func lineMapping(o, n Node) (func(int) (int, bool), bool) {
	if o.CoverageDirectives != n.CoverageDirectives {
		return nil, false
	}
	oldKept, okOld := o.keptLines(0)
	newKept, okNew := n.keptLines(0)
	switch {
	case !okOld || !okNew:
		return nil, false
	case o.LineRev == "" && n.LineRev == "":
		return nil, false // every line in place: nothing to move
	case o.LineRev == "":
		oldKept, _ = o.keptLines(len(newKept))
	case n.LineRev == "":
		newKept, _ = n.keptLines(len(oldKept))
	}
	if len(oldKept) != len(newKept) {
		return nil, false
	}
	to := make(map[int]int, len(oldKept))
	for i, l := range oldKept {
		to[l] = newKept[i]
	}
	return func(l int) (int, bool) { nl, ok := to[l]; return nl, ok }, true
}

// remapCoverage moves the node's coverage measured at `from` to the revision `to`, each line
// to its new number: the lines each suite instrumented and covered, and the branches no suite
// took. The mutation is not moved — a comment may be a directive of the mutation tool, and its
// score is measured again.
func (g *Graph) remapCoverage(id, from, to string, shift func(int) (int, bool)) {
	i := g.nodeIndex(id)
	if i < 0 || g.Nodes[i].Signal == nil {
		return
	}
	s := g.Nodes[i].Signal
	moved := false
	for k, c := range s.CoverageBySuite {
		if c.AtRev != from {
			continue
		}
		c.Instrumented = encodeRanges(shiftLines(decodeRanges(c.Instrumented), shift))
		c.Covered = encodeRanges(shiftLines(decodeRanges(c.Covered), shift))
		c.BranchMissed = shiftBranches(c.BranchMissed, shift)
		c.AtRev = to
		s.CoverageBySuite[k] = c
		moved = true
	}
	if s.CoverageRev == from {
		s.CoverageRev = to
		moved = true
	}
	if !moved && s.TotalLines == 0 {
		return
	}
	s.BranchMissed = shiftBranches(s.BranchMissed, shift)
	if s.AtRev == from {
		// The shared revision stamps the mutation too when it has none of its own: it keeps
		// the old one, and stays stale.
		if s.MutationAtRev == "" && (s.MutantsKilled+s.MutantsSurvived+s.MutantsNoCoverage+s.MutantsIgnored+s.MutantsStatic > 0 || len(s.MutationByScope) > 0) {
			s.MutationAtRev = from
		}
		s.AtRev = to
	}
}

func shiftLines(lines []int, shift func(int) (int, bool)) []int {
	var out []int
	for _, l := range lines {
		if nl, ok := shift(l); ok {
			out = append(out, nl)
		}
	}
	sort.Ints(out)
	return out
}

// shiftBranches moves each `<line>:<block>:<branch>` to the line's new number.
func shiftBranches(s string, shift func(int) (int, bool)) string {
	if strings.TrimSpace(s) == "" {
		return s
	}
	var out []string
	for _, b := range strings.Fields(s) {
		line, rest, ok := strings.Cut(b, ":")
		n, err := strconv.Atoi(line)
		if !ok || err != nil {
			out = append(out, b)
			continue
		}
		if nl, ok := shift(n); ok {
			out = append(out, strconv.Itoa(nl)+":"+rest)
		}
	}
	return strings.Join(out, " ")
}

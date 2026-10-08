// @anchors
//   code: INMPN
//   ref: GRINC

package mapx

import (
	"os"
	"path"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// ADDING NEW FILES WITHOUT READING THE TREE.
//
// A file created after the last `map build` had no node, and whatever reached the map for it
// was dropped: `anchors test` proved a new spec's rules and had nowhere to write the proof,
// so `scenario-coverage` blocked a spec whose tests passed. Rebuilding the whole map fixes
// it and reads every file of the tree for the sake of one.
//
// The map already knows most of what a new file relates to: every path, kind, layer and
// declared code. So the new files are read, and every other file STANDS IN as what its node
// keeps; the same `Build` runs over that set, in memory, and only what involves the new
// files is merged. Files whose content decides the new files' relations — the siblings the
// derivation links them to, and the anchor of their unit — are read too, in a second pass.
//
// What stays for the next `map build`: a relation an EXISTING file declares toward the new
// one in its text — a row of its dependencies table, a plan's seed list, an `@realizes`, a
// scenario code cited in a test of another directory. The map does not keep those
// declarations, and finding them means reading the files that make them. The commit hook
// and the CI still demand the full map.

// Reader reads files by path, as the scan reads them in the tree (`scan.ScanPaths`).
type Reader func(paths []string) ([]scan.File, error)

// AddFiles brings the files at `paths` into the map, if they are not in it yet, with the
// relations they take part in, reading only them and the files their unit links them to.
// It returns the ids of the nodes it added. The rest of the map — nodes, revs, signals,
// stamps — is left as it was, except the derivation links inside the re-read unit, which
// are replaced by what the derivation says now (a new feature, for one, takes the place of
// the spec→test link a unit with no feature had).
func (g *Graph) AddFiles(paths []string, read Reader, cfg *config.Config, updatedAt map[string]string) ([]string, error) {
	known := map[string]bool{}
	for _, n := range g.Nodes {
		known[n.ID] = true
	}
	var fresh []string
	for _, p := range paths {
		if !known[p] {
			fresh = append(fresh, p)
		}
	}
	if len(fresh) == 0 {
		return nil, nil
	}
	newFiles, err := read(fresh)
	if err != nil {
		return nil, err
	}
	isNew := map[string]bool{}
	for _, f := range newFiles {
		isNew[f.Path] = true
	}
	if len(isNew) == 0 {
		return nil, nil // ignored, or outside every layer: nothing for the map
	}

	// Pass 1: which existing files does the derivation link to the new ones?
	first := Build(g.standIns(newFiles), cfg, updatedAt)
	var reread []string
	seen := map[string]bool{}
	want := func(id string) {
		if !isNew[id] && !seen[id] && known[id] {
			seen[id] = true
			reread = append(reread, id)
		}
	}
	// The whole UNIT, not only the files linked straight to a new one: a new anchor also
	// decides links between existing siblings (its feature→test), so the derivation links
	// are followed until the unit closes.
	inReach := func(id string) bool { return isNew[id] || seen[id] }
	for grew := true; grew; {
		grew = false
		for _, e := range first.Edges {
			if e.Origin == OriginConvention && (inReach(e.From) || inReach(e.To)) && (!inReach(e.From) || !inReach(e.To)) {
				want(e.From)
				want(e.To)
				grew = true
			}
		}
	}
	// The anchors that may own a new file decide its unit even when pass 1, reading them as
	// stand-ins, drew no link: an override chosen by the anchor's HEADER layer is only in its
	// text, and without it the stand-in takes the default templates. An anchor may own the
	// new file when it sits beside it, or when the new file carries its name — the templates
	// keep `{{name}}` wherever an override puts the file (`Tela.tsx` →
	// `__tests__/Tela.test.tsx`).
	if cfg != nil && cfg.Derived != nil {
		for _, n := range g.Nodes {
			if string(n.Kind) != cfg.Derived.Anchor {
				continue
			}
			name, _ := StemOfAnchor(n.ID)
			for id := range isNew {
				b := path.Base(id)
				if path.Dir(id) == path.Dir(n.ID) || b == name || strings.HasPrefix(b, name+".") {
					want(n.ID)
				}
			}
		}
	}
	unit, err := read(reread)
	if err != nil {
		return nil, err
	}

	// Pass 2: the new files and their unit read, everything else as the map keeps it.
	read2 := append(append([]scan.File{}, newFiles...), unit...)
	second := Build(g.standIns(read2), cfg, updatedAt)
	inUnit := map[string]bool{}
	for _, f := range read2 {
		inUnit[f.Path] = true
	}

	var added []string
	for _, n := range second.Nodes {
		switch {
		case isNew[n.ID]:
			g.Nodes = append(g.Nodes, n)
			added = append(added, n.ID)
		case inUnit[n.ID]:
			// A new anchor gives its derived siblings its declared code.
			if i := g.nodeIndex(n.ID); i >= 0 && !g.Nodes[i].CodeDeclarado {
				g.Nodes[i].Code = n.Code
			}
		}
	}
	kept := g.Edges[:0]
	for _, e := range g.Edges {
		if e.Origin == OriginConvention && inUnit[e.From] && inUnit[e.To] {
			continue // replaced below by what the derivation says now
		}
		kept = append(kept, e)
	}
	g.Edges = kept
	have := map[string]bool{}
	for _, e := range g.Edges {
		have[edgeIdentity(e)] = true
	}
	for _, e := range second.Edges {
		inside := e.Origin == OriginConvention && inUnit[e.From] && inUnit[e.To]
		touches := isNew[e.From] || isNew[e.To]
		if (inside || touches) && !have[edgeIdentity(e)] {
			g.Edges = append(g.Edges, e)
			have[edgeIdentity(e)] = true
		}
	}
	sortGraph(g)
	return added, nil
}

// standIns is the file set a partial build reads: the files given, and every other node of
// the map as what it keeps — path, kind, layer, declared code, support and upstream flags.
func (g *Graph) standIns(read []scan.File) []scan.File {
	out := append([]scan.File{}, read...)
	isRead := map[string]bool{}
	for _, f := range read {
		isRead[f.Path] = true
	}
	for _, n := range g.Nodes {
		if isRead[n.ID] {
			continue
		}
		f := scan.File{Path: n.ID, Kind: string(n.Kind), Layer: n.Layer, Rev: n.Rev,
			Support: n.Support, Upstream: n.Upstream, Parent: n.Parent, Needs: n.Needs, Revises: n.Revises, OutRows: n.OutRows}
		if n.CodeDeclarado {
			f.HeaderCode = n.Code
		}
		for _, r := range n.Resources {
			if kind, name, ok := strings.Cut(r, ":"); ok {
				f.CodeDeps = append(f.CodeDeps, scan.CodeDep{Kind: kind, Name: name})
			}
		}
		out = append(out, f)
	}
	return out
}

func (g *Graph) nodeIndex(id string) int {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return i
		}
	}
	return -1
}

func edgeIdentity(e Edge) string {
	return string(e.Type) + "\x00" + e.From + "\x00" + e.To + "\x00" + e.Method + "\x00" + e.Dep
}

// AddFilesAt brings new files into the map on disk: under the map's lock, it reads them —
// and their unit — from the tree and adds their nodes and relations. A project with no map
// yet is left alone: the first `map build` makes it whole. It returns the nodes added.
func AddFilesAt(root, mapPath string, paths []string, cfg *config.Config) ([]string, error) {
	if _, err := os.Stat(mapPath); err != nil {
		return nil, nil
	}
	read := func(ps []string) ([]scan.File, error) { return scan.ScanPaths(root, cfg, ps) }
	var added []string
	err := Update(mapPath, func(g *Graph) error {
		var err error
		added, err = g.AddFiles(paths, read, cfg, nil)
		return err
	})
	return added, err
}

// AddMissingAt adds to the map every governed file of the tree it does not have yet —
// listing the tree's paths, reading only the missing files and their units. It is the net
// under the other two ways in (`anchors new`, the watcher) for a file created by hand.
func AddMissingAt(root, mapPath string, cfg *config.Config) ([]string, error) {
	paths, err := scan.GovernedPaths(root, cfg)
	if err != nil {
		return nil, err
	}
	return AddFilesAt(root, mapPath, paths, cfg)
}

// @anchors
//   code: MXRSV
//   ref: GRBLG

package mapx

import (
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
)

// Modules resolves an import path — or a double's specifier — into the map's files, the way
// the dependency chain does: relative to the importer or through an alias the project
// declares, as written or with each of its extensions. It is the map's own, built once per
// graph, so the map can resolve what its files name while it is built, and a gate resolves
// the same way.
type Modules struct {
	kind map[string]Kind
	dirs map[string]bool // directories that hold a code file
	code []string        // the code files, sorted
}

var (
	modulesMu  sync.Mutex
	modulesFor *Graph
	modulesLen int
	modulesIdx *Modules
)

// Modules is the graph's index of files for resolving imports, cached while the graph keeps
// the same nodes.
func (g *Graph) Modules() *Modules {
	modulesMu.Lock()
	defer modulesMu.Unlock()
	if modulesFor == g && modulesLen == len(g.Nodes) && modulesIdx != nil {
		return modulesIdx
	}
	m := &Modules{kind: map[string]Kind{}, dirs: map[string]bool{}}
	for _, n := range g.Nodes {
		m.kind[n.ID] = n.Kind
		if n.Kind == KindCode {
			m.dirs[path.Dir(n.ID)] = true
			m.code = append(m.code, n.ID)
		}
	}
	sort.Strings(m.code)
	modulesFor, modulesLen, modulesIdx = g, len(g.Nodes), m
	return m
}

// Resolve turns an import's path into a file of the map: relative to the importer (`from`),
// or through an alias the project declares; then as written, or with each extension. A path
// that names a directory of governed files is a package (dir true). Anything else is outside
// the project ("").
func (m *Modules) Resolve(from, p string, relative bool, d config.Dialect) (string, bool) {
	if r := d.ImportResolve; r != nil && r.Separator != "" && r.Separator != "/" {
		p = withSlashes(p, r.Separator)
	}
	var base string
	switch {
	// A pattern may say its paths are relative to the importer without `./` — a flow names
	// the flow it runs that way.
	case strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") || p == "." || p == ".." || relative:
		base = path.Clean(path.Join(path.Dir(from), p))
	default:
		aliases := d.ImportResolve.AliasesFor(from)
		best := ""
		for prefix := range aliases {
			if strings.HasPrefix(p, prefix) && len(prefix) > len(best) {
				best = prefix
			}
		}
		if best == "" {
			return "", false
		}
		base = path.Clean(aliases[best] + strings.TrimPrefix(p, best))
	}
	if _, ok := m.kind[base]; ok {
		return base, false
	}
	if r := d.ImportResolve; r != nil {
		for _, e := range r.Extensions {
			if _, ok := m.kind[base+e]; ok {
				return base + e, false
			}
		}
	}
	if m.dirs[base] {
		return base, true
	}
	return "", false
}

// withSlashes writes an import path whose directories another separator divides with `/`:
// leading separators climb from the importer (`.x` → `./x`, `..x` → `../x`).
func withSlashes(p, sep string) string {
	n := 0
	for strings.HasPrefix(p, sep) {
		p, n = strings.TrimPrefix(p, sep), n+1
	}
	p = strings.ReplaceAll(p, sep, "/")
	switch {
	case n == 1:
		p = "./" + p
	case n > 1:
		p = strings.Repeat("../", n-1) + p
	}
	return p
}

// Candidates are the code files of the map a double's specifier may name: the one its import
// resolution gives, or every one whose path ends as the specifier does.
func (m *Modules) Candidates(spec, from string, d config.Dialect) []string {
	if target, dir := m.Resolve(from, spec, false, d); target != "" && !dir && m.kind[target] == KindCode {
		return []string{target}
	}
	target := strings.TrimPrefix(spec, "./")
	for strings.HasPrefix(target, "../") {
		target = strings.TrimPrefix(target, "../")
	}
	// An alias the project declares maps the head (`@/` → `apps/mobile/`). When the file it
	// gives is not in the map, the tail is matched alone: the same alias names another
	// workspace's files too (`@/lib/x` in `apps/landing-page/src/`), and those are still the
	// project's modules.
	var tail string
	aliases := d.ImportResolve.AliasesFor(from)
	best := ""
	for prefix := range aliases {
		if strings.HasPrefix(target, prefix) && len(prefix) > len(best) {
			best = prefix
		}
	}
	if best != "" {
		tail = strings.TrimPrefix(target, best)
		target = strings.TrimPrefix(aliases[best]+tail, "/")
	}
	if cands := m.codeEndingIn(target); len(cands) > 0 || tail == "" {
		return cands
	}
	return m.codeEndingIn(tail)
}

// File is the ONE code file a double's specifier names: its only candidate, or the one that
// shares the most folders with the importer. Several the folder does not decide between, or
// none, give "" — with the candidates, so the caller can say why.
func (m *Modules) File(spec, from string, d config.Dialect) (string, []string) {
	cands := m.Candidates(spec, from, d)
	switch len(cands) {
	case 0:
		return "", nil
	case 1:
		return cands[0], cands
	}
	best, bestLen, tie := "", -1, false
	for _, c := range cands {
		l := sharedDirs(c, from)
		switch {
		case l > bestLen:
			best, bestLen, tie = c, l, false
		case l == bestLen:
			tie = true
		}
	}
	if tie {
		return "", cands
	}
	return best, cands
}

// codeEndingIn are the code files of the map whose path, without its extension, is the
// target or ends in it after a `/`.
func (m *Modules) codeEndingIn(target string) []string {
	if target == "" {
		return nil
	}
	var cands []string
	for _, id := range m.code {
		x := WithoutExtension(id)
		if x == target || strings.HasSuffix(x, "/"+target) {
			cands = append(cands, id)
		}
	}
	return cands
}

// sharedDirs counts the leading directory segments two paths share.
func sharedDirs(a, b string) int {
	sa, sb := strings.Split(filepath.ToSlash(a), "/"), strings.Split(filepath.ToSlash(b), "/")
	n := 0
	for n < len(sa)-1 && n < len(sb)-1 && sa[n] == sb[n] {
		n++
	}
	return n
}

// WithoutExtension drops the path's last extension (`x/y.ts` → `x/y`), so a node of the map
// compares with an import's specifier, which never carries one.
func WithoutExtension(id string) string {
	if i := strings.LastIndex(id, "."); i > strings.LastIndex(id, "/") {
		return id[:i]
	}
	return id
}

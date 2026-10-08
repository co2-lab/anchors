// @anchors
//   code: NCGNA
//   ref: NCGNV

package gate

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// THE NAVIGATION CHAIN, flagged in the code and confronted with the screens' specs
// (DESIGN-dependencies-and-navigation.md):
//
//	nav-annotated     every navigation call carries `@navigates:` (or `@no-nav:`), naming
//	                  the screen its route leads to
//	nav-matches-spec  the screen's Out table and its code's `@navigates:` name the same screens
//	nav-symmetric     A's Out leads to B exactly when B's In comes from A
//	nav-reachable     every screen is reachable from the app's entry routes
//
// Every navigation counts — `goBack`, `reset`, `popToTop` too: a back navigation names
// the screens it returns to. Anchors reads the calls with a pattern (the dialect's
// `navigation_call`, or the family's) and the screens by their specs' route and name.

// Screen is a screen of the app, as its spec declares it.
type Screen struct {
	Spec  string   // the spec's path
	Code  string   // the unit's code
	Codes []string // every code that names it: the unit's, and its files' own
	Name  string   // the spec's stem — `GoalEditScreen`
	Route string   // the route the spec declares — `GoalEdit`
	Out   []NavRow // where it leads (the Out table)
	In    []NavRow // where it comes from (the In table)
}

// NavRow is a row of a navigation table: the screen it names, and the rule of the row.
type NavRow struct {
	Name string // as written: a screen's name or route
	Rule string
}

// familyNavigationCall is the navigation call a family knows, its destination in `route`.
func familyNavigationCall(family string) string {
	switch strings.ToLower(family) {
	case "ts", "js", "typescript", "javascript":
		return `\b(?:navigation|navigator|router|nav)\s*\.\s*(?:navigate|push|replace|reset|goBack|back|pop|popToTop|dispatch)\s*\(\s*(?:\{\s*(?:name|pathname)\s*:\s*)?(?:['"](?P<route>[A-Za-z0-9_/.:-]+)['"])?`
	}
	return ""
}

// navigationCallRE is the pattern of a navigation call: the project's, or the family's.
func navigationCallRE(d config.Dialect) *regexp.Regexp {
	if p := strings.TrimSpace(d.NavigationCall); p != "" {
		if re := d.Compile(p); re != nil {
			return re
		}
	}
	if p := familyNavigationCall(d.Family); p != "" {
		return regexp.MustCompile(p)
	}
	return nil
}

var (
	navInHeadingRE  = regexp.MustCompile(`(?i)^#{2,4}\s+(?:Entrada|Entry|In|Incoming|Origem|Origin)\b`)
	navOutHeadingRE = regexp.MustCompile(`(?i)^#{2,4}\s+(?:Sa[íi]da|Exit|Out|Outgoing|Destino|Destination)\b`)
	navToCols       = []string{"destino", "destination", "to", "target", "alvo"}
	navFromCols     = []string{"origem", "origin", "from", "source", "fonte"}
	navRuleCols     = []string{"regra", "rule", "regla"}
)

// navTables reads a screen spec's In and Out tables.
func navTables(content string) (in, out []NavRow) {
	var cur *[]NavRow
	var nameCols []string
	var header []string
	for _, line := range strings.Split(content, "\n") {
		t := strings.TrimSpace(line)
		switch {
		case navInHeadingRE.MatchString(t):
			cur, nameCols, header = &in, navFromCols, nil
			continue
		case navOutHeadingRE.MatchString(t):
			cur, nameCols, header = &out, navToCols, nil
			continue
		case strings.HasPrefix(t, "#"):
			cur = nil
			continue
		}
		if cur == nil || !strings.HasPrefix(t, "|") {
			continue
		}
		cells := cellsOf(t)
		if header == nil {
			for _, c := range cells {
				header = append(header, strings.ToLower(strings.TrimSpace(c)))
			}
			continue
		}
		if strings.Trim(strings.Join(cells, ""), ":- ") == "" {
			continue
		}
		row := map[string]string{}
		for i, c := range cells {
			if i < len(header) {
				row[header[i]] = strings.TrimSpace(c)
			}
		}
		name, _ := col(row, nameCols...)
		rule, _ := col(row, navRuleCols...)
		name = strings.Trim(name, "` *")
		if name == "" || blank(name) {
			continue
		}
		*cur = append(*cur, NavRow{Name: name, Rule: strings.Trim(citedCodeRE.FindString(rule), "`")})
	}
	return in, out
}

// screenIndex is the app's screens, read once per map.
type screenIndex struct {
	g       *mapx.Graph
	screens []*Screen
	byName  map[string]*Screen // by stem and by route, lower-cased
	byCode  map[string]*Screen
	bySpec  map[string]*Screen
}

var (
	screensMu    sync.Mutex
	screensCache = map[string]*screenIndex{}
)

// screensOf is the index of the app's screens: the specs whose layer is `screen`.
func screensOf(root string, g *mapx.Graph) *screenIndex {
	screensMu.Lock()
	defer screensMu.Unlock()
	if idx, ok := screensCache[root]; ok && idx.g == g {
		return idx
	}
	idx := &screenIndex{g: g, byName: map[string]*Screen{}, byCode: map[string]*Screen{}, bySpec: map[string]*Screen{}}
	if g == nil {
		return idx
	}
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := readFile(root, n.ID)
		if err != nil || layerOf(n, string(b)) != "screen" {
			continue
		}
		s := &Screen{Spec: n.ID, Code: n.Code, Name: strings.TrimSuffix(path.Base(n.ID), ".spec.md"), Route: declaredRoute(string(b))}
		s.In, s.Out = navTables(string(b))
		if s.Code != "" {
			s.Codes = append(s.Codes, s.Code)
		}
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type == mapx.EdgeSpecifies {
				if c := g.Node(e.To); c != nil && c.FileCode != "" {
					s.Codes = append(s.Codes, c.FileCode)
				}
			}
		}
		idx.screens = append(idx.screens, s)
		idx.bySpec[s.Spec] = s
		for _, c := range s.Codes {
			idx.byCode[c] = s
		}
		for _, k := range []string{s.Name, strings.TrimSuffix(s.Name, "Screen"), s.Route} {
			if k != "" {
				idx.byName[strings.ToLower(strings.Trim(k, "/"))] = s
			}
		}
	}
	sort.Slice(idx.screens, func(i, j int) bool { return idx.screens[i].Spec < idx.screens[j].Spec })
	screensCache[root] = idx
	return idx
}

// screenNamed is the screen a table cell or a route names.
func (idx *screenIndex) screenNamed(name string) *Screen {
	return idx.byName[strings.ToLower(strings.Trim(strings.Trim(name, "` "), "/"))]
}

// screenOfNode is the screen a code file belongs to: the screen whose spec specifies it.
func (idx *screenIndex) screenOfNode(id string) *Screen {
	if idx.g == nil {
		return nil
	}
	for _, e := range idx.g.Neighbors(id).In {
		if e.Type == mapx.EdgeSpecifies {
			if s := idx.bySpec[e.From]; s != nil {
				return s
			}
		}
	}
	return nil
}

// NavCall is a navigation call of a file, with the flag that answers it.
type NavCall struct {
	Line  int
	Route string // the destination the call names; empty for a back or a dynamic one
	Flag  *scan.Navigation
}

// navCallsOf are the navigation calls of a file, each with its flag (on its line, or alone
// on the line above).
func navCallsOf(content string, d config.Dialect) []NavCall {
	re := navigationCallRE(d)
	if re == nil {
		return nil
	}
	route := re.SubexpIndex("route")
	flags := map[int]scan.Navigation{}
	for _, f := range scan.NavigatesIn([]byte(content)) {
		flags[f.Line] = f
	}
	var out []NavCall
	for i, ln := range strings.Split(content, "\n") {
		if commentLine(ln) {
			continue
		}
		m := re.FindStringSubmatch(cutInlineComment(ln))
		if m == nil {
			continue
		}
		c := NavCall{Line: i + 1}
		if route >= 0 {
			c.Route = m[route]
		}
		if f, ok := flags[i+1]; ok {
			f := f
			c.Flag = &f
		}
		out = append(out, c)
	}
	return out
}

// checkNavAnnotated: does every navigation call carry its flag, naming the screen its
// route leads to?
func checkNavAnnotated(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if !chainUnit(n) {
		return Skip, i18n.T("gate.dep_chain.skip_not_code")
	}
	d := cfg.DialectFor()
	if navigationCallRE(d) == nil {
		return Pending, i18n.T("gate.nav_chain.pending_no_pattern")
	}
	idx := screensOf(root, g)
	var gaps []string
	for _, c := range navCallsOf(content, d) {
		if c.Flag == nil {
			gaps = append(gaps, i18n.T("gate.nav_chain.call_unflagged", c.Line, nonEmpty(c.Route)))
			continue
		}
		if c.Flag.Waiver != "" || c.Route == "" {
			continue
		}
		target := idx.screenNamed(c.Route)
		if target == nil {
			continue // a route no screen spec declares: route-exists asks about it
		}
		if !anyIn(c.Flag.Codes, target.Codes) {
			gaps = append(gaps, i18n.T("gate.nav_chain.call_wrong_screen", c.Line, strings.Join(c.Flag.Codes, ", "), c.Route, target.Code))
		}
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ") + i18n.T("gate.nav_chain.guide")
}

func nonEmpty(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func anyIn(xs, ys []string) bool {
	for _, x := range xs {
		if contains(ys, x) {
			return true
		}
	}
	return false
}

// flaggedTargets are the screens the `@navigates:` flags of a screen's code name.
func flaggedTargets(s *Screen, root string, idx *screenIndex) map[string]bool {
	out := map[string]bool{}
	if idx.g == nil {
		return out
	}
	for _, e := range idx.g.Neighbors(s.Spec).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		b, err := readFile(root, e.To)
		if err != nil {
			continue
		}
		for _, f := range scan.NavigatesIn(b) {
			for _, c := range f.Codes {
				if t := idx.byCode[c]; t != nil {
					out[t.Spec] = true
				}
			}
		}
	}
	return out
}

// screenSpec reads a screen spec node into its screen, or nil when the node is no screen.
func screenSpec(n mapx.Node, root string, g *mapx.Graph) (*Screen, *screenIndex) {
	if n.Kind != mapx.KindSpec {
		return nil, nil
	}
	idx := screensOf(root, g)
	return idx.bySpec[n.ID], idx
}

// checkNavMatchesSpec: do the screen's Out table and its code's flags name the same screens?
func checkNavMatchesSpec(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	s, idx := screenSpec(n, root, g)
	if s == nil {
		return Skip, i18n.T("gate.nav_chain.skip_not_screen")
	}
	flagged := flaggedTargets(s, root, idx)
	declared := map[string]bool{}
	var gaps []string
	for _, r := range s.Out {
		t := idx.screenNamed(r.Name)
		if t == nil {
			gaps = append(gaps, i18n.T("gate.nav_chain.out_unknown", r.Name))
			continue
		}
		declared[t.Spec] = true
		if !flagged[t.Spec] {
			gaps = append(gaps, i18n.T("gate.nav_chain.out_unflagged", t.Name, t.Code))
		}
	}
	for spec := range flagged {
		if !declared[spec] {
			t := idx.bySpec[spec]
			gaps = append(gaps, i18n.T("gate.nav_chain.flag_undeclared", t.Name, t.Code))
		}
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	sort.Strings(gaps)
	return Fail, strings.Join(gaps, "; ") + i18n.T("gate.nav_chain.guide")
}

// checkNavSymmetric: does every Out of the screen meet an In of its destination, and every
// In meet an Out of its origin?
func checkNavSymmetric(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	s, idx := screenSpec(n, root, g)
	if s == nil {
		return Skip, i18n.T("gate.nav_chain.skip_not_screen")
	}
	names := func(rows []NavRow) map[string]bool {
		out := map[string]bool{}
		for _, r := range rows {
			if t := idx.screenNamed(r.Name); t != nil {
				out[t.Spec] = true
			}
		}
		return out
	}
	var gaps []string
	for _, r := range s.Out {
		t := idx.screenNamed(r.Name)
		if t != nil && t.Spec != s.Spec && !names(t.In)[s.Spec] {
			gaps = append(gaps, i18n.T("gate.nav_chain.out_without_in", t.Name, s.Name))
		}
	}
	for _, r := range s.In {
		o := idx.screenNamed(r.Name)
		if o != nil && o.Spec != s.Spec && !names(o.Out)[s.Spec] {
			gaps = append(gaps, i18n.T("gate.nav_chain.in_without_out", o.Name, s.Name))
		}
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	sort.Strings(gaps)
	return Fail, strings.Join(gaps, "; ") + i18n.T("gate.nav_chain.guide")
}

// NavEdges are the app's navigation, screen to screen: the Out tables and the code's flags.
func NavEdges(root string, g *mapx.Graph) map[string]map[string]bool {
	idx := screensOf(root, g)
	out := map[string]map[string]bool{}
	add := func(from, to string) {
		if out[from] == nil {
			out[from] = map[string]bool{}
		}
		out[from][to] = true
	}
	for _, s := range idx.screens {
		for _, r := range s.Out {
			if t := idx.screenNamed(r.Name); t != nil {
				add(s.Spec, t.Spec)
			}
		}
		for t := range flaggedTargets(s, root, idx) {
			add(s.Spec, t)
		}
	}
	return out
}

// checkNavReachable: is the screen reachable from the app's entry routes?
func checkNavReachable(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	s, idx := screenSpec(n, root, g)
	if s == nil {
		return Skip, i18n.T("gate.nav_chain.skip_not_screen")
	}
	if cfg == nil || cfg.Navigation == nil || len(cfg.Navigation.Entry) == 0 {
		return Pending, i18n.T("gate.nav_chain.pending_no_entry")
	}
	edges := NavEdges(root, g)
	seen := map[string]bool{}
	var queue []string
	for _, e := range cfg.Navigation.Entry {
		if t := idx.screenNamed(e); t != nil && !seen[t.Spec] {
			seen[t.Spec] = true
			queue = append(queue, t.Spec)
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for next := range edges[cur] {
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	if seen[s.Spec] {
		return Pass, ""
	}
	return Fail, i18n.T("gate.nav_chain.unreachable", s.Name, strings.Join(cfg.Navigation.Entry, ", ")) + i18n.T("gate.nav_chain.guide")
}

// fixNavFlags writes the `@navigates:` of each unflagged navigation call whose route
// resolves to a single screen — a back navigation and a dynamic route are the author's.
func fixNavFlags(content string, n mapx.Node, root string, g *mapx.Graph) (string, bool, string) {
	if !chainUnit(n) || g == nil {
		return content, false, ""
	}
	d := fixConfig.DialectFor()
	idx := screensOf(root, g)
	lines := strings.Split(content, "\n")
	marker := config.LineCommentFor(n.ID)
	var done []string
	for _, c := range navCallsOf(content, d) {
		if c.Flag != nil || c.Route == "" {
			continue
		}
		t := idx.screenNamed(c.Route)
		if t == nil || t.Code == "" {
			continue
		}
		l := lines[c.Line-1]
		if code := strings.TrimSpace(cutInlineComment(l)); strings.HasSuffix(code, ">") {
			// A line ending in a JSX tag: what follows `>` is text the app renders, and a
			// `// @navigates:` there shows on screen (reported from MIF). The flag goes as a
			// block comment right after the call, inside its expression; a call whose end is
			// not found is left to the author.
			at := navCallEnd(l, navigationCallRE(d))
			if at < 0 {
				continue
			}
			lines[c.Line-1] = l[:at] + " /* @navigates: " + t.Code + " */" + l[at:]
		} else {
			lines[c.Line-1] = strings.TrimRight(l, " \t") + " " + marker + " @navigates: " + t.Code
		}
		done = append(done, fmt.Sprintf("%d %s", c.Line, t.Code))
	}
	if len(done) == 0 {
		return content, false, ""
	}
	return strings.Join(lines, "\n"), true, i18n.T("gate.fix.nav_written", strings.Join(done, ", "))
}

// navCallEnd is the index right after the closing parenthesis of the navigation call on a
// line, or -1: from the call's match, the parentheses are counted outside the strings.
func navCallEnd(line string, re *regexp.Regexp) int {
	if re == nil {
		return -1
	}
	m := re.FindStringIndex(line)
	if m == nil {
		return -1
	}
	open := strings.IndexByte(line[m[0]:], '(')
	if open < 0 {
		return -1
	}
	depth := 0
	var quote byte
	for i := m[0] + open; i < len(line); i++ {
		c := line[i]
		if quote != 0 {
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\'', '"', '`':
			quote = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

// @anchors
//   code: DCGDA
//   ref: DCGDP

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

// THE DEPENDENCY CHAIN, flagged in the code and confronted with it
// (DESIGN-dependencies-and-navigation.md):
//
//	dep-declared      every real import of a governed file carries `@dep: <its code>`
//	dep-honored       every `@dep:` names the code of the file its import resolves to
//	used-by-declared  every imported symbol carries `@used-by:` with exactly who imports it
//
// Anchors does not parse the language. It reads the import's PATH with a pattern — the
// dialect's `import_pattern`, with its path in the group `path`, or the family's —, and
// resolves it to a file of the map with the dialect's `import_resolve`. An import of a
// package outside the project resolves to nothing, and is no dependency of the chain.

// RealImport is one import statement of a file, as the code writes it.
type RealImport struct {
	First   int    // 1-based, the statement's first line
	Line    int    // 1-based, the line that carries the path — where the flag goes
	Path    string // as written
	Target  string // the file it resolves to; empty when outside the project
	Dir     bool   // it resolves to a directory — a package of several files
	Symbols []string
}

// familyImportPattern is the import pattern a family knows, its path in the group `path`.
func familyImportPattern(family string) string {
	switch strings.ToLower(family) {
	case "ts", "js", "typescript", "javascript":
		return `(?:\bfrom\s+|^\s*import\s+|\brequire\(\s*|\bimport\(\s*)['"](?P<path>[^'"]+)['"]`
	case "go":
		return `^\s*(?:import\s+)?(?:[\w.]+\s+)?"(?P<path>[^"]+)"\s*$`
	}
	return ""
}

// familyExtensions are the extensions a family's import paths leave out.
func familyExtensions(family string) []string {
	switch strings.ToLower(family) {
	case "ts", "js", "typescript", "javascript":
		return []string{".ts", ".tsx", ".js", ".jsx", "/index.ts", "/index.tsx", "/index.js"}
	}
	return nil
}

// importPathRE is the pattern that reads an import's path: the project's `import_pattern`
// when it captures `path`, or the family's.
func importPathRE(d config.Dialect) *regexp.Regexp {
	if p := strings.TrimSpace(d.ImportPattern); p != "" {
		if re := d.Compile(p); re != nil && re.SubexpIndex("path") >= 0 {
			return re
		}
	}
	if p := familyImportPattern(d.Family); p != "" {
		return regexp.MustCompile(p)
	}
	return nil
}

var (
	importKeywordRE   = regexp.MustCompile(`^\s*import\b`)
	importBlockOpenRE = regexp.MustCompile(`^\s*import\s*\(\s*$`)
	// listOpenRE opens a list of names that spans lines: `import {`, `import type {`,
	// `import X, {`, `export {`.
	listOpenRE = regexp.MustCompile(`^(?:import|export)\b[^'"]*\{[^}]*$`)
)

// ImportsOf are the import statements of a file's content, resolved against the map's files.
// Without a pattern for the dialect, there is nothing to read.
func ImportsOf(content, rel string, d config.Dialect, g *mapx.Graph) []RealImport {
	re := importPathRE(d)
	if re == nil {
		return nil
	}
	group := re.SubexpIndex("path")
	lines := strings.Split(content, "\n")
	goBlocks := strings.EqualFold(d.Family, "go")
	inBlock := false
	var out []RealImport
	for i, ln := range lines {
		if goBlocks {
			// Go: a string alone on a line is an import only inside an import block, or
			// after the keyword — a slice literal of strings is not one.
			t := strings.TrimSpace(ln)
			if importBlockOpenRE.MatchString(ln) {
				inBlock = true
				continue
			}
			if inBlock && t == ")" {
				inBlock = false
				continue
			}
			if !inBlock && !importKeywordRE.MatchString(ln) {
				continue
			}
		}
		code := cutInlineComment(ln)
		m := re.FindStringSubmatch(code)
		if m == nil || m[group] == "" {
			continue
		}
		imp := RealImport{First: i + 1, Line: i + 1, Path: m[group]}
		// The statement's first line: only a line that CLOSES a list of names (`} from '…'`)
		// opens above, at the `import {` or `export {` that opened the list. A `require(…)`
		// or an `export … from` on its own line is a statement of its own — climbing from
		// them took the names and the flag of the import above.
		if !goBlocks && strings.HasPrefix(strings.TrimSpace(code), "}") {
			for j := i - 1; j >= 0 && i-j <= 60; j-- {
				t := strings.TrimSpace(cutInlineComment(lines[j]))
				if listOpenRE.MatchString(t) {
					imp.First = j + 1
					break
				}
				if strings.Contains(t, " from ") || strings.HasSuffix(t, ";") || t == "" {
					break // another statement: this one opens no list above
				}
			}
		}
		imp.Symbols = scan.ImportSymbols(strings.Join(lines[imp.First-1:imp.Line], " "))
		imp.Target, imp.Dir = resolveImport(rel, imp.Path, d, g)
		out = append(out, imp)
	}
	return out
}

// resolveImport turns an import's path into a file of the map: relative to the importer,
// or through an alias the project declares; then as written, or with each extension. A
// path that names a directory of governed files is a package. Anything else is outside
// the project.
func resolveImport(rel, p string, d config.Dialect, g *mapx.Graph) (string, bool) {
	if g == nil {
		return "", false
	}
	var base string
	switch {
	case strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") || p == "." || p == "..":
		base = path.Clean(path.Join(path.Dir(rel), p))
	default:
		var aliases map[string]string
		if d.ImportResolve != nil {
			aliases = d.ImportResolve.Aliases
		}
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
	exts := familyExtensions(d.Family)
	if d.ImportResolve != nil && len(d.ImportResolve.Extensions) > 0 {
		exts = d.ImportResolve.Extensions
	}
	if n := g.Node(base); n != nil {
		return base, false
	}
	for _, e := range exts {
		if n := g.Node(base + e); n != nil {
			return base + e, false
		}
	}
	for _, n := range g.Nodes {
		if path.Dir(n.ID) == base && n.Kind == mapx.KindCode {
			return base, true
		}
	}
	return "", false
}

// codesOfTarget are the codes a flag may name for an import: the target file's own code,
// or, for a package, the own code of any of its code files.
func codesOfTarget(imp RealImport, g *mapx.Graph) []string {
	if g == nil || imp.Target == "" {
		return nil
	}
	if !imp.Dir {
		if n := g.Node(imp.Target); n != nil && n.FileCode != "" {
			return []string{n.FileCode}
		}
		return nil
	}
	var out []string
	for _, n := range g.Nodes {
		if path.Dir(n.ID) == imp.Target && n.Kind == mapx.KindCode && n.FileCode != "" {
			out = append(out, n.FileCode)
		}
	}
	sort.Strings(out)
	return out
}

// flagsByStatement maps each import statement to the dependency flag on one of its lines.
func flagsByStatement(content string, imps []RealImport) map[int]*scan.CodeDep {
	byLine := map[int]scan.CodeDep{}
	for _, f := range scan.CodeDepsIn([]byte(content)) {
		byLine[f.Line] = f
	}
	out := map[int]*scan.CodeDep{}
	for i, imp := range imps {
		for l := imp.First; l <= imp.Line; l++ {
			if f, ok := byLine[l]; ok {
				f := f
				out[i] = &f
				break
			}
		}
	}
	return out
}

// chainUnit says whether a node takes part in the chain: code of the project, never a test
// (a test is tied by its `ref:`), a support file nor a vendored one.
func chainUnit(n mapx.Node) bool {
	return n.Kind == mapx.KindCode && !n.Support && !n.Upstream
}

// checkDepDeclared: does every import of a governed file carry `@dep:` with a code, or a
// waiver?
func checkDepDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if !chainUnit(n) {
		return Skip, i18n.T("gate.dep_chain.skip_not_code")
	}
	d := cfg.DialectFor()
	if importPathRE(d) == nil {
		return Pending, i18n.T("gate.dep_chain.pending_no_pattern")
	}
	imps := ImportsOf(content, n.ID, d, g)
	flags := flagsByStatement(content, imps)
	var missing []string
	for i, imp := range imps {
		codes := codesOfTarget(imp, g)
		if len(codes) == 0 || flags[i] != nil {
			continue
		}
		missing = append(missing, fmt.Sprintf("%d `%s` → %s", imp.Line, imp.Path, strings.Join(codes, "|")))
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.dep_chain.undeclared", len(missing), strings.Join(missing, "; "))
}

// checkDepHonored: does every `@dep:` name the code of the file its import resolves to?
func checkDepHonored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if !chainUnit(n) {
		return Skip, i18n.T("gate.dep_chain.skip_not_code")
	}
	all := scan.CodeDepsIn([]byte(content))
	if len(all) == 0 {
		return Skip, i18n.T("gate.dep_chain.skip_no_flags")
	}
	d := cfg.DialectFor()
	imps := ImportsOf(content, n.ID, d, g)
	var wrong []string
	for _, f := range all {
		if f.Code == "" {
			continue // a waiver names no code
		}
		var imp *RealImport
		for i := range imps {
			if f.Line >= imps[i].First && f.Line <= imps[i].Line {
				imp = &imps[i]
				break
			}
		}
		if imp == nil {
			wrong = append(wrong, i18n.T("gate.dep_chain.flag_without_import", f.Line, f.Code))
			continue
		}
		codes := codesOfTarget(*imp, g)
		if !contains(codes, f.Code) {
			want := strings.Join(codes, "|")
			if want == "" {
				want = "—"
			}
			wrong = append(wrong, i18n.T("gate.dep_chain.flag_wrong_code", f.Line, f.Code, imp.Path, want))
		}
	}
	if len(wrong) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(wrong, "; ")
}

// importerIndex is, per file, who imports which of its symbols — read once per run.
type importerIndex struct {
	g     *mapx.Graph
	usage map[string]map[string]map[string]bool // target → symbol → importer codes
}

var (
	importersMu    sync.Mutex
	importersCache = map[string]*importerIndex{}
)

// importersOf is the index of the imports of the project's code, built once per map.
func importersOf(root string, g *mapx.Graph, d config.Dialect) *importerIndex {
	importersMu.Lock()
	defer importersMu.Unlock()
	if idx, ok := importersCache[root]; ok && idx.g == g {
		return idx
	}
	idx := &importerIndex{g: g, usage: map[string]map[string]map[string]bool{}}
	for _, n := range g.Nodes {
		if !chainUnit(n) || n.FileCode == "" {
			continue
		}
		b, err := readFile(root, n.ID)
		if err != nil {
			continue
		}
		for _, imp := range ImportsOf(string(b), n.ID, d, g) {
			if imp.Target == "" || imp.Dir {
				continue
			}
			for _, s := range imp.Symbols {
				if idx.usage[imp.Target] == nil {
					idx.usage[imp.Target] = map[string]map[string]bool{}
				}
				if idx.usage[imp.Target][s] == nil {
					idx.usage[imp.Target][s] = map[string]bool{}
				}
				idx.usage[imp.Target][s][n.FileCode] = true
			}
		}
	}
	importersCache[root] = idx
	return idx
}

// UsedByOf is, for a file, who imports each of its symbols — by the importers' codes,
// sorted.
func UsedByOf(root, id string, g *mapx.Graph, d config.Dialect) map[string][]string {
	out := map[string][]string{}
	for sym, codes := range importersOf(root, g, d).usage[id] {
		for c := range codes {
			out[sym] = append(out[sym], c)
		}
		sort.Strings(out[sym])
	}
	return out
}

// checkUsedByDeclared: does every symbol another file imports carry `@used-by:` with
// exactly the codes of who imports it?
func checkUsedByDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if !chainUnit(n) {
		return Skip, i18n.T("gate.dep_chain.skip_not_code")
	}
	d := cfg.DialectFor()
	if importPathRE(d) == nil || g == nil {
		return Pending, i18n.T("gate.dep_chain.pending_no_pattern")
	}
	want := UsedByOf(root, n.ID, g, d)
	have := map[string][]string{}
	for _, u := range scan.UsedByIn([]byte(content)) {
		if u.Symbol != "" {
			codes := append([]string(nil), u.Codes...)
			sort.Strings(codes)
			have[u.Symbol] = codes
		}
	}
	var gaps []string
	syms := make([]string, 0, len(want))
	for s := range want {
		syms = append(syms, s)
	}
	sort.Strings(syms)
	for _, s := range syms {
		got, ok := have[s]
		switch {
		case !ok:
			gaps = append(gaps, i18n.T("gate.dep_chain.used_by_missing", s, strings.Join(want[s], ", ")))
		case strings.Join(got, ",") != strings.Join(want[s], ","):
			gaps = append(gaps, i18n.T("gate.dep_chain.used_by_wrong", s, strings.Join(got, ", "), strings.Join(want[s], ", ")))
		}
	}
	for s := range have {
		if _, ok := want[s]; !ok {
			gaps = append(gaps, i18n.T("gate.dep_chain.used_by_stale", s))
		}
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	sort.Strings(gaps)
	return Fail, strings.Join(gaps, "; ")
}

// fixDepFlags writes the `@dep:` of each import of a governed file that has none, and
// corrects the code of one that names another file — the mechanical half of the chain.
// An import of a package of several files is left as it is: which file it uses is the
// author's to say.
func fixDepFlags(content string, n mapx.Node, root string, g *mapx.Graph) (string, bool, string) {
	if !chainUnit(n) || g == nil {
		return content, false, ""
	}
	d := fixConfig.DialectFor()
	imps := ImportsOf(content, n.ID, d, g)
	if len(imps) == 0 {
		return content, false, ""
	}
	flags := flagsByStatement(content, imps)
	lines := strings.Split(content, "\n")
	marker := config.LineCommentFor(n.ID)
	var done []string
	for i, imp := range imps {
		codes := codesOfTarget(imp, g)
		if len(codes) != 1 {
			continue
		}
		want := codes[0]
		f := flags[i]
		switch {
		case f == nil:
			lines[imp.Line-1] = strings.TrimRight(lines[imp.Line-1], " \t") + " " + marker + " @dep: " + want
			done = append(done, fmt.Sprintf("%d %s", imp.Line, want))
		case f.Code != "" && f.Code != want:
			l := lines[f.Line-1]
			lines[f.Line-1] = strings.Replace(l, "@dep: "+f.Code, "@dep: "+want, 1)
			done = append(done, fmt.Sprintf("%d %s→%s", f.Line, f.Code, want))
		}
	}
	if len(done) == 0 {
		return content, false, ""
	}
	return strings.Join(lines, "\n"), true, i18n.T("gate.fix.dep_written", strings.Join(done, ", "))
}

// fixUsedBy writes, above each symbol another file imports, the `@used-by:` with exactly
// who imports it, replacing a flag that says otherwise; a flag on a symbol nobody imports
// is removed.
func fixUsedBy(content string, n mapx.Node, root string, g *mapx.Graph) (string, bool, string) {
	if !chainUnit(n) || g == nil {
		return content, false, ""
	}
	d := fixConfig.DialectFor()
	if importPathRE(d) == nil {
		return content, false, ""
	}
	want := UsedByOf(root, n.ID, g, d)
	lines := strings.Split(content, "\n")
	flagLine := map[string]int{} // symbol → 0-based line of its flag
	for _, u := range scan.UsedByIn([]byte(content)) {
		if u.Symbol != "" {
			flagLine[u.Symbol] = u.Line - 1
		}
	}
	marker := config.LineCommentFor(n.ID)
	type edit struct {
		at     int
		insert bool
		text   string
		remove bool
	}
	var edits []edit
	var done []string
	for sym, codes := range want {
		text := marker + " @used-by: " + strings.Join(codes, ", ")
		if at, ok := flagLine[sym]; ok {
			indent := lines[at][:len(lines[at])-len(strings.TrimLeft(lines[at], " \t"))]
			if strings.TrimSpace(lines[at]) != text {
				edits = append(edits, edit{at: at, text: indent + text})
				done = append(done, sym)
			}
			continue
		}
		at := declarationLine(lines, sym)
		if at < 0 {
			continue
		}
		indent := lines[at][:len(lines[at])-len(strings.TrimLeft(lines[at], " \t"))]
		edits = append(edits, edit{at: at, insert: true, text: indent + text})
		done = append(done, sym)
	}
	for sym, at := range flagLine {
		if _, ok := want[sym]; !ok {
			edits = append(edits, edit{at: at, remove: true})
			done = append(done, "-"+sym)
		}
	}
	if len(edits) == 0 {
		return content, false, ""
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].at > edits[j].at })
	for _, e := range edits {
		switch {
		case e.remove:
			lines = append(lines[:e.at], lines[e.at+1:]...)
		case e.insert:
			lines = append(lines[:e.at], append([]string{e.text}, lines[e.at:]...)...)
		default:
			lines[e.at] = e.text
		}
	}
	sort.Strings(done)
	return strings.Join(lines, "\n"), true, i18n.T("gate.fix.used_by_written", strings.Join(done, ", "))
}

// declarationLine is the 0-based line that declares a symbol at the top level of the file,
// or -1.
func declarationLine(lines []string, sym string) int {
	if sym == "default" {
		for i, l := range lines {
			if defaultDeclRE.MatchString(l) {
				return i
			}
		}
		return -1
	}
	re := regexp.MustCompile(`^(?:export\s+)?(?:default\s+)?(?:async\s+)?(?:function\*?|const|let|var|class|interface|type|enum|func|def)\s+(?:\([^)]*\)\s*)?` + regexp.QuoteMeta(sym) + `\b`)
	for i, l := range lines {
		if re.MatchString(l) {
			return i
		}
	}
	return -1
}

// defaultDeclRE is the line that declares a module's default export.
var defaultDeclRE = regexp.MustCompile(`^\s*export\s+default\b|^\s*module\.exports\s*=`)

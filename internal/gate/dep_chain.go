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
// Anchors does not parse the language. It reads each import statement with the project's
// patterns — the dialect's `imports`, the path in the group `path`; `import_names` for what
// it brings —, and resolves the path to a file of the map with the dialect's
// `import_resolve`. The core knows no syntax of an import (DESIGN-chain-patterns.md). An import of a
// package outside the project resolves to nothing, and is no dependency of the chain.

// RealImport is one import statement of a file, as the code writes it.
type RealImport struct {
	First   int    // 1-based, the statement's first line
	Line    int    // 1-based, the line that carries the path — where the flag goes
	Path    string // as written
	Target  string // the file it resolves to; empty when outside the project
	Dir     bool   // it resolves to a directory — a package of several files
	Symbols []string
	// TypeOnly: the statement brings types only (the dialect's type_import_pattern).
	TypeOnly bool
	// Relative: the pattern that read it says its path is relative to the importer.
	Relative bool
	// Runs: the pattern that read it says the statement plainly runs.
	Runs bool
}

// importResolve is the project's import resolution: nil when it declares none.
func importResolve(d config.Dialect) *config.ImportResolve { return d.ImportResolve }

// byName says whether the dialect's imports are module names read anywhere in the code,
// flagged per file (config.ImportResolve.ByName).
func byName(d config.Dialect) bool {
	r := importResolve(d)
	return r != nil && r.ByName
}

// fileName is a module name as its file is named.
func fileName(name, nameCase string) string {
	switch nameCase {
	case "lower_first":
		if name == "" {
			return name
		}
		return strings.ToLower(name[:1]) + name[1:]
	case "lower":
		return strings.ToLower(name)
	}
	return name
}

// blockCommentRE is a block comment that opens and closes on one line, in the languages
// whose only comment is a block: OCaml's.
var blockCommentRE = regexp.MustCompile(`\(\*.*?\*\)`)

// importsByName are the modules a file names, one import per module, at its first mention,
// each resolved to the code file named after it. The file's comments name nothing.
func importsByName(content, rel string, d config.Dialect, g *mapx.Graph) []RealImport {
	r := importResolve(d)
	seen := map[string]bool{}
	var out []RealImport
	depth := 0
	for i, ln := range strings.Split(content, "\n") {
		ln = blockCommentRE.ReplaceAllString(ln, " ")
		if depth > 0 {
			if j := strings.Index(ln, "*)"); j >= 0 {
				depth, ln = 0, ln[j+2:]
			} else {
				continue
			}
		}
		if j := strings.Index(ln, "(*"); j >= 0 {
			depth, ln = 1, ln[:j]
		}
		for _, m := range d.Imports.FindAll(ln, rel) {
			name := m.Groups["path"].Value
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			target := resolveByName(rel, name, r, g)
			out = append(out, RealImport{First: i + 1, Line: i + 1, Path: name, Target: target})
		}
	}
	return out
}

// resolveByName is the code file named after the module: the one beside the importer when
// there is one, or the only one in the project. A name several files carry, none of them
// beside the importer, is not told apart — no target rather than a wrong one.
func resolveByName(rel, name string, r *config.ImportResolve, g *mapx.Graph) string {
	if g == nil || r == nil {
		return ""
	}
	exts := r.Extensions
	if len(exts) == 0 {
		exts = []string{""}
	}
	stem := fileName(name, r.NameCase)
	var near, all []string
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode || n.ID == rel {
			continue
		}
		for _, e := range exts {
			if path.Base(n.ID) == stem+e {
				all = append(all, n.ID)
				if path.Dir(n.ID) == path.Dir(rel) {
					near = append(near, n.ID)
				}
			}
		}
	}
	switch {
	case len(near) == 1:
		return near[0]
	case len(all) == 1:
		return all[0]
	}
	return ""
}

// hasImports says whether the project declares how its imports are written: without it
// the chain gates have nothing to read the code with (DESIGN-chain-patterns.md).
func hasImports(d config.Dialect) bool { return len(d.Imports) > 0 }

// blankComments is the content with each comment — a line of comment, or what follows the
// code on a line — turned into spaces, its offsets and lines kept: an import commented out
// is no import.
func blankComments(content string) string {
	lines := strings.Split(content, "\n")
	for i, ln := range lines {
		switch {
		case commentLine(ln):
			lines[i] = strings.Repeat(" ", len(ln))
		default:
			if code := cutInlineComment(ln); len(code) < len(ln) {
				lines[i] = code + strings.Repeat(" ", len(ln)-len(code))
			}
		}
	}
	return strings.Join(lines, "\n")
}

// ImportsOf are the import statements of a file's content, resolved against the map's files.
// Each is a match of one of the dialect's `imports` patterns: the match is the statement, its
// `path` group the path and the line its flag goes on. Without patterns there is nothing to
// read.
func ImportsOf(content, rel string, d config.Dialect, g *mapx.Graph) []RealImport {
	if !hasImports(d) {
		return nil
	}
	if byName(d) {
		return importsByName(content, rel, d, g)
	}
	var typeRE *regexp.Regexp
	if p := strings.TrimSpace(d.TypeImportPattern); p != "" {
		typeRE = d.Compile(p)
	}
	code := blankComments(content)
	seen := map[int]bool{} // a path captured by two patterns is one import
	var out []RealImport
	for _, m := range d.Imports.FindAll(code, rel) {
		p, ok := m.Groups["path"]
		if !ok || p.Value == "" || seen[p.Start] {
			continue
		}
		seen[p.Start] = true
		statement := code[m.Start:m.End]
		imp := RealImport{
			First:    config.LineOf(code, m.Start+len(statement)-len(strings.TrimLeft(statement, " \t\n"))),
			Line:     config.LineOf(code, p.Start),
			Path:     p.Value,
			Relative: m.Relative,
			Runs:     m.Runs,
			Symbols:  d.NamesBrought(statement, rel),
			TypeOnly: typeRE != nil && typeRE.MatchString(statement),
		}
		imp.Target, imp.Dir = resolveImport(rel, imp, d, g)
		out = append(out, imp)
	}
	return out
}

// resolveImport turns an import's path into a file of the map: relative to the importer,
// or through an alias the project declares; then as written, or with each extension. A
// path that names a directory of governed files is a package. Anything else is outside
// the project.
func resolveImport(rel string, imp RealImport, d config.Dialect, g *mapx.Graph) (string, bool) {
	if g == nil {
		return "", false
	}
	p := imp.Path
	var base string
	switch {
	// A pattern may say its paths are relative to the importer without `./` — a flow names
	// the flow it runs that way.
	case strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") || p == "." || p == ".." || imp.Relative:
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
	var exts []string
	if d.ImportResolve != nil {
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

// chainUnit says whether a node takes part in the chain: every artifact of the project that
// imports — code, tests, test support, flows —, never a vendored one. A test imports what it
// tests, its fixtures and its helpers, and a moved fixture breaks it as a moved hook breaks a
// screen (DESIGN-dependencies-out-of-the-spec.md, DOOSD-D01).
func chainUnit(n mapx.Node) bool {
	return (n.Kind == mapx.KindCode || n.Kind == mapx.KindTest) && !n.Upstream
}

// checkDepDeclared: does every import of a governed file carry `@dep:` with a code, or a
// waiver?
func checkDepDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if !chainUnit(n) {
		return Skip, i18n.T("gate.dep_chain.skip_not_code")
	}
	d := cfg.DialectFor()
	if !hasImports(d) {
		return Pending, i18n.T("gate.dep_chain.pending_no_pattern")
	}
	imps := ImportsOf(content, n.ID, d, g)
	flags := flagsByStatement(content, imps)
	// By name, the flags are the file's: each module it uses is flagged once, anywhere.
	fileFlags := map[string]bool{}
	if byName(d) {
		for _, f := range scan.CodeDepsIn([]byte(content)) {
			fileFlags[f.Code] = true
		}
	}
	var missing []string
	for i, imp := range imps {
		codes := codesOfTarget(imp, g)
		if len(codes) == 0 || flags[i] != nil || anyFlagged(codes, fileFlags) {
			continue
		}
		missing = append(missing, fmt.Sprintf("%d `%s` → %s", imp.Line, imp.Path, strings.Join(codes, "|")))
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.dep_chain.undeclared", len(missing), strings.Join(missing, "; ")) + i18n.T("gate.dep_chain.guide")
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
		if f.Kind != "" && f.Kind != scan.TypeDepKind {
			// What the code reaches that is no import: the kind must be one the project
			// declares; the call is the author's to name.
			if !contains(cfg.DependencyKinds, f.Kind) {
				wrong = append(wrong, i18n.T("gate.dep_chain.kind_undeclared", f.Line, f.Kind, f.Name))
			}
			continue
		}
		if f.Code == "" {
			continue // a waiver names no code
		}
		if byName(d) {
			if !usesCode(imps, g, f.Code) {
				wrong = append(wrong, i18n.T("gate.dep_chain.flag_unused_module", f.Line, f.Code))
			}
			continue
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
	return Fail, strings.Join(wrong, "; ") + i18n.T("gate.dep_chain.guide")
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
func importersOf(root string, g *mapx.Graph, d config.Dialect, withTests bool) *importerIndex {
	importersMu.Lock()
	defer importersMu.Unlock()
	key := fmt.Sprint(root, "\x00", withTests)
	if idx, ok := importersCache[key]; ok && idx.g == g {
		return idx
	}
	idx := &importerIndex{g: g, usage: map[string]map[string]map[string]bool{}}
	for _, n := range g.Nodes {
		if !chainUnit(n) || n.FileCode == "" || n.Kind == mapx.KindTest && !withTests {
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
	importersCache[key] = idx
	return idx
}

// testsInChain says whether the project put its tests in the chain: the used-by gate
// confronts tests (`on:` holds `test`). Then a symbol a test imports lists the test; with the
// gate on code alone, @used-by lists the code that imports it, as before the tests could take
// part — adopting the chain on the tests is the project's decision, not a version's (reported
// from MIF: 768 symbols asked for their tests after an upgrade).
func testsInChain(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	for _, g := range cfg.Gates {
		if g.Check == "used-by-declared" || g.Name == "used-by-declared" {
			for _, k := range g.On {
				if k == "test" {
					return true
				}
			}
		}
	}
	return false
}

// UsedByOf is, for a file, who imports each of its symbols — by the importers' codes,
// sorted.
func UsedByOf(root, id string, g *mapx.Graph, d config.Dialect, withTests bool) map[string][]string {
	out := map[string][]string{}
	for sym, codes := range importersOf(root, g, d, withTests).usage[id] {
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
	if !hasImports(d) || g == nil {
		return Pending, i18n.T("gate.dep_chain.pending_no_pattern")
	}
	want := UsedByOf(root, n.ID, g, d, testsInChain(cfg))
	have := map[string][]string{}
	for _, u := range scan.UsedByIn([]byte(content), n.ID, d) {
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
	return Fail, strings.Join(gaps, "; ") + i18n.T("gate.dep_chain.guide")
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
	if byName(d) {
		return fixDepFlagsByName(content, n, imps, g)
	}
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
		// An import of types only is flagged `@dep[type]:` — nothing of it runs, and the
		// project's evidence may say so (`evidence.no_side_effect.flags`).
		flag := "@dep: "
		f := flags[i]
		// A type flag the author wrote stays unless the statement plainly runs (its pattern
		// says `runs`): an inline `import('…').T` in a type position reads like any other
		// import, and only the author knows which it is.
		if imp.TypeOnly || (f != nil && f.Kind == scan.TypeDepKind && !imp.Runs) {
			flag = "@dep[" + scan.TypeDepKind + "]: "
		}
		switch {
		case f == nil:
			lines[imp.Line-1] = strings.TrimRight(lines[imp.Line-1], " \t") + " " + marker + " " + flag + want + config.CommentCloseFor(marker)
			done = append(done, fmt.Sprintf("%d %s", imp.Line, want))
		case f.Code != "" && (f.Code != want || (imp.TypeOnly && f.Kind != scan.TypeDepKind) ||
			(!imp.TypeOnly && f.Kind == scan.TypeDepKind && imp.Runs)):
			l := lines[f.Line-1]
			had := "@dep: " + f.Code
			if f.Kind != "" {
				had = "@dep[" + f.Kind + "]: " + f.Code
			}
			lines[f.Line-1] = strings.Replace(l, had, flag+want, 1)
			done = append(done, fmt.Sprintf("%d %s→%s%s", f.Line, f.Code, strings.TrimSuffix(flag, " "), want))
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
	if !hasImports(d) {
		return content, false, ""
	}
	want := UsedByOf(root, n.ID, g, d, testsInChain(fixConfig))
	lines := strings.Split(content, "\n")
	flagLine := map[string]int{} // symbol → 0-based line of its flag
	for _, u := range scan.UsedByIn([]byte(content), n.ID, d) {
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
	syms := make([]string, 0, len(want))
	for sym := range want {
		syms = append(syms, sym)
	}
	sort.Strings(syms)
	for _, sym := range syms {
		codes := want[sym]
		at, listed := declarationLine(content, n.ID, d, sym)
		text := marker + " @used-by: " + strings.Join(codes, ", ") + config.CommentCloseFor(marker)
		if listed {
			text += " (" + sym + ")"
		}
		if at, ok := flagLine[sym]; ok {
			indent := lines[at][:len(lines[at])-len(strings.TrimLeft(lines[at], " \t"))]
			if strings.TrimSpace(lines[at]) != text {
				edits = append(edits, edit{at: at, text: indent + text})
				done = append(done, sym)
			}
			continue
		}
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
	// From the bottom up; flags inserted above the same line land in the order of their text.
	sort.SliceStable(edits, func(i, j int) bool {
		if edits[i].at != edits[j].at {
			return edits[i].at > edits[j].at
		}
		return edits[i].text > edits[j].text
	})
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

// declarationLine is the 0-based line that declares a symbol, by the dialect's
// `declarations` — where the statement of the declaration starts —, or -1; listed when the
// declaration is one of a list, where a flag must name its symbol.
func declarationLine(content, rel string, d config.Dialect, sym string) (int, bool) {
	code := blankComments(content)
	for _, m := range d.Declarations.FindAll(code, rel) {
		if (sym == "default" && m.Has("default")) || (sym != "default" && m.Groups["name"].Value == sym) {
			return config.LineOf(code, m.Region) - 1, m.Listed
		}
	}
	return -1, false
}

// anyFlagged says whether one of the codes is in the set.
func anyFlagged(codes []string, set map[string]bool) bool {
	for _, c := range codes {
		if set[c] {
			return true
		}
	}
	return false
}

// usesCode says whether one of the file's imports resolves to a file of that code.
func usesCode(imps []RealImport, g *mapx.Graph, code string) bool {
	for _, imp := range imps {
		if contains(codesOfTarget(imp, g), code) {
			return true
		}
	}
	return false
}

// fixDepFlagsByName writes, below the file's header, a flag line for each module the file
// uses that none flags, and removes a flag line whose module it no longer uses: by name, the
// flags are the file's, one per line, in its comment.
func fixDepFlagsByName(content string, n mapx.Node, imps []RealImport, g *mapx.Graph) (string, bool, string) {
	want := map[string]bool{}
	for _, imp := range imps {
		if codes := codesOfTarget(imp, g); len(codes) == 1 {
			want[codes[0]] = true
		}
	}
	have := map[string]bool{}
	lines := strings.Split(content, "\n")
	var kept []string
	var done []string
	for i, l := range lines {
		if scan.ChainFlagOnlyRE.MatchString(l) {
			if fs := scan.CodeDepsIn([]byte(l)); len(fs) == 1 && fs[0].Code != "" && fs[0].Kind == "" {
				if !want[fs[0].Code] {
					done = append(done, fmt.Sprintf("%d -%s", i+1, fs[0].Code))
					continue
				}
				have[fs[0].Code] = true
			}
		}
		kept = append(kept, l)
	}
	for _, f := range scan.CodeDepsIn([]byte(content)) {
		have[f.Code] = true
	}
	var add []string
	for c := range want {
		if !have[c] {
			add = append(add, c)
		}
	}
	sort.Strings(add)
	if len(add) == 0 && len(done) == 0 {
		return content, false, ""
	}
	marker := config.LineCommentFor(n.ID)
	var flagLines []string
	for _, c := range add {
		flagLines = append(flagLines, marker+" @dep: "+c+config.CommentCloseFor(marker))
		done = append(done, c)
	}
	at := 0
	if block := scan.AnchorsHeader([]byte(strings.Join(kept, "\n"))); block != nil {
		at = strings.Count(strings.SplitN(strings.Join(kept, "\n"), string(block), 2)[0], "\n") + strings.Count(string(block), "\n") + 1
	}
	out := append(append(append([]string{}, kept[:at]...), flagLines...), kept[at:]...)
	return strings.Join(out, "\n"), true, i18n.T("gate.fix.dep_written", strings.Join(done, ", "))
}

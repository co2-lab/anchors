// @anchors
//   code: MSGGM
//   ref: MKSTP

package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// --- writing the stamps `mock-stamped` checks ---
//
// A project adopting `mock-stamped` has every double of a governed module unstamped at
// once — measured in the project that asked for this: 1278 `jest.mock`/`vi.mock` calls,
// none stamped. Stamping by hand is not viable, and a stamp written by hand is also the
// one most likely to point at the wrong snippet.
//
// THE GENERATOR ONLY WRITES STAMPS THAT ARE MISSING. It never rewrites an existing one,
// and that is the whole point of the mechanism: the gate recomputes the stamp against the
// real module, and a divergence says the double may have drifted. A tool that refreshed
// divergent stamps would let whoever edits the test regenerate the stamp to match their
// own mock — the stamp would certify itself, which the gate's own comment names as the
// thing it exists to prevent. When a stamp diverges, a person looks at the double.
//
// Everything it writes is built with the gate's own functions (`declaredStamps`,
// `snippetHash`, `moduleHasStamp`), so a stamp it writes is exactly the stamp the gate
// recomputes.

// StampWritten is one stamp the generator added.
type StampWritten struct {
	Module string // the double's specifier, as written in the test
	File   string // the real module, relative to the root
	Anchor string
	Count  int
	Hash   string
}

// StampSkipped is a double the generator did not stamp, and why.
type StampSkipped struct {
	Module string
	Reason string
}

// GenerateStamps returns the test content with a stamp above every double of a governed
// module that has none, plus what it wrote and what it had to skip.
//
// Per double:
//   - the module is resolved to ONE file of the map; when the path suffix matches more
//     than one, the file sharing the longest directory prefix with the test wins, and a
//     tie is skipped — guessing would stamp a file the test never imports;
//   - each EXPORT of the module (by the project's `export_detect`) the double's call names
//     gets its own stamp, anchored on that export's line and
//     covering its block — so a change to that member is what makes the stamp diverge;
//   - with no export named in the call, one stamp covers the whole module.
func GenerateStamps(content, testID, root string, g *mapx.Graph, cfg *config.Config) (string, []StampWritten, []StampSkipped, error) {
	detector, err := doubleDetector(cfg)
	if err != nil {
		return content, nil, nil, err
	}
	if detector == nil {
		return content, nil, nil, fmt.Errorf("the project does not declare `derived.mock_detect` — without it no double can be found")
	}
	exportRE := exportDetectDe(cfg)

	existing := declaredStamps(content)
	lines := strings.Split(content, "\n")
	lineOf := lineIndexer(content)
	comment := cfg.CommentsFor(testID).LineMarker()
	if comment == "" {
		return content, nil, nil, fmt.Errorf("the project declares no comments for %s (`dialect.comments`): a stamp is written in one", testID)
	}

	var written []StampWritten
	var skipped []StampSkipped
	inserts := map[int][]string{} // line index → stamp lines to insert above it
	done := map[string]bool{}
	// The test as the call is read: its strings and comments blanked, as its syntax says.
	blanked := cfg.CommentsFor(testID).Blank(content)

	for _, m := range detector.FindAllStringSubmatchIndex(content, -1) {
		if len(m) < 4 || m[2] < 0 {
			continue
		}
		module := strings.TrimSpace(content[m[2]:m[3]])
		if module == "" || done[module] || !isGovernedModule(module, testID, cfg.DialectFor(), g) {
			continue
		}
		done[module] = true
		// A stamped double is never rewritten — see the file comment. It can still GAIN a
		// stamp: a key added to its factory later names an export no stamp covers.
		stamped := moduleHasStamp(module, testID, cfg.DialectFor(), g, existing)

		file, reason := moduleFile(module, testID, cfg.DialectFor(), g)
		if file == "" {
			skipped = append(skipped, StampSkipped{module, reason})
			continue
		}
		body, err := readFile(root, file)
		if err != nil {
			skipped = append(skipped, StampSkipped{module, "the module could not be read: " + file})
			continue
		}
		// The module as the stamp reads it: without the chain's flags (see stampLines).
		modCS := cfg.CommentsFor(file)
		modLines := stampLines(strings.Split(string(body), "\n"), modCS)
		starts := declarationStarts(strings.Join(modLines, "\n"), file, cfg.DialectFor())

		var stamps []StampWritten
		if exportRE != nil {
			for _, key := range factoryNames(blanked, m[0], exportNames(modLines, exportRE)) {
				if s, ok := stampForExport(module, file, key, modLines, exportRE, starts); ok {
					if stamped && coveredByExisting(s.Anchor, file, modLines, existing) {
						continue
					}
					stamps = append(stamps, s)
				}
			}
		}
		if stamped && len(stamps) == 0 {
			continue
		}
		if len(stamps) == 0 {
			s, ok := stampForWholeModule(module, file, modLines, modCS)
			if !ok {
				skipped = append(skipped, StampSkipped{module, "no line of the module can anchor a stamp unambiguously"})
				continue
			}
			stamps = append(stamps, s)
		}

		at := lineOf(m[0])
		indent := leadingSpace(lines[at])
		for _, s := range stamps {
			inserts[at] = append(inserts[at], fmt.Sprintf("%s%s @contract: %s | %s | %d | %s",
				indent, comment, s.File, s.Anchor, s.Count, s.Hash))
			written = append(written, s)
		}
	}

	if len(inserts) == 0 {
		return content, written, skipped, nil
	}
	var out []string
	for i, l := range lines {
		out = append(out, inserts[i]...)
		out = append(out, l)
	}
	return strings.Join(out, "\n"), written, skipped, nil
}

// coveredByExisting reports whether an existing stamp on `file` already covers the line
// `anchor` — the same anchor, or a range (a whole-module stamp, a wider block) that
// contains it. Such a member needs no new stamp.
func coveredByExisting(anchor, file string, modLines []string, existing []declaredStamp) bool {
	at := -1
	for i, l := range modLines {
		if strings.TrimRight(l, "\r") == anchor {
			at = i
			break
		}
	}
	for _, c := range existing {
		if withoutExtension(c.file) != withoutExtension(file) {
			continue
		}
		if c.anchor == anchor {
			return true
		}
		for i, l := range modLines {
			if strings.TrimRight(l, "\r") == c.anchor && at >= i && at < i+c.count {
				return true
			}
		}
	}
	return false
}

// moduleFile finds the ONE code file of the map a double's specifier names. The specifier
// is an import path: resolved as the chain resolves one (`import_resolve` — relative paths,
// the project's aliases, its extensions); a path that does not resolve is matched by its
// tail against the map's code files, past its `./` and `../` — the test's folder deciding
// between several. A specifier naming no file of the project is a third party's.
func moduleFile(spec, testID string, d config.Dialect, g *mapx.Graph) (string, string) {
	if g == nil {
		return "", "no code file of the map matches the specifier"
	}
	file, cands := g.Modules().File(spec, testID, d)
	switch {
	case file != "":
		return file, ""
	case len(cands) == 0:
		return "", "no code file of the map matches the specifier"
	}
	return "", fmt.Sprintf("the specifier matches %d files and the test's folder does not decide: %s", len(cands), strings.Join(cands, ", "))
}

// moduleCandidates are the code files of the map a double's specifier may name: the one its
// import resolution gives, or every one whose path ends as the specifier does.
func moduleCandidates(spec, testID string, d config.Dialect, g *mapx.Graph) []string {
	if g == nil {
		return nil
	}
	return g.Modules().Candidates(spec, testID, d)
}

// factoryNames are the module's exports the double's call names, in the order it names them
// — from the `(` after the match to the parenthesis that closes it, in a test whose strings
// and comments are blanked. A word is read as a member of the double only when it is the
// name of an export of the real module: the core reads no object literal of a language.
func factoryNames(blanked string, from int, exports []string) []string {
	open := strings.Index(blanked[from:], "(")
	if open < 0 {
		return nil
	}
	start := from + open
	depth, end := 0, -1
	for i := start; i < len(blanked) && end < 0; i++ {
		switch blanked[i] {
		case '(':
			depth++
		case ')':
			if depth--; depth == 0 {
				end = i
			}
		}
	}
	if end < 0 {
		return nil
	}
	call := blanked[start:end]
	type at struct {
		name string
		pos  int
	}
	var found []at
	for _, e := range exports {
		if loc := regexp.MustCompile(`(^|[^\p{L}\p{N}_$])` + regexp.QuoteMeta(e) + `($|[^\p{L}\p{N}_$])`).FindStringIndex(call); loc != nil {
			found = append(found, at{e, loc[0]})
		}
	}
	sort.SliceStable(found, func(i, j int) bool { return found[i].pos < found[j].pos })
	out := make([]string, len(found))
	for i, f := range found {
		out[i] = f.name
	}
	return out
}

// exportNames are the names the module exports, by the project's export pattern read line by
// line — as stampForExport anchors them —, each once.
func exportNames(lines []string, exportRE *regexp.Regexp) []string {
	seen := map[string]bool{}
	var out []string
	for _, l := range lines {
		for _, m := range exportRE.FindAllStringSubmatch(l, -1) {
			for _, g := range m[1:] {
				if g != "" && !seen[g] {
					seen[g] = true
					out = append(out, g)
					break
				}
			}
		}
	}
	return out
}

// stampForExport anchors a stamp on the line that exports `key`, covering its block.
func stampForExport(module, file, key string, lines []string, exportRE *regexp.Regexp, starts map[int]bool) (StampWritten, bool) {
	for i, l := range lines {
		m := exportRE.FindStringSubmatch(l)
		if len(m) < 2 || m[1] != key {
			continue
		}
		if !usableAnchor(lines, i) {
			return StampWritten{}, false
		}
		n := blockLength(lines, i, starts)
		return StampWritten{Module: module, File: file, Anchor: lines[i], Count: n,
			Hash: snippetHash(strings.Join(lines[i:i+n], "\n"))}, true
	}
	return StampWritten{}, false
}

// stampForWholeModule anchors a stamp on the first usable line AFTER the `@anchors`
// header and covers to the end.
//
// The header stays out of the hash because `check --fix` rewrites its `updated_at:` on
// every edit: a stamp covering it diverged on ANY change to the module, including one that
// never touched what the double replaces — measured in the project that adopted the
// generator, where a 1037-line module flagged its doubles on every edit.
func stampForWholeModule(module, file string, lines []string, cs *config.Comments) (StampWritten, bool) {
	for i := headerEnd(lines, cs); i < len(lines); i++ {
		if !usableAnchor(lines, i) {
			continue
		}
		n := len(lines) - i
		return StampWritten{Module: module, File: file, Anchor: lines[i], Count: n,
			Hash: snippetHash(strings.Join(lines[i:], "\n"))}, true
	}
	return StampWritten{}, false
}

// headerSearchLines bounds where the `@anchors` header may start: it sits at the top, after
// at most a directive (`'use client'`) or a shebang. A mention further down is prose.
const headerSearchLines = 10

// headerEnd returns the index of the first line after the file's `@anchors` header, as the
// map reads it (scan.HeaderRange), or 0 when the file has none.
func headerEnd(lines []string, cs *config.Comments) int {
	if _, to, ok := scan.HeaderRange(lines, cs); ok {
		return to
	}
	return 0
}

// usableAnchor: the line can anchor a stamp the gate will find again.
//
// It must be non-blank, carry no surrounding whitespace (the stamp's pattern trims it, and
// the gate compares the line exactly), contain no `|` (the stamp's separator), and occur
// ONCE in the module — the gate refuses an ambiguous anchor, and would fail the stamp.
func usableAnchor(lines []string, i int) bool {
	l := strings.TrimRight(lines[i], "\r")
	if strings.TrimSpace(l) == "" || strings.TrimSpace(l) != l || strings.Contains(l, "|") {
		return false
	}
	count := 0
	for _, other := range lines {
		if strings.TrimRight(other, "\r") == l {
			count++
		}
	}
	return count == 1
}

// declarationStarts are the lines where a declaration of the module starts, by the dialect's
// `declarations` (where each statement starts): the next one ends the block before it.
func declarationStarts(content, rel string, d config.Dialect) map[int]bool {
	out := map[int]bool{}
	for _, m := range d.Declarations.FindAll(content, rel) {
		out[config.LineOf(content, m.Region)-1] = true
	}
	return out
}

// blockLength counts the lines of the declaration opening at `i`: up to the next line where
// a declaration starts (`starts`, by the dialect), or the end of the module; trailing blank
// lines do not count. With no declarations declared, the block runs to the end.
func blockLength(lines []string, i int, starts map[int]bool) int {
	end := len(lines)
	for j := i + 1; j < len(lines); j++ {
		if starts[j] {
			end = j
			break
		}
	}
	for end > i+1 && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return end - i
}

func leadingSpace(l string) string {
	return l[:len(l)-len(strings.TrimLeft(l, " \t"))]
}

// lineIndexer maps a byte offset to its 0-based line.
func lineIndexer(content string) func(int) int {
	var starts []int
	starts = append(starts, 0)
	for i, c := range content {
		if c == '\n' {
			starts = append(starts, i+1)
		}
	}
	return func(off int) int {
		lo, hi := 0, len(starts)-1
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if starts[mid] <= off {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		return lo
	}
}

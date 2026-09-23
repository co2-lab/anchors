package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
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
//   - each key of the double's factory that names an EXPORT of the module (by the
//     project's `export_detect`) gets its own stamp, anchored on that export's line and
//     covering its block — so a change to that member is what makes the stamp diverge;
//   - with no factory key naming an export, one stamp covers the whole module.
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
	comment := commentPrefixFor(testID)

	var written []StampWritten
	var skipped []StampSkipped
	inserts := map[int][]string{} // line index → stamp lines to insert above it
	done := map[string]bool{}

	for _, m := range detector.FindAllStringSubmatchIndex(content, -1) {
		if len(m) < 4 || m[2] < 0 {
			continue
		}
		module := strings.TrimSpace(content[m[2]:m[3]])
		if module == "" || done[module] || !isGovernedModule(module, g) {
			continue
		}
		done[module] = true
		if moduleHasStamp(module, existing) {
			continue // never rewritten — see the file comment
		}

		file, reason := resolveModuleFile(module, testID, g)
		if file == "" {
			skipped = append(skipped, StampSkipped{module, reason})
			continue
		}
		body, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			skipped = append(skipped, StampSkipped{module, "the module could not be read: " + file})
			continue
		}
		modLines := strings.Split(string(body), "\n")

		var stamps []StampWritten
		if exportRE != nil {
			for _, key := range factoryKeys(content, m[0]) {
				if s, ok := stampForExport(module, file, key, modLines, exportRE); ok {
					stamps = append(stamps, s)
				}
			}
		}
		if len(stamps) == 0 {
			s, ok := stampForWholeModule(module, file, modLines)
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

// resolveModuleFile finds the ONE code file of the map a double's specifier names.
func resolveModuleFile(spec, testID string, g *mapx.Graph) (string, string) {
	target := strings.TrimPrefix(spec, "./")
	for strings.HasPrefix(target, "../") {
		target = strings.TrimPrefix(target, "../")
	}
	if strings.HasPrefix(target, "@") {
		if i := strings.Index(target, "/"); i >= 0 {
			target = target[i+1:]
		}
	}
	var cands []string
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode {
			continue
		}
		id := withoutExtension(n.ID)
		if id == target || strings.HasSuffix(id, "/"+target) {
			cands = append(cands, n.ID)
		}
	}
	switch len(cands) {
	case 0:
		return "", "no code file of the map matches the specifier"
	case 1:
		return cands[0], ""
	}
	best, bestLen, tie := "", -1, false
	for _, c := range cands {
		l := sharedDirs(c, testID)
		switch {
		case l > bestLen:
			best, bestLen, tie = c, l, false
		case l == bestLen:
			tie = true
		}
	}
	if tie {
		return "", fmt.Sprintf("the specifier matches %d files and the test's folder does not decide: %s", len(cands), strings.Join(cands, ", "))
	}
	return best, ""
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

// factoryKeyRE finds the keys of an object literal (`useX:`, `default:`).
var factoryKeyRE = regexp.MustCompile(`([A-Za-z_$][\w$]*)\s*:`)

// factoryKeys returns the keys written inside the double's call — from the `(` after the
// match to the parenthesis that closes it. Keys of nested objects come too; they are
// filtered later by having to name an EXPORT of the real module.
func factoryKeys(content string, from int) []string {
	open := strings.Index(content[from:], "(")
	if open < 0 {
		return nil
	}
	start := from + open
	depth, end := 0, -1
	var quote byte
	for i := start; i < len(content); i++ {
		c := content[i]
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
				end = i
			}
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		return nil
	}
	seen := map[string]bool{}
	var keys []string
	for _, k := range factoryKeyRE.FindAllStringSubmatch(content[start:end], -1) {
		if !seen[k[1]] {
			seen[k[1]] = true
			keys = append(keys, k[1])
		}
	}
	return keys
}

// stampForExport anchors a stamp on the line that exports `key`, covering its block.
func stampForExport(module, file, key string, lines []string, exportRE *regexp.Regexp) (StampWritten, bool) {
	for i, l := range lines {
		m := exportRE.FindStringSubmatch(l)
		if len(m) < 2 || m[1] != key {
			continue
		}
		if !usableAnchor(lines, i) {
			return StampWritten{}, false
		}
		n := blockLength(lines, i)
		return StampWritten{Module: module, File: file, Anchor: lines[i], Count: n,
			Hash: snippetHash(strings.Join(lines[i:i+n], "\n"))}, true
	}
	return StampWritten{}, false
}

// stampForWholeModule anchors a stamp on the first usable line and covers to the end.
func stampForWholeModule(module, file string, lines []string) (StampWritten, bool) {
	for i := range lines {
		if !usableAnchor(lines, i) {
			continue
		}
		n := len(lines) - i
		return StampWritten{Module: module, File: file, Anchor: lines[i], Count: n,
			Hash: snippetHash(strings.Join(lines[i:], "\n"))}, true
	}
	return StampWritten{}, false
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

// blockLength counts the lines of the top-level declaration opening at `i`: up to the
// next line that starts a NEW top-level statement (column 0, not a closing bracket).
// Closing lines and indented lines belong to the block; trailing blank lines do not.
func blockLength(lines []string, i int) int {
	end := len(lines)
	for j := i + 1; j < len(lines); j++ {
		l := lines[j]
		if l == "" || l[0] == ' ' || l[0] == '\t' {
			continue
		}
		if strings.ContainsRune("})]", rune(l[0])) {
			continue
		}
		end = j
		break
	}
	for end > i+1 && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return end - i
}

// commentPrefixFor picks the line-comment marker by the test file's extension.
func commentPrefixFor(testID string) string {
	switch strings.ToLower(filepath.Ext(testID)) {
	case ".py", ".rb", ".sh", ".yaml", ".yml", ".toml", ".r":
		return "#"
	}
	return "//"
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

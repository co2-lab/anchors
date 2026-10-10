// @anchors
//   code: EVRSC
//   ref: RPSCR

package scan

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
)

// A proof goes stale when what it proves changes, not when the file does
// (DESIGN-evidence-by-what-it-proves.md). The file's revision (`Rev`) says whether the map
// matches the tree; the revisions here say whether what a test proved still stands. Each
// leaves out of the file what changes no behaviour, and the map carries a file's evidence
// across an edit that left them as they were.

// Evidence is what a file's proofs are measured against.
type Evidence struct {
	// Rev is the hash of the content without its `@anchors` header — the date the hook
	// writes —; in a spec, also without its navigation and change-history sections, with its
	// spacing normalized; in any other file, without the chain's flags (`@dep`, `@used-by`,
	// `@navigates`, `@no-dep`, `@no-nav`) and the lines that are only a comment — the lines
	// that only carry a flag go with them.
	Rev string
	// LineRev is, outside specs, the hash of the content with every line kept in its place —
	// the header and the flag-only lines blank: coverage and mutation name lines by number,
	// and a line inserted moves them even when it says nothing.
	LineRev string
	// Rules are, in a spec, the hash of each rule's definition — its heading and what is under
	// it, its table rows, its bold bullet —, by code; Rest is the hash of everything else the
	// evidence reads. A spec whose Rest held and some rules changed staled those rules' proofs
	// alone.
	Rules map[string]string
	Rest  string
}

var (
	// ChainFlagOnlyRE is a line that is only a comment carrying a flag of the chains;
	// chainFlagEndRE, such a comment at the end of a line of code; chainFlagBlockRE, one in a
	// block comment, wherever it sits in the line — what follows it is code.
	ChainFlagOnlyRE  = regexp.MustCompile(`^\s*(?://|#|--|/\*|<!--)\s*@(?:dep(?:\[[^\]]*\])?|no-dep|used-by|navigates|no-nav):.*$`)
	chainFlagEndRE   = regexp.MustCompile(`\s*(?://|#|--|<!--)\s*@(?:dep(?:\[[^\]]*\])?|no-dep|used-by|navigates|no-nav):.*$`)
	chainFlagBlockRE = regexp.MustCompile(`\s*/\*\s*@(?:dep(?:\[[^\]]*\])?|no-dep|used-by|navigates|no-nav):[^*]*\*/`)
	evHeadingRE      = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)
	evNoteRE         = regexp.MustCompile(`\s*\([^)]*\)\s*$`)
)

// StripChainFlag is a line without the chain's flags: a block comment carrying one, and a
// line comment carrying one at its end.
func StripChainFlag(l string) string {
	l = chainFlagBlockRE.ReplaceAllString(strings.TrimRight(l, "\r"), "")
	return chainFlagEndRE.ReplaceAllString(l, "")
}

// EvidenceOf reads the evidence revisions of one file, at rel — its extension says how its
// comments are written.
func EvidenceOf(kind, rel string, content []byte, cfg *config.Config) Evidence {
	lines := strings.Split(string(content), "\n")
	header := headerSpan(lines)
	if kind == "spec" {
		return specEvidence(lines, header, cfg)
	}
	comment := commentLines(lines, rel)
	var kept, placed []string
	for i, l := range lines {
		switch {
		case header[i] || comment[i] || ChainFlagOnlyRE.MatchString(l):
			placed = append(placed, "")
		default:
			l = StripChainFlag(l)
			kept = append(kept, l)
			placed = append(placed, l)
		}
	}
	return Evidence{
		Rev:     shortHash([]byte(strings.Join(kept, "\n"))),
		LineRev: shortHash([]byte(strings.Join(placed, "\n"))),
	}
}

// commentLines marks the lines that are nothing but a comment, in the markers of the file's
// language: a line comment alone on its line, and a block comment (`/* … */`, `{/* … */}`)
// whose lines it fills, in the languages that write them. Nothing a comment says runs — a
// rule cited above the code that answers it staled every proof of the file (reported from
// MIF). A comment at the end of a line of code stays: a `//` inside a string reads the same.
func commentLines(lines []string, rel string) map[int]bool {
	ext := ""
	if i := strings.LastIndex(rel, "."); i >= 0 {
		ext = strings.ToLower(rel[i:])
	}
	markers := config.MarkersFor(ext)
	if len(markers) == 0 || markers[0] == "<!--" {
		return nil
	}
	blocks := markers[0] == "//"
	out := map[int]bool{}
	inBlock := false
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if inBlock {
			if end := strings.Index(t, "*/"); end >= 0 {
				inBlock = false
				if rest := strings.TrimSpace(t[end+2:]); rest != "" && rest != "}" {
					continue // code after the comment's end: the line runs
				}
			}
			out[i] = true
			continue
		}
		if blocks && (strings.HasPrefix(t, "/*") || strings.HasPrefix(t, "{/*")) {
			if end := strings.Index(t, "*/"); end >= 0 {
				if rest := strings.TrimSpace(t[end+2:]); rest == "" || rest == "}" {
					out[i] = true
				}
				continue
			}
			inBlock = true
			out[i] = true
			continue
		}
		for _, m := range markers {
			if strings.HasPrefix(t, m) {
				out[i] = true
				break
			}
		}
	}
	return out
}

// headerSpan marks the lines of the `@anchors` header: the comment that opens with it near
// the top of the file, as far as the comment goes — to the `-->` of an HTML comment, or
// while the lines keep the line comment's mark.
func headerSpan(lines []string) map[int]bool {
	out := map[int]bool{}
	for i := 0; i < len(lines) && i < 10; i++ {
		if !strings.Contains(lines[i], "@anchors") {
			continue
		}
		t := strings.TrimSpace(lines[i])
		switch {
		case strings.HasPrefix(t, "<!--"):
			for j := i; j < len(lines); j++ {
				out[j] = true
				if strings.Contains(lines[j], "-->") {
					break
				}
			}
		default:
			mark := ""
			for _, m := range []string{"//", "#", "--"} {
				if strings.HasPrefix(t, m) {
					mark = m
					break
				}
			}
			if mark == "" {
				return out
			}
			for j := i; j < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[j]), mark); j++ {
				out[j] = true
			}
		}
		return out
	}
	return out
}

// proveNothing are the titles of a spec's sections with no side effect on the tests, as the
// project declares them; by default its navigation —
// each Out row has a revision of its own, which the flows asserting it read — and its
// change history, in every language and as the project names them. An Out table under
// another title is skipped by its heading.
func proveNothing(cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	// The project's own list replaces the defaults (`evidence.no_side_effect.sections`).
	if own := cfg.NoSideEffectSections(); own != nil {
		for _, t := range own {
			out[strings.ToLower(strings.TrimSpace(evNoteRE.ReplaceAllString(t, "")))] = true
		}
		return out
	}
	for _, k := range []string{"navigation", "history"} {
		for _, t := range i18n.AllTranslations("section.title." + k) {
			out[strings.ToLower(t)] = true
		}
		for _, t := range cfg.SectionTitlesFor(k) {
			out[strings.ToLower(t)] = true
		}
	}
	return out
}

func specEvidence(lines []string, header map[int]bool, cfg *config.Config) Evidence {
	skip := proveNothing(cfg)
	ruleRE := localRuleRE()
	type block struct {
		depth int
		code  string
		skip  bool
	}
	var stack []block
	var kept, rest []string
	byRule := map[string][]string{}
	blank := false
	for i, l := range lines {
		if header[i] {
			continue
		}
		l = strings.TrimRight(l, " \t\r")
		if m := evHeadingRE.FindStringSubmatch(l); m != nil {
			depth := len(m[1])
			for len(stack) > 0 && stack[len(stack)-1].depth >= depth {
				stack = stack[:len(stack)-1]
			}
			title := strings.ToLower(strings.TrimSpace(evNoteRE.ReplaceAllString(m[2], "")))
			// An Out table under any title is skipped too: its rows are read by OutRows.
			b := block{depth: depth, skip: skip[title] || skip[strings.ToLower(m[2])] || outHeadingRE.MatchString(l)}
			if c := ruleRE.FindStringSubmatch(l); c != nil {
				b.code = c[1]
			}
			stack = append(stack, b)
		}
		code, skipped := "", false
		for _, b := range stack {
			skipped = skipped || b.skip
			if b.code != "" {
				code = b.code
			}
		}
		if skipped {
			continue
		}
		// Spacing proves nothing: no run of blank lines counts more than one.
		if l == "" {
			if blank {
				continue
			}
			blank = true
		} else {
			blank = false
		}
		kept = append(kept, l)
		if !evHeadingRE.MatchString(l) {
			if c := ruleRE.FindStringSubmatch(l); c != nil {
				code = c[1]
			}
		}
		if code != "" && l != "" {
			byRule[code] = append(byRule[code], l)
		} else {
			rest = append(rest, l)
		}
	}
	ev := Evidence{
		Rev:  shortHash([]byte(strings.Join(kept, "\n"))),
		Rest: shortHash([]byte(strings.Join(rest, "\n"))),
	}
	if len(byRule) > 0 {
		codes := make([]string, 0, len(byRule))
		for c := range byRule {
			codes = append(codes, c)
		}
		sort.Strings(codes)
		ev.Rules = make(map[string]string, len(codes))
		for _, c := range codes {
			ev.Rules[c] = shortHash([]byte(strings.Join(byRule[c], "\n")))
		}
	}
	return ev
}

// captureCodeRE is a visual-regression or contract code, with its state when it has one.
var captureCodeRE = regexp.MustCompile(`\b[A-Z0-9]{3,}-(?:VR|CT)(?:-S\d{2})?\b`)

// captureCodes are the capture codes a file names outside its comment lines, sorted.
func captureCodes(rel string, content []byte) []string {
	lines := strings.Split(string(content), "\n")
	comment := commentLines(lines, rel)
	seen := map[string]bool{}
	var out []string
	for i, l := range lines {
		if comment[i] {
			continue
		}
		for _, c := range captureCodeRE.FindAllString(l, -1) {
			if !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	sort.Strings(out)
	return out
}

// CaptureCodesIn exposes the capture codes a file names outside its comment lines.
func CaptureCodesIn(rel string, content []byte) []string { return captureCodes(rel, content) }

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
	// chainFlagRE is a flag of the chains — Anchors' syntax, inside a comment of the file's.
	chainFlagRE = regexp.MustCompile(`@(?:dep(?:\[[^\]]*\])?|no-dep|used-by|navigates|no-nav):`)
	evHeadingRE = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)
	evNoteRE    = regexp.MustCompile(`\s*\([^)]*\)\s*$`)
)

// ChainFlagOnly says whether a line is only a comment, in the file's syntax, carrying a flag
// of the chains.
func ChainFlagOnly(l string, cs *config.Comments) bool {
	return cs.CommentOnly([]string{l})[0] && chainFlagRE.MatchString(l)
}

// StripChainFlag is a line without the comments that carry a flag of the chains — each with
// the blanks before it; the code around them, and a comment that carries none, stay.
func StripChainFlag(l string, cs *config.Comments) string {
	l = strings.TrimRight(l, "\r")
	if !chainFlagRE.MatchString(l) {
		return l
	}
	var b strings.Builder
	at := 0
	for _, sp := range cs.Spans(l) {
		if !chainFlagRE.MatchString(l[sp[0]:sp[1]]) {
			continue
		}
		b.WriteString(strings.TrimRight(l[at:sp[0]], " \t"))
		at = sp[1]
	}
	b.WriteString(l[at:])
	return b.String()
}

// EvidenceOf reads the evidence revisions of one file, at rel — the comments its syntax
// declares (`dialect.comments`) are no evidence: a line that is only a comment, and a
// comment at the end of a line of code, leave it. Nothing a comment says runs — a rule
// cited above the code that answers it staled every proof of the file, and a `// @no-handle`
// at the end of a call cost MIF 56 flows to run again. A file whose comments nobody declared
// keeps them all.
func EvidenceOf(kind, rel string, content []byte, cfg *config.Config) Evidence {
	cs := cfg.CommentsFor(rel)
	lines := strings.Split(string(content), "\n")
	header := headerSpan(lines, cs)
	if kind == "spec" {
		return specEvidence(lines, header, cfg)
	}
	blanked := strings.Split(cs.Blank(string(content)), "\n")
	var kept, placed []string
	for i, l := range lines {
		code := strings.TrimRight(blanked[i], " \t\r")
		switch {
		case header[i] || (strings.TrimSpace(code) == "" && strings.TrimSpace(l) != ""):
			placed = append(placed, "")
		default:
			kept = append(kept, code)
			placed = append(placed, code)
		}
	}
	return Evidence{
		Rev:     shortHash([]byte(strings.Join(kept, "\n"))),
		LineRev: shortHash([]byte(strings.Join(placed, "\n"))),
	}
}

// headerSpan marks the lines of the `@anchors` header (AnchorsHeader): the date the hook
// writes proves nothing.
func headerSpan(lines []string, cs *config.Comments) map[int]bool {
	out := map[int]bool{}
	from, to, ok := HeaderRange(lines, cs)
	for i := from; ok && i < to; i++ {
		out[i] = true
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
func captureCodes(rel string, content []byte, cs *config.Comments) []string {
	lines := strings.Split(string(content), "\n")
	comment := cs.CommentOnly(lines)
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
func CaptureCodesIn(rel string, content []byte, cs *config.Comments) []string {
	return captureCodes(rel, content, cs)
}

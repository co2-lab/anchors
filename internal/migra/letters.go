package migra

import (
	"regexp"
	"sort"
)

// codeRefRE matches a code with a two-digit item: `PLANA-F01`, `FLOWA-P04`. Exactly two
// digits and no letter or digit after them: a plan's revision block (`PLTFR-R0004`) has
// four, and is not an item whose letter changed.
var codeRefRE = regexp.MustCompile(`\b([A-Z0-9]{2,8})-([A-Z])([0-9]{2})\b`)

// LetterRenamesBetween gathers the letter renames of the steps that take a project from
// `from` to `to`, in order.
func LetterRenamesBetween(from, to int) []LetterRename {
	var out []LetterRename
	for _, s := range steps {
		if s.To > from && s.To <= to {
			out = append(out, s.RenameLetters...)
		}
	}
	return out
}

// RewriteCodeLetters rewrites, in one file's text, the codes whose letter changed: a code
// is rewritten when its unit is of a renamed kind (`kindOf` maps a unit code to its kind)
// and its letter is that kind's old one. Every other code is left as it is — the same
// letter means something else in a spec. It returns the text and, per rewritten code, how
// many times it was replaced (`PLANA-F01 → PLANA-W01`).
func RewriteCodeLetters(content string, kindOf map[string]string, renames []LetterRename) (string, map[string]int) {
	byKind := map[string]map[string]string{}
	for _, r := range renames {
		if byKind[r.Kind] == nil {
			byKind[r.Kind] = map[string]string{}
		}
		byKind[r.Kind][r.From] = r.To
	}
	counts := map[string]int{}
	out := codeRefRE.ReplaceAllStringFunc(content, func(m string) string {
		p := codeRefRE.FindStringSubmatch(m)
		to, ok := byKind[kindOf[p[1]]][p[2]]
		if !ok {
			return m
		}
		n := p[1] + "-" + to + p[3]
		counts[m+" → "+n]++
		return n
	})
	return out, counts
}

// SortedCounts lists the counts of a rewrite in a stable order, for the command's report.
func SortedCounts(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

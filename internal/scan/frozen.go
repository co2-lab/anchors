// @anchors
//   code: FRMRF
//   ref: RPSCR

package scan

import (
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// --- frozen parts: what the text marks as kept as it is (DESIGN-frozen-units.md) ---
//
// A freeze is written in two places: its RECORD, an entry of a doctrine document coded with
// the freeze letter (`QLDTQ-F01 — the plan mode ships next quarter`), and its MARKS, on what
// it covers — `frozen: QLDTQ-F01` in a header, `@frozen: QLDTQ-F01` on a rule's or a
// heading's line. A project with no doctrine writes the reason in the mark itself, and the
// mark is then the whole freeze.
//
// The scan reads both as they are written. Which parts a mark covers beyond its own file —
// a unit's derived files, the spec rules that realize a frozen doctrine rule — is the map's
// to say: it is the one that knows who is whose.

// headerFrozenRE is the header's `frozen:` line, with its value: a freeze code or a reason.
var headerFrozenRE = regexp.MustCompile(`(?m)` + config.HeaderLinePrefix + `frozen:[^\S\n]*(\S[^\n]*)$`)

// frozenMarkRE is a `@frozen:` mark with its value. A mark inside backticks is a citation —
// a document explaining the mark —, not a mark: the same guard `@TBD` has.
var frozenMarkRE = regexp.MustCompile("(?:^|[^`])@frozen:[^\\S\\n]*([^`\\s][^\\n]*)")

// headingRE is a markdown heading, capturing its hashes.
var headingRE = regexp.MustCompile(`^(#{1,6})\s`)

// markValue is what a mark says, without the punctuation of the form around it: the cell's
// pipe in a table row, the comment's closer in a header.
func markValue(v string) string {
	for _, end := range []string{"|", "-->", "*/"} {
		if i := strings.Index(v, end); i >= 0 {
			v = v[:i]
		}
	}
	return strings.TrimSpace(v)
}

// extractFrozen is the value of the header's `frozen:` line, or empty.
func extractFrozen(content string) string {
	if m := headerFrozenRE.FindStringSubmatch(string(AnchorsHeader([]byte(content)))); m != nil {
		return markValue(m[1])
	}
	return ""
}

// FrozenOf is the freeze a file's header declares (`frozen: <code or reason>`), or empty.
func FrozenOf(content []byte) string { return extractFrozen(string(content)) }

// extractFrozenRules reads the rules a spec or a doctrine document marks as frozen, each with
// its mark's value. A mark on a rule's line — or on a line of the rule's own block, as
// `@realizes` is written — freezes that rule; a mark on a heading that defines no rule
// freezes every rule under it, down to the next heading of its level or above.
func extractFrozenRules(kind, content string) map[string]string {
	if kind != "spec" && kind != "product" {
		return nil
	}
	out := map[string]string{}
	type section struct {
		level int
		value string
	}
	var open []section
	current := ""
	ruleRE := localRuleRE()
	for _, line := range strings.Split(content, "\n") {
		value := ""
		if m := frozenMarkRE.FindStringSubmatch(line); m != nil {
			value = markValue(m[1])
		}
		level := 0
		if h := headingRE.FindStringSubmatch(line); h != nil {
			level = len(h[1])
			for len(open) > 0 && open[len(open)-1].level >= level {
				open = open[:len(open)-1]
			}
		}
		if m := ruleRE.FindStringSubmatch(line); m != nil {
			current = m[1]
			switch {
			case IsFreezeCode(current):
				// A freeze's own entry is the record, not something it covers.
			case value != "":
				out[current] = value
			case len(open) > 0:
				out[current] = open[len(open)-1].value
			}
			continue
		}
		if strings.TrimSpace(line) == "" {
			current = ""
		}
		switch {
		case value == "":
		case level > 0:
			open = append(open, section{level, value})
		case current != "":
			out[current] = value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// freezeCodeRE is a freeze's code: a rule code of the freeze letter.
func freezeCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`^[A-Z0-9]` + config.CodeLengthPattern() + `-` + config.FreezeLetter + `[0-9]{2}$`)
}

// IsFreezeCode says whether a mark's value is a freeze's code rather than a reason.
func IsFreezeCode(v string) bool { return freezeCodeRE().MatchString(v) }

// extractFreezes reads a doctrine document's freeze entries: each rule of the freeze letter,
// in any of the three forms a rule is written in, with the text after its code as the
// reason.
func extractFreezes(kind, content string) map[string]string {
	if kind != "product" {
		return nil
	}
	out := map[string]string{}
	ruleRE := localRuleRE()
	for _, line := range strings.Split(content, "\n") {
		m := ruleRE.FindStringSubmatch(line)
		if m == nil || !IsFreezeCode(m[1]) {
			continue
		}
		rest := line[strings.Index(line, m[1])+len(m[1]):]
		out[m[1]] = strings.TrimSpace(strings.Trim(strings.ReplaceAll(rest, "**", ""), "`|— -"))
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// Package flagx reads the CONDITION of a feature-flag scenario.
//
// The grammar is FIXED, and that is the whole point of the decision behind it: prose in
// the "when the value" column covers any case and can be confronted with none. A fixed
// grammar refuses the case it did not foresee — and refusing loudly is what lets a gate
// say "this scenario is not the one the code implements".
//
// The operator set is not invented here. It is the union of what the market's flag tools
// actually offer, taken from their own documentation:
//
//	LaunchDarkly  in, endsWith, startsWith, matches, contains, lessThan, lessThanOrEqual,
//	              greaterThan, greaterThanOrEqual, before, after, segmentMatch,
//	              semVerEqual, semVerLessThan, semVerGreaterThan
//	Unleash       IN, NOT_IN, STR_CONTAINS, STR_STARTS_WITH, STR_ENDS_WITH, NUM_EQ,
//	              NUM_GT, NUM_GTE, NUM_LT, NUM_LTE, DATE_AFTER, DATE_BEFORE, SEMVER_EQ,
//	              SEMVER_GT, SEMVER_GTE, SEMVER_LT, SEMVER_LTE, REGEX
//	Flagsmith     Equal, Not Equal, Greater Than, Less Than, (Inclusive variants),
//	              Contains, Not Contains, Regex, Percentage Split, Is Set, Is Not Set,
//	              Modulo
//
// Three findings from that comparison shaped what follows.
//
// The first: every tool has equality, ordering, substring and regex, and they disagree
// only on SPELLING. So the grammar keeps one operator per meaning and accepts the
// spellings as aliases — a project writing `>=` and a project writing `NUM_GTE` are
// making the same statement, and the gate should not care which house style they learned.
//
// The second: only Flagsmith has `Is Set` / `Is Not Set`, and it is the operator this
// framework needs MOST. The absent flag — flag service down, new environment, local test —
// is the case that breaks in production and the case nobody writes down. The minority
// operator is the load-bearing one here, so it is first-class (`OpAbsent`, `OpPresent`),
// not an afterthought.
//
// The third: `segmentMatch` (LaunchDarkly) and `Modulo` (Flagsmith) are deliberately NOT
// here. They evaluate against a service this framework has no access to and must not
// pretend to have — a segment lives in LaunchDarkly's database, not in the repository. A
// scenario that depends on one is written as the VALUE it produces, which is what the
// code branches on anyway.
package flagx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
)

// Op is the meaning of a comparison, independent of how a given tool spells it.
type Op string

const (
	OpEq         Op = "eq"       // = "on"
	OpNe         Op = "ne"       // != "on"
	OpGt         Op = "gt"       // > 50
	OpGte        Op = "gte"      // >= 50
	OpLt         Op = "lt"       // < 50
	OpLte        Op = "lte"      // <= 50
	OpContains   Op = "contains" // contains "beta"
	OpNotContain Op = "not-contains"
	OpStartsWith Op = "starts-with"
	OpEndsWith   Op = "ends-with"
	OpMatches    Op = "matches" // regex
	OpIn         Op = "in"      // in ["a", "b"]
	OpNotIn      Op = "not-in"
	OpBefore     Op = "before" // dates
	OpAfter      Op = "after"
	OpRollout    Op = "rollout" // percentage rollout
	OpAbsent     Op = "absent"  // the flag does not answer
	OpPresent    Op = "present" // it answers, any value
)

// Condition is a parsed scenario condition: what is compared, and against what.
//
// Operand is empty exactly for OpAbsent and OpPresent, which compare against nothing.
type Condition struct {
	Op      Op
	Operand string
	Raw     string
}

// NeedsOperand reports whether this operator compares against a value at all.
func (o Op) NeedsOperand() bool { return o != OpAbsent && o != OpPresent }

// aliases maps every accepted spelling to its meaning. Lowercased, spaces and underscores
// normalised, so `STR_STARTS_WITH`, `starts with` and `startsWith` all arrive the same.
//
// Symbols come first because they are what a person writes by hand; the named forms are
// the vocabulary of whichever tool the project already uses.
var aliases = map[string]Op{
	"=": OpEq, "==": OpEq, "eq": OpEq, "equal": OpEq, "equals": OpEq,
	"numeq": OpEq, "semvereq": OpEq, "is": OpEq,

	"!=": OpNe, "<>": OpNe, "ne": OpNe, "notequal": OpNe, "isnot": OpNe,

	">": OpGt, "gt": OpGt, "numgt": OpGt, "greaterthan": OpGt, "semvergt": OpGt,
	">=": OpGte, "gte": OpGte, "numgte": OpGte, "greaterthanorequal": OpGte,
	"greaterthaninclusive": OpGte, "semvergte": OpGte,
	"<": OpLt, "lt": OpLt, "numlt": OpLt, "lessthan": OpLt, "semverlt": OpLt,
	"<=": OpLte, "lte": OpLte, "numlte": OpLte, "lessthanorequal": OpLte,
	"lessthaninclusive": OpLte, "semverlte": OpLte,

	"contains": OpContains, "strcontains": OpContains, "includes": OpContains,
	"notcontains": OpNotContain, "doesnotcontain": OpNotContain,

	"startswith": OpStartsWith, "strstartswith": OpStartsWith,
	"endswith": OpEndsWith, "strendswith": OpEndsWith,

	"matches": OpMatches, "regex": OpMatches, "matchesregex": OpMatches,

	"in": OpIn, "oneof": OpIn,
	"notin": OpNotIn, "noneof": OpNotIn,

	"before": OpBefore, "datebefore": OpBefore,
	"after": OpAfter, "dateafter": OpAfter,

	"rollout": OpRollout, "percentage": OpRollout, "percentagesplit": OpRollout,
	"percentagerollout": OpRollout,
}

// normalise folds a written operator to its lookup key: case, spaces, underscores and
// hyphens all disappear, so the three house styles collapse onto one entry.
func normalise(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.NewReplacer(" ", "", "_", "", "-", "").Replace(s)
}

// valuelessOps are the spellings of the two operators that compare against nothing.
//
// The English forms are always accepted (they are the engine's own vocabulary, and a
// project may write in one language and read in another), and the catalog adds whatever
// each supported language calls them.
func valuelessOps() map[string]Op {
	out := map[string]Op{
		"absent": OpAbsent, "isnotset": OpAbsent, "notset": OpAbsent,
		"unset": OpAbsent, "missing": OpAbsent, "undefined": OpAbsent, "null": OpAbsent,
		"present": OpPresent, "isset": OpPresent, "set": OpPresent,
		"anyvalue": OpPresent, "any": OpPresent,
	}
	for _, w := range i18n.AllTranslations("flag.keyword.absent") {
		out[normalise(w)] = OpAbsent
	}
	for _, w := range i18n.AllTranslations("flag.keyword.present") {
		out[normalise(w)] = OpPresent
	}
	return out
}

// symbolRE pulls a leading symbolic operator off the condition. Two-character symbols are
// listed first: `>` would otherwise swallow the `>=` before it could match.
var symbolRE = regexp.MustCompile(`^\s*(>=|<=|!=|<>|==|=|>|<)\s*(.*)$`)

// Parse reads one condition cell.
//
// The error is deliberately specific about what was not understood. A grammar that
// refuses without saying what it wanted just moves the guessing to the person writing the
// scenario.
func Parse(raw string) (Condition, error) {
	s := strings.TrimSpace(raw)
	s = strings.Trim(s, "`")
	s = strings.TrimSpace(s)
	if s == "" {
		return Condition{}, fmt.Errorf("empty condition — a scenario states when it applies")
	}

	// The valueless operators, checked before anything else: their spelling is a whole
	// phrase, and feeding them through the operand split would strip half of it away.
	//
	// Measured, and the reason this lookup exists at all: with the words nailed in
	// English, a flag file written in Portuguese saying `ausente` did not parse as the
	// absent case — it fell through to the word split and came out as the operator `ne`
	// carrying the operand `set`. A WRONG answer, silently, on precisely the scenario
	// that matters most. The keywords come from the TRANSLATION CATALOG for the same
	// reason `flow.keyword.fits` does.
	if op, ok := valuelessOps()[normalise(s)]; ok {
		return Condition{Op: op, Raw: raw}, nil
	}

	if m := symbolRE.FindStringSubmatch(s); m != nil {
		op := aliases[m[1]]
		operand := strings.TrimSpace(m[2])
		if operand == "" {
			return Condition{}, fmt.Errorf("operator %q with nothing to compare against", m[1])
		}
		return Condition{Op: op, Operand: unquote(operand), Raw: raw}, nil
	}

	// Word operators: try the longest word prefix first, so `starts with` is not read as
	// the unknown operator `starts` carrying the operand `with "beta"`.
	fields := strings.Fields(s)
	for n := len(fields); n > 0; n-- {
		key := normalise(strings.Join(fields[:n], ""))
		op, ok := aliases[key]
		if !ok {
			continue
		}
		operand := strings.TrimSpace(strings.Join(fields[n:], " "))
		if operand == "" {
			return Condition{}, fmt.Errorf("operator %q with nothing to compare against",
				strings.Join(fields[:n], " "))
		}
		return Condition{Op: op, Operand: unquote(operand), Raw: raw}, nil
	}

	return Condition{}, fmt.Errorf(
		"unrecognised condition %q — write a comparison (`= \"on\"`, `>= 50`, "+
			"`contains \"beta\"`, `matches \"^v2\"`) or the absent case (`absent`)", s)
}

// unquote strips one layer of surrounding quotes, leaving inner ones alone.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

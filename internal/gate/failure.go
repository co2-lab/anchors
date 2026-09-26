package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- the FAILURE declared, handled and recorded ---
//
// The spec already catalogued how the unit fails, with its own letter (`-E`) and a table
// of condition and effect. And nothing confronted it: the rules crossed the whole pipeline
// and nobody asked whether they were handled, whether they logged, or whether they happen.
//
// These three gates close the first of the concept's three layers — the one that is PURE
// STATIC confrontation, and therefore depends on neither production nor any log format:
//
//	failure-handled    does the declared failure have a path that handles it?
//	failure-logged     does that path RECORD the occurrence?
//	failure-declared   does the handling that exists answer some declared failure?
//
// The third is the inverse of the first, and catches the commonest case: somebody wrote a
// defence and never declared what it prevents.
//
// HANDLING IS NOT "HAVING A CATCH". Nailing `catch` would nail the syntax of one family of
// languages — `if (x == null) { return refuse() }` handles just as much, and so does a
// `match` in Rust. What the shapes have in common is not the look, it is the EFFECT: the
// failure becomes part of the flow and the application carries on. Who knows the local
// dialect's shape is the project, in `handle_patterns` — the same mechanism as
// `guard_patterns`, which already existed.

// failureRuleRE finds a FAILURE rule (`-E`) in the three catalogued forms.
var failureRuleRE = regexp.MustCompile(`(?m)(?:^#{1,6}\s+|^\s*\|\s*` + "`?" + `|^\s*-\s+\*\*)([A-Z0-9]{3,6}-E[0-9]{2})`)

// resilientRE — the failure UNDERSTOOD and absorbed by the flow.
//
// It is an assertion different from the two that already existed, and so it gets its own
// marker:
//
//	@no-<thing>: <reason>   "it will never have one"   — permanent waiver
//	@TBD: <reason>          "it does not have one yet"  — debt, keeps showing up
//	@resilient: <reason>    "it happens, I know why, and it is handled"
//
// It is not a waiver (the failure is real and happens) nor debt (nothing is pending): it is
// knowledge acquired. The mandatory reason is what tells it apart from silencing an alert —
// it answers the question somebody will ask in six months, "why do we ignore this?".
var resilientRE = regexp.MustCompile("(?i)(^|[^`])@resilient[^\\S\\n]*:[^\\S\\n]*\\S+")

// observingRE — the failure that OCCURS and whose cause is not yet known.
//
// It is the third conclusion of the observation layer, and the one that no observability
// tool records: the difference between "nobody investigated" and "we investigated and
// still do not know". The second is KNOWLEDGE, and today it is lost — the next person to
// look starts from zero, ruling out what somebody already ruled out.
//
// The mandatory text is what carries that: `@observing: ruled out timeout and partner
// retry; happens only on accounts in migration`.
var observingRE = regexp.MustCompile("(?i)(^|[^`])@observing[^\\S\\n]*:[^\\S\\n]*\\S+")

// FailureConclusion is what a spec says about a failure it declared, after seeing it
// happen.
//
// The three are progress, and they differ in what they ask of whoever reads next:
//
//	cause found    the rule carries it in prose — the handling can stop being generic
//	resilient      `@resilient: <reason>` — it leaves the radar without leaving the record
//	observing      `@observing: <what was ruled out>` — the question stays open, with history
type FailureConclusion struct {
	Rule      string
	Resilient string
	Observing string
}

// FailureConclusions reads, per rule, what the spec concluded about each failure.
func FailureConclusions(content string) map[string]FailureConclusion {
	out := map[string]FailureConclusion{}
	for _, line := range strings.Split(content, "\n") {
		for _, m := range failureRuleRE.FindAllStringSubmatch(line, -1) {
			c := out[m[1]]
			c.Rule = m[1]
			if mm := reasonOf(resilientRE, line); mm != "" {
				c.Resilient = mm
			}
			if mm := reasonOf(observingRE, line); mm != "" {
				c.Observing = mm
			}
			out[m[1]] = c
		}
	}
	return out
}

// reasonOf extracts the written reason of a marker — everything after the colon, to the
// end of the line or the cell.
//
// The reason is what separates a conclusion from silencing: a bare marker would say "stop
// asking" without saying why, and in six months nobody would know whether it still holds.
//
// It is read from the LINE and not from the marker's own match, because the marker pattern
// stops at the first token — it only has to prove the reason EXISTS. Reading the reason
// from it truncated `@resilient: the partner restarts at 3am` down to `the`, which throws
// away exactly the part that is worth keeping.
func reasonOf(re *regexp.Regexp, line string) string {
	loc := re.FindStringIndex(line)
	if loc == nil {
		return ""
	}
	_, after, ok := strings.Cut(line[loc[0]:], ":")
	if !ok {
		return ""
	}
	// A table cell ends at the pipe: the reason is what the author wrote in THIS column,
	// and swallowing the next one would attribute to the conclusion a text that belongs to
	// another field.
	if i := strings.Index(after, "|"); i >= 0 {
		after = after[:i]
	}
	return strings.TrimSpace(after)
}

// declaredFailures reads the `-E` rules from the spec, separating the ones marked resilient.
func declaredFailures(content string) (all, resilient []string) {
	seen := map[string]bool{}
	for _, line := range strings.Split(content, "\n") {
		for _, m := range failureRuleRE.FindAllStringSubmatch(line, -1) {
			if seen[m[1]] {
				continue
			}
			seen[m[1]] = true
			all = append(all, m[1])
			if resilientRE.MatchString(line) {
				resilient = append(resilient, m[1])
			}
		}
	}
	return all, resilient
}

// governedCode reads the non-comment content of the code this spec specifies.
func governedCode(n mapx.Node, root string, g *mapx.Graph) (string, bool) {
	if g == nil {
		return "", false
	}
	var body strings.Builder
	found := false
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.To))
		if err != nil {
			continue
		}
		found = true
		body.WriteString(stripLineComments(string(b)))
		body.WriteString("\n")
	}
	return body.String(), found
}

// anyMatch says whether any of the declared patterns matches the code.
func anyMatch(d config.Dialect, patterns []string, code string) bool {
	for _, p := range patterns {
		if re := d.Compile(p); re != nil && re.MatchString(code) {
			return true
		}
	}
	return false
}

// --- does the declared failure have a path that handles it? ---
func checkFailureHandled(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.failure_handled.skip_not_spec")
	}
	all, resilient := declaredFailures(content)
	if len(all) == 0 {
		return Skip, i18n.T("gate.failure_handled.skip_none")
	}
	d := cfg.DialectFor()
	if len(d.HandlePatterns) == 0 {
		// Pending and not Pass: without knowing how to recognise a handling path, the gate
		// verified nothing — and a ✓ here would stamp what was never measured.
		return Pending, i18n.T("gate.failure_handled.skip_no_dialect")
	}
	code, ok := governedCode(n, root, g)
	if !ok {
		return Pending, i18n.T("gate.failure_handled.pending_no_code")
	}

	// A RESILIENT failure leaves the charge: it was understood, the flow absorbs it, and
	// the reason is written beside it.
	isResilient := map[string]bool{}
	for _, r := range resilient {
		isResilient[r] = true
	}

	// The ruler is over the SET, not over each rule: the code does not cite the failure's
	// code (`CRED-E01` does not appear in an `if`), so there is no way to tie one rule to a
	// specific path. What can be asserted is whether the unit has NO handling at all while
	// declaring failures — and that is exactly the case that matters.
	if anyMatch(d, d.HandlePatterns, code) {
		return Pass, ""
	}
	var charged []string
	for _, f := range all {
		if !isResilient[f] {
			charged = append(charged, f)
		}
	}
	if len(charged) == 0 {
		return Pass, ""
	}
	sort.Strings(charged)
	return Fail, fmt.Sprintf(i18n.T("gate.failure_handled.untreated"), len(charged), strings.Join(charged, ", "))
}

// --- does the handling RECORD the occurrence? ---
func checkFailureLogged(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.failure_handled.skip_not_spec")
	}
	all, resilient := declaredFailures(content)
	if len(all) == 0 {
		return Skip, i18n.T("gate.failure_handled.skip_none")
	}
	d := cfg.DialectFor()
	if len(d.HandlePatterns) == 0 || len(d.LogPatterns) == 0 {
		return Pending, i18n.T("gate.failure_handled.skip_no_dialect")
	}
	code, ok := governedCode(n, root, g)
	if !ok {
		return Pending, i18n.T("gate.failure_handled.pending_no_code")
	}
	// With no handling at all there is nothing to charge here: the sibling gate charges it,
	// and two gates accusing the same defect turn into noise.
	if !anyMatch(d, d.HandlePatterns, code) {
		return Skip, ""
	}
	if anyMatch(d, d.LogPatterns, code) {
		return Pass, ""
	}
	isResilient := map[string]bool{}
	for _, r := range resilient {
		isResilient[r] = true
	}
	var charged []string
	for _, f := range all {
		if !isResilient[f] {
			charged = append(charged, f)
		}
	}
	if len(charged) == 0 {
		return Pass, ""
	}
	sort.Strings(charged)
	return Fail, fmt.Sprintf(i18n.T("gate.failure_logged.unlogged"), len(charged), strings.Join(charged, ", "))
}

// --- does the handling that exists answer some declared failure? ---
//
// The inverse of `failure-handled`, and the commonest case: somebody wrote a defence and
// never declared what it prevents. The defence may be perfectly right — what is missing is
// the spec saying which failure it answers, so whoever reads it later knows if it still
// applies.
func checkFailureDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.failure_handled.skip_not_spec")
	}
	d := cfg.DialectFor()
	if len(d.HandlePatterns) == 0 {
		return Pending, i18n.T("gate.failure_handled.skip_no_dialect")
	}
	code, ok := governedCode(n, root, g)
	if !ok {
		return Pending, i18n.T("gate.failure_handled.pending_no_code")
	}
	hits := 0
	for _, p := range d.HandlePatterns {
		if re := d.Compile(p); re != nil {
			hits += len(re.FindAllString(code, -1))
		}
	}
	if hits == 0 {
		return Skip, i18n.T("gate.failure_declared.skip_none")
	}
	all, _ := declaredFailures(content)
	if len(all) > 0 {
		return Pass, ""
	}
	if reason := errorsClosedAsNone(content); reason != "" {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(i18n.T("gate.failure_declared.undeclared"), hits, len(all))
}

// errorsSectionRE opens the spec's failure section, in the languages the catalogue writes.
var errorsSectionRE = regexp.MustCompile(`(?im)^#{2,4}\s*(?:errors?|failures?|erros?|falhas?|errores?|fallos?)\b[^\n]*\n`)

// errorsNoneRE is the section closed as "this unit handles no failure", with the reason.
var errorsNoneRE = regexp.MustCompile(`(?i)^(?:none|nenhum[a]?|ningun[oa]?)\s*(?:—|–|-|:)\s*(\S.*)$`)

// errorsClosedAsNone returns the reason when the failure section is closed with `none`.
//
// The handling patterns are textual, and in Go `== nil {` is also a lazy map init
// (`if m[k] == nil { m[k] = … }`) and a regex that did not match (`if m == nil {
// continue }`) — normal flow, not a defence. A unit whose only matches are those has no
// failure to declare and nothing to alias, and the honest answer is to say so: `none —
// <why>`, like a decisions section closed with `none`. The reason is mandatory, so the
// claim is written where a reviewer reads it; a bare `none` does not close the section.
func errorsClosedAsNone(content string) string {
	loc := errorsSectionRE.FindStringIndex(content)
	if loc == nil {
		return ""
	}
	for _, line := range strings.Split(content[loc[1]:], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if m := errorsNoneRE.FindStringSubmatch(line); m != nil {
			return strings.TrimSpace(m[1])
		}
		return ""
	}
	return ""
}

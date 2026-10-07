// @anchors
//   code: FLGTF
//   ref: FLRAI

package gate

import (
	"fmt"
	"path"
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
//
// Compiled per CALL, as defineRuleCaptureRE: the code length comes from `code_lengths`,
// loaded after the package globals. As a `var` with a fixed `{3,6}` it never saw a code
// of another declared length (a 7- or 8-character code was invisible to the gate).
func failureRuleRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)(?:^#{1,6}\s+|^\s*\|\s*` + "`?" + `|^\s*-\s+\*\*)([A-Z0-9]` +
		config.CodeLengthPattern() + `-E[0-9]{2})`)
}

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
	ruleRE := failureRuleRE()
	for _, line := range strings.Split(content, "\n") {
		for _, m := range ruleRE.FindAllStringSubmatch(line, -1) {
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
	ruleRE := failureRuleRE()
	for _, line := range strings.Split(content, "\n") {
		if retiredLine(line) {
			continue // a retired failure is neither handled nor logged any more
		}
		for _, m := range ruleRE.FindAllStringSubmatch(line, -1) {
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
		b, err := readFile(root, e.To)
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
	// EACH FALLIBLE CALL has its handling beside it. Read over the whole file, one `catch`
	// anywhere answered a unit with five fetches, and the screen that reads `data` from a
	// query and never looks at its error passed (reported from MIF: an endless spinner, an
	// empty list, a form of defaults that erased the data on save).
	if open := unhandledCalls(n, root, g, cfg.DialectFor()); len(open) > 0 {
		shown := open
		if len(shown) > 5 {
			shown = append(append([]string(nil), shown[:5]...), fmt.Sprintf("… +%d", len(open)-5))
		}
		return Fail, i18n.T("gate.failure_handled.call_unhandled", strings.Join(shown, "; "))
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
	// A unit that CAN fail declares how. The other questions start from what somebody
	// wrote — a failure declared, a handling coded —, and a unit that consumes a fallible
	// source and wrote neither passed all of them in silence (reported from MIF: six detail
	// screens showed "not found", an endless spinner, an empty list, or a form of defaults
	// that erased the data on save, when the fetch failed).
	//
	// Each source is answered by a failure that NAMES it — its `DEPn` or its name, in the
	// failure's row or in its row of the rules' uses. "At least one failure" was not enough:
	// the screens that showed "not found" declared exactly that failure, and the load that
	// failed was no failure of theirs.
	if srcs := uncoveredSources(content, fallibleSources(n, root, g, cfg, d)); len(srcs) > 0 {
		if errorsClosedAsNone(content) == "" && !noFailureRE.MatchString(content) {
			shown := srcs
			if len(shown) > 5 {
				shown = append(append([]string(nil), shown[:5]...), fmt.Sprintf("… +%d", len(srcs)-5))
			}
			return Fail, i18n.T("gate.failure_declared.fallible_undeclared", strings.Join(shown, "; "))
		}
	}
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

// noFailureRE is the spec's waiver of the fallible-source question, with its reason.
var noFailureRE = regexp.MustCompile(`@no-failure:\s*\S`)

// FallibleCall is a call that can fail, in a file of the unit (see FalliblePatterns).
type FallibleCall struct {
	File    string
	Line    int
	Text    string
	Pattern config.FalliblePattern
}

// fallibleCalls are the calls of the unit's code (the files its spec `specifies`) that a
// declared pattern recognises as fallible. Comment lines do not count, and a trailing
// comment is cut before matching; the line is the file's own.
func fallibleCalls(n mapx.Node, root string, g *mapx.Graph, d config.Dialect) []FallibleCall {
	if g == nil || len(d.FalliblePatterns) == 0 {
		return nil
	}
	type compiled struct {
		p  config.FalliblePattern
		re *regexp.Regexp
	}
	var pats []compiled
	for _, p := range d.FalliblePatterns {
		if re := d.Compile(p.Call); re != nil {
			pats = append(pats, compiled{p, re})
		}
	}
	var out []FallibleCall
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		b, err := readFile(root, e.To)
		if err != nil {
			continue
		}
		for i, ln := range strings.Split(string(b), "\n") {
			if commentLine(ln) {
				continue
			}
			code := cutInlineComment(ln)
			for _, c := range pats {
				if c.re.MatchString(code) {
					out = append(out, FallibleCall{File: e.To, Line: i + 1, Text: strings.TrimSpace(code), Pattern: c.p})
					break
				}
			}
		}
	}
	return out
}

// commentLine says whether a line is only a comment, by the same reading stripLineComments
// makes.
func commentLine(ln string) bool {
	return strings.TrimSpace(stripLineComments(ln)) == "" && strings.TrimSpace(ln) != ""
}

// FallibleSource is what makes a unit able to fail, and the names a failure cites it by.
type FallibleSource struct {
	Label string   // how the verdict names it: `file:line` or `DEP2 path`
	Names []string // what a failure cites to answer it: the call's name, the DEPn, the file's stem
}

// fallibleSources names what makes the unit able to fail: each fallible call of its code
// (`file:line`, cited by the name called), and each dependency its spec declares on a file
// of a layer the project marks `fallible: true` (`DEP2 path`, cited by `DEP2` or the
// file's stem).
func fallibleSources(n mapx.Node, root string, g *mapx.Graph, cfg *config.Config, d config.Dialect) []FallibleSource {
	var out []FallibleSource
	for _, c := range fallibleCalls(n, root, g, d) {
		src := FallibleSource{Label: fmt.Sprintf("%s:%d", c.File, c.Line)}
		if re := d.Compile(c.Pattern.Call); re != nil {
			if m := calledNameRE.FindStringSubmatch(re.FindString(c.Text)); m != nil {
				src.Names = append(src.Names, m[1])
			}
		}
		out = append(out, src)
	}
	if g != nil && cfg != nil {
		seen := map[string]bool{}
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type != mapx.EdgeDependsOn || seen[e.To] {
				continue
			}
			t := g.Node(e.To)
			if t == nil {
				continue
			}
			if l, ok := cfg.Layers[t.Layer]; ok && l.Fallible {
				seen[e.To] = true
				src := FallibleSource{Label: e.To}
				if e.Dep != "" {
					src.Label = e.Dep + " " + e.To
					src.Names = append(src.Names, e.Dep)
				}
				base := path.Base(e.To)
				if i := strings.Index(base, "."); i > 0 {
					base = base[:i]
				}
				src.Names = append(src.Names, base)
				out = append(out, src)
			}
		}
	}
	return out
}

// calledNameRE is the name a matched call calls: the identifier right before its parenthesis.
var calledNameRE = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

// uncoveredSources are the sources no declared failure names: a source is answered when a
// line that carries a failure code of the spec — its row, or its row of the rules' uses —
// cites one of the source's names as a whole word.
func uncoveredSources(content string, srcs []FallibleSource) []string {
	if len(srcs) == 0 {
		return nil
	}
	all, _ := declaredFailures(content)
	var failureLines strings.Builder
	if len(all) > 0 {
		codes := map[string]bool{}
		for _, c := range all {
			codes[c] = true
		}
		codeRE := regexp.MustCompile(`[A-Z0-9]` + config.CodeLengthPattern() + `-E[0-9]{2}`)
		for _, ln := range strings.Split(content, "\n") {
			for _, c := range codeRE.FindAllString(ln, -1) {
				if codes[c] {
					failureLines.WriteString(ln)
					failureLines.WriteString("\n")
					break
				}
			}
		}
	}
	text := failureLines.String()
	var out []string
	for _, src := range srcs {
		named := false
		for _, name := range src.Names {
			if name != "" && regexp.MustCompile(`\b`+regexp.QuoteMeta(name)+`\b`).MatchString(text) {
				named = true
				break
			}
		}
		if !named {
			out = append(out, src.Label)
		}
	}
	return out
}

// noHandleRE waives one fallible call, with its reason, on its line or the line above.
var noHandleRE = regexp.MustCompile(`@no-handle:\s*\S`)

// handleBefore is the most lines above a fallible call its handling may stand: a
// destructuring that reads the error spans the lines before the call it receives.
const handleBefore = 10

// statementStart is the first line of the statement that ends at line i: it climbs while
// the line above continues it, and stops at a blank line or one that ends a statement
// (`;`, `}`, `)`) — the previous statement's handling is not this call's.
func statementStart(lines []string, i int) int {
	j := i
	for j > 0 && i-j < handleBefore {
		t := strings.TrimSpace(cutInlineComment(lines[j-1]))
		if t == "" || strings.HasSuffix(t, ";") || strings.HasSuffix(t, "}") || strings.HasSuffix(t, ")") {
			break
		}
		j--
	}
	return j
}

// unhandledCalls are the fallible calls of the unit with no handling in their window — the
// call's statement from its first line (a destructuring above it), and the pattern's window
// after it
// —, named `file:line \`call\“. The handling is the pattern's own `handled`, or the
// project's `handle_patterns` when the pattern declares none.
func unhandledCalls(n mapx.Node, root string, g *mapx.Graph, d config.Dialect) []string {
	calls := fallibleCalls(n, root, g, d)
	if len(calls) == 0 {
		return nil
	}
	files := map[string][]string{}
	var out []string
	for _, c := range calls {
		lines, ok := files[c.File]
		if !ok {
			b, err := readFile(root, c.File)
			if err != nil {
				continue
			}
			lines = strings.Split(string(b), "\n")
			files[c.File] = lines
		}
		i := c.Line - 1
		if noHandleRE.MatchString(lines[i]) || (i > 0 && noHandleRE.MatchString(lines[i-1])) {
			continue
		}
		var handled []*regexp.Regexp
		if c.Pattern.Handled != "" {
			if re := d.Compile(c.Pattern.Handled); re != nil {
				handled = append(handled, re)
			}
		} else {
			for _, p := range d.HandlePatterns {
				if re := d.Compile(p); re != nil {
					handled = append(handled, re)
				}
			}
		}
		if len(handled) == 0 {
			continue // nothing says what handling looks like: nothing to confront
		}
		window := c.Pattern.Window
		if window <= 0 {
			window = config.DefaultFallibleWindow
		}
		from, to := statementStart(lines, i), i+window
		if to >= len(lines) {
			to = len(lines) - 1
		}
		found := false
		for j := from; j <= to && !found; j++ {
			if commentLine(lines[j]) {
				continue
			}
			code := cutInlineComment(lines[j])
			for _, re := range handled {
				if re.MatchString(code) {
					found = true
					break
				}
			}
		}
		if !found {
			out = append(out, fmt.Sprintf("%s:%d `%s`", c.File, c.Line, c.Text))
		}
	}
	return out
}

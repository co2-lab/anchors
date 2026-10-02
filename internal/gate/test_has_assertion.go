// @anchors
//   ref: THSAS

package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testlist"
)

// test-has-assertion: a test that asserts nothing passes whatever the code does.
//
// It carries the scenario's code, its title matches, the scenario counts as proved — and
// the body calls the code and checks nothing. Every gate that reads titles is green over
// it. What an assertion looks like is the test library's, so the project declares it
// (`dialect.tests.assertion`, with a default for the families that have one), and the gate
// only asks whether the test's body holds one.
//
// Where a test's body ends is also the language's. A source that knows it says it (the
// script's `end`); otherwise the body is read by its layout, which every language keeps
// whatever its syntax: a block opened by a bracket closes on the first line back at the
// opener's indentation that starts with a closer, a block opened by anything else ends on
// the first line back at that indentation. Lines inside a multi-line literal (a backtick or
// triple-quoted string) are text and do not count; a literal delimited any other way can
// end the block early — the script's `end` is exact.
func checkTestHasAssertion(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest || n.Support {
		return Skip, i18n.T("gate.test_has_assertion.skip_not_test")
	}
	src := testsSource(cfg)
	if !src.Declared() {
		return Skip, i18n.T("gate.test_has_assertion.skip_no_source")
	}
	assertion := ""
	if t := cfg.DialectFor().Tests; t != nil {
		assertion = t.Assertion
	}
	if strings.TrimSpace(assertion) == "" {
		return Skip, i18n.T("gate.test_has_assertion.skip_no_assertion")
	}
	re, err := regexp.Compile(assertion)
	if err != nil {
		return Fail, fmt.Sprintf("`dialect.tests.assertion`: %v", err)
	}
	tests, _, err := projectTests(root, g, cfg)
	if err != nil {
		return Fail, err.Error()
	}
	mine := testsIn(tests, []string{n.ID})
	if len(mine) == 0 {
		return Skip, i18n.T("gate.test_has_assertion.skip_no_tests")
	}
	labels := gateEntry(cfg, "test-has-assertion").Labels
	lines := strings.Split(content, "\n")
	var empty []string
	for _, t := range mine {
		body := testBody(lines, t, labels)
		if t.Line >= 2 && t.Line-2 < len(lines) {
			body = lines[t.Line-2] + "\n" + body
		}
		if noAssertRE.MatchString(body) {
			continue // declared: the test asserts nothing on purpose, and says why
		}
		if !re.MatchString(body) {
			empty = append(empty, fmt.Sprintf("%d: %s", t.Line, t.Title))
		}
	}
	if len(empty) > 0 {
		return Fail, i18n.T("gate.test_has_assertion.fail", len(empty), strings.Join(empty, "; "))
	}
	return Pass, ""
}

// testBody is the text of a test: its lines as the source bounds them, or as the layout
// does. With `labels`, a test whose body is empty stands for the block around it.
func testBody(lines []string, t testlist.Test, labels bool) string {
	i := t.Line - 1
	if i < 0 || i >= len(lines) {
		return ""
	}
	first, last := i, blockEnd(lines, i)
	if t.End > 0 {
		last = min(t.End-1, len(lines)-1)
	}
	if labels && last == first && emptyBodyRE.MatchString(lines[i]) {
		if k := enclosingLine(lines, i); k >= 0 {
			first, last = k, blockEnd(lines, k)
		}
	}
	return strings.Join(lines[first:last+1], "\n")
}

// noAssertRE is the waiver of a test with no assertion, with its reason, in the test's
// body or on the line above it: `@no-assert: <why>`.
var noAssertRE = regexp.MustCompile(`@no-assert[^\S\n]*:[^\S\n]*\S+`)

var (
	// opensBlockRE: a line that leaves a block open at its end — a bracket, or a `do`.
	opensBlockRE = regexp.MustCompile(`(?:[{(\[]|\bdo)\s*$`)
	// closesBlockRE: a line that starts by closing one.
	closesBlockRE = regexp.MustCompile(`^\s*(?:[}\])]|end\b)`)
	// emptyBodyRE: a block opened and closed with nothing inside.
	emptyBodyRE = regexp.MustCompile(`\{\s*\}`)
)

func indentOf(l string) int { return len(l) - len(strings.TrimLeft(l, " \t")) }

// blockEnd is the last line of the block that line i opens, by its layout. A line that
// starts inside a multi-line literal is text, not layout, and never ends the block.
func blockEnd(lines []string, i int) int {
	base := indentOf(lines[i])
	bracket := opensBlockRE.MatchString(lines[i])
	next := i + 1
	for next < len(lines) && strings.TrimSpace(lines[next]) == "" {
		next++
	}
	if !bracket && (next >= len(lines) || indentOf(lines[next]) <= base) {
		return i // opened and closed on its own line
	}
	var lit literalState
	lit.feed(lines[i])
	for j := i + 1; j < len(lines); j++ {
		inside := lit.open()
		lit.feed(lines[j])
		if inside || strings.TrimSpace(lines[j]) == "" || indentOf(lines[j]) > base {
			continue
		}
		if !bracket {
			return j - 1
		}
		// A line that closes and opens again goes on with the block: the table of an
		// `it.each([` closes into the test's own body (`])('title', () => {`), and a
		// `} else {` into the next branch. Only a closer that opens nothing ends it.
		if closesBlockRE.MatchString(lines[j]) && !opensBlockRE.MatchString(lines[j]) {
			return j
		}
	}
	return len(lines) - 1
}

// literalState follows the delimiters most languages give a multi-line literal — the
// backtick (Go, JavaScript) and the triple quote (Python, Kotlin, Swift) — line by line.
// A backtick quoted on its own (`"`"`) is a character, not a delimiter.
type literalState struct{ backtick, triple2, triple1 bool }

var loneBacktickRE = regexp.MustCompile("[\"']`[\"']")

func (s *literalState) feed(line string) {
	line = loneBacktickRE.ReplaceAllString(line, "")
	if strings.Count(line, "`")%2 == 1 {
		s.backtick = !s.backtick
	}
	if strings.Count(line, `"""`)%2 == 1 {
		s.triple2 = !s.triple2
	}
	if strings.Count(line, "'''")%2 == 1 {
		s.triple1 = !s.triple1
	}
}

func (s literalState) open() bool { return s.backtick || s.triple2 || s.triple1 }

// enclosingLine is the nearest line above i with less indentation: the opener of the block
// that holds line i. -1 when i is at the top level.
func enclosingLine(lines []string, i int) int {
	base := indentOf(lines[i])
	for k := i - 1; k >= 0; k-- {
		if strings.TrimSpace(lines[k]) != "" && indentOf(lines[k]) < base {
			return k
		}
	}
	return -1
}

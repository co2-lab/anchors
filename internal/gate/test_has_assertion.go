// @anchors
//   code: THAGT
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
		body := testBody(lines, t, labels, newLayout(lines, n.ID, cfg))
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

// layout is how a test file's blocks are read: the dialect's `tests.block_open`,
// `block_close` and `empty_body`, and the lines that start inside a string spanning lines
// (the file's comment syntax, `multiline`) — text, never a block's end.
type layout struct {
	open, close, empty *regexp.Regexp
	inString           []bool
}

func newLayout(lines []string, rel string, cfg *config.Config) layout {
	var l layout
	if t := cfg.DialectFor().Tests; t != nil {
		l.open, l.close, l.empty = compileOrNil(t.BlockOpen), compileOrNil(t.BlockClose), compileOrNil(t.EmptyBody)
	}
	l.inString = cfg.CommentsFor(rel).InString(lines)
	return l
}

// testBody is the text of a test: its lines as the source bounds them, or as the layout
// does. With `labels`, a test whose body is empty stands for the block around it.
func testBody(lines []string, t testlist.Test, labels bool, lay layout) string {
	i := t.Line - 1
	if i < 0 || i >= len(lines) {
		return ""
	}
	first, last := i, blockEnd(lines, i, lay)
	if t.End > 0 {
		last = min(t.End-1, len(lines)-1)
	}
	if labels && last == first && lay.empty != nil && lay.empty.MatchString(lines[i]) {
		if k := enclosingLine(lines, i); k >= 0 {
			first, last = k, blockEnd(lines, k, lay)
		}
	}
	return strings.Join(lines[first:last+1], "\n")
}

// noAssertRE is the waiver of a test with no assertion, with its reason, in the test's
// body or on the line above it: `@no-assert: <why>`.
var noAssertRE = regexp.MustCompile(`@no-assert[^\S\n]*:[^\S\n]*\S+`)

func indentOf(l string) int { return len(l) - len(strings.TrimLeft(l, " \t")) }

// blockEnd is the last line of the block that line i opens, by its layout: a block the
// dialect's `block_open` opens closes on the first line back at the opener's indentation
// that `block_close` reads and that opens nothing again; any other block ends on the first
// line back at that indentation. A line that starts inside a string spanning lines is text,
// and never ends the block.
func blockEnd(lines []string, i int, lay layout) int {
	base := indentOf(lines[i])
	bracket := lay.open != nil && lay.open.MatchString(lines[i])
	next := i + 1
	for next < len(lines) && strings.TrimSpace(lines[next]) == "" {
		next++
	}
	if !bracket && (next >= len(lines) || indentOf(lines[next]) <= base) {
		return i // opened and closed on its own line
	}
	for j := i + 1; j < len(lines); j++ {
		inside := j < len(lay.inString) && lay.inString[j]
		if inside || strings.TrimSpace(lines[j]) == "" || indentOf(lines[j]) > base {
			continue
		}
		if !bracket {
			return j - 1
		}
		// A line that closes and opens again goes on with the block: the table of an
		// `it.each([` closes into the test's own body (`])('title', () => {`), and a
		// `} else {` into the next branch. Only a closer that opens nothing ends it.
		if lay.close != nil && lay.close.MatchString(lines[j]) && !lay.open.MatchString(lines[j]) {
			return j
		}
	}
	return len(lines) - 1
}

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

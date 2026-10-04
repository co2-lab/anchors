// @anchors
//   code: THATT
//   ref: THSAS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testlist"
)

// assertCfg reads Go tests by the family's `t.Run(` and asserts with `t.Errorf`/`t.Fatal`.
func assertCfg(labels bool, tests *config.TestsSource) *config.Config {
	return &config.Config{
		Dialect: &config.Dialect{Family: "go", Tests: tests},
		Gates:   []config.Gate{{Name: "a", Check: "test-has-assertion", On: []string{"test"}, Labels: labels}},
	}
}

// runAssertion writes one test file and confronts it.
func runAssertion(t *testing.T, cfg *config.Config, src string) (Verdict, string) {
	t.Helper()
	resetProjectTestsCache()
	root := t.TempDir()
	writeFile(t, root, "a_test.go", src)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a_test.go", Kind: mapx.KindTest}}}
	return checkTestHasAssertion(src, g.Nodes[0], root, g, cfg)
}

func TestTestHasAssertion_aTestWithNoAssertionFails(t *testing.T) {
	t.Run("THSAS-B01: A test with no assertion fails, named by its line and title", func(t *testing.T) {})
	src := "func TestX(t *testing.T) {\n" +
		"\tt.Run(\"verifies\", func(t *testing.T) {\n\t\tif f() != 1 {\n\t\t\tt.Errorf(\"no\")\n\t\t}\n\t})\n" +
		"\tt.Run(\"only calls\", func(t *testing.T) {\n\t\tf()\n\t})\n}\n"
	v, msg := runAssertion(t, assertCfg(false, nil), src)
	if v != Fail || !strings.Contains(msg, "7: only calls") || strings.Contains(msg, "verifies") || !strings.Contains(msg, "1 test") {
		t.Fatalf("the test that only calls fails by line and title, got %v: %s", v, msg)
	}
	ok := "\tt.Run(\"a\", func(t *testing.T) {\n\t\tt.Fatal(\"x\")\n\t})\n"
	if v, msg := runAssertion(t, assertCfg(false, nil), ok); v != Pass {
		t.Errorf("a file whose tests assert passes, got %v: %s", v, msg)
	}
}

func TestTestHasAssertion_theBodyIsTheBlock(t *testing.T) {
	t.Run("THSAS-B02: The body is the block the test's line opens", func(t *testing.T) {})
	lines := strings.Split("it('a', () => {\n  x()\n  }\n})\nafter()\n", "\n")
	if got := blockEnd(lines, 0); got != 3 {
		t.Errorf("a bracket block ends on its closer at the opener's indentation, got line %d", got)
	}
	lines = strings.Split("  f(\n    a,\n  )\n", "\n")
	if got := blockEnd(lines, 0); got != 2 {
		t.Errorf("an indented closer at the opener's indentation ends it, got line %d", got)
	}
	lines = strings.Split("def test_a():\n    assert x\n\n    y()\ndef test_b():\n", "\n")
	if got := blockEnd(lines, 0); got != 3 {
		t.Errorf("a block opened by anything else ends before the line back at its indentation, got line %d", got)
	}
	lines = strings.Split("t.Run(\"a\", f)\nt.Run(\"b\", g)\n", "\n")
	if got := blockEnd(lines, 0); got != 0 {
		t.Errorf("a line whose next line is not deeper is a block by itself, got line %d", got)
	}
	lines = strings.Split("t.Run(\"a\", f)\n\n\n", "\n")
	if got := blockEnd(lines, 0); got != 0 {
		t.Errorf("a line followed only by blank lines is a block by itself, got line %d", got)
	}
	lines = strings.Split("describe('x', () => {\n  it('a')\n", "\n")
	if got := blockEnd(lines, 0); got != len(lines)-1 {
		t.Errorf("a block never closed runs to the end, got line %d", got)
	}
	lines = strings.Split("loop do\n  x\nend\n", "\n")
	if got := blockEnd(lines, 0); got != 2 {
		t.Errorf("a `do` block ends on its `end`, got line %d", got)
	}
	lines = strings.Split("f(\n  a\n)\n", "\n")
	if body := testBody(lines, testlist.Test{Line: 1}, false); body != "f(\n  a\n)" {
		t.Errorf("the body runs from the test's line to its end, got %q", body)
	}
	if body := testBody(lines, testlist.Test{Line: len(lines) + 1}, false); body != "" {
		t.Errorf("the line right after the file has no body, got %q", body)
	}
	lines = strings.Split("def test_a():\n\n    assert x\ndef b():\n", "\n")
	if got := blockEnd(lines, 0); got != 2 {
		t.Errorf("blank lines before a deeper body do not close the block, got line %d", got)
	}
	if body := testBody(lines, testlist.Test{Line: 9}, false); body != "" {
		t.Errorf("a line outside the file has no body, got %q", body)
	}
}

func TestTestHasAssertion_literalsAreText(t *testing.T) {
	t.Run("THSAS-B03: A multi-line literal is text, not layout", func(t *testing.T) {})
	src := "func TestX(t *testing.T) {\n\tsrc := `\nfunc f() {\n}\n`\n\tpy := \"\"\"\n}\n\"\"\"\n\tq := \"`\"\n\tt.Errorf(q)\n}\nfunc other() {}\n"
	lines := strings.Split(src, "\n")
	if got := blockEnd(lines, 0); got != 10 {
		t.Fatalf("the block runs past the literals to its closer, got line %d", got)
	}
	single := strings.Split("f(x, (\n'''\n)\n'''\n  y\n)\n", "\n")
	if got := blockEnd(single, 0); got != 5 {
		t.Errorf("a triple single quote is a literal too, got line %d", got)
	}
}

func TestTestHasAssertion_labels(t *testing.T) {
	t.Run("THSAS-B04: An empty test is a label only when the gate says so", func(t *testing.T) {})
	src := "func TestX(t *testing.T) {\n\tt.Run(\"CODE-B01: labelled\", func(t *testing.T) {})\n\tif f() != 1 {\n\t\tt.Errorf(\"no\")\n\t}\n}\n"
	if v, msg := runAssertion(t, assertCfg(true, nil), src); v != Pass {
		t.Errorf("with labels the block around asserts for the label, got %v: %s", v, msg)
	}
	if v, _ := runAssertion(t, assertCfg(false, nil), src); v != Fail {
		t.Errorf("without labels the empty test has no assertion, got %v", v)
	}
	two := "func TestX(t *testing.T) {\n\tt.Run(\"A-B01: one\", func(t *testing.T) {})\n\tt.Run(\"A-B02: two\", func(t *testing.T) {})\n\tt.Errorf(\"no\")\n}\n"
	if v, msg := runAssertion(t, assertCfg(true, nil), two); v != Pass {
		t.Errorf("two labels in a row both stand for the function, got %v: %s", v, msg)
	}
	last := "func TestX(t *testing.T) {\n\tt.Errorf(\"no\")\n\tt.Run(\"A-B01: last\", func(t *testing.T) {})\n}\n"
	if v, msg := runAssertion(t, assertCfg(true, nil), last); v != Pass {
		t.Errorf("the block around is found above the label, not below, got %v: %s", v, msg)
	}
	top := "t.Run(\"top\", func(t *testing.T) {})\nt.Errorf(\"x\")\n"
	if v, _ := runAssertion(t, assertCfg(true, nil), top); v != Fail {
		t.Errorf("a label at the top level has no block around it, got %v", v)
	}
	notEmpty := "func TestX(t *testing.T) {\n\tt.Run(\"a\", run)\n\tt.Errorf(\"x\")\n}\n"
	if v, _ := runAssertion(t, assertCfg(true, nil), notEmpty); v != Fail {
		t.Errorf("a one-line test that is not empty is not a label, got %v", v)
	}
}

func TestTestHasAssertion_theSourcesEnd(t *testing.T) {
	t.Run("THSAS-B05: The source's end bounds the body", func(t *testing.T) {})
	src := "t.Run(\"a\", func(t *testing.T) {\n\tf()\n\tt.Errorf(\"x\")\n})\n"
	script := `printf '{"version":1,"tests":[{"file":"a_test.go","line":1,"end":2,"title":"a"}]}'`
	cfg := assertCfg(false, &config.TestsSource{Script: script, Assertion: `\bt\.Errorf\b`})
	if v, msg := runAssertion(t, cfg, src); v != Fail || !strings.Contains(msg, "1: a") {
		t.Errorf("the body ends where the script says, got %v: %s", v, msg)
	}
	script = `printf '{"version":1,"tests":[{"file":"a_test.go","line":1,"end":99,"title":"a"}]}'`
	cfg = assertCfg(false, &config.TestsSource{Script: script, Assertion: `\bt\.Errorf\b`})
	if v, msg := runAssertion(t, cfg, src); v != Pass {
		t.Errorf("an end past the file stops at its last line, got %v: %s", v, msg)
	}
}

func TestTestHasAssertion_skips(t *testing.T) {
	t.Run("THSAS-B06: Nothing to measure is skipped", func(t *testing.T) {})
	cfg := assertCfg(false, nil)
	if v, _ := checkTestHasAssertion("", mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}, "", nil, cfg); v != Skip {
		t.Errorf("a spec is skipped, got %v", v)
	}
	if v, _ := checkTestHasAssertion("", mapx.Node{ID: "h_test.go", Kind: mapx.KindTest, Support: true}, "", nil, cfg); v != Skip {
		t.Errorf("a support file is skipped, got %v", v)
	}
	if v, _ := runAssertion(t, &config.Config{Dialect: &config.Dialect{Family: "rust"}}, "x"); v != Skip {
		t.Errorf("no tests source is skipped, got %v", v)
	}
	noAssert := &config.Config{Dialect: &config.Dialect{Tests: &config.TestsSource{Pattern: `\bt\.Run\(`}}}
	if v, _ := runAssertion(t, noAssert, "t.Run(\"a\", f)\n"); v != Skip {
		t.Errorf("no assertion declared is skipped, got %v", v)
	}
	if v, _ := runAssertion(t, cfg, "package a\n"); v != Skip {
		t.Errorf("a file with no listed test is skipped, got %v", v)
	}
}

func TestTestHasAssertion_unlistableTestsFail(t *testing.T) {
	t.Run("THSAS-E01: Tests that cannot be listed fail the gate with the reason", func(t *testing.T) {})
	cfg := assertCfg(false, &config.TestsSource{Script: "exit 3", Assertion: `x`})
	if v, msg := runAssertion(t, cfg, "t.Run(\"a\", f)\n"); v != Fail || !strings.Contains(msg, "exit 3") {
		t.Errorf("a failing script fails the gate naming it, got %v: %s", v, msg)
	}
	cfg = assertCfg(false, &config.TestsSource{Pattern: `\bt\.Run\(`, Assertion: `expect(`})
	if v, msg := runAssertion(t, cfg, "t.Run(\"a\", f)\n"); v != Fail || !strings.Contains(msg, "dialect.tests.assertion") {
		t.Errorf("an assertion that does not compile fails naming the field, got %v: %s", v, msg)
	}
}

func TestTestHasAssertion_declared(t *testing.T) {
	t.Run("THSAS-B07: A test that declares it asserts nothing is left out", func(t *testing.T) {})
	src := "\t// @no-assert: only proves the call does not panic\n\tt.Run(\"above\", func(t *testing.T) {\n\t\tf()\n\t})\n" +
		"\tt.Run(\"inside\", func(t *testing.T) {\n\t\tf() // @no-assert: smoke run of the fixture\n\t})\n" +
		"\tt.Run(\"bare\", func(t *testing.T) {\n\t\tf() // @no-assert:\n\t})\n"
	v, msg := runAssertion(t, assertCfg(false, nil), src)
	if v != Fail || !strings.Contains(msg, "1 test(s)") || !strings.Contains(msg, "8: bare") {
		t.Errorf("only the bare waiver fails, got %v: %s", v, msg)
	}
	first := "t.Run(\"top\", func(t *testing.T) {\n\tf()\n})\n"
	if v, _ := runAssertion(t, assertCfg(false, nil), first); v != Fail {
		t.Errorf("a test on the first line has no line above to waive it, got %v", v)
	}
}

func TestTestHasAssertion_aLineBeyondTheContent(t *testing.T) {
	t.Run("THSAS-B08: A test at a line the content does not have", func(t *testing.T) {})
	lines := []string{"one", "two"}
	for _, line := range []int{3, 4, 50} {
		if got := testBody(lines, testlist.Test{Line: line}, false); got != "" {
			t.Errorf("line %d reads an empty body, got %q", line, got)
		}
	}
	resetProjectTestsCache()
	defer resetProjectTestsCache()
	root := t.TempDir()
	writeFile(t, root, "a_test.go", strings.Repeat("\n", 40)+"\tt.Run(\"far\", func(t *testing.T) {})\n")
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a_test.go", Kind: mapx.KindTest}}}
	if v, msg := checkTestHasAssertion("short\n", g.Nodes[0], root, g, assertCfg(false, nil)); v != Fail || !strings.Contains(msg, "far") {
		t.Errorf("the test beyond the content is read as empty, not a break, got %v: %s", v, msg)
	}
}

func TestTestHasAssertion_aCloserThatOpensGoesOn(t *testing.T) {
	t.Run("THSAS-B09: A closer that opens again goes on with the block", func(t *testing.T) {})
	src := "it.each([\n  [1, 2],\n  [3, 4],\n])('adds %i', async (a, b) => {\n  expect(a + 1).toBe(b)\n})\nafter()\n"
	lines := strings.Split(src, "\n")
	if got := blockEnd(lines, 0); got != 5 {
		t.Errorf("the table closes into the body, and the block ends at its closer, got line %d", got)
	}
	asConst := strings.Replace(src, "])(", "] as const)(", 1)
	if got := blockEnd(strings.Split(asConst, "\n"), 0); got != 5 {
		t.Errorf("a table closed `as const` goes on too, got line %d", got)
	}
	branch := strings.Split("if (x) {\n  a()\n} else {\n  b()\n}\nc()\n", "\n")
	if got := blockEnd(branch, 0); got != 4 {
		t.Errorf("an else goes on with the block, got line %d", got)
	}
	if body := testBody(lines, testlist.Test{Line: 1}, false); !strings.Contains(body, "expect(a + 1)") || strings.Contains(body, "after()") {
		t.Errorf("the body holds the assertion after the table, and stops at its closer, got:\n%s", body)
	}
}

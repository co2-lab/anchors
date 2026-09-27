package testlist

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// write puts files under a temporary root and returns it.
func write(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, body := range files {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const goRun = `\bt\.Run\(`

func TestPatternReadsCallsFollowedByALiteral(t *testing.T) {
	t.Run("TSTLS-B01: A pattern reads each call followed by a literal as a test", func(t *testing.T) {})
	root := write(t, map[string]string{"a_test.go": "package a\n\nfunc TestA(t *testing.T) {\n\tt.Run(\"ABCDX-B01: first\", func(t *testing.T) {})\n\tt.Run( \"ABCDX-B02: second\", nil)\n}\n"})
	got, err := List(root, []string{"a_test.go"}, Source{Pattern: goRun})
	want := []Test{{File: "a_test.go", Line: 4, Title: "ABCDX-B01: first"}, {File: "a_test.go", Line: 5, Title: "ABCDX-B02: second"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v (%v)", want, got, err)
	}
}

func TestTitlesQuotedThreeWays(t *testing.T) {
	t.Run("TSTLS-B02: Titles may be quoted three ways, with escapes", func(t *testing.T) {})
	src := "it('single \\'quoted\\'', f)\nit(\"double \\\"x\\\"\", f)\nit(`back\ntick`, f)\nit('broken\nline', f)\n"
	root := write(t, map[string]string{"a.test.ts": src})
	got, err := List(root, []string{"a.test.ts"}, Source{Pattern: `\bit\(`})
	if err != nil {
		t.Fatal(err)
	}
	var titles []string
	for _, x := range got {
		titles = append(titles, x.Title)
	}
	want := []string{"single 'quoted'", `double "x"`, "back\ntick"}
	if !reflect.DeepEqual(titles, want) {
		t.Fatalf("want %q, got %q", want, titles)
	}
}

func TestCallWithoutLiteralIsLeftOut(t *testing.T) {
	t.Run("TSTLS-B03: A call without a literal title is left out", func(t *testing.T) {})
	root := write(t, map[string]string{"a.test.ts": "it(name, f)\nit(\n  'real', f)\nit("})
	// The pattern swallows the blanks after the parenthesis, across the line break: the
	// line is still the one where the call starts.
	got, err := List(root, []string{"a.test.ts"}, Source{Pattern: `\bit\(\s*`})
	if err != nil || len(got) != 1 || got[0].Title != "real" || got[0].Line != 2 {
		t.Fatalf("only the literal title counts, at the line its call starts; got %v (%v)", got, err)
	}
}

func TestPatternOrder(t *testing.T) {
	t.Run("TSTLS-B04: The pattern's tests come in file order then position", func(t *testing.T) {})
	root := write(t, map[string]string{
		"b.test.ts": "it('b1', f)\nit('b2', f)\n",
		"a.test.ts": "it('a1', f)\nit('a2', f)\n",
	})
	got, err := List(root, []string{"b.test.ts", "a.test.ts"}, Source{Pattern: `\bit\(`})
	var titles []string
	for _, x := range got {
		titles = append(titles, x.File+":"+x.Title)
	}
	want := []string{"a.test.ts:a1", "a.test.ts:a2", "b.test.ts:b1", "b.test.ts:b2"}
	if err != nil || !reflect.DeepEqual(titles, want) {
		t.Fatalf("want %v, got %v (%v)", want, titles, err)
	}
}

func TestScriptOutputIsTheList(t *testing.T) {
	t.Run("TSTLS-B05: A script's contract output is the list of tests", func(t *testing.T) {})
	root := write(t, map[string]string{"marker": "here"})
	// The script proves it runs at the root by reading a file only the root has.
	script := `test -f marker && printf '{"version":1,"tests":[{"file":"a.test.ts","line":3,"title":"one"},{"file":"b.test.ts","line":9,"title":"two"}]}'`
	got, err := List(root, nil, Source{Script: script})
	want := []Test{{File: "a.test.ts", Line: 3, Title: "one"}, {File: "b.test.ts", Line: 9, Title: "two"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v (%v)", want, got, err)
	}
}

func TestOutputOutsideTheContract(t *testing.T) {
	t.Run("TSTLS-B06: Output outside the contract is refused naming what is wrong", func(t *testing.T) {})
	ok := `{"file":"a.ts","line":1,"title":"t"}`
	for _, c := range []struct{ out, why string }{
		{`not json`, "not the contract's JSON"},
		{`{"version":1,"tests":[]} {}`, "more than one JSON value"},
		{`{"version":1,"tests":[],"extra":1}`, "unknown field"},
		{`{"version":1,"tests":[{"file":"a.ts","line":1,"title":"t","kind":"x"}]}`, "unknown field"},
		{`{"tests":[` + ok + `]}`, "`version` must be 1"},
		{`{"version":2,"tests":[` + ok + `]}`, "`version` must be 1"},
		{`{"version":1}`, "`tests` is missing"},
		{`{"version":1,"tests":[{"file":"","line":1,"title":"t"}]}`, "tests[0]: `file` is empty"},
		{`{"version":1,"tests":[{"file":"/abs/a.ts","line":1,"title":"t"}]}`, "not relative"},
		{`{"version":1,"tests":[{"file":"../out.ts","line":1,"title":"t"}]}`, "not relative"},
		{`{"version":1,"tests":[` + ok + `,{"file":"a.ts","line":0,"title":"t"}]}`, "tests[1]: `line` must count from 1"},
		{`{"version":1,"tests":[{"file":"a.ts","line":1,"title":"  "}]}`, "`title` is empty"},
	} {
		if _, err := Parse([]byte(c.out)); err == nil || !strings.Contains(err.Error(), c.why) {
			t.Errorf("%s: want an error naming %q, got %v", c.out, c.why, err)
		}
	}
	if got, err := Parse([]byte(`{"version":1,"tests":[]}`)); err != nil || len(got) != 0 {
		t.Errorf("an empty list is inside the contract, got %v (%v)", got, err)
	}
	root := write(t, map[string]string{})
	if _, err := List(root, nil, Source{Script: `echo '{"version":2,"tests":[]}'`}); err == nil ||
		!strings.Contains(err.Error(), "answered outside the contract") || !strings.Contains(err.Error(), "`version` must be 1") {
		t.Errorf("a script outside the contract must be named with the violation, got %v", err)
	}
}

func TestScriptPathsAreNormalised(t *testing.T) {
	t.Run("TSTLS-B07: A script's file paths are normalised to the map's form", func(t *testing.T) {})
	got, err := Parse([]byte(`{"version":1,"tests":[{"file":"src/./x/../a.test.ts","line":1,"title":"t"}]}`))
	if err != nil || len(got) != 1 || got[0].File != "src/a.test.ts" {
		t.Fatalf("want src/a.test.ts, got %v (%v)", got, err)
	}
}

func TestNothingDeclaredListsNothing(t *testing.T) {
	t.Run("TSTLS-B08: A source that declares nothing lists nothing", func(t *testing.T) {})
	if (Source{}).Declared() {
		t.Error("an empty source declares nothing")
	}
	if !(Source{Pattern: "x"}).Declared() || !(Source{Script: "x"}).Declared() {
		t.Error("a pattern or a script is a declaration")
	}
	got, err := List(t.TempDir(), []string{"a_test.go"}, Source{})
	if err != nil || got != nil {
		t.Fatalf("want no test and no error, got %v (%v)", got, err)
	}
}

func TestBothSourcesRefused(t *testing.T) {
	t.Run("TSTLS-E01: A source with both a pattern and a script is refused", func(t *testing.T) {})
	if _, err := List(t.TempDir(), nil, Source{Pattern: "x", Script: "true"}); err == nil || !strings.Contains(err.Error(), "declare one") {
		t.Fatalf("want the conflict named, got %v", err)
	}
}

func TestPatternThatDoesNotCompile(t *testing.T) {
	t.Run("TSTLS-E02: A pattern that does not compile is refused", func(t *testing.T) {})
	if _, err := List(t.TempDir(), nil, Source{Pattern: "it("}); err == nil || !strings.Contains(err.Error(), "does not compile") {
		t.Fatalf("want the compile error named, got %v", err)
	}
}

func TestFailingScript(t *testing.T) {
	t.Run("TSTLS-E03: A failing script is an error naming why", func(t *testing.T) {})
	_, err := List(t.TempDir(), nil, Source{Script: `echo noise >&2; echo "jest not installed" >&2; exit 3`})
	if err == nil || !strings.Contains(err.Error(), "jest not installed") || !strings.Contains(err.Error(), "exit status 3") {
		t.Fatalf("want the command, the exit and the reason named, got %v", err)
	}
}

// @anchors
//   code: ICTNA
//   ref: INCTN

package initx

import (
	"strings"
	"testing"
)

func TestTestConventionOf(t *testing.T) {
	t.Run("INCTN-B06: A file is a test when its name carries a convention's prefix and suffix", func(t *testing.T) {})
	for name, want := range map[string]bool{
		"foo_test.go": true, "test_foo.py": true, "Login.spec.tsx": true, "UserTest.java": true,
		"_test.go": false, "foo.go": false,
	} {
		if _, got := testConventionOf(name); got != want {
			t.Errorf("%s: test = %v, want %v", name, got, want)
		}
	}
	if c, _ := testConventionOf("Login.spec.tsx"); c.Suffix != ".spec.tsx" {
		t.Errorf("Login.spec.tsx follows %q, want .spec.tsx", c.Suffix)
	}
}

func TestConventionGlobAndTemplate(t *testing.T) {
	t.Run("INCTN-B07: A convention gives its glob and its template", func(t *testing.T) {})
	t.Run("INCTN-B08: A test file's unit name drops the convention's prefix and suffix", func(t *testing.T) {})
	goC := TestConvention{Suffix: "_test.go"}
	py := TestConvention{Prefix: "test_", Suffix: ".py"}
	if goC.Glob() != "**/*_test.go" || goC.Template() != "{{dir}}/{{name}}_test.go" {
		t.Errorf("go: %s %s", goC.Glob(), goC.Template())
	}
	if py.Glob() != "**/test_*.py" || py.Template() != "{{dir}}/test_{{name}}.py" {
		t.Errorf("python: %s %s", py.Glob(), py.Template())
	}
	if py.stemOfTest("test_foo.py") != "foo" || goC.stemOfTest("foo_test.go") != "foo" {
		t.Error("the unit name should be foo")
	}
}

func TestConventionsFallBackToTheFamily(t *testing.T) {
	t.Run("INCTN-B09: With no convention read, the family default is used", func(t *testing.T) {})
	if got := (&Proposal{Family: "python"}).TestGlobs(); len(got) != 1 || got[0] != "**/test_*.py" {
		t.Errorf("python with no test: %v", got)
	}
	p := &Proposal{}
	if len(p.TestGlobs()) != 0 || p.TestTemplate() != "" {
		t.Errorf("no family, no test: %v %q", p.TestGlobs(), p.TestTemplate())
	}
}

func TestCoverageHintByFamily(t *testing.T) {
	t.Run("INCTN-B10: A family's coverage hint names the reports ingest reads", func(t *testing.T) {})
	if h := CoverageHint("go"); !strings.Contains(h, "lcov") || !strings.Contains(h, "junit") {
		t.Errorf("go hint = %q", h)
	}
	if h := CoverageHint("python"); !strings.Contains(h, "pytest") {
		t.Errorf("python hint = %q", h)
	}
	if CoverageHint("php") != "" || CoverageHint("") != "" {
		t.Error("a family with no hint, or no family, gets none")
	}
}

func TestGateScriptExampleByFamily(t *testing.T) {
	t.Run("INCTN-B11: A family's gate script runs in its own language", func(t *testing.T) {})
	for fam, want := range map[string]string{
		"go": "go run ./tools/gates <gate>", "ts": "node tools/gates.mjs <gate>", "python": "python -m tools.gates <gate>",
		"cobol": "a program in the project's own language",
	} {
		if got := GateScriptExample(fam); got != want {
			t.Errorf("%s: %q, want %q", fam, got, want)
		}
	}
}

func TestNoConventionIsShadowed(t *testing.T) {
	t.Run("INCTN-I03: No convention of the catalog is shadowed by a shorter one", func(t *testing.T) {})
	for i, a := range testConventions {
		for _, b := range testConventions[i+1:] {
			// b is shadowed when every name b matches is also matched by a, listed first.
			if strings.HasSuffix(b.Suffix, a.Suffix) && strings.HasPrefix(b.Prefix, a.Prefix) && a != b {
				t.Errorf("%s*%s comes before %s*%s and shadows it", a.Prefix, a.Suffix, b.Prefix, b.Suffix)
			}
		}
	}
}

func TestFamilyDefaultsAreInTheCatalog(t *testing.T) {
	t.Run("INCTN-I04: Every family default is a convention of the catalog", func(t *testing.T) {})
	for fam, d := range familyDefault {
		found := false
		for _, c := range testConventions {
			if c == d && c.Family == fam {
				found = true
			}
		}
		if !found {
			t.Errorf("the default of %s (%s*%s) is not in the catalog", fam, d.Prefix, d.Suffix)
		}
	}
}

func TestCatalogNamesNoFolder(t *testing.T) {
	t.Run("INCTN-X02: The catalog names no folder", func(t *testing.T) {})
	for _, c := range testConventions {
		if strings.Contains(c.Prefix+c.Suffix, "/") {
			t.Errorf("convention %s*%s names a folder", c.Prefix, c.Suffix)
		}
	}
	for _, m := range manifestFamily {
		if strings.Contains(m.file, "/") {
			t.Errorf("manifest %s names a folder", m.file)
		}
	}
}

// The instruction SAYS WHAT TO ANSWER, and says NOT to pass.
//
// A text that only mentioned `@TBD` would satisfy the judgment-gate test without solving
// the problem: the whole point is that the answer is not `pass`.
func TestTBDInstructionForbidsPassAndNamesTheAbsence(t *testing.T) {
	t.Run("INCTN-B04: The @TBD instruction forbids pass, orders a waiver naming the absence, and names the piece asked about", func(t *testing.T) {})
	t.Run("INCTN-I02: The @TBD instruction demands checking that the @TBD is still true", func(t *testing.T) {})
	got := tbdInstruction("the code")
	for _, required := range []string{"@TBD", "`waived`", "`pass`", "naming the absence"} {
		if !strings.Contains(got, required) {
			t.Errorf("the instruction does not mention %q:\n%s", required, got)
		}
	}
	// The piece goes into the text: without it the instruction would speak of "the code"
	// in a gate that asks about a test.
	if !strings.Contains(tbdInstruction("the test"), "the test") {
		t.Error("the instruction does not use the piece it received")
	}
	// The stale `@TBD` is the other half: a piece that came to exist with the marker still
	// in the file makes every gate that reads it waive what it should charge.
	if !strings.Contains(got, "stale") {
		t.Errorf("the instruction does not cover the `@TBD` that stopped being true:\n%s", got)
	}
}

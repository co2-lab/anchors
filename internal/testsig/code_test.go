package testsig

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// The package's DEFAULT has to match the letters `config` declares as canonical.
//
// Measured in the reference app: `testsig` had `SRVAXBNMD` and `config` had `SRVAXBNMDEIQF`. A
// project without `rule_types:` uses the `config` default — which includes `I` for
// Invariant — and this package did not recognize `GLCGL-I01` in the JUnit case name.
//
// The symptom: 11 green cases, three scenarios proven, and the three invariants showing
// "no green test" with a passing test. False traceability in the most dangerous direction —
// the requirement looks uncovered when it is proven, and whoever looks writes the test again.
func TestCodesInCase_recognizesInvariant(t *testing.T) {
	got := CodesInCase("GLCGL-I01: item `done` with the artifact missing fails")

	if len(got) != 1 || got[0] != "GLCGL-I01" {
		t.Errorf("codes = %v — the `-I01` (Invariant) was not recognized", got)
	}
}

// The `config` letters that were missing here: E, I, Q, F. All four.
func TestCodesInCase_recognizesTheMissingLetters(t *testing.T) {
	t.Run("RCGRL-B01: A rule code of each canonical letter is recognized in a test name", func(t *testing.T) {})
	for _, c := range []struct{ name, code string }{
		{"ABCDX-E01: example", "ABCDX-E01"},
		{"ABCDX-I01: invariant", "ABCDX-I01"},
		{"ABCDX-Q01: open decision", "ABCDX-Q01"},
		{"ABCDX-W01: phase", "ABCDX-W01"},
	} {
		got := CodesInCase(c.name)
		if len(got) != 1 || got[0] != c.code {
			t.Errorf("%q → %v, want [%s]", c.name, got, c.code)
		}
	}
}

// And the canonical ones keep working: the fix must not have swapped the list.
func TestCodesInCase_keepsTheCanonicalOnes(t *testing.T) {
	for _, code := range []string{"ABCDX-B01", "ABCDX-R01", "ABCDX-S01", "ABCDX-V01"} {
		if got := CodesInCase(code + ": something"); len(got) != 1 || got[0] != code {
			t.Errorf("%s was not recognized: %v", code, got)
		}
	}
}

// And the GUARD against the divergence coming back.
//
// Two copies of the same list diverge at the first new letter. The `config` comment already
// records that "it is the third time the list falls behind a new letter" — and this time it
// fell four behind. This test is the ruler that prevents the fifth.
func TestRuleLetters_doesNotDivergeFromConfig(t *testing.T) {
	t.Run("RCGRL-I01: The default rule letters equal the configuration's default letters", func(t *testing.T) {})
	if ruleLetters != config.DefaultRuleLetters {
		t.Errorf("ruleLetters = %q, config declares %q\n"+
			"  Each letter missing here is a scenario that passes its test and shows\n"+
			"  'no green test' — false traceability in the direction that makes someone\n"+
			"  write the test again.", ruleLetters, config.DefaultRuleLetters)
	}
	for _, l := range config.DefaultRuleLetters {
		code := "ABCDX-" + string(l) + "01"
		if got := CodesInCase(code + ": x"); len(got) != 1 || got[0] != code {
			t.Errorf("the canonical letter %c is not recognized: %v", l, got)
		}
	}
}

// The project's vocabulary replaces the canonical one in the grammar builder: a declared
// letter is recognised, and an empty value (a config that did not declare it) keeps what
// was there.
func TestSetRuleLetters(t *testing.T) {
	t.Run("RCGRL-B04: Declared rule letters replace the vocabulary", func(t *testing.T) {})
	t.Run("RCGRL-B06: An empty declaration keeps the vocabulary in place", func(t *testing.T) {})
	saved := ruleLetters
	t.Cleanup(func() { ruleLetters = saved })

	SetRuleLetters("BZ")
	if got := mustCodeRE().FindAllString("ABCDX-Z01: project letter", -1); len(got) != 1 || got[0] != "ABCDX-Z01" {
		t.Errorf("a declared letter must be recognised, got %v", got)
	}
	if got := mustCodeRE().FindAllString("ABCDX-I01: not declared", -1); len(got) != 0 {
		t.Errorf("a letter outside the project's vocabulary must not match, got %v", got)
	}
	SetRuleLetters("")
	if ruleLetters != "BZ" {
		t.Errorf("an empty value must keep the letters, got %q", ruleLetters)
	}
}

// The code length is the project's too: a 3-character code matches only once declared.
func TestSetCodeLenPattern(t *testing.T) {
	t.Run("RCGRL-B05: A declared code length replaces the accepted identity length", func(t *testing.T) {})
	saved := codeLenPattern
	t.Cleanup(func() { codeLenPattern = saved })

	if got := mustCodeRE().FindAllString("ABC-B01: short code", -1); len(got) != 0 {
		t.Fatalf("with the default length a 3-char code must not match, got %v", got)
	}
	SetCodeLenPattern("{3}")
	if got := mustCodeRE().FindAllString("ABC-B01: short code", -1); len(got) != 1 || got[0] != "ABC-B01" {
		t.Errorf("after SetCodeLenPattern({3}) the code must match, got %v", got)
	}
	SetCodeLenPattern("")
	if codeLenPattern != "{3}" {
		t.Errorf("an empty value must keep the pattern, got %q", codeLenPattern)
	}
}

// A slug, a design-system item and the visual-regression marker are codes too.
func TestCodeForms(t *testing.T) {
	t.Run("RCGRL-B02: A rule code with a lowercase slug is recognized with the slug", func(t *testing.T) {
		if got := CodesInCase("ABCDX-B01-some-slug proves it"); len(got) != 1 || got[0] != "ABCDX-B01-some-slug" {
			t.Errorf("got %v, want [ABCDX-B01-some-slug]", got)
		}
	})
	t.Run("RCGRL-B03: Design-system and visual-regression codes are recognized", func(t *testing.T) {
		if got := CodesInCase("ABCDX-DS-Button-primary renders"); len(got) != 1 || got[0] != "ABCDX-DS-Button-primary" {
			t.Errorf("design-system: got %v", got)
		}
		if got := CodesInCase("ABCDX-VR matches"); len(got) != 1 || got[0] != "ABCDX-VR" {
			t.Errorf("visual regression: got %v", got)
		}
	})
	t.Run("RCGRL-X01: An identity too long or glued to a longer word is not a code", func(t *testing.T) {
		for _, name := range []string{"XABCDXY-B01 too long", "abcABCDX-B01 glued"} {
			if got := CodesInCase(name); len(got) != 0 {
				t.Errorf("%q must yield no code, got %v", name, got)
			}
		}
	})
}

func TestScenarioCodes_keepTheVariant(t *testing.T) {
	t.Run("RCGRL-B07: A scenario code keeps its variant", func(t *testing.T) {})
	if got := ScenarioCodesInCase("RDCHX-B02#02: under integration; and RDCHX-B03"); len(got) != 2 || got[0] != "RDCHX-B02#02" || got[1] != "RDCHX-B03" {
		t.Errorf("the variant is kept, got %v", got)
	}
	if ScenarioRoot("RDCHX-B02#02") != "RDCHX-B02" || ScenarioRoot("RDCHX-B03") != "RDCHX-B03" {
		t.Error("the rule is the code without its variant")
	}
	feature := "  @RDCHX-B02#01 @unit-level\n  @RDCHX-B02#02 @integration-level\n  @RDCHX-B03\n  @RDCHX-B03\n"
	if got := FeatureScenarios(feature); len(got) != 3 || got[1] != "RDCHX-B02#02" {
		t.Errorf("each tag once, with its variant, got %v", got)
	}
}

func TestRulesProven_eachScenario(t *testing.T) {
	t.Run("RCGRL-B08: A rule is proven only when each of its scenarios is", func(t *testing.T) {})
	declared := []string{"RDCHX-B02#01", "RDCHX-B02#02", "RDCHX-B03#01", "RDCHX-B03#02"}
	proven := []string{"RDCHX-B02#01", "RDCHX-B03#01", "RDCHX-B03#02", "RDCHX-B04"}
	got := RulesProven(proven, declared)
	if got["RDCHX-B02"] || !got["RDCHX-B03"] || !got["RDCHX-B04"] {
		t.Errorf("B02 misses a variant, B03 has both, B04 is declared by no feature, got %v", got)
	}
	if m := UnprovenScenarios("RDCHX-B02", proven, declared); len(m) != 1 || m[0] != "RDCHX-B02#02" {
		t.Errorf("the missing variant is named, got %v", m)
	}
	if RulesProven([]string{"RDCHX-B02"}, declared)["RDCHX-B02"] {
		t.Error("the rule's own code does not prove the variants its features declare")
	}
}

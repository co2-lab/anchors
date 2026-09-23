package flagx

import "testing"

// The three house styles say the same thing, and the grammar has to hear one statement.
// This is the finding that shaped the alias table: the tools disagree only on spelling.
func TestParse_theThreeSpellingsSayTheSame(t *testing.T) {
	groups := map[Op][]string{
		OpGte:        {">= 50", "gte 50", "NUM_GTE 50", "greaterThanOrEqual 50", "Greater Than Inclusive 50"},
		OpStartsWith: {`startsWith "beta"`, `STR_STARTS_WITH "beta"`, `starts with "beta"`},
		OpMatches:    {`matches "^v2"`, `REGEX "^v2"`, `matches regex "^v2"`},
		OpContains:   {`contains "beta"`, `STR_CONTAINS "beta"`, `includes "beta"`},
		OpEq:         {`= "on"`, `== "on"`, `eq "on"`, `Equal "on"`, `NUM_EQ "on"`},
	}
	for want, spellings := range groups {
		for _, raw := range spellings {
			c, err := Parse(raw)
			if err != nil {
				t.Errorf("Parse(%q) refused: %v", raw, err)
				continue
			}
			if c.Op != want {
				t.Errorf("Parse(%q).Op = %q, want %q", raw, c.Op, want)
			}
		}
	}
}

// The operand loses its quotes, and only one layer of them.
func TestParse_operand(t *testing.T) {
	cases := map[string]string{
		`= "on"`:        "on",
		`= 'on'`:        "on",
		">= 50":         "50",
		`contains "be"`: "be",
		`= "a \"b\" c"`: `a \"b\" c`,
	}
	for raw, want := range cases {
		c, err := Parse(raw)
		if err != nil {
			t.Fatalf("Parse(%q): %v", raw, err)
		}
		if c.Operand != want {
			t.Errorf("Parse(%q).Operand = %q, want %q", raw, c.Operand, want)
		}
	}
}

// The ABSENT case is the one that breaks in production and the one nobody writes. It
// compares against nothing, and so it does not go through the operand split.
func TestParse_absentAndPresent(t *testing.T) {
	for _, raw := range []string{"absent", "is not set", "not set", "unset", "missing", "ABSENT"} {
		c, err := Parse(raw)
		if err != nil {
			t.Errorf("Parse(%q) refused the absent case: %v", raw, err)
			continue
		}
		if c.Op != OpAbsent {
			t.Errorf("Parse(%q).Op = %q, want %q", raw, c.Op, OpAbsent)
		}
		if c.Operand != "" {
			t.Errorf("Parse(%q).Operand = %q, the absent case compares against nothing", raw, c.Operand)
		}
	}
	for _, raw := range []string{"present", "is set", "any value"} {
		c, err := Parse(raw)
		if err != nil || c.Op != OpPresent {
			t.Errorf("Parse(%q) = (%v, %v), want OpPresent", raw, c.Op, err)
		}
	}
}

// `>` must not swallow `>=`, and `starts` must not be read as an unknown operator carrying
// `with "beta"` as its operand.
func TestParse_longestPrefixWins(t *testing.T) {
	c, err := Parse(">= 50")
	if err != nil || c.Op != OpGte || c.Operand != "50" {
		t.Errorf(`Parse(">= 50") = (%q, %q, %v), want gte/50`, c.Op, c.Operand, err)
	}
	c, err = Parse(`starts with "beta"`)
	if err != nil || c.Op != OpStartsWith || c.Operand != "beta" {
		t.Errorf(`Parse("starts with ...") = (%q, %q, %v), want starts-with/beta`, c.Op, c.Operand, err)
	}
	c, err = Parse(`not contains "beta"`)
	if err != nil || c.Op != OpNotContain {
		t.Errorf(`Parse("not contains ...") = (%q, %v), want not-contains`, c.Op, err)
	}
}

// The grammar exists to REFUSE. One that accepts everything confronts nothing — that is
// the difference between a fixed grammar and prose, the decision that produced this
// package.
func TestParse_refusesWhatItDoesNotUnderstand(t *testing.T) {
	refuse := []string{
		"",
		"   ",
		"when the user is a beta tester", // prose
		"segmentMatch beta-users",        // depends on an external service, out of scope
		"modulo 3",                       // same
		">=",                             // an operator with nothing to compare against
		`contains`,                       // same
		"maybe",                          // a loose word
	}
	for _, raw := range refuse {
		if c, err := Parse(raw); err == nil {
			t.Errorf("Parse(%q) accepted it as %q/%q — it should refuse", raw, c.Op, c.Operand)
		}
	}
}

// NeedsOperand separates the two valueless operators from the rest.
func TestNeedsOperand(t *testing.T) {
	if OpAbsent.NeedsOperand() || OpPresent.NeedsOperand() {
		t.Error("absent and present compare against nothing")
	}
	for _, op := range []Op{OpEq, OpGte, OpContains, OpMatches, OpRollout} {
		if !op.NeedsOperand() {
			t.Errorf("%q compares against a value", op)
		}
	}
}

// The word for the ABSENT case comes from the translation catalog, never hardcoded. The
// first version knew only the English word, and a flag file written in Portuguese saying
// `ausente` came out as `ne` with operand `set` — a WRONG answer, silently, on precisely
// the scenario that matters most. (`ausente` and `presente` below are catalog data, not
// prose of this project.)
func TestParse_absentInEveryLanguage(t *testing.T) {
	for _, raw := range []string{"absent", "ausente", "present", "presente"} {
		c, err := Parse(raw)
		if err != nil {
			t.Errorf("Parse(%q) refused: %v", raw, err)
			continue
		}
		if c.Op != OpAbsent && c.Op != OpPresent {
			t.Errorf("Parse(%q).Op = %q — it fell into the operand split", raw, c.Op)
		}
		if c.Operand != "" {
			t.Errorf("Parse(%q).Operand = %q, it should be empty", raw, c.Operand)
		}
	}
}

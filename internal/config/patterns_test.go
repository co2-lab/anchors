// @anchors
//   ref: DRPTD

package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// The field accepts BOTH shapes: a string for the spec that governs one file (most of them)
// and a list for a configuration spec, which governs several. Requiring a list from all of
// them would be noise in exchange for nothing.
func TestPadroesAcceptsStringOrList(t *testing.T) {
	t.Run("DRPTD-B01: A single pattern written as text becomes a list of one", func(t *testing.T) {})
	t.Run("DRPTD-B02: A list of patterns is kept whole and in order", func(t *testing.T) {})
	var d Derived
	if err := yaml.Unmarshal([]byte(`
anchor: spec
files:
  code: "{{dir}}/{{name}}.ts"
  test:
    - "{{dir}}/{{name}}.test.ts"
    - "__tests__/{{name}}.test.ts"
`), &d); err != nil {
		t.Fatal(err)
	}
	if got := d.PadroesDe()["code"]; len(got) != 1 || got[0] != "{{dir}}/{{name}}.ts" {
		t.Errorf("a string should become a list of one, got %v", got)
	}
	if got := d.PadroesDe()["test"]; len(got) != 2 || got[0] != "{{dir}}/{{name}}.test.ts" || got[1] != "__tests__/{{name}}.test.ts" {
		t.Errorf("a list should keep both, in order, got %v", got)
	}
}

// The loader refuses what is neither text nor a non-empty list of text, and says why.
func TestPadroesRejectsOtherShapes(t *testing.T) {
	t.Run("DRPTD-B03: An empty list is refused, naming the cause", func(t *testing.T) {})
	t.Run("DRPTD-B04: Any shape other than text or a list of text is refused", func(t *testing.T) {})
	for name, doc := range map[string]string{
		"empty list":       "code: []\n",
		"mapping":          "code: {a: b}\n",
		"list of mappings": "code:\n  - {a: b}\n",
	} {
		var m map[string]Padroes
		if err := yaml.Unmarshal([]byte(doc), &m); err == nil {
			t.Errorf("%s: want an error, got %v", name, m)
		}
	}
	var m map[string]Padroes
	err := yaml.Unmarshal([]byte("code: []\n"), &m)
	if err == nil || !strings.Contains(err.Error(), "empty pattern list") {
		t.Errorf("an empty list must name its cause, got %v", err)
	}
}

// Writing back keeps the simplest shape: one pattern is a string, several are a list.
func TestPadroesMarshalsToTheSimplestShape(t *testing.T) {
	t.Run("DRPTD-B05: Written back, one pattern is text and several are a list", func(t *testing.T) {})
	out, err := yaml.Marshal(map[string]Padroes{
		"code": {"x.ts"},
		"test": {"a.test.ts", "b.test.ts"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := "code: x.ts\ntest:\n    - a.test.ts\n    - b.test.ts\n"; string(out) != want {
		t.Errorf("Marshal = %q, want %q", out, want)
	}
}

func TestPadroes_roundTrip(t *testing.T) {
	t.Run("DRPTD-I01: What is written back reads back as the same patterns", func(t *testing.T) {})
	for _, p := range []Padroes{{"x.ts"}, {"a.test.ts", "b.test.ts", "c.test.ts"}} {
		out, err := yaml.Marshal(map[string]Padroes{"code": p})
		if err != nil {
			t.Fatal(err)
		}
		var back map[string]Padroes
		if err := yaml.Unmarshal(out, &back); err != nil {
			t.Fatal(err)
		}
		if got := back["code"]; strings.Join(got, "|") != strings.Join(p, "|") {
			t.Errorf("%v came back as %v", p, got)
		}
	}
}

func TestPadroes_keptAsWritten(t *testing.T) {
	t.Run("DRPTD-X01: The patterns are kept as written, neither expanded nor checked as globs", func(t *testing.T) {})
	var m map[string]Padroes
	doc := "code:\n  - \"src/[bad\"\n  - \" spaced/*.ts \"\n"
	if err := yaml.Unmarshal([]byte(doc), &m); err != nil {
		t.Fatalf("reading the shape must not judge the glob, got %v", err)
	}
	if got := m["code"]; len(got) != 2 || got[0] != "src/[bad" || got[1] != " spaced/*.ts " {
		t.Errorf("patterns = %q, want them exactly as written", got)
	}
}

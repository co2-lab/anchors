package gate

import (
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// `non-empty` counts SCENARIOS for a feature, and a scenario opens in any language of the
// official Gherkin table — including the synonyms (`Example:` for `Scenario:`).
//
// The first version reused a regex with five keywords nailed in Portuguese and English,
// and it failed valid features on a BLOCKING gate: a Spanish `Escenario:` and an English
// `Rule:` + `Example:` both came out as "feature with no scenario".
func TestNonEmpty_featureScenarioInAnyGherkinLanguage(t *testing.T) {
	n := mapx.Node{Kind: mapx.KindFeature}
	pass := map[string]string{
		"pt":            "Funcionalidade: X\n\n  Cenário: a\n",
		"pt outline":    "Funcionalidade: X\n\n  Esquema do Cenário: a\n",
		"en outline":    "Feature: X\n\n  Scenario Outline: a\n",
		"en example":    "Feature: X\n  Rule: r\n    Example: a\n",
		"es":            "Característica: X\n\n  Escenario: a\n",
		"fr":            "Fonctionnalité: X\n\n  Scénario: a\n",
		"pt unaccented": "Funcionalidade: X\n\n  Cenario: a\n",
	}
	for name, content := range pass {
		if v, msg := checkNonEmpty(content, n); v != Pass {
			t.Errorf("%s: a valid scenario was not recognised (%v: %s)", name, v, msg)
		}
	}

	fail := map[string]string{
		// `Examples:` is the data table of an outline, not a scenario: the `:` right after
		// the keyword keeps `Example` from matching it.
		"examples table only": "Feature: X\n\n  Examples:\n    | a |\n",
		"header only":         "Funcionalidade: X\n# nothing else\n",
	}
	for name, content := range fail {
		if v, _ := checkNonEmpty(content, n); v != Fail {
			t.Errorf("%s: a feature with no scenario passed", name)
		}
	}
}

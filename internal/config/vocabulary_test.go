// @anchors
//   ref: GTVCG

package config

import (
	"reflect"
	"strings"
	"testing"
)

// THE ALIAS TABLE this file used to test WAS REMOVED.
//
// It accepted the Portuguese gate names and converted them at load, forever. The tests here
// proved the conversion worked — and it did; the problem was the design: the file never got
// fixed, and the map piled up stamps in both forms. Measured in the reference app: 40 judgments
// recorded as `regra-cumprida` living next to 2 as `rule-fulfilled`.
//
// There was also an asymmetry: READING normalised (`mapx.mesmoGate`), WRITING did not — a
// project that renamed the gate got a SECOND stamp instead of updating the first.
//
// The conversion became migration step `1→2` (`internal/migra/formato_2.go`), and the
// rulers worth keeping went with it:
//
//   - `internal/migra` tests the conversion itself (renames key and value, leaves comments
//     alone, is idempotent, respects the target file);
//   - `internal/initx/vocabulario_test.go` checks that every target of the mapping is a
//     gate that EXISTS, and that no default gate was born with a Portuguese name.
//
// This file stays as a record: whoever looks for the alias tests finds why they are gone.

func TestDefaultGateNamesForTest(t *testing.T) {
	t.Run("GTVCG-B01: With no source registered, the default gate names are absent, not a failure", func(t *testing.T) {})
	t.Run("GTVCG-B02: With a source registered, the default gate names are the ones it gives", func(t *testing.T) {})
	saved := defaultGateNames
	t.Cleanup(func() { defaultGateNames = saved })

	RegisterGateNames(nil)
	if got := DefaultGateNamesForTest(); got != nil {
		t.Fatalf("with nothing registered, want nil, got %v", got)
	}
	RegisterGateNames(func() []string { return []string{"spec-complete", "rule-fulfilled"} })
	if got, want := DefaultGateNamesForTest(), []string{"spec-complete", "rule-fulfilled"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("DefaultGateNamesForTest = %v, want %v", got, want)
	}
}

func TestRegisterGateNames_lastRegistrationWins(t *testing.T) {
	t.Run("GTVCG-I01: The answer always comes from the source registered last", func(t *testing.T) {})
	saved := defaultGateNames
	t.Cleanup(func() { defaultGateNames = saved })

	RegisterGateNames(func() []string { return []string{"first"} })
	RegisterGateNames(func() []string { return []string{"second"} })
	if got := DefaultGateNamesForTest(); !reflect.DeepEqual(got, []string{"second"}) {
		t.Fatalf("DefaultGateNamesForTest = %v, want the second source's list", got)
	}
}

func TestDefaultGateNamesForTest_asksTheSourceEachTime(t *testing.T) {
	t.Run("GTVCG-X01: The names are not cached: each question asks the registered source again", func(t *testing.T) {})
	saved := defaultGateNames
	t.Cleanup(func() { defaultGateNames = saved })

	calls := 0
	RegisterGateNames(func() []string { calls++; return []string{"g"} })
	DefaultGateNamesForTest()
	DefaultGateNamesForTest()
	if calls != 2 {
		t.Fatalf("the source was asked %d times for two questions, want 2", calls)
	}
}

func TestArtifactLetters(t *testing.T) {
	t.Run("GTVCG-B03: The letters of plans, flows and actions", func(t *testing.T) {})
	if PhaseLetter != "W" || StepLetter != "T" || OutcomeLetter != "O" {
		t.Fatalf("the letters are W, T and O, got %s %s %s", PhaseLetter, StepLetter, OutcomeLetter)
	}
	if !strings.Contains(DefaultRuleLetters, PhaseLetter) {
		t.Errorf("specs cite phases in needs:, so %s must be canonical: %s", PhaseLetter, DefaultRuleLetters)
	}
	for _, l := range []string{StepLetter, OutcomeLetter} {
		if strings.Contains(DefaultRuleLetters, l) {
			t.Errorf("%s would mean something else in a spec: %s", l, DefaultRuleLetters)
		}
	}
}

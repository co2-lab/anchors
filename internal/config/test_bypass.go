// @anchors
//   code: TBYCF
//   ref: TBYSA

package config

import (
	"fmt"
	"regexp"
	"strings"
)

// TestBypass is how the project's code skips an effect while a test runs — an early return
// under an end-to-end flag, a write left out under CI — and the suite layers whose proofs reach
// the path it skips.
//
// Measured in MIF (MIF-1S): a hook returned early under `EXPO_PUBLIC_E2E` so the visual captures
// stayed clean, and the only suite against the real backend never ran the write it skipped —
// the write that reached the store review broken.
type TestBypass struct {
	// Patterns recognise a bypass on a line of code, each a regular expression.
	Patterns []string `yaml:"patterns"`
	// ProvenBy are the suite layers (`tests[].layer`) whose proofs run the skipped path.
	ProvenBy []string `yaml:"proven_by"`
}

// PatternREs are the compiled patterns; one that does not compile is left out, which the
// configuration's load refuses.
func (b *TestBypass) PatternREs() []*regexp.Regexp {
	if b == nil {
		return nil
	}
	var out []*regexp.Regexp
	for _, p := range b.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			out = append(out, re)
		}
	}
	return out
}

// validateTestBypass refuses a declaration with no pattern, a pattern that does not compile,
// or no layer that proves the skipped path.
func validateTestBypass(b *TestBypass) error {
	if b == nil {
		return nil
	}
	if len(b.Patterns) == 0 {
		return fmt.Errorf("`test_bypass.patterns` is empty: say how a bypass is recognised")
	}
	for _, p := range b.Patterns {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("`test_bypass.patterns` holds an empty pattern")
		}
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("`test_bypass.patterns`: %q is no regular expression: %v", p, err)
		}
	}
	if len(b.ProvenBy) == 0 {
		return fmt.Errorf("`test_bypass.proven_by` names no suite layer: say which suites run the path a bypass skips")
	}
	return nil
}

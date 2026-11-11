// @anchors
//   code: EXSCF
//   ref: EXSSA

package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExternalSystem is a real system the project's code talks to — a database, an API, a store —
// whose doubles in the unit tests accept whatever shape they are given. A code file that
// calls it is proven only by a test that ran in a suite reaching the real system.
//
// Measured in MIF: fourteen defects crossed a boundary — an index that did not exist, a
// condition outside the key schema, an owner field the API refused —, each proven green
// against a double that accepted the wrong thing.
type ExternalSystem struct {
	// Name is the project's name for the system (appsync, stripe, s3).
	Name string `yaml:"name"`
	// Calls is the regular expression that recognises a call to it in a code file.
	Calls string `yaml:"calls"`
	// ProvenBy are the layers of the project's suites (`tests[].layer`) that reach the real
	// system: integration, live, contract.
	ProvenBy []string `yaml:"proven_by"`
	// Version is the version of the system's API the project speaks today (`2025-09-30`,
	// `v2`): a response recorded from another version is a double of a system that is gone.
	Version string `yaml:"version,omitempty"`
	// MaxAge is how old a recorded response may be before it is recorded again, in days
	// (`90d`); empty, a recording does not age.
	MaxAge string `yaml:"max_age,omitempty"`
}

// MaxAgeDays is MaxAge in days, 0 when it is not declared.
func (b ExternalSystem) MaxAgeDays() int {
	n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(b.MaxAge), "d"))
	return n
}

// CallsRE is the system's compiled pattern; nil when it does not compile, which the
// configuration's load refuses.
func (b ExternalSystem) CallsRE() *regexp.Regexp {
	re, err := regexp.Compile(b.Calls)
	if err != nil {
		return nil
	}
	return re
}

// validateExternalSystems refuses a system with no name or a repeated one, a pattern that does
// not compile, and one no suite layer proves.
func validateExternalSystems(bs []ExternalSystem) error {
	seen := map[string]bool{}
	for i, b := range bs {
		name := strings.TrimSpace(b.Name)
		if name == "" {
			return fmt.Errorf("`external_systems[%d]` has no name", i)
		}
		if seen[name] {
			return fmt.Errorf("`external_systems`: %q is declared twice", name)
		}
		seen[name] = true
		if strings.TrimSpace(b.Calls) == "" {
			return fmt.Errorf("`external_systems[%s].calls` is empty: say how a call to %s is recognised", name, name)
		}
		if _, err := regexp.Compile(b.Calls); err != nil {
			return fmt.Errorf("`external_systems[%s].calls` is no regular expression: %v", name, err)
		}
		if a := strings.TrimSpace(b.MaxAge); a != "" && !maxAgeRE.MatchString(a) {
			return fmt.Errorf("`external_systems[%s].max_age` must be a number of days (90d), got %q", name, b.MaxAge)
		}
		if len(b.ProvenBy) == 0 {
			return fmt.Errorf("`external_systems[%s].proven_by` names no suite layer: say which suites reach the real %s", name, name)
		}
	}
	return nil
}

var maxAgeRE = regexp.MustCompile(`^[1-9][0-9]*d$`)

// @anchors
//   code: EXSTA
//   ref: EXSSA

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadSystems(t *testing.T, systems string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), DefaultFile)
	if err := os.WriteFile(p, []byte("version: 1\nlayers: {}\nexternal_systems:\n"+systems), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestExternalSystems_declared(t *testing.T) {
	t.Run("EXSSA-B01: A system declares its calls, the layers that reach it, its version and a recording's age", func(t *testing.T) {})
	c, err := loadSystems(t, "  - name: stripe\n    calls: 'stripe\\.charges\\.create\\('\n    proven_by: [live, contract]\n    version: 2025-09-30.dahlia\n    max_age: 90d\n")
	if err != nil {
		t.Fatal(err)
	}
	s := c.ExternalSystems[0]
	if s.Name != "stripe" || s.CallsRE() == nil || !s.CallsRE().MatchString("stripe.charges.create({})") ||
		strings.Join(s.ProvenBy, ",") != "live,contract" || s.Version != "2025-09-30.dahlia" || s.MaxAgeDays() != 90 {
		t.Errorf("the system: %+v", s)
	}
}

func TestExternalSystems_refused(t *testing.T) {
	t.Run("EXSSA-E01: A system that cannot hold refuses the configuration", func(t *testing.T) {})
	for want, systems := range map[string]string{
		"has no name":           "  - calls: x\n    proven_by: [live]\n",
		"declared twice":        "  - name: a\n    calls: x\n    proven_by: [live]\n  - name: a\n    calls: y\n    proven_by: [live]\n",
		"calls` is empty":       "  - name: a\n    proven_by: [live]\n",
		"no regular expression": "  - name: a\n    calls: '('\n    proven_by: [live]\n",
		"names no suite layer":  "  - name: a\n    calls: x\n",
		"number of days":        "  - name: a\n    calls: x\n    proven_by: [live]\n    max_age: 3 months\n",
	} {
		if _, err := loadSystems(t, systems); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: want %q, got %v", systems, want, err)
		}
	}
}

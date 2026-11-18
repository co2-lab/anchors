// @anchors
//   code: TBYTA
//   ref: TBYSA

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadBypass(t *testing.T, block string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), DefaultFile)
	if err := os.WriteFile(p, []byte("version: 1\nlayers: {}\ntest_bypass:\n"+block), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestTestBypass_declared(t *testing.T) {
	t.Run("TBYSA-B01: The project declares its bypasses and the layers that run the skipped path", func(t *testing.T) {})
	c, err := loadBypass(t, "  patterns: ['E2E === .true.\\)\\s*return']\n  proven_by: [integration]\n")
	if err != nil {
		t.Fatal(err)
	}
	res := c.TestBypass.PatternREs()
	if len(res) != 1 || !res[0].MatchString("if (process.env.E2E === 'true') return") || strings.Join(c.TestBypass.ProvenBy, ",") != "integration" {
		t.Errorf("the bypass: %+v", c.TestBypass)
	}
}

func TestTestBypass_refused(t *testing.T) {
	t.Run("TBYSA-E01: A bypass declaration that cannot hold refuses the configuration", func(t *testing.T) {})
	for want, block := range map[string]string{
		"patterns` is empty":    "  proven_by: [integration]\n",
		"no regular expression": "  patterns: ['(']\n  proven_by: [integration]\n",
		"names no suite layer":  "  patterns: ['x']\n",
	} {
		if _, err := loadBypass(t, block); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: want %q, got %v", block, want, err)
		}
	}
}

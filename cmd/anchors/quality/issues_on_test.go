package quality

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// Which check writes issues, per mode: local always; github only in CI or when asked;
// manual only when asked.
func TestIssuesOnFor(t *testing.T) {
	manual := &config.Config{Workflow: &config.Workflow{Mode: config.ModeManual}}
	gh := &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub, Repo: "o/r", Labels: []string{"anchors"}}}
	for _, c := range []struct {
		name         string
		cfg          *config.Config
		record, inCI bool
		want         bool
	}{
		{"local", &config.Config{}, false, false, true},
		{"manual", manual, false, false, false},
		{"manual in CI", manual, false, true, false},
		{"manual --record-issues", manual, true, false, true},
		{"github local run", gh, false, false, false},
		{"github in CI", gh, false, true, true},
		{"github --record-issues", gh, true, false, true},
	} {
		if got := issuesOnFor(c.cfg, c.record, c.inCI); got != c.want {
			t.Errorf("%s: issuesOnFor = %v, want %v", c.name, got, c.want)
		}
	}
}

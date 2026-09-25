package mapcmd

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// In manual mode the judge writes no issue unless asked; every other mode keeps writing it.
func TestJudgeWritesIssue(t *testing.T) {
	manual := &config.Config{Workflow: &config.Workflow{Mode: config.ModeManual}}
	local := &config.Config{Workflow: &config.Workflow{Mode: config.ModeLocal}}
	for _, c := range []struct {
		name   string
		cfg    *config.Config
		record bool
		want   bool
	}{
		{"manual, no flag", manual, false, false},
		{"manual, --record-issues", manual, true, true},
		{"local", local, false, true},
		{"no workflow block (local)", &config.Config{}, false, true},
	} {
		if got := judgeWritesIssue(c.cfg, c.record); got != c.want {
			t.Errorf("%s: judgeWritesIssue = %v, want %v", c.name, got, c.want)
		}
	}
}

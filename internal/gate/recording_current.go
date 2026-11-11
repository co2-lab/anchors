// @anchors
//   code: RCCRC
//   ref: RCCRS

package gate

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// recording-current: a response recorded from a real system is a double of the system the
// project speaks to today.
//
// Measured in MIF: nine defects passed against doubles whose responses had been right once —
// an API moved to a new version (Stripe's `dahlia`), a capability changed, a presign
// condition was added — and the double kept answering the old shape. The project marks each
// recording with where it came from, `@recorded: stripe 2025-09-30.dahlia 2026-10-02`, and
// declares the version it speaks and how long a recording lives (`external_systems:`): a
// recording of another version fails, one older than its system's `max_age` diverges until it
// is recorded again.
func checkRecordingCurrent(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	marks := recordedRE.FindAllStringSubmatch(content, -1)
	if len(marks) == 0 {
		return Skip, i18n.T("gate.recording_current.skip_none")
	}
	systems := map[string]config.ExternalSystem{}
	if cfg != nil {
		for _, s := range cfg.ExternalSystems {
			systems[s.Name] = s
		}
	}
	now, _ := time.Parse("2006-01-02", today())
	var failed, aged []string
	for _, m := range marks {
		fields := strings.Fields(m[1])
		if len(fields) != 3 {
			failed = append(failed, i18n.T("gate.recording_current.malformed", strings.TrimSpace(m[0])))
			continue
		}
		name, version, date := fields[0], fields[1], fields[2]
		s, ok := systems[name]
		if !ok {
			failed = append(failed, i18n.T("gate.recording_current.undeclared", name))
			continue
		}
		at, err := time.Parse("2006-01-02", date)
		if err != nil {
			failed = append(failed, i18n.T("gate.recording_current.malformed", strings.TrimSpace(m[0])))
			continue
		}
		if s.Version != "" && version != s.Version {
			failed = append(failed, i18n.T("gate.recording_current.other_version", name, version, s.Version))
			continue
		}
		if max := s.MaxAgeDays(); max > 0 && !now.IsZero() {
			if days := int(now.Sub(at).Hours() / 24); days > max {
				aged = append(aged, i18n.T("gate.recording_current.aged", name, date, fmt.Sprint(days), s.MaxAge))
			}
		}
	}
	switch {
	case len(failed) > 0:
		return Fail, strings.Join(append(failed, aged...), "; ")
	case len(aged) > 0:
		return Diverge, strings.Join(aged, "; ")
	}
	return Pass, ""
}

// recordedRE reads a recording's mark: the system, the version and the date it was recorded.
var recordedRE = regexp.MustCompile(`(?m)@recorded[^\S\n]*:[^\S\n]*([^\n]*?)[^\S\n]*(?:\*/|-->|$)`)

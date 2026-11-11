// @anchors
//   code: RCCRT
//   ref: RCCRS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func recordingCfg(maxAge string) *config.Config {
	return &config.Config{ExternalSystems: []config.ExternalSystem{{Name: "stripe", Calls: "x", ProvenBy: []string{"live"}, Version: "2025-09-30.dahlia", MaxAge: maxAge}}}
}

func atDay(t *testing.T, d string) {
	t.Helper()
	old := today
	t.Cleanup(func() { today = old })
	today = func() string { return d }
}

var fixtureNode = mapx.Node{ID: "src/__fixtures__/charge.test.ts", Kind: mapx.KindTest}

func TestRecordingCurrent_theVersion(t *testing.T) {
	t.Run("RCCRS-B01: A recording of another version fails", func(t *testing.T) {})
	atDay(t, "2026-11-10")
	if v, msg := checkRecordingCurrent("// @recorded: stripe 2025-09-30.dahlia 2026-11-01\n", fixtureNode, "", nil, recordingCfg("")); v != Pass {
		t.Errorf("the current version: %v %s", v, msg)
	}
	v, msg := checkRecordingCurrent("/* @recorded: stripe 2024-06-20 2026-11-01 */\n", fixtureNode, "", nil, recordingCfg(""))
	if v != Fail || !strings.Contains(msg, "2024-06-20") || !strings.Contains(msg, "2025-09-30.dahlia") {
		t.Errorf("an older version: %v %s", v, msg)
	}
}

func TestRecordingCurrent_theAge(t *testing.T) {
	t.Run("RCCRS-B02: A recording past its system's age diverges", func(t *testing.T) {})
	atDay(t, "2026-11-10")
	old := "// @recorded: stripe 2025-09-30.dahlia 2026-07-13\n"
	v, msg := checkRecordingCurrent(old, fixtureNode, "", nil, recordingCfg("90d"))
	if v != Diverge || !strings.Contains(msg, "2026-07-13") || !strings.Contains(msg, "120 days") {
		t.Errorf("past its age: %v %s", v, msg)
	}
	if v, msg := checkRecordingCurrent(old, fixtureNode, "", nil, recordingCfg("")); v != Pass {
		t.Errorf("no max age, no aging: %v %s", v, msg)
	}
}

func TestRecordingCurrent_noneAndMalformed(t *testing.T) {
	t.Run("RCCRS-B03: A file with no recording is skipped", func(t *testing.T) {})
	t.Run("RCCRS-E01: A mark that cannot be placed fails", func(t *testing.T) {})
	if v, _ := checkRecordingCurrent("it('x', () => {})\n", fixtureNode, "", nil, recordingCfg("")); v != Skip {
		t.Errorf("no recording: %v", v)
	}
	if v, msg := checkRecordingCurrent("// @recorded: stripe 2025-09-30.dahlia\n", fixtureNode, "", nil, recordingCfg("")); v != Fail || !strings.Contains(msg, "@recorded") {
		t.Errorf("no date: %v %s", v, msg)
	}
	if v, msg := checkRecordingCurrent("// @recorded: paypal v1 2026-11-01\n", fixtureNode, "", nil, recordingCfg("")); v != Fail || !strings.Contains(msg, "paypal") {
		t.Errorf("undeclared system: %v %s", v, msg)
	}
}

package quality

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
)

// The header counts everything the check found, by kind: failures (blocking and
// informative) and divergences. It said "N issue(s) — divergences recorded:" over the
// failures alone, which read as N issues plus the divergences, with the divergences
// counted nowhere.
func TestFindingsSummaryCountsEveryKind(t *testing.T) {
	block := gate.Result{Gate: "g1", Target: "a", Verdict: gate.Fail, Blocking: true}
	info := gate.Result{Gate: "g2", Target: "b", Verdict: gate.Fail}
	drift := gate.Result{Gate: "g3", Target: "c", Verdict: gate.Pending, Detail: "diverged"}
	p := gate.Profile{Results: []gate.Result{block, info, info, drift}, Failures: []gate.Result{block, info, info}}

	out := captureStdout(t, func() { printProfile(p, false, false) })
	for _, want := range []string{"3 failure(s)", "1 blocking", "2 informative", "1 divergence(s)", "--show-drift"} {
		if !strings.Contains(out, want) {
			t.Errorf("the summary should say %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "issue(s) — divergences recorded") {
		t.Errorf("the old ambiguous header is back:\n%s", out)
	}

	// Only divergences: still counted, with no failure list.
	out = captureStdout(t, func() { printProfile(gate.Profile{Results: []gate.Result{drift}}, false, false) })
	if !strings.Contains(out, "0 failure(s)") || !strings.Contains(out, "1 divergence(s)") {
		t.Errorf("divergences alone must still be counted:\n%s", out)
	}
}

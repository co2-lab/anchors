// @anchors
//   code: LDCMT
//   ref: LDCMS

package ops

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/load"
)

func TestLoadCmd(t *testing.T) {
	t.Run("LDCMS-B01: anchors load shows the capacity, what is in use, the runs holding the machine and the runs waiting, and which waited past their limit", func(t *testing.T) {
		now := time.Now()
		st := load.State{
			Leases:  []load.Lease{{Project: "/p/app", Task: "test unit", Level: "heavy", PID: os.Getpid(), Cost: load.Cost{CPU: 4, Memory: 6 << 30}, Since: now.Add(-90 * time.Second)}},
			Waiters: []load.Lease{{Project: "/p/web", Task: "mutation unit", Level: "medium", PID: 7, Cost: load.Cost{CPU: 2, Memory: 2 << 30}, Since: now.Add(-10 * time.Minute), Limit: 300}},
		}
		var b bytes.Buffer
		printLoad(&b, st, load.Capacity{CPU: 5, Memory: 8 << 30}, now)
		out := b.String()
		for _, want := range []string{
			"capacity   5 CPU, 8GB",
			"in use     4 CPU, 6GB",
			"running    app test unit (heavy: 4 CPU, 6GB), for 1m30s",
			"waiting    web mutation unit (medium: 2 CPU, 2GB), for 10m0s, pid 7 — past its limit",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in\n%s", want, out)
			}
		}
		var empty bytes.Buffer
		printLoad(&empty, load.State{}, load.Capacity{CPU: 5}, now)
		if !strings.Contains(empty.String(), "running    —") || !strings.Contains(empty.String(), "waiting    —") {
			t.Errorf("an idle machine: %s", empty.String())
		}
	})
}

// @anchors
//   code: TRGTA
//   ref: RTSRT

package gate

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testsig"
)

// retriedWindow is how far back the gate looks for passes on retry, and retriedUnstable how
// many of them, for one test, make it unstable rather than unlucky.
const (
	retriedWindow   = 7 * 24 * time.Hour
	retriedUnstable = 3
)

// checkTestsRetried lists, for a test file, its tests that passed only on a retry under load
// in the last 7 days (DESIGN-retry-on-load.md): a retry lets a run pass, and this keeps it
// in sight. A test that needed 3 retries is unstable: under load it runs out of time.
func checkTestsRetried(_ string, n mapx.Node, root string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	since := time.Now().Add(-retriedWindow)
	count := map[string]int{}
	last := map[string]testsig.Retried{}
	for _, r := range testsig.LoadRetried(root) {
		if filepath.ToSlash(r.File) != n.ID || r.At.Before(since) {
			continue
		}
		count[r.Name]++
		if r.At.After(last[r.Name].At) {
			last[r.Name] = r
		}
	}
	if len(count) == 0 {
		return Pass, ""
	}
	names := make([]string, 0, len(count))
	for name := range count {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		r := last[name]
		part := fmt.Sprintf("%s — %d time(s), last under load %.1f/%d", name, count[name], r.Load, r.CPUs)
		if count[name] >= retriedUnstable {
			part += " — " + i18n.T("gate.tests_retried.unstable")
		}
		parts = append(parts, part)
	}
	return Fail, i18n.T("gate.tests_retried.fail", strings.Join(parts, "; "))
}

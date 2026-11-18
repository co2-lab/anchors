// @anchors
//   code: RTTRT
//   ref: RTSRT

package gate

import (
	"github.com/co2-lab/anchors/internal/testkit"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testsig"
)

func retriedProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	now := time.Now()
	var rec []testsig.Retried
	add := func(name, file string, ago time.Duration, load float64) {
		rec = append(rec, testsig.Retried{Name: name, File: file, At: now.Add(-ago), Load: load, CPUs: 10})
	}
	add("IPSBI-V06#01", "a.test.ts", time.Hour, 31)
	for i := 1; i <= 3; i++ {
		add("UIRNU-B01#01", "a.test.ts", time.Duration(i)*time.Hour, float64(20+i))
	}
	add("OLD-B01", "a.test.ts", 8*24*time.Hour, 50)
	add("OTHER-B01", "b.test.ts", time.Hour, 50)
	if err := testsig.AddRetried(root, rec); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestTestsRetried_none(t *testing.T) {
	t.Run("RTSRT-B01: A test file with no retry passes", func(t *testing.T) {})
	root := retriedProject(t)
	if v, d := checkTestsRetried("tests-retried", mapx.Node{ID: "c.test.ts", Kind: mapx.KindTest}, root, nil, testkit.WithComments(nil)); v != Pass {
		t.Errorf("no retry of this file: %v %s", v, d)
	}
}

func TestTestsRetried_listed(t *testing.T) {
	t.Run("RTSRT-B02: A test file with tests passed on retry is listed", func(t *testing.T) {})
	root := retriedProject(t)
	v, d := checkTestsRetried("tests-retried", mapx.Node{ID: "a.test.ts", Kind: mapx.KindTest}, root, nil, testkit.WithComments(nil))
	if v != Fail || !strings.Contains(d, "IPSBI-V06#01 — 1 time(s), last under load 31.0/10") ||
		!strings.Contains(d, "UIRNU-B01#01 — 3 time(s), last under load 21.0/10 — unstable") ||
		strings.Contains(d, "OLD-B01") || strings.Contains(d, "OTHER-B01") {
		t.Errorf("the retried tests, with counts, load and unstable: %v %s", v, d)
	}
}

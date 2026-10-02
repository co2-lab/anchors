// @anchors
//   ref: BRCOV

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// branchNode is a code file whose lcov listed 4 branches and missed two, on lines 3 and 7.
func branchNode() mapx.Node {
	return mapx.Node{ID: "a.ts", Kind: mapx.KindCode, Rev: "r1", Signal: &mapx.TestSignal{
		TotalLines: 10, CoveredLines: 10, AtRev: "r1", BranchTotal: 4, BranchMissed: "3:0:1 7:0:0",
	}}
}

func branchCfg(floor *float64) *config.Config {
	return &config.Config{Gates: []config.Gate{{Name: "b", Check: "branch-coverage", On: []string{"code"}, MinPercent: floor}}}
}

func pct(v float64) *float64 { return &v }

func TestBranchCoverage_belowTheFloor(t *testing.T) {
	t.Run("BRCOV-B01: Below the floor fails naming the lines", func(t *testing.T) {})
	v, msg := checkBranchCoverage("", branchNode(), "", nil, branchCfg(nil))
	if v != Fail || !strings.Contains(msg, "50.0% of the branches taken, below the floor of 100%") ||
		!strings.Contains(msg, "2 of 4 never taken, on line(s) 3, 7") || strings.Contains(msg, "dead") {
		t.Errorf("with no floor every branch is asked, got %v: %s", v, msg)
	}
	if v, msg := checkBranchCoverage("", branchNode(), "", nil, branchCfg(pct(50))); v != Pass {
		t.Errorf("at the floor it passes, got %v: %s", v, msg)
	}
	if v, _ := checkBranchCoverage("", branchNode(), "", nil, branchCfg(pct(60))); v != Fail {
		t.Errorf("below the floor it fails, got %v", v)
	}
	two := branchNode()
	two.Signal.BranchMissed = "3:0:1 3:0:2"
	if _, msg := checkBranchCoverage("", two, "", nil, branchCfg(nil)); !strings.Contains(msg, "2 of 4 never taken, on line(s) 3.") {
		t.Errorf("two branches of one line count two and name the line once, got %s", msg)
	}
}

func TestBranchCoverage_waived(t *testing.T) {
	t.Run("BRCOV-B02: A waived branch is left out", func(t *testing.T) {})
	src := "a\n// @no-branch: the fallback only a broken device reaches\nif (x) {\n\n\n\nif (y) { // @no-branch:\n"
	v, msg := checkBranchCoverage(src, branchNode(), "", nil, branchCfg(nil))
	if v != Fail || !strings.Contains(msg, "1 of 4 never taken, on line(s) 7.") {
		t.Errorf("only the unwaived line is missed, got %v: %s", v, msg)
	}
	same := "a\nb\nif (x) { // @no-branch: guarded upstream\n\n\n\nif (y) {\n"
	if _, msg := checkBranchCoverage(same, branchNode(), "", nil, branchCfg(nil)); !strings.Contains(msg, "on line(s) 7.") {
		t.Errorf("a waiver on the branch's own line waives it, got %s", msg)
	}
	lines := noBranchLines("x\n@no-branch: why\n")
	if !lines[2] || !lines[3] || lines[1] || len(lines) != 2 {
		t.Errorf("the waiver covers its line and the next, got %v", lines)
	}
}

func TestBranchCoverage_likelyDead(t *testing.T) {
	t.Run("BRCOV-B03: A branch no test reaches is likely dead", func(t *testing.T) {})
	n := branchNode()
	n.Signal.MutationAtRev = "r1"
	n.Signal.NoCoverageAt = "5,7,9"
	v, msg := checkBranchCoverage("", n, "", nil, branchCfg(pct(0)))
	if v != Fail || !strings.Contains(msg, "likely dead branch on line(s) 7:") || strings.Contains(msg, "below the floor") {
		t.Errorf("a dead branch fails whatever the floor, got %v: %s", v, msg)
	}
	n.Signal.NoCoverageAt = "3,7"
	if _, msg := checkBranchCoverage("", n, "", nil, branchCfg(pct(0))); !strings.Contains(msg, "line(s) 3, 7:") {
		t.Errorf("every dead line is named, in order, got %s", msg)
	}
	n.Signal.MutationAtRev = "r0"
	if v, msg := checkBranchCoverage("", n, "", nil, branchCfg(pct(0))); v != Pass {
		t.Errorf("a mutation of another revision says nothing, got %v: %s", v, msg)
	}
}

func TestBranchCoverage_nothingToMeasure(t *testing.T) {
	t.Run("BRCOV-B04: Nothing to measure is skipped or pending", func(t *testing.T) {})
	if v, _ := checkBranchCoverage("", mapx.Node{Kind: mapx.KindSpec}, "", nil, nil); v != Skip {
		t.Errorf("a spec is skipped, got %v", v)
	}
	if v, _ := checkBranchCoverage("", mapx.Node{Kind: mapx.KindCode}, "", nil, nil); v != Pending {
		t.Errorf("no coverage is pending, got %v", v)
	}
	if v, _ := checkBranchCoverage("", mapx.Node{Kind: mapx.KindCode, Signal: &mapx.TestSignal{}}, "", nil, nil); v != Pending {
		t.Errorf("a signal with no line is pending, got %v", v)
	}
	stale := branchNode()
	stale.Rev = "r2"
	if v, _ := checkBranchCoverage("", stale, "", nil, nil); v != Pending {
		t.Errorf("stale coverage is pending, got %v", v)
	}
	flat := branchNode()
	flat.Signal.BranchTotal, flat.Signal.BranchMissed = 0, ""
	if v, _ := checkBranchCoverage("", flat, "", nil, nil); v != Skip {
		t.Errorf("coverage with no branch is skipped, got %v", v)
	}
}

func TestBranchCoverage_nothingToCover(t *testing.T) {
	t.Run("BRCOV-B05: Branch coverage reads a file with nothing to cover as line coverage does", func(t *testing.T) {})
	listed := mapx.Node{Kind: mapx.KindCode, Rev: "r1", Signal: &mapx.TestSignal{CoverageRev: "r1"}}
	omitted := mapx.Node{Kind: mapx.KindCode, Rev: "r1", Signal: &mapx.TestSignal{CoverageOmitted: "r1"}}
	if v, _ := checkBranchCoverage("", listed, "", nil, nil); v != Skip {
		t.Errorf("listed with no line is skipped, got %v", v)
	}
	if v, _ := checkBranchCoverage("", omitted, "", nil, nil); v != Diverge {
		t.Errorf("omitted is a divergence, got %v", v)
	}
}

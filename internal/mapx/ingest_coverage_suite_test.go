package mapx

import "testing"

func lines(covered []int, uncovered []int) map[int]bool {
	m := map[int]bool{}
	for _, l := range covered {
		m[l] = true
	}
	for _, l := range uncovered {
		m[l] = false
	}
	return m
}

// The integration suite, ingested after the unit one, must not erase the unit coverage of
// the same file: a line covered by ANY suite is covered.
func TestIngestCoverageSuite_unionOfLines(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "infra/db.ts", Kind: KindCode, Rev: "r1"}}}

	// unit covers 1-2 of 1-4; integration covers 3-4 of 1-4.
	g.IngestCoverageSuite(map[string]FileCov{"infra/db.ts": {Covered: 2, Total: 4, Lines: lines([]int{1, 2}, []int{3, 4})}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"infra/db.ts": {Covered: 2, Total: 4, Lines: lines([]int{3, 4}, []int{1, 2})}}, "integration.info", "t2")

	sig := g.Nodes[0].Signal
	if sig.CoveredLines != 4 || sig.TotalLines != 4 || sig.LineCoverage != 100 {
		t.Errorf("union should be 4/4 = 100%%, got %d/%d = %.0f%%", sig.CoveredLines, sig.TotalLines, sig.LineCoverage)
	}
}

// A line covered by BOTH suites is counted once.
func TestIngestCoverageSuite_overlapCountedOnce(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 2, Total: 4, Lines: lines([]int{1, 2}, []int{3, 4})}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 2, Total: 4, Lines: lines([]int{1, 2}, []int{3, 4})}}, "integration.info", "t2")

	if sig := g.Nodes[0].Signal; sig.CoveredLines != 2 || sig.TotalLines != 4 {
		t.Errorf("the same two lines from two suites must count once: got %d/%d", sig.CoveredLines, sig.TotalLines)
	}
}

// Re-ingesting the SAME suite replaces it — it still erases what that suite stopped
// covering.
func TestIngestCoverageSuite_sameSuiteReplaces(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 4, Total: 4, Lines: lines([]int{1, 2, 3, 4}, nil)}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 4, Lines: lines([]int{1}, []int{2, 3, 4})}}, "unit.info", "t2")

	if sig := g.Nodes[0].Signal; sig.CoveredLines != 1 {
		t.Errorf("the unit suite now covers 1 line and the node says %d", sig.CoveredLines)
	}
}

// The coverage one suite measured on an OLDER rev keeps the node stale — the lesson of
// ProvenBySuite, applied here from the start.
func TestIngestCoverageSuite_staleSuiteKeepsNodeStale(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 2, Lines: lines([]int{1}, []int{2})}}, "unit.info", "t1")
	g.Nodes[0].Rev = "r2"
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 2, Lines: lines([]int{2}, []int{1})}}, "integration.info", "t2")

	if !g.Nodes[0].SignalStale() {
		t.Error("the unit coverage was measured at r1 and the node was reported fresh at r2")
	}
}

// A report with only totals (LF/LH, no DA lines) cannot join a line union. The fallback is
// the suite with the most covered lines — an under-estimate, never a double count.
func TestIngestCoverageSuite_totalsOnlyFallback(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 3, Total: 10}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 5, Total: 10}}, "integration.info", "t2")

	if sig := g.Nodes[0].Signal; sig.CoveredLines != 5 || sig.TotalLines != 10 {
		t.Errorf("without line detail the best suite should be used (5/10), got %d/%d", sig.CoveredLines, sig.TotalLines)
	}
}

func TestRanges_roundTrip(t *testing.T) {
	in := []int{12, 1, 2, 3, 9, 13, 5, 4, 20}
	enc := encodeRanges(append([]int(nil), in...))
	if enc != "1-5,9,12-13,20" {
		t.Errorf("encodeRanges = %q", enc)
	}
	if got := decodeRanges(enc); len(got) != len(in) {
		t.Errorf("decodeRanges(%q) = %v", enc, got)
	}
}

// A suite recorded before the file was edited speaks of another numbering: joining its
// lines to the fresh suite's mixes two texts. The reference case: unit fresh at 73/75,
// integration stale at 0/76 — the union read 73/106 = 69% on a file covered at 97%.
func TestIngestCoverageSuite_staleSuiteLinesDoNotJoinTheUnion(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "handler.ts", Kind: KindCode, Rev: "old"}}}
	var oldLines []int
	for l := 1; l <= 76; l++ {
		oldLines = append(oldLines, l+30) // the old numbering, shifted
	}
	g.IngestCoverageSuite(map[string]FileCov{"handler.ts": {Covered: 0, Total: 76, Lines: lines(nil, oldLines)}}, "integration.info", "t1")

	g.Nodes[0].Rev = "new" // the file was edited
	var hit, miss []int
	for l := 1; l <= 73; l++ {
		hit = append(hit, l)
	}
	miss = []int{74, 75}
	g.IngestCoverageSuite(map[string]FileCov{"handler.ts": {Covered: 73, Total: 75, Lines: lines(hit, miss)}}, "unit.info", "t2")

	sig := g.Nodes[0].Signal
	if sig.CoveredLines != 73 || sig.TotalLines != 75 {
		t.Errorf("only the fresh suite speaks of the current lines: want 73/75, got %d/%d (%.0f%%)",
			sig.CoveredLines, sig.TotalLines, sig.LineCoverage)
	}
	if !g.Nodes[0].SignalStale() {
		t.Error("the integration suite has not re-measured the new rev, and the node reads fresh")
	}
}

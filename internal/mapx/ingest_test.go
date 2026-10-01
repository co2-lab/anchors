package mapx

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func ingestGraph() *Graph {
	return &Graph{Nodes: []Node{
		{ID: "src/A.spec.md", Kind: KindSpec, Rev: "r1"},
		{ID: "src/A.tsx", Kind: KindCode, Rev: "r1"},
		{ID: "src/A.test.tsx", Kind: KindTest, Rev: "r1"},
	}}
}

// ── Execution ───────────────────────────────────────────────────────────────────────────────

func TestIngestExecutionCrossesScenarios(t *testing.T) {
	t.Run("SGINA-B03: A full run records the proven rules and erases the lost ones", func(t *testing.T) {})
	g := ingestGraph()
	byFile := map[string]ExecByFile{"src/A.test.tsx": {Passed: 2, Failed: 1}}
	proven := map[string]bool{"AAAAX-V01": true, "AAAAX-V02": true} // V03 failed → not proven
	declared := map[string][]string{"src/A.spec.md": {"AAAAX-V01", "AAAAX-V02", "AAAAX-V03"}}
	mf, mc := g.IngestExecution(byFile, proven, declared, "unit", "now")
	if mf != 1 {
		t.Fatalf("expected 1 matched test, got %d", mf)
	}
	if mc != 2 {
		t.Fatalf("expected 2 proven scenarios, got %d", mc)
	}
	for _, n := range g.Nodes {
		if n.ID == "src/A.test.tsx" && (n.Signal == nil || n.Signal.Failed != 1) {
			t.Error("the test node should have Failed=1")
		}
		if n.ID == "src/A.spec.md" {
			if n.Signal == nil || len(n.Signal.ProvenCodes) != 2 {
				t.Errorf("the spec should have 2 proven scenarios, got %+v", n.Signal)
			}
		}
	}
}

func TestIngestExecutionAccumulatesLayers(t *testing.T) {
	t.Run("SGINA-B01: A test node sums its layers and records its revisions", func(t *testing.T) {})
	t.Run("SGINA-B02: Re-ingesting a layer replaces only that layer", func(t *testing.T) {})
	g := ingestGraph()
	g.Nodes = append(g.Nodes, Node{ID: "src/util.ts", Kind: KindCode, Rev: "u1"})
	g.Edges = []Edge{{From: "src/A.test.tsx", To: "src/util.ts", Type: EdgeDependsOn}}
	test := func() *TestSignal { return g.Nodes[2].Signal }

	g.IngestExecution(map[string]ExecByFile{"src/A.test.tsx": {Passed: 2}}, nil, nil, "", "t1")
	g.IngestExecution(map[string]ExecByFile{"src/A.test.tsx": {Failed: 1}}, nil, nil, "e2e", "t2")
	s := test()
	if s.ByLayer["unit"].Passed != 2 || s.ByLayer["e2e"].Failed != 1 {
		t.Fatalf("no layer named is `unit`, and each layer is kept: %+v", s.ByLayer)
	}
	if s.Passed != 2 || s.Failed != 1 || s.AtRev != "r1" || s.IngestedAt != "t2" || s.ClosureRev["src/util.ts"] != "u1" {
		t.Errorf("totals, rev, closure and date should be recorded: %+v", s)
	}

	g.IngestExecution(map[string]ExecByFile{"src/A.test.tsx": {Passed: 5}}, nil, nil, "unit", "t3")
	s = test()
	if s.ByLayer["unit"].Passed != 5 || s.ByLayer["e2e"].Failed != 1 || s.Passed != 5 || s.Failed != 1 {
		t.Errorf("re-ingesting unit replaces unit only: %+v / totals %d passed %d failed", s.ByLayer, s.Passed, s.Failed)
	}
}

// ── Proven rules ────────────────────────────────────────────────────────────────────────────

// An ingestion that only ADDS lets the map claim proof that stopped existing.
//
// MEASURED in the reference app: after fixing the filter that made a spec declare its neighbour's rule,
// seven nodes kept carrying someone else's `proven_codes`. `map build` and `ingest` ran again and
// cleaned nothing — because `len(pc) > 0` meant "nothing to write", when the right answer is "no
// proof any more", which is information, not the absence of it.
func TestIngestErasesProofThatStoppedExisting(t *testing.T) {
	t.Run("SGINA-B03: A full run records the proven rules and erases the lost ones", func(t *testing.T) {})
	g := ingestGraph()

	// 1st ingestion: the spec has two proven scenarios.
	g.IngestExecution(
		map[string]ExecByFile{"src/A.test.tsx": {Passed: 2}},
		map[string]bool{"AAAAX-V01": true, "AAAAX-V02": true},
		map[string][]string{"src/A.spec.md": {"AAAAX-V01", "AAAAX-V02"}},
		"unit", "before")

	spec := func() *Node { return &g.Nodes[0] }
	if n := spec(); n.Signal == nil || len(n.Signal.ProvenCodes) != 2 {
		t.Fatalf("setup failed: expected 2 proven, got %+v", n)
	}

	// 2nd ingestion: no scenario passes any more.
	g.IngestExecution(
		map[string]ExecByFile{"src/A.test.tsx": {Passed: 0, Failed: 2}},
		map[string]bool{},
		map[string][]string{"src/A.spec.md": {"AAAAX-V01", "AAAAX-V02"}},
		"unit", "after")

	if n := spec(); n.Signal != nil && len(n.Signal.ProvenCodes) > 0 {
		t.Errorf("the map still claims %v proven after the proof vanished — `stale` will not "+
			"charge the unit, because the old proof answers for it", n.Signal.ProvenCodes)
	}
}

// Monorepo: each suite is ingested alone. Mobile's must not speak for the backend's —
// MEASURED in the reference app: the mobile ingestion wrote EMPTY on 131 backend specs.
func TestIngestPerSuiteDoesNotEraseAnotherSuitesProof(t *testing.T) {
	t.Run("SGINA-B04: One suite never erases another suite's proof", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{
		{ID: "back/A.spec.md", Kind: KindSpec, Rev: "r1"},
		{ID: "mob/B.spec.md", Kind: KindSpec, Rev: "r1"},
	}}
	declared := map[string][]string{
		"back/A.spec.md": {"AAAA-B01"},
		"mob/B.spec.md":  {"BBBB-B01"},
	}
	proven := func(id string) []string {
		for _, n := range g.Nodes {
			if n.ID == id && n.Signal != nil {
				return n.Signal.ProvenCodes
			}
		}
		return nil
	}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAA-B01": true}, nil, declared, "unit", "back/junit.xml", "t1")
	g.IngestExecutionSuite(nil, map[string]bool{"BBBB-B01": true}, nil, declared, "unit", "mob/junit.xml", "t2")

	if got := proven("back/A.spec.md"); len(got) != 1 || got[0] != "AAAA-B01" {
		t.Errorf("the mobile suite erased the backend's proof: %v", got)
	}
	if got := proven("mob/B.spec.md"); len(got) != 1 || got[0] != "BBBB-B01" {
		t.Errorf("mobile's proof not recorded: %v", got)
	}

	// re-ingesting the SAME suite without the proof still erases — the rule of the test above holds
	g.IngestExecutionSuite(nil, map[string]bool{}, nil, declared, "unit", "back/junit.xml", "t3")
	if got := proven("back/A.spec.md"); len(got) != 0 {
		t.Errorf("the proof the suite itself stopped giving stayed in the map: %v", got)
	}
}

func TestIngestSuite_provenCodesAreTheSortedUnion(t *testing.T) {
	t.Run("SGINA-I01: The proven rules are the sorted union of the suites", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "r1"}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B03", "AAAAA-B02", "AAAAA-B01"}}
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true, "AAAAA-B01": true}, nil, decl, "", "one.xml", "t1")
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B03": true, "AAAAA-B01": true}, nil, decl, "", "two.xml", "t2")
	if got := g.Nodes[0].Signal.ProvenCodes; !slices.Equal(got, []string{"AAAAA-B01", "AAAAA-B02", "AAAAA-B03"}) {
		t.Errorf("the proven codes should be the sorted union with no repetition, got %v", got)
	}
}

// A proof from one suite must not be laundered as fresh by ANOTHER suite running.
//
// Suite A proves B01 at rev1, the spec changes, suite B proves B02 at rev2. With one `AtRev` for
// the whole node, the map said the union was measured at rev2 — B01 included, which was never
// measured there. Probed before the fix: `stale=false`.
func TestIngestSuite_staleSuiteKeepsTheNodeStale(t *testing.T) {
	t.Run("SGINA-B06: The union is as fresh as its oldest contributor", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev1"}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01", "AAAAA-B02"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, nil, decl, "", "mobile.xml", "t1")
	g.Nodes[0].Rev = "rev2"
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true}, nil, decl, "", "backend.xml", "t2")

	if !g.Nodes[0].SignalStale() {
		t.Fatalf("suite mobile.xml measured rev1 and did not run again: the node must be stale (at_rev=%s)",
			g.Nodes[0].Signal.AtRev)
	}

	// Once the stale suite runs again at the current rev, the union is fresh.
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, nil, decl, "", "mobile.xml", "t3")
	if g.Nodes[0].SignalStale() {
		t.Errorf("every contributing suite measured rev2 and the node is still stale (at_rev=%s)",
			g.Nodes[0].Signal.AtRev)
	}
}

// A suite that stops proving anything leaves the union and leaves no empty entry behind.
func TestIngestSuite_suiteThatProvesNothingLeaves(t *testing.T) {
	t.Run("SGINA-B05: A suite that proves nothing leaves the union", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev1"}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, nil, decl, "", "mobile.xml", "t1")
	g.IngestExecutionSuite(nil, map[string]bool{}, nil, decl, "", "mobile.xml", "t2")

	sig := g.Nodes[0].Signal
	if len(sig.ProvenCodes) != 0 {
		t.Errorf("the only suite stopped proving B01 and the union still has %v", sig.ProvenCodes)
	}
	if _, ok := sig.ProvenBySuite["mobile.xml"]; ok {
		t.Error("an empty suite entry was left in the versioned map")
	}
}

// An entry written before per-suite revs existed has unknown freshness — and unknown is not
// fresh.
func TestIngestSuite_legacyEntryWithoutRevIsStale(t *testing.T) {
	t.Run("SGINA-B06: The union is as fresh as its oldest contributor", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev2",
		Signal: &TestSignal{ProvenBySuite: map[string][]string{"old.xml": {"AAAAA-B01"}}}}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01", "AAAAA-B02"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true}, nil, decl, "", "new.xml", "t1")
	if !g.Nodes[0].SignalStale() {
		t.Error("a suite with no recorded rev was treated as fresh")
	}
}

// A PARTIAL run (`anchors test --changed`) measured only its cut. Reported from the reference app: a run of 4
// test files erased the proof of MoneyDetailScreen, whose test it never executed.
func TestIngestPartialRunKeepsWhatItDidNotSee(t *testing.T) {
	t.Run("SGINA-B07: A partial run changes only what it saw", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{
		{ID: "money.spec.md", Kind: KindSpec, Rev: "r1"},
		{ID: "member.spec.md", Kind: KindSpec, Rev: "r1"},
	}}
	declared := map[string][]string{
		"money.spec.md":  {"MNDT-B01", "MNDT-B02"},
		"member.spec.md": {"MBDT-B01", "MBDT-B02"},
	}
	codes := func(id string) []string {
		for _, n := range g.Nodes {
			if n.ID == id && n.Signal != nil {
				return n.Signal.ProvenCodes
			}
		}
		return nil
	}
	// Full run: everything proven.
	all := map[string]bool{"MNDT-B01": true, "MNDT-B02": true, "MBDT-B01": true, "MBDT-B02": true}
	g.IngestExecutionSuite(nil, all, nil, declared, "unit", "mobile", "t0")

	// Partial run: only member's test ran; MBDT-B01 passed, MBDT-B02 failed.
	g.IngestExecutionSuite(nil, map[string]bool{"MBDT-B01": true},
		map[string]bool{"MBDT-B01": true, "MBDT-B02": true}, declared, "unit", "mobile", "t1")
	if got := codes("money.spec.md"); len(got) != 2 {
		t.Errorf("a spec outside the cut must keep its proof, has %v", got)
	}
	if got := codes("member.spec.md"); len(got) != 1 || got[0] != "MBDT-B01" {
		t.Errorf("a code the run saw fail must lose its proof, the one it saw pass keeps it: %v", got)
	}

	// Partial run that saw MNDT-B01 only: MNDT-B02, not seen, keeps its proof.
	g.IngestExecutionSuite(nil, map[string]bool{"MNDT-B01": true},
		map[string]bool{"MNDT-B01": true}, declared, "unit", "mobile", "t2")
	if got := codes("money.spec.md"); len(got) != 2 {
		t.Errorf("a code the partial run did not see keeps its proof: %v", got)
	}

	// Control: a FULL run is the whole measurement — what it did not prove is gone.
	g.IngestExecutionSuite(nil, map[string]bool{"MBDT-B01": true}, nil, declared, "unit", "mobile", "t3")
	if got := codes("money.spec.md"); len(got) != 0 {
		t.Errorf("a full run that proved nothing of money erases its proof, has %v", got)
	}
}

// A report ingested from outside the repository is an ad-hoc run that never runs again. Its proof
// must not keep the node stale forever: a full run of an in-repo suite drops it, and the in-repo
// proof stands alone and fresh.
func TestDropExternalSuites(t *testing.T) {
	t.Run("SGINA-B08: External suites are dropped and the union recomputed", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{
		{ID: "a.spec.md", Kind: KindSpec, Rev: "rev1"},
		{ID: "a.go", Kind: KindCode, Rev: "c1"},
	}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01", "AAAAA-B02"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true}, nil, decl, "", ExternalSuitePrefix+"report.xml", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.go": {Covered: 1, Total: 2, Lines: map[int]bool{1: true, 2: false}}}, ExternalSuitePrefix+"lcov.info", "t1")
	g.Nodes[0].Rev = "rev2"
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, nil, decl, "", ".anchors/junit.xml", "t2")
	if !g.Nodes[0].SignalStale() {
		t.Fatal("precondition: the external suite measured rev1 and keeps the node stale")
	}

	dropped := g.DropExternalSuites()
	if len(dropped) != 2 || dropped[0] != "external/lcov.info" || dropped[1] != "external/report.xml" {
		t.Fatalf("both external suites are dropped and named: %v", dropped)
	}
	spec := g.Nodes[0].Signal
	if g.Nodes[0].SignalStale() {
		t.Errorf("with the external suite gone, the in-repo suite is fresh (at_rev=%s)", spec.AtRev)
	}
	if len(spec.ProvenCodes) != 1 || spec.ProvenCodes[0] != "AAAAA-B01" {
		t.Errorf("only the in-repo proof remains: %v", spec.ProvenCodes)
	}
	if _, ok := spec.ProvenBySuite[".anchors/junit.xml"]; !ok {
		t.Errorf("the in-repo suite must stay: %v", spec.ProvenBySuite)
	}
	code := g.Nodes[1].Signal
	if len(code.CoverageBySuite) != 0 || code.TotalLines != 0 || code.LineCoverage != 0 {
		t.Errorf("the external coverage leaves no number behind: %+v", code)
	}
	if again := g.DropExternalSuites(); len(again) != 0 {
		t.Errorf("nothing left to drop, dropped %v", again)
	}
}

// ── Coverage ────────────────────────────────────────────────────────────────────────────────

func TestIngestCoverageAndStale(t *testing.T) {
	t.Run("SGINA-B09: Coverage with no suite records the percentage and the baseline", func(t *testing.T) {})
	t.Run("SGINA-B16: A signal goes stale when its file moves", func(t *testing.T) {})
	g := ingestGraph()
	g.IngestCoverage(map[string]FileCov{"src/A.tsx": {Covered: 3, Total: 5}}, "now")
	code := &g.Nodes[1]
	if code.Signal == nil || code.Signal.LineCoverage != 60 || code.Signal.AtRev != "r1" {
		t.Fatalf("expected 60%% coverage at r1, got %+v", code.Signal)
	}
	if code.Signal.PrevLineCoverage != 0 {
		t.Errorf("the first measurement has no baseline, got %v", code.Signal.PrevLineCoverage)
	}
	g.IngestCoverage(map[string]FileCov{"src/A.tsx": {Covered: 4, Total: 5}}, "later")
	if code.Signal.LineCoverage != 80 || code.Signal.PrevLineCoverage != 60 {
		t.Errorf("the earlier measurement becomes the baseline: got %v with baseline %v", code.Signal.LineCoverage, code.Signal.PrevLineCoverage)
	}
	if code.SignalStale() {
		t.Error("freshly ingested, it should not be stale")
	}
	code.Rev = "r2" // the file changed
	if !code.SignalStale() {
		t.Error("after the rev changed, the signal should be stale")
	}
	if (Node{Rev: "r2", Signal: &TestSignal{Passed: 1}}).SignalStale() {
		t.Error("a signal with no recorded rev is not stale")
	}
}

func coverageLines(covered []int, uncovered []int) map[int]bool {
	m := map[int]bool{}
	for _, l := range covered {
		m[l] = true
	}
	for _, l := range uncovered {
		m[l] = false
	}
	return m
}

// The integration suite, ingested after the unit one, must not erase the unit coverage of the
// same file: a line covered by ANY suite is covered.
func TestIngestCoverageSuite_unionOfLines(t *testing.T) {
	t.Run("SGINA-B10: Suite coverage is the union of lines", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "infra/db.ts", Kind: KindCode, Rev: "r1"}}}

	// unit covers 1-2 of 1-4; integration covers 3-4 of 1-4.
	g.IngestCoverageSuite(map[string]FileCov{"infra/db.ts": {Covered: 2, Total: 4, Lines: coverageLines([]int{1, 2}, []int{3, 4})}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"infra/db.ts": {Covered: 2, Total: 4, Lines: coverageLines([]int{3, 4}, []int{1, 2})}}, "integration.info", "t2")

	sig := g.Nodes[0].Signal
	if sig.CoveredLines != 4 || sig.TotalLines != 4 || sig.LineCoverage != 100 {
		t.Errorf("union should be 4/4 = 100%%, got %d/%d = %.0f%%", sig.CoveredLines, sig.TotalLines, sig.LineCoverage)
	}
}

// A line covered by BOTH suites is counted once.
func TestIngestCoverageSuite_overlapCountedOnce(t *testing.T) {
	t.Run("SGINA-B10: Suite coverage is the union of lines", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 2, Total: 4, Lines: coverageLines([]int{1, 2}, []int{3, 4})}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 2, Total: 4, Lines: coverageLines([]int{1, 2}, []int{3, 4})}}, "integration.info", "t2")

	if sig := g.Nodes[0].Signal; sig.CoveredLines != 2 || sig.TotalLines != 4 {
		t.Errorf("the same two lines from two suites must count once: got %d/%d", sig.CoveredLines, sig.TotalLines)
	}
}

// Re-ingesting the SAME suite replaces it — it still erases what that suite stopped covering.
func TestIngestCoverageSuite_sameSuiteReplaces(t *testing.T) {
	t.Run("SGINA-B10: Suite coverage is the union of lines", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 4, Total: 4, Lines: coverageLines([]int{1, 2, 3, 4}, nil)}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 4, Lines: coverageLines([]int{1}, []int{2, 3, 4})}}, "unit.info", "t2")

	if sig := g.Nodes[0].Signal; sig.CoveredLines != 1 {
		t.Errorf("the unit suite now covers 1 line and the node says %d", sig.CoveredLines)
	}
}

// The coverage one suite measured on an OLDER rev keeps the node stale — the lesson of
// ProvenBySuite, applied here from the start.
func TestIngestCoverageSuite_staleSuiteKeepsNodeStale(t *testing.T) {
	t.Run("SGINA-B12: A suite measured before the edit stays out of the union", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 2, Lines: coverageLines([]int{1}, []int{2})}}, "unit.info", "t1")
	g.Nodes[0].Rev = "r2"
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 2, Lines: coverageLines([]int{2}, []int{1})}}, "integration.info", "t2")

	if !g.Nodes[0].SignalStale() {
		t.Error("the unit coverage was measured at r1 and the node was reported fresh at r2")
	}
}

// A report with only totals (LF/LH, no DA lines) cannot join a line union. The fallback is the
// suite with the most covered lines — an under-estimate, never a double count.
func TestIngestCoverageSuite_totalsOnlyFallback(t *testing.T) {
	t.Run("SGINA-B11: Totals-only suites fall back to the best suite", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 3, Total: 10}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 5, Total: 10}}, "integration.info", "t2")

	if sig := g.Nodes[0].Signal; sig.CoveredLines != 5 || sig.TotalLines != 10 {
		t.Errorf("without line detail the best suite should be used (5/10), got %d/%d", sig.CoveredLines, sig.TotalLines)
	}
}

func TestRanges_roundTrip(t *testing.T) {
	t.Run("SGINA-B14: Line ranges round-trip", func(t *testing.T) {})
	in := []int{12, 1, 2, 3, 9, 13, 5, 4, 20}
	enc := encodeRanges(append([]int(nil), in...))
	if enc != "1-5,9,12-13,20" {
		t.Errorf("encodeRanges = %q", enc)
	}
	got := decodeRanges(enc)
	want := slices.Clone(in)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("decodeRanges(%q) = %v, want %v", enc, got, want)
	}
}

// A suite recorded before the file was edited speaks of another numbering: joining its lines to
// the fresh suite's mixes two texts. The reference case: unit fresh at 73/75, integration stale at
// 0/76 — the union read 73/106 = 69% on a file covered at 97%.
func TestIngestCoverageSuite_staleSuiteLinesDoNotJoinTheUnion(t *testing.T) {
	t.Run("SGINA-B12: A suite measured before the edit stays out of the union", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "handler.ts", Kind: KindCode, Rev: "old"}}}
	var oldLines []int
	for l := 1; l <= 76; l++ {
		oldLines = append(oldLines, l+30) // the old numbering, shifted
	}
	g.IngestCoverageSuite(map[string]FileCov{"handler.ts": {Covered: 0, Total: 76, Lines: coverageLines(nil, oldLines)}}, "integration.info", "t1")

	g.Nodes[0].Rev = "new" // the file was edited
	var hit, miss []int
	for l := 1; l <= 73; l++ {
		hit = append(hit, l)
	}
	miss = []int{74, 75}
	g.IngestCoverageSuite(map[string]FileCov{"handler.ts": {Covered: 73, Total: 75, Lines: coverageLines(hit, miss)}}, "unit.info", "t2")

	sig := g.Nodes[0].Signal
	if sig.CoveredLines != 73 || sig.TotalLines != 75 {
		t.Errorf("only the fresh suite speaks of the current lines: want 73/75, got %d/%d (%.0f%%)",
			sig.CoveredLines, sig.TotalLines, sig.LineCoverage)
	}
	if !g.Nodes[0].SignalStale() {
		t.Error("the integration suite has not re-measured the new rev, and the node reads fresh")
	}
}

// A report written before the file's current content is ingested with no rev: it neither joins
// the union nor passes for fresh, even when ingested today.
func TestIngestCoverageSuite_reportThatPredatesTheFileIsNotFresh(t *testing.T) {
	t.Run("SGINA-B13: A report older than the file is kept without a revision", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "handler.ts", Kind: KindCode, Rev: "new"}}}
	g.IngestCoverageSuite(map[string]FileCov{"handler.ts": {Covered: 0, Total: 4, Lines: coverageLines(nil, []int{31, 32, 33, 34}), Predates: true}}, "integration.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"handler.ts": {Covered: 3, Total: 4, Lines: coverageLines([]int{1, 2, 3}, []int{4})}}, "unit.info", "t2")

	sig := g.Nodes[0].Signal
	if got := sig.CoverageBySuite["integration.info"].AtRev; got != "" {
		t.Errorf("a report older than the file was stamped with rev %q", got)
	}
	if sig.CoveredLines != 3 || sig.TotalLines != 4 {
		t.Errorf("the predating report joined the union: %d/%d", sig.CoveredLines, sig.TotalLines)
	}
	if !g.Nodes[0].SignalStale() {
		t.Error("with a suite that has not measured this text, the node reads fresh")
	}
}

// ── Mutation ────────────────────────────────────────────────────────────────────────────────

// Without a scope, the second ingestion of the same file erased the first, and the most useful
// number — the DIFFERENCE between isolated and full — could not be computed.
func TestIngestMutationScoped(t *testing.T) {
	t.Run("SGINA-B15: Mutation totals and the per-scope measurement", func(t *testing.T) {})
	g := ingestGraph()
	mu := FileMutation{Killed: 8, Survived: 2, NoCoverage: 1, Ignored: 3, Score: 80}
	if m := g.IngestMutationScoped(map[string]FileMutation{"src/A.tsx": mu}, "isolated", "t1", 60, 80); m != 1 {
		t.Fatalf("expected 1 matched code node, got %d", m)
	}
	s := g.Nodes[1].Signal
	if s.MutantsKilled != 8 || s.MutantsSurvived != 2 || s.MutantsNoCoverage != 1 || s.MutantsIgnored != 3 ||
		s.MutationScore != 80 || s.MutationLow != 60 || s.MutationHigh != 80 || s.AtRev != "r1" {
		t.Errorf("the totals and thresholds should be recorded: %+v", s)
	}
	want := MutationScope{Killed: 8, Survived: 2, NoCoverage: 1, Ignored: 3, Score: 80, AtRev: "r1"}
	if got := s.MutationByScope["isolated"]; got != want {
		t.Errorf("the isolated scope should hold the same numbers at r1, got %+v", got)
	}
	// with no scope, only the totals change
	g.IngestMutation(map[string]FileMutation{"src/A.tsx": {Killed: 1}}, "t2")
	if s.MutantsKilled != 1 || len(s.MutationByScope) != 1 {
		t.Errorf("no scope named writes only the totals: %+v", s)
	}
}

// ── Matching ────────────────────────────────────────────────────────────────────────────────

func TestPathMatches(t *testing.T) {
	t.Run("SGINA-B17: Paths match at a path boundary in either direction", func(t *testing.T) {})
	if !pathMatches("apps/mobile/src/A.tsx", "/abs/apps/mobile/src/A.tsx") {
		t.Error("a longer report path should match by path suffix")
	}
	if !pathMatches("apps/mobile/src/A.tsx", "src/A.tsx") {
		t.Error("a shorter report path should match the node that ends with it")
	}
	if pathMatches("a/b.go", "xa/b.go") {
		t.Error("it should not match without a path boundary")
	}
}

// Monorepo: the runner writes paths relative to the workspace, and the same file exists in two
// workspaces. Measured in the reference app: the landing's coverage landed on mobile's SectionLabel.
func monorepoGraph() *Graph {
	return &Graph{Nodes: []Node{
		{ID: "apps/landing-page/src/atoms/Label.tsx", Kind: KindCode, Rev: "r1"},
		{ID: "apps/mobile/src/atoms/Label.tsx", Kind: KindCode, Rev: "r1"},
		{ID: "apps/mobile/src/atoms/Only.tsx", Kind: KindCode, Rev: "r1"},
	}}
}

func TestResolveReportPathsTieBrokenByTheReportFolder(t *testing.T) {
	t.Run("SGINA-B18: Report paths resolve by the report's folder, or stay unowned", func(t *testing.T) {})
	g := monorepoGraph()
	got, amb := g.ResolveReportPaths(KindCode,
		[]string{"src/atoms/Label.tsx", "src/atoms/Only.tsx", "src/nothing.ts"},
		"apps/landing-page/test-output/coverage/lcov.info")
	if len(amb) != 0 {
		t.Fatalf("the report folder decides; nothing should stay ambiguous: %v", amb)
	}
	if got["src/atoms/Label.tsx"] != "apps/landing-page/src/atoms/Label.tsx" {
		t.Errorf("Label should go to the landing, went to %q", got["src/atoms/Label.tsx"])
	}
	if got["src/atoms/Only.tsx"] != "apps/mobile/src/atoms/Only.tsx" {
		t.Errorf("a unique path should resolve to its only owner, got %q", got["src/atoms/Only.tsx"])
	}
	if got["src/nothing.ts"] != "src/nothing.ts" {
		t.Errorf("an ownerless path is kept as it came, got %q", got["src/nothing.ts"])
	}

	// resolved to the exact ID, ingestion ties ONE node, not two
	cov := map[string]FileCov{got["src/atoms/Label.tsx"]: {Covered: 1, Total: 2}}
	if m := g.IngestCoverage(cov, "now"); m != 1 {
		t.Fatalf("expected 1 matched node, got %d", m)
	}
	for _, n := range g.Nodes {
		if n.ID == "apps/mobile/src/atoms/Label.tsx" && n.Signal != nil {
			t.Error("mobile's namesake must not receive the landing's coverage")
		}
	}
}

func TestResolveReportPathsTieLeavesNoOwner(t *testing.T) {
	t.Run("SGINA-B18: Report paths resolve by the report's folder, or stay unowned", func(t *testing.T) {})
	g := monorepoGraph()
	// report at the root: neither folder is closer
	got, amb := g.ResolveReportPaths(KindCode, []string{"src/atoms/Label.tsx"}, "coverage/lcov.info")
	if len(amb) != 1 || amb[0] != "src/atoms/Label.tsx" {
		t.Fatalf("a tie should come back as ambiguous, got %v", amb)
	}
	if _, ok := got["src/atoms/Label.tsx"]; ok {
		t.Error("an ambiguous path must not be assigned on a guess")
	}
}

func TestIngestLandsOnlyOnItsKindOfNode(t *testing.T) {
	t.Run("SGINA-X01: Each measurement lands only on its kind of node", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{
		{ID: "specs/A.ts", Kind: KindSpec, Rev: "r"},
		{ID: "code/A.ts", Kind: KindCode, Rev: "r"},
		{ID: "tests/A.ts", Kind: KindTest, Rev: "r"},
	}}
	g.IngestExecution(map[string]ExecByFile{"A.ts": {Passed: 1}}, nil, nil, "", "t")
	g.IngestCoverage(map[string]FileCov{"A.ts": {Covered: 1, Total: 1}}, "t")
	g.IngestMutation(map[string]FileMutation{"A.ts": {Killed: 1}}, "t")
	if s := g.Nodes[0].Signal; s != nil {
		t.Errorf("the spec takes no execution, coverage or mutation: %+v", s)
	}
	if s := g.Nodes[1].Signal; s == nil || s.Passed != 0 || s.TotalLines != 1 || s.MutantsKilled != 1 {
		t.Errorf("the code file takes coverage and mutation only: %+v", s)
	}
	if s := g.Nodes[2].Signal; s == nil || s.Passed != 1 || s.TotalLines != 0 || s.MutantsKilled != 0 {
		t.Errorf("the test file takes the execution only: %+v", s)
	}
}

// A report that instruments no line of the file is a measurement too: nothing to cover.
// The percentage was only written when the total was positive, so the node kept the
// previous ingestion's percentage next to a total of zero.
func TestIngestCoverage_zeroTotalClearsThePercentage(t *testing.T) {
	t.Run("SGINA-B19: A report with no instrumented line leaves no percentage behind", func(t *testing.T) {})
	for _, suite := range []string{"", "unit"} {
		g := ingestGraph()
		g.IngestCoverageSuite(map[string]FileCov{"src/A.tsx": {Covered: 4, Total: 5, Lines: coverageLines([]int{1, 2, 3, 4}, []int{5})}}, suite, "t1")
		g.IngestCoverageSuite(map[string]FileCov{"src/A.tsx": {Covered: 0, Total: 0}}, suite, "t2")
		sig := g.Nodes[1].Signal
		if sig.TotalLines != 0 || sig.LineCoverage != 0 {
			t.Errorf("suite %q: after a report of zero lines the node reads %d lines at %.0f%%, want 0 lines at 0%%",
				suite, sig.TotalLines, sig.LineCoverage)
		}
	}
}

// A test file whose path several report entries match (the report's own path and a
// shorter suffix of it) used to take whichever the map iteration reached first, so the
// same report recorded different counts from run to run.
func TestIngest_severalMatchingReportPathsChooseTheSameOne(t *testing.T) {
	t.Run("SGINA-B20: Among several report paths matching a node, the exact one, then the closest in length, is always chosen", func(t *testing.T) {})
	exec := map[string]ExecByFile{
		"A.test.tsx":           {Passed: 1},
		"src/A.test.tsx":       {Passed: 2},
		"app/src/A.test.tsx":   {Passed: 3},
		"x/app/src/A.test.tsx": {Passed: 4},
	}
	cov := map[string]FileCov{
		"A.tsx":           {Covered: 1, Total: 10},
		"p/src/A.tsx":     {Covered: 3, Total: 10},
		"b/app/src/A.tsx": {Covered: 5, Total: 10},
	}
	mut := map[string]FileMutation{"A.tsx": {Killed: 1}, "z/src/A.tsx": {Killed: 2}, "y/src/A.tsx": {Killed: 3}}
	for i := 0; i < 50; i++ {
		g := ingestGraph()
		g.IngestExecution(exec, nil, nil, "unit", "t")
		g.IngestCoverage(cov, "t")
		g.IngestMutation(mut, "t")
		if p := g.Nodes[2].Signal.Passed; p != 2 {
			t.Fatalf("run %d: the test node took %d passed, want 2 (the exact path)", i, p)
		}
		if c := g.Nodes[1].Signal.CoveredLines; c != 3 {
			t.Fatalf("run %d: the code node took %d covered lines, want 3 (the closest in length)", i, c)
		}
		if k := g.Nodes[1].Signal.MutantsKilled; k != 3 {
			t.Fatalf("run %d: the code node took %d killed mutants, want 3 (a length tie goes to the first path in order)", i, k)
		}
	}
}

func TestIngestMutation_RecordsTheTimedOut(t *testing.T) {
	t.Run("SGINA-B21: A mutation ingestion records the timed-out apart", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "src/a.go", Kind: KindCode, Rev: "r1"}}}
	g.IngestMutation(map[string]FileMutation{"src/a.go": {Killed: 10, Survived: 2, TimedOut: 7, Score: 83}}, "now")
	if s := g.Nodes[0].Signal; s == nil || s.MutantsKilled != 10 || s.MutantsTimedOut != 7 {
		t.Fatalf("want 10 killed with 7 timed out recorded apart, got %+v", s)
	}
}

func TestIngestExecution_RecordsRunSeconds(t *testing.T) {
	t.Run("SGINA-B22: An execution ingestion records each test file's run time under its suite", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.test.ts", Kind: KindTest, Rev: "r"}, {ID: "b.test.ts", Kind: KindTest, Rev: "r"}}}
	exec := map[string]ExecByFile{"a.test.ts": {Passed: 2, Seconds: 1.5}, "b.test.ts": {Passed: 1}}
	g.IngestExecutionSuite(exec, nil, nil, nil, "unit", "out/junit.xml", "now")
	if s := g.Nodes[0].Signal; s == nil || s.SecondsBySuite["out/junit.xml"] != 1.5 {
		t.Fatalf("a.test.ts must record 1.5s under the suite, got %+v", s)
	}
	if s := g.Nodes[1].Signal; s == nil || len(s.SecondsBySuite) != 0 {
		t.Fatalf("a file whose cases carry no time records none, got %+v", s)
	}
	g.IngestExecutionSuite(exec, nil, nil, nil, "e2e", "", "now")
	if s := g.Nodes[0].Signal; s.SecondsBySuite["e2e"] != 1.5 {
		t.Fatalf("with no suite the time goes under the layer, got %+v", s.SecondsBySuite)
	}
}

func TestRecordRunSeconds(t *testing.T) {
	t.Run("SGINA-B23: A run time Anchors measured is recorded on the node", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.go", Kind: KindCode}}}
	g.RecordRunSeconds("a.go", "out/mutation.json", 12.5)
	g.RecordRunSeconds("gone.go", "out/mutation.json", 3)
	if s := g.Nodes[0].Signal; s == nil || s.SecondsBySuite["out/mutation.json"] != 12.5 || len(g.Nodes) != 1 {
		t.Fatalf("a.go must carry 12.5s under the suite and no node be added, got %+v (%d nodes)", s, len(g.Nodes))
	}
}

func TestRefreshRevs(t *testing.T) {
	t.Run("SGINA-B24: The tree's revs replace the map's", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Rev: "old"}, {ID: "a.go", Rev: "same"}, {ID: "a_test.go", Rev: "t1"}}}
	if n := g.RefreshRevs(map[string]string{"a.spec.md": "new", "a.go": "same", "gone.go": "x"}); n != 1 {
		t.Errorf("only the changed node is counted, got %d", n)
	}
	want := []string{"new", "same", "t1"}
	for i, n := range g.Nodes {
		if n.Rev != want[i] {
			t.Errorf("%s: rev %q, want %q", n.ID, n.Rev, want[i])
		}
	}
	if len(g.Nodes) != 3 {
		t.Errorf("a path the map does not have adds no node, got %d nodes", len(g.Nodes))
	}
}

func TestMutationStale(t *testing.T) {
	t.Run("SGINA-B25: The mutation result keeps its own rev", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.go", Kind: KindCode, Rev: "r1"}}}
	g.IngestMutationScoped(map[string]FileMutation{"a.go": {Killed: 9, Survived: 1, Score: 90}}, "", "now", 60, 80)
	if g.Nodes[0].Signal.MutationAtRev != "r1" || g.Nodes[0].MutationStale() {
		t.Fatalf("a fresh mutation records its rev and is not stale, got %+v", g.Nodes[0].Signal)
	}
	g.Nodes[0].Rev = "r2"
	g.IngestCoverage(map[string]FileCov{"a.go": {Total: 2, Covered: 2}}, "now")
	if g.Nodes[0].Signal.AtRev != "r2" {
		t.Fatalf("the probe needs the coverage to move the shared rev, got %q", g.Nodes[0].Signal.AtRev)
	}
	if !g.Nodes[0].MutationStale() {
		t.Error("a coverage ingestion after the change must not make the old mutation current")
	}
	old := Node{Rev: "r2", Signal: &TestSignal{AtRev: "r1"}}
	if !old.MutationStale() || (Node{Rev: "r1", Signal: &TestSignal{AtRev: "r1"}}).MutationStale() || (Node{}).MutationStale() {
		t.Error("with no mutation rev, the shared one decides; with no signal, nothing is stale")
	}
}

// A branch the unit suite missed and the integration suite took is taken; a stale suite
// speaks of other lines and does not take part.
func TestIngestCoverageSuite_branchUnion(t *testing.T) {
	t.Run("SGINA-B26: The branches are the union of the fresh suites", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 1, Lines: map[int]bool{1: true},
		Branches: map[string]bool{"10:0:1": false, "3:0:0": true, "3:0:1": false, "3:1:0": true}}}, "unit.info", "t1")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 1, Lines: map[int]bool{1: true},
		Branches: map[string]bool{"3:0:1": true, "10:0:1": false}}}, "it.info", "t2")
	g.IngestCoverageSuite(map[string]FileCov{"a.ts": {Covered: 1, Total: 1, Lines: map[int]bool{1: true}}}, "lines.info", "t3")
	s := g.Nodes[0].Signal
	if s.BranchTotal != 4 || s.BranchMissed != "10:0:1" {
		t.Errorf("want 4 branches with only 10:0:1 missed, got %d %q", s.BranchTotal, s.BranchMissed)
	}
	if got := s.CoverageBySuite["unit.info"].BranchMissed; got != "3:0:1 10:0:1" {
		t.Errorf("a suite keeps its missed branches in line order, got %q", got)
	}
	s.CoverageBySuite["unit.info"] = SuiteCoverage{BranchTotal: 9, BranchMissed: "1:0:0", AtRev: "old"}
	if total, missed := unionBranches(s.CoverageBySuite, "r1"); total != 2 || missed != "10:0:1" {
		t.Errorf("a stale suite does not take part, got %d %q", total, missed)
	}
	if total, missed := unionBranches(nil, "r1"); total != 0 || missed != "" {
		t.Errorf("no suite, no branch, got %d %q", total, missed)
	}
	g2 := &Graph{Nodes: []Node{{ID: "b.ts", Kind: KindCode, Rev: "r1"}}}
	g2.IngestCoverage(map[string]FileCov{"b.ts": {Covered: 1, Total: 1, Branches: map[string]bool{"2:0:0": false, "2:0:1": true}}}, "t")
	if s := g2.Nodes[0].Signal; s.BranchTotal != 2 || s.BranchMissed != "2:0:0" {
		t.Errorf("with no suite the report's own branches, got %d %q", s.BranchTotal, s.BranchMissed)
	}
	ids := []string{"2:10:0", "10:0:0", "2:9:1", "2:9:0"}
	sortBranchIDs(ids)
	if !reflect.DeepEqual(ids, []string{"2:9:0", "2:9:1", "2:10:0", "10:0:0"}) {
		t.Errorf("branch ids sort by number, part by part, got %v", ids)
	}
}

func TestIngestMutation_noCoverageLines(t *testing.T) {
	t.Run("SGINA-B27: The lines of mutants no test ran are recorded", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.ts", Kind: KindCode, Rev: "r1"}}}
	g.IngestMutation(map[string]FileMutation{"a.ts": {Killed: 1, NoCoverage: 3, NoCoverageAt: []int{9, 4, 5}}}, "t")
	if got := g.Nodes[0].Signal.NoCoverageLines(); !reflect.DeepEqual(got, []int{4, 5, 9}) {
		t.Errorf("want 4 5 9, got %v", got)
	}
	var none *TestSignal
	if none.NoCoverageLines() != nil {
		t.Error("no signal, no line")
	}
}

func TestIngestSuite_keepsTheVariant(t *testing.T) {
	t.Run("SGINA-B28: A proven variant is kept on the node that declares its rule", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "r1"}}}
	declared := map[string][]string{"a.spec.md": {"AAAAX-B02"}}
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAX-B02#01": true, "BBBBX-B01": true}, nil, declared, "unit", "junit.xml", "t1")
	if got := g.Nodes[0].Signal; got == nil || len(got.ProvenCodes) != 1 || got.ProvenCodes[0] != "AAAAX-B02#01" {
		t.Errorf("the spec keeps its proven variant and nothing else, got %+v", got)
	}
}

func TestIngestCoverage_listedAndOmitted(t *testing.T) {
	t.Run("SGINA-B29: A file the coverage report lists, or leaves out, says so", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{
		{ID: "src/a.ts", Kind: KindCode, Rev: "ra"}, {ID: "src/empty.ts", Kind: KindCode, Rev: "re"},
		{ID: "src/types.ts", Kind: KindCode, Rev: "rt"}, {ID: "other/x.ts", Kind: KindCode, Rev: "rx"},
	}}
	byFile := map[string]FileCov{"src/a.ts": {Covered: 1, Total: 2}, "src/empty.ts": {}}
	g.IngestCoverageSuite(byFile, "unit", "t1")
	covers := func(id string) bool { return strings.HasPrefix(id, "src/") }
	if n := g.MarkCoverageOmitted(byFile, covers); n != 1 {
		t.Errorf("only the covered file the report left out is marked, got %d", n)
	}
	by := map[string]*TestSignal{}
	for i := range g.Nodes {
		by[g.Nodes[i].ID] = g.Nodes[i].Signal
	}
	if by["src/a.ts"].CoverageRev != "ra" || by["src/empty.ts"].CoverageRev != "re" || by["src/empty.ts"].TotalLines != 0 {
		t.Errorf("a listed file records its revision, with lines or with none: %+v %+v", by["src/a.ts"], by["src/empty.ts"])
	}
	if by["src/types.ts"] == nil || by["src/types.ts"].CoverageOmitted != "rt" || by["other/x.ts"] != nil {
		t.Errorf("the omitted file the suite covers records its revision, and nothing else: %+v %+v", by["src/types.ts"], by["other/x.ts"])
	}
	g.IngestCoverageSuite(map[string]FileCov{"src/types.ts": {Covered: 1, Total: 1}}, "integration", "t2")
	if by := g.Node("src/types.ts").Signal; by.CoverageOmitted != "" || by.CoverageRev != "rt" {
		t.Errorf("a later report that lists it clears the omission, got %+v", by)
	}
}

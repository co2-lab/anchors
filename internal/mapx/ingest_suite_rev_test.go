package mapx

import "testing"

// A proof from one suite must not be laundered as fresh by ANOTHER suite running.
//
// Suite A proves B01 at rev1, the spec changes, suite B proves B02 at rev2. With one
// `AtRev` for the whole node, the map said the union was measured at rev2 — B01 included,
// which was never measured there. Probed before the fix: `stale=false`.
func TestIngestSuite_staleSuiteKeepsTheNodeStale(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev1"}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01", "AAAAA-B02"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, decl, "", "mobile.xml", "t1")
	g.Nodes[0].Rev = "rev2"
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true}, decl, "", "backend.xml", "t2")

	if !g.Nodes[0].SignalStale() {
		t.Fatalf("suite mobile.xml measured rev1 and did not run again: the node must be stale (at_rev=%s)",
			g.Nodes[0].Signal.AtRev)
	}

	// Once the stale suite runs again at the current rev, the union is fresh.
	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, decl, "", "mobile.xml", "t3")
	if g.Nodes[0].SignalStale() {
		t.Errorf("every contributing suite measured rev2 and the node is still stale (at_rev=%s)",
			g.Nodes[0].Signal.AtRev)
	}
}

// A suite that stops proving anything leaves the union and leaves no empty entry behind.
func TestIngestSuite_suiteThatProvesNothingLeaves(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev1"}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, decl, "", "mobile.xml", "t1")
	g.IngestExecutionSuite(nil, map[string]bool{}, decl, "", "mobile.xml", "t2")

	sig := g.Nodes[0].Signal
	if len(sig.ProvenCodes) != 0 {
		t.Errorf("the only suite stopped proving B01 and the union still has %v", sig.ProvenCodes)
	}
	if _, ok := sig.ProvenBySuite["mobile.xml"]; ok {
		t.Error("an empty suite entry was left in the versioned map")
	}
}

// An entry written before per-suite revs existed has unknown freshness — and unknown is
// not fresh.
func TestIngestSuite_legacyEntryWithoutRevIsStale(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev2",
		Signal: &TestSignal{ProvenBySuite: map[string][]string{"old.xml": {"AAAAA-B01"}}}}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01", "AAAAA-B02"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true}, decl, "", "new.xml", "t1")
	if !g.Nodes[0].SignalStale() {
		t.Error("a suite with no recorded rev was treated as fresh")
	}
}

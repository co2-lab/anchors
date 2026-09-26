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

// An entry written before per-suite revs existed has unknown freshness — and unknown is
// not fresh.
func TestIngestSuite_legacyEntryWithoutRevIsStale(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "a.spec.md", Kind: KindSpec, Rev: "rev2",
		Signal: &TestSignal{ProvenBySuite: map[string][]string{"old.xml": {"AAAAA-B01"}}}}}}
	decl := map[string][]string{"a.spec.md": {"AAAAA-B01", "AAAAA-B02"}}

	g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B02": true}, nil, decl, "", "new.xml", "t1")
	if !g.Nodes[0].SignalStale() {
		t.Error("a suite with no recorded rev was treated as fresh")
	}
}

// A PARTIAL run (`anchors test --changed`) measured only its cut. Reported from MIF: a run
// of 4 test files erased the proof of MoneyDetailScreen, whose test it never executed.
func TestIngestPartialRunKeepsWhatItDidNotSee(t *testing.T) {
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

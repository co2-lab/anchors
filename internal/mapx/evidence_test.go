// @anchors
//   ref: EVFRA

package mapx

import "testing"

// evidenceGraph — a script composing a util, in the real shape: SS-03 → login.yaml.
func evidenceGraph() *Graph {
	return &Graph{
		Nodes: []Node{
			{ID: "suites/SS-03.yaml", Kind: KindTest, Rev: "t1"},
			{ID: "utils/login.yaml", Kind: KindTest, Rev: "u1"},
			{ID: "utils/launchApp.yaml", Kind: KindTest, Rev: "l1"},
		},
		Edges: []Edge{
			{From: "suites/SS-03.yaml", To: "utils/login.yaml", Type: EdgeDependsOn},
			{From: "utils/login.yaml", To: "utils/launchApp.yaml", Type: EdgeDependsOn},
		},
	}
}

func TestEvidenceClosureDescendsTransitively(t *testing.T) {
	t.Run("EVFRA-B06: The closure descends transitively with current revisions", func(t *testing.T) {})
	t.Run("EVFRA-I01: The closure never contains the test itself", func(t *testing.T) {})
	// The closure must reach the util of the util: SS-03 → login → launchApp. Stopping at the
	// first level would leave out the deep dependencies, which are the ones nobody remembers to
	// revalidate by hand.
	g := evidenceGraph()
	closure := g.EvidenceClosure("suites/SS-03.yaml")
	if len(closure) != 2 {
		t.Fatalf("expected 2 nodes in the closure, got %v", closure)
	}
	if closure["utils/login.yaml"] != "u1" || closure["utils/launchApp.yaml"] != "l1" {
		t.Errorf("closure with wrong revs: %v", closure)
	}
}

func TestEvidenceExpiresWhenAComposedUtilChanges(t *testing.T) {
	t.Run("EVFRA-B03: A composed script that changed is named as the culprit", func(t *testing.T) {})
	t.Run("EVFRA-B04: Nothing moved means no verdict", func(t *testing.T) {})
	// THE CASE THAT MOTIVATED IT ALL: login.yaml was touched after 8 suites had passed. The
	// suite's signal does not change (its script is the same), so SignalStale() says it is
	// fresh — and it is wrong.
	g := evidenceGraph()
	g.Nodes[0].Signal = &TestSignal{
		Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1", "utils/launchApp.yaml": "l1"},
	}
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Fatalf("nothing changed yet — it should be fresh, got %+v", ev)
	}
	// the util changes; the suite's script does NOT
	g.Nodes[1].Rev = "u2"
	if g.Nodes[0].SignalStale() {
		t.Fatal("SignalStale() of the isolated node cannot catch this — that is why the closure exists")
	}
	ev := g.EvidenceStaleFor("suites/SS-03.yaml")
	if ev == nil {
		t.Fatal("the evidence HAD to expire: the composed util changed rev")
	}
	if ev.Own {
		t.Error("the script itself did not change — Own should be false")
	}
	if len(ev.Culprit) != 1 || ev.Culprit[0] != "utils/login.yaml" {
		t.Errorf("the culprit should be only login.yaml, got %v", ev.Culprit)
	}
}

func TestEvidenceWithoutSignalIsNotStale(t *testing.T) {
	t.Run("EVFRA-B01: A test that was never ingested has no verdict", func(t *testing.T) {})
	// Absence of proof is not stale proof: a test that never ran has no evidence to expire,
	// and reporting it here would mix "never measured" with "old measurement".
	g := evidenceGraph()
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Errorf("no signal = nothing to judge, got %+v", ev)
	}
	g.Nodes[0].Signal = &TestSignal{Passed: 1} // a signal with no ingestion rev
	g.Nodes[0].Rev = "t2"
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Errorf("a signal with no ingestion rev = nothing to judge, got %+v", ev)
	}
}

func TestEvidenceOldSignalWithoutClosureFallsBackToTheNode(t *testing.T) {
	t.Run("EVFRA-B05: A signal with no recorded closure is judged by its own file only", func(t *testing.T) {})
	t.Run("EVFRA-B02: The test's own change expires its evidence", func(t *testing.T) {})
	// A signal ingested by a version that did not record the closure: judge by the isolated
	// node, instead of declaring every earlier evidence stale — which would turn a precision
	// improvement into hundreds of false expirations the day it shipped.
	g := evidenceGraph()
	g.Nodes[0].Signal = &TestSignal{Passed: 1, AtRev: "t1"} // no ClosureRev
	g.Nodes[1].Rev = "u2"                                   // the util changed
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Errorf("with no recorded closure, the util does not count — got %+v", ev)
	}
	g.Nodes[0].Rev = "t2" // now the test ITSELF changed
	ev := g.EvidenceStaleFor("suites/SS-03.yaml")
	if ev == nil || !ev.Own {
		t.Errorf("the file itself changed: it should expire with Own=true, got %+v", ev)
	}
}

func TestEvidenceClosureStopsAtNoPropagation(t *testing.T) {
	t.Run("EVFRA-B07: A no-propagation node is in the closure but not walked through", func(t *testing.T) {})
	// A @noPropagation node states that changes in it do not flow down. Honouring it keeps a
	// deliberately volatile file from expiring the evidence of half the suite.
	g := evidenceGraph()
	g.Nodes[1].NoPropagation = true // login.yaml does not propagate
	closure := g.EvidenceClosure("suites/SS-03.yaml")
	if _, hasLaunch := closure["utils/launchApp.yaml"]; hasLaunch {
		t.Errorf("it should not descend THROUGH the noPropagation node, got %v", closure)
	}
	if _, hasLogin := closure["utils/login.yaml"]; !hasLogin {
		t.Error("the noPropagation node itself is in the closure (it IS a direct dependency)")
	}
}

func TestEvidenceIgnoresAVanishedClosureNode(t *testing.T) {
	t.Run("EVFRA-B08: A recorded node that left the graph is not a culprit", func(t *testing.T) {})
	g := evidenceGraph()
	g.Nodes[0].Signal = &TestSignal{Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1", "utils/deleted.yaml": "d1"}}
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Errorf("a file that left the graph has no current rev to compare, got %+v", ev)
	}
}

func TestEvidenceClosureNeverClimbs(t *testing.T) {
	t.Run("EVFRA-X01: The closure never climbs to the spec above the test", func(t *testing.T) {})
	g := evidenceGraph()
	g.Nodes = append(g.Nodes, Node{ID: "SS-03.spec.md", Kind: KindSpec, Rev: "s1"})
	g.Edges = append(g.Edges, Edge{From: "SS-03.spec.md", To: "suites/SS-03.yaml", Type: EdgeTestedBy})
	closure := g.EvidenceClosure("suites/SS-03.yaml")
	if _, ok := closure["SS-03.spec.md"]; ok {
		t.Errorf("the spec the test proves is not an input of the test, got %v", closure)
	}
	if _, ok := closure["suites/SS-03.yaml"]; ok {
		t.Errorf("the test is not part of its own closure, got %v", closure)
	}
}

package mapx

import "testing"

func ingestGraph() *Graph {
	return &Graph{Nodes: []Node{
		{ID: "src/A.spec.md", Kind: KindSpec, Rev: "r1"},
		{ID: "src/A.tsx", Kind: KindCode, Rev: "r1"},
		{ID: "src/A.test.tsx", Kind: KindTest, Rev: "r1"},
	}}
}

func TestIngestExecutionCrossesScenarios(t *testing.T) {
	g := ingestGraph()
	byFile := map[string]ExecByFile{"src/A.test.tsx": {Passed: 2, Failed: 1}}
	proven := map[string]bool{"AAAAX-V01": true, "AAAAX-V02": true} // V03 falhou → não provado
	declared := map[string][]string{"src/A.spec.md": {"AAAAX-V01", "AAAAX-V02", "AAAAX-V03"}}
	mf, mc := g.IngestExecution(byFile, proven, declared, "unit", "now")
	if mf != 1 {
		t.Fatalf("esperava 1 teste casado, veio %d", mf)
	}
	if mc != 2 {
		t.Fatalf("esperava 2 cenários provados, veio %d", mc)
	}
	// o nó de teste tem o resultado
	for _, n := range g.Nodes {
		if n.ID == "src/A.test.tsx" && (n.Signal == nil || n.Signal.Failed != 1) {
			t.Error("nó de teste deveria ter Failed=1")
		}
		if n.ID == "src/A.spec.md" {
			if n.Signal == nil || len(n.Signal.ProvenCodes) != 2 {
				t.Errorf("spec deveria ter 2 cenários provados, veio %+v", n.Signal)
			}
		}
	}
}

func TestIngestCoverageAndStale(t *testing.T) {
	g := ingestGraph()
	g.IngestCoverage(map[string]FileCov{"src/A.tsx": {Covered: 3, Total: 5}}, "now")
	var code *Node
	for i := range g.Nodes {
		if g.Nodes[i].ID == "src/A.tsx" {
			code = &g.Nodes[i]
		}
	}
	if code.Signal == nil || code.Signal.LineCoverage != 60 {
		t.Fatalf("esperava 60%% de cobertura, veio %+v", code.Signal)
	}
	if code.SignalStale() {
		t.Error("recém-ingerido não deveria ser stale")
	}
	code.Rev = "r2" // o arquivo mudou
	if !code.SignalStale() {
		t.Error("após mudar de rev, o sinal deveria ser stale")
	}
}

func TestPathMatches(t *testing.T) {
	if !pathMatches("apps/mobile/src/A.tsx", "/abs/apps/mobile/src/A.tsx") {
		t.Error("deveria casar por sufixo de caminho")
	}
	if pathMatches("a/b.go", "xa/b.go") {
		t.Error("não deveria casar sem fronteira de caminho")
	}
}

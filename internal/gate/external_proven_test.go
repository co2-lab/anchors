// @anchors
//   code: EXPRT
//   ref: EXPRS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func externalCase(layer string) (*mapx.Graph, *config.Config) {
	cfg := &config.Config{ExternalSystems: []config.ExternalSystem{{Name: "appsync", Calls: `client\.graphql\(`, ProvenBy: []string{"integration", "live"}}}}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "src/alerts.ts", Kind: mapx.KindCode},
			{ID: "src/alerts.test.ts", Kind: mapx.KindTest, Signal: &mapx.TestSignal{ByLayer: map[string]mapx.LayerExec{layer: {Passed: 3}}}},
		},
		Edges: []mapx.Edge{{From: "src/alerts.test.ts", To: "src/alerts.ts", Type: mapx.EdgeDependsOn}},
	}
	resetTestedUnits()
	return g, cfg
}

const callsAppSync = "export const mark = (id) => client.graphql({ query: updateAuditItem, variables: { id } })\n"

func TestExternalProven_aFileCallingARealSystem(t *testing.T) {
	t.Run("EXPRS-B01: A file calling a real system with no test passed in a suite that reaches it fails", func(t *testing.T) {})
	n := mapx.Node{ID: "src/alerts.ts", Kind: mapx.KindCode}
	g, cfg := externalCase("unit")
	v, msg := checkExternalProven(callsAppSync, n, "", g, cfg)
	if v != Fail || !strings.Contains(msg, "appsync") || !strings.Contains(msg, "integration, live") || !strings.Contains(msg, "1 test") {
		t.Errorf("a unit-only proof does not reach the real system: %v %s", v, msg)
	}
	g, cfg = externalCase("integration")
	if v, msg := checkExternalProven(callsAppSync, n, "", g, cfg); v != Pass {
		t.Errorf("an integration proof reaches it: %v %s", v, msg)
	}
}

func TestExternalProven_outsideItsSubject(t *testing.T) {
	t.Run("EXPRS-B02: Outside its subject the gate is skipped", func(t *testing.T) {})
	t.Run("EXPRS-E01: Without a map the gate is pending", func(t *testing.T) {})
	n := mapx.Node{ID: "src/alerts.ts", Kind: mapx.KindCode}
	g, cfg := externalCase("unit")
	if v, _ := checkExternalProven(callsAppSync, n, "", g, &config.Config{}); v != Skip {
		t.Errorf("no system declared: %v", v)
	}
	if v, _ := checkExternalProven(callsAppSync, mapx.Node{ID: "a.test.ts", Kind: mapx.KindTest}, "", g, cfg); v != Skip {
		t.Errorf("a test file: %v", v)
	}
	if v, _ := checkExternalProven("export const x = 1\n", n, "", g, cfg); v != Skip {
		t.Errorf("no call: %v", v)
	}
	waived := "// @no-external-proof: the call runs only behind the beta flag\n" + callsAppSync
	if v, _ := checkExternalProven(waived, n, "", g, cfg); v != Skip {
		t.Errorf("waived: %v", v)
	}
	if v, _ := checkExternalProven(callsAppSync, n, "", nil, cfg); v != Pending {
		t.Errorf("no map: %v", v)
	}
}

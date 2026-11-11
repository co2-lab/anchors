// @anchors
//   code: BYPRT
//   ref: BYPRS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func bypassCase(layer string) (*mapx.Graph, *config.Config) {
	cfg := &config.Config{TestBypass: &config.TestBypass{Patterns: []string{`E2E === 'true'\)\s*return`}, ProvenBy: []string{"integration"}}}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "src/useMarkSeen.ts", Kind: mapx.KindCode},
			{ID: "src/useMarkSeen.test.ts", Kind: mapx.KindTest, Signal: &mapx.TestSignal{ByLayer: map[string]mapx.LayerExec{layer: {Passed: 2}}}},
		},
		Edges: []mapx.Edge{{From: "src/useMarkSeen.test.ts", To: "src/useMarkSeen.ts", Type: mapx.EdgeDependsOn}},
	}
	resetTestedUnits()
	return g, cfg
}

const hookWithBypass = "export function useMarkSeen() {\n  if (process.env.EXPO_PUBLIC_E2E === 'true') return\n  markAuditSeen(id)\n}\n"

func TestBypassProven_aSkippedPath(t *testing.T) {
	t.Run("BYPRS-B01: A bypass with no test passed in a suite that runs its path fails", func(t *testing.T) {})
	n := mapx.Node{ID: "src/useMarkSeen.ts", Kind: mapx.KindCode}
	g, cfg := bypassCase("unit")
	if v, msg := checkBypassProven(hookWithBypass, n, "", g, cfg); v != Fail || !strings.Contains(msg, "line 2") || !strings.Contains(msg, "integration") {
		t.Errorf("only the unit suite ran: %v %s", v, msg)
	}
	g, cfg = bypassCase("integration")
	if v, msg := checkBypassProven(hookWithBypass, n, "", g, cfg); v != Pass {
		t.Errorf("the integration suite ran the path: %v %s", v, msg)
	}
}

func TestBypassProven_outsideItsSubject(t *testing.T) {
	t.Run("BYPRS-B02: Outside its subject the gate is skipped", func(t *testing.T) {})
	t.Run("BYPRS-E01: Without a map the gate is pending", func(t *testing.T) {})
	n := mapx.Node{ID: "src/useMarkSeen.ts", Kind: mapx.KindCode}
	g, cfg := bypassCase("unit")
	if v, _ := checkBypassProven(hookWithBypass, n, "", g, &config.Config{}); v != Skip {
		t.Errorf("none declared: %v", v)
	}
	if v, _ := checkBypassProven(hookWithBypass, mapx.Node{ID: "a.test.ts", Kind: mapx.KindTest}, "", g, cfg); v != Skip {
		t.Errorf("a test file: %v", v)
	}
	waived := strings.Replace(hookWithBypass, "  if (process", "  // @no-bypass-proof: the write is proven by the backend's own integration suite\n  if (process", 1)
	if v, _ := checkBypassProven(waived, n, "", g, cfg); v != Skip {
		t.Errorf("waived on the line above: %v", v)
	}
	if v, _ := checkBypassProven(hookWithBypass, n, "", nil, cfg); v != Pending {
		t.Errorf("no map: %v", v)
	}
}

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func splitCfg() *config.Config {
	return &config.Config{Derived: &config.Derived{Anchor: "code", Files: map[string]config.Padroes{
		"test": {"{{dir}}/{{name}}.test.ts", "{{dir}}/{{name}}.test.tsx", "{{dir}}/{{name}}.it.ts"},
	}}}
}

func splitGraph(tests map[string]string) *mapx.Graph {
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "src/pay.ts", Kind: mapx.KindCode}, {ID: "src/lone.ts", Kind: mapx.KindCode}}}
	for id, layer := range tests {
		g.Nodes = append(g.Nodes, mapx.Node{ID: id, Kind: mapx.KindTest, Layer: layer})
	}
	return g
}

func TestSingleTestPerUnit_aSplitFails(t *testing.T) {
	t.Run("SNGTU-B01: Two test files of one unit in one layer fail", func(t *testing.T) {})
	resetTestedUnits()
	root := t.TempDir()
	writeFile(t, root, "src/pay.test.ts", "test('a')\n")
	writeFile(t, root, "src/pay.test.tsx", "test('b')\n")
	g := splitGraph(map[string]string{"src/pay.test.ts": "unit", "src/pay.test.tsx": "unit"})
	v, msg := checkSingleTestPerUnit("", mapx.Node{ID: "src/pay.ts", Kind: mapx.KindCode}, root, g, splitCfg())
	if v != Fail || !strings.Contains(msg, "unit") || !strings.Contains(msg, "src/pay.test.ts, src/pay.test.tsx") {
		t.Fatalf("the split fails naming the layer and the files, got %v: %s", v, msg)
	}
	resetTestedUnits()
	g = splitGraph(map[string]string{"src/pay.test.ts": "unit", "src/pay.it.ts": "integration"})
	if v, msg := checkSingleTestPerUnit("", mapx.Node{ID: "src/pay.ts", Kind: mapx.KindCode}, root, g, splitCfg()); v != Pass {
		t.Errorf("one file per layer passes, got %v: %s", v, msg)
	}
}

func TestSingleTestPerUnit_aDeclaredSplitPasses(t *testing.T) {
	t.Run("SNGTU-B02: A declared split passes", func(t *testing.T) {})
	resetTestedUnits()
	root := t.TempDir()
	writeFile(t, root, "src/pay.test.ts", "test('a')\n")
	writeFile(t, root, "src/pay.test.tsx", "// @split-test: the native SDK needs its own mocks\ntest('b')\n")
	g := splitGraph(map[string]string{"src/pay.test.ts": "unit", "src/pay.test.tsx": "unit"})
	if v, msg := checkSingleTestPerUnit("", mapx.Node{ID: "src/pay.ts", Kind: mapx.KindCode}, root, g, splitCfg()); v != Pass {
		t.Fatalf("a declared split passes, got %v: %s", v, msg)
	}
}

func TestSingleTestPerUnit_skips(t *testing.T) {
	t.Run("SNGTU-B03: Not code, no test, no map", func(t *testing.T) {})
	resetTestedUnits()
	g := splitGraph(nil)
	if v, _ := checkSingleTestPerUnit("", mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}, "", g, splitCfg()); v != Skip {
		t.Errorf("a spec is skipped, got %v", v)
	}
	if v, _ := checkSingleTestPerUnit("", mapx.Node{ID: "src/lone.ts", Kind: mapx.KindCode}, "", g, splitCfg()); v != Skip {
		t.Errorf("a unit with no test is skipped, got %v", v)
	}
	if v, _ := checkSingleTestPerUnit("", mapx.Node{ID: "src/lone.ts", Kind: mapx.KindCode}, "", nil, splitCfg()); v != Pending {
		t.Errorf("without a map the gate is pending, got %v", v)
	}
}

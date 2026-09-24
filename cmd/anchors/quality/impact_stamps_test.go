package quality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// `check --changed <module>` brings in the tests whose `@contract` stamps point at the
// module: the pre-commit of whoever changes a function then runs `mock-stamped` on the
// doubles of that function, and a stale double blocks the change that made it stale.
func TestImpactOfBringsTheTestsThatStampTheChangedFile(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"src/hooks/balance.ts":       "export function useBalance(id) {\n  return id\n}\n",
		"src/screens/Home.test.tsx":  "// @contract: src/hooks/balance.ts | export function useBalance(id) { | 3 | deadbeef\njest.mock('@/src/hooks/balance')\n",
		"src/screens/Other.test.tsx": "it('x', () => {})\n",
	}
	g := &mapx.Graph{}
	for name, body := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		kind := mapx.KindCode
		if filepath.Ext(filepath.Base(name)) == ".tsx" {
			kind = mapx.KindTest
		}
		g.Nodes = append(g.Nodes, mapx.Node{ID: name, Kind: kind})
	}

	ids, err := impactOf(g, &config.Config{}, filepath.Join(root, "src/hooks/balance.ts"), root)
	if err != nil {
		t.Fatal(err)
	}
	has := map[string]bool{}
	for _, id := range ids {
		has[id] = true
	}
	if !has["src/screens/Home.test.tsx"] {
		t.Errorf("the test stamping the changed module is not in its impact: %v", ids)
	}
	if has["src/screens/Other.test.tsx"] {
		t.Errorf("a test with no stamp on the module entered its impact: %v", ids)
	}
}

// @anchors
//   code: MDTMP
//   ref: MDCMP

package mapcmd

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func depsGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "ui/Arena.tsx", Kind: mapx.KindCode, FileCode: "ARSCR", Code: "ARNAA"},
			{ID: "theme/tokens.ts", Kind: mapx.KindCode, FileCode: "TOKNS"},
			{ID: "theme/colors.ts", Kind: mapx.KindCode, FileCode: "CLRST"},
		},
		Edges: []mapx.Edge{
			{From: "ui/Arena.tsx", To: "theme/tokens.ts", Type: mapx.EdgeDependsOn},
			{From: "theme/tokens.ts", To: "theme/colors.ts", Type: mapx.EdgeDependsOn},
			{From: "theme/colors.ts", To: "theme/tokens.ts", Type: mapx.EdgeDependsOn},
		},
	}
}

func TestDepsTree(t *testing.T) {
	t.Run("MDCMP-B01: The tree goes down what a file uses, and up who uses it, a cycle marked once", func(t *testing.T) {})
	g := depsGraph()
	down := strings.Join(DepsTree(g, "ui/Arena.tsx", false, 0), "\n")
	if down != "ARSCR ui/Arena.tsx\n  TOKNS theme/tokens.ts\n    CLRST theme/colors.ts\n      TOKNS theme/tokens.ts ↺" {
		t.Errorf("down:\n%s", down)
	}
	if up := strings.Join(DepsTree(g, "theme/tokens.ts", true, 1), "\n"); up != "TOKNS theme/tokens.ts\n  CLRST theme/colors.ts\n  ARSCR ui/Arena.tsx" {
		t.Errorf("up, one level:\n%s", up)
	}
	t.Run("MDCMP-B02: A file is named by its own code, its unit's code or its path", func(t *testing.T) {})
	for arg, want := range map[string]string{"TOKNS": "theme/tokens.ts", "ARNAA": "ui/Arena.tsx", "theme/colors.ts": "theme/colors.ts", "NOPEX": ""} {
		if got := resolveDepsStart(g, arg); got != want {
			t.Errorf("%s → %q, want %q", arg, got, want)
		}
	}
}

func TestResourceUsers(t *testing.T) {
	t.Run("MDCMP-B03: The resources of a kind and the files that reach them", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.ts", FileCode: "AAAAA", Resources: []string{"db:transactions"}},
		{ID: "b.ts", FileCode: "BBBBB", Resources: []string{"api:stripe", "db:transactions"}},
	}}
	if got := strings.Join(ResourceUsers(g, "db"), "\n"); got != "db:transactions\n  AAAAA a.ts\n  BBBBB b.ts" {
		t.Errorf("got:\n%s", got)
	}
}

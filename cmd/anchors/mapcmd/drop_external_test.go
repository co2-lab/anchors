package mapcmd

import (
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// The external proof is dropped only by a FULL run of a suite inside the repository: a
// partial run measured only its cut, and an external report does not drop its own kind.
func TestDropExternal(t *testing.T) {
	graph := func() *mapx.Graph {
		g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "r1"}}}
		g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, nil,
			map[string][]string{"a.spec.md": {"AAAAA-B01"}}, "", "external/report.xml", "t")
		return g
	}
	for _, c := range []struct {
		name    string
		key     string
		partial bool
		drops   bool
	}{
		{"full in-repo run", ".anchors/junit.xml", false, true},
		{"partial in-repo run", ".anchors/junit.xml", true, false},
		{"another external report", "external/other.xml", false, false},
	} {
		g := graph()
		dropExternal(g, c.key, c.partial)
		_, kept := g.Nodes[0].Signal.ProvenBySuite["external/report.xml"]
		if kept == c.drops {
			t.Errorf("%s: external suite kept=%v, want dropped=%v", c.name, kept, c.drops)
		}
	}
}

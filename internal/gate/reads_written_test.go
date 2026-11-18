// @anchors
//   code: RDWRT
//   ref: RDWRS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testkit"
)

func TestReadsWritten(t *testing.T) {
	t.Run("RDWRS-B01: A data read that ties the file to no writer fails, naming it", func(t *testing.T) {})
	t.Run("RDWRS-B02: A file whose header reads nothing is skipped", func(t *testing.T) {})
	t.Run("RDWRS-E01: Without a map the gate is pending", func(t *testing.T) {})
	cfg := testkit.WithComments(nil)
	n := mapx.Node{ID: "reader.ts", Kind: mapx.KindCode}
	g := &mapx.Graph{Edges: []mapx.Edge{{From: "reader.ts", To: "writer.ts", Type: mapx.EdgeReadsFrom, Method: "AuditItem"}}}
	two := "// @anchors\n//   reads: AuditItem, Nobody\n\nexport const x = 1\n"
	if v, msg := checkReadsWritten(two, n, "", g, cfg); v != Fail || !strings.Contains(msg, "Nobody") || strings.Contains(msg, "AuditItem") {
		t.Errorf("the unwritten data is named: %v %s", v, msg)
	}
	one := "// @anchors\n//   reads: AuditItem\n\nexport const x = 1\n"
	if v, msg := checkReadsWritten(one, n, "", g, cfg); v != Pass {
		t.Errorf("every read tied: %v %s", v, msg)
	}
	if v, _ := checkReadsWritten("// @anchors\n//   code: XXXXX\n\nexport const x = 1\n", n, "", g, cfg); v != Skip {
		t.Errorf("no reads: %v", v)
	}
	if v, _ := checkReadsWritten(one, n, "", nil, cfg); v != Pending {
		t.Errorf("no map: %v", v)
	}
}

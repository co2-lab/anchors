// @anchors
//   code: DHTDP
//   ref: DEPHN

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestDependencyHonored_theTableIsGone(t *testing.T) {
	t.Run("DEPHN-B09: A spec still carrying a Dependencies table diverges, pointing at the migration", func(t *testing.T) {})
	spec := mapx.Node{ID: "pay.spec.md", Kind: mapx.KindSpec}
	table := "# Pay\n\n## Dependencies\n\n| Code | File | Method |\n| --- | --- | --- |\n| DEP1 | `a.ts` | `x` |\n| `DEP2` | `b.ts` | `y` |\n"
	v, msg := checkDependencyHonored(table, spec, "", nil, nil)
	if v != Diverge || !strings.Contains(msg, "2") || !strings.Contains(msg, "anchors migrate") {
		t.Errorf("the table diverges, naming its rows and the migration: %v %s", v, msg)
	}
	uses := "# Pay\n\n## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `PAYMT-B01` | DEP1 |\n"
	if v, msg := checkDependencyHonored(uses, spec, "", nil, nil); v != Pass {
		t.Errorf("a DEPn cited in a rule's uses is no table: %v %s", v, msg)
	}
	if v, _ := checkDependencyHonored(table, mapx.Node{ID: "pay.ts", Kind: mapx.KindCode}, "", nil, nil); v != Skip {
		t.Error("a code file leaves without a verdict")
	}
}

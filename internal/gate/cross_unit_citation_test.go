// @anchors
//   code: CUCTC
//   ref: CRUCT

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestCrossUnitCitation(t *testing.T) {
	t.Run("CRUCT-B01: A spec citing a rule of another unit that realizes no product rule fails, naming the code and its lines; a cited rule realizing the product passes, and a navigation's reference, a realized product rule, an alias, a revision and a retired line are no citation", func(t *testing.T) {})
	root := t.TempDir()
	wallet := "<!-- @anchors\n  code: WLLTW\n-->\n# Wallet\n\n| `WLLTW-B03` | computes the balance |\n| `WLLTW-B07` | rounds |\n"
	writeFile(t, root, "wallet.spec.md", wallet)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "home.spec.md", Kind: mapx.KindSpec, Code: "HOMEH"},
		{ID: "wallet.spec.md", Kind: mapx.KindSpec, Code: "WLLTW"},
		{ID: "budget.doctrine.md", Kind: mapx.KindProduct, Code: "BUDGT"},
	}}
	spec := "<!-- @anchors\n  code: HOMEH\n-->\n# Home\n\n## Rules\n\n" +
		"| `HOMEH-B01` | shows the balance, as `WLLTW-B03` computes it |\n" +
		"| `HOMEH-B02` | caps the budget @realizes BUDGT-B01 |\n" +
		"| `HOMEH-B03` | REF[HOMEH-B01]: the same |\n" +
		"| `HOMEH-B04` | @retired: HOMEH-R0001 — moved to WLLTW-B07 |\n\n" +
		"## Navigation\n\n### Out\n\n| Rule | Destination |\n| --- | --- |\n| `HOMEH-A01` | WalletScreen (`WLLTW-S01`) |\n\n" +
		"## Revisions\n\nHOMEH-R0001: the balance follows WLLTW-B03\n"
	v, msg := checkCrossUnitCitation(spec, mapx.Node{ID: "home.spec.md", Kind: mapx.KindSpec}, root, g, nil)
	if v != Fail || !strings.Contains(msg, "`WLLTW-B03` (8)") {
		t.Errorf("the rule citing the wallet's rule is named with its line: %v %s", v, msg)
	}
	for _, not := range []string{"BUDGT", "WLLTW-S01", "WLLTW-B07", "HOMEH-B01"} {
		if strings.Contains(msg, not) {
			t.Errorf("%s is no citation: %s", not, msg)
		}
	}
	// The balance rule moves up to the product: the wallet realizes it, and the citation is
	// of a shared rule that lives where it should.
	writeFile(t, root, "wallet.spec.md", strings.Replace(wallet, "computes the balance |", "computes the balance @realizes BUDGT-B02 |", 1))
	if v, msg := checkCrossUnitCitation(spec, mapx.Node{ID: "home.spec.md", Kind: mapx.KindSpec}, root, g, nil); v != Pass {
		t.Errorf("a cited rule realizing the product passes: %v %s", v, msg)
	}
	if v, _ := checkCrossUnitCitation(spec, mapx.Node{ID: "home.ts", Kind: mapx.KindCode}, root, g, nil); v != Skip {
		t.Error("a code file leaves without a verdict")
	}
}

// @anchors
//   code: VRPRT
//   ref: VRPRS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const variantsSpec = "<!-- @anchors\n  code: UAUDU\n-->\n# Alerts\n\n## Effects\n\n| Effect | Description |\n| --- | --- |\n" +
	"| `UAUDU-B02` | a tap marks the item seen @variants: audit, notification |\n" +
	"| `UAUDU-E02` | a failed write propagates @variants: B02 |\n"

func variantsCase(t *testing.T, spec, feature string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "alerts.spec.md", spec)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "alerts.spec.md", Kind: mapx.KindSpec}}}
	if feature != "" {
		writeFile(t, root, "alerts.feature", feature)
		g.Nodes = append(g.Nodes, mapx.Node{ID: "alerts.feature", Kind: mapx.KindFeature})
		g.Edges = []mapx.Edge{{From: "alerts.spec.md", To: "alerts.feature", Type: "covered-by"}}
	}
	return checkVariantsProven(spec, mapx.Node{ID: "alerts.spec.md", Kind: mapx.KindSpec}, root, g, nil)
}

func TestVariantsProven_eachVariantOfARule(t *testing.T) {
	t.Run("VRPRS-B01: A variant no scenario of its rule names fails, naming it", func(t *testing.T) {})
	spec := strings.Replace(variantsSpec, " @variants: B02", "", 1)
	audit := "Feature: Alerts\n\n  @UAUDU-B02 @variant:audit\n  Scenario: tapping an audit marks it\n    Given x\n    Then y\n"
	if v, msg := variantsCase(t, spec, audit); v != Fail || !strings.Contains(msg, "UAUDU-B02: notification") || strings.Contains(msg, "audit,") {
		t.Errorf("the notification has no scenario: %v %s", v, msg)
	}
	outline := audit + "\n  @UAUDU-B02#02\n  Scenario Outline: tapping an item of <kind>\n    Given x\n    Then y\n\n    Examples:\n      | kind |\n      | Notification |\n"
	if v, msg := variantsCase(t, spec, outline); v != Pass {
		t.Errorf("an Examples row names the notification, without case: %v %s", v, msg)
	}
}

func TestVariantsProven_aRuleFollowingAnother(t *testing.T) {
	t.Run("VRPRS-B02: A rule that follows another's variants diverges while one of them is unproven", func(t *testing.T) {})
	feature := "Feature: Alerts\n\n  @UAUDU-B02 @variant:audit @variant:notification\n  Scenario: tapping marks either kind\n    Given x\n    Then y\n\n" +
		"  @UAUDU-E02 @variant:audit\n  Scenario: a failed audit write propagates\n    Given x\n    Then y\n"
	v, msg := variantsCase(t, variantsSpec, feature)
	if v != Diverge || !strings.Contains(msg, "UAUDU-E02") || !strings.Contains(msg, "notification") || !strings.Contains(msg, "UAUDU-B02") {
		t.Errorf("the error rule is proven only on the audit: %v %s", v, msg)
	}
	both := strings.Replace(feature, "@UAUDU-E02 @variant:audit", "@UAUDU-E02 @variant:audit @variant:notification", 1)
	if v, msg := variantsCase(t, variantsSpec, both); v != Pass {
		t.Errorf("both proven: %v %s", v, msg)
	}
}

func TestVariantsProven_outsideItsSubject(t *testing.T) {
	t.Run("VRPRS-B03: A spec naming no variants, with no feature, or no spec at all, is skipped", func(t *testing.T) {})
	plain := strings.NewReplacer(" @variants: audit, notification", "", " @variants: B02", "").Replace(variantsSpec)
	if v, _ := variantsCase(t, plain, "Feature: x\n"); v != Skip {
		t.Errorf("no variants named: %v", v)
	}
	if v, _ := variantsCase(t, variantsSpec, ""); v != Skip {
		t.Errorf("no feature: %v", v)
	}
	if v, _ := checkVariantsProven("x", mapx.Node{ID: "a.go", Kind: mapx.KindCode}, "", nil, nil); v != Skip {
		t.Errorf("no spec: %v", v)
	}
}

func TestVariantsProven_errors(t *testing.T) {
	t.Run("VRPRS-E01: A rule following a rule that names no variants fails", func(t *testing.T) {})
	t.Run("VRPRS-E02: Without a map the gate is pending", func(t *testing.T) {})
	dangling := strings.Replace(variantsSpec, " @variants: audit, notification", "", 1)
	if v, msg := variantsCase(t, dangling, "Feature: x\n"); v != Fail || !strings.Contains(msg, "UAUDU-E02") || !strings.Contains(msg, "UAUDU-B02") {
		t.Errorf("following nothing: %v %s", v, msg)
	}
	if v, _ := checkVariantsProven(variantsSpec, mapx.Node{ID: "alerts.spec.md", Kind: mapx.KindSpec}, t.TempDir(), nil, nil); v != Pending {
		t.Errorf("no map: %v", v)
	}
}

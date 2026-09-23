package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// THE REAL CASE: a revert deleted `DTSTD-B10` — spec, code and test — and the `.feature`
// kept its scenario, because at the point of the revert that scenario did not exist yet
// and a parallel PR reintroduced it with no git conflict.
//
// Probed before this gate existed: `spec-feature-match` answered Pass with an EMPTY
// message.

const featureWithOrphan = "@UNITX\nFeature: X\n\n" +
	"  @UNITX-B01\n  Scenario: the one that stayed\n\n" +
	"  @UNITX-B10\n  Scenario: the one whose rule was REVERTED\n"

func featNodeRev() mapx.Node {
	return mapx.Node{ID: "x.feature", Kind: mapx.KindFeature, Code: "UNITX"}
}

// withSpec builds a project whose spec is linked to the feature by `covered-by`.
func withSpec(t *testing.T, spec string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.spec.md"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{
		Nodes: []mapx.Node{featNodeRev()},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy}},
	}
}

// --- feature-spec-match ---

func TestFeatureSpecMatch_scenarioWithoutRuleIsReported(t *testing.T) {
	t.Run("FSPMT-B01: a scenario whose rule the spec no longer declares is reported", func(t *testing.T) {})
	root, g := withSpec(t, "### UNITX-B01 — the rule that stayed\n")

	v, msg := checkFeatureSpecMatch(featureWithOrphan, featNodeRev(), root, g, nil)
	if v != Fail {
		t.Fatalf("scenario B10 has no rule and the verdict was %v", v)
	}
	if !strings.Contains(msg, "UNITX-B10") {
		t.Errorf("the verdict does not name the orphan: %q", msg)
	}
	if strings.Contains(msg, "UNITX-B01") {
		t.Errorf("it accused the scenario that DOES have a rule: %q", msg)
	}
}

func TestFeatureSpecMatch_allBackedByRulesPasses(t *testing.T) {
	t.Run("FSPMT-B02: every scenario backed by a declared rule passes", func(t *testing.T) {})
	root, g := withSpec(t, "### UNITX-B01 — stayed\n\n### UNITX-B10 — also stayed\n")

	if v, msg := checkFeatureSpecMatch(featureWithOrphan, featNodeRev(), root, g, nil); v != Pass {
		t.Errorf("both rules exist and the verdict was %v: %s", v, msg)
	}
}

// A scenario cites ANOTHER unit's rule to say what it runs against. Charging that to the
// local spec would have the gate ask the impossible — the mistake `scenario-coverage`
// already measured, with 18 scenarios charged to a spec that defined 6.
func TestFeatureSpecMatch_doesNotChargeAnotherUnitsCode(t *testing.T) {
	t.Run("FSPMT-X01: a code from another unit is not charged of this spec", func(t *testing.T) {})
	root, g := withSpec(t, "### UNITX-B01 — the rule\n")
	withNeighbour := "@UNITX\nFeature: X\n\n  @UNITX-B01 @OTHER-B07\n  Scenario: runs against the neighbour\n"

	if v, msg := checkFeatureSpecMatch(withNeighbour, featNodeRev(), root, g, nil); v != Pass {
		t.Errorf("the neighbour's code was charged: %v / %s", v, msg)
	}
}

func TestFeatureSpecMatch_skips(t *testing.T) {
	notFeature := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFeatureSpecMatch(featureWithOrphan, notFeature, "", nil, nil); v != Skip {
		t.Errorf("not a feature and the verdict was %v", v)
	}
	// No linked spec: Pending, because whoever charges the spec's existence is co-location.
	noEdge := &mapx.Graph{Nodes: []mapx.Node{featNodeRev()}}
	if v, _ := checkFeatureSpecMatch(featureWithOrphan, featNodeRev(), t.TempDir(), noEdge, nil); v != Pending {
		t.Errorf("with no linked spec Pending was expected, got %v", v)
	}
}

// --- test-feature-match ---

func testNodeRev() mapx.Node {
	return mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest, Code: "UNITX"}
}

func withFeature(t *testing.T, feat string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.feature"), []byte(feat), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{
		Nodes: []mapx.Node{testNodeRev()},
		Edges: []mapx.Edge{{From: "x.feature", To: "x.test.ts", Type: mapx.EdgeTestedBy}},
	}
}

// A green test over a reverted rule is worse than a missing test, because it ATTESTS.
func TestTestFeatureMatch_testProvesCodeWithoutScenario(t *testing.T) {
	t.Run("TFTMT-B01: a code the test claims to prove and no scenario declares is reported", func(t *testing.T) {})
	root, g := withFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: the one that stayed\n")
	test := "it('UNITX-B01 ok', () => {})\nit('UNITX-B10 orphan', () => {})\n"

	v, msg := checkTestFeatureMatch(test, testNodeRev(), root, g, nil)
	if v != Fail {
		t.Fatalf("the test proves B10 without a scenario and the verdict was %v", v)
	}
	if !strings.Contains(msg, "UNITX-B10") {
		t.Errorf("the verdict does not name the orphan: %q", msg)
	}
}

// COMMENTS OUT: a code cited in a comment is a reference, not proof.
func TestTestFeatureMatch_commentDoesNotCount(t *testing.T) {
	t.Run("TFTMT-X01: a code cited only in a comment is not a claim of proof", func(t *testing.T) {})
	root, g := withFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: the one that stayed\n")
	test := "// UNITX-B10 was reverted, see #811\nit('UNITX-B01 ok', () => {})\n"

	if v, msg := checkTestFeatureMatch(test, testNodeRev(), root, g, nil); v != Pass {
		t.Errorf("the code in a COMMENT was charged: %v / %s", v, msg)
	}
}

func TestTestFeatureMatch_doesNotChargeAnotherUnitsCode(t *testing.T) {
	root, g := withFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: the one that stayed\n")
	test := "it('UNITX-B01 ok', () => { buildFixture(OTHER-B07) })\n"

	if v, msg := checkTestFeatureMatch(test, testNodeRev(), root, g, nil); v != Pass {
		t.Errorf("the neighbour's code was charged: %v / %s", v, msg)
	}
}

func TestTestFeatureMatch_skips(t *testing.T) {
	notTest := mapx.Node{ID: "x.feature", Kind: mapx.KindFeature}
	if v, _ := checkTestFeatureMatch("", notTest, "", nil, nil); v != Skip {
		t.Errorf("not a test and the verdict was %v", v)
	}
	noEdge := &mapx.Graph{Nodes: []mapx.Node{testNodeRev()}}
	if v, _ := checkTestFeatureMatch("it('x')", testNodeRev(), t.TempDir(), noEdge, nil); v != Pending {
		t.Errorf("with no linked feature Pending was expected, got %v", v)
	}
}

// A VARIANT is not another rule.
//
// A feature numbers variants of the same requirement — `@DTTBD-B01#01`, `#02` — when a
// rule needs more than one scenario to be exercised. The spec declares the rule ONCE.
//
// MEASURED in the reference app: without handling the suffix, the gate reported 69
// features, ALL because of variants — eight scenarios of a feature whose rule exists and is
// declared. A gate that accuses what is right teaches people to ignore it, and would have
// buried the real case (the reverted `DTSTD-B10`) in the noise.
func TestFeatureSpecMatch_variantIsNotAnotherRule(t *testing.T) {
	t.Run("FSPMT-X02: a numbered variant resolves to the rule it varies", func(t *testing.T) {})
	root, g := withSpec(t, "### UNITX-B01 — the rule\n")
	withVariants := "@UNITX\nFeature: X\n\n" +
		"  @UNITX-B01#01\n  Scenario: first slice\n\n" +
		"  @UNITX-B01#02\n  Scenario: second slice\n"

	if v, msg := checkFeatureSpecMatch(withVariants, featNodeRev(), root, g, nil); v != Pass {
		t.Errorf("the variants belong to B01, which exists: %v / %s", v, msg)
	}
}

// And an ORPHAN variant is still reported — handling the suffix must not become amnesty.
func TestFeatureSpecMatch_variantOfMissingRuleIsReported(t *testing.T) {
	root, g := withSpec(t, "### UNITX-B01 — the rule that stayed\n")
	withOrphan := "@UNITX\nFeature: X\n\n  @UNITX-B10#01\n  Scenario: variant of the reverted one\n"

	v, msg := checkFeatureSpecMatch(withOrphan, featNodeRev(), root, g, nil)
	if v != Fail {
		t.Fatalf("B10 does not exist and the verdict was %v", v)
	}
	// The code AS WRITTEN, so whoever reads the verdict can find it in the file.
	if !strings.Contains(msg, "UNITX-B10#01") {
		t.Errorf("the verdict does not carry the code as written: %q", msg)
	}
}

func TestTestFeatureMatch_variantIsNotAnotherRule(t *testing.T) {
	t.Run("TFTMT-X02: a numbered variant resolves to the rule it varies", func(t *testing.T) {})
	root, g := withFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01#01\n  Scenario: slice\n")
	test := "it('UNITX-B01 exercises the rule', () => {})\n"

	if v, msg := checkTestFeatureMatch(test, testNodeRev(), root, g, nil); v != Pass {
		t.Errorf("the test proves B01, which the feature declares as a variant: %v / %s", v, msg)
	}
}

// A REVISION IS NOT A RULE, and no scenario is charged for it.
//
// `JDDTJ-R0002` is a revision — four digits —, and `anyCodeRE` matches `R00` because `R`
// is one of the canonical letters (for Rule) and the pattern reads two digits. The rest is
// left over, and the gate reported a code nobody wrote.
//
// MEASURED in the reference app: two of the three findings were this.
func TestTestFeatureMatch_revisionIsNotReported(t *testing.T) {
	t.Run("TFTMT-X03: a revision code is not charged as a rule", func(t *testing.T) {})
	root, g := withFeature(t, "@UNITX\nFeature: X\n\n  @UNITX-B01\n  Scenario: the rule\n")
	test := "describe('UNITX-R0002 — the decision', () => {\n  it('UNITX-B01 ok', () => {})\n})\n"

	if v, msg := checkTestFeatureMatch(test, testNodeRev(), root, g, nil); v != Pass {
		t.Errorf("the revision was charged as a rule: %v / %s", v, msg)
	}
}

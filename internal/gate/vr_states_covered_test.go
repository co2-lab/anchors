// @anchors
//   ref: VRSTC

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const specVR = `<!-- @anchors
  code: BUTTN
-->
# Button

## States

### BUTTN-S01: Enabled

### BUTTN-S02: Disabled
`

// featureVRStates is a feature with a VR scenario per given state, in the short form.
func featureVRStates(states ...string) string {
	var b strings.Builder
	b.WriteString("Feature: Button\n\n  @BUTTN-B01 @unit-level\n  Scenario: Clicks\n    Given x\n")
	for _, s := range states {
		b.WriteString("\n  @BUTTN-" + s + " @vr-level\n  Scenario: " + s + " looks like its baseline\n    Given x\n")
	}
	return b.String()
}

// vrFixture writes a Button unit under ui/ with its spec and the given files, and returns
// the root and the map, whose test nodes are the given paths.
func vrFixture(t *testing.T, files map[string]string, tests ...string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if _, ok := files["ui/Button.spec.md"]; !ok {
		files["ui/Button.spec.md"] = specVR
	}
	for rel, body := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		must(t, os.MkdirAll(filepath.Dir(p), 0o755))
		must(t, os.WriteFile(p, []byte(body), 0o644))
	}
	g := &mapx.Graph{}
	for _, id := range tests {
		g.Nodes = append(g.Nodes, mapx.Node{ID: id, Kind: mapx.KindTest})
	}
	return root, g
}

var buttonCode = mapx.Node{ID: "ui/Button.tsx", Kind: mapx.KindCode}

var vrChecks = map[string]func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string){
	"vr-states-covered": checkVRStatesCovered, "vr-scenarios-tested": checkVRScenariosTested,
	"vr-scenarios-of-states": checkVRScenariosOfStates, "vr-tests-of-scenarios": checkVRTestsOfScenarios,
}

func TestVR_nothingToConfront(t *testing.T) {
	t.Run("VRSTC-B01: Nothing to confront leaves every gate without a verdict", func(t *testing.T) {})
	t.Run("VRSTC-E01: A code file with no spec beside it leaves without a verdict", func(t *testing.T) {})
	root, g := vrFixture(t, map[string]string{"ui/Other.spec.md": "# Other\n### BUTTN-S01\n"})
	for name, check := range vrChecks {
		if v, _ := check("", mapx.Node{ID: "ui/Button.spec.md", Kind: mapx.KindSpec}, root, g, nil); v != Skip {
			t.Errorf("%s on a spec node: %v", name, v)
		}
		if v, d := check("", mapx.Node{ID: "ui/Button.styles.ts", Kind: mapx.KindCode}, root, g, nil); v != Skip || !strings.Contains(d, "spec") {
			t.Errorf("%s on a part of the unit: %v %s", name, v, d)
		}
		if v, _ := check("", mapx.Node{ID: "ui/Other.tsx", Kind: mapx.KindCode}, root, g, nil); v != Skip {
			t.Errorf("%s on a spec with no code: %v", name, v)
		}
	}
}

func TestVR_everyStateHasAScenario(t *testing.T) {
	t.Run("VRSTC-B02: Every state needs a VR scenario", func(t *testing.T) {})
	t.Run("VRSTC-E02: A unit with no feature has no VR scenario", func(t *testing.T) {})
	root, g := vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01")})
	if v, d := checkVRStatesCovered("", buttonCode, root, g, nil); v != Fail || !strings.Contains(d, "BUTTN-S02") ||
		strings.Contains(d, "BUTTN-S01") || !strings.Contains(d, "@vr-level") {
		t.Errorf("S02 has no VR scenario: %v %s", v, d)
	}
	root, g = vrFixture(t, map[string]string{})
	if v, d := checkVRStatesCovered("", buttonCode, root, g, nil); v != Fail || !strings.Contains(d, "BUTTN-S01, BUTTN-S02") {
		t.Errorf("no feature, no VR scenario for any state: %v %s", v, d)
	}
	root, g = vrFixture(t, map[string]string{"ui/Button.spec.md": "<!-- @anchors\n  code: BUTTN\n-->\n### BUTTN-B01\n"})
	if v, _ := checkVRStatesCovered("", buttonCode, root, g, nil); v != Skip {
		t.Errorf("a spec with no state: %v", v)
	}
}

func TestVR_everyScenarioHasATestAndAnImage(t *testing.T) {
	t.Run("VRSTC-B03: Every VR scenario needs a VR test", func(t *testing.T) {})
	t.Run("VRSTC-B04: Every VR scenario needs its baseline image, in any image format", func(t *testing.T) {})
	t.Run("VRSTC-E03: A baseline that cannot be found counts as none", func(t *testing.T) {})
	root, g := vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01", "S02"),
		"ui/Button.vr.test.ts": "toHaveScreenshot('BUTTN-VR-S01')", "ui/Button.BUTTN-VR-S01.svg": "<svg/>"}, "ui/Button.vr.test.ts")
	v, d := checkVRScenariosTested("", buttonCode, root, g, nil)
	if v != Fail || !strings.Contains(d, "no VR test: BUTTN-VR-S02") || !strings.Contains(d, "baseline image: BUTTN-VR-S02") ||
		strings.Contains(d, "BUTTN-VR-S01") || !strings.Contains(d, "Button.BUTTN-VR-<state>.png") {
		t.Errorf("S02 has no test and no image: %v %s", v, d)
	}
	root, g = vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01", "S02"),
		"ui/Button.vr.test.ts": "BUTTN-VR-S01 BUTTN-VR-S02", "ui/Button.BUTTN-VR-S01.png": "x", "ui/Button.BUTTN-VR-S02-dark.webp": "x"}, "ui/Button.vr.test.ts")
	if v, d := checkVRScenariosTested("", buttonCode, root, g, nil); v != Pass {
		t.Errorf("every scenario tested and captured: %v %s", v, d)
	}
	root, g = vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates()})
	if v, _ := checkVRScenariosTested("", buttonCode, root, g, nil); v != Skip {
		t.Errorf("no VR scenario: %v", v)
	}
}

func TestVR_everyScenarioIsOfAState(t *testing.T) {
	t.Run("VRSTC-B05: A VR scenario of a state the spec does not register", func(t *testing.T) {})
	t.Run("VRSTC-B06: A VR scenario of no state", func(t *testing.T) {})
	root, g := vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01", "S03") +
		"\n  @BUTTN-VR @vr-level\n  Scenario: every state\n    Given x\n"})
	v, d := checkVRScenariosOfStates("", buttonCode, root, g, nil)
	if v != Fail || !strings.Contains(d, "BUTTN-S03") || !strings.Contains(d, "of no state: BUTTN-VR") || strings.Contains(d, "BUTTN-S01") {
		t.Errorf("S03 is not a state, and BUTTN-VR is of none: %v %s", v, d)
	}
	root, g = vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01", "S02")})
	if v, d := checkVRScenariosOfStates("", buttonCode, root, g, nil); v != Pass {
		t.Errorf("every scenario of a state: %v %s", v, d)
	}
}

func TestVR_everyTestIsOfAScenario(t *testing.T) {
	t.Run("VRSTC-B07: A VR test of a scenario the feature does not declare", func(t *testing.T) {})
	root, g := vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01"),
		"ui/Button.vr.test.ts": "BUTTN-VR-S01 BUTTN-VR-S03"}, "ui/Button.vr.test.ts")
	if v, d := checkVRTestsOfScenarios("", buttonCode, root, g, nil); v != Fail || !strings.Contains(d, "BUTTN-VR-S03 (ui/Button.vr.test.ts)") ||
		strings.Contains(d, "BUTTN-VR-S01") {
		t.Errorf("S03's test has no scenario: %v %s", v, d)
	}
	root, g = vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01")})
	if v, _ := checkVRTestsOfScenarios("", buttonCode, root, g, nil); v != Skip {
		t.Errorf("no VR test: %v", v)
	}
}

func TestVR_readingTheUnit(t *testing.T) {
	t.Run("VRSTC-B08: A VR scenario carries the regime tag and the state's code, in either form", func(t *testing.T) {})
	t.Run("VRSTC-B09: Which tests are of the unit", func(t *testing.T) {})
	t.Run("VRSTC-E04: A test that cannot be read is read by its path", func(t *testing.T) {})
	feature := "Feature: Button\n  @BUTTN-S01 @vr-level\n  Scenario: a\n  @BUTTN-VR-S02 @vr-level\n  Scenario: b\n  @BUTTN-S01 @unit-level\n  Scenario: c\n"
	root, g := vrFixture(t, map[string]string{"ui/Button.feature": feature,
		".maestro/screens/Button/capture.yaml": "takeScreenshot: BUTTN-VR-S01",
		"ui/Button_vr_test.go":                 "// BUTTN-VR-S02",
		"other/Other.test.ts":                  "BUTTN-VR-S01"},
		".maestro/screens/Button/capture.yaml", "ui/Button_vr_test.go", "flows/BUTTN-VR-S02.yaml", "other/Other.test.ts", "ui/Button.BUTTN-VR-S01.png")
	u, _, ok := readVRUnit(buttonCode, root, g, nil)
	if !ok {
		t.Fatal("the unit is read")
	}
	if len(u.scenarios) != 2 || !u.scenarios["S01"] || !u.scenarios["S02"] || len(u.stateless) != 0 {
		t.Errorf("the VR scenarios are of S01 and S02: %v %v", u.scenarios, u.stateless)
	}
	if strings.Join(u.tests["S01"], ",") != ".maestro/screens/Button/capture.yaml" ||
		strings.Join(u.tests["S02"], ",") != "ui/Button_vr_test.go,flows/BUTTN-VR-S02.yaml" {
		t.Errorf("tests of the unit: %v", u.tests)
	}
}

func TestVR_oneStateNeverAnswersForAnother(t *testing.T) {
	t.Run("VRSTC-I01: One state never answers for another", func(t *testing.T) {})
	root, g := vrFixture(t, map[string]string{"ui/Button.feature": featureVRStates("S01"),
		"ui/Button.vr.test.ts": "BUTTN-VR-S01", "ui/Button.BUTTN-VR-S01.png": "x"}, "ui/Button.vr.test.ts")
	for name, check := range vrChecks {
		v, d := check("", buttonCode, root, g, nil)
		if strings.Contains(d, "S01") {
			t.Errorf("%s names S01, which is covered: %v %s", name, v, d)
		}
	}
	if v, d := checkVRStatesCovered("", buttonCode, root, g, nil); v != Fail || !strings.Contains(d, "BUTTN-S02") {
		t.Errorf("S02 is still uncovered: %v %s", v, d)
	}
}

func TestVR_stateLetter(t *testing.T) {
	t.Run("VRSTC-B10: The State letter comes from the project's rule types", func(t *testing.T) {})
	if l := stateLetter(nil); l != "S" {
		t.Errorf("default = %q", l)
	}
	cfg := &config.Config{RuleTypes: []config.RuleType{{Letter: "R", Term: "Rule"}, {Letter: "e", Term: "Estado"}}}
	if l := stateLetter(cfg); l != "E" {
		t.Errorf("from the rule types = %q", l)
	}
}

func TestVRBaseline_otherImageFormats(t *testing.T) {
	t.Run("VRSTC-B11: vr-baseline accepts other image formats", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "tela"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, "tela", "T.TCDTX-VR-loaded.jpg"), []byte("jpg"), 0o644))
	if v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil); v != Pass {
		t.Fatalf("a jpg baseline counts: %v (%s)", v, msg)
	}
}

func TestVR_aStateIsExemptedWithAReason(t *testing.T) {
	t.Run("VRSTC-B12: A state with no visual value is exempted with @no-vr and a reason", func(t *testing.T) {})
	spec := "<!-- @anchors\n  code: BUTTN\n-->\n# Button\n\n| State | Name | Note |\n| --- | --- | --- |\n" +
		"| BUTTN-S01 | Enabled | |\n| BUTTN-S02 | Loading | @no-vr: transient, the spinner has its own capture |\n\n" +
		"### BUTTN-S03: Pressed @no-vr\n"
	root, g := vrFixture(t, map[string]string{"ui/Button.spec.md": spec, "ui/Button.feature": featureVRStates("S02")})
	v, d := checkVRStatesCovered("", buttonCode, root, g, nil)
	if v != Fail || !strings.Contains(d, "BUTTN-S01, BUTTN-S03") || !strings.Contains(d, "with no reason: BUTTN-S03") ||
		strings.Contains(d, "BUTTN-S02") {
		t.Errorf("S02 is exempt; S03's exemption has no reason, so it is still asked: %v %s", v, d)
	}
	if v, d := checkVRScenariosOfStates("", buttonCode, root, g, nil); v != Fail || !strings.Contains(d, "exempts with `@no-vr`: BUTTN-S02") {
		t.Errorf("a VR scenario of an exempt state contradicts the spec: %v %s", v, d)
	}
}

func TestVR_messagesAndRegisteredStates(t *testing.T) {
	t.Run("VRSTC-B13: Every message is captured like a state", func(t *testing.T) {})
	t.Run("VRSTC-B14: The states registered are those of the States section", func(t *testing.T) {})
	withMessage := specVR + "\n## User Messages\n| Rule | Condition | Message |\n| --- | --- | --- |\n| `BUTTN-M01` | refused | \"Try again\" |\n"
	root, g := vrFixture(t, map[string]string{"ui/Button.spec.md": withMessage, "ui/Button.feature": featureVRStates("S01", "S02")})
	if v, d := checkVRStatesCovered("", buttonCode, root, g, nil); v != Fail || !strings.Contains(d, "BUTTN-M01") || strings.Contains(d, "BUTTN-S0") {
		t.Errorf("the message is asked for its capture: %v %s", v, d)
	}
	sectioned := "<!-- @anchors\n  code: BUTTN\n-->\n# Button\n\n## States (Estados da Tela)\n\n### BUTTN-S01: Enabled\n\n" +
		"## State Flow\n\n| From | Trigger | To |\n| --- | --- | --- |\n| `BUTTN-S01` | click | `BUTTN-S09` |\n"
	root, g = vrFixture(t, map[string]string{"ui/Button.spec.md": sectioned, "ui/Button.feature": featureVRStates("S01")})
	if v, d := checkVRStatesCovered("", buttonCode, root, g, nil); v != Pass {
		t.Errorf("S09 is cited, not registered: %v %s", v, d)
	}
}

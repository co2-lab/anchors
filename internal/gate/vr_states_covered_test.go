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

const featureVR = `Feature: Button

  @BUTTN-VR @vr-level
  Scenario: Each state of Button looks like its baseline
    Given each state is captured
`

// vrUnit writes a Button unit under ui/ with the given files and returns the root and
// the map, whose test nodes are the given test paths.
func vrUnit(t *testing.T, files map[string]string, tests ...string) (string, *mapx.Graph) {
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

func vrCheck(root string, g *mapx.Graph) (Verdict, string) {
	return checkVRStatesCovered("", mapx.Node{ID: "ui/Button.tsx", Kind: mapx.KindCode}, root, g, nil)
}

func TestVRStates_nothingToCover(t *testing.T) {
	t.Run("VRSTC-B01: Specs that have nothing to cover leave without a verdict", func(t *testing.T) {})
	root, _ := vrUnit(t, map[string]string{})
	if v, _ := checkVRStatesCovered("", mapx.Node{ID: "ui/Button.spec.md", Kind: mapx.KindSpec}, root, nil, nil); v != Skip {
		t.Errorf("a node that is not code: %v", v)
	}
	if v, _ := checkVRStatesCovered("", mapx.Node{ID: "ui/Button.styles.ts", Kind: mapx.KindCode}, root, nil, nil); v != Skip {
		t.Errorf("a part of the unit, with no spec of its own: %v", v)
	}
	for name, spec := range map[string]string{"no code": "# Button\n### BUTTN-S01\n",
		"no state": "<!-- @anchors\n  code: BUTTN\n-->\n# Button\n### BUTTN-B01\n"} {
		root, _ := vrUnit(t, map[string]string{"ui/Button.spec.md": spec})
		if v, _ := vrCheck(root, nil); v != Skip {
			t.Errorf("a spec with %s: %v", name, v)
		}
	}
}

func TestVRStates_needsTheScenario(t *testing.T) {
	t.Run("VRSTC-B02: The unit's feature must declare the VR scenario", func(t *testing.T) {})
	root, g := vrUnit(t, map[string]string{"ui/Button.feature": "Feature: Button\n  @BUTTN-S01\n  Scenario: x\n",
		"ui/Button.BUTTN-VR-S01.png": "x", "ui/Button.BUTTN-VR-S02.png": "x"}, "ui/BUTTN-VR.yaml")
	v, d := vrCheck(root, g)
	if v != Fail || !strings.Contains(d, "BUTTN-VR") || !strings.Contains(d, "@vr-level") {
		t.Errorf("no VR scenario fails naming it and the tag: %v %s", v, d)
	}
}

func TestVRStates_everyStateNeedsItsImage(t *testing.T) {
	t.Run("VRSTC-B03: Every state needs its baseline image, in any image format", func(t *testing.T) {})
	t.Run("VRSTC-I01: One state's image never covers another", func(t *testing.T) {})
	root, g := vrUnit(t, map[string]string{"ui/Button.feature": featureVR, "ui/Button.BUTTN-VR-S01.svg": "<svg/>"}, "ui/BUTTN-VR.yaml")
	v, d := vrCheck(root, g)
	if v != Fail || !strings.Contains(d, "BUTTN-S02") || strings.Contains(d, "BUTTN-S01") || !strings.Contains(d, "Button.BUTTN-VR-<state>.png") {
		t.Errorf("S02 has no image and S01's does not cover it: %v %s", v, d)
	}
}

func TestVRStates_theCapture(t *testing.T) {
	t.Run("VRSTC-B04: A capture test is a flow named by the VR code or a screenshot test beside the unit", func(t *testing.T) {})
	images := map[string]string{"ui/Button.feature": featureVR, "ui/Button.BUTTN-VR-S01.png": "x", "ui/Button.BUTTN-VR-S02-dark.jpg": "x"}
	root, g := vrUnit(t, images, "ui/Button.BUTTN-VR-S01.png")
	if v, d := vrCheck(root, g); v != Fail || !strings.Contains(d, "no test captures `BUTTN-VR`") {
		t.Errorf("an image is not the capture: %v %s", v, d)
	}
	root, g = vrUnit(t, images, ".maestro/ui/BUTTN-VR.yaml")
	if v, d := vrCheck(root, g); v != Pass {
		t.Errorf("a flow named by the VR code captures it: %v %s", v, d)
	}
	withShot := map[string]string{"ui/Button.vr.test.ts": "expect(page).toHaveScreenshot('BUTTN-VR-S01')"}
	for k, v := range images {
		withShot[k] = v
	}
	root, g = vrUnit(t, withShot, "ui/Button.vr.test.ts")
	if v, d := vrCheck(root, g); v != Pass {
		t.Errorf("a screenshot test beside the unit naming the VR captures it: %v %s", v, d)
	}
}

func TestVRStates_coveredPasses(t *testing.T) {
	t.Run("VRSTC-B05: Scenario, images and capture pass", func(t *testing.T) {})
	root, g := vrUnit(t, map[string]string{"ui/Button.feature": featureVR, "ui/Button.BUTTN-VR-S01.png": "x",
		"ui/Button.BUTTN-VR-S02.webp": "x"}, "ui/BUTTN-VR.yaml")
	if v, d := vrCheck(root, g); v != Pass {
		t.Errorf("all covered passes: %v %s", v, d)
	}
}

func TestVRStates_stateLetter(t *testing.T) {
	t.Run("VRSTC-B06: The State letter comes from the project's rule types", func(t *testing.T) {})
	if l := stateLetter(nil); l != "S" {
		t.Errorf("default = %q", l)
	}
	cfg := &config.Config{RuleTypes: []config.RuleType{{Letter: "R", Term: "Rule"}, {Letter: "e", Term: "Estado"}}}
	if l := stateLetter(cfg); l != "E" {
		t.Errorf("from the rule types = %q", l)
	}
}

func TestVRBaseline_otherImageFormats(t *testing.T) {
	t.Run("VRSTC-B07: vr-baseline accepts other image formats", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "tela"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, "tela", "T.TCDTX-VR-loaded.jpg"), []byte("jpg"), 0o644))
	if v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil); v != Pass {
		t.Fatalf("a jpg baseline counts: %v (%s)", v, msg)
	}
}

func TestVRStates_failures(t *testing.T) {
	t.Run("VRSTC-E01: A code file with no spec beside it leaves without a verdict", func(t *testing.T) {})
	t.Run("VRSTC-E02: A unit with no feature has no VR scenario", func(t *testing.T) {})
	t.Run("VRSTC-E03: A state whose image cannot be found counts as without one", func(t *testing.T) {})
	t.Run("VRSTC-E04: An unreadable test beside the unit is not the capture", func(t *testing.T) {})
	root, g := vrUnit(t, map[string]string{"ui/Button.BUTTN-VR-S01.png": "x"}, "ui/Button.vr.test.ts")
	if v, d := checkVRStatesCovered("", mapx.Node{ID: "ui/Button.styles.ts", Kind: mapx.KindCode}, root, g, nil); v != Skip || !strings.Contains(d, "spec") {
		t.Errorf("no spec beside it: %v %s", v, d)
	}
	v, d := vrCheck(root, g)
	if v != Fail || !strings.Contains(d, "`BUTTN-VR`") || !strings.Contains(d, "BUTTN-S02") || !strings.Contains(d, "no test captures") {
		t.Errorf("no feature, no S02 image and an unreadable test are each a gap: %v %s", v, d)
	}
}


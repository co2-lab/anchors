// @anchors
//   ref: VRCPT

package mapx

import (
	"sort"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/scan"
)

func captureFiles() []scan.File {
	return []scan.File{
		{Path: "ui/Button.spec.md", Kind: "spec", HeaderCode: "BUTTN", Rev: "s1"},
		{Path: "ui/Button.tsx", Kind: "code", Rev: "c1"},
		{Path: "ui/Icon.tsx", Kind: "code", Rev: "i1"},
		{Path: "ui/Button.BUTTN-VR-S01.png", Kind: "test", Rev: "p1"},
		{Path: "ui/Button.BUTTN-VR-S02-dark.svg", Kind: "test", Rev: "p2"},
		{Path: ".maestro/BUTTN-VR-S01.yaml", Kind: "test", Rev: "f1"},
		{Path: "ui/Button.vr.test.ts", Kind: "test", Codes: []string{"BUTTN-VR-S02"}, Rev: "t1"},
		{Path: "ui/Button.test.ts", Kind: "test", Codes: []string{"BUTTN-B01"}, Rev: "t2"},
		{Path: ".maestro/ZZZZZ-VR.yaml", Kind: "test", Rev: "f2"},
	}
}

func capturesFrom(edges []Edge, from string) []string {
	var out []string
	for _, e := range edges {
		if e.Type == EdgeCaptures && e.From == from {
			out = append(out, e.To)
		}
	}
	sort.Strings(out)
	return out
}

func TestCaptureEdges_aVRTestCapturesItsUnit(t *testing.T) {
	t.Run("VRCPT-B01: A VR test captures its unit's code file and images", func(t *testing.T) {})
	t.Run("VRCPT-B02: What captures nothing gets no edge", func(t *testing.T) {})
	edges := captureEdges(captureFiles())
	want := "ui/Button.BUTTN-VR-S01.png,ui/Button.BUTTN-VR-S02-dark.svg,ui/Button.tsx"
	for _, test := range []string{".maestro/BUTTN-VR-S01.yaml", "ui/Button.vr.test.ts"} {
		if got := strings.Join(capturesFrom(edges, test), ","); got != want {
			t.Errorf("%s captures %s, want %s", test, got, want)
		}
	}
	for _, none := range []string{"ui/Button.test.ts", "ui/Button.BUTTN-VR-S01.png", ".maestro/ZZZZZ-VR.yaml"} {
		if got := capturesFrom(edges, none); len(got) != 0 {
			t.Errorf("%s captures nothing, got %v", none, got)
		}
	}
}

// captureGraph is the graph of the files, with the screen depending on the component.
func captureGraph() *Graph {
	g := &Graph{}
	for _, f := range captureFiles() {
		g.Nodes = append(g.Nodes, Node{ID: f.Path, Kind: Kind(f.Kind), Rev: f.Rev})
	}
	g.Edges = append(captureEdges(captureFiles()), Edge{From: "ui/Button.tsx", To: "ui/Icon.tsx", Type: EdgeDependsOn})
	return g
}

func TestEvidenceClosure_aCaptureIsOneLevel(t *testing.T) {
	t.Run("VRCPT-B03: The closure of a capture is one level", func(t *testing.T) {})
	closure := captureGraph().EvidenceClosure(".maestro/BUTTN-VR-S01.yaml")
	if _, ok := closure["ui/Button.tsx"]; !ok {
		t.Errorf("the screen is in the closure: %v", closure)
	}
	if _, ok := closure["ui/Button.BUTTN-VR-S01.png"]; !ok {
		t.Errorf("the baseline is in the closure: %v", closure)
	}
	if _, ok := closure["ui/Icon.tsx"]; ok {
		t.Errorf("the component the screen uses has its own capture, and is not in this one's: %v", closure)
	}
}

func TestScan_aVRCodeIsReadWithItsState(t *testing.T) {
	t.Run("VRCPT-B04: A VR code is read whole with its state", func(t *testing.T) {})
	if got := scan.ScenarioCodeRE().FindString("expect(page).toHaveScreenshot('BUTTN-VR-S01')"); got != "BUTTN-VR-S01" {
		t.Errorf("scanned %q", got)
	}
}

func TestCapture_screenChangeStalesComponentChangeDoesNot(t *testing.T) {
	t.Run("VRCPT-I01: The screen's change stales its capture, a component's does not", func(t *testing.T) {})
	g := captureGraph()
	flow := g.node(".maestro/BUTTN-VR-S01.yaml")
	flow.Signal = &TestSignal{AtRev: flow.Rev, ClosureRev: g.EvidenceClosure(flow.ID)}
	g.node("ui/Icon.tsx").Rev = "i2"
	if ev := g.EvidenceStaleFor(flow.ID); ev != nil {
		t.Errorf("a component's change leaves the screen's capture fresh: %+v", ev)
	}
	g.node("ui/Button.tsx").Rev = "c2"
	if ev := g.EvidenceStaleFor(flow.ID); ev == nil || strings.Join(ev.Culprit, ",") != "ui/Button.tsx" {
		t.Errorf("the screen's change stales its capture: %+v", ev)
	}
}

func TestCaptureEdges_aContractTestCapturesItsAPI(t *testing.T) {
	t.Run("VRCPT-B05: A contract test captures its API unit, its spec and the OpenAPI document", func(t *testing.T) {})
	files := []scan.File{
		{Path: "api/generate.spec.md", Kind: "spec", HeaderCode: "GENAP"},
		{Path: "api/generate.go", Kind: "code"},
		{Path: "docs/openapi.yaml", Kind: "doc"},
		{Path: "api/generate_contract_test.go", Kind: "test", Codes: []string{"GENAP-CT"}},
	}
	got := strings.Join(capturesFrom(captureEdges(files), "api/generate_contract_test.go"), ",")
	if got != "api/generate.go,api/generate.spec.md,docs/openapi.yaml" {
		t.Errorf("a contract test captures %s", got)
	}
}

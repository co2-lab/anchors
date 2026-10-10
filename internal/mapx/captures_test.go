// @anchors
//   code: CPTSC
//   ref: VRCPT

package mapx

import (
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
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
		{Path: "ui/Button.vr.test.ts", Kind: "test", CaptureCodes: []string{"BUTTN-VR-S02"}, Rev: "t1"},
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
	for test, want := range map[string]string{
		".maestro/BUTTN-VR-S01.yaml": "ui/Button.BUTTN-VR-S01.png,ui/Button.tsx",
		"ui/Button.vr.test.ts":       "ui/Button.BUTTN-VR-S02-dark.svg,ui/Button.tsx",
	} {
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
	// The component has a capture of its own: that is what keeps it out of the screen's.
	g.Nodes = append(g.Nodes, Node{ID: ".maestro/ICONX-VR-S01.yaml", Kind: KindTest, Rev: "x1"})
	g.Edges = append(g.Edges, Edge{From: ".maestro/ICONX-VR-S01.yaml", To: "ui/Icon.tsx", Type: EdgeCaptures})
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
		{Path: "api/generate_contract_test.go", Kind: "test", CaptureCodes: []string{"GENAP-CT"}},
	}
	got := strings.Join(capturesFrom(captureEdges(files), "api/generate_contract_test.go"), ",")
	if got != "api/generate.go,api/generate.spec.md,docs/openapi.yaml" {
		t.Errorf("a contract test captures %s", got)
	}
}

// chainGraph: a captured screen whose spec depends on a hook (which depends on a store) and
// composes a captured component.
func chainGraph() *Graph {
	g := &Graph{}
	for _, n := range []Node{
		{ID: "ui/Arena.spec.md", Kind: KindSpec, Rev: "a"}, {ID: "ui/Arena.tsx", Kind: KindCode, Rev: "b"},
		{ID: "hooks/useMatches.ts", Kind: KindCode, Rev: "c"}, {ID: "stores/league.ts", Kind: KindCode, Rev: "d"},
		{ID: "ui/Sheet.tsx", Kind: KindCode, Rev: "e"},
		{ID: "flows/ARENA-VR-S01.yaml", Kind: KindTest, Rev: "f"}, {ID: "flows/SHEET-VR-S01.yaml", Kind: KindTest, Rev: "g"},
	} {
		g.Nodes = append(g.Nodes, n)
	}
	g.Edges = []Edge{
		{From: "ui/Arena.spec.md", To: "ui/Arena.tsx", Type: EdgeSpecifies},
		{From: "ui/Arena.spec.md", To: "hooks/useMatches.ts", Type: EdgeDependsOn},
		{From: "hooks/useMatches.ts", To: "stores/league.ts", Type: EdgeDependsOn},
		{From: "ui/Arena.spec.md", To: "ui/Sheet.tsx", Type: EdgeComposes},
		{From: "ui/Arena.spec.md", To: "ui/Sheet.tsx", Type: EdgeDependsOn},
		{From: "flows/ARENA-VR-S01.yaml", To: "ui/Arena.tsx", Type: EdgeCaptures},
		{From: "flows/SHEET-VR-S01.yaml", To: "ui/Sheet.tsx", Type: EdgeCaptures},
	}
	return g
}

func TestCaptureClosure_reachesUncapturedDependencies(t *testing.T) {
	t.Run("VRCPT-B06: A capture's closure reaches the uncaptured dependencies, and stops at a captured one", func(t *testing.T) {})
	c := chainGraph().EvidenceClosure("flows/ARENA-VR-S01.yaml")
	for _, want := range []string{"ui/Arena.tsx", "hooks/useMatches.ts", "stores/league.ts"} {
		if _, ok := c[want]; !ok {
			t.Errorf("the closure lacks %s: %v", want, c)
		}
	}
	if _, ok := c["ui/Sheet.tsx"]; ok {
		t.Errorf("the component has its own capture, and is not in the screen's: %v", c)
	}
}

func TestComposesEdges(t *testing.T) {
	t.Run("VRCPT-B07: Parts Used names become composes edges", func(t *testing.T) {})
	files := []scan.File{
		{Path: "ui/Arena.spec.md", Kind: "spec", Composes: []string{"Sheet", "Missing"}},
		{Path: "ui/Sheet.tsx", Kind: "code"},
	}
	edges := composesEdges(files)
	if len(edges) != 1 || edges[0].To != "ui/Sheet.tsx" || edges[0].Type != EdgeComposes {
		t.Errorf("edges = %+v", edges)
	}
}

func TestCapturesReaching(t *testing.T) {
	t.Run("VRCPT-B08: The captures a changed file reaches", func(t *testing.T) {})
	g := chainGraph()
	if got := strings.Join(g.CapturesReaching([]string{"stores/league.ts"}), ","); got != "flows/ARENA-VR-S01.yaml" {
		t.Errorf("the store reaches the screen's capture: %s", got)
	}
	if got := strings.Join(g.CapturesReaching([]string{"ui/Sheet.tsx"}), ","); got != "flows/SHEET-VR-S01.yaml" {
		t.Errorf("the component reaches only its own capture: %s", got)
	}
}

// flaggedApp is a captured screen whose code imports a hook, which imports a store, each
// import carrying its `@dep:` flag — no spec declares the hook or the store.
func flaggedApp() *Graph {
	return Build([]scan.File{
		{Path: "ui/Arena.spec.md", Kind: "spec", HeaderCode: "ARENA", Rev: "a"},
		{Path: "ui/Arena.tsx", Kind: "code", HeaderCode: "ARNAA", Rev: "b", CodeDeps: []scan.CodeDep{{Code: "USMTC", Line: 1}}},
		{Path: "hooks/useMatches.ts", Kind: "code", HeaderCode: "USMTC", Rev: "c", CodeDeps: []scan.CodeDep{{Code: "LEAGS", Line: 1}}},
		{Path: "stores/league.ts", Kind: "code", HeaderCode: "LEAGS", Rev: "d"},
		{Path: "flows/ARENA-VR-S01.yaml", Kind: "test", Rev: "f"},
	}, &config.Config{}, nil)
}

func TestFlaggedChain_feedsTheCaptureAndTheImpact(t *testing.T) {
	t.Run("VRCPT-B09: A dependency the code's flags declare reaches the capture's closure and the impact of a change, transitively", func(t *testing.T) {})
	g := flaggedApp()
	c := g.EvidenceClosure("flows/ARENA-VR-S01.yaml")
	for _, want := range []string{"ui/Arena.tsx", "hooks/useMatches.ts", "stores/league.ts"} {
		if _, ok := c[want]; !ok {
			t.Errorf("the capture's closure lacks %s, flagged in the chain: %v", want, c)
		}
	}
	flow := g.node("flows/ARENA-VR-S01.yaml")
	flow.Signal = &TestSignal{AtRev: flow.Rev, ClosureRev: c}
	g.node("stores/league.ts").Rev = "d2"
	if ev := g.EvidenceStaleFor(flow.ID); ev == nil || !slices.Contains(ev.Culprit, "stores/league.ts") {
		t.Errorf("a change to the store stales the screen's capture: %+v", ev)
	}
	up := g.AnalyzeImpact("stores/league.ts").Validate
	for _, want := range []string{"hooks/useMatches.ts", "ui/Arena.tsx"} {
		if !slices.Contains(up, want) {
			t.Errorf("a change to the store reaches %s, who imports it through the chain: %v", want, up)
		}
	}
}

func TestCaptureEdges_whatTheTestComparesAgainst(t *testing.T) {
	t.Run("VRCPT-B10: A test captures the states it names, and no code a comment mentions", func(t *testing.T) {})
	files := append(captureFiles(),
		scan.File{Path: ".maestro/BUTTN-VR.yaml", Kind: "test"},
		scan.File{Path: ".maestro/ALRTS-VR.yaml", Kind: "test", CaptureCodes: nil},
	)
	edges := captureEdges(files)
	if got := strings.Join(capturesFrom(edges, ".maestro/BUTTN-VR.yaml"), ","); got != "ui/Button.BUTTN-VR-S01.png,ui/Button.BUTTN-VR-S02-dark.svg,ui/Button.tsx" {
		t.Errorf("a code with no state captures every state: %s", got)
	}
	if got := capturesFrom(edges, ".maestro/ALRTS-VR.yaml"); len(got) != 0 {
		t.Errorf("a code no spec declares captures nothing: %v", got)
	}
	// The comment is the scan's: a code on a comment line is no capture code.
	flow := []byte("appId: x\n# the Button has its own flow (BUTTN-VR)\n- takeScreenshot: ALRTS-VR-S01\n")
	if got := scan.CaptureCodesIn(".maestro/ALRTS-VR.yaml", flow); strings.Join(got, ",") != "ALRTS-VR-S01" {
		t.Errorf("only what the flow compares against: %v", got)
	}
}

func TestEvidenceStaleFor_anotherTestsBaseline(t *testing.T) {
	t.Run("VRCPT-B10: A test captures the states it names, and no code a comment mentions", func(t *testing.T) {})
	g := Build(captureFiles(), &config.Config{}, nil)
	flow := g.node(".maestro/BUTTN-VR-S01.yaml")
	// Stamped when the flow captured every Button image.
	flow.Signal = &TestSignal{AtRev: flow.Rev, ClosureRev: map[string]string{
		"ui/Button.tsx": "c1", "ui/Button.BUTTN-VR-S01.png": "p1", "ui/Button.BUTTN-VR-S02-dark.svg": "p2"}}
	g.node("ui/Button.BUTTN-VR-S02-dark.svg").Rev = "p3"
	if ev := g.EvidenceStaleFor(flow.ID); ev != nil {
		t.Errorf("another state's baseline stales nothing: %+v", ev)
	}
	g.node("ui/Button.BUTTN-VR-S01.png").Rev = "p4"
	if ev := g.EvidenceStaleFor(flow.ID); ev == nil || strings.Join(ev.Culprit, ",") != "ui/Button.BUTTN-VR-S01.png" {
		t.Errorf("its own baseline does: %+v", ev)
	}
}

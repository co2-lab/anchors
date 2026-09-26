package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/spf13/cobra"
)

// qProject writes a throwaway project: anchors.yaml (when yaml is not empty), the given
// files, and the map (when g is not nil). Returns the root.
func qProject(t *testing.T, yaml string, files map[string]string, g *mapx.Graph) string {
	t.Helper()
	dir := t.TempDir()
	if yaml != "" {
		if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for name, body := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if g != nil {
		if err := mapx.Save(g, filepath.Join(dir, mapx.DefaultPath)); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// runQ executes a command built by a constructor and returns what it printed to stdout.
func runQ(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	return out, err
}

func readQ(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

const reportYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
  spec:
    kind: spec
    pattern: "*.spec.md"
governs:
  - from: GUIDE.md
    governs: code
gates:
  - name: always-fails
    on: [code]
    blocking: true
    run: "echo the tool said no; exit 1"
  - name: needs-a-judge
    on: [code]
    measures: judgment
`

// reportGraph carries every signal the test report reads: execution by layer (one
// stale), a measured spec with a gap, an unmeasured spec, and a code file below the
// line threshold that also regressed.
func reportGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "GUIDE.md", Kind: mapx.KindGuide},
			{ID: "LONE.md", Kind: mapx.KindGuide},
			{ID: "x.go", Kind: mapx.KindCode, Rev: "r1", Signal: &mapx.TestSignal{
				TotalLines: 10, CoveredLines: 4, LineCoverage: 40, PrevLineCoverage: 60,
			}},
			{ID: "y.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 90}},
			{ID: "x_test.go", Kind: mapx.KindTest, Rev: "new", Signal: &mapx.TestSignal{
				AtRev: "old",
				ByLayer: map[string]mapx.LayerExec{
					"unit": {Passed: 5, Failed: 1},
					"e2e":  {Passed: 2, Skipped: 3},
				},
			}},
			{ID: "b.spec.md", Kind: mapx.KindSpec, Code: "BBBBB", Signal: &mapx.TestSignal{
				ProvenCodes: []string{"BBBBB-B01"},
			}},
			{ID: "c.spec.md", Kind: mapx.KindSpec, Code: "CCCCC"},
			{ID: "nocode.spec.md", Kind: mapx.KindSpec},
		},
		Edges: []mapx.Edge{
			{From: "GUIDE.md", To: "x.go", Type: mapx.EdgeGoverns},
			{From: "GUIDE.md", To: "y.go", Type: mapx.EdgeGoverns},
		},
	}
}

func reportFiles() map[string]string {
	return map[string]string{
		"GUIDE.md":       "# guide\n",
		"LONE.md":        "# lone\n",
		"x.go":           "package x\n",
		"y.go":           "package x\n",
		"x_test.go":      "package x\n",
		"b.spec.md":      "# B\n\n## BBBBB-B01 — one\n\n## BBBBB-B02 — two\n",
		"c.spec.md":      "# C\n\n## CCCCC-B01 — one\n",
		"nocode.spec.md": "# no code\n",
		// issues: one waits on the user, one is agent work in progress, one is deferred
		"issues/todo/decide-pricing.md": "# decide pricing\n\n- **owner:** usuario\n",
		"issues/doing/fix-x.md":         "# fix x\n",
		"issues/future/later.md":        "# later\n",
	}
}

func TestRenderTestsMergesLayersAndSeparatesMeasured(t *testing.T) {
	dir := qProject(t, "", reportFiles(), nil)
	out := renderTests(reportCtx{g: reportGraph(), root: dir, when: "2026-01-02 03:04"})

	for _, want := range []string{
		"# Test report\n\n> Generated on 2026-01-02 03:04",
		"| e2e | 2 | 0 | 3 |\n| unit | 5 | 1 | 0 |\n",
		"| **total** | **7** | **1** | **3** |",
		"⚠ **1 failing test(s)**",
		"⚠ 1 test file(s) with a STALE signal",
		// only the measured spec (b) enters the proven count; c is reported apart
		"**1/2 requirements proven** by a green test, across 1 measured spec(s). 1 spec(s) with a gap:",
		"- `b.spec.md` — 1/2 unproven: BBBBB-B02",
		"_1 spec(s) still NOT measured_",
		"1 file(s) below 70%:\n\n- `x.go` — 40%\n",
		"- `x.go` — 60% → 40% (-20)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "`y.go` — 90%") {
		t.Errorf("a file above the threshold is not listed:\n%s", out)
	}
}

func TestRenderTestsWithNothingIngested(t *testing.T) {
	dir := qProject(t, "", reportFiles(), nil)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "c.spec.md", Kind: mapx.KindSpec, Code: "CCCCC"},
		{ID: "x.go", Kind: mapx.KindCode},
	}}
	out := renderTests(reportCtx{g: g, root: dir})
	for _, want := range []string{
		"_No execution result ingested._",
		"_No spec measured yet_ — 1 spec(s) with scenarios await test ingestion.",
		"0 file(s) below 70%:\n\n- (none, or line coverage not ingested)",
		"- (no drop since the previous ingestion)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "STALE") || strings.Contains(out, "failing test") {
		t.Errorf("nothing ingested: no stale and no failure can be claimed:\n%s", out)
	}
}

func TestReportAllWritesEveryPerspectiveAndTheIndex(t *testing.T) {
	dir := qProject(t, reportYAML, reportFiles(), reportGraph())
	if _, err := queue.Enqueue(dir, queue.Task{ID: "001-code-x", Changed: "x.go", Kind: "code", SuggestedNext: "test"}); err != nil {
		t.Fatal(err)
	}

	out, err := runQ(t, newReportCmd(), "all", "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "6 report(s) + index generated in docs/anchors/") {
		t.Errorf("unexpected summary:\n%s", out)
	}
	docs := filepath.Join(dir, "docs", "anchors")
	index := readQ(t, filepath.Join(docs, "index.md"))
	for _, sp := range reportSpecs {
		if !strings.Contains(index, "- ["+sp.name+"]("+sp.file+")") {
			t.Errorf("index should link %s:\n%s", sp.name, index)
		}
		if _, err := os.Stat(filepath.Join(docs, sp.file)); err != nil {
			t.Errorf("%s not written: %v", sp.file, err)
		}
	}

	quality := readQ(t, filepath.Join(docs, "anchors-quality-report.md"))
	for _, want := range []string{
		"| always-fails | **blocking** | 0 | 2 | 0 |",
		"| needs-a-judge | informative | 0 | 0 | 2 |",
		"✗ **Barred** — ",
		"- [BLOCKS] `always-fails` @ `x.go` — the tool said no",
		"## Awaiting AI judgment (2)",
	} {
		if !strings.Contains(quality, want) {
			t.Errorf("quality report missing %q:\n%s", want, quality)
		}
	}

	structure := readQ(t, filepath.Join(docs, "anchors-structure-report.md"))
	for _, want := range []string{
		"| code | 2 |", "| spec | 3 |",
		"8 nodes, 2 edges.",
		"| `GUIDE.md` | 2 |",
		"## Missing identity (1)\n\n- `nocode.spec.md`",
	} {
		if !strings.Contains(structure, want) {
			t.Errorf("structure report missing %q:\n%s", want, structure)
		}
	}

	cfgRep := readQ(t, filepath.Join(docs, "anchors-config-report.md"))
	for _, want := range []string{
		"## Declared layers (2)",
		"| code | code |",
		"- `GUIDE.md` governs the tag `code`",
		"## Declared gates (2)",
		"- `always-fails` (external: echo the tool said no; exit 1, blocking) over [code]",
		"- `needs-a-judge` (AI-judgment, informative) over [code]",
		"## Derived (co-location): not configured",
		"- `LONE.md` —",
		"## Kinds with no gate (coverage hole)",
		"- `spec` —",
	} {
		if !strings.Contains(cfgRep, want) {
			t.Errorf("config report missing %q:\n%s", want, cfgRep)
		}
	}

	issues := readQ(t, filepath.Join(docs, "anchors-issues-report.md"))
	for _, want := range []string{
		"| future | 1 |\n| todo | 1 |\n| doing | 1 |\n| done | 0 |",
		"### Waiting on YOU (the agent cannot resolve it)",
		"- decide-pricing\n",
		"### Open (the work of now)\n\n- fix-x\n",
		"### Deferred (adopted debt — it comes due, it does not vanish)\n\n- later\n",
		"1 live task(s):",
		"`x.go` — test (code)",
	} {
		if !strings.Contains(issues, want) {
			t.Errorf("issues report missing %q:\n%s", want, issues)
		}
	}
	// the user's issue is in the user's list only
	if strings.Count(issues, "- decide-pricing") != 1 {
		t.Errorf("the issue that waits on the user must not be repeated as agent work:\n%s", issues)
	}

	inc := readQ(t, filepath.Join(docs, "anchors-inconsistencies-report.md"))
	for _, want := range []string{
		"### guide-sem-governo (1)",
		"- ⚠ `LONE.md`",
		"## From the quality gates (2 failures)",
		"- `always-fails` @ `x.go` — the tool said no",
	} {
		if !strings.Contains(inc, want) {
			t.Errorf("inconsistencies report missing %q:\n%s", want, inc)
		}
	}
}

func TestReportSubcommandWritesToOut(t *testing.T) {
	dir := qProject(t, reportYAML, reportFiles(), reportGraph())
	dest := filepath.Join(dir, "custom", "tests.md")

	out, err := runQ(t, newReportCmd(), "tests", "--root", dir, "--out", dest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "report generated: custom/tests.md") {
		t.Errorf("the path printed is relative to the root:\n%s", out)
	}
	if !strings.Contains(readQ(t, dest), "# Test report") {
		t.Error("the tests perspective was not written to --out")
	}

	// default destination: docs/<file>
	if _, err := runQ(t, newReportCmd(), "structure", "--root", dir); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readQ(t, filepath.Join(dir, "docs", "anchors-structure-report.md")), "# Structure report") {
		t.Error("the default destination is docs/<file>")
	}
}

func TestReportWithoutMapFails(t *testing.T) {
	dir := qProject(t, reportYAML, nil, nil)
	_, err := runQ(t, newReportCmd(), "tests", "--root", dir)
	if err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: must fail pointing at map build; got %v", err)
	}
	_, err = runQ(t, newReportCmd(), "all", "--root", dir)
	if err == nil || !strings.Contains(err.Error(), "load map") {
		t.Errorf("no map: `all` must fail too; got %v", err)
	}
}

// Without anchors.yaml the renderers that depend on it say so instead of inventing.
func TestRenderersWithoutConfig(t *testing.T) {
	ctx := reportCtx{g: &mapx.Graph{}, root: t.TempDir(), when: "now"}
	if out := renderQuality(ctx); !strings.Contains(out, "_No gate declared in anchors.yaml._") {
		t.Errorf("quality without config:\n%s", out)
	}
	if out := renderConfig(ctx); !strings.Contains(out, "_anchors.yaml not found._") {
		t.Errorf("config without anchors.yaml:\n%s", out)
	}
	out := renderIssues(ctx)
	if !strings.Contains(out, "| todo | 0 |") || !strings.Contains(out, "_Empty queue._") {
		t.Errorf("issues without anything:\n%s", out)
	}
	if strings.Contains(out, "Waiting on YOU") || strings.Contains(out, "Deferred") {
		t.Errorf("empty lists are not printed:\n%s", out)
	}
}

// findingSection cuts the health report by check prefix and caps a long list.
func TestFindingSectionCutsByPrefixAndCaps(t *testing.T) {
	var rep health.Report
	for i := 0; i < 27; i++ {
		rep.Findings = append(rep.Findings, health.Finding{Check: "orfao-identidade", Subject: fmt.Sprintf("n%02d", i), Detail: "d"})
	}
	rep.Findings = append(rep.Findings, health.Finding{Check: "other", Subject: "zz", Detail: "d"})

	out := findingSection(rep, "orfao", "Orphans")
	if !strings.HasPrefix(out, "## Orphans (27)\n") {
		t.Errorf("the title counts only the matching findings:\n%s", out)
	}
	if !strings.Contains(out, "- `n24` — d\n- … and 2 more\n") || strings.Contains(out, "n25") {
		t.Errorf("the list stops at 25 and says how many are left:\n%s", out)
	}
	if strings.Contains(out, "zz") {
		t.Errorf("another check does not enter the section:\n%s", out)
	}
	if findingSection(rep, "absent", "X") != "" {
		t.Error("no finding of the check: no section at all")
	}
}

func TestReportHelpers(t *testing.T) {
	if got := firstLineOf("one\ntwo"); got != "one" {
		t.Errorf("firstLineOf = %q", got)
	}
	if !contémNome([]string{"a.md", "b.md"}, "b.md") || contémNome([]string{"a.md"}, "b.md") {
		t.Error("contémNome must match the exact name only")
	}
}

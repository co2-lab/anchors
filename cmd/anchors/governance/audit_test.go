package governance

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/mapx"
)

const auditYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
  spec:
    kind: spec
    pattern: "*.spec.md"
gates:
  - name: always-fails
    on: [code]
    run: "echo broken line one; echo second line; exit 1"
  - name: needs-a-judge
    on: [code]
    measures: judgment
`

func auditGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "a.go", Kind: mapx.KindCode, Layer: "code"},
			{ID: "a.spec.md", Kind: mapx.KindSpec, Layer: "spec"},
			{ID: "other.go", Kind: mapx.KindCode, Layer: "code"},
		},
		Edges: []mapx.Edge{
			{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies},
		},
	}
}

func auditFiles() map[string]string {
	return map[string]string{"a.go": "package a\n", "a.spec.md": "# A\n", "other.go": "package a\n"}
}

// Without --impact the dossier is the file alone: the gates that apply to its kind,
// grouped under the target, and nothing about the other nodes.
func TestAuditFileOnly(t *testing.T) {
	dir := govProject(t, auditYAML, auditFiles(), auditGraph())

	out, err := runCmd(t, newAuditCmd(), "--root", dir, "a.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "audit: a.go — the file\n") {
		t.Errorf("header should name the target and the file-only scope:\n%s", out)
	}
	if !strings.Contains(out, "● a.go\n") {
		t.Errorf("the target should head its own block:\n%s", out)
	}
	if !strings.Contains(out, "  ✗ [gate] always-fails — broken line one\n") {
		t.Errorf("a failing gate is listed with its FIRST detail line only:\n%s", out)
	}
	if strings.Contains(out, "second line") {
		t.Errorf("only the first line of the detail belongs in the dossier:\n%s", out)
	}
	if !strings.Contains(out, "⏳ [gate] needs-a-judge") {
		t.Errorf("a judgment gate is listed as awaiting judgment:\n%s", out)
	}
	if strings.Contains(out, "other.go") || strings.Contains(out, "(impact)") {
		t.Errorf("without --impact no other node belongs in the dossier:\n%s", out)
	}
	// only the Fail counts as actionable; the pending judgment does not
	if !strings.Contains(out, "1 actionable pending item(s)") {
		t.Errorf("one failing gate = one actionable item:\n%s", out)
	}
}

// With --impact the unit's other nodes join the scope and print after the target.
func TestAuditWithImpactIncludesTheUnit(t *testing.T) {
	dir := govProject(t, auditYAML, auditFiles(), auditGraph())

	out, err := runCmd(t, newAuditCmd(), "--root", dir, "--impact", "a.spec.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "the unit (") {
		t.Errorf("--impact should announce the unit scope:\n%s", out)
	}
	if !strings.Contains(out, "○ a.go (impact)") {
		t.Errorf("the code the spec propagates to should appear as an impact node:\n%s", out)
	}
	if strings.Contains(out, "other.go") {
		t.Errorf("a node outside the impact path must stay out:\n%s", out)
	}
}

func TestAuditRefusesAFileOutsideTheMap(t *testing.T) {
	files := auditFiles()
	files["new.go"] = "package a\n"
	dir := govProject(t, auditYAML, files, auditGraph())

	_, err := runCmd(t, newAuditCmd(), "--root", dir, "new.go")
	if err == nil || !strings.Contains(err.Error(), `"new.go" is not in the map`) {
		t.Errorf("a file the map does not know cannot be audited; got %v", err)
	}
}

func TestAuditWithoutConfigOrMapFails(t *testing.T) {
	dir := t.TempDir()
	if _, err := runCmd(t, newAuditCmd(), "--root", dir, "a.go"); err == nil || !strings.Contains(err.Error(), "load config") {
		t.Errorf("no anchors.yaml must fail loading the config; got %v", err)
	}
	dir = govProject(t, auditYAML, auditFiles(), nil)
	if _, err := runCmd(t, newAuditCmd(), "--root", dir, "a.go"); err == nil || !strings.Contains(err.Error(), "load map") {
		t.Errorf("no map must fail loading the map; got %v", err)
	}
}

// printAudit is the dossier's shape: pass/skip vanish, a warn doctor finding counts, an
// info one does not, and a finding about a node outside the scope is dropped.
func TestPrintAuditCountsOnlyActionableItems(t *testing.T) {
	ids := map[string]bool{"a.go": true, "a.spec.md": true}
	results := []gate.Result{
		{Gate: "ok", Target: "a.go", Verdict: gate.Pass},
		{Gate: "skipped", Target: "a.go", Verdict: gate.Skip},
		{Gate: "diverges", Target: "a.spec.md", Verdict: gate.Pending, Detail: "one\ntwo"},
	}
	rep := health.Report{Findings: []health.Finding{
		{Check: "orphan", Subject: "a.go", Severity: health.Warn, Detail: "no edges"},
		{Check: "hint", Subject: "a.go", Severity: health.Info, Detail: "just a note"},
		{Check: "orphan", Subject: "elsewhere.go", Severity: health.Warn, Detail: "out of scope"},
	}}

	out := captureStdout(t, func() {
		if err := printAudit("a.go", true, []mapx.Node{{ID: "a.go"}, {ID: "a.spec.md"}}, results, rep, ids); err != nil {
			t.Error(err)
		}
	})

	for _, want := range []string{
		"the unit (2 nodes on the impact path)",
		"  ⚠ [doctor:orphan] no edges",
		"  ℹ [doctor:hint] just a note",
		"○ a.spec.md (impact)",
		"  ~ [gate] diverges — one\n",
		"1 actionable pending item(s)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	for _, not := range []string{"[gate] ok", "[gate] skipped", "out of scope"} {
		if strings.Contains(out, not) {
			t.Errorf("%q must not be in the dossier:\n%s", not, out)
		}
	}
	if strings.Index(out, "● a.go") > strings.Index(out, "○ a.spec.md") {
		t.Errorf("the target prints before the impact nodes:\n%s", out)
	}
}

func TestPrintAuditNothingPending(t *testing.T) {
	out := captureStdout(t, func() {
		_ = printAudit("a.go", false, nil, []gate.Result{{Gate: "ok", Target: "a.go", Verdict: gate.Pass}}, health.Report{}, map[string]bool{"a.go": true})
	})
	if !strings.Contains(out, "✓ nothing pending") {
		t.Errorf("a clean file says so:\n%s", out)
	}
}

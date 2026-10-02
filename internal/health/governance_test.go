package health

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func hasOpportunity(fs []Finding, check, subject string) bool {
	for _, f := range fs {
		if f.Check == check && f.Subject == subject {
			return true
		}
	}
	return false
}

func TestCheckGovernanceOpportunities_EvidenceFreshSuggestion(t *testing.T) {
	t.Run("GVOPG-B01: Tests and code without evidence freshness get the suggestion", func(t *testing.T) {})
	// A project with tests and code, but no evidence-fresh
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "DataTable.tsx", Kind: mapx.KindCode},
			{ID: "DataTable.test.tsx", Kind: mapx.KindTest},
		},
	}
	cfg := &config.Config{
		Gates: []config.Gate{
			{Name: "tests-green", Check: "tests-pass"},
		},
	}

	found := false
	for _, f := range checkGovernanceOpportunities(g, cfg) {
		if f.Check == "sugestao-gate" && f.Subject == "evidence-fresh" {
			found = true
			if f.Severity != Info {
				t.Errorf("severity should be Info, got %v", f.Severity)
			}
		}
	}
	if !found {
		t.Error("expected the evidence-fresh suggestion for a project with tests and code")
	}
}

func TestCheckGovernanceOpportunities_NoSecretLeakedNotBlocking(t *testing.T) {
	t.Run("GVOPG-B03: A secret gate that does not block is suboptimal", func(t *testing.T) {})
	// no-secret-leaked declared with blocking: false (suboptimal)
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "main.go", Kind: mapx.KindCode},
		},
	}
	blockingFalse := false
	cfg := &config.Config{
		Gates: []config.Gate{
			{Name: "no-secret-leaked", Blocking: &blockingFalse},
		},
	}

	fs := checkGovernanceOpportunities(g, cfg)
	if !hasOpportunity(fs, "gate-subotimo", "no-secret-leaked") {
		t.Error("expected the suboptimal-gate finding for a non-blocking no-secret-leaked")
	}
	if hasOpportunity(fs, "sugestao-gate", "no-secret-leaked") {
		t.Error("a declared gate is not suggested again")
	}
}

func TestCheckGovernanceOpportunities_TestsWithoutJunit(t *testing.T) {
	t.Run("GVOPG-B05: Tests without JUnit output are a suboptimal configuration", func(t *testing.T) {})
	// A project with tests but no junit configured
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "auth.test.ts", Kind: mapx.KindTest},
		},
	}
	cfg := &config.Config{
		Tests: []config.Suite{
			{Layer: "unit", Run: "pnpm test"},
		},
	}

	if !hasOpportunity(checkGovernanceOpportunities(g, cfg), "config-subotima", "tests.junit") {
		t.Error("expected the suboptimal-configuration finding for tests without junit")
	}
}

// Code alone is offered the three code gates, each once.
func TestCheckGovernanceOpportunities_CodeGates(t *testing.T) {
	t.Run("GVOPG-B02: Code without the secret gate gets the suggestion", func(t *testing.T) {})
	t.Run("GVOPG-B04: Code without dependency audit or duplication gates gets both suggestions", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "main.go", Kind: mapx.KindCode}}}
	fs := checkGovernanceOpportunities(g, &config.Config{})
	for _, gate := range []string{"no-secret-leaked", "dependency-vulnerable", "no-duplication"} {
		if !hasOpportunity(fs, "sugestao-gate", gate) {
			t.Errorf("expected a suggestion of %s for a project with code: %+v", gate, fs)
		}
	}
	// Code alone has no tests: neither evidence-fresh nor the junit setting apply.
	if len(fs) != 3 {
		t.Errorf("code alone gives exactly the three code suggestions, got %+v", fs)
	}
}

func TestCheckGovernanceOpportunities_NilInputs(t *testing.T) {
	t.Run("GVOPG-B06: A nil map or configuration gives nothing", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "main.go", Kind: mapx.KindCode}}}
	if fs := checkGovernanceOpportunities(nil, &config.Config{}); fs != nil {
		t.Errorf("a nil map gives nothing: %+v", fs)
	}
	if fs := checkGovernanceOpportunities(g, nil); fs != nil {
		t.Errorf("a nil configuration gives nothing: %+v", fs)
	}
}

func everyOpportunity() (*mapx.Graph, *config.Config) {
	return &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.go", Kind: mapx.KindCode},
		{ID: "a_test.go", Kind: mapx.KindTest},
	}}, &config.Config{}
}

func TestQuickGovernanceHints_firstTwo(t *testing.T) {
	t.Run("GVOPG-B07: The quick hints are the first two opportunities", func(t *testing.T) {})
	g, cfg := everyOpportunity()
	all := checkGovernanceOpportunities(g, cfg)
	if len(all) != 5 {
		t.Fatalf("setup: five opportunities expected, got %+v", all)
	}
	hints := QuickGovernanceHints(g, cfg)
	if len(hints) != 2 || hints[0] != all[0] || hints[1] != all[1] {
		t.Fatalf("the hints must be the first two of %+v, got %+v", all, hints)
	}
	// Fewer than two: returned whole.
	blocking := true
	onlyOne := &config.Config{Gates: []config.Gate{
		{Name: "no-secret-leaked", Blocking: &blocking}, {Name: "dependency-vulnerable"},
	}}
	code := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}}
	if hints := QuickGovernanceHints(code, onlyOne); len(hints) != 1 || hints[0].Subject != "no-duplication" {
		t.Fatalf("a single opportunity is returned whole, got %+v", hints)
	}
}

func TestCheckGovernanceOpportunities_allInfo(t *testing.T) {
	t.Run("GVOPG-I01: Every opportunity is informational", func(t *testing.T) {})
	g, _ := everyOpportunity()
	notBlocking := false
	cfg := &config.Config{Gates: []config.Gate{{Name: "no-secret-leaked", Blocking: &notBlocking}}}
	fs := checkGovernanceOpportunities(g, cfg)
	if len(fs) != 5 {
		t.Fatalf("setup: five opportunities expected, got %+v", fs)
	}
	for _, f := range fs {
		if f.Severity != Info {
			t.Errorf("an opportunity is never a warning: %+v", f)
		}
	}
}

func TestCheckGovernanceOpportunities_adoptedIsNotSuggested(t *testing.T) {
	t.Run("GVOPG-X01: What the project already adopted is not suggested", func(t *testing.T) {})
	g, _ := everyOpportunity()
	blocking := true
	cfg := &config.Config{
		Gates: []config.Gate{
			{Name: "evidence-fresh"}, {Name: "no-secret-leaked", Blocking: &blocking},
			{Name: "dependency-vulnerable"}, {Name: "no-duplication"},
		},
		Tests: []config.Suite{{Layer: "unit", Run: "go test", JUnit: "junit.xml"}},
	}
	if fs := checkGovernanceOpportunities(g, cfg); len(fs) != 0 {
		t.Fatalf("a project that adopted everything gets no suggestion: %+v", fs)
	}
}

func TestCheckGovernanceOpportunities_applicableCatalogGates(t *testing.T) {
	t.Run("GVOPG-B08: Every applicable undeclared catalog gate is suggested", func(t *testing.T) {})
	prev := config.SetGateCatalog(func() []config.Gate {
		return []config.Gate{
			{Name: "feature-test-match", On: []string{"feature"}, Measures: "every scenario is implemented"},
			{Name: "scenario-letter-declared", On: []string{"feature"}, Presupposes: []string{"rule_types"}, Measures: "letters are declared"},
		}
	})
	t.Cleanup(func() { config.SetGateCatalog(prev) })
	cfg := &config.Config{Layers: map[string]config.Layer{"feature": {Kind: "feature"}}}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.feature", Kind: mapx.KindFeature}}}
	var ftm, sld *Finding
	fs := checkGovernanceOpportunities(g, cfg)
	for i := range fs {
		switch fs[i].Subject {
		case "feature-test-match":
			ftm = &fs[i]
		case "scenario-letter-declared":
			sld = &fs[i]
		}
	}
	if ftm == nil || ftm.Severity != Info || !strings.Contains(ftm.Detail, "- name: feature-test-match") {
		t.Errorf("feature-test-match is suggested with how to declare it: %+v", ftm)
	}
	if sld == nil || !strings.Contains(sld.Detail, "rule_types") {
		t.Errorf("scenario-letter-declared's suggestion names rule_types: %+v", sld)
	}
}

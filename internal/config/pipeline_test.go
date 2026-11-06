// @anchors
//   code: PCTPP
//   ref: PCSPP

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadPipeline(t *testing.T, body string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 8\npipeline:\n"+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

const mifPipeline = `  unguarded:
    build-android.yml: only builds an artifact
  stages:
    review-release:
      branches: [review]
      requires: [flow, { gates: ci }]
      verify: node scripts/smoke-review.mjs
    production-release:
      branches: [production]
      requires:
        - flow
        - any:
            - passed: review-release
            - from: "hotfix/*"
      verify: node scripts/smoke-prod.mjs
`

func TestPipeline_stagesAndTheirRequirements(t *testing.T) {
	t.Run("PCSPP-B01: A stage declares its branches, its requirements and the command that runs or validates it", func(t *testing.T) {})
	c, err := loadPipeline(t, mifPipeline)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(c.Pipeline.StageNames(), ","); got != "production-release,review-release" {
		t.Errorf("the stages, sorted: %s", got)
	}
	review := c.Pipeline.Stages["review-release"]
	if len(review.Requires) != 2 || !review.Requires[0].Flow || review.Requires[1].Gates != "ci" {
		t.Errorf("review requires the flow and the ci gates: %+v", review.Requires)
	}
	prod := c.Pipeline.Stages["production-release"].Requires
	if len(prod) != 2 || len(prod[1].Any) != 2 || prod[1].Any[0].Passed != "review-release" || prod[1].Any[1].From != "hotfix/*" {
		t.Errorf("production requires the flow and one of two: %+v", prod)
	}
	if prod[1].String() != "any: [passed: review-release, from: hotfix/*]" {
		t.Errorf("a requirement is written as declared: %s", prod[1])
	}
	if c.Pipeline.Unguarded["build-android.yml"] == "" {
		t.Error("an unguarded workflow keeps its reason")
	}
	if g := (*Pipeline)(nil).WorkflowGlobs(); len(g) != 2 || g[0] != ".github/workflows/*.yml" {
		t.Errorf("the workflows default to GitHub's: %v", g)
	}
}

func TestPipeline_aStageThatGuardsNothingDoesNotLoad(t *testing.T) {
	t.Run("PCSPP-E01: A stage with nothing that runs or validates it, a passed: naming no stage, or an unknown requirement refuses the configuration", func(t *testing.T) {})
	for bad, want := range map[string]string{
		"  stages:\n    a:\n      branches: [x]\n":                                           "run:",
		"  stages:\n    a:\n      verify: x\n      requires: [{ passed: b }]\n":              "no stage",
		"  stages:\n    a:\n      verify: x\n      requires: [review]\n":                     "requirement",
		"  stages:\n    a:\n      verify: x\n      requires: [{ when: x }]\n":                "requirement",
		"  stages:\n    a:\n      verify: x\n      branches: [\"re:(\"]\n":                   "regular expression",
		"  stages:\n    a:\n      verify: x\n      requires: [{ any: [{ from: \"[\" }] }]\n": "branch pattern",
	} {
		if _, err := loadPipeline(t, bad); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q refuses the configuration naming %q: %v", bad, want, err)
		}
	}
}

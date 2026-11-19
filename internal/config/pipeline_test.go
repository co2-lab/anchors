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

const releasePipeline = `  unguarded:
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
	c, err := loadPipeline(t, releasePipeline)
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

const versionedPipeline = `versioning:
  app:
    scheme: semver
    source: [{ file: app/version.json, path: version }]
    released: { file: policy.json, keys: versions, where: { target: production } }
pipeline:
  stages:
    review-release:
      branches: [review]
      requires:
        - contains: production
        - version: { product: app, fresh: true, listed: { where: { target: review } }, expected: true }
      verify: x
    production-release:
      branches: [production]
      requires:
        - version: { same_as_passed: review-release }
      verify: x
`

func loadRaw(t *testing.T, body string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 8\n"+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestPipeline_versionAndContains(t *testing.T) {
	t.Run("PCSPP-B02: A stage asks of the version its sha carries, and of a branch whose head the sha contains", func(t *testing.T) {
		c, err := loadRaw(t, versionedPipeline)
		if err != nil {
			t.Fatal(err)
		}
		rev := c.Pipeline.Stages["review-release"].Requires
		if rev[0].Contains != "production" || rev[0].String() != "contains: production" {
			t.Errorf("contains: %+v %s", rev[0], rev[0])
		}
		v := rev[1].Version
		if v == nil || v.Product != "app" || !v.Fresh || !v.Expected || v.Listed == nil || v.Listed.Where["target"] != "review" || v.Listed.File != "" {
			t.Fatalf("version: %+v", v)
		}
		if got := rev[1].String(); got != "version: app fresh, listed, expected" {
			t.Errorf("string: %q", got)
		}
		if got := c.Pipeline.Stages["production-release"].Requires[0].String(); got != "version: same_as_passed review-release" {
			t.Errorf("a single product needs no name: %q", got)
		}
	})
	t.Run("PCSPP-E02: A version asked of no product, with nothing asked, of a stage or a listing that do not exist, or a contains: pattern refuses the configuration", func(t *testing.T) {
		stage := func(req string) string {
			return strings.Replace(versionedPipeline, "        - version: { same_as_passed: review-release }", "        - "+req, 1)
		}
		for bad, want := range map[string]string{
			stage("version: { product: web, fresh: true }"):             `no product "web"`,
			stage("version: { product: app }"):                          "ask fresh, listed, same_as_passed or expected",
			stage("version: { same_as_passed: qa }"):                    "qa is no stage",
			stage("version: { fresh: true, latest: true }"):             "latest",
			stage("version: { listed: { where: { a: b }, match: x } }"): "match",
			stage("contains: \"hotfix/*\""):                             "name one branch",
			strings.Replace(stage("version: { listed: { where: { a: b } } }"), "    released: { file: policy.json, keys: versions, where: { target: production } }\n", "", 1): "declare versioning.app.released",
			"pipeline:\n  stages:\n    a:\n      verify: x\n      requires: [{ version: { fresh: true } }]\n":                                                                 "no `versioning:`",
		} {
			if _, err := loadRaw(t, bad); err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("refuses naming %q: %v\n%s", want, err, bad)
			}
		}
		two := strings.Replace(versionedPipeline, "pipeline:", "  api: { scheme: number, next: none }\npipeline:", 1)
		if _, err := loadRaw(t, two); err == nil || !strings.Contains(err.Error(), "name one with product:") {
			t.Errorf("several products need a name: %v", err)
		}
	})
}

func TestPipeline_release(t *testing.T) {
	t.Run("PCSPP-B03: A stage names the product it releases", func(t *testing.T) {
		c, err := loadRaw(t, strings.Replace(versionedPipeline, "      branches: [production]\n", "      branches: [production]\n      release: app\n", 1))
		if err != nil || c.Pipeline.Stages["production-release"].Release != "app" {
			t.Errorf("%v", err)
		}
	})
	t.Run("PCSPP-E03: A stage that releases a product versioning: does not declare refuses the configuration", func(t *testing.T) {
		_, err := loadRaw(t, strings.Replace(versionedPipeline, "      branches: [production]\n", "      branches: [production]\n      release: web\n", 1))
		if err == nil || !strings.Contains(err.Error(), "production-release: release:") || !strings.Contains(err.Error(), `no product "web"`) {
			t.Errorf("%v", err)
		}
	})
}

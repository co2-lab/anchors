// @anchors
//   code: PWTPP
//   ref: PWSPP

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func wiredCfg() *config.Config {
	return &config.Config{Pipeline: &config.Pipeline{
		Unguarded: map[string]string{"build-android.yml": "only builds an artifact"},
		Stages: map[string]config.Stage{
			"review-release":     {Verify: "x"},
			"production-release": {Verify: "x"},
		},
	}}
}

const reviewWorkflow = "jobs:\n  guard:\n    steps:\n      - run: anchors pipeline can-run review-release\n  done:\n    steps:\n      - run: anchors pipeline done review-release\n"
const prodWorkflow = "jobs:\n  guard:\n    steps:\n      - run: anchors pipeline can-run production-release\n  done:\n    steps:\n      - run: anchors pipeline done production-release\n"

func TestPipelineWired_everyStageCalledAndEveryCallDeclared(t *testing.T) {
	t.Run("PWSPP-B01: The gate passes when every call names a declared stage, every stage is called and recorded, and every workflow calls one or is unguarded", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		".github/workflows/review-release.yml": reviewWorkflow,
		".github/workflows/mobile-release.yml": prodWorkflow,
		".github/workflows/build-android.yml":  "jobs: {}\n",
	})
	if v, d := checkPipelineWired(config.Gate{}, root, nil, wiredCfg()); v != Pass || !strings.Contains(d, "3 workflow(s), 2 stage(s)") {
		t.Errorf("a wired pipeline passes: %v %s", v, d)
	}
}

func TestPipelineWired_whatIsNotWiredFails(t *testing.T) {
	t.Run("PWSPP-B02: A call to an undeclared stage, a stage no workflow calls, a stage asked and never recorded, and a workflow calling none fail the gate", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		".github/workflows/review-release.yml": strings.Replace(reviewWorkflow, "done review-release", "done deploy", 1),
		".github/workflows/admin.yaml":         "jobs: {}\n",
	})
	v, d := checkPipelineWired(config.Gate{}, root, nil, wiredCfg())
	if v != Fail {
		t.Fatalf("a pipeline not wired fails: %v %s", v, d)
	}
	for _, want := range []string{
		"review-release.yml calls the stage deploy",
		"production-release is declared and no workflow calls it",
		"review-release is asked with can-run and never recorded",
		"admin.yaml calls no stage",
	} {
		if !strings.Contains(d, want) {
			t.Errorf("the failure names %q:\n%s", want, d)
		}
	}
}

func TestPipelineWired_nothingToWire(t *testing.T) {
	t.Run("PWSPP-B03: With no stage declared, or no workflow found, the gate is pending", func(t *testing.T) {})
	if v, _ := checkPipelineWired(config.Gate{}, t.TempDir(), nil, &config.Config{}); v != Pending {
		t.Errorf("no pipeline: pending, got %v", v)
	}
	if v, d := checkPipelineWired(config.Gate{}, t.TempDir(), nil, wiredCfg()); v != Pending || !strings.Contains(d, ".github/workflows/*.yml") {
		t.Errorf("no workflow: pending, naming the globs, got %v %s", v, d)
	}
}

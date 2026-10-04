// @anchors
//   code: CNTSB
//   ref: CTTST

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const apiUnitSpec = "<!-- @anchors\n  code: QRGEN\n-->\n# Generate\n\n## Endpoint\n| Method | Path | Operation | Deprecated |\n| --- | --- | --- | --- |\n| POST | /v1/qrcodes | generateQrCode | no |\n"

const contractFeature = "Feature: Generate\n\n  @QRGEN-CT @contract-level\n  Scenario: The API keeps its OpenAPI contract\n    Given the document\n"

func apiFixture(t *testing.T, files map[string]string, tests ...string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if _, ok := files["api/generate.spec.md"]; !ok {
		files["api/generate.spec.md"] = apiUnitSpec
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

var apiCode = mapx.Node{ID: "api/generate.go", Kind: mapx.KindCode}

func TestContractTested_noAPIUnit(t *testing.T) {
	t.Run("CTTST-B01: What is no API unit leaves without a verdict", func(t *testing.T) {})
	t.Run("CTTST-E01: A code file with no spec beside it leaves without a verdict", func(t *testing.T) {})
	root, g := apiFixture(t, map[string]string{"api/model.spec.md": "<!-- @anchors\n  code: MODEL\n-->\n# Model\n"})
	if v, _ := checkContractTested("", mapx.Node{ID: "api/generate.spec.md", Kind: mapx.KindSpec}, root, g, nil); v != Skip {
		t.Errorf("a spec node: %v", v)
	}
	if v, d := checkContractTested("", mapx.Node{ID: "api/helpers.go", Kind: mapx.KindCode}, root, g, nil); v != Skip || !strings.Contains(d, "spec") {
		t.Errorf("no spec beside it: %v %s", v, d)
	}
	if v, d := checkContractTested("", mapx.Node{ID: "api/model.go", Kind: mapx.KindCode}, root, g, nil); v != Skip || !strings.Contains(d, "Endpoint") {
		t.Errorf("a spec with no Endpoint: %v %s", v, d)
	}
}

func TestContractTested_theThreeParts(t *testing.T) {
	t.Run("CTTST-B02: The feature must carry the contract scenario", func(t *testing.T) {})
	t.Run("CTTST-B03: A test of the unit must name the contract code", func(t *testing.T) {})
	t.Run("CTTST-B04: The contract test must load the OpenAPI document", func(t *testing.T) {})
	t.Run("CTTST-B06: Scenario and a contract test that loads the document pass", func(t *testing.T) {})
	root, g := apiFixture(t, map[string]string{})
	if v, d := checkContractTested("", apiCode, root, g, nil); v != Fail || !strings.Contains(d, "`QRGEN-CT` (tagged `@contract-level`)") ||
		!strings.Contains(d, "no test of the unit names `QRGEN-CT`") {
		t.Errorf("no scenario, no test: %v %s", v, d)
	}
	root, g = apiFixture(t, map[string]string{"api/generate.feature": contractFeature,
		"api/generate_contract_test.go": "// QRGEN-CT\nfunc TestContract(t *testing.T) {}"}, "api/generate_contract_test.go")
	if v, d := checkContractTested("", apiCode, root, g, nil); v != Fail || !strings.Contains(d, "does not load the OpenAPI document") {
		t.Errorf("a test that does not load the document: %v %s", v, d)
	}
	root, g = apiFixture(t, map[string]string{"api/generate.feature": contractFeature,
		"api/generate_contract_test.go": "// QRGEN-CT\nvar doc = load(\"../docs/openapi.yaml\")"}, "api/generate_contract_test.go")
	if v, d := checkContractTested("", apiCode, root, g, nil); v != Pass {
		t.Errorf("the three parts pass: %v %s", v, d)
	}
}

func TestContractTested_theRegimeTag(t *testing.T) {
	t.Run("CTTST-B05: The contract regime tag comes from the project", func(t *testing.T) {})
	if tag := contractRegimeTag(nil); tag != "contract-level" {
		t.Errorf("default = %q", tag)
	}
	cfg := &config.Config{Derived: &config.Derived{Regimes: map[string]string{"nivel-unit": "unit", "nivel-contrato": "contract"}}}
	if tag := contractRegimeTag(cfg); tag != "nivel-contrato" {
		t.Errorf("from the project = %q", tag)
	}
}

func TestContractTested_anUnreadTestIsReadByItsPath(t *testing.T) {
	t.Run("CTTST-E02: A test that cannot be read is read by its path", func(t *testing.T) {})
	root, g := apiFixture(t, map[string]string{"api/generate.feature": contractFeature}, "flows/QRGEN-CT-openapi.yaml")
	if v, d := checkContractTested("", apiCode, root, g, nil); v != Pass {
		t.Errorf("a flow named by the code and the document counts: %v %s", v, d)
	}
}

// @anchors
//   code: ASTPS
//   ref: APISP

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const apiHead = "<!-- @anchors\n  code: QRGEN\n-->\n# Generate\n\n## Endpoint\n| Method | Path | Operation | Deprecated |\n| --- | --- | --- | --- |\n| POST | /v1/qrcodes | generateQrCode | no |\n\n"

// apiProject writes the files, and a map whose specs carry their header code and whose
// unit spec specifies the given code files.
func apiProject(t *testing.T, files map[string]string, specifies ...string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	g := &mapx.Graph{}
	for rel, body := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		must(t, os.MkdirAll(filepath.Dir(p), 0o755))
		must(t, os.WriteFile(p, []byte(body), 0o644))
		if strings.HasSuffix(rel, ".spec.md") {
			code := ""
			if m := specCodeRE().FindStringSubmatch(body); m != nil {
				code = m[1]
			}
			g.Nodes = append(g.Nodes, mapx.Node{ID: rel, Kind: mapx.KindSpec, Code: code})
		}
	}
	for _, s := range specifies {
		g.Edges = append(g.Edges, mapx.Edge{From: "api/generate.spec.md", To: s, Type: mapx.EdgeSpecifies})
	}
	return root, g
}

var apiChecks = map[string]func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string){
	"api-contracts-resolve": checkAPIContractsResolve, "api-errors-declared": checkAPIErrorsDeclared,
	"error-codes-honored": checkErrorCodesHonored,
}

func TestAPISpec_noAPIUnit(t *testing.T) {
	t.Run("APISP-B01: What is no API unit leaves every gate without a verdict", func(t *testing.T) {})
	t.Run("APISP-E01: A code file with no spec beside it leaves without a verdict", func(t *testing.T) {})
	root, g := apiProject(t, map[string]string{"api/model.spec.md": "<!-- @anchors\n  code: MODEL\n-->\n# Model\n"})
	for name, check := range apiChecks {
		if v, _ := check("", mapx.Node{ID: "api/model.spec.md", Kind: mapx.KindSpec}, root, g, nil); v != Skip {
			t.Errorf("%s on a spec node: %v", name, v)
		}
		if v, _ := check("", mapx.Node{ID: "api/model.go", Kind: mapx.KindCode}, root, g, nil); v != Skip {
			t.Errorf("%s on a spec with no Endpoint: %v", name, v)
		}
		if v, d := check("", mapx.Node{ID: "api/other.go", Kind: mapx.KindCode}, root, g, nil); v != Skip || !strings.Contains(d, "spec") {
			t.Errorf("%s with no spec beside it: %v %s", name, v, d)
		}
	}
}

func TestAPISpec_contractsResolve(t *testing.T) {
	t.Run("APISP-B02: Every contract cited resolves to a spec with a Domain", func(t *testing.T) {})
	spec := apiHead + "## Request Body\n| Content type | Contract | Required |\n| --- | --- | --- |\n| application/json | `NODOM` | yes |\n\n" +
		"## Responses\n| Status | When | Content type | Contract |\n| --- | --- | --- | --- |\n" +
		"| 201 | created | application/json | `GHOST` |\n| 400 | refused | application/json | |\n| 204 | nothing | — | — |\n"
	root, g := apiProject(t, map[string]string{"api/generate.spec.md": spec,
		"api/nodom.spec.md": "<!-- @anchors\n  code: NODOM\n-->\n# No domain\n\n## Overview\n\nx\n"})
	v, d := checkAPIContractsResolve("", mapx.Node{ID: "api/generate.go", Kind: mapx.KindCode}, root, g, nil)
	if v != Fail || !strings.Contains(d, "not said for: 400") || !strings.Contains(d, "no spec of the project: GHOST") ||
		!strings.Contains(d, "no Domain table to give the fields: NODOM") || strings.Contains(d, "204") {
		t.Errorf("each kind of gap is named: %v %s", v, d)
	}
}

func TestAPISpec_errorsDeclared(t *testing.T) {
	t.Run("APISP-B03: Every error response is complete and under a declared status", func(t *testing.T) {})
	spec := apiHead + "## Responses\n| Status | When | Content type | Contract |\n| --- | --- | --- | --- |\n| 201 | ok | application/json | — |\n| 4xx | refused | application/json | — |\n\n" +
		"## Error Responses\n| Rule | When | Status | Error code | Message |\n| --- | --- | --- | --- | --- |\n" +
		"| `QRGEN-E01` | no wallet | 404 | WALLET_NOT_FOUND | \"Wallet not found\" |\n" +
		"| `QRGEN-E02` | provider down | 503 | PROVIDER_DOWN | \"Try again\" |\n" +
		"| `QRGEN-E03` | forbidden | 403 | WALLET_FORBIDDEN | |\n"
	root, g := apiProject(t, map[string]string{"api/generate.spec.md": spec})
	v, d := checkAPIErrorsDeclared("", mapx.Node{ID: "api/generate.go", Kind: mapx.KindCode}, root, g, nil)
	if v != Fail || !strings.Contains(d, "QRGEN-E02 (503)") || !strings.Contains(d, "without an error code or a message: QRGEN-E03") ||
		strings.Contains(d, "QRGEN-E01") {
		t.Errorf("503 is undeclared, E03 has no message, 404 is in 4xx: %v %s", v, d)
	}
}

func TestAPISpec_errorCodesHonored(t *testing.T) {
	t.Run("APISP-B04: Every declared error code is emitted by the unit's code", func(t *testing.T) {})
	t.Run("APISP-E02: A code file the spec specifies that cannot be read emits nothing", func(t *testing.T) {})
	spec := apiHead + "## Error Responses\n| Rule | When | Status | Error code | Message |\n| --- | --- | --- | --- | --- |\n" +
		"| `QRGEN-E01` | no wallet | 404 | WALLET_NOT_FOUND | \"Wallet not found\" |\n" +
		"| `QRGEN-E02` | forbidden | 403 | WALLET_FORBIDDEN | \"Not yours\" |\n" +
		"| `QRGEN-E03` | down | 503 | PROVIDER_DOWN | \"Try again\" |\n"
	root, g := apiProject(t, map[string]string{"api/generate.spec.md": spec,
		"api/generate.go": "package api\n// WALLET_FORBIDDEN is handled upstream\nvar errNoWallet = apiError(404, \"WALLET_NOT_FOUND\")\n",
		"api/errors.go":   "package api\nvar errDown = apiError(503, \"PROVIDER_DOWN\")\n",
	}, "api/errors.go", "api/missing.go")
	v, d := checkErrorCodesHonored("", mapx.Node{ID: "api/generate.go", Kind: mapx.KindCode}, root, g, nil)
	if v != Fail || !strings.Contains(d, "WALLET_FORBIDDEN") || strings.Contains(d, "WALLET_NOT_FOUND,") || strings.Contains(d, "PROVIDER_DOWN") {
		t.Errorf("only the code in a comment is absent: %v %s", v, d)
	}
}

func TestAPISpec_anyLanguage(t *testing.T) {
	t.Run("APISP-B05: Sections and columns are read in any language", func(t *testing.T) {})
	spec := apiHead + "## Respostas\n| Status | Quando | Content type | Contrato |\n| --- | --- | --- | --- |\n| 4xx | recusada | application/json | — |\n\n" +
		"## Respostas de Erro\n| Regra | Quando | Status | Código de erro | Mensagem |\n| --- | --- | --- | --- | --- |\n" +
		"| `QRGEN-E01` | sem carteira | 404 | WALLET_NOT_FOUND | \"Carteira não encontrada\" |\n"
	root, g := apiProject(t, map[string]string{"api/generate.spec.md": spec})
	if v, d := checkAPIErrorsDeclared("", mapx.Node{ID: "api/generate.go", Kind: mapx.KindCode}, root, g, nil); v != Pass {
		t.Errorf("a Portuguese spec is read: %v %s", v, d)
	}
}

func TestAPISpec_projectSectionTitles(t *testing.T) {
	t.Run("APISP-B06: A section renamed by the project is found", func(t *testing.T) {})
	SetProjectSectionTitles(&config.Config{Layers: map[string]config.Layer{"api": {SectionTitles: config.SectionTitles{"error-responses": "Recusas da API"}}}})
	t.Cleanup(func() { SetProjectSectionTitles(nil) })
	spec := apiHead + "## Responses\n| Status | When | Content type | Contract |\n| --- | --- | --- | --- |\n| 4xx | refused | application/json | — |\n\n" +
		"## Recusas da API\n| Rule | When | Status | Error code | Message |\n| --- | --- | --- | --- | --- |\n| `QRGEN-E01` | x | 503 | DOWN | \"Down\" |\n"
	root, g := apiProject(t, map[string]string{"api/generate.spec.md": spec})
	if v, d := checkAPIErrorsDeclared("", mapx.Node{ID: "api/generate.go", Kind: mapx.KindCode}, root, g, nil); v != Fail || !strings.Contains(d, "QRGEN-E01 (503)") {
		t.Errorf("the renamed section is read: %v %s", v, d)
	}
}

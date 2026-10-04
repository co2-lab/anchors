// @anchors
//   ref: ENVDC

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func envSpec(rows ...string) string {
	s := "<!-- @anchors\n  code: PAYMT\n-->\n# Pay\n\n## Environment Variables\n| Variable | Type | Required | Default | Values | Deprecated | Description |\n| --- | --- | --- | --- | --- | --- | --- |\n"
	for _, r := range rows {
		s += r + "\n"
	}
	return s
}

// envRun confronts the spec with the given code files, which the spec specifies.
func envRun(t *testing.T, spec string, cfg *config.Config, code map[string]string, missing ...string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	g := &mapx.Graph{}
	for rel, body := range code {
		p := filepath.Join(root, filepath.FromSlash(rel))
		must(t, os.MkdirAll(filepath.Dir(p), 0o755))
		must(t, os.WriteFile(p, []byte(body), 0o644))
		g.Edges = append(g.Edges, mapx.Edge{From: "pay/pay.spec.md", To: rel, Type: mapx.EdgeSpecifies})
	}
	for _, m := range missing {
		g.Edges = append(g.Edges, mapx.Edge{From: "pay/pay.spec.md", To: m, Type: mapx.EdgeSpecifies})
	}
	return checkEnvDeclared(spec, mapx.Node{ID: "pay/pay.spec.md", Kind: mapx.KindSpec}, root, g, cfg)
}

var goDialect = &config.Config{Dialect: &config.Dialect{Family: "go"}}

func TestEnvDeclared_nothingToConfront(t *testing.T) {
	t.Run("ENVDC-B01: Nothing declared and nothing read leaves without a verdict", func(t *testing.T) {})
	if v, _ := envRun(t, "# Pay\n", goDialect, map[string]string{"pay/pay.go": "package pay\n"}); v != Skip {
		t.Errorf("nothing declared, nothing read: %v", v)
	}
	if v, _ := checkEnvDeclared("", mapx.Node{ID: "pay/pay.go", Kind: mapx.KindCode}, t.TempDir(), &mapx.Graph{}, nil); v != Skip {
		t.Errorf("not a spec: %v", v)
	}
}

func TestEnvDeclared_bothSides(t *testing.T) {
	t.Run("ENVDC-B02: A variable read and not declared is named", func(t *testing.T) {})
	t.Run("ENVDC-B03: A variable declared and not read is named, unless deprecated", func(t *testing.T) {})
	t.Run("ENVDC-B04: A read in a comment is no read", func(t *testing.T) {})
	code := "package pay\nimport \"os\"\n// os.Getenv(\"DEBUG\") was here\nvar url = os.Getenv(\"PAYMENTS_URL\")\nvar key, _ = os.LookupEnv(\"API_KEY\")\n"
	spec := envSpec("| PAYMENTS_URL | string | yes | — | — | no | base URL |",
		"| OLD_HOST | string | no | — | — | no | the old host |",
		"| LEGACY_TOKEN | string | no | — | — | yes: use API_KEY | the old token |")
	v, d := envRun(t, spec, goDialect, map[string]string{"pay/pay.go": code})
	if v != Fail || !strings.Contains(d, "does not declare: API_KEY") || !strings.Contains(d, "does not read: OLD_HOST") ||
		strings.Contains(d, "LEGACY_TOKEN") || strings.Contains(d, "DEBUG") || strings.Contains(d, "PAYMENTS_URL") {
		t.Errorf("API_KEY undeclared, OLD_HOST unread, the deprecated and the comment left alone: %v %s", v, d)
	}
}

func TestEnvDeclared_everyFamilyWhenNoneIsDeclared(t *testing.T) {
	t.Run("ENVDC-B05: Every family's reads are recognised when none is declared", func(t *testing.T) {})
	spec := envSpec("| LOG_LEVEL | string | no | info | debug, info | no | how much is logged |")
	if v, d := envRun(t, spec, nil, map[string]string{"pay/pay.ts": "const level = process.env.LOG_LEVEL ?? 'info'\n"}); v != Pass {
		t.Errorf("TypeScript read with no dialect: %v %s", v, d)
	}
}

func TestEnvDeclared_anUnreadableFileReadsNothing(t *testing.T) {
	t.Run("ENVDC-E01: A specified file that cannot be read reads nothing", func(t *testing.T) {})
	spec := envSpec("| PORT | integer | yes | 8080 | — | no | the port |")
	if v, d := envRun(t, spec, goDialect, nil, "pay/missing.go"); v != Fail || !strings.Contains(d, "does not read: PORT") {
		t.Errorf("the declaration stands alone: %v %s", v, d)
	}
}

// @anchors
//   code: EDTNV
//   ref: ENVDC

package gate

import (
	"github.com/co2-lab/anchors/internal/testkit"
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

var goDialect = testkit.WithComments(&config.Config{Dialect: &config.Dialect{Family: "go"}})

func TestEnvDeclared_nothingToConfront(t *testing.T) {
	t.Run("ENVDC-B01: Nothing declared and nothing read leaves without a verdict", func(t *testing.T) {})
	if v, _ := envRun(t, "# Pay\n", goDialect, map[string]string{"pay/pay.go": "package pay\n"}); v != Skip {
		t.Errorf("nothing declared, nothing read: %v", v)
	}
	if v, _ := checkEnvDeclared("", mapx.Node{ID: "pay/pay.go", Kind: mapx.KindCode}, t.TempDir(), &mapx.Graph{}, testkit.WithComments(nil)); v != Skip {
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

func TestEnvDeclared_pendingWithoutEnvRead(t *testing.T) {
	t.Run("ENVDC-B05: With no env_read declared the gate is pending", func(t *testing.T) {})
	spec := envSpec("| LOG_LEVEL | string | no | info | debug, info | no | how much is logged |")
	cfg := &config.Config{Dialect: &config.Dialect{Comments: config.FamilyComments("ts")}}
	if v, d := envRun(t, spec, cfg, map[string]string{"pay/pay.ts": "const level = process.env.LOG_LEVEL ?? 'info'\n"}); v != Pending || !strings.Contains(d, "dialect.env_read") {
		t.Errorf("no env_read leaves the gate pending, naming the field: %v %s", v, d)
	}
}

func TestEnvDeclared_anUnreadableFileReadsNothing(t *testing.T) {
	t.Run("ENVDC-E01: A specified file that cannot be read reads nothing", func(t *testing.T) {})
	spec := envSpec("| PORT | integer | yes | 8080 | — | no | the port |")
	if v, d := envRun(t, spec, goDialect, nil, "pay/missing.go"); v != Fail || !strings.Contains(d, "does not read: PORT") {
		t.Errorf("the declaration stands alone: %v %s", v, d)
	}
}

func TestEnvDeclared_eachFileByItsDialect(t *testing.T) {
	t.Run("ENVDC-B06: Each code file is read with its own dialect's env_read, and one whose dialect declares none leaves the gate pending naming that dialect", func(t *testing.T) {})
	spec := envSpec("| `DB_URL` | string | yes | | | | db |", "| `API_KEY` | string | yes | | | | key |")
	code := map[string]string{
		"services/api/pay.go": "package pay\n\nvar u = os.Getenv(\"DB_URL\")\n",
		"apps/web/pay.ts":     "export const k = process.env.API_KEY\n",
	}
	if v, msg := envRun(t, spec, testkit.TwoDialects(), code); v != Pass {
		t.Errorf("the Go read by the Go dialect, the TS by the TS one: %v (%s)", v, msg)
	}
	cfg := testkit.TwoDialects()
	cfg.Dialects["go"].EnvRead = ""
	if v, msg := envRun(t, spec, cfg, code); v != Pending || !strings.Contains(msg, "dialects.go.env_read") {
		t.Errorf("the Go dialect declares no env_read: pending naming it, got %v (%s)", v, msg)
	}
}

func TestEnvDeclared_readByADependency(t *testing.T) {
	t.Run("ENVDC-B07: A declared variable the unit's dependencies read is read", func(t *testing.T) {})
	root := t.TempDir()
	for rel, body := range map[string]string{
		"pay/pay.go":     "package pay\nvar url = os.Getenv(\"PAYMENTS_URL\")\n",
		"quota/quota.go": "package quota\n",
		"usage/usage.go": "package usage\nvar table = os.Getenv(\"USAGE_TABLE\")\n",
		"other/other.go": "package other\nvar k = os.Getenv(\"STRAY_KEY\")\n",
	} {
		p := filepath.Join(root, filepath.FromSlash(rel))
		must(t, os.MkdirAll(filepath.Dir(p), 0o755))
		must(t, os.WriteFile(p, []byte(body), 0o644))
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "pay/pay.go", Kind: mapx.KindCode}, {ID: "quota/quota.go", Kind: mapx.KindCode},
			{ID: "usage/usage.go", Kind: mapx.KindCode}, {ID: "other/other.go", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{
			{From: "pay/pay.spec.md", To: "pay/pay.go", Type: mapx.EdgeSpecifies},
			{From: "pay/pay.go", To: "quota/quota.go", Type: mapx.EdgeDependsOn},
			{From: "quota/quota.go", To: "usage/usage.go", Type: mapx.EdgeDependsOn},
		},
	}
	n := mapx.Node{ID: "pay/pay.spec.md", Kind: mapx.KindSpec}
	spec := envSpec("| PAYMENTS_URL | string | yes | — | — | no | base URL |",
		"| USAGE_TABLE | string | yes | — | — | no | read by the quota helper's model |")
	if v, msg := checkEnvDeclared(spec, n, root, g, goDialect); v != Pass {
		t.Errorf("a variable a dependency of a dependency reads is the unit's to declare: %v %s", v, msg)
	}
	stray := envSpec("| PAYMENTS_URL | string | yes | — | — | no | base URL |",
		"| STRAY_KEY | string | yes | — | — | no | read by no dependency |")
	if v, msg := checkEnvDeclared(stray, n, root, g, goDialect); v != Fail || !strings.Contains(msg, "STRAY_KEY") {
		t.Errorf("a file the unit does not depend on reads nothing for it: %v %s", v, msg)
	}
}

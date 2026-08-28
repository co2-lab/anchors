package recode

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// projeto sintético mínimo p/ exercitar BuildPlan com o bloco recode. Cria um
// anchors.yaml com uma layer que casa .tsx/.spec.md e o bloco recode declarado.
func writeProject(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func testCfg() *config.Config {
	return &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec": {Pattern: "**/*.spec.md", Kind: "spec"},
			"code": {Pattern: "**/*.tsx", Kind: "code"},
			"e2e":  {Pattern: "**/*.yaml", Kind: "test"},
		},
		Recode: &config.Recode{
			TestID:       "lower",
			FilePatterns: []string{"**/{{code}}-*.yaml"},
		},
	}
}

func TestBuildPlan_dialect_conformingProject(t *testing.T) {
	// projeto que SEGUE a convenção: testID = código minúsculo (abcdx-).
	root := writeProject(t, map[string]string{
		"Foo.spec.md":          "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-S01 x\n",
		"Foo.tsx":              "// @anchors\n//   ref: ABCDX\nconst x = <View testID=\"abcdx-root\" />\n",
		"flows/ABCDX-S01.yaml": "tags:\n  - ABCDX-S01\n",
	})
	plan, err := BuildPlan(root, testCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	if plan.TestIDs != 1 {
		t.Errorf("esperava 1 testID trocado (abcdx-→wxyzx-), veio %d", plan.TestIDs)
	}
	if len(plan.Renames) != 1 || plan.Renames[0].To != "flows/WXYZX-S01.yaml" {
		t.Errorf("esperava rename ABCDX-S01.yaml→WXYZX-S01.yaml, veio %+v", plan.Renames)
	}
	if plan.TestIDLegacy != "" {
		t.Errorf("projeto conforme NÃO deveria ter aviso de legado: %q", plan.TestIDLegacy)
	}
}

func TestBuildPlan_dialect_legacyWarns(t *testing.T) {
	// projeto LEGADO: código ABCDX mas testID usa prefixo divergente (zzzz-).
	root := writeProject(t, map[string]string{
		"Foo.spec.md": "<!-- @anchors\n  code: ABCDX\n-->\n### ABCDX-S01 x\n",
		"Foo.tsx":     "// @anchors\n//   ref: ABCDX\nconst x = <View testID=\"zzzz-root\" />\n",
	})
	plan, err := BuildPlan(root, testCfg(), "ABCDX", "WXYZX")
	if err != nil {
		t.Fatal(err)
	}
	if plan.TestIDs != 0 {
		t.Errorf("prefixo esperado ausente → 0 trocas, veio %d", plan.TestIDs)
	}
	if plan.TestIDLegacy == "" {
		t.Error("legado deveria emitir o aviso de testID divergente")
	}
}

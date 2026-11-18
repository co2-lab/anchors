// @anchors
//   code: SCRQT
//   ref: SCRQR

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testkit"
)

// sectionCfg is a project whose components call for Test IDs when they expose a handle, a
// data contract when they depend on a store, Navigation when they carry the flag and
// Messages when they show a toast.
func sectionCfg() *config.Config {
	return testkit.WithComments(&config.Config{
		Layers: map[string]config.Layer{
			"component": {Pattern: "src/components/*.tsx", Kind: "code", RequiresSections: []config.SectionRequirement{
				{Section: "Test IDs", When: config.SectionTrigger{Exposes: "test_handles"}},
				{Section: "Contrato de Dados", When: config.SectionTrigger{DependsOnLayer: []string{"store"}}},
				{Section: "Navegação", When: config.SectionTrigger{Flag: "@navigates"}},
				{Section: "Mensagens ao Usuário", When: config.SectionTrigger{CodeMatches: `showToast\(`}},
			}},
			"store": {Pattern: "src/stores/*.ts", Kind: "code"},
		},
		Dialect: &config.Dialect{TestHandles: config.FamilyTestHandles("ts", "testID")},
	})
}

// sectionUnit writes the component's code and returns its spec node, the map and the root.
func sectionUnit(t *testing.T, code string, dependsOnStore bool) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "src/components/Login.tsx", code)
	spec := mapx.Node{ID: "src/components/Login.spec.md", Kind: mapx.KindSpec}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "src/components/Login.tsx", Kind: mapx.KindCode, Layer: "component"},
			{ID: "src/stores/auth.ts", Kind: mapx.KindCode, Layer: "store"}},
		Edges: []mapx.Edge{{From: spec.ID, To: "src/components/Login.tsx", Type: mapx.EdgeSpecifies}},
	}
	if dependsOnStore {
		g.Edges = append(g.Edges, mapx.Edge{From: "src/components/Login.tsx", To: "src/stores/auth.ts", Type: mapx.EdgeDependsOn})
	}
	return spec, g, root
}

func TestSectionRequired_nothingToConfront(t *testing.T) {
	t.Run("SCRQR-B01: Nothing to confront leaves without a verdict", func(t *testing.T) {})
	cfg := sectionCfg()
	if v, _ := checkSectionRequired("", mapx.Node{ID: "a.tsx", Kind: mapx.KindCode}, "", &mapx.Graph{}, cfg); v != Skip {
		t.Errorf("not a spec: %v", v)
	}
	lone := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkSectionRequired("# X\n", lone, "", &mapx.Graph{Nodes: []mapx.Node{lone}}, cfg); v != Skip {
		t.Errorf("a spec with no code is unit-complete's: %v", v)
	}
	spec, g, root := sectionUnit(t, `<View testID=":login-submit" />`, false)
	g.Nodes[1].Layer = "store"
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, cfg); v != Skip {
		t.Errorf("a layer that declares nothing calls for nothing: %v %s", v, msg)
	}
}

func TestSectionRequired_missingSectionsFail(t *testing.T) {
	t.Run("SCRQR-B02: A triggered section the spec lacks fails, named once", func(t *testing.T) {})
	t.Run("SCRQR-B03: A section present satisfies the requirement", func(t *testing.T) {})
	spec, g, root := sectionUnit(t, `export const Login = () => <View testID=":login-submit" />`+"\n", true)
	v, msg := checkSectionRequired("# Login\n\n## Overview\n", spec, root, g, sectionCfg())
	if v != Fail || !strings.Contains(msg, "`## Test IDs` — `Login.tsx` exposes the handle `:login-submit`") ||
		!strings.Contains(msg, "`## Contrato de Dados` — `Login.tsx` depends on `src/stores/auth.ts` (layer store)") ||
		strings.Contains(msg, "Navegação") || strings.Contains(msg, "Mensagens") {
		t.Fatalf("both triggered sections named, the others not: %v\n%s", v, msg)
	}
	ok := "# Login\n\n### test ids\n\nnone yet\n\n## Contrato de Dados\n"
	if v, msg := checkSectionRequired(ok, spec, root, g, sectionCfg()); v != Pass {
		t.Errorf("headings at any level and in any case satisfy: %v %s", v, msg)
	}
}

func TestSectionRequired_triggers(t *testing.T) {
	t.Run("SCRQR-B04: Each trigger reads what the core measures", func(t *testing.T) {})
	quiet := "export const Login = () => {\n  const s = '@navigates' // showToast( is not called here\n  return <View />\n}\n"
	spec, g, root := sectionUnit(t, quiet, false)
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, sectionCfg()); v != Pass {
		t.Fatalf("nothing triggered — the flag only in a string, the pattern only in a comment: %v %s", v, msg)
	}
	loud := "export const Login = () => {\n  showToast('ok')\n  return <View /> // @navigates: HOMEX\n}\n"
	spec, g, root = sectionUnit(t, loud, false)
	v, msg := checkSectionRequired("# Login\n", spec, root, g, sectionCfg())
	if v != Fail || !strings.Contains(msg, "`## Navegação` — `Login.tsx` carries `@navigates`") ||
		!strings.Contains(msg, "`## Mensagens ao Usuário` — `Login.tsx` matches the project's pattern at `showToast(`") {
		t.Errorf("the flag in a comment and the pattern in the code call for their sections: %v\n%s", v, msg)
	}
}

func TestSectionRequired_unmeasurableIsPending(t *testing.T) {
	t.Run("SCRQR-B05: A trigger that cannot be measured is pending, not silent", func(t *testing.T) {})
	t.Run("SCRQR-E01: A code file that cannot be read calls for nothing", func(t *testing.T) {})
	cfg := sectionCfg()
	cfg.Dialect.TestHandles = nil
	spec, g, root := sectionUnit(t, `<View testID=":login-submit" />`, false)
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, cfg); v != Pending || !strings.Contains(msg, "dialect.test_handles.exposed") {
		t.Errorf("no exposed patterns: pending, naming the field: %v %s", v, msg)
	}
	g.Nodes[1].ID, g.Edges[0].To = "src/components/Gone.tsx", "src/components/Gone.tsx"
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, sectionCfg()); v != Pass {
		t.Errorf("a file that cannot be read calls for nothing: %v %s", v, msg)
	}
	if v, _ := checkSectionRequired("# Login\n", spec, root, nil, sectionCfg()); v != Pending {
		t.Errorf("no map: pending, %v", v)
	}
	spec, g, root = sectionUnit(t, "<Pressable testID={`${prefix}${r.k}`} />\n", false)
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, sectionCfg()); v != Pending || !strings.Contains(msg, "Login.tsx:1") {
		t.Errorf("a handle of values alone: pending, %v %s", v, msg)
	}
}

func TestSectionRequired_typeOnlyDependency(t *testing.T) {
	t.Run("SCRQR-B06: A dependency on types only calls for nothing unless asked", func(t *testing.T) {})
	spec, g, root := sectionUnit(t, "export const Login = () => null\n", true)
	g.Edges[1].Kind = "type"
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, sectionCfg()); v != Pass {
		t.Errorf("the shape of the data is not the data: %v %s", v, msg)
	}
	cfg := sectionCfg()
	l := cfg.Layers["component"]
	l.RequiresSections[1].When.Kinds = []string{"value", "type"}
	cfg.Layers["component"] = l
	if v, msg := checkSectionRequired("# Login\n", spec, root, g, cfg); v != Fail || !strings.Contains(msg, "depends on the types of `src/stores/auth.ts` (layer store)") {
		t.Errorf("asked for, a type dependency calls for the section: %v %s", v, msg)
	}
}

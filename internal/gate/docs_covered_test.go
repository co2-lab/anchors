package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
)

// projectWithOrphan builds a project where ONE spec is out of the templates' reach: the
// template asks for `layer=gate`, and spec B belongs to another layer.
func projectWithOrphan(t *testing.T) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	spec := func(code, layer string) string {
		return "---\ncode: " + code + "\nlayer: " + layer + "\n---\n\n# T — t\n\n## Visão Geral\n\nx.\n"
	}
	os.WriteFile(filepath.Join(root, "pkg/A.spec.md"), []byte(spec("AAAAA", "gate")), 0o644)
	os.WriteFile(filepath.Join(root, "pkg/B.spec.md"), []byte(spec("BBBBB", "orfa")), 0o644)
	os.MkdirAll(filepath.Join(root, doct.Dir), 0o755)
	os.WriteFile(filepath.Join(root, doct.Dir, "g.md.tmpl"),
		[]byte(`{{range specs "layer=gate"}}{{section . "Visão Geral"}}{{end}}`), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/A.spec.md", Kind: mapx.KindSpec, Code: "AAAAA", Layer: "gate"},
		{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB", Layer: "orfa"},
	}}
	return root, g
}

// The defect this gate exists to catch: the unit has a spec, a triad, passes every
// relational gate — and is documented nowhere.
func TestDocsCovered_flagsASpecOutsideTheTemplates(t *testing.T) {
	t.Run("DCCVD-B03: A spec no template reaches fails, naming the spec", func(t *testing.T) {})
	resetDocsCoverageCache()
	root, g := projectWithOrphan(t)
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB"}
	v, d := checkDocsCovered("", n, root, g, nil)
	if v != Fail {
		t.Fatalf("the orphan spec should fail, got %v (%s)", v, d)
	}
	if !strings.Contains(d, "pkg/B.spec.md") {
		t.Errorf("the verdict must NAME the orphan spec: %s", d)
	}
}

// The verdict is about THIS target. Reporting the other orphans would charge one file for
// what another lacks — and A's author has nothing to do about B's problem.
func TestDocsCovered_reachedSpecPasses(t *testing.T) {
	t.Run("DCCVD-B04: A spec a template reaches passes", func(t *testing.T) {})
	t.Run("DCCVD-I01: Another spec's orphan status never fails the confronted spec", func(t *testing.T) {})
	t.Run("DCCVD-X01: A reached spec passes without any page having been built", func(t *testing.T) {})
	resetDocsCoverageCache()
	root, g := projectWithOrphan(t)
	n := mapx.Node{ID: "pkg/A.spec.md", Kind: mapx.KindSpec, Code: "AAAAA"}
	if v, d := checkDocsCovered("", n, root, g, nil); v != Pass {
		t.Errorf("the spec the template reaches should pass, got %v (%s)", v, d)
	}
}

func TestDocsCovered_onlyConfrontsSpecs(t *testing.T) {
	t.Run("DCCVD-B01: An artifact that is not a spec is skipped", func(t *testing.T) {})
	resetDocsCoverageCache()
	root, g := projectWithOrphan(t)
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindCode, Code: "BBBBB"}
	if v, _ := checkDocsCovered("", n, root, g, nil); v != Skip {
		t.Errorf("expected Skip for a non-spec, got %v", v)
	}
}

// Without `doct/` there are no templates and no coverage to charge: demanding
// documentation from a project that declared none would invent a duty.
func TestDocsCovered_skipsAProjectWithoutTemplates(t *testing.T) {
	t.Run("DCCVD-B02: A project with no templates directory is skipped", func(t *testing.T) {})
	resetDocsCoverageCache()
	root, g := projectWithOrphan(t)
	os.RemoveAll(filepath.Join(root, doct.Dir))
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB"}
	if v, _ := checkDocsCovered("", n, root, g, nil); v != Skip {
		t.Errorf("expected Skip without templates, got %v", v)
	}
}

// A template that does not compile is docs-fresh's finding, with the compiler's report.
// Repeating it here would make the author think there are two defects.
func TestDocsCovered_aTemplateThatDoesNotCompileIsLeftToTheSibling(t *testing.T) {
	t.Run("DCCVD-B05: A template that does not compile skips with no message", func(t *testing.T) {})
	resetDocsCoverageCache()
	root, g := projectWithOrphan(t)
	if err := os.WriteFile(filepath.Join(root, doct.Dir, "g.md.tmpl"), []byte(`{{range specs "layer=gate"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB"}
	if v, d := checkDocsCovered("", n, root, g, nil); v != Skip || d != "" {
		t.Errorf("a broken template is not this gate's finding: %v (%q)", v, d)
	}
}

// The answer is the same for every target of a scan, so the templates are compiled once
// per (root, map): a later change on disk is not seen until the map changes.
func TestDocsCovered_compilesOncePerRootAndMap(t *testing.T) {
	t.Run("DCCVD-B06: The coverage is computed once per project root and map", func(t *testing.T) {})
	resetDocsCoverageCache()
	root, g := projectWithOrphan(t)
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB"}
	if v, _ := checkDocsCovered("", n, root, g, nil); v != Fail {
		t.Fatalf("setup: B starts as an orphan, got %v", v)
	}
	// widen the template so it reaches every spec
	if err := os.WriteFile(filepath.Join(root, doct.Dir, "g.md.tmpl"),
		[]byte(`{{range specs ""}}{{section . "Visão Geral"}}{{end}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkDocsCovered("", n, root, g, nil); v != Fail {
		t.Errorf("the same root and map reuse the first answer, got %v", v)
	}
	fresh := &mapx.Graph{Nodes: append([]mapx.Node{}, g.Nodes...)}
	if v, d := checkDocsCovered("", n, root, fresh, nil); v != Pass {
		t.Errorf("a new map recomputes, and the widened template reaches B: %v (%s)", v, d)
	}
}

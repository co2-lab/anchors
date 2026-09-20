package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
)

// projetoComOrfa monta um projeto onde UMA spec fica fora do alcance dos templates: o
// template pede `layer=gate`, e a spec B e' de outra camada.
func projetoComOrfa(t *testing.T) (string, *mapx.Graph) {
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

// O defeito que este gate existe para pegar: a unidade tem spec, tem trinca, passa por
// todos os gates relacionais — e nao esta documentada em lugar nenhum.
func TestDocsCovered_acusaSpecForaDosTemplates(t *testing.T) {
	resetDocsCoverageCache()
	root, g := projetoComOrfa(t)
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB"}
	v, d := checkDocsCovered("", n, root, g, nil)
	if v != Fail {
		t.Fatalf("a spec orfa devia reprovar, veio %v (%s)", v, d)
	}
	if !strings.Contains(d, "pkg/B.spec.md") {
		t.Errorf("o veredito tem de NOMEAR a spec orfa: %s", d)
	}
}

// O veredito e' sobre ESTE alvo. Reportar as orfas das outras acusaria um arquivo pelo
// que falta noutro — e o autor de A nao tem o que fazer com o problema de B.
func TestDocsCovered_specAlcancadaPassa(t *testing.T) {
	resetDocsCoverageCache()
	root, g := projetoComOrfa(t)
	n := mapx.Node{ID: "pkg/A.spec.md", Kind: mapx.KindSpec, Code: "AAAAA"}
	if v, d := checkDocsCovered("", n, root, g, nil); v != Pass {
		t.Errorf("a spec alcancada pelo template devia passar, veio %v (%s)", v, d)
	}
}

func TestDocsCovered_soConfrontaSpec(t *testing.T) {
	resetDocsCoverageCache()
	root, g := projetoComOrfa(t)
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindCode, Code: "BBBBB"}
	if v, _ := checkDocsCovered("", n, root, g, nil); v != Skip {
		t.Errorf("esperava Skip para nao-spec, veio %v", v)
	}
}

// Sem `doct/` nao ha templates, e nao ha cobertura a cobrar: exigir documentacao de um
// projeto que nao declarou nenhuma seria inventar um dever.
func TestDocsCovered_pulaProjetoSemTemplates(t *testing.T) {
	resetDocsCoverageCache()
	root, g := projetoComOrfa(t)
	os.RemoveAll(filepath.Join(root, doct.Dir))
	n := mapx.Node{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB"}
	if v, _ := checkDocsCovered("", n, root, g, nil); v != Skip {
		t.Errorf("esperava Skip sem templates, veio %v", v)
	}
}

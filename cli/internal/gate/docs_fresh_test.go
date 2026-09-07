package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
)

const specParaDoc = `---
code: GLCGL
layer: infra
---

# GoLiveChecklist — a régua

## Visão Geral

O conteúdo original.
`

func projetoComDoc(t *testing.T) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(specParaDoc), 0o644)
	os.MkdirAll(filepath.Join(root, doct.Dir), 0o755)
	os.WriteFile(filepath.Join(root, doct.Dir, "infra.md.tmpl"),
		[]byte(`{{range specs "layer=infra"}}{{section . "Visão Geral"}}{{end}}`), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/GoLive.spec.md", Kind: mapx.KindSpec, Code: "GLCGL", Layer: "spec"},
	}}
	return root, g
}

func noSpec() mapx.Node {
	return mapx.Node{ID: "pkg/GoLive.spec.md", Kind: mapx.KindSpec, Code: "GLCGL"}
}

// O defeito que o gate existe para pegar: a spec mudou, o compilado não foi refeito, e a
// documentação segue afirmando a versão antiga — com conteúdo real, e por isso convincente.
func TestDocsFresh_acusaCompiladoDefasado(t *testing.T) {
	root, g := projetoComDoc(t)
	c, _ := doct.New(root, g)
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}

	// A spec muda; ninguém recompila.
	novo := strings.Replace(specParaDoc, "O conteúdo original.", "O conteúdo REVISADO.", 1)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(novo), 0o644)

	v, msg := checkDocsFresh(novo, noSpec(), root, g, nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s) — a doc afirma a versão antiga e o gate passou", v, msg)
	}
	if !strings.Contains(msg, "infra.md") {
		t.Errorf("a mensagem não diz QUAL documento: %s", msg)
	}
}

func TestDocsFresh_passaComDocEmDia(t *testing.T) {
	root, g := projetoComDoc(t)
	c, _ := doct.New(root, g)
	c.Build(false)

	if v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a doc acabou de ser compilada", v, msg)
	}
}

// Projeto sem `doct/` não usa doc compilada: cobrar dele seria cobrar de quem não optou
// pelo mecanismo.
func TestDocsFresh_pulaProjetoSemTemplates(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(specParaDoc), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "pkg/GoLive.spec.md", Kind: mapx.KindSpec}}}

	if v, _ := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Skip {
		t.Errorf("verdict = %v, queria Skip", v)
	}
}

// O gate NÃO escreve. Se ele consertasse o que aponta, a segunda execução sempre passaria
// e o defeito só apareceria em quem clonasse o repositório.
func TestDocsFresh_naoEscreveNoDisco(t *testing.T) {
	root, g := projetoComDoc(t)
	destino := filepath.Join(root, doct.OutDir, "infra.md")

	checkDocsFresh(specParaDoc, noSpec(), root, g, nil)

	if _, err := os.Stat(destino); err == nil {
		t.Error("o gate compilou a doc — ele aponta, não conserta")
	}
}

// Doc escrita à mão (sem marcador) não é cobrada: o `Build` se recusa a sobrescrevê-la, e
// mandar rodar um comando que não muda nada seria um aviso que ninguém consegue resolver.
func TestDocsFresh_ignoraDocEscritaAMao(t *testing.T) {
	root, g := projetoComDoc(t)
	os.MkdirAll(filepath.Join(root, doct.OutDir), 0o755)
	os.WriteFile(filepath.Join(root, doct.OutDir, "infra.md"),
		[]byte("# Escrito à mão\n"), 0o644)

	if v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Pass {
		t.Errorf("verdict = %v (%s) — o Build não sobrescreve este arquivo", v, msg)
	}
}

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

// compila a doc do projeto e devolve o conteúdo da spec REVISADA — o estado em que o
// compilado ficou para trás.
func defasa(t *testing.T, root string, g *mapx.Graph) string {
	t.Helper()
	c, err := doct.New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	novo := strings.Replace(specParaDoc, "O conteúdo original.", "O conteúdo REVISADO.", 1)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(novo), 0o644)
	return novo
}

// O defeito que o gate existe para pegar: a spec mudou, o compilado não foi refeito, e a
// documentação segue afirmando a versão antiga — com conteúdo real, e por isso convincente.
func TestDocsFresh_acusaCompiladoDefasado(t *testing.T) {
	t.Run("DCFRD-B01: A compiled document that no longer matches the spec fails", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	novo := defasa(t, root, g)

	v, msg := checkDocsFresh(novo, noSpec(), root, g, nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s) — a doc afirma a versão antiga e o gate passou", v, msg)
	}
}

// Apontar sem DIZER QUAL manda o autor comparar o diretório inteiro à mão.
func TestDocsFresh_dizQualDocumento(t *testing.T) {
	t.Run("DCFRD-B02: The verdict names the stale documents", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	novo := defasa(t, root, g)

	_, msg := checkDocsFresh(novo, noSpec(), root, g, nil)
	if !strings.Contains(msg, "infra.md") {
		t.Errorf("a mensagem não diz QUAL documento: %s", msg)
	}
	if !strings.Contains(msg, doct.OutDir) {
		t.Errorf("a mensagem não diz ONDE o documento mora: %s", msg)
	}
}

func TestDocsFresh_passaComDocEmDia(t *testing.T) {
	t.Run("DCFRD-B03: A compiled document that matches the templates passes", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	c, _ := doct.New(root, g)
	c.Build(false)

	if v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a doc acabou de ser compilada", v, msg)
	}
}

// A cobrança parte da FONTE. Código, teste e feature não alimentam template nenhum.
func TestDocsFresh_soConfrontaSpec(t *testing.T) {
	t.Run("DCFRD-B04: An artifact that is not a spec leaves without a verdict", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	defasa(t, root, g)
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature, mapx.KindGuide} {
		if v, _ := checkDocsFresh("", mapx.Node{ID: "pkg/GoLive.go", Kind: k}, root, g, nil); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}

// Projeto sem `doct/` não usa doc compilada: cobrar dele seria cobrar de quem não optou
// pelo mecanismo.
func TestDocsFresh_pulaProjetoSemTemplates(t *testing.T) {
	t.Run("DCFRD-B05: A project with no template directory is not charged", func(t *testing.T) {})
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(specParaDoc), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "pkg/GoLive.spec.md", Kind: mapx.KindSpec}}}

	if v, _ := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Skip {
		t.Errorf("verdict = %v, queria Skip", v)
	}
}

// Template declarado e documento NUNCA produzido: a ausência é defasagem. Tratá-la como
// satisfeita deixaria o projeto inteiro sem documentação e verde.
func TestDocsFresh_ausenciaEhDefasagem(t *testing.T) {
	t.Run("DCFRD-B06: A template whose document was never produced counts as stale", func(t *testing.T) {})
	root, g := projetoComDoc(t) // template existe, nada foi compilado

	v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s) — não há nada afirmando o que a spec diz hoje", v, msg)
	}
	if !strings.Contains(msg, "infra.md") {
		t.Errorf("a mensagem não nomeia o documento ausente: %s", msg)
	}
}

// Doc escrita à mão (sem marcador) não é cobrada: o `Build` se recusa a sobrescrevê-la, e
// mandar rodar um comando que não muda nada seria um aviso que ninguém consegue resolver.
func TestDocsFresh_ignoraDocEscritaAMao(t *testing.T) {
	t.Run("DCFRD-B07: A document written by hand is not charged", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	os.MkdirAll(filepath.Join(root, doct.OutDir), 0o755)
	os.WriteFile(filepath.Join(root, doct.OutDir, "infra.md"),
		[]byte("# Escrito à mão\n"), 0o644)

	if v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Pass {
		t.Errorf("verdict = %v (%s) — o Build não sobrescreve este arquivo", v, msg)
	}
}

// Template que não compila é DEFEITO. Silenciar aqui esconderia a documentação inteira
// atrás de um erro que ninguém veria.
func TestDocsFresh_templateQuebradoReprova(t *testing.T) {
	t.Run("DCFRD-B08: A template that cannot be compiled fails", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	os.WriteFile(filepath.Join(root, doct.Dir, "infra.md.tmpl"),
		[]byte(`{{range specs "camada=naoexiste"}}{{end}}`), 0o644)

	v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s) — template quebrado é defeito, não silêncio", v, msg)
	}
	if strings.TrimSpace(msg) == "" {
		t.Error("a mensagem não carrega o erro do compilador")
	}
}

// O mapa conhece uma spec que o disco não tem: o gate NÃO consegue olhar. Aprovar aqui
// carimbaria o que não foi medido.
func TestDocsFresh_semConseguirLerAsSpecsNaoAprova(t *testing.T) {
	t.Run("DCFRD-B09: A project whose specs cannot be read leaves without a verdict", func(t *testing.T) {})
	root, _ := projetoComDoc(t)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/GoLive.spec.md", Kind: mapx.KindSpec, Code: "GLCGL"},
		{ID: "pkg/Fantasma.spec.md", Kind: mapx.KindSpec, Code: "FANTX"},
	}}

	v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil)
	if v == Pass {
		t.Fatalf("aprovou sem conseguir ler as specs (%s)", msg)
	}
	if v != Skip {
		t.Fatalf("verdict = %v (%s), queria Skip com a razão", v, msg)
	}
	if strings.TrimSpace(msg) == "" {
		t.Error("o Skip não carrega a razão de não ter conseguido olhar")
	}
}

// O gate NÃO escreve. Se ele consertasse o que aponta, a segunda execução sempre passaria
// e o defeito só apareceria em quem clonasse o repositório.
func TestDocsFresh_naoEscreveNoDisco(t *testing.T) {
	t.Run("DCFRD-B10: The comparison happens in memory", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	destino := filepath.Join(root, doct.OutDir, "infra.md")

	checkDocsFresh(specParaDoc, noSpec(), root, g, nil)

	if _, err := os.Stat(destino); err == nil {
		t.Error("o gate compilou a doc — ele aponta, não conserta")
	}
}

// Rodar duas vezes tem de dar o MESMO veredito: um gate que consertasse passaria na
// segunda, e o defeito só apareceria em quem clonasse o repositório.
func TestDocsFresh_segundaExecucaoDaOMesmoVeredito(t *testing.T) {
	t.Run("DCFRD-I01: The gate never repairs what it points at", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	novo := defasa(t, root, g)

	primeira, _ := checkDocsFresh(novo, noSpec(), root, g, nil)
	segunda, msg := checkDocsFresh(novo, noSpec(), root, g, nil)
	if primeira != Fail || segunda != Fail {
		t.Fatalf("primeira %s, segunda %s (%s) — a segunda passou porque o gate consertou",
			primeira, segunda, msg)
	}
}

// A cobrança pousa na SPEC. O compilado é saída de compilador e não é nó do mapa.
func TestDocsFresh_aCobrancaPousaNaSpec(t *testing.T) {
	t.Run("DCFRD-I02: The charge starts from the spec and never from the compiled document", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	novo := defasa(t, root, g)

	naSpec, _ := checkDocsFresh(novo, noSpec(), root, g, nil)
	noCompilado, _ := checkDocsFresh(novo,
		mapx.Node{ID: filepath.Join(doct.OutDir, "infra.md"), Kind: mapx.KindGuide}, root, g, nil)
	if naSpec != Fail {
		t.Errorf("a spec deveria receber a cobrança, foi %s", naSpec)
	}
	if noCompilado != Skip {
		t.Errorf("o compilado não é unidade declarada, foi %s", noCompilado)
	}
}

// A régua é IGUALDADE com o que os templates produzem. Um template que escolhe mal passa
// — a escolha é de quem o escreveu.
func TestDocsFresh_naoJulgaAQualidadeDoCompilado(t *testing.T) {
	t.Run("DCFRD-X01: The gate does not judge whether the compiled document is good", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	// um template que não seleciona nada de útil, mas compila
	os.WriteFile(filepath.Join(root, doct.Dir, "infra.md.tmpl"), []byte("nada de útil aqui\n"), 0o644)
	c, _ := doct.New(root, g)
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}

	if v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a régua é igualdade, não qualidade", v, msg)
	}
}

// O gate não roda o build mesmo sabendo o conserto: o `check` roda em hook e em CI, e um
// gate que altera arquivos torna o resultado dependente de já ter rodado.
func TestDocsFresh_naoRodaOBuild(t *testing.T) {
	t.Run("DCFRD-X02: The gate does not run the build even knowing the fix", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	novo := defasa(t, root, g)
	destino := filepath.Join(root, doct.OutDir, "infra.md")
	antes, err := os.ReadFile(destino)
	if err != nil {
		t.Fatal(err)
	}

	checkDocsFresh(novo, noSpec(), root, g, nil)

	depois, err := os.ReadFile(destino)
	if err != nil {
		t.Fatal(err)
	}
	if string(antes) != string(depois) {
		t.Error("o gate reescreveu o compilado — ele aponta, não conserta")
	}
}

// Escrita à mão CONTRA um template que produziria outra coisa: ainda assim não é cobrada,
// porque o compilador se recusa a escrever aquele arquivo.
func TestDocsFresh_naoCobraDocEscritaAMaoNemContraTemplateDivergente(t *testing.T) {
	t.Run("DCFRD-X03: The gate does not charge documents written by hand", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	os.MkdirAll(filepath.Join(root, doct.OutDir), 0o755)
	os.WriteFile(filepath.Join(root, doct.OutDir, "infra.md"),
		[]byte("# Uma página completamente diferente do que o template produziria\n"), 0o644)

	if v, msg := checkDocsFresh(specParaDoc, noSpec(), root, g, nil); v != Pass {
		t.Errorf("verdict = %v (%s) — o compilador não escreve este arquivo", v, msg)
	}
}

// O compilado não é artefato próprio: dar nó a ele cobraria spec, feature e teste de uma
// saída de build.
func TestDocsFresh_oCompiladoNaoEhArtefatoProprio(t *testing.T) {
	t.Run("DCFRD-X04: The compiled document is not confronted as an artifact of its own", func(t *testing.T) {})
	root, g := projetoComDoc(t)
	c, _ := doct.New(root, g)
	c.Build(false)

	no := mapx.Node{ID: filepath.Join(doct.OutDir, "infra.md"), Kind: mapx.KindGuide}
	if v, _ := checkDocsFresh("", no, root, g, nil); v != Skip {
		t.Errorf("o compilado recebeu veredito (%s) — ele é saída de build", v)
	}
}

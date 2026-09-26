package doct

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const specDeExemplo = `---
code: GLCGL
layer: infra
---

# GoLiveChecklist — a régua que registra o que não foi feito

## Visão Geral

A régua confronta o que foi prometido contra o que existe.

## Regras

### GLCGL-B01 — item sem artefato não passa

Um item que não cita artefato não é conferível.

### GLCGL-B02 — dívida aberta bloqueia

Enquanto a dívida está aberta, a régua não libera.

## Invariantes

### GLCGL-I01 — o resumo nunca mente

O resumo conta o que a régua viu.
`

// monta um projeto de mentira com uma spec no disco e um mapa que a conhece.
func projetoDeTeste(t *testing.T, specs map[string]string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	g := &mapx.Graph{}
	for caminho, conteudo := range specs {
		full := filepath.Join(root, caminho)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
		code := ""
		if m := headerCodeDeTeste.FindStringSubmatch(conteudo); m != nil {
			code = m[1]
		}
		g.Nodes = append(g.Nodes, mapx.Node{ID: caminho, Kind: mapx.KindSpec, Code: code, Layer: "spec"})
	}
	return root, g
}

var headerCodeDeTeste = regexp.MustCompile(`(?m)^\s*code:\s*([A-Z0-9]+)`)

func escreveTemplate(t *testing.T, root, nome, corpo string) {
	t.Helper()
	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, nome), []byte(corpo), 0o644); err != nil {
		t.Fatal(err)
	}
}

// O CONTEÚDO da spec entra no compilado — não um link para ela.
//
// É a razão de o mecanismo existir: uma documentação que só aponta para os arquivos deixa
// de ser documentação e vira indexação, e quem lê precisa abrir link por link.
func TestBuild_conteudoDaSpecEntraNoCompilado(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "infra.md.tmpl",
		`{{range specs "layer=infra"}}## {{.Titulo}}

{{section . "Visão Geral"}}
{{end}}`)

	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(filepath.Join(root, OutDir, "infra.md"))
	if err != nil {
		t.Fatalf("nada em `%s/infra.md`: %v", OutDir, err)
	}
	out := string(b)
	if !strings.Contains(out, "A régua confronta o que foi prometido") {
		t.Errorf("o corpo da seção não entrou — o compilado virou índice:\n%s", out)
	}
	if !strings.Contains(out, "GoLiveChecklist — a régua") {
		t.Errorf("o título não entrou:\n%s", out)
	}
}

// A camada vem do HEADER da spec, não do mapa.
//
// Medido no projeto de referência: as 84 specs têm `layer: spec` no mapa (casam
// `**/*.spec.md`), e a camada real — screen, lambdas, infra — está no header. Agrupar pela
// do mapa daria UM grupo com tudo, que é o oposto da visão por camada.
func TestSpecs_camadaVemDoHeaderDaSpec(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.fnSpecs("layer=infra")
	if err != nil || len(got) != 1 {
		t.Fatalf("filtro por `infra` (do header) achou %d (%v), queria 1", len(got), err)
	}
	if _, err := c.fnSpecs("layer=spec"); err == nil {
		t.Error("filtrar pela camada do MAPA casou algo — a camada do header é que vale")
	}
}

// A seção `## Regras` contém `###` de regra, e o recorte não pode parar no primeiro deles.
func TestSection_naoParaNoSubHeading(t *testing.T) {
	s := Spec{raw: specDeExemplo}
	got := fnSection(s, "Regras")
	if !strings.Contains(got, "GLCGL-B02") {
		t.Errorf("o recorte parou antes da segunda regra:\n%s", got)
	}
	if strings.Contains(got, "GLCGL-I01") {
		t.Errorf("o recorte invadiu `## Invariantes`:\n%s", got)
	}
}

// O `rules` devolve as regras SOLTAS — a visão "todas as regras do sistema" é uma lista de
// itens, não a concatenação de seções `## Regras`.
func TestRules_devolveCadaRegraSeparada(t *testing.T) {
	rs := fnRules(Spec{raw: specDeExemplo})
	if len(rs) != 3 {
		t.Fatalf("achou %d regras, queria 3 (B01, B02, I01)", len(rs))
	}
	if rs[0].Code != "GLCGL-B01" || rs[0].Titulo != "item sem artefato não passa" {
		t.Errorf("primeira regra = %+v", rs[0])
	}
	if strings.Contains(rs[1].Corpo, "Invariantes") {
		t.Errorf("o corpo de B02 invadiu o heading seguinte: %q", rs[1].Corpo)
	}
	if rs[2].Code != "GLCGL-I01" {
		t.Errorf("invariante não entrou como regra: %+v", rs[2])
	}
}

// O compilado carrega o MARCADOR: quem abre sabe que editar é perder o trabalho.
func TestBuild_escreveOMarcador(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", "nada")
	c, _ := New(root, g)
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(root, OutDir, "x.md"))
	if !markerRE.Match(b) {
		t.Errorf("o compilado saiu sem marcador:\n%s", b)
	}
	if !strings.Contains(string(b), "doct/x.md.tmpl") {
		t.Errorf("o marcador não diz de onde veio:\n%s", b)
	}
}

// A PROTEÇÃO que evita apagar trabalho: um `.md` escrito à mão não é sobrescrito porque
// alguém criou um template com o mesmo nome.
//
// O caso é real e previsível: a doc de produto é escrita direto em `docs/`, sem template.
// Um `doct/produto.md.tmpl` criado por engano destruiria o arquivo inteiro, e o conteúdo
// não estaria em spec nenhuma para ser recompilado.
func TestBuild_naoSobrescreveDocEscritaAMao(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "produto.md.tmpl", "gerado")

	docs := filepath.Join(root, OutDir)
	os.MkdirAll(docs, 0o755)
	aMao := "# Produto\n\nEscrito à mão, com conteúdo que não está em spec nenhuma.\n"
	os.WriteFile(filepath.Join(docs, "produto.md"), []byte(aMao), 0o644)

	c, _ := New(root, g)
	res, err := c.Build(false)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := os.ReadFile(filepath.Join(docs, "produto.md"))
	if string(b) != aMao {
		t.Fatalf("o compilador apagou trabalho escrito à mão:\n%s", b)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "produto.md" {
		t.Errorf("o arquivo foi pulado em SILÊNCIO: ignorados = %v", res.Skipped)
	}
}

// `--dry-run` não escreve: é o modo que o gate usa para comparar sem tocar no disco.
func TestBuild_dryRunNaoEscreve(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", "nada")
	c, _ := New(root, g)
	if _, err := c.Build(true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, OutDir, "x.md")); err == nil {
		t.Error("o dry-run escreveu no disco")
	}
}

// Template quebrado FALHA: um erro de sintaxe não pode virar um doc pela metade.
func TestBuild_templateQuebradoFalha(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", `{{range specs "layer=infra"}}sem end`)
	c, _ := New(root, g)
	if _, err := c.Build(false); err == nil {
		t.Error("template sem `end` compilou sem erro")
	}
}

// `layers` só lista camada que TEM spec — uma lista com camadas vazias mandaria quem lê
// procurar seção que não existe.
func TestLayers_soAsQueTemSpec(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	c, _ := New(root, g)
	got := c.fnLayers()
	if len(got) != 1 || got[0] != "infra" {
		t.Errorf("layers = %v, queria [infra]", got)
	}
}

// O FILTRO ERRADO PARA O BUILD. Sem isso a falha seria invisível: o documento compila
// verde, com título e com as outras seções, e a que interessa simplesmente não está lá.
// Ninguém procura o que não sabe que falta.
func TestSpecs_filtroErradoFalha(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	c, _ := New(root, g)

	casos := []struct{ nome, filtro, esperaNaMsg string }{
		{"campo com typo", "layar=infra", "unknown field"},
		{"camada inexistente", "layer=screen", "existing"},
		{"sem o `=`", "infra", "`=`"},
		{"código inexistente", "code=NAOEX", "NAOEX"},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			_, err := c.fnSpecs(caso.filtro)
			if err == nil {
				t.Fatalf("filtro %q passou em silêncio — a seção sairia vazia", caso.filtro)
			}
			if !strings.Contains(err.Error(), caso.esperaNaMsg) {
				t.Errorf("a mensagem não ajuda a consertar: %v", err)
			}
		})
	}
}

// O erro do filtro sobe até o BUILD e o derruba — não vira um documento pela metade.
func TestBuild_filtroErradoDerrubaOBuild(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", `{{range specs "layer=naoexiste"}}{{.Code}}{{end}}`)
	c, _ := New(root, g)
	if _, err := c.Build(false); err == nil {
		t.Fatal("o build passou com filtro que não casa nada")
	}
	if _, err := os.Stat(filepath.Join(root, OutDir, "x.md")); err == nil {
		t.Error("escreveu um documento incompleto antes de falhar")
	}
}

// Spec no MAPA e não no disco para o build: ela sairia da documentação sem uma palavra, e
// o compilado ficaria verde sem ela.
func TestNew_specNoMapaSemArquivoFalha(t *testing.T) {
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/Fantasma.spec.md", Kind: mapx.KindSpec, Code: "FNTSM"},
	}}
	if _, err := New(root, g); err == nil {
		t.Fatal("spec ausente do disco passou em silêncio")
	}
}

// O RÓTULO DO LINK NÃO CARREGA A DISPENSA.
//
// O heading de uma regra pode trazer o marcador na própria linha — e o gate o aceita ali de
// propósito (`FRMTT-I02`: à direita do símbolo é onde o formatador não o move). Mas o
// compilado usa esse texto como rótulo entre colchetes, e o resultado era ilegível:
//
//   - [DSHBR-B01 — a escolha de fixar <!-- @no-mark: ... -->](camadas/shell.md#...)
//
// MEDIDO no projeto de referência (#647): 57 ocorrências no `docs/regras.md`, espalhadas
// por todas as camadas — não é defeito de uma spec, é do gerador.
func TestTituloDaRegraNaoLevaComentarioHTML(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{
		"pkg/Com.spec.md": `<!-- @anchors
code: ABCDE
-->

# Com

## Regras

### ABCDE-B01 — a regra diz isto <!-- @no-mark: não há símbolo a marcar -->

O corpo da regra.
`,
	})
	escreveTemplate(t, root, "r.md.tmpl",
		`{{range specs "layer=spec"}}{{range rules .}}- [{{.Code}} — {{.Titulo}}](x)
{{end}}{{end}}`)

	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, OutDir, "r.md"))
	if err != nil {
		t.Fatal(err)
	}
	out := string(b)

	// SÓ A LINHA DO LINK. O compilado SEMPRE abre com `<!-- anchors:generated ... -->`, e
	// procurar `<!--` no arquivo inteiro casaria esse cabeçalho — a primeira versão desta
	// asserção reprovava a correção certa por isso.
	var linha string
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, "- [ABCDE-B01") {
			linha = l
			break
		}
	}
	if linha == "" {
		t.Fatalf("não achei a linha do link no compilado:\n%s", out)
	}
	if strings.Contains(linha, "<!--") || strings.Contains(linha, "-->") {
		t.Errorf("o rótulo do link levou o comentário de dispensa:\n  %s", linha)
	}
	// E O TÍTULO TEM DE SOBREVIVER: cortar o comentário não pode cortar a regra junto.
	if !strings.Contains(linha, "a regra diz isto") {
		t.Errorf("o título da regra sumiu do rótulo — o corte levou o texto junto:\n  %s", linha)
	}
}

// A memorizacao do `fnSize` e' otimizacao pura: a saida tem de ser IDENTICA com e sem
// cache. Um teste de comportamento nao a pegaria (e' o mesmo resultado por definicao), e
// e' justamente por isso que ela precisa de guarda propria — o risco nao e' o cache
// errar, e' ele devolver a resposta de OUTRO recorte.
//
// Confronta os dois caminhos no mesmo compilador: a primeira chamada calcula e memoriza,
// a segunda le do cache, e recortes diferentes tem de continuar diferentes.
func TestFnSize_cachePorRecorte(t *testing.T) {
	root, g := projetoDoisRecortes(t)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	a1, err := c.fnSize("layer=alfa")
	if err != nil {
		t.Fatal(err)
	}
	b1, err := c.fnSize("layer=beta")
	if err != nil {
		t.Fatal(err)
	}
	if a1.Units == b1.Units {
		t.Fatalf("o cenario de teste nao distingue os recortes: alfa=%d beta=%d", a1.Units, b1.Units)
	}
	// Segunda rodada: agora vem do cache, e tem de bater com a primeira.
	a2, _ := c.fnSize("layer=alfa")
	b2, _ := c.fnSize("layer=beta")
	if a2 != a1 {
		t.Errorf("recorte alfa mudou ao vir do cache: %+v vs %+v", a2, a1)
	}
	if b2 != b1 {
		t.Errorf("recorte beta mudou ao vir do cache: %+v vs %+v", b2, b1)
	}
}

// projetoDoisRecortes monta duas camadas com contagens DIFERENTES, para que um cache que
// confundisse os recortes fosse visivel.
func projetoDoisRecortes(t *testing.T) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	spec := func(code, layer string) string {
		return "---\ncode: " + code + "\nlayer: " + layer + "\n---\n\n# T — t\n\n## Visão Geral\n\nx.\n"
	}
	os.WriteFile(filepath.Join(root, "pkg/A.spec.md"), []byte(spec("AAAAA", "alfa")), 0o644)
	os.WriteFile(filepath.Join(root, "pkg/B.spec.md"), []byte(spec("BBBBB", "beta")), 0o644)
	os.WriteFile(filepath.Join(root, "pkg/C.spec.md"), []byte(spec("CCCCC", "beta")), 0o644)
	return root, &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/A.spec.md", Kind: mapx.KindSpec, Code: "AAAAA", Layer: "alfa"},
		{ID: "pkg/B.spec.md", Kind: mapx.KindSpec, Code: "BBBBB", Layer: "beta"},
		{ID: "pkg/C.spec.md", Kind: mapx.KindSpec, Code: "CCCCC", Layer: "beta"},
	}}
}

// The header's date is not an input of the document: the pre-commit dates staged headers
// before the gates run, and hashing the date failed `docs-fresh` on docs just built.
func TestInputsHashIgnoresTheHeaderDate(t *testing.T) {
	hash := func(raw string) string {
		c := &Compiler{specs: []Spec{{Path: "a.spec.md", raw: raw}}}
		return c.inputsHash([]byte("tmpl"))
	}
	spec := "<!-- @anchors\n  code: AAAAA\n  updated_at: 2026-09-01\n-->\n# A\n\nbody\n"
	base := hash(spec)
	if got := hash(strings.Replace(spec, "2026-09-01", "2026-09-26", 1)); got != base {
		t.Errorf("only the header date changed and the docs were declared out of date")
	}
	if got := hash(strings.Replace(spec, "body", "new body", 1)); got == base {
		t.Error("the content changed and the hash did not")
	}
	// Below the header an `updated_at:` line is content, not the header's date.
	deep := func(date string) string { return spec + strings.Repeat("x\n", 12) + "updated_at: " + date + "\n" }
	if hash(deep("2026-09-01")) == hash(deep("2026-09-26")) {
		t.Error("an updated_at line far below the header is content and must change the hash")
	}
}

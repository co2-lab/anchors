package ops

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// O COMPILADO SAI DA ÁRVORE, NÃO DO MAPA EM DISCO.
//
// `docs build` compila a partir do MAPA. Ele o lia como estava, e um mapa que não conhece
// uma unidade produz um compilado SEM os cenários dela — sem erro, sem aviso: o arquivo
// encolhe e o commit parece normal.
//
// MEDIDO no projeto de referência (#717, #743): o compilado no `develop` tinha 687 entradas
// contra as 718 que as specs produzem. Faltavam quatro unidades inteiras — DTSTD 11,
// SRMTS 8, NTCNN 8, HLCHH 4 —, e as specs das quatro estavam na árvore havia mais de uma
// semana quando o compilado foi reescrito sem elas.
//
// A causa a montante era o mapa daquele clone, mutilado por um merge (o driver perdia nós
// do lado incoming — corrigido na v0.1.129). Este comando não podia depender de o mapa
// estar íntegro: ele reconstrói da árvore e preserva os carimbos, que são estado de
// trabalho e não derivado.
//
// `--no-map-rebuild` mantém o comportamento antigo para quem PRECISA de um mapa específico.
// É opt-out porque o modo perigoso tem de ser o pedido: quem não sabe que a garantia existe
// recebe a garantia.
func TestDocsBuildReconstroiOMapaEmVezDeLerODisco(t *testing.T) {
	b, err := os.ReadFile("docs.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	// A RECONSTRUÇÃO: o mesmo caminho do `map build` — scan da árvore, grafo novo.
	for _, peca := range []string{"scan.Walk", "mapx.Build"} {
		if !strings.Contains(s, peca) {
			t.Errorf("o `docs build` não chama %q — sem reconstruir, ele compila o que o "+
				"mapa em disco souber, e o que faltar lá some do compilado em silêncio", peca)
		}
	}
	// OS CARIMBOS não podem ser perdidos na reconstrução: o laudo de cada julgamento vive
	// no `--reason` do comando que o gravou, não no arquivo.
	if !strings.Contains(s, "mapx.PreserveStamps") {
		t.Error("a reconstrução não preserva os carimbos — refazer um julgamento " +
			"adversarial é caro, e o mapa voltaria a acusar tudo como nunca validado")
	}
	// O OPT-OUT tem de existir, e ser opt-out: o padrão é a garantia.
	if !strings.Contains(s, "no-map-rebuild") {
		t.Error("falta `--no-map-rebuild`: compilar contra um mapa específico é legítimo " +
			"(conferir uma revisão antiga, um `--map` de outro lugar)")
	}
	if strings.Contains(s, `"no-map-rebuild", true`) {
		t.Error("`--no-map-rebuild` está ligado por padrão — o modo que produziu o " +
			"compilado mutilado tem de ser o PEDIDO, nunca o herdado")
	}
}

// A MESMA GARANTIA, MEDIDA PELO RESULTADO.
//
// O teste acima confronta o FONTE — que as chamadas estão lá. Isto confronta o COMPILADO:
// com um mapa em disco que não conhece a unidade, o cenário dela tem de aparecer assim
// mesmo, e com `--no-map-rebuild` NÃO pode aparecer.
//
// Os dois juntos porque cada um sozinho mente de um jeito: o de fonte passaria se as
// chamadas existissem num ramo morto, e este passaria se alguém trocasse a reconstrução por
// qualquer outra coisa que produzisse o mesmo arquivo.
func TestDocsBuildCompilaAUnidadeQueOMapaNaoConhece(t *testing.T) {
	root := t.TempDir()

	escreve := func(rel, conteudo string) {
		t.Helper()
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	escreve("anchors.yaml", "lang: pt-BR\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	escreve("pkg/Coisa.spec.md", `<!-- @anchors
code: ABCDE
-->

# Coisa

## Visão Geral

O que a unidade faz.
`)
	// O MAPA EM DISCO NÃO CONHECE A SPEC — é exatamente o estado que produziu o compilado
	// com 31 entradas a menos: um merge tinha perdido os nós, e nada acusava.
	escreve("anchors.graph.yaml", "version: 3\nnodes: []\nedges: []\n")
	escreve(filepath.Join("doct", "camadas.md.tmpl"),
		`{{range specs "layer=spec"}}## {{.Titulo}}

{{section . "Visão Geral"}}
{{end}}`)

	compila := func(args ...string) string {
		t.Helper()
		cmd := newDocsBuildCmd()
		cmd.SetArgs(append([]string{"--root", root}, args...))
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("docs build %v: %v", args, err)
		}
		b, err := os.ReadFile(filepath.Join(root, "docs", "camadas.md"))
		if err != nil {
			t.Fatalf("nada compilado em docs/camadas.md: %v", err)
		}
		return string(b)
	}

	if out := compila(); !strings.Contains(out, "O que a unidade faz") {
		t.Errorf("o compilado NÃO tem o cenário da spec que o mapa em disco desconhece — "+
			"é o defeito medido: 687 entradas onde as specs produziam 718.\ncompilado:\n%s", out)
	}

	// E O OPT-OUT tem de continuar valendo: quem pede o mapa em disco recebe o mapa em
	// disco, com o que ele não conhece de fora. Sem esta metade, `--no-map-rebuild` seria
	// uma flag que não faz nada — e o teste acima passaria mesmo se a reconstrução fosse
	// incondicional.
	//
	// Aqui o mapa está VAZIO, então "a spec fica de fora" se manifesta como erro do
	// template ("nenhuma spec na camada"), e não como compilado sem a linha. As duas
	// formas são o mesmo fato: o que o mapa não conhece não chega ao compilado.
	semRebuild := newDocsBuildCmd()
	semRebuild.SetArgs([]string{"--root", root, "--no-map-rebuild"})
	semRebuild.SetOut(io.Discard)
	semRebuild.SetErr(io.Discard)
	err := semRebuild.Execute()
	if err == nil {
		b, _ := os.ReadFile(filepath.Join(root, "docs", "camadas.md"))
		if strings.Contains(string(b), "O que a unidade faz") {
			t.Error("com `--no-map-rebuild` o compilado trouxe a spec que o mapa em disco " +
				"NÃO tem — a flag não está sendo respeitada, e quem precisa de um mapa " +
				"específico recebe outro")
		}
		return
	}
	if !strings.Contains(err.Error(), "no spec in layer") {
		t.Errorf("com `--no-map-rebuild` e mapa vazio, o erro devia ser a ausência da spec "+
			"— veio outro: %v", err)
	}
}

// --- duties, init, and what `build` reports ---

const dutiesConfig = `layers:
  api:
    pattern: "api/**/*.go"
    kind: code
docs:
  required:
    - kind: openapi
      path: docs/api.yaml
      trigger: [api]
      why: the contract of the endpoints
    - kind: c4
      path: docs/c4.md
`

func TestDocsDutiesAnswersPerLayerAndPerUnit(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", dutiesConfig)
	writeFile(t, root, "api/handler.go", "package api\n")

	err, out := runCmd(t, newDocsCmd(), "duties", "--root", root)
	if err != nil || !strings.Contains(out, "Required documentation:") ||
		!strings.Contains(out, "docs/api.yaml") || !strings.Contains(out, "docs/c4.md") {
		t.Errorf("all duties: %v\n%s", err, out)
	}

	err, out = runCmd(t, newDocsCmd(), "duties", "--root", root, "--layer", "api")
	if err != nil || !strings.Contains(out, "Changing `api` requires touching") ||
		!strings.Contains(out, "docs/api.yaml") || strings.Contains(out, "docs/c4.md") {
		t.Errorf("--layer api must name only the triggered doc: %v\n%s", err, out)
	}

	err, out = runCmd(t, newDocsCmd(), "duties", "--root", root, "--layer", "web")
	if err != nil || !strings.Contains(out, "No documentation is required when changing `web`") {
		t.Errorf("--layer web: %v\n%s", err, out)
	}

	// The UNIT is resolved to its layer: the handler lives in `api`.
	err, out = runCmd(t, newDocsCmd(), "duties", "--root", root, "--unit", "api/handler.go")
	if err != nil || !strings.Contains(out, "Changing `api/handler.go` requires touching") ||
		!strings.Contains(out, "docs/api.yaml") {
		t.Errorf("--unit: %v\n%s", err, out)
	}
}

func TestDocsDutiesWithoutDeclarationTeachesTheKinds(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "version: 1\n")
	err, out := runCmd(t, newDocsCmd(), "duties", "--root", root)
	if err != nil || !strings.Contains(out, "declares no required documentation") || !strings.Contains(out, "openapi, c4") {
		t.Errorf("no docs declared: %v\n%s", err, out)
	}
	if err, _ := runCmd(t, newDocsCmd(), "duties", "--root", t.TempDir()); err == nil {
		t.Error("duties without anchors.yaml must fail")
	}
}

// `docs init` writes the skeleton once; a second run skips what exists, --force rewrites.
func TestDocsInitWritesTheSkeletonOnce(t *testing.T) {
	root := t.TempDir()
	if err, _ := runCmd(t, newDocsCmd(), "init", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("init without a map must point at `map build`: %v", err)
	}
	writeFile(t, root, "anchors.graph.yaml", "version: 4\nnodes: []\nedges: []\n")

	err, out := runCmd(t, newDocsCmd(), "init", "--root", root)
	if err != nil {
		t.Fatalf("docs init: %v", err)
	}
	tmpls, _ := filepath.Glob(filepath.Join(root, "doct", "*.md.tmpl"))
	if len(tmpls) == 0 || !strings.Contains(out, "✓ doct/") {
		t.Fatalf("no template written:\n%s", out)
	}
	edited := tmpls[0]
	if err := os.WriteFile(edited, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, out = runCmd(t, newDocsCmd(), "init", "--root", root)
	if !strings.Contains(out, "already exists (use --force to overwrite)") {
		t.Errorf("the second run does not say what it skipped:\n%s", out)
	}
	if b, _ := os.ReadFile(edited); string(b) != "mine" {
		t.Error("a second run without --force overwrote an edited template")
	}
	runCmd(t, newDocsCmd(), "init", "--root", root, "--force")
	if b, _ := os.ReadFile(edited); string(b) == "mine" {
		t.Error("--force did not rewrite the template")
	}
}

// What `build` reports: a hand-written page with a template is SKIPPED and named; a dry
// run says nothing was written; with no template, it says there was nothing to compile.
func TestDocsBuildReportsSkippedDryRunAndNothing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "version: 1\n")
	if err := os.MkdirAll(filepath.Join(root, "doct"), 0o755); err != nil {
		t.Fatal(err)
	}
	err, out := runCmd(t, newDocsBuildCmd(), "--root", root)
	if err != nil || !strings.Contains(out, "no template in `doct/`") {
		t.Errorf("no template: %v\n%s", err, out)
	}

	writeFile(t, root, "doct/guia.md.tmpl", "# Guide\n")
	writeFile(t, root, "doct/manual.md.tmpl", "# Manual\n")
	writeFile(t, root, "docs/manual.md", "written by hand\n")

	err, out = runCmd(t, newDocsBuildCmd(), "--root", root, "--dry-run")
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}
	if !strings.Contains(out, "docs/guia.md compiled (not written)") {
		t.Errorf("the dry run does not list the page:\n%s", out)
	}
	if _, serr := os.Stat(filepath.Join(root, "docs", "guia.md")); serr == nil {
		t.Error("--dry-run wrote the page")
	}
	if !strings.Contains(out, "1 file(s) NOT generated") || !strings.Contains(out, "docs/manual.md") {
		t.Errorf("the hand-written page is not reported as skipped:\n%s", out)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "docs", "manual.md")); string(b) != "written by hand\n" {
		t.Error("the hand-written page was overwritten")
	}

	if err, _ := runCmd(t, newDocsBuildCmd(), "--root", root, "--no-map-rebuild"); err == nil ||
		!strings.Contains(err.Error(), "load map") {
		t.Errorf("--no-map-rebuild without a map must say so: %v", err)
	}
}

func TestDocsFreshHintNamesEveryStalePage(t *testing.T) {
	got := docsFreshHint([]string{"docs/a.md", "docs/b.md"})
	if !strings.Contains(got, "2 documento(s)") || !strings.Contains(got, "docs/a.md, docs/b.md") ||
		!strings.Contains(got, "anchors docs build") {
		t.Errorf("hint = %q", got)
	}
}

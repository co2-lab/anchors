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

// THE SAME GUARANTEE, MEASURED BY THE RESULT.
//
// The test above confronts the SOURCE — that the calls are there. This one confronts the
// COMPILED output: with a map on disk that does not know the unit, its scenario must appear
// anyway, and with `--no-map-rebuild` it must NOT.
//
// Both together, because each alone lies in its own way: the source test would pass if the
// calls lived in a dead branch, and this one would pass if someone swapped the rebuild for
// anything else producing the same file.
func TestDocsBuildCompilesTheUnitTheMapDoesNotKnow(t *testing.T) {
	t.Run("DCCMD-B01: Build compiles from the tree, even what the map on disk does not know", func(t *testing.T) {})
	t.Run("DCCMD-B02: No-map-rebuild compiles against the map on disk and warns", func(t *testing.T) {})
	t.Run("DCCMD-X01: Build never writes the map", func(t *testing.T) {})
	root := t.TempDir()

	write := func(rel, content string) {
		t.Helper()
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("anchors.yaml", "lang: pt-BR\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	write("pkg/Coisa.spec.md", `<!-- @anchors
code: ABCDE
-->

# Coisa

## Visão Geral

O que a unidade faz.
`)
	// THE MAP ON DISK DOES NOT KNOW THE SPEC — exactly the state that produced the compiled
	// output with 31 entries missing: a merge had lost the nodes, and nothing complained.
	const emptyMap = "version: 3\nnodes: []\nedges: []\n"
	write("anchors.graph.yaml", emptyMap)
	write(filepath.Join("doct", "camadas.md.tmpl"),
		`{{range specs "layer=spec"}}## {{.Titulo}}

{{section . "Visão Geral"}}
{{end}}`)

	compile := func(args ...string) string {
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
			t.Fatalf("nothing compiled into docs/camadas.md: %v", err)
		}
		return string(b)
	}

	if out := compile(); !strings.Contains(out, "O que a unidade faz") {
		t.Errorf("the compiled output does NOT have the scenario of the spec the map on disk does not know — "+
			"the measured defect: 687 entries where the specs produced 718.\ncompiled:\n%s", out)
	}
	// The rebuild stays in memory: the map on disk is the map command's to write.
	if b, _ := os.ReadFile(filepath.Join(root, "anchors.graph.yaml")); string(b) != emptyMap {
		t.Errorf("docs build wrote the map:\n%s", b)
	}

	// AND THE OPT-OUT must keep working: whoever asks for the map on disk gets the map on
	// disk, with what it does not know left out. Without this half, `--no-map-rebuild` would
	// be a flag that does nothing — and the check above would pass even if the rebuild were
	// unconditional.
	//
	// Here the map is EMPTY, so "the spec is left out" shows up as a template error ("no spec
	// in layer"), not as compiled output without the line. Both are the same fact: what the
	// map does not know does not reach the compiled output.
	noRebuild := newDocsBuildCmd()
	noRebuild.SetArgs([]string{"--root", root, "--no-map-rebuild"})
	noRebuild.SetOut(io.Discard)
	var stderr strings.Builder
	noRebuild.SetErr(&stderr)
	err := noRebuild.Execute()
	if !strings.Contains(stderr.String(), "the compiled output comes from the map on disk") {
		t.Errorf("--no-map-rebuild does not warn on stderr:\n%s", stderr.String())
	}
	if err == nil {
		b, _ := os.ReadFile(filepath.Join(root, "docs", "camadas.md"))
		if strings.Contains(string(b), "O que a unidade faz") {
			t.Error("with `--no-map-rebuild` the compiled output brought the spec the map on disk " +
				"does NOT have — the flag is not honored, and whoever needs a specific map gets another")
		}
		return
	}
	if !strings.Contains(err.Error(), "no spec in layer") {
		t.Errorf("with `--no-map-rebuild` and an empty map, the error should be the missing spec "+
			"— got another: %v", err)
	}
}

// Without anchors.yaml there is no tree to rebuild the map from, and the error says how to
// create one.
func TestDocsBuildWithoutConfigPointsAtInit(t *testing.T) {
	t.Run("DCCMD-E01: Build without a config points at init", func(t *testing.T) {})
	err, _ := runCmd(t, newDocsBuildCmd(), "--root", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "anchors init") {
		t.Errorf("want the load error pointing at `anchors init`, got %v", err)
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
	t.Run("DCCMD-B07: Duties answers for the project, a layer or a unit", func(t *testing.T) {})
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
	t.Run("DCCMD-B08: Duties without a declaration teaches the known kinds", func(t *testing.T) {})
	t.Run("DCCMD-E02: Duties without a config fails", func(t *testing.T) {})
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
	t.Run("DCCMD-B06: Init writes the skeleton once and needs a map", func(t *testing.T) {})
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
	t.Run("DCCMD-B03: A dry run compiles without writing", func(t *testing.T) {})
	t.Run("DCCMD-B04: A hand-written page is skipped, kept and named", func(t *testing.T) {})
	t.Run("DCCMD-B05: Build with no template says there is nothing to compile", func(t *testing.T) {})
	t.Run("DCCMD-E03: No-map-rebuild without a map says so", func(t *testing.T) {})
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

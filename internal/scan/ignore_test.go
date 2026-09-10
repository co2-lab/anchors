package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// TestAncoragemDoGitignore guarda o bug que apagou 129 unidades do mapa em silêncio.
//
// A linha `/data/` do projeto significa, para o git, "o diretório `data` NA RAIZ". Ao
// descartar a barra inicial, o Anchors a lia como "qualquer segmento chamado `data`" — e
// `amplify/data/models/` inteiro (as specs e os modelos do schema) saía da varredura. O
// mapa não ficava errado com barulho: ficava menor, e nada acusava a ausência.
func TestAncoragemDoGitignore(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(strings.Join([]string{
		"/data/",       // ancorado: só na raiz
		"node_modules", // sem barra: qualquer nível
		"*.log",
		"probe*.ts",
		"docs/build", // barra no meio: ancorado
		"docs/saida", // idem, com nome fora da lista universal
	}, "\n")), 0o644))

	ig := LoadIgnore(dir)
	casos := []struct {
		rel   string
		isDir bool
		quer  bool
		por   string
	}{
		{"data", true, true, "`/data/` casa o diretório na raiz"},
		{"data/dump.json", false, true, "abaixo de um diretório ignorado"},
		{"amplify/data", true, false, "`/data/` é ANCORADO — não casa `data` aninhado"},
		{"amplify/data/models/AiUsage.ts", false, false, "o caso real que sumiu do mapa"},
		{"a/node_modules", true, true, "sem barra casa em qualquer nível"},
		{"a/b/c.log", false, true, "`*.log` casa o basename em qualquer nível"},
		{"probe1.ts", false, true, "a sonda de revisor que virava task"},
		{"docs/build", true, true, "barra no meio: ancorado, e casa"},
		// `apps/docs/build` NÃO é caso de ancoragem: `build` está na lista universal, que
		// vale em qualquer projeto e é aplicada antes do `.gitignore`. Um par aninhado que
		// prove ancoragem precisa de um nome fora dessa lista.
		{"docs/saida", true, true, "ancorado, na raiz: casa"},
		{"apps/docs/saida", true, false, "ancorado: `docs/saida` não casa aninhado"},
	}
	for _, c := range casos {
		var got bool
		if c.isDir {
			got = ig.SkipDir(filepath.Base(c.rel), c.rel)
		} else {
			got = ig.SkipFile(c.rel)
		}
		if got != c.quer {
			t.Errorf("%q (dir=%v): ignorado=%v, esperava %v — %s", c.rel, c.isDir, got, c.quer, c.por)
		}
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// TestListaEmbutidaEhDerrotavel: `build`/`dist` são saída de compilação na maioria dos
// projetos e PASTA DE CÓDIGO em alguns. Uma lista embutida que não pode ser contestada é o
// framework decidindo, por um projeto que não conhece, o que nele é descartável — e o erro
// é silencioso: a camada some do mapa e nenhum gate acusa a ausência.
func TestListaEmbutidaEhDerrotavel(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.log\n"), 0o644))

	// sem declaração: o default vale
	if !LoadIgnore(dir).SkipDir("build", "build") {
		t.Error("sem declaração, `build` deve seguir o default (ignorado)")
	}

	// a ESTRUTURA declara que ali há código
	cfg := &config.Config{Layers: map[string]config.Layer{
		"core": {Pattern: "build/**/*.ts", Kind: "code"},
	}}
	if LoadIgnoreFor(dir, cfg).SkipDir("build", "build") {
		t.Error("camada apontando para dentro de `build` deve DERROTAR o default")
	}

	// um catch-all NÃO reabilita — senão `**/*.ts` traria `node_modules` de volta
	amplo := &config.Config{Layers: map[string]config.Layer{
		"tudo": {Pattern: "**/*.ts", Kind: "code"},
	}}
	if !LoadIgnoreFor(dir, amplo).SkipDir("node_modules", "node_modules") {
		t.Error("catch-all não pode reabilitar `node_modules`")
	}

	// a negação no `.gitignore` também derrota
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("!build/\n"), 0o644))
	if LoadIgnore(dir).SkipDir("build", "build") {
		t.Error("`!build/` no .gitignore deve derrotar o default")
	}

	// a maquinaria NÃO é derrotável, sob nenhuma declaração
	must(t, os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("!.git/\n"), 0o644))
	if !LoadIgnore(dir).SkipDir(".git", ".git") {
		t.Error("`.git` é maquinaria — nenhuma declaração o reabilita")
	}
}

// TestEfemerosNaoViramTrabalho guarda o ruído medido num E2E real: o watcher enfileirou
// task para `amplify/data/.!21662!resource.spec.md` — arquivo que o editor cria por
// milissegundos durante um salvamento atômico e que casa o glob de spec. Três tasks
// tiveram de ser descartadas à mão.
//
// O `.gitignore` do projeto não cobre isso, e não deveria: não é decisão do projeto, é
// ruído do sistema de arquivos.
func TestEfemerosNaoViramTrabalho(t *testing.T) {
	ig := LoadIgnore(t.TempDir())
	ruido := []string{
		"amplify/data/.!21662!resource.spec.md", // o caso real
		"src/a.ts.swp", "src/a.ts~", "src/.#a.ts", "src/#a.ts#",
		"src/x.tmp", ".DS_Store",
	}
	for _, r := range ruido {
		if !ig.SkipFile(r) {
			t.Errorf("%q é ruído de editor/sistema e não pode virar trabalho", r)
		}
	}
	material := []string{
		"src/a.ts", "src/a.spec.md", "amplify/data/resource.spec.md",
		"src/tmp/util.ts", // `tmp` no CAMINHO não é `.tmp` no nome
	}
	for _, m := range material {
		if ig.SkipFile(m) {
			t.Errorf("%q é material do projeto e foi descartado", m)
		}
	}
}

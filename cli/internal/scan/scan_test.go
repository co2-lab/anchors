package scan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestExtractCodesIgnoresComments(t *testing.T) {
	src := `// vê o cenário FOOOX-VR de outra tela
const x = 1  // BARRX-S01 também aqui
export const y = "BAZZX-B01"  // este está em código (string), conta`
	codes := extractCodes([]byte(src))
	has := func(c string) bool {
		for _, x := range codes {
			if x == c {
				return true
			}
		}
		return false
	}
	if has("FOOOX-VR") || has("BARRX-S01") {
		t.Errorf("códigos em comentário não deveriam contar: %v", codes)
	}
	if !has("BAZZX-B01") {
		t.Errorf("código em código (fora de comentário) deveria contar: %v", codes)
	}
}

func TestStripBlockComment(t *testing.T) {
	src := "/* AAAAX-S01 no bloco */ const z = \"BBBBX-V01\""
	codes := extractCodes([]byte(src))
	for _, c := range codes {
		if c == "AAAAX-S01" {
			t.Error("código em comentário de bloco não deveria contar")
		}
	}
}

func TestExtractDepsParsesTable(t *testing.T) {
	// cria os arquivos-alvo no tempdir (resolveDepPath faz os.Stat)
	root := t.TempDir()
	mk := func(rel string) {
		p := root + "/" + rel
		if err := writeDeep(p, "x"); err != nil {
			t.Fatal(err)
		}
	}
	mk("apps/mobile/src/stores/auth.store.ts")
	mk("apps/mobile/src/hooks/useAuth.ts")

	spec := `# LoginScreen

> ` + "`" + `code: LOGIX` + "`" + `

## Dependências de Dados

| Cód  | Arquivo                  | Método         | Camada |
| ---- | ------------------------ | -------------- | ------ |
| DEP1 | ` + "`stores/auth.store.ts`" + ` | ` + "`useAuthStore`" + ` | store  |
| DEP2 | ` + "`hooks/useAuth.ts`" + `     | ` + "`signIn`" + `       | hook   |

## Data Contract

| Campo | Origem | Obrigatório |
| ----- | ------ | ----------- |
| ` + "`isLoading`" + ` | DEP1 | ✅ |
`
	specRel := "apps/mobile/src/features/auth/screens/LoginScreen.spec.md"
	deps := extractDeps([]byte(spec), root, specRel)
	if len(deps) != 2 {
		t.Fatalf("esperava 2 deps, veio %d: %+v", len(deps), deps)
	}
	// Method PRESERVA as crases do autor (o sinal de "isto é um SÍMBOLO", não prosa);
	// as demais colunas seguem sem markdown inline.
	if deps[0].Code != "DEP1" || deps[0].Method != "`useAuthStore`" || deps[0].Layer != "store" {
		t.Errorf("DEP1 mal parseada: %+v", deps[0])
	}
	// resolveDepPath deve ter resolvido para o caminho relativo à raiz (via /src/ da spec)
	if deps[0].File != "apps/mobile/src/stores/auth.store.ts" {
		t.Errorf("DEP1.File não resolveu p/ raiz: %q", deps[0].File)
	}
	if deps[1].Code != "DEP2" || deps[1].File != "apps/mobile/src/hooks/useAuth.ts" {
		t.Errorf("DEP2 mal parseada: %+v", deps[1])
	}
}

func TestResolveDepPathAcceptsSrcPrefixed(t *testing.T) {
	// o autor escreve "src/hooks/x.ts" (convenção @/src/… do projeto); deve resolver
	// para o path relativo à raiz sem duplicar o src/ ("apps/mobile/src/src/…").
	root := t.TempDir()
	writeDeep(root+"/apps/mobile/src/hooks/x.ts", "x")
	specRel := "apps/mobile/src/features/f/screens/S.spec.md"
	got := resolveDepPath(root, specRel, "src/hooks/x.ts")
	if got != "apps/mobile/src/hooks/x.ts" {
		t.Errorf("decl com src/ deveria resolver p/ raiz, veio %q", got)
	}
	// e a variante SEM src/ continua funcionando
	got2 := resolveDepPath(root, specRel, "hooks/x.ts")
	if got2 != "apps/mobile/src/hooks/x.ts" {
		t.Errorf("decl sem src/ deveria resolver p/ raiz, veio %q", got2)
	}
}

func TestExtractDepsNoneWithoutSection(t *testing.T) {
	spec := "# Algo\n\n## Data Contract\n\n| Campo | Origem |\n|--|--|\n| x | y |\n"
	if d := extractDeps([]byte(spec), t.TempDir(), "a/src/b.spec.md"); d != nil {
		t.Errorf("sem seção Dependências deveria ser nil, veio %+v", d)
	}
}

func TestExtractDepsIgnoresDocumentalTable(t *testing.T) {
	// tabela "Dependências" documental (Arquivo|Descrição, SEM coluna Cód) — como a de
	// um guide. Não é Tabela de Dependências de reúso; não deve gerar deps.
	root := t.TempDir()
	writeDeep(root+"/src/x.tsx", "x")
	doc := "# Guia\n## Dependências\n| Arquivo | Descrição |\n|--|--|\n| `src/x.tsx` | registrar rotas |\n"
	if d := extractDeps([]byte(doc), root, "guides/G.md"); d != nil {
		t.Errorf("tabela documental (sem coluna Cód) não deveria virar deps: %+v", d)
	}
}

func TestDepsForOnlySpec(t *testing.T) {
	root := t.TempDir()
	writeDeep(root+"/src/x.ts", "x")
	tbl := "# T\n## Dependências\n| Cód | Arquivo | Método |\n|--|--|--|\n| DEP1 | `src/x.ts` | m |\n"
	if d := depsFor("guide", []byte(tbl), root, "guides/G.md"); d != nil {
		t.Errorf("kind guide não deveria extrair deps: %+v", d)
	}
	if d := depsFor("spec", []byte(tbl), root, "a/src/s.spec.md"); len(d) != 1 {
		t.Errorf("kind spec deveria extrair 1 dep, veio %+v", d)
	}
}

func TestExtractDepsSkipsMalformedCode(t *testing.T) {
	root := t.TempDir()
	writeDeep(root+"/src/x.ts", "x")
	spec := "# S\n## Dependências\n| Cód | Arquivo | Método |\n|--|--|--|\n| notdep | `src/x.ts` | m |\n| DEP1 | `src/x.ts` | n |\n"
	deps := extractDeps([]byte(spec), root, "src/s.spec.md")
	if len(deps) != 1 || deps[0].Code != "DEP1" {
		t.Errorf("linha com código malformado deveria ser pulada: %+v", deps)
	}
}

// writeDeep cria um arquivo criando os diretórios do caminho.
func writeDeep(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func TestExtractHeaderDeps(t *testing.T) {
	root := t.TempDir()
	writeDeep(root+"/apps/mobile/src/theme/tokens.ts", "x")
	writeDeep(root+"/apps/mobile/src/utils/cn.ts", "x")
	// arquivo de presentation (sem spec) declarando dep no header
	content := []byte(`// @anchors
//   layer: presentation
//   dep: theme/tokens.ts, utils/cn.ts
//   updated_at: 2026-08-09
export const x = 1
`)
	rel := "apps/mobile/src/features/audit/presentation/auditFeed.ts"
	deps := extractHeaderDeps(content, root, rel)
	if len(deps) != 2 {
		t.Fatalf("esperava 2 header-deps, veio %d: %+v", len(deps), deps)
	}
	if deps[0].File != "apps/mobile/src/theme/tokens.ts" || deps[1].File != "apps/mobile/src/utils/cn.ts" {
		t.Errorf("header-deps mal resolvidas: %+v", deps)
	}
}

func TestDepsForHeaderVsTable(t *testing.T) {
	root := t.TempDir()
	writeDeep(root+"/apps/mobile/src/theme/tokens.ts", "x")
	hdr := []byte("// @anchors\n//   layer: presentation\n//   dep: theme/tokens.ts\n")
	// código (não-spec) → lê header dep
	if d := depsFor("code", hdr, root, "apps/mobile/src/f/presentation/p.ts"); len(d) != 1 {
		t.Errorf("código deveria ler dep do header, veio %+v", d)
	}
	// spec → NÃO lê header dep (usa tabela)
	if d := depsFor("spec", hdr, root, "apps/mobile/src/f/x.spec.md"); d != nil {
		t.Errorf("spec não deveria ler dep do header (só tabela), veio %+v", d)
	}
}

// TestVocabularioDoProjetoEhEnxergadoPeloScan guarda o modo de falha SILENCIOSO: quando a
// classe de letras estava soldada no regex, um projeto que declarasse a sua (`rule_types`)
// tinha os gates enxergando o código e o scan não. O resultado não era erro — era a
// unidade ausente do mapa, sem ninguém acusar.
func TestVocabularioDoProjetoEhEnxergadoPeloScan(t *testing.T) {
	defer SetRuleLetters(config.DefaultRuleLetters)

	const texto = "cenário KVALX-P01: a chave é imutável"
	if got := extractCodes([]byte(texto)); len(got) != 0 {
		t.Fatalf("com o vocabulário canônico, `I` não é letra válida — veio %v", got)
	}

	SetRuleLetters("SRVAXBNMDEIP") // o projeto declarou `P` de Política
	got := extractCodes([]byte(texto))
	if len(got) != 1 || got[0] != "KVALX-P01" {
		t.Errorf("o scan deve enxergar a letra DECLARADA pelo projeto; veio %v", got)
	}
}

// TestDesempateDeCamadaEhEstavelEDeclaravel: `cfg.Layers` é um map, cuja ordem Go não
// define. Sem desempate total, dois patterns de mesmo comprimento sorteariam a camada a
// cada execução — e classificação que muda entre duas rodadas envenena todo gate.
func TestDesempateDeCamadaEhEstavelEDeclaravel(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"amplo":      {Pattern: "pkg/**/*.ts", Kind: "code"},
		"especifico": {Pattern: "pkg/models/**/*.ts", Kind: "spec"},
	}}
	const alvo = "pkg/models/a.ts"

	for i := 0; i < 50; i++ {
		if l, _ := classify(alvo, cfg); l != "especifico" {
			t.Fatalf("heurística instável: rodada %d deu %q", i, l)
		}
	}

	// Onde a heurística erra, o projeto DECLARA — e a declaração vence.
	cfg.Layers["amplo"] = config.Layer{Pattern: "pkg/**/*.ts", Kind: "code", Priority: 10}
	for i := 0; i < 50; i++ {
		if l, _ := classify(alvo, cfg); l != "amplo" {
			t.Fatalf("`priority` declarada deve vencer a heurística; rodada %d deu %q", i, l)
		}
	}
}

// TestSeedIgnoraGlobEmProsa: um plano escreve `*.spec.md` para falar de um CONJUNTO
// ("todo `*.spec.md` precisa de header"), não para citar um arquivo que vai nascer.
// Tratá-lo como seed fazia o plano parecer eternamente não-cumprido — o "arquivo" nunca
// existiria — e a partida a frio da fila passava a semear fases já concluídas.
func TestSeedIgnoraGlobEmProsa(t *testing.T) {
	plano := "Todo `*.spec.md` precisa de header.\n" +
		"A spec de `apps/x/Tela.spec.md` nasce nesta fase.\n" +
		"Não copie o `_TEMPLATE_SCREEN.spec.md`.\n"

	got := extractSeeds("plan", plano)
	if len(got) != 1 || got[0] != "apps/x/Tela.spec.md" {
		t.Errorf("só o caminho concreto é seed; veio %v", got)
	}
}

// TestClassifyNormalizaSeparadorDoWindows: o `classify` recebe caminho de mais de dez
// chamadores, e no Windows o `filepath.Rel` devolve `\`. Sem normalizar, o
// `doublestar.Match` não casa NENHUM pattern que tenha diretório, e o efeito não é uma
// camada errada — é camada NENHUMA: o arquivo some do mapa e nenhum gate volta a
// confrontá-lo. MEDIDO em 24/08 no app de referência: o `map build` no Windows entregava 1857 nós e
// 5416 arestas onde o mesmo commit dava 2799/11302, sem uma linha de erro.
func TestClassifyNormalizaSeparadorDoWindows(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"code": {Pattern: "apps/mobile/src/**/*.ts", Kind: "code"},
	}}
	for _, rel := range []string{
		"apps/mobile/src/business-logic/a.ts",
		`apps\mobile\src\business-logic\a.ts`,
	} {
		if l, k := classify(rel, cfg); l != "code" || k != "code" {
			t.Fatalf("classify(%q) = (%q, %q); as duas formas do mesmo caminho têm de casar", rel, l, k)
		}
	}
}

// TestRevIgnoraFimDeLinha: o `rev` é a identidade do CONTEÚDO — é ele que decide se um
// sinal ingerido ainda vale. Com `core.autocrlf=true`, o mesmo commit tem bytes diferentes
// no Windows e no macOS; se o fim de linha entrasse no hash, o rev de todo arquivo
// divergiria entre as duas máquinas e um `map build` de um lado descartaria os sinais
// acumulados do outro. MEDIDO em 24/08 no app de referência: um `map build` no Windows apagava os 1061
// sinais do mapa gerado no macOS, sem avisar.
func TestRevIgnoraFimDeLinha(t *testing.T) {
	lf := []byte("linha um\nlinha dois\n")
	crlf := []byte("linha um\r\nlinha dois\r\n")
	if shortHash(lf) != shortHash(crlf) {
		t.Fatalf("mesmo conteúdo com fim de linha diferente deu revs diferentes: %s vs %s",
			shortHash(lf), shortHash(crlf))
	}
	// E conteúdo REALMENTE diferente continua com rev diferente — a normalização não pode
	// virar colisão.
	if shortHash(lf) == shortHash([]byte("linha um\nlinha tres\n")) {
		t.Fatal("conteúdos distintos não podem compartilhar rev")
	}
}

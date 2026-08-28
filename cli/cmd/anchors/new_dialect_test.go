package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O esqueleto de teste NASCE na linguagem do projeto. Antes desta correção, `anchors new
// test` emitia `describe/it` com comentário `//` para todas as stacks — um `.py` com
// JavaScript dentro, sintaxe inválida, arquivo que não roda. E contradizia o `anchors
// init`, que oferece 18 stacks.
func TestTestBodyPorFamilia(t *testing.T) {
	casos := []struct {
		family string
		// deve conter: a construção idiomática da linguagem
		contém []string
		// não pode conter: sintaxe de OUTRA linguagem
		nãoContém []string
	}{
		{"python", []string{"def test_calc_total(", "[XXXXX-B01]"}, []string{"describe(", "=>", "//"}},
		{"go", []string{"func TestCalcTotal(t *testing.T)", "[XXXXX-B01]"}, []string{"describe(", "=>", "def "}},
		{"java", []string{"@Test", "void calcTotal()", "[XXXXX-B01]"}, []string{"describe(", "def "}},
		{"kotlin", []string{"@Test", "[XXXXX-B01]"}, []string{"describe(", "def "}},
		{"csharp", []string{"[Fact]", "public void CalcTotal()"}, []string{"describe(", "def "}},
		{"rust", []string{"#[test]", "fn calc_total()"}, []string{"describe(", "def "}},
		{"ruby", []string{"describe '", "do", "end"}, []string{"=>", "def test_"}},
		{"php", []string{"public function testCalcTotal(): void"}, []string{"describe(", "def "}},
		{"ts", []string{"describe(", "it(", "=>"}, []string{"def ", "@Test"}},
	}
	for _, c := range casos {
		t.Run(c.family, func(t *testing.T) {
			got := testBody(c.family, "calcTotal", "XXXXX")
			for _, want := range c.contém {
				if !strings.Contains(got, want) {
					t.Errorf("não contém %q:\n%s", want, got)
				}
			}
			for _, proibido := range c.nãoContém {
				if strings.Contains(got, proibido) {
					t.Errorf("vazou sintaxe de outra linguagem (%q):\n%s", proibido, got)
				}
			}
		})
	}
}

// O scenario-code no NOME do caso é o que o `anchors ingest --junit` cruza com a spec —
// sem ele, o gate scenario-coverage nunca sabe que o requisito tem teste verde. É a
// parte UNIVERSAL do esqueleto e não pode faltar em família alguma.
func TestTestBodySempreCarregaOScenarioCode(t *testing.T) {
	famílias := []string{"python", "go", "java", "kotlin", "csharp", "rust", "ruby", "php", "ts", "", "elixir-que-nao-conhecemos"}
	for _, f := range famílias {
		t.Run("família="+f, func(t *testing.T) {
			if got := testBody(f, "algo", "ABCDX"); !strings.Contains(got, "ABCDX-B01") {
				t.Fatalf("esqueleto sem o scenario-code — o teste nasce órfão do requisito:\n%s", got)
			}
		})
	}
}

// Família declarada que não conhecemos: melhor uma INSTRUÇÃO do que um chute de sintaxe.
// Emitir `describe/it` como "default" é justamente o viés que se está corrigindo.
func TestTestBodyFamiliaDesconhecidaNaoChuta(t *testing.T) {
	got := testBody("cobol", "algo", "ABCDX")
	if strings.Contains(got, "describe(") || strings.Contains(got, "def ") {
		t.Fatalf("chutou sintaxe para família desconhecida:\n%s", got)
	}
	if !strings.Contains(got, "TODO") {
		t.Fatalf("não instruiu o autor:\n%s", got)
	}
}

// A tag de REGIME vem do de-para do projeto. Cravar `@nivel-unit` gerava cenário que
// nenhum gate confronta em projeto com outro vocabulário — e o `anchors work` já
// ensinava que "a tag é do PROJETO e não é traduzível": o template contradizia a régua.
func TestUnitRegimeTagVemDoProjeto(t *testing.T) {
	casos := []struct {
		nome     string
		cfg      *config.Config
		esperado string
	}{
		{"vocabulário em português", &config.Config{Derived: &config.Derived{
			Regimes: map[string]string{"nivel-unit": "unit", "nivel-e2e": "e2e"}}}, "@nivel-unit"},
		{"vocabulário em inglês", &config.Config{Derived: &config.Derived{
			Regimes: map[string]string{"level-unit": "unit"}}}, "@level-unit"},
		{"vocabulário próprio", &config.Config{Derived: &config.Derived{
			Regimes: map[string]string{"fast": "unit", "slow": "e2e"}}}, "@fast"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := unitRegimeTag(c.cfg); got != c.esperado {
				t.Fatalf("tag = %q, queria %q", got, c.esperado)
			}
		})
	}
}

// Sem de-para declarado: um TODO VISÍVEL, não um chute. Uma tag errada passa
// despercebida (nenhum gate a confronta); um TODO é corrigido na hora.
func TestUnitRegimeTagSemDeParaNaoChuta(t *testing.T) {
	for nome, cfg := range map[string]*config.Config{
		"config nula":     nil,
		"sem derived":     {},
		"derived vazio":   {Derived: &config.Derived{}},
		"sem regime unit": {Derived: &config.Derived{Regimes: map[string]string{"lento": "e2e"}}},
	} {
		t.Run(nome, func(t *testing.T) {
			got := unitRegimeTag(cfg)
			if !strings.Contains(got, "TODO") {
				t.Fatalf("chutou uma tag em vez de pedir a declaração: %q", got)
			}
			if strings.Contains(got, "nivel-unit") {
				t.Fatalf("chutou o vocabulário de um projeto específico: %q", got)
			}
		})
	}
}

// A ORDEM das seções é o fio de leitura da spec, e é do preset: numa spec de store,
// "Shape do Estado" precisa vir antes de "Invariantes" — não se enuncia invariante sobre
// um estado que o leitor ainda não conhece. Antes, a ordem vinha sempre do catálogo e
// saía invertida.
func TestOrdemDasSecoesVemDoPreset(t *testing.T) {
	chosen, ordem, err := resolveSectionsWithPreset(specTemplate, "store", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	var chaves []string
	for _, s := range ordenaSecoes(specTemplate, chosen, ordem) {
		chaves = append(chaves, s.Key)
	}
	pos := func(k string) int {
		for i, c := range chaves {
			if c == k {
				return i
			}
		}
		t.Fatalf("seção %q não foi emitida: %v", k, chaves)
		return -1
	}
	if pos("state-shape") > pos("invariants") {
		t.Errorf("Shape do Estado depois de Invariantes — o leitor encontra a invariante "+
			"antes de saber o que é o estado: %v", chaves)
	}
	if pos("actions") < pos("state-shape") {
		t.Errorf("Actions antes do Shape: %v", chaves)
	}
	if chaves[len(chaves)-1] != "open" {
		t.Errorf("Decisões em aberto deveria fechar a spec: %v", chaves)
	}
}

// Seção adicionada por --with que o preset não previu entra ao final, na ordem do
// catálogo — é o único critério disponível para ela.
func TestOrdemComSecaoExtra(t *testing.T) {
	chosen, ordem, err := resolveSectionsWithPreset(specTemplate, "validation", []string{"auth"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := ordenaSecoes(specTemplate, chosen, ordem)
	var achouAuth bool
	for _, s := range got {
		if s.Key == "auth" {
			achouAuth = true
		}
	}
	if !achouAuth {
		t.Fatal("a seção pedida com --with não foi emitida")
	}
}

// A identidade de uma feature/teste é LIDA da spec irmã — eles referenciam (`ref:`), não
// possuem. Gerar um código novo produz um órfão POR CONSTRUÇÃO: o `ref:` aponta para uma
// spec que não existe. Aconteceu: `anchors new feature metadataVersioning` cunhou `MTVA`
// ao lado de uma spec que declarava `MTVRX`.
func TestRefVemDaSpecIrma(t *testing.T) {
	dir := t.TempDir()
	spec := filepath.Join(dir, "metadataVersioning.spec.md")
	if err := os.WriteFile(spec, []byte("<!-- @anchors\n  code: MTVRX\n-->\n# x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	casos := map[string]string{
		"feature ao lado":      "metadataVersioning.feature",
		"teste ao lado":        "metadataVersioning.test.ts",
		"teste go ao lado":     "metadataVersioning_test.go",
		"teste python ao lado": "metadataVersioning_test.py",
	}
	for nome, arquivo := range casos {
		t.Run(nome, func(t *testing.T) {
			got, origem := refDaSpecIrma(dir, filepath.Join(dir, arquivo), "metadataVersioning")
			if got != "MTVRX" {
				t.Fatalf("ref = %q, queria MTVRX (a spec irmã está ao lado); origem=%q", got, origem)
			}
		})
	}
}

// Sem spec irmã, NÃO se inventa uma identidade em silêncio: quem chama gera com aviso.
// O caso mais provável não é "a spec vem depois" — é o --out estar errado.
func TestRefSemSpecIrmaNaoInventa(t *testing.T) {
	dir := t.TempDir()
	if got, _ := refDaSpecIrma(dir, filepath.Join(dir, "semSpec.feature"), "semSpec"); got != "" {
		t.Fatalf("sem spec irmã deveria devolver vazio, veio %q", got)
	}
	// spec de OUTRA unidade no mesmo diretório não serve
	os.WriteFile(filepath.Join(dir, "outraCoisa.spec.md"), []byte("<!-- @anchors\n  code: XXXXX\n-->\n"), 0o644)
	if got, _ := refDaSpecIrma(dir, filepath.Join(dir, "semSpec.feature"), "semSpec"); got != "" {
		t.Fatalf("a spec de outra unidade não é irmã, veio %q", got)
	}
}

package testsig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Fixture no formato do gremlins — a forma exata que `gremlins unleash --output` grava
// (ver internal/report/internal/structure.go do gremlins). Cobre cada status do enum
// (internal/mutator/mutator.go): KILLED, LIVED, NOT COVERED, NOT VIABLE, TIMED OUT.
const fixtureGremlins = `{
  "go_module": "servico-exemplo",
  "test_efficacy": 66.6,
  "mutations_coverage": 75.0,
  "mutants_total": 6,
  "mutants_killed": 2,
  "mutants_lived": 1,
  "mutants_not_viable": 1,
  "mutants_not_covered": 1,
  "elapsed_time": 12.5,
  "files": [
    {
      "file_name": "src/domain/bankaccount.go",
      "mutations": [
        {"type":"CONDITIONALS_NEGATION","status":"KILLED","line":10,"column":4},
        {"type":"ARITHMETIC_BASE","status":"TIMED OUT","line":11,"column":8},
        {"type":"INVERT_NEGATIVES","status":"LIVED","line":42,"column":2},
        {"type":"CONDITIONALS_BOUNDARY","status":"NOT COVERED","line":88,"column":6},
        {"type":"INCREMENT_DECREMENT","status":"NOT VIABLE","line":99,"column":1},
        {"type":"INVERT_LOGICAL","status":"RUNNABLE","line":120,"column":3}
      ]
    }
  ]
}`

func TestParseMutationFormat_Gremlins(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "gremlins.json")
	if err := os.WriteFile(p, []byte(fixtureGremlins), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutationFormat(p, "gremlins")
	if err != nil {
		t.Fatal(err)
	}
	fm, ok := rep.Files["src/domain/bankaccount.go"]
	if !ok {
		t.Fatalf("arquivo não casou; veio %v", rep.Files)
	}
	// KILLED + TIMED OUT = 2 mortos (timeout é morte por travamento).
	if fm.Killed != 2 {
		t.Errorf("Killed = %d, esperado 2", fm.Killed)
	}
	// LIVED + NOT COVERED = 2 sobreviventes (não coberto é, por definição, não provado).
	if fm.Survived != 2 {
		t.Errorf("Survived = %d, esperado 2", fm.Survived)
	}
	// NOT VIABLE e RUNNABLE ficam fora do denominador: 2/(2+2) = 50%.
	if fm.Score != 50 {
		t.Errorf("Score = %v, esperado 50", fm.Score)
	}
	// As linhas dos sobreviventes são o que o autor precisa ver para agir.
	if len(fm.SurvivedAt) != 2 || fm.SurvivedAt[0] != 42 || fm.SurvivedAt[1] != 88 {
		t.Errorf("SurvivedAt = %v, esperado [42 88]", fm.SurvivedAt)
	}
	// O gremlins não escreve limiar no relatório — Low/High ficam ausentes (zero) e o
	// engine cai no default, em vez de herdar régua inventada aqui.
	if rep.Low != 0 || rep.High != 0 {
		t.Errorf("limiares = %v/%v, esperado 0/0 (gremlins não os emite)", rep.Low, rep.High)
	}
}

// O default (formato vazio) continua sendo o canônico: nenhum projeto existente muda de
// comportamento por causa desta feature.
func TestParseMutationFormat_VazioEhMTE(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "mte.json")
	if err := os.WriteFile(p, []byte(fixtureMT), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutationFormat(p, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := rep.Files["src/business-logic/pricing.ts"]; !ok {
		t.Fatalf("o formato vazio deveria ler MTE; veio %v", rep.Files)
	}
}

// Trocar o formato pelo outro tem de FALHAR com mensagem que ensina — é o erro provável
// de quem declara `format:` errado no anchors.yaml, e um silêncio aqui viraria "0
// arquivos" sem explicação.
func TestParseMutationFormat_FormatoTrocado(t *testing.T) {
	dir := t.TempDir()
	mte := filepath.Join(dir, "mte.json")
	if err := os.WriteFile(mte, []byte(fixtureMT), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMutationFormat(mte, "gremlins"); err == nil {
		t.Error("ler MTE como gremlins deveria falhar")
	} else if !strings.Contains(err.Error(), "format") {
		t.Errorf("a mensagem deveria apontar o `format:`; veio: %v", err)
	}

	grem := filepath.Join(dir, "grem.json")
	if err := os.WriteFile(grem, []byte(fixtureGremlins), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMutationFormat(grem, ""); err == nil {
		t.Error("ler gremlins como MTE deveria falhar")
	}
}

func TestParseMutationFormat_Desconhecido(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x.json")
	if err := os.WriteFile(p, []byte(fixtureGremlins), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ParseMutationFormat(p, "stryker4s-xml")
	if err == nil {
		t.Fatal("formato desconhecido deveria falhar")
	}
	// A mensagem precisa LISTAR os aceitos — senão o operador adivinha.
	if !strings.Contains(err.Error(), "gremlins") {
		t.Errorf("a mensagem deveria listar os formatos aceitos; veio: %v", err)
	}
}

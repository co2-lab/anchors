package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// O campo aceita as DUAS formas: string para a spec que governa um arquivo (a maioria) e
// lista para a de configuração, que governa vários. Exigir lista de todas seria ruído em
// troca de nada.
func TestPadroesAceitaStringOuLista(t *testing.T) {
	var d Derived
	if err := yaml.Unmarshal([]byte(`
anchor: spec
files:
  code: "{{dir}}/{{name}}.ts"
  test:
    - "{{dir}}/{{name}}.test.ts"
    - "__tests__/{{name}}.test.ts"
`), &d); err != nil {
		t.Fatal(err)
	}
	if got := d.PadroesDe()["code"]; len(got) != 1 || got[0] != "{{dir}}/{{name}}.ts" {
		t.Errorf("string deveria virar lista de um, veio %v", got)
	}
	if got := d.PadroesDe()["test"]; len(got) != 2 {
		t.Errorf("lista deveria manter os dois, veio %v", got)
	}
}

// `patterns` é uma chave DENTRO de `files`, irmã de `code`/`feature`/`test`, e substitui
// o `code` quando declarada — mas o resto de `files` sobrevive.
//
// Uma spec de configuração pode ter `patterns` para o código e continuar querendo o
// `feature`/`test` da co-location. Descartar o mapa inteiro obrigaria a repetir o que não
// mudou, e repetição em config é onde a divergência começa.
func TestPatternsSubstituiCodeEMantemOResto(t *testing.T) {
	var d Derived
	if err := yaml.Unmarshal([]byte(`
anchor: spec
files:
  code: "{{dir}}/{{name}}.ts"
  feature: "{{dir}}/{{name}}.feature"
  test: "{{dir}}/{{name}}.test.ts"
  patterns:
    - "tsconfig.base.json"
    - "packages/*/tsconfig.json"
`), &d); err != nil {
		t.Fatal(err)
	}
	got := d.PadroesDe()
	// O código vem de `patterns`.
	if len(got["code"]) != 2 || got["code"][0] != "tsconfig.base.json" {
		t.Errorf("`patterns` deveria substituir `code`, veio %v", got["code"])
	}
	// E o resto continua vindo de `files`.
	if len(got["feature"]) != 1 || got["feature"][0] != "{{dir}}/{{name}}.feature" {
		t.Errorf("`feature` deveria sobreviver, veio %v", got["feature"])
	}
	// `patterns` NÃO é uma camada: não pode virar um derivado com esse nome, ou o mapa
	// passaria a procurar um arquivo "patterns" que ninguém escreveu.
	if _, virouCamada := got[PatternKey]; virouCamada {
		t.Error("`patterns` não é camada e não pode sobrar no mapa de derivados")
	}
}

// Sem `patterns`, o comportamento é o de sempre: todo anchors.yaml que já existe continua
// funcionando sem tocar em nada.
func TestSemPatternsUsaFiles(t *testing.T) {
	var d Derived
	if err := yaml.Unmarshal([]byte("anchor: spec\nfiles:\n  code: \"x.ts\"\n"), &d); err != nil {
		t.Fatal(err)
	}
	if got := d.PadroesDe()["code"]; len(got) != 1 || got[0] != "x.ts" {
		t.Errorf("sem patterns, `files` manda: %v", got)
	}
}

// The loader refuses what is neither text nor a non-empty list of text, and says why.
func TestPadroesRejectsOtherShapes(t *testing.T) {
	for name, doc := range map[string]string{
		"empty list":       "code: []\n",
		"mapping":          "code: {a: b}\n",
		"list of mappings": "code:\n  - {a: b}\n",
	} {
		var m map[string]Padroes
		if err := yaml.Unmarshal([]byte(doc), &m); err == nil {
			t.Errorf("%s: want an error, got %v", name, m)
		}
	}
	var m map[string]Padroes
	err := yaml.Unmarshal([]byte("code: []\n"), &m)
	if err == nil || !strings.Contains(err.Error(), "empty pattern list") {
		t.Errorf("an empty list must name its cause, got %v", err)
	}
}

// Writing back keeps the simplest shape: one pattern is a string, several are a list.
func TestPadroesMarshalsToTheSimplestShape(t *testing.T) {
	out, err := yaml.Marshal(map[string]Padroes{
		"code": {"x.ts"},
		"test": {"a.test.ts", "b.test.ts"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := "code: x.ts\ntest:\n    - a.test.ts\n    - b.test.ts\n"; string(out) != want {
		t.Errorf("Marshal = %q, want %q", out, want)
	}
}

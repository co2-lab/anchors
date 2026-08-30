package config

import (
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

// `patterns` VENCE `files` quando os dois declaram a mesma camada — mas MESCLA o resto.
//
// Um override pode querer padrões só para o código e continuar usando o `feature`/`test`
// da co-location. Substituir o mapa inteiro obrigaria a repetir o que não mudou, e
// repetição em config é onde a divergência começa.
func TestPatternsVenceFilesEMesclaOResto(t *testing.T) {
	var d Derived
	if err := yaml.Unmarshal([]byte(`
anchor: spec
files:
  code: "{{dir}}/{{name}}.ts"
  feature: "{{dir}}/{{name}}.feature"
  test: "{{dir}}/{{name}}.test.ts"
patterns:
  code:
    - "tsconfig.base.json"
    - "packages/*/tsconfig.json"
`), &d); err != nil {
		t.Fatal(err)
	}
	got := d.PadroesDe()
	// O código vem de `patterns`.
	if len(got["code"]) != 2 || got["code"][0] != "tsconfig.base.json" {
		t.Errorf("`patterns` deveria vencer para `code`, veio %v", got["code"])
	}
	// E o que `patterns` não declarou continua vindo de `files`.
	if len(got["feature"]) != 1 || got["feature"][0] != "{{dir}}/{{name}}.feature" {
		t.Errorf("`feature` deveria vir de `files`, veio %v", got["feature"])
	}
	if len(got["test"]) != 1 {
		t.Errorf("`test` deveria vir de `files`, veio %v", got["test"])
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

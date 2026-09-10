package main

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O NOME DA UNIDADE não é o basename menos a última extensão.
//
// Medido no blue-eyes: `anchors work code --for packages/infra/GoLiveChecklist.spec.md`
// mandava escrever `GoLiveChecklist.spec.feature` e `GoLiveChecklist.spec.test.ts`.
//
// A causa era `strings.TrimSuffix(base, filepath.Ext(rel))`: `filepath.Ext("X.spec.md")`
// é `.md`, então o `{{name}}` saía `X.spec`. O mapa nunca teve esse problema — ele usa
// `mapx.StemOfAnchor`, que corta o sufixo COMPOSTO — e as duas implementações discordavam.
//
// O dano é que o prompt e o mapa apontam para arquivos diferentes: quem seguisse o prompt
// criaria `X.spec.ts`, e o `triad-complete` continuaria procurando `X.ts`.
func TestDerivedPaths_naoDeixaSpecNoNome(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "spec",
		Files: map[string]config.Padroes{
			"code":    {"{{dir}}/{{name}}.ts"},
			"feature": {"{{dir}}/{{name}}.feature"},
			"test":    {"{{dir}}/{{name}}.test.ts"},
		},
	}}

	files, _ := derivedPaths("packages/infra/GoLiveChecklist.spec.md", "spec", cfg)

	quer := map[string]string{
		"code":    "packages/infra/GoLiveChecklist.ts",
		"feature": "packages/infra/GoLiveChecklist.feature",
		"test":    "packages/infra/GoLiveChecklist.test.ts",
	}
	for k, esperado := range quer {
		if got := files[k]; got != esperado {
			t.Errorf("%s: %q — esperava %q\n  (o `.spec` no meio faz o prompt apontar "+
				"para arquivo que nenhum gate encontra)", k, got, esperado)
		}
	}
}

// E a FEATURE também é âncora de sufixo composto: `X.feature` → `X`.
func TestDerivedPaths_cortaOSufixoDaFeature(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "spec",
		Files:  map[string]config.Padroes{"test": {"{{dir}}/{{name}}.test.ts"}},
	}}

	files, _ := derivedPaths("apps/mobile/src/Login.feature", "feature", cfg)

	if got := files["test"]; got != "apps/mobile/src/Login.test.ts" {
		t.Errorf("test: %q — esperava apps/mobile/src/Login.test.ts", got)
	}
}

package main

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func cfgComDeclarativa() *config.Config {
	return &config.Config{
		Layers: map[string]config.Layer{
			"dao":   {Pattern: "backend/models/**/*.ts", Kind: "code", Regime: "declarativo"},
			"logic": {Pattern: "backend/business-logic/**/*.ts", Kind: "code"},
			"spec":  {Pattern: "**/*.spec.md", Kind: "spec"},
		},
		// `derived:` é o que resolve os caminhos da trinca; sem ele o prompt não tem o que
		// derivar (e diz isso). A fixture precisa dele para exercitar o caminho comum.
		Derived: &config.Derived{
			Anchor: "code",
			Files:  map[string]config.Padroes{"spec": {"{{dir}}/{{name}}.spec.md"}},
		},
	}
}

// O prompt se contradizia: declarava "Camada: dao (regime: declarativo)" e três linhas
// abaixo listava `metadata.spec.md` entre "as peças e onde nascem" — a spec que a camada
// PROÍBE. O procedimento ainda mandava "leia a spec inteira; ela é a régua". Um agente que
// confiasse nisso criaria a peça proibida, e o `anchors new` depois a recusaria sem que
// nada explicasse a contradição.
func TestWorkCamadaDeclarativaNaoPedeSpec(t *testing.T) {
	out, err := composeWorkPrompt(t.TempDir(), "backend/models/metadata.ts", "code", cfgComDeclarativa(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out, "metadata.spec.md") {
		t.Errorf("listou a spec que a camada declarativa proíbe:\n%s", out)
	}
	if strings.Contains(out, "Leia a spec inteira") {
		t.Errorf("mandou ler uma spec que não existe:\n%s", out)
	}
	// e diz o que vale no lugar
	for _, esperado := range []string{"não tem spec", "não origina regra", "não decida"} {
		if !strings.Contains(out, esperado) {
			t.Errorf("o prompt não explica a camada declarativa (falta %q)", esperado)
		}
	}
}

// A camada REGIDA continua recebendo a trinca e o procedimento normal — a correção não
// pode ter esvaziado o caminho comum.
func TestWorkCamadaRegidaMantemATrinca(t *testing.T) {
	out, err := composeWorkPrompt(t.TempDir(), "backend/business-logic/recurrence.ts", "code", cfgComDeclarativa(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "recurrence.spec.md") {
		t.Errorf("camada regida deveria listar a spec da trinca:\n%s", out)
	}
	if !strings.Contains(out, "Leia a spec inteira") {
		t.Errorf("camada regida deveria manter o procedimento normal")
	}
}

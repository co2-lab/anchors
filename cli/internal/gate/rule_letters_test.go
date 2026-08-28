package gate

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O vocabulário de letras é do PROJETO (`rule_types`). Um regex de código preso às
// canônicas torna INVISÍVEL todo cenário de uma letra declarada — e o gate reporta verde
// sobre o que não conferiu, o modo de falha mais perigoso. Aconteceu de verdade com a
// letra `I` (Invariant): o cenário existia, o gate não o via, e ninguém percebeu.
func TestLetraDeclaradaPeloProjetoEhVista(t *testing.T) {
	defer SetRuleLetters(config.DefaultRuleLetters) // não vazar para outros testes

	feature := "@ABCDX-I01 @nivel-unit\n  Cenário: a invariante vale sempre\n"

	// antes de declarar: a letra não é do vocabulário, então não é vista — correto.
	SetRuleLetters(config.DefaultRuleLetters)
	if got := parseFeatureScenarios(feature); len(got) != 0 {
		t.Fatalf("letra NÃO declarada não deveria ser reconhecida, veio %+v", got)
	}

	// depois de declarar: passa a ser vista.
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "B", Term: "Behavior"}, {Letter: "I", Term: "Invariant"},
	}}
	SetRuleLetters(cfg.RuleLetters())
	got := parseFeatureScenarios(feature)
	if len(got) != 1 || got[0].Code != "ABCDX-I01" {
		t.Fatalf("cenário da letra declarada continua invisível: %+v", got)
	}
}

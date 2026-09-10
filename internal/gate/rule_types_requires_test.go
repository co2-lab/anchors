package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func cfgComRequires() *config.Config {
	return &config.Config{RuleTypes: []config.RuleType{
		{Letter: "B", Term: "Behavior",
			Sections:     []string{"Eventos / Callbacks", "Comportamentos"},
			RequiresCode: []string{"Eventos / Callbacks"}},
		{Letter: "S", Term: "State",
			Sections: []string{"Variantes", "Estados Visuais"}},
	}}
}

// A seção declarada `requires_code` com tabela preenchida e nenhum código é o furo
// por onde o cenário fica sem âncora — e acaba emprestando o código de outra seção.
func TestSecaoQueExigeCodigoSemCodigoEhAchado(t *testing.T) {
	v, msg := checkRuleTypes(`# Spec

## Eventos / Callbacks

| Evento    | Quando       | Payload |
| --------- | ------------ | ------- |
| `+"`onPress`"+` | Tap na linha | —       |
`, specNode(), "", nil, cfgComRequires())

	if v != Pending {
		t.Fatalf("veredito %v (%s), queria Pending", v, msg)
	}
	if !strings.Contains(msg, "Eventos / Callbacks") {
		t.Errorf("mensagem não nomeia a seção: %s", msg)
	}
}

// Com o código na tabela, não há o que cobrar.
func TestSecaoComCodigoPassa(t *testing.T) {
	v, msg := checkRuleTypes(`# Spec

## Eventos / Callbacks

| Regra      | Evento    | Quando       |
| ---------- | --------- | ------------ |
| `+"`BGCRX-B01`"+` | `+"`onPress`"+` | Tap na linha |
`, specNode(), "", nil, cfgComRequires())

	if v == Pending {
		t.Errorf("cobrou seção que JÁ tem código: %s", msg)
	}
}

// "Variantes" está declarada sob a letra S, mas NÃO em `requires_code`: ela apenas
// enumera valores. Cobrar código dela seria exigir regra de um índice.
func TestSecaoDeclaradaSemRequiresCodeNaoEhCobrada(t *testing.T) {
	v, msg := checkRuleTypes(`# Spec

## Variantes

| `+"`seatType`"+` | Emoji | Label |
| ---------- | ----- | ----- |
| `+"`individual`"+` | 👤    | Individual |
| `+"`familia`"+`    | 👨‍👩‍👧  | Família    |
`, specNode(), "", nil, cfgComRequires())

	if v == Pending {
		t.Errorf("cobrou seção que só enumera valores: %s", msg)
	}
}

// Projeto que não usa `requires_code` não muda de comportamento — a régua nasce
// opt-in, senão o gate acusaria toda base existente de uma vez.
func TestSemRequiresCodeNadaMuda(t *testing.T) {
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "B", Term: "Behavior", Sections: []string{"Eventos / Callbacks"}},
	}}
	v, msg := checkRuleTypes(`# Spec

## Eventos / Callbacks

| Evento    | Quando       |
| --------- | ------------ |
| `+"`onPress`"+` | Tap na linha |
`, specNode(), "", nil, cfg)

	if v == Pending {
		t.Errorf("cobrou sem o projeto ter declarado `requires_code`: %s", msg)
	}
}

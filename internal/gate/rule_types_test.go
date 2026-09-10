package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func vocab() *config.Config {
	return &config.Config{RuleTypes: []config.RuleType{
		{Letter: "S", Term: "State", Sections: []string{"Fluxo de Estados"}},
		{Letter: "B", Term: "Behavior", Sections: []string{"Comportamentos"}},
		{Letter: "E", Term: "Error", Sections: []string{"Erros / Falhas"}},
	}}
}

func TestRuleTypes_letraNaoDeclarada(t *testing.T) {
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ok |\n## Planos\n| `HOMEX-P01` | limite |\n"
	v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab())
	if v != Fail {
		t.Fatalf("esperava Fail p/ letra não declarada, got %v", v)
	}
	if !strings.Contains(msg, "P") {
		t.Errorf("mensagem deveria citar a letra P: %q", msg)
	}
}

func TestRuleTypes_letraDeclaradaPassa(t *testing.T) {
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ok |\n## Erros / Falhas\n| `HOMEX-E01` | falhou |\n"
	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("esperava Pass, got %v (%s)", v, msg)
	}
}

func TestRuleTypes_secaoSemLetraDeclarada(t *testing.T) {
	// letra declarada (S), mas sob uma seção que ninguém reivindica
	content := "## Regras Inventadas\n| `HOMEX-S01` | ok |\n"
	v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab())
	if v != Fail {
		t.Fatalf("esperava Fail p/ seção não reivindicada, got %v", v)
	}
	if !strings.Contains(msg, "Regras Inventadas") {
		t.Errorf("mensagem deveria citar a seção: %q", msg)
	}
}

func TestRuleTypes_conflitoDeLetra(t *testing.T) {
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "E", Term: "Error", Sections: []string{"Erros"}},
		{Letter: "E", Term: "Estado", Sections: []string{"Estados"}},
	}}
	v, msg := checkRuleTypes("## Erros\n| `HOMEX-E01` | x |\n", mapx.Node{}, "", nil, cfg)
	if v != Fail {
		t.Fatalf("esperava Fail p/ conflito de letra, got %v", v)
	}
	if !strings.Contains(msg, "CONFLITO") {
		t.Errorf("mensagem deveria sinalizar CONFLITO: %q", msg)
	}
}

// Este teste travava o Skip — o gate canônico que nunca media nada. Passou a confrontar
// as letras canônicas: `P` não está em SRVAXBNMD, e uma letra que o engine não reconhece é
// invisível para a rastreabilidade, com ou sem vocabulário declarado.
func TestRuleTypes_semVocabularioUsaAsCanonicas(t *testing.T) {
	v, msg := checkRuleTypes("## X\n| `HOMEX-P01` | x |\n", mapx.Node{}, "", nil, &config.Config{})
	if v != Fail {
		t.Errorf("`P` está fora das canônicas e deveria reprovar, got %v", v)
	}
	if !strings.Contains(msg, "canônico") {
		t.Errorf("a mensagem deveria dizer que confronta o vocabulário canônico: %q", msg)
	}
}

func TestRuleTypes_tituloQueEhOProprioCodigoNaoConta(t *testing.T) {
	// "### OCHIX-S01: ..." é cabeçalho da regra, não seção-categoria
	content := "### HOMEX-S01: Estado inicial\nAlguma prosa com `HOMEX-S01` citado.\n"
	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("cabeçalho de regra não deveria ser cobrado como seção: %v (%s)", v, msg)
	}
}

func TestRuleLetters_derivaDoVocabulario(t *testing.T) {
	if got := vocab().RuleLetters(); got != "SBE" {
		t.Errorf("RuleLetters = %q, quer SBE", got)
	}
	var nilCfg *config.Config
	if got := nilCfg.RuleLetters(); got != config.DefaultRuleLetters {
		t.Errorf("cfg nil deveria cair nas canônicas, got %q", got)
	}
}

func TestRuleTypes_secaoQueSoCITAcodigoNaoConta(t *testing.T) {
	// "Test IDs (Maestro)" referencia códigos de OUTRAS seções numa coluna interna —
	// não define regra, logo não precisa reivindicar letra.
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ok |\n" +
		"## Test IDs (Maestro)\n| testID | Elemento | Usado em |\n| `home-root` | Raiz | `HOMEX-S01`, `HOMEX-VR` |\n"
	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("seção que só cita código não deveria reprovar: %v (%s)", v, msg)
	}
}

func TestRuleTypes_secaoQueDEFINEcodigoEmTabelaConta(t *testing.T) {
	// primeira célula da tabela = definição → seção precisa de letra declarada
	content := "## Regras Inventadas\n| Regra | Descrição |\n| `HOMEX-S02` | define algo |\n"
	if v, _ := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Fail {
		t.Errorf("seção que DEFINE regra sob título não declarado deveria reprovar, got %v", v)
	}
}

// Sem `rule_types` declarado o gate media NADA: fazia Skip e reportava indeterminado para
// sempre — um gate canônico, semeado pelo `init` em todo projeto, que nunca confrontou
// coisa alguma. É a impressão de defesa que não existe.
func TestRuleTypesConfrontaCanonicasSemVocabulario(t *testing.T) {
	semVocabulario := &config.Config{}

	// Letra canônica: passa.
	spec := "## Regras\n\n### ABCDX-B01 — regra\n\nTexto.\n"
	if v, msg := checkRuleTypes(spec, mapx.Node{}, "", nil, semVocabulario); v != Pass {
		t.Errorf("`B` é canônica e deveria passar: %v — %s", v, msg)
	}

	// Letra fora das canônicas: reprova, e nomeia a letra. `P` (política) é um tipo que o
	// framework não conhece — o exemplo antes era `I`, que passou a ser canônica quando se
	// descobriu que o `anchors new` já a emitia.
	fora := "## Regras\n\n### ABCDX-P01 — política\n\nTexto.\n"
	v, msg := checkRuleTypes(fora, mapx.Node{}, "", nil, semVocabulario)
	if v != Fail {
		t.Errorf("`P` não está em %s e deveria reprovar, veio %v", config.DefaultRuleLetters, v)
	}
	if !strings.Contains(msg, "P") || !strings.Contains(msg, "rule_types") {
		t.Errorf("a mensagem deveria nomear a letra E onde declará-la: %s", msg)
	}

	// Sem código nenhum não há o que confrontar — e isso NÃO é falha deste gate: quem
	// cobra a existência de regra catalogada é o `spec-completa`.
	if v, _ := checkRuleTypes("## Visão Geral\n\nProsa.\n", mapx.Node{}, "", nil, semVocabulario); v != Pass {
		t.Errorf("spec sem código não é problema do rule-types: %v", v)
	}
}

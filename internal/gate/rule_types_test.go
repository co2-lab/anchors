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
	t.Run("RLTYR-B01: A letter that is not declared in the vocabulary fails", func(t *testing.T) {})
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
	t.Run("RLTYR-B02: A declared letter under a claimed section passes", func(t *testing.T) {})
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ok |\n## Erros / Falhas\n| `HOMEX-E01` | falhou |\n"
	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("esperava Pass, got %v (%s)", v, msg)
	}
}

func TestRuleTypes_secaoSemLetraDeclarada(t *testing.T) {
	t.Run("RLTYR-B03: A section cataloguing rules under a title no letter claims fails", func(t *testing.T) {})
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
	t.Run("RLTYR-B04: The same letter claimed by two terms is a conflict in the vocabulary", func(t *testing.T) {})
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
	t.Run("RLTYR-B05: With no vocabulary declared the gate confronts the canonical letters", func(t *testing.T) {})
	v, msg := checkRuleTypes("## X\n| `HOMEX-P01` | x |\n", mapx.Node{}, "", nil, &config.Config{})
	if v != Fail {
		t.Errorf("`P` está fora das canônicas e deveria reprovar, got %v", v)
	}
	if !strings.Contains(msg, "canônico") && !strings.Contains(msg, "canonical") {
		t.Errorf("a mensagem deveria dizer que confronta o vocabulário canônico: %q", msg)
	}
}

func TestRuleTypes_tituloQueEhOProprioCodigoNaoConta(t *testing.T) {
	t.Run("RLTYR-B06: A heading that is the rule code itself is not a category section", func(t *testing.T) {})
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
	t.Run("RLTYR-B07: A section that only cites other sections' codes claims no letter", func(t *testing.T) {})
	// "Test IDs (Maestro)" referencia códigos de OUTRAS seções numa coluna interna —
	// não define regra, logo não precisa reivindicar letra.
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ok |\n" +
		"## Test IDs (Maestro)\n| testID | Elemento | Usado em |\n| `home-root` | Raiz | `HOMEX-S01`, `HOMEX-VR` |\n"
	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("seção que só cita código não deveria reprovar: %v (%s)", v, msg)
	}
}

func TestRuleTypes_secaoQueDEFINEcodigoEmTabelaConta(t *testing.T) {
	t.Run("RLTYR-B08: A section that defines a code in the first table cell is charged", func(t *testing.T) {})
	// primeira célula da tabela = definição → seção precisa de letra declarada
	content := "## Regras Inventadas\n| Regra | Descrição |\n| `HOMEX-S02` | define algo |\n"
	if v, _ := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Fail {
		t.Errorf("seção que DEFINE regra sob título não declarado deveria reprovar, got %v", v)
	}
}

// Sem código nenhum não há o que confrontar — e isso NÃO é falha deste gate: quem cobra
// a existência de regra catalogada é o `spec-completa`. Acusar nos dois duplicaria a régua.
func TestRuleTypesSpecSemCodigoNaoEhProblemaDaqui(t *testing.T) {
	t.Run("RLTYR-I01: A spec with no rule code at all is not this gate's problem", func(t *testing.T) {})
	semVocabulario := &config.Config{}

	// Letra canônica presente: passa, e confirma que o gate ESTÁ medindo.
	if v, msg := checkRuleTypes("## Regras\n\n### ABCDX-B01 — regra\n\nTexto.\n", mapx.Node{}, "", nil, semVocabulario); v != Pass {
		t.Errorf("`B` é canônica e deveria passar: %v — %s", v, msg)
	}

	// Prosa pura: nada a confrontar.
	if v, msg := checkRuleTypes("## Visão Geral\n\nProsa.\n", mapx.Node{}, "", nil, semVocabulario); v != Pass {
		t.Errorf("spec sem código não é problema do rule-types: %v (%s)", v, msg)
	}
}

// O veredito que reprova sem dizer O QUÊ e ONDE DECLARAR transfere o trabalho de
// diagnóstico para quem lê.
func TestRuleTypesVereditoNomeiaLetraEOndeDeclarar(t *testing.T) {
	t.Run("RLTYR-I02: The verdict names the letter and where to declare it", func(t *testing.T) {})
	// `P` (política) é um tipo que o framework não conhece — o exemplo antes era `I`, que
	// passou a ser canônica quando se descobriu que o `anchors new` já a emitia.
	fora := "## Regras\n\n### ABCDX-P01 — política\n\nTexto.\n"
	v, msg := checkRuleTypes(fora, mapx.Node{}, "", nil, &config.Config{})
	if v != Fail {
		t.Fatalf("`P` não está em %s e deveria reprovar, veio %v", config.DefaultRuleLetters, v)
	}
	if !strings.Contains(msg, "P") {
		t.Errorf("a mensagem não nomeia a LETRA: %s", msg)
	}
	if !strings.Contains(msg, "rule_types") {
		t.Errorf("a mensagem não diz ONDE declarar: %s", msg)
	}
}

// O vocabulário é EXTENSÍVEL por desenho: uma letra fora das canônicas, mas DECLARADA
// com sua seção, passa. As canônicas são o fallback, não um teto.
func TestRuleTypesVocabularioEExtensivel(t *testing.T) {
	t.Run("RLTYR-X01: The gate does not decide which letters exist", func(t *testing.T) {})
	if strings.Contains(config.DefaultRuleLetters, "P") {
		t.Fatalf("o teste depende de `P` estar FORA das canônicas (%s)", config.DefaultRuleLetters)
	}
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "P", Term: "Política", Sections: []string{"Políticas"}},
	}}
	content := "## Políticas\n| `HOMEX-P01` | limite |\n"

	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, cfg); v != Pass {
		t.Errorf("letra declarada fora das canônicas deveria passar — o vocabulário é do "+
			"projeto, e as canônicas são fallback: %v (%s)", v, msg)
	}
}

// Sem vocabulário declarado só a LETRA é cobrada. Seções e termos só existem depois que o
// projeto declara; cobrá-los contra um vocabulário implícito inventaria regra.
func TestRuleTypesSemVocabularioNaoCobraSecao(t *testing.T) {
	t.Run("RLTYR-X02: Without a declared vocabulary only the letter is charged", func(t *testing.T) {})
	// Letra canônica (`B`) sob um título que vocabulário NENHUM reivindica.
	content := "## Um Título Que Ninguém Reivindica\n| `HOMEX-B01` | faz algo |\n"

	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, &config.Config{}); v != Pass {
		t.Errorf("sem vocabulário não há seção a reivindicar — cobrar inventaria regra "+
			"que ninguém escreveu: %v (%s)", v, msg)
	}
}

// Que letra uma regra MERECE é julgamento editorial. A régua aqui é a rastreabilidade
// enxergar o código, não a escolha ser a melhor.
func TestRuleTypesNaoJulgaSeALetraCabe(t *testing.T) {
	t.Run("RLTYR-X03: The gate does not judge whether the letter suits the rule", func(t *testing.T) {})
	// Uma regra claramente de COMPORTAMENTO catalogada sob a letra de ESTADO.
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ao tocar no botão, dispara o envio |\n"

	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("a escolha da letra é editorial; a régua é a rastreabilidade: %v (%s)", v, msg)
	}
}

// O gate cobra RASTREABILIDADE, não formato: uma seção não reivindicada que cataloga
// coisa nenhuma não tem rastreabilidade a defender, qualquer que seja sua forma.
func TestRuleTypesNaoCobraFormato(t *testing.T) {
	t.Run("RLTYR-X04: The gate charges traceability, not format", func(t *testing.T) {})
	content := "## Fluxo de Estados\n| `HOMEX-S01` | ok |\n" +
		"## Notas Soltas De Formato Estranho\n\nProsa sem tabela, sem bullet, sem código.\n"

	if v, msg := checkRuleTypes(content, mapx.Node{}, "", nil, vocab()); v != Pass {
		t.Errorf("seção sem código não cataloga nada e não deve ser cobrada: %v (%s)", v, msg)
	}
}

// A seção declarada `requires_code` com tabela preenchida e nenhum código é o furo
// por onde o cenário fica sem âncora — e acaba emprestando o código de outra seção.
func TestSecaoQueExigeCodigoSemCodigoEhAchado(t *testing.T) {
	t.Run("RLTYR-B09: A section declared as rule-cataloguing and filled without a code is Pending", func(t *testing.T) {})
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
	t.Run("RLTYR-B10: A section whose table already carries the code is not charged", func(t *testing.T) {})
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
	t.Run("RLTYR-B11: A declared section outside sections_require_code is not charged", func(t *testing.T) {})
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
	t.Run("RLTYR-B12: A project that does not use sections_require_code changes no behaviour", func(t *testing.T) {})
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

// O FURO por onde o cenário fica sem âncora. Medido num projeto real: 48 specs com
// "Eventos / Callbacks" preenchida e sem um único código, e os cenários que provavam
// esses eventos emprestaram o código do estado vizinho — um `-S` regendo comportamento.
//
// A prova é a ASSIMETRIA: a mesma seção, com e sem código na tabela. Sem código o gate
// reporta; com código cala. Se calasse nos dois, o cenário órfão seguiria invisível.
func TestSecaoSemCodigoEhOndeOCenarioPerdeAAncora(t *testing.T) {
	t.Run("RLTYR-I03: A filled section with no code is the gap where the scenario loses its anchor", func(t *testing.T) {})

	semCodigo := "# Spec\n\n## Eventos / Callbacks\n\n| Evento | Quando |\n| --- | --- |\n" +
		"| `onPress` | Tap na linha |\n"
	v, msg := checkRuleTypes(semCodigo, specNode(), "", nil, cfgComRequires())
	if v != Pending {
		t.Fatalf("a seção preenchida sem código deveria ser reportada, veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "Eventos / Callbacks") {
		t.Errorf("o achado não nomeia a seção onde o cenário perde a âncora: %s", msg)
	}

	// A mesma seção, agora com o código: o achado some. É o que prova que o gate mede a
	// AUSÊNCIA do código, e não a existência da seção.
	comCodigo := "# Spec\n\n## Eventos / Callbacks\n\n| Regra | Evento | Quando |\n| --- | --- | --- |\n" +
		"| `BGCRX-B01` | `onPress` | Tap na linha |\n"
	if v, msg := checkRuleTypes(comCodigo, specNode(), "", nil, cfgComRequires()); v == Pending {
		t.Errorf("com o código na tabela não há âncora perdida a reportar: %s", msg)
	}
}

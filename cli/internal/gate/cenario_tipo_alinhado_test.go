package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func cfgComTags() *config.Config {
	return &config.Config{RuleTypes: []config.RuleType{
		{Letter: "S", Term: "State", Tags: []string{"estado"}},
		{Letter: "B", Term: "Behavior", Tags: []string{"comportamento"}},
	}}
}

func TestCenarioTipoAlinhado_acusaTagQueDiscordaDaLetra(t *testing.T) {
	feat := `# language: pt
Funcionalidade: Seção

  @comportamento @LGSTX-S01#02 @nivel-integration @P2
  Cenário: Intro presente é renderizada em itálico
    Quando o componente é renderizado
    Então devo ver o texto de intro
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	v, detail := checkCenarioTipoAlinhado(feat, n, "", nil, cfgComTags())
	if v != Pending {
		t.Fatalf("esperava Pending, veio %v: %s", v, detail)
	}
	if !strings.Contains(detail, "LGSTX-S01") || !strings.Contains(detail, "comportamento") {
		t.Errorf("o detalhe não nomeia o código nem a tag: %s", detail)
	}
}

func TestCenarioTipoAlinhado_passaQuandoConcordam(t *testing.T) {
	feat := `# language: pt
Funcionalidade: Seção

  @estado @LGSTX-S01 @P2
  Cenário: Seção padrão exibe número e título
    Então devo ver o título

  @comportamento @LGSTX-B01 @P2
  Cenário: Tocar no link abre os termos
    Quando eu toco no link
    Então devo ver os termos
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkCenarioTipoAlinhado(feat, n, "", nil, cfgComTags()); v != Pass {
		t.Errorf("esperava Pass, veio %v: %s", v, detail)
	}
}

// Sem `tags:` no vocabulário, o gate não pode adivinhar que `@estado` é S — e
// adivinhar acusaria todo projeto que usa outro idioma nas tags.
func TestCenarioTipoAlinhado_semMapaFicaEmSilencio(t *testing.T) {
	cfg := &config.Config{RuleTypes: []config.RuleType{{Letter: "S", Term: "State"}}}
	feat := "# language: pt\n\n  @comportamento @LGSTX-S01 @P2\n  Cenário: x\n    Então y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkCenarioTipoAlinhado(feat, n, "", nil, cfg); v != Skip {
		t.Errorf("esperava Skip sem mapa tag→letra, veio %v", v)
	}
}

// Códigos que não carregam natureza (`-DS-…`, `-VR`) não têm letra a confrontar.
func TestLetraDoCodigoIgnoraCodigosSemNatureza(t *testing.T) {
	casos := map[string]string{
		"BUGEX-S02":                "S",
		"BUGEX-S02#01":             "S",
		"BUGEX-DS-mode-standalone": "",
		"BUGEX-VR":                 "",
	}
	for entrada, quer := range casos {
		if got := letraDoCodigo(entrada); got != quer {
			t.Errorf("letraDoCodigo(%q) = %q, queria %q", entrada, got, quer)
		}
	}
}

// A mesma tag sob duas letras não é ambiguidade: `@estado-dado` marca cenários de S
// ("a tela exibe X quando o dado é Y") e de V, cuja seção cataloga aparência-por-prop.
// Acusar a segunda obrigaria a escolher entre duas classificações corretas.
func TestCenarioTipoAlinhado_tagPodeCaberEmMaisDeUmaLetra(t *testing.T) {
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "S", Term: "State", Tags: []string{"estado", "estado-dado"}},
		{Letter: "V", Term: "Validation", Tags: []string{"validacao", "estado-dado"}},
	}}
	feat := `# language: pt
Funcionalidade: Badge

  @estado-dado @PRBDX-V01 @P2
  Cenário: Variante high define fundo, cor e rótulo
    Então devo ver o texto "Alta"

  @estado-dado @PRBDX-S01 @P2
  Cenário: Sem prioridade não exibe a pílula
    Então não devo ver a pílula
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkCenarioTipoAlinhado(feat, n, "", nil, cfg); v != Pass {
		t.Errorf("esperava Pass — a tag cabe nas duas letras: %v %s", v, detail)
	}
}

// Um cenário pode provar mais de um requisito, e a tag pode descrever o segundo:
// `@mensagem @IPSBX-V02 @IPSBX-M01` classifica pelo M, não pelo V. Acusar isso trocaria
// uma classificação certa por outra.
func TestCenarioTipoAlinhado_tagPodeDescreverCodigoSecundario(t *testing.T) {
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "V", Term: "Validation", Tags: []string{"validacao"}},
		{Letter: "M", Term: "Message", Tags: []string{"mensagem"}},
	}}
	feat := `# language: pt
Funcionalidade: Banner

  @mensagem @IPSBX-V02 @IPSBX-M01 @P3
  Cenário: Processing usa spinner e rótulo "Processando..."
    Então devo ver o rótulo "Processando..."
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkCenarioTipoAlinhado(feat, n, "", nil, cfg); v != Pass {
		t.Errorf("esperava Pass — a tag descreve o código secundário: %v %s", v, detail)
	}
}

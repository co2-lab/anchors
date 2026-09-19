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

func TestScenarioTypeAligned_NaoFeaturePula(t *testing.T) {
	t.Run("STASC-B01: Non-feature artifacts skip confrontation", func(t *testing.T) {})
	t.Run("STASC-I01: Classification alignment is evaluated only on feature artifacts", func(t *testing.T) {})
	t.Run("STASC-X03: Specifications and test files are not inspected or modified", func(t *testing.T) {})
	cfg := cfgComTags()
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindTest} {
		n := mapx.Node{ID: "arquivo", Kind: k}
		v, msg := checkScenarioTypeAligned("qualquer", n, "", nil, cfg)
		if v != Skip {
			t.Fatalf("esperava Skip para kind %v, veio %v (%s)", k, v, msg)
		}
	}
}

func TestScenarioTypeAligned_ConfigSemRuleTypesPula(t *testing.T) {
	t.Run("STASC-B02: Missing or empty rule types in configuration skip confrontation", func(t *testing.T) {})
	feat := "Feature: X\nScenario: y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkScenarioTypeAligned(feat, n, "", nil, nil); v != Skip {
		t.Fatalf("esperava Skip com cfg nil, veio %v", v)
	}
	if v, _ := checkScenarioTypeAligned(feat, n, "", nil, &config.Config{}); v != Skip {
		t.Fatalf("esperava Skip com RuleTypes vazias, veio %v", v)
	}
}

func TestScenarioTypeAligned_SemMapaFicaEmSilencio(t *testing.T) {
	t.Run("STASC-B03: Configuration without tag mappings skips confrontation", func(t *testing.T) {})
	t.Run("STASC-I02: Absence of tag mappings prevents speculative enforcement", func(t *testing.T) {})
	cfg := &config.Config{RuleTypes: []config.RuleType{{Letter: "S", Term: "State"}}}
	codeS01 := "LGSTX" + "-S01"
	feat := "# language: pt\n\n  @comportamento @" + codeS01 + " @P2\n  Cenário: x\n    Então y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkScenarioTypeAligned(feat, n, "", nil, cfg); v != Skip {
		t.Errorf("esperava Skip sem mapa tag→letra, veio %v", v)
	}
}

func TestScenarioTypeAligned_FeatureSemCenariosComCodigoPula(t *testing.T) {
	t.Run("STASC-B04: Features without coded scenarios skip confrontation", func(t *testing.T) {})
	feat := "# language: pt\nFuncionalidade: Vazia\n  Cenário: sem código\n    Então y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags()); v != Skip {
		t.Errorf("esperava Skip para feature sem cenários com código, veio %v", v)
	}
}

func TestScenarioTypeAligned_CodigosSemLetraSaoIgnorados(t *testing.T) {
	t.Run("STASC-B05: Codes lacking a recognized rule letter are ignored", func(t *testing.T) {})
	codeS02 := "BUGEX" + "-S02"
	codeS02Var := "BUGEX" + "-S02#01"
	codeDS := "BUGEX" + "-DS-mode-standalone"
	codeVR := "BUGEX" + "-VR"
	casos := map[string]string{
		codeS02:    "S",
		codeS02Var: "S",
		codeDS:     "",
		codeVR:     "",
	}
	for entrada, quer := range casos {
		if got := codeLetter(entrada); got != quer {
			t.Errorf("codeLetter(%q) = %q, queria %q", entrada, got, quer)
		}
	}
	feat := "# language: pt\nFuncionalidade: X\n  @comportamento @" + codeDS + "\n  Cenário: ds\n    Então x\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags()); v != Pass {
		t.Fatalf("código sem letra deve ser ignorado e passar: %v", v)
	}
}

func TestScenarioTypeAligned_TagsNaoRegistradasIgnoradas(t *testing.T) {
	t.Run("STASC-B06: Unregistered scenario tags are ignored", func(t *testing.T) {})
	t.Run("STASC-X01: The gate does not mandate classification tags on every scenario", func(t *testing.T) {})
	codeS01 := "LGSTX" + "-S01"
	feat := "# language: pt\nFuncionalidade: X\n  @" + codeS01 + " @unit-level @P2 @sem-tag-de-tipo\n  Cenário: apenas tags de nivel\n    Então y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags()); v != Pass {
		t.Fatalf("esperava Pass ignorando tags não registradas, veio %v: %s", v, detail)
	}
}

func TestScenarioTypeAligned_PassaQuandoConcordam(t *testing.T) {
	t.Run("STASC-B07: Scenarios whose classification tags align with code letters pass", func(t *testing.T) {})
	codeS01 := "LGSTX" + "-S01"
	codeB01 := "LGSTX" + "-B01"
	feat := `# language: pt
Funcionalidade: Seção

  @estado @` + codeS01 + ` @P2
  Cenário: Seção padrão exibe número e título
    Então devo ver o título

  @comportamento @` + codeB01 + ` @P2
  Cenário: Tocar no link abre os termos
    Quando eu toco no link
    Então devo ver os termos
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags()); v != Pass {
		t.Errorf("esperava Pass, veio %v: %s", v, detail)
	}
}

func TestScenarioTypeAligned_TagPodeCaberEmMaisDeUmaLetra(t *testing.T) {
	t.Run("STASC-B08: A tag mapped to multiple rule letters passes if any matches", func(t *testing.T) {})
	t.Run("STASC-X04: Tags are permitted to map across multiple rule letters", func(t *testing.T) {})
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "S", Term: "State", Tags: []string{"estado", "estado-dado"}},
		{Letter: "V", Term: "Validation", Tags: []string{"validacao", "estado-dado"}},
	}}
	codeV01 := "PRBDX" + "-V01"
	codeS01 := "PRBDX" + "-S01"
	feat := `# language: pt
Funcionalidade: Badge

  @estado-dado @` + codeV01 + ` @P2
  Cenário: Variante high define fundo, cor e rótulo
    Então devo ver o texto "Alta"

  @estado-dado @` + codeS01 + ` @P2
  Cenário: Sem prioridade não exibe a pílula
    Então não devo ver a pílula
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfg); v != Pass {
		t.Errorf("esperava Pass — a tag cabe nas duas letras: %v %s", v, detail)
	}
}

func TestScenarioTypeAligned_TagPodeDescreverCodigoSecundario(t *testing.T) {
	t.Run("STASC-B09: A multi-coded scenario passes when a tag matches any code", func(t *testing.T) {})
	t.Run("STASC-I04: Secondary requirement codes prevent false mismatch reporting", func(t *testing.T) {})
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "V", Term: "Validation", Tags: []string{"validacao"}},
		{Letter: "M", Term: "Message", Tags: []string{"mensagem"}},
	}}
	codeV02 := "IPSBX" + "-V02"
	codeM01 := "IPSBX" + "-M01"
	feat := `# language: pt
Funcionalidade: Banner

  @mensagem @` + codeV02 + ` @` + codeM01 + ` @P3
  Cenário: Processing usa spinner e rótulo "Processando..."
    Então devo ver o rótulo "Processando..."
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfg); v != Pass {
		t.Errorf("esperava Pass — a tag descreve o código secundário: %v %s", v, detail)
	}
}

func TestScenarioTypeAligned_AcusaTagQueDiscordaDaLetra(t *testing.T) {
	t.Run("STASC-B10: A scenario tag disagreeing with the code rule letter returns Pending", func(t *testing.T) {})
	t.Run("STASC-I03: Mismatched scenario types return Pending rather than Fail", func(t *testing.T) {})
	t.Run("STASC-X02: The gate does not decide whether tag or code is erroneous", func(t *testing.T) {})
	codeS01 := "LGSTX" + "-S01"
	feat := `# language: pt
Funcionalidade: Seção

  @comportamento @` + codeS01 + `#02 @nivel-integration @P2
  Cenário: Intro presente é renderizada em itálico
    Quando o componente é renderizado
    Então devo ver o texto de intro
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags())
	if v != Pending {
		t.Fatalf("esperava Pending, veio %v: %s", v, detail)
	}
	if !strings.Contains(detail, codeS01) || !strings.Contains(detail, "comportamento") {
		t.Errorf("o detalhe não nomeia o código nem a tag: %s", detail)
	}
}

func TestScenarioTypeAligned_DetalhesDoVereditoPending(t *testing.T) {
	t.Run("STASC-B11: The Pending verdict cites details of the type divergence", func(t *testing.T) {})
	codeS01 := "LGSTX" + "-S01"
	tituloLongo := "Um titulo deliberadamente longo para testar a abreviacao de texto pelo formatador"
	feat := `# language: pt
Funcionalidade: Detalhes

  @comportamento @` + codeS01 + `
  Cenário: ` + tituloLongo + `
    Quando o componente é renderizado
    Então devo ver o texto
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags())
	if v != Pending {
		t.Fatalf("esperava Pending, veio %v: %s", v, detail)
	}
	if !strings.Contains(detail, codeS01) {
		t.Errorf("deve conter o código: %s", detail)
	}
	if !strings.Contains(detail, "S") {
		t.Errorf("deve conter a letra do código S: %s", detail)
	}
	if !strings.Contains(detail, "comportamento") {
		t.Errorf("deve conter o nome da tag: %s", detail)
	}
	if !strings.Contains(detail, "B") {
		t.Errorf("deve conter as letras permitidas para a tag: %s", detail)
	}
	if !strings.Contains(detail, "…") {
		t.Errorf("deve encurtar o título longo com reticências: %s", detail)
	}
}

func TestScenarioTypeAligned_MultiplosAchadosOrdenados(t *testing.T) {
	t.Run("STASC-B12: Multiple mismatch findings are sorted deterministically", func(t *testing.T) {})
	codeS01 := "LGSTX" + "-S01"
	codeS02 := "LGSTX" + "-S02"
	feat := `# language: pt
Funcionalidade: Multiplos

  @comportamento @` + codeS02 + `
  Cenário: ZZZ segundo cenario
    Então z

  @comportamento @` + codeS01 + `
  Cenário: AAA primeiro cenario
    Então a
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	v, detail := checkScenarioTypeAligned(feat, n, "", nil, cfgComTags())
	if v != Pending {
		t.Fatalf("esperava Pending, veio %v: %s", v, detail)
	}
	if !strings.Contains(detail, "2") {
		t.Errorf("esperava contagem de 2 achados: %s", detail)
	}
	idx1 := strings.Index(detail, codeS01)
	idx2 := strings.Index(detail, codeS02)
	if idx1 == -1 || idx2 == -1 || idx1 > idx2 {
		t.Errorf("achados devem ser ordenados deterministicamente: %s", detail)
	}
}

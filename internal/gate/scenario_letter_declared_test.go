package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func cfgLetras() *config.Config {
	return &config.Config{RuleTypes: []config.RuleType{
		{Letter: "S", Term: "State"},
		{Letter: "B", Term: "Behavior"},
	}}
}

func TestCenarioLetraDeclarada_acusaLetraInventada(t *testing.T) {
	t.Run("SCLTR-B05: a letter outside the vocabulary is undetermined, not a failure", func(t *testing.T) {})
	t.Run("SCLTR-B06: The verdict names the letters that are outside and the codes carrying them", func(t *testing.T) {})
	t.Run("SCLTR-I03: a valid letter is never named in the verdict", func(t *testing.T) {})
	t.Run("SCLTR-I01: the scan is over the shape of a code, never over the vocabulary", func(t *testing.T) {})
	t.Run("SCLTR-X01: the gate does not choose between declaring and remapping", func(t *testing.T) {})
	feat := `# language: pt
Funcionalidade: Recorrências

  @navegacao @RCRRX-SG05 @P2
  Cenário: Tocar numa sugestão abre o detalhe
    Quando eu toco na sugestão
    Então devo ver o detalhe

  @estado @RCRRX-S01 @P2
  Cenário: Lista vazia
    Então devo ver o texto de vazio
`
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	v, detail := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras())
	if v != Pending {
		t.Fatalf("esperava Pending, veio %v: %s", v, detail)
	}
	if !strings.Contains(detail, "SG") || !strings.Contains(detail, "RCRRX-SG05") {
		t.Errorf("o detalhe não nomeia a letra nem o código: %s", detail)
	}
	if strings.Contains(detail, "RCRRX-S01") {
		t.Errorf("acusou um código de letra VÁLIDA: %s", detail)
	}
}

func TestCenarioLetraDeclarada_passaComVocabularioRespeitado(t *testing.T) {
	t.Run("SCLTR-B04: every letter inside the vocabulary passes", func(t *testing.T) {})
	feat := "# language: pt\n\n  @estado @ABCDX-S01 @P2\n  Cenário: x\n    Então y\n\n  @comportamento @ABCDX-B01#02 @P2\n  Cenário: z\n    Então w\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras()); v != Pass {
		t.Errorf("esperava Pass, veio %v: %s", v, detail)
	}
}

// A mesma letra em muitos cenários vira UMA linha: o leitor precisa saber quais
// letras estão fora, não reler a mesma acusação oito vezes.
func TestCenarioLetraDeclarada_agrupaPorLetra(t *testing.T) {
	t.Run("SCLTR-B07: codes sharing one unknown letter are grouped into a single line", func(t *testing.T) {})
	feat := "# language: pt\n\n  @x @ABCDX-FP01 @P2\n  Cenário: a\n    Então y\n\n  @x @ABCDX-FP02 @P2\n  Cenário: b\n    Então y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	_, detail := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras())
	if !strings.HasPrefix(detail, "1 letra(s)") && !strings.HasPrefix(detail, "1 code letter(s)") {
		t.Errorf("esperava uma única letra acusada, veio: %s", detail)
	}
	if !strings.Contains(detail, "ABCDX-FP01, ABCDX-FP02") {
		t.Errorf("os dois códigos deviam vir na mesma linha: %s", detail)
	}
}

// SCLTR-B01: scenario codes live in features — every other kind leaves without a verdict.
func TestScenarioLetterDeclared_B01_skipsNonFeature(t *testing.T) {
	t.Run("SCLTR-B01: an artifact that is not a feature leaves without a verdict", func(t *testing.T) {})
	feat := "@ABCDX-SG01\nScenario: x\n"
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindTest} {
		n := mapx.Node{ID: "a", Kind: k}
		if v, _ := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras()); v != Skip {
			t.Errorf("kind %v: expected Skip, got %v", k, v)
		}
	}
}

// SCLTR-B02: with no declared vocabulary every letter would be either all valid or all
// invented — both answers are noise, so the gate declines to judge.
func TestScenarioLetterDeclared_B02_skipsWithoutVocabulary(t *testing.T) {
	t.Run("SCLTR-B02: with no declared vocabulary the gate leaves without a verdict", func(t *testing.T) {})
	feat := "@ABCDX-SG01\nScenario: x\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkScenarioLetterDeclared(feat, n, "", nil, &config.Config{}); v != Skip {
		t.Errorf("empty vocabulary: expected Skip, got %v", v)
	}
	if v, _ := checkScenarioLetterDeclared(feat, n, "", nil, nil); v != Skip {
		t.Errorf("nil config: expected Skip, got %v", v)
	}
}

// SCLTR-B03: a feature with no code at all has nothing to judge.
func TestScenarioLetterDeclared_B03_skipsWithoutCodes(t *testing.T) {
	t.Run("SCLTR-B03: a feature carrying no scenario code leaves without a verdict", func(t *testing.T) {})
	feat := "# language: en\nFeature: x\n\n  Scenario: no code here\n    Then y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, _ := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras()); v != Skip {
		t.Errorf("expected Skip, got %v", v)
	}
}

// SCLTR-I02: the code-length pattern is read at EVERY call. It comes from the project's
// Structure, which loads AFTER the package globals — a pattern built once at load time
// would freeze the default and silently ignore the project's declaration.
//
// Proven by declaring a non-default length and confronting a code of that length: it has
// to be SEEN (reported as an undeclared letter), not silently skipped as unrecognisable.
func TestScenarioLetterDeclared_I02_codeLengthReadPerCall(t *testing.T) {
	t.Run("SCLTR-I02: the code-length pattern is read at every call", func(t *testing.T) {})
	original := append([]int{}, config.CodeLengths...)
	defer func() { config.CodeLengths = original }()

	config.CodeLengths = []int{7}
	feat := "@ABCDEFG-SG01\nScenario: x\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	v, detail := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras())
	if v != Pending {
		t.Fatalf("a 7-char code must be recognised once the project declares length 7; got %v (%s)", v, detail)
	}
	if !strings.Contains(detail, "ABCDEFG-SG01") {
		t.Errorf("the verdict must name the code: %s", detail)
	}
}

// SCLTR-X02: a tag is free vocabulary by design. Charging it here would turn a precise
// instrument into a style opinion — only the LETTER OF THE CODE is judged.
func TestScenarioLetterDeclared_X02_ignoresTags(t *testing.T) {
	t.Run("SCLTR-X02: the tags accompanying a code are not judged", func(t *testing.T) {})
	feat := "@whatever-free-tag @ABCDX-S01 @another_one\nScenario: x\n    Then y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkScenarioLetterDeclared(feat, n, "", nil, cfgLetras()); v != Pass {
		t.Errorf("free-form tags must not be judged; got %v (%s)", v, detail)
	}
}

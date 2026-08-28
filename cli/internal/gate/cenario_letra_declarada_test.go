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
	v, detail := checkCenarioLetraDeclarada(feat, n, "", nil, cfgLetras())
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
	feat := "# language: pt\n\n  @estado @ABCDX-S01 @P2\n  Cenário: x\n    Então y\n\n  @comportamento @ABCDX-B01#02 @P2\n  Cenário: z\n    Então w\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	if v, detail := checkCenarioLetraDeclarada(feat, n, "", nil, cfgLetras()); v != Pass {
		t.Errorf("esperava Pass, veio %v: %s", v, detail)
	}
}

// A mesma letra em muitos cenários vira UMA linha: o leitor precisa saber quais
// letras estão fora, não reler a mesma acusação oito vezes.
func TestCenarioLetraDeclarada_agrupaPorLetra(t *testing.T) {
	feat := "# language: pt\n\n  @x @ABCDX-FP01 @P2\n  Cenário: a\n    Então y\n\n  @x @ABCDX-FP02 @P2\n  Cenário: b\n    Então y\n"
	n := mapx.Node{ID: "a.feature", Kind: mapx.KindFeature}
	_, detail := checkCenarioLetraDeclarada(feat, n, "", nil, cfgLetras())
	if !strings.HasPrefix(detail, "1 letra(s)") {
		t.Errorf("esperava uma única letra acusada, veio: %s", detail)
	}
	if !strings.Contains(detail, "ABCDX-FP01, ABCDX-FP02") {
		t.Errorf("os dois códigos deviam vir na mesma linha: %s", detail)
	}
}

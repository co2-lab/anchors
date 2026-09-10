package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaAsserts(t *testing.T, feature, lang string) (Verdict, string) {
	t.Helper()
	cfg := &config.Config{Dialect: &config.Dialect{GherkinLanguage: lang}}
	return checkScenarioAsserts(feature, mapx.Node{Kind: mapx.KindFeature}, "", nil, cfg)
}

// A tautologia real que motivou o gate: 12 de 12 cenários de uma feature terminavam em
// "Então o efeito XXXXX-B0n se verifica" — o passo repete o código da regra e não afirma
// resultado nenhum. Medido: 359 cenários assim num projeto real, todos vindos do template.
func TestScenarioAssertsPegaTautologia(t *testing.T) {
	feature := `@MTVRX
Funcionalidade: versionamento

  @MTVRX-B01 @nivel-unit
  Cenário: carry-forward
    Dado versões em janeiro
    Quando leio maio
    Então o efeito MTVRX-B01 se verifica

  @MTVRX-B05 @nivel-unit
  Cenário: cada chave versiona sozinha
    Dado duas chaves
    Quando edito uma
    Então o efeito MTVRX-B05 se verifica
`
	v, d := rodaAsserts(t, feature, "pt")
	if v != Fail {
		t.Fatalf("tautologia deveria reprovar, foi %s (%s)", v, d)
	}
	for _, c := range []string{"MTVRX-B01", "MTVRX-B05"} {
		if !strings.Contains(d, c) {
			t.Errorf("não nomeou o cenário %s: %s", c, d)
		}
	}
}

// O que NÃO pode acontecer: acusar um passo que afirma resultado de verdade. Um gate que
// gera falso positivo é desligado, e aí não protege nada.
func TestScenarioAssertsNaoAcusaPassoQueAfirma(t *testing.T) {
	casos := map[string]string{
		"valor concreto":        "    Então o valor lido é \"ACME\"",
		"resultado observável":  "    Então a chave não aparece no conjunto do mês",
		"cita a regra E afirma": "    Então o valor de maio é \"ACME\", como manda MTVRX-B01",
		"comparação":            "    Então a leitura de março é igual à de antes da edição",
		"negativa":              "    Então nenhuma versão de remoção sobra no histórico",
	}
	for nome, passo := range casos {
		t.Run(nome, func(t *testing.T) {
			feature := "@X\nFuncionalidade: y\n\n  @MTVRX-B01 @nivel-unit\n  Cenário: z\n" + passo + "\n"
			if v, d := rodaAsserts(t, feature, "pt"); v != Pass {
				t.Fatalf("passo que AFIRMA foi acusado (%s): %s", v, d)
			}
		})
	}
}

// As variações da mesma tautologia — o autor troca as palavras de ligação, a forma
// continua vazia.
func TestScenarioAssertsVariacoesDaTautologia(t *testing.T) {
	for _, passo := range []string{
		"    Então o efeito MTVRX-B01 se verifica",
		"    Então a regra MTVRX-B01 é aplicada",
		"    Então MTVRX-B01",
		"    Então o comportamento MTVRX-B01 vale",
		"    Então MTVRX-B01 se verifica",
	} {
		t.Run(strings.TrimSpace(passo), func(t *testing.T) {
			feature := "@X\nFuncionalidade: y\n\n  @MTVRX-B01 @nivel-unit\n  Cenário: z\n" + passo + "\n"
			if v, _ := rodaAsserts(t, feature, "pt"); v != Fail {
				t.Fatalf("variação da tautologia passou: %q", passo)
			}
		})
	}
}

// Funciona em qualquer idioma do Gherkin: um repositório pode ter features herdadas
// noutra língua, e ficar cego nelas é reportar verde sobre o que não se enxerga.
func TestScenarioAssertsEntreIdiomas(t *testing.T) {
	casos := map[string]struct{ lang, tautologico, bom string }{
		"en": {"en", "    Then the effect MTVRX-B01 is verified", "    Then the value read is \"ACME\""},
		"es": {"es", "    Entonces el efecto MTVRX-B01 se verifica", "    Entonces el valor es \"ACME\""},
		"pt": {"pt", "    Então o efeito MTVRX-B01 se verifica", "    Então o valor lido é \"ACME\""},
	}
	for nome, c := range casos {
		t.Run(nome+" — tautológico reprova", func(t *testing.T) {
			f := "@X\nFuncionalidade: y\n\n  @MTVRX-B01 @unit\n  Cenário: z\n" + c.tautologico + "\n"
			if v, d := rodaAsserts(t, f, c.lang); v != Fail {
				t.Fatalf("não pegou em %s: %s (%s)", nome, v, d)
			}
		})
		t.Run(nome+" — afirmativo passa", func(t *testing.T) {
			f := "@X\nFuncionalidade: y\n\n  @MTVRX-B01 @unit\n  Cenário: z\n" + c.bom + "\n"
			if v, d := rodaAsserts(t, f, c.lang); v != Pass {
				t.Fatalf("falso positivo em %s: %s (%s)", nome, v, d)
			}
		})
	}
}

// Só a feature descreve cenário; os demais artefatos não têm passo de resultado.
func TestScenarioAssertsSoFeature(t *testing.T) {
	f := "  Então o efeito MTVRX-B01 se verifica\n"
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindTest} {
		cfg := &config.Config{Dialect: &config.Dialect{GherkinLanguage: "pt"}}
		if v, _ := checkScenarioAsserts(f, mapx.Node{Kind: k}, "", nil, cfg); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}

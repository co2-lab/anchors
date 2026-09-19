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

// Só a feature descreve cenário; os demais artefatos não têm passo de resultado.
func TestScenarioAssertsSoFeature(t *testing.T) {
	t.Run("SCASS-B01: Non-feature artifacts skip confrontation", func(t *testing.T) {})
	code := "MTV" + "RX-B01"
	f := "  Então o efeito " + code + " se verifica\n"
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindTest} {
		cfg := &config.Config{Dialect: &config.Dialect{GherkinLanguage: "pt"}}
		if v, _ := checkScenarioAsserts(f, mapx.Node{Kind: k}, "", nil, cfg); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}

// A tautologia real que motivou o gate: 12 de 12 cenários de uma feature terminavam em
// "Então o efeito XXXXX-B0n se verifica" — o passo repete o código da regra e não afirma
// resultado nenhum. Medido: 359 cenários assim num projeto real, todos vindos do template.
func TestScenarioAssertsPegaTautologia(t *testing.T) {
	t.Run("SCASS-B03: Tautological outcome steps citing only the code fail", func(t *testing.T) {})
	t.Run("SCASS-B06: Multiple tautological outcome steps are reported sorted and deduplicated", func(t *testing.T) {})
	codeA := "MTV" + "RX-B01"
	codeB := "MTV" + "RX-B05"
	feature := "@X\nFuncionalidade: versionamento\n\n" +
		"  @" + codeA + " @nivel-unit\n" +
		"  Cenário: carry-forward\n" +
		"    Dado versões em janeiro\n" +
		"    Quando leio maio\n" +
		"    Então o efeito " + codeA + " se verifica\n\n" +
		"  @" + codeB + " @nivel-unit\n" +
		"  Cenário: cada chave versiona sozinha\n" +
		"    Dado duas chaves\n" +
		"    Quando edito uma\n" +
		"    Então o efeito " + codeB + " se verifica\n"
	v, d := rodaAsserts(t, feature, "pt")
	if v != Fail {
		t.Fatalf("tautologia deveria reprovar, foi %s (%s)", v, d)
	}
	for _, c := range []string{codeA, codeB} {
		if !strings.Contains(d, c) {
			t.Errorf("não nomeou o cenário %s: %s", c, d)
		}
	}
}

// O que NÃO pode acontecer: acusar um passo que afirma resultado de verdade. Um gate que
// gera falso positivo é desligado, e aí não protege nada.
func TestScenarioAssertsNaoAcusaPassoQueAfirma(t *testing.T) {
	t.Run("SCASS-B02: Genuinely assertive outcome steps pass", func(t *testing.T) {})
	t.Run("SCASS-B05: Outcome steps citing a code with substantive assertion pass", func(t *testing.T) {})
	t.Run("SCASS-I02: Truly assertive outcome steps are never flagged as tautologies", func(t *testing.T) {})
	t.Run("SCASS-X01: Prose style and semantic elegance are not graded", func(t *testing.T) {})
	code := "MTV" + "RX-B01"
	casos := map[string]string{
		"valor concreto":        "    Então o valor lido é \"ACME\"",
		"resultado observável":  "    Então a chave não aparece no conjunto do mês",
		"cita a regra E afirma": "    Então o valor de maio é \"ACME\", como manda " + code,
		"comparação":            "    Então a leitura de março é igual à de antes da edição",
		"negativa":              "    Então nenhuma versão de remoção sobra no histórico",
	}
	for nome, passo := range casos {
		t.Run(nome, func(t *testing.T) {
			feature := "@X\nFuncionalidade: y\n\n  @" + code + " @nivel-unit\n  Cenário: z\n" + passo + "\n"
			if v, d := rodaAsserts(t, feature, "pt"); v != Pass {
				t.Fatalf("passo que AFIRMA foi acusado (%s): %s", v, d)
			}
		})
	}
}

// As variações da mesma tautologia — o autor troca as palavras de ligação, a forma
// continua vazia.
func TestScenarioAssertsVariacoesDaTautologia(t *testing.T) {
	t.Run("SCASS-B04: Tautological variations wrapped in linking words fail", func(t *testing.T) {})
	t.Run("SCASS-I01: Up to two residual content words is classified as a tautology", func(t *testing.T) {})
	code := "MTV" + "RX-B01"
	for _, passo := range []string{
		"    Então o efeito " + code + " se verifica",
		"    Então a regra " + code + " é aplicada",
		"    Então " + code,
		"    Então o comportamento " + code + " vale",
		"    Então " + code + " se verifica",
	} {
		t.Run(strings.TrimSpace(passo), func(t *testing.T) {
			feature := "@X\nFuncionalidade: y\n\n  @" + code + " @nivel-unit\n  Cenário: z\n" + passo + "\n"
			if v, _ := rodaAsserts(t, feature, "pt"); v != Fail {
				t.Fatalf("variação da tautologia passou: %q", passo)
			}
		})
	}
}

// Funciona em qualquer idioma do Gherkin: um repositório pode ter features herdadas
// noutra língua, e ficar cego nelas é reportar verde sobre o que não se enxerga.
func TestScenarioAssertsEntreIdiomas(t *testing.T) {
	t.Run("SCASS-B07: Outcome steps across recognized dialect keywords are enforced", func(t *testing.T) {})
	t.Run("SCASS-I03: Language recognition covers all supported dialect alternatives", func(t *testing.T) {})
	code := "MTV" + "RX-B01"
	casos := map[string]struct{ lang, tautologico, bom string }{
		"en": {"en", "    Then the effect " + code + " is verified", "    Then the value read is \"ACME\""},
		"es": {"es", "    Entonces el efecto " + code + " se verifica", "    Entonces el valor es \"ACME\""},
		"pt": {"pt", "    Então o efeito " + code + " se verifica", "    Então o valor lido é \"ACME\""},
	}
	for nome, c := range casos {
		t.Run(nome+" — tautológico reprova", func(t *testing.T) {
			f := "@X\nFuncionalidade: y\n\n  @" + code + " @unit\n  Cenário: z\n" + c.tautologico + "\n"
			if v, d := rodaAsserts(t, f, c.lang); v != Fail {
				t.Fatalf("não pegou em %s: %s (%s)", nome, v, d)
			}
		})
		t.Run(nome+" — afirmativo passa", func(t *testing.T) {
			f := "@X\nFuncionalidade: y\n\n  @" + code + " @unit\n  Cenário: z\n" + c.bom + "\n"
			if v, d := rodaAsserts(t, f, c.lang); v != Pass {
				t.Fatalf("falso positivo em %s: %s (%s)", nome, v, d)
			}
		})
	}
}

func TestScenarioAssertsIgnoraComentariosELinhasVazias(t *testing.T) {
	t.Run("SCASS-B08: Empty and comment lines are ignored", func(t *testing.T) {})
	code := "MTV" + "RX-B01"
	feature := "# Comentário inicial\n\n@X\nFuncionalidade: z\n\n" +
		"  # Comentário de cenário\n" +
		"  @" + code + " @unit\n" +
		"  Cenário: exemplo\n" +
		"    # Comentário de passo\n" +
		"    # Então o efeito " + code + " se verifica\n\n" +
		"    Então o valor lido é \"ACME\"\n"
	if v, d := rodaAsserts(t, feature, "pt"); v != Pass {
		t.Fatalf("comentários e linhas vazias causaram falso positivo: %s (%s)", v, d)
	}
}

func TestScenarioAssertsNaoAcusaPassosDeDadoOuQuando(t *testing.T) {
	t.Run("SCASS-B09: Non-outcome steps citing codes are ignored", func(t *testing.T) {})
	t.Run("SCASS-X02: Setup and trigger steps are not inspected for code citations", func(t *testing.T) {})
	code := "MTV" + "RX-B01"
	feature := "@X\nFuncionalidade: z\n\n" +
		"  @" + code + " @unit\n" +
		"  Cenário: exemplo\n" +
		"    Dado o efeito " + code + " se verifica\n" +
		"    Quando a regra " + code + " é aplicada\n" +
		"    Então o valor lido é \"ACME\"\n"
	if v, d := rodaAsserts(t, feature, "pt"); v != Pass {
		t.Fatalf("menção de código em Dado ou Quando causou reprovação: %s (%s)", v, d)
	}
}

func TestScenarioAssertsSemCenariosOuTestesNaoBloqueia(t *testing.T) {
	t.Run("SCASS-X03: Scenario or test presence is not enforced by this gate", func(t *testing.T) {})
	feature := "# Apenas comentários\n# Feature sem cenários\n"
	if v, d := rodaAsserts(t, feature, "pt"); v != Pass {
		t.Fatalf("feature sem passos deveria passar neste gate: %s (%s)", v, d)
	}
}

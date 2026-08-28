package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// monta spec+feature em disco e o grafo que as liga, como o `check` faz de verdade.
func rodaSpecFeature(t *testing.T, spec, feature string) (Verdict, string) {
	t.Helper()
	dir := t.TempDir()
	if feature != "" {
		if err := os.WriteFile(filepath.Join(dir, "x.feature"), []byte(feature), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &mapx.Graph{}
	if feature != "" {
		g.Edges = []mapx.Edge{{From: "x.spec.md", To: "x.feature", Type: "covered-by"}}
	}
	return checkSpecFeatureMatch(spec, mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}, dir, g, nil)
}

// O caso real que motivou o gate: uma spec do app de referência declarava OCHIX-X02 e a feature irmã não
// tinha cenário nenhum com essa tag. Todos os 12 gates ficavam VERDES — a spec tem código,
// a feature existe, a feature bate com o teste. O requisito não era de ninguém.
func TestSpecFeatureMatchRequisitoSemCenario(t *testing.T) {
	spec := `# Spec

## Regras
### OCHIX-B01 — lista as ocorrências
### OCHIX-X02 — não busca dado, recebe por prop
`
	feature := "@OCHIX\nFuncionalidade: histórico\n\n  @OCHIX-B01 @nivel-unit\n  Cenário: lista\n"

	v, d := rodaSpecFeature(t, spec, feature)
	if v != Fail {
		t.Fatalf("requisito sem cenário deveria reprovar, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "OCHIX-X02") {
		t.Errorf("não nomeou o requisito descoberto: %s", d)
	}
	if strings.Contains(d, "OCHIX-B01") {
		t.Errorf("acusou um requisito que TEM cenário: %s", d)
	}
}

func TestSpecFeatureMatchTudoCoberto(t *testing.T) {
	spec := "# Spec\n\n### AAAAX-B01 — x\n### AAAAX-B02 — y\n"
	feature := "@AAAAX\nFuncionalidade: x\n\n  @AAAAX-B01 @nivel-unit\n  Cenário: a\n\n  @AAAAX-B02 @nivel-unit\n  Cenário: b\n"
	if v, d := rodaSpecFeature(t, spec, feature); v != Pass {
		t.Fatalf("tudo coberto deveria passar, foi %s (%s)", v, d)
	}
}

// DEFINIR ≠ CITAR. Uma spec cita códigos de outras unidades o tempo todo (Tabela de
// Dependências, notas, referências cruzadas) e não contrai obrigação por isso. Sem essa
// distinção o gate viraria um gerador de ruído e seria desligado.
func TestSpecFeatureMatchCitacaoNaoObriga(t *testing.T) {
	spec := `# Spec

### AAAAX-B01 — o único requisito desta spec

## Notas
Esta unidade é consumida pela tela que implementa ` + "`BBBBX-S01`" + ` e depende do
comportamento ` + "`CCCCX-B09`" + ` descrito noutra spec.

## Dependências
| Cód | Arquivo | Método | Camada |
| --- | --- | --- | --- |
| DEP1 | x/y.ts | ` + "`faz`" + ` | logic |
`
	feature := "@AAAAX\nFuncionalidade: x\n\n  @AAAAX-B01 @nivel-unit\n  Cenário: a\n"
	v, d := rodaSpecFeature(t, spec, feature)
	if v != Pass {
		t.Fatalf("códigos CITADOS não são obrigação desta spec, foi %s (%s)", v, d)
	}
}

// Opt-out honesto por REQUISITO (CONCEPT §5.1): vale com razão, não vale nu.
func TestSpecFeatureMatchDispensaExigeRazao(t *testing.T) {
	feature := "@AAAAX\nFuncionalidade: x\n\n  @AAAAX-B01 @nivel-unit\n  Cenário: a\n"

	comRazão := "# Spec\n\n### AAAAX-B01 — x\n### AAAAX-X02 — limite de camada @no-scenario: restrição estrutural, provada por check-arch e não por cenário\n"
	if v, d := rodaSpecFeature(t, comRazão, feature); v != Pass {
		t.Errorf("dispensa COM razão deveria passar, foi %s (%s)", v, d)
	}

	nu := "# Spec\n\n### AAAAX-B01 — x\n### AAAAX-X02 — limite @no-scenario:\n"
	if v, _ := rodaSpecFeature(t, nu, feature); v != Fail {
		t.Errorf("marcador NU deveria continuar reprovando, foi %s", v)
	}
}

// `@no-feature` na spec ARRASTA a dispensa para todos os requisitos: sem feature, nenhum
// deles pode ter cenário. Sem este arrasto o autor precisa repetir `@no-scenario` em cada
// linha da tabela para dizer o que a spec já disse uma vez — e as marcações podem divergir.
func TestSpecFeatureMatchNoFeatureArrastaTodosOsRequisitos(t *testing.T) {
	// Feature vazia (só cabeçalho) + requisitos SEM `@no-scenario`: sem o arrasto isto
	// reprovaria acusando os dois requisitos como descobertos.
	spec := "# Spec\n\n@no-feature: gateway que só repassa — nada observável por cenário\n\n### AAAAX-B01 — x\n### AAAAX-B02 — y\n"
	feature := "@AAAAX\nFuncionalidade: x\n"

	if v, d := rodaSpecFeature(t, spec, feature); v != Skip {
		t.Errorf("`@no-feature` dispensa o cenário de todo requisito, foi %s (%s)", v, d)
	}
}

// O arrasto NÃO pode virar um jeito silencioso de calar o gate: sem a tag, os mesmos
// requisitos descobertos continuam reprovando. É o par negativo do teste acima — sem ele,
// um `Skip` incondicional passaria despercebido.
func TestSpecFeatureMatchSemNoFeatureContinuaCobrando(t *testing.T) {
	spec := "# Spec\n\n### AAAAX-B01 — x\n### AAAAX-B02 — y\n"
	feature := "@AAAAX\nFuncionalidade: x\n"

	if v, _ := rodaSpecFeature(t, spec, feature); v != Fail {
		t.Errorf("sem a dispensa o requisito descoberto deve reprovar, foi %s", v)
	}
}

// A dispensa exige RAZÃO, igual ao `@no-scenario`: um `@no-feature` nu não arrasta nada,
// senão o marcador vira um interruptor para desligar o gate sem prestar contas.
func TestSpecFeatureMatchNoFeatureNuNaoArrasta(t *testing.T) {
	spec := "# Spec\n\n@no-feature:\n\n### AAAAX-B01 — x\n"
	feature := "@AAAAX\nFuncionalidade: x\n"

	if v, _ := rodaSpecFeature(t, spec, feature); v != Fail {
		t.Errorf("marcador NU não deveria dispensar, foi %s", v)
	}
}

// Cada gate acusa UMA coisa: a ausência da feature é do trinca-completa. Acusar aqui
// também faria o mesmo defeito aparecer duas vezes no relatório.
func TestSpecFeatureMatchSemFeatureEhDoOutroGate(t *testing.T) {
	spec := "# Spec\n\n### AAAAX-B01 — x\n"
	if v, d := rodaSpecFeature(t, spec, ""); v != Skip {
		t.Fatalf("spec sem feature deveria ser Skip (é do trinca-completa), foi %s (%s)", v, d)
	}
}

// Uma spec pode ser coberta por MAIS DE UMA feature; o requisito só precisa estar em alguma.
func TestSpecFeatureMatchUneVariasFeatures(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.feature"), []byte("@AAAAX-B01 @nivel-unit\n  Cenário: a\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "b.feature"), []byte("@AAAAX-B02 @nivel-e2e\n  Cenário: b\n"), 0o644)
	g := &mapx.Graph{Edges: []mapx.Edge{
		{From: "x.spec.md", To: "a.feature", Type: "covered-by"},
		{From: "x.spec.md", To: "b.feature", Type: "covered-by"},
	}}
	spec := "# Spec\n\n### AAAAX-B01 — x\n### AAAAX-B02 — y\n"
	v, d := checkSpecFeatureMatch(spec, mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}, dir, g, nil)
	if v != Pass {
		t.Fatalf("requisitos cobertos por features DIFERENTES deveriam passar, foi %s (%s)", v, d)
	}
}

// Bug PREEXISTENTE, achado ao escrever este gate: o parser de cenários só reconhecia
// `Cenário|Scenario` com um ` Outline` opcional. Em português o Scenario Outline é
// `Esquema do Cenário` — uma expressão própria, não "Cenário + sufixo". Resultado: 108
// cenários de um projeto real eram INVISÍVEIS para o `feature-test-match`, que é
// BLOQUEANTE e reportava verde sobre o que não enxergava.
func TestParseFeatureReconheceEsquemaEmQualquerIdioma(t *testing.T) {
	casos := map[string]string{
		"pt — Cenário":               "@AAAAX-B01 @nivel-unit\n  Cenário: faz algo\n",
		"pt — Esquema do Cenário":    "@AAAAX-B01 @nivel-unit\n  Esquema do Cenário: faz algo\n",
		"en — Scenario":              "@AAAAX-B01 @nivel-unit\n  Scenario: does something\n",
		"en — Scenario Outline":      "@AAAAX-B01 @nivel-unit\n  Scenario Outline: does something\n",
		"es — Escenario":             "@AAAAX-B01 @nivel-unit\n  Escenario: hace algo\n",
		"es — Esquema del escenario": "@AAAAX-B01 @nivel-unit\n  Esquema del escenario: hace algo\n",
		"de — Szenariogrundriss":     "@AAAAX-B01 @nivel-unit\n  Szenariogrundriss: macht etwas\n",
		"fr — Plan du scénario":      "@AAAAX-B01 @nivel-unit\n  Plan du scénario: fait quelque chose\n",
	}
	for nome, feature := range casos {
		t.Run(nome, func(t *testing.T) {
			got := parseFeatureScenarios(feature)
			if len(got) != 1 {
				t.Fatalf("cenário invisível para o parser — o gate reportaria verde sobre ele: %+v", got)
			}
			if got[0].Code != "AAAAX-B01" {
				t.Errorf("código errado: %+v", got[0])
			}
		})
	}
}

// A forma de Esquema precisa ser tentada ANTES da de Cenário: em vários idiomas ela
// CONTÉM a palavra "cenário" (`Esquema do Cenário`), e a alternativa curta casaria o
// prefixo e engoliria o título errado.
func TestParseFeatureEsquemaNaoEhEngolidoPeloPrefixo(t *testing.T) {
	got := parseFeatureScenarios("@AAAAX-B01 @nivel-unit\n  Esquema do Cenário: variante define a cor\n")
	if len(got) != 1 {
		t.Fatalf("não reconheceu: %+v", got)
	}
	if got[0].Title != "variante define a cor" {
		t.Fatalf("título capturado errado (%q) — a alternativa curta casou o prefixo", got[0].Title)
	}
}

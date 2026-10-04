// @anchors
//   code: FTTSF
//   ref: FTMFT

package gate

import (
	"github.com/co2-lab/anchors/internal/testlist"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cfg de teste com o de-para de regime do app de referência: @nivel-unit→unit, @nivel-integration→
// integration (superfície test); @nivel-e2e→e2e (outra superfície).
func regimeCfg() *config.Config {
	return &config.Config{
		Derived: &config.Derived{
			Regimes:  map[string]string{"nivel-unit": "unit", "nivel-integration": "integration", "nivel-e2e": "e2e"},
			Surfaces: map[string]string{"unit": "test", "integration": "test", "e2e": "e2e"},
		},
		// The fixtures are TypeScript tests; the family says how a test opens.
		Dialect: &config.Dialect{Family: "ts"},
	}
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const featureSrc = `# language: pt
# @anchors
` + "# " + `  ref: DDTDX
@backend @business-logic @dedup
Funcionalidade: Dedup

  @DDTDX-B01 @nivel-unit
  Cenário: Duplicata automática quando descrição e valor idênticos
    Dado tx idêntica
    Então classifica como duplicata

  @DDTDX-B02 @nivel-unit
  Cenário: Repetição real quando valor distinto
    Dado tx com valor distinto
    Então classifica como repeticao
`

func featureGraph(featPath, testPath string) *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: featPath, Kind: mapx.KindFeature},
			{ID: testPath, Kind: mapx.KindTest},
		},
		Edges: []mapx.Edge{
			{From: featPath, To: testPath, Type: mapx.EdgeTestedBy},
		},
	}
}

func TestFeatureTestMatch_pass(t *testing.T) {
	t.Run("FTMFT-B08: Tests implementing scenario codes with exact titles pass", func(t *testing.T) {})
	t.Run("FTMFT-X01: Static analysis does not run tests or inspect execution results", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	// teste implementa AMBOS os códigos com a descrição EXATA do cenário — a regra
	// não é "parecido o bastante": a divergência de uma palavra ("com" no lugar de
	// "quando") é o que o gate existe para acusar.
	writeFile(t, root, test, `
describe('dedup', () => {
  it('DDTDX-B01: Duplicata automática quando descrição e valor idênticos', () => {})
  it('DDTDX-B02: Repetição real quando valor distinto', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Pass {
		t.Errorf("esperava Pass, veio %v: %s", v, detail)
	}
}

func TestFeatureTestMatch_missingCode(t *testing.T) {
	t.Run("FTMFT-B06: A scenario code completely absent from tests fails", func(t *testing.T) {})
	t.Run("FTMFT-I01: Missing scenario code is always a failure", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	// implementa só B01 — B02 foi PULADO
	writeFile(t, root, test, `
describe('dedup', () => {
  it('DDTDX-B01: duplicata automática com descrição e valor idênticos', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Fail {
		t.Errorf("esperava Fail (B02 ausente), veio %v: %s", v, detail)
	}
}

func TestFeatureTestMatch_codeInCommentDoesNotCount(t *testing.T) {
	t.Run("FTMFT-B07: A scenario code appearing only in comments fails", func(t *testing.T) {})
	t.Run("FTMFT-I03: Code presence ignores comments while description matching reads them", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	// B02 aparece só em COMENTÁRIO — não conta como implementação
	writeFile(t, root, test, `
describe('dedup', () => {
  it('DDTDX-B01: duplicata automática com descrição e valor idênticos', () => {})
  // TODO: DDTDX-B02 ainda não implementado
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, _ := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Fail {
		t.Errorf("esperava Fail (B02 só em comentário), veio %v", v)
	}
}

func TestFeatureTestMatch_e2eAndVRSkipped(t *testing.T) {
	t.Run("FTMFT-B05: Scenarios belonging to non-test surfaces are skipped", func(t *testing.T) {})
	t.Run("FTMFT-X02: Non-unit surfaces are left to their respective gates", func(t *testing.T) {})
	root := t.TempDir()
	feat := "screens/Login.feature"
	test := "screens/Login.test.tsx"
	// LGN-S01 é @nivel-unit (cobrado); LGN-R01 é só @nivel-e2e (Maestro); LGN-VR é visual.
	featSrc := `# language: pt
# @anchors
` + "# " + `  ref: LGN0X
@screen @login
Funcionalidade: Login

  @LGN0X-S01 @nivel-unit
  Cenário: Render inicial mostra o formulário
    Então vejo o formulário

  @LGN0X-R01 @nivel-e2e
  Cenário: Login completo navega para Home
    Então chego na Home

  @LGN0X-VR @nivel-e2e
  Cenário: Regressão visual da tela
    Então a tela bate o baseline
`
	writeFile(t, root, feat, featSrc)
	// o teste só implementa o cenário unit; e2e e VR NÃO devem ser cobrados
	writeFile(t, root, test, `
describe('Login', () => {
  it('LGN0X-S01: render inicial mostra o formulário', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featSrc, n, root, g, regimeCfg())
	if v != Pass {
		t.Errorf("esperava Pass (e2e/VR não cobrados no .test), veio %v: %s", v, detail)
	}
}

func TestFeatureTestMatch_descriptionDrift(t *testing.T) {
	t.Run("FTMFT-B09: Tests with matching codes but drifting descriptions issue a warning", func(t *testing.T) {})
	t.Run("FTMFT-I02: Descriptive divergence is always an informative warning", func(t *testing.T) {})
	t.Run("FTMFT-X03: Minor description drift is a divergence, not a failure", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	// códigos presentes, mas descrições totalmente diferentes do cenário
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: xyz qwe abc', () => {})
  it('DDTDX-B02: foo bar baz', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Diverge {
		t.Errorf("esperava Pending (drift de descrição), veio %v: %s", v, detail)
	}
}

// O título entre aspas SIMPLES que cita algo entre aspas DUPLAS tem de ser lido
// inteiro. Antes o corpo da captura excluía as três quotes de uma vez, o título
// era truncado no primeiro `"`, e o gate acusava divergência de descrição num par
// que dizia exatamente a mesma coisa.
func TestFeatureTestMatch_tituloComAspasInternas(t *testing.T) {
	t.Run("FTMFT-B10: Test titles containing quotes are parsed without truncation", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: Duplicata automática quando descrição e valor idênticos', () => {})
  it('DDTDX-B02: Repetição real quando valor "distinto"', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}

	// O B02 do cenário não tem as aspas, então o par segue divergente — mas o que
	// importa aqui é COMO: o título lido tem de ser a frase inteira.
	if got, ok := testTitleFor(`it('DDTDX-B02: Repetição real quando valor "distinto"', () => {})`, "DDTDX-B02"); !ok ||
		got != `Repetição real quando valor "distinto"` {
		t.Fatalf("título lido = %q (ok=%v), queria a frase inteira com as aspas internas", got, ok)
	}

	// E com o cenário idêntico ao título, o gate passa.
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: Duplicata automática quando descrição e valor idênticos', () => {})
  it('DDTDX-B02: Repetição real quando valor distinto', () => {})
})`)
	if v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg()); v != Pass {
		t.Errorf("esperava Pass, veio %v: %s", v, detail)
	}
}

// Um teste que prova VÁRIOS cenários cita todos no título. O título é o texto
// depois de todos os códigos, e vale para QUALQUER um deles — o primeiro da lista
// e os seguintes.
func TestFeatureTestMatch_tituloComCodigosIrmaos(t *testing.T) {
	t.Run("FTMFT-B11: Sibling scenario codes in composite test titles are extracted", func(t *testing.T) {})
	corpo := `it('DDTDX-B01 / DDTDX-B02: duplicata e repetição saem do mesmo confronto', () => {})`

	for _, cod := range []string{"DDTDX-B01", "DDTDX-B02"} {
		got, ok := testTitleFor(corpo, cod)
		if !ok {
			t.Fatalf("%s: título não encontrado no título composto", cod)
		}
		if got != "duplicata e repetição saem do mesmo confronto" {
			t.Errorf("%s: título = %q, queria o texto DEPOIS de todos os códigos", cod, got)
		}
	}

	// Também na forma com colchetes e vírgula.
	if got, ok := testTitleFor(`it('[DDTDX-B01], [DDTDX-B02]: texto', () => {})`, "DDTDX-B02"); !ok ||
		got != "texto" {
		t.Errorf("forma com colchetes: título = %q (ok=%v)", got, ok)
	}

	// O irmão pode ser NOMINAL (`ABCDX-DS-<nome>`), não só numérico. Deixá-lo de
	// fora fazia o título do PRIMEIRO código vir com o prefixo do irmão grudado —
	// e o par aparecia como "similar 100%": mesmo texto, comparação diferente.
	nominal := `it('[DDTDX-B01] [DDTDX-DS-fatura-marcado] texto do caso', () => {})`
	for _, cod := range []string{"DDTDX-B01", "DDTDX-DS-fatura-marcado"} {
		if got, ok := testTitleFor(nominal, cod); !ok || got != "texto do caso" {
			t.Errorf("irmão nominal (%s): título = %q (ok=%v)", cod, got, ok)
		}
	}
}

// Título COMPARTILHADO por vários cenários não é confrontado por igualdade.
//
// Com N códigos num título, comparar cada cenário com o MESMO texto condenaria
// N-1 deles por construção — no máximo um pode ser idêntico. A pergunta certa ali
// é a da régua de corpo: o miolo do cenário está no teste?
func TestFeatureTestMatch_tituloCompartilhadoNaoExigeIgualdade(t *testing.T) {
	t.Run("FTMFT-B12: Shared test titles verify scenario presence through test body", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	// Um teste só, citando os dois códigos, com o miolo dos DOIS cenários no corpo.
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01 / DDTDX-B02: duplicata e repetição saem do mesmo confronto', () => {
    const idênticos = { descrição: 'x', valor: 10 }
    expect(classifica(idênticos)).toBe('duplicata')
    const distinto = { descrição: 'x', valor: 11 }
    expect(classifica(distinto)).toBe('repeticao')
  })
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	if v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg()); v != Pass {
		t.Errorf("esperava Pass (título conjunto, miolo dos dois no corpo), veio %v: %s", v, detail)
	}

	// E o compartilhamento é detectado como tal.
	corpo := `it('DDTDX-B01 / DDTDX-B02: texto', () => {})`
	if !sharedTitle(corpo, "DDTDX-B01") {
		t.Error("DDTDX-B01: título com dois códigos deveria contar como compartilhado")
	}
	if sharedTitle(`it('DDTDX-B01: texto', () => {})`, "DDTDX-B01") {
		t.Error("título com um só código NÃO é compartilhado")
	}
}

// Duas leituras do teste, duas perguntas. O CÓDIGO em comentário não conta como
// implementação; a DESCRIÇÃO em comentário conta como cobertura.
//
// É o comentário que liga o vocabulário do cenário ("duplicata automática") ao do
// código (`classifica`, `'duplicata'`). Sem ele nesta régua, o gate cobraria do
// TypeScript uma palavra portuguesa que ele nunca vai conter.
func TestFeatureTestMatch_comentarioCobreDescricaoMasNaoImplementa(t *testing.T) {
	t.Run("FTMFT-B13: Test comments contribute to descriptive match but not code presence", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)

	// O código do B02 aparece SÓ em comentário → segue faltando implementação.
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: Duplicata automática quando descrição e valor idênticos', () => {
    expect(classifica(a, b)).toBe('duplicata')
  })
  // DDTDX-B02: Repetição real quando valor distinto
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Fail {
		t.Errorf("código só em comentário deveria FALTAR: veio %v (%s)", v, detail)
	}

	// Agora o B02 tem `it` próprio, e a descrição de UM teste conjunto vem do
	// comentário: o gate aceita, porque a pergunta ali é de cobertura.
	writeFile(t, root, test, `
describe('x', () => {
  // Duplicata automática quando descrição e valor idênticos;
  // Repetição real quando valor distinto.
  it('DDTDX-B01 / DDTDX-B02: os dois vereditos saem do mesmo confronto', () => {
    expect(classifica(a, b)).toBe('duplicata')
    expect(classifica(a, c)).toBe('repeticao')
  })
})`)
	if v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg()); v != Pass {
		t.Errorf("descrição no comentário deveria COBRIR: veio %v (%s)", v, detail)
	}
}

// O código do cenário tem de TERMINAR onde termina: sem a fronteira,
// `ABCDX-DS-delta-up` casava dentro de `ABCDX-DS-delta-up-high`, e o gate comparava
// o cenário de "até 20%" com a prova de "acima de 20%" — divergência inventada.
func TestFeatureTestMatch_codigoNaoCasaPrefixoDeOutro(t *testing.T) {
	t.Run("FTMFT-B14: Scenario codes match with exact word boundaries", func(t *testing.T) {})
	casos := []struct{ corpo, cod, quer string }{
		{`it('SNBDX-DS-delta-up: até 20 em âmbar', () => {`, "SNBDX-DS-delta-up", "até 20 em âmbar"},
		{`it('SNBDX-DS-delta-up-high: acima de 20 em vermelho', () => {`, "SNBDX-DS-delta-up", ""},
		{`it('[DDTDX-B01] texto', () => {`, "DDTDX-B01", "texto"},
		{`it('DDTDX-B01: texto', () => {`, "DDTDX-B01", "texto"},
		{`it('DDTDX-B01 / DDTDX-B02: texto', () => {`, "DDTDX-B01", "texto"},
	}
	for _, c := range casos {
		got, ok := testTitleFor(c.corpo, c.cod)
		if c.quer == "" {
			if ok {
				t.Errorf("%s NÃO devia casar em %q, veio %q", c.cod, c.corpo, got)
			}
			continue
		}
		if !ok || got != c.quer {
			t.Errorf("%s: título = %q (ok=%v), queria %q", c.cod, got, ok, c.quer)
		}
	}
}

func TestFeatureTestMatch_tRunSuportado(t *testing.T) {
	t.Run("FTMFT-B15: Go t.Run declarations are recognized as valid test titles", func(t *testing.T) {})
	corpo := `t.Run("DDTDX-B01: duplicata encontrada", func(t *testing.T) {})`
	got, ok := testTitleFor(corpo, "DDTDX-B01")
	if !ok || got != "duplicata encontrada" {
		t.Errorf("t.Run em Go devia casar título: got=%q (ok=%v)", got, ok)
	}
}

func TestStripLineCommentsSuportaHash(t *testing.T) {
	t.Run("FTMFT-B16: Script comment markers are stripped when verifying code presence", func(t *testing.T) {})
	src := "val = 1 # comentario\n# linha inteira de comentario\nval2 = 2\n"
	res := stripLineComments(src)
	if strings.Contains(res, "comentario") {
		t.Errorf("comentários # deveriam ser removidos: %q", res)
	}
	if !strings.Contains(res, "val = 1") || !strings.Contains(res, "val2 = 2") {
		t.Errorf("código deve ser preservado: %q", res)
	}
}

func TestFeatureTestMatch_nonFeatureSkips(t *testing.T) {
	t.Run("FTMFT-B01: Non-feature artifacts skip confrontation", func(t *testing.T) {})
	n := mapx.Node{ID: "foo.spec.md", Kind: mapx.KindSpec}
	v, _ := checkFeatureTestMatch("", n, "", nil, nil)
	if v != Skip {
		t.Fatalf("esperava Skip para kind != feature, veio %v", v)
	}
}

func TestFeatureTestMatch_nilGraphPending(t *testing.T) {
	t.Run("FTMFT-B02: A nil graph returns pending without approving", func(t *testing.T) {})
	n := mapx.Node{ID: "foo.feature", Kind: mapx.KindFeature}
	v, _ := checkFeatureTestMatch(featureSrc, n, "", nil, nil)
	if v != Pending {
		t.Fatalf("esperava Pending para grafo nil, veio %v", v)
	}
}

func TestFeatureTestMatch_noScenariosSkips(t *testing.T) {
	t.Run("FTMFT-B03: A feature declaring no scenarios skips confrontation", func(t *testing.T) {})
	n := mapx.Node{ID: "empty.feature", Kind: mapx.KindFeature}
	g := &mapx.Graph{}
	v, _ := checkFeatureTestMatch("# Apenas comentários\n", n, "", g, nil)
	if v != Skip {
		t.Fatalf("esperava Skip para feature sem cenários, veio %v", v)
	}
}

func TestFeatureTestMatch_noLinkedTestsPending(t *testing.T) {
	t.Run("FTMFT-B04: A feature with no linked tests returns pending", func(t *testing.T) {})
	n := mapx.Node{ID: "feature_sem_teste.feature", Kind: mapx.KindFeature}
	g := &mapx.Graph{
		Nodes: []mapx.Node{n},
	}
	v, _ := checkFeatureTestMatch(featureSrc, n, "", g, nil)
	if v != Pending {
		t.Fatalf("esperava Pending para feature sem teste ligado, veio %v", v)
	}
}

func TestRootCode(t *testing.T) {
	t.Run("FTMFT-B17: RootCode returns the root requirement code without scenario sub-index", func(t *testing.T) {})
	if r := RootCode("ABCDX-B01#02"); r != "ABCDX-B01" {
		t.Fatalf("esperava ABCDX-B01, obteve %s", r)
	}
	if r := RootCode("ABCDX-B01"); r != "ABCDX-B01" {
		t.Fatalf("esperava ABCDX-B01, obteve %s", r)
	}
}

// An UNMAPPED regime tag used to pass for "another regime": `@nivel-compilacao` is not under
// the reference app's `regimes:`, and every scenario carrying it was skipped by this gate silently.
func TestFeatureTestMatch_unmappedRegimeTagIsStillConfronted(t *testing.T) {
	t.Run("FTMFT-B18: an unmapped regime tag does not exempt a scenario", func(t *testing.T) {})
	root := t.TempDir()
	feat := "screens/Form.feature"
	test := "screens/Form.test.tsx"
	featSrc := `# language: pt
# @anchors
` + "# " + `  ref: FRMXX
@screen
Funcionalidade: Form

  @FRMXX-B01 @nivel-unit
  Cenário: valida o campo
    Então vejo o erro

  @FRMXX-B02 @nivel-compilacao
  Cenário: o tipo recusa a chave extra
    Então o tsc recusa
`
	writeFile(t, root, feat, featSrc)
	writeFile(t, root, test, `
describe('Form', () => {
  it('FRMXX-B01: valida o campo', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featSrc, n, root, g, regimeCfg())
	if v == Pass || !strings.Contains(detail, "FRMXX-B02") {
		t.Errorf("the @nivel-compilacao scenario (unmapped) was not confronted: %v %s", v, detail)
	}
}

func TestFeatureTestMatch_Errors(t *testing.T) {
	t.Run("FTMFT-E01: A linked test gone from disk implements nothing while the others still count", func(t *testing.T) {
		root := t.TempDir()
		feat := "packages/backend/dedup.feature"
		present := "packages/backend/dedup.test.ts"
		writeFile(t, root, feat, featureSrc)
		writeFile(t, root, present, "it('DDTDX-B01: Duplicata automática quando descrição e valor idênticos', () => {})\n")
		g := featureGraph(feat, present)
		gone := "packages/backend/gone.test.ts"
		g.Nodes = append(g.Nodes, mapx.Node{ID: gone, Kind: mapx.KindTest})
		g.Edges = append(g.Edges, mapx.Edge{From: feat, To: gone, Type: mapx.EdgeTestedBy})
		v, d := checkFeatureTestMatch(featureSrc, mapx.Node{ID: feat, Kind: mapx.KindFeature}, root, g, regimeCfg())
		if v != Fail || !strings.Contains(d, "DDTDX-B02") || strings.Contains(d, "DDTDX-B01") {
			t.Fatalf("only the scenario the missing test would prove is charged; got %v: %s", v, d)
		}
	})
}

// A marker inside a string, or glued to a name, is not a comment.
func TestStripLineCommentsRespectsQuotes(t *testing.T) {
	t.Run("FTMFT-B19: Comment markers inside strings or glued to a name keep the line", func(t *testing.T) {})
	src := "run(\"gh\", \"--add-label\", initx.LabelManual)\n" +
		"get(\"https://x\", client.Fetch)\n" +
		"for i := n; i-- > 0; { step.Apply() }\n" +
		"val := 1 -- a SQL-style comment\n" +
		"x := Real() // a real comment\n"
	res := stripLineComments(src)
	for _, keep := range []string{"LabelManual", "client.Fetch", "step.Apply", "val := 1", "Real()"} {
		if !strings.Contains(res, keep) {
			t.Errorf("%q was cut away: %q", keep, res)
		}
	}
	for _, gone := range []string{"a SQL-style comment", "a real comment"} {
		if strings.Contains(res, gone) {
			t.Errorf("the trailing comment %q must be removed: %q", gone, res)
		}
	}
}

// The similarity verdict in the drift message is translated. It was the raw
// `similarity.Verdict.String`, in Portuguese ("divergente", "limítrofe") whatever the
// project's language.
func TestFeatureTestMatch_verdictIsTranslated(t *testing.T) {
	t.Run("FTMFT-B20: The drift verdict is written in the project's language", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: xyz qwe abc', () => {})
  it('DDTDX-B02: foo bar baz', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	for _, tc := range []struct{ lang, want, not string }{
		{"en", "(divergent,", "divergente"},
		{"pt-BR", "(divergente,", "(divergent,"},
	} {
		i18n.Set(tc.lang)
		_, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
		if !strings.Contains(detail, tc.want) || strings.Contains(detail, tc.not) {
			t.Errorf("[%s] detail %q: want the verdict %q, not %q", tc.lang, detail, tc.want, tc.not)
		}
	}
	i18n.Set(i18n.Default)
}

func TestTestTitleForReadsModifiedCalls(t *testing.T) {
	t.Run("FTMFT-B21: A parametrised, focused or skipped test is read by its own title", func(t *testing.T) {})
	next := "\n  it('RMIHX-E11 nota ausente é aceita como vazia', () => {})\n"
	for _, decl := range []string{
		"it.each([[\"a\", 1], [fn(2), {x: [3]}]])",
		"test.each([1, 2])",
		"it.each(\n    cases,\n  )",
		"it.only",
		"test.skip",
	} {
		body := decl + "('RMIHX-E11 relato fora da forma (%p) responde 400', (x) => {})" + next
		got, ok := testTitleFor(body, "RMIHX-E11")
		if !ok || got != "relato fora da forma (%p) responde 400" {
			t.Errorf("%s: want the title of that test, got %q (%v)", decl, got, ok)
		}
	}
}

// testTitleFor reads `body` as one test file through the ts and go families' patterns and
// returns the title of the test of `code`, as the gate does through the project's source.
func testTitleFor(body, code string) (string, bool) {
	title, _, ok := titleFor(listBody(body), code)
	return title, ok
}

// sharedTitle says whether the test of `code` in `body` cites other codes too.
func sharedTitle(body, code string) bool {
	_, shared, ok := titleFor(listBody(body), code)
	return ok && shared
}

func listBody(body string) []testlist.Test {
	dir, err := os.MkdirTemp("", "titles-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "x.test.ts"), []byte(body), 0o644); err != nil {
		panic(err)
	}
	ts := (&config.Config{Dialect: &config.Dialect{Family: "ts"}}).DialectFor().Tests.Pattern
	golang := (&config.Config{Dialect: &config.Dialect{Family: "go"}}).DialectFor().Tests.Pattern
	tests, err := testlist.List(dir, []string{"x.test.ts"}, testlist.Source{Pattern: ts + "|" + golang})
	if err != nil {
		panic(err)
	}
	return tests
}

func TestFeatureTestMatch_TitlesComeFromTheSource(t *testing.T) {
	t.Run("FTMFT-B21: A parametrised, focused or skipped test is read by its own title", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	// The file itself cites the codes with titles that diverge from the feature; the
	// project's script lists other titles. The script is what the gate must read.
	writeFile(t, root, test, "it('DDTDX-B01: xyz qwe abc', f)\nit('DDTDX-B02: foo bar baz', f)\n")
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	cfg := regimeCfg()
	cfg.Dialect = &config.Dialect{Tests: &config.TestsSource{}}
	titles := featureTitles(t, featureSrc)
	out := `{"version":1,"tests":[{"file":"` + test + `","line":1,"title":"DDTDX-B01: ` + titles[0] + `"},` +
		`{"file":"` + test + `","line":2,"title":"DDTDX-B02: ` + titles[1] + `"}]}`
	cfg.Dialect.Tests.Script = "printf '%s' '" + out + "'"
	if v, detail := checkFeatureTestMatch(featureSrc, n, root, g, cfg); v != Pass {
		t.Fatalf("the script's titles match the feature and must be the ones read, got %v (%s)", v, detail)
	}
	resetProjectTestsCache()
	if v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg()); v == Pass || !strings.Contains(detail, "DDTDX-B01") {
		t.Fatalf("through the ts pattern the file's own diverging titles are read, got %v (%s)", v, detail)
	}
}

func TestFeatureTestMatch_FailingSource(t *testing.T) {
	t.Run("FTMFT-E02: A failing tests source fails the gate naming the error", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	writeFile(t, root, test, "it('DDTDX-B01: a', f)\nit('DDTDX-B02: b', f)\n")
	cfg := regimeCfg()
	cfg.Dialect = &config.Dialect{Tests: &config.TestsSource{Script: `echo '{"version":9,"tests":[]}'`}}
	v, detail := checkFeatureTestMatch(featureSrc, mapx.Node{ID: feat, Kind: mapx.KindFeature}, root, featureGraph(feat, test), cfg)
	if v != Fail || !strings.Contains(detail, "`version` must be 1") {
		t.Fatalf("a source outside the contract must fail the gate naming the violation, got %v (%s)", v, detail)
	}
}

// featureTitles lists the scenario titles of a feature, in order.
func featureTitles(t *testing.T, src string) []string {
	t.Helper()
	var out []string
	for _, sc := range parseFeatureScenarios(src) {
		out = append(out, sc.Title)
	}
	if len(out) < 2 {
		t.Fatalf("the fixture needs two scenarios, has %v", out)
	}
	return out
}

func TestFeatureTestMatch_SupportIsNotItsTest(t *testing.T) {
	t.Run("FTMFT-B22: A support file linked to a feature is not confronted as its test", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	writeFile(t, root, test, "it('DDTDX-B01: x', f)\nit('DDTDX-B02: y', f)\n")
	g := featureGraph(feat, test)
	for i := range g.Nodes {
		if g.Nodes[i].ID == test {
			g.Nodes[i].Support = true
		}
	}
	v, detail := checkFeatureTestMatch(featureSrc, mapx.Node{ID: feat, Kind: mapx.KindFeature}, root, g, regimeCfg())
	if v != Pending || !strings.Contains(detail, i18n.T("gate.feature_test_match.pending_no_linked_tests")) {
		t.Fatalf("with only a support file linked there is no test to confront, got %v (%s)", v, detail)
	}
}

// A rule legitimately has more than one scenario — the happy path and the alternatives.
// Without a suffix, the N scenarios carry the same code and nothing tells them apart: the
// gate compares the N titles with the same test and at most one matches.
func TestFeatureScenarios_suffixGivesEachScenarioItsIdentity(t *testing.T) {
	t.Run("FTMFT-B23: A scenario suffix gives each case of a rule its own identity", func(t *testing.T) {})
	feature := `
  @USBPX-B01#01 @USBPX-B02#01 @nivel-unit
  Cenário: busca pontos por userId+month
    Então o repository é consultado

  @USBPX-B01#02 @nivel-unit
  Cenário: sem usuário logado nada é buscado
    Então o repository não é consultado
`
	scenarios := parseFeatureScenarios(feature)
	if len(scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(scenarios))
	}
	if scenarios[0].Code != "USBPX-B01#01" {
		t.Errorf("code of the 1st: %q, want USBPX-B01#01", scenarios[0].Code)
	}
	if scenarios[1].Code != "USBPX-B01#02" {
		t.Errorf("code of the 2nd: %q, want USBPX-B01#02", scenarios[1].Code)
	}
	// The second code of the tag line carries its suffix too.
	if len(scenarios[0].Codes) != 2 || scenarios[0].Codes[1] != "USBPX-B02#01" {
		t.Errorf("codes of the 1st: %v", scenarios[0].Codes)
	}
}

// Backward compatible: the suffix is optional. A project that never adopts it keeps
// working — and that is the condition for the gate to change without breaking anyone.
func TestFeatureScenarios_codeWithoutSuffixStillHolds(t *testing.T) {
	t.Run("FTMFT-B23: A scenario suffix gives each case of a rule its own identity", func(t *testing.T) {})
	scenarios := parseFeatureScenarios(`
  @SAUTX-B01 @nivel-unit
  Cenário: Hidratar carrega a sessão
    Então o usuário fica disponível
`)
	if len(scenarios) != 1 || scenarios[0].Code != "SAUTX-B01" {
		t.Fatalf("a code without a suffix broke: %+v", scenarios)
	}
}

// The gates that speak of a RULE need the root; the ones that speak of a SCENARIO, the
// whole code. Confusing the two would make the rule `USBPX-B01` look like three rules.
func TestRootCode_separatesRuleFromScenario(t *testing.T) {
	t.Run("FTMFT-B17: RootCode returns the root requirement code without scenario sub-index", func(t *testing.T) {})
	cases := map[string]string{
		"USBPX-B01#02": "USBPX-B01",
		"USBPX-B01":    "USBPX-B01",
		"ATLNX-VR":     "ATLNX-VR",
		"MNMTX-DS-kv":  "MNMTX-DS-kv",
	}
	for input, want := range cases {
		if got := RootCode(input); got != want {
			t.Errorf("RootCode(%q) = %q, want %q", input, got, want)
		}
	}
}

// The case from the reference app: a logo easter egg carried the code of "open the alerts",
// behind a test that did open them — and only the first test was ever compared.
func TestFeatureTestMatch_everyCitingTestIsCompared(t *testing.T) {
	t.Run("FTMFT-B24: Every test that cites the code is compared, not only the first", func(t *testing.T) {})
	root := t.TempDir()
	feat := "business-logic/dedup.feature"
	test := "__tests__/dedup.test.ts"
	writeFile(t, root, feat, featureSrc)
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: Duplicata automática quando descrição e valor idênticos', () => {})
  it('DDTDX-B01: duplicata automática quando descrição e valor idênticos em outra conta', () => {})
  it('DDTDX-B01: o logo gira ao tocar três vezes no cabeçalho', () => {})
  it('DDTDX-B02: Repetição real quando valor distinto', () => {})
})`)
	g := featureGraph(feat, test)
	n := mapx.Node{ID: feat, Kind: mapx.KindFeature}
	v, detail := checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Diverge || !strings.Contains(detail, "o logo gira") || strings.Contains(detail, "outra conta") {
		t.Fatalf("only the unrelated test is named, as a warning; got %v: %s", v, detail)
	}

	// The first test keeps its own ruler (any drift), and is not named again as "another";
	// a title shared by several codes describes none of them alone.
	writeFile(t, root, test, `
describe('x', () => {
  it('DDTDX-B01: o logo gira ao tocar três vezes no cabeçalho', () => {})
  it('DDTDX-B01 / DDTDX-B02: tabela de casos variados sem relação', () => {})
  it('DDTDX-B02: Repetição real quando valor distinto', () => {})
})`)
	v, detail = checkFeatureTestMatch(featureSrc, n, root, g, regimeCfg())
	if v != Diverge || strings.Count(detail, "o logo gira") != 0 || strings.Contains(detail, "tabela de casos") || !strings.Contains(detail, "DDTDX-B01 (") {
		t.Fatalf("the first test drifts under its own ruler only, the shared title is not named; got %v: %s", v, detail)
	}
}

func TestFeatureTestMatch_titleAndMore(t *testing.T) {
	t.Run("FTMFT-B25: A test title that says the scenario's title and more matches it", func(t *testing.T) {})
	for test, want := range map[string]bool{
		"O login trava após três tentativas — sandbox": true,
		"O login trava após três tentativas (backend)": true,
		"O login trava após três tentativas: com %s":   true,
		"O login trava após $count tentativas":         false,
		"Após três tentativas o login trava":           false,
		"O login":                                      false,
	} {
		if got := titleCovers("O login trava após três tentativas", test); got != want {
			t.Errorf("%q covers the scenario: got %v, want %v", test, got, want)
		}
	}
	if !titleCovers("Cada <formato> é aceito", "Cada %s é aceito") {
		t.Error("an outline's parameter and a table placeholder are not words of either title")
	}
}

func TestFeatureCodes_aVRCodeCarriesItsState(t *testing.T) {
	t.Run("FTMFT-B26: A VR code carries its state, and each state's scenario is its own code", func(t *testing.T) {})
	feature := "Feature: Button\n\n" +
		"  @state @BUTTN-VR-S01 @vr-level\n  Scenario: The button in Enabled looks like its baseline\n    Given x\n\n" +
		"  @state @BUTTN-VR-S02 @vr-level\n  Scenario: The button in Disabled looks like its baseline\n    Given x\n"
	var codes []string
	for _, s := range parseFeatureScenarios(feature) {
		codes = append(codes, s.Codes...)
	}
	if strings.Join(codes, ",") != "BUTTN-VR-S01,BUTTN-VR-S02" {
		t.Errorf("codes = %v", codes)
	}
	if v, d := checkScenarioIdentity(feature, mapx.Node{ID: "ui/Button.feature", Kind: mapx.KindFeature}, "", nil, nil); v == Fail {
		t.Errorf("two states' VR scenarios are two codes: %v %s", v, d)
	}
}

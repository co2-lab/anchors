package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cfg de teste com o de-para de regime do app de referência: @nivel-unit→unit, @nivel-integration→
// integration (superfície test); @nivel-e2e→e2e (outra superfície).
func regimeCfg() *config.Config {
	return &config.Config{Derived: &config.Derived{
		Regimes:  map[string]string{"nivel-unit": "unit", "nivel-integration": "integration", "nivel-e2e": "e2e"},
		Surfaces: map[string]string{"unit": "test", "integration": "test", "e2e": "e2e"},
	}}
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
#   ref: DDTDX
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
	root := t.TempDir()
	feat := "screens/Login.feature"
	test := "screens/Login.test.tsx"
	// LGN-S01 é @nivel-unit (cobrado); LGN-R01 é só @nivel-e2e (Maestro); LGN-VR é visual.
	featSrc := `# language: pt
# @anchors
#   ref: LGN0X
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
	if v != Pending {
		t.Errorf("esperava Pending (drift de descrição), veio %v: %s", v, detail)
	}
}

// O título entre aspas SIMPLES que cita algo entre aspas DUPLAS tem de ser lido
// inteiro. Antes o corpo da captura excluía as três quotes de uma vez, o título
// era truncado no primeiro `"`, e o gate acusava divergência de descrição num par
// que dizia exatamente a mesma coisa.
func TestFeatureTestMatch_tituloComAspasInternas(t *testing.T) {
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
	if !tituloCompartilhado(corpo, "DDTDX-B01") {
		t.Error("DDTDX-B01: título com dois códigos deveria contar como compartilhado")
	}
	if tituloCompartilhado(`it('DDTDX-B01: texto', () => {})`, "DDTDX-B01") {
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

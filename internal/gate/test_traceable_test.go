package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaRastreavel(t *testing.T, teste, feature string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.feature"), []byte(feature), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "u.feature", Kind: mapx.KindFeature}, {ID: "u.test.ts", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{{From: "u.feature", To: "u.test.ts", Type: mapx.EdgeTestedBy}},
	}
	return checkTestTraceable(teste, mapx.Node{ID: "u.test.ts", Kind: mapx.KindTest}, root, g, nil)
}

func TestTestTraceable_NaoTestePula(t *testing.T) {
	t.Run("TSTRT-B01: Non-test artifacts skip confrontation", func(t *testing.T) {})
	t.Run("TSTRT-I01: Traceability is strictly charged on test artifacts", func(t *testing.T) {})
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindFeature, mapx.KindCode} {
		n := mapx.Node{ID: "arquivo", Kind: k}
		v, msg := checkTestTraceable("qualquer conteudo", n, "", &mapx.Graph{}, nil)
		if v != Skip {
			t.Fatalf("esperava Skip para kind %v, veio %v (%s)", k, v, msg)
		}
	}
}

func TestTestTraceable_GrafoNilPendente(t *testing.T) {
	t.Run("TSTRT-B02: Confrontation without a dependency graph returns Pending", func(t *testing.T) {})
	t.Run("TSTRT-I04: Missing graph structure produces Pending rather than approving", func(t *testing.T) {})
	n := mapx.Node{ID: "u.test.ts", Kind: mapx.KindTest}
	v, _ := checkTestTraceable("conteudo", n, "", nil, nil)
	if v != Pending {
		t.Fatalf("sem grafo deve retornar Pending, veio %v", v)
	}
}

func TestTestTraceable_SemFeatureLigadaPula(t *testing.T) {
	t.Run("TSTRT-B03: A test with no linked feature skips confrontation", func(t *testing.T) {})
	t.Run("TSTRT-I02: Tests without a linked feature are never charged", func(t *testing.T) {})
	t.Run("TSTRT-X03: Scenario codes are not required in standalone test files", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.test.ts", Kind: mapx.KindTest}}}
	v, _ := checkTestTraceable("it('helper sem feature', () => {})",
		mapx.Node{ID: "u.test.ts", Kind: mapx.KindTest}, t.TempDir(), g, nil)
	if v != Skip {
		t.Fatalf("teste sem feature ligada deve pular: %v", v)
	}
}

func TestTestTraceable_FeatureInexistentePula(t *testing.T) {
	t.Run("TSTRT-B04: A test whose linked feature cannot be read skips confrontation", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "sumiu.feature", Kind: mapx.KindFeature}, {ID: "u.test.ts", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{{From: "sumiu.feature", To: "u.test.ts", Type: mapx.EdgeTestedBy}},
	}
	v, _ := checkTestTraceable("conteudo", mapx.Node{ID: "u.test.ts", Kind: mapx.KindTest}, root, g, nil)
	if v != Skip {
		t.Fatalf("feature que nao existe no disco deve retornar Skip, veio %v", v)
	}
}

func TestTestTraceable_FeatureSemCodigosPula(t *testing.T) {
	t.Run("TSTRT-B05: A test whose linked feature declares no scenario codes skips confrontation", func(t *testing.T) {})
	v, _ := rodaRastreavel(t, "it('x', () => {})", "Feature: Sem codigos\nScenario: sem tag\n")
	if v != Skip {
		t.Fatalf("feature sem codigos de cenario deve pular, veio %v", v)
	}
}

func TestTestTraceable_CodigoDeclaradoPassa(t *testing.T) {
	t.Run("TSTRT-B06: A test containing a declared scenario code passes", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	codeA01 := "PHA1X" + "-A01"
	feat := "@" + codeS01 + " @unit-level\nScenario: a\n@" + codeA01 + " @unit-level\nScenario: b\n"
	teste := "it('" + codeS01 + ": renderiza', () => { expect(1).toBe(1) })"
	v, msg := rodaRastreavel(t, teste, feat)
	if v != Pass {
		t.Fatalf("teste contendo codigo valido deve passar: %v (%s)", v, msg)
	}
}

func TestTestTraceable_UmCodigoBasta(t *testing.T) {
	t.Run("TSTRT-B07: A single declared scenario code is sufficient to pass", func(t *testing.T) {})
	t.Run("TSTRT-I03: The acceptance threshold requires only one scenario code present", func(t *testing.T) {})
	t.Run("TSTRT-X01: The gate does not enforce one-to-one scenario coverage", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	codeA01 := "PHA1X" + "-A01"
	codeB02 := "PHA1X" + "-B02"
	feat := "@" + codeS01 + "\nScenario: a\n@" + codeA01 + "\nScenario: b\n@" + codeB02 + "\nScenario: c\n"
	teste := "it('" + codeB02 + ": apenas este cenario', () => {})"
	v, msg := rodaRastreavel(t, teste, feat)
	if v != Pass {
		t.Fatalf("um unico codigo ja torna o teste rastreavel: %v (%s)", v, msg)
	}
}

func TestTestTraceable_SemCodigoReprova(t *testing.T) {
	t.Run("TSTRT-B08: A test containing none of the linked feature scenario codes fails", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	codeA01 := "PHA1X" + "-A01"
	feat := "@" + codeS01 + "\nScenario: a\n@" + codeA01 + "\nScenario: b\n"
	teste := "it('adiciona com id gerado sem citar codigo', () => {})"
	v, _ := rodaRastreavel(t, teste, feat)
	if v != Fail {
		t.Fatalf("teste sem codigo deve falhar: %v", v)
	}
}

func TestTestTraceable_CodigoTrocadoReprova(t *testing.T) {
	t.Run("TSTRT-B09: A test with transposed or misspelled codes fails exact substring confrontation", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	feat := "@" + codeS01 + "\nScenario: a\n"
	teste := "it('PH1A: codigo com letra e digito trocados', () => {})"
	v, _ := rodaRastreavel(t, teste, feat)
	if v != Fail {
		t.Fatalf("codigo trocado nao deve passar na busca por substring: %v", v)
	}
}

func TestTestTraceable_MensagemCitaFeatureECodigos(t *testing.T) {
	t.Run("TSTRT-B10: A failing verdict names the linked feature and lists expected codes", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	codeA01 := "PHA1X" + "-A01"
	codeB01 := "PHA1X" + "-B01"
	codeC01 := "PHA1X" + "-C01"
	feat := "@" + codeS01 + "\nScenario: a\n@" + codeA01 + "\nScenario: b\n@" + codeB01 + "\nScenario: c\n@" + codeC01 + "\nScenario: d\n"
	v, msg := rodaRastreavel(t, "it('sem codigo', () => {})", feat)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, "u.feature") {
		t.Errorf("mensagem deve citar o arquivo da feature: %s", msg)
	}
	if !strings.Contains(msg, codeS01) {
		t.Errorf("mensagem deve citar o primeiro codigo como sugestao: %s", msg)
	}
	if !strings.Contains(msg, codeA01) || !strings.Contains(msg, codeB01) {
		t.Errorf("mensagem deve listar codigos esperados: %s", msg)
	}
}

func TestTestTraceable_CodigoEmQualquerLugarPassa(t *testing.T) {
	t.Run("TSTRT-B11: Scenario codes appearing anywhere in test content satisfy traceability", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	feat := "@" + codeS01 + "\nScenario: a\n"
	testeComComentario := "// Cobertura do requisito " + codeS01 + "\ndescribe('modulo', () => {})"
	v, msg := rodaRastreavel(t, testeComComentario, feat)
	if v != Pass {
		t.Fatalf("codigo em comentario deve ser aceito: %v (%s)", v, msg)
	}
}

func TestTestTraceable_NaoExecutaTestes(t *testing.T) {
	t.Run("TSTRT-X02: Test execution results are not inspected by this gate", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	feat := "@" + codeS01 + "\nScenario: a\n"
	testeComAssertFalso := "it('" + codeS01 + ": assert falho', () => { throw new Error('boom'); })"
	v, msg := rodaRastreavel(t, testeComAssertFalso, feat)
	if v != Pass {
		t.Fatalf("rastreabilidade nao executa o teste, deve passar estaticamente: %v (%s)", v, msg)
	}
}

func TestTestTraceable_NaoAvaliaSemanticaDoTeste(t *testing.T) {
	t.Run("TSTRT-X04: Test assertion semantics and quality are not evaluated", func(t *testing.T) {})
	codeS01 := "PHA1X" + "-S01"
	feat := "@" + codeS01 + "\nScenario: a\n"
	testeVazio := "test('" + codeS01 + ": vazio', () => {});"
	v, msg := rodaRastreavel(t, testeVazio, feat)
	if v != Pass {
		t.Fatalf("validacao de semantica de assercao nao pertence a este gate: %v (%s)", v, msg)
	}
}

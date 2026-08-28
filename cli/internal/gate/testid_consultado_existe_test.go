package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// fixtureConsultado monta um projeto mínimo com as duas pontas que o gate confronta:
// o CÓDIGO (que expõe) e a superfície E2E (que consulta).
func fixtureConsultado(t *testing.T, codigo, flow string) (string, *config.Config) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "flows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tela.tsx"), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "flows", "F-A01.yaml"), []byte(flow), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{Derived: &config.Derived{
		TestHandle: "testID",
		Surfaces:   map[string]string{"e2e": "e2e"},
		Files:      map[string]string{"e2e": "flows"},
	}}
	return root, cfg
}

func TestConsultadoExiste_handleQueExistePassa(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<View testID=":abcd-tela" />`,
		"- assertVisible:\n    id: ':abcd-tela'\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("handle exposto pelo código deveria passar, veio %v: %s", v, d)
	}
}

// O caso que motivou o gate: o flow procura um id que ninguém expõe.
func TestConsultadoExiste_handleInventadoReprova(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<View testID=":abcd-tela" />`,
		"- tapOn:\n    id: ':abcd-botao-que-ninguem-expoe'\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Fail {
		t.Fatalf("handle inexistente deveria reprovar, veio %v: %s", v, d)
	}
	if !strings.Contains(d, "abcd-botao-que-ninguem-expoe") {
		t.Fatalf("o laudo deve NOMEAR o handle ausente; veio: %s", d)
	}
	if !strings.Contains(d, "F-A01.yaml") {
		t.Fatalf("o laudo deve nomear o flow que consulta, senão o achado não é acionável; veio: %s", d)
	}
}

// O caso mais grave, e o que o runner nunca dá: num assertNotVisible o id inexistente
// faz o teste PASSAR. É verde por vacuidade.
func TestConsultadoExiste_assertNotVisibleVacuoReprova(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<View testID=":abcd-tela" />`,
		"- assertNotVisible:\n    id: ':abcd-controles-de-edicao'\n")
	v, _ := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Fail {
		t.Fatalf("assertNotVisible sobre id inexistente é verde-por-vacuidade e deve reprovar, veio %v", v)
	}
}

// Template no código (`:item-${id}`) cobre a instância que o flow procura (`:item-3`).
// Sem isto todo id com sufixo dinâmico viraria achado.
func TestConsultadoExiste_templateCobreInstancia(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		"<View testID={`:abcd-item-${id}`} />",
		"- tapOn:\n    id: ':abcd-item-3'\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("instância de template exposto deveria passar, veio %v: %s", v, d)
	}
}

// O inverso: o flow procura por PADRÃO (regex do Maestro) e o código expõe a instância.
func TestConsultadoExiste_regexDoFlowCasaCabeca(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		"<View testID={`:abcd-linha-${i}`} />",
		"- tapOn:\n    id: ':abcd-linha-.*maestro-001'\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("regex do flow deve casar a cabeça exposta, veio %v: %s", v, d)
	}
}

// Id COMPOSTO em runtime pelo próprio flow: o valor final não existe até a execução, e
// acusá-lo seria inventar defeito.
func TestConsultadoExiste_interpolacaoDoFlowNaoAcusa(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<View testID=":abcd-tela" />`,
		"- tapOn:\n    id: \"${':abcd-x-' + output.data.ns}\"\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("id interpolado pelo runner não é confrontável e não deve acusar, veio %v: %s", v, d)
	}
}

// Handle guardado em TABELA de consulta (não colado ao atributo). O reconhecedor
// estrito não o vê; aqui ele tem de contar, senão o gate acusa quem cumpre.
func TestConsultadoExiste_handleEmTabelaConta(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		"const rotas = { Family: ':abcd-navigate-to-family' }",
		"- tapOn:\n    id: ':abcd-navigate-to-family'\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("handle em tabela de consulta existe no código e não deve acusar, veio %v: %s", v, d)
	}
}

// Template atrás de fallback (`?? `) — a forma do AlertSheet do app de referência.
func TestConsultadoExiste_templateAtrasDeFallbackConta(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		"<Btn testID={btn.testID ?? `:abcd-alert-sheet-button-${i}`} />",
		"- tapOn:\n    id: ':abcd-alert-sheet-button-1'\n")
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("template atrás de `??` existe no código e não deve acusar, veio %v: %s", v, d)
	}
}

// Sufixo composto num componente FILHO, a partir de prop: o id final
// (`:ctdt-list-row-0`) não existe em arquivo nenhum — a cabeça vem de quem renderiza
// (`testID=":ctdt-list"`) e o sufixo nasce no filho. Sem reconhecer isto o gate acusava
// 8 flows corretos do app de referência de uma vez.
func TestConsultadoExiste_sufixoCompostoNoFilhoConta(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<List testID=":abcd-list" />`,
		"- tapOn:\n    id: ':abcd-list-row-0'\n")
	// O componente filho compõe o sufixo a partir da prop recebida.
	if err := os.WriteFile(filepath.Join(root, "List.tsx"),
		[]byte("<Row testID={testID ? `${testID}-row-${index}` : undefined} />"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("sufixo composto no filho existe em runtime e não deve acusar, veio %v: %s", v, d)
	}
}

// Sufixo TERMINAL composto da prop de mesmo nome: `${testID}-toggle` no atom Input
// gera `:reis-register-input-password-toggle` — que não CONTÉM "-toggle-", só termina
// em "-toggle". O gate acusava esse handle de inexistente, e eu cheguei a "consertar"
// um flow correto por causa disso (o toque no olho virou toque em ponto neutro, que
// apagou a senha já digitada).
func TestConsultadoExiste_sufixoTerminalConta(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<Input testID=":abcd-input-password" />`,
		"- tapOn:\n    id: ':abcd-input-password-toggle'\n")
	if err := os.WriteFile(filepath.Join(root, "Input.tsx"),
		[]byte("<Btn testID={testID ? `${testID}-toggle` : undefined} />"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, d := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("sufixo terminal composto existe em runtime e não deve acusar, veio %v: %s", v, d)
	}
}

// Arquivo de TESTE não expõe handle — ele consulta. Foi assim que um id fantasma
// sobreviveu no app de referência, referenciado só por um teste de unidade.
func TestConsultadoExiste_testeNaoContaComoExposicao(t *testing.T) {
	root, cfg := fixtureConsultado(t,
		`<View testID=":abcd-tela" />`,
		"- tapOn:\n    id: ':abcd-fantasma'\n")
	if err := os.WriteFile(filepath.Join(root, "tela.test.tsx"),
		[]byte(`getByTestId(':abcd-fantasma')`), 0o644); err != nil {
		t.Fatal(err)
	}
	v, _ := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Fail {
		t.Fatalf("id presente só em arquivo de teste não está exposto e deve reprovar, veio %v", v)
	}
}

// Sem `test_handle` declarado o gate PULA — inferir `testID` faria ele reportar verde
// sobre o que não conferiu.
func TestConsultadoExiste_semHandleDeclaradoPula(t *testing.T) {
	root, cfg := fixtureConsultado(t, `<View testID=":abcd-tela" />`, "- tapOn:\n    id: ':abcd-x'\n")
	cfg.Derived.TestHandle = ""
	v, _ := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Skip {
		t.Fatalf("sem atributo declarado o gate deve pular, veio %v", v)
	}
}

// Sem superfície e2e declarada não há o que varrer: Skip, nunca Pass — "não medi" não
// é "está limpo".
func TestConsultadoExiste_semSuperficieE2EPula(t *testing.T) {
	root, cfg := fixtureConsultado(t, `<View testID=":abcd-tela" />`, "- tapOn:\n    id: ':abcd-x'\n")
	cfg.Derived.Files = map[string]string{}
	v, _ := checkTestIDConsultadoExiste("", mapx.Node{}, root, nil, cfg)
	if v != Skip {
		t.Fatalf("sem superfície declarada o gate deve pular, veio %v", v)
	}
}

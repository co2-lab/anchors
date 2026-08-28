package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const moduloBase = `import x from 'y'

export function useAvailableMonths() {
  return 1
}

export function useMonthlySummary(
  month: string,
  userId?: string,
) {
  return { month, userId }
}
`

// O detector é do PROJETO: aqui, o dialeto jest/vitest. Sem ele o gate pula.
const detectJS = `(?:jest|vi)\.mock\(['"]([^'"]+)`

func cfgComCarimbo() *config.Config {
	return &config.Config{Derived: &config.Derived{MockDetect: detectJS}}
}

// escreveModulo grava o módulo e devolve a raiz temporária.
func escreveModulo(t *testing.T, corpo string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "src/mod.ts"), []byte(corpo), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// carimboDe calcula o valor CERTO para o trecho — é o que o gate vai recalcular.
func carimboDe(t *testing.T, corpo, ancora string, qtd int) string {
	t.Helper()
	linhas := strings.Split(corpo, "\n")
	for i, l := range linhas {
		if l == ancora {
			fim := i + qtd
			if fim > len(linhas) {
				fim = len(linhas)
			}
			return hashDoTrecho(strings.Join(linhas[i:fim], "\n"))
		}
	}
	t.Fatalf("âncora não achada no fixture: %q", ancora)
	return ""
}

func rodaCarimbo(t *testing.T, root, teste string) (Verdict, string) {
	t.Helper()
	n := mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}
	return checkMockCarimbado(teste, n, root, &mapx.Graph{}, cfgComCarimbo())
}

const ancora = "export function useMonthlySummary("

func TestMockCarimbado_carimboQueBatePassa(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	h := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	if v, msg := rodaCarimbo(t, root, teste); v != Pass {
		t.Errorf("carimbo que corresponde ao módulo deve passar: %v (%s)", v, msg)
	}
}

// O caso que o gate existe para pegar: o trecho mudou e o dublê ficou para trás.
func TestMockCarimbado_trechoMudouReprova(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	// carimbo tirado de uma versão ANTIGA (com um parâmetro a menos)
	antigo := strings.Replace(moduloBase, "  userId?: string,\n", "", 1)
	h := carimboDe(t, antigo, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("trecho alterado deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "contrato hoje") {
		t.Errorf("a mensagem deve mostrar o valor atual: %s", msg)
	}
}

// A razão de a âncora ser a LINHA e não o número: editar acima não pode invalidar o
// carimbo de quem não mudou. Medido no app de referência — com número de linha, um comentário na
// linha 5 quebrava o carimbo da função da linha 70.
func TestMockCarimbado_imuneADeslocamento(t *testing.T) {
	h := carimboDe(t, moduloBase, ancora, 5)
	deslocado := "// comentário novo no topo\n// e outro\n" + moduloBase
	root := escreveModulo(t, deslocado)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	if v, msg := rodaCarimbo(t, root, teste); v != Pass {
		t.Errorf("deslocar o trecho não muda o contrato dele: %v (%s)", v, msg)
	}
}

// Âncora que some (renome/remoção) é ACHADO, não erro de ferramenta: o dublê
// certamente está desatualizado, e falhar explícito é melhor que silêncio.
func TestMockCarimbado_ancoraSumidaReprova(t *testing.T) {
	renomeado := strings.Replace(moduloBase, ancora, "export function useMonthlySummaryV2(", 1)
	root := escreveModulo(t, renomeado)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | deadbeef\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("âncora ausente deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "renomeado ou removido") {
		t.Errorf("a mensagem deve explicar o que houve: %s", msg)
	}
}

// Âncora repetida torna o alvo ambíguo. O gate prefere ACUSAR a escolher uma: um
// carimbo que aponta para "alguma das duas" não prova nada.
func TestMockCarimbado_ancoraAmbiguaReprova(t *testing.T) {
	root := escreveModulo(t, moduloBase+"\n"+ancora+"\n  x: number,\n)\n")
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | deadbeef\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("âncora ambígua deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "ambígua") {
		t.Errorf("a mensagem deve dizer que é ambígua: %s", msg)
	}
}

// `qtd` é a janela, e ela é FIXA no carimbo — o alcance fica à vista de quem lê, e o
// gate não precisa adivinhar onde o bloco termina (o que exigiria parser por linguagem).
func TestMockCarimbado_qtdDelimitaAJanela(t *testing.T) {
	// muda a ÚLTIMA linha do corpo, fora de uma janela de 2 linhas
	mudado := strings.Replace(moduloBase, "  return { month, userId }", "  return { month }", 1)
	root := escreveModulo(t, mudado)

	h2 := carimboDe(t, moduloBase, ancora, 2)
	teste2 := "// @contract: src/mod.ts | " + ancora + " | 2 | " + h2 + "\njest.mock('src/mod')"
	if v, _ := rodaCarimbo(t, root, teste2); v != Pass {
		t.Errorf("janela de 2 linhas não alcança a mudança: %v", v)
	}

	h9 := carimboDe(t, moduloBase, ancora, 9)
	teste9 := "// @contract: src/mod.ts | " + ancora + " | 9 | " + h9 + "\njest.mock('src/mod')"
	if v, _ := rodaCarimbo(t, root, teste9); v != Fail {
		t.Errorf("janela de 9 linhas deve alcançar a mudança: %v", v)
	}
}

// Sem `derived.mock_stamp` o gate pula: adotar o carimbo é decisão do projeto.
func TestMockCarimbado_semDeclaracaoPula(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	n := mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | deadbeef"

	v, msg := checkMockCarimbado(teste, n, root, &mapx.Graph{}, &config.Config{})
	if v != Skip {
		t.Errorf("sem declaração não há o que confrontar: %v", v)
	}
	if !strings.Contains(msg, "mock_detect") {
		t.Errorf("a mensagem deve dizer o que declarar: %s", msg)
	}
}

// Sem módulo REGIDO dublado não há carimbo a cobrar — aqui o grafo está vazio, então
// `src/mod` conta como externo. (A ausência de carimbo em módulo regido é Fail; ver
// TestMockCarimbado_ausenciaDeCarimboReprova.)
func TestMockCarimbado_semModuloRegidoPula(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	if v, _ := rodaCarimbo(t, root, "jest.mock('src/mod', () => ({}))"); v != Skip {
		t.Errorf("dublê de módulo não-regido não exige carimbo: %v", v)
	}
}

// Módulo que não existe mais é achado com mensagem própria — não um crash.
func TestMockCarimbado_moduloInexistenteReprova(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	teste := "// @contract: src/sumiu.ts | " + ancora + " | 5 | deadbeef"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("módulo ausente deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "não encontrado") {
		t.Errorf("a mensagem deve nomear o problema: %s", msg)
	}
}

// grafoComMod — o módulo dublado é REGIDO pelo projeto, condição para a cobrança.
func grafoComMod() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{{ID: "src/mod.ts", Kind: mapx.KindCode}}}
}

// A metade mais importante do gate: AUSÊNCIA de carimbo é acusada, não pulada.
//
// Carimbo divergente ACUSA; carimbo ausente é SILÊNCIO — o mesmo "falha aberto" que o
// `trinca-completa` existe para fechar. Se a ausência passasse, o carimbo viraria
// opcional na prática e o mecanismo protegeria só quem já escolheu ser protegido.
func TestMockCarimbado_ausenciaDeCarimboReprova(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	teste := "jest.mock('src/mod', () => ({ useMonthlySummary: jest.fn() }))"

	v, msg := checkMockCarimbado(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo())
	if v != Fail {
		t.Fatalf("dublê sem carimbo deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "sem carimbo") {
		t.Errorf("a mensagem deve nomear a ausência: %s", msg)
	}
}

// Terceiro segue de fora também aqui — o recorte é o mesmo do `mock-tipado`, senão o
// gate acusaria todo dublê de biblioteca e viraria ruído.
func TestMockCarimbado_terceiroSemCarimboNaoEhCobrado(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	teste := "jest.mock('@gorhom/bottom-sheet', () => ({ BottomSheet: 'View' }))"

	if v, msg := checkMockCarimbado(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo()); v != Skip {
		t.Errorf("dublê de fora do projeto não exige carimbo: %v (%s)", v, msg)
	}
}

// Carimbo presente e correto satisfaz a cobrança de ausência E a de correspondência.
func TestMockCarimbado_carimboPresenteSatisfazAmbas(t *testing.T) {
	root := escreveModulo(t, moduloBase)
	h := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h +
		"\njest.mock('src/mod', () => ({ useMonthlySummary: jest.fn() }))"

	if v, msg := checkMockCarimbado(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo()); v != Pass {
		t.Errorf("dublê carimbado e correspondente deve passar: %v (%s)", v, msg)
	}
}

// Regex inválido é erro de CONFIGURAÇÃO e falha ALTO. Silenciá-lo faria o gate varrer
// zero dublês e reportar verde — o pior desfecho possível num medidor.
func TestMockCarimbado_detectorInvalidoReprova(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{MockDetect: `jest\.mock\(([`}}
	v, msg := checkMockCarimbado("jest.mock('src/mod')",
		mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}, t.TempDir(), grafoComMod(), cfg)
	if v != Fail {
		t.Fatalf("padrão que não compila deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "não compila") {
		t.Errorf("a mensagem deve dizer que o padrão é inválido: %s", msg)
	}
}

// Sem grupo de captura o gate não sabe QUAL módulo foi dublado — é config incompleta,
// não arquivo defeituoso.
func TestMockCarimbado_detectorSemCapturaReprova(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{MockDetect: `jest\.mock`}}
	v, msg := checkMockCarimbado("jest.mock('src/mod')",
		mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}, t.TempDir(), grafoComMod(), cfg)
	if v != Fail {
		t.Fatalf("padrão sem captura deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "grupo de captura") {
		t.Errorf("a mensagem deve dizer o que falta: %s", msg)
	}
}

// O detector é do PROJETO: um dialeto diferente (Python) é reconhecido igual, desde que
// declarado. É o que torna o carimbo de fato agnóstico.
func TestMockCarimbado_detectorDeOutroDialeto(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{
		MockDetect: `(?:mock\.)?patch\(['"]([^'"]+)`,
	}}
	root := escreveModulo(t, moduloBase)
	v, msg := checkMockCarimbado("@patch('src/mod')\ndef test_x(): pass",
		mapx.Node{ID: "x_test.py", Kind: mapx.KindTest}, root, grafoComMod(), cfg)
	if v != Fail {
		t.Fatalf("dublê Python sem carimbo deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "sem carimbo") {
		t.Errorf("mesma cobrança, outro dialeto: %s", msg)
	}
}

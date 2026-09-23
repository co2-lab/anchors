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
			return snippetHash(strings.Join(linhas[i:fim], "\n"))
		}
	}
	t.Fatalf("âncora não achada no fixture: %q", ancora)
	return ""
}

func rodaCarimbo(t *testing.T, root, teste string) (Verdict, string) {
	t.Helper()
	n := mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}
	return checkMockStamped(teste, n, root, &mapx.Graph{}, cfgComCarimbo())
}

const ancora = "export function useMonthlySummary("

func TestMockCarimbado_carimboQueBatePassa(t *testing.T) {
	t.Run("MCSTM-B02: A stamp that matches the module today passes", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	h := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	if v, msg := rodaCarimbo(t, root, teste); v != Pass {
		t.Errorf("carimbo que corresponde ao módulo deve passar: %v (%s)", v, msg)
	}
}

// O caso que o gate existe para pegar: o trecho mudou e o dublê ficou para trás.
func TestMockCarimbado_trechoMudouReprova(t *testing.T) {
	t.Run("MCSTM-B03: A snippet that changed since the stamp was written fails", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	// carimbo tirado de uma versão ANTIGA (com um parâmetro a menos)
	antigo := strings.Replace(moduloBase, "  userId?: string,\n", "", 1)
	h := carimboDe(t, antigo, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("trecho alterado deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "contrato hoje") && !strings.Contains(msg, "contract today") {
		t.Errorf("a mensagem deve mostrar o valor atual: %s", msg)
	}
}

// A razão de a âncora ser a LINHA e não o número: editar acima não pode invalidar o
// carimbo de quem não mudou. Medido no app de referência — com número de linha, um comentário na
// linha 5 quebrava o carimbo da função da linha 70.
func TestMockCarimbado_imuneADeslocamento(t *testing.T) {
	t.Run("MCSTM-B04: The stamp is immune to displacement", func(t *testing.T) {})
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
	t.Run("MCSTM-B05: An anchor that vanished fails with its own message", func(t *testing.T) {})
	renomeado := strings.Replace(moduloBase, ancora, "export function useMonthlySummaryV2(", 1)
	root := escreveModulo(t, renomeado)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | deadbeef\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("âncora ausente deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "renomeado ou removido") && !strings.Contains(msg, "renamed or removed") {
		t.Errorf("a mensagem deve explicar o que houve: %s", msg)
	}
}

// Âncora repetida torna o alvo ambíguo. O gate prefere ACUSAR a escolher uma: um
// carimbo que aponta para "alguma das duas" não prova nada.
func TestMockCarimbado_ancoraAmbiguaReprova(t *testing.T) {
	t.Run("MCSTM-B06: An anchor occurring more than once is ambiguous and fails", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase+"\n"+ancora+"\n  x: number,\n)\n")
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | deadbeef\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("âncora ambígua deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "ambígua") && !strings.Contains(msg, "ambiguous") {
		t.Errorf("a mensagem deve dizer que é ambígua: %s", msg)
	}
}

// `qtd` é a janela, e ela é FIXA no carimbo — o alcance fica à vista de quem lê, e o
// gate não precisa adivinhar onde o bloco termina (o que exigiria parser por linguagem).
func TestMockCarimbado_qtdDelimitaAJanela(t *testing.T) {
	t.Run("MCSTM-B07: The declared line count delimits the window", func(t *testing.T) {})
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
	t.Run("MCSTM-B08: Without the dialect declared the gate goes quiet", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	n := mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | deadbeef"

	v, msg := checkMockStamped(teste, n, root, &mapx.Graph{}, &config.Config{})
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
	t.Run("MCSTM-B10: A stamp whose module no longer exists is a finding, not a crash", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	teste := "// @contract: src/sumiu.ts | " + ancora + " | 5 | deadbeef"

	v, msg := rodaCarimbo(t, root, teste)
	if v != Fail {
		t.Fatalf("módulo ausente deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "não encontrado") && !strings.Contains(msg, "not found") {
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
	t.Run("MCSTM-B11: The absence of a stamp on a governed double is accused", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	teste := "jest.mock('src/mod', () => ({ useMonthlySummary: jest.fn() }))"

	v, msg := checkMockStamped(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo())
	if v != Fail {
		t.Fatalf("dublê sem carimbo deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "sem carimbo") && !strings.Contains(msg, "without contract stamp") {
		t.Errorf("a mensagem deve nomear a ausência: %s", msg)
	}
}

// Terceiro segue de fora também aqui — o recorte é o mesmo do `mock-tipado`, senão o
// gate acusaria todo dublê de biblioteca e viraria ruído.
func TestMockCarimbado_terceiroSemCarimboNaoEhCobrado(t *testing.T) {
	t.Run("MCSTM-B09: A double of a module the project does not govern is not charged", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	teste := "jest.mock('@gorhom/bottom-sheet', () => ({ BottomSheet: 'View' }))"

	if v, msg := checkMockStamped(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo()); v != Skip {
		t.Errorf("dublê de fora do projeto não exige carimbo: %v (%s)", v, msg)
	}
}

// Carimbo presente e correto satisfaz a cobrança de ausência E a de correspondência.
func TestMockCarimbado_carimboPresenteSatisfazAmbas(t *testing.T) {
	t.Run("MCSTM-B12: A stamp present and correct satisfies both charges", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	h := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h +
		"\njest.mock('src/mod', () => ({ useMonthlySummary: jest.fn() }))"

	if v, msg := checkMockStamped(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo()); v != Pass {
		t.Errorf("dublê carimbado e correspondente deve passar: %v (%s)", v, msg)
	}
}

// Regex inválido é erro de CONFIGURAÇÃO e falha ALTO. Silenciá-lo faria o gate varrer
// zero dublês e reportar verde — o pior desfecho possível num medidor.
func TestMockCarimbado_detectorInvalidoReprova(t *testing.T) {
	t.Run("MCSTM-B13: A dialect regex that does not compile fails loudly", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{MockDetect: `jest\.mock\(([`}}
	v, msg := checkMockStamped("jest.mock('src/mod')",
		mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}, t.TempDir(), grafoComMod(), cfg)
	if v != Fail {
		t.Fatalf("padrão que não compila deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "não compila") && !strings.Contains(msg, "does not compile") {
		t.Errorf("a mensagem deve dizer que o padrão é inválido: %s", msg)
	}
}

// Sem grupo de captura o gate não sabe QUAL módulo foi dublado — é config incompleta,
// não arquivo defeituoso.
func TestMockCarimbado_detectorSemCapturaReprova(t *testing.T) {
	t.Run("MCSTM-B14: A dialect regex with no capture group fails", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{MockDetect: `jest\.mock`}}
	v, msg := checkMockStamped("jest.mock('src/mod')",
		mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}, t.TempDir(), grafoComMod(), cfg)
	if v != Fail {
		t.Fatalf("padrão sem captura deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "grupo de captura") && !strings.Contains(msg, "capture group") {
		t.Errorf("a mensagem deve dizer o que falta: %s", msg)
	}
}

// O detector é do PROJETO: um dialeto diferente (Python) é reconhecido igual, desde que
// declarado. É o que torna o carimbo de fato agnóstico.
func TestMockCarimbado_detectorDeOutroDialeto(t *testing.T) {
	t.Run("MCSTM-B15: Another ecosystem's dialect is charged the same way", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{
		MockDetect: `(?:mock\.)?patch\(['"]([^'"]+)`,
	}}
	root := escreveModulo(t, moduloBase)
	v, msg := checkMockStamped("@patch('src/mod')\ndef test_x(): pass",
		mapx.Node{ID: "x_test.py", Kind: mapx.KindTest}, root, grafoComMod(), cfg)
	if v != Fail {
		t.Fatalf("dublê Python sem carimbo deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "sem carimbo") && !strings.Contains(msg, "without contract stamp") {
		t.Errorf("mesma cobrança, outro dialeto: %s", msg)
	}
}

// Nó que não é teste não é assunto deste gate: só o teste declara dublê.
func TestMockCarimbado_naoTesteNaoTemVeredito(t *testing.T) {
	t.Run("MCSTM-B01: An artifact that is not a test leaves without a verdict", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	teste := "jest.mock('src/mod', () => ({}))"
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindFeature} {
		v, _ := checkMockStamped(teste, mapx.Node{ID: "x.ts", Kind: k}, root, grafoComMod(), cfgComCarimbo())
		if v != Skip {
			t.Errorf("kind %v: quer Skip, obteve %v", k, v)
		}
	}
}

// O gate RECALCULA — não valida formato. O carimbo aqui está bem formado e estava
// certo quando foi escrito; o módulo mudou depois. Se o gate só conferisse formato,
// ele passaria, e o carimbo passaria a certificar a si mesmo.
func TestMockCarimbado_recalculaEmVezDeValidarFormato(t *testing.T) {
	t.Run("MCSTM-I01: The gate recomputes the hash instead of validating the stamp's format", func(t *testing.T) {})
	h := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	// Com o módulo original: passa — o carimbo é bem formado E corresponde.
	if v, msg := rodaCarimbo(t, escreveModulo(t, moduloBase), teste); v != Pass {
		t.Fatalf("carimbo correto devia passar: %v (%s)", v, msg)
	}

	// MESMO carimbo, mesmo formato, módulo alterado depois: reprova.
	mudado := strings.Replace(moduloBase, "return { month, userId }", "return { month, userId, extra: 1 }", 1)
	v, msg := rodaCarimbo(t, escreveModulo(t, mudado), teste)
	if v != Fail {
		t.Fatalf("o carimbo bem formado não basta — o gate recalcula: %v (%s)", v, msg)
	}
}

// A ligação dublê↔carimbo é pelo CAMINHO, casado por sufixo sem extensão: o dublê
// vem por alias (`@/mod`) e o carimbo por caminho de disco (`src/mod.ts`). São o
// mesmo arquivo, e exigir um resolvedor de alias seria específico do ecossistema.
func TestMockCarimbado_ligaDubleAoCarimboPeloCaminho(t *testing.T) {
	t.Run("MCSTM-I02: The double and the stamp are tied by the module's path", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	h := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h +
		"\njest.mock('@app/src/mod', () => ({ useMonthlySummary: jest.fn() }))"

	v, msg := checkMockStamped(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo())
	if v != Pass {
		t.Fatalf("alias e caminho de disco descrevem o mesmo arquivo: %v (%s)", v, msg)
	}
}

// Falha de CONFIGURAÇÃO reprova; decisão do PROJETO pula. A diferença é se alguém
// ESCOLHEU o silêncio — regex quebrado e regex sem captura são defeito, dialeto não
// declarado é escolha.
func TestMockCarimbado_configDefeituosaReprovaEscolhaPula(t *testing.T) {
	t.Run("MCSTM-I03: A configuration fault fails and a project decision skips", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	teste := "jest.mock('src/mod', () => ({}))"
	n := mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}

	casos := []struct {
		nome string
		cfg  *config.Config
		quer Verdict
	}{
		{"regex que não compila", &config.Config{Derived: &config.Derived{MockDetect: `jest\.mock\(([`}}, Fail},
		{"regex sem captura", &config.Config{Derived: &config.Derived{MockDetect: `jest\.mock`}}, Fail},
		{"dialeto não declarado", &config.Config{Derived: &config.Derived{}}, Skip},
	}
	for _, c := range casos {
		if v, msg := checkMockStamped(teste, n, root, grafoComMod(), c.cfg); v != c.quer {
			t.Errorf("%s: quer %v, obteve %v (%s)", c.nome, c.quer, v, msg)
		}
	}
}

// O carimbo é HASH DE TEXTO, e é por isso que ele pega o que um extrator de
// assinatura não pegaria: uma constante mudada DENTRO do corpo, sem tocar a
// assinatura. E pega sem precisar existir uma vez por linguagem.
func TestMockCarimbado_pegaMudancaSoNoCorpo(t *testing.T) {
	t.Run("MCSTM-X01: The gate does not interpret the code of the stamped module", func(t *testing.T) {})
	h := carimboDe(t, moduloBase, ancora, 5)
	// A assinatura fica IDÊNTICA; só o corpo muda.
	soCorpo := strings.Replace(moduloBase, "  return { month, userId }", "  return { month, userId, v: 2 }", 1)
	if !strings.Contains(soCorpo, ancora) {
		t.Fatal("a assinatura devia continuar idêntica no fixture")
	}
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + h + "\njest.mock('src/mod')"

	v, msg := rodaCarimbo(t, escreveModulo(t, soCorpo), teste)
	if v != Fail {
		t.Fatalf("mudança só no corpo devia ser acusada — hash de texto alcança: %v (%s)", v, msg)
	}
}

// Sem dialeto declarado o gate PULA, não passa. Num ecossistema onde o dublê não é
// uma chamada detectável (em Go é uma interface satisfeita, não há o que detectar),
// reportar verde sobre o que não se conferiu seria a pior falha de um medidor.
func TestMockCarimbado_semDialetoPulaEmVezDePassar(t *testing.T) {
	t.Run("MCSTM-X02: The gate carries no built-in dialect for detecting doubles", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	// Um dublê de módulo REGIDO, que sob o dialeto declarado seria Fail por ausência.
	teste := "jest.mock('src/mod', () => ({ useMonthlySummary: jest.fn() }))"
	n := mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest}

	v, _ := checkMockStamped(teste, n, root, grafoComMod(), &config.Config{})
	if v != Skip {
		t.Fatalf("sem dialeto declarado o gate se cala; verde seria mentira: %v", v)
	}
	// E com o dialeto declarado, o MESMO arquivo é cobrado — prova de que o Skip
	// vem da ausência de declaração, não da ausência de dublê.
	if v2, _ := checkMockStamped(teste, n, root, grafoComMod(), cfgComCarimbo()); v2 != Fail {
		t.Fatalf("com dialeto declarado o mesmo arquivo devia ser cobrado: %v", v2)
	}
}

// A ausência de carimbo NÃO é afrouxada para acomodar legado. Desenhar pelo projeto
// antigo transformaria um problema de migração em propriedade permanente do
// framework: todo projeto futuro herdaria a frouxidão.
func TestMockCarimbado_ausenciaNaoEhAfrouxadaPorLegado(t *testing.T) {
	t.Run("MCSTM-X03: The gate does not skip the absence of a stamp to accommodate legacy code", func(t *testing.T) {})
	root := escreveModulo(t, moduloBase)
	// Vários dublês sem carimbo de uma vez — a forma que o legado tem.
	teste := "jest.mock('src/mod', () => ({}))\njest.mock('src/mod', () => ({}))\n"
	v, msg := checkMockStamped(teste, mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest},
		root, grafoComMod(), cfgComCarimbo())
	if v != Fail {
		t.Fatalf("volume de achados não afrouxa a régua: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "src/mod") {
		t.Errorf("o achado devia nomear o módulo sem carimbo; msg = %q", msg)
	}
}

// O hash é TRUNCADO de propósito: ele mora numa linha de comentário e é lido por
// humano. A comparação é sobre o valor truncado, não sobre o digest inteiro.
func TestMockCarimbado_hashEhTruncadoParaLeituraHumana(t *testing.T) {
	t.Run("MCSTM-X04: The gate does not guarantee cryptographic strength", func(t *testing.T) {})
	h := snippetHash("qualquer trecho")
	if len(h) != 8 {
		t.Fatalf("o carimbo devia ser curto o bastante para caber num comentário: %q", h)
	}
	root := escreveModulo(t, moduloBase)
	certo := carimboDe(t, moduloBase, ancora, 5)
	teste := "// @contract: src/mod.ts | " + ancora + " | 5 | " + certo + "\njest.mock('src/mod')"
	if v, msg := rodaCarimbo(t, root, teste); v != Pass {
		t.Fatalf("a comparação é sobre o valor truncado: %v (%s)", v, msg)
	}
	// E um valor truncado diferente reprova — o truncamento não cega o gate.
	errado := "// @contract: src/mod.ts | " + ancora + " | 5 | 00000000\njest.mock('src/mod')"
	if v, _ := rodaCarimbo(t, root, errado); v != Fail {
		t.Fatalf("truncar não pode cegar a comparação: %v", v)
	}
}

// A MODULE WITH A DOT IN ITS NAME (`auth.store`). The gate used to strip an "extension"
// from the specifier and cut part of the name, so a stamp on `auth.store.ts` never matched
// the double of `@/src/stores/auth.store`. Measured after stamping a real project: 93
// tests still failed as "double without stamp", all like this.
func TestGenerateStamps_moduleWithADotInItsName(t *testing.T) {
	t.Run("MCSTM-B16: a module whose name has a dot is matched to its stamp", func(t *testing.T) {})
	test := "jest.mock('@/src/stores/auth.store')\n"
	root, g := stampProject(t, map[string]string{"src/stores/auth.store.ts": realHooks, "src/Home.test.tsx": test})
	out, written, _, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil || len(written) != 1 {
		t.Fatalf("expected one stamp: %v %v", written, err)
	}
	if v, msg := checkMockStamped(out, mapx.Node{ID: "src/Home.test.tsx", Kind: mapx.KindTest}, root, g, stampCfg()); v != Pass {
		t.Errorf("the stamp on auth.store.ts does not satisfy the double of auth.store: %v / %s", v, msg)
	}
}

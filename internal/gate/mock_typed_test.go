package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func testNode() mapx.Node { return mapx.Node{ID: "x.test.ts", Kind: mapx.KindTest} }

func cfgComContrato() *config.Config {
	return &config.Config{Derived: &config.Derived{MockContract: "Partial<typeof {{module}}>"}}
}

// grafoDoProjeto — os módulos que o projeto REGE. Só eles são cobrados; a biblioteca
// de terceiro não aparece no mapa e por isso fica de fora (ver `ehModuloRegido`).
func grafoDoProjeto() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{
		{ID: "packages/backend/repositories/lotes.ts", Kind: mapx.KindCode},
		{ID: "src/a.ts", Kind: mapx.KindCode},
		{ID: "src/b.ts", Kind: mapx.KindCode},
	}}
}

func rodaMock(t *testing.T, content string, cfg *config.Config) (Verdict, string) {
	t.Helper()
	return checkMockTyped(content, testNode(), "", grafoDoProjeto(), cfg)
}

// O caso que o gate existe para pegar: dublê SEM amarra continua verde depois que o
// módulo muda de assinatura. É prova falsa, não ausência de prova.
func TestMockTipado_dubleSemAmarraReprova(t *testing.T) {
	t.Run("MCTYM-B01: A double with no tie fails and the verdict names the loose module", func(t *testing.T) {})
	src := `jest.mock('@backend/repositories/lotes', () => ({
  listar: jest.fn(),
}))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Fail {
		t.Fatalf("dublê sem amarra deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "@backend/repositories/lotes") {
		t.Errorf("a mensagem deve NOMEAR o módulo solto: %s", msg)
	}
}

// A amarra é uma anotação de tipo na fábrica — é ela que faz o compilador conferir
// nome, assinatura e retorno contra o módulo real.
func TestMockTipado_dubleAmarradoPassa(t *testing.T) {
	t.Run("MCTYM-B02: A double whose factory carries the declared tie passes", func(t *testing.T) {})
	src := `import type * as Real from '@backend/repositories/lotes'
jest.mock('@backend/repositories/lotes', (): Partial<typeof Real> => ({
  listar: jest.fn(),
}))`
	if v, msg := rodaMock(t, src, cfgComContrato()); v != Pass {
		t.Errorf("dublê anotado deveria passar: %v (%s)", v, msg)
	}
}

// A cobrança é por MÓDULO, não por arquivo: tipar três e esquecer o quarto é
// exatamente o buraco a pegar, e um veredito por arquivo o daria como resolvido.
func TestMockTipado_cobraPorModuloNaoPorArquivo(t *testing.T) {
	t.Run("MCTYM-B03: The charge is per module, not per file", func(t *testing.T) {})
	src := `import type * as A from 'a'
jest.mock('a', (): Partial<typeof A> => ({ f: jest.fn() }))
jest.mock('b', () => ({ g: jest.fn() }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Fail {
		t.Fatalf("um dublê solto entre vários basta para reprovar: %v", v)
	}
	// A contagem é o que separa "um solto" de "todos soltos": nomear `b` sem dizer
	// quantos deixaria o leitor sem saber se `a` também está na lista.
	if !strings.Contains(msg, "1 dublê(s) sem amarra") && !strings.Contains(msg, "1 double(s) without binding") {
		t.Errorf("deve contar apenas o solto: %s", msg)
	}
	// O módulo amarrado NÃO entra na lista de soltos. (A verificação é sobre a lista,
	// não sobre a mensagem inteira: o texto do contrato declarado — que contém
	// `typeof` — aparece legitimamente na instrução de conserto ao final.)
	var lista string
	if parts := strings.SplitN(msg, "real: ", 2); len(parts) > 1 {
		lista = strings.SplitN(parts[1], ".", 2)[0]
	}
	if strings.Contains(lista, "a") {
		t.Errorf("o módulo amarrado não deve ser acusado: lista=%q", lista)
	}
}

// `jest.mock('x')` sem fábrica usa o automock, que deriva do módulo real por
// construção — não drifta, e cobrá-lo seria ruído que se aprende a ignorar.
func TestMockTipado_automockNaoEhCobrado(t *testing.T) {
	t.Run("MCTYM-B04: A double with no factory is not charged", func(t *testing.T) {})
	src := `jest.mock('@backend/repositories/lotes')`
	if v, msg := rodaMock(t, src, cfgComContrato()); v != Skip {
		t.Errorf("automock não tem fábrica que possa mentir: %v (%s)", v, msg)
	}
}

// Sem `derived.mock_contract` o gate PULA. Inferir `Partial<typeof …>` assumiria
// TypeScript e reportaria verde sobre o que não se conferiu em qualquer outro
// ecossistema — a pior falha possível num medidor.
func TestMockTipado_semContratoDeclaradoPula(t *testing.T) {
	t.Run("MCTYM-B05: Without the tie shape declared the gate goes quiet", func(t *testing.T) {})
	src := `jest.mock('a', () => ({ f: jest.fn() }))`
	v, msg := rodaMock(t, src, &config.Config{})
	if v != Skip {
		t.Errorf("sem contrato declarado o gate não tem o que cobrar: %v", v)
	}
	if !strings.Contains(msg, "mock_contract") {
		t.Errorf("a mensagem deve dizer o que declarar: %s", msg)
	}
}

// A cobrança é do TESTE — é lá que o dublê vive. Rodar sobre a spec acusaria o
// arquivo errado.
func TestMockTipado_soRodaSobreTeste(t *testing.T) {
	t.Run("MCTYM-B06: An artifact that is not a test leaves without a verdict", func(t *testing.T) {})
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkMockTyped("", n, "", grafoDoProjeto(), cfgComContrato()); v != Skip {
		t.Errorf("o dublê é cobrado do teste: %v", v)
	}
}

// Teste que não dubla ninguém não tem o que amarrar — Skip, não Pass: não houve
// verificação, e um Pass aqui inflaria a contagem de verdes com nada.
func TestMockTipado_semDubleNadaACobrar(t *testing.T) {
	t.Run("MCTYM-B07: A test that doubles nobody leaves without a verdict", func(t *testing.T) {})
	src := `it('soma', () => { expect(1+1).toBe(2) })`
	if v, _ := rodaMock(t, src, cfgComContrato()); v != Skip {
		t.Errorf("sem dublê não há amarra a cobrar: %v", v)
	}
}

// O dialeto do runner é do projeto: `vi.mock` (vitest) conta igual a `jest.mock`.
func TestMockTipado_reconheceVitest(t *testing.T) {
	t.Run("MCTYM-B08: Another runner of the same ecosystem is recognised the same way", func(t *testing.T) {})
	src := `vi.mock('a', () => ({ f: vi.fn() }))`
	if v, _ := rodaMock(t, src, cfgComContrato()); v != Fail {
		t.Errorf("`vi.mock` é dublê como qualquer outro: %v", v)
	}
}

// Biblioteca de TERCEIRO não é cobrada. O drift que o gate persegue é "o vizinho mudou
// e o dublê não soube", e o vizinho que muda toda semana é o módulo próprio: a
// dependência externa tem versão travada no lockfile, e o dublê dela troca um
// componente por um stub em vez de reproduzir um contrato.
//
// Medido: sem este recorte o gate acusou 305 arquivos num projeto real, 245 deles por
// dublê de terceiro. Gate que acusa tudo não é lido — é desligado, e leva junto os
// achados legítimos.
func TestMockTipado_terceiroNaoEhCobrado(t *testing.T) {
	t.Run("MCTYM-B09: A third-party library double is not charged", func(t *testing.T) {})
	src := `jest.mock('@gorhom/bottom-sheet', () => ({ BottomSheet: 'View' }))
jest.mock('@react-navigation/native', () => ({ useNavigation: jest.fn() }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Skip {
		t.Errorf("dublê de fora do projeto não é cobrado: %v (%s)", v, msg)
	}
}

// O recorte não pode virar escape: com um dublê PRÓPRIO solto no meio dos de terceiro,
// o gate continua reprovando — e nomeia só o próprio.
func TestMockTipado_terceiroNaoEsconde0Proprio(t *testing.T) {
	t.Run("MCTYM-B10: The third-party exemption is not an escape hatch", func(t *testing.T) {})
	src := `jest.mock('@gorhom/bottom-sheet', () => ({ BottomSheet: 'View' }))
jest.mock('@backend/repositories/lotes', () => ({ listar: jest.fn() }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Fail {
		t.Fatalf("o próprio solto deve reprovar mesmo cercado de terceiros: %v", v)
	}
	if strings.Contains(msg, "gorhom") {
		t.Errorf("terceiro não deve ser acusado: %s", msg)
	}
	if !strings.Contains(msg, "1 dublê(s) sem amarra") && !strings.Contains(msg, "1 double(s) without binding") {
		t.Errorf("deve contar só o próprio: %s", msg)
	}
}

// Import relativo (`../stores/auth.store`) é módulo próprio como qualquer outro — o
// critério é resolver para um nó do mapa, não a forma do especificador.
func TestMockTipado_relativoEhProprio(t *testing.T) {
	t.Run("MCTYM-B11: A relative import that resolves in the map is an own module", func(t *testing.T) {})
	src := `jest.mock('../a', () => ({ f: jest.fn() }))`
	if v, _ := rodaMock(t, src, cfgComContrato()); v != Fail {
		t.Errorf("import relativo que resolve no mapa é próprio: %v", v)
	}
}

func TestMockTipado_mockDetectGenerico(t *testing.T) {
	t.Run("MCTYM-B12: Another ecosystem's dialect is charged the same way", func(t *testing.T) {})
	cfg := &config.Config{
		Derived: &config.Derived{
			MockDetect:   `(?:mock\.)?patch\(['"]([^'"]+)`,
			MockContract: "autospec=True",
		},
	}
	srcSemAmarra := `@patch('src/a')
def test_foo():
    pass`
	v, msg := rodaMock(t, srcSemAmarra, cfg)
	if v != Fail {
		t.Fatalf("patch sem autospec=True devia reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "src/a") {
		t.Errorf("devia acusar src/a: %s", msg)
	}

	srcComAmarra := `@patch('src/a', autospec=True)
def test_foo():
    pass`
	v2, msg2 := rodaMock(t, srcComAmarra, cfg)
	if v2 != Pass {
		t.Fatalf("patch com autospec=True devia passar: %v (%s)", v2, msg2)
	}
}

// O que é REGIDO é decidido pelo GRAFO, não por lista de prefixos na config. As três
// formas de especificador — alias, relativo e nome nu — apontam para o MESMO nó do
// mapa, e as três são cobradas. Isso não pede configuração nova e acompanha o projeto
// sozinho: código que nasce entra no mapa e passa a ser cobrado.
func TestMockTipado_regidoEhDecididoPeloGrafo(t *testing.T) {
	t.Run("MCTYM-I01: What is governed is decided by the graph, never by a prefix list", func(t *testing.T) {})
	formas := []struct {
		nome string
		spec string
	}{
		{"alias", "@app/src/a"},
		{"relativo", "../src/a"},
		{"nome nu", "src/a"},
	}
	for _, f := range formas {
		src := "jest.mock('" + f.spec + "', () => ({ f: jest.fn() }))"
		v, msg := rodaMock(t, src, cfgComContrato())
		if v != Fail {
			t.Errorf("%s: resolve para um nó do mapa e devia ser cobrado; obteve %v (%s)", f.nome, v, msg)
		}
	}
	// E o que NÃO resolve no mapa não é cobrado — a diferença vem do grafo, não da forma.
	if v, _ := rodaMock(t, "jest.mock('@terceiro/nada', () => ({ f: jest.fn() }))", cfgComContrato()); v != Skip {
		t.Errorf("o que não resolve no grafo não é regido; obteve %v", v)
	}
}

// A forma de amarra não declarada PULA em vez de adivinhar. Inferir `Partial<typeof …>`
// assumiria TypeScript e reportaria verde sobre o que não se conferiu em qualquer
// outro ecossistema.
func TestMockTipado_formaNaoDeclaradaPulaEmVezDeAdivinhar(t *testing.T) {
	t.Run("MCTYM-I02: An undeclared tie shape skips instead of guessing one", func(t *testing.T) {})
	src := `jest.mock('src/a', () => ({ f: jest.fn() }))`
	// Sem forma declarada: Skip.
	v, _ := rodaMock(t, src, &config.Config{Derived: &config.Derived{}})
	if v != Skip {
		t.Fatalf("sem forma declarada o gate se cala; obteve %v", v)
	}
	// Com a forma declarada: o MESMO arquivo é cobrado. Prova de que o Skip vem da
	// ausência de declaração, não da ausência de dublê solto.
	if v2, _ := rodaMock(t, src, cfgComContrato()); v2 != Fail {
		t.Fatalf("com a forma declarada o mesmo arquivo devia reprovar; obteve %v", v2)
	}
}

// O veredito CONTA os soltos, não só nomeia. Nomear um sem dizer quantos deixaria o
// leitor sem saber se os outros também estão na lista.
func TestMockTipado_veredictoContaOsSoltos(t *testing.T) {
	t.Run("MCTYM-I03: The verdict counts the loose doubles, not just names them", func(t *testing.T) {})
	src := `import type * as A from 'src/a'
jest.mock('src/a', (): Partial<typeof A> => ({ f: jest.fn() }))
jest.mock('src/b', () => ({ g: jest.fn() }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Fail {
		t.Fatalf("um solto entre dois devia reprovar; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "1 dublê(s) sem amarra") && !strings.Contains(msg, "1 double(s) without binding") {
		t.Errorf("o veredito devia dizer QUANTOS estão soltos; msg = %q", msg)
	}
	if strings.Contains(msg, "2 dublê") || strings.Contains(msg, "2 double") {
		t.Errorf("o amarrado não entra na conta; msg = %q", msg)
	}
}

// O gate não confere se o tipo anotado CORRESPONDE ao módulo — quem confere é o
// compilador da linguagem, que já faz isso e faz melhor. Aqui a régua é que a amarra
// FOI ESCRITA, e é escrevê-la que entrega a conferência ao compilador.
func TestMockTipado_naoConfereSeOTipoCorresponde(t *testing.T) {
	t.Run("MCTYM-X01: The gate does not check whether the annotated type matches the real module", func(t *testing.T) {})
	src := `import type * as Outro from 'totalmente/outra/coisa'
jest.mock('src/a', (): Partial<typeof Outro> => ({ f: jest.fn() }))`
	if v, msg := rodaMock(t, src, cfgComContrato()); v != Pass {
		t.Fatalf("a régua é a presença da amarra, não a sua correção; obteve %v (%s)", v, msg)
	}
}

// LIMITE explícito: o gate cobre drift de FORMA, não de COMPORTAMENTO. Uma fábrica
// amarrada que devolve a forma certa com o sentido errado passa — e prometer mais
// seria vender um verde que este gate não dá.
func TestMockTipado_naoCobreDriftDeComportamento(t *testing.T) {
	t.Run("MCTYM-X02: The gate does not cover drift of behaviour", func(t *testing.T) {})
	src := `import type * as A from 'src/a'
jest.mock('src/a', (): Partial<typeof A> => ({
  calcularSaldo: jest.fn(() => -1),
}))`
	if v, msg := rodaMock(t, src, cfgComContrato()); v != Pass {
		t.Fatalf("forma certa e sentido errado ainda passa — é o limite declarado; obteve %v (%s)", v, msg)
	}
}

// O gate não carrega forma nem ecossistema embutidos. Num projeto onde o dublê não é
// uma chamada a detectar (em Go é uma interface satisfeita), sem nada declarado ele
// PULA: o veredito certo é "não se aplica", não um verde sobre o que não se conferiu.
func TestMockTipado_semFormaNemDialetoNaoSeAplica(t *testing.T) {
	t.Run("MCTYM-X03: The gate carries no built-in tie shape and no built-in ecosystem", func(t *testing.T) {})
	// Um teste Go: o dublê é um tipo que satisfaz a interface — não há chamada.
	src := `type repoFalso struct{}
func (repoFalso) Listar() []int { return nil }
func TestX(t *testing.T) { usar(repoFalso{}) }`
	v, msg := checkMockTyped(src, mapx.Node{ID: "x_test.go", Kind: mapx.KindTest},
		"", grafoDoProjeto(), &config.Config{Derived: &config.Derived{}})
	if v != Skip {
		t.Fatalf("sem forma nem dialeto declarados o gate não se aplica; obteve %v (%s)", v, msg)
	}
	if v == Pass {
		t.Error("verde sobre o que não se conferiu é a pior falha de um medidor")
	}
}

// Terceiro não é cobrado, e a razão é medida: sem este recorte o gate acusou 305
// arquivos de uma vez num projeto real, 245 deles por dublê de biblioteca. Gate que
// acusa tudo é desligado, e leva junto os achados legítimos.
func TestMockTipado_terceirosEmVolumeNaoSaoCobrados(t *testing.T) {
	t.Run("MCTYM-X04: The gate does not charge third-party library doubles", func(t *testing.T) {})
	src := `jest.mock('@react-navigation/native', () => ({ useNavigation: jest.fn() }))
jest.mock('@aws-sdk/client-s3', () => ({ S3: jest.fn() }))
jest.mock('@gorhom/bottom-sheet', () => ({ BottomSheet: 'View' }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Skip {
		t.Fatalf("três dublês de terceiro não produzem achado nenhum; obteve %v (%s)", v, msg)
	}
	for _, lib := range []string{"react-navigation", "aws-sdk", "gorhom"} {
		if strings.Contains(msg, lib) {
			t.Errorf("terceiro não devia ser nomeado: %s em %q", lib, msg)
		}
	}
}

// O modo spy do Vitest (`vi.mock('x', { spy: true })`) carrega o módulo REAL e só envolve
// as funções — não tem fábrica, e não há o que anotar. O `[^=]*` de antes atravessava
// quebras de linha e atribuía ao spy a fábrica de um mock POSTERIOR; medido no app de
// referência, 3 de 4 "dublês soltos" eram exatamente isso.
func TestMockTipado_spyNaoEhFabricaDeOutroMock(t *testing.T) {
	t.Run("MCTYM-B06: A spy mock is not charged with a later mock's factory", func(t *testing.T) {})
	src := `vi.mock('src/a', { spy: true })
vi.mock('src/b', (): Partial<typeof import('src/b')> => ({
  algo: vi.fn(),
}))`
	if v, msg := rodaMock(t, src, cfgComContrato()); v != Pass {
		t.Errorf("spy seguido de fábrica anotada deveria passar: %v (%s)", v, msg)
	}
}

// A fábrica na linha SEGUINTE (como o prettier quebra) continua sendo lida — e sua
// anotação continua contando. Sem o `\s*` antes do grupo, a quebra derrubaria o casamento
// e o dublê anotado passaria por não ter fábrica.
func TestMockTipado_fabricaNaLinhaSeguinte(t *testing.T) {
	t.Run("MCTYM-B07: A factory on the next line is still read, tie and all", func(t *testing.T) {})
	anotado := `vi.mock(
  'src/a',
  (): Partial<typeof import('src/a')> => ({ x: vi.fn() }),
)`
	if v, msg := rodaMock(t, anotado, cfgComContrato()); v != Pass {
		t.Errorf("fábrica anotada na linha seguinte deveria passar: %v (%s)", v, msg)
	}
	solto := `vi.mock(
  'src/a',
  () => ({ x: vi.fn() }),
)`
	if v, _ := rodaMock(t, solto, cfgComContrato()); v != Fail {
		t.Errorf("fábrica SEM anotação na linha seguinte deveria reprovar: %v", v)
	}
}

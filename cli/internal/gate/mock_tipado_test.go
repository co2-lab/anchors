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
	return checkMockTipado(content, testNode(), "", grafoDoProjeto(), cfg)
}

// O caso que o gate existe para pegar: dublê SEM amarra continua verde depois que o
// módulo muda de assinatura. É prova falsa, não ausência de prova.
func TestMockTipado_dubleSemAmarraReprova(t *testing.T) {
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
	src := `import type * as A from 'a'
jest.mock('a', (): Partial<typeof A> => ({ f: jest.fn() }))
jest.mock('b', () => ({ g: jest.fn() }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Fail {
		t.Fatalf("um dublê solto entre vários basta para reprovar: %v", v)
	}
	// A contagem é o que separa "um solto" de "todos soltos": nomear `b` sem dizer
	// quantos deixaria o leitor sem saber se `a` também está na lista.
	if !strings.Contains(msg, "1 dublê(s) sem amarra") {
		t.Errorf("deve contar apenas o solto: %s", msg)
	}
	// O módulo amarrado NÃO entra na lista de soltos. (A verificação é sobre a lista,
	// não sobre a mensagem inteira: o texto do contrato declarado — que contém
	// `typeof` — aparece legitimamente na instrução de conserto ao final.)
	lista := strings.SplitN(strings.SplitN(msg, "real: ", 2)[1], ".", 2)[0]
	if strings.Contains(lista, "a") {
		t.Errorf("o módulo amarrado não deve ser acusado: lista=%q", lista)
	}
}

// `jest.mock('x')` sem fábrica usa o automock, que deriva do módulo real por
// construção — não drifta, e cobrá-lo seria ruído que se aprende a ignorar.
func TestMockTipado_automockNaoEhCobrado(t *testing.T) {
	src := `jest.mock('@backend/repositories/lotes')`
	if v, msg := rodaMock(t, src, cfgComContrato()); v != Skip {
		t.Errorf("automock não tem fábrica que possa mentir: %v (%s)", v, msg)
	}
}

// Sem `derived.mock_contract` o gate PULA. Inferir `Partial<typeof …>` assumiria
// TypeScript e reportaria verde sobre o que não se conferiu em qualquer outro
// ecossistema — a pior falha possível num medidor.
func TestMockTipado_semContratoDeclaradoPula(t *testing.T) {
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
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkMockTipado("", n, "", grafoDoProjeto(), cfgComContrato()); v != Skip {
		t.Errorf("o dublê é cobrado do teste: %v", v)
	}
}

// Teste que não dubla ninguém não tem o que amarrar — Skip, não Pass: não houve
// verificação, e um Pass aqui inflaria a contagem de verdes com nada.
func TestMockTipado_semDubleNadaACobrar(t *testing.T) {
	src := `it('soma', () => { expect(1+1).toBe(2) })`
	if v, _ := rodaMock(t, src, cfgComContrato()); v != Skip {
		t.Errorf("sem dublê não há amarra a cobrar: %v", v)
	}
}

// O dialeto do runner é do projeto: `vi.mock` (vitest) conta igual a `jest.mock`.
func TestMockTipado_reconheceVitest(t *testing.T) {
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
	src := `jest.mock('@gorhom/bottom-sheet', () => ({ BottomSheet: 'View' }))
jest.mock('@backend/repositories/lotes', () => ({ listar: jest.fn() }))`
	v, msg := rodaMock(t, src, cfgComContrato())
	if v != Fail {
		t.Fatalf("o próprio solto deve reprovar mesmo cercado de terceiros: %v", v)
	}
	if strings.Contains(msg, "gorhom") {
		t.Errorf("terceiro não deve ser acusado: %s", msg)
	}
	if !strings.Contains(msg, "1 dublê(s) sem amarra") {
		t.Errorf("deve contar só o próprio: %s", msg)
	}
}

// Import relativo (`../stores/auth.store`) é módulo próprio como qualquer outro — o
// critério é resolver para um nó do mapa, não a forma do especificador.
func TestMockTipado_relativoEhProprio(t *testing.T) {
	src := `jest.mock('../a', () => ({ f: jest.fn() }))`
	if v, _ := rodaMock(t, src, cfgComContrato()); v != Fail {
		t.Errorf("import relativo que resolve no mapa é próprio: %v", v)
	}
}

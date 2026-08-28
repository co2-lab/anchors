package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// mock-tipado: todo dublê de teste precisa DERIVAR do módulo que substitui.
//
// O buraco que ele fecha é o mais traiçoeiro do conjunto de gates, porque não é
// ausência de prova — é prova FALSA. Um teste que dubla o vizinho continua passando
// depois que o vizinho muda de assinatura, de retorno ou de nome: o dublê virou uma
// cópia congelada de um contrato que já não existe. O verde certifica a versão antiga,
// e ninguém vai procurar defeito onde há teste verde.
//
// Os outros gates da família não alcançam isto:
//   - `trinca-completa` pergunta se o teste EXISTE — e ele existe;
//   - `feature-test-match` confronta cenário contra caso de teste — e eles casam;
//   - `tests-green` lê o resultado da execução — e ela passou.
//
// Todos respondem "sim" sobre um teste que mente.
//
// A defesa não é o gate reimplementar checagem de tipo: é exigir que o dublê seja
// AMARRADO ao original por um mecanismo que a linguagem já sabe conferir. Em
// TypeScript, anotar a fábrica com `Partial<typeof Real>` faz o compilador acusar tanto
// método inexistente quanto retorno divergente; em Python é `autospec=True`; em Go, a
// interface. O framework não sabe qual é — o projeto declara em
// `derived.mock_contract`, e este gate cobra a presença da amarra.
//
// LIMITE, explícito: isto cobre drift de FORMA (nome, assinatura, tipo de retorno).
// NÃO cobre drift de COMPORTAMENTO — se o módulo real passa a lançar num caso, ou muda
// de semântica mantendo a assinatura, o dublê segue mentindo e nenhum compilador vê.
// Para esse resto existe julgamento (`mock-nao-replica-regra`) e teste de integração na
// borda; prometer mais do que a forma seria vender o verde que este gate não dá.
func checkMockTipado(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, "o dublê vive no teste — é dele que a amarra é cobrada"
	}
	forma := contratoDeMock(cfg)
	if forma == "" {
		// Sem a forma declarada não há o que procurar. Inferir uma (`Partial<typeof …>`)
		// seria assumir TypeScript e reportar VERDE sobre o que não se conferiu em
		// qualquer outro ecossistema — a pior falha possível num medidor.
		return Skip, "o projeto não declara `derived.mock_contract` — não há amarra de dublê a confrontar"
	}

	todos := dublesDeclarados(content)
	// Só o que o PROJETO rege é cobrado. Ver `ehModuloRegido`.
	var dubles []dubleDeclarado
	for _, d := range todos {
		if ehModuloRegido(d.modulo, g) {
			dubles = append(dubles, d)
		}
	}
	if len(dubles) == 0 {
		if len(todos) > 0 {
			return Skip, "o teste só dubla módulo de fora do projeto — a amarra é cobrada do que o projeto rege"
		}
		return Skip, "o teste não dubla módulo nenhum — nada a amarrar"
	}

	// A amarra é conferida por MÓDULO, não pelo arquivo: um teste que tipa três dublês
	// e esquece o quarto tem exatamente o buraco que este gate existe para pegar, e um
	// veredito por arquivo o daria como resolvido.
	var soltos []string
	for _, d := range dubles {
		if !d.amarrado {
			soltos = append(soltos, d.modulo)
		}
	}
	if len(soltos) == 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(
		"%d dublê(s) sem amarra ao módulo real: %s.\n\nUm dublê solto continua passando "+
			"depois que o módulo muda — o teste vira uma cópia congelada do contrato "+
			"antigo, e o verde certifica a versão que não existe mais. Amarre-o à forma "+
			"que este projeto declara (`derived.mock_contract`: %s), para o compilador "+
			"conferir nome, assinatura e retorno em vez de ninguém conferir",
		len(soltos), strings.Join(soltos, ", "), forma)
}

// contratoDeMock lê a forma de amarra declarada pelo projeto.
func contratoDeMock(cfg *config.Config) string {
	// `Derived` é opcional na config — um projeto que não declara superfície de trinca
	// o deixa nil, e desreferenciá-lo derrubaria o check inteiro em vez de pular.
	if cfg == nil || cfg.Derived == nil {
		return ""
	}
	return strings.TrimSpace(cfg.Derived.MockContract)
}

type dubleDeclarado struct {
	modulo   string
	amarrado bool
}

// ehModuloRegido diz se o especificador do dublê aponta para código DO PROJETO.
//
// A cobrança vale para o que o projeto rege, e não para biblioteca de terceiro. O drift
// que este gate persegue é "o vizinho mudou e o dublê não soube" — e o vizinho que muda
// toda semana é o módulo próprio. A dependência externa tem versão travada no lockfile,
// e o dublê dela costuma trocar um componente por um stub, não reproduzir um contrato:
// exigir amarra ali seria burocracia sem defeito correspondente.
//
// Há uma razão mais dura, medida: no app de referência o gate acusou 305 arquivos de uma vez, dos
// quais 245 dublês eram de terceiros (`@react-navigation/native`, `@aws-sdk/*`,
// `@gorhom/bottom-sheet`). Um gate que acusa tudo não é lido — é desligado, e leva
// junto os achados legítimos.
//
// O critério é o GRAFO, não uma lista de prefixos na config: um módulo é do projeto se
// resolve para um nó do mapa. Isso não pede configuração nova (o mapa já existe), não
// assume convenção de alias (`@/`, `~/`, `src/` variam por ecossistema) e acompanha o
// projeto sozinho — código que nasce entra no mapa e passa a ser cobrado.
func ehModuloRegido(spec string, g *mapx.Graph) bool {
	if g == nil {
		return false
	}
	// O especificador é um caminho de import (`@backend/repositories/lotes`,
	// `../stores/auth.store`), sem extensão; o nó do mapa é um caminho de arquivo com
	// ela. Casar pelo SUFIXO do caminho sem extensão cobre os dois formatos sem
	// precisar resolver alias.
	alvo := strings.TrimPrefix(spec, "./")
	for strings.HasPrefix(alvo, "../") {
		alvo = strings.TrimPrefix(alvo, "../")
	}
	if i := strings.Index(alvo, "/"); i >= 0 {
		// Descarta o primeiro segmento quando ele é um alias (`@backend`, `@`): o que
		// identifica o módulo é a cauda do caminho.
		if strings.HasPrefix(alvo, "@") {
			alvo = alvo[i+1:]
		}
	}
	if alvo == "" {
		return false
	}
	for _, n := range g.Nodes {
		if semExtensao(n.ID) == alvo || strings.HasSuffix(semExtensao(n.ID), "/"+alvo) {
			return true
		}
	}
	return false
}

// semExtensao tira a extensão final do caminho (`x/y.ts` → `x/y`), para o nó do mapa
// poder ser comparado com um especificador de import, que nunca a carrega.
func semExtensao(id string) string {
	if i := strings.LastIndex(id, "."); i > strings.LastIndex(id, "/") {
		return id[:i]
	}
	return id
}

// mockComFabricaRE casa a declaração de dublê COM fábrica — a única forma que pode
// mentir. `jest.mock('x')` sem fábrica usa o automock, que deriva do módulo real por
// construção e portanto não drifta; cobrá-lo seria ruído.
//
// O grupo 2 captura o que vem entre a vírgula e a seta da fábrica: é onde a anotação de
// tipo apareceria (`(): Partial<typeof Real> =>`). Casar a fábrica inteira exigiria
// equilibrar chaves, que regex não faz — e não é preciso: a amarra, quando existe, está
// sempre antes do corpo.
var mockComFabricaRE = regexp.MustCompile(
	`(?:jest|vi)\s*\.\s*mock\s*\(\s*['"` + "`" + `]([^'"` + "`" + `]+)['"` + "`" + `]\s*,([^=]*)=>`)

// dublesDeclarados inventaria os dublês do arquivo e diz quais têm amarra.
//
// "Ter amarra" é o trecho entre a vírgula e a seta mencionar o módulo dublado por um
// tipo — o que a forma declarada descreve. A checagem é pela presença de uma ANOTAÇÃO
// (`:` seguido de tipo) que referencie o módulo, não pela igualdade literal com o
// template: o projeto escreve `Partial<typeof Real>` com o alias que quiser, e exigir o
// texto exato transformaria o gate num verificador de estilo.
func dublesDeclarados(content string) []dubleDeclarado {
	var out []dubleDeclarado
	for _, m := range mockComFabricaRE.FindAllStringSubmatch(content, -1) {
		modulo, cabeca := m[1], m[2]
		out = append(out, dubleDeclarado{
			modulo:   modulo,
			amarrado: temAnotacaoDeTipo(cabeca),
		})
	}
	return out
}

// anotacaoRE — a fábrica anotada tem `): <Tipo>` entre os parâmetros e a seta. Aceita
// qualquer tipo: quem confere se ele bate com o módulo é o compilador da linguagem, não
// este gate. Aqui só se verifica que a amarra FOI ESCRITA.
var anotacaoRE = regexp.MustCompile(`\)\s*:\s*\S`)

func temAnotacaoDeTipo(cabeca string) bool {
	return anotacaoRE.MatchString(cabeca)
}

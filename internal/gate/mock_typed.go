// @anchors
//   code: MTGMC
//   ref: MCTYM

package gate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
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
//   - `unit-complete` pergunta se o teste EXISTE — e ele existe;
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
func checkMockTyped(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, i18n.T("gate.mock_typed.skip_not_test")
	}
	forma := mockContract(cfg)
	if forma == "" {
		// Sem a forma declarada não há o que procurar. Inferir uma (`Partial<typeof …>`)
		// seria assumir TypeScript e reportar VERDE sobre o que não se conferiu em
		// qualquer outro ecossistema — a pior falha possível num medidor.
		// Pending until someone decides, as `mock-stamped`: the opt-out is the decision.
		if cfg != nil && cfg.DialectFor().WaivedField("mock_contract") {
			return Skip, i18n.T("gate.mock_typed.skip_opt_out")
		}
		return Pending, i18n.T("gate.mock_typed.pending_no_mock_contract")
	}

	todos := declaredDoubles(content, cfg, forma)
	// With no map, which doubles are the PROJECT's cannot be told (MCTYM-E02): answering
	// "it only mocks modules outside the project" would state a cause nobody measured.
	if g == nil && len(todos) > 0 {
		return pendingNoMap()
	}
	// Só o que o PROJETO rege é cobrado. Ver `ehModuloRegido`.
	var dubles []declaredDouble
	for _, d := range todos {
		if isGovernedModule(d.module, n.ID, cfg.DialectFor(), g) {
			dubles = append(dubles, d)
		}
	}
	if len(dubles) == 0 {
		if len(todos) > 0 {
			return Skip, i18n.T("gate.mock_typed.skip_external_only")
		}
		return Skip, i18n.T("gate.mock_typed.skip_no_doubles")
	}

	// A amarra é conferida por MÓDULO, não pelo arquivo: um teste que tipa três dublês
	// e esquece o quarto tem exatamente o buraco que este gate existe para pegar, e um
	// veredito por arquivo o daria como resolvido.
	var soltos []string
	for _, d := range dubles {
		if !d.bound {
			soltos = append(soltos, d.module)
		}
	}
	if len(soltos) == 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(i18n.T("gate.mock_typed.fail_unbound"),
		len(soltos), strings.Join(soltos, ", "), forma)
}

// mockContract lê a forma de amarra declarada pelo projeto.
func mockContract(cfg *config.Config) string {
	// `Derived` é opcional na config — um projeto que não declara superfície de unidade
	// o deixa nil, e desreferenciá-lo derrubaria o check inteiro em vez de pular.
	if cfg == nil || cfg.Derived == nil {
		return ""
	}
	return strings.TrimSpace(cfg.Derived.MockContract)
}

type declaredDouble struct {
	module string
	bound  bool
}

// isGovernedModule says whether a double's specifier names one code file of the map — a
// module of the project, not a third-party one (moduleFile).
func isGovernedModule(spec, testID string, d config.Dialect, g *mapx.Graph) bool {
	return len(moduleCandidates(spec, testID, d, g)) > 0
}

// withoutExtension tira a extensão final do caminho (`x/y.ts` → `x/y`), para o nó do mapa
// poder ser comparado com um especificador de import, que nunca a carrega.
func withoutExtension(id string) string {
	if i := strings.LastIndex(id, "."); i > strings.LastIndex(id, "/") {
		return id[:i]
	}
	return id
}

// declaredDoubles inventaria os dublês do arquivo e diz quais têm amarra, pelo que o projeto
// declara (`derived.mock_detect`, `mock_factory`, `mock_typed`): o núcleo não lê sintaxe de
// framework de teste nenhum — sem `mock_detect` não há dublê a achar, e o gate pressupõe o
// campo.
//
// A double with no factory — an automock, Vitest's spy mode — mirrors the real module by
// construction: there is no hand-written contract to go stale. `mock_factory` says what a
// double WITH a factory looks like (measured in the project that declared `mock_detect` for
// `mock-stamped`: 51 doubles charged, 42 automocks and 9 spy mocks, not one hand-written
// factory); with none declared, every double is charged. `mock_typed` says what binds a
// factory besides the contract's own form (a TS return annotation, `(): T =>`).
func declaredDoubles(content string, cfg *config.Config, forma string) []declaredDouble {
	re, err := doubleDetector(cfg)
	if err != nil || re == nil {
		return nil
	}
	var factory, typed *regexp.Regexp
	if cfg.Derived != nil {
		// Anchored at the double's own call: a later double's factory is not this one's.
		if strings.TrimSpace(cfg.Derived.MockFactory) != "" {
			factory = compileOrNil(`\A(?:` + cfg.Derived.MockFactory + `)`)
		}
		typed = compileOrNil(cfg.Derived.MockTyped)
	}
	var out []declaredDouble
	linhas := strings.Split(content, "\n")
	// Over the whole file: a call broken across lines (prettier's `vi.mock(\n  'x',`) is
	// one double.
	for _, loc := range re.FindAllStringSubmatchIndex(content, -1) {
		if len(loc) < 4 || loc[2] < 0 || strings.TrimSpace(content[loc[2]:loc[3]]) == "" {
			continue
		}
		modulo := strings.TrimSpace(content[loc[2]:loc[3]])
		i := config.LineOf(content, loc[0]) - 1
		// From the line the call opens on: the factory pattern reads the call whole, a
		// prettier-broken one included.
		lineStart := strings.LastIndex(content[:loc[0]], "\n") + 1
		if strings.Contains(forma, "{{module}}") && factory != nil && !factoryAt(factory, content, lineStart, loc[0]) {
			continue
		}
		// Contexto: a linha do mock mais as próximas 3 linhas (para chamadas multi-linha)
		fim := i + 4
		if fim > len(linhas) {
			fim = len(linhas)
		}
		ctx := strings.Join(linhas[i:fim], " ")
		bound := false
		if strings.Contains(forma, "{{module}}") {
			specFmt := strings.ReplaceAll(forma, "{{module}}", filepath.Base(modulo))
			bound = strings.Contains(ctx, specFmt) || (typed != nil && typed.MatchString(ctx))
		} else {
			bound = strings.Contains(ctx, forma)
		}
		out = append(out, declaredDouble{module: modulo, bound: bound})
	}
	return out
}

// factoryAt says whether the factory pattern, anchored, matches from the start of the
// double's line or from the double itself — the pattern may read what precedes the call
// on its line, or start at it.
func factoryAt(factory *regexp.Regexp, content string, lineStart, at int) bool {
	return factory.MatchString(content[at:]) || factory.MatchString(strings.TrimLeft(content[lineStart:], " \t"))
}

// compileOrNil is the pattern compiled, or nil when empty or broken.
func compileOrNil(p string) *regexp.Regexp {
	if strings.TrimSpace(p) == "" {
		return nil
	}
	re, err := regexp.Compile(p)
	if err != nil {
		return nil
	}
	return re
}

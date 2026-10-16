// @anchors
//   code: TQEGT
//   ref: TQETS

package gate

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkQueriedTestIDExists — a QUARTA aresta do contrato de testID, a que nenhum
// gate cobria: o flow de ponta a ponta CONSULTA um handle que o código não expõe.
//
// Por que não cabe no `testid-coerente`: aquele gate parte de UMA spec, e a superfície
// e2e é a árvore INTEIRA de flows do projeto — um id de outra tela aparece nela do
// mesmo jeito. Cobrar da spec do InventoryMove um `:acnw-account-new-*` seria acusá-la
// por um handle de contas bancárias; medido ao tentar, deu 829 achados numa spec só,
// todos de outras telas. A pergunta é do PROJETO ("existe flow procurando handle que
// ninguém expõe?"), então o escopo é `project` e os dois lados são varridos de uma vez.
//
// O que este gate pega, medido no app de referência (2026-08-25): 13 ids inventados nos flows, com
// prefixo errado (`review-*` onde a tela emite `revi-*`). A premissa de que "um id
// inexistente já falha ao rodar" é FALSA por dois motivos independentes:
//
//   - o flow precisa CHEGAR à linha. Numa suíte cujo passo 3 já falhava, os ids
//     fantasma dos passos seguintes eram invisíveis — apareceram um por rodada, ao
//     longo de sete medições.
//   - num `assertNotVisible`, o id inexistente faz o teste PASSAR. O REVI-R03
//     asseverava `:review-edit-controls` (que nunca existiu) e provava o oposto do que
//     aparentava: verde por VACUIDADE. É o caso mais grave, e o runner nunca o dá.
//
// Direção ÚNICA (flow → código), deliberadamente. O contrário (código expõe handle que
// nenhum flow consulta) NÃO é ofensa: nem todo elemento marcado precisa de cenário
// automatizado, e o `testid-coerente` já reporta o consumo por spec, onde a pergunta
// tem dono.
func checkQueriedTestIDExists(_ string, _ mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	// The core knows no syntax of a handle, in the code or in a flow: with no pattern
	// declared, there is nothing to read either side with — and assuming one would answer
	// green over what it never read.
	th := testHandles(cfg)
	switch {
	case th == nil || len(th.Exposed) == 0:
		return Pending, i18n.T("gate.testid.no_patterns", "dialect.test_handles.exposed")
	case len(th.Queried) == 0:
		return Pending, i18n.T("gate.testid.no_patterns", "dialect.test_handles.queried")
	case g == nil:
		return pendingNoMap()
	}
	flows := e2eSurfaceFiles(root, cfg)
	if len(flows) == 0 {
		// Sem superfície declarada não há o que confrontar. Skip é honesto: "não medi"
		// não é "está limpo".
		return Skip, i18n.T("gate.testid_queried_exists.skip_no_e2e_surface")
	}

	// O universo do que o CÓDIGO expõe: every code file of the map, not only the ones a
	// spec claims — a handle may be born in a component no screen spec claims (the
	// navigator's `tabBarButtonTestID`: `:tab-import` belongs to no screen), and charging
	// its absence would accuse who complies.
	expostos, codigoIlegivel := projectExposedHandles(root, g, th)
	if len(expostos) == 0 && len(codigoIlegivel) == 0 {
		return Skip, i18n.T("gate.testid_queried_exists.skip_no_handles_found")
	}

	type achado struct{ id, onde string }
	var achados []achado
	vistos := map[string]bool{}

	// What could not be READ is never approved (TQETS-E02): an unreadable flow may be
	// the one querying the missing handle, and an unreadable source may be the one that
	// exposes it. Skipping them made the gate answer Pass over a flow it never read.
	var fluxoIlegivel []string
	for _, f := range flows {
		rel, _ := filepath.Rel(root, f.path)
		rel = filepath.ToSlash(rel) // named as the project writes its paths, on every system
		b, err := os.ReadFile(f.path)
		if err != nil {
			fluxoIlegivel = append(fluxoIlegivel, rel)
			continue
		}
		for _, id := range queriedHandles(string(b), rel, th.Queried) {
			if vistos[id+"|"+rel] {
				continue
			}
			if handleExists(expostos, id) {
				continue
			}
			vistos[id+"|"+rel] = true
			achados = append(achados, achado{id: id, onde: rel})
		}
	}

	if ilegiveis := append(fluxoIlegivel, codigoIlegivel...); len(ilegiveis) > 0 {
		sort.Strings(ilegiveis)
		return Pending, i18n.T("gate.testid_queried_exists.unreadable", len(ilegiveis), strings.Join(ilegiveis, ", "))
	}
	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Slice(achados, func(i, j int) bool {
		if achados[i].id != achados[j].id {
			return achados[i].id < achados[j].id
		}
		return achados[i].onde < achados[j].onde
	})
	// Agrupa por id: o mesmo handle inventado costuma aparecer em vários flows, e um
	// achado por arquivo faria um defeito parecer sete.
	porID := map[string][]string{}
	var ordem []string
	for _, a := range achados {
		if _, ok := porID[a.id]; !ok {
			ordem = append(ordem, a.id)
		}
		porID[a.id] = append(porID[a.id], a.onde)
	}
	var linhas []string
	for _, id := range ordem {
		linhas = append(linhas, i18n.T("gate.testid_queried_exists.line_queried_in", id, strings.Join(porID[id], ", ")))
	}
	return Fail, i18n.T("gate.testid_queried_exists.handles_missing",
		len(ordem), strings.Join(linhas, "\n"))
}

// e2eFile é um flow da superfície de ponta a ponta.
type e2eFile struct{ path string }

// e2eSurfaceFiles devolve os CAMINHOS dos flows (o `lerSuperficieE2E` irmão
// devolve só o conteúdo, e aqui o caminho é o que o laudo precisa nomear para o
// achado ser acionável).
//
// A resolução é em DOIS passos e confundi-los custa caro: `surfaces[e2e]` devolve a
// CHAVE da superfície (ex.: "e2e"), não um caminho — quem tem o caminho é
// `files[chave]`.
func e2eSurfaceFiles(root string, cfg *config.Config) []e2eFile {
	if cfg == nil || cfg.Derived == nil {
		return nil
	}
	chave := cfg.Derived.Surfaces["e2e"]
	if chave == "" {
		return nil
	}
	// O PRIMEIRO padrão da camada: este gate resolve UM caminho esperado, e uma spec com
	// vários padrões (a de configuração) não muda o que ele pergunta.
	var padrao string
	if ps := cfg.Derived.PadroesDe()[chave]; len(ps) > 0 {
		padrao = ps[0]
	}
	if padrao == "" {
		for _, ov := range cfg.Derived.Overrides {
			if ps := ov.PadroesDe()[chave]; len(ps) > 0 {
				padrao = ps[0]
				break
			}
		}
	}
	if padrao == "" {
		return nil
	}
	dir := filepath.Join(root, firstStaticSegment(padrao))
	var out []e2eFile
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// Every file of the surface: what a flow is, and how it queries, is the project's
		// `queried` patterns' — a utility of the runner beside the flows matches none.
		out = append(out, e2eFile{path: p})
		return nil
	})
	return out
}

// queriedHandles extracts the handles a flow queries, by the project's `queried` patterns:
// a `testid` group is the handle, a `prefix` group the head of a query by pattern
// (`id: ':revi-.*'`), confronted as a head. A query the patterns do not capture — one the
// runner composes at run time — is not confronted: the gate cannot know its final value.
func queriedHandles(src, rel string, patterns config.Patterns) []string {
	var out []string
	for _, h := range handlesIn(src, rel, patterns) {
		id := strings.TrimPrefix(h, ":")
		if strings.HasSuffix(id, "-*") {
			id = strings.TrimSuffix(id, "*")
		}
		id = strings.TrimRight(id, "-")
		// Cabeça curta demais não identifica nada: confrontá-la daria falso positivo em
		// qualquer projeto. O piso é o mesmo do resto do framework: um segmento com
		// prefixo de unidade + pelo menos uma palavra.
		if len(id) < 4 || !strings.Contains(id, "-") {
			continue
		}
		out = append(out, id)
	}
	return out
}

// projectExposedHandles reads every code file of the map with the `exposed` and the
// `mentioned` patterns, and returns every handle they find.
//
// The map's code, not a walk by extension: which files are source and which are tests is
// the project's layers', not a list of extensions in the core. A handle only a test
// mentions does not count — a test queries a handle, it does not create it; that is how a
// phantom id (`:recent-import-item` in the reference app) survived, referenced by a unit
// test. A component no spec claims is still a code node, so its handles count.
//
// The `mentioned` patterns answer a looser question than `exposed`: not "who owns this
// handle", but "does it exist anywhere in the code?" — a fallback (`btn.testID ?? …`), a
// table of handles, a suffix composed in a child. In a gate that ACCUSES absence, being
// loose (letting an orphan pass) costs far less than being strict (accusing who complies).
func projectExposedHandles(root string, g *mapx.Graph, th *config.TestHandles) (out, unreadable []string) {
	patterns := append(append(config.Patterns(nil), th.Exposed...), th.Mentioned...)
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(n.ID)))
		if err != nil {
			unreadable = append(unreadable, n.ID)
			continue
		}
		out = append(out, handlesIn(string(b), n.ID, patterns)...)
	}
	return out, unreadable
}

// handleExists confronta um handle consultado contra o universo exposto.
//
// Não é igualdade de string: o exposto pode ser um TEMPLATE (`item-*`, de
// “ testID={`:item-${id}`} “) e o consultado uma instância — ou o inverso, quando o
// flow procura por padrão. Ambos os lados podem carregar curinga, e os dois casam.
func handleExists(expostos []string, id string) bool {
	if coveredBy(expostos, id) {
		return true
	}
	// O consultado é a CABEÇA de um id composto no código: `:revi-review-tx` casa
	// `:revi-review-tx-*`. Sem isto, todo flow que procura por prefixo viraria achado.
	for _, e := range expostos {
		nu := strings.TrimPrefix(e, ":")
		if strings.HasPrefix(nu, id) {
			return true
		}
		// SUFIXO composto em componente filho (`*-row-*`): a cabeça vem de quem
		// renderiza, então só o miolo pode ser confrontado. Casa `ctdt-list-row-0`
		// contra o `-row-` que o TransactionList compõe.
		//
		// `Contains` E `HasSuffix`, e a segunda não é redundante: o id pode TERMINAR no
		// segmento composto, sem nada depois. `:reis-register-input-password-toggle`
		// nasce de `${testID}-toggle` no atom Input e não contém "-toggle-" — só
		// termina em "-toggle". Sem esta metade o gate acusava de inexistente um handle
		// que existe, e eu cheguei a "consertar" um flow correto por causa disso: troquei
		// o toque no olho por um toque em ponto neutro, que APAGOU a senha já digitada.
		if strings.HasPrefix(nu, "*") {
			if meio := strings.Trim(nu, "*"); meio != "" {
				if strings.Contains(id, meio) || strings.HasSuffix(id, strings.TrimSuffix(meio, "-")) {
					return true
				}
			}
		}
	}
	return false
}

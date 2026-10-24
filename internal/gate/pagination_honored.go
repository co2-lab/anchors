// @anchors
//   code: PHGPG
//   ref: PGNHN

package gate

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// pagination-honored: uma função que promete devolver um CONJUNTO ("liste todos", "os
// pendentes") não pode devolver a primeira página em silêncio.
//
// Esta é a classe de defeito de CUSTO/ESCALA — a que não aparece em teste nenhum porque o
// teste roda com 3 registros e a produção roda com 3 mil. Nada no código está "errado":
// a query é válida, o tipo é o esperado, a suíte passa. O defeito é a diferença entre o
// que o NOME promete e o que a função entrega quando o dado cresce.
//
// Os três casos, medidos num mesmo arquivo real (packages/backend/models/commissionLedger.ts):
//
//  1. `listCommissionsByReseller(id, limit)` — tem `Limit`, o chamador escolheu. É uma
//     página deliberada e o gate NÃO acusa: quem passou o limite sabe que há mais.
//  2. `listReleasableCommissions(nowIso, limit = 100)` — o limite é um DEFAULT que o
//     chamador não vê. O nome promete "as liberáveis"; devolve no máximo 100. A 101ª
//     comissão nunca é liberada, e ninguém é notificado disso. O consumidor é um lambda
//     agendado cujo teste mocka a função — o limite nunca aparece em teste.
//  3. `listCommissionsByOrganization(id)` — sem limite e sem paginação: trunca no 1 MB
//     do DynamoDB, silenciosamente, quando a organização crescer.
//
// O gate acusa 2 e 3, e usa uma prova por ASSIMETRIA quando ela existe: se o módulo tem
// irmãs que paginam (o loop `do { } while (lastKey)`), então o autor conhecia o padrão e
// a que não pagina é esquecimento, não decisão. Sem irmãs paginando, o veredito é mais
// fraco — só acusa quando o nome promete conjunto de forma inequívoca.
//
// Como dispensar honestamente: dê ao parâmetro de limite um nome que apareça na
// assinatura do chamador (`pageSize`), devolva o cursor, ou renomeie a função para o que
// ela faz (`firstNCommissions`). Um comentário `// @no-paginate: <razão>` na função
// também dispensa — opt-out explícito, como manda o CONCEPT §5.1 —, acima da assinatura ou
// dentro do corpo.
//
// AGNOSTICISMO — o que é do framework e o que é do projeto:
//   - do FRAMEWORK: a verdade confrontada ("nome promete conjunto, retorno é parcial"),
//     a distinção limite-do-chamador vs limite-escondido, e a prova por assimetria.
//     Nada disso é de linguagem alguma.
//   - do PROJETO: como se reconhece uma função exportada, um laço, um cursor e uma
//     consulta de coleção. Vem de `dialect` no anchors.yaml (internal/config/dialect.go).
//     A consulta de coleção é a mais específica de stack e NÃO tem default: um projeto
//     DynamoDB e um projeto SQL não se parecem, então o Anchors não chuta — sem
//     `dialect.collection_query`, este gate se declara Pendente.
func checkPaginationHonored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, i18n.T("gate.pagination_honored.skip_not_code")
	}
	d := cfg.DialectFor()
	// Cada padrão ausente é um olho fechado. Pendente nomeia QUAL falta, para a correção
	// ser uma linha de YAML e não uma investigação.
	var faltando []string
	if d.ExportedFunc == "" && !d.WaivedField("exported_func") {
		faltando = append(faltando, i18n.T("gate.pagination_honored.missing_exported_func"))
	}
	if d.CollectionQuery == "" && !d.WaivedField("collection_query") {
		faltando = append(faltando, i18n.T("gate.pagination_honored.missing_collection_query"))
	}
	if len(faltando) > 0 {
		// Pending, não Skip: o gate NÃO verificou, e dizer isso é o ponto. Silenciar aqui
		// faria o relatório afirmar conformidade sobre o que ninguém mediu.
		//
		// Quem decidiu não configurar tem o opt-out explícito, e a mensagem o oferece —
		// senão a única saída visível seria conviver com o aviso para sempre, e aviso
		// permanente é aviso que se aprende a ignorar (levando os outros com ele).
		return Pending, i18n.T("gate.pagination_honored.pending_dialect_missing",
			strings.Join(faltando, " / "),
			strings.Join(config.KnownDialectFamilies(), ", "))
	}
	// Campo dispensado = o projeto afirmou que não tem. Sem o padrão não há como confrontar,
	// e o Skip aqui é honesto: diferente do Pending acima, alguém decidiu.
	if d.CollectionQuery == "" || d.ExportedFunc == "" {
		return Skip, i18n.T("gate.pagination_honored.skip_opt_out")
	}

	// Comments are no code: a `do` in a Portuguese comment ("o uso do mês") is no loop, and a
	// query a comment mentions queries nothing. The code is read with them blanked, offsets
	// kept; the waiver, which is a comment, is read on the content.
	code := blankComments(content, n.ID, d)
	consulta := d.Compile(d.CollectionQuery)
	if consulta == nil || !consulta.MatchString(code) {
		return Skip, i18n.T("gate.pagination_honored.skip_no_collection_queries")
	}

	fns := exportedFuncs(code, d)
	if len(fns) == 0 {
		return Skip, i18n.T("gate.pagination_honored.skip_no_exported_funcs")
	}

	// A function drains the cursor itself, or calls one of the module's own that does — the
	// `readAll` that walks every page for the functions that read through it.
	helpers := drainingHelpers(code, d)
	drains := func(f exportedFunc) bool { return paginatesAll(f.body, d) || callsAny(f.body, f.name, helpers) }

	// O módulo conhece o padrão de paginação? Se alguma irmã pagina, a que não pagina é
	// assimetria — o mesmo raciocínio do gate sibling-guard.
	sistersPaginate := 0
	for _, f := range fns {
		if drains(f) {
			sistersPaginate++
		}
	}

	var achados []string
	for i, f := range fns {
		if !promisesSet(f.name, d) || !consulta.MatchString(f.body) {
			continue
		}
		if drains(f) || explicitWaiver(ownText(content, code, fns, i)) {
			continue
		}
		// NÃO existe regra de "cursor descartado" aqui, e a ausência é uma decisão medida.
		//
		// A ideia (do Adriel) é certa: uma função que recebe cursor do provedor e não o
		// devolve impede QUALQUER paginação — a informação de que há mais morre ali. O caso
		// real era `cognitoListDevices`, que fazia `return res.Devices ?? []` e jogava fora
		// o `res.PaginationToken`.
		//
		// Mas um gate que lê TEXTO não consegue ver isso: o `PaginationToken` nunca aparece
		// no código original. Descartar um campo é precisamente NÃO mencioná-lo, e não há
		// como detectar por regex a ausência de algo que nunca foi escrito — só um analisador
		// que conheça o tipo de retorno do SDK saberia que aquele campo existia para ser
		// repassado. Isso exigiria parsear a linguagem e resolver tipos, o que o Anchors não
		// faz por desenho (STRUCTURE/D2: lê texto, nunca parseia código).
		//
		// O que SOBRA e funciona: o gate pega este caso pela regra do limite — reprova
		// `cognitoListDevices` por truncar em `Limit` sem devolver cursor. A mensagem cita o
		// limite em vez do cursor, o que é menos preciso, mas aponta a mesma função e o
		// conserto (devolver o cursor) está entre as saídas que ela oferece.
		// limite que o CHAMADOR escolhe é decisão dele; default escondido não é.
		if lim, temDefault := hiddenLimit(f.body, n.ID, f.params, d); lim != "" {
			if !temDefault {
				continue // o chamador passa o limite: página deliberada, sem promessa quebrada
			}
			achados = append(achados, i18n.T("gate.pagination_honored.item_hidden_default", f.name, lim))
			continue
		}
		achados = append(achados, i18n.T("gate.pagination_honored.item_no_pagination", f.name))
	}

	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	var sisterSuffix string
	if sistersPaginate > 0 {
		var countStr string
		if sistersPaginate == 1 {
			countStr = i18n.T("gate.pagination_honored.sister_count_singular")
		} else {
			countStr = i18n.T("gate.pagination_honored.sister_count_plural", sistersPaginate)
		}
		sisterSuffix = i18n.T("gate.pagination_honored.sister_functions_paginate", countStr)
	}
	return Fail, i18n.T("gate.pagination_honored.fail_broken_promise", strings.Join(achados, "; "), sisterSuffix)
}

// promisesSet: o NOME diz que devolve tudo? É o único lugar onde a promessa está
// escrita — o tipo de retorno é idêntico para uma página e para o conjunto
// (`Promise<T[]>`, `List<T>`, `[]T`… em toda linguagem).
//
// O padrão vem do dialeto, mas aqui o default do Anchors é forte: prefixos como `list`/
// `getAll`/`find_all` são convenção de NOMENCLATURA difundida, não sintaxe de linguagem.
// Um projeto que nomeia noutro idioma sobrepõe com `dialect.set_promise`.
func promisesSet(name string, d config.Dialect) bool {
	promessa := d.Compile(d.SetPromise)
	if promessa == nil || !promessa.MatchString(name) {
		return false
	}
	// `listFirst…` / `listRecent…` / `listTop…` prometem um RECORTE, não o conjunto.
	// Quem nomeou assim já disse que não devolve tudo.
	if recorte := d.Compile(d.SetSlice); recorte != nil && recorte.MatchString(name) {
		return false
	}
	return true
}

// paginatesAll: o corpo drena o cursor até o fim? O sinal é o cursor do provedor
// realimentando um laço — em TS `do { … } while (lastKey)`, em Python
// `while next_token:`, em Go `for { … if tok == "" { break } }`. A FORMA é do dialeto; o
// que se procura (cursor realimentando laço) é universal.
func paginatesAll(body string, d config.Dialect) bool {
	cursor, laço := d.Compile(d.Cursor), d.Compile(d.Loop)
	if cursor == nil || laço == nil {
		return false
	}
	// Cursor sem laço não drena nada: ou é repassado ao chamador (que aí é o responsável),
	// ou está lá sem uso. Cursor DENTRO de laço é a prova de que drena.
	return cursor.MatchString(body) && laço.MatchString(body)
}

// hiddenLimit devolve o limite e se ele é ESCONDIDO do chamador.
//
// Esta é a distinção central do gate, e ela é conceitual (não de linguagem): um limite que
// o chamador PASSA é contrato — ele sabe que existe mais. Um limite que mora no código
// (default na assinatura, literal na query, constante do módulo) é truncamento que o
// chamador não tem como saber que existe.
//
// COMO um limite se escreve é do projeto (`dialect.limits`, DESIGN-core-without-syntax.md):
// cada padrão captura o nome do limite (`name`) e o seu valor — um literal (`value`): está
// no código, escondido; um identificador (`param`): contrato quando é parâmetro da função,
// escondido quando não é (uma constante do módulo). Um padrão com o grupo `type` lê uma
// anotação de tipo (`limit: number`): não é a passagem do valor, e o que outro padrão leia
// no mesmo lugar é ignorado — lê-la como valor fazia o gate concluir "escondido" justamente
// no caso em que o chamador escolheu.
func hiddenLimit(body, rel string, params []string, d config.Dialect) (string, bool) {
	matches := d.Limits.FindAll(body, rel)
	typed := map[int]bool{}
	for _, m := range matches {
		if c, ok := m.Groups["type"]; ok {
			for i := m.Start; i < c.End; i++ {
				typed[i] = true
			}
		}
	}
	name := func(m config.Match) string {
		if c, ok := m.Groups["name"]; ok {
			return c.Value
		}
		return strings.TrimSpace(body[m.Start:m.End])
	}
	// Na ordem do código: um default ou um literal antes de tudo esconde; um parâmetro do
	// chamador é contrato; um identificador que não é parâmetro vem de fora — escondido.
	for _, m := range matches {
		if _, isType := m.Groups["type"]; isType {
			continue
		}
		if v, ok := m.Groups["value"]; ok && !typed[v.Start] {
			return name(m) + " = " + v.Value, true
		}
	}
	for _, m := range matches {
		p, ok := m.Groups["param"]
		if !ok || typed[p.Start] {
			continue
		}
		for _, x := range params {
			if x == p.Value {
				return name(m), false
			}
		}
		return name(m) + ": " + p.Value, true
	}
	return "", false
}

// dispensaExplicita: o opt-out honesto. Exige razão escrita depois dos dois-pontos —
// um marcador nu não dispensa nada, só esconde melhor.
// `[^\S\n]*` = espaço/tab mas NÃO quebra de linha: sem isso o `\s*` atravessaria o
// fim da linha e a razão seria "achada" na linha de código seguinte — o marcador nu
// passaria, que é exatamente o que a dispensa não pode permitir.
var noPaginateRE = regexp.MustCompile(`@no-paginate[^\S\n]*:[^\S\n]*\S+`)

func explicitWaiver(body string) bool { return noPaginateRE.MatchString(body) }

// ownText is the text that is the i-th function's own: the comments right above its
// declaration, and its body up to the comments of the next one. A function's body runs to
// the next declaration, so the mark above that one's signature was read as the previous
// function's, which passed by it, while the function it was written for stayed accused.
func ownText(content, code string, fns []exportedFunc, i int) string {
	end := len(content)
	if i+1 < len(fns) {
		end = leadStart(content, code, fns[i+1].decl)
	}
	return content[leadStart(content, code, fns[i].decl):end]
}

// leadStart is where the comment lines right above the declaration at `at` begin: lines with
// text in the content that are blank in the code. A line of code or an empty line ends them.
func leadStart(content, code string, at int) int {
	start := strings.LastIndexByte(content[:at], '\n') + 1
	for start > 0 {
		prev := strings.LastIndexByte(content[:start-1], '\n') + 1
		if strings.TrimSpace(content[prev:start-1]) == "" || strings.TrimSpace(code[prev:start-1]) != "" {
			break
		}
		start = prev
	}
	return start
}

// drainingHelpers are the names the module defines (`dialect.definition`) whose code drains a
// cursor: a function that walks every page and returns them all. A definition's code runs to
// the next one indented no deeper, so what it defines inside stays its own.
func drainingHelpers(code string, d config.Dialect) []string {
	defRE := d.Compile(d.Definition)
	if defRE == nil {
		return nil
	}
	type def struct {
		name          string
		start, indent int
	}
	var defs []def
	for _, loc := range defRE.FindAllStringSubmatchIndex(code, -1) {
		for g := 1; g*2+1 < len(loc); g++ {
			if loc[g*2] < 0 || defRE.SubexpNames()[g] == "owner" {
				continue
			}
			line := strings.LastIndexByte(code[:loc[g*2]], '\n') + 1
			indent := len(code[line:]) - len(strings.TrimLeft(code[line:], " \t"))
			defs = append(defs, def{code[loc[g*2]:loc[g*2+1]], line, indent})
			break
		}
	}
	var out []string
	for i, df := range defs {
		end := len(code)
		for _, nx := range defs[i+1:] {
			if nx.start > df.start && nx.indent <= df.indent {
				end = nx.start
				break
			}
		}
		if paginatesAll(code[df.start:end], d) {
			out = append(out, df.name)
		}
	}
	return out
}

// callsAny says whether the body calls one of the names, other than the function's own.
func callsAny(body, self string, names []string) bool {
	for _, name := range names {
		if name != self && regexp.MustCompile(`\b`+regexp.QuoteMeta(name)+`\s*(?:<[^()]*>)?\(`).MatchString(body) {
			return true
		}
	}
	return false
}

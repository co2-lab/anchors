package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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
// também dispensa — opt-out explícito, como manda o CONCEPT §5.1.
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
		return Skip, "não é código — a promessa se lê na assinatura da função"
	}
	d := cfg.DialectFor()
	// Cada padrão ausente é um olho fechado. Pendente nomeia QUAL falta, para a correção
	// ser uma linha de YAML e não uma investigação.
	var faltando []string
	if d.ExportedFunc == "" && !d.Dispensado("exported_func") {
		faltando = append(faltando, "`exported_func` (como se reconhece uma função exportada)")
	}
	if d.CollectionQuery == "" && !d.Dispensado("collection_query") {
		faltando = append(faltando, "`collection_query` (como se reconhece uma consulta que "+
			"devolve muitos registros — não tem default porque depende do seu provedor de dados: "+
			"`QueryCommand|ScanCommand` no DynamoDB, `SELECT` em SQL, `.find()` no Mongo)")
	}
	if len(faltando) > 0 {
		// Pending, não Skip: o gate NÃO verificou, e dizer isso é o ponto. Silenciar aqui
		// faria o relatório afirmar conformidade sobre o que ninguém mediu.
		//
		// Quem decidiu não configurar tem o opt-out explícito, e a mensagem o oferece —
		// senão a única saída visível seria conviver com o aviso para sempre, e aviso
		// permanente é aviso que se aprende a ignorar (levando os outros com ele).
		return Pending, "o projeto não declarou " + strings.Join(faltando, " nem ") +
			". Declare em `dialect:` no anchors.yaml (`family:` resolve parte disso: " +
			strings.Join(config.KnownDialectFamilies(), ", ") + "). Se não se aplica a este " +
			"projeto, dispense com `dialect.opt_out: [<campo>]` — o aviso sai do relatório e " +
			"a decisão fica escrita."
	}
	// Campo dispensado = o projeto afirmou que não tem. Sem o padrão não há como confrontar,
	// e o Skip aqui é honesto: diferente do Pending acima, alguém decidiu.
	if d.CollectionQuery == "" || d.ExportedFunc == "" {
		return Skip, "dialeto dispensado por `opt_out` — o projeto declarou que este padrão não se aplica"
	}

	consulta := d.Compile(d.CollectionQuery)
	if consulta == nil || !consulta.MatchString(content) {
		return Skip, "nenhuma consulta de coleção neste arquivo (pelo `dialect.collection_query` do projeto)"
	}

	fns := exportedFuncs(content, d)
	if len(fns) == 0 {
		return Skip, "nenhuma função exportada reconhecida pelo dialeto do projeto"
	}

	// O módulo conhece o padrão de paginação? Se alguma irmã pagina, a que não pagina é
	// assimetria — o mesmo raciocínio do gate sibling-guard.
	irmãsPaginam := 0
	for _, f := range fns {
		if paginaTudo(f.body, d) {
			irmãsPaginam++
		}
	}

	var achados []string
	for _, f := range fns {
		if !prometeConjunto(f.name, d) || !consulta.MatchString(f.body) {
			continue
		}
		if paginaTudo(f.body, d) || dispensaExplicita(f.body) {
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
		if lim, temDefault := limiteEscondido(f.body, f.params); lim != "" {
			if !temDefault {
				continue // o chamador passa o limite: página deliberada, sem promessa quebrada
			}
			achados = append(achados, fmt.Sprintf(
				"`%s` promete o conjunto mas trunca em `%s` por default — o chamador não "+
					"escolheu esse limite e não recebe cursor, então o que passa dele é perdido "+
					"em silêncio", f.name, lim))
			continue
		}
		achados = append(achados, fmt.Sprintf(
			"`%s` promete o conjunto e consulta sem paginar — a consulta trunca na página "+
				"do provedor quando o dado crescer, devolvendo um resultado parcial que "+
				"parece completo", f.name))
	}

	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	msg := "promessa de conjunto não honrada: " + strings.Join(achados, "; ")
	if irmãsPaginam > 0 {
		msg += fmt.Sprintf(". %s neste módulo pagina(m) com loop de cursor — o padrão é "+
			"conhecido aqui, o que torna a omissão esquecimento e não decisão", plural(
			irmãsPaginam, "1 função", fmt.Sprintf("%d funções", irmãsPaginam)))
	}
	msg += ". Para dispensar: devolva o cursor, nomeie o limite como página " +
		"(`pageSize`), renomeie a função para o que ela faz, ou marque `// @no-paginate: <razão>`"
	return Fail, msg
}

// prometeConjunto: o NOME diz que devolve tudo? É o único lugar onde a promessa está
// escrita — o tipo de retorno é idêntico para uma página e para o conjunto
// (`Promise<T[]>`, `List<T>`, `[]T`… em toda linguagem).
//
// O padrão vem do dialeto, mas aqui o default do Anchors é forte: prefixos como `list`/
// `getAll`/`find_all` são convenção de NOMENCLATURA difundida, não sintaxe de linguagem.
// Um projeto que nomeia noutro idioma sobrepõe com `dialect.set_promise`.
func prometeConjunto(name string, d config.Dialect) bool {
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

// paginaTudo: o corpo drena o cursor até o fim? O sinal é o cursor do provedor
// realimentando um laço — em TS `do { … } while (lastKey)`, em Python
// `while next_token:`, em Go `for { … if tok == "" { break } }`. A FORMA é do dialeto; o
// que se procura (cursor realimentando laço) é universal.
func paginaTudo(body string, d config.Dialect) bool {
	cursor, laço := d.Compile(d.Cursor), d.Compile(d.Loop)
	if cursor == nil || laço == nil {
		return false
	}
	// Cursor sem laço não drena nada: ou é repassado ao chamador (que aí é o responsável),
	// ou está lá sem uso. Cursor DENTRO de laço é a prova de que drena.
	return cursor.MatchString(body) && laço.MatchString(body)
}

// limiteEscondido devolve o nome do limite e se ele é ESCONDIDO do chamador.
//
// Esta é a distinção central do gate, e ela é conceitual (não de linguagem): um limite que
// o chamador PASSA é contrato — ele sabe que existe mais. Um limite que mora no código
// (default na assinatura, literal na query, constante do módulo) é truncamento que o
// chamador não tem como saber que existe.
//
// Os nomes procurados (`limit`, `max`, `pageSize`, `take`, `top`) são vocabulário de
// API de dados, difundido em toda stack — a mesma classe de universalidade de `list`/
// `getAll`. A ASSOCIAÇÃO nome↔valor é que varia por sintaxe, e vem em duas formas que
// cobrem as linguagens correntes: `nome = valor` (default de parâmetro) e `nome: valor`
// ou `nome=valor` (campo de objeto / argumento nomeado).
var limiteNome = `(?i)\b(limit|max|maxItems|maxKeys|maxResults|pageSize|page_size|take|top|first)\b`

// `[^=\n]*?` deixa passar uma anotação de tipo (`limit: number = 100`, `limit: int = 100`)
// sem atravessar a linha. Sem esse cuidado, `limit: number` — parâmetro OBRIGATÓRIO, que o
// chamador escolhe — é lido como default e vira falso positivo (aconteceu, e o teste pegou).
var limiteDefaultRE = regexp.MustCompile(limiteNome + `[^=\n]*?=\s*(\d+)`)
var limiteLiteralRE = regexp.MustCompile(limiteNome + `\s*[:=]\s*(\d+)`)

// O valor tem de ser um IDENTIFICADOR que o chamador possa ter passado — não um nome de
// TIPO. `limit: number` (TS), `limit: int` (Python) e `limit: Int` (Kotlin) são anotações
// de tipo na assinatura, não a passagem do limite para a consulta; lê-las como valor faz o
// gate concluir "escondido" justamente no caso em que o chamador escolheu. O teste pegou.
var tiposComuns = regexp.MustCompile(`(?i)^(number|int|integer|long|short|byte|float|double|decimal|uint\d*|int\d*|usize|i\d+|u\d+|Int|Long|Integer|Number|BigInt|size_t)$`)

var limiteParamRE = regexp.MustCompile(limiteNome + `\s*[:=]\s*(\w+)`)

func limiteEscondido(body string, params []string) (string, bool) {
	// default na assinatura: `limit = 100` — o chamador pode omitir e nem saber do corte.
	if m := limiteDefaultRE.FindStringSubmatch(body); m != nil {
		return m[1] + " = " + m[2], true
	}
	// literal cravado na query: `Limit: 50` — ninguém escolheu, está no código.
	if m := limiteLiteralRE.FindStringSubmatch(body); m != nil {
		return m[1] + ": " + m[2], true
	}
	// limite que vem de um parâmetro SEM default: o chamador escolheu, é contrato dele.
	// Varre TODAS as ocorrências: a primeira pode ser a anotação de tipo na assinatura
	// (`limit: number`) e a que interessa vir depois, na consulta (`Limit: limit`).
	for _, m := range limiteParamRE.FindAllStringSubmatch(body, -1) {
		if tiposComuns.MatchString(m[2]) {
			continue // anotação de tipo, não a passagem do valor
		}
		for _, p := range params {
			if p == m[2] {
				return m[1], false
			}
		}
		// vem de uma constante do módulo, não do chamador → escondido.
		return m[1] + ": " + m[2], true
	}
	return "", false
}

// dispensaExplicita: o opt-out honesto. Exige razão escrita depois dos dois-pontos —
// um marcador nu não dispensa nada, só esconde melhor.
// `[^\S\n]*` = espaço/tab mas NÃO quebra de linha: sem isso o `\s*` atravessaria o
// fim da linha e a razão seria "achada" na linha de código seguinte — o marcador nu
// passaria, que é exatamente o que a dispensa não pode permitir.
var noPaginateRE = regexp.MustCompile(`@no-paginate[^\S\n]*:[^\S\n]*\S+`)

func dispensaExplicita(body string) bool { return noPaginateRE.MatchString(body) }

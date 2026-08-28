package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// contract-status-declared: a tabela "Contrato de Saída" de uma spec de INTERFACE
// (handler, rota) tem de listar os status que o código realmente devolve — e só eles.
//
// POR QUE ESTE GATE EXISTE. Numa auditoria de 51 divergências spec×código no app de referência
// (2026-08), este foi o padrão MAIS REPETIDO: oito handlers declaravam um contrato
// que o código não cumpria. E o status omitido era quase sempre o de SEGURANÇA —
// 403 de ownership, 409 de conflito — porque quem escreve a tabela pensa no caminho
// felizardo e nos erros "de negócio", não nas recusas de acesso.
//
// Os dois lados do erro importam, e são diferentes:
//
//   - status EMITIDO e não declarado: o cliente programado pela tabela não trata a
//     recusa. O usuário vê erro genérico onde havia uma razão específica. Foi o caso
//     de `accept-org-invite`, que declarava 3 status e emitia 8 — os dois 403 e o 409
//     ausentes eram a defesa contra sequestro de convite.
//
//   - status DECLARADO e nunca emitido: pior de outro jeito — é código morto no
//     cliente, e some sem ninguém notar. `reanalyse-metadata` declarava 402 para cota
//     excedida; nenhum caminho do handler emite 402 (a cota responde 429). Um cliente
//     que tratasse 402 como "precisa pagar" nunca dispararia esse ramo.
//
// O gate NÃO exige as faixas genéricas (`5xx`) nem inventa semântica: só compara os
// números concretos que aparecem na tabela com os que aparecem no código.
// contratoSecaoRE isola a seção "Contrato de Saída" até o próximo cabeçalho `##`.
// O nome da seção varia um pouco entre specs (`## Contrato de Saída`, ou com sufixo
// como "(por action + erro→status)"), então casa o prefixo.
var contratoSecaoRE = regexp.MustCompile(`(?s)##\s*Contrato de Sa[íi]da[^\n]*\n(.*?)(?:\n##|\z)`)

// statusNaTabelaRE casa o número de status numa linha de tabela: `| 200 | …`.
// Aceita `4xx`/`5xx` na captura para poder IGNORÁ-los depois (são faixas, não
// status concretos — cobrar `500` porque a spec diz `5xx` seria falso-positivo).
var statusNaTabelaRE = regexp.MustCompile(`(?m)^\s*\|\s*\*{0,2}(\d{3}|\d[xX]{2})\b`)

// statusHTTPPadraoTS é o léxico de FALLBACK, usado só quando o projeto não declara
// `dialect:`. Não é o léxico do gate: o gate lê `Dialect.HTTPStatus`, e é isso que o
// torna agnóstico de linguagem. Este default existe para não exigir configuração de
// um projeto TS/Lambda — a stack em que o gate nasceu.
const statusHTTPPadraoTS = `statusCode:\s*(\d{3})|\b[A-Za-z_$][\w$]*\(\s*(\d{3})\s*[,)]`
const statusDinamicoPadraoTS = `statusCode:\s*[A-Za-z_$]`

// statusHTTPNomeados traduz a constante nomeada em número, para os léxicos que a
// usam (`http.StatusForbidden` em Go, `:forbidden` em Rails). Só os status que
// aparecem em contrato de API — não é a tabela inteira da RFC.
var statusHTTPNomeados = map[string]string{
	"OK": "200", "Created": "201", "Accepted": "202", "NoContent": "204",
	"BadRequest": "400", "Unauthorized": "401", "PaymentRequired": "402",
	"Forbidden": "403", "NotFound": "404", "MethodNotAllowed": "405",
	"Conflict": "409", "Gone": "410", "UnprocessableEntity": "422",
	"TooManyRequests": "429", "InternalServerError": "500", "BadGateway": "502",
	"ServiceUnavailable": "503",
}

// statusExcecaoRE: status que o gate não cobra por AUSÊNCIA na tabela, porque vêm
// de infraestrutura compartilhada e não de uma decisão do handler.
//   - 500: quase todo handler tem o try/catch do topo, e as specs declaram `5xx`.
//   - 401: quando o authorizer do API Gateway barra antes do handler.
var statusExcecao = map[string]bool{"500": true}

func checkContractStatusDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "não é uma spec — só spec tem Contrato de Saída"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	m := contratoSecaoRE.FindStringSubmatch(content)
	if m == nil {
		return Skip, "spec sem seção `## Contrato de Saída` — nada a confrontar"
	}
	secao := m[1]

	// Declarados: só os concretos. `5xx`/`4xx` são faixas — a spec usa para "qualquer
	// falha", e tratá-los como status exigiria adivinhar quais números cobrem.
	declarados := map[string]bool{}
	temFaixa := false
	for _, mm := range statusNaTabelaRE.FindAllStringSubmatch(secao, -1) {
		s := mm[1]
		if strings.ContainsAny(s, "xX") {
			temFaixa = true
			continue
		}
		declarados[s] = true
	}

	// O código que esta spec `specifies`.
	var codePaths []string
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type == mapx.EdgeSpecifies {
			codePaths = append(codePaths, e.To)
		}
	}
	if len(codePaths) == 0 {
		return Pending, "spec sem código ligado (specifies) — nada a confrontar ainda"
	}

	var corpo strings.Builder
	for _, cp := range codePaths {
		b, err := os.ReadFile(filepath.Join(root, cp))
		if err != nil {
			continue
		}
		// Sem comentários: um `// devolve 404` num comentário não é um status emitido,
		// e é exatamente o tipo de texto que descreve o comportamento ANTIGO.
		corpo.WriteString(stripLineComments(string(b)))
		corpo.WriteString("\n")
	}
	code := corpo.String()
	if strings.TrimSpace(code) == "" {
		return Pending, "código ligado não pôde ser lido — nada a confrontar"
	}

	// O LÉXICO vem do dialeto do projeto — é o que mantém o gate agnóstico. Um
	// projeto Go declara `http_status` casando `WriteHeader(403)`; um Rails, o
	// `render status:`. Sem `dialect:`, cai no default TS/Lambda.
	d := cfg.DialectFor()
	padrao := d.HTTPStatus
	if padrao == "" {
		padrao = statusHTTPPadraoTS
	}
	reStatus := d.Compile(padrao)
	if reStatus == nil {
		return Pending, "o `dialect.http_status` do projeto não compila como regex — " +
			"sem léxico não há como ler os status do código"
	}

	emitidos := map[string]bool{}
	for _, mm := range reStatus.FindAllStringSubmatch(code, -1) {
		// Qualquer grupo que casou serve: o padrão tem uma alternativa por forma
		// (envelope, helper, constante nomeada), e só uma casa por vez.
		for _, g := range mm[1:] {
			if g == "" {
				continue
			}
			if num, ok := statusHTTPNomeados[g]; ok {
				emitidos[num] = true
			} else if len(g) == 3 && g[0] >= '1' && g[0] <= '5' {
				emitidos[g] = true
			}
		}
	}
	if len(emitidos) == 0 {
		return Skip, "o código não devolve status por forma reconhecível " +
			"(`statusCode:`, `jsonResponse(`, `cors(`) — pode não ser uma interface HTTP"
	}

	var faltando, fantasma []string
	for s := range emitidos {
		if declarados[s] || statusExcecao[s] {
			continue
		}
		// Faixa declarada cobre o 5xx correspondente: `5xx` na tabela ⇒ 503 ok.
		if temFaixa && len(s) == 3 && (s[0] == '5' || s[0] == '4') {
			// Só a faixa 5xx é aceita como cobertura genérica: um `4xx` genérico
			// esconderia justamente as recusas de acesso que este gate persegue.
			if s[0] == '5' {
				continue
			}
		}
		faltando = append(faltando, s)
	}
	// O lado "fantasma" só é afirmável quando TODO status do código é literal. Com
	// um helper que recebe o status por parâmetro (`statusCode: status`), um valor
	// declarado pode muito bem ser emitido por uma chamada que a leitura textual não
	// alcança — acusar seria falso-positivo. Este é o caso do `log-consent`, que
	// define `fail(status, message)` e chama `fail(400, …)`: o gate vê o literal na
	// chamada, mas se o helper fosse indireto não veria.
	padraoDin := d.HTTPStatusDynamic
	if padraoDin == "" && d.HTTPStatus == "" {
		padraoDin = statusDinamicoPadraoTS // só com o default TS; dialeto declarado manda
	}
	reDin := d.Compile(padraoDin)
	temDinamico := reDin != nil && reDin.MatchString(code)
	if !temDinamico {
		for s := range declarados {
			if !emitidos[s] {
				fantasma = append(fantasma, s)
			}
		}
	}
	sort.Strings(faltando)
	sort.Strings(fantasma)

	var partes []string
	if len(faltando) > 0 {
		partes = append(partes, fmt.Sprintf(
			"o código devolve %s e o Contrato de Saída não declara. Quem programa pela "+
				"tabela não trata essa resposta — e o status omitido costuma ser o de "+
				"SEGURANÇA (403 de ownership, 409 de conflito), porque quem escreve a "+
				"tabela pensa no caminho felizardo",
			strings.Join(faltando, ", ")))
	}
	if len(fantasma) > 0 {
		partes = append(partes, fmt.Sprintf(
			"o Contrato declara %s e NENHUM caminho do código emite. É código morto no "+
				"cliente: o ramo que trata esse status nunca dispara, e ninguém percebe",
			strings.Join(fantasma, ", ")))
	}
	if len(partes) > 0 {
		return Fail, strings.Join(partes, "; ") + ". Confronte a tabela com os status " +
			"reais do handler — a lista de emitidos sai de um grep dos `statusCode:`"
	}

	return Pass, ""
}

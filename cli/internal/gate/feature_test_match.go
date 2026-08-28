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
	"github.com/co2-lab/anchors/internal/similarity"
)

// checkFeatureTestMatch — gate RELACIONAL da trinca (STRUCTURE/TRACEABILITY): confronta
// uma `.feature` contra o(s) teste(s) que a realizam (arestas `tested-by`). Garante que
// cada CENÁRIO documentado na feature está IMPLEMENTADO no teste — por CÓDIGO e por
// DESCRIÇÃO. Pega o agente que pula um cenário, renomeia ou muda os passos, divergindo
// do que a feature documentou.
//
// Régua (estática, sem execução):
//   - CÓDIGO: todo scenario-code `XXXX-Y##` de um `Cenário:` da feature deve aparecer
//     no conteúdo (não-comentário) de ao menos um teste ligado.
//   - DESCRIÇÃO: o título do `Cenário:` deve corresponder (match TOLERANTE por termos
//     significativos) ao texto do teste que cita aquele código — o desvio de descrição
//     é sinalizado como aviso no detalhe (não derruba se o código casa), porque texto
//     livre varia; a ausência do CÓDIGO é o que reprova.
func checkFeatureTestMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, "" // só confronta features
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}
	// O de-para de regime (tag-do-projeto → regime-canônico) vem da Estrutura. Este gate
	// confronta a superfície `test` (o teste ligado por tested-by): só os cenários cujo
	// regime canônico é confrontado por ela (unit/integration). Cenários de outros regimes
	// (e2e, vr) ou sem regime mapeado são verificados noutra superfície — aqui, pulados.
	testRegimes := regimesForTestSurface(cfg)
	scenarios := parseFeatureScenarios(content)
	if len(scenarios) == 0 {
		return Skip, "" // feature sem cenários codificados — nada a confrontar
	}

	// testes ligados (tested-by saindo da feature)
	var testPaths []string
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type == mapx.EdgeTestedBy {
			testPaths = append(testPaths, e.To)
		}
	}
	if len(testPaths) == 0 {
		// sem teste ligado: se a feature declara cenários, isto é uma lacuna da trinca —
		// mas quem cobra a EXISTÊNCIA do teste é a co-location/tested-by; aqui só
		// confrontamos correspondência quando há teste. Pending para não duplicar.
		return Pending, "nenhum teste ligado (tested-by) — nada a confrontar ainda"
	}

	// une o conteúdo (não-comentário) de todos os testes ligados
	var testBody, bodyComComentarios strings.Builder
	for _, tp := range testPaths {
		b, err := os.ReadFile(filepath.Join(root, tp))
		if err != nil {
			continue
		}
		testBody.WriteString(stripLineComments(string(b)))
		testBody.WriteString("\n")
		bodyComComentarios.Write(b)
		bodyComComentarios.WriteString("\n")
	}
	body := testBody.String()
	// Duas leituras do mesmo teste, para duas perguntas diferentes.
	//
	// `body` (sem comentários) responde "o código está IMPLEMENTADO aqui?" — e um
	// código citado em comentário não é implementação, é referência.
	//
	// `bodyNorm` (COM comentários) responde "a descrição do cenário está coberta?".
	// Aqui o comentário conta, e tem de contar: o cenário fala o vocabulário do
	// domínio ("recolore a linha", "mascaram os montantes", "não navega") e o código
	// fala identificadores em inglês (`fill`, `••••`, ausência de `useNavigation`).
	// A ponte entre os dois é exatamente o comentário que o autor escreveu acima do
	// `it`. Sem ele nesta régua, o gate cobrava do teste uma palavra portuguesa que
	// nenhum TypeScript vai conter — e a prova, que existe e é específica, contava
	// como ausente.
	bodyNorm := normalizeText(bodyComComentarios.String())

	var missingCode []string
	var driftDesc []string
	// O corpus é o conjunto de títulos DESTA feature: é ele que define o que é
	// palavra comum aqui dentro. "Componente" pesa 0 num arquivo em que todos os
	// cenários começam assim, e muito num em que só um o menciona.
	corpus := make([]string, 0, len(scenarios))
	for _, sc := range scenarios {
		corpus = append(corpus, sc.Title)
	}
	pesos := similarity.Weights(corpus)

	for _, sc := range scenarios {
		if !scenarioHitsTestSurface(sc, testRegimes) {
			continue // cenário de outro regime (e2e/vr) ou sem regime → outra superfície
		}
		if !strings.Contains(body, sc.Code) {
			missingCode = append(missingCode, sc.Code)
			continue
		}
		// Código presente — agora a descrição. A régua é IGUALDADE: o título do
		// cenário e o texto do teste têm de dizer a mesma coisa, não "60% da mesma
		// coisa". Um limiar fracionário deixa passar o teste que descreve outro
		// caso e ainda cita o código certo.
		//
		// A SIMILARIDADE não relaxa a régua — ela classifica o achado, para quem
		// vai consertar saber o que fazer:
		//   similar    → mesmo assunto, palavras diferentes: reescreva um dos lados.
		//   divergente → assuntos diferentes: decida qual dos dois está velho.
		titulo, temTitulo := testTitleFor(body, sc.Code)
		switch {
		case temTitulo && !tituloCompartilhado(body, sc.Code):
			if v, score := similarity.Classifica(sc.Title, titulo, pesos); v != similarity.Identico {
				driftDesc = append(driftDesc, fmt.Sprintf("%s (%s, %.0f%%)", sc.Code, v, score*100))
			}
		case !descriptionMatches(sc.Title, bodyNorm):
			// Duas situações caem aqui, e nas duas não existe UM título para comparar
			// um-a-um:
			//
			//   - o código aparece só em comentário, ou em `describe`;
			//   - o título é COMPARTILHADO por vários cenários
			//     (`it('A / B / C: sem iniciais exibe o fallback')`).
			//
			// O segundo caso é o que forçava falso positivo: com N códigos num
			// título, a régua de igualdade compara os N cenários com o MESMO texto, e
			// no máximo um pode ser idêntico — os outros N-1 divergiam sempre, por
			// construção. E o teste conjunto é legítimo: três cenários que descrevem a
			// mesma situação por eixos diferentes do vocabulário (estado,
			// comportamento, mensagem) têm uma prova só.
			//
			// Aqui a pergunta certa é a da régua de corpo: o miolo do cenário está no
			// teste? Se está, o cenário é provado por ele — mesmo que o título diga
			// respeito a outro dos irmãos.
			driftDesc = append(driftDesc, sc.Code)
		}
	}

	if len(missingCode) > 0 {
		sort.Strings(missingCode)
		detail := fmt.Sprintf("%d cenário(s) da feature SEM implementação no teste (código ausente): %s",
			len(missingCode), strings.Join(missingCode, ", "))
		return Fail, detail
	}
	if len(driftDesc) > 0 {
		sort.Strings(driftDesc)
		// descrição divergente é AVISO (Pending), não Fail: o código casa (rastreável),
		// mas o texto do teste não reflete o cenário — o autor pode ter mudado os passos.
		return Pending, fmt.Sprintf("código casa, mas a descrição de %d cenário(s) diverge do teste: %s",
			len(driftDesc), strings.Join(driftDesc, ", "))
	}
	return Pass, ""
}

type featureScenario struct {
	Code  string
	Title string
	Tags  []string // as tags de regime do projeto na tag-line (sem @), ex.: ["nivel-unit","smoke"]
	// Codes: TODOS os códigos de cenário da tag-line, não só o primeiro. Um cenário pode
	// provar mais de um requisito, e o dialeto do projeto co-etiqueta:
	// `@acao @TCDTX-A01 @TCDTX-M03 @nivel-integration`.
	//
	// Enquanto só o primeiro era lido, `spec-feature-match` acusava o segundo como "sem
	// cenário" — medido numa spec real, 2 requisitos denunciados que estavam escritos,
	// testados e passando. O gate mandava escrever um cenário que já existia, e o autor
	// obediente duplicaria o cenário para calar o gate.
	Codes []string
}

// featScenarioCodeRE é RECONFIGURADO por SetRuleLetters na carga da config — a letra do
// código vem do vocabulário do PROJETO (`rule_types`), não de uma lista fixa. Cravar as
// canônicas aqui tornava INVISÍVEL todo cenário de uma letra declarada pelo projeto: o
// gate não via o cenário e reportava verde sobre o que não conferiu — o modo de falha
// mais perigoso que existe. Aconteceu com a letra `I` (Invariant), adicionada por um
// projeto real, e nenhum teste pegou.
var (
	featScenarioCodeRE = featCodeREFor(config.DefaultRuleLetters)
	// Reconhece a abertura de cenário em QUALQUER idioma do Gherkin — inclusive as formas
	// de Esquema/Outline, que em vários idiomas não são "<Cenário> + sufixo" e sim uma
	// expressão própria (`Esquema do Cenário`, `Plan du scénario`, `Szenariogrundriss`).
	// O regex anterior só casava `Cenário|Scenario` com um ` Outline` opcional em inglês:
	// 108 `Esquema do Cenário` de um projeto real eram INVISÍVEIS para este gate, que é
	// BLOQUEANTE — reportava verde sobre o que não enxergava.
	featTitleRE = regexp.MustCompile(`(?m)^\s*(?:` +
		strings.Join(quoteAll(config.GherkinScenarioAlternatives()), "|") +
		`)(?:\s+Outline)?:\s*(.+?)\s*$`)
	featTagRE = regexp.MustCompile(`@([a-z0-9][a-z0-9-]*)`)
)

// featCodeREFor: o código do cenário, com o SUFIXO de cenário opcional (`#NN`).
//
// O sufixo existe porque uma regra legitimamente tem mais de um cenário — caminho
// feliz e alternativos. Sem ele, os N cenários de `USBPX-B01` carregam o mesmo
// código, e nada distingue um do outro: o gate compara os N títulos com o mesmo
// teste e no máximo um casa. Com `#01`/`#02`, o código passa a identificar o
// CENÁRIO (a regra continua legível no prefixo), e o par cenário↔teste volta a ser
// um-para-um.
//
// Retrocompatível: o sufixo é opcional, e código sem ele continua casando.
func featCodeREFor(letters string) *regexp.Regexp {
	return regexp.MustCompile(`@([A-Z0-9]` + config.CodeLengthPattern() + `-(?:[` + regexp.QuoteMeta(letters) + `]\d{2}|DS-[A-Za-z0-9-]+|VR))(#\d{2})?\b`)
}

// codeRaizRE separa a raiz (`USBPX-B01`) do sufixo de cenário (`#02`). Os gates que
// falam de REGRA usam a raiz; os que falam de CENÁRIO usam o código inteiro.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func codeRaizRE() *regexp.Regexp {
	return regexp.MustCompile(`^([A-Z0-9]` + config.CodeLengthPattern() + `-[A-Za-z0-9-]+?)(#\d{2})?$`)
}

// CodeRaiz devolve o código sem o sufixo de cenário.
func CodeRaiz(code string) string {
	if m := codeRaizRE().FindStringSubmatch(code); m != nil {
		return m[1]
	}
	return code
}

// parseFeatureScenarios lê os blocos de cenário: a linha de tags (@CODE + @tags-de-regime)
// seguida da linha `Cenário:`. Casa o código ao título e captura as TAGS da tag-line
// (nomes de regime do projeto). O roteamento por regime é decidido depois, com o de-para
// da config (STRUCTURE §2.3) — este parser não conhece a nomenclatura local.
func parseFeatureScenarios(content string) []featureScenario {
	lines := strings.Split(content, "\n")
	var out []featureScenario
	var pendingCode string
	var pendingCodes []string
	var pendingTags []string
	flush := func(title string) {
		if pendingCode == "" {
			return
		}
		out = append(out, featureScenario{Code: pendingCode, Title: title, Tags: pendingTags, Codes: pendingCodes})
		pendingCode, pendingTags, pendingCodes = "", nil, nil
	}
	for _, ln := range lines {
		if ms := featScenarioCodeRE.FindAllStringSubmatch(ln, -1); ms != nil {
			// O PRIMEIRO código é a identidade do cenário (o que o mapa e os demais gates
			// usam); os seguintes são requisitos que o mesmo cenário também prova.
			// O código do cenário inclui o sufixo `#NN` quando presente: é ele que
			// dá identidade a cada cenário de uma regra com vários.
			pendingCode = ms[0][1] + ms[0][2]
			pendingCodes = nil
			for _, m := range ms {
				pendingCodes = append(pendingCodes, m[1]+m[2])
			}
			pendingTags = nil // tags são as da tag-line DESTE código — não vazar do cenário anterior
			for _, tm := range featTagRE.FindAllStringSubmatch(ln, -1) {
				pendingTags = append(pendingTags, tm[1])
			}
		}
		if m := featTitleRE.FindStringSubmatch(ln); m != nil {
			flush(m[1])
		}
	}
	return out
}

// canonicalTestSurfaceRegimes: os regimes canônicos que a superfície `test` confronta.
// (unit e integration moram no arquivo de teste; e2e/vr moram em outras superfícies.)
var canonicalTestSurfaceRegimes = map[string]bool{"unit": true, "integration": true}

// regimesForTestSurface devolve o conjunto de TAGS-do-projeto cujo regime canônico é
// confrontado pela superfície `test`, a partir do de-para da config (derived.regimes +
// derived.surfaces). Se a config não declara de-para, cai num default sensato que
// reconhece as próprias tags canônicas (unit/integration) — assim um projeto do zero
// funciona sem de-para.
func regimesForTestSurface(cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	// quais regimes canônicos a superfície "test" cobre? (via derived.surfaces, se houver)
	testRegimes := map[string]bool{}
	if cfg != nil && cfg.Derived != nil && len(cfg.Derived.Surfaces) > 0 {
		for regime, surface := range cfg.Derived.Surfaces {
			if surface == "test" {
				testRegimes[regime] = true
			}
		}
	} else {
		testRegimes = canonicalTestSurfaceRegimes // default: unit+integration → test
	}
	// traduz os regimes-test de volta para as TAGS do projeto (de-para reverso)
	if cfg != nil && cfg.Derived != nil && len(cfg.Derived.Regimes) > 0 {
		for tag, regime := range cfg.Derived.Regimes {
			if testRegimes[regime] {
				out[tag] = true
			}
		}
	} else {
		// sem de-para: a própria tag é o regime canônico
		for regime := range testRegimes {
			out[regime] = true
		}
	}
	return out
}

// scenarioHitsTestSurface: o cenário é confrontado pela superfície `test`? Verdadeiro se
// ALGUMA de suas tags mapeia para um regime-test. Cenário SEM nenhuma tag de regime
// declarada é confrontado (default conservador — cobra presença). Cenário cujas tags são
// todas de outros regimes (e2e/vr) NÃO é confrontado aqui (mora noutra superfície).
func scenarioHitsTestSurface(sc featureScenario, testRegimes map[string]bool) bool {
	sawRegimeTag := false
	for _, t := range sc.Tags {
		if testRegimes[t] {
			return true
		}
		// é uma tag de regime conhecida (mapeada a QUALQUER regime)? marca que houve
		// declaração de regime — para distinguir "sem regime" de "outro regime".
		if isRegimeTag(t) {
			sawRegimeTag = true
		}
	}
	return !sawRegimeTag // sem nenhuma tag de regime → conservador (confronta)
}

// isRegimeTag: heurística barata p/ reconhecer uma tag que É de regime (mesmo que não
// caia na superfície test) — as convenções comuns começam por "nivel-" ou são um regime
// canônico. Usada só para distinguir "cenário sem regime" (confronta) de "cenário de
// outro regime" (pula). Conservadora: na dúvida, NÃO é regime.
func isRegimeTag(t string) bool {
	if strings.HasPrefix(t, "nivel-") {
		return true
	}
	switch t {
	case "unit", "integration", "e2e", "vr":
		return true
	}
	return false
}

// stripLineComments remove comentários de linha `//` e de bloco simples para que os
// códigos citados em comentário do teste NÃO contem como implementação (coerente com
// extractCodes: comentário é referência, não posse). Barato e suficiente aqui.
func stripLineComments(s string) string {
	var b strings.Builder
	for _, ln := range strings.Split(s, "\n") {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "*") || strings.HasPrefix(t, "/*") {
			continue
		}
		// corta comentário inline
		if i := strings.Index(ln, "//"); i >= 0 {
			ln = ln[:i]
		}
		b.WriteString(ln)
		b.WriteString("\n")
	}
	return b.String()
}

var nonWordRE = regexp.MustCompile(`[^\p{L}\p{N}]+`)

// stopwords PT/EN comuns + termos de Gherkin/teste — ruído que não distingue cenários.
var descStopwords = map[string]bool{
	"a": true, "o": true, "os": true, "as": true, "de": true, "da": true, "do": true,
	"e": true, "que": true, "um": true, "uma": true, "no": true, "na": true, "em": true,
	"com": true, "por": true, "para": true, "quando": true, "então": true, "entao": true,
	"dado": true, "deve": true, "the": true, "and": true, "to": true, "of": true, "is": true,
	"when": true, "then": true, "given": true, "should": true, "cenário": true, "cenario": true,
}

func normalizeText(s string) string {
	return " " + strings.ToLower(nonWordRE.ReplaceAllString(s, " ")) + " "
}

// descriptionMatches: ≥60% dos termos SIGNIFICATIVOS do título do cenário aparecem no
// corpo (normalizado) do teste. Tolerante a reformulação — só exige que o miolo esteja
// lá. Título sem termos significativos (raro) passa (não há o que confrontar).
func descriptionMatches(title, bodyNorm string) bool {
	terms := significantTerms(title)
	if len(terms) == 0 {
		return true
	}
	hit := 0
	for _, t := range terms {
		if strings.Contains(bodyNorm, " "+t+" ") || strings.Contains(bodyNorm, " "+t) {
			hit++
		}
	}
	return float64(hit)/float64(len(terms)) >= 0.6
}

func significantTerms(title string) []string {
	var out []string
	seen := map[string]bool{}
	for _, w := range strings.Fields(strings.ToLower(nonWordRE.ReplaceAllString(title, " "))) {
		if len(w) < 3 || descStopwords[w] || seen[w] {
			continue
		}
		seen[w] = true
		out = append(out, w)
	}
	return out
}

// quoteAll escapa cada alternativa para uso literal dentro do regex.
func quoteAll(xs []string) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = regexp.QuoteMeta(x)
	}
	return out
}

// testTitleReCache evita recompilar o regex por cenário (features grandes têm 40+).
var testTitleReCache = map[string]*regexp.Regexp{}

// testTitleFor extrai o título do teste que cita `code`.
//
// Casa as formas usuais — `it('CODE: título')`, `it('[CODE] título')`,
// `it("CODE — título")` — e devolve ok=false quando o código aparece só em
// comentário ou num teste que prova vários cenários de uma vez: nesses casos não
// há UM título para comparar, e forçar a comparação inventaria divergência.
// tituloCompartilhado diz se o `it` que cita `code` cita OUTRO código também.
//
// Um título com vários códigos descreve o conjunto, não cada um: comparar o
// título com cada cenário por igualdade condenaria N-1 deles sempre.
func tituloCompartilhado(body, code string) bool {
	re, ok := tituloIrmaosReCache[code]
	if !ok {
		cod := regexp.QuoteMeta(code)
		re = regexp.MustCompile(
			`(?:it|test)\s*\(\s*['"` + "`" + `][^'"` + "`" + `]*` + cod + `[^'"` + "`" + `]*['"` + "`" + `]`)
		tituloIrmaosReCache[code] = re
	}
	m := re.FindString(body)
	if m == "" {
		return false
	}
	// quantos códigos DISTINTOS o título cita?
	achados := map[string]bool{}
	for _, c := range codigoNoTituloRE.FindAllString(m, -1) {
		achados[c] = true
	}
	return len(achados) > 1
}

var (
	tituloIrmaosReCache = map[string]*regexp.Regexp{}
	codigoNoTituloRE    = regexp.MustCompile(`[A-Z0-9]` + config.CodeLengthPattern() + `-[A-Z]{1,2}\d{2}(?:#\d{2})?`)
)

// testTitleFor extrai o título do `it` que cita `code`.
//
// Um teste pode provar VÁRIOS cenários e citar todos no título — a forma usada no
// projeto é `it('AATAX-S03 / AATAX-B02 / AATAX-M01: sem iniciais exibe o fallback')`,
// e ela é legítima: três cenários que descrevem a mesma situação por eixos
// diferentes do vocabulário (estado, comportamento, mensagem) têm uma prova só.
//
// O título desse teste é o texto DEPOIS de todos os códigos. Antes o prefixo de
// códigos era casado como `\[?CODE\]?`, um código exato, e isso errava nos dois
// sentidos: para o PRIMEIRO código o título capturado vinha com o resto do prefixo
// grudado ("/ AATAX-B02 / AATAX-M01: sem iniciais…"), e para os SEGUINTES nada
// casava — caindo na régua de corpo, que responde outra pergunta. As duas falhas
// viravam divergência de descrição relatada onde os dois lados diziam o mesmo.
func testTitleFor(body, code string) (string, bool) {
	re, ok := testTitleReCache[code]
	if !ok {
		// RE2 não tem backreference, então "fecha com o mesmo delimitador que abriu"
		// não cabe numa expressão só: são três ramos, unidos por alternância.
		//
		// `prefixo` cobre a lista de códigos irmãos que pode vir antes ou depois do
		// nosso, separados por `/`, `,` ou espaço. O corpo capturado começa só
		// depois do último deles.
		// Cada código pode vir entre colchetes por conta própria (`[A], [B]`) ou o
		// colchete pode envolver a lista inteira (`[A / B]`): o `\[?`/`\]?` fica em
		// cada peça, e não numa volta só.
		// O código tem de TERMINAR aqui: sem a fronteira, `ABCDX-DS-delta-up` casava
		// dentro de `ABCDX-DS-delta-up-high` e o gate lia o título do teste vizinho —
		// comparando o cenário de "até 20%" com a prova de "acima de 20%".
		cod := `\[?` + regexp.QuoteMeta(code) + `\]?(?:[^\w#-]|$)`
		// O irmão pode ser numérico (`ABCDX-B01`, `ABCDX-B01#02`) ou NOMINAL
		// (`ABCDX-DS-fatura-marcado`) — o vocabulário aceita as duas formas, e deixar a
		// segunda de fora fazia o título do primeiro código vir com o prefixo do irmão
		// grudado, o que produzia "similar 100%": mesmo texto, comparação diferente.
		outro := `\[?[A-Z0-9]` + config.CodeLengthPattern() + `-(?:[A-Z]{1,2}\d{2}(?:#\d{2})?|DS-[\w-]+)\]?`
		irmaos := `(?:\s*[/,]?\s*` + outro + `)*`
		abre := `(?:it|test)\s*\(\s*`
		meio := `\[?` + irmaos + `\s*[/,]?\s*` + cod + irmaos + `\]?\s*[:—-]?\s*`
		re = regexp.MustCompile(
			abre + `'` + meio + `([^']*)'` + `|` +
				abre + `"` + meio + `([^"]*)"` + `|` +
				abre + "`" + meio + "([^`]*)`",
		)
		testTitleReCache[code] = re
	}
	m := re.FindStringSubmatch(body)
	if m == nil {
		return "", false
	}
	// Só um dos três ramos casou; os outros grupos vêm vazios.
	for _, g := range m[1:] {
		if g != "" {
			return strings.TrimSpace(g), true
		}
	}
	return "", false
}

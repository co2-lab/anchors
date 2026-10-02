package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/similarity"
	"github.com/co2-lab/anchors/internal/testlist"
)

// checkFeatureTestMatch — gate RELACIONAL da unidade (STRUCTURE/TRACEABILITY): confronta
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
		return pendingNoMap()
	}
	// O de-para de regime (tag-do-projeto → regime-canônico) vem da Estrutura. Este gate
	// confronta a superfície `test` (o teste ligado por tested-by): só os cenários cujo
	// regime canônico é confrontado por ela (unit/integration). Cenários de outros regimes
	// (e2e, vr) ou sem regime mapeado são verificados noutra superfície — aqui, pulados.
	testRegimes := regimesForTestSurface(cfg)
	knownRegimes := knownRegimeTags(cfg)
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
	testPaths = realTests(g, testPaths)
	if len(testPaths) == 0 {
		// sem teste ligado: se a feature declara cenários, isto é uma lacuna da unidade —
		// mas quem cobra a EXISTÊNCIA do teste é a co-location/tested-by; aqui só
		// confrontamos correspondência quando há teste. Pending para não duplicar.
		return Pending, i18n.T("gate.feature_test_match.pending_no_linked_tests")
	}

	// une o conteúdo (não-comentário) de todos os testes ligados
	var testBody, bodyComComentarios strings.Builder
	for _, tp := range testPaths {
		b, err := readFile(root, tp)
		if err != nil {
			continue
		}
		testBody.WriteString(stripLineComments(string(b)))
		testBody.WriteString("\n")
		bodyComComentarios.Write(b)
		bodyComComentarios.WriteString("\n")
	}
	body := testBody.String()
	// The titles of the linked tests, read the way the project says its tests are written.
	all, _, err := projectTests(root, g, cfg)
	if err != nil {
		return Fail, i18n.T("gate.tests_source.failed", err)
	}
	mine := testsIn(all, testPaths)
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
		if !scenarioHitsTestSurface(sc, testRegimes, knownRegimes) {
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
		titulo, compartilhado, temTitulo := titleFor(mine, sc.Code)
		// EVERY test that cites the code, not only the first. A test tagged with a code whose
		// scenario it does not exercise passed unseen as long as another test came first: in
		// the reference app, a logo easter egg carried the code of "open the alerts". The
		// others are held to a looser ruler than the first — only a DIVERGENT title is said —,
		// since several tests of one rule legitimately name its variations.
		for _, other := range otherTitlesFor(mine, sc.Code) {
			if titleCovers(sc.Title, other) {
				continue
			}
			if v, score := similarity.Classify(sc.Title, withoutPlaceholders(other), pesos); v == similarity.Divergente {
				driftDesc = append(driftDesc, i18n.T("gate.feature_test_match.another_test_diverges", sc.Code, other, score*100))
			}
		}
		switch {
		case temTitulo && !compartilhado:
			if titleCovers(sc.Title, titulo) {
				break
			}
			if v, score := similarity.Classify(sc.Title, withoutPlaceholders(titulo), pesos); v != similarity.Identico {
				driftDesc = append(driftDesc, fmt.Sprintf("%s (%s, %.0f%%)", sc.Code, verdictLabel(v), score*100))
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
		detail := i18n.T("gate.feature_test_match.fail_missing_implementation",
			len(missingCode), strings.Join(missingCode, ", "))
		return Fail, detail
	}
	if len(driftDesc) > 0 {
		sort.Strings(driftDesc)
		// descrição divergente é AVISO (Pending), não Fail: o código casa (rastreável),
		// mas o texto do teste não reflete o cenário — o autor pode ter mudado os passos.
		return Diverge, i18n.T("gate.feature_test_match.pending_description_diverges",
			len(driftDesc), strings.Join(driftDesc, ", "))
	}
	return Pass, ""
}

// titleCovers says whether a test's title says what the scenario's title says: the same
// words, in order — case and punctuation aside —, optionally followed by detail
// (`<title> — <detail>`, `<title> (<detail>)`, `<title>: <detail>`). A table-driven
// test's placeholders (`%s`, `%d`, `$name`, `${name}`) are not words of the title. The
// ruler stays equality of what the scenario asserts: a test that says it and more says it;
// one that says something else does not.
func titleCovers(scenario, test string) bool {
	want, got := titleWords(withoutPlaceholders(scenario)), titleWords(withoutPlaceholders(test))
	if len(want) == 0 || len(got) < len(want) {
		return false
	}
	for i, w := range want {
		if got[i] != w {
			return false
		}
	}
	return true
}

// tablePlaceholderRE is a placeholder of a parameterised title: printf verbs (`%s`, `%d`,
// `%#`), `$name` / `${name}` interpolations, and an outline's `<name>`.
var tablePlaceholderRE = regexp.MustCompile(`%[#a-zA-Z%]|\$\{[^}]*\}|\$[A-Za-z_][A-Za-z0-9_.]*|<[A-Za-z_][\w -]*>`)

// withoutPlaceholders drops a table-driven test's placeholders from its title.
func withoutPlaceholders(title string) string { return tablePlaceholderRE.ReplaceAllString(title, " ") }

var titleWordRE = regexp.MustCompile(`[\p{L}\p{N}]+`)

// titleWords are a title's runs of letters and digits, in lower case.
func titleWords(s string) []string {
	ws := titleWordRE.FindAllString(strings.ToLower(s), -1)
	return ws
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

// rootCodeRE separa a raiz (`USBPX-B01`) do sufixo de cenário (`#02`). Os gates que
// falam de REGRA usam a raiz; os que falam de CENÁRIO usam o código inteiro.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func rootCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`^([A-Z0-9]` + config.CodeLengthPattern() + `-[A-Za-z0-9-]+?)(#\d{2})?$`)
}

// RootCode devolve o código sem o sufixo de cenário.
func RootCode(code string) string {
	if m := rootCodeRE().FindStringSubmatch(code); m != nil {
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

// scenarioHitsTestSurface: is the scenario confronted by the `test` surface? True if ANY of
// its tags maps to a test regime. A scenario with no KNOWN regime tag is confronted
// (conservative default — presence is charged). A scenario whose known regime tags all map
// to other regimes (e2e/vr) is NOT confronted here (it lives on another surface).
//
// "Known" means MAPPED by the project (or a canonical regime name), not "looks like a
// regime". The old test was the `nivel-` prefix, and an UNMAPPED tag passed for "another
// regime": measured in the reference app, `@nivel-compilacao` is not under `regimes:`, and every
// scenario carrying it was skipped by this gate without a word.
func scenarioHitsTestSurface(sc featureScenario, testRegimes, knownRegimes map[string]bool) bool {
	sawRegimeTag := false
	for _, t := range sc.Tags {
		if testRegimes[t] {
			return true
		}
		if knownRegimes[t] {
			sawRegimeTag = true
		}
	}
	return !sawRegimeTag // no known regime tag → conservative (confront)
}

// knownRegimeTags: the tags the project maps under `regimes:`, plus the canonical regime
// names (a project with no mapping uses them directly).
func knownRegimeTags(cfg *config.Config) map[string]bool {
	out := map[string]bool{"unit": true, "integration": true, "e2e": true, "vr": true}
	if cfg != nil && cfg.Derived != nil {
		for tag := range cfg.Derived.Regimes {
			out[tag] = true
		}
	}
	return out
}

// stripLineComments remove comentários de linha (`//`, `#`, `--`) e de bloco simples para que os
// códigos citados em comentário do teste NÃO contem como implementação (coerente com
// extractCodes: comentário é referência, não posse). Barato e suficiente aqui.
func stripLineComments(s string) string {
	var b strings.Builder
	for _, ln := range strings.Split(s, "\n") {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "*") || strings.HasPrefix(t, "/*") ||
			strings.HasPrefix(t, "--") {
			continue
		}
		if strings.HasPrefix(t, "#") && !strings.HasPrefix(t, "#include") &&
			!strings.HasPrefix(t, "#define") && !strings.HasPrefix(t, "#pragma") &&
			!strings.HasPrefix(t, "#if") && !strings.HasPrefix(t, "#else") &&
			!strings.HasPrefix(t, "#endif") {
			continue
		}
		b.WriteString(cutInlineComment(ln))
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

// titleCodeRE matches a scenario code inside a title. Compiled per CALL and not in a
// `var` — the same rule as `codeRE` (rule_implemented.go): the code length comes from
// `code_lengths`, loaded AFTER the globals. In a `var` this regex froze the default `[5]`,
// and in a `[4]` project it matched no requirement at all — a title citing three
// scenarios counted zero.
func titleCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`[A-Z0-9]` + config.CodeLengthPattern() + `-[A-Z]{1,2}\d{2}(?:#\d{2})?`)
}

// titleFor finds, among the listed tests, the first whose title LEADS with `code` — alone
// or in a list of sibling codes — and returns the text after the codes, and whether that
// title cites more than one code.
//
// A test may prove several scenarios and cite all of them in the title — the form used in
// the project is `it('AATAX-S03 / AATAX-B02 / AATAX-M01: sem iniciais exibe o fallback')`,
// and it is legitimate: three scenarios describing the same situation along different
// axes of the vocabulary (state, behaviour, message) have one proof. Such a title
// describes the set, not each one, and `shared` tells the caller so.
//
// The title of that test is the text AFTER all the codes. The prefix used to be matched
// as `\[?CODE\]?`, an exact code, and that was wrong both ways: for the FIRST code the
// captured title came with the rest of the prefix stuck on ("/ AATAX-B02 / AATAX-M01: sem
// iniciais…"), and for the following ones nothing matched. Both failures became
// description drift reported where the two sides said the same.
//
// The tests come from the project's `tests` source (see projectTests): how a test opens
// is the test library's business, and reading the call here would tie the gate to one.
func titleFor(tests []testlist.Test, code string) (title string, shared, ok bool) {
	re := leadingCodesRE(code)
	for _, t := range tests {
		m := re.FindStringSubmatch(t.Title)
		if m == nil {
			continue
		}
		distinct := map[string]bool{}
		for _, c := range titleCodeRE().FindAllString(t.Title, -1) {
			distinct[c] = true
		}
		return strings.TrimSpace(m[1]), len(distinct) > 1, true
	}
	return "", false, false
}

// leadingCodesRE matches a title that opens with a list of codes containing `code`, and
// captures the text after the list.
// otherTitlesFor is the titles of the tests that lead with `code` after the first one —
// the one titleFor answers —, each without its code, leaving out a title shared with other
// codes (a table-driven name carries several and describes none alone).
func otherTitlesFor(tests []testlist.Test, code string) []string {
	re := leadingCodesRE(code)
	var out []string
	first := true
	for _, t := range tests {
		m := re.FindStringSubmatch(t.Title)
		if m == nil {
			continue
		}
		if first {
			first = false
			continue
		}
		distinct := map[string]bool{}
		for _, c := range titleCodeRE().FindAllString(t.Title, -1) {
			distinct[c] = true
		}
		if len(distinct) == 1 {
			out = append(out, strings.TrimSpace(m[1]))
		}
	}
	return out
}

func leadingCodesRE(code string) *regexp.Regexp {
	if re, ok := leadingCodesCache[code]; ok {
		return re
	}
	// The code has to END here: without the boundary, `ABCDX-DS-delta-up` matched inside
	// `ABCDX-DS-delta-up-high` and the gate read the neighbouring test's title — comparing
	// the "up to 20%" scenario with the proof of "above 20%".
	cod := `\[?` + regexp.QuoteMeta(code) + `\]?(?:[^\w#-]|$)`
	// A sibling may be numeric (`ABCDX-B01`, `ABCDX-B01#02`) or NOMINAL
	// (`ABCDX-DS-fatura-marcado`); leaving the second out stuck the sibling's prefix onto
	// the first code's title, which produced "similar 100%": same text, different
	// comparison. Each code may carry its own brackets (`[A], [B]`), or one bracket may
	// wrap the whole list (`[A / B]`).
	outro := `\[?[A-Z0-9]` + config.CodeLengthPattern() + `-(?:[A-Z]{1,2}\d{2}(?:#\d{2})?|DS-[\w-]+)\]?`
	irmaos := `(?:\s*[/,]?\s*` + outro + `)*`
	re := regexp.MustCompile(`(?s)^\s*\[?` + irmaos + `\s*[/,]?\s*` + cod + irmaos + `\]?\s*[:—-]?\s*(.*)$`)
	leadingCodesCache[code] = re
	return re
}

// leadingCodesCache avoids recompiling per scenario (large features have 40+).
var leadingCodesCache = map[string]*regexp.Regexp{}

// cutInlineComment drops a trailing comment from a code line, and only a comment.
//
// It scanned for the markers anywhere, strings included: `"--label"` cut the line at the
// dashes, so every symbol after a CLI flag argument vanished and dependency-honored
// accused dependencies the code does use; `"https://…"` cut a URL; and a Go/TS `i--`
// cut the loop body that followed it. Markers now count only OUTSIDE quotes, and `--`
// and `#` only after whitespace — the SQL/Lua and shell forms — never glued to a name.
func cutInlineComment(ln string) string {
	var quote byte
	for i := 0; i < len(ln); i++ {
		c := ln[i]
		if quote != 0 {
			if c == '\\' && quote != '`' {
				i++
			} else if c == quote {
				quote = 0
			}
			continue
		}
		switch {
		case c == '"' || c == '\'' || c == '`':
			quote = c
		case c == '/' && i+1 < len(ln) && ln[i+1] == '/':
			return ln[:i]
		case (c == '-' && i+1 < len(ln) && ln[i+1] == '-') || c == '#':
			if i > 0 && (ln[i-1] == ' ' || ln[i-1] == '\t') {
				return ln[:i]
			}
		}
	}
	return ln
}

// verdictLabel is the similarity verdict in the project's language. The verdict's own
// `String` used to be printed here, and it was Portuguese ("divergente", "limítrofe") in
// every project, whatever `lang:` said.
func verdictLabel(v similarity.Verdict) string {
	switch v {
	case similarity.Identico:
		return i18n.T("gate.feature_test_match.verdict_identical")
	case similarity.Similar:
		return i18n.T("gate.feature_test_match.verdict_similar")
	case similarity.Limitrofe:
		return i18n.T("gate.feature_test_match.verdict_borderline")
	default:
		return i18n.T("gate.feature_test_match.verdict_divergent")
	}
}

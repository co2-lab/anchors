// @anchors
//   code: DLCNA
//   ref: DLCTI

package config

import (
	"path"
	"regexp"
	"sort"
	"strings"
)

// Dialect é o LÉXICO da linguagem do projeto — a ponte entre um gate agnóstico e um
// código concreto.
//
// Por que existe: alguns gates confrontam uma verdade que vale em qualquer linguagem
// ("uma função que promete o conjunto não pode devolver a primeira página"), mas para
// enxergá-la precisam reconhecer uma FUNÇÃO EXPORTADA no arquivo. Isso é sintaxe, e
// sintaxe é do projeto, não do framework.
//
// Cravar o dialeto no gate tem um custo que já pagamos, medido: a primeira versão do
// `sibling-guard` exigia `export function` e ficou cega para 291 de 298 funções dos
// models do app de referência (97,6% são `export async function`). Um dialeto errado NÃO dá erro — dá
// silêncio, e o gate passa verde por não ver nada. Gate cego é pior que gate ausente,
// porque ocupa a linha do relatório com um ✓.
//
// Como o projeto declara, no anchors.yaml:
//
//	dialect:
//	  # family resolve tudo de uma vez para as linguagens conhecidas
//	  family: python
//	  # ou, para uma linguagem que o Anchors não conhece, declare os padrões:
//	  exported_func: '^\s*pub fn (\w+)'
//	  block_end: '^\}'
//
// `family` é atalho, não exclusividade: qualquer campo declarado explicitamente vence o
// que a família traz. E se nada for declarado, o gate que depende do dialeto se declara
// PENDENTE dizendo o que falta — nunca finge que verificou (QUALITY §7, o terceiro estado).
type Dialect struct {
	// Family resolve o léxico de uma linguagem conhecida de uma vez. Ver dialectFamilies.
	Family string `yaml:"family,omitempty"`
	// OptOut lista os campos que o projeto declara NÃO ter — e não pretende ter
	// (`opt_out: [collection_query]`). O gate que depende de um campo dispensado sai como
	// Skip silencioso em vez de Pending.
	//
	// Existe porque "não declarei ainda" e "não se aplica a mim" são estados diferentes com
	// o mesmo sintoma. Um projeto sem banco nenhum não tem `collection_query` para declarar,
	// e cobrar dele eternamente transforma o relatório em ruído que se aprende a ignorar —
	// o que apaga também os avisos que importam.
	//
	// O opt-out é EXPLÍCITO de propósito: some do relatório porque alguém decidiu, e a
	// decisão fica escrita no anchors.yaml onde a próxima pessoa a encontra. É o oposto de
	// desaparecer por omissão, que é o defeito que este campo consertaria se fosse default.
	OptOut []string `yaml:"opt_out,omitempty"`
	// ExportedFunc casa a declaração de uma função visível de fora do módulo. Deve ter ao
	// menos um grupo de captura com o NOME.
	ExportedFunc string `yaml:"exported_func,omitempty"`
	// ParamName casa o nome de um parâmetro na assinatura (um grupo de captura).
	ParamName string `yaml:"param_name,omitempty"`
	// Loop casa uma construção de repetição — o sinal de que um cursor está sendo drenado.
	Loop string `yaml:"loop,omitempty"`
	// CollectionQuery casa uma consulta que devolve MUITOS registros (o alvo do gate de
	// paginação). É a parte mais específica de stack: um projeto DynamoDB e um projeto SQL
	// não se parecem.
	CollectionQuery string `yaml:"collection_query,omitempty"`
	// Cursor casa o token de continuação do provedor de dados.
	Cursor string `yaml:"cursor,omitempty"`
	// SetPromise casa os PREFIXOS de nome que prometem o conjunto inteiro; SetSlice casa
	// os que declaram um recorte. São convenções de nomenclatura, e por isso o default
	// serve à maioria dos projetos em qualquer linguagem.
	SetPromise string `yaml:"set_promise,omitempty"`
	SetSlice   string `yaml:"set_slice,omitempty"`
	// HTTPStatus casa o STATUS HTTP que o código devolve, com um grupo de captura
	// para o número. É o léxico mais dependente de stack que existe: um handler
	// Lambda/API Gateway devolve `{ statusCode: 403 }`, um Go net/http chama
	// `w.WriteHeader(http.StatusForbidden)`, um Rails faz `render status: :forbidden`.
	// O gate `contract-status-declared` é agnóstico; este campo é o que o ensina a ler.
	HTTPStatus string `yaml:"http_status,omitempty"`
	// HTTPStatusDynamic casa a devolução de status por VARIÁVEL (`statusCode: s`),
	// o sinal de um helper que recebe o status por parâmetro. Com ele presente, o
	// gate para de afirmar que um status declarado é fantasma — os valores passados
	// ao helper podem incluí-lo por um caminho que a leitura textual não alcança.
	HTTPStatusDynamic string `yaml:"http_status_dynamic,omitempty"`
	// EnvRead casa a LEITURA de uma variável de ambiente no código, capturando o nome
	// (`os.Getenv("X")`, `process.env.X`, `System.getenv("X")`). O gate `env-declared` é
	// agnóstico; este campo é o que o ensina a ler. Um projeto com leitor próprio
	// (`config.get("X")`) declara o seu.
	EnvRead string `yaml:"env_read,omitempty"`
	// GherkinLanguage é o código de idioma do Gherkin (`en`, `pt`, `es`, `fr`…) — o que
	// vai na linha `# language:` e decide as palavras-chave da feature. Default `en`, o
	// idioma nativo do Gherkin: cravar um idioma específico obrigaria todo projeto a
	// escrever feature na língua de quem fez o framework.
	GherkinLanguage string `yaml:"gherkin_language,omitempty"`
	// GuardPatterns lista padrões adicionais ou customizados que caracterizam uma
	// guarda de parâmetro no sibling-guard. Use `{{param}}` onde o nome do parâmetro deve casar.
	GuardPatterns []string `yaml:"guard_patterns,omitempty"`
	// Imports are the import statements of this project's files, each pattern capturing the
	// imported path in its `path` group; the match is the statement, the line of `path` is
	// where its `@dep:` goes. A statement may span lines — the pattern says how. The chain
	// gates (dep-declared, dep-honored, used-by-declared) read imports only through them: the
	// core knows no syntax of an import (DESIGN-chain-patterns.md).
	Imports Patterns `yaml:"imports,omitempty"`
	// ImportNames are what an import statement brings, read on the statement's text: a
	// `name` group per name brought, or a `default` group for the module's default.
	ImportNames Patterns `yaml:"import_names,omitempty"`
	// Declarations are what a line of a file declares for others to import, read below a
	// `@used-by:` flag and by its fixer: a `name` group per name, or a `default` group for the
	// module's default. A region (`within`/`each`) is a list: its names need the flag to
	// say which one it speaks of.
	Declarations Patterns `yaml:"declarations,omitempty"`
	// TypeImportPattern recognises an import statement that brings types only — nothing of
	// the imported file runs with the importer —, read on the statement's text. The
	// dependency fixer flags such an import `@dep[type]:`.
	TypeImportPattern string `yaml:"type_import_pattern,omitempty"`
	// HandlePatterns recognise a path that HANDLES a failure instead of letting it escape.
	//
	// It is NOT "find the catch". Nailing `catch` would nail the syntax of one family of
	// languages, and handling takes many shapes that look nothing alike:
	//
	//	if (x == null) { return refuse() }        handles
	//	try { ... } catch (e) { ... }             handles
	//	match result { Err(e) => ... }            handles
	//	if err != nil { return fmt.Errorf(...) }  handles
	//
	// What they have in common is not the shape, it is the EFFECT: the failure becomes part
	// of the flow and the application carries on — which is what makes it resilient. So the
	// patterns come from the project, like the guard ones: whoever knows the shape in the
	// local dialect is whoever writes in it.
	HandlePatterns []string `yaml:"handle_patterns,omitempty"`
	// LogPatterns recognise the RECORDING of the occurrence.
	//
	// It is the piece that sustains everything else, and the one nobody charges: a handling
	// path that swallows the failure without recording anything is the perfect silence —
	// the failure happens, nothing knows, and no tool downstream has anything to read.
	LogPatterns []string `yaml:"log_patterns,omitempty"`
	// FalliblePatterns are the calls that CAN FAIL — a query, a fetch, a read from a
	// device —, each with what counts as handling it and how many lines after the call
	// the handling may stand. A unit whose code makes one is a unit that can fail, whether
	// or not anybody wrote it down: the failure gates ask its spec for the failure, and
	// each call for its handling.
	FalliblePatterns []FalliblePattern `yaml:"fallible_patterns,omitempty"`
	// SignatureEnd matches where a function's signature ends after its exported_func match —
	// the `{` of the brace languages, Python's `:` —: the parameters are read up to it.
	// Empty: the signature is what exported_func matched.
	SignatureEnd string `yaml:"signature_end,omitempty"`
	// Limits recognise a limit on a query in the code (pagination-honored): each pattern
	// captures the limit's `name` and its `value` (a literal: hidden from the caller) or its
	// `param` (an identifier: the caller's contract when it is a parameter of the function);
	// a pattern with a `type` group reads a type annotation, which is no limit passed.
	Limits Patterns `yaml:"limits,omitempty"`
	// StatementEnd matches a line of code (comments blanked) that ends a statement: the
	// failure gates climb from a call to its statement's first line until one (`;`, `}`, `)`
	// in the ts example). Empty: a statement is the call's line.
	StatementEnd string `yaml:"statement_end,omitempty"`
	// Comments say how the project's files write a comment, by glob (see Comments): the
	// header, the chain's flags and the evidence are read through them, and a file no entry
	// covers has none — the gates that need its comments say so.
	Comments []Comments `yaml:"comments,omitempty"`
	// TestHandles say how this project's code exposes a test handle, and how its flows query
	// one — the patterns the handle gates read (testid-consistent, testid-queried-exists,
	// identity-consistent). The core knows no syntax of its own: with none declared, those
	// gates measure nothing. FamilyTestHandles gives a family's example to start from.
	TestHandles *TestHandles `yaml:"test_handles,omitempty"`
	// ImportResolve says how an import's path becomes a file of the project: the
	// extensions and index files tried, and the aliases the project's build maps. The
	// dependency gates resolve each real import with it, to confront it with its `@dep:`.
	ImportResolve *ImportResolve `yaml:"import_resolve,omitempty"`
	// NavigationCall recognises a navigation call and captures its destination — the route
	// name — in the group `route`. The navigation gates confront each call with its
	// `@navigates:` flag.
	NavigationCall string `yaml:"navigation_call,omitempty"`
	// Tests says how the project's tests are written, so the gates that read a test's
	// title (feature-test-match, test-traceable, scenario-coverage) can find it. It is the
	// project's knowledge, not the engine's: the test library decides how a test opens.
	// Either a fixed `pattern` (the call that opens a test, up to its parenthesis; the
	// title is the string literal after it) or a `script` the project provides, which
	// prints the tests under the contract of the `testlist` package. Never both.
	Tests *TestsSource `yaml:"tests,omitempty"`
	// Definition recognises a name the code DEFINES — a function, a type, a class — with
	// the name in a capture group (the first non-empty one). `test-exercises-unit` reads it
	// on both sides: the names a unit defines are what its test must reach, and the same
	// name defined again in the test is the unit copied instead of exercised. Which lines
	// define a name is the language's, so it comes from the family or from the project.
	Definition string `yaml:"definition,omitempty"`
	// GenericBasenames are the file names that are a role in their folder and not the unit —
	// `handler`, `index` —: a code generated for such a file is made from its folder's name
	// (functions/manage-metadata/handler.ts → manage-metadata). Compared ignoring case.
	GenericBasenames []string `yaml:"generic_basenames,omitempty"`
}

// TestsSource is how the project's tests are read: a fixed pattern or a script.
type TestsSource struct {
	Pattern string `yaml:"pattern,omitempty"`
	Script  string `yaml:"script,omitempty"`
	// Assertion recognises an assertion inside a test's body, for `test-has-assertion`:
	// the test library's own (`expect(`, `t.Errorf`) and the project's helpers that assert.
	// A test with none passes whatever the code does.
	Assertion string `yaml:"assertion,omitempty"`
	// BlockOpen matches a line that leaves a block open at its end (`{`, `(`, `[`, Ruby's
	// `do`), BlockClose a line that starts by closing one: a test's body is read by them and
	// by its indentation. Empty: by its indentation alone.
	BlockOpen  string `yaml:"block_open,omitempty"`
	BlockClose string `yaml:"block_close,omitempty"`
	// EmptyBody matches a test opened and closed with nothing inside (`{}`): with labels, the
	// block around it stands for its body.
	EmptyBody string `yaml:"empty_body,omitempty"`
	// Placeholder matches a placeholder of a table-driven test's title (Jest's `%s`, `$name`,
	// `${name}`): not a word of the title.
	Placeholder string `yaml:"placeholder,omitempty"`
}

// GherkinKeywords são as palavras-chave da feature no idioma do projeto. Só os idiomas
// que o Gherkin oficialmente suporta entram aqui — não é uma tradução nossa, é a mesma
// tabela que Cucumber/Behave/SpecFlow usam.
type GherkinKeywords struct {
	Feature, Scenario, Outline, Given, When, Then, Examples string
	// ScenarioSynonyms are the other ways the official Gherkin table lets a scenario open
	// in this language — `Example:` in English, `Exemplo:` and the unaccented `Cenario:`
	// in Portuguese. They are READ, never written: the scaffold keeps using `Scenario`.
	//
	// Without them a reader goes blind on valid features. Measured when `non-empty` began
	// counting scenarios: a feature written with `Rule:` + `Example:` (standard Gherkin 6)
	// failed a BLOCKING gate as "feature with no scenario".
	ScenarioSynonyms []string
}

var gherkinByLang = map[string]GherkinKeywords{
	"en":    {"Feature", "Scenario", "Scenario Outline", "Given", "When", "Then", "Examples", []string{"Example", "Scenario Template"}},
	"pt":    {"Funcionalidade", "Cenário", "Esquema do Cenário", "Dado", "Quando", "Então", "Exemplos", []string{"Cenario", "Exemplo", "Esquema do Cenario"}},
	"es":    {"Característica", "Escenario", "Esquema del escenario", "Dado", "Cuando", "Entonces", "Ejemplos", []string{"Ejemplo"}},
	"fr":    {"Fonctionnalité", "Scénario", "Plan du scénario", "Soit", "Quand", "Alors", "Exemples", []string{"Exemple"}},
	"de":    {"Funktionalität", "Szenario", "Szenariogrundriss", "Angenommen", "Wenn", "Dann", "Beispiele", []string{"Beispiel"}},
	"it":    {"Funzionalità", "Scenario", "Schema dello scenario", "Dato", "Quando", "Allora", "Esempi", []string{"Esempio"}},
	"nl":    {"Functionaliteit", "Scenario", "Abstract Scenario", "Gegeven", "Als", "Dan", "Voorbeelden", []string{"Voorbeeld"}},
	"ru":    {"Функция", "Сценарий", "Структура сценария", "Дано", "Когда", "Тогда", "Примеры", []string{"Пример"}},
	"ja":    {"機能", "シナリオ", "シナリオテンプレート", "前提", "もし", "ならば", "例", nil},
	"zh-CN": {"功能", "场景", "场景大纲", "假如", "当", "那么", "例子", []string{"剧本"}},
}

// GherkinScenarioAlternatives são TODAS as formas de abrir um cenário, em todos os
// idiomas da tabela — inclusive as de Esquema/Outline. Serve a quem LÊ features: um
// parser precisa reconhecer o que qualquer projeto escreveu, não só o idioma de um.
//
// Por que a lista completa, e não o idioma configurado: um parser que só reconhecesse o
// idioma declarado ficaria cego num repositório com features herdadas noutro idioma — e
// ficar cego aqui significa reportar VERDE sobre cenários que não se enxerga. Aconteceu:
// `Esquema do Cenário` (a forma de Scenario Outline em português) não era reconhecida, e
// 108 cenários de um projeto real eram invisíveis para um gate BLOQUEANTE.
func GherkinScenarioAlternatives() []string {
	seen := map[string]bool{}
	var out []string
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for _, kw := range gherkinByLang {
		// o Esquema ANTES do Cenário: em vários idiomas o esquema CONTÉM a palavra
		// "cenário" (`Esquema do Cenário`), então a alternativa mais longa precisa ser
		// tentada primeiro — senão casa o prefixo e perde o resto.
		add(kw.Outline)
		add(kw.Scenario)
		for _, syn := range kw.ScenarioSynonyms {
			add(syn)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if len(out[i]) != len(out[j]) {
			return len(out[i]) > len(out[j]) // mais longo primeiro
		}
		return out[i] < out[j]
	})
	return out
}

// GherkinThenAlternatives são todas as formas do passo de RESULTADO, em todos os idiomas
// da tabela. Quem LÊ features precisa reconhecer o que qualquer projeto escreveu — ficar
// cego num idioma herdado significaria reportar verde sobre o que não se enxerga.
func GherkinThenAlternatives() []string {
	seen := map[string]bool{}
	var out []string
	for _, kw := range gherkinByLang {
		if kw.Then != "" && !seen[kw.Then] {
			seen[kw.Then] = true
			out = append(out, kw.Then)
		}
	}
	sort.Strings(out)
	return out
}

// GherkinExamplesAlternatives are every way to open an outline's examples table, in every
// language of the table, sorted. A reader must recognise what any project wrote.
func GherkinExamplesAlternatives() []string {
	seen := map[string]bool{}
	var out []string
	for _, kw := range gherkinByLang {
		if kw.Examples != "" && !seen[kw.Examples] {
			seen[kw.Examples] = true
			out = append(out, kw.Examples)
		}
	}
	sort.Strings(out)
	return out
}

// GherkinFor devolve as palavras-chave do idioma configurado (ou `en`).
//
// The lookup ignores case and answers the table's own spelling of the code. It used to
// lower-case the declaration and look that up, so `zh-CN` (a key with capitals) was never
// reached: a project that declared it got English keywords.
func (d Dialect) GherkinFor() (string, GherkinKeywords) {
	lang := d.GherkinLanguage
	if lang == "" {
		lang = "en"
	}
	for code, kw := range gherkinByLang {
		if strings.EqualFold(code, lang) {
			return code, kw
		}
	}
	// Idioma que o Gherkin suporta mas não temos na tabela: preserva a declaração do
	// projeto na linha `# language:` (o parser dele resolve) e usa o inglês no esqueleto,
	// em vez de traduzir por conta própria.
	return lang, gherkinByLang["en"]
}

// dialectFamilies traz o léxico das famílias sintáticas comuns. Não é uma lista fechada:
// é o atalho para quem está numa delas. Fora delas, declare os campos.
//
// A regra de inclusão aqui é a mesma dos formatos ingeridos (JUnit, lcov, Mutation
// Testing Elements): entra o que é amplamente difundido, não o que um vendor específico
// usa. Nomes de VENDOR (QueryCommand, do SDK da AWS) ficam fora — vêm de
// `collection_query`, declarado pelo projeto.
var dialectFamilies = map[string]Dialect{
	// C-like com `export`: TS, JS. Cobre `export async function`, `export const f = (…)`,
	// e os métodos exportados por classe.
	"ts": {
		// What can fail in a TS app and what handles it: React Query (the hooks and the
		// mutation), SWR, fetch and axios. A project with its own data hooks adds them, or
		// declares `fallible_patterns: []` to turn these off.
		FalliblePatterns: []FalliblePattern{
			// Each names the function a failure cites in its `name` group, takes a type
			// argument (`useQuery<Budget>(`): the generic call is the common one in TS, and
			// says in its `propagates` group that a call `return`ed or an arrow's body, with or
			// without `await`, hands its failure to the caller.
			{Call: `(?P<propagates>(?:\breturn|=>)\s*(?:await\s+)?)?\b(?P<name>use(?:Query|SuspenseQuery|InfiniteQuery|SuspenseInfiniteQuery|Queries|Mutation))(?:<[^()]*>)?\(`,
				Handled: `\bisError\b|\berror\b|\bonError\b|\bthrowOnError\b|\bstatus\b|\bcatch\b`},
			{Call: `(?P<propagates>(?:\breturn|=>)\s*(?:await\s+)?)?\b(?P<name>useSWR(?:Infinite|Immutable|Mutation)?)(?:<[^()]*>)?\(`, Handled: `\berror\b|\bonError\b`},
			{Call: `(?P<propagates>(?:\breturn|=>)\s*(?:await\s+)?)?\b(?P<name>fetch)\(`, Handled: `\.catch\(|\bcatch\b|\.ok\b|\bthrow\b`},
			{Call: `(?P<propagates>(?:\breturn|=>)\s*(?:await\s+)?)?\b(?P<name>axios)(?:\.\w+)?(?:<[^()]*>)?\(`, Handled: `\.catch\(|\bcatch\b`},
		},
		EnvRead:      "process\\.env\\.([A-Za-z_][A-Za-z0-9_]*)|process\\.env\\[\\s*['\"`]([A-Za-z_][A-Za-z0-9_]*)['\"`]\\s*\\]|import\\.meta\\.env\\.([A-Za-z_][A-Za-z0-9_]*)",
		ExportedFunc: `(?m)^export\s+(?:async\s+)?function\s+(\w+)\s*\(|^export\s+const\s+(\w+)\s*=\s*(?:async\s+)?[\(<]`,
		// `(?:[(,]|^)` ancora no abre-parêntese ou na vírgula — os parâmetros podem estar
		// todos na MESMA linha (`(userId: string, limit: number)`), então ancorar só em
		// início de linha perderia todos menos o primeiro.
		ParamName: `(?:[(,]|^)\s*(\w+)\s*[:)]`,
		Loop:      `\b(do|while|for|forEach)\b`,
		Cursor:    `(?i)(nextToken|continuationToken|nextCursor|nextPageToken|hasNextPage|offset)`,
		// Envelope do API Gateway/Lambda (`statusCode: 403`) e helpers de resposta
		// que recebem o status como 1º argumento (`jsonResponse(409, …)`, `fail(400)`).
		HTTPStatus:        `statusCode:\s*(\d{3})|\b[A-Za-z_$][\w$]*\(\s*(\d{3})\s*[,)]`,
		HTTPStatusDynamic: `statusCode:\s*[A-Za-z_$]`,
		HandlePatterns: []string{
			// `catch (e) {` and `catch {`: the binding is optional since ES2019, and a
			// handler without it was invisible (reported from the reference app).
			`\bcatch\s*[({]`,
			`\.catch\(`,
			`if\s*\([^)]*(?:==|===)\s*null`,
			`if\s*\(\s*!`,
			`\bthrow\s+new\b`,
		},
		LogPatterns: []string{
			`\b(?:logger|log|console)\s*\.\s*(?:error|warn|fatal|exception)\b`,
			`\bcaptureException\b`,
		},
		// Jest, Vitest and Mocha: `it`/`test`/`describe`, also modified by `.only`,
		// `.skip`, or `.each(table)` with the table's parentheses nested one level.
		// A test opens with `it`, `test` or `describe` — also `.only`, `.skip` or
		// `.each(table)` —, its title quoted by `'`, `"` or a backtick.
		Tests: &TestsSource{Pattern: `\b(?:it|test|describe)(?:\.each\s*\((?:[^()]|\([^()]*\))*\)|\.only|\.skip)?\s*\(\s*(?:'(?P<title>(?:[^'\\\n]|\\.)*)'|"(?P<title2>(?:[^"\\\n]|\\.)*)"|` + "`(?P<title3>[^`]*)`)",
			BlockOpen: `[{(\[]\s*$`, BlockClose: `^\s*[}\])]`, EmptyBody: `\{\s*\}`,
			Placeholder: `%[#a-zA-Z%]|\$\{[^}]*\}|\$[A-Za-z_][A-Za-z0-9_.]*`,
			Assertion:   `\bexpect(?:\.\w+)?\s*\(|\bassert(?:\.\w+)?\s*\(`},
		// Functions, classes, and constants bound to an arrow function, at any depth.
		Definition: `(?m)^\s*(?:export\s+(?:default\s+)?)?(?:async\s+)?(?:function\*?\s+(\w+)|class\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?(?:\([^)]*\)|\w+)\s*=>)`,
	},
	"go": {
		// What can fail in Go and what handles it: an HTTP request and a database call, both
		// answered by reading their error.
		FalliblePatterns: []FalliblePattern{
			// Each names the function a failure cites in its `name` group.
			{Call: `(?P<propagates>\breturn\s+)?\bhttp\.(?P<name>Get|Post|Head|PostForm)\(`, Handled: `\berr\s*!=\s*nil|\berrors\.`},
			{Call: `\.(?P<name>Do)\(`, Handled: `\berr\s*!=\s*nil|\berrors\.`},
			{Call: `\.(?P<name>(?:Query|QueryRow|Exec)(?:Context)?)\(`, Handled: `\berr\s*!=\s*nil|\.Err\(\)|\berrors\.`},
		},
		EnvRead: "os\\.(?:Getenv|LookupEnv)\\(\\s*\"([A-Za-z_][A-Za-z0-9_]*)\"\\s*\\)",
		// Em Go a exportação é a MAIÚSCULA inicial — não uma palavra-chave.
		ExportedFunc: `(?m)^func\s+(?:\([^)]*\)\s+)?([A-Z]\w*)\s*\(`,
		ParamName:    `(\w+)\s+\*?\w`,
		Loop:         `\bfor\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken|cursor|offset)`,
		// net/http: `w.WriteHeader(403)` ou a constante nomeada; e o httptest/gin
		// `c.JSON(409, …)`. A constante vira número pelo mapa do gate.
		HTTPStatus:        `WriteHeader\(\s*(\d{3})\s*\)|\bJSON\(\s*(\d{3})\s*,|http\.Status(\w+)`,
		HTTPStatusDynamic: `WriteHeader\(\s*[a-z]`,
		HandlePatterns: []string{
			// Both shapes: `if err != nil` and `if err := f(); err != nil`. The second is
			// the commoner one in Go (986 of them in Anchors itself) and the first pattern
			// alone missed it, so the handling it does was invisible to the failure gates.
			`\bif\b[^\n{]*\berr\s*!=\s*nil`,
			`\brecover\(\)`,
			`==\s*nil\s*\{`,
			`\berrors\.(?:Is|As)\(`,
			// An error held in a FIELD, not a variable: GORM's `if result.Error != nil`.
			// The `err` word above never saw it, and a repository that handled every
			// failure of its queries failed `failure-handled` (reported from baas-proxy).
			`\.Error\s*!=\s*nil`,
			// A sentinel error returned: `return ErrNotFound`, `return nil, ErrX`. The
			// failure becomes a value the caller decides on — handled here.
			`\breturn\b[^\n]*\bErr[A-Z]\w*`,
		},
		LogPatterns: []string{
			`\b(?:log|logger|slog)\.(?:Error|Warn|Fatal|Printf|Print)\b`,
			`\bfmt\.Errorf\(`,
			// Propagating is recording in Go, like `fmt.Errorf`: the error goes up to whoever
			// logs it. Two shapes the patterns above missed — the error field returned
			// (`return nil, result.Error`, never `e.Error()`, which turns it into text) and
			// the sentinel returned.
			`\breturn\b[^\n]*\.Error\s*(?:$|[,;)}])`,
			`\breturn\b[^\n]*\bErr[A-Z]\w*`,
		},
		// The standard library: a test names its cases with `t.Run("title", …)`.
		Tests: &TestsSource{Pattern: `\bt\.Run\(\s*(?:"(?P<title>(?:[^"\\\n]|\\.)*)"|` + "`(?P<title2>[^`]*)`)",
			Assertion: `\bt\.(?:Error|Errorf|Fatal|Fatalf|Fail|FailNow)\b`,
			BlockOpen: `[{(\[]\s*$`, BlockClose: `^\s*[}\])]`, EmptyBody: `\{\s*\}`},
		// Functions, methods and types, at the top level.
		// A method is named WITH its receiver's type (the `owner` group): a test's fake that
		// implements an interface defines `fakePinger.Ping`, and the unit's `GormPinger.Ping`
		// is another definition — the same bare name read them as a copy.
		Definition: `(?m)^(?:func\s+\(\s*(?:\w+\s+)?\*?(?P<owner>\w+)(?:\[[^\]]*\])?\s*\)\s+(\w+)|func\s+(\w+)|type\s+(\w+))`,
	},
	"python": {
		// What can fail in Python and what handles it: an HTTP request through `requests`.
		FalliblePatterns: []FalliblePattern{
			{Call: `(?P<propagates>\breturn\s+)?\brequests\.(?P<name>get|post|put|patch|delete|head|request)\(`, Handled: `\bexcept\b|raise_for_status|status_code`},
		},
		EnvRead: "os\\.(?:getenv|environ\\.get)\\(\\s*['\"]([A-Za-z_]\\w*)['\"]|os\\.environ\\[\\s*['\"]([A-Za-z_]\\w*)['\"]\\s*\\]",
		// Sem palavra-chave de exportação: convenção é o underscore inicial marcar o
		// privado, então "exportada" = def no nível do módulo sem `_`.
		ExportedFunc: `(?m)^(?:async\s+)?def\s+([a-zA-Z]\w*)\s*\(`,
		Definition:   `(?m)^\s*(?:async\s+)?def\s+(\w+)|^\s*class\s+(\w+)`,
		ParamName:    `(\w+)\s*[:=,)]`,
		Loop:         `\b(for|while)\b`,
		Cursor:       `(?i)(next_token|continuation_token|next_cursor|next_page_token|offset)`,
	},
	"java": {
		EnvRead:      "System\\.getenv\\(\\s*\"([A-Za-z_]\\w*)\"\\s*\\)",
		ExportedFunc: `(?m)^\s*public\s+(?:static\s+)?(?:async\s+)?[\w<>\[\],\s]+\s+(\w+)\s*\(`,
		ParamName:    `\w[\w<>\[\]]*\s+(\w+)\s*[,)]`,
		Loop:         `\b(do|while|for)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken|pageable)`,
	},
	"rust": {
		EnvRead:      "env::var(?:_os)?\\(\\s*\"([A-Za-z_]\\w*)\"\\s*\\)",
		ExportedFunc: `(?m)^\s*pub\s+(?:async\s+)?fn\s+(\w+)\s*[\(<]`,
		ParamName:    `(\w+)\s*:`,
		Loop:         `\b(for|while|loop)\b`,
		Cursor:       `(?i)(next_token|continuation_token|next_cursor|offset)`,
	},
	"ruby": {
		EnvRead:      "ENV\\[\\s*['\"]([A-Za-z_]\\w*)['\"]\\s*\\]|ENV\\.fetch\\(\\s*['\"]([A-Za-z_]\\w*)['\"]",
		ExportedFunc: `(?m)^\s*def\s+([a-z]\w*)`,
		ParamName:    `(\w+)\s*[:,)]`,
		Loop:         `\b(each|while|until|for|loop)\b`,
		Cursor:       `(?i)(next_token|continuation_token|next_cursor|offset)`,
	},
	"php": {
		EnvRead:      "getenv\\(\\s*['\"]([A-Za-z_]\\w*)['\"]\\s*\\)|\\$_ENV\\[\\s*['\"]([A-Za-z_]\\w*)['\"]\\s*\\]",
		ExportedFunc: `(?m)^\s*(?:public\s+)?function\s+(\w+)\s*\(`,
		ParamName:    `\$(\w+)`,
		Loop:         `\b(do|while|for|foreach)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextCursor|offset)`,
	},
	"csharp": {
		EnvRead:      "Environment\\.GetEnvironmentVariable\\(\\s*\"([A-Za-z_]\\w*)\"\\s*\\)",
		ExportedFunc: `(?m)^\s*public\s+(?:static\s+)?(?:async\s+)?[\w<>\[\],\s]+\s+(\w+)\s*\(`,
		ParamName:    `\w[\w<>\[\]]*\s+(\w+)\s*[,)]`,
		Loop:         `\b(do|while|for|foreach)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken)`,
	},
	// OCaml: a top-level `let` is visible from outside the module unless an `.mli` hides
	// it; `open M` is the import line (reported from the bavel adoption, where these
	// patterns were first declared by hand).
	"ocaml": {
		EnvRead:      "Sys\\.getenv(?:_opt)?\\s+\"([A-Za-z_][A-Za-z0-9_]*)\"",
		ExportedFunc: `(?m)^let\s+(?:rec\s+)?([a-z_][A-Za-z0-9_']*)`,
		Loop:         `\b(for|while)\b`,
	},
	"kotlin": {
		EnvRead:      "System\\.getenv\\(\\s*\"([A-Za-z_]\\w*)\"\\s*\\)",
		ExportedFunc: `(?m)^\s*(?:public\s+)?(?:suspend\s+)?fun\s+(\w+)\s*[\(<]`,
		ParamName:    `(\w+)\s*:`,
		Loop:         `\b(do|while|for|forEach)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken)`,
	},
}

// Estes dois defaults NÃO são de linguagem — são convenções de NOMENCLATURA, difundidas o
// bastante para servirem de piso em qualquer stack (`listUsers`, `list_users`,
// `GetAllItems`). O projeto sobrepõe se o vocabulário dele for outro (outro idioma, por
// exemplo: `listar`, `buscarTodos`).
const (
	// O verbo de conjunto pode vir DEPOIS de um prefixo de provedor: `cognitoListDevices`,
	// `s3ListObjects`, `dynamo_query_all` prometem o conjunto tanto quanto `listDevices`.
	// Só o `^` deixava o prefixo ESCONDER a promessa — medido no caso que originou a regra
	// do cursor descartado: `cognitoListDevices` jogava fora o PaginationToken e o gate nem
	// olhava, porque o nome não casava.
	//
	// O verbo tem de COMEÇAR e TERMINAR palavra, senão `all` interno gera falso positivo:
	// `allocateSlot` e `callbackUrl` casavam com a versão frouxa. Início = começo da string,
	// `_`, ou minúscula/dígito antes (camelCase); fim = final, `_`, ou próxima maiúscula.
	//
	// Sem lookbehind/lookahead: Go (RE2) não os tem, e o regex simplesmente NÃO COMPILA —
	// o gate trata o padrão como nil e passa a aprovar tudo em silêncio. Aconteceu aqui, e
	// o teste do caso real foi o que pegou. Em RE2 a fronteira se expressa consumindo o
	// caractere vizinho num grupo, o que basta porque só o VERBO é consultado.
	defaultSetPromise = `(?:^|_|[a-z0-9])(?:[Ll]ist|[Gg]etAll|get_all|[Ff]indAll|find_all|[Ff]etchAll|fetch_all|[Qq]ueryAll|query_all|[Aa]ll|[Tt]odos|[Ll]istar)(?:$|_|[A-Z]|\b)`
	defaultSetSlice   = `(?i)^(list|get|find|fetch|query|listar|buscar)_?(First|Recent|Latest|Top|Page|Some|Sample|Next|Primeiros|Recentes)`
)

// DialectFor devolve o dialeto EFETIVO: a família (se declarada) com os campos
// explícitos do projeto por cima, e os defaults de nomenclatura no piso.
func (c *Config) DialectFor() Dialect {
	var d Dialect
	if c != nil && c.Dialect != nil {
		d = *c.Dialect
	}
	// Nothing comes from the family: what a language looks like is the project's to declare
	// (DESIGN-core-without-syntax.md). The family's lexicon is an example — FamilyLexicon —,
	// written by `anchors init` and `anchors migrate`, and offered by the doctor.
	return d
}

// IsGenericBasename says whether a file's name is one the project calls generic — a role in
// its folder, the folder being the unit (`dialect.generic_basenames`, compared ignoring case).
func (d Dialect) IsGenericBasename(name string) bool {
	for _, g := range d.GenericBasenames {
		if strings.EqualFold(g, name) {
			return true
		}
	}
	return false
}

// ExampleGenericBasenames are the example of `generic_basenames`: the roles a file plays in
// the folder of a serverless function, a package or a route, never applied by itself.
func ExampleGenericBasenames() []string {
	return []string{"handler", "index", "resource", "mod", "main", "types", "route", "routes", "schema", "config"}
}

// NamingExample is the example of the two naming conventions — set_promise and set_slice —,
// a list of verbs in English and Portuguese: every family's example carries it.
func NamingExample() (promise, slice string) { return defaultSetPromise, defaultSetSlice }

// FamilyLexicon is the family's example of the lexicon the gates read — exported_func,
// param_name, loop, cursor, http_status, env_read, handle_patterns, log_patterns, tests,
// definition, statement_end — with the naming conventions set_promise and set_slice: what
// `anchors init` and `anchors migrate` write and the doctor offers, never applied by itself.
// An unknown family gives nil.
func FamilyLexicon(family string) *Dialect {
	base, ok := dialectFamilies[strings.ToLower(family)]
	if !ok {
		return nil
	}
	ex := Dialect{
		ExportedFunc: base.ExportedFunc, ParamName: base.ParamName, Loop: base.Loop,
		Cursor: base.Cursor, HTTPStatus: base.HTTPStatus, HTTPStatusDynamic: base.HTTPStatusDynamic,
		EnvRead: base.EnvRead, Definition: base.Definition, Tests: base.Tests,
		HandlePatterns: append([]string(nil), base.HandlePatterns...),
		LogPatterns:    append([]string(nil), base.LogPatterns...),
		SetPromise:     defaultSetPromise, SetSlice: defaultSetSlice,
		StatementEnd:     familyStatementEnd(family),
		Limits:           exampleLimits(),
		SignatureEnd:     familySignatureEnd(family),
		GuardPatterns:    exampleGuards(family),
		GenericBasenames: ExampleGenericBasenames(),
	}
	return &ex
}

// familySignatureEnd is where a function's signature ends in the family's code.
func familySignatureEnd(family string) string {
	switch strings.ToLower(family) {
	case "python":
		return `:\s*(?:#.*)?\n`
	case "ruby":
		return `\n`
	}
	return `\{`
}

// exampleGuards are the guards most languages write over a parameter, `{{param}}` its name:
// a filter by it, an early exit or a throw because of it, a default for it.
func exampleGuards(family string) []string {
	return []string{
		`\.filter\([^)]*{{param}}\b`,
		`\b{{param}}\s*\.\s*filter\(`,
		`if\s*\(?[^\n{:]*\b{{param}}\b[^\n{:]*\)?\s*[:{]?\s*(throw|return|raise|panic)`,
		`\b{{param}}\b[^\n]*\?\?`,
		`(throw|assert|raise|panic)[^\n]*\b{{param}}\b`,
	}
}

// exampleLimits are the limits most data APIs write — `limit`, `max`, `pageSize`, `take`,
// `top` —, as a default (`limit = 100`, past a type annotation), a literal (`Limit: 50`) or
// an identifier passed (`Limit: limit`), with the type names a signature annotates them with.
func exampleLimits() Patterns {
	name := `(?i)\b(?P<name>limit|max|maxItems|maxKeys|maxResults|pageSize|page_size|take|top|first)\b`
	return Patterns{
		{Re: name + `[^=\n]*?=\s*(?P<value>\d+)`},
		{Re: name + `\s*[:=]\s*(?P<value>\d+)`},
		{Re: name + `\s*[:=]\s*(?P<param>\w+)`},
		{Re: name + `\s*:\s*(?P<type>(?i:number|int|integer|long|short|byte|float|double|decimal|uint\d*|int\d*|usize|i\d+|u\d+|bigint|size_t))\b`},
	}
}

// familyStatementEnd is the line ending a statement in the family's code: the brace and
// semicolon languages end one at `;`, `}` or `)`; the others have no example.
func familyStatementEnd(family string) string {
	switch strings.ToLower(family) {
	case "ts", "go", "java", "kotlin", "csharp", "rust", "php":
		return `[;})]$`
	}
	return ""
}

// KnownDialectFamilies lista as famílias com léxico embutido (para mensagens de erro e
// para o `anchors init` perguntar).
func KnownDialectFamilies() []string {
	out := make([]string, 0, len(dialectFamilies))
	for k := range dialectFamilies {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Compile compila um padrão do dialeto, devolvendo nil se vazio ou inválido — quem chama
// trata nil como "o projeto não declarou este padrão" e responde Pendente.
func (d Dialect) Compile(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re
}

// WaivedField diz se o projeto declarou opt-out para um campo do dialeto.
//
// A comparação é sobre o NOME YAML do campo (`collection_query`), não o do Go: quem escreve
// o opt-out está lendo o anchors.yaml, e pedir a tradução mental para `CollectionQuery`
// convidaria ao erro silencioso — um nome que não casa nada é opt-out que não vale, e o
// gate seguiria pendente sem explicar por quê.
func (d Dialect) WaivedField(campoYAML string) bool {
	for _, c := range d.OptOut {
		if strings.EqualFold(strings.TrimSpace(c), campoYAML) {
			return true
		}
	}
	return false
}

// EnvReadPattern is how this project's code reads an environment variable: its own
// `env_read`, or "" — no family's is read by itself, nor every family's together.
func (d Dialect) EnvReadPattern() string { return d.EnvRead }

// FalliblePattern is a call that can fail, and what handles it (see FalliblePatterns). The
// families know their own (React Query, SWR, fetch and axios for `ts`; HTTP and SQL for
// `go`; `requests` for `python`) — see FamilyFalliblePatterns.
type FalliblePattern struct {
	// Call recognises the call that can fail (a regular expression). Its `name` group, when
	// it has one — `\b(?P<name>apiCall)(?:<[^()]*>)?\(` —, is the name a failure cites the
	// call by: the pattern knows its language, and the core reads no syntax of its own. Its
	// `propagates` group, when it takes part in the match —
	// `(?P<propagates>\breturn\s+)?\bhttp\.(?P<name>Get)\(` —, says the call's result is the
	// unit's own, so its failure is the caller's to handle.
	Call string `yaml:"call"`
	// Handled recognises its handling — the error read, a catch, an error callback.
	Handled string `yaml:"handled,omitempty"`
	// Window is how many lines after the call the handling may stand; 0 is the default.
	Window int `yaml:"window,omitempty"`
}

// DefaultFallibleWindow is the window when a pattern declares none: a destructuring and
// the early returns of a component fit in it, the next component does not.
const DefaultFallibleWindow = 25

// FamilyFalliblePatterns are what the family knows can fail, and what handles it. They are
// not applied by themselves: `anchors init` writes them into a new project's dialect, and
// the governance tips suggest them to a project that declared none — a project's failure
// gates are measured with them before they bar with them.
func FamilyFalliblePatterns(family string) []FalliblePattern {
	base, ok := dialectFamilies[strings.ToLower(family)]
	if !ok {
		return nil
	}
	return append([]FalliblePattern(nil), base.FalliblePatterns...)
}

// TestHandles are the patterns that find test handles in code and in flows (see
// Dialect.TestHandles). Each pattern names what it captures in a group:
//
//   - `testid`: a handle, whole — `testID=":abcd-save"`.
//   - `prefix`: the head of a family of handles the rest of which is runtime data — a
//     template, a helper composing `${base}-${option}` —; it reads as `head-*`. A value one
//     pattern captures as a prefix is not also a whole handle by another.
//   - `suffix` (Mentioned only): a segment composed onto a head that comes from outside the
//     file — `${testID}-row`; it reads as `*-row-*`.
type TestHandles struct {
	// Exposed: where the code exposes a handle. What these find is what a unit owns — the
	// handles its spec declares, the acronyms its identity is confronted by.
	Exposed Patterns `yaml:"exposed,omitempty"`
	// Mentioned: what else in the code makes a handle exist without exposing it at the
	// attribute — a fallback, a table of handles, a suffix composed in a child. Only
	// testid-queried-exists reads it, to know a queried handle exists somewhere.
	Mentioned Patterns `yaml:"mentioned,omitempty"`
	// Queried: how an end-to-end flow queries a handle — testid-queried-exists reads it.
	Queried Patterns `yaml:"queried,omitempty"`
}

// validate says which pattern is wrong (see Pattern.Validate).
func (th *TestHandles) validate() error {
	for _, l := range []struct {
		field    string
		patterns Patterns
		groups   []string
	}{
		{"exposed", th.Exposed, []string{"testid", "prefix", "suffix"}},
		{"mentioned", th.Mentioned, []string{"testid", "prefix", "suffix"}},
		{"queried", th.Queried, []string{"testid", "prefix"}},
	} {
		if err := l.patterns.Validate("dialect.test_handles."+l.field, l.groups); err != nil {
			return err
		}
	}
	return nil
}

// FamilyTestHandles is the family's example of TestHandles for the attribute the project
// marks its elements with (`testID`, `data-testid`): what `anchors init` writes and the
// governance tips suggest. It is never applied by itself — the gates read only what the
// project declares. A family with no example, or no attribute, gives none.
func FamilyTestHandles(family, attr string) *TestHandles {
	if strings.ToLower(family) != "ts" || strings.TrimSpace(attr) == "" {
		return nil
	}
	// The attribute, or a prop derived from it (`backTestID`, `confirmTestID`), as a JSX
	// attribute or an object key.
	a := `(?i)\b[a-z]*` + regexp.QuoteMeta(attr) + `[a-z]*`
	id := `:?[a-zA-Z][a-zA-Z0-9._-]*`
	re := func(s string) Pattern { return Pattern{Re: s} }
	return &TestHandles{
		Exposed: Patterns{
			re(a + `\s*[=:]\s*\{?["'](?P<testid>` + id + `)["']`),
			// A prop ending in Prefix gives a child the head of its handles.
			re(a + `prefix\s*[=:]\s*\{?["'](?P<prefix>` + id + `)["']`),
			re(a + "\\s*[=:]\\s*\\{?`(?P<testid>" + id + ")`"),
			re(a + "\\s*[=:]\\s*\\{?`(?P<prefix>" + id + ")\\$\\{"),
			// A handle whose head is the caller's — `${testID}-option-${v}` —: what the unit
			// adds to it is its own segment, the head an id it was given.
			re(a + "\\s*[=:]\\s*\\{?`\\$\\{[^}]*\\}(?P<suffix>[^`$]*)"),
			// A conditional handle: every literal right after a `?` or a `:` of it — its
			// branches, nested or not —, never a condition's (`kind === 'email'`). The region
			// starts at the first `?`, and may hold templates (`${i}`).
			{Within: a + "\\s*=\\s*\\{[^?{}]*(?P<within>\\?(?:[^{}]|\\$\\{[^{}]*\\})*)\\}",
				Each: `[?:]\s*["'](?P<testid>` + id + `)["']|[?:]\s*` + "`(?P<prefix>" + id + ")\\$\\{"},
		},
		Mentioned: Patterns{
			// A handle marked with `:` anywhere: a fallback, a table of handles.
			re("[\"'`](?P<testid>:[a-zA-Z][a-zA-Z0-9._-]*[a-zA-Z0-9._])[\"'`]"),
			re("[\"'`](?P<prefix>:[a-zA-Z][a-zA-Z0-9._-]*)\\$\\{"),
			// A segment a child composes onto the head its parent gives it.
			re("`\\$\\{[a-zA-Z_][a-zA-Z0-9_.]*\\}(?P<suffix>-[a-zA-Z][a-zA-Z0-9-]*)"),
		},
		// Maestro flows: `id: ':abcd-save'`, or a regular expression whose head is static.
		Queried: Patterns{
			re(`(?m)^\s*id:\s*["'](?P<testid>` + id + `)["']\s*$`),
			re(`(?m)^\s*id:\s*["'](?P<prefix>:?[a-zA-Z][a-zA-Z0-9_-]*)[.*+?\[(]`),
		},
	}
}

// FamilyChain is the family's example of the chain patterns — Imports, ImportNames,
// Declarations, TypeImportPattern, ImportResolve —, as a dialect holding only them: what
// `anchors init` writes and the governance tips suggest. It is never applied by itself: the
// chain gates read only what the project declares. A family with no example gives nil.
func FamilyChain(family string) *Dialect {
	ident := `[\p{L}_$][\p{L}\p{N}_$]*`
	re := func(s string) Pattern { return Pattern{Re: s} }
	switch strings.ToLower(family) {
	case "ts":
		code := "**/*.{ts,tsx,js,jsx,mjs,cjs,mts,cts}"
		return &Dialect{
			Imports: Patterns{
				// `import … from '…'` and `export … from '…'`: a list of names in braces may
				// span lines; anything else stays on the keyword's line.
				// Both plainly run, unless `import type`.
				{Re: `(?m)^[ \t]*(?:import|export)\b(?:[^'";{}\n]|\{[^}]*\})*?\bfrom\s*['"](?P<path>[^'"]+)['"]`, Files: code, Runs: true},
				// An import for its effect alone: `import './polyfills'`.
				{Re: `(?m)^[ \t]*import\s*['"](?P<path>[^'"]+)['"]`, Files: code, Runs: true},
				// `require('…')` and `import('…')`, with what is destructured from them or the
				// member read from them.
				{Re: `(?:\{[^{}]*\}\s*=\s*(?:await\s+)?)?\b(?:require|import)\(\s*['"](?P<path>[^'"]+)['"]\s*\)(?:\.` + ident + `)?`, Files: code},
				// A Maestro flow running another, relative to itself.
				{Re: `(?m)^\s*(?:-\s*)?(?:runFlow:|file:)\s*['"]?(?P<path>[^'"\s#{}]+\.ya?ml)\b`, Files: "**/*.{yaml,yml}", Relative: true},
			},
			ImportNames: Patterns{
				// The names in braces: `{ A, B as C }` brings A and B, `type A` brings A.
				{Within: `^[ \t]*(?:import|export)\b[^{'"]*\{(?P<within>[^}]*)\}`, Each: `(?:^|,)\s*(?:type\s+)?(?P<name>` + ident + `)`},
				// The default import: `import X from`, `import X, { … } from`.
				re(`^[ \t]*import\s+(?:type\s+)?(?P<default>` + ident + `)\s*(?:,|\s+from\b)`),
				// A destructuring of a dynamic import: `const { a, b: c } = await import('…')`.
				{Within: `\{(?P<within>[^{}]*)\}\s*=\s*(?:await\s+)?(?:require|import)\(`, Each: `(?:^|,)\s*(?P<name>` + ident + `)`},
				// A member read from it: `import('…').Name`.
				re(`\b(?:require|import)\([^)]*\)\.(?P<name>` + ident + `)`),
			},
			Declarations: Patterns{
				// At the top level of the file: a name declared inside a function is no one's to
				// import.
				re(`(?m)^export\s+default\b(?P<default>)`),
				re(`(?m)^module\.exports\s*=(?P<default>)`),
				re(`(?m)^(?:export\s+)?(?:declare\s+)?(?:default\s+)?(?:abstract\s+)?(?:async\s+)?(?:function\*?|const|let|var|class|interface|type|enum)\s+(?P<name>` + ident + `)`),
				// An export list, on one line or several: `export { A, B as C } from '…'`
				// declares A and C.
				{Within: `(?m)^export\s+(?:type\s+)?\{(?P<within>[^}]*)\}`, Each: `(?P<name>` + ident + `)\s*(?:,|$)`},
			},
			TypeImportPattern: `^\s*(?:import|export)\s+type\b`,
			ImportResolve:     &ImportResolve{Extensions: []string{".ts", ".tsx", ".js", ".jsx", "/index.ts", "/index.tsx", "/index.js"}},
		}
	case "go":
		goIdent := `[\p{L}_][\p{L}\p{N}_]*`
		return &Dialect{
			Imports: Patterns{
				{Re: `(?m)^import\s+(?:[\w.]+\s+)?"(?P<path>[^"]+)"`, Files: "**/*.go"},
				// An import block: each path on its own line.
				{Within: `(?m)^import\s*\((?P<within>[^)]*)\)`, Each: `(?m)^[ \t]*(?:[\w.]+\s+)?"(?P<path>[^"]+)"`, Files: "**/*.go"},
			},
			Declarations: Patterns{
				re(`(?m)^func\s+(?:\([^)]*\)\s*)?(?P<name>` + goIdent + `)`),
				re(`(?m)^(?:type|var|const)\s+(?P<name>` + goIdent + `)`),
			},
		}
	case "python":
		pyIdent := `[\p{L}_][\p{L}\p{N}_]*`
		return &Dialect{
			Imports: Patterns{
				{Re: `(?m)^[ \t]*from\s+(?P<path>[\w.]+)\s+import\b[^\n(]*(?:\([^)]*\))?`, Files: "**/*.py"},
				{Re: `(?m)^[ \t]*import\s+(?P<path>[\w.]+)`, Files: "**/*.py"},
			},
			ImportNames: Patterns{
				{Within: `\bimport\s+\(?(?P<within>[^)]*)\)?`, Each: `(?:^|,)\s*(?P<name>` + pyIdent + `)`},
			},
			// `app.core` is `app/core.py` or the package `app/core/`; a project maps its
			// top-level packages with aliases (`"app/": "src/app/"`).
			ImportResolve: &ImportResolve{Separator: ".", Extensions: []string{".py", "/__init__.py"}},
			Declarations: Patterns{
				re(`(?m)^(?:async\s+)?def\s+(?P<name>` + pyIdent + `)`),
				re(`(?m)^class\s+(?P<name>` + pyIdent + `)`),
				re(`(?m)^(?P<name>` + pyIdent + `)\s*(?::[^=\n]*)?=`),
			},
		}
	case "ocaml":
		return &Dialect{
			Imports: Patterns{
				// A module named after `open`, or qualifying a name anywhere: `Eval.apply`. A
				// constructor — `Some x` — qualifies nothing, and is not read.
				{Re: `\bopen!?\s+(?P<path>[A-Z][A-Za-z0-9_']*)`, Files: "**/*.ml"},
				{Re: `\b(?P<path>[A-Z][A-Za-z0-9_']*)\.[A-Za-z_(]`, Files: "**/*.ml"},
			},
			// dune names a module's file after it, first letter lowered.
			ImportResolve: &ImportResolve{ByName: true, NameCase: "lower_first", Extensions: []string{".ml"}},
		}
	}
	return nil
}

// FamilyNavigationCall is the family's example of `navigation_call`, its destination in the
// group `route`: what `anchors init` writes and the governance tips suggest, never applied by
// itself. A family with none gives "".
func FamilyNavigationCall(family string) string {
	switch strings.ToLower(family) {
	case "ts":
		// React Navigation and the routers that look like it: `navigation.navigate('Goal')`,
		// `router.push({ pathname: '/goal' })`, the back and reset calls with no route.
		return `\b(?:navigation|navigator|router|nav)\s*\.\s*(?:navigate|push|replace|reset|goBack|back|pop|popToTop|dispatch)\s*\(\s*(?:\{\s*(?:name|pathname)\s*:\s*)?(?:['"](?P<route>[A-Za-z0-9_/.:-]+)['"])?`
	}
	return ""
}

// FamilyMocks is the family's example of how its tests write a double, for the mock gates —
// `derived.mock_detect` (the double and the module it replaces), `mock_contract` (what binds
// it to the module), `mock_factory` (a double with a hand-written factory), `mock_typed` (a
// type binding the factory): what `anchors init` and `anchors migrate` write and the doctor
// offers, never applied by itself. A family with none gives nil.
func FamilyMocks(family string) *Derived {
	switch strings.ToLower(family) {
	case "ts":
		return &Derived{
			MockDetect:   "(?:jest|vi)\\s*\\.\\s*mock\\s*\\(\\s*['\"`]([^'\"`]+)",
			MockContract: "Partial<typeof import('{{module}}')>",
			// A second argument that is no options object: the factory. An automock
			// (`jest.mock('x')`) or Vitest's `{ spy: true }` mirrors the module.
			MockFactory: "(?:jest|vi)\\s*\\.\\s*mock\\s*\\(\\s*(?:['\"`][^'\"`]*['\"`]|[^,()]*\\([^()]*\\))\\s*,\\s*[^\\s{)]",
			// A return annotation between the factory's parameters and its arrow.
			MockTyped: "\\)\\s*:\\s*\\S[^\\n]*=>",
		}
	case "python":
		return &Derived{MockDetect: "(?:mock\\.)?patch\\(['\"]([^'\"]+)", MockContract: "autospec=True"}
	}
	return nil
}

// FamilyRoutePattern is the family's example of `derived.route_pattern` — how the app
// registers a route, its name in a capture group —: React Navigation's and Expo Router's
// `name="X"`, a stack's param list `X: undefined`, a CDK API's `addResource('x')`. Never
// applied by itself; a family with none gives "".
func FamilyRoutePattern(family string) string {
	if strings.ToLower(family) == "ts" {
		return `name="([A-Za-z][A-Za-z0-9_]*)"|(?m)^\s{2,}([A-Za-z][A-Za-z0-9_]*)\s*:\s*(?:undefined|\{)|addResource\('([a-z0-9][a-z0-9/_-]*)'`
	}
	return ""
}

// ImportResolve is how an import's path becomes a file (see Dialect.ImportResolve).
type ImportResolve struct {
	// Extensions are tried in order after the path: `.ts`, `.tsx`, `/index.ts`.
	Extensions []string `yaml:"extensions,omitempty"`
	// Aliases map a path prefix to a directory of the project: `"@/": "apps/mobile/src/"`.
	Aliases map[string]string `yaml:"aliases,omitempty"`
	// ByName: the import pattern captures a module NAME, not a path, wherever it stands on a
	// line — a qualified reference in the middle of the code (OCaml's `Eval.apply`), with no
	// import line to carry a flag. The name resolves to the project's code file named after
	// it (NameCase, Extensions), and the file's `@dep:` flags are its own, one per module it
	// uses, on lines of their own below the header.
	ByName bool `yaml:"by_name,omitempty"`
	// NameCase is how a module name becomes a file name: `lower_first` (`Js_pilha` →
	// `js_pilha`), `lower`, or empty for the name as written.
	NameCase string `yaml:"name_case,omitempty"`
	// Separator is what separates the directories of an import path when it is not `/`:
	// Python's `.` (`app.core.money`). Leading separators climb from the importer — `.x` is
	// beside it, `..x` one directory up.
	Separator string `yaml:"separator,omitempty"`
	// Scoped are aliases that hold only for the importers a glob matches — each workspace's
	// own `paths`, `@/` meaning `apps/mobile/` in one app and `apps/landing-page/src/` in
	// another. An importer's scoped aliases win over Aliases for the same prefix.
	Scoped []ScopedAliases `yaml:"scoped,omitempty"`
}

// ScopedAliases are aliases for the importers its glob matches (see ImportResolve.Scoped).
type ScopedAliases struct {
	Files   string            `yaml:"files"`
	Aliases map[string]string `yaml:"aliases"`
}

// AliasesFor are the aliases an import of the file at rel resolves through: Aliases, with
// those of every scoped entry whose glob matches rel over them — the first entry wins a
// prefix two entries declare.
func (r *ImportResolve) AliasesFor(rel string) map[string]string {
	if r == nil {
		return nil
	}
	if len(r.Scoped) == 0 {
		return r.Aliases
	}
	out := map[string]string{}
	for k, v := range r.Aliases {
		out[k] = v
	}
	rel = strings.TrimPrefix(path.Clean(strings.ReplaceAll(rel, "\\", "/")), "./")
	for i := len(r.Scoped) - 1; i >= 0; i-- {
		if (&Comments{Files: r.Scoped[i].Files}).matches(rel) {
			for k, v := range r.Scoped[i].Aliases {
				out[k] = v
			}
		}
	}
	return out
}

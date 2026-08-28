package config

import (
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
	// GherkinLanguage é o código de idioma do Gherkin (`en`, `pt`, `es`, `fr`…) — o que
	// vai na linha `# language:` e decide as palavras-chave da feature. Default `en`, o
	// idioma nativo do Gherkin: cravar um idioma específico obrigaria todo projeto a
	// escrever feature na língua de quem fez o framework.
	GherkinLanguage string `yaml:"gherkin_language,omitempty"`
}

// GherkinKeywords são as palavras-chave da feature no idioma do projeto. Só os idiomas
// que o Gherkin oficialmente suporta entram aqui — não é uma tradução nossa, é a mesma
// tabela que Cucumber/Behave/SpecFlow usam.
type GherkinKeywords struct {
	Feature, Scenario, Outline, Given, When, Then, Examples string
}

var gherkinByLang = map[string]GherkinKeywords{
	"en":    {"Feature", "Scenario", "Scenario Outline", "Given", "When", "Then", "Examples"},
	"pt":    {"Funcionalidade", "Cenário", "Esquema do Cenário", "Dado", "Quando", "Então", "Exemplos"},
	"es":    {"Característica", "Escenario", "Esquema del escenario", "Dado", "Cuando", "Entonces", "Ejemplos"},
	"fr":    {"Fonctionnalité", "Scénario", "Plan du scénario", "Soit", "Quand", "Alors", "Exemples"},
	"de":    {"Funktionalität", "Szenario", "Szenariogrundriss", "Angenommen", "Wenn", "Dann", "Beispiele"},
	"it":    {"Funzionalità", "Scenario", "Schema dello scenario", "Dato", "Quando", "Allora", "Esempi"},
	"nl":    {"Functionaliteit", "Scenario", "Abstract Scenario", "Gegeven", "Als", "Dan", "Voorbeelden"},
	"ru":    {"Функция", "Сценарий", "Структура сценария", "Дано", "Когда", "Тогда", "Примеры"},
	"ja":    {"機能", "シナリオ", "シナリオテンプレート", "前提", "もし", "ならば", "例"},
	"zh-CN": {"功能", "场景", "场景大纲", "假如", "当", "那么", "例子"},
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

// GherkinFor devolve as palavras-chave do idioma configurado (ou `en`).
func (d Dialect) GherkinFor() (string, GherkinKeywords) {
	lang := strings.ToLower(d.GherkinLanguage)
	if lang == "" {
		lang = "en"
	}
	if kw, ok := gherkinByLang[lang]; ok {
		return lang, kw
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
	},
	"go": {
		// Em Go a exportação é a MAIÚSCULA inicial — não uma palavra-chave.
		ExportedFunc: `(?m)^func\s+(?:\([^)]*\)\s+)?([A-Z]\w*)\s*\(`,
		ParamName:    `(\w+)\s+\*?\w`,
		Loop:         `\bfor\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken|cursor|offset)`,
		// net/http: `w.WriteHeader(403)` ou a constante nomeada; e o httptest/gin
		// `c.JSON(409, …)`. A constante vira número pelo mapa do gate.
		HTTPStatus:        `WriteHeader\(\s*(\d{3})\s*\)|\bJSON\(\s*(\d{3})\s*,|http\.Status(\w+)`,
		HTTPStatusDynamic: `WriteHeader\(\s*[a-z]`,
	},
	"python": {
		// Sem palavra-chave de exportação: convenção é o underscore inicial marcar o
		// privado, então "exportada" = def no nível do módulo sem `_`.
		ExportedFunc: `(?m)^(?:async\s+)?def\s+([a-zA-Z]\w*)\s*\(`,
		ParamName:    `(\w+)\s*[:=,)]`,
		Loop:         `\b(for|while)\b`,
		Cursor:       `(?i)(next_token|continuation_token|next_cursor|next_page_token|offset)`,
	},
	"java": {
		ExportedFunc: `(?m)^\s*public\s+(?:static\s+)?(?:async\s+)?[\w<>\[\],\s]+\s+(\w+)\s*\(`,
		ParamName:    `\w[\w<>\[\]]*\s+(\w+)\s*[,)]`,
		Loop:         `\b(do|while|for)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken|pageable)`,
	},
	"rust": {
		ExportedFunc: `(?m)^\s*pub\s+(?:async\s+)?fn\s+(\w+)\s*[\(<]`,
		ParamName:    `(\w+)\s*:`,
		Loop:         `\b(for|while|loop)\b`,
		Cursor:       `(?i)(next_token|continuation_token|next_cursor|offset)`,
	},
	"ruby": {
		ExportedFunc: `(?m)^\s*def\s+([a-z]\w*)`,
		ParamName:    `(\w+)\s*[:,)]`,
		Loop:         `\b(each|while|until|for|loop)\b`,
		Cursor:       `(?i)(next_token|continuation_token|next_cursor|offset)`,
	},
	"php": {
		ExportedFunc: `(?m)^\s*(?:public\s+)?function\s+(\w+)\s*\(`,
		ParamName:    `\$(\w+)`,
		Loop:         `\b(do|while|for|foreach)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextCursor|offset)`,
	},
	"csharp": {
		ExportedFunc: `(?m)^\s*public\s+(?:static\s+)?(?:async\s+)?[\w<>\[\],\s]+\s+(\w+)\s*\(`,
		ParamName:    `\w[\w<>\[\]]*\s+(\w+)\s*[,)]`,
		Loop:         `\b(do|while|for|foreach)\b`,
		Cursor:       `(?i)(nextToken|continuationToken|nextPageToken)`,
	},
	"kotlin": {
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
	if base, ok := dialectFamilies[strings.ToLower(d.Family)]; ok {
		// o explícito do projeto sempre vence a família
		if d.ExportedFunc == "" {
			d.ExportedFunc = base.ExportedFunc
		}
		if d.ParamName == "" {
			d.ParamName = base.ParamName
		}
		if d.Loop == "" {
			d.Loop = base.Loop
		}
		if d.Cursor == "" {
			d.Cursor = base.Cursor
		}
		if d.CollectionQuery == "" {
			d.CollectionQuery = base.CollectionQuery
		}
		if d.HTTPStatus == "" {
			d.HTTPStatus = base.HTTPStatus
		}
		if d.HTTPStatusDynamic == "" {
			d.HTTPStatusDynamic = base.HTTPStatusDynamic
		}
	}
	if d.SetPromise == "" {
		d.SetPromise = defaultSetPromise
	}
	if d.SetSlice == "" {
		d.SetSlice = defaultSetSlice
	}
	return d
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

// Dispensado diz se o projeto declarou opt-out para um campo do dialeto.
//
// A comparação é sobre o NOME YAML do campo (`collection_query`), não o do Go: quem escreve
// o opt-out está lendo o anchors.yaml, e pedir a tradução mental para `CollectionQuery`
// convidaria ao erro silencioso — um nome que não casa nada é opt-out que não vale, e o
// gate seguiria pendente sem explicar por quê.
func (d Dialect) Dispensado(campoYAML string) bool {
	for _, c := range d.OptOut {
		if strings.EqualFold(strings.TrimSpace(c), campoYAML) {
			return true
		}
	}
	return false
}

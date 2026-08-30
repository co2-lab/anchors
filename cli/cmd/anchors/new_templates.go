package main

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// templates é o catálogo de molduras por kind. As seções são AGNÓSTICAS (o piso
// universal do Anchors); um projeto especializa via seu próprio guide. O objetivo é
// só nascer conforme: header @anchors + identidade + as seções que a régua pede.
//
// Parametrização: cada seção tem Default (entra sempre, salvo --without) ou não
// (entra com --with). A ordem do slice é a ordem no arquivo.
var templates = map[string]template{
	"spec":    specTemplate,
	"feature": featureTemplate,
	"test":    testTemplate,
	"plan":    planTemplate,
}

// ─── header helpers (dialeto de comentário por artefato, HEADER_GUIDE) ──────────

func mdHeader(idField string) func(id, outPath string) string {
	// A spec é sempre Markdown (é o formato do Anchors, não da linguagem do projeto):
	// o dialeto de comentário aqui é fixo com razão.
	return func(id, _ string) string {
		return fmt.Sprintf("<!-- @anchors\n  %s: %s\n  updated_at: TODO\n  layer: TODO\n-->\n", idField, id)
	}
}

// feature usa dialeto Gherkin: comentário `#` após a linha de linguagem.
// A linha `# language:` e as palavras-chave vêm do IDIOMA declarado pelo projeto.
// Cravar `pt` obrigava todo projeto a escrever feature em português — a língua de quem
// fez o framework, não a de quem o usa. Default `en`, o idioma nativo do Gherkin.
// Ver dialect.GherkinFor().
func featureHeader(id string, lang string) string {
	if lang == "" {
		lang = "en"
	}
	return fmt.Sprintf("# language: %s\n# @anchors\n#   ref: %s\n#   updated_at: TODO\n#   layer: feature\n", lang, id)
}

// lineHeader emite o header no dialeto de comentário do ARQUIVO que está sendo criado —
// deduzido da extensão do caminho de saída, via config.CommentMarkers (que já cobre 20+
// linguagens). Fixar `//` aqui gerava um `.py` com comentário de JavaScript: sintaxe
// inválida, arquivo que nem roda. O `anchors init` oferece 18 stacks (Java, Python, Go,
// Rust, PHP, Ruby, Elixir, Dart, C++…); emitir só um dialeto contradizia a própria
// promessa do init.
func lineHeader(idField string) func(id, outPath string) string {
	return func(id, outPath string) string {
		c := config.LineCommentFor(outPath)
		return fmt.Sprintf("%s @anchors\n%s   %s: %s\n%s   updated_at: TODO\n%s   layer: test\n",
			c, c, idField, id, c, c)
	}
}

// ─── SPEC ───────────────────────────────────────────────────────────────────────

var specTemplate = template{
	kind:     "spec",
	ext:      ".spec.md",
	idField:  "code",
	headerFn: mdHeader("code"),
	sections: []section{
		{Key: "title", Title: "Cabeçalho com código + propósito", Default: true,
			Body: "# {name} — TODO propósito em uma frase\n\n> **Código**: `{id}`\n\n"},
		// `route` vem logo após o título: a rota pertence ao BLOCO DE CABEÇALHO da spec
		// (é onde o gate route-declared a procura), não ao corpo.
		{Key: "route", Title: "Rota (spec de tela)", Default: false,
			Purpose: "TELA navegável: declara como se chega até ela. Sem isso o grafo de navegação fica com um nó solto.",
			Body:    "> **Rota**: `TODO`\n\n"},
		{Key: "overview", Title: "Visão geral (o que faz e para quem)", Default: true,
			Body: "## Visão Geral\nTODO: o que a unidade faz e para quem.\n\n"},
		{Key: "rules", Title: "Regras/efeitos catalogados (cada um com ID)", Default: true, Realizes: "B",
			Body: "## Regras\n\n### {id}-B01 — TODO regra\nDescreva o comportamento (não a implementação).\n\n"},
		{Key: "contract", Title: "Contrato de entrada/saída", Default: false,
			Purpose:  "Unidade com entrada e saída DISTINTAS que merecem descrição separada (regra de negócio, validador). Prefira `signature` quando a assinatura da função já diz tudo.",
			Variants: []string{"signature"},
			Body:     "## Contrato\n**Entrada**: TODO.\n**Saída**: TODO.\n\n"},
		// NÃO emitir bloco de código aqui: a spec descreve COMPORTAMENTO, e guides de
		// projeto costumam proibir TypeScript/tipos dentro da spec (o app de referência proíbe
		// explicitamente). A assinatura é descrita em tabela — entrada, saída, e o que
		// cada uma significa — sem prender a spec a uma linguagem.
		{Key: "signature", Title: "Assinatura (contrato condensado, em tabela)", Default: false,
			Purpose:  "Função PURA cuja entrada/saída cabem em uma tabela curta. Comum em business-logic de backend. Prefira `contract` quando entrada e saída pedirem descrição longa.",
			Variants: []string{"contract"},
			Body:     "## Signature\n| Parâmetro | Tipo | Descrição |\n| --- | --- | --- |\n| TODO | TODO | TODO |\n\n**Retorna**: TODO.\n\n"},
		{Key: "effects", Title: "Efeitos (o que a unidade provoca)", Default: false, Realizes: "B",
			Purpose:  "Alternativa enxuta a `rules` para unidade sem ramificação: lista o que ela provoca, cada item com seu código. Prefira `rules` quando houver decisão/condicional a catalogar.",
			Variants: []string{"rules"},
			Body:     "## Efeitos\n| Efeito | Descrição |\n| --- | --- |\n| `{id}-B01` | TODO |\n\n"},
		// A seção da FRONTEIRA DE ENTRADA. Nasce da medição: 71% das specs de um projeto
		// real tinham seção de regras/efeitos e apenas 14% diziam o que a unidade ACEITA.
		// Todo defeito de borda encontrado em três rodadas de review adversarial morava no
		// que ninguém tinha declarado.
		//
		// Não confundir com `constraints`: aquela diz o que a unidade NÃO FAZ (delimita
		// RESPONSABILIDADE, empurra o dever para fora); esta diz o que ela ACEITA (delimita
		// ENTRADA, e nomeia quem garante). São complementares — e só a segunda fecha o
		// circuito. Medido: três specs escreveram "não valido a chave" em `constraints`,
		// cada uma correta, e o dever ficou órfão; o resultado foi perda silenciosa de dado
		// do usuário com todos os gates verdes.
		//
		// A coluna QUEM GARANTE é a peça central. Sem ela, "fora do domínio" vira mais uma
		// forma de dizer "não é meu problema" — e o problema não é de ninguém. Um `?` ali é
		// visível ao escrever; um dever órfão espalhado por três arquivos, não.
		{Key: "domain", Title: "Domínio (o que a unidade ACEITA, e quem garante)", Default: false,
			Purpose: "A fronteira de ENTRADA: que valores existem, quais estão fora, e QUEM garante que o inválido não chega. Diferente de `constraints` (o que a unidade não faz): aqui se declara o que ela recebe. Toda unidade com entrada externa precisa — é onde moram os defeitos de borda que nenhum gate pega.",
			Body: "## Domínio\n\n| Entrada | Aceita | Fora do domínio | Quem garante |\n" +
				"| --- | --- | --- | --- |\n" +
				"| TODO | TODO: os valores válidos | TODO: o que NÃO pode chegar aqui | TODO: quem barra (esta unidade? o chamador? a interface?) |\n\n" +
				"> `Quem garante` não pode ficar vazio nem dizer só \"não é meu\": se ninguém\n" +
				"> garante, o dever é órfão — e é exatamente aí que a entrada inválida passa.\n\n"},
		{Key: "invariants", Title: "Invariantes (o que vale SEMPRE)", Default: false, Realizes: "I",
			Purpose: "Verdade que sobrevive a QUALQUER sequência de operações — \"escrever V no mês M ⟹ ler M devolve V\", \"remover e repor volta ao estado anterior\". Diferente de efeito (o que uma chamada faz) e de restrição (o limite da camada). Toda unidade com PRODUTOR e CONSUMIDOR precisa de ao menos um: é onde moram os bugs que o teste-por-regra não pega.",
			Body:    "## Invariantes\n\n| Regra | Vale sempre | Como se prova |\n| --- | --- | --- |\n| `{id}-I01` | TODO: a verdade que não pode ser violada | ciclo fechado: aplica o produtor, verifica pelo consumidor |\n\n"},
		{Key: "constraints", Title: "Restrições (o que a unidade NÃO faz)", Default: false, Realizes: "X",
			Purpose: "O limite da camada — o que é responsabilidade de outro. Vale sempre que a fronteira for confundível (ex.: regra pura que não busca dado).",
			Body:    "## Restrições\n| Regra | Limite | Por quê |\n| --- | --- | --- |\n| `{id}-X01` | TODO | TODO |\n\n"},
		{Key: "errors", Title: "Erros / Falhas", Default: false, Realizes: "E",
			Purpose: "Como a unidade falha e quem sinaliza. Use quando houver entrada inválida, dependência externa ou estado impossível a tratar.",
			Body:    "## Erros / Falhas\n| Regra | Condição | Falha |\n| --- | --- | --- |\n| `{id}-E01` | TODO | TODO |\n\n"},
		{Key: "constants", Title: "Constantes de negócio", Default: false,
			Purpose: "Números/limites que são DECISÃO de negócio (teto de plano, janela de meses) — evita a constante virar mágica no código.",
			Body:    "## Constantes de Negócio\n| Constante | Valor | Por quê |\n| --- | --- | --- |\n| TODO | TODO | TODO |\n\n"},
		// As três seções de ESTADO GLOBAL. Vieram de um projeto real que as havia
		// desenhado sozinho num template próprio — sinal de que o preset faltava, não de
		// que o projeto era peculiar: shape/actions/selectors é o vocabulário de Redux,
		// Zustand, Pinia, MobX e NgRx, não de um deles.
		{Key: "state-shape", Title: "Shape do estado (os campos e seus defaults)", Default: false, Realizes: "S",
			Purpose: "Store/estado global: QUAIS campos existem, de que tipo, com que valor inicial. Sem isto, o consumidor descobre o shape lendo a implementação.",
			Body:    "## Shape do Estado\n| Campo | Tipo | Default | Descrição |\n| --- | --- | --- | --- |\n| TODO | TODO | TODO | TODO |\n\n"},
		{Key: "actions", Title: "Actions (o que muda o estado)", Default: false, Realizes: "A",
			Purpose: "Store: cada ação com o EFEITO sobre o estado. A regra que decide a transição mora na business-logic; aqui fica o que a ação faz ao estado.",
			Body:    "## Actions\n| Regra | Ação | Efeito no estado |\n| --- | --- | --- |\n| `{id}-B01` | TODO | TODO |\n\n"},
		{Key: "hydration", Title: "Hidratação / Persistência", Default: false, Realizes: "B",
			Purpose: "Store: o que sobrevive ao fechar o app, de onde o estado é reidratado e quando. É a fonte de bugs de estado obsoleto que nenhum teste de unidade pega.",
			Body:    "## Hidratação / Persistência\n| Regra | O que persiste | De onde hidrata |\n| --- | --- | --- |\n| `{id}-B02` | TODO | TODO |\n\n"},
		{Key: "messages", Title: "Mensagens ao usuário", Default: false, Realizes: "M",
			Purpose: "Validação/formulário: o texto que o usuário LÊ quando a regra reprova. Fica na spec porque é decisão de produto, não de implementação — e é o que o teste asserta.",
			Body:    "## Mensagens ao Usuário\n| Regra | Condição | Mensagem |\n| --- | --- | --- |\n| `{id}-M01` | TODO | TODO |\n\n"},
		{Key: "states", Title: "Estados e transições (unidade de fluxo)", Default: false, Realizes: "S",
			Purpose: "Unidade com ESTADOS observáveis (tela, máquina de estado, hook com loading/erro/vazio). Não use em função pura.",
			Body:    "## Estados\n| Estado | Quando | O que mostra |\n| --- | --- | --- |\n| TODO | TODO | TODO |\n\n"},
		// A pergunta que nenhuma outra secao responde: quando o dado CRESCER, esta unidade
		// continua correta? `states` diz o que a tela mostra ENQUANTO carrega (skeleton,
		// vazio); esta diz COMO o conjunto chega — de uma vez, paginado ou por scroll
		// infinito.
		//
		// Nao da para inferir do codigo: uma lista de 3 categorias fixas e uma de transacoes
		// se escrevem com o mesmo `.map()`. So quem conhece o dominio sabe qual e qual, e e
		// por isso que se DECLARA.
		//
		// Medido no projeto que originou a secao: 68 telas iteravam colecao com `.map()`
		// direto e 3 tinham scroll infinito. Uma delas listava tudo sem paginar e nenhum
		// gate viu — o `pagination-honored` mede a FUNCAO de dados (procura a consulta de
		// colecao), e a tela nao tem consulta nenhuma: chama um hook.
		{Key: "loading", Title: "Carregamento da colecao", Default: false, Realizes: "B",
			Purpose: "Unidade que lista uma COLECAO que cresce. Declare a estrategia: tudo de uma vez (e por que a colecao e limitada), paginada, ou scroll infinito — e qual peca a implementa.",
			Body:    "## Carregamento\n| Regra | Estrategia | Origem | Comportamento |\n| --- | --- | --- | --- |\n| `{id}-B01` | TODO: tudo de uma vez / paginada / scroll infinito | TODO: o hook ou consulta | TODO: se limitada, POR QUE; se paginada, o tamanho da pagina |\n\n"},
		{Key: "auth", Title: "Auth/Acesso", Default: false, Realizes: "R",
			Purpose: "Unidade cujo acesso depende de quem é o usuário (permissão, plano, dono do dado).",
			Body:    "## Auth/Acesso\nTODO: quem pode; regra de acesso.\n\n"},
		{Key: "deps", Title: "Tabela de dependências", Default: false,
			Purpose: "O que a unidade consome. Símbolos entre `crases` viram CONTRATO verificável (gate dependency-honored); descrição em prosa não é cobrada.",
			Body:    "## Dependências\n| Cód | Arquivo | Método | Camada |\n| --- | --- | --- | --- |\n| DEP1 | TODO | `TODO` | TODO |\n\n"},
		// A seção da AMBIGUIDADE. Default: true — é a única seção cuja ausência esconde
		// justamente o que ela existe para revelar. Quem não tem dúvida gasta uma palavra
		// ("nenhuma"); quem tem, ganha um lugar declarado para ela em vez de um comentário
		// de PR que morre no merge.
		{Key: "open", Title: "Decisões em aberto (o que a spec ainda NÃO decide)", Default: true, Realizes: "Q",
			Purpose: "O que quem implementa vai precisar e a spec não responde. Enquanto houver item aqui, implementar é adivinhar — e a adivinhação não é confrontada por gate nenhum, porque todas as peças existem. Cada pergunta ganha um CÓDIGO (`{CODE}-Q01`): sem identidade ela não vira issue rastreável nem sobrevive a uma reescrita da spec. Quando a resposta vier, PROMOVA-A a regra (com código) e marque o item como resolvido citando a regra que nasceu dela; a pergunta fica no rastro, não é varrida.",
			Body:    "## Decisões em aberto\n\n| Código | Pergunta | Quem decide | Vira |\n| --- | --- | --- | --- |\n\nnenhuma\n\n"},
		{Key: "notes", Title: "Notas de implementação", Default: false,
			Purpose: "Decisão técnica não-óbvia que o código não explica sozinho (por que assim, e não do jeito esperado).",
			Body:    "## Notas de Implementação\nTODO: o que um leitor futuro precisa saber.\n\n"},
	},
}

// ─── FEATURE ─────────────────────────────────────────────────────────────────────

var featureTemplate = template{
	kind:     "feature",
	ext:      ".feature",
	idField:  "ref",
	headerFn: featureHeader,
	sections: []section{
		// {FEATURE}/{SCENARIO}/{GIVEN}… são resolvidos pelo IDIOMA do projeto no render.
		{Key: "funcionalidade", Title: "Funcionalidade + tags", Default: true,
			Body: "\n@{id}\n{FEATURE}: {name}\n\n"},
		// O `Então` do esqueleto NÃO pode ser "o efeito {id}-B01 se verifica": isso é a
		// tautologia "dado X, então X" — não afirma comportamento observável nenhum e
		// delega ao teste a definição do que "se verifica" significa. Dois cenários de
		// regras diferentes ficam indistinguíveis, e o gate `feature-test-match` (que
		// confronta CÓDIGO e descrição) passa satisfeito.
		//
		// Medido: um esqueleto assim gerou 359 cenários tautológicos num projeto real, e um
		// revisor externo identificou isso como a razão ESTRUTURAL pela qual uma regra
		// (`-B05`) ficou sem caso discriminante — nada, entre a spec e o teste, obrigou
		// alguém a escrever o resultado esperado.
		//
		// O TODO explícito custa uma edição e impede o pior desfecho: um cenário que
		// atravessa o pipeline sem nunca ter dito o que deveria acontecer.
		{Key: "cenario", Title: "Cenário-esqueleto com scenario-code", Default: true,
			Body: "  @{id}-B01 {UNIT_TAG}\n  {SCENARIO}: TODO — o que este cenário exercita\n" +
				"    {GIVEN} TODO: o estado de partida (dados concretos, não \"um estado válido\")\n" +
				"    {WHEN} TODO: a ação\n" +
				"    {THEN} TODO: o RESULTADO OBSERVÁVEL, com o valor esperado.\n" +
				"    # Não escreva \"o efeito {id}-B01 se verifica\" — isso não afirma nada e\n" +
				"    # passa nos gates sem dizer o que deveria acontecer.\n\n"},
		{Key: "esquema", Title: "Esquema de cenário (variações por exemplo)", Default: false,
			Body: "  @{id}-B02 {UNIT_TAG}\n  {OUTLINE}: TODO — a variação que os exemplos percorrem\n" +
				"    {GIVEN} <entrada>\n    {THEN} o resultado é <saida>\n\n    {EXAMPLES}:\n" +
				"      | entrada | saida |\n      | TODO    | TODO  |\n\n"},
	},
}

// ─── TEST ────────────────────────────────────────────────────────────────────────

var testTemplate = template{
	kind:     "test",
	ext:      ".test.ts", // só o exemplo da mensagem de --out; o arquivo real usa o caminho dado
	idField:  "ref",
	headerFn: lineHeader("ref"),
	sections: []section{
		{Key: "describe", Title: "caso de teste com o scenario-code no nome", Default: true,
			// O corpo é resolvido pelo DIALETO do projeto — ver testBody. Este Body é o
			// fallback (família não declarada/desconhecida).
			Body: "{TEST_BODY}"},
	},
}

// testBody devolve o esqueleto de caso de teste na família do projeto.
//
// O que é do FRAMEWORK e o que é do projeto: a estrutura é universal — um caso nomeado
// com o SCENARIO-CODE (`[{id}-B01]`), porque é isso que o gate `scenario-coverage` lê do
// JUnit para saber que aquele requisito tem teste verde. A SINTAXE é do projeto.
//
// Emitir `describe/it` para todo mundo criava um `.py` com JavaScript dentro — sintaxe
// inválida, arquivo que não roda. E contradizia o próprio `anchors init`, que oferece 18
// stacks (Java, Python, Go, Rust, PHP, Ruby, Elixir, Dart, C++…).
//
// O que NÃO se emite: uma escolha de framework de teste. Dentro da família, o esqueleto
// usa a construção da BIBLIOTECA PADRÃO ou a convenção dominante da linguagem
// (`func TestX` em Go vale para testify e para o testing puro; `def test_x` em Python
// vale para pytest e unittest). Onde a linguagem não decide sozinha, o comentário no
// esqueleto diz o que ajustar.
func testBody(family, name, id string) string {
	code := id + "-B01"
	switch family {
	case "python":
		return fmt.Sprintf("\n\ndef test_%s():\n    \"\"\"[%s] TODO\"\"\"\n    # TODO: arrange / act / assert\n    pass\n",
			toSnake(name), code)
	case "go":
		return fmt.Sprintf("\nfunc Test%s(t *testing.T) {\n\t// [%s] TODO\n\t// TODO: arrange / act / assert\n}\n",
			toPascal(name), code)
	case "java", "kotlin":
		return fmt.Sprintf("\n@Test\nvoid %s() {\n    // [%s] TODO\n    // TODO: arrange / act / assert\n}\n",
			toCamel(name), code)
	case "csharp":
		return fmt.Sprintf("\n[Fact]\npublic void %s()\n{\n    // [%s] TODO\n    // TODO: arrange / act / assert\n}\n",
			toPascal(name), code)
	case "rust":
		return fmt.Sprintf("\n#[test]\nfn %s() {\n    // [%s] TODO\n    // TODO: arrange / act / assert\n}\n",
			toSnake(name), code)
	case "ruby":
		return fmt.Sprintf("\ndescribe '%s' do\n  it '[%s] TODO' do\n    # TODO: arrange / act / assert\n  end\nend\n",
			name, code)
	case "php":
		return fmt.Sprintf("\npublic function test%s(): void\n{\n    // [%s] TODO\n    // TODO: arrange / act / assert\n}\n",
			toPascal(name), code)
	case "ts", "":
		return fmt.Sprintf("\ndescribe('%s', () => {\n  it('[%s] TODO', () => {\n    // TODO: arrange / act / assert\n  })\n})\n",
			name, code)
	default:
		// Família declarada que não conhecemos: o esqueleto vira instrução, não um chute
		// de sintaxe. O que importa é o scenario-code no NOME do caso — é o que o
		// `anchors ingest --junit` cruza com a spec.
		return fmt.Sprintf("\nTODO: um caso de teste com `[%s]` no NOME (é o que amarra o "+
			"teste ao requisito da spec, via o JUnit que `anchors ingest` consome).\n", code)
	}
}

func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + 32)
			continue
		}
		if r == '-' || r == ' ' || r == '.' {
			b.WriteByte('_')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func toPascal(s string) string {
	var b strings.Builder
	up := true
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			up = true
			continue
		}
		if up && r >= 'a' && r <= 'z' {
			b.WriteRune(r - 32)
		} else {
			b.WriteRune(r)
		}
		up = false
	}
	return b.String()
}

func toCamel(s string) string {
	p := toPascal(s)
	if p == "" {
		return p
	}
	r := []rune(p)
	if r[0] >= 'A' && r[0] <= 'Z' {
		r[0] += 32
	}
	return string(r)
}

// ─── PRESETS ─────────────────────────────────────────────────────────────────────
//
// Um preset é um CONJUNTO NOMEADO de seções para um tipo de unidade — o atalho de quem
// já sabe o que a camada exige. Substitui a escolha manual de --with/--without (que
// continua valendo, e pode refinar o preset).
//
// Presets existem porque a mesma camada tem dialetos legítimos: uma função pura de
// backend quer `signature`+`effects`; uma regra de negócio de mobile quer
// `contract`+`rules`+`constraints`; uma tela quer `route`+`states`+`auth`. Forçar um
// formato único geraria seções vazias (ou, pior, prosa inventada para preenchê-las).
var specPresets = map[string]presetDef{
	"backend-logic": {
		Desc:     "função PURA de backend — a assinatura é o contrato",
		Sections: []string{"title", "overview", "signature", "domain", "effects", "invariants", "constraints", "deps", "open"},
	},
	"mobile-logic": {
		Desc:     "regra de negócio do app — entrada/saída + regras catalogadas",
		Sections: []string{"title", "overview", "contract", "domain", "rules", "invariants", "constraints", "deps", "open"},
	},
	"screen": {
		Desc: "tela navegável — rota, estados observáveis, dados e acesso",
		// Medido em 96 specs de tela de um projeto real: `messages` aparece em 50 e
		// `notes` em 92, e o preset não emitia nenhuma das duas. Faltar a de mensagens
		// numa tela de FORMULÁRIO é o pior caso — é onde mora o texto que o usuário lê
		// quando erra, e a seção que ninguém emite é a seção que ninguém escreve.
		//
		// `domain` sai: a tela não é a fronteira de entrada do dado (2 das 96 a usam);
		// quem valida é o hook ou o handler abaixo dela.
		Sections: []string{"title", "route", "overview", "states", "loading", "rules", "messages", "auth", "deps", "notes", "open"},
	},
	"component": {
		Desc:     "componente de UI — props e estados visuais, sem rota",
		Sections: []string{"title", "overview", "contract", "states", "rules", "open"},
	},
	"handler": {
		Desc:     "interface do backend (Lambda/rota) — request/response, auth e erro",
		Sections: []string{"title", "overview", "contract", "domain", "rules", "auth", "errors", "deps", "open"},
	},
	"schema": {
		Desc: "interface do DADO — modelos, índices e autorização (quem lê/escreve)",
		// Sem `contract`: um modelo de dado não tem entrada/saída — ele TEM FORMA, que é o
		// que `domain` descreve. Com `notes`: índice, migração e limite do provedor são
		// exatamente o que um leitor futuro precisa e não cabe em regra.
		//
		// Medido contra 50 specs de modelo de um projeto real: com este preset, 8 de 8
		// seções coincidem; com o anterior, 5 de 8 — e as 3 divergentes eram as centrais.
		Sections: []string{"title", "overview", "domain", "rules", "auth", "constraints", "deps", "notes", "open"},
	},
	"hook": {
		Desc: "hook/composable — o que ele faz acontecer e o limite da camada",
		// Sem `states`: medido em 37 specs de hook de um projeto real, "Fluxo de Estados"
		// aparece em UMA. O estado observável é da TELA (`screen` tem `states`); o hook
		// entrega dado e efeito. Emiti-la por padrão produzia seção vazia — ou, pior,
		// prosa inventada para preenchê-la.
		//
		// `signature` fica (13 de 37 a usam), mas depois de `effects`: o que os 37 têm em
		// comum é dizer O QUE PROVOCAM; a assinatura é detalhe de quem tem contrato
		// complexo. A ordem das seções é o fio de leitura.
		Sections: []string{"title", "overview", "effects", "signature", "constraints", "deps", "open"},
	},
	"store": {
		Desc:     "estado global (Redux/Zustand/Pinia/MobX) — shape, actions e hidratação",
		Sections: []string{"title", "overview", "state-shape", "actions", "hydration", "invariants", "constraints", "deps", "open"},
	},
	"validation": {
		Desc:     "regra de validação — critérios e a mensagem que o usuário lê",
		Sections: []string{"title", "overview", "contract", "domain", "rules", "messages", "deps", "open"},
	},
	"service": {
		Desc:     "serviço — operação, dependência externa e como falha",
		Sections: []string{"title", "overview", "contract", "domain", "rules", "errors", "constraints", "deps", "open"},
	},
	"repository": {
		Desc:     "acesso a dado — operações e limites da camada",
		Sections: []string{"title", "overview", "domain", "effects", "constraints", "errors", "deps", "open"},
	},
}

// presetDef é um preset nomeado: a lista ordenada de seções e para que serve.
type presetDef struct {
	Desc     string
	Sections []string
}

// planTemplate — o esqueleto de um PLANO.
//
// Faltava, e a ausência tinha custo: quem escreve um plano copia o anterior, e o que o
// anterior esqueceu se propaga. As três coisas que este esqueleto força são as que os
// planos escritos à mão mais omitem — as fases com código, a ordem entre elas, e o que
// fica de fora.
var planTemplate = template{
	kind:     "plan",
	ext:      ".md",
	idField:  "code",
	headerFn: mdHeader("code"),
	sections: []section{
		{Key: "title", Title: "Cabeçalho com código + objetivo", Default: true,
			Body: "# {name}\n\n> **Código**: `{id}`\n\n"},

		{Key: "objective", Title: "Objetivo (o que este plano entrega)", Default: true,
			Purpose: "O ESTADO que o repositório alcança quando o plano fecha — não a lista de tarefas. " +
				"Um objetivo que descreve atividade (\"criar os arquivos X e Y\") não diz quando parar; " +
				"um que descreve estado (\"o teste roda e o sinal chega ao mapa\") diz.",
			Body: "## Objetivo\n\nTODO: o estado que o repositório alcança quando este plano fechar.\n\n"},

		{Key: "why", Title: "Motivo (por que agora, e por que assim)", Default: true,
			Purpose: "O que acontece se este plano NÃO for feito, ou for feito depois. É o que " +
				"permite despriorizá-lo com consciência em vez de por esquecimento.",
			Body: "## Motivo\n\nTODO: o que quebra sem isto, ou o que fica mais caro depois.\n\n"},

		{Key: "phases", Title: "Fases (com código e ordem)", Default: true,
			Realizes: "F",
			Purpose: "Cada fase é um item CATALOGADO (`{id}-F01`), e a ordem entre elas é declarada — " +
				"`(depende de {id}-F01)`. Sem o código, a ordem vive em prosa e o pipeline não a " +
				"confronta: as specs semeadas nascem todas disponíveis, e o agente pega a da fase 3 " +
				"com a fase 1 em aberto. Cada spec declara `parent: {id}-F0N` (a que fase pertence) " +
				"e `needs:` (de qual depende).",
			Body: "## Fases\n\n### {id}-F01 — TODO nome da primeira fase\n\n" +
				"TODO: o que esta fase entrega, e as specs que ela semeia.\n\n" +
				"- [ ] `caminho/Unidade.spec.md` — TODO o que ela descreve\n\n" +
				"### {id}-F02 — TODO nome da segunda fase (depende de {id}-F01)\n\n" +
				"TODO: por que esta fase só começa depois da anterior.\n\n"},

		{Key: "out-of-scope", Title: "Fora de escopo (o que este plano NÃO faz)", Default: true,
			Purpose: "O que alguém razoavelmente esperaria daqui e não vai encontrar, com o lugar " +
				"onde está. Sem isto, o plano seguinte é escrito assumindo que o anterior cobriu.",
			Body: "## Fora de escopo\n\n- TODO: o que fica para outro plano, e qual.\n\n"},

		{Key: "done", Title: "Definição de pronto", Default: true,
			Purpose: "O que se confere para dizer que o plano fechou, em termos VERIFICÁVEIS — um " +
				"comando que roda, um gate que passa. \"Está funcionando\" não é definição de pronto.",
			Body: "## Definição de pronto\n\n- TODO: o comando que passa, ou o gate que fica verde.\n\n"},

		{Key: "revision", Title: "Revisão de outro plano (quando este revisa)", Default: true,
			Purpose: "Só quando este plano REVISA outro. Declare `supersedes: plans/00XX-nome.md` no " +
				"header, e escreva aqui o que muda e por quê. O plano revisado NÃO é editado — ele " +
				"continua sendo o registro do que se decidiu na época —, mas ganha no topo um aviso " +
				"apontando para cá, senão quem o lê fora de ordem segue uma decisão revista.",
			Body: "## O que este plano revisa\n\nRevisa `plans/00XX-nome.md`.\n\n" +
				"**O que muda:** TODO\n\n**Por quê:** TODO — o que se aprendeu depois de escrever aquele.\n\n" +
				"**O que continua valendo:** TODO — para quem leu o anterior saber o que não mudou.\n\n" +
				"> ⚠ **Falta um passo, e ele é no OUTRO arquivo.**\n" +
				">\n" +
				"> Declare `supersedes: plans/00XX-nome.md` no header DESTE plano, e escreva no\n" +
				"> **topo daquele**:\n" +
				">\n" +
				"> ```\n" +
				"> > **Revisado por** `plans/00YY-este.md` — <o que mudou e por quê>\n" +
				"> ```\n" +
				">\n" +
				"> Sem isso, quem abrir o plano antigo segue uma decisão que foi revista: ele\n" +
				"> continua parecendo coerente, porque É o registro coerente do que se decidiu na\n" +
				"> época. O gate `plano-revisado` reprova até o aviso existir — e apagar esta\n" +
				"> caixa depois de cumprir os dois passos.\n\n"},
	},
}

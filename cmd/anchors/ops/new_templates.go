package ops

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
	"product": productTemplate,
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
			Feeds:   []string{"route-declared"},
			Body:    "> **Rota**: `TODO`\n\n"},
		{Key: "overview", Title: "Visão geral (o que faz e para quem)", Default: true,
			Body: "## Visão Geral\nTODO: o que a unidade faz e para quem.\n\n"},
		{Key: "rules", Title: "Regras/efeitos catalogados (cada um com ID)", Default: true, Realizes: "B",
			Body: "## Regras\n\n### {id}-B01 — TODO regra\nDescreva o comportamento (não a implementação).\n\n"},
		{Key: "contract", Title: "Contrato de entrada/saída", Default: false,
			Purpose:  "Unidade com entrada e saída DISTINTAS que merecem descrição separada (regra de negócio, validador). Prefira `signature` quando a assinatura da função já diz tudo.",
			Variants: []string{"signature"},
			Feeds:    []string{"contract-status-declared"},
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
			Feeds:   []string{"domain-declared"},
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
		// A tabela de três colunas (`Estado | Quando | O que mostra`) DESCREVIA um estado;
		// não era confrontável. A forma abaixo é a do template que originou o conceito: por
		// estado, a CONDIÇÃO DE ENTRADA e a MATRIZ DE ELEMENTOS com as tags de visibilidade.
		//
		// A diferença é de natureza, não de volume. `[H]` (oculto) é uma asserção NEGATIVA
		// que um teste pode provar — "o botão de excluir não está visível para quem não é
		// dono" —, e a coluna "O que mostra" não sabe expressá-la: quem escreve lista o que
		// aparece e cala sobre o que não pode aparecer. É exatamente aí que mora o defeito
		// de permissão que passa por todos os gates.
		{Key: "states", Title: "Estados e transições (unidade de fluxo)", Default: false, Realizes: "S",
			Purpose: "Unidade com ESTADOS observáveis (tela, máquina de estado, hook com loading/erro/vazio). Não use em função pura.",
			Feeds:   []string{"scenario-asserts", "feature-test-match"},
			Body: "## Estados\n\n### {id}-S01 — TODO nome do estado\n\n" +
				"**Quando**: TODO: a condição de ENTRADA neste estado (sem ela o estado não é testável).\n\n" +
				"**Elementos** — `[V]` visível · `[D]` dinâmico · `[X]` desabilitado · `[H]` oculto · `[O]` opcional:\n\n" +
				"- `[V]` TODO: elemento sempre presente\n" +
				"- `[H]` TODO: o que NÃO pode aparecer neste estado\n\n"},
		// O DIAGRAMA das transições. Separado de `states` de propósito: aquela declara cada
		// estado isoladamente, esta declara o que LIGA um ao outro — e é onde aparecem os
		// estados inalcançáveis e os sem saída, que a lista por estado não revela.
		{Key: "state-flow", Title: "Fluxo de estados (as transições)", Default: false,
			Purpose: "Unidade com mais de um estado: o que leva de um a outro. Revela estado inalcançável e estado sem saída — que a lista estado-a-estado não mostra.",
			Feeds:   []string{"phase-ordered"},
			Body: "## Fluxo de Estados\n\n| De | Gatilho | Para |\n| --- | --- | --- |\n" +
				"| `{id}-S01` | TODO: o que dispara | `{id}-S02` |\n\n"},
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
			Feeds:   []string{"pagination-honored"},
			Purpose: "Unidade que lista uma COLECAO que cresce. Declare a estrategia: tudo de uma vez (e por que a colecao e limitada), paginada, ou scroll infinito — e qual peca a implementa.",
			Body:    "## Carregamento\n| Regra | Estrategia | Origem | Comportamento |\n| --- | --- | --- | --- |\n| `{id}-B01` | TODO: tudo de uma vez / paginada / scroll infinito | TODO: o hook ou consulta | TODO: se limitada, POR QUE; se paginada, o tamanho da pagina |\n\n"},
		{Key: "auth", Title: "Auth/Acesso", Default: false, Realizes: "R",
			Purpose: "Unidade cujo acesso depende de quem é o usuário (permissão, plano, dono do dado).",
			Body:    "## Auth/Acesso\nTODO: quem pode; regra de acesso.\n\n"},
		// ─── AS SEÇÕES DE CRUZAMENTO DA UNIDADE DE INTERFACE ──────────────────────────
		//
		// As sete seções abaixo voltaram ao catálogo depois de terem sido cortadas na
		// reescrita AGNÓSTICA do template (ver FINDINGS-spec-sections.md). O corte foi feito
		// pela APARÊNCIA de acoplamento: "Data Contract", "Navegação" e "Test IDs" parecem
		// específicos de app móvel, e por isso saíram junto com as tags de UI.
		//
		// Não eram. "Que dado aparece, de onde vem, em que formato, com que default" e "de
		// onde se chega a esta tela e para onde ela leva" são universais para QUALQUER
		// interface — web, desktop, TUI, voz. O que era específico do projeto de origem era
		// o EXEMPLO (React Native, Maestro), nunca a decisão.
		//
		// O custo do corte, medido num projeto que nasceu com o catálogo empobrecido: 85
		// specs colapsaram para SEIS títulos de seção, e estas sete aparecem ZERO vezes. O
		// conteúdo continuou denso e bom; o que se perdeu foi ONDE ele mora — e portanto o
		// que pode ser cruzado contra ele. Conteúdo fora de seção não é confrontável.
		//
		// A regra que sobrevive a isso: ao tornar algo agnóstico, a pergunta não é "isto
		// cita uma ferramenta?" e sim "isto é CRUZÁVEL?".
		{Key: "data-contract", Title: "Contrato de dados (todo dado dinâmico exibido)", Default: false,
			Purpose: "Unidade que EXIBE dado que vem de fora (consulta, parâmetro, store). Declara origem, obrigatoriedade, formato e default de cada campo. É onde se declara TIMEZONE de data e a tradução de enum — os dois defeitos de exibição que nenhum gate pega e que só aparecem em produção.",
			Feeds:   []string{"dependency-honored", "value-anchored"},
			Body: "## Contrato de Dados\n\n| Campo | Origem | Obrigatório | Formato | Default |\n" +
				"| --- | --- | --- | --- | --- |\n" +
				"| TODO | TODO: a consulta, o parâmetro ou o store | sempre / condicional / opcional | TODO: inclusive o FUSO, se for data | TODO: o valor quando ausente |\n\n" +
				"> Data sem fuso declarado e enum sem tradução declarada são os dois defeitos\n" +
				"> de exibição que atravessam todos os gates e só aparecem para o usuário.\n\n"},
		{Key: "data-states", Title: "Estados dos dados (um ramo por enum/condicional)", Default: false,
			Purpose: "Campo com valores enumerados ou origem condicional. Um item por VARIANTE — é o que garante um cenário por ramo em vez de um cenário para o caminho feliz. Sem isto, o enum com cinco valores vira um teste só.",
			Feeds:   []string{"scenario-asserts", "feature-test-match"},
			Body: "## Estados dos Dados\n\n### `TODO: campo`\n\n| Valor | Condição | O que a interface mostra |\n" +
				"| --- | --- | --- |\n| TODO | TODO | TODO |\n\n"},
		// A seção que o `route-declared` cobra. A doutrina do gate é explícita: com termos
		// genéricos ("Próxima tela", "Menu principal") a aresta de navegação não aponta para
		// lugar nenhum, e o grafo fica com nós soltos. Por isso o corpo insiste em NOME
		// CONCRETO — é a diferença entre uma aresta e uma frase.
		{Key: "navigation", Title: "Navegação (de onde se chega, para onde leva)", Default: false, Realizes: "N",
			Purpose: "Unidade NAVEGÁVEL: quem leva até ela e para onde ela leva, cada uma por NOME CONCRETO. É o que forma o grafo de navegação — termo genérico (\"a próxima tela\") não vira aresta.",
			Feeds:   []string{"route-declared", "route-exists"},
			Body: "## Navegação\n\n**Entrada** — quem leva até aqui:\n\n| Origem | Gatilho | Parâmetros |\n" +
				"| --- | --- | --- |\n| TODO: o NOME da unidade de origem | TODO | TODO |\n\n" +
				"**Saída** — para onde esta leva:\n\n| Destino | Gatilho | Condição |\n" +
				"| --- | --- | --- |\n| TODO: o NOME da unidade de destino | TODO | TODO |\n\n"},
		// O CONTRATO com o teste de ponta a ponta. Existe porque o identificador é a única
		// coisa que o teste E2E conhece da interface: renomeá-lo quebra o flow sem que
		// compilador, linter ou type-checker vejam nada. Declarado aqui, o rename fica
		// visível na spec — que é onde o revisor olha.
		{Key: "testids", Title: "Identificadores de teste (o contrato com o E2E)", Default: false,
			Purpose: "Unidade exercitada por teste de ponta a ponta. Os identificadores estáveis que o teste procura — renomear um quebra o flow e NENHUMA ferramenta de linguagem avisa. Declará-los aqui torna o rename visível em revisão.",
			Feeds:   []string{"testid-consistent", "testid-queried-exists"},
			Body: "## Identificadores de Teste\n\n| Identificador | Elemento | Usado em | Ação no teste |\n" +
				"| --- | --- | --- | --- |\n| TODO | TODO | TODO: o cenário | TODO |\n\n"},
		{Key: "components", Title: "Peças utilizadas (o que a unidade compõe)", Default: false,
			Purpose: "Unidade composta por outras já especificadas. Cada peça entre `crases` vira CONTRATO verificável; prosa não é cobrada. Diferente de `deps` (o que ela consome para funcionar): aqui é do que ela é FEITA.",
			Feeds:   []string{"dependency-honored", "sibling-guard"},
			Body: "## Peças Utilizadas\n\n| Peça | Papel | Origem |\n| --- | --- | --- |\n" +
				"| `TODO` | TODO | TODO |\n\n"},
		{Key: "a11y", Title: "Acessibilidade", Default: false,
			Purpose: "Unidade com formulário ou interação complexa: rótulo e dica de cada elemento, ordem de foco e comportamento de teclado. É requisito legal em boa parte dos contextos — e é decisão de produto, não detalhe de implementação.",
			Body: "## Acessibilidade\n\n| Elemento | Rótulo | Dica |\n| --- | --- | --- |\n| TODO | TODO | TODO |\n\n" +
				"**Foco inicial**: TODO.\n**Teclado**: TODO.\n\n"},
		// ─── A VARIAÇÃO POR ENTRADA (a unidade de COMPOSIÇÃO) ────────────────────────
		//
		// A distinção que organiza estas quatro: uma unidade de FLUXO varia por ESTADO
		// interno e navega; uma de COMPOSIÇÃO varia por ENTRADA (as propriedades que
		// recebe) e não navega — ela EMITE, e quem decide o que fazer é quem a compõe.
		//
		// O preset `component` não tinha nenhuma seção para isso: catalogava `states` e
		// `rules`, que são o vocabulário da unidade de fluxo. Ou seja, o preset da camada
		// de composição não capturava a ORIGEM DA VARIAÇÃO da própria camada — o defeito
		// era o mesmo do preset de tela, uma camada ao lado.
		//
		// Nota sobre `props` e TIPO: o template que originou o conceito escrevia a coluna
		// de tipo na sintaxe da linguagem (`'primary' | 'secondary'`). Aqui não: o que a
		// spec decide é o CONJUNTO DE VALORES aceitos, e ele sobrevive à troca de
		// linguagem — a sintaxe da união não.
		// Sem `Realizes`: a tabela de propriedades é um CONTRATO (o que entra), não um
		// catálogo de regras. Dar-lhe uma letra própria (`P`) fez o teste de letras
		// canônicas reprovar — e com razão: letra fora do vocabulário nasce invisível para
		// o `feature-test-match`, que é o pior tipo de furo (parece coberto e não está).
		{Key: "props", Title: "Propriedades (a interface pública da unidade)", Default: false,
			Purpose: "Unidade de COMPOSIÇÃO: tudo que entra e afeta o que ela mostra ou faz. Cada propriedade com os valores que aceita, se é obrigatória e o default. É o contrato com quem a compõe — e o análogo do `domain` para uma unidade de interface.",
			Feeds:   []string{"dependency-honored"},
			Body: "## Propriedades\n\n| Propriedade | Valores aceitos | Obrigatória | Default | Descrição |\n" +
				"| --- | --- | --- | --- | --- |\n" +
				"| TODO | TODO: o CONJUNTO de valores (não a sintaxe de tipo da linguagem) | sim / condicional / não | TODO | TODO |\n\n"},
		{Key: "variants", Title: "Variantes (a matriz de aparência e comportamento)", Default: false,
			Purpose: "Unidade com EIXOS de variação (aparência, tamanho, ênfase). Um quadro por eixo, com o que cada valor muda. Sem isto, a combinação de dois eixos nunca é enumerada — e é aí que mora a combinação que ninguém desenhou.",
			Feeds:   []string{"scenario-asserts"},
			Body: "## Variantes\n\n### `TODO: eixo de variação`\n\n| Valor | O que muda | Quando usar |\n" +
				"| --- | --- | --- |\n| TODO | TODO | TODO |\n\n"},
		// Substitui `navigation` na unidade de composição: ela não navega, emite. Manter as
		// duas separadas é o que impede o preset de componente de herdar a tabela de
		// navegação — que é falso-positivo garantido no `route-declared`.
		{Key: "callbacks", Title: "Eventos emitidos (o que a unidade avisa a quem a compõe)", Default: false, Realizes: "A",
			Purpose: "Unidade de COMPOSIÇÃO: o que ela emite para fora, com o gatilho e a CONDIÇÃO de disparo. É o equivalente da navegação para quem não navega — e a condição é a parte que se esquece (\"só dispara se não estiver desabilitado\").",
			Body: "## Eventos Emitidos\n\n| Evento | Dados | Gatilho | Condição de disparo |\n" +
				"| --- | --- | --- | --- |\n| `{id}-A01` | TODO | TODO | TODO: quando NÃO dispara |\n\n"},
		{Key: "slots", Title: "Composição (o que a unidade aceita dentro de si)", Default: false,
			Purpose: "Unidade que aceita CONTEÚDO de quem a compõe (filhos, regiões nomeadas). Declara quais regiões existem e o que cada uma espera. Omita se a unidade não aceita conteúdo externo.",
			Body: "## Composição\n\n| Região | O que aceita | Obrigatória |\n| --- | --- | --- |\n" +
				"| TODO | TODO | sim / não |\n\n"},
		// O par de `actions` na unidade de estado: aquela ESCREVE, esta LÊ. Sem declarar o
		// que se deriva, cada consumidor deriva por conta própria — e a mesma pergunta
		// passa a ter duas respostas que divergem quando a regra muda num lugar só.
		{Key: "selectors", Title: "Leituras derivadas (o que se calcula a partir do estado)", Default: false,
			Purpose: "Estado global: as leituras DERIVADAS que o consumidor usa em vez de recalcular. Declará-las aqui é o que impede a mesma derivação de ser reescrita em cada consumidor e divergir.",
			Body: "## Leituras Derivadas\n\n| Leitura | Devolve | Deriva de |\n| --- | --- | --- |\n" +
				"| `TODO` | TODO | TODO: os campos do estado que ela combina |\n\n"},
		{Key: "history", Title: "Histórico de alterações", Default: false,
			Purpose: "O rastro de quem mudou a spec e por quê. Vale onde a spec é contrato entre times e a mudança precisa ser atribuível.",
			Body:    "## Histórico de Alterações\n\n| Data | Autor | Alteração |\n| --- | --- | --- |\n| TODO | TODO | criação |\n\n"},
		{Key: "deps", Title: "Tabela de dependências", Default: false,
			Purpose: "O que a unidade consome. Símbolos entre `crases` viram CONTRATO verificável (gate dependency-honored); descrição em prosa não é cobrada.",
			Feeds:   []string{"dependency-honored"},
			Body:    "## Dependências\n| Cód | Arquivo | Método | Camada |\n| --- | --- | --- | --- |\n| DEP1 | TODO | `TODO` | TODO |\n\n"},
		// A seção da AMBIGUIDADE. Default: true — é a única seção cuja ausência esconde
		// justamente o que ela existe para revelar. Quem não tem dúvida gasta uma palavra
		// ("nenhuma"); quem tem, ganha um lugar declarado para ela em vez de um comentário
		// de PR que morre no merge.
		{Key: "open", Title: "Decisões em aberto (o que a spec ainda NÃO decide)", Default: true, Realizes: "Q",
			Feeds:   []string{"open-questions-resolved"},
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
		return fmt.Sprintf("\nTODO: a test case with `[%s]` in the NAME (it is what binds the "+
			"test to the spec requirement, through the JUnit that `anchors ingest` consumes).\n", code)
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
		Desc: "função PURA de backend — a assinatura é o contrato",
		// `constants` já existia no catálogo e NENHUM preset a emitia — seção que ninguém
		// emite é seção que ninguém escreve, e o limite de negócio vira número mágico.
		Sections: []string{"title", "overview", "signature", "domain", "effects", "invariants", "constants", "constraints", "deps", "open"},
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
		//
		// `data-contract`, `data-states`, `navigation` e `testids` VOLTARAM (ver
		// FINDINGS-spec-sections.md). São as seções de CRUZAMENTO da tela: sem elas o
		// `route-declared` fica cego, o enum de cinco valores vira um cenário só, e o
		// rename de um identificador de teste quebra o E2E sem que nada avise.
		//
		// Entram como DEFAULT pelo princípio que o parágrafo acima já mediu — a seção que
		// ninguém emite é a seção que ninguém escreve. A prova: num projeto nascido com o
		// preset SEM elas, 25 specs `layer: screen` foram escritas e 24 ficaram sem rota.
		// Quem não precisar de uma delas tira com `--without`; o caro é o contrário.
		Sections: []string{"title", "route", "overview", "states", "state-flow", "loading", "rules", "data-contract", "data-states", "messages", "navigation", "auth", "components", "deps", "testids", "a11y", "notes", "open"},
	},
	"component": {
		Desc: "componente de UI — props e estados visuais, sem rota",
		// Sem `route` e sem `navigation`: componente não é navegável — ele emite callback,
		// e quem navega é a tela que o compõe. Mas os identificadores de teste e a
		// acessibilidade valem igual: o componente é exercitado pelo E2E da tela, e é nele
		// que o rótulo e a dica de fato moram.
		// `props` entra no lugar de `contract`: são o mesmo papel (o que entra) no dialeto
		// da unidade de composição — e `props` cataloga valor a valor, que é o que a
		// matriz de `variants` depois cruza.
		//
		// `callbacks` é o que substitui `navigation` aqui: componente não navega, emite.
		Sections: []string{"title", "overview", "props", "variants", "states", "callbacks", "slots", "rules", "data-states", "messages", "components", "testids", "a11y", "open"},
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
		Desc: "estado global (Redux/Zustand/Pinia/MobX) — shape, actions e hidratação",
		// `selectors` entre `actions` e `hydration`: escreve, lê, persiste — a ordem em que
		// se pensa um estado global.
		Sections: []string{"title", "overview", "state-shape", "actions", "selectors", "hydration", "invariants", "constraints", "deps", "open"},
	},
	"validation": {
		Desc:     "regra de validação — critérios e a mensagem que o usuário lê",
		Sections: []string{"title", "overview", "contract", "domain", "rules", "constants", "messages", "deps", "open"},
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
				"e `needs:` (de qual depende). O PROGRESSO não mora aqui: ele vive no " +
				"`-progress.md` ao lado, porque o plano é DECISÃO e alterá-lo tem de significar " +
				"que a decisão mudou.",
			Body: "## Fases\n\n### {id}-F01 — TODO nome da primeira fase\n\n" +
				"TODO: o que esta fase entrega, e as specs que ela semeia.\n\n" +
				"- `caminho/Unidade.spec.md` — TODO o que ela descreve\n\n" +
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
			Purpose: "Só quando este plano REVISA outro. Declare `revises: plans/00XX-nome.md` no " +
				"header, e escreva aqui o que muda e por quê. O plano revisado NÃO é editado — ele " +
				"continua sendo o registro do que se decidiu na época —, mas ganha no topo um aviso " +
				"apontando para cá, senão quem o lê fora de ordem segue uma decisão revista.",
			Body: "## O que este plano revisa\n\nRevisa `plans/00XX-nome.md`.\n\n" +
				"**O que muda:** TODO\n\n**Por quê:** TODO — o que se aprendeu depois de escrever aquele.\n\n" +
				"**O que continua valendo:** TODO — para quem leu o anterior saber o que não mudou.\n\n" +
				"> ⚠ **Falta um passo, e ele é no OUTRO arquivo.**\n" +
				">\n" +
				"> Declare `revises: plans/00XX-nome.md` no header DESTE plano, e escreva no\n" +
				"> **topo daquele**:\n" +
				">\n" +
				"> ```\n" +
				"> > [!IMPORTANT]\n" +
				"> > Revisado por `plans/00YY-este.md` — <o que mudou e por quê>\n" +
				"> >\n" +
				"> > (ou `@revised-by: plans/00YY-este.md`, a forma buscável por grep)\n" +
				"> ```\n" +
				">\n" +
				"> Sem isso, quem abrir o plano antigo segue uma decisão que foi revista: ele\n" +
				"> continua parecendo coerente, porque É o registro coerente do que se decidiu na\n" +
				"> época. O gate `plano-revisado` reprova até o aviso existir.\n" +
				">\n" +
				"> **E marque as PARTES atingidas, não só o topo.** O aviso de topo diz que o\n" +
				"> plano mudou; ele não diz ONDE. Quem lê a fase 3 não sabe se ela é uma das que\n" +
				"> mudaram — e o custo de descobrir é reler o plano inteiro procurando.\n" +
				">\n" +
				"> A marcação depende do que já aconteceu:\n" +
				">\n" +
				"> - **parte JÁ IMPLEMENTADA** — o texto FICA como está: ele descreve o que foi\n" +
				">   feito, e reescrevê-lo faria o registro mentir sobre o passado. Acrescente\n" +
				">   abaixo dele:\n" +
				">\n" +
				">   ```\n" +
				">   > [!WARNING]\n" +
				">   > Alterado por `plans/00YY-este.md` — <o que muda daqui em diante>.\n" +
				">   > O texto acima descreve o que FOI implementado.\n" +
				">   ```\n" +
				">\n" +
				"> - **parte AINDA NÃO implementada** — reescreva o texto com o comportamento\n" +
				">   novo, e preserve o antigo ao lado:\n" +
				">\n" +
				">   ```\n" +
				">   > [!WARNING]\n" +
				">   > Alterado por `plans/00YY-este.md`.\n" +
				">   > **Era:** <o texto original>\n" +
				">   > **Por quê:** <o que se aprendeu>\n" +
				">   ```\n" +
				">\n" +
				"> A diferença é o que o registro precisa preservar: no primeiro caso, o que se\n" +
				"> fez; no segundo, o que se pretendia fazer e por que mudou.\n" +
				">\n" +
				"> Apague esta caixa depois de cumprir os passos.\n\n"},
	},
}

// productTemplate — a DOUTRINA DE PRODUTO: a regra que atravessa alvos.
//
// Poucas seções de propósito. A doutrina não é uma spec: ela não tem alvo, não tem
// contrato, não tem estado. Ela decide, e diz por quê — e cada seção a mais é um convite
// a preencher formulário em vez de escrever a decisão.
var productTemplate = template{
	kind:     "product",
	ext:      ".doctrine.md",
	idField:  "code",
	headerFn: productHeader,
	// As chaves levam o prefixo `doctrine_` de proposito.
	//
	// `sectionBody` resolve o corpo por `section.body.<chave>` no catalogo de traducao, e
	// o literal do template so' vale quando a chave NAO existe la'. Com `rules` ou
	// `overview` nus, a doutrina herdava silenciosamente o corpo da SPEC — medido: o
	// `anchors new product` emitia `### CRLMC-B01 — TODO rule` e "o que a unidade faz",
	// texto de spec num artefato que nao tem unidade nenhuma.
	sections: []section{
		{Key: "doctrine_title", Title: "section.title.title", Default: true,
			Body: "# {name}\n\n> **Code**: `{id}`\n\n"},

		{Key: "doctrine_overview", Title: "section.title.overview", Default: true,
			Purpose: "section.purpose.doctrine_overview",
			Body:    "## Overview\n\nTODO: what this doctrine decides, and what it costs when it is broken.\n\n"},

		{Key: "doctrine_rules", Title: "section.title.rules", Default: true,
			Purpose: "section.purpose.doctrine_rules",
			Body: "## Rules\n\n### {id}-R01 — TODO: the rule, stated so a spec can realize it\n\n" +
				"### {id}-R02 — TODO\n\n"},

		{Key: "doctrine_constraints", Title: "section.title.constraints", Default: false,
			Purpose: "section.purpose.doctrine_constraints",
			Body:    "## Constraints\n\n| Rule | Boundary | Why |\n| --- | --- | --- |\n| `{id}-X01` | TODO | TODO |\n\n"},

		{Key: "doctrine_open", Title: "section.title.open", Default: true,
			Purpose: "section.purpose.open",
			Body:    "## Open Decisions\n\n| Code | Question | Who decides | Becomes |\n| --- | --- | --- | --- |\n\nnone\n\n"},
	},
}

// productHeader — a doutrina NÃO declara `layer:`.
//
// A camada é do alvo que uma spec descreve, e a doutrina não tem alvo: ela é o artefato.
// O `mdHeader` emite `layer: TODO`, que aqui seria um campo que ninguém pode preencher —
// e um placeholder eterno é exatamente o que o `placeholder-filled` existe para acusar.
func productHeader(id, outPath string) string {
	return fmt.Sprintf("<!-- @anchors\n  code: %s\n  updated_at: TODO\n-->\n", id)
}

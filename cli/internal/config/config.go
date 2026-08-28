package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultFile é o nome do arquivo de configuração geral do Anchors na raiz do projeto.
const DefaultFile = "anchors.yaml"

// Config é a configuração geral do projeto Anchors. A primeira seção é a Estrutura
// de Projeto (o grafo virtual, STRUCTURE §2.1); gates, providers etc. crescem aqui
// depois. Ver DECISIONS.md D1.
type Config struct {
	Version  int                 `yaml:"version"`
	Comments map[string][]string `yaml:"comments,omitempty"` // override/extensão dos marcadores (D4)
	Layers   map[string]Layer    `yaml:"layers"`             // as camadas (Estrutura)
	Derived  *Derived            `yaml:"derived,omitempty"`  // co-location dos derivados
	Governs  []GovernRule        `yaml:"governs,omitempty"`  // dimensão vertical (arestas de alto grau)
	Gates    []Gate              `yaml:"gates,omitempty"`    // os gates de qualidade (QUALITY §3-§5)
	Recode   *Recode             `yaml:"recode,omitempty"`   // convenções de projeto p/ `anchors recode`
	// Tests e Mutation declaram COMO este projeto produz sinal de teste: o comando é
	// do projeto, a amarração ao mapa é do Anchors. Ver Suite.
	Tests    []Suite `yaml:"tests,omitempty"`
	Mutation []Suite `yaml:"mutation,omitempty"`
	// RuleTypes declara o VOCABULÁRIO de tipos de regra do projeto: cada letra do
	// código (`{CODE}-<letra><NN>`) é a inicial do TERMO que nomeia a seção da spec.
	// Vazio → o engine usa as letras canônicas (ver DefaultRuleLetters).
	RuleTypes []RuleType `yaml:"rule_types,omitempty"`
	// CodeLengths: os comprimentos de código de identidade que ESTE projeto usa
	// (`code_lengths: [4]`, `[5]`, ou `[4, 5]` durante uma migração). Vazio → o default
	// do engine, que aceita 4 e 5. Ver CodeLengths para o porquê de ser do projeto.
	CodeLengths []int `yaml:"code_lengths,omitempty"`
	// SectionTitles: o léxico do projeto para as seções de spec, por chave do catálogo
	// (`contract: "Modelo de Dado"`). Sem isto, o `anchors new` emite os títulos do
	// framework e as specs nascem em dialeto diferente do dos vizinhos.
	SectionTitles SectionTitles `yaml:"section_titles,omitempty"`
	// RouteRegistry: onde o projeto REGISTRA suas rotas de navegação (globs). O Anchors
	// não pode adivinhar — cada stack tem o seu lugar e a sua sintaxe. Sem isto, o gate
	// `route-exists` fica Pendente em vez de afirmar que uma rota existe sem ter olhado.
	RouteRegistryGlobs []string `yaml:"route_registry,omitempty"`
	// Obligations são as OBRIGAÇÕES TRANSVERSAIS do projeto: "todo nó que carrega o
	// atributo P deve aparecer em Q". Ver Obligation.
	Obligations []Obligation `yaml:"obligations,omitempty"`
	// Dialect é o LÉXICO da linguagem do projeto — o que permite a um gate agnóstico
	// reconhecer uma função/consulta no código concreto. Ver dialect.go.
	Dialect *Dialect `yaml:"dialect,omitempty"`
	// Boundaries são as FRONTEIRAS entre camadas: o que cada uma não pode alcançar.
	// Declarar `layers:` sem isto é desenhar a arquitetura sem defendê-la. Ver Boundary.
	Boundaries []Boundary `yaml:"boundaries,omitempty"`
	// Packs são conjuntos de obrigações DISTRIBUÍVEIS (LGPD, GDPR, PCI-DSS, WCAG…) que o
	// projeto adota em vez de reescrever. Ver internal/pack.
	//
	// Sem eles, cada projeto declarava os mesmos deveres do zero, com vocabulário próprio
	// — e não havia como responder "quais deveres se aplicam a mim". Um pack traz o dever
	// e a norma que o origina; o PROJETO traz onde isso vive nele (`pack_values`).
	Packs []string `yaml:"packs,omitempty"`
	// PackValues resolve os placeholders que os packs usam: `{{purge_handler}}` vira o
	// caminho real deste projeto. É a divisão que impede o pack de adivinhar a estrutura
	// alheia — o dever é universal, o endereço é local.
	PackValues map[string]string `yaml:"pack_values,omitempty"`
	// AutoJudgment permite que a IA DECIDA sozinha as sugestões de julgamento, em vez de
	// deixá-las aguardando aprovação humana.
	//
	// Desligado por default, e o default é a decisão importante: aprovar automaticamente
	// significa que uma correção entra no projeto sem ninguém ter olhado. Isso é
	// aceitável em algumas casas (repositório pessoal, projeto experimental) e
	// inaceitável na maioria — e o silêncio tem de cair no lado seguro.
	//
	// Quando ligado, a decisão fica MARCADA como automática no arquivo da sugestão
	// (`por: IA (auto_judgment)`). Apagar essa distinção faria "ninguém olhou isto"
	// parecer "alguém aprovou", que é a informação mais cara de perder numa auditoria.
	AutoJudgment bool `yaml:"auto_judgment,omitempty"`

	// Jurisdictions são os mercados onde a aplicação opera (`br`, `eu`, `us-ca`). Filtram
	// quais packs valem: um pack de jurisdição não declarada não é carregado, e o aviso
	// aparece — declarar GDPR num app que só opera no Brasil costuma ser engano, e o
	// silêncio esconderia isso.
	Jurisdictions []string `yaml:"jurisdictions,omitempty"`

	// Workflow — ONDE a fila de trabalho mora. Ver o tipo.
	Workflow *Workflow `yaml:"workflow,omitempty"`
}

// Workflow declara o modo de gestão do trabalho, e é EXCLUDENTE de propósito.
//
//	workflow:
//	  mode: local     # a fila é .anchors/tasks/ (o watcher enfileira, `anchors next` puxa)
//	  mode: github    # a fila são as issues/cards do repositório
//
// Um modo OU outro, nunca um com o outro de reserva. A alternativa — tentar o GitHub e
// cair no local quando falha — parece robustez e é a pior escolha possível: passa a
// existir a pergunta "de qual fila veio esta task?", que ninguém consegue responder depois
// do fato, e dois agentes podem trabalhar no mesmo item achando que reivindicaram.
//
// É a mesma armadilha de um util de login que "às vezes loga, às vezes herda a sessão":
// enquanto o comportamento é condicional e implícito, o sintoma aparece longe da causa.
// Declarar o modo faz o próprio config responder, e faz as validações serem simples: no
// modo github, `.anchors/tasks/` não deve existir; no modo local, nenhum comando toca a
// rede.
type Workflow struct {
	// Mode: "local" | "github". Vazio = local (o comportamento que sempre existiu, para
	// não quebrar projeto que nunca declarou nada).
	Mode string `yaml:"mode"`

	// Repo no formato "owner/nome". Só no modo github, e OBRIGATÓRIO nele: inferir do
	// remote do git faria o Anchors escrever em outro repositório quando alguém trabalha
	// num fork, e escrita em lugar errado é o erro que não se desfaz com um revert.
	Repo string `yaml:"repo,omitempty"`

	// Labels que marcam uma issue como trabalho DO ANCHORS. Sem isto, `anchors next`
	// puxaria qualquer issue do repositório — inclusive as de produto, que não têm a
	// forma que o ciclo espera. Vazio no modo github é erro de configuração, não default.
	Labels []string `yaml:"labels,omitempty"`
}

// ModoGitHub diz se o projeto declarou a gestão no GitHub.
func (c *Config) ModoGitHub() bool {
	return c != nil && c.Workflow != nil && c.Workflow.Mode == ModeGitHub
}

const (
	ModeLocal  = "local"
	ModeGitHub = "github"
)

// Boundary — uma fronteira de camada: "arquivos DESTA camada não podem conter ISTO".
//
// É o mecanismo por trás de toda regra arquitetural de import ("tela não fala com
// repositório", "model não importa service") e também das proibições de padrão ("sem cor
// literal", "sem relógio cru"). São a mesma forma, e por isso um só mecanismo: o projeto
// declara O QUE não pode, o engine confronta ONDE.
//
// O engine não sabe o que é uma "tela" nem o que é um "repositório" — quem sabe é o
// projeto, que já declarou suas camadas em `layers:`. Aqui ele só diz o que cada uma
// não alcança.
type Boundary struct {
	// Layer é a camada a quem a regra se aplica (o nome em `layers:`). VAZIO = vale para
	// todo código — é assim que se declara uma proibição global.
	Layer string `yaml:"layer,omitempty"`
	// Forbid é o padrão (regex) que não pode aparecer. Quem o escreve é o projeto: é aqui
	// que mora o dialeto de import da linguagem, que o engine não tem por que conhecer.
	Forbid string `yaml:"forbid"`
	// Because é o PORQUÊ da regra — vai na mensagem do gate. Uma proibição sem motivo
	// vira ritual: quem esbarra nela não sabe se está corrigindo um erro ou contornando
	// uma regra arbitrária.
	Because string `yaml:"because,omitempty"`
	// Severity: "error" (default, reprova) ou "warn" (registra sem reprovar). É a
	// maturação (QUALITY §7) POR REGRA — permite travar a fronteira nova sem desligar o
	// gate por causa do backlog da fronteira antiga.
	Severity string `yaml:"severity,omitempty"`
}

// Obligation — uma obrigação TRANSVERSAL: o dever que um artefato contrai com um lugar
// que ele NÃO conhece. É a classe de erro que nenhuma spec-por-unidade cobre, porque a
// obrigação mora FORA da unidade: a spec de um modelo de dado não tem por que falar do
// script de exclusão de conta, e quem escreveu o script não sabe que nasceu um modelo.
//
// Instâncias em qualquer projeto:
//   - LGPD/GDPR: modelo que carrega dado pessoal → precisa estar no purge e no export
//   - i18n:      string exibida ao usuário       → precisa estar no arquivo de traduções
//   - a11y:      componente interativo           → precisa de rótulo acessível
//   - auditoria: operação financeira             → precisa emitir trilha
//
// O ANCHORS não sabe o que é LGPD. Sabe que existem obrigações, e oferece como
// declará-las e cobrá-las. QUAIS obrigações, e o vocabulário de atributos, é do PROJETO.
//
// Por que o gatilho é um ATRIBUTO e não uma camada: deduzir por camada produz
// falso-positivo em massa. Medido num projeto real: 38 de 50 modelos não estavam no
// script de purge, e a MAIORIA legitimamente (dado compartilhado, dado com TTL,
// métrica agregada). Um gate que erra 76% das vezes é desligado no primeiro dia.
type Obligation struct {
	Name string `yaml:"name"` // identificador da obrigação (ex.: "pii-purgavel")
	// When: o ATRIBUTO que dispara a obrigação, como aparece no header `@anchors` do
	// nó (ex.: "carries: pii"). Vocabulário do projeto.
	When string `yaml:"when"`
	// MustAppearIn: os arquivos onde o nó DEVE aparecer (globs). É o dever.
	MustAppearIn []string `yaml:"must_appear_in"`
	// IdentifiedBy: como o nome do nó aparece no destino. O nome do artefato e a forma
	// como ele é referenciado raramente coincidem (`MetadataEntry` vira
	// `METADATA_ENTRIES_TABLE_NAME`). Valores: "as-is" (default), "screaming-snake",
	// "snake", "kebab", ou um template com {{NAME}} (ex.: "{{NAME}}_TABLE_NAME").
	IdentifiedBy string `yaml:"identified_by,omitempty"`
	// IdentifiedAsForm SOBREPÕE o `identified_as` declarado no NÓ, para esta obrigação.
	// Existe porque deveres diferentes procuram FORMAS diferentes do mesmo artefato: o
	// dever de LGPD procura a env var (`METADATA_ENTRY_TABLE_NAME`) nos handlers; o de
	// provisionamento procura o nome do modelo (`tables['MetadataEntry']`) na infra.
	// Sem isso, um dos dois passa a procurar o token errado e acusa quem está correto.
	IdentifiedAsForm string `yaml:"identified_as_form,omitempty"`
	// Because: POR QUE a obrigação existe. Não é decorativo — entra na mensagem do
	// gate. "Falta em Q" é ignorado; "LGPD: dado pessoal precisa ser apagável" é
	// obedecido.
	Because string `yaml:"because,omitempty"`
}

// SectionTitles renomeia as seções do catálogo para o LÉXICO do projeto, por chave
// estável (`contract`, `domain`, `rules`…). A chave é do framework; o título é do
// projeto.
//
// Existe porque o catálogo emitia títulos fixos ("Contrato", "Domínio", "Efeitos") e um
// projeto real usava outros ("Modelo de Dado", "Comportamentos", "Notas de
// Implementação"). Medido: das seções que o preset `schema` emitia, NENHUMA das três
// centrais coincidia com as 50 specs vizinhas do mesmo projeto. O `anchors new` e os
// vizinhos ensinavam formatos diferentes, e dois agentes escreveram specs em dialetos
// opostos — cada um obedecendo a uma fonte, os dois passando nos gates.
//
// Título de seção é CONTEÚDO, não mecanismo: o Anchors não tem por que impô-lo, do mesmo
// jeito que não impõe idioma nem nomes de vendor (ver `dialect`).
type SectionTitles map[string]string

// TituloDaSecao devolve o título que o projeto usa para uma chave do catálogo, ou o
// padrão do framework quando ninguém renomeia.
//
// A precedência é CAMADA > projeto > framework, e a camada existe porque o dialeto quase
// nunca é global. Medido num projeto real: "Modelo de Dado" e "Comportamentos" aparecem
// em 50 specs de uma única camada e em 1 de todas as outras 588 — renomear no nível do
// projeto teria trocado o título de 588 specs que estavam certas para consertar 50.
func (c *Config) TituloDaSecao(chave, padrao, camada string) string {
	if c == nil {
		return padrao
	}
	if camada != "" {
		if l, ok := c.Layers[camada]; ok {
			if t := strings.TrimSpace(l.SectionTitles[chave]); t != "" {
				return t
			}
		}
	}
	if t := strings.TrimSpace(c.SectionTitles[chave]); t != "" {
		return t
	}
	return padrao
}

// RuleType liga uma LETRA de código ao termo que a origina e às seções de spec que a
// usam. O vocabulário é extensível por projeto; o que não pode é CONFLITO — duas seções
// reivindicando a mesma letra, ou um código usando letra não declarada (essa regra fica
// invisível para a rastreabilidade).
type RuleType struct {
	Letter   string   `yaml:"letter"`             // a letra (1 char, maiúsculo)
	Term     string   `yaml:"term"`               // o termo que a origina (ex.: "State")
	Sections []string `yaml:"sections,omitempty"` // títulos de seção que a usam
	// RequiresCode: as seções (deste subconjunto de `Sections`) em que uma tabela
	// preenchida SEM código é achado.
	//
	// Nem toda seção declarada cataloga regra. Medido num projeto real: "Eventos /
	// Callbacks" afirma comportamento ("`onPress` — Tap na linha") e precisa de
	// código para o cenário citar; "Variantes" apenas ENUMERA o que existe
	// (`individual`, `familia` e de onde vêm seus rótulos) e não afirma nada
	// verificável. As duas estão declaradas sob letras, e só a primeira deve cobrar.
	//
	// Fica na config, e não numa heurística do gate, porque a distinção é semântica
	// do PROJETO: o mesmo título pode enumerar num repositório e normatizar noutro.
	RequiresCode []string `yaml:"requires_code,omitempty"`

	// Tags: as tags de CENÁRIO que declaram esta natureza (`@estado` para a letra S).
	//
	// A letra vive no código (`ABCDX-S01`) e a tag vive na linha acima do cenário. As
	// duas dizem a mesma coisa, e nada as confrontava: medido num projeto real, 36
	// cenários afirmavam `@comportamento` sob código `-S` (e vice-versa) — o texto
	// dizia "renderiza em itálico", que é estado, e a tag dizia comportamento.
	//
	// Sem o mapa, o gate `cenario-tipo-alinhado` não tem como saber que `@estado`
	// corresponde a S neste projeto: a tag é vocabulário do repositório, não do
	// framework (um projeto em inglês escreveria `@state`).
	Tags []string `yaml:"tags,omitempty"`
}

// LetrasDaTag devolve TODAS as letras que declaram a tag de cenário, e se a tag é
// conhecida. Tags fora do vocabulário não são erro: o projeto usa `@smoke`, `@P1`,
// `@nivel-e2e` e outras que não falam de natureza de regra.
//
// São várias, e não uma, porque a mesma tag pode caber sob mais de uma natureza sem
// ambiguidade. Medido num projeto real: `@estado-dado` marca tanto cenários de `S`
// ("a tela exibe X quando o dado é Y") quanto de `V`, cuja seção "Validações Visuais"
// cataloga aparência-por-prop — a seção classifica pelo que a spec normatiza, a tag
// pelo que o cenário observa, e as duas leituras são verdadeiras ao mesmo tempo.
//
// Devolver só a primeira faria o gate acusar a segunda para sempre, e a saída seria
// escolher entre duas classificações corretas.
func (c *Config) LetrasDaTag(tag string) ([]string, bool) {
	if c == nil {
		return nil, false
	}
	var out []string
	for _, rt := range c.RuleTypes {
		for _, t := range rt.Tags {
			if normalizaTitulo(t) == normalizaTitulo(tag) {
				out = append(out, strings.ToUpper(strings.TrimSpace(rt.Letter)))
			}
		}
	}
	return out, len(out) > 0
}

// ExigeCodigo diz se a seção (pelo título) foi declarada como catalogadora de regra.
func (r RuleType) ExigeCodigo(titulo string) bool {
	for _, s := range r.RequiresCode {
		if normalizaTitulo(s) == normalizaTitulo(titulo) {
			return true
		}
	}
	return false
}

// normalizaTitulo compara títulos ignorando caixa e espaço de borda — o mesmo
// critério que o gate `rule-types` usa para casar seção com letra.
func normalizaTitulo(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// DefaultRuleLetters são as letras canônicas, usadas quando o projeto não declara
// `rule_types`. Mantém compatibilidade com projetos anteriores ao registro.
var DefaultRuleLetters = "SRVAXBNMD"

// CodeLengths são os comprimentos de código de identidade que o engine reconhece.
//
// Era FIXO — e trocar o número fixo de 4 para 5 (ou o contrário) só move o problema de
// lugar: o engine passa a servir apenas os projetos que escolheram o mesmo tamanho. Medido:
// o app de referência tem **4046 códigos de 4 caracteres** e zero de 5; com o engine em `{5}`, nenhum
// deles é reconhecido — as trincas ficam invisíveis e os gates de identidade não têm o que
// confrontar, sem que nada acuse a causa.
//
// O comprimento é do PROJETO, como o vocabulário de letras (`rule_types`) já é. Um projeto
// que declara 4 continua válido; um que declara 5 também; e um que aceita os dois durante
// uma migração — que é o caso real de quem já tem specs escritas — declara ambos.
//
// Default `[5]`: é o comprimento canônico do framework, e a escolha é do Adriel — o
// Anchors FORÇA 5 em si mesmo. Um projeto que usa outro tamanho declara
// (`code_lengths: [4]`), e quem está migrando declara os dois (`[4, 5]`), que é o caso do
// app de referência: 4046 códigos de 4 caracteres escritos, migrando para 5.
//
// O default é restritivo de propósito. Ele governa o projeto NOVO, que ainda não tem código
// escrito e por isso pode nascer no formato canônico — e um default permissivo (`[4,5]`)
// faria todo projeto novo aceitar dois formatos para sempre, sem nunca ter decidido qual
// usa. Quem tem código legado declara e o engine obedece.
var CodeLengths = []int{5}

// SetCodeLengths reconfigura os comprimentos aceitos. Chamada na porta de entrada (o Walk),
// como o SetRuleLetters — os pacotes que leem código precisam concordar sobre o que É um
// código, senão cada um reconhece um conjunto diferente e a cola se parte em silêncio.
func SetCodeLengths(ls []int) {
	if len(ls) > 0 {
		CodeLengths = ls
	}
}

// CodeLengthPattern devolve o trecho de regex que casa os comprimentos aceitos —
// `{4}` para um só, `{4,5}` para uma faixa contígua, `{4}|{6}` alternado para o resto.
//
// Existe para que os 25 regexes espalhados pelo engine deixem de repetir o número: cada
// cópia era um lugar a esquecer numa mudança, e foi exatamente o que aconteceu — a troca de
// 4 para 5 pegou uns arquivos e não outros.
func CodeLengthPattern() string {
	ls := append([]int{}, CodeLengths...)
	sort.Ints(ls)
	if len(ls) == 1 {
		return fmt.Sprintf("{%d}", ls[0])
	}
	if ls[len(ls)-1]-ls[0] == len(ls)-1 {
		return fmt.Sprintf("{%d,%d}", ls[0], ls[len(ls)-1])
	}
	partes := make([]string, len(ls))
	for i, l := range ls {
		partes[i] = fmt.Sprintf("[A-Z0-9]{%d}", l)
	}
	return "(?:" + strings.Join(partes, "|") + ")"
}

// RuleLetters devolve a classe de letras válidas do projeto (para compor os regexes de
// código de cenário). Cai nas canônicas se o projeto não declarar o vocabulário.
func (c *Config) RuleLetters() string {
	if c == nil || len(c.RuleTypes) == 0 {
		return DefaultRuleLetters
	}
	var b strings.Builder
	seen := map[string]bool{}
	for _, rt := range c.RuleTypes {
		l := strings.ToUpper(strings.TrimSpace(rt.Letter))
		if len(l) != 1 || seen[l] {
			continue
		}
		seen[l] = true
		b.WriteString(l)
	}
	if b.Len() == 0 {
		return DefaultRuleLetters
	}
	return b.String()
}

// Recode declara o DIALETO de projeto que o `anchors recode` usa além do genérico
// (header/scenario-codes/refs). É opcional — sem ele, o recode só faz o genérico.
type Recode struct {
	// TestID: como o prefixo de testID DERIVA do código. Convenção ideal do projeto
	// (o testID é a projeção do código na borda observável). Valores: "lower" (código
	// em minúsculo: TCDT→tcdt), "kebab" (reservado), ou vazio (não deriva testID).
	// Se o prefixo derivado não existir mas o arquivo tiver testIDs (prefixo divergente,
	// de um recode manual anterior), o recode AVISA e NÃO os toca — não adivinha.
	TestID string `yaml:"test_id,omitempty"`
	// FilePatterns: globs de NOMES de arquivo que contêm o código, com a variável
	// {{code}} (ex.: "**/{{code}}-*.yaml", "**/*.{{code}}-*.png"). O recode renomeia
	// esses arquivos (git mv) ao aplicar, e atualiza as referências a eles.
	FilePatterns []string `yaml:"file_patterns,omitempty"`
}

// Gate é um gate de qualidade declarado (QUALITY §4: "ligar um gate é declarativo").
// Confronta os alvos que casam `on` contra um critério. Dois modos de execução:
//   - EXTERNO: `run` é um comando (jest, eslint, tsc…); exit 0 = passou. O CLI o
//     invoca e lê o exit code — não reimplementa a ferramenta.
//   - INTERNO: `check` nomeia uma verificação embutida do CLI (ex.: "spec-sections",
//     "scenario-covered") — o CLI a executa lendo texto.
type Gate struct {
	Name string `yaml:"name"` // identificador do gate
	// ID é o identificador ESTÁVEL da regra que este gate verifica.
	//
	// Diferente de `Name`, que é o rótulo legível e pode ser renomeado, o ID é o que
	// outros lugares citam: o relatório o imprime, e é por ele que se dispensa UMA
	// verificação sem descartar as demais.
	//
	// Existe porque um gate cobra mais de uma coisa — `spec-completa` cobra "sem
	// placeholder" E "tem regra catalogada" — e sem ID as duas eram indistinguíveis:
	// quem quisesse dispensar a primeira perdia a segunda junto.
	//
	// Os gates CANÔNICOS têm ID fixo (ver o catálogo em initx). Um gate declarado pelo
	// projeto precisa de um ID ÚNICO — a validação o cobra na carga, porque dois gates
	// com o mesmo ID fazem toda citação apontar para dois lugares.
	ID   string   `yaml:"id,omitempty"`
	On   []string `yaml:"on"`             // kinds de nó que este gate confronta (spec, code, feature…)
	Tags []string `yaml:"tags,omitempty"` // se presente, restringe aos nós com QUALQUER destas tags
	// ExcludeTags: o filtro NEGATIVO que faltava. `tags` diz a quem o gate se aplica; não
	// havia como dizer a quem ele NÃO se aplica, e a diferença aparece quando a exceção é
	// uma camada inteira no meio de um universo grande — listar todas as outras tags no
	// `tags` positivo é uma lista que envelhece a cada camada nova.
	//
	// O caso que motivou: `mutation-score` cobrado de camadas que a ferramenta de mutação
	// nem consegue rodar (declaração de infra, schema) — nenhum teste as alcança, então o
	// Stryker sai com "No tests were found" e não há relatório para ingerir. Sem sinal, o
	// gate ficava Pending eterno pedindo uma medição impossível de produzir.
	//
	// Um nó com QUALQUER destas tags fica fora do gate. Vazio = sem exclusão, para que o
	// silêncio nunca desligue gate existente.
	ExcludeTags []string `yaml:"exclude_tags,omitempty"`
	Run         string   `yaml:"run,omitempty"`   // comando externo (modo externo)
	Check       string   `yaml:"check,omitempty"` // verificação interna nomeada (modo interno)
	// Blocking: bloqueante vs. informativo (maturação, QUALITY §7).
	//
	// PONTEIRO de propósito. Como `bool`, o campo omitido e o `blocking: false`
	// explícito chegavam idênticos ao parser — e o merge canônico não tinha como
	// distinguir "não declarei" de "decidi que é informativo". A consequência seria
	// grave na direção errada: um gate mantido informativo DE PROPÓSITO (porque a
	// unidade ainda tem débito) seria promovido a bloqueante ao herdar do canônico, e
	// passaria a barrar commits que o autor havia liberado conscientemente.
	//
	// Com `*bool`, `nil` é exatamente "a chave não veio no yaml" e o merge funciona
	// como para todos os outros campos. Leia sempre por `IsBlocking()`.
	Blocking *bool `yaml:"blocking,omitempty"`
	// Measures declara o TIPO de medidor (QUALITY §5.2). Vazio/qualquer descrição =
	// determinístico (usa run/check). O valor especial "judgment" marca o GATE DE
	// JULGAMENTO POR IA: o CLI não computa — enfileira o alvo para uma IA confrontar
	// e reportar o veredito via `anchors judge`. Mesma dupla saída (carimbo + issue).
	Measures string `yaml:"measures,omitempty"`
	Guide    string `yaml:"guide,omitempty"` // (judgment) o arquivo-régua que a IA deve ler
	Ask      string `yaml:"ask,omitempty"`   // (judgment) a pergunta específica do gate para a IA

	// Requires restringe o gate aos alvos cujo CONTEÚDO contém este texto.
	//
	// `on` filtra por KIND de nó e `tags` por tag declarada; nenhum dos dois alcança
	// "só as specs que usaram tal marcação". O `no-test-prova-real` é o caso: ele
	// pergunta se a dispensa `@no-test` aponta prova real, e sem filtro enfileirava
	// julgamento de IA para TODA spec do projeto — 583 alvos para 16 que de fato
	// dispensaram. O ⏳ mentia sobre o tamanho do trabalho, e o custo era real: cada
	// alvo enfileirado é uma pergunta a responder.
	//
	// Vazio = sem restrição, para que o silêncio nunca desligue gate existente.
	Requires string `yaml:"requires,omitempty"`

	// NeedsTool é o BINÁRIO externo sem o qual este gate não tem como medir nada.
	//
	// Sem ele, ferramenta ausente virava REPROVAÇÃO: o `sh` sai 127, o gate não produz
	// laudo e o operador lê "violação" onde há problema de ambiente. Isso é pior que
	// ruído — inverte o sentido do veredito, porque o vermelho passa a significar
	// "não instalei" em vez de "achei". Um gate canônico com essa falha reprovaria
	// todo projeto recém-criado, que é a razão de nenhum canônico usar `run:` hoje.
	//
	// Declarado, o motor confere o PATH antes de rodar e devolve Skip ("não se aplica")
	// em vez de Fail. A ausência não some: o `doctor` a levanta como Warn, com o
	// InstallHint. Vazio = sem exigência, para que o silêncio nunca desligue gate algum.
	NeedsTool string `yaml:"needs_tool,omitempty"`

	// InstallHint é o comando que instala o NeedsTool — mostrado pelo `doctor` junto do
	// aviso. Existe para que o aviso termine em AÇÃO: "falta gitleaks" manda o leitor
	// procurar; "falta gitleaks (brew install gitleaks)" resolve.
	InstallHint string `yaml:"install_hint,omitempty"`

	// Scope diz sobre QUANTO o gate roda de uma vez. O default histórico (e o único que
	// existia) é por NÓ: o comando roda uma vez para cada alvo que casa `on`/`tags`.
	// Isso está certo para o que mede um arquivo isolado, e errado para toda ferramenta
	// de projeto: um `tsc` declarado assim rodaria 63 vezes num commit de 63 arquivos —
	// e ainda daria a resposta errada, porque typecheck só faz sentido sobre o programa
	// inteiro. Ver ScopeNode/ScopeBatch/ScopeProject.
	Scope string `yaml:"scope,omitempty"`

	// ScopeFull é o escopo a usar quando a varredura é o PROJETO INTEIRO (`check --all`).
	// Existe porque `batch` responde a uma pergunta que, no full, deixa de fazer sentido:
	// "eis o recorte, confronte-o". No full o recorte é tudo, e continuar passando a lista
	// é só custo — no Windows, custo que estoura a linha de comando; em qualquer sistema,
	// N invocações da ferramenta onde uma bastava.
	//
	// Vale só para quem SABE varrer sozinho. Um script que sai 0 quando não recebe
	// argumento (padrão comum: `[ $# -eq 0 ] && exit 0`) declarado aqui como `project`
	// passaria a certificar sem olhar nada — verde por omissão, o pior estado possível.
	// Por isso é opt-in explícito, e não inferido do fato de a varredura ser completa.
	ScopeFull string `yaml:"scope_full,omitempty"`

	// When declara em que MOMENTOS o gate é cobrado (`pre-commit`, `pre-push`, `ci`,
	// `manual`). Vazio = todas as fases — o silêncio não pode desligar gate que já existe.
	// É o que permite ao hook pedir uma FASE em vez de listar gates: o gate caro declara
	// `when: [ci]` e sai do commit sem sair do pipeline.
	When []string `yaml:"when,omitempty"`

	// Cost é a ordem de grandeza do gate (`fast`, `slow`). Eixo independente de `when`:
	// responde "posso rodar isto num loop apertado?", enquanto `when` responde "isto é
	// cobrado nesta fase?". Vazio = fast (é o que todo gate interno é).
	Cost string `yaml:"cost,omitempty"`

	// SkipOn desliga o gate numa PERSPECTIVA (`change`, `all`) — ver SkipsOn.
	//
	// Lista de EXCLUSÃO, e não de inclusão, para que o silêncio signifique "vale nas
	// duas": um campo omitido nunca pode desligar verificação, senão todo gate escrito
	// antes deste campo existir passaria a não rodar em lugar nenhum.
	SkipOn []string `yaml:"skip_on,omitempty"`

	// Category é a NATUREZA do que o gate mede (`traceability`, `types`, `style`,
	// `test`, `security`…). Não decide execução — organiza o relatório e permite rodar
	// uma família só (`--category types`). Vazio = sem categoria.
	Category string `yaml:"category,omitempty"`

	// Format declara o FORMATO do relatório que este gate ingere, quando a ferramenta
	// do stack não emite o formato canônico.
	//
	// Hoje só o `mutation-score` o usa. O canônico da mutação é o Mutation Testing
	// Elements (Stryker/PIT/Infection/mutmut) e continua sendo o DEFAULT — vazio
	// significa MTE, para que nenhum projeto existente mude de comportamento.
	//
	// Existe porque em Go o canônico não está disponível na prática: o runner dominante
	// (`gremlins`, 391★) emite formato próprio, não há conversor público, e o único
	// runner Go com MTE nativo (`gomutants`) tem 6★ e nasceu em abr/2026. Sem esta
	// chave, o conselho que o `doctor` dá a toda stack — "rode a ferramenta de mutação
	// e ingira o relatório" — era, em Go, um beco: a ferramenta roda e o relatório não
	// entra.
	//
	// A declaração fica AQUI, junto do gate, e não numa flag de linha de comando,
	// porque o formato é propriedade da STACK do projeto — decidido uma vez, não
	// relembrado a cada ingestão. Ver FormatMTE/FormatGremlins.
	Format string `yaml:"format,omitempty"`
}

// Formatos de relatório de mutação aceitos por `Gate.Format`.
const (
	// FormatMTE é o canônico: Mutation Testing Elements, emitido por Stryker (JS/TS/
	// C#/Scala), PIT (Java), Infection (PHP), mutmut (Python). É o default quando
	// `format:` é omitido.
	FormatMTE = "mutation-testing-elements"
	// FormatGremlins é o formato próprio do github.com/go-gremlins/gremlins, o runner
	// de mutação dominante em Go.
	FormatGremlins = "gremlins"
)

// MutationFormat devolve o formato de relatório de mutação declarado no gate
// `mutation-score`, ou o canônico (MTE) quando nada foi declarado.
//
// O silêncio resolve para o canônico DE PROPÓSITO: é a mesma regra dos outros campos
// do gate — uma chave omitida nunca pode mudar o comportamento de um projeto que já
// funcionava.
func (c *Config) MutationFormat() string {
	for _, g := range c.Gates {
		if g.Check != "mutation-score" && g.Name != "mutation-score" {
			continue
		}
		if f := strings.TrimSpace(strings.ToLower(g.Format)); f != "" {
			return f
		}
	}
	return FormatMTE
}

// MeasuresJudgment é o valor de Gate.Measures que marca um gate de julgamento por IA.
const MeasuresJudgment = "judgment"

// IsJudgment diz se este é um gate de julgamento por IA (não computável).
func (g Gate) IsJudgment() bool { return g.Measures == MeasuresJudgment }

// IsBlocking devolve se o gate barra a promoção. O default é INFORMATIVO: um gate cuja
// severidade ninguém declarou não pode barrar commit — a maturação (QUALITY §7) é
// escolha explícita do projeto, não algo que o silêncio decide.
func (g Gate) IsBlocking() bool { return g.Blocking != nil && *g.Blocking }

// Bool devolve um ponteiro para o literal — açúcar para declarar gates.
func Bool(b bool) *bool { return &b }

// Os escopos de execução de um gate.
const (
	// ScopeNode: uma execução POR ALVO. Default — o comportamento que todos os gates
	// tinham antes de `scope:` existir.
	ScopeNode = "node"
	// ScopeBatch: UMA execução recebendo os alvos que casaram como argumentos
	// posicionais ("$@"). É o escopo do eslint/prettier — a ferramenta aceita lista de
	// arquivos e o custo de subir o processo é pago uma vez só.
	ScopeBatch = "batch"
	// ScopeProject: UMA execução, sem alvos. É o escopo do tsc/knip/madge — a
	// ferramenta olha o projeto inteiro e não sabe responder "sobre este arquivo".
	// Roda se HOUVER ao menos um alvo relevante: senão, um commit só de README
	// dispararia um typecheck do monorepo inteiro.
	ScopeProject = "project"
)

// EffectiveScope devolve o escopo do gate, com o default explícito.
func (g Gate) EffectiveScope() string {
	switch g.Scope {
	case ScopeBatch, ScopeProject:
		return g.Scope
	default:
		return ScopeNode
	}
}

// ScopeParaVarredura devolve o escopo do gate para o tipo de varredura em curso: o
// `scope_full` quando o recorte é o projeto inteiro e o gate declarou um, o `scope` de
// sempre no resto. Só batch/project são aceitos como `scope_full` — `node` no full
// significaria uma execução por arquivo do projeto, que é o oposto da intenção.
func (g Gate) ScopeParaVarredura(completa bool) string {
	if !completa {
		return g.EffectiveScope()
	}
	switch g.ScopeFull {
	case ScopeBatch, ScopeProject:
		return g.ScopeFull
	default:
		return g.EffectiveScope()
	}
}

// As PERSPECTIVAS de execução — o eixo "sobre quanto do projeto se pergunta".
//
// Ortogonal a `Scope` (quantas execuções o gate faz) e a `When` (em que momento ele é
// cobrado). Um gate `scope: node` roda uma vez por alvo nas DUAS perspectivas; o que
// muda é o conjunto de alvos que o comando entrega.
const (
	// PerspectiveChange: `--changed` — só o caminho de impacto do que mudou.
	PerspectiveChange = "change"
	// PerspectiveAll: `--all` — a foto completa do projeto.
	PerspectiveAll = "all"
)

// SkipsOn diz se o gate se declara inútil NESTA perspectiva.
//
// O default é participar das duas: `skip_on` é lista de EXCLUSÃO, e o silêncio não pode
// desligar verificação — mesma regra do `When` vazio. Serve para o gate cuja resposta
// só tem valor num dos lados: um detector de órfãos perguntado sobre um recorte de um
// arquivo responde "órfão" para tudo que o recorte não alcança (`skip_on: [change]`);
// um gate de custo proibitivo na varredura completa sai do `--all` sem sair do commit.
func (g Gate) SkipsOn(perspective string) bool {
	for _, p := range g.SkipOn {
		if p == perspective {
			return true
		}
	}
	return false
}

// As fases em que um gate pode ser cobrado.
const (
	PhasePreCommit = "pre-commit"
	PhasePrePush   = "pre-push"
	PhaseCI        = "ci"
	PhaseManual    = "manual"
)

// RunsIn diz se o gate é cobrado na fase dada. Gate sem `when:` roda em TODAS as
// fases — o silêncio não pode desligar um gate que já existia.
//
// `manual` é a exceção deliberada: declarar `when: [manual]` tira o gate das fases
// automáticas (é como se declara "isto é caro demais para o loop, rode sob demanda").
func (g Gate) RunsIn(phase string) bool {
	if phase == "" || len(g.When) == 0 {
		return true
	}
	for _, w := range g.When {
		if w == phase {
			return true
		}
	}
	return false
}

// Os custos declaráveis.
const (
	CostFast = "fast"
	CostSlow = "slow"
)

// IsSlow diz se o gate foi declarado caro.
func (g Gate) IsSlow() bool { return g.Cost == CostSlow }

// IsExternal diz se o gate delega a uma ferramenta de terceiro (`run:`).
func (g Gate) IsExternal() bool { return g.Run != "" }

// Layer — uma camada declarada: como reconhecer seus arquivos, o kind do nó, e as
// tags (rótulos de agrupamento transversais) pelas quais um guide pode regê-la.
type Layer struct {
	Pattern string   `yaml:"pattern"`           // glob que reconhece arquivos da camada
	Kind    string   `yaml:"kind"`              // spec|feature|test|code|doc|guide|plan
	Tags    []string `yaml:"tags,omitempty"`    // rótulos p/ governs (ex.: frontend, mobile)
	Regime  string   `yaml:"regime,omitempty"`  // comportamental|declarativo|misto
	Exclude []string `yaml:"exclude,omitempty"` // globs a excluir (derivados que casam o glob amplo)
	// Priority: desempate DECLARADO quando dois patterns casam o mesmo arquivo. Maior
	// vence; ausente = 0. Sem isso, o Anchors desempata por comprimento do pattern —
	// uma heurística que mede verbosidade, não precisão, e que já classificou errado
	// (um pattern com muitas alternativas vencia outro que apontava para um subconjunto
	// seu). Declare quando a heurística errar; o `check` avisa onde ela decidiu sozinha.
	Priority int `yaml:"priority,omitempty"`
	// CodePrefix: prefixo de módulo (Camada 1 do código de identidade). Se a unidade
	// cai nesta camada, o código começa por este prefixo + chars distintivos. É a
	// ligação do gerador de identidade com a ESTRUTURA — o módulo do arquivo (auth,
	// ir, family…) define seu prefixo, em vez de uma tabela hardcoded.
	CodePrefix string `yaml:"code_prefix,omitempty"`
	// TrincaOpcional: peças da trinca que ESTA camada dispensa, por aresta
	// (`specifies` | `covered-by` | `tested-by`). É o opt-out HONESTO do gate
	// `trinca-completa`: fica declarado na Estrutura, à vista, em vez de escondido
	// num Skip do gate. Ex.: uma camada provada só por teste de integração central
	// declara `trinca_opcional: [tested-by]`.
	TrincaOpcional []string `yaml:"trinca_opcional,omitempty"`
	// Work: passos EXTRA que esta camada exige, por artefato (spec|code|feature|test).
	// O `anchors work` já compõe um procedimento universal a partir da Estrutura; isto
	// acrescenta o que só o projeto sabe ("rode o seed antes", "o teste desta camada é
	// de integração e exige AWS"). Opcional — sem isso, vale só o universal.
	Work map[string][]string `yaml:"work,omitempty"`
	// SectionTitles: o léxico de seções DESTA camada, por chave do catálogo. Vence o
	// `section_titles` do projeto. É o nível certo para a maioria dos dialetos — eles
	// costumam ser da camada (um modelo de dado fala "Modelo de Dado"; uma tela, não).
	SectionTitles SectionTitles `yaml:"section_titles,omitempty"`
}

// Derived — a superfície da trinca (STRUCTURE.md §2.2): onde as peças derivadas de um
// "anchor" (spec/feature/test) DEVEM morar. Por default é co-location (partilham o stem
// {{dir}}/{{name}}); `overrides` declara padrões de localização por camada-âncora quando
// a peça NÃO é co-localizada (ex.: testes centralizados num backend). A ligação material
// segue por CÓDIGO (TRACEABILITY) — isto é só a superfície de descoberta/validação.
// Variáveis nos templates: {{dir}} (dir da âncora), {{name}} (basename sem ext), {{ext}}
// (extensão da âncora), {{module}} (o dir-pai — útil quando o basename é genérico, ex.
// `handler` em functions/<module>/handler.ts → o módulo é <module>).
type Derived struct {
	Anchor    string            `yaml:"anchor"`              // a camada-âncora (ex.: "code")
	Files     map[string]string `yaml:"files"`               // default (co-location): camada-derivada → template
	Overrides []DerivedOverride `yaml:"overrides,omitempty"` // padrões por camada-âncora (não co-localizado)
	// Regimes — o de-para do vocabulário de REGIME de teste do PROJETO para o regime
	// CANÔNICO do framework (unit|integration|e2e|vr). Chave = a tag do projeto (sem @,
	// ex.: "nivel-unit"); valor = o regime canônico. Um cenário de feature declara sua
	// tag de regime; o gate de correspondência a traduz p/ o canônico e confronta contra
	// a SUPERFÍCIE do regime. Tag sem mapeamento = não confrontada (opt-out). STRUCTURE §2.3.
	Regimes map[string]string `yaml:"regimes,omitempty"`
	// Surfaces — a superfície de cada regime canônico: qual chave de `Files`/`Overrides`
	// localiza a peça de teste daquele regime (unit→"test", e2e→"e2e", vr→"vr"…). Regime
	// sem superfície mapeada = não confrontado por arquivo (verificado fora do grafo).
	Surfaces map[string]string `yaml:"surfaces,omitempty"`
	// TestHandle — como ESTE projeto marca um elemento para ser alcançado de fora
	// (pelo teste, pelo flow de ponta a ponta). O nome do atributo é do ecossistema,
	// não do framework: React Native usa `testID`, a web `data-testid`, Android
	// `contentDescription`. Um backend não marca nada e deixa vazio.
	//
	// É o que torna os gates de inventário de handle OPT-IN por natureza: sem esta
	// declaração eles não têm o que procurar e pulam (Skip, não Fail) — mesmo que o
	// projeto os liste em `gates`. Declarar o atributo é o ato de dizer "aqui a
	// ancoragem de teste é contrato, cobre-a de mim".
	TestHandle string `yaml:"test_handle,omitempty"`

	// MockContract — como ESTE projeto amarra um dublê de teste ao módulo que ele
	// substitui.
	//
	// O problema que isto ataca: um teste que mocka o vizinho continua VERDE depois que
	// o vizinho muda. O mock virou uma cópia congelada de um contrato que já não existe,
	// e o verde passa a certificar a versão antiga — pior que a ausência de teste,
	// porque ninguém procura defeito onde há prova.
	//
	// A defesa é fazer o dublê DERIVAR do original em vez de repeti-lo. Como isso se
	// escreve é do ecossistema, não do framework: em TypeScript, anotar a fábrica com
	// `Partial<typeof Real>` faz o compilador conferir nome, assinatura e retorno; em
	// Python é `autospec=True`; em Go, satisfazer a interface. O projeto declara a forma
	// que vale nele, com `{{module}}` onde entra o módulo dublado:
	//
	//	mock_contract: "Partial<typeof {{module}}>"
	//
	// Vazio = o gate `mock-tipado` não tem forma a cobrar e PULA (Skip, não Fail),
	// mesmo que o projeto o declare em `gates`. Declarar é o ato de dizer "aqui o dublê
	// é contrato, cobre-o de mim".
	//
	// COBRE APENAS drift de FORMA no nível do MÓDULO (que chaves existem). Não alcança a
	// forma do VALOR devolvido — medido: 202 testes reprovaram por dublar um hook do
	// React Query com dois campos onde o real devolve ~15, e `Partial<typeof X>` não
	// distingue "dublê parcial deliberado" de "dublê defasado". Para esse resto existe
	// o CARIMBO (`MockStamp`), que é agnóstico de linguagem.
	MockContract string `yaml:"mock_contract,omitempty"`

	// MockStamp — o CARIMBO DE CONTRATO: a marca que um dublê carrega do trecho de
	// código que ele substitui, para o gate recalcular e detectar quando ele mudou.
	//
	// Existe porque `MockContract` depende de a linguagem ter tipos estruturais, e a
	// maior parte não tem (Python com dict, JS puro, Ruby). O carimbo é hash de TEXTO:
	// funciona em qualquer linguagem e pega qualquer mudança — assinatura, corpo, tipo,
	// constante — porque não interpreta o código, só o lê.
	//
	// O formato é `<âncora>|<qtd>|<hash>`, e cada parte resolve um problema medido:
	//
	//	âncora  a LINHA INTEIRA que abre o trecho (`export function useX(`), não o
	//	        número dela. Número é POSIÇÃO: um comentário no topo do arquivo
	//	        deslocaria tudo abaixo e quebraria o carimbo de quem não mudou nada
	//	        (medido: +1 linha na linha 5 invalidava a função da linha 70). Âncora
	//	        por conteúdo é imune a deslocamento e, quando SOME, acusa explicitamente
	//	        — renome vira erro, não silêncio.
	//	qtd     quantas linhas a partir da âncora entram no hash. FIXA e escrita no
	//	        carimbo: o alcance fica à vista de quem lê, e o gate não precisa
	//	        adivinhar onde o bloco termina — o que exigiria um parser por linguagem
	//	        e mataria o agnosticismo.
	//	hash    do trecho. É o que muda quando o contrato muda.
	//
	// O gate RECALCULA — não valida formato. Um carimbo bem-formatado que ninguém
	// confronta é teatro: quem edita o teste o regeneraria para casar com o próprio
	// mock, e ele passaria a certificar a si mesmo em vez do módulo. Recalcular é
	// barato (ler N linhas a partir de uma âncora não pede parser) e é o que torna o
	// mecanismo à prova de quem o escreve.
	//
	// Vazio = o gate `mock-carimbado` pula. Quem o habilita é `MockDetect`: sem saber
	// RECONHECER um dublê, não há o que carimbar — uma flag separada seria uma segunda
	// chave para a mesma decisão, e duas chaves divergem.

	// RuleMarking declara que ESTE projeto EXIGE a marcação regra↔código: cada regra da
	// spec aparece como comentário no trecho que a realiza (`// ABCDX-B01: …`).
	//
	// Existe porque a prática não pode ser cobrada por default. O `regra-implementada`
	// nasceu num projeto com 3.114 regras não marcadas: exigi-las de saída reprovaria
	// 98% das unidades no primeiro dia, e um gate assim é desligado — levando junto os
	// que funcionam. Por isso a unidade que nunca marcou NADA cai numa pendência de
	// migração, que nomeia a dívida sem barrar.
	//
	// O problema é que essa pendência NUNCA VENCE. Medido: uma spec afirmava "mantém o
	// cartão selecionado" enquanto o código era um CRUD sem seleção nenhuma — o gate viu
	// (a regra `SBNKX-B01` não estava no código), caiu no ramo de migração e devolveu
	// pendência. Pendência não barra e ninguém a revisita: o defeito atravessou 44 gates.
	//
	// Declarar `rule_marking: required` é o ato de dizer "aqui a migração acabou": a
	// pendência vira reprovação, e a unidade que não marcou passa a ser cobrada como
	// qualquer outra. Projeto NOVO deve nascer com isto ligado — não há dívida a
	// acomodar, e o primeiro código é escrito depois da regra existir.
	RuleMarking string `yaml:"rule_marking,omitempty"`

	// MockDetect — como ESTE projeto ESCREVE um dublê de teste. Regex com um grupo de
	// captura: o módulo dublado.
	//
	// É a chave que faltava para o carimbo ser de fato agnóstico. O hash de texto não
	// depende de linguagem, mas RECONHECER o dublê depende: `jest.mock('x', () => …)`
	// é dialeto de jest/vitest, e nada nele casa `unittest.mock.patch('x')` do Python,
	// `allow(X).to receive(...)` do Ruby ou `Mockito.mock(X.class)` do Java. Em Go nem
	// há chamada a detectar — o dublê é uma interface satisfeita.
	//
	// Sem esta declaração o gate PULA. É deliberado, e a alternativa é pior: embutir o
	// padrão jest/vitest como default faria o gate rodar num projeto Python, encontrar
	// zero dublês e reportar VERDE sobre o que não conferiu — a pior falha possível num
	// medidor, e ainda com viés de ecossistema escondido no silêncio.
	//
	// Declarar o regex É habilitar o gate: presença significa "aqui o dublê é contrato,
	// cobre-o de mim". Não há flag separada porque duas chaves para uma decisão
	// divergem — uma ligada sem a outra vira comportamento indefinido.
	//
	// Exemplos:
	//
	//	mock_detect: "(?:jest|vi)\\.mock\\(['\"]([^'\"]+)"        # TS/JS
	//	mock_detect: "(?:mock\\.)?patch\\(['\"]([^'\"]+)"          # Python
	//	mock_detect: "Mockito\\.mock\\(([A-Za-z0-9_.]+)\\.class"   # Java
	//
	// Um regex MAL FORMADO para o dialeto do projeto falha em silêncio (casa nada, o
	// gate passa). Quem cobre esse ponto cego é o gate de julgamento
	// `mock-detect-cobre-o-dialeto`, que lê os testes e diz se o padrão declarado
	// alcança o que o projeto realmente escreve.
	MockDetect string `yaml:"mock_detect,omitempty"`
}

// DerivedOverride — para âncoras cuja tag casa `When`, os templates em `Files`
// SOBRESCREVEM os do default (só as camadas presentes; as demais herdam o default).
type DerivedOverride struct {
	When  string            `yaml:"when"`  // a tag/camada da âncora a que este override se aplica (ex.: handler)
	Files map[string]string `yaml:"files"` // camada-derivada → template de região (aceita {{module}})
}

// GovernRule — uma aresta vertical declarada: a régua `from` rege os nós de TODAS
// as layers que carregam a tag `governs`. É sempre por TAG (nunca por nome de
// layer) — um só mecanismo, sem ambiguidade. O escopo vem dos patterns das layers
// marcadas com a tag, então não há glob duplicado aqui (DRY).
//
// `from` é o caminho do arquivo-régua (ex.: guides/SCREEN_GUIDE.md).
type GovernRule struct {
	From    string `yaml:"from"`    // caminho do arquivo-régua (o guide)
	Governs string `yaml:"governs"` // a tag que ele rege (casa layers com essa tag)
}

// canonicalGate resolve a declaração canônica de um gate pelo nome. É INJETADO por
// `initx` (via SetCanonicalGateResolver) porque o catálogo vive lá junto da semente do
// `init` — e `config` não pode importar `initx` sem criar ciclo.
//
// Nil quando ninguém registrou: aí o merge é no-op e o comportamento é o de antes.
var canonicalGate func(name string) (Gate, bool)

// SetCanonicalGateResolver registra a fonte das declarações canônicas de gate.
func SetCanonicalGateResolver(f func(name string) (Gate, bool)) { canonicalGate = f }

// Load lê o anchors.yaml da raiz do projeto.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	// KnownFields(true): chave desconhecida no anchors.yaml é ERRO, não silêncio.
	//
	// Com o Unmarshal padrão, uma chave que o Config não conhece é simplesmente ignorada —
	// e o efeito é o pior possível para quem está adotando: a declaração não faz nada e
	// nada avisa. Medido ao pôr o Anchors no próprio Anchors: escrevi o bloco `governs`
	// com `guide:`/`tags:` (plural) em vez de `from:`/`governs:`, o `map build` respondeu
	// "222 nós, 0 arestas" e segui adiante achando que o projeto não tinha relações. A
	// declaração inteira tinha sido descartada.
	//
	// O modo de falha é traiçoeiro porque o arquivo PARECE certo: carrega sem erro, o
	// comando roda, o resultado é só mais pobre. Quem não conhece o formato de memória não
	// tem como distinguir "declarei errado" de "meu projeto realmente não tem isso".
	//
	// Erro com o nome da chave e a linha é barato de consertar; um mapa vazio sem
	// explicação custa uma sessão de investigação — foi o que custou aqui.
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	for i := range c.Gates {
		c.Gates[i] = mergeCanonical(c.Gates[i])
	}
	if err := c.validarIDsDeGate(); err != nil {
		return nil, err
	}
	if err := c.validarWorkflow(); err != nil {
		return nil, err
	}
	if err := c.validarEnumsDeGate(); err != nil {
		return nil, err
	}
	for _, l := range c.CodeLengths {
		// 2 é o mínimo que ainda distingue unidades; acima de 8 o código deixa de ser
		// legível de relance, que é a razão de ele ser curto. Fora disso é typo.
		if l < 2 || l > 8 {
			return nil, fmt.Errorf("code_lengths: %d fora da faixa (2..8) — o código é "+
				"identificador curto, lido de relance; comprimento assim costuma ser typo", l)
		}
	}
	SetCodeLengths(c.CodeLengths)
	// O GERADOR de código também segue o projeto. Sem isto, `anchors code Login` respondia
	// `LGNOI` (5 slots, o canônico) num projeto que declara 4 e tem 659 códigos de 4 — o CLI
	// sugeria identidade que o próprio engine do projeto não reconhece.
	//
	// Injetado em vez de importado (mesmo padrão do SetCanonicalGateResolver): `config` fica
	// sem depender de `code`, e quem liga é o main.
	if setSlots != nil {
		setSlots(CodeLengths)
	}
	return &c, nil
}

// setSlots é o gancho do gerador de código. Registrado por SetSlotsHook no main.
var setSlots func([]int)

// SetSlotsHook registra a função que ajusta o comprimento GERADO ao do projeto.
func SetSlotsHook(f func([]int)) { setSlots = f }

// validarEnumsDeGate cobra os VALORES dos campos enumerados de gate.
//
// `KnownFields` pega chave inexistente; não pega valor inexistente em chave válida. E o
// segundo caso é igualmente silencioso: declarei `scope: repo` (o certo é `project`), o
// YAML carregou sem reclamar, e três gates entraram como "sem nada para medir". Levei três
// tentativas para achar — procurando no `on:`, no mapa, no catálogo — porque nada apontava
// para o valor.
//
// Enum errado é sempre gate DESLIGADO em silêncio: um `scope` que não casa nenhum dos três
// cai no default por nó; um `when` com typo tira o gate da fase que o autor queria cobrar;
// um `cost` inválido faz o gate ser tratado como rápido e entrar num loop apertado que
// deveria evitar. Nos três, o efeito é proteção que o arquivo afirma e não existe.
// validarIDsDeGate cobra ID único em cada gate — na CARGA, não no uso.
//
// O ID é o que outros lugares citam: o relatório o imprime, e é por ele que se dispensa
// uma regra. Dois gates com o mesmo ID fazem toda citação apontar para dois lugares, e
// uma dispensa desligar o que ninguém pediu.
//
// A ausência é tolerada por compatibilidade — um projeto que já funcionava não pode
// quebrar por um campo novo —, e aí o nome do gate serve de ID. Mas a COLISÃO nunca é
// tolerada: ela é sempre erro, e falhar na carga é o que impede o achado errado.
func (c *Config) validarIDsDeGate() error {
	vistos := map[string]string{} // id → nome do primeiro gate que o usou
	for _, g := range c.Gates {
		id := g.ID
		if id == "" {
			id = g.Name // sem ID declarado, o nome identifica
		}
		if anterior, ok := vistos[id]; ok {
			return fmt.Errorf("gate %q: `id: %q` já é usado pelo gate %q — o ID é o que o "+
				"relatório imprime e o que uma dispensa cita; repetido, ele aponta para dois "+
				"lugares e desliga o que ninguém pediu", g.Name, id, anterior)
		}
		vistos[id] = g.Name
	}
	return nil
}

func (c *Config) validarEnumsDeGate() error {
	escopos := map[string]bool{ScopeNode: true, ScopeBatch: true, ScopeProject: true}
	custos := map[string]bool{CostFast: true, CostSlow: true}
	fases := map[string]bool{
		PhasePreCommit: true, PhasePrePush: true, PhaseCI: true, PhaseManual: true,
	}
	for _, g := range c.Gates {
		if g.Scope != "" && !escopos[g.Scope] {
			return fmt.Errorf("gate %q: `scope: %q` desconhecido — use %q (uma execução por alvo), "+
				"%q (uma execução com os alvos como argumentos) ou %q (uma execução, sem alvos)",
				g.Name, g.Scope, ScopeNode, ScopeBatch, ScopeProject)
		}
		if g.Cost != "" && !custos[g.Cost] {
			return fmt.Errorf("gate %q: `cost: %q` desconhecido — use %q ou %q",
				g.Name, g.Cost, CostFast, CostSlow)
		}
		for _, f := range g.When {
			if !fases[f] {
				return fmt.Errorf("gate %q: `when: [%q]` não é uma fase — use %q, %q, %q ou %q",
					g.Name, f, PhasePreCommit, PhasePrePush, PhaseCI, PhaseManual)
			}
		}
	}
	return nil
}

// validarWorkflow cobra a coerência do bloco `workflow` na CARGA, não no uso.
//
// Um modo mal declarado tem de falhar no primeiro comando que lê a config, com a linha do
// que falta. Deixar para o momento do `anchors next` significaria descobrir a configuração
// quebrada no meio de um trabalho, e — pior — depois de o agente já ter feito alterações
// achando que sabia de onde vinha a task.
func (c *Config) validarWorkflow() error {
	if c.Workflow == nil {
		return nil
	}
	w := c.Workflow
	switch w.Mode {
	case "", ModeLocal:
		// O modo local não usa Repo nem Labels. Declará-los aqui não é inofensivo: quem lê
		// o arquivo conclui que a integração está ativa, e ela não está.
		if w.Repo != "" || len(w.Labels) > 0 {
			return fmt.Errorf("workflow: `repo`/`labels` só valem em `mode: github` — " +
				"no modo local eles não são lidos, e deixá-los declarados faz o arquivo " +
				"afirmar uma integração que não existe")
		}
		return nil
	case ModeGitHub:
		if w.Repo == "" {
			return fmt.Errorf("workflow: `mode: github` exige `repo: owner/nome` — " +
				"o Anchors não infere do remote do git de propósito: num fork, inferir " +
				"faria a escrita cair no repositório errado")
		}
		if !strings.Contains(w.Repo, "/") {
			return fmt.Errorf("workflow: `repo: %q` fora do formato `owner/nome`", w.Repo)
		}
		if len(w.Labels) == 0 {
			return fmt.Errorf("workflow: `mode: github` exige ao menos uma label em " +
				"`labels` — sem ela, `anchors next` puxaria qualquer issue do repositório, " +
				"inclusive as de produto, que não têm a forma que o ciclo espera")
		}
		return nil
	default:
		return fmt.Errorf("workflow: `mode: %q` desconhecido — use `local` ou `github` "+
			"(não há fallback entre eles: o modo é declarado, não adivinhado)", w.Mode)
	}
}

// mergeCanonical completa os campos OMITIDOS de um gate canônico a partir da declaração
// do framework. O yaml do projeto sempre vence onde declara — o merge só preenche buraco.
//
// Existe porque um gate canônico declarado no projeto precisava REPETIR a redação inteira
// (`check`, `measures`, `ask`, `guide`…). A redação passava a viver em dois lugares e
// divergia no primeiro ajuste: o framework corrigia a pergunta de um gate de julgamento e
// os projetos seguiam perguntando a versão velha, sem ninguém perceber. Com o merge,
// `- name: no-test-prova-real` basta.
//
// `Blocking` entra no merge como os demais porque é `*bool`: `nil` é "a chave não veio",
// distinto do `blocking: false` explícito. Fosse `bool`, herdá-lo promoveria a bloqueante
// um gate que o projeto manteve informativo DE PROPÓSITO — barrando commits liberados
// conscientemente. É o campo que motivou o ponteiro.
func mergeCanonical(g Gate) Gate {
	if canonicalGate == nil || g.Name == "" {
		return g
	}
	base, ok := canonicalGate(g.Name)
	if !ok {
		return g
	}
	if len(g.On) == 0 {
		g.On = base.On
	}
	if len(g.Tags) == 0 {
		g.Tags = base.Tags
	}
	if g.Run == "" {
		g.Run = base.Run
	}
	if g.Check == "" {
		g.Check = base.Check
	}
	if g.Measures == "" {
		g.Measures = base.Measures
	}
	if g.Guide == "" {
		g.Guide = base.Guide
	}
	if g.Requires == "" {
		g.Requires = base.Requires
	}
	if g.Ask == "" {
		g.Ask = base.Ask
	}
	// A ferramenta exigida vem do canônico junto com o `run:` que a usa — os dois
	// descrevem o MESMO comando, e herdar um sem o outro produziria o pior estado
	// possível: o gate roda o comando canônico e reprova com "command not found",
	// que é exatamente o que `needs_tool` existe para impedir.
	if g.NeedsTool == "" {
		g.NeedsTool = base.NeedsTool
	}
	if g.InstallHint == "" {
		g.InstallHint = base.InstallHint
	}
	if g.Scope == "" {
		g.Scope = base.Scope
	}
	// ScopeFull acompanha o Scope: são um par (o recorte no incremental, o projeto
	// inteiro no `--all`). Herdar só o primeiro faria o gate canônico varrer em lotes
	// no full — N invocações onde uma bastava, e no Windows a linha de comando estoura.
	if g.ScopeFull == "" {
		g.ScopeFull = base.ScopeFull
	}
	if len(g.When) == 0 {
		g.When = base.When
	}
	if g.Cost == "" {
		g.Cost = base.Cost
	}
	if g.Category == "" {
		g.Category = base.Category
	}
	if len(g.SkipOn) == 0 {
		g.SkipOn = base.SkipOn
	}
	if g.Blocking == nil {
		g.Blocking = base.Blocking
	}
	return g
}

// Save escreve o anchors.yaml.
func Save(c *Config, path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	header := []byte("# anchors.yaml — configuração do projeto Anchors\n" +
		"# A seção `layers` é a Estrutura de Projeto (o grafo virtual).\n")
	return os.WriteFile(path, append(header, data...), 0o644)
}

// RouteRegistry devolve os globs onde o projeto registra rotas de navegação.
func (c *Config) RouteRegistry() []string {
	if c == nil {
		return nil
	}
	return c.RouteRegistryGlobs
}

// ─────────────────────────────────────────────────────────────────────────────
// SUÍTES — como o PROJETO produz sinal de teste
// ─────────────────────────────────────────────────────────────────────────────

// Suite é um comando do projeto que produz um relatório de teste, mais o caminho do
// relatório que ele deixa. O Anchors não sabe rodar teste — não conhece jest, pytest
// nem go test, e não deve conhecer: a stack é do projeto. Mas ele sabe AMARRAR o
// relatório ao mapa, e é aí que o par se separava.
//
// Antes disto, produzir sinal era dois passos manuais: rode a ferramenta, depois
// lembre o caminho do relatório e chame `anchors ingest`. O modo de falha é
// silencioso e comum — quem roda e esquece de ingerir vê o gate continuar acusando
// "sem sinal ingerido", como se a rodada não tivesse acontecido. Medido num projeto
// real: 1008 das 1011 pendências do `check --all` eram sinal ausente, não qualidade
// ruim.
//
// A assimetria era só do lado do sinal: um `gate` já declara seu `run:` e o Anchors o
// executa sem saber o que é eslint. Suite dá ao sinal a mesma cidadania.
type Suite struct {
	// Workspace é ONDE a suíte roda (backend, mobile, landing…). Em monorepo é o
	// segundo eixo, e ele é independente da camada: `unit` existe no backend E no
	// mobile, então "rode a unit" e "rode o backend" são perguntas diferentes —
	// colapsar os dois num campo só obrigaria a inventar camadas falsas ("mobile" como
	// se fosse um nível de teste) e tornaria impossível pedir a unit de todos.
	//
	// Vazio é legítimo: projeto de pacote único não tem esse eixo, e aí `--workspace`
	// não filtra nada. O vocabulário é do projeto, como o da camada.
	Workspace string `yaml:"workspace,omitempty"`
	// Layer é o vocabulário do PROJETO (unit, integration, e2e, contract, smoke…). É
	// por ele que `anchors test unit e2e` filtra, e é ele que vai para o `--layer` da
	// ingestão — o Anchors não valida a lista, porque não é dono dela.
	Layer string `yaml:"layer"`
	// Run é o comando do modo COMPLETO (o projeto inteiro), executado via `sh -c` na
	// raiz. `{{target}}`, se presente, recebe o alvo passado em `--target`.
	Run string `yaml:"run"`
	// RunChanged é o comando do modo INCREMENTAL, usado por `--changed`. `{{files}}`
	// recebe os arquivos do caminho de impacto que o Anchors calculou.
	//
	// É um campo à parte, e não o mesmo `run:` com um placeholder opcional, porque em
	// quase todo runner as duas coisas são invocações DIFERENTES, não a mesma com uma
	// lista no fim: jest quer `--findRelatedTests`, go test quer a lista de pacotes,
	// pytest quer os caminhos. Espremer as duas num campo só obrigaria o projeto a
	// escrever um comando que funciona nos dois modos, e não existe um.
	//
	// Ausente: `--changed` recusa a suíte em vez de rodar a completa. Rodar tudo quando
	// se pediu o incremental é caro e, pior, mente sobre o que rodou.
	RunChanged string `yaml:"run_changed,omitempty"`
	// JUnit e Lcov são os relatórios que a suíte DEIXA, ingeridos ao final se o
	// comando passar. Vazios: o Anchors roda e não ingere nada — o que ainda é útil
	// (um comando só) mas mantém o gate reclamando, e é melhor dizer isso do que
	// fingir que ingeriu.
	JUnit string `yaml:"junit,omitempty"`
	Lcov  string `yaml:"lcov,omitempty"`
	// Report é o relatório de MUTAÇÃO (formato Mutation Testing Elements) e só vale
	// nas suítes de `mutation:`.
	Report string `yaml:"report,omitempty"`
	// Scope diz se os mutantes rodaram contra o teste da própria unidade (`isolated`)
	// ou também contra o dos dependentes (`full`). Importa: com os dois ingeridos, o
	// gate julga pelo ISOLADO, que é bem mais duro.
	Scope string `yaml:"scope,omitempty"`
}

// SelecionaSuites filtra pelos DOIS eixos e devolve o que sobra na ordem em que foi
// DECLARADO — não na ordem em que o usuário digitou. Quem escreveu o arquivo decidiu
// que unit vem antes de e2e, e isso costuma ser dependência real; respeitar a digitação
// faria o mesmo pedido se comportar diferente conforme quem digita.
//
// Os eixos se COMBINAM (interseção): `camadas=[unit] workspaces=[backend]` é "a unit do
// backend". Cada eixo vazio significa "todos" — é o que faz `anchors test unit` valer
// para todos os workspaces, e `--workspace backend` para todas as camadas dele.
//
// As ausências saem já rotuladas ("camada \"smoke\"", "workspace \"web\"") e juntas: quem
// errou o nome precisa ver de uma vez o que existe, não descobrir por tentativa.
func SelecionaSuites(suites []Suite, camadas, workspaces, escopos []string) (sel []Suite, ausentes []string) {
	normaliza := func(vs []string) map[string]bool {
		m := map[string]bool{}
		for _, v := range vs {
			if k := strings.ToLower(strings.TrimSpace(v)); k != "" {
				m[k] = true
			}
		}
		return m
	}
	pedidaC, pedidaW, pedidaE := normaliza(camadas), normaliza(workspaces), normaliza(escopos)

	for _, s := range suites {
		c, w, e := strings.ToLower(s.Layer), strings.ToLower(s.Workspace), strings.ToLower(s.Scope)
		if len(pedidaC) > 0 && !pedidaC[c] {
			continue
		}
		if len(pedidaW) > 0 && !pedidaW[w] {
			continue
		}
		if len(pedidaE) > 0 && !pedidaE[e] {
			continue
		}
		sel = append(sel, s)
	}
	// A ausência é aferida contra o que EXISTE no arquivo, não contra o que sobreviveu
	// à interseção: pedir `unit --workspace mobile` num projeto que tem os dois, mas
	// não a unit do mobile, não é nome errado — é combinação vazia, e dizer "camada
	// unit não existe" mandaria corrigir o que está certo.
	existeC, existeW, existeE := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, s := range suites {
		existeC[strings.ToLower(s.Layer)] = true
		existeW[strings.ToLower(s.Workspace)] = true
		existeE[strings.ToLower(s.Scope)] = true
	}
	for _, c := range camadas {
		if k := strings.ToLower(strings.TrimSpace(c)); k != "" && !existeC[k] {
			ausentes = append(ausentes, "camada "+strconvQuote(c))
		}
	}
	for _, w := range workspaces {
		if k := strings.ToLower(strings.TrimSpace(w)); k != "" && !existeW[k] {
			ausentes = append(ausentes, "workspace "+strconvQuote(w))
		}
	}
	for _, e := range escopos {
		if k := strings.ToLower(strings.TrimSpace(e)); k != "" && !existeE[k] {
			ausentes = append(ausentes, "escopo "+strconvQuote(e))
		}
	}
	return sel, ausentes
}

func strconvQuote(s string) string { return `"` + s + `"` }

// CamadasDeclaradas e WorkspacesDeclarados listam o vocabulário do projeto, sem
// repetir, na ordem de declaração — é o que se mostra a quem pediu um nome que não
// existe, e a ordem do arquivo é como a pessoa vai reencontrá-los lá.
func CamadasDeclaradas(suites []Suite) []string {
	return distintos(suites, func(s Suite) string { return s.Layer })
}

func WorkspacesDeclarados(suites []Suite) []string {
	return distintos(suites, func(s Suite) string { return s.Workspace })
}

// EscoposDeclarados é o terceiro eixo, e ele só existe na mutação: a MESMA unidade
// medida contra o teste dela (`isolated`) e contra o de quem a importa (`full`). Não é
// camada nem workspace — é a mesma suíte, com abrangência diferente, e o gate julga
// pelo isolado quando os dois existem.
func EscoposDeclarados(suites []Suite) []string {
	return distintos(suites, func(s Suite) string { return s.Scope })
}

func distintos(suites []Suite, campo func(Suite) string) []string {
	visto := map[string]bool{}
	out := []string{}
	for _, s := range suites {
		v := campo(s)
		if v == "" || visto[strings.ToLower(v)] {
			continue
		}
		visto[strings.ToLower(v)] = true
		out = append(out, v)
	}
	return out
}

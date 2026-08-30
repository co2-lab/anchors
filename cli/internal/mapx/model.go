// Package mapx modela e persiste o mapa de dependências — o anchors.graph.yaml.
//
// O mapa é a materialização do grafo (CONCEPT §3) e o artefato-dono da
// Rastreabilidade (TRACEABILITY §4). É a projeção da Estrutura (o grafo virtual,
// STRUCTURE §2.1) sobre os arquivos reais: a Estrutura diz que TIPO depende de que
// tipo; o mapa diz QUAL arquivo depende de qual, com versão e carimbo por aresta.
package mapx

// Kind é o tipo de um nó — os valores literais canônicos de CONCEPT §3.
type Kind string

const (
	KindSpec    Kind = "spec"
	KindFeature Kind = "feature"
	KindTest    Kind = "test"
	KindCode    Kind = "code"
	KindDoc     Kind = "doc"
	KindGuide   Kind = "guide"
	KindPlan    Kind = "plan"
)

// EdgeType é o tipo semântico da aresta (CONCEPT §3). O tipo carrega a força
// (blocking vs. informativa) e a pergunta de confronto.
type EdgeType string

const (
	EdgeGoverns    EdgeType = "governs"    // o pai rege o filho (bloqueante)
	EdgeSpecifies  EdgeType = "specifies"  // spec descreve o código (bloqueante)
	EdgeCoveredBy  EdgeType = "covered-by" // feature cobre a spec (bloqueante)
	EdgeTestedBy   EdgeType = "tested-by"  // teste exercita a feature (bloqueante)
	EdgeReferences EdgeType = "references" // só rastreabilidade (informativa)
	// EdgeDependsOn — a aresta de REÚSO entre camadas (SPEC_TYPES §5, TRACEABILITY
	// §"arestas declared"): a spec consumidora depende de outra camada (tela→hook/store,
	// usecase→repository). É por ela que a Propagação desce pelas camadas de dados.
	EdgeDependsOn EdgeType = "depends-on"
	// EdgeSeeds — o PLANO semeia a spec. É a aresta que liga a decisão de fazer ao
	// artefato que a realiza, e sem ela o plano fica órfão do grafo: `anchors impact`
	// respondia "nenhum filho depende dele" sobre um plano que nomeia 11 specs, e
	// mudar o plano não propagava para nada.
	//
	// Informativa, não bloqueante: um plano cita specs que ainda VÃO nascer, e cobrar a
	// existência delas transformaria todo plano em erro até a última fase terminar.
	EdgeSeeds EdgeType = "seeds"

	// EdgeNeeds — este PLANO precisa que outro TERMINE antes de começar (`needs:` no
	// header). Diferente de `depends-on`, que diz "consumo este arquivo": aqui a
	// pergunta é de ORDEM DE TRABALHO, não de reúso.
	//
	// Existe porque vários agentes puxam de uma fila comum. Sem ela, alguém pega o card
	// de uma feature antes de a fundação existir, e descobre isso ao tentar criar o
	// primeiro arquivo — depois de já ter reivindicado o trabalho.
	EdgeNeeds EdgeType = "needs"
)

// Origin — como a aresta entrou no mapa (TRACEABILITY §4).
type Origin string

const (
	OriginConvention Origin = "convention" // co-location de nomes
	OriginDeclared   Origin = "declared"   // declarada na spec/no mapa
	OriginInferred   Origin = "inferred"   // inferida (imports/símbolos)
)

// Node — um arquivo que participa do grafo.
type Node struct {
	ID        string `yaml:"id"` // caminho relativo à raiz do projeto
	Kind      Kind   `yaml:"kind"`
	Rev       string `yaml:"rev"`                  // revisão do conteúdo (hash) — a "versão do arquivo"
	UpdatedAt string `yaml:"updated_at,omitempty"` // carimbo de alteração (ISO8601)
	Code      string `yaml:"code,omitempty"`       // código de cenário/identidade, se houver
	// Layer é a camada que este arquivo casou na Estrutura (`layers:` do anchors.yaml).
	//
	// Existe porque a config resolve decisões POR CAMADA — `section_titles` é a primeira,
	// e o léxico de seções quase nunca é global (medido: "Modelo de Dado" em 50 specs de
	// uma camada e em 1 das outras 588). Quem lê o grafo depois do scan — o `doctor`, um
	// relatório — não tinha como resolver isso e caía no nível do projeto, contando zero
	// numa spec cuja seção se chama outra coisa naquela camada.
	Layer string `yaml:"layer,omitempty"`
	// CodeDeclarado distingue identidade DECLARADA (o `code:` do header @anchors, ou a
	// tabela/cabeçalho de regra numa spec) de identidade INFERIDA (o primeiro código que
	// apareceu no texto, quando não havia declaração).
	//
	// Sem a distinção, um fixture de teste que usa `AAAAX-B01` como string de exemplo entra
	// no mapa como DONO de `AAAA` — indistinguível de uma unidade real. Medido no próprio
	// Anchors: 29 códigos "em uso", todos vindos de fixture e de exemplo em documentação,
	// zero de unidade (o projeto tem 0 specs). Um comando que confira identidade sobre esse
	// conjunto produz 29 falsos positivos e nenhum achado.
	//
	// Quem PROMETE identidade declara; o resto é citação, e citação não se confere.
	CodeDeclarado bool     `yaml:"code_declarado,omitempty"`
	Tags          []string `yaml:"tags,omitempty"` // tags da camada do nó (p/ gates escopados por tag)
	// Regime da camada do nó (comportamental|declarativo|misto), copiado da Estrutura
	// (config.Layer.Regime). `declarativo` = camada RECONHECIDA (não origina regra): o
	// gate header-conforme aceita `layer:` como identidade mínima, sem exigir code/ref.
	// Vazio quando a camada não declara regime (fallback aos nomes canônicos reconhecidos).
	Regime string `yaml:"regime,omitempty"`
	// NoPropagation: o filho declarou `@noPropagation` — ele NÃO depende do pai, então
	// a onda de propagação (descida) NÃO passa por ele. Opt-out honesto, tag no filho.
	NoPropagation bool `yaml:"no_propagation,omitempty"`
	// SharedCode: o arquivo declarou `@anchors-shared-code` — seus códigos de cenário
	// pertencem a outra unidade de propósito; não contam como colisão de identidade.
	SharedCode bool `yaml:"shared_code,omitempty"`

	// Needs — os planos que precisam TERMINAR antes deste começar (`needs:` no header).
	// Guardado no NÓ, e não só como aresta, porque um `needs` para plano inexistente não
	// vira aresta nenhuma — e é justamente esse caso que o doctor precisa reportar.
	Needs []string `yaml:"needs,omitempty"`
	// Parent — o CÓDIGO do artefato que CONTÉM este (`parent:` no header).
	//
	// Distinto de `Needs`, que é ordem: a fase 2 vem DEPOIS da fase 1 e não está DENTRO
	// dela. Com só uma das duas relações, quem monta a árvore precisa inferir a outra — e
	// inferir pertencimento a partir de ordem encaixa as fases de um plano numa escada.
	Parent string `yaml:"parent,omitempty"`
	// Revises — os planos que ESTE revisa (`revises:` no header).
	Revises []string `yaml:"revises,omitempty"`
	// Signal — sinais de qualidade INGERIDOS do runner (o Anchors não roda o teste;
	// consome o artefato que o projeto já gera). Preenchido por `anchors ingest`.
	Signal *TestSignal `yaml:"signal,omitempty"`
}

// TestSignal são os sinais de qualidade de teste amarrados a um nó (código ou teste).
// Ingeridos de artefatos padrão (JUnit para execução, lcov para cobertura). Levam a
// rev do nó no momento da ingestão — ficam STALE se o arquivo mudar depois.
type TestSignal struct {
	// Execução (de JUnit): quantos casos passaram/falharam sobre este nó de teste.
	// Os totais são a soma de TODAS as camadas (ByLayer) — mantidos para os checks
	// que não se importam com camada.
	Passed  int `yaml:"passed,omitempty"`
	Failed  int `yaml:"failed,omitempty"`
	Skipped int `yaml:"skipped,omitempty"`
	// ByLayer: execução por CAMADA de teste (unit/integration/e2e…), para o app de referência e
	// afins que rodam suites separadas. A chave é o --layer da ingestão. Permite o
	// relatório mergeado rotulado por camada, e a regra "crítico exige e2e".
	ByLayer map[string]LayerExec `yaml:"by_layer,omitempty"`
	// Cobertura de linha (de lcov): % de linhas cobertas deste nó de código.
	CoveredLines int     `yaml:"covered_lines,omitempty"`
	TotalLines   int     `yaml:"total_lines,omitempty"`
	LineCoverage float64 `yaml:"line_coverage,omitempty"` // 0..100
	// PrevLineCoverage: a cobertura de linha da ingestão ANTERIOR — o baseline para o
	// delta ("a cobertura caiu?"). Preservado ao ingerir por cima. -1 = sem baseline.
	PrevLineCoverage float64 `yaml:"prev_line_coverage,omitempty"`
	// MUTAÇÃO (do formato Mutation Testing Elements, schemaVersion 1.x): quantos
	// mutantes o teste MATOU. É a única medida objetiva de "o teste prova algo": um
	// mutante SOBREVIVENTE é uma alteração no código que os testes não perceberam —
	// logo, uma linha que ninguém prova. Complementa a cobertura de linha, que só diz
	// que a linha EXECUTOU, não que alguém verificou o resultado.
	//
	// Ingerido de artefato padrão, como JUnit e lcov: o Anchors não roda mutação e não
	// conhece ferramenta. Stryker (JS), PIT (Java), Infection (PHP) e mutmut (Python)
	// emitem o mesmo schema.
	MutantsKilled   int     `yaml:"mutants_killed,omitempty"`
	MutantsSurvived int     `yaml:"mutants_survived,omitempty"`
	MutationScore   float64 `yaml:"mutation_score,omitempty"` // 0..100
	// MutantsNoCoverage são os que NENHUM teste executou. Fora do score: "existe teste
	// que execute esta linha?" é pergunta do gate de COBERTURA. Ficam gravados porque a
	// informação é útil — só não é deste gate.
	MutantsNoCoverage int `yaml:"mutants_no_coverage,omitempty"`
	// MutantsIgnored são os que a ferramenta descartou antes de rodar (`ignoreStatic`,
	// `disable`). Ficam FORA do score — não houve experimento —, e existem no sinal para
	// separar dois 100% que significam coisas diferentes: "tudo foi provado" e "não havia
	// o que provar". Um arquivo de tabela de constantes cai no segundo.
	MutantsIgnored int `yaml:"mutants_ignored,omitempty"`
	// MutationByScope: o mesmo arquivo medido em dois ESCOPOS de suíte.
	//
	//   isolated — só o teste da própria unidade
	//   full     — os testes de todos os que a importam
	//
	// Os dois números juntos dizem o que nenhum diz sozinho. Medido em dois
	// projetos reais: um átomo de UI marcou 8% isolado e 77% completo, e uma função
	// pura de negócio marcou 27% e 77%. Nos dois, a maior parte dos mutantes só
	// morre nos DEPENDENTES — ou seja, o teste da unidade executa o código e quase
	// não verifica o resultado; quem prova são os outros.
	//
	// A diferença (`full - isolated`) mede o quanto a unidade depende de terceiros
	// para se provar. Num módulo bem testado ela tende a zero.
	MutationByScope map[string]MutationScope `yaml:"mutation_by_scope,omitempty"`
	// Os limiares que a FERRAMENTA DO PROJETO declarou, ingeridos junto com a medida.
	// `thresholds` e campo obrigatorio do schema Mutation Testing Elements (com `high`
	// e `low` obrigatorios dentro), entao lê-los daqui e tao agnostico quanto ler o
	// status dos mutantes — e evita que o mesmo numero exista em dois lugares.
	//
	//   MutationLow  — o minimo ACEITAVEL. Abaixo, o gate reprova.
	//   MutationHigh — o DESEJAVEL. Entre os dois, passa e aparece no relatorio.
	//
	// Zero = a ferramenta nao declarou; o gate cai no default do engine.
	MutationLow  float64 `yaml:"mutation_low,omitempty"`
	MutationHigh float64 `yaml:"mutation_high,omitempty"`
	// Códigos de cenário PROVADOS: os que aparecem em um caso de teste que PASSOU.
	// É a cobertura semântica (qual requisito da spec tem teste verde).
	ProvenCodes []string `yaml:"proven_codes,omitempty"`
	// AtRev: a rev do nó quando o sinal foi ingerido (para detectar staleness).
	AtRev string `yaml:"at_rev,omitempty"`
	// ClosureRev: a rev de cada nó do FECHO deste teste no momento da ingestão — o que ele
	// alcança descendo pelas arestas de saída (utils que compõe, unidade que importa).
	//
	// Sem isto, o frescor só enxergava o próprio arquivo de teste, e o custo foi medido:
	// `utils/login.yaml` é composto por 290 roteiros; ao tocá-lo, a evidência dos 290 devia
	// vencer, e nenhum sinal deles muda — o teste não mudou, mudou o que ele executa.
	ClosureRev map[string]string `yaml:"closure_rev,omitempty"`
	IngestedAt string            `yaml:"ingested_at,omitempty"`
}

// LayerExec é o resultado de execução de UMA camada de teste sobre um nó.
// MutationScope é o resultado de mutação de UM escopo de suíte sobre um nó.
type MutationScope struct {
	Killed     int     `yaml:"killed,omitempty"`
	Survived   int     `yaml:"survived,omitempty"`
	NoCoverage int     `yaml:"no_coverage,omitempty"`
	Ignored    int     `yaml:"ignored,omitempty"`
	Score      float64 `yaml:"score,omitempty"` // 0..100
	// AtRev é a rev do nó quando ESTE escopo foi medido — e ele é por escopo, não do
	// sinal inteiro, porque os dois são ingeridos em momentos diferentes.
	//
	// MEDIDO em 23/08: com um carimbo só, reingerir o `full` renovava o carimbo do
	// sinal e o `isolated` de uma medição anterior pegava carona como se fosse atual.
	// O gate julga PELO ISOLADO, então o veredito saía de um número velho — no caso,
	// medido contra código que nem existia mais.
	AtRev string `yaml:"at_rev,omitempty"`
}

// Stale diz se ESTE escopo foi medido contra outra versão do arquivo. Sem AtRev
// (sinal antigo, gravado antes deste campo existir) não dá para afirmar que está
// velho — e acusar staleness sem base faria o gate pedir remedição do que talvez
// esteja correto.
func (m MutationScope) Stale(rev string) bool {
	return m.AtRev != "" && m.AtRev != rev
}

type LayerExec struct {
	Passed  int `yaml:"passed,omitempty"`
	Failed  int `yaml:"failed,omitempty"`
	Skipped int `yaml:"skipped,omitempty"`
}

// Stamp — o carimbo de validação de uma aresta. Conjunto COMPLETO (fonte-única
// aqui, TRACEABILITY §4): a Propagação lê estes campos para calcular staleness.
type Stamp struct {
	ValidatedFromRev string `yaml:"validated_from_rev,omitempty"`
	ValidatedToRev   string `yaml:"validated_to_rev,omitempty"`
	// ChangedAt: DESDE QUANDO esta relação está como está — a data em que a rev de uma
	// ponta ou o veredito mudaram pela última vez.
	//
	// Não é "quando foi validado": o confronto acontece muitas vezes por dia, e registrar
	// cada um fazia o mapa mudar sozinho sem dizer nada. Confrontar de novo e achar o
	// mesmo resultado não é fato novo; a MUDANÇA é.
	ChangedAt string `yaml:"changed_at,omitempty"`
	Verdict   string `yaml:"verdict,omitempty"` // ok | issue | pending

	// Gate: qual gate produziu este veredito. Vazio no carimbo do `check` (que
	// confronta muitos gates de uma vez e resume o pior veredito); preenchido pelo
	// `anchors judge`, que julga UM gate sobre UM alvo.
	Gate string `yaml:"gate,omitempty"`
}

// Julgamento é o veredito de UM gate de julgamento sobre esta aresta.
//
// Campo SEPARADO do `Stamp`, e a separação é necessária: o `check` reescreve o Stamp
// inteiro a cada rodada (ele resume o pior veredito de todos os gates daquele nó),
// então um veredito de IA guardado ali era apagado no primeiro check seguinte — e o
// julgamento voltava a ser perguntado como se ninguém tivesse lido.
//
// Guarda as revs das pontas pelo mesmo motivo do Stamp: o veredito envelhece se o
// alvo mudar, e aí volta a ser pergunta.
type Julgamento struct {
	Gate             string `yaml:"gate"`
	Verdict          string `yaml:"verdict"` // ok | issue
	ValidatedFromRev string `yaml:"validated_from_rev,omitempty"`
	ValidatedToRev   string `yaml:"validated_to_rev,omitempty"`
	ChangedAt        string `yaml:"changed_at,omitempty"`
}

// Edge — uma relação dirigida (identidade ou dependência).
type Edge struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	// Julgamentos: um por gate de julgamento que já respondeu sobre esta aresta.
	// Preservado entre rodadas do `check` — ver o tipo Julgamento.
	Julgamentos []Julgamento `yaml:"julgamentos,omitempty"`
	Type        EdgeType     `yaml:"type"`
	Origin      Origin       `yaml:"origin"`
	// Method — metadado da aresta `depends-on`: o método/símbolo consumido do alvo
	// (SPEC_TYPES §5). Vazio nas demais arestas. Habilita impacto fino ("mudou signIn
	// → só as telas que usam signIn").
	Method string `yaml:"method,omitempty"`
	// Dep — o código local (DEPn) da linha da Tabela de Dependências que originou esta
	// aresta. Vazio nas demais arestas. Liga a aresta de volta à célula do data contract.
	Dep   string `yaml:"dep,omitempty"`
	Stamp *Stamp `yaml:"stamp,omitempty"`
}

// Graph — o mapa material inteiro. É o anchors.graph.yaml.
type Graph struct {
	Version int `yaml:"version"`
	// GeradoPor: a versão do binário que escreveu este mapa. Ver mapx.GeradoPor —
	// é o que permite ao `check` acusar um binário mais VELHO que o mapa, que
	// silenciosamente desfaz o que a versão nova escreveu.
	GeradoPor string `yaml:"gerado_por,omitempty"`
	Nodes     []Node `yaml:"nodes"`
	Edges     []Edge `yaml:"edges"`
}

// Stale devolve true se a aresta está desatualizada dado o estado atual dos nós.
// Regra de staleness (PROPAGATION §3): stale quando qualquer ponta avançou de rev
// desde o último confronto, ou quando nunca foi validada.
func (g *Graph) Stale(e Edge) bool {
	if e.Stamp == nil {
		return true // nunca validada
	}
	from := g.nodeRev(e.From)
	to := g.nodeRev(e.To)
	return from != e.Stamp.ValidatedFromRev || to != e.Stamp.ValidatedToRev
}

func (g *Graph) nodeRev(id string) string {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n.Rev
		}
	}
	return ""
}

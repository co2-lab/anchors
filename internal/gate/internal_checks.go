package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// Checkers internos: verificações que o CLI faz lendo TEXTO (não invoca ferramenta
// externa). São a instância "determinística interna" do medidor (QUALITY §3). Cada
// um recebe o conteúdo do alvo e devolve veredito + detalhe.
//
// Este registro é extensível: um gate no anchors.yaml com `check: <nome>` roteia
// para o checker de nome correspondente aqui.
var internalCheckers = map[string]func(content string, n mapx.Node) (Verdict, string){
	"non-empty":           checkNonEmpty,
	"has-code":            checkHasScenarioCode,
	"guide-has-checklist": checkGuideHasChecklist,
	"line-coverage":       checkLineCoverage,
	"coverage-delta":      checkCoverageDelta,
	"mutation-score":      checkMutationScore,
	"tests-pass":          checkTestsPass,
	"header-valid":        checkHeaderConforms,
	"route-declared":      checkRouteDeclared,
}

// checkersWithRoot são checkers que precisam da RAIZ do projeto (ex.: para invocar
// git), além do conteúdo. Roteados à parte dos checkers puros de conteúdo.
var checkersWithRoot = map[string]func(content string, n mapx.Node, root string) (Verdict, string){
	"updated-at-atual": checkUpdatedAt,
}

// checkersWithGraph são checkers RELACIONAIS: confrontam um nó contra seus vizinhos no
// mapa (arestas). Recebem o grafo, a raiz (para ler os arquivos vizinhos) e a config (p/
// o de-para de regimes e as superfícies da trinca — STRUCTURE §2.3). É a classe de gates
// que atravessa a trinca — ex.: feature↔test (cada cenário da feature está implementado
// no teste ligado, roteado pelo regime do cenário?).
var checkersWithGraph = map[string]func(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string){
	"progress-honest":          checkProgressHonest,
	"plan-doctrine-exists":     checkPlanDoctrineExists,
	"doctrine-realized":        checkDoctrineRealized,
	"scenario-coverage":        checkScenarioCoverage,
	"revision-orphans":         checkRevisionOrphans,
	"feature-spec-match":       checkFeatureSpecMatch,
	"test-feature-match":       checkTestFeatureMatch,
	"flag-scenario-grammar":    checkFlagScenarioGrammar,
	"flag-scenarios-complete":  checkFlagScenariosComplete,
	"flag-scenario-exists":     checkFlagScenarioExists,
	"flag-covered":             checkFlagCovered,
	"spec-doctrine-exists":     checkSpecDoctrineExists,
	"doctrine-not-duplicated":  checkDoctrineNotDuplicated,
	"spec-realizes-doctrine":   checkSpecRealizesDoctrine,
	"failure-handled":          checkFailureHandled,
	"failure-logged":           checkFailureLogged,
	"failure-declared":         checkFailureDeclared,
	"feature-test-match":       checkFeatureTestMatch,
	"scenario-identity":        checkScenarioIdentity,
	"scenario-type-aligned":    checkScenarioTypeAligned,
	"scenario-letter-declared": checkScenarioLetterDeclared,
	"spec-feature-match":       checkSpecFeatureMatch,
	"code-reference-valid":     checkCodeReferenceValid,
	"scenario-asserts":         checkScenarioAsserts,
	"domain-declared":          checkDomainDeclared,
	// Precisa de `cfg` para ler a própria opção `enforce_section_language` — ver checkSpecSections.
	"spec-sections":            checkSpecSections,
	"count-honored":            checkCountHonored,
	"trigger-declared":         checkTriggerDeclared,
	"route-exists":             checkRouteExists,
	"placeholder-filled":       checkPlaceholderFilled,
	"rule-implemented":         checkRuleImplemented,
	"vr-baseline":              checkVRBaseline,
	"ref-resolves":             checkRefResolves,
	"layer-boundary":           checkLayerBoundary,
	"dependency-honored":       checkDependencyHonored,
	"contract-status-declared": checkContractStatusDeclared,
	"proof-crosses-boundary":   checkProofCrossesBoundary,
	"triad-complete":           checkTriadComplete,
	"plan-seeds-valid":         checkPlanSeedsValid,
	"plan-source-declared":     checkPlanSourceDeclared,
	"phase-ordered":            checkPhaseOrdered,
	"phase-exists":             checkPhaseExists,
	"parent-valid":             checkParentValid,
	"plan-revised":             checkPlanRevised,
	"plan-change-justified":    checkPlanChangeJustified,
	"obligation-honored":       checkObligationHonored,
	"sibling-guard":            checkSiblingGuard,
	"pagination-honored":       checkPaginationHonored,
	"open-questions-resolved":  checkOpenQuestions,
	"rule-types":               checkRuleTypes,
	"identity-consistent":      checkIdentityConsistent,
	"region-pair-honored":      checkRegionPairHonored,
	"evidence-fresh":           checkEvidenceFresh,
	"testid-consistent":        checkTestIDCoherent,
	"testid-queried-exists":    checkQueriedTestIDExists,
	"mock-typed":               checkMockTyped,
	"mock-stamped":             checkMockStamped,
	"test-traceable":           checkTestTraceable,
	"code-cataloged":           checkCodeCataloged,
	"value-anchored":           checkValueAnchored,
	"docs-fresh":               checkDocsFresh,
	"docs-covered":             checkDocsCovered,
	"doc-required":             checkDocRequired,
	"doc-self-contained":       checkDocSelfContained,
}

func runInternal(name string, n mapx.Node, root string, graph *mapx.Graph, cfg *config.Config) (Verdict, string) {
	content, err := os.ReadFile(filepath.Join(root, n.ID))
	if err != nil {
		return Fail, i18n.T("gate.read_file_failed", err.Error())
	}
	if fn, ok := checkersWithGraph[name]; ok {
		return fn(string(content), n, root, graph, cfg)
	}
	if fn, ok := checkersWithRoot[name]; ok {
		return fn(string(content), n, root)
	}
	fn, ok := internalCheckers[name]
	if !ok {
		return Pending, i18n.T("gate.checker_not_implemented", name)
	}
	return fn(string(content), n)
}

// checkersComGate: os agregados que precisam saber QUAL instância os invocou.
//
// Os demais checkers respondem a mesma pergunta sempre, e por isso `cfg` lhes basta. Um
// gate GENÉRICO não: `marker-parity` só sabe o que procurar depois de ler o próprio
// `marker_prefix`, e um projeto declara vários deles. Sem o `config.Gate` aqui, todas as
// instâncias veriam a mesma configuração — ou nenhuma.
var checkersComGate = map[string]func(g config.Gate, root string, graph *mapx.Graph, cfg *config.Config) (Verdict, string){
	"marker-parity": checkMarkerParity,
	// AGREGADO: um veredito por DOCUMENTO, não por unidade. A versão por-nó produziu 37
	// cards que escreviam nos mesmos dois arquivos, e cada merge invalidava os outros 36.
	"doc-required": checkDocRequiredAggregate,
}

// runInternalAggregate executa um checker interno de escopo batch/project.
//
// Difere do `runInternal` num ponto que não é detalhe: NÃO há arquivo para ler. O
// escopo é o conjunto, então o checker recebe um nó vazio e se orienta por `root` e
// `cfg`. Tentar ler o conteúdo de um alvo aqui (como o runInternal faz) devolveria
// erro de leitura e o gate reprovaria por um arquivo que nunca existiu.
func runInternalAggregate(g config.Gate, root string, graph *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if fn, ok := checkersComGate[g.Check]; ok {
		return fn(g, root, graph, cfg)
	}
	fn, ok := checkersWithGraph[g.Check]
	if !ok {
		return Pending, i18n.T("gate.aggregate_checker_unknown", g.Check)
	}
	return fn("", mapx.Node{}, root, graph, cfg)
}

// non-empty: o arquivo não é vazio nem só espaço. Trivial, mas pega placeholders.
func checkNonEmpty(content string, _ mapx.Node) (Verdict, string) {
	if strings.TrimSpace(content) == "" {
		return Fail, i18n.T("gate.empty_file")
	}
	return Pass, ""
}

// spec-sections: uma spec deve CATALOGAR ao menos uma regra/estado com código — o
// sinal mínimo de que não é um esqueleto vazio (SPEC §8). Um item catalogado conta em
// qualquer forma ESTRUTURADA com identidade:
//   - cabeçalho:      `### CODE-...`
//   - linha de tabela: `| CODE-... |`  (ou `| \`CODE-...\` |`)
//   - bullet-negrito:  `- **CODE-...**` (marcação deliberada, não prosa)
//
// O que NÃO conta é uma menção SOLTA em prosa (sem heading, tabela nem negrito). O
// gate mede substância (a regra tem código e um lugar), não a forma exata — o app de referência usa
// os três formatos e todos cumprem o espírito.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func specHeadingRE() *regexp.Regexp {
	return regexp.MustCompile("(?m)^###\\s+[A-Z0-9]" + config.CodeLengthPattern() + "-")
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func specTableRE() *regexp.Regexp {
	return regexp.MustCompile("(?m)^\\|\\s*`?[A-Z0-9]" + config.CodeLengthPattern() + "-")
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func specBoldBulletRE() *regexp.Regexp {
	return regexp.MustCompile("(?m)^\\s*[-*]\\s*\\*\\*`?[A-Z0-9]" + config.CodeLengthPattern() + "-")
}

// updated-at-atual: o `updated_at` do header bate com a data do ÚLTIMO COMMIT que
// tocou o arquivo (só ANO-MÊS-DIA — nunca hora; editar num horário e commitar noutro
// do mesmo dia é OK). A manutenção do campo é de quem alterou; este gate confronta e
// força a correção quando o dia diverge — a âncora não pode mentir sobre quando mudou.
// Reparável por `anchors check --fix` (escreve a data do commit no header).
var updatedAtRE = regexp.MustCompile(`updated_at:\s*(\d{4}-\d{2}-\d{2})`)

func checkUpdatedAt(content string, n mapx.Node, root string) (Verdict, string) {
	m := updatedAtRE.FindStringSubmatch(content)
	if m == nil {
		return Skip, "" // sem updated_at declarado — o gate header-valid cobra o header
	}
	declared := m[1]

	// EDIÇÃO NÃO-COMMITADA: a "última alteração" é AGORA, não o último commit (que é
	// de um dia anterior). Desempate pela data do SISTEMA: o header deve dizer HOJE
	// (o dia da alteração em curso). Se disser, ok; senão, o autor esqueceu de
	// atualizar para o dia em que está mexendo.
	mudou, sabido := gitmeta.UncommittedChanges(root, n.ID)
	if !sabido {
		// Sem git não há como saber se o arquivo mudou NEM qual foi o último commit.
		// Antes isto caía no ramo "commitado" e devolvia "arquivo sem commit
		// (novo/untracked)" — uma afirmação FALSA e específica, que mandava o autor
		// investigar o arquivo quando o que falta é o repositório.
		return Skip, i18n.T("gate.updated_at.no_git")
	}
	if mudou {
		today := gitmeta.Today()
		if declared == today {
			return Pass, "" // atualizado para hoje (a alteração em curso)
		}
		return Fail, i18n.T("gate.updated_at.uncommitted_diff", declared, today, today)
	}

	// COMMITADO: compara com a data do último commit que tocou o arquivo (só o dia).
	gitDate, ok := gitmeta.LastCommitDate(root, n.ID)
	if !ok {
		return Pending, i18n.T("gate.updated_at.no_commit")
	}
	if declared != gitDate {
		return Fail, i18n.T("gate.updated_at.diff", declared, gitDate, gitDate)
	}
	return Pass, ""
}

// header-valid: o arquivo tem o BLOCO DE CABEÇALHO do Anchors (`@anchors`) com o
// mínimo obrigatório — a identidade (`code:`). O guide de header (anchors guide
// header) é a régua; este gate verifica presença + mínimo. Agnóstico de dialeto de
// comentário (só procura os marcadores no texto). Detalhes extras (updated_at, tags)
// são opcionais — o gate só cobra o piso.
var headerBlockRE = regexp.MustCompile(`@anchors\b`)

// identidade no header: `code:` (posse — o dono, ex.: a spec) OU `ref:` (referência —
// o resto da trinca aponta o código da unidade que realiza/cobre/prova). Um dos dois
// é obrigatório para camadas REGIDAS; qual depende do papel do arquivo. Camadas
// RECONHECIDAS (sem spec) usam `layer:` como identidade mínima (ver `anchors guide header`).
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func headerCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*[A-Z0-9]` + config.CodeLengthPattern() + `\b`)
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func headerRefRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*ref:\s*[A-Z0-9]` + config.CodeLengthPattern() + `\b`)
}

var headerLayerRE = regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*layer:\s*\S+`)

// recognizedLayers são as tags de camadas RECONHECIDAS — declaradas na Estrutura só
// para sair do escrutínio de spec (dao/infra/presentation/domain-types e afins). Um
// arquivo dessas camadas não tem spec dona nem irmã a referenciar; sua identidade
// mínima honesta é o `layer:` (a que layer pertence), não um code/ref inventado.
var recognizedLayers = map[string]bool{
	"dao": true, "infra": true, "presentation": true, "domain-types": true,
}

// headerLayerValueRE captura o VALOR do `layer:` declarado no header.
var headerLayerValueRE = regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*layer:\s*(\S+)`)

// isRecognizedLayer decide se o arquivo pertence a uma layer reconhecida — pela tag
// do nó OU pelo `layer:` DECLARADO no header. O header importa porque um TESTE de um
// arquivo de layer reconhecida é classificado como `kind: test` (o pattern de test é
// mais específico que o da layer), perdendo a tag da layer; mas ele declara
// `layer: <reconhecida>` no header, e o header é a fonte da verdade da identidade.
func isRecognizedLayer(n mapx.Node, content string) bool {
	// A Estrutura do projeto é a fonte da verdade: uma layer com `regime: declarativo`
	// é RECONHECIDA (não origina regra), qualquer que seja seu nome — o vocabulário de
	// camadas pertence ao projeto (anchors.yaml), não ao engine. Ver STRUCTURE.md.
	if n.Regime == "declarativo" {
		return true
	}
	// REGIME DECLARADO E DIFERENTE encerra a decisão: o projeto disse o que a layer é, e
	// o fallback por nome não pode contradizê-lo.
	//
	// Sem esta linha, a Estrutura era fonte da verdade só na direção que DISPENSA. Uma
	// layer `infra` que declara `regime: regra` continuava caindo no fallback — e ali
	// `infra` é um dos nomes canônicos de layer declarativa (DAO, adaptador, o que não
	// decide nada).
	//
	// Medido no projeto de referência: `packages/infra/` é CDK COM regra — três unidades
	// com 19 regras e invariantes somados, 61 testes, toda mutação detectada — e o
	// `triad-complete` respondia INDETERMINADO nas três, por causa do NOME da layer. Nem
	// passa nem acusa, e parece cobertura.
	if n.Regime != "" {
		return false
	}
	// Fallback aos nomes canônicos, p/ projetos que ainda não declaram `regime:`.
	for _, t := range n.Tags {
		if recognizedLayers[t] {
			return true
		}
	}
	if m := headerLayerValueRE.FindStringSubmatch(content); m != nil {
		return recognizedLayers[strings.TrimSpace(m[1])]
	}
	return false
}

// isRecognizedLayerCfg decide o mesmo, consultando a ESTRUTURA para a layer que a spec
// declara no header.
//
// A distinção existe porque o nó de uma SPEC tem `layer: spec` — ela casa `**/*.spec.md`,
// e o `Regime` copiado para o nó é o da layer `spec`, não o da unidade. O regime que o
// projeto declarou para `infra` nunca chegava aqui, e a decisão caía no fallback por nome.
//
// Medido: `packages/infra/` do projeto de referência é CDK COM regra — três unidades, 19
// regras e invariantes somados, 61 testes —, o projeto declarou `regime: comportamental`
// na layer, e o `triad-complete` seguiu respondendo INDETERMINADO. Nem passa nem acusa, e
// parece cobertura.
func isRecognizedLayerCfg(n mapx.Node, content string, cfg *config.Config) bool {
	if n.Regime == "declarativo" {
		return true
	}
	// A layer da UNIDADE, que a spec declara no header — e o regime QUE O PROJETO deu a
	// ela. Declarado e diferente de `declarativo` encerra a decisão: a Estrutura é fonte
	// da verdade nas duas direções, não só na que dispensa.
	if cfg != nil {
		if m := headerLayerValueRE.FindStringSubmatch(content); m != nil {
			if l, ok := cfg.Layers[strings.TrimSpace(m[1])]; ok && l.Regime != "" {
				return l.Regime == "declarativo"
			}
		}
	}
	return isRecognizedLayer(n, content)
}

// isExecutableScript — o nó é um roteiro do runner e2e (não um arquivo de teste em
// linguagem de programação). Reconhecido pela extensão de dado + kind test: um `.test.tsx`
// é código nosso e carrega cabeçalho; um `.yaml` é entrada de um runner externo.
func isExecutableScript(n mapx.Node) bool {
	if n.Kind != mapx.KindTest {
		return false
	}
	return strings.HasSuffix(n.ID, ".yaml") || strings.HasSuffix(n.ID, ".yml")
}

func checkHeaderConforms(content string, n mapx.Node) (Verdict, string) {
	// Arquivo BINÁRIO não carrega cabeçalho — não há sintaxe de comentário num PNG.
	// A identidade dele está no NOME (`<Unidade>.<CODE>-VR-<variante>.png`), que é o
	// que o `identity-consistent` confronta. Cobrar header aqui exigiria o impossível
	// e barraria todo commit de baseline visual.
	if isBinary(content) {
		return Skip, i18n.T("gate.header.binary")
	}
	// ROTEIRO de teste executável (.yaml do runner e2e): mesma razão do binário acima,
	// por um caminho diferente. A identidade dele está no NOME do arquivo
	// (`<CODE>-A01.yaml`, que é como o `tested-by` o encontra) e nas `tags:`, que o
	// próprio runner consome — e o formato é do RUNNER, não nosso: um bloco `@anchors`
	// no topo de um YAML de flow é comentário morto para quem o executa.
	//
	// A alternativa era retrofitar 717 arquivos com um cabeçalho que nenhum deles jamais
	// teve — e o gate só passou a vê-los porque a layer `e2e-flow` os trouxe para o
	// grafo (para que a execução de E2E pudesse deixar carimbo). Cobrar deles um contrato
	// escrito depois seria transformar a chegada ao mapa em 717 defeitos retroativos.
	if isExecutableScript(n) {
		return Skip, i18n.T("gate.header.executable_script")
	}
	if !headerBlockRE.MatchString(content) {
		return Fail, i18n.T("gate.header.missing_block")
	}
	// Camada RECONHECIDA (sem spec): `layer:` é a identidade mínima suficiente.
	if isRecognizedLayer(n, content) {
		if !headerLayerRE.MatchString(content) &&
			!headerCodeRE().MatchString(content) && !headerRefRE().MatchString(content) {
			return Fail, i18n.T("gate.header.recognized_missing_id")
		}
		return Pass, ""
	}
	// Camada REGIDA (ou spec): exige posse ou referência.
	if !headerCodeRE().MatchString(content) && !headerRefRE().MatchString(content) {
		return Fail, i18n.T("gate.header.governed_missing_id")
	}
	return Pass, ""
}

func checkSpecSections(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	// O PLACEHOLDER não é cobrado aqui, e antes era: este gate casava as frases dos
	// templates em português ("[Descreva aqui]", "TODO: descrever"), o que quebraria na
	// tradução — e, pior, duplicava o gate `placeholder-preenchido`, que faz o mesmo
	// confronto com o vocabulário universal (`TODO`/`FIXME`/`XXX`/`campo: <valor>`) e
	// sabe distinguir o marcador real do exemplo de sintaxe.
	//
	// Dois gates sobre o mesmo defeito produzem dois achados para um problema, e quem
	// corrige um continua barrado pelo outro sem entender por quê.
	catalogued := specHeadingRE().MatchString(content) ||
		specTableRE().MatchString(content) ||
		specBoldBulletRE().MatchString(content)
	if !catalogued {
		return Fail, i18n.T("gate.spec_sections.no_rules")
	}
	// "Ao menos uma" é o piso, e sozinho ele deixa passar o caso mais comum: a spec que
	// cataloga a primeira regra e escreve as outras nove em prosa. As nove ficam
	// invisíveis para TODOS os gates de identidade — não por estarem erradas, mas por não
	// terem código, e um gate que não enxerga a regra reporta verde sobre o que não
	// conferiu.
	if msg := siblingsWithoutCode(content); msg != "" {
		return Fail, msg
	}
	if msg := sectionsInWrongLanguage(content, cfg); msg != "" {
		return Fail, msg
	}
	return Pass, ""
}

// sectionsInWrongLanguage acha títulos de seção escritos em idioma DIFERENTE do `lang:`
// do projeto.
//
// Por que é um defeito, já que os confrontos de seção aceitam todos os idiomas: aceitar
// todos é o que impede a tradução do catálogo de quebrar acervo existente — decisão
// deliberada, e ela não deve virar licença para o acervo MISTO. Metade das specs com
// `## Visão Geral` e metade com `## Overview` é uma spec sem forma única, e nenhum outro
// gate reclama disso (cada um, sozinho, acha o que procura).
//
// Só reclama do que ele RECONHECE: um título fora do catálogo (`## Fora de escopo`, o
// léxico próprio do projeto) não é idioma errado, é seção que o framework não nomeia —
// acusá-la seria cobrar do projeto o vocabulário do engine.
func sectionsInWrongLanguage(content string, cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	// A opção mora no gate que hospeda este check. Sem instância declarada (o caso do
	// `--all` sintético), o default do framework vale: cobra.
	for _, g := range cfg.Gates {
		if g.Check == "spec-sections" && !g.ChecksSectionLanguage() {
			return ""
		}
	}
	lang := cfg.Lang
	if lang == "" {
		lang = i18n.Default
	}
	var fora []string
	for _, titulo := range headingTitles(content) {
		chave, idioma := i18n.SectionKeyFor(titulo)
		if chave == "" || idioma == lang {
			continue
		}
		esperado := i18n.TIn(lang, chave)
		if esperado == "" || strings.EqualFold(esperado, titulo) {
			continue
		}
		fora = append(fora, fmt.Sprintf("%q (%s) → %q", titulo, idioma, esperado))
	}
	if len(fora) == 0 {
		return ""
	}
	return i18n.T("gate.spec_sections.wrong_language", lang, strings.Join(fora, "; "))
}

// headingTitles devolve o texto dos títulos `##`..`####` do conteúdo.
func headingTitles(content string) []string {
	var out []string
	for linha := range strings.SplitSeq(content, "\n") {
		t := strings.TrimSpace(linha)
		if !strings.HasPrefix(t, "##") {
			continue
		}
		if t = strings.TrimSpace(strings.TrimLeft(t, "#")); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// seçãoVazia diz se a seção não tem CONTEÚDO — só o esqueleto.
//
// Conta como vazio: linha em branco, separador de tabela (`|---|---|`), cabeçalho de
// tabela (a linha de rótulos que precede o separador) e o divisor `---`. Qualquer outra
// coisa é conteúdo, e aí a ausência de código volta a ser o defeito que o gate pega.
func emptySection(linhas []string) bool {
	for i, l := range linhas {
		t := strings.TrimSpace(l)
		if t == "" || t == "---" {
			continue
		}
		if sepTabelaRE.MatchString(t) {
			continue
		}
		// Cabeçalho de tabela: linha de `|` cujo SUCESSOR imediato é o separador. Sem
		// olhar o sucessor, uma linha de dados seria confundida com rótulo.
		if strings.HasPrefix(t, "|") && i+1 < len(linhas) &&
			sepTabelaRE.MatchString(strings.TrimSpace(linhas[i+1])) {
			continue
		}
		return false
	}
	return true
}

// noContentRE: a dispensa honesta da SEÇÃO. Exige motivo — `@no-content` sozinho não
// conta, pelo mesmo princípio de `@no-rule`: dispensa sem porquê é a dispensa que ninguém
// revisa depois.
var noContentRE = regexp.MustCompile(`@no-content[^\S\n]*:[^\S\n]*\S+`)

var sepTabelaRE = regexp.MustCompile(`^\|[\s:|-]+\|?$`)

// leveledSectionRE casa um cabeçalho e captura o nível (para agrupar por pai) e o título.
var leveledSectionRE = regexp.MustCompile(`^(#{2,4})\s+(.+?)\s*$`)

// siblingsWithoutCode acha a seção que DEVERIA catalogar e não cataloga.
//
// O difícil é separar "seção de regra sem código" de "seção de prosa" — `## Visão Geral`
// e `## Restrições` não catalogam nada e não devem ser cobradas. A pergunta não se
// responde pelo título (cada projeto nomeia como quer), e sim pela VIZINHANÇA: entre as
// seções irmãs sob o mesmo pai, se algumas catalogam com código e outras não, as sem
// código destoam do padrão que a própria spec estabeleceu.
//
// É o mesmo princípio do `unclaimedSections` no `rule-types`: a spec declara sua forma, e
// o gate cobra coerência com ela — não com um formato que o Anchors imponha.
func siblingsWithoutCode(content string) string {
	type sec struct {
		titulo string
		nivel  int
		linhas []string
	}
	var secoes []sec
	for _, l := range strings.Split(content, "\n") {
		if m := leveledSectionRE.FindStringSubmatch(l); m != nil {
			secoes = append(secoes, sec{titulo: m[2], nivel: len(m[1])})
			continue
		}
		if len(secoes) > 0 {
			secoes[len(secoes)-1].linhas = append(secoes[len(secoes)-1].linhas, l)
		}
	}
	// Agrupa por PAI: o `##` que precede cada bloco de `###`.
	//
	// O `####` NÃO entra no grupo — ele é conteúdo do `###` que o precede, não irmão
	// dele. Tratá-los como irmãos produzia um falso positivo silencioso: a linha de
	// conteúdo vai sempre para a ÚLTIMA seção vista, então o `####` roubava as linhas
	// do pai e o `###` aparecia vazio.
	//
	// Medido em `SplashScreen.spec.md`: `### Caminhos Condicionais` tem três linhas de
	// tabela, todas dentro de `#### destination`. O gate acusava o pai de estar vazio
	// enquanto o conteúdo estava logo abaixo.
	pai := ""
	grupos := map[string][]sec{}
	for i, s := range secoes {
		if s.nivel == 2 {
			pai = s.titulo
			continue
		}
		if s.nivel > 3 {
			continue // conteúdo do `###`, não irmão dele
		}
		// O `###` herda as linhas dos `####` que o seguem: são o corpo dele.
		filhas := s
		for j := i + 1; j < len(secoes) && secoes[j].nivel > 3; j++ {
			filhas.linhas = append(filhas.linhas, secoes[j].titulo)
			filhas.linhas = append(filhas.linhas, secoes[j].linhas...)
		}
		grupos[pai] = append(grupos[pai], filhas)
	}
	for _, irmas := range grupos {
		if len(irmas) < 2 {
			continue // sem irmã não há padrão a comparar
		}
		var comCodigo, semCodigo, semConteudo []string
		for _, s := range irmas {
			// O código pode estar no TÍTULO (`### CODE-B01 — ...`) ou no corpo (tabela,
			// bullet). As duas formas contam: o gate cobra identidade, não posição.
			corpo := s.titulo + "\n" + strings.Join(s.linhas, "\n")
			if anyCodeRE.MatchString(corpo) {
				comCodigo = append(comCodigo, s.titulo)
				continue
			}
			// SEÇÃO VAZIA não é seção sem código: não há regra alguma a codificar, e
			// mandar "dê um código a cada uma" pede o impossível.
			//
			// Mas vazia sozinha NÃO absolve — são dois estados que se parecem e não são
			// o mesmo:
			//   • vazia porque NÃO SE APLICA  → decisão, e o autor a declara
			//   • vazia porque ninguém escreveu → esquecimento, e o gate tem de pegar
			//
			// Aceitar as duas em silêncio trocaria um falso positivo por um falso
			// negativo, e o segundo é pior: o gate passaria a reportar verde sobre
			// seção que ninguém preencheu.
			//
			// A saída é a dispensa HONESTA que o resto do vocabulário já usa
			// (`@no-rule`, `@no-code`, `@no-scenario`, `@no-mark`): a seção fica, e
			// quem decidiu que ela não tem conteúdo escreve o motivo.
			//
			//     ### Comportamentos Automáticos
			//
			//     @no-content: tela estática — não há efeito, timer nem carga.
			//
			// A seção PERMANECE, que é o ponto: há seções obrigatórias por exigência
			// regulatória, e apagá-las para calar o gate destruiria a obrigação junto
			// com o aviso.
			// A declaração vem PRIMEIRO: ela é texto, e por isso a seção que a carrega
			// não conta como vazia. Perguntar "está vazia?" antes rejeitaria justamente
			// quem fez a coisa certa.
			if noContentRE.MatchString(strings.Join(s.linhas, "\n")) {
				continue
			}
			if emptySection(s.linhas) {
				semConteudo = append(semConteudo, s.titulo)
				continue
			}
			semCodigo = append(semCodigo, s.titulo)
		}
		// Só acusa quando a MAIORIA das irmãs cataloga: um grupo em que só uma tem código
		// não estabeleceu padrão nenhum, e cobrar as outras seria inventar um.
		// Seção VAZIA e sem declaração: o autor não disse se decidiu ou esqueceu, e o
		// gate não tem como saber. Pede a declaração — que custa uma linha e resolve
		// para sempre.
		if len(comCodigo) > 0 && len(semConteudo) > 0 && len(semCodigo) == 0 {
			return i18n.T("gate.siblings.empty_sections",
				strings.Join(semConteudo, ", "), strings.Join(comCodigo, ", "))
		}
		if len(comCodigo) > len(semCodigo) && len(semCodigo) > 0 {
			return i18n.T("gate.siblings.missing_code",
				strings.Join(semCodigo, ", "), strings.Join(comCodigo, ", "))
		}
	}
	return ""
}

// has-code: o arquivo carrega um código de cenário (a identidade). Sem código, a
// peça é órfã invisível (TRACEABILITY §6). Aqui só verificamos presença.
// A classe de letras vem do vocabulário do projeto (`rule_types`) — ver SetRuleLetters.
var anyCodeRE = anyCodeREFor(config.DefaultRuleLetters)

func anyCodeREFor(letters string) *regexp.Regexp {
	return regexp.MustCompile(`\b[A-Z0-9]` + config.CodeLengthPattern() + `-(?:[` + regexp.QuoteMeta(letters) + `]\d{2}|DS-|VR)`)
}

// SetRuleLetters reconfigura a gramática de código dos gates para o vocabulário do
// projeto. Chamado no início do pipeline de gates (RunWithConfig).
// SetRuleLetters reconfigura TODOS os regexes deste pacote que dependem do vocabulário
// de letras do projeto. Adicionar um regex de código sem registrá-lo aqui o deixa preso
// às letras canônicas — e um cenário de letra declarada pelo projeto vira invisível para
// aquele gate, que então reporta verde sobre o que não conferiu.
func SetRuleLetters(letters string) {
	anyCodeRE = anyCodeREFor(letters)
	featScenarioCodeRE = featCodeREFor(letters)
}

func checkHasScenarioCode(content string, _ mapx.Node) (Verdict, string) {
	if anyCodeRE.MatchString(content) {
		return Pass, ""
	}
	return Fail, i18n.T("gate.missing_scenario_code")
}

// guide-has-checklist: um GUIDE de governança deve destilar suas regras em PONTOS DE
// CONFORMIDADE verificáveis (a seção "## Pontos de conformidade" com itens CK1, CK2…).
// Sem eles, o gate de julgamento por IA recai em heurística vaga. Este checker é
// DETERMINÍSTICO — só verifica a PRESENÇA da seção e de ao menos um item CK; a
// QUALIDADE de cada ponto é julgamento (ver `anchors guide guide`). Ver a sugestão do
// meta-guide: "a presença da checklist é um check determinístico".
// A seção é reconhecida em QUALQUER idioma do catálogo, e não só em português.
//
// O acoplamento que isto desfaz: o `anchors init` semeia o HEADER_GUIDE.md com o título
// traduzido pelo `lang:` do projeto, e a regex antiga só casava a forma portuguesa — um
// projeto `lang: en` nasceria reprovando o `guide-checklist` no primeiro `check`, por um
// guide que o próprio init acabara de escrever. Aceitar todos os idiomas é o que mantém a
// tradução do guia e o gate falando da mesma coisa.
var checklistHeadingRE = regexp.MustCompile(`(?mi)^##+\s+(?:` +
	strings.Join(escapeAll(i18n.AllTranslations("section.title.compliance_points")), "|") + `)\b`)

// escapeAll prepara títulos traduzidos para entrar numa alternância de regex.
func escapeAll(xs []string) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		if x != "" {
			out = append(out, regexp.QuoteMeta(x))
		}
	}
	if len(out) == 0 {
		out = append(out, `pontos de conformidade`)
	}
	return out
}

var checklistItemRE = regexp.MustCompile(`(?m)\bCK\d+\b`)

// --- sinais de teste ingeridos (execução, cobertura) — QUALITY §5 Execução ---

// scenario-coverage: cada código de cenário que a spec DECLARA tem um teste que
// PASSOU (está em Signal.ProvenCodes)? Fecha o gate de Rastreabilidade — não basta
// existir teste, cada requisito precisa estar provado. Pending se nada foi ingerido.
func checkScenarioCoverage(content string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	// Só os requisitos DEFINIDOS por esta spec, não toda menção de código no texto.
	//
	// O `anyCodeRE` sobre o conteúdo inteiro casa também o que a spec CITA ao justificar
	// as regras dela — e uma spec bem escrita cita muito. Medido no blue-eyes: a
	// `GoLiveChecklist` tinha 6 requisitos e o gate cobrava 18 cenários, 15 deles de
	// outras unidades (`CRPNC-B03`, `MTTLM-B02`, `DTSTD-B06`…). Nenhum daqueles cenários
	// poderia ser provado por um teste desta unidade — o gate pedia o impossível, e a
	// mensagem sugeria que a spec estava mal coberta.
	//
	// O `definedRequirements` é a mesma função que o `spec-feature-match` usa: ela casa
	// TÍTULO DE SEÇÃO (`### CODE-B01 — …`) e respeita `@no-scenario`. Duas leituras do que
	// é "um requisito desta spec" divergiriam — e divergiam.
	declared := definedRequirements(content)
	if len(declared) == 0 {
		return Skip, "" // spec sem requisito definido — nada a cobrir
	}
	if n.SignalStale() {
		return Pending, i18n.T("gate.stale_test_signal")
	}

	// DUAS PERGUNTAS, e juntá-las perde a resposta das duas.
	//
	//   ESCRITO  algum teste NOMEIA este código de cenário   (estático, sempre respondível)
	//   VERDE    esse teste RODOU e PASSOU                   (exige execução ingerida)
	//
	// A versão anterior só perguntava a segunda, e num projeto que nunca ingeriu relatório
	// respondia Pending para tudo — "ninguém mediu" —, escondendo os cenários que ninguém
	// testou. E só a estática seria o erro oposto: teste escrito pode nunca ter rodado.
	//
	// É a mesma régua do `flag-covered`, e pela mesma razão: o conserto de "sem teste" é
	// escrever um; o de "escrito e não executado" é rodar a suíte. Um veredito que não
	// distingue os dois manda a pessoa pelo caminho errado metade das vezes.
	proven := map[string]bool{}
	ingested := n.Signal != nil
	if ingested {
		for _, c := range n.Signal.ProvenCodes {
			proven[c] = true
		}
	}
	written := codesNamedByTests(declared, root, g, n.ID)

	var noTest, notGreen []string
	seen := map[string]bool{}
	for _, code := range declared {
		if seen[code] {
			continue
		}
		seen[code] = true
		switch {
		case proven[code]:
			// provado: nada a cobrar
		case written[code]:
			notGreen = append(notGreen, code)
		default:
			noTest = append(noTest, code)
		}
	}
	if len(noTest) == 0 && len(notGreen) == 0 {
		return Pass, ""
	}

	var b strings.Builder
	if len(noTest) > 0 {
		b.WriteString(i18n.T("gate.scenario_missing_test", len(noTest), strings.Join(noTest, ", ")))
	}
	if len(notGreen) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		if !ingested {
			b.WriteString(i18n.T("gate.scenario_written_not_ingested", len(notGreen), strings.Join(notGreen, ", ")))
		} else {
			b.WriteString(i18n.T("gate.scenario_written_not_green", len(notGreen), strings.Join(notGreen, ", ")))
		}
	}
	return Fail, strings.TrimRight(b.String(), "\n")
}

// codesNamedByTests responde a metade ESTÁTICA: quais destes códigos algum teste nomeia.
//
// Percorre os testes LIGADOS a este nó (arestas `tested-by`) e, na falta delas, os testes
// do grafo — a flag não tem aresta para os seus testes, e a spec nem sempre tem. Comentário
// não conta, pela mesma régua do `feature-test-match`: código citado em comentário é
// REFERÊNCIA a outra unidade, não implementação.
func codesNamedByTests(codes []string, root string, g *mapx.Graph, id string) map[string]bool {
	written := map[string]bool{}
	if g == nil || root == "" {
		return written
	}
	var paths []string
	for _, e := range g.Neighbors(id).Out {
		if e.Type == mapx.EdgeTestedBy {
			paths = append(paths, e.To)
		}
	}
	if len(paths) == 0 {
		for _, node := range g.Nodes {
			if node.Kind == mapx.KindTest {
				paths = append(paths, node.ID)
			}
		}
	}
	for _, tp := range paths {
		b, err := os.ReadFile(filepath.Join(root, tp))
		if err != nil {
			continue
		}
		body := stripLineComments(string(b))
		for _, c := range codes {
			if strings.Contains(body, c) {
				written[c] = true
			}
		}
	}
	return written
}

// line-coverage: a cobertura de linha do nó de código está >= 70%? (limiar fixo por
// ora — poderia vir da config). Pending se não ingerida.
func checkLineCoverage(_ string, n mapx.Node) (Verdict, string) {
	if n.Signal == nil || n.Signal.TotalLines == 0 {
		return Pending, i18n.T("gate.no_line_coverage")
	}
	if n.SignalStale() {
		return Pending, i18n.T("gate.stale_coverage")
	}
	const threshold = 70.0
	if n.Signal.LineCoverage < threshold {
		return Fail, i18n.T("gate.line_coverage_low", n.Signal.LineCoverage, threshold)
	}
	return Pass, ""
}

// coverage-delta: a cobertura de linha deste nó de código CAIU desde a ingestão
// anterior? Pega a regressão de cobertura (código pode estar coberto no diff mas ter
// derrubado a cobertura de outra parte). Pending se não há baseline.
func checkCoverageDelta(_ string, n mapx.Node) (Verdict, string) {
	if n.Signal == nil || n.Signal.PrevLineCoverage == 0 {
		return Pending, i18n.T("gate.no_coverage_baseline")
	}
	if n.SignalStale() {
		return Pending, i18n.T("gate.stale_coverage")
	}
	delta := n.Signal.LineCoverage - n.Signal.PrevLineCoverage
	if delta < -0.01 {
		return Fail, i18n.T("gate.coverage_dropped", n.Signal.PrevLineCoverage, n.Signal.LineCoverage, delta)
	}
	return Pass, ""
}

// mutation-score: dos mutantes que rodaram neste arquivo, o teste matou o suficiente?
//
// Este é o único gate que responde "o teste PROVA a linha?" — todos os outros respondem
// "a linha executou?" ou "a peça existe?". Um mutante SOBREVIVENTE é uma alteração no
// código que a suíte não percebeu; logo, uma linha que ninguém prova. Já vimos, medido:
// arquivo com 16 testes verdes e 100% de cobertura onde apagar a guarda `!removed &&`
// não derrubava NENHUM teste.
//
// PENDING (não Fail) quando o projeto não ingeriu mutação: exigir a ferramenta seria o
// framework decidindo pelo projeto — nem toda linguagem/stack tem uma que rode em tempo
// viável. Mas Pending não é silêncio: aparece no status e no doctor com o que falta e o
// que se perde por não ter. A alavanca do projeto é declarar o gate como blocking no
// anchors.yaml quando ele já tiver a ferramenta rodando.
func checkMutationScore(_ string, n mapx.Node) (Verdict, string) {
	// `MutantsIgnored` entra na condição, e é o que separa "a ferramenta nunca rodou aqui"
	// de "ela rodou e não havia o que medir". Sem ele, um arquivo cujos mutantes foram
	// TODOS ignorados recebia "rode a ferramenta de mutação" — sobre um arquivo em que ela
	// já tinha rodado e feito a coisa certa. O pedido era impossível de atender: rodar de
	// novo daria o mesmo resultado.
	if n.Signal == nil || (n.Signal.MutantsKilled == 0 && n.Signal.MutantsSurvived == 0 &&
		n.Signal.MutantsIgnored == 0 && n.Signal.MutantsNoCoverage == 0) {
		return Pending, i18n.T("gate.mutation.no_signal")
	}
	if n.SignalStale() {
		return Pending, i18n.T("gate.mutation.stale")
	}
	// Nenhum mutante EXECUTADO: o score é 100 por construção (ver ParseMutation) e não há
	// veredito a dar. Passa em silêncio, e o silêncio é a decisão — os dois motivos
	// possíveis pertencem a outros donos:
	//
	//	tudo ignorado      → não há regra a provar (tabela, tipo, reexport)
	//	tudo sem cobertura → não há TESTE que execute o arquivo, e quem cobra isso é o
	//	                     gate de cobertura de linha, não este
	//
	// Cobrar aqui duplicaria o achado no melhor caso e, no pior — o medido no app de referência —,
	// afogaria os zeros REAIS: 187 arquivos marcavam 0% quase todos por falta de teste, e
	// no meio deles se perdiam os poucos em que o teste roda e não verifica nada.
	if n.Signal.MutantsKilled == 0 && n.Signal.MutantsSurvived == 0 {
		return Pass, ""
	}
	// A régua é do PROJETO e chega junto com a medida: `thresholds` é campo OBRIGATÓRIO
	// do schema Mutation Testing Elements, com `low` e `high` obrigatórios dentro dele.
	// Declará-los também no anchors.yaml criaria duas fontes para o mesmo número, que
	// divergem em silêncio no dia em que uma mudar. O default do engine só entra quando
	// a ferramenta não declarou.
	//
	// Os dois separam perguntas que um limiar só não distingue:
	//
	//	< low        reprova — abaixo do ACEITÁVEL
	//	low … high   passa, e APARECE no relatório: aceitável, ainda não DESEJÁVEL
	//	>= high      passa limpo
	//
	// A faixa do meio não bloqueia nada, de propósito — um aviso que barra é um erro com
	// outro nome. Ela existe para impedir que "aceitável" seja lido como "pronto". Com um
	// limiar só é preciso escolher entre dois erros: posto no desejável, reprova trabalho
	// bom o bastante e vira ruído que se aprende a ignorar; posto no aceitável, a unidade
	// de 71% e a de 95% aparecem iguais, e some a informação de onde ainda há o que ganhar.
	threshold := n.Signal.MutationLow
	if threshold <= 0 {
		threshold = 70.0
	}
	desejavel := n.Signal.MutationHigh

	// Quando os DOIS escopos foram ingeridos, o veredito é sobre o ISOLADO — é ele
	// que responde "o teste desta unidade prova o que ela faz?". O completo entra na
	// mensagem como contexto, e a diferença entre os dois é o achado:
	//
	// Medido em dois projetos reais: um átomo de UI marcou 8% isolado / 77% completo,
	// e uma função pura de negócio marcou 27% / 77%. Nos dois, a maior parte dos
	// mutantes só morre nos DEPENDENTES. Olhando só o completo, os dois pareciam
	// saudáveis; olhando o par, o que se vê é acoplamento — quem prova a unidade são
	// os outros, e um refactor nela não é protegido pelos próprios testes.
	// Um escopo medido contra OUTRA versão do arquivo não entra no veredito. O gate
	// julga pelo ISOLADO, então herdar um isolado velho — renovado de carona quando só
	// o `full` foi reingerido — faria o veredito sair de um número medido sobre código
	// que já mudou. Fora do par, o total continua valendo e a leitura cai no caminho
	// de um escopo só, que é o comportamento honesto para "só um deles é atual".
	iso, temIso := n.Signal.MutationByScope["isolated"]
	full, temFull := n.Signal.MutationByScope["full"]
	if temIso && iso.Stale(n.Rev) {
		temIso = false
	}
	if temFull && full.Stale(n.Rev) {
		temFull = false
	}
	if temIso && temFull {
		delta := full.Score - iso.Score
		ctx := fmt.Sprintf("isolado %.0f%%, completo %.0f%% (delta %.0fp)",
			iso.Score, full.Score, delta)
		if iso.Score < threshold {
			// A leitura do delta sai INTEIRA do comparaDelta. Antes havia uma frase fixa
			// ("Delta alto significa que quem prova esta unidade são os dependentes")
			// concatenada aqui, que aparecia mesmo com delta baixo — e o laudo se
			// contradizia: "os dois escopos concordam … Delta alto significa …".
			return Fail, i18n.T("gate.mutation.isolated_failed",
				iso.Score, threshold, ctx, iso.Survived,
				compareDelta(delta))
		}
		if faixa := middleRange(iso.Score, threshold, desejavel, iso.Survived); faixa != "" {
			return Pending, faixa + " — " + ctx
		}
		return Pass, ""
	}

	if n.Signal.MutationScore < threshold {
		return Fail, i18n.T("gate.mutation.failed",
			n.Signal.MutationScore, threshold, n.Signal.MutantsSurvived)
	}
	if faixa := middleRange(n.Signal.MutationScore, threshold, desejavel, n.Signal.MutantsSurvived); faixa != "" {
		return Pending, faixa
	}
	return Pass, ""
}

// middleRange devolve o laudo de ACEITÁVEL-MAS-NÃO-DESEJÁVEL, ou "" quando não se
// aplica. Devolver Pending (não Fail) é o ponto: a unidade passou, e o que se está
// dizendo é "dá para ir além", não "está errado".
//
// Sem `high` declarado, a faixa não existe e o gate volta ao comportamento de um limiar
// só — nenhum projeto é obrigado a adotar o conceito para continuar usando o gate.
func middleRange(score, minimo, desejavel float64, sobreviventes int) string {
	if desejavel <= 0 || desejavel <= minimo || score >= desejavel {
		return ""
	}
	return i18n.T("gate.mutation.middle_range",
		score, minimo, desejavel, desejavel-score, sobreviventes)
}

// compareDelta traduz a diferença entre os escopos numa frase acionável.
func compareDelta(delta float64) string {
	switch {
	case delta >= 40:
		return i18n.T("gate.mutation.delta_high")
	case delta >= 15:
		return i18n.T("gate.mutation.delta_medium")
	default:
		return i18n.T("gate.mutation.delta_low")
	}
}

// tests-pass: o nó de teste tem 0 falhas (do resultado de execução ingerido)?
func checkTestsPass(_ string, n mapx.Node) (Verdict, string) {
	if n.Signal == nil || (n.Signal.Passed == 0 && n.Signal.Failed == 0 && n.Signal.Skipped == 0) {
		return Pending, i18n.T("gate.tests_pass.no_signal")
	}
	if n.SignalStale() {
		return Pending, i18n.T("gate.tests_pass.stale")
	}
	if n.Signal.Failed > 0 {
		return Fail, i18n.T("gate.tests_pass.failed", n.Signal.Failed, n.Signal.Passed)
	}
	return Pass, ""
}

func checkGuideHasChecklist(content string, _ mapx.Node) (Verdict, string) {
	if !checklistHeadingRE.MatchString(content) {
		return Fail, i18n.T("gate.guide_checklist.missing_section")
	}
	if !checklistItemRE.MatchString(content) {
		return Fail, i18n.T("gate.guide_checklist.missing_items")
	}
	return Pass, ""
}

// isBinary decide se o conteúdo é binário pela presença de byte NUL nos primeiros
// 8 KB — a mesma heurística que o git usa para decidir se mostra o diff. Os gates
// que leem TEXTO (header, spellcheck, asserções) não têm o que fazer com esses
// arquivos, e cobrá-los produz uma exigência impossível de cumprir.
func isBinary(content string) bool {
	n := len(content)
	if n > 8000 {
		n = 8000
	}
	return strings.IndexByte(content[:n], 0) >= 0
}

// placeholderProsaRE: o marcador solto no corpo da spec — `TODO`, `FIXME`, `XXX` ou um
// `<placeholder entre sinais>`, com ou sem o texto que o acompanha.
//
// O vocabulário é o mesmo do gate `placeholder-preenchido`, e é universal: nenhum projeto
// traduz `TODO`. A frase que vem depois dele ("descrever", "describe", "escribir") pode
// mudar com o idioma do template, e por isso não entra na régua.

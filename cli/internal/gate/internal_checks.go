package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
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
	"spec-sections":       checkSpecSections,
	"has-code":            checkHasScenarioCode,
	"guide-has-checklist": checkGuideHasChecklist,
	"scenario-coverage":   checkScenarioCoverage,
	"line-coverage":       checkLineCoverage,
	"coverage-delta":      checkCoverageDelta,
	"mutation-score":      checkMutationScore,
	"tests-pass":          checkTestsPass,
	"header-conforme":     checkHeaderConforme,
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
	"feature-test-match":         checkFeatureTestMatch,
	"cenario-identidade":         checkCenarioIdentidade,
	"cenario-tipo-alinhado":      checkCenarioTipoAlinhado,
	"cenario-letra-declarada":    checkCenarioLetraDeclarada,
	"spec-feature-match":         checkSpecFeatureMatch,
	"code-reference-valid":       checkCodeReferenceValid,
	"scenario-asserts":           checkScenarioAsserts,
	"domain-declared":            checkDomainDeclared,
	"count-honored":              checkCountHonored,
	"trigger-declared":           checkTriggerDeclared,
	"route-exists":               checkRouteExists,
	"placeholder-preenchido":     checkPlaceholderPreenchido,
	"regra-implementada":         checkRegraImplementada,
	"vr-baseline":                checkVRBaseline,
	"ref-resolves":               checkRefResolves,
	"layer-boundary":             checkLayerBoundary,
	"dependency-honored":         checkDependencyHonored,
	"contract-status-declared":   checkContractStatusDeclared,
	"prova-cruza-fronteira":      checkProvaCruzaFronteira,
	"trinca-completa":            checkTrincaCompleta,
	"plan-seeds-valid":           checkPlanSeedsValid,
	"fase-ordenada":              checkFaseOrdenada,
	"fase-existe":                checkFaseExiste,
	"parent-valido":              checkParentValido,
	"plano-revisado":             checkPlanoRevisado,
	"plano-alterado-justificado": checkPlanoAlteradoJustificado,
	"obligation-honored":         checkObligationHonored,
	"sibling-guard":              checkSiblingGuard,
	"pagination-honored":         checkPaginationHonored,
	"open-questions-resolved":    checkOpenQuestions,
	"rule-types":                 checkRuleTypes,
	"identity-consistent":        checkIdentityConsistent,
	"region-pair-honored":        checkRegionPairHonored,
	"evidence-fresh":             checkEvidenceFresh,
	"testid-coerente":            checkTestIDCoerente,
	"testid-consultado-existe":   checkTestIDConsultadoExiste,
	"mock-tipado":                checkMockTipado,
	"mock-carimbado":             checkMockCarimbado,
	"teste-rastreavel":           checkTesteRastreavel,
	"codigo-catalogado":          checkCodigoCatalogado,
}

func runInternal(name string, n mapx.Node, root string, graph *mapx.Graph, cfg *config.Config) (Verdict, string) {
	content, err := os.ReadFile(filepath.Join(root, n.ID))
	if err != nil {
		return Fail, "não foi possível ler o arquivo: " + err.Error()
	}
	if fn, ok := checkersWithGraph[name]; ok {
		return fn(string(content), n, root, graph, cfg)
	}
	if fn, ok := checkersWithRoot[name]; ok {
		return fn(string(content), n, root)
	}
	fn, ok := internalCheckers[name]
	if !ok {
		return Pending, "checker interno desconhecido: " + name
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
}

// runInternalAgregado executa um checker interno de escopo batch/project.
//
// Difere do `runInternal` num ponto que não é detalhe: NÃO há arquivo para ler. O
// escopo é o conjunto, então o checker recebe um nó vazio e se orienta por `root` e
// `cfg`. Tentar ler o conteúdo de um alvo aqui (como o runInternal faz) devolveria
// erro de leitura e o gate reprovaria por um arquivo que nunca existiu.
func runInternalAgregado(g config.Gate, root string, graph *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if fn, ok := checkersComGate[g.Check]; ok {
		return fn(g, root, graph, cfg)
	}
	fn, ok := checkersWithGraph[g.Check]
	if !ok {
		return Pending, "checker interno agregado desconhecido: " + g.Check
	}
	return fn("", mapx.Node{}, root, graph, cfg)
}

// non-empty: o arquivo não é vazio nem só espaço. Trivial, mas pega placeholders.
func checkNonEmpty(content string, _ mapx.Node) (Verdict, string) {
	if strings.TrimSpace(content) == "" {
		return Fail, "arquivo vazio"
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
		return Skip, "" // sem updated_at declarado — o gate header-conforme cobra o header
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
		return Skip, "sem repositório git — o updated_at não tem contra o que ser conferido"
	}
	if mudou {
		today := gitmeta.Today()
		if declared == today {
			return Pass, "" // atualizado para hoje (a alteração em curso)
		}
		return Fail, fmt.Sprintf("updated_at do header (%s) ≠ hoje (%s), e o arquivo tem "+
			"alterações não-commitadas — atualize para o dia da alteração (%s), ou `anchors check --fix`.",
			declared, today, today)
	}

	// COMMITADO: compara com a data do último commit que tocou o arquivo (só o dia).
	gitDate, ok := gitmeta.LastCommitDate(root, n.ID)
	if !ok {
		return Pending, "arquivo sem commit no git (novo/untracked) — nada a comparar"
	}
	if declared != gitDate {
		return Fail, fmt.Sprintf("updated_at do header (%s) ≠ data do último commit (%s). "+
			"Atualize para %s (ou rode `anchors check --fix`).", declared, gitDate, gitDate)
	}
	return Pass, ""
}

// header-conforme: o arquivo tem o BLOCO DE CABEÇALHO do Anchors (`@anchors`) com o
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
// mínima honesta é o `layer:` (a que camada pertence), não um code/ref inventado.
var recognizedLayers = map[string]bool{
	"dao": true, "infra": true, "presentation": true, "domain-types": true,
}

// headerLayerValueRE captura o VALOR do `layer:` declarado no header.
var headerLayerValueRE = regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*layer:\s*(\S+)`)

// isRecognizedLayer decide se o arquivo pertence a uma camada reconhecida — pela tag
// do nó OU pelo `layer:` DECLARADO no header. O header importa porque um TESTE de um
// arquivo de camada reconhecida é classificado como `kind: test` (o pattern de test é
// mais específico que o da camada), perdendo a tag da camada; mas ele declara
// `layer: <reconhecida>` no header, e o header é a fonte da verdade da identidade.
func isRecognizedLayer(n mapx.Node, content string) bool {
	// A Estrutura do projeto é a fonte da verdade: uma camada com `regime: declarativo`
	// é RECONHECIDA (não origina regra), qualquer que seja seu nome — o vocabulário de
	// camadas pertence ao projeto (anchors.yaml), não ao engine. Ver STRUCTURE.md.
	if n.Regime == "declarativo" {
		return true
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

// ehRoteiroExecutavel — o nó é um roteiro do runner e2e (não um arquivo de teste em
// linguagem de programação). Reconhecido pela extensão de dado + kind test: um `.test.tsx`
// é código nosso e carrega cabeçalho; um `.yaml` é entrada de um runner externo.
func ehRoteiroExecutavel(n mapx.Node) bool {
	if n.Kind != mapx.KindTest {
		return false
	}
	return strings.HasSuffix(n.ID, ".yaml") || strings.HasSuffix(n.ID, ".yml")
}

func checkHeaderConforme(content string, n mapx.Node) (Verdict, string) {
	// Arquivo BINÁRIO não carrega cabeçalho — não há sintaxe de comentário num PNG.
	// A identidade dele está no NOME (`<Unidade>.<CODE>-VR-<variante>.png`), que é o
	// que o `identity-consistent` confronta. Cobrar header aqui exigiria o impossível
	// e barraria todo commit de baseline visual.
	if ehBinario(content) {
		return Skip, "arquivo binário — a identidade está no nome, não em cabeçalho"
	}
	// ROTEIRO de teste executável (.yaml do runner e2e): mesma razão do binário acima,
	// por um caminho diferente. A identidade dele está no NOME do arquivo
	// (`<CODE>-A01.yaml`, que é como o `tested-by` o encontra) e nas `tags:`, que o
	// próprio runner consome — e o formato é do RUNNER, não nosso: um bloco `@anchors`
	// no topo de um YAML de flow é comentário morto para quem o executa.
	//
	// A alternativa era retrofitar 717 arquivos com um cabeçalho que nenhum deles jamais
	// teve — e o gate só passou a vê-los porque a camada `e2e-flow` os trouxe para o
	// grafo (para que a execução de E2E pudesse deixar carimbo). Cobrar deles um contrato
	// escrito depois seria transformar a chegada ao mapa em 717 defeitos retroativos.
	if ehRoteiroExecutavel(n) {
		return Skip, "roteiro de teste executável — a identidade está no nome do arquivo e nas tags do runner"
	}
	if !headerBlockRE.MatchString(content) {
		return Fail, "sem bloco de cabeçalho `@anchors` — todo arquivo do grafo deve ter " +
			"o cabeçalho padrão no topo (veja `anchors guide header`)."
	}
	// Camada RECONHECIDA (sem spec): `layer:` é a identidade mínima suficiente.
	if isRecognizedLayer(n, content) {
		if !headerLayerRE.MatchString(content) &&
			!headerCodeRE().MatchString(content) && !headerRefRE().MatchString(content) {
			return Fail, "cabeçalho `@anchors` sem identidade — arquivo de camada reconhecida " +
				"(dao/infra/presentation/domain-types) declara ao menos `layer: <camada>` " +
				"(sua identidade mínima, pois não tem spec). Ver `anchors guide header`."
		}
		return Pass, ""
	}
	// Camada REGIDA (ou spec): exige posse ou referência.
	if !headerCodeRE().MatchString(content) && !headerRefRE().MatchString(content) {
		return Fail, "cabeçalho `@anchors` sem identidade — declare `code: <CÓDIGO>` se este " +
			"arquivo é o DONO (a spec), ou `ref: <CÓDIGO>` se ele referencia a unidade (o " +
			"resto da trinca: code/feature/test). Ver `anchors guide header`."
	}
	return Pass, ""
}

func checkSpecSections(content string, _ mapx.Node) (Verdict, string) {
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
		return Fail, "spec sem nenhuma regra/estado CATALOGADO. Cada regra precisa de um " +
			"código e um lugar estruturado — cabeçalho `### ABCDX-V01`, linha de tabela " +
			"`| ABCDX-V01 | ... |`, ou bullet-negrito `- **ABCDX-V01** ...`. (Menção solta " +
			"em prosa não conta.)"
	}
	// "Ao menos uma" é o piso, e sozinho ele deixa passar o caso mais comum: a spec que
	// cataloga a primeira regra e escreve as outras nove em prosa. As nove ficam
	// invisíveis para TODOS os gates de identidade — não por estarem erradas, mas por não
	// terem código, e um gate que não enxerga a regra reporta verde sobre o que não
	// conferiu.
	if msg := irmasSemCodigo(content); msg != "" {
		return Fail, msg
	}
	return Pass, ""
}

// seçãoVazia diz se a seção não tem CONTEÚDO — só o esqueleto.
//
// Conta como vazio: linha em branco, separador de tabela (`|---|---|`), cabeçalho de
// tabela (a linha de rótulos que precede o separador) e o divisor `---`. Qualquer outra
// coisa é conteúdo, e aí a ausência de código volta a ser o defeito que o gate pega.
func seçãoVazia(linhas []string) bool {
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

var sepTabelaRE = regexp.MustCompile(`^\|[\s:|-]+\|?$`)

// secaoComNivelRE casa um cabeçalho e captura o nível (para agrupar por pai) e o título.
var secaoComNivelRE = regexp.MustCompile(`^(#{2,4})\s+(.+?)\s*$`)

// irmasSemCodigo acha a seção que DEVERIA catalogar e não cataloga.
//
// O difícil é separar "seção de regra sem código" de "seção de prosa" — `## Visão Geral`
// e `## Restrições` não catalogam nada e não devem ser cobradas. A pergunta não se
// responde pelo título (cada projeto nomeia como quer), e sim pela VIZINHANÇA: entre as
// seções irmãs sob o mesmo pai, se algumas catalogam com código e outras não, as sem
// código destoam do padrão que a própria spec estabeleceu.
//
// É o mesmo princípio do `unclaimedSections` no `rule-types`: a spec declara sua forma, e
// o gate cobra coerência com ela — não com um formato que o Anchors imponha.
func irmasSemCodigo(content string) string {
	type sec struct {
		titulo string
		nivel  int
		linhas []string
	}
	var secoes []sec
	for _, l := range strings.Split(content, "\n") {
		if m := secaoComNivelRE.FindStringSubmatch(l); m != nil {
			secoes = append(secoes, sec{titulo: m[2], nivel: len(m[1])})
			continue
		}
		if len(secoes) > 0 {
			secoes[len(secoes)-1].linhas = append(secoes[len(secoes)-1].linhas, l)
		}
	}
	// Agrupa por PAI: o `##` que precede cada bloco de `###`.
	pai := ""
	grupos := map[string][]sec{}
	for _, s := range secoes {
		if s.nivel == 2 {
			pai = s.titulo
			continue
		}
		grupos[pai] = append(grupos[pai], s)
	}
	for _, irmas := range grupos {
		if len(irmas) < 2 {
			continue // sem irmã não há padrão a comparar
		}
		var comCodigo, semCodigo []string
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
			// São dois estados opostos que a busca por código não distingue:
			//   • a seção com regras escritas em prosa   → o defeito que este gate pega
			//   • a seção só com o esqueleto da tabela   → nada foi escrito ainda
			//
			// O caso real: seis telas de onboarding estáticas herdaram o cabeçalho de
			// "Comportamentos Automáticos" do modelo de spec, e não têm comportamento
			// automático nenhum (zero efeitos no código). O gate as barrava pedindo
			// códigos para uma tabela sem linhas.
			//
			// Quem cobra o esqueleto vazio é `placeholder-preenchido`, que sabe
			// distinguir o marcador real do exemplo de sintaxe. Dois gates sobre o mesmo
			// defeito produzem dois achados para um problema.
			if seçãoVazia(s.linhas) {
				continue
			}
			semCodigo = append(semCodigo, s.titulo)
		}
		// Só acusa quando a MAIORIA das irmãs cataloga: um grupo em que só uma tem código
		// não estabeleceu padrão nenhum, e cobrar as outras seria inventar um.
		if len(comCodigo) > len(semCodigo) && len(semCodigo) > 0 {
			return fmt.Sprintf("seção(ões) sem código sob irmãs que catalogam: %s. "+
				"As irmãs (%s) têm código e estas não — uma regra sem código é invisível "+
				"para os gates de identidade, e eles reportam verde sobre o que não "+
				"conferiram. Dê um código a cada uma, ou mova o texto para uma seção de "+
				"prosa (`## Restrições`, `## Notas`) se não for regra.",
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
	return Fail, "sem código de cenário (identidade ausente)"
}

// guide-has-checklist: um GUIDE de governança deve destilar suas regras em PONTOS DE
// CONFORMIDADE verificáveis (a seção "## Pontos de conformidade" com itens CK1, CK2…).
// Sem eles, o gate de julgamento por IA recai em heurística vaga. Este checker é
// DETERMINÍSTICO — só verifica a PRESENÇA da seção e de ao menos um item CK; a
// QUALIDADE de cada ponto é julgamento (ver `anchors guide guide`). Ver a sugestão do
// meta-guide: "a presença da checklist é um check determinístico".
var checklistHeadingRE = regexp.MustCompile(`(?mi)^##+\s+pontos de conformidade\b`)
var checklistItemRE = regexp.MustCompile(`(?m)\bCK\d+\b`)

// --- sinais de teste ingeridos (execução, cobertura) — QUALITY §5 Execução ---

// scenario-coverage: cada código de cenário que a spec DECLARA tem um teste que
// PASSOU (está em Signal.ProvenCodes)? Fecha o gate de Rastreabilidade — não basta
// existir teste, cada requisito precisa estar provado. Pending se nada foi ingerido.
func checkScenarioCoverage(content string, n mapx.Node) (Verdict, string) {
	declared := anyCodeRE.FindAllString(content, -1)
	if len(declared) == 0 {
		return Skip, "" // spec sem cenários — nada a cobrir
	}
	if n.Signal == nil {
		return Pending, "sem sinal de teste ingerido (rode `anchors ingest --junit`)"
	}
	if n.SignalStale() {
		return Pending, "sinal de teste stale (a spec mudou desde a ingestão — reingira)"
	}
	proven := map[string]bool{}
	for _, c := range n.Signal.ProvenCodes {
		proven[c] = true
	}
	var missing []string
	seen := map[string]bool{}
	for _, code := range declared {
		if seen[code] {
			continue
		}
		seen[code] = true
		if !proven[code] {
			missing = append(missing, code)
		}
	}
	if len(missing) > 0 {
		return Fail, fmt.Sprintf("%d cenário(s) sem teste verde: %s", len(missing), strings.Join(missing, ", "))
	}
	return Pass, ""
}

// line-coverage: a cobertura de linha do nó de código está >= 70%? (limiar fixo por
// ora — poderia vir da config). Pending se não ingerida.
func checkLineCoverage(_ string, n mapx.Node) (Verdict, string) {
	if n.Signal == nil || n.Signal.TotalLines == 0 {
		return Pending, "sem cobertura de linha ingerida (rode `anchors ingest --lcov`)"
	}
	if n.SignalStale() {
		return Pending, "cobertura stale (o arquivo mudou desde a ingestão — reingira)"
	}
	const threshold = 70.0
	if n.Signal.LineCoverage < threshold {
		return Fail, fmt.Sprintf("cobertura de linha %.0f%% < %.0f%%", n.Signal.LineCoverage, threshold)
	}
	return Pass, ""
}

// coverage-delta: a cobertura de linha deste nó de código CAIU desde a ingestão
// anterior? Pega a regressão de cobertura (código pode estar coberto no diff mas ter
// derrubado a cobertura de outra parte). Pending se não há baseline.
func checkCoverageDelta(_ string, n mapx.Node) (Verdict, string) {
	if n.Signal == nil || n.Signal.PrevLineCoverage == 0 {
		return Pending, "sem baseline de cobertura (ingira duas vezes para comparar)"
	}
	if n.SignalStale() {
		return Pending, "cobertura stale (o arquivo mudou desde a ingestão — reingira)"
	}
	delta := n.Signal.LineCoverage - n.Signal.PrevLineCoverage
	if delta < -0.01 {
		return Fail, fmt.Sprintf("cobertura caiu %.0f%% → %.0f%% (%.0f pontos)", n.Signal.PrevLineCoverage, n.Signal.LineCoverage, delta)
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
		return Pending, "sem sinal de mutação ingerido — a cobertura de linha diz que a " +
			"linha EXECUTOU, não que alguém verificou o resultado. Rode a ferramenta de " +
			"mutação do seu stack (Stryker/PIT/Infection/mutmut) e " +
			"`anchors ingest --mutation <relatório>`"
	}
	if n.SignalStale() {
		return Pending, "sinal de mutação stale (o arquivo mudou desde a ingestão — reingira)"
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
			return Fail, fmt.Sprintf("score de mutação ISOLADO %.0f%% < %.0f%% — %s. "+
				"%d mutante(s) sobrevivem ao teste da própria unidade; %s",
				iso.Score, threshold, ctx, iso.Survived,
				comparaDelta(delta))
		}
		if faixa := faixaDoMeio(iso.Score, threshold, desejavel, iso.Survived); faixa != "" {
			return Pending, faixa + " — " + ctx
		}
		return Pass, ""
	}

	if n.Signal.MutationScore < threshold {
		return Fail, fmt.Sprintf("score de mutação %.0f%% < %.0f%% — %d mutante(s) sobreviveram: "+
			"alterações no código que os testes não perceberam",
			n.Signal.MutationScore, threshold, n.Signal.MutantsSurvived)
	}
	if faixa := faixaDoMeio(n.Signal.MutationScore, threshold, desejavel, n.Signal.MutantsSurvived); faixa != "" {
		return Pending, faixa
	}
	return Pass, ""
}

// faixaDoMeio devolve o laudo de ACEITÁVEL-MAS-NÃO-DESEJÁVEL, ou "" quando não se
// aplica. Devolver Pending (não Fail) é o ponto: a unidade passou, e o que se está
// dizendo é "dá para ir além", não "está errado".
//
// Sem `high` declarado, a faixa não existe e o gate volta ao comportamento de um limiar
// só — nenhum projeto é obrigado a adotar o conceito para continuar usando o gate.
func faixaDoMeio(score, minimo, desejavel float64, sobreviventes int) string {
	if desejavel <= 0 || desejavel <= minimo || score >= desejavel {
		return ""
	}
	return fmt.Sprintf("score de mutação %.0f%%: aceitável (>= %.0f%%), ainda não desejável (%.0f%%) — "+
		"faltam %.0f ponto(s), %d mutante(s) ainda sobrevivem",
		score, minimo, desejavel, desejavel-score, sobreviventes)
}

// comparaDelta traduz a diferença entre os escopos numa frase acionável.
func comparaDelta(delta float64) string {
	switch {
	case delta >= 40:
		return "a suíte completa cobre muito mais, então o buraco é no teste da unidade: " +
			"quem prova esta unidade são os dependentes, e um refactor aqui não é protegido " +
			"pelos testes daqui"
	case delta >= 15:
		return "parte da prova vem dos dependentes — nessa fatia, um refactor aqui não é " +
			"protegido pelos testes daqui"
	default:
		return "os dois escopos concordam — o buraco não é de acoplamento, é de asserção"
	}
}

// tests-pass: o nó de teste tem 0 falhas (do resultado de execução ingerido)?
func checkTestsPass(_ string, n mapx.Node) (Verdict, string) {
	if n.Signal == nil || (n.Signal.Passed == 0 && n.Signal.Failed == 0 && n.Signal.Skipped == 0) {
		return Pending, "sem resultado de execução ingerido (rode `anchors ingest --junit`)"
	}
	if n.SignalStale() {
		return Pending, "resultado stale (o teste mudou desde a ingestão — reingira)"
	}
	if n.Signal.Failed > 0 {
		return Fail, fmt.Sprintf("%d teste(s) falhando (%d passam)", n.Signal.Failed, n.Signal.Passed)
	}
	return Pass, ""
}

func checkGuideHasChecklist(content string, _ mapx.Node) (Verdict, string) {
	if !checklistHeadingRE.MatchString(content) {
		return Fail, "guide sem a seção '## Pontos de conformidade' — sem ela o julgamento " +
			"vira heurística. Destile as regras em itens CK verificáveis (veja `anchors guide guide`)."
	}
	if !checklistItemRE.MatchString(content) {
		return Fail, "a seção de pontos de conformidade existe mas não tem itens 'CKn' — " +
			"liste cada ponto verificável com seu código (CK1, CK2, …)."
	}
	return Pass, ""
}

// ehBinario decide se o conteúdo é binário pela presença de byte NUL nos primeiros
// 8 KB — a mesma heurística que o git usa para decidir se mostra o diff. Os gates
// que leem TEXTO (header, spellcheck, asserções) não têm o que fazer com esses
// arquivos, e cobrá-los produz uma exigência impossível de cumprir.
func ehBinario(content string) bool {
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

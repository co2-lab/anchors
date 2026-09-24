package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// trinca-completa: uma spec de camada REGIDA precisa das três peças que a realizam —
// o CÓDIGO (`specifies`), a FEATURE (`covered-by`) e o TESTE (`tested-by`).
//
// Existe porque os gates relacionais FALHAM ABERTO por construção: sem teste ligado, o
// feature-test-match devolve Pending ("nada a confrontar ainda") em vez de reprovar, e
// o `dependency-honored` idem sem código. O efeito colateral é grave — uma spec sozinha,
// sem nenhuma implementação, atravessa TODOS os gates e o pipeline conclui
// "✓ pode promover". Ou seja: o verde certifica trabalho que não existe.
//
// Este gate fecha esse buraco pelo lado positivo: em vez de perguntar "as peças casam?"
// (o que exige que elas existam), pergunta "as peças existem?".
//
// NÃO se aplica a:
//   - camadas RECONHECIDAS (regime declarativo — dao/infra/resource): não têm spec nem
//     trinca por definição;
//   - specs cuja camada dispensa alguma peça por de-para do projeto (ex.: repository, que
//     no app de referência é provado por teste de integração central, não por teste co-localizado).
//     Isso é declarado com `trinca_opcional` na camada do anchors.yaml.
func checkTriadComplete(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.triad.skip_not_spec")
	}
	if g == nil {
		return pendingNoMap()
	}
	// Camada RECONHECIDA não tem trinca a cobrar.
	if isRecognizedLayerCfg(n, content, cfg) {
		return Skip, i18n.T("gate.triad.skip_recognized")
	}

	// A trinca NÃO é uma estrela: a spec aponta o código (`specifies`) e a feature
	// (`covered-by`), mas quem aponta o teste é a FEATURE (`tested-by`, ver mapx.Build).
	// Então o teste se alcança em DOIS saltos — spec → feature → teste. Verificar
	// `tested-by` direto na spec acusaria falta de teste em todo o projeto.
	have := map[mapx.EdgeType]bool{}
	for _, e := range g.Neighbors(n.ID).Out {
		have[e.Type] = true
		if e.Type == mapx.EdgeCoveredBy {
			for _, fe := range g.Neighbors(e.To).Out {
				if fe.Type == mapx.EdgeTestedBy {
					have[mapx.EdgeTestedBy] = true
				}
			}
		}
	}

	// Peças exigidas, na ordem em que o dev as produz.
	required := []struct {
		edge  mapx.EdgeType
		peca  string
		ondeE string
	}{
		{mapx.EdgeSpecifies, i18n.T("gate.triad.code_piece"), i18n.T("gate.triad.code_desc")},
		{mapx.EdgeCoveredBy, i18n.T("gate.triad.feature_piece"), i18n.T("gate.triad.feature_desc")},
		{mapx.EdgeTestedBy, i18n.T("gate.triad.test_piece"), i18n.T("gate.triad.test_desc")},
	}
	optional := optionalPieces(n, cfg, g)
	// Dispensa POR UNIDADE, declarada na própria spec (`@no-test:`/`@no-feature:`).
	//
	// A dispensa por CAMADA (`trinca_opcional`) isenta em bloco: ou toda `service`
	// precisa de teste, ou nenhuma. Mas dentro da mesma camada convivem o gateway de
	// 9 linhas que só repassa a chamada e o módulo de 150 com regra de verdade —
	// isentar os dois junto apaga a cobrança justamente onde ela vale.
	//
	// Aqui a decisão é da UNIDADE e fica escrita nela, com razão obrigatória: quem lê
	// a spec vê por que aquele arquivo não tem teste, em vez de descobrir num
	// `trinca_opcional` distante que removeu a exigência da camada inteira.
	waived := specWaivers(content)
	for peca := range waived {
		optional[peca] = true
	}
	// DEBT (`@TBD`) is counted APART from the waivers, and the distinction decides the
	// verdict: `@no-<piece>: <reason>` asserts "this unit WILL NEVER HAVE that piece",
	// and the gate goes quiet for good; `@TBD: code,test` asserts "I have not written it
	// yet", which is another thing entirely — pending work, not a decision taken.
	//
	// Until now the two shared a bucket and both became Pass. The effect was to erase
	// from the radar exactly what remains to be done: a spec with `@TBD: code,feature,test`
	// came out green, indistinguishable from a complete triad.
	//
	// Now the deferred piece does not fail (it was declared, with a written reason) but
	// does not pass either: it becomes Pending, which keeps showing up until someone
	// pays it.
	toDevelop := piecesToDevelop(content)

	// A dispensa CONTRADIZ a feature? `@no-test` diz "esta unidade não precisa de
	// teste"; um cenário na feature diz "este comportamento se verifica assim". As
	// duas afirmações não podem conviver: ou o cenário é real e alguém precisa
	// prová-lo, ou ele não deveria existir.
	//
	// Sem esta checagem a contradição fica MUDA — o `trinca-completa` passa (a
	// dispensa o satisfaz) e o `feature-test-match` também (ele exige teste só
	// quando há teste a confrontar). O resultado é um cenário escrito que ninguém
	// prova, com o pipeline inteiro verde.
	// `@no-test` sozinho tem de dizer ONDE a prova está, e o lugar tem de existir.
	//
	// Só se cobra quando a dispensa vem do `@no-test`: se a spec declara `@no-feature`,
	// não há comportamento observável nenhum e não existe prova a apontar — exigir
	// referência ali seria pedir o endereço de algo que a spec acabou de dizer que não
	// existe.
	if unitWaiver(noTestRE, content) && !unitWaiver(noFeatureRE, content) {
		alvo, temRef := proofPointedByNoTest(content)
		if !temRef {
			return Fail, i18n.T("gate.triad.no_test_no_pointer")
		}
		if _, achou := provingTest(alvo, root, g); !achou {
			return Fail, i18n.T("gate.triad.no_test_unresolved", alvo)
		}
	}

	if waived[string(mapx.EdgeTestedBy)] {
		if qtd, feat := linkedFeatureScenarios(n, root, g); qtd > 0 {
			return Fail, i18n.T("gate.triad.no_test_conflict", feat, qtd)
		}
	}

	var missing []string
	var owed []string
	for _, r := range required {
		if have[r.edge] || optional[string(r.edge)] {
			continue
		}
		if toDevelop[string(r.edge)] {
			owed = append(owed, r.peca)
			continue
		}
		missing = append(missing, r.peca+" ("+r.ondeE+")")
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Fail, i18n.T("gate.triad.incomplete", strings.Join(missing, "; "))
	}
	if len(owed) > 0 {
		sort.Strings(owed)
		return Pending, i18n.T("gate.triad.to_be_developed", len(owed), strings.Join(owed, ", "))
	}
	return Pass, ""
}

// optionalPieces lê da camada do nó quais peças da trinca o projeto dispensa
// (`trinca_opcional: [tested-by]`, por exemplo). É o opt-out HONESTO: fica declarado na
// Estrutura, não escondido num Skip do gate.
func optionalPieces(n mapx.Node, cfg *config.Config, g *mapx.Graph) map[string]bool {
	out := map[string]bool{}
	if cfg == nil {
		return out
	}
	// A dispensa é declarada na camada do CÓDIGO (`schema-model`, `dao`…), mas este
	// gate roda sobre a SPEC — e uma `.spec.md` casa a layer `spec`, nunca a do alvo.
	// Então resolvemos a camada pelo ALVO que a spec descreve (aresta `specifies`).
	tags := append([]string{}, n.Tags...)
	if g != nil {
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type == mapx.EdgeSpecifies {
				for _, alvo := range g.Nodes {
					if alvo.ID == e.To {
						tags = append(tags, alvo.Tags...)
					}
				}
			}
		}
	}
	for _, t := range tags {
		l, ok := cfg.Layers[t]
		if !ok {
			continue
		}
		for _, p := range l.OptionalTriadEdges {
			out[p] = true
		}
	}
	// A aresta `specifies` só existe DEPOIS que o código nasce — e na etapa `spec` ele
	// ainda não nasceu. Sem esta segunda via, o gate cobrava justamente a peça que a
	// camada dispensa, no único momento em que o autor não tem como provar o contrário:
	// o `work` dizia "DISPENSADA: NÃO crie" e o gate respondia "falta a feature".
	//
	// Resolver pelo CAMINHO do alvo cobre esse intervalo. É o mesmo alvo que a aresta
	// apontaria — só que deduzido do nome, que existe desde o primeiro instante.
	if len(out) == 0 {
		if base := strings.TrimSuffix(n.ID, ".spec.md"); base != n.ID {
			// A extensão certa é a que resolve para uma camada COM dispensa declarada —
			// tentar todas evita chutar `.ts` numa tela (`.tsx`) e cair no catch-all.
			for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
				layer, _ := scan.Classify(base+ext, cfg)
				if layer == "" {
					continue
				}
				if disp := cfg.Layers[layer].OptionalTriadEdges; len(disp) > 0 {
					for _, p := range disp {
						out[p] = true
					}
					break
				}
			}
		}
	}
	return out
}

// ── Dispensa por unidade ─────────────────────────────────────────────────────
//
// `@no-test: <razão>` e `@no-feature: <razão>` no cabeçalho da spec dispensam
// aquela peça DAQUELA unidade. Mesmo padrão do `@no-scenario`/`@no-code`/
// `@no-paginate` (CONCEPT §5.1): a razão é OBRIGATÓRIA, e é ela que separa decisão
// de esquecimento — um marcador nu passaria a ser um jeito silencioso de calar o
// gate, que é exatamente o que a dispensa não pode virar.
//
// O `@no-test` ainda exige uma REFERÊNCIA VERIFICÁVEL — ver `provaApontadaPeloNoTest`.
var (
	noTestRE    = regexp.MustCompile(`@no-test[^\S\n]*:[^\S\n]*\S+`)
	noFeatureRE = regexp.MustCompile(`@no-feature[^\S\n]*:[^\S\n]*\S+`)
	// `@no-code: <razão>` — a unidade não tem MÓDULO a escrever.
	//
	// O comentário acima citava este marcador desde sempre e ele nunca existiu. A lacuna
	// só apareceu quando o gate passou a confrontar de verdade: uma unidade cuja
	// implementação É configuração — os workflows do CI, no projeto de referência — não
	// tem arquivo de código, e a única saída era deixar o gate reprovando para sempre ou
	// fingir um módulo que "valida a configuração" e só saberia dizer que ela existe.
	//
	// A razão é obrigatória, como nas outras duas: um marcador nu seria um jeito
	// silencioso de calar o gate.
	noCodeRE = regexp.MustCompile(`@no-code[^\S\n]*:[^\S\n]*\S+`)

	// `@TBD` — TO BE DEVELOPED: a peça está decidida e AINDA NÃO foi escrita.
	//
	// A diferença para `@no-test` é o TEMPO, e ela é o motivo de o marcador existir.
	// `@no-test` afirma "esta unidade NÃO PRECISA de teste" — decisão permanente, que
	// fica escrita para quem ler a spec depois. `@TBD` afirma "falta escrever", e vence
	// sozinho: no instante em que a peça aparece no mapa, a dispensa deixa de valer sem
	// ninguém remover nada.
	//
	// Também não é "a decidir": o que a peça vai fazer já está na spec — é justamente
	// isso que a spec É. O que falta é o arquivo.
	//
	// Sem esta distinção, a spec que nasce ANTES do código — que é o fluxo normal do
	// Anchors, já que a spec é a âncora — tinha duas saídas, ambas ruins: barrar o
	// commit de todo trabalho em andamento, ou declarar `@no-test` mentindo, e aí a
	// cobrança some para sempre justamente na unidade que mais vai precisar dela.
	//
	// O ALVO é obrigatório (`@TBD: code`, `@TBD: code,test`): "está em andamento" sem
	// dizer o quê viraria um interruptor geral do gate.
	// FORA DE CRASES: uma revisão que EXPLICA a remoção da dispensa cita o marcador — "a
	// dispensa `@TBD: code,feature,test` do cabeçalho SAIU" — e sem esta guarda a citação
	// a reativa. A spec passaria a declarar uma ausência que o texto ao lado diz ter
	// deixado de existir.
	//
	// A crase é o sinal certo porque é como a doutrina cita qualquer marcador: em código,
	// dentro do parágrafo. Um marcador ATIVO nunca está entre crases — ele é a declaração,
	// não a menção a ela.
	tbdRE = regexp.MustCompile(`(?i)(^|[^` + "`" + `])@TBD[^\S\n]*:[^\S\n]*([a-z,\s]+)`)
)

// referenciaRE — o CÓDIGO do cenário que prova esta unidade, entre crases, dentro do
// bloco do `@no-test`: `@no-test: provado por `SGHBX-B01` no teste do handler`.
//
// A referência é o código, não o caminho do arquivo, por duas razões:
//
//  1. O código sobrevive a mover e renomear arquivo; um caminho apodrece no primeiro
//     refactor, e uma referência podre é PIOR que nenhuma — ela passa a impressão de
//     que a prova foi conferida.
//  2. O código permite conferir que o teste cobre ESTE comportamento, não apenas que
//     algum arquivo existe. É a diferença entre "há um teste lá" e "há um teste DISTO".
//
// A gramática do código NÃO é escrita aqui: vem de `scan.ScenarioCodeRE()`, que honra as
// `rule_types` do projeto. Um regex próprio com as letras canônicas fixas rejeitaria o
// código de uma letra que o projeto declarou (ex.: `-I01`, de Invariant) — o Anchors
// passaria a exigir uma referência que ele mesmo se recusa a reconhecer.
func referenceInBlock(bloco string) (string, bool) {
	for _, m := range crasesRE.FindAllStringSubmatch(bloco, -1) {
		if code := scan.ScenarioCodeRE().FindString(m[1]); code != "" {
			return code, true
		}
	}
	return "", false
}

// crasesRE isola o conteúdo entre crases; a validação de que aquilo é um código fica
// com a gramática do projeto, não com este regex.
var crasesRE = regexp.MustCompile("`([^`]+)`")

// blocoNoTestRE isola o parágrafo do `@no-test` — a referência tem de estar NELE, não
// em qualquer lugar da spec. Sem isso, uma crase solta na Visão Geral satisfaria a
// exigência e o gate voltaria a aceitar prosa.
var blocoNoTestRE = regexp.MustCompile(`(?ms)^@no-test:.*?(?:\n\n|\n@|\z)`)

// proofPointedByNoTest extrai o caminho que a dispensa alega conter a prova.
//
// `@no-test` afirma algo mais forte que as outras dispensas: NÃO que o comportamento
// seja inobservável (isso é `@no-feature`), mas que a prova dele existe EM OUTRO LUGAR —
// no teste de integração central, no flow e2e, no teste do handler que consome esta
// unidade. Enquanto esse "outro lugar" for prosa, a afirmação não é conferível: dá para
// escrever "provado na integração" sem que exista integração alguma, e o gate aceita.
// A razão obrigatória garante que houve DECISÃO, não que a decisão seja VERDADEIRA.
//
// LIMITE DESTE GATE, explícito: ele confere que o arquivo apontado EXISTE. Não confere
// que aquele teste prova ESTE comportamento — isso é julgamento, terreno do
// `feature-test-match`. É a diferença entre "a referência não é fantasia" e "a prova
// está correta"; só a primeira é determinística, e é só ela que este gate promete.
func proofPointedByNoTest(content string) (codigo string, temReferencia bool) {
	bloco := blocoNoTestRE.FindString(content)
	if bloco == "" {
		return "", false
	}
	return referenceInBlock(bloco)
}

// provingTest procura, entre os TESTES do mapa, algum que mencione o código alegado.
//
// Determinístico de ponta a ponta: o conjunto de testes vem do mapa (não de um glob que
// eu chutaria aqui) e a busca é textual pelo código. O que ele responde é "existe um
// teste que se declara prova DESTE código?" — e essa pergunta tem resposta binária.
//
// A busca é TEXTUAL de propósito, e isso é uma decisão de agnosticismo, não preguiça:
// casar `it('CODE: …')` amarraria o Anchors ao dialeto de um framework (jest/vitest),
// e um projeto Go, Python ou Rust — que nomeia caso como `func TestX` ou `def test_x`,
// ou cita o código num atributo — deixaria de conseguir usar a dispensa. O framework
// não pode presumir a sintaxe de teste do projeto que ele governa.
//
// O preço é uma folga: um código citado só em COMENTÁRIO resolve a referência sem
// provar nada. Fechá-la aqui exigiria justamente o parser por dialeto que acabamos de
// recusar — por isso ela é coberta pelo gate de JULGAMENTO `no-test-prova-real`, que
// pergunta se o teste exercita o comportamento em vez de só mencioná-lo. Determinístico
// prova que a referência resolve; julgamento prova que ela vale.
func provingTest(codigo, root string, g *mapx.Graph) (arquivo string, achou bool) {
	if g == nil {
		return "", false
	}
	for _, node := range g.Nodes {
		if node.Kind != mapx.KindTest {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, node.ID))
		if err != nil {
			continue
		}
		if strings.Contains(string(b), codigo) {
			return node.ID, true
		}
	}
	return "", false
}

// specWaivers devolve as ARESTAS dispensadas pela própria spec.
// unitWaiver reports whether a UNIT-level waiver is declared: the marker outside a table
// row. A line that starts with `|` is a rule's row, and a waiver written there is that
// RULE's (the per-rule shape `rule-implemented` reads) — not the unit's. Matching the whole
// spec turned every per-rule `@no-code:` into a waiver of code, feature and test for the
// unit: measured in MIF, 293 of 308 triad failures were that one misreading.
func unitWaiver(re *regexp.Regexp, content string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			continue
		}
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

func specWaivers(content string) map[string]bool {
	out := map[string]bool{}
	if unitWaiver(noTestRE, content) {
		out[string(mapx.EdgeTestedBy)] = true
	}
	if unitWaiver(noFeatureRE, content) {
		// Sem feature não há o que provar por cenário — a dispensa da feature
		// arrasta a do teste, senão o gate cobraria um teste de cenário nenhum.
		out[string(mapx.EdgeCoveredBy)] = true
		out[string(mapx.EdgeTestedBy)] = true
	}
	if unitWaiver(noCodeRE, content) {
		// Sem módulo, não há o que a feature exercitar nem o que o teste provar: a
		// dispensa do código arrasta as outras duas, pelo mesmo motivo que a do
		// `@no-feature` arrasta a do teste.
		//
		// Quem quiser dispensar SÓ o código — uma unidade cuja implementação é
		// configuração mas que ainda assim tem comportamento observável — declara as três
		// separadamente, e a razão de cada uma fica escrita.
		out[string(mapx.EdgeSpecifies)] = true
		out[string(mapx.EdgeCoveredBy)] = true
		out[string(mapx.EdgeTestedBy)] = true
	}
	return out
}

// linkedFeatureScenarios conta os cenários da feature que a spec cobre, e devolve
// o caminho dela. Zero quando não há feature ligada — aí não existe contradição.
func linkedFeatureScenarios(n mapx.Node, root string, g *mapx.Graph) (int, string) {
	if g == nil {
		return 0, ""
	}
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeCoveredBy {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.To))
		if err != nil {
			continue
		}
		// A contagem é dos CENÁRIOS, não das tags: uma feature pode existir com
		// cabeçalho e nenhum cenário (esqueleto), e isso não contradiz nada.
		return len(scenarioRE.FindAllString(string(b), -1)), e.To
	}
	return 0, ""
}

// scenarioRE matches the line that OPENS a scenario, in any language of the official
// Gherkin table (`config.GherkinScenarioAlternatives`).
//
// It used to carry five keywords nailed in Portuguese and English. That was a partial
// copy of the table the dialect already keeps for exactly this purpose, and it went blind
// on valid features: measured when `non-empty` began counting scenarios, a Spanish feature
// (`Escenario:`) and an English one written as `Rule:` + `Example:` both failed a
// BLOCKING gate as "feature with no scenario".
var scenarioRE = regexp.MustCompile(`(?m)^\s*(?:` +
	strings.Join(quoteAll(config.GherkinScenarioAlternatives()), "|") + `):`)

// piecesToDevelop lê o `@TBD:` e devolve as arestas cuja peça ainda não foi escrita.
//
// O vocabulário é o do TRABALHO (`code`, `feature`, `test`), e não o das arestas do mapa:
// quem escreve a spec pensa em peças, não em `covered_by`.
func piecesToDevelop(content string) map[string]bool {
	out := map[string]bool{}
	m := tbdRE.FindStringSubmatch(content)
	if m == nil {
		return out
	}
	// `m[2]`: o grupo 1 é o caractere antes do `@TBD`, que a guarda de crase introduziu.
	for _, peca := range strings.Split(m[2], ",") {
		switch strings.TrimSpace(peca) {
		case "code", "codigo", "código":
			out[string(mapx.EdgeSpecifies)] = true
		case "feature":
			out[string(mapx.EdgeCoveredBy)] = true
		case "test", "teste":
			out[string(mapx.EdgeTestedBy)] = true
		}
	}
	return out
}

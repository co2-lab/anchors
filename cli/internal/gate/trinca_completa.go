package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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
func checkTrincaCompleta(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "não é uma spec — a trinca é cobrada da spec (a dona do código)"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}
	// Camada RECONHECIDA não tem trinca a cobrar.
	if isRecognizedLayer(n, content) {
		return Skip, "camada reconhecida (declarativa) — não tem trinca"
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
		{mapx.EdgeSpecifies, "código", "o arquivo que a spec descreve"},
		{mapx.EdgeCoveredBy, "feature", "os cenários em Gherkin (`.feature`)"},
		{mapx.EdgeTestedBy, "teste", "o teste que prova os cenários"},
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
	dispensas := dispensasDaSpec(content)
	for peca := range dispensas {
		optional[peca] = true
	}

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
	if noTestRE.MatchString(content) && !noFeatureRE.MatchString(content) {
		alvo, temRef := provaApontadaPeloNoTest(content)
		if !temRef {
			return Fail, "a spec declara `@no-test` mas não aponta QUAL cenário prova esta " +
				"unidade. A dispensa afirma que o comportamento é provado em outro lugar " +
				"(integração central, flow e2e, teste do handler que consome esta unidade) — " +
				"e enquanto esse lugar for prosa, a afirmação não é conferível: dá para " +
				"escrever \"provado na integração\" sem que exista integração alguma.\n\n" +
				"Cite o CÓDIGO do cenário entre crases no bloco do `@no-test`, ex.: " +
				"`@no-test: repassa ao handler, provado por `SGHBX-B01``"
		}
		if _, achou := testeQueProva(alvo, root, g); !achou {
			return Fail, fmt.Sprintf(
				"o `@no-test` aponta a prova em `%s`, mas nenhum teste do projeto menciona "+
					"esse código. Uma referência que não resolve é pior que nenhuma: ela passa "+
					"a impressão de que a prova foi conferida.\n\nCorrija o código, ou escreva "+
					"o teste que falta", alvo)
		}
	}

	if dispensas[string(mapx.EdgeTestedBy)] {
		if qtd, feat := cenariosDaFeatureLigada(n, root, g); qtd > 0 {
			return Fail, fmt.Sprintf(
				"a spec declara `@no-test` mas a feature ligada (`%s`) tem %d cenário(s). "+
					"São afirmações contraditórias: a dispensa diz que não há o que provar, e o "+
					"cenário diz o contrário — e a contradição fica MUDA, porque a dispensa "+
					"satisfaz este gate e o `feature-test-match` só cobra quando há teste a "+
					"confrontar.\n\nEscolha uma: escreva o teste que os cenários pedem, ou "+
					"remova os cenários (e a feature, com `@no-feature: <razão>`) se o "+
					"comportamento não é observável",
				feat, qtd)
		}
	}

	var missing []string
	for _, r := range required {
		if have[r.edge] || optional[string(r.edge)] {
			continue
		}
		missing = append(missing, r.peca+" ("+r.ondeE+")")
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	sort.Strings(missing)
	return Fail, "spec sem trinca completa — falta: " + strings.Join(missing, "; ") +
		". Uma spec sozinha atravessa os gates relacionais (eles não têm o que confrontar) " +
		"e o pipeline passa a certificar trabalho que não existe."
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
		for _, p := range l.TrincaOpcional {
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
				if disp := cfg.Layers[layer].TrincaOpcional; len(disp) > 0 {
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
	tbdRE = regexp.MustCompile(`(?i)@TBD[^\S\n]*:[^\S\n]*([a-z,\s]+)`)
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
func referenciaNoBloco(bloco string) (string, bool) {
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

// provaApontadaPeloNoTest extrai o caminho que a dispensa alega conter a prova.
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
func provaApontadaPeloNoTest(content string) (codigo string, temReferencia bool) {
	bloco := blocoNoTestRE.FindString(content)
	if bloco == "" {
		return "", false
	}
	return referenciaNoBloco(bloco)
}

// testeQueProva procura, entre os TESTES do mapa, algum que mencione o código alegado.
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
func testeQueProva(codigo, root string, g *mapx.Graph) (arquivo string, achou bool) {
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

// dispensasDaSpec devolve as ARESTAS dispensadas pela própria spec.
func dispensasDaSpec(content string) map[string]bool {
	out := map[string]bool{}
	if noTestRE.MatchString(content) {
		out[string(mapx.EdgeTestedBy)] = true
	}
	if noFeatureRE.MatchString(content) {
		// Sem feature não há o que provar por cenário — a dispensa da feature
		// arrasta a do teste, senão o gate cobraria um teste de cenário nenhum.
		out[string(mapx.EdgeCoveredBy)] = true
		out[string(mapx.EdgeTestedBy)] = true
	}
	for peca := range pecasPorDesenvolver(content) {
		out[peca] = true
	}
	return out
}

// cenariosDaFeatureLigada conta os cenários da feature que a spec cobre, e devolve
// o caminho dela. Zero quando não há feature ligada — aí não existe contradição.
func cenariosDaFeatureLigada(n mapx.Node, root string, g *mapx.Graph) (int, string) {
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
		return len(cenarioRE.FindAllString(string(b), -1)), e.To
	}
	return 0, ""
}

// cenarioRE — o Gherkin do projeto pode estar em pt ou en; ambos abrem o cenário no
// início da linha.
var cenarioRE = regexp.MustCompile(`(?m)^\s*(?:Cenário|Cenario|Scenario|Esquema do Cenário|Scenario Outline):`)

// pecasPorDesenvolver lê o `@TBD:` e devolve as arestas cuja peça ainda não foi escrita.
//
// O vocabulário é o do TRABALHO (`code`, `feature`, `test`), e não o das arestas do mapa:
// quem escreve a spec pensa em peças, não em `covered_by`.
func pecasPorDesenvolver(content string) map[string]bool {
	out := map[string]bool{}
	m := tbdRE.FindStringSubmatch(content)
	if m == nil {
		return out
	}
	for _, peca := range strings.Split(m[1], ",") {
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

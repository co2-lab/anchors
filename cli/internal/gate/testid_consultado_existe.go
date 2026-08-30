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
)

// checkTestIDConsultadoExiste — a QUARTA aresta do contrato de testID, a que nenhum
// gate cobria: o flow de ponta a ponta CONSULTA um handle que o código não expõe.
//
// Por que não cabe no `testid-coerente`: aquele gate parte de UMA spec, e a superfície
// e2e é a árvore INTEIRA de flows do projeto — um id de outra tela aparece nela do
// mesmo jeito. Cobrar da spec do InventoryMove um `:acnw-account-new-*` seria acusá-la
// por um handle de contas bancárias; medido ao tentar, deu 829 achados numa spec só,
// todos de outras telas. A pergunta é do PROJETO ("existe flow procurando handle que
// ninguém expõe?"), então o escopo é `project` e os dois lados são varridos de uma vez.
//
// O que este gate pega, medido no app de referência (2026-08-25): 13 ids inventados nos flows, com
// prefixo errado (`review-*` onde a tela emite `revi-*`). A premissa de que "um id
// inexistente já falha ao rodar" é FALSA por dois motivos independentes:
//
//   - o flow precisa CHEGAR à linha. Numa suíte cujo passo 3 já falhava, os ids
//     fantasma dos passos seguintes eram invisíveis — apareceram um por rodada, ao
//     longo de sete medições.
//   - num `assertNotVisible`, o id inexistente faz o teste PASSAR. O REVI-R03
//     asseverava `:review-edit-controls` (que nunca existiu) e provava o oposto do que
//     aparentava: verde por VACUIDADE. É o caso mais grave, e o runner nunca o dá.
//
// Direção ÚNICA (flow → código), deliberadamente. O contrário (código expõe handle que
// nenhum flow consulta) NÃO é ofensa: nem todo elemento marcado precisa de cenário
// automatizado, e o `testid-coerente` já reporta o consumo por spec, onde a pergunta
// tem dono.
func checkTestIDConsultadoExiste(_ string, _ mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	attr := handleDeTeste(cfg)
	if attr == "" {
		// Inferir `testID` por default faria o gate reportar VERDE sobre o que não
		// conferiu — a pior falha possível num medidor.
		return Skip, "o projeto não declara `derived.test_handle` — não há atributo de ancoragem a confrontar"
	}
	flows := arquivosDaSuperficieE2E(root, cfg)
	if len(flows) == 0 {
		// Sem superfície declarada não há o que confrontar. Skip é honesto: "não medi"
		// não é "está limpo".
		return Skip, "o projeto não declara a superfície e2e (`derived.files[...]`) — nada a varrer"
	}

	// O universo do que o CÓDIGO expõe. Varremos o projeto inteiro, não o grafo: um
	// handle pode nascer num componente que nenhuma spec de tela reivindica (o
	// `tabBarButtonTestID` do navegador é o caso real — `:tab-import` não pertence a
	// tela nenhuma), e cobrar dele por ausência no grafo acusaria quem cumpre.
	expostos := handlesExpostosNoProjeto(root, attr, cfg)
	if len(expostos) == 0 {
		return Skip, "nenhum `" + attr + "` encontrado no código — sem lado para confrontar"
	}

	type achado struct{ id, onde string }
	var achados []achado
	vistos := map[string]bool{}

	for _, f := range flows {
		b, err := os.ReadFile(f.path)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(root, f.path)
		for _, id := range handlesConsultados(string(b)) {
			if vistos[id+"|"+rel] {
				continue
			}
			if handleExiste(expostos, id) {
				continue
			}
			vistos[id+"|"+rel] = true
			achados = append(achados, achado{id: id, onde: rel})
		}
	}

	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Slice(achados, func(i, j int) bool {
		if achados[i].id != achados[j].id {
			return achados[i].id < achados[j].id
		}
		return achados[i].onde < achados[j].onde
	})
	// Agrupa por id: o mesmo handle inventado costuma aparecer em vários flows, e um
	// achado por arquivo faria um defeito parecer sete.
	porID := map[string][]string{}
	var ordem []string
	for _, a := range achados {
		if _, ok := porID[a.id]; !ok {
			ordem = append(ordem, a.id)
		}
		porID[a.id] = append(porID[a.id], a.onde)
	}
	var linhas []string
	for _, id := range ordem {
		linhas = append(linhas, fmt.Sprintf("  %s\n      consultado em: %s", id, strings.Join(porID[id], ", ")))
	}
	return Fail, fmt.Sprintf("%d handle(s) consultados por flow e AUSENTES do código:\n%s",
		len(ordem), strings.Join(linhas, "\n"))
}

// arquivoE2E é um flow da superfície de ponta a ponta.
type arquivoE2E struct{ path string }

// arquivosDaSuperficieE2E devolve os CAMINHOS dos flows (o `lerSuperficieE2E` irmão
// devolve só o conteúdo, e aqui o caminho é o que o laudo precisa nomear para o
// achado ser acionável).
//
// A resolução é em DOIS passos e confundi-los custa caro: `surfaces[e2e]` devolve a
// CHAVE da superfície (ex.: "e2e"), não um caminho — quem tem o caminho é
// `files[chave]`.
func arquivosDaSuperficieE2E(root string, cfg *config.Config) []arquivoE2E {
	if cfg == nil || cfg.Derived == nil {
		return nil
	}
	chave := cfg.Derived.Surfaces["e2e"]
	if chave == "" {
		return nil
	}
	// O PRIMEIRO padrão da camada: este gate resolve UM caminho esperado, e uma spec com
	// vários padrões (a de configuração) não muda o que ele pergunta.
	var padrao string
	if ps := cfg.Derived.Files[chave]; len(ps) > 0 {
		padrao = ps[0]
	}
	if padrao == "" {
		for _, ov := range cfg.Derived.Overrides {
			if ps := ov.Files[chave]; len(ps) > 0 {
				padrao = ps[0]
				break
			}
		}
	}
	if padrao == "" {
		return nil
	}
	dir := filepath.Join(root, primeiroSegmentoEstatico(padrao))
	var out []arquivoE2E
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// Só os flows. `.js` da mesma árvore são utilitários de runner (dataLoader),
		// que não consultam handle.
		if ext := strings.ToLower(filepath.Ext(p)); ext != ".yaml" && ext != ".yml" {
			return nil
		}
		out = append(out, arquivoE2E{path: p})
		return nil
	})
	return out
}

// reHandleConsultado captura o valor de `id:` nos flows — a forma como o Maestro (e
// runners equivalentes) referenciam o handle.
//
// Aceita aspas simples e duplas em ALTERNÂNCIA, não por retrovisor: a RE2 do Go não
// tem `\1`, e usá-lo aborta o binário no init (medido). Os dois grupos são
// mutuamente exclusivos — o que casou é o que vem preenchido.
// A aspa dupla é o caso da INTERPOLAÇÃO (`id: "${':bgcr-' + output.data.ns}"`),
// descartada adiante.
var reHandleConsultado = regexp.MustCompile(`(?m)^\s*id:\s*(?:'([^']*)'|"([^"]*)")\s*$`)

// handlesConsultados extrai os handles que um flow procura, descartando as formas em
// que o id é COMPOSTO em runtime — nelas o gate não tem como saber o valor final, e
// acusar seria inventar defeito.
func handlesConsultados(src string) []string {
	var out []string
	for _, m := range reHandleConsultado.FindAllStringSubmatch(src, -1) {
		// Grupo 1 = aspa simples, grupo 2 = aspa dupla; só um vem preenchido.
		bruto := strings.TrimSpace(m[1])
		if bruto == "" {
			bruto = strings.TrimSpace(m[2])
		}
		if bruto == "" {
			continue
		}
		// INTERPOLADO (`${...}`): o valor sai de expressão JS avaliada pelo runner. A
		// parte estática dá para extrair, mas o SUFIXO é dado (namespace do mutante, id
		// de registro) — confrontá-lo exigiria executar o flow.
		if strings.Contains(bruto, "${") {
			continue
		}
		id := strings.TrimPrefix(bruto, ":")
		// REGEX (`.*`, `.+`): o flow procura por PADRÃO, de propósito — o id real leva
		// sufixo dinâmico. Cortamos no primeiro metacaractere e confrontamos a cabeça,
		// que é a parte que o código tem de expor.
		if i := strings.IndexAny(id, ".*+?[](){}|^$\\"); i >= 0 {
			id = id[:i]
		}
		id = strings.TrimRight(id, "-")
		// Cabeça curta demais não identifica nada: confrontá-la daria falso positivo em
		// qualquer projeto. O piso é o mesmo do resto do framework: um segmento com
		// prefixo de unidade + pelo menos uma palavra.
		if len(id) < 4 || !strings.Contains(id, "-") {
			continue
		}
		out = append(out, id)
	}
	return out
}

// handlesExpostosNoProjeto varre o código do projeto e devolve tudo que ele marca com
// o atributo declarado.
//
// Varre por EXTENSÃO de fonte e ignora as pastas que não são código do projeto
// (dependências, artefatos de build). Um handle presente só num arquivo de teste
// (`.test.tsx`) NÃO conta como exposto: teste consulta handle, não o cria — e foi
// exatamente assim que um id fantasma (`:recent-import-item` no app de referência) sobreviveu
// referenciado por um teste de unidade.
func handlesExpostosNoProjeto(root, attr string, cfg *config.Config) []string {
	ignorar := map[string]bool{
		"node_modules": true, ".git": true, "dist": true, "build": true,
		".anchors": true, "coverage": true, "ios": true, "android": true,
		"test-output": true, ".next": true,
	}
	var out []string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if ignorar[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		switch strings.ToLower(filepath.Ext(p)) {
		case ".tsx", ".ts", ".jsx", ".js", ".kt", ".swift", ".vue", ".svelte":
		default:
			return nil
		}
		// Arquivo de teste não EXPÕE handle — ele consulta.
		if strings.Contains(info.Name(), ".test.") || strings.Contains(info.Name(), ".spec.") {
			return nil
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return nil
		}
		src := string(b)
		out = append(out, testIDsExpostos(src, attr)...)
		// O reconhecedor compartilhado é CONSERVADOR de propósito: ele responde "quem é
		// o dono deste handle", e para isso só conta o que está colado ao atributo. Aqui
		// a pergunta é outra — "este handle existe em ALGUM lugar do código?" — e as duas
		// formas abaixo escapavam dele, gerando falso positivo (medido: 2 dos 18 achados
		// da primeira rodada):
		//
		//   testID={btn.testID ?? `:aerx-alert-sheet-button-${i}`}   fallback com template
		//   Family: ':navigate-to-family',                            handle em tabela
		//
		// Num gate que ACUSA ausência, o erro de ser frouxo (deixar passar um id órfão)
		// é muito menos grave que o de ser estrito (acusar quem cumpre) — por isso aqui
		// colhemos todo literal com cara de handle do arquivo.
		out = append(out, handlesLiteraisSoltos(src)...)
		out = append(out, sufixosCompostosDeProp(src)...)
		return nil
	})
	return out
}

// handleExiste confronta um handle consultado contra o universo exposto.
//
// Não é igualdade de string: o exposto pode ser um TEMPLATE (`item-*`, de
// “ testID={`:item-${id}`} “) e o consultado uma instância — ou o inverso, quando o
// flow procura por padrão. Ambos os lados podem carregar curinga, e os dois casam.
func handleExiste(expostos []string, id string) bool {
	if cobertoPor(expostos, id) {
		return true
	}
	// O consultado é a CABEÇA de um id composto no código: `:revi-review-tx` casa
	// `:revi-review-tx-*`. Sem isto, todo flow que procura por prefixo viraria achado.
	for _, e := range expostos {
		nu := strings.TrimPrefix(e, ":")
		if strings.HasPrefix(nu, id) {
			return true
		}
		// SUFIXO composto em componente filho (`*-row-*`): a cabeça vem de quem
		// renderiza, então só o miolo pode ser confrontado. Casa `ctdt-list-row-0`
		// contra o `-row-` que o TransactionList compõe.
		//
		// `Contains` E `HasSuffix`, e a segunda não é redundante: o id pode TERMINAR no
		// segmento composto, sem nada depois. `:reis-register-input-password-toggle`
		// nasce de `${testID}-toggle` no atom Input e não contém "-toggle-" — só
		// termina em "-toggle". Sem esta metade o gate acusava de inexistente um handle
		// que existe, e eu cheguei a "consertar" um flow correto por causa disso: troquei
		// o toque no olho por um toque em ponto neutro, que APAGOU a senha já digitada.
		if strings.HasPrefix(nu, "*") {
			if meio := strings.Trim(nu, "*"); meio != "" {
				if strings.Contains(id, meio) || strings.HasSuffix(id, strings.TrimSuffix(meio, "-")) {
					return true
				}
			}
		}
	}
	return false
}

// reHandleLiteralSolto casa qualquer literal MARCADO (`':foo-bar'`) no fonte, sem exigir
// que esteja colado ao atributo de teste.
//
// A marca (`:`) é o que torna isso seguro: no app de referência todo handle a carrega por convenção, e
// exigi-la evita colher string de domínio ('card_invoice', 'pt-BR') como se fosse handle.
// Num projeto que não marque, este reconhecedor simplesmente não acha nada — e o gate cai
// no reconhecedor estrito, que é o comportamento correto por omissão.
var reHandleLiteralSolto = regexp.MustCompile("[\"'`](:[a-zA-Z][a-zA-Z0-9._-]*)")

// handlesLiteraisSoltos colhe os literais marcados do arquivo, normalizando o template
// (`:foo-${i}` vira `foo-*`) do mesmo jeito que o reconhecedor estrito faz.
func handlesLiteraisSoltos(src string) []string {
	var out []string
	for _, m := range reHandleLiteralSolto.FindAllStringSubmatch(src, -1) {
		id := m[1]
		// Template interrompido pelo `${`: o regex para no `$`, então um id que termine
		// em `-` é cabeça de composição e entra como curinga.
		if strings.HasSuffix(id, "-") {
			out = append(out, strings.TrimSuffix(id, "-")+"-*")
			continue
		}
		out = append(out, id)
	}
	return out
}

// reSufixoCompostoDeProp casa a composição de handle a partir de uma PROP recebida:
//
//	testID={testID ? `${testID}-row-${index}` : undefined}
//
// O id final (`:ctdt-list-row-0`) só existe em runtime: a CABEÇA vem de quem renderiza
// (`testID=":ctdt-list"`) e o SUFIXO nasce aqui, num componente filho que não sabe de
// que tela é. Nenhum dos dois arquivos contém a string inteira.
//
// Sem reconhecer isto o gate acusava todo flow que toca linha de lista — 8 flows do app de referência
// numa tacada, todos corretos. O que colhemos é o SUFIXO (`-row-*`), que adiante é
// casado contra qualquer cabeça exposta.
var reSufixoCompostoDeProp = regexp.MustCompile("`\\$\\{[a-zA-Z_][a-zA-Z0-9_.]*\\}(-[a-zA-Z][a-zA-Z0-9-]*)")

// sufixosCompostosDeProp devolve os sufixos de composição na forma `*-row-*`: casam
// qualquer cabeça, porque a cabeça vem de fora do arquivo.
func sufixosCompostosDeProp(src string) []string {
	var out []string
	for _, m := range reSufixoCompostoDeProp.FindAllStringSubmatch(src, -1) {
		// `-row-${index}` → `-row-`; o que interessa é o segmento ESTÁTICO.
		seg := strings.TrimRight(m[1], "-")
		if seg == "" {
			continue
		}
		out = append(out, "*"+seg+"-*")
	}
	return out
}

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

// spec-feature-match: todo REQUISITO declarado na spec tem ao menos um CENÁRIO na feature.
//
// É a ponta que faltava na trinca. O `feature-test-match` confronta feature→teste; o
// `trinca-completa` confronta que as PEÇAS existem. Ninguém confrontava spec→feature — e
// é aí que mora um buraco silencioso: a spec declara `XXXXX-X02`, a feature não tem
// cenário nenhum com essa tag, e o requisito atravessa o pipeline inteiro sem que nada o
// verifique. Todos os gates ficam verdes: a spec tem código, a feature existe, a feature
// bate com o teste. O requisito simplesmente não é de ninguém.
//
// Medido num projeto real: 11 de 287 specs com feature tinham requisito sem cenário.
//
// A régua é a mesma dos outros gates relacionais — CÓDIGO, não prosa: cada
// `{CODE}-{letra}{NN}` que a spec DEFINE precisa aparecer como tag de cenário na feature.
// Definir é diferente de citar: uma spec que menciona o código de outra unidade (numa
// Tabela de Dependências, por exemplo) não contrai obrigação nenhuma.
//
// Opt-out honesto (CONCEPT §5.1): `@no-scenario: <razão>` na linha do requisito dispensa
// aquele requisito específico, com a razão escrita. Serve para o que é verdadeiramente
// não-observável por cenário — e deixa o rastro de que foi decisão, não esquecimento.
func checkSpecFeatureMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a obrigação é da spec — é ela que declara os requisitos"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	// A feature se alcança pela aresta `covered-by` (spec → feature).
	var features []string
	for _, e := range g.Edges {
		if e.From == n.ID && e.Type == "covered-by" {
			features = append(features, e.To)
		}
	}
	if len(features) == 0 {
		// Ausência de feature é problema do `trinca-completa`, não deste gate — cada um
		// acusa uma coisa, senão o mesmo defeito aparece duas vezes no relatório.
		return Skip, "spec sem feature ligada (`covered-by`) — a ausência da peça é do gate trinca-completa"
	}

	// `@no-feature` ARRASTA a dispensa de cenário de TODOS os requisitos.
	//
	// A implicação é lógica, não convenção: a tag afirma que esta spec não tem feature, e
	// sem feature nenhum requisito dela pode ter cenário. Cobrar mesmo assim obrigaria o
	// autor a repetir `@no-scenario` em cada linha da tabela para dizer o que a spec já
	// disse uma vez — três marcações (`@no-test`, `@no-feature`, `@no-scenario`) para uma
	// única decisão.
	//
	// E a repetição não é só verbosa, é frágil: nada impede alguém remover o
	// `@no-scenario` de uma linha e deixar o `@no-feature`, e aí a spec afirma duas coisas
	// contraditórias. Uma decisão, um lugar.
	//
	// O `trinca-completa` já faz esse mesmo arrasto para o teste (`@no-feature` implica
	// `@no-test`, porque sem cenário não há o que provar); aqui ele se completa.
	if noFeatureRE.MatchString(content) {
		return Skip, "a spec declara `@no-feature` — sem feature não há cenário a cobrar de requisito nenhum"
	}

	declarados := requisitosDefinidos(content)
	if len(declarados) == 0 {
		return Skip, "a spec não define requisito com código — nada a cobrir"
	}

	// une os códigos citados como TAG em qualquer feature ligada (uma spec pode ser
	// coberta por mais de uma feature; o requisito só precisa estar em alguma).
	cobertos := map[string]bool{}
	for _, f := range features {
		b, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			continue
		}
		for _, sc := range parseFeatureScenarios(string(b)) {
			// TODOS os códigos da tag-line contam como cobertos — um cenário pode provar
			// mais de um requisito, e o dialeto do projeto co-etiqueta.
			// A RAIZ também conta: `@CATRX-B02#01` é um cenário do requisito `CATRX-B02`,
			// e é a raiz que a spec cataloga. O sufixo dá identidade ao cenário sem
			// criar requisito novo — sem esta linha, numerar cenários faria o gate
			// acusar de repente todo requisito que ganhou mais de um caso.
			for _, c := range sc.Codes {
				cobertos[c] = true
				cobertos[CodeRaiz(c)] = true
			}
			cobertos[sc.Code] = true
			cobertos[CodeRaiz(sc.Code)] = true
		}
	}

	var faltando []string
	for _, code := range declarados {
		if !cobertos[code] {
			faltando = append(faltando, code)
		}
	}
	if len(faltando) == 0 {
		return Pass, ""
	}
	sort.Strings(faltando)

	mostra := faltando
	sufixo := ""
	if len(mostra) > 8 {
		mostra, sufixo = mostra[:8], fmt.Sprintf(" (e mais %d)", len(faltando)-8)
	}
	return Fail, fmt.Sprintf("%d requisito(s) declarado(s) sem cenário na feature: %s%s. "+
		"Um requisito sem cenário atravessa o pipeline inteiro sem ser verificado — e todos "+
		"os outros gates ficam VERDES, porque a spec tem código, a feature existe e ela bate "+
		"com o teste. Escreva o cenário, ou dispense na linha do requisito com "+
		"`@no-scenario: <razão>`",
		len(faltando), strings.Join(mostra, ", "), sufixo)
}

// requisitosDefinidos extrai os códigos que a spec DEFINE — não os que ela cita.
//
// A distinção é a mesma que o gate `rule-types` já faz, e é o que separa um gate útil de
// um gerador de ruído: uma spec cita códigos de outras unidades o tempo todo (na Tabela
// de Dependências, em notas, em referências cruzadas) e não contrai obrigação alguma por
// isso. Define quem coloca o código no INÍCIO de uma linha, de um item de lista, de um
// título de seção, ou na PRIMEIRA célula de uma tabela — as formas em que uma régua
// enuncia um requisito.
func requisitosDefinidos(content string) []string {
	vistos := map[string]bool{}
	var out []string
	for _, linha := range strings.Split(content, "\n") {
		if dispensadoPorNoScenario(linha) {
			continue
		}
		m := defineRuleCaptureRE.FindStringSubmatch(linha)
		if m == nil {
			continue
		}
		code := m[1]
		if !vistos[code] {
			vistos[code] = true
			out = append(out, code)
		}
	}
	return out
}

// defineRuleCaptureRE é o `definesRuleRE` do gate rule-types com o código CAPTURADO —
// mesma gramática de "definir" (código no início da linha/item/título, ou primeira célula
// de tabela), porque as duas perguntas dependem da mesma distinção: definir ≠ citar.
var defineRuleCaptureRE = regexp.MustCompile(
	"(?m)^\\s*(?:#{2,6}\\s+|[-*]\\s+\\**|\\|\\s*)`?\\*{0,2}([A-Z0-9]" + config.CodeLengthPattern() + "-[A-Z]\\d{2})")

// noScenarioRE — o opt-out por requisito, com razão obrigatória depois dos dois-pontos.
var noScenarioRE = regexp.MustCompile(`@no-scenario[^\S\n]*:[^\S\n]*\S+`)

// dispensadoPorNoScenario: o opt-out precisa de RAZÃO escrita. `[^\S\n]` = espaço/tab mas
// não quebra de linha, senão a razão seria "achada" na linha seguinte e um marcador nu
// passaria — que é justamente o que a dispensa não pode permitir.
func dispensadoPorNoScenario(linha string) bool {
	return noScenarioRE.MatchString(linha)
}

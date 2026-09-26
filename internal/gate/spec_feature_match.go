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
		return Skip, i18n.T("gate.spec_feature.skip_not_spec")
	}
	// A rule ALIAS must resolve before anything else is read: a dangling one would drop a
	// requirement from every scenario check without standing for any other rule.
	if bad := invalidRuleAliases(content); len(bad) > 0 {
		return Fail, i18n.T("gate.spec_feature.alias_invalid", len(bad), strings.Join(bad, "; "))
	}
	if g == nil {
		return pendingNoMap()
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
		return Skip, i18n.T("gate.spec_feature.skip_no_feature")
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
	if unitWaiver(noFeatureRE, content) {
		return Skip, i18n.T("gate.spec_feature.skip_waived")
	}

	declarados := definedRequirements(content)
	if len(declarados) == 0 {
		return Skip, i18n.T("gate.spec_feature.skip_no_requirements")
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
				cobertos[RootCode(c)] = true
			}
			cobertos[sc.Code] = true
			cobertos[RootCode(sc.Code)] = true
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
		mostra, sufixo = mostra[:8], i18n.T("gate.spec_feature.and_more", len(faltando)-8)
	}
	return Fail, i18n.T("gate.spec_feature.missing", len(faltando), strings.Join(mostra, ", "), sufixo)
}

// definedRequirements extrai os códigos que a spec DEFINE — não os que ela cita.
//
// A distinção é a mesma que o gate `rule-types` já faz, e é o que separa um gate útil de
// um gerador de ruído: uma spec cita códigos de outras unidades o tempo todo (na Tabela
// de Dependências, em notas, em referências cruzadas) e não contrai obrigação alguma por
// isso. Define quem coloca o código no INÍCIO de uma linha, de um item de lista, de um
// título de seção, ou na PRIMEIRA célula de uma tabela — as formas em que uma régua
// enuncia um requisito.
func definedRequirements(content string) []string {
	vistos := map[string]bool{}
	var out []string
	re := defineRuleCaptureRE()
	for _, linha := range strings.Split(content, "\n") {
		if waivedByNoScenario(linha) {
			continue
		}
		// An ALIAS (`REF[CODE-B05]: reason`) has no scenario of its own: the rule it
		// points at is the one the scenario and the test prove.
		if _, ok := ruleAliasTarget(linha); ok {
			continue
		}
		m := re.FindStringSubmatch(linha)
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

// defineRuleCaptureRE is the rule-types gate's `definesRuleRE` with the code CAPTURED —
// the same grammar of "defining" (code at the start of the line/item/heading, or the
// first table cell), because both questions depend on the same distinction: define ≠ cite.
//
// Compiled per CALL and not in a `var` — the same rule as `codeRE` (rule_implemented.go):
// the code length comes from `code_lengths`, loaded AFTER the globals. In a `var` this
// regex froze the default `[5]`, and in a `[4]` project it matched no requirement at all —
// MEASURED in MIF (2026-09-23): `spec-feature-match` and `scenario-coverage` undetermined
// on all 691 specs, blocking and blind.
func defineRuleCaptureRE() *regexp.Regexp {
	return regexp.MustCompile(
		"(?m)^\\s*(?:#{2,6}\\s+|[-*]\\s+\\**|\\|\\s*)`?\\*{0,2}([A-Z0-9]" + config.CodeLengthPattern() + "-[A-Z]\\d{2})")
}

// ruleAliasRE is a rule declared as an ALIAS of another rule of the same spec:
//
//	| `DCFRD-E01` | REF[DCFRD-B09]: a spec missing from disk is the stale map B09 answers |
//
// It exists for the rule that must be catalogued under one letter while its behaviour is
// already stated under another — a failure (`-E`) the spec had declared as a behaviour
// (`-B`). Writing the rule twice would put one decision in two places, and the copies
// would drift; waiving the letter would lose the catalogue. The alias keeps both: the
// code exists where its letter is looked for, and the statement lives once, in the
// target. The reason is mandatory, as for every waiver: it says why the target answers.
var ruleAliasRE = regexp.MustCompile("^REF\\[`?([A-Z0-9]{3,6}-[A-Z]\\d{2})`?\\]([^\\S\\n]*:[^\\S\\n]*[^\\s|])?")

// ruleAliasTarget returns the rule a DEFINING line aliases, if it is an alias.
//
// The alias is the FIRST thing after the rule's code — the whole statement of the rule.
// Anywhere else on the line it is prose: a rule that explains the syntax by quoting
// `REF[CODE-B05]` would otherwise turn into an alias of that example.
func ruleAliasTarget(linha string) (string, bool) {
	m := aliasOf(linha)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// aliasOf matches the alias that opens a defining line's statement, or returns nil.
func aliasOf(linha string) []string {
	loc := defineRuleCaptureRE().FindStringSubmatchIndex(linha)
	if loc == nil {
		return nil
	}
	rest := strings.TrimLeft(linha[loc[3]:], "`* \t|—–-:")
	return ruleAliasRE.FindStringSubmatch(rest)
}

// invalidRuleAliases lists every alias that does not stand for a real rule: a target the
// spec does not define, a target that is itself an alias (a chain hides where the
// statement lives), or no written reason.
func invalidRuleAliases(content string) []string {
	re := defineRuleCaptureRE()
	defined := map[string]bool{}
	aliased := map[string]bool{}
	for _, linha := range strings.Split(content, "\n") {
		if m := re.FindStringSubmatch(linha); m != nil {
			defined[m[1]] = true
			if _, ok := ruleAliasTarget(linha); ok {
				aliased[m[1]] = true
			}
		}
	}
	var bad []string
	for _, linha := range strings.Split(content, "\n") {
		m := re.FindStringSubmatch(linha)
		if m == nil {
			continue
		}
		a := aliasOf(linha)
		if a == nil {
			continue
		}
		switch {
		case !defined[a[1]]:
			bad = append(bad, i18n.T("gate.spec_feature.alias_unknown", m[1], a[1]))
		case aliased[a[1]]:
			bad = append(bad, i18n.T("gate.spec_feature.alias_chain", m[1], a[1]))
		case a[2] == "":
			bad = append(bad, i18n.T("gate.spec_feature.alias_no_reason", m[1], a[1]))
		}
	}
	return bad
}

// noScenarioRE — o opt-out por requisito, com razão obrigatória depois dos dois-pontos.
var noScenarioRE = regexp.MustCompile(`@no-scenario[^\S\n]*:[^\S\n]*\S+`)

// waivedByNoScenario: o opt-out precisa de RAZÃO escrita. `[^\S\n]` = espaço/tab mas
// não quebra de linha, senão a razão seria "achada" na linha seguinte e um marcador nu
// passaria — que é justamente o que a dispensa não pode permitir.
func waivedByNoScenario(linha string) bool {
	return noScenarioRE.MatchString(linha)
}

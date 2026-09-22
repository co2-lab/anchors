package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/flagx"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- the FLAG AXIS: the scenarios a flag's value opens ---
//
// A feature flag multiplies the code's paths without multiplying the spec. `if
// flag("new-checkout")` creates two behaviours, and the spec describes ONE — or, worse,
// describes both mixed into a sentence that does not say which holds when.
//
// The cost shows up in three places, and all three are silences no other gate sees:
// in review, where nobody knows which branch is current and which is dying; in the test,
// where the ON path gets covered and the OFF path gets no proof; and in removal, where the
// "temporary" flag turns three years old and deleting it becomes archaeology.
//
// The gates here answer distinct questions:
//
//	is the condition written in the grammar?    flag-scenario-grammar   (fails)
//	does the cited scenario exist?              flag-scenario-exists    (fails)
//	does the flag declare the ABSENT case?      flag-scenarios-complete (fails)
//	does every scenario have a test?            flag-covered            (informs)
//
// What is deliberately NOT here is reading the flag's REAL value. Anchors has no access
// to the flag service, must not have, and the value changes per user and per minute. What
// it confronts is the DECLARED scenario.

// noAbsentRE — the declared waiver for a flag that cannot be absent.
//
// `@no-absent` and not `@TBD` because the two assert different things: a flag read from a
// local constant will NEVER have the absent case, which is a permanent waiver, not a debt
// somebody still owes. A bare marker does not count, by the same rule as every other
// opt-out here — the reason is what answers the question somebody asks six months later.
var noAbsentRE = regexp.MustCompile(`(?i)(^|[^` + "`" + `])@no-absent[^\S\n]*:[^\S\n]*\S+`)

// gatedByRE finds `@gated-by` citations in a spec. Mirrors `scan.gatedByRE`, including
// the fixed `G`: `@gated-by CRED-B03` is not a flag, it is a typo, and matching any letter
// would have the gate hunt for a scenario that can never exist.
var gatedByRE = regexp.MustCompile("@gated-by\\s+`?([A-Z0-9]{3,6}-G[0-9]{2})`?")

// --- is the condition written in the grammar? ---
//
// The first gate of the axis, and the one the others depend on: a condition the grammar
// refuses is a scenario nothing downstream can confront. The parser keeps the row rather
// than dropping it precisely so this gate can report it — a line that vanishes is the
// silence the gates exist to end.
func checkFlagScenarioGrammar(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFlag {
		return Skip, i18n.T("gate.flag_scenario_grammar.skip_not_flag")
	}
	f := flagx.ParseContent(content, n.ID)
	if len(f.Scenarios) == 0 {
		return Skip, i18n.T("gate.flag_scenarios_complete.no_scenarios")
	}
	var bad []string
	for _, s := range f.Scenarios {
		if s.Err != nil {
			bad = append(bad, fmt.Sprintf("%s (line %d): %v", s.Code, s.Line, s.Err))
		}
	}
	if len(bad) == 0 {
		return Pass, ""
	}
	sort.Strings(bad)
	return Fail, fmt.Sprintf(i18n.T("gate.flag_scenario_grammar.unparsed"),
		len(bad), strings.Join(bad, "\n  "))
}

// --- does the flag declare the ABSENT case? ---
//
// The scenario everyone forgets and the one that breaks hardest: the flag that does not
// answer — flag service down, a new environment, a local test. Without a declared
// behaviour, what happens is whatever the library's default happens to be, and that is
// discovered in production.
//
// Measured against the market's own tools: every one of them can EXPRESS this case
// (Flagsmith calls it `Is Not Set`, and LaunchDarkly falls back to the serving default),
// and not one of them REQUIRES it. That gap is precisely where this gate earns its keep.
func checkFlagScenariosComplete(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFlag {
		return Skip, i18n.T("gate.flag_scenarios_complete.skip_not_flag")
	}
	f := flagx.ParseContent(content, n.ID)
	if len(f.Scenarios) == 0 {
		return Skip, i18n.T("gate.flag_scenarios_complete.no_scenarios")
	}
	if f.Absent() {
		return Pass, ""
	}
	// The declared way out, and it is `@no-` rather than `@TBD`: a flag that genuinely
	// cannot be absent (read from a local constant, say) will never have the case, which
	// is a permanent waiver and not a debt.
	if noAbsentRE.MatchString(content) {
		return Pass, ""
	}
	return Fail, i18n.T("gate.flag_scenarios_complete.no_absent")
}

// --- does the cited scenario exist? ---
//
// This axis's `ref-resolves`. A spec declaring `@gated-by CHKUT-G02` against a scenario
// nobody wrote is a rule governed by a condition that does not exist — and the map does
// not turn it into an edge, on purpose, so that it is reported here with the code in hand
// instead of making the relational gates confront emptiness.
func checkFlagScenarioExists(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.flag_scenario_exists.skip_not_spec")
	}
	cited := gatedByRE.FindAllStringSubmatch(content, -1)
	if len(cited) == 0 {
		return Skip, i18n.T("gate.flag_scenario_exists.no_citation")
	}
	flags, err := flagx.Load(root)
	if err != nil {
		return pendingNoMap()
	}
	known := flagx.ByCode(flags)

	missing := map[string]bool{}
	for _, m := range cited {
		if _, ok := known[m[1]]; !ok {
			missing[m[1]] = true
		}
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	var codes []string
	for c := range missing {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	return Fail, fmt.Sprintf(i18n.T("gate.flag_scenario_exists.unknown"),
		len(codes), strings.Join(codes, ", "))
}

// --- does every scenario govern some rule? ---
//
// A VOLTA do `flag-scenario-exists`, e a pergunta é a inversa: aquele confronta a spec
// que cita um cenário inexistente; este confronta o cenário que ninguém cita.
//
// Um cenário que nenhuma regra invoca é caminho DECLARADO e não governado: alguém
// escreveu "quando o valor for X, então Y" e nenhuma spec diz que regra vale sob ele. O
// código ramifica na flag de qualquer forma — o que falta é o registro de quem depende
// disso, e é justamente esse registro que o eixo existe para manter.
//
// É o mesmo par assimétrico que os dois sentidos da trinca mostraram hoje, e que o
// carimbo da documentação já tinha mostrado: perguntar só a ida deixa o resto invisível.
//
// INFORMATIVO, e por uma razão de ordem: a flag nasce ANTES das specs que a citam — quem
// escreve a flag está decidindo os caminhos, e as regras vêm depois. Um gate bloqueante
// aqui exigiria que as duas pontas nascessem no mesmo commit, que não é como o trabalho
// acontece.
func checkFlagScenarioGoverns(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFlag {
		return Skip, i18n.T("gate.flag_scenario_governs.skip_not_flag")
	}
	if g == nil {
		return pendingNoMap()
	}
	f := flagx.ParseContent(content, n.ID)
	if len(f.Scenarios) == 0 {
		return Skip, i18n.T("gate.flag_scenarios_complete.no_scenarios")
	}

	// Quem cita: as arestas `gated-by` que CHEGAM aqui carregam, no Method, o código do
	// cenário invocado — a mesma mecânica que o `doctrine-realized` usa para saber quem
	// realiza uma regra de produto.
	governados := map[string]bool{}
	for _, e := range g.Neighbors(n.ID).In {
		if e.Type == mapx.EdgeGatedBy {
			governados[e.Method] = true
		}
	}

	var soltos []string
	for _, s := range f.Scenarios {
		if governados[s.Code] {
			continue
		}
		// A dispensa declarada, por cenário: há caminho que existe e nenhuma regra
		// precisa nomear — o `off` de uma flag cujo desligado é simplesmente o
		// comportamento antigo, já governado pelas regras que sempre valeram.
		if noGovernRE.MatchString(s.Then) {
			continue
		}
		soltos = append(soltos, s.Code)
	}
	if len(soltos) == 0 {
		return Pass, ""
	}
	sort.Strings(soltos)
	return Fail, fmt.Sprintf(i18n.T("gate.flag_scenario_governs.ungoverned"),
		len(soltos), len(f.Scenarios), strings.Join(soltos, ", "))
}

// noGovernRE — a dispensa por cenário, com razão obrigatória.
//
// `@no-govern: <razão>` e não `@TBD`, porque as duas afirmações são diferentes: este
// cenário NUNCA vai ter regra própria (o desligado que devolve ao comportamento antigo),
// e isso é dispensa permanente, não dívida.
var noGovernRE = regexp.MustCompile(`(?i)@no-govern[^\S\n]*:[^\S\n]*\S+`)

// --- does every scenario have a test? ---
//
// Per SCENARIO, not per flag, and that is the whole point: the flag whose ON path is
// tested and whose OFF path is not passes a per-flag rule while leaving exactly the branch
// that will break unproven.
//
// TWO QUESTIONS, and collapsing them loses the answer to both:
//
//	ESCRITO   some test file names this scenario code     (static, always available)
//	VERDE     that test RAN and PASSED                    (needs ingested execution)
//
// The first version asked only the second, and on a project that has never ingested a
// report it accused every scenario — including the ones with tests written and passing.
// Measured on this very repository: zero `proven_codes` in the whole graph.
//
// But answering only the first would be the opposite error, and a worse one: a test can
// be written and never run, or run and fail. "Somebody wrote it" is not "it works".
//
// So the gate reports BOTH, and names which scenarios are in which state. The distinction
// is the finding: a scenario with a test that never ran is a different problem, with a
// different fix, than a scenario nobody tested.
//
// INFORMATIVE, because the scenario written today and tested next commit is ordinary
// work, not a defect.
func checkFlagCovered(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFlag {
		return Skip, i18n.T("gate.flag_covered.skip_not_flag")
	}
	if g == nil {
		return pendingNoMap()
	}
	f := flagx.ParseContent(content, n.ID)
	if len(f.Scenarios) == 0 {
		return Skip, i18n.T("gate.flag_scenarios_complete.no_scenarios")
	}

	// O VERDE vem do sinal DESTE no', como o `scenario-coverage` le' o da spec: a
	// ingestao cruza os codigos provados com os que o no' DECLARA, e a flag declara os
	// seus. Procurar o sinal nos nos de teste era a leitura errada — nenhum arquivo de
	// teste "declara" o cenario da flag, e por isso nunca havia sinal onde eu olhava.
	green := map[string]bool{}
	ingested := n.Signal != nil
	if ingested {
		for _, c := range n.Signal.ProvenCodes {
			green[c] = true
		}
	}
	written := codesNamedByTests(scenarioCodes(f), root, g, n.ID)

	var noTest, notGreen []string
	for _, s := range f.Scenarios {
		switch {
		case green[s.Code]:
			// provado: nada a cobrar
		case written[s.Code]:
			notGreen = append(notGreen, s.Code)
		default:
			noTest = append(noTest, s.Code)
		}
	}
	sort.Strings(noTest)
	sort.Strings(notGreen)

	if len(noTest) == 0 && len(notGreen) == 0 {
		return Pass, ""
	}

	var b strings.Builder
	if len(noTest) > 0 {
		fmt.Fprintf(&b, i18n.T("gate.flag_covered.no_test"),
			len(noTest), len(f.Scenarios), strings.Join(noTest, ", "))
	}
	if len(notGreen) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		// A EXECUÇÃO é o que falta, não o teste — e a saída tem de dizer qual das duas,
		// porque o conserto é outro: escrever um teste versus rodar o que já existe.
		if !ingested {
			fmt.Fprintf(&b, i18n.T("gate.flag_covered.written_not_ingested"),
				len(notGreen), strings.Join(notGreen, ", "))
		} else {
			fmt.Fprintf(&b, i18n.T("gate.flag_covered.written_not_green"),
				len(notGreen), strings.Join(notGreen, ", "))
		}
	}
	return Fail, strings.TrimRight(b.String(), "\n")
}

// scenarioCodes lista os codigos que esta flag declara.
func scenarioCodes(f flagx.Flag) []string {
	out := make([]string, 0, len(f.Scenarios))
	for _, s := range f.Scenarios {
		out = append(out, s.Code)
	}
	return out
}

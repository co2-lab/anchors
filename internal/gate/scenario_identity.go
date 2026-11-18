// @anchors
//   code: SIGSC
//   ref: SCIDS

package gate

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cenario-identidade: dois cenários da MESMA feature não podem ter o mesmo código.
//
// Uma regra legitimamente tem vários cenários — caminho feliz e alternativos. O que
// não pode é os dois serem indistinguíveis: com o mesmo código, nada liga UM cenário
// a UM teste, e os gates relacionais comparam N títulos contra o mesmo teste (no
// máximo um casa; os outros viram divergência que ninguém consegue resolver).
//
// A saída é o SUFIXO de cenário: `@USBPX-B01#01` e `@USBPX-B01#02` mantêm a regra
// legível no prefixo e dão identidade a cada caso.
//
// Medido no projeto que originou o gate: 204 códigos repetidos, e 65 deles com tags
// de TIPO conflitantes no mesmo código (`@estado` num cenário, `@comportamento` no
// outro) — sinal de código emprestado, não de regra com dois caminhos. Um desses
// (`GLFLX-S02`) provou ser defeito: a spec o define como "Foco", e o segundo cenário
// falava de `onChangeText`, comportamento que não tinha código próprio.
//
// PENDING e não FAIL: numerar cenários é migração, e o gate nasce sobre uma base que
// não conhecia a notação. Quem já migrou fica verde; quem não, vê o que falta.
func checkScenarioIdentity(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, "" // só confronta features
	}
	cenarios := parseFeatureScenarios(content)
	if len(cenarios) == 0 {
		return Skip, i18n.T("gate.scenario_identity.skip_no_scenarios_with_code")
	}

	// Agrupa por código COMPLETO (com sufixo, se houver): é ele que identifica o
	// cenário. Dois cenários com `#01` e `#02` são distintos; dois sem sufixo, não.
	porCodigo := map[string][]string{}
	linhas := map[string][]int{}
	for _, c := range cenarios {
		porCodigo[c.Code] = append(porCodigo[c.Code], c.Title)
		linhas[c.Code] = append(linhas[c.Code], c.Line)
	}

	// The repeats are this gate's own, and not the duplicates mechanism's: its message names
	// the scenarios by their titles and teaches the suffix with the project's code — what the
	// author needs to tell which is which. It names the lines too.
	var repetidos []string
	for cod, titulos := range porCodigo {
		if len(titulos) < 2 {
			continue
		}
		repetidos = append(repetidos, i18n.T("gate.scenario_identity.repeated_item",
			cod, len(titulos), strings.Join(summarizeTitles(titulos), " / ")+" — "+joinInts(linhas[cod])))
	}
	if len(repetidos) > 0 {
		sort.Strings(repetidos)
		return Diverge, i18n.T("gate.scenario_identity.pending_repeated_codes",
			len(repetidos), strings.Join(repetidos, "; "),
			firstCode(repetidos), firstCode(repetidos))
	}
	if copies := copiedBodies(content); len(copies) > 0 {
		return Diverge, i18n.T("gate.scenario_identity.pending_copied_body", len(copies), strings.Join(copies, "; "))
	}
	return Pass, ""
}

// copiedBodies names the scenarios whose steps are the same as an earlier scenario's under
// another title. Indistinguishable by their steps, two scenarios prove one thing twice —
// and in the reference app nine did so because their body was copied from the first and
// never rewritten: "… → does not alert" ending in "Then a finding must come out". The
// steps are compared case and spacing aside, and nothing else: scenarios that differ only
// in their quoted values are the variations of a rule, not copies, and a scenario of one
// step only states an outcome many triggers share.
func copiedBodies(content string) []string {
	first := map[string]string{}
	var out []string
	for _, sc := range scenarioBodies(content) {
		// Two steps at least: a scenario that only states its outcome ("Then the user is
		// available") shares it with every trigger that leads there, and that is not a copy.
		if strings.Count(sc.body, "\n") < 1 {
			continue
		}
		if prev, ok := first[sc.body]; ok && prev != sc.code {
			out = append(out, i18n.T("gate.scenario_identity.copied_item", sc.code, prev))
			continue
		}
		if _, ok := first[sc.body]; !ok {
			first[sc.body] = sc.code
		}
	}
	return out
}

type scenarioBody struct{ code, title, body string }

// scenarioBodies reads each coded scenario with its steps — the lines after its title up to
// the next tag line or title —, normalized: comments out, lower case, spacing collapsed.
func scenarioBodies(content string) []scenarioBody {
	var out []scenarioBody
	var cur *scenarioBody
	var steps []string
	pendingCode := ""
	flush := func() {
		if cur != nil {
			cur.body = strings.Join(steps, "\n")
			out = append(out, *cur)
		}
		cur, steps = nil, nil
	}
	for _, ln := range strings.Split(content, "\n") {
		t := strings.TrimSpace(ln)
		if ms := featScenarioCodeRE.FindAllStringSubmatch(ln, -1); ms != nil {
			flush()
			pendingCode = ms[0][1] + ms[0][2]
			continue
		}
		if m := featTitleRE.FindStringSubmatch(ln); m != nil {
			flush()
			if pendingCode != "" {
				cur = &scenarioBody{code: pendingCode, title: m[1]}
			}
			pendingCode = ""
			continue
		}
		if cur == nil || t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "@") {
			continue // a tag line of no code belongs to the next scenario, whose title closes this one
		}
		steps = append(steps, strings.Join(strings.Fields(strings.ToLower(t)), " "))
	}
	flush()
	return out
}

// summarizeTitles encurta os títulos para a mensagem caber — o endereço é o código,
// o título só ajuda a reconhecer qual cenário é qual.
func summarizeTitles(ts []string) []string {
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		if len([]rune(t)) > 40 {
			t = string([]rune(t)[:40]) + "…"
		}
		out = append(out, t)
	}
	return out
}

// firstCode extrai o código do primeiro achado, para o exemplo da mensagem
// falar do caso REAL do projeto em vez de um genérico.
func firstCode(repetidos []string) string {
	if len(repetidos) == 0 {
		return "XXXXX-B01"
	}
	if i := strings.Index(repetidos[0], " "); i > 0 {
		return repetidos[0][:i]
	}
	return repetidos[0]
}

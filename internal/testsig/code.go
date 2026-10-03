// @anchors
//   ref: RCGRL

package testsig

import (
	"regexp"
	"strings"
)

// ruleLetters são as letras de tipo de regra em vigor. Começam nas canônicas e são
// substituídas por SetRuleLetters quando a config do projeto é carregada — o vocabulário
// é do PROJETO (`rule_types`), não uma lista fixa do engine.
//
// Por que importa aqui: este regex extrai o scenario-code dos NOMES dos casos no relatório
// JUnit. Preso às canônicas, um cenário de letra declarada pelo projeto (ex.: `-I01`, de
// Invariant) nunca é reconhecido como provado — e o requisito aparece "sem teste verde"
// mesmo tendo um teste que passa.
// O DEFAULT tem de ser o mesmo do `config.DefaultRuleLetters`, e ele divergiu.
//
// Medido no app de referência: um projeto sem `rule_types:` declarado usa o default do `config`
// (`SRVAXBNMDEIQF`, que inclui `I` de Invariant). O `SetRuleLetters` é chamado com esse
// valor e a divergência não apareceria — MAS o `ingest` só o chama quando a config
// carrega, e qualquer caminho que leia o relatório antes disso usa esta constante.
//
// O sintoma foi o que o comentário acima descreve: 11 casos verdes, três cenários
// reconhecidos como provados, e os invariantes (`GLCGL-I01`, `I02`, `I03`) aparecendo
// "sem teste verde" com teste passando.
//
// Duas cópias da mesma lista divergem na primeira letra nova — e esta ficou três atrás
// (`E`, `I`, `Q`, `F`). O comentário do `config` já registra que "é a terceira vez que a
// lista fica para trás de uma letra nova".
var ruleLetters = "SRVAXBNMDEIQWGP"

// codeLenPattern espelha `config.CodeLengthPattern()`. Duplicado pelo mesmo motivo que
// `ruleLetters`: o pacote testsig não depende de scan nem de config, e o comprimento do
// código é do PROJETO — preso a um número fixo, um relatório com códigos de outro tamanho
// não casaria nenhum caso e todo requisito apareceria "sem teste verde".
var codeLenPattern = "{4,5}"

// SetCodeLenPattern ajusta o padrão de comprimento (ver config.CodeLengthPattern).
func SetCodeLenPattern(p string) {
	if p != "" {
		codeLenPattern = p
	}
}

// SetRuleLetters ajusta o vocabulário de letras usado ao ler o relatório de execução.
func SetRuleLetters(letters string) {
	if letters != "" {
		ruleLetters = letters
	}
}

// mustCodeRE devolve o regex de código de cenário — MESMA gramática do scan
// (TRACEABILITY §3), duplicada aqui para o pacote testsig não depender de scan.
func mustCodeRE() *regexp.Regexp {
	return regexp.MustCompile(codeCore() + `\b`)
}

func codeCore() string {
	return `\b[A-Z0-9]` + codeLenPattern + `-(?:[` + regexp.QuoteMeta(ruleLetters) + `]\d{2}(?:-[a-z][a-z0-9-]*)?|DS-[A-Za-z0-9-]+|VR(?:-[` + regexp.QuoteMeta(ruleLetters) + `]\d{2})?)`
}

// mustScenarioRE is the scenario code WITH its variant suffix, when it carries one:
// `RDCH-B02#02` is a scenario of its own, and the rule `RDCH-B02` is proven only when each
// of its scenarios is. Read without the suffix, a skipped `#02` was "proven" by the green
// `#01` beside it.
func mustScenarioRE() *regexp.Regexp {
	return regexp.MustCompile(codeCore() + `(?:#\d{2})?\b`)
}

// ScenarioRoot is the rule a scenario code belongs to: `RDCH-B02#02` → `RDCH-B02`. A code
// without a variant is its own root.
func ScenarioRoot(code string) string {
	if i := strings.IndexByte(code, '#'); i > 0 {
		return code[:i]
	}
	return code
}

// RulesProven says which rules the proven scenario codes prove, given the scenarios the
// features declare (`declared`, with their variants). A rule with declared scenarios is
// proven only when EVERY one of them is — the rule's own code when a scenario carries no
// variant, each `#NN` when they do. A rule no feature declares a scenario of is proven by
// its own code or any variant of it.
func RulesProven(proven, declared []string) map[string]bool {
	got := make(map[string]bool, len(proven))
	for _, c := range proven {
		got[c] = true
	}
	byRule := map[string][]string{}
	for _, d := range declared {
		r := ScenarioRoot(d)
		byRule[r] = append(byRule[r], d)
	}
	out := map[string]bool{}
	for _, c := range proven {
		if r := ScenarioRoot(c); len(byRule[r]) == 0 {
			out[r] = true
		}
	}
	for r, scenarios := range byRule {
		all := true
		for _, sc := range scenarios {
			if !got[sc] {
				all = false
				break
			}
		}
		if all {
			out[r] = true
		}
	}
	return out
}

// UnprovenScenarios are the declared scenarios of the rule that no proven code proves, in
// the order declared.
func UnprovenScenarios(rule string, proven, declared []string) []string {
	got := make(map[string]bool, len(proven))
	for _, c := range proven {
		got[c] = true
	}
	var out []string
	for _, d := range declared {
		if ScenarioRoot(d) == rule && !got[d] {
			out = append(out, d)
		}
	}
	return out
}

// featureTagRE is a scenario tag of a feature: `@RDCH-B02`, `@RDCH-B02#02`.
func featureTagRE() *regexp.Regexp {
	return regexp.MustCompile(`@(` + strings.TrimPrefix(codeCore(), `\b`) + `(?:#\d{2})?)\b`)
}

// FeatureScenarios are the scenario codes a feature's tags declare, with their variants,
// each once, in the order they appear.
func FeatureScenarios(text string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range featureTagRE().FindAllStringSubmatch(text, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

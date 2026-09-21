package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- a FALHA declarada, tratada e registrada ---
//
// A spec já catalogava como a unidade falha, com letra própria (`-E`) e uma tabela de
// condição e efeito. E nada confrontava isso: as regras atravessavam o pipeline inteiro,
// e ninguém perguntava se eram tratadas, se logavam, ou se aconteciam.
//
// Estes três gates fecham a primeira das três camadas do conceito — a que é CONFRONTO
// ESTÁTICO puro, e por isso não depende de produção nem de formato de log nenhum:
//
//	failure-handled    a falha declarada tem caminho que a trata?
//	failure-logged     esse caminho REGISTRA a ocorrência?
//	failure-declared   o tratamento que existe responde a alguma falha declarada?
//
// O terceiro é o inverso do primeiro, e pega o caso mais comum: alguém escreveu uma defesa
// e nunca declarou o que ela previne.
//
// TRATAR NÃO É "TER UM CATCH". Cravar `catch` seria cravar a sintaxe de uma família de
// linguagens — `if (x == null) { return recusa() }` trata tanto quanto, e um `match` em
// Rust também. O que as formas têm em comum não é a aparência, é o EFEITO: a falha vira
// parte do fluxo e a aplicação segue. Quem sabe reconhecer a forma no dialeto local é o
// projeto, em `handle_patterns` — o mesmo mecanismo do `guard_patterns`, que já existia.

// failureRuleRE acha uma regra de FALHA (`-E`) nas três formas catalogadas.
var failureRuleRE = regexp.MustCompile(`(?m)(?:^#{1,6}\s+|^\s*\|\s*` + "`?" + `|^\s*-\s+\*\*)([A-Z0-9]{3,6}-E[0-9]{2})`)

// resilientRE — a falha COMPREENDIDA e absorvida pelo fluxo.
//
// É uma afirmação diferente das duas que já existiam, e por isso tem marcador próprio:
//
//	@no-<coisa>: <razão>   "não vai ter"        — dispensa permanente
//	@TBD: <razão>          "ainda não tem"      — dívida, continua aparecendo
//	@resilient: <razão>    "acontece, sei por quê, e está tratada"
//
// Não é dispensa (a falha é real e acontece) nem dívida (não há nada pendente): é
// conhecimento adquirido. A razão obrigatória é o que a distingue de calar um alerta — ela
// responde a pergunta que alguém vai fazer em seis meses, "por que ignoramos isso?".
var resilientRE = regexp.MustCompile("(?i)(^|[^`])@resilient[^\\S\\n]*:[^\\S\\n]*\\S+")

// declaredFailures lê as regras `-E` da spec, separando as marcadas como resilientes.
func declaredFailures(content string) (all, resilient []string) {
	seen := map[string]bool{}
	for _, line := range strings.Split(content, "\n") {
		for _, m := range failureRuleRE.FindAllStringSubmatch(line, -1) {
			if seen[m[1]] {
				continue
			}
			seen[m[1]] = true
			all = append(all, m[1])
			if resilientRE.MatchString(line) {
				resilient = append(resilient, m[1])
			}
		}
	}
	return all, resilient
}

// governedCode reads the non-comment content of the code this spec specifies.
func governedCode(n mapx.Node, root string, g *mapx.Graph) (string, bool) {
	if g == nil {
		return "", false
	}
	var body strings.Builder
	found := false
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.To))
		if err != nil {
			continue
		}
		found = true
		body.WriteString(stripLineComments(string(b)))
		body.WriteString("\n")
	}
	return body.String(), found
}

// anyMatch says whether any of the declared patterns matches the code.
func anyMatch(d config.Dialect, patterns []string, code string) bool {
	for _, p := range patterns {
		if re := d.Compile(p); re != nil && re.MatchString(code) {
			return true
		}
	}
	return false
}

// --- a falha declarada tem caminho que a trata? ---
func checkFailureHandled(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.failure_handled.skip_not_spec")
	}
	all, resilient := declaredFailures(content)
	if len(all) == 0 {
		return Skip, i18n.T("gate.failure_handled.skip_none")
	}
	d := cfg.DialectFor()
	if len(d.HandlePatterns) == 0 {
		// Pendente e não Pass: sem saber reconhecer um tratamento, o gate não verificou
		// nada — e um ✓ aqui seria carimbar o que nunca foi medido.
		return Pending, i18n.T("gate.failure_handled.skip_no_dialect")
	}
	code, ok := governedCode(n, root, g)
	if !ok {
		return Pending, i18n.T("gate.failure_handled.pending_no_code")
	}

	// A falha RESILIENTE sai da cobrança: ela foi compreendida, está absorvida pelo fluxo,
	// e a razão está escrita ao lado dela.
	isResilient := map[string]bool{}
	for _, r := range resilient {
		isResilient[r] = true
	}

	// A régua é do CONJUNTO, e não de cada regra: o código não cita o código da falha
	// (`CRED-E01` não aparece num `if`), então não há como amarrar uma regra a um caminho
	// específico. O que dá para afirmar é se a unidade tem tratamento nenhum enquanto
	// declara falhas — e é justamente esse o caso que interessa.
	if anyMatch(d, d.HandlePatterns, code) {
		return Pass, ""
	}
	var charged []string
	for _, f := range all {
		if !isResilient[f] {
			charged = append(charged, f)
		}
	}
	if len(charged) == 0 {
		return Pass, ""
	}
	sort.Strings(charged)
	return Fail, fmt.Sprintf(i18n.T("gate.failure_handled.untreated"), len(charged), strings.Join(charged, ", "))
}

// --- o tratamento REGISTRA a ocorrência? ---
func checkFailureLogged(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.failure_handled.skip_not_spec")
	}
	all, resilient := declaredFailures(content)
	if len(all) == 0 {
		return Skip, i18n.T("gate.failure_handled.skip_none")
	}
	d := cfg.DialectFor()
	if len(d.HandlePatterns) == 0 || len(d.LogPatterns) == 0 {
		return Pending, i18n.T("gate.failure_handled.skip_no_dialect")
	}
	code, ok := governedCode(n, root, g)
	if !ok {
		return Pending, i18n.T("gate.failure_handled.pending_no_code")
	}
	// Sem tratamento nenhum não há o que cobrar aqui: quem cobra é o gate irmão, e dois
	// gates acusando o mesmo defeito viram ruído.
	if !anyMatch(d, d.HandlePatterns, code) {
		return Skip, ""
	}
	if anyMatch(d, d.LogPatterns, code) {
		return Pass, ""
	}
	isResilient := map[string]bool{}
	for _, r := range resilient {
		isResilient[r] = true
	}
	var charged []string
	for _, f := range all {
		if !isResilient[f] {
			charged = append(charged, f)
		}
	}
	if len(charged) == 0 {
		return Pass, ""
	}
	sort.Strings(charged)
	return Fail, fmt.Sprintf(i18n.T("gate.failure_logged.unlogged"), len(charged), strings.Join(charged, ", "))
}

// --- o tratamento que existe responde a alguma falha declarada? ---
//
// O inverso do `failure-handled`, e o caso mais comum: alguém escreveu uma defesa e nunca
// declarou o que ela previne. A defesa pode estar certíssima — o que falta é a spec dizer
// a qual falha ela responde, para quem ler depois saber se ainda vale.
func checkFailureDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.failure_handled.skip_not_spec")
	}
	d := cfg.DialectFor()
	if len(d.HandlePatterns) == 0 {
		return Pending, i18n.T("gate.failure_handled.skip_no_dialect")
	}
	code, ok := governedCode(n, root, g)
	if !ok {
		return Pending, i18n.T("gate.failure_handled.pending_no_code")
	}
	hits := 0
	for _, p := range d.HandlePatterns {
		if re := d.Compile(p); re != nil {
			hits += len(re.FindAllString(code, -1))
		}
	}
	if hits == 0 {
		return Skip, i18n.T("gate.failure_declared.skip_none")
	}
	all, _ := declaredFailures(content)
	if len(all) > 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(i18n.T("gate.failure_declared.undeclared"), hits, len(all))
}

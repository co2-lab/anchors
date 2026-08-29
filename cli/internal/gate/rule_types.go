package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// rule-types-consistentes: o VOCABULÁRIO de tipos de regra é extensível — cada letra do
// código (`{CODE}-<letra><NN>`) é a inicial do termo que nomeia a seção da spec —, mas
// precisa ser DECLARADO (`rule_types` no anchors.yaml). Uma letra não declarada é
// invisível para a rastreabilidade: a regra aparece na spec, na `.feature` e no teste,
// e mesmo assim o `feature-test-match` não a enxerga (o regex de código não a casa).
// Esse é o pior tipo de furo — parece coberto e não está.
//
// Confronta três coisas (nesta ordem de gravidade):
//  1. LETRA NÃO DECLARADA — a spec usa `{CODE}-P01` e `P` não consta no vocabulário.
//  2. SEÇÃO SEM LETRA — a spec cataloga regras sob um título que nenhuma letra reivindica.
//  3. CONFLITO DE LETRA — duas seções DIFERENTES reivindicam a MESMA letra (erro de
//     configuração; detectado no vocabulário, não no arquivo).
//
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func ruleCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`\b[A-Z0-9]` + config.CodeLengthPattern() + `-([A-Z])\d{2}\b`)
}

// letrasCanonicas confronta a spec contra `config.DefaultRuleLetters` — o que o gate faz
// quando o projeto não declarou vocabulário próprio.
//
// Uma letra fora das canônicas é invisível para a rastreabilidade: o `feature-test-match`
// não a enxerga, mesmo com feature e teste escritos. O achado é o mesmo do modo declarado;
// muda só o conserto sugerido, porque aqui o projeto ainda não tem onde declarar.
func letrasCanonicas(content string) (Verdict, string) {
	canonicas := map[string]bool{}
	for _, l := range config.DefaultRuleLetters {
		canonicas[string(l)] = true
	}
	var fora []string
	visto := map[string]bool{}
	for _, m := range ruleCodeRE().FindAllStringSubmatch(content, -1) {
		l := m[1]
		if canonicas[l] || visto[l] {
			continue
		}
		visto[l] = true
		fora = append(fora, l)
	}
	if len(fora) == 0 {
		return Pass, ""
	}
	sort.Strings(fora)
	return Fail, fmt.Sprintf("usa tipo(s) de regra fora do vocabulário canônico (%s): %s. "+
		"Uma letra que o engine não reconhece é invisível para a rastreabilidade — o "+
		"feature-test-match não a enxerga, mesmo com feature e teste escritos. Declare o "+
		"vocabulário do projeto em `rule_types` no anchors.yaml (letra + termo + seções), "+
		"ou use uma letra canônica.",
		config.DefaultRuleLetters, strings.Join(fora, ", "))
}

// specSectionRE captura os títulos de seção (## / ### / ####) de uma spec.
var specSectionRE = regexp.MustCompile(`(?m)^#{2,4}\s+(.+?)\s*$`)

func checkRuleTypes(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	// SEM `rule_types` declarado o gate confronta as LETRAS CANÔNICAS.
	//
	// Ele fazia `Skip` aqui, e o efeito era que um gate CANÔNICO — semeado pelo `init` em
	// todo projeto — nunca media nada: aparecia na tabela do `check`, ocupava linha, e
	// reportava indeterminado para sempre. É o silêncio que o Anchors combate em todo
	// resto: um gate declarado que não confronta nada dá a impressão de defesa que não
	// existe.
	//
	// Só a verificação (1) — letra fora do vocabulário — tem o que fazer nesse modo. As
	// outras duas dependem de SEÇÕES e TERMOS, que só existem quando o projeto declara o
	// vocabulário; cobrá-las contra um vocabulário implícito seria inventar regra que
	// ninguém escreveu.
	if cfg == nil || len(cfg.RuleTypes) == 0 {
		return letrasCanonicas(content)
	}

	// (3) CONFLITO no próprio vocabulário: mesma letra reivindicada por seções distintas.
	if msg := conflictingLetters(cfg.RuleTypes); msg != "" {
		return Fail, msg
	}

	declared := map[string]bool{}
	for _, rt := range cfg.RuleTypes {
		declared[strings.ToUpper(strings.TrimSpace(rt.Letter))] = true
	}

	// (1) LETRA NÃO DECLARADA usada neste arquivo.
	var undeclared []string
	seen := map[string]bool{}
	for _, m := range ruleCodeRE().FindAllStringSubmatch(content, -1) {
		l := m[1]
		if declared[l] || seen[l] {
			continue
		}
		seen[l] = true
		undeclared = append(undeclared, l)
	}
	if len(undeclared) > 0 {
		sort.Strings(undeclared)
		return Fail, fmt.Sprintf("usa tipo(s) de regra NÃO declarado(s): %s. Uma letra fora do "+
			"vocabulário é invisível para a rastreabilidade (o feature-test-match não a "+
			"enxerga, mesmo com feature e teste escritos). Declare em `rule_types` no "+
			"anchors.yaml (letra + termo + seções) ou use uma letra existente.",
			strings.Join(undeclared, ", "))
	}

	// (2) SEÇÃO QUE CATALOGA REGRAS mas cujo título nenhuma letra reivindica.
	if msg := unclaimedSections(content, cfg.RuleTypes); msg != "" {
		return Fail, msg
	}

	// (4) SEÇÃO DECLARADA `requires_code` mas PREENCHIDA SEM CÓDIGO.
	//
	// O inverso do (2): ali a seção tem código sob título não declarado; aqui ela
	// tem CONTEÚDO e nenhum código, num título que o projeto declarou como
	// catalogador de regra. É o furo por onde o cenário fica sem âncora — medido num
	// projeto real: 48 specs com "Eventos / Callbacks" preenchida e sem um único
	// código, e os cenários que provavam esses eventos emprestaram o código do
	// estado vizinho (um `-S` regendo comportamento).
	if msg := secoesSemCodigo(content, cfg.RuleTypes); msg != "" {
		return Pending, msg
	}
	return Pass, ""
}

// conflictingLetters acha a mesma letra reivindicada por SEÇÕES diferentes. Repetir a
// letra com as mesmas seções é só redundância; o conflito é a disputa.
func conflictingLetters(types []config.RuleType) string {
	byLetter := map[string][]string{}
	for _, rt := range types {
		l := strings.ToUpper(strings.TrimSpace(rt.Letter))
		byLetter[l] = append(byLetter[l], rt.Term)
	}
	var bad []string
	for l, terms := range byLetter {
		if len(terms) > 1 {
			bad = append(bad, fmt.Sprintf("%s (%s)", l, strings.Join(terms, " vs ")))
		}
	}
	if len(bad) == 0 {
		return ""
	}
	sort.Strings(bad)
	return "CONFLITO no vocabulário: letra(s) reivindicada(s) por mais de um termo — " +
		strings.Join(bad, "; ") + ". Cada letra pertence a UM termo/seção."
}

// definesRuleRE casa uma linha que DEFINE uma regra: o código é a primeira coisa da
// linha (cabeçalho `### CODEX-B01:`, bullet `- **CODEX-B01**`) ou a primeira célula de
// uma linha de tabela (`| \`CODEX-B01\` | … |`). Uma seção que apenas CITA códigos de
// outras seções (ex.: "Test IDs (Maestro)", que referencia `BUTOX-A01` na coluna "Usado
// em") NÃO cataloga regra — e portanto não precisa reivindicar letra.
var definesRuleRE = regexp.MustCompile(
	"(?m)^\\s*(?:#{2,6}\\s+|[-*]\\s+\\**|\\|\\s*)`?\\*{0,2}[A-Z0-9]" + config.CodeLengthPattern() + "-[A-Z]\\d{2}")

// unclaimedSections acha seções que catalogam regras (DEFINEM ao menos um código) sob um
// título que nenhuma entrada do vocabulário declara. Seções sem código, ou que só citam
// códigos alheios, são ignoradas — o gate cobra rastreabilidade, não formato.
func unclaimedSections(content string, types []config.RuleType) string {
	claimed := map[string]bool{}
	for _, rt := range types {
		for _, s := range rt.Sections {
			claimed[normalizeSection(s)] = true
		}
	}
	lines := strings.Split(content, "\n")
	cur := ""
	var orphan []string
	seen := map[string]bool{}
	for _, line := range lines {
		if m := specSectionRE.FindStringSubmatch(line); m != nil {
			cur = m[1]
			continue
		}
		if cur == "" || !definesRuleRE.MatchString(line) {
			continue
		}
		key := normalizeSection(cur)
		if claimed[key] || seen[key] {
			continue
		}
		// Um título que É o próprio código (ex.: "### OCHIX-P01: Limite…") não é uma
		// seção-categoria; é o cabeçalho da regra. Não cobramos esses.
		if ruleCodeRE().MatchString(cur) {
			continue
		}
		// Idem para o cabeçalho de ENUMERADOR de Data State (`### \`campo\` — descrição`):
		// é a declaração de um campo, não uma categoria de regra.
		if strings.HasPrefix(strings.TrimSpace(cur), "`") {
			continue
		}
		seen[key] = true
		orphan = append(orphan, cur)
	}
	if len(orphan) == 0 {
		return ""
	}
	sort.Strings(orphan)
	return fmt.Sprintf("seção(ões) que catalogam regras sem letra declarada no vocabulário: %q. "+
		"Declare a seção em `rule_types` (na entrada da letra que ela usa) no anchors.yaml.",
		strings.Join(orphan, ", "))
}

func normalizeSection(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// linhaDeTabelaRE: linha de corpo de tabela markdown (não o cabeçalho nem o
// separador). Duas barras bastam para ser célula; o separador é `| --- |`.
var linhaDeTabelaRE = regexp.MustCompile(`^\s*\|[^|]*\|`)
var separadorTabelaRE = regexp.MustCompile(`^\s*\|[\s:|-]+\|?\s*$`)

// secoesSemCodigo acha seção declarada `requires_code` que tem tabela preenchida e
// nenhum código de regra.
//
// Só olha TABELA: uma seção pode ter prosa explicativa sem catalogar nada. O que
// caracteriza catálogo é a linha de tabela — e é lá que o código deveria estar, na
// primeira célula, como as outras seções fazem.
func secoesSemCodigo(content string, types []config.RuleType) string {
	exige := map[string]string{} // título normalizado → letra
	for _, rt := range types {
		for _, s := range rt.RequiresCode {
			exige[normalizeSection(s)] = rt.Letter
		}
	}
	if len(exige) == 0 {
		return ""
	}

	lines := strings.Split(content, "\n")
	cur := ""
	linhas := 0
	temCodigo := false
	var achados []string
	fecha := func() {
		if cur == "" {
			return
		}
		if letra, ok := exige[normalizeSection(cur)]; ok && linhas > 0 && !temCodigo {
			achados = append(achados, fmt.Sprintf("%q (letra %s, %d linha(s))", cur, letra, linhas))
		}
		linhas, temCodigo = 0, false
	}
	for _, line := range lines {
		if m := specSectionRE.FindStringSubmatch(line); m != nil {
			fecha()
			cur = m[1]
			continue
		}
		if cur == "" {
			continue
		}
		if ruleCodeRE().MatchString(line) {
			temCodigo = true
		}
		if linhaDeTabelaRE.MatchString(line) && !separadorTabelaRE.MatchString(line) &&
			!strings.Contains(strings.ToLower(line), "---") {
			linhas++
		}
	}
	fecha()

	if len(achados) == 0 {
		return ""
	}
	sort.Strings(achados)
	return fmt.Sprintf("%d seção(ões) que o projeto declarou como catalogadoras de regra "+
		"(`requires_code`) estão preenchidas SEM código: %s. Cada linha dessas tabelas afirma "+
		"algo verificável, e sem código o cenário que a prova não tem o que citar — na prática, "+
		"ele empresta o código de outra seção e passa a reger o que não é dele. "+
		"Acrescente uma coluna `Regra` com o código, ou tire a seção de `requires_code` se ela "+
		"apenas enumera valores",
		len(achados), strings.Join(achados, "; "))
}

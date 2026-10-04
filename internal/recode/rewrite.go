// @anchors
//   code: RWINR
//   ref: RCRWR

// Package recode renomeia um CÓDIGO de identidade (ex.: TCDTX → TCTXX) e o propaga por
// todas as superfícies textuais onde ele aparece: o header @anchors (code:/ref:), os
// scenario-codes derivados (CODEX-B01, CODEX-S02, CODE-DS-*, CODEX-VR…) e as menções nuas
// do código em referências cruzadas de outras unidades.
//
// É a base da estabilidade REVERSÍVEL da identidade: hoje o código é "estável" só
// porque não há como renomeá-lo com segurança. Este motor é puro (texto→texto); o
// planner e o comando lidam com I/O, escopo e o mapa.
package recode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// codeRE valida um código de identidade: alfanumérico maiúsculo, no comprimento que o
// projeto declara (`code_lengths`). Compilado por CHAMADA porque a config é carregada
// depois dos globais — um `var` congelaria o default.
func codeRE() *regexp.Regexp {
	return regexp.MustCompile(`^[A-Z0-9]` + config.CodeLengthPattern() + `$`)
}

// ValidCode diz se s é um código de identidade bem-formado.
func ValidCode(s string) bool { return codeRE().MatchString(s) }

// Occurrence é uma ocorrência do código num texto, classificada por superfície.
type Occurrence struct {
	Kind  string // "header-code" | "header-ref" | "scenario-code" | "bare-ref"
	Match string // o texto exato casado (ex.: "code: TCDT", "TCDTX-S02", "TCDT")
	Line  int    // 1-indexed
}

// buildPatterns monta, para um código OLD, as três famílias de padrão que o recode
// reconhece. Todas ancoram o OLD por FRONTEIRA para nunca casar um código maior
// (TCDT não casa dentro de TCDTX) nem o miolo de um scenario-code de outro código.
func buildPatterns(old string) (headerRE, scenarioRE, bareRE *regexp.Regexp) {
	o := regexp.QuoteMeta(old)
	// header: `code:`/`ref:` seguido do código; ref pode ser LISTA (o replace troca só
	// as ocorrências do OLD, preservando os demais códigos da lista).
	headerRE = regexp.MustCompile(`(?m)((?:code|ref):[^\n]*?\b)` + o + `\b`)
	// scenario-code: OLD- seguido do sufixo (1-2 letras + dígitos + sufixo textual/alnum,
	// ou DS-*, ou VR). Mais abrangente que a regex do scan (cobre FP/RA/RC/sufixo-b).
	scenarioRE = regexp.MustCompile(`\b` + o + `(-[A-Za-z0-9][A-Za-z0-9-]*)`)
	// menção nua do código (referência cruzada), fora de header e scenario-code.
	// Go (RE2) não tem lookahead, então casamos OLD + o caractere-limite seguinte
	// (não-alfanumérico e não '-'), num grupo, para recolocá-lo. `(?:$)` cobre o fim.
	// Os scenario-codes (OLD-…) já foram trocados antes, então um OLD-… não sobra aqui.
	//
	// `_` is no boundary: `GOAL_TABLE_NAME` is an identifier that holds the word, not a
	// citation of the code `GOAL` — rewritten, it broke an env var its consumers still read
	// by the old name (reported from MIF).
	bareRE = regexp.MustCompile(`\b` + o + `([^A-Za-z0-9_-]|$)`)
	return
}

// Rewrite reescreve todas as ocorrências do código OLD para NEW num texto, em ordem
// segura: primeiro os scenario-codes (OLD-…), depois o header (code:/ref:), por fim as
// menções nuas — sem que uma etapa corrompa a outra. Devolve o texto novo e quantas
// substituições fez.
func Rewrite(content, old, new string) (string, int) {
	headerRE, scenarioRE, bareRE := buildPatterns(old)
	n := 0

	// 1. scenario-codes: OLD-XXX → NEW-XXX (preserva o sufixo).
	content = scenarioRE.ReplaceAllStringFunc(content, func(m string) string {
		n++
		suffix := m[len(old):] // "-S02", "-DS-x", "-VR"…
		return new + suffix
	})

	// 2. header code:/ref: (inclui listas — só o token OLD é trocado).
	content = headerRE.ReplaceAllStringFunc(content, func(m string) string {
		n++
		// o grupo 1 é o prefixo "code: …"; o OLD está no fim do match.
		return m[:len(m)-len(old)] + new
	})

	// 3. menções nuas (refs cruzadas). Já trocamos scenario-codes, então um OLD que
	// sobrou aqui é uma menção pura do código. O grupo 1 é o caractere-limite seguinte
	// (ou vazio, no fim) — recolocado intacto.
	content = bareRE.ReplaceAllStringFunc(content, func(m string) string {
		n++
		boundary := m[len(old):] // o caractere-limite capturado
		return new + boundary
	})

	return content, n
}

// Find lista as ocorrências do código OLD num texto (para o dry-run), classificadas.
// Espelha a ordem/gramática do Rewrite.
func Find(content, old string) []Occurrence {
	headerRE, scenarioRE, bareRE := buildPatterns(old)
	var occ []Occurrence

	add := func(kind string, idx [][]int) {
		for _, loc := range idx {
			occ = append(occ, Occurrence{
				Kind:  kind,
				Match: content[loc[0]:loc[1]],
				Line:  lineOf(content, loc[0]),
			})
		}
	}
	// scenario-codes e header podem se sobrepor a bare; classificamos na mesma ordem do
	// Rewrite, e cada ocorrência já contada é MASCARADA no texto-sombra para não recontar.
	//
	// Mascarar, e não remover: o shadow guarda o comprimento do original, então as
	// posições dos kinds seguintes continuam valendo no `content`. Removendo, o shadow
	// encurtava, e o `Match` e a `Line` de um header ou de uma menção depois de um
	// scenario-code saíam de outro trecho do texto. Só o código é mascarado, como o
	// Rewrite só troca o código: o sufixo e o prefixo `code:` ficam, e a máscara é
	// alfanumérica como o código novo seria — as fronteiras que os padrões seguintes
	// leem são as mesmas do Rewrite.
	shadow := content
	scen := scenarioRE.FindAllStringIndex(shadow, -1)
	add("scenario-code", scen)
	for _, loc := range scen {
		shadow = maskCode(shadow, loc[0], loc[0]+len(old))
	}

	hdr := headerRE.FindAllStringIndex(shadow, -1)
	add("header", hdr)
	for _, loc := range hdr {
		shadow = maskCode(shadow, loc[1]-len(old), loc[1])
	}

	bare := bareRE.FindAllStringIndex(shadow, -1)
	add("bare-ref", bare)

	return occ
}

// maskCode overwrites s[from:to] with lower-case letters, keeping the length. Lower case
// never forms a code, and a letter keeps the boundaries a replaced code would.
func maskCode(s string, from, to int) string {
	return s[:from] + strings.Repeat("x", to-from) + s[to:]
}

// lineOf gives the 1-indexed line of a byte offset. The offset never passes the end of s:
// it comes from a match in s or in the shadow text, which has the same length as s.
func lineOf(s string, byteIdx int) int {
	return strings.Count(s[:byteIdx], "\n") + 1
}

// String — util p/ debug.
func (o Occurrence) String() string {
	return fmt.Sprintf("L%d %s: %s", o.Line, o.Kind, o.Match)
}

// RewriteRuleCodes rewrites only the rule and scenario codes of OLD — `OLD-B08`,
// `OLD-S06#02`, `OLD-VR-S01`, `OLD-CT` — leaving the bare code alone. It is the rewrite
// for a file the project does not govern (a runner script, a lint config, an env
// example): there a bare `DATA` is as likely an ordinary word or `DATA_URL` as a code,
// and only the code's rule shape says for certain that it is one.
func RewriteRuleCodes(content, old, new string) (string, int) {
	n := 0
	out := ruleCodeRE(old).ReplaceAllStringFunc(content, func(m string) string {
		n++
		return new + m[len(old):]
	})
	return out, n
}

// ruleCodeRE matches the rule and scenario codes of OLD, with their suffix.
func ruleCodeRE(old string) *regexp.Regexp {
	return regexp.MustCompile(`\b` + regexp.QuoteMeta(old) + `(-(?:VR(?:-[A-Z]{1,2}\d{2})?|CT|[A-Z]{1,2}\d{2,}[a-z]?)(?:#\d+)?)\b`)
}

// RewriteCited rewrites OLD only where the text cites it AS A CODE: the rule and scenario
// codes (`OLD-B08`, `OLD-VR-S01`, `OLD-CT`), the header fields that hold codes (`code:`,
// `ref:`/`refs:`, `dep:`, `needs:`), the code in backticks (“ `OLD` “) and a Gherkin tag
// (`@OLD`). A bare word is left alone: a code is short and upper-case, and the same letters
// are a word of the language ("CNPJ inválido", "na HOME") or part of an identifier
// (`SEAT_PRICE`) — rewriting those broke user-facing text and an env var the code reads by
// name (reported from MIF). What is left is listed by CitedSet.Rewrite, for a person to judge.
func RewriteCited(content, old, new string) (string, int) {
	out, counts, _ := NewCitedSet(map[string]string{old: new}).Rewrite(content)
	return out, counts[old]
}

// CitedSet rewrites many codes at once where each is cited as a code (see RewriteCited),
// in one pass over the text: a migration renaming hundreds of codes reads every file —
// and the map, the largest — once, not once per code.
type CitedSet struct {
	to map[string]string
}

// NewCitedSet is the rewrite of each old code to its new one.
func NewCitedSet(to map[string]string) *CitedSet { return &CitedSet{to: to} }

var (
	citedWordRE   = regexp.MustCompile(`[A-Za-z0-9_]+`)
	citedSuffixRE = regexp.MustCompile(`^-(?:VR(?:-[A-Z]{1,2}\d{2})?|CT|[A-Z]{1,2}\d{2,}[a-z]?)(?:#\d+)?\b`)
	citedFieldRE  = regexp.MustCompile(`(?:code|refs?|dep|needs)\s*:`)
)

// RewriteRuleCodes rewrites only the rule and scenario codes of the set's codes (see
// RewriteRuleCodes), in one pass, and says how many of each.
func (c *CitedSet) RewriteRuleCodes(content string) (string, map[string]int) {
	counts := map[string]int{}
	var b strings.Builder
	last := 0
	for _, loc := range citedWordRE.FindAllStringIndex(content, -1) {
		s, e := loc[0], loc[1]
		nw, ok := c.to[content[s:e]]
		if !ok || !citedSuffixRE.MatchString(content[e:]) {
			continue
		}
		b.WriteString(content[last:s])
		b.WriteString(nw)
		last = e
		counts[content[s:e]]++
	}
	b.WriteString(content[last:])
	return b.String(), counts
}

// Rewrite returns the text with every cited old code rewritten, how many of each it
// rewrote, and the bare mentions of each it left — a whole word with the code's letters,
// not followed by a hyphen, outside the cited forms.
func (c *CitedSet) Rewrite(content string) (string, map[string]int, map[string][]BareMention) {
	counts := map[string]int{}
	var bare map[string][]BareMention
	var b strings.Builder
	last, line, lineAt := 0, 1, 0
	for _, loc := range citedWordRE.FindAllStringIndex(content, -1) {
		s, e := loc[0], loc[1]
		w := content[s:e]
		nw, ok := c.to[w]
		if !ok {
			continue
		}
		cited := citedSuffixRE.MatchString(content[e:]) ||
			(s > 0 && e < len(content) && content[s-1] == '`' && content[e] == '`') ||
			(s > 0 && content[s-1] == '@' && (s == 1 || strings.ContainsRune(" \t\n([,", rune(content[s-2]))) &&
				(e == len(content) || content[e] != '-'))
		if !cited {
			ls := strings.LastIndexByte(content[:s], '\n') + 1
			cited = citedFieldRE.MatchString(content[ls:s])
		}
		if cited {
			b.WriteString(content[last:s])
			b.WriteString(nw)
			last = e
			counts[w]++
			continue
		}
		if e < len(content) && content[e] == '-' {
			continue // a testID or a word joined by a hyphen, not a mention
		}
		line += strings.Count(content[lineAt:s], "\n")
		lineAt = s
		ls := strings.LastIndexByte(content[:s], '\n') + 1
		le := strings.IndexByte(content[s:], '\n')
		if le < 0 {
			le = len(content) - s
		}
		if bare == nil {
			bare = map[string][]BareMention{}
		}
		if ms := bare[w]; len(ms) == 0 || ms[len(ms)-1].Line != line {
			bare[w] = append(ms, BareMention{Line: line, Text: strings.TrimSpace(content[ls : s+le])})
		}
	}
	b.WriteString(content[last:])
	return b.String(), counts, bare
}

// BareMention is an occurrence of a code that a cited rewrite leaves alone, for review.
type BareMention struct {
	Line int
	Text string // the line, trimmed
}

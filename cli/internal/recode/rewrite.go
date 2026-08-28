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
	bareRE = regexp.MustCompile(`\b` + o + `([^A-Za-z0-9-]|$)`)
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
	// Rewrite e removemos do texto-sombra para não recontar. NOTA: o shadow encurta ao
	// remover, então Occurrence.Line dos kinds posteriores é aproximado — a CONTAGEM por
	// kind (o que o dry-run usa) é exata; a linha é só para debug.
	shadow := content
	scen := scenarioRE.FindAllStringIndex(shadow, -1)
	add("scenario-code", scen)
	shadow = scenarioRE.ReplaceAllString(shadow, "")

	hdr := headerRE.FindAllStringIndex(shadow, -1)
	add("header", hdr)
	shadow = headerRE.ReplaceAllString(shadow, "")

	bare := bareRE.FindAllStringIndex(shadow, -1)
	add("bare-ref", bare)

	return occ
}

func lineOf(s string, byteIdx int) int {
	line := 1
	for i := 0; i < byteIdx && i < len(s); i++ {
		if s[i] == '\n' {
			line++
		}
	}
	return line
}

// String — util p/ debug.
func (o Occurrence) String() string {
	return fmt.Sprintf("L%d %s: %s", o.Line, o.Kind, o.Match)
}

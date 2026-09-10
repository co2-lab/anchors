package recode

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// TestIDPrefix deriva o prefixo de testID de um código, segundo a convenção declarada.
// "lower" → o código em minúsculo (TCDT → "tcdt"). Vazio/desconhecido → "".
// A convenção IDEAL do projeto: o testID é a projeção do código na borda observável.
func TestIDPrefix(code, convention string) string {
	switch convention {
	case "lower":
		return strings.ToLower(code)
	default:
		return ""
	}
}

// testIDToken casa um prefixo de testID usado como token: dentro de aspas/backtick
// seguido do prefixo e de '-' (ex.: testID="tcdt-amount", 'tcdt-row'). Devolve os
// índices para troca. Usa fronteira à esquerda (não-alfanumérico) para não casar meio
// de palavra.
func testIDTokenRE(prefix string) *regexp.Regexp {
	// prefixo seguido de '-' e de um caractere de identificador; à esquerda, um
	// delimitador de string (aspas simples/dupla/backtick) — o testID vive em literais.
	return regexp.MustCompile(`(["'` + "`" + `])` + regexp.QuoteMeta(prefix) + `(-[A-Za-z0-9])`)
}

// RewriteTestIDs troca o prefixo de testID oldPrefix→newPrefix num texto, só dentro de
// literais de string. Devolve o texto novo e a contagem. Se oldPrefix == "" não faz nada.
func RewriteTestIDs(content, oldPrefix, newPrefix string) (string, int) {
	if oldPrefix == "" || oldPrefix == newPrefix {
		return content, 0
	}
	re := testIDTokenRE(oldPrefix)
	n := 0
	out := re.ReplaceAllStringFunc(content, func(m string) string {
		n++
		// m = <quote><oldPrefix><-x>; recompõe com o novo prefixo.
		quote := m[:1]
		rest := m[1+len(oldPrefix):]
		return quote + newPrefix + rest
	})
	return out, n
}

// CountTestIDPrefix conta ocorrências de um prefixo de testID (para detectar se o
// esperado existe, ou se há um candidato divergente — o aviso do legado).
func CountTestIDPrefix(content, prefix string) int {
	if prefix == "" {
		return 0
	}
	return len(testIDTokenRE(prefix).FindAllStringIndex(content, -1))
}

// anyTestIDRE casa QUALQUER testID (testID="algo") — para detectar que o arquivo TEM
// testIDs mesmo que o prefixo esperado não bata (o aviso de legado não conhece o
// prefixo antigo; só sabe que o esperado sumiu e há testIDs).
var anyTestIDRE = regexp.MustCompile(`\btestID\s*=\s*["'` + "`" + `][a-zA-Z]`)

// CountAnyTestID conta quantos atributos testID= existem no texto.
func CountAnyTestID(content string) int {
	return len(anyTestIDRE.FindAllStringIndex(content, -1))
}

// FileMatchesCode diz se um caminho casa algum dos file_patterns para o código dado,
// substituindo {{code}} pelo código. Case-sensitive no código (TCDT ≠ tcdt).
func FileMatchesCode(relPath, code string, patterns []string) bool {
	for _, p := range patterns {
		glob := strings.ReplaceAll(p, "{{code}}", code)
		if ok, _ := doublestar.Match(glob, filepath.ToSlash(relPath)); ok {
			return true
		}
	}
	return false
}

// RenameFilePath devolve o novo caminho de um arquivo cujo NOME contém o código old,
// trocando old→new no basename (por fronteira, para não tocar código maior). Só o
// basename muda; o diretório permanece.
func RenameFilePath(relPath, old, new string) string {
	dir := filepath.Dir(relPath)
	base := filepath.Base(relPath)
	// troca o código no nome por fronteira (delimitadores: início, '.', '-', fim).
	re := regexp.MustCompile(`(^|[.\-])` + regexp.QuoteMeta(old) + `($|[.\-])`)
	nb := re.ReplaceAllString(base, "${1}"+new+"${2}")
	if nb == base {
		return relPath // não continha o código no nome
	}
	// ToSlash porque o que se devolve é um id de nó, não um caminho de disco: no Windows
	// o Join daria "a\b\X.yaml" e o id deixaria de casar com o mapa (gravado com "/").
	return filepath.ToSlash(filepath.Join(dir, nb))
}

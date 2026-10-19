// @anchors
//   code: PTINP
//   ref: CFPCD

package code

import (
	"path/filepath"
)

// GenerateFromPath é a porta de alto nível: dado o CAMINHO de uma unidade, o NOME dela — o
// basename sem o que as layers acrescentam, que o projeto diz (`config.UnitName`) — e os
// códigos já tomados, devolve um código único e o MAIS SEMÂNTICO possível, sem prefixo de
// módulo declarado. A régua:
//
//   - basename GENÉRICO (handler/index/…) → o dir-pai É a identidade: comprime o
//     dir-pai (manage-metadata → MNMT). É DETERMINÍSTICO: independe de ordem, pois o
//     nome real (o dir) é o mesmo sempre.
//   - basename normal → algoritmo canônico (Camada 2), com o desempate cego de colisão
//     (GenerateUnique). NOTA: o desempate cego é determinístico DADO O MAPA, mas quem
//     leva o código "limpo" depende da ordem de criação (o primeiro fica com o
//     canônico, o segundo recebe a variação). Um desempate SIMÉTRICO (ambos
//     prefixados) exigiria reescrever o código do primeiro — uma operação de RENAME
//     (recode) que propaga por toda a unidade. Isso é uma feature própria (planejada à
//     parte), não um efeito colateral da geração. Aqui ficamos no cego determinístico.
//
// Quando a layer TEM code_prefix declarado, o chamador passa o prefixo e NÃO usa esta
// função (Camada 1 tem precedência — é dialeto curado do projeto).
// Whether the name is GENERIC — a role in its folder, not the unit (`handler`, `index`) — is
// the project's to say (`dialect.generic_basenames`); the caller asks it and passes it.
func GenerateFromPath(path, base string, generic bool, taken map[string]bool) string {
	parent := parentDir(path)

	// basename genérico: a pasta é o nome (determinístico).
	if generic && parent != "" {
		return GenerateUnique(parent, taken)
	}

	// basename normal: canônico + desempate cego de colisão.
	return GenerateUnique(base, taken)
}

// parentDir devolve o nome da pasta imediatamente acima do arquivo ("" se na raiz).
func parentDir(path string) string {
	d := filepath.Dir(path)
	if d == "." || d == "/" || d == "" {
		return ""
	}
	return filepath.Base(d)
}

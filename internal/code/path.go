package code

import (
	"path/filepath"
	"strings"
)

// genericBasenames são nomes de ARQUIVO que não distinguem a unidade: o que a
// identifica é a PASTA-PAI (functions/manage-metadata/handler.ts → a unidade é
// "manage-metadata", não "handler"). Nesses, o dir-pai vira o nome.
var genericBasenames = map[string]bool{
	"handler": true, "index": true, "resource": true, "mod": true,
	"main": true, "types": true, "route": true, "routes": true,
	"schema": true, "config": true,
}

// IsGenericBasename diz se um basename (sem extensão) não carrega identidade própria —
// é um papel dentro da pasta, e a pasta é o nome real.
func IsGenericBasename(name string) bool {
	return genericBasenames[strings.ToLower(name)]
}

// GenerateFromPath é a porta de alto nível: dado o CAMINHO de uma unidade e os códigos
// já tomados, devolve um código único e o MAIS SEMÂNTICO possível, sem prefixo de
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
//     (recode) que propaga por toda a trinca. Isso é uma feature própria (planejada à
//     parte), não um efeito colateral da geração. Aqui ficamos no cego determinístico.
//
// Quando a layer TEM code_prefix declarado, o chamador passa o prefixo e NÃO usa esta
// função (Camada 1 tem precedência — é dialeto curado do projeto).
func GenerateFromPath(path string, taken map[string]bool) string {
	base := unitBase(path)
	parent := parentDir(path)

	// basename genérico: a pasta é o nome (determinístico).
	if IsGenericBasename(base) && parent != "" {
		return GenerateUnique(parent, taken)
	}

	// basename normal: canônico + desempate cego de colisão.
	return GenerateUnique(base, taken)
}

// unitBase extrai o nome distintivo do caminho: basename sem extensão nem sufixos de
// artefato (Login.spec.md → Login; metadata.ts → metadata).
func unitBase(path string) string {
	b := filepath.Base(path)
	for _, suf := range []string{".spec.md", ".feature"} {
		if s, ok := strings.CutSuffix(b, suf); ok {
			return s
		}
	}
	if s, _, ok := strings.Cut(b, ".test."); ok {
		return s
	}
	return strings.TrimSuffix(b, filepath.Ext(b))
}

// parentDir devolve o nome da pasta imediatamente acima do arquivo ("" se na raiz).
func parentDir(path string) string {
	d := filepath.Dir(path)
	if d == "." || d == "/" || d == "" {
		return ""
	}
	return filepath.Base(d)
}

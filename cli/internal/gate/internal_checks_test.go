package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestHeaderConforme_binarioNaoCarregaCabecalho(t *testing.T) {
	// O baseline de VR é peça de prova e entrou no mapa como `kind: test` — mas é PNG.
	// Exigir cabeçalho `@anchors` num binário é impossível de cumprir, e barrava todo
	// commit de baseline visual. A identidade dele está no NOME do arquivo, que é o
	// que o `identity-consistent` confronta.
	png := "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00"
	if v, msg := checkHeaderConforme(png, mapx.Node{ID: "X.ABCDX-VR-loaded.png", Kind: mapx.KindTest}); v != Skip {
		t.Errorf("binário não carrega cabeçalho: %v (%s)", v, msg)
	}
	// Texto sem header continua reprovando — a dispensa é só para binário.
	if v, _ := checkHeaderConforme("const x = 1\n", mapx.Node{ID: "x.ts", Kind: mapx.KindCode}); v != Fail {
		t.Errorf("arquivo de texto sem header deve reprovar: %v", v)
	}
}

// "Ao menos uma regra catalogada" é o piso, e sozinho deixa passar o caso mais comum: a
// spec cataloga a primeira regra e escreve as outras em prosa. As outras ficam invisíveis
// para os gates de identidade — que então reportam verde sobre o que não conferiram.
func TestSpecSectionsCobraIrmaSemCodigo(t *testing.T) {
	// Três irmãs sob `## Regras`: duas com código, uma sem. A spec estabeleceu o padrão
	// e uma seção destoa.
	comFuro := `# Unidade

## Regras

### ABCDX-B01 — primeira

Texto.

### ABCDX-B02 — segunda

Texto.

### A terceira regra

Esta não tem código, e é regra igual às irmãs.
`
	v, msg := checkSpecSections(comFuro, mapx.Node{})
	if v != Fail {
		t.Errorf("a irmã sem código deveria reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "A terceira regra") {
		t.Errorf("a mensagem deveria NOMEAR a seção que falta: %s", msg)
	}

	// Sem furo: todas as irmãs catalogam.
	semFuro := `# Unidade

## Regras

### ABCDX-B01 — primeira

Texto.

### ABCDX-B02 — segunda

Texto.
`
	if v, msg := checkSpecSections(semFuro, mapx.Node{}); v != Pass {
		t.Errorf("todas as irmãs têm código; não havia o que cobrar: %v — %s", v, msg)
	}
}

// A seção de PROSA não pode ser cobrada: `## Visão Geral` e `## Restrições` não catalogam
// regra, e exigir código delas transformaria o gate num cobrador de formato.
func TestSpecSectionsNaoCobraProsa(t *testing.T) {
	spec := `# Unidade

## Visão Geral

Texto livre, sem código nenhum.

## Regras

### ABCDX-B01 — a regra

Texto.

## Restrições

- Uma restrição em prosa.
- Outra.

## Notas de Implementação

Mais prosa.
`
	if v, msg := checkSpecSections(spec, mapx.Node{}); v != Pass {
		t.Errorf("seções de prosa não catalogam regra e não podem ser cobradas: %v — %s", v, msg)
	}
}

// Uma irmã sozinha com código NÃO estabelece padrão. Cobrar as outras a partir de um
// único exemplo inventaria uma regra que a spec não declarou.
func TestSpecSectionsNaoInventaPadraoComUmaSo(t *testing.T) {
	spec := `# Unidade

## Regras

### ABCDX-B01 — a única com código

Texto.

### Uma seção

Texto.

### Outra seção

Texto.
`
	if v, msg := checkSpecSections(spec, mapx.Node{}); v != Pass {
		t.Errorf("uma irmã só não estabelece padrão: %v — %s", v, msg)
	}
}

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
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

// ── seção VAZIA × seção sem código ───────────────────────────────────────────
//
// São dois estados opostos que a busca por código não distingue sozinha, e confundi-los
// inverte o gate: ele passa a pedir "dê um código a cada uma" para uma tabela que não tem
// nenhuma linha.
//
// O caso real que originou: seis telas de onboarding ESTÁTICAS (zero efeitos no código)
// herdaram o cabeçalho de "Comportamentos Automáticos" do modelo de spec. O gate as
// barrava por não catalogarem regras que não existem.

func comCodigosDe4(t *testing.T) {
	t.Helper()
	orig := config.CodeLengths
	config.CodeLengths = []int{4}
	SetRuleLetters(config.DefaultRuleLetters)
	t.Cleanup(func() {
		config.CodeLengths = orig
		SetRuleLetters(config.DefaultRuleLetters)
	})
}

const specComSecaoVazia = `# Tela

## Rules

### Permissões e Acesso

| Regra | Descrição |
| ---------- | --------- |
| ` + "`PHA1-R01`" + ` | Pública |

### Ações Permitidas

| Regra | Ação | Resultado |
| ---------- | ---- | --------- |
| ` + "`PHA1-A01`" + ` | Toque | Navega |

### Comportamentos Automáticos

| Regra | Gatilho | Ação Automática |
| ---------- | ------- | --------------- |

---
`

func TestIrmasSemCodigoIgnoraSecaoVazia(t *testing.T) {
	comCodigosDe4(t)
	if d := irmasSemCodigo(specComSecaoVazia); d != "" {
		t.Fatalf("seção VAZIA não tem regra a codificar; não devia acusar. Veio: %s", d)
	}
}

// O defeito REAL continua pego: a seção que tem regra escrita em prosa, sem código.
func TestIrmasSemCodigoAindaPegaRegraSemCodigo(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n| — | Abertura | Carrega o perfil |\n",
		1)
	d := irmasSemCodigo(spec)
	if d == "" {
		t.Fatal("regra SEM código tem de ser acusada — é o defeito que o gate existe para pegar")
	}
	if !strings.Contains(d, "Comportamentos Automáticos") {
		t.Fatalf("o laudo tem de NOMEAR a seção; veio: %s", d)
	}
}

// Prosa também é conteúdo: a seção com texto e sem código segue acusada.
func TestIrmasSemCodigoPegaSecaoComProsa(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"### Comportamentos Automáticos\n\n| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"### Comportamentos Automáticos\n\nAo abrir, a tela carrega o perfil do usuário.\n",
		1)
	if d := irmasSemCodigo(spec); d == "" {
		t.Fatal("regra em prosa sem código tem de ser acusada")
	}
}

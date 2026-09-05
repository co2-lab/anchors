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

// Vazia SEM declaração: o gate não sabe se foi decisão ou esquecimento, e pede que
// alguém diga. Trocar o falso positivo antigo por um falso negativo seria pior — o gate
// passaria a reportar verde sobre seção que ninguém preencheu.
func TestSecaoVaziaSemDeclaracaoPedeQueAlguemDiga(t *testing.T) {
	comCodigosDe4(t)
	d := irmasSemCodigo(specComSecaoVazia)
	if d == "" {
		t.Fatal("vazia sem declaração tem de ser cobrada — senão o esquecimento passa")
	}
	if !strings.Contains(d, "@no-content") {
		t.Fatalf("o laudo tem de ENSINAR a saída; veio: %s", d)
	}
	// e NÃO pode pedir código para uma tabela sem linhas
	if strings.Contains(d, "Dê um código a cada uma") {
		t.Fatalf("não há 'cada uma' numa seção vazia; veio: %s", d)
	}
}

// Vazia COM declaração: a seção FICA (importa quando é obrigatória por regulação) e o
// gate absolve, porque alguém decidiu e escreveu o porquê.
func TestSecaoVaziaComNoContentEAceita(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"### Comportamentos Automáticos\n\n| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"### Comportamentos Automáticos\n\n@no-content: tela estática — não há efeito, timer nem carga.\n",
		1)
	if d := irmasSemCodigo(spec); d != "" {
		t.Fatalf("declarada, a seção vazia é aceita. Veio: %s", d)
	}
}

// `@no-content` SEM motivo não conta — dispensa sem porquê é a que ninguém revisa depois.
func TestNoContentExigeMotivo(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"### Comportamentos Automáticos\n\n| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"### Comportamentos Automáticos\n\n@no-content:\n",
		1)
	if d := irmasSemCodigo(spec); d == "" {
		t.Fatal("`@no-content` sem motivo não pode absolver")
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

// O `####` e CONTEUDO do `###` que o precede, nao irmao dele.
//
// Tratando-os como irmaos, a linha de conteudo ia para a ULTIMA secao vista — o `####`
// roubava as linhas do pai, e o `###` aparecia VAZIO. Medido em `SplashScreen.spec.md`:
// `### Caminhos Condicionais` tem tres linhas de tabela dentro de `#### destination`, e o
// gate acusava o pai de vazio com o conteudo logo abaixo.
func TestSubsecaoNaoEsvaziaOPai(t *testing.T) {
	comCodigosDe4(t)
	spec := "# Tela\n\n## Rules\n\n" +
		"### Comportamentos Automáticos\n\n" +
		"| Regra | Gatilho |\n| --- | --- |\n| `SPAX-B01` | Abertura |\n\n" +
		"### Caminhos Condicionais\n\n" +
		"#### `destination` — rota de destino\n\n" +
		"| Data State | Condição |\n| --- | --- |\n| `DS-dest-main` | autenticado |\n\n---\n"

	if d := irmasSemCodigo(spec); d != "" {
		t.Fatalf("o `###` tem conteúdo no `####` filho; não devia acusar. Veio: %s", d)
	}
}

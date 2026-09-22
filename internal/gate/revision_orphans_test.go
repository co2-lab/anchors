package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A spec-modelo reproduz a forma que originou o gate: uma revisão que trocou o badge de
// CONTAGEM para PONTO, nomeando B03 e B04 — e deixando o invariante I02, que fala do
// mesmo badge, sem menção.
const specComRevisao = `<!-- @anchors
  code: NTCNN
-->
# NotificationCenter

> **NTCNN-R0002:** o badge é um PONTO, não um número.
>
> **Revises:** ` + "`B03`, `B04`" + `

### NTCNN-B01 — A lista é de entregas, e uma entrega não muda de estado

### NTCNN-B03 — O badge é um PONTO, e não conta ao abrir

### NTCNN-B04 — Sem notificação nova, o sino fica limpo

### NTCNN-I02 — O badge nunca conta o que a lista não mostra
`

func specNodeRev() mapx.Node {
	return mapx.Node{ID: "n.spec.md", Kind: mapx.KindSpec, Code: "NTCNN"}
}

func TestRevisionOrphans_naoSpec(t *testing.T) {
	t.Run("RVORP-B01: Non-spec artifacts skip confrontation", func(t *testing.T) {})
	n := mapx.Node{ID: "p.md", Kind: mapx.KindPlan}
	if v, _ := checkRevisionOrphans(specComRevisao, n, "", nil, nil); v != Skip {
		t.Errorf("não é spec e veio %v", v)
	}
}

func TestRevisionOrphans_semRevisao(t *testing.T) {
	t.Run("RVORP-B02: A spec with no revision has nothing to confront", func(t *testing.T) {})
	semRev := "### NTCNN-B01 — uma regra\n\n### NTCNN-B02 — outra\n"
	if v, _ := checkRevisionOrphans(semRev, specNodeRev(), "", nil, nil); v != Skip {
		t.Errorf("sem revisão esperava Skip, veio %v", v)
	}
}

// Uma revisão que não nomeia nada não pode ser confrontada contra irmã nenhuma — e é a
// irmã que segue afirmando o que a revisão revogou.
func TestRevisionOrphans_semRevises(t *testing.T) {
	t.Run("RVORP-B03: A revision naming no rule cannot be confronted", func(t *testing.T) {})
	semCampo := strings.Replace(specComRevisao, "> **Revises:** `B03`, `B04`\n", "", 1)
	// PENDING, não Fail: o campo é novo, e 439 revisões escritas antes dele não podem
	// nascer todas acusadas. Pending diz a verdade — não há como confrontar — sem cobrar
	// de quem escreveu antes da régua.
	v, msg := checkRevisionOrphans(semCampo, specNodeRev(), "", nil, nil)
	if v != Pending {
		t.Fatalf("revisão sem `Revises:` esperava Pending, veio %v", v)
	}
	if !strings.Contains(msg, "Revises") {
		t.Errorf("o veredito não nomeia o campo que falta: %q", msg)
	}
}

func TestRevisionOrphans_codigoInexistente(t *testing.T) {
	t.Run("RVORP-B04: A revision naming a rule the spec does not define fails", func(t *testing.T) {})
	comFantasma := strings.Replace(specComRevisao, "`B03`, `B04`", "`B03`, `B99`", 1)
	v, msg := checkRevisionOrphans(comFantasma, specNodeRev(), "", nil, nil)
	if v != Fail {
		t.Fatalf("código inexistente e veio %v", v)
	}
	if !strings.Contains(msg, "B99") {
		t.Errorf("o veredito não traz o código: %q", msg)
	}
}

// O CASO QUE ORIGINOU O GATE: a I02 fala do mesmo badge e não foi mencionada.
func TestRevisionOrphans_irmaOrfaEAcusada(t *testing.T) {
	t.Run("RVORP-B05: A sibling sharing vocabulary and left unmentioned is reported", func(t *testing.T) {})
	v, msg := checkRevisionOrphans(specComRevisao, specNodeRev(), "", nil, nil)
	if v != Fail {
		t.Fatalf("a I02 fala do badge e não foi mencionada — esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, "I02") {
		t.Errorf("o veredito não nomeia a órfã: %q", msg)
	}
	if !strings.Contains(msg, "badge") {
		t.Errorf("o veredito não diz QUAL vocabulário liga as duas: %q", msg)
	}
	// A B01 fala de entrega e lista, não de badge: não pode entrar.
	if strings.Contains(msg, "B01") {
		t.Errorf("acusou regra que não compartilha o assunto: %q", msg)
	}
}

func TestRevisionOrphans_tudoMencionadoPassa(t *testing.T) {
	t.Run("RVORP-B06: Every vocabulary-sharing sibling accounted for passes", func(t *testing.T) {})
	completo := strings.Replace(specComRevisao,
		"> **Revises:** `B03`, `B04`", "> **Revises:** `B03`, `B04`\n>\n> **Checked:** `I02`", 1)
	if v, msg := checkRevisionOrphans(completo, specNodeRev(), "", nil, nil); v != Pass {
		t.Errorf("toda irmã contabilizada e veio %v: %s", v, msg)
	}
}

// `Checked` não afirma que a regra está certa — afirma que alguém olhou. É barato de
// escrever DEPOIS de ler, e impossível de escrever honestamente sem ler.
func TestRevisionOrphans_checkedTiraDaAcusacao(t *testing.T) {
	t.Run("RVORP-B07: Checked clears the accusation without asserting correctness", func(t *testing.T) {})
	comChecked := strings.Replace(specComRevisao,
		"> **Revises:** `B03`, `B04`", "> **Revises:** `B03`, `B04`\n>\n> **Checked:** `I02`", 1)
	_, msg := checkRevisionOrphans(comChecked, specNodeRev(), "", nil, nil)
	if strings.Contains(msg, "I02") {
		t.Errorf("a I02 foi conferida e continua acusada: %q", msg)
	}
}

func TestRevisionOrphans_regraNaoAcusaASiMesma(t *testing.T) {
	t.Run("RVORP-I01: A revised rule is never its own orphan", func(t *testing.T) {})
	_, msg := checkRevisionOrphans(specComRevisao, specNodeRev(), "", nil, nil)
	// A B03 é a própria regra revisada, e fala de badge: se se acusasse, apareceria aqui.
	for _, linha := range strings.Split(msg, "\n") {
		if strings.Contains(linha, "B03 —") {
			t.Errorf("a regra revisada acusou a si mesma: %q", linha)
		}
	}
}

// A NEGAÇÃO não liga duas regras.
//
// Metade das regras de uma spec bem escrita diz o que a unidade NÃO faz, e "não" ligava
// toda regra a toda outra. Medido na spec real: com a negação contando, três regras vinham
// acusadas; sem ela, uma — a que de fato contradizia.
func TestRevisionOrphans_negacaoNaoLiga(t *testing.T) {
	t.Run("RVORP-X01: Terms that do not name the subject do not link rules", func(t *testing.T) {})
	fixture := `> **NTCNN-R0001:** mudou.
>
> **Revises:** ` + "`B01`" + `

### NTCNN-B01 — O badge não some sozinho
### NTCNN-B02 — A ordem não muda por filtro
### NTCNN-I02 — O badge engana quando conta
`
	v, msg := checkRevisionOrphans(fixture, specNodeRev(), "", nil, nil)
	if v != Fail || !strings.Contains(msg, "I02") {
		t.Fatalf("a I02 divide `badge` com a revisada: %v / %s", v, msg)
	}
	// A B02 só divide a negação com a B01.
	if strings.Contains(msg, "B02 —") {
		t.Errorf("a negação ligou duas regras: %q", msg)
	}
}

// A DISPENSA EM COMENTÁRIO não é parte do que a regra afirma.
//
// Medido na spec real: a `B07` carrega dois `@no-*` com razão escrita e saía com 34
// termos, contra 4 a 6 das irmãs — compartilhando vocabulário com quase todas por
// acidente de prosa, não por assunto.
func TestRevisionOrphans_comentarioNaoEntraNoVocabulario(t *testing.T) {
	fixture := `> **NTCNN-R0001:** mudou.
>
> **Revises:** ` + "`B01`" + `

### NTCNN-B01 — O badge fica limpo
### NTCNN-B02 — A ordem segue a chegada <!-- @no-scenario: o badge decide, e esta unidade não desenha nada -->
`
	_, msg := checkRevisionOrphans(fixture, specNodeRev(), "", nil, nil)
	if strings.Contains(msg, "B02 —") {
		t.Errorf("o `badge` citado na DISPENSA virou assunto da regra: %q", msg)
	}
}

// E o caso REAL, reproduzido: uma palavra de domínio basta quando o título está limpo.
func TestRevisionOrphans_umaPalavraDeDominioBasta(t *testing.T) {
	fixture := `> **NTCNN-R0001:** o badge é um ponto.
>
> **Revises:** ` + "`B03`" + `

### NTCNN-B03 — O badge é um PONTO, e some ao abrir
### NTCNN-B04 — Sem notificação nova, o sino fica limpo
### NTCNN-I02 — O badge nunca conta o que a lista não mostra
`
	v, msg := checkRevisionOrphans(fixture, specNodeRev(), "", nil, nil)
	if v != Fail || !strings.Contains(msg, "I02") {
		t.Fatalf("`badge` sozinho liga quando o título está limpo: %v / %s", v, msg)
	}
	if strings.Contains(msg, "B04 —") {
		t.Errorf("a B04 não fala de badge e foi acusada: %q", msg)
	}
}

package quality

import (
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
)

// A CONTAGEM NAO E' ENDERECO.
//
// A versao anterior imprimia so' "%d alvos aguardam julgamento", e o revisor sabia que
// havia trabalho sem saber qual, onde, nem o que perguntar. Medido no app de referencia:
// 59 das 85 specs nunca receberam veredito, com o CI rodando em todos os PRs.
func judgedFixture() []gate.Result {
	return []gate.Result{
		{Gate: "regra-cumprida", Target: "a/Um.spec.md", Verdict: gate.Judge},
		{Gate: "regra-cumprida", Target: "b/Dois.spec.md", Verdict: gate.Judge},
		{Gate: "teste-prova", Target: "c/Tres.test.ts", Verdict: gate.Judge},
	}
}

func brief(t *testing.T) string {
	t.Helper()
	return capturaSaida(t, func() {
		printJudgmentBrief(judgedFixture(),
			map[string]string{"regra-cumprida": "SPEC.md"},
			map[string]string{"regra-cumprida": "O trecho REALIZA o que a regra descreve?"})
	})
}

// O ALVO: sem ele o revisor nao sabe onde olhar.
func TestBriefDeJulgamento_nomeiaCadaAlvo(t *testing.T) {
	saida := brief(t)
	for _, alvo := range []string{"a/Um.spec.md", "b/Dois.spec.md", "c/Tres.test.ts"} {
		if !strings.Contains(saida, alvo) {
			t.Errorf("o brief nao nomeia o alvo %q:\n%s", alvo, saida)
		}
	}
}

// A PERGUNTA: e' o `ask` que o gate declara, e sem ela o revisor inventa o criterio.
func TestBriefDeJulgamento_trazAPerguntaEOGuia(t *testing.T) {
	saida := brief(t)
	if !strings.Contains(saida, "REALIZA o que a regra descreve") {
		t.Errorf("o brief nao traz a pergunta do gate:\n%s", saida)
	}
	if !strings.Contains(saida, "SPEC.md") {
		t.Errorf("o brief nao diz onde esta' a regua:\n%s", saida)
	}
}

// AGRUPADO POR GATE: um gate pergunta a mesma coisa de varios alvos, e repetir a
// pergunta por alvo faria o revisor le-la tres vezes para responder uma.
func TestBriefDeJulgamento_agrupaPorGate(t *testing.T) {
	saida := brief(t)
	if n := strings.Count(saida, "REALIZA o que a regra descreve"); n != 1 {
		t.Errorf("a pergunta aparece %d vezes, esperava 1 (agrupada por gate):\n%s", n, saida)
	}
	if !strings.Contains(saida, "regra-cumprida") || !strings.Contains(saida, "teste-prova") {
		t.Errorf("o brief nao nomeia os dois gates:\n%s", saida)
	}
}

// QUEM JULGA: o pipeline entrega a lista, o revisor decide. Sem dizer isso, a lista
// parece pendencia da maquina — e ninguem a percorre.
func TestBriefDeJulgamento_dizQueOPipelineNaoJulga(t *testing.T) {
	saida := brief(t)
	achatado := strings.Join(strings.Fields(saida), " ")
	if !strings.Contains(achatado, "NAO julga") && !strings.Contains(achatado, "NOT judge") {
		t.Errorf("o brief nao diz que o pipeline nao julga:\n%s", saida)
	}
	if !strings.Contains(saida, "REV-CK5") {
		t.Errorf("o brief nao aponta o item do checklist de review:\n%s", saida)
	}
}

// ATE' DEZ, e o resto contado: sessenta caminhos afogam a pergunta que vem antes deles.
func TestBriefDeJulgamento_truncaListaLonga(t *testing.T) {
	var muitos []gate.Result
	for i := 0; i < 25; i++ {
		muitos = append(muitos, gate.Result{Gate: "regra-cumprida",
			Target: "u/" + string(rune('a'+i)) + ".spec.md", Verdict: gate.Judge})
	}
	saida := capturaSaida(t, func() {
		printJudgmentBrief(muitos, map[string]string{}, map[string]string{})
	})
	if !strings.Contains(saida, "15") {
		t.Errorf("o brief nao conta os que ficaram de fora:\n%s", saida)
	}
	if strings.Count(saida, ".spec.md") > 11 {
		t.Errorf("o brief listou mais de dez alvos:\n%s", saida)
	}
}

// O BRIEF É RELATÓRIO, NÃO REGISTRO — e a chamada tem de viver FORA do bloco de
// registro.
//
// Ele nasceu dentro do `if !noRecord`, e o efeito foi medido no app de referência: o
// pipeline roda `check --all --no-record` (deliberado — registrar dali abriria card a
// cada push), que é EXATAMENTE o modo em que o revisor precisa da lista. O único lugar
// onde o brief importava era o único em que ele não saía.
//
// A guarda é sobre a ORDEM do código-fonte porque é ali que o defeito mora: o brief
// chamado depois de `if !noRecord {` volta a sumir no pipeline, e nenhum teste de saída
// pegaria isso sem montar um projeto inteiro em modo github.
func TestBriefDeJulgamento_ficaForaDoBlocoDeRegistro(t *testing.T) {
	src, err := os.ReadFile("check.go")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(src)
	chamada := strings.Index(texto, "printJudgmentBrief(profile.Judged")
	if chamada < 0 {
		t.Fatal("a chamada do brief sumiu do fluxo do check")
	}
	registro := strings.Index(texto, "if !noRecord {")
	if registro < 0 {
		t.Fatal("o bloco de registro sumiu")
	}
	if chamada > registro {
		t.Error("o brief voltou para DENTRO do bloco de registro — ele some no " +
			"`--no-record`, que é o modo do pipeline e o único em que ele importa")
	}
}

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O gate de mutação só tem valor se distinguir três situações; um gate que sempre passa
// (ou sempre reprova) não informa nada. Cada caso abaixo fixa uma delas.
func TestMutationScore(t *testing.T) {
	casos := []struct {
		nome     string
		sig      *mapx.TestSignal
		esperado Verdict
		contem   string
	}{
		{"sem sinal ingerido → Pending, dizendo o que falta e o que se perde",
			nil, Pending, "ingest --mutation"},
		{"score acima do limiar → passa",
			&mapx.TestSignal{MutantsKilled: 90, MutantsSurvived: 5, MutationScore: 94.7}, Pass, ""},
		{"sobreviventes demais → reprova nomeando quantos",
			&mapx.TestSignal{MutantsKilled: 5, MutantsSurvived: 15, MutationScore: 25}, Fail, "15 mutante(s) sobreviveram"},
		{"limiar exato → passa (o limite não reprova)",
			&mapx.TestSignal{MutantsKilled: 7, MutantsSurvived: 3, MutationScore: 70}, Pass, ""},
		// A ferramenta RODOU e ignorou tudo — tabela de constantes com `ignoreStatic`. Nada
		// sobreviveu, entao 100 e o veredito e Pass. Antes disto o gate mandava "rode a
		// ferramenta de mutacao" sobre um arquivo em que ela ja tinha rodado: um pedido que
		// rodar de novo nao atenderia, e o tipo de ruido que ensina a ignorar o gate.
		{"tudo ignorado → passa, sem pedir rodada nova",
			&mapx.TestSignal{MutantsIgnored: 12, MutationScore: 100}, Pass, ""},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			n := mapx.Node{Kind: mapx.KindCode, Rev: "r1", Signal: c.sig}
			if c.sig != nil {
				c.sig.AtRev = "r1" // sinal fresco; o stale tem caso próprio
			}
			v, d := checkMutationScore("", n)
			if v != c.esperado {
				t.Fatalf("veredito = %s, queria %s (detalhe: %s)", v, c.esperado, d)
			}
			if c.contem != "" && !strings.Contains(d, c.contem) {
				t.Fatalf("detalhe %q não menciona %q", d, c.contem)
			}
		})
	}
}

// Sinal de mutação medido numa versão ANTERIOR do arquivo não prova nada sobre a atual —
// é a mesma regra dos outros sinais ingeridos, e a que mais engana: o número parece bom.
func TestMutationScoreStale(t *testing.T) {
	n := mapx.Node{Kind: mapx.KindCode, Rev: "r2",
		Signal: &mapx.TestSignal{MutantsKilled: 100, MutationScore: 100, AtRev: "r1"}}
	v, d := checkMutationScore("", n)
	if v != Pending {
		t.Fatalf("score perfeito de uma revisão antiga deveria ser Pending, foi %s", v)
	}
	if !strings.Contains(d, "stale") {
		t.Fatalf("detalhe não explica o stale: %q", d)
	}
}

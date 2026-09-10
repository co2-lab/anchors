package main

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// PERDA DE CARIMBO é silenciosa, e o mapa continua VÁLIDO.
//
// O `PreserveStamps` preserva o que está no arquivo anterior — e num merge o arquivo
// anterior é o do OUTRO lado. Medido no blue-eyes (co2-lab/anchors#12): resolvi um
// conflito com `checkout --theirs` + `map build`, e o carimbo de `review` de uma spec
// desapareceu. Descobri por acaso, contando: 18 onde eu esperava 19.
//
// O custo não é o carimbo — é o LAUDO. Ele vive no `--reason` do comando, não no arquivo,
// e refazer uma revisão adversarial é caro.
func TestStampLossWarning(t *testing.T) {
	comJulgamentos := func(pares ...string) *mapx.Graph {
		g := &mapx.Graph{}
		for _, gate := range pares {
			g.Edges = append(g.Edges, mapx.Edge{
				Julgamentos: []mapx.Judgment{{Gate: gate}},
			})
		}
		return g
	}

	t.Run("perda AVISA, e conta por gate", func(t *testing.T) {
		antigo := comJulgamentos("review", "review", "rule-fulfilled")
		novo := comJulgamentos("review", "rule-fulfilled")

		aviso := stampLossWarning(antigo, novo)
		if aviso == "" {
			t.Fatal("perdeu um carimbo de review e não avisou")
		}
		if !strings.Contains(aviso, "review") || !strings.Contains(aviso, "2 → 1") {
			t.Errorf("o aviso não diz QUAL gate e quanto:\n%s", aviso)
		}
		// contar por gate e não o total: "perdeu 1 carimbo" não diz o que refazer
		if strings.Contains(aviso, "rule-fulfilled") {
			t.Errorf("acusou um gate que não perdeu nada:\n%s", aviso)
		}
		// o aviso tem de dizer o que se perde de fato
		if !strings.Contains(aviso, "laudo") {
			t.Errorf("o aviso não diz que o LAUDO se perde:\n%s", aviso)
		}
	})

	// Um aviso que sempre aparece deixa de ser lido.
	t.Run("sem perda, silêncio", func(t *testing.T) {
		igual := comJulgamentos("review", "rule-fulfilled")
		if got := stampLossWarning(igual, comJulgamentos("review", "rule-fulfilled")); got != "" {
			t.Errorf("avisou sem perda: %q", got)
		}
		// GANHAR carimbo também é silêncio: é o caso normal de julgar algo novo
		mais := comJulgamentos("review", "review", "rule-fulfilled")
		if got := stampLossWarning(igual, mais); got != "" {
			t.Errorf("avisou ao GANHAR carimbo: %q", got)
		}
	})

	t.Run("mapa anterior vazio não avisa", func(t *testing.T) {
		if got := stampLossWarning(&mapx.Graph{}, comJulgamentos("review")); got != "" {
			t.Errorf("avisou com anterior vazio: %q", got)
		}
	})
}

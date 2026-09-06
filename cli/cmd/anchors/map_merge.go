package main

import (
	"fmt"
	"os"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// --- o MERGE DRIVER do mapa ---
//
// O `anchors.graph.yaml` é derivado, e o git não sabe disso: ele o trata como texto e
// mescla linha a linha. Quando os dois lados tocam o mesmo trecho, ele pede resolução
// manual — e um grafo com centenas de arestas e carimbos aninhados não tem resolução
// manual confiável. Quando NÃO tocam, ele resolve sozinho, e é aí que o dano acontece.
//
// Medido no blue-eyes (co2-lab/anchors#12): um `git merge origin/develop` mesclou o mapa
// sem conflito e apagou SESSENTA E DOIS carimbos de julgamento:
//
//	git show HEAD --stat | grep graph
//	anchors.graph.yaml | 1212 -------------------------------------------------
//
// O aviso do `map build` (v0.1.43) não pega este caso: ele compara com o arquivo ANTERIOR,
// e a perda acontece no `git merge`, antes de o Anchors ser chamado. O `map build`
// seguinte compara vazio com vazio e não tem o que acusar.
//
// A saída é pôr o Anchors no caminho, e é o que este comando faz. No `.gitattributes`:
//
//	anchors.graph.yaml merge=anchors-map
//
// e o driver registrado no git config. A partir daí o git chama `anchors map merge` em vez
// de tentar o merge textual.
//
// O QUE ELE FAZ é o que eu fiz à mão ao consertar o estrago: unir os carimbos dos dois
// lados. Um carimbo é ESTADO DO TRABALHO — alguém julgou, e o laudo vive no `--reason` do
// `anchors judge`, não no arquivo. Perder o carimbo manda o julgamento de volta para a
// fila e a evidência tem de ser refeita.
//
// União e não escolha: se um lado julgou a aresta A e o outro a aresta B, o resultado tem
// as duas. Quando os DOIS julgaram a mesma aresta, vence o mais recente pela rev — que é a
// regra que o `PreserveStamps` já aplica.

func newMapMergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge <base> <nosso> <deles>",
		Short: "Merge driver do mapa: une os carimbos dos dois lados (git merge=anchors-map)",
		Long: `Resolve o merge do ` + "`anchors.graph.yaml`" + ` unindo os carimbos de julgamento
dos dois lados, em vez de mesclar texto.

Não é para rodar à mão. Configure o git para chamá-lo:

  # .gitattributes (no repositório)
  anchors.graph.yaml merge=anchors-map

  # uma vez por clone
  git config merge.anchors-map.name "anchors: une carimbos do mapa"
  git config merge.anchors-map.driver "anchors map merge %O %A %B"

POR QUÊ: o mapa é derivado, e o git o trata como texto. Medido: um merge sem conflito
apagou 62 carimbos de julgamento — o laudo de cada um vive no comando que o gravou, não
no arquivo, e refazer uma revisão adversarial é caro.

O resultado é escrito em <nosso>, que é o que o git espera.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// O git passa três caminhos: %O (base comum), %A (nosso, e é onde o
			// resultado vai), %B (deles). A base não é usada aqui — a união não
			// precisa saber o que havia antes, porque carimbo não se apaga por
			// ausência: se um lado tem e o outro não, o que tem vence.
			nosso, deles := args[1], args[2]

			gNosso, err := mapx.Load(nosso)
			if err != nil {
				return fmt.Errorf("ler o nosso lado (%s): %w", nosso, err)
			}
			gDeles, err := mapx.Load(deles)
			if err != nil {
				return fmt.Errorf("ler o outro lado (%s): %w", deles, err)
			}

			antes := countJudgments(gNosso)

			// PRIMEIRO as arestas que só existem do outro lado.
			//
			// O `PreserveStamps` copia carimbo para aresta que JÁ ESTÁ no destino — é o
			// que ele foi feito para fazer (preservar entre reconstruções do mesmo
			// projeto). Num merge os dois lados podem ter arestas diferentes: uma spec
			// nova de cada branch.
			//
			// Medido no teste: a aresta `c.spec.md` existia só do outro lado, e o
			// julgamento dela sumia — o `PreserveStamps` não tinha onde pousá-lo.
			addMissingEdges(gNosso, gDeles)

			// E DEPOIS os carimbos das arestas comuns. Duas passadas porque
			// `PreserveStamps(novo, antigo)` leva numa direção só; rodar nos dois
			// sentidos garante que nenhum lado perca o que era só dele.
			mapx.PreserveStamps(gNosso, gDeles)
			mapx.PreserveStamps(gDeles, gNosso)
			mapx.PreserveStamps(gNosso, gDeles)

			if err := mapx.Save(gNosso, nosso); err != nil {
				return fmt.Errorf("gravar o resultado: %w", err)
			}
			depois := countJudgments(gNosso)
			fmt.Fprintf(os.Stderr, "anchors: mapa mesclado — %d julgamento(s) preservado(s)",
				depois)
			if depois > antes {
				fmt.Fprintf(os.Stderr, " (%d vieram do outro lado)", depois-antes)
			}
			fmt.Fprintln(os.Stderr)
			return nil
		},
	}
	return cmd
}

func countJudgments(g *mapx.Graph) int {
	if g == nil {
		return 0
	}
	n := 0
	for _, e := range g.Edges {
		n += len(e.Julgamentos)
	}
	return n
}

// addMissingEdges copia para `destino` as arestas que só existem em `origem`.
//
// Um merge junta dois branches que podem ter criado unidades diferentes — e a aresta que
// só existe de um lado leva o julgamento dela junto. Sem isto, o `PreserveStamps` não tem
// onde pousar o carimbo, e ele some sem aviso.
//
// Só acrescenta o que FALTA: aresta comum é resolvida pelos `PreserveStamps` seguintes,
// que sabem comparar rev.
func addMissingEdges(destino, origem *mapx.Graph) {
	if destino == nil || origem == nil {
		return
	}
	tem := make(map[string]bool, len(destino.Edges))
	for _, e := range destino.Edges {
		tem[edgeKey(e)] = true
	}
	for _, e := range origem.Edges {
		if !tem[edgeKey(e)] {
			destino.Edges = append(destino.Edges, e)
		}
	}
}

func edgeKey(e mapx.Edge) string {
	return string(e.Type) + "\x00" + e.From + "\x00" + e.To
}

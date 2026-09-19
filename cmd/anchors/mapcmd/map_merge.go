package mapcmd

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
		Use:   "merge <base> <ours> <theirs>",
		Short: "Map merge driver: unites the stamps from both sides (git merge=anchors-map)",
		Long: `Resolves the merge of ` + "`anchors.graph.yaml`" + ` by uniting the judgment stamps
from both sides, instead of merging text.

Not meant to be run by hand. Configure git to call it:

  # .gitattributes (in the repository)
  anchors.graph.yaml merge=anchors-map

  # once per clone
  git config merge.anchors-map.name "anchors: unites the map's stamps"
  git config merge.anchors-map.driver "anchors map merge %O %A %B"

WHY: the map is derived, and git treats it as text. Measured: a conflict-free merge
erased 62 judgment stamps — the report of each one lives in the command that wrote it, not
in the file, and redoing an adversarial review is expensive.

The result is written to <ours>, which is what git expects.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// O git passa três caminhos: %O (base comum), %A (nosso, e é onde o
			// resultado vai), %B (deles). A base não é usada aqui — a união não
			// precisa saber o que havia antes, porque carimbo não se apaga por
			// ausência: se um lado tem e o outro não, o que tem vence.
			nosso, deles := args[1], args[2]

			gNosso, err := mapx.Load(nosso)
			if err != nil {
				return fmt.Errorf("read our side (%s): %w", nosso, err)
			}
			gDeles, err := mapx.Load(deles)
			if err != nil {
				return fmt.Errorf("read the other side (%s): %w", deles, err)
			}

			antes := countJudgments(gNosso)

			// PRIMEIRO os NÓS que só existem do outro lado.
			//
			// `Graph` guarda `Nodes` e `Edges` em listas separadas, e por três versões
			// este driver reconciliou só a segunda. O resultado saía com a lista de nós
			// do lado `nosso`, e todo nó criado só no outro branch desaparecia.
			//
			// MEDIDO no blue-eyes (#730): base 329 nós, nosso 330, deles 332, resultado
			// 330 — e o git reporta "Automatic merge went well".
			//
			// O dano é calado por construção: o arquivo continua REGIDO, o `check` o
			// reconhece, e ele não está no mapa — então nenhum gate o confronta. O
			// `map build` seguinte reinsere os nós, o que esconde o defeito de quem
			// mescla e reconstrói, e o publica para quem mescla e empurra.
			addMissingNodes(gNosso, gDeles)

			// DEPOIS as arestas que só existem do outro lado.
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
				return fmt.Errorf("write the result: %w", err)
			}
			depois := countJudgments(gNosso)
			fmt.Fprintf(os.Stderr, "anchors: map merged — %d judgment(s) preserved",
				depois)
			if depois > antes {
				fmt.Fprintf(os.Stderr, " (%d came from the other side)", depois-antes)
			}
			fmt.Fprintln(os.Stderr)
			return nil
		},
	}
	return cmd
}

// addMissingNodes copia para `destino` os nós que só existem em `origem`.
//
// Espelha o `addMissingEdges`: um merge junta branches que criaram unidades diferentes, e
// o nó que só existe de um lado tem de sobreviver. Nó comum NÃO é tocado — a rev dele é
// derivada do conteúdo do arquivo, e quem a resolve é o `map build` sobre a árvore
// mesclada, não este driver, que só vê os dois mapas.
func addMissingNodes(destino, origem *mapx.Graph) {
	if destino == nil || origem == nil {
		return
	}
	tem := make(map[string]bool, len(destino.Nodes))
	for _, n := range destino.Nodes {
		tem[n.ID] = true
	}
	for _, n := range origem.Nodes {
		if !tem[n.ID] {
			destino.Nodes = append(destino.Nodes, n)
		}
	}
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

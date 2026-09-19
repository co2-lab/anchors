package flow

import (
	"fmt"
	"os"

	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// --- o MERGE DRIVER do progresso ---
//
// O irmão do `anchors map merge`, pelo mesmo motivo e com a mesma forma.
//
// Medido no blue-eyes: TRÊS PRs consecutivos conflitaram no `*-progress.md` do mesmo
// plano (#217, #221, e o seguinte teria), sempre porque duas branches marcaram
// checkboxes vizinhos. No mesmo período o `anchors.graph.yaml` — que TEM driver desde a
// v0.1.48 — não conflitou nenhuma vez.
//
// A diferença entre os dois arquivos não é a complexidade: é que um tinha driver.
//
// Aqui o dano do merge textual é menor que no mapa (o git PEDE resolução em vez de
// apagar em silêncio), mas o atrito é diário: a resolução é sempre a mesma, é mecânica,
// e interrompe o fluxo no pior momento — entre o commit e o PR.
//
// A REGRA é uma só: `[x]` vence `[ ]`. Marcar concluído é um FATO — a spec foi entregue,
// o `anchors check` a carimbou, o PR mergeou. Um merge que desmarca apaga o fato, e o
// gate `progress-honest` passa a acusar como não feito um arquivo que existe.
func newProgressMergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-progress <base> <ours> <theirs>",
		Short: "Progress merge driver: unites both sides, `[x]` beats `[ ]`",
		Long: `Resolves the merge of a ` + "`*-progress.md`" + ` by uniting the items of both sides.

Not meant to be run by hand. Configure git to call it:

  # .gitattributes (no repositório)
  *-progress.md merge=anchors-progress

  # uma vez por clone
  git config merge.anchors-progress.name "anchors: unites the progress of the plans"
  git config merge.anchors-progress.driver "anchors merge-progress %O %A %B"

WHY: two branches delivering specs of the SAME plan mark neighboring checkboxes, and
git asks for manual resolution every time. Measured: three consecutive PRs, identical
resolution in all three. The rule is ` + "`[x]` beats `[ ]`" + ` — unmarking by merge would
erase a delivery that already happened.

The result is written to <ours>, which is what git expects.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// %O (base), %A (nosso, onde o resultado vai), %B (deles). A base não é
			// usada: não existe "desmarcar" legítimo vindo de um merge — quem quer
			// desmarcar edita o arquivo, e a edição chega como um lado só.
			nosso, deles := args[1], args[2]

			bNosso, err := os.ReadFile(nosso)
			if err != nil {
				return fmt.Errorf("read our side (%s): %w", nosso, err)
			}
			bDeles, err := os.ReadFile(deles)
			if err != nil {
				return fmt.Errorf("read the other side (%s): %w", deles, err)
			}

			antes := scan.ProgressDone(string(bNosso))
			res := scan.MergeProgress(string(bNosso), string(bDeles))

			if err := os.WriteFile(nosso, []byte(res), 0o644); err != nil {
				return fmt.Errorf("write the result: %w", err)
			}

			depois := scan.ProgressDone(res)
			fmt.Fprintf(os.Stderr, "anchors: progress merged — %d done", depois)
			if depois > antes {
				fmt.Fprintf(os.Stderr, " (%d from the other side)", depois-antes)
			}
			fmt.Fprintln(os.Stderr)
			return nil
		},
	}
	return cmd
}

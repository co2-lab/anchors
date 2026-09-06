package main

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
		Use:   "merge-progress <base> <nosso> <deles>",
		Short: "Merge driver do progresso: une os dois lados, `[x]` vence `[ ]`",
		Long: `Resolve o merge de um ` + "`*-progress.md`" + ` unindo os itens dos dois lados.

Não é para rodar à mão. Configure o git para chamá-lo:

  # .gitattributes (no repositório)
  *-progress.md merge=anchors-progress

  # uma vez por clone
  git config merge.anchors-progress.name "anchors: une o progresso dos planos"
  git config merge.anchors-progress.driver "anchors merge-progress %O %A %B"

POR QUÊ: duas branches que entregam specs do MESMO plano marcam checkboxes vizinhos, e o
git pede resolução manual toda vez. Medido: três PRs consecutivos, resolução idêntica nos
três. A regra é ` + "`[x]` vence `[ ]`" + ` — desmarcar por merge apagaria uma entrega
que já aconteceu.

O resultado é escrito em <nosso>, que é o que o git espera.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// %O (base), %A (nosso, onde o resultado vai), %B (deles). A base não é
			// usada: não existe "desmarcar" legítimo vindo de um merge — quem quer
			// desmarcar edita o arquivo, e a edição chega como um lado só.
			nosso, deles := args[1], args[2]

			bNosso, err := os.ReadFile(nosso)
			if err != nil {
				return fmt.Errorf("ler o nosso lado (%s): %w", nosso, err)
			}
			bDeles, err := os.ReadFile(deles)
			if err != nil {
				return fmt.Errorf("ler o outro lado (%s): %w", deles, err)
			}

			antes := scan.ProgressDone(string(bNosso))
			res := scan.MergeProgress(string(bNosso), string(bDeles))

			if err := os.WriteFile(nosso, []byte(res), 0o644); err != nil {
				return fmt.Errorf("gravar o resultado: %w", err)
			}

			depois := scan.ProgressDone(res)
			fmt.Fprintf(os.Stderr, "anchors: progresso mesclado — %d concluído(s)", depois)
			if depois > antes {
				fmt.Fprintf(os.Stderr, " (%d do outro lado)", depois-antes)
			}
			fmt.Fprintln(os.Stderr)
			return nil
		},
	}
	return cmd
}

package flow

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// newDiscardCmd tira do board o card que não faz mais sentido.
//
// FECHAR NÃO BASTA, e é o que obrigou este comando. O board mostra os cards fechados porque
// o roadmap precisa deles — é como ele desenha o que já foi entregue. Um card de teste, ou
// um achado sobre arquivo que não existe mais, fica lá para sempre: fechado, sem pai
// possível, ocupando uma linha na raiz da árvore.
//
// MEDIDO no projeto de referência: das 27 raízes do roadmap, DOZE eram ruído permanente.
//
// SOFT-DELETE, e não remoção. Apagar perderia o rastro de que a pergunta existiu — e é isso
// que distingue "resolvido" de "descartado". O card continua no GitHub, com a razão em
// comentário; quem procurar o número o acha.
func newDiscardCmd() *cobra.Command {
	var root, motivo string
	cmd := &cobra.Command{
		Use:   "discard <card>...",
		Short: "Take off the board the card that no longer makes sense, without deleting it",
		Long: `Marks one or more cards as DISCARDED.

The card leaves the board — columns, tree, roadmap, pending strip — and stays on
GitHub with the reason recorded. It is for what no longer makes sense, not for what
finished: delivered work CLOSES, and the roadmap needs it to draw the past.

The cases that motivated it, in the reference project:

  · test cards of the flow itself ("[test] state lock — delete")
  · findings about files that no longer exist (` + "`Orfa.spec.md`, `ZzTeste.spec.md`" + `)
  · questions the project answered by another path

DO NOT use it for work that got done — that is closing. Nor for what still matters and
nobody took: a forgotten card is still work, and hiding it from the board is the
quietest way of losing it.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if strings.TrimSpace(motivo) == "" {
				return fmt.Errorf("`--reason` is REQUIRED: a card that disappears from the board without " +
					"saying why is indistinguishable from one that was lost")
			}

			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return err
			}
			if !cfg.GitHubMode() {
				return fmt.Errorf("`discard` exists in github mode (the card is an issue). " +
					"In local mode, delete the card file")
			}

			// A LABEL É CRIADA SOB DEMANDA. O `doctor --fix` a garante, mas quem roda o
			// `discard` num projeto que ainda não rodou o doctor não deveria ser barrado
			// por isso — e `gh issue edit` com label inexistente falha o comando INTEIRO.
			_ = exec.Command("gh", "label", "create", initx.LabelDiscarded,
				"--repo", cfg.Workflow.Repo,
				"--color", "d4d4d4",
				"--description", "off the board: no longer makes sense",
			).Run()

			var falhas []string
			for _, a := range args {
				card := strings.TrimPrefix(strings.TrimSpace(a), "#")
				if card == "" {
					continue
				}
				if err := descarta(cfg.Workflow.Repo, card, motivo); err != nil {
					falhas = append(falhas, fmt.Sprintf("#%s: %v", card, err))
					continue
				}
				fmt.Printf("· #%s discarded — off the board, still on GitHub\n", card)
			}
			if len(falhas) > 0 {
				return fmt.Errorf("%d card(s) were not discarded:\n  %s",
					len(falhas), strings.Join(falhas, "\n  "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&motivo, "reason", "", "why this card no longer makes sense")
	common.AliasDeFlag(cmd, "reason", "motivo")
	return cmd
}

// descarta aplica a label e registra a razão — nessa ordem.
//
// A RAZÃO PRIMEIRO seria pior: se a label falhar, o card fica com um comentário dizendo que
// foi descartado e continua no board, o que é a pior das duas inconsistências possíveis.
func descarta(repo, card, motivo string) error {
	out, err := exec.Command("gh", "issue", "edit", card,
		"--repo", repo, "--add-label", initx.LabelDiscarded).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}

	corpo := "🗑 **Discarded** — this card left the board.\n\n" + motivo +
		"\n\nIt stays here: discarding does not erase, because the trace that the question " +
		"existed has value. To bring it back, remove the label `" +
		initx.LabelDiscarded + "`."
	if out, err := exec.Command("gh", "issue", "comment", card,
		"--repo", repo, "--body", corpo).CombinedOutput(); err != nil {
		return fmt.Errorf("record the reason: %w: %s", err, strings.TrimSpace(string(out)))
	}

	// FECHAR VEM POR ÚLTIMO, e só se ainda estiver aberto: um card descartado que
	// continua aberto seria servido pelo `claim`.
	_ = exec.Command("gh", "issue", "close", card, "--repo", repo).Run()
	return nil
}

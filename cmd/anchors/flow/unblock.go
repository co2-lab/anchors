package flow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- a decisão que GERA TRABALHO ---
//
// `needs-user` diz que um card espera uma pessoa, e o claim não o entrega enquanto a label
// estiver lá. A instrução que o pipeline escreve no card é "decida, e depois remova a
// label" — e ela supõe que decidir é só escrever.
//
// Muitas decisões não são. A pessoa decide, e a decisão exige que alguém MUDE alguma coisa
// antes de o card original seguir: uma spec que precisa ganhar uma regra, um contrato que
// precisa mudar, um defeito que precisa ser consertado noutro lugar.
//
// SEM VÍNCULO, os dois caminhos são ruins:
//
//   - a pessoa abre o card da mudança e deixa o `needs-user` no original. Ele fica parado
//     para sempre, porque a condição de remover a label ("decida") já foi cumprida e
//     ninguém sabe que falta o trabalho.
//
//   - a pessoa remove a label achando que decidir bastava. O card volta à fila, o agente
//     que o pega encontra o mesmo impasse, e escala de novo — o ciclo se repete.
//
// `anchors unblock` fecha isso: cria o card da mudança já ligado ao bloqueado, e diz ao
// bloqueado que ele espera aquele card. Quando o card da mudança fecha, o outro volta à
// fila sozinho.

func newUnblockCmd() *cobra.Command {
	var root, motivo, sobre string
	cmd := &cobra.Command{
		Use:   "unblock <card>",
		Short: "Open the work card that unblocks a card stuck in `needs-user`",
		Long: `Creates the card of the change a decision demanded, linked to the blocked card.

The card in ` + "`needs-user`" + ` waits for a person. When that person's decision
GENERATES WORK — a spec that gains a rule, a contract that changes, a defect
somewhere else —, the work becomes a card of its own, and the blocked one waits for
IT instead of waiting indefinitely for someone who has already decided.

The new card is born with ` + "`anchors:desbloqueia-<n>`" + `. When it closes, the
blocked one returns to the queue.

    anchors unblock 311 --reason "the dispatch loop needs try/catch per token"
    anchors unblock 311 --reason "..." --about packages/lambdas/push/NotificationDispatcher.ts`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			bloqueado := strings.TrimPrefix(args[0], "#")
			if strings.TrimSpace(motivo) == "" {
				return fmt.Errorf("`--reason` is REQUIRED: the new card needs to say WHAT to change, " +
					"and whoever takes it does not have the context of the decision")
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
				return fmt.Errorf("`unblock` exists in github mode: in local mode the block " +
					"is not a label, and the card moves between folders")
			}

			corpo, err := corpoDoDesbloqueio(bloqueado, motivo, sobre)
			if err != nil {
				return err
			}
			defer os.Remove(corpo)

			// A label é criada SOB DEMANDA, como a `sob-<n>` do `escalate`: ela é uma por
			// card bloqueado, e pré-criar todas seria impossível. O erro é ignorado porque
			// "já existe" é o caso comum a partir do segundo card de desbloqueio.
			_ = exec.Command("gh", "label", "create", initx.LabelDesbloqueia(bloqueado),
				"--repo", cfg.Workflow.Repo,
				"--color", "d4c5f9",
				"--description", "the delivery of this card unblocks #"+bloqueado,
			).Run()

			argv := []string{"issue", "create",
				"--repo", cfg.Workflow.Repo,
				"--title", "[unblocks #" + bloqueado + "] " + firstLineOfReason(motivo),
				"--body-file", corpo,
				"--label", cfg.Workflow.Labels[0],
				"--label", "anchors:to-do",
				"--label", initx.LabelDesbloqueia(bloqueado),
			}
			out, err := exec.Command("gh", argv...).CombinedOutput()
			if err != nil {
				return fmt.Errorf("create the card: %w\n%s", err, out)
			}
			url := strings.TrimSpace(string(out))
			fmt.Printf("unblock card created: %s\n", url)

			// O CARD BLOQUEADO precisa saber por quem espera. Sem este comentário, quem
			// abre o #311 vê `needs-user` e a instrução "decida e remova a label" — e não
			// tem como descobrir que a decisão já saiu e virou trabalho.
			aviso := fmt.Sprintf(
				"⏸ **This card is waiting for the delivery of %s.**\n\n"+
					"The decision was made and generated work: %s\n\n"+
					"The label `anchors:needs-user` STAYS until then — the claim does not deliver this card "+
					"while it is here, and that is what keeps another agent from taking it and "+
					"running into the same impasse.\n\n"+
					"When the card above is delivered, remove the label and this one goes back to the queue.",
				url, motivo)
			if err := exec.Command("gh", "issue", "comment", bloqueado,
				"--repo", cfg.Workflow.Repo, "--body", aviso).Run(); err != nil {
				// O card foi criado; falhar aqui não desfaz isso. Avisar é melhor que
				// abortar e deixar o usuário sem saber o que existe e o que não existe.
				fmt.Printf("⚠ the card was created, but I could not comment on #%s: %v\n", bloqueado, err)
			}
			fmt.Printf("#%s remains stopped, and now SAYS who it is waiting for\n", bloqueado)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&motivo, "reason", "", "REQUIRED — what needs to change")
	cmd.Flags().StringVar(&sobre, "about", "", "file the change touches")
	return cmd
}

// corpoDoDesbloqueio escreve o corpo do card num arquivo temporário.
//
// Arquivo e não `--body`: o motivo tem parágrafos, e um argumento de linha de comando com
// quebras de linha atravessa shells de jeitos diferentes.
func corpoDoDesbloqueio(bloqueado, motivo, sobre string) (string, error) {
	var b strings.Builder
	b.WriteString("🔓 **The delivery of this card unblocks #" + bloqueado + ".**\n\n")
	b.WriteString(motivo + "\n\n")
	if sobre != "" {
		b.WriteString("**Where:** `" + sobre + "`\n\n")
	}
	b.WriteString("**Where it came from:** #" + bloqueado + " stopped at `anchors:needs-user` — " +
		"it was waiting on a decision no agent could make. The decision came out, and demanded this " +
		"change.\n\n")
	// O QUE DISTINGUE este card dos outros: quem o pega não precisa entender o impasse
	// original, só fazer o que está escrito. Dizer isso evita que ele vá ler o #311 inteiro.
	b.WriteString("**What to do:** what is described above, through the normal flow. You do NOT " +
		"need to reopen the discussion of #" + bloqueado + " — it was already resolved, and what " +
		"is left is this work.\n\n")
	b.WriteString("When this card is delivered, remove the label `anchors:needs-user` from #" +
		bloqueado + " and it goes back to the queue.\n")

	tmp, err := os.CreateTemp("", "anchors-unblock-*.md")
	if err != nil {
		return "", err
	}
	if _, err := tmp.WriteString(b.String()); err != nil {
		tmp.Close()
		return "", err
	}
	return tmp.Name(), tmp.Close()
}

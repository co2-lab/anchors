package main

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
		Short: "Abre o card de trabalho que destrava um card parado em `needs-user`",
		Long: `Cria o card da mudança que uma decisão exigiu, ligado ao card bloqueado.

O card em ` + "`needs-user`" + ` espera uma pessoa. Quando a decisão dessa pessoa
GERA TRABALHO — uma spec que ganha regra, um contrato que muda, um defeito
noutro lugar —, o trabalho vira um card próprio, e o bloqueado espera por ELE em
vez de esperar indefinidamente por alguém que já decidiu.

O card novo nasce com ` + "`anchors:desbloqueia-<n>`" + `. Quando ele fecha, o
bloqueado volta à fila.

    anchors unblock 311 --reason "o laço de dispara precisa de try/catch por token"
    anchors unblock 311 --reason "..." --about packages/lambdas/push/NotificationDispatcher.ts`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			bloqueado := strings.TrimPrefix(args[0], "#")
			if strings.TrimSpace(motivo) == "" {
				return fmt.Errorf("`--reason` é obrigatório: o card novo precisa dizer O QUE mudar, " +
					"e quem o pegar não tem o contexto da decisão")
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
				return fmt.Errorf("`unblock` existe no modo github: no modo local o bloqueio " +
					"não é uma label, e o card se move de pasta")
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
				"--description", "a entrega deste card destrava o #"+bloqueado,
			).Run()

			argv := []string{"issue", "create",
				"--repo", cfg.Workflow.Repo,
				"--title", "[destrava #" + bloqueado + "] " + firstLineOfReason(motivo),
				"--body-file", corpo,
				"--label", cfg.Workflow.Labels[0],
				"--label", "anchors:to-do",
				"--label", initx.LabelDesbloqueia(bloqueado),
			}
			out, err := exec.Command("gh", argv...).CombinedOutput()
			if err != nil {
				return fmt.Errorf("criar o card: %w\n%s", err, out)
			}
			url := strings.TrimSpace(string(out))
			fmt.Printf("card de desbloqueio criado: %s\n", url)

			// O CARD BLOQUEADO precisa saber por quem espera. Sem este comentário, quem
			// abre o #311 vê `needs-user` e a instrução "decida e remova a label" — e não
			// tem como descobrir que a decisão já saiu e virou trabalho.
			aviso := fmt.Sprintf(
				"⏸ **Este card espera a entrega de %s.**\n\n"+
					"A decisão foi tomada e gerou trabalho: %s\n\n"+
					"A label `anchors:needs-user` FICA até lá — o claim não entrega este card "+
					"enquanto ela estiver aqui, e é isso que impede outro agente de pegá-lo e "+
					"esbarrar no mesmo impasse.\n\n"+
					"Quando o card acima for entregue, remova a label e este volta à fila.",
				url, motivo)
			if err := exec.Command("gh", "issue", "comment", bloqueado,
				"--repo", cfg.Workflow.Repo, "--body", aviso).Run(); err != nil {
				// O card foi criado; falhar aqui não desfaz isso. Avisar é melhor que
				// abortar e deixar o usuário sem saber o que existe e o que não existe.
				fmt.Printf("⚠ o card foi criado, mas não consegui comentar no #%s: %v\n", bloqueado, err)
			}
			fmt.Printf("#%s continua parado, e agora DIZ por quem espera\n", bloqueado)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&motivo, "reason", "", "OBRIGATÓRIO — o que precisa mudar")
	cmd.Flags().StringVar(&sobre, "about", "", "arquivo que a mudança toca")
	return cmd
}

// corpoDoDesbloqueio escreve o corpo do card num arquivo temporário.
//
// Arquivo e não `--body`: o motivo tem parágrafos, e um argumento de linha de comando com
// quebras de linha atravessa shells de jeitos diferentes.
func corpoDoDesbloqueio(bloqueado, motivo, sobre string) (string, error) {
	var b strings.Builder
	b.WriteString("🔓 **A entrega deste card destrava o #" + bloqueado + ".**\n\n")
	b.WriteString(motivo + "\n\n")
	if sobre != "" {
		b.WriteString("**Onde:** `" + sobre + "`\n\n")
	}
	b.WriteString("**De onde veio:** o #" + bloqueado + " parou em `anchors:needs-user` — " +
		"esperava uma decisão que nenhum agente podia tomar. A decisão saiu, e exigiu esta " +
		"mudança.\n\n")
	// O QUE DISTINGUE este card dos outros: quem o pega não precisa entender o impasse
	// original, só fazer o que está escrito. Dizer isso evita que ele vá ler o #311 inteiro.
	b.WriteString("**O que fazer:** o que está descrito acima, pelo fluxo normal. Você NÃO " +
		"precisa reabrir a discussão do #" + bloqueado + " — ela já foi resolvida, e o que " +
		"sobrou é este trabalho.\n\n")
	b.WriteString("Quando este card for entregue, remova a label `anchors:needs-user` do #" +
		bloqueado + " e ele volta à fila.\n")

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

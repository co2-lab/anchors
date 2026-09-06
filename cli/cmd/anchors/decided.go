package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- a outra metade do `escalate --for-user` ---
//
// O `escalate --for-user` grava um estado que ele não sabia reverter: põe
// `anchors:needs-user` no card e o claim para de entregá-lo. Isso é o comportamento certo
// — é o que impede o agente seguinte de refazer o caminho até a mesma dúvida.
//
// O que faltava era o caminho de volta. Medido no blue-eyes (co2-lab/anchors#10): a
// decisão saiu, as revisões foram aplicadas nos 4 arquivos, a issue de decisão foi
// fechada — e o card continuou parado. Removi a label à mão.
//
// A instrução de remover existia, mas dentro do corpo do card que o WORKFLOW abre. Quem
// escala pelo CLI recebia só `· card #4 parado até a decisão`, e descobrir como retomar
// exigia grepar o YAML do workflow.
//
// O modo de falha é o silencioso de sempre: o card fica parado, o `anchors next` o pula, e
// ninguém percebe até alguém perguntar por que o plano não avança.

func newDecidedCmd() *cobra.Command {
	var root, card, resolucao string
	cmd := &cobra.Command{
		Use:   "decided --card <n> --resolution <revisão>",
		Short: "Libera o card que o `escalate --for-user` parou, depois que a decisão saiu",
		Long: `Fecha o ciclo do ` + "`escalate --for-user`" + `: a decisão saiu, foi promovida a
regra, e o card volta para a fila.

  anchors decided --card 4 --resolution "PLTFR-R0004: o mTLS é construído na F03"

O que faz:
  · remove a label ` + "`" + initx.LabelNeedsUser + "`" + ` do card (é ela que faz o claim pular)
  · comenta a resolução no card, para quem o pegar depois saber o que mudou
  · fecha as issues de decisão abertas SOB este card, com a resolução

A ` + "`--resolution`" + ` é obrigatória, e não é burocracia: a saída de uma decisão em
aberto é UMA — a resposta vira REGRA, com código. Liberar o card sem dizer qual revisão
nasceu dela deixaria a decisão sem rastro, e o próximo a ler o plano não saberia por que
ele diz o que diz.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if card == "" {
				return fmt.Errorf("informe o card com --card (ex: `anchors decided --card 4 --resolution \"ABCDE-R0001: ...\"`)")
			}
			if strings.TrimSpace(resolucao) == "" {
				return fmt.Errorf("informe a --resolution: qual revisão nasceu da decisão " +
					"(ex: `--resolution \"PLTFR-R0004: o mTLS é construído na F03\"`).\n" +
					"   A saída de uma decisão em aberto é UMA: a resposta vira REGRA, com código. " +
					"Sem isso o card volta à fila e a decisão fica sem rastro")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar %s: %w", config.DefaultFile, err)
			}
			if !cfg.GitHubMode() {
				return fmt.Errorf("este comando é do modo github — no modo local a decisão " +
					"mora em `issues/`, e resolvê-la é mover o arquivo para `issues/done/`")
			}

			// A LABEL primeiro: é ela que trava. Se o resto falhar, o card já está livre —
			// a ordem inversa deixaria o card parado com a issue fechada, que é
			// exatamente o estado inconsistente que este comando existe para desfazer.
			out, err := exec.Command("gh", "issue", "edit", card,
				"--repo", cfg.Workflow.Repo,
				"--remove-label", initx.LabelNeedsUser,
			).CombinedOutput()
			if err != nil {
				return fmt.Errorf("liberar o card #%s: %w\n%s", card, err, out)
			}
			fmt.Printf("✓ card #%s liberado — o claim volta a entregá-lo\n", card)

			_ = exec.Command("gh", "issue", "comment", card,
				"--repo", cfg.Workflow.Repo,
				"--body", "▶ Liberado: a decisão saiu e virou regra.\n\n**Resolução:** "+resolucao,
			).Run()

			// As issues de decisão abertas SOB este card. A label `under-<n>` é o que liga
			// as duas pontas — uma frase no corpo ("descoberto durante o card #4") não se
			// consulta, e foi por isso que o `escalate` a usa.
			fechadas := closeDecisionsUnder(cfg.Workflow.Repo, card, resolucao)
			switch fechadas {
			case 0:
				fmt.Printf("  (nenhuma issue de decisão aberta sob o card #%s — "+
					"se havia, ela já estava fechada)\n", card)
			case 1:
				fmt.Println("  1 issue de decisão fechada com a resolução")
			default:
				fmt.Printf("  %d issues de decisão fechadas com a resolução\n", fechadas)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&card, "card", "", "o card que o `escalate --for-user` parou")
	cmd.Flags().StringVar(&resolucao, "resolution", "",
		"a revisão que nasceu da decisão (ex: \"PLTFR-R0004: o mTLS é construído na F03\")")
	return cmd
}

// closeDecisionsUnder fecha as issues de decisão abertas sob um card, e devolve quantas.
//
// Filtra por `needs-user` E pela label `under-<n>`: a primeira diz que é decisão, a
// segunda que é DESTE card. Sem a segunda, o comando fecharia decisões de outro trabalho
// — e uma decisão fechada por engano volta a ser tomada por omissão.
func closeDecisionsUnder(repo, card, resolucao string) int {
	out, err := exec.Command("gh", "issue", "list",
		"--repo", repo,
		"--state", "open",
		"--label", initx.LabelNeedsUser,
		"--label", initx.LabelSob(card),
		"--json", "number",
	).Output()
	if err != nil {
		return 0
	}
	var achadas []struct{ Number int }
	if json.Unmarshal(out, &achadas) != nil {
		return 0
	}
	n := 0
	for _, i := range achadas {
		if exec.Command("gh", "issue", "close", fmt.Sprint(i.Number),
			"--repo", repo,
			"--comment", "✓ Decidido, e a resposta virou regra.\n\n**Resolução:** "+resolucao,
		).Run() == nil {
			n++
		}
	}
	return n
}

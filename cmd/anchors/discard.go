package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

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
		Short: "Tira do board o card que não faz mais sentido, sem apagá-lo",
		Long: `Marca um ou mais cards como DESCARTADOS.

O card sai do board — colunas, árvore, roadmap, faixa de pendências — e continua no
GitHub com a razão registrada. É para o que não tem mais sentido, não para o que
terminou: o trabalho entregue FECHA, e o roadmap precisa dele para desenhar o passado.

Os casos que o motivaram, no projeto de referência:

  · cards de teste do próprio fluxo ("[teste] trava de estado — apagar")
  · achados sobre arquivos que não existem mais (` + "`Orfa.spec.md`, `ZzTeste.spec.md`" + `)
  · perguntas que o projeto respondeu por outro caminho

NÃO use para trabalho que ficou pronto — isso é fechar. Nem para o que ainda importa e
ninguém pegou: um card esquecido continua sendo trabalho, e escondê-lo do board é a
forma mais silenciosa de perdê-lo.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if strings.TrimSpace(motivo) == "" {
				return fmt.Errorf("`--reason` é obrigatório: um card que some do board sem " +
					"dizer por quê é indistinguível de um que se perdeu")
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
				return fmt.Errorf("`discard` existe no modo github (o card é uma issue). " +
					"No modo local, apague o arquivo do card")
			}

			// A LABEL É CRIADA SOB DEMANDA. O `doctor --fix` a garante, mas quem roda o
			// `discard` num projeto que ainda não rodou o doctor não deveria ser barrado
			// por isso — e `gh issue edit` com label inexistente falha o comando INTEIRO.
			_ = exec.Command("gh", "label", "create", initx.LabelDiscarded,
				"--repo", cfg.Workflow.Repo,
				"--color", "d4d4d4",
				"--description", "fora do board: não faz mais sentido",
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
				fmt.Printf("· #%s descartado — fora do board, ainda no GitHub\n", card)
			}
			if len(falhas) > 0 {
				return fmt.Errorf("%d card(s) não foram descartados:\n  %s",
					len(falhas), strings.Join(falhas, "\n  "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&motivo, "reason", "", "por que este card não faz mais sentido")
	aliasDeFlag(cmd, "reason", "motivo")
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

	corpo := "🗑 **Descartado** — este card saiu do board.\n\n" + motivo +
		"\n\nEle continua aqui: descartar não apaga, porque o rastro de que a pergunta " +
		"existiu tem valor. Para trazê-lo de volta, remova a label `" +
		initx.LabelDiscarded + "`."
	if out, err := exec.Command("gh", "issue", "comment", card,
		"--repo", repo, "--body", corpo).CombinedOutput(); err != nil {
		return fmt.Errorf("registrar a razão: %w: %s", err, strings.TrimSpace(string(out)))
	}

	// FECHAR VEM POR ÚLTIMO, e só se ainda estiver aberto: um card descartado que
	// continua aberto seria servido pelo `claim`.
	_ = exec.Command("gh", "issue", "close", card, "--repo", repo).Run()
	return nil
}

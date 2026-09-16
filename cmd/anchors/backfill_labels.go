package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- o vínculo que já existe nas issues abertas ---
//
// Uma label nova não alcança o passado. O `escalate --for-user` passou a marcar o card de
// origem com `anchors:blocked-by-<n>`, e os cards que pararam ANTES disso só têm o
// `needs-user`: o board diz que esperam e não diz por quem, e o `claim` não tem como
// conferir se a decisão saiu.
//
// MEDIDO no projeto de referência: 29 decisões abertas, e destravar as dependentes exigia
// reler card por card para descobrir quem esperava o quê.
//
// DE ONDE SAI O VÍNCULO, e por que ele é recuperável: o `escalate` sempre escreveu o
// número do card de origem na label `anchors:under-<n>` do card NOVO. Então a relação
// existe — na direção contrária. Este comando a lê e escreve a que faltava.
//
// O QUE ELE NÃO FAZ: adivinhar. Card em `needs-user` sem nenhum `under-<n>` apontando para
// ele fica como está, e o comando diz quantos foram esses. Inventar um bloqueador seria
// pior que não ter a label: o `claim` passaria a segurar um card por causa de uma decisão
// que ninguém ligou a ele.

func newBackfillLabelsCmd() *cobra.Command {
	var root string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "backfill-labels",
		Short: "Escreve nas issues abertas os vínculos que a versão nova passou a marcar",
		Long: `Recupera o vínculo ` + "`anchors:blocked-by-<n>`" + ` nos cards que pararam antes de
ele existir.

O ` + "`escalate --for-user`" + ` marca o card de origem com ` + "`blocked-by-<n>`" + `, e os que
pararam antes só têm o ` + "`needs-user`" + `: o board diz que esperam e não diz por quem.

O vínculo é recuperável porque o card NOVO sempre carregou ` + "`anchors:under-<n>`" + ` com o
número do card de origem — a relação existe na direção contrária, e este comando escreve a
que faltava.

Card em ` + "`needs-user`" + ` sem ninguém apontando para ele NÃO é tocado: inventar um
bloqueador faria o claim segurar o card por uma decisão que ninguém ligou a ele.

    anchors backfill-labels --dry-run    # mostra o que faria
    anchors backfill-labels`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(absRoot + "/" + config.DefaultFile)
			if err != nil {
				return err
			}
			if !cfg.GitHubMode() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`backfill-labels` existe no modo github: no modo local " +
					"o vínculo vive na pasta de `issues/`, não em label")
			}
			repo := cfg.Workflow.Repo

			// TODAS as decisões abertas, com as labels. Uma chamada, não uma por card: o
			// limite secundário da API derruba a sequência longa, e já derrubou aqui.
			out, err := exec.Command("gh", "issue", "list",
				"--repo", repo,
				"--state", "open",
				"--label", initx.LabelNeedsUser,
				"--limit", "200",
				"--json", "number,labels",
			).Output()
			if err != nil {
				return fmt.Errorf("listar as decisões abertas: %w", err)
			}
			var decisoes []struct {
				Number int `json:"number"`
				Labels []struct {
					Name string `json:"name"`
				} `json:"labels"`
			}
			if err := json.Unmarshal(out, &decisoes); err != nil {
				return fmt.Errorf("ler a lista: %w", err)
			}

			// O MAPA INVERSO: para cada card de origem, qual decisão o segura.
			//
			// Um card pode ter várias decisões apontando para ele — dois achados no mesmo
			// trabalho. Todas entram: o `claim` segura enquanto QUALQUER uma estiver
			// aberta, e omitir as demais faria o card ser servido quando a primeira
			// fechasse.
			segurados := map[string][]string{}
			for _, d := range decisoes {
				for _, l := range d.Labels {
					origem := ""
					switch {
					case strings.HasPrefix(l.Name, initx.PrefixoLabelSob):
						origem = strings.TrimPrefix(l.Name, initx.PrefixoLabelSob)
					case strings.HasPrefix(l.Name, initx.PrefixoLabelSobAntigo):
						origem = strings.TrimPrefix(l.Name, initx.PrefixoLabelSobAntigo)
					}
					if origem == "" {
						continue
					}
					n := fmt.Sprintf("%d", d.Number)
					// A DECISÃO NÃO SEGURA A SI MESMA. Um card pode ser ao mesmo tempo
					// decisão e origem de outro achado, e ligá-lo a si o tornaria
					// eternamente bloqueado.
					if origem != n {
						segurados[origem] = append(segurados[origem], n)
					}
				}
			}

			if len(segurados) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(),
					"nenhum vínculo a recuperar: %d decisão(ões) aberta(s), e nenhuma "+
						"aponta para um card de origem\n", len(decisoes))
				return nil
			}

			escritos, pulados := 0, 0
			for origem, donos := range segurados {
				// SÓ SE O CARD DE ORIGEM ESTÁ ABERTO: pôr bloqueio em card fechado não
				// muda nada e enche o histórico.
				st, err := exec.Command("gh", "issue", "view", origem,
					"--repo", repo, "--json", "state,labels").Output()
				if err != nil {
					pulados++
					continue
				}
				var estado struct {
					State  string `json:"state"`
					Labels []struct {
						Name string `json:"name"`
					} `json:"labels"`
				}
				if err := json.Unmarshal(st, &estado); err != nil || estado.State != "OPEN" {
					pulados++
					continue
				}
				jaTem := map[string]bool{}
				for _, l := range estado.Labels {
					jaTem[l.Name] = true
				}
				for _, dono := range donos {
					rotulo := initx.LabelBlockedBy(dono)
					if jaTem[rotulo] {
						continue
					}
					if dryRun {
						fmt.Fprintf(cmd.OutOrStdout(), "· #%s receberia `%s`\n", origem, rotulo)
						escritos++
						continue
					}
					_ = exec.Command("gh", "label", "create", rotulo,
						"--repo", repo,
						"--color", "b60205",
						"--description", "este card espera a decisão do #"+dono,
					).Run()
					if o, err := exec.Command("gh", "issue", "edit", origem,
						"--repo", repo, "--add-label", rotulo,
					).CombinedOutput(); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "· #%s: %v — %s\n",
							origem, err, strings.TrimSpace(string(o)))
						pulados++
						continue
					}
					fmt.Fprintf(cmd.OutOrStdout(), "· #%s ← bloqueado pelo #%s\n", origem, dono)
					escritos++
				}
			}
			verbo := "escrito(s)"
			if dryRun {
				verbo = "a escrever (nada foi tocado)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d vínculo(s) %s · %d pulado(s)\n",
				escritos, verbo, pulados)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "mostra o que faria, sem tocar em nada")
	return cmd
}

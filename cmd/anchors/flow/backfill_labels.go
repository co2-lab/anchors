package flow

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
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
		Short: "Write into open issues the links the new version began to mark",
		Long: `Recovers the link ` + "`anchors:blocked-by-<n>`" + ` on the cards that stopped before
it existed.

The ` + "`escalate --for-user`" + ` marks the origin card with ` + "`blocked-by-<n>`" + `, and the ones that
stopped before only have the ` + "`needs-user`" + `: the board says they are waiting and does not say for whom.

The link is recoverable because the NEW card always carried ` + "`anchors:under-<n>`" + ` with the
number of the origin card — the relation exists in the opposite direction, and this
command writes the one that was missing.

A card in ` + "`needs-user`" + ` with nobody pointing at it is NOT touched: inventing a
blocker would make the claim hold the card for a decision nobody linked to it.

    anchors backfill-labels --dry-run    # shows what it would do
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
				return fmt.Errorf("`backfill-labels` exists in github mode: in local mode " +
					"the link lives in the `issues/` folder, not in a label")
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
				return fmt.Errorf("list the open decisions: %w", err)
			}
			var decisoes []struct {
				Number int `json:"number"`
				Labels []struct {
					Name string `json:"name"`
				} `json:"labels"`
			}
			if err := json.Unmarshal(out, &decisoes); err != nil {
				return fmt.Errorf("read the list: %w", err)
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
					if strings.HasPrefix(l.Name, initx.PrefixoLabelSob) {
						origem = strings.TrimPrefix(l.Name, initx.PrefixoLabelSob)
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

			// A PROCEDÊNCIA DO PR REVISADO, para quem nasceu sem card.
			//
			// Um achado que nasce revisando trabalho alheio não tem `anchors-owner`, e até
			// a v0.1.139 nascia sem label nenhuma — citando o PR em PROSA, que é o que a
			// doutrina do `under-<n>` condena.
			//
			// MEDIDO no projeto de referência: 23 decisões sem `under-`, e várias citando
			// "PR #N" no corpo. O vínculo existe, escrito à mão, e não se consulta.
			//
			// O QUE ESTE PASSO FAZ: lê o PR citado, deriva o card do `Refs`/`Closes` dele,
			// e grava as DUAS labels — `from-pr-<n>` (onde foi visto) e `under-<card>`
			// (a que trabalho pertence).
			//
			// E NÃO ADIVINHA: só a primeira menção de `PR #N` no corpo, e só quando o PR
			// de fato declara um card. Um número solto na prosa não vira vínculo — é o
			// mesmo cuidado que o `pr-checks` toma ao ler o `Refs` só no início da linha.
			recuperadosDePR := 0
			for _, d := range decisoes {
				temUnder := false
				for _, l := range d.Labels {
					if strings.HasPrefix(l.Name, initx.PrefixoLabelSob) {
						temUnder = true
						break
					}
				}
				if temUnder {
					continue
				}
				corpo, err := exec.Command("gh", "issue", "view", fmt.Sprintf("%d", d.Number),
					"--repo", repo, "--json", "body", "--jq", ".body // \"\"").Output()
				if err != nil {
					continue
				}
				m := prMentionRE.FindStringSubmatch(string(corpo))
				if m == nil {
					continue
				}
				prCitado := m[1]
				cardDoPr := cardDoPR(repo, prCitado)
				if cardDoPr == "" {
					continue
				}
				rotulos := []string{initx.LabelDePR(prCitado), initx.LabelSob(cardDoPr)}
				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(),
						"· #%d would receive `%s` and `%s` (PR cited in the body)\n",
						d.Number, rotulos[0], rotulos[1])
					recuperadosDePR++
					continue
				}
				for _, r := range rotulos {
					cor, desc := "d4c5f9", "finding seen while reviewing PR #"+prCitado
					if strings.HasPrefix(r, initx.PrefixoLabelSob) {
						cor, desc = "c5def5", "finding born during the work of card #"+cardDoPr
					}
					_ = exec.Command("gh", "label", "create", r,
						"--repo", repo, "--color", cor, "--description", desc).Run()
				}
				if o, err := exec.Command("gh", "issue", "edit", fmt.Sprintf("%d", d.Number),
					"--repo", repo, "--add-label", strings.Join(rotulos, ","),
				).CombinedOutput(); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "· #%d: %v — %s\n",
						d.Number, err, strings.TrimSpace(string(o)))
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "· #%d ← PR #%s, card #%s\n",
					d.Number, prCitado, cardDoPr)
				recuperadosDePR++
			}
			if recuperadosDePR > 0 {
				verbo := "recovered from the cited PR"
				if dryRun {
					verbo = "to recover from the cited PR (nothing was touched)"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\n%d provenance(s) %s\n", recuperadosDePR, verbo)
			}

			if len(segurados) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(),
					"no link to recover: %d open decision(s), and none "+
						"points to an origin card\n", len(decisoes))
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

				// O BACKFILL NÃO INFERE PRECEDÊNCIA ENTRE DECISÕES.
				//
				// `needs-user` + `blocked-by-<n>` juntas são LEGÍTIMAS: uma decisão pode
				// depender de outra, e isso é ordem de corretude — "responda o #701
				// primeiro, porque a resposta dele condiciona a do #647". Proibir a
				// combinação apagaria informação real.
				//
				// O QUE ESTE COMANDO NÃO PODE É ADIVINHAR QUAL DAS DUAS É. Ele lê o
				// `under-<n>`, que significa "nasci do trabalho do #n" — procedência, não
				// precedência. Quando o #n é um card de TRABALHO, o vínculo é direto: o
				// trabalho para até a decisão sair. Quando o #n é OUTRA DECISÃO, `under`
				// não diz se uma condiciona a outra ou se as duas só nasceram no mesmo
				// lugar, e são perguntas diferentes.
				//
				// MEDIDO no projeto de referência: quatro cards (#647, #706, #721, #555)
				// receberam bloqueio por esta inferência, e nos quatro o "bloqueador" era
				// outra decisão que apenas nasceu junto — não precedência. O único vínculo
				// que o `under` acertou foi o #198 → #647: card em `in-review`, trabalho
				// de verdade parado.
				//
				// ENTÃO O BACKFILL RECUPERA SÓ O CASO INEQUÍVOCO, e quem quiser declarar
				// precedência entre decisões põe a label à mão — é julgamento de quem lê
				// as duas, do mesmo jeito que decidir se são o mesmo assunto.
				if jaTem[initx.LabelNeedsUser] {
					fmt.Fprintf(cmd.OutOrStdout(),
						"· #%s is also an open decision — precedence between decisions is not "+
							"inferred from `under-<n>`; if #%s conditions this one, put the label by hand\n",
						origem, strings.Join(donos, ", #"))
					pulados++
					continue
				}
				for _, dono := range donos {
					rotulo := initx.LabelBlockedBy(dono)
					if jaTem[rotulo] {
						continue
					}
					if dryRun {
						fmt.Fprintf(cmd.OutOrStdout(), "· #%s would receive `%s`\n", origem, rotulo)
						escritos++
						continue
					}
					_ = exec.Command("gh", "label", "create", rotulo,
						"--repo", repo,
						"--color", "b60205",
						"--description", "this card is waiting for the decision of #"+dono,
					).Run()
					if o, err := exec.Command("gh", "issue", "edit", origem,
						"--repo", repo, "--add-label", rotulo,
					).CombinedOutput(); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "· #%s: %v — %s\n",
							origem, err, strings.TrimSpace(string(o)))
						pulados++
						continue
					}
					fmt.Fprintf(cmd.OutOrStdout(), "· #%s ← blocked by #%s\n", origem, dono)
					escritos++
				}
			}
			verbo := "written"
			if dryRun {
				verbo = "to write (nothing was touched)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d link(s) %s · %d skipped\n",
				escritos, verbo, pulados)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what it would do, without touching anything")
	return cmd
}

// prMentionRE casa a primeira menção a um PR no corpo do card.
//
// `PR #N` e não só `#N`: um card cita muitos números — outros cards, commits, regras — e
// só a forma com a palavra diz que aquele é um pull request. É o mínimo para não
// transformar prosa em vínculo.
var prMentionRE = regexp.MustCompile(`(?i)\bPR #(\d+)`)

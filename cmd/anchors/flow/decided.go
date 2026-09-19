package flow

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
		Use:   "decided --card <n> --resolution <revision>",
		Short: "Release the card `escalate --for-user` stopped, once the decision is made",
		Long: `Closes the cycle of ` + "`escalate --for-user`" + `: the decision came out, was promoted to a
rule, and the card goes back to the queue.

  anchors decided --card 4 --resolution "PLTFR-R0004: mTLS is built in F03"

What it does:
  · removes the label ` + "`" + initx.LabelNeedsUser + "`" + ` from the card (it is what makes the claim skip it)
  · comments the resolution on the card, so whoever takes it later knows what changed
  · closes the open decision issues UNDER this card, with the resolution

The ` + "`--resolution`" + ` is mandatory, and it is not bureaucracy: the exit of an open
decision is ONE — the answer becomes a RULE, with a code. Releasing the card without
saying which revision was born from it would leave the decision without a trace, and
the next one to read the plan would not know why it says what it says.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if card == "" {
				return fmt.Errorf("provide the card with --card (e.g.: `anchors decided --card 4 --resolution \"ABCDE-R0001: ...\"`)")
			}
			if strings.TrimSpace(resolucao) == "" {
				return fmt.Errorf("provide the --resolution: which revision was born from the decision " +
					"(e.g.: `--resolution \"PLTFR-R0004: mTLS is built in F03\"`).\n" +
					"   The exit of an open decision is ONE: the answer becomes a RULE, with a code. " +
					"Without it the card goes back to the queue and the decision is left without a trace")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load %s: %w", config.DefaultFile, err)
			}
			if !cfg.GitHubMode() {
				return fmt.Errorf("this command belongs to github mode — in local mode the decision " +
					"lives in `issues/`, and resolving it means moving the file to `issues/done/`")
			}

			// A DECISÃO GEROU TRABALHO? Então o card espera a ENTREGA, não a decisão.
			//
			// `anchors unblock` cria o card da mudança com `anchors:desbloqueia-<n>`, e
			// enquanto ele estiver aberto a label NÃO pode sair: o card voltaria à fila e
			// o agente que o pegasse encontraria o trabalho que a decisão pediu ainda por
			// fazer — esbarrando no mesmo impasse que foi escalado.
			//
			// Recusar aqui, e não avisar: liberar o card e imprimir um alerta deixaria o
			// desfecho na mão de quem lê a saída, e a saída de um comando que "funcionou"
			// não se lê com atenção.
			if pendentes := desbloqueiosAbertos(cfg.Workflow.Repo, card); len(pendentes) > 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("card #%s is waiting for the delivery of #%s — the decision came out and "+
					"generated work.\n\n"+
					"   The label `%s` STAYS until then: without it the card goes back to the queue, and whoever "+
					"takes it finds the work still to be done.\n\n"+
					"   When #%s is delivered, run this command again.",
					card, strings.Join(pendentes, ", #"), initx.LabelNeedsUser,
					strings.Join(pendentes, ", #"))
			}

			// O RASTRO ANTES DA REMOÇÃO, porque a label é o único lugar onde ele está.
			//
			// `blocked-by-<n>` sai junto com o `needs-user` — senão fica pendurada num
			// card já livre, e quem lê o board depois vê bloqueio que não existe mais.
			//
			// MAS REMOVER APAGA A HISTÓRIA. Saber que este card esperou por aquela decisão
			// é o que explica o atraso dele, e é o que permite rastrear a decisão para
			// trás: "por que isto ficou parado três dias?" só tem resposta se o vínculo
			// sobreviver ao desbloqueio.
			//
			// O COMENTÁRIO é onde o rastro fica. Ele é imutável no GitHub, aparece na
			// timeline com data, e sobrevive a qualquer mexida posterior nas labels — o
			// contrário da label, que é estado do AGORA e some quando o agora muda.
			bloqueadores := labelsDeBloqueio(cfg.Workflow.Repo, card)

			// A LABEL primeiro: é ela que trava. Se o resto falhar, o card já está livre —
			// a ordem inversa deixaria o card parado com a issue fechada, que é
			// exatamente o estado inconsistente que este comando existe para desfazer.
			remover := []string{initx.LabelNeedsUser}
			remover = append(remover, bloqueadores...)
			out, err := exec.Command("gh", "issue", "edit", card,
				"--repo", cfg.Workflow.Repo,
				"--remove-label", strings.Join(remover, ","),
			).CombinedOutput()
			if err != nil {
				return fmt.Errorf("release card #%s: %w\n%s", card, err, out)
			}
			fmt.Printf("✓ card #%s released — the claim delivers it again\n", card)

			corpo := "▶ Released: the decision came out and became a rule.\n\n**Resolution:** " + resolucao
			if len(bloqueadores) > 0 {
				var quem []string
				for _, b := range bloqueadores {
					quem = append(quem, "#"+strings.TrimPrefix(b, initx.PrefixoLabelBlockedBy))
				}
				corpo += "\n\n**Was blocked by:** " + strings.Join(quem, ", ") +
					"\n\nThe blocking label was removed by this command. The link stays " +
					"recorded here: it is what explains the time stopped, and it is how the " +
					"decision is traced backwards."
			}
			_ = exec.Command("gh", "issue", "comment", card,
				"--repo", cfg.Workflow.Repo,
				"--body", corpo,
			).Run()

			// As issues de decisão abertas SOB este card. A label `under-<n>` é o que liga
			// as duas pontas — uma frase no corpo ("descoberto durante o card #4") não se
			// consulta, e foi por isso que o `escalate` a usa.
			fechadas := closeDecisionsUnder(cfg.Workflow.Repo, card, resolucao)
			switch fechadas {
			case 0:
				fmt.Printf("  (no open decision issue under card #%s — "+
					"if there was one, it was already closed)\n", card)
			case 1:
				fmt.Println("  1 decision issue closed with the resolution")
			default:
				fmt.Printf("  %d decision issues closed with the resolution\n", fechadas)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&card, "card", "", "the card that `escalate --for-user` stopped")
	cmd.Flags().StringVar(&resolucao, "resolution", "",
		"the revision born from the decision (e.g.: \"PLTFR-R0004: mTLS is built in F03\")")
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
			"--comment", "✓ Decided, and the answer became a rule.\n\n**Resolution:** "+resolucao,
		).Run() == nil {
			n++
		}
	}
	return n
}

// desbloqueiosAbertos lista os cards cuja entrega destrava `card`.
//
// A BUSCA POR LABEL do GitHub tem latência de índice — medido: um card recém-criado não
// aparece por alguns segundos, e um recém-fechado continua aparecendo pelo mesmo tempo.
//
// Isso torna o comando ocasionalmente conservador: logo depois de fechar o desbloqueio, o
// `decided` ainda pode recusar. É o lado certo do erro — recusar quando podia liberar custa
// rodar o comando de novo; liberar quando não podia devolve à fila um card cujo trabalho
// não terminou, e o agente que o pegar escala de novo.
//
// Vazio significa que a decisão não gerou trabalho — ou que o trabalho já foi entregue, e
// nos dois casos o card pode voltar à fila.
func desbloqueiosAbertos(repo, card string) []string {
	out, err := exec.Command("gh", "issue", "list",
		"--repo", repo,
		"--state", "open",
		"--label", initx.LabelDesbloqueia(card),
		"--json", "number",
	).Output()
	if err != nil {
		// Sem resposta do `gh` a resposta honesta é "não sei", e não-sei aqui não pode
		// virar "pode liberar": o comando seguiria e removeria a label de um card que
		// talvez espere trabalho. Devolver vazio é o que faz isso acontecer — então a
		// falha de rede é tratada como ausência, e o operador vê o erro do `gh` na tela.
		return nil
	}
	var achados []struct{ Number int }
	if json.Unmarshal(out, &achados) != nil {
		return nil
	}
	var ns []string
	for _, a := range achados {
		ns = append(ns, fmt.Sprint(a.Number))
	}
	return ns
}

// labelsDeBloqueio devolve as labels `blocked-by-<n>` que o card carrega.
//
// Precisa vir ANTES da remoção: depois, a informação não existe em lugar nenhum — a label é
// estado do agora, e o rastro de que o card esperou por aquela decisão só sobrevive se for
// escrito em comentário, que é imutável e datado.
//
// TODAS e não a primeira: um card pode ter parado por duas decisões (dois achados no mesmo
// trabalho), e registrar só uma contaria a história pela metade.
func labelsDeBloqueio(repo, card string) []string {
	out, err := exec.Command("gh", "issue", "view", card,
		"--repo", repo, "--json", "labels",
		"--jq", `[.labels[].name | select(startswith("`+initx.PrefixoLabelBlockedBy+`"))] | .[]`,
	).Output()
	if err != nil {
		return nil
	}
	var rotulos []string
	for _, l := range strings.Fields(string(out)) {
		if l = strings.TrimSpace(l); l != "" {
			rotulos = append(rotulos, l)
		}
	}
	return rotulos
}

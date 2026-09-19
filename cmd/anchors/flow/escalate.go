package flow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- o plano precisa mudar: abrir a issue certa ---
//
// Quem implementa é quem descobre que o plano precisa mudar. E o gatilho não é só
// incoerência: o plano pode estar perfeitamente coerente e INCOMPLETO — a falta não
// aparece no texto, porque o que falta não está escrito em lugar nenhum. Medido: o plano
// da fundação não cobria configuração e execução de migrations, e ninguém veria isso
// lendo o plano.
//
// Incoerência e lacuna são descobertas diferentes com o MESMO fluxo: quem achou
// interpreta o impacto, e a interpretação escolhe a saída.
//
//	--for-user  a mudança impacta a DIREÇÃO do projeto. Vira decisão de quem o
//	                planejou, com `anchors:needs-user`, e o claim não entrega
//	                o card enquanto ela não sair.
//
//	(padrão)        não impacta a direção. Vira card comum: nasce em `to-do`, entra na
//	                fila, um agente pega. Não para ninguém.
//
// A TERCEIRA saída não é deste comando: quando a correção é trivial E o agente já está
// mexendo naquele arquivo, ele corrige e registra a revisão (`{CODIGO}-R0001`). Abrir
// card para trocar uma palavra seria burocracia.
//
// O que este comando fecha é ERGONOMIA, e aqui ela decide o resultado: enquanto abrir a
// issue certa for mais trabalhoso que corrigir em silêncio, o agente corrige em silêncio.
//
// O escalonamento por EXAUSTÃO (dez revisões sem convergir) já existia no pipeline de
// claim. Este é o por JUÍZO: ninguém está travado, alguém percebeu algo.
func newEscalateCmd() *cobra.Command {
	var root, sobre, card, revisandoPR string
	var paraUsuario, incerto bool
	cmd := &cobra.Command{
		Use:   "escalate <reason>",
		Short: "Open the issue for a change needed in the plan or the spec",
		Long: `You found out that the plan or the spec needs to change — by incoherence (the
text contradicts itself) or by GAP (the plan is coherent and did not cover something).

Whoever found it interprets the impact, and the interpretation chooses the exit:

  --for-user   you ASSERT that the change impacts the DIRECTION of the project.
                   There is more than one defensible answer, and choosing between
                   them changes what the product does. The issue is born with
                   'anchors:needs-user' and the claim does not hand over the card
                   while the decision has not come out.

  --unsure     you DO NOT KNOW whether it impacts. It is born just like
                   '--for-user', and the card says the first question is the
                   framing: if it does not impact, whoever reads it returns it to
                   the queue instead of deciding.

  (default)        it does not impact the direction. It becomes an ordinary card:
                   born in 'to-do', enters the queue, an agent takes it. It stops
                   nobody.

WHY TWO EXITS FOR WHAT STOPS THE CARD: whoever reads 30 open decisions needs to
know what each one asks for. "Decide between A and B" and "check whether this is
yours" are different jobs, and mixing them makes the second cost as much as the first.

HOW TO CHOOSE. The question is not "am I sure?" — it is "is there more than one
defensible answer, and does choosing between them change what the product DOES?".

  Impacts the direction          Does not impact
  ---------------------          ---------------
  which value the app offers     the rule is right and the gate that enforces it is missing
  whether the UI shows X or Y    the text contradicts itself and one of the sides is the right one
  what counts as proof           the waiver came out and the ruler did not come in
  changing something in production   the path written on the card is wrong

If the exit is "fix what is wrong" and nobody would defend the current state,
it is an ordinary card: there is no choice to make, only work.

AND IF YOU DO NOT KNOW, use '--unsure' instead of '--for-user'. The cost of
escalating for safety does not show up for whoever escalates: each card in
'needs-user' waits for a person, and one person decides more slowly than N agents
escalate. Measured in a real project: 11 escalations in an hour, and a triage
concluded that the 23 open ones were SEVEN decisions. Saying "I do not know" is
more honest — and cheaper to read — than asserting an impact you did not measure.

Do not use it for what is trivial AND is in the file you are already editing: there,
fix it and record the revision ('{CODIGO}-R0001: what changed and why'). Opening a
card to change one word is bureaucracy.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return err
			}
			// FALHA CEDO e com a razão: no modo local não há card nem label, e abrir uma
			// issue que ninguém lê seria pior que dizer que não dá.
			// Sem label do fluxo a issue nasceria órfã: o claim filtra por ela, e um card
			// que ninguém enxerga é pior que nenhum card. O `Load` já exige isto no modo
			// github — a conferência aqui é para o caso de a validação mudar.
			if cfg.GitHubMode() && len(cfg.Workflow.Labels) == 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("`workflow.labels` is empty: the issue would be born without the " +
					"label the claim pipeline uses to find it, and would be orphaned")
			}
			if !cfg.GitHubMode() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`escalate` exists in github mode (the card is an issue). " +
					"In local mode, write the doubt in the plan and stop the work — there is no " +
					"shared queue to take the card from")
			}

			// JÁ HÁ ALGUÉM NESTE ARQUIVO? — a pergunta que faltava.
			//
			// Dois agentes entregaram o MESMO trabalho no mesmo dia: um pegou o card do
			// gate para `MetricCard.spec.md` às 11:38, outro encontrou o mesmo problema
			// às 12:02 e abriu um card novo. Os dois PRs acrescentaram a mesma seção ao
			// mesmo documento.
			//
			// O `claim` impede dois agentes de pegarem o mesmo card. Não impedia um
			// agente de CRIAR um card para trabalho que já está em andamento noutro.
			if sobre != "" && len(cfg.Workflow.Labels) > 0 {
				if outros := openCardsAbout(sobre, cfg.Workflow.Labels[0]); len(outros) > 0 {
					fmt.Fprintln(os.Stderr)
					fmt.Fprintln(os.Stderr, i18n.T("escalate.already_open", sobre))
					for i, o := range outros {
						if i == 3 {
							fmt.Fprintln(os.Stderr, i18n.T("escalate.already_open_more", len(outros)-3))
							break
						}
						fmt.Fprintln(os.Stderr, "    "+o)
					}
					fmt.Fprintln(os.Stderr)
				}
			}

			// O CARD DE ORIGEM, DESCOBERTO se não foi informado.
			//
			// O `--card` era opcional e silencioso: sem ele o achado nasce sem procedência,
			// e ninguém consegue perguntar "de onde veio isto?" nem "o que este trabalho
			// gerou?". A doutrina do `under-<n>` diz que o vínculo é LABEL justamente
			// porque frase no corpo não se consulta — e sem o `--card` não há nem frase.
			//
			// MEDIDO no projeto de referência: 18 de 27 decisões abertas sem amarração
			// alguma. Todas traziam `**Onde:** <arquivo>` (procedência pelo ARQUIVO), e
			// uma chegou a escrever "o achado veio do #690" em prosa — exatamente o que a
			// doutrina condena.
			//
			// O `pr-body` já descobria o card pelo `ANCHORS_AGENT` (é o que o
			// `requestedCards` faz), e este comando não. A assimetria não tinha razão: os
			// dois perguntam a mesma coisa — "qual card este agente pegou?".
			// O CARD SAI DO PR QUE SE REVISAVA, quando é esse o caminho.
			//
			// Um achado que nasce revisando trabalho alheio não tem `anchors-owner` — o
			// agente não pegou card nenhum, está lendo o de outro. A descoberta abaixo
			// devolve vazio, e o achado nascia sem procedência: medido nos cards #786 e
			// #788, que citam o PR revisado em PROSA e não têm label nenhuma.
			//
			// O VÍNCULO NÃO É COM O PR, É COM O CARD. Um PR já declara a issue dele
			// (`Refs`/`Closes`), e guardar as duas pontas duplicaria o que a plataforma
			// relaciona: quem abre o PR vê a issue. O que faltava era o comando LER essa
			// declaração em vez de o agente procurar.
			//
			// Conferido nos dois casos reais: o PR #556 declara `Closes #198`, e o #783
			// declara `Refs #735` — que é o card cuja decisão gerou a contradição que o
			// #788 escalou.
			if strings.TrimSpace(card) == "" && strings.TrimSpace(revisandoPR) != "" {
				if n := cardDoPR(cfg.Workflow.Repo, revisandoPR); n != "" {
					card = n
					fmt.Fprintf(os.Stderr,
						"anchors: PR #%s declares card #%s — it is under it that the finding is born\n",
						strings.TrimPrefix(revisandoPR, "#"), card)
				} else {
					fmt.Fprintf(os.Stderr,
						"anchors: PR #%s declares no card (`Refs`/`Closes`) — the finding "+
							"is born without provenance\n", strings.TrimPrefix(revisandoPR, "#"))
				}
			}

			if strings.TrimSpace(card) == "" {
				if achados := requestedCards("", cfg); len(achados) == 1 {
					card = achados[0]
					fmt.Fprintf(os.Stderr,
						"anchors: the finding is born under card #%s (from `anchors-owner`)\n", card)
				} else if len(achados) > 1 {
					// AMBIGUIDADE NÃO SE RESOLVE POR CHUTE: com dois cards em mãos, só o
					// agente sabe em qual estava mexendo quando achou isto.
					cmd.SilenceUsage = true
					return fmt.Errorf("you have %d cards in hand (#%s) — provide `--card <n>` "+
						"to say under which one this finding was born",
						len(achados), strings.Join(achados, ", #"))
				} else {
					// SEM CARD NENHUM o achado nasce solto, e isso é legítimo: alguém pode
					// escalar fora de um trabalho. O que não pode é acontecer sem aviso.
					fmt.Fprintln(os.Stderr,
						"anchors: without `--card` and without `anchors-owner` — the finding is born WITHOUT "+
							"provenance.")
					fmt.Fprintln(os.Stderr,
						"  nobody will be able to ask where it came from, nor what this "+
							"work generated. If it was born from a card, pass `--card <n>`.")
				}
			}

			// `--unsure` PARA O CARD como o `--for-user`, e a diferença está no CORPO.
			//
			// As duas afirmam coisas diferentes. "Impacta a direção" pede uma decisão;
			// "não sei se impacta" pede um juízo sobre o próprio enquadramento — e quem
			// lê trinta decisões abertas precisa saber qual das duas chegou, porque
			// "decida entre A e B" e "confira se isto é seu" são trabalhos distintos.
			//
			// Misturá-los faz o segundo custar como o primeiro: quem lê gasta o esforço
			// de decidir antes de descobrir que só precisava devolver o card à fila.
			paraUsuario = paraUsuario || incerto
			motivo := strings.Join(args, " ")
			corpoTexto := escalationBody(motivo, sobre, card, paraUsuario, incerto)

			tmp, err := os.CreateTemp("", "anchors-escalate-*.md")
			if err != nil {
				return err
			}
			defer os.Remove(tmp.Name())
			if _, err := tmp.WriteString(corpoTexto); err != nil {
				return err
			}
			tmp.Close()

			// AS LABELS decidem quem pega a issue, e errar aqui a torna órfã: sem a
			// label do fluxo E um estado, o pipeline de claim não a enxerga (ele filtra
			// por `--label $LABEL --label anchors:to-do`), e ela fica no repositório sem
			// nunca chegar a ninguém.
			titulo := "[plan] " + firstLineOfReason(motivo)
			labels := []string{cfg.Workflow.Labels[0], "anchors:to-do"}
			if paraUsuario {
				titulo = "[decision] " + firstLineOfReason(motivo)
				labels = append(labels, initx.LabelNeedsUser)
			}
			// A DÚVIDA GANHA TÍTULO E LABEL PRÓPRIOS. Quem abre a fila precisa distinguir
			// de relance "decida entre A e B" de "confira se isto é seu" — a segunda tem
			// saída barata (trocar a label por `to-do`) e não exige decidir o mérito.
			if incerto {
				titulo = "[framing] " + firstLineOfReason(motivo)
				labels = append(labels, initx.LabelNeedsFraming)
			}
			// SOB o card de origem, como LABEL — o que permite listar o que pende sob um
			// trabalho (`--label anchors:under-44`) e entregá-lo no mesmo PR. Uma frase no
			// corpo ("descoberto durante o card #44") não se consulta.
			// O PR REVISADO É PROCEDÊNCIA TAMBÉM, e coexiste com o card.
			//
			// Os dois respondem perguntas diferentes: o `under-<n>` diz a qual trabalho o
			// achado pertence (por onde ele se entrega), e o `from-pr-<n>` diz onde alguém
			// o viu (por onde se rastreia a revisão).
			//
			// NÃO É REDUNDANTE com o `Refs` do PR, embora o card seja DERIVADO dele. Ler o
			// `Refs` é uma consulta no momento do escalate; se o PR for reescrito depois, a
			// declaração muda e o vínculo histórico se perde. A label registra o que foi
			// lido, quando foi lido.
			if revisandoPR != "" {
				n := strings.TrimPrefix(strings.TrimSpace(revisandoPR), "#")
				rotulo := initx.LabelDePR(n)
				// SOB DEMANDA, como as outras de vínculo: uma por PR, e pré-criar todas é
				// impossível. "Já existe" é o caso comum do segundo achado no mesmo PR.
				_ = exec.Command("gh", "label", "create", rotulo,
					"--repo", cfg.Workflow.Repo,
					"--color", "d4c5f9",
					"--description", "finding seen while reviewing PR #"+n,
				).Run()
				labels = append(labels, rotulo)
			}
			if card != "" {
				labels = append(labels, initx.LabelSob(card))
			}
			argv := []string{"issue", "create",
				"--repo", cfg.Workflow.Repo,
				"--title", titulo,
				"--body-file", tmp.Name(),
			}
			for _, l := range labels {
				argv = append(argv, "--label", l)
			}
			// A label `sob-<n>` é criada SOB DEMANDA: ela é uma por card, e pré-criar
			// todas seria impossível. `gh issue create` falha se a label não existe, então
			// a criação vem antes — e o erro é ignorado de propósito, porque "já existe" é
			// o caso comum a partir do segundo achado do mesmo card.
			if card != "" {
				_ = exec.Command("gh", "label", "create", initx.LabelSob(card),
					"--repo", cfg.Workflow.Repo,
					"--color", "c5def5",
					"--description", "finding born during the work of card #"+card,
				).Run()
			}
			out, err := exec.Command("gh", argv...).CombinedOutput()
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("open the issue: %v — %s", err, strings.TrimSpace(string(out)))
			}
			url := strings.TrimSpace(string(out))

			// A palavra acompanha a SAÍDA: "decisão" para o que espera uma pessoa,
			// "achado" para o card comum. Dizer "decisão aberta" nos dois casos fez eu
			// mesmo conferir a label achando que tinha escalado sem querer.
			que := "finding recorded"
			if paraUsuario {
				que = "decision opened"
			}
			fmt.Printf("%s: %s\n", que, url)

			// O CARD só é parado quando a decisão é do usuário. Numa issue comum o
			// trabalho SEGUE: a mudança não impacta a direção, e travar o card seria
			// efeito colateral que ninguém pediu.
			//
			// Quando para, a label é o que impede o agente seguinte de pegar o card e
			// refazer o mesmo caminho até a mesma dúvida.
			if paraUsuario && card != "" {
				// A LABEL DIZ POR QUAL CARD ELE ESPERA, e não só que espera.
				//
				// O `needs-user` sozinho para o card — e não diz QUEM o segura. O board
				// mostrava "esperando você" sem número, e quem responde a decisão não
				// sabe o que acabou de soltar: cada card tem de ser reencontrado à mão, e
				// o que não for reencontrado segue parado depois de a decisão já ter saído.
				//
				// `anchors:blocked-by-<n>` fecha isso nas três pontas: o board desenha o
				// vínculo, o `claim` confere se o #n ainda está aberto antes de servir o
				// card, e quem decide lista tudo que a sua resposta libera
				// (`--label anchors:blocked-by-<n>`).
				//
				// É AUTOMÁTICO porque não há julgamento a fazer: se este card parou por
				// causa daquela decisão, o vínculo é fato, não escolha. O julgamento que
				// existe — se a decisão impede o trabalho — já foi feito quando o agente
				// escolheu `--for-user`.
				// AQUI O VÍNCULO É DECLARADO, e por isso vale mesmo quando o card de
				// origem JÁ era uma decisão aberta.
				//
				// `needs-user` + `blocked-by-<n>` juntas são legítimas: uma decisão pode
				// depender de outra, e isso é ordem de corretude — "responda o #701
				// primeiro, porque a resposta dele condiciona a do #647".
				//
				// E é exatamente o que acontece quando um agente trabalhando num card já
				// escalado descobre que há uma pergunta ANTES daquela: ele escala de novo,
				// e a nova decisão precede a que já estava lá. Suprimir a label nesse caso
				// apagaria a ordem que ele acabou de estabelecer.
				//
				// A DIFERENÇA COM O `backfill` é a fonte. Lá o vínculo é INFERIDO do
				// `under-<n>` (que diz procedência, não precedência) e a inferência erra:
				// medido, quatro cards receberam bloqueio de uma decisão que apenas nasceu
				// junto. Aqui quem escalou está declarando, no ato, que o trabalho para
				// por causa desta decisão.
				rotuloBloqueio := ""
				if n := numeroDaIssue(url); n != "" {
					rotuloBloqueio = initx.LabelBlockedBy(n)
					// SOB DEMANDA, como a `under-<n>`: é uma label por card, e pré-criar
					// todas é impossível. O erro é ignorado porque "já existe" é o caso
					// comum a partir do segundo card que espera a mesma decisão.
					_ = exec.Command("gh", "label", "create", rotuloBloqueio,
						"--repo", cfg.Workflow.Repo,
						"--color", "b60205",
						"--description", "this card is waiting for the decision of #"+n,
					).Run()
				}
				rotulos := initx.LabelNeedsUser
				if rotuloBloqueio != "" {
					rotulos += "," + rotuloBloqueio
				}
				if _, err := exec.Command("gh", "issue", "edit", card,
					"--repo", cfg.Workflow.Repo,
					"--add-label", rotulos,
				).CombinedOutput(); err != nil {
					fmt.Printf("· warning: could not label card #%s — label it by hand, "+
						"otherwise another agent takes the card and redoes the path\n", card)
				} else {
					// O CAMINHO DE VOLTA sai junto com o aviso de parada.
					//
					// Sem isto o comando gravava um estado e não dizia como revertê-lo:
					// medido em blue-eyes#139, a decisão saiu, as revisões foram
					// aplicadas, a issue fechada — e o card ficou parado, porque a
					// instrução de remover a label só existia no corpo do card que o
					// WORKFLOW abre. Descobrir exigia grepar o YAML.
					fmt.Printf("· card #%s stopped until the decision\n", card)
					fmt.Printf("  to resume, after the decision becomes a rule:\n"+
						"    anchors decided --card %s --resolution \"<CODE>-R000N: what changed\"\n", card)
				}
				_ = exec.Command("gh", "issue", "comment", card,
					"--repo", cfg.Workflow.Repo,
					"--body", "⏸ Stopped: there is an open decision — "+url,
				).Run()
			} else if card != "" {
				// Sem parar, mas deixando o rastro: quem for revisar este card precisa
				// saber que a mudança de plano nasceu daqui.
				_ = exec.Command("gh", "issue", "comment", card,
					"--repo", cfg.Workflow.Repo,
					"--body", "📋 Plan change recorded from this work — "+url,
				).Run()
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&sobre, "about", "", "the plan or spec where the incoherence is")
	cmd.Flags().StringVar(&card, "card", "", "card number where the need was discovered")
	cmd.Flags().StringVar(&revisandoPR, "reviewing-pr", "",
		"the PR you were reviewing: the card comes out of its Refs/Closes, without you searching")
	cmd.Flags().BoolVar(&paraUsuario, "for-user", false,
		"you ASSERT the change impacts the project's DIRECTION: it becomes a decision and stops the card")
	cmd.Flags().BoolVar(&incerto, "unsure", false,
		"you DO NOT KNOW whether it impacts the direction: it stops the card, and the first question is the framing")
	common.AliasDeFlag(cmd, "about", "sobre")
	common.AliasDeFlag(cmd, "for-user", "para-usuario")
	cmd.PreRunE = func(c *cobra.Command, _ []string) error {
		return common.ResolveAliases(c, map[string]string{"about": "sobre", "for-user": "para-usuario"})
	}
	return cmd
}

// firstLineOfReason faz o título da issue, que é uma linha.
func firstLineOfReason(s string) string {
	return common.FirstLineOfReason(s)
}

// escalationBody monta o texto da issue.
//
// Separado do comando porque é ELE o que precisa ser confrontado: o valor está em dizer
// por que o trabalho parou e como destravar. Um teste que precisasse do `gh` para ler
// isso não rodaria em máquina nenhuma, e o texto ficaria sem régua.
func escalationBody(motivo, sobre, card string, paraUsuario, incerto bool) string {
	var b strings.Builder
	if incerto {
		b.WriteString("❓ **I do not know whether this decision is mine.**\n\n")
	} else if paraUsuario {
		b.WriteString("🛑 **This decision is not the agent's.**\n\n")
	} else {
		b.WriteString("📋 **The plan needs to change.**\n\n")
	}
	b.WriteString(motivo + "\n\n")
	if sobre != "" {
		b.WriteString(fmt.Sprintf("**Where:** `%s`\n\n", sobre))
	}
	if incerto {
		b.WriteString("**THE FIRST QUESTION IS THE FRAMING, not the merit.** Whoever " +
			"discovered this could not say whether it changes the project's DIRECTION, and preferred " +
			"to declare the doubt rather than assert an impact they did not measure.\n\n")
		b.WriteString("**If it does NOT impact the direction:** swap the label `" + initx.LabelNeedsUser +
			"` for `anchors:to-do` and the card goes back to the queue — an agent takes it. Do not spend the " +
			"effort of deciding the merit: the exit here is \"this is work, not a choice\".\n\n")
		b.WriteString("**If it does impact:** proceed as with any decision — record it where it holds, " +
			"in the plan or in the spec, as a revision (`{CODE}-R0001: what changed and why`), and " +
			"then remove the label `" + initx.LabelNeedsUser + "`.\n\n")
		b.WriteString("> The criterion: is there more than one defensible answer, and does choosing " +
			"between them change what the product DOES? If the exit is \"fix what is " +
			"wrong\" and nobody would defend the current state, it is an ordinary card.\n\n")
		if card != "" {
			b.WriteString(fmt.Sprintf("Work stopped on card #%s.\n", card))
		}
		return b.String()
	}
	if paraUsuario {
		b.WriteString("**Why it stopped here:** whoever discovered it read this change as " +
			"impacting the project's DIRECTION, and that is a decision for whoever planned it. Changing it " +
			"on one's own would make the project walk toward a destination nobody chose — " +
			"and the altered plan would stay valid, so no gate would flag it.\n\n")
		b.WriteString("**How to unblock:** decide, and record the decision where it holds — in the " +
			"plan or in the spec, as a revision (`{CODE}-R0001: what changed and why`). If the " +
			"change is large, a new plan with `revises:`. Then remove the label `" +
			initx.LabelNeedsUser + "`.\n\n")
		if card != "" {
			b.WriteString(fmt.Sprintf("Work stopped on card #%s.\n", card))
		}
		return b.String()
	}
	b.WriteString("**Why it is an ordinary card:** whoever discovered it read the change as NOT " +
		"impacting the project's direction — it is the plan becoming correct, not changing course. " +
		"It enters the queue like any other work.\n\n")
	b.WriteString("**What to do:** change the plan or the spec and record the revision in the " +
		"file itself (`{CODE}-R0001: what changed and why`). If while working you " +
		"conclude that this CHANGES THE DIRECTION, do not proceed: `anchors escalate ... --for-user`.\n\n")
	if card != "" {
		b.WriteString(fmt.Sprintf("Born UNDER card #%s (label `%s`), which proceeds normally. "+
			"The two are delivered in the same PR: the finding appeared while doing that work, and "+
			"separating them would make one of the two wait for no reason.\n", card, initx.LabelSob(card)))
	}
	return b.String()
}

// numeroDaIssue tira o número da URL que o `gh issue create` imprime.
//
// O `gh` devolve a URL completa, e o que se precisa é o número — para montar a label que
// liga os dois cards. Vazio se a saída não tiver a forma esperada: melhor não rotular do
// que rotular com um pedaço de URL.
func numeroDaIssue(url string) string {
	i := strings.LastIndex(url, "/")
	if i < 0 || i+1 >= len(url) {
		return ""
	}
	n := strings.TrimSpace(url[i+1:])
	for _, r := range n {
		if r < '0' || r > '9' {
			return ""
		}
	}
	return n
}

// cardDoPR lê o card que um PR declara no corpo, pelo `Refs`/`Closes`.
//
// É o mesmo vínculo que o `pr-body` escreve e que o `pr-checks` lê para mover o card — a
// diferença é a direção: aqui se pergunta "de qual card é este PR?", para que um achado
// nascido ao revisá-lo herde a procedência certa.
//
// SÓ A PRIMEIRA linha de vínculo. Um PR que entrega vários cards declara vários, e escolher
// entre eles é julgamento de quem leu o achado — o comando não adivinha: devolve o primeiro
// e quem discordar passa `--card`.
func cardDoPR(repo, pr string) string {
	pr = strings.TrimPrefix(strings.TrimSpace(pr), "#")
	if repo == "" || pr == "" {
		return ""
	}
	out, err := exec.Command("gh", "pr", "view", pr,
		"--repo", repo, "--json", "body", "--jq", ".body // \"\"").Output()
	if err != nil {
		return ""
	}
	m := vinculoNoCorpoRE.FindStringSubmatch(string(out))
	if m == nil {
		return ""
	}
	return m[1]
}

// vinculoNoCorpoRE casa a linha de vínculo que o `pr-body` gera e as que se escrevem à mão.
//
// `(?mi)` porque a linha pode estar em qualquer ponto do corpo e a grafia varia. As quatro
// palavras são as que o `pr-checks` também aceita — manter as listas iguais é o que impede
// que um PR seja vínculo para um comando e não para o outro.
var vinculoNoCorpoRE = regexp.MustCompile(`(?mi)^\s*(?:refs|closes|fixes|resolves)\s+#(\d+)`)

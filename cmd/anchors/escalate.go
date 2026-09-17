package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	var root, sobre, card string
	var paraUsuario, incerto bool
	cmd := &cobra.Command{
		Use:   "escalate <motivo>",
		Short: "Abre a issue de uma mudança necessária no plano ou na spec",
		Long: `Você descobriu que o plano ou a spec precisa mudar — por incoerência (o
texto se contradiz) ou por LACUNA (o plano está coerente e não cobriu algo).

Quem descobriu interpreta o impacto, e a interpretação escolhe a saída:

  --for-user   você AFIRMA que a mudança impacta a DIREÇÃO do projeto. Existe
                   mais de uma resposta defensável, e escolher entre elas muda
                   o que o produto faz. A issue nasce com 'anchors:needs-user'
                   e o claim não entrega o card enquanto a decisão não sair.

  --unsure     você NÃO SABE se impacta. Nasce igual ao '--for-user', e o card
                   diz que a primeira pergunta é o enquadramento: se não
                   impacta, quem lê devolve à fila em vez de decidir.

  (padrão)         não impacta a direção. Vira card comum: nasce em 'to-do',
                   entra na fila, um agente pega. Não para ninguém.

POR QUE DUAS SAÍDAS PARA O QUE PARA O CARD: quem lê 30 decisões abertas precisa
saber o que cada uma pede. "Decida entre A e B" e "confira se isto é seu" são
trabalhos diferentes, e misturá-los faz o segundo custar como o primeiro.

COMO ESCOLHER. A pergunta não é "tenho certeza?" — é "existe mais de uma
resposta defensável, e escolher entre elas muda o que o produto FAZ?".

  Impacta a direção              Não impacta
  -----------------              -----------
  qual valor o app oferece       a regra está certa e falta o gate que a cobra
  se a UI mostra X ou Y          o texto se contradiz e um dos lados é o certo
  o que conta como prova         a dispensa saiu e a régua não entrou
  trocar algo já em produção     o caminho escrito no card está errado

Se a saída é "arrumar o que está errado" e ninguém defenderia o estado atual,
é card comum: não há escolha a fazer, só trabalho.

E SE NÃO SOUBER, use '--unsure' em vez de '--for-user'. O custo de escalar por
segurança não aparece para quem escala: cada card em 'needs-user' espera uma
pessoa, e uma pessoa decide mais devagar do que N agentes escalam. Medido num
projeto real: 11 escaladas numa hora, e uma triagem concluiu que as 23 abertas
eram SETE decisões. Dizer "não sei" é mais honesto — e mais barato de ler — do
que afirmar impacto que você não mediu.

Não use para o que é trivial E está no arquivo que você já está editando: aí
corrija e registre a revisão ('{CODIGO}-R0001: o que mudou e por quê'). Abrir
card para trocar uma palavra é burocracia.`,
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
				return fmt.Errorf("`workflow.labels` está vazio: a issue nasceria sem a " +
					"label que o pipeline de claim usa para achá-la, e ficaria órfã")
			}
			if !cfg.GitHubMode() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`escalate` existe no modo github (o card é uma issue). " +
					"No modo local, escreva a dúvida no plano e pare o trabalho — não há " +
					"fila compartilhada de onde tirar o card")
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
			if strings.TrimSpace(card) == "" {
				if achados := requestedCards("", cfg); len(achados) == 1 {
					card = achados[0]
					fmt.Fprintf(os.Stderr,
						"anchors: o achado nasce sob o card #%s (do `anchors-owner`)\n", card)
				} else if len(achados) > 1 {
					// AMBIGUIDADE NÃO SE RESOLVE POR CHUTE: com dois cards em mãos, só o
					// agente sabe em qual estava mexendo quando achou isto.
					cmd.SilenceUsage = true
					return fmt.Errorf("você tem %d cards em mãos (#%s) — informe `--card <n>` "+
						"para dizer sob qual este achado nasceu",
						len(achados), strings.Join(achados, ", #"))
				} else {
					// SEM CARD NENHUM o achado nasce solto, e isso é legítimo: alguém pode
					// escalar fora de um trabalho. O que não pode é acontecer sem aviso.
					fmt.Fprintln(os.Stderr,
						"anchors: sem `--card` e sem `anchors-owner` — o achado nasce SEM "+
							"procedência.")
					fmt.Fprintln(os.Stderr,
						"  ninguém vai poder perguntar de onde ele veio, nem o que este "+
							"trabalho gerou. Se ele nasceu de um card, passe `--card <n>`.")
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
			titulo := "[plano] " + firstLineOfReason(motivo)
			labels := []string{cfg.Workflow.Labels[0], "anchors:to-do"}
			if paraUsuario {
				titulo = "[decisão] " + firstLineOfReason(motivo)
				labels = append(labels, initx.LabelNeedsUser)
			}
			// A DÚVIDA GANHA TÍTULO E LABEL PRÓPRIOS. Quem abre a fila precisa distinguir
			// de relance "decida entre A e B" de "confira se isto é seu" — a segunda tem
			// saída barata (trocar a label por `to-do`) e não exige decidir o mérito.
			if incerto {
				titulo = "[enquadramento] " + firstLineOfReason(motivo)
				labels = append(labels, initx.LabelNeedsFraming)
			}
			// SOB o card de origem, como LABEL — o que permite listar o que pende sob um
			// trabalho (`--label anchors:under-44`) e entregá-lo no mesmo PR. Uma frase no
			// corpo ("descoberto durante o card #44") não se consulta.
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
					"--description", "achado que nasceu durante o trabalho do card #"+card,
				).Run()
			}
			out, err := exec.Command("gh", argv...).CombinedOutput()
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("abrir a issue: %v — %s", err, strings.TrimSpace(string(out)))
			}
			url := strings.TrimSpace(string(out))

			// A palavra acompanha a SAÍDA: "decisão" para o que espera uma pessoa,
			// "achado" para o card comum. Dizer "decisão aberta" nos dois casos fez eu
			// mesmo conferir a label achando que tinha escalado sem querer.
			que := "achado registrado"
			if paraUsuario {
				que = "decisão aberta"
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
						"--description", "este card espera a decisão do #"+n,
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
					fmt.Printf("· aviso: não consegui rotular o card #%s — rotule à mão, "+
						"senão outro agente pega o card e refaz o caminho\n", card)
				} else {
					// O CAMINHO DE VOLTA sai junto com o aviso de parada.
					//
					// Sem isto o comando gravava um estado e não dizia como revertê-lo:
					// medido em blue-eyes#139, a decisão saiu, as revisões foram
					// aplicadas, a issue fechada — e o card ficou parado, porque a
					// instrução de remover a label só existia no corpo do card que o
					// WORKFLOW abre. Descobrir exigia grepar o YAML.
					fmt.Printf("· card #%s parado até a decisão\n", card)
					fmt.Printf("  para retomar, depois que a decisão virar regra:\n"+
						"    anchors decided --card %s --resolution \"<CODIGO>-R000N: o que mudou\"\n", card)
				}
				_ = exec.Command("gh", "issue", "comment", card,
					"--repo", cfg.Workflow.Repo,
					"--body", "⏸ Parado: há uma decisão em aberto — "+url,
				).Run()
			} else if card != "" {
				// Sem parar, mas deixando o rastro: quem for revisar este card precisa
				// saber que a mudança de plano nasceu daqui.
				_ = exec.Command("gh", "issue", "comment", card,
					"--repo", cfg.Workflow.Repo,
					"--body", "📋 Mudança de plano registrada a partir deste trabalho — "+url,
				).Run()
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&sobre, "about", "", "o plano ou spec onde está a incoerência")
	cmd.Flags().StringVar(&card, "card", "", "número do card onde a necessidade foi descoberta")
	cmd.Flags().BoolVar(&paraUsuario, "for-user", false,
		"você AFIRMA que a mudança impacta a DIREÇÃO do projeto: vira decisão e para o card")
	cmd.Flags().BoolVar(&incerto, "unsure", false,
		"você NÃO SABE se impacta a direção: para o card, e a primeira pergunta é o enquadramento")
	aliasDeFlag(cmd, "about", "sobre")
	aliasDeFlag(cmd, "for-user", "para-usuario")
	cmd.PreRunE = func(c *cobra.Command, _ []string) error {
		return resolveAliases(c, map[string]string{"about": "sobre", "for-user": "para-usuario"})
	}
	return cmd
}

// firstLineOfReason faz o título da issue, que é uma linha.
func firstLineOfReason(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if len(s) > 70 {
		return s[:67] + "..."
	}
	return s
}

// escalationBody monta o texto da issue.
//
// Separado do comando porque é ELE o que precisa ser confrontado: o valor está em dizer
// por que o trabalho parou e como destravar. Um teste que precisasse do `gh` para ler
// isso não rodaria em máquina nenhuma, e o texto ficaria sem régua.
func escalationBody(motivo, sobre, card string, paraUsuario, incerto bool) string {
	var b strings.Builder
	if incerto {
		b.WriteString("❓ **Não sei se esta decisão é minha.**\n\n")
	} else if paraUsuario {
		b.WriteString("🛑 **Esta decisão não é do agente.**\n\n")
	} else {
		b.WriteString("📋 **O plano precisa mudar.**\n\n")
	}
	b.WriteString(motivo + "\n\n")
	if sobre != "" {
		b.WriteString(fmt.Sprintf("**Onde:** `%s`\n\n", sobre))
	}
	if incerto {
		b.WriteString("**A PRIMEIRA PERGUNTA É O ENQUADRAMENTO, não o mérito.** Quem " +
			"descobriu isto não soube dizer se muda a DIREÇÃO do projeto, e preferiu " +
			"declarar a dúvida a afirmar um impacto que não mediu.\n\n")
		b.WriteString("**Se NÃO impacta a direção:** troque a label `" + initx.LabelNeedsUser +
			"` por `anchors:to-do` e o card volta à fila — um agente pega. Não gaste o " +
			"esforço de decidir o mérito: a saída aqui é \"isto é trabalho, não escolha\".\n\n")
		b.WriteString("**Se impacta:** siga como qualquer decisão — registre onde ela vale, " +
			"no plano ou na spec, como revisão (`{CODIGO}-R0001: o que mudou e por quê`), e " +
			"depois remova a label `" + initx.LabelNeedsUser + "`.\n\n")
		b.WriteString("> O critério: existe mais de uma resposta defensável, e escolher " +
			"entre elas muda o que o produto FAZ? Se a saída é \"arrumar o que está " +
			"errado\" e ninguém defenderia o estado atual, é card comum.\n\n")
		if card != "" {
			b.WriteString(fmt.Sprintf("Trabalho parado no card #%s.\n", card))
		}
		return b.String()
	}
	if paraUsuario {
		b.WriteString("**Por que parou aqui:** quem descobriu interpretou que esta mudança " +
			"impacta a DIREÇÃO do projeto, e isso é decisão de quem o planejou. Mudar por " +
			"conta própria faria o projeto caminhar para um destino que ninguém escolheu — " +
			"e o plano alterado ficaria válido, então nenhum gate acusaria.\n\n")
		b.WriteString("**Como destravar:** decida, e registre a decisão onde ela vale — no " +
			"plano ou na spec, como revisão (`{CODIGO}-R0001: o que mudou e por quê`). Se a " +
			"mudança for grande, um plano novo com `revises:`. Depois remova a label `" +
			initx.LabelNeedsUser + "`.\n\n")
		if card != "" {
			b.WriteString(fmt.Sprintf("Trabalho parado no card #%s.\n", card))
		}
		return b.String()
	}
	b.WriteString("**Por que é card comum:** quem descobriu interpretou que a mudança NÃO " +
		"impacta a direção do projeto — é o plano ficando correto, não mudando de rumo. " +
		"Entra na fila como qualquer outro trabalho.\n\n")
	b.WriteString("**O que fazer:** altere o plano ou a spec e registre a revisão no " +
		"próprio arquivo (`{CODIGO}-R0001: o que mudou e por quê`). Se ao mexer você " +
		"concluir que isto MUDA A DIREÇÃO, não siga: `anchors escalate ... --for-user`.\n\n")
	if card != "" {
		b.WriteString(fmt.Sprintf("Nasceu SOB o card #%s (label `%s`), que segue normalmente. "+
			"Os dois se entregam no mesmo PR: o achado apareceu fazendo aquele trabalho, e "+
			"separá-los faria um dos dois esperar sem razão.\n", card, initx.LabelSob(card)))
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

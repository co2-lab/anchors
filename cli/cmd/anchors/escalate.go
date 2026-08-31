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
//	--para-usuario  a mudança impacta a DIREÇÃO do projeto. Vira decisão de quem o
//	                planejou, com `anchors:precisa-do-usuario`, e o claim não entrega
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
	var paraUsuario bool
	cmd := &cobra.Command{
		Use:   "escalate <motivo>",
		Short: "Abre a issue de uma mudança necessária no plano ou na spec",
		Long: `Você descobriu que o plano ou a spec precisa mudar — por incoerência (o
texto se contradiz) ou por LACUNA (o plano está coerente e não cobriu algo).

Quem descobriu interpreta o impacto, e a interpretação escolhe a saída:

  --para-usuario   a mudança impacta a DIREÇÃO do projeto, ou você tem dúvida se
                   impacta. Vira decisão de quem planejou: a issue nasce com
                   'anchors:precisa-do-usuario', e o claim não entrega o card
                   enquanto a decisão não sair.

  (padrão)         não impacta a direção. Vira card comum: nasce em 'to-do',
                   entra na fila, um agente pega. Não para ninguém.

Não use para o que é trivial E está no arquivo que você já está editando: aí
corrija e registre a revisão ('{CODIGO}-R0001: o que mudou e por quê'). Abrir
card para trocar uma palavra é burocracia.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
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
			if cfg.ModoGitHub() && len(cfg.Workflow.Labels) == 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("`workflow.labels` está vazio: a issue nasceria sem a " +
					"label que o pipeline de claim usa para achá-la, e ficaria órfã")
			}
			if !cfg.ModoGitHub() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`escalate` existe no modo github (o card é uma issue). " +
					"No modo local, escreva a dúvida no plano e pare o trabalho — não há " +
					"fila compartilhada de onde tirar o card")
			}

			motivo := strings.Join(args, " ")
			corpoTexto := corpoDaEscalada(motivo, sobre, card, paraUsuario)

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
			titulo := "[plano] " + primeiraLinhaDoMotivo(motivo)
			labels := []string{cfg.Workflow.Labels[0], "anchors:to-do"}
			if paraUsuario {
				titulo = "[decisão] " + primeiraLinhaDoMotivo(motivo)
				labels = append(labels, initx.LabelPrecisaDoUsuario)
			}
			// SOB o card de origem, como LABEL — o que permite listar o que pende sob um
			// trabalho (`--label anchors:sob-44`) e entregá-lo no mesmo PR. Uma frase no
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
				if _, err := exec.Command("gh", "issue", "edit", card,
					"--repo", cfg.Workflow.Repo,
					"--add-label", initx.LabelPrecisaDoUsuario,
				).CombinedOutput(); err != nil {
					fmt.Printf("· aviso: não consegui rotular o card #%s — rotule à mão, "+
						"senão outro agente pega o card e refaz o caminho\n", card)
				} else {
					fmt.Printf("· card #%s parado até a decisão\n", card)
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
	cmd.Flags().StringVar(&sobre, "sobre", "", "o plano ou spec onde está a incoerência")
	cmd.Flags().StringVar(&card, "card", "", "número do card onde a necessidade foi descoberta")
	cmd.Flags().BoolVar(&paraUsuario, "para-usuario", false,
		"a mudança impacta a DIREÇÃO do projeto: vira decisão do usuário e para o card")
	return cmd
}

// primeiraLinhaDoMotivo faz o título da issue, que é uma linha.
func primeiraLinhaDoMotivo(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if len(s) > 70 {
		return s[:67] + "..."
	}
	return s
}

// corpoDaEscalada monta o texto da issue.
//
// Separado do comando porque é ELE o que precisa ser confrontado: o valor está em dizer
// por que o trabalho parou e como destravar. Um teste que precisasse do `gh` para ler
// isso não rodaria em máquina nenhuma, e o texto ficaria sem régua.
func corpoDaEscalada(motivo, sobre, card string, paraUsuario bool) string {
	var b strings.Builder
	if paraUsuario {
		b.WriteString("🛑 **Esta decisão não é do agente.**\n\n")
	} else {
		b.WriteString("📋 **O plano precisa mudar.**\n\n")
	}
	b.WriteString(motivo + "\n\n")
	if sobre != "" {
		b.WriteString(fmt.Sprintf("**Onde:** `%s`\n\n", sobre))
	}
	if paraUsuario {
		b.WriteString("**Por que parou aqui:** quem descobriu interpretou que esta mudança " +
			"impacta a DIREÇÃO do projeto, e isso é decisão de quem o planejou. Mudar por " +
			"conta própria faria o projeto caminhar para um destino que ninguém escolheu — " +
			"e o plano alterado ficaria válido, então nenhum gate acusaria.\n\n")
		b.WriteString("**Como destravar:** decida, e registre a decisão onde ela vale — no " +
			"plano ou na spec, como revisão (`{CODIGO}-R0001: o que mudou e por quê`). Se a " +
			"mudança for grande, um plano novo com `revises:`. Depois remova a label `" +
			initx.LabelPrecisaDoUsuario + "`.\n\n")
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
		"concluir que isto MUDA A DIREÇÃO, não siga: `anchors escalate ... --para-usuario`.\n\n")
	if card != "" {
		b.WriteString(fmt.Sprintf("Nasceu SOB o card #%s (label `%s`), que segue normalmente. "+
			"Os dois se entregam no mesmo PR: o achado apareceu fazendo aquele trabalho, e "+
			"separá-los faria um dos dois esperar sem razão.\n", card, initx.LabelSob(card)))
	}
	return b.String()
}

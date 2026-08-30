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

// --- escalar uma decisão que não é do agente ---
//
// Quem implementa é quem descobre o erro do plano. E aí ele avalia o impacto: uma
// redação ambígua se corrige na hora, com a revisão registrada. Mas uma correção que
// muda PARA ONDE o projeto vai não é dele — nem do gate, que só valida que a decisão
// ficou escrita.
//
// O buraco que este comando fecha é de ERGONOMIA, e ergonomia aqui decide o resultado:
// enquanto escalar for mais trabalhoso que corrigir, o agente corrige. Ele já podia abrir
// a issue e pôr a label na mão — em três comandos e sabendo o nome exato da label. O
// caminho barato era corrigir em silêncio, e o caminho barato é o que acontece.
//
// O escalonamento por EXAUSTÃO (dez revisões sem convergir) já existia no pipeline de
// claim. Este é o por JUÍZO: ninguém está travado, alguém percebeu algo. São gatilhos
// diferentes para a mesma saída, e é por isso que usam a mesma label.
func newEscalateCmd() *cobra.Command {
	var root, sobre, card string
	cmd := &cobra.Command{
		Use:   "escalate <motivo>",
		Short: "Abre uma decisão para o usuário e para de trabalhar no card",
		Long: `Abre uma issue rotulada 'anchors:precisa-do-usuario' e devolve a decisão a
quem pode tomá-la.

Use quando a correção que você faria MUDA A DIREÇÃO do projeto — ou quando você
tem dúvida se muda. Enquanto a label estiver no card, o pipeline de claim não o
entrega a ninguém: o trabalho para onde está, em vez de seguir para um destino
que ninguém escolheu.

Para a incoerência INÓCUA (redação, exemplo, typo) não use isto: corrija e
registre a revisão no próprio arquivo ('{CODIGO}-R0001: o que mudou e por quê').`,
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
			if !cfg.ModoGitHub() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`escalate` existe no modo github (o card é uma issue). " +
					"No modo local, escreva a dúvida no plano e pare o trabalho — não há " +
					"fila compartilhada de onde tirar o card")
			}

			motivo := strings.Join(args, " ")
			corpoTexto := corpoDaEscalada(motivo, sobre, card)

			tmp, err := os.CreateTemp("", "anchors-escalate-*.md")
			if err != nil {
				return err
			}
			defer os.Remove(tmp.Name())
			if _, err := tmp.WriteString(corpoTexto); err != nil {
				return err
			}
			tmp.Close()

			titulo := "[decisão] " + primeiraLinhaDoMotivo(motivo)
			out, err := exec.Command("gh", "issue", "create",
				"--repo", cfg.Workflow.Repo,
				"--title", titulo,
				"--body-file", tmp.Name(),
				"--label", initx.LabelPrecisaDoUsuario,
			).CombinedOutput()
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("abrir a issue: %v — %s", err, strings.TrimSpace(string(out)))
			}
			url := strings.TrimSpace(string(out))
			fmt.Printf("decisão aberta: %s\n", url)

			// O CARD onde o trabalho parou também recebe a label. Sem isso o agente
			// seguinte pegaria o card e refaria o mesmo caminho até a mesma dúvida.
			if card != "" {
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
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&sobre, "sobre", "", "o plano ou spec onde está a incoerência")
	cmd.Flags().StringVar(&card, "card", "", "número do card onde o trabalho parou")
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
func corpoDaEscalada(motivo, sobre, card string) string {
	var b strings.Builder
	b.WriteString("🛑 **Esta decisão não é do agente.**\n\n")
	b.WriteString(motivo + "\n\n")
	if sobre != "" {
		b.WriteString(fmt.Sprintf("**Onde:** `%s`\n\n", sobre))
	}
	b.WriteString("**Por que parou aqui:** a correção mudaria para onde o projeto vai, e " +
		"isso é decisão de quem o planejou. Corrigir por conta própria faria o projeto " +
		"caminhar para um destino que ninguém escolheu — e o plano corrigido ficaria " +
		"válido, então nenhum gate acusaria.\n\n")
	b.WriteString("**Como destravar:** decida, e registre a decisão onde ela vale — no " +
		"plano ou na spec, como revisão (`{CODIGO}-R0001: o que mudou e por quê`). Se a " +
		"mudança for grande, um plano novo com `revises:`. Depois remova a label `" +
		initx.LabelPrecisaDoUsuario + "`.\n\n")
	if card != "" {
		b.WriteString(fmt.Sprintf("Trabalho parado no card #%s.\n", card))
	}
	return b.String()
}

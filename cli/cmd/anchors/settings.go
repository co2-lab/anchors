package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/settings"
	"github.com/spf13/cobra"
)

func newSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "A configuração LOCAL deste agente — o que é dele, e não do projeto",
		Long: `O ` + "`anchors.yaml`" + ` é a Estrutura: versionada, revisada, igual para todo mundo.
O que fica em ` + "`.anchors/settings.yaml`" + ` é o oposto — vale para UM agente numa máquina,
e não vai para o git.

Hoje há uma decisão aqui: se este agente atua nos cards ESCALONADOS
(` + "`needs-user`" + `), que esperam decisão de quem conhece o produto.`,
	}
	cmd.AddCommand(newSettingsShowCmd())
	cmd.AddCommand(newSettingsUserIssuesCmd())
	return cmd
}

func newSettingsShowCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Mostra a configuração local deste agente",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			s, err := settings.Load(absRoot)
			if err != nil {
				return err
			}
			fmt.Printf("configuração local: %s\n\n", settings.Path(absRoot))
			fmt.Printf("  cards escalonados: %s\n", s.Describe())
			if !s.Decided() {
				fmt.Println()
				fmt.Println("  Para decidir:")
				fmt.Println("      anchors settings user-issues sim   # eu decido o produto")
				fmt.Println("      anchors settings user-issues nao   # deixo para quem decide")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	return cmd
}

func newSettingsUserIssuesCmd() *cobra.Command {
	var root, date string
	cmd := &cobra.Command{
		Use:   "user-issues [sim|nao]",
		Short: "Declara se este agente atua nos cards escalonados (`needs-user`)",
		Long: `Os cards ESCALONADOS esperam decisão de quem conhece o produto: o agente achou
algo que muda a direção e escalou em vez de decidir sozinho.

Num projeto com vários devs, cada um pode rodar o seu agente — e nem todos podem
decidir pelo produto. Um agente que pega um card desses e pergunta a quem o está
rodando obtém uma resposta, e a resposta pode não ser a do dono do projeto.

O padrão é NÃO atuar. Quem pode decidir declara que pode.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if date == "" {
				return fmt.Errorf("informe --date AAAA-MM-DD (o Anchors não lê o relógio: " +
					"a data é carimbada por quem declara)")
			}
			var decisao *bool
			if len(args) == 1 {
				if decisao = settings.ParseAnswer(args[0]); decisao == nil {
					return fmt.Errorf("não entendi %q — responda `sim` ou `nao`", args[0])
				}
			} else {
				if decisao, err = askUserIssues(); err != nil {
					return err
				}
			}
			s, err := settings.Load(absRoot)
			if err != nil {
				return err
			}
			s.UserIssues = decisao
			s.Agent = agentID()
			s.DecidedAt = date
			if err := settings.Save(absRoot, s); err != nil {
				return err
			}
			fmt.Printf("✓ %s\n", s.Describe())
			fmt.Printf("  registrado em %s (local, fora do git)\n", settings.Path(absRoot))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&date, "date", "", "OBRIGATÓRIO — AAAA-MM-DD, a data da decisão")
	return cmd
}

// askUserIssues faz a pergunta no terminal.
//
// Repete até entender, e não assume: assumir "não" no que não entendeu seria conveniente e
// errado — quem digitou algo estava respondendo, e descartar a resposta em silêncio faz o
// agente decidir por conta.
func askUserIssues() (*bool, error) {
	in := bufio.NewReader(os.Stdin)
	for tentativa := 0; tentativa < 3; tentativa++ {
		fmt.Println("Você decide o rumo deste produto?")
		fmt.Println()
		fmt.Println("  Os cards ESCALONADOS (`needs-user`) esperam decisão de quem conhece o")
		fmt.Println("  produto — o agente achou algo que muda a direção e escalou.")
		fmt.Println()
		fmt.Println("  Responda `sim` se você pode tomar essas decisões; `nao` se elas são de")
		fmt.Println("  outra pessoa. No `nao`, este agente segue nos cards comuns e os")
		fmt.Println("  escalonados esperam quem os resolve.")
		fmt.Print("\n  [sim/nao] ")

		linha, err := in.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("ler a resposta: %w", err)
		}
		if d := settings.ParseAnswer(linha); d != nil {
			return d, nil
		}
		fmt.Printf("\n  não entendi %q — responda `sim` ou `nao`.\n\n", strings.TrimSpace(linha))
	}
	return nil, fmt.Errorf("sem resposta reconhecida — declare com " +
		"`anchors settings user-issues sim|nao --date AAAA-MM-DD`")
}

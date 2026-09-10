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
	cmd.AddCommand(newSettingsRoleCmd())
	cmd.AddCommand(newSettingsUserIssuesCmd())
	return cmd
}

// newSettingsRoleCmd declara o PERFIL de quem opera este agente.
func newSettingsRoleCmd() *cobra.Command {
	var root, date string
	cmd := &cobra.Command{
		Use:   "role [perfil]",
		Short: "Declara o PERFIL de quem opera este agente",
		Long: `O perfil diz que trabalho é seu, e as capacidades derivam dele.

Não é hierarquia: o arquiteto não manda no dev, ele responde outra pergunta. E não
é permissão de repositório — o git cuida disso. É sobre QUE TRABALHO o agente puxa
do board e COMO ele se comporta diante do que não sabe.

A diferença que mais aparece: só ` + "`product-owner`" + ` e ` + "`architect`" + ` atuam nos cards
ESCALONADOS (` + "`needs-user`" + `), que esperam decisão de quem conhece o produto. Os
outros perfis escalam e seguem para o próximo card.`,
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
			var perfil settings.Role
			if len(args) == 1 {
				if perfil = settings.ParseRole(args[0]); perfil == "" {
					return fmt.Errorf("perfil %q não reconhecido.%s", args[0], roleList())
				}
			} else {
				if perfil, err = askRole(); err != nil {
					return err
				}
			}
			s, err := settings.Load(absRoot)
			if err != nil {
				return err
			}
			s.Role = perfil
			s.Agent = agentID()
			s.DecidedAt = date
			// O campo antigo sai: manter os dois faria a próxima leitura ter duas fontes
			// para a mesma pergunta, e a resposta viria da que alguém esquecesse de mudar.
			s.UserIssues = nil
			if err := settings.Save(absRoot, s); err != nil {
				return err
			}
			printRole(perfil)
			fmt.Printf("\n  registrado em %s (local, fora do git)\n", settings.Path(absRoot))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&date, "date", "", "OBRIGATÓRIO — AAAA-MM-DD, a data da declaração")
	return cmd
}

// roleList monta a lista para a mensagem de erro e para a pergunta.
func roleList() string {
	var b strings.Builder
	b.WriteString("\n\nOs perfis:\n\n")
	for _, r := range settings.KnownRoles() {
		fmt.Fprintf(&b, "  %-22s %s\n", string(r), r.Does())
	}
	return b.String()
}

// printRole mostra o que o perfil declarado significa.
//
// Mostra as CAPACIDADES e a LENTE, e não só o nome: quem acabou de declarar precisa saber o
// que mudou — e o perfil que revisa precisa saber com que pergunta ler o código, senão faz
// a revisão que sabe fazer em vez da que falta.
func printRole(r settings.Role) {
	fmt.Printf("✓ perfil: %s\n", r.Title())
	fmt.Printf("  %s\n", r.Does())

	if r.Can(settings.CapDecideProduct) {
		fmt.Println("\n  Você atua nos cards ESCALONADOS (`needs-user`) — os que esperam")
		fmt.Println("  decisão de quem conhece o produto. Resolver um deles é registrar a")
		fmt.Println("  decisão no card, não respondê-la numa conversa.")
	} else {
		fmt.Println("\n  Você NÃO atua nos cards escalonados. Diante de ambiguidade:")
		fmt.Println("      anchors escalate \"<o que precisa ser decidido>\" --about <arquivo> --for-user")
		fmt.Println("  E siga para o próximo card — não pergunte a quem está rodando você.")
	}
	if lente := r.Lens(); lente != "" {
		fmt.Printf("\n  A LENTE deste perfil, ao revisar:\n  %s.\n", lente)
	}
}

// askRole pergunta o perfil no terminal.
func askRole() (settings.Role, error) {
	in := bufio.NewReader(os.Stdin)
	for tentativa := 0; tentativa < 3; tentativa++ {
		fmt.Println("Qual é o seu perfil neste projeto?")
		fmt.Print(roleList())
		fmt.Print("\n  perfil: ")

		linha, err := in.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("ler a resposta: %w", err)
		}
		if r := settings.ParseRole(linha); r != "" {
			return r, nil
		}
		fmt.Printf("\n  não reconheci %q.\n\n", strings.TrimSpace(linha))
	}
	return "", fmt.Errorf("sem perfil reconhecido — declare com " +
		"`anchors settings role <perfil> --date AAAA-MM-DD`")
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
			fmt.Printf("  perfil: %s\n", s.Describe())
			if s.Role != "" {
				fmt.Println()
				printRole(s.Role)
				fmt.Println("\n  capacidades:")
				for _, c := range s.Role.Caps() {
					fmt.Printf("    %s\n", c)
				}
			} else {
				fmt.Println()
				fmt.Println("  Para declarar:")
				fmt.Println("      anchors settings role --date AAAA-MM-DD")
				fmt.Print(roleList())
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

package ops

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/settings"
	"github.com/spf13/cobra"
)

func newSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "This agent's LOCAL configuration — what is its own, not the project's",
		Long: `The ` + "`anchors.yaml`" + ` is the Structure: versioned, reviewed, the same for everyone.
What lives in ` + "`.anchors/settings.yaml`" + ` is the opposite — it holds for ONE agent on one machine,
and it does not go to git.

Today there is one decision here: whether this agent acts on ESCALATED cards
(` + "`needs-user`" + `), which await a decision from whoever knows the product.`,
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
		Short: "Declare the PROFILE of whoever operates this agent",
		Long: `The profile says which work is yours, and the capabilities derive from it.

It is not hierarchy: the architect does not command the dev, he answers another question. And it
is not repository permission — git takes care of that. It is about WHAT WORK the agent pulls
from the board and HOW it behaves in the face of what it does not know.

The difference that shows up most: only ` + "`product-owner`" + ` and ` + "`architect`" + ` act on
ESCALATED cards (` + "`needs-user`" + `), which await a decision from whoever knows the product. The
other profiles escalate and move on to the next card.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if date == "" {
				return fmt.Errorf("%s", i18n.T("settings.date_required"))
			}
			var perfil settings.Role
			if len(args) == 1 {
				if perfil = settings.ParseRole(args[0]); perfil == "" {
					return fmt.Errorf("%s", i18n.T("settings.role_unrecognized", args[0], common.RoleList()))
				}
			} else {
				if perfil, err = common.AskRole(); err != nil {
					return err
				}
			}
			s, err := settings.Load(absRoot)
			if err != nil {
				return err
			}
			s.Role = perfil
			s.Agent = common.AgentID()
			s.DecidedAt = date
			// O campo antigo sai: manter os dois faria a próxima leitura ter duas fontes
			// para a mesma pergunta, e a resposta viria da que alguém esquecesse de mudar.
			s.UserIssues = nil
			if err := settings.Save(absRoot, s); err != nil {
				return err
			}
			common.PrintRole(perfil)
			fmt.Print(i18n.T("settings.recorded_at", settings.Path(absRoot)))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&date, "date", "", "REQUIRED — YYYY-MM-DD, the date of the declaration")
	return cmd
}

func newSettingsShowCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show this agent's local configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			s, err := settings.Load(absRoot)
			if err != nil {
				return err
			}
			fmt.Print(i18n.T("settings.show.header", settings.Path(absRoot)))
			fmt.Print(i18n.T("settings.show.role", s.Describe()))
			if s.Role != "" {
				fmt.Println()
				common.PrintRole(s.Role)
				fmt.Println(i18n.T("settings.show.capabilities"))
				for _, c := range s.Role.Caps() {
					fmt.Printf("    %s\n", c)
				}
			} else {
				fmt.Println()
				fmt.Println(i18n.T("settings.show.how_to_declare"))
				fmt.Print(common.RoleList())
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

func newSettingsUserIssuesCmd() *cobra.Command {
	var root, date string
	cmd := &cobra.Command{
		Use:   "user-issues [sim|nao]",
		Short: "Declare whether this agent acts on escalated cards (`needs-user`)",
		Long: `ESCALATED cards await a decision from whoever knows the product: the agent found
something that changes the direction and escalated instead of deciding alone.

In a project with several devs, each one can run their own agent — and not all of them can
decide for the product. An agent that picks up one of these cards and asks whoever is
running it gets an answer, and that answer may not be the project owner's.

The default is NOT to act. Whoever can decide declares that they can.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if date == "" {
				return fmt.Errorf("%s", i18n.T("settings.date_required"))
			}
			var decisao *bool
			if len(args) == 1 {
				if decisao = settings.ParseAnswer(args[0]); decisao == nil {
					return fmt.Errorf("%s", i18n.T("settings.user_issues.invalid_answer", args[0]))
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
			s.Agent = common.AgentID()
			s.DecidedAt = date
			if err := settings.Save(absRoot, s); err != nil {
				return err
			}
			fmt.Printf("✓ %s\n", s.Describe())
			fmt.Print(i18n.T("settings.recorded_at", settings.Path(absRoot)))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&date, "date", "", "REQUIRED — YYYY-MM-DD, the date of the decision")
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
		fmt.Print(i18n.T("settings.user_issues.prompt"))

		linha, err := in.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("role.read_response_error"), err)
		}
		if d := settings.ParseAnswer(linha); d != nil {
			return d, nil
		}
		fmt.Printf("\n  %s\n\n", i18n.T("settings.user_issues.invalid_answer", strings.TrimSpace(linha)))
	}
	return nil, fmt.Errorf("%s", i18n.T("settings.user_issues.unrecognized_error"))
}

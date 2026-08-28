package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/suggestion"
	"github.com/spf13/cobra"
)

// `anchors suggest` é o verbo das CORREÇÕES PROPOSTAS: listar o que está aguardando
// decisão, ver o diff, aplicar, aprovar ou recusar.
//
// O comando não PROPÕE nada — quem propõe são os gates e o julgamento por IA. Aqui
// mora só a decisão, e é de propósito: a proposta é automática, a aceitação não.
func newSuggestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest",
		Short: "Lista, aplica e decide as correções propostas",
		Long: `As correções que os gates e o julgamento por IA propõem, aguardando decisão.

Uma sugestão é um patch do git mais o porquê. O ESTADO é a pasta em que ela vive
(pending/approved/rejected): mover é decidir.

  anchors suggest list                          o que aguarda decisão
  anchors suggest show <id>                     o porquê e o diff
  anchors suggest apply <id> --reason "..."     aplica o patch e aprova
  anchors suggest reject <id> --reason "..."    recusa, com o motivo

A razão é OBRIGATÓRIA ao decidir: sem ela, escolha e descuido ficam
indistinguíveis — a mesma regra do '@no-test' nu.`,
	}
	cmd.AddCommand(newSuggestListCmd(), newSuggestShowCmd(),
		newSuggestApplyCmd(), newSuggestRejectCmd())
	return cmd
}

func newSuggestListCmd() *cobra.Command {
	var root, estado string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lista as sugestões de um estado (default: pending)",
		RunE: func(_ *cobra.Command, _ []string) error {
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}
			st := suggestion.State(estado)
			ids, err := suggestion.List(abs, st)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				fmt.Printf("nenhuma sugestão em %s\n", st)
				return nil
			}
			fmt.Printf("%d sugestão(ões) em %s:\n\n", len(ids), st)
			for _, id := range ids {
				fmt.Printf("  %s\n", id)
			}
			if st == suggestion.Pending {
				fmt.Printf("\n`anchors suggest show <id>` para ver o diff.\n")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&estado, "state", string(suggestion.Pending), "pending|approved|rejected")
	return cmd
}

func newSuggestShowCmd() *cobra.Command {
	var root, estado string
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Mostra o porquê e o diff de uma sugestão",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}
			p := filepath.Join(abs, suggestion.Dir, estado, args[0]+".md")
			b, err := os.ReadFile(p)
			if err != nil {
				return fmt.Errorf("sugestão %q não encontrada em %s", args[0], estado)
			}
			fmt.Println(string(b))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&estado, "state", string(suggestion.Pending), "pending|approved|rejected")
	return cmd
}

func newSuggestApplyCmd() *cobra.Command {
	var root, reason string
	var auto, dryRun bool
	cmd := &cobra.Command{
		Use:   "apply <id>",
		Short: "Aplica o patch e move a sugestão para approved",
		Long: `Aplica o patch com 'git apply' e, se der certo, aprova a sugestão.

A aplicação vem ANTES da aprovação de propósito: aprovar algo que não aplica
deixaria o registro dizendo que a correção entrou quando ela não entrou.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}
			id := args[0]
			patch, err := suggestion.PatchOf(abs, id, suggestion.Pending)
			if err != nil {
				return err
			}
			// `--check` primeiro: um patch que não casa mais (o arquivo mudou desde a
			// proposta) precisa falhar ANTES de mexer no working tree, e com uma
			// mensagem que diga o que houve.
			if err := gitApply(abs, patch, true); err != nil {
				return fmt.Errorf("o patch não casa mais com o arquivo — ele mudou desde "+
					"a proposta. Rode o check de novo para gerar uma sugestão atual: %w", err)
			}
			if dryRun {
				fmt.Printf("✓ o patch de %s aplica limpo (nada foi alterado — dry-run)\n", id)
				return nil
			}
			if err := gitApply(abs, patch, false); err != nil {
				return err
			}
			if strings.TrimSpace(reason) == "" {
				reason = "patch aplicado sem ressalvas"
			}
			if err := suggestion.Decide(abs, id, suggestion.Approved, reason, auto); err != nil {
				return err
			}
			fmt.Printf("✓ %s aplicado e aprovado\n", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&reason, "reason", "", "por que aceitar (entra no registro)")
	cmd.Flags().BoolVar(&auto, "auto", false, "marca a decisão como automática (auto_judgment)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "só confere se o patch ainda casa")
	return cmd
}

func newSuggestRejectCmd() *cobra.Command {
	var root, reason string
	var auto bool
	cmd := &cobra.Command{
		Use:   "reject <id>",
		Short: "Recusa uma sugestão, com o motivo",
		Long: `Recusa a sugestão. A razão é OBRIGATÓRIA — e a recusada NÃO é apagada:
é o registro que impede a mesma proposta de voltar na varredura seguinte como
novidade, e que preserva o motivo para quem vier depois.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}
			if err := suggestion.Decide(abs, args[0], suggestion.Rejected, reason, auto); err != nil {
				return err
			}
			fmt.Printf("✓ %s recusada\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&reason, "reason", "", "por que recusar (OBRIGATÓRIO)")
	cmd.Flags().BoolVar(&auto, "auto", false, "marca a decisão como automática (auto_judgment)")
	return cmd
}

// gitApply roda `git apply` lendo o patch da entrada padrão.
func gitApply(root, patch string, check bool) error {
	args := []string{"apply"}
	if check {
		args = append(args, "--check")
	}
	c := exec.Command("git", args...)
	c.Dir = root
	c.Stdin = strings.NewReader(patch + "\n")
	if out, err := c.CombinedOutput(); err != nil {
		// Uma sugestão É um patch: sem git não há como aplicá-la, e o erro do git sobre
		// um repo inexistente não menciona nem o patch nem a sugestão.
		if msg := gitmeta.Explica(gitmeta.Verifica(root), "aplicar o patch da sugestão"); msg != "" {
			return errors.New(msg)
		}
		return fmt.Errorf("git apply: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

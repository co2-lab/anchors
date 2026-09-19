package ops

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
		Short: "List, apply and decide the proposed fixes",
		Long: `The fixes that the gates and the AI judgment propose, awaiting a decision.

A suggestion is a git patch plus the why. The STATE is the folder it lives in
(pending/approved/rejected): moving it is deciding.

  anchors suggest list                          what awaits a decision
  anchors suggest show <id>                     the why and the diff
  anchors suggest apply <id> --reason "..."     applies the patch and approves
  anchors suggest reject <id> --reason "..."    rejects, with the reason

The reason is MANDATORY when deciding: without it, choice and carelessness become
indistinguishable — the same rule as the bare '@no-test'.`,
	}
	cmd.AddCommand(newSuggestListCmd(), newSuggestShowCmd(),
		newSuggestApplyCmd(), newSuggestRejectCmd())
	return cmd
}

func newSuggestListCmd() *cobra.Command {
	var root, estado string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the suggestions in a state (default: pending)",
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
				fmt.Printf("no suggestion in %s\n", st)
				return nil
			}
			fmt.Printf("%d suggestion(s) in %s:\n\n", len(ids), st)
			for _, id := range ids {
				fmt.Printf("  %s\n", id)
			}
			if st == suggestion.Pending {
				fmt.Printf("\n`anchors suggest show <id>` to see the diff.\n")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&estado, "state", string(suggestion.Pending), "pending|approved|rejected")
	return cmd
}

func newSuggestShowCmd() *cobra.Command {
	var root, estado string
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show the why and the diff of a suggestion",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}
			p := filepath.Join(abs, suggestion.Dir, estado, args[0]+".md")
			b, err := os.ReadFile(p)
			if err != nil {
				return fmt.Errorf("suggestion %q not found in %s", args[0], estado)
			}
			fmt.Println(string(b))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&estado, "state", string(suggestion.Pending), "pending|approved|rejected")
	return cmd
}

func newSuggestApplyCmd() *cobra.Command {
	var root, reason string
	var auto, dryRun bool
	cmd := &cobra.Command{
		Use:   "apply <id>",
		Short: "Apply the patch and move the suggestion to approved",
		Long: `Applies the patch with 'git apply' and, if it works, approves the suggestion.

The application comes BEFORE the approval on purpose: approving something that does not apply
would leave the record saying the fix went in when it did not.`,
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
				return fmt.Errorf("the patch no longer matches the file — it changed since "+
					"the proposal. Run the check again to generate a current suggestion: %w", err)
			}
			if dryRun {
				fmt.Printf("✓ the patch of %s applies cleanly (nothing was changed — dry-run)\n", id)
				return nil
			}
			if err := gitApply(abs, patch, false); err != nil {
				return err
			}
			if strings.TrimSpace(reason) == "" {
				reason = "patch applied without reservations"
			}
			if err := suggestion.Decide(abs, id, suggestion.Approved, reason, auto); err != nil {
				return err
			}
			fmt.Printf("✓ %s applied and approved\n", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&reason, "reason", "", "why accept it (goes into the record)")
	cmd.Flags().BoolVar(&auto, "auto", false, "mark the decision as automatic (enable_auto_judgment)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "only check whether the patch still applies")
	return cmd
}

func newSuggestRejectCmd() *cobra.Command {
	var root, reason string
	var auto bool
	cmd := &cobra.Command{
		Use:   "reject <id>",
		Short: "Reject a suggestion, with the reason",
		Long: `Rejects the suggestion. The reason is MANDATORY — and the rejected one is NOT erased:
it is the record that keeps the same proposal from coming back in the next sweep as
news, and that preserves the reason for whoever comes later.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}
			if err := suggestion.Decide(abs, args[0], suggestion.Rejected, reason, auto); err != nil {
				return err
			}
			fmt.Printf("✓ %s rejected\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&reason, "reason", "", "why refuse it (REQUIRED)")
	cmd.Flags().BoolVar(&auto, "auto", false, "mark the decision as automatic (enable_auto_judgment)")
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
		if msg := gitmeta.Explain(gitmeta.Check(root), "apply the suggestion patch"); msg != "" {
			return errors.New(msg)
		}
		return fmt.Errorf("git apply: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

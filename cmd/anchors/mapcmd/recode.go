package mapcmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/recode"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// `anchors recode <antigo> <novo>` renomeia um código de identidade e propaga a troca
// por todas as superfícies textuais (header code:/ref:, scenario-codes, refs cruzadas)
// em TODO o projeto. Dry-run por padrão; --apply grava e reconstrói o mapa.
//
// É a base da estabilidade REVERSÍVEL do código: antes, renomear à mão espalhava
// resíduo (o app de referência tem um TXDT→TCDT incompleto que prova isso). Este comando fecha essa
// lacuna. (Fase 1: o motor genérico — testIDs e rename de ARQUIVOS vêm depois, guiados
// por convenção declarada no anchors.yaml.)
func newRecodeCmd() *cobra.Command {
	var root string
	var apply bool
	cmd := &cobra.Command{
		Use:   "recode <old> <new>",
		Short: "Rename an identity code and propagate it across the project",
		Long: `Swaps an identity code (e.g. TCDT → TCTX) across every surface:
  • the @anchors header (code:/ref:, comma-separated lists included);
  • derived scenario-codes (CODEX-B01, CODEX-S02, CODE-DS-*, CODEX-VR, CODE-FP01b…);
  • bare mentions of the code in cross-references from other units.

DRY-RUN by default: shows each occurrence classified, writes nothing. --apply writes and
rebuilds the map (the graph is reindexed from the headers, the source of truth).

Validates: NEW must not already belong to another unit (collision); OLD must exist.

  anchors recode TCDT TCTX            # dry-run: what would change
  anchors recode TCDT TCTX --apply    # applies + rebuilds the map

Out of scope in this version (phase 2): swapping the testID prefix and renaming FILES
whose name contains the code — both depend on a convention declared in anchors.yaml.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load %s: %w", config.DefaultFile, err)
			}
			old, new := strings.ToUpper(args[0]), strings.ToUpper(args[1])

			plan, err := recode.BuildPlan(absRoot, cfg, old, new)
			if err != nil {
				return err
			}

			// Relatório (dry-run e apply mostram o plano).
			fmt.Printf("recode %s → %s: %d file(s), %d content substitution(s)", plan.Old, plan.New, len(plan.Files), plan.Total)
			if plan.TestIDs > 0 {
				fmt.Printf(", %d testID(s)", plan.TestIDs)
			}
			if len(plan.Renames) > 0 {
				fmt.Printf(", %d file(s) to rename", len(plan.Renames))
			}
			fmt.Print("\n\n")
			for _, fc := range plan.Files {
				byKind := map[string]int{}
				for _, o := range fc.Occurrences {
					byKind[o.Kind]++
				}
				fmt.Printf("  %s\n", fc.Path)
				for _, k := range []string{"header", "scenario-code", "bare-ref"} {
					if byKind[k] > 0 {
						fmt.Printf("      %-14s %d\n", k, byKind[k])
					}
				}
			}
			if len(plan.Renames) > 0 {
				fmt.Println("\n  rename (git mv):")
				for _, r := range plan.Renames {
					fmt.Printf("      %s → %s\n", r.From, r.To)
				}
			}
			if plan.TestIDLegacy != "" {
				fmt.Printf("\n  ⚠ %s\n", plan.TestIDLegacy)
			}
			fmt.Println()

			if !apply {
				fmt.Println("(dry-run — nothing was written. Review the above and run with --apply to write.)")
				return nil
			}

			n, err := plan.Apply(absRoot)
			if err != nil {
				// O recode aplica em MASSA, e o Apply para no primeiro problema. Sair só
				// com o erro esconderia o que já aconteceu: o disco ficou num estado
				// intermediário — parte aplicada, parte não, e o mapa apontando para a
				// forma antiga. Quem não souber disso vai reexecutar sobre um projeto
				// meio convertido.
				if n > 0 {
					fmt.Printf("⚠ %d file(s) had ALREADY been changed before the failure.\n", n)
					fmt.Println("  the project is half converted and the map still describes the old form.")
					fmt.Println("  review with `git status` before running again.")
				}
				return err
			}
			fmt.Printf("✓ %d file(s) rewritten.\n", n)

			// Reconstrói o mapa a partir dos headers (fonte da verdade) — não edita o
			// graph por regex, para não perpetuar divergências header↔graph.
			files, serr := scan.Walk(absRoot, cfg)
			if serr != nil {
				return fmt.Errorf("map rebuild (scan): %w", serr)
			}
			g := mapx.Build(files, cfg, gitmeta.AllCommitDates(absRoot))
			if serr := mapx.Save(g, filepath.Join(absRoot, mapx.DefaultPath)); serr != nil {
				return fmt.Errorf("map rebuild (save): %w", serr)
			}
			fmt.Printf("✓ map rebuilt (%d nodes).\n", len(g.Nodes))
			fmt.Println("  check it with `anchors check --all`.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&apply, "apply", false, "writes the changes (without this, it is a dry-run)")
	return cmd
}

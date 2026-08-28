package main

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
		Use:   "recode <antigo> <novo>",
		Short: "Renomeia um código de identidade e propaga por todo o projeto",
		Long: `Troca um código de identidade (ex.: TCDT → TCTX) em todas as superfícies:
  • header @anchors (code:/ref:, inclusive listas com vírgula);
  • scenario-codes derivados (CODEX-B01, CODEX-S02, CODE-DS-*, CODEX-VR, CODE-FP01b…);
  • menções nuas do código em referências cruzadas de outras unidades.

Por padrão é DRY-RUN: mostra cada ocorrência classificada, não grava. --apply grava e
reconstrói o mapa (o graph é reindexado a partir dos headers, a fonte da verdade).

Valida: NEW não pode já ser de outra unidade (colisão); OLD tem de existir.

  anchors recode TCDT TCTX            # dry-run: o que mudaria
  anchors recode TCDT TCTX --apply    # aplica + rebuild do mapa

Fora do escopo desta versão (fase 2): trocar prefixo de testID e renomear ARQUIVOS
cujo nome contém o código — dependem de convenção declarada no anchors.yaml.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar %s: %w", config.DefaultFile, err)
			}
			old, new := strings.ToUpper(args[0]), strings.ToUpper(args[1])

			plan, err := recode.BuildPlan(absRoot, cfg, old, new)
			if err != nil {
				return err
			}

			// Relatório (dry-run e apply mostram o plano).
			fmt.Printf("recode %s → %s: %d arquivo(s), %d subst. de conteúdo", plan.Old, plan.New, len(plan.Files), plan.Total)
			if plan.TestIDs > 0 {
				fmt.Printf(", %d testID(s)", plan.TestIDs)
			}
			if len(plan.Renames) > 0 {
				fmt.Printf(", %d arquivo(s) a renomear", len(plan.Renames))
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
				fmt.Println("\n  renomear (git mv):")
				for _, r := range plan.Renames {
					fmt.Printf("      %s → %s\n", r.From, r.To)
				}
			}
			if plan.TestIDLegacy != "" {
				fmt.Printf("\n  ⚠ %s\n", plan.TestIDLegacy)
			}
			fmt.Println()

			if !apply {
				fmt.Println("(dry-run — nada foi escrito. Reveja acima e rode com --apply para gravar.)")
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
					fmt.Printf("⚠ %d arquivo(s) JÁ foram alterados antes da falha.\n", n)
					fmt.Println("  o projeto está meio convertido e o mapa ainda descreve a forma antiga.")
					fmt.Println("  revise com `git status` antes de rodar de novo.")
				}
				return err
			}
			fmt.Printf("✓ %d arquivo(s) reescritos.\n", n)

			// Reconstrói o mapa a partir dos headers (fonte da verdade) — não edita o
			// graph por regex, para não perpetuar divergências header↔graph.
			files, serr := scan.Walk(absRoot, cfg)
			if serr != nil {
				return fmt.Errorf("rebuild do mapa (scan): %w", serr)
			}
			g := mapx.Build(files, cfg, gitmeta.AllCommitDates(absRoot))
			if serr := mapx.Save(g, filepath.Join(absRoot, mapx.DefaultPath)); serr != nil {
				return fmt.Errorf("rebuild do mapa (save): %w", serr)
			}
			fmt.Printf("✓ mapa reconstruído (%d nós).\n", len(g.Nodes))
			fmt.Println("  confira com `anchors check --all`.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&apply, "apply", false, "grava as mudanças (sem isto, é dry-run)")
	return cmd
}

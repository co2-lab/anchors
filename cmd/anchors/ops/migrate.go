package ops

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
	"github.com/spf13/cobra"
)

// O comando que a mensagem de erro do formato PROMETE.
//
// Um erro que diz "rode `anchors migrate`" e não tem o comando é pior que nenhum erro:
// quem lê tenta, falha, e passa a desconfiar da próxima mensagem.
func newMigrateCmd() *cobra.Command {
	var root string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Bring the Anchors files up to the current format",
		Long: `Rewrites the ` + "`anchors.graph.yaml`" + ` and the ` + "`anchors.yaml`" + ` into the format
that this binary reads.

The format goes up when a change makes the file UNREADABLE for the previous
version — a renamed key, a value that changed shape. It does not go up for a new
optional field: an old binary simply ignores it.

The migration runs ONCE and leaves the project ready to commit. It is idempotent:
running it again changes nothing.

    anchors migrate              # migrates
    anchors migrate --dry-run    # says what it would do, without writing`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true

			alvos := []string{
				filepath.Join(absRoot, mapx.DefaultPath),
				filepath.Join(absRoot, config.DefaultFile),
			}

			var mudou bool
			for _, alvo := range alvos {
				r, err := migra.MigrateFile(alvo, mapx.FormatoAtual, dryRun)
				if err != nil {
					// Arquivo ausente não é erro: um projeto pode não ter mapa ainda, e
					// recusar a migração inteira por causa disso deixaria o outro arquivo
					// no formato velho.
					fmt.Printf("· %s: %v\n", filepath.Base(alvo), err)
					continue
				}
				if !r.Changed {
					fmt.Printf("· %s is already in format %d\n", filepath.Base(alvo), r.To)
					continue
				}
				mudou = true
				verbo := "migrated"
				if dryRun {
					verbo = "would be migrated"
				}
				fmt.Printf("✓ %s %s: format %d → %d\n", filepath.Base(alvo), verbo, r.De, r.To)

				// AS CHAVES, uma a uma e com a contagem. Um "migrado" seco não diz o que
				// mudou, e quem revisa o diff precisa saber o que esperar antes de abri-lo.
				chaves := make([]string, 0, len(r.Replaced))
				for k := range r.Replaced {
					chaves = append(chaves, k)
				}
				sort.Strings(chaves)
				for _, k := range chaves {
					fmt.Printf("    %s  (%d occurrence(s))\n", k, r.Replaced[k])
				}
			}

			if !mudou {
				return nil
			}
			fmt.Println()
			if dryRun {
				fmt.Println("  nothing was written — run without `--dry-run` to apply.")
				return nil
			}
			// O COMMIT é de quem rodou, e dizer isso importa: o mapa é versionado, e uma
			// migração que fica só na máquina faz o próximo agente reencontrar o formato
			// velho — e migrar de novo, gerando o mesmo diff outra vez.
			fmt.Println("  the files are VERSIONED: commit the migration so that the")
			fmt.Println("  other agents do not redo it.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "say what it would do, without writing")
	return cmd
}

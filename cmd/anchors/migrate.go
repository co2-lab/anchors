package main

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
		Short: "Leva os arquivos do Anchors ao formato atual",
		Long: `Reescreve o ` + "`anchors.graph.yaml`" + ` e o ` + "`anchors.yaml`" + ` no formato
que este binário lê.

O formato sobe quando uma mudança torna o arquivo ILEGÍVEL para a versão
anterior — uma chave renomeada, um valor que mudou de forma. Não sobe por campo
novo opcional: um binário velho simplesmente o ignora.

A migração roda UMA VEZ e deixa o projeto pronto para commit. Ela é idempotente:
rodar de novo não muda nada.

    anchors migrate              # migra
    anchors migrate --dry-run    # diz o que faria, sem escrever`,
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
					fmt.Printf("· %s já está no formato %d\n", filepath.Base(alvo), r.To)
					continue
				}
				mudou = true
				verbo := "migrado"
				if dryRun {
					verbo = "seria migrado"
				}
				fmt.Printf("✓ %s %s: formato %d → %d\n", filepath.Base(alvo), verbo, r.De, r.To)

				// AS CHAVES, uma a uma e com a contagem. Um "migrado" seco não diz o que
				// mudou, e quem revisa o diff precisa saber o que esperar antes de abri-lo.
				chaves := make([]string, 0, len(r.Replaced))
				for k := range r.Replaced {
					chaves = append(chaves, k)
				}
				sort.Strings(chaves)
				for _, k := range chaves {
					fmt.Printf("    %s  (%d ocorrência(s))\n", k, r.Replaced[k])
				}
			}

			if !mudou {
				return nil
			}
			fmt.Println()
			if dryRun {
				fmt.Println("  nada foi escrito — rode sem `--dry-run` para aplicar.")
				return nil
			}
			// O COMMIT é de quem rodou, e dizer isso importa: o mapa é versionado, e uma
			// migração que fica só na máquina faz o próximo agente reencontrar o formato
			// velho — e migrar de novo, gerando o mesmo diff outra vez.
			fmt.Println("  os arquivos são VERSIONADOS: commite a migração para que os")
			fmt.Println("  outros agentes não a refaçam.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "diz o que faria, sem escrever")
	return cmd
}

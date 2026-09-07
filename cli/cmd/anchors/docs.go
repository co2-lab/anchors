package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Compila a documentação a partir dos templates de `doct/`",
		Long: `A documentação tem três camadas, e o conteúdo mora em UMA:

  *.spec.md        a fonte — o conteúdo mora aqui, e só aqui
  doct/*.md.tmpl   o template — a moldura, que REFERENCIA trechos das specs
  docs/*.md        o compilado — conteúdo real, gerado, ninguém edita

Markdown não tem inclusão nativa: sem build, a única saída seria duplicar à mão, e
duplicação escrita à mão desatualiza sem ninguém ver. Aqui a duplicação existe só no
compilado, onde o ` + "`docs-fresh`" + ` acusa quando ela envelhece.`,
	}
	cmd.AddCommand(newDocsBuildCmd())
	cmd.AddCommand(newDocsInitCmd())
	cmd.AddCommand(newDocsDutiesCmd())
	return cmd
}

func newDocsBuildCmd() *cobra.Command {
	var root, mapPath string
	var dryRun bool
	var maxUnits, maxLines int
	padrao := doct.DefaultLayout()
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Compila `doct/*.md.tmpl` em `docs/*.md`",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}
			c, err := doct.New(absRoot, g)
			if err != nil {
				return err
			}
			// O LAYOUT é resolvido AQUI, uma vez, e vale para todas as páginas: a de
			// camada, o índice de regras e o de comportamento perguntam ao mesmo objeto.
			// Deixar cada template decidir foi o que produziu 483 links quebrados — a
			// página resumia por um critério e o link era montado por outro.
			c.Layout = doct.Layout{MaxUnits: maxUnits, MaxLines: maxLines}
			res, err := c.Build(dryRun)
			if err != nil {
				return err
			}
			return printDocsResult(res, dryRun)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa (padrão: <root>/"+mapx.DefaultPath+")")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "compila sem escrever (é o que o gate usa)")
	// O CORTE entre a página completa e a resumida. Vale para o projeto inteiro: a
	// documentação em que uma camada segue uma lógica e a vizinha outra é a que obriga
	// quem lê a descobrir a lógica antes de achar o que procura.
	cmd.Flags().IntVar(&maxUnits, "max-units", padrao.MaxUnits,
		"acima de N unidades, a página da camada traz o RESUMO de cada uma")
	cmd.Flags().IntVar(&maxLines, "max-lines", padrao.MaxLines,
		"acima de N linhas de spec somadas, idem — cinco specs longas pesam mais que vinte curtas")
	return cmd
}

func printDocsResult(res doct.Result, dryRun bool) error {
	verbo := "escrito"
	if dryRun {
		verbo = "compilado (não escrito)"
	}
	for _, f := range res.Written {
		fmt.Printf("  ✓ %s/%s %s\n", doct.OutDir, f, verbo)
	}
	// Os ignorados NÃO são um detalhe: o template existe e o documento não foi gerado.
	// Silenciar isso faria o autor acreditar que a doc reflete as specs quando ela é outra
	// coisa — o mesmo defeito que o marcador existe para evitar, pela porta de trás.
	if len(res.Skipped) > 0 {
		fmt.Printf("\n  %d arquivo(s) NÃO gerado(s) — existem em `%s/` sem o marcador:\n",
			len(res.Skipped), doct.OutDir)
		for _, f := range res.Skipped {
			fmt.Printf("    %s/%s\n", doct.OutDir, f)
		}
		fmt.Println("  Foram escritos à mão, e o compilador não apaga trabalho de ninguém.")
		fmt.Println("  Ou apague o `.tmpl` (o documento é manual), ou apague o `.md`")
		fmt.Println("  (o conteúdo já está nas specs e o template o reconstrói).")
	}
	if len(res.Written) == 0 && len(res.Skipped) == 0 {
		fmt.Printf("  nenhum template em `%s/` — nada a compilar\n", doct.Dir)
	}
	return nil
}

// docsFreshHint é o que o gate diz quando o compilado está defasado.
func docsFreshHint(arquivos []string) string {
	return fmt.Sprintf("%d documento(s) fora de data: %s.\n"+
		"  Rode `anchors docs build`.",
		len(arquivos), strings.Join(arquivos, ", "))
}

// newDocsInitCmd escreve o ESQUELETO da documentação.
//
// Existe para que ninguém precise inventar a organização do zero. Inventar a cada projeto
// é como se chega a uma documentação onde cada página segue uma lógica diferente — e
// quem lê tem de descobrir a lógica antes de achar o que procura.
func newDocsInitCmd() *cobra.Command {
	var root string
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Escreve o esqueleto de `" + doct.Dir + "/` (arquitetura, comportamento, regras, camadas)",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}
			c, err := doct.New(absRoot, g)
			if err != nil {
				return err
			}
			escritos, pulados, err := c.InitScaffolds(force)
			if err != nil {
				return err
			}
			for _, f := range escritos {
				fmt.Printf("  ✓ %s/%s\n", doct.Dir, f)
			}
			// PULADO não é silêncio: um template que já existe e não foi tocado é o caso
			// normal de rodar de novo, e dizer isso evita que alguém procure por que a
			// mudança não apareceu.
			for _, f := range pulados {
				fmt.Printf("  · %s/%s já existe (use --force para sobrescrever)\n", doct.Dir, f)
			}
			fmt.Printf("\nEdite os templates e rode `anchors docs build`.\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&force, "force", false, "sobrescreve templates existentes")
	return cmd
}

// newDocsDutiesCmd lista as documentações que o projeto declarou como obrigatórias, e o
// que cada uma pede.
//
// É o comando que um agente roda antes de mexer numa unidade: "o que eu tenho de
// documentar, e o que essa documentação precisa responder". Sem ele, a instrução seria
// "atualize a doc", que produz o endpoint acrescentado ao OpenAPI sem os erros.
func newDocsDutiesCmd() *cobra.Command {
	var root, layer string
	cmd := &cobra.Command{
		Use:   "duties",
		Short: "Lista as documentações obrigatórias do projeto e o que cada uma exige",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return err
			}
			docs := cfg.AllRequiredDocs()
			if layer != "" {
				docs = cfg.RequiredFor(layer)
			}
			if len(docs) == 0 {
				if layer != "" {
					fmt.Printf("Nenhuma documentação é obrigatória ao alterar `%s`.\n", layer)
					return nil
				}
				fmt.Println("O projeto não declara documentações obrigatórias.")
				fmt.Println("Declare-as em `docs:` no `" + config.DefaultFile + "` — tipos que o")
				fmt.Println("Anchors sabe instruir: " + strings.Join(doct.KnownKinds(), ", ") + ".")
				return nil
			}
			if layer != "" {
				fmt.Printf("Alterar `%s` obriga tocar:\n\n", layer)
			} else {
				fmt.Println("Documentações obrigatórias:")
				fmt.Println()
			}
			for _, d := range docs {
				fmt.Print(doct.Duty(d))
				fmt.Println()
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&layer, "layer", "", "só as obrigatórias ao alterar esta camada")
	return cmd
}

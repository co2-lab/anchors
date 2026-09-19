package ops

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

func newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Compile the documentation from the templates in `doct/`",
		Long: `The documentation has three layers, and the content lives in ONE:

  *.spec.md        the source — the content lives here, and only here
  doct/*.md.tmpl   the template — the frame, which REFERENCES excerpts of the specs
  docs/*.md        the compiled — real content, generated, nobody edits it

Markdown has no native inclusion: without a build, the only way out would be to duplicate by hand, and
duplication written by hand goes stale without anyone seeing. Here the duplication exists only in the
compiled, where the ` + "`docs-fresh`" + ` flags it when it ages.`,
	}
	cmd.AddCommand(newDocsBuildCmd())
	cmd.AddCommand(newDocsInitCmd())
	cmd.AddCommand(newDocsDutiesCmd())
	return cmd
}

func newDocsBuildCmd() *cobra.Command {
	var root, mapPath string
	var dryRun, semRebuild bool
	var maxUnits, maxLines int
	padrao := doct.DefaultLayout()
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Compile `doct/*.md.tmpl` into `docs/*.md`",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}

			// O MAPA É RECONSTRUÍDO AQUI, e não lido do disco.
			//
			// Este comando compila a partir do MAPA, não da árvore. Ele lia o
			// `anchors.graph.yaml` como estava — e um mapa mais velho que as specs não
			// conhece as unidades novas, então o compilado saía SEM os cenários delas.
			// Sem erro: o arquivo encolhia e o commit parecia normal.
			//
			// MEDIDO no projeto de referência (#717, #743): o compilado no `develop` tinha
			// 31 entradas a menos do que as specs produzem — DTSTD 11, HLCHH 4, NTCNN 8,
			// SRMTS 8, quatro unidades inteiras. O commit que as apagou é de um PR que não
			// toca nenhuma delas; quem o escreveu rodou `docs build` antes do `map build`
			// e não tinha como saber.
			//
			// A ordem certa não podia ser responsabilidade de quem chama. Recusar o mapa
			// velho seria melhor que aceitar, mas ainda devolveria o trabalho para quem
			// não causou o problema — e a régua do projeto induzia a ordem errada.
			//
			// Reconstruir é o mesmo que o `map build` faz (`scan.Walk` → `mapx.Build`), e
			// os CARIMBOS do mapa em disco são preservados: eles são estado de trabalho,
			// não derivado, e o laudo de cada um vive no `--reason` do `anchors judge`.
			//
			// Este comando NÃO grava o mapa. Reconstruir para compilar é dele; decidir o
			// que o `anchors.graph.yaml` guarda é do `map build`.
			// `--no-map-rebuild` existe para quem PRECISA compilar contra um mapa
			// específico: conferir o que uma revisão antiga produzia, compilar contra um
			// `--map` de outro lugar, ou rodar onde o scan não vale (uma árvore parcial).
			//
			// É opt-out e não padrão porque o modo perigoso tem de ser o PEDIDO. Quem
			// passa a flag está dizendo "sei que este mapa é a entrada que eu quero"; quem
			// não passa recebe a garantia sem precisar saber que ela existe.
			var g *mapx.Graph
			if semRebuild {
				var err error
				if g, err = mapx.Load(mapPath); err != nil {
					return fmt.Errorf("load map: %w (run `anchors map build`, or "+
						"drop the `--no-map-rebuild` so the compiled output comes from the tree)", err)
				}
				fmt.Fprintln(cmd.ErrOrStderr(),
					"anchors: `--no-map-rebuild` — the compiled output comes from the map on disk, and what "+
						"it does not know does NOT go in")
			} else {
				cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
				if err != nil {
					return fmt.Errorf("load %s: %w (run `anchors init` to create it)",
						config.DefaultFile, err)
				}
				files, err := scan.Walk(absRoot, cfg)
				if err != nil {
					return fmt.Errorf("scan: %w", err)
				}
				g = mapx.Build(files, cfg, gitmeta.AllCommitDates(absRoot))
				if anterior, err := mapx.Load(mapPath); err == nil {
					mapx.PreserveStamps(g, anterior)
				}
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
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "map path (default: <root>/"+mapx.DefaultPath+")")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "compile without writing (this is what the gate uses)")
	cmd.Flags().BoolVar(&semRebuild, "no-map-rebuild", false,
		"compile against the map on disk, without rebuilding it — what it does not know does not get in")
	// O CORTE entre a página completa e a resumida. Vale para o projeto inteiro: a
	// documentação em que uma camada segue uma lógica e a vizinha outra é a que obriga
	// quem lê a descobrir a lógica antes de achar o que procura.
	cmd.Flags().IntVar(&maxUnits, "max-units", padrao.MaxUnits,
		"above N units, the layer page carries the SUMMARY of each one")
	cmd.Flags().IntVar(&maxLines, "max-lines", padrao.MaxLines,
		"above N spec lines in total, likewise — five long specs weigh more than twenty short ones")
	return cmd
}

func printDocsResult(res doct.Result, dryRun bool) error {
	verbo := "written"
	if dryRun {
		verbo = "compiled (not written)"
	}
	for _, f := range res.Written {
		fmt.Printf("  ✓ %s/%s %s\n", doct.OutDir, f, verbo)
	}
	// Os ignorados NÃO são um detalhe: o template existe e o documento não foi gerado.
	// Silenciar isso faria o autor acreditar que a doc reflete as specs quando ela é outra
	// coisa — o mesmo defeito que o marcador existe para evitar, pela porta de trás.
	if len(res.Skipped) > 0 {
		fmt.Printf("\n  %d file(s) NOT generated — they exist in `%s/` without the marker:\n",
			len(res.Skipped), doct.OutDir)
		for _, f := range res.Skipped {
			fmt.Printf("    %s/%s\n", doct.OutDir, f)
		}
		fmt.Println("  They were written by hand, and the compiler does not erase anyone's work.")
		fmt.Println("  Either delete the `.tmpl` (the document is manual), or delete the `.md`")
		fmt.Println("  (the content is already in the specs and the template rebuilds it).")
	}
	if len(res.Written) == 0 && len(res.Skipped) == 0 {
		fmt.Printf("  no template in `%s/` — nothing to compile\n", doct.Dir)
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
		Short: "Write the skeleton of `" + doct.Dir + "/` (arquitetura, comportamento, regras, camadas)",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
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
				fmt.Printf("  · %s/%s already exists (use --force to overwrite)\n", doct.Dir, f)
			}
			fmt.Printf("\nEdit the templates and run `anchors docs build`.\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing templates")
	return cmd
}

// newDocsDutiesCmd lista as documentações que o projeto declarou como obrigatórias, e o
// que cada uma pede.
//
// É o comando que um agente roda antes de mexer numa unidade: "o que eu tenho de
// documentar, e o que essa documentação precisa responder". Sem ele, a instrução seria
// "atualize a doc", que produz o endpoint acrescentado ao OpenAPI sem os erros.
func newDocsDutiesCmd() *cobra.Command {
	var root, layer, unit string
	cmd := &cobra.Command{
		Use:   "duties",
		Short: "List the project's mandatory documents and what each requires",
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
			// A UNIDADE é a pergunta que o agente faz de verdade: ele não vai mexer em
			// "a camada infra", vai mexer numa unidade dela. E a camada erra sozinha —
			// no projeto de referência ela tem nove unidades e uma só toca esquema.
			if unit != "" {
				// A camada da UNIDADE, e não a do arquivo: uma spec casa `**/*.spec.md`
				// e o `ClassifyPath` devolveria `spec`, que não deve documentação
				// nenhuma. O card aponta a spec, e é o caminho que o agente usa.
				camada := scan.LayerOfUnit(absRoot, common.RelTo(absRoot, unit), cfg)
				docs = cfg.RequiredFor(camada, common.CodeOfUnit(absRoot, common.RelTo(absRoot, unit)))
				layer = unit
			} else if layer != "" {
				docs = cfg.RequiredFor(layer)
			}
			if len(docs) == 0 {
				if layer != "" {
					fmt.Printf("No documentation is required when changing `%s`.\n", layer)
					return nil
				}
				fmt.Println("The project declares no required documentation.")
				fmt.Println("Declare them under `docs:` in `" + config.DefaultFile + "` — kinds that")
				fmt.Println("Anchors knows how to instruct: " + strings.Join(doct.KnownKinds(), ", ") + ".")
				return nil
			}
			if layer != "" {
				fmt.Printf("Changing `%s` requires touching:\n\n", layer)
			} else {
				fmt.Println("Required documentation:")
				fmt.Println()
			}
			for _, d := range docs {
				fmt.Print(doct.Duty(d))
				fmt.Println()
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&layer, "layer", "", "only those required when changing this layer")
	cmd.Flags().StringVar(&unit, "unit", "", "only those required when changing THIS unit (more precise than --layer)")
	return cmd
}

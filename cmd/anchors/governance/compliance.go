package governance

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/pack"
	"github.com/spf13/cobra"
)

// `anchors compliance` responde a pergunta que nenhum comando respondia: **quais deveres
// se aplicam a este projeto, e onde estamos em relação a cada um.**
//
// Antes disto, a informação existia espalhada — o `check` dizia se UM arquivo violava UMA
// obrigação, e descobrir o estado de um regime inteiro exigia percorrer o repositório à
// mão. É justamente a pergunta que um auditor faz, e a que se responde na véspera.
//
// O relatório é por DEVER, não por arquivo: "Art. 17 — 42 nós sujeitos, 41 cumprem" diz
// algo; "arquivo X viola" repetido 42 vezes não diz.
func newComplianceCmd() *cobra.Command {
	var root string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "The state of each regulatory duty: how many nodes subject, how many comply",
		Long: `Shows, per DUTY, how many nodes are subject to it and how many comply.

The duties come from the adopted packs (` + "`packs:`" + ` in anchors.yaml) and from the
obligations declared inline. Each line cites the norm that originates the duty — that is what
makes the report presentable to whoever audits, and not only to whoever programs.

A duty with subject nodes and none complying is usually a disconnection (the target changed
path and the pack points at the old place), not 100% violation — and the report
signals that instead of letting you conclude wrongly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}

			packs, avisos, err := pack.LoadAll(absRoot, cfg.Packs, cfg.PackValues, cfg.Jurisdictions)
			if err != nil {
				return err
			}
			for _, a := range avisos {
				fmt.Fprintf(os.Stderr, "warning: %s\n", a)
			}

			obrigacoes := gate.ObligationsInForce(absRoot, cfg)
			if len(obrigacoes) == 0 {
				fmt.Println("No duty declared.")
				printAvailable(cfg)
				return nil
			}

			// de qual pack veio cada dever, para agrupar o relatório pela NORMA
			origem := map[string]*pack.Pack{}
			artigo := map[string]string{}
			for _, p := range packs {
				for _, ob := range p.Obligations {
					origem[ob.Name] = p
					artigo[ob.Name] = ob.Article
				}
			}

			if len(cfg.Jurisdictions) > 0 {
				fmt.Printf("declared jurisdictions: %s\n\n", strings.Join(cfg.Jurisdictions, ", "))
			}

			res := gate.EvaluateObligations(absRoot, cfg, g, obrigacoes)
			porNorma := map[string][]gate.ObligationStatus{}
			for _, r := range res {
				norma := "declared in the project"
				if p, ok := origem[r.Name]; ok {
					norma = p.Authority
					if norma == "" {
						norma = p.Name
					}
				}
				porNorma[norma] = append(porNorma[norma], r)
			}

			var normas []string
			for n := range porNorma {
				normas = append(normas, n)
			}
			sort.Strings(normas)

			totalSujeitos, totalCumpre := 0, 0
			for _, norma := range normas {
				fmt.Printf("── %s\n", norma)
				linhas := porNorma[norma]
				sort.Slice(linhas, func(i, j int) bool { return linhas[i].Name < linhas[j].Name })
				for _, r := range linhas {
					marca := "✓"
					switch {
					case r.Subject == 0:
						marca = "·" // nenhum nó dispara este dever
					case r.Fulfilled < r.Subject:
						marca = "✗"
					}
					art := artigo[r.Name]
					if art != "" {
						art = " " + art
					}
					fmt.Printf("   %s %-26s%-14s %3d subject, %3d complying",
						marca, r.Name, art, r.Subject, r.Fulfilled)
					if r.Debt > 0 {
						fmt.Printf(", %d assumed debt(s)", r.Debt)
					}
					if r.Waived > 0 {
						fmt.Printf(", %d waiver(s)", r.Waived)
					}
					fmt.Println()
					// Zero cumprindo com nós sujeitos quase nunca é 100% de violação: é o
					// alvo apontando para caminho que não existe mais. Dizer isso evita a
					// conclusão errada — e a correção é uma linha de `pack_values`.
					if r.Subject > 0 && r.Fulfilled == 0 && r.Debt == 0 {
						fmt.Printf("       ⚠ NONE complies — check whether `%s` is still the right path "+
							"in `pack_values:`; a disconnected target looks like total violation\n",
							strings.Join(r.Targets, ", "))
					}
					if verbose && len(r.Missing) > 0 {
						for _, m := range r.Missing {
							fmt.Printf("       ✗ %s\n", m)
						}
					}
					totalSujeitos += r.Subject
					totalCumpre += r.Fulfilled
				}
				fmt.Println()
			}

			fmt.Printf("total: %d subject duty(ies) across %d node(s), %d fulfilled\n",
				len(res), totalSujeitos, totalCumpre)
			if !verbose && totalCumpre < totalSujeitos {
				fmt.Println("(use --verbose to see which nodes are missing in each duty)")
			}
			printAvailable(cfg)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "lists the nodes missing in each duty")
	return cmd
}

// printAvailable mostra os packs que existem no disco e o projeto NÃO adotou. Não é
// sugestão de adotar — é a diferença entre "não se aplica a mim" e "esqueci", que só o
// projeto sabe qual é, mas precisa poder ver.
func printAvailable(cfg *config.Config) {
	adotado := map[string]bool{}
	for _, p := range cfg.Packs {
		adotado[strings.TrimSuffix(strings.TrimPrefix(p, "./packs/"), ".yaml")] = true
	}
	var falta []string
	for _, nomes := range initx.AvailablePacks() {
		for _, n := range nomes {
			if !adotado[n] {
				falta = append(falta, n)
			}
		}
	}
	if len(falta) == 0 {
		return
	}
	sort.Strings(falta)
	fmt.Printf("\npacks available and not adopted: %s\n", strings.Join(falta, ", "))
	fmt.Println("  (if any applies to your product, declare it in `packs:` in anchors.yaml)")
}

package main

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
		Short: "O estado de cada dever regulatório: quantos nós sujeitos, quantos cumprem",
		Long: `Mostra, por DEVER, quantos nós estão sujeitos a ele e quantos o cumprem.

Os deveres vêm dos packs adotados (` + "`packs:`" + ` no anchors.yaml) e das obrigações
declaradas inline. Cada linha cita a norma que origina o dever — é o que torna o
relatório apresentável a quem audita, em vez de só a quem programa.

Um dever com nós sujeitos e nenhum cumprindo costuma ser desconexão (o alvo mudou de
caminho e o pack aponta para o lugar antigo), não 100% de violação — e o relatório
sinaliza isso em vez de deixar você concluir errado.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar config: %w", err)
			}
			g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}

			packs, avisos, err := pack.LoadAll(absRoot, cfg.Packs, cfg.PackValues, cfg.Jurisdictions)
			if err != nil {
				return err
			}
			for _, a := range avisos {
				fmt.Fprintf(os.Stderr, "aviso: %s\n", a)
			}

			obrigacoes := gate.ObligationsInForce(absRoot, cfg)
			if len(obrigacoes) == 0 {
				fmt.Println("Nenhum dever declarado.")
				printDisponiveis(cfg)
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
				fmt.Printf("jurisdições declaradas: %s\n\n", strings.Join(cfg.Jurisdictions, ", "))
			}

			res := gate.EvaluateObligations(absRoot, cfg, g, obrigacoes)
			porNorma := map[string][]gate.ObligationStatus{}
			for _, r := range res {
				norma := "declarados no projeto"
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
					fmt.Printf("   %s %-26s%-14s %3d sujeito(s), %3d cumpre(m)",
						marca, r.Name, art, r.Subject, r.Fulfilled)
					if r.Debt > 0 {
						fmt.Printf(", %d dívida(s) assumida(s)", r.Debt)
					}
					if r.Waived > 0 {
						fmt.Printf(", %d dispensa(s)", r.Waived)
					}
					fmt.Println()
					// Zero cumprindo com nós sujeitos quase nunca é 100% de violação: é o
					// alvo apontando para caminho que não existe mais. Dizer isso evita a
					// conclusão errada — e a correção é uma linha de `pack_values`.
					if r.Subject > 0 && r.Fulfilled == 0 && r.Debt == 0 {
						fmt.Printf("       ⚠ NENHUM cumpre — confira se `%s` ainda é o caminho certo "+
							"em `pack_values:`; alvo desconectado parece violação total\n",
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

			fmt.Printf("total: %d dever(es) sujeito(s) em %d nó(s), %d cumprido(s)\n",
				len(res), totalSujeitos, totalCumpre)
			if !verbose && totalCumpre < totalSujeitos {
				fmt.Println("(use --verbose para ver quais nós faltam em cada dever)")
			}
			printDisponiveis(cfg)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "lista os nós que faltam em cada dever")
	return cmd
}

// printDisponiveis mostra os packs que existem no disco e o projeto NÃO adotou. Não é
// sugestão de adotar — é a diferença entre "não se aplica a mim" e "esqueci", que só o
// projeto sabe qual é, mas precisa poder ver.
func printDisponiveis(cfg *config.Config) {
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
	fmt.Printf("\npacks disponíveis e não adotados: %s\n", strings.Join(falta, ", "))
	fmt.Println("  (se algum se aplica ao seu produto, declare em `packs:` no anchors.yaml)")
}

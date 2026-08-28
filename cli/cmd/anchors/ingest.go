package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/co2-lab/anchors/internal/config"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testsig"
	"github.com/spf13/cobra"
)

// `anchors ingest` consome os artefatos que o RUNNER do projeto já gera (o Anchors
// não roda o teste) e os amarra ao grafo: execução (JUnit → passou/falhou) e
// cobertura de linha (lcov → % por arquivo). Daí deriva a cobertura por CENÁRIO
// (cada código de cenário da spec tem um teste que PASSOU?). O sinal leva a rev do
// nó — fica stale se o arquivo mudar depois (não confie em teste de versão antiga).
func newIngestCmd() *cobra.Command {
	var root, mapPath, junit, lcov, mutation, layer, scope string
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingere sinais de teste (JUnit/lcov) que o projeto gerou e os amarra ao mapa",
		Long: `Consome os artefatos do test runner e grava os sinais no mapa:

  anchors ingest --junit results.xml      execução: quais testes passaram/falharam
  anchors ingest --lcov coverage.info     cobertura de linha por arquivo de código
  anchors ingest --mutation mutation.json  mutação: o teste PROVA a linha, ou só a executa?
  anchors ingest --mutation m.json --scope isolated   só o teste da unidade
  anchors ingest --mutation m.json --scope full       com os testes dos dependentes
  anchors ingest --junit r.xml --lcov c.info   vários numa passada

O Anchors NÃO roda o teste — você roda (jest --coverage, go test -coverprofile,
pytest-cov…) e passa o relatório. Da execução, o Anchors também deriva a cobertura por
CENÁRIO: um código de cenário (SPCRX-V01) está PROVADO se aparece num caso que passou.
Rode 'anchors coverage' depois para ver os requisitos de spec sem teste verde.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if junit == "" && lcov == "" && mutation == "" {
				return fmt.Errorf("informe --junit <arquivo>, --lcov <arquivo> e/ou --mutation <arquivo>")
			}
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			return ingereArtefatos(absRoot, mapPath, junit, lcov, mutation, layer, scope)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().StringVar(&junit, "junit", "", "arquivo JUnit XML (resultado de execução)")
	cmd.Flags().StringVar(&lcov, "lcov", "", "arquivo lcov .info (cobertura de linha)")
	cmd.Flags().StringVar(&mutation, "mutation", "", "relatório JSON de mutação. O FORMATO é do projeto: declare no gate mutation-score do anchors.yaml, chave 'format' (default 'mutation-testing-elements' — Stryker/PIT/Infection/mutmut; 'gremlins' para Go)")
	cmd.Flags().StringVar(&layer, "layer", "", "camada de teste desta suite (unit|integration|e2e…); default unit — mergeia várias")
	cmd.Flags().StringVar(&scope, "scope", "", "escopo da suíte que rodou os mutantes: `isolated` (só o teste da unidade) ou `full` (com os dependentes). Ingerir os dois permite ler a DIFERENÇA — quanto a unidade depende de terceiros para se provar")
	return cmd
}

// ingereArtefatos é o miolo da ingestão, separado do comando para que `anchors test` e
// `anchors mutation` possam ingerir o que acabaram de produzir sem reimplementar nada
// nem invocar o próprio binário de novo. É o que fecha o par "rodar" / "ingerir" que
// antes exigia um humano no meio.
func ingereArtefatos(absRoot, mapPath, junit, lcov, mutation, layer, scope string) error {
	{
		{
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}
			// O vocabulário de letras é do PROJETO: sem isto, um cenário de letra
			// declarada (ex.: `-I01`) não é reconhecido no nome do caso JUnit, e o
			// requisito aparece sem teste verde mesmo tendo um teste que passa.
			// O formato do relatório de MUTAÇÃO é do PROJETO, declarado junto do gate
			// (`gates: - name: mutation-score / format: gremlins`). Vazio = o canônico
			// Mutation Testing Elements, que é o default para todo projeto existente.
			mutationFormat := config.FormatMTE
			if cfg, cerr := config.Load(filepath.Join(absRoot, config.DefaultFile)); cerr == nil {
				testsig.SetRuleLetters(cfg.RuleLetters())
				mutationFormat = cfg.MutationFormat()
			}
			now := time.Now().Format(time.RFC3339)

			if junit != "" {
				rep, err := testsig.ParseJUnit(junit)
				if err != nil {
					return fmt.Errorf("parse JUnit: %w", err)
				}
				byFile := map[string]mapx.ExecByFile{}
				for _, c := range rep.Cases {
					if c.File == "" {
						continue // sem arquivo, não dá para amarrar ao nó
					}
					e := byFile[c.File]
					switch {
					case c.Failed:
						e.Failed++
					case c.Skipped:
						e.Skipped++
					default:
						e.Passed++
					}
					byFile[c.File] = e
				}
				proven := rep.PassedCodes()
				// monta os cenários DECLARADOS por cada spec (lendo o arquivo — o
				// mapx não toca disco), para o cruzamento de cobertura semântica.
				declaredByNode := map[string][]string{}
				for _, n := range g.Nodes {
					if n.Kind == mapx.KindSpec {
						if codes, err := codesInFile(filepath.Join(absRoot, n.ID)); err == nil && len(codes) > 0 {
							declaredByNode[n.ID] = codes
						}
					}
				}
				mf, mc := g.IngestExecution(byFile, proven, declaredByNode, layer, now)
				fmt.Printf("execução: %d caso(s), %d arquivo(s) de teste casados, %d cenário(s) provado(s)\n",
					len(rep.Cases), mf, mc)
				if len(byFile) > 0 && mf == 0 {
					fmt.Println("  aviso: nenhum arquivo de teste casou — o JUnit tem o atributo 'file'? (use um reporter que o emita)")
				}
			}

			if lcov != "" {
				rep, err := testsig.ParseLCOV(lcov)
				if err != nil {
					return fmt.Errorf("parse lcov: %w", err)
				}
				byFile := map[string]mapx.FileCov{}
				for _, fc := range rep.Files {
					byFile[fc.File] = mapx.FileCov{Covered: fc.CoveredLines, Total: fc.TotalLines}
				}
				m := g.IngestCoverage(byFile, now)
				fmt.Printf("cobertura: %d arquivo(s) no lcov, %d nó(s) de código casados\n", len(rep.Files), m)
			}

			if mutation != "" {
				rep, err := testsig.ParseMutationFormat(mutation, mutationFormat)
				if err != nil {
					return fmt.Errorf("parse mutação: %w", err)
				}
				byFile := map[string]mapx.FileMutation{}
				sobreviventes := 0
				for file, fm := range rep.Files {
					byFile[file] = mapx.FileMutation{
						Killed: fm.Killed, Survived: fm.Survived,
						NoCoverage: fm.NoCoverage, Ignored: fm.Ignored, Score: fm.Score,
					}
					sobreviventes += fm.Survived
				}
				m := g.IngestMutationScoped(byFile, scope, now, rep.Low, rep.High)
				fmt.Printf("mutação (%s): %d arquivo(s) no relatório, %d nó(s) de código casados, %d mutante(s) sobrevivente(s)\n",
					mutationFormat, len(rep.Files), m, sobreviventes)
				if sobreviventes > 0 {
					fmt.Println("  cada sobrevivente é uma alteração no código que os testes NÃO perceberam.")
				}
			}

			if err := mapx.Save(g, mapPath); err != nil {
				return fmt.Errorf("salvar mapa: %w", err)
			}
			fmt.Println("sinais gravados no mapa. Veja com `anchors coverage`.")
			return nil
		}
	}
}

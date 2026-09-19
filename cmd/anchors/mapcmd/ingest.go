package mapcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
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
// common.ViaAnchorsTest é ligado pelo `anchors test`/`anchors mutation` antes de ingerir.
//
// A distinção é o ponto: aqueles comandos rodam a suíte E ingerem numa operação só, então
// o sinal corresponde à execução que acabou de acontecer. Quem chama `ingest` direto pode
// estar ingerindo um relatório de uma hora atrás — e o mapa passa a afirmar uma cobertura
// que já não vale.
var ViaAnchorsTest bool

func newIngestCmd() *cobra.Command {
	var root, mapPath, junit, lcov, mutation, layer, scope string
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest test signals (JUnit/lcov) the project generated and bind them to the map",
		Long: `Consumes the test runner's artifacts and writes the signals into the map:

  anchors ingest --junit results.xml      execution: which tests passed/failed
  anchors ingest --lcov coverage.info     line coverage per code file
  anchors ingest --mutation mutation.json  mutation: does the test PROVE the line, or only run it?
  anchors ingest --mutation m.json --scope isolated   the unit's test only
  anchors ingest --mutation m.json --scope full       with the dependents' tests
  anchors ingest --junit r.xml --lcov c.info   several in one pass

Anchors does NOT run the test — you run it (jest --coverage, go test -coverprofile,
pytest-cov…) and hand over the report. From the execution, Anchors also derives coverage per
SCENARIO: a scenario code (SPCRX-V01) is PROVEN if it appears in a case that passed.
Run 'anchors coverage' afterwards to see the spec requirements with no green test.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if junit == "" && lcov == "" && mutation == "" {
				return fmt.Errorf("provide --junit <file>, --lcov <file> and/or --mutation <file>")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			return IngestArtifacts(absRoot, mapPath, junit, lcov, mutation, layer, scope)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringVar(&junit, "junit", "", "JUnit XML file (execution result)")
	cmd.Flags().StringVar(&lcov, "lcov", "", "lcov .info file (line coverage)")
	cmd.Flags().StringVar(&mutation, "mutation", "", "mutation JSON report. The FORMAT belongs to the project: declare it in the mutation-score gate of anchors.yaml, key 'format' (default 'mutation-testing-elements' — Stryker/PIT/Infection/mutmut; 'gremlins' for Go)")
	cmd.Flags().StringVar(&layer, "layer", "", "test layer of this suite (unit|integration|e2e…); default unit — merges several")
	cmd.Flags().StringVar(&scope, "scope", "", "scope of the suite that ran the mutants: `isolated` (only the unit's test) or `full` (with the dependents). Ingesting both allows reading the DIFFERENCE — how much the unit depends on third parties to prove itself")
	return cmd
}

// IngestArtifacts é o miolo da ingestão, separado do comando para que `anchors test` e
// `anchors mutation` possam ingerir o que acabaram de produzir sem reimplementar nada
// nem invocar o próprio binário de novo. É o que fecha o par "rodar" / "ingerir" que
// antes exigia um humano no meio.
func IngestArtifacts(absRoot, mapPath, junit, lcov, mutation, layer, scope string) error {
	if err := warnIfManualIngest(absRoot); err != nil {
		return err
	}
	{
		{
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
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
						if codes, err := common.CodesInFileOfUnit(filepath.Join(absRoot, n.ID), n.Code); err == nil && len(codes) > 0 {
							declaredByNode[n.ID] = codes
						}
					}
				}
				mf, mc := g.IngestExecution(byFile, proven, declaredByNode, layer, now)
				fmt.Printf("execution: %d case(s), %d test file(s) matched, %d scenario(s) proven\n",
					len(rep.Cases), mf, mc)
				if len(byFile) > 0 && mf == 0 {
					fmt.Println("  warning: no test file matched — does the JUnit have the 'file' attribute? (use a reporter that emits it)")
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
				fmt.Printf("coverage: %d file(s) in the lcov, %d code node(s) matched\n", len(rep.Files), m)
			}

			if mutation != "" {
				rep, err := testsig.ParseMutationFormat(mutation, mutationFormat)
				if err != nil {
					return fmt.Errorf("parse mutation: %w", err)
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
				fmt.Printf("mutation (%s): %d file(s) in the report, %d code node(s) matched, %d surviving mutant(s)\n",
					mutationFormat, len(rep.Files), m, sobreviventes)
				if sobreviventes > 0 {
					fmt.Println("  each survivor is a change to the code that the tests did NOT notice.")
				}
			}

			if err := mapx.Save(g, mapPath); err != nil {
				return fmt.Errorf("save map: %w", err)
			}
			fmt.Println("signals written into the map. See them with `anchors coverage`.")
			return nil
		}
	}
}

// warnIfManualIngest reclama de uma ingestão feita fora do `anchors test`.
//
// AVISA por padrão e só BARRA quando o projeto declara `manual_ingest_blocks: true`. A
// razão de não barrar sempre é que há usos legítimos — um CI que rodou a suíte noutro
// job, uma ferramenta que o `tests:` não cobre —, e derrubá-los tiraria a saída de quem
// tem razão.
func warnIfManualIngest(absRoot string) error {
	if ViaAnchorsTest {
		return nil
	}
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		return nil // sem config não há o que exigir
	}
	// Sem `tests:` declarado o `anchors test` não roda, e exigir que se use o que não
	// existe seria mandar o projeto para um beco sem saída.
	if len(cfg.Tests) == 0 && len(cfg.Mutation) == 0 {
		return nil
	}
	if cfg.Workflow.IngestManualBarra() {
		return fmt.Errorf("MANUAL ingestion refused — this project declares " +
			"`manual_ingest_blocks: true`.\n" +
			"  Use `anchors test` (or `anchors mutation`): it runs the suite AND ingests in a " +
			"single operation, and that is what guarantees the signal in the map corresponds to the last " +
			"run.\n" +
			"  Ingesting by hand, you can run the suite, edit the code and ingest the " +
			"old report — and the map starts asserting a coverage that no longer holds")
	}
	fmt.Fprintln(os.Stderr, "⚠ MANUAL ingestion: the signal may not correspond to the last run.")
	fmt.Fprintln(os.Stderr, "  `anchors test` runs the suite and ingests in a single operation — that is what")
	fmt.Fprintln(os.Stderr, "  guarantees the map reflects what has just run. To REFUSE manual")
	fmt.Fprintln(os.Stderr, "  ingestion, declare `manual_ingest_blocks: true` under `workflow:`.")
	return nil
}

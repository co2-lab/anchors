package mapcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"

	"github.com/co2-lab/anchors/internal/logscan"
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
	var logs bool
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
			if junit == "" && lcov == "" && mutation == "" && !logs {
				return fmt.Errorf("provide --junit <file>, --lcov <file>, --mutation <file> and/or --logs")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if junit != "" || lcov != "" || mutation != "" {
				if err := IngestArtifacts(absRoot, mapPath, junit, lcov, mutation, layer, scope); err != nil {
					return err
				}
			}
			if logs {
				return ingestLogs(absRoot, mapPath)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&logs, "logs", false,
		"scan the logs declared in `logs.paths` and bind the failure occurrences to the specs that declare them")
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
					// SPEC e FLAG: os dois declaram cenarios que um teste pode provar.
					//
					// A flag entrou depois, e sem ela os cenarios de flag (`-G`) nunca
					// chegavam ao mapa: o teste nomeava o codigo, o relatorio o trazia, e o
					// cruzamento o descartava por o no' nao ser spec. O `flag-covered`
					// respondia "nenhuma execucao ingerida" com a suite verde na mao.
					if n.Kind != mapx.KindSpec && n.Kind != mapx.KindFlag {
						continue
					}
					if codes, err := common.CodesInFileOfUnit(filepath.Join(absRoot, n.ID), n.Code); err == nil && len(codes) > 0 {
						declaredByNode[n.ID] = codes
					}
				}
				byFile = resolveByFile(g, mapx.KindTest, byFile, absRoot, junit)
				// A suíte é o RELATÓRIO: num monorepo cada workspace ingere o seu, e sem a
				// chave a prova de uma suíte apagava a das outras (ver ProvenBySuite).
				mf, mc := g.IngestExecutionSuite(byFile, proven, declaredByNode, layer, relToRoot(absRoot, junit), now)
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
				byFile = resolveByFile(g, mapx.KindCode, byFile, absRoot, lcov)
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
				byFile = resolveByFile(g, mapx.KindCode, byFile, absRoot, mutation)
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
// resolveByFile reescreve as chaves do relatório para o ID exato do nó dono (ver
// mapx.ResolveReportPaths) e descarta, com aviso, as que continuam ambíguas — caminho
// relativo ao workspace que casa nós de dois workspaces e o diretório do relatório não
// desempata. Sem isto os DOIS nós recebiam o sinal.
func resolveByFile[T any](g *mapx.Graph, kind mapx.Kind, byFile map[string]T, absRoot, report string) map[string]T {
	paths := make([]string, 0, len(byFile))
	for p := range byFile {
		paths = append(paths, p)
	}
	resolved, ambiguous := g.ResolveReportPaths(kind, paths, relToRoot(absRoot, report))
	out := make(map[string]T, len(byFile))
	for p, v := range byFile {
		if id, ok := resolved[p]; ok {
			out[id] = v
		}
	}
	if len(ambiguous) > 0 {
		sort.Strings(ambiguous)
		fmt.Printf("  warning: %d path(s) in %s match more than one node and the report's folder does not decide — left WITHOUT signal:\n", len(ambiguous), filepath.Base(report))
		for _, a := range ambiguous {
			fmt.Printf("    %s\n", a)
		}
		fmt.Println("  make the runner write paths from the repository root (lcov: projectRoot; jest-junit: filePathPrefix).")
	}
	return out
}

// relToRoot devolve o caminho do relatório relativo à raiz do projeto, com `/` — é a
// identidade da suíte no mapa, e precisa ser a mesma em qualquer máquina.
func relToRoot(absRoot, report string) string {
	rel := report
	if abs, err := filepath.Abs(report); err == nil {
		if r, err := filepath.Rel(absRoot, abs); err == nil {
			rel = r
		}
	}
	return filepath.ToSlash(rel)
}

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

// ingestLogs varre os logs do projeto e amarra as ocorrências às specs que as declaram.
//
// O Anchors VARRE — ele não pede o resultado pronto, ao contrário do que a primeira versão
// deste eixo fazia. O que torna isso possível sem ditar formato é que o log carrega o
// CÓDIGO da falha: `CRED-E01` é a mesma sequência de caracteres em JSON, em texto puro ou
// em syslog, e procurar o código dispensa entender o formato.
func ingestLogs(absRoot, mapPath string) error {
	if mapPath == "" {
		mapPath = filepath.Join(absRoot, mapx.DefaultPath)
	}
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		return err
	}
	if cfg.Logs == nil || len(cfg.Logs.Paths) == 0 {
		return fmt.Errorf("no log declared — add `logs.paths` to anchors.yaml with the globs of your log files")
	}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return err
	}
	res, err := logscan.Scan(absRoot, cfg, g)
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}

	// A ocorrência vai para a SPEC que declara a regra — é dela a falha, e é nela que a
	// conclusão sobre a causa vai ser escrita depois.
	byRule := map[string]*mapx.Node{}
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(absRoot, n.ID))
		if err != nil {
			continue
		}
		for _, m := range logscan.SpecFailureCodes(string(b)) {
			byRule[m] = n
		}
		n.Failures = nil
	}
	bound := 0
	for _, o := range res.Occurrences {
		n, ok := byRule[o.Rule]
		if !ok {
			continue
		}
		// A REV da spec no momento da ingestão: a ocorrência envelhece se a spec mudar,
		// porque a falha observada era da versão anterior da regra.
		o.Rev = n.Rev
		n.Failures = append(n.Failures, o)
		bound++
	}
	if err := mapx.Save(g, mapPath); err != nil {
		return err
	}

	fmt.Printf("logs: %d file(s), %d line(s) — %d occurrence(s) bound to %d spec rule(s)\n",
		res.Files, res.Lines, bound, len(res.Occurrences))
	if len(res.Unknown) == 0 {
		return nil
	}
	// O ACHADO da camada, e o que nenhuma observability dá: o log mostra um erro que a
	// spec não previu. Ele não entra no mapa (amarrá-lo inventaria dono para uma falha
	// órfã), mas some em silêncio se não for dito aqui.
	fmt.Printf("\n⚠ %d failure code(s) in the log that NO spec declares:\n", len(res.Unknown))
	codes := make([]string, 0, len(res.Unknown))
	for c := range res.Unknown {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	for _, c := range codes {
		fmt.Printf("    %-12s %d occurrence(s)\n", c, res.Unknown[c])
	}
	fmt.Println("  Somebody is handling and logging a failure the spec never declared —")
	fmt.Println("  or the code is a typo. Both are worth a look.")
	return nil
}

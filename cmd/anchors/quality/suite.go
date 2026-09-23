package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// `anchors test` e `anchors mutation` são PROXIES: o comando é do projeto, declarado no
// anchors.yaml, e o Anchors só o executa e amarra o relatório ao mapa.
//
// Por que existem, se o projeto já sabe rodar `yarn test`: porque produzir sinal eram
// dois passos, e o segundo era esquecível. Quem roda a ferramenta e não chama
// `anchors ingest` vê o gate continuar acusando "sem sinal ingerido", como se a rodada
// não tivesse acontecido — falha silenciosa e no lado errado (parece pendência de
// qualidade, é pendência de processo). Medido num projeto real: 1008 das 1011
// pendências do `check --all` eram sinal ausente, não teste ruim.
//
// O que o Anchors NÃO faz aqui, deliberadamente: conhecer a stack. Ele não sabe o que é
// jest, pytest ou go test, não monta linha de comando, não interpreta saída. Se o
// projeto não declarou, ele não adivinha — mostra como declarar e sai. A alternativa
// (embutir convenções de runner) seria o framework decidindo a stack do projeto, que é
// exatamente o que os gates externos já provaram desnecessário: o `run:` é do projeto e
// o Anchors o executa sem saber o que é eslint.

func newTestCmd() *cobra.Command {
	return newSuiteCommand(suiteCommand{
		nome:  "test",
		curto: "Run the test suites declared in anchors.yaml and ingest the reports",
		secao: "tests",
		exemplo: `tests:
  - workspace: backend
    layer: unit
    run: "yarn test:unit"
    run_changed: "yarn test:unit --findRelatedTests {{files}}"
    junit: "packages/backend/test-output/junit/unit.xml"
    lcov: "packages/backend/test-output/coverage/lcov.info"
  - layer: integration
    run: "yarn test:integration"
    junit: "packages/backend/test-output/junit/integration.xml"
  - layer: e2e
    run: "yarn e2e"
    junit: "apps/mobile/test-output/junit/e2e.xml"`,
		usoLongo: `Runs what the PROJECT declared in ` + "`tests:`" + ` and ingests what the run left behind.

  anchors test                    every declared layer, in the file's order
  anchors test unit               unit only
  anchors test unit integration   more than one, in the file's order
  anchors test -w backend         filters by workspace (monorepo)
  anchors test --changed a.ts     INCREMENTAL: only the impact path of what changed
  anchors test --then check       on passing, charges the gates with the signal still fresh

Two modes, the same as check's: FULL (without --changed) runs the "run:"; INCREMENTAL
runs the "run_changed:", which receives in {{files}} the files of the impact path -
the SAME cut that "check --changed" uses, so the two do not disagree about what
your commit moves.

Anchors does not know how to run a test — the stack is the project's. It executes the declared ` + "`run:`" + `
and binds the ` + "`junit:`" + `/` + "`lcov:`" + ` to the map, which is the step that usually gets forgotten.`,
	})
}

func newMutationCmd() *cobra.Command {
	return newSuiteCommand(suiteCommand{
		nome:  "mutation",
		curto: "Run the mutation suites declared in anchors.yaml and ingest the reports",
		secao: "mutation",
		exemplo: `mutation:
  - layer: unit
    run: "cd packages/backend && npx stryker run --mutate {{target}}"
    run_changed: "cd packages/backend && npx stryker run --mutate {{files}}"
    report: "packages/backend/reports/mutation/mutation.json"
    scope: full`,
		usoLongo: `Runs what the PROJECT declared in ` + "`mutation:`" + ` and ingests the report.

  anchors mutation                          every declared layer
  anchors mutation unit                     unit only
  anchors mutation unit --target x/y.ts     fills {{target}} in the declared command
  anchors mutation --changed x/y.ts         INCREMENTAL: mutates the impact path
  anchors mutation unit --then check        on passing, charges the gates next

In incremental mode, {{files}} receives only the CODE nodes: mutation alters the rule, and mutating
the test would invert the experiment - the test is the measuring instrument, not the object.

Mutation answers what coverage does not: change the line — does the test notice? A
SURVIVING mutant is a line nobody proves.

Mind the ` + "`scope:`" + `: with ` + "`isolated`" + ` and ` + "`full`" + ` both ingested, the gate judges by the
ISOLATED one, which is far harsher — it is the one that says whether the unit's test proves the unit,
instead of it being proven by the dependents.`,
	})
}

// suiteCommand descreve o que muda entre `test` e `mutation` — o resto é idêntico, e
// duplicar os dois faria a mensagem de "não configurado" divergir com o tempo.
type suiteCommand struct {
	nome     string
	curto    string
	secao    string // o nome da seção no anchors.yaml, usado nas mensagens
	exemplo  string // o bloco YAML que se mostra a quem não configurou
	usoLongo string
}

func newSuiteCommand(cs suiteCommand) *cobra.Command {
	var root, target, then string
	var workspaces, changed, escopos []string
	cmd := &cobra.Command{
		Use:   cs.nome + " [layers...]",
		Short: cs.curto,
		Long:  cs.usoLongo,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load %s: %w", config.DefaultFile, err)
			}

			declaradas := cfg.Tests
			if cs.secao == "mutation" {
				declaradas = cfg.Mutation
			}
			if len(declaradas) == 0 {
				printHowToConfigure(cs)
				return fmt.Errorf("no suite declared in `%s:` in %s", cs.secao, config.DefaultFile)
			}

			sel, ausentes := config.SelecionaSuites(declaradas, args, workspaces, escopos)
			if len(ausentes) > 0 {
				return fmt.Errorf("not declared in `%s:`: %s\n  declared layers:     %s\n  declared workspaces: %s\n  declared scopes:     %s",
					cs.secao, strings.Join(ausentes, ", "),
					strings.Join(config.DeclaredLayers(declaradas), ", "),
					joinOrDash(config.DeclaredWorkspaces(declaradas)),
					joinOrDash(config.DeclaredScopes(declaradas)))
			}
			// Combinação válida mas vazia não é erro de digitação — é "não existe essa
			// suíte". Dizer isso com clareza evita a leitura de que o comando rodou e
			// tudo passou, que é o que um silêncio com exit 0 comunicaria.
			if len(sel) == 0 {
				return fmt.Errorf("no suite matches layer(s) %s, workspace(s) %s and scope(s) %s in `%s:`",
					joinOrDash(args), joinOrDash(workspaces), joinOrDash(escopos), cs.secao)
			}

			// O caminho de impacto sai da MESMA função que o `check --changed` usa: se
			// as duas divergissem, "os gates que o meu commit move" e "os testes que o
			// meu commit move" passariam a ser conjuntos diferentes, e o incremental
			// deixaria de ser confiável exatamente onde ele é a única defesa.
			var alvos []string
			if len(changed) > 0 {
				g, mapErr := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
				if mapErr != nil {
					return fmt.Errorf("load map: %w (run `anchors map build`)", mapErr)
				}
				nodes, _, selErr := selectNodes(g, cfg, false, changed, absRoot)
				if selErr != nil {
					return selErr
				}
				alvos = impactFiles(nodes, cs.secao, absRoot)
				if len(alvos) == 0 {
					fmt.Printf("the impact path reaches no %s file — nothing to run.\n",
						map[bool]string{true: "code", false: "code or test"}[cs.secao == "mutation"])
					return nil
				}
				fmt.Printf("incremental: %d file(s) on the impact path\n\n", len(alvos))
			}

			if err := runSuites(cs, sel, absRoot, target, alvos); err != nil {
				return err
			}
			// O encadeamento é OPT-IN e só acontece depois do sucesso: um `check` sobre
			// sinal que não foi produzido diria o mesmo de antes, e um sobre suíte que
			// falhou culparia o gate por um teste vermelho.
			return runChained(then, absRoot)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringSliceVarP(&workspaces, "workspace", "w", nil, "filters by the declared WORKSPACE (backend, mobile…) — repeatable or comma-separated. Combines with the layers: `anchors "+cs.nome+" unit -w backend`")
	// Sem `--scope`, roda TODOS os escopos declarados — e é o default certo: é o par
	// isolado+completo que produz a leitura de acoplamento do gate, e pedir dois comandos
	// para obtê-la devolveria o passo manual que estes comandos vieram eliminar.
	cmd.Flags().StringSliceVar(&escopos, "scope", nil, "filters by the declared SCOPE (`isolated`, `full`) — only on mutation. Without this, runs whichever are declared")
	cmd.Flags().StringVar(&target, "target", "", "target that replaces `{{target}}` in the declared `run:` (e.g. the file to mutate)")
	cmd.Flags().StringSliceVar(&changed, "changed", nil, "INCREMENTAL mode: changed file(s) — runs the `run_changed:` over the union of the impact paths, the same slice as `check --changed`")
	cmd.Flags().StringVar(&then, "then", "", "on PASSING, chains Anchors commands: `check`, `coverage` (separate by comma). Opt-in — without this, it runs and stops")
	return cmd
}

// runSuites executa cada suíte e ingere o que ela deixou. Para na primeira que falhar:
// as camadas costumam depender umas das outras (não faz sentido rodar e2e depois de a
// unit quebrar), e seguir adiante só produziria ruído sobre uma base já vermelha.
func runSuites(cs suiteCommand, suites []config.Suite, absRoot, target string, alvos []string) error {
	for _, s := range suites {
		linha, err := pickCommand(s, alvos, target)
		if err != nil {
			return fmt.Errorf("layer %q: %w", s.Layer, err)
		}
		fmt.Printf("━━━ %s [%s%s%s] ━━━\n%s\n\n", cs.nome, workspaceLabel(s), s.Layer, scopeLabel(s), linha)

		if teto := suiteArgvLimit(); len(linha) > teto {
			return fmt.Errorf("layer %q: the command line came out with %d characters (ceiling %d on this platform).\n"+
				"  The impact path is too big for a single invocation. Run the full mode,\n"+
				"  narrow the --changed, or adjust ANCHORS_ARGV_MAX if you know your shell can take it",
				s.Layer, len(linha), teto)
		}

		inicio := time.Now()
		errRun := execAtRoot(linha, absRoot)

		junit, lcov, mutation := absPath(absRoot, s.JUnit), absPath(absRoot, s.Lcov), absPath(absRoot, s.Report)
		// A ingestão acontece MESMO se o comando saiu != 0, e essa é a regra menos
		// óbvia daqui. Um runner sai != 0 exatamente quando há o que reportar: o jest
		// quando um teste falha, o Stryker quando o score fica abaixo do próprio
		// `thresholds.break`. Pular a ingestão nesse caso deixaria justamente a unidade
		// PROBLEMÁTICA sem sinal no mapa — o gate seguiria dizendo "sem sinal ingerido"
		// sobre a única que já se sabe ruim, e o vermelho pareceria ausência de medida.
		//
		// O relatório é a evidência do que aconteceu, não um prêmio por ter passado.
		if err := ingestIfRecent(absRoot, junit, lcov, mutation, s, inicio); err != nil {
			return fmt.Errorf("layer %q: ingest report: %w", s.Layer, err)
		}
		if errRun != nil {
			return fmt.Errorf("layer %q failed: %w", s.Layer, errRun)
		}

		if junit == "" && lcov == "" && mutation == "" {
			// Dizer isto em voz alta importa: sem relatório declarado o comando vira um
			// atalho de shell, e o gate continua acusando "sem sinal ingerido" — o
			// usuário precisa saber que rodou mas nada foi amarrado ao mapa.
			fmt.Printf("  [%s] passed, but the suite declares no report — nothing was ingested.\n"+
				"  Declare `junit:`/`lcov:` (or `report:`, on mutation) for the signal to reach the map.\n\n", s.Layer)
		}
		fmt.Println()
	}
	return nil
}

// ingestIfRecent só ingere o relatório que ESTA rodada produziu. A checagem de mtime
// não é zelo: quando o comando morre antes de escrever (erro de config, dependência
// faltando), o relatório da rodada ANTERIOR continua no disco — ingeri-lo gravaria no
// mapa um número velho como se fosse o de agora, que é pior que não ter número nenhum.
func ingestIfRecent(absRoot, junit, lcov, mutation string, s config.Suite, inicio time.Time) error {
	// A ingestão vem do `anchors test`: a suíte ACABOU de rodar, e o sinal corresponde a
	// ela. É o que distingue esta chamada de um `ingest` à mão.
	mapcmd.ViaAnchorsTest = true
	defer func() { mapcmd.ViaAnchorsTest = false }()
	recente := func(p string) string {
		if p == "" {
			return ""
		}
		fi, err := os.Stat(p)
		if err != nil {
			fmt.Printf("  [%s] the declared report did not show up: %s\n", s.Layer, p)
			return ""
		}
		if fi.ModTime().Before(inicio) {
			fmt.Printf("  [%s] report OLDER than this run (%s) — not ingested, so as not to record an earlier number as if it were from now.\n",
				s.Layer, p)
			return ""
		}
		return p
	}
	j, l, m := recente(junit), recente(lcov), recente(mutation)
	if j == "" && l == "" && m == "" {
		return nil
	}
	return mapcmd.IngestArtifacts(absRoot, "", j, l, m, s.Layer, s.Scope, "")
}

// workspaceLabel prefixa o workspace no cabeçalho quando existe. Sem ele, duas suítes
// `unit` de workspaces diferentes imprimem o mesmo título e a saída fica ilegível.
func workspaceLabel(s config.Suite) string {
	if s.Workspace == "" {
		return ""
	}
	return s.Workspace + "/"
}

// scopeLabel sufixa o escopo no cabeçalho. Sem ele, as duas rodadas da mesma unidade
// (isolada e completa) imprimem títulos idênticos e a saída fica ilegível — justamente
// quando as duas rodam em sequência, que é o default.
func scopeLabel(s config.Suite) string {
	if s.Scope == "" {
		return ""
	}
	return " " + s.Scope
}

// pickCommand decide entre os dois modos que o Anchors já tem em toda parte:
// COMPLETO (o projeto inteiro) e INCREMENTAL (só o caminho de impacto do que mudou).
// A simetria com o `check --all` / `check --changed` é o ponto: sem ela, o ciclo de
// quem alterou um arquivo teria um passo barato para os gates e um caro para os testes.
func pickCommand(s config.Suite, alvos []string, target string) (string, error) {
	if len(alvos) == 0 {
		return buildCommand(s.Run, target)
	}
	if strings.TrimSpace(s.RunChanged) == "" {
		// Cair para a rodada completa aqui seria o pior dos mundos: caro, e mentindo
		// sobre o que rodou — o usuário leria "passou" achando que foi o recorte dele.
		return "", fmt.Errorf("the suite declares no `run_changed:` — without it there is no incremental mode.\n"+
			"  Declare the command that receives the files, with {{files}}:\n"+
			"    run_changed: \"%s --findRelatedTests {{files}}\"   (jest example)", firstWord(s.Run))
	}
	linha := strings.ReplaceAll(s.RunChanged, "{{files}}", strings.Join(alvos, " "))
	return buildCommand(linha, target)
}

func firstWord(s string) string {
	if i := strings.IndexByte(strings.TrimSpace(s), ' '); i > 0 {
		return strings.TrimSpace(s)[:i]
	}
	return strings.TrimSpace(s)
}

// impactFiles traduz os nós do caminho de impacto nos ARQUIVOS que fazem sentido
// para cada comando. A filtragem por kind não é opinião sobre a stack: passar um
// `.spec.md` para um runner de teste ou para um mutador não significa nada, e o
// caminho de impacto do Anchors carrega a trinca inteira (spec, feature, teste, código).
//
//   - test:     código e teste — é o que `--findRelatedTests` e equivalentes esperam.
//   - mutation: só código — mutação altera a REGRA; mutar o teste inverteria o
//     experimento (o teste é o instrumento de medida, não o objeto medido).
//
// Os caminhos saem ABSOLUTOS, e isso e decisao, nao descuido. O id do no e relativo a
// raiz, mas o comando declarado quase sempre entra num workspace antes de rodar
// ("cd packages/backend && jest ..."), e ai um caminho relativo a raiz aponta para o
// lugar errado a partir do diretorio do comando. MEDIDO: jest --findRelatedTests com o
// caminho relativo a raiz, rodando DENTRO de packages/backend, achou o teste assim
// mesmo — por leniencia do jest. Outro runner falharia calado, ou rodaria a suite
// inteira achando que nao recebeu recorte. O Anchors nao tem como saber para onde o
// comando vai fazer cd, entao entrega o caminho que vale de qualquer lugar.
func impactFiles(nodes []mapx.Node, secao, absRoot string) []string {
	var out []string
	// ToSlash no fim, e não é cosmético: o comando roda dentro de `sh -c`, onde a barra
	// invertida é ESCAPE. Medido: com o separador nativo, "C:\Users\...\dedup.ts" chega
	// ao jest desfeito e a rodada acha 0 teste — sem erro, só um recorte vazio que se
	// lê como "não havia o que rodar". Com "/", o sh repassa intacto e as ferramentas
	// do Windows aceitam. Mesma regra do resto do CLI: "/" para nomear, separador
	// nativo só na API de arquivo.
	juntar := func(id string) string {
		if absRoot == "" {
			return id
		}
		return filepath.ToSlash(filepath.Join(absRoot, filepath.FromSlash(id)))
	}
	for _, n := range nodes {
		switch n.Kind {
		case mapx.KindCode:
			out = append(out, juntar(n.ID))
		case mapx.KindTest:
			if secao != "mutation" {
				out = append(out, juntar(n.ID))
			}
		}
	}
	return out
}

// suiteArgvLimit e o mesmo teto que os gates externos usam, pela mesma razao: no
// Windows o CreateProcess corta em 32767 caracteres, e o sh do MSYS2 ainda converte
// relativos em absolutos antes de chamar o .exe, multiplicando o que medimos aqui. Um
// caminho de impacto grande estoura isso com facilidade — e o modo como estoura e o
// pior possivel: o comando falha sem escrever nada, e a falha se parece com "o teste
// reprovou". Melhor dizer o que aconteceu e o que fazer.
func suiteArgvLimit() int {
	if v := os.Getenv("ANCHORS_ARGV_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if runtime.GOOS == "windows" {
		return 6000
	}
	return 100000
}

// buildCommand substitui `{{target}}`. O erro quando falta o alvo é deliberado: rodar
// com o placeholder vazio faria o Stryker mutar o projeto INTEIRO em vez do arquivo
// pedido — caro e nada do que se pediu.
func buildCommand(run, target string) (string, error) {
	if !strings.Contains(run, "{{target}}") {
		return run, nil
	}
	if strings.TrimSpace(target) == "" {
		return "", fmt.Errorf("the declared `run:` uses {{target}} — provide --target <target>")
	}
	return strings.ReplaceAll(run, "{{target}}", target), nil
}

// execAtRoot roda via `sh -c`, como os gates externos: o comando é do projeto e pode ter
// pipe, `&&`, variável — interpretá-lo aqui seria reimplementar um shell pela metade.
// A saída vai direto para o terminal, sem captura: quem roda teste quer ver o teste
// rodando, e engolir a saída para reimprimir no fim quebra qualquer barra de progresso.
func execAtRoot(linha, absRoot string) error {
	cmd := exec.Command("sh", "-c", linha) //nolint:gosec // o comando é declarado pelo projeto, como no `run:` dos gates
	cmd.Dir = absRoot
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

// absPath resolve o relatório contra a raiz. Vazio continua vazio — é o sinal de
// "esta suíte não declara este artefato".
func absPath(absRoot, p string) string {
	if strings.TrimSpace(p) == "" {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(absRoot, filepath.FromSlash(p))
}

// runChained roda os comandos do Anchors pedidos em `--then`, no processo atual
// (não re-invoca o binário: o mapa acabou de ser salvo, e um subprocesso só pagaria
// carregamento de novo).
func runChained(then, absRoot string) error {
	for _, nome := range strings.Split(then, ",") {
		nome = strings.ToLower(strings.TrimSpace(nome))
		if nome == "" {
			continue
		}
		var sub *cobra.Command
		switch nome {
		case "check":
			sub = newCheckCmd()
		case "coverage":
			sub = newCoverageCmd()
		default:
			return fmt.Errorf("--then %q is not chainable; use `check` or `coverage`", nome)
		}
		fmt.Printf("━━━ then: anchors %s ━━━\n", nome)
		sub.SetArgs([]string{"--root", absRoot})
		if err := sub.Execute(); err != nil {
			return err
		}
	}
	return nil
}

// printHowToConfigure é a resposta a "não configurado": mostrar o que declarar, não
// reclamar que falta. Quem chega aqui não sabe que a seção existe — mandá-lo para a
// documentação seria transferir o trabalho de descobrir.
func printHowToConfigure(cs suiteCommand) {
	fmt.Printf(`
No %s suite declared — Anchors does not guess how this project runs tests,
because the stack is yours. Declare it in %s:

%s

Each entry has three parts, and all three matter:
  layer:  the layer name, YOUR vocabulary — it is what `+"`anchors %s unit e2e`"+` filters by
  run:    the command, run via sh at the project root
  %s

Once that is done, `+"`anchors %s`"+` runs and INGESTS in one pass — the step that usually
gets left behind and makes the gate keep reporting "no signal ingested" even after the
suite has passed.
`, cs.secao, config.DefaultFile, cs.exemplo, cs.nome, reportLine(cs.secao), cs.nome)
}

func reportLine(secao string) string {
	if secao == "mutation" {
		return "report: the mutation JSON the run leaves behind (Mutation Testing Elements)"
	}
	return "junit:/lcov: the reports the run leaves behind, for Anchors to bind to the map"
}

// joinOrDash imprime uma lista de filtros, ou "—" quando o eixo não foi filtrado. O
// traço é melhor que a string vazia numa mensagem de erro: "camada(s) — com
// workspace(s) web" se lê como "qualquer camada do workspace web", que é o pedido real.
func joinOrDash(vs []string) string {
	if len(vs) == 0 {
		return "—"
	}
	return strings.Join(vs, ", ")
}

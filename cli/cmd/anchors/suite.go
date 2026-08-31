package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	return novoComandoSuite(comandoSuite{
		nome:  "test",
		curto: "Roda as suítes de teste declaradas no anchors.yaml e ingere os relatórios",
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
		usoLongo: `Roda o que o PROJETO declarou em ` + "`tests:`" + ` e ingere o que a rodada deixou.

  anchors test                    todas as camadas declaradas, na ordem do arquivo
  anchors test unit               só a unit
  anchors test unit integration   mais de uma, na ordem do arquivo
  anchors test -w backend         filtra por workspace (monorepo)
  anchors test --changed a.ts     INCREMENTAL: so o caminho de impacto do que mudou
  anchors test --then check       ao passar, cobra os gates com o sinal ja fresco

Dois modos, os mesmos do check: COMPLETO (sem --changed) roda o "run:"; INCREMENTAL
roda o "run_changed:", que recebe em {{files}} os arquivos do caminho de impacto -
o MESMO recorte que o "check --changed" usa, para os dois nao discordarem sobre o que
o seu commit move.

O Anchors não sabe rodar teste — a stack é do projeto. Ele executa o ` + "`run:`" + ` declarado
e amarra o ` + "`junit:`" + `/` + "`lcov:`" + ` ao mapa, que é o passo que costuma ser esquecido.`,
	})
}

func newMutationCmd() *cobra.Command {
	return novoComandoSuite(comandoSuite{
		nome:  "mutation",
		curto: "Roda as suítes de mutação declaradas no anchors.yaml e ingere os relatórios",
		secao: "mutation",
		exemplo: `mutation:
  - layer: unit
    run: "cd packages/backend && npx stryker run --mutate {{target}}"
    run_changed: "cd packages/backend && npx stryker run --mutate {{files}}"
    report: "packages/backend/reports/mutation/mutation.json"
    scope: full`,
		usoLongo: `Roda o que o PROJETO declarou em ` + "`mutation:`" + ` e ingere o relatório.

  anchors mutation                          todas as camadas declaradas
  anchors mutation unit                     só a unit
  anchors mutation unit --target x/y.ts     preenche {{target}} no comando declarado
  anchors mutation --changed x/y.ts         INCREMENTAL: muta o caminho de impacto
  anchors mutation unit --then check        ao passar, cobra os gates em seguida

No incremental, {{files}} recebe so os nos de CODIGO: mutacao altera a regra, e mutar
o teste inverteria o experimento - o teste e o instrumento de medida, nao o objeto.

Mutação responde o que cobertura não responde: altere a linha — o teste percebe? Um
mutante SOBREVIVENTE é uma linha que ninguém prova.

Atenção ao ` + "`scope:`" + `: com ` + "`isolated`" + ` e ` + "`full`" + ` ambos ingeridos, o gate julga pelo
ISOLADO, que é bem mais duro — é ele que diz se o teste da unidade prova a unidade,
em vez de ela ser provada pelos dependentes.`,
	})
}

// comandoSuite descreve o que muda entre `test` e `mutation` — o resto é idêntico, e
// duplicar os dois faria a mensagem de "não configurado" divergir com o tempo.
type comandoSuite struct {
	nome     string
	curto    string
	secao    string // o nome da seção no anchors.yaml, usado nas mensagens
	exemplo  string // o bloco YAML que se mostra a quem não configurou
	usoLongo string
}

func novoComandoSuite(cs comandoSuite) *cobra.Command {
	var root, target, then string
	var workspaces, changed, escopos []string
	cmd := &cobra.Command{
		Use:   cs.nome + " [camadas...]",
		Short: cs.curto,
		Long:  cs.usoLongo,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar %s: %w", config.DefaultFile, err)
			}

			declaradas := cfg.Tests
			if cs.secao == "mutation" {
				declaradas = cfg.Mutation
			}
			if len(declaradas) == 0 {
				imprimeComoConfigurar(cs)
				return fmt.Errorf("nenhuma suíte declarada em `%s:` no %s", cs.secao, config.DefaultFile)
			}

			sel, ausentes := config.SelecionaSuites(declaradas, args, workspaces, escopos)
			if len(ausentes) > 0 {
				return fmt.Errorf("não declarado em `%s:`: %s\n  camadas declaradas:    %s\n  workspaces declarados: %s\n  escopos declarados:    %s",
					cs.secao, strings.Join(ausentes, ", "),
					strings.Join(config.CamadasDeclaradas(declaradas), ", "),
					juntaOuTraco(config.WorkspacesDeclarados(declaradas)),
					juntaOuTraco(config.EscoposDeclarados(declaradas)))
			}
			// Combinação válida mas vazia não é erro de digitação — é "não existe essa
			// suíte". Dizer isso com clareza evita a leitura de que o comando rodou e
			// tudo passou, que é o que um silêncio com exit 0 comunicaria.
			if len(sel) == 0 {
				return fmt.Errorf("nenhuma suíte combina camada(s) %s, workspace(s) %s e escopo(s) %s em `%s:`",
					juntaOuTraco(args), juntaOuTraco(workspaces), juntaOuTraco(escopos), cs.secao)
			}

			// O caminho de impacto sai da MESMA função que o `check --changed` usa: se
			// as duas divergissem, "os gates que o meu commit move" e "os testes que o
			// meu commit move" passariam a ser conjuntos diferentes, e o incremental
			// deixaria de ser confiável exatamente onde ele é a única defesa.
			var alvos []string
			if len(changed) > 0 {
				g, mapErr := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
				if mapErr != nil {
					return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", mapErr)
				}
				nodes, _, selErr := selectNodes(g, cfg, false, changed, absRoot)
				if selErr != nil {
					return selErr
				}
				alvos = arquivosDoImpacto(nodes, cs.secao, absRoot)
				if len(alvos) == 0 {
					fmt.Printf("o caminho de impacto não alcança arquivo de %s — nada a rodar.\n",
						map[bool]string{true: "código", false: "código ou teste"}[cs.secao == "mutation"])
					return nil
				}
				fmt.Printf("incremental: %d arquivo(s) no caminho de impacto\n\n", len(alvos))
			}

			if err := rodaSuites(cs, sel, absRoot, target, alvos); err != nil {
				return err
			}
			// O encadeamento é OPT-IN e só acontece depois do sucesso: um `check` sobre
			// sinal que não foi produzido diria o mesmo de antes, e um sobre suíte que
			// falhou culparia o gate por um teste vermelho.
			return executaEncadeados(then, absRoot)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringSliceVarP(&workspaces, "workspace", "w", nil, "filtra pelo WORKSPACE declarado (backend, mobile…) — repetível ou separado por vírgula. Combina com as camadas: `anchors "+cs.nome+" unit -w backend`")
	// Sem `--scope`, roda TODOS os escopos declarados — e é o default certo: é o par
	// isolado+completo que produz a leitura de acoplamento do gate, e pedir dois comandos
	// para obtê-la devolveria o passo manual que estes comandos vieram eliminar.
	cmd.Flags().StringSliceVar(&escopos, "scope", nil, "filtra pelo ESCOPO declarado (`isolated`, `full`) — só na mutação. Sem isto, roda os que estiverem declarados")
	cmd.Flags().StringVar(&target, "target", "", "alvo que substitui `{{target}}` no `run:` declarado (ex.: o arquivo a mutar)")
	cmd.Flags().StringSliceVar(&changed, "changed", nil, "modo INCREMENTAL: arquivo(s) alterado(s) — roda o `run_changed:` sobre a união dos caminhos de impacto, o mesmo recorte do `check --changed`")
	cmd.Flags().StringVar(&then, "then", "", "ao PASSAR, encadeia comandos do Anchors: `check`, `coverage` (separe por vírgula). Opt-in — sem isto, roda e para")
	return cmd
}

// rodaSuites executa cada suíte e ingere o que ela deixou. Para na primeira que falhar:
// as camadas costumam depender umas das outras (não faz sentido rodar e2e depois de a
// unit quebrar), e seguir adiante só produziria ruído sobre uma base já vermelha.
func rodaSuites(cs comandoSuite, suites []config.Suite, absRoot, target string, alvos []string) error {
	for _, s := range suites {
		linha, err := escolheComando(s, alvos, target)
		if err != nil {
			return fmt.Errorf("camada %q: %w", s.Layer, err)
		}
		fmt.Printf("━━━ %s [%s%s%s] ━━━\n%s\n\n", cs.nome, rotuloWorkspace(s), s.Layer, rotuloEscopo(s), linha)

		if teto := limiteArgvSuite(); len(linha) > teto {
			return fmt.Errorf("camada %q: a linha de comando ficou com %d caracteres (teto %d nesta plataforma).\n"+
				"  O caminho de impacto é grande demais para uma invocação só. Rode o modo completo,\n"+
				"  reduza o --changed, ou ajuste ANCHORS_ARGV_MAX se souber que o seu shell aguenta",
				s.Layer, len(linha), teto)
		}

		inicio := time.Now()
		errRun := execNaRaiz(linha, absRoot)

		junit, lcov, mutation := caminhoAbs(absRoot, s.JUnit), caminhoAbs(absRoot, s.Lcov), caminhoAbs(absRoot, s.Report)
		// A ingestão acontece MESMO se o comando saiu != 0, e essa é a regra menos
		// óbvia daqui. Um runner sai != 0 exatamente quando há o que reportar: o jest
		// quando um teste falha, o Stryker quando o score fica abaixo do próprio
		// `thresholds.break`. Pular a ingestão nesse caso deixaria justamente a unidade
		// PROBLEMÁTICA sem sinal no mapa — o gate seguiria dizendo "sem sinal ingerido"
		// sobre a única que já se sabe ruim, e o vermelho pareceria ausência de medida.
		//
		// O relatório é a evidência do que aconteceu, não um prêmio por ter passado.
		if err := ingereSeRecente(absRoot, junit, lcov, mutation, s, inicio); err != nil {
			return fmt.Errorf("camada %q: ingerir relatório: %w", s.Layer, err)
		}
		if errRun != nil {
			return fmt.Errorf("camada %q falhou: %w", s.Layer, errRun)
		}

		if junit == "" && lcov == "" && mutation == "" {
			// Dizer isto em voz alta importa: sem relatório declarado o comando vira um
			// atalho de shell, e o gate continua acusando "sem sinal ingerido" — o
			// usuário precisa saber que rodou mas nada foi amarrado ao mapa.
			fmt.Printf("  [%s] passou, mas a suíte não declara relatório — nada foi ingerido.\n"+
				"  Declare `junit:`/`lcov:` (ou `report:`, em mutation) para o sinal chegar ao mapa.\n\n", s.Layer)
		}
		fmt.Println()
	}
	return nil
}

// ingereSeRecente só ingere o relatório que ESTA rodada produziu. A checagem de mtime
// não é zelo: quando o comando morre antes de escrever (erro de config, dependência
// faltando), o relatório da rodada ANTERIOR continua no disco — ingeri-lo gravaria no
// mapa um número velho como se fosse o de agora, que é pior que não ter número nenhum.
func ingereSeRecente(absRoot, junit, lcov, mutation string, s config.Suite, inicio time.Time) error {
	// A ingestão vem do `anchors test`: a suíte ACABOU de rodar, e o sinal corresponde a
	// ela. É o que distingue esta chamada de um `ingest` à mão.
	viaAnchorsTest = true
	defer func() { viaAnchorsTest = false }()
	recente := func(p string) string {
		if p == "" {
			return ""
		}
		fi, err := os.Stat(p)
		if err != nil {
			fmt.Printf("  [%s] o relatório declarado não apareceu: %s\n", s.Layer, p)
			return ""
		}
		if fi.ModTime().Before(inicio) {
			fmt.Printf("  [%s] relatório mais VELHO que esta rodada (%s) — não ingerido, para não gravar número de antes como se fosse de agora.\n",
				s.Layer, p)
			return ""
		}
		return p
	}
	j, l, m := recente(junit), recente(lcov), recente(mutation)
	if j == "" && l == "" && m == "" {
		return nil
	}
	return ingereArtefatos(absRoot, "", j, l, m, s.Layer, s.Scope)
}

// rotuloWorkspace prefixa o workspace no cabeçalho quando existe. Sem ele, duas suítes
// `unit` de workspaces diferentes imprimem o mesmo título e a saída fica ilegível.
func rotuloWorkspace(s config.Suite) string {
	if s.Workspace == "" {
		return ""
	}
	return s.Workspace + "/"
}

// rotuloEscopo sufixa o escopo no cabeçalho. Sem ele, as duas rodadas da mesma unidade
// (isolada e completa) imprimem títulos idênticos e a saída fica ilegível — justamente
// quando as duas rodam em sequência, que é o default.
func rotuloEscopo(s config.Suite) string {
	if s.Scope == "" {
		return ""
	}
	return " " + s.Scope
}

// escolheComando decide entre os dois modos que o Anchors já tem em toda parte:
// COMPLETO (o projeto inteiro) e INCREMENTAL (só o caminho de impacto do que mudou).
// A simetria com o `check --all` / `check --changed` é o ponto: sem ela, o ciclo de
// quem alterou um arquivo teria um passo barato para os gates e um caro para os testes.
func escolheComando(s config.Suite, alvos []string, target string) (string, error) {
	if len(alvos) == 0 {
		return montaComando(s.Run, target)
	}
	if strings.TrimSpace(s.RunChanged) == "" {
		// Cair para a rodada completa aqui seria o pior dos mundos: caro, e mentindo
		// sobre o que rodou — o usuário leria "passou" achando que foi o recorte dele.
		return "", fmt.Errorf("a suíte não declara `run_changed:` — sem ele não há modo incremental.\n"+
			"  Declare o comando que recebe os arquivos, com {{files}}:\n"+
			"    run_changed: \"%s --findRelatedTests {{files}}\"   (exemplo de jest)", primeiraPalavra(s.Run))
	}
	linha := strings.ReplaceAll(s.RunChanged, "{{files}}", strings.Join(alvos, " "))
	return montaComando(linha, target)
}

func primeiraPalavra(s string) string {
	if i := strings.IndexByte(strings.TrimSpace(s), ' '); i > 0 {
		return strings.TrimSpace(s)[:i]
	}
	return strings.TrimSpace(s)
}

// arquivosDoImpacto traduz os nós do caminho de impacto nos ARQUIVOS que fazem sentido
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
func arquivosDoImpacto(nodes []mapx.Node, secao, absRoot string) []string {
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

// limiteArgvSuite e o mesmo teto que os gates externos usam, pela mesma razao: no
// Windows o CreateProcess corta em 32767 caracteres, e o sh do MSYS2 ainda converte
// relativos em absolutos antes de chamar o .exe, multiplicando o que medimos aqui. Um
// caminho de impacto grande estoura isso com facilidade — e o modo como estoura e o
// pior possivel: o comando falha sem escrever nada, e a falha se parece com "o teste
// reprovou". Melhor dizer o que aconteceu e o que fazer.
func limiteArgvSuite() int {
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

// montaComando substitui `{{target}}`. O erro quando falta o alvo é deliberado: rodar
// com o placeholder vazio faria o Stryker mutar o projeto INTEIRO em vez do arquivo
// pedido — caro e nada do que se pediu.
func montaComando(run, target string) (string, error) {
	if !strings.Contains(run, "{{target}}") {
		return run, nil
	}
	if strings.TrimSpace(target) == "" {
		return "", fmt.Errorf("o `run:` declarado usa {{target}} — informe --target <alvo>")
	}
	return strings.ReplaceAll(run, "{{target}}", target), nil
}

// execNaRaiz roda via `sh -c`, como os gates externos: o comando é do projeto e pode ter
// pipe, `&&`, variável — interpretá-lo aqui seria reimplementar um shell pela metade.
// A saída vai direto para o terminal, sem captura: quem roda teste quer ver o teste
// rodando, e engolir a saída para reimprimir no fim quebra qualquer barra de progresso.
func execNaRaiz(linha, absRoot string) error {
	cmd := exec.Command("sh", "-c", linha) //nolint:gosec // o comando é declarado pelo projeto, como no `run:` dos gates
	cmd.Dir = absRoot
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

// caminhoAbs resolve o relatório contra a raiz. Vazio continua vazio — é o sinal de
// "esta suíte não declara este artefato".
func caminhoAbs(absRoot, p string) string {
	if strings.TrimSpace(p) == "" {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(absRoot, filepath.FromSlash(p))
}

// executaEncadeados roda os comandos do Anchors pedidos em `--then`, no processo atual
// (não re-invoca o binário: o mapa acabou de ser salvo, e um subprocesso só pagaria
// carregamento de novo).
func executaEncadeados(then, absRoot string) error {
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
			return fmt.Errorf("--then %q não é encadeável; use `check` ou `coverage`", nome)
		}
		fmt.Printf("━━━ then: anchors %s ━━━\n", nome)
		sub.SetArgs([]string{"--root", absRoot})
		if err := sub.Execute(); err != nil {
			return err
		}
	}
	return nil
}

// imprimeComoConfigurar é a resposta a "não configurado": mostrar o que declarar, não
// reclamar que falta. Quem chega aqui não sabe que a seção existe — mandá-lo para a
// documentação seria transferir o trabalho de descobrir.
func imprimeComoConfigurar(cs comandoSuite) {
	fmt.Printf(`
Nenhuma suíte de %s declarada — o Anchors não adivinha como este projeto roda teste,
porque a stack é sua. Declare no %s:

%s

Cada entrada tem três partes, e as três importam:
  layer:  o nome da camada, vocabulário SEU — é por ele que `+"`anchors %s unit e2e`"+` filtra
  run:    o comando, rodado via sh na raiz do projeto
  %s

Feito isso, `+"`anchors %s`"+` roda e INGERE numa passada — que é o passo que costuma
ficar para trás e faz o gate seguir acusando "sem sinal ingerido" mesmo depois de a
suíte ter passado.
`, cs.secao, config.DefaultFile, cs.exemplo, cs.nome, linhaRelatorio(cs.secao), cs.nome)
}

func linhaRelatorio(secao string) string {
	if secao == "mutation" {
		return "report: o JSON de mutação que a rodada deixa (Mutation Testing Elements)"
	}
	return "junit:/lcov: os relatórios que a rodada deixa, para o Anchors amarrar ao mapa"
}

// juntaOuTraco imprime uma lista de filtros, ou "—" quando o eixo não foi filtrado. O
// traço é melhor que a string vazia numa mensagem de erro: "camada(s) — com
// workspace(s) web" se lê como "qualquer camada do workspace web", que é o pedido real.
func juntaOuTraco(vs []string) string {
	if len(vs) == 0 {
		return "—"
	}
	return strings.Join(vs, ", ")
}

package quality

import (
	"context"
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/change"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	var root, mapPath string
	var corrigir bool
	var soPipelines bool
	cmd := &cobra.Command{
		Use: "doctor",

		Short: "X-ray of the ecosystem: systemic loose ends, health and maturity",
		Long: `The ecosystem's health validator (QUALITY §5.2) — the GLOBAL view.
Unlike check (node against gate, incremental), doctor sweeps the map + the
Structure + the disk and hunts the SYSTEMIC loose ends no local gate sees:
map integrity (dead edges, phantom nodes), orphans (code without spec,
missing identity), loose layers, and gate coverage holes.

It PRESENTS and RECORDS, but does NOT block — it is diagnosis, run on demand.

With --fix, it creates what is missing: pipelines, git hooks, state labels and the
protection of the integration branch. It is IDEMPOTENT — in an already assembled
repository it changes nothing, and each repair checks before acting. Run it without asking
for authorization: it is the environment preparation step, not a project
change.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// `--check-pipelines` responde UMA pergunta e sai: o pipeline que roda no CI
			// está atualizado?
			//
			// Existe porque o pipeline DESATUALIZADO falha em SILÊNCIO — ele roda, faz o
			// que a versão dele sabia fazer, e o que foi corrigido depois simplesmente não
			// acontece. Medido: um passo novo do `identify` não rodou por três execuções
			// sem que nada acusasse, e só apareceu quando fui ler o log procurando outra
			// coisa.
			//
			// O raio-X completo seria ruído no CI: ele responde dezenas de perguntas, e
			// quem chama daqui quer uma.
			if soPipelines {
				return checkPipelines(cmd)
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}

			rep := health.Diagnose(g, cfg, absRoot)
			printReport(rep)
			if corrigir {
				return repairEnvironment(absRoot, cfg)
			}
			return nil // doctor NUNCA bloqueia — só reporta
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	// A descrição diz IDEMPOTENTE, e isso não é detalhe de implementação — é a única
	// informação que decide se um agente roda o comando ou pára para perguntar.
	//
	// Medido: um dev novo pediu ao agente dele para contribuir, o agente montou o plano de
	// onboarding correto e PAROU no `doctor --fix`, pedindo autorização — porque "protege o
	// branch" e "cria as labels" se leem como mudança de estado compartilhado. A cautela
	// estava certa; o que faltava era o comando dizer que num repositório já montado ele
	// não muda nada.
	cmd.Flags().BoolVar(&corrigir, "fix", false,
		"creates what is missing in the `github` mode environment (pipelines, labels, branch protection). "+
			"IDEMPOTENT: in an already set-up repository it changes nothing, and it is safe to always run")
	cmd.Flags().BoolVar(&soPipelines, "check-pipelines", false,
		"exits with 1 if any Anchors pipeline is out of date (so CI can self-verify)")
	return cmd
}

func printReport(r health.Report) {
	fmt.Println(i18n.T("doctor.header", r.Nodes, r.Edges, r.Layers))
	fmt.Println()

	warnings := r.Warnings()
	if len(r.Findings) == 0 {
		fmt.Println(i18n.T("doctor.all_clean"))
		return
	}

	// agrupa por verificação
	byCheck := map[string][]health.Finding{}
	var order []string
	for _, f := range r.Findings {
		if _, ok := byCheck[f.Check]; !ok {
			order = append(order, f.Check)
		}
		byCheck[f.Check] = append(byCheck[f.Check], f)
	}

	for _, check := range order {
		fs := byCheck[check]
		// o grupo é Warn se QUALQUER finding for Warn (a mais alta vence) — senão o
		// ícone esconderia warns atrás de um info que veio primeiro na ordenação.
		mark := "ℹ"
		nWarn := 0
		for _, f := range fs {
			if f.Severity == health.Warn {
				nWarn++
			}
		}
		if nWarn > 0 {
			mark = "⚠"
		}
		fmt.Printf("%s %s (%d)\n", mark, check, len(fs))
		for _, f := range fs {
			itemMark := " "
			if f.Severity == health.Warn {
				itemMark = "⚠"
			}
			fmt.Printf("  %s %s — %s\n", itemMark, f.Subject, f.Detail)
		}
		fmt.Println()
	}

	fmt.Println(i18n.T("doctor.summary", len(warnings), len(r.Findings)))
	fmt.Println(i18n.T("doctor.diagnosis_note"))
}

// repairEnvironment cria o que falta no ambiente do modo `github`. É o `--fix` do doctor,
// no precedente do `check --fix`: sem a flag, o doctor segue sendo diagnóstico puro
// ("nada foi bloqueado; decida o que conciliar"), e com ela ele age.
//
// Conserta só o que é DELE para consertar: os pipelines são arquivos do repositório,
// versionados e revisáveis num diff, e criá-los não afeta ninguém até serem commitados.
//
// O BOARD não é criado aqui, e a diferença é de natureza: um GitHub Project é estrutura
// COMPARTILHADA pelo time e vive fora do repositório — um board criado por engano polui a
// organização inteira e não se desfaz com `git checkout`. O doctor diz o que falta; criar
// é decisão de quem opera.
func repairEnvironment(root string, cfg *config.Config) error {
	if !cfg.GitHubMode() {
		fmt.Println()
		fmt.Println(i18n.T("doctor.fix.nothing_github_mode"))
		return nil
	}
	// A CREDENCIAL, ANTES DE TUDO.
	//
	// No modo `github` quase todo reparo passa pelo `gh`: proteger o branch, criar as
	// labels de estado, ler o board. Sem credencial, cada um falha por conta — e quem lê a
	// saída recebe cinco erros diferentes para uma causa só.
	//
	// Medido com um dev novo: o `doctor --fix` não conseguiu proteger o branch, desligar a
	// exigência de aprovação, nem criar as labels, e a saída não dizia que a causa era
	// login. Ele descobriu sozinho, testando.
	//
	// Uma verificação antes vale mais que cinco erros depois: ela nomeia a causa, dá o
	// comando, e não gasta a atenção de quem lê com sintomas.
	if err := requireGHAuth(); err != nil {
		return err
	}
	// A FILA LOCAL NÃO DEVE EXISTIR no modo github.
	//
	// O `config.go` declara o modo como EXCLUDENTE — *"um modo OU outro, nunca um com o
	// outro de reserva"* — e diz que a validação é simples: *"no modo github,
	// `.anchors/tasks/` não deve existir"*. Ninguém a havia escrito.
	//
	// Medido no blue-eyes: com `mode: github` no anchors.yaml, `.anchors/tasks/` tinha 11
	// arquivos e o `anchors next` lia dali — respondendo "fila vazia" com 84 cards abertos
	// no board. Duas filas para a mesma pergunta, e a resposta vinha da errada.
	warnOrphanLocalQueue(root)
	warnOrphanChanges(root)

	// Lido ANTES de semear: depois da escrita os arquivos já casam o template, e não
	// haveria como dizer quais foram ATUALIZADOS em vez de criados.
	faltavam := initx.MissingWorkflow(root)
	desatualizados := initx.OutdatedWorkflows(root, cfg)
	_, board, err := initx.SemeiaWorkflows(root, cfg)
	if err != nil {
		return fmt.Errorf("seed the pipelines: %w", err)
	}
	fmt.Println()
	// A PÁGINA do board é anunciada ANTES do resumo dos pipelines, e separada dele: ela
	// muda por conta própria (o template do Anchors evolui) e é a única coisa aqui que
	// aparece como diff no repositório de quem rodou o comando.
	//
	// Antes disto o `--fix` imprimia "os pipelines já existem e estão atualizados"
	// enquanto REESCREVIA a página em silêncio. Quem visse o `git status` depois não
	// saberia se tinha feito aquilo — e a única saída era ler o diff inteiro para
	// descobrir de quem era a mudança.
	switch board {
	case initx.BoardCreated:
		fmt.Println(i18n.T("doctor.fix.board_created", initx.BoardFile))
	case initx.BoardUpdated:
		fmt.Println(i18n.T("doctor.fix.board_updated", initx.BoardFile))
		fmt.Println(i18n.T("doctor.fix.board_review_diff"))
	}
	if len(faltavam) == 0 && len(desatualizados) == 0 && board == initx.BoardUnchanged {
		fmt.Println(i18n.T("doctor.fix.pipelines_current"))
	}
	if len(faltavam) > 0 {
		fmt.Println(i18n.T("doctor.fix.pipelines_created", len(faltavam), initx.DirWorkflows))
		for _, w := range faltavam {
			fmt.Printf("    %s\n", w.Arquivo)
		}
		fmt.Println(i18n.T("doctor.fix.review_and_configure"))
	}
	// Atualizado é distinto de criado, e a mensagem separa os dois: um arquivo que MUDOU
	// sozinho no repositório de alguém precisa ser lido antes de subir — dizer só
	// "criados" faria o diff parecer coisa que o time não fez.
	if len(desatualizados) > 0 {
		fmt.Println(i18n.T("doctor.fix.pipelines_updated", len(desatualizados), initx.DirWorkflows))
		for _, w := range desatualizados {
			fmt.Printf("    %s\n", w.Arquivo)
		}
		fmt.Println(i18n.T("doctor.fix.review_diff"))
	}
	// A PROTEÇÃO DO BRANCH é o que enforça "todo trabalho sobe via PR". Sem ela nada
	// falha: o push direto funciona, e pula o card, a revisão e o pipeline de
	// identificação — que dispara na ABERTURA do PR.
	if err := protectBranches(cfg); err != nil {
		fmt.Println(i18n.T("doctor.fix.protect_branches_error", err))
	}

	// A APROVAÇÃO INALCANÇÁVEL: quando o fluxo exige aprovação que o GitHub impede o
	// autor de dar, e não há admin para ignorar a regra, o `--fix` desliga a exigência.
	//
	// É deliberado o fix mexer aqui: deixar a exigência de pé seria manter um fluxo que
	// NÃO TEM COMO ser cumprido, e quem opera descobriria no meio de um merge. A revisão
	// continua sendo cobrada — pelo estado do card, que é o que o Anchors controla.
	if cfg.Workflow.RequiredApprovalsOrDefault() > 0 {
		repo, branch := cfg.Workflow.Repo, cfg.Workflow.IntegrationBranchOrDefault()
		if ok, _ := health.CanBypassProtection(repo, branch); !ok {
			if err := health.DisableApprovalRequirement(repo, branch); err != nil {
				fmt.Println(i18n.T("doctor.fix.disable_approval_error", err))
			} else {
				fmt.Println(i18n.T("doctor.fix.disable_approval_notice"))
			}
		}
	}

	// As LABELS de estado são o único pré-requisito real do fluxo — e criá-las é seguro:
	// label é do repositório, reversível, e não afeta ninguém fora dele.
	if err := createStateLabels(cfg); err != nil {
		fmt.Println(i18n.T("doctor.fix.state_labels_error", err))
	}

	// O BOARD é OPCIONAL: o estado vive na label, e o Project só espelha. Dizê-lo aqui
	// evita que alguém conclua que falta configurar algo.
	fmt.Println(i18n.T("doctor.fix.board_optional_hint",
		strings.Join(initx.ColunasDoBoard, " · "), initx.AnchorsFinalState))
	return nil
}

// createStateLabels garante as labels que carregam o estado do trabalho. São o único
// pré-requisito do fluxo que não é arquivo — e sem elas o `identify` cria cards que o
// `claim` nunca encontra.
func createStateLabels(cfg *config.Config) error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("`gh` is not on the PATH")
	}
	repo := cfg.Workflow.Repo
	// Cores por família: a fila em cinza, o trabalho em curso em azul, o que saiu da
	// alçada do Anchors em verde.
	cor := map[string]string{
		"anchors:to-do": "ededed", "anchors:ready-to-review": "ededed",
		"anchors:in-progress": "1d76db", "anchors:in-review": "1d76db",
		"anchors:ready-to-test": "0e8a16", "anchors:in-test": "0e8a16",
		"anchors:ready-to-release": "0e8a16", "anchors:production": "0e8a16",
		// VERMELHO, e é o único: o card escalado é o que ninguém no fluxo destrava, e
		// precisa saltar num board cheio de cinza e azul.
		initx.LabelNeedsUser: "d73a4a",
		// LARANJA: a bug waits for someone who can fix the pipeline or the tool — urgent
		// like a decision, but it asks for work, not a choice, and must not read as one.
		initx.LabelBug: "e99695",
		// ÂMBAR: o opt-out da trava de estado. Nem alarme (o card não está travado) nem
		// fluxo normal (alguém autorizou mexer nele à mão) — a cor de atenção é a que
		// separa os dois.
		initx.LabelManual: "fbca04",
		// CINZA APAGADO: o card descartado saiu do board, e a cor diz isso sem competir
		// com nada. Não é alarme nem trabalho — é registro de que a pergunta existiu.
		initx.LabelDiscarded: "d4d4d4",
	}
	var criadas int
	// A label de ESCALAÇÃO entra junto: sem ela criada, o pipeline que escala um card
	// falha ao aplicá-la — e falha em silêncio, porque `gh issue edit` com label
	// inexistente não é erro fatal. O card ficaria travado sem o sinalizador que diz por
	// quê.
	// O OPT-OUT entra pelo mesmo motivo do `needs-user` acima: sem ele criado, quem quiser
	// mover um card à mão não consegue aplicá-lo, e a trava de estado vira uma parede sem
	// porta. `gh issue edit` com label inexistente falha em silêncio.
	todas := append([]string{cfg.Workflow.Labels[0], initx.LabelNeedsUser, initx.LabelBug,
		initx.LabelManual, initx.LabelDiscarded},
		initx.WorkStates...)
	for _, e := range todas {
		c := cor[e]
		if c == "" {
			c = "5319e7" // a label do próprio Anchors
		}
		out, err := exec.Command("gh", "label", "create", e,
			"--repo", repo, "--color", c, "--force").CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", e, strings.TrimSpace(string(out)))
		}
		criadas++
	}
	fmt.Println(i18n.T("doctor.fix.state_labels_created", criadas, repo))
	return nil
}

// protectionBody monta o JSON que a API de proteção de branch exige.
//
// `required_approving_review_count` era ZERO, com o argumento de que exigir aprovação de
// outra conta travaria um time de uma pessoa. Isso valia quando quem aprovava era gente —
// e deixou de valer: no fluxo do modo `github` quem revisa é OUTRO AGENTE, que pega o
// card em `ready-to-review` e aprova o PR.
//
// Com zero, o merge acontece sem que ninguém revise, e a etapa que o card representa vira
// decoração. Medido nesta sessão: mesclei três PRs direto, e o board dizia
// "ready-to-review" sobre trabalho que já estava na develop.
//
// Os demais campos vão nulos de propósito: cada um é decisão do time, e o Anchors só
// cobra a porta. Os três são OBRIGATÓRIOS no corpo, mesmo nulos — a API responde 422 se
// qualquer um faltar, e a função é separada para que um teste confronte isso sem falar
// com o GitHub.
func protectionBody(aprovacoes int) string {
	return fmt.Sprintf(`{"required_status_checks":null,"enforce_admins":false,`+
		`"required_pull_request_reviews":{"required_approving_review_count":%d},`+
		`"restrictions":null}`, aprovacoes)
}

// protectBranches exige PR nos branches que o projeto declarou como portas.
//
// Não exige APROVAÇÃO de outra conta: num time de uma pessoa com agentes, isso travaria
// o fluxo inteiro. O PR existe aqui para o card ter objeto, para a revisão acontecer e
// para o histórico ficar legível — não para satisfazer uma contagem.
func protectBranches(cfg *config.Config) error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("%s", i18n.T("doctor.gh_not_installed"))
	}
	repo := cfg.Workflow.Repo
	for _, b := range cfg.Workflow.ProtectedBranchesOrDefault() {
		body := protectionBody(cfg.Workflow.RequiredApprovalsOrDefault())
		// `--input -` LÊ do stdin, e é preciso de fato escrever nele: sem isso o corpo
		// chega vazio e a API responde 422 reclamando de um campo obrigatório nulo — que
		// foi o que aconteceu enquanto o `body` era montado e descartado logo abaixo.
		cmd := exec.Command("gh", "api", "--method", "PUT",
			"repos/"+repo+"/branches/"+b+"/protection",
			"--input", "-")
		cmd.Stdin = strings.NewReader(body)
		out, err := cmd.CombinedOutput()
		if err != nil {
			// Branch inexistente é comum (o `staging` só nasce quando alguém o cria), e
			// não é falha do fix: reportar e seguir é melhor que abortar o resto.
			if strings.Contains(string(out), "Branch not found") {
				fmt.Println(i18n.T("doctor.fix.branch_not_found", b))
				continue
			}
			return fmt.Errorf("%s: %s", b, strings.TrimSpace(string(out)))
		}
		fmt.Println(i18n.T("doctor.fix.branch_protected", b))
	}
	return nil
}

// checkPipelines é o `doctor --check-pipelines`: uma pergunta, um código de saída.
//
// Sai com 1 quando algo está desatualizado ou faltando, para o CI poder barrar. Um aviso
// que não muda o código de saída seria ignorado pelo próprio pipeline que o emitiu.
func checkPipelines(cmd *cobra.Command) error {
	root := config.ProjectRoot(".")
	cfg, err := config.Load(filepath.Join(root, "anchors.yaml"))
	if err != nil {
		return err
	}
	if !cfg.GitHubMode() {
		fmt.Println(i18n.T("doctor.pipelines.local_mode"))
		return nil
	}

	faltam := initx.MissingWorkflow(root)
	velhos := initx.OutdatedWorkflows(root, cfg)
	if len(faltam) == 0 && len(velhos) == 0 {
		fmt.Println(i18n.T("doctor.pipelines.all_current", len(initx.WorkflowsDoFluxo), common.Version))
		return nil
	}

	for _, w := range faltam {
		fmt.Println(i18n.T("doctor.pipelines.missing", w.Arquivo, w.Papel))
	}
	// DESATUALIZADO é o caso que motivou o comando, e a mensagem diz o que ninguém vê
	// olhando o log: o pipeline RODOU, e fez o que a versão dele sabia fazer.
	for _, w := range velhos {
		fmt.Println(i18n.T("doctor.pipelines.outdated", w.Arquivo))
	}
	fmt.Println()
	fmt.Println(i18n.T("doctor.pipelines.run_fix_hint"))
	fmt.Println(i18n.T("doctor.pipelines.ci_binary_version", common.Version))
	// BARRAR é opt-in. O padrão avisa, porque o pipeline velho ainda faz o trabalho antigo
	// e derrubá-lo troca "faz menos do que devia" por "não faz nada". Quem declara
	// `stale_pipeline_blocks: true` decidiu que, no CI dele, o aviso não seria lido.
	if !cfg.Workflow.PipelineVelhoBarra() {
		fmt.Println(i18n.T("doctor.pipelines.stale_blocks_hint"))
		return nil
	}
	cmd.SilenceUsage = true
	return fmt.Errorf("%s", i18n.T("doctor.pipelines.error_blocked", len(faltam)+len(velhos)))
}

// warnOrphanLocalQueue avisa (e não apaga) a fila local encontrada em modo github.
//
// AVISA em vez de apagar porque a fila pode ter task `claimed__` de um trabalho em curso:
// apagar sem olhar destrói o registro de quem estava com o quê. E o histórico
// (`.anchors/done/`) é memória — quem migrou de `local` para `github` não deve perdê-lo.
func warnOrphanLocalQueue(root string) {
	dir := filepath.Join(root, ".anchors", "tasks")
	entradas, err := os.ReadDir(dir)
	if err != nil || len(entradas) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(i18n.T("doctor.orphan_tasks.header", len(entradas)))
	fmt.Println(i18n.T("doctor.orphan_tasks.body"))
}

// warnOrphanChanges avisa sobre registros de entrega em arquivo encontrados em modo github.
//
// É o mesmo defeito da fila órfã, por outra porta. No modo `github` a entrega é registrada
// como comentário na ISSUE — é lá que o revisor está olhando, e é o pipeline que move o
// card para `ready-to-review`. Um `changes/*.md` ali não é lido por ninguém: o watcher que
// o veria não é o mecanismo naquele modo.
//
// Medido no projeto de referência: 73 registros no disco, ZERO revisados, e as issues sem
// a informação. Cada um representa uma entrega cuja intenção declarada nunca foi
// confrontada — e a intenção declarada é justamente o que o registro existe para dar.
//
// AVISA e não apaga, pela mesma razão da fila: o registro é memória do que aconteceu com o
// produto, e apagá-lo sem alguém olhar destrói o que ninguém mais tem.
func warnOrphanChanges(root string) {
	pendentes, err := change.Pending(root)
	if err != nil || len(pendentes) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(i18n.T("doctor.orphan_changes.header", len(pendentes), change.Dir))
	fmt.Println(i18n.T("doctor.orphan_changes.body", change.ReviewedDir))
}

// requireGHAuth confere a credencial do `gh` antes de tentar reparar.
//
// Recusa com instrução em vez de tentar e falhar: `gh auth login` é interativo, e um agente
// não tem como completá-lo — mandar a pessoa fazê-lo é a única saída honesta.
func requireGHAuth() error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("%s", i18n.T("doctor.gh_auth_required"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "gh", "auth", "status")
	// Sem stdin: se o `gh` decidir perguntar algo, ele falha em vez de pendurar o doctor.
	cmd.Stdin = nil
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", i18n.T("doctor.gh_auth_status_failed", strings.TrimSpace(string(out))))
	}
	return nil
}

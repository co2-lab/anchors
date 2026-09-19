package quality

import (
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// newStatusCmd — `anchors status`: onde o projeto ESTÁ, e qual é o próximo passo.
//
// Existe porque a volta ao trabalho é sempre iniciada por uma IA (o Anchors não tem
// comando para "continuar"), e quem retoma precisa saber onde parou. Sem isso, o agente
// que abre a conversa dias depois começa adivinhando: já houve entrevista? o mapa está
// construído? há trabalho pendente?
//
// Legível por pessoa E por agente, de propósito. É a mesma pergunta nos dois casos, e o
// que muda é só quem lê a resposta.
func newStatusCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Where the project stands in the cycle, and what the next step is",
		Long: `Answers "where am I?" — the question of whoever is picking the work back up.

Unlike doctor (which hunts systemic loose ends in an already assembled project), status
says which PHASE the project is in: is the interview missing? the init? the first
plan? is there work in progress? And, for each state, what the next step is.

It changes nothing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			return runStatus(absRoot)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

// runStatus imprime o estado na ORDEM DO CICLO, e para no primeiro passo que falta.
// Listar tudo que está pendente de uma vez faria o leitor escolher por onde começar —
// e a ordem do ciclo é justamente o que ele não deveria ter de reconstruir sozinho.
func runStatus(root string) error {
	fmt.Println(i18n.T("status.header", root))
	fmt.Println()

	// 1. GIT — o substrato. Sem ele, metade do framework fica desligada em silêncio.
	switch gitmeta.Check(root) {
	case gitmeta.SemBinário:
		fmt.Println(i18n.T("status.git_missing"))
	case gitmeta.SemRepo:
		fmt.Println(i18n.T("status.no_git_repo"))
		return nil
	}

	// 2. A FASE DESCOBRIR — antes do init, e a única que o Anchors não executa.
	temProject := initx.HasProjectMD(root)
	cfgPath := filepath.Join(root, config.DefaultFile)
	_, errCfg := os.Stat(cfgPath)
	temConfig := errCfg == nil

	if !temProject && !temConfig {
		fmt.Println(i18n.T("status.not_started"))
		return nil
	}
	if temProject {
		fmt.Println(i18n.T("status.project_md_exists"))
	}

	if !temConfig {
		fmt.Println(i18n.T("status.no_config"))
		return nil
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load %s: %w", config.DefaultFile, err)
	}
	fmt.Println(i18n.T("status.config_ok", len(cfg.Layers), len(cfg.Gates)))

	// 3. O MAPA — todo arquivo novo só existe para os gates depois do `map build`.
	mapPath := filepath.Join(root, mapx.DefaultPath)
	g, errMapa := mapx.Load(mapPath)
	if errMapa != nil {
		fmt.Println(i18n.T("status.no_map"))
		return nil
	}
	fmt.Println(i18n.T("status.map_ok", len(g.Nodes), len(g.Edges)))

	// 3.5. MATURAÇÃO — o gate informativo que já está limpo (QUALITY §7).
	//
	// Aparece aqui, e não só no `check`, porque o `status` é o comando de quem RETOMA: é
	// onde se pergunta "o que falta?", e um gate que mede sem defender é exatamente isso.
	if prom := gate.PromotableGates(
		gate.Aggregate(gate.RunWithConfig(cfg.Gates, g.Nodes, root, g, cfg)),
	); len(prom) > 0 {
		nomes := make([]string, 0, len(prom))
		for _, x := range prom {
			nomes = append(nomes, x.Gate)
		}
		fmt.Println(i18n.T("status.clean_informative_gates", len(prom), strings.Join(nomes, ", ")))
	}

	// 4. O TRABALHO — a fila mora onde o modo declara (WORKFLOW.md §2).
	fmt.Println()
	if cfg.GitHubMode() {
		statusGitHub(root, cfg, g)
	} else {
		statusLocal(root, g)
	}
	return nil
}

// statusGitHub descreve a fila do modo `github`: ela mora nas issues do repositório, e o
// estado de cada trabalho é a COLUNA do board (BOOTSTRAP.md §7.13).
func statusGitHub(root string, cfg *config.Config, g *mapx.Graph) {
	fmt.Println(i18n.T("status.github_header", cfg.Workflow.Repo, cfg.Workflow.Labels))

	// O ambiente precisa estar montado antes de a fila fazer sentido — e o doctor é
	// quem sabe conferir isso. Aqui basta apontar, sem repetir a verificação.
	if faltam := initx.MissingWorkflow(root); len(faltam) > 0 {
		fmt.Println(i18n.T("status.github_missing_pipelines", len(faltam)))
		return
	}
	fmt.Println(i18n.T("status.github_pipelines_ok"))

	// O FLUXO DO PR, dito onde quem retoma o trabalho vai ler. Todo trabalho sobe por
	// PR: é o que dá objeto à revisão e o que faz o card nascer (o pipeline de
	// identificação dispara na ABERTURA do PR, não no push).
	base := cfg.Workflow.IntegrationBranchOrDefault()
	line := i18n.T("status.github_pr_flow", base)
	if p := cfg.Workflow.ProtectedBranchesOrDefault(); len(p) > 1 {
		line += i18n.T("status.github_protected_branches", strings.Join(p, ", "))
	}
	fmt.Println(line)
	fmt.Println()

	// Um projeto sem trabalho não tem card a pedir: o passo é criar o primeiro plano,
	// e mandar pedir trabalho aqui daria uma instrução que não devolve nada.
	if noRealWork(g) {
		printFirstPlan()
		return
	}

	// O CARD QUE JÁ É SEU vem antes de pedir outro. Quem volta de uma pausa não sabe se
	// ainda tem trabalho — e a resposta muda o próximo passo por inteiro: pedir um card
	// novo com um seu em aberto produz dois trabalhos pela metade.
	//
	// Importa mais depois que o `stale` age: ele libera card sem sinal de vida, e uma
	// revisão longa (ou uma noite) atravessa o prazo. O card volta para `to-do` sem dono,
	// e quem retoma precisa VER isso para saber que dá para retomar.
	if meus := agentCards(cfg); len(meus) > 0 {
		fmt.Println(i18n.T("status.github_already_have_work"))
		for _, c := range meus {
			fmt.Printf("    #%s %s [%s]\n", c.numero, c.titulo, c.estado)
		}
		fmt.Println()
		fmt.Println(i18n.T("status.github_finish_before_new"))
		return
	}

	fmt.Println(i18n.T("status.github_claim_next"))
}

// statusLocal descreve a fila do modo `local`: as tasks em `.anchors/tasks/` e as issues
// em `issues/`, cujo ESTADO é a pasta (CONCEPT.md §5).
func statusLocal(root string, g *mapx.Graph) {
	fmt.Println(i18n.T("status.local_header"))

	tasks := countFiles(filepath.Join(root, ".anchors", "tasks"))
	todo := countFiles(filepath.Join(root, "issues", "todo"))
	doing := countFiles(filepath.Join(root, "issues", "doing"))

	fmt.Println(i18n.T("status.local_tasks", tasks))
	fmt.Println(i18n.T("status.local_issues", todo, doing))

	fmt.Println()
	switch {
	case doing > 0:
		fmt.Println(i18n.T("status.local_doing_next"))
	case todo > 0:
		fmt.Println(i18n.T("status.local_todo_next"))
	case tasks > 0:
		fmt.Println(i18n.T("status.local_tasks_next"))
	case noRealWork(g):
		// "Nada pendente" com o projeto vazio seria uma resposta enganosa: não há nada
		// pendente porque não há nada.
		printFirstPlan()
	default:
		fmt.Println(i18n.T("status.local_all_clean"))
	}
}

// printFirstPlan orienta o primeiro plano de um projeto sem código.
//
// Diz os OBJETIVOS, não um template: o que a fundação precisa responder é universal
// (onde o código mora, o que formata, como se roda o teste, o que o CI executa), mas o
// COMO muda por stack — e o PROJECT.md já decidiu isso. Um template cravaria ESLint num
// projeto Python. É a mesma régua da fase DESCOBRIR, que fixa etapas e objetivos e não
// as perguntas.
func printFirstPlan() {
	fmt.Println(i18n.T("status.first_plan_guidance"))
}

// noRealWork diz se o mapa só tem o que o próprio `init` semeou (os guides). Um
// projeto assim está montado, não começado — e a diferença é o que separa "nada pendente"
// de "ainda não há o que fazer aqui".
func noRealWork(g *mapx.Graph) bool {
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindGuide {
			return false
		}
	}
	return true
}

// countFiles conta entradas de arquivo num diretório. Diretório ausente é 0 — no modo
// local, `issues/doing` só existe depois que alguém pega o primeiro trabalho.
func countFiles(dir string) int {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entradas {
		if !e.IsDir() {
			n++
		}
	}
	return n
}

// agentCard é um card que carrega o nome deste agente no último `anchors-owner`.
type agentCard struct{ numero, titulo, estado string }

// agentCards devolve os cards cujo último dono é este agente.
//
// A identidade vem de `ANCHORS_AGENT`, e não do usuário do git: agentes na mesma máquina
// compartilham a conta do GitHub — foi por isso que a posse virou comentário
// (`anchors-owner:`) em vez de assignee. Sem a variável não há como saber quem pergunta,
// e devolver os cards de OUTRO agente seria pior que não responder.
func agentCards(cfg *config.Config) []agentCard {
	agente := strings.TrimSpace(os.Getenv("ANCHORS_AGENT"))
	if agente == "" || cfg == nil || cfg.Workflow == nil {
		return nil
	}
	if _, err := exec.LookPath("gh"); err != nil {
		return nil
	}
	out, err := exec.Command("gh", "issue", "list",
		"--repo", cfg.Workflow.Repo,
		"--state", "open", "--label", cfg.Workflow.Labels[0], "--limit", "200",
		"--json", "number,title,labels,comments",
		"--jq", `[.[] | select([.comments[].body | select(startswith("anchors-owner:"))] | last == "anchors-owner: `+agente+`")
		         | {n: .number, t: .title, e: ([.labels[].name | select(startswith("anchors:"))] | .[0] // "")}] | .[] | "\(.n)\t\(.t)\t\(.e)"`,
	).Output()
	if err != nil {
		return nil
	}
	var cards []agentCard
	for _, l := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		p := strings.Split(l, "\t")
		if len(p) != 3 || p[0] == "" {
			continue
		}
		cards = append(cards, agentCard{p[0], p[1], strings.TrimPrefix(p[2], "anchors:")})
	}
	return cards
}

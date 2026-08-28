package main

import (
	"fmt"
	"os"
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
		Short: "Onde o projeto está no ciclo, e qual é o próximo passo",
		Long: `Responde "onde estou?" — a pergunta de quem retoma o trabalho.

Diferente do doctor (que caça pontas sistêmicas num projeto já montado), o status
diz em que FASE o projeto está: falta a entrevista? falta o init? falta o primeiro
plano? há trabalho em andamento? E, para cada estado, qual é o passo seguinte.

Não altera nada.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			return runStatus(absRoot)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	return cmd
}

// runStatus imprime o estado na ORDEM DO CICLO, e para no primeiro passo que falta.
// Listar tudo que está pendente de uma vez faria o leitor escolher por onde começar —
// e a ordem do ciclo é justamente o que ele não deveria ter de reconstruir sozinho.
func runStatus(root string) error {
	fmt.Printf("anchors status — %s\n\n", root)

	// 1. GIT — o substrato. Sem ele, metade do framework fica desligada em silêncio.
	switch gitmeta.Verifica(root) {
	case gitmeta.SemBinário:
		fmt.Println("⚠ git não instalado — o carimbo de alteração, `coverage --diff` e os")
		fmt.Println("  hooks ficam desligados. Instale o git.")
	case gitmeta.SemRepo:
		fmt.Println("⚠ sem repositório git.")
		fmt.Println("  → PRÓXIMO PASSO: `git init` (ou rode `anchors init`, que oferece fazê-lo)")
		return nil
	}

	// 2. A FASE DESCOBRIR — antes do init, e a única que o Anchors não executa.
	temProject := initx.TemProjectMD(root)
	cfgPath := filepath.Join(root, config.DefaultFile)
	_, errCfg := os.Stat(cfgPath)
	temConfig := errCfg == nil

	if !temProject && !temConfig {
		fmt.Println("○ projeto ainda não iniciado: sem PROJECT.md e sem anchors.yaml.")
		fmt.Println()
		fmt.Println("  → PRÓXIMO PASSO: a fase DESCOBRIR — uma entrevista de 5 etapas que decide")
		fmt.Println("    stack, arquitetura, estrutura e convenções, e escreve PROJECT.md.")
		fmt.Println("    Quem conduz é uma IA: rode `anchors guide project` para a régua,")
		fmt.Println("    ou `anchors init`, que reconhece o estado e instrui.")
		return nil
	}
	if temProject {
		fmt.Println("✓ PROJECT.md existe (a fase DESCOBRIR aconteceu)")
	}

	if !temConfig {
		fmt.Println("○ sem anchors.yaml")
		fmt.Println()
		fmt.Println("  → PRÓXIMO PASSO: `anchors init` — as perguntas dele são respondidas pelo")
		fmt.Println("    que o PROJECT.md decidiu.")
		return nil
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("carregar %s: %w", config.DefaultFile, err)
	}
	fmt.Printf("✓ anchors.yaml (%d camadas, %d gates)\n", len(cfg.Layers), len(cfg.Gates))

	// 3. O MAPA — todo arquivo novo só existe para os gates depois do `map build`.
	mapPath := filepath.Join(root, mapx.DefaultPath)
	g, errMapa := mapx.Load(mapPath)
	if errMapa != nil {
		fmt.Println("○ sem mapa")
		fmt.Println()
		fmt.Println("  → PRÓXIMO PASSO: `anchors map build` — sem o mapa, nenhum arquivo existe")
		fmt.Println("    para os gates.")
		return nil
	}
	fmt.Printf("✓ mapa (%d nós, %d arestas)\n", len(g.Nodes), len(g.Edges))

	// 3.5. MATURAÇÃO — o gate informativo que já está limpo (QUALITY §7).
	//
	// Aparece aqui, e não só no `check`, porque o `status` é o comando de quem RETOMA: é
	// onde se pergunta "o que falta?", e um gate que mede sem defender é exatamente isso.
	if prom := gate.GatesPromoviveis(
		gate.Aggregate(gate.RunWithConfig(cfg.Gates, g.Nodes, root, g, cfg)),
	); len(prom) > 0 {
		nomes := make([]string, 0, len(prom))
		for _, x := range prom {
			nomes = append(nomes, x.Gate)
		}
		fmt.Printf("○ %d gate(s) informativo(s) LIMPO(s): %s\n", len(prom), strings.Join(nomes, ", "))
		fmt.Println("  medem e não defendem — promova a `blocking: true` no anchors.yaml")
	}

	// 4. O TRABALHO — a fila mora onde o modo declara (WORKFLOW.md §2).
	fmt.Println()
	if cfg.ModoGitHub() {
		statusGitHub(root, cfg, g)
	} else {
		statusLocal(root, g)
	}
	return nil
}

// statusGitHub descreve a fila do modo `github`: ela mora nas issues do repositório, e o
// estado de cada trabalho é a COLUNA do board (BOOTSTRAP.md §7.13).
func statusGitHub(root string, cfg *config.Config, g *mapx.Graph) {
	fmt.Printf("fila: GitHub (%s, label %v)\n", cfg.Workflow.Repo, cfg.Workflow.Labels)

	// O ambiente precisa estar montado antes de a fila fazer sentido — e o doctor é
	// quem sabe conferir isso. Aqui basta apontar, sem repetir a verificação.
	if faltam := initx.FaltaWorkflow(root); len(faltam) > 0 {
		fmt.Printf("⚠ %d pipeline(s) do fluxo ausentes — sem eles o ciclo não avança sozinho.\n", len(faltam))
		fmt.Println("  → `anchors doctor --fix` cria os que faltam")
		return
	}
	fmt.Println("✓ pipelines do fluxo no lugar")

	// O FLUXO DO PR, dito onde quem retoma o trabalho vai ler. Todo trabalho sobe por
	// PR: é o que dá objeto à revisão e o que faz o card nascer (o pipeline de
	// identificação dispara na ABERTURA do PR, não no push).
	base := cfg.Workflow.BranchDeIntegracao()
	fmt.Printf("  trabalho entra por PR para `%s`", base)
	if p := cfg.Workflow.BranchesProtegidos(); len(p) > 1 {
		fmt.Printf(" · protegidos: %s", strings.Join(p, ", "))
	}
	fmt.Println()
	fmt.Println()

	// Um projeto sem trabalho não tem card a pedir: o passo é criar o primeiro plano,
	// e mandar pedir trabalho aqui daria uma instrução que não devolve nada.
	if semTrabalhoReal(g) {
		imprimePrimeiroPlano()
		return
	}

	fmt.Println("  → PRÓXIMO PASSO: peça trabalho ao pipeline de claim —")
	fmt.Println("    `gh workflow run anchors-claim.yml -f agent=$(hostname)/<sessao>`")
	fmt.Println()
	fmt.Println("  a prioridade do board é da DIREITA para a ESQUERDA: termine o que está")
	fmt.Println("  mais adiantado antes de pegar coisa nova.")
}

// statusLocal descreve a fila do modo `local`: as tasks em `.anchors/tasks/` e as issues
// em `issues/`, cujo ESTADO é a pasta (CONCEPT.md §5).
func statusLocal(root string, g *mapx.Graph) {
	fmt.Println("fila: local")

	tasks := contaArquivos(filepath.Join(root, ".anchors", "tasks"))
	todo := contaArquivos(filepath.Join(root, "issues", "todo"))
	doing := contaArquivos(filepath.Join(root, "issues", "doing"))

	fmt.Printf("  tasks pendentes: %d\n", tasks)
	fmt.Printf("  issues: %d em todo, %d em doing\n", todo, doing)

	fmt.Println()
	switch {
	case doing > 0:
		fmt.Println("  → PRÓXIMO PASSO: há trabalho em `issues/doing` — termine o que já começou")
		fmt.Println("    antes de pegar coisa nova.")
	case todo > 0:
		fmt.Println("  → PRÓXIMO PASSO: `issues/todo` tem trabalho — leia e escolha um.")
	case tasks > 0:
		fmt.Println("  → PRÓXIMO PASSO: `anchors next` puxa a próxima task da fila.")
	case semTrabalhoReal(g):
		// "Nada pendente" com o projeto vazio seria uma resposta enganosa: não há nada
		// pendente porque não há nada.
		imprimePrimeiroPlano()
	default:
		fmt.Println("  → nada pendente. `anchors doctor` mostra as pontas sistêmicas.")
	}
}

// imprimePrimeiroPlano orienta o primeiro plano de um projeto sem código.
//
// Diz os OBJETIVOS, não um template: o que a fundação precisa responder é universal
// (onde o código mora, o que formata, como se roda o teste, o que o CI executa), mas o
// COMO muda por stack — e o PROJECT.md já decidiu isso. Um template cravaria ESLint num
// projeto Python. É a mesma régua da fase DESCOBRIR, que fixa etapas e objetivos e não
// as perguntas.
func imprimePrimeiroPlano() {
	fmt.Println("  → PRÓXIMO PASSO: o projeto está montado e ainda sem trabalho.")
	fmt.Println()
	fmt.Println("  O primeiro plano é o de FUNDAÇÃO, e vem antes de qualquer feature: sem")
	fmt.Println("  estrutura de pastas, linter, formatador e teste rodando, o primeiro plano")
	fmt.Println("  de produto esbarra nisso no primeiro arquivo — e cada agente resolve à sua")
	fmt.Println("  maneira. Ele precisa responder, com o que o PROJECT.md decidiu:")
	fmt.Println()
	fmt.Println("    · onde o código mora — a árvore que o `anchors.yaml` declara")
	fmt.Println("    · o que formata e o que linta, e a configuração de cada um")
	fmt.Println("    · como se roda o teste, e como o sinal chega ao `anchors ingest`")
	fmt.Println("    · o que o CI executa a cada push")
	fmt.Println()
	fmt.Println("  Rode `anchors guide plan` para a régua.")
	fmt.Println()
	fmt.Println("  Vale escrever a SEQUÊNCIA inteira de planos de uma vez — é o que permite")
	fmt.Println("  reconciliar escopo entre eles — e liberar aos poucos: declare `needs:` no")
	fmt.Println("  header de cada um e commite só o que já pode ser trabalhado. Um plano cujo")
	fmt.Println("  `needs` não terminou não vira card.")
}

// semTrabalhoReal diz se o mapa só tem o que o próprio `init` semeou (os guides). Um
// projeto assim está montado, não começado — e a diferença é o que separa "nada pendente"
// de "ainda não há o que fazer aqui".
func semTrabalhoReal(g *mapx.Graph) bool {
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindGuide {
			return false
		}
	}
	return true
}

// contaArquivos conta entradas de arquivo num diretório. Diretório ausente é 0 — no modo
// local, `issues/doing` só existe depois que alguém pega o primeiro trabalho.
func contaArquivos(dir string) int {
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

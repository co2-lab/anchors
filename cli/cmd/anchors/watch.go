package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/co2-lab/anchors/internal/change"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/daemon"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "O watcher em background: vê mudanças e ENFILEIRA trabalho",
		Long: `O watcher (PROPAGATION §6) observa o repo e, a cada arquivo que muda,
classifica a mudança e ENFILEIRA uma task em .anchors/tasks/ (com o próximo passo
sugerido). Ele NÃO executa nem chama IA — quem trabalha puxa a task com 'anchors
next'. É a inversão: o watcher só transforma "mudou" em "há trabalho".

Roda em BACKGROUND — o terminal fica livre. Controle com os subcomandos:
  anchors watch start    inicia o watcher (daemoniza, retorna o terminal)
  anchors watch status   está rodando? desde quando? pausado?
  anchors watch stop     encerra
  anchors watch pause    pausa (para de reagir, sem encerrar)
  anchors watch resume   retoma
  anchors watch logs     mostra o que o watcher reportou`,
	}
	cmd.AddCommand(newWatchStartCmd(), newWatchStatusCmd(), newWatchStopCmd(),
		newWatchPauseCmd(), newWatchResumeCmd(), newWatchLogsCmd(), newWatchRunCmd())
	return cmd
}

// rootFlag resolve a raiz absoluta a partir da flag --root de um subcomando.
func absRootFlag(cmd *cobra.Command) (string, error) {
	root, _ := cmd.Flags().GetString("root")
	return config.AbsRaiz(root)
}

func addRootFlag(cmd *cobra.Command) {
	cmd.Flags().String("root", ".", "raiz do projeto")
}

// --- start: daemoniza (re-exec desanexado) e retorna o terminal ---

func newWatchStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Inicia o watcher em background",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := absRootFlag(cmd)
			if err != nil {
				return err
			}
			p := daemon.PathsFor(root)
			if pid := daemon.Running(p); pid != 0 {
				return fmt.Errorf("watcher já está rodando (pid %d)", pid)
			}
			if err := os.MkdirAll(p.Dir, 0o755); err != nil {
				return err
			}

			// re-executa a si mesmo no modo `run` (foreground), desanexado da
			// sessão (Setsid) e com stdout/stderr no log. O pai retorna já.
			logf, err := os.OpenFile(p.Log, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
			if err != nil {
				return err
			}
			defer logf.Close()

			self, err := os.Executable()
			if err != nil {
				return err
			}
			child := exec.Command(self, "watch", "run", "--root", root)
			child.Stdout = logf
			child.Stderr = logf
			detach(child) // desanexa da sessão do terminal (implementação por plataforma)
			if err := child.Start(); err != nil {
				return err
			}
			pid := child.Process.Pid
			if err := daemon.WritePID(p, pid); err != nil {
				return err
			}
			_ = daemon.WriteMeta(p, time.Now(), root)
			// não espera o filho — desanexa
			fmt.Printf("watcher iniciado em background (pid %d)\n", pid)
			fmt.Printf("  logs:   anchors watch logs\n  status: anchors watch status\n  parar:  anchors watch stop\n")
			return nil
		},
	}
	addRootFlag(cmd)
	return cmd
}

// --- run: o loop em foreground (usado pelo daemon; útil p/ debug direto) ---

func newWatchRunCmd() *cobra.Command {
	var debounceMs int
	cmd := &cobra.Command{
		Use:    "run",
		Short:  "Roda o loop do watcher em foreground (uso interno do daemon)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := absRootFlag(cmd)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar config: %w", err)
			}
			g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("carregar mapa: %w", err)
			}
			return runWatchLoop(root, cfg, g, time.Duration(debounceMs)*time.Millisecond)
		},
	}
	addRootFlag(cmd)
	cmd.Flags().IntVar(&debounceMs, "debounce", 300, "janela de debounce em ms")
	return cmd
}

// --- controles ---

func newWatchStatusCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "status", Short: "Estado do watcher",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := absRootFlag(cmd)
			p := daemon.PathsFor(root)
			pid := daemon.Running(p)
			if pid == 0 {
				fmt.Println("watcher: parado")
				return nil
			}
			state := "rodando"
			if daemon.IsPaused(p) {
				state = "pausado"
			}
			fmt.Printf("watcher: %s (pid %d)\n", state, pid)
			if meta := daemon.ReadMeta(p); meta != "" {
				fmt.Print(meta)
			}
			return nil
		}}
	addRootFlag(cmd)
	return cmd
}

func newWatchStopCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "stop", Short: "Encerra o watcher",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := absRootFlag(cmd)
			if err := daemon.Stop(daemon.PathsFor(root)); err != nil {
				return err
			}
			fmt.Println("watcher encerrado")
			return nil
		}}
	addRootFlag(cmd)
	return cmd
}

func newWatchPauseCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "pause", Short: "Pausa o watcher (sem encerrar)",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := absRootFlag(cmd)
			p := daemon.PathsFor(root)
			if daemon.Running(p) == 0 {
				return fmt.Errorf("watcher não está rodando")
			}
			if err := daemon.Pause(p); err != nil {
				return err
			}
			fmt.Println("watcher pausado (resume para retomar)")
			return nil
		}}
	addRootFlag(cmd)
	return cmd
}

func newWatchResumeCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "resume", Short: "Retoma o watcher pausado",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := absRootFlag(cmd)
			_ = daemon.Resume(daemon.PathsFor(root))
			fmt.Println("watcher retomado")
			return nil
		}}
	addRootFlag(cmd)
	return cmd
}

func newWatchLogsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "logs", Short: "Mostra o log do watcher",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := absRootFlag(cmd)
			data, err := os.ReadFile(daemon.PathsFor(root).Log)
			if err != nil {
				return fmt.Errorf("sem log (watcher já rodou?): %w", err)
			}
			fmt.Print(string(data))
			return nil
		}}
	addRootFlag(cmd)
	return cmd
}

// --- o loop ---

// watchIgnore é carregado do `.gitignore` do projeto + a lista universal, uma vez ao
// iniciar o daemon. Ver internal/scan/ignore.go.
var watchIgnore *scan.Ignore

func runWatchLoop(root string, cfg *config.Config, g *mapx.Graph, debounce time.Duration) error {
	// Carregado uma vez: o `.gitignore` não muda no meio de uma sessão de trabalho, e
	// relê-lo a cada evento tornaria o watcher mais caro sem ganho.
	// `LoadIgnoreFor` e não `LoadIgnore`: a Estrutura pode DERROTAR a lista embutida (um
	// projeto onde `build/` é código-fonte), e o watcher precisa enxergar o mesmo que o
	// `check`. Havia aqui uma segunda cópia da lista universal, já divergente da viva —
	// código morto que ainda contradizia a régua.
	watchIgnore = scan.LoadIgnoreFor(root, cfg)
	p := daemon.PathsFor(root)
	defer daemon.Cleanup(p)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	added := addTreeToWatcher(w, root)
	fmt.Printf("anchors watch — observando %d diretórios sob %s\n", added, root)

	// SIGTERM/SIGINT → encerra limpo (o defer Cleanup roda).
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	pending := map[string]*time.Timer{}
	for {
		select {
		case sig := <-sigc:
			fmt.Printf("recebido %s — encerrando\n", sig)
			return nil
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}
			// Diretório NOVO: passa a observá-lo (recursivamente). Sem isto, um
			// arquivo criado numa pasta nascida após o start (ex.: plans/) nunca
			// gera evento — a fila silenciava para trabalho novo em pastas novas.
			if ev.Op&fsnotify.Create != 0 {
				if fi, err := os.Stat(ev.Name); err == nil && fi.IsDir() {
					addTreeToWatcher(w, ev.Name)
					// corrida: arquivos podem ter nascido na pasta ANTES do Add
					// efetivar (não gerariam evento). Varre e processa os que já
					// existem, para não perder trabalho na janela de registro.
					for _, f := range filesBornIn(ev.Name) {
						rel, _ := filepath.Rel(root, f)
						handleChange(root, cfg, g, rel)
					}
				}
			}
			if daemon.IsPaused(p) {
				continue // pausado: ignora eventos
			}
			rel, _ := filepath.Rel(root, ev.Name)
			if pending[rel] != nil {
				pending[rel].Stop()
			}
			relCopy := rel
			pending[relCopy] = time.AfterFunc(debounce, func() {
				handleChange(root, cfg, g, relCopy)
			})
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintln(os.Stderr, "watch erro:", err)
		}
	}
}

// addTreeToWatcher registra `base` e todos os seus subdiretórios (fora os ignorados)
// no fsnotify, devolvendo quantos foram adicionados. Usado no start (a árvore toda) e
// quando um diretório NOVO nasce durante a execução (recursão para pegar netos).
func addTreeToWatcher(w *fsnotify.Watcher, base string) int {
	added := 0
	_ = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			rel, _ := filepath.Rel(base, path)
			if watchIgnore.SkipDir(d.Name(), rel) {
				return filepath.SkipDir
			}
			if w.Add(path) == nil {
				added++
			}
		}
		return nil
	})
	return added
}

// filesBornIn lista os arquivos (não-dirs) diretamente sob dir — usado para varrer
// uma pasta recém-criada e não perder arquivos nascidos na janela de corrida entre
// o Create do dir e o Add ao fsnotify. Só o primeiro nível (subdirs geram seu
// próprio evento Create quando o Add já está ativo).
func filesBornIn(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out
}

// handleChange: o watcher ENFILEIRA. Classifica a mudança, deriva o próximo passo
// do ciclo de vida (SuggestNext), e escreve uma task em .anchors/tasks/. Não executa
// check nem chama IA — isso é trabalho do WORKER, que puxa a task com `anchors next`.
// Essa é a inversão: o watcher só transforma "algo mudou" em "há trabalho na fila";
// quem trabalha (a IA, via subagente ou sessão bloqueante) puxa quando puder.
func handleChange(root string, cfg *config.Config, g *mapx.Graph, rel string) {
	// O que o projeto declara descartável no `.gitignore` não vira trabalho.
	//
	// Medido: um revisor criou sondas (`probe1.ts`) para atacar a unidade por execução, e
	// o watcher as enfileirou — cinco tasks descartadas à mão pelo orquestrador. Eram
	// instrumento de review, não entrega, e o `.gitignore` já dizia isso.
	if watchIgnore != nil && watchIgnore.SkipFile(rel) {
		return
	}
	// Um REGISTRO DE ENTREGA (`changes/*.md`) é o gatilho do review, e não depende da
	// Estrutura do projeto: é artefato do próprio Anchors, como as issues. Sem este
	// caminho, o ciclo terminava quando o código nascia e ninguém confrontava o que foi
	// entregue — medido: 7 defeitos graves atravessaram três rodadas de um E2E com todos
	// os gates verdes, e os 7 saíram de revisão adversarial que alguém teve de lembrar
	// de pedir.
	if unidade, ok := changeDelivered(root, rel); ok {
		// A entrega de um PLANO fecha o conjunto, e o review dela é de outra natureza.
		//
		// Medido na rodada 4 de um E2E real: dos 6 achados de um review adversarial, apenas
		// 2 cabiam no escopo de UMA unidade. Os outros 4 atravessavam — o DAO que trunca e
		// quebra o carry-forward da regra que o consome; o campo `owner` que a spec do
		// modelo manda plantar, o DAO não planta e o handler não preenche; a obrigação de
		// exclusão que passa estruturalmente e some na mutação. Em todos, CADA PEÇA está
		// correta sozinha, e é por isso que o review por unidade não os alcança.
		//
		// Os dois níveis são complementares, não redundantes: o de unidade garante a
		// parcial (e chega cedo, quando corrigir é barato); o de plano garante o conjunto
		// (e é o único que vê a costura entre as peças).
		if planoDeEntrega(root, rel) {
			enqueueTask(root, rel, "change", "review-plan",
				"um PLANO foi entregue — REVISE O CONJUNTO (`anchors work review-plan --for "+
					unidade+"`). O review por unidade já garantiu cada peça; este garante a "+
					"COSTURA: o que uma camada promete e outra não cumpre, a regra que "+
					"contradiz a de cima, o dever que passa estruturalmente em cada peça e "+
					"não funciona no todo. Medido: 4 de 6 achados de um review real "+
					"atravessavam unidades, com cada peça correta sozinha.")
			return
		}
		// O review de unidade ATACA por execução — mutar a regra, rodar com entrada de
		// borda, ver se o teste cai. Isso exige que exista código e teste. Disparado na
		// entrega da SPEC, o prompt manda mutar um arquivo que ainda não nasceu, e quem
		// recebe a task fica sem chão (medido: aconteceu na primeira rodada em que o ciclo
		// completo funcionou; o orquestrador teve de segurar o review na mão).
		//
		// Então a entrega de uma peça isolada REGISTRA e espera; o review entra quando a
		// trinca fecha. O sinal de que fechou é a peça que nasce por último: o TESTE.
		if faltaPecaParaRevisar(root, unidade, cfg) {
			fmt.Printf("● %s [change] — entrega registrada; o review espera a trinca fechar "+
				"(falta código e/ou teste em `%s`)\n", rel, unidade)
			return
		}
		enqueueTask(root, rel, "change", "review",
			"uma entrega foi registrada e a TRINCA FECHOU — REVISE a unidade `"+unidade+"`: "+
				"os gates verdes não provam que está certo. Ataque por execução (mute a regra "+
				"e veja se o teste cai; rode com entrada de borda), não por leitura.")
		return
	}

	// O que o projeto (ou o sistema) declara descartável não vira trabalho. O `Classify`
	// sozinho não basta: um `.!21662!resource.spec.md`, que o editor cria por
	// milissegundos durante um salvamento atômico, casa o glob de spec — e virava task.
	// Medido num E2E real: 3 tasks-fantasma descartadas à mão (o temporário do editor e
	// duas sondas de revisor já apagadas).
	if watchIgnore.SkipFile(rel) {
		return
	}
	layer, kind := scan.Classify(rel, cfg)
	if layer == "" {
		return // arquivo fora da Estrutura — ignora
	}
	// Arquivo que não existe mais não é trabalho: entre o evento e o processamento, o
	// temporário já foi renomeado e a sonda já foi apagada.
	if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
		return
	}
	updateNodeRev(g, root, rel)

	next, reason := queue.SuggestNext(kind)
	// A camada pode DISPENSAR a peça que a fila sugere (`trinca_opcional`). Enfileirá-la
	// é a fila contradizendo o `anchors work`, que para o mesmo alvo responde `PARE — NÃO
	// crie`. Medido em duas execuções: tasks de `feature` para camada `dao` e para
	// `schema-model`, descartadas à mão pelo orquestrador.
	//
	// Pular para a etapa seguinte preserva a cadeia — quem dispensa a feature ainda quer
	// o teste (ou o review), e parar a corrente na peça dispensada deixaria a unidade sem
	// próxima etapa nenhuma.
	if camada := camadaDoAlvoDaPeca(root, rel, cfg); camada != "" {
		disp := pecasDispensadas(camada, cfg)
		// Pula a peça DISPENSADA pela camada e a que JÁ EXISTE no disco.
		//
		// A segunda foi relatada em três execuções: a fila mandava "descreva o
		// comportamento em cenários" para uma unidade cuja `.feature` estava escrita,
		// ligada e verde. Um worker obediente reescreve o que já estava pronto; um
		// atento descarta a task à mão — e nos dois casos a fila perdeu a confiança de
		// quem a puxa.
		//
		// O caso legítimo (a peça existe e PRECISA mudar porque a spec mudou) já é
		// coberto: quem enfileira ali é o `stale`, pela aresta com rev avançada, não
		// esta sugestão de cadeia.
		for i := 0; i < 3 && next != ""; i++ {
			if !disp[next] && !pecaJaExiste(root, rel, next, cfg) {
				break
			}
			next, reason = queue.SuggestNext(next)
		}
	}
	if next == "" {
		return
	}
	enqueueTask(root, rel, kind, next, reason)
}

// camadaDoAlvoDaPeca resolve a camada da UNIDADE a partir de qualquer peça dela.
//
// É preciso porque a dispensa é declarada na camada do CÓDIGO (`schema-model`, `dao`),
// enquanto o evento chega numa peça derivada — e um `.spec.md` casa a camada `spec`, que
// não declara dispensa nenhuma. Classificar o caminho recebido devolveria sempre a camada
// errada.
func camadaDoAlvoDaPeca(root, rel string, cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	base := rel
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx"} {
		base = strings.TrimSuffix(base, suf)
	}
	if base == rel {
		if l, _ := scan.Classify(rel, cfg); l != "" {
			return l
		}
		return ""
	}
	for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
		if _, err := os.Stat(filepath.Join(root, base+ext)); err != nil {
			continue
		}
		if l, _ := scan.Classify(base+ext, cfg); l != "" {
			return l
		}
	}
	return ""
}

// taskID gera um ID estável e determinístico (sem time/random, que o projeto evita):
// <slug-do-arquivo>-<passo>-<hash-curto>. O hash torna o ID único por conteúdo do
// caminho+passo; o dedup do Enqueue evita duplicatas vivas de todo modo.
func taskID(rel, next string) string {
	// Os dois separadores, como o slug das issues já fazia: o rel chega com barra normal
	// (é id de nó) e no Windows o filepath.Separator sozinho não trocaria nada, deixando
	// uma barra dentro de um id que vira nome de arquivo na fila.
	slug := strings.ReplaceAll(rel, "/", "-")
	slug = strings.ReplaceAll(slug, string(filepath.Separator), "-")
	slug = strings.TrimSuffix(slug, filepath.Ext(slug))
	h := scan.ShortHash([]byte(rel + "|" + next))
	return fmt.Sprintf("%s-%s-%s", slug, next, h)
}

// nowStamp: o daemon é processo normal (não script de workflow), então time.Now é
// permitido aqui — só o carimbo de criação, informativo.
func nowStamp() string { return time.Now().Format(time.RFC3339) }

func updateNodeRev(g *mapx.Graph, root, rel string) {
	content, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		return
	}
	rev := scan.ShortHash(content)
	for i := range g.Nodes {
		if g.Nodes[i].ID == rel {
			g.Nodes[i].Rev = rev
			return
		}
	}
}

// changeDelivered diz se o arquivo é um REGISTRO DE ENTREGA pendente de review, e devolve
// a unidade que ele descreve. Só conta o que está na raiz de `changes/` — o que já foi
// revisado vive em `changes/reviewed/` e não deve reabrir a fila.
func changeDelivered(root, rel string) (string, bool) {
	dir := filepath.ToSlash(filepath.Dir(rel))
	if dir != change.Dir || !strings.HasSuffix(rel, ".md") {
		return "", false
	}
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		return "", false
	}
	if m := changeUnitRE.FindStringSubmatch(string(b)); m != nil {
		return m[1], true
	}
	return "", true // registro sem unidade legível: ainda assim é entrega, revise
}

var changeUnitRE = regexp.MustCompile(`(?m)^\s*unit:\s*(\S+)`)

// enqueueTask põe uma task na fila e reporta — o caminho comum de todo gatilho do watcher.
func enqueueTask(root, rel, kind, next, reason string) {
	t := queue.Task{
		ID:            taskID(rel, next),
		Changed:       rel,
		Kind:          kind,
		Origin:        "watch",
		SuggestedNext: next,
		Reason:        reason,
		CreatedAt:     nowStamp(),
	}
	created, err := queue.Enqueue(root, t)
	if err != nil {
		fmt.Printf("● %s [%s] — erro ao enfileirar: %v\n", rel, kind, err)
		return
	}
	if !created {
		fmt.Printf("● %s [%s] — já na fila (%s)\n", rel, kind, next)
		return
	}
	fmt.Printf("● %s [%s] → task enfileirada: %s (%s)\n", rel, kind, t.ID, next)
	if n, _ := queue.PendingCount(root); n > 0 {
		fmt.Printf("   fila: %d task(s) — um worker pode puxar com `anchors next`\n", n)
	}
}

// planoDeEntrega diz se o registro é o fecho de um PLANO (não de uma unidade).
func planoDeEntrega(root, rel string) bool {
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		return false
	}
	return regexp.MustCompile(`(?m)^\s*stage:\s*plan\b`).Match(b)
}

// faltaPecaParaRevisar diz se a unidade ainda não tem o material que o review ATACA:
// código e teste. Sem eles, o procedimento do review (mutar a regra, rodar a suíte) não
// tem alvo — e mandar alguém mutar um arquivo inexistente é pior que não disparar.
//
// Não usa o mapa: o arquivo pode ter nascido depois do último `map build`, e o review
// atrasaria por um detalhe de sincronização. Olha o disco, que é a verdade.
func faltaPecaParaRevisar(root, unidade string, cfg *config.Config) bool {
	if unidade == "" {
		return false // sem unidade legível: não bloqueia o review, que decidirá o que fazer
	}
	if _, err := os.Stat(filepath.Join(root, unidade)); err != nil {
		return true // o código da unidade ainda não existe
	}
	// A camada pode DISPENSAR o teste (`trinca_opcional: [tested-by]`). Ali a trinca
	// nunca fecha pelo sinal do teste — e esperar por ele significa que a unidade NUNCA é
	// revisada.
	//
	// Medido: numa camada de modelo com o teste dispensado, o review não foi enfileirado
	// nenhuma vez em duas execuções; a revisão manual dessa mesma unidade produziu 2
	// achados críticos e 6 sérios, todos invisíveis a 25 gates verdes. São exatamente as
	// camadas onde o gate tem menos a confrontar e o review mais a achar — e eram as
	// únicas que o ciclo não alcançava.
	if cfg != nil {
		if layer, _ := scan.Classify(unidade, cfg); layer != "" {
			for _, aresta := range cfg.Layers[layer].TrincaOpcional {
				if aresta == "tested-by" {
					return false // o código existe e o teste é dispensado: pode revisar
				}
			}
		}
	}
	// o teste é a peça que nasce por último — é o sinal de que a trinca fechou.
	base := strings.TrimSuffix(unidade, filepath.Ext(unidade))
	for _, suf := range []string{".test.ts", ".test.tsx", ".spec.ts", "_test.go", "_test.py", "_spec.rb", ".test.js"} {
		if _, err := os.Stat(filepath.Join(root, base+suf)); err == nil {
			return false
		}
	}
	return true
}

// pecaJaExiste diz se a peça que a fila ia sugerir já está escrita. Resolve o caminho pelo
// `derived:` do projeto — a mesma fonte que o `anchors work` usa para dizer onde a peça
// nasce, para que fila e prompt não discordem sobre o mesmo arquivo.
func pecaJaExiste(root, rel, peca string, cfg *config.Config) bool {
	if peca != "spec" && peca != "feature" && peca != "test" {
		return false // `code` e `review` não são peça derivada com caminho previsível
	}
	alvo := rel
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	for _, suf := range []string{".spec.md", ".feature", ".test", ".spec"} {
		base = strings.TrimSuffix(base, suf)
	}
	// o alvo da unidade é o arquivo de código; a peça deriva dele
	for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
		if _, err := os.Stat(filepath.Join(root, base+ext)); err == nil {
			alvo = base + ext
			break
		}
	}
	layer, _ := scan.Classify(alvo, cfg)
	files, _ := caminhosDerivados(alvo, layer, cfg)
	caminho := files[peca]
	if caminho == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(root, caminho))
	return err == nil
}

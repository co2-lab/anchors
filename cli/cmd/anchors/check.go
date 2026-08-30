package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/checklog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

func newCheckCmd() *cobra.Command {
	var root, mapPath, phase, category string
	var changed []string
	var all, noRecord, fix, deterministic, skipSlow, onlyIssues, showDrift bool
	var skipRegras string
	var msgPath string
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Roda os gates de qualidade (o pipeline)",
		Long: `Confronta os nós contra os gates de qualidade declarados no anchors.yaml
e reporta o perfil de vereditos (QUALITY §5-§6).

Por padrão é INCREMENTAL: roda sobre o caminho de impacto de um arquivo alterado
(--changed <arquivo>). Com --all, varre todos os nós (caro; a foto completa).

Cada gate que reprova gera uma issue; um gate bloqueante que reprova barra a
promoção. Gates informativos entram no perfil mas não barram (maturação, §7).

--deterministic PULA os gates de julgamento por IA (measures: judgment) e roda só
os computáveis. É o modo do PRE-COMMIT: um gate judge não pode barrar um commit
(não dá para esperar a IA) nem deve enfileirar julgamento a cada commit (lixo
repetido). Sem esse modo, judge fica invisível (nem barra, nem registra).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar config: %w", err)
			}
			// No modo github o achado de gate vira CARD, não arquivo: o `issues/` é a fila do
			// modo local (mover pasta à mão), e manter os dois faz o board esconder o que os
			// gates encontraram.
			if cfg != nil && cfg.ModoGitHub() && len(cfg.Workflow.Labels) > 0 {
				issue.UsarGitHub(cfg.Workflow.Repo, cfg.Workflow.Labels[0])
			}
			if len(cfg.Gates) == 0 {
				return fmt.Errorf("nenhum gate declarado no anchors.yaml (seção `gates:`)")
			}
			// --deterministic: remove os gates de julgamento por IA do conjunto.
			if deterministic {
				kept := cfg.Gates[:0:0]
				for _, gt := range cfg.Gates {
					if !gt.IsJudgment() {
						kept = append(kept, gt)
					}
				}
				cfg.Gates = kept
				if len(cfg.Gates) == 0 {
					fmt.Println("nenhum gate determinístico a rodar (todos são de julgamento).")
					return nil
				}
			}
			// Categorização (fase/custo/natureza): filtra QUAIS gates são cobrados agora.
			perspective := config.PerspectiveChange
			if all {
				perspective = config.PerspectiveAll
			}
			// A dispensa vem da flag OU do ambiente — o hook do git passa por variável,
			// porque `git commit` não repassa flags ao pre-commit.
			bruto := skipRegras
			if bruto == "" {
				bruto = os.Getenv("ANCHORS_SKIP_RULES")
			}
			dispensa, erros := gate.ParseDispensa(bruto)
			// A MENSAGEM DE COMMIT também dispensa: `[skip-trinca-completa@WRKSP: motivo]`.
			//
			// O caminho vem por `--commit-msg`, e não de `.git/COMMIT_EDITMSG`: MEDIDO, o
			// git NÃO grava esse arquivo antes do `pre-commit` — nem com `-m`. Lê-lo ali
			// devolve a mensagem do commit ANTERIOR, e a dispensa passa a valer para o
			// commit errado, em silêncio. É o hook `commit-msg` que recebe o arquivo, e é
			// ele quem passa este caminho.
			if msgPath != "" {
				if msg, err := os.ReadFile(msgPath); err == nil {
					daMsg, errosMsg := gate.DispensaDaMensagem(string(msg))
					dispensa = dispensa.Mescla(daMsg)
					erros = append(erros, errosMsg...)
				}
			}
			if len(erros) > 0 {
				return fmt.Errorf("--skip-rule inválido:\n  %s\n\nO motivo é obrigatório: uma "+
					"dispensa sem justificativa escrita é indistinguível de alguém fugindo de um "+
					"gate que achou defeito", strings.Join(erros, "\n  "))
			}
			cfg.Gates = filtrarGates(cfg.Gates, phase, category, skipSlow, perspective, dispensa)
			if len(cfg.Gates) == 0 {
				fmt.Printf("nenhum gate a rodar para este recorte (fase=%q categoria=%q).\n", phase, category)
				return nil
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}

			// Os arquivos que de fato MUDARAM, distintos do raio de impacto que o
			// `selectNodes` devolve. Um gate que julga a mudança precisa dos primeiros;
			// os demais, do raio. Ver Config.Alterados.
			cfg.Alterados = normalizaAlterados(changed, absRoot)
			nodes, scope, err := selectNodes(g, cfg, all, changed, absRoot)
			if err != nil {
				return err
			}

			// --fix: aplica os reparos automáticos (self-healer) ANTES de confrontar,
			// para que o check seguinte já reflita o conserto.
			if fix {
				fixes := gate.Fix(cfg.Gates, nodes, absRoot)
				n := 0
				for _, fr := range fixes {
					if fr.Fixed {
						fmt.Printf("  ✔ corrigido: %s @ %s — %s\n", fr.Gate, fr.Target, fr.Detail)
						n++
					}
				}
				fmt.Printf("%d reparo(s) aplicado(s).\n\n", n)
			}

			// Espelha a saída num arquivo para que ela possa ser RELIDA sem
			// re-executar. Abre AQUI, e não no topo: o que vem antes é erro de
			// invocação (config ausente, mapa ilegível), que não é relatório —
			// gravá-lo sobrescreveria uma foto boa com uma mensagem de erro.
			head, assunto, _ := gitmeta.Head(absRoot)
			espelho := checklog.Abrir(absRoot, all, checklog.Cabecalho(
				"anchors "+strings.Join(os.Args[1:], " "),
				head, assunto, gitmeta.DirtyCount(absRoot), time.Now(),
			))
			defer espelho.Fechar()

			warnIfMapStale(absRoot, mapPath, cfg)
			fmt.Printf("check %s — %d nós, %d gates\n\n", scope, len(nodes), len(cfg.Gates))

			// `all` chega até os gates: é o que permite ao gate que sabe varrer sozinho
			// (`scope_full`) rodar UMA vez sem receber a lista, em vez de receber o projeto
			// inteiro em lotes.
			results := gate.RunComDispensa(cfg.Gates, nodes, absRoot, g, cfg, all, dispensa)
			profile := gate.Aggregate(results)
			printProfile(profile, onlyIssues, showDrift)
			avisarGatesSemAlvo(cfg.Gates, profile)

			// O LOOP: check → carimbo → issue. Deixa de "reportar" e passa a
			// "registrar": grava o veredito por aresta no mapa (destrava stale) e
			// abre uma issue de violation por fail bloqueante (sobrevive à sessão).
			// Opt-out honesto: --no-record só reporta, não registra.
			if !noRecord {
				if err := recordCheck(absRoot, mapPath, g, profile); err != nil {
					fmt.Fprintf(os.Stderr, "aviso: falha ao registrar (carimbo/issue): %v\n", err)
				}
				// gates de JULGAMENTO: enfileira uma task `judge` por alvo pendente,
				// para uma IA confrontar e reportar com `anchors judge`.
				if n := enqueueJudgments(absRoot, cfg, profile); n > 0 {
					fmt.Printf("%d alvo(s) aguardam julgamento de IA — rode `anchors next` (ou `anchors judge --pending`)\n", n)
				}
			}

			if c := espelho.Caminho(); c != "" {
				rel, err := filepath.Rel(absRoot, c)
				if err != nil {
					rel = c
				}
				fmt.Printf("\nsaída espelhada em %s — releia daqui em vez de rodar de novo.\n", rel)
			}

			if !profile.Passed {
				// `os.Exit` não roda os `defer`: sem fechar aqui, o espelho perderia
				// o fim do relatório exatamente no caso em que ele mais importa — o
				// da reprovação.
				espelho.Fechar()
				os.Exit(1) // barra: há fail bloqueante
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().StringSliceVar(&changed, "changed", nil, "arquivo(s) alterado(s) — repetível ou separado por vírgula; os gates rodam UMA vez sobre a união dos caminhos de impacto")
	cmd.Flags().BoolVar(&all, "all", false, "varre todos os nós (a foto completa; caro)")
	cmd.Flags().BoolVar(&noRecord, "no-record", false, "só reporta: não carimba o mapa nem abre issues")
	cmd.Flags().BoolVar(&fix, "fix", false, "self-healer: aplica os reparos automáticos (ex.: corrige updated_at) antes de confrontar")
	cmd.Flags().BoolVar(&deterministic, "deterministic", false, "roda só os gates computáveis (pula os de julgamento por IA) — o modo do pre-commit")
	cmd.Flags().StringVar(&phase, "phase", "", "cobra só os gates desta fase (pre-commit|pre-push|ci|manual)")
	cmd.Flags().StringVar(&category, "category", "", "cobra só os gates desta natureza (types|style|traceability…)")
	cmd.Flags().BoolVar(&skipSlow, "skip-slow", false, "pula os gates declarados `cost: slow`")
	cmd.Flags().StringVar(&skipRegras, "skip-rule", "",
		"dispensa regras nesta execução: `id=motivo[,id=motivo]`. O motivo é obrigatório")
	cmd.Flags().StringVar(&msgPath, "commit-msg", "",
		"caminho do arquivo de mensagem de commit, de onde ler os marcadores `[skip-regra@CODIGO: motivo]`")
	cmd.Flags().BoolVar(&onlyIssues, "only-issues", false, "omite da tabela os gates que passaram em tudo e não deixaram pendência (o total continua no rodapé)")
	cmd.Flags().BoolVar(&showDrift, "show-drift", false, "lista TODAS as pendências (⚠) com o endereço de cada uma; sem a flag, só o contador da tabela")
	return cmd
}

// filtrarGates aplica a CATEGORIZAÇÃO: fase, natureza, custo e perspectiva. Cada eixo é
// independente — `when` diz em que momento o gate é cobrado, `cost` diz se cabe num
// loop apertado, `category` diz o que ele mede, `skip_on` diz sobre QUANTO do projeto a
// resposta dele tem valor.
//
// Todos os filtros são permissivos por omissão: gate sem `when`/`cost`/`category`/
// `skip_on` continua rodando como antes. É o que torna a categorização adotável aos
// poucos — declarar um eixo num gate não muda o comportamento dos outros 30.
func filtrarGates(gates []config.Gate, phase, category string, skipSlow bool, perspective string, dispensa gate.Dispensa) []config.Gate {
	out := gates[:0:0]
	for _, g := range gates {
		if !g.RunsIn(phase) {
			continue
		}
		if category != "" && g.Category != category {
			continue
		}
		if skipSlow && g.IsSlow() {
			continue
		}
		// A perspectiva NÃO é escolha do usuário: vem de COMO o check foi chamado
		// (`--all` ou `--changed`). O gate declara em qual delas ele se abstém.
		if g.SkipsOn(perspective) {
			continue
		}
		// DISPENSA por ID: o usuário declarou, com motivo, que esta regra não deve ser
		// confrontada nesta execução. Diferente de `skip_on`, que é decisão do gate sobre
		// perspectiva, esta é decisão de quem chama sobre uma regra específica — e é o que
		// permite commitar a primeira spec de uma unidade (a feature ainda é um card) sem
		// desligar os gates que verificam outra coisa.
		if motivo, ok := dispensa.Dispensou(gate.RegraID(idDoGateCmd(g))); ok {
			fmt.Printf("○ dispensado: %s — %s\n", idDoGateCmd(g), motivo)
			continue
		}
		out = append(out, g)
	}
	return out
}

// enqueueJudgments cria uma task `judge` por alvo que um gate de julgamento marcou
// como pendente (verdict Judge). A task carrega o gate e o guide, para o worker (a
// IA) saber o que ler e confrontar. Reusa a fila. Idempotente pelo dedup da fila.
func enqueueJudgments(root string, cfg *config.Config, p gate.Profile) int {
	// índice gate → (guide, ask) para enriquecer a task
	guideOf := map[string]string{}
	askOf := map[string]string{}
	gatesDeJulgamentoConhecidos = gatesDeJulgamentoConhecidos[:0]
	for _, g := range cfg.Gates {
		if g.IsJudgment() {
			guideOf[g.Name] = g.Guide
			askOf[g.Name] = g.Ask
			gatesDeJulgamentoConhecidos = append(gatesDeJulgamentoConhecidos, g.Name)
		}
	}
	n := 0
	for _, r := range p.Judged {
		reason := fmt.Sprintf("gate '%s'", r.Gate)
		if gd := guideOf[r.Gate]; gd != "" {
			reason += " — leia " + gd
		}
		if ask := askOf[r.Gate]; ask != "" {
			reason += " — pergunta: " + strings.TrimSpace(ask)
		}
		t := queue.Task{
			ID:            "judge-" + r.Gate + "-" + strings.ReplaceAll(relSlug(r.Target), "/", "-"),
			Changed:       r.Target,
			Kind:          "judgment",
			Origin:        "check",
			SuggestedNext: "judge",
			Reason:        reason,
			CreatedAt:     time.Now().Format(time.RFC3339),
		}
		if created, _ := queue.Enqueue(root, t); created {
			n++
		}
	}
	descartaJulgamentosObsoletos(root, cfg, p)
	return n
}

// descartaJulgamentosObsoletos tira da fila as tasks de julgamento cujo alvo o gate
// NÃO enfileirou nesta rodada.
//
// A fila é persistente e o conjunto de alvos aplicáveis não é: um gate que passa a
// declarar `requires`, um alvo que perde a marcação que o tornava aplicável, uma spec
// renomeada — em todos os casos a task velha sobrevive e a fila passa a mentir. Foi o
// que aconteceu no app de referência: o `no-test-prova-real` ganhou o filtro, o check passou a
// reportar 16 alvos, e `judge --pending` seguia listando 583 porque as tasks de dias
// antes continuavam ali.
//
// Só mexe em task de julgamento vinda do check (`kind: judgment`, `origin: check`), e
// só do gate que rodou nesta rodada — não toca trabalho de outra origem nem de gate
// que não foi cobrado agora.
func descartaJulgamentosObsoletos(root string, cfg *config.Config, p gate.Profile) {
	rodou := map[string]bool{}
	for _, g := range cfg.Gates {
		if g.IsJudgment() {
			rodou[g.Name] = true
		}
	}
	vivos := map[string]bool{} // "gate\x00alvo" que o check acabou de enfileirar
	for _, r := range p.Judged {
		vivos[r.Gate+"\x00"+r.Target] = true
	}
	tasks, err := queue.List(root)
	if err != nil {
		return
	}
	for _, t := range tasks {
		if t.Kind != "judgment" || t.Origin != "check" {
			continue
		}
		gateName := gateDaTaskJudge(t.ID)
		if gateName == "" || !rodou[gateName] {
			continue
		}
		if vivos[gateName+"\x00"+t.Changed] {
			continue
		}
		_ = queue.Drop(root, t.ID)
	}
}

// gateDaTaskJudge extrai o nome do gate do ID `judge-<gate>-<slug-do-alvo>`.
//
// O ID é montado com `-` como separador e o nome do gate também os contém
// (`no-test-prova-real`), então não dá para partir pelo separador: a leitura é por
// PREFIXO conhecido, comparando com os gates de julgamento declarados.
func gateDaTaskJudge(id string) string {
	const pref = "judge-"
	if !strings.HasPrefix(id, pref) {
		return ""
	}
	resto := id[len(pref):]
	for _, g := range gatesDeJulgamentoConhecidos {
		if strings.HasPrefix(resto, g+"-") {
			return g
		}
	}
	return ""
}

// gatesDeJulgamentoConhecidos é preenchido por enqueueJudgments a cada rodada — o
// nome do gate vem da config, não de uma lista fixa.
var gatesDeJulgamentoConhecidos []string

// relSlug reduz um caminho a algo usável em ID de task.
func relSlug(p string) string { return strings.TrimSuffix(p, filepath.Ext(p)) }

// recordCheck fecha o loop: carimba as arestas confrontadas no mapa (persistindo o
// veredito, o que destrava a detecção de stale) e abre uma issue de violation por
// fail bloqueante. É a passagem de "reporta" para "registra" (QUALITY §5).
func recordCheck(root, mapPath string, g *mapx.Graph, p gate.Profile) error {
	now := time.Now()
	day := now.Format("2006-01-02")

	// 1) CARIMBO — converte os vereditos por nó (do gate) para o mapa e grava.
	gv := p.NodeVerdicts()
	verdicts := make([]mapx.NodeVerdict, len(gv))
	for i, v := range gv {
		verdicts[i] = mapx.NodeVerdict{ID: v.ID, Failed: v.Failed}
	}
	stamped := g.StampEdges(verdicts, now.Format(time.RFC3339))
	if err := mapx.Save(g, mapPath); err != nil {
		return fmt.Errorf("salvar mapa carimbado: %w", err)
	}

	// 2) ISSUES — o loop completo, abre E fecha:
	//   • fail bloqueante → ABRE uma violation (se ainda não existe).
	//   • pass de um (gate, alvo) que tinha issue → RESOLVE (move para done/).
	// A chave estável (gate+alvo+kind) é o que liga o pass de hoje à issue de ontem.
	opened, resolved := 0, 0
	adiadas := 0
	for _, r := range p.Results {
		if r.Verdict == gate.Skip {
			continue // gate não se aplica — não mexe na issue
		}
		iss := issue.Issue{
			Kind: issue.Violation, Target: r.Target, Gate: r.Gate,
			Detail: r.Detail, Date: day,
		}
		// DÍVIDA ASSUMIDA vira issue em `future/`. Enquanto era só `Pending`, ela ficava
		// como uma linha no cabeçalho de um arquivo — visível para quem o abrisse, sem
		// estado, sem como ser paga, sem como vencer. O gate dizia que quem declara
		// afirma três coisas ("conhece o dever, ele vale, e QUANDO será pago"), e o
		// "quando" era prosa livre que nada confrontava.
		//
		// Materializada como issue, ela ganha o ciclo de vida que já existe: mover para
		// `todo/` quando chega a hora, e o próprio confronto a fecha quando o dever é
		// cumprido. `future/` e não `todo/` porque quem lê `todo/` pergunta "o que faço
		// AGORA" — afogar essa lista com o que só vence depois é o caminho mais curto
		// para ninguém mais olhar nenhuma das duas.
		if r.Verdict == gate.Pending {
			// DECISÃO EM ABERTO vira issue em `todo/`, e não em `future/`: ela só é
			// resolvida se alguém a VIR e a levar a quem decide. `future/` é o que vence
			// depois — a pergunta não vence, ela trava quem for implementar.
			//
			// É a issue que o usuário fecha: respondendo direto nela, ou pedindo à IA que
			// liste as perguntas abertas. Quando a resposta virar regra e o item sair da
			// spec, o `Pass` do próximo confronto a resolve sozinho, pelo mesmo caminho
			// que já fecha as violações.
			if r.Decisão {
				issDec := iss
				issDec.Kind = issue.Decision
				// DO USUÁRIO: a resposta não está no código, e o agente não a tem. Deixá-la
				// como dele faria o agente retentá-la para sempre — ou, pior, decidir por
				// conta própria, que é exatamente o que este gate existe para evitar.
				issDec.Dono = issue.DonoUsuário
				created, at, err := issue.Open(root, issDec)
				if err != nil {
					return fmt.Errorf("abrir issue de decisão: %w", err)
				}
				if created {
					opened++
				} else if at != issue.Todo {
					fmt.Printf("   (decisão de %s já em %s/)\n", r.Target, at)
				}
				continue
			}
			if !r.Divida {
				continue // indeterminado, não é dívida de ninguém
			}
			iss.Prazo = r.Prazo
			created, at, err := issue.OpenAt(root, iss, issue.Future)
			if err != nil {
				return fmt.Errorf("abrir issue de dívida: %w", err)
			}
			switch {
			case created:
				adiadas++
			case at == issue.Todo || at == issue.Doing:
				// Alguém puxou a dívida para o trabalho de agora. Dizer isso importa: a
				// declaração no header continua dizendo "depois", e o repositório já diz
				// "agora" — quem lê só o header não saberia.
				fmt.Printf("   (dívida de %s @ %s já está em %s/ — virou trabalho de agora)\n",
					r.Gate, r.Target, at)
			}
			continue
		}
		switch r.Verdict {
		case gate.Fail:
			if r.Blocking { // só fail BLOQUEANTE vira issue (informativo não barra nem registra)
				created, at, err := issue.Open(root, iss)
				if err != nil {
					return fmt.Errorf("abrir issue: %w", err)
				}
				if created {
					opened++
				} else if at != issue.Todo {
					fmt.Printf("   (issue de %s @ %s já em %s/)\n", r.Gate, r.Target, at)
				}
			}
		case gate.Pass:
			// passou: se havia issue aberta para este (gate, alvo), fecha-a.
			//
			// O KIND faz parte da chave, e o `iss` acima nasce como `violation`. Um gate
			// que abre issue de outro kind precisa fechá-la pelo kind com que abriu — sem
			// isto a issue de decisão nunca fechava: o gate voltava a ✓ e a pergunta
			// continuava em `todo/` para sempre, que é pior que não ter aberto.
			if d := iss; d.Gate == "open-questions-resolved" {
				d.Kind = issue.Decision
				if ok, err := issue.Resolve(root, d.Key()); err != nil {
					return fmt.Errorf("resolver issue de decisão: %w", err)
				} else if ok {
					resolved++
					fmt.Printf("   ✓ decisão resolvida: %s → issues/done/\n", r.Target)
				}
			}
			ok, err := issue.Resolve(root, iss.Key())
			if err != nil {
				return fmt.Errorf("resolver issue: %w", err)
			}
			if ok {
				resolved++
				fmt.Printf("   ✓ issue resolvida: %s @ %s → %s/done/\n", r.Gate, r.Target, issue.Dir)
			}
		}
	}

	fmt.Printf("\nregistrado: %d aresta(s) carimbada(s); %d issue(s) nova(s), %d resolvida(s) em %s/\n",
		stamped, opened, resolved, issue.Dir)
	if adiadas > 0 {
		fmt.Printf("  %d dívida(s) assumida(s) registrada(s) em %s/future/ — não é para agora, "+
			"e não some sozinha\n", adiadas, issue.Dir)
	}
	return nil
}

// ExitNaoRegido é o código de saída para "este caminho não é regido pela Estrutura".
// Não é uma reprovação — é a resposta "não tenho jurisdição sobre isto".
//
// Existe como CÓDIGO, e não como texto a ser grepado, porque o pre-commit precisa
// distinguir isto de uma reprovação real. O hook antes fazia `grep "não está no mapa"`
// na saída: um casamento de string frágil que passava a valer para as DUAS situações
// assim que elas compartilharam uma mensagem — e o furo entrou por aí.
const ExitNaoRegido = 3

// errNaoRegido sinaliza o caminho não-regido. Carrega o alvo para a mensagem, e é
// reconhecida em main() para virar ExitNaoRegido em vez do exit 1 genérico.
type errNaoRegido struct{ target string }

func (e errNaoRegido) Error() string {
	return fmt.Sprintf("%q não é regido pela Estrutura (não casa nenhuma camada do `layers:`) — nada a confrontar", e.target)
}

// selectNodes decide o conjunto de nós a confrontar: todos (--all) ou o caminho
// de impacto de um arquivo alterado (--changed, incremental).
func selectNodes(g *mapx.Graph, cfg *config.Config, all bool, changed []string, root string) (nodes []mapx.Node, scope string, err error) {
	if all {
		return g.Nodes, "--all", nil
	}
	if len(changed) == 0 {
		return nil, "", fmt.Errorf("informe --changed <arquivo> (incremental) ou --all (tudo)")
	}

	// UNIÃO dos caminhos de impacto de todos os arquivos, resolvida numa passagem só.
	// Antes o `--changed` aceitava um arquivo, e quem tinha muitos (o pre-commit) era
	// obrigado a chamar o binário N vezes — recarregando config e mapa a cada volta,
	// ~1,2s por arquivo. Pior que lento: os gates relacionais (feature-test-match,
	// trinca-completa) confrontam a UNIDADE, então rodavam repetidos sobre o mesmo
	// conjunto, uma vez por peça dela.
	vistos := map[string]bool{}
	var ordem []string
	var naoRegidos int
	for _, c := range changed {
		ids, err := impactoDe(g, cfg, c, root)
		if err != nil {
			var nr errNaoRegido
			if errors.As(err, &nr) {
				// Não-regido não contamina o lote: o Anchors só não tem jurisdição sobre
				// ele. Com TODOS não-regidos, o erro sobe (o chamador decide o exit 3).
				naoRegidos++
				continue
			}
			return nil, "", err
		}
		for _, id := range ids {
			if !vistos[id] {
				vistos[id] = true
				ordem = append(ordem, id)
			}
		}
	}
	if len(ordem) == 0 && naoRegidos == len(changed) {
		// Todos fora da jurisdição. A mensagem nomeia UM alvo e conta o resto: juntar
		// os caminhos numa string só produzia um erro ilegível de 8 KB num commit de 80
		// arquivos, e ainda dava a impressão de existir um arquivo com aquele nome.
		alvo := changed[0]
		if len(changed) > 1 {
			alvo = fmt.Sprintf("%s (e mais %d)", changed[0], len(changed)-1)
		}
		return nil, "", errNaoRegido{target: alvo}
	}
	for _, n := range g.Nodes {
		if vistos[n.ID] {
			nodes = append(nodes, n)
		}
	}
	if len(changed) == 1 {
		return nodes, "--changed " + changed[0], nil
	}
	return nodes, fmt.Sprintf("--changed (%d arquivos)", len(changed)), nil
}

// impactoDe resolve UM arquivo alterado nos ids de nó que ele arrasta (o alvo, o que
// propaga a partir dele, o que o valida, e as peças da mesma unidade).
func impactoDe(g *mapx.Graph, cfg *config.Config, changed, root string) ([]string, error) {
	target := relTo(root, changed)
	if !nodeExists(g, target) {
		if _, statErr := os.Stat(filepath.Join(root, target)); statErr != nil {
			return nil, fmt.Errorf("%q não existe no disco nem no mapa — confira o caminho", target)
		}
		// O arquivo EXISTE e não está no mapa. Duas situações COMPLETAMENTE distintas
		// se escondiam aqui sob uma mensagem só, e o pre-commit tratava ambas como
		// benignas — deixando passar exatamente o que o Anchors existe para barrar:
		// um arquivo REGIDO novo (hook, tela, service) commitado sem spec/feature/teste.
		// A distinção é a Estrutura: se o caminho casa uma camada do `layers:`, ele é
		// regido; senão, é um arquivo qualquer do repo (package.json, lockfile, CI).
		//
		// O IGNORE vem ANTES da classificação, e a ordem importa: `issues/` e `changes/`
		// casam a camada `doc` (`**/*.md`), mas o scanner NUNCA os indexa — são a saída
		// do próprio Anchors, não material a trabalhar. Classificar sem consultar o
		// ignore criava um impasse sem saída: o arquivo era dito "regido", não estava no
		// mapa, e `map build` não o acrescentava nunca — o commit ficava barrado para
		// sempre. Medido ao commitar as issues que o próprio `check` tinha resolvido.
		// `SkipDir` (e não `SkipFile`) porque é lá que mora `registroDoAnchors` — a
		// decisão é sobre o DIRETÓRIO-RAIZ do caminho (`issues/…`, `changes/…`).
		ig := scan.LoadIgnoreFor(root, cfg)
		if raiz, _, achou := strings.Cut(filepath.ToSlash(target), "/"); achou && ig.SkipDir(raiz, raiz) {
			return nil, errNaoRegido{target: target}
		}
		if ig.SkipFile(target) {
			return nil, errNaoRegido{target: target}
		}
		if layer, _ := scan.Classify(target, cfg); layer == "" {
			return nil, errNaoRegido{target: target}
		}
		return nil, fmt.Errorf("%q é REGIDO (camada do `layers:`) mas não está no mapa — "+
			"rode `anchors map build` primeiro (é o passo que registra o arquivo novo). "+
			"Enquanto ele estiver fora do mapa, NENHUM gate o confronta: a trinca não é "+
			"cobrada e o pipeline certifica trabalho que não foi verificado", target)
	}
	imp := g.AnalyzeImpact(target)
	// o alvo + o que propaga a partir dele (os filhos a refazer). A validação
	// (subir) é confronto contra os pais — entra como alvos também, para os gates
	// que se aplicam a eles rodarem.
	ids := map[string]bool{target: true}
	for _, n := range imp.Propagate {
		ids[n] = true
	}
	for _, n := range imp.Validate {
		ids[n] = true
	}
	// As PEÇAS DA MESMA UNIDADE entram junto, sempre. O impacto é calculado por arestas, e
	// spec/feature/teste/código de uma unidade nem sempre estão todos ligados por elas —
	// então `--changed <alvo>.ts` deixava de fora a `.feature` e o `.test.ts` irmãos.
	//
	// A consequência é um verde que não significa nada: medido, o `--changed` do código
	// não dispara `feature-test-match`, `feature-nao-vazia` nem `teste-nao-vazio`, e um
	// worker que siga o comando prescrito pelo `anchors work` declara "todos os
	// bloqueantes ✗0" sem jamais ter rodado o bloqueante mais importante da sua etapa.
	for _, peca := range pecasDaUnidade(target) {
		if nodeExists(g, peca) {
			ids[peca] = true
		}
	}
	out := make([]string, 0, len(ids))
	for id := range ids {
		out = append(out, id)
	}
	return out, nil
}

// larguraDoNome é a coluna do nome do gate: o maior nome presente, com um piso.
//
// Era `%-20s` fixo, e cinco gates passam disso (`handler-ddb-inline-passivo` tem
// 26). O nome mais longo empurrava a coluna do veredito e desalinhava a tabela
// inteira — numa lista de 49 linhas, o desalinhamento é o que faz o olho perder
// a coluna que importa.
func larguraDoNome(nomes []string) int {
	w := 18
	for _, n := range nomes {
		if len([]rune(n)) > w {
			w = len([]rune(n))
		}
	}
	return w
}

// larguras dos contadores: uma por COLUNA, pelo maior número daquela coluna.
//
// Uma largura única para todas desperdiça espaço onde ele não é preciso: numa
// varredura real o `~` chega a 582 e o `✗` fica em 0 ou 1, e a largura comum
// obrigaria a coluna dos fails a reservar três casas para nada. Cada coluna com
// a sua mantém os números alinhados à direita — que é o que permite compará-los
// a olho — sem esticar a tabela.
type largurasContador struct{ pass, fail, drift, skip, judge int }

func casas(n int) int { return len(fmt.Sprint(n)) }

func calcularLarguras(p gate.Profile) largurasContador {
	// `drift` nasce em 0 — e continua 0 se nenhum gate tiver drift, que é o
	// sinal para a coluna inteira não existir. Os outros têm piso 1: eles sempre
	// aparecem, e `%*d` com largura 0 imprimiria colado no símbolo.
	w := largurasContador{pass: 1, fail: 1, drift: 0, skip: 1, judge: 1}
	max := func(atual, n int) int {
		if c := casas(n); c > atual {
			return c
		}
		return atual
	}
	for nome, s := range p.ByGate {
		d := driftCount(p, nome)
		w.pass = max(w.pass, s.Pass)
		w.fail = max(w.fail, s.Fail)
		// `max` não serve para o drift: `casas(0)` é 1, e a coluna ganharia
		// largura mesmo sem drift nenhum — que é justamente o sinal de que ela
		// não deve existir. Só um drift REAL abre a coluna.
		if d > 0 {
			w.drift = max(w.drift, d)
		}
		w.skip = max(w.skip, s.Skip+s.Pending-d)
		w.judge = max(w.judge, s.Judge)
	}
	return w
}

// colunaDrift devolve a célula do ⚠. A decisão é da TABELA, não da linha:
//
//   - se algum gate tem drift, a coluna existe em TODAS as linhas — vazia vira
//     branco do mesmo tamanho. Omiti-la só nas linhas sem drift empurraria o `~`
//     para a esquerda nelas, e a mesma coluna passaria a existir em dois lugares.
//   - se NENHUM tem, a coluna não existe em lugar nenhum. Reservá-la aí deixaria
//     um buraco no meio de todas as linhas sem nada que o justificasse — e o
//     caso é o comum, não a exceção: no `--changed` a tabela costuma sair
//     inteira sem drift.
//
// `largura == 0` é o sinal de "a tabela não tem drift nenhum".
func colunaDrift(drift, largura int) string {
	if largura == 0 {
		return ""
	}
	if drift == 0 {
		// `len("⚠")` são 3 BYTES, mas o símbolo ocupa 1 coluna no terminal.
		// Usar `len` aqui reservava dois espaços a mais e desalinhava justamente
		// as linhas que o branco existia para alinhar.
		return strings.Repeat(" ", largura+1)
	}
	return fmt.Sprintf("⚠%*d", largura, drift)
}

// separador some junto com a coluna: sem isso, a tabela sem drift ficaria com
// dois espaços a mais entre `✗` e `~`.
func separadorDrift(largura int) string {
	if largura == 0 {
		return ""
	}
	return "  "
}

// gateLimpo: passou em tudo que olhou e não deixou nada pendente. É o gate que
// não pede nada de ninguém — o candidato a sumir sob `--only-issues`.
func gateLimpo(s gate.GateSummary, drift int) bool {
	return s.Fail == 0 && drift == 0 && s.Skip+s.Pending == 0 && s.Judge == 0
}

// printDrift lista as pendências agrupadas por GATE e, dentro dele, por MOTIVO.
//
// A lista crua tratava 2.430 pendências como 2.430 problemas, e elas não são:
// medido no app de referência, 832 delas têm o motivo IDÊNTICO ("sem sinal de mutação
// ingerido") e são um problema só — ninguém rodou a ingestão. Repetir o mesmo
// parágrafo 832 vezes esconde essa leitura em vez de revelá-la.
//
// Onde o motivo se repete, ele é escrito UMA vez e os alvos vêm compactados numa
// linha; onde cada motivo é único (`feature-test-match`, em que a divergência é
// específica de cada unidade), a lista continua alvo a alvo. O agrupamento não
// pode custar o endereço — ele é a parte acionável.
func printDrift(drifts []gate.Result) {
	if len(drifts) == 0 {
		return
	}
	// Agrupa preservando a ordem de aparição: gate → motivo → alvos.
	type grupo struct {
		motivo string
		alvos  []string
	}
	ordemGates := []string{}
	porGate := map[string][]*grupo{}
	total := map[string]int{}
	for _, r := range drifts {
		if _, visto := porGate[r.Gate]; !visto {
			ordemGates = append(ordemGates, r.Gate)
		}
		total[r.Gate]++
		var g *grupo
		for _, cand := range porGate[r.Gate] {
			if cand.motivo == r.Detail {
				g = cand
				break
			}
		}
		if g == nil {
			g = &grupo{motivo: r.Detail}
			porGate[r.Gate] = append(porGate[r.Gate], g)
		}
		g.alvos = append(g.alvos, r.Target)
	}

	fmt.Printf("\n⚠ %d pendência(s) em %d gate(s) — o gate olhou, algo divergiu, e não barrou:\n",
		len(drifts), len(ordemGates))
	for _, nome := range ordemGates {
		fmt.Printf("\n  %s — %d\n", nome, total[nome])
		for _, g := range porGate[nome] {
			if len(g.alvos) == 1 {
				fmt.Printf("    ⚠ %s\n", g.alvos[0])
				if g.motivo != "" {
					fmt.Print(indent(g.motivo, "        "))
				}
				continue
			}
			// Motivo repetido: escrito UMA vez, e os alvos listados um por LINHA.
			//
			// Juntá-los com vírgula economizava altura e custava legibilidade: no
			// projeto que originou isto, 583 alvos viraram uma linha de 36 MIL
			// caracteres — o endereço estava lá e ninguém conseguia lê-lo, nem
			// grepar por um arquivo específico. Uma linha por alvo é o formato que
			// as ferramentas de texto esperam.
			if g.motivo != "" {
				fmt.Print(indent(g.motivo, "    "))
			}
			fmt.Printf("      %d alvo(s):\n", len(g.alvos))
			for _, alvo := range g.alvos {
				fmt.Printf("        %s\n", alvo)
			}
		}
	}
}

// printLegenda explica os símbolos da tabela — e SÓ os que ela usou. Uma legenda
// fixa listaria `⚠` e `⏳` em varreduras que não os têm, e explicar coluna
// ausente é ruído que ensina a pular a legenda inteira.
//
// A distinção que ela precisa carregar é entre `✗` e `~`: falha é o gate tendo
// confrontado e divergido; indeterminado é ele não ter tido o que confrontar.
// Sem isso, `~582` parece um débito de 582 itens, quando é o contrário — é a
// medida de quanto daquele gate não se aplica ali.
func printLegenda(p gate.Profile, w largurasContador) {
	temJudge := false
	for _, s := range p.ByGate {
		if s.Judge > 0 {
			temJudge = true
			break
		}
	}
	partes := []string{
		"✓  passou",
		"✗  reprovou — vira issue",
	}
	if w.drift > 0 {
		partes = append(partes, "⚠  divergiu, mas não barra")
	}
	partes = append(partes, "~  indeterminado — o gate não teve o que confrontar")
	if temJudge {
		partes = append(partes, "⏳ aguarda julgamento de IA")
	}
	// Uma por linha. Em coluna única os cinco símbolos ficam empilhados e
	// comparáveis; na mesma linha a legenda passa de 100 colunas e quebra sozinha
	// num terminal estreito — o oposto do que ela serve.
	fmt.Println()
	for _, l := range partes {
		fmt.Printf("  %s\n", l)
	}
}

// avisarGatesSemAlvo relata os gates DECLARADOS que não apareceram na tabela.
//
// A tabela é montada de `profile.GateNames()`, que só conhece gate AVALIADO. Um gate
// declarado cujo `on:` não casa nenhum nó do mapa nunca é avaliado, então desaparece por
// completo — e o rodapé ainda diz "✓ pode promover, todos os gates passaram".
//
// Para quem está adotando o framework, isso lê como "está tudo certo". Medido ao pôr o
// Anchors no próprio Anchors: declarei 8 gates, o cabeçalho anunciou "8 gates", a tabela
// mostrou 2 e o rodapé deu ✓. Os 6 ausentes eram os de trinca, e o projeto tem 0 spec e 0
// feature — a informação útil ("declarei regra para artefato que não existe aqui") era
// exatamente a que não aparecia.
//
// Não é falha: é declaração sem alvo, e o conserto é do projeto (criar o artefato ou
// remover o gate). Por isso avisa, não barra.
func avisarGatesSemAlvo(declarados []config.Gate, p gate.Profile) {
	avaliados := map[string]bool{}
	for _, n := range p.GateNames() {
		avaliados[n] = true
	}
	var semAlvo []string
	for _, gt := range declarados {
		if gt.Name != "" && !avaliados[gt.Name] {
			semAlvo = append(semAlvo, gt.Name)
		}
	}
	if len(semAlvo) == 0 {
		return
	}
	fmt.Printf("\n⚠ %d gate(s) declarado(s) sem nada para medir — nenhum nó do mapa casa o `on:` deles:\n",
		len(semAlvo))
	for _, n := range semAlvo {
		fmt.Printf("    %s\n", n)
	}
	fmt.Println("  Não é falha: o gate está declarado e o artefato que ele mede não existe (ainda).")
	fmt.Println("  Crie o artefato, ou remova o gate — enquanto isso ele não protege nada.")
}

func printProfile(p gate.Profile, onlyIssues, showDrift bool) {
	nomes := p.GateNames()
	wn := larguraDoNome(nomes)
	w := calcularLarguras(p)
	limpos := 0

	// perfil por gate
	for _, name := range nomes {
		s := p.ByGate[name]
		tag := "informativo"
		if s.Blocking {
			tag = "bloqueante"
		}
		if s.Judge > 0 {
			tag = "julgamento"
			fmt.Printf("  %-*s  %-11s  ⏳%*d pendente(s) de IA\n", wn, name, tag, w.judge, s.Judge)
			continue
		}
		// O `~` fundia TRÊS coisas: "não se aplica", "não havia o que confrontar" e
		// DRIFT REAL rebaixado a aviso (ex.: o código do cenário casa mas a descrição
		// do teste divergiu). Somar tudo num número faz o drift parecer benigno — e é
		// o oposto: é a única categoria acionável do balde. Separamos em ⚠.
		drift := driftCount(p, name)
		if onlyIssues && gateLimpo(s, drift) {
			limpos++
			continue
		}
		// A coluna do ⚠ é reservada SEMPRE, mesmo vazia. Omiti-la nas linhas sem
		// drift empurrava o `~` para a esquerda só nelas, e a coluna passava a
		// existir em dois lugares na mesma tabela — que é o pior caso: o olho
		// desce a lista comparando números que não estão na mesma vertical.
		fmt.Printf("  %-*s  %-11s  ✓%*d  ✗%*d%s%s  ~%*d\n", wn, name, tag,
			w.pass, s.Pass, w.fail, s.Fail,
			separadorDrift(w.drift), colunaDrift(drift, w.drift),
			w.skip, s.Skip+s.Pending-drift)
	}
	// O gate omitido continua tendo rodado, e o número diz isso. Sem esta linha o
	// `--only-issues` pareceria uma varredura menor, e não a mesma varredura com a
	// listagem enxuta.
	if limpos > 0 {
		fmt.Printf("\n  + %d gate(s) sem nada a reportar (omitidos por --only-issues)\n", limpos)
	}
	printLegenda(p, w)

	// DRIFT detalhado: só sob `--show-drift`, e aí SEM TETO.
	//
	// Truncar a lista devolvia o problema que ela existe para resolver: quem pede
	// os endereços quer agir sobre eles, e "… e mais 2405" deixa 2405 sem
	// endereço. Ou a lista é completa, ou o contador da tabela já bastava.
	//
	// Por isso o padrão é o contador: numa varredura grande são milhares de
	// linhas, e despejá-las sem ninguém ter pedido enterraria as issues logo
	// abaixo. A tabela continua anunciando `⚠407`, e a flag é o caminho do número
	// para os endereços.
	if showDrift {
		printDrift(driftResults(p))
	}

	// Motivo dos `~` (skip/pending). O contador sozinho não diz nada — quem lê fica em
	// dúvida se `~1` é problema dele. No modo incremental (poucos nós) o motivo cabe na
	// tela e é acionável ("o código ainda não existe"), então mostramos.
	if reasons := skipReasons(p); len(reasons) > 0 && len(p.Results) <= maxResultsForSkipDetail {
		fmt.Printf("\n~ %d indeterminado(s) — não é falha; o gate não teve o que confrontar:\n", len(reasons))
		for _, r := range reasons {
			fmt.Printf("  ~ %s @ %s\n", r.Gate, r.Target)
			if r.Detail != "" {
				fmt.Print(indent(r.Detail, "      "))
			}
		}
	}

	if len(p.Failures) > 0 {
		fmt.Printf("\n%d issue(s) — divergências registradas:\n", len(p.Failures))
		for _, r := range p.Failures {
			mark := "informativo"
			if r.Blocking {
				mark = "BLOQUEIA"
			}
			fmt.Printf("  ✗ [%s] %s @ %s\n", mark, r.Gate, r.Target)
			if r.Detail != "" {
				fmt.Print(indent(r.Detail, "      "))
			}
		}
	}

	fmt.Println()
	// O ACHADO INFORMATIVO precisa sobreviver ao veredito.
	//
	// A última linha é a que se lê, e "✓ pode promover" apagava tudo que veio antes.
	// Medido em três rodadas de um E2E: um gate informativo diagnosticou um defeito real
	// (uma função que truncava a consulta e quebrava a regra que a consome), o autor leu
	// o verde, e registrou por escrito que o arquivo estava correto — contra o gate que
	// ele mesmo citou pelo nome. O sinal existia e sumiu no meio.
	//
	// Informativo não barra, e não deve barrar: é a maturação (QUALITY §7). Mas "não
	// barra" é diferente de "não aconteceu", e o veredito precisa dizer as duas coisas.
	informativos := 0
	// DRIFT é `Pending` COM motivo: o gate olhou, algo divergiu, e ele não barra. Contá-lo
	// aqui é o que evita a linha "sem achado em aberto" aparecer sobre uma spec que declara
	// 3 decisões que ninguém tomou — o veredito diria o contrário do que o gate viu.
	drift := 0
	// naoConfrontado são os Skip PUROS — o gate rodou e o nó não lhe dizia nada. Ficavam
	// fora de toda contagem, então um projeto onde NADA foi medido lia "sem achado em
	// aberto" no rodapé. Contá-los à parte é o que separa "conforme" de "não medido".
	naoConfrontado := 0
	for _, r := range p.Results {
		if r.Verdict == gate.Fail && !r.Blocking {
			informativos++
		}
		if r.Verdict == gate.Pending && r.Detail != "" {
			drift++
		}
		if r.Verdict == gate.Skip {
			naoConfrontado++
		}
	}
	if p.Passed {
		switch {
		case informativos > 0:
			fmt.Printf("✓ pode promover (bloqueantes passaram) — mas %d achado(s) INFORMATIVO(s) "+
				"em aberto\n", informativos)
			if drift > 0 {
				fmt.Printf("  (+%d pendência(s) — `anchors check` por unidade mostra quais)\n", drift)
			}
			fmt.Println("  Informativo não barra a promoção; também não desaparece. Leia antes de " +
				"seguir — é onde mora o que os bloqueantes não confrontam.")
		case drift > 0:
			fmt.Printf("✓ pode promover (nenhum gate reprovou) — mas %d PENDÊNCIA(s) em aberto\n", drift)
			fmt.Println("  Pendência não é defeito: é o gate dizendo que algo foi declarado e " +
				"ainda não decidido. Não barra, e não some sozinha.")
		case naoConfrontado > 0:
			// `Skip` PURO: o gate rodou, o nó não se aplicava, e nada foi medido. Não é
			// achado nem pendência — é COBERTURA AUSENTE, e o veredito não pode chamar
			// isso de "sem achado em aberto".
			//
			// Medido ao adotar o Anchors no próprio Anchors: 222 nós, dois gates com ~67 e
			// ~85, e o rodapé dizia "todos os gates passaram, sem achado em aberto". 152
			// confrontos não aconteceram, e a última linha — a que se lê — afirmava o
			// oposto. Para quem está adotando, é a diferença entre "meu projeto está
			// conforme" e "eu ainda não declarei o que fazia isto ser medido".
			fmt.Printf("\u2713 pode promover (nenhum gate reprovou) \u2014 mas %d confronto(s) N\u00c3O ACONTECERAM\n",
				naoConfrontado)
			fmt.Println("  O gate rodou e não teve o que medir naquele nó (falta o artefato do outro " +
				"lado da relação). Não é falha nem pendência: é cobertura que ainda não existe.")
		default:
			fmt.Println("✓ pode promover — todos os gates passaram, sem achado em aberto")
		}
	} else {
		fmt.Printf("✗ barrado — %d gate(s) bloqueante(s) reprovaram", len(p.Blocked))
		if informativos > 0 {
			fmt.Printf(" (+%d achado(s) informativo(s))", informativos)
		}
		fmt.Println()
	}

	lembraMaturacao(p, onlyIssues)
}

// lembraMaturacao avisa sobre gate informativo que já está LIMPO.
//
// A maturação (QUALITY §7) tem uma metade que o Anchors não cobrava: um gate nasce
// informativo porque o projeto ainda não cumpre o limiar, e quando passa a cumprir,
// ninguém volta ao anchors.yaml para promovê-lo. O gate fica medindo sem defender — e o
// projeto acha que está protegido por algo que não barra nada.
//
// A promoção continua sendo decisão humana; o que muda é que ela deixa de depender de
// alguém lembrar sozinho. E o lembrete aparece AQUI, onde a pessoa acabou de ler os
// vereditos — um aviso que exige rodar outro comando é um aviso que ninguém vê.
func lembraMaturacao(p gate.Profile, onlyIssues bool) {
	prom := gate.GatesPromoviveis(p)
	if len(prom) == 0 {
		return
	}
	// Com `--only-issues` a pessoa pediu para NÃO ver o que está limpo. Nomear os gates
	// aqui contradiria a flag — mas omiti-los por completo esconderia a maturação
	// pendente, que é justamente o que ela precisa saber. Então: conta, não lista.
	if onlyIssues {
		fmt.Printf("\n○ %d gate(s) informativo(s) limpo(s) — `anchors status` mostra quais\n", len(prom))
		return
	}
	fmt.Printf("\n○ %d gate(s) informativo(s) LIMPO(s) — prontos para virar bloqueantes:\n", len(prom))
	for i, g := range prom {
		if i == 6 {
			fmt.Printf("    … e mais %d\n", len(prom)-6)
			break
		}
		fmt.Printf("    %-26s %d nó(s) aprovado(s), 0 reprovado(s)\n", g.Gate, g.Passou)
	}
	fmt.Println("  Informativo mede e não defende. Enquanto ele não for `blocking: true` no")
	fmt.Println("  anchors.yaml, nada impede o próximo commit de desfazer o que já está conforme.")
}

// indent prefixa cada linha de s com pad — para o detalhe do gate (que pode ensinar
// o formato esperado em várias linhas) sair legível sob a issue.
func indent(s, pad string) string {
	var b strings.Builder
	for _, line := range strings.Split(strings.TrimRight(quebraOcorrencias(s), "\n"), "\n") {
		b.WriteString(pad)
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// quebraOcorrencias põe cada ocorrência do detalhe na sua própria linha.
//
// Dez gates montam a mensagem com `strings.Join(achados, "; ")`, e cinco
// violações no mesmo arquivo saíam numa linha de 800 caracteres: o leitor não
// distingue onde uma acaba e a outra começa, e a repetição do texto da regra
// (idêntico nas cinco) afoga o único dado que varia, que é o número da linha.
//
// A quebra é no separador que os gates JÁ usam, então nenhum deles precisa
// mudar. Basta UM separador para quebrar: duas ocorrências já se confundem numa
// linha só, e é o caso mais comum.
func quebraOcorrencias(s string) string {
	s = strings.ReplaceAll(s, "; ", ";\n")

	// Listas longas separadas por vírgula também quebram — 17 gates montam a
	// mensagem com `strings.Join(lista, ", ")` SEM cortar, e uma lista grande vira
	// uma linha que ninguém lê nem grepa.
	//
	// O que separa LISTA de PROSA não é o tamanho da linha: é o que vem entre as
	// vírgulas. Item de lista é um token — um caminho, um código, um símbolo — e
	// não tem espaço interno; oração tem. A primeira versão usava "linha longa com
	// 3+ vírgulas" e picou a própria mensagem deste gate ("a spec existe,\n o
	// código existe,\n os dois se referenciam") — a frase que o comentário citava
	// como o caso a NÃO quebrar.
	var out []string
	for _, linha := range strings.Split(s, "\n") {
		if ehListaDeTokens(linha) {
			linha = strings.ReplaceAll(linha, ", ", ",\n")
		}
		out = append(out, linha)
	}
	return strings.Join(out, "\n")
}

// limiarQuebraLista — abaixo disto a linha cabe na tela e não vale quebrar, seja
// ela lista ou não. Escolhido para caber num terminal de 120 colunas com a
// indentação do detalhe (6 a 8 espaços).
const limiarQuebraLista = 110

// ehListaDeTokens: a linha é uma enumeração de itens sem espaço interno?
//
// `a.spec.md, b.spec.md, c.spec.md` é lista; `a spec existe, o código existe, e
// os dois se referenciam` é prosa. A distinção é o espaço DENTRO do item: nome de
// arquivo, código de regra e símbolo não têm; oração tem. Exige maioria dos itens
// sem espaço para tolerar o último ("… e mais 3") e o texto que abre a lista.
func ehListaDeTokens(linha string) bool {
	if len([]rune(linha)) <= limiarQuebraLista {
		return false
	}
	partes := strings.Split(linha, ", ")
	if len(partes) < 4 {
		return false
	}
	tokens := 0
	for _, p := range partes {
		if !strings.Contains(strings.TrimSpace(p), " ") {
			tokens++
		}
	}
	return tokens*2 > len(partes)
}

// maxResultsForSkipDetail — acima disso (varredura --all) a lista de motivos vira
// ruído; o contador basta. No incremental o motivo é o que o usuário precisa.
const maxResultsForSkipDetail = 40

// driftResults devolve os DRIFTs: `Pending` COM motivo — o gate olhou, algo
// divergiu, e ele não barrou.
//
// Eles saíam do mesmo balde que os `~` e por isso desapareciam no `--all`, junto
// com o corte de ruído: a tabela anunciava `⚠407` e o detalhe abaixo não trazia
// nenhum. Contar uma categoria como acionável e depois escondê-la é pior que não
// a separar — o número vira uma acusação sem endereço.
func driftResults(p gate.Profile) []gate.Result {
	var out []gate.Result
	for _, r := range p.Results {
		if r.Verdict == gate.Pending && r.Detail != "" {
			out = append(out, r)
		}
	}
	return out
}

// skipReasons devolve os resultados indeterminados (skip/pending) QUE TÊM motivo
// escrito. Skip silencioso (gate simplesmente não se aplica àquele kind) não vira
// linha — só polui.
func skipReasons(p gate.Profile) []gate.Result {
	var out []gate.Result
	for _, r := range p.Results {
		// `Pending` COM motivo é DRIFT, e tem bloco próprio (`driftResults`).
		// Deixá-lo aqui o faria herdar o corte de ruído do `~` e sumir no --all.
		if r.Verdict == gate.Skip && r.Detail != "" {
			out = append(out, r)
		}
	}
	return out
}

// warnIfMapStale avisa quando o MAPA é mais velho que os arquivos que ele descreve.
// Sem isso, um arquivo criado/alterado depois do último `map build` é INVISÍVEL ao
// check — e o verde passa a não significar nada (o gate não roda sobre o que não está
// no mapa). É a diferença entre "passou" e "não foi olhado".
//
// Heurística barata e honesta: compara o mtime do anchors.graph.yaml com o mtime dos
// arquivos-fonte varridos. Não tenta ser exato (renames, deleções) — só levanta a mão.
func warnIfMapStale(root, mapPath string, cfg *config.Config) {
	mi, err := os.Stat(mapPath)
	if err != nil {
		return
	}
	mapTime := mi.ModTime()
	files, err := scan.Walk(root, cfg)
	if err != nil {
		return
	}
	newer := 0
	var sample string
	for _, f := range files {
		fi, err := os.Stat(filepath.Join(root, f.Path))
		if err != nil {
			continue
		}
		if fi.ModTime().After(mapTime) {
			newer++
			if sample == "" {
				sample = f.Path
			}
		}
	}
	if newer == 0 {
		return
	}
	extra := ""
	if newer > 1 {
		extra = fmt.Sprintf(" (e mais %d)", newer-1)
	}
	fmt.Printf("⚠ mapa DESATUALIZADO: %d arquivo(s) mudaram desde o último `map build` — ex.: %s%s\n"+
		"  Os gates só enxergam o que está no mapa; rode `anchors map build` antes de confiar neste resultado.\n\n",
		newer, sample, extra)
}

// driftCount conta, dentro do balde `~` de um gate, os resultados que são DRIFT REAL
// (indeterminados COM motivo escrito) — em oposição aos que apenas não se aplicam.
// Um Pending com detalhe é o gate dizendo "olhei e algo divergiu, mas não vou barrar";
// somá-lo aos "não se aplica" esconde a única parte acionável do número.
func driftCount(p gate.Profile, gateName string) int {
	n := 0
	for _, r := range p.Results {
		if r.Gate != gateName {
			continue
		}
		// DRIFT é PENDING com motivo: o gate OLHOU e algo divergiu, mas não barra.
		// Skip com motivo é o oposto — o gate não se aplica ("não é uma tela"), e
		// contá-lo como drift inflaria o número com ruído benigno.
		if r.Verdict == gate.Pending && r.Detail != "" {
			n++
		}
	}
	return n
}

// pecasDaUnidade devolve os caminhos das outras peças da trinca de um alvo — spec,
// feature, teste e código. É a vizinhança que a IDENTIDADE define, e não a que as arestas
// registram: as duas costumam coincidir, mas a segunda depende de o mapa já ter ligado as
// pontas, e o `check` precisa funcionar antes disso.
func pecasDaUnidade(target string) []string {
	base := strings.TrimSuffix(target, filepath.Ext(target))
	for _, suf := range []string{".spec.md", ".feature", ".test", ".spec"} {
		base = strings.TrimSuffix(base, suf)
	}
	var out []string
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".ts", ".tsx"} {
		if cand := base + suf; cand != target {
			out = append(out, cand)
		}
	}
	return out
}

// idDoGateCmd devolve o ID de um gate (ou o nome, quando o ID falta). Duplica a lógica
// de `internal/gate` de propósito: exportá-la só para isto acoplaria os dois pacotes por
// uma função de três linhas.
func idDoGateCmd(g config.Gate) string {
	if g.ID != "" {
		return g.ID
	}
	return g.Name
}

// normalizaAlterados põe os caminhos na mesma forma que os IDs do mapa (relativos à raiz),
// para que a comparação não dependa de como o chamador escreveu o caminho — o pre-commit
// passa relativo, e quem roda à mão costuma passar absoluto.
func normalizaAlterados(changed []string, root string) []string {
	out := make([]string, 0, len(changed))
	for _, c := range changed {
		p := c
		if filepath.IsAbs(p) {
			if rel, err := filepath.Rel(root, p); err == nil {
				p = rel
			}
		}
		out = append(out, filepath.ToSlash(filepath.Clean(p)))
	}
	return out
}

// Package gate executa os gates de qualidade sobre um conjunto de nós — o pipeline
// (QUALITY §5). Cada gate confronta os alvos que casam seu `on` e devolve um
// veredito. A dupla saída (issue + bloqueio) é derivada do veredito + do campo
// blocking (QUALITY §2, §7). O `check` orquestra; este pacote é a mecânica pura +
// a invocação de comandos externos.
package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// Verdict é o resultado de rodar um gate contra um alvo.
type Verdict string

const (
	Pass    Verdict = "pass"    // passou
	Fail    Verdict = "fail"    // reprovou → issue (bloqueia se o gate é blocking)
	Skip    Verdict = "skip"    // não se aplica a este alvo
	Pending Verdict = "pending" // gate interno ainda não implementado / indeterminado
	Judge   Verdict = "judge"   // gate de JULGAMENTO POR IA: aguarda veredito de uma IA
)

// Result é o veredito de um gate sobre UM alvo.
type Result struct {
	Gate string // nome do gate
	// Regra identifica QUAL verificação dentro do gate produziu este veredito.
	//
	// Um gate cobra mais de uma coisa (`spec-completa` cobra "sem placeholder" E "tem
	// regra catalogada"), e sem este campo dois defeitos diferentes eram indistinguíveis
	// — o leitor tinha de inferir pela mensagem, que muda. É também o que permite
	// dispensar UMA verificação sem descartar o resto do gate.
	//
	// Vazio quando o gate faz uma verificação só, e aí o nome do gate já a identifica.
	Regra    string
	Target   string // nó confrontado (ID)
	Verdict  Verdict
	Blocking bool   // se este gate bloqueia (da config)
	Detail   string // mensagem (por que falhou, etc.)
	// Divida marca o Pending que é DÍVIDA ASSUMIDA — um dever conhecido, ainda válido,
	// com um momento declarado para ser pago (`obligation_pending: <nome> — <quando>`).
	//
	// É um campo, e não uma inferência sobre o texto de Detail, porque a distinção
	// decide o DESTINO do achado: dívida vira issue em `future/` (trabalho adiado, que
	// vence), e os demais Pending não viram issue alguma. Deduzir isso de uma frase é
	// exatamente o tipo de acoplamento a texto que envelhece na primeira reescrita da
	// mensagem.
	Divida bool
	// Decisão marca o Pending que é DECISÃO EM ABERTO — a spec declarou que não decidiu
	// algo de que o código precisa.
	//
	// Campo próprio, e não um caso de `Divida`: dívida tem dono e vencimento ("sei o que
	// devo e quando pago"), e a decisão em aberto é o oposto — não se sabe a resposta,
	// não se sabe quando, e depende de outra pessoa. Os destinos também diferem: dívida
	// vai para `future/` (o que vence depois), e a pergunta vai para `todo/`, porque só
	// é resolvida se alguém a VIR e a levar a quem decide.
	Decisão bool
	// Prazo é só o "quando" de cada dívida — o que vem depois de "DÍVIDA ASSUMIDA:".
	// Separado de Detail porque quem lê a issue quer saber QUANDO ela vence sem
	// reprocessar o laudo inteiro do gate.
	Prazo string
	// Impede marca a PENDÊNCIA que impede a promoção — "há decisão por tomar", não "não
	// tive o que confrontar". Os dois casos usam o veredito `Pending`, e só o primeiro
	// deveria barrar: medido num repositório real, tratar todos igual reprovou 411 nós
	// de uma vez (410 eram gates sem sinal ingerido).
	//
	// Quem sabe a diferença é o gate, que sabe o que mediu.
	Impede bool
}

// applies: o gate se aplica a um nó se o kind está no `on` do gate, E o nó não carrega
// nenhuma tag de `exclude_tags`, E (se o gate declara tags) carrega ao menos uma delas,
// E (se declara `requires`) o conteúdo do alvo contém aquele texto. Sem nada disso → só
// o kind conta.
func applies(g config.Gate, n mapx.Node, root string) bool {
	if !slices.Contains(g.On, string(n.Kind)) {
		return false
	}
	// A exclusão vem ANTES do filtro positivo, e a ordem é a decisão: um nó excluído fica
	// fora mesmo que também case uma tag do `tags`. Camadas costumam carregar rótulos
	// transversais (`backend`, `code`) junto com o seu próprio, então a exceção precisa
	// vencer — senão declarar a exclusão não teria efeito nenhum onde ela importa.
	// A exclusão vem ANTES do filtro positivo, e a ordem é a decisão: um nó excluído fica
	// fora mesmo que também case uma tag do `tags`. Camadas costumam carregar rótulos
	// transversais (`backend`, `code`) junto com o seu próprio, então a exceção precisa
	// vencer — senão declarar a exclusão não teria efeito nenhum onde ela importa.
	for _, t := range g.ExcludeTags {
		if slices.Contains(n.Tags, t) {
			return false
		}
	}
	if len(g.Tags) > 0 {
		achou := false
		for _, t := range g.Tags {
			if slices.Contains(n.Tags, t) {
				achou = true
				break
			}
		}
		if !achou {
			return false
		}
	}
	if g.Requires == "" {
		return true
	}
	// `requires` filtra por CONTEÚDO: o gate que pergunta sobre uma marcação só se
	// aplica a quem a usou. Sem isto, um gate de julgamento declarado `on: [spec]`
	// enfileira uma pergunta de IA para toda spec do projeto — e o contador de
	// pendências passa a medir o tamanho do projeto, não o do trabalho.
	//
	// Arquivo ilegível NÃO aplica: melhor deixar de cobrar do que cobrar às cegas de
	// um alvo cujo conteúdo não se conhece.
	b, err := os.ReadFile(filepath.Join(root, n.ID))
	if err != nil {
		return false
	}
	return strings.Contains(string(b), g.Requires)
}

// Run roda todos os gates aplicáveis contra os nós dados. root é a raiz do projeto
// (para comandos externos e leitura de arquivos internos). Devolve um Result por
// (gate, alvo aplicável).
// Run confronta os gates. `graph` é o mapa (arestas + nós) — passado aos checkers
// RELACIONAIS (que confrontam um nó contra seus vizinhos, ex.: feature↔test). Pode ser
// nil quando não há mapa carregado; checkers relacionais devem tratar nil como Pending.
func Run(gates []config.Gate, nodes []mapx.Node, root string, graph *mapx.Graph) []Result {
	return RunWithConfig(gates, nodes, root, graph, nil)
}

// RunWithConfig é o Run que também recebe a config completa — necessária aos checkers
// relacionais que consultam a Estrutura (de-para de regimes, superfícies da trinca).
// `Run` delega a ela com cfg nil (checkers relacionais tratam nil como sem-de-para).
func RunWithConfig(gates []config.Gate, nodes []mapx.Node, root string, graph *mapx.Graph, cfg *config.Config) []Result {
	return RunCompleto(gates, nodes, root, graph, cfg, false)
}

// RunCompleto é o Run que também sabe se a varredura é o PROJETO INTEIRO (`check --all`).
// Só isso permite honrar o `scope_full` do gate: no full, quem sabe varrer sozinho roda
// UMA vez sem receber a lista, em vez de receber os milhares de alvos em lotes. Os demais
// chamadores seguem por RunWithConfig, que passa `completa: false` — o comportamento de
// sempre para recorte incremental.
func RunCompleto(gates []config.Gate, nodes []mapx.Node, root string, graph *mapx.Graph, cfg *config.Config, completa bool) []Result {
	return RunComDispensa(gates, nodes, root, graph, cfg, completa, Dispensa{})
}

// RunComDispensa é o Run que honra dispensa POR ALVO.
//
// A dispensa por regra era aplicada FILTRANDO o gate da lista, e isso bastava enquanto
// ela valia para tudo. Uma dispensa restrita a caminhos não pode sair por ali: o gate
// precisa RODAR e confrontar os outros alvos — senão dispensar 4 specs novas apagaria o
// gate para o repositório inteiro, e uma trinca quebrada por descuido passaria junto.
func RunComDispensa(gates []config.Gate, nodes []mapx.Node, root string, graph *mapx.Graph, cfg *config.Config, completa bool, disp Dispensa) []Result {
	// A gramática do código de cenário segue o vocabulário do projeto (`rule_types`).
	SetRuleLetters(cfg.RuleLetters())
	// índice kind por nó já vem em node.Kind
	var results []Result
	for _, g := range gates {
		// Os alvos deste gate — o mesmo filtro `on`/`tags` para qualquer escopo.
		var alvos []mapx.Node
		for _, n := range nodes {
			if applies(g, n, root) {
				alvos = append(alvos, n)
			}
		}
		if len(alvos) == 0 {
			// Nenhum alvo relevante: o gate não roda. Vale inclusive para `scope: project`
			// — um commit só de README não dispara o typecheck do monorepo inteiro.
			continue
		}
		// Ferramenta exigida e ausente: Skip, nunca Fail. O gate não mediu — e "não medi"
		// não é "está limpo" nem "está sujo". Um Fail aqui diria que o projeto violou algo,
		// quando quem falta é o binário. O aviso não se perde: o `doctor` levanta a ausência
		// (health.checkFerramentasAusentes), que é onde ela pode ser lida uma vez e resolvida,
		// em vez de repetida a cada alvo de cada varredura.
		if faltando, ok := ferramentaAusente(g); ok {
			results = append(results, Result{
				Gate: g.Name, Target: "(" + string(g.ScopeParaVarredura(completa)) + ")",
				Verdict: Skip, Blocking: g.IsBlocking(),
				Detail: "ferramenta ausente: " + faltando + " — gate não executado",
			})
			continue
		}
		switch g.ScopeParaVarredura(completa) {
		case config.ScopeBatch, config.ScopeProject:
			results = append(results, runAgregado(g, alvos, root, completa, graph, cfg))
		default:
			for _, n := range alvos {
				// DISPENSA POR ALVO: o gate roda, e só este nó é poupado. O veredito é
				// `Skip` com o motivo escrito — some do placar de reprovações sem sumir
				// do relatório, que é a diferença entre dispensar e esconder.
				if motivo, ok := disp.DispensouAlvo(RegraID(idDoGate(g)), n.Code); ok {
					results = append(results, Result{
						Gate: g.Name, Regra: idDoGate(g), Target: n.ID,
						Verdict: Skip, Blocking: g.IsBlocking(),
						Detail: "dispensado: " + motivo,
					})
					continue
				}
				results = append(results, runOne(g, n, root, graph, cfg))
			}
		}
	}
	return results
}

// runAgregado executa UMA vez um gate de escopo batch/project. O veredito é único:
// a ferramenta olhou o conjunto e respondeu sobre ele.
//
// O alvo reportado é o próprio gate (não um arquivo): atribuir a falha do `tsc` a um
// dos 63 arquivos seria mentira — o erro pode estar em qualquer um, ou na relação
// entre eles. O laudo (stdout da ferramenta) é que nomeia arquivo e linha.
func runAgregado(g config.Gate, alvos []mapx.Node, root string, completa bool, graph *mapx.Graph, cfg *config.Config) Result {
	escopo := g.ScopeParaVarredura(completa)
	r := Result{Gate: g.Name, Regra: idDoGate(g), Target: "(" + escopo + ")", Blocking: g.IsBlocking()}
	// Um gate agregado pode ser INTERNO: a pergunta é sobre o conjunto, mas quem
	// responde é o próprio CLI, não uma ferramenta de fora. É o caso de
	// `testid-consultado-existe`, que confronta as duas pontas do projeto de uma vez —
	// exigir `run:` dele obrigaria a empacotar lógica de Go num script externo só para
	// caber na assinatura.
	//
	// O alvo não é um nó (o escopo é o conjunto), então passamos um nó VAZIO: os
	// checkers agregados leem `root`/`cfg`, nunca o conteúdo de um arquivo — e é por
	// isso que este caminho não tenta lê-lo, ao contrário do `runInternal`.
	if g.Check != "" {
		r.Verdict, r.Detail = runInternalAgregado(g, root, graph, cfg)
		return r
	}
	if g.Run == "" {
		r.Verdict, r.Detail = Pending, "gate de escopo "+escopo+" exige `run:` (comando externo) ou `check:` (interno)"
		return r
	}
	var args []string
	if escopo == config.ScopeBatch {
		// batch: a ferramenta recebe os arquivos que casaram ("$@").
		args = make([]string, 0, len(alvos))
		for _, n := range alvos {
			args = append(args, n.ID)
		}
	}
	// project: sem argumentos — a ferramenta já sabe o que olhar.
	r.Verdict, r.Detail = RunExternalArgs(g.Run, args, root)
	return r
}

// runOne executa um gate contra um alvo — despacha para interno ou externo.
func runOne(g config.Gate, n mapx.Node, root string, graph *mapx.Graph, cfg *config.Config) Result {
	r := Result{Gate: g.Name, Regra: idDoGate(g), Target: n.ID, Blocking: g.IsBlocking()}
	switch {
	case g.IsJudgment():
		// gate de julgamento por IA: o CLI NÃO computa. Ele só sabe se ALGUÉM já
		// respondeu — e essa resposta está no carimbo que o `anchors judge` deixou.
		//
		// Sem esta consulta o gate pedia julgamento para sempre: o `judge` gravava o
		// veredito no grafo e o `check` seguinte emitia Judge de novo, então o contador
		// de pendências nunca descia e o trabalho já feito era invisível. O carimbo leva
		// as revs das pontas, então um veredito envelhece se o alvo mudar depois — e
		// nesse caso volta a ser pergunta, que é o comportamento certo.
		if graph != nil {
			if v, ok := graph.JulgadoPor(n.ID, g.Name); ok {
				switch v {
				case "issue":
					r.Verdict = Fail
					r.Detail = "julgado por IA: reprovado — ver a issue do laudo"
				case "ok":
					r.Verdict = Pass
					r.Detail = "julgado por IA: aprovado"
				case "dispensado":
					// `Skip`, e não `Pass`: o gate NÃO MEDIU — o alvo da pergunta não
					// existe (a spec o declara `@TBD`). `Pass` afirmaria aprovação, que é
					// a mentira que o veredito `dispensado` existe para evitar; `Pending`
					// o classificaria como divergência, que também não é o caso.
					//
					// `Skip` é a categoria certa e já existe: "não se aplica a este alvo".
					r.Verdict = Skip
					r.Detail = "julgado por IA: DISPENSADO — o alvo da pergunta não existe " +
						"(peça declarada `@TBD`)"
				default:
					r.Verdict = Pending
					r.Detail = "julgado por IA com veredito " + v
				}
				return r
			}
		}
		r.Verdict = Judge
		r.Detail = "aguarda julgamento de IA"
		if g.Guide != "" {
			r.Detail += " (guide: " + g.Guide + ")"
		}
	case g.Check != "":
		r.Verdict, r.Detail = runInternal(g.Check, n, root, graph, cfg)
		// DÍVIDA ASSUMIDA é o Pending que tem dono e vencimento — e o único que vira
		// trabalho registrado (issue em `future/`). Só o gate de obrigações a produz; os
		// demais Pending são "não tive o que confrontar", que não é dívida de ninguém.
		if r.Verdict == Pending && g.Check == "obligation-honored" {
			r.Divida = true
			r.Prazo = prazosDeclarados(r.Detail)
		}
		// Decisão em aberto IMPEDE: a spec declara que não decidiu algo que o código vai
		// precisar, e quem implementar vai adivinhar. Dívida assumida NÃO impede — ela
		// diz "sei o que devo e quando pago", que é o oposto de não saber.
		// Só a pendência que diz "há decisão POR TOMAR" impede a promoção. A que diz "a
		// seção nem existe" é dívida de migração — 586 specs do repositório nasceram antes
		// da prática, e barrá-las reprovaria o projeto inteiro por algo que ninguém
		// escreveu errado.
		//
		// A distinção sai do próprio gate, pelo MARCADOR que ele põe no laudo. Antes era
		// por prosa ("que a spec ainda NÃO tomou"), e isso quebraria na tradução do
		// laudo — sem erro, sem aviso: o gate simplesmente pararia de barrar.
		if r.Verdict == Pending && g.Check == "open-questions-resolved" &&
			strings.Contains(r.Detail, MarcaDecisaoEmAberto) {
			r.Impede = true
			// E vira ISSUE. A mesma distinção decide as duas coisas: "há decisão por
			// tomar" é achado que precisa sobreviver à sessão; "a spec nasceu antes da
			// prática" é dívida de migração, e abrir issue para as 586 specs que nunca
			// tiveram a seção afogaria `todo/` com o que ninguém escreveu errado.
			r.Decisão = true
		}
	case g.Run != "":
		r.Verdict, r.Detail = runExternal(g.Run, n, root)
	default:
		r.Verdict, r.Detail = Pending, "gate sem `run` nem `check` declarado"
	}
	return r
}

// prazoRE captura o "quando" de cada dívida na mensagem do gate de obrigações, que as
// concatena com ";". O nome da obrigação vem entre colchetes no início de cada trecho.
var prazoRE = regexp.MustCompile(`\[([a-z0-9-]+)\][^;]*?DÍVIDA ASSUMIDA: ([^;]+)`)

// prazosDeclarados extrai, do laudo do gate, apenas os vencimentos — um por obrigação.
func prazosDeclarados(detail string) string {
	var out []string
	for _, m := range prazoRE.FindAllStringSubmatch(detail, -1) {
		out = append(out, "`"+m[1]+"` — "+strings.TrimSpace(m[2]))
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, "\n- ")
}

// ferramentaAusente diz se o gate exige um binário que não está no PATH.
//
// `exec.LookPath` é a mesma resolução que o `sh` faria ao rodar o comando, então a
// resposta aqui e o comportamento real do gate não podem divergir.
func ferramentaAusente(g config.Gate) (string, bool) {
	if g.NeedsTool == "" {
		return "", false
	}
	if _, err := exec.LookPath(g.NeedsTool); err != nil {
		return g.NeedsTool, true
	}
	return "", false
}

// idDoGate devolve o ID declarado, ou o nome quando ele falta.
//
// A ausência é tolerada por compatibilidade — um projeto que já funcionava não pode
// quebrar por um campo novo —, e o nome do gate é um identificador razoável enquanto o ID
// não foi declarado. A validação da carga garante que, de um jeito ou de outro, ele é
// único.
func idDoGate(g config.Gate) string {
	if g.ID != "" {
		return g.ID
	}
	return g.Name
}

// Package issue materializa as issues do Anchors (CONCEPT §5): o registro de uma
// divergência que sobrevive à sessão. Uma issue é um arquivo markdown; seu ESTADO é
// a pasta em que vive (todo/doing/done — imutável, mover = mudar de estado). Três
// KINDS, cada um de uma origem: stale (uma ponta avançou), conflict (âncoras
// bloqueantes discordam), violation (o confronto rodou e o alvo viola a âncora).
//
// As issues vivem em issues/ (conteúdo de PROJETO, versionável), não em .anchors/
// (estado efêmero do daemon).
package issue

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Dir é a raiz das issues, relativa à raiz do projeto.
const Dir = "issues"

// Kind — a origem da issue (CONCEPT §5).
type Kind string

const (
	Stale     Kind = "stale"     // uma ponta da aresta avançou de rev
	Conflict  Kind = "conflict"  // âncoras bloqueantes discordam do mesmo alvo
	Violation Kind = "violation" // o confronto rodou e o alvo viola a âncora
	// Decision é a DECISÃO EM ABERTO: a spec declara que não decidiu algo que o código
	// vai precisar. Não é violação — ninguém escreveu nada errado —, e não é dívida
	// assumida, que tem dono e vencimento ("sei o que devo e quando pago"). É o oposto:
	// alguém PRECISA decidir, não se sabe quando, e depende de outra pessoa.
	//
	// Vira issue porque é o achado que mais precisa sobreviver à sessão: depende de humano
	// e pode levar semanas. Enquanto era só uma linha no `check`, a pergunta existia
	// apenas para quem abrisse a spec.
	Decision Kind = "decision"
)

// State — o ciclo de vida, codificado na pasta.
type State string

const (
	// Future é o que se DEVE e ainda não se paga: dívida assumida, trabalho
	// deliberadamente adiado. Não é `todo` — quem olha `todo/` está perguntando "o que
	// faço agora", e afogar essa lista com o que só vence depois é o caminho mais curto
	// para ninguém mais olhar. Não é `done`, porque não foi feito.
	//
	// Nasce da distinção que o `obligation_pending` já fazia no header e que nada
	// materializava: um dever conhecido, ainda válido, com um momento declarado para ser
	// cumprido. Ele vivia como uma linha no cabeçalho de um arquivo — visível só para
	// quem o abrisse, sem estado, sem como ser pago, sem como vencer.
	Future State = "future"
	Todo   State = "todo"  // detectada, ninguém pegou
	Doing  State = "doing" // alguém está resolvendo
	Done   State = "done"  // tratada (fato datado)
)

// Estados é a ordem de leitura do ciclo de vida — usada na busca e nos relatórios.
var Estados = []State{Future, Todo, Doing, Done}

// Issue é uma divergência registrada.
type Issue struct {
	Kind Kind
	// A aresta/alvo em confronto. Anchor é a âncora (a régua); Target é o regido.
	// Para uma violation de gate: Anchor pode ser vazio (o gate É a régra) e Target é
	// o nó reprovado.
	Anchor string
	Target string
	Gate   string // o gate que reprovou (quando kind=violation vinda de gate)
	Detail string // o corpo: por que, qual invariante quebrou
	Date   string // AAAA-MM-DD — carimbado por quem cria (evita Date.now interno)
	// Prazo é o "quando" de uma DÍVIDA ASSUMIDA (`obligation_pending: <nome> — <quando>`).
	// Preenchido, a issue nasce em `future/` e o corpo diz quando ela vence.
	Prazo string
	// Dono é DE QUEM é a issue: do agente, que a resolve confrontando código, ou do
	// USUÁRIO, que precisa decidir ou responder algo que o agente não pode.
	//
	// É um eixo INDEPENDENTE do kind, e por isso um campo e não um kind novo. Uma
	// `violation` costuma ser do agente e uma `decision` é sempre do usuário, mas a
	// correspondência não é fixa: o agente que tenta resolver uma violação e esbarra numa
	// pergunta que só o usuário responde precisa PASSAR a issue adiante, sem que ela deixe
	// de ser a violação que era.
	//
	// Filtrável (`anchors issues --dono usuario`): é a lista que se leva para a conversa
	// com quem decide, e misturá-la com o trabalho do agente faz as duas serem ignoradas.
	Dono Dono
}

// Dono diz quem consegue resolver a issue.
type Dono string

const (
	// DonoAgente é o padrão: a issue se resolve mexendo no repositório, e o próximo
	// confronto a fecha sozinho.
	DonoAgente Dono = "agente"
	// DonoUsuário é a issue que ESPERA UMA PESSOA: uma decisão de produto, uma resposta
	// que não está em lugar nenhum do código. O agente não a resolve por mais que tente —
	// e tentar é o defeito, porque decidir por conta própria é o que a decisão em aberto
	// existe para evitar.
	DonoUsuário Dono = "usuario"
)

// DonoDe devolve o dono declarado, com o padrão aplicado. Uma issue antiga, gravada antes
// do campo existir, não tem o valor no disco — e é do agente, que era o único caso.
func DonoDe(d Dono) Dono {
	if d == "" {
		return DonoAgente
	}
	return d
}

// Key é a IDENTIDADE ESTÁVEL da issue no tempo: <kind>--<aresta>, sem data nem hash.
// É por ela que uma issue é encontrada — para não duplicar entre dias e para poder
// ser RESOLVIDA quando o confronto volta a passar. Duas detecções do mesmo problema
// (mesma aresta+kind), em dias diferentes, têm a MESMA Key.
func (i Issue) Key() string {
	edge := slug(i.Target)
	if i.Anchor != "" {
		edge = slug(i.Anchor) + "--vs--" + slug(i.Target)
	}
	// o gate faz parte da identidade de uma violation (o MESMO alvo pode violar dois
	// gates distintos — são duas issues).
	if i.Gate != "" {
		edge = slug(i.Gate) + "--" + edge
	}
	return fmt.Sprintf("%s--%s", i.Kind, edge)
}

// ID é o nome do arquivo (sem pasta): <data>--<key>.md — legível e único, com a data
// só para leitura humana. A BUSCA é sempre pela Key (via byKey), não pelo ID inteiro,
// porque a data varia entre detecções.
func (i Issue) ID() string {
	return fmt.Sprintf("%s--%s.md", i.Date, i.Key())
}

func slug(path string) string {
	s := strings.ReplaceAll(path, "/", "-")
	return strings.ReplaceAll(s, string(filepath.Separator), "-")
}

// Body renderiza o markdown da issue — o template por kind (CONCEPT §5).
func (i Issue) Body() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: %s\n\n", strings.ToUpper(string(i.Kind)), i.Target)
	fmt.Fprintf(&b, "- **kind:** %s\n", i.Kind)
	if i.Anchor != "" {
		fmt.Fprintf(&b, "- **âncora (régua):** %s\n", i.Anchor)
	}
	fmt.Fprintf(&b, "- **alvo (regido):** %s\n", i.Target)
	if i.Gate != "" {
		fmt.Fprintf(&b, "- **gate:** %s\n", i.Gate)
	}
	// O DONO no cabeçalho é o que permite filtrar sem abrir cada arquivo — e é a
	// pergunta que quem lê a pasta faz primeiro: "o que aqui é meu?".
	fmt.Fprintf(&b, "- **dono:** %s\n", DonoDe(i.Dono))
	fmt.Fprintf(&b, "- **detectada em:** %s\n\n", i.Date)
	// Cabeçalho da seção do corpo — omitido quando o Detail já traz seus próprios
	// cabeçalhos markdown (um laudo estruturado da IA), para não duplicar.
	detailHasHeadings := strings.Contains(i.Detail, "\n#") || strings.HasPrefix(i.Detail, "#")
	if !detailHasHeadings {
		switch i.Kind {
		case Violation:
			b.WriteString("## Invariante violada\n\n")
		case Stale:
			b.WriteString("## Quem está atrás de quem\n\n")
		case Conflict:
			b.WriteString("## A discordância\n\n")
		case Decision:
			b.WriteString("## A pergunta em aberto\n\n")
		}
	}
	if i.Detail != "" {
		b.WriteString(i.Detail)
		b.WriteString("\n")
	}
	b.WriteString("\n---\n")
	if i.Kind == Decision {
		// O RODAPÉ diz como fechar. É a issue que o usuário resolve — sozinho, ou pedindo
		// à IA que liste as perguntas abertas —, e sem isto ele teria a pergunta sem o
		// procedimento: o caminho errado (apagar o item) é mais curto que o certo.
		b.WriteString("_Decisão em aberto, registrada pelo Anchors. Não é defeito: é a spec " +
			"dizendo a verdade sobre o que ainda não sabe.\n\n" +
			"**Como fechar:** leve a pergunta a quem decide e PROMOVA a resposta a regra " +
			"(com código) na spec. Marque o item como resolvido (`[x]`) citando a regra que " +
			"nasceu dele — o `anchors check` fecha esta issue sozinho no próximo confronto.\n\n" +
			"**O que NÃO fazer:** apagar o item sem regra nova. O gate não distingue isso de " +
			"ter decidido, e quem apaga assume a decisão silenciosamente._\n")
		return b.String()
	}
	if i.Prazo != "" {
		fmt.Fprintf(&b, "**Quando será paga:**\n\n- %s\n\n", i.Prazo)
		b.WriteString("_Dívida ASSUMIDA, registrada pelo Anchors. Não é defeito: é um dever " +
			"conhecido, ainda válido, com um momento declarado para ser cumprido. Quando " +
			"chegar a hora, mova para `todo/`; ao pagar, o confronto a fecha sozinho._\n")
		return b.String()
	}
	b.WriteString("_Registrada pelo Anchors. Resolva movendo para `doing/` e, ao tratar, para `done/`._\n")
	return b.String()
}

// pathFor devolve o caminho absoluto de uma issue num estado.
func pathFor(root string, state State, id string) string {
	return filepath.Join(root, Dir, string(state), id)
}

// byKey procura, em todos os estados, um arquivo de issue cujo nome termine em
// `--<key>.md` (a identidade estável). Devolve o estado, o nome do arquivo e se achou.
// É o coração da busca por identidade: a data no nome varia, a Key não.
func byKey(root, key string) (State, string, bool) {
	suffix := "--" + key + ".md"
	for _, st := range Estados {
		names, _ := List(root, st)
		for _, name := range names {
			if strings.HasSuffix(name, suffix) {
				return st, name, true
			}
		}
	}
	return "", "", false
}

// Exists diz se uma issue com esta Key já existe em QUALQUER estado — para não
// reabrir uma issue que já está sendo tratada ou foi resolvida (dedup por aresta,
// estável no tempo).
func Exists(root, key string) (State, bool) {
	st, _, ok := byKey(root, key)
	return st, ok
}

// Resolve move uma issue viva (todo/ ou doing/) para done/ — chamado quando o
// confronto que a gerou volta a PASSAR. Idempotente: se já está em done/ (ou não
// existe), não faz nada e devolve resolved=false. Este é o fechamento do loop: a
// issue deixa de mentir sobre o estado quando o problema é corrigido.
func Resolve(root, key string) (resolved bool, err error) {
	st, name, ok := byKey(root, key)
	if !ok || st == Done {
		return false, nil
	}
	from := pathFor(root, st, name)
	doneDir := filepath.Join(root, Dir, string(Done))
	if err := os.MkdirAll(doneDir, 0o755); err != nil {
		return false, err
	}
	if err := os.Rename(from, filepath.Join(doneDir, name)); err != nil {
		return false, err
	}
	return true, nil
}

// Open grava a issue em todo/ — mas só se ela ainda não existe em nenhum estado.
// Devolve (criada?, estado-atual). Idempotente: reconfrontar o mesmo alvo não
// duplica nem "ressuscita" uma issue já em doing/done.
func Open(root string, i Issue) (created bool, at State, err error) {
	return OpenAt(root, i, Todo)
}

// OpenAt grava a issue no estado em que ela NASCE. A maioria nasce em `todo/`; a dívida
// assumida nasce em `future/`, porque foi declarada com um momento para ser paga e
// colocá-la na fila do agora seria contar a mesma mentira ao contrário.
//
// A dedup continua valendo entre TODOS os estados: uma dívida que virou trabalho de agora
// (movida para `todo/`) não é recriada em `future/` no confronto seguinte.
func OpenAt(root string, i Issue, nasce State) (created bool, at State, err error) {
	if st, ok := Exists(root, i.Key()); ok {
		return false, st, nil // já existe (em qualquer estado) — não duplica
	}
	id := i.ID()
	dir := filepath.Join(root, Dir, string(nasce))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, "", err
	}
	if err := os.WriteFile(pathFor(root, nasce, id), []byte(i.Body()), 0o644); err != nil {
		return false, "", err
	}
	return true, nasce, nil
}

// List devolve os IDs das issues num estado (leitura para relatórios).
func List(root string, state State) ([]string, error) {
	dir := filepath.Join(root, Dir, string(state))
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

// Reabrir move uma issue de `done/` de volta para `todo/` e ACRESCENTA o novo laudo ao
// corpo, preservando o anterior.
//
// Existe porque um achado NOVO sobre uma unidade já revisada sumia em silêncio: o `Open`
// é idempotente por identidade (mesmo gate + mesmo alvo = mesma issue), então o segundo
// veredito só imprimia "issue já registrada" e descartava o `--reason`. Medido: dois
// laudos distintos foram perdidos assim, e a issue seguiu em `done/` com o texto antigo.
//
// Idempotência é a política certa para o MESMO problema detectado duas vezes; não é para
// um problema DIFERENTE no mesmo lugar. A diferença está no corpo, e por isso ele é
// comparado antes de decidir.
func Reabrir(root string, i Issue) (reaberta bool, err error) {
	st, name, ok := byKey(root, i.Key())
	if !ok {
		return false, nil
	}
	caminho := pathFor(root, st, name)
	antigo, rerr := os.ReadFile(caminho)
	if rerr != nil {
		return false, rerr
	}
	// laudo idêntico: é a mesma detecção, não um achado novo — nada a fazer.
	if i.Detail != "" && strings.Contains(string(antigo), strings.TrimSpace(i.Detail)) {
		return false, nil
	}
	corpo := string(antigo)
	if i.Detail != "" {
		corpo += "\n\n---\n\n## Achado adicional — " + i.Date + "\n\n" + i.Detail + "\n"
	}
	destino := pathFor(root, Todo, name)
	if st != Todo {
		if err := os.MkdirAll(filepath.Join(root, Dir, string(Todo)), 0o755); err != nil {
			return false, err
		}
		if err := os.Remove(caminho); err != nil {
			return false, err
		}
	}
	if err := os.WriteFile(destino, []byte(corpo), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// donoRE lê o dono do cabeçalho de uma issue já gravada.
var donoRE = regexp.MustCompile(`(?m)^- \*\*dono:\*\*\s*(\S+)\s*$`)

// DonoDoArquivo lê de quem é a issue, sem carregar o resto.
//
// Uma issue gravada ANTES do campo existir não o traz, e é do agente — que era o único
// caso. Devolver "usuário" nesse caso encheria a lista de quem decide com trabalho que
// não é dele, e o efeito seria a lista deixar de ser lida.
func DonoDoArquivo(caminho string) Dono {
	b, err := os.ReadFile(caminho)
	if err != nil {
		return DonoAgente
	}
	if m := donoRE.FindSubmatch(b); m != nil {
		return DonoDe(Dono(m[1]))
	}
	return DonoAgente
}

// ListaPorDono filtra as issues de um estado por dono. É a lista que se leva para a
// conversa com quem decide (`usuario`), ou a fila de trabalho do agente.
func ListaPorDono(root string, st State, dono Dono) ([]string, error) {
	nomes, err := List(root, st)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, n := range nomes {
		if DonoDoArquivo(filepath.Join(root, Dir, string(st), n)) == DonoDe(dono) {
			out = append(out, n)
		}
	}
	return out, nil
}

// Reatribui muda o dono de uma issue JÁ ABERTA, preservando o resto do arquivo.
//
// É o caminho para o agente que tentou resolver uma issue sua e esbarrou em algo que só
// uma pessoa responde: a issue continua sendo a violação que era, e passa a esperar quem
// pode resolvê-la. Sem isto, ela ficaria em `todo/` do agente sendo retentada para sempre,
// ou seria fechada sem que o problema tivesse sido resolvido.
func Reatribui(root string, st State, nome string, para Dono, porque string) error {
	caminho := filepath.Join(root, Dir, string(st), nome)
	b, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}
	texto := string(b)
	linha := "- **dono:** " + string(DonoDe(para))
	if donoRE.MatchString(texto) {
		texto = donoRE.ReplaceAllString(texto, linha)
	} else {
		// Issue anterior ao campo: insere depois do kind, que sempre existe.
		texto = strings.Replace(texto, "\n- **alvo", "\n"+linha+"\n- **alvo", 1)
	}
	// O PORQUÊ é obrigatório na prática: uma issue que muda de dono sem explicação chega
	// a quem decide sem contexto, e a primeira coisa que essa pessoa faz é perguntar o
	// que já se tentou.
	if porque != "" {
		texto = strings.TrimRight(texto, "\n") + "\n\n---\n**Passou a ser do " +
			string(DonoDe(para)) + ":** " + porque + "\n"
	}
	return os.WriteFile(caminho, []byte(texto), 0o644)
}

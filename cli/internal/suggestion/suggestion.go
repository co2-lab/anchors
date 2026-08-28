// Package suggestion materializa uma CORREÇÃO PROPOSTA que ainda não foi aplicada.
//
// Os gates até aqui só sabem ACUSAR: dizem que o carimbo divergiu, que falta amarra,
// que o `mock_detect` não alcança o dialeto — e param aí. Quem lê precisa reconstruir à
// mão o que a ferramenta já sabia. Pior no caso do julgamento por IA: ela leu o alvo
// inteiro para emitir o veredito, sabe exatamente qual seria a correção, e o formato de
// saída só permite descrevê-la em prosa.
//
// Uma sugestão é um PATCH do git mais o porquê. Isso muda o que o Anchors pode entregar:
// em vez de "o padrão declarado não casa o dialeto deste projeto", ele entrega o diff
// que conserta o `anchors.yaml` — revisável, aplicável com `git apply`, e recusável sem
// custo.
//
// Por que patch e não o arquivo já editado: uma correção aplicada direto some no meio do
// trabalho de quem estava editando, e some sem pedir licença. O patch preserva a ordem
// certa — o Anchors propõe, alguém decide. Também é o único formato que sobrevive a
// revisão: `git apply --check` diz se ainda casa, e o diff mostra exatamente o que muda.
//
// O ESTADO é a pasta, como nas issues (`issue.State`): mover é decidir. Uma sugestão
// aprovada e uma recusada são fatos distintos e ambos valem registro — a recusada evita
// que a mesma proposta volte na próxima varredura sem que ninguém saiba que já foi
// rejeitada.
package suggestion

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Dir é a raiz das sugestões, relativa à raiz do projeto.
//
// Vive em `suggestions/` (conteúdo de PROJETO, versionável), não em `.anchors/` (estado
// efêmero do daemon): uma proposta de correção é material de revisão — entra no diff,
// é discutida no PR e fica no histórico. O mesmo critério que colocou `issues/` fora do
// `.anchors/`.
const Dir = "suggestions"

// State — o ciclo de vida, codificado na pasta.
type State string

const (
	// Pending: proposta, aguardando decisão. É o estado em que toda sugestão nasce.
	Pending State = "pending"
	// Approved: aceita. O patch foi aplicado (ou está liberado para aplicar).
	Approved State = "approved"
	// Rejected: recusada, com razão. NÃO é lixo a apagar — é o registro que impede a
	// mesma proposta de voltar na próxima varredura como se fosse novidade, e que
	// preserva o motivo de quem decidiu para quem vier depois.
	Rejected State = "rejected"
)

// Origin — de onde a sugestão veio. Um gate determinístico que sabe a correção exata e
// uma IA que a inferiu merecem pesos diferentes na hora de revisar.
type Origin string

const (
	FromGate     Origin = "gate"     // gate determinístico: a correção é computada
	FromJudgment Origin = "judgment" // julgamento por IA: a correção é inferida
)

// Suggestion — uma correção proposta.
type Suggestion struct {
	ID     string // identificador estável; deriva de gate+alvo (dedup entre varreduras)
	Gate   string // quem propôs
	Target string // o arquivo/nó a que se refere
	Origin Origin
	// Why é o PORQUÊ, em prosa. Não é decorativo: um patch sem motivo é indistinguível
	// de uma mudança arbitrária, e quem revisa não tem como decidir sem ele.
	Why string
	// Patch é o diff unificado, aplicável com `git apply`. Vazio é permitido — há
	// achados cuja correção não é mecânica, e uma sugestão sem patch ainda vale como
	// diagnóstico registrado.
	Patch string
	// AutoJudged marca que a decisão foi tomada por IA sob `auto_judgment`, não por
	// pessoa. Fica no arquivo porque a diferença importa numa auditoria: "ninguém
	// olhou isto" é informação, e apagá-la faria decisão automática e decisão humana
	// parecerem a mesma coisa.
	AutoJudged bool
	// Reason preenche-se ao decidir: por que foi aprovada ou recusada.
	Reason string
}

// Open grava a sugestão em `suggestions/pending/`. Idempotente pelo ID: a mesma
// proposta reaparecendo numa varredura seguinte não gera duplicata — e, se já foi
// decidida (approved/rejected), não volta para pending.
func Open(root string, s Suggestion) (created bool, path string, err error) {
	if s.ID == "" {
		return false, "", fmt.Errorf("sugestão sem ID")
	}
	// Já decidida? Não reabre. Reabrir apagaria a decisão de quem já olhou, e a mesma
	// proposta recusada voltaria a cada varredura como se fosse nova.
	for _, st := range []State{Approved, Rejected, Pending} {
		p := filepath.Join(root, Dir, string(st), s.ID+".md")
		if _, err := os.Stat(p); err == nil {
			return false, p, nil
		}
	}
	dir := filepath.Join(root, Dir, string(Pending))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, "", err
	}
	p := filepath.Join(dir, s.ID+".md")
	if err := os.WriteFile(p, []byte(render(s)), 0o644); err != nil {
		return false, "", err
	}
	return true, p, nil
}

// Decide move a sugestão para approved/rejected, gravando a razão.
//
// Mover é a decisão — o mesmo princípio das issues. Não há campo de status dentro do
// arquivo que possa divergir da pasta em que ele está.
func Decide(root, id string, to State, reason string, autoJudged bool) error {
	if to != Approved && to != Rejected {
		return fmt.Errorf("estado de decisão inválido: %s", to)
	}
	if strings.TrimSpace(reason) == "" {
		// Decisão sem razão é a mesma falha que o `@no-test` nu: some o rastro de por
		// que alguém escolheu, e a escolha vira indistinguível de descuido.
		return fmt.Errorf("decisão exige razão escrita")
	}
	origem := filepath.Join(root, Dir, string(Pending), id+".md")
	b, err := os.ReadFile(origem)
	if err != nil {
		return fmt.Errorf("sugestão pendente não encontrada: %s", id)
	}
	destDir := filepath.Join(root, Dir, string(to))
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	carimbo := fmt.Sprintf("\n\n## Decisão\n\n- **estado:** %s\n- **em:** %s\n- **por:** %s\n\n%s\n",
		to, time.Now().Format("2006-01-02"), decisor(autoJudged), reason)
	if err := os.WriteFile(filepath.Join(destDir, id+".md"), append(b, []byte(carimbo)...), 0o644); err != nil {
		return err
	}
	return os.Remove(origem)
}

func decisor(auto bool) string {
	if auto {
		return "IA (auto_judgment)"
	}
	return "pessoa"
}

// List devolve os IDs num estado, ordenados.
func List(root string, st State) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, Dir, string(st)))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			out = append(out, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	sort.Strings(out)
	return out, nil
}

// PatchOf extrai o diff de uma sugestão, para `git apply`.
func PatchOf(root, id string, st State) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, Dir, string(st), id+".md"))
	if err != nil {
		return "", err
	}
	_, patch, ok := strings.Cut(string(b), "```diff\n")
	if !ok {
		return "", fmt.Errorf("sugestão %s não tem patch", id)
	}
	patch, _, _ = strings.Cut(patch, "\n```")
	return patch, nil
}

func render(s Suggestion) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# SUGGESTION: %s\n\n", s.Target)
	fmt.Fprintf(&b, "- **gate:** %s\n", s.Gate)
	fmt.Fprintf(&b, "- **origem:** %s\n", s.Origin)
	fmt.Fprintf(&b, "- **alvo:** %s\n", s.Target)
	fmt.Fprintf(&b, "- **proposta em:** %s\n\n", time.Now().Format("2006-01-02"))
	fmt.Fprintf(&b, "## Por quê\n\n%s\n", strings.TrimSpace(s.Why))
	if strings.TrimSpace(s.Patch) != "" {
		fmt.Fprintf(&b, "\n## Patch\n\n```diff\n%s\n```\n", strings.TrimRight(s.Patch, "\n"))
		b.WriteString("\nAplicar: `anchors suggest apply " + s.ID + "`" +
			" — ou `git apply` no diff acima.\n")
	} else {
		b.WriteString("\n## Patch\n\n_A correção não é mecânica: este achado precisa de decisão humana._\n")
	}
	return b.String()
}

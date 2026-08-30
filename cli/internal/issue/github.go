package issue

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// --- no modo github, o achado vira CARD, não arquivo ---
//
// O `issues/` é a fila do modo LOCAL: o protocolo é mover a pasta à mão (`todo/` →
// `doing/` → `done/`), e é assim que se pega trabalho sem GitHub.
//
// No modo github isso é duplicação — e a parte grave não é o arquivo sobrando, é o que
// NÃO acontece: medido no projeto de referência, 11 achados de gate (`trinca-completa`,
// `rule-types`, `guide-checklist`) foram para arquivo local e NENHUM virou card. O board
// mostrava o trabalho planejado e escondia o que os gates encontraram, e quem olhasse o
// board concluiria que não havia nada a corrigir.
//
// O ciclo de vida é o mesmo, com outro substrato:
//
//	arquivo em todo/      →  issue aberta, com a label do fluxo e `anchors:to-do`
//	mover para done/      →  issue FECHADA
//	reabrir de done/      →  issue REABERTA, com o novo laudo em comentário
//
// A DEDUPLICAÇÃO, que no local vem de procurar o arquivo pelas pastas, aqui vem de
// procurar a issue pela chave — que vai no corpo, num marcador estável. Sem isso cada
// execução do `check` abriria um card novo para o mesmo achado.

// MarcadorChave identifica a issue de um achado de gate no corpo do card.
//
// Vai no CORPO e não no título: o título é o que a pessoa lê, e prendê-lo ao formato da
// chave o tornaria ilegível ou frágil a qualquer mudança de redação.
const MarcadorChave = "<!-- anchors-issue-key: %s -->"

// GitHub é o repositório e a label do fluxo, de `workflow:`.
type GitHub struct {
	Repo  string
	Label string
}

func (g GitHub) gh(args ...string) ([]byte, error) {
	args = append(args, "--repo", g.Repo)
	out, err := exec.Command("gh", args...).CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("gh %s: %v — %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return out, nil
}

type cardEncontrado struct {
	Number int    `json:"number"`
	State  string `json:"state"`
}

// acha procura o card de um achado pela chave no corpo.
//
// Busca em `--state all` de propósito: uma issue FECHADA ainda é memória — é o que
// distingue "achado novo" de "achado que voltou", e sem ela o segundo perderia o laudo
// anterior.
func (g GitHub) acha(key string) (cardEncontrado, bool, error) {
	marca := fmt.Sprintf(MarcadorChave, key)
	out, err := g.gh("issue", "list", "--state", "all", "--limit", "500",
		"--search", key, "--json", "number,state,body")
	if err != nil {
		return cardEncontrado{}, false, err
	}
	var todas []struct {
		Number int    `json:"number"`
		State  string `json:"state"`
		Body   string `json:"body"`
	}
	if err := json.Unmarshal(out, &todas); err != nil {
		return cardEncontrado{}, false, fmt.Errorf("ler a busca: %w", err)
	}
	// A busca do GitHub é por texto e devolve aproximações: a CONFIRMAÇÃO é o marcador
	// exato no corpo. Sem ela um achado sobre `Foo.spec.md` casaria o card de
	// `FooBar.spec.md`, e o Anchors fecharia o card errado.
	for _, c := range todas {
		if strings.Contains(c.Body, marca) {
			return cardEncontrado{Number: c.Number, State: c.State}, true, nil
		}
	}
	return cardEncontrado{}, false, nil
}

// Open abre o card do achado, ou não faz nada se ele já existe.
//
// Se o card existe FECHADO, reabre e acrescenta o novo laudo — o achado voltou, e abrir
// um card novo perderia o histórico de quando ele apareceu antes.
func (g GitHub) Open(i Issue, nasce State) (created bool, at State, err error) {
	key := i.Key()
	c, existe, err := g.acha(key)
	if err != nil {
		return false, "", err
	}
	if existe {
		if c.State == "OPEN" {
			return false, Todo, nil
		}
		if _, err := g.gh("issue", "reopen", fmt.Sprint(c.Number)); err != nil {
			return false, "", err
		}
		if _, err := g.gh("issue", "comment", fmt.Sprint(c.Number),
			"--body", "⟲ **O achado voltou.**\n\n"+i.Detail); err != nil {
			return false, "", err
		}
		return true, Todo, nil
	}

	corpo := i.Body() + "\n\n" + fmt.Sprintf(MarcadorChave, key) + "\n"
	argv := []string{"issue", "create",
		"--title", g.titulo(i),
		"--body", corpo,
		"--label", g.Label,
	}
	// A DÍVIDA nasce diferente: ela vence depois, e misturá-la com o que se faz agora é o
	// caminho mais curto para ninguém olhar nenhuma das duas listas. Sem estado de fluxo,
	// ela fica fora do board até alguém decidir que chegou a hora.
	if nasce != Future {
		argv = append(argv, "--label", "anchors:to-do")
	}
	if i.Dono == DonoUsuário {
		argv = append(argv, "--label", "anchors:precisa-do-usuario")
	}
	if _, err := g.gh(argv...); err != nil {
		return false, "", err
	}
	return true, nasce, nil
}

// Resolve FECHA o card quando o confronto que o gerou volta a passar.
func (g GitHub) Resolve(key string) (bool, error) {
	c, existe, err := g.acha(key)
	if err != nil || !existe || c.State != "OPEN" {
		return false, err
	}
	_, err = g.gh("issue", "close", fmt.Sprint(c.Number),
		"--comment", "✓ Resolvido: o confronto que gerou este achado voltou a passar.")
	return err == nil, err
}

// titulo nomeia o card pelo que ele é, sem depender do formato da chave.
func (g GitHub) titulo(i Issue) string {
	que := "Violação"
	switch i.Kind {
	case Decision:
		que = "Decisão"
	case Stale:
		que = "Desatualizado"
	case Conflict:
		que = "Conflito"
	}
	return fmt.Sprintf("[%s] %s @ %s", i.Gate, que, i.Target)
}

// @anchors
//   ref: INQSN

package initx

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
)

// Question é uma decisão HUMANA do `init`, descrita de forma que um agente possa
// respondê-la sem ver a TUI.
//
// Existe porque o `init` é interativo e a guarda de TTY o aborta fora de um terminal —
// o que deixava um agente sem como iniciar projeto nenhum, justamente no fluxo em que o
// usuário pediu a ele para fazer isso (BOOTSTRAP.md §5).
//
// O contrato é de duas chamadas: a primeira devolve as perguntas (esta struct, em JSON),
// a segunda traz as respostas em flags e devolve o veredito de cada uma.
type Question struct {
	// ID é o nome da flag que responde esta pergunta (`--artifacts`, `--colocation`).
	ID string `json:"id"`
	// Texto é a pergunta como a TUI a faria — o agente a usa para explicar a escolha
	// ao usuário, que é quem decide de verdade.
	Texto string `json:"pergunta"`
	// Tipo diz a forma da resposta: "confirm", "select" ou "multiselect".
	Tipo string `json:"tipo"`
	// Opcoes são os valores aceitos (vazio em "confirm"). Responder fora desta lista é
	// erro — e é reportado como erro, não corrigido em silêncio.
	Opcoes []string `json:"opcoes,omitempty"`
	// Default é o que o Anchors INFERIU do disco. Um agente que não tenha base para
	// discordar deve aceitá-lo: é a leitura do projeto real, não um chute.
	Default any `json:"default"`
	// PorQue explica o que a resposta MUDA no projeto. Sem isto o agente escolhe pelo
	// nome da opção, que é como se escolhe errado.
	PorQue string `json:"por_que"`
}

// Respostas são os valores que a segunda chamada traz. Campos nulos (ponteiro nil ou
// slice vazia com a flag ausente) significam "não respondi" — e aí vale o default.
//
// A distinção entre "não respondi" e "respondi vazio" é o motivo dos ponteiros: para um
// `--artifacts=""` deliberado (nenhum artefato) não ser confundido com a flag ausente.
type Respostas struct {
	Header       *bool
	Contributing *bool
	Artifacts    *[]string
	Gates        *bool
	Colocation   *bool
	Layers       *[]string
	Governs      map[string][]string
	Workflow     *string
	Repo         *string
	Labels       *[]string
}

// StatusResposta é o veredito de UMA resposta, na saída da segunda chamada. O agente
// precisa saber não só que falhou, mas qual das sete respostas falhou e por quê.
type StatusResposta struct {
	ID       string `json:"id"`
	Valor    any    `json:"valor"`
	Aceita   bool   `json:"aceita"`
	Detalhe  string `json:"detalhe,omitempty"`
	UsouPada bool   `json:"usou_default,omitempty"`
}

// Questions monta a lista a partir do que foi inferido do disco. A ordem é a mesma da
// TUI: cada resposta restringe a seguinte, e apresentá-las fora de ordem faria o agente
// decidir camadas antes de saber se há co-location.
func Questions(p *Proposal) []Question {
	artefatosDetectados := []string{}
	if p != nil {
		for nome, sim := range p.DetectedArtifacts() {
			if sim {
				artefatosDetectados = append(artefatosDetectados, nome)
			}
		}
	}
	// Map order is random: unsorted, the same disk gave a default in a different order
	// on each call.
	sort.Strings(artefatosDetectados)
	colocado := false
	var camadas []string
	if p != nil {
		colocado = p.Colocated
		// `Config` pode vir nulo (inferência que não chegou a montar nada). Pedir as
		// camadas dali derrubaria o comando inteiro — e um `init --questions` que entra
		// em panic é pior do que um que devolve lista vazia: a lista vazia é a resposta
		// CERTA num projeto sem código.
		if p.Config != nil {
			camadas = CodeLayerNames(p.Config)
		}
	}

	// Texts and reasons go through i18n: the agent relays them to the user, who reads them
	// in the project's `lang:`. They were hard-coded in Portuguese.
	qs := []Question{
		{
			ID:      "header",
			Texto:   i18n.T("init.question.header.text"),
			Tipo:    "confirm",
			Default: true,
			PorQue:  i18n.T("init.question.header.why"),
		},
		{
			ID:      "contributing",
			Texto:   i18n.T("init.question.contributing.text"),
			Tipo:    "confirm",
			Default: true,
			PorQue:  i18n.T("init.question.contributing.why"),
		},
		{
			ID:      "artifacts",
			Texto:   i18n.T("init.question.artifacts.text"),
			Tipo:    "multiselect",
			Opcoes:  ArtifactNames(),
			Default: artefatosDetectados,
			PorQue:  i18n.T("init.question.artifacts.why"),
		},
		{
			ID:      "gates",
			Texto:   i18n.T("init.question.gates.text"),
			Tipo:    "confirm",
			Default: true,
			PorQue:  i18n.T("init.question.gates.why"),
		},
		{
			ID:      "colocation",
			Texto:   i18n.T("init.question.colocation.text"),
			Tipo:    "confirm",
			Default: colocado,
			PorQue:  i18n.T("init.question.colocation.why"),
		},
		{
			ID:      "layers",
			Texto:   i18n.T("init.question.layers.text"),
			Tipo:    "multiselect",
			Opcoes:  camadas,
			Default: camadas,
			PorQue:  i18n.T("init.question.layers.why"),
		},
		{
			ID:      "workflow",
			Texto:   i18n.T("init.question.workflow.text"),
			Tipo:    "select",
			Opcoes:  []string{"local", "manual", "github"},
			Default: "local",
			PorQue:  i18n.T("init.question.workflow.why"),
		},
		{
			ID:      "repo",
			Texto:   i18n.T("init.question.repo.text"),
			Tipo:    "texto",
			Default: "",
			PorQue:  i18n.T("init.question.repo.why"),
		},
		{
			ID:      "labels",
			Texto:   i18n.T("init.question.labels.text"),
			Tipo:    "multiselect",
			Default: []string{"anchors"},
			PorQue:  i18n.T("init.question.labels.why"),
		},
		{
			ID:      "governs",
			Texto:   i18n.T("init.question.governs.text"),
			Tipo:    "multiselect",
			Opcoes:  nil, // guide=tag1,tag2 — depende dos guides que existirem
			Default: map[string][]string{},
			PorQue:  i18n.T("init.question.governs.why"),
		},
	}
	return qs
}

// ValidateAnswers confere cada resposta contra as opções da pergunta e devolve o status
// de TODAS — não só das inválidas.
//
// Reportar as sete, e não apenas os erros, é o que permite ao agente conferir que o
// Anchors entendeu o que ele quis dizer. Uma resposta silenciosamente ignorada (flag
// escrita errada, por exemplo) seria indistinguível de uma aceita.
func ValidateAnswers(qs []Question, r Respostas) []StatusResposta {
	var out []StatusResposta
	for _, q := range qs {
		st := StatusResposta{ID: q.ID, Aceita: true}
		switch q.ID {
		case "header":
			st.Valor, st.UsouPada = valorBool(r.Header, q.Default)
		case "contributing":
			st.Valor, st.UsouPada = valorBool(r.Contributing, q.Default)
		case "gates":
			st.Valor, st.UsouPada = valorBool(r.Gates, q.Default)
		case "colocation":
			st.Valor, st.UsouPada = valorBool(r.Colocation, q.Default)
		case "artifacts":
			st.Valor, st.UsouPada, st.Aceita, st.Detalhe = listValue(r.Artifacts, q)
		case "layers":
			st.Valor, st.UsouPada, st.Aceita, st.Detalhe = listValue(r.Layers, q)
		case "workflow":
			if r.Workflow == nil {
				st.Valor, st.UsouPada = q.Default, true
			} else {
				st.Valor = *r.Workflow
				if !contem(q.Opcoes, *r.Workflow) {
					st.Aceita = false
					st.Detalhe = i18n.T("init.question.detail.unknown_mode", strings.Join(q.Opcoes, ", "))
				}
			}
		case "repo":
			st.Valor, st.UsouPada = valorTexto(r.Repo, q.Default)
			// A exigência é do MODO, não do campo: no `local` um `repo` declarado faz quem
			// lê o arquivo concluir que a integração está ativa (WORKFLOW.md §2).
			if modoGitHub(r) && empty(r.Repo) {
				st.Aceita = false
				st.Detalhe = i18n.T("init.question.detail.repo_required")
			}
			if !modoGitHub(r) && !empty(r.Repo) {
				st.Aceita = false
				st.Detalhe = i18n.T("init.question.detail.repo_only_github")
			}
		case "labels":
			st.Valor, st.UsouPada, st.Aceita, st.Detalhe = listValue(r.Labels, q)
			if st.Aceita && modoGitHub(r) && (r.Labels == nil || len(*r.Labels) == 0) {
				st.Aceita = false
				st.Detalhe = i18n.T("init.question.detail.labels_required")
			}
		case "governs":
			st.Valor = r.Governs
			st.UsouPada = len(r.Governs) == 0
		}
		out = append(out, st)
	}
	return out
}

// TudoAceito diz se a segunda chamada pode escrever. Uma resposta inválida recusa o
// conjunto INTEIRO: escrever as válidas produziria um anchors.yaml que ninguém decidiu
// por completo — a mesma régua da guarda de TTY, que prefere não escrever a escrever
// pela metade.
func TudoAceito(st []StatusResposta) bool {
	for _, s := range st {
		if !s.Aceita {
			return false
		}
	}
	return true
}

// modoGitHub diz se a resposta escolheu o modo `github` — é o que torna `repo` e
// `labels` obrigatórios.
func modoGitHub(r Respostas) bool {
	return r.Workflow != nil && *r.Workflow == "github"
}

func empty(p *string) bool { return p == nil || *p == "" }

func valorTexto(p *string, def any) (any, bool) {
	if p == nil {
		return def, true
	}
	return *p, false
}

func valorBool(p *bool, def any) (any, bool) {
	if p == nil {
		return def, true
	}
	return *p, false
}

func listValue(p *[]string, q Question) (valor any, usouDefault, aceita bool, detalhe string) {
	if p == nil {
		return q.Default, true, true, ""
	}
	for _, v := range *p {
		if len(q.Opcoes) > 0 && !contem(q.Opcoes, v) {
			return *p, false, false, i18n.T("init.question.detail.invalid_value", v, strings.Join(q.Opcoes, ", "))
		}
	}
	return *p, false, true, ""
}

func contem(lista []string, v string) bool {
	for _, x := range lista {
		if x == v {
			return true
		}
	}
	return false
}

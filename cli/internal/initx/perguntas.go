package initx

import "strings"

// Pergunta é uma decisão HUMANA do `init`, descrita de forma que um agente possa
// respondê-la sem ver a TUI.
//
// Existe porque o `init` é interativo e a guarda de TTY o aborta fora de um terminal —
// o que deixava um agente sem como iniciar projeto nenhum, justamente no fluxo em que o
// usuário pediu a ele para fazer isso (BOOTSTRAP.md §5).
//
// O contrato é de duas chamadas: a primeira devolve as perguntas (esta struct, em JSON),
// a segunda traz as respostas em flags e devolve o veredito de cada uma.
type Pergunta struct {
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
	Preset     *string
	Header     *bool
	Artifacts  *[]string
	Gates      *bool
	Colocation *bool
	Layers     *[]string
	Governs    map[string][]string
	Workflow   *string
	Repo       *string
	Labels     *[]string
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

// Perguntas monta a lista a partir do que foi inferido do disco. A ordem é a mesma da
// TUI: cada resposta restringe a seguinte, e apresentá-las fora de ordem faria o agente
// decidir camadas antes de saber se há co-location.
func Perguntas(p *Proposal, presets []string) []Pergunta {
	artefatosDetectados := []string{}
	if p != nil {
		for nome, sim := range p.DetectedArtifacts() {
			if sim {
				artefatosDetectados = append(artefatosDetectados, nome)
			}
		}
	}
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

	qs := []Pergunta{
		{
			ID:      "preset",
			Texto:   "Qual preset de stack usar?",
			Tipo:    "select",
			Opcoes:  append([]string{"nenhum"}, presets...),
			Default: "nenhum",
			PorQue: "o preset preenche as camadas de código com uma estrutura consagrada da stack. " +
				"`nenhum` mantém a inferência do disco — escolha um só se ele casar com o que o " +
				"PROJECT.md decidiu, porque camada errada faz gate reprovar arquivo certo.",
		},
		{
			ID:      "header",
			Texto:   "Semear guides/HEADER_GUIDE.md?",
			Tipo:    "confirm",
			Default: true,
			PorQue: "é o padrão MANDATÓRIO do bloco @anchors dos arquivos. Sem ele, cada arquivo " +
				"novo negocia o formato do cabeçalho de novo, e o gate de header não tem régua.",
		},
		{
			ID:      "artifacts",
			Texto:   "Quais tipos de âncora o projeto usa (ou vai usar)?",
			Tipo:    "multiselect",
			Opcoes:  ArtifactNames(),
			Default: artefatosDetectados,
			PorQue: "define as camadas de artefato e quais gates padrão fazem sentido. Num projeto " +
				"novo, marque o que PRETENDE usar — o default vem do que existe no disco, que " +
				"ainda é nada.",
		},
		{
			ID:      "gates",
			Texto:   "Semear os gates padrão (informativos) dos artefatos escolhidos?",
			Tipo:    "confirm",
			Default: true,
			PorQue: "os gates nascem INFORMATIVOS (não bloqueiam). É o que amarra os sinais de teste " +
				"ao ciclo sem escrever tudo à mão; promover para bloqueante é decisão posterior.",
		},
		{
			ID:      "colocation",
			Texto:   "Os derivados (spec/feature/teste) ficam AO LADO do código?",
			Tipo:    "confirm",
			Default: colocado,
			PorQue: "decide o glob de cada camada de artefato: ao lado do código (co-location) ou em " +
				"árvore separada. Responder o contrário do que o projeto faz deixa os artefatos " +
				"fora do mapa, e o que está fora do mapa não existe para os gates.",
		},
		{
			ID:      "layers",
			Texto:   "Quais diretórios de código tratar como camadas?",
			Tipo:    "multiselect",
			Opcoes:  camadas,
			Default: camadas,
			PorQue: "só o que for camada é REGIDO pelos gates. Num projeto vazio a lista vem vazia — " +
				"declare as camadas no anchors.yaml quando o código existir.",
		},
		{
			ID:      "workflow",
			Texto:   "Onde a fila de trabalho mora?",
			Tipo:    "select",
			Opcoes:  []string{"local", "github"},
			Default: "local",
			PorQue: "os modos são EXCLUDENTES (WORKFLOW.md §2): `local` guarda a fila em " +
				"`.anchors/tasks/` e as issues em `issues/`; `github` a põe nas issues do " +
				"repositório, com o estado de cada trabalho na coluna de um Project. Nunca um " +
				"com o outro de reserva — de qual fila veio esta task? é pergunta que ninguém " +
				"consegue responder depois do fato. O modo `github` EXIGE --repo e --labels.",
		},
		{
			ID:      "repo",
			Texto:   "Qual repositório do GitHub (owner/nome)?",
			Tipo:    "texto",
			Default: "",
			PorQue: "obrigatório no modo `github`, e NUNCA inferido do remote: num fork, inferir " +
				"faria a escrita cair no repositório errado — e escrita em lugar errado não se " +
				"desfaz com revert.",
		},
		{
			ID:      "labels",
			Texto:   "Quais labels marcam os cards do Anchors?",
			Tipo:    "multiselect",
			Default: []string{"anchors"},
			PorQue: "obrigatório no modo `github`. O board é COMPARTILHADO — carrega issues de " +
				"produto e de infra —, e a label é o que separa o que é do Anchors. Sem ela, um " +
				"agente pegaria uma issue de produto e a moveria para `IN PROGRESS`.",
		},
		{
			ID:      "governs",
			Texto:   "Qual guide rege quais tags?",
			Tipo:    "multiselect",
			Opcoes:  nil, // guide=tag1,tag2 — depende dos guides que existirem
			Default: map[string][]string{},
			PorQue: "a regra `governs` liga um guide aos arquivos que ele governa. Um guide que não " +
				"rege ninguém é débito que o doctor reporta; declare quando os guides existirem.",
		},
	}
	return qs
}

// ValidaRespostas confere cada resposta contra as opções da pergunta e devolve o status
// de TODAS — não só das inválidas.
//
// Reportar as sete, e não apenas os erros, é o que permite ao agente conferir que o
// Anchors entendeu o que ele quis dizer. Uma resposta silenciosamente ignorada (flag
// escrita errada, por exemplo) seria indistinguível de uma aceita.
func ValidaRespostas(qs []Pergunta, r Respostas) []StatusResposta {
	var out []StatusResposta
	for _, q := range qs {
		st := StatusResposta{ID: q.ID, Aceita: true}
		switch q.ID {
		case "preset":
			if r.Preset == nil {
				st.Valor, st.UsouPada = q.Default, true
			} else {
				st.Valor = *r.Preset
				if !contem(q.Opcoes, *r.Preset) {
					st.Aceita = false
					st.Detalhe = "preset desconhecido; aceitos: " + strings.Join(q.Opcoes, ", ")
				}
			}
		case "header":
			st.Valor, st.UsouPada = valorBool(r.Header, q.Default)
		case "gates":
			st.Valor, st.UsouPada = valorBool(r.Gates, q.Default)
		case "colocation":
			st.Valor, st.UsouPada = valorBool(r.Colocation, q.Default)
		case "artifacts":
			st.Valor, st.UsouPada, st.Aceita, st.Detalhe = valorLista(r.Artifacts, q)
		case "layers":
			st.Valor, st.UsouPada, st.Aceita, st.Detalhe = valorLista(r.Layers, q)
		case "workflow":
			if r.Workflow == nil {
				st.Valor, st.UsouPada = q.Default, true
			} else {
				st.Valor = *r.Workflow
				if !contem(q.Opcoes, *r.Workflow) {
					st.Aceita = false
					st.Detalhe = "modo desconhecido; aceitos: " + strings.Join(q.Opcoes, ", ")
				}
			}
		case "repo":
			st.Valor, st.UsouPada = valorTexto(r.Repo, q.Default)
			// A exigência é do MODO, não do campo: no `local` um `repo` declarado faz quem
			// lê o arquivo concluir que a integração está ativa (WORKFLOW.md §2).
			if modoGitHub(r) && vazio(r.Repo) {
				st.Aceita = false
				st.Detalhe = "obrigatório no modo `github` — sem ele, o fluxo não sabe de qual repositório puxar"
			}
			if !modoGitHub(r) && !vazio(r.Repo) {
				st.Aceita = false
				st.Detalhe = "só vale no modo `github`; no `local` este campo faz o arquivo mentir sobre a integração estar ativa"
			}
		case "labels":
			st.Valor, st.UsouPada, st.Aceita, st.Detalhe = valorLista(r.Labels, q)
			if st.Aceita && modoGitHub(r) && (r.Labels == nil || len(*r.Labels) == 0) {
				st.Aceita = false
				st.Detalhe = "obrigatório no modo `github` — sem ela, o fluxo pegaria qualquer issue do repositório, inclusive as de produto"
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

func vazio(p *string) bool { return p == nil || *p == "" }

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

func valorLista(p *[]string, q Pergunta) (valor any, usouDefault, aceita bool, detalhe string) {
	if p == nil {
		return q.Default, true, true, ""
	}
	for _, v := range *p {
		if len(q.Opcoes) > 0 && !contem(q.Opcoes, v) {
			return *p, false, false, "valor inválido: " + v + "; aceitos: " + strings.Join(q.Opcoes, ", ")
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

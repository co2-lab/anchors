package settings

import (
	"sort"
	"strings"
)

// --- os PERFIS: quem é quem no projeto ---
//
// Um projeto com vários contribuidores não tem um único tipo de trabalho, e nem todos podem
// fazer tudo. O `user_issues` booleano resolvia UMA pergunta — quem decide o produto — e o
// resto ficava implícito: o dev que revisava por segurança e o que revisava por performance
// liam a mesma régua, e o arquiteto e o QA recebiam o mesmo card.
//
// O perfil é a declaração de PAPEL, e as capacidades derivam dele. Duas razões para não
// declarar as capacidades diretamente:
//
//  1. o papel é o que a pessoa sabe dizer sobre si. "Sou dev" é uma resposta; "atuo em
//     needs-user, não escrevo plano, reviso código" é um formulário que ninguém preenche
//     com cuidado;
//  2. as capacidades mudam com o Anchors, o papel não. Quando uma capacidade nova aparecer,
//     ela nasce mapeada aos papéis existentes — e quem declarou "sou QA" não precisa
//     redeclarar nada.
//
// O QUE O PERFIL NÃO É: hierarquia. O arquiteto não "manda" no dev; ele responde outra
// pergunta. E não é permissão de repositório — o git e o GitHub cuidam disso. É sobre QUE
// TRABALHO o agente puxa e COMO ele se comporta diante do que não sabe.

// Role é um perfil de contribuição.
type Role string

const (
	// RoleDev implementa: código, feature, teste, documentação.
	RoleDev Role = "dev"
	// RoleQA escreve o que prova — cenários e testes — e ataca o que já existe.
	RoleQA Role = "qa"
	// RoleReviewerSec revisa com a lente de segurança.
	RoleReviewerSec Role = "reviewer-security"
	// RoleReviewerPerf revisa com a lente de performance.
	RoleReviewerPerf Role = "reviewer-performance"
	// RoleReviewer revisa sem lente declarada — o adversarial geral.
	RoleReviewer Role = "reviewer"
	// RolePO decide o rumo do produto: o que entra, o que fica fora, o que a régua não diz.
	RolePO Role = "product-owner"
	// RoleArchitect decide a estrutura: camadas, fronteiras, o que depende de quê.
	RoleArchitect Role = "architect"
)

// Capability é o que um perfil pode fazer no fluxo.
type Capability string

const (
	// CapDecideProduct — atuar nos cards escalonados (`needs-user`) e resolver a dúvida.
	//
	// É a capacidade que motivou os perfis existirem. Um agente que a tem pode responder
	// "o que a spec deveria decidir"; um que não a tem escala e segue para o próximo card.
	CapDecideProduct Capability = "decide-product"
	// CapDecideStructure — decidir camada, fronteira e dependência (o `anchors.yaml`).
	//
	// Separada da anterior porque as perguntas são diferentes: "esta tela deve mostrar o
	// valor?" é de produto; "esta camada pode importar daquela?" é de estrutura. Quem
	// responde uma pode não ter contexto para a outra.
	CapDecideStructure Capability = "decide-structure"
	// CapWritePlan — escrever e revisar plano.
	CapWritePlan Capability = "write-plan"
	// CapWriteSpec — escrever spec.
	CapWriteSpec Capability = "write-spec"
	// CapImplement — escrever código.
	CapImplement Capability = "implement"
	// CapProve — escrever feature e teste.
	CapProve Capability = "prove"
	// CapReview — revisar o trabalho de outro.
	CapReview Capability = "review"
)

// roleDef é a definição de um perfil.
type roleDef struct {
	// Title é como o perfil se apresenta.
	Title string
	// Does é o que a pessoa daquele perfil faz, em uma frase — o que o comando mostra na
	// hora de escolher.
	Does string
	// Caps são as capacidades.
	Caps []Capability
	// Lens é o foco da revisão, quando o perfil revisa. Vazio nos que não revisam.
	//
	// Existe porque "revisar" não é uma coisa só: quem procura vazamento de dado e quem
	// procura consulta em laço leem o mesmo código com perguntas diferentes, e um revisor
	// sem lente declarada tende a fazer a revisão que sabe fazer — não a que falta.
	Lens string
}

var roles = map[Role]roleDef{
	RoleDev: {
		Title: "Dev",
		Does:  "implementa a spec: código, cenários, teste e a documentação que acompanha",
		Caps:  []Capability{CapImplement, CapProve, CapReview},
	},
	RoleQA: {
		Title: "QA",
		Does:  "escreve o que prova — cenários e testes — e ataca o que já passou",
		Caps:  []Capability{CapProve, CapReview},
		Lens: "o que o teste NÃO prova. Mute a regra e rode a suíte: se ela continuar " +
			"verde, aquele teste não prova aquela linha. Ataque por execução, com entrada " +
			"de borda — não por leitura",
	},
	RoleReviewerSec: {
		Title: "Revisor — segurança",
		Does:  "revisa com a lente de segurança: o que sai, o que entra, o que se guarda",
		Caps:  []Capability{CapReview},
		Lens: "o dado que ESCAPA e a superfície que ABRE. Siga cada valor até onde ele " +
			"sai do sistema — resposta, log, métrica, push — e pergunte se ele precisava " +
			"sair. Confronte cada entrada: quem pode chamá-la, com que credencial, e o " +
			"que acontece sem ela. Desconfie do campo novo num tipo que atravessa " +
			"fronteira, e da dimensão de métrica com nome livre",
	},
	RoleReviewerPerf: {
		Title: "Revisor — performance",
		Does:  "revisa com a lente de custo: o que roda em laço, o que espera, o que cresce",
		Caps:  []Capability{CapReview},
		Lens: "o que roda MAIS VEZES do que parece. Procure a consulta dentro do laço, o " +
			"retry sem espera, a lista sem paginação e o dado que cresce sem limite. Num " +
			"serverless o custo é por invocação: uma função que chama outra por item de " +
			"uma lista de mil é mil invocações. Meça a ordem de grandeza, não o " +
			"milissegundo",
	},
	RoleReviewer: {
		Title: "Revisor",
		Does:  "revisa adversarialmente, sem lente fixa",
		Caps:  []Capability{CapReview},
		Lens: "o que os gates NÃO alcançam: regra que se contradiz com a vizinha, prosa " +
			"que afirma algo sobre outro artefato, teste que passa sem provar. Ataque por " +
			"execução — comando e saída, não dedução",
	},
	RolePO: {
		Title: "Product Owner",
		Does:  "decide o rumo do produto, e resolve as dúvidas escalonadas",
		Caps: []Capability{
			CapDecideProduct, CapWritePlan, CapWriteSpec, CapReview,
		},
	},
	RoleArchitect: {
		Title: "Arquiteto",
		Does:  "decide a estrutura — camadas, fronteiras, dependências — e resolve o que é dela",
		Caps: []Capability{
			CapDecideStructure, CapDecideProduct, CapWritePlan, CapWriteSpec,
			CapImplement, CapReview,
		},
	},
}

// KnownRoles devolve os perfis, em ordem estável.
func KnownRoles() []Role {
	out := make([]Role, 0, len(roles))
	for r := range roles {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// ParseRole reconhece um perfil pelo nome.
//
// Aceita abreviações comuns porque a pergunta é feita no terminal: quem responde digita
// `po`, não `product-owner`. Devolve vazio para o que não reconhece — e aí quem chama
// pergunta de novo, em vez de assumir um perfil.
func ParseRole(s string) Role {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "dev", "developer", "desenvolvedor":
		return RoleDev
	case "qa", "tester", "teste":
		return RoleQA
	case "reviewer-security", "security", "seguranca", "segurança", "sec":
		return RoleReviewerSec
	case "reviewer-performance", "performance", "perf":
		return RoleReviewerPerf
	case "reviewer", "review", "revisor":
		return RoleReviewer
	case "product-owner", "po", "product", "produto":
		return RolePO
	case "architect", "arquiteto", "arch":
		return RoleArchitect
	}
	return ""
}

// Title do perfil.
func (r Role) Title() string {
	if d, ok := roles[r]; ok {
		return d.Title
	}
	return string(r)
}

// Does devolve o que o perfil faz, em uma frase.
func (r Role) Does() string {
	if d, ok := roles[r]; ok {
		return d.Does
	}
	return ""
}

// Lens devolve o foco de revisão do perfil, ou vazio.
func (r Role) Lens() string {
	if d, ok := roles[r]; ok {
		return d.Lens
	}
	return ""
}

// Can diz se o perfil tem uma capacidade.
//
// Perfil desconhecido não pode NADA, e é deliberado: um `settings.yaml` escrito à mão com
// um perfil inventado não deve destravar capacidade nenhuma — o padrão fechado é o mesmo do
// booleano que este mecanismo substituiu.
func (r Role) Can(c Capability) bool {
	d, ok := roles[r]
	if !ok {
		return false
	}
	for _, cap := range d.Caps {
		if cap == c {
			return true
		}
	}
	return false
}

// Caps devolve as capacidades do perfil, em ordem estável.
func (r Role) Caps() []Capability {
	d, ok := roles[r]
	if !ok {
		return nil
	}
	out := append([]Capability(nil), d.Caps...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

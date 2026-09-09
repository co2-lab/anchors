package settings

import (
	"strings"
	"testing"
)

// SÓ O PO E O ARQUITETO decidem o produto. É a regra que motivou os perfis existirem:
// num projeto com vários contribuidores, nem todos podem decidir pelo produto, e um agente
// que pergunta a quem o está rodando obtém uma resposta que pode não ser a do dono.
func TestRole_quemDecideOProduto(t *testing.T) {
	decidem := map[Role]bool{RolePO: true, RoleArchitect: true}

	for _, r := range KnownRoles() {
		got := r.Can(CapDecideProduct)
		if got != decidem[r] {
			t.Errorf("%s: decide o produto = %v, queria %v", r, got, decidem[r])
		}
	}
}

// A ESTRUTURA é do arquiteto, e não do PO. As perguntas são diferentes: "esta tela deve
// mostrar o valor?" é de produto; "esta camada pode importar daquela?" é de estrutura — e
// quem responde uma pode não ter contexto para a outra.
func TestRole_aEstruturaEhDoArquiteto(t *testing.T) {
	if !RoleArchitect.Can(CapDecideStructure) {
		t.Error("o arquiteto tem de decidir a estrutura")
	}
	if RolePO.Can(CapDecideStructure) {
		t.Error("o PO não decide a estrutura — é outra pergunta")
	}
}

// QA E REVISORES executam; não decidem. Foi o pedido literal, e é o que separa quem tira
// dúvida de quem faz tarefa.
func TestRole_qaEReviewersNaoDecidem(t *testing.T) {
	for _, r := range []Role{RoleQA, RoleReviewer, RoleReviewerSec, RoleReviewerPerf, RoleDev} {
		if r.Can(CapDecideProduct) || r.Can(CapDecideStructure) {
			t.Errorf("%s não deve decidir nada — ele executa", r)
		}
		if !r.Can(CapReview) {
			t.Errorf("%s deve poder revisar", r)
		}
	}
}

// CADA PERFIL QUE REVISA TEM UMA LENTE, e elas são DISTINTAS.
//
// "Revisar" não é uma coisa só: quem procura vazamento de dado e quem procura consulta em
// laço leem o mesmo código com perguntas diferentes. Um revisor sem lente tende a fazer a
// revisão que SABE fazer — não a que falta.
func TestRole_asLentesSaoDistintas(t *testing.T) {
	comLente := []Role{RoleQA, RoleReviewer, RoleReviewerSec, RoleReviewerPerf}
	vistas := map[string]Role{}

	for _, r := range comLente {
		lens := r.Lens()
		if strings.TrimSpace(lens) == "" {
			t.Errorf("%s revisa e não declara lente", r)
			continue
		}
		if outro, repetida := vistas[lens]; repetida {
			t.Errorf("%s e %s têm a MESMA lente — então uma delas não precisa existir", r, outro)
		}
		vistas[lens] = r
	}

	// E o dev NÃO tem lente: ele revisa como parte da entrega, sem foco declarado. Dar
	// uma lente a ele faria o perfil de revisor perder a razão de existir.
	if RoleDev.Lens() != "" {
		t.Error("o dev não deve ter lente — é o perfil de revisor que a tem")
	}
}

// PERFIL DESCONHECIDO NÃO PODE NADA. Um `settings.yaml` escrito à mão com um perfil
// inventado não destrava capacidade nenhuma — o padrão fechado é o mesmo do booleano que
// este mecanismo substituiu.
func TestRole_desconhecidoNaoPodeNada(t *testing.T) {
	for _, r := range []Role{"", "tech-lead", "gerente", "DEV"} {
		for _, c := range []Capability{
			CapDecideProduct, CapDecideStructure, CapWritePlan,
			CapWriteSpec, CapImplement, CapProve, CapReview,
		} {
			if Role(r).Can(c) {
				t.Errorf("perfil %q destravou %q", r, c)
			}
		}
	}
}

// O nome é lido nas duas línguas e nas abreviações: a pergunta é no terminal, e quem
// responde digita `po`, não `product-owner`.
func TestParseRole(t *testing.T) {
	casos := map[string]Role{
		"dev": RoleDev, "developer": RoleDev, "desenvolvedor": RoleDev,
		"qa": RoleQA, "tester": RoleQA,
		"po": RolePO, "product-owner": RolePO, "produto": RolePO,
		"arch": RoleArchitect, "arquiteto": RoleArchitect,
		"sec": RoleReviewerSec, "seguranca": RoleReviewerSec, "segurança": RoleReviewerSec,
		"perf": RoleReviewerPerf, "performance": RoleReviewerPerf,
		"reviewer": RoleReviewer, "revisor": RoleReviewer,
		" DEV \n": RoleDev,
	}
	for entrada, esperado := range casos {
		if got := ParseRole(entrada); got != esperado {
			t.Errorf("ParseRole(%q) = %q, queria %q", entrada, got, esperado)
		}
	}
	// O que não se reconhece devolve vazio — e quem chama pergunta de novo, em vez de
	// assumir um perfil. Assumir aqui daria capacidades a quem não as pediu.
	for _, invalido := range []string{"", "tech-lead", "gerente", "sim"} {
		if got := ParseRole(invalido); got != "" {
			t.Errorf("ParseRole(%q) = %q, queria vazio", invalido, got)
		}
	}
}

// Todo perfil declara TÍTULO e o que FAZ — é o que o comando mostra na hora de escolher, e
// escolher sem saber o que cada um faz é escolher no escuro.
func TestKnownRoles_todosSeApresentam(t *testing.T) {
	rs := KnownRoles()
	if len(rs) < 7 {
		t.Fatalf("perfis = %d, esperava ao menos 7", len(rs))
	}
	for _, r := range rs {
		if strings.TrimSpace(r.Title()) == "" {
			t.Errorf("%s sem título", r)
		}
		if strings.TrimSpace(r.Does()) == "" {
			t.Errorf("%s não diz o que faz", r)
		}
		if len(r.Caps()) == 0 {
			t.Errorf("%s não tem capacidade nenhuma — então ele não faz nada", r)
		}
	}
}

// O CAMPO ANTIGO continua valendo para quem já declarou.
//
// Ignorá-lo faria quem declarou `user_issues: true` antes dos perfis perder a capacidade
// sem nada avisar, no meio de um trabalho.
func TestSettings_compatibilidadeComOCampoAntigo(t *testing.T) {
	antigo := Settings{UserIssues: Bool(true)}
	if !antigo.Can(CapDecideProduct) {
		t.Error("`user_issues: true` sem perfil deve continuar decidindo o produto")
	}
	if !antigo.Decided() {
		t.Error("quem declarou pelo campo antigo já decidiu — não deve ser perguntado de novo")
	}
	// Mas ele só responde a UMA pergunta: era um booleano, e não tinha como dizer mais.
	if antigo.Can(CapWritePlan) || antigo.Can(CapReview) {
		t.Error("o campo antigo não pode destravar capacidade que ele não declarava")
	}

	// E o perfil VENCE quando os dois existem: é a fonte nova, e manter os dois faria a
	// resposta vir da que alguém esquecesse de mudar.
	comOsDois := Settings{Role: RoleDev, UserIssues: Bool(true)}
	if comOsDois.Can(CapDecideProduct) {
		t.Error("o perfil declarado tem de vencer o campo antigo")
	}
}

package mapx

import "testing"

// grafoEvidencia — um roteiro que compõe um util, no molde real: SS-03 → login.yaml.
func grafoEvidencia() *Graph {
	return &Graph{
		Nodes: []Node{
			{ID: "suites/SS-03.yaml", Kind: KindTest, Rev: "t1"},
			{ID: "utils/login.yaml", Kind: KindTest, Rev: "u1"},
			{ID: "utils/launchApp.yaml", Kind: KindTest, Rev: "l1"},
		},
		Edges: []Edge{
			{From: "suites/SS-03.yaml", To: "utils/login.yaml", Type: EdgeDependsOn},
			{From: "utils/login.yaml", To: "utils/launchApp.yaml", Type: EdgeDependsOn},
		},
	}
}

func TestEvidenceClosureDesceTransitivamente(t *testing.T) {
	// O fecho tem de alcançar o util do util: SS-03 → login → launchApp. Parar no
	// primeiro nível deixaria de fora justamente as dependências profundas, que são as
	// que ninguém lembra de revalidar à mão.
	g := grafoEvidencia()
	fecho := g.EvidenceClosure("suites/SS-03.yaml")
	if len(fecho) != 2 {
		t.Fatalf("esperava 2 nós no fecho, veio %v", fecho)
	}
	if fecho["utils/login.yaml"] != "u1" || fecho["utils/launchApp.yaml"] != "l1" {
		t.Errorf("fecho com revs erradas: %v", fecho)
	}
}

func TestEvidenciaVenceQuandoUtilCompostoMuda(t *testing.T) {
	// O CASO QUE MOTIVOU TUDO: toquei o login.yaml depois de 8 suítes já terem passado.
	// O sinal da suíte não muda (o roteiro dela é o mesmo), então SignalStale() diz que
	// está fresco — e está errado.
	g := grafoEvidencia()
	g.Nodes[0].Signal = &TestSignal{
		Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1", "utils/launchApp.yaml": "l1"},
	}
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Fatalf("nada mudou ainda — devia estar fresco, veio %+v", ev)
	}
	// o util muda; o roteiro da suíte NÃO
	g.Nodes[1].Rev = "u2"
	if g.Nodes[0].SignalStale() {
		t.Fatal("SignalStale() do nó isolado não pode pegar isto — é por isso que o fecho existe")
	}
	ev := g.EvidenceStaleFor("suites/SS-03.yaml")
	if ev == nil {
		t.Fatal("a evidência TINHA de vencer: o util composto mudou de rev")
	}
	if ev.Own {
		t.Error("o próprio roteiro não mudou — Own devia ser falso")
	}
	if len(ev.Culprit) != 1 || ev.Culprit[0] != "utils/login.yaml" {
		t.Errorf("o culpado devia ser só o login.yaml, veio %v", ev.Culprit)
	}
}

func TestEvidenciaSemSinalNaoEhVencida(t *testing.T) {
	// Ausência de prova não é prova vencida: um teste que nunca rodou não tem evidência
	// para expirar, e reportá-lo aqui misturaria "nunca medido" com "medição velha".
	g := grafoEvidencia()
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Errorf("sem sinal = nada a julgar, veio %+v", ev)
	}
}

func TestEvidenciaSinalAntigoSemFechoCaiNoNoIsolado(t *testing.T) {
	// Sinal ingerido por uma versão que não gravava o fecho: julgar pelo nó isolado, em
	// vez de declarar vencida toda evidência anterior — o que transformaria uma melhoria
	// de precisão em centenas de falsos vencidos no dia em que entrasse.
	g := grafoEvidencia()
	g.Nodes[0].Signal = &TestSignal{Passed: 1, AtRev: "t1"} // sem ClosureRev
	g.Nodes[1].Rev = "u2"                                   // o util mudou
	if ev := g.EvidenceStaleFor("suites/SS-03.yaml"); ev != nil {
		t.Errorf("sem fecho carimbado, o util não conta — veio %+v", ev)
	}
	g.Nodes[0].Rev = "t2" // agora o PRÓPRIO teste mudou
	ev := g.EvidenceStaleFor("suites/SS-03.yaml")
	if ev == nil || !ev.Own {
		t.Errorf("o próprio arquivo mudou: devia vencer com Own=true, veio %+v", ev)
	}
}

func TestEvidenciaPodaNoPropagation(t *testing.T) {
	// Um nó @noPropagation afirma que mudanças nele não descem. Respeitar isso evita que
	// um arquivo deliberadamente volátil vença a evidência de metade da suíte.
	g := grafoEvidencia()
	g.Nodes[1].NoPropagation = true // login.yaml não propaga
	fecho := g.EvidenceClosure("suites/SS-03.yaml")
	if _, temLaunch := fecho["utils/launchApp.yaml"]; temLaunch {
		t.Errorf("não devia descer PASSANDO pelo noPropagation, veio %v", fecho)
	}
	if _, temLogin := fecho["utils/login.yaml"]; !temLogin {
		t.Error("o próprio nó noPropagation entra no fecho (ele É dependência direta)")
	}
}

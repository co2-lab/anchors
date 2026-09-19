package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// mapaComAdaptador monta o grafo mínimo do caso: um plano semeia o adaptador de uma
// fonte, e outro plano a consome.
func mapaComAdaptador(planoAdaptador, semente string) *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{{ID: planoAdaptador, Kind: mapx.KindPlan}},
		Edges: []mapx.Edge{{From: planoAdaptador, To: semente, Type: mapx.EdgeSeeds}},
	}
}

// O CASO REAL, medido no blue-eyes.
//
// O plano 0008 dizia "Fonte: **GA4**" e declarava `needs: plans/0005-home-e-indice.md`.
// O adaptador vinha do 0002, e essa dependência existia SÓ NA PROSA. Quando o 0002
// removeu o `Ga4Adapter` numa revisão, o 0008 ficou dependendo de uma fonte que ninguém
// ia construir — e nada acusou, porque os dois planos seguiram internamente coerentes.
func TestPlanSourceDeclared_acusaAFonteQueSoExisteNaProsa(t *testing.T) {
	t.Run("PSDPL-B01: A source that lives only in the prose is failed", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/Ga4Adapter.spec.md")
	n := mapx.Node{
		ID:    "plans/0008-frontend-web.md",
		Kind:  mapx.KindPlan,
		Needs: []string{"plans/0005-home-e-indice.md"},
	}

	v, msg := checkPlanSourceDeclared("Fonte: **GA4**, e ela é a exceção arquitetural.", n, "", g, nil)

	if v != Fail {
		t.Fatalf("o gate não acusou a dependência que só existe na prosa: %v — %s", v, msg)
	}
	if !strings.Contains(msg, "GA4") || !strings.Contains(msg, "0002") {
		t.Errorf("a mensagem não diz QUAL fonte nem ONDE está o adaptador:\n%s", msg)
	}
}

// Declarado o `needs:` do plano que constrói o adaptador, o gate passa. É a correção que
// a `FRWBF-R0002` fez à mão.
func TestPlanSourceDeclared_passaComONeedsDeclarado(t *testing.T) {
	t.Run("PSDPL-B03: With the owning plan declared in needs the gate passes", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/Ga4Adapter.spec.md")
	n := mapx.Node{
		ID:    "plans/0008-frontend-web.md",
		Kind:  mapx.KindPlan,
		Needs: []string{"plans/0005-home-e-indice.md", "plans/0002-plataforma.md"},
	}

	if v, msg := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", g, nil); v != Pass {
		t.Errorf("com o `needs:` declarado o gate devia passar: %v — %s", v, msg)
	}
}

// Várias fontes numa linha só — a forma que o plano 0007 usa: "Fontes: **Prometheus**,
// **ELK** e **CloudWatch**". Cada uma é confrontada.
func TestPlanSourceDeclared_confrontaCadaFonteDaLinha(t *testing.T) {
	t.Run("PSDPL-B04: Every source of the line is confronted on its own", func(t *testing.T) {})
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "plans/0002-plataforma.md", Kind: mapx.KindPlan}},
		Edges: []mapx.Edge{
			{From: "plans/0002-plataforma.md", To: "packages/lambdas/fontes/ElkAdapter.spec.md", Type: mapx.EdgeSeeds},
			{From: "plans/0002-plataforma.md", To: "packages/lambdas/fontes/PrometheusAdapter.spec.md", Type: mapx.EdgeSeeds},
		},
	}
	n := mapx.Node{ID: "plans/0007-infraestrutura.md", Kind: mapx.KindPlan}

	v, msg := checkPlanSourceDeclared("Fontes: **Prometheus** (nó e pod), **ELK** (o lado cliente).", n, "", g, nil)

	if v != Fail {
		t.Fatalf("esperava Fail com duas fontes não declaradas: %v — %s", v, msg)
	}
	for _, fonte := range []string{"Prometheus", "ELK"} {
		if !strings.Contains(msg, fonte) {
			t.Errorf("a fonte %q não foi acusada:\n%s", fonte, msg)
		}
	}
}

// O nome da fonte casa com o do arquivo INDEPENDENTE de caixa e pontuação: `GA4` casa com
// `Ga4Adapter.spec.md`, e `CloudWatch` com `CloudWatchAdapter.spec.md`.
func TestPlanSourceDeclared_casaNomeIndependenteDeCaixa(t *testing.T) {
	t.Run("PSDPL-B05: The source name matches the adapter regardless of case", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/CloudWatchAdapter.spec.md")
	n := mapx.Node{ID: "plans/0009-database.md", Kind: mapx.KindPlan}

	if v, _ := checkPlanSourceDeclared("Fonte: **cloudwatch**.", n, "", g, nil); v != Fail {
		t.Errorf("`cloudwatch` devia casar com `CloudWatchAdapter.spec.md`: %v", v)
	}
}

// Fonte cujo adaptador NINGUÉM semeia não é acusada: pode ser fonte de um plano futuro, e
// o gate não pode inventar uma dependência que não existe.
func TestPlanSourceDeclared_ignoraFonteSemAdaptadorSemeado(t *testing.T) {
	t.Run("PSDPL-B06: A source whose adapter nobody seeds is not charged", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/ElkAdapter.spec.md")
	n := mapx.Node{ID: "plans/0013-status-page.md", Kind: mapx.KindPlan}

	if v, msg := checkPlanSourceDeclared("Fonte: **Datadog**.", n, "", g, nil); v == Fail {
		t.Errorf("fonte sem adaptador semeado não devia reprovar: %s", msg)
	}
}

// O plano que semeia o próprio adaptador não precisa declarar `needs:` para si mesmo.
func TestPlanSourceDeclared_naoCobraOPlanoDeSiMesmo(t *testing.T) {
	t.Run("PSDPL-B07: The plan that seeds the adapter is not charged for itself", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/ElkAdapter.spec.md")
	n := mapx.Node{ID: "plans/0002-plataforma.md", Kind: mapx.KindPlan}

	if v, msg := checkPlanSourceDeclared("Fonte: **ELK**.", n, "", g, nil); v == Fail {
		t.Errorf("o plano que SEMEIA o adaptador não depende de si mesmo: %s", msg)
	}
}

// Sem linha de fonte não há o que confrontar — Skip, não Pass: o gate não mediu nada.
func TestPlanSourceDeclared_semLinhaDeFonteEhSkip(t *testing.T) {
	t.Run("PSDPL-B08: A plan with no source line returns Skip", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/ElkAdapter.spec.md")
	n := mapx.Node{ID: "plans/0016-endurecimento.md", Kind: mapx.KindPlan}

	if v, _ := checkPlanSourceDeclared("Este plano não nomeia fonte nenhuma.", n, "", g, nil); v != Skip {
		t.Errorf("plano sem linha de fonte devia ser Skip, veio %v", v)
	}
}

// O gate é de PLANO: uma spec que mencione "Fonte:" não é confrontada.
func TestPlanSourceDeclared_soValeParaPlano(t *testing.T) {
	t.Run("PSDPL-B09: An artifact that is not a plan returns Skip", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/Ga4Adapter.spec.md")
	n := mapx.Node{ID: "packages/lambdas/frontend/WebVitals.spec.md", Kind: mapx.KindSpec}

	if v, _ := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", g, nil); v != Skip {
		t.Errorf("o gate confronta PLANO; numa spec devia ser Skip, veio %v", v)
	}
}

// O laudo tem de dizer QUAL fonte e ONDE esta o adaptador: a correcao e uma
// declaracao, nao uma investigacao.
func TestPlanSourceDeclared_laudoNomeiaFonteEDono(t *testing.T) {
	t.Run("PSDPL-B02: The verdict names which source and where its adapter lives", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/Ga4Adapter.spec.md")
	n := mapx.Node{ID: "plans/0008-frontend-web.md", Kind: mapx.KindPlan}

	v, msg := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", g, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v: %s", v, msg)
	}
	if !strings.Contains(msg, "GA4") {
		t.Errorf("o laudo tem de nomear a FONTE; veio:\n%s", msg)
	}
	if !strings.Contains(msg, "plans/0002-plataforma.md") {
		t.Errorf("o laudo tem de nomear o PLANO dono do adaptador; veio:\n%s", msg)
	}
}

// O que nao foi medido nunca e aprovado: sem grafo e sem nenhum adaptador semeado, o
// veredito e Pending — aprovar ali carimbaria um confronto que nao aconteceu.
func TestPlanSourceDeclared_naoAprovaOQueNaoMediu(t *testing.T) {
	t.Run("PSDPL-I01: What was not measured is never approved", func(t *testing.T) {})
	n := mapx.Node{ID: "plans/0008-frontend-web.md", Kind: mapx.KindPlan}

	if v, msg := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", nil, nil); v != Pending {
		t.Errorf("sem grafo esperava Pending, veio %v: %s", v, msg)
	}

	semAdaptador := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "plans/0002-plataforma.md", Kind: mapx.KindPlan}},
		Edges: []mapx.Edge{{
			From: "plans/0002-plataforma.md",
			To:   "packages/lambdas/Qualquer.spec.md",
			Type: mapx.EdgeSeeds,
		}},
	}
	if v, msg := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", semAdaptador, nil); v != Pending {
		t.Errorf("sem adaptador nenhum semeado esperava Pending, veio %v: %s", v, msg)
	}
}

// A convencao de posse falha para o lado SEGURO: arquivo semeado fora do padrao
// `...Adapter.spec.md` nao e dono de fonte nenhuma, e nada e cobrado.
func TestPlanSourceDeclared_nomeForaDoPadraoNaoEhDono(t *testing.T) {
	t.Run("PSDPL-I02: A seeded file off the naming pattern owns nothing", func(t *testing.T) {})
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "plans/0002-plataforma.md", Kind: mapx.KindPlan}},
		Edges: []mapx.Edge{
			// fora do padrao: `Ga4Client`, nao `Ga4Adapter`
			{From: "plans/0002-plataforma.md", To: "packages/lambdas/fontes/Ga4Client.spec.md", Type: mapx.EdgeSeeds},
			// um adaptador de verdade, para o gate ter o que medir
			{From: "plans/0002-plataforma.md", To: "packages/lambdas/fontes/ElkAdapter.spec.md", Type: mapx.EdgeSeeds},
		},
	}
	n := mapx.Node{ID: "plans/0008-frontend-web.md", Kind: mapx.KindPlan}

	if v, msg := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", g, nil); v == Fail {
		t.Errorf("acusacao errada custa mais que a perdida — ensina a ignorar o gate: %s", msg)
	}
}

// A ORDEM das fases e regua do `fase-ordenada`: com o `needs:` declarado, este gate
// passa mesmo que o plano dono venha depois.
func TestPlanSourceDeclared_naoConfrontaOrdemDeFase(t *testing.T) {
	t.Run("PSDPL-X01: The gate does not confront the order of the phases", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0019-plataforma-tardia.md", "packages/lambdas/fontes/Ga4Adapter.spec.md")
	n := mapx.Node{
		ID:    "plans/0008-frontend-web.md",
		Kind:  mapx.KindPlan,
		Needs: []string{"plans/0019-plataforma-tardia.md"},
	}

	if v, msg := checkPlanSourceDeclared("Fonte: **GA4**.", n, "", g, nil); v != Pass {
		t.Errorf("ordem de fase e de outro gate; esperava Pass, veio %v: %s", v, msg)
	}
}

// O gate nao pede um `needs:` apontando para o nada: fonte que nenhum plano constroi
// nao e cobrada, senao o gate estaria pedindo uma mentira.
func TestPlanSourceDeclared_naoPedeNeedsParaONada(t *testing.T) {
	t.Run("PSDPL-X02: The gate does not demand a needs pointing at nothing", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/ElkAdapter.spec.md")
	n := mapx.Node{ID: "plans/0013-status-page.md", Kind: mapx.KindPlan}

	v, msg := checkPlanSourceDeclared("Fonte: **Datadog**.", n, "", g, nil)
	if v == Fail {
		t.Fatalf("cobrar isso pediria uma declaracao apontando para o nada: %s", msg)
	}
	if strings.Contains(msg, "Datadog") {
		t.Errorf("a fonte sem dono nao podia ser acusada no laudo: %s", msg)
	}
}

// A regua e o NOME em negrito da linha de fonte. O gate nao interpreta a prosa em volta:
// uma fonte nomeada para dizer que foi descartada e cobrada do mesmo jeito.
func TestPlanSourceDeclared_naoInterpretaAProsa(t *testing.T) {
	t.Run("PSDPL-X03: The gate does not interpret what the source is for", func(t *testing.T) {})
	g := mapaComAdaptador("plans/0002-plataforma.md", "packages/lambdas/fontes/Ga4Adapter.spec.md")
	n := mapx.Node{ID: "plans/0008-frontend-web.md", Kind: mapx.KindPlan}

	v, msg := checkPlanSourceDeclared("Fonte: **GA4**, que foi descartada nesta revisao.", n, "", g, nil)
	if v != Fail {
		t.Fatalf("interpretar a prosa nao e regua de gate bloqueante; veio %v: %s", v, msg)
	}
	if !strings.Contains(msg, "GA4") {
		t.Errorf("o laudo tem de nomear a fonte que a linha declarou: %s", msg)
	}
}

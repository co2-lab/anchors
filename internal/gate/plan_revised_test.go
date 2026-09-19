package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O plano REVISADO precisa avisar quem o lê.
//
// Editar o plano antigo perderia o registro — um plano é ÂNCORA, e editá-lo depois de
// implementado faz o documento descrever algo que não foi o que aconteceu. Por isso a
// revisão é um plano NOVO. O custo dessa escolha é a leitura fora de ordem: quem abre o
// antigo segue uma decisão que foi revista, e é isso que o gate impede.

func TestPlanRevised_NonPlanSkips(t *testing.T) {
	t.Run("PLRVP-B01: Non-plan artifacts skip confrontation", func(t *testing.T) {})
	n := mapx.Node{ID: "src/calc.go", Kind: mapx.KindCode}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	v, _ := checkPlanRevised("package main\n", n, "", g, nil)
	if v != Skip {
		t.Fatalf("não-plano deveria pular; veio %v", v)
	}
}

func TestPlanRevised_NoMapPending(t *testing.T) {
	t.Run("PLRVP-B02: Confronting without a map graph returns pending", func(t *testing.T) {})
	n := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	v, _ := checkPlanRevised("# Plano\n", n, "", nil, nil)
	if v != Pending {
		t.Fatalf("sem mapa deveria retornar Pending; veio %v", v)
	}
}

func TestPlanoSemRevisaoPula(t *testing.T) {
	t.Run("PLRVP-B03: A plan with neither revisions nor revisers skips confrontation", func(t *testing.T) {})
	t.Run("PLRVP-X03: Section amendment markers are not demanded on unrevised plans", func(t *testing.T) {})
	n := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, _ := checkPlanRevised("# Plano\n", n, "", g, nil); v != Skip {
		t.Errorf("plano sem revisão deveria pular; veio %v", v)
	}
}

func TestRevisesParaPlanoInexistente(t *testing.T) {
	t.Run("PLRVP-B04: Declaring a revision target that does not exist in the map fails", func(t *testing.T) {})
	orfao := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Revises: []string{"plans/nao-existe.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{orfao}}
	v, msg := checkPlanRevised("# Plano\n", orfao, "", g, nil)
	if v != Fail {
		t.Fatalf("revises órfão deveria reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "nao-existe.md") {
		t.Errorf("a mensagem deveria nomear o alvo ausente: %s", msg)
	}
}

func TestQuemRevisaEhLembradoDoAviso(t *testing.T) {
	t.Run("PLRVP-B05: A revising plan receives a pending reminder when the target lacks a top notice", func(t *testing.T) {})
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	// O plano REVISOR recebe o lembrete — é pendência, não falha: ele não fez nada de
	// errado, só falta um passo que está noutro arquivo.
	v, msg := checkPlanRevised("# Plano 0007\n", novo, "", g, nil)
	if v != Pending {
		t.Fatalf("quem revisa deveria receber PENDING com o lembrete; veio %v", v)
	}
	if !strings.Contains(msg, "plans/0001.md") || !strings.Contains(msg, "@revised-by") {
		t.Errorf("o lembrete deveria dizer ONDE escrever e O QUÊ: %s", msg)
	}
}

func TestLembreteSomeQuandoOAvisoExiste(t *testing.T) {
	t.Run("PLRVP-B06: The pending reminder on a revising plan clears once the target carries the notice", func(t *testing.T) {})
	t.Run("PLRVP-I03: The pending reminder clears once the revised plan is notified", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	// SEM o aviso no arquivo: o revisor recebe o lembrete.
	if err := os.WriteFile(filepath.Join(root, "plans/0001.md"),
		[]byte("# Plano 0001\n\n## Objetivo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkPlanRevised("# Plano 0007\n", novo, root, g, nil); v != Pending {
		t.Fatalf("sem o aviso, o revisor deveria receber o lembrete; veio %v", v)
	}

	// COM o aviso: o lembrete some e retorna Skip.
	if err := os.WriteFile(filepath.Join(root, "plans/0001.md"),
		[]byte("# Plano 0001\n\n> `@revised-by: plans/0007.md` — mudou.\n\n> `@amended-by: plans/0007.md` — a F01.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, msg := checkPlanRevised("# Plano 0007\n", novo, root, g, nil); v != Skip {
		t.Errorf("com o aviso escrito, o lembrete tem de sumir e retornar Skip; veio %v (%s)", v, msg)
	}
}

func TestPlanoRevisadoSemAvisoFalha(t *testing.T) {
	t.Run("PLRVP-B07: A revised plan lacking a top revision notice fails", func(t *testing.T) {})
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "ANT" + "GO"}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan, Code: "NOV" + "OO",
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	semAviso := "# Plano 0001\n\n## Objetivo\n\nFazer algo.\n"
	v, msg := checkPlanRevised(semAviso, antigo, "", g, nil)
	if v != Fail {
		t.Fatalf("plano revisado sem aviso deveria reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "plans/0007.md") {
		t.Errorf("a mensagem deveria NOMEAR quem revisa: %s", msg)
	}
}

func TestPlanoRevisadoAvisoTardeFalha(t *testing.T) {
	t.Run("PLRVP-B08: A revised plan placing the revision notice after line 40 fails", func(t *testing.T) {})
	t.Run("PLRVP-I01: Revision notices must be placed within the first 40 lines", func(t *testing.T) {})
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "ANT" + "GO"}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan, Code: "NOV" + "OO",
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	tarde := "# Plano 0001\n" + strings.Repeat("\ntexto\n", 30) +
		"\n> `@revised-by: plans/0007.md` — tarde demais.\n" +
		"> `@amended-by: plans/0007.md` — a parte está marcada; o problema é o topo.\n"
	if v, _ := checkPlanRevised(tarde, antigo, "", g, nil); v != Fail {
		t.Errorf("aviso fora do topo não avisa ninguém; veio %v", v)
	}
}

func TestPlanoRevisadoSoTopoEhPendencia(t *testing.T) {
	t.Run("PLRVP-B09: A revised plan with top notice but no section amendment markers returns pending", func(t *testing.T) {})
	t.Run("PLRVP-I02: Missing section amendment markers yield pending rather than failure", func(t *testing.T) {})
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "ANT" + "GO"}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan, Code: "NOV" + "OO",
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	soTopo := "# Plano 0001\n\n> `@revised-by: plans/0007.md` — a fase 2 mudou de ordem.\n\n## Objetivo\n"
	if v, msg := checkPlanRevised(soTopo, antigo, "", g, nil); v != Pending {
		t.Errorf("aviso de topo sem marcar as partes é pendência; veio %v (%s)", v, msg)
	}
}

func TestPlanoRevisadoCompletoPassa(t *testing.T) {
	t.Run("PLRVP-B10: A revised plan with top notice and marked section amendments passes", func(t *testing.T) {})
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "ANT" + "GO"}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan, Code: "NOV" + "OO",
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	completo := "# Plano 0001\n\n> `@revised-by: plans/0007.md` — a fase 2 mudou de ordem.\n\n" +
		"## Fases\n\n### ANT" + "GO-F02 — a régua\n\n" +
		"> `@amended-by: plans/0007.md` — esta fase passou a vir depois da F03.\n"
	if v, msg := checkPlanRevised(completo, antigo, "", g, nil); v != Pass {
		t.Errorf("com aviso e parte marcada deveria passar; veio %v (%s)", v, msg)
	}
}

func TestMarcacaoDeParteAceitaAsDuasFormas(t *testing.T) {
	t.Run("PLRVP-B11: Markdown alerts and metadata directives are both accepted as valid markers", func(t *testing.T) {})
	t.Run("PLRVP-X01: Language neutrality allows markdown alerts and metadata directives", func(t *testing.T) {})
	t.Run("PLRVP-X02: Prose description quality accompanying revision markers is not evaluated", func(t *testing.T) {})
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "ANT" + "GO"}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan, Code: "NOV" + "OO",
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}
	topo := "# Plano 0001\n\n> `@revised-by: plans/0007.md` — mudou.\n\n"

	// JÁ IMPLEMENTADA: o texto original permanece, com o aviso ABAIXO.
	implementada := topo + "## Fases\n\n### ANT" + "GO-F01 — a árvore\n\nCria os pacotes.\n\n" +
		"> `@amended-by: plans/0007.md` — daqui em diante os pacotes ganham mutação.\n" +
		"> O texto acima descreve o que FOI implementado.\n"
	if v, msg := checkPlanRevised(implementada, antigo, "", g, nil); v != Pass {
		t.Errorf("parte implementada com aviso abaixo deveria passar; veio %v (%s)", v, msg)
	}

	// NÃO IMPLEMENTADA: o texto é o novo, e o original fica preservado no aviso.
	futura := topo + "## Fases\n\n### ANT" + "GO-F04 — o CI com mutação\n\n" +
		"> `@amended-by: plans/0007.md`.\n" +
		"> **Era:** o CI roda instalar, lint, build e teste.\n" +
		"> **Por quê:** cobertura não responde se o teste PROVA a linha.\n"
	if v, msg := checkPlanRevised(futura, antigo, "", g, nil); v != Pass {
		t.Errorf("parte futura com o original preservado deveria passar; veio %v (%s)", v, msg)
	}

	// Markdown alerts (> [!IMPORTANT] e > [!WARNING]) também são aceitos.
	markdownAlerts := "# Plano 0001\n\n> [!IMPORTANT]\n> Revised by " + novo.Code + " with new scope\n\n" +
		"## Phases\n\n> [!WARNING]\n> Amended by " + novo.Code + " short note\n"
	if v, msg := checkPlanRevised(markdownAlerts, antigo, "", g, nil); v != Pass {
		t.Errorf("alerta markdown com código deveria passar; veio %v (%s)", v, msg)
	}
}

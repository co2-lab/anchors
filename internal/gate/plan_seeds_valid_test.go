package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func planCfg() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{
		"dao":            {Pattern: "packages/backend/models/**/*.ts", Regime: "declarativo"},
		"business-logic": {Pattern: "packages/backend/business-logic/**/*.ts"},
		"service-go":     {Pattern: "services/**/*.go"},
		"script-py":      {Pattern: "scripts/**/*.py"},
	}}
}

func planNode() mapx.Node { return mapx.Node{ID: "plans/sample.md", Kind: mapx.KindPlan} }

func setupPlanRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "packages"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "services"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestPlanSeeds_B01_NaoPlanoPula(t *testing.T) {
	t.Run("PSVPL-B01: Non-plan artifacts skip confrontation", func(t *testing.T) {})
	t.Run("PSVPL-I01: Plan seed validation applies exclusively to plan artifacts", func(t *testing.T) {})
	content := "- `packages/backend/models/metadata.spec.md`\n"
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindFeature, mapx.KindTest} {
		v, _ := checkPlanSeedsValid(content, mapx.Node{ID: "x", Kind: k}, "/tmp", nil, planCfg())
		if v != Skip {
			t.Fatalf("esperava Skip para kind %v, veio %v", k, v)
		}
	}
}

func TestPlanSeeds_B02_SemConfigFicaPendente(t *testing.T) {
	t.Run("PSVPL-B02: Missing project configuration returns Pending", func(t *testing.T) {})
	c := "- `packages/backend/business-logic/calc.spec.md`\n"
	v, msg := checkPlanSeedsValid(c, planNode(), "/tmp", nil, nil)
	if v != Pending {
		t.Fatalf("esperava Pending com cfg nil, veio %v (%s)", v, msg)
	}
}

func TestPlanSeeds_B03_SemSementesPula(t *testing.T) {
	t.Run("PSVPL-B03: Plans without seeded specifications skip confrontation", func(t *testing.T) {})
	c := "# Plano sem sementes de spec\nRealizar tarefas de infraestrutura e deployment.\n"
	v, _ := checkPlanSeedsValid(c, planNode(), "/tmp", nil, planCfg())
	if v != Skip {
		t.Fatalf("plano sem sementes deve retornar Skip, veio %v", v)
	}
}

func TestPlanSeeds_B04_TemplateIgnorado(t *testing.T) {
	t.Run("PSVPL-B04: Template specification references are ignored as templates", func(t *testing.T) {})
	t.Run("PSVPL-I04: Casual prose citations and templates are not treated as seeded paths", func(t *testing.T) {})
	c := "Copie o molde de `packages/backend/_TEMPLATE_SCREEN.spec.md` para criar a tela.\n"
	v, _ := checkPlanSeedsValid(c, planNode(), "/tmp", nil, planCfg())
	if v != Skip {
		t.Fatalf("plano apenas com templates deve retornar Skip, veio %v", v)
	}
}

func TestPlanSeeds_B05_NomeSoltoIgnorado(t *testing.T) {
	t.Run("PSVPL-B05: Bare specification file names without directory paths are ignored", func(t *testing.T) {})
	c := "Siga o modelo adotado em `SubscriptionScreen.spec.md` para orientar a implementação.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), "/tmp", nil, planCfg())
	if v == Fail {
		t.Fatalf("nome solto sem barra deve ser ignorado como prosa e não reprovar, veio %v (%s)", v, msg)
	}
}

func TestPlanSeeds_B06_FragmentoProsaSemDiretorioRealIgnorado(t *testing.T) {
	t.Run("PSVPL-B06: Informal path abbreviations without real top-level directories are ignored", func(t *testing.T) {})
	root := t.TempDir() // sem diretório 'features' na raiz
	c := "Consulte o rascunho em `features/subscription/Screen.spec.md`.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), root, nil, planCfg())
	if v == Fail {
		t.Fatalf("caminho com primeiro segmento inexistente no disco deve ser ignorado e não reprovar, veio %v (%s)", v, msg)
	}
}

func TestPlanSeeds_B07_CamadaGovernadaPassa(t *testing.T) {
	t.Run("PSVPL-B07: Seeded specifications targeting valid governed layers pass", func(t *testing.T) {})
	t.Run("PSVPL-I03: Seed validation evaluates structural layer validity without requiring file existence", func(t *testing.T) {})
	t.Run("PSVPL-X01: Existing file presence is not required for seeded specifications", func(t *testing.T) {})
	t.Run("PSVPL-X02: Plan progress synchronization is not evaluated by this gate", func(t *testing.T) {})
	t.Run("PSVPL-X03: Specification content and scenarios within seeded files are not verified", func(t *testing.T) {})
	root := setupPlanRoot(t)
	c := "- `packages/backend/business-logic/calc.spec.md` — **nasce**.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), root, nil, planCfg())
	if v != Pass {
		t.Fatalf("camada governada válida deve passar, veio %v (%s)", v, msg)
	}
}

func TestPlanSeeds_B08_MultiplasExtensoesDeLinguagem(t *testing.T) {
	t.Run("PSVPL-B08: Multiple target source file extensions resolve the governed layer", func(t *testing.T) {})
	root := setupPlanRoot(t)
	c := "- `services/order/handler.spec.md` — **nasce** (Go)\n" +
		"- `scripts/migrate/data.spec.md` — **nasce** (Python)\n"
	v, msg := checkPlanSeedsValid(c, planNode(), root, nil, planCfg())
	if v != Pass {
		t.Fatalf("extensões .go e .py devem resolver camadas governadas e passar; veio %v (%s)", v, msg)
	}
}

func TestPlanSeeds_B09_CamadaDeclarativaReprova(t *testing.T) {
	t.Run("PSVPL-B09: Seeded specifications targeting declarative layers fail", func(t *testing.T) {})
	t.Run("PSVPL-I02: Declarative layers reject specification seeding", func(t *testing.T) {})
	root := setupPlanRoot(t)
	c := "- `packages/backend/models/metadata.spec.md` — **nasce**.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), root, nil, planCfg())
	if v != Fail {
		t.Fatalf("semeadura em camada declarativa deve reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "RECONHECIDA") && !strings.Contains(msg, "RECOGNIZED") && !strings.Contains(msg, "declarativo") {
		t.Errorf("mensagem deve explicar que camada é declarativa: %s", msg)
	}
}

func TestPlanSeeds_B10_SemCamadaEmDiretorioRealReprova(t *testing.T) {
	t.Run("PSVPL-B10: Seeded specifications matching no declared layer in a real directory fail", func(t *testing.T) {})
	root := setupPlanRoot(t)
	c := "- `packages/backend/unknown-domain/orphan.spec.md` — **nasce**.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), root, nil, planCfg())
	if v != Fail {
		t.Fatalf("semeadura em diretório real sem camada correspondente deve reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "packages/backend/unknown-domain/orphan.spec.md") {
		t.Errorf("mensagem deve citar o caminho sem camada: %s", msg)
	}
}

func TestPlanSeeds_B11_MultiplosDefeitosAgregados(t *testing.T) {
	t.Run("PSVPL-B11: Multiple seed defects across declarative and undeclared layers are aggregated", func(t *testing.T) {})
	root := setupPlanRoot(t)
	c := "- `packages/backend/models/meta.spec.md` — **nasce**.\n" +
		"- `packages/backend/unknown-domain/orphan.spec.md` — **nasce**.\n"
	v, msg := checkPlanSeedsValid(c, planNode(), root, nil, planCfg())
	if v != Fail {
		t.Fatalf("múltiplos defeitos devem reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "meta.spec.md") || !strings.Contains(msg, "orphan.spec.md") {
		t.Errorf("mensagem deve conter ambos os caminhos defeituosos: %s", msg)
	}
	if !strings.Contains(msg, ";") {
		t.Errorf("defeitos de categorias distintas devem ser separados por ponto e vírgula: %s", msg)
	}
}

// A DOUTRINA DE PRODUTO semeada por um plano tem UMA regra: morar em `product/`.
//
// Nao se cobra dela camada de alvo (doutrina nao tem alvo — ela E' o artefato) nem
// existencia (semear significa que vai nascer). Cobra-se o lugar, porque o kind
// `product` vem do CAMINHO: nascida fora, ela e' lida como doc comum, e a spec que a
// realizar aponta para um arquivo que o mapa nao reconhece como doutrina.
func TestPlanSeedsValid_doutrinaForaDeProductReprova(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "src"), 0o755)
	n := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	v, msg := checkPlanSeedsValid("semeia `src/limite.doctrine.md`\n", n, root, nil, &config.Config{})
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "src/limite.doctrine.md") {
		t.Errorf("o veredito tem de NOMEAR a doutrina: %s", msg)
	}
}

func TestPlanSeedsValid_doutrinaEmProductPassa(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "product"), 0o755)
	n := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	if v, msg := checkPlanSeedsValid("semeia `product/limite.doctrine.md`\n", n, root, nil, &config.Config{}); v == Fail {
		t.Errorf("doutrina em product/ nao devia reprovar: %s", msg)
	}
}

// Um plano que semeia APENAS doutrinas (nenhuma spec) e' legitimo, e o guarda de
// "nenhuma spec semeada" nao pode engoli-lo: com ele antes da checagem de doutrina, o
// plano saia Skip — sem veredito, com a doutrina fora de `product/` passando em silencio.
func TestPlanSeedsValid_planoSoDeDoutrinaNaoEscapaPeloGuarda(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "src"), 0o755)
	n := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	if v, _ := checkPlanSeedsValid("so doutrina: `src/x.doctrine.md`\n", n, root, nil, &config.Config{}); v != Fail {
		t.Errorf("esperava Fail, veio %v — o guarda de seeds vazias engoliu o plano", v)
	}
}

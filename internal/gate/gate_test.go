package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestApplies(t *testing.T) {
	g := config.Gate{On: []string{"spec", "feature"}}
	if !applies(g, mapx.Node{Kind: mapx.KindSpec}, "") {
		t.Error("gate on=[spec,feature] deveria aplicar a spec")
	}
	if applies(g, mapx.Node{Kind: mapx.KindCode}, "") {
		t.Error("gate on=[spec,feature] NÃO deveria aplicar a code")
	}
}

func TestAppliesWithTags(t *testing.T) {
	g := config.Gate{On: []string{"code"}, Tags: []string{"screen"}}
	// nó code com tag screen → aplica
	if !applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"frontend", "screen"}}, "") {
		t.Error("gate tags=[screen] deveria aplicar a code com tag screen")
	}
	// nó code SEM a tag → não aplica
	if applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend"}}, "") {
		t.Error("gate tags=[screen] NÃO deveria aplicar a code sem a tag")
	}
	// kind errado, mesmo com a tag → não aplica
	if applies(g, mapx.Node{Kind: mapx.KindSpec, Tags: []string{"screen"}}, "") {
		t.Error("gate on=[code] NÃO deveria aplicar a spec")
	}
}

func TestAppliesComExcludeTags(t *testing.T) {
	// O caso real: `mutation-score` cobrado de camadas que a ferramenta de mutação nem
	// consegue rodar (declaração de infra, schema). Sem teste que as alcance, o Stryker sai
	// com "No tests were found" e não há relatório para ingerir — o gate ficava Pending
	// eterno pedindo uma medição impossível de produzir.
	g := config.Gate{On: []string{"code"}, ExcludeTags: []string{"resource", "schema-model"}}

	if applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "code", "resource"}}, "") {
		t.Error("nó com tag excluída NÃO deveria entrar no gate")
	}
	if applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "code", "schema-model"}}, "") {
		t.Error("basta UMA das tags excluídas para ficar de fora")
	}
	if !applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "code", "dao"}}, "") {
		t.Error("nó sem tag excluída deveria seguir sendo cobrado")
	}

	// A ordem importa: a exclusão vence o filtro positivo. Camadas carregam rótulos
	// transversais (`backend`, `code`) junto com o próprio, então um gate que declare os
	// dois filtros ainda tem de deixar a exceção de fora — senão declará-la não teria
	// efeito nenhum justamente onde ela importa.
	g2 := config.Gate{On: []string{"code"}, Tags: []string{"backend"}, ExcludeTags: []string{"resource"}}
	if applies(g2, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "resource"}}, "") {
		t.Error("exclude_tags deveria vencer o filtro positivo de tags")
	}
	if !applies(g2, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "dao"}}, "") {
		t.Error("o filtro positivo deveria continuar valendo para quem não é exceção")
	}
}

func TestJudgmentGateEmitsJudge(t *testing.T) {
	g := config.Gate{Name: "atomic", On: []string{"code"}, Measures: config.MeasuresJudgment, Guide: "guides/SCREEN_GUIDE.md"}
	r := runOne(g, mapx.Node{ID: "Home.tsx", Kind: mapx.KindCode}, t.TempDir(), nil, nil)
	if r.Verdict != Judge {
		t.Fatalf("gate judgment deveria emitir Judge, veio %q", r.Verdict)
	}
}

func TestAggregate_promotion(t *testing.T) {
	// um fail bloqueante barra a promoção; um fail informativo não.
	results := []Result{
		{Gate: "lint", Target: "a", Verdict: Pass, Blocking: true},
		{Gate: "spec-sections", Target: "b", Verdict: Fail, Blocking: true}, // barra
		{Gate: "coverage", Target: "c", Verdict: Fail, Blocking: false},     // não barra
	}
	p := Aggregate(results)

	if p.Passed {
		t.Error("com um fail bloqueante, Passed deveria ser false")
	}
	if len(p.Failures) != 2 {
		t.Errorf("esperava 2 failures (viram issues), got %d", len(p.Failures))
	}
	if len(p.Blocked) != 1 || p.Blocked[0].Gate != "spec-sections" {
		t.Errorf("esperava 1 bloqueio (spec-sections), got %+v", p.Blocked)
	}
}

func TestAggregate_allPass(t *testing.T) {
	p := Aggregate([]Result{
		{Gate: "lint", Verdict: Pass, Blocking: true},
		{Gate: "coverage", Verdict: Pass, Blocking: false},
	})
	if !p.Passed {
		t.Error("todos pass → Passed deveria ser true")
	}
	if len(p.Failures) != 0 {
		t.Errorf("nenhum fail esperado, got %d", len(p.Failures))
	}
}

func TestAggregate_informativeFailDoesNotBlock(t *testing.T) {
	// só fail informativo: gera issue (Failures) mas NÃO barra (Passed=true).
	p := Aggregate([]Result{{Gate: "coverage", Verdict: Fail, Blocking: false}})
	if !p.Passed {
		t.Error("fail informativo não deveria barrar a promoção")
	}
	if len(p.Failures) != 1 {
		t.Error("fail informativo ainda gera issue (registro)")
	}
}

func TestInternalCheckers(t *testing.T) {
	cases := []struct {
		checker string
		content string
		want    Verdict
	}{
		{"non-empty", "   \n  ", Fail},
		{"non-empty", "algo", Pass},
		{"spec-sections", "> **Código**: `LOGIX`\n### LOGIX-S01: ok", Pass},
		{"spec-sections", "| Regra | Descrição |\n| `BADEX-X01` | não-interativo |", Pass}, // tabela conta
		{"spec-sections", "| `BADEX-B01` | children |", Pass},                              // tabela sem backtick de coluna
		{"spec-sections", "- **HRMCX-S01** rótulo do mês no topo", Pass},                   // bullet-negrito conta
		{"spec-sections", "só um título\nsem estados", Fail},
		{"spec-sections", "menção solta BADEX-X01 em prosa", Fail}, // prosa não conta
		// O placeholder saiu daqui: quem o confronta é o `placeholder-preenchido`, com o
		// vocabulário universal. Aqui a régua é a SEÇÃO catalogada.
		{"spec-sections", "### LOGIX-S01\nUma regra escrita de verdade.", Pass}, // placeholder
		{"has-code", "it('LOGIX-A01: ...')", Pass},
		{"has-code", "sem identidade nenhuma", Fail},
		{"header-conforme", "// @anchors\n//   code: LGNNX\nconst x = 1", Pass}, // dono (code)
		{"header-conforme", "// @anchors\n//   ref: LGNNX\nconst x = 1", Pass},  // referência (ref) também conta
		{"header-conforme", "<!-- @anchors\n  code: SPCRX\n-->\n# spec", Pass},  // dialeto markdown
		{"header-conforme", "const x = 1 // nada aqui", Fail},                   // sem bloco
		{"header-conforme", "// @anchors\n//   layer: screen\nconst x=1", Fail}, // bloco sem code NEM ref
	}
	for _, c := range cases {
		fn := internalCheckers[c.checker]
		if fn == nil {
			t.Fatalf("checker %q não registrado", c.checker)
		}
		got, _ := fn(c.content, mapx.Node{})
		if got != c.want {
			t.Errorf("%s(%q) = %v, want %v", c.checker, c.content, got, c.want)
		}
	}
}

func TestHeaderConformeRecognizedLayer(t *testing.T) {
	presentation := mapx.Node{Kind: mapx.KindCode, Tags: []string{"frontend", "presentation"}}
	regida := mapx.Node{Kind: mapx.KindCode, Tags: []string{"frontend", "business-logic"}}
	fn := internalCheckers["header-conforme"]

	// reconhecida com só `layer:` → PASSA
	if v, _ := fn("// @anchors\n//   layer: presentation\n", presentation); v != Pass {
		t.Error("presentation com layer deveria passar")
	}
	// reconhecida SEM identidade nenhuma → FALHA
	if v, _ := fn("// @anchors\n//   updated_at: x\n", presentation); v != Fail {
		t.Error("reconhecida sem layer/code/ref deveria falhar")
	}
	// regida com só `layer:` (sem code/ref) → FALHA (exige posse/referência)
	if v, _ := fn("// @anchors\n//   layer: business-logic\n", regida); v != Fail {
		t.Error("regida só com layer (sem code/ref) deveria falhar")
	}
	// regida com ref → PASSA
	if v, _ := fn("// @anchors\n//   ref: BLBLX\n", regida); v != Pass {
		t.Error("regida com ref deveria passar")
	}
}

func TestHeaderConformeTestOfRecognizedLayer(t *testing.T) {
	// um TESTE de arquivo de camada reconhecida é classificado kind:test (perde a tag
	// da camada), mas declara layer:presentation no header → deve passar com layer.
	testNode := mapx.Node{Kind: mapx.KindTest, Tags: []string{"test"}}
	fn := internalCheckers["header-conforme"]
	if v, _ := fn("// @anchors\n//   layer: presentation\nimport x", testNode); v != Pass {
		t.Error("teste de presentation (kind:test) com layer no header deveria passar")
	}
	// mas um teste normal (sem layer reconhecida) segue exigindo code/ref
	if v, _ := fn("// @anchors\n//   updated_at: x\nimport x", testNode); v != Fail {
		t.Error("teste sem layer reconhecida nem code/ref deveria falhar")
	}
}

// `requires` filtra por CONTEÚDO do alvo: o gate que pergunta sobre uma marcação só
// se aplica a quem a usou. Sem isso, um gate de julgamento `on: [spec]` enfileira uma
// pergunta de IA para toda spec do projeto.
func TestAppliesWithRequires(t *testing.T) {
	root := t.TempDir()
	comMarca := "com.spec.md"
	semMarca := "sem.spec.md"
	if err := os.WriteFile(filepath.Join(root, comMarca), []byte("| `X-A01` @no-test: gateway de repasse |\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, semMarca), []byte("| `X-A01` | tem teste próprio |\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := config.Gate{On: []string{"spec"}, Requires: "@no-test"}

	if !applies(g, mapx.Node{ID: comMarca, Kind: mapx.KindSpec}, root) {
		t.Error("spec COM @no-test deveria aplicar")
	}
	if applies(g, mapx.Node{ID: semMarca, Kind: mapx.KindSpec}, root) {
		t.Error("spec SEM @no-test NÃO deveria aplicar")
	}
	// Alvo ilegível não aplica: melhor deixar de cobrar do que cobrar às cegas.
	if applies(g, mapx.Node{ID: "nao-existe.spec.md", Kind: mapx.KindSpec}, root) {
		t.Error("alvo ilegível NÃO deveria aplicar")
	}
	// Sem `requires`, o silêncio não desliga nada.
	semReq := config.Gate{On: []string{"spec"}}
	if !applies(semReq, mapx.Node{ID: semMarca, Kind: mapx.KindSpec}, root) {
		t.Error("gate sem requires deveria aplicar a qualquer spec do kind")
	}
}

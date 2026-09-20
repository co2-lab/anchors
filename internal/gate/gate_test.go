package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestApplies(t *testing.T) {
	t.Run("GTENG-B01: A gate reaches only the kinds it declares", func(t *testing.T) {})

	g := config.Gate{On: []string{"spec", "feature"}}
	if !applies(g, mapx.Node{Kind: mapx.KindSpec}, "") {
		t.Error("gate on=[spec,feature] deveria aplicar a spec")
	}
	if applies(g, mapx.Node{Kind: mapx.KindCode}, "") {
		t.Error("gate on=[spec,feature] NÃO deveria aplicar a code")
	}
}

func TestAppliesWithTags(t *testing.T) {
	t.Run("GTENG-B02: A gate that names labels reaches only the nodes carrying one", func(t *testing.T) {})

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
	t.Run("GTENG-B03: One excluded label is enough to keep a node out", func(t *testing.T) {})

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
}

// A ORDEM importa: a exclusão vence o filtro positivo. Camadas carregam rótulos
// transversais (`backend`, `code`) junto com o próprio, então um gate que declare os dois
// filtros ainda tem de deixar a exceção de fora — senão declará-la não teria efeito
// nenhum justamente onde ela importa.
func TestExclusaoVenceOFiltroPositivo(t *testing.T) {
	t.Run("GTENG-B04: Exclusion wins over the positive label filter", func(t *testing.T) {})

	g := config.Gate{On: []string{"code"}, Tags: []string{"backend"}, ExcludeTags: []string{"resource"}}
	if applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "resource"}}, "") {
		t.Error("exclude_tags deveria vencer o filtro positivo de tags")
	}
	if !applies(g, mapx.Node{Kind: mapx.KindCode, Tags: []string{"backend", "dao"}}, "") {
		t.Error("o filtro positivo deveria continuar valendo para quem não é exceção")
	}
}

func TestJudgmentGateEmitsJudge(t *testing.T) {
	t.Run("GTENG-B11: A judgment gate asks for judgment when nothing has answered it", func(t *testing.T) {})

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
		{"header-valid", "// @anchors\n//   code: LGNNX\nconst x = 1", Pass}, // dono (code)
		{"header-valid", "// @anchors\n//   ref: LGNNX\nconst x = 1", Pass},  // referência (ref) também conta
		{"header-valid", "<!-- @anchors\n  code: SPCRX\n-->\n# spec", Pass},  // dialeto markdown
		{"header-valid", "const x = 1 // nada aqui", Fail},                   // sem bloco
		{"header-valid", "// @anchors\n//   layer: screen\nconst x=1", Fail}, // bloco sem code NEM ref
	}
	for _, c := range cases {
		// O checker pode estar em qualquer um dos dois registros: os puros de conteúdo e
		// os que recebem grafo+config. Procurar só o primeiro fazia o teste reprovar quando
		// um checker passava a precisar de `cfg` — mudança de assinatura, não de
		// comportamento, que é justamente o que este teste NÃO deveria acusar.
		fn := internalCheckers[c.checker]
		if fn == nil {
			if withGraph := checkersWithGraph[c.checker]; withGraph != nil {
				fn = func(content string, n mapx.Node) (Verdict, string) {
					return withGraph(content, n, "", nil, nil)
				}
			}
		}
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
	fn := internalCheckers["header-valid"]

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
	fn := internalCheckers["header-valid"]
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
	t.Run("GTENG-B05: A gate demanding a mark reaches only the targets that carry it", func(t *testing.T) {})

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
	// Sem `requires`, o silêncio não desliga nada.
	semReq := config.Gate{On: []string{"spec"}}
	if !applies(semReq, mapx.Node{ID: semMarca, Kind: mapx.KindSpec}, root) {
		t.Error("gate sem requires deveria aplicar a qualquer spec do kind")
	}
}

// Alvo ILEGÍVEL não aplica: melhor deixar de cobrar do que cobrar às cegas de um alvo
// cujo conteúdo não se conhece.
func TestAlvoIlegivelNaoAplica(t *testing.T) {
	t.Run("GTENG-B06: An unreadable target does not apply", func(t *testing.T) {})

	g := config.Gate{On: []string{"spec"}, Requires: "@no-test"}
	if applies(g, mapx.Node{ID: "nao-existe.spec.md", Kind: mapx.KindSpec}, t.TempDir()) {
		t.Error("alvo ilegível NÃO deveria aplicar")
	}
}

// ── ROTEAMENTO E AGREGAÇÃO: o que o motor decide sozinho ─────────────────────

// Nenhum alvo relevante: o gate NÃO RODA. Vale inclusive para escopo de projeto — um
// commit só de README não dispara o typecheck do monorepo inteiro.
func TestGateSemAlvoNaoRoda(t *testing.T) {
	t.Run("GTENG-B07: A gate with no applicable target does not run at all", func(t *testing.T) {})

	g := config.Gate{Name: "typecheck", ID: "typecheck", On: []string{"code"},
		Scope: config.ScopeProject, Run: "exit 1"}
	// O recorte só tem spec: o gate declara `code` e não encontra alvo.
	res := RunWithConfig([]config.Gate{g},
		[]mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec}}, t.TempDir(), nil, nil)
	if len(res) != 0 {
		t.Fatalf("sem alvo o gate não roda; vieram %d resultado(s): %+v", len(res), res)
	}
}

// Ferramenta exigida e AUSENTE: Skip, nunca Fail. O gate não mediu — e "não medi" não é
// "está limpo" nem "está sujo". Um Fail diria que o projeto violou algo, quando quem
// falta é o binário.
func TestFerramentaAusenteEhSkipNuncaFail(t *testing.T) {
	t.Run("GTENG-B08: A gate whose required binary is absent steps aside", func(t *testing.T) {})

	const inexistente = "binario-que-ninguem-instalou-jamais"
	g := config.Gate{Name: "semgrep", ID: "semgrep", On: []string{"code"},
		Run: inexistente, NeedsTool: inexistente}
	res := RunWithConfig([]config.Gate{g},
		[]mapx.Node{{ID: "a.ts", Kind: mapx.KindCode}}, t.TempDir(), nil, nil)

	if len(res) != 1 {
		t.Fatalf("esperava um veredito, vieram %d", len(res))
	}
	if res[0].Verdict != Skip {
		t.Errorf("ferramenta ausente é Skip, nunca Fail; veio %v", res[0].Verdict)
	}
	if !strings.Contains(res[0].Detail, inexistente) {
		t.Errorf("o laudo tem de NOMEAR o binário que falta: %s", res[0].Detail)
	}
}

// DISPENSA POR ALVO: o gate roda, e só este nó é poupado. O veredito é Skip com o motivo
// escrito — some do placar de reprovações sem sumir do relatório.
func TestDispensaPorAlvoPoupaUmESoUm(t *testing.T) {
	t.Run("GTENG-B09: A waiver by target spares one node and confronts the rest", func(t *testing.T) {})

	root := t.TempDir()
	for _, f := range []string{"poupado.md", "confrontado.md"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("   \n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := config.Gate{Name: "não-vazio", ID: "non-empty", On: []string{"spec"}, Check: "non-empty"}
	disp, erros := ParseWaiver("non-empty@POUPA=rascunho declarado do plano 0007")
	if len(erros) > 0 {
		t.Fatalf("dispensa malformada no fixture: %v", erros)
	}

	res := RunWithWaiver([]config.Gate{g}, []mapx.Node{
		{ID: "poupado.md", Kind: mapx.KindSpec, Code: "POUPA"},
		{ID: "confrontado.md", Kind: mapx.KindSpec, Code: "CONFR"},
	}, root, nil, nil, false, disp)

	if len(res) != 2 {
		t.Fatalf("o gate tem de RODAR e produzir veredito para os dois; vieram %d", len(res))
	}
	porAlvo := map[string]Result{}
	for _, r := range res {
		porAlvo[r.Target] = r
	}
	if v := porAlvo["poupado.md"].Verdict; v != Skip {
		t.Errorf("o alvo dispensado é Skip, veio %v", v)
	}
	if v := porAlvo["confrontado.md"].Verdict; v != Fail {
		t.Errorf("o alvo NÃO dispensado continua confrontado e reprova; veio %v", v)
	}
}

// O gate AGREGADO roda UMA vez, e o alvo reportado é o próprio escopo. Atribuir a falha
// a um dos arquivos seria mentira — o erro pode estar em qualquer um, ou na relação
// entre eles.
func TestGateAgregadoRodaUmaVezContraOEscopo(t *testing.T) {
	t.Run("GTENG-B10: An aggregate gate runs once and reports against the scope", func(t *testing.T) {})

	g := config.Gate{Name: "typecheck", ID: "typecheck", On: []string{"code"},
		Scope: config.ScopeProject, Run: "true"}
	alvos := []mapx.Node{
		{ID: "a.ts", Kind: mapx.KindCode},
		{ID: "b.ts", Kind: mapx.KindCode},
		{ID: "c.ts", Kind: mapx.KindCode},
	}
	res := RunWithConfig([]config.Gate{g}, alvos, t.TempDir(), nil, nil)

	if len(res) != 1 {
		t.Fatalf("o agregado responde UMA vez sobre o conjunto; vieram %d", len(res))
	}
	for _, n := range alvos {
		if res[0].Target == n.ID {
			t.Fatalf("o alvo do agregado é o ESCOPO, nunca um dos arquivos; veio %q", res[0].Target)
		}
	}
	if !strings.Contains(res[0].Target, config.ScopeProject) {
		t.Errorf("o alvo tem de nomear o escopo; veio %q", res[0].Target)
	}
}

// Sem esta consulta o gate PEDIA JULGAMENTO PARA SEMPRE: o `judge` gravava o veredito no
// grafo e o `check` seguinte emitia Judge de novo, então o contador de pendências nunca
// descia e o trabalho já feito era invisível.
func TestJulgamentoJaRespondidoViraVeredito(t *testing.T) {
	t.Run("GTENG-B12: A judgment gate reads the stamp an earlier judgement left", func(t *testing.T) {})

	g := config.Gate{Name: "atomic", On: []string{"code"}, Measures: config.MeasuresJudgment}
	alvo := mapx.Node{ID: "Home.tsx", Kind: mapx.KindCode, Rev: "r1"}

	for _, c := range []struct {
		carimbo string
		quer    Verdict
	}{
		{"ok", Pass},
		{"issue", Fail},
	} {
		graph := grafoComJulgamento(alvo, g.Name, c.carimbo)
		r := runOne(g, alvo, t.TempDir(), graph, nil)
		if r.Verdict != c.quer {
			t.Errorf("o carimbo %q deveria virar %v, veio %v (%s)", c.carimbo, c.quer, r.Verdict, r.Detail)
		}
	}
}

// `dispensado` vira Skip, e NÃO Pass: o gate NÃO MEDIU — o alvo da pergunta não existe.
// Pass afirmaria aprovação, que é a mentira que este veredito existe para evitar.
func TestJulgamentoDispensadoEhSkipNuncaPass(t *testing.T) {
	t.Run("GTENG-B13: A judgement recorded as waived becomes Skip and never Pass", func(t *testing.T) {})

	g := config.Gate{Name: "atomic", On: []string{"code"}, Measures: config.MeasuresJudgment}
	alvo := mapx.Node{ID: "Home.tsx", Kind: mapx.KindCode, Rev: "r1"}

	// A forma nova e a ANTIGA: o valor antigo está em mapas já commitados, e ignorá-lo
	// faria o gate repreguntar um julgamento que alguém respondeu.
	for _, carimbo := range []string{"waived", "dispensado"} {
		r := runOne(g, alvo, t.TempDir(), grafoComJulgamento(alvo, g.Name, carimbo), nil)
		if r.Verdict != Skip {
			t.Errorf("o carimbo %q tem de virar Skip, veio %v (%s)", carimbo, r.Verdict, r.Detail)
		}
	}
}

// grafoComJulgamento monta um mapa mínimo com UM julgamento vivo sobre o alvo.
func grafoComJulgamento(alvo mapx.Node, gate, veredito string) *mapx.Graph {
	vizinho := mapx.Node{ID: "Home.spec.md", Kind: mapx.KindSpec, Rev: "r0"}
	return &mapx.Graph{
		Nodes: []mapx.Node{alvo, vizinho},
		Edges: []mapx.Edge{{
			From: vizinho.ID, To: alvo.ID, Type: mapx.EdgeSpecifies,
			Julgamentos: []mapx.Judgment{{
				Gate: gate, Verdict: veredito,
				ValidatedFromRev: vizinho.Rev, ValidatedToRev: alvo.Rev,
			}},
		}},
	}
}

// Um gate que não declara NEM comando NEM check não tem como responder. Passar em
// silêncio carimbaria verde sobre uma pergunta que ninguém sabe fazer.
func TestGateSemComandoNemCheckEhIndeterminado(t *testing.T) {
	t.Run("GTENG-B14: A gate declaring neither a command nor a check is undetermined", func(t *testing.T) {})

	g := config.Gate{Name: "mudo", ID: "mudo", On: []string{"code"}}
	r := runOne(g, mapx.Node{ID: "a.ts", Kind: mapx.KindCode}, t.TempDir(), nil, nil)
	if r.Verdict != Pending {
		t.Fatalf("sem run nem check o veredito é indeterminado, veio %v", r.Verdict)
	}
	if r.Detail == "" {
		t.Error("o laudo tem de NOMEAR a omissão, e não calar")
	}
}

// Só a pendência que diz "há decisão POR TOMAR" impede a promoção. A que diz "não tive o
// que confrontar" não: medido num repositório real, tratar todos igual reprovou 411 nós
// de uma vez, e 410 eram gates sem sinal ingerido.
func TestSoADecisaoPorTomarImpedeAPromocao(t *testing.T) {
	t.Run("GTENG-B15: Only the pending item that says a decision is still to take bars promotion", func(t *testing.T) {})

	root := t.TempDir()
	comDecisao, semDecisao := "com.spec.md", "sem.spec.md"

	// Um checker de teste no lugar do real: o que se confronta aqui é a LEITURA que o
	// motor faz do marcador, não o gate de decisões em aberto.
	orig := checkersWithGraph["open-questions-resolved"]
	checkersWithGraph["open-questions-resolved"] = func(_ string, n mapx.Node, _ string,
		_ *mapx.Graph, _ *config.Config) (Verdict, string) {
		if n.ID == comDecisao {
			return Pending, "há pergunta aberta " + OpenDecisionMarker
		}
		return Pending, "a seção nem existe — dívida de migração"
	}
	t.Cleanup(func() { checkersWithGraph["open-questions-resolved"] = orig })

	for _, f := range []string{comDecisao, semDecisao} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("texto\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := config.Gate{Name: "decisões", ID: "open-questions-resolved", On: []string{"spec"},
		Check: "open-questions-resolved"}
	res := RunWithConfig([]config.Gate{g}, []mapx.Node{
		{ID: comDecisao, Kind: mapx.KindSpec},
		{ID: semDecisao, Kind: mapx.KindSpec},
	}, root, nil, nil)

	porAlvo := map[string]Result{}
	for _, r := range res {
		porAlvo[r.Target] = r
	}
	if !porAlvo[comDecisao].Impede {
		t.Error("`há decisão por tomar` tem de IMPEDIR a promoção")
	}
	if porAlvo[semDecisao].Impede {
		t.Error("`não tive o que confrontar` NÃO pode impedir — foi o que reprovou 411 nós")
	}
	// E os dois usam o MESMO veredito: a distinção sai do campo, não do veredito.
	if porAlvo[comDecisao].Verdict != Pending || porAlvo[semDecisao].Verdict != Pending {
		t.Error("os dois casos são Pending; quem os separa é o campo, não o veredito")
	}
}

// DÍVIDA ASSUMIDA é o Pending que tem dono e vencimento — e o único que vira trabalho
// registrado. Os demais Pending são "não tive o que confrontar", que não é dívida de
// ninguém.
func TestSoOGateDeObrigacoesProduzDivida(t *testing.T) {
	t.Run("GTENG-B16: Only the obligations gate produces assumed debt", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "a.spec.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("texto\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	laudo := "[migracao-de-schema] pendente — DÍVIDA ASSUMIDA: até a virada do trimestre"

	devolvePending := func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string) {
		return Pending, laudo
	}
	origObr := checkersWithGraph["obligation-honored"]
	origOutro := checkersWithGraph["docs-fresh"]
	checkersWithGraph["obligation-honored"] = devolvePending
	checkersWithGraph["docs-fresh"] = devolvePending
	t.Cleanup(func() {
		checkersWithGraph["obligation-honored"] = origObr
		checkersWithGraph["docs-fresh"] = origOutro
	})

	obrig := runOne(config.Gate{Name: "obrigações", ID: "obligation-honored",
		On: []string{"spec"}, Check: "obligation-honored"},
		mapx.Node{ID: alvo, Kind: mapx.KindSpec}, root, nil, nil)
	outro := runOne(config.Gate{Name: "docs", ID: "docs-fresh",
		On: []string{"spec"}, Check: "docs-fresh"},
		mapx.Node{ID: alvo, Kind: mapx.KindSpec}, root, nil, nil)

	if !obrig.Divida {
		t.Error("o gate de obrigações produz DÍVIDA ASSUMIDA")
	}
	if obrig.Prazo == "" {
		t.Errorf("a dívida carrega o VENCIMENTO para o registro; veio vazio (laudo: %s)", obrig.Detail)
	}
	if outro.Divida {
		t.Error("outro gate com o MESMO laudo não pode virar dívida — o que decide é o gate")
	}
}

// A gramática do código de cenário segue o vocabulário do PROJETO, e o motor a
// reconfigura antes de rodar qualquer coisa. Um gate lendo o vocabulário velho reporta
// verde sobre o que nunca olhou.
func TestOMotorReconfiguraAGramaticaAntesDeRodar(t *testing.T) {
	t.Run("GTENG-B17: The engine reconfigures the code grammar before running anything", func(t *testing.T) {})

	t.Cleanup(func() { SetRuleLetters(config.DefaultRuleLetters) })
	SetRuleLetters(config.DefaultRuleLetters)

	const letra = "Z"
	if strings.Contains(config.DefaultRuleLetters, letra) {
		t.Fatalf("o teste precisa de uma letra FORA do vocabulário canônico")
	}
	root := t.TempDir()
	alvo := "a.spec.md"
	conteudo := "| `ABCDX-" + letra + "01` | uma regra da letra do projeto |\n"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: letra, Term: "algo que o projeto nomeia"},
	}}
	g := config.Gate{Name: "tem-código", ID: "has-code", On: []string{"spec"}, Check: "has-code"}

	res := RunWithConfig([]config.Gate{g}, []mapx.Node{{ID: alvo, Kind: mapx.KindSpec}}, root, nil, cfg)
	if len(res) != 1 {
		t.Fatalf("esperava um veredito, vieram %d", len(res))
	}
	if res[0].Verdict != Pass {
		t.Errorf("a letra declarada pelo projeto tem de estar viva quando o gate roda; "+
			"veio %v — %s", res[0].Verdict, res[0].Detail)
	}
}

// A porta simples: gates, nós e mapa, sem Estrutura. Os checkers relacionais leem a
// ausência como "sem de-para declarado" — e não como erro.
func TestPortaSimplesRodaSemEstrutura(t *testing.T) {
	t.Run("GTENG-B18: The plain entry point runs with the map and no Structure", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "a.spec.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("texto de verdade\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := config.Gate{Name: "não-vazio", ID: "non-empty", On: []string{"spec"}, Check: "non-empty"}
	res := Run([]config.Gate{g}, []mapx.Node{{ID: alvo, Kind: mapx.KindSpec}}, root, &mapx.Graph{})
	if len(res) != 1 || res[0].Verdict != Pass {
		t.Fatalf("a porta simples tem de rodar sem Estrutura; veio %+v", res)
	}
}

// A porta que CARREGA A ESTRUTURA entrega-a ao checker relacional: os regimes e as
// superfícies da trinca são declarados lá.
func TestPortaComEstruturaEntregaAEstrutura(t *testing.T) {
	t.Run("GTENG-B19: The entry point that carries the Structure hands it to the checkers", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "a.spec.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("texto\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const nome = "espia-a-estrutura"
	var recebeu bool
	checkersWithGraph[nome] = func(_ string, _ mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
		recebeu = cfg != nil
		return Pass, ""
	}
	t.Cleanup(func() { delete(checkersWithGraph, nome) })

	g := config.Gate{Name: nome, ID: nome, On: []string{"spec"}, Check: nome}
	nos := []mapx.Node{{ID: alvo, Kind: mapx.KindSpec}}

	RunWithConfig([]config.Gate{g}, nos, root, nil, &config.Config{Lang: "en"})
	if !recebeu {
		t.Error("a porta com Estrutura tem de entregá-la ao checker relacional")
	}
	recebeu = true
	Run([]config.Gate{g}, nos, root, nil)
	if recebeu {
		t.Error("a porta simples não inventa Estrutura: o checker recebe a ausência")
	}
}

// Só a porta que SABE se a varredura é o projeto inteiro pode honrar o escopo de full.
// No full, quem sabe varrer sozinho roda UMA vez sem receber a lista, em vez de receber
// os milhares de alvos em lotes.
func TestPortaQueSabeAVarreduraHonraOEscopoDeFull(t *testing.T) {
	t.Run("GTENG-B20: The entry point that knows the sweep kind honours the full-sweep scope", func(t *testing.T) {})

	g := config.Gate{Name: "lint", ID: "lint", On: []string{"code"},
		Scope: config.ScopeBatch, ScopeFull: config.ScopeProject, Run: "true"}
	alvos := []mapx.Node{{ID: "a.ts", Kind: mapx.KindCode}, {ID: "b.ts", Kind: mapx.KindCode}}
	root := t.TempDir()

	recorte := RunFull([]config.Gate{g}, alvos, root, nil, nil, false)
	completa := RunFull([]config.Gate{g}, alvos, root, nil, nil, true)

	if len(recorte) != 1 || len(completa) != 1 {
		t.Fatalf("o agregado responde uma vez em cada modo; vieram %d e %d", len(recorte), len(completa))
	}
	if recorte[0].Target == completa[0].Target {
		t.Errorf("o escopo tem de MUDAR entre o recorte e o full; os dois vieram %q", recorte[0].Target)
	}
	if !strings.Contains(recorte[0].Target, config.ScopeBatch) {
		t.Errorf("no recorte vale o `scope`; veio %q", recorte[0].Target)
	}
	if !strings.Contains(completa[0].Target, config.ScopeProject) {
		t.Errorf("no full vale o `scope_full`; veio %q", completa[0].Target)
	}
}

// A porta que honra a dispensa POR ALVO produz veredito para TODOS os nós — um deles
// poupado. A dispensa não pode sair filtrando o gate da lista.
func TestPortaComDispensaProduzVereditoParaTodos(t *testing.T) {
	t.Run("GTENG-B21: The entry point that honours a waiver by target keeps the gate running", func(t *testing.T) {})

	root := t.TempDir()
	nomes := []string{"um.md", "dois.md", "tres.md"}
	for _, f := range nomes {
		if err := os.WriteFile(filepath.Join(root, f), []byte("   \n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	disp, _ := ParseWaiver("non-empty@AAAAA=rascunho declarado")
	g := config.Gate{Name: "não-vazio", ID: "non-empty", On: []string{"spec"}, Check: "non-empty"}
	res := RunWithWaiver([]config.Gate{g}, []mapx.Node{
		{ID: nomes[0], Kind: mapx.KindSpec, Code: "AAAAA"},
		{ID: nomes[1], Kind: mapx.KindSpec, Code: "BBBBB"},
		{ID: nomes[2], Kind: mapx.KindSpec, Code: "CCCCC"},
	}, root, nil, nil, false, disp)

	if len(res) != len(nomes) {
		t.Fatalf("o gate tem de produzir veredito para os %d nós; vieram %d", len(nomes), len(res))
	}
	var poupados int
	for _, r := range res {
		if r.Verdict == Skip {
			poupados++
		}
	}
	if poupados != 1 {
		t.Errorf("exatamente um nó foi dispensado; vieram %d Skip", poupados)
	}
}

// A dispensa por alvo é aplicada COM O GATE RODANDO, nunca removendo-o da lista.
// Removê-lo apagaria a régua para o repositório inteiro, e um defeito noutro lugar
// passaria de carona.
func TestDispensaNaoRemoveOGateDaLista(t *testing.T) {
	t.Run("GTENG-I01: A waiver by target never removes the gate from the list", func(t *testing.T) {})

	root := t.TempDir()
	for _, f := range []string{"dispensado.md", "quebrado.md"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("   \n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	disp, _ := ParseWaiver("non-empty@DISPE=spec nova do plano 0007")
	g := config.Gate{Name: "não-vazio", ID: "non-empty", On: []string{"spec"}, Check: "non-empty"}
	res := RunWithWaiver([]config.Gate{g}, []mapx.Node{
		{ID: "dispensado.md", Kind: mapx.KindSpec, Code: "DISPE"},
		{ID: "quebrado.md", Kind: mapx.KindSpec, Code: "QBRDA"},
	}, root, nil, nil, false, disp)

	var achouReprovacao bool
	for _, r := range res {
		if r.Target == "quebrado.md" && r.Verdict == Fail {
			achouReprovacao = true
		}
	}
	if !achouReprovacao {
		t.Fatalf("o defeito NÃO dispensado tem de ser acusado — é o mascaramento que a "+
			"dispensa por alvo existe para impedir; vieram %+v", res)
	}
}

// Skip, Pending e Fail são TRÊS respostas, e nunca colapsam em duas. Cada uma diz o que
// as outras não dizem: não se aplica, não consegui medir, medi e reprovou.
func TestSkipPendingEFailNaoColapsam(t *testing.T) {
	t.Run("GTENG-I02: Stepping aside, not measuring and failing are three different answers", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "vazio.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nos := []mapx.Node{{ID: alvo, Kind: mapx.KindSpec}}

	semFerramenta := RunWithConfig([]config.Gate{{Name: "ext", ID: "ext", On: []string{"spec"},
		Run: "nada", NeedsTool: "binario-que-ninguem-instalou-jamais"}}, nos, root, nil, nil)
	mediuEReprovou := RunWithConfig([]config.Gate{{Name: "não-vazio", ID: "non-empty",
		On: []string{"spec"}, Check: "non-empty"}}, nos, root, nil, nil)
	naoMediu := RunWithConfig([]config.Gate{{Name: "mudo", ID: "mudo", On: []string{"spec"}}},
		nos, root, nil, nil)

	if len(semFerramenta) != 1 || len(mediuEReprovou) != 1 || len(naoMediu) != 1 {
		t.Fatal("cada gate deveria ter produzido um veredito")
	}
	tres := []Verdict{semFerramenta[0].Verdict, mediuEReprovou[0].Verdict, naoMediu[0].Verdict}
	vistos := map[Verdict]bool{}
	for _, v := range tres {
		if vistos[v] {
			t.Fatalf("os três casos colapsaram: %v", tres)
		}
		vistos[v] = true
	}
	if tres[0] != Skip || tres[1] != Fail || tres[2] != Pending {
		t.Errorf("esperava Skip/Fail/Pending nessa ordem, veio %v", tres)
	}
}

// O veredito do alvo dispensado é Skip COM O MOTIVO ESCRITO, nunca silêncio. Some do
// placar de reprovações sem sumir do relatório — a diferença entre dispensar e esconder.
func TestOAlvoDispensadoCarregaOMotivo(t *testing.T) {
	t.Run("GTENG-I03: A waived target leaves the failure tally without leaving the report", func(t *testing.T) {})

	const motivo = "rascunho declarado do plano 0007"
	root := t.TempDir()
	alvo := "poupado.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	disp, _ := ParseWaiver("non-empty@POUPA=" + motivo)
	res := RunWithWaiver([]config.Gate{{Name: "não-vazio", ID: "non-empty",
		On: []string{"spec"}, Check: "non-empty"}},
		[]mapx.Node{{ID: alvo, Kind: mapx.KindSpec, Code: "POUPA"}},
		root, nil, nil, false, disp)

	if len(res) != 1 {
		t.Fatalf("esperava um veredito, vieram %d", len(res))
	}
	if res[0].Verdict != Skip {
		t.Fatalf("o dispensado é Skip, veio %v", res[0].Verdict)
	}
	if !strings.Contains(res[0].Detail, motivo) {
		t.Errorf("o laudo tem de carregar o MOTIVO escrito, e não calar; veio %q", res[0].Detail)
	}
}

// O alvo reportado pelo agregado é o ESCOPO. Culpar um de muitos arquivos por uma
// resposta sobre o conjunto seria afirmação que o motor não sustenta.
func TestAlvoDoAgregadoEhOEscopo(t *testing.T) {
	t.Run("GTENG-I04: The reported target of an aggregate gate is the scope", func(t *testing.T) {})

	g := config.Gate{Name: "lint", ID: "lint", On: []string{"code"},
		Scope: config.ScopeBatch, Run: "true"}
	alvos := []mapx.Node{
		{ID: "a.ts", Kind: mapx.KindCode},
		{ID: "b.ts", Kind: mapx.KindCode},
	}
	res := RunWithConfig([]config.Gate{g}, alvos, t.TempDir(), nil, nil)
	if len(res) != 1 {
		t.Fatalf("esperava um veredito agregado, vieram %d", len(res))
	}
	if !strings.Contains(res[0].Target, config.ScopeBatch) {
		t.Errorf("o alvo tem de nomear o escopo; veio %q", res[0].Target)
	}
	for _, n := range alvos {
		if strings.Contains(res[0].Target, n.ID) {
			t.Errorf("o alvo não pode nomear um dos arquivos; veio %q", res[0].Target)
		}
	}
}

// A régua é do checker. Um motor que julgasse poria a medição dentro do roteador, e o
// mesmo defeito teria dois donos.
func TestOMotorNaoDecideSeOAlvoEstaCerto(t *testing.T) {
	t.Run("GTENG-X01: The engine does not decide whether a target is correct", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "absurdo.md"
	// Conteúdo que qualquer revisor reprovaria; o checker aprova, e é ele quem manda.
	if err := os.WriteFile(filepath.Join(root, alvo),
		[]byte("esta spec diz exatamente o oposto do que a unidade faz\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const nome = "aprova-sempre"
	checkersWithGraph[nome] = func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string) {
		return Pass, ""
	}
	t.Cleanup(func() { delete(checkersWithGraph, nome) })

	res := RunWithConfig([]config.Gate{{Name: nome, ID: nome, On: []string{"spec"}, Check: nome}},
		[]mapx.Node{{ID: alvo, Kind: mapx.KindSpec}}, root, nil, nil)
	if len(res) != 1 || res[0].Verdict != Pass {
		t.Fatalf("o veredito é o do CHECKER; veio %+v", res)
	}
}

// O motor NÃO computa o veredito de um gate de julgamento. Tudo o que ele sabe é se
// ALGUÉM já respondeu — e a resposta está no carimbo.
func TestOMotorNaoComputaJulgamento(t *testing.T) {
	t.Run("GTENG-X02: The engine does not compute the verdict of a judgment gate", func(t *testing.T) {})

	g := config.Gate{Name: "atomic", On: []string{"code"}, Measures: config.MeasuresJudgment}
	alvo := mapx.Node{ID: "Home.tsx", Kind: mapx.KindCode, Rev: "r1"}

	// Um grafo COM arestas, mas SEM carimbo deste gate: não há resposta a ler, e o motor
	// não inventa uma por conta própria.
	graph := grafoComJulgamento(alvo, "outro-gate", "ok")
	r := runOne(g, alvo, t.TempDir(), graph, nil)
	if r.Verdict != Judge {
		t.Fatalf("sem carimbo DESTE gate, o motor só pode pedir julgamento; veio %v (%s)",
			r.Verdict, r.Detail)
	}
}

// Ausência é CASO, não erro. Fabricar mapa, Estrutura ou dispensa faria a rodada
// responder sobre um estado de projeto que não existe.
func TestOMotorNaoInventaMapaEstruturaNemDispensa(t *testing.T) {
	t.Run("GTENG-X03: The engine invents neither a map nor a Structure nor a waiver", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "vazio.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const nome = "espia-a-ausencia"
	var viuGrafo, viuEstrutura bool
	checkersWithGraph[nome] = func(_ string, _ mapx.Node, _ string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
		viuGrafo, viuEstrutura = g != nil, cfg != nil
		return Pass, ""
	}
	t.Cleanup(func() { delete(checkersWithGraph, nome) })

	// Sem mapa, sem Estrutura e sem dispensa: o checker recebe a ausência tal como é.
	RunWithWaiver([]config.Gate{{Name: nome, ID: nome, On: []string{"spec"}, Check: nome}},
		[]mapx.Node{{ID: alvo, Kind: mapx.KindSpec, Code: "AUSEN"}},
		root, nil, nil, false, Waiver{})
	if viuGrafo || viuEstrutura {
		t.Errorf("o motor não pode fabricar o que não recebeu: grafo=%v estrutura=%v",
			viuGrafo, viuEstrutura)
	}

	// E a dispensa vazia não dispensa ninguém: o alvo quebrado continua acusado.
	res := RunWithWaiver([]config.Gate{{Name: "não-vazio", ID: "non-empty",
		On: []string{"spec"}, Check: "non-empty"}},
		[]mapx.Node{{ID: alvo, Kind: mapx.KindSpec, Code: "AUSEN"}},
		root, nil, nil, false, Waiver{})
	if len(res) != 1 || res[0].Verdict != Fail {
		t.Errorf("dispensa vazia não pode poupar ninguém; veio %+v", res)
	}
}

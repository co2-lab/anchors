package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const flagCompleta = "| Cenário | Quando o valor | Então |\n" +
	"| --- | --- | --- |\n" +
	"| `CHKUT-G01` | `= \"off\"` | o checkout antigo responde |\n" +
	"| `CHKUT-G02` | `= \"on\"` | o checkout novo responde |\n" +
	"| `CHKUT-G03` | ausente | vale o default do cabeçalho |\n"

func flagNode() mapx.Node {
	return mapx.Node{ID: "flags/checkout.flag.md", Kind: mapx.KindFlag, Code: "CHKUT"}
}

// --- flag-scenarios-complete ---

// O caso AUSENTE é o que quebra em produção e o que ninguém escreve.
func TestFlagScenariosComplete(t *testing.T) {
	v, _ := checkFlagScenariosComplete(flagCompleta, flagNode(), "", nil, nil)
	if v != Pass {
		t.Errorf("a flag declara `ausente` e veio %v", v)
	}

	semAusente := "| `CHKUT-G01` | `= \"on\"` | liga |\n"
	v, msg := checkFlagScenariosComplete(semAusente, flagNode(), "", nil, nil)
	if v != Fail {
		t.Errorf("a flag NÃO declara o ausente e veio %v", v)
	}
	if !strings.Contains(msg, "@no-absent") {
		t.Errorf("o veredito não nomeia a dispensa: %q", msg)
	}
}

// A dispensa é `@no-absent` com razão — marcador nu não vale, pela mesma regra dos outros
// opt-outs: a razão é o que responde a pergunta daqui a seis meses.
func TestFlagScenariosComplete_dispensa(t *testing.T) {
	semAusente := "| `CHKUT-G01` | `= \"on\"` | liga |\n"

	comRazao := semAusente + "\n@no-absent: lida de constante local, nunca falta\n"
	if v, _ := checkFlagScenariosComplete(comRazao, flagNode(), "", nil, nil); v != Pass {
		t.Errorf("a dispensa tem razão escrita e veio %v", v)
	}

	nu := semAusente + "\n@no-absent\n"
	if v, _ := checkFlagScenariosComplete(nu, flagNode(), "", nil, nil); v != Fail {
		t.Error("o marcador NU foi aceito — a razão é obrigatória")
	}
}

func TestFlagScenariosComplete_skips(t *testing.T) {
	naoEFlag := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFlagScenariosComplete(flagCompleta, naoEFlag, "", nil, nil); v != Skip {
		t.Errorf("não é flag e veio %v", v)
	}
	if v, _ := checkFlagScenariosComplete("# sem tabela\n", flagNode(), "", nil, nil); v != Skip {
		t.Error("flag sem cenário nenhum deveria dar Skip")
	}
}

// --- flag-scenario-grammar ---

// A gramática é fixa para RECUSAR. Prosa passa por qualquer régua e não confronta nada.
func TestFlagScenarioGrammar(t *testing.T) {
	if v, _ := checkFlagScenarioGrammar(flagCompleta, flagNode(), "", nil, nil); v != Pass {
		t.Errorf("todas as condições estão na gramática e veio %v", v)
	}

	prosa := "| `CHKUT-G01` | quando o usuário for beta | liga |\n"
	v, msg := checkFlagScenarioGrammar(prosa, flagNode(), "", nil, nil)
	if v != Fail {
		t.Errorf("a condição é prosa e veio %v", v)
	}
	if !strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("o veredito não diz QUAL cenário: %q", msg)
	}
}

// --- flag-scenario-exists ---

func projetoComFlag(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "flags"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, "flags", "checkout.flag.md")
	if err := os.WriteFile(p, []byte(flagCompleta), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestFlagScenarioExists(t *testing.T) {
	root := projetoComFlag(t)
	spec := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}

	existe := "### CRED-V01 — valida o limite   @gated-by CHKUT-G02\n"
	if v, _ := checkFlagScenarioExists(existe, spec, root, nil, nil); v != Pass {
		t.Errorf("o cenário existe e veio %v", v)
	}

	naoExiste := "### CRED-V01 — valida o limite   @gated-by CHKUT-G99\n"
	v, msg := checkFlagScenarioExists(naoExiste, spec, root, nil, nil)
	if v != Fail {
		t.Errorf("o cenário NÃO existe e veio %v", v)
	}
	if !strings.Contains(msg, "CHKUT-G99") {
		t.Errorf("o veredito não traz o código: %q", msg)
	}

	if v, _ := checkFlagScenarioExists("### CRED-V01 — sem citação\n", spec, root, nil, nil); v != Skip {
		t.Error("spec sem citação deveria dar Skip")
	}
}

// --- flag-covered ---

// Por CENÁRIO, não por flag: é justamente o ramo desligado que fica sem prova numa régua
// por flag, e é ele que vai quebrar.
func TestFlagCovered_porCenarioNaoPorFlag(t *testing.T) {
	// Um teste prova só o G01. Numa régua por FLAG isso passaria.
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	v, msg := checkFlagCovered(flagCompleta, n, "", g, nil)
	if v != Fail {
		t.Fatalf("dois cenários sem teste e veio %v", v)
	}
	for _, c := range []string{"CHKUT-G02", "CHKUT-G03"} {
		if !strings.Contains(msg, c) {
			t.Errorf("o veredito não nomeia %s: %q", c, msg)
		}
	}
	if strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("o veredito acusa o cenário que TEM teste: %q", msg)
	}
}

func TestFlagCovered_todosProvados(t *testing.T) {
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, msg := checkFlagCovered(flagCompleta, n, "", g, nil); v != Pass {
		t.Errorf("todos os cenários têm teste e veio %v: %s", v, msg)
	}
}

// Sem mapa não há como saber o que foi provado — e Pass aqui seria mentira.
func TestFlagCovered_semMapaNaoAfirmaPass(t *testing.T) {
	if v, _ := checkFlagCovered(flagCompleta, flagNode(), "", nil, nil); v == Pass {
		t.Error("sem mapa o gate afirmou Pass — não tinha como saber")
	}
}

// AS DUAS PERGUNTAS, e por que separá-las.
//
// "Tem teste escrito?" é estática e sempre respondível. "O teste passou?" exige execução
// ingerida. Juntar as duas perde a resposta das duas: num projeto que nunca ingeriu
// relatório (o caso deste repositório — zero `proven_codes` no grafo), a régua só de
// execução acusava TODO cenário, inclusive os que tinham teste escrito e passando.
//
// E só a estática seria o erro oposto, pior: teste escrito pode nunca ter rodado.
func TestFlagCovered_semTesteNenhumAcusaMesmoSemIngestao(t *testing.T) {
	root := t.TempDir()
	// Um nó de teste que existe e NÃO nomeia cenário nenhum.
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte("func TestNada(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}

	v, msg := checkFlagCovered(flagCompleta, flagNode(), root, g, nil)
	if v != Fail {
		t.Fatalf("nenhum teste nomeia os cenários — esperava Fail, veio %v", v)
	}
	for _, c := range []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"} {
		if !strings.Contains(msg, c) {
			t.Errorf("o veredito não nomeia %s: %q", c, msg)
		}
	}
}

// ESCRITO mas NUNCA RODADO: o gate tem de dizer isso, e não "sem teste".
//
// É o caso que o pedido nomeou: o teste pode ter sido desenvolvido e nunca ter colhido
// resultado. O conserto é outro — rodar a suíte, não escrever um teste.
func TestFlagCovered_escritoMasSemExecucaoDizQualDosDois(t *testing.T) {
	root := t.TempDir()
	corpo := "func TestX(t *testing.T) { /* CHKUT-G01 */ }\nconst c = \"CHKUT-G01\"\n"
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte(corpo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}

	v, msg := checkFlagCovered(flagCompleta, flagNode(), root, g, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	// O G01 tem teste ESCRITO: não pode aparecer como "nenhum teste nomeia".
	if !strings.Contains(msg, "ingest") {
		t.Errorf("o veredito não distingue escrito-sem-execução: %q", msg)
	}
	if !strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("o veredito não nomeia o cenário escrito e não executado: %q", msg)
	}
}

// O CÓDIGO EM COMENTÁRIO não conta como teste escrito — mesma régua do
// `feature-test-match`: citação é referência, não implementação.
func TestFlagCovered_codigoEmComentarioNaoContaComoTeste(t *testing.T) {
	root := t.TempDir()
	soComentario := "// CHKUT-G01 é tratado noutro lugar\nfunc TestX(t *testing.T) {}\n"
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte(soComentario), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}

	_, msg := checkFlagCovered(flagCompleta, flagNode(), root, g, nil)
	if strings.Contains(msg, "ingest") {
		t.Errorf("o código só em COMENTÁRIO passou por teste escrito: %q", msg)
	}
}

// E com execução ingerida, o cenário provado sai da acusação.
func TestFlagCovered_ingeridoEVerdePassa(t *testing.T) {
	root := t.TempDir()
	// O sinal vive no no' DA FLAG, como o do `scenario-coverage` vive no da spec: a
	// ingestao cruza os provados com os que o no' DECLARA.
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, msg := checkFlagCovered(flagCompleta, n, root, g, nil); v != Pass {
		t.Errorf("todos provados e veio %v: %s", v, msg)
	}
}

// --- flag-scenario-governs: A VOLTA do `flag-scenario-exists` ---
//
// Aquele confronta a spec que cita um cenário inexistente; este, o cenário que ninguém
// cita. Um cenário que nenhuma regra invoca é caminho DECLARADO e não governado.

// A aresta `gated-by` CHEGA na flag e carrega, no Method, o código do cenário invocado.
func citando(codigos ...string) *mapx.Graph {
	g := &mapx.Graph{Nodes: []mapx.Node{flagNode()}}
	for _, c := range codigos {
		g.Edges = append(g.Edges, mapx.Edge{
			From: "a.spec.md", To: "flags/checkout.flag.md",
			Type: mapx.EdgeGatedBy, Method: c,
		})
	}
	return g
}

func TestFlagScenarioGoverns_cenarioSemRegraEAcusado(t *testing.T) {
	t.Run("CHKUT-G0X: a scenario no rule invokes is reported", func(t *testing.T) {})
	// Só o G01 é citado; G02 e G03 ficam soltos.
	v, msg := checkFlagScenarioGoverns(flagCompleta, flagNode(), "", citando("CHKUT-G01"), nil)
	if v != Fail {
		t.Fatalf("dois cenários sem regra e veio %v", v)
	}
	for _, c := range []string{"CHKUT-G02", "CHKUT-G03"} {
		if !strings.Contains(msg, c) {
			t.Errorf("o veredito não nomeia %s: %q", c, msg)
		}
	}
	if strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("acusou o cenário que É citado: %q", msg)
	}
}

func TestFlagScenarioGoverns_todosCitadosPassa(t *testing.T) {
	g := citando("CHKUT-G01", "CHKUT-G02", "CHKUT-G03")
	if v, msg := checkFlagScenarioGoverns(flagCompleta, flagNode(), "", g, nil); v != Pass {
		t.Errorf("todos citados e veio %v: %s", v, msg)
	}
}

// A dispensa por cenário: há caminho que existe e nenhuma regra precisa nomear — o `off`
// que devolve ao comportamento antigo, já governado pelas regras que sempre valeram.
func TestFlagScenarioGoverns_dispensaPorCenario(t *testing.T) {
	comDispensa := "| Cenário | Quando o valor | Então |\n| --- | --- | --- |\n" +
		"| `CHKUT-G01` | `= \"off\"` | vale o comportamento antigo @no-govern: as regras de sempre já o governam |\n"
	if v, msg := checkFlagScenarioGoverns(comDispensa, flagNode(), "", citando(), nil); v != Pass {
		t.Errorf("a dispensa tem razão escrita e veio %v: %s", v, msg)
	}

	nu := strings.Replace(comDispensa, "@no-govern: as regras de sempre já o governam", "@no-govern", 1)
	if v, _ := checkFlagScenarioGoverns(nu, flagNode(), "", citando(), nil); v != Fail {
		t.Error("o marcador NU foi aceito — a razão é obrigatória")
	}
}

func TestFlagScenarioGoverns_skips(t *testing.T) {
	naoEFlag := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFlagScenarioGoverns(flagCompleta, naoEFlag, "", citando(), nil); v != Skip {
		t.Errorf("não é flag e veio %v", v)
	}
	if v, _ := checkFlagScenarioGoverns(flagCompleta, flagNode(), "", nil, nil); v == Pass {
		t.Error("sem mapa o gate afirmou Pass — não tinha como saber quem cita")
	}
}

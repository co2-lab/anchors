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
	g := &mapx.Graph{Nodes: []mapx.Node{{
		ID: "t.test.ts", Kind: mapx.KindTest,
		Signal: &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01"}},
	}}}
	v, msg := checkFlagCovered(flagCompleta, flagNode(), "", g, nil)
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
	g := &mapx.Graph{Nodes: []mapx.Node{{
		ID: "t.test.ts", Kind: mapx.KindTest,
		Signal: &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"}},
	}}}
	if v, msg := checkFlagCovered(flagCompleta, flagNode(), "", g, nil); v != Pass {
		t.Errorf("todos os cenários têm teste e veio %v: %s", v, msg)
	}
}

// Sem mapa não há como saber o que foi provado — e Pass aqui seria mentira.
func TestFlagCovered_semMapaNaoAfirmaPass(t *testing.T) {
	if v, _ := checkFlagCovered(flagCompleta, flagNode(), "", nil, nil); v == Pass {
		t.Error("sem mapa o gate afirmou Pass — não tinha como saber")
	}
}

// NADA INGERIDO não é "sem teste" — e confundir os dois é acusação FALSA.
//
// Medido neste próprio repositório: zero `proven_codes` no grafo inteiro, e a primeira
// versão do gate acusou três cenários cujos testes ela não tinha como enxergar. É a
// mesma distinção que o `scenario-coverage` já faz, e o comentário dele registra o custo
// de errá-la: "o gate pedia o impossível, e a mensagem sugeria que a spec estava mal
// coberta".
func TestFlagCovered_semIngestaoEPendingNaoFail(t *testing.T) {
	// Um nó de teste EXISTE e nada foi ingerido: Signal == nil.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}
	v, msg := checkFlagCovered(flagCompleta, flagNode(), "", g, nil)
	if v == Fail {
		t.Errorf("sem ingestão o gate ACUSOU (%q) — não tinha como saber se há teste", msg)
	}
	if v != Pending {
		t.Errorf("sem ingestão esperava Pending, veio %v", v)
	}
}

// E com sinal ingerido ele volta a acusar de verdade: a distinção não pode virar
// desculpa para nunca cobrar nada.
func TestFlagCovered_comIngestaoVoltaACobrar(t *testing.T) {
	g := &mapx.Graph{Nodes: []mapx.Node{{
		ID: "t_test.go", Kind: mapx.KindTest,
		Signal: &mapx.TestSignal{ProvenCodes: []string{"OUTRO-G01"}},
	}}}
	if v, _ := checkFlagCovered(flagCompleta, flagNode(), "", g, nil); v != Fail {
		t.Errorf("houve ingestão e nenhum cenário provado — esperava Fail, veio %v", v)
	}
}

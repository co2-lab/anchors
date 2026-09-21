package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// AS DUAS PERGUNTAS do `scenario-coverage` — a mesma régua do `flag-covered`.
//
// A versão anterior só perguntava "passou?", e num projeto que nunca ingeriu relatório
// respondia Pending para tudo: "ninguém mediu", escondendo os cenários que ninguém
// testou. E só a estática seria o erro oposto — teste escrito pode nunca ter rodado.
//
// O conserto de cada estado é diferente, e é por isso que o veredito tem de separá-los:
// "sem teste" pede que alguém escreva um; "escrito e não executado" pede que alguém rode.

const specComDoisRequisitos = "### CREDX-B01 — valida o limite\n\n### CREDX-B02 — recusa o saldo\n"

func specNodeCobertura() mapx.Node {
	return mapx.Node{ID: "credx.spec.md", Kind: mapx.KindSpec, Code: "CREDX"}
}

func rootComTeste(t *testing.T, corpo string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "credx_test.go"), []byte(corpo), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{Nodes: []mapx.Node{{ID: "credx_test.go", Kind: mapx.KindTest}}}
}

// SEM TESTE NENHUM: acusa, mesmo sem execução ingerida. É o caso que a versão anterior
// escondia atrás de um Pending.
func TestScenarioCoverage_semTesteAcusaMesmoSemIngestao(t *testing.T) {
	root, g := rootComTeste(t, "func TestNada(t *testing.T) {}\n")

	v, msg := checkScenarioCoverage(specComDoisRequisitos, specNodeCobertura(), root, g, nil)
	if v != Fail {
		t.Fatalf("nenhum teste nomeia os requisitos — esperava Fail, veio %v", v)
	}
	for _, c := range []string{"CREDX-B01", "CREDX-B02"} {
		if !strings.Contains(msg, c) {
			t.Errorf("o veredito não nomeia %s: %q", c, msg)
		}
	}
}

// ESCRITO E NUNCA RODADO: o veredito tem de dizer qual dos dois problemas é.
func TestScenarioCoverage_escritoSemExecucaoDizQualDosDois(t *testing.T) {
	root, g := rootComTeste(t, "const c = \"CREDX-B01\"\nfunc TestX(t *testing.T) {}\n")

	v, msg := checkScenarioCoverage(specComDoisRequisitos, specNodeCobertura(), root, g, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, "ingest") {
		t.Errorf("o veredito não distingue escrito-sem-execução: %q", msg)
	}
	// O B02 continua sem teste nenhum, e os dois estados aparecem no mesmo veredito.
	if !strings.Contains(msg, "CREDX-B02") {
		t.Errorf("o veredito perdeu o requisito sem teste: %q", msg)
	}
}

// CÓDIGO EM COMENTÁRIO não conta — citação é referência, não implementação.
func TestScenarioCoverage_comentarioNaoContaComoTeste(t *testing.T) {
	root, g := rootComTeste(t, "// CREDX-B01 é coberto noutro lugar\nfunc TestX(t *testing.T) {}\n")

	_, msg := checkScenarioCoverage(specComDoisRequisitos, specNodeCobertura(), root, g, nil)
	if strings.Contains(msg, "ingest") {
		t.Errorf("o código só em COMENTÁRIO passou por teste escrito: %q", msg)
	}
}

// E o provado por EXECUÇÃO sai da acusação, que é o comportamento original.
func TestScenarioCoverage_provadoPassa(t *testing.T) {
	root, _ := rootComTeste(t, "func TestNada(t *testing.T) {}\n")
	n := specNodeCobertura()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CREDX-B01", "CREDX-B02"}}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "credx_test.go", Kind: mapx.KindTest}}}

	if v, msg := checkScenarioCoverage(specComDoisRequisitos, n, root, g, nil); v != Pass {
		t.Errorf("os dois requisitos provados e veio %v: %s", v, msg)
	}
}

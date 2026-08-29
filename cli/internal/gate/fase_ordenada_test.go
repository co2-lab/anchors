package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A ordem das fases vivia em prosa ("Fase 3 — depende da Fase 2"), e prosa não é
// confrontável: as specs de um plano nasciam todas disponíveis, e o agente pegava a da
// fase 3 com a fase 1 em aberto. Medido no primeiro uso real.
func TestFaseOrdenada(t *testing.T) {
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}

	ok := "### FNDTN-F01 — a árvore\n\n### FNDTN-F02 — a régua (depende de FNDTN-F01)\n"
	if v, msg := checkFaseOrdenada(ok, plano, "", nil, nil); v != Pass {
		t.Errorf("ordem válida deveria passar: %v — %s", v, msg)
	}

	// Depender do que vem DEPOIS é impossível de cumprir.
	invertida := "### FNDTN-F01 — a árvore (depende de FNDTN-F02)\n\n### FNDTN-F02 — a régua\n"
	v, msg := checkFaseOrdenada(invertida, plano, "", nil, nil)
	if v != Fail || !strings.Contains(msg, "DEPOIS") {
		t.Errorf("depender do que vem depois deveria reprovar: %v — %s", v, msg)
	}

	// Fase que não existe no plano.
	fantasma := "### FNDTN-F01 — a árvore (depende de FNDTN-F09)\n"
	if v, msg := checkFaseOrdenada(fantasma, plano, "", nil, nil); v != Fail || !strings.Contains(msg, "não está catalogada") {
		t.Errorf("fase inexistente deveria reprovar: %v — %s", v, msg)
	}

	// Código repetido: duas fases com o mesmo código tornam impossível dizer de qual uma
	// spec depende.
	repetida := "### FNDTN-F01 — a árvore\n\n### FNDTN-F01 — outra\n"
	if v, _ := checkFaseOrdenada(repetida, plano, "", nil, nil); v != Fail {
		t.Errorf("código repetido deveria reprovar, veio %v", v)
	}

	// PROSA sem código: pendência, não falha — é dívida de quem quiser a ordem
	// confrontável, e não erro de quem escreveu um plano pequeno.
	prosa := "## Fases\n\n### Fase 1 — a árvore\n\n### Fase 2 — depende da Fase 1\n"
	if v, msg := checkFaseOrdenada(prosa, plano, "", nil, nil); v != Pending || !strings.Contains(msg, "não cataloga") {
		t.Errorf("fase em prosa é pendência: %v — %s", v, msg)
	}

	// Plano sem fase nenhuma não está errado: plano pequeno não precisa de fase.
	if v, _ := checkFaseOrdenada("## Objetivo\n\nTexto.\n", plano, "", nil, nil); v != Skip {
		t.Errorf("plano sem fases não deveria ser cobrado, veio %v", v)
	}
}

// Uma spec pode depender de MAIS DE UMA fase — a de teste precisa da árvore E da régua
// de tipos, não de uma delas. O `needs:` aceita lista separada por vírgula, e cada item é
// confrontado por si: apontar duas fases e ter uma errada não pode passar por acaso.
func TestFaseExisteAceitaVarias(t *testing.T) {
	planoID := "plans/0001.md"
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: planoID, Kind: mapx.KindPlan}}}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plano := "### FNDTN-F01 — a árvore\n\n### FNDTN-F02 — a régua\n\n### FNDTN-F03 — o teste\n"
	if err := os.WriteFile(filepath.Join(dir, planoID), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}

	// DUAS fases, ambas existentes.
	spec := mapx.Node{Kind: mapx.KindSpec, Needs: []string{"FNDTN-F01", "FNDTN-F02"}}
	if v, msg := checkFaseExiste("", spec, dir, g, nil); v != Pass {
		t.Errorf("duas fases existentes deveriam passar: %v — %s", v, msg)
	}

	// Uma existe e a outra não: reprova, e NOMEIA só a que falta — dizer "alguma está
	// errada" obrigaria a conferir as duas à mão.
	meio := mapx.Node{Kind: mapx.KindSpec, Needs: []string{"FNDTN-F01", "FNDTN-F09"}}
	v, msg := checkFaseExiste("", meio, dir, g, nil)
	if v != Fail {
		t.Fatalf("uma fase inexistente entre duas deveria reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "FNDTN-F09") {
		t.Errorf("a mensagem deveria nomear a fase que falta: %s", msg)
	}
	if strings.Contains(msg, "FNDTN-F01") {
		t.Errorf("a fase que EXISTE não pode aparecer como ausente: %s", msg)
	}
}

// `parent` é PERTENCIMENTO e `needs` é ORDEM. Um pai que não existe não falha em lugar
// nenhum — o item SOME da árvore, e é por isso que o gate precisa dizer.
func TestParentValido(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plano := "### FNDTN-F01 — a árvore\n\n### FNDTN-F02 — a régua\n"
	if err := os.WriteFile(filepath.Join(dir, "plans/0001.md"), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "FNDTN"},
		{ID: "a.spec.md", Kind: mapx.KindSpec, Code: "WRKSP", Parent: "FNDTN-F01"},
	}}

	// Pai que É uma fase catalogada.
	ok := mapx.Node{Kind: mapx.KindSpec, Code: "WRKSP", Parent: "FNDTN-F01"}
	if v, msg := checkParentValido("", ok, dir, g, nil); v != Pass {
		t.Errorf("fase catalogada é pai válido: %v — %s", v, msg)
	}

	// Pai que é o CÓDIGO de um artefato.
	if v, msg := checkParentValido("", mapx.Node{Code: "X", Parent: "FNDTN"}, dir, g, nil); v != Pass {
		t.Errorf("artefato do mapa é pai válido: %v — %s", v, msg)
	}

	// Pai INEXISTENTE: o item sumiria da árvore em silêncio.
	v, msg := checkParentValido("", mapx.Node{Code: "X", Parent: "FNDTN-F09"}, dir, g, nil)
	if v != Fail {
		t.Fatalf("pai inexistente deveria reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "SOME da árvore") {
		t.Errorf("a mensagem deveria dizer o que acontece na prática: %s", msg)
	}

	// Pai de si mesmo: quem monta a árvore entra em laço.
	if v, _ := checkParentValido("", mapx.Node{Code: "WRKSP", Parent: "WRKSP"}, dir, g, nil); v != Fail {
		t.Errorf("ser pai de si mesmo deveria reprovar, veio %v", v)
	}

	// Sem `parent` não há o que confrontar — e isso não é falha: a maioria dos artefatos
	// não pertence a nada.
	if v, _ := checkParentValido("", mapx.Node{Code: "X"}, dir, g, nil); v != Skip {
		t.Errorf("sem parent deveria pular, veio %v", v)
	}
}

// CICLO na cadeia: A contém B contém A. Quem percorre nunca chega à raiz.
func TestParentSemCiclo(t *testing.T) {
	dir := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.md", Code: "AAAAA", Parent: "BBBBB"},
		{ID: "b.md", Code: "BBBBB", Parent: "AAAAA"},
	}}
	v, msg := checkParentValido("", g.Nodes[0], dir, g, nil)
	if v != Fail {
		t.Fatalf("ciclo deveria reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "ciclo") {
		t.Errorf("a mensagem deveria nomear o ciclo: %s", msg)
	}
}

package suggestion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func nova(id string) Suggestion {
	return Suggestion{
		ID: id, Gate: "mock-detect-cobre-o-dialeto", Target: "anchors.yaml",
		Origin: FromJudgment,
		Why:    "o padrão declarado não casa `unittest.mock.patch`, usado em 12 testes",
		Patch:  "--- a/anchors.yaml\n+++ b/anchors.yaml\n@@\n-  mock_detect: \"x\"\n+  mock_detect: \"y\"",
	}
}

func TestOpenGravaEmPending(t *testing.T) {
	root := t.TempDir()
	criada, p, err := Open(root, nova("s1"))
	if err != nil || !criada {
		t.Fatalf("deveria criar: %v %v", criada, err)
	}
	if !strings.Contains(p, filepath.Join("suggestions", "pending")) {
		t.Errorf("sugestão nasce em pending: %s", p)
	}
	b, _ := os.ReadFile(p)
	if !strings.Contains(string(b), "```diff") {
		t.Errorf("o patch deve estar no arquivo: %s", b)
	}
}

// Idempotente: a mesma proposta reaparecendo numa varredura seguinte não duplica.
func TestOpenNaoDuplica(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, nova("s1"))
	criada, _, err := Open(root, nova("s1"))
	if err != nil {
		t.Fatal(err)
	}
	if criada {
		t.Error("a mesma sugestão não pode ser criada duas vezes")
	}
}

// Proposta já RECUSADA não volta para pending. Reabrir apagaria a decisão de quem já
// olhou, e a mesma sugestão voltaria a cada varredura como se fosse nova.
func TestOpenNaoReabreDecidida(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, nova("s1"))
	if err := Decide(root, "s1", Rejected, "o dialeto é intencional aqui", false); err != nil {
		t.Fatal(err)
	}
	criada, _, _ := Open(root, nova("s1"))
	if criada {
		t.Error("sugestão já decidida não pode reabrir")
	}
	pend, _ := List(root, Pending)
	if len(pend) != 0 {
		t.Errorf("pending deveria estar vazio: %v", pend)
	}
}

// Mover é decidir — o estado é a pasta, como nas issues. Não há campo de status dentro
// do arquivo que possa divergir de onde ele está.
func TestDecideMoveERegistra(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, nova("s1"))
	if err := Decide(root, "s1", Approved, "o padrão estava mesmo errado", false); err != nil {
		t.Fatal(err)
	}
	ap, _ := List(root, Approved)
	if len(ap) != 1 || ap[0] != "s1" {
		t.Fatalf("deveria estar em approved: %v", ap)
	}
	b, _ := os.ReadFile(filepath.Join(root, Dir, string(Approved), "s1.md"))
	if !strings.Contains(string(b), "o padrão estava mesmo errado") {
		t.Errorf("a razão deve ficar registrada: %s", b)
	}
	if !strings.Contains(string(b), "por:** pessoa") {
		t.Errorf("quem decidiu deve ficar registrado: %s", b)
	}
}

// Decisão sem razão é a mesma falha do `@no-test` nu: some o rastro de por que alguém
// escolheu, e a escolha vira indistinguível de descuido.
func TestDecideExigeRazao(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, nova("s1"))
	if err := Decide(root, "s1", Approved, "   ", false); err == nil {
		t.Error("decisão sem razão deveria falhar")
	}
}

// Sob `auto_judgment` a decisão é da IA, e isso fica MARCADO. Apagar a distinção faria
// "ninguém olhou isto" parecer "alguém aprovou".
func TestDecideMarcaJulgamentoAutomatico(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, nova("s1"))
	if err := Decide(root, "s1", Approved, "padrão claramente incompleto", true); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(root, Dir, string(Approved), "s1.md"))
	if !strings.Contains(string(b), "IA (auto_judgment)") {
		t.Errorf("decisão automática deve ser distinguível da humana: %s", b)
	}
}

func TestPatchOfExtraiODiff(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, nova("s1"))
	p, err := PatchOf(root, "s1", Pending)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(p, "--- a/anchors.yaml") || strings.Contains(p, "```") {
		t.Errorf("o patch deve sair limpo, aplicável com git apply: %q", p)
	}
}

// Achado cuja correção NÃO é mecânica ainda vale como diagnóstico registrado — a
// sugestão sem patch diz isso explicitamente em vez de fingir uma correção.
func TestSugestaoSemPatchEhDiagnostico(t *testing.T) {
	root := t.TempDir()
	s := nova("s2")
	s.Patch = ""
	_, p, _ := Open(root, s)
	b, _ := os.ReadFile(p)
	if strings.Contains(string(b), "```diff") {
		t.Error("sem patch não deve haver bloco de diff")
	}
	if !strings.Contains(string(b), "precisa de decisão humana") {
		t.Errorf("deve dizer que a correção não é mecânica: %s", b)
	}
}

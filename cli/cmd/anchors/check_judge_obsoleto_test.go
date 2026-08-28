package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/queue"
)

// A fila é persistente e o conjunto de alvos aplicáveis não é. Um gate que passa a
// declarar `requires`, ou um alvo que perde a marcação que o tornava aplicável, deixa
// a task velha órfã — e a fila passa a mentir sobre o tamanho do trabalho.
func TestDescartaJulgamentosObsoletos(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "no-test-prova-real", Measures: config.MeasuresJudgment},
	}}
	// os alvos têm de existir: `queue.List` já descarta task de alvo apagado, e o que
	// se testa aqui é o outro caso — alvo que existe e gate que não se aplica mais.
	for _, f := range []string{"a.spec.md", "b.spec.md", "c.spec.md"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("# spec\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	vivo := queue.Task{
		ID: "judge-no-test-prova-real-a", Changed: "a.spec.md",
		Kind: "judgment", Origin: "check",
	}
	obsoleto := queue.Task{
		ID: "judge-no-test-prova-real-b", Changed: "b.spec.md",
		Kind: "judgment", Origin: "check",
	}
	// trabalho de outra origem NÃO pode ser descartado por este caminho
	alheio := queue.Task{
		ID: "judge-no-test-prova-real-c", Changed: "c.spec.md",
		Kind: "judgment", Origin: "humano",
	}
	for _, tk := range []queue.Task{vivo, obsoleto, alheio} {
		if _, err := queue.Enqueue(root, tk); err != nil {
			t.Fatal(err)
		}
	}

	// o check desta rodada só enfileirou o `a`
	p := gate.Profile{Judged: []gate.Result{
		{Gate: "no-test-prova-real", Target: "a.spec.md"},
	}}
	gatesDeJulgamentoConhecidos = []string{"no-test-prova-real"}
	descartaJulgamentosObsoletos(root, cfg, p)

	restou := map[string]bool{}
	tasks, err := queue.List(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range tasks {
		restou[tk.ID] = true
	}
	if !restou[vivo.ID] {
		t.Error("a task do alvo ainda aplicável NÃO devia ser descartada")
	}
	if restou[obsoleto.ID] {
		t.Error("a task do alvo que o gate já não enfileira devia ter saído da fila")
	}
	if !restou[alheio.ID] {
		t.Error("task de outra origem não é deste caminho — devia ficar")
	}
}

// O ID é `judge-<gate>-<slug>`, e o nome do gate contém `-`: a leitura é por prefixo
// conhecido, não por partir no separador.
func TestGateDaTaskJudge(t *testing.T) {
	gatesDeJulgamentoConhecidos = []string{"no-test-prova-real", "atomic-design"}
	casos := map[string]string{
		"judge-no-test-prova-real-apps-x-y.spec": "no-test-prova-real",
		"judge-atomic-design-apps-x.tsx":         "atomic-design",
		"judge-gate-que-nao-existe-x":            "",
		"outra-coisa":                            "",
	}
	for id, quer := range casos {
		if got := gateDaTaskJudge(id); got != quer {
			t.Errorf("%s → %q, queria %q", id, got, quer)
		}
	}
}

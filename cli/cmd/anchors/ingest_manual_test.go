package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func projetoCom(t *testing.T, yaml string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

const comSuite = `derived:
  anchor: spec
tests:
  - layer: unit
    run: pnpm test
    junit: .test-results/junit.xml
`

// O `anchors test` roda a suíte E ingere numa operação só — é isso que garante que o
// sinal no mapa corresponde à execução que acabou de acontecer.
//
// Quem chama `ingest` direto quebra a garantia sem perceber: dá para rodar a suíte,
// editar o código, ingerir o relatório velho, e o mapa afirmar uma cobertura que já não
// vale. Medido — o `line-coverage` aprovava 1 nó quando havia 2 cobertos, e o gate estava
// certo: ele reporta o que o mapa sabe.
func TestIngestManualAvisaMasNaoBarra(t *testing.T) {
	viaAnchorsTest = false
	if err := avisaSeIngestManual(projetoCom(t, comSuite)); err != nil {
		t.Errorf("o padrão AVISA e deixa seguir: há usos legítimos (um CI que rodou a "+
			"suíte noutro job), e barrá-los tiraria a saída de quem tem razão; veio: %v", err)
	}
}

// Quando o projeto declara, a ingestão manual é RECUSADA — e a mensagem tem de dizer o
// que usar no lugar, senão quem foi barrado não sabe como seguir.
func TestIngestManualBarraQuandoOProjetoPede(t *testing.T) {
	viaAnchorsTest = false
	yaml := strings.Replace(comSuite, "derived:", "workflow:\n  manual_ingest_blocks: true\nderived:", 1)
	err := avisaSeIngestManual(projetoCom(t, yaml))
	if err == nil {
		t.Fatal("com `manual_ingest_blocks: true` a ingestão manual deve ser recusada")
	}
	if !strings.Contains(err.Error(), "anchors test") {
		t.Errorf("o erro deve dizer O QUE USAR no lugar; veio: %v", err)
	}
}

// Vindo do `anchors test`, nunca reclama — nem quando o projeto declara o bloqueio. Se
// reclamasse, o comando correto seria barrado pela regra que existe para promovê-lo.
func TestViaAnchorsTestNuncaReclama(t *testing.T) {
	viaAnchorsTest = true
	defer func() { viaAnchorsTest = false }()
	yaml := strings.Replace(comSuite, "derived:", "workflow:\n  manual_ingest_blocks: true\nderived:", 1)
	if err := avisaSeIngestManual(projetoCom(t, yaml)); err != nil {
		t.Errorf("o `anchors test` é o caminho CERTO — não pode ser barrado: %v", err)
	}
}

// SEM `tests:` declarado o `anchors test` não roda, e exigir o que não existe mandaria o
// projeto para um beco sem saída: o `ingest` é a única forma de o sinal chegar ao mapa.
func TestSemSuiteDeclaradaNaoExige(t *testing.T) {
	viaAnchorsTest = false
	semSuite := "derived:\n  anchor: spec\n"
	if err := avisaSeIngestManual(projetoCom(t, semSuite)); err != nil {
		t.Errorf("sem suíte declarada não há alternativa a exigir; veio: %v", err)
	}
}

package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// O AVISO É SOBRE CONTEÚDO, NÃO SOBRE RELÓGIO.
//
// A primeira versão comparava o mtime do mapa com o dos arquivos-fonte. Barato e errado
// em CI: o `checkout` escreve TODO o repositório no instante do clone, e cada arquivo
// fica microssegundos mais novo que o mapa. Medido no primeiro run do pipeline de gates —
// 26 arquivos "mudaram" num repositório onde nada mudara, e o log do mesmo job dizia,
// duas linhas acima, que o mapa correspondia ao repositório.
//
// Ruído assim custa mais do que parece: um aviso que aparece em todo PR verde ensina a
// ignorá-lo, e aí ele não serve quando for verdadeiro.
func TestMapaStaleNaoDependeDeMtime(t *testing.T) {
	dir := t.TempDir()
	spec := filepath.Join(dir, "a.spec.md")
	if err := os.WriteFile(spec, []byte("# A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mapPath := filepath.Join(dir, "anchors.graph.yaml")
	if err := os.WriteFile(mapPath, []byte("version: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// O MAPA mais VELHO que o arquivo — a situação que a heurística antiga acusava.
	antigo := mustTime(t, "2020-01-01T00:00:00Z")
	if err := os.Chtimes(mapPath, antigo, antigo); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Layers: map[string]config.Layer{
		"spec": {Kind: "spec", Pattern: "*.spec.md"},
	}}
	// O arquivo ESTÁ no mapa: mesmo com o mapa "mais velho", não há o que avisar.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md"}}}

	saida := capturaSaida(t, func() { warnIfMapStale(dir, mapPath, cfg, g) })
	if strings.Contains(saida, "DESATUALIZADO") {
		t.Errorf("o mapa conhece o arquivo — o mtime não pode gerar aviso.\nveio: %s", saida)
	}

	// Agora o caso REAL: um arquivo governado que o mapa não conhece.
	if err := os.WriteFile(filepath.Join(dir, "b.spec.md"), []byte("# B\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	saida = capturaSaida(t, func() { warnIfMapStale(dir, mapPath, cfg, g) })
	if !strings.Contains(saida, "DESATUALIZADO") {
		t.Error("arquivo governado fora do mapa é INVISÍVEL ao check — tem de avisar")
	}
	// A mensagem diz o que foi medido: dizer "mudou" mandaria procurar uma alteração
	// que não existe — o caso comum é arquivo NOVO, que nunca esteve no mapa.
	if strings.Contains(saida, "mudaram") {
		t.Errorf("a mensagem deve dizer 'não estão no mapa', não 'mudaram'.\nveio: %s", saida)
	}
}

// Sem grafo o aviso se cala: não há com o que comparar, e adivinhar seria pior.
func TestMapaStaleSemGrafoSeCala(t *testing.T) {
	if s := capturaSaida(t, func() {
		warnIfMapStale(t.TempDir(), "inexistente.yaml", &config.Config{}, nil)
	}); s != "" {
		t.Errorf("sem grafo não há o que conferir; veio: %s", s)
	}
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// capturaSaida coleta o que a função escreve em stdout.
func capturaSaida(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = orig
	b, _ := io.ReadAll(r)
	return string(b)
}

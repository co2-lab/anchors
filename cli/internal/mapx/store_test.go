package mapx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// O CAMPO `gerado_por` NÃO PODE CAUSAR A OSCILAÇÃO QUE EXISTE PARA DENUNCIAR.
//
// O build local grava "dev" e o CI grava a versão publicada. Se cada gravação reescrevesse
// o campo, os dois desfariam o trabalho um do outro a cada execução — exatamente o
// defeito que o campo foi criado para acusar.
//
// A primeira implementação comparava só o YAML contra o arquivo em disco, que tem o
// header: nunca eram iguais, e a guarda não guardava nada. Só um teste isolado mostrou.
func TestGeradoPorSozinhoNaoReescreveOMapa(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.graph.yaml")
	g := &Graph{Version: 1, Nodes: []Node{{ID: "a", Rev: "r1"}}}

	GeradoPor = "0.1.10"
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	antes, _ := os.ReadFile(p)

	// Outro binário, MESMO mapa: o arquivo tem de ficar intacto.
	GeradoPor = "dev"
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	depois, _ := os.ReadFile(p)
	if string(antes) != string(depois) {
		t.Error("só o `gerado_por` mudou — reescrever faz o mapa oscilar entre o build " +
			"local e o publicado, que é o defeito que o campo existe para acusar")
	}

	// Mas mudança DE VERDADE grava, inclusive o gerado_por novo.
	g.Nodes = append(g.Nodes, Node{ID: "b", Rev: "r1"})
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	final, _ := os.ReadFile(p)
	if !strings.Contains(string(final), "gerado_por: dev") {
		t.Error("quando o mapa muda de verdade, quem gravou é quem está rodando")
	}
}

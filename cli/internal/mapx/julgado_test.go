package mapx

import "testing"

// O gate de julgamento não computa nada: é o carimbo que diz se alguém já leu. E o
// carimbo leva as revs das pontas, então o veredito envelhece se o alvo mudar — nesse
// caso volta a ser pergunta, que é o anti-drift do julgamento.
func TestJulgadoPor(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Rev: "r1"},
			{ID: "a.test.ts", Rev: "r1"},
		},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}

	// antes de julgar: ninguém respondeu
	if _, ok := g.JulgadoPor("a.spec.md", "meu-gate"); ok {
		t.Error("sem carimbo, JulgadoPor deveria dizer que não há veredito")
	}

	// julgado: o veredito vale
	if n := g.StampNodeByGate("a.spec.md", "ok", "agora", "meu-gate"); n != 1 {
		t.Fatalf("carimbou %d arestas, queria 1", n)
	}
	v, ok := g.JulgadoPor("a.spec.md", "meu-gate")
	if !ok || v != "ok" {
		t.Errorf("veredito = %q (ok=%v), queria ok/true", v, ok)
	}

	// carimbo de OUTRO gate não responde por este
	if _, ok := g.JulgadoPor("a.spec.md", "outro-gate"); ok {
		t.Error("o carimbo de um gate NÃO responde por outro")
	}

	// o alvo mudou depois → o veredito envelhece e volta a ser pergunta
	g.Nodes[0].Rev = "r2"
	if _, ok := g.JulgadoPor("a.spec.md", "meu-gate"); ok {
		t.Error("alvo alterado depois do carimbo: o veredito deveria ter envelhecido")
	}
}

// O `check` reescreve o Stamp a cada rodada (resume o pior veredito de todos os
// gates). O julgamento NÃO pode morrer nisso — foi o que apagava 16 vereditos e
// mantinha o contador em 16.
func TestJulgamentoSobreviveAoCarimboDoCheck(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Rev: "r1"},
			{ID: "a.test.ts", Rev: "r1"},
		},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}
	g.StampNodeByGate("a.spec.md", "ok", "t0", "meu-gate")

	// o check passa e reescreve o Stamp das duas pontas confrontadas
	g.StampEdges([]NodeVerdict{{ID: "a.spec.md"}, {ID: "a.test.ts"}}, "t1")

	v, ok := g.JulgadoPor("a.spec.md", "meu-gate")
	if !ok || v != "ok" {
		t.Errorf("o julgamento deveria sobreviver ao carimbo do check: veredito=%q ok=%v", v, ok)
	}
}

// `map build` reconstrói o grafo, e o `check --all` roda build ANTES de confrontar:
// se o rebuild apaga o julgamento, o veredito morre no caminho entre um comando e o
// seguinte.
func TestJulgamentoSobreviveAoRebuild(t *testing.T) {
	antigo := &Graph{
		Nodes: []Node{{ID: "a.spec.md", Rev: "r1"}, {ID: "a.test.ts", Rev: "r1"}},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}
	antigo.StampNodeByGate("a.spec.md", "ok", "t0", "meu-gate")

	// o build monta um grafo novo, com as mesmas arestas e sem carimbo
	novo := &Graph{
		Nodes: []Node{{ID: "a.spec.md", Rev: "r1"}, {ID: "a.test.ts", Rev: "r1"}},
		Edges: []Edge{{From: "a.spec.md", To: "a.test.ts"}},
	}
	PreservarCarimbos(novo, antigo)

	if v, ok := novo.JulgadoPor("a.spec.md", "meu-gate"); !ok || v != "ok" {
		t.Errorf("o julgamento deveria sobreviver ao rebuild: veredito=%q ok=%v", v, ok)
	}
}

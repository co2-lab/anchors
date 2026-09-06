package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O git trata o `anchors.graph.yaml` como TEXTO, e ele é derivado.
//
// Medido no blue-eyes (co2-lab/anchors#12): um `git merge origin/develop` mesclou o mapa
// SEM CONFLITO e apagou 62 carimbos de julgamento. O aviso do `map build` (v0.1.43) não
// pega esse caso — ele compara com o arquivo anterior, e a perda acontece no `git merge`,
// antes de o Anchors ser chamado.
//
// Um carimbo é estado do trabalho, e o LAUDO vive no `--reason` do `anchors judge`, não no
// arquivo: perder o carimbo manda o julgamento de volta para a fila e a evidência tem de
// ser refeita.
func TestMapMerge_uneOsCarimbosDosDoisLados(t *testing.T) {
	dir := t.TempDir()

	aresta := func(from, to string, gates ...string) mapx.Edge {
		e := mapx.Edge{Type: "specifies", From: from, To: to}
		for _, g := range gates {
			e.Julgamentos = append(e.Julgamentos, mapx.Judgment{Gate: g, Verdict: "ok"})
		}
		return e
	}
	grava := func(nome string, g *mapx.Graph) string {
		p := filepath.Join(dir, nome)
		if err := mapx.Save(g, p); err != nil {
			t.Fatal(err)
		}
		return p
	}

	// os dois lados compartilham a aresta A (julgada nos dois) e têm uma exclusiva cada
	base := grava("base.yaml", &mapx.Graph{Edges: []mapx.Edge{aresta("a.spec.md", "a.ts")}})
	nosso := grava("nosso.yaml", &mapx.Graph{Edges: []mapx.Edge{
		aresta("a.spec.md", "a.ts", "review"),
		aresta("b.spec.md", "b.ts", "rule-fulfilled"), // só nós
	}})
	deles := grava("deles.yaml", &mapx.Graph{Edges: []mapx.Edge{
		aresta("a.spec.md", "a.ts", "review"),
		aresta("c.spec.md", "c.ts", "review"), // só eles
	}})

	cmd := newMapMergeCmd()
	cmd.SetArgs([]string{base, nosso, deles})
	cmd.SetOut(os.Stderr)
	cmd.SetErr(os.Stderr)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("map merge: %v", err)
	}

	// o resultado vai em `nosso`, que é o que o git espera
	g, err := mapx.Load(nosso)
	if err != nil {
		t.Fatal(err)
	}
	por := map[string][]string{}
	for _, e := range g.Edges {
		for _, j := range e.Julgamentos {
			por[e.From] = append(por[e.From], j.Gate)
		}
	}
	for _, caso := range []struct{ from, gate string }{
		{"a.spec.md", "review"},         // dos dois
		{"b.spec.md", "rule-fulfilled"}, // só nosso — não pode sumir
		{"c.spec.md", "review"},         // só deles — tem de vir
	} {
		achou := false
		for _, g := range por[caso.from] {
			if g == caso.gate {
				achou = true
			}
		}
		if !achou {
			t.Errorf("o julgamento %q de %s sumiu na união: %v", caso.gate, caso.from, por)
		}
	}
}

// União e não substituição: o lado que só tem carimbo NÃO pode perdê-lo para um lado que
// não tem nenhum. É o caso do merge que apagou os 62 — um lado vazio venceria.
func TestMapMerge_ladoVazioNaoApagaOOutro(t *testing.T) {
	dir := t.TempDir()
	comCarimbo := &mapx.Graph{Edges: []mapx.Edge{{
		Type: "specifies", From: "a.spec.md", To: "a.ts",
		Julgamentos: []mapx.Judgment{{Gate: "review", Verdict: "ok"}},
	}}}
	semCarimbo := &mapx.Graph{Edges: []mapx.Edge{{
		Type: "specifies", From: "a.spec.md", To: "a.ts",
	}}}

	p := func(nome string, g *mapx.Graph) string {
		caminho := filepath.Join(dir, nome)
		if err := mapx.Save(g, caminho); err != nil {
			t.Fatal(err)
		}
		return caminho
	}

	for _, caso := range []struct {
		nome         string
		nosso, deles *mapx.Graph
	}{
		{"o carimbo está do nosso lado", comCarimbo, semCarimbo},
		{"o carimbo está do outro lado", semCarimbo, comCarimbo},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			base := p("base-"+caso.nome+".yaml", semCarimbo)
			nosso := p("nosso-"+caso.nome+".yaml", caso.nosso)
			deles := p("deles-"+caso.nome+".yaml", caso.deles)

			cmd := newMapMergeCmd()
			cmd.SetArgs([]string{base, nosso, deles})
			cmd.SetOut(os.Stderr)
			cmd.SetErr(os.Stderr)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			g, err := mapx.Load(nosso)
			if err != nil {
				t.Fatal(err)
			}
			if n := countJudgments(g); n != 1 {
				t.Errorf("resultado com %d julgamento(s), queria 1 — o lado vazio apagou o outro", n)
			}
		})
	}
}

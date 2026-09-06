package main

import (
	"strings"
	"testing"
)

// `--decision` e `--uncovered` recebem PROSA, e prosa tem vírgula.
//
// Com `StringSliceVar` o pflag divide o valor na vírgula. Medido no blue-eyes, entregando
// a spec do DataStore:
//
//	--decision "oito regras e dois invariantes, na letra B/I que as vizinhas usam"
//
// virou DUAS decisões no registro — a segunda começando com espaço e sem sujeito. De
// cinco decisões declaradas saíram nove itens, quatro deles fragmentos.
//
// O dano é sobre o que o registro existe para fazer: o revisor confronta cada decisão
// contra o disco, e meia frase não é confrontável. E a contagem infla — "nove decisões"
// descreve um trabalho que tomou cinco.
func TestDeliver_decisaoComVirgulaNaoEhDividida(t *testing.T) {
	cmd := newDeliverCmd()
	const prosa = "oito regras e dois invariantes, na letra B/I que as vizinhas do F03 usam"

	if err := cmd.Flags().Parse([]string{
		"--decision", prosa,
		"--uncovered", "o I01 se prova lendo o código de escrita, que ainda não existe",
	}); err != nil {
		t.Fatal(err)
	}

	for _, nome := range []string{"decision", "uncovered"} {
		v, err := cmd.Flags().GetStringArray(nome)
		if err != nil {
			t.Fatalf("--%s não é StringArray: %v — prosa com vírgula será dividida", nome, err)
		}
		if len(v) != 1 {
			t.Errorf("--%s virou %d itens: %q", nome, len(v), v)
		}
	}

	d, _ := cmd.Flags().GetStringArray("decision")
	if d[0] != prosa {
		t.Errorf("a decisão chegou alterada:\n  quer: %q\n  veio: %q", prosa, d[0])
	}
	// o sintoma que apareceu no registro: um item começando com espaço
	for _, item := range d {
		if strings.HasPrefix(item, " ") {
			t.Errorf("item com espaço à esquerda (%q) — sinal de divisão na vírgula", item)
		}
	}
}

// `--file` CONTINUA dividindo: caminho de arquivo não tem vírgula, e ali a divisão é
// conveniência real (`--file a.ts,b.ts`). A correção não é trocar tudo por StringArray.
func TestDeliver_fileAindaDivideNaVirgula(t *testing.T) {
	cmd := newDeliverCmd()
	if err := cmd.Flags().Parse([]string{"--file", "a.ts,b.ts"}); err != nil {
		t.Fatal(err)
	}
	f, err := cmd.Flags().GetStringSlice("file")
	if err != nil {
		t.Fatalf("--file deixou de ser StringSlice: %v", err)
	}
	if len(f) != 2 {
		t.Errorf("--file a.ts,b.ts virou %d item(ns): %q", len(f), f)
	}
}

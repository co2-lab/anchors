package main

import (
	"strings"
	"testing"
)

// O BOARD PUBLICADO E UMA FOTO, e a foto envelhece.
//
// O pipeline roda, gera o `board.json`, publica. Entre uma execucao e outra o que se ve e'
// passado -- e num projeto com nove agentes entregando, minutos bastam para a foto mentir.
//
// O `board serve` sobe o MESMO html apontando para dado vivo. O "mesmo" nao e' detalhe:
// duas renderizacoes do board divergiriam, e a local passaria a mentir de outro jeito.
func TestBoardServeExiste(t *testing.T) {
	root := newRootCmd()
	var achou bool
	for _, c := range root.Commands() {
		if c.Name() == "board" {
			achou = true
			var temServe bool
			for _, s := range c.Commands() {
				if s.Name() == "serve" {
					temServe = true
				}
			}
			if !temServe {
				t.Error("`anchors board` existe mas nao tem `serve` — sem ele o board " +
					"local nao sobe")
			}
		}
	}
	if !achou {
		t.Fatal("`anchors board` nao existe")
	}
}

// A ESCUTA e' o desenho B: o custo acompanha a MUDANCA, nao o tempo.
//
// Um poll de 30s custa 1800 chamadas por hora, mudando algo ou nao -- e o rate limit do
// GitHub ja estourou neste projeto com seis agentes. A escuta incremental pergunta "o que
// mudou desde X?", e o custo cai para o que de fato mudou.
//
// O FALLBACK e' o desenho A, e existe porque a escuta pode nao estar disponivel (sem
// webhook, sem permissao, API recusando). Sem ele o comando falharia onde o poll
// funcionaria -- e um board que nao sobe e' pior que um board com alguns segundos de
// atraso.
func TestBoardServeEscutaComFallbackParaPoll(t *testing.T) {
	root := newRootCmd()
	var longo string
	var temIntervalo bool
	for _, c := range root.Commands() {
		if c.Name() != "board" {
			continue
		}
		for _, sub := range c.Commands() {
			if sub.Name() != "serve" {
				continue
			}
			longo = sub.Long
			if sub.Flags().Lookup("interval") != nil {
				temIntervalo = true
			}
		}
	}
	if longo == "" {
		t.Fatal("`board serve` nao existe ou nao tem doc")
	}

	// O intervalo do fallback e' configuravel: projeto grande e projeto pequeno nao
	// toleram a mesma frequencia.
	if !temIntervalo {
		t.Error("falta a flag `--interval` — o fallback precisa ser calibravel pelo projeto")
	}

	// A doc precisa NOMEAR o fallback. Um comando que silenciosamente troca de estrategia
	// e' um comando cujo comportamento ninguem consegue prever.
	if !strings.Contains(strings.ToLower(longo), "fallback") &&
		!strings.Contains(strings.ToLower(longo), "recua") {
		t.Error("a doc nao diz que ha fallback — trocar de estrategia em silencio faz o " +
			"comportamento ficar imprevisivel para quem depende dele")
	}
}

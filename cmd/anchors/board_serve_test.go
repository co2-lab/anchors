package main

import (
	"os"
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

// O INCREMENTAL PRECISA PERGUNTAR PELA MAIS RECENTE, e nao pela de maior numero.
//
// MEDIDO: o `gh issue list --limit 1` devolve a issue de maior NUMERO, nao a tocada por
// ultimo. Onze cards foram fechados e o board NAO releu, porque o topo da lista (`#683`,
// 19:22) era mais velho que o card de fato tocado (`#630`, 19:29).
//
// O defeito e' silencioso do pior jeito: o board segue servindo, com cara de fresco, e
// mostra passado. E' a foto de novo -- agora local, e sem ninguem saber.
//
// Duas exigencias, e as duas vem de medicao:
//
//	sort=updated   sem isso a lista vem por numero, e a pergunta responde errado
//	gh api (REST)  o `gh issue list` usa GraphQL, que tem limite SECUNDARIO invisivel
//	               nos contadores -- derruba a consulta enquanto `rate_limit` diz cheio
func TestIncrementalPerguntaPelaMaisRecente(t *testing.T) {
	b, err := os.ReadFile("board_serve.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)

	i := strings.Index(src, "func (f *boardSource) algoMudou()")
	if i < 0 {
		t.Fatal("nao achei o `algoMudou` — o incremental sumiu")
	}
	corpo := src[i:]
	if j := strings.Index(corpo, "\n}"); j > 0 {
		corpo = corpo[:j]
	}

	if !strings.Contains(corpo, "sort=updated") {
		t.Error("o incremental nao ordena por `updated` — a consulta devolve a issue de " +
			"maior NUMERO, e uma issue antiga tocada agora fica invisivel. O board segue " +
			"servindo com cara de fresco, mostrando passado")
	}
	if strings.Contains(corpo, `"issue", "list"`) {
		t.Error("o incremental usa `gh issue list`, que vai por GraphQL — e o limite " +
			"SECUNDARIO dele derruba a consulta enquanto `gh api rate_limit` ainda " +
			"reporta cota cheia. O REST nao tem esse problema")
	}
}

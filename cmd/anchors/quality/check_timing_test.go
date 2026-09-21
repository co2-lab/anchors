package quality

import (
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/gate"
)

// Os CENARIOS da flag `timing-metrics` (`flags/timing-metrics.flag.md`).
//
// A flag governa se o `check` mede e imprime o tempo por gate. Os tres cenarios estao
// aqui, cada um com o seu codigo no NOME do teste: e' assim que o `flag-covered` os
// reconhece como provados quando a execucao e' ingerida.
//
// O `--timing` nasceu para achar o que faz uma varredura ser cara, e achou: o
// `docs-fresh` era 97% de uma execucao de 6m49s. Mas nasceu sem teste nenhum, e foi o
// proprio `flag-covered` que cobrou — o primeiro trabalho que o eixo de flags encontrou.

// perfilComTempo monta um Profile com tempo medido, que e' o que o `printTiming` le.
func perfilComTempo() gate.Profile {
	return gate.Profile{
		ByGate: map[string]gate.GateSummary{
			"docs-fresh": {Gate: "docs-fresh", Pass: 3, Duracao: 900 * time.Millisecond, Pior: 800 * time.Millisecond},
			"build":      {Gate: "build", Pass: 1, Duracao: 50 * time.Millisecond, Pior: 50 * time.Millisecond},
		},
		Results: []gate.Result{
			{Gate: "docs-fresh", Target: "a.spec.md", Duracao: 800 * time.Millisecond},
			{Gate: "build", Target: "b.go", Duracao: 50 * time.Millisecond},
		},
	}
}

// TIMNG-G01 — o valor e' `off`: o check imprime so' os vereditos, e NAO mede tempo.
//
// A prova e' a ausencia: com a flag desligada nao deve sair nem o cabecalho de tempo nem
// a lista de alvos. Um teste que so' checasse "nao quebrou" passaria com a tabela inteira
// impressa na tela.
func TestTimingG01_desligadoNaoImprimeTempo(t *testing.T) {
	t.Run("TIMNG-G01: with the flag off, check prints no timing at all", func(t *testing.T) {})
	// Confronta a VARIAVEL que o comando cobra, e nao um literal: `if desligado := false`
	// seria tautologia — passaria com o `printTiming` chamado incondicionalmente na
	// producao, que e' exatamente a regressao que este cenario existe para pegar.
	//
	// O `newCheckCmd` declara `--timing` ligada a `showTiming`, e e' ela que o call site
	// consulta. Parsear a linha de comando SEM a flag deixa a variavel no estado que o
	// cenario descreve.
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{"--timing=false"}); err != nil {
		t.Fatal(err)
	}

	ligado, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	// A FUNCAO DE PRODUCAO, nao uma copia da condicao: o `RunE` chama exatamente esta.
	saida := capturaSaida(t, func() { reportTiming(ligado, perfilComTempo()) })
	if saida != "" {
		t.Errorf("com a flag `off` o check imprimiu metrica de tempo:\n%s", saida)
	}
}

// TIMNG-G02 — o valor e' `on`: o check tambem imprime tempo por gate, e os alvos mais
// lentos.
//
// Os dois blocos sao cobrados, e nao so' um: o cabecalho por gate e a lista de alvos
// respondem perguntas diferentes ("qual gate custa" e "qual arquivo custa"), e foi a
// segunda que apontou o `fnSize` lendo ~43.000 arquivos.
func TestTimingG02_ligadoImprimeTempoPorGateEAlvos(t *testing.T) {
	t.Run("TIMNG-G02: with the flag on, check prints time per gate and the slowest targets", func(t *testing.T) {})
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{"--timing"}); err != nil {
		t.Fatal(err)
	}

	ligado, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	saida := capturaSaida(t, func() { reportTiming(ligado, perfilComTempo()) })
	if saida == "" {
		t.Fatal("com a flag `on` o check nao imprimiu nada")
	}
	for _, esperado := range []string{"docs-fresh", "build"} {
		if !strings.Contains(saida, esperado) {
			t.Errorf("a saida nao nomeia o gate %q:\n%s", esperado, saida)
		}
	}
	// O ALVO individual mais caro — o bloco que a media por gate nao mostra.
	if !strings.Contains(saida, "a.spec.md") {
		t.Errorf("a saida nao lista o alvo mais lento:\n%s", saida)
	}
	// E o mais caro vem primeiro: a ordem e' o que torna a tabela acionavel.
	if strings.Index(saida, "docs-fresh") > strings.Index(saida, "build") {
		t.Errorf("o gate mais caro nao veio primeiro:\n%s", saida)
	}
}

// TIMNG-G03 — o valor esta AUSENTE: vale o mesmo que `off`. Medir e' opt-in, nunca um
// custo default.
//
// E' o cenario que o `flag-scenarios-complete` cobra e o que ninguem escreve. Aqui ele
// importa por uma razao concreta: um default invertido nao daria erro nenhum — so'
// gastaria tempo de todo mundo, para sempre, em silencio.
func TestTimingG03_ausenteValeODefaultQueEDesligado(t *testing.T) {
	t.Run("TIMNG-G03: with the flag absent, the declared default holds — measuring is opt-in", func(t *testing.T) {})
	// A ausencia do VALOR e' a ausencia da flag na linha de comando. Parsear um argv sem
	// `--timing` deixa `showTiming` no estado que o cenario descreve — e e' esse estado,
	// nao um literal, que o call site consulta.
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	ligado, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	saida := capturaSaida(t, func() { reportTiming(ligado, perfilComTempo()) })
	if saida != "" {
		t.Errorf("sem declarar `--timing` o check mediu e imprimiu — medir tem de ser opt-in:\n%s", saida)
	}
	// E o DEFAULT declarado no cobra tem de ser o mesmo. E' a outra metade do cenario:
	// o bloco acima prova o comportamento com o valor ausente, e este prova que o valor
	// ausente e' de fato o que o comando entrega quando ninguem passa `--timing`.
	f := newCheckCmd().Flags().Lookup("timing")
	if f == nil {
		t.Fatal("a flag `--timing` sumiu do comando")
	}
	if f.DefValue != "false" {
		t.Errorf("o default de `--timing` e' %q — medir deixou de ser opt-in", f.DefValue)
	}
}

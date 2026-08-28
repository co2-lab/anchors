package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
)

// perfil monta um Profile com os contadores que interessam ao formato.
func perfil(gates ...gate.GateSummary) gate.Profile {
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}, Passed: true}
	for _, g := range gates {
		p.ByGate[g.Gate] = g
	}
	return p
}

// A tabela alinha para ser COMPARADA a olho: numa lista de 49 gates, a coluna
// torta é o que faz o leitor perder o número que importa.
func TestColunasAlinhamComNomeLongo(t *testing.T) {
	nomes := []string{"eslint", "handler-ddb-inline-passivo", "circular"}
	w := larguraDoNome(nomes)
	if w < len("handler-ddb-inline-passivo") {
		t.Fatalf("largura %d corta o maior nome (%d)", w, len("handler-ddb-inline-passivo"))
	}
	// Piso: com nomes curtos a coluna não encolhe a ponto de colar no veredito.
	if got := larguraDoNome([]string{"a", "bb"}); got < 18 {
		t.Errorf("larguraDoNome sem piso: %d", got)
	}
}

// Uma largura ÚNICA para todas as colunas desperdiça espaço: se `~` chega a 582
// e `✗` nunca passa de 1, a coluna dos fails reservaria três casas para nada.
func TestLarguraEPorColuna(t *testing.T) {
	p := perfil(
		gate.GateSummary{Gate: "a", Pass: 1116, Fail: 0, Skip: 0},
		gate.GateSummary{Gate: "b", Pass: 1, Fail: 1, Skip: 582},
	)
	w := calcularLarguras(p)

	if w.pass != 4 {
		t.Errorf("pass: %d, queria 4 (por causa de 1116)", w.pass)
	}
	if w.fail != 1 {
		t.Errorf("fail: %d, queria 1 — a coluna não deve herdar a largura de outra", w.fail)
	}
	if w.skip != 3 {
		t.Errorf("skip: %d, queria 3 (por causa de 582)", w.skip)
	}
}

// As colunas SEMPRE presentes têm piso 1: `%*d` com largura 0 imprimiria colado
// no símbolo. A do drift é a exceção — ela nasce 0 e só abre com drift real.
func TestLarguraMinimaEUm(t *testing.T) {
	w := calcularLarguras(perfil(gate.GateSummary{Gate: "a"}))
	for nome, got := range map[string]int{
		"pass": w.pass, "fail": w.fail, "skip": w.skip, "judge": w.judge,
	} {
		if got != 1 {
			t.Errorf("%s: %d, queria 1", nome, got)
		}
	}
	// Sem drift a coluna não existe: reservá-la deixaria um buraco no meio de
	// todas as linhas sem nada que o justificasse — e a tabela sem drift é o
	// caso comum, não a exceção.
	if w.drift != 0 {
		t.Errorf("drift: %d, queria 0 — a coluna não deve existir sem drift", w.drift)
	}
}

// A coluna do ⚠ some da tabela inteira quando nenhum gate tem drift, e o
// separador some com ela: senão sobrariam dois espaços entre `✗` e `~`.
func TestSemDriftNenhumAColunaNaoExiste(t *testing.T) {
	p := perfil(
		gate.GateSummary{Gate: "um", Pass: 1},
		gate.GateSummary{Gate: "dois", Pass: 583, Skip: 12},
	)
	out := capturarStdout(t, func() { printProfile(p, false, false) })

	for _, linha := range linhasDeGate(out) {
		// Entre o contador de `✗` e o `~` só pode haver o separador de 2 espaços.
		i := indiceDeRuna([]rune(linha), '✗')
		j := indiceDeRuna([]rune(linha), '~')
		if i < 0 || j < 0 {
			continue
		}
		meio := string([]rune(linha)[i:j])
		if strings.Count(meio, " ") > 2+len("0") {
			t.Errorf("buraco entre ✗ e ~ sem drift na tabela: %q", linha)
		}
	}
}

// `--only-issues` omite quem passou em tudo E não deixou pendência. Um gate com
// `~` não é limpo: ele não confrontou nada, e isso é informação.
func TestGateLimpoExigeAusenciaDePendencia(t *testing.T) {
	casos := []struct {
		nome  string
		s     gate.GateSummary
		drift int
		limpo bool
	}{
		{"tudo zerado com passes", gate.GateSummary{Pass: 10}, 0, true},
		{"com fail", gate.GateSummary{Pass: 10, Fail: 1}, 0, false},
		{"com drift", gate.GateSummary{Pass: 10}, 3, false},
		{"com skip", gate.GateSummary{Pass: 10, Skip: 2}, 0, false},
		{"com pending", gate.GateSummary{Pass: 10, Pending: 2}, 0, false},
		{"aguardando IA", gate.GateSummary{Pass: 10, Judge: 1}, 0, false},
	}
	for _, c := range casos {
		if got := gateLimpo(c.s, c.drift); got != c.limpo {
			t.Errorf("%s: gateLimpo = %v, queria %v", c.nome, got, c.limpo)
		}
	}
}

// O default NÃO esconde nada: a tabela cheia é a prova de que os 49 gates
// rodaram. Esconder por padrão trocaria essa prova por brevidade.
func TestDefaultMostraGateLimpo(t *testing.T) {
	p := perfil(
		gate.GateSummary{Gate: "limpo", Pass: 5},
		gate.GateSummary{Gate: "sujo", Pass: 1, Fail: 1},
	)
	out := capturarStdout(t, func() { printProfile(p, false, false) })

	if !strings.Contains(out, "limpo") {
		t.Errorf("gate limpo sumiu do default:\n%s", out)
	}
	if strings.Contains(out, "omitidos") {
		t.Errorf("rodapé de omissão apareceu sem a flag:\n%s", out)
	}
}

func TestOnlyIssuesOmiteLimpoMasContaNoRodape(t *testing.T) {
	p := perfil(
		gate.GateSummary{Gate: "limpo-um", Pass: 5},
		gate.GateSummary{Gate: "limpo-dois", Pass: 7},
		gate.GateSummary{Gate: "sujo", Pass: 1, Fail: 1},
	)
	out := capturarStdout(t, func() { printProfile(p, true, false) })

	if strings.Contains(out, "limpo-um") || strings.Contains(out, "limpo-dois") {
		t.Errorf("--only-issues não omitiu o gate limpo:\n%s", out)
	}
	if !strings.Contains(out, "sujo") {
		t.Errorf("--only-issues omitiu quem tem achado:\n%s", out)
	}
	// O gate omitido RODOU. Sem o número, a saída pareceria uma varredura menor.
	if !strings.Contains(out, "2 gate(s) sem nada a reportar") {
		t.Errorf("rodapé não contou os omitidos:\n%s", out)
	}
}

// capturarStdout coleta o que `fn` imprime. O `printProfile` escreve direto no
// stdout (é um comando de terminal, não uma biblioteca), então testar o FORMATO
// exige interceptá-lo.
func capturarStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = original }()

	pronto := make(chan string, 1)
	go func() {
		var b strings.Builder
		io.Copy(&b, r)
		pronto <- b.String()
	}()

	fn()
	w.Close()
	return <-pronto
}

// A coluna do `~` tem de cair na MESMA posição com e sem `⚠`. Antes, a linha sem
// drift omitia a coluna inteira e o `~` migrava para a esquerda só nela — a
// mesma coluna existindo em dois lugares na mesma tabela.
func TestColunaDoSkipNaoMigraComOuSemDrift(t *testing.T) {
	p := perfil(
		gate.GateSummary{Gate: "com-drift", Pass: 110, Pending: 423},
		gate.GateSummary{Gate: "sem-drift", Pass: 2537, Skip: 103},
	)
	// `driftCount` lê de Results: sem eles a coluna do ⚠ nem existe, e o teste
	// mediria a tabela errada.
	p.Results = []gate.Result{
		{Gate: "com-drift", Verdict: gate.Pending, Detail: "divergiu", Target: "x"},
	}
	// `driftCount` lê de p.Results; sem resultados o drift é 0 nas duas. Para
	// exercitar o formato, medimos a função da célula diretamente.
	if semDrift, comDrift := colunaDrift(0, 3), colunaDrift(407, 3); len([]rune(semDrift)) != len([]rune(comDrift)) {
		t.Errorf("célula do ⚠ com larguras diferentes: sem=%d runas, com=%d runas (%q vs %q)",
			len([]rune(semDrift)), len([]rune(comDrift)), semDrift, comDrift)
	}

	out := capturarStdout(t, func() { printProfile(p, false, false) })
	var posicoes []int
	// `TrimSpace` no bloco inteiro comeria a indentação da PRIMEIRA linha e só
	// dela — que é exatamente a coluna sendo medida. O corte é por linha.
	for _, linha := range linhasDeGate(out) {
		// A posição é a COLUNA do terminal, contada em runas. `strings.Index`
		// devolve deslocamento em BYTES, e `✓`/`✗`/`~` são multi-byte: recortar
		// por byte e depois contar runas mistura as duas unidades e acusa
		// desalinhamento onde as colunas coincidem.
		if i := indiceDeRuna([]rune(linha), '~'); i >= 0 {
			posicoes = append(posicoes, i)
		}
	}
	if len(posicoes) < 2 {
		t.Fatalf("esperava 2 linhas com ~:\n%s", out)
	}
	for _, pos := range posicoes[1:] {
		if pos != posicoes[0] {
			t.Errorf("o ~ mudou de coluna entre linhas (%v):\n%s", posicoes, out)
		}
	}
}

// `⚠` ocupa 1 coluna no terminal e 3 bytes em Go. Reservar o branco com `len()`
// desalinha justamente o que ele existe para alinhar.
func TestCelulaDoDriftMedeEmColunasNaoEmBytes(t *testing.T) {
	if got := len([]rune(colunaDrift(0, 3))); got != 4 {
		t.Errorf("célula vazia com %d colunas, queria 4 (símbolo + 3 dígitos)", got)
	}
}

// indiceDeRuna devolve a posição de `alvo` contada em runas (colunas), não em
// bytes.
func indiceDeRuna(runas []rune, alvo rune) int {
	for i, r := range runas {
		if r == alvo {
			return i
		}
	}
	return -1
}

// linhasDeGate filtra as linhas da TABELA. A legenda e os blocos de detalhe
// também contêm `✓`/`✗`/`~`, e medi-los como se fossem colunas acusaria
// desalinhamento onde não há tabela nenhuma.
func linhasDeGate(out string) []string {
	var out2 []string
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "bloqueante") || strings.Contains(l, "informativo") || strings.Contains(l, "julgamento") {
			out2 = append(out2, l)
		}
	}
	return out2
}

// Com a flag, o drift NÃO herda o corte de ruído do `~`: o `~` é suprimido em
// varredura grande porque milhares de "não se aplica" são ruído, mas quem pediu
// os endereços das pendências os quer todos, e o volume da varredura não muda
// isso.
func TestShowDriftNaoHerdaOCorteDoSkip(t *testing.T) {
	p := perfil(gate.GateSummary{Gate: "layer-boundary", Pass: 1, Pending: 1})
	// Volume acima do corte que suprime o detalhe dos `~`.
	for i := 0; i < maxResultsForSkipDetail*3; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "outro", Verdict: gate.Skip, Detail: "não se aplica", Target: "n",
		})
	}
	p.Results = append(p.Results, gate.Result{
		Gate: "layer-boundary", Verdict: gate.Pending,
		Detail: "fronteira em migração", Target: "AlertSheet.tsx",
	})

	out := capturarStdout(t, func() { printProfile(p, false, true) })

	if !strings.Contains(out, "pendência(s)") {
		t.Errorf("bloco de drift ausente numa varredura grande:\n%s", out)
	}
	if !strings.Contains(out, "AlertSheet.tsx") {
		t.Errorf("drift sem endereço — é a única parte acionável:\n%s", out)
	}
}

// O `~` continua cortado em varredura grande: 2.430 linhas de "não se aplica"
// são ruído, e o contador basta.
func TestSkipContinuaSuprimidoEmVarreduraGrande(t *testing.T) {
	p := perfil(gate.GateSummary{Gate: "g", Pass: 1, Skip: 1})
	for i := 0; i < maxResultsForSkipDetail*3; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "g", Verdict: gate.Skip, Detail: "não é uma tela", Target: "n",
		})
	}
	out := capturarStdout(t, func() { printProfile(p, false, false) })

	if strings.Contains(out, "indeterminado(s) — não é falha") {
		t.Errorf("detalhe dos ~ deveria estar suprimido no volume grande:\n%s", out)
	}
}

// Quem pede os endereços quer agir sobre eles: truncar a lista devolveria o
// problema que ela existe para resolver. Ou é completa, ou o contador da tabela
// já bastava.
func TestShowDriftListaTodasSemTeto(t *testing.T) {
	const n = 120
	p := perfil(gate.GateSummary{Gate: "g", Pass: 1, Pending: n})
	for i := 0; i < n; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "g", Verdict: gate.Pending, Detail: "divergiu",
			Target: fmt.Sprintf("alvo-%03d", i),
		})
	}
	out := capturarStdout(t, func() { printProfile(p, false, true) })

	for _, i := range []int{0, n / 2, n - 1} {
		alvo := fmt.Sprintf("alvo-%03d", i)
		if !strings.Contains(out, alvo) {
			t.Errorf("%s não foi listado — a lista está truncada", alvo)
		}
	}
	if strings.Contains(out, "e mais") {
		t.Errorf("lista truncada com --show-drift:\n%s", out)
	}
}

// Sem a flag, o contador da tabela é tudo: numa varredura grande são milhares de
// linhas, e despejá-las sem ninguém pedir enterraria as issues logo abaixo.
func TestSemShowDriftNaoListaDetalhe(t *testing.T) {
	p := perfil(gate.GateSummary{Gate: "g", Pass: 1, Pending: 1})
	p.Results = append(p.Results, gate.Result{
		Gate: "g", Verdict: gate.Pending, Detail: "divergiu", Target: "alvo-unico",
	})
	out := capturarStdout(t, func() { printProfile(p, false, false) })

	if strings.Contains(out, "alvo-unico") {
		t.Errorf("detalhe listado sem a flag:\n%s", out)
	}
	// O sinal não some: a tabela continua contando.
	if !strings.Contains(out, "⚠1") {
		t.Errorf("o contador do ⚠ sumiu da tabela:\n%s", out)
	}
}

// A legenda explica só o que a tabela usou: listar `⚠` numa varredura sem drift
// ensina a pular a legenda inteira.
func TestLegendaSoMostraOsSimbolosUsados(t *testing.T) {
	semDrift := capturarStdout(t, func() {
		printProfile(perfil(gate.GateSummary{Gate: "g", Pass: 1}), false, false)
	})
	if strings.Contains(semDrift, "⚠  divergiu") {
		t.Errorf("legenda explicou ⚠ numa tabela sem drift:\n%s", semDrift)
	}
	if !strings.Contains(semDrift, "✓  passou") || !strings.Contains(semDrift, "~  indeterminado") {
		t.Errorf("legenda incompleta:\n%s", semDrift)
	}

	comJudge := capturarStdout(t, func() {
		printProfile(perfil(gate.GateSummary{Gate: "g", Judge: 2}), false, false)
	})
	if !strings.Contains(comJudge, "⏳ aguarda julgamento") {
		t.Errorf("legenda sem ⏳ numa tabela que o usa:\n%s", comJudge)
	}
}

// 832 pendências com o motivo IDÊNTICO não são 832 problemas — são um só (a
// ingestão que ninguém rodou). Repetir o mesmo parágrafo 832 vezes esconde essa
// leitura em vez de revelá-la.
func TestDriftAgrupaMotivoRepetido(t *testing.T) {
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}}
	for i := 0; i < 50; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "mutation-score", Verdict: gate.Pending,
			Detail: "sem sinal de mutação ingerido", Target: fmt.Sprintf("alvo-%02d", i),
		})
	}
	out := capturarStdout(t, func() { printDrift(driftResults(p)) })

	if n := strings.Count(out, "sem sinal de mutação ingerido"); n != 1 {
		t.Errorf("motivo repetido %d vezes, queria 1:\n%s", n, out)
	}
	if !strings.Contains(out, "50 alvo(s):") {
		t.Errorf("cabeçalho do grupo ausente:\n%s", out)
	}
	// Um alvo por LINHA: juntá-los com vírgula produzia uma linha de 36 mil
	// caracteres no projeto real — o endereço estava lá e ninguém o lia.
	for _, l := range strings.Split(out, "\n") {
		if len(l) > 200 {
			t.Errorf("linha de %d caracteres — os alvos voltaram a ser concatenados", len(l))
			break
		}
	}
	// O agrupamento não pode custar o ENDEREÇO — ele é a parte acionável.
	for _, alvo := range []string{"alvo-00", "alvo-25", "alvo-49"} {
		if !strings.Contains(out, alvo) {
			t.Errorf("%s sumiu no agrupamento", alvo)
		}
	}
	if !strings.Contains(out, "mutation-score — 50") {
		t.Errorf("cabeçalho do gate sem o total:\n%s", out)
	}
}

// Onde cada motivo é ÚNICO (a divergência é específica da unidade), a lista
// continua alvo a alvo: compactar ali perderia justamente o que diferencia um
// achado do outro.
func TestDriftNaoAgrupaMotivoUnico(t *testing.T) {
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}}
	for i := 0; i < 3; i++ {
		p.Results = append(p.Results, gate.Result{
			Gate: "feature-test-match", Verdict: gate.Pending,
			Detail: fmt.Sprintf("cenário %d diverge", i), Target: fmt.Sprintf("f%d.feature", i),
		})
	}
	out := capturarStdout(t, func() { printDrift(driftResults(p)) })

	for i := 0; i < 3; i++ {
		if !strings.Contains(out, fmt.Sprintf("cenário %d diverge", i)) {
			t.Errorf("motivo %d sumiu:\n%s", i, out)
		}
	}
	if strings.Contains(out, "alvos:") {
		t.Errorf("compactou motivos que são diferentes entre si:\n%s", out)
	}
}

// O cabeçalho diz em quantos GATES as pendências estão: é a primeira leitura
// útil — 2.430 espalhadas por 5 gates é um retrato diferente de 2.430 num só.
func TestDriftCabecalhoContaGates(t *testing.T) {
	p := gate.Profile{ByGate: map[string]gate.GateSummary{}}
	for _, g := range []string{"a", "b", "a", "c"} {
		p.Results = append(p.Results, gate.Result{
			Gate: g, Verdict: gate.Pending, Detail: "x", Target: "t",
		})
	}
	out := capturarStdout(t, func() { printDrift(driftResults(p)) })
	if !strings.Contains(out, "4 pendência(s) em 3 gate(s)") {
		t.Errorf("cabeçalho errado:\n%s", out)
	}
}

// Dez gates montam a mensagem com `strings.Join(achados, "; ")`. Cinco violações
// no mesmo arquivo saíam numa linha de 800 caracteres: o leitor não distingue
// onde uma acaba e a outra começa, e o texto da regra — idêntico nas cinco —
// afoga o único dado que varia, que é o número da linha.
func TestOcorrenciasDoMesmoArquivoQuebramEmLinhas(t *testing.T) {
	detalhe := "linha 8: cor hex; linha 10: cor hex; linha 12: cor hex"
	got := indent(detalhe, "    ")

	linhas := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(linhas) != 3 {
		t.Fatalf("esperava 3 linhas, veio %d:\n%s", len(linhas), got)
	}
	// A indentação vale para TODAS: sem isso a segunda ocorrência em diante
	// encostaria na margem e sairia do bloco a que pertence.
	for _, l := range linhas {
		if !strings.HasPrefix(l, "    ") {
			t.Errorf("linha sem indentação: %q", l)
		}
	}
}

// Basta UM separador: duas ocorrências já se confundem numa linha só, e é o caso
// mais comum. (Com o corte em >= 2 separadores, o par não quebrava.)
func TestDuasOcorrenciasJaQuebram(t *testing.T) {
	got := quebraOcorrencias("linha 8: erro; linha 10: erro")
	if !strings.Contains(got, "\n") {
		t.Errorf("par não quebrou: %q", got)
	}
}

// Um detalhe de uma frase só continua inteiro — não há o que separar.
func TestOcorrenciaUnicaNaoQuebra(t *testing.T) {
	got := quebraOcorrencias("a spec não declara `## Decisões em aberto`")
	if strings.Contains(got, "\n") {
		t.Errorf("quebrou onde não havia separador: %q", got)
	}
}

// A quebra é no separador `"; "` (com espaço), não em qualquer `;`. As mensagens
// carregam regex e trechos de código — o `;` de um `#[0-9A-Fa-f]{3};` não separa
// ocorrência nenhuma, e quebrar ali picaria a mensagem no meio.
func TestNaoQuebraPontoEVirgulaColado(t *testing.T) {
	got := quebraOcorrencias("a camada não pode conter `a;b;c` — use token")
	if strings.Contains(got, "\n") {
		t.Errorf("quebrou num `;` colado, que não é separador de ocorrência: %q", got)
	}
}

// 17 gates montam a mensagem com `strings.Join(lista, ", ")` sem cortar. Uma
// lista grande vira uma linha que ninguém lê nem grepa.
func TestListaLongaComVirgulaQuebra(t *testing.T) {
	itens := make([]string, 30)
	for i := range itens {
		itens[i] = fmt.Sprintf("src/features/algum/caminho/Arquivo%02d.tsx", i)
	}
	got := indent("símbolos sem catálogo: "+strings.Join(itens, ", "), "    ")

	for _, l := range strings.Split(strings.TrimRight(got, "\n"), "\n") {
		if len([]rune(l)) > limiarQuebraLista+20 {
			t.Errorf("linha de %d runas — a lista não quebrou:\n%s", len([]rune(l)), l)
		}
	}
	if !strings.Contains(got, "Arquivo29.tsx") {
		t.Errorf("a quebra perdeu itens:\n%s", got)
	}
}

// Em texto normal a vírgula separa oração: picá-la destruiria a frase.
func TestFraseCurtaComVirgulaNaoQuebra(t *testing.T) {
	frase := "a spec existe, o código existe, e os dois se referenciam"
	if got := quebraOcorrencias(frase); strings.Contains(got, "\n") {
		t.Errorf("quebrou uma frase normal: %q", got)
	}
}

// A quebra de lista não pode picar PROSA. A primeira versão usava "linha longa
// com 3+ vírgulas" e quebrou a própria mensagem do `regra-implementada` em
// "a spec existe,\n o código existe,\n os dois se referenciam".
func TestProsaLongaComVirgulasNaoQuebra(t *testing.T) {
	prosa := "Uma regra catalogada sem implementador atravessa o pipeline inteiro — a spec existe, " +
		"o código existe, os dois se referenciam pelo header, e todos os gates ficam verdes " +
		"sobre trabalho que não foi feito."
	if got := quebraOcorrencias(prosa); strings.Contains(got, "\n") {
		t.Errorf("picou uma frase em orações:\n%s", got)
	}
}

// E continua quebrando lista de verdade — itens sem espaço interno.
func TestListaDeCaminhosAindaQuebra(t *testing.T) {
	itens := make([]string, 20)
	for i := range itens {
		itens[i] = fmt.Sprintf("apps/mobile/src/features/Modulo%02d.spec.md", i)
	}
	if got := quebraOcorrencias(strings.Join(itens, ", ")); !strings.Contains(got, "\n") {
		t.Errorf("lista de caminhos não quebrou:\n%s", got)
	}
}

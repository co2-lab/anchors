package initx

import (
	"embed"
	"strings"
	"testing"
)

// O CONTRATO DO `board.json` E UNICO, e o `board serve` le a expressao do PIPELINE em
// vez de reimplementa-la.
//
// Reimplementar criaria uma segunda versao que diverge no primeiro campo novo: o `owner`
// entra no pipeline e o board local nao o mostra -- ou pior, mostra diferente, e ninguem
// sabe qual dos dois esta certo.
func TestBoardCollectJQSaiDoPipeline(t *testing.T) {
	jq, err := BoardCollectJQ()
	if err != nil {
		t.Fatalf("nao extraiu a expressao: %v", err)
	}
	// Os campos que o board precisa — se o regex casar o trecho errado, algum sai.
	for _, campo := range []string{"number:", "title:", "state:", "owner:", "ownership:"} {
		if !strings.Contains(jq, campo) {
			t.Errorf("a expressao extraida nao tem `%s` — o regex casou o trecho errado, "+
				"e o board local serviria um contrato incompleto", campo)
		}
	}
	// O descartado sai na ORIGEM, e isso e regra do board: se ele saiu do board, saiu do
	// dado. O `serve` tem de herdar isso.
	if !strings.Contains(jq, "anchors:discarded") {
		t.Error("a expressao nao filtra `anchors:discarded` — o board local mostraria " +
			"card descartado, que o publicado nao mostra")
	}

	// A EXPRESSAO PRECISA SER JQ VALIDO, e a regua anterior nao media isso.
	//
	// A primeira versao do regex fechava em `'` nao-guloso e engolia o
	// `> _board/board.json` que vem depois no shell. Os campos todos apareciam -- as
	// asserções acima passavam -- e o `jq` morria com `unexpected token "'"` so quando o
	// comando rodava de verdade.
	//
	// E' a armadilha de conferir a FORMA e nao o que a forma serve para fazer.
	// O `>` sozinho NAO serve de sinal: ele aparece legitimamente na captura nomeada do
	// jq (`(?<c>...)`). A regua primeira versao acusava isso como defeito -- grosseira
	// demais, e teria mandado consertar o que estava certo.
	//
	// O que denuncia o corte errado e' o REDIRECIONAMENTO, que so existe no shell.
	if strings.Contains(jq, "_board/board.json") {
		t.Errorf("a expressao extraida carrega redirecionamento do shell — o jq vai "+
			"receber `>` como sintaxe e morrer:\n%s", ultimos(jq, 80))
	}
	if !strings.HasPrefix(strings.TrimSpace(jq), "[") || !strings.HasSuffix(strings.TrimSpace(jq), "]") {
		t.Errorf("a expressao nao abre e fecha o array — o regex cortou no lugar errado:\n%s",
			ultimos(jq, 80))
	}
}

// O HTML e o MESMO — nao uma copia que envelhece em paralelo.
func TestBoardHTMLEhOMesmoDoPipeline(t *testing.T) {
	h, err := BoardHTML()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(h, "board.json") {
		t.Error("o HTML nao busca `board.json` — o `serve` serviria uma pagina que nao le nada")
	}
	// O QUE ESTA REGUA MEDE e' que o `serve` le o HTML DO PRODUTO -- nao que ele tenha
	// uma correcao especifica.
	//
	// A primeira versao exigia `g-rot[data-number]`, que e' de outro PR ainda aberto. Ela
	// reprovava aqui por um motivo que nao tem a ver com o que ela diz medir, e forcaria
	// quem a lesse a acoplar dois trabalhos independentes.
	//
	// As marcas abaixo sao da ESTRUTURA do board, nao de uma correcao: se elas somem, o
	// `serve` passou a ler outro arquivo.
	for _, marca := range []string{"id=\"colunas\"", "id=\"detalhe\"", "id=\"v-roadmap\""} {
		if !strings.Contains(h, marca) {
			t.Errorf("o HTML nao tem `%s` — o `serve` esta lendo outro arquivo", marca)
		}
	}
}

// ultimos devolve o fim de uma string, para a mensagem de erro mostrar ONDE cortou.
func ultimos(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return "…" + s[len(s)-n:]
}

// Without the embedded files, both readers fail loudly instead of serving an empty board.
func TestBoardReadersFailWithoutTheEmbeddedFiles(t *testing.T) {
	savedBoard, savedWorkflows := boardFS, workflowsFS
	t.Cleanup(func() { boardFS, workflowsFS = savedBoard, savedWorkflows })
	boardFS, workflowsFS = embed.FS{}, embed.FS{}

	if h, err := BoardHTML(); err == nil || h != "" || !strings.Contains(err.Error(), "the board HTML is not embedded") {
		t.Errorf("BoardHTML = %q, %v; want the not-embedded error", h, err)
	}
	if jq, err := BoardCollectJQ(); err == nil || jq != "" || !strings.Contains(err.Error(), "the board pipeline is not embedded") {
		t.Errorf("BoardCollectJQ = %q, %v; want the not-embedded error", jq, err)
	}
}

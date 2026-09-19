package initx

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"
)

// BoardHTML devolve o HTML do board — o MESMO que o pipeline publica.
//
// Não é conveniência: duas renderizações do board divergiriam com o tempo, e aí ninguém
// saberia qual das duas está certa. O `board serve` serve este arquivo, não uma cópia.
func BoardHTML() (string, error) {
	b, err := fs.ReadFile(boardFS, "board/anchors-board.html")
	if err != nil {
		return "", fmt.Errorf("the board HTML is not embedded: %w", err)
	}
	return string(b), nil
}

// jqDaColeta extrai o `--jq` que o pipeline usa para montar cada item do board.
//
// O CONTRATO DO `board.json` É ÚNICO, e reimplementá-lo em Go criaria uma segunda
// versão que diverge no primeiro campo novo — o `owner` entra no pipeline, e o board
// local não o mostra; ou pior, mostra diferente.
//
// A âncora é o `--json` da coleta, que é único no arquivo.
// A ÂNCORA DO FIM importa tanto quanto a do começo. A primeira versão fechava em `'`
// não-guloso e engolia o `> _board/board.json` que vem depois — o jq recebia a
// redireção do shell como expressão e morria com `unexpected token "'"`.
//
// O fim real é `}]'` seguido do redirecionamento: o `]` fecha o array que o `[.[] | ...]`
// abriu, e o `'` fecha a aspa do shell.
var jqDaColeta = regexp.MustCompile(`(?s)--jq '(\[\.\[\].*?\}\])' *>`)

// BoardCollectJQ devolve a expressão jq da coleta, lida do pipeline embutido.
func BoardCollectJQ() (string, error) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-board.yml")
	if err != nil {
		return "", fmt.Errorf("the board pipeline is not embedded: %w", err)
	}
	m := jqDaColeta.FindSubmatch(b)
	if m == nil {
		return "", fmt.Errorf("could not find the collect expression in `anchors-board.yml` — " +
			"the pipeline changed shape, and `board serve` would read a contract different from " +
			"the one it publishes")
	}
	// O YAML indenta a expressão; o jq não se importa, mas a leitura fica melhor sem.
	return strings.TrimSpace(string(m[1])), nil
}

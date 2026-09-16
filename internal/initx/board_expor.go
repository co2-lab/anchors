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
		return "", fmt.Errorf("o HTML do board não está embutido: %w", err)
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
var jqDaColeta = regexp.MustCompile(`(?s)--json number,title,url,labels,comments,updatedAt,body,createdAt,closedAt,author \\\n\s*--jq '(.*?)'\s*\n`)

// BoardCollectJQ devolve a expressão jq da coleta, lida do pipeline embutido.
func BoardCollectJQ() (string, error) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-board.yml")
	if err != nil {
		return "", fmt.Errorf("o pipeline do board não está embutido: %w", err)
	}
	m := jqDaColeta.FindSubmatch(b)
	if m == nil {
		return "", fmt.Errorf("não achei a expressão de coleta no `anchors-board.yml` — " +
			"o pipeline mudou de forma, e o `board serve` leria um contrato diferente do " +
			"que ele publica")
	}
	// O YAML indenta a expressão; o jq não se importa, mas a leitura fica melhor sem.
	return strings.TrimSpace(string(m[1])), nil
}

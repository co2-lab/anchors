package scan

import "strings"

// --- o arquivo de PROGRESSO de um plano ---
//
// Um plano é DECISÃO; o progresso é ESTADO. Enquanto os dois viviam no mesmo arquivo,
// marcar uma fase como concluída era ALTERAR o plano — e isso fazia "terminei a fase 1"
// ser indistinguível de "mudei a direção do projeto", que é justamente a diferença que o
// `plano-alterado-justificado` existe para preservar.
//
// A separação vive aqui, na porta de entrada do mapa, porque é o mapa que precisa não
// vê-lo: um arquivo que existe para mudar não pode ser confrontado por gates que cobram
// justificativa de mudança. Ele VAI para o git (é o histórico do trabalho); o que não vai
// é para o mapa.
const sufixoProgresso = "-progress.md"

// EhArquivoDeProgresso diz se o caminho é o companheiro de estado de um plano.
func EhArquivoDeProgresso(caminho string) bool {
	return strings.HasSuffix(caminho, sufixoProgresso)
}

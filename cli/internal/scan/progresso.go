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

// IsProgressFile diz se o caminho é o companheiro de estado de um plano.
func IsProgressFile(caminho string) bool {
	return strings.HasSuffix(caminho, sufixoProgresso)
}

// ProgressPathFor devolve o caminho do companheiro de progresso de um plano.
//
// Existe aqui, e não em quem consome, para que o sufixo tenha UMA definição. O `scan` é
// quem precisa manter o arquivo fora do mapa; uma segunda constante em outro pacote
// poderia divergir desta em silêncio — e o consumidor passaria a procurar um arquivo que
// não existe, ou a confrontar um que o scanner indexa.
func ProgressPathFor(plano string) string {
	ext := ""
	if i := strings.LastIndex(plano, "."); i > strings.LastIndex(plano, "/") {
		ext = plano[i:]
	}
	return strings.TrimSuffix(plano, ext) + sufixoProgresso
}

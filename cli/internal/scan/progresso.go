package scan

import (
	"regexp"
	"strings"
)

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

// --- o MERGE de dois lados de um progress ---
//
// Medido no blue-eyes: TRÊS PRs consecutivos conflitaram no mesmo arquivo (#217, #221 e o
// seguinte), sempre pelo mesmo motivo — duas branches marcando checkboxes vizinhos do
// mesmo plano. O `anchors.graph.yaml` tem driver desde a v0.1.48 e não conflitou nenhuma
// vez no mesmo período; este arquivo não tinha, e a resolução à mão é sempre idêntica.
//
// A REGRA é uma só: `[x]` vence `[ ]`. Marcar concluído é um FATO — a spec foi entregue,
// o `anchors check` a carimbou, o PR mergeou. Desmarcar por merge apagaria o fato, e o
// gate `progress-honest` passaria a acusar um arquivo que existe como não feito.

// reItemProgresso casa uma linha marcável e separa a MARCA da identidade.
//
// A identidade é o texto depois do `- [x] `: dois lados que marcaram o mesmo alvo têm
// textos idênticos ali, porque o texto vem do plano e ninguém o reescreve ao marcar.
var reItemProgresso = regexp.MustCompile(`^(\s*)- \[([ xX])\] (.*)$`)

// MergeProgress une dois lados de um `*-progress.md`.
//
// A base do merge não é usada, pelo mesmo motivo do driver do mapa: a união não precisa
// saber o que havia antes. Não existe "desmarcar" legítimo vindo de um merge — quem quer
// desmarcar edita o arquivo, e essa edição chega como um lado só.
func MergeProgress(nosso, deles string) string {
	marcado := map[string]bool{}
	for _, l := range strings.Split(deles, "\n") {
		if m := reItemProgresso.FindStringSubmatch(l); m != nil && m[2] != " " {
			marcado[m[3]] = true
		}
	}

	// Percorre o NOSSO lado preservando ordem e formatação, promovendo a marca quando o
	// outro lado marcou. Sair do nosso lado (em vez de reconstruir) mantém tudo que não
	// é item — cabeçalho, prosa, seções — exatamente como está.
	var out []string
	visto := map[string]bool{}
	for _, l := range strings.Split(nosso, "\n") {
		m := reItemProgresso.FindStringSubmatch(l)
		if m == nil {
			out = append(out, l)
			continue
		}
		visto[m[3]] = true
		if m[2] == " " && marcado[m[3]] {
			l = m[1] + "- [x] " + m[3]
		}
		out = append(out, l)
	}

	// O que existe SÓ do outro lado entra depois do último item.
	//
	// No fim, e não no ponto "certo": não há ponto certo derivável — a ordem é a do
	// plano, que não está neste arquivo. Perder o item seria pior que posicioná-lo mal,
	// e quem o reordena é o `anchors new progress`.
	var faltando []string
	for _, l := range strings.Split(deles, "\n") {
		if m := reItemProgresso.FindStringSubmatch(l); m != nil && !visto[m[3]] {
			faltando = append(faltando, l)
			visto[m[3]] = true
		}
	}
	if len(faltando) == 0 {
		return strings.Join(out, "\n")
	}

	corte := -1
	for i := len(out) - 1; i >= 0; i-- {
		if reItemProgresso.MatchString(out[i]) {
			corte = i
			break
		}
	}
	if corte < 0 {
		return strings.Join(append(out, faltando...), "\n")
	}
	res := make([]string, 0, len(out)+len(faltando))
	res = append(res, out[:corte+1]...)
	res = append(res, faltando...)
	res = append(res, out[corte+1:]...)
	return strings.Join(res, "\n")
}

// ProgressDone conta os itens marcados — o número que o driver reporta ao git.
func ProgressDone(s string) int {
	n := 0
	for _, l := range strings.Split(s, "\n") {
		if m := reItemProgresso.FindStringSubmatch(l); m != nil && m[2] != " " {
			n++
		}
	}
	return n
}

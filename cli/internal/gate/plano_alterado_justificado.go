package gate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- o plano ALTERADO diz por que mudou ---
//
// Planejar erra, e quem descobre o erro é quem implementa. O problema não é o erro: é o
// que acontece depois. Um agente que acha inconsistência no plano tem duas saídas, e as
// duas são ruins — corrigir por conta própria (e o projeto passa a caminhar para um
// destino que ninguém escolheu) ou implementar o que está escrito sabendo que está errado.
//
// A deriva é o risco maior, e é SILENCIOSA por construção: nenhum gate de ESTADO consegue
// vê-la, porque o plano corrigido fica perfeitamente válido — a inconsistência foi
// removida. O que denuncia não é o estado do arquivo, é a MUDANÇA sem justificativa.
//
// Por isso este gate olha o diff, não o conteúdo: plano ou spec que aparece entre os
// arquivos alterados tem de trazer uma revisão declarada. Sem ela, barra.
//
// A régua é o julgamento de quem alterou, e ele é explícito:
//
//   - correção INÓCUA (redação, exemplo, typo, uma ambiguidade que só tinha uma leitura
//     possível) — corrige e registra a revisão. O gate confere que a revisão existe.
//
//   - correção que MUDA A DIREÇÃO, ou dúvida sobre se muda — não corrige. Abre issue com
//     `anchors:precisa-do-usuario`, e o `claim` para de entregar o card até alguém decidir.
//
// O gate não sabe distinguir os dois casos, e não é para saber: essa é a decisão que se
// quer que um humano ou um agente TOME, com o contexto na mão. O que ele garante é que a
// decisão foi tomada por alguém e ficou escrita — em vez de acontecer por omissão.

// revisaoRE casa a revisão registrada no arquivo: `FNDTN-R0001: o que mudou e por quê`.
//
// O formato segue o vocabulário que já existe (`FNDTN-F04` para fase), e a NUMERAÇÃO é o
// que uma marca solta não daria: dá para ver quantas vezes o documento mudou, e em que
// ordem. Um `@plan-fix` solto responderia "mudou"; `-R0003` responde "mudou três vezes".
func revisaoRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^[^\S\n]*>?[^\S\n]*(?:\*\*)?([A-Z0-9]` +
		config.CodeLengthPattern() + `)-R(\d{4})(?:\*\*)?[^\S\n]*:[^\S\n]*(\S.*)$`)
}

// Revisao é uma alteração registrada no próprio documento.
type Revisao struct {
	Codigo     string // o código do arquivo revisado (`FNDTN`)
	Numero     int    // sequencial: 1, 2, 3...
	Explicacao string
}

// RevisoesDe devolve as revisões declaradas no conteúdo, na ordem em que aparecem.
func RevisoesDe(content string) []Revisao {
	var out []Revisao
	for _, m := range revisaoRE().FindAllStringSubmatch(content, -1) {
		n, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		out = append(out, Revisao{Codigo: m[1], Numero: n, Explicacao: strings.TrimSpace(m[3])})
	}
	return out
}

// checkPlanoAlteradoJustificado confronta o plano/spec ALTERADO com a revisão declarada.
//
// O gate se abstém em `--all` (via `skip_on: [all]` no anchors.yaml). Ali não existe
// "alterado": reprovar todo plano que nunca precisou de revisão seria acusar quem acertou
// de primeira. Rodando com `--changed`, todo nó que ele recebe JÁ é um arquivo alterado —
// por isso não precisa da lista, e não a recebe.
func checkPlanoAlteradoJustificado(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	// SÓ o que de fato MUDOU. O `--changed X` entrega o RAIO DE IMPACTO de X — todo nó
	// que depende dele —, e isso é certo para quase todo gate: quem quebrou por tabela tem
	// de ser confrontado. Aqui não: um plano que não mudou não tem o que justificar.
	//
	// Medido no blue-eyes: sem esta conferência, alterar UM plano acusava 8 arquivos, 7
	// deles intocados. Um gate bloqueante que acusa inocente é pior que gate nenhum — a
	// saída barata vira desligá-lo.
	if cfg == nil || !mudouDeFato(n.ID, cfg.Alterados) {
		return Skip, "não está entre os arquivos alterados (só foi alcançado pelo raio de impacto)"
	}

	codigo := n.Code
	if codigo == "" {
		return Skip, "o arquivo não tem código — quem cobra isso é o `spec-tem-codigo`"
	}

	// Só contam as revisões DESTE documento. Um plano pode citar a revisão de outro ao
	// explicar o contexto, e isso não justifica a própria mudança.
	var minhas []Revisao
	for _, r := range RevisoesDe(content) {
		if r.Codigo == codigo {
			minhas = append(minhas, r)
		}
	}

	if len(minhas) == 0 {
		return Fail, fmt.Sprintf(
			"foi ALTERADO e não diz por quê. Planejar erra, e quem descobre o erro é quem "+
				"implementa — mas corrigir o plano em silêncio faz o projeto caminhar para um "+
				"destino que ninguém escolheu, e nenhum gate de estado vê isso (o plano "+
				"corrigido fica válido). Registre a revisão no próprio arquivo:\n"+
				"    > **%s-R0001:** <o que mudou e por quê>\n"+
				"Se a correção MUDA A DIREÇÃO do projeto, ou se você tem dúvida se muda, NÃO "+
				"corrija — a avaliação do impacto é sua, e esta saída existe para quando "+
				"ela dá `muda`:\n"+
				"    anchors escalate \"<o que está incoerente>\" --sobre %s --card <n>", codigo, n.ID)
	}

	// A numeração tem de ser sequencial a partir de 1. Sem isso ela não responderia
	// "quantas vezes mudou" — que é a única coisa que ela dá a mais que uma marca solta.
	maior := 0
	for _, r := range minhas {
		if r.Numero > maior {
			maior = r.Numero
		}
	}
	if maior != len(minhas) {
		return Fail, fmt.Sprintf(
			"as revisões não são sequenciais: há %d declarada(s) e a maior é `-R%04d`. A "+
				"numeração é o que diz QUANTAS vezes o documento mudou; com buraco ou "+
				"repetição ela deixa de responder isso.", len(minhas), maior)
	}

	ult := minhas[len(minhas)-1]
	return Pass, fmt.Sprintf("alterado, e a revisão `%s-R%04d` diz por quê: %s",
		ult.Codigo, ult.Numero, primeiraLinha(ult.Explicacao))
}

// primeiraLinha encurta a explicação para o laudo, que é uma linha.
func primeiraLinha(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 90 {
		return s[:87] + "..."
	}
	return s
}

// mudouDeFato diz se o nó está na lista dos que mudaram.
func mudouDeFato(id string, alterados []string) bool {
	for _, a := range alterados {
		if a == id {
			return true
		}
	}
	return false
}

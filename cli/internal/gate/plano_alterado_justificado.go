package gate

import (
	"fmt"
	"os/exec"
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

	// E O GIT CONFIRMA. `--changed X` é uma AFIRMAÇÃO de quem chama, não uma medição:
	// quem roda à mão passa o caminho para testar algo, e o pre-commit passa TODO arquivo
	// staged — inclusive os que entraram por rebase ou merge sem ninguém os ter editado.
	//
	// Cobrar justificativa de quem não mexeu é o pior defeito possível num gate
	// bloqueante: ele barra trabalho correto, e a saída barata vira desligá-lo.
	// ARQUIVO NOVO não tem o que justificar: ele não existia, então nada foi ALTERADO.
	// A revisão registra por que o texto mudou — e num arquivo que nasce agora, o texto
	// inteiro é a decisão. Cobrar `-R0001` aqui obrigaria toda spec nova a declarar uma
	// revisão de si mesma no primeiro commit, que é ruído puro.
	if gitDizQueEhNovo(root, n.ID) {
		return Skip, "arquivo novo — não há alteração a justificar"
	}
	if !gitDizQueMudou(root, n.ID) {
		return Skip, "o git não vê mudança neste arquivo — `--changed` o incluiu, mas o " +
			"conteúdo é igual ao do último commit"
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

	// JÁ SE EXPLICA por outro mecanismo? Então não há o que cobrar.
	//
	// Medido no PR do plano de mutação: os dois planos alterados diziam por que mudaram —
	// o revisado com `@revised-by`/`@amended-by`, o que revisa com `revises:` no header —
	// e este gate reprovou os dois, exigindo que dissessem de novo em outra notação.
	//
	// Exigir a mesma informação duas vezes não protege nada: ensina a satisfazer o gate
	// em vez de comunicar, que é o oposto do que ele existe para fazer.
	if seExplicaPorRevisao(content) {
		return Pass, "a mudança já está explicada pelo mecanismo de revisão de planos " +
			"(`revises:` / `@revised-by`), e cobrar a mesma coisa em duas notações não " +
			"protegeria nada"
	}

	if len(minhas) == 0 {
		return Fail, fmt.Sprintf(
			"foi ALTERADO e não diz por quê. Planejar erra, e quem descobre o erro é quem "+
				"implementa — mas corrigir o plano em silêncio faz o projeto caminhar para um "+
				"destino que ninguém escolheu, e nenhum gate de estado vê isso (o plano "+
				"corrigido fica válido). Registre a revisão no próprio arquivo:\n"+
				"    > **%s-R0001:** <o que mudou e por quê>\n"+
				"Se a mudança IMPACTA A DIREÇÃO do projeto, ou se você tem dúvida, não a "+
				"faça aqui — a interpretação do impacto é sua, e ela escolhe a saída:\n"+
				"    anchors escalate \"<o que precisa mudar>\" --sobre %s --para-usuario\n"+
				"Se não impacta a direção mas também não é para agora, vira card comum:\n"+
				"    anchors escalate \"<o que precisa mudar>\" --sobre %s", codigo, n.ID, n.ID)
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

// seExplicaPorRevisao diz se o arquivo já declara a mudança pelo mecanismo de revisão
// entre planos — o `revises:` de quem revisa, e o aviso de quem foi revisado.
//
// São marcadores ESTÁVEIS, não prosa: o `plano-revisado` já os usa, e casar texto corrido
// quebraria em projeto escrito noutro idioma.
func seExplicaPorRevisao(content string) bool {
	for _, marca := range []string{"revises:", "@revised-by", "@amended-by"} {
		if strings.Contains(content, marca) {
			return true
		}
	}
	return false
}

// gitDizQueMudou confronta a lista recebida com o que o git de fato vê.
//
// Conta o que está no índice E na árvore de trabalho: o pre-commit roda com o arquivo já
// staged, e olhar só um dos dois deixaria passar metade dos casos.
//
// Sem git (ou fora de repositório), devolve `true` e deixa a decisão com quem chamou —
// negar ali silenciaria o gate onde ele não tem como medir.
func gitDizQueMudou(root, path string) bool {
	cmd := exec.Command("git", "status", "--porcelain", "--", path)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return len(strings.TrimSpace(string(out))) > 0
}

// gitDizQueEhNovo diz se o arquivo ainda não existe no histórico.
//
// `git log -1 -- <path>` vazio significa que nenhum commit o tocou — é a diferença entre
// "mudou" e "nasceu". O status porcelain não serve aqui: ele marca `??` para não
// rastreado e `A ` para staged, e um arquivo novo já adicionado ao índice apareceria como
// alteração.
func gitDizQueEhNovo(root, path string) bool {
	// `git ls-files` e não `git log`: num repositório SEM NENHUM COMMIT o log falha
	// ("does not have any commits yet") em vez de devolver vazio, e tratar o erro como
	// "não é novo" fazia o gate cobrar justificativa de todo arquivo do primeiro commit —
	// justamente onde nada pode ter sido alterado.
	//
	// `ls-files` responde o que interessa: o arquivo é RASTREADO? Se não é, ele nasce
	// agora, com ou sem histórico no repositório.
	//
	// FORA de repositório a resposta é NÃO: ali o gate não tem como medir, e afirmar
	// "é novo" o silenciaria em todo projeto sem git — que é o caso dos testes de unidade
	// e de quem roda o Anchors fora de um repositório.
	if !emRepositorio(root) {
		return false
	}
	cmd := exec.Command("git", "ls-files", "--error-unmatch", "--", path)
	cmd.Dir = root
	return cmd.Run() != nil // erro = não rastreado = nasce agora
}

// emRepositorio diz se `root` está dentro de um repositório git.
func emRepositorio(root string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = root
	return cmd.Run() == nil
}

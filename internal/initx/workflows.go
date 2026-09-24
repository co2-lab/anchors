package initx

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// workflowsFS carrega os pipelines que o modo `github` do fluxo de trabalho pressupõe.
//
// Viajam no binário e são COPIADOS para `.github/workflows/` — a mesma régua dos packs.
// A diferença importa: no projeto, o pipeline é material, versionado e revisável num
// diff. Um time que precise ajustar o ritmo do `stale`, trocar a label ou mudar a
// permissão edita o arquivo; um pipeline que vivesse só no binário obrigaria a esperar um
// release para mudar o que é decisão do time.
//
//go:embed workflows
var workflowsFS embed.FS

// boardFS carrega a PÁGINA do board. Fica fora de `workflows/` porque não é um pipeline:
// o GitHub executa tudo que está em `.github/workflows/`, e um HTML ali seria um arquivo
// de workflow inválido — o repositório passa a exibir um erro de sintaxe permanente.
//
//go:embed board
var boardFS embed.FS

// BoardFile é onde a página mora, relativa à raiz do projeto. Em `.github/` e não em
// `docs/`: é infraestrutura do fluxo, não documentação do produto, e misturá-la com o que
// o time escreve convida a alguém a editá-la sem saber que o `--fix` a mantém.
const BoardFile = ".github/anchors-board.html"

// Workflow é um pipeline do fluxo GitHub, com o que o Anchors precisa saber para
// verificar (doctor) e semear (doctor --fix).
type Workflow struct {
	// Arquivo é o nome em `.github/workflows/`.
	Arquivo string
	// Papel é o que ele faz, para a mensagem do doctor dizer o que fica sem acontecer.
	Papel string
	// ExigeSerial diz se este pipeline PRECISA de `concurrency` sem cancelamento. Para
	// os que atribuem ou criam cards, a serialização não é otimização: é o que impede
	// duas execuções de atribuírem o mesmo card ou criarem o card duas vezes.
	ExigeSerial bool
}

// WorkflowsDoFluxo são os pipelines que o modo `github` pressupõe. A lista é a fonte da
// verdade tanto do que o doctor confere quanto do que o `--fix` semeia — uma lista só,
// para que verificar e consertar nunca discordem sobre o que deveria existir.
var WorkflowsDoFluxo = []Workflow{
	{
		Arquivo:     "anchors-identify.yml",
		Papel:       "creates the card for every artifact that reached the repository without one",
		ExigeSerial: true,
	},
	{
		Arquivo: "anchors-gates.yml",
		Papel:   "confronts the gates on every PR (it is the boundary: nothing lands without passing)",
		// NÃO serializa: só lê e reporta. Serializar faria cada PR esperar a fila dos
		// outros sem necessidade — o `concurrency` dele é por PR, não global.
		ExigeSerial: false,
	},
	{
		Arquivo:     "anchors-claim.yml",
		Papel:       "assigns work to the agents (it is what removes the race for a card)",
		ExigeSerial: true,
	},
	{
		Arquivo: "anchors-pr-checks.yml",
		Papel: "moves the card to review when the PR checks pass (and ONLY when " +
			"they pass)",
		ExigeSerial: true,
	},
	{
		Arquivo: "anchors-resolve-queue.yml",
		Papel: "resolves the GENERATED-file conflict when something merges, and opens a synthesis " +
			"card when the conflict is about content",
		// SERIAL, e a razão é o push: dois runs resolvendo a mesma branch empurrariam
		// resoluções concorrentes. O grupo é global de propósito — diferente do
		// `pr-checks`, aqui o recurso disputado é a FILA inteira, não um card.
		ExigeSerial: true,
	},
	{
		Arquivo: "anchors-board.yml",
		Papel:   "publishes the board on Pages from the issues (no Projects and no PAT)",
		// O ÚNICO que não exige serialização sem cancelamento: o board é estado DERIVADO,
		// e a execução mais nova sempre produz uma foto melhor que a que está no meio do
		// caminho. Cancelar aqui não perde trabalho — evita publicar uma foto velha por
		// cima de uma nova.
		ExigeSerial: false,
	},
	{
		Arquivo:     "anchors-stale.yml",
		Papel:       "releases cards whose owner vanished, preserving the history",
		ExigeSerial: true,
	},
	{
		Arquivo: "anchors-guard.yml",
		Papel:   "undoes the state change that did not come from the flow",
		// SERIAL, e por card: uma reversão escreve label, e o evento que ela mesma dispara
		// chegaria antes de a primeira execução terminar.
		ExigeSerial: true,
	},
	{
		Arquivo: "anchors-decided.yml",
		Papel:   "devolve à fila o card cujos desbloqueios foram todos entregues",
		// SERIAL: ele escreve label, e duas execuções sobre o mesmo card — uma vinda do
		// comentário, outra do cron — removeriam a mesma label duas vezes e comentariam
		// duas vezes no card.
		ExigeSerial: true,
	},
}

// BranchProtection é o que o modo `github` exige da `main`: nada entra sem PR.
//
// É o que torna o ciclo de revisão possível. A §7.9 do BOOTSTRAP diz que o agente sobe o
// código e ABRE O PR, e o card vai para `ready-to-review` — sem PR não há o que revisar,
// e o estado `in-review` fica sem objeto.
//
// E é o que o pipeline de identificação pressupõe: ele dispara na abertura do PR, porque
// o push na main acontece DEPOIS do merge, quando o trabalho já terminou. Push direto na
// main pula o card, pula a revisão, e pula o pipeline.
type BranchProtection struct {
	// ExigePR — nada entra na main sem pull request.
	ExigePR bool
	// RevisoesNecessarias — quantas aprovações. Zero é legítimo num time de uma pessoa
	// com agentes: o PR existe para o card ter objeto e para o histórico ficar legível,
	// e exigir aprovação de outra conta travaria o fluxo inteiro.
	RevisoesNecessarias int
}

// RequiredProtection é o mínimo que o fluxo pressupõe.
var RequiredProtection = BranchProtection{ExigePR: true, RevisoesNecessarias: 0}

// DirWorkflows é onde os pipelines moram no projeto. One constant with the scan, which
// reads a marker-carrying file here as upstream-owned: two copies would drift.
const DirWorkflows = scan.UpstreamDir

// WorkStates são as LABELS que carregam o estado de um card, na ordem do fluxo.
//
// O estado é uma LABEL, e não a coluna do Project (ver BOOTSTRAP.md §7.13). A escolha
// anterior foi a coluna, e ela cobrava um preço que só apareceu no uso: escrever num
// Project de organização exige um PAT com escopo `project`, que o `GITHUB_TOKEN` da
// Action não tem — então todo projeto que adotasse o fluxo precisaria criar e manter um
// token pessoal antes de o primeiro card se mover. Atrito de adoção por uma decisão de
// visualização.
//
// Com label, o `GITHUB_TOKEN` basta e nada precisa ser configurado. E o board não some:
// o GitHub Projects tem automação nativa que move o card quando a label muda — a
// sincronia acontece do lado deles, e só UM lado escreve (nós na label, eles no board),
// que era a preocupação original de ter estado em dois lugares.
//
// O par `ready-to-x` / `in-x` é o que torna a fila legível: um diz "disponível para
// alguém pegar", o outro "alguém está fazendo".
var WorkStates = []string{
	LabelToDo,
	"anchors:in-progress",
	"anchors:ready-to-review",
	"anchors:in-review",
	"anchors:ready-to-test",
	"anchors:in-test",
	"anchors:ready-to-release",
	"anchors:production",
}

// PrefixoLabelSob liga um card ao trabalho de onde ele NASCEU: `anchors:under-44`.
//
// O achado que aparece enquanto se implementa outra coisa precisa de duas coisas ao mesmo
// tempo — existir por si (para não se perder) e estar amarrado ao trabalho em curso (para
// ser entregue junto). Sem a amarra, ele vira um card solto que ninguém relaciona; sem a
// existência própria, vira uma frase no corpo de outra issue.
//
// É LABEL, e não texto no corpo. A primeira versão escrevia "Descoberto durante o card
// #44" na descrição, e isso não se consulta: não dá para listar o que pende sob um card,
// nem para o board desenhar a relação. Label é filtrável (`--label anchors:under-44`),
// aparece na lista de issues e sobrevive a qualquer reescrita do texto.
//
// Também NÃO é a sub-issue nativa do GitHub: ela só aceita um nível, e a hierarquia deste
// fluxo tem mais (plano → fase → spec → achado). A árvore do board já se monta pelo
// `parent:` do artefato; esta label é o que amarra o achado que NÃO tem artefato — uma
// config, um pipeline, um arquivo que nenhuma spec governa.
const PrefixoLabelSob = "anchors:under-"

// PrefixoLabelSobAntigo é o nome anterior, em português.
//
// Ele está em ISSUES do GitHub, não só em configuração: renomear a constante não renomeia
// as labels que já existem. Os pipelines aceitam os dois enquanto durar a migração.
//
// NÃO HÁ RENOMEAÇÃO AUTOMÁTICA, e este comentário afirmava que havia ("o `anchors doctor
// --fix` renomeia as labels no board"). Não renomeia, nunca renomeou — a promessa saiu em
// vez de ganhar uma implementação porque aceitar as duas grafias já basta: renomear label
// no GitHub reescreve o histórico de quem a usou, e o ganho seria cosmético.
//
// Quem LÊ as duas: os pipelines, e o `anchors backfill-labels`, que recupera vínculo de
// card antigo.
const PrefixoLabelSobAntigo = "anchors:sob-"

// LabelSob devolve a label que liga um card ao trabalho de origem.
func LabelSob(card string) string { return PrefixoLabelSob + card }

// LabelNeedsUser marca o card que ESPERA UMA PESSOA.
//
// Não é um estado do fluxo, e por isso não entra em `EstadosDoTrabalho`: o card continua
// onde está (`in-review`, `to-do`), e o que muda é QUEM pode destravá-lo. Fosse estado,
// um card escalado sairia da coluna onde o trabalho realmente parou, e o board deixaria
// de mostrar onde o fluxo travou.
//
// É a mesma distinção que `issue.DonoUsuário` faz para as issues em `issues/`: o dono é
// um eixo independente do estado.
const LabelNeedsUser = "anchors:needs-user"

// LabelNeedsFraming marca o card em que a dúvida é o ENQUADRAMENTO, não o mérito.
//
// São duas filas do usuário, com pesos diferentes, e misturá-las faz a mais barata custar
// como a mais cara:
//
//	needs-user     "decida entre A e B" — existe mais de uma resposta defensável, e
//	               escolher entre elas muda o que o produto faz
//	needs-framing  "confira se isto é seu" — quem escalou não soube dizer se muda a
//	               direção, e preferiu declarar a dúvida a afirmar impacto que não mediu
//
// A SEGUNDA TEM SAÍDA BARATA: se não impacta, quem lê troca a label por `to-do` e o card
// volta à fila. Ninguém precisa decidir o mérito. Sem a distinção, quem abre a fila gasta
// o esforço de decidir antes de descobrir que só precisava devolver o card.
//
// MEDIDO no projeto de referência: das 26 decisões abertas, uma triagem concluiu que eram
// SETE decisões reais. Boa parte do resto era enquadramento — card aberto como decisão
// porque o agente teve dúvida, e a doutrina do `escalate` dizia "use `--for-user` (...) ou
// se você tem dúvida se impacta", convidando a escalar por segurança.
//
// NÃO É ESTADO, como o `needs-user`: o card continua na coluna onde o trabalho parou.
const LabelNeedsFraming = "anchors:needs-framing"

// PrefixoLabelDesbloqueia marca o card cuja ENTREGA destrava outro.
//
// `needs-user` diz que um card espera uma PESSOA, e o claim já não o entrega enquanto a
// label estiver lá. O que faltava era o caso em que a decisão da pessoa GERA TRABALHO: ela
// decide, e a decisão exige que alguém mude alguma coisa antes de o card original seguir.
//
// Sem o vínculo, o que acontece é isto: a pessoa abre um card novo para a mudança, e o card
// bloqueado fica com `needs-user` para sempre — porque a instrução diz "decida e remova a
// label", e ela decidiu mas o trabalho ainda não foi feito. Ou pior: ela remove a label
// achando que decidir bastava, o card volta à fila, e o agente que o pega encontra o mesmo
// impasse.
//
// A label liga os dois: `anchors:desbloqueia-311` num card diz "quando eu for entregue, o
// #311 pode voltar à fila". O pipeline de board a lê e o `anchors next` a respeita.
const PrefixoLabelDesbloqueia = "anchors:desbloqueia-"

// MarcadorDeReversao abre o comentário que a trava de estado escreve ao desfazer uma
// mudança manual.
//
// Existe como CONSTANTE porque duas pontas precisam dele: o pipeline o escreve, e o
// `anchors task-status` o procura para dizer ao agente que sua mudança foi revertida.
//
// A razão de a segunda ponta existir: um agente fechou o card à mão, a trava reverteu no
// mesmo minuto, e ele escreveu "issue closed e resolvida, nada mais a fazer" — sem saber.
// O comentário estava lá e estava correto; quem já saiu da conversa não o lê.
const MarcadorDeReversao = "🔒"

// LabelManual é o OPT-OUT da trava: o card pode ser movido à mão.
//
// O pipeline `anchors-guard` desfaz qualquer mudança de estado que não tenha vindo dele —
// e precisa disso porque os agentes usam a MESMA conta que a pessoa: o `actor` de um evento
// distingue `github-actions[bot]` de humano, mas não distingue agente de dono do projeto.
//
// A label é o ato deliberado que diz "este card sou eu que movo". Ela fica NO CARD, à
// vista, e some do histórico só quando alguém a remove — diferente de uma frase em
// comentário, que se perde no meio da conversa.
//
// O CUSTO é pôr a label antes de mexer, e ele é o ponto: mover um card é raro, e o que a
// trava impede é justamente o movimento que ninguém pensou duas vezes antes de fazer.
const LabelManual = "anchors:manual"

// LabelDiscarded marca o card que não faz mais sentido.
//
// SOFT-DELETE, e a escolha é deliberada. Fechar não basta: o board mostra os fechados porque
// o roadmap precisa deles — é como ele desenha o que já foi entregue. Apagar de verdade
// perderia o rastro de que a pergunta existiu, e é isso que distingue "resolvido" de
// "descartado".
//
// MEDIDO no projeto de referência: das 27 raízes do roadmap, DOZE eram ruído permanente —
// cards de teste (`[teste] … apagar`) e achados sobre arquivos que não existem mais. Todos
// fechados, todos ainda desenhados, e nenhum jamais terá um pai a encontrar.
//
// O QUE ELA FAZ: o card sai do board (colunas, árvore, roadmap, faixa de pendências) e
// continua no GitHub, com a razão registrada em comentário. Quem procurar o número o acha.
const LabelDiscarded = "anchors:discarded"

// LabelDesbloqueia é a label do card que destrava `card`.
func LabelDesbloqueia(card string) string { return PrefixoLabelDesbloqueia + card }

// PrefixoLabelBlockedBy marca o card que NÃO PODE SER TRABALHADO até outro sair.
//
// É a direção que faltava. As duas labels existentes ligam os cards, e nenhuma para o
// trabalho de quem ficou esperando:
//
//	anchors:under-44          o card NASCEU do trabalho do #44 — procedência
//	anchors:desbloqueia-311   a ENTREGA deste card destrava o #311 — do destravador
//	anchors:blocked-by-44     ESTE card espera o #44 — do bloqueado
//
// O QUE FALTAVA: o `escalate --for-user` já PARA o card de origem — põe `needs-user` nele
// e comenta "⏸ Parado". O que ele não dizia é POR QUAL card ele espera.
//
// Sem o número, o board mostra "esperando você" sem vínculo, e quem responde a decisão não
// sabe o que acabou de soltar: cada card tem de ser reencontrado à mão, e o que não for
// reencontrado segue parado depois de a decisão já ter saído.
//
// MEDIDO no projeto de referência: 23 decisões abertas eram SETE perguntas, e destravar as
// dependentes exigia reler card por card para descobrir quem esperava o quê.
//
// O QUE A LABEL FECHA, nas três pontas:
//
//   - o `claim` confere se o bloqueador ainda está aberto antes de servir o card — a
//     ninguém, nem ao dono atual
//   - o board desenha o vínculo, com o número de quem segura
//   - quem decide lista tudo que a resposta libera (`--label anchors:blocked-by-<n>`)
//
// É AUTOMÁTICA e não flag: se o card parou por causa daquela decisão, o vínculo é fato.
// O julgamento — se a decisão impede o trabalho — já foi feito quando o agente escolheu
// `--for-user`.
const PrefixoLabelBlockedBy = "anchors:blocked-by-"

// PrefixoLabelDePR liga o achado ao PULL REQUEST em que ele foi visto.
//
// É procedência, como o `under-<n>`, e os dois COEXISTEM porque respondem perguntas
// diferentes:
//
//	anchors:under-198     o achado pertence ao trabalho do card #198 — é por onde se
//	                      entrega junto, e é o que o `claim` e o `pr-body` leem
//	anchors:from-pr-556   foi lendo o PR #556 que alguém viu — é por onde se rastreia a
//	                      REVISÃO, e ele sobrevive ao PR ser revertido ou reescrito
//
// O card #647 tem os dois na história e só um virou label: a prosa diz "ao revisar o PR
// #556" e a label diz `under-198` — que é o card que aquele PR fecha. Guardar só o card
// perde por onde o achado apareceu; guardar só o PR perde onde ele se entrega.
//
// NÃO É REDUNDANTE com o `Refs` do PR. O PR declara a issue dele, e é daí que o `escalate`
// DERIVA o card quando recebe `--reviewing-pr`. Mas derivar é ler uma vez: se o PR for
// reescrito depois, a declaração muda e o vínculo histórico se perde. A label registra o
// que foi lido, no momento em que foi lido.
const PrefixoLabelDePR = "anchors:from-pr-"

// LabelDePR é a label que liga o achado ao PR revisado.
func LabelDePR(pr string) string { return PrefixoLabelDePR + pr }

// LabelBlockedBy é a label do card que espera `card`.
func LabelBlockedBy(card string) string { return PrefixoLabelBlockedBy + card }

// A DECISÃO é registrada por COMANDO, não por tag em comentário.
//
// Houve aqui uma `TagDeDecisao = "#solution"`, para um pipeline reconhecer a resposta do
// usuário num comentário. Ela saiu porque o `anchors decided --card N --resolution "..."`
// já existia e faz mais: exige a REVISÃO que nasceu da decisão, e uma decisão que libera o
// card sem dizer qual regra nasceu dela fica sem rastro.
//
// A tag teria criado um segundo caminho — o mais barato dos dois, e o que não exige a
// resolução. Dois jeitos de destravar o mesmo card, um deles sem rastro.

// LabelToDo é o estado em que um card NASCE.
//
// Constante e não literal pelo mesmo motivo do `LabelNeedsUser`: quem cria a issue
// (`internal/issue`) e quem cria o label (`anchors init`) precisam concordar, e um
// literal repetido nos dois lados foi exatamente como o `needs-user` divergiu — o `gh`
// recusa o comando inteiro por um label inexistente, então o achado do gate deixa de ser
// registrado sem que nada além do aviso apareça.
const LabelToDo = "anchors:to-do"

// LabelNeedsUserLegacy é o nome anterior. Ver PrefixoLabelSobAntigo.
const LabelNeedsUserLegacy = "anchors:precisa-do-usuario"

// ColunasDoBoard são os nomes das colunas do Project que ESPELHAM os estados acima.
// O board é opcional: quem o quiser cria as colunas com estes nomes e liga a automação
// nativa do Projects (label adicionada → move para a coluna). Quem não quiser trabalha
// só com issues, e o fluxo funciona igual.
var ColunasDoBoard = []string{
	"TO DO",
	"IN PROGRESS",
	"READY TO REVIEW",
	"IN REVIEW",
	"READY TO TEST",
	"IN TEST",
	"READY TO RELEASE",
	"PRODUCTION",
}

// ColunasDisponiveis são as colunas de onde um agente TIRA trabalho — as `READY TO ...`
// e o `TO DO` inicial, em ORDEM DE PRIORIDADE: da direita para a esquerda do board.
//
// A regra é essa e vale para todo o fluxo: o trabalho mais ADIANTADO vem primeiro.
// Terminar o que está quase pronto antes de começar coisa nova é o que impede o board de
// encher de trabalho pela metade — e trabalho pela metade não entrega nada, ocupa revisor,
// e envelhece até o contexto de quem o escreveu se perder.
//
// Daí a ordem invertida em relação a ColunasDoBoard: `READY TO REVIEW` (mais à direita)
// antes de `TO DO` (mais à esquerda).
var ColunasDisponiveis = []string{
	"READY TO REVIEW",
	"TO DO",
}

// AnchorsFinalState é o último estado que o Anchors ESCREVE.
const AnchorsFinalState = "anchors:ready-to-test"

// RecyclableStates são os únicos estados de onde o `stale` tira um card.
//
// A distinção é entre trabalho TRAVADO e trabalho ESPERANDO. Um card em `in-progress`
// sem sinal de vida é um agente que morreu — reciclar devolve o trabalho à fila. Um card
// em `ready-to-review` está parado POR DEFINIÇÃO: ele terminou e aguarda um revisor
// humano, e ninguém comenta enquanto isso. Reciclá-lo joga trabalho pronto de volta em
// `to-do`, como se nunca tivesse sido feito — e o agente seguinte o refaz.
//
// Só entram aqui os estados de trabalho ATIVO, aqueles em que alguém deveria estar
// mexendo agora. Os `ready-to-*` são filas de espera, e esperar não é estar travado.
var RecyclableStates = []string{
	"anchors:in-progress",
	"anchors:in-review",
	"anchors:in-test",
}

// AvailableStates são os estados de onde um agente TIRA trabalho, em ORDEM DE
// PRIORIDADE: da direita para a esquerda do fluxo. O trabalho mais ADIANTADO vem
// primeiro — terminar o que está quase pronto antes de começar coisa nova é o que impede
// o board de encher de trabalho pela metade.
var AvailableStates = []string{
	"anchors:ready-to-review",
	"anchors:to-do",
}

// ColunaFinalDoAnchors é a última coluna que o Anchors ESCREVE. Da seguinte em diante
// (`IN TEST`, `READY TO RELEASE`, `PRODUCTION`) quem move são os pipelines de entrega do
// projeto — cada time tem o seu, e o Anchors não tem o que dizer sobre quando um teste de
// aceitação passou ou um deploy aconteceu. Ele continua LENDO essas colunas (é o que o
// `anchors status` mostra), mas não escreve nelas.
const ColunaFinalDoAnchors = "READY TO TEST"

// ColumnsAnchorsWrites são as que os pipelines do Anchors movem. Serve ao doctor:
// uma coluna ausente aqui quebra o fluxo; uma ausente depois é problema do time.
func ColumnsAnchorsWrites() []string {
	for i, c := range ColunasDoBoard {
		if c == ColunaFinalDoAnchors {
			return ColunasDoBoard[:i+1]
		}
	}
	return ColunasDoBoard
}

// MissingWorkflow diz quais dos pipelines do fluxo não existem no projeto. Só presença —
// a coerência do conteúdo é outra pergunta, respondida por `SemConcurrency`.
func MissingWorkflow(root string) []Workflow {
	var faltam []Workflow
	for _, w := range WorkflowsDoFluxo {
		if _, err := os.Stat(filepath.Join(root, DirWorkflows, w.Arquivo)); err != nil {
			faltam = append(faltam, w)
		}
	}
	return faltam
}

// SemConcurrency diz quais pipelines existentes estão sem a serialização que o desenho
// pressupõe. Um `anchors-claim.yml` sem `concurrency` roda em paralelo consigo mesmo e
// atribui o mesmo card a dois agentes — o problema exato que o pipeline existe para
// eliminar, de volta e agora invisível, porque o arquivo ESTÁ lá.
//
// Verifica também `cancel-in-progress: false`: cancelar a execução em curso poderia
// matá-la depois de ela já ter criado metade dos cards.
func SemConcurrency(root string) []Workflow {
	var quebrados []Workflow
	for _, w := range WorkflowsDoFluxo {
		if !w.ExigeSerial {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, DirWorkflows, w.Arquivo))
		if err != nil {
			continue // ausente é outro achado (FaltaWorkflow); não contar duas vezes
		}
		texto := string(b)
		if !strings.Contains(texto, "concurrency:") || !strings.Contains(texto, "cancel-in-progress: false") {
			quebrados = append(quebrados, w)
		}
	}
	return quebrados
}

// MarcadorDeTemplate identifica um pipeline que ainda é o do Anchors — não editado pelo
// time. É o que separa "customizado" de "desatualizado".
//
// Sem essa distinção, `--fix` só sabia que o arquivo EXISTIA, e tratava as duas
// situações igual: preservava as duas. O efeito era que uma correção no template nunca
// alcançava quem já tinha instalado — o projeto ficava com o pipeline defeituoso para
// sempre, e nada avisava. Foi assim que um `stale` que reciclava trabalho pronto
// continuou rodando depois de corrigido na fonte.
// Sem o `#`: o marcador vive tanto num YAML (`# anchors:template`) quanto num HTML
// (`<!-- anchors:template -->`), e prendê-lo à sintaxe de comentário de uma linguagem
// faria a página do board nunca casar — o `--fix` a trataria como editada pelo time e
// jamais a atualizaria, em silêncio.
//
// The value lives in the scan (`scan.UpstreamMarker`), which reads the same marker to keep
// the map from giving a vendored pipeline a local code: one marker, one definition.
const MarcadorDeTemplate = scan.UpstreamMarker

// ÉTemplateIntacto diz se o pipeline instalado ainda é o do Anchors.
//
// A prova é o marcador que o template carrega. Editar o arquivo pede que se remova a
// linha — e quem edita sem removê-la perde a customização na próxima atualização. É uma
// troca deliberada: o marcador é explícito no arquivo e diz exatamente isso, e a
// alternativa (hash do conteúdo) marcaria como "customizado" qualquer arquivo que o
// próprio Anchors gerou com um branch de integração diferente.
func ÉTemplateIntacto(root, arquivo string) bool {
	b, err := os.ReadFile(filepath.Join(root, DirWorkflows, arquivo))
	if err != nil {
		return false
	}
	return strings.Contains(string(b), MarcadorDeTemplate)
}

// OutdatedWorkflows diz quais pipelines instalados são do Anchors (marcador
// intacto) e diferem do template atual — os que `--fix` pode e deve atualizar.
//
// Um pipeline SEM o marcador não entra aqui mesmo que difira: é do time, e a diferença é
// a customização dele.
func OutdatedWorkflows(root string, cfg *config.Config) []Workflow {
	var velhos []Workflow
	for _, w := range WorkflowsDoFluxo {
		if !ÉTemplateIntacto(root, w.Arquivo) {
			continue
		}
		atual, err := os.ReadFile(filepath.Join(root, DirWorkflows, w.Arquivo))
		if err != nil {
			continue
		}
		esperado, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			continue
		}
		esperado = applyIntegrationBranch(esperado, cfg.Workflow.IntegrationBranchOrDefault())
		if !bytes.Equal(atual, esperado) {
			velhos = append(velhos, w)
		}
	}
	return velhos
}

// SemeiaWorkflows escreve em `.github/workflows/` os pipelines que faltam, e ATUALIZA os
// que ainda são do Anchors e ficaram para trás. Devolve os nomes escritos.
//
// Nunca sobrescreve um pipeline que o time editou — outro ritmo de stale, uma permissão a
// mais, um passo de build próprio é trabalho deliberado, e reescrevê-lo apagaria a
// customização sem avisar. É a régua que o `install-hooks` já usa com um pre-commit
// alheio. A diferença é que agora "editado pelo time" é uma pergunta que se responde
// (o marcador), e não uma suposição a partir da mera existência do arquivo.
// Recebe o CONFIG, e não o branch já extraído: as regras de branch moram no anchors.yaml,
// e quem precisa delas as lê de lá. Passar o valor pronto espalharia a decisão por cada
// chamador, e bastaria um deles ler de outro lugar para o projeto ter dois fluxos.
func SemeiaWorkflows(root string, cfg *config.Config) ([]string, BoardOutcome, error) {
	dir := filepath.Join(root, DirWorkflows)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, BoardUnchanged, fmt.Errorf("creating %s: %w", DirWorkflows, err)
	}
	var escritos []string
	for _, w := range WorkflowsDoFluxo {
		dest := filepath.Join(dir, w.Arquivo)
		if _, err := os.Stat(dest); err == nil && !ÉTemplateIntacto(root, w.Arquivo) {
			continue // existe e o marcador saiu: é do time, não nosso
		}
		conteudo, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			return escritos, BoardUnchanged, fmt.Errorf("reading the template %s: %w", w.Arquivo, err)
		}
		conteudo = applyIntegrationBranch(conteudo, cfg.Workflow.IntegrationBranchOrDefault())
		if err := os.WriteFile(dest, conteudo, 0o644); err != nil {
			return escritos, BoardUnchanged, fmt.Errorf("writing %s: %w", dest, err)
		}
		escritos = append(escritos, w.Arquivo)
	}
	// A PÁGINA do board acompanha o pipeline que a publica: semear um sem o outro deixa o
	// fluxo pela metade — o workflow roda e falha ao copiar um arquivo que não existe.
	board, err := semeiaBoard(root)
	if err != nil {
		return escritos, board, err
	}
	sort.Strings(escritos)
	return escritos, board, nil
}

// BoardOutcome diz o que aconteceu com a PÁGINA do board numa semeadura.
type BoardOutcome int

const (
	// BoardUnchanged — a página já era idêntica ao template, ou é do time (marcador
	// removido) e o Anchors não a toca.
	BoardUnchanged BoardOutcome = iota
	// BoardCreated — não existia.
	BoardCreated
	// BoardUpdated — existia, era do Anchors, e ficou para trás.
	BoardUpdated
)

// semeiaBoard escreve a página do board, e a mantém atualizada pela mesma régua dos
// pipelines: intocada pelo marcador, o Anchors a atualiza; editada, ela passa a ser do
// time.
//
// Devolve O QUE FEZ, e não só o erro. A primeira versão devolvia apenas `error`, e o
// `doctor --fix` imprimia "os pipelines já existem e estão atualizados" enquanto
// REESCREVIA a página em silêncio — a mudança aparecia no `git status` sem nada tê-la
// anunciado, e quem visse o diff não saberia se tinha feito aquilo.
//
// Um arquivo que muda sozinho no repositório de alguém precisa ser dito.
func semeiaBoard(root string) (BoardOutcome, error) {
	dest := filepath.Join(root, BoardFile)
	conteudo, err := fs.ReadFile(boardFS, "board/anchors-board.html")
	if err != nil {
		return BoardUnchanged, fmt.Errorf("reading the board template: %w", err)
	}

	estado := BoardCreated
	if b, err := os.ReadFile(dest); err == nil {
		if !strings.Contains(string(b), MarcadorDeTemplate) {
			return BoardUnchanged, nil // existe e é do time
		}
		// Idêntico não é atualização: dizer "atualizei" quando nada mudou treina quem lê
		// a ignorar o aviso, e aí ele deixa de servir quando a mudança for real.
		if bytes.Equal(b, conteudo) {
			return BoardUnchanged, nil
		}
		estado = BoardUpdated
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return BoardUnchanged, err
	}
	if err := os.WriteFile(dest, conteudo, 0o644); err != nil {
		return BoardUnchanged, err
	}
	return estado, nil
}

// applyIntegrationBranch troca o branch cravado no template pelo que o projeto
// declarou. A linha alvo é marcada com `# anchors:integration-branch` — um marcador, e
// não uma busca por "main", porque "main" aparece em comentário e em outros contextos, e
// substituir a ocorrência errada quebraria o pipeline de um jeito difícil de ver.
func applyIntegrationBranch(conteudo []byte, branch string) []byte {
	if branch == "" || branch == "main" {
		return conteudo
	}
	linhas := strings.Split(string(conteudo), "\n")
	for i, l := range linhas {
		if strings.Contains(l, "# anchors:integration-branch") {
			indent := l[:len(l)-len(strings.TrimLeft(l, " "))]
			linhas[i] = indent + "branches: [" + branch + "] # anchors:integration-branch"
		}
	}
	return []byte(strings.Join(linhas, "\n"))
}

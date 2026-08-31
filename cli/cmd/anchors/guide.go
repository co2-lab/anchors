package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agentGuide é o playbook que ensina uma IA a OPERAR o Anchors. A inversão: o
// Anchors não invoca IA; a IA (no CLI dela) invoca o Anchors como ferramenta. Ela
// roda `anchors guide`, aprende o fluxo, e opera com os comandos existentes.
const agentGuide = `# Operando o Anchors (guia para agentes de IA)

Você é um agente de IA e o Anchors é uma FERRAMENTA de linha de comando que você
usa para desenvolver com rigor. Você lê e escreve os arquivos (specs, código,
features, testes) — o Anchors NÃO gera conteúdo; ele diz o que fazer, mantém o mapa
de dependências, enfileira o trabalho, e verifica o que você fez.

## O modelo

- Uma ÂNCORA é um documento que guia o desenvolvimento e confronta o que foi feito
  (spec, feature, guide, doc). Elas vivem no repositório.
- O MAPA (anchors.graph.yaml) liga as âncoras aos arquivos: quem depende de quem.
- Um GUIDE é a régua de como produzir cada artefato (spec, código, feature, teste,
  plano). SEMPRE leia o guide relevante antes de escrever.
- A FILA (.anchors/tasks/) é o trabalho pendente. O WATCHER a alimenta: cada vez que
  um arquivo muda, ele enfileira uma task com o próximo passo sugerido. Quem trabalha
  PUXA a task — o watcher nunca chama você.

## Dois papéis (leia com atenção — é o que te mantém livre)

Há DOIS modos de você agir, e não misture:

- CONVERSA: a sessão em que você fala com o usuário. Ela NUNCA deve ficar presa
  executando uma fila de trabalho. Aqui você planeja, dispara, e reporta.
- WORKER: quem realmente executa uma task da fila (escreve a spec, implementa,
  testa). Um worker pega UMA task, faz UM passo, fecha, e acaba.

### Como rodar um worker SEM prender a conversa

Quando há trabalho na fila e é hora de executá-lo:

  • SE você tem como delegar a um subagente/processo em BACKGROUND (ex.: uma tool de
    subagente): delegue o worker e VOLTE a conversar com o usuário. Ele fica livre.
  • SE você NÃO tem background (muitos clientes de IA não têm): NÃO sequestre a
    conversa em silêncio. Diga ao usuário, por exemplo:
      "Há N tasks na fila (anchors queue). Posso processá-las agora, mas isso vai me
       ocupar por um tempo — quer que eu vá em frente, ou prefere abrir outra sessão/
       terminal para rodar 'anchors next' em paralelo?"
    A decisão é do usuário. A fila é arquivo em disco e o claim é atômico, então
    várias sessões podem rodar 'anchors next' ao mesmo tempo sem colisão.

Regra de ouro: trabalho da fila roda em background se der; senão, só com o aval
explícito do usuário. A conversa é dele, não sua para monopolizar.

## O fluxo de desenvolvimento

### 0. PREPARE O MAPA E LIGUE O WATCHER (primeira coisa)
Rode:  anchors map build     (o watcher e o impact/check precisam do mapa existir)
Rode:  anchors watch start
A partir daqui, toda mudança de arquivo vira uma task na fila — inclusive as SUAS.
Sim, é redundante quando foi você quem mudou o arquivo (você já sabia). Mas é o MESMO
caminho de quando a mudança vem de fora (o usuário editou no editor, um 'git pull'
trouxe algo). Um mecanismo só, para toda origem. Confie na fila, não na sua memória.

### 0.5. DESCOBRIR — SÓ se o projeto ainda NÃO EXISTE (na CONVERSA)
Se o diretório está vazio (ou quase — sem código, sem anchors.yaml), há uma passada
ANTES do plano: descobrir o que o projeto é tecnicamente. Rode 'anchors guide project'
e siga a régua. Em resumo: VOCÊ entrevista o usuário em 5 etapas (propósito e forma →
linguagem → arquitetura e paradigma → estrutura macro e convenções → ferramental e
formatação), UMA pergunta por vez, esperando a resposta de cada uma. No fim, e só no
fim, escreve dois arquivos na raiz:
  • PROJECT.md   o resumo TÉCNICO decidido (stack, paradigma, estrutura, indentação,
                 extensões, editores) — é o que se lê antes de escrever cada arquivo
  • INSIGHTS.md  a transcrição: cada pergunta, cada resposta, e o que foi DESCARTADO
                 e por quê — é o que responde "por que esta escolha?" depois
Pule esta etapa num projeto que já tem código: lá o 'anchors init' infere do disco.
(Mesmo assim vale escrever o PROJECT.md do que você OBSERVOU — veja o guide.)

### 1. PLANEJAR (na CONVERSA)
Leia o guia de plano:  anchors guide plan
Com o usuário, rascunhe EM TEXTO na conversa um PLANO que decide QUAIS specs precisam
nascer ou mudar (semeia specs, nunca código direto). Itere até ele APROVAR o escopo.
Só ENTÃO escreva o arquivo do plano — salvá-lo já faz o watcher enfileirar a task
"especificar" (a esteira arranca). Escrever o arquivo antes do "sim" liga a máquina
sem aprovação. → aprovado? escreva o plano.

### 2+. TRABALHE PUXANDO DA FILA
Do plano em diante, o trabalho flui pela fila. O ciclo de um WORKER é sempre o mesmo:

  a) anchors next               puxa e reivindica a próxima task (atômico)
  b) anchors map build          incorpore ao mapa o arquivo que a task cita (se ele é
                                NOVO, ele ainda não está no mapa — sem este passo,
                                'impact' e 'check' respondem "não está no mapa")
  c) anchors impact <arquivo>   agora sim: o detalhe fino do que propagar/validar
  d) execute o passo, lendo o GUIDE certo antes de escrever:
       • task "specify"   → escreva os .spec.md (guide de spec). Dê a cada requisito
                            seu código de cenário (identidade) e seu regime.
       • task "implement" → escreva código + feature (guides de código/feature).
       • task "test"      → escreva os testes (guide de teste).
       • task "verify"    → só confronte (não há artefato novo a escrever).
       • task "verify-tests" → RODE a suíte de testes gerando os artefatos de
                            cobertura (veja o CoverageHint do seu stack: jest
                            --coverage + jest-junit, go test -coverprofile + go-junit-
                            report, pytest --cov --junitxml…), depois:
                              anchors ingest --junit <r.xml> --lcov <cov.info>
                            para o Anchors saber quais testes passaram e quais
                            requisitos da spec ficaram provados. Depois pergunte as
                            duas coisas que pegam bug NOVO:
                              anchors coverage --diff <base> --lcov <cov>  (o que você
                                mudou está coberto? — pega a linha nova sem teste)
                              anchors coverage --delta                     (a cobertura
                                caiu vs. antes? — pega regressão de cobertura)
                            SÓ ENTÃO confronte.
  e) anchors map build          incorpore os arquivos que VOCÊ acabou de escrever
  f) anchors check --changed <arquivo>   confronte (veja CONFRONTAR abaixo)
  g) anchors done <id>          feche a task (vai para .anchors/done/)
  h) volte para (a) até 'anchors next' dizer "fila vazia"

REGRA: 'impact' e 'check' leem o MAPA, não o disco. Todo arquivo novo ou movido só
existe para eles depois de um 'anchors map build'. Na dúvida, rode 'map build' antes.

Cada arquivo que você salva em (d)/(e) faz o watcher enfileirar a PRÓXIMA task
(spec→implement, feature→test, ...). Por isso o ciclo se sustenta: você não precisa
lembrar o que vem depois; a fila te diz.

### CONFRONTAR (o passo (f), em detalhe)
  anchors check --changed <arquivo>   incremental, o que a mudança tocou
  anchors check --all                 o projeto todo
  • Gate BLOQUEANTE reprovado impede a promoção (exit 1); informativo só registra.
  • REPORTE ao usuário os gates que falharam e por quê. Corrija e rode de novo.
  • NUNCA feche a task (done) com gate bloqueante ainda vermelho.

### JULGAR (o gate que um script não computa — VOCÊ é o medidor)
Alguns gates medem o que nenhum script sabe: "esta tela quebra em atomic design?",
"a spec descreve comportamento e não implementação?", "todo export público desta
unidade rastreia a uma regra da spec — ou é código morto/scope-creep?" (o inverso do
feature-test-match: aqui um SÍMBOLO grepável não basta, porque a spec descreve
comportamento e não cita nomes de impl — julgar exige entender que a regra X justifica
o export Y). São gates 'measures: judgment'.
O 'anchors check' não os computa — marca os alvos como ⏳ e os enfileira. Você:
  a) anchors judge --pending            veja os alvos e, em cada task, o guide + a pergunta
  b) LEIA o guide indicado. Vá direto à seção "## Pontos de conformidade" — a lista de
     itens CK é o que você deve verificar. NÃO julgue pela prosa inteira "no olho";
     julgue o alvo contra CADA ponto. Se a checklist é AGRUPADA POR TARGET (subseções
     "### Para <camada>"), aplique só os pontos do grupo da camada DESTE alvo + os do
     grupo "Para todos" — um ponto de tela não vale para um model. Isso torna o
     veredito focado, objetivo e reproduzível. (Se o guide não tem a seção, é débito —
     o gate determinístico 'guide-has-checklist' já avisa; julgue pelo que houver.)
  c) anchors judge <alvo> --gate <g> --verdict pass|fail --reason "<LAUDO>"
Você avaliou o alvo contra cada CK; NÃO devolva uma frase — devolva o LAUDO por item.
Para cada CK: conforme, ou a não-conformidade com o quê, ONDE (arquivo:linha), qual CK
violou, e COMO corrigir. Esse texto vira o corpo da issue — assim ninguém reprocessa o
alvo depois só para saber o que arrumar.
'fail' abre a issue com o laudo; 'pass' resolve a issue anterior (se houver). O
veredito é carimbado com a rev do alvo, então ENVELHECE (stale) se o alvo mudar.

### SAÚDE (quando quiser o panorama)
  anchors doctor    pontas SISTÊMICAS que o check não vê: órfãos, buracos de
                    cobertura, camadas frouxas, arestas mortas. Não bloqueia.

## ADOÇÃO — trazer o Anchors para um projeto que JÁ EXISTE

Um projeto que NASCE com o Anchors não tem este problema: cada arquivo é confrontado
e carimbado quando é criado, então nunca acumula dívida. O cold start pesado é só de
QUEM ADOTA — um projeto grande e antigo cujo mapa nasce com milhares de pares
(guide, arquivo) todos sem carimbo. Se for o seu caso, NÃO tente auditar tudo de uma
vez (pode custar milhões de tokens). Faça em LOTES, nesta ordem:

1. CORTE ESTRUTURAL primeiro (grátis, sem IA). Rode 'anchors governs' e 'anchors
   doctor'. Se vários guides regem o MESMO conjunto (redundância), afine as tags no
   anchors.yaml para cada guide reger só o seu escopo. Isso corta o trabalho na raiz
   antes de gastar um token. Complete também as regras 'governs' que faltam (um guide
   sem governança não é guide — o doctor avisa).

2. BATCH POR GUIDE, não por alvo. O guide é caro de ler (centenas de linhas) e o
   alvo é barato. Então NÃO releia o guide a cada arquivo: leia o guide UMA vez e
   julgue muitos alvos daquele guide em sequência. Use 'anchors governs <guide>' para
   pegar a lista de alvos, leia o guide, e julgue o lote. Isso amortiza o custo fixo
   do guide entre os alvos (corte de ~4×). Se você delega a subagentes para
   paralelizar, dê UM GUIDE INTEIRO a cada subagente (nunca um alvo por subagente —
   isso faz cada um reler o guide, o oposto da economia). Julgamento é leitura: não
   precisa de worktree/workspace isolado.

3. O CARIMBO É O CURSOR. A primeira passada NÃO precisa terminar de uma vez. Cada
   'anchors judge' carimba o alvo; reconfrontar pula o que já tem carimbo válido na
   rev atual. Então faça um lote hoje, outro amanhã — o Anchors só re-enfileira o que
   ainda não foi julgado (ou o que mudou desde então). Reporte ao usuário o que cobriu
   e o que falta; nunca finja cobertura total que não houve.

## CORREÇÃO EM LOTE — varrer o projeto arquivo a arquivo (paralelo, sem retrabalho)

Quando há MUITAS pendências espalhadas (ex.: aplicar o cabeçalho @anchors a todo o
projeto, completar specs), a forma eficiente é varrer por ARQUIVO, com vários
workers em paralelo — mas há duas regras que evitam desperdício:

1. UM ARQUIVO POR VEZ, TODAS AS PENDÊNCIAS. Se você vai abrir um arquivo, conserte
   TUDO que pende nele de uma vez — não só o header. Use:
     anchors audit <arquivo>            as pendências DELE (gates + doctor)
     anchors audit <arquivo> --impact   a UNIDADE toda (a trinca spec↔code↔feature↔
                                        test) — conserta a unidade num worker só.
   É desperdício um worker mexer no header do .spec e outro no .tsx irmão; pegue a
   unidade (--impact) e resolva a trinca inteira de uma vez.

2. DE CIMA PRA BAIXO NA ÁRVORE. Consertar um FILHO e depois o PAI (que o rege)
   obriga a refazer o filho. Processe na ordem topológica — réguas/specs (pais)
   ANTES de código/testes (regidos). O Anchors dá a ordem pronta:
     anchors map show --worklist --pending   os arquivos com pendência, JÁ ordenados
                                             de cima pra baixo. Processe nessa ordem.
   Ao paralelizar: distribua FATIAS da worklist mantendo a ordem (um worker por
   feature/subárvore), nunca um pai e seu filho em workers concorrentes.

## O que reportar ao usuário

- Ao ligar o fluxo: "watcher ligado; vou planejar com você e depois processar a fila
  em background / ou pedir seu aval para rodar, conforme minha capacidade."
- Após 'check': quais gates falharam, quais bloqueiam, o débito que ficou.
- Após 'doctor': pontas sistêmicas.
- Divergências que você NÃO resolve sozinho → apresente as opções: corrigir o alvo,
  atualizar a âncora (spec/guide), ou virar um novo plano.
- NUNCA "conserte" fazendo a âncora mentir sobre o código. Se o código diverge da
  spec, ou o código está errado (corrija) ou a spec envelheceu (atualize a spec, o
  que pode gerar novo trabalho — reporte isso).

## Referência de comandos

  anchors init                      configura o projeto (anchors.yaml) — uma vez
  anchors install-hooks             instala o git pre-commit que roda os gates nos staged
  anchors new <kind> <nome>         emite o esqueleto de um artefato (spec|feature|test)
  anchors new <kind> --list-sections  as seções (default/opcional) do kind
  anchors recode <antigo> <novo>    renomeia um código e propaga (dry-run; --apply grava)
  anchors watch start               liga o watcher (enfileira tasks em background)
  anchors watch status|stop|pause|resume|logs   controla o watcher
  anchors queue                     lista as tasks vivas (leitura)
  anchors next                      puxa+reivindica a próxima task (o worker)
  anchors done <id>                 fecha uma task concluída
  anchors map build                 (re)constrói o mapa de dependências
  anchors map show <arquivo>        vizinhança de um nó (↑regido / ↓propaga)
  anchors impact <arquivo>          o que uma alteração atinge (↓propaga / ↑valida)
  anchors check --changed <arq>     roda os gates sobre o caminho de impacto
  anchors check --all               roda os gates sobre tudo
  anchors ingest --junit/--lcov     ingere sinais de teste do runner (execução/cobertura)
  anchors coverage [<spec>]         cobertura por cenário (requisito provado?) e por linha
  anchors coverage --diff <base> --lcov <c>   o que MUDEI está coberto? (patch coverage)
  anchors coverage --delta          a cobertura CAIU desde a última ingestão?
  anchors ingest --junit r --layer <unit|integration|e2e>   ingere por CAMADA (mergeia)
  anchors report tests|quality|structure|config|issues|inconsistencies   relatórios em docs/
  anchors report all                gera todas as perspectivas + índice em docs/anchors/
  anchors doctor                    saúde do ecossistema (pontas sistêmicas)
  anchors guide                     este guia
  anchors guide project             a régua da fase DESCOBRIR (projeto novo → PROJECT.md)
  anchors guide plan                a régua da fase PLANEJAR

## Onde estão os guides do projeto

Os guides (as réguas de cada artefato) vivem no repositório — veja a seção do
anchors.yaml para a camada 'guide' (tipicamente em guides/). Leia o guide certo
ANTES de escrever cada tipo de artefato. Se o projeto não tem um guide para algo
que você vai produzir, avise o usuário — é uma ponta aberta.
`

func newGuideCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guide",
		Short: "Imprime os guias do Anchors para agentes de IA",
		Long: `Sem argumentos, imprime o playbook de operação: o fluxo de desenvolvimento
(planejar → especificar → mapear → implementar → testar → confrontar), os comandos,
e o que reportar ao usuário. A IA roda isto primeiro e opera o Anchors como ferramenta.

Subcomandos imprimem os guias das réguas específicas:
  anchors guide project  como descobrir um projeto que ainda não existe (PROJECT.md)
  anchors guide plan     como estruturar um plano (a origem do movimento)
  anchors guide spec     como escrever uma spec (a origem da verdade)
  anchors guide code     como implementar o código guiado pela spec
  anchors guide feature  como escrever a feature (os cenários de comportamento)
  anchors guide test     como escrever os testes (a régua executável)
  anchors guide guide    como escrever um guide (a régua de uma régua)
  anchors guide header   o bloco de cabeçalho de cada arquivo (transversal, mandatório)`,
		// sem subcomando → o playbook de operação
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(agentGuide)
			return nil
		},
	}
	cmd.AddCommand(
		newGuideSubCmd("work", "como trabalhar um card: a ordem, o que fazer com o achado que não é dele", workGuide),
		newGuideSubCmd("project", "como descobrir um projeto que ainda não existe (PROJECT.md + INSIGHTS.md)", projectGuide),
		newGuidePlanCmd(),
		newGuideSubCmd("spec", "como escrever uma spec (a origem da verdade)", specGuide),
		newGuideSubCmd("code", "como implementar o código guiado pela spec", codeGuide),
		newGuideSubCmd("feature", "como escrever a feature (cenários de comportamento)", featureGuide),
		newGuideSubCmd("test", "como escrever os testes (a régua executável)", testGuide),
		newGuideSubCmd("guide", "como escrever um guide (a régua de uma régua)", guideGuide),
		newGuideSubCmd("header", "o bloco de cabeçalho de cada arquivo (transversal, mandatório)", headerGuide),
	)
	return cmd
}

// newGuideSubCmd fabrica um subcomando de guia que só imprime um texto embutido.
// As quatro réguas (spec/code/feature/test) compartilham essa casca fina.
func newGuideSubCmd(use, short, body string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(body)
			return nil
		},
	}
}

func newGuidePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Imprime o guia de plano (como um plano se estrutura e semeia specs)",
		Long: `O guia de plano é a régua da fase PLANEJAR: como a IA e o usuário
escrevem um plano que decide QUAIS specs nascem ou mudam — nunca código direto.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(planGuide)
			return nil
		},
	}
}

package main

// workGuide é a régua de quem PEGA UM CARD. Ele responde ao que o card não cabe dizer
// sem se repetir em cada issue — e o custo da repetição não é só ruído: instrução copiada
// CONGELA. Um card criado hoje carrega o texto de hoje, e quando a instrução muda, os
// cards antigos passam a ensinar o errado sem que nada acuse.
//
// Centralizado no binário, o guia acompanha a versão: quem roda `anchors guide work` lê a
// instrução ATUAL, não a do dia em que o card nasceu.
const workGuide = `# Guia de trabalho (a régua de quem pegou um card)

## Antes de começar

` + "`anchors status`" + ` diz onde o projeto está, e ` + "`anchors guide <artefato>`" + ` diz
como escrever o que você vai escrever (spec, code, feature, test).

## A ordem: da DIREITA para a ESQUERDA

Termine o que está mais adiantado antes de pegar coisa nova. Trabalho pela metade não
entrega nada e envelhece até o contexto se perder — quem o retoma paga de novo o custo de
entender, e às vezes descobre que a decisão mudou no meio.

## Achou algo errado que NÃO é este card?

Acontece o tempo todo: uma config que contradiz a doutrina do projeto, um caminho que
ninguém documentou, um arquivo no lugar errado. **Nenhum gate vê isso** — o gate abre
issue do que ELE detecta, e o resto depende de você.

O caminho barato é consertar na hora e seguir. E aí o conserto some do histórico: quem
vier depois não sabe que aquilo já foi problema, nem por que a solução é aquela.

    anchors escalate "<o que está errado>" --about <arquivo> --card <este card>

O achado nasce com a label ` + "`anchors:sob-<número>`" + `, e os dois se entregam no MESMO
PR: você já está com o contexto na mão, e separá-los faria um dos dois esperar sem razão.

### Quando NÃO usar

Se a correção é trivial E está no arquivo que você já está editando, corrija e registre a
revisão no próprio arquivo (` + "`{CODIGO}-R0001: o que mudou e por quê`" + `). Abrir card
para trocar uma palavra é burocracia.

Se a mudança **impacta a direção do projeto** — ou se você tem dúvida —, não a faça:

    anchors escalate "<o que precisa mudar>" --about <arquivo> --for-user

Isso vira decisão de quem planejou, e o card para até ela sair. A interpretação do impacto
é sua: você é quem tem o contexto do que descobriu.

## Antes de commitar

O ` + "`anchors check`" + ` barra enquanto houver julgamento pendente — ele é trabalho DESTE
commit, e quem mexeu no arquivo é quem tem contexto para responder:

    anchors judge --pending
    anchors judge <alvo> --gate <g> --verdict pass|fail --reason "..."

## Ao abrir o PR: o Anchors escreve as linhas de fechamento

    anchors pr-body

Ele imprime as linhas que fecham o card que você pegou E os achados que nasceram sob
ele (` + "`anchors:sob-<n>`" + `). Cole no corpo do PR.

Você NÃO precisa saber a sintaxe da plataforma — e é justamente ela que se erra em
silêncio. O GitHub reconhece a palavra de fechamento só em INGLÊS: escrever "Fecha
#44" num projeto em português é ignorado sem erro nenhum, o PR mescla, e o card fica
aberto. Aconteceu: um PR dizia "Fecha #44, #49, #50" e os três continuaram abertos.

O vínculo é declarado no vocabulário do Anchors — o card que você pegou e os achados
sob ele. A palavra que a plataforma entende é DERIVADA disso, e muda com a
plataforma, não com o idioma do projeto.

## Depois de abrir o PR: o card não está entregue

Abrir o PR não fecha o card, e **empurrar um commit não é um ponto de parada**. O CI roda
depois do push, e o resultado dele é parte do seu trabalho — não notícia que alguém traz.

    anchors work review          (mostra o que falta no PR que é seu)
    gh pr checks <n> --watch     (BLOQUEIA até o CI concluir, e devolve o veredito)

O ` + "`--watch`" + ` existe para isso: ele espera pelo processo, você não. Sem ele, a única forma
de saber o resultado é perguntar de novo mais tarde — e "aguardando a nova rodada" é uma
frase que encerra o turno sem entregar nada. O card continua ` + "`in-progress`" + `, com seu nome
nele, e quem o retoma paga de novo o custo de entender.

### O CI reprovou

O vermelho é trabalho DESTE card, não um card novo. Leia a falha, conserte, empurre e
**espere de novo** — quantas rodadas forem necessárias. Só há duas saídas legítimas do
ciclo:

- o CI ficou verde e o card avançou no board; ou
- a falha exige uma decisão que não é sua, e aí ela vira escalonamento
  (` + "`anchors escalate ... --for-user`" + `), com o card parado explicitamente.

Uma terceira coisa NÃO é saída: relatar o diagnóstico e parar. O diagnóstico correto é
metade do trabalho; a outra metade é o veredito do check que você disparou.

## A spec nasce antes do código

É o fluxo normal: a spec é a âncora. Enquanto as peças não existem, declare o que falta:

    > **@TBD: code,feature,test** — as peças são a fase em andamento.

` + "`@TBD`" + ` é *to be developed*, e vence sozinho: quando a peça aparece no mapa, o gate
volta a confrontá-la. É diferente de ` + "`@no-test`" + `, que afirma "esta unidade NÃO
PRECISA de teste" — permanente, e que apagaria a cobrança para sempre.
`

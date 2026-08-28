package main

// planGuide é a régua da fase PLANEJAR (Planejamento — a origem do movimento). Um
// plano no Anchors NÃO cria código; ele decide QUAIS specs precisam nascer ou mudar,
// e por quê. É a âncora que semeia as outras âncoras. Este guia embutido casa com a
// versão do binário — a IA o lê antes de escrever qualquer plano.
const planGuide = `# Guia de plano (a régua da fase PLANEJAR)

Um plano é a ORIGEM do movimento no Anchors. Ele não escreve código nem features —
ele decide QUAIS specs precisam nascer ou mudar, e por quê. O plano semeia as specs;
as specs semeiam o resto. Se você está prestes a escrever código a partir de um
plano, o plano falhou: ele deveria ter apontado uma spec.

## REGRA DE OURO: aprove o escopo ANTES de escrever o arquivo

Escrever o arquivo do plano é um ATO IRREVERSÍVEL de escopo: assim que ele é salvo na
camada de plano, o watcher o detecta e ENFILEIRA a primeira task ("especificar") — a
máquina começa a executar. Por isso a aprovação do usuário vem ANTES do arquivo existir,
não depois:

1. Na CONVERSA, rascunhe o escopo (objetivo, motivo, a LISTA de specs, fases, fora de
   escopo, definição de pronto) e apresente ao usuário EM TEXTO na resposta.
2. ITEREM até ele aprovar o escopo explicitamente.
3. SÓ ENTÃO escreva o arquivo do plano em plans/.

Não escreva o arquivo para "mostrar como ficou" e depois perguntar — isso já liga a
esteira. Se o watcher está ligado, salvar = começar. Trate o arquivo do plano como o
commit do escopo: só existe depois do "sim".

## O que um plano É e o que NÃO é

- É: uma decisão de escopo. "Estas N specs vão nascer/mudar, nesta ordem, por estes
  motivos, e é assim que saberemos que terminou."
- NÃO é: um design detalhado (isso é a spec), nem uma lista de tarefas de código.
- Um plano pode ser pequeno (uma spec) ou grande (uma fatia de produto). O tamanho
  do plano é o tamanho da mudança de escopo — não do esforço de implementação.

## A estrutura de um plano

Escreva o plano como um documento na camada de plano do projeto (veja o anchors.yaml
— tipicamente plans/). Um plano tem:

### 1. Objetivo
Uma ou duas frases: o que muda no produto e para quem. Sem jargão de implementação.

### 2. Motivo
Por que agora. O que está errado, faltando, ou pedido. É o que sobrevive quando
alguém reabre o plano em seis meses perguntando "por que fizemos isso?".

### 3. Specs semeadas
O coração do plano. Uma LISTA das specs que vão nascer ou mudar. Para cada uma:
  • o arquivo da spec (novo ou existente) — ex: features/larder/AddItem.spec.md
  • nasce ou muda? se muda, o que muda nela
  • uma linha do que ela deve cobrir (o detalhe fica na spec, não aqui)
Não liste código, telas, endpoints. Se você sente falta de listar isso, é sinal de
que a spec correspondente ainda não foi pensada — adicione a spec à lista.

### 4. Fases (ordem, dependências e progresso)
Agrupe as specs em FASES pequenas e independentemente entregáveis, na ordem em que
devem nascer. Se uma spec depende de outra antes, ela vem numa fase posterior. O
plano é uma sequência, não um balde. Marque o progresso com caixas — a IA (o worker)
vira '[ ]' em '[x]' ao concluir cada spec, então o plano é também o placar do
trabalho:

  ## Fase 1 — <nome>
  - [ ] features/larder/AddItem.spec.md — cadastro de item
  - [ ] features/larder/ItemCard.spec.md — cartão do item
  ## Fase 2 — <nome> (depende da Fase 1)
  - [ ] features/larder/ExpiryAlert.spec.md — alerta de validade

### 5. Fora de escopo
O que este plano deliberadamente NÃO faz. Fecha a porta para escopo que vaza. Cada
item aqui é uma tentação nomeada e recusada.

### 6. Definição de pronto
Como saber que o plano terminou. Normalmente: "todas as specs semeadas existem,
implementadas, com feature e teste, e 'anchors check' passa nos gates bloqueantes".
Se há um débito aceito (opt-out honesto), nomeie-o aqui — não o deixe implícito.

## Regras do plano

- SEMPRE por spec. Cada linha de trabalho no plano termina numa spec, nunca num
  arquivo de código. Essa é a disciplina spec-first: o plano semeia specs, a spec é
  a origem da verdade.
- O QUE, não COMO. O plano diz o que precisa existir e por quê. O COMO (comportamento,
  contratos, invariantes) é da spec; a implementação é do código. Não descreva
  algoritmo nem estrutura de dados no plano — se sentiu vontade, é conteúdo de spec.
- RECONCILIE antes de semear (anti-retrabalho). Antes de listar uma spec nova, veja o
  que JÁ existe — outros planos e specs. Não planeje o que outro plano já entrega, não
  crie escopo sobreposto nem ordenação impossível. Se você muda uma decisão porque ela
  conflita com algo existente, anote no plano o PORQUÊ — a razão precisa sobreviver.
- RESPEITE a estrutura. As specs semeadas caem nas camadas que o anchors.yaml
  declara. Se uma spec não tem camada onde morar, ou a estrutura está incompleta
  (reporte ao usuário) ou a spec está no lugar errado.
- IDENTIDADE cedo. Se o projeto usa códigos de cenário (identidade — TRACEABILITY),
  a spec vai carregá-los; o plano não precisa inventá-los, mas deve saber que cada
  requisito na spec vai receber um.
- AMBIGUIDADE vira TODO, nunca chute. Se o escopo de uma spec está incerto, escreva
  um TODO explícito no plano em vez de adivinhar. Um TODO honesto é melhor que uma
  decisão inventada que o worker vai implementar errado.
- MAS o TODO não pode governar a EXISTÊNCIA de um item. Esta é a distinção que decide
  se o plano é executável:

    ✓ delegar ONDE se decide — "os índices ficam na spec do SCHEMA, decididos ao ver
      o schema atual". O item existe; o CONTEÚDO de uma decisão dele mora noutro
      artefato, que vai nascer. Quem executa sabe onde a resposta vai estar.

    ✗ adiar SE o item existe — "nasce SE o mobile chamar via API Gateway; senão
      marcar fora de escopo". Aqui quem executa precisa RESOLVER a condição antes de
      saber o que fazer.

  O segundo tem duas saídas, e as duas são ruins: o agente decide sozinho (e a decisão
  de produto vira escolha de implementação, sem passar por quem devia), ou ele para e
  pergunta (e o fluxo automatizado morre ali). Medido num projeto real: o agente
  decidiu — acertou, e ninguém tinha pedido que ele decidisse.

  Antes de semear um item, pergunte: "quem executa isto precisa DECIDIR alguma coisa
  para começar?". Se sim, decida agora, ou semeie primeiro o item que RESPONDE a
  pergunta — a decisão vira trabalho, não pré-requisito invisível.
- PERGUNTA ABERTA fica FORA do checklist. O plano PODE (e deve) registrar o que ainda
  não se sabe — em prosa, numa seção própria, com quem decide. O que ele não pode é
  esconder a pergunta DENTRO de um item semeado, onde ela parece trabalho e é decisão.
- "A OU B" É UM CHUTE ADIADO — decida no plano, não na spec. Um plano que diz "usa
  repository OU service", "grava aqui OU ali", "síncrono OU async" NÃO é um TODO
  honesto: é uma bifurcação de ARQUITETURA que a spec vai resolver no silêncio, sem
  ninguém pesar a consequência (ex.: "repository" pode significar duplicar uma regra
  fora do lugar dela). Feche cada "ou" AGORA — com o usuário — ou marque-o como um TODO
  NOMEADO que bloqueia descer para a spec até ser resolvido. Nunca deixe a escolha para
  quem especifica.
- HONESTIDADE sobre débito. Se o plano decide pular um gate (ex: sem teste agora),
  isso é uma decisão explícita na Definição de pronto — nunca um silêncio.

## Depois de escrever o plano (já aprovado)

O arquivo só nasce depois do "sim" do usuário (ver a REGRA DE OURO acima). Ao salvá-lo,
o watcher enfileira a primeira task ("especificar"). Daí em diante o trabalho flui pela
fila: você (ou um worker em background) roda 'anchors next' e segue o ciclo do
playbook. O plano é o contrato de escopo; as specs são a execução dele. Se o escopo
mudar no meio, isso é uma EDIÇÃO consciente do plano (e nova conversa), não uma deriva
silenciosa das specs.

## Anti-padrões (recuse-os)

- Escrever o arquivo do plano ANTES da aprovação → liga a esteira sem "sim"; apresente
  o escopo em texto na conversa primeiro. O arquivo é o commit do escopo.
- Plano que lista arquivos de código → falta a spec; suba um nível.
- Plano sem "fora de escopo" → vai vazar; nomeie o que não faz.
- Plano sem "definição de pronto" → não há como saber quando parar; adicione.
- Plano que muda a régua (um guide) de passagem → mudar régua é um plano próprio,
  não um efeito colateral. Separe.
`

---
title: O fluxo de trabalho
description: Como um dia de trabalho acontece com o Anchors — do pedido de trabalho ao merge.
---

Os pilares dizem **o que** o Anchors defende. Esta página diz **o que você faz**,
na ordem, num dia normal de trabalho.

Ela vale para quem trabalha sozinho e para um time. A diferença entre os dois
casos aparece onde importa, e está marcada.

## O board é o estado

Todo trabalho vive num **card** — uma issue com a label `anchors`. O estado dele
é outra label, e a sequência é sempre a mesma:

```
to-do → in-progress → ready-to-review → in-review → ready-to-test → ...
└──────── alçada do agente ────────────┘ └── alçada dos pipelines de entrega ──┘
```

Duas coisas importam nessa divisão:

- **até `ready-to-review`, quem move é quem trabalha.** O card é seu enquanto
  você o tem.
- **de `ready-to-test` em diante, quem move são os pipelines de entrega** do
  projeto. O Anchors não vai além disso — ele não sabe nada sobre o seu
  processo de release.

Além do estado, duas labels atravessam qualquer coluna:

| label | o que significa |
| --- | --- |
| `anchors:precisa-do-usuario` | o card espera uma **decisão humana** e não é entregue a ninguém até ela sair |
| `anchors:sob-<n>` | este card é um **achado** que nasceu enquanto alguém trabalhava no card `<n>` |

## 1. Peça trabalho — não escolha

```sh
gh workflow run anchors-claim.yml -f agent=<maquina>/<sessao>
```

Você **pede**, o pipeline **decide**. Isso não é cerimônia: é o que impede dois
agentes de pegarem o mesmo card.

O motivo é técnico e não tem volta por outro caminho. Se cada um reivindicasse
direto, dois poderiam ler "sem dono" antes de qualquer um escrever, e ambos se
achariam donos — a API do GitHub não oferece compare-and-swap. O pipeline de
claim é serializado, então "há card livre?" e "atribua a ele" acontecem sem
ninguém no meio.

**A prioridade é da direita para a esquerda.** Termine o que está mais adiantado
antes de pegar coisa nova: trabalho pela metade não entrega nada, ocupa revisor,
e envelhece até o contexto de quem o escreveu se perder.

E o **seu** card vem antes de tudo. Se você já tem um em andamento, o claim o
devolve em vez de dar outro.

## 2. Leia a régua antes de escrever

```sh
anchors status          # onde o projeto está, e o próximo passo
anchors guide work      # a régua de quem pegou um card
anchors guide spec      # (ou code/feature/test) como escrever o artefato
```

O `guide work` responde três coisas que você vai precisar: a ordem do board, o
que fazer com um achado que **não** é do seu card, e o que o `anchors check`
cobra antes do commit.

Não decore — o guia está no binário justamente para não precisar.

## 3. Trabalhe, e confronte antes de commitar

```sh
anchors check --changed <arquivo>
```

O `check` roda os gates sobre o que você mexeu. Ele tem três desfechos, e
distingui-los evita muita frustração:

| símbolo | significa | barra? |
| --- | --- | --- |
| `✓` | passou | — |
| `✗ [BLOQUEIA]` | reprovou um gate bloqueante | **sim** |
| `✗ [informativo]` | achou algo que não impede a entrega | não |
| `~` | **indeterminado** — o gate não teve o que confrontar | não |

O `~` é o mais mal-entendido. Ele não é falha: é o gate dizendo que o artefato do
outro lado da relação não existe. Um gate que mede cobertura de teste num
arquivo sem teste não reprova — ele não mede.

Se um gate bloqueante reprovar e a reprovação for **deliberada**, declare na
mensagem do commit:

```
[skip-<regra>@<CODIGO>: por que a reprovação é aceitável aqui]
```

Isso não desliga o gate: registra a dispensa com a razão, no lugar onde quem
vier depois vai lê-la.

## 4. Achou algo que não é deste card?

Acontece o tempo todo. Uma config que contradiz a doutrina, um caminho que
ninguém documentou, um arquivo no lugar errado.

**Nenhum gate vê isso** — o gate abre issue do que ELE detecta; o resto depende
de você.

```sh
anchors escalate "<o que está errado>" --sobre <arquivo> --card <este card>
```

O achado nasce com a label `anchors:sob-<n>`, e os dois se entregam no **mesmo
PR**: você já está com o contexto na mão.

### Quando NÃO escalar

Se a correção é trivial **e** está no arquivo que você já está editando, corrija
e registre a revisão no próprio arquivo (`{CODIGO}-R0001: o que mudou e por
quê`). Abrir card para trocar uma palavra é burocracia.

### Quando a decisão não é sua

Se a mudança **impacta a direção do projeto** — ou se você tem dúvida:

```sh
anchors escalate "<o que precisa mudar>" --sobre <arquivo> --para-usuario
```

Isso vira decisão de quem planejou, e **o card para até ela sair**. A
interpretação do impacto é sua: você é quem tem o contexto do que descobriu.

## 5. Abra o PR — sem inventar a sintaxe

```sh
anchors pr-body
```

Ele imprime as linhas que fecham o card que você pegou **e** os achados que
nasceram sob ele. Cole no corpo do PR.

Você não precisa saber a sintaxe da plataforma — e é justamente ela que se erra
em silêncio. O GitHub só reconhece a palavra de fechamento em **inglês**:
escrever "Fecha #44" num projeto em português é ignorado sem erro nenhum, o PR
mescla, e o card fica aberto.

Isso aconteceu de verdade: um PR dizia "Fecha #44, #49, #50" e os três
continuaram abertos.

## 5.1 Espere o veredito — empurrar não é entregar

Abrir o PR não fecha o card, e **empurrar um commit não é um ponto de parada**. O
CI roda depois do push, e o resultado dele é parte do seu trabalho.

```sh
gh pr checks <n> --watch
```

O `--watch` **bloqueia** até o CI concluir: ele espera pelo processo, você não.
Sem ele, a única forma de saber o resultado é perguntar de novo mais tarde — e
"aguardando a nova rodada" é uma frase que encerra o turno sem entregar nada. O
card continua `in-progress`, com seu nome nele.

Se o CI reprovar, o vermelho é trabalho **deste** card, não um card novo: leia a
falha, conserte, empurre e espere de novo. Há exatamente duas saídas do ciclo:

- o CI ficou verde e o card avançou no board; ou
- a falha exige uma decisão que não é sua, e ela vira escalonamento
  (`anchors escalate ... --for-user`), com o card parado explicitamente.

Relatar o diagnóstico e parar **não é** uma terceira saída. O diagnóstico correto
é metade do trabalho; a outra metade é o veredito do check que você disparou.

## 5.2 O relato da rodada tem formato

```sh
anchors task-status
```

Quem lê o seu relato decide se continua, se revisa, ou se responde uma pergunta —
e o que decide isso não é a narrativa do que você fez. É o **estado**: onde o card
está, se o veredito do CI foi lido, e o que espera uma pessoa.

O comando descobre o que a máquina sabe (o card e seu estado, o PR e os checks, o
que não foi enviado, as decisões paradas em `needs-user`) e deixa **duas lacunas**,
que são suas:

- **O que provei** — as regras que a suíte confronta, e o que a mutação matou.
  Testes que passam não são prova; prova é a mutação que morreu.
- **O que ficou de fora** — nada, ou o que você deixou e por quê. Reduzir escopo é
  decisão de quem pediu: se algo não entrou, é aqui que ele descobre.

Sem formato, cada rodada relata o que o agente achou importante, e o que se omite
primeiro é justamente o estado.

## 6. A revisão

O claim entrega cards em `ready-to-review` **antes** dos `to-do` — revisar vem
antes de começar coisa nova.

```sh
anchors guide review
```

### O passo zero: os checks EXISTEM?

Antes de qualquer julgamento:

```sh
gh pr checks <n>
```

Um check que não rodou **não aparece como falha — ele não aparece**. O PR fica
com a mesma cara de um PR aprovado, e a plataforma não distingue "passou" de
"nunca existiu".

Aconteceu: um PR com conflito de merge não dispara workflow de `pull_request` no
GitHub, em silêncio. Nenhum gate rodou, e nada acusou.

Se os checks do Anchors não estão na lista, dispare antes de revisar:

```sh
gh workflow run anchors-gates.yml --ref <branch-do-PR>
```

**Não aprove um PR cujos checks nunca rodaram.** Você é a fronteira que sobra
quando a automação falha em silêncio.

### O que é seu, e o que não é

O que o script já confrontou **não** é trabalho de review: o mapa em dia, os
gates bloqueantes, os cards declarados no corpo. Reconferir à mão gasta você no
que a máquina faz melhor.

O que é seu é o que exige julgamento:

- **A spec decide o que precisava decidir?** Uma spec que diz "o sistema deve ser
  rápido" passa em todos os gates e não decide nada.
- **O código realiza a regra, ou só a cita?** Um veredito `pass` sobre marcação
  num lugar genérico (topo do arquivo, import) é o defeito que a automação erra
  com mais frequência.
- **O teste PROVA, ou só executa?** Leia as asserções, não o número.
- **O que o PR mudou sem dizer?** Diff que o autor não menciona na descrição é
  onde mora o que ele não percebeu que mudou.

### Os três desfechos

| o que você achou | o que fazer |
| --- | --- |
| nada | mova o card para `ready-to-test` |
| defeito de **execução** (marcação errada, teste que não cobre o caso) | **corrija você mesmo**, no mesmo PR, e devolva o card para `ready-to-review` liberando a posse |
| defeito de **entendimento** (a spec foi lida errado, a abordagem não serve) | **devolva**: card para `to-do`, posse para quem implementou |

A diferença entre os dois últimos não é de tamanho, é de natureza. Corrigir um
defeito de entendimento esconderia que o autor entendeu errado, e o mesmo erro
volta no próximo card que ele pegar.

> **Quem corrige não aprova o próprio conserto.** Ao corrigir, você virou autor
> daquele trecho.
>
> **Time de um:** a regra pressupõe uma segunda cabeça, e às vezes ela não
> existe. Quando o card volta para o MESMO agente, o claim **avisa e entrega** —
> travar não produziria o revisor que falta. O aviso fica no card, para quem ler
> depois saber que a segunda revisão não foi independente.

## 7. O merge

O card fecha pelo `Closes #N` do corpo, e o pipeline o move para
`ready-to-test` — o fim da alçada do Anchors.

Daí em diante, o processo de entrega é do seu projeto.

<!-- @anchors
  code: WORKR
  updated_at: 2026-09-28
-->
# Worker — o ciclo de quem executa uma tarefa da fila

> **Código**: `WORKR`

O fluxo mais percorrido do Anchors: toda tarefa passa por ele.

Um FLUXO não redesenha o trabalho — ele ENCAIXA ações, e diz o que fazer com cada
resultado que elas oferecem. As ações vivem em `flows/actions/` e são reusáveis: o
`map build` daqui é a mesma peça que o fluxo de adoção encaixa, escrita uma vez.

O que o fluxo acrescenta é a LIGAÇÃO — e é nela que as regras deixam de depender de
memória. Uma regra que hoje é prosa ("NUNCA feche com bloqueante vermelho") vira a
ausência de encaixe: o resultado BARRADO não tem ligação para `done`.

## Montagem

### WORKR-T01 — puxar a próxima tarefa

Encaixa: `ACNXT` (`anchors next`)

Resultados:
- `ACNXT-O01` TAREFA PUXADA → `WORKR-T02`
- `ACNXT-O02` FILA VAZIA → `WORKR-T07`

### WORKR-T02 — pôr no mapa o arquivo que a tarefa cita

Encaixa: `ACMAP` (`anchors map build`)

`impact` e `check` leem o MAPA, não o disco: um arquivo novo só existe para eles depois
desta peça. Sem ela, os dois respondem "não está no mapa" e a rodada se perde procurando
o motivo.

Resultados:
- `ACMAP-O01` MAPA EM DIA → `WORKR-T03`

### WORKR-T03 — escrever o artefato da etapa

Encaixa: `ACWRK` (`anchors work <etapa> --for <alvo>`)

A etapa diz o que se escreve — `specify` a spec, `implement` código e feature, `test` os
testes, `verify` nada (só confronta).

Resultados:
- `ACWRK-O01` ARTEFATO ESCRITO → `WORKR-T04`
- `ACWRK-O02` ALVO FORA DA ESTRUTURA → `WORKR-T08`

### WORKR-T04 — confrontar

Encaixa: `ACHCK` (`anchors check --changed`)

Resultados:
- `ACHCK-O01` PROMOVÍVEL → `WORKR-T06`
- `ACHCK-O02` BARRADO → `WORKR-T03` (volta a escrever — e a issue aberta segue no fluxo `ISSUE`)
- `ACHCK-O03` JULGAMENTO PENDENTE → `WORKR-T05`
- `ACHCK-O04` FORA DA ESTRUTURA → `WORKR-T08`
- `ACHCK-O05` MAPA DESATUALIZADO → `WORKR-T02` (refaz o mapa e confronta de novo)

### WORKR-T05 — julgar o que nenhum script computa

Encaixa: o fluxo `JUDGE` inteiro (um fluxo encaixa como peça)

Resultados:
- `JUDGE-O01` VEREDITO DADO → `WORKR-T04` (confronta de novo)

### WORKR-T06 — fechar a tarefa

Encaixa: `ACDON` (`anchors done <id>`)

Resultados:
- `ACDON-O01` FECHADA → `WORKR-T09`
- `ACDON-O02` SEM ID → `WORKR-T06` (a tarefa puxada tem id; use o dele)

### WORKR-T09 — o vigia já enfileirou a próxima etapa

Encaixa: `ACWTC` (o vigia, que dispara sozinho)

Este passo não é trabalho de ninguém — ele já aconteceu. Quando o artefato foi salvo em
`WORKR-T03`, o vigia classificou a mudança e enfileirou a etapa seguinte
(spec→implement, feature→test).

Está no fluxo justamente porque é o que faz o ciclo se sustentar. Sem ele desenhado, o
`worker` pareceria depender de alguém lembrar o que vem depois — e o mecanismo que
substitui essa memória ficaria invisível.

Resultados:
- `ACWTC-O01` TAREFA ENFILEIRADA → `WORKR-T01` (volta a puxar)
- `ACWTC-O02` IGNORADO → `WORKR-T01` (a mudança não pedia trabalho; a fila decide)

### WORKR-T07 — nada a fazer

Não é erro: `anchors next` sai com código 0. Fim legítimo da rodada.

> @terminal

### WORKR-T08 — o alvo não é regido pela Estrutura

Não é "passou". Ou falta declarar a camada, ou o arquivo não devia estar ali — e as duas
são decisão de quem conhece o projeto, não do worker.

> @terminal

<!-- @anchors
  code: WORKR
  updated_at: 2026-09-21
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

### WORKR-P01 — puxar a próxima tarefa

Encaixa: `ACNXT` (`anchors next`)

Resultados:
- `ACNXT-R01` TAREFA PUXADA → `WORKR-P02`
- `ACNXT-R02` FILA VAZIA → `WORKR-P07`

### WORKR-P02 — pôr no mapa o arquivo que a tarefa cita

Encaixa: `ACMAP` (`anchors map build`)

`impact` e `check` leem o MAPA, não o disco: um arquivo novo só existe para eles depois
desta peça. Sem ela, os dois respondem "não está no mapa" e a rodada se perde procurando
o motivo.

Resultados:
- `ACMAP-R01` MAPA EM DIA → `WORKR-P03`

### WORKR-P03 — escrever o artefato da etapa

Encaixa: `ACWRK` (`anchors work <etapa> --for <alvo>`)

A etapa diz o que se escreve — `specify` a spec, `implement` código e feature, `test` os
testes, `verify` nada (só confronta).

Resultados:
- `ACWRK-R01` ARTEFATO ESCRITO → `WORKR-P04`
- `ACWRK-R02` ALVO FORA DA ESTRUTURA → `WORKR-P08`

### WORKR-P04 — confrontar

Encaixa: `ACHCK` (`anchors check --changed`)

Resultados:
- `ACHCK-R01` PROMOVÍVEL → `WORKR-P06`
- `ACHCK-R02` BARRADO → `WORKR-P03` (volta a escrever — e a issue aberta segue no fluxo `ISSUE`)
- `ACHCK-R03` JULGAMENTO PENDENTE → `WORKR-P05`
- `ACHCK-R04` FORA DA ESTRUTURA → `WORKR-P08`
- `ACHCK-R05` MAPA DESATUALIZADO → `WORKR-P02` (refaz o mapa e confronta de novo)

### WORKR-P05 — julgar o que nenhum script computa

Encaixa: o fluxo `JUDGE` inteiro (um fluxo encaixa como peça)

Resultados:
- `JUDGE-R01` VEREDITO DADO → `WORKR-P04` (confronta de novo)

### WORKR-P06 — fechar a tarefa

Encaixa: `ACDON` (`anchors done <id>`)

Resultados:
- `ACDON-R01` FECHADA → `WORKR-P09`
- `ACDON-R02` SEM ID → `WORKR-P06` (a tarefa puxada tem id; use o dele)

### WORKR-P09 — o vigia já enfileirou a próxima etapa

Encaixa: `ACWTC` (o vigia, que dispara sozinho)

Este passo não é trabalho de ninguém — ele já aconteceu. Quando o artefato foi salvo em
`WORKR-P03`, o vigia classificou a mudança e enfileirou a etapa seguinte
(spec→implement, feature→test).

Está no fluxo justamente porque é o que faz o ciclo se sustentar. Sem ele desenhado, o
`worker` pareceria depender de alguém lembrar o que vem depois — e o mecanismo que
substitui essa memória ficaria invisível.

Resultados:
- `ACWTC-R01` TAREFA ENFILEIRADA → `WORKR-P01` (volta a puxar)
- `ACWTC-R02` IGNORADO → `WORKR-P01` (a mudança não pedia trabalho; a fila decide)

### WORKR-P07 — nada a fazer

Não é erro: `anchors next` sai com código 0. Fim legítimo da rodada.

> @terminal

### WORKR-P08 — o alvo não é regido pela Estrutura

Não é "passou". Ou falta declarar a camada, ou o arquivo não devia estar ali — e as duas
são decisão de quem conhece o projeto, não do worker.

> @terminal

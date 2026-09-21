<!-- @anchors
  code: PLANO
  updated_at: 2026-09-21
-->
# Plano — do escopo ao arquivo

> **Código**: `PLANO`

O único fluxo com REGRA DE OURO explícita, e ela existe porque gravar o arquivo já dispara
a máquina: o vigia detecta o plano e enfileira a primeira tarefa.

O guia diz, em maiúsculas: *"não escreva o arquivo para mostrar como ficou e pergunte
depois"*. Como fluxo, isso deixa de ser disciplina — `PLANO-P01` não tem saída para
`PLANO-P03`. Só se chega ao arquivo passando pela aprovação.

## Montagem

### PLANO-P01 — rascunhar o escopo NA CONVERSA

Encaixa: `ACGDE` (`anchors guide plan`)

Com o usuário, rascunhe EM TEXTO: objetivo, razão, a LISTA de specs que nascem ou mudam,
as fases, o que está fora de escopo, e a definição de pronto. Nada de arquivo ainda.

Resultados:
- `ACGDE-R01` RÉGUA LIDA → `PLANO-P02`

### PLANO-P02 — iterar até o usuário APROVAR

O escopo é do usuário, não de quem escreve. Iterar aqui é barato; iterar depois que a
esteira partiu custa uma rodada de trabalho enfileirado sobre a decisão errada.

Resultados:
- `PLANO-P02` (o usuário pediu ajuste) → `PLANO-P01`
- `PLANO-P03` (o usuário APROVOU o escopo) → `PLANO-P03`

### PLANO-P03 — escrever o arquivo do plano

Encaixa: `ACNPL` (`anchors new plan`)

Resultados:
- `ACNPL-R01` PLANO SEMEADO → `PLANO-P04`

### PLANO-P04 — a esteira partiu

O vigia enfileirou a primeira tarefa. Daqui em diante o trabalho flui pela fila, e quem
executa é o fluxo `WORKR`.

> @terminal

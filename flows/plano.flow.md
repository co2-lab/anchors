<!-- @anchors
  code: PLANO
  updated_at: 2026-09-28
-->
# Plano — do escopo ao arquivo

> **Código**: `PLANO`

O único fluxo com REGRA DE OURO explícita, e ela existe porque gravar o arquivo já dispara
a máquina: o vigia detecta o plano e enfileira a primeira tarefa.

O guia diz, em maiúsculas: *"não escreva o arquivo para mostrar como ficou e pergunte
depois"*. Como fluxo, isso deixa de ser disciplina — `PLANO-T01` não tem saída para
`PLANO-T03`. Só se chega ao arquivo passando pela aprovação.

## Montagem

### PLANO-T01 — rascunhar o escopo NA CONVERSA

Encaixa: `ACGDE` (`anchors guide plan`)

Com o usuário, rascunhe EM TEXTO: objetivo, razão, a LISTA de specs que nascem ou mudam,
as fases, o que está fora de escopo, e a definição de pronto. Nada de arquivo ainda.

Resultados:
- `ACGDE-O01` RÉGUA LIDA → `PLANO-T02`

### PLANO-T02 — iterar até o usuário APROVAR

O escopo é do usuário, não de quem escreve. Iterar aqui é barato; iterar depois que a
esteira partiu custa uma rodada de trabalho enfileirado sobre a decisão errada.

Resultados:
- `PLANO-T02` (o usuário pediu ajuste) → `PLANO-T01`
- `PLANO-T03` (o usuário APROVOU o escopo) → `PLANO-T03`

### PLANO-T03 — escrever o arquivo do plano

Encaixa: `ACNPL` (`anchors new plan`)

Resultados:
- `ACNPL-O01` PLANO SEMEADO → `PLANO-T04`

### PLANO-T04 — a esteira partiu

O vigia enfileirou a primeira tarefa. Daqui em diante o trabalho flui pela fila, e quem
executa é o fluxo `WORKR`.

> @terminal

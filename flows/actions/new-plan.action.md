<!-- @anchors
  code: ACNPL
  updated_at: 2026-09-21
-->
# Ação: `anchors new plan` — escrever o arquivo do plano

> **Código**: `ACNPL`

Emite o esqueleto do plano e o grava.

Gravar é um ATO IRREVERSÍVEL de escopo: assim que o arquivo existe na camada de plano, o
vigia o detecta e ENFILEIRA a primeira tarefa. A máquina começa a andar.

Por isso a aprovação do usuário vem ANTES de o arquivo existir, e não depois.

## Comando

    anchors new plan <Nome> --out plans/<nnnn>-<nome>.md

## Resultados

### ACNPL-R01 — PLANO SEMEADO: o arquivo existe e a esteira partiu

O vigia enfileirou a primeira tarefa ("specify"). A partir daqui o trabalho flui pela
fila, e quem executa é o `worker`.

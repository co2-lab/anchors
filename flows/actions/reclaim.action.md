<!-- @anchors
  code: ACRCL
  updated_at: 2026-09-21
-->
# Ação: `anchors reclaim` — devolver as tarefas do worker que morreu

> **Código**: `ACRCL`

Uma tarefa reivindicada por uma sessão que morreu sem fechar fica presa: ninguém a tem, e
ninguém pode pegá-la. Esta ação as devolve à fila, e o `anchors next` volta a alcançá-las.

É a única causa de travamento que não passa por decisão nenhuma — é mecânica, e o conserto
também.

## Comando

    anchors reclaim

## Resultados

### ACRCL-R01 — TAREFAS DEVOLVIDAS: a fila volta a alcançá-las

### ACRCL-R02 — NADA PRESO: nenhuma reivindicação órfã

Não é erro: é a ausência do problema.

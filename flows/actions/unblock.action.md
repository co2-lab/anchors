<!-- @anchors
  code: ACUNB
  updated_at: 2026-09-21
-->
# Ação: `anchors unblock` — a decisão que GERA trabalho

> **Código**: `ACUNB`

Às vezes a decisão do usuário não libera o card: ela CRIA trabalho antes — uma spec que
ganha regra, um contrato que muda, um defeito noutro lugar.

O trabalho vira card próprio, e o bloqueado passa a esperar POR ELE em vez de esperar
indefinidamente por alguém que já decidiu.

## Comando

    anchors unblock --card <n> ...

## Resultados

### ACUNB-R01 — CARD DE DESBLOQUEIO ABERTO: o bloqueado espera trabalho, não pessoa

Nasce com `anchors:desbloqueia-<n>`. Quando ele fecha, o bloqueado volta à fila.

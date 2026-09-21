<!-- @anchors
  code: ACNXT
  updated_at: 2026-09-21
-->
# Ação: `anchors next` — puxar a próxima tarefa da fila

> **Código**: `ACNXT`

Reivindica ATOMICAMENTE a próxima tarefa pendente e a imprime. Dois workers em terminais
diferentes nunca pegam a mesma — por isso dá para rodar em paralelo.

## Comando

    anchors next

## Resultados

### ACNXT-R01 — TAREFA PUXADA: reivindicada, e é só desta sessão

Quem puxou é dono até fechar (`done`) ou devolver (`reclaim`, se a sessão morrer).

### ACNXT-R02 — FILA VAZIA: não há trabalho pendente

Sai com código 0 — medido. Não é erro nem falha: é o fim legítimo da rodada.

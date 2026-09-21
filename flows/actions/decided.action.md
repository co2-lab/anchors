<!-- @anchors
  code: ACDEC
  updated_at: 2026-09-21
-->
# Ação: `anchors decided` — soltar o card que esperava a decisão

> **Código**: `ACDEC`

A decisão saiu. Tira o rótulo `anchors:needs-user` (é ele que faz a reivindicação pular o
card), comenta a resolução, e o card volta à fila.

A resolução é comentada no card, e não só anunciada: quem pegar depois precisa saber o que
mudou, sem ter de reconstruir a conversa.

## Comando

    anchors decided --card <n> --resolution "<o que foi decidido>"

## Resultados

### ACDEC-R01 — CARD LIBERADO: volta à fila com a resolução escrita

<!-- @anchors
  code: ACJPN
  updated_at: 2026-09-21
-->
# Ação: `anchors judge --pending` — ver o que aguarda julgamento

> **Código**: `ACJPN`

Lista os alvos que o `check` marcou com `⏳` e enfileirou, e em cada um o guia a ler e a
pergunta a responder.

## Comando

    anchors judge --pending

## Resultados

### ACJPN-R01 — HÁ ALVO AGUARDANDO: o guia e a pergunta vêm junto

### ACJPN-R02 — NADA AGUARDANDO: nenhum gate de julgamento pendente

Medido: sai com código 0 e manda rodar o `check` para descobrir. Não é erro — é a
ausência de trabalho deste tipo.

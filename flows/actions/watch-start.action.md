<!-- @anchors
  code: ACWCH
  updated_at: 2026-09-21
-->
# Ação: `anchors watch start` — pôr a esteira no ar

> **Código**: `ACWCH`

O vigia em segundo plano: toda mudança de arquivo vira tarefa na fila — inclusive as suas.

É redundante quando foi você quem mudou (você já sabia). Mas é o MESMO caminho de quando a
mudança vem de fora: o usuário editou no editor, um `git pull` trouxe algo. Um mecanismo
só, para toda origem — e é o que permite confiar na fila em vez da memória.

## Comando

    anchors watch start

## Resultados

### ACWCH-R01 — VIGIA NO AR: as mudanças viram tarefa

### ACWCH-R02 — JÁ ESTAVA RODANDO: nada a fazer

Não é erro: dois vigias sobre o mesmo projeto enfileirariam a mesma tarefa duas vezes.

<!-- @anchors
  code: ACDON
  updated_at: 2026-09-21
-->
# Ação: `anchors done` — fechar a tarefa

> **Código**: `ACDON`

Move a tarefa reivindicada para o histórico (`.anchors/done/`).

## Comando

    anchors done <id>

## Resultados

### ACDON-R01 — FECHADA: a tarefa foi para o histórico

Salvar o artefato já fez o watcher enfileirar a PRÓXIMA etapa (spec→implement,
feature→test). O ciclo se sustenta sozinho: ninguém precisa lembrar o que vem depois.

### ACDON-R02 — SEM ID: a ação não sabe o que fechar

Medido: sem `<id>` a ação erra pedindo o identificador ou um filtro (`--file`, `--kind`,
`--all`). Fechar a tarefa errada é pior que não fechar — a fila passaria a mentir.

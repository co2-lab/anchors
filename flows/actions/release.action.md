<!-- @anchors
  code: ACREL
  updated_at: 2026-09-21
-->
# Ação: o RELEASE — publicar a versão

> **Código**: `ACREL`

Dispara no push de uma tag `v*`. Compila os binários e publica.

Não tem ramo: ou publica, ou falha. É o fim do fluxo de entrega, e o único gatilho aqui
que não decide nada — ele executa o que a tag já decidiu.

## Gatilho

    git push origin v<versão>     (.github/workflows/release-cli.yml)

## Resultados

### ACREL-R01 — PUBLICADO: os binários estão disponíveis

### ACREL-R02 — FALHOU: o build da release quebrou

A tag existe e a release não. O conserto é corrigir e re-tagar — apagar a tag publicada
quebraria quem já a baixou.

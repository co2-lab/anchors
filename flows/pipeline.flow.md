<!-- @anchors
  code: PIPEL
  updated_at: 2026-09-28
-->
# Pipeline — do PR aberto à versão publicada

> **Código**: `PIPEL`

O que a MÁQUINA faz. Nenhum passo aqui é comando digitado: tudo dispara sozinho, a partir
de um PR aberto ou de uma tag empurrada.

É o fluxo com a propriedade mais incômoda: quando ele barra, quem precisa agir não está
olhando. Por isso cada resultado que barra aponta de volta para o fluxo de quem trabalha.

## Montagem

### PIPEL-T01 — o CI confronta o PR

Encaixa: `ACCIP` (o CI, que dispara em `pull_request`)

Resultados:
- `ACCIP-O01` PR VERDE → `PIPEL-T02`
- `ACCIP-O02` MAPA VELHO → `PIPEL-T05`
- `ACCIP-O03` GATE REPROVOU → `PIPEL-T05`
- `ACCIP-O04` ISSUE PARA TRÁS → `PIPEL-T05`

### PIPEL-T02 — o PR está pronto para revisão

O confronto automático passou. O que falta é humano, e não é deste fluxo.

Resultados:
- `PIPEL-T03` (o PR foi aprovado e mergeado) → `PIPEL-T03`

### PIPEL-T03 — mergeado na main

Resultados:
- `PIPEL-T04` (a tag de versão foi empurrada) → `PIPEL-T04`

### PIPEL-T04 — publicar

Encaixa: `ACREL` (o release, que dispara no push da tag)

Resultados:
- `ACREL-O01` PUBLICADO → `PIPEL-T06`
- `ACREL-O02` FALHOU → `PIPEL-T07`

### PIPEL-T05 — a bola voltou para quem trabalha

O pipeline barrou, e ninguém está olhando para ele. Quem tem de agir é quem abriu o PR — e
o caminho de volta é o fluxo `ENTRG`, corrigindo e commitando de novo.

> @terminal

### PIPEL-T06 — versão publicada

> @terminal

### PIPEL-T07 — a tag existe e a release não

O conserto é corrigir e re-tagar. Apagar a tag publicada quebraria quem já a baixou.

> @terminal

<!-- @anchors
  code: PIPEL
  updated_at: 2026-09-21
-->
# Pipeline — do PR aberto à versão publicada

> **Código**: `PIPEL`

O que a MÁQUINA faz. Nenhum passo aqui é comando digitado: tudo dispara sozinho, a partir
de um PR aberto ou de uma tag empurrada.

É o fluxo com a propriedade mais incômoda: quando ele barra, quem precisa agir não está
olhando. Por isso cada resultado que barra aponta de volta para o fluxo de quem trabalha.

## Montagem

### PIPEL-P01 — o CI confronta o PR

Encaixa: `ACCIP` (o CI, que dispara em `pull_request`)

Resultados:
- `ACCIP-R01` PR VERDE → `PIPEL-P02`
- `ACCIP-R02` MAPA VELHO → `PIPEL-P05`
- `ACCIP-R03` GATE REPROVOU → `PIPEL-P05`
- `ACCIP-R04` ISSUE PARA TRÁS → `PIPEL-P05`

### PIPEL-P02 — o PR está pronto para revisão

O confronto automático passou. O que falta é humano, e não é deste fluxo.

Resultados:
- `PIPEL-P03` (o PR foi aprovado e mergeado) → `PIPEL-P03`

### PIPEL-P03 — mergeado na main

Resultados:
- `PIPEL-P04` (a tag de versão foi empurrada) → `PIPEL-P04`

### PIPEL-P04 — publicar

Encaixa: `ACREL` (o release, que dispara no push da tag)

Resultados:
- `ACREL-R01` PUBLICADO → `PIPEL-P06`
- `ACREL-R02` FALHOU → `PIPEL-P07`

### PIPEL-P05 — a bola voltou para quem trabalha

O pipeline barrou, e ninguém está olhando para ele. Quem tem de agir é quem abriu o PR — e
o caminho de volta é o fluxo `ENTRG`, corrigindo e commitando de novo.

> @terminal

### PIPEL-P06 — versão publicada

> @terminal

### PIPEL-P07 — a tag existe e a release não

O conserto é corrigir e re-tagar. Apagar a tag publicada quebraria quem já a baixou.

> @terminal

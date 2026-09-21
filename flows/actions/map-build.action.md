<!-- @anchors
  code: ACMAP
  updated_at: 2026-09-21
-->
# Ação: `anchors map build` — pôr o disco no mapa

> **Código**: `ACMAP`

Varre o projeto e reconstrói `anchors.graph.yaml` — os nós, as arestas, e os carimbos de
validação que o rebuild preserva.

É a peça que mais se repete entre fluxos, e a razão é sempre a mesma: `impact` e `check`
leem o MAPA, não o disco. Um arquivo novo só existe para eles depois desta ação.

## Comando

    anchors map build

## Resultados

### ACMAP-R01 — MAPA EM DIA: o grafo reflete o disco

### ACMAP-R02 — CARIMBO PERDIDO: o rebuild não achou um nó que tinha validação

Aviso, não reprovação: remover um nó legitimamente remove os carimbos dele. O que a ação
faz é transformar perda silenciosa em perda visível — o laudo de uma revisão adversarial é
caro de refazer.

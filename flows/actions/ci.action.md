<!-- @anchors
  code: ACCIP
  updated_at: 2026-09-21
-->
# Ação: o CI — o confronto que ninguém pode pular

> **Código**: `ACCIP`

Dispara em `pull_request`. O pre-commit pode ser contornado (`--no-verify`); o CI não.

Oito passos, e os três últimos são o Anchors sobre o projeto: se o mapa está em dia, se os
gates passam, e se o check deixou issue para trás.

## Gatilho

    pull_request        (.github/workflows/ci.yml)

## Resultados

### ACCIP-R01 — PR VERDE: compila, testa, e os gates passam

### ACCIP-R02 — MAPA VELHO: o mapa não reflete o que o PR traz

Os gates deste PR estariam confrontando uma foto vencida — o veredito pode estar errado
nas duas direções.

Sugere: `anchors map build`, e commitar o mapa junto.

### ACCIP-R03 — GATE REPROVOU: um bloqueante falhou no confronto completo

### ACCIP-R04 — ISSUE PARA TRÁS: o check abriu achado que o PR não fechou

O trabalho está feito e algo ficou aberto. Não é o mesmo que gate vermelho: aqui o
confronto passou e deixou pendência registrada.

Sugere: tratar o que ficou, ou declarar por que fica.

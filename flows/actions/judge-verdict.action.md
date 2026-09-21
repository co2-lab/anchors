<!-- @anchors
  code: ACJVD
  updated_at: 2026-09-21
-->
# Ação: `anchors judge --verdict` — registrar o veredito

> **Código**: `ACJVD`

Grava o veredito de um gate de julgamento sobre um alvo. O `--reason` NÃO é uma frase de
veredito: é o RELATÓRIO inteiro, e ele vira o corpo da issue, literal.

O carimbo leva a revisão do alvo — então o veredito ENVELHECE se o alvo mudar, e o gate
volta a cobrar. É o que impede um "pass" de valer para sempre sobre um arquivo que já é
outro.

## Comando

    anchors judge <alvo> --gate <g> --verdict pass|fail|waived --reason "<RELATÓRIO>"

## Resultados

### ACJVD-R01 — APROVADO: o alvo cumpre os pontos de conformidade

Resolve a issue anterior, se havia.

### ACJVD-R02 — REPROVADO: o alvo viola um ou mais pontos

Abre a issue com o relatório. É por isso que o relatório importa: sem ele, alguém vai
reprocessar o alvo depois só para descobrir o que consertar.

### ACJVD-R03 — DISPENSADO: o ponto não se aplica a este alvo

O terceiro veredito. Existe porque um guia governa vários tipos de alvo, e um ponto
escrito para tela não se aplica a um modelo — forçar `pass` ou `fail` ali seria afirmar
sobre o que não foi medido.

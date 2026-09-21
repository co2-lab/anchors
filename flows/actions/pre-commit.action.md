<!-- @anchors
  code: ACPRC
  updated_at: 2026-09-21
-->
# Ação: o PRE-COMMIT — o gate que roda sem ninguém chamar

> **Código**: `ACPRC`

Instalado por `anchors install-hooks`, e a partir daí ninguém o digita: ele dispara no
`git commit`, sobre os arquivos em stage.

Roda `check --changed --deterministic --no-record` em cada um. Incremental: valida só o
que o commit toca. Não escreve no mapa nem abre issue.

## Gatilho

    git commit          sobre os arquivos em stage

## Resultados

### ACPRC-R01 — COMMIT PASSOU: nenhum bloqueante vermelho

### ACPRC-R02 — BARRADO POR GATE: um bloqueante reprovou no que está em stage

Sugere: corrigir, ou declarar a dispensa na mensagem do commit com `[skip-<regra>@<CÓDIGO>: por quê]`.

### ACPRC-R03 — BARRADO POR FALTA DE MAPA: arquivo REGIDO e fora do mapa

O caso que pega quem não conhece. Um arquivo novo que casa uma camada, mas que o
`map build` ainda não registrou: fora do mapa NENHUM gate o confronta, e o commit seguiria
certificando trabalho que ninguém verificou.

Sugere: `anchors map build`, e commitar de novo.

### ACPRC-R04 — IGNORADO: arquivo não regido pela Estrutura

`package.json`, lockfile, arquivo de CI. O Anchors não tem jurisdição sobre eles, e barrar
aqui seria cobrar régua de quem nunca a aceitou.

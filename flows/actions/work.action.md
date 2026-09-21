<!-- @anchors
  code: ACWRK
  updated_at: 2026-09-21
-->
# Ação: `anchors work` — compor o procedimento da etapa

> **Código**: `ACWRK`

Entrega a QUEM TRABALHA o que ler, em que ordem, onde as peças nascem, o que não é escopo,
e os gates que virão. Composto do `anchors.yaml` — nada é inventado aqui.

## Comando

    anchors work <etapa> --for <alvo>       etapa: spec|code|feature|test|review

## Resultados

### ACWRK-R01 — PROCEDIMENTO COMPOSTO: quem trabalha sabe o que fazer

Escrever o artefato é trabalho de quem lê — o Anchors não gera conteúdo.

### ACWRK-R02 — ALVO FORA DA ESTRUTURA: o procedimento sai genérico

Medido: `work spec --for /tmp/nao-existe.go` compõe procedimento para um arquivo que não
existe e não casa camada. O texto sai, mas sem camada resolvida ele perde o que tem de
mais útil — os guias daquela camada, os caminhos das peças, os gates que vão cobrar.

Um procedimento genérico parece procedimento. É o resultado a tratar, não a ignorar.

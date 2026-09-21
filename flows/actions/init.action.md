<!-- @anchors
  code: ACINI
  updated_at: 2026-09-21
-->
# Ação: `anchors init` — configurar o projeto

> **Código**: `ACINI`

Varre o projeto DETERMINISTICAMENTE (sem IA), propõe uma Estrutura, e confirma com o
usuário por perguntas — chegando a um `anchors.yaml` correto.

O grosso é inferido do disco. As perguntas cobrem só as decisões HUMANAS: co-locação,
granularidade de camada, e quais guias regem quais tags.

É o gatilho do fluxo de adoção: antes dele não há decisão nenhuma a tomar, e todas as
outras nascem dele.

## Comando

    anchors init                    interativo
    anchors init --non-interactive  emite as perguntas como JSON, ou aplica as respostas

## Resultados

### ACINI-R01 — ESTRUTURA INFERIDA: o projeto tem código e o disco respondeu

### ACINI-R02 — DIRETÓRIO VAZIO: não há de onde inferir

O projeto ainda não existe. Não há camada a propor, nem dialeto a detectar — e a
Estrutura que saísse daqui seria adivinhação com cara de configuração.

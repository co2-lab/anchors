<!-- @anchors
  code: ACGDE
  updated_at: 2026-09-21
-->
# Ação: `anchors guide` — ler a régua antes de escrever

> **Código**: `ACGDE`

Imprime a doutrina de um artefato: o que ele é, o que não é, e os pontos de conformidade
contra os quais ele será julgado.

## Comando

    anchors guide <spec|code|feature|test|plan|product|header|project|guide>

## Resultados

### ACGDE-R01 — RÉGUA LIDA: quem vai escrever sabe o critério

### ACGDE-R02 — GUIA SEM PONTOS DE CONFORMIDADE: só há prosa a interpretar

O guia existe mas não distilou as regras em itens verificáveis. Julgar por prosa inteira
"no olho" é o que torna o veredito irreprodutível — dois julgadores, dois resultados.

É dívida declarada, e o gate `guide-checklist` já a acusa. Quem julga faz o que dá com o
que está escrito, e diz que julgou sobre prosa.

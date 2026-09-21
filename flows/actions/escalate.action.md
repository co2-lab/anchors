<!-- @anchors
  code: ACESC
  updated_at: 2026-09-21
-->
# Ação: `anchors escalate` — abrir a questão que o trabalho encontrou

> **Código**: `ACESC`

Quem executa encontra uma divergência: o plano se contradiz, ou não cobriu algo. Quem
encontrou INTERPRETA o impacto, e a interpretação escolhe a saída.

A pergunta não é "tenho certeza?" — é **"há mais de uma resposta defensável, e escolher
entre elas muda o que o produto FAZ?"**.

## Comando

    anchors escalate --for-user   impacta a DIREÇÃO do produto
    anchors escalate --unsure     não sei se impacta
    anchors escalate              não impacta: vira card comum

## Resultados

### ACESC-R01 — CARD COMUM: não impacta a direção

Nasce em "a fazer", entra na fila, um agente pega. Não para ninguém.

### ACESC-R02 — AGUARDANDO O USUÁRIO: impacta a direção

Nasce com `anchors:needs-user`, e a reivindicação não entrega o card enquanto a decisão
não sair.

### ACESC-R03 — AGUARDANDO ENQUADRAMENTO: não se sabe se impacta

Nasce igual ao anterior, mas a primeira pergunta do card é outra: *é seu?* Quem lê, se
concluir que não impacta, devolve à fila em vez de decidir.

Duas saídas para o que para o card, e a razão está escrita no comando: quem lê trinta
decisões abertas precisa saber o que cada uma PEDE. "Decida entre A e B" e "veja se isto
é seu" são trabalhos diferentes, e misturá-los faz o segundo custar como o primeiro.

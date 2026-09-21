<!-- @anchors
  code: JUDGE
  updated_at: 2026-09-21
-->
# Judge — o veredito que nenhum script computa

> **Código**: `JUDGE`

Alguns gates medem o que nenhum script sabe: "esta tela se decompõe em atomic design?",
"a spec descreve comportamento e não implementação?". O `check` não os computa — marca o
alvo com `⏳` e enfileira. Quem julga é uma IA, contra os pontos de conformidade do guia.

Este fluxo encaixa no `worker` como peça: `WORKR-P05` é ele inteiro.

O que a topologia garante aqui é uma coisa só, e é a que mais se esquece: **ler o guia é o
único caminho até o veredito**. Hoje isso é uma advertência em prosa ("NÃO julgue pela
prosa inteira no olho — julgue contra CADA ponto"), e advertência só funciona se alguém a
ler antes de julgar.

## Montagem

### JUDGE-P01 — ver o que aguarda julgamento

Encaixa: `ACJPN` (`anchors judge --pending`)

Resultados:
- `ACJPN-R01` HÁ ALVO AGUARDANDO → `JUDGE-P02`
- `ACJPN-R02` NADA AGUARDANDO → `JUDGE-P05`

### JUDGE-P02 — ler a régua do que rege este alvo

Encaixa: `ACGDE` (`anchors guide <o que rege>`)

Vá direto à seção de pontos de conformidade: é a lista de itens a verificar. Se a
checklist estiver AGRUPADA POR ALVO, aplique só os pontos do grupo desta camada mais os
do grupo "para todos" — um ponto de tela não se aplica a um modelo.

Resultados:
- `ACGDE-R01` RÉGUA LIDA → `JUDGE-P03`
- `ACGDE-R02` GUIA SEM PONTOS DE CONFORMIDADE → `JUDGE-P03`

### JUDGE-P03 — julgar o alvo contra CADA ponto

Não devolva uma frase: devolva o RELATÓRIO por item. Para cada ponto — cumpre, ou a não
conformidade com o quê, ONDE (`arquivo:linha`), qual ponto violou, e como consertar.

Esse texto vira o corpo da issue. Sem ele, alguém vai reprocessar o alvo depois só para
descobrir o que consertar — e o custo do julgamento se paga duas vezes.

Resultados:
- `JUDGE-P03` (a decisão é de quem julga) → `JUDGE-P04`

### JUDGE-P04 — registrar o veredito

Encaixa: `ACJVD` (`anchors judge --verdict`)

O carimbo leva a revisão do alvo: o veredito ENVELHECE se o alvo mudar, e o gate volta a
cobrar. Um "pass" não vale para sempre sobre um arquivo que já é outro.

Resultados:
- `ACJVD-R01` APROVADO → `JUDGE-P06`
- `ACJVD-R02` REPROVADO → `JUDGE-P06`
- `ACJVD-R03` DISPENSADO → `JUDGE-P06`

### JUDGE-P05 — nada a julgar

> @terminal

### JUDGE-P06 — veredito registrado

O `check` roda de novo: o gate de julgamento agora tem resposta, e o alvo volta ao fluxo
de quem o confrontou.

> @terminal

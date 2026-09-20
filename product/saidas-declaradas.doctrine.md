<!-- @anchors
  code: SAIDA
  updated_at: 2026-09-20
-->
# Saídas declaradas — como um gate cala, e o que o silêncio significa

> **Código**: `SAIDA`

## Overview

Todo gate do Anchors precisa de uma saída: existe sempre o caso legítimo que a régua não
previu, e um gate sem saída honesta é um gate que o time aprende a contornar — desligando,
ignorando, ou escrevendo o que ele quer ouvir.

Mas "calar o gate" não é uma coisa só. São **duas afirmações diferentes**, e confundi-las
faz o gate perder a informação mais cara que ele tem: a diferença entre *decidi que não
precisa* e *ainda não fiz*.

Esta doutrina é transversal por construção. Ela não pertence a `triad-complete` nem a
`plan-doctrine-exists` nem a `doctrine-realized` — vale para os três, e para todo gate que
venha a oferecer saída. Até aqui vivia duplicada no comentário de cada um, que é
exatamente a divergência que o eixo de produto existe para acabar.

## Rules

### SAIDA-R01 — a dispensa PERMANENTE se declara com `@no-<coisa>: <razão>`

Afirma que aquela exigência **nunca** se aplicará a esta unidade: o código é configuração
pura, o comportamento não é observável, a prova vive noutro lugar. É decisão tomada, e o
gate passa — `Pass`, definitivo.

### SAIDA-R02 — a dívida se declara com `@TBD: <o quê> — <razão>`

Afirma que a peça **ainda não existe**, e que alguém vai escrevê-la. É trabalho pendente,
não decisão: o gate devolve `Pending`, e o achado continua aparecendo até ser pago.

### SAIDA-R03 — dívida NUNCA vira `Pass`

Medido neste repositório antes da correção: `triad-complete` jogava `@TBD` no mesmo balde
do `@no-*`, e uma spec com `@TBD: code,feature,test` saía **verde**, indistinguível de uma
tríade completa. O trabalho que falta desaparecia do radar por causa da declaração honesta
de quem o assumiu — que é o pior incentivo possível.

### SAIDA-R04 — marcador nu não dispensa nada

`@no-code` sozinho, sem `:` e sem razão escrita, não é saída: é o silêncio que os gates
existem para acabar. A razão é obrigatória em ambos os marcadores, e é verificada por
regex — quem lê a unidade descobre ali por que ela está dispensada, sem caçar a decisão
noutro lugar.

### SAIDA-R05 — a saída vale onde está escrita

Um marcador numa linha dispensa aquela linha; no header, a unidade. Não há saída que
valha para o projeto inteiro escrita num arquivo distante — a decisão fica onde quem lê
vai encontrá-la.

### SAIDA-R06 — marcador entre crases é MENÇÃO, não declaração

Uma revisão que explica a remoção de uma dispensa cita o marcador (*"a dispensa `@TBD: code`
saiu"*), e sem essa distinção a citação **reativa** a dispensa que o texto diz ter
acabado. Marcador ativo nunca está entre crases.

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SAIDA-X01` | Nenhum gate infere saída pelo texto da mensagem. | A saída é um campo/marcador declarado, e não uma frase reconhecida — deduzir intenção de prosa envelhece na primeira reescrita da mensagem. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

nenhuma

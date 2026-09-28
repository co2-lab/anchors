<!-- @anchors
  code: ISSUE
  updated_at: 2026-09-28
-->
# Issue — o trabalho que um gate reprovado deixa para trás

> **Código**: `ISSUE`

Este fluxo existe porque o resultado de uma ação nem sempre morre no passo que o produziu.
Um gate bloqueante que reprova não faz só voltar o trabalho: ele ABRE UMA ISSUE, e ela
sobrevive à rodada, com dono, estado e fim próprios.

Os quatro estados são os do código (`internal/issue`), e a pasta em que a issue vive É o
estado — mover é mudar de estado.

A distinção entre `future` e `todo` é a que sustenta o resto, e a doutrina do código a
explica: *"quem olha `todo/` está perguntando 'o que faço agora', e afogar essa lista com
o que só vence depois é o caminho mais curto para ninguém mais olhar"*.

Encaixa como peça: `WORKR-T04` aponta para cá quando o `check` reprova.

## Montagem

### ISSUE-T01 — detectada, ninguém pegou

A issue nasceu de um gate reprovado ou de um veredito de julgamento. Está em `todo/`, e a
pergunta que ela responde é "o que faço agora".

Resultados:
- `ISSUE-T02` (alguém vai resolver agora) → `ISSUE-T02`
- `ISSUE-T05` (é dívida com prazo, não vence agora) → `ISSUE-T05`

### ISSUE-T02 — alguém está resolvendo

Encaixa: `ACAUD` (`anchors audit <arquivo>`)

Se vai abrir o arquivo, conserte TUDO o que ele tem aberto. Voltar três vezes ao mesmo
arquivo por três achados custa três leituras do contexto inteiro.

Resultados:
- `ACAUD-O01` HÁ PENDÊNCIA → `ISSUE-T03`
- `ACAUD-O02` NADA ABERTO → `ISSUE-T04`

### ISSUE-T03 — corrigir, e confrontar de novo

Encaixa: `ACHCK` (`anchors check --changed`)

Resultados:
- `ACHCK-O01` PROMOVÍVEL → `ISSUE-T04`
- `ACHCK-O02` BARRADO → `ISSUE-T03` (ainda não)

### ISSUE-T04 — tratada

Fato datado: a issue foi para `done/`. O `check` que a abriu, ao rodar de novo, a resolve.

> @terminal

### ISSUE-T05 — dívida assumida, com prazo

Não é `todo` (não é o que se faz agora) nem `done` (não foi feito). É um dever conhecido,
ainda válido, com um momento declarado para ser cumprido.

Nasce da distinção que o `obligation_pending` já fazia no cabeçalho e que nada
materializava: uma linha visível só para quem abrisse o arquivo, sem estado, sem como ser
paga, sem como vencer.

Resultados:
- `ISSUE-T01` (o prazo venceu: vira trabalho de agora) → `ISSUE-T01`

> @terminal

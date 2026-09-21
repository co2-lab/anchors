<!-- @anchors
  code: ISSUE
  updated_at: 2026-09-21
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

Encaixa como peça: `WORKR-P04` aponta para cá quando o `check` reprova.

## Montagem

### ISSUE-P01 — detectada, ninguém pegou

A issue nasceu de um gate reprovado ou de um veredito de julgamento. Está em `todo/`, e a
pergunta que ela responde é "o que faço agora".

Resultados:
- `ISSUE-P02` (alguém vai resolver agora) → `ISSUE-P02`
- `ISSUE-P05` (é dívida com prazo, não vence agora) → `ISSUE-P05`

### ISSUE-P02 — alguém está resolvendo

Encaixa: `ACAUD` (`anchors audit <arquivo>`)

Se vai abrir o arquivo, conserte TUDO o que ele tem aberto. Voltar três vezes ao mesmo
arquivo por três achados custa três leituras do contexto inteiro.

Resultados:
- `ACAUD-R01` HÁ PENDÊNCIA → `ISSUE-P03`
- `ACAUD-R02` NADA ABERTO → `ISSUE-P04`

### ISSUE-P03 — corrigir, e confrontar de novo

Encaixa: `ACHCK` (`anchors check --changed`)

Resultados:
- `ACHCK-R01` PROMOVÍVEL → `ISSUE-P04`
- `ACHCK-R02` BARRADO → `ISSUE-P03` (ainda não)

### ISSUE-P04 — tratada

Fato datado: a issue foi para `done/`. O `check` que a abriu, ao rodar de novo, a resolve.

> @terminal

### ISSUE-P05 — dívida assumida, com prazo

Não é `todo` (não é o que se faz agora) nem `done` (não foi feito). É um dever conhecido,
ainda válido, com um momento declarado para ser cumprido.

Nasce da distinção que o `obligation_pending` já fazia no cabeçalho e que nada
materializava: uma linha visível só para quem abrisse o arquivo, sem estado, sem como ser
paga, sem como vencer.

Resultados:
- `ISSUE-P01` (o prazo venceu: vira trabalho de agora) → `ISSUE-P01`

> @terminal

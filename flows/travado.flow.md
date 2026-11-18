<!-- @anchors
  code: TRAVA
  updated_at: 2026-09-28
-->
# Travado — o que fazer quando o trabalho para

> **Código**: `TRAVA`

O fluxo que NÃO existia escrito em lugar nenhum — nem no guia, nem nos comandos. Os
comandos existem (`escalate`, `decided`, `unblock`, `discard`, `reclaim`), cada um com a
sua doutrina no `--help`, e a SEQUÊNCIA entre eles nunca foi desenhada.

E travar é justamente quando ninguém lembra do procedimento: quem está travado está
travado porque algo saiu do previsto.

> **Este fluxo é DESENHO, não transcrição.** Os quatro anteriores foram lidos de um
> procedimento que já existia; este foi composto a partir do que cada comando declara
> fazer. Vale confrontá-lo com o uso real antes de tratá-lo como régua.

## Montagem

### TRAVA-T01 — o trabalho parou: por quê?

O primeiro enquadramento, e ele decide tudo o que vem depois. Três causas, e elas pedem
coisas diferentes:

- uma DIVERGÊNCIA no plano ou na spec (contradição, ou lacuna) → `TRAVA-T02`
- o card não faz mais SENTIDO → `TRAVA-T07`
- a sessão MORREU sem fechar a tarefa → `TRAVA-T08`

Resultados:
- `TRAVA-T02` (divergência no plano ou na spec) → `TRAVA-T02`
- `TRAVA-T07` (o card não faz mais sentido) → `TRAVA-T07`
- `TRAVA-T08` (o worker morreu sem fechar) → `TRAVA-T08`

### TRAVA-T02 — abrir a questão, escolhendo a saída

Encaixa: `ACESC` (`anchors escalate`)

A pergunta do enquadramento: *há mais de uma resposta defensável, e escolher entre elas
muda o que o produto FAZ?*

Resultados:
- `ACESC-O01` CARD COMUM → `TRAVA-T06`
- `ACESC-O02` AGUARDANDO O USUÁRIO → `TRAVA-T03`
- `ACESC-O03` AGUARDANDO ENQUADRAMENTO → `TRAVA-T03`

### TRAVA-T03 — o card espera uma pessoa

Enquanto o rótulo `anchors:needs-user` estiver no card, a reivindicação o pula. Ninguém
pega por engano, e ninguém decide no lugar de quem deve decidir.

Resultados:
- `TRAVA-T04` (a decisão saiu e LIBERA o card) → `TRAVA-T04`
- `TRAVA-T05` (a decisão saiu e GERA trabalho antes) → `TRAVA-T05`
- `TRAVA-T07` (a decisão foi: isto não faz mais sentido) → `TRAVA-T07`

### TRAVA-T04 — soltar o card com a resolução escrita

Encaixa: `ACDEC` (`anchors decided`)

Resultados:
- `ACDEC-O01` CARD LIBERADO → `TRAVA-T06`

### TRAVA-T05 — abrir o card que desbloqueia

Encaixa: `ACUNB` (`anchors unblock`)

O bloqueado passa a esperar TRABALHO, não pessoa. É a diferença entre uma espera com dono
e uma espera indefinida.

Resultados:
- `ACUNB-O01` CARD DE DESBLOQUEIO ABERTO → `TRAVA-T06`

### TRAVA-T06 — de volta à fila

O card é pegável de novo. Quem o executa é o fluxo `WORKR`.

> @terminal

### TRAVA-T07 — tirar do quadro

Encaixa: `ACDIS` (`anchors discard`)

Só para o que não faz mais sentido. Trabalho entregue se FECHA — descartar o que foi feito
apagaria o passado que o roadmap desenha.

Resultados:
- `ACDIS-O01` CARD DESCARTADO → `TRAVA-T09`

### TRAVA-T08 — devolver as tarefas do worker morto

Encaixa: `ACRCL` (`anchors reclaim`)

Tarefa reivindicada por uma sessão que morreu fica presa: ninguém a tem, e ninguém pode
pegá-la. Esta é a única causa de travamento que não passa por decisão nenhuma — é
mecânica, e o conserto também.

Resultados:
- `ACRCL-O01` TAREFAS DEVOLVIDAS → `TRAVA-T06`
- `ACRCL-O02` NADA PRESO → `TRAVA-T10`

### TRAVA-T10 — não era isso

Nenhuma reivindicação órfã: o travamento tem outra causa, e o enquadramento de `TRAVA-T01`
errou. Volte a ele com o que agora se sabe.

Resultados:
- `TRAVA-T01` (reenquadrar com o que se descobriu) → `TRAVA-T01`

### TRAVA-T09 — encerrado sem entrega

O card saiu do quadro com a razão registrada. Não é fracasso nem entrega: é uma decisão
tomada, e ela fica escrita.

> @terminal

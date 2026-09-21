<!-- @anchors
  code: ACWTC
  updated_at: 2026-09-21
-->
# Ação: o VIGIA — transformar "mudou" em "há trabalho"

> **Código**: `ACWTC`

Roda em segundo plano depois de `anchors watch start`. A cada mudança de arquivo,
classifica e ENFILEIRA uma tarefa, com a próxima etapa sugerida.

Não executa nem chama IA — essa é a inversão. Ele só diz que há trabalho; quem trabalha
puxa com `anchors next`.

É a peça que faz o ciclo do `worker` se fechar sozinho: quando uma tarefa fecha, foi o
vigia que já enfileirou a seguinte (spec→implement, feature→test). É por isso que ninguém
precisa lembrar o que vem depois.

## Gatilho

    (qualquer arquivo muda)     enquanto `anchors watch start` estiver no ar

## Resultados

### ACWTC-R01 — TAREFA ENFILEIRADA: a mudança virou trabalho pendente

Inclusive quando foi você quem mudou. É redundante — você já sabia —, e é o MESMO caminho
de quando a mudança vem de fora (o usuário editou, um `git pull` trouxe algo). Um mecanismo
só, para toda origem: é o que permite confiar na fila em vez da memória.

### ACWTC-R02 — IGNORADO: a mudança não pede trabalho

Arquivo fora das camadas regidas, ou mudança que não altera o que o mapa confronta.

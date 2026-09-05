---
title: Congelar o projeto
description: O botão de pânico — como parar todo o trabalho quando um problema precisa ser resolvido antes.
---

Às vezes um problema aparece e **ninguém deve trabalhar até ele ser resolvido**:
o plano aponta para uma spec que não existe, uma decisão de arquitetura se
revelou errada, uma credencial vazou.

Sem um freio, o time continua produzindo contra uma base que vai mudar — e o
trabalho de horas envelhece antes de ser entregue.

```sh
anchors freeze --motivo "o plano 0002 aponta para uma spec que não existe — ver #42"
```

Para liberar:

```sh
anchors thaw
```

## Quatro camadas, e nenhuma sozinha basta

O `freeze` liga as quatro numa operação só. Separá-las produz estado
inconsistente — um ruleset ativo sem issue explicando, ou um `enabled: false`
que ninguém empurrou.

| camada | o que impede | quem pode contornar |
| --- | --- | --- |
| `enabled: false` no `anchors.yaml` | commit e push na máquina de quem já clonou | `--no-verify` |
| **o CLI** | todo comando que produz estado | — |
| **o ruleset** no remoto | push e merge no servidor | admin |
| **o claim** | atribuição de card novo | — |

### Por que o freio da plataforma não basta

Um ruleset que barra push e merge não alcança a máquina de quem **já clonou**.
Ali o `check`, o `judge` e o `ingest` continuam rodando e gravando no mapa.

Congelar significa **parar de produzir estado**, não só de entregá-lo.

### Por que os hooks locais não bastam

Um hook é contornável (`--no-verify`), e quem clona depois do congelamento não
passa por ele até instalar os hooks. Por isso ele é a **segunda** camada — a
primeira é o ruleset, que ninguém contorna de dentro.

Os hooks existem para que a pessoa **descubra cedo**, não para serem
invioláveis.

## O motivo é obrigatório

```sh
anchors freeze --motivo "..."     # sem isto, o comando recusa
```

Não é burocracia. É o texto que **toda recusa** vai mostrar — no commit
bloqueado, no push bloqueado, no claim recusado, na issue.

Um congelamento sem razão escrita é indistinguível de configuração quebrada, e
quem esbarra nele tenta contornar em vez de ler.

## Quem for consertar passa

Isso é deliberado, e está escrito em cada mensagem de recusa:

- `--no-verify` contorna os hooks
- o **admin** contorna o ruleset

Um freio que impede o próprio conserto vira o problema. O congelamento existe
para impedir trabalho **por inércia**, não a correção que o destrava.

## O que continua funcionando congelado

A régua: roda congelado o que **não produz estado do projeto**.

```sh
anchors status      # onde o projeto está
anchors doctor      # o raio-X do ecossistema
anchors guide ...   # os guias
anchors coverage    # o que já foi medido
anchors impact      # o que uma mudança atingiria
anchors thaw        # senão o congelamento seria irreversível
```

Quem está investigando o problema precisa desses. O que é recusado é o que
**escreve**: `check`, `map build`, `judge`, `ingest`, `new`, `escalate`.

## O que o freeze faz, passo a passo

1. escreve `enabled: false` e o motivo no `anchors.yaml`
2. **commita e empurra** — com `--no-verify`, porque os hooks que ele acabou de
   ativar recusariam o próprio congelamento
3. cria o ruleset `anchors-freeze` no remoto, barrando push e merge em todos os
   branches
4. abre a issue `[congelado]` com o motivo

O passo 2 usa `--no-verify` por uma razão a mais: o `pre-commit` roda os gates, e
um congelamento urgente **não pode depender de a suíte estar verde** — o motivo
do congelamento pode ser justamente que ela não está.

### As opções

```sh
anchors freeze --motivo "..." --sem-ruleset   # só o freio local
anchors freeze --motivo "..." --sem-push      # não commita nem empurra
```

## Depois do thaw

O `thaw` devolve o `anchors.yaml` ao estado exato de antes — comentários
incluídos.

Uma coisa a saber: os hooks guardam a resposta em cache por até **10 minutos**,
para não pagar um `git fetch` a cada commit. Quem tiver o cache quente pode
levar esse tempo para perceber a liberação, ou forçar com um `git fetch`.

O cache é **assimétrico** de propósito: ele guarda "não congelado", mas quando
diz "congelado" consulta a rede de novo. O motivo é o custo do erro — um cache
que mantivesse o congelamento deixaria a pessoa barrada com o projeto já
liberado, sem entender por quê e sem nada que pudesse fazer.

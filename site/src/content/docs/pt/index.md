---
title: Anchors
description: Um framework spec-first de continuidade para desenvolvimento assistido por IA.
template: splash
hero:
  tagline: Um framework spec-first — a spec vem antes do código, e continua sendo a verdade depois dele. Mantém um projeto coerente ao longo do tempo através de âncoras que guiam o trabalho na ida e o confrontam na volta, e não podem mentir sobre o código.
  actions:
    - text: Ler o conceito
      link: /docs/conceito/
      icon: right-arrow
      variant: primary
    - text: Ver no GitHub
      link: https://github.com/co2-lab/anchors
      icon: external
---

## Spec-first

A spec nasce antes do código e continua sendo a verdade depois dele — não é
documentação escrita a posteriori. Tudo mais (plano, mapa de dependências,
testes, gates de qualidade) deriva da spec e é confrontado contra ela. Veja o
pilar [Spec](/docs/spec/) para a disciplina completa.

## Três partes

**A doutrina** — o conceito, agnóstico de ferramenta: os 6 pilares e o mecanismo
comum.

**A operação** — como um dia de trabalho acontece, como configurar, e o que cada
gate mede.

**O CLI** — a ferramenta em Go que uma IA opera para exercitar o ciclo.

## Por onde começar

**Se você vai USAR o Anchors num projeto**, comece pela operação — a doutrina
faz mais sentido depois de você ter visto o ciclo rodar:

1. [**O CLI**](/docs/cli/) — instale, e rode os primeiros cinco minutos.
2. [**O fluxo de trabalho**](/docs/fluxo-de-trabalho/) — como um dia de trabalho
   acontece, do pedido de trabalho ao merge.
3. [**O anchors.yaml**](/docs/anchors-yaml/) — o que cada bloco de configuração
   decide.
4. [**Os gates**](/docs/gates/) — o catálogo, e o que ligar em cada tipo de
   projeto.

**Se você quer ENTENDER o Anchors**, comece pela doutrina:

1. [**Conceito**](/docs/conceito/) — a fundação: a âncora, a maturidade, o
   grafo, a sincronia, as issues. Todos os pilares o pressupõem.
2. **Os pilares**, na ordem da rota: [Estrutura](/docs/estrutura/) →
   [Planejamento](/docs/planejamento/) → [Spec](/docs/spec/) →
   [Rastreabilidade](/docs/rastreabilidade/) → [Propagação](/docs/propagacao/) →
   [Qualidade](/docs/qualidade/).

## Quando algo der errado

- [**Congelar o projeto**](/docs/congelar/) — o botão de pânico: como parar todo
  o trabalho quando um problema precisa ser resolvido antes.

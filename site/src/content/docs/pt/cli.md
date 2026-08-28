---
title: O CLI
---

CLI único do framework Anchors, em Go. É a ferramenta que uma IA opera para
exercitar o ciclo — a IA não precisa saber o Anchors de cor, ela pergunta ao
binário (`anchors guide`), aprende o fluxo, e opera com os comandos. O Anchors
**não embute IA**: ele é a ferramenta que a IA usa, em qualquer cliente
(Claude Code, GPT, Gemini…).

## Instalação

```sh
git clone https://github.com/co2-lab/anchors.git
cd anchors/cli
go build -o anchors ./cmd/anchors
./anchors --help
```

## Uso num projeto

```sh
anchors init            # configura o anchors.yaml (por P&R; sugere presets de stack)
anchors map build       # constrói o mapa de dependências a partir dos arquivos
anchors doctor          # saúde do ecossistema: órfãos, colisões, buracos de cobertura
anchors check --all     # roda os gates de qualidade; abre issues; carimba o mapa
anchors report all      # gera os relatórios em docs/anchors/
```

O CLI só lê **texto** — nunca parseia código. As anotações que ele entende
(identidade, `@noPropagation`…) vivem em comentários; os marcadores por
linguagem são configuráveis. Assim ele é agnóstico de stack.

## O fluxo, dirigido pela IA

A IA lê `anchors guide` e opera este ciclo — o watcher **enfileira** o
trabalho, a IA **puxa** da fila (a conversa nunca fica presa):

```
  planejar ──▶ especificar ──▶ mapear ──▶ implementar ──▶ testar ──▶ confrontar
     │             │             │            │             │            │
 guia de plano   .spec.md    map build   código+feature  testes    check / doctor
                                                                        │
                                          issue ◀── divergência que a IA não resolve
```

Cada arquivo salvo faz o watcher enfileirar a próxima task (spec→implementar,
feature→testar…). A IA não precisa lembrar o que vem depois; a fila diz.

## O que o CLI faz hoje

| área | comandos | o que entrega |
|---|---|---|
| **Estrutura** | `init` | configura o `anchors.yaml`; presets de estrutura para ~17 stacks |
| **Mapa** | `map build`, `map show`, `governs` | o grafo de dependências; quem rege quem |
| **Propagação** | `impact`, `stale` | a onda de uma mudança; o que ficou desatualizado |
| **Fila** | `watch`, `queue`, `next`, `done`, `drop`, `reclaim` | o watcher em background enfileira; a IA puxa |
| **Qualidade** | `check`, `judge`, `doctor` | gates determinísticos **e de julgamento por IA**; saúde sistêmica |
| **Identidade** | `code` | gera/valida um código de cenário único (evita colisões) |
| **Confiança** | `ingest`, `coverage` | ingere JUnit/lcov do runner; cobertura por **cenário**, do **diff**, e **delta** |
| **Relatórios** | `report` | 6 perspectivas em `docs/`: tests, quality, structure, config, issues, inconsistencies |
| **A ponte IA** | `guide` (+ `guide plan/spec/code/feature/test/guide`) | o playbook e as réguas embutidas que a IA lê para operar |

### O que dá confiança no entregável

O melhor indicador de um ciclo bem-sucedido é **não ter bugs no final**. Além
de exigir testes, o Anchors mede a *qualidade* deles a partir do artefato que
o runner já gera:

- **por cenário** — cada requisito da spec (`SPCR-V01`…) tem um teste que
  **passou**? (semântico, não cobertura de linha)
- **do diff** — as linhas que você **mudou** estão cobertas? (pega o bug na
  linha nova)
- **delta** — a cobertura **caiu** desde a última medição? (pega regressão)

E gates que um script não computa ("esta tela respeita a arquitetura?") viram
**gates de julgamento por IA**: a IA lê os *pontos de conformidade* do guide,
confronta o alvo item a item, e o veredito entra na mesma mecânica (carimbo +
issue) — envelhecendo se o alvo mudar.

## Estado do projeto

Em construção, e honesto sobre isso. A **doutrina** dos 6 pilares está escrita
e revisada; o **CLI** exercita o ciclo inteiro e foi validado contra uma
prova de conceito real — um app mobile com backend serverless, com um grafo
de porte não trivial.

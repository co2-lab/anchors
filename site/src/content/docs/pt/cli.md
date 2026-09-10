---
title: O CLI
description: A referência de todos os comandos, agrupados pelo que você está tentando fazer.
---

CLI único do framework Anchors, em Go. É a ferramenta que uma IA opera para
exercitar o ciclo — a IA não precisa saber o Anchors de cor, ela pergunta ao
binário (`anchors guide`), aprende o fluxo, e opera com os comandos.

O Anchors **não embute IA**: ele é a ferramenta que a IA usa, em qualquer
cliente.

O CLI só lê **texto** — nunca parseia código. As anotações que ele entende vivem
em comentários, e os marcadores por linguagem são configuráveis. É assim que ele
é agnóstico de stack.

## Instalação

```sh
brew install co2-lab/tap/anchors
```

Ou compilando:

```sh
git clone https://github.com/co2-lab/anchors.git
cd anchors/cli && go install ./cmd/anchors
anchors --help
```

## Os primeiros cinco minutos

```sh
anchors init            # configura o anchors.yaml, por perguntas e respostas
anchors install-hooks   # instala os hooks de pre-commit e pre-push
anchors map build       # constrói o mapa a partir dos arquivos
anchors status          # onde o projeto está, e o próximo passo
```

O `status` é o comando a rodar quando você não sabe o que fazer. Ele responde
onde o projeto está no ciclo e **qual é o próximo passo** — não é um relatório,
é uma instrução.

---

## Começar um projeto

| comando | o que faz |
| --- | --- |
| `anchors init` | gera o `anchors.yaml` por P&R, sugerindo preset conforme a stack |
| `anchors install-hooks` | instala o `pre-commit`, o `commit-msg` e o `pre-push` |
| `anchors map build` | constrói o `anchors.graph.yaml` a partir dos arquivos |

O `init` tem modo não-interativo, para script e CI:

```sh
anchors init --non-interactive                              # devolve as decisões em JSON
anchors init --non-interactive --artifacts=spec,test --colocation
```

---

## Saber onde você está

| comando | responde |
| --- | --- |
| `anchors status` | onde o projeto está no ciclo, e o próximo passo |
| `anchors doctor` | o raio-X: órfãos, colisões, sinais ausentes, buracos |
| `anchors coverage` | cobertura por cenário, por linha, e mutação |
| `anchors stale` | o que mudou e ainda não foi reconfrontado |
| `anchors impact <arquivo>` | o que uma alteração aqui atinge |
| `anchors governs` | quem cada guia rege, e quantos |
| `anchors compliance` | o estado de cada dever regulatório |

O `doctor` é o que dizer para alguém que herdou o projeto. Ele não mede
qualidade de código — mede se o **ecossistema** está saudável: spec sem código,
código sem spec, sinal que ninguém ingeriu, gate declarado que não protege nada.

---

## Escrever

| comando | o que faz |
| --- | --- |
| `anchors guide <artefato>` | a régua de como escrever (spec, code, feature, test, plan, review, work) |
| `anchors new <kind> <nome>` | emite o esqueleto conforme a régua |
| `anchors code` | gera um código de identidade único |
| `anchors recode <de> <para>` | renomeia um código e propaga por todo o projeto |
| `anchors work <etapa> --for <alvo>` | emite o prompt de trabalho de uma etapa |

Os guias são a documentação **executável**: em vez de a IA decorar o Anchors,
ela pergunta.

```sh
anchors guide work      # a régua de quem pegou um card
anchors guide spec      # como escrever uma spec
anchors guide review    # a régua de quem revisa um PR
```

---

## Confrontar

| comando | o que faz |
| --- | --- |
| `anchors check --changed <arq>` | roda os gates sobre o que mudou |
| `anchors check --all` | roda sobre o projeto inteiro |
| `anchors verify --phase <fase>` | roda TUDO o que a fase cobra (gates + ferramentas externas) |
| `anchors audit <arquivo>` | o dossiê de pendências de um arquivo, para correção em lote |
| `anchors judge <alvo> --gate <g>` | registra o veredito de um gate de julgamento |
| `anchors suggest` | lista, aplica e decide as correções propostas |

### `--changed` vs `--all`

O `--changed` entrega o **raio de impacto**: o arquivo e tudo que depende dele.
É o certo para quase todo gate — quem quebrou por tabela precisa ser
confrontado.

O `--all` é a foto do projeto inteiro, e é o que o CI roda.

### `--no-record`

Roda os gates **sem** gravar no mapa. É o que o CI usa: ele confronta, mas não
deve produzir um mapa diferente do que foi commitado.

---

## Os sinais de teste

| comando | o que faz |
| --- | --- |
| `anchors test` | roda as suítes declaradas e ingere os relatórios |
| `anchors mutation` | roda as suítes de mutação e ingere os relatórios |
| `anchors ingest --junit <x> --lcov <y>` | ingere relatórios que o projeto já gerou |

**Prefira `anchors test` a `anchors ingest`.** Os dois numa operação só é o que
garante que o mapa reflete o que **acabou de rodar** — ingerir à mão pode amarrar
ao mapa o resultado de uma execução anterior, e nada acusaria.

O Anchors não roda mutação nem conhece ferramenta: ele consome o formato aberto
**Mutation Testing Elements** (`schemaVersion 1.x`), que Stryker, PIT, Infection
e mutmut emitem.

---

## O trabalho (modo GitHub)

| comando | o que faz |
| --- | --- |
| `anchors escalate "..."` | abre a issue de uma mudança necessária no plano ou na spec |
| `anchors pr-body` | escreve as linhas que fecham os cards, na sintaxe da plataforma |
| `anchors deliver` | registra a entrega de uma etapa — o gatilho do review |
| `anchors task-status` | o relato da rodada: onde o card está, o veredito dos checks, e o que vem |

Ver [O fluxo de trabalho](/docs/fluxo-de-trabalho/) para como esses comandos se
encaixam num dia de trabalho.

### `escalate`: as duas saídas

```sh
anchors escalate "<o que está errado>" --sobre <arquivo> --card <n>
anchors escalate "<o que precisa mudar>" --sobre <arquivo> --para-usuario
```

A primeira abre um achado que se entrega junto com o card. A segunda abre uma
**decisão**, e o card para até ela sair.

A escolha entre as duas é sua, e o critério é um só: **isto muda a direção do
projeto?** Se muda, ou se você tem dúvida, é decisão de quem planejou.

---

## A fila (modo local)

| comando | o que faz |
| --- | --- |
| `anchors watch` | o watcher em background: vê mudanças e ENFILEIRA trabalho |
| `anchors queue` | lista as tasks vivas |
| `anchors next` | puxa e reivindica a próxima |
| `anchors done` | fecha task(s) reivindicada(s) |
| `anchors drop` | descarta uma task sem concluí-la |
| `anchors reclaim` | devolve à fila as tasks presas (worker morto) |

O watcher **enfileira**, a IA **puxa**. É o que impede a conversa de ficar presa
esperando o trabalho terminar.

---

## Parar tudo

| comando | o que faz |
| --- | --- |
| `anchors freeze --motivo "..."` | congela o projeto inteiro |
| `anchors thaw` | libera |

Ver [Congelar o projeto](/docs/congelar/).

---

## Relatórios

```sh
anchors report all       # gera os relatórios em docs/anchors/
```

Recortes do que o Anchors mede, por perspectiva — para quem precisa do estado
sem rodar comando.

---

## Códigos de saída

| código | significa |
| --- | --- |
| `0` | passou |
| `1` | um gate bloqueante reprovou, ou há julgamento pendente |
| `3` | **não regido**: nenhum arquivo casa uma camada do `layers:` |

O `3` existe para o hook distinguir "reprovou" de "não tenho jurisdição sobre
isto". Tratá-lo como falha impediria commitar mudança só de configuração — que é
trabalho legítimo que a Estrutura deliberadamente não rege.

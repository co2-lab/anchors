---
title: "Referência de Comandos do CLI"
description: "Índice completo e pesquisável de todos os comandos disponíveis no CLI do Anchors."
---

O **CLI do Anchors** fornece uma suíte de ferramentas integrada para gerenciar o ciclo de desenvolvimento, os gates de qualidade, o grafo de dependências e o trabalho de agentes de IA.

Abaixo está o índice completo de todos os comandos organizados por área funcional.

---

## 1. Esteira de Qualidade e Verificação

Comandos para rodar gates, validar código e diagnosticar a saúde do projeto:

| Comando | Uso | Descrição | Guia Detalhado |
| :--- | :--- | :--- | :--- |
| **`check`** | `anchors check [--all]` | Roda os gates de qualidade declarados contra arquivos modificados. | [Guia: check](/pt/docs/cli/commands/check/) |
| **`verify`** | `anchors verify` | Executa TUDO que a fase exige: gates do Anchors + ferramentas externas. | [Guia: check](/pt/docs/cli/commands/check/) |
| **`doctor`** | `anchors doctor [--verbose]` | Raio-X da saúde do projeto, pontas soltas, orfandade e maturidade. | [Guia: doctor](/pt/docs/cli/commands/doctor/) |
| **`audit`** | `anchors audit <arquivo>` | Dossiê detalhado de pendências e gates de um único arquivo/módulo. | [Guia: doctor](/pt/docs/cli/commands/doctor/) |
| **`test`** | `anchors test` | Roda as suítes de teste declaradas e amarra os sinais ao grafo. | [Guia: check](/pt/docs/cli/commands/check/) |
| **`mutation`** | `anchors mutation` | Executa motores de teste de mutação e afere o score de mutantes. | [Gate: mutation-score](/pt/docs/gates/mutation-score/) |
| **`coverage`** | `anchors coverage` | Exibe a cobertura de testes por cenário, linha e delta da diff. | [Gate: line-coverage](/pt/docs/gates/line-coverage/) |
| **`judge`** | `anchors judge <gate>` | Registra o veredito semântico de uma IA para gates de julgamento. | [Julgamento por IA](/pt/docs/concepts/ai-judgment/) |
| **`review`** | `anchors review <alvo> --gate <g> --by <quem>` | Registra um segundo olhar sobre um alvo de um gate revisado, com os achados; `--pending` lista o que está a revisar. | [Guia: review](/pt/docs/cli/commands/review/) |

---

## 2. Grafo e Rastreabilidade

Comandos que inspecionam, alteram e calculam arestas em `anchors.graph.yaml`:

| Comando | Uso | Descrição | Guia Detalhado |
| :--- | :--- | :--- | :--- |
| **`map`** | `anchors map <build\|show>` | Opera o grafo de dependências; inspeciona nós e arestas direcionadas. | [Guia: map](/pt/docs/cli/commands/map/) |
| **`impact`** | `anchors impact <alvo>` | Análise de impacto: calcula a onda de alcance se o arquivo mudar. | [Guia: map](/pt/docs/cli/commands/map/) |
| **`stale`** | `anchors stale` | Lista arestas obsoletas: dependências que mudaram sem reconfronto. | [Propagação](/pt/docs/concepts/propagation-and-impact/) |
| **`code`** | `anchors code <camada>` | Gera o próximo código de identidade único (ex: `AUTH-B03`). | [Rastreabilidade](/pt/docs/concepts/traceability-and-codes/) |
| **`recode`** | `anchors recode <velho> <novo>` | Renomeia um código atomicamente em specs, features, testes e código. | [Guia: map](/pt/docs/cli/commands/map/) |
| **`stamp`** | `anchors stamp <arquivo>` | Escreve carimbos `@contract` faltantes em dublês de teste (mocks). | [Gate: mock-stamped](/pt/docs/gates/mock-stamped/) |
| **`failures`** | `anchors failures` | Audita falhas observadas que ainda não foram explicadas em spec. | [Gate: failure-declared](/pt/docs/gates/failure-declared/) |
| **`compliance`**| `anchors compliance` | Estado dos deveres regulatórios: quantos nós sujeitos vs conformes. | [Gate: obligation-honored](/pt/docs/gates/obligation-honored/) |
| **`governs`** | `anchors governs` | Exibe quais arquivos cada guia de doutrina rege no repositório. | [Doutrina de Produto](/pt/docs/concepts/doctrine/) |

---

## 3. Fluxo e Fila de Trabalho para IAs

Comandos que sustentam o ciclo contínuo de pair programming assistido por IA:

| Comando | Uso | Descrição | Guia Detalhado |
| :--- | :--- | :--- | :--- |
| **`watch`** | `anchors watch` | Monitor de arquivos em background que detecta alterações e enfileira. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`queue`** | `anchors queue` | Lista tarefas ativas na fila aguardando execução por um worker. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`next`** | `anchors next` | Puxa e assume a próxima tarefa pendente na fila (worker call). | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`work`** | `anchors work <alvo>` | Emite o prompt de trabalho otimizado para o estágio correspondente. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`deliver`** | `anchors deliver` | Registra a entrega de uma etapa e dispara o gatilho de revisão. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`done`** | `anchors done <id>` | Conclui tarefas e move para o histórico de arquivamento. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`discard`** | `anchors discard <id>` | Tira o card do quadro sem deletar o registro histórico. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`drop`** | `anchors drop <id>` | Descarta imediatamente uma tarefa da fila sem arquivar. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`reclaim`** | `anchors reclaim` | Devolve à fila tarefas travadas por agentes/workers desconectados. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`suggest`** | `anchors suggest` | Lista, aplica ou descarta sugestões automatizadas de correção. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`escalate`** | `anchors escalate` | Abre issue de escalonamento para decisões que exigem humano. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`unblock`** | `anchors unblock` | Abre card de trabalho para destravar um item em `needs-user`. | [Guia: flow](/pt/docs/cli/commands/flow/) |
| **`decided`** | `anchors decided <card>` | Libera um card escalonado após a decisão humana ser tomada. | [Guia: flow](/pt/docs/cli/commands/flow/) |

---

## 4. Setup, Operação e Governança

Comandos de inicialização, configuração e congelamento do projeto:

| Comando | Uso | Descrição | Guia Detalhado |
| :--- | :--- | :--- | :--- |
| **`init`** | `anchors init` | Assistente interativo para configurar `anchors.yaml` e descobrir camadas. | [Guia: init](/pt/docs/cli/commands/init/) |
| **`new`** | `anchors new <tipo>` | Emite esqueletos de artefatos (spec, feature, teste) nos padrões. | [Guia: init](/pt/docs/cli/commands/init/) |
| **`freeze`** | `anchors freeze` | Congela o projeto: bloqueia trabalhos até `anchors thaw`. | [Guia: freeze](/pt/docs/cli/commands/freeze/) |
| **`thaw`** | `anchors thaw` | Descongela o projeto: retoma os fluxos de desenvolvimento. | [Guia: freeze](/pt/docs/cli/commands/freeze/) |
| **`install-hooks`** | `anchors install-hooks` | Instala hook `pre-commit` do git para rodar gates sobre o stage. | [Guia: init](/pt/docs/cli/commands/init/) |
| **`changelog`** | `anchors changelog` | Compila o changelog técnico a partir dos commits release a release. | [Pilares](/pt/docs/structure/) |
| **`commit-msg`** | `anchors commit-msg` | Valida a mensagem de commit no formato consumido pelo changelog. | [Gate: header-valid](/pt/docs/gates/header-valid/) |
| **`migrate`** | `anchors migrate` | Atualiza arquivos de configuração legados para a sintaxe moderna. | [O anchors.yaml](/pt/docs/anchors-yaml/) |
| **`settings`** | `anchors settings` | Configuração LOCAL do agente atual sem sujar o repositório. | [Guia: init](/pt/docs/cli/commands/init/) |
| **`guide`** | `anchors guide` | Imprime os guias autossuficientes destinados à ingestão por IAs. | [Conceitos](/pt/docs/concept/) |
| **`board`** | `anchors board [--live]`| Publica ou serve o quadro Kanban ao vivo com o estado do projeto. | [Guia: flow](/pt/docs/cli/commands/flow/) |

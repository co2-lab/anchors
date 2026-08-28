# Usando o Anchors

> Documentação de USO. O `README.md` e os 11 documentos de pilar (`SPEC.md`,
> `QUALITY.md`, `TRACEABILITY.md`…) explicam o CONCEITO — por que o framework existe e
> como ele pensa. Este arquivo explica o que você faz na prática: que arquivos o projeto
> precisa ter, o que cada um contém, e o que cada gate cobra.
>
> Ele nasceu de uma falha concreta: uma sessão de trabalho travou por não saber o formato
> de uma spec. A informação existia — em `anchors new spec --list-sections` e na mensagem
> de erro de um gate — mas nada apontava para lá, e o guide embutido chegava a instruir a
> descobrir errando ("rode `anchors check`; se reprovar, a mensagem diz o formato").
> Documentação atrás de subcomando só é achada por quem já sabe que ela existe.

## 1. Os arquivos que o projeto precisa ter

| arquivo | quem escreve | serve para |
| --- | --- | --- |
| `PROJECT.md` | o agente, entrevistando você (`anchors guide project`) | o resumo TÉCNICO do projeto: stack, paradigma, estrutura macro, extensões, indentação, ferramental. Lido antes de escrever cada arquivo |
| `INSIGHTS.md` | o mesmo, na mesma passada | a transcrição da entrevista: cada pergunta, cada resposta, e o que foi DESCARTADO e por quê |
| `anchors.yaml` | `anchors init` (e você refina) | declara as CAMADAS (que arquivo é o quê), os GATES ligados, e o dialeto do projeto |
| `anchors.graph.yaml` | `anchors map build` | o mapa: quem depende de quem. Gerado, versionado, nunca editado à mão |
| `guides/*.md` | você | a régua de cada artefato. O `init` semeia `HEADER_GUIDE.md` e `SPEC_GUIDE.md` |
| `<Unidade>.spec.md` | você (via `anchors new spec`) | o que a unidade faz, com cada regra catalogada e identificada |
| `<Unidade>.feature` | você (via `anchors new feature`) | os cenários de comportamento, cada um citando o código da regra |
| `<Unidade>.test.*` | você | a prova executável, citando o código no nome do caso |
| `.anchors/` | o CLI | fila de trabalho, espelho da última verificação, issues. Não versionado |

A TRINCA é `spec` + `feature` + `test` amarradas pelo mesmo código de identidade. É o que
permite perguntar "esta regra tem cenário? tem teste? o teste passou?" sem ler o código.

## 2. Começando num projeto

**Projeto que ainda NÃO existe** (diretório vazio): antes de tudo, há a fase DESCOBRIR.
O `anchors init` INFERE a Estrutura do que está no disco — num diretório vazio não há o
que inferir. Então o agente que opera o Anchors lê `anchors guide project` e entrevista
você em cinco etapas (propósito e forma → linguagem → arquitetura e paradigma → estrutura
macro e convenções → ferramental e formatação), uma pergunta por vez. No fim ele escreve
`PROJECT.md` (as decisões) e `INSIGHTS.md` (a transcrição, com o que foi descartado).
Só então o `init` tem o que perguntar — e você responde com o que está no PROJECT.md.

```sh
anchors init          # perguntas → anchors.yaml + guides semeados (precisa de terminal)
anchors map build     # constrói o mapa; TODO arquivo novo só existe para os gates depois disto
anchors check --all   # roda os gates sobre tudo
anchors doctor        # saúde sistêmica: órfãos, camadas frouxas, arestas mortas
```

Num projeto que JÁ EXISTE, não tente conformar tudo de uma vez. Ligue os gates de
ESTRUTURA primeiro (eles medem o código como está, sem exigir spec), deixe os de trinca
informativos, e promova conforme as specs nascem. Ligar tudo de uma vez produz centenas de
issues no primeiro dia e ninguém lê a lista.

## 3. Escrevendo uma spec

Não escreva do zero — o CLI emite o esqueleto já conforme:

```sh
anchors new spec <Nome> --out <caminho>/<Nome>.spec.md
anchors new spec --list-sections     # as 25 seções, com QUANDO usar cada uma
```

Uma regra é **catalogada** quando tem código E lugar estruturado. Menção solta em prosa
não conta:

```md
### LOGIX-B01 — descrição da regra        <- cabeçalho (preferido)
| `LOGIX-B02` | descrição |               <- linha de tabela
- **LOGIX-B03** descrição                 <- bullet-negrito
```

A LETRA diz a natureza da regra (`S` estado, `A` ação, `V` validação, `R` permissão, `X`
restrição…). O projeto declara as suas em `rule_types`; sem declaração valem as canônicas.
O COMPRIMENTO do código é do projeto (`code_lengths`) — 4, 5, ou os dois durante uma
migração.

## 4. Os gates

Um gate CONFRONTA dois artefatos e devolve um veredito. **Bloqueante** barra a promoção;
**informativo** registra e não barra (é a maturação: liga-se informativo, promove-se depois).

Quatro vereditos, e a diferença importa:

- **✓ passou** — confrontou e está conforme.
- **✗ reprovou** — vira issue. Bloqueante: barra.
- **⚠ divergiu** — olhou, algo não bate, mas não barra.
- **~ indeterminado** — o gate rodou e NÃO TEVE O QUE MEDIR. Não é aprovação: é
  cobertura que não existe. Um projeto sem specs tem quase tudo aqui.

### Referência

| `cenario-identidade` | dois cenários da MESMA feature não podem ter o mesmo código |
| `cenario-letra-declarada` | a letra do código de cenário existe no vocabulário |
| `cenario-tipo-alinhado` | a tag de natureza do cenário concorda com a letra do código |
| `code-reference-valid` | um código citado por uma spec precisa EXISTIR |
| `codigo-catalogado` | o que o código EXPORTA precisa estar na spec — ou dispensado nele |
| `contract-status-declared` | a tabela "Contrato de Saída" de uma spec de INTERFACE |
| `count-honored` | um NÚMERO que a spec afirma sobre o código precisa bater com o código |
| `coverage-delta` | a cobertura de linha deste nó de código CAIU desde a ingestão |
| `domain-declared` | a seção `## Domínio` declara o que a unidade ACEITA — e quem garante |
| `evidence-fresh` | o placar deste teste vale contra o código de HOJE |
| `guide-has-checklist` | um GUIDE de governança deve destilar suas regras em PONTOS DE |
| `has-code` | o arquivo carrega um código de cenário (a identidade). Sem código, a |
| `header-conforme` | o arquivo tem o BLOCO DE CABEÇALHO do Anchors (`@anchors`) com o |
| `layer-boundary` | uma camada não alcança o que não é dela |
| `line-coverage` | a cobertura de linha do nó de código está >= 70%? (limiar fixo por |
| `mock-carimbado` | o dublê carrega a marca do trecho que ele substitui, e o gate |
| `mock-tipado` | todo dublê de teste precisa DERIVAR do módulo que substitui |
| `mutation-score` | dos mutantes que rodaram neste arquivo, o teste matou o suficiente? |
| `non-empty` | o arquivo não é vazio nem só espaço. Trivial, mas pega placeholders |
| `obligation-honored` | confronta as OBRIGAÇÕES TRANSVERSAIS declaradas (`obligations:` |
| `open-questions-resolved` | uma spec com pergunta em aberto NÃO está pronta para implementar |
| `pagination-honored` | uma função que promete devolver um CONJUNTO ("liste todos", "os |
| `placeholder-preenchido` | o esqueleto que o `anchors new` emite tem de ser PREENCHIDO |
| `plan-seeds-valid` | um PLANO semeia artefatos (as specs que vão nascer). Este gate |
| `prova-cruza-fronteira` | quando uma REGRA afirma relação com outra unidade — "espelha |
| `ref-resolves` | o `ref:` de um artefato aponta para a spec que REALMENTE o descreve |
| `region-pair-honored` | toda região aberta fecha, e fecha com o próprio código |
| `regra-implementada` | a spec cataloga regras; o código tem de mostrar que as realizou |
| `route-declared` | uma spec de TELA (layer: screen) declara COMO se chega até ela — a |
| `route-exists` | a rota que a spec DECLARA tem de existir no app |
| `scenario-asserts` | o passo de RESULTADO de um cenário precisa afirmar um resultado |
| `scenario-coverage` | cada código de cenário que a spec DECLARA tem um teste que |
| `sibling-guard` | funções IRMÃS (exportadas do mesmo módulo, recebendo o mesmo |
| `spec-feature-match` | todo REQUISITO declarado na spec tem ao menos um CENÁRIO na feature |
| `spec-sections` | uma spec deve CATALOGAR ao menos uma regra/estado com código — o |
| `teste-rastreavel` | um teste ligado a uma feature precisa DIZER o que prova |
| `tests-pass` | o nó de teste tem 0 falhas (do resultado de execução ingerido)? |
| `trigger-declared` | um gatilho de obrigação CITADO tem de existir |
| `trinca-completa` | uma spec de camada REGIDA precisa das três peças que a realizam — |
| `updated-at-atual` | o `updated_at` do header bate com a data do ÚLTIMO COMMIT que |
| `vr-baseline` | o cenário de regressão VISUAL prometido tem imagem de referência |

Gates externos (`run:`) chamam a ferramenta do seu stack (`tsc`, `eslint`, `go vet`) — o
Anchors não reimplementa o que a ferramenta já faz, só lê o veredito. O `scope:` decide se
o comando roda por arquivo (`node`), com os arquivos como argumento (`batch`) ou uma vez
sobre o projeto (`project`).

## 5. O ciclo de trabalho

```
(descobrir) → plano → spec → código+feature → teste → confronto
```

O `descobrir` entre parênteses roda UMA vez, só em projeto novo: é o que fixa a linguagem,
o paradigma e as convenções antes de o primeiro plano decidir quais specs nascem.

Cada passo produz o artefato que o próximo confronta. O `anchors watch` enfileira o
próximo passo a cada arquivo salvo; `anchors next` puxa da fila e `anchors done` fecha.
Com `workflow.mode: github`, a fila são as issues do repositório (ver `WORKFLOW.md`).

## 6. Quando algo não bate

| sintoma | o que é |
| --- | --- |
| `não é regido pela Estrutura` | o arquivo não casa nenhuma camada do `layers:` — declare a camada ou é arquivo que não deveria estar ali |
| gate `~` em tudo | falta o artefato do outro lado da relação (spec sem feature, feature sem teste) |
| `X gate(s) declarado(s) sem nada para medir` | o `on:` do gate não casa nenhum nó — crie o artefato ou remova o gate |
| `evidência VENCIDA` | o teste passou contra código que já mudou. Rode de novo |
| mapa DESATUALIZADO | `anchors map build` — os gates leem o mapa, não o disco |

Nunca "conserte" fazendo a âncora mentir sobre o código. Se o código diverge da spec, ou o
código está errado (corrija), ou a spec envelheceu (atualize — e isso pode gerar trabalho
novo, que é informação, não obstáculo).

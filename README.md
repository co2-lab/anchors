<div align="center">

# ⚓ Anchors

**Um framework spec-first de continuidade para desenvolvimento assistido por IA.**

Mantém um projeto coerente ao longo do tempo — através de âncoras: documentos que
guiam o trabalho na ida e o confrontam na volta, e não podem mentir sobre o código.

A spec vem antes do código, e continua sendo a verdade depois dele — não é
documentação escrita a posteriori. Tudo mais (plano, mapa, testes, gates) deriva
dela e é confrontado contra ela.

</div>

---

## O problema

Uma sessão de IA é sem memória entre invocações. A sessão acaba, o contexto some, e
a próxima — outro agente, outro dia, outro humano — recomeça sem as regras, sem
saber o que já foi decidido, sem saber onde o trabalho parou. O projeto cresce em
código, mas a **direção** não persiste. E os testes ficam verdes enquanto requisitos
inteiros seguem sem prova.

O Anchors ataca isso projetando os artefatos certos no repositório — as **âncoras** —
e dando a uma IA uma ferramenta para operá-las: manter o mapa de dependências,
propagar mudanças, medir a qualidade, e registrar o que diverge.

O nome vem da escalada: para subir uma parede você crava âncoras, pontos fixos que
dizem **para onde ir**, **seguram a corda para você não cair** e **demarcam por onde
passou**. No Anchors, âncoras são os documentos que nos acompanham na escalada do
projeto — specs, planos, features, guides, docs.

---

## Duas partes

| | | |
|---|---|---|
| **A doutrina** | o *conceito*, agnóstico de ferramenta — os 6 pilares e o mecanismo comum | os `*.md` desta raiz |
| **O CLI** | a *ferramenta* em Go que uma IA opera para exercitar o ciclo | [`cli/`](./cli) |

A IA não precisa saber o Anchors de cor — ela pergunta ao binário (`anchors guide`),
aprende o fluxo, e opera com os comandos. O Anchors **não embute IA**: ele é a
ferramenta que a IA usa, em qualquer cliente (Claude Code, GPT, Gemini…).

---

## Começando (o CLI)

```sh
cd cli
go build -o anchors ./cmd/anchors
./anchors --help
```

Num projeto qualquer:

```sh
anchors init            # configura o anchors.yaml (por P&R; sugere presets de stack)
anchors map build       # constrói o mapa de dependências a partir dos arquivos
anchors doctor          # saúde do ecossistema: órfãos, colisões, buracos de cobertura
anchors check --all     # roda os gates de qualidade; abre issues; carimba o mapa
anchors report all      # gera os relatórios em docs/anchors/
```

O CLI só lê **texto** — nunca parseia código. As anotações que ele entende
(identidade, `@noPropagation`…) vivem em comentários; os marcadores por linguagem são
configuráveis. Assim ele é agnóstico de stack.

---

## O fluxo, dirigido pela IA

A IA lê `anchors guide` e opera este ciclo — o watcher **enfileira** o trabalho, a IA
**puxa** da fila (a conversa nunca fica presa):

```
  planejar ──▶ especificar ──▶ mapear ──▶ implementar ──▶ testar ──▶ confrontar
     │             │             │            │             │            │
 guia de plano   .spec.md    map build   código+feature  testes    check / doctor
                                                                        │
                                          issue ◀── divergência que a IA não resolve
```

Cada arquivo salvo faz o watcher enfileirar a próxima task (spec→implementar,
feature→testar…). A IA não precisa lembrar o que vem depois; a fila diz.

---

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

Detalhes e decisões de arquitetura em [`cli/README.md`](./cli/README.md) e
[`cli/DECISIONS.md`](./cli/DECISIONS.md).

### O que dá confiança no entregável

O melhor indicador de um ciclo bem-sucedido é **não ter bugs no final**. Além de exigir
testes, o Anchors mede a *qualidade* deles a partir do artefato que o runner já gera:

- **por cenário** — cada requisito da spec (`SPCRX-V01`…) tem um teste que **passou**?
  (semântico, não cobertura de linha)
- **do diff** — as linhas que você **mudou** estão cobertas? (pega o bug na linha nova)
- **delta** — a cobertura **caiu** desde a última medição? (pega regressão)

E gates que um script não computa ("esta tela respeita a arquitetura?") viram **gates
de julgamento por IA**: a IA lê os *pontos de conformidade* do guide, confronta o alvo
item a item, e o veredito entra na mesma mecânica (carimbo + issue) — envelhecendo se o
alvo mudar.

---

## Os pilares

Na ordem da rota — da origem ao acabamento. Cada um responde a uma pergunta que uma
sessão futura faz ao chegar no projeto.

| # | Pilar | Pergunta que responde | Documento |
|---|---|---|---|
| 1 | **Estrutura de Projeto** | quais camadas existem e como se organizam? | [`STRUCTURE.md`](./STRUCTURE.md) |
| 2 | **Planejamento** | para onde vamos, em que ordem, onde paramos? | [`PLANNING.md`](./PLANNING.md) |
| 3 | **Spec** | o que esta coisa deve ser? | [`SPEC.md`](./SPEC.md) |
| 4 | **Rastreabilidade** | como as peças se conectam? | [`TRACEABILITY.md`](./TRACEABILITY.md) |
| 5 | **Propagação** | quando algo muda, como a mudança percorre tudo até fechar? | [`PROPAGATION.md`](./PROPAGATION.md) |
| 6 | **Qualidade** | isto está bom o suficiente? | [`QUALITY.md`](./QUALITY.md) |

O **fluxo de trabalho** que os interliga **não é um pilar** — é um desenho, que pode
mudar ou ter várias perspectivas. Os componentes que operam dentro dele é que são os
pilares.

O pilar **Spec** é o que faz o Anchors ser **spec-first**: nada se implementa sem
uma spec que diga o que construir primeiro, e ela continua sendo a verdade contra a
qual o código é confrontado — nunca documentação escrita depois do fato.

### As duas ideias que sustentam tudo

- **Maturidade é o topo.** Um projeto é maduro quando tem todos os seus pilares
  implementados e vigorosos — não é um número num artefato, é a presença e o vigor dos
  pilares no todo.
- **Âncoras não podem mentir.** Cada âncora é confrontada contra a realidade. Quando
  diverge, ou o trabalho está errado (corrige) ou a âncora ficou para trás (atualiza).
  O framework **detecta e apresenta** a divergência — nunca arbitra sozinho. Vale a
  cada mudança, de forma **incremental**: só se reconfere o mínimo que a mudança tocou.

---

## O mecanismo comum

Todos os pilares compartilham um mecanismo, definido em [`CONCEPT.md`](./CONCEPT.md):

- **Âncora** — documento que aponta, segura a corda e/ou demarca. As **estruturais**
  regem para dentro; as **de consumo** (como docs) são geradas para fora — âncoras, mas
  não pilares.
- **Grafo** — as âncoras e suas dependências formam um grafo material, versionado, de
  arestas tipadas (muitos-para-muitos). É o `anchors.graph.yaml`.
- **Sincronia** — cada aresta sabe se está em dia (carimbo por relação); a Propagação
  usa isso para espalhar a invalidação de forma incremental.
- **Issue** — quando um confronto falha e o Anchors não pode resolver sozinho, registra
  uma issue: arquivo em pasta-estado (`todo`/`doing`/`done`), imutável, aberta só por
  confronto, que nunca reabre (recorrência é issue nova).
- **Opt-out honesto** — dispensar uma exigência de forma explícita, registrada e datada
  (`@no-test`, `--no-block`, maturação). Dispensa a exigência, nunca o registro.
- **Dois planos** — o **vivo** (âncoras + grafo, a verdade atual) e o **histórico**
  (issues resolvidas, laudos, datados e imutáveis). Nunca se contradizem.

---

## Por onde começar a ler

**Vai USAR o Anchors num projeto? Comece por [`USING.md`](./USING.md)** — que arquivos o
projeto precisa ter, como escrever uma spec, o que cada gate cobra e o que fazer quando
algo não bate. Os documentos abaixo explicam o CONCEITO; o `USING.md` explica a prática.

1. **[`CONCEPT.md`](./CONCEPT.md)** — a fundação: a âncora, a maturidade, o grafo, a
   sincronia, as issues. Todos os pilares o pressupõem.
2. **Os pilares**, na ordem da rota: [`STRUCTURE`](./STRUCTURE.md) →
   [`PLANNING`](./PLANNING.md) → [`SPEC`](./SPEC.md) →
   [`TRACEABILITY`](./TRACEABILITY.md) → [`PROPAGATION`](./PROPAGATION.md) →
   [`QUALITY`](./QUALITY.md).
3. **O CLI** — [`cli/README.md`](./cli/README.md) para o estado da ferramenta,
   comando a comando, com o que já foi validado contra um projeto real.

Cada pilar é auto-contido, mas aponta para os outros onde se tocam.

---

## Estado do projeto

Em construção, e honesto sobre isso. A **doutrina** dos 6 pilares está escrita e
revisada; o **CLI** exercita o ciclo inteiro e foi validado contra uma prova de
conceito real, não-trivial — um app mobile + backend serverless com um grafo de
dependências de porte substancial. Os comandos acima existem e rodam; o roadmap
e os furos conhecidos estão em [`cli/README.md`](./cli/README.md).

O Anchors é o **conceito** — independente de qualquer ferramenta. Duas instâncias reais
o exercitam: uma ferramenta/IDE que o aplica e o workspace de prova de conceito real
(incompleto e imaturo de propósito — prova o mecanismo, não a completude). Elas ilustram
o conceito; não o definem.

> Nada aqui está cravado na pedra. Estes documentos são o conceito saindo da cabeça e
> encontrando seu lugar — refinados conforme o framework amadurece.

### Furo conhecido: o idioma

O Anchors é agnóstico de **linguagem** (`dialect:`), de **stack** (presets), de
**vocabulário** (`rule_types`) e de **jurisdição** — mas ainda não de **idioma**. Medido:
**752 strings em português contra 11 em inglês** no CLI. Tudo que o usuário lê — mensagem
de gate, prompt do `work`, guide embutido — fala a língua de quem escreveu o framework.

O que já é agnóstico: o Gherkin da feature (`dialect.gherkin_language`, 10 idiomas, default
`en`), as letras de `rule_types` (iniciais de termos em inglês: State, Rule, Constraint…)
e todo o conteúdo escrito pelo projeto.

O que falta: `language:` no `anchors.yaml` (default `en`) governando CLI, guides embutidos
e as seções do template de spec — que hoje saem **mescladas** (`## Actions` convivendo com
`## Efeitos`), sinal de que ninguém está decidindo o idioma, ele só acontece.

Está registrado como pendência deliberada, não como descuido. E não é tradução mecânica: as
mensagens dos gates explicam **por que** cada gate existe, não só o que falhou — dois
revisores independentes as chamaram do melhor material de documentação do framework.
Traduzir mal destruiria exatamente isso.

### Índice completo dos documentos

| Documento | Conteúdo |
|---|---|
| [`CONCEPT.md`](./CONCEPT.md) | O mecanismo comum: âncora, maturidade, grafo, sincronia, issues, vivo vs. histórico |
| [`STRUCTURE.md`](./STRUCTURE.md) | Pilar — a planta da casa; as camadas e sua ordem; o gabarito |
| [`PLANNING.md`](./PLANNING.md) | Pilar — a origem do movimento; semeia specs; o norte entre sessões |
| [`SPEC.md`](./SPEC.md) | Pilar — a origem da verdade; disciplina spec-first; guide→template→spec; regimes |
| [`SPEC_TYPES.md`](./SPEC_TYPES.md) | Catálogo (apoio ao pilar Spec) — tipos de spec por família |
| [`TRACEABILITY.md`](./TRACEABILITY.md) | Pilar — identidade contínua + mapa de dependências; órfãos |
| [`PROPAGATION.md`](./PROPAGATION.md) | Pilar — a onda incremental; staleness; quiescência |
| [`QUALITY.md`](./QUALITY.md) | Pilar — gates que medem; features → testes; maturação informativo → bloqueante |
| [`cli/`](./cli) | A ferramenta em Go: comandos, arquitetura, roadmap |
| [`simulation/`](./simulation) | A simulação Larder — o ciclo de vida exercitado numa app fictícia |

---

<div align="center">
<sub>O conceito é independente de qualquer ferramenta. Ele permanece verdadeiro sem elas.</sub>
</div>

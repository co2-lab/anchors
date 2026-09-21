# Os fluxos que o Anchors já tem — levantamento

> Levantamento anterior à escrita. O que o Anchors executa hoje, onde há DECISÃO, e o que
> vira fluxo. O fluxo de destravar de um projeto acompanhado NÃO entra aqui: é dele, e o Anchors
> apenas oferece o mecanismo para o projeto escrever o seu.

## O critério, corrigido

A primeira versão deste levantamento separou os comandos em "vira fluxo" e "não vira".
A pergunta estava errada.

**Nem todo comando vira um fluxo, mas todo comando FAZ PARTE de um.** Um comando pode ser:

- o GATILHO de um fluxo — `init` é o caso claro: não há decisão antes dele, e as decisões
  nascem dele (o próprio `--help` diz que "as questões cobrem apenas as decisões humanas");
- um PASSO dentro de um estado — `map build` no ciclo do worker;
- o que provoca a TRANSIÇÃO — `check` é o que faz CONFRONTADA virar FECHADA ou voltar;
- uma CONSULTA que informa a decisão sem mudar o estado — `impact`, `coverage`, `doctor`.

A diferença importa porque um comando que não é gatilho de nada, nem passo de nada, é um
comando órfão do processo — e isso é achado, não organização.

## Os fluxos

### 1. `adocao` — do diretório ao primeiro plano

Gatilho: **`anchors init`**. Antes dele não há decisão; a partir dele há uma, e o guia a
nomeia: *"DISCOVER — só se o projeto NÃO EXISTE ainda... pule esta etapa num projeto que
já tem código"*.

    (vazio)                    init
      └─ projeto NOVO      →   DESCOBERTA (entrevista em 5 etapas, PROJECT.md + INSIGHTS.md)
      └─ projeto EXISTENTE →   ESTRUTURA INFERIDA (o init lê do disco)
                           →   MAPA E WATCHER NO AR (map build + watch start)
                           →   PRONTO PARA PLANEJAR

O ramo é a decisão, e hoje ela é uma frase que alguém tem de ler e aplicar. Como fluxo, a
entrevista de 5 etapas é o único caminho quando o diretório está vazio.

Os comandos: `init` (gatilho), `guide project` (o passo da descoberta), `map build` e
`watch start` (os passos do último estado).

### 2. `plan` — do escopo ao arquivo

Gatilho: a conversa com o usuário. O único fluxo com REGRA DE OURO explícita: aprovar o
escopo ANTES de escrever o arquivo, porque salvar o plano já dispara o watcher.

    RASCUNHO (na conversa) → APROVADO → SEMEADO

O guia diz, em maiúsculas, *"não escreva o arquivo para mostrar como ficou e pergunte
depois"*. Como fluxo, RASCUNHO não tem saída para SEMEADO — só passando por APROVADO.

Os comandos: `guide plan` (passo), `new plan` (a transição para SEMEADO).

### 3. `worker` — o ciclo de quem executa uma tarefa da fila

O mais percorrido: toda tarefa passa por ele. Gatilho: **`anchors next`**.

Hoje é a lista (a)→(h) do guia, com regras em prosa ao redor. Os estados saem dela quase
diretos, e as decisões que hoje são nota de rodapé viram saídas:

    PUXADA (next) → MAPEADA (map build) → EXECUTADA (o artefato escrito)
                  → CONFRONTADA (check)
                       ├─ bloqueante vermelho  → volta a EXECUTANDO
                       ├─ julgamento pendente  → entra no fluxo `judge`
                       └─ verde                → FECHADA (done)

Duas regras do guia deixam de depender de memória:

- *"NUNCA feche a tarefa com bloqueante vermelho"* vira a AUSÊNCIA de saída de CONFRONTADA
  para FECHADA;
- *"impact e check leem o MAPA, não o disco"* vira o estado MAPEADA, que hoje é o passo (b)
  com uma nota de rodapé explicando por que ele existe.

Os comandos: `next` (gatilho), `map build`, `impact`, `ingest`, `coverage` (passos),
`check` (a transição), `done` (o fechamento).

### 4. `judge` — o veredito que nenhum script computa

Gatilho: **`anchors check`** marcando um alvo como pendente de julgamento. É subfluxo do
worker, e por isso tem gatilho de dentro.

    PENDENTE (judge --pending) → GUIA LIDO → JULGADO
                                   ├─ fail → issue aberta com o relatório
                                   └─ pass → issue anterior resolvida

O guia adverte *"NÃO julgue pela prosa inteira no olho — julgue contra CADA ponto"*. Como
fluxo, ler o guia é o único caminho para JULGADO.

Os comandos: `judge --pending` (gatilho), `guide <o que rege>` (passo), `judge --verdict`
(a transição).

## Os comandos e onde cada um entra

| comando | papel | fluxo |
| --- | --- | --- |
| `init` | gatilho | adocao |
| `guide project` | passo | adocao |
| `watch start` | passo | adocao |
| `new plan` | transição | plan |
| `next` | gatilho | worker |
| `map build` | passo | adocao, worker |
| `impact`, `ingest`, `coverage` | passo | worker |
| `check` | transição | worker (e gatilho de judge) |
| `done` | transição | worker |
| `judge --pending` | gatilho | judge |
| `judge --verdict` | transição | judge |

**Os que ainda não têm fluxo** — e isso é achado, não sobra:

`escalate`, `decided`, `unblock`, `synthesize`, `discard`, `freeze`, `thaw`, `reclaim`,
`drop`. Todos tratam do trabalho que TRAVA ou PARA, e a sequência entre eles não está
escrita em lugar nenhum — nem no guia. É o quinto fluxo, e provavelmente o que mais se
beneficia, porque travar é justamente quando ninguém lembra do procedimento.

Os de consulta e visualização (`doctor`, `governs`, `compliance`, `board`, `report`,
`task-status`, `audit`) informam decisões sem mudar estado — participam sem serem passo.

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |
| `Q01` | Os fluxos nascem juntos, ou o `worker` primeiro? Ele é o mais percorrido — e escrevê-lo sozinho exercita o formato enquanto mudá-lo ainda é barato. | você | a ordem de implementação |
| `Q03` | O quinto fluxo (o do trabalho que TRAVA: `escalate`, `decided`, `unblock`, `synthesize`, `discard`) entra nesta rodada? A sequência entre esses comandos não existe escrita em lugar nenhum, então escrevê-la é desenhar, não transcrever. | você | se o levantamento vira quatro fluxos ou cinco |
| `Q02` | O fluxo semeado pelo `init` é FIXO (o Anchors traz o seu e o projeto sobrescreve) ou é COPIADO para `flows/` na adoção (o projeto passa a ser dono desde o dia um)? A primeira mantém o default atualizável pelo Anchors; a segunda deixa o projeto livre sem precisar sobrescrever nada. | você | como `flows/` nasce num projeto novo |

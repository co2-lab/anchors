# Anchors — Catálogo de Tipos de Spec

> Catálogo de referência que acompanha o pilar de [`SPEC.md`](./SPEC.md). Enquanto o
> SPEC.md define o **modelo** (guide → template → spec; esqueleto comum + bloco de
> variação; regime comportamental/declarativo), este catálogo levanta os **tipos
> concretos** de spec por família de projeto e o que o bloco de variação de cada um
> contém.
>
> É um **esboço vivo**, não uma lista fechada. A taxonomia de tipos vem da Estrutura
> de Projeto de cada projeto (`STRUCTURE.md`) — este catálogo reúne os padrões
> principais das linguagens/stacks para *padronizar*, não para catalogar tudo.
> Destilado de um levantamento sobre frontend, backend e infra/CLI/lib, ancorado em
> exemplos reais: um app mobile com backend serverless, e uma ferramenta desktop
> multi-stack (Go no backend, React na interface).

---

## 1. O esqueleto comum (todo tipo carrega)

Independente do tipo, toda spec tem este esqueleto. O que muda entre tipos é só o
**vocabulário** que preenche entrada/saída/contrato, e o **bloco de variação** (§3).

| bloco comum | papel |
|---|---|
| **Identidade** | código estável + arquivo + tipo/camada + data. A nascente da rastreabilidade. |
| **Propósito / responsabilidade** | uma frase: o que faz e o que **não** faz (o limite da camada). |
| **Contrato de entrada** | o que recebe. |
| **Contrato de saída / efeitos** | o que retorna ou muda no mundo. |
| **Regras / comportamento** | pré/pós-condições, invariantes. |
| **Erros / falhas** | como e por que falha; quem sinaliza. |
| **Dependências** | de que camadas depende — e quais é **proibido** chamar (limite de camada). |
| **Rastreabilidade** | os códigos que descem para feature/teste. |

---

## 2. Os regimes (ortogonal ao tipo)

Cada tipo tende a um regime, mas o eixo é ortogonal — um artefato pode ter partes
dos dois (ver [`SPEC.md`](./SPEC.md) §6).

- **Comportamental** — descreve ação ("dado X, acontece Y"); propaga para
  features → testes (Gherkin).
- **Declarativo** — descreve estado desejado ("estes recursos existirão");
  propaga para conformidade de forma (diff / drift / validação estrutural).

---

## 2.5 A INTERFACE: o ponto de entrada (um conceito, muitos dialetos)

Antes da tabela por família, um conceito que **atravessa todas elas** e é a nascente de
quase todo processo: a **interface**.

**Interface é o ponto onde um trigger — interno ou externo — dá entrada no processo para
interação.** É onde o comportamento começa. Toda aplicação tem uma, e o *papel* é sempre
o mesmo; só o **nome e o dialeto** mudam conforme o tipo de aplicação:

| Aplicação | A interface é o(a)… | Quem/o que dispara |
|---|---|---|
| CLI | **comando** | o usuário ou outro programa |
| Web / Desktop / Mobile | **tela** (screen/page) | o usuário |
| API HTTP | **rota** (handler) | uma requisição externa |
| Job / Worker / serverless de evento | **trigger** (cron/evento) | um agendador ou evento (interno/externo) |

`screen`, `handler`, `command`, `trigger` **não são tipos independentes — são
DIALETOS de um mesmo tipo conceitual: a interface.** Confundir o dialeto com o conceito
é o erro que faz um projeto tratar "tela" e "handler HTTP" como coisas sem parentesco,
quando são a mesma coisa arquitetural: a porta de entrada do processo.

### Por que unificar importa

- **Mesma régua de conteúdo.** Toda interface, em qualquer dialeto, tem a mesma
  estrutura de spec: **contrato de entrada** (o que chega — props+navegação numa tela,
  path/query/body/headers numa rota, args+flags num comando, payload do evento num
  trigger), **contrato de saída** (o que volta ou muda), **estados** (loading/erro/dado,
  ou os status codes), **erros** (e seu mapeamento), e **auth/acesso**. O vocabulário
  concreto muda; os *aspectos regidos* não. Por isso os guides dialetais
  (`SCREEN_GUIDE`, `HANDLER_GUIDE`, …) são **irmãos**: mesmo conteúdo, dialetos
  diferentes — cada projeto tem o guide do dialeto que usa, mas todos descrevem uma
  interface.
- **Mesma posição na propagação.** A interface é a **origem** da onda de dados: ela
  declara o que consome e a propagação desce dela para as camadas de dados
  (interface → usecase/business-logic → repository/service). `STRUCTURE.md` §4 fala em
  "spec de interface demanda alteração numa spec de usecase" — é este conceito, não um
  tipo backend-específico. A Tabela de Dependências (§5) da interface é o que liga essa
  cadeia, seja a interface uma tela ou uma rota.
- **Mesmo regime.** Comportamental: "dado a entrada X, o processo faz Y e responde Z" —
  propaga naturalmente para feature → teste, em qualquer dialeto.

### Como um projeto materializa

A **Estrutura de Projeto** (`STRUCTURE.md`) declara qual(is) dialeto(s) de interface o
projeto tem, com o nome que o time usa. Um projeto que é mobile **e** backend
serverless, por exemplo, tem dois: `screen` (a interface do app) e
`handler`/`auth-trigger`/`job-trigger` (as interfaces do backend — rota HTTP, evento
de autenticação, cron/evento). São **tags distintas no
grafo** (o dialeto tem seu nome), mas **regidas pela mesma régua** (a de interface). O
guide de cada uma pode existir separado (um `HANDLER_GUIDE` além do `SCREEN_GUIDE`), e
seu conteúdo espelha o do irmão porque ambos descrevem uma interface — só troca o
dialeto do contrato de entrada/saída.

> Regra prática: ao encontrar um "ponto de entrada" num projeto novo — uma tela, uma
> rota, um comando, um trigger de job — trate-o como uma **interface**. Dê-lhe o nome do
> dialeto do projeto, mas escreva a spec pela régua de interface (entrada, saída,
> estados, erros, auth), e ligue a Tabela de Dependências para a propagação descer dela.

---

## 3. Tipos por família

Para cada tipo: a **fonte da variação** (o que o distingue), o **bloco de variação**
(as seções extras sobre o esqueleto comum) e o **regime** predominante.

### Frontend (React / React Native — ancorado num app mobile real)

| tipo | fonte da variação | bloco de variação | regime |
|---|---|---|---|
| **tela / page** _(interface)_ | estado + navegação | states (loading/error/data), data contract, data states, navegação, mensagens | comportamental |
| **componente** | props + variantes | props API, variantes (por eixo), estados visuais, eventos/callbacks, slots | comportamental |
| **hook** | efeitos | signature (in/out), efeitos (queries/mutations, queryKeys, invalidations), estratégia de erro | comportamental |
| **store (estado global)** | estado + ações | shape do estado, actions, selectors, hidratação/persistência, invariantes | comportamental |
| **service** | operação remota | operação/endpoint, input/output, efeitos no servidor, erros (throw) | comportamental |
| **repository** | acesso a dado | model acessado, queries/mutations (CRUD), shape do dado, erros (throw) | comportamental |
| **lib / helper puro** | transformação pura | signature, regra/cálculo, pureza (sem efeitos) | comportamental |

> Mobile nativo (iOS SwiftUI/MVVM, Android Compose/MVVM): o "estado da tela"
> (diluído em screen+hook+store no React) concentra-se num **ViewModel** com um
> **UiState** tipado. Adiciona um tipo `viewmodel` (bloco: UiState, intents/actions,
> efeitos); a spec de View fica mais fina (bindings + estados visuais).

### Backend (Clean Architecture / serverless — ancorado em exemplos reais)

| tipo | fonte da variação | bloco de variação | regime |
|---|---|---|---|
| **entity / domain** | invariantes | campos e tipos, invariantes, máquina de estados, validações intrínsecas | misto (campos = declarativo; invariantes/transições = comportamental) |
| **usecase / interactor** | regras de negócio | input/output (DTO), pré-condições, regras passo a passo, erros de domínio, portas consumidas | comportamental |
| **repository (porta+adapter)** | acesso a dado | operações (CRUD), chaves/queries/índices, formato de retorno, paginação/consistência, erros de dados | comportamental |
| **service** | operação externa / orquestração | operação exposta, input/output, dependências externas, idempotência, retry/timeout | comportamental |
| **handler / trigger** _(interface)_ | rota **ou** evento + status | rota/método (ou tipo de evento), request (path/query/body/headers), response, status codes, auth, mapeamento erro→status | comportamental |
| **schema de dados** _(interface do dado)_ | modelo + acesso | modelos e campos, chaves e ÍNDICES, relações, AUTORIZAÇÃO (quem lê/escreve o quê), migração/versão | comportamental |
| **infrastructure / resource** | recursos + wiring | recursos provisionados, env/secrets, permissões (IAM), agrupamento de stack, cron/schedule | declarativo |

> Nuances reais: **"service" é termo sobrecarregado** — Clean clássico = serviço de
> domínio; em outros projetos = acesso remoto (chama uma função serverless). O
> projeto escolhe a definição e a spec a segue. **Handler serverless** descreve
> *evento*, não rota (API GW / AppSync / EventBridge / eventos de autenticação têm
> shapes distintos). **Repository** pode ser partido em dois (`repositories/` remoto
> + `models/` de acesso direto ao banco, por exemplo). **Go** materializa
> porta (interface de código) e adapter em arquivos separados — a spec pode cobrir os dois.

> **A `tela` (frontend) e o `handler`/`trigger` (backend) são o MESMO conceito: a
> interface (§2.5)** — o ponto de entrada do processo. Aparecem em famílias diferentes
> desta tabela só porque têm dialetos diferentes (props+navegação vs request/response/
> status), mas a régua de conteúdo da spec é a mesma. Um `HANDLER_GUIDE` espelha o
> `SCREEN_GUIDE`; ambos descrevem uma interface. Na CLI, o dialeto é o **comando** (linha
> "CLI command" abaixo) — também uma interface.

> **O `schema de dados` é REGIDO, não declarativo — e é fácil confundi-lo com
> `infrastructure/resource`.** A distinção: infra provisiona (uma tabela existe, com tal
> capacidade); o schema DECIDE (quais campos, quais índices, e **quem pode ler ou
> escrever cada modelo**). Um índice ausente torna uma query impossível meses depois; um
> `allow.owner()` é regra de negócio sobre visibilidade. Nada disso é dedutível do código
> que consome o dado — por isso a decisão precisa de spec.
>
> Não confundir com o **DAO** (`repository` partido em dois): o DAO *traduz* tabela↔objeto
> e não decide nada — é camada RECONHECIDA, sem spec. Quem decide é o schema. Dialetos:
> `amplify/data/resource.ts` (Amplify), `schema.prisma` (Prisma), migrations (Rails/
> Django), `CREATE TABLE` versionado (SQL puro), IaC de tabela + índices (DynamoDB/CDK).

### Infra / CLI / Lib

| tipo | fonte da variação | bloco de variação | regime |
|---|---|---|---|
| **módulo Terraform / IaC** | inputs → recursos | variáveis (nome/tipo/default/obrigatório), outputs, recursos criados, providers, preconditions, custo/blast-radius | declarativo |
| **manifesto K8s / Helm** | workload + rede | workload (tipo, réplicas), recursos (cpu/mem), rede (service/ingress), config/secrets, health probes, storage, políticas | declarativo |
| **Dockerfile / imagem** | ambiente de execução | base image, artefatos incluídos, entrypoint/cmd, portas, env vars, volumes, usuário/permissões | declarativo |
| **pipeline CI/CD** | stages + gates | gatilhos, stages/jobs e ordem, gates (lint/testes/aprovação), inputs/secrets, artefatos, alvos de deploy, rollback | misto (stages declarativos, gates comportamentais) |
| **config de deploy** | alvo + serviços | ambiente, serviços/réplicas, dependências, env/secrets, portas/volumes, estratégia (rolling/blue-green), healthcheck | declarativo |
| **CLI command** _(interface)_ | assinatura de invocação | args posicionais, flags/opções (default, env equivalente, exclusão mútua), comportamento, stdin, stdout (formato), exit codes, precondições | comportamental |
| **biblioteca / package** | API pública | superfície exportada (funções/tipos), contrato por símbolo (assinatura, erros, pré/pós), invariantes (thread-safety, complexidade), SemVer/deprecações | comportamental |

> Um artefato **declarativo** descreve forma/estado (inputs → recursos garantidos,
> sem sequência temporal) e casa mal com cenários `.feature` de step temporal — ele
> propaga para gates de conformidade estrutural, não para testes Gherkin. Um artefato
> **comportamental** casa naturalmente com dado/quando/então.

---

## 4. Como usar este catálogo

- A **Estrutura de Projeto** (`STRUCTURE.md`) do projeto declara *quais* dessas
  camadas/tipos existem naquele projeto.
- Para cada tipo declarado, o projeto tem **um template** = esqueleto comum (§1) +
  o bloco de variação da linha correspondente (§3).
- O **guide de spec** do projeto define as regras universais e aponta os templates.
- Cada `.spec.md` de um target concreto segue o template do seu tipo, e declara seu
  **regime** (§2), que determina o *tipo de teste* que o requisito exige
  (comportamental ou de conformidade). Testável por default; `@no-test` dispensa. A
  Propagação segue as arestas normalmente — o regime não roteia a onda.

Este catálogo é ponto de partida — cada projeto ajusta, remove e acrescenta tipos
conforme sua realidade. O que o Anchors padroniza é o **modelo** (esqueleto +
variação + regime), não a lista.

---

## 5. A declaração de dependências: como uma spec aponta as camadas que consome

O esqueleto (§1) tem o bloco **Dependências** ("de que camadas depende — e quais é
proibido chamar"). Este é o **formato concreto** desse bloco, e é o que materializa a
aresta de reúso que os pilares pressupõem: `STRUCTURE.md` §4 ("quando uma spec de
interface demanda alteração numa spec de usecase, é a Estrutura que diz que essa
dependência existe e em que direção") e `TRACEABILITY.md` §"Como as arestas entram no
mapa" (a aresta **`declared`** — a dependência conceitual que a inferência não pega).

### Por que declarar, e não inferir

Uma camada (usecase, repository, service, hook, store) existe porque é **reutilizada**:
consumida por N unidades. Comportamento reutilizável precisa ser **regulado por spec
própria** — senão cada consumidor o redefine e diverge. É a mesma razão pela qual um
componente reusado por N telas ganha spec própria (Atomic), só que no eixo de **dados**
em vez do **visual**. A dependência é **declarada por quem a tem** (a spec consumidora),
não inferida do import: a declaração é honesta (o autor afirma "eu dependo disto"),
estável (não quebra ao refatorar imports) e agnóstica de linguagem (o Anchors só lê
texto).

### O mecanismo: duas tabelas ligadas por código

A spec consumidora (uma tela, um usecase) já tem uma tabela de **contrato de dados**
(entrada/saída). Para não inchar essa tabela com origem de cada campo, a dependência
vive em **tabela própria, codificada**, e o contrato apenas **aponta o código**:

**Tabela de Dependências** — uma linha por origem de dado, com código local à spec:

| Cód  | Arquivo                   | Método         | Camada     |
| ---- | ------------------------- | -------------- | ---------- |
| DEP1 | `stores/auth.store.ts`    | `useAuthStore` | store      |
| DEP2 | `hooks/useAuth.ts`        | `signIn`       | hook       |

- **Cód** — `DEPn`, **local à spec** (numera de DEP1; único dentro da spec). A aresta
  real do grafo vem do **Arquivo**, não do código — o `DEPn` é só o ponteiro que o
  contrato usa.
- **Arquivo** e **Método** são colunas **separadas**: o mesmo método (`signIn`, `list`,
  `get`) pode morar em arquivos diferentes conforme a situação; juntá-los perderia a
  referência. A aresta liga ao **arquivo** (o nó do grafo); o **método** é **metadado
  da aresta** — habilita impacto fino ("mudou `signIn` → só as telas que usam `signIn`").
- **Camada** — o tipo (§3) da dependência, para os gates de limite de camada ("uma
  tela pode depender de hook/store/service; um repository não pode depender de tela").

**Contrato de dados** — a coluna de origem passa a referenciar o `DEPn` em vez de texto
livre:

| Campo       | Origem | Obrigatório | … |
| ----------- | ------ | ----------- | - |
| `isLoading` | DEP1   | ✅          | … |
| `user`      | DEP2   | ✅          | … |

### A aresta que nasce disso

Cada linha da Tabela de Dependências vira uma aresta **`depends-on`** (declared) da
spec consumidora para o **arquivo** referenciado, carregando o **método** como
metadado. É por essas arestas que a Propagação **desce** pelas camadas de dados: mudou
a spec do repository → as telas que o consomem ficam stale (a onda de reúso). É o
trilho que faltava para a propagação de dados descrita em `PROPAGATION.md` §7 fluir
para além da trinca de uma única unidade.

> Profundidade variável por projeto. A cadeia de dependência tem o comprimento que a
> Estrutura daquele projeto declara: um projeto vai `tela → repository/service` direto;
> outro vai `tela → usecase → repository → service`. A tabela é a mesma em cada elo —
> um usecase declara suas dependências de repository/service do mesmo jeito que a tela
> declara as dele. A recursão é a propagação seguindo a planta (`STRUCTURE.md` §4).

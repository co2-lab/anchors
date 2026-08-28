# Anchors — Pilar de Estrutura de Projeto

> Este documento define o **pilar de Estrutura de Projeto** do Anchors. Ele
> pressupõe o mecanismo de [`CONCEPT.md`](./CONCEPT.md) e é o **gabarito** sobre o
> qual os outros pilares operam: define quais camadas o projeto tem, como se
> organizam e em que ordem se relacionam.
>
> É teoria de framework, já validada por uma ferramenta que implementa
> (decompõe o projeto em camadas — infra, arquitetura, interface, usecase,
> repository, service, devops — e a cadeia de agentes respeita essa ordem)
> e por uma prova de conceito real.

---

## 1. A planta da casa

Para que o framework funcione, é preciso um acordo prévio: **quais camadas o
projeto tem e como elas se relacionam.** Sem esse acordo, cada plano e cada agente
inventa suas próprias camadas, e o projeto vira um amontoado sem forma comum.

A **Estrutura de Projeto** é esse acordo. Se a Rastreabilidade é a fiação e a
Propagação é a corrente, a Estrutura de Projeto é a **planta da casa** — onde ficam
os cômodos e como se conectam. Você pode fiar e passar corrente, mas é a planta que
diz quais cômodos existem e como se ligam. Sem planta, cada um constrói um
puxadinho.

Ela **sugere ou dita** a estrutura, e essa estrutura deve ser **respeitada e
documentada**. Como toda âncora, ela guia (o plano e os agentes consultam-na para
saber onde as coisas moram) e é confrontável (uma spec está na camada certa? a
propagação seguiu a ordem certa?).

---

## 2. O que a Estrutura de Projeto define

- **Quais camadas existem** — o vocabulário de camadas do projeto (ex.: infra,
  arquitetura, interface, usecase, repository, service, devops). Não é uma lista
  universal; cada projeto declara a sua.
- **A ordem / dependência entre camadas** — qual vem antes de qual, qual pode
  depender de qual. É isto que diz "usecase vem antes de interface" ou "interface
  depende de usecase".
- **Onde cada tipo de âncora mora** — em que camada vive uma spec, uma feature, um
  teste; a convenção de organização (co-location, pastas por camada, etc.),
  expressa como **patterns de arquivo/pasta** (`screens/**/*.spec.md` → spec de
  screen; `repositories/**/*.ts` → repository).
- **As regras de fronteira** — o que uma camada pode ou não pode fazer/importar. É
  o gabarito contra o qual se valida se algo respeitou os limites da sua camada.

> **Uma camada que atravessa aplicações: a INTERFACE.** Entre as camadas que um projeto
> declara, quase sempre há a **interface** — o ponto de entrada onde um trigger (interno
> ou externo) inicia o processo. Seu dialeto muda com o tipo de app (CLI: comando; GUI:
> tela; API: rota/handler; Job: trigger), mas é sempre a mesma camada conceitual, com a
> mesma régua de spec e a mesma posição na propagação (é a origem da onda; declara o que
> consome e desce para usecase/repository/service). Um projeto pode ter mais de um
> dialeto de interface (um projeto mobile com backend, por exemplo, tem `screen`
> no app e `handler`/`trigger` no backend) — são tags distintas, mesma régua.
> Ver `SPEC_TYPES.md` §2.5.

---

## 2.1 A Estrutura é o grafo virtual (o bootstrap do mapa)

Os patterns de arquivo/pasta por camada não servem só para organizar — eles dão à
Estrutura um papel que resolve o **paradoxo do mapa vazio**. Se o mapa de
dependências (`TRACEABILITY.md` §4) é construído incrementalmente por cada agente,
no início ele está vazio: no primeiro arquivo, não há aresta nenhuma. Como o watcher
sabe o que fazer?

A resposta: a Estrutura é o **grafo virtual** — o grafo *esperado*, antes de
qualquer aresta material existir. Pelos patterns, ao ver um arquivo, o framework já
sabe, sem consultar o mapa:

- **a que camada ele pertence** (o pattern classifica: `screens/Login.spec.md` → spec
  de screen);
- **que estrutura ele deveria ter ao redor** (uma spec de screen *deveria* ter um
  `.tsx`, um `.feature`, um `.test.tsx` co-localizados);
- **qual agente deve tratá-lo** (é uma spec → o agente de propagação).

Ou seja, há dois grafos complementares:

| | fonte | responde |
|---|---|---|
| **grafo virtual** (Estrutura) | patterns de arquivo/pasta por camada | "a que camada este arquivo pertence e que estrutura ele *deveria* ter" |
| **grafo material** (o mapa, Rastreabilidade) | arestas registradas por cada agente | "o que de fato depende de quê, agora" |

O grafo virtual dá o **esqueleto esperado**; o material dá as **relações realizadas**.
Um projeto novo tem grafo virtual completo (as regras já existem) e grafo material
vazio (nada foi criado ainda) — e é o virtual que permite ao watcher classificar o
primeiro arquivo e disparar o primeiro agente, que então começa a preencher o
material. Sem paradoxo de bootstrap: a Estrutura sabe a forma antes de haver conteúdo.

**A síntese:** a **Estrutura já é o diagrama de dependência** — as regras de quem
depende de quem, por camada, que existem sem nenhum arquivo. O **mapa** é a projeção
dessas regras sobre os arquivos reais, e serve para traçar o **caminho mínimo** de
impacto quando um arquivo altera (`TRACEABILITY.md` §4). A Estrutura dá o esqueleto;
o mapa dá o caminho.

## 2.2 A superfície da trinca: onde as peças *devem* morar

O grafo virtual diz que uma âncora *deveria* ter sua trinca ao redor — spec, feature,
teste. Mas **onde** essas peças moram é uma decisão de layout do projeto, e há duas
convenções legítimas:

- **Co-localizada** — as peças são irmãs no mesmo diretório (`Login.tsx`,
  `Login.spec.md`, `Login.feature`, `Login.test.tsx`). É o caso simples: o stem
  `{{dir}}/{{name}}` resolve tudo.
- **Centralizada / por região** — a peça mora numa árvore própria, longe da âncora
  (ex.: os testes de todas as Lambdas em `__tests__/unit/lambdas/`, não ao lado de
  cada `handler.ts`). Comum em backends.

Em ambos os casos, **a ligação material continua sendo por CÓDIGO** (`TRACEABILITY.md`:
"o código é a chave, não o nome do arquivo") — quem realiza qual requisito é o
scenario-code compartilhado, não a localização. O que muda é a **descoberta**: quando
não é co-localizado, o framework precisa saber *onde esperar* cada peça para poder
**confrontar a ausência**. É aí que entra o **padrão de localização da trinca**.

> **O caminho não é identidade — mas deve respeitar um padrão.** O caminho de um
> arquivo nunca identifica a unidade (isso é o código). Porém ele **precisa obedecer a
> um padrão estrutural declarado**: a spec de tal camada mora *aqui*, a feature *ali*,
> o teste *acolá*. Esse padrão é a **superfície de validação** dos gates de estrutura —
> a régua contra a qual se pergunta "as peças da trinca desta unidade estão nos lugares
> que o padrão manda?". Um teste que existe e liga por código, mas mora fora do padrão,
> é um **achado de estrutura** (não um furo de rastreabilidade, mas uma violação de
> layout) — a Qualidade decide se acusa.

Por isso a declaração de derivados (`Derived` na config) admite, além do stem
co-localizado (o default), **padrões de localização por camada-âncora**: para uma
âncora daquela camada, o teste/feature *deve* casar tal template de região. O template
pode capturar partes do caminho da âncora (ex.: o **módulo** — o diretório-pai quando o
arquivo é um `handler.ts`) para compor a região esperada. Assim o grafo virtual sabe a
forma esperada da trinca mesmo quando ela não é co-localizada, e o gate tem superfície
para confrontar — sem jamais tratar o caminho como identidade.

> É o mesmo classificador do gate "camada declarada" (§6): perguntar "a que camada
> esta peça pertence?" é aplicar os patterns do grafo virtual.

## 2.3 Regimes de verificação: cada cenário confronta a superfície do seu regime

A "peça de teste" da trinca não é uma só — um requisito pode ser verificado em **regimes
de teste** diferentes, e cada regime vive numa **superfície** própria:

| Regime (canônico) | O que verifica | Superfície típica |
|---|---|---|
| `unit` | regra/função pura, sem UI nem I/O | arquivo de teste unitário |
| `integration` | comportamento de uma unidade com suas dependências próximas | arquivo de teste de integração |
| `e2e` | fluxo do usuário atravessando telas/serviços | roteiro de ponta a ponta (ex.: `.yaml`) |
| `vr` | aparência/regressão visual | baseline de screenshot |

Uma **feature mistura regimes**: seus cenários declaram, cada um, em que regime serão
verificados (uma tag por cenário; um cenário pode ter mais de um regime). O gate de
correspondência **não joga todos os cenários contra o arquivo de teste** — ele roteia
**cada cenário para a superfície do regime que aquele cenário declara**. Um cenário `e2e`
é confrontado contra o roteiro e2e, não contra o teste unitário; um cenário `vr` contra o
baseline visual. Confrontar um cenário visual contra o teste unitário é um falso-negativo
estrutural — a peça existe, só mora noutra superfície.

> **O vocabulário de regime é do PROJETO; o regime canônico é do framework.** Um projeto
> pode chamar seus regimes como quiser (`@nivel-unit`, `@smoke`, `@wip`…). A Estrutura
> declara um **de-para** `tag-do-projeto → regime-canônico` (unit/integration/e2e/vr) para
> que o engine roteie sem conhecer a nomenclatura local — o mesmo princípio das camadas
> reconhecidas (`regime: declarativo`) e dos padrões de localização (§2.2): o mecanismo é
> universal, o vocabulário é local. Um projeto novo já nomearia seus cenários com os
> regimes canônicos e dispensaria o de-para. Uma tag de cenário sem regime mapeado não é
> confrontada por nenhuma superfície (fica fora do escrutínio — opt-out honesto).

Cada regime resolve para uma **superfície** (uma chave do `Derived`: onde aquela peça de
teste mora, co-localizada ou por padrão de localização §2.2). Assim o grafo virtual sabe
não só *que* um requisito deve ser testado, mas *em que regime* e *onde* essa
verificação mora — e o gate confronta cada cenário contra a superfície certa.

---

## 3. Meta-nível: a Estrutura define as camadas; as specs vivem dentro delas

Há uma distinção fina, mas essencial, entre a Estrutura de Projeto e as âncoras que
descrevem cada parte:

- A **Estrutura de Projeto** é o **meta-nível**: ela declara que *existe* uma camada
  de arquitetura, uma de usecase, uma de interface — e a ordem entre elas.
- Uma **spec** (inclusive a spec de arquitetura) é **conteúdo dentro** de uma das
  camadas que a Estrutura definiu.

Ou seja: a spec de arquitetura descreve *a arquitetura daquele projeto*; a Estrutura
de Projeto é o gabarito que diz *que existe uma camada de arquitetura, e onde ela
se encaixa na ordem*. O gabarito precede o conteúdo. Confundir os dois é confundir a
planta da casa com a decoração de um cômodo.

---

## 4. Os outros pilares operam sobre a Estrutura

A Estrutura de Projeto é transversal — os demais pilares a pressupõem:

- **O Planejamento** (`PLANNING.md`) organiza segundo ela. O plano *pode* se
  organizar por camada, mas as camadas em si vêm da Estrutura. A exigência não é do
  plano — é da Estrutura: ela é que diz quais camadas existem e sua ordem; o plano
  apenas as usa.
- **A Propagação** (`PROPAGATION.md`) flui na ordem que a Estrutura define. Quando
  uma spec de interface demanda alteração numa spec de usecase, é a Estrutura que
  diz que essa dependência existe e em que direção — a onda segue essa planta.
- **A Rastreabilidade** (`TRACEABILITY.md`) liga camadas que a Estrutura declara. O
  mapa de dependências respeita as fronteiras da planta.
- **A Qualidade** (`QUALITY.md`) tem gates que confrontam contra a Estrutura: "esta
  peça respeita os limites da sua camada?".

A Estrutura é o gabarito; os outros pilares desenham sobre ele.

---

## 5. A falha característica: a anarquia estrutural

Cada pilar tem sua falha. A da Estrutura de Projeto é a **anarquia estrutural** —
sem uma estrutura definida e respeitada, cada plano e cada agente inventa suas
próprias camadas e fronteiras. O projeto perde a forma comum: peças moram em
lugares inconsistentes, camadas se misturam, a ordem de dependência é violada sem
que ninguém saiba, porque não havia acordo do que era a ordem.

É uma falha distinta:

- não é desconexão (**Rastreabilidade**) — as peças podem estar ligadas, só que
  ligadas de qualquer jeito;
- não é onda incompleta (**Propagação**) — a onda pode propagar, mas por caminhos
  que não deveriam existir;
- não é trabalho sem norte (**Planejamento**) — pode haver um plano claro, mas
  executado sobre uma estrutura que muda a cada passo.

É **ausência de gabarito**. Um projeto sem estrutura respeitada não tem uma forma
que uma sessão futura possa reconhecer e continuar — cada sessão redesenharia a
planta. A Estrutura de Projeto é o que garante que todos constroem a mesma casa.

---

## 6. A régua que a Estrutura dá aos gates de conformidade

A Estrutura não tem mecânica de gate própria — **a mecânica de gate pertence à
Qualidade** (`QUALITY.md` §2). O que a Estrutura fornece é a **régua**: o critério
de conformidade à planta que o gate de arquitetura da Qualidade (`QUALITY.md` §3)
confronta. Assim como a Rastreabilidade fornece a régua de *conexão*, a Estrutura
fornece a régua de *conformidade estrutural* — quem executa o confronto é a
Qualidade.

O que a Estrutura dá para confrontar:

| dimensão medida | pergunta de confronto | falha = |
|---|---|---|
| **camada declarada** | toda peça pertence a uma camada que a Estrutura define? | peça fora da planta |
| **fronteira respeitada** | a peça respeita o que sua camada pode/não pode fazer? | violação de camada |
| **ordem honrada** | as dependências seguem a ordem que a Estrutura define? | dependência invertida |
| **organização honrada** | as âncoras estão onde a convenção manda (co-location, pastas)? | peça no lugar errado |

Como todo gate de Qualidade, o gate de arquitetura roda sobre o caminho de impacto
da Propagação (`PROPAGATION.md` §3), emite issues materiais quando falha, e matura
de `informativo` a `bloqueante` (`QUALITY.md` §7) — um projeto pode começar com a
estrutura frouxa (report-only) e endurecê-la conforme amadurece. A Estrutura é dona
do *critério*; a Qualidade, da *execução*.

---

## 7. Relação com os outros pilares

A Estrutura de Projeto é o gabarito de base:

- **Antecede o Planejamento**: o plano só pode organizar por camada porque a
  Estrutura já definiu as camadas.
- **Guia a Propagação**: a onda segue a ordem de dependência que a Estrutura
  declara.
- **Emoldura a Rastreabilidade**: o mapa de dependências respeita as fronteiras da
  planta.
- **Dá régua à Qualidade**: os gates de fronteira/camada confrontam contra a
  Estrutura.

Num roteiro de amadurecimento, a Estrutura tende a vir muito cedo — junto ou antes
do Planejamento —, porque é o gabarito que todos os outros pressupõem.

---

## 8. Resumo do pilar

- **Estrutura de Projeto = a planta da casa.** Define quais camadas existem, sua
  ordem/dependência, onde cada âncora mora e as regras de fronteira. Sugere ou dita,
  e deve ser respeitada e documentada.
- **É o meta-nível.** Declara *que existem* as camadas; as specs (inclusive a de
  arquitetura) são conteúdo *dentro* das camadas. O gabarito precede o conteúdo.
- **Os outros pilares operam sobre ela.** Planejamento organiza por ela, Propagação
  flui na ordem dela, Rastreabilidade liga camadas que ela declara, Qualidade
  confronta contra ela.
- **Falha característica = anarquia estrutural.** Sem gabarito respeitado, cada
  sessão redesenha a planta e o projeto perde a forma comum.
- **Dá a régua aos gates de conformidade.** Camada declarada, fronteira respeitada,
  ordem honrada, organização honrada. A mecânica de gate é da Qualidade; a Estrutura
  fornece o critério que o gate de arquitetura confronta.

O que o pilar entrega: um projeto com **forma reconhecível** — uma planta que toda
sessão respeita, de modo que o que uma constrói a próxima entende e continua. É o
gabarito que impede o projeto de virar um amontoado de puxadinhos.

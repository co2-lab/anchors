---
title: Planejamento
---

> Este documento define o **pilar de Planejamento** do Anchors. Ele pressupõe o
> mecanismo de [`CONCEPT.md`](/docs/conceito/) — âncora (na metáfora da escalada),
> grafo, propagação, issues. O Planejamento é a **origem**: a âncora que aponta a
> rota, o input do fluxo de trabalho, de onde vem a primeira alteração.
>
> É teoria de framework, já validada por uma ferramenta que implementa (o
> agente *Planner* produz o *Plan* que dispara toda a cadeia) e por uma
> prova de conceito real.

---

## 1. A leitura da rota antes de subir

Todo escalador lê a parede antes de tocar nela: onde estão os apoios, por onde
subir, em que ordem cravar as âncoras. O **Planejamento** é essa leitura. É a
âncora que **aponta para onde ir** — a primeira das três funções da escalada
(`CONCEPT.md` §2).

Ele é o pilar da **origem**. Todo o resto do framework é reativo: a Propagação
reage a uma alteração, a Qualidade reage a um confronto, a Rastreabilidade liga o
que já existe. Mas nada disso responde à pergunta anterior a todas: **de onde vem
a primeira alteração?** O fluxo de trabalho começa em "uma alteração acontece" —
como se ela caísse do céu. O Planejamento é a resposta: ele **semeia as specs de
partida**. É a origem da onda.

Ponto crucial: o plano **não planeja código** — planeja **specs**. Sua saída é o
conjunto de specs que precisam nascer ou mudar. O código nunca é produto do plano;
é produto da *propagação das specs* (§3). O Planejamento é a origem, mas opera **um
nível acima do código**: semeia as âncoras-base, e são elas que descem até o código.

Sem Planejamento, você tem um sistema que **reage** perfeitamente mas não **avança
com direção**. Ele propaga, valida, mede — mas responde a mudanças avulsas em vez
de a um plano. É a diferença entre escalar uma rota e agarrar pedras ao acaso.

---

## 2. O plano é uma âncora

O plano cumpre as funções de âncora como qualquer outra:

- **Aponta** — diz o que construir ou evoluir, decomposto e ordenado. É a rota.
- **Segura a corda** — é confrontável: *"o que foi feito corresponde ao que foi
  planejado? o plano foi cumprido?"*. Um plano que ninguém confronta é só uma lista
  de desejos; um plano-âncora é o critério contra o qual o progresso é medido.
- **Demarca** — registra a intenção e sua ordem, deixando o rastro do *porquê*
  estamos subindo por aqui.

Como âncora, o plano vive no grafo: ele **rege** as specs que dele derivam (uma
aresta `governs` do plano para cada spec que ele manda criar), propaga quando muda,
e fica stale quando a realidade diverge dele. Repara que o plano rege **specs**, não
código — a fronteira do plano é a spec (§3).

---

## 3. O plano é uma semente, não um blueprint completo

Um plano não é uma lista de tarefas solta, nem uma enumeração exaustiva de tudo que
o projeto precisa. Ele **inicia** as specs necessárias — o mínimo para dar partida
— e deixa a propagação fazer o resto.

O que um plano captura:

- **Semeadura** — as specs que precisam nascer ou mudar para dar partida. É o que
  liga o Planejamento à Spec: o plano diz *quais specs de partida*, a spec diz *o
  quê* cada uma.
- **Ordem / fases (opcional)** — a sequência em que as partes são construídas. Um
  plano *pode* se organizar por camada e fase, mas **não é exigência**: as camadas e
  sua ordem vêm da **Estrutura de Projeto** (`STRUCTURE.md`), não do plano. O plano
  usa as camadas que a Estrutura define; organizar-se por elas é escolha de quem
  planeja.
- **Progresso** — o estado de cada parte (a fazer / em andamento / feito). É o que
  torna o plano confrontável e o que responde "onde paramos?".
- **Direção / o porquê** — a intenção que justifica a rota. É isto que carrega o
  norte entre sessões (§6).

### Propagação sob demanda

O plano **não precisa** enumerar toda a cascata de camadas. Uma spec de interface,
ao ser implementada, pode demandar alteração numa spec de usecase — e essa
propagação **acontece na hora, pelo agente**, seguindo a Estrutura de Projeto, sem
precisar estar no plano. Uma spec propaga para outras specs, para código, para
features, para testes, para execução — tudo pela mesma máquina de **Propagação**
(`PROPAGATION.md`). É propagação **sob demanda**: o plano semeia; a cascata inteira
emerge da propagação a partir da semente.

Por isso Planejamento e Spec são adjacentes mas distintos: o plano semeia as specs
*de partida*; as specs (via Propagação + Estrutura) geram as *demais* sob demanda.
O plano não precisa saber a árvore inteira — só o suficiente para a propagação
assumir.

A prova de conceito materializa a semeadura: o *Planner* produz um *Plan*; o *Plan
Executor* **gera as specs** nas pastas correspondentes (nunca código-fonte) e
atualiza o progresso a cada spec gerada. Daí em diante, são as specs que propagam.

### Dois modos: semente e reforço (ponto aberto)

O plano tem dois modos de operação:

- **Modo semente (padrão)** — inicia as specs mínimas e deixa a propagação e os
  agentes resolverem o resto sob demanda.
- **Modo reforço (exceção)** — às vezes o agente interpretaria uma spec de forma
  errada ao propagar. Isso costuma ser sinal de uma **falha de especificação ou de
  guide** (algo não estava claro o bastante — e nem tudo dá para prever). O contorno
  é o plano **fixar/ditar** detalhes de comportamento da cadeia de propagação, para
  forçar o resultado certo e evitar a interpretação errada.

> **Abertura para refino futuro.** O modo reforço ainda não está refinado. Ficam
> registradas as perguntas em aberto, sem resposta por ora: *quando reforçar vs.
> deixar propagar? um reforço no plano é um cheiro de spec/guide fraca (e deveria
> virar uma correção na spec/guide, não uma fixação no plano)? como o reforço
> interage com a propagação sob demanda?* São pontos a lapidar quando o uso real
> mostrar o padrão.

---

## 4. Quando o plano compensa (e de onde mais ele vem)

O Planejamento não é obrigatório para toda mudança — é uma **recomendação de
eficiência**, não uma regra. Duas coisas o alimentam além da intenção inicial do
usuário: mudanças de alto grau e issues convertidas.

### Mudança de alto grau: o plano é o caminho mais eficiente, não o único

Mudar um nó de **alto grau** — um guide, uma spec de arquitetura, algo que rege
muitos arquivos — pode ser feito de dois jeitos, e **ambos funcionam**:

- **Direto** — edita-se o guide, e a Propagação faz o resto: onda global, todas as
  specs regidas ficam stale, o impacto é descoberto *na marcha*, âncora por âncora
  (`PROPAGATION.md` §3). Chega ao resultado certo, mas é **caro** — reprocessa a
  árvore reativamente.
- **Via plano** — o plano já **traz as revisões de spec** que a mudança implica e
  **mapeia a árvore afetada de antemão**. A propagação que segue é dirigida, não
  cega. Mesmo resultado, **mais eficiente**.

Por isso o Anchors **recomenda** (não exige) que mudanças de alto grau nasçam de um
plano: não porque o direto está errado, mas porque o plano evita o custo da onda
global às cegas. O grau do nó (`PROPAGATION.md` §3) não decide só o *tamanho da
onda* — sinaliza *quando o plano compensa*. Alto grau → o plano paga por si.

### A issue como origem: a terceira rota do usuário

O plano também nasce de uma **issue convertida**. Diante de uma issue, o usuário tem
três rotas (`CONCEPT.md` §5): resolver ele mesmo, delegar a um agente, ou **converter
em plano** quando a resolução é trabalho estruturado, não um retoque. A conversão
realimenta o fluxo: a issue encerra (`done/`, "convertida no plano 00XX") e o plano
nasce com o propósito ("aberto para resolver a issue 00XX") — referência cruzada
bidirecional, o débito transferido com rastro, não apagado.

Isso fecha o ciclo do framework: o Planejamento é a **origem** (semeia specs) e
também o **destino de reentrada** — uma issue grande volta ao começo virando plano,
por escolha de quem opera. O fluxo não é linear com recuo; é genuinamente cíclico.

---

## 5. A falha característica: trabalho sem norte

Cada pilar tem sua falha. A do Planejamento é o **trabalho sem norte** — alterações
que acontecem, propagam, passam nos gates, e ainda assim o projeto não vai a lugar
nenhum coerente porque ninguém decidiu *o que* construir e *em que ordem*.

É uma falha distinta de todas as outras:

- não é desconexão (**Rastreabilidade**) — as peças podem estar todas ligadas;
- não é âncora que mente (anti-drift) — cada âncora pode estar verdadeira;
- não é baixa qualidade (**Qualidade**) — cada gate pode passar;
- não é onda incompleta (**Propagação**) — tudo pode ter propagado direitinho.

É **ausência de intenção estruturada**. Um projeto pode ter todos os outros pilares
vigorosos e ainda assim ser um moto-perpétuo sem destino: muito movimento, nenhum
avanço. O Planejamento é o pilar que dá **vetor** ao que seria só agitação.

---

## 6. O Planejamento carrega o norte entre sessões

A dor original que motivou o Anchors: *as sessões futuras precisam continuar com o
mesmo rigor, seguindo a mesma direção das sessões anteriores.* O Planejamento é
onde essa **direção** vive.

Uma sessão futura — outro agente, outro dia, outro humano — que lê o plano sabe
três coisas que nenhum outro artefato entrega juntas:

- **para onde o projeto vai** (a intenção, a rota);
- **em que ordem** (as fases);
- **onde paramos** (o progresso).

O "o que não posso quebrar" e o "onde retomar" nascem do plano. Por isso o
Planejamento **absorve a continuidade de direção entre sessões**: ela não precisa
ser um pilar separado — é uma *consequência* de o plano ser uma âncora bem feita,
viva e confrontada. Um plano atualizado é o handoff: a próxima sessão pega a corda
exatamente onde a anterior a deixou.

---

## 7. Relação com os outros pilares

- **Usa a Estrutura de Projeto** (`STRUCTURE.md`): as camadas em que o plano se
  organiza — e a ordem entre elas — vêm da Estrutura, não do plano. A Estrutura é o
  gabarito; o plano semeia specs *dentro* dele.
- **Precede a Spec** (`SPEC.md`): o plano diz *quais* specs de partida nascem; a
  spec diz *o quê* cada uma é. Planejamento é o vetor (a direção do movimento); Spec
  é o alvo (o destino verdadeiro). O plano inicia a rota; a spec dá a cada ponto um
  destino confrontável.
- **Dá origem à Propagação** (`PROPAGATION.md`): a primeira alteração de qualquer
  ciclo nasce de um plano — que semeia specs. O Planejamento é o gatilho a montante
  da onda; a Propagação então espalha a partir das specs semeadas, sob demanda.
- **Define o escopo da Qualidade** (`QUALITY.md`): a prioridade e o alvo de cada
  parte do plano informam o que precisa estar pronto para promover — o plano diz o
  que é crítico; a Qualidade defende o limiar.
- **É confrontado como qualquer âncora** (`CONCEPT.md` §2): "o feito corresponde ao
  planejado?" é o confronto do plano. Divergência gera issue, pelo mecanismo comum.
- **Dá a régua ao gate de validação de plano** (`QUALITY.md` §5.1): o plano
  recém-criado é coerente/exequível? — review de IA, na borda de entrada do fluxo. A
  Qualidade executa; o Planejamento dá o critério. (Já o *fechamento do laço*
  issue→plano não é um gate pontual — é verificado pelo **validador de saúde do
  ecossistema**, `QUALITY.md` §5.2, que varre laços que nunca fecharam.)

---

## 8. Resumo do pilar

- **Planejamento = a origem.** Semeia as specs de partida; é o input do fluxo, de
  onde vem a primeira onda. Nada no framework responde "de onde vem a mudança?" a
  não ser ele.
- **O plano gera specs, não código.** Sua fronteira é a spec — opera um nível acima
  do código. O código é produto da propagação das specs, nunca do plano.
- **O plano é uma âncora.** Aponta a rota, segura a corda (é confrontável: "o feito
  corresponde ao planejado?"), e demarca a intenção. Vive no grafo, rege as specs
  que dele derivam, propaga.
- **É semente, não blueprint.** Inicia as specs mínimas; a cascata emerge por
  **propagação sob demanda** (uma spec propaga para outras specs/código/features/
  testes, seguindo a Estrutura). Modo **semente** (padrão) vs. **reforço** (fixar
  detalhes no plano quando a propagação erraria — ponto aberto para refino).
- **Captura** semeadura, ordem/fases (opcional, vinda da Estrutura), progresso e
  direção.
- **Recomendado (não exigido) para mudanças de alto grau.** Editar um guide direto
  também propaga e funciona — mas o plano, que traz as revisões e mapeia a árvore, é
  *mais eficiente* que a onda global às cegas. Custo, não correção.
- **Origem e destino de reentrada.** Nasce da intenção, de mudanças de alto grau, e
  de **issues convertidas** (a 3ª rota do usuário, `CONCEPT.md` §5): a issue encerra
  apontando o plano, o plano nasce com o propósito. O fluxo é cíclico, não linear.
- **Falha característica = trabalho sem norte.** Muito movimento, nenhum avanço. Dá
  vetor ao que seria só agitação.
- **Carrega o norte entre sessões.** Para onde vamos + em que ordem + onde paramos.
  Absorve a continuidade de direção — o plano atualizado é o handoff.
- **Precede a Spec.** Vetor (Planejamento) e alvo (Spec) são as duas origens: uma
  do movimento, outra da verdade.

O que o pilar entrega: um projeto que **avança com direção**, não só reage. Cada
ciclo de trabalho começa numa intenção estruturada, e qualquer sessão futura sabe a
rota, a ordem e o ponto de retomada. É o que impede um projeto de ser um organismo
saudável que anda em círculos.

---
title: Rastreabilidade
---

> Este documento define o **pilar de Rastreabilidade** (Traceability) do Anchors.
> Ele pressupõe o mecanismo geral de [`CONCEPT.md`](/docs/conceito/) — âncora, grafo,
> sincronia, issues — e o especializa para um propósito: garantir que cada peça do
> projeto tenha uma **identidade contínua** e esteja **conectada**, de modo que o
> projeto seja um organismo único, não um monte de arquivos que por acaso moram na
> mesma pasta.
>
> É teoria de framework, já validada por uma ferramenta que implementa o
> Anchors e por uma prova de conceito real (incompleta e imatura de
> propósito). Ambas aparecem apenas como instâncias.

---

## 1. A floresta, não as árvores

Numa floresta, o que faz dela um único organismo vivo não são as árvores — é a
**trama de fungos** que corre por baixo do solo e conecta as raízes de todas
elas. Essa rede (a rede micorrízica) leva nutrientes e sinais de uma árvore a
outra. Sem ela, você não tem uma floresta; tem árvores isoladas dividindo o mesmo
terreno.

Um projeto de software é igual. As **âncoras** (spec, feature, teste, código) são
as árvores. A **Rastreabilidade** é a trama por baixo: o que faz cada requisito
manter a mesma identidade através de todas as suas formas, e o que garante que
nenhuma peça vira uma ilha desconectada. É a cola que transforma arquivos
co-localizados num tecido único.

> Enquanto o **grafo** (`CONCEPT.md` §3) é a estrutura visível — quais nós
> existem e como se ligam —, a **Rastreabilidade** é o que faz cada ligação ser
> confiável, única e sem furos. O grafo é o mapa; a Rastreabilidade é o que
> garante que o mapa corresponde à realidade e que nenhuma árvore ficou fora dele.

### As duas metades da Rastreabilidade

A Rastreabilidade tem **dois lados**, e um projeto pode ter um sem o outro:

- **Identidade** — "este teste realiza *aquele* requisito". Liga um requisito às
  suas encarnações (spec → feature → teste → código) por uma chave estável (§3).
  Responde: *o mesmo requisito, através de suas formas.*
- **Dependência** — "este arquivo depende *daqueles*; se mudar, *estes* ficam
  impactados". É o **mapa de dependências** entre os arquivos (§4). Responde: *o que
  se liga a quê.*

Os dois são a mesma cola vista de ângulos diferentes: a identidade conecta as
*formas de um requisito*; a dependência conecta os *arquivos entre si*. E os dois se
materializam no **mesmo artefato** — o mapa (§4) —, onde um nó pode ser alvo de uma
aresta de identidade (spec→teste) e de uma aresta de dependência (arquivo→arquivo)
ao mesmo tempo.

> Sem a metade de dependência, não há sobre o que a **Propagação** operar — ela
> percorre o mapa para calcular o caminho de impacto (`PROPAGATION.md`). A
> Rastreabilidade **mantém** o mapa (a estrutura); a Propagação o **consome** (a
> dinâmica). Um artefato, dois pilares.

---

## 2. Por que Rastreabilidade é um pilar

A Rastreabilidade é **transversal**: todos os outros mecanismos do Anchors a
pressupõem.

- O **grafo** só existe porque há como ligar `spec ↔ feature ↔ teste ↔ código`.
  Sem rastreabilidade, não há arestas para desenhar.
- O pilar de **Qualidade** (`QUALITY.md` §4) depende inteiramente de saber "este
  cenário virou teste?". Sem a cola, o gate feature→teste é cego.
- O **anti-drift** (`CONCEPT.md` §2) só é verificável porque se consegue rastrear
  "este identificador citado na âncora → esta peça real".
- As **issues** nomeiam a aresta que quebrou — de novo, rastreabilidade.

Quando a cola falha, tudo desliga em silêncio: o grafo fica furado, os gates
ficam cegos, e peças se acumulam sem par. Por ser pressuposta por todos, se ficar
implícita ela apodrece sem ninguém perceber. É o mesmo critério que elevou a
Qualidade a pilar: **nomear para que não seja tratada com descuido.**

E ela passa no teste de "ser um pilar de verdade": tem **gates próprios** e uma
**falha característica** que não pertence a nenhum outro pilar (§6). Há uma
simetria com a Qualidade: **Qualidade mede se o trabalho é bom; Rastreabilidade
mede se o trabalho está conectado.** Os dois são "medir" — um mede o nível, o
outro mede se as peças formam um tecido único.

---

## 3. Identidade contínua: o código estável

O coração da Rastreabilidade é **preservar e propagar** a **identidade estável**
que cada requisito recebe na spec (`SPEC.md` §4), de modo que ela sobreviva através
de todas as suas encarnações. Um requisito nasce numa spec, é descrito como
comportamento numa feature, é implementado como teste, e existe como código. São quatro formas — mas é **a mesma coisa**. O que as amarra como a mesma
coisa é um **código estável** carregado por todas.

```
spec.md   declara   →  CODE-S01     (o requisito nasce com identidade)
feature   marca     →  @CODE-S01     (a mesma identidade, outra forma)
teste     titula    →  CODE-S01: …   (a mesma identidade, executável)
código    realiza   →  (rastreável via a cadeia acima)
```

A cadeia é **universal** — vale para specs comportamentais e declarativas. O que
varia é o *tipo* de teste na ponta: comportamental (cenário Gherkin) para specs de
comportamento; **de conformidade** (o recurso declarado bate com a forma?) para
specs declarativas (`SPEC.md` §6). "Teste" no Anchors abrange os dois. Um requisito
que o usuário não quer testar é marcado **`@no-test`** — a tag o dispensa da
exigência de teste sem quebrar a cadeia (a identidade continua rastreável nas demais
formas).

Propriedades que fazem a identidade funcionar (destiladas do POC):

- **O código é a chave, não o nome do arquivo.** Ligar por nome de arquivo é
  frágil (renomeou, quebrou; dois arquivos parecidos confundem). O código estável
  atravessa renomeações e reorganizações — a identidade não mora na localização.
- **A relação é N:M.** Um requisito pode ter várias encarnações do mesmo tipo (um
  cenário coberto por vários testes) e várias formas (o mesmo código aparece na
  spec, na feature e em testes de níveis diferentes). O código as reconcilia todas.
- **O código tem uma gramática de fonte única.** Existe *um* lugar que define como
  um código é formado e lido, importado por todos os validadores de todas as
  linguagens. É o ponto de acoplamento que faz peças de mundos diferentes (uma
  spec markdown, um teste em jest, um teste em go) falarem a **mesma língua de
  identidade**. Sem fonte única, cada validador inventa sua própria noção de
  código e a cola se parte.

### Quando a identidade diverge: o drift silencioso

A falha de identidade não é abstrata. No POC há um caso real: a tela de login foi
renomeada, e seu código estável passou de `LGNN` para `LOGI`. A spec, a feature e a
maioria dos testes foram atualizados — mas **um teste ficou para trás**, ainda
citando `LGNN-S02`. A feature declara `@LOGI-S02`; o gate procura um teste que cite
`LOGI-S02`; o `LGNN-S02` obsoleto **não casa**. Resultado: a feature parece
descoberta, embora o teste exista — a cola quebrou silenciosamente porque a string
de identidade divergiu entre duas encarnações do mesmo requisito.

É exatamente por isso que a gramática precisa ser **fonte única** e a identidade
precisa ser **a chave, não o nome**: a divergência de uma string de identidade é
invisível a olho nu, mas visível ao gate — que confronta por interseção exata de
códigos. O drift de identidade é o análogo, no lado identidade, do que o anti-drift
é para o conteúdo da âncora.

> **Fronteira com o grafo:** o grafo (`CONCEPT.md`) define *o que é uma aresta* —
> nós, direção, tipo, sincronia. A Rastreabilidade define *como as arestas passam
> a existir e se mantêm honestas* — a chave estável que permite afirmar "este nó é
> a encarnação daquele", a unicidade dessa chave, e a ausência de nós órfãos. O
> grafo usa a Rastreabilidade para se materializar; a Rastreabilidade garante que
> o grafo não tem furos nem duplicatas.

---

## 4. O mapa: o artefato material

A metade de dependência da Rastreabilidade (§1) precisa de um lugar onde morar. Esse
lugar é **o mapa** — o artefato material, versionado no repositório, que lista os
arquivos do projeto e as arestas entre eles. Sem esse arquivo, "o projeto sabe suas
dependências" é uma afirmação sem sustento: a informação existe espalhada nos
imports, mas não há um lugar único, consultável e confrontável que a materialize. O
mapa é esse lugar.

O mapa **é a materialização do grafo** (`CONCEPT.md` §3) — não um segundo artefato.
O CONCEPT define o grafo em abstrato (nós, arestas tipadas, versionado); o mapa é o
arquivo concreto onde esse grafo vive (o `anchors.graph.yaml` do exemplo em
`CONCEPT.md` §3), e é a Rastreabilidade que o mantém.

O mapa é a **projeção da Estrutura sobre os arquivos reais**. A Estrutura de Projeto
(`STRUCTURE.md` §2.1) já é o *diagrama* de dependência — as regras de quem depende de
quem, por camada, que existem sem nenhum arquivo. O mapa instancia essas regras nos
arquivos que de fato existem, com versão e carimbo por aresta — e é isso que permite
traçar o **caminho mínimo** de impacto quando um arquivo altera (a Estrutura diz *que
tipo* depende de *que tipo*; só o mapa diz *qual arquivo* exatamente). A Estrutura dá
o esqueleto; o mapa dá o caminho.

> É o furo que o POC ainda não preencheu: ele faz muito bem a identidade (o código
> de cenário atravessando spec→feature→teste), mas **não tem um arquivo com o mapa
> de arquivos e dependências**, nem versão por arquivo, nem carimbo de quando cada
> relação foi validada. O mapa é o que materializa a Rastreabilidade de dependência.

### Anatomia do mapa

O mapa carrega três informações que hoje faltam — cada uma responde a uma pergunta
que a Propagação precisa fazer:

**Nós** — um por arquivo que participa do grafo:

| campo do nó | responde |
|---|---|
| `id` · `kind` | qual arquivo, de que tipo (os valores literais de `kind` — `spec`/`feature`/`test`/`code`/`doc`/`guide` — em `CONCEPT.md` §3) |
| **`rev`** (versão do conteúdo) | *"esta é a mesma versão de quando validei?"* — a **versão do arquivo** |
| **`updated_at`** | *"quando este arquivo mudou pela última vez?"* — o **carimbo de alteração** |

**Arestas** — uma por relação, de identidade ou de dependência:

| campo da aresta | responde |
|---|---|
| `from` · `to` · `type` | quem depende/rege quem (`CONCEPT.md` §3) |
| `origin` | como a aresta entrou no mapa (abaixo) |
| **carimbo de validação** (`validated_from_rev`, `validated_to_rev`, `last_validated`, `verdict`) | *"contra que versões de cada ponta esta relação foi confrontada, quando, e com que veredito?"* |

Este é o conjunto completo de campos do carimbo — a Rastreabilidade é a fonte-única
dele. O `verdict` (`ok`/`issue`/`pending`) é escrito pelo gate que confronta a
aresta (`QUALITY.md`) e lido pela Propagação; os campos `validated_*` e
`last_validated` guardam contra o quê e quando.

A **versão do nó** e o **carimbo da aresta** são a matéria-prima da sincronia: a
Propagação (`PROPAGATION.md` §3) compara `rev` atual do nó com a `rev` que o carimbo
guarda — se avançou, a relação está stale. A Rastreabilidade **descreve e mantém**
esses campos no mapa; a Propagação os **lê** para calcular o caminho de impacto. É a
mesma fronteira estrutura/dinâmica das duas metades.

### Como as arestas entram no mapa

As arestas nascem de três origens, que **convivem** no mesmo mapa material — com a
**inferência como motor**:

- **`inferred` (o motor)** — o framework lê imports/símbolos do código e **propõe** a
  maioria das arestas automaticamente. É o que torna o mapa viável num projeto
  grande: ninguém desenha o grafo à mão.
- **`convention`** — a co-location gera as arestas óbvias (`login.tsx` ↔
  `login.spec.md` ↔ `login.feature`), sem precisar de inferência nem declaração.
- **`declared`** — a spec (ou o próprio mapa) declara explicitamente as arestas que a
  inferência não pega — uma dependência conceitual que não aparece como import, uma
  spec de arquitetura que rege N interfaces.

Inferência e convenção *propõem*; o mapa material é onde as arestas passam a existir
de fato, versionadas e revisáveis em PR. O mapa é a fonte de verdade; qualquer
índice em memória é cache reconstruível a partir dele (`CONCEPT.md` §3).

**A aresta de reúso entre camadas (`depends-on`) é `declared`.** Quando uma unidade
consome outra camada — uma tela que usa um hook/store, um usecase que usa um
repository — essa dependência não é co-location (as peças não moram juntas) nem sempre
é inferível (a relação é conceitual). Ela é **declarada na spec consumidora**, pelo
formato codificado de `SPEC_TYPES.md` §5 (a Tabela de Dependências: `DEPn · Arquivo ·
Método · Camada`). Cada linha vira uma aresta `depends-on` da spec para o **arquivo**
referenciado, com o **método** como metadado (para impacto fino). É por essas arestas
que a Propagação desce pelas camadas de dados — sem elas, a mudança numa spec de
repository nunca alcançaria as telas que dela dependem (uma **dependência invisível**,
abaixo). Declará-las é o que dá trilho à propagação de reúso.

### O mapa também não pode mentir

Como toda âncora, o mapa tem seu anti-drift. Uma aresta que aponta para um arquivo
que não existe mais é uma mentira estrutural; uma dependência real que o mapa não
registrou é uma **dependência invisível** — o análogo, no lado do mapa, do órfão de
identidade (§6). A dependência invisível é perigosa porque quebra a Propagação em
silêncio: a onda não percorre uma aresta que não está no mapa, então um arquivo que
*deveria* ter ficado stale nunca é reconferido. Manter o mapa honesto (arestas
inferidas re-propostas quando o código muda, arestas mortas removidas) é uma
obrigação do pilar.

### O mapa é uma âncora auto-referente

Uma pergunta natural: se toda âncora tem dependências registradas no mapa, quais são
as dependências *do próprio mapa*? A resposta dissolve a regressão: **as dependências
do mapa são todos os documentos contidos nele.** O mapa não precisa de uma lista de
dependências sobre si mesmo — ele *é* a lista. É auto-referente por natureza: depende
de tudo que mapeia, e isso já está escrito nele. Não há "quem mapeia o mapa"; ele se
contém.

### A lei de manutenção: quem mexe num arquivo atualiza o mapa

O mapa não é *construído* por um agente-dono — ele é **mantido por toda execução de
agente**, como obrigação transversal. A regra é simples e vale para agentes de
qualquer pilar:

> **Todo agente que cria, move ou remove um arquivo atualiza o mapa no mesmo ato** —
> registrando os nós e arestas que nasceram, movendo os que mudaram de lugar,
> removendo os que sumiram.

É por isso que a criação do mapa não precisava de um mecanismo especial (a suposta
ponta "o framework percorre o grafo mas não o cria"): a criação é **distribuída**.
Quando o agente de planejamento gera uma spec, ele registra a aresta `plano→spec`;
quando o agente de propagação cria o código e a feature de uma spec, ele registra
esses nós e suas arestas. A aresta nasce do **ato de criação**, não de uma fonte
externa — é o agente que fez a mudança declarando o que fez.

Isso é o que mantém o mapa **sempre verdadeiro**: ele é atualizado no mesmo instante
em que os arquivos mudam, nunca ficando atrasado à espera de uma reconstrução. É a
manutenção contínua que dá à Propagação um mapa em que ela pode confiar. (No POC,
essa manutenção é um side-effect determinístico de cada disparo — ver a cascata
concreta em `PROPAGATION.md` §6.)

---

## 5. Rastreabilidade é o que roteia

A identidade estável não serve só para *ligar* — ela **roteia**. No POC, cada
cenário declara não só seu código mas também seu **nível** e sua **prioridade**, e
esses atributos, presos à identidade, dizem ao framework o que fazer com aquela
peça:

- o **nível** (`@nivel-unit`, `@nivel-integration`, `@nivel-e2e`) roteia o cenário
  ao runner de teste correto — é a Rastreabilidade dizendo "esta identidade deve
  ter uma encarnação *aqui*";
- a **prioridade** governa o que é exigido e o que bloqueia release.

Ou seja: a mesma cola que conecta também **carrega o contrato** de onde cada peça
precisa existir. É a Rastreabilidade que permite ao pilar de Qualidade perguntar
"onde deveria haver um teste para esta identidade, e ele existe?".

---

## 6. A falha característica: o órfão

Cada pilar tem uma falha que é sua. A da Rastreabilidade é o **órfão**: uma peça
que ficou desconectada da rede — uma árvore fora da trama de fungos. Como o pilar
tem duas metades (§1), o órfão tem duas famílias.

**Órfãos de identidade** (o requisito e suas formas não se ligam):

- **Requisito sem encarnação** — uma spec ou cenário que ninguém implementa; a
  identidade existe mas não tem realização.
- **Encarnação sem requisito** — um teste ou código que não mapeia a nenhuma
  identidade; existe mas ninguém sabe *o que* ele garante.
- **Identidade ausente** — uma peça sem código estável. Este é o órfão mais
  perigoso do lado identidade porque é **invisível aos gates**: sem código, nenhum
  validador consegue cobrá-la ou ligá-la. O POC documenta isto como a raiz dos
  órfãos — um cenário sem código nunca é cobrado, então some do radar sem gerar
  alarme.

**Órfãos de dependência** (o arquivo e o mapa não batem):

- **Dependência invisível** — uma dependência real entre arquivos que o mapa **não**
  registrou. É o análogo, no lado do mapa, da identidade ausente: sem a aresta, a
  Propagação não percorre aquele caminho, e um arquivo que deveria ter ficado stale
  nunca é reconferido. Quebra a onda em silêncio.
- **Aresta morta** — uma aresta no mapa que aponta para um arquivo que não existe
  mais. O mapa mente sobre uma dependência que já não há.

O órfão não é falha de qualidade (o teste até pode passar) nem de arquitetura (o
código até pode respeitar as camadas). É falha de **conexão** — a peça não faz
parte do organismo, ou o mapa não a enxerga. Por isso a Rastreabilidade merece pilar
próprio: sua falha tem natureza própria, e sem um pilar que a cace, ela cresce
silenciosamente até o grafo ser mais buraco que tecido.

---

## 7. Os gates de Rastreabilidade

Como todo pilar, a Rastreabilidade se impõe por **gates** (que são âncoras de
confronto, `QUALITY.md` §2) — só que o que eles medem é *conexão*, não nível de
qualidade:

| gate | pergunta de confronto | falha = |
|---|---|---|
| **identidade presente** | toda peça rastreável tem um código estável? | órfão invisível |
| **identidade única** | cada código identifica uma só coisa (sem colisão)? | ambiguidade |
| **requisito realizado** | toda identidade declarada tem as encarnações que seu contrato exige? (teste comportamental ou de conformidade — salvo o que está `@no-test`) | requisito órfão |
| **encarnação ancorada** | todo teste/código mapeia a uma identidade conhecida? | encarnação órfã |
| **fonte única honrada** | todos os validadores usam a mesma gramática de código? | cola partida |
| **mapa fiel** | toda aresta do mapa aponta para arquivo existente, e toda dependência real está no mapa? | dependência invisível / aresta morta |

Esses gates seguem o mesmo modelo do resto do framework: rodam sobre o **caminho de
impacto** (`PROPAGATION.md` §3 — a Propagação percorre o mapa de dependências que a
Rastreabilidade mantém e reconfere só o que a mudança tocou), emitem **issues
materiais** quando falham (`CONCEPT.md` §5), e têm o mesmo ciclo
de **maturação `informativo → bloqueante`** (`QUALITY.md` §7) — um gate de
rastreabilidade pode nascer report-only ("temos 66 lacunas conhecidas") e ser
promovido a bloqueante quando a cobertura fecha.

A issue de rastreabilidade é tipicamente de kind `violation` (uma peça viola o
contrato de conexão) e nomeia exatamente qual identidade ficou órfã e onde — de
modo que resolver é reconectar (dar código, implementar a encarnação faltante, ou
ancorar a encarnação solta).

---

## 8. Relação com os outros pilares

A Rastreabilidade é a base sobre a qual os outros pilares operam:

- **Habilita o grafo** (`CONCEPT.md`): sem identidade estável, não há como afirmar
  que dois nós são encarnações do mesmo requisito. A cola é o que permite desenhar
  a aresta.
- **Habilita a Qualidade** (`QUALITY.md`): o gate feature→teste só sabe "este
  cenário virou teste?" porque a identidade do cenário aparece no teste. Sem
  Rastreabilidade, os gates de Qualidade medem no vazio.
- **Habilita o anti-drift** (`CONCEPT.md` §2): "todo path/identificador citado
  numa âncora deve existir" é uma checagem de rastreabilidade — seguir o rastro da
  citação até a peça real.
- **Respeita a Estrutura de Projeto** (`STRUCTURE.md`): o mapa de dependências liga
  camadas que a Estrutura declara e respeita as fronteiras da planta. A Estrutura diz
  quais camadas existem; a Rastreabilidade liga as peças dentro e entre elas.

Por isso, num roteiro de amadurecimento, a Rastreabilidade tende a vir cedo: ela é
o solo em que os outros pilares fincam raiz.

---

## 9. Resumo do pilar

- **Rastreabilidade = a cola.** É a trama que conecta as peças e faz do projeto um
  organismo único, não árvores isoladas. O grafo é a estrutura visível; a
  Rastreabilidade é o que a torna confiável, única e sem furos.
- **Duas metades.** *Identidade* (o mesmo requisito através de suas formas) e
  *dependência* (o que se liga a quê). A mesma cola, dois ângulos, materializados no
  mesmo mapa.
- **Identidade contínua via código estável.** Cada requisito carrega um código que
  atravessa suas encarnações (spec → feature → teste → código). A chave é o código,
  não o nome do arquivo; a relação é N:M; a gramática do código é fonte única.
  Quando a string diverge (o drift `LGNN`/`LOGI` do POC), a cola quebra em silêncio.
- **O mapa é o artefato material.** Versionado, lista arquivos e arestas. Cada nó
  carrega sua **versão** (`rev`) e o **carimbo de alteração** (`updated_at`); cada
  aresta, o **carimbo de validação**. A Rastreabilidade mantém o mapa; a Propagação o
  consome. Arestas entram por inferência (motor) + convenção + declaração.
- **A identidade roteia.** O código, com nível e prioridade, carrega o contrato de
  onde cada peça precisa existir — é o que permite à Qualidade saber onde cobrar um
  teste.
- **Falha característica = órfão, em duas famílias.** Identidade (requisito/encarnação
  sem par; identidade ausente) e dependência (dependência invisível; aresta morta). O
  invisível — sem código, ou aresta faltante — é o mais perigoso, porque escapa dos
  gates e quebra a onda em silêncio.
- **Gates de conexão.** Identidade presente, única, realizada, ancorada, fonte única
  honrada, e mapa fiel. Mesmo modelo: caminho de impacto, issues materiais, maturação
  informativo→bloqueante.
- **Base dos outros pilares.** Habilita o grafo, a Qualidade e o anti-drift. Tende
  a amadurecer cedo — é o solo dos demais.

O que o pilar entrega: um projeto onde nada se perde e nada é uma ilha. Cada
requisito pode ser seguido do desejo à realização e de volta; cada peça sabe a que
pertence; e o organismo permanece um só à medida que cresce. É o que impede que um
projeto grande vire um monte de arquivos que ninguém sabe mais como se conectam.

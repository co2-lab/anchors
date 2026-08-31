# Anchors — Pilar de Qualidade

> Este documento define o **pilar de Qualidade** do Anchors. Ele pressupõe o
> mecanismo geral estabelecido em [`CONCEPT.md`](./CONCEPT.md) — âncora, grafo,
> sincronia incremental, issues, e a separação vivo/histórico — e o especializa
> para um propósito: fazer com que o trabalho não seja apenas bem feito, mas que
> esse "bem feito" seja **medido**.
>
> É teoria de framework, já validada por uma ferramenta que implementa o
> Anchors e por uma prova de conceito real (incompleta e imatura de
> propósito — prova o mecanismo, não a completude). Ambas aparecem apenas
> como instâncias.

---

## 1. Por que qualidade é um pilar

Frameworks spec-first já existem. Eles garantem que você constrói a coisa certa
(a spec guia) e, com anti-drift, que a âncora não mente sobre o código. O que
eles não têm — e o que o Anchors adiciona como pilar de primeira classe — é a
garantia de que o trabalho atingiu um **nível medido de qualidade** antes de
avançar.

Nenhum projeto vai para produção sem um certo nível de qualidade ou maturidade.
Mas "qualidade" tratada como cuidado subjetivo do desenvolvedor no dia é frágil:
é uma sensação, não sobrevive entre sessões, e some quando o autor troca. Um
projeto assim pode funcionar e ainda ser inefetivo — fica indo e voltando em
falhas.

O pilar de Qualidade transforma essa sensação em uma **propriedade medida e
versionada do artefato**, que qualquer sessão futura pode reavaliar. É por isso
que ele é nomeado pilar isoladamente: para reforçar que qualidade não é diluída
nos outros mecanismos nem tratada com descuido. Sem esse pilar bem estruturado, o
projeto não é maduro (ver [`CONCEPT.md` §1.1](./CONCEPT.md)).

---

## 2. Um gate é uma âncora de qualidade

Um **gate de qualidade** é uma âncora (`CONCEPT.md` §2) cujo alvo é confrontado
contra um **limiar medido**, não contra uma descrição.

- **Guia (ida):** o gate declara o alvo de qualidade — "features deste arquivo
  devem ter testes que passam", "cobertura ≥ 80%", "zero violações de camada",
  "toda tela tem `.feature`".
- **Confronta (volta):** o gate **executa a medição** — roda os testes, calcula a
  cobertura, roda o linter, valida a estrutura — e emite uma issue se o resultado
  está abaixo do limiar.

A diferença entre um gate e uma âncora comum é o tipo de pergunta do confronto.
Uma spec pergunta *"o código faz o que descrevi?"* (correspondência). Um gate
pergunta *"o alvo atingiu o nível medido?"* (limiar). Mecanicamente, é o mesmo
confronto: o resultado vira o `verdict` no carimbo da aresta, e a falha vira uma
issue de kind `violation` (`CONCEPT.md` §5). O pilar de Qualidade **não inventa
mecânica nova** — reusa a âncora, o grafo e a issue.

### A dupla saída: issue e bloqueio

Quando um gate reprova, ele produz **duas** coisas, não uma:

- **Uma issue** — o registro material da falha, para correção (`CONCEPT.md` §5).
  Rastreável, datada, append-only.
- **Um bloqueio** — a barreira que impede o avanço: nada sobe enquanto a falha
  existe.

E há uma **assimetria deliberada** entre as duas: a **issue é inegociável**, o
**bloqueio é negociável**. A dessincronia fica registrada de qualquer jeito — a
verdade não se apaga. Mas o bloqueio tem opt-out (`--no-block`, §7): você pode
*escolher deixar passar*, conscientemente, sem apagar a issue. O trabalho avança, o
débito fica visível. É o mesmo espírito do `@no-test` da Spec e do "done não pode
mentir" das issues (`CONCEPT.md` §5): **pode-se optar por não bloquear, nunca por
não registrar.**

O que o pilar **adiciona** ao mecanismo comum é:

1. **A ligação com medidores externos** — como plugar test runners, linters e
   validadores da linguagem como gates (§4).
2. **O pipeline de qualidade** — a orquestração do conjunto de gates (§5).
3. **A completude** — fechar todas as pontas: garantir que nada escapa da validação
   (§5.1).
4. **A agregação** — qualidade como o conjunto de vereditos, não um gate isolado (§6).
5. **A maturação do gate** — o ciclo `informativo → bloqueante`, e o opt-out pontual
   `--no-block` — que ligam o pilar ao conceito de maturidade (§7).

---

## 3. "Gate de qualidade" é um conjunto, não um único gate

"Gate de qualidade" é um nome guarda-chuva. Na prática é um **conjunto de tools e
prompts que medem e avaliam** o trabalho. **A qualidade emerge do conjunto** —
cada gate mede uma coisa específica, e nenhum gate isolado *é* a qualidade.

O pilar se organiza em **categorias de gate**, cada uma confrontando um tipo de
âncora diferente. (Não confundir com as *camadas de arquitetura* do projeto, que
são domínio da Estrutura de Projeto, `STRUCTURE.md` — aqui "categoria" é o tipo de
coisa que o gate mede.) A prova de conceito real já exercita quatro categorias,
que servem de referência — não de lista fechada:

| categoria de gate | o que o gate mede | tipo de medidor |
|---|---|---|
| **Arquitetura** | o alvo respeita as regras de camada/estrutura? (régua da Estrutura, `STRUCTURE.md`) | determinístico (analisador estático) |
| **Artefato** | as âncoras existem e estão completas? cobrem o que a spec declara? | determinístico (validador de spec/feature) |
| **Execução** | os testes derivados das features passam? | determinístico (test runner nativo) |
| **Rastreabilidade** | cada cenário de feature virou teste no runner certo? (régua da Rastreabilidade) | determinístico (validador de cobertura cenário↔teste) |

### Dois tipos de medidor: determinístico e julgamento por IA

O que muda entre gates não é só *o que* medem — é **como** medem. Há dois tipos de
medidor, e o pilar acomoda ambos **com o mesmo modelo** (dupla saída issue+bloqueio,
§2; força da aresta; maturação):

- **Determinístico** — roda um comando e lê o código de saída. Objetivo, barato,
  binário. Mede o que se pode *computar*: lint passou, spellcheck limpo, cobertura ≥
  X%, build verde, spec tem as seções obrigatórias, cada cenário virou teste.
- **Julgamento por IA** — um agente confronta o alvo e emite o mesmo tipo de
  veredito. Mede o que *nenhum script computa*: a spec está **completa** e
  coerente? o código está **legível**? a solução respeita o espírito da arquitetura,
  não só a letra? É o **gate de review por IA** — o exemplar do lado qualitativo.

A diferença é só *quem mede*. A força da aresta (bloqueante/informativa), a dupla
saída (issue + bloqueio) e o kind da issue são idênticos nos dois. Um gate de review
por IA que reprova gera a mesma issue e o mesmo bloqueio (com `--no-block`) que um
lint que reprova — muda a natureza do medidor, não a mecânica do gate.

> **E a conformidade?** O teste de conformidade do regime declarativo (`SPEC.md` §6 —
> "o recurso declarado bate com a forma?") **não é um terceiro tipo de medidor**. Ele
> se distribui pelos dois acima conforme o *critério esteja fixado*: comparar o
> `EXPOSE 8080` do Dockerfile com a spec, ou o GSI declarado no Terraform com a query
> que a spec pede, é **determinístico** (forma ↔ forma, computável). Só vira
> **julgamento** quando o critério é vago ("3 réplicas satisfazem 'alta
> disponibilidade'?" — a menos que a spec fixe "HA = ≥3"). Conformidade é *o que* o
> gate confronta (uma forma contra outra), não *como* mede.

### O eixo ortogonal: gate local vs. gate com estado externo

Independente do tipo de medidor, há um segundo eixo — **de onde o gate tira a
verdade contra a qual confronta**:

- **Gate local** — precisa só do **repositório**. Lint, spellcheck, validate-spec,
  comparar Dockerfile↔spec. Barato, reproduzível, roda em qualquer lugar.
- **Gate com estado externo** — precisa olhar o **mundo fora do repo**. "O GSI
  declarado no Terraform existe *de fato* na nuvem?" exige consultar a AWS. É o
  confronto *declarado ↔ realizado* — a base da detecção de **drift** de infra.

Este eixo importa porque gates com estado externo são mais caros, menos
reproduzíveis e podem falhar por razões fora do código (a nuvem mudou, a credencial
expirou). Um projeto tipicamente roda os locais sempre (baratos) e os de estado
externo na promoção ou sob demanda (§8). O eixo é ortogonal ao tipo de medidor: um
gate de estado externo ainda é determinístico ou julgamento — só que a verdade que
ele confronta mora fora do repositório.

### Catálogo de gates comuns (referência, não lista fechada)

| gate | mede | medidor | fonte |
|---|---|---|---|
| **lint** | estilo/padrões de código além do compilador | determinístico | local |
| **spellcheck** | ortografia em código/docs | determinístico | local |
| **format** | formatação canônica | determinístico | local |
| **typecheck / build** | compila, tipos batem | determinístico | local |
| **validate-spec** | spec presente, completa, spec-first (`SPEC.md` §8) | determinístico | local |
| **validate-feature** | cobertura spec→feature; feature válida | determinístico | local |
| **test / coverage** | testes passam; cobertura ≥ limiar (`§4`) | determinístico | local |
| **traceability** | cada requisito virou encarnação (`TRACEABILITY.md` §7) | determinístico | local |
| **architecture** | respeita as fronteiras de camada (régua da `STRUCTURE.md`) | determinístico | local |
| **conformidade (forma)** | o declarado (Dockerfile/`.tf`) bate com a spec? | determinístico | local |
| **conformidade (drift)** | o recurso declarado existe *de fato* no ambiente? | determinístico | **estado externo** |
| **review por IA** | coerência, legibilidade, spec completa, aderência ao espírito | **julgamento por IA** | local |
| **best-practices por IA** | segue os patterns/anti-patterns do domínio | **julgamento por IA** | local |

| **obrigação transversal** | o dever que um artefato contrai com um lugar que ele não conhece (§3.1) | determinístico | local |

Um projeto declara os gates que quer (§4, *ligar um gate é declarativo*) — esta
tabela é ponto de partida, ancorada no que o POC exercita, não uma lista fixa. O que
importa é que **todo tipo de coisa que precisa de validação tem um gate cobrindo-a**
(§5.1) — e onde um script não alcança, o gate de review por IA fecha a ponta.

---

### 3.1 A obrigação transversal: o dever que mora fora da unidade

Quase toda validação olha para DENTRO de uma unidade: a spec descreve o que ela faz, o
teste prova, o gate confronta. Há uma classe de erro que nenhum desses vê, porque o dever
não é da unidade — é do **sistema**:

> **Todo artefato que satisfaz P deve aparecer em Q.**

O artefato novo contrai um dever com um lugar que ele não conhece, e o dono desse lugar
não sabe que o artefato nasceu. Nenhum dos dois lados tem como perceber sozinho.

| domínio | P (o gatilho) | Q (o dever) |
|---|---|---|
| privacidade (LGPD/GDPR) | modelo que carrega dado pessoal | script de exclusão e de exportação de conta |
| i18n | string exibida ao usuário | arquivo de traduções |
| acessibilidade | componente interativo | ter rótulo acessível |
| observabilidade | operação de negócio | emitir métrica/log |
| auditoria regulatória | operação financeira | trilha imutável |

**O caso que originou o conceito:** um modelo de dado novo — o campo de texto livre onde
o usuário anota o que quiser, portanto o de maior risco de dado sensível — ficou de fora
do script de exclusão de conta. Dado pessoal que nunca seria apagado. O projeto **já
tinha sido mordido por isso antes** (o próprio script carrega um comentário chamando de
"violação silenciosa"), e ainda assim repetiu: nove gates verdes, e nenhum viu, porque
todos olhavam para dentro da unidade.

#### O gatilho precisa ser DECLARADO, não deduzido

A tentação é deduzir da camada: *"todo modelo de dado vai no script de exclusão"*.
Medição num projeto real: **38 de 50 modelos não estavam lá — e a maioria
legitimamente** (dado compartilhado entre pessoas, dado com expiração automática,
métrica agregada sem identificação). A régua deduzida erraria 76% das vezes, e um gate
assim é desligado no primeiro dia.

Por isso o gatilho é um **atributo declarado no cabeçalho** do artefato (`carries: pii`),
com vocabulário do projeto. O Anchors não sabe o que é dado pessoal; sabe que existem
obrigações, e oferece como declará-las (`obligations:`) e como cobrá-las.

#### A exceção é honesta (§5.1 do CONCEPT)

Um artefato pode se eximir — com **motivo obrigatório**. Dispensa sem justificativa é
tratada como ausente, e reprova. É o que impede o time de desligar o gate quando aparece
a primeira exceção legítima: ela fica registrada e localizada, em vez de o rigor inteiro
cair.

## 4. A ponte central: feature → teste → gate

A ligação que fecha o loop spec-first com a qualidade é a cadeia
**spec → feature → teste**. É o núcleo do pilar.

```
spec  ──gera──►  feature (cenários)  ──implementada como──►  teste (runner nativo)
 (o quê)          (comportamento               (jest / go test / pytest / …)
                   esperado, executável)              │
                                                é um GATE de qualidade
```

A **feature não é só documento** — ela é a especificação executável do
comportamento, que **deve ser implementada como testes** nos frameworks de teste
nativos da linguagem. Esses testes **são um gate de qualidade**. O framework deve
oferecer a forma de **ligar esses testes como gate** e de montar um **pipeline de
qualidade** que os inclui.

> **Teste abrange conformidade, não só comportamento.** A cadeia acima é a do
> regime **comportamental** (`SPEC.md` §6) — cenário Gherkin → teste executável. Uma
> spec **declarativa** (Terraform, Dockerfile, K8s) é testada por **conformidade**:
> o recurso declarado bate com a forma especificada? É um teste de tipo diferente,
> mas entra no mesmo pipeline como gate. Toda spec é testável por default; uma spec
> que o usuário não quer testar é marcada **`@no-test`** e o gate a dispensa (ver
> `TRACEABILITY.md`).

### A relação feature↔teste é N:M, mediada por um código de cenário

A ligação não pode ser por convenção de nome de arquivo (frágil) nem 1:1
(insuficiente). O padrão provado na prática é mais robusto: **cada cenário carrega
um código estável**, e esse código é a chave que amarra os quatro artefatos.

```
spec.md   declara   →  CODEX-S01 (um estado/regra/requisito)
feature   marca     →  @CODEX-S01  @nivel-integration  @P2
                                   │  (o nível roteia o cenário ao runner certo)
                    ┌──────────────┴───────────────┐
                    ▼                              ▼
          @nivel-unit / @nivel-integration    @nivel-e2e
          teste unitário/componente:          teste de fluxo:
          it('CODEX-S01: …')  no runner nativo   arquivo E2E nomeado CODEX-S01
```

Propriedades desse modelo (destiladas do POC):

- **Um cenário pode declarar múltiplos níveis** (`@nivel-unit @nivel-e2e`) → gera
  teste em mais de um runner. Daí a relação ser N:M, não 1:1.
- **O nível roteia o cenário ao medidor** — cada nível de teste corresponde a um
  runner: regra pura → teste unitário; componente isolado → teste de componente;
  fluxo entre telas → teste E2E.
- **A chave de ligação é o código**, não o nome do arquivo. "O cenário virou
  teste?" se responde procurando o código no artefato de teste do nível
  declarado. Um cenário **sem código é invisível ao gate** — nunca é cobrado, e
  isso é a raiz de órfãos: o código é o que torna a rastreabilidade verificável.
- **O código é fonte única.** Um único módulo define a gramática do código de
  cenário, importado por todos os validadores (de todas as linguagens). É o ponto
  de acoplamento que mantém os gates de linguagens diferentes falando a mesma
  língua.

### Ligar um gate ao framework é declarativo

Como o teste é implementado no runner nativo (que varia por linguagem), o Anchors
define o **conceito** de gate de teste, mas a **ligação concreta** é declarada por
projeto. O framework oferece o mecanismo; o projeto declara o gate:

- **qual comando** roda o medidor (o runner, o linter, o validador);
- **em que escopo** ele roda (qual pacote/diretório/linguagem);
- **como o resultado é lido** (código de saída, relatório de cobertura, contagem
  de violações);
- **qual o limiar** e se o gate é **bloqueante ou informativo**.

O POC hoje codifica essa lista **imperativamente** em dois lugares (o workflow de
CI e um script agregador) — não há um manifesto declarativo único. **Esse é
exatamente o buraco que o pilar preenche:** o Anchors deve oferecer o manifesto
declarativo de gates, para que "estes comandos são os gates" seja uma âncora
versionada e não código espalhado.

---

## 5. O pipeline de qualidade

O **pipeline de qualidade** é a orquestração do conjunto de gates sobre um alvo.
Ele transforma a spec em qualidade medida:

```
                          alvo alterado
                               │
        caminho de impacto: gates a rodar (via Propagação, PROPAGATION §3)
                               │
      ┌───────────┬────────────┼────────────┬───────────────┐
      ▼           ▼            ▼             ▼               ▼
 arquitetura  artefato    execução      cobertura     (gates do projeto)
 (camadas)    (spec/feat) (testes)      (cenário↔teste)
      │           │            │             │               │
      └───────────┴────────────┼─────────────┴───────────────┘
                               ▼
                    agregação = perfil de qualidade
                               │
              algum gate bloqueante falhou? → issue(s) violation
                               │
                    (senão) alvo passa o pipeline
```

Dois pontos herdados do mecanismo comum tornam o pipeline barato e honesto:

- **Incremental** (`PROPAGATION.md` §3): o pipeline não roda todos os gates sobre
  todo o projeto a cada mudança. A **análise de impacto** da Propagação, sobre o
  mapa de dependências da Rastreabilidade, dá o **caminho de impacto** — só os gates
  cujas arestas ficaram stale entram. É o que torna o rigor sustentável num projeto
  que cresce: caro-demais nunca é feito, e não-feito vira drift.
- **Issues materiais** (`CONCEPT.md` §5): a falha de um gate bloqueante não é um
  log que some — é uma issue `violation`, arquivo em `issues/todo/`, imutável,
  que só sai quando a dessincronia é resolvida (corrige o alvo ou baixa o limiar
  conscientemente, o que é atualizar a âncora do gate).

---

## 5.1 Fechar todas as pontas: nada sem validação

A qualidade emerge do *conjunto* de gates (§3) — mas isso só é verdade se o conjunto
**cobre todas as pontas**. Se existe uma âncora, uma aresta ou um requisito que
**nenhum gate confronta**, há um buraco por onde trabalho não-validado passa. Fechar
todas as pontas é a meta: cada coisa que existe tem alguém que a valida.

### A falha característica: o buraco de cobertura

A falha própria da Qualidade **não é "um gate reprovou"** — isso é o sistema
funcionando (o gate fez seu trabalho). A falha é o oposto: **existe algo que nenhum
gate cobre** — a *ponta aberta*. E ela é mais perigosa que uma reprovação, porque é
**silenciosa**: ninguém reprova, mas ninguém validou. O trabalho parece são só
porque não foi olhado.

É uma falha distinta das outras:

- não é o **órfão** da Rastreabilidade (a peça pode estar perfeitamente conectada e
  ainda não ter um gate que meça sua qualidade);
- não é a **onda incompleta** da Propagação (a onda pode ter alcançado a peça, mas
  não havia gate na ponta para confrontá-la).

É **ausência de cobertura**. A porta *legítima* de sair da cobertura é o **opt-out
honesto** (`CONCEPT.md` §5.1) — a dispensa explícita, registrada e justificada
(`@no-test`, `--no-block`, maturação informativa). O buraco de cobertura é a saída
*ilegítima* — o oposto ponto a ponto: implícito, não-registrado, invisível; escapar
sem ninguém decidir. É essa comparação que dá ao meta-gate seu critério: uma ponta
sem gate só é aceitável se houver um opt-out honesto a cobrindo; senão, é buraco.

### Gates de artefato e gates de fluxo

Fechar todas as pontas exige distinguir *o que* se valida. Há duas famílias de gate:

- **Gates de artefato** — validam uma **âncora** (o catálogo do §3: lint, validate-
  spec, test, architecture…). "Este artefato está bom?"
- **Gates de fluxo** — validam um **passo de execução ou propagação**: a entrada, a
  saída e os passos transversais da cascata. "Este *passo* foi validado?"

O miolo da cascata (`spec → código → feature → teste`) é coberto por gates de
artefato. Mas as **bordas** (entrada, saída) e os **passos transversais** precisam
de gates de fluxo — senão são pontas abertas. Cada passo de execução tem o seu:

| gate de fluxo | valida o passo | medidor | régua |
|---|---|---|---|
| **validação de plano** | o plano recém-criado é coerente, completo, exequível? | review de IA | Planejamento |
| **frescor da doc** | a doc reflete a fonte agora? (a sincronia dispara; a IA confronta) | review de IA | — (âncora de consumo) |
| **mapa-fiel (incremental)** | a cada execução de agente: as arestas que o ato deveria criar estão no mapa? (script que percorre o trecho tocado) | determinístico | Rastreabilidade |

Padrão observado: os gates de fluxo das **bordas** (plano, doc) tendem a **review de
IA** — validar intenção e prosa é julgamento, não computação; os **transversais**
(mapa) tendem a determinístico. Todos são gates de Qualidade com régua do pilar
dono, como os de artefato — mesma mecânica (dupla saída, maturação), nenhuma
novidade. A auditoria que os encontra é o próprio meta-gate rodando sobre o fluxo
(exemplo em `simulation/larder/cobertura.md`). Nem toda ponta, porém, fecha por um
gate no momento do passo — algumas são **sistêmicas** e só o validador de saúde
(§5.2) as pega.

> **O opt-out não precisa de gate.** Cogitou-se um gate que auditasse opt-outs
> (`@no-test`/`--no-block`) velhos. Mas o opt-out é decisão *do usuário* — ele valida
> se ainda faz sentido e reage quando não for como esperava. Um gate que o cobra
> seria o framework arbitrando uma decisão que é do operador. O opt-out já é honesto
> por ser explícito, registrado e datado; o rastro basta para o usuário revisar
> *quando quiser* — o framework não cobra.

### O meta-gate de completude

Por isso a Qualidade precisa de um gate que confronta a **própria completude do
conjunto**: *"toda âncora/aresta que deveria ter um gate tem um?"*. Ele mede
cobertura de *gates*, não de código — caça as pontas abertas antes que virem o
caminho por onde o drift entra. Quando acha uma ponta aberta, produz a dupla saída
como qualquer gate (§2): uma issue (registrar a ponta) e um bloqueio (com opt-out).
Fechar a ponta é ou ligar um gate que a cubra, ou marcá-la explicitamente como
dispensada (`@no-test` / sem exigência) — nunca deixá-la aberta em silêncio.

### A regra sem régua: toda regra nova pede um gate

O §5.1 caça a **âncora** sem gate. Há um segundo buraco, do lado oposto: a **regra**
sem gate — um dever que alguém declarou (no guide, na spec, num combinado verbal) e
que nenhum medidor confronta. Enquanto ela vive só em prosa, ela vale enquanto alguém
lembra; e lembrar não escala.

**Dois momentos obrigam o agente a perguntar "isto cabe num gate?":**

1. **Quando uma regra é ADICIONADA** — o usuário escreve uma regra num guide, numa
   spec, ou pede um novo padrão. Antes de considerar o trabalho feito, o agente
   propõe a régua que a cobre.
2. **Quando uma regra é PEGA FALHANDO** — apareceu uma violação de algo que já era
   regra. A correção da violação é metade do trabalho; a outra metade é a régua que
   impede a reincidência. Uma regra que foi violada uma vez prova, por construção,
   que a memória humana não a sustenta.

**A ordem de preferência do medidor é a mesma do §3:**

| tentativa | quando serve | custo |
|---|---|---|
| **1. Determinístico declarativo** | a regra é reconhecível por padrão no artefato (import proibido, campo obrigatório, nome fora de convenção) — cabe numa entrada de configuração do framework, sem código novo | quase zero |
| **2. Determinístico com código** | a regra exige atravessar estrutura (resolver uma aresta, comparar dois artefatos) — pede um gate próprio | médio |
| **3. Julgamento por IA** | a regra é semântica e não tem forma fixa ("o comentário descreve o que o código faz") | alto, e não determinístico |

**Só se as três não servirem** é que a regra fica em prosa — e aí ela é declarada
como dívida explícita, não esquecida em silêncio.

O caso que originou esta seção: um projeto tinha o combinado "todo sheet usa a lib de
bottom-sheet, nunca o modal nativo da plataforma". Não era gate. Uma varredura achou
**20 usos** do modal nativo em 11 arquivos, invisíveis por meses — e o custo não era
estético: as duas famílias de modal, aninhadas, produzem falha SILENCIOSA (o sheet
simplesmente não abre, sem erro). A regra existia; a régua, não. Virou uma entrada
declarativa de fronteira em minutos — a tentativa 1 da tabela — com duas exceções
marcadas na própria linha do código.

> **A regra de bolso:** ao terminar de corrigir uma violação, pergunte *"o que impede
> isto de voltar amanhã?"*. Se a resposta é "alguém vai lembrar", o trabalho não
> terminou.

---

## 5.2 O validador de saúde do ecossistema

Os gates — de artefato e de fluxo — cobrem pontas **locais**: um artefato, um passo.
Mas há pontas que são **sistêmicas** — não moram num passo, moram na *integração
entre os pilares e no estado do ecossistema como um todo*:

- um **laço issue→plano** que nunca fechou (o plano concluiu mas a aresta que gerou
  a issue continua stale);
- a **integridade global do mapa** (não "esta aresta que acabei de criar", mas "o
  grafo inteiro é consistente, sem ilhas nem ciclos espúrios");
- um **pilar frouxo** (a Estrutura declara camadas que nenhuma spec usa; há specs sem
  nenhum gate; a Rastreabilidade tem regiões sem identidade);
- **gates que deveriam existir e não foram declarados** (o meta-gate de completude,
  §5.1, no nível do ecossistema).

Essas pontas não fecham por um gate no momento de um passo — nenhum passo isolado as
vê. Elas precisam de uma **visão global**. Por isso o Anchors tem o **validador de
saúde do ecossistema**: um meta-gate elevado do nível do *passo* para o nível do
*framework inteiro*. Ele confronta *"os pilares estão íntegros e integrados? o
projeto está maduro?"* — varrendo o estado de todos os pilares e as pontas
sistêmicas de uma vez.

Diferente dos outros gates, ele:

- roda **periodicamente / sob demanda** (visão global), não a cada alteração
  (os gates incrementais já cobrem o local);
- mede **integração e maturidade**, não um alvo — é a **encarnação executável da
  maturidade** (`CONCEPT.md` §1.1). "Maturidade = presença e vigor dos pilares" era
  descritivo; o validador de saúde é o que *mede* isso e o torna um veredito;
- **apresenta e registra, mas não bloqueia.** É um caso especial na mecânica de gate
  (§2): como roda sob demanda e fora do caminho de um merge, ele não *trava* nada —
  **reporta** o estado de saúde e **abre issues** para as pontas sistêmicas que
  encontra (reusando o registro material). O bloqueio continua sendo dos gates
  locais, no momento do passo. É diagnóstico global, fiel ao "detecta e apresenta,
  não arbitra" (`CONCEPT.md` §2): quem decide agir sobre o veredito é o operador.

O validador de saúde é o que **fecha as últimas pontas** — as que nenhum gate local
alcança. É o guardião de que o framework aplicado a um projeto continua um organismo
íntegro, não um conjunto de pilares que funcionam isolados mas não se integram.

### A forma concreta: `doctor` / `status` no CLI

Materializa-se num **comando de CLI** — algo como `anchors doctor`, `status` ou
`validate` — que roda a varredura do ecossistema e reporta: problemas encontrados,
estado dos pilares, pontas sistêmicas abertas, alertas, e o nível de maturidade
atual. É o "raio-X" do framework aplicado ao projeto, invocado sob demanda por quem
opera. (Uma ferramenta que implementa o Anchors tem o precedente: um `/doctor` que
verifica o índice e oferece reindex; o validador de saúde generaliza isso para
*todos* os pilares.) Como o watcher
(`PROPAGATION.md` §6), o comando é *uma* materialização — o Anchors define o validador
de saúde em abstrato; o executor o expõe como CLI.

---

## 6. Agregação: qualidade é o conjunto dos vereditos

A qualidade de um alvo **não é um único número**. É o **conjunto dos vereditos**
dos seus gates — um perfil por categoria de gate:

```
login.tsx  →  arquitetura ✔   artefato ✔   execução ✔   cobertura ⚠(72%)
```

O Anchors modela qualidade como esse perfil, não como um score 0–100. A decisão de
promoção ("este alvo pode avançar?") é **todos os gates bloqueantes passaram** —
gates informativos entram no perfil mas não travam. Um score agregado, se
desejado, é derivável do perfil, mas o perfil é a fonte de verdade: ele diz
*exatamente o que* está abaixo, não só *quanto*.

Isso preserva a regra do pilar: cada gate mede uma coisa específica; a qualidade
é dada pelo conjunto. Nenhum gate isolado carrega o veredito de qualidade sozinho.

---

## 7. Maturação do gate: `informativo → bloqueante`

Este é o elo entre o pilar de Qualidade e o conceito guarda-chuva de maturidade.

Um gate não nasce bloqueante. Um projeto real quase nunca cumpre, no dia um, o
limiar que quer atingir — impor o gate como bloqueante imediatamente pararia o
projeto. Por isso todo gate tem um **estado de maturação**:

| estado | comportamento | quando |
|---|---|---|
| **informativo** (report-only) | mede e reporta, mas **não bloqueia** | o projeto ainda não cumpre o limiar; está no caminho |
| **bloqueante** | mede e **impede o avanço** se abaixo do limiar | o projeto alcançou o limiar; agora ele é defendido |

A **promoção** de informativo para bloqueante é uma **decisão explícita e
versionada** — muda-se a declaração do gate no manifesto, revisável em PR. Ela
acontece quando a realidade do projeto alcança o limiar do gate ("a fase que zera
as lacunas fechou"). A partir daí, o gate deixa de ser uma meta e passa a ser uma
fronteira defendida: nada regride abaixo dele.

Este ciclo **é** a maturidade em ação, no nível do gate:

- Um gate **informativo** é uma promessa ("queremos chegar aqui").
- Um gate **bloqueante** é uma garantia ("não descemos daqui").
- Um projeto **maduro** é aquele cujos gates críticos foram promovidos a
  bloqueantes — os pilares estão implementados e defendidos.
- Um projeto **imaturo** ainda tem seus gates em report-only — mede, sabe onde
  está, mas ainda não defende o patamar.

O POC demonstra isso rodando: seus gates nascem em modo report-only e são
promovidos a bloqueantes um a um, conforme cada frente amadurece. A POC ser
"incompleta e imatura" não é um defeito a esconder — é a **demonstração viva de
que o pilar tem estágios de maturação**. Um projeto acabado teria todos os gates
já bloqueantes; um POC os tem em transição.

### O cold start da adoção: maturação em lotes

Há uma assimetria entre nascer com o Anchors e adotá-lo depois, e ela importa
sobretudo para o **gate de julgamento por IA** (§5.2), cujo medidor é caro:

- **Projeto que nasce com o Anchors** amadurece de graça: cada artefato é confrontado
  e **carimbado** (`PROPAGATION.md` §3) quando é criado. A cobertura de julgamento
  cresce junto com o código; nunca há dívida acumulada.
- **Projeto que adota o Anchors** herda um mapa onde *tudo* nasce sem carimbo. Uma
  auditoria de julgamento completa pode custar milhões de tokens de IA — inviável de
  uma vez. Isto **não é um defeito do pilar**; é o mesmo cold start que a maturação
  informativo→bloqueante já pressupõe, agora no eixo da *cobertura*, não do *limiar*.

A estratégia de adoção é **maturar a cobertura em lotes**, e o próprio mecanismo do
pilar a sustenta:

1. **Corte estrutural antes de medir.** Redundância de governança (vários guides
   regendo o mesmo conjunto) e guides sem governança inflam o trabalho sem gerar
   valor. Afinar a Estrutura (as tags) corta o total na raiz, sem gastar medição.
2. **Batch por régua, não por alvo.** O custo do julgamento é dominado por *ler a
   régua*, não o alvo. Confrontar muitos alvos de um mesmo guide numa passada
   amortiza esse custo fixo. Paralelizar, se for o caso, é **por régua** — nunca por
   alvo (que faria reler a régua a cada um).
3. **O carimbo é o cursor.** A primeira passada não precisa terminar de uma vez: o
   carimbo por aresta registra o que já foi confrontado, então retomar pula o que tem
   veredito válido na rev atual. A auditoria vira incremental e retomável — e, uma vez
   feita, só o *drift* (o que muda depois) volta à fila.

A honestidade da §7 vale aqui: uma cobertura parcial deve ser **reportada como
parcial** (o validador de saúde do ecossistema, §5.1, mostra o que ainda não foi
confrontado), nunca disfarçada de auditoria completa.

### `--no-block`: o override pontual, não a política

A maturação é a **política permanente** do gate — informativo ou bloqueante, valendo
para todo o projeto até ser mudada em PR. O **`--no-block` é diferente**: é um
override **pontual**, de uma execução, sobre um gate que *é* bloqueante. "Sei que
isto bloquearia; deixa passar desta vez."

A diferença importa:

| | maturação (informativo) | `--no-block` (pontual) |
|---|---|---|
| escopo | permanente, todo o projeto | uma execução |
| onde vive | declaração do gate, versionada | invocação, registrada |
| a issue | nasce (é registro, não bloqueio) | nasce igual |

Nos dois casos a **issue nasce** — a dupla saída (§2) garante que a verdade é sempre
registrada. O que o `--no-block` dispensa é só o bloqueio daquela vez, e de forma
**explícita e datada** com o *porquê* (não um silêncio). Tanto o `--no-block` quanto
a maturação informativa são instâncias do **opt-out honesto** (`CONCEPT.md` §5.1) — a
dispensa explícita, registrada e justificada, irmã do `@no-test` da Spec: deixar
passar é uma decisão visível e auditável, nunca um buraco. Baixar um gate
permanentemente é maturação (mexe na política); pular o bloqueio uma vez é
`--no-block` (não mexe na política, e o gate continua bloqueante na próxima).

### As classes que passam por todos os outros gates

Um gate confronta uma promessa contra o que existe. Isso cobre a maioria dos defeitos —
mas há classes em que **todas as peças existem e se referenciam**, e ainda assim o
produto está errado. São elas que motivam os gates mais recentes, e cada uma tem uma
assinatura própria:

| classe | por que os outros gates não pegam | o que a confronta |
|---|---|---|
| **invariante/interação** | cada regra tem teste verde; o defeito está entre elas | invariante provada em CICLO FECHADO (produtor → consumidor) |
| **contrato de fronteira** | as irmãs são consistentes, uma não é | assimetria entre funções irmãs |
| **obrigação transversal** | a obrigação mora FORA da unidade | declarar `obligations:`, não deduzir |
| **custo/escala** | o teste roda com 3 registros; a produção com 3 mil | promessa do NOME vs. completude do retorno |
| **teste que não prova** | 100% verde e 100% coberto sem verificar nada | mutação: apague a linha e veja se cai |
| **ambiguidade não resolvida** | ninguém tomou a decisão | `## Decisões em aberto` vazia como condição |

Duas observações que valem para todas:

**A assimetria é a prova mais barata que existe.** Não exige entender o domínio: se três
funções irmãs recebem o mesmo parâmetro e duas o guardam, a terceira é esquecimento e não
decisão. As duas leituras não podem estar certas ao mesmo tempo. O mesmo raciocínio pega
a função que não pagina num módulo onde outras catorze paginam — o padrão é conhecido ali.

**Custo/escala e teste-que-não-prova se reforçam.** Medido: uma função que devolve no
máximo 100 comissões é consumida por um lambda agendado cujo teste MOCKA a função. O
limite nunca aparece em teste algum — e não apareceria mesmo com cobertura de 100%.

---

### A régua contra si mesma: quando duas fontes discordam

O modo de falha mais frequente medido em cinco rodadas de um E2E real não foi regra
ausente — foi **duas fontes da própria régua dizendo coisas opostas**. Quatro casos, e em
todos o agente teve de escolher qual obedecer:

| fonte A | fonte B | quem estava certo |
|---|---|---|
| guide: "handler não consome o DAO" | `boundaries:` permite, e proíbe o contrário | a config (64 handlers consomem o DAO, 0 consomem repository) |
| preset emite "Auth/Acesso" | `rule_types` declara "Permissões e Acesso" | a config (50 specs vizinhas) |
| guide se declara "taxonomia canônica" | `rule_types` tem letra que o guide não lista | a config |
| plano manda editar `resource.ts` | a estrutura já separou os modelos | a estrutura |

O padrão: **o guide é PROSA, a config é DECLARAÇÃO, e ninguém confronta os dois.** Um
projeto real tinha 10 guides com proibição escrita em texto e 16 fronteiras declaradas em
`boundaries:` — dois sistemas de regra sobre o mesmo assunto, sem ligação.

Isso custa mais que uma regra ausente. Regra ausente produz pergunta; regra contraditória
produz **desobediência sem consciência** — quem obedece a fonte errada não sabe que
desobedeceu, e o gate fica verde porque confronta a fonte que ele conhece.

Três consequências para o desenho:

1. **A config é a fonte de verdade sobre o que ela declara.** Onde `boundaries:`,
   `rule_types` ou `layers:` decidem, o guide EXPLICA e não contradiz — e o que o guide
   ensina deve ser derivado da declaração, não escrito ao lado dela. O `anchors new` já
   faz isso com os títulos de seção (o preset traduz pelo `rule_types` do projeto).
2. **Onde a régua não pode ser derivada, ela precisa dizer quem desempata.** A frase que
   faltava nos quatro casos é a mesma: "se isto divergir do `anchors.yaml`, o YAML vence".
3. **Contradição detectável é gate.** Três dos quatro casos eram confrontáveis por
   máquina: um guide que proíbe o que o `boundaries:` permite, um preset que emite seção
   que o `rule_types` não reconhece, um guide que se declara canônico sobre um vocabulário
   que a config estende.

---

### Dois níveis de review: a parcial e o conjunto

O review fecha o ciclo, e ele tem DOIS níveis — não por rigor, mas porque cada um alcança
uma classe que o outro não alcança por definição de escopo.

**Por unidade**, disparado pela entrega (`anchors deliver`). Garante a parcial, e chega
cedo: o defeito aparece quando corrigir ainda é barato, antes de as peças seguintes se
apoiarem nele.

**Por conjunto**, disparado pela entrega do PLANO. Garante a costura.

A medição que separa os dois: num E2E real, de 6 achados de um review adversarial,
**apenas 2 cabiam no escopo de uma unidade**. Os outros 4 atravessavam — e em todos eles
*cada peça estava correta sozinha*:

| achado | por que a peça isolada passa |
|---|---|
| o DAO trunca a consulta | truncar é válido; quem quebra é a regra que consome e escolhe "a versão de maior mês" |
| campo de propriedade não plantado | a spec do modelo manda; o DAO não planta; o handler não preenche — três peças, nenhuma errada sozinha |
| obrigação de exclusão sem prova | a tabela aparece no handler (estrutura verde); mutar o bloco não derruba teste nenhum |
| registro afirma o que não fez | a divergência é entre o declarado e o disco, não dentro de nenhuma peça |

Nenhum desses é alcançável por um revisor com escopo de um arquivo, e nenhum é
confrontável por gate — porque em todos a contradição mora ENTRE artefatos que,
separadamente, cumprem o que prometem.

Daí a regra: **"cada parte passou" não é prova do todo.** É a premissa que os 4 achados
desmentem, e é por isso que o segundo nível existe.

---

### O que a medição derruba: contradição entre regras

A **contradição entre duas regras da mesma spec** foi o defeito mais caro de uma rodada de
review adversarial — duas regras se anulavam, e o bug crítico morava exatamente ali. O
gate parece óbvio: leia as regras, confronte umas com as outras.

Ele não existe, e a razão é a pergunta que antecede todo gate: *isto é confrontável por
padrão de texto, ou exige entender o que a frase significa?* Contradição é semântica. Duas
regras que se anulam usam as mesmas palavras que duas regras que se complementam — o que
muda é o sentido. Medido no repositório real: 2.540 linhas de regra, 47 citando outra
regra. Nenhum padrão separa "B09 depende de B05" de "B09 contradiz B05".

O que sobrou desta tentativa foi melhor que o gate: ao escrever o detector, ele acusou 37
specs com regras citadas que não existiam. **15 eram falso positivo meu** (a regra estava
numa lista, não numa tabela). Das restantes, 8 eram um defeito real e inesperado — a
tabela tinha `\n` LITERAL no lugar da quebra de linha, então três regras ocupavam uma
linha física. O renderizador mostrava uma; qualquer gate que lesse linha a linha via uma;
as outras duas existiam apenas no arquivo. E as 25 finais não eram defeito nenhum: eram
tabelas de testID citando o cenário que provam.

A lição não é "não tente". É que **a tentativa paga mesmo quando o gate não nasce** — o
detector descartável achou 8 defeitos reais que nenhum gate teria achado, porque para
escrevê-lo foi preciso medir o repositório de um jeito que ninguém tinha medido.

### Onde o gate determinístico ACABA

Nem toda classe de defeito vira gate, e insistir produz coisa pior que a ausência: um
gate que erra na maioria é desligado no primeiro dia, e leva junto a credibilidade dos
que funcionam.

A **coerência interna da unidade** — "o código faz algo que nenhuma regra da spec
governa" — parece o próximo gate óbvio e não é. Medido ao tentar escrevê-lo:

- **32%** dos campos declarados no código aparecem literalmente na spec. Não é
  descuido: a spec descreve o campo em PROSA ("o dono é um lançamento OU uma
  recorrência, há um discriminador de tipo") porque a régua manda ser **agnóstica ao
  framework**. Um gate que exigisse o nome literal geraria 138 falsos positivos e
  brigaria com a própria doutrina.
- A tentativa seguinte — confrontar as citações de regra que o CÓDIGO faz — encontrou
  **zero** citações da própria spec em código de produção: 4429 das 4451 estão em
  testes, onde o `feature-test-match` já as confronta.

O que sobrou dessa classe é o que **já tem dono**: recurso prometido e não provisionado
é `obligation-honored`; símbolo prometido e não usado é `dependency-honored`. O resto —
campo sem leitor, regra que contradiz outra regra — é **julgamento**, e o lugar dele é
o review, não um regex.

A regra prática: antes de escrever um gate, **meça o falso-positivo contra o repositório
real**. Se ele acusa a maioria, o defeito não é do projeto — é do gate.

---

### O terceiro estado: `pendente` — o que o framework não pode exigir

Informativo e bloqueante são as duas pontas da escala de **política**: o gate *sabe* a
resposta, e o projeto decide se ela barra o merge. Existe uma terceira situação, de
natureza diferente: o gate **não tem como saber**, porque o insumo depende de uma
capacidade que o projeto pode legitimamente não ter.

O caso concreto: o Anchors pode exigir que a spec decida, que a trinca exista, que a
dependência prometida seja honrada — tudo isso ele lê do repositório. Não pode exigir
que o projeto rode uma ferramenta de **mutação**, emita **JUnit XML** ou gere **lcov**.
Isso depende de stack, de linguagem, de tempo de CI. Exigir seria o framework decidindo
pelo projeto — o oposto de *"a complexidade do projeto não sobe para o framework"*
(`CONCEPT.md` §5.2).

Mas **não exigir não é ficar calado**. Um gate que fica `pendente` em silêncio é pior
que ausente: ele ocupa a linha do relatório com um veredito que parece benigno, e o
pipeline segue verde certificando trabalho que ninguém mediu. É o mesmo mecanismo do
buraco de cobertura (§5.1), com o disfarce mais convincente.

A resolução é separar **veredito** de **visibilidade**:

| | o gate faz | o relatório faz |
|---|---|---|
| **bloqueante** | reprova e barra | conta como falha |
| **informativo** | reprova e registra | conta como issue |
| **pendente** | não julga | **nomeia o que falta E o que se perde sem isso** |

O que torna o `pendente` honesto é o *detalhe*. Não basta "sem sinal ingerido": o gate
diz qual comando produz o insumo, e **qual risco concreto fica descoberto** enquanto o
insumo não existir. Um `pendente` sem essa frase é um buraco com nome bonito.

E há uma assimetria a respeitar no relatório agregado: **"nada abaixo do limiar" e "nada
medido" são conclusões opostas** e não podem imprimir o mesmo número. `0 arquivos abaixo
de 70%` lido num projeto que nunca ingeriu cobertura é uma afirmação falsa produzida por
um relatório tecnicamente correto. Quando o denominador é zero, o relatório diz *"não
medido"*, nunca *"zero problemas"* — é a regra de ouro (**ausência de prova não é prova
de ausência**) aplicada à própria saída da ferramenta.

Onde o alerta aparece, e por quê em três lugares:

- no **`check`**, como indeterminado com motivo — quem está mexendo no arquivo vê;
- no **`doctor`/`status`**, como achado sistêmico — quem olha o projeto inteiro vê que
  uma família de gates está inerte;
- no **`coverage`**, no lugar onde o número seria lido — porque é ali que a ausência de
  medição vira conclusão errada.

A alavanca fica com o projeto: assim que a ferramenta existir e o sinal for ingerido, o
mesmo gate passa a julgar de verdade, e a maturação normal (§7) decide quando ele barra.

---

## 8. Onde a qualidade é imposta

O pilar distingue **onde** cada gate roda, porque isso decide o que é barato e o
que é definitivo:

- **Localmente**, antes de registrar trabalho (ex.: um hook de commit): gates
  baratos e incrementais sobre o que mudou — feedback rápido, escopo pequeno.
- **Na promoção** (ex.: antes de integrar à linha principal): o pipeline completo
  dos gates bloqueantes — a fronteira real, "nada sobe se não passar".

O POC prova a política "nada sobe sem passar nos gates" com um hook local barato
(arquitetura sobre arquivos alterados) e um pipeline de integração completo que
trava a promoção. A distinção importa: o gate local dá agilidade; o gate de
promoção dá a garantia. O Anchors define ambos como parte do pilar, com
granularidades diferentes.

---

## 8.0 O sinal chega pelo `anchors test`, não à mão

O gate de cobertura reporta o que o **mapa** sabe. Se o sinal no mapa for de uma
execução anterior, o gate fica verde sobre o que não mediu — e não há como
distinguir isso de "passou".

```yaml
tests:
  - layer: unit
    run: pnpm test
    run_changed: pnpm test -- --findRelatedTests {{files}}
    junit: .test-results/junit.xml
    lcov: .test-results/coverage/lcov.info
```

Declarado, `anchors test` **roda a suíte e ingere numa operação só**. É a operação
única que garante a correspondência: o sinal é daquela execução, e de nenhuma outra.

`anchors ingest` chamado direto quebra essa garantia sem avisar — dá para rodar a
suíte, editar o código, ingerir o relatório velho, e o mapa passa a afirmar uma
cobertura que já não vale. Aconteceu: o `line-coverage` aprovava um nó quando havia
dois cobertos, e o gate estava certo — o sinal é que estava velho.

Por isso o `ingest` manual **avisa**, e pode ser **recusado**:

```yaml
workflow:
  manual_ingest_blocks: true   # padrão: false
```

O padrão avisa porque há usos legítimos — um CI que rodou a suíte noutro job, uma
ferramenta que o `tests:` não cobre. Barrar todos por causa do caso comum tiraria a
saída de quem tem razão. E sem `tests:` declarado nada é exigido: o `ingest` é a
única forma de o sinal chegar ao mapa, e exigir o que não existe seria um beco sem
saída.

---

## 8.1 A mensagem de commit também é confrontada

O `commit-msg` é o segundo ponto de imposição local, e ele existe por uma razão que
não é estética: **o changelog nasce dos commits**.

Um assunto fora do formato não some do histórico — some do **changelog**. E isso só
se descobre quando alguém gera a primeira versão e o que faltou já está a centenas
de commits de distância. Commit já feito não se conserta.

### O caso que motivou

O squash de um PR entrou como:

```
[MTUAO] Plano 0017 — mutação, revisando o plano 0001
```

O GitHub usa o **título do PR** como mensagem do squash, e o título estava no
formato do card. Esse commit — o que introduzia o plano inteiro — não apareceria no
changelog.

Isso mostra por que a régua fica no hook e não num prompt interativo (`commitizen`
e afins): a mensagem do squash **não passa por terminal nenhum**. Um prompt ajuda a
escrever; só o `commit-msg` garante, porque é por onde toda mensagem passa.

### A régua

Conventional Commits, e não uma convenção nossa — é o que as ferramentas de
changelog já sabem ler. Inventar formato daria trabalho duas vezes.

```
tipo(escopo)!: o que mudou
```

| conferência | o que protege |
|---|---|
| tipo da lista fechada | tipo livre vira sinônimo (`bugfix`, `hotfix`) e desfaz o agrupamento |
| tipo em minúsculas | `Feat` e `feat` viram grupos **separados** na mesma versão |
| escopo não vazio | `feat(): x` parece que alguém ia dizer algo e parou |
| assunto ≤ 100 | acima disso é cortado na lista de commits e no changelog |
| sem ponto final | o changelog emenda o assunto a marcadores, e o ponto sobra |

Cada defeito tem **laudo próprio**. Um "formato inválido" genérico obrigaria quem
foi barrado a adivinhar qual regra quebrou — e adivinhar três vezes é o que faz
alguém desligar o hook.

O que o **git** gera passa (`Merge`, `Revert`, `fixup!`, `squash!`): ninguém
escreveu aquilo, e barrá-lo quebraria operações normais em vez de melhorar o
histórico. Mas só o formato exato escapa — `Mergeando o trabalho` é humano e não
passa, senão o prefixo vira brecha.

### Relação com o commitlint

O [commitlint](https://commitlint.js.org) é a ferramenta madura desta régua, e a
implementação embutida foi confrontada com ela caso a caso: **9 de 10 assuntos
recebem o mesmo veredito**.

O único divergente é `feat(): x` — o commitlint aceita escopo vazio, o Anchors
barra. A **direção** dessa divergência é deliberada e está fixada por teste: ser
mais estrito é seguro, porque o que passa aqui passa lá. Um projeto que troque a
régua embutida pelo commitlint não descobre um histórico que a ferramenta nova
reprova; o contrário produziria exatamente isso.

Uma regra do commitlint fica de fora: `subject-case`, que barra assunto começando
com maiúscula. Ela reprova `feat: SBOM sai da pasta ignorada` — não distingue sigla
legítima de frase capitalizada. Num projeto que fala de SBOM, CI e PR isso é atrito
sem ganho, e barrar sem razão empurra para o mesmo lugar de sempre: desligar o hook.

**Onde há Node, use o commitlint** — mais regras, mais configurável, e é o que o
ecossistema de changelog espera encontrar. A régua embutida existe para o projeto
que não tem: exigir um runtime inteiro para validar uma linha de texto contradiria
a proposta de funcionar em qualquer stack.

---

## 9. Resumo do pilar

- **Gate = âncora de qualidade.** Guia um limiar medido, confronta executando a
  medição, emite issue `violation` se abaixo. Reusa âncora + grafo + issue do
  CONCEPT.
- **Qualidade = conjunto.** "Gate de qualidade" é um guarda-chuva de muitos gates,
  cada um medindo uma coisa; a qualidade emerge do conjunto, não de um gate.
- **Dois tipos de medidor.** Determinístico (roda comando, lê exit code: lint,
  spellcheck, validators, testes) e **julgamento por IA** (um agente mede o que
  script não computa: coerência, legibilidade, spec completa — o gate de review por
  IA). Mesma mecânica, medidor diferente.
- **Dupla saída: issue e bloqueio.** Todo gate que reprova gera as duas — mas a
  **issue é inegociável** e o **bloqueio é negociável** (`--no-block`). Pode-se optar
  por não bloquear, nunca por não registrar.
- **Ponte feature → teste.** A feature é especificação executável; vira teste no
  runner nativo; o teste é gate. A relação é **N:M mediada por código de cenário**,
  a chave verificável que amarra spec ↔ feature ↔ teste.
- **Ligação declarativa.** O framework oferece o mecanismo de ligar medidores
  externos (runners, linters, validadores, agentes) como gates; o projeto declara
  comando, escopo, leitura de resultado, limiar e força. Substitui a lista imperativa
  do POC.
- **Fechar todas as pontas.** A falha própria é o **buraco de cobertura** — algo que
  nenhum gate cobre (a ponta aberta, silenciosa), distinto de "um gate reprovou". Um
  meta-gate confronta a completude; onde o script não alcança, o review por IA fecha.
- **Pipeline incremental.** O conjunto de gates roda sobre o caminho de impacto
  (Propagação); falhas viram issues materiais. Barato o bastante para ser lei.
- **Agregação por perfil.** Qualidade é o perfil de vereditos por categoria de gate,
  não um score. Promoção = todos os bloqueantes passaram.
- **Maturação `informativo → bloqueante`.** Todo gate tem estado de maturação
  (política permanente); a promoção é decisão versionada. `--no-block` é o override
  *pontual* distinto — deixa passar uma vez, sem mexer na política nem apagar a issue.

O que o pilar entrega: qualidade que não depende da memória nem do cuidado de uma
sessão. Ela está ancorada em gates versionados, medida a cada mudança de forma
incremental, e o nível de maturidade do projeto é legível no estado desses gates.
O "bem feito" para de ser sensação e vira propriedade verificável.

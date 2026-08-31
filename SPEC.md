# Anchors — Pilar de Spec

> Este documento define o **pilar de Spec** do Anchors. Ele pressupõe o mecanismo
> de [`CONCEPT.md`](./CONCEPT.md) e se conecta a todos os outros pilares — porque
> a spec é onde eles se amarram. A Spec é a **âncora-base**: a origem da verdade e
> o ponto de sustentação do qual o resto pende.
>
> É teoria de framework, já validada por uma ferramenta que implementa
> (é uma ferramenta *spec-first*) e por uma prova de conceito real. (Não
> confundir este documento — o pilar conceitual — com o `SPEC.md` que uma
> ferramenta concreta pode ter na própria raiz, que seria a especificação
> daquela ferramenta.)

---

## 1. O safepoint do qual tudo pende

Numa escalada, algumas âncoras são passagem e outras são **safepoint** — o ponto
fixo, bem cravado, do qual a corda inteira depende. Se ele falha, a queda é longa.
A **Spec** é o safepoint do Anchors. É a âncora que **segura a corda** com mais
força, e é a ela que todo o resto se prende.

Enquanto o **Planejamento** (`PLANNING.md`) é a origem do *movimento* — de onde vem
a alteração —, a **Spec** é a origem da *verdade* — o critério do que a coisa deve
ser. O plano é o vetor; a spec é o alvo. O plano diz *o que fazer e em que ordem*;
a spec diz *o que a coisa deve ser*. Juntos são as duas origens: o plano inicia a
onda, a spec dá a ela um destino verdadeiro.

---

## 2. A disciplina spec-first é o pilar

O pilar não é o *arquivo* de spec. É a **disciplina spec-first**: a lei de que
**toda coisa nasce de uma spec**, de que a spec é a fonte de verdade, e de que nada
existe no sistema sem uma spec que o governe e o amarre aos outros pilares.

Essa distinção importa. Um arquivo `.spec.md` solto é só um formato. O que *opera*
— o que faz a Spec ser pilar — é a disciplina:

- a spec vem **antes** do código (o código realiza a spec, não o contrário);
- a spec é a **fonte de verdade** (quando há divergência, é ela que arbitra o que é
  certo, até que se decida atualizá-la);
- **nada entra no sistema sem passar por uma spec** que o ligue à rede.

Essa lei é o que dá aos outros pilares um lugar onde se ancorar. É por isso que a
Spec é o pilar mais **central**, ainda que não seja o mais dinâmico: ela é o nó de
maior grau no grafo, o ponto de convergência.

---

## 3. A spec amarra os pilares

O que você disse ao propor este pilar: *as specs amarram a propagação com a
rastreabilidade e servem como âncora para os gates de qualidade.* A spec é onde os
pilares se **encontram**.

- **Amarra a Rastreabilidade** (`TRACEABILITY.md`): é na spec que os requisitos
  nascem com sua **identidade estável** (o código de cenário/estado). Essa
  identidade é o que a rastreabilidade propaga por toda a cadeia (spec → feature →
  teste → código). A spec é a nascente da identidade contínua.
- **Amarra a Propagação** (`PROPAGATION.md`): a spec é o nó por onde a onda passa
  em ambas as direções. Uma mudança no código torna a spec stale; uma mudança na
  spec propaga para tudo que dela deriva. Ela é o pivô da propagação — sobe e desce
  a partir dela.
- **É a âncora dos gates de Qualidade** (`QUALITY.md`): os gates confrontam o
  trabalho *contra a spec*. É a spec que declara o comportamento que os testes
  devem cobrir e que o código deve realizar. Sem a spec como régua, os gates medem
  no vazio.

Por isso a spec é o **ponto de amarração**: tire-a, e a rastreabilidade não tem
nascente, a propagação não tem pivô, e a qualidade não tem régua. Os outros pilares
não deixam de existir — deixam de ter *onde se prender*.

---

## 4. O que uma spec é (e o que não é)

A spec captura **o quê**, não **o como** — comportamento, contratos e invariantes,
nunca detalhes de implementação. Isto é o anti-drift aplicado à spec (`CONCEPT.md`
§2): detalhes de implementação driftam no primeiro refactor; contratos sobrevivem.
Uma spec que descreve implementação vira uma segunda fonte de verdade que compete
com o código e apodrece.

Consequências da regra (destiladas do POC):

- **Descreve comportamento, não código.** Sem blocos de código na spec — eles
  criariam um alvo que drifta. Prosa, tabelas, contratos, invariantes.
- **Cada requisito ganha identidade na spec.** Estados, regras e requisitos
  recebem seu código estável aqui — é a nascente da rastreabilidade.
- **Preserva o estável ao evoluir.** Ao atualizar, preserva identificadores,
  histórico e as partes ainda corretas — não reescreve por atacado.
- **Sem placeholders.** Uma seção ou está preenchida ou é removida. Um safepoint
  meio-cravado é pior que safepoint nenhum: dá confiança falsa.

---

## 5. Três níveis: guide, template, spec

A Spec não é um artefato solto — é uma **cadeia de âncoras regendo âncoras**, com
três níveis:

```
GUIDE de spec   → as regras universais: como toda spec se comporta            [1 por projeto]
  └─ governs
TEMPLATE        → o esqueleto de UM tipo de spec (por camada/tipo de arquivo)  [1 por tipo]
  └─ governs
SPEC            → a instância concreta que descreve um target                  [N no projeto]
```

O **guide** rege os templates; os templates regem as specs — o mesmo `governs` do
grafo (`CONCEPT.md` §3), aplicado dentro do próprio pilar. Isso liga a Spec à
**Estrutura de Projeto** (`STRUCTURE.md`): os *tipos de template* correspondem às
*camadas que a Estrutura declara*. Tipicamente uma camada corresponde a um tipo de
spec e um template — mas a relação nem sempre é 1:1 (uma camada pode render dois
tipos, um tipo pode cobrir dois arquivos; ver [`SPEC_TYPES.md`](./SPEC_TYPES.md)). A
Estrutura diz "existe a camada X"; a Spec provê o(s) template(s) para X.

A prova de conceito prova os três níveis: tem um guide de spec e dois templates
(tela e componente); cada `.spec.md` de um target segue o template do seu tipo.

### O template = esqueleto comum + bloco de variação

Levantar os padrões de spec por linguagem/camada (frontend, backend, infra, CLI,
lib) revela que os templates **não são estruturas independentes** — são **um
esqueleto comum + um bloco de variação**:

**Esqueleto comum** — toda spec, de qualquer tipo, carrega:

| bloco comum | papel |
|---|---|
| **Identidade** | código estável + arquivo + tipo/camada (a nascente da rastreabilidade) |
| **Propósito / responsabilidade** | o que faz e o que **não** faz (o limite da camada) |
| **Contrato de entrada** | o que recebe (params / props / request / evento / args / inputs) |
| **Contrato de saída / efeitos** | o que retorna ou muda no mundo |
| **Regras / comportamento** | pré/pós-condições, invariantes |
| **Erros / falhas** | como e por que falha (e quem sinaliza) |
| **Dependências** | de que camadas depende — e quais é **proibido** chamar (as fronteiras vêm da Estrutura, `STRUCTURE.md`; a spec apenas as instancia/estreita para este target) |
| **Rastreabilidade** | os códigos que descem para feature/teste |

**Bloco de variação** — cada tipo acrescenta *um* bloco, ditado pela sua "fonte da
variação". O que muda não é a estrutura; é o **vocabulário que preenche
entrada/saída/contrato**:

| tipo de spec | fonte da variação → bloco que distingue |
|---|---|
| tela | estado + navegação |
| componente | props + variantes + eventos |
| hook | efeitos (queries/mutations) |
| repository | queries / índices / acesso a dado |
| handler / interface | rota (ou evento) + status |
| usecase | regras de negócio |
| entity | invariantes + máquina de estados |
| CLI command | args / flags / exit codes |
| biblioteca | API pública exportada |
| infra / devops | recursos / secrets / permissões |

Por isso o guide define **um template base** e os tipos são **overlays finos** sobre
ele — não N templates do zero. A lista concreta de tipos e o conteúdo de cada bloco
de variação vive num catálogo à parte ([`SPEC_TYPES.md`](./SPEC_TYPES.md)), porque
cresce por projeto/linguagem e incharia o pilar.

---

## 6. Dois regimes: spec comportamental e spec declarativa

Nem toda spec descreve *comportamento*. Há dois **regimes**, ortogonais ao tipo:

| | **comportamental** | **declarativa** |
|---|---|---|
| descreve | ação: "dado input X, acontece Y" | estado desejado: "estes recursos existirão, assim conectados" |
| exemplos | tela, hook, usecase, CLI command, lib | Dockerfile, manifesto K8s, módulo Terraform, schema de dados |
| tipo de teste natural | comportamental (cenários dado/quando/então) | **de conformidade** (o recurso declarado existe e está conforme?) |
| modo de falha | runtime, caso-limite, exit code | reconciliação, config inválida, drift |

O regime não bifurca o caminho da propagação — a cadeia é a mesma (spec → feature →
teste → código, `TRACEABILITY.md`). O que ele muda é o **tipo de teste** que faz
sentido na ponta. Para isso, o Anchors **expande a noção de teste**: teste não é só
o cenário Gherkin comportamental — inclui o **teste de conformidade** (validação
estrutural: o recurso declarado bate com a forma especificada?). Uma spec
declarativa é testável — por conformidade.

### Testável por default; `@no-test` como opt-out

**Toda spec é testável por default**, inclusive a declarativa. Mas testar
conformidade de infra nem sempre é desejado — depende do usuário. Então o
mecanismo é opt-out, via **tag**, do mesmo jeito que se tagam nível e prioridade. É
a instância, na Spec, do **opt-out honesto** (`CONCEPT.md` §5.1): explícito,
registrado e com justificativa datada (a tag carrega o *porquê* não se testa).

- uma spec (ou um requisito dentro dela) marcada **`@no-test`** é dispensada da
  exigência de teste — o gate de rastreabilidade "requisito realizado"
  (`TRACEABILITY.md`) **não** a trata como órfã;
- sem a tag, a spec é testável e cobra sua encarnação de teste (comportamental ou
  de conformidade, conforme o regime).

Assim a cadeia da Rastreabilidade permanece **universal** — não há ramo especial
para declarativas; o que existe é um teste de tipo diferente e um opt-out explícito.

### `@TBD` — o que ainda não foi desenvolvido

`@no-test` afirma **"esta unidade não precisa de teste"**. É decisão permanente, e
fica escrita para quem ler a spec depois.

Mas a spec **nasce antes do código** — é o fluxo normal, já que ela é a âncora. E
aí existe uma segunda pergunta, que o opt-out não responde: *"e enquanto a peça não
existe?"*

```markdown
> **@TBD: code,feature,test** — as peças que realizam esta spec são a fase em andamento.
```

`@TBD` é *to be developed*: a peça **está decidida e ainda não foi escrita**. Não é
"a decidir" — o que ela vai fazer já está na spec, que é justamente o que a spec
**é**. Falta o arquivo.

A diferença entre os dois é o **tempo**, e é o que torna o `@TBD` honesto:

| marca | afirma | vence? |
|---|---|---|
| `@no-test:` | esta unidade não precisa de teste | **não** — é permanente |
| `@TBD: test` | falta escrever | **sim** — no instante em que a peça aparece no mapa |

O `@TBD` **vence sozinho**: ninguém precisa lembrar de removê-lo. No momento em que
a peça nasce, o gate volta a confrontá-la — e se a marca continuar lá, ela não
protege mais nada, porque a peça existe.

**O alvo é obrigatório.** `@TBD` sem dizer o quê viraria um interruptor geral do
gate, e é exatamente isso que ele não pode ser. Declarar `code` dispensa o código e
**continua cobrando** feature e teste.

Sem esta distinção, quem escreve uma spec nova tem duas saídas e ambas são ruins:
barrar o commit de todo trabalho em andamento, ou declarar `@no-test` mentindo — e
aí a cobrança some **para sempre**, justamente na unidade que mais vai precisar
dela.

O eixo é ortogonal ao tipo: um mesmo artefato pode ter partes dos dois regimes (um
pipeline de CI é declarativo na definição de stages, comportamental nos gates). A
spec declara seu regime — e o regime determina *qual tipo de teste* o requisito
exige; a Propagação segue as arestas normalmente, sem roteamento paralelo.

### O regime é por requisito, não por spec

O regime não é um rótulo da spec inteira — é uma propriedade **de cada requisito**,
declarada como **tag**, do mesmo jeito que nível e prioridade (`TRACEABILITY.md`). O
requisito comportamental leva sua tag de regime comportamental; o declarativo, a sua.
Assim uma spec **mista** não é um terceiro regime — é uma spec cujos *requisitos*
têm regimes diferentes.

O exemplo canônico é a **entity**: seus *campos* são declarativos (a forma existe e
está conforme?) e suas *invariantes / máquina de estados* são comportamentais (dado
X, a transição acontece?). A spec da entity não escolhe um regime — cada requisito
seu carrega o próprio, e cada um puxa o tipo de teste que lhe cabe (conformidade para
os campos, comportamental para as invariantes). Reusa o mecanismo de tag que já
existe; não inventa nada.

---

## 7. A falha característica: o que existe sem spec que o governe

A falha própria do pilar é a **coisa sem spec** — código, comportamento ou artefato
que existe no sistema sem uma spec que o governe. É escalar sem cravar o safepoint:
a coisa está lá, funciona talvez, mas **nada a sustenta** — não há fonte de verdade
contra a qual medi-la, nem identidade que a ligue à rede, nem régua que a
confronte.

Essa falha é distinta:

- não é órfão de **Rastreabilidade** (uma peça pode ter arestas e ainda não ter
  spec que a governe — está ligada, mas a nada que a *manda*);
- não é baixa **Qualidade** (a coisa pode passar em todo gate genérico e ainda não
  ter uma spec própria que diga o que *ela* deveria ser);
- não é falta de **Planejamento** (pode ter sido planejada e ainda pular a etapa da
  spec).

É **ausência de governo**. Uma coisa sem spec é uma coisa que o projeto não
controla — ela evolui sem critério, e quando diverge, não há contra o quê julgar.
O pilar de Spec é o que garante que **tudo que existe é governado por algo**.

---

## 8. Os gates da Spec

A disciplina spec-first se impõe por gates (âncoras de confronto, `QUALITY.md` §2)
que medem *governo e completude*, não outra dimensão:

| gate | pergunta de confronto | falha = |
|---|---|---|
| **spec presente** | toda coisa que exige governo tem sua spec? | coisa sem governo |
| **spec completa** | a spec tem as seções obrigatórias, sem placeholder? | safepoint meio-cravado |
| **spec-first honrada** | a spec descreve o quê (não o como), sem código? | segunda fonte de verdade |
| **cobertura declarada** | tudo que a spec declara (estados, regras) virou identidade rastreável? | requisito que a spec promete mas não ancora |

A **trinca completa** (spec + código + prova) é cobrada pelo mesmo modelo, com dois
opt-outs que dizem coisas diferentes: `@no-test`/`@no-feature` para a peça que **não
existirá**, e `@TBD` para a que **ainda não foi escrita** (§6). O segundo vence
sozinho quando a peça nasce; o primeiro é permanente.

Esses gates seguem o modelo comum: rodam sobre o caminho de impacto da Propagação
(`PROPAGATION.md` §3), emitem issues materiais quando falham, e maturam de
`informativo` a `bloqueante`
(`QUALITY.md` §7) — um projeto pode começar com muitas coisas sem spec (gate
report-only, "temos N lacunas") e promover o gate a bloqueante quando o governo
fecha.

---

## 9. Relação com os outros pilares

A Spec é o centro de gravidade:

- **Vive dentro da Estrutura de Projeto** (`STRUCTURE.md`): a Estrutura declara *que
  existe* a camada; a spec (inclusive a de arquitetura) é o conteúdo *dentro* dela.
  O gabarito precede o conteúdo.
- **Depende do Planejamento** (`PLANNING.md`): o plano decide *quais* specs de
  partida nascem; a spec preenche *o quê* cada uma é (a ordem das camadas vem da
  Estrutura, não do plano). Planejamento é o vetor, Spec é o alvo.
- **Nascente da Rastreabilidade**: a identidade estável dos requisitos nasce aqui.
- **Pivô da Propagação**: a onda sobe e desce a partir da spec.
- **Régua da Qualidade**: os gates confrontam o trabalho contra a spec.
- **Governa a Documentação**: a doc de consumo (`CONCEPT.md` §2) é derivada, entre
  outras fontes, da spec — a spec dá o conteúdo verdadeiro que a doc apresenta para
  fora.

Num roteiro de amadurecimento, a Spec vem logo após o Planejamento na rota (depois
da Estrutura e do Planejamento): é o safepoint que se crava antes de confiar o peso
à corda.

---

## 10. Resumo do pilar

- **Spec = a âncora-base, o safepoint.** A âncora que segura a corda com mais força;
  a ela todo o resto se prende. Origem da *verdade* (o alvo), como o Planejamento é
  a origem do *movimento* (o vetor).
- **O pilar é a disciplina spec-first**, não o arquivo: toda coisa nasce de uma
  spec, a spec é fonte de verdade, nada existe sem uma spec que o governe.
- **Amarra os pilares.** Nascente da Rastreabilidade (identidade), pivô da
  Propagação (a onda sobe e desce por ela), régua da Qualidade (os gates confrontam
  contra ela). É o ponto de convergência.
- **Captura o quê, não o como.** Comportamento e contratos; sem código; sem
  placeholder; preserva o estável.
- **Três níveis: guide → template → spec.** O guide rege os templates; os templates
  regem as specs. Os tipos de template correspondem às camadas que a Estrutura
  declara. O template é **esqueleto comum + um bloco de variação** por tipo — não N
  estruturas do zero. (Catálogo de tipos em [`SPEC_TYPES.md`](./SPEC_TYPES.md).)
- **Dois regimes: comportamental e declarativo.** O regime é **por requisito** (tag,
  como nível/prioridade), não da spec inteira — uma spec mista (ex.: entity) tem
  requisitos de regimes diferentes. O regime muda como propaga: comportamental →
  features/testes; declarativa → conformidade de forma. Eixo ortogonal ao tipo.
- **Falha característica = coisa sem spec que a governe.** Ausência de governo — a
  coisa existe mas nada a sustenta nem a julga.
- **Gates de governo e completude.** Spec presente, completa, spec-first honrada,
  cobertura declarada. Mesmo modelo: caminho de impacto, issues, maturação.

O que o pilar entrega: um projeto onde **tudo que existe é governado por uma fonte
de verdade**, e onde essa fonte é o ponto em que rastreabilidade, propagação e
qualidade se encontram. É o safepoint que faz a escalada inteira ser segura — o nó
que, bem cravado, sustenta todo o peso do resto.

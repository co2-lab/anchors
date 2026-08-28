---
title: Conceito
---

> **Anchors** é um framework de continuidade para desenvolvimento assistido por IA.
> Este documento define o **conceito**: o que é uma âncora, como as âncoras se
> relacionam, como o sistema sabe o que está em dia, e como a dessincronia é
> tratada. É teoria pura — independente de qualquer ferramenta.
>
> Existe uma ferramenta que implementa o Anchors, e uma prova de conceito
> real onde o conceito é validado na prática. Ambas aparecem aqui apenas
> como instâncias — o Anchors não depende de nenhuma delas.

---

## 1. O problema

Todo trabalho assistido por IA sofre de **amnésia estrutural**. Uma sessão de IA
é sem estado entre invocações: a sessão acaba, o contexto some, e a próxima
sessão — outro agente, outro dia, outro humano — recomeça sem as regras, sem
saber o que já foi decidido e por quê, sem saber onde o trabalho parou. O projeto
cresce em código, mas a **direção** não persiste.

Ferramentas de agente já sabem gerenciar a amnésia *dentro* de uma sessão
(compressão de histórico, sumarização, reset intencional de contexto). O que
falta é a continuidade *entre* sessões e *ao longo da vida* do projeto — manter,
com rigor, os mesmos patterns, a mesma estrutura e a mesma direção que as sessões
anteriores estabeleceram.

O que persiste entre sessões não é memória. É **artefato ancorado no
repositório**. O Anchors é sobre projetar esses artefatos: quais são, que papel
têm, como não apodrecem, e como uma sessão futura é obrigada a respeitá-los.

---

## 1.1 Maturidade: o conceito guarda-chuva

Entregar não é só construir a coisa certa — é entregar com um nível de qualidade
que a torne efetiva. Um projeto pode *funcionar* e ainda assim ser inefetivo:
ficar indo e voltando em falhas, sem nunca atingir um patamar confiável. A
diferença entre um projeto que funciona e um projeto que está **pronto** é a
**maturidade**.

No Anchors, maturidade é o conceito no topo. Um projeto é maduro quando tem
**todos os seus pilares implementados e vigorosos**. Um projeto imaturo é aquele
a que faltam pilares, ou cujos pilares são frouxos. Maturidade não é um número
medido num artefato isolado — é a **presença e o vigor dos pilares** no projeto
como um todo.

Isso não é só descritivo: a maturidade é **medida** pelo **validador de saúde do
ecossistema** (`QUALITY.md` §5.2) — um meta-gate de visão global (materializado num
comando de CLI, tipo `anchors doctor`) que varre o estado dos pilares e as pontas
sistêmicas, e responde "o framework aplicado a este projeto está íntegro e maduro?".
É o que transforma "presença e vigor dos pilares" de uma noção num veredito.

Os pilares são os mecanismos estruturais que o Anchors define. Este documento
especifica o mecanismo comum a todos eles — a **âncora** (§2), o **grafo** (§3),
a **sincronia incremental** (§4), o **rastreio de dessincronia** (§5) e a
**separação vivo/histórico** (§6). Cada pilar temático é uma aplicação
especializada desse mecanismo.

Pilares nomeados até agora, na ordem da rota (da origem ao acabamento):

- **Estrutura de Projeto** — [`STRUCTURE.md`](/docs/estrutura/). A planta da casa:
  define quais camadas existem, sua ordem/dependência e onde cada âncora mora. É o
  gabarito sobre o qual os outros pilares operam — o meta-nível que declara as
  camadas antes de qualquer spec preenchê-las.
- **Planejamento** — [`PLANNING.md`](/docs/planejamento/). A origem do *movimento*: semeia
  as specs de partida (nunca código), é o input do fluxo e carrega o norte entre
  sessões (para onde vamos, em que ordem, onde paramos). Sem ele, o projeto reage
  mas não avança com direção.
- **Spec** — [`SPEC.md`](/docs/spec/). A origem da *verdade*: a âncora-base, o
  safepoint do qual tudo pende. A disciplina spec-first amarra os outros pilares —
  é nascente da Rastreabilidade, pivô da Propagação e régua da Qualidade.
- **Rastreabilidade** — [`TRACEABILITY.md`](/docs/rastreabilidade/). A cola, em duas
  metades: dá a cada requisito uma identidade contínua através de suas formas (spec →
  feature → teste → código) e mantém o mapa de dependências entre os arquivos,
  garantindo que nenhuma peça vira uma ilha. É o solo em que os outros pilares
  fincam raiz.
- **Propagação** — [`PROPAGATION.md`](/docs/propagacao/). O motor: faz uma alteração
  num ponto percorrer o organismo pela Rastreabilidade, marcando o que ficou stale,
  até tudo voltar a ser coerente. É a propagação das alterações que faz o
  desenvolvimento avançar.
- **Qualidade** — [`QUALITY.md`](/docs/qualidade/). Sem qualidade *medida*, o "bem
  feito" é uma sensação que não sobrevive entre sessões. Define gates que medem se
  o trabalho atingiu um limiar, e como esses gates compõem a maturidade.

Outros pilares serão nomeados à medida que o framework amadurece. Nada aqui está
cravado na pedra — os pilares são ideias encontrando seu lugar, refinadas
conforme o conceito evolui.

A relação entre os documentos: o CONCEPT define o que é uma âncora (na metáfora da
escalada) e como ela aponta, segura a corda e demarca; cada pilar aplica isso a um
tipo de âncora ou dinâmica essencial — e define como compõem a maturidade do
projeto. A **documentação** é âncora, mas não pilar: é de consumo, não estrutural
(§2).

---

## 2. A âncora

O nome do framework vem da **escalada**. Para subir uma parede você crava
âncoras — pontos fixos que fazem três coisas: dizem **para onde ir** (o próximo
ponto na superfície sem apoio), **seguram a corda para você não cair** (o
safepoint), e **demarcam por onde você passou** (o registro da rota). No Anchors é
igual.

> Uma **âncora** é qualquer documento que redigimos e que nos acompanha na
> escalada do projeto — apontando o próximo passo, segurando a corda para não
> cairmos, ou marcando por onde passamos.

Essa é a definição-raiz, e é **generosa de propósito**: todo documento que
geramos e que é usado pela aplicação ou pelo framework é uma âncora. Specs,
features, planos, guides, docs — todos são âncoras, porque todos nos acompanham na
subida. A pergunta não é "este documento confronta algo?"; é "este documento é um
ponto fixo do qual dependemos para subir com segurança?".

### As três funções de uma âncora

Toda âncora cumpre uma ou mais das três funções da escalada:

- **Aponta (diz para onde ir).** A âncora guia: a IA a lê e gera/altera o trabalho
  seguindo-a. Um plano aponta a rota; uma spec aponta o que construir.
- **Segura (a corda, o safepoint).** A âncora confronta: review, validação, lint e
  check rodam *contra* ela. Se o trabalho diverge, ela sustenta — impede a queda,
  força a correção. Esta é a função **bloqueante**.
- **Demarca (o rastro).** A âncora registra por onde passamos — a rota subida, as
  decisões, a cara da aplicação para quem chega depois. Esta função é
  **informativa**: demarca sem bloquear.

Uma mesma âncora pode cumprir várias funções, e é a **força** com que ela segura a
corda que decide seu peso no confronto — o que já aparece nos tipos de aresta do
grafo (§3): arestas `blocking` seguram a corda; arestas `references` só demarcam o
rastro. A metáfora da escalada é o modelo; o grafo é sua materialização.

### Guiar + confrontar: a âncora que segura a corda

A forma mais rigorosa de uma âncora — a que **segura a corda** — combina apontar e
segurar num ciclo fechado: ela é **executável como critério**, bússola na ida e
régua na volta.

```
        ÂNCORA (spec / feature / plano / guide)
         ▲                              │
         │ segura a corda               │ aponta
         │ (review, validate,           │ (a IA gera/altera
         │  lint, check)                │  seguindo a âncora)
         │                              ▼
        ALVO (código / outro artefato) ── diverge? ──┐
                                                      │
                            ┌─────────────────────────┴───────────┐
                      corrige o alvo                     atualiza a âncora
                    (âncora certa,                     (âncora superada,
                     feito errado)                      feito intencional)
```

O Anchors nunca resolve essa divergência sozinho. Ele a **detecta e a
apresenta**. A decisão entre "corrigir o alvo" e "atualizar a âncora" é sempre de
quem opera. O framework é detector e rastreador de dessincronia, não juiz.

### Âncoras estruturais e âncoras de consumo

As âncoras se distinguem pela **direção** em que operam:

- **Estruturais** — regem *para dentro*, sustentam o sistema. A spec rege o código,
  o plano rege a spec, a arquitetura rege as camadas. São as que têm coisas
  penduradas nelas. É onde vivem os pilares.
- **De consumo (terminais)** — são geradas *para fora*, a partir do sistema, para
  quem chega de fora. A **documentação** é o exemplar: ela dá uma **cara** à
  aplicação para quem vai conhecê-la sem exercitá-la. É uma **folha** do grafo —
  nada depende dela, ela depende de tudo; a onda de propagação sempre *termina*
  nela, nunca dispara nada rio acima. Sua única obrigação de confronto é o
  **frescor**: "está atualizada com o que descreve?". Se o que ela descreve mudou
  e ela não, ela mente para o observador externo.

A documentação é âncora de primeira classe — cidadã do grafo, propaga, fica stale,
gera issue quando desatualiza — mas **não é um pilar**: sua importância é grande
para o usuário final e pequena para a mecânica do framework, e pilar é uma medida
de mecânica. Ela é sustentada; não sustenta.

### Anti-drift é a lei

Uma âncora que segura a corda ou aponta o caminho precisa de um mecanismo que a
impeça de **mentir** sobre o que descreve. Sem isso, âncoras apodrecem: a spec
descreve um comportamento que o código não tem mais, o guide cita um arquivo que
foi renomeado, o doc promete uma API que sumiu. Uma âncora que mente é pior que
ausência de âncora — é um safepoint solto, que dá confiança falsa e te derruba
quando você mais confia nela.

O anti-drift não é uma boa prática opcional; é o que mantém uma âncora confiável —
o que a impede de virar uma nota morta. Ele se manifesta em três regras que valem
para toda âncora:

- **Fundamentar na realidade.** Todo caminho, símbolo ou identificador citado
  numa âncora deve existir no momento em que ela é escrita. Se não existe, a
  âncora está errada.
- **Descrever o QUÊ, não o COMO.** Âncoras capturam comportamento, contratos e
  invariantes — nunca detalhes de implementação. Detalhes de implementação
  driftam no primeiro refactor; contratos sobrevivem.
- **Preservar o que é estável.** Ao atualizar uma âncora, preserve
  identificadores estáveis, histórico e as partes ainda corretas. Reescrever por
  atacado destrói rastreabilidade e quebra referências cruzadas.

---

## 3. O grafo de âncoras

Âncoras não vivem isoladas. Uma spec de arquitetura rege muitas specs de tela;
uma tela é regida por sua spec, pela arquitetura e por um guide de camada. A
relação é **muitos-para-muitos e dirigida** — é um **grafo**, não uma lista de
pares.

O grafo é **material**: um artefato versionado no repositório, revisável em Pull
Request, que sobrevive sem nenhum banco de dados. Qualquer índice em memória ou
`.db` é apenas cache reconstruível a partir dele. O grafo é fonte de verdade;
o cache, conveniência.

### Nós

Todo arquivo que participa do grafo é um nó. Muitos nós são âncora *e* alvo ao
mesmo tempo (uma spec é alvo do guide que a rege *e* âncora do código que ela
rege).

| campo | significado |
|---|---|
| `id` | caminho do arquivo (`login.spec.md`) |
| `kind` | `spec` \| `feature` \| `test` \| `code` \| `doc` \| `guide` |
| `rev` | revisão atual do conteúdo do nó |
| `updated_at` | quando o conteúdo mudou pela última vez |

### Arestas

Cada relação "rege / depende de" é **uma aresta dirigida própria**. É aqui que o
muitos-para-muitos vive: um nó aparece em quantas arestas forem necessárias, como
origem (`from`) e como destino (`to`).

```yaml
# anchors.graph.yaml (material, versionado)
- from: architecture.spec.md
  to:   login.spec.md
  type: governs        # arquitetura rege a spec de tela
  origin: declared

- from: login.spec.md
  to:   login.tsx
  type: specifies      # a spec descreve o código
  origin: convention   # veio da co-location de nomes

- from: login.spec.md
  to:   login.feature
  type: covered-by     # a feature cobre a spec

- from: login.feature
  to:   login.test.tsx
  type: tested-by      # o teste exercita a feature

- from: docs/auth.md
  to:   login.tsx
  type: references     # o doc só aponta — não confronta
```

### Aresta tipada: o tipo carrega a semântica

Cada aresta tem um **tipo semântico**. O tipo carrega duas coisas de uma vez: a
**força** (se a divergência bloqueia) e a **pergunta de confronto** (o que o
verificador pergunta ao comparar alvo e âncora).

| tipo | pergunta de confronto | força |
|---|---|---|
| `governs` | "o filho respeita os limites do pai?" | bloqueante |
| `specifies` | "o código faz o que a spec descreve?" | bloqueante |
| `covered-by` | "a feature cobre todos os cenários da spec?" | bloqueante |
| `tested-by` | "o teste exercita a feature?" | bloqueante |
| `references` | (nenhuma — só rastreabilidade) | informativa |

A força é **derivada do tipo**, não um campo separado. Uma aresta bloqueante que
falha no confronto vira uma issue (§5); uma aresta informativa nunca gera issue —
no máximo um aviso.

O tipo diz **o que confrontar**, nunca **quem manda**. Não há hierarquia embutida:
dois `governs` que discordam do mesmo alvo não se resolvem por precedência — viram
uma issue de conflito (§5). O Anchors não arbitra qual âncora vence.

### Origem da aresta

Cada aresta registra **como entrou no grafo**, porque as três fontes convivem no
mesmo grafo material:

| origem | como a aresta surge |
|---|---|
| `convention` | co-location de nomes (`login.tsx` → `login.spec.md`) |
| `declared` | a âncora declara explicitamente o que rege |
| `inferred` | derivada de imports/símbolos por uma ferramenta |

Independente da origem, a aresta é materializada e versionada no grafo. Inferência
e convenção *propõem* arestas; o grafo material é onde elas passam a existir de
fato.

---

## 4. Sincronia

O grafo diz *quem* depende de quem. Falta saber *se* algo saiu de sincronia. Cada
nó carrega uma **revisão** (`rev`) do seu conteúdo, que avança quando o arquivo
muda; e cada **aresta** carrega um carimbo de "validado contra qual rev de cada
ponta". Uma aresta está **stale** quando alguma ponta avançou desde o último
confronto — a âncora mudou, ou o alvo mudou, ou a aresta nunca foi validada.

Detectar staleness sem revalidar o projeto inteiro, e fazer uma mudança percorrer
o grafo até tudo voltar a ser coerente, é a **dinâmica** do sistema — e ela é um
pilar próprio: a **Propagação**, especificada em [`PROPAGATION.md`](/docs/propagacao/).

Aqui basta reter a estrutura: **o carimbo de sincronia vive na aresta** (não no
nó), porque "estar em dia" é propriedade de uma *relação*, não de um arquivo — a
mesma spec pode estar sincronizada com o código que descreve e atrasada em relação
à arquitetura que a rege. Como a sincronia é calculada, como a **análise de
impacto** percorre o mapa de dependências e produz o **caminho de impacto**, e o
que significa uma mudança "terminar de propagar" (quiescência) são o conteúdo do
pilar de Propagação.

---

## 5. Issues — o registro da dessincronia

Quando um confronto de aresta bloqueante falha, ou quando o grafo revela uma
dessincronia que o Anchors não pode resolver sozinho, o resultado é uma
**issue**: o registro material de uma divergência que precisa de decisão humana.

O Anchors **não arbitra**. Toda divergência não-resolvível mecanicamente vira
issue e para. Quem opera decide.

### Três kinds

Uma issue nasce de uma das três origens, e cada kind tem seu próprio template
porque cada um pede um corpo diferente:

| kind | gatilho | o que o corpo descreve |
|---|---|---|
| `stale` | uma ponta da aresta avançou de rev; a revalidação não é automática | quem está atrás de quem (`anchor_rev` vs `target_rev`) |
| `conflict` | duas ou mais âncoras bloqueantes discordam do mesmo alvo | "âncora A diz X, âncora B diz Y — decida qual muda" |
| `violation` | o confronto rodou e o alvo **viola** a âncora | qual invariante da âncora foi quebrada |

### Issues são arquivos; o estado é a pasta

Uma issue é um arquivo markdown. Seu estado é a **pasta** onde ele está. Mover de
pasta = mudar de estado — um `git mv` visível no PR. O estado vive no filesystem,
não num banco: impossível dessincronizar do repositório.

```
issues/
  _templates/
    stale.md
    conflict.md
    violation.md
  todo/       ← detectada, ninguém pegou
  doing/      ← alguém está resolvendo
  done/       ← tratada (fato datado)
```

O nome do arquivo codifica data, kind e a aresta, para ser único e legível:

```
issues/todo/2026-08-05-a1--stale--login.spec.md--vs--login.tsx.md
```

Exemplo de issue `stale`:

```markdown
---
id: 2026-08-05-a1
kind: stale
edge:
  from: login.spec.md
  to:   login.tsx
  type: specifies
detected_at: 2026-08-05
detected_by: <confronto que abriu — lint | check | agente>
anchor_rev: 4
target_rev: 9          # o alvo avançou → spec ficou atrás
report: <ponteiro para o laudo completo, o "porquê">
---

## O que está dessincronizado
login.tsx (rev 9) foi alterado depois da última validação da spec (rev 4).
A spec descreve login por email; o código agora aceita telefone.

## Como resolver (uma das duas)
- [ ] **Corrigir o alvo** — reverter telefone em login.tsx (a spec manda)
- [ ] **Atualizar a âncora** — descrever telefone em login.spec.md
      (⚠ pode propagar stale para quem rege a spec)

## Resolução
<!-- preenchido ao mover para done/: qual ação foi tomada, por quem, quando -->
```

### Resolução é sempre uma das duas ações

Uma issue nunca é resolvida escrevendo "resolvido". Ela é resolvida executando
uma das duas ações do fluxo bidirecional (§2):

- **corrigir o alvo** → o alvo re-sincroniza com a âncora, ou
- **atualizar a âncora** → a âncora passa a refletir a nova realidade (e isso pode
  propagar stale para cima no grafo).

Em ambos os casos, um novo confronto roda, o carimbo da aresta alcança as revs
atuais, e a issue vai para `done/`.

### A issue é intervenção do usuário: três rotas

A issue é o ponto onde o **humano intervém** — o Anchors detecta e apresenta, mas
quem decide o que fazer é quem opera. Diante de uma issue, o usuário escolhe entre
**três rotas**:

1. **Resolver ele mesmo** — executa uma das duas ações à mão (o retoque manual).
2. **Delegar a um agente** — manda um agente executar a correção (o retoque
   automatizado).
3. **Converter em um plano** — quando a resolução é trabalho estruturado, não um
   retoque, a issue **realimenta o fluxo de desenvolvimento** virando um plano
   (`PLANNING.md`).

A conversão em plano **nunca é automática** — nem por tamanho, nem por heurística. É
uma escolha do operador, uma das três. Quando ele converte:

- a **issue encerra** (vai para `done/`) com a mensagem *"convertida no plano 00XX"*;
- o **plano nasce com o propósito** *"aberto para resolver a issue 00XX"* —
  referência cruzada bidirecional.

Isso **não** viola "done não pode mentir": a issue não está resolvida no sentido de
sincronizada — está **encaminhada**. O débito não sumiu; mudou de dono (agora o plano
o carrega e o decompõe), e o rastro mostra para onde foi. É a diferença entre
*resolver* e *transferir com registro*.

### Issues são imutáveis. `done` é done para sempre.

Uma issue é um **evento datado**: "nesta data, detectou-se esta dessincronia; foi
tratada assim". `todo → doing → done` é um caminho de mão única. Uma issue **nunca
reabre** — reabrir apagaria a história e faria `done` mentir sobre o passado.

Se a mesma dessincronia **volta** (o `done` era falso, ou o alvo regrediu depois),
isso não é a mesma issue de novo — é um **novo evento**. O próximo confronto
(lint, check ou agente) que detecta o conflito abre uma **issue nova** em `todo/`,
com nova data e novo id, opcionalmente apontando "reincidência de <id anterior>".

```
issues/done/2026-08-05-a1--stale--login.spec.md--vs--login.tsx.md   (tratada em 05/08)
                    │
              (regrediu / done era falso)
                    ▼
   próximo confronto detecta de novo
                    ▼
issues/todo/2026-08-09-c7--stale--login.spec.md--vs--login.tsx.md   (novo fato em 09/08)
```

Não existe estado "estamos convivendo com isso". Uma dessincronia que
conscientemente não vai ser resolvida agora simplesmente **fica em `todo/`,
incomodando** — e esse é o comportamento correto. O incômodo persistente é a
feature, não o bug: o único jeito de silenciá-lo é resolver, nunca arquivar.

Isso dá, de graça, um sinal de qualidade: como o nome do arquivo codifica a
aresta, varrer `done/` revela dessincronias **crônicas** ("esta aresta já gerou
três issues") — sinal de que a âncora pode estar mal desenhada.

---

## 5.1 Opt-out honesto: dispensar sem esconder

Há momentos em que uma exigência de validação **não** deve valer para uma peça: uma
spec declarativa que o time não quer testar, um bloqueio que precisa ser pulado desta
vez, um gate que ainda não é cobrado com força. O Anchors permite dispensar — mas só
de um jeito: o **opt-out honesto**.

Um opt-out honesto é uma **dispensa explícita, registrada e datada** de uma exigência
de validação. Tem três invariantes:

1. **Explícito** — alguém escreveu a marca (uma tag, uma flag, uma mudança de
   política). Nunca emerge de silêncio ou de esquecimento.
2. **Registrado** — dispensa *a exigência* (o teste, o bloqueio, o rigor), **nunca o
   registro**. O débito permanece visível. Pode-se optar por não bloquear, jamais por
   não registrar.
3. **Datado e justificado** — carrega *por quê* foi dispensado, quando e por quem, num
   artefato versionado. É auditável: dá para listar "tudo que dispensamos e a razão",
   e reavaliar se ainda faz sentido.

Isto é o que separa a **saída legítima** da validação do **buraco de cobertura**
(`QUALITY.md` §5.1). O buraco é o oposto ponto a ponto: *implícito* (ninguém decidiu),
*não-registrado* (nada aparece), *invisível* (não dá para achar). O opt-out honesto é
a porta da frente — decidida, à vista; o buraco é a janela quebrada.

Não confundir com "estamos convivendo com isso" (§5): aquilo é uma *dessincronia
detectada* que fica incomodando até ser resolvida. O opt-out honesto é anterior — é
dispensar a *exigência* de validar aquela peça, de modo que a dessincronia nem chega
a ser cobrada. Um dispensa o confronto; o outro é o confronto pendente.

O conceito é transversal — cada pilar o instancia sobre o que dispensa:

| instância | dispensa | onde | escopo |
|---|---|---|---|
| **`@no-test`** (`SPEC.md` §6) | a exigência de teste | tag na spec | permanente, por âncora |
| **`--no-block`** (`QUALITY.md` §7) | o bloqueio, desta vez | flag de invocação | pontual, uma execução |
| **maturação informativa** (`QUALITY.md` §7) | o rigor (mede, não trava) | política do gate | permanente, por gate |

Todas as três são o mesmo conceito — dispensa honesta — aplicado a coisas diferentes
(o teste, o bloqueio, o rigor). Em todas, a marca é explícita, o registro nasce, e o
porquê fica datado.

---

## 6. Dois planos: o vivo e o histórico

O Anchors mantém dois planos separados, com naturezas temporais opostas. Nunca se
contradizem porque cada um guarda uma verdade diferente.

| plano | natureza | guarda a verdade… | contém |
|---|---|---|---|
| **âncoras + grafo** | **vivo** — reflete o presente, reescreve | …*atual* da sincronia | specs, features, docs, guides, arestas, revs, carimbos |
| **registro** | **append-only** — fatos datados, imutáveis | …*histórica* dos eventos | issues resolvidas, laudos, patches |

A verdade *atual* sobre se algo está em sincronia vive no **grafo** (calculada de
revs e carimbos). A verdade *histórica* sobre o que aconteceu vive no
**registro** (issues e laudos datados, imutáveis).

Por isso `done` não precisa ser verdade eterna sobre o grafo — precisa ser verdade
sobre *o que aconteceu naquele dia*. A honestidade de uma issue não é "esta
dessincronia está resolvida para sempre" (impossível garantir); é "esta
dessincronia foi tratada nesta data" (fato imutável). Cada plano é honesto no seu
próprio eixo de tempo.

Misturar os dois planos — deixar decisões datadas apodrecerem dentro de um
artefato vivo, ou reescrever o histórico para refletir o presente — é o erro que
faz qualquer sistema de memória apodrecer. O Anchors os mantém rigorosamente
separados.

---

## 7. Resumo do modelo

- **Âncora** = documento que nos acompanha na escalada: **aponta** (diz para onde
  ir), **segura a corda** (confronta/bloqueia) ou **demarca** (deixa rastro).
  Definição generosa — toda âncora usada é âncora, inclusive a documentação (âncora
  de consumo, não pilar). "Guiar + confrontar" é a forma **bloqueante** (a que
  segura a corda), não o teste de existência. Anti-drift é lei: fundamentar na
  realidade, descrever o quê não o como, preservar o estável.

- **Grafo** = nós + arestas dirigidas, muitos-para-muitos, **material e
  versionado**. A aresta é **tipada** (o tipo carrega força + pergunta de
  confronto); a força é derivada do tipo. Sem hierarquia embutida — o tipo diz o
  que confrontar, nunca quem manda.

- **Sincronia** = carimbo de validação por **aresta** (não por nó), porque "estar
  em dia" é propriedade da relação. Como a onda percorre o mapa de dependências, faz
  a **análise de impacto** e produz o **caminho de impacto** (o mínimo a reconferir)
  é o pilar de **Propagação** (`PROPAGATION.md`).

- **Issue** = registro material de dessincronia, em **três kinds** (stale,
  conflict, violation), como **arquivos em pastas-estado** (todo/doing/done). Aberta
  só por confronto de aresta bloqueante. **Imutável e de mão única** — nunca reabre;
  recorrência é issue nova. Resolução = uma das duas ações do fluxo bidirecional.

- **Opt-out honesto** = dispensar uma exigência de validação de forma **explícita,
  registrada e datada** (`@no-test`, `--no-block`, maturação). Dispensa a exigência,
  nunca o registro. É a porta legítima; o oposto — dispensa implícita e silenciosa —
  é o buraco de cobertura (`QUALITY.md` §5.1).

- **Dois planos** = o **vivo** (âncoras + grafo, guarda a verdade atual) e o
  **histórico** (registro append-only, guarda os fatos datados). Separados por
  princípio; nunca se contradizem.

O que tudo isso entrega é a resposta à dor original: um projeto que cresce sem
perder o rigor. As sessões futuras não reconstroem o contexto — elas leem âncoras
que as guiam, confrontam o que fazem contra essas âncoras, e o custo desse rigor
permanece baixo porque a validação é sempre incremental. A direção persiste
porque está ancorada, e as âncoras não podem mentir.

---

## 8. Instâncias (não normativo)

O conceito acima é independente de ferramenta. Duas instâncias reais o exercitam:

- Uma ferramenta que implementa o Anchors com agentes de IA: watchers
  detectam alterações, agentes confrontam artefatos contra guides/specs, e toda
  alteração passa por PR. O "guide" dessa ferramenta é uma âncora; sua cadeia de
  reação é o ciclo guia+confronta; seus laudos e patches datados são o plano
  histórico. O que o Anchors adiciona a ela é o **grafo material de dependências**
  e a **sincronia incremental** — hoje ela reconfere por convenção de nomes, não
  por grafo com carimbo de frescor.

- Um workspace de prova de conceito real, que já pratica âncoras à mão
  (`guides/` com política anti-drift própria, memórias tipadas, planos com
  progresso). É onde o conceito é validado antes de ser generalizado.

Estas instâncias ilustram o conceito; não o definem. O Anchors permanece verdadeiro
sem elas.

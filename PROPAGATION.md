# Anchors — Pilar de Propagação

> Este documento define o **pilar de Propagação** do Anchors. Ele pressupõe o
> mecanismo geral de [`CONCEPT.md`](./CONCEPT.md) — âncora, grafo, issues — e o
> pilar de [`TRACEABILITY.md`](./TRACEABILITY.md), sobre o qual opera. A Propagação
> é o motor: é ela que faz uma alteração num ponto percorrer o organismo e é ela
> que faz o desenvolvimento avançar.
>
> É teoria de framework, já validada por uma ferramenta que implementa o
> Anchors e por uma prova de conceito real que exercita o conceito.

---

## 1. O motor, não o mapa

Os outros pilares descrevem **propriedades** do organismo: como as peças se
conectam (Rastreabilidade), se estão boas (Qualidade), se dizem a verdade
(o anti-drift, lei do `CONCEPT.md` §2). A Propagação é diferente — ela é o
**metabolismo**. Não descreve um estado; produz movimento. É o que faz o projeto sair de um estado coerente e
alcançar o próximo depois de cada mudança.

Se a **Rastreabilidade** é a trama que conecta as raízes da floresta, a
**Propagação** é o sinal que corre por essa trama — o nutriente que flui de uma
raiz até alcançar toda a rede. Uma é a fiação; a outra é a corrente passando por
ela. Uma é estrutura; a outra é vida.

A frase que define o pilar: **é a propagação das alterações que faz o
desenvolvimento avançar.** Sem ela, o grafo é um mapa parado de dependências. Com
ela, o projeto é um sistema que **se completa** a cada alteração — a mudança num
ponto empurra tudo que dela depende até o organismo voltar a ser coerente.

---

## 2. Sincronia é diagnóstico; propagação é movimento

Vale separar duas coisas que parecem uma só:

- **Sincronia** responde *"o que está desatualizado agora?"*. É um **estado**, uma
  foto. Diagnóstica.
- **Propagação** responde *"quando isto muda, o que mais fica desatualizado, e como
  a onda percorre o grafo?"*. É um **movimento**, um fluxo. Generativa.

A sincronia é o que a propagação **produz e lê** ao longo do caminho. A propagação
percorre o **mapa de dependências** que a Rastreabilidade mantém para *calcular* a
sincronia: em cada aresta que toca, determina se aquela relação ficou stale. O
carimbo de sincronia é estrutural (mora na aresta, `CONCEPT.md` §4); a propagação o
**recalcula** ao passar. Detecção de stale não é um mecanismo à parte — é o
**efeito** da onda percorrendo o mapa.

> **Escopo enxuto do pilar.** A Propagação **aponta**, não resolve. Seu trabalho
> termina em "estas relações ficaram stale e precisam de atenção". O que acontece
> depois — rodar os gates, confrontar, gerar issue — não é propagação; é a
> continuação do fluxo de trabalho (§6). A Propagação é *um passo* de um movimento
> maior, não o movimento inteiro.

---

## 3. Como a onda funciona

### O gatilho: uma revisão avança

Cada nó do grafo tem uma **revisão** (`rev`) do seu conteúdo. Quando um arquivo
muda, sua `rev` avança. Esse é o estímulo que inicia a onda.

### O carimbo vive na aresta

"Estar em dia" não é propriedade de um arquivo isolado — é propriedade de **uma
relação**. Uma spec pode estar em sincronia com o código que ela descreve e ao
mesmo tempo atrasada em relação à arquitetura que a rege. Dois vereditos, mesma
spec, duas arestas. Por isso o carimbo de validação pendura na **aresta**, não no
nó.

Os campos do carimbo (`validated_from_rev`, `validated_to_rev`, `changed_at`,
`verdict`) são definidos e mantidos pela Rastreabilidade, que é a dona do mapa —
ver a enumeração completa em `TRACEABILITY.md` §4. A Propagação os **lê** para
calcular staleness; é esse uso que este pilar descreve.

### A regra de staleness

Uma aresta está **stale** quando qualquer ponta avançou desde o último confronto,
ou quando nunca foi validada:

```
stale(edge) :=
     node[edge.from].rev > edge.validated_from_rev    -- a âncora mudou
  OR node[edge.to].rev   > edge.validated_to_rev       -- o alvo mudou
  OR validação ausente                                 -- nunca validada
```

```
código           spec (âncora)
 rev 7   ◄─specifies─  validated_to_rev 7   → IN SYNC  ✔
 rev 9   ◄─specifies─  validated_to_rev 7   → STALE    ⚠  (alvo 2 revs à frente)
```

### A análise de impacto: o caminho de impacto, nunca o projeto todo

Quando um nó avança de rev, a Propagação faz a **análise de impacto** — percorre o
mapa de dependências da Rastreabilidade a partir do que mudou e computa
**exatamente** o que ficou stale. O resultado é o **caminho de impacto**: o
conjunto exato de confrontos que precisam rodar, e só ele. Revalidar o projeto
inteiro a cada mudança é caro; caro não é feito; não feito vira drift. O rigor só
escala se a onda tocar **o mínimo necessário**.

É o mesmo princípio de um build system incremental (Make, Bazel): não recompile
tudo, recompile só o que ficou desatualizado em relação às suas dependências. A
análise de impacto:

- **Alcança a vizinhança direta primeiro** — todas as arestas onde o nó alterado é
  `to` (quem o rege precisa reconferir) *e* onde é `from` (o que ele rege pode ter
  ficado stale). As duas direções, só os vizinhos imediatos.
- **Propaga em largura, podando o que não ficou stale** — uma busca pelo grafo a
  partir do nó alterado, que **para** em toda aresta que *não* ficou stale. A onda
  anda só pela frente das mudanças; nunca varre o repositório.

### A onda sobe e desce pelo grafo

A propagação segue a direção das dependências, e isso naturalmente a faz **subir e
descer**:

```
código (rev↑)
   │ specifies
   ▼
 spec  ── fica stale ──►  é atualizada (rev↑)
   │ governs                      │
   ▼                              ▼ agora a spec avançou
arquitetura / guide  ── pode ficar stale ──►  reconfere
```

Mexeu no código → a spec que o descreve fica stale. Quando a spec é atualizada
para reabsorver a mudança, a *própria spec* avança de rev — e então quem **rege a
spec** (a arquitetura, um guide) pode ficar stale por sua vez. A onda sobe. Mas só
os nós realmente tocados entram na fila: um ramo do grafo que a mudança não
alcançou permanece em sincronia e é podado.

### O tamanho da onda emerge do grau do nó

Uma mesma máquina de propagação produz ondas de tamanhos radicalmente diferentes —
e o tamanho **não** é um segundo mecanismo, é consequência de **quantas arestas o nó
alterado tem**:

- Mudar uma spec de tela (um nó de **baixo grau de saída**, que rege só seu código,
  feature e teste) gera uma onda **local**: alcança o punhado de derivados
  co-localizados e assenta.
- Mudar um **guide** ou uma spec de arquitetura (um nó de **alto grau de saída** —
  *fan-out* —, que *rege* dezenas ou centenas de arquivos por arestas `governs`) gera
  uma onda **global**: quando a régua muda, *todos* os nós que ela rege ficam stale de
  uma vez, e a onda é um fan-out por todo o domínio.

> O grau que importa aqui é o de **saída** (quantas arestas `governs` partem do nó).
> É diferente do "nó de maior grau de convergência" que a Spec é (`SPEC.md` §2) —
> lá, muitas encarnações *pendem* de uma spec (grau de entrada). Fan-out e
> convergência são dimensões distintas do grau.

A distinção "mudança de dado" (local) vs. "mudança de critério/régua" (global) que
o POC materializa é, no Anchors, **a mesma onda vista em dois nós de grau diferente**.
Não há propagação especial para guides — há um nó que rege muitos, e a propagação
faz o que sempre faz: alcança quem ficou stale. É por isso que mudar um guide é
**caro**: o custo vem do fan-out — uma mudança de régua reconfere o domínio inteiro.
(E é por isso que a **maturação** — decidir *quando* o gate passa a bloquear,
`QUALITY.md` §7 — pesa tanto para um gate que confronta um guide: bloqueia
proporcionalmente ao fan-out.)

---

## 4. A falha característica: a onda que morre no meio

Cada pilar tem uma falha que é sua. A da Propagação é a **onda incompleta** — a
ripple que para antes de alcançar tudo que deveria.

O resultado é um projeto num estado **meio-coerente**: parte alcançou o novo
estado, parte ficou no antigo. E é a falha mais insidiosa de todas porque é
**invisível ponto a ponto** — cada peça, olhada isoladamente, parece ok. Só quem
enxerga a onda inteira percebe que a mudança não terminou de se espalhar.

Essa falha é distinta de todas as outras:

- não é falha de **Rastreabilidade** — as arestas existem, as conexões estão lá;
- não é falha de anti-drift — cada âncora, isolada, pode estar dizendo a verdade;
- não é falha de **Qualidade** — cada gate local pode estar passando.

É falha de **completude da propagação**: a alteração não terminou de percorrer o
grafo. O desenvolvimento parou no meio de uma transição sem que ninguém percebesse.
Sem um pilar que vigie a onda como um todo, o projeto acumula transições
inacabadas até virar um mosaico de estados inconsistentes que "localmente parecem
certos".

Por isso a Propagação precisa de pilar próprio: sua falha só é detectável por quem
olha o movimento inteiro, não as peças.

---

## 5. Quiescência: como se sabe que uma mudança terminou

A Propagação traz um critério objetivo que o Anchors não teria sem ela: **quando
uma alteração está de fato completa.**

Não é quando você salvou o arquivo. É quando a onda **assentou** — quando percorreu
o grafo e não há mais nenhuma aresta stale que a mudança tenha deixado para trás.
Esse estado de repouso é a **quiescência**.

```
alteração → onda percorre → arestas ficam stale → são reconferidas/atualizadas
          → geram novas ondas → ... → nenhuma aresta stale restante
          → QUIESCÊNCIA (a mudança terminou de se propagar)
```

> **Fronteira importante.** A quiescência *da propagação* significa apenas "a onda
> de stale-ness assentou" — nada mais está desatualizado por *dependência*. Isso
> **não** é o mesmo que "o trabalho terminou". Depois que a propagação assenta, o
> fluxo continua: os gates de qualidade executam, podem achar divergências, gerar
> issues, e iniciar novas alterações — que propagam de novo. A quiescência
> *completa do projeto* (onda assentada **e** gates passados **e** zero issues
> abertas) é uma propriedade do **fluxo de trabalho** (§6), não da Propagação
> isolada. A Propagação garante só a sua parte: nada ficou stale por dependência.

---

## 6. Propagação dentro do fluxo de trabalho

A Propagação é **um passo** de um movimento maior. O fluxo de trabalho é o desenho
que interliga os passos — e a Propagação é o primeiro deles, não o único:

```
   alteração
      │
      ▼
 ┌──────────────┐
 │ PROPAGAÇÃO   │  ← análise de impacto sobre o mapa da Rastreabilidade
 │ (este pilar) │     marca o que ficou stale — APONTA, não resolve
 └──────┬───────┘
        │
        ▼
 ┌──────────────┐
 │ GATES        │  ← executam, medem, confrontam
 │ (Qualidade,  │     NÃO é propagação — é a continuação do fluxo
 │  anti-drift) │
 └──────┬───────┘
        │
   divergência? ──► ISSUE ──► nova alteração ──┐
        │                                      │
        └──────────────────────────────────────┘
                    (e a nova alteração PODE propagar de novo)
```

Quando um gate é chamado, aquilo **não é mais propagação de alteração** — é
execução, um passo de natureza diferente. A Propagação entrega o gatilho ("estas
arestas ficaram stale, confronte-as") e sai de cena; os gates assumem. Se um gate
acha divergência, gera uma issue; a issue pode iniciar uma nova alteração; e essa
alteração entra de novo no topo da onda. O ciclo é o fluxo; a Propagação é o
trecho que espalha.

### Barreiras intermediárias: o gate entre os elos, não só no fim

O diagrama acima simplifica ao pôr "os gates" como um bloco único depois da
Propagação. Na prática, quando a onda percorre uma **cadeia** de âncoras (spec →
código → feature → teste), o gate pode ficar **entre cada elo**, não só no fim. Um
fluxo pode exigir que a propagação de spec→feature só **avance** para feature→teste
quando a feature foi confrontada e passou — uma **barreira** por elo.

```
spec ──►│gate│──► feature ──►│gate│──► teste ──►│gate│──► ...
        (avança só se passar)  (idem)            (idem)
```

Isso não muda o que a Propagação *é* — ela continua só apontando o stale. Muda o
*desenho do fluxo*: onde as barreiras ficam é escolha do fluxo (que não é pilar).
O POC materializa isso pondo um gate (aprovação/merge) entre cada elo da cadeia — a
onda de um elo só dispara o próximo quando o anterior é aceito. O Anchors reconhece
que a propagação pode ser **gated por elo**, não só ao final.

> Esta é a posição da barreira *na cadeia* (entre-elos vs. no-fim) — eixo ortogonal
> à distinção *local vs. promoção* de `QUALITY.md` §8, que é sobre o *momento* da
> imposição (pré-registro vs. pré-integração). Um não particiona o outro: pode-se ter
> barreiras entre-elos tanto no gate local quanto no de promoção.

### Propagação inversa: a onda que gera issue, não mutação

A onda não corre só no sentido "âncora → alvo" (spec desce para código). Ela também
corre **inverso**: quando o *alvo* muda fora do fluxo (alguém edita o código direto,
sem tocar a spec), a aresta `specifies` fica stale pela ponta de baixo, e a onda
sobe até a spec.

Mas a propagação inversa tem um comportamento próprio: ela **não auto-resolve**.
Diferente da onda direta (onde a spec *manda* e o código a segue), aqui o framework
não sabe se o código divergente está certo (a spec ficou velha) ou errado (o código
violou a spec). Então a onda inversa **para e gera uma issue `stale`** — a ponta de
baixo avançou de rev, e a revalidação não é automática (`CONCEPT.md` §5). Não é uma
issue `violation`: a propagação inversa *para antes de confrontar*, então o
framework ainda não sabe se houve violação — só registra que as pontas
dessincronizaram. A decisão entre "corrigir o alvo" e "atualizar a âncora" é de quem
opera (`CONCEPT.md` §2, o fluxo bidirecional). É uma onda que produz *dessincronia
registrada*, não *mutação automática*. O Anchors nunca reescreve uma âncora para
legitimar uma mudança de alvo sem decisão humana.

### Como a onda se materializa: um watcher que sabe o mapa e a inter-relação

Abstratamente, a onda "percorre o grafo". Concretamente, isso se materializa num
**watcher** — e o ideal não é um enxame de watchers cegos (um por tipo de arquivo),
mas **um watcher inteligente** que conhece três coisas:

- **o grafo virtual** (`STRUCTURE.md` §2.1) — pelos patterns de arquivo/pasta, sabe
  *a que camada* o arquivo alterado pertence e que estrutura ele deveria ter, mesmo
  quando o mapa material ainda está vazio (o bootstrap);
- **o mapa de dependências** (`TRACEABILITY.md` §4) — ao ver que um arquivo mudou,
  sabe *o que de fato depende dele*, sem precisar de padrões separados por tipo;
- **a inter-relação dos pilares** — sabe *qual agente* chamar para aquela alteração
  (spec mudou → agente de propagação, que cria código + feature; feature mudou →
  agentes de teste e de doc).

O watcher é a análise de impacto encarnada: lê o mapa para saber *o que* ficou stale,
e a inter-relação dos pilares para saber *quem* despachar. A onda vira então uma
cascata — cada agente disparado cria/atualiza artefatos, o que o watcher percebe e
roteia para o próximo agente:

```
alteração de arquivo
   │  watcher consulta o mapa + a inter-relação dos pilares
   ▼
spec criada/alterada   → agente de propagação → cria código + feature
   │
feature criada         → agente de testes  (cria/atualiza testes)
   │                    → agente de docs    (gera doc p/ devs + stakeholders — elo terminal)
   ▼
(a cada elo) o agente ATUALIZA O MAPA — registra nós/arestas criados, movidos, removidos
```

Dois pontos que fecham as costuras:

- **Cada agente cria E mapeia.** A propagação não só *detecta* stale — o agente
  disparado **cria** os artefatos faltantes (o código, a feature, o teste) e, no
  mesmo ato, atualiza o mapa (a lei de manutenção, `TRACEABILITY.md` §4). É assim que
  o grafo *cresce*, não só é percorrido.
- **A doc é o elo terminal.** O agente de documentação é disparado no fim da cascata,
  gera a doc de consumo para stakeholders (devs e usuários), e a onda assenta ali — a
  doc é folha do grafo, não dispara nada rio acima (`CONCEPT.md` §2).

Isto é *uma* materialização (a que o POC adota); o Anchors define a onda em abstrato,
e o watcher-que-sabe-o-mapa é a forma concreta de executá-la. Outra ferramenta
poderia disparar de outro jeito — o conceito de propagação não depende do watcher.

O **fluxo em si não é um pilar** — é um desenho de como as coisas se interligam,
que pode mudar ou ter várias perspectivas. Os componentes que operam dentro dele
(Estrutura de Projeto, Planejamento, Spec, Rastreabilidade, Propagação, Qualidade)
é que são os pilares: sem eles, o fluxo seria só um fluxo. A Propagação é um desses
componentes — o que dá movimento ao desenho.

---

### Regras de convivência entre workers

A fila já garante EXCLUSIVIDADE: o claim é um rename atômico, e dois workers nunca pegam
a mesma task. Isso resolve a disputa pela *task* — e não resolve a convivência, que é
outro problema: **o disco é compartilhado, e o Anchors não intermedia a escrita.**

O modo de falha real, medido: um subagente rodou 90 minutos numa única etapa. A quem
observava de fora, pareceu travado (nenhum arquivo mudando). Alguém deu `reclaim`,
disparou outro worker, e os dois passaram a escrever no mesmo arquivo. Deu certo por
sorte — as regras acrescentadas por um não colidiram com as do outro.

Repare que nada disso é falta de trava. É falta de três combinados:

**1. `reclaim` respeita quem pegou há pouco.** Devolver "toda task claimed" pressupõe que
claim parado significa worker morto — e não significa: significa worker TRABALHANDO. O
critério é TEMPO, não processo vivo, porque *o worker do Anchors não é um processo*: o
`anchors next` imprime a task e morre; quem trabalha depois é o agente, que o CLI não
observa. Checar o PID do `next` mediria sempre "morto" e voltaria a roubar trabalho ativo.

A janela é generosa (horas, não minutos) porque os custos são assimétricos: devolver cedo
duplica trabalho e só se descobre depois; devolver tarde custa espera, e `--force` resolve
para quem tem certeza.

**2. Trabalhar FORA da fila é o que quebra o combinado.** Quem pega trabalho direto — sem
`next`, sem claim — não aparece no `queue`, e nenhum outro worker tem como saber. A fila
não é burocracia: é o único lugar onde "estou nisto" fica visível. Um orquestrador que
delega deve fazê-lo pela fila, não por instrução direta ao subagente.

**3. Entrega registrada é o que torna o trabalho auditável depois.** O `changes/` responde
"quem fez o quê" quando dois agentes tocaram a mesma unidade — e foi assim que a colisão
acima foi diagnosticada, não por log de processo.

## 7. Relação com os outros pilares

- **Opera sobre a Rastreabilidade** (`TRACEABILITY.md`): a onda só percorre o que
  está conectado. A Propagação *usa* a cola rastreável para saber por onde se
  espalhar — sem identidade contínua, a onda não tem trilhos. Se a Rastreabilidade
  é a fiação, a Propagação é a corrente.
- **Segue a Estrutura de Projeto** (`STRUCTURE.md`): a *ordem* em que a onda sobe e
  desce pelas dependências é a ordem que a Estrutura declara. A Rastreabilidade diz
  quais peças se ligam; a Estrutura diz em que ordem as camadas dependem umas das
  outras — e a onda respeita essa planta.
- **Alimenta a Qualidade** (`QUALITY.md`): o caminho de impacto que a Propagação
  calcula é exatamente o conjunto de gates que o pipeline de qualidade precisa
  rodar. A Propagação diz *o que* reconferir; a Qualidade *reconfere*.
- **Alimenta o anti-drift** (`CONCEPT.md` §2): a onda aponta quais âncoras podem ter
  passado a mentir (o alvo mudou, a âncora não acompanhou) — dando ao confronto de
  anti-drift o seu caminho de impacto.
- **Gera as issues via o mecanismo comum** (`CONCEPT.md` §5): quando o confronto
  disparado pela onda falha, o resultado é uma issue material — o mesmo modelo de
  registro do resto do framework.

A Propagação é o pilar que **conecta os outros no tempo**: ela é o que transforma
"aqui está um conjunto de propriedades desejáveis" em "aqui está um sistema que se
mantém coerente enquanto muda".

---

## 8. Resumo do pilar

- **Propagação = o motor.** Não descreve um estado; produz movimento. É a
  propagação das alterações que faz o desenvolvimento avançar.
- **Sincronia é diagnóstico; propagação é movimento.** A propagação percorre o mapa
  de dependências da Rastreabilidade para *calcular* a sincronia — detectar stale é
  o efeito da onda passando, não um mecanismo separado.
- **Análise de impacto = caminho de impacto.** Um nó avança de rev → só as arestas
  realmente tocadas ficam stale → busca em largura podando o que não mudou. Nunca o
  projeto todo. Incremental como um build system.
- **A onda sobe e desce** pelo grafo, seguindo as dependências, mas só pelos ramos
  que a mudança alcançou.
- **O tamanho da onda emerge do grau do nó.** Mudar uma spec de tela (baixo grau) =
  onda local; mudar um guide/arquitetura (alto grau, rege muitos) = onda global. Não
  é mecanismo especial — é a mesma onda alcançando quem ficou stale.
- **Falha característica = onda incompleta.** Estado meio-coerente, invisível ponto
  a ponto. Só quem vê o movimento inteiro detecta.
- **Barreiras e inversão.** O fluxo pode pôr um gate *entre cada elo* da cadeia
  (avança só se passar), não só no fim. E a onda inversa (alvo muda fora do fluxo)
  **para e gera issue** em vez de auto-mutar — a decisão é humana.
- **Quiescência** = a onda assentou (nada stale por dependência). É o critério de
  "a alteração terminou de propagar" — mas **não** de "o trabalho terminou": isso é
  do fluxo, que continua pelos gates.
- **Escopo enxuto: aponta, não resolve.** Entrega o caminho de impacto e sai de
  cena; os gates assumem. É *um passo* do fluxo, não o fluxo.
- **O fluxo não é pilar** — é o desenho; os pilares são os componentes que operam
  nele.

O que o pilar entrega: um projeto que **evolui como um organismo**. Cada alteração
não fica presa onde nasceu — ela percorre a rede, torna visível tudo que precisa
acompanhar, e o sistema se completa. É o que impede um projeto de virar um
acúmulo de mudanças locais que nunca terminam de se integrar ao todo.

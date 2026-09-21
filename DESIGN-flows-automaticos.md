# Os gatilhos que não são comando digitado — levantamento

> Correção ao `DESIGN-flows-do-anchors.md`. Aquele levantamento mapeou os 50 comandos e
> parou aí. Falta uma categoria inteira: o que dispara SOZINHO.

## A lacuna

Os cinco fluxos montados têm gatilhos que alguém digita: `init`, `next`, `check`, a
conversa. Mas boa parte do que o Anchors faz num projeto real não passa por ninguém
digitando nada.

O pre-commit roda quando você faz `git commit`. O CI roda quando você abre um PR. O vigia
enfileira quando um arquivo muda. Nenhum desses aparece nos fluxos que escrevi, e os três
têm DECISÃO — o pre-commit barra o commit, o CI barra o PR, o vigia escolhe qual tarefa
enfileirar.

A consequência é a mesma que motivou o conceito: quem trabalha encontra um barramento que
não está em fluxo nenhum, e improvisa.

## As quatro superfícies

### 1. O PRE-COMMIT (`anchors install-hooks`)

Dispara em `git commit`, sobre os arquivos em stage. E tem uma decisão fina, já escrita no
`--help`, que ninguém memoriza:

- arquivo NÃO REGIDO (não casa camada — `package.json`, lockfile, CI): **ignorado**, o
  Anchors não tem jurisdição;
- arquivo REGIDO mas FORA DO MAPA (novo, e o `map build` não rodou): **BARRA** — fora do
  mapa nenhum gate o confronta, e o commit passaria sem confronto nenhum;
- gate bloqueante vermelho: **barra**.

O terceiro caso é o esperado. O segundo é o que pega quem não conhece, e é exatamente um
resultado que merece peça.

### 2. O CI (`.github/workflows/ci.yml`)

Dispara em `pull_request`. Oito passos, e os três últimos são do Anchors sobre o Anchors:

    Compila · Vet · Formatação · Testes · Instala o anchors desta árvore
    O mapa está em dia?              ← barra o PR se o mapa está velho
    Os gates do Anchors, sobre o Anchors
    O check deixou issue para trás?  ← barra se o check abriu issue e ninguém tratou

Os dois marcados são decisão, não verificação mecânica: o primeiro manda rodar `map build`
e commitar; o segundo diz que há trabalho aberto que o PR não fechou.

### 3. O RELEASE (`.github/workflows/release-cli.yml`)

Dispara em `push` de tag `v*`. Não tem ramo — é linear — mas é um gatilho que ninguém
digita como comando do Anchors, e o fluxo de entrega termina nele.

### 4. O VIGIA (`anchors watch`)

Dispara em mudança de arquivo. Não executa nem chama IA: transforma "mudou" em "há
trabalho", e enfileira com a próxima etapa sugerida.

É o que sustenta o ciclo do `worker` — o passo `WORKR-P06` (fechar) diz "o vigia já
enfileirou a próxima etapa", e o vigia não está em fluxo nenhum. Ele é a peça que faz o
ciclo se fechar sozinho, e está invisível.

## A lacuna MAIOR: o resultado que gera trabalho novo

As quatro superfícies acima são o caso fácil — falta gatilho, e gatilho é peça. A lacuna
que importa é outra, e ela atravessa tudo o que já escrevi.

**Um resultado não termina o passo: ele muitas vezes CRIA trabalho que sobrevive à rodada.**

O caso que a expõe é o `check`. No fluxo `worker` eu escrevi:

    ACHCK-R02  BARRADO → WORKR-P03 (volta a escrever)

Como se corrigir fosse imediato e o assunto morresse ali. Não morre: cada gate que reprova
**ABRE UMA ISSUE**, e a issue tem ciclo de vida próprio, com quatro estados declarados no
código:

    future   dívida assumida, com prazo — ainda não vence
    todo     detectada, ninguém pegou
    doing    alguém está resolvendo
    done     tratada (fato datado)

A doutrina do `future` diz por que a distinção existe: *"quem olha `todo/` está perguntando
'o que faço agora', e afogar essa lista com o que só vence depois é o caminho mais curto
para ninguém mais olhar"*.

Nada disso está nos cinco fluxos. O `worker` volta a escrever e segue; a issue fica.

### E não é só o `check`

O mesmo padrão aparece em vários lugares, e em todos eu modelei o resultado como se fosse
o fim:

| ação | resultado | o trabalho que ele CRIA |
| --- | --- | --- |
| `check` | BARRADO | uma issue por gate, com ciclo próprio |
| `check` | JULGAMENTO PENDENTE | um alvo enfileirado para a IA julgar |
| `judge --verdict` | REPROVADO | a issue com o relatório inteiro no corpo |
| `escalate` | AGUARDANDO O USUÁRIO | um card que para a fila |
| `unblock` | CARD ABERTO | um card que o bloqueado espera |
| o vigia | arquivo mudou | uma tarefa na fila |
| o CI | issue para trás | o PR barrado até alguém tratar |

Sete lugares onde o resultado não fecha o passo — ele **ramifica**, e o ramo novo tem dono,
estado e fim próprios.

### E a terceira categoria: a reação SUGERIDA

Nem todo resultado gera trabalho automaticamente. Mas muitos SUGEREM uma reação — e a
sugestão hoje vive na prosa da mensagem, onde depende de alguém ler e lembrar.

O veredito real do `spec-feature-match`:

    1 declared requirement(s) without scenario in feature: TRCMT-B08.
    Write scenario, OR waive on requirement line with `@no-scenario: <reason>`

Duas reações. Nenhuma automática, nenhuma arbitrária — e nenhuma no fluxo. Medido no
catálogo de mensagens: **18 ocorrências** de "run `anchors <comando>`" dentro de vereditos
de gate, das quais `doctor --fix` (5), `map build` (3), `check --fix` (2) e `ingest` (3).

Cada uma dessas é uma ARESTA DE FLUXO escondida em prosa. Quem lê a mensagem sabe o que
fazer; quem não lê, improvisa — e é o mesmo problema que o conceito inteiro existe para
resolver, um nível abaixo.

**As três categorias, e a diferença que importa:**

| categoria | exemplo | quem age |
| --- | --- | --- |
| reação AUTOMÁTICA | `check` reprova → issue nasce | a máquina, sozinha |
| reação SUGERIDA | "escreva o cenário, ou dispense com `@no-scenario`" | quem trabalha, orientado |
| sem reação | `next` diz "fila vazia" | ninguém — acabou |

A do meio é a mais comum e a menos modelada. E ela tem uma propriedade que as outras não
têm: **costuma ser mais de uma**, e escolher entre elas é a decisão. "Escreva o cenário" e
"dispense com razão" levam a lugares diferentes, e a escolha depende de saber se o
requisito é real.

### O que isso significa para o modelo
### O que isso significa para o modelo

Hoje uma transição diz "deste passo, com este resultado, vá para aquele passo". Falta
dizer que um resultado também pode **abrir um trabalho paralelo**, que segue vivo depois de
o passo terminar.

Dois candidatos, e a escolha muda o que o `flow next` responde:

Agora são DOIS problemas, não um:

- o resultado que **gera** trabalho paralelo (automático);
- o resultado que **sugere** uma ou mais reações (a escolha é de quem trabalha).

Para o segundo, o formato já quase serve: uma sugestão é uma saída como as outras, com a
diferença de que quem decide é humano. O que falta é dizê-lo — uma saída sugerida não é
uma ordem, e apresentá-la como transição comum faria o fluxo mentir sobre quem manda.

**(a) o resultado aponta para outro FLUXO.** `ACHCK-R02` abriria o fluxo `issue`
(`todo → doing → done`), e o worker seguiria em paralelo. Simples, e reusa o que existe —
um fluxo já encaixa como peça (`WORKR-P05` encaixa o `JUDGE` inteiro).

**(b) a transição ganha um tipo**: `→` para "vá para", e algo como `⊕` para "isto também
abre". O `flow next` passaria a responder duas coisas: para onde EU vou, e o que ficou
aberto atrás de mim.

A (a) é menos mecanismo e responde a maior parte. A (b) diz mais, e custa um conceito novo
no formato.

## O que muda no modelo

Nada estrutural, e isso é bom sinal: uma ação já declara **o que faz e o que oferece**, e
um gatilho automático é uma ação cujo comando não é digitado.

O que falta é DISTINGUIR os dois, porque a diferença muda o que quem lê precisa saber:

    ## Command          → o que se digita
    ## Trigger          → o que dispara, e quando

Um `## Trigger` diz "roda em `git commit`, sobre os arquivos em stage" em vez de uma linha
de comando. Os resultados continuam iguais.

## O que entra

1. a seção `## Trigger` nas ações que disparam sozinhas
2. quatro peças novas: `pre-commit`, `ci`, `release`, `watch`
3. o fluxo `entrega` — do trabalho pronto ao release, atravessando commit, PR e tag
4. o `worker` passa a encaixar o vigia no passo que fecha o ciclo

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |
| `Q01` | O fluxo de ENTREGA é um só (commit → PR → merge → tag → release) ou dois (o do autor, que vai até o PR; e o do pipeline, que vai do merge ao release)? São donos diferentes — uma pessoa e uma máquina — e misturá-los pode esconder de quem espera onde a bola está. | você | um ou dois arquivos de fluxo |
| `Q03` | Resultado que gera trabalho: aponta para outro FLUXO (a), ou a transição ganha um TIPO que distingue "vá para" de "isto também abre" (b)? | você | se o formato ganha um conceito novo |
| `Q02` | O CI de um projeto que ADOTA o Anchors é peça semeada ou exemplo no guia? O `ci.yml` daqui é do Anchors sobre o Anchors; num projeto que o adota, os passos são outros. | você | se `flows/actions/ci.action.md` nasce com o `init` |

# Feature flags — o cenário por valor, e o que ele governa

> IMPLEMENTADO. Este documento fica como o registro do desenho e das decisões —
> inclusive das que mudaram no caminho.

## O problema

Uma feature flag multiplica os caminhos do código sem multiplicar a spec. `if
flag('novo-checkout')` cria dois comportamentos, e a spec descreve um — ou, pior, descreve
os dois misturados numa frase que não diz qual vale quando.

O custo aparece depois, em três lugares:

- **na revisão** — ninguém sabe qual ramo é o atual e qual é o que vai morrer;
- **no teste** — cobre-se o caminho ligado e o desligado fica sem prova;
- **na remoção** — a flag "temporária" fica três anos, e removê-la vira arqueologia.

## O desenho

Uma pasta `flags/` no root, um arquivo por flag, e dentro dele os CENÁRIOS por valor:

    flags/<nome>.flag.md

Cada cenário é uma linha de tabela com código próprio, na mesma gramática das regras:

    ## Cenários

    | Cenário | Quando o valor | Então |
    | --- | --- | --- |
    | `CHKUT-G01` | `= "off"` | o checkout antigo responde, e o novo não é chamado |
    | `CHKUT-G02` | `= "on"` | o checkout novo responde, com o mesmo contrato |
    | `CHKUT-G03` | `>= 50` | a porcentagem de rollout decide por usuário, de forma estável |
    | `CHKUT-G04` | ausente | vale o default declarado no cabeçalho — nunca um erro |

O `CHKUT-G04` é o que mais se esquece e o que mais quebra: a flag ausente (serviço de
flags fora do ar, ambiente novo, teste local) tem de ter comportamento declarado. Hoje
isso é descoberto em produção.

## A letra: `G`

Livre. As canônicas são `SRVAXBNMDEIQF`, e `G` não está entre elas.

**Alerta que o próprio código deixa:** o comentário do `DefaultRuleLetters` registra que
**três vezes** a lista ficou para trás de uma letra nova — *"são as duas metades da mesma
decisão, e elas divergiram sem que nada acusasse"*.

E são **três** metades, não duas: além da lista de letras e do catálogo de seções, o
`internal/testsig/code.go` carrega uma CÓPIA da lista (o pacote não importa `config` de
propósito), e o comentário dela registra que ficou *quatro* atrás. As três foram
atualizadas na mesma mudança, e os dois testes de guarda foram confrontados por mutação:
desalinhar qualquer uma reprova em `TestRuleLetters_naoDivergeDoConfig` ou em
`TestCatalogoNaoUsaLetraForaDasCanonicas`. A quarta letra não repetiu a história das três.

## A aresta: `gated-by`

A regra da SPEC declara que só vale sob um cenário de flag:

    ### CRED-V01 — valida o limite antes de submeter   @gated-by CHKUT-G02

Direção igual à do `realizes`: quem sabe que depende da flag é a spec, e um arquivo de
flag que listasse seus dependentes viraria índice que envelhece a cada spec nova.

## Os gates

| gate | pergunta | por que importa |
| --- | --- | --- |
| `flag-scenario-exists` | o `@gated-by` aponta para um cenário que existe? | o `ref-resolves` deste eixo |
| `flag-scenarios-complete` | a flag declara o caso AUSENTE? | é o que quebra em produção, e ninguém escreve |
| `flag-covered` | cada cenário tem teste que o exercita? | o caminho desligado é o que fica sem prova |
| `flag-stale` | a flag passou do prazo declarado? | é o que faz a flag temporária virar permanente |

O `flag-stale` depende de a flag declarar quando deve morrer — e isso é uma decisão de
formato: `expires: 2026-12-31` no cabeçalho, ou `@TBD` com prazo, que já existe.

## O que NÃO fazer

- **Ler o valor real da flag.** O Anchors não tem acesso ao serviço de flags, não deve ter,
  e o valor muda por usuário e por minuto. O que ele confronta é o CENÁRIO declarado.
- **Ditar a biblioteca.** `LaunchDarkly`, `Unleash`, um `if` num arquivo de config — o
  Anchors reconhece o cenário declarado, não a chamada. Quem sabe a forma da chamada no
  dialeto local é o projeto, como já acontece em `handle_patterns`.

## Decisões — todas fechadas

| Código | Decisão | Como ficou |
| --- | --- | --- |
| `Q01` | KIND novo, seguindo o precedente do `product`. | `mapx.KindFlag` + aresta `EdgeGatedBy`, na mesma mecânica do `realizes`: índice por código no `map build`, e cenário inexistente NÃO vira aresta morta — quem reporta é o `flag-scenario-exists`, com o código na mão. |
| `Q02` | GRAMÁTICA FIXA, a partir de um benchmark das ferramentas do mercado. | `internal/flagx/grammar.go`. Ver a seção abaixo. |
| `Q03` | Teste por CENÁRIO. | `flag-covered` cobra cada cenário contra os `ProvenCodes` do mapa. INFORMATIVO, não bloqueante: o cenário escrito hoje e testado no commit seguinte é trabalho normal, não defeito. |
| `Q04` | FORA desta rodada. | O prazo de morte da flag depende de um campo novo no cabeçalho e é decisão independente. O eixo nasce com quatro gates; o `flag-stale` entra quando o campo existir. |

## A gramática, e o benchmark que a produziu

A coluna "Quando o valor" é PARSEADA, não prosa. Prosa cobre qualquer caso e não dá para
confrontar — e um gate que não confronta é documentação com mais passos.

O conjunto de operadores não foi inventado. É a união do que as ferramentas do mercado
oferecem, lida da documentação de cada uma:

| família | LaunchDarkly | Unleash | Flagsmith |
| --- | --- | --- | --- |
| igualdade | `in` | `IN` / `NOT_IN` | `Equal` / `Not Equal` |
| ordem | `lessThan`, `greaterThanOrEqual`… | `NUM_EQ`, `NUM_GTE`… | `Greater Than`, `Less Than` (+ `Inclusive`) |
| texto | `startsWith`, `endsWith`, `contains` | `STR_STARTS_WITH`, `STR_ENDS_WITH`, `STR_CONTAINS` | `Contains`, `Not Contains` |
| regex | `matches` | `REGEX` | `Regex` |
| data | `before`, `after` | `DATE_BEFORE`, `DATE_AFTER` | — |
| semver | `semVerEqual`, `semVerLessThan`… | `SEMVER_EQ`, `SEMVER_GTE`… | — |
| presença | — | — | **`Is Set` / `Is Not Set`** |
| rollout | percentage rollout | rollout % | `Percentage Split` |
| conjunto | `segmentMatch` | — | `Modulo` |

Três achados desse cruzamento decidiram o desenho.

**Primeiro: as três discordam apenas na GRAFIA.** Todas têm igualdade, ordem, substring e
regex. Então a gramática guarda um operador por SIGNIFICADO e aceita as grafias como
alias: quem escreve `>=` e quem escreve `NUM_GTE` está fazendo a mesma afirmação, e o gate
não deve se importar com qual delas a pessoa aprendeu.

**Segundo: só a Flagsmith tem `Is Set`/`Is Not Set` — e é o operador de que este framework
mais precisa.** A flag ausente (serviço fora do ar, ambiente novo, teste local) é o caso
que quebra em produção e o que ninguém escreve. Todas as ferramentas conseguem EXPRESSAR
esse caso; nenhuma o EXIGE. É exatamente essa lacuna que o `flag-scenarios-complete`
fecha — e por isso o operador minoritário virou cidadão de primeira classe aqui.

**Terceiro: `segmentMatch` e `Modulo` ficaram DE FORA, de propósito.** Os dois avaliam
contra um serviço a que o Anchors não tem (nem deve ter) acesso: um segmento vive no banco
da LaunchDarkly, não no repositório. O cenário que depende de um deles se escreve como o
VALOR que ele produz — que é sobre o que o código ramifica de qualquer forma.

## O que a implementação ensinou

**A palavra do caso ausente não podia ser cravada.** A primeira versão conhecia só
`absent`, em inglês. Um arquivo de flag escrito em português dizendo `ausente` não caía no
caso ausente: caía na divisão de operando e saía como o operador `ne` com operando `set`
— resposta ERRADA, em silêncio, justamente no cenário que mais importa. As palavras vêm do
catálogo de traduções (`flag.keyword.absent`), como o `flow.keyword.fits` já fazia, e o
teste que cobre isso morre quando a chave é trocada por uma inexistente.

**A condição que a gramática recusa vira ACHADO, não some.** O parser guarda a linha com o
erro em vez de descartá-la. Uma linha que desaparece é o silêncio que os gates existem para
acabar — e é o `flag-scenario-grammar` que a reporta, com o código e a linha.

**O eixo achou trabalho já no primeiro arquivo real.** `flags/timing-metrics.flag.md`
descreve o `--timing` construído nesta mesma sessão, e o `flag-covered` acusou os três
cenários sem teste verde.

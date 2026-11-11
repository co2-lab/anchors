---
title: O anchors.yaml
description: "A referência completa do arquivo de configuração — o que cada bloco decide, e por quê."
---

O `anchors.yaml` é onde o projeto declara **as suas próprias regras**. O Anchors
não traz um padrão embutido de como o seu código deve ser organizado — ele traz
o mecanismo, e este arquivo diz como aplicá-lo aqui.

Você não escreve este arquivo do zero: `anchors init` o gera por perguntas e
respostas, lendo as pastas do próprio projeto como camadas e a linguagem como dialeto — sem propor estrutura.

> **Chave desconhecida é ERRO, não silêncio.** O Anchors recusa carregar um
> `anchors.yaml` com uma chave que não conhece, dizendo o nome e a linha.
>
> O motivo veio de uma medição: um bloco escrito com `guide:`/`tags:` em vez de
> `from:`/`governs:` foi silenciosamente descartado, o `map build` respondeu
> "222 nós, 0 arestas" e o trabalho seguiu por um dia achando que o projeto não
> tinha relações. Um erro com a linha custa um minuto; um mapa vazio sem
> explicação custa uma sessão.

## O esqueleto

```yaml
version: 1

layers:      # QUAIS arquivos o Anchors rege, e o que cada um é
derived:     # como achar o código/feature/teste de uma spec
governs:     # quais guias regem quais arquivos
boundaries:  # quem pode importar quem
gates:       # o que é confrontado, e o que barra
workflow:    # onde o trabalho vive (local ou GitHub)
```

Só `version` e `layers` são obrigatórios. O resto entra conforme o projeto
precisa.

---

## `layers` — o que o Anchors rege

É o bloco mais importante, e o único sem o qual nada funciona. Ele responde:
**quais arquivos deste repositório são regidos, e o que cada um é.**

```yaml
layers:
    spec:
        pattern: '**/*.spec.md'
        kind: spec
        tags: [spec]
    shared:
        pattern: 'packages/shared/**/*.ts'
        kind: code
        tags: [shared, contrato]
    test:
        pattern: '**/*.test.*'
        kind: test
```

| campo | o que decide |
| --- | --- |
| `pattern` | o glob que reconhece os arquivos da camada |
| `kind` | `spec` · `feature` · `test` · `code` · `doc` · `guide` · `plan` |
| `tags` | rótulos livres, usados por `governs` para mirar grupos |
| `exclude` | globs a excluir (derivados que casam o glob amplo) |
| `regime` | `comportamental` · `declarativo` · `misto` — o que se espera daquela camada |
| `priority` | desempate declarado quando dois patterns casam o mesmo arquivo |
| `code_prefix` | prefixo de módulo no código de identidade |

**Um arquivo fora de toda camada é invisível ao Anchors.** Nenhum gate o
confronta, ele não entra no mapa, e o pipeline certifica trabalho que ninguém
verificou. Se você criar um arquivo de configuração novo (um `stryker.config.js`,
por exemplo) e o `check` reclamar que ele "não é regido", é este bloco que
precisa aprendê-lo.

### `priority`, e quando você vai precisar dele

Quando dois patterns casam o mesmo arquivo, o Anchors desempata pelo comprimento
do pattern — uma heurística que mede verbosidade, não precisão. Ela já
classificou errado: um pattern com muitas alternativas venceu outro que apontava
para um subconjunto seu.

Declare `priority` quando ela errar. O `check` avisa onde decidiu sozinho.


### `support` — arquivos de uma camada de teste que não são testes

O pattern de uma camada de teste costuma pegar arquivos que servem aos testes sem ser
testes: sub-flows que outros flows chamam (`runFlow`), agregadores de suíte, arquivos de
helper. Liste-os:

```yaml
layers:
  e2e-flow:
    pattern: "apps/mobile/.maestro/**/*.yaml"
    kind: test
    support: ["apps/mobile/.maestro/utils/**", "apps/mobile/.maestro/suites/**"]
```

Eles continuam no mapa, então mudar um helper ainda alcança todo teste que o usa. Os gates
que julgam um teste como teste ([`tests-pass`](/pt/docs/gates///tests-pass/), `test-feature-match`, [`test-traceable`](/pt/docs/gates///test-traceable/))
os pulam, e eles não contam como teste que nomeia um cenário.

---

## `derived` — como achar as peças de uma unidade

A doutrina diz que da spec nascem o código, a feature e o teste. Este bloco diz
**onde procurá-los**.

```yaml
derived:
    anchor: spec
    files:
        code: '{{dir}}/{{name}}.ts'
        feature: '{{dir}}/{{name}}.feature'
        test: '{{dir}}/{{name}}.test.ts'
```

É por convenção de nome: a spec `AreaStatus.spec.md` procura um `AreaStatus.ts`
ao lado. É o que faz a unidade ser conferida sem ninguém declarar aresta à mão.

### `overrides` — quando a convenção não serve

Arquivos de configuração quebram a convenção: uma spec chamada
`TypeScriptConfig.spec.md` governa `tsconfig.json`, `tsconfig.base.json` e os
`tsconfig.json` de cada pacote — nenhum deles se chama `TypeScriptConfig`.

```yaml
derived:
    overrides:
        - code: TSCTY
          files:
              patterns:
                  - 'tsconfig.base.json'
                  - 'tsconfig.json'
                  - 'packages/*/tsconfig.json'
```

Sem isto, a spec fica eternamente "sem código ligado" e o gate [`unit-complete`](/pt/docs/gates/unit-complete/)
reprova com razão — o arquivo existe, mas nada os liga.

---

## `aliases` — nomes curtos para caminhos longos

Um caminho escrito no `anchors.yaml` pode ser longo, e o mesmo prefixo se
repete em camadas, guias, templates de derivados e exceções de gates. O
`aliases` o nomeia uma vez, como o `paths` do tsconfig ou o
`moduleNameMapper` do Jest:

```yaml
aliases:
    "@app": apps/mobile/src
    "@guias": docs/engenharia/guias
layers:
    ui:
        pattern: "@app/screens/**/*.tsx"
        kind: code
governs:
    - from: "@guias/TELA.md"
      governs: [ui]
```

Um alias se chama `@` e uma palavra, e vale por um caminho relativo à raiz.
Ele é traduzido só como o **primeiro segmento** de um caminho, em todo campo
que guarda um — camadas, guias, `derived`, `frozen_paths`, suítes, exceções de
gates. Um caminho que cita um alias não declarado recusa a configuração,
nomeando o campo; uma pasta real com `@` no nome (`@scope/pkg`) fica como
escrita.

### Arquivos externos

Um campo que **lê** um arquivo — `governs[].from`, o `guide` de um gate de
julgamento, `changelog.template`, `packs` — pode nomear uma URL, direto ou por
um alias. Um guia compartilhado por vários repositórios mora num lugar só:

```yaml
aliases:
    "@std": https://raw.githubusercontent.com/acme/padroes/{{branch|develop}}
governs:
    - from: "@std/TELA.md"
      governs: [ui]
```

`{{branch}}` é a branch atual do projeto (o head do pull request no CI): um
projeto em `develop` lê a `develop` do repositório de padrões;
`{{branch|develop}}` lê `develop` quando a espelhada não existe. Uma URL que
não responde — nem a da branch nem a do default — faz falhar todo comando que a
lê. Cada URL fica guardada em `.anchors/cache/`, e sem rede a última cópia
responde. Um repositório privado do GitHub é lido com `GH_TOKEN` ou
`GITHUB_TOKEN`.

Um glob (`layers`, `frozen_paths`, …) ou um arquivo gravado (`derived`,
`changelog.path`, os relatórios das suítes) não pode ser URL: a configuração o
recusa.

---

## `governs` — quais guias regem quais arquivos

Um **guia** é um documento de régua transversal: o padrão de acessibilidade, a
política de segredos, a convenção de nomes. Ele não descreve uma unidade — ele
atravessa muitas.

```yaml
governs:
    - from: guides/seguranca.md
      governs: [lambdas, infra]
```

O `governs` aponta **tags** das camadas, não caminhos. É o que permite dizer "a
régua de segurança vale para toda lambda" sem listar arquivo por arquivo.

---

## `boundaries` — quem pode importar quem

```yaml
boundaries:
    - from: shared
      forbid: [lambdas, infra]
      because: 'o contrato não conhece quem o consome'
```

O `because` não é enfeite: é o texto que aparece quando o gate [`layer-boundary`](/pt/docs/gates///layer-boundary/)
reprova. Uma fronteira sem razão escrita vira "regra que alguém pôs", e a
primeira reação de quem esbarra nela é removê-la.

---

## `external_systems` — os sistemas reais com que o código fala

```yaml
external_systems:
  - name: appsync
    calls: "\\bclient\\.graphql\\("        # como uma chamada a ele é reconhecida no código
    proven_by: [integration, live]           # as camadas de suíte que alcançam o sistema real
  - name: stripe
    calls: "\\bstripe\\.[a-z]+\\.(create|update)\\("
    proven_by: [live]
    version: 2025-09-30.dahlia               # a versão da API que o projeto fala
    max_age: 90d                             # quanto tempo uma resposta gravada vale
```

Um teste de unidade prova o arquivo contra dublês, e um dublê aceita a forma que receber.
[`external-proven`](/pt/docs/gates/external-proven/) pede, a cada arquivo que chama um sistema
declarado, um teste dele que passou numa das camadas de `proven_by` (`tests[].layer`); o arquivo
dispensa com `@no-external-proof: <motivo>`. [`recording-current`](/pt/docs/gates/recording-current/)
lê a marca de cada resposta gravada, `@recorded: stripe 2025-09-30.dahlia 2026-10-02`, e reprova a
de outra versão, divergindo na mais velha que `max_age`.

---

## `required_states` — os estados que toda spec de uma camada registra

```yaml
required_states:
  screen: [loading, error, empty, offline]
```

Uma spec registra os estados em que o autor pensou.
[`states-required`](/pt/docs/gates/states-required/) pede, a cada spec da camada (o `layer:` do
header), um estado com cada nome — uma linha com o código do estado e o nome — ou
`@no-state: <nome> — <motivo>`. Os nomes são seus.

---

## `test_bypass` — o que o código pula enquanto um teste roda

```yaml
test_bypass:
  patterns:
    - "process\\.env\\.EXPO_PUBLIC_E2E === 'true'\\)\\s*return"
  proven_by: [integration]              # as camadas de suíte que rodam o caminho pulado
```

Um efeito pulado enquanto um teste roda não é provado por nenhum teste que o pula.
[`bypass-proven`](/pt/docs/gates/bypass-proven/) pede, a cada arquivo de código com um desvio, um
teste dele que passou numa camada de `proven_by`; `@no-bypass-proof: <motivo>` na linha do desvio,
ou na de cima, dispensa.

---

## `gates` — o que é confrontado

Este bloco tem [página própria](/pt/docs/gates///), porque é o mais extenso. O
formato mínimo:

```yaml
gates:
    - name: unit-complete
      on: [spec]
      check: unit-complete
      blocking: true
      measures: 'a spec tem as peças da sua unidade: código, feature e teste'
```

Um gate canônico precisa só do nome — `- name: unit-complete` — e herda o resto. Campos que uma entrada de gate também pode declarar:

| Campo | O que faz |
| --- | --- |
| `when: [manual]` | o gate roda só com `anchors check --phase manual`; um check sem fase o deixa de fora |
| `severity:` | o que cada nível de veredito faz (`fail`, `divergence`, `pending` → `block`, `inform`, `ignore`); sem ele, um gate bloqueante barra só as falhas. Um `severity:` no topo é o padrão do projeto |
| `legacy:` | a severidade dos alvos alterados pela última vez antes de uma data: `before:` a data, `severity:` o que os vereditos deles fazem (nunca mais forte que o do gate), `until:` o prazo a partir do qual o legado acaba e todo alvo responde com a severidade do gate — um gate novo barra o código novo já e o antigo a partir de uma data |
| `presupposes:` | os campos de configuração que o gate dá como declarados (`derived.mock_detect`); faltando um, o gate fica pendente e não pergunta nada, e com ele em `dialect.opt_out` é pulado |
| `review:` | os alvos do gate também ficam marcados para revisão, à parte do veredito (`review: {ask: "..."}`) — veja [julgamento e review](/pt/docs/concepts/ai-judgment/) |
| `min_coverage:` | para o `line-coverage`: o piso, 70 quando não declarado |
| `coverage_floors:` | para o `line-coverage`: um piso por glob, cada um com o seu `why` — `{"cmd/**": {min: 40, why: "..."}}` |
| `no_signal:` | alvos sem nada a medir, cada um com o motivo |

O `anchors doctor` lista todo gate do catálogo que cobre suas camadas e não está declarado, e o `anchors check --all` diz quantos numa linha.

---

## `workflow` — onde o trabalho vive

```yaml
workflow:
    mode: github
    repo: 'org/projeto'
    labels: [anchors]
    integration_branch: develop
    required_approvals: 1
```

| campo | o que decide |
| --- | --- |
| `mode` | `local` (tasks em `.anchors/`) ou `github` (cards como issues) |
| `repo` | `owner/nome` — **obrigatório** no modo github |
| `labels` | o que marca uma issue como trabalho do Anchors |
| `integration_branch` | o branch para onde o trabalho vai |
| `stale_pipeline_blocks` | um pipeline desatualizado BARRA o CI, em vez de só avisar |
| `manual_ingest_blocks` | `anchors ingest` chamado à mão é RECUSADO, em vez de só avisar |

O `repo` é obrigatório e não é inferido do remote do git de propósito: inferir
faria o Anchors escrever em outro repositório quando alguém trabalha num fork, e
escrita em lugar errado é o erro que não se desfaz com um revert.

---

## `branch_flow` — qual branch entra em qual

```yaml
branch_flow:
  merges:                      # base: as heads que um pull request pode mergear nela
    review:     [develop, "hotfix/*"]
    production: [review, "hotfix/*"]
    develop:    [review, production, "hotfix/*"]
  protected: [review, production]   # não recebem commit próprio, só merge
```

Uma branch é um glob cujo `*` não cruza `/` (`hotfix/*`), ou uma expressão regular depois de
`re:`. Uma base que nenhuma entrada nomeia recebe merge de qualquer branch, a menos que seja
protegida.

- **No remoto**, `anchors branch check` num job do evento de pull request recusa um par que o
  fluxo não permite. No evento de push, recusa um push numa branch guardada que não seja o merge
  de um pull request permitido: um repositório sem proteção no servidor ainda aceita push direto.
  A consulta ao GitHub é pelo `gh` (`GH_TOKEN`).
- **Na máquina**, os hooks pre-commit e pre-merge-commit recusam commit numa branch protegida
  (`anchors branch guard`). Merge não é confrontado ali: ele chega ao remoto como pull request.

```yaml
# .github/workflows/branch-flow.yml
on:
  pull_request:
  push:
    branches: [review, production]
jobs:
  flow:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: anchors branch check
        env:
          GH_TOKEN: ${{ github.token }}
```

---

## `pipeline` — o que cada etapa do pipeline pergunta antes de rodar

```yaml
pipeline:
  unguarded:                       # workflows que não chamam etapa, com o motivo
    build-android.yml: só gera um artefato
  stages:
    review-release:
      branches: [review]
      requires: [flow, { gates: ci }]
      verify: node scripts/smoke-review.mjs
    production-release:
      branches: [production]
      requires:
        - flow
        - any:
            - passed: review-release
            - from: "hotfix/*"
      verify: node scripts/smoke-prod.mjs
```

Um requisito é um destes:

- `flow`: o sha é o merge de um pull request que o `branch_flow` permite.
- `gates: <fase>`: os gates da fase passam.
- `passed: <etapa>`: o sha, ou então a head de onde ele foi mergeado, passou por aquela etapa.
- `from: <branch>`: o sha foi mergeado daquela head.
- `any:`: vale um de vários.

Uma etapa precisa de `verify:`, de `run:`, ou dos dois. Etapa que a configuração não declara não
roda.

```yaml
jobs:
  guard:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: anchors pipeline can-run review-release
        env: { GH_TOKEN: "${{ github.token }}" }
  deploy:
    needs: guard
    # … os passos da própria etapa
  done:
    needs: [deploy]
    permissions: { statuses: write }
    steps:
      - uses: actions/checkout@v4
      - run: anchors pipeline record review-release
        env: { GH_TOKEN: "${{ github.token }}" }
```

O `record` (apelido `done`) roda o `verify:` e registra o resultado no commit status `anchors/pipeline/<etapa>`. É esse
registro que o `passed:` lê. `anchors pipeline run <etapa>` faz can-run, `run:`, `verify:` e o
registro num passo só. O gate `pipeline-wired` lê os workflows. Ele reprova chamada a etapa não
declarada, etapa declarada que nenhum workflow chama, etapa nunca registrada, e workflow que não
chama etapa e não está em `unguarded`.

---

## `enabled` — o botão de pânico

```yaml
enabled: false
freeze_reason: 'o plano 0002 aponta para uma spec que não existe — ver #42'
```

Congela o projeto inteiro. Ver [Congelar o projeto](/pt/docs/freeze/).

**A ausência do campo significa HABILITADO.** Só o `false` explícito congela —
senão todo projeto que nunca declarou o campo nasceria parado.

---

## `dialect.tests` — como os seus testes são escritos

Os gates que leem o título de um teste ([`feature-test-match`](/pt/docs/gates///feature-test-match/), [`test-traceable`](/pt/docs/gates///test-traceable/),
[`scenario-coverage`](/pt/docs/gates///scenario-coverage/), [`flag-covered`](/pt/docs/gates///flag-covered/)) precisam saber como um teste abre na sua biblioteca de
teste. Esse conhecimento é do projeto, não do Anchors, então você o declara — com um padrão
fixo ou com um script, o que servir melhor:

```yaml
dialect:
  family: ts
  tests:
    pattern: '\b(?:it|test)(?:\.only|\.skip)?\s*\('   # a chamada que abre um teste
    # ou
    script: "node scripts/anchors-tests.mjs"            # um listador que você fornece
```

- **`pattern`** — a chamada que abre um teste, até o parêntese. O título é o literal de string
  logo depois.
- **`script`** — um comando que o Anchors roda na raiz do projeto, uma vez por varredura. Ele
  pode perguntar à própria biblioteca de teste (`jest --listTests`, `pytest --collect-only`,
  `go/ast`), então quando a biblioteca muda, a resposta muda junto. Ele imprime um objeto JSON:

  ```json
  {"version": 1, "tests": [{"file": "src/a.test.ts", "line": 12, "title": "ABCDX-B01: …"}]}
  ```

  `file` relativo à raiz, `line` a partir de 1, `title` como escrito. Qualquer outra coisa —
  outro campo, outra versão, um script que falha — é acusada pelos gates, nomeando o problema.

Declare um, nunca os dois. Sem nenhum, as famílias `go` e `ts` trazem um default
(`t.Run(`; `it`/`test`/`describe` com `.only`, `.skip`, `.each(tabela)`). Sem fonte alguma, os
gates voltam a procurar o código em qualquer ponto do arquivo de teste.

---

## `dialects` — um projeto em várias linguagens

Um repositório pode ter apps numa linguagem e serviços em outra, para que os gates continuem
confrontando a spec, o contrato e o código entre elas. Declare o léxico de cada linguagem uma
vez, por nome, em `dialects:`, e faça cada camada nomear aquele em que seus arquivos são
escritos:

```yaml
dialect:
  gherkin_language: pt        # o que não é de linguagem nenhuma fica aqui
dialects:
  ts:
    family: ts
    exported_func: '...'
    tests: { pattern: '...' }
  go:
    family: go
    exported_func: '...'
    tests: { pattern: '...' }
layers:
  screen:  { pattern: "apps/**/*.tsx",          kind: code, dialect: ts }
  handler: { pattern: "services/**/handler.go", kind: code, dialect: go }
  go-test: { pattern: "services/**/*_test.go",  kind: test, dialect: go }
```

- **Um arquivo é lido só com o dialeto da sua camada.** Um arquivo já pertence a uma camada
  (`priority` desempata dois padrões que casam), e essa camada diz em que ele é escrito.
- **Ao lado de `dialects:`, `dialect:` guarda só `gherkin_language`, `set_promise` e
  `set_slice`**: toda feature é escrita num idioma, e um nome promete o conjunto inteiro com as
  mesmas palavras em qualquer arquivo.
- **Com vários dialetos, toda camada `code` e `test` nomeia um**, ou a configuração não
  carrega: um padrão leria Go com o léxico de TS e o aprovaria sem ler. Camadas de spec, feature
  e doc não nomeiam nenhum.
- Um projeto de uma linguagem só mantém o seu bloco `dialect:`, sem nada a mudar.
- O `anchors doctor` aponta a camada com arquivos de uma família diferente da do seu dialeto,
  com o exemplo dessa família; o `anchors init` propõe um dialeto por família que encontra.

---

## `changelog` — como o `anchors changelog --write` escreve

```yaml
changelog:
  mode: incremental          # ou per_version
  path: CHANGELOG.md         # o arquivo; per_version: o diretório (padrão changelog/)
  template: doct/changelog.md.tmpl   # opcional: o seu template para uma versão
```

- **`incremental`** (o padrão) — um arquivo só; as versões que ele ainda não tem entram no
  topo, abaixo de um título `# ` inicial, cada uma depois de um marcador
  `<!-- anchors:changelog vX -->`. O que está sob um marcador é seu para editar: nunca é
  reescrito. A exceção é o bloco do que não foi lançado — ele é substituído a cada escrita,
  porque cresce até a próxima tag.
- **`per_version`** — um arquivo por versão, `<path>/<versão>.md`; o arquivo de uma versão que
  já existe é mantido, e o `unreleased.md` é reescrito.
- **`template`** — um `text/template` do Go que renderiza uma versão, como os templates da
  documentação. Recebe `.Version`, `.Date` e as listas `.Breaking`, `.Features`, `.Bugs`,
  `.Fixes`, cada entrada com `.Type`, `.Scope`, `.Subject`, `.Hash` e `.Bug`;
  `{{t "changelog.features"}}` e as outras chaves `changelog.*` traduzem os títulos para o
  `lang`.

Este é o changelog técnico; o de produto está em `anchors guide changelog`.

---

## Blocos avançados

Você provavelmente não vai precisar destes no começo.

| bloco | para quê |
| --- | --- |
| `comments` | os marcadores de comentário por linguagem, quando o projeto usa um dialeto que o Anchors não conhece |
| `rule_types` | o vocabulário de letras do código de identidade (`B` de regra, `I` de invariante…) |
| `code_lengths` | quantas letras tem um código de identidade (default: 5) |
| `obligations` | deveres regulatórios (LGPD, retenção) e quais nós estão sujeitos a eles |
| `contracts` | contratos externos cujo status precisa ser declarado |
| `regimes` | o que se espera de cada regime de camada |
| `tools` | ferramentas externas que o `anchors verify` roda por fase |

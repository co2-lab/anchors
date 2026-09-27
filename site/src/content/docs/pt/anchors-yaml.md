---
title: O anchors.yaml
description: A referência completa do arquivo de configuração — o que cada bloco decide, e por quê.
---

O `anchors.yaml` é onde o projeto declara **as suas próprias regras**. O Anchors
não traz um padrão embutido de como o seu código deve ser organizado — ele traz
o mecanismo, e este arquivo diz como aplicá-lo aqui.

Você não escreve este arquivo do zero: `anchors init` o gera por perguntas e
respostas, e sugere um preset conforme a stack que encontrar no disco.

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
que julgam um teste como teste (`tests-pass`, `test-feature-match`, `test-traceable`)
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
ao lado. É o que faz a trinca ser conferida sem ninguém declarar aresta à mão.

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

Sem isto, a spec fica eternamente "sem código ligado" e o gate `trinca-completa`
reprova com razão — o arquivo existe, mas nada os liga.

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

O `because` não é enfeite: é o texto que aparece quando o gate `layer-boundary`
reprova. Uma fronteira sem razão escrita vira "regra que alguém pôs", e a
primeira reação de quem esbarra nela é removê-la.

---

## `gates` — o que é confrontado

Este bloco tem [página própria](/docs/gates/), porque é o mais extenso. O
formato mínimo:

```yaml
gates:
    - name: trinca-completa
      on: [spec]
      check: trinca-completa
      blocking: true
      measures: 'a spec tem código, feature e teste que a realizam'
```

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

## `enabled` — o botão de pânico

```yaml
enabled: false
freeze_reason: 'o plano 0002 aponta para uma spec que não existe — ver #42'
```

Congela o projeto inteiro. Ver [Congelar o projeto](/docs/congelar/).

**A ausência do campo significa HABILITADO.** Só o `false` explícito congela —
senão todo projeto que nunca declarou o campo nasceria parado.

---

## `dialect.tests` — como os seus testes são escritos

Os gates que leem o título de um teste (`feature-test-match`, `test-traceable`,
`scenario-coverage`, `flag-covered`) precisam saber como um teste abre na sua biblioteca de
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

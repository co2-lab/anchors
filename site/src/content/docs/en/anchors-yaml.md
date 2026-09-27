---
title: The anchors.yaml
description: The complete configuration reference — what each block decides, and why.
---

`anchors.yaml` is where a project declares **its own rules**. Anchors ships no
built-in standard for how your code should be organized — it ships the
mechanism, and this file says how to apply it here.

You don't write this file from scratch: `anchors init` generates it through
questions and answers, suggesting a preset based on the stack it finds on disk.

> **An unknown key is an ERROR, not silence.** Anchors refuses to load an
> `anchors.yaml` with a key it doesn't know, naming the key and the line.
>
> The reason came from a measurement: a block written with `guide:`/`tags:`
> instead of `from:`/`governs:` was silently discarded, `map build` answered
> "222 nodes, 0 edges", and work went on for a day believing the project had no
> relations. An error with a line number costs a minute; an empty map with no
> explanation costs a session.

## The skeleton

```yaml
version: 1

layers:      # WHICH files Anchors governs, and what each one is
derived:     # how to find a spec's code/feature/test
governs:     # which guides govern which files
boundaries:  # who may import whom
gates:       # what gets confronted, and what blocks
workflow:    # where the work lives (local or GitHub)
```

Only `version` and `layers` are required. The rest comes in as the project needs
it.

---

## `layers` — what Anchors governs

The most important block, and the only one without which nothing works. It
answers: **which files in this repository are governed, and what is each one.**

```yaml
layers:
    spec:
        pattern: '**/*.spec.md'
        kind: spec
        tags: [spec]
    shared:
        pattern: 'packages/shared/**/*.ts'
        kind: code
        tags: [shared, contract]
    test:
        pattern: '**/*.test.*'
        kind: test
```

| field | decides |
| --- | --- |
| `pattern` | the glob that recognizes the layer's files |
| `kind` | `spec` · `feature` · `test` · `code` · `doc` · `guide` · `plan` |
| `tags` | free labels, used by `governs` to target groups |
| `exclude` | globs to exclude (derived files matching a broad glob) |
| `regime` | `comportamental` · `declarativo` · `misto` — what that layer is expected to be |
| `priority` | declared tie-break when two patterns match the same file |
| `code_prefix` | module prefix in the identity code |

**A file outside every layer is invisible to Anchors.** No gate confronts it, it
doesn't enter the map, and the pipeline certifies work nobody verified. If you
create a new config file (a `stryker.config.js`, say) and `check` complains it
"isn't governed", this is the block that needs to learn about it.

### `priority`, and when you'll need it

When two patterns match the same file, Anchors breaks the tie by pattern length —
a heuristic that measures verbosity, not precision. It has already classified
wrongly: a pattern with many alternatives beat one pointing at a subset of it.

Declare `priority` when it gets it wrong. `check` tells you where it decided on
its own.


### `support` — files of a test layer that are not tests

A test layer's pattern often catches files that serve the tests without being tests:
sub-flows other flows call (`runFlow`), suite aggregators, helper files. List them:

```yaml
layers:
  e2e-flow:
    pattern: "apps/mobile/.maestro/**/*.yaml"
    kind: test
    support: ["apps/mobile/.maestro/utils/**", "apps/mobile/.maestro/suites/**"]
```

They stay in the map, so a change to a helper still reaches every test that uses it. The
gates that judge a test as a test (`tests-pass`, `test-feature-match`, `test-traceable`)
skip them, and they don't count as tests that name a scenario.

---

## `derived` — how to find a unit's pieces

The doctrine says code, feature and test are born from the spec. This block says
**where to look for them**.

```yaml
derived:
    anchor: spec
    files:
        code: '{{dir}}/{{name}}.ts'
        feature: '{{dir}}/{{name}}.feature'
        test: '{{dir}}/{{name}}.test.ts'
```

It's by naming convention: the spec `AreaStatus.spec.md` looks for an
`AreaStatus.ts` beside it. That's what lets the triad be checked without anyone
declaring edges by hand.

### `overrides` — when the convention doesn't fit

Config files break the convention: a spec named `TypeScriptConfig.spec.md`
governs `tsconfig.json`, `tsconfig.base.json` and each package's
`tsconfig.json` — none of which is named `TypeScriptConfig`.

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

Without this, the spec stays forever "with no code linked" and the
`trinca-completa` gate rightly fails — the file exists, but nothing connects
them.

---

## `governs` — which guides govern which files

A **guide** is a cross-cutting rule document: the accessibility standard, the
secrets policy, the naming convention. It doesn't describe a unit — it cuts
across many.

```yaml
governs:
    - from: guides/security.md
      governs: [lambdas, infra]
```

`governs` points at layer **tags**, not paths. That's what lets you say "the
security rule applies to every lambda" without listing files one by one.

---

## `boundaries` — who may import whom

```yaml
boundaries:
    - from: shared
      forbid: [lambdas, infra]
      because: 'the contract does not know its consumers'
```

`because` isn't decoration: it's the text shown when the `layer-boundary` gate
fails. A boundary with no written reason becomes "a rule someone added", and the
first reaction of whoever hits it is to remove it.

---

## `gates` — what gets confronted

This block has [its own page](/en/docs/gates/), being the longest. The minimal
form:

```yaml
gates:
    - name: trinca-completa
      on: [spec]
      check: trinca-completa
      blocking: true
      measures: 'the spec has code, feature and test fulfilling it'
```

---

## `workflow` — where the work lives

```yaml
workflow:
    mode: github
    repo: 'org/project'
    labels: [anchors]
    integration_branch: develop
    required_approvals: 1
```

| field | decides |
| --- | --- |
| `mode` | `local` (tasks in `.anchors/`) or `github` (cards as issues) |
| `repo` | `owner/name` — **required** in github mode |
| `labels` | what marks an issue as Anchors work |
| `integration_branch` | the branch work goes to |
| `stale_pipeline_blocks` | an outdated pipeline BLOCKS CI instead of merely warning |
| `manual_ingest_blocks` | `anchors ingest` called by hand is REFUSED instead of merely warning |

`repo` is required and deliberately not inferred from the git remote: inferring
would make Anchors write to another repository when someone works on a fork, and
writing to the wrong place is the mistake a revert doesn't undo.

---

## `enabled` — the panic button

```yaml
enabled: false
freeze_reason: 'plan 0002 points at a spec that does not exist — see #42'
```

Freezes the whole project. See [Freezing the project](/en/docs/congelar/).

**An absent field means ENABLED.** Only an explicit `false` freezes — otherwise
every project that never declared the field would be born stopped.

---

## `dialect.tests` — how your tests are written

The gates that read a test's title (`feature-test-match`, `test-traceable`,
`scenario-coverage`, `flag-covered`) need to know how a test opens in your test library.
That is your project's knowledge, not Anchors', so you declare it — with a fixed pattern or
with a script, whichever fits:

```yaml
dialect:
  family: ts
  tests:
    pattern: '\b(?:it|test)(?:\.only|\.skip)?\s*\('   # the call that opens a test
    # or
    script: "node scripts/anchors-tests.mjs"            # a lister you provide
```

- **`pattern`** — the call that opens a test, up to its parenthesis. The title is the string
  literal right after it.
- **`script`** — a command Anchors runs at the project root, once per scan. It can ask your
  test library itself (`jest --listTests`, `pytest --collect-only`, `go/ast`), so when the
  library changes, the answer changes with it. It prints one JSON object:

  ```json
  {"version": 1, "tests": [{"file": "src/a.test.ts", "line": 12, "title": "ABCDX-B01: …"}]}
  ```

  `file` relative to the root, `line` from 1, `title` as written. Anything else — another
  field, another version, a script that fails — is reported by the gates, naming the problem.

Declare one, never both. Without either, the `go` and `ts` families bring a default
(`t.Run(`; `it`/`test`/`describe` with `.only`, `.skip`, `.each(table)`). With no source at
all, the gates fall back to finding the code anywhere in the test file.

---

## Advanced blocks

You probably won't need these early on.

| block | for |
| --- | --- |
| `comments` | comment markers per language, when the project uses a dialect Anchors doesn't know |
| `rule_types` | the identity-code letter vocabulary (`B` for rule, `I` for invariant…) |
| `code_lengths` | how many letters an identity code has (default: 5) |
| `obligations` | regulatory duties (GDPR, retention) and which nodes are subject to them |
| `contracts` | external contracts whose status must be declared |
| `regimes` | what each layer regime is expected to be |
| `tools` | external tools `anchors verify` runs per phase |

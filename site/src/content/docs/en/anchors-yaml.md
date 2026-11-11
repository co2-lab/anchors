---
title: The anchors.yaml
description: "The complete configuration reference — what each block decides, and why."
---

`anchors.yaml` is where a project declares **its own rules**. Anchors ships no
built-in standard for how your code should be organized — it ships the
mechanism, and this file says how to apply it here.

You don't write this file from scratch: `anchors init` generates it through
questions and answers, reading the project's own folders as layers and its language as dialect — it proposes no structure.

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
`AreaStatus.ts` beside it. That's what lets the unit be checked without anyone
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
`unit-complete` gate rightly fails — the file exists, but nothing connects
them.

---

## `aliases` — short names for long paths

A path written in `anchors.yaml` can be long, and the same prefix repeats in
layers, guides, derived templates and gate exceptions. `aliases` names it once,
as a tsconfig's `paths` or Jest's `moduleNameMapper` do:

```yaml
aliases:
    "@app": apps/mobile/src
    "@guides": docs/engineering/guides
layers:
    ui:
        pattern: "@app/screens/**/*.tsx"
        kind: code
governs:
    - from: "@guides/SCREEN_GUIDE.md"
      governs: [ui]
```

An alias is named `@` and a word, and stands for a path relative to the root.
It is translated only as the **first segment** of a path, in every field that
holds one — layers, guides, `derived`, `frozen_paths`, suites, gate exceptions.
A path citing an alias nobody declared refuses the configuration, naming the
field; a real folder named with `@` (`@scope/pkg`) is kept as written.

### External files

A field that **reads** a file — `governs[].from`, a judgment gate's `guide`,
`changelog.template`, `packs` — may name a URL, directly or through an alias.
A guide shared by several repositories lives in one place:

```yaml
aliases:
    "@std": https://raw.githubusercontent.com/acme/standards/{{branch|develop}}
governs:
    - from: "@std/SCREEN_GUIDE.md"
      governs: [ui]
```

`{{branch}}` is the project's current branch (a pull request's head in CI), so
a project on `develop` reads the `develop` of the standards repository;
`{{branch|develop}}` reads `develop` when the mirror does not exist. A URL that
answers nothing — neither the branch nor the default — fails every command
that reads it. Each URL is kept in `.anchors/cache/`, and offline the last copy
answers. A private GitHub repository is read with `GH_TOKEN` or `GITHUB_TOKEN`.

A glob (`layers`, `frozen_paths`, …) or a file written (`derived`,
`changelog.path`, the suites' reports) cannot be a URL: the configuration
refuses it.

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
    - layer: shared
      forbid: "from ['\"](\\.\\./)*(lambdas|infra)/"
      because: 'the contract does not know its consumers'
    - layer: screen
      forbid: "<Text[^>]*>\\{[^}]*\\.(kind|category|action|status)\\}"
      because: 'a value the backend writes is a key, not a label: map it to its text'
```

Each boundary names a layer (empty: all code), a pattern its files may not contain — an
import, or any other construct, such as a backend's enum shown raw on a screen —, and why.
`because` isn't decoration: it's the text shown when the `layer-boundary` gate
fails. A boundary with no written reason becomes "a rule someone added", and the
first reaction of whoever hits it is to remove it.

---

## `external_systems` — the real systems your code talks to

```yaml
external_systems:
  - name: appsync
    calls: "\\bclient\\.graphql\\("        # how a call to it is recognised in the code
    proven_by: [integration, live]           # the suite layers that reach the real system
  - name: stripe
    calls: "\\bstripe\\.[a-z]+\\.(create|update)\\("
    proven_by: [live]
    version: 2025-09-30.dahlia               # the API version the project speaks
    max_age: 90d                             # how long a recorded response lives
```

A unit test proves a file against doubles, and a double accepts whatever shape it is given.
[`external-proven`](/docs/gates/external-proven/) asks each file that calls a declared system for
a test of its own that passed in one of its `proven_by` layers (`tests[].layer`); a file waives it
with `@no-external-proof: <reason>`. [`recording-current`](/docs/gates/recording-current/) reads
each recorded response's mark, `@recorded: stripe 2025-09-30.dahlia 2026-10-02`, and fails one of
another version, diverging on one older than `max_age`.

---

## `required_states` — the states every spec of a layer registers

```yaml
required_states:
  screen: [loading, error, empty, offline]
```

A spec registers the states its author thought of.
[`states-required`](/docs/gates/states-required/) asks each spec of the layer (its header's
`layer:`) for a state by each name — a line with the state's code and the name — or
`@no-state: <name> — <reason>`. The names are yours.

---

## `test_bypass` — what the code skips while a test runs

```yaml
test_bypass:
  patterns:
    - "process\\.env\\.EXPO_PUBLIC_E2E === 'true'\\)\\s*return"
  proven_by: [integration]              # the suite layers that run the skipped path
```

An effect skipped while a test runs is proven by no test that skips it.
[`bypass-proven`](/docs/gates/bypass-proven/) asks each code file with a bypass for a test of its
own that passed in a `proven_by` layer; `@no-bypass-proof: <reason>` on the bypass's line, or the
line above, waives it.

---

## `gates` — what gets confronted

This block has [its own page](/docs/gates/), being the longest. The minimal
form:

```yaml
gates:
    - name: unit-complete
      on: [spec]
      check: unit-complete
      blocking: true
      measures: 'the spec has the pieces of its unit: code, feature and test'
```

A canonical gate needs only its name — `- name: unit-complete` — and inherits the rest. Fields a gate entry can also declare:

| Field | What it does |
| --- | --- |
| `when: [manual]` | the gate runs only under `anchors check --phase manual`; a check with no phase leaves it out |
| `severity:` | what each verdict level does (`fail`, `divergence`, `pending` → `block`, `inform`, `ignore`); with none, a blocking gate blocks only its failures. A top-level `severity:` is the project's default |
| `legacy:` | the severity of the targets last changed before a date: `before:` the date, `severity:` what their verdicts do (never stronger than the gate's), `until:` the deadline from which the legacy ends and every target answers with the gate's own — a new gate can bar the new code at once and the old code from a date on |
| `presupposes:` | the configuration fields the gate takes as declared (`derived.mock_detect`); with one missing the gate is pending and asks nothing, and with it in `dialect.opt_out` it skips |
| `review:` | the gate's targets are also marked to review, apart from its verdict (`review: {ask: "..."}`) — see [judgment and review](/docs/concepts/ai-judgment/) |
| `min_coverage:` | for `line-coverage`: the floor, 70 when undeclared |
| `coverage_floors:` | for `line-coverage`: a floor per glob, each with its `why` — `{"cmd/**": {min: 40, why: "..."}}` |
| `no_signal:` | targets with nothing to measure, each with its reason |

`anchors doctor` names every catalog gate that covers your layers and is not declared, and `anchors check --all` says how many in one line.

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

## `branch_flow` — which branch merges into which

```yaml
branch_flow:
  merges:                      # base: the heads a pull request may merge into it
    review:     [develop, "hotfix/*"]
    production: [review, "hotfix/*"]
    develop:    [review, production, "hotfix/*"]
  protected: [review, production]   # take no commit of their own, only merges
```

A branch is a glob whose `*` does not cross `/` (`hotfix/*`), or a regular expression after
`re:`. A base no entry names takes a merge from anywhere, unless it is protected.

- **On the remote**, `anchors branch check` in a job on the pull request event refuses a pair the
  flow does not allow. On the push event it refuses a push to a guarded branch that is not the
  merge of an allowed pull request: a repository with no server-side protection still takes a
  direct push. It asks GitHub through the `gh` CLI (`GH_TOKEN`).
- **On the machine**, the pre-commit and pre-merge-commit hooks refuse a commit on a protected
  branch (`anchors branch guard`). Merges are not confronted there: they reach the remote as pull
  requests.

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

## `pipeline` — what each stage of the pipeline asks before it runs

```yaml
pipeline:
  unguarded:                       # workflows that call no stage, with the reason
    build-android.yml: only builds an artifact
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

A requirement is one of these:

- `flow`: the sha is the merge of a pull request `branch_flow` allows.
- `gates: <phase>`: the phase's gates pass.
- `passed: <stage>`: the sha, or else the head it was merged from, passed that stage.
- `from: <branch>`: the sha was merged from that head.
- `any:`: one of several holds.

A stage needs `verify:`, `run:`, or both. A stage the configuration does not declare does not run.

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
    # … the stage's own steps
  done:
    needs: [deploy]
    permissions: { statuses: write }
    steps:
      - uses: actions/checkout@v4
      - run: anchors pipeline record review-release
        env: { GH_TOKEN: "${{ github.token }}" }
```

`record` (alias `done`) runs `verify:` and records the outcome as the commit status `anchors/pipeline/<stage>`.
That record is what `passed:` reads. `anchors pipeline run <stage>` does can-run, `run:`,
`verify:` and the record in one step. The `pipeline-wired` gate reads the workflows. It fails on a
call to an undeclared stage, a declared stage no workflow calls, a stage never recorded, and a
workflow that calls no stage and is not unguarded.

---

## `enabled` — the panic button

```yaml
enabled: false
freeze_reason: 'plan 0002 points at a spec that does not exist — see #42'
```

Freezes the whole project. See [Freezing the project](/docs/freeze/).

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

## `dialects` — a project in several languages

A repository can hold apps in one language and services in another, so the gates keep
confronting the spec, the contract and the code across them. Declare each language's lexicon
once, by name, under `dialects:`, and have each layer name the one its files are written in:

```yaml
dialect:
  gherkin_language: en        # what is about no language stays here
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

- **A file is read with its layer's dialect alone.** A file already belongs to one layer
  (`priority` settles two matching patterns), and that layer says what it is written in.
- **Beside `dialects:`, `dialect:` holds only `gherkin_language`, `set_promise` and
  `set_slice`**: every feature is written in one language, and a name promises the whole set
  in the same words in any file.
- **With several dialects, every `code` and `test` layer names one**, or the configuration
  does not load: a default would read Go with the TS lexicon and pass it green. Spec, feature
  and doc layers name none.
- A project of one language keeps its one `dialect:` block, with nothing to change.
- `anchors doctor` names a layer holding files of a family its dialect's is not, with that
  family's example; `anchors init` proposes one dialect per family it finds.

---

## `changelog` — how `anchors changelog --write` writes

```yaml
changelog:
  mode: incremental          # or per_version
  path: CHANGELOG.md         # the file; per_version: the directory (default changelog/)
  template: doct/changelog.md.tmpl   # optional: your own template for one release
```

- **`incremental`** (the default) — one file; the releases it does not hold yet go on its top,
  below a leading `# ` title, each after a `<!-- anchors:changelog vX -->` marker. What is
  under a marker is yours to edit: it is never rewritten. The unreleased block is the one
  exception — it is replaced on each write, since it grows until the next tag.
- **`per_version`** — one file per release, `<path>/<version>.md`; a release's file that
  exists is kept, and `unreleased.md` is rewritten.
- **`template`** — a Go `text/template` rendering one release, like the docs templates. It
  receives `.Version`, `.Date` and the lists `.Breaking`, `.Features`, `.Bugs`, `.Fixes`, each
  entry with `.Type`, `.Scope`, `.Subject`, `.Hash` and `.Bug`; `{{t "changelog.features"}}`
  and the other `changelog.*` keys translate the headings to `lang`.

This is the technical changelog; see `anchors guide changelog` for the product one.

---

## `push_review` — the evidence each push carries

```yaml
push_review:
  required: true          # each pushed commit is covered by a review whose verdict is ok
  file: REVIEWS.md        # where the reviews are written (the default)
  retention: 30d          # writing a review drops those older; empty keeps them all
  require_judgment: true  # each judgment gate the pushed commits reach has a verdict
  require_gates: true     # the blocking gates pass over what the push sends, no waiver
```

Off unless declared. The pre-push hook runs `anchors push-review verify` with the refs the push
sends (`anchors install-hooks --force` brings the hook up to date):

- **`required`** — whoever pushes reviews the commits it sends (`anchors guide review` is the
  ruler) and records the verdict: `anchors push-review --ok "<what was validated, and how>" --by
  <who>`, or `--fail "<what failed>"`, with `--note` for each observation. The review goes to the
  top of `file` and is committed alone. A commit no review covers stops the push, and so does one
  whose newest review failed. A review covers a commit by its hash, or by its patch id after a
  rebase; a commit that changes only the reviews' file and the map owes no review.
- **`require_judgment`** — each judgment gate the pushed commits reach has a verdict
  (`anchors judge`) in the map the push sends; an ok review is not recorded over one still
  awaiting. It works with or without `required`.
- **`require_gates`** — the blocking gates a commit or a push charges run over what the pushed
  commits reach, each file read as the pushed tip holds it, and no waiver a commit message
  declared (`[skip-<gate>@<CODE>: why]`) counts: it frees the commit in the middle of the work,
  and is settled before the push. A failure refuses the push naming the gate, the target and the
  commit that waived it. A waiver meant to reach the remote goes in the artifact (`@no-*`).
- **`retention`** — days a review stays in the file; older ones are dropped when one is written.

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

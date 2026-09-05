---
title: The gates
description: The complete catalog, how to configure them, and what to enable per project type.
---

A **gate** is a question the project asks itself, whose answer Anchors records.
[Quality](/en/docs/qualidade/) explains the doctrine; this page is the practical
reference: which gates exist, what each measures, and which to enable.

## How a gate is declared

```yaml
gates:
    - name: trinca-completa      # how it appears in output
      on: [spec]                 # which file kind it runs on
      check: trinca-completa     # the internal checker fulfilling it
      blocking: true             # does failing BLOCK the commit/CI?
      measures: 'the spec has code, feature and test fulfilling it'
```

| field | decides |
| --- | --- |
| `on` | the file `kind`s the gate runs on |
| `check` | the **internal** (deterministic) checker |
| `run` | an **external** command, when the gate is a project tool |
| `blocking` | `true` blocks; `false` informs |
| `measures` | what it measures, in one sentence — shown in output |
| `requires` | runs only on targets whose content contains this text |
| `when` | the phases it runs in: `pre-commit`, `pre-push`, `ci` |
| `ask` | the question of an **AI judgment** gate |

### The three outcomes

| symbol | means |
| --- | --- |
| `✓` | passed |
| `✗` | failed — becomes an issue; blocks if `blocking: true` |
| `~` | **indeterminate**: the gate had nothing to confront |

`~` isn't failure. A coverage gate on a file with no test doesn't fail — it
doesn't measure. Confusing the two makes someone "fix" what isn't broken.

---

## Start informative

**Every new gate should be born with `blocking: false`.**

Not timidity: it's the only way to learn what it will flag before it blocks
someone's work. A gate born blocking in a project that already carries debt
fails everything on day one, and the cheap way out becomes turning it off —
which is worse than never enabling it.

The path is: enable informative → measure for a week → fix what it found →
promote to blocking.

---

## The catalog

### The triad — does the spec have the pieces fulfilling it?

| gate | measures |
| --- | --- |
| `trinca-completa` | the spec has code, feature and test linked |
| `spec-tem-codigo` | the spec carries an identity code |
| `spec-completa` | the spec has at least one rule, with no placeholder |
| `spec-feature-match` | each spec rule has a scenario in the feature |
| `feature-nao-vazia` | the feature has a real scenario |

`trinca-completa` prevents Anchors' quietest defect: a lone spec passes **every**
relational gate — they fail *open*, with no test linked there's nothing to
confront — and the pipeline concludes "ready to promote" over work that doesn't
exist.

### Identity — can the rule be found?

| gate | measures |
| --- | --- |
| `codigo-catalogado` | every exported symbol has a rule in the spec, or a written waiver |
| `code-reference-valid` | code references point at rules that exist |
| `rule-types` | code letters belong to the declared vocabulary |
| `teste-rastreavel` | the test cites the code of the scenario it proves |
| `scenario-asserts` | the scenario asserts something instead of merely running |

### Planning — does the plan still describe reality?

| gate | measures |
| --- | --- |
| `fase-existe` | the phases cited exist in the plan |
| `fase-ordenada` | the declared order between phases is coherent |
| `plan-seeds-valid` | the specs the plan seeds exist |
| `plano-alterado-justificado` | a plan/spec that CHANGED says why |
| `plano-revisado` | the revision is numbered and explained |
| `parent-valido` | `parent:` points at a phase that exists |
| `open-questions-resolved` | open decisions have been decided |

`plano-alterado-justificado` deserves a note. It looks at the **diff**, not the
content, because drift is silent by construction: a plan quietly corrected is
perfectly valid — the inconsistency was removed. What exposes it isn't the file's
state, it's the change without justification.

### Proof — does the test prove, or merely execute?

| gate | measures |
| --- | --- |
| `tests-green` | the suite passes |
| `line-coverage` | the line executed during the test |
| `coverage-delta` | coverage didn't drop with this change |
| `mutation-score` | if the line changed, would a test break |
| `scenario-coverage` | each spec scenario has a green test |

**Coverage and mutation are not the same thing**, and confusing them is the
defect `mutation-score` exists to catch. A recently measured project had **100%
line coverage** and 47 surviving mutants — 47 code changes no test noticed.

### Boundaries — who knows whom

| gate | measures |
| --- | --- |
| `layer-boundary` | nobody imports what the layer forbids |
| `sibling-guard` | a module doesn't reach a sibling the wrong way |
| `prova-cruza-fronteira` | a test crossing a boundary declares it |
| `dependency-honored` | the declared dependency is the one that exists |

### Test doubles — does the mock tell the truth?

| gate | measures |
| --- | --- |
| `mock-carimbado` | every double declares what it pretends to be |
| `mock-tipado` | the double honors the contract of what it replaces |
| `mock-detect-cobre-o-dialeto` | the double-detection regex reaches the forms the project uses |

### Security and compliance

| gate | measures |
| --- | --- |
| `secret-nao-vazado` | no secret enters the repository |
| `dependencia-vulneravel` | how many known CVEs the dependencies carry |
| `sbom-gerado` | the dependency inventory is published |
| `obligation-honored` | the declared regulatory duties are honored |
| `contract-status-declared` | each external contract's status is stated |

### AI judgment

Three gates have no deterministic answer — they **ask**:

| gate | question |
| --- | --- |
| `regra-cumprida` | does the marked passage FULFILL what the rule describes? |
| `no-test-prova-real` | does the proof pointed at by `@no-test` really exercise the behavior? |
| `mock-detect-cobre-o-dialeto` | does the declared pattern reach every double? |

The verdict is recorded with `anchors judge`, and has three values:

```sh
anchors judge <target> --gate <g> --verdict pass|fail|dispensado --reason "..."
```

**`dispensado`** exists because the other two lie when the target doesn't exist.
A spec declaring `@TBD: code` states the code hasn't been written yet — so `pass`
would assert it fulfills the rule, and `fail` would fail work nobody got wrong.

---

## What to enable per project type

There's no universal set. What follows are starting points measured in real
projects.

### New project, starting from the spec

Enable **everything as informative** and promote nothing in week one. You need to
see what the project has before deciding what to block.

The first to promote, once the project has 3–4 specs:

```yaml
- name: spec-tem-codigo      # without identity nothing is traceable
- name: spec-completa        # a spec with a placeholder decides nothing
- name: fase-existe          # the plan points at a phase that exists
```

### Backend / serverless

What matters is **boundary** and **secret** — the two things that break in
production and don't show up in tests.

```yaml
- name: layer-boundary            blocking: true   # the contract doesn't know its consumers
- name: secret-nao-vazado         blocking: true   when: [pre-commit, pre-push, ci]
- name: licenca-compativel        blocking: true   when: [pre-push, ci]
- name: trinca-completa           blocking: true
- name: codigo-catalogado         blocking: true
- name: mutation-score            blocking: false  # informative until you know the project's number
- name: dependencia-vulneravel    blocking: false
```

`mutation-score` stays informative **on purpose**: demanding 80% before knowing
what the project has today would pick a number in the dark.

### Application with a UI

Add the ones confronting what the screen promises:

```yaml
- name: contract-status-declared  blocking: true   # the screen declares each datum's state
- name: domain-declared           blocking: true   # a value arriving wrong has handling
- name: pagination-honored        blocking: false  # a list that paginates really paginates
- name: count-honored             blocking: false  # the displayed count is the real one
```

### Project with regulated data (GDPR, health, finance)

```yaml
- name: obligation-honored        blocking: true
- name: marker-parity             blocking: true   # the same rule appears on BOTH ends
- name: secret-nao-vazado         blocking: true
```

`marker-parity` prevents the costliest divergence of that kind of project: the
screen promises to delete one datum, the backend deletes another, and nothing
flags it — because each side is internally coherent.

### Library / CLI

Boundaries matter less; **proof** matters more.

```yaml
- name: tests-green               blocking: true
- name: mutation-score            blocking: false → true once it stabilizes
- name: codigo-catalogado         blocking: true
- name: sem-duplicacao            blocking: false
```

---

## External gates: when the tool is the project's

Not every gate is internal. A project command becomes a gate via `run`:

```yaml
- name: licenca-compativel
  on: [code]
  scope: project
  run: 'bash scripts/anchors-licencas.sh'
  needs_tool: go-licenses
  install_hint: 'go install github.com/google/go-licenses@latest'
  blocking: true
  when: [pre-push, ci]
```

`needs_tool` and `install_hint` matter: without them, a gate depending on a
missing tool fails with a shell error, and whoever sees it can't tell whether the
code is wrong or something needs installing.

> **A declared gate that never runs is worse than no gate**, because it sits in
> the configuration as if it protected something. If the gate declares
> `when: [ci]` and your CI doesn't execute it, it's decoration. Check with
> `anchors doctor`.

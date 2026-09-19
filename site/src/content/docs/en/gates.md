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
| `region-pair-honored` | `#region` and `#endregion` markers close with matching code |

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
| `evidence-fresh` | the test evidence remains fresh (code/closure has not advanced) |
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

### Security, compliance, and external governance

| gate | type | measures |
| --- | --- | --- |
| `no-secret-leaked` | external (`run`) | no secret enters the repository |
| `dependency-vulnerable` | external (`run`) | how many known CVEs the dependencies carry |
| `sbom-generated` | external (`run`) | the CycloneDX/SPDX dependency inventory is published |
| `license-compatible` | external (`run`) | absence of dependencies with strong copyleft or incompatible licenses |
| `no-duplication` | external (`run`) | no code block appears duplicated across files |
| `spellcheck` | external (`run`) | no spelling mistakes in text, code, and identifiers |
| `circular` | external (`run`) | no circular dependencies between modules |
| `deadcode` | external (`run`) | no dead exports, functions, or orphan files |
| `obligation-honored` | internal | the declared regulatory duties are honored |
| `contract-status-declared` | internal | each external contract's status is stated |

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

## External gates: when the tool belongs to the project

Not every gate is built into the Anchors binary. Anchors is **strictly language-agnostic**: it decides **when** to check (`when`), the **scope** (`batch` vs `project`), and whether a failure **blocks** the commit or PR (`blocking`), while your project provides the appropriate tool via `run:`.

```yaml
- name: license-compatible
  on: [code]
  scope: project
  run: 'go-licenses check ./... --disallowed_types=forbidden,restricted'
  needs_tool: go-licenses
  install_hint: 'go install github.com/google/go-licenses@latest'
  blocking: false
  when: [ci]
```

`needs_tool` and `install_hint` are critical: if the tool is absent from the host, Anchors emits **`Skip`** and `anchors doctor` warns what needs installing, rather than breaking the build with a cryptic shell error.

---

### Tooling guide across ecosystems

The table below summarizes recommended tools for each external gate across major platforms:

| Gate | Universal (Native Binary) | Go | Node.js / TypeScript | Python | Rust |
|---|---|---|---|---|---|
| `no-secret-leaked` | `gitleaks` | `gitleaks` | `gitleaks` | `detect-secrets` | `gitleaks` |
| `dependency-vulnerable` | `osv-scanner` / `trivy` | `govulncheck` | `pnpm audit` / `osv-scanner` | `pip-audit` | `cargo-audit` |
| `sbom-generated` | `syft` | `syft` / `cyclonedx-gomod` | `syft` / `@cyclonedx/cyclonedx-npm` | `cyclonedx-py` | `cargo-cyclonedx` |
| `no-duplication` | `pmd cpd` | `dupl` | `jscpd` | `pylint --enable=similarities` | `flcl` / `pmd cpd` |
| `spellcheck` | `typos` | `typos` | `typos` or `cspell` | `typos` or `codespell` | `typos` |
| `license-compatible` | — | `go-licenses` | `license-checker` | `pip-licenses` | `cargo-deny` |
| `circular` | — | Go compiler / `go vet` | `madge` | `import-linter` | Rust compiler |
| `deadcode` | — | `deadcode` (x/tools) | `knip` | `vulture` | `cargo-udeps` |

#### 1. `no-secret-leaked` (Secrets detection)
Blocking from day one — a committed credential remains in git history indefinitely.
- **Universal (Recommended)**: `gitleaks git --no-banner --redact -v` (`brew install gitleaks`)

#### 2. `dependency-vulnerable` (Vulnerabilities in dependencies)
Audits lockfiles against known CVE databases.
- **Universal (Recommended)**: `osv-scanner scan source -r .` (`brew install osv-scanner`)
- **Go**: `govulncheck ./...`
- **Node / TS**: `pnpm audit --prod` or `npm audit --omit=dev`
- **Python**: `pip-audit`
- **Rust**: `cargo-audit`

#### 3. `sbom-generated` (Software Bill of Materials)
Generates CycloneDX/SPDX inventories for audits and compliance.
- **Universal (Recommended)**: `syft scan dir:. -o cyclonedx-json=sbom.json -q` (`brew install syft`)
- **Go**: `cyclonedx-gomod app -json -output sbom.json`
- **Node / TS**: `npx @cyclonedx/cyclonedx-npm --output-file sbom.json`

#### 4. `no-duplication` (Copy-paste detection)
Flags identical blocks of logic duplicated across multiple files.
- **Universal**: `pmd cpd --minimum-tokens 70 --dir . --language <lang>` (`brew install pmd`)
- **Node / TS / Multi-lang**: `npx --yes jscpd . --reporters console --silent`
- **Go**: `dupl -t 70`

#### 5. `spellcheck` (Spelling mistakes)
Eliminates spelling errors in text, error messages, and code identifiers (camelCase, snake_case).
- **Universal (Recommended)**: `typos` (`brew install typos`) — native, ultra-fast binary in Rust with zero runtime dependencies.
- **Node / TS**: `npx cspell --no-progress --no-summary {{files}}`

#### 6. `license-compatible` (License compliance)
Prevents accidental introduction of strong copyleft (AGPL/SSPL) or restricted dependencies into proprietary code.
- **Go**: `go-licenses check ./... --disallowed_types=forbidden,restricted`
- **Node / TS**: `npx license-checker --production --onlyAllow 'MIT;Apache-2.0;BSD-2-Clause;BSD-3-Clause;ISC'`
- **Python**: `pip-licenses`
- **Rust**: `cargo-deny check bans licenses`

#### 7. `circular` (Circular dependencies)
Detects cycles in the module or package import graph.
- **Node / TS**: `npx madge --circular --extensions ts,tsx src/`
- **Python**: `lint-imports` (via `import-linter`)
- **Go / Rust**: Enforced natively by the language compiler.

#### 8. `deadcode` (Dead code detection)
Identifies unconsumed exports, functions, types, and orphan files.
- **Go**: `deadcode ./...` (`go install golang.org/x/tools/cmd/deadcode@latest`)
- **Node / TS**: `npx knip`
- **Python**: `vulture src/`
- **Rust**: `cargo +nightly udeps`

> **See the full guide:** Check [`guides/GATES_ECOSYSTEM_GUIDE.md`](https://github.com/co2-lab/anchors/blob/main/guides/GATES_ECOSYSTEM_GUIDE.md) in the repository for complete installation instructions, allowlists, and configuration files.

> **A declared gate that never runs is worse than no gate**, because it sits in
> the configuration as if it protected something. If the gate declares
> `when: [ci]` and your CI doesn't execute it, it's decoration. Check with
> `anchors doctor`.

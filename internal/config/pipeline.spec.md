<!-- @anchors
  code: PCSPP
  updated_at: 2026-11-19
  layer: config
-->
# Pipeline — the stages a workflow asks about, what each requires, and what validates it

> **Code**: `PCSPP`

## Overview

A release stage ran on a sha no review had seen, because nothing the workflow consulted said it should
not. `pipeline:` declares the stages a pipeline asks the Anchors about (`anchors pipeline`): the
branches each runs on, the conditions its sha must meet, and the command that runs or validates it.

```yaml
pipeline:
  unguarded:
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

A requirement is `flow` (the sha is the merge of a pull request `branch_flow:` allows), `gates: <phase>`
(the phase's gates pass), `passed: <stage>` (the head the sha was merged from passed that stage),
`from: <branch pattern>` (the sha was merged from that head), `version:` (what the version the sha
carries must be, against the product `versioning:` declares, VRSNN), `contains: <branch>` (the
sha has the head of that branch in it), or `any:` (one of several).

```yaml
    review-release:
      branches: [review]
      requires:
        - contains: production
        - version: { product: app, fresh: true, listed: { where: { target: review } } }
    production-release:
      branches: [production]
      requires:
        - version: { product: app, same_as_passed: review-release }
```

`version:` asks one or more of: `fresh` (no release has the version), `listed` (a listing gives it
an entry, VRSNN-B09), `same_as_passed: <stage>` (the version the commit that passed that stage
carried) and `expected` (the version the branch's strategy computes). The product may go unnamed
when `versioning:` declares one.

A list of requirements is one path of an `any:`, and it holds when all of its requirements hold:
a stage reached by two paths asks each path its own requirements.

```yaml
    production-release:
      requires:
        - any:
            - [{ passed: review-release }, { version: { same_as_passed: review-release } }]
            - [{ from: "hotfix/*" }, { version: { fresh: true } }]
```

`release: <product>` makes the stage write that product's tag
when it passes (RLSTA). `workflows:`
are the globs `pipeline-wired` reads, GitHub's by default; `unguarded:` names the workflows that call
no stage, each with its reason.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a stage | `branches`, `requires`, `run`, `verify` | a stage with neither `run` nor `verify` | this unit: refuses the configuration (`PCSPP-E01`) |
| a requirement | `flow`, a map of one key: `gates`, `passed`, `from`, `version`, `contains`, `any`, or a list of requirements | any other word or key, an empty list | this unit: refuses the configuration (`PCSPP-E01`, `PCSPP-E04`) |
| `version:` | `product`, `fresh`, `listed` (`file`, `keys`, `where`), `same_as_passed`, `expected`, at least one asked | a product not declared, several and none named, nothing asked, another key, a stage or a listing file that do not exist | this unit: refuses the configuration (`PCSPP-E02`) |
| `contains:` | one branch's name | a pattern | this unit: refuses the configuration (`PCSPP-E02`) |
| `release:` | a product of `versioning:` | any other name | this unit: refuses the configuration (`PCSPP-E03`) |
| `passed:` | a declared stage | a name no stage has | this unit: refuses the configuration (`PCSPP-E01`) |
| a branch pattern | as `branch_flow:` reads it (BFSBR) | neither glob nor regular expression | this unit: refuses the configuration (`PCSPP-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `PCSPP-B01` | A stage declares the branches it runs on, its requirements in order — `flow`, `gates:`, `passed:`, `from:`, `any:`, each written back as declared —, and the commands that run and validate it; the stages are listed sorted, the unguarded workflows keep their reasons, and the workflows default to `.github/workflows/*.yml` and `*.yaml`. (`StageNames`, `WorkflowGlobs`, `Requirement.String`) |
| `PCSPP-B02` | A stage asks of the version its sha carries — the product, and any of `fresh`, `listed` with its fields, `same_as_passed` and `expected` — and of a branch whose head the sha contains; each is written back as declared, the product left out when it was. (`VersionRequirement.String`) |
| `PCSPP-B03` | A stage names the product whose release it records (`release:`). |
| `PCSPP-B04` | A list of requirements is read as one path, all of whose requirements must hold, and written back in brackets; inside an `any:` it pairs the requirements of each path. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PCSPP-E01` | A stage has neither `run` nor `verify`; a requirement is an unknown word, an unknown key or empty; `passed:` names no stage; a branch pattern is neither glob nor regular expression. | The configuration does not load, naming the stage. | Each would let a stage run unconfronted, or be required by nothing that exists. |
| `PCSPP-E02` | `version:` names a product `versioning:` does not declare, names none among several, asks nothing, has a key it does not know, `same_as_passed:` names no stage, or `listed:` names no file while the product declares no `released:`; `contains:` is a pattern. | The configuration does not load, naming the stage and the field. | A version asked of nothing, or a branch no one can fetch, would let the stage run unconfronted. |
| `PCSPP-E03` | `release:` names a product `versioning:` does not declare. | The configuration does not load, naming the stage and the products. | A release of no product would write a tag no one reads. |
| `PCSPP-E04` | An empty list of requirements, or a requirement inside a list that would itself refuse the configuration. | The configuration does not load, naming the stage. | An empty path requires nothing, and would let the stage run by it. |

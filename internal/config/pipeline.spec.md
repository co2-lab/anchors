<!-- @anchors
  code: PCSPP
  updated_at: 2026-11-06
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
`from: <branch pattern>` (the sha was merged from that head), or `any:` (one of several). `workflows:`
are the globs `pipeline-wired` reads, GitHub's by default; `unguarded:` names the workflows that call
no stage, each with its reason.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a stage | `branches`, `requires`, `run`, `verify` | a stage with neither `run` nor `verify` | this unit: refuses the configuration (`PCSPP-E01`) |
| a requirement | `flow`, or a map of one key: `gates`, `passed`, `from`, `any` | any other word or key | this unit: refuses the configuration (`PCSPP-E01`) |
| `passed:` | a declared stage | a name no stage has | this unit: refuses the configuration (`PCSPP-E01`) |
| a branch pattern | as `branch_flow:` reads it (BFSBR) | neither glob nor regular expression | this unit: refuses the configuration (`PCSPP-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `PCSPP-B01` | A stage declares the branches it runs on, its requirements in order — `flow`, `gates:`, `passed:`, `from:`, `any:`, each written back as declared —, and the commands that run and validate it; the stages are listed sorted, the unguarded workflows keep their reasons, and the workflows default to `.github/workflows/*.yml` and `*.yaml`. (`StageNames`, `WorkflowGlobs`, `Requirement.String`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PCSPP-E01` | A stage has neither `run` nor `verify`; a requirement is an unknown word, an unknown key or empty; `passed:` names no stage; a branch pattern is neither glob nor regular expression. | The configuration does not load, naming the stage. | Each would let a stage run unconfronted, or be required by nothing that exists. |

<!-- @anchors
  code: PCSPA
  updated_at: 2026-11-06
  layer: comando
-->
# Pipeline — ask whether a stage may run, and record that it passed

> **Code**: `PCSPA`

## Overview

Each stage of a pipeline asks the Anchors before it runs, and tells it when it passed, against the
stages `pipeline:` declares (PCSPP). The first job of a stage runs `anchors pipeline can-run <stage>`,
and every other job needs it: a refusal runs nothing. The last job runs `anchors pipeline record
<stage>` (alias `done`), which runs the stage's `verify:` and records the outcome on the sha. A stage the Anchors runs
whole uses `anchors pipeline run <stage>`: can-run, its `run:`, its `verify:` and the record, in one
step.

What passed is recorded where the next stage reads it — on the sha, as the commit status
`anchors/pipeline/<stage>`, written and read through the `gh` CLI. A production stage requiring
`passed: <stage>` asks it of the sha itself — a later stage of the same branch reads the stage before
it there —, and else of the head the sha was merged from: the review commit that came into
production, whatever merge commit the pull request made.

A stage the configuration does not declare does not run.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the stage | a name `pipeline.stages` declares | any other | this unit: refuses (`PCSPA-E01`) |
| the branch and the sha | `--branch`, `--sha`, or `$GITHUB_REF_NAME`, `$GITHUB_SHA` | missing | this unit: refuses (`PCSPA-E02`) |
| the pull request of a sha, the statuses of a commit | what `gh api` answers | no `gh`, no answer | this unit: the requirement fails (`PCSPA-E03`) |
| `run:`, `verify:` | a shell command, run in the project's root | — | the project |
| `gates: <phase>` | the phase `anchors verify --phase <phase> --all` confronts | — | the configuration |

## Effects

| Effect | Description |
| --- | --- |
| `PCSPA-B01` | `can-run` passes a stage on a branch its `branches` match when every requirement holds, writing one line per requirement with why; on another branch it fails naming the stage's branches, and with a requirement that does not hold it fails naming it. `flow` holds for a branch the flow does not guard, or for the merge of a pull request the flow allows; `gates:` holds when the phase's gates pass. (`canRun`) |
| `PCSPA-B02` | `passed: <stage>` holds when the sha, or else the head of the merged pull request that made it, carries the status `anchors/pipeline/<stage>` as success; `from: <pattern>` holds when that pull request's head matches; `any:` holds when one of its requirements does. |
| `PCSPA-B03` | `record` — alias `done` — runs the stage's `verify:` and records the outcome on the sha as the commit status `anchors/pipeline/<stage>`: success when it passed; failure, and the command fails, when it did not. (`finish`) |
| `PCSPA-B04` | `run` confronts the stage as `can-run` does, runs its `run:`, then its `verify:`, and records the outcome; a `run:` that fails stops before `verify:` and records failure. |
| `PCSPA-B05` | `passed:` reads the sha itself before the head it was merged from: a stage recorded on the same sha — the production release before the version pointer flips — satisfies a later stage of that branch. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PCSPA-E01` | The stage is not declared in `pipeline.stages`. | The command fails naming the declared stages. | A stage the configuration does not declare is a stage no rule confronts. |
| `PCSPA-E02` | No branch or sha is given, by flag or environment. | The command fails naming the flags and the variables. | A stage with no commit named would be confronted against nothing. |
| `PCSPA-E03` | GitHub cannot be asked — no `gh`, no answer —, or `record` names a stage with no `verify:`. | The requirement fails, and the stage with it; `record` fails pointing to `pipeline run`. | A requirement no one could answer would let the stage run unconfronted. |

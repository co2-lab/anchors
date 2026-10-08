<!-- @anchors
  code: BDGRN
  updated_at: 2026-10-08
  layer: comando
-->
# BudgetRun — run a suite's files fastest first until a time budget is spent

> **Code**: `BDGRN`

## Overview

`anchors test --budget 60s` and `anchors mutation --budget 10m` run the selected suites' files
fastest first, in batches through each suite's `run_changed:`, until the budget is spent; what did
not fit is left for a later run, like a smoke run of whatever fits. There is no priority besides time:
the project chooses the budget.

The order comes from the times earlier runs recorded in the map, per suite: a test file's from its
JUnit cases, a code file's from Anchors timing its mutation run — a mutation report carries no time
per file, so each file's mutation runs alone. Files no suite has timed yet go last, one per batch, and
their run is what times them.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the budget | a positive duration | zero, which means no budget | the command: without `--budget` the suite runs as always |
| the suite | one that declares `run_changed:` and a report | none of them | this unit: refuses, saying which is missing |
| the map | the project's map, with times from earlier runs or without | no map | this unit: refuses, saying to build it |

## Effects

| Effect | Description |
| --- | --- |
| `BDGRN-B01` | The plan is the suite's timed files, fastest first with ties by path, then the files no suite has timed, by path; a file only other suites timed, a support file, and a file of another kind are left out. |
| `BDGRN-B02` | A batch takes, from the front of the plan, the timed files whose times fit in what remains; when even the fastest does not fit, nothing more runs. After the timed files, the untimed ones run one per batch. |
| `BDGRN-B03` | A batch still running at the deadline is stopped with its whole process group — first a TERM, so a tool that works in place can restore the source, then a KILL of whatever is left after a grace of 10 seconds — and nothing of it is ingested. On Windows, which has no TERM to trap, the process is killed at once. |
| `BDGRN-B04` | Each batch runs through `run_changed:` with its files and is ingested as a partial run; the end reports how many files ran, how many were left, and whether the deadline cut a batch. |
| `BDGRN-B05` | A mutation budget runs one file per batch, and records in the map, under the suite, how long that file's run took. |
| `BDGRN-B06` | A batch that fails does not stop the budget: the next batches run, and the command fails at the end naming the suites with a failed batch. |
| `BDGRN-B07` | A suite with no `run_changed:` or no report, or a budget together with `--changed`, is refused before anything runs. |
| `BDGRN-B08` | The plan holds only the files the suite's `paths:` cover. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BDGRN-E01` | REF[BDGRN-B07]: a suite the budget cannot run in batches | — | — |
| `BDGRN-E02` | The map cannot be loaded | an error saying to build it | Without the recorded times there is no order to follow |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

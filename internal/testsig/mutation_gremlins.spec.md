<!-- @anchors
  code: GRING
  updated_at: 2026-09-28
  layer: infra
-->
# GremlinsIngest — the mutation score per file, read from a gremlins report

> **Code**: `GRING`

## Overview

In Go the canonical mutation format does not exist in practice: the dominant runner, gremlins, writes its own
report, and no converter to the canonical format exists. Without reading it, the advice to "run the mutation
tool and ingest the report" led a Go project nowhere. This unit reads the gremlins report into the same
per-file mutation result the canonical reading produces, so that nothing downstream knows which tool it came
from. It is chosen when the mutation gate declares `format: gremlins`.

The two formats differ in shape: gremlins lists its files, while the canonical format indexes them by path.
The shapes cannot be confused by the decoder, which is what lets a wrong `format:` fail with a message that
explains it. Gremlins writes its status as human text, and it keeps its thresholds in its own configuration,
not in the report.

## Effects

| Effect | Description |
| --- | --- |
| `GRING-B01` | The report's listed files are read into one mutation result per file. |
| `GRING-B02` | A killed mutant and a timed-out mutant both count as killed. The timed-out are also counted apart, since under load they inflate the score. |
| `GRING-B03` | A mutant that lived counts as survived, and its line is recorded. |
| `GRING-B04` | A mutant that was not viable, was only runnable or was skipped, or has a status the unit does not know, stays out of the score. |
| `GRING-B05` | The status is matched ignoring case and spaces. |
| `GRING-B06` | The thresholds stay zero, which means absent: the engine's default decides. |
| `GRING-B07` | A file's path is normalized as in the canonical reading: a leading `./` is dropped, and a path containing `/src/` is cut to start at `src/`. |
| `GRING-B08` | A mutant not covered by any test is counted apart as no-coverage and does not enter the score, and its line is recorded, as in the canonical reading. |
| `GRING-B09` | A file where no mutant ran scores 100, as in the canonical reading; its no-coverage count stays recorded. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the report file | a gremlins JSON report with at least one file | a canonical-format report, or no JSON at all | this unit: refuses it with an error that names `format:` |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GRING-I01` | When any mutant ran, the score is killed over killed plus survived, times one hundred. | reads files of different mixes and verifies each score against its counts |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GRING-X01` | The score is recomputed from the mutants' statuses; the report's own efficacy figure is not used. | Every format is scored by the same arithmetic, so the gate compares like with like. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GRING-E01` | The report is not a gremlins report — a canonical-format report under `format: gremlins`, or no JSON. | An error that names the expected shape and the `format:` declaration. | A wrong `format:` is the likely mistake, and silence would read as a project with no file. |
| `GRING-E02` | The report has no file. | An error pointing at the gremlins run. | An empty report is a run that produced no mutant, not a proven project. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/testsig/mutation.go` | `MutationReport`, `FileMutation`, `normalizeMutationPath` | infra — the shared mutation result and path normalization |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

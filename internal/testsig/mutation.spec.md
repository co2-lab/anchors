<!-- @anchors
  code: MTINM
  updated_at: 2026-10-03
  layer: infra
-->
# MutationIngest — the mutation score per file, read from a Mutation Testing Elements report

> **Code**: `MTINM`

## Overview

Line coverage says a line RAN; it does not say a test CHECKED its result. Mutation answers that question:
change the line, and if the tests stay green they do not prove it. Anchors does not run mutation and knows no
tool: it reads the report. The canonical format is Mutation Testing Elements, an open schema several tools
emit; the format is chosen by the mutation gate's configuration, and the empty choice means this one, so a
project that never declared a format keeps its behaviour.

Per source file, the unit sorts the mutants by what they say about the tests. A killed mutant, or one that
timed out, was noticed. A survivor ran under a test that did not notice it, and its line is kept so the author
can act. A mutant no test executed is counted apart and left out of the score: whether a test runs the line is
the coverage gate's question, and charging it here buried the real survivors under files no test reaches. A
mutant the tool ignored is counted apart too, and mutants that failed to compile or run say nothing about the
tests. The score is the share of killed among the mutants that ran; a file where none ran scores 100, and the
separate counters are what tell that 100 apart from a fully proven file.

The project's thresholds travel in the report itself, so the unit reads them there rather than asking the
project to declare them twice.

## Effects

| Effect | Description |
| --- | --- |
| `MTINM-B01` | `ParseMutationFormat` reads the canonical format for the names `mutation-testing-elements`, `mte`, `stryker` and the empty name, ignoring case and surrounding spaces; `ParseMutation` reads it with no name given. |
| `MTINM-B02` | A killed mutant and a timed-out mutant both count as killed. The timed-out are also counted apart, since under load they inflate the score. |
| `MTINM-B03` | A surviving mutant counts as survived, and its starting line is recorded. |
| `MTINM-B04` | A mutant no test covered is counted apart and does not enter the score, and its starting line is recorded. |
| `MTINM-B05` | A mutant the tool ignored is counted apart and does not enter the score. |
| `MTINM-B06` | A mutant that failed to compile or to run is not counted and does not enter the score. |
| `MTINM-B07` | The score is killed over killed plus survived, times one hundred; a file where no mutant ran has no score (zero, not written to the map): nothing was measured. |
| `MTINM-B08` | The low and high thresholds are read from the report when it carries them. |
| `MTINM-B10` | A report in the format — it carries its schema version — with no file is a run that had nothing to mutate: it is read as a report of no file, not refused. |
| `MTINM-B09` | A file's path is normalized: a leading `./` is dropped, and a path containing `/src/` is cut to start at `src/`. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the report file | a Mutation Testing Elements JSON report with at least one file | any other JSON, or no JSON at all | this unit: refuses it with an error |
| the format name | the canonical names, `gremlins`, or empty | any other name | this unit: refuses it naming the accepted formats |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MTINM-I01` | A file where nothing ran is told apart from a proven one: it has no score, and its uncovered or ignored mutants stay counted. | reads a file whose mutants were all ignored and verifies it has no score, with the ignored count kept |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MTINM-X01` | Thresholds absent from the report stay zero: the unit never invents a threshold. | Zero means absent, and the engine's default decides; a made-up threshold would be a rule nobody declared. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MTINM-E01` | The format name is not one the unit knows. | An error naming the accepted formats and where to declare them. | Otherwise the operator guesses the name. |
| `MTINM-E02` | The report is not JSON of the canonical format. | An error naming the expected format. | Accepting it silently would record a phantom score. |
| `MTINM-E03` | The report has no file and no schema version. | An error pointing at the tool's output format. | Without the version it is a tool emitting another format; a report WITH it and no file is a run with nothing to mutate, read as an empty report (`MTINM-B10`). |
| `MTINM-E04` | The report file cannot be read. | The read error is returned. | A missing report must not look like a clean run. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/testsig/mutation_gremlins.go` | `parseGremlins` | infra — the reading of the gremlins format, chosen by name |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

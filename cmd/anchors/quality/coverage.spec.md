<!-- @anchors
  code: CVCMC
  updated_at: 2026-09-26
  layer: comando
-->
# CoverageCommand — answers the confidence questions from the ingested signals: by scenario, by line, of the diff and the delta

> **Code**: `CVCMC`

## Overview

The test reports ingested into the map say which scenario codes a green test proved, how many lines of
each code file ran, and how many mutants survived. The coverage command reads those signals and answers
four questions.

For one spec: is each scenario it declares proven by a green test? Only the codes the spec declares for its
own unit count; a code of another unit that the prose merely cites is not a scenario of this spec, and
counting it would let the map claim a proof no test gave. When the spec changed since the signal was
ingested the answer is flagged as stale.

For the whole project, the panorama: which specs have an unproven scenario, which code files are below the
line threshold, and which are below it by mutation score. "Nothing was measured" and "everything is fine"
are opposite conclusions and are never printed the same way. Each section is capped, and a closing summary
gives the full counts.

For a change: are the lines I changed covered? The diff (from git or from a diff file) is crossed with a
coverage report; only instrumented changed lines count, and the command fails when the covered share is
below the threshold. This question needs no map.

Against the previous ingestion: did any file lose line coverage? A drop fails the command.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the spec argument | a spec path, relative to the root or not | a path that is not a node of the map | this unit: it refuses (`CVCMC-E02`) |
| the diff | a git reference, or a unified diff file | a reference in a project not under git; an unreadable diff file | this unit: it refuses and explains (`CVCMC-E05`, `CVCMC-E06`) |
| the coverage report | an lcov file, mandatory for the diff question | no report, or one that cannot be read | this unit: it refuses (`CVCMC-E03`, `CVCMC-E04`) |
| the threshold | a percentage, 70 when not given | — | the flag parser rejects a non-number |

## Effects

| Effect | Description |
| --- | --- |
| `CVCMC-B01` | For one spec, lists each scenario code it declares, in order, as proven by a green test or not, then how many are proven out of how many and which are missing. |
| `CVCMC-B02` | A code another unit owns, cited in prose, is not a declared scenario of the spec. |
| `CVCMC-B03` | A spec with no unit code keeps every code it declares. |
| `CVCMC-B04` | When the spec changed since its signal was ingested, the proven count is flagged as a stale signal. |
| `CVCMC-B05` | A spec that declares no scenario code is reported as having nothing to cover. |
| `CVCMC-B06` | The panorama prints three sections — specs with an unproven scenario, code files below the line threshold, code files below the threshold by mutation score with their survivors — marking stale signals, and closes with a summary of the counts. |
| `CVCMC-B07` | When every measured file is above the threshold the panorama says so with the count, and still names the surviving mutants when there are any. |
| `CVCMC-B08` | Each panorama section shows at most fifteen entries and counts the rest, while the summary counts all of them. |
| `CVCMC-B09` | The diff coverage lists, per changed file that has coverage, its instrumented changed lines, the share covered and the uncovered line numbers, then the total, and passes at or above the threshold; changed lines that are not instrumented and changed files without coverage are left out, and a file matches its coverage by exact path or by path suffix. |
| `CVCMC-B10` | When no changed line is instrumented, the diff coverage says there is nothing to cover and passes. |
| `CVCMC-B11` | Changed lines covered below the threshold fail the diff coverage with exit status 1. |
| `CVCMC-B12` | The delta compares each code file with a previous line coverage against it, ignores files with none, and when nothing dropped says so with how many improved. |
| `CVCMC-B13` | A file that lost more than a hundredth of a point of line coverage is listed with its old and new coverage, and the delta fails with exit status 1 naming the worst drop. |
| `CVCMC-B14` | The diff coverage lists the changed files in path order, and a changed file whose path suffix matches several coverage entries is crossed with the first of them in path order, so the same inputs print the same report on every run. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CVCMC-I01` | The panorama never reports "nothing measured" as "nothing below the threshold": with no line or mutation signal it says NOT measured, in its section and in the summary. | a map whose one code file has no signal prints the NOT measured lines and no clean verdict |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CVCMC-X01` | The diff coverage does not read the map. | It crosses a diff with a coverage report; demanding a map would refuse the question in CI jobs that never built one. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CVCMC-E01` | The map does not exist or cannot be read, for any question but the diff. | Error pointing at `anchors map build`. | The scenario, panorama and delta answers are read from the map's signals. |
| `CVCMC-E02` | The spec asked for is not a node of the map. | Error saying it is not in the map. | Its declared codes and proven codes are both read from the map's node. |
| `CVCMC-E03` | The diff question is asked without a coverage report. | Error saying the report is mandatory with the diff. | The diff alone says what changed, not whether it ran. |
| `CVCMC-E04` | The coverage report cannot be read. | Error naming the report parse. | An unread report would pass every changed line as not instrumented. |
| `CVCMC-E05` | The diff is asked against a git reference and the project is not under git. | Error explaining why the diff could not be taken and pointing at the diff file alternative. | Without history there is nothing to compare against; the diff file serves any other source. |
| `CVCMC-E06` | The diff file cannot be read. | The diff coverage fails. | An empty diff would pass as "nothing to cover". |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/testsig/lcov.go` | `ParseLCOV` | infra — reading the coverage report |
| DEP2 | `internal/testsig/diff.go` | `GitDiff`, `ParseDiffFile` | infra — the changed lines |
| DEP3 | `internal/testsig/junit.go` | `CodesInCase` | infra — the scenario codes a text carries |
| DEP4 | `internal/gitmeta/availability.go` | `Check`, `Explain` | infra — why git could not answer |
| DEP5 | `internal/mapx/store.go` | `Load` | mapa — the ingested signals |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

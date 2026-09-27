<!-- @anchors
  code: RPRTS
  updated_at: 2026-09-26
  layer: comando
-->
# Reports — markdown perspectives on what Anchors already measures, written into docs

> **Code**: `RPRTS`

## Overview

The map, the ingested test signals, the gates, the health validator, the issues and the task queue each
answer a question about the project, but reading them means knowing several commands. The report command
cuts those sources into six perspectives and writes each one as a markdown document: tests (the confidence
panel of the deliverable), quality (the gates' verdict and the debt), structure (kinds, governance and
identity), configuration (what the project declares and what it misses), issues (the tracked work) and
inconsistencies (everything the validators flag). None of them invents data: each is a recut of something
Anchors already measured.

The reports are terminal documents for people, versioned to give history; they are regenerated on every
run and never become anchors of their own. A single perspective is written under the docs folder, or to a
chosen file; the all form writes every perspective and an index into one folder. They all share one header,
which says when and by what they were generated and that they must not be edited by hand.

This spec coordinates the two files that make the command: one holds the table of perspectives, the
command shell and the tests perspective; the other holds the renderers of the other five perspectives and
the helpers they share.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the perspective | one of the six names, or all | any other name | the command tree: an unknown subcommand is not found |
| the map | the project's map at its default path, or the one given | no map | this unit: it refuses (`RPRTS-E01`) |
| the configuration | the project's anchors.yaml, or none | — | this unit: every perspective says what is missing without it (`RPRTS-B14`) |
| the destination | a file path that can be created | a path under an existing file | this unit: it fails (`RPRTS-E02`) |

## Effects

| Effect | Description |
| --- | --- |
| `RPRTS-B01` | A single perspective is written to docs or to the chosen file, and the path written is printed relative to the root. |
| `RPRTS-B02` | The all command writes every perspective and an index into docs/anchors, the index linking each perspective with its description. |
| `RPRTS-B03` | Every perspective of a configured project opens with the same header — title, generation time, the source, the do-not-edit note — and closes with an Anchors footer line. |
| `RPRTS-B04` | The tests perspective merges execution by layer with a total row, warns when tests fail, and warns about test files whose signal is stale. |
| `RPRTS-B05` | The tests perspective counts only measured specs: it gives the requirements proven out of those declared, lists each spec with a gap and its unproven codes, and reports the unmeasured specs apart. |
| `RPRTS-B06` | The tests perspective lists the code files below 70% of lines, at most fifteen, and every file whose line coverage dropped since the previous ingestion. |
| `RPRTS-B07` | The tests perspective with nothing ingested says so in each section and claims no failure and no stale signal. |
| `RPRTS-B08` | The quality perspective gives, per gate, its force and its pass, fail and pending counts, the promotable or barred verdict, the divergences marked as blocking or informative, and the targets awaiting AI judgment. |
| `RPRTS-B09` | The structure perspective counts nodes by kind, gives the totals, lists the guides with how many nodes each governs, and the identity collisions, the orphans (specs with no implementation) and the missing identities the health validator found. |
| `RPRTS-B10` | The configuration perspective lists the declared layers, governance rules, gates with their kind and force, the co-location setting, the guides with no governance and the kinds with no gate. |
| `RPRTS-B11` | The issues perspective counts issues per state, then lists the open issues waiting on the user, the open agent work and the deferred ones, each list only when not empty, then the live tasks or an empty queue. |
| `RPRTS-B12` | The inconsistencies perspective lists every health finding grouped by check, the failures of the quality gates, and the total. |
| `RPRTS-B13` | A finding section takes the findings whose check starts with its name, titles itself with their count, caps the list at 25 saying how many are left, and is absent when there is none. |
| `RPRTS-B14` | Without anchors.yaml the perspectives say what is missing instead of inventing: the quality one that no gate is declared, the configuration one that the file was not found, the issues one its zero counts and empty queue, and the structure and inconsistencies ones that the health validator needs the file; the all form still writes every perspective. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RPRTS-I01` | An issue waiting on the user is listed only as the user's, never also as agent work. | the issues perspective of a project with a user-owned issue in todo names it exactly once |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RPRTS-X01` | Generating the reports leaves the map as it was: they write only the report files. | The reports describe the measured state; a report that changed what it reads would describe something else by the time it is read. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RPRTS-E01` | The map does not exist or cannot be read, for a single perspective or for all. | Error naming the map load and pointing at `anchors map build`. | Every perspective is a cut of the map; without it there is nothing to report. |
| `RPRTS-E02` | The destination cannot be created or written. | The report fails with the write error. | A report that says "generated" and wrote nothing would pass for an up-to-date document. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/store.go` | `Load` | mapa — the map and the ingested signals |
| DEP2 | `internal/gate/profile.go` | `Aggregate` | gate — the quality verdict |
| DEP3 | `internal/health/health.go` | `Diagnose` | infra — the structural findings |
| DEP4 | `internal/issue/issue.go` | `List`, `FileOwner` | infra — the issues by state and owner |
| DEP5 | `internal/queue/queue.go` | `List` | infra — the task queue |
| DEP6 | `internal/config/config.go` | `Load` | config — the declared layers, governance and gates |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

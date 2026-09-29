<!-- @anchors
  code: PRJTS
  updated_at: 2026-09-29
  layer: gate
-->
# ProjectTests — the gates read the project's tests through the source the project declares

> **Code**: `PRJTS`

## Overview

The gates that read a test's title — `feature-test-match`, `test-traceable`, `scenario-coverage` and
`flag-covered` — get the project's tests from here. How a test is written belongs to the project's
test library, not to the engine: the project declares it in `dialect.tests` (a pattern or a script, see
the `testlist` package), or its dialect family does. The gates used to carry Jest's and Go's call
syntax themselves, and a project on any other library had no title read at all.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | a project configuration, or nil | — | this unit: nil declares no source |
| the map | the scan's map, or nil | — | this unit: with no map there are no test files to read |

## Effects

| Effect | Description |
| --- | --- |
| `PRJTS-B01` | The source is the project's own `dialect.tests`, or its family's when it declares none; with neither, nothing is declared and the gates fall back to what they can do without titles. |
| `PRJTS-B02` | A pattern reads the map's test files, and only those. |
| `PRJTS-B03` | The tests are read once per scan: the same root, map and configuration reuse the reading, and a new map reads again. |
| `PRJTS-B04` | The tests of a set of files come in the order the files are given, and in each file in the order of their lines; tests of other files are left out. |
| `PRJTS-B05` | An error of the source is returned with the declaration, never taken for a project without tests. |
| `PRJTS-B06` | Support files are not read as tests, and are dropped from any list of test paths the gates confront. |
| `PRJTS-B07` | Under `--index` the tests are listed from what the commit records — the same copy of the index a `workdir: index` gate runs in —, so a test's line is a line of the content the gates read. (`projectTests`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRJTS-E01` | REF[PRJTS-B05]: the source fails or answers outside its contract | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `DialectFor` | core — the project's dialect and its `tests` source |
| DEP2 | `internal/mapx/model.go` | `Graph`, `KindTest` | core — the map's test files |
| DEP3 | `internal/testlist/testlist.go` | `List`, `Source`, `Test` | infra — reading the tests through a pattern or a script |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

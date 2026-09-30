<!-- @anchors
  code: SLCTN
  updated_at: 2026-09-29
  layer: comando
-->
# RunSelection — a run takes only what is stale and below the minimum, unless told otherwise

> **Code**: `SLCTN`

## Overview

`anchors test` and `anchors mutation` run, by default, only the files whose last result is both
stale and below the minimum, plus the files never measured: what is known and current is skipped,
so the everyday run is light. Two questions place each file in one of four boxes — is its result
**fresh** (measured at the current version) and is it **passing** (a test file with no failing case,
a code file whose mutation score is at the floor or above) — and each flag opens one side:

| | fresh | stale |
| --- | --- | --- |
| **passing** | `--include-fresh --include-passing` | `--include-passing` |
| **below the minimum** | `--include-fresh` | run by default |

`--skip-unmeasured` leaves out the files never measured; `--all` runs every file through the
suite's `run:`, as before this selection. `--changed` and `--target` choose the files their own
way and are not filtered by state. A suite with no `run_changed:` cannot run a subset and runs
whole, saying so.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | the project's map, with the signals of earlier runs or none | no map | this unit: refuses, saying to build it or to use `--all` |
| the flags | any combination of the state flags, or `--all`, or `--changed` | `--all` or `--changed` together with the state flags | the command: refuses before anything runs |

## Effects

| Effect | Description |
| --- | --- |
| `SLCTN-B01` | A file is placed by two answers — fresh or stale, passing or below the minimum — into one of four boxes; a file with no result is never measured. |
| `SLCTN-B02` | By default a run takes the files stale and below the minimum and the ones never measured; `--include-fresh` adds the fresh ones below the minimum, `--include-passing` the stale passing ones, both together the fresh passing ones too, and `--skip-unmeasured` leaves out the never measured. |
| `SLCTN-B03` | A test file is stale when it or anything it exercises changed since its result; its passing is read from the suite's own layer when the file ran in several; a file whose cases were all skipped is not passing. |
| `SLCTN-B04` | A code file's mutation result is stale when measured at another revision of the file — its own, not the one a later coverage ingestion moves —, and one taken under load — more timed-out mutants than the gate's ceiling — counts as stale; it passes at the report's floor or 70% when the report declares none, and a file where no mutant ran passes. |
| `SLCTN-B05` | A test file belongs to the suite that timed it; before any time, to the suite whose layer ran it; a file no suite ran is offered to every suite. |
| `SLCTN-B06` | Support files and targets the matching gate declares in `no_signal` are never run. |
| `SLCTN-B07` | The run says how many files it selected and, by state, how many it left out and the flag that would take them. |
| `SLCTN-B08` | A suite with no `run_changed:` runs whole, saying it cannot run a subset; with a budget it is refused before anything runs. |
| `SLCTN-B09` | When the selection takes no file, the suite says so and runs nothing. |
| `SLCTN-B10` | The selected files run in as few batches as the command-line ceiling allows, each ingested as a partial run; a file longer than the ceiling still runs alone. |
| `SLCTN-B11` | `--all` runs the suites whole through `run:`; it is refused with `--changed` or with a state flag, and the state flags are refused with `--changed`. |
| `SLCTN-B12` | A suite that declares `paths:` is handed only the files they cover, test and code files alike; the files outside them are neither run nor counted as left out. |
| `SLCTN-B13` | A file the matching gate does not confront — outside its `on:` kinds, carrying one of its `exclude_tags`, lacking all of its `tags` — is never run for it; a gate entry that declares no kinds filters nothing. |
| `SLCTN-B14` | A test run also takes the test files of the rules a contract field changed since the last commit reaches, when the suite runs them — saying which were added and why —, whatever state their own file is in. |
| `SLCTN-B15` | The map is read as a build would leave it now: a file edited since the last build has no result — the build drops the result of another version — and is selected as never measured; an unchanged file keeps its result. (`currentRevs`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SLCTN-E01` | The map cannot be loaded for a selective run | an error saying to build it, or to run with `--all` | Without the recorded results there is no state to select by |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Suite`, `Gate`, `DefaultTimeoutCeiling`, `TimeoutCeilingOrDefault`, `NoSignalFor` | core — the suites and the gates' ceiling and `no_signal` |
| DEP2 | `internal/mapx/model.go` | `Graph`, `Node`, `SignalStale`, `KindTest`, `KindCode` | core — the files and their last results |
| DEP3 | `internal/mapx/evidence.go` | `EvidenceStaleFor` | core — a test stale through what it exercises |
| DEP4 | `cmd/anchors/mapcmd/ingest.go` | `SuiteKey` | comando — the key a suite's times are kept under |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

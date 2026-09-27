<!-- @anchors
  code: STPRS
  updated_at: 2026-09-27
  layer: comando
-->
# SuiteProxy — runs the test and mutation suites the project declared, and binds their reports to the map

> **Code**: `STPRS`

## Overview

Producing a test signal used to take two steps, and the second was the one people forgot: the suite ran,
nobody ingested the report, and the gates went on saying "no signal ingested" as if the run had never
happened. The test and mutation commands close that gap. They are proxies: the command belongs to the
project, declared in its configuration (the tests section for `test`, the mutation section for
`mutation`), and Anchors only runs it at the project root and binds the reports it left to the map.

Anchors does not know the stack. It never builds a runner command and never guesses one; when nothing is
declared it shows what to declare and stops. The positional arguments select layers, and a workspace filter
(and, for mutation, a scope filter) narrows further.

There are two modes, the same as the check's. The full mode runs each suite's declared command. The
incremental mode takes the changed files, computes the same impact path the check's incremental mode
uses, and runs the suite's incremental command with those files in its placeholder — code and tests for
the test command, code only for mutation, since mutating the test would mutate the measuring instrument.
Without an incremental command declared it refuses rather than silently running everything.

Each report the run wrote is ingested even when the suite failed, because a failing run is exactly the one
with something to say; a report older than the run is never ingested, because it would record an earlier
number as if it were now. The first failing suite stops the ones after it. On success the run can chain
coverage, or the check over the run's own scope: the same changed files, or the full sweep.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared suites | entries with a layer, a command, and optionally a workspace, a scope, an incremental command and report paths | a configuration that declares none | this unit: it shows how to declare and fails (`STPRS-E02`) |
| the layer, workspace and scope filters | values that some declared suite carries | a value no suite declares; a combination no suite matches | this unit: it refuses (`STPRS-E03`, `STPRS-E04`) |
| the changed files | files of the project | a project without a map | this unit: it refuses (`STPRS-E09`) |
| the target | any text, used only where the command declares its placeholder | an empty target for a command that needs one | this unit: it refuses (`STPRS-E05`) |
| the chain | `coverage` or `check`, comma separated, or nothing | any other name | this unit: it refuses (`STPRS-E08`) |

## Effects

| Effect | Description |
| --- | --- |
| `STPRS-B01` | Runs each selected suite's declared command through the shell at the project root, in the file's order, under a header naming the command, the workspace, the layer and the scope. |
| `STPRS-B02` | The reports a suite declares, written by this run, are ingested into the map, resolved against the project root. |
| `STPRS-B03` | A suite that passes and declares no report says it passed and that nothing was ingested. |
| `STPRS-B04` | A failing suite still has its report ingested, fails the command naming its layer, and the suites after it do not run. |
| `STPRS-B05` | Without changed files the full command runs, even when an incremental one is declared. |
| `STPRS-B06` | With changed files the incremental command runs with the files of the impact path in its placeholder, where the command declares it, as absolute paths with forward slashes. |
| `STPRS-B07` | The test command receives the code and test files of the impact path; the mutation command receives only the code files. |
| `STPRS-B08` | A command declaring a target placeholder gets the given target in its place; a command without the placeholder is left as declared. |
| `STPRS-B09` | After every suite passed, a chain naming coverage runs the coverage answer; an empty chain does nothing. |
| `STPRS-B10` | An impact path with no code file runs nothing. |
| `STPRS-B11` | After every suite passed, a chain naming check runs the check over the run's own scope: the changed files of an incremental run, the full sweep otherwise. |
| `STPRS-B12` | An incremental run hands each suite only the impact files its `paths:` cover; a suite the impact path does not reach says so and runs nothing. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `STPRS-I01` | Only a report written during this run reaches the map; one older than the run's start is left out and named. | a suite whose report predates the run leaves the map with the earlier signal and prints that the report is older than the run |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `STPRS-X01` | The command that runs is exactly the declared one with only its placeholders filled: nothing is appended, reordered or guessed. | The stack is the project's; a proxy that edited the command line would be the framework choosing how the project runs its tests. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `STPRS-E01` | The configuration does not load. | Error naming the configuration file. | The suites are declared there. |
| `STPRS-E02` | The section for this command declares no suite. | Prints how to declare the section, with the report it needs, and fails. | Anchors does not guess how a project runs its tests; the answer is what to declare, not a complaint. |
| `STPRS-E03` | A layer, workspace or scope filter names something no suite declares. | Error naming it and listing the declared layers, workspaces and scopes. | A typo must not pass for "ran nothing and passed". |
| `STPRS-E04` | Declared filters that match no suite together are refused. | Error naming the filters. | A silent success with nothing run would read as a green suite. |
| `STPRS-E05` | The command declares a target placeholder and no target was given. | Error asking for the target. | An empty target would make a mutation tool mutate the whole project. |
| `STPRS-E06` | Changed files were given and the suite declares no incremental command. | Error showing how to declare one. | Falling back to the full run would be expensive and would lie about what ran. |
| `STPRS-E07` | The command line is longer than the platform's ceiling (overridable by an environment variable). | Error before running, with the ceiling and the ways out. | Past the ceiling the command fails without writing anything, and that failure looks like a red test. |
| `STPRS-E08` | The chain names a command other than check or coverage. | Error naming what is accepted. | The chain is not a disguised shell runner. |
| `STPRS-E09` | The incremental mode without a map points at the map build. | Error pointing at `anchors map build`. | The impact path is read from the map. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `SelecionaSuites`, `DeclaredLayers` | config — the declared suites and their filters |
| DEP2 | `cmd/anchors/mapcmd/ingest.go` | `IngestArtifacts` | comando — binding the reports to the map |
| DEP3 | `cmd/anchors/quality/check.go` | `selectNodes` | comando — the impact path the check uses |
| DEP4 | `internal/mapx/store.go` | `Load` | mapa — the map of the impact path |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

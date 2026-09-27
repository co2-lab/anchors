<!-- @anchors
  code: NGSTI
  updated_at: 2026-09-27
  layer: comando
-->
# Ingest — binds the test and log signals the project produced to the nodes of the map

> **Code**: `NGSTI`

## Overview

Anchors does not run the project's tests. The project's runner produces reports, and this
unit reads them and writes what they say into the map: the execution report (which test
cases passed, failed or were skipped), the line coverage report, and the mutation report.
From the execution it also derives coverage per scenario: a rule code a spec declares is
proven when a passing case names it. The same logic is what the test and mutation commands
call right after running a suite, so running and ingesting become one operation there.

Each report is a suite, keyed so that the key is the same on every machine: a report inside
the repository is keyed by its path from the root, one outside it by its file name. A full
run of a suite inside the repository replaces what an ad-hoc report from outside left in the
map. Report paths are resolved to exactly one node, deterministically.

Ingesting by hand can put an old report in the map as if it were the last run. So a manual
ingestion warns, and the project can declare that it must refuse instead; a project that has
no suite declared, or no configuration, is never asked to use a command it cannot run.

The unit also ingests the application's logs: it scans the declared log files for failure
codes and binds each occurrence to the spec that declares that failure, stamped with the
spec's revision, so the failures review can later ask for a conclusion. A failure code no
spec declares is reported, never bound to an invented owner.

## Effects

| Effect | Description |
| --- | --- |
| `NGSTI-B01` | In one pass of `IngestArtifacts` (the entry the ingest command and the test commands share), the execution report reaches the test node (passed, failed and skipped counts), the passing case codes are proven on the spec that declares them, and the coverage and mutation reports reach the code node. |
| `NGSTI-B02` | An execution report whose cases name no test node of the map warns that no test file matched. |
| `NGSTI-B03` | When the project declares its code length, the case names of the execution report are read with it, so a rule code of that length is proven by its passing case. |
| `NGSTI-B04` | An ingestion not run by the test commands warns on the error stream and proceeds, by default. |
| `NGSTI-B05` | An ingestion run by the test commands never complains, even when the project refuses manual ingestion. |
| `NGSTI-B06` | A project with no test or mutation suite declared, or with no configuration at all, is not asked to use the test commands. |
| `NGSTI-B07` | The suite key is the report's path from the project root; a report outside the repository is keyed by the external prefix and its file name, the same on every machine. |
| `NGSTI-B08` | A full run of a suite inside the repository drops from the map the suites ingested from outside it, and says so. |
| `NGSTI-B09` | A partial run, or another report from outside the repository, drops no external suite. |
| `NGSTI-B10` | When two paths of one report name the same node, the path equal to the node id wins, whatever order the report lists them in. |
| `NGSTI-B11` | A line coverage entry for a file changed after the report was written is marked as predating the file; a file older than the report, or one not on disk, is not. |
| `NGSTI-B12` | Log ingestion binds each occurrence of a declared failure code to the spec that declares it, with the number of occurrences, stamped with the spec's current revision. |
| `NGSTI-B13` | Log ingestion reports the failure codes found in the logs that no spec declares. |
| `NGSTI-B14` | A test file's run time in the report is the sum of its cases' times, kept in the map under the report's suite. |
| `NGSTI-B16` | An ingestion that follows a run of `anchors test` or `anchors mutation` first takes the revs of the tree as it is now, so a proof of a spec edited after the last `map build` stays fresh once the map is rebuilt; a manual ingestion keeps the map's revs, since its report may be older than the tree. |
| `NGSTI-B15` | A report's signals are kept under its path relative to the root; a report outside the repository is kept under `external/` and its file name. (`SuiteKey`) |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the report options | at least one of an execution, coverage or mutation report, or the log option | none of them | this unit: refuses with the usage |
| the reports | readable reports in the formats the parsers read; the mutation format is the one the project declares | a missing or malformed report | the report parsers: this unit fails naming the format |
| the map | the project's map | a missing map | this unit: fails asking for the map build |
| the configuration | the project's configuration, or none | — | this unit: without configuration the permissive defaults hold and nothing is demanded |
| the log paths | the globs declared under the logs configuration | a project that declares none | this unit: refuses naming the missing declaration |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `NGSTI-I01` | Ingesting the logs replaces the occurrences bound before on every spec: the same logs ingested twice leave the same counts. | ingest the same log twice with the real command and verify the counts did not double |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `NGSTI-X01` | A failure code no spec declares is bound to no node of the map. | Binding it would invent an owner for an orphan failure; it is reported instead. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `NGSTI-E01` | No report and no log option is given. | The command refuses, naming the options it expects. | There is nothing to ingest. |
| `NGSTI-E02` | A manual ingestion in a project that declares manual ingestion must be refused. | The ingestion is refused, and the message names the test command to use instead. | The project chose that the map only carries signals of the run that just happened. |
| `NGSTI-E03` | Log ingestion in a project that declares no log paths. | The command refuses, naming the declaration to add. | There is no log to scan. |
| `NGSTI-E04` | A report that is missing or cannot be parsed. | The command fails naming which report format could not be read. | A signal read from a broken report would be invented. |
| `NGSTI-E05` | The map cannot be loaded while ingesting reports. | The command fails with a hint to build the map. | The signals are written onto the map's nodes. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/testsig/junit.go` | `ParseJUnit` | infra — reads the execution report |
| DEP2 | `internal/testsig/lcov.go` | `ParseLCOV` | infra — reads the line coverage report |
| DEP3 | `internal/testsig/mutation.go` | `ParseMutationFormat` | infra — reads the mutation report in the declared format |
| DEP4 | `internal/mapx/ingest.go` | `IngestExecutionSuite`, `IngestCoverageSuite`, `IngestMutationScoped`, `ResolveReportPaths`, `DropExternalSuites` | mapa — writes the signals onto the nodes |
| DEP5 | `internal/logscan/scan.go` | `Scan`, `SpecFailureCodes` | apoio — finds the failure codes in the logs and in the specs |
| DEP6 | `cmd/anchors/common/spec.go` | `CodesInFileOfUnit` | comando — the scenario codes a spec declares |
| DEP7 | `internal/config/config.go` | `Load` | config — the project's letters, code length, suites and log paths |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

<!-- @anchors
  code: HLDCH
  updated_at: 2026-09-26
  layer: comando
-->
# DoctorCommand — the global health x-ray, and the repair of the github-mode environment

> **Code**: `HLDCH`

## Overview

The check confronts nodes against gates one target at a time. The doctor is the global view: it sweeps the
map, the configuration and the disk for the systemic loose ends no local gate sees (dead edges, phantom
nodes, missing identities, loose layers, holes in gate coverage) and prints them grouped by the check
that found them. It is diagnosis on demand: it presents what it found and never blocks.

With the fix flag it prepares the environment of the github workflow mode, and only of that mode. It
checks the `gh` credential before anything else, so a missing login is named once instead of surfacing as
five unrelated failures; warns (without deleting) about a local task queue or delivery records that must
not exist in that mode; seeds the workflow pipelines; protects the declared branches with a body that
requires a pull request; disables the approval requirement when the author could never satisfy it; and
ensures the state labels. It never creates the board, which is structure shared by the team and outside
the repository. Running it again on a repository already set up changes nothing.

The pipelines flag answers one question for CI: are the pipelines in place and current? It names the
missing and outdated ones, and fails only when the project opted into blocking on a stale pipeline.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration and the map | the project's configuration and map | a project with neither | this unit: it refuses (`HLDCH-E01`, `HLDCH-E02`) |
| the workflow mode | local, manual or github | — | the configuration package: the mode is validated when the configuration loads |
| the `gh` tool | installed and authenticated, in github mode | absent, or not logged in | this unit: the fix refuses before any repair (`HLDCH-E03`, `HLDCH-E04`) |

## Effects

| Effect | Description |
| --- | --- |
| `HLDCH-B01` | Prints a header with the number of nodes, edges and layers, then the findings grouped by the check that raised them, each group with its count, then a summary of attention points and findings and the note that nothing was blocked. |
| `HLDCH-B02` | A group is marked as a warning when any of its findings is a warning, even when an informational one came first; each warning finding is marked on its own line. |
| `HLDCH-B03` | With no finding, it prints that no systemic loose end was found. |
| `HLDCH-B04` | The fix outside the github mode says there is nothing to do and seeds nothing. |
| `HLDCH-B05` | The fix in github mode seeds the missing workflow pipelines and names them. |
| `HLDCH-B06` | The fix protects each declared branch by sending the protection body to the API, reporting a branch that does not exist yet and going on to the next. |
| `HLDCH-B07` | The fix disables the approval requirement when the protection cannot be bypassed by the current account, and says so. |
| `HLDCH-B08` | The fix ensures the state labels in the repository, the project's own label among them, and closes saying the board is optional. |
| `HLDCH-B09` | The protection body carries every field the API requires, with the pull-request review object holding the number of approvals passed. |
| `HLDCH-B10` | The pipelines check in local mode says there is no workflow pipeline to check. |
| `HLDCH-B11` | The pipelines check names each missing or outdated pipeline and, by default, warns and lets CI continue; when every pipeline is current it says so. |
| `HLDCH-B12` | The pipelines check fails, counting the missing and outdated pipelines, when the project declared that a stale pipeline blocks. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `HLDCH-I01` | The diagnosis never fails the command, whatever it finds. | runs the doctor on a map with warning findings and gets no error |
| `HLDCH-I02` | The fix is idempotent: a second run on a repository it already set up reports the pipelines up to date and creates none. | runs the fix twice and reads the second run's output |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `HLDCH-X01` | The fix never creates the board; it only says the board is optional. | A board is structure shared by the team and lives outside the repository: one created by mistake pollutes the organization and is not undone by a checkout. |
| `HLDCH-X02` | The fix warns about a local task queue and delivery records found in github mode and never deletes them. | The queue may hold claimed work in progress and the records are the memory of what was delivered; deleting them unseen destroys what nobody else has. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `HLDCH-E01` | The project has no configuration. | Error naming the configuration load. | The checks and the mode come from the configuration. |
| `HLDCH-E02` | The map does not exist or cannot be read. | Error naming the map load and pointing at `anchors map build`. | The diagnosis is a sweep over the map. |
| `HLDCH-E03` | The fix runs in github mode and `gh` is not installed. | Refuses naming the tool, before seeding or changing anything. | Every repair goes through `gh`; one refusal naming the cause beats five failures naming symptoms. |
| `HLDCH-E04` | The fix runs in github mode and `gh` is installed but not authenticated. | Refuses with the instruction to log in, before seeding or changing anything. | Logging in is interactive and an agent cannot complete it; telling the person is the only honest way out. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/health/health.go` | `Diagnose` | infra — the systemic findings |
| DEP2 | `internal/initx/workflows.go` | `MissingWorkflow`, `OutdatedWorkflows`, `SemeiaWorkflows` | infra — the workflow pipelines |
| DEP3 | `internal/config/config.go` | `Load` | config — the mode, the repository and the protected branches |
| DEP4 | `internal/mapx/store.go` | `Load` | mapa — the map the diagnosis sweeps |
| DEP5 | `internal/change/change.go` | `Pending` | infra — the delivery records |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

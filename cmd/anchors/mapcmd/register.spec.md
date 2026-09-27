<!-- @anchors
  code: MPRGM
  updated_at: 2026-09-26
  layer: comando
-->
# MapRegister — hangs the map domain's commands on the root command

> **Code**: `MPRGM`

## Overview

The command line is assembled by domain: each domain package registers its own commands on
the root, and the root knows only the domains. This unit is the map domain's registration.
It hangs the eight commands that read or write the map on the root: building and showing the
map, the impact query, the ingestion of test and log signals, the judgment record, the code
rename, the revision renumbering, the work flows and the review of observed failures.

The registration holds no logic of its own; what it guarantees is that each of those
commands is reachable by its name from the root, and that the map domain adds nothing else.

## Effects

| Effect | Description |
| --- | --- |
| `MPRGM-B01` | After registration, each of map, impact, ingest, judge, recode, renumber, flow and failures is reachable from the root by its name. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the root command | the program's root command, before it runs | a command that is not the root | the program entry point: it registers every domain on its own root |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPRGM-X01` | The map domain registers no command outside its own eight. | Another domain's command registered here would be registered twice, or owned by the wrong package. |

## Errors

none — registration adds commands to the root and handles no failure; a command's own failures are its spec's

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/mapcmd/map.go` | `newMapCmd` | comando — the map command |
| DEP2 | `cmd/anchors/mapcmd/impact.go` | `newImpactCmd` | comando — the impact command |
| DEP3 | `cmd/anchors/mapcmd/ingest.go` | `newIngestCmd` | comando — the ingest command |
| DEP4 | `cmd/anchors/mapcmd/judge.go` | `newJudgeCmd` | comando — the judge command |
| DEP5 | `cmd/anchors/mapcmd/recode.go` | `newRecodeCmd` | comando — the recode command |
| DEP6 | `cmd/anchors/mapcmd/renumber.go` | `newRenumberCmd` | comando — the renumber command |
| DEP7 | `cmd/anchors/mapcmd/flow.go` | `newFlowCmd` | comando — the flow command |
| DEP8 | `cmd/anchors/mapcmd/failures.go` | `newFailuresCmd` | comando — the failures command |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

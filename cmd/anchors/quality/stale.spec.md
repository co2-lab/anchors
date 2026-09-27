<!-- @anchors
  code: STEDS
  updated_at: 2026-09-26
  layer: comando
-->
# StaleEdges — lists the confrontation debt: expired test evidence and stale edges

> **Code**: `STEDS`

## Overview

Once the check stamps the map's edges, "validated" and "stale" mean something. An edge is stale when it
was never confronted, or when one of its ends moved to a new revision after the last stamp. A test's
evidence expires when the test passed against code that has since changed: the test file itself, or one of
the files it depends on.

The stale command reads that state from the map and says what has to be confronted again. It lists the
expired test evidence first, because it is the most actionable of the two and the one that hides itself:
the old report keeps saying the test is green. Then it lists the stale edges, one per line, each labelled
as never validated or as drifted, and closes with how many of each there are. When no edge is stale it says
so with the number of edges it looked at.

It only reads. Reconciling the debt is the check's job, and the command says so in its help.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | the project's map at its default path, or the one given by the map flag | no map, or a map that cannot be read | this unit: it refuses and points at the map build (`STEDS-E01`) |
| the project root | a directory, resolved to an absolute path | a root that cannot be resolved | the configuration package: resolving the root returns the error, which this unit returns |

## Effects

| Effect | Description |
| --- | --- |
| `STEDS-B01` | The expired test evidence is listed before the stale edges, with a heading counting it and a closing note that expired evidence is not a test defect. |
| `STEDS-B02` | Each expired evidence names its reason: the test file itself changed, or how many dependencies changed naming the first one, adding that the test itself changed too when both happened. |
| `STEDS-B03` | Each stale edge is printed as its origin, type and destination with its reason, never validated when it has no stamp and advanced rev otherwise, followed by a summary counting each reason. |
| `STEDS-B04` | When no edge is stale, it prints the clean message carrying the total number of edges. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `STEDS-I01` | An edge stamped at the current revisions of both of its ends is never listed as stale. | a map with a drifted edge, a never-validated edge and an edge stamped at the current revisions lists only the first two |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `STEDS-X01` | Reads the map and never writes it: no edge is stamped and no evidence is refreshed. | Reconfronting is the check's job; a listing that re-stamped what it lists would erase the debt it exists to show. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `STEDS-E01` | The map does not exist or cannot be read. | Error pointing at `anchors map build`. | Without a map there are no edges and no stamps; an empty listing would pass for a clean project. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/store.go` | `Load` | mapa — reading the map |
| DEP5 | `internal/mapx/stamp.go` | `StaleEdges` | mapa — the edges whose stamp is missing or behind |
| DEP2 | `internal/mapx/evidence.go` | `EvidenceStaleFor` | mapa — the test evidence against its closure |
| DEP3 | `internal/config/config.go` | `AbsRoot` | config — the project root |
| DEP4 | `internal/i18n/i18n.go` | `T` | apoio — the localized lines |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none

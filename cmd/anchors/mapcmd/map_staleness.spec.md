<!-- @anchors
  code: MPSTM
  updated_at: 2026-09-26
  layer: comando
-->
# MapStaleness — names the files of the map whose content changed after the map was built

> **Code**: `MPSTM`

## Overview

Presence and age are two different questions about the map. The pre-commit already refuses a governed
file that is not in the map at all — the new file of someone who never rebuilt it. What it did not see is
the file that IS in the map but carries the revision of an older version: the map was built, and the
work went on after it. In the reference project, seven of twelve rejections of the gates pipeline were
exactly this, and in all seven the map had been committed; it had simply been built too early.

This unit answers the second question. It walks the project as the map build does and returns every
node whose current content revision differs from the revision the map recorded. The comparison is
content hash against content hash: a date moves on a checkout without the content moving, and stays
put on an edit that preserves the modification time, so a date would give the wrong answer both ways.

It is an auxiliary check consumed by `check`, which prints the result as a warning. For that reason it
never fails: when the project cannot be walked, it has nothing to compare and says nothing, rather than
taking down the report the user came for.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | a loaded map, or no map at all | — | this unit: no map answers nothing stale |
| the project configuration | a loaded configuration | an absent configuration | the caller: `check` loads the configuration before asking |
| the project root | a directory that can be walked | a root that cannot be walked | this unit: an unwalkable root answers nothing stale |

## Effects

| Effect | Description |
| --- | --- |
| `MPSTM-B01` | `StaleMapNodes` names a file of the map edited after the map was built as stale; a map freshly built names nothing. |
| `MPSTM-B02` | A node whose file no longer exists on disk is not named: a missing file is a different problem from an old one. |
| `MPSTM-B03` | Staleness is decided by the content revision, never by the modification time: a touched file with the same content is not stale, and a changed file with an old modification time is. |
| `MPSTM-B04` | A node whose recorded revision is empty is not named: an empty revision says nothing to compare. |
| `MPSTM-B05` | With no map, nothing is named. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MPSTM-I01` | Rebuilding the map from the current tree makes the answer empty: stale means the map is behind the files, nothing else. | builds the map, confirms nothing is stale, edits a file, confirms it is named |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPSTM-X01` | Does not rebuild or write the map; it only compares. | Rebuilding is the user's action (`map build`); a check that silently rebuilt would hide the drift it exists to show. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPSTM-E01` | The project root cannot be walked. | Nothing is named and no error is raised. | The check is auxiliary: failing here would take down the `check` report that prints it. <!-- @resilient: the staleness notice is a warning beside the check, never its verdict; a root that cannot be walked only loses the warning --> |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/scan.go` | `Walk` | scan — the current content revision of each file |
| DEP2 | `internal/mapx/model.go` | `Graph` | mapa — the revision the map recorded |
| DEP3 | `internal/config/config.go` | `Config` | config — the layers that decide what is walked |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
